# Chapter 5 — Reflection (səh. 83-101)

## Bu fəsil nədən bəhs edir?

`reflect` paketi ilə runtime introspection (kod öz strukturunu yoxlaması):
`reflect.Type` (tip naviqasiyası), `reflect.Value` (dəyər oxu/yaz), runtime
funksiya yaradılması (`MakeFunc`), struct tag-lər və Rob Pike-un refleksiyanın
üç qanunu.

## Əsas fikirlər

### 1. reflect.Type — tip introspection
**Nədir:** Proqramın runtime-da öz tip strukturlarını yoxlamaq (metaproqramlaşdırma
formasısı). `reflect.TypeOf(i interface{})` — boş interfeys qəbul edir, çünki
`interface{}` hər tipi saxlaya bilər.

**Əsas tiplərlə:**
```go
var a int32 = 42
var b string = "forty two"
fmt.Println(reflect.TypeOf(a))   // int32
fmt.Println(reflect.TypeOf(b))   // string
```

**Struct sahələrinin naviqasiyası:**
```go
type T struct {
    A int32
    B string
}

t := T{42, "forty two"}
typeT := reflect.TypeOf(t)
fmt.Println(typeT)                 // main.T

for i := 0; i < typeT.NumField(); i++ {
    field := typeT.Field(i)
    fmt.Println(field.Name, field.Type)   // A int32 / B string
}
```

**Sub-kod izahı:**
- `NumField()` → struct sahələrinin sayı
- `Field(i)` → i-ci `StructField` (adı, tipi, tag-i oxunur)

**Interfeys implementasiyasının yoxlanılması:**
```go
type Adder interface{ Add(int, int) int }
type Calculator struct{}
func (c *Calculator) Add(a int, b int) int { return a + b }

var ptrAdder *Adder
adderType := reflect.TypeOf(ptrAdder).Elem()  // Adder interfeysinin tipi

c := Calculator{}
calcType := reflect.TypeOf(c)       // main.Calculator
calcTypePtr := reflect.TypeOf(&c)   // *main.Calculator

calcType.Implements(adderType)      // false — Add *Calculator-dədir!
calcTypePtr.Implements(adderType)   // true
```
- `Implements()` → pointer receiver metodları yalnız pointer tipinə aid olur
  (Chapter 4 qaydasının runtime təsdiqi).

**Rekursiv struct inspector:** `Kind()` == `reflect.Struct` yoxlaması ilə
nested struct-lar rekursiv gəzilə bilər (`typeT.Kind()`, `f.Type.Kind()`).

### 2. reflect.Value — dəyər çıxarma
**Nədir:** `reflect.Type` dəyərlərə çıxış vermir; dəyərlər `reflect.Value`
ilə oxunur. `Value`-dən `Type` alınır, əksi yox.

```go
var a int32 = 42
valueOfA := reflect.ValueOf(a)
fmt.Println(valueOfA.Interface())   // 42 — Interface() ilə real dəyər
```

**Kind + switch ilə tipə görə emal:**
```go
func ValuePrint(i interface{}) {
    v := reflect.ValueOf(i)
    switch v.Kind() {
    case reflect.Int32:
        fmt.Println("Int32 with value", v.Int())
    case reflect.String:
        fmt.Println("String with value", v.String())
    default:
        fmt.Println("unknown type")
    }
}
```
- `Kind()` → Go tiplərindən birini təmsil edən rəqəm (Int32, String, Map...)
- `v.Int()`, `v.String()`, `v.Float()` → tipə uyğun dəyər çıxarır

**Struct dəyərləri:**
```go
t := T{42, "forty two", 3.14}
valueT := reflect.ValueOf(t)
fmt.Println(valueT.Kind(), valueT)   // struct {42 forty two 3.14}

for i := 0; i < valueT.NumField(); i++ {
    field := valueT.Field(i)
    fmt.Println(field.Kind(), field.String(), field.Interface())
}
```
- Rəqəmsal sahələrdə `field.String()` → `<int32 Value>` kimi tip göstəricisi
  verir; real dəyər üçün `field.Interface()` və ya `fmt.Println(field)`.

### 3. Dəyərlərin dəyişdirilməsi (Set)
**Nədir:** `SetString`, `SetInt` və s. ilə sahə dəyərləri runtime-da dəyişilir.
Tələb: pointer üzərindən `Elem()` ilə.

**Kitabdan kod nümunəsi (string-ləri böyük hərfə çevir):**
```go
t := T{"hello", 42, "bye"}
fmt.Println(t)   // {hello 42 bye}

valueOfT := reflect.ValueOf(&t).Elem()   // pointer + Elem → settable
for i := 0; i < valueOfT.NumField(); i++ {
    f := valueOfT.Field(i)
    if f.Kind() == reflect.String {
        current := f.String()
        f.SetString(strings.ToUpper(current))
    }
}
fmt.Println(t)   // {HELLO 42 BYE}
```

**CanSet() — unexported sahələrin qorunması:**
```go
type T struct {
    A string
    B int32
    c string   // unexported (kiçik hərf)
}
// ...
if f.CanSet() {
    f.SetString(strings.ToUpper(current))
}
```
- Unexported sahəyə Set → **panic**; `CanSet()` yoxlaması bunu ağlabatan keçir
  (`{HELLO 42 bye}` — c dəyişməz qalır).

### 4. reflect.MakeFunc — runtime funksiya yaratma
**Nəyə lazımdır:** Go-da funksiya overloading və (o dövrdə) generics yoxdur.
Eyni "toplama" məntiqini müxtəlif tiplər üçün bir dəfə yazmaq üçün `MakeFunc`
funksiya faktoriyası qurulur.

