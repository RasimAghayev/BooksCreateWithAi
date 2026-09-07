# Chapter 6 — Functions (Funksiyalar)

## Bu fəsil nədən bəhs edir?

Generic funksiyaların real tətbiqi — container tipləri (slice/map) üzərində
utility funksiyalar: Contains (comparable), Reverse, Sort (cmp.Ordered),
birinci dərəcəli funksiyaların (first-class functions) generic gücü:
mapFunc/keepFunc/reduceFunc tip adları, Map/Filter/Reduce funksional
triodu, type inference-in işə düşməsi (mapFunc[int] mismatch xətası),
generic keepFunc (IsEven — % operatoru constraints.Integer tələb edir),
Reduce-da any-nin niyə kifayət etməsi (+ operatoru BİZİM funksiyamızda
istifadə olunur — yazıla bilməzsə Reduce-ə ötürülə bilməz), konkurent
Map/Filter imkanları (embarrassingly parallel), FuncMap (dynamic dispatch)
və Compose (funksiya zənciri, 3 type parameter) məşqləri.

## Əsas fikirlər

### 1. Contains — comparable İlə Axtarış
```go
func Contains[E comparable](s []E, v E) bool {
    for _, vs := range s {
        if v == vs {
            return true
        }
    }
    return false
}
```
- **Nədir:** slice-də element axtarışı — generic olmadan hər tip üçün ayrıca
  yazılmalı idi
- **Necə işləyir:** E comparable → == zəmanətli; range + == müqayisəsi
- **Nəyə lazımdır:** generics-in əsas qazancı — utility paketləri (sonradan
  slices paketi)

### 2. Reverse və Sort
```go
func Reverse[S ~[]E, E any](s S) S {
    result := make(S, 0, len(s))
    for i := len(s) - 1; i >= 0; i-- {
        result = append(result, s[i])
    }
    return result
}

func Sort[S ~[]E, E cmp.Ordered](s S) S {
    result := make(S, len(s))
    copy(result, s)
    sort.Slice(result, func(i, j int) bool {
        return result[i] < result[j]
    })
    return result
}
```
- **Reverse:** müqayisə YOXDUR → any kifayətdir; ~[]E → derived slice qorunur
- **Sort:** < operatoru → cmp.Ordered; sort.Slice altında reflection var
  (yavaş) — generics ilə Quicksort birbaşa yazıla bilər; kitab variantı
  kopya üzərində işləyir (orijinal dəyişmir)

### 3. Map — Funksiyanı Kolleksiyaya Tətbiq Et
```go
type mapFunc[E any] func(E) E

func Map[S ~[]E, E any](s S, f mapFunc[E]) S {
    result := make(S, len(s))
    for i := range s {
        result[i] = f(s[i])
    }
    return result
}

s := []string{"a", "b", "c"}
fmt.Println(Map(s, strings.ToUpper))
// [A B C]

s := []int{1, 2, 3}
fmt.Println(Map(s, strings.ToUpper))
// type func(s string) string of strings.ToUpper does not
// match inferred type mapFunc[int] for mapFunc[E]
```
- **Nədir:** hər elementə funksiya tətbiqi — nəticə eyni tipdə transformasiya
  olunmuş slice
- **Type inference dərsi:** Go []int-dən E=int çıxarır → 2-ci parametr
  func(int) int olmalıdır; strings.ToUpper func(string) string → MISMATCH.
  Go inference-i parametrlər üzrə paylaşır — bir parametr tipi digərinə
  təsir edir.
- **İmza aydınlığı:** mapFunc[E] ADI imzanı oxunaqlı edir — funksiya tipləri
  üçün type alias idiomatikdir.

### 4. Filter — Seçmə
```go
type keepFunc[E any] func(E) bool

func Filter[S ~[]E, E any](s S, f keepFunc[E]) S {
    result := S{}
    for _, v := range s {
        if f(v) {
            result = append(result, v)
        }
    }
    return result
}

s := []int{1, 2, 3, 4}
fmt.Println(Filter(s, func(v int) bool {
    return v%2 == 0
}))
// [2 4]
```
- **Nədir:** keep funksiyası true qaytaran elementlər saxlanılır
- **Generic keepFunc — maraqlı hal:**
```go
func IsEven[T any](v T) bool {
    return v%2 == 0
}
// % yalnız tam ədədlərdə → any İŞLƏMƏZ

func IsEven[T constraints.Integer](v T) bool {
    return v%2 == 0
}
fmt.Println(Filter(s, IsEven[int]))   // [2 4] — EXPLICIT instantiate!
```
  Constraint operatora görə seçilir; generic funksiyanı Filter-ə ötürəndə
  `[int]` ilə instantiate etmək lazımdır.

