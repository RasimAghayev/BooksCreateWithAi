# Chapter 11 — Reflection and code generation (Technique 66-70)

## Bu chapter nədən bəhs edir?

Reflect paketi: value/type/kind üçlüyü, type switch vs kind switch, interfeys implementasiya yoxlaması (nil pointer trick), struct sahələrinin rekursiv gəzilməsi, custom struct tag-lər (INI marshal/unmarshal), go generate + text/template ilə kod generasiyası (typed queue generator).

## Əsas fikirlər

### 1. Reflection üçlüyü: Value, Type, Kind
| Anlayış | Təsvir | reflect |
|---|---|---|
| **Value** | dəyişənin işarə etdiyi DATA (`x := 5` → 5; `myFunc := strings.Split` → funksiya) | `reflect.Value` |
| **Type** | dəyərin tipi (`bytes.Buffer`) | `reflect.Type` (interfeys) |
| **Kind** | PRİMİTİV növ (struct, ptr, int, float64, string, slice, func...) | `reflect.Kind` |

Reflection adətən value ilə başlayır → type və kind kəşf olunur.

### TECHNIQUE 66: Type switch vs Kind switch
**Type switch (type-a bağlı):**
```go
func sum(v ...interface{}) float64 {
    var res float64 = 0
    for _, val := range v {
        switch val.(type) {                 // TİP üzərində
        case int:
            res += float64(val.(int))
        case int64:
            res += float64(val.(int64))
        case uint8:
            res += float64(val.(uint8))
        case string:
            a, err := strconv.ParseFloat(val.(string), 64)
            if err != nil {
                panic(err)
            }
            res += a
        default:
            fmt.Printf("Unsupported type %T. Ignoring.\n", val)
        }
    }
    return res
}
```
**MƏHDUDİYYƏT:** `type MyInt int64` — MyInt int64 DEYİL (fərqli TİP)! → default-a düşür:
```
Unsupported type main.MyInt. Ignoring.
```
Hər case 1 tip (case 1, 2 kimi birləşdirmə YOX) — 10 integer tipi = 10 case.

**Kind switch (reflect ilə — həll):**
```go
func sum(v ...interface{}) float64 {
    var res float64 = 0
    for _, val := range v {
        ref := reflect.ValueOf(val)          // interface{} → reflect.Value
        switch ref.Kind() {                  // KİND üzərində
        case reflect.Int, reflect.Int64:     // birləşdirmə MÜMKÜN!
            res += float64(ref.Int())        // Int() — ən böyük formaya çevirir
        case reflect.Uint8:
            res += float64(ref.Uint())
        case reflect.String:
            a, err := strconv.ParseFloat(ref.String(), 64)
            if err != nil {
                panic(err)
            }
            res += a
        default:
            fmt.Printf("Unsupported type %T. Ignoring.\n", val)
        }
    }
    return res
}
```
- `reflect.ValueOf(val)` → Value; `.Kind()` → primitiv
- **`ref.Int()` / `ref.Uint()`** — int8/int16/... hamısını int64-ə, uint-ləri uint64-ə çevirir → 10 case 2-yə düşür!
- MyInt-in kind-i int64-dür → tutulur ✓

**Type vs Kind seçim:** numeric emal → kind; xüsusi tiplər (error handling) → type switch.

### TECHNIQUE 67: İnterfeys implementasiya yoxlaması
**Üsul 1 — type assertion (dəyər üçün):**
```go
func isStringer(v interface{}) bool {
    _, ok := v.(fmt.Stringer)     // dönüşdür + yoxla bir addımda
    return ok
}
```

**Üsul 2 — reflect (runtime-da interfeys SEÇİMİ üçün):**
```go
func implements(concrete interface{}, target interface{}) bool {
    iface := reflect.TypeOf(target).Elem()    // nil pointer → interfeys TİPİ
    v := reflect.ValueOf(concrete)
    t := v.Type()
    if t.Implements(iface) {
        fmt.Printf("%T is a %s\n", concrete, iface.Name())
        return true
    }
    fmt.Printf("%T is not a %s\n", concrete, iface.Name())
    return false
}

func main() {
    n := &Name{First: "Inigo", Last: "Montoya"}
    stringer := (*fmt.Stringer)(nil)     // İRADİ nil pointer — TİPİ saxlama triki!
    implements(n, stringer)               // *main.Name is a Stringer
    writer := (*io.Writer)(nil)
    implements(n, writer)                // *main.Name is not a Writer
}
```