**Kitabdan kod nümunəsi (BuildAdder):**
```go
func BuildAdder(i interface{}) {
    fn := reflect.ValueOf(i).Elem()

    newF := reflect.MakeFunc(fn.Type(), func(in []reflect.Value) []reflect.Value {
        if len(in) > 2 { return []reflect.Value{} }
        a, b := in[0], in[1]
        if a.Kind() != b.Kind() { return []reflect.Value{} }

        var result reflect.Value
        switch a.Kind() {
        case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
            result = reflect.ValueOf(a.Int() + b.Int())
        case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
            result = reflect.ValueOf(a.Uint() + b.Uint())
        case reflect.Float32, reflect.Float64:
            result = reflect.ValueOf(a.Float() + b.Float())
        case reflect.String:
            result = reflect.ValueOf(a.String() + b.String())
        default:
            result = reflect.ValueOf(interface{}(nil))
        }
        return []reflect.Value{result}
    })
    fn.Set(newF)
}

func main() {
    var intAdder func(int64, int64) int64
    var floatAdder func(float64, float64) float64
    var strAdder func(string, string) string

    BuildAdder(&intAdder)
    BuildAdder(&floatAdder)
    BuildAdder(&strAdder)

    fmt.Println(intAdder(1, 2))        // 3
    fmt.Println(floatAdder(3.0, 2.423)) // 5.423
    fmt.Println(strAdder("hello", " go")) // hello go
}
```

**Sub-kod izahı:**
- `reflect.MakeFunc(tip, impl)` → verilmiş imzalı funksiya yaradır; impl
  `[]reflect.Value` alır və qaytarır
- `fn.Set(newF)` → pointer arqumentə yeni funksiya mənimsədilir
- Nəticə: bir implementasiya — int, float və string üçün işləyir

### 5. Struct Tag-lər
**Nədir:** Struct sahələrinin metadatası; reflection ilə runtime-da oxunur
(JSON, validation, ORM açar sözləri buradan gəlir).

```go
type User struct {
    UserId   string `tagA:"valueA1" tagB:"valueA2"`
    Email    string `tagB:"value"`
    Password string `tagC:"v1 v2"`
    Others   string `something a b`   // konvensiyaya uyğun deyil
}
```
- Konvensiya: `key:"value"` cütlükləri, boşluqla ayrılır; `key: "value"`
  (iki nöqtədən sonra boşluq) düzgün deyil.

**Tag-lərin oxunması:**
```go
T := reflect.TypeOf(User{})

fieldUserId := T.Field(0)
t := fieldUserId.Tag
fmt.Println("StructTag is:", t)        // tam tag sətri
v, _ := t.Lookup("tagA")
fmt.Printf("tagA: %s\n", v)            // valueA1

fieldEmail, _ := T.FieldByName("Email")
fmt.Println("email tagB:", fieldEmail.Tag.Get("tagB"))  // value
```
- `Tag.Lookup(key)` → (dəyər, var/yox) qaytarır
- `Tag.Get(key)` → dəyəri verir

### 6. Refleksiyanın üç qanunu (Rob Pike)
1. **Refleksiya interface dəyərindən refleksiya obyektinə gedir.**
   `TypeOf`/`ValueOf` interface{} qəbul edir — hər dəyər əvvəlcə boş
   interfeysə "paketlənir", sonra tip/dəyər çıxarılır.
2. **Refleksiya obyektindən interface dəyərinə gedir.**
   `v.Interface()` → Value-dan real dəyəri qaytarır:
   ```go
   v := reflect.ValueOf(a)
   fmt.Println(v)             // 42 — Println xüsusi hallarda Value-ni açır
   fmt.Println(v.Interface()) // 42 — burada artıq int32 çap olunur
   ```
3. **Refleksiya obyektini dəyişmək üçün dəyər settable olmalıdır.**
   ```go
   var a int32 = 42
   v := reflect.ValueOf(a)
   v.SetInt(16)              // PANIC — kopya, settable deyil

   v = reflect.ValueOf(&a)
   v.SetInt(16)              // PANIC — bu, pointer-in özüdür

   v = reflect.ValueOf(&a).Elem()
   v.SetInt(16)              // İşləyir — Elem() göstərilən dəyəri verir
   ```
   Settablelik `CanSet()` ilə yoxlanılır.

## Əsas terminlər
- Reflection (refleksiya) — proqramın öz strukturunu runtime-da yoxlaması
- Introspection (introspeksiya) — tip/sahə məlumatının oxunması
- Kind — tipin kateqoriyası (Int32, String, Struct, Map...)
- Settable (təyin edilə bilən) — yalnız pointer+Elem ilə yazmağa icazə
- Struct tag — sahə metadatası; `key:"value"` konvensiyası
- MakeFunc — runtime funksiya konstruktoru

## Praktik nəticə
Reflection güclü, amma performans/oxunaqlılıq/idiarəetmə xərcli alətdir —
yalnız başqa yox dərmanın olduqda (serialization, ORM, generic-məntiq)
işlədin. Yazma əməliyyatı üçün mütləq `ValueOf(&x).Elem()` + `CanSet()`
pattern-i; unexported sahələrə yazmaq panic verir. Tag-lər `Lookup`/`Get`
ilə oxunur — JSON və DB alanlarının mənbəyi məhz bu mexanizmdir.

## Mənbə
Pages: 83-101 (PDF 83-101)
