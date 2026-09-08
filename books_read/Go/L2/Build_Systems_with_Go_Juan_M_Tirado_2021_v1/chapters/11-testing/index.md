# Chapter 11 — Testing (səh. 221-251)

## Bu fəsil nədən bəhs edir?

Go-nun inteqrasiya olunmuş test alətləri: `TestXXX` funksiyaları, `go test`,
subtests (`t.Run`), test skipping, Example funksiyaları (sənədləşmə + output
yoxlaması), benchmarking (`BenchmarkXXX`, `RunParallel`), coverage
(`-cover`, `go tool cover`) və profiling (`go tool pprof`).

## Əsas fikirlər

### 1. Əsas testlər
**Nədir:** İmzası `TestXXX(t *testing.T)` olan istənilən funksiya test kimi
işə düşür. Uğursuzluq `t.Error(...)` ilə bildirilir.

**Kitabdan kod nümunəsi:**
```go
package example_01

import "testing"

func TestMe(t *testing.T) {
    r := 2 + 2
    if r != 4 {
        t.Error("expected 2 got", r)
    }
}
```

**İcra:**
```bash
>>> go test          # PASS + vaxt
>>> go test -v       # hər test ayrıca: === RUN TestMe / — PASS
# uğursuz halda: FAIL: TestMe ... expected 2 got 3
```

**Funksiya testi (bug tutma nümunəsi):**
```go
func MovieQuote(movie Movie) Quote {
    switch movie {
    case Conan:
        return Crush
    case Terminator2:
        return T1000
    default:
        return Unknown     // Predator üçün quote YOXDUR → test bunu tutur
    }
}

func TestMovieQuote(t *testing.T) {
    movies := []Movie{Conan, Predator, Terminator2}
    for _, m := range movies {
        if q := MovieQuote(m); q == Unknown {
            t.Error("unknown quote for movie", m)
        }
    }
}
```

### 2. Test skipping (t.Skip)
**Nədir:** Vaxt məhdudiyyətli testlərdə uzun işlərin ötürülməsi.
`testing.Short()` → `-short` flag-i qoyulubsa true.

```go
func TestSum(t *testing.T) {
    var i int64
    done := make(chan bool)
    for i = 1000; i < math.MaxInt64; i += 100000 {
        go Sum(i, done)
        timeout := time.NewTimer(time.Millisecond)
        select {
        case <-timeout.C:
            t.Skip(fmt.Sprintf("%d took longer than 1 millisecond", i))
        case <-done:
        }
    }
}
// go test -v → SKIP: TestSum (0.02s)
```

### 3. Subtests (t.Run)
**Nədir:** Test daxilində testlər; `key=value` adlandırma ilə süzgəclənə
qruplar qurmaq mümkündür.

```go
func TestFoo(t *testing.T) {
    t.Run("A=1", func(t *testing.T) { /* ... */ })
    t.Run("A=2", func(t *testing.T) { /* ... */ })
    t.Run("B=1", func(t *testing.T) { /* ... */ })
}
```