### 5. Reduce — Birləşdirmə (yığma)
```go
type reduceFunc[E any] func(cur, next E) E

func Reduce[E any](s []E, init E, f reduceFunc[E]) E {
    cur := init
    for _, v := range s {
        cur = f(cur, v)
    }
    return cur
}

s := []int{1, 2, 3, 4}
sum := Reduce(s, 0, func(cur, next int) int { return cur + next })  // 10
p := Reduce(s, 1, func(cur, next int) int { return cur * next })     // 24
j := Reduce([]string{"a", "b", "c"}, "", func(c, n string) string { return c + n })  // abc
```
- **Nədir:** ardıcıllığı TƏK qiymətə yığma (fold/inject kimi də tanınır)
- **Necə işləyir:** cur (cari nəticə) + next (element) → yeni cur; init
  başlanğıc qiymətdir (cəm üçün 0, hasil üçün 1, konkatenasiya üçün "")
- **Niyə any kifayətdir?** + operatoru Reduce-un ÖZÜNDƏ deyil, BİZİM
  reduceFunc-da istifadə olunur. Struct üçün `cur + next` yazılmaz (compile
  xətası) → belə funksiya ÜMUMİYYƏTLƏ yazıla bilməz → Reduce-ə ötürülmə
  problemi yaranmır. Operator funksiya daxilində olduqda constraint
  istifadəçinin öz funksiyasına aiddir.

### 6. Konkurrent Map/Filter Perspektivi
- **Embarrassingly parallel (tam paralel):** elementlər müstəqildir →
  paralel filtrləmə xətti sürətlənmə verir
- **Nümunə:** 1000 URL filtrləməsi — serial ~250 saniyə (hər biri 250ms)
  vs konkurrent ~250ms (hamısı eyni anda)
- **Nəyə lazımdır:** concurrency kodu mürəkkəbdir, amma generic olduğu üçün
  CƏMİ BİR DƏFƏ yazılır — bəzi üçüncü tərəf paketlər konkurent Filter təklif
  edir.

### 7. FuncMap — Dynamic Dispatch Məşqi
```go
type FuncMap[T, U any] map[string]func(T) U

func (fm FuncMap[T, U]) Apply(name string, val T) U {
    return fm[name](val)
}

fm := funcmap.FuncMap[int, int]{
    "double": func(i int) int { return i * 2 },
    "addOne": func(i int) int { return i + 1 },
}
fmt.Println(fm.Apply("double", 2))  // 4
```
- **Nədir:** runtime-da string adı ilə funksiya çağırma — map[string]func(T) U
- **Qeyd:** metod FuncMap[T, U] receiver-i üzərində yazılır (düz
  "FuncMap" üzərində YOX — instantiate edilməmiş generic tipə metod
  yazılmaz)

### 8. Compose — Funksiya Zənciri Məşqi
```go
func Compose[T, U, V any](f func(U) T, g func(V) U, v V) T {
    return f(g(v))
}

fmt.Println(compose.Compose(double, addOne, 1))  // double(addOne(1)) = 4
fmt.Println(compose.Compose(utf8.RuneCountInString, strings.ToUpper, "HeLlO, wOrLd"))
// RuneCount(ToUpper("HeLlO, wOrLd")) = 12
```
- **Nədir:** iki funksiyanı zəncirə salma — g tətbiq olunur, sonra f
- **Sıra qaydası:** "last named, first applied" — Compose(f, g, v) = f(g(v));
  f = outer, g = inner
- **3 type parameter:** V (g-nin girişi), U (g-nin nəticəsi = f-in girişi),
  T (f-in nəticəsi). f-in parametri g-nin qaytarığı ilə UYĞUN olmalıdır.

## Əsas terminlər
- First-class function (birinci dərəcəli funksiya) — dəyər kimi ötürülən/qaytarılan funksiya
- Map/Filter/Reduce — transformasiya/seçmə/yığma triodu
- Fold/Inject — Reduce-un digər dillərdəki adları
- Type inference (tip çıxarımı) — Go-nun E-ni parametrlərdən tapması
- Embarrassingly parallel (tam paralel) — müstəqil tapşırıqların xətti sürətlənməsi
- Dynamic dispatch (dinamik yönləndirmə) — runtime-da funksiya seçimi
- Function composition (funksiya birləşməsi) — f(g(v)) zənciri

## Praktik nəticə

1. **Utility funksiyalar üçün hazır idiomlar:** Contains → [E comparable];
   Sort/Greatest → [S ~[]E, E cmp.Ordered]; Map/Filter → [S ~[]E, E any] +
   mapFunc/keepFunc tip adları.
2. **Funksiya tipinə AD ver:** type mapFunc[E any] func(E) E — imzalar
   oxunaqlı olur, səhvlər aydın görünür.
3. **Operator funksiyada olduqda** constraint-i istifadəçinin funksiyası
   müəyyən edir (Reduce-da any; IsEven-də constraints.Integer).
4. **Generic funksiyanı arqument kimi ötürəndə** instantiate et: IsEven[int].
5. **Funksional üslub:** Map+Filter+Reduce qısa və ifadəlidir; performans
   kritik olduqda konkurrent variant bir dəfə yazılıb hər yerdə istifadə
   oluna bilər.

## Mənbə
Pages: 104-124 (PDF 105-125)
