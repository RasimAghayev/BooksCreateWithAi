# Chapter 15 — A Look at the Future: Generics in Go (Gələcəyə Baxış: Go-da Generics)

## Bu chapter nədən bəhs edir?

Generics-in (type parameters) səbəbləri, draft dizaynı: `[T any]` sintaksisi, `any` və
`comparable` universe identifier-ləri, generic interfeyslər, type list-lər (operator
məhdudiyyətləri), generic funksiyalar (Map/Reduce/Filter), nəyin DAXİL EDİLMƏDİYİ və
idiomatik Go-ya təsiri.

## Əsas fikirlər

### 1. Niyə Generics?
**Problem:** Builtin tiplər (map, slice, len...) tipli generik işləyir, amma
istifadəçi tipləri/funksiyaları YOX. Nəticə: hər tip üçün ayrı binary tree; `interface{}`
ilə type-safety itkisi; `[]string` → `[]interface{}` təyin olunmur; slice funksiyaları hər
tip üçün təkrarlanır (yalnız reflection ilə ümumiləşir — Ch14-də gördük: 30-70x yavaş).

**interface{} həllinin tələsi (Orderable nümunəsi):** `Order(interface{}) int` metodu
ilə tree tipləri qarışdıra bilərsən — compile KEÇİR, runtime-da panic:
```go
it = it.Insert(OrderableInt(5))
it = it.Insert(OrderableString("nope"))   // compile OK → runtime panic!
```
Russ Cox (2009): generics-i istisna etdilər, çünki sürətli kompilyator + oxunaqlı kod +
yaxşı icra vaxtının ÜÇÜ birlikdə mövcud həllərdə mümkün deyildi. 10 ildən sonra — Go 1.18
üçün Type Parameters Draft Design.

### 2. Generic Sintaksis — Stack Nümunəsi
**Kitabdan kod nümunəsi:**
```go
type Stack[T any] struct {
    vals []T
}
func (s *Stack[T]) Push(val T) {
    s.vals = append(s.vals, val)
}
func (s *Stack[T]) Pop() (T, bool) {
    if len(s.vals) == 0 {
        var zero T          // zero value almaq üçün var hiyləsi (nil value tipinə OLMAZ!)
        return zero, false
    }
    top := s.vals[len(s.vals)-1]
    s.vals = s.vals[:len(s.vals)-1]
    return top, true
}

var s Stack[int]      // instansiya tip parametri ilə
s.Push("nope")        // COMPILE XƏTASI: cannot convert "nope" to int
v, ok := s.Pop()      // v int tipində — assertion YOX!
```

**Qaydalar:**
- Tip parametrləri `[]`-də, ad + bound: `[T any]`; adətən böyük hərflər.
- `any` — universe block-da yeni identifier; `interface{}`-in bərabəri, amma yalnız tip
  constraint kontekstində.
- Receiver: `Stack[T]` (yalnız `Stack` YOX).
- Zero value: `var zero T` (həmişə işləyir).
- Metodlar üstündə tip parametri receiver-dən gəlir.

### 3. comparable — == İcazəsi
`any` operatorlardan xəbər vermir: `v == val` → compile xətası ("operator == not defined
for T"). Həll — universe block-da yeni `comparable` interfeysi (==/!= mümkün olan bütün
tiplər — struct-ların kəsişimi xaric):
```go
type Stack[T comparable] struct { vals []T }
func (s Stack[T]) Contains(val T) bool {
    for _, v := range s.vals {
        if v == val { return true }
    }
    return false
}
```

### 4. Generic İnterfeys və Ağac Nümunəsi
```go
type Orderable[T any] interface {
    Order(T) int    // parametr tipi T — type-safe!
}
func (oi OrderableInt) Order(val OrderableInt) int { return int(oi - val) }  // assertion YOX

type Tree[T Orderable[T]] struct {
    val         T
    left, right *Tree[T]
}
func (t *Tree[T]) Insert(val T) *Tree[T] {
    if t == nil {
        return &Tree[T]{val: val}    // instansiya Tree[T] ilə
    }
    switch comp := val.Order(t.val); { ... }
}
var it *Tree[OrderableInt]
it.Insert(OrderableString("nope"))   // COMPILE XƏTASI — runtime panic mərhələsi GƏLDİ
```

### 5. Type List-lər — Operatorlar üçün
Wrapper tiplərə ehtiyac duymadan builtin tiplərlə işləmək — interface daxilində tip
siyahısı; icazəli operatorlar = siyahıdakı HAMISINDA işləyənlər:
```go
type BuiltInOrdered interface {
    string, int, int8, int16, int32, int64,
    float32, float64, uint, uint8, uint16, uint32, uint64, uintptr
}   // (draft sintaksisi; Go 1.18-də ~ göstərilir)
type Tree[T BuiltInOrdered] struct { ... }
// <, >, ==, !=, + istifadə oluna bilər
var it *Tree[int]
it = it.Insert(5)   // birbaşa int — wrapper YOX
```