**nil pointer trick (vacib!):** reflect-də birbaşa interfeysə reflect etmək MÜMKÜN DEYİL — `(*fmt.Stringer)(nil)` yalnız TİPİ daşıyan yer tutucu yaradır; `.Elem()` interfeys tipinə çıxarır. `Type.Implements(ifaceType)` — yoxlamanın özü.

### TECHNIQUE 68: Struct sahələrinin gəzilməsi
**Kitabdan kod nümunəsi — rekursiv walk:**
```go
type Person struct {
    Name    *Name
    Address *Address
}
type Name struct {
    Title, First, Last string
}

func walk(u interface{}, depth int) {
    val := reflect.Indirect(reflect.ValueOf(u))    // pointer-i İZLƏ (deref)
    t := val.Type()
    tabs := strings.Repeat("\t", depth+1)

    fmt.Printf("%sValue is type %q (%s)\n", tabs, t, val.Kind())

    if val.Kind() == reflect.Struct {
        for i := 0; i < t.NumField(); i++ {        // sahə sayı
            field := t.Field(i)                     // reflect.StructField (TİPDƏN)
            fieldVal := reflect.Indirect(val.Field(i))  // dəyər (VALUEDAN, eyni index!)
            ...
            if fieldVal.Kind() == reflect.Struct {
                walk(fieldVal.Interface(), depth+1)     // REKURSİV
            }
        }
    }
}

// Output:
// Field "Name" is type "*main.Name" (struct)
//   Value is type "main.Name" (struct)      ← Indirect pointer-i açdı
//     Field "Title" is type "string" (string)
```

**Ayrı-seçkilik (kitabın xülasəsi):**
- `reflect.ValueOf(i)` → interface{}-dən Value
- `reflect.Indirect(v)` → pointer-i izlə (pointer deyilsə eyni qaytarır — TƏHLÜKƏSİZ)
- Value-dan type/kind birbaşa alınır
- Struct üçün: **TİPDƏN** `NumField()` + `Field(i)` (StructField: ad, tip, tag); **VALUEDAN** `Field(i)` (dəyər) — eyni i!

**Diqqət:** reflect-in çox funksiyası ERROR YOX — PANIC edir.

### 2. Struct tag-lər əsasları
**Annotations = backtick içində istənilən string** — kompilyasiyada ROLU YOX, runtime-da reflect ilə oxunur:
```go
type Name struct {
    First string `json:"firstName" xml:"FirstName"`     // hər encode özünükini tapır
    Last  string `json:"lastName,omitempty"`            // ",omitempty" → data
    Other string `not,even.a=tag`                       // LEGAL — sadəcə mənasız
}
```
**Tag konvensiyası (de-facto standart):** `NAME:"VALUE,DATA"` — `json:"myField,omitempty"`, `xml:"href,attr"`.
**İgnore idiomu:** `json:"-"` (dash) — encoder sahəni ötürsün.
**Validation istifadəsi:** `validate:"^[a-z]+$"` (Deis Router nümunəsi).

### TECHNIQUE 69: Custom tag-lər — INI marshal/unmarshal
**Tapşırıq:** `total=247 / running=2` formatı ↔ struct.

**Struct + tag:**
```go
type Processes struct {
    Total    int     `ini:"total"`
    Running  int     `ini:"running"`
    Sleeping int     `ini:"sleeping"`
    Threads  int     `ini:"threads"`
    Load     float64 `ini:"load"`
}
```

**Tag oxunu:**
```go
func fieldName(field reflect.StructField) string {
    if t := field.Tag.Get("ini"); t != "" {    // GO AVTOMATİK TAG PARSE EDİR!
        return t
    }
    return field.Name                            // tag yoxdursa sahə adı
}
```
`StructField.Tag.Get("ini")` — `NAME:"VALUE"` formatını avtomatik pars edir.

