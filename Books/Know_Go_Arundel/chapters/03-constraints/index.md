# Chapter 3 — Constraints (Məhdudiyyətlər)

## Bu fəsil nədən bəhs edir?

Generics-də constraint (məhdudiyyət) anlayışı: any kasıblıq problemindən
çıxış yolu. Basic interface (yalnız metod elementləri), type set constraint
(düz type elementləri), union (| ayırıcı), intersection (çoxsətirli interface),
empty type set xətası, composite/struct type literal-ləri, struct sahə
müraciətinin limitasiyası, "constraints are not classes" prinsipi,
approximation (~ — derived tipləri də qəbul edir), interface literal-i
constraint kimi (interface{} qısaltması, ad keçirmə — yalnız TİK bir type
element olanda), type parameter-ə constraint daxilində müraciət
(Equal(T) bool — adlandırılmış interface ilə MÜMKÜNSÜZ, yalnız literal ilə).

## Əsas fikirlər

### 1. any-nin Limitasiyası — Operator Qadağası
```go
func AddAnything[T any](x, y T) T {
    return x + y  // compile xətası: operator + not defined on T
}
```
- **Nədir:** `any` bütün tipləri qəbul edir, amma Go üçün + operatorunun
  T üzərində işləyəcəyinə ZƏMANƏT vermir
- **Go proverb:** "The bigger the interface, the weaker the abstraction" —
  constraint nə qədər genişdirsə, zəmanət o qədər azdır
- **Nəyə lazımdır:** any ilə yalnız dəyişən elan etmək, mənimsətmək, qaytarmaq
  olar — HESAB YOX

### 2. Basic Interface — Metod Dəsti Constraint
```go
func Stringify[T fmt.Stringer](s T) string {
    return s.String()
}
fmt.Println(Stringify(1))
// int does not implement Stringer (missing method String)
```
- **Nədir:** yalnız metod elementlərindən ibarət interface (String() string)
- **Necə işləyir:** T yalnız Stringer implement edən tiplərlə instantiate olunur
- **Vacib:** burada generic yazmağın ÜSTÜNLÜYÜ YOX — adi funksiya
  (parametr tipi Stringer) eyni işi görür; amma basic interface constraint
  kimi VALID-dır
- **StringifyTo məşqi:**
```go
func StringifyTo[T fmt.Stringer](w io.Writer, p T) {
    fmt.Fprintln(w, p.String())
}
```

### 3. Type Set Constraint — İcazələnən Tiplərin Siyahısı
```go
type OnlyInt interface {
    int
}
func Double[T OnlyInt](v T) T {
    return v * 2
}
```
- **Nədir:** metod YOX — interface daxilində birbaşa tip adı (type element)
- **Necə işləyir:** type set = ancaq int → Go * operatoruna ZƏMANƏT verə
  bilir (int * int həmişə işləyir)
- **Çatışmamazlığı:** yalnız int → int-dən başqa * destekleyən tiplər kənar
  qalır

### 4. Union — Bir Neçə Tipin Birləşməsi
```go
type Integer interface {
    int | int8 | int16 | int32 | int64
}
type Float interface {
    float32 | float64
}
type Complex interface {
    complex64 | complex128
}
type Number interface {
    Integer | Float | Complex   // constraint-lərdən kompozisiya!
}
```
- **Nədir:** `|` (pipe) ilə ayrılmış tip elementləri — "və ya" mənası
- **Necə işləyir:** tip int VE YA int8 VE YA ... olarsa constraint-i
  qane edir
- **Üstünlüyü:** constraint-lərin kompozisiyası — yeni Number constraint-i
  üç hazır constraint-dən yığılır
- **Type set:** Number-in type set-i = Integer ∪ Float ∪ Complex

### 5. Intersection — Çoxsətirli Interface
```go
type ReaderStringer interface {
    io.Reader
    fmt.Stringer
}
// qısa forma: interface { io.Reader; fmt.Stringer }
```
- **Nədir:** hər sətir ayrıca element — type set = elementlərin KƏSİŞMƏSİ
- **Necə işləyir:** tip hər İKİSİNİ implement etməlidir (Reader VƏ Stringer),
  yalnız biri yetərli DEYİL
- **Boş type set:**
```go
type Unpossible interface {
    int
    string
}
// cannot implement Unpossible (empty type set)
```
  Heç bir tip eyni anda int və string ola bilməz → instantiate xətası.
  Adətən səhv nəticəsidir, amma xəta mesajını tanımaq vacibdir.

### 6. Composite Type Literals — Anında Tip Yaratma
```go
type Pointish interface {
    struct{ X, Y int }
}
```
- **Nədir:** adlandırılmış tip YOX — type literal (sintaksisin özü) element
- **Necə işləyir:** type set tam olaraq bir tipdən ibarətdir: struct{ X, Y int }
- **Limitasiya:** struct sahələrinə çıxış YOXDUR:
```go
func GetX[T Pointish](p T) int {
    return p.X
}
// p.X undefined (type T has no field or method X)
```
  Constraint X sahəsini zəmanət etsə belə compiler bunu hələ dəstəkləmir.
- **İkinci limitasiya:** type elementli interface yalnız CONSTRAINT kimi
  işləyə bilər — adi parametr tipi kimi YOX (`func Double(p Number) Number` →
  "interface contains type constraints" xətası)