**Praktik nümunə — XML/JSON encode/read matrisi:**
```go
// Test olunan kod: User XML/JSON formatda fayla yazılır/oxunur
type Encoding int
const (
    XML Encoding = iota
    JSON
)

type User struct {
    UserId string `xml:"id" json:"userId"`
    Email  string `xml:"email" json:"email"`
    Score  int    `xml:"score" json:"score"`
}

// --- test faylı ---
func testWriteXML(t *testing.T) {
    tmpDir := os.TempDir()
    for _, u := range Users {
        f := tmpDir + u.UserId + ".xml"
        if err := u.ToEncodedFile(XML, f); err != nil {
            t.Error(err)
        }
    }
}

func testReadXML(t *testing.T) {
    tmpDir := os.TempDir()
    for _, u := range Users {
        f := tmpDir + "/" + u.UserId + ".xml"
        newUser := User{}
        if err := newUser.FromEncodedFile(XML, f); err != nil {
            t.Error(err)
        }
        if !newUser.Equal(u) {
            t.Error(fmt.Sprintf("found %v, expected %v", newUser, u))
        }
    }
}

func testXML(t *testing.T) {
    t.Run("Action=Write", testWriteXML)
    t.Run("Action=Read", testReadXML)
    // cleanup — tmp faylları sil
    for _, u := range Users {
        _ = os.Remove(os.TempDir() + "/" + u.UserId + ".xml")
    }
}

func testJSON(t *testing.T) {
    t.Run("Action=Write", testWriteJSON)
    t.Run("Action=Read", testReadJSON)
    // cleanup ...
}

func TestEncoding(t *testing.T) {     // yalnız BÜYÜK hərflə başlayanlar testdir
    t.Run("Encoding=XML", testXML)
    t.Run("Encoding=JSON", testJSON)
}

func TestMain(m *testing.M) {         // setup + icra + exit
    UserA := User{"UserA", "usera@email.org", 42}
    UserB := User{"UserB", "userb@email.org", 333}
    Users = []User{UserA, UserB}
    os.Exit(m.Run())
}
```

**Subtest süzgəcləri (-run + regex):**
```bash
>>> go test -v -run /Encoding=JSON       # yalnız JSON qolu
>>> go test -v -run /./Action=Write      # bütün formatlarda Write
```
- `testWriteXML` (kiçik hərf) tək başına işə düşmür — yalnız TestEncoding
  qrupundan çağrılır

### 4. Example funksiyaları (sənədləşmə testləri)
**Nədir:** `Example` prefiksli funksiyalar — çıxışı `// Output:` şərhi ilə
müqayisə olunur və sənədləşməyə düşür (godoc).

**Adlandırma qaydası:**
```go
func Example()        // paket üçün
func ExampleF()       // F funksiyası üçün
func ExampleT()       // T tipi üçün
func ExampleT_M()     // T tipinin M metodu üçün
func ExampleT_M_suffix()  // bir neçə variant varsa
```

**Kitabdan kod nümunəsi:**
```go
func ExampleUser() {
    j := User{"John", nil}
    m := User{"Mary", []User{j}}
    fmt.Println(m)
    // Output:
    // {Mary [{John []}]}
}

func ExampleUser_GetUserId() {
    u := User{"John", nil}
    fmt.Println(u.GetUserId())
    // Output:
    // JOHN
}

func ExampleCommonFriend() {
    a := User{"a", nil}
    b := User{"b", []User{a}}
    c := User{"c", []User{a, b}}
    fmt.Println(CommonFriend(&b, &c))
    // Output:
    // &{a []}
}
```
- Sıra zəmanəti yoxdursa `// Unordered output:` istifadə edin

### 5. Benchmarking
**Nədir:** `BenchmarkXXX(b *testing.B)` — `b.N` dəfə təkrarlanaraq performans
ölçülür. `-bench` flag-i ilə işə düşür; Go `b.N`-i avtomatik kalibrləyir.

```go
func Sum(n int64) int64 {
    var result int64 = 0
    var i int64
    for i = 0; i < n; i++ {
        result = result + i
    }
    return result
}

func BenchmarkSum(b *testing.B) {
    fmt.Println("b.N:", b.N)
    for i := 0; i < b.N; i++ {
        Sum(1000000)
    }
}
```

```bash
>>> go test -v -bench .
BenchmarkSum-16    3817    265858 ns/op
# -16 → 16 goroutine; 3817 → iterasiya sayı; 265858 ns/op → vaxt/iterasiya
```

**Paralel benchmark (RunParallel):**
```go
func BenchmarkSumParallel(b *testing.B) {
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            Sum(1000000)
        }
    })
}
```

```bash
>>> go test -v -bench . -cpu 1,2,4,8,16,32
BenchmarkSumParallel             4324    249886 ns/op
BenchmarkSumParallel-2          9462    127147 ns/op
BenchmarkSumParallel-4         18202     66514 ns/op
BenchmarkSumParallel-8         31191     33927 ns/op
# -cpu → goroutine sayı dəyişəndə ns/op azalır (paralellik qazancı görünür)
```