**Marshal:**
```go
func Marshal(v interface{}) ([]byte, error) {
    var b bytes.Buffer
    val := reflect.Indirect(reflect.ValueOf(v))   // deref
    if val.Kind() != reflect.Struct {
        return []byte{}, errors.New("unmarshal can only take structs")
    }
    t := val.Type()
    for i := 0; i < t.NumField(); i++ {
        f := t.Field(i)                            // TİPDƏN sahə
        name := fieldName(f)                       // tag-dən ad
        raw := val.Field(i).Interface()            // VALUEDAN dəyər
        fmt.Fprintf(&b, "%s=%v\n", name, raw)      // Fprintf çevirməni ÖZÜ edir
    }
    return b.Bytes(), nil
}
```

**Unmarshal:**
```go
func Unmarshal(data []byte, v interface{}) error {
    val := reflect.Indirect(reflect.ValueOf(v))
    t := val.Type()
    b := bytes.NewBuffer(data)
    scanner := bufio.NewScanner(b)                  // sətir-sətir
    for scanner.Scan() {
        line := scanner.Text()
        pair := strings.SplitN(line, "=", 2)         // ad=dəyər
        if len(pair) < 2 {
            continue                                 // xətalı sətir — ötür
        }
        setField(pair[0], pair[1], t, val)
    }
    return nil
}
```

**setField — kind switch + SET:**
```go
func setField(name, value string, t reflect.Type, v reflect.Value) {
    for i := 0; i < t.NumField(); i++ {
        field := t.Field(i)
        if name == fieldName(field) {                // uyğun sahə tapıldı
            var dest reflect.Value
            switch field.Type.Kind() {              // KİND-ə görə çevir
            case reflect.Int:
                ival, err := strconv.Atoi(value)
                if err != nil { ...; continue }
                dest = reflect.ValueOf(ival)
            case reflect.Float64:
                fval, err := strconv.ParseFloat(value, 64)
                ...
                dest = reflect.ValueOf(fval)
            case reflect.String:
                dest = reflect.ValueOf(value)
            case reflect.Bool:
                bval, err := strconv.ParseBool(value)
                ...
                dest = reflect.ValueOf(bval)
            default:
                continue                             // dəstəklənməyən — skip (error DEYİL)
            }
            v.Field(i).Set(dest)                     // SAHƏNİ SET ET!
        }
    }
}
```
- `v.Field(i).Set(reflect.ValueOf(...))` — runtime-da struct-a yazmaq
- Dəstəklənməyən kind → sadəcə ötür; genişləndirmə asan amma təkrarlanan

**Dərs:** reflect ilə tip çevirmələri çox BOILERPLATE — bəzən generasiya daha yaxşıdır.

### TECHNIQUE 70: go generate + template — kod yaradan kod
**go generate mexanizmi:** faylların başındakı xüsusi şərh icra edir:
```go
//go:generate COMMAND [ARGUMENT...]

//go:generate echo hello       → $ go generate simple.go → "hello" çap
```

**Problem:** generics yox (kitab dövrü) — type-safe queue hər tip üçün əl iləmi?

**Template (listing 11.19):**
```go
var tpl = `package {{.Package}}
type {{.MyType}}Queue struct {
    q []{{.MyType}}
}
func New{{.MyType}}Queue() *{{.MyType}}Queue {
    return &{{.MyType}}Queue{
        q: []{{.MyType}}{},
    }
}
func (o *{{.MyType}}Queue) Insert(v {{.MyType}}) {
    o.q = append(o.q, v)
}
func (o *{{.MyType}}Queue) Remove() {{.MyType}} {
    if len(o.q) == 0 {
        panic("Oops.")
    }
    first := o.q[0]
    o.q = o.q[1:]
    return first
}
`

func main() {
    tt := template.Must(template.New("queue").Parse(tpl))
    for i := 1; i < len(os.Args); i++ {                       // hər arqument = 1 tip
        dest := strings.ToLower(os.Args[i]) + "_queue.go"      // myint_queue.go
        file, err := os.Create(dest)
        if err != nil { ...; continue }
        vals := map[string]string{
            "MyType":  os.Args[i],
            "Package": os.Getenv("GOPACKAGE"),                 // go generate ENV qoyur!
        }
        tt.Execute(file, vals)
        file.Close()
    }
}
```

**İstifadə:**
```go
//go:generate ./queue MyInt MyFloat64     ← sətir!
package main

type MyInt int

func main() {
    q := NewMyIntQueue()                   // GENERATED funksiya!
    q.Insert(1)
    fmt.Println(q.Remove())
}
```

```bash
$ go generate            # myint_queue.go yaranır
$ go run myint.go myint_queue.go
First value: 1
```

