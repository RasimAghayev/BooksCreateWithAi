# Chapter 3 — Преобразование данных и композиция (Data transformasiyası və kompozisiya)

## Bu chapter nədən bəhs edir?

Tip çevirmələri və interface type-assertion, math/math/big böyük ədədlər,
valyuta (float64 problemləri), SQL Null tipləri, gob/base64 encoding,
struct tag-lər + reflection, closure-larla funksional kolleksiyalar.

## Əsas fikirlər

### 1. Tip çevirmələri və interface type-assertion
**Nədir:** Go-da tip çevirmə (conversion) və interface-dən konkret tipə
keçid (type assertion / type switch).

**Kitabdan kod nümunəsi:**
```go
// Numerik çevirmə:
var a = 24        // int
var b = 2.0       // float64
c := float64(a) * b          // int → float64 çevrilməsi MÜTLƏQ
precision := fmt.Sprintf("%.2f", b)   // string formatlama

// strconv ilə string ↔ ədəd:
res, err := strconv.ParseInt("1234", 10, 64)  // base 10, 64 bit
res, err = strconv.ParseInt("FF", 16, 64)     // hex!
val, err := strconv.ParseBool("true")

// Type switch — interface tipinə görə:
func CheckType(s interface{}) {
    switch s.(type) {
    case string:
        fmt.Println("It's a string!")
    case int:
        fmt.Println("It's an int!")
    default:
        fmt.Println("not sure what it is...")
    }
}

// Tək tip yoxlaması (comma-ok idiomu):
var i interface{} = "test"
if val, ok := i.(string); ok {
    fmt.Println("val is", val)
}
```

**Sub-kod izahı:**
- `float64(a)` → müxtəlif numerik tiplər avtomatik çevrilmir — açıq cast lazımdır
- `strconv.ParseInt(s, base, bitSize)` → base: 10, 16 (hex) və s.
- `i.(string)` assertion uğursuzsa `ok=false` — panic YOX (comma-ok)
- Type switch reflection-in yüngül alternatividir

**Vacib:** Sırf conversion xətaları compile-time yaxalanır; type assertion
xətaları RUNTIME-da baş verir — comma-ok istifadə et!

### 2. math və math/big
**Nədir:** `math` — float64 üzərində qabaqcıl riyaziyyat; `math/big` — int64-ə
sığmayan dəyərlər üçün dəyişən uzunluqlu ədədlər.

**Kitabdan kod nümunəsi:**
```go
// math paketi:
result := math.Sqrt(float64(25))   // 5 — int DƏYİL, float64 qaytarır
result = math.Ceil(9.5)            // 10 — yuxarı yuvarlaqlaşdır
result = math.Floor(9.5)           // 9  — aşağı
fmt.Println(math.Pi, math.E)       // konstantlar

// math/big + memoization ilə Fibonaççi:
var memoize map[int]*big.Int
func init() {
    memoize = make(map[int]*big.Int)   // init() avtomatik çağrılır
}
func Fib(n int) *big.Int {
    if n < 0 {
        return big.NewInt(1)
    }
    if n < 2 {
        memoize[n] = big.NewInt(1)    // base case
    }
    if val, ok := memoize[n]; ok {
        return val                     // hesablanmışsa bir daha hesablama
    }
    memoize[n] = big.NewInt(0)
    memoize[n].Add(memoize[n], Fib(n-1))   // Add dəyişənin özünə yazır
    memoize[n].Add(memoize[n], Fib(n-2))
    return memoize[n]
}
```

**Sub-kod izahı:**
- `big.Int` immutable DEYİL — `Add(a, b)` nəticəni receiver-ə yazır
- `init()` funksiyası → package yüklənəndə avtomatik icra olunur
- Memoization → rekursiv Fib-in eksponensial təkrar hesabını map ilə kəsir
- int64 Fib-də ~90-cu addımda overflow edir — big.Int məcburidir

### 3. Valyuta və float64 problemi
**Nədir:** Pul float64 kimi saxlanırsa yuvarlaqlaşdırma xətaları YARANIR —
pulu SENTLƏR kimi int64-də saxlamaq düzgün haldır.

**Necə işləyir:** "15.93" stringini float64-ə keçmədən birbaşa 1593 sentə
çevirir — string parçalanması ilə.

