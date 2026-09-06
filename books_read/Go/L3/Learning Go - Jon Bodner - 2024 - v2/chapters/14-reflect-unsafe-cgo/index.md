# Chapter 14 — Here There Be Dragons: Reflect, Unsafe, and Cgo (Əjdahalar Ölkəsi: Reflect, Unsafe və Cgo)

## Bu chapter nədən bəhs edir?

Go-nun "qaçış qapıları": reflect paketi (runtime tip dəyərləndirməsi, dəyər yaratma/dəyişmə,
marshaler və funksiya generatoru nümunələri), unsafe (yaddaş manipulyasiyası, binary data),
cgo (C inteqrasiyası) və nə vaxt istifadə EDİLMƏMƏLİ haqqında məsləhətlər.

## Əsas fikirlər

### 1. Reflection — Runtime-da Tiplərlə İş
**Nədir:** Compile vaxtı məlum olmayan tiplərlə runtime-da işləmək; dəyişənləri yoxlamaq,
dəyişdirmək, yaratmaq.

**Standart kitabxanada istifadə yerləri:** database/sql, text/template, html/template,
fmt, errors.Is/As, sort.Slice*, encoding/* (JSON/XML struct tag oxunuşu). **Ortak xətt:
proqramın sərhədlərində data giriş/çıxışı.** Əlavə: `reflect.DeepEqual` — == ilə
müqayisə olunmayan (slice/map) dəyərlərin "dərin bərabərliyi" (testdə istifadə).

**3 anlayış:**
- **Type** — dəyişənin tipi: `reflect.TypeOf(v)`; `Name()` (ad — pointer/slice üçün boş),
  `Kind()` (quruluş: reflect.Struct, reflect.Ptr, reflect.Int...). Qayda: `type Foo struct`
  → kind = reflect.Struct, type adı = "Foo".
- **Kind-ə uyğun metodlar:** səhv kind-də metod çağırmaq PANIC edir (`NumIn` yalnız
  funksiya tipində). `Elem()` — pointer/slice/map/channel-də tərkib tipi.
- **Value** — dəyişənin dəyəri: `reflect.ValueOf(v)`; `.Type()`, `.Kind()`, `.Interface()`.

**Struct introspeksiya:**
```go
type Foo struct {
    A int    `myTag:"value"`
    B string `myTag:"value2"`
}
ft := reflect.TypeOf(Foo{})
for i := 0; i < NumField; i++ {
    curField := ft.Field(i)
    fmt.Println(curField.Name, curField.Type.Name(), curField.Tag.Get("myTag"))
    // A int value / B string value2
}
```

### 2. Dəyərlərin Oxunması/Yazılması
**Oxu:** `sv.Interface().([]string)` (type assertion ilə geri) və ya kind-specific:
`Bool()/Int()/String()/Float()/Bytes()` — səhv kind-də panic.

**Yazı (3 addım):** pointer ötür → `Elem()` → `Set` metodu:
```go
i := 10
iv := reflect.ValueOf(&i)     // 1. pointer
ivv := iv.Elem()              // 2. göstərilən dəyər
ivv.SetInt(20)               // 3. təyin et → i = 20
```
Pointer-siz `reflect.ValueOf` ilə Set* çağırmaq panic (call-by-value məntiqi — funksiyada
`*i = 20` ilə eyni). Primitivlər üçün SetBool/SetInt/SetUint/SetFloat/SetString; digərlər
üçün `Set(reflect.Value)`.

**Yeni dəyərlər:** `reflect.New(T)` (new kimi), `reflect.MakeSlice/MakeMap/
MakeMapWithSize/MakeChan` (make kimi). Tip yoxdursa — nil pointer hiyləsi:
```go
var stringType = reflect.TypeOf((*string)(nil)).Elem()  // nil-i *string-ə çevir → Elem
var stringSliceType = reflect.TypeOf([]string(nil))

ssv := reflect.MakeSlice(stringSliceType, 0, 10)
sv := reflect.New(stringType).Elem()
sv.SetString("hello")
ssv = reflect.Append(ssv, sv)
ss := ssv.Interface().([]string)
```

### 3. Interface-dəki nil-i Yoxlamaq
Interface nil = tip DƏYƏR hər ikisi nil (Ch7 tələsi). Dəyərin nil-ini reflection ilə:
```go
func hasNoValue(i interface{}) bool {
    iv := reflect.ValueOf(i)
    if !iv.IsValid() {
        return true                       // boş interface
    }
    switch iv.Kind() {
    case reflect.Ptr, reflect.Slice, reflect.Map, reflect.Func, reflect.Interface:
        return iv.IsNil()
    default:
        return false
    }
}
```
Sıra vacibdir: `IsValid` false ikən digər metodlar panic; `IsNil` yalnız nil-olabilən
kind-lərdə. Amma ideal — interface nil-dəyərlə düzgün işləyən kod yazmaq; bunu çarəsiz
hallara saxla.

### 4. Marshaler Nümunəsi — Praktik Reflection
CSV struct-larına mapping (`csv:"name"` tag-ləri ilə) — standard library-nin
JSON encoder-inin necə işlədiyinin miniatürü:

**Marshal:** `interface{}` parametr (oxumaq üçün pointer YOX) → Kind yoxlamaları
(slice olmalı, elem-i struct) → `NumField` + `Field(i)` + `Tag.Lookup("csv")` ilə
header → hər element `marshalOne`: kind-switch (`reflect.Int` → `fieldVal.Int()` →
strconv.FormatInt...).

**Unmarshal:** pointer tələb olunur (dəyişmə!) → `sliceValPtr.Elem()` → header-dən
`namePos` map-i → hər sətir üçün `reflect.New(structType).Elem()` → `unmarshalOne`:
tag → pozisiya → kind-switch → `SetInt/SetString/SetBool` → `sliceVal.Set(reflect.Append(...))`.

**Dərslər:** (1) dəyişəndə pointer, oxuyanda value; (2) kind-a görə switch; (3) error
halları `fmt.Errorf("cannot handle field of kind %v")`; (4) reflect kodu verbose olur —
real JSON encoder milyonlarla sətirdə eyni pattern-i izləyir.

### 5. Funksiya Yaratma — reflect.MakeFunc
**Pattern:** mövcud funksiyanı ortaq funksionallıqla wrap etmək (timing nümunəsi):
```go
func MakeTimedFunction(f interface{}) interface{} {
    ft := reflect.TypeOf(f)
    fv := reflect.ValueOf(f)
    wrapperF := reflect.MakeFunc(ft, func(in []reflect.Value) []reflect.Value {
        start := time.Now()
        out := fv.Call(in)
        fmt.Println(time.Since(start))
        return out
    })
    return wrapperF.Interface()
}
timed := MakeTimedFunction(timeMe).(func(int) int)
```
Real istifadə: müəllifin Proteus kitabxanası (SQL sorğusundan typesafe funksiya generasiyası).
**Qayda:** generasiya olunan funksiya aydın sənədlənsin; reflection yavaşdır — yalnız
zaten yavaş əməliyyatların (network çağırışı) wrap-inə layiqdir. `reflect.StructOf` —
runtime struct yaratmaq (yalnız interface{}-ə mənimsədilə bilən) — akademik maraq.
**Reflection METOD yarada BİLMƏZ** → runtime-da interface implement etmək mümkün deyil.

### 6. "Yalnız Dəyər Olduqda" — Filter BENCHMARK-ı
Reflection ilə universal Filter vs xüsusi funksiya (i7-8700, Go 1.14, 1000 element):
```
BenchmarkFilterReflectString-12   4822   229099 ns/op  87361 B/op  2219 allocs/op
BenchmarkFilterString-12       158197     7795 ns/op  16384 B/op     1 allocs/op
BenchmarkFilterReflectInt-12      4962   232885 ns/op  72256 B/op  2503 allocs/op
BenchmarkFilterInt-12           348441     3440 ns/op   8192 B/op     1 allocs/op
```
**30-70x yavaş** + minlərlə alloc + **type-safety itkisi** (yanlış tip runtime-da crash).
Nəticə: bir neçə sətir iqtisadına reflection xərci dəyməz — reflection sərhəd
(marshal/unmarshal) kodu üçündür.

### 7. unsafe — Yaddaş Səviyyəli Əməliyyatlar
**Tərkib:** `Sizeof/Offsetof/Alignof` (konstant qaytarırlar) + `unsafe.Pointer` (hər tip
pointer-ilə qarşılıqlı convert olunan körpü) + uintptr riyaziyyatı.

**Niyə var (2020 araşdırması, 2438 layihə):** 24% layihədə istifadə; 45.7% OS/C
interoperabilitesi; 23.6% performans.

**Binary data convert (wire protocol nümunəsi):**
```go
type Data struct {
    Value  uint32    // 4 bayt
    Label  [10]byte  // 10 bayt
    Active bool       // 1 bayt + 1 padding
}
func DataFromBytesUnsafe(b [16]byte) Data {
    data := *(*Data)(unsafe.Pointer(&b))   // bayt array → struct (bit_copy)
    if isLE {
        data.Value = bits.ReverseBytes32(data.Value)  // network big-endian düzəlişi
    }
    return data
}
```
Safe versiyadan ~2x sürətli (10.4 → 4.01 ns/op). Endianness yoxlaması özü də unsafe + init
(təcrid olunmaz dəyər — init üçün məqbul istifadə):
```go
var isLE bool
func init() {
    var x uint16 = 0xFF00
    xb := *(*[2]byte)(unsafe.Pointer(&x))
    isLE = (xb[0] == 0x00)
}
```

**String/slice header-ləri:** `reflect.StringHeader` (Data uintptr + Len) və
`reflect.SliceHeader` (+ Cap) — bayt-bayt gəzmək üçün; uintptr GC-ə qarşı kövrəkdir →
istifadənin sonunda **`runtime.KeepAlive(s)`** — GC-nin s-i vaxtından əvvəl toplamasının
qarşısını alır.

**Alət:** `-gcflags=-d=checkptr` — unsafePointer istismarını runtime-da yoxlayır (race
checker kimi, hamısını tutmur, yavaşıldırır — testdə işlət).

### 8. cgo — C İntegrasiyası
**Sintaksis:** C kodu `/* ... */` şərhində, dərhal ardınca `import "C"`; `C.add`,
`C.sqrt`, `C.CString` pseudo-paketi. Go→C export: `//export doubler` şərhi + C tərəfdə
`#include "_cgo_export.h"`.