### 6. Coverage (kod örtüyü)
**Nədir:** Kodun neçə faizi test olunub. `-cover` ilə ölçülür.

```go
// test olunan kod: 8 budaqlı Periods funksiyası
func Periods(year int) string { /* switch ... */ }

// test: yalnız bir budaq (Periods(333)) sınanır
func TestOptions(t *testing.T) {
    Periods(333)
}
```

```bash
>>> go test -v -cover .
coverage: 22.2% of statements
```

**Dərin analiz alətləri:**
```bash
>>> go test -coverprofile=cover.out .
>>> go tool cover -func=cover.out   # funksiya-funksiya örtük
>>> go tool cover -html=cover.out   # brauzerdə kod bölgələri (rəngli)
```

**-covermode rejimləri:**
- `set` → icra olundumu (bool)
- `count` → neçə dəfə icra olundu
- `atomic` → count kimi, multithread-safe

### 7. Profiling (runtime/pprof)
**Nədir:** CPU/ yaddaş istifadəsinin səviyyə-səviyyə analizi. Benchmark ilə
birlikdə işlədilir.

**Profil fayllarının yığılması:**
```bash
>>> go test -bench=. -benchmem -memprofile mem.out -cpuprofile cpu.out
BenchmarkGraph-16   1365   827536 ns/op   411520 B/op   901 allocs/op
```

**Vizual hesabatlar:**
```bash
>>> go tool pprof -pdf -output cpu.pdf cpu.out
>>> go tool pprof -pdf -output mem.pdf mem.out
```

**İnteraktiv analiz:**
```bash
>>> go tool pprof cpu.out
(pprof) top
# flat  flat%   sum%    cum   cum%
# 250ms 20.83% 20.83%  250ms 20.83%  sync.(*Mutex).Unlock
# 220ms 18.33% 39.17%  530ms 44.17%  math/rand.(*lockedSource).Int63
# 150ms 12.50% 69.17%  690ms 57.50%  math/rand.(*Rand).Int31n
# ...
(pprof) list BuildGraph
# .    390ms   15:  from := rand.Intn(vertices)
# .    370ms   16:  to := rand.Intn(vertices)
```

**Sub-kod izahı:**
- `top` → ən çox vaxt aparan funksiyalar (azalan sıra ilə)
- `flat` → funksiyanın özündə keçən vaxt; `cum` → uşaqları ilə birgə
- `list Funksiya` → sətir-sətir vaxt bölgüsü
- Nümunə nəticə: BuildGraph-in 72.5%-i `rand.Intn` çağırışlarına gedir →
  təsadüfi say generatoru bottle-neck-dür → daha sürətli alternativ axtar

## Əsas terminlər
- Test function — `TestXXX(t *testing.T)` imzalı funksiya
- go test — test icra aləti (-v verbose, -run regex, -short)
- Subtest — `t.Run("ad", fn)` ilə qruplaşdırılmış alt-test
- TestMain — setup/teardown üçün giriş nöqtəsi (`m.Run()` + `os.Exit`)
- Example — sənədləşən, çıxışı yoxlanan nümunə
- Benchmark — `BenchmarkXXX(b *testing.B)`, `b.N` iterasiyası
- RunParallel — paralel benchmark rejimi
- Coverage — test olunan kod faizi
- Profiling — CPU/mem istifadə analizi (pprof)

## Praktik nəticə
Go-da test ayrıca framework tələb etmir: standart `testing` paketi + `go
test`. Böyük test dəstlərini `t.Run` + `key=value` subtest adları ilə
qurun və `-run /Regex` ilə süzün. Nümunələr üçün `Example...` funksiyaları
`// Output:` ilə həm sənəd, həm testdir. Performans üçün `-bench` +
`-cpu` matrisi; keyfiyyət üçün `-coverprofile` + `go tool cover`;
 darboğaz tapmaq üçün `-cpuprofile/-memprofile` + `go tool pprof` (`top`,
 `list`).

## Mənbə
Pages: 221-251 (PDF 221-251)
