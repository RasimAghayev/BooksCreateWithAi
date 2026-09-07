# Chapter 15 — Changes in Recent Go Versions (Son Go Versiyalarındakı Dəyişikliklər)

## Bu chapter nədən bəhs edir?

rand.Seed-in 1.20-dən etibarən lazımsızlığı (New/NewSource替代i), Go 1.21 yenilikləri
(sync.OnceFunc, clear), Go 1.22 yenilikləri (loop dəyişəni paylaşımının ləğvi, slices
funksiyalarının silinən elementləri sıfırlaması, slices.Concat, integer range, math/rand/v2
generic rand.N) və ServeMux-un metod+wildcard dəstəyi.

## Əsas fikirlər

### 1. rand.Seed() Artıq Lazım Deyil (Go 1.20+)
**Nə dəyişdi:** rand.Seed(randomDəyər) çağırmaya ehtiyac YOXDUR — global generator
avtomatik random başlanğıclanır. Mövcud kodu POZMUR.

**Xüsusi ardıcıllıq lazımdırsa (test üçün):**
```go
src := rand.NewSource(seed)     // seed-dən asılı mənbə
r := rand.New(src)              // yeni *rand.Rand
for i := 0; i < times; i++ {
    fmt.Println(r.Uint64())     // eyni seed → EYNİ ardıcıllıq
}
```
randSeed.go 1 ilə 2 dəfə işə salındı → hər ikisində eyni 2 ədəd.

### 2. Go 1.21 — sync.OnceFunc
**Nədir:** `func OnceFunc(f func()) func()` — funksiyanı YALNIZ BİR DƏFƏ icra edən
wrapper qaytarır. İlk çağırış f-i işlədir; sonrakılar NO-OP.

**Kitabdan kod nümunəsi:**
```go
var x = 0
func initializeValue() { x = 5 }

function := sync.OnceFunc(initializeValue)
for i := 0; i < 10; i++ {
    go function()        // 10 goroutine çağırır — amma initializeValue 1 DƏFƏ
}
time.Sleep(time.Second)
// x = 5 (yalnız bir dəfə initialize edildi)
```
**İstifadə halları:** dəyişən/bağlantı/fayl inicializasiyasının thread-safe tək icrası —
race condition olmadan lazy init.

### 3. Go 1.21 — clear Funksiyası
**Kitabdan kod nümunəsi:**
```go
m := map[string]int{"One": 1}
m["Two"] = 2
clear(m)      // map: BUTÜN elementlər SİLİNİR → map[]

s := []int{0, 1, 2, 3, 4}   // len 5, cap 10
clear(s)      // slice: elementlər ZERO VALUE → [0 0 0 0 0], len/cap EYNİ qalır
```
**Vacib fərq:** map-də silmə, slice-da SIFIRLAMA (uzunluq/tutum saxlanılır).

### 4. Go 1.22 — Loop Dəyişəni Paylaşılmır
**Ən böyük dəyişiklik:** hər iterasiya öz dəyişən NÜSXƏSİ alır.

**Nə deməkdir:**
```go
values := []int{1, 2, 3, 4, 5}
for _, val := range values {
    go func() {
        fmt.Printf("%d ", val)    // GO 1.22-DƏN ƏVVƏL: RACE — hamısı son dəyəri görə bilərdi
    }()                            // 1.22+: hər goroutine ÖZ val nüsxəsi
}
```
Ch8-dəki goClosure.go closure tələsi bu versiyadan etibarən avtomatik həll olunur.
**Müəllim qeydi:** aydın kod üçün parametr ötürmə hələ də good practice-dir.

### 5. Go 1.22 — slices Paketinin Dəyişiklikləri
**Silən funksiyalar silinən elementləri SIFIRLAYIR** (Delete, DeleteFunc, Compact,
CompactFunc, Replace — yeni len ilə köhnə len arası zero value):
```go
v1 := []int{-1, 1, 2, 3, 4}
v2 := slices.Delete(v1, 1, 3)
// v1: [-1 3 4 0 0]  ← silinən yerlər sıfırlandı (memory leak qorunması!)
// v2: [-1 3 4]      ← qaytarılan slice
```
**Niyə vacib:** silinən elementlər underlying array-də yaşamır — pointer saxlayan
tiplərdə GC-yə maneə yoxdur (Ch14 leak mövzusu ilə əlaqə).

**Yeni: slices.Concat:**
```go
conCat := slices.Concat(s1, s2, s3)   // [1 2 -1 -2 10 20]
```

### 6. Go 1.22 — range over Integers
```go
for x := range 5 {
    fmt.Print(" ", x)     // 0 1 2 3 4
}
```
Klassik `for i := 0; i < n; i++` üçün qısa alternativ.

### 7. Go 1.22 — math/rand/v2
**Əsas dəyişikliklər:**
- `rand.Read()` DEPRECATED — öz Read-imizi Uint64 ilə qura bilərik
- **Generic funksiyalar:** `rand.N()` — hər hansı integer tipi ilə işləyir

**Kitabdan kod nümunəsi:**
```go
import "math/rand/v2"

// Custom Read — deprecated metodu əvəz edir:
func Read(p []byte) (n int, err error) {
    for i := 0; i < len(p); {
        val := rand.Uint64()
        for j := 0; j < 8 && i < len(p); j++ {
            p[i] = byte(val & 0xff)
            val >>= 8
            i++
        }
    }
    return len(p), nil
}

// Generic rand.N — int VƏ uint (VƏ time.Duration!) ilə:
var max int = 100
n := rand.N(max)            // int qaytarır
var uMax uint = 100
uN := rand.N(uMax)          // uint qaytarır
```
`rand.N(max)` — parametrin tipi nəticənin tipini müəyyən edir; [0, max) intervalı.

### 8. Go 1.22 — ServeMux Təkmilləşməsi
net/http.ServeMux pattern-ləri artıq **metodlar** (GET /x) və **wildcard-lar**
(/items/{id}) qəbul edir — üçüncü tərəf routerlərə ehtiyac azalır (kitabın Ch9/Ch11
mövzularına təsir edən modernləşmə).

## Əsas terminlələr
- rand.New/NewSource — seed-ə nəzarət edən ayrıca generator
- sync.OnceFunc — tək icra garantili funksiya wrapper-i
- clear — map silmə / slice sıfırlama built-in
- Loop Variable Semantics (1.22) — hər iterasiya müstəqil dəyişən
- slices.Concat — çoxslice birləşdirmə
- Zeroing on Delete — silinən elementlərin sıfırlanması (leak qorunması)
- range over int — tam ədəd üzərində iterasiya
- math/rand/v2 — generic rand.N, Read deprecated
- ServeMux Method+Wildcard — standart router-in gücləndirilməsi

## Praktik nətidə

(1) Yeni kodda rand.Seed(seed) YAZMA — test üçün rand.New(NewSource(seed)). (2) Lazy
init = sync.OnceFunc — init race-dən qorunur. (3) clear: map üçün təmizləmə, slice üçün
sıfırlama — fərqli semantika, len/cap dəyişmir. (4) Go 1.22-də loop dəyişəni goroutine-ə
təhlükəsizdir — amma oxunaqlıq üçün parametr ötürməyi saxla. (5) slices.Delete family
artıq aralı elementləri sıfırlayır — pointer saxlayan slice-larda leak riski azalır.
(6) rand.N() generic — int/uint/duration hamısı. (7) Yeni layihələrdə metod+wildcard
ServeMux pattern-lərini nəzərdən keçir — router asılılığı azala bilər.

## Mənbə
Pages: 661-671 (PDF 692-703)
