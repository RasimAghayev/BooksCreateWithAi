# Chapter 4 — Operations (Əməliyyatlar)

## Bu fəsil nədən bəhs edir?

Type parameter üzərində mümkün əməliyyatlar və onları mümkün edən
constraint-lər: arifmetika (Integer/Float/Complex → Number; ~ ilə derived
tiplər daxil), ordered types (Real → Ordered → cmp.Ordered standart
kitabxanada), multiple type parameters (T, U — HƏR BİRİ ayrıca tipdir,
operatorlar aralarında işləmir), derived slice tip problemi ([S ~[]E, E any]
idiomu — slices paketinin Clone imzası), comparable (predeclared — sonsuz
sayda comparable tip var, interface ilə ifadə OLMAZ), abstract type anlayışı
(type parameter-in body-də istifadəsi, var top E, zero value E(0)), type
switch-in generic-də QADAĞASI və `switch any(x).(type)` həlli.

## Əsas fikirlər

### 1. Arifmetika — Number Constraint
```go
type Integer interface {
    ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 |
        ~uint16 | ~uint32 | ~uint64
}
type Float interface {
    ~float32 | ~float64
}
type Complex interface {
    ~complex64 | ~complex128
}
type Number interface {
    Integer | Float | Complex
}

func AddAnything[T Number](x, y T) T {
    return x + y
}
fmt.Println(AddAnything(1 + 2i, 3 + 4i))
// Output: (4+6i)
```
- **Nədir:** bütün built-in rəqəm tipləri + ~ (derived) → + operatoru
  zəmanətli
- **Necə işləyir:** type set-in HƏR elementi + dəstəkləyir → Go compile
  edir
- **Product məşqi:** `func Product[T Number](x, y T) T { return x * y }` —
  eyni Number constraint * üçün də kifayətdir.

### 2. Ordered Types — > Operatoru
```go
func Greater[T Number](x, y T) T {
    if x > y { return x }
    return y
}
// invalid operation: x > y (type parameter T is not comparable with >)
```
- **Problem:** Number-a Complex daxildir; kompleks ədədlərin təbii
  SIRALANMASI yoxdur → > işləmir
- **Həll — Real constraint:**
```go
type Real interface {
    Integer | Float
}
func Greater[T Real](x, y T) T { ... }
fmt.Println(Greater(1, 2))  // 2
```
- **Ordered anlayışı:** tip ordered-dir, əgər onun qiymətlərini kiçikdən
  böyüyə müəyyən qayda ilə düzmək mümkündürsə. String-lər də ordered-dir
  ("b" > "a").
```go
type Ordered interface {
    Integer | Float | ~string
}
func Greater[T Ordered](x, y T) T {
    return max(x, y)   // built-in max/min də eyni constraint tələb edir
}
```
- **cmp.Ordered:** standart kitabxanada HAZIR: `import "cmp"` →
  `func Greater[T cmp.Ordered](x, y T) T`. Gələcəkdə yeni ordered tipləri
  əlavə olunsa belə cmp.Ordered aktual qalır.

### 3. Multiple Type Parameters — T və U FƏRLİ TİPLƏRDİR
```go
func Identity[T, U any](x T, y U) (T, U) {
    return x, y
}
fmt.Println(Identity(1, "hello"))  // 1 hello
```
- **Nədir:** iki (və daha çox) type parameter — vergüllə ayrılır
- **SIRTI QAYDA:** T və U eyni constraint-i qane etsə belə AYRI TİPLƏRDİR.
  `x > y` və `x + y` yalnız eyni tip arasında işləyir → T və U arasında
  operator İŞLƏMİR. int və float64 hər ikisi cmp.Ordered-dir, amma
  müqayisə oluna bilməz. Operator lazımdırsa → tək T istifadə et.

### 4. Derived Slice Problemi və [S ~[]E, E any] İdiomu
```go
type StringList []string

func Identity[E any](s []E) []E {   // PİS variant
    return s
}
result := Identity(StringList{"a", "b", "c"})
fmt.Printf("result is a %T\n", result)
// result is a []string  — StringList İTİRİLDİ!

func (s StringList) Len() int { return len(s) }
result.Len()
// result.Len undefined (type []string has no field or method Len)
```
- **Nədir:** `[]E` constraint-i yalnız adlandırılmamış slice-ları qəbul edir;
  StringList ondan TÖRƏMİP tipdir — funksiya `[]string` qaytarır
- **Həll — S type parameter-i ~ ilə:**
```go
func Identity[S ~[]E, E any](s S) S {
    return s
}
result := Identity(StringList{"a", "b", "c"})
fmt.Println(result.Len())  // 3 — tip SAXLANILDI
```
- **Vacib:** slice qaytaran və ya slice elementi qaytaran generic
  funksiyalarda HAMİŞƏ `[S ~[]E, E any]` stilindən istifadə et — slices
  paketinin `func Clone[S ~[]E, E any](s S) S` imzası da belədir.

