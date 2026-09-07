# Chapter 5 — Types (Tiplər)

## Bu fəsil nədən bəhs edir?

Generic TİPLƏR: named basic tiplər (Age/HeightCM — məna qarışıqlığının
qarşısını almaq), generic basic tipin QADAĞASI (`type MyT[T any] T` — RHS-də
type parameter olmaz), generic slice (Bunch — instantiate olunmadan YOXDUR,
"there are no generic types in a running program"), []any ilə Bunch[E]-nin
FƏRQİ (heterogen vs homogen), generic map (Index[K comparable, V any] — map
açarı comparable TƏLƏB EDİR), generic struct (NamedThing[T], self-reference
rekursiyası mümkün deyil — daxili instantiate tələb olunur), generic tiplərin
METODLARI (receiver-da E saxlanılır; amma metodun ÖZ type parameter-i ola
BİLMƏZ), generic interface (Equaler[T any] — previous chapter probleminin
həlli), generic channel (myChan[E]).

## Əsas fikirlər

### 1. Named Basic Tiplər — Məna Mismatch-inin Qarşısı
```go
type (
    Age      int
    HeightCM int
)
```
- **Nədir:** eyni underlying tipdən (int) fərqli ADLI tiplər
- **Nəyə lazımdır:** 30 sm boy + 180 yaş = type DEYİL, MƏNA qarışıqlığıdır —
  named tiplər bunu compile xətasına çevirir
- **QADAĞA:** generic basic tip YOXDUR:
```go
type MyT[T any] T
// cannot use a type parameter as RHS in type declaration
```
  T hələ bilinmir → yeni ad vermək mənasızdır. Composite tiplərdə isə (slice,
  map, struct) type parameter MÜMKÜNDÜR.

### 2. Generic Slice — Bunch[E]
```go
type Bunch[E any] []E
b := Bunch[int]{1, 2, 3}
b = append(b, "hello")
// cannot use "hello" as int value in argument to append
```
- **Vacib fərq — []any DEYİL:**
```go
type Stuff []any
s := Stuff{1, 2, 3}
s = append(s, "hello")  // OK — heterogen!
```
  - `[]any` — fərqli konkret tiplər qarışıq (heterogen collection)
  - `Bunch[E]` — instantiate olunmuş HAMISI EYNSİ tip (homogen);
    runtime-da "sadəcə Bunch" yoxdur — yalnız Bunch[int], Bunch[string]...
- **"There are no generic types" prinsipi:** compile-da hər parameterised tip
  konkret tipə instantiate olunmalıdır; runtime-da yalnız specific tiplər var.
- Interface tipinə instantiate: `var b Bunch[error]` — mümkündür (error
  constraint-i qane edir).

### 3. Generic Map — Açar comparable Tələbi
```go
type Catalog[V any] map[string]V
cat := Catalog[int]{}
fmt.Println(cat["bogus"])  // 0

type Index[K, V any] map[K]V
// invalid map key type K (missing comparable constraint)

type Index[K comparable, V any] map[K]V
age := Index[string, int]{}
```
- **Nədir:** map tipləri də generic ola bilər — V (value) hər hansı, K (key)
  isə comparable OLMALIDIR
- **Səbəb:** map eyni açarını tapmaq üçün == istifadə edir — K
  comparable deyilsə açar yoxlamaq mümkün deyil
- **Multi-instantiate:** `Index[string, int]` — kvadrat mötərizədə, vergüllə.

### 4. Generic Struct + Self-Reference
```go
type NamedThing[T any] struct {
    Name  string
    Thing T
}
n := NamedThing[float64]{"Latitude", 50.406}

var n NamedThing[NamedThing]
// cannot use generic type NamedThing[T any] without instantiation

p := NamedThing[NamedThing[float64]]{
    Name: "Position",
    Thing: NamedThing[float64]{"Latitude", 50.406},
}
// {Position {Latitude 50.406}}
```
- **Nədir:** struct sahəsi arbitrary tip T ola bilər
- **Self-reference qaydası:** daxili generic tip DƏ instantiate olunmalıdır —
  "hall of mirrors" rekursiyası mümkün deyil, çünku NamedThing tipli özü
  tip deyil, tip FABRİKİdir.

### 5. Generic Tipin Metodları
```go
type Bunch[E any] []E

func (b Bunch[E]) First() E {
    return b[0]
}
b := Bunch[string]{"a", "b", "c"}
fmt.Println(b.First())  // a
```
- **Nədir:** receiver `Bunch[E]` — E receiver-də saxlanılır, metod özü
  generic DƏQİL
- **Üstünlük:** First-i N tip üçün N dəfə YOX, bir dəfə yazırıq — bütün
  Bunch instantiate-lərində işləyir
- **QADAĞA — parameterised metod:**
```go
func (b Bunch[E]) PrintWith[T any](v T) { ... }
// methods cannot have type parameters
```
  Metodun öz type parameter-i OLMAZ — əlavə tip lazımdırsa, ayrıca generic
  FUNKSIYA yaz.
- **Empty məşqi:**
```go
type Sequence[E any] []E
func (s Sequence[E]) Empty() bool {
    return len(s) == 0
}
```

### 6. Generic Interface — Equaler[T]
```go
type Equaler[T any] interface {
    Equal(T) bool
}
```
- **Nədir:** interface-in ÖZÜ type parameter ala bilər
- **Nəyə lazımdır:** Chapter 3-dəki problem — `Contains[T interface{ Equal(T) bool }]`
  literal tələb edirdi, çünku ADLI interface T-yə istinad edə bilmirdi.
  İndi `Equaler[T]` adlandırılmış həlldir: "hər T üçün Equal(T) bool
  metoduna malik interface".

### 7. Generic Channel
```go
type myChan[E any] chan E
ch := make(myChan[error])
go func() {
    ch <- errors.New("oh no")
}()
fmt.Println(<-ch)  // oh no
```
- **Nədir:** channel element tipi də generic ola bilər
- **Nəyə lazımdır:** Drain[E], Merge[E] kimi funksiyalar element tipindən
  asılı olmayacaq əməliyyatlar üçün.

## Əsas terminlər
- Named type (adlı tip) — mövcud tipdən yaradılan yeni ad (məna ayırdetməsi)
- Underlying type (altlying/taban tip) — tipin əsas strukturu
- Heterogeneous collection (qarışıq kolleksiya) — []any: fərqli tiplər
- Homogeneous collection (vahid kolleksiya) — Bunch[int]: hamısı eyni tip
- Instantiation (instansiasiya) — generic tipin konkret tipə "çevrilməsi"
- Type factory (tip fabriki) — generic tip: adı tək, konkretləri çox

## Praktik nəticə

1. **Məna qarışıqlığı riski:** eyni underlying tipdən fərqli mənalı named
   tiplər yarat (Age, HeightCM) — compiler səni qoruyur.
2. **"Generic collection" termininə diqqət:** Bunch[E] heterogen DEYİL;
   heterogen lazımdırsa []any istifadə et (amma type assertion xərcləri ilə).
3. **Map açarı:** həmişə comparable constraint — `[K comparable, V any]`.
4. **Metodlar:** receiver-da E saxla (Bunch[E]), amma metodun yeni type
   parameter-i YAZA BILMƏZSƏN — çıxış yolu: səviyyəli generic funksiya.
5. **Nested generic:** daxili tipi həmişə instantiate et:
   NamedThing[NamedThing[float64]].

## Mənbə
Pages: 91-103 (PDF 92-104)