### 7. Constraints Are Not Classes
```go
type Cow struct{ moo string }
type Chicken struct{ cluck string }
type Animal interface {
    Cow | Chicken
}
type Farm[T Animal] []T

dairy := Farm[Cow]{}
poultry := Farm[Chicken]{}
mixed := Farm[Animal]{}   // ERROR: interface contains type constraints
```
- **Nədir:** class-hiyerarxiyalı dillərdən fərqli olaraq constraint = tip DEYİL
- **Nəyə lazımdır:** `Farm[Animal]` yazmaq olmaz — Animal özü type set-in
  elementi deyil, sadəcə constraint-dir. Class sistemindən gələn reflekslərdən
  qorunmaq lazımdır.

### 8. Approximation (~) — Derived Tipləri Qəbul Et
```go
type MyInt int
type Integer interface {
    int | int8 | int16 | int32 | int64
}
fmt.Println(Double(MyInt(1)))
// MyInt does not implement Integer

type ApproximatelyInt interface {
    ~int
}
```
- **Nədir:** `~` (tilde) işarəsi — "underlying type int olan HƏR ŞEY"
- **Necə işləyir:** `~int` type set-inə int ÖZÜ + underlying type-ı int olan
  bütün derived tiplər (MyInt) daxildir
- **Nəyə lazımdır:** derived tip (type MyInt int) adlı tipdən FƏRLİDİR —
  named-only constraint onu rədd edir; ~ bunu həll edir
- **Struct halı:**
```go
type Pointish interface {
    struct{ x, y int }
}
p := Point{1, 2}
Plot(p)  // Point does not implement Pointish
//        (possibly missing ~ for struct{x int; y int} in constraint Pointish)

type Pointish interface {
    ~struct{ x, y int }   // → Point artıq QƏBUL edilir
}
```
  Go özü "hint" verir: ~ yazmağı unutmusan!
- **Intish məşqi:**
```go
type Intish interface {
    ~int
}
func IsPositive[T Intish](v T) bool {
    return v > 0
}
```

### 9. Interface Literal — Constraint-i Yerində Yaz
```go
func Identity[T interface{}](v T) T { ... }  // = [T any]

func Stringify[T interface{ String() string }](s T) string {
    return s.String()
}

[T interface{ ~int }]
[T ~int]        // interface açar sözü ODA BİLƏR — yalnız 1 type element olanda
```
- **Nədir:** adlandırılmış constraint YOX — literal sintaksis
- **Qaydalar:**
  - `interface{}` = `any` (predeclared ad)
  - interface açar sözü yalnız TİK type element olanda ixtisar olunur:
    `[T ~int]` OK, `[T ~int; ~float64]` SYNTAX ERROR
  - metod elementi ilə ixtisar OLMAZ: `[T String() string]` SYNTAX ERROR
  - ixtisar yalnız constraint LITERAL-də; adlandırılmış constraint
    tərifində `type Intish ~int` SYNTAX ERROR
- **Nəyə lazımdır:** bəzən ad vermək əskidir, literal daha aydındır

### 10. Type Parameter-ə Constraint Daxilində Müraciət
```go
func Contains[T interface{ Equal(T) bool }](s []T, v T) bool {
    ...
}
```
- **Nədir:** constraint-in metod imzası T-yə İSTİNAD edir
- **Necə işləyir:** `Equal(T) bool` — hər konkret T üçün Equal konkret tip
  qəbul edir
- **MÜMKÜNSÜZ adlandırılmış variant:**
```go
type Equaler interface {
    Equal(???) bool  // T burada mövcud deyil — bilmirik nə yazaq!
}
```
  Interface literal constraint daxilində T-yə istinad etməyin YEGƏNƏ yoludur.
  (Generic interface — növbəti fəsil.)
- **Greater məşqi:**
```go
func IsGreater[T interface{ Greater(T) bool }](x, y T) bool {
    return x.Greater(y)
}
```

## Əsas terminlər
- Constraint (məhdudiyyət) — type parameter üçün icazəli tiplər dəsti
- Type set (tip dəsti) — constraint-i qane edən bütün tiplər
- Union (birləşmə) — `|` ilə ayrılan alternatives type elementləri
- Intersection (kəsişmə) — çoxsətirli interface — hamısını tələb edir
- Type approximation (tip yaxınlaşması) — `~T`: underlying type T olanlar
- Basic interface (əsas interfeys) — yalnız metod elementli interface
- Type element (tip elementi) — interface daxilində tip (metod deyil)
- Derived type (törəmə tip) — `type MyInt int` kimi mövcud tipdən yaradılan
- Interface literal (interfeys literali) — ad verilməmiş inline interface

## Praktik nəticə

1. **Constraint seçim qaydası:** operator lazımdırsa → number union-ları
   (~int, ~float64); metod lazımdırsa → basic interface; hər ikisi üçün
   kompozisiya (Number = Integer | Float | Complex).
2. **~ yazmaq qərarı:** öz tipini (type MyInt int) başqalarının da
   istifadə edəcəyi constraint-də istəyirsənsə ~ mütləqdir; named-only
   type set-lər derived tipləri rədd edir.
3. **Adlandırma vs literal:** sadə, təkrar istifadə olunan → adlandırılmış
   constraint; T-ə istinad edən (Equal(T) bool) → yalnız literal mümkündür.
4. **Xəta tanınması:** "possibly missing ~" = derived tip problemi;
   "empty type set" = kəsişmə boşdur; "interface contains type constraints" =
   constraint-i adi tip kimi istifadə etmək istəyirsən — QADAĞANDIR.

## Mənbə
Pages: 44-65 (PDF 45-66)
