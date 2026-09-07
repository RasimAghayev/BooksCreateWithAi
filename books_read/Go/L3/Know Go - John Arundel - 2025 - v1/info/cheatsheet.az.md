# Know Go — Cheat Sheet (Azərbaycanca)

## 1. Constraint seçim cədvəli (operatora görə)

| İstədiyin operator | Constraint | Nümunə |
|---|---|---|
| `+ - * /` | `Number` (Integer \| Float \| Complex) | `func AddAnything[T Number](x, y T) T` |
| `< > <= >=`, `min`/`max` | `cmp.Ordered` (import "cmp") | `func Greater[T cmp.Ordered](x, y T) T` |
| `==` `!=` | `comparable` (predeclared) | `func Equal[T comparable](x, y T) bool` |
| `%` (tam ədədlər) | `constraints.Integer` / öz Integer-in | `func IsEven[T Integer](v T) bool` |
| Heç biri | `any` | `func PrintAnything[T any](v T)` |

```go
type Integer interface {
    ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 |
        ~uint16 | ~uint32 | ~uint64
}
type Float interface { ~float32 | ~float64 }
type Complex interface { ~complex64 | ~complex128 }
type Number interface { Integer | Float | Complex }
type Real interface { Integer | Float }
type Ordered interface { Integer | Float | ~string }   // = cmp.Ordered
```

## 2. Əsas generic sintaksis formaları

```go
// Funksiya:
func F[T any](v T) T { return v }              // 1 parametr
func F[T, U any](x T, y U) (T, U) { ... }       // 2 parametr (AYRI tiplər!)
func F[T interface{ Equal(T) bool }](s []T, v T) bool  // literal + T istinadı

// Tip:
type Bunch[E any] []E
type Index[K comparable, V any] map[K]V          // map açarı comparable!
type NamedThing[T any] struct { Name string; Thing T }
type Equaler[T any] interface { Equal(T) bool }  // generic interface
type myChan[E any] chan E

// Metodlar (receiver-da E saxlanır; metodun ÖZ type parametri OLMAZ):
func (b Bunch[E]) First() E { return b[0] }

// Constraint ixtisarı (yalnız TİK type element olanda):
func Increment[T ~int](v T) T { return v + 1 }
```

## 3. Derived tiplər üçün ~ idiomları

```go
// Slice qaytaran generic — DERİV SLİCE İTİRMƏSİN:
func Clone[S ~[]E, E any](s S) S            // slices.Clone imzası belədir
func Identity[S ~[]E, E any](s S) S
func Reverse[S ~[]E, E any](s S) S

// Map-lar üçün:
func Merge[M ~map[K]V, K comparable, V any](ms ...M) M

// Yalnız ~int named-deyil, derived də qəbul etsin:
type Intish interface { ~int }
```

## 4. slices paketi (əsas komandalar)

```go
slices.Equal(s1, s2)                    // == YOXDUR — bunu işləd
slices.EqualFunc(s1, s2, myEq)          // fərqli element tipləri belə OK
slices.Compare(s1, s2)                  // -1/0/+1
slices.Index(s, v)                      // tapılmadısa -1
slices.Contains(s, v)                   // bool
slices.Max(s) / slices.Min(s)
slices.Insert(s, i, v...)               // i indeksinə daxil et
slices.Delete(s, i, j)                  // [i, j) — j DAHIL DEYIL
slices.Clone(s)                          // make+copy əvəzinin 1 sətiri
slices.Compact(s)                        // ARDICIL dublikatlar (uniq)
slices.Grow(s, n)                        // kapasiteni bir dəfə böyüt
slices.Clip(s)                           // kapasiteyi length-ə kəs
slices.Sort(s)                           // in-place, cmp.Ordered
slices.SortFunc(s, cmp.Compare)          // custom müqayisə
slices.SortStableFunc(s, cmp)            // stabil
slices.IsSorted(s)
slices.Reverse(s)                        // in-place
slices.Replace(s, i, j, v...)            // [i,j) əvəz et
slices.BinarySearch(sorted, v)           // (index, found); tapılmasa insertion point
slices.Collect(seq)                      // iterator → slice
slices.All(s) / slices.Values(s)         // slice → iterator
```

## 5. maps paketi

```go
maps.Equal(m1, m2)                       // açar həmişə ==
maps.EqualFunc(m1, m2, valueEq)
maps.DeleteFunc(m, func(k, v) bool { return ... })
maps.Clear(m)
maps.Clone(m)
maps.Copy(to, from)                      // SIRA: destination birinci!
```