**Kitabdan kod nümunəsi:**
```go
func ConvertStringDollarsToPennies(amount string) (int64, error) {
    _, err := strconv.ParseFloat(amount, 64)   // yalnız VALIDASİYA üçün
    if err != nil {
        return 0, err
    }
    groups := strings.Split(amount, ".")       // "15.93" → ["15","93"]
    result := groups[0]
    r := ""
    if len(groups) == 2 {
        if len(groups[1]) != 2 {
            return 0, errors.New("invalid cents")
        }
        r = groups[1]
    }
    for len(r) < 2 {          // "5" → "50", "" → "00"
        r += "0"
    }
    result += r               // "15" + "93" = "1593"
    return strconv.ParseInt(result, 10, 64)
}

func ConvertPenniesToDollarString(amount int64) string {
    result := strconv.FormatInt(amount, 10)
    negative := false
    if result[0] == '-' { result = result[1:]; negative = true }
    for len(result) < 3 { result = "0" + result }   // 5 → "005"
    length := len(result)
    result = result[0:length-2] + "." + result[length-2:]  // "1608" → "16.08"
    if negative { result = "-" + result }
    return result
}
```

**Sub-kod izahı:**
- Float-a heç vaxt keçmir — yalnız ParseFloat ilə validlik yoxlanılır
- Sentlər həmişə 2 rəqəmli olur; mənfi dəyər ayrıca idarə olunur
- Aritmetika int64 üzərində — dəqiqdir (pennies += 15 → 16.08)

### 4. Pointer-lər və SQL NullTypes
**Nədir:** JSON-da 0 və "" default dəyərləri "yoxlanmayan" dəyərlərdən seçilə
bilmir; `omitempty` 0-lığı da silir. Həll: pointer və ya sql.NullInt64.

**Kitabdan kod nümunəsi:**
```go
// Problem: age=0 və age yoxdur eyni görünür:
type Example struct {
    Age int `json:"age,omitempty"`   // 0 YOX SAYILIR!
}

// Həll 1 — pointer:
type ExamplePointer struct {
    Age *int `json:"age,omitempty"`  // nil = qeyd olunmayıb, &0 = sıfır dəyər
    Name string
}

// Həll 2 — sql.NullInt64 wrapper:
type nullInt64 sql.NullInt64
type ExampleNullInt struct {
    Age *nullInt64 `json:"age,omitempty"`
    Name string
}
func (v *nullInt64) MarshalJSON() ([]byte, error) {
    if v.Valid {
        return json.Marshal(v.Int64)
    }
    return json.Marshal(nil)        // invalid → JSON null
}
func (v *nullInt64) UnmarshalJSON(b []byte) error {
    v.Valid = false
    if b != nil {
        v.Valid = true
        return json.Unmarshal(b, &v.Int64)
    }
    return nil
}
```

**Sub-kod izahı:**
- `*int` → nil (yoxdur) ilə 0 (sıfır) fərqlənir; təyinat: `a := 0; e.Age = &a`
- `sql.NullInt64{Int64, Valid}` → SQL NULL-in Go modeli
- MarshalJSON/UnmarshalJSON realizə edərək hər hansı tipi JSON-ə uyğunlaşdır
- Pointer saxlanır ki, `omitempty` nil halında düzgün işləsin

### 5. gob və base64 encoding
**Nədir:** gob — Go tipləri üçün binary stream formatı; base64 — binarı
string/URL-yə uyğun formaya çevirir.

**Kitabdan kod nümunəsi:**
```go
// gob:
func GobExample() error {
    buffer := bytes.Buffer{}
    p := pos{X: 10, Y: 15, Object: "wrench"}
    e := gob.NewEncoder(&buffer)
    if err := e.Encode(&p); err != nil {   // struct → binary
        return err
    }
    p2 := pos{}
    d := gob.NewDecoder(&buffer)
    if err := d.Decode(&p2); err != nil {  // binary → struct
        return err
    }
    return nil
}

// base64:
value := base64.URLEncoding.EncodeToString([]byte("encoding some data!"))
decoded, err := base64.URLEncoding.DecodeString(value)

// base64 stream (encoder/decoder):
encoder := base64.NewEncoder(base64.StdEncoding, &buffer)
encoder.Write([]byte("data"))
encoder.Close()                            // Close MÜTLƏQ!
decoder := base64.NewDecoder(base64.StdEncoding, &buffer)
results, err := ioutil.ReadAll(decoder)
```

**Sub-kod izahı:**
- gob → Go-dan Go-yə RPC/wire protokolu üçün; interface istifadəsində
  `gob.Register` lazımdır; ardıcıl çox elementdə ən effektivdir
- base64 → GET URL-ləri, JSON payload-ları; StdEncoding vs URLEncoding
  (- _ simvolları URL-də təhlükəsizdir)
- Tək element üçün JSON daha portable; gob Go-ya spesifikdir

### 6. Struct tag-lər və reflection
**Nədir:** Struct tag-lər key-value string-ləridir; reflect paketi interface
içindəki tipi, sahələri və tag-ləri dindirir.

**Necə işləyir:** `reflect.TypeOf/ValueOf` → `NumField()` → `Field(i)` →
`Tag.Lookup("serialize")`.