**Vacib prinsiplər:**
- `$GOPACKAGE` (və $GOFILE, $GOLINE, $DOLLAR) — go generate tərəfindən təyin edilir
- Generator BUILD-dan ƏVVƏL, DEV lifecycle-da icra olunmur; build/runtime-da YOX
- **Generasiya olunan kod VCS-ə COMMIT olunmalı** — istifadəçilər generator işlətməməlidir
- generator build üçün əlçatan olmalıdır ($PATH və ya `./queue`)
- Alternativlər: go/ast (AST manipulyasiyası), yacc; template-lər ən sadə

**Müəllifin story:** bu pattern-dən ilhamlanıb **Helm**-də (Kubernetes package manager) template → manifest çevirici yazdılar.

**Trade-off xülasəsi:**
| Reflection | Code Generation |
|---|---|
| Runtime çevirmə — universal | Compile-time — type-safe |
| Yavaş + panik-yönlü | Sürətli; amma debug çətin |
| Boilerplate-yox | Dev prosesinə addım əlavə |

İkisi də VASİTƏDİR — "panacea deyil, strong typing workaround-u da deyil".

## Reflection API qısa xülasəsi
| Əməliyyat | Kod |
|---|---|
| Value alma | `reflect.ValueOf(x)` |
| Pointer izlə | `reflect.Indirect(v)` |
| Type alma | `v.Type()` / `reflect.TypeOf(x)` |
| Kind | `v.Kind()` |
| Struct sahə sayı/ad/dəyər | `t.NumField() / t.Field(i) / v.Field(i)` |
| Tag oxu | `field.Tag.Get("ad")` |
| İnterfeys implement? | `t.Implements(iface)` (nil ptr + Elem()) |
| Sahə yaz | `v.Field(i).Set(reflect.ValueOf(val))` |
| Böyük formaya çevir | `v.Int() / v.Uint() / v.Float()` |

## Əsas terminlər
- Reflection (öz strukturunu runtime-da imtahan)
- reflect.Value / reflect.Type / reflect.Kind
- Type Switch vs Kind Switch (MyInt dərsi)
- reflect.ValueOf / TypeOf / Indirect
- Int()/Uint() — kind-birləşdirmə çeviricilər
- Type Assertion interfeys yoxlaması
- Nil Pointer Trick (`(*fmt.Stringer)(nil)` + Elem())
- Type.Implements()
- NumField / Field (StructField) / Value.Field
- Struct Walking (rekursiv)
- Annotation (backtick, kompilyasiyada rolu yox)
- Struct Tag (`NAME:"VALUE,DATA"` konvensiyası)
- Tag.Get ("ini" kimi custom açarlar)
- `json:"-"` (ignore idiomu)
- validate tag (regex validation)
- Marshal/Unmarshal pattern (encoding/json kimi)
- v.Field(i).Set (runtime yazma)
- go generate / `//go:generate` direktivi
- $GOPACKAGE env
- text/template ilə kod generasiyası
- Typed Queue Generator (generics əvəzi — kitab dövrü)
- go/ast (AST əsaslı generasiya)
- Metaprogramming
- Helm (Kubernetes package manager — nümunə istifadə)

## Praktik nəticə
- Numeric/structural emal → kind switch (10 case 2-yə düşür); konkret tiplər (error handling) → type switch.
- İnterfeys yoxlaması runtime-da seçilirsə: nil pointer + reflect.TypeOf().Elem() + Implements() — triki əzbərlə.
- Struct sahələrini gəzərkən tipdən Field(i), valuedan Field(i) — EYNİ i İŞLƏT; pointer-lər reflect.Indirect ilə aç.
- Custom formatlar üçün tag-lər yaz: `ini:"ad"` kimi; Tag.Get avtomatik pars edir; tag yoxdursa sahə adına fallback.
- setField kind-switch + strconv çevirmələri + Field(i).Set — marshal/unmarshal generatorların əsası.
- Repetitive type-safe kod (collections və s.) → go generate + text/template; nəticəni VCS-ə commit et — build deyil, DEV alətidir.
- Reflection hər yerdə YOX — generics (Go 1.18+) bir çox bu kitabın generator istifadələrini əvəz edir (müəllim qeydi).

## Mənbə
Pages: 276-303 (PDF), book pages 253-280