**Məhdudiyyətlər (GC vs manual memory):**
- Pointer ötürülür, amma **pointer saxlayan dəyər YOX** (string, slice, funksiya, içində
  pointer olan struct).
- C funksiyası Go pointer-inin kopyasını return-dan sonraya saxlaya BİLMƏZ.
- Pozulma compile/runtime xətası verməyə bilər — sadəcə crash/yanlış davranış.
- Variadic C funksiyası çağırıla bilmir; C union → []byte; C funksiya pointer-i çağrılmır
  (dəyişənə mənimsədilə bilər).

**Performans reallığı:** Go→C çağırışı C→C-dən ~29x yavaş (Filippo Valsorda "Why cgo
is slow") — Python/Ruby-nin əksinə Go-da sürət üçün C-ə enmək məntiqi İŞLEMİR.
**Yeganə səbəb:** əvəzsiz C kitabxanası (wrapper axtar: SQLite, ImageMagick mövcuddur).
Oxu: Tobias Grieger "The Cost and Complexity of Cgo".

## Əsas terminlər
- Reflection — runtime tip/dəyər manipulyasiyası
- reflect.Type/Kind/Value — tip metaforasının 3 sütunu
- Elem — pointer/slice/map-in tərkibinə keçid
- reflect.MakeFunc — runtime funksiya generasiyası
- unsafe.Pointer — hər tip pointer üçün körpü
- uintptr — riyaziyyat üçün tam ədəd (GC-dən qorunmur)
- runtime.KeepAlive — GC-nin obyekti vaxtından əvvəl yığmasının bloku
- cgo — Go↔C FFI
- Wire format — şəbəkə bayt sırası (big-endian)

## Praktik nəticə

Bu üç alət "sərhəd alətləridir": reflect — xarici TEXT data ↔ Go tipləri (marshaler-lər);
unsafe — OS/şəbəkə BINARY data və miqyaslı performans (yalnız benchmark sübut edəndə);
cgo — yalnız əvəzsiz C kitabxanası. Hamısı üçün ortaq qaydalar: (1) default — safe Go;
(2) reflection core business məntiqində YOX (30-70x yavaş + type-safety itkisi);
(3) unsafe-da KeepAlive + checkptr; (4) cgo pointer qaydalarını pozmaq "bəzən işləyən"
crash-lər yaradır; (5) copy-paste etdiyin "aqıllı" həlləri başa düş — əjdahalar xəritənin
kənarında güzgüdür.

## Mənbə
Pages: 405-438