**Kitabdan kod nümunəsi:**
```go
type Person struct {
    Name string `serialize:"name"`
    City string `serialize:"city"`
    State string                    // tag YOXDUR → sahə adı istifadə olunur
    Misc string `serialize:"-"`     // ignore
    Year int    `serialize:"year"`  // int olduğundan YƏNƏ ignore (yalnız string)
}

func SerializeStructStrings(s interface{}) (string, error) {
    r := reflect.TypeOf(s)
    value := reflect.ValueOf(s)
    if r.Kind() == reflect.Ptr {       // pointer-i dereference et
        r = r.Elem()
        value = value.Elem()
    }
    result := ""
    for i := 0; i < NumField sınırı; i++ {
        field := r.Field(i)
        key := field.Name
        if serialize, ok := field.Tag.Lookup("serialize"); ok {
            if serialize == "-" { continue }
            key = serialize
        }
        if value.Field(i).Kind() == reflect.String {
            result += key + ":" + value.Field(i).String() + ";"
        }
    }
    return result, nil
}

// Deserialize: pointer TƏLƏB OLUNUR (dəyişmək üçün):
if r.Kind() != reflect.Ptr {
    return errors.New("res must be a pointer")
}
value.Field(i).SetString(val)     // sahəyə yazmaq
```

**Sub-kod izahı:**
- `field.Tag.Lookup("serialize")` → tag varsa ok=true
- `value.Field(i).SetString(...)` → yalnız pointer üzərindən mümkün
- Qaydalar: yalnız string sahələr; "-" ignore; tag yoxdursa field adı
- RESTRIKSİYA: reflection export olunmayan (kiçik hərf) sahələrlə tam işləmir

### 7. Closure-larla kolleksiyalar (map/filter)
**Nədir:** Generics olmadan funksional Map/Filter — closure-lar və funksiya
tipləri ilə.

**Kitabdan kod nümunəsi:**
```go
type WorkWith struct {
    Data    string
    Version int
}

func Filter(ws []WorkWith, f func(w WorkWith) bool) []WorkWith {
    result := make([]WorkWith, 0)      // len 0 — böyümək üçün
    for _, w := range ws {
        if f(w) {
            result = append(result, w)
        }
    }
    return result
}

func Map(ws []WorkWith, f func(w WorkWith) WorkWith) []WorkWith {
    result := make([]WorkWith, len(ws))   // eyni uzunluq
    for pos, w := range ws {
        result[pos] = f(w)
    }
    return result
}

// Dəyişən funksiyalar:
func LowerCaseData(w WorkWith) WorkWith { w.Data = strings.ToLower(w.Data); return w }
func IncrementVersion(w WorkWith) WorkWith { w.Version++; return w }

// Closure generator — arqumenti "yadda saxlayan" funksiya:
func OldVersion(v int) func(w WorkWith) bool {
    return func(w WorkWith) bool {
        return w.Version >= v
    }
}

// İstifadə (zəncirləmə):
ws = collections.Map(ws, collections.LowerCaseData)
ws = collections.Map(ws, collections.IncrementVersion)
ws = collections.Filter(ws, collections.OldVersion(3))
```

**Sub-kod izahı:**
- `Filter` → predikat true olanları saxlayır; `Map` → hamısını transform edir
- `OldVersion(3)` → closure: v=3 dəyəri tutulmuş predikat qaytarır
- Funksiyalar safdır (pure) — error yoxdur, yan təsir yoxdur; test asandır

## Əsas terminlər

- Type Conversion / Type Assertion (tip çevirməsi/təsdiqi)
- Type Switch (tip seçimi)
- comma-ok idiom
- big.Int (dəyişən uzunluqlu tam ədəd)
- Memoization (nəticələrin yadda saxlanması)
- Pointer Semantics (göstərici mənası)
- SQL NullTypes (sql.NullInt64 və s.)
- gob Encoding (Go binary format)
- Base64 (StdEncoding / URLEncoding)
- Struct Tag (struktur etiketi)
- Reflection (reflect paketi)
- Closure (qapanış)
- Map / Filter (funksional kolleksiya əməliyyatları)

## Praktik nəticə

- Pul float64 ilə YOX — string→sent (int64) çevirməsi pattern
- 0 vs "yoxdur" probleminə pointer və ya NullType + custom MarshalJSON
- Böyük ədədlər üçün math/big; Add metodları receiver-ə yazır
- Öz seriya formatını yazmaq üçün: reflect + struct tag-lər kifayətdir
- Generics olmadan map/filter — closure generatorləri ilə zərif həll
- gob yalnız Go↔Go üçün; inter-platform lazımdırsa JSON/base64

## Mənbə

Pages: 93-132 (Chapter 3, Go Programming Cookbook 2nd ed)