### 5. Comparable — == Operatoru
```go
func Equal[T any](x, y T) bool {
    return x == y
}
// invalid operation: x == y (incomparable types in type set)
```
- **Problem:** slice, map, channel, funksiya == DƏSTƏKLƏMİR → any çox genişdir
- **Ordered də yaramır:** bool == dəstəkləyir amma ordered DEYİL (true > false
  mənasızdır!). Comparable-amma-ordered tiplər: bool, struct-lar, complex,
  pointer, channel, array, interface.
- **Sonsuz dəstə problemi:** comparable tiplərin siyahısı SONSUZDUR — struct
  tiplərini adlarla sadalamaq mümkün deyil; metod dəsti də yaramır (==
  metod ÇAĞIRMIR).
- **comparable — predeclared:**
```go
func Equal[T comparable](x, y T) bool {
    return x == y
}
fmt.Println(Equal(1, 1))  // true
```
- **Niyə predeclared?** cmp paketi təmiz Go ilə yazılıb; comparable-ı Go-da
  ifadə etmək MÜMKÜN DEYİL → compiler daxili implementasiyasıdır.
- **Dupes məşqi:**
```go
func Dupes[E comparable](s []E) bool {
    seen := map[E]bool{}
    for _, v := range s {
        if seen[v] {
            return true
        }
        seen[v] = true
    }
    return false
}
```
  (map-in açarı E — comparable tələbi buradan da gəlir; nəticə bool olduğu
  üçün ~[]E lazım DEYİL.)

### 6. Abstract Type — Body-də İstifadə
```go
func Greatest[S ~[]E, E cmp.Ordered](s S) E {
    if len(s) < 1 {
        panic("Greatest: empty slice")
    }
    var top E          // type parameter BODY-də — abstract type!
    for _, e := range s {
        top = max(top, e)
    }
    return top
}
```
- **Nədir:** type parameter-in funksiya body-də istifadəsi — abstract type;
  scope funksiya sonunadək
- **Instantiate olunanda** sadəcə konkret tip olur (`var top string` kimi)
- **Standart alternativ:** `slices.Max(s)` — slices paketində hazır var.
- **Zero value almağın 2 yolu:**
  - `var top E` — həmişə işləyir
  - `return E(0)` — yalnız 0 konstantı E-nin BÜTÜN mümkün tiplərinə
    çevrilə bildikdə (məs. Number üçün YOX — complex-in 0-ı ayrıdır)

### 7. Type Switch — Generic-də QADAĞA və Həlli
```go
func Identify[T any](x T) {
    switch x.(type) { ... }
}
// cannot use type switch on type parameter value x
```
- **Səbəb:** type switch YALNIZ interface tipləri üzərində işləyir; T
  zəmanətən interface DEYİL (int kimi konkret ola bilər — konkret tip
  üzərində switch MƏNASIZDIR)
- **Həll — any-yə çevir:**
```go
func Identify[T any](x T) {
    switch any(x).(type) {
    case int:
        fmt.Println("Looks like an int")
    case float64:
        fmt.Println("Looks like a float")
    default:
        fmt.Println("Doesn't look like anything to me")
    }
}
Identify(99)  // Looks like an int
```
- **Arzuolunmazlıq qeydi:** generic funksiyada type switch adətən lazımsızdır —
  hər tip üçün ayrıca funksiya yaz, ya da adi `any`-parametrli funksiya yaz.

## Əsas terminlər
- Ordered type (sıralana bilən tip) — <, >, <=, >= dəstəklənir
- Comparable type (müqayisə olunan tip) — == dəstəklənir (ordered olmaya bilər!)
- Abstract type (abstrakt tip) — body-də istifadə olunan type parameter
- Zero value (sıfır qiymət) — tipin default qiyməti; `var top E` və ya `E(0)`
- Predeclared (əvvəlcədən elan edilmiş) — comparable kimi compiler daxili
- Type set (tip dəsti) — constraint-i qane edən bütün tiplər

## Praktik nəticə

1. **Constraint seçim cədvəli (operator → constraint):**
   - `+ * - /` → `Number` (Integer | Float | Complex, hamısı ~ ilə)
   - `< > <= >=`, `min`/`max` → `cmp.Ordered` (Real + ~string)
   - `== !=` → `comparable`
   - heç bir operator → `any`
2. **Slice qaytaran generic:** `[S ~[]E, E any]` — əks halda derived slice
   tipi itir, metodları çağırılmır.
3. **Birden çox type parameter:** operatorlar yalnız eyni T daxilində —
   T, U arası operator QADAĞANDIR.
4. **Zero value:** `var x E` universal; `E(0)` yalnız çevirmə mümkünsə.
5. **Type switch lazımdırsa:** `any(x).(type)` ilə; amma əvvəlcə düşün —
   ayrıca funksiyalar daha idiomatikdir.

## Mənbə
Pages: 66-90 (PDF 67-91)