## 6. cmp paketi

```go
import "cmp"
cmp.Or(userValue, "default")             // ilk non-zero
cmp.Compare(x, y)                         // -1/0/+1 (SortFunc ilə)
// cmp.Ordered — constraint (yuxarıda)
```

## 7. Set implementation (Ch7-8)

```go
type Set[E comparable] map[E]struct{}

func NewSet[E comparable](vals ...E) Set[E]
func (s Set[E]) Add(vals ...E)                    // s[v] = struct{}{}
func (s Set[E]) Contains(v E) bool                // _, ok := s[v]
func (s Set[E]) All() []E
func (s Set[E]) Union(s2 Set[E]) Set[E]
func (s Set[E]) Intersection(s2 Set[E]) Set[E]

// Concurrency-safe variant:
type SetC[E comparable] struct {
    mutex *sync.RWMutex                            // POINTER — kopya qorunur!
    data  map[E]struct{}
}
// Add → mutex.Lock + defer Unlock
// All/Contains → mutex.RLock + defer RUnlock
```

## 8. Concurrency (Ch8)

```go
// Mutex idiomu:
s.mutex.Lock()
defer s.mutex.Unlock()          // yazı
s.mutex.RLock()
defer s.mutex.RUnlock()         // oxu (paylaşılan)

// Atomic sayaç (mutex-siz):
type Channel[T any] struct {
    ch              chan T
    sends, receives atomic.Uint64
}
c.sends.Add(1)                  // atomik artım
c.sends.Load()                  // atomik oxu

// Test: həmişə race detector ilə:
go test -race
go run -race main.go
```

## 9. Map/Filter/Reduce (Ch6)

```go
type mapFunc[E any] func(E) E
type keepFunc[E any] func(E) bool
type reduceFunc[E any] func(cur, next E) E

func Map[S ~[]E, E any](s S, f mapFunc[E]) S
func Filter[S ~[]E, E any](s S, f keepFunc[E]) S
func Reduce[E any](s []E, init E, f reduceFunc[E]) E

// Nümunələr:
Map([]string{"a"}, strings.ToUpper)
Filter([]int{1,2,3,4}, func(v int) bool { return v%2 == 0 })  // [2 4]
Reduce([]int{1,2,3,4}, 0, func(c, n int) int { return c + n })  // 10

// Compose:
func Compose[T, U, V any](f func(U) T, g func(V) U, v V) T { return f(g(v)) }
```

## 10. İteratorlar (Ch11, Go 1.23+)

```go
import "iter"

func Items() iter.Seq[Item] {           // 1 dəyər
    return func(yield func(Item) bool) {
        for _, v := range items {
            if !yield(v) { return }     // FALSE → dayan + təmizlə
        }
    }
}
func Items2() iter.Seq2[int, Item]       // 2 dəyər (indeks/error üçün)

for v := range Items() { ... }
for i, v := range Items2() { ... }

// Xətalarla:
func Lines(file string) iter.Seq2[string, error]

// Kompozisiya (sonsuz + filtrlər):
for p := range Primes(Integers()) { ... }

// Standart: slices.All/Values, maps.All/Keys/Values, slices.Collect
```

## 11. Səhv mesajları tərcüməsi

| Xəta mesajı | Mənası | Həll |
|---|---|---|
| `operator + not defined on T (constrained by any)` | any operator zəmanət etmir | Number/Ordered kimi dar constraint |
| `cannot infer T` | nəticə T, daxil olmadı | `Something[int]()` explicit |
| `MyInt does not implement Integer` | derived tip named set-də deyil | constraint-ə ~ əlavə et |
| `possibly missing ~ for ...` | Go özü işarə edir | ~ yaz |
| `interface contains type constraints` | constraint-i ADİ TİP kimi istifadə | yalnız [T C] mövqeyində |
| `empty type set` | kəsişmə boşdur (int ∩ string) | constraint elementlərini yoxla |
| `methods cannot have type parameters` | metodda [T] yazmaq olmaz | generic FUNKSIYA yaz |
| `incomparable types in type set` | == üçün any genişdir | comparable |
| `concurrent map iteration and map write` | data race! | mutex və ya -race ilə tut |
| `range function continued iteration after exit` | false-dan sonra yield | if !yield → return |
| `T is not comparable with >` | Complex ordered deyil | Real/Ordered (Complex-siz) |