### 6. Generic Funksiyalar — Map/Reduce/Filter
**Kitabdan kod nümunəsi:**
```go
func Map[T1, T2 any](s []T1, f func(T1) T2) []T2 {
    r := make([]T2, len(s))
    for i, v := range s {
        r[i] = f(v)
    }
    return r
}
func Reduce[T1, T2 any](s []T1, initializer T2, f func(T2, T1) T2) T2 {
    r := initializer
    for _, v := range s {
        r = f(r, v)
    }
    return r
}
func Filter[T any](s []T, f func(T) bool) []T {
    var r []T
    for _, v := range s {
        if f(v) { r = append(r, v) }
    }
    return r
}
```
Tip çıxarışı (inference) `:=` kimi işləyir; yalnız return-də istifadə olunan parametr
çıxarıla bilməz → açıq göstər: `Convert[int, int64](a)`.

### 7. Type List-lərin Dəqiqliyi
- **Konstantlar:** hamısına uyğun konstant YOXDURSA təyin OLMAZ (`Integer` siyahısında
  `in + 1_000` — int8-ə sığmır → YANLIŞ; `+ 100` — hər siyahı üzvünə sığır → DÜZGÜN).
- **Type literal-lar:** `type []T, map[E]T` kimi slice/map şablonları da siyahıya düşə
  bilər (CopyVals nümunəsi — hər hansı map/slice-dən dəyərləri çıxarır).
- **Metodlar:** siyahıdakı user-defined tiplərin metodları İGNOR olunur — metod lazımdırsa
  interfeysdə ELAN edilməli; elan olunarsa siyahıdakı HAMISI onu implement etməli.
  `type MyInt, MyFloat` + `String() string` → yalnız bu tiplər; `int, float64` + `String()`
  → heç bir builtin uyğun gəlmir (metodları yoxdur).
- Type list-li interfeys YALNIZ tip parametri kimi istifadə oluna bilər.

### 8. Nə DAXİL EDİLMƏYİB (və niyə)
- **Operator overloading YOX** — user-defined tip üçün `<` təyin etmək olmaz; range/[]
  da user container-lər üçün işləməz. Səbəb: oxunaqlıq (C++ `<<` fəlakəti), overload
  mexanizmi yoxdur.
- **Parameterized metodlar YOX:** `func (fs functionalSlice[T]) Map[E any](...)` —
  MÜMKÜNSÜZ; funksional zəncir (`xs.Map(f).Reduce(...)`) yoxdur — ya nested çağırışlar,
  ya da addım-addım aralıq dəyişənlər.
- **Variadic tip parametrləri YOX** — istənilən imzalı funksiya wrap etmək üçün hələ
  reflection lazımdır.
- **Specialization, currying, metaprogramming YOX** — Go kiçik qalır.

### 9. İdiomatik Təsir və Gələcək
- `float64`-ün "universal numeric" rolu bitir; `interface{}` data strukturlarında/parametrlərdə
  generics-ə yerini verir; slice funksiyaları təkrarlanmaz.
- Mövcud kodu dərhal göçürmə — YOX; yeni pattern-lər təkmilləşəcək; benchmark/profile
  hələ vacibdir (generics-in performans təsiri real layihələrdə ölçüləcək).
- Gözlənilən standart kitabxana dəyişiklikləri: Orderable kimi interfeyslər, set/tree/
  ordered map tipləri, yeni funksiyalar. `any` gələcəkdə ümumi `interface{}` əvəzedici
  ola bilər.
- **Sum types perspektivi:** type list-lərin interfeys parametrlərinə şamil edilməsi —
  JSON-un "tək dəyər VƏ YA siyahı" problemindən Rust/Swift tipli enum-lara yol.

## Əsas terminlələr
- Type parameter (tip parametri) — `[T any]` formalı generik dəyişən
- `any` — constraint kontekstində interface{} qarşılığı
- `comparable` — ==/!= dəstəyi olan tiplər üçün universe interfeysi
- Type list (tip siyahısı) — interfeysdaxili tip sadalaması (operator icazəsi)
- Type inference — tip arqumentlərinin avtomatik çıxarılması
- Operator overloading (operator üstəyükləməsi) — Go-da OLMAYAN xüsusiyyət
- Parameterized method (parametrləşdirilmiş metod) — Go-da İCAZƏSİZ konstruksiya
- Sum type (cəm tipi) — gələcək üçün müzakirə olunan məhdud tip dəsti

## Praktik nəticə

Generics Go-nu dəyişir, amma fəlsəfəsini YOX: (1) əvəzlədikləri — interface{} data
strukturları, reflection-lu sort/filter, hər-tip-üçün-kopyala-yaz; (2) saxladıqları —
explicit, oxunaqlı, kiçik dil; (3) zero value üçün `var zero T` hiyləsini bil; (4) operator
lazımdırsa type list, metod lazımdırsa generic interfeys; (5) metodda YENİ tip parametri
mümkünsüzdür — funksiya yaz; (6) köhnə kod işləməyə davam edir — panik edib hər şeyi
generic-ləşdirmə. Sonda kitabın vədi: düzgün yazılmış Go "sıxdırıcı"dır — bu xüsusiyyətdir,
qüsur yox.

## Mənbə
Pages: 439-460
