# Chapters 8-14 — Testing, Benchmarks, Tooling, Security, CI, Deployment, Monitoring

## Bu bölmələr nədən bəhs edir?

Test suite (table-driven, HTTP handler testləri, mocking, coverage, examples), benchmarking (Fibonacci optimallaşdırma hekayəsi, benchcmp, ResetTimer, -benchmem, escape analysis, modulo vs bitwise-and), tooling (godoc, Go Guru, race detector, Go Report Card), security (CSRF, HSTS, CSP), CI (Travis, Drone, Makefile) və qeyri-tamam Monitoring/Optimization bölmələri.

**PDF səhifələr:** 88-145 (Testing: 88-106, Benchmarks: 107-118, Tooling: 119-133, Security: 134-137, CI: 138-140, Deployment: 141-142, Monitoring+: 142-145)

## Əsas fikirlər

### Ch 8 — Testing

**Niyə test? — double-entry bookkeeping analogiyası:** hər tranzaksiya 2 hesabda əks olunur → səhv görünür. Hər funksiya üçün test = sistemin checks-and-balances-i. Production dəyişəcək — manual test davam edə bilməz.

**Test əsasları:**
```go
// prime_test.go — EYNİ paketdə, _test.go suffiksi
func TestIsPrime(t *testing.T) {
    got := IsPrime(19)
    if got != true {
        t.Errorf("IsPrime(%d) = %t, want %t", 19, got, true)
    }
}
```
- `go test` / `go test -v` (həmişə -v tövsiyə olunur)
- Go assertion funksiyaları VERMİYİB — adi kontrol axını (if + Errorf) — öyrənmə əyrisi minimal

**Table-driven test (ƏSAS pattern):**
```go
func TestIsPrimeTD(t *testing.T) {
    cases := []struct {
        give int
        want bool
    }{
        {19, true}, {21, false}, {10007, true},
        {1, false}, {0, false}, {-1, false},
    }
    for _, c := range cases {
        got := IsPrime(c.give)
        if got != c.want {
            t.Errorf("IsPrime(%d) = %t, want %t", c.give, got, c.want)
        }
    }
}
```
Edge case-lər (0, 1, mənfi) bir sətirlə əlavə olunur.

**Error message konvensiyası (CodeReviewComments):**
- `actual != expected` sırası
- Mesaj: FUNKSİYA(parametr) = actual, want expected
- `IsPrime(17) = true, want false` → nə çağırıldı, nə oldu, nə gözlənilirdi — 3 sual bir sətirdə

**t.Error vs t.Fatal:**
- `t.Error/t.Errorf` — logla + FAIL + DAVAM ET
- `t.Fatal/t.Fatalf` — logla + FAIL + DAYANDIR (setup xətaları üçün)

**HTTP handler testi (httptest):**
```go
func TestHTTPHandler(t *testing.T) {
    req, err := http.NewRequest("GET", "/hello", nil)
    if err != nil {
        t.Fatal(err)
    }

    // ResponseRecorder — ResponseWriter interfeysini qeyd edir
    r := httptest.NewRecorder()
    handler := http.HandlerFunc(helloHandler)     // adi funksiya → Handler adapter

    handler.ServeHTTP(r, req)                    // server OLMADAN icra!

    if r.Code != http.StatusOK { ... }           // status yoxla
    want := "Hello, friend :)\n"
    got := r.Body.String()                       // body yoxla
    if got != want { ... }
}
```

**Mocking — öz interfeysin:**
```go
// randIntGenerator — 1 metodlu interfeys, math/rand.Rand-ın Intn imzası
type randIntGenerator interface {
    Intn(int) int
}

type EightBall struct {
    rand randIntGenerator          // unexported — istifadəçi GÖRMÜR
}

func New() *EightBall {
    return &EightBall{rand: rand.New(rand.NewSource(time.Now().UnixNano()))}
}

func (e EightBall) Answer(s string) string {
    n := e.rand.Intn(3)           // interfeysdən istifadə → MOCK MÜMKÜN
    switch n { case 0: return "Definitely not"; ... }
}

// Test — FIXED generator:
type fixedRandIntGenerator struct {
    randomNum   int               // "random" nəticə
    calledWithN int               // Intn nə ilə çağırıldı (qeyd et!)
}

func (g *fixedRandIntGenerator) Intn(n int) int {
    g.calledWithN = n             // arqumenti SAXLA — çağırış yoxlaması
    return g.randomNum            // İSTƏDİYİMİZ dəyəri qaytar
}

// Test case: {randomNum → want} xəritəsi + calledWithN == 3 assert
```
- Java-dan fərq: vendor-dan asılılıq YOX — ÖZ interfeysini yazırsan, kitabxana AVTOMATİK implement edir (implicit)
- Random seed fiksle ≠ hər nəticəni SPESİFİK test etmək

**Coverage:**
```bash
$ go test -cover                          # coverage: 62.5% of statements
$ go test -coverprofile=coverage.out      # profil faylı
$ go tool cover -func=coverage.out        # funksiya-xətti faiz
$ go tool cover -html=coverage.out        # BRAUZERDƏ sətir-sətir yaşıllı/qırmızı!
```
Username validasiya nümunəsi: 1 test → 62.5%; 4 case (boş, normal, $, 31 hərf) → **100%**.

**Examples:**
```go
func ExampleUsername() {
    ...
    fmt.Printf("%q: %t\n", tt.in, valid)
    // Output:          ← go test BUNU YOXLAYIR!
    // "": false
    // "gopher": true
    // "gopher$": false
}
```
- Adlandırma: `ExampleFunksiya`, `ExampleTip_Metod` (ExampleUser_ValidateName), `Example_funksiya_second`
- godoc-da render olunur — sənəd + test bir yerdə

### Ch 9 — Benchmarks

**Fibonacci optimallaşdırma hekayəsi:**
```go
// Recursive (dərslik versiyası):
func F(n int) int {
    if n <= 0 { return 0 } else if n == 1 { return 1 }
    return F(n-1) + F(n-2)
}

// İterativ (3 dəyişən, 0 allokasiya):
func FastF(n int) int {
    var a, b int = 0, 1
    for i := 0; i < n; i++ {
        a, b = b, a+b          // paralel təyinat!
    }
    return a
}
```

**Benchmark:**
```go
var numbers = []int{0, 10, 20, 30}

func BenchmarkF(b *testing.B) {
    m := len(numbers)
    for n := 0; n < b.N; n++ {
        F(numbers[n%m])
    }
}
```
```bash
$ go test -bench=.
BenchmarkF-4      1000     1255534 ns/op     # 1.2 ms
BenchmarkFastF-4  50000000     20.3 ns/op     # ~60,000x FƏRQ!
```
- b.N avtomatik kalibrlənir; -bench regex
- -4 → 4 CPU; goos/goarch başlıq

**benchcmp (köhnə/yeni müqayisə):**
```bash
$ go get golang.org/x/tools/cmd/benchcmp
$ go test -bench . > old.txt     # dəyişiklikdən ƏVVƏL
$ go test -bench . > new.txt     # dəyişiklikdən SONRA
$ benchcmp old.txt new.txt
# BenchmarkF-4  1965113  25.0  -100.00%
```

**b.ResetTimer (setup-i çıxart):**
```go
func BenchmarkEncrypt(b *testing.B) {
    tt := encryptTests[0]
    c, err := NewCipher(tt.key)     // SETUP
    if err != nil { b.Fatal("NewCipher:", err) }
    out := make([]byte, len(tt.in))
    b.SetBytes(int64(len(out)))      # bayt sayı → B/s hesabı
    b.ResetTimer()                   # taymeri İNDİ sıfırla — setup sayılmasın
    for i := 0; i < b.N; i++ {
        c.Encrypt(out, tt.in)
    }
}
```

**-benchmem (allokasiya statistikası):**
```bash
$ go test -bench=. -benchmem
BenchmarkF-4           1000   1241017 ns/op    0 B/op   0 allocs/op
BenchmarkFastHighMemF-4 20000000  72.0 ns/op  132 B/op   0 allocs/op
```
- `B/op` — əməliyyat başına bayt; `allocs/op` — heap allokasiyası
- HighMem variantı (`make([]int, n+1)`) 132 B/op — slice başlığı + orta n=15 × 8 bayt

**Escape analysis (-gcflags=-m):**
```bash
$ go test -bench=. -benchmem -gcflags=-m
./fibonacci_bench_test.go:10:20: BenchmarkF b does not escape
./fibonacci_test.go:22:38: tt.n escapes to heap
```
- Kaçmayan dəyişənlər → STACK (heap + GC xərcindən azad)
- fibonacci.go-un YOXLUĞU outputda = heç bir dəyişən escape etmədi

**⚠️ Modulo vs Bitwise-and (benchmark dəqiqliyi):**
```go
FastF(nums[n%m])     // IDIVQ asm əmri — YAVAŞ
m := len(nums)-1
FastF(nums[n&m])     // ANDQ asm — SÜRƏTLİ (m+1 = 2-nin qüvvəti!)
```
**Qanun:** `n % m == n & (m-1)` — m 2-nin qüvvətidirsə!
```
BenchmarkFastFModulo-4      100000000   16.7 ns/op
BenchmarkFastFBitwiseAnd-4  200000000    7.40 ns/op   # 2x sürətli!
```
→ Benchmark-un ÖZ overhead-i nəticəni pozur — kiçik kodlarda & istifadə et. `go tool compile -S` ASM çıxarır.

### Ch 10 — Tooling

**Godoc:**
```bash
$ go get golang.org/x/tools/cmd/godoc
$ godoc encoding/json           # CLI
$ godoc encoding/json Marshal   # konkret funksiya
$ godoc -http:6060              # lokal HTML server
```
- Şərh konvensiyası: elanın DƏRHAL üstündə, boş sətir YOX
- Paket şərhi: "Package [name] ..." — yoxsa ayrı doc.go faylı
- godoc.org — GitHub-da host olunan paketlər üçün

**Go Guru** (`go get -u golang.org/x/tools/cmd/guru`):
- Vim inteqrasiyası: `:GoImplements` (bu interfeysi kim implement edir), `:GoReferrers`, `:GoCallees`, `:GoCallers`

**Race detector — 2 tam nümunə:**

**Nümunə 1 — Cat:**
```go
func updateCat(c *Cat) {
    go c.SetNoise("にゃん")       // FONDA yaz
    log.Println(c.Noise())       // EYNİ ANDA oxu → RACE!
}
```
```bash
$ go test -race
WARNING: DATA RACE
Write at 0x... by goroutine 7:  cat.(*Cat).SetNoise()  cat.go:14
Previous read ... by goroutine 6: cat.updateCat()       cat.go:19
```
→ DƏQİQ sətirlər! Həll: **struct-a Mutex embed + hər oxu/yazıda Lock:**
```go
type Cat struct {
    mu    sync.Mutex
    noise string
}
func (c *Cat) SetNoise(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.noise = n }
func (c *Cat) Noise() string      { c.mu.Lock(); defer c.mu.Unlock(); return c.noise }
```

**Nümunə 2 — API (go run -race ilə production testi):**
```go
go func(e *Entry) { normalizeCountry(e) }(&entry)   // fondda normalize
log.Printf("INFO: received Entry %v", entry)         // oxu → RACE
```
Curl-də göstərilir: `{enGLaND ...}` — normalize İŞLƏMƏDİNİ göstərir (race-in reallığı!). Həll: **WaitGroup ilə gözlə:**
```go
var wg sync.WaitGroup
wg.Add(1)
go func(e *Entry) {
    defer wg.Done()
    normalizeCountry(e)
}(&entry)
wg.Wait()                     // bitməsini gözlə → race YOX + England QAYTARIR
```

**Go Report Card:** goreportcard.com — paketə GPA verir (gofmt, golint, vet, misspell, ...); müəlliflərin ÖZ alətləri! Minio müəlliflərindəndir.

### Ch 11 — Security

**Keep Go up to date:** golang-announce siyahısı — `[security]` prefiksli elanlar (məs. 1.11.3: `go get -u` remote execution bug).

**CSRF (Cross-Site Request Forgery):** İstifadəçi bank saytında login — zərərli səhifə formu AVTOMATİK submit edir:
```html
<body onload="document.forms[0].submit()">
<form action="http://insecurebank.com/transfer" method="POST">
    <input type="hidden" name="account" value="5555555555" />
```
**Həll:** per-session CSRF token (hidden field); POST-da yoxsа/yoxsa uyğunsusa → DENY. Paketlər: **nosurf** (justinas), **gorilla/csrf**.

**HSTS (Strict-Transport-Security):** protocol downgrade + cookie hijacking qarşısı:
```go
func headerWrap(h http.HandlerFunc) http.HandlerFunc {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Add("Strict-Transport-Security",
            "max-age=31536000; includeSubDomains")     // 1 il
        h.ServeHTTP(w, r)
    })
}

http.HandleFunc("/", headerWrap(homeHandler))    // BÜTÜN handler-ları wrap et
```
Curl yoxlama: `curl -i localhost:8000` → `Strict-Transport-Security: max-age=...`

**CSP (Content-Security-Policy):** XSS qarşısı — məzmun mənbələrini məhdudlaşdırır. (Kitabda yarımçıq.)

**See Also:** bluemonday (sanitizer), unrolled/secure (header middleware), Mozilla security guidelines.

### Ch 12 — Continuous Integration

| CI | Open Source | Config | Pulsuz |
|---|---|---|---|
| Travis CI | YOX | .travis.yml | OSS üçün |
| **Drone** | BƏLİ | .drone.yml (Docker!) | self-host + enterprise |
| Jenkins | BƏLİ | Jenkinsfile | self-host |
| TeamCity | YOX | UI | freemium |

Tövsiyə: **Travis və ya Drone** (sadəlik).

**.travis.yml (Go Report Card-dan):**
```yaml
language: go

go:
  - 1.8.x
  - 1.9.x
  - 1.10.x
  - 1.11.x
  - tip

install:
  - make install

script:
  - make lint
  - make test
```
→ 5 Go versiyasında lint + test.

**Makefile (Go Report Card-dan):**
```makefile
all: lint build test

build:
	go build ./...

install:
	./scripts/make-install.sh

lint:
	gometalinter --exclude=vendor --exclude=repos --disable-all --enable=golint --enable=vet --enable=gofmt ./...
	find . -name '*.go' | xargs gofmt -w -s

test:
	go test -cover ./check ./handlers

start:
	go run main.go

misspell:
	find . -name '*.go' -not -path './vendor/*' -not -path './_repos/*' | xargs misspell -error
```

### Ch 13-14 — Deployment, Monitoring (yarımçıq)

- **Deployment:** mövcud infrastrukturdan asılı — kitabda tamamlanmamış
- **Prometheus:** Go-da yazılmış monitoring sistemi; **Grafana** qrafiklər; Alerts (gözlənilən)
- **Optimization:** TODO — sadəcə başlıq
- **Common Gotchas — Nil interface** (Ch 3-də əhatə olundu)
- **Further Reading:** Spec, Effective Go, Golang Weekly, Go Blog, Gophers Slack (#golang-newbies, #performance)

## Əsas terminlər
- Table-driven test / anonymous struct cases
- actual != expected konvensiyası
- t.Error vs t.Fatal
- httptest.NewRecorder / http.HandlerFunc adapter / ServeHTTP
- randIntGenerator (mock interfeys) / fixedRandIntGenerator (calledWithN qeydi)
- go test -cover / -coverprofile / go tool cover -func / -html
- ExampleX / ExampleT_M / Example_suffix / // Output:
- b.N / -bench regex / -benchmem / B/op / allocs/op
- benchcmp (old vs new delta)
- b.ResetTimer / b.SetBytes
- Escape Analysis (-gcflags=-m) / stack vs heap
- n%m == n&(m-1) (2-nin qüvvəti) / go tool compile -S
- godoc / doc.go / "Package [name] ..."
- Go Guru (:GoImplements, :GoReferrers)
- go test/run/build -race / WARNING: DATA RACE
- WaitGroup ilə race həlli (handler-da)
- Go Report Card / goreportcard.com
- CSRF / per-session token / nosurf / gorilla/csrf
- HSTS / Strict-Transport-Security / max-age / includeSubDomains / headerWrap
- CSP / bluemonday / unrolled/secure
- .travis.yml / .drone.yml / Makefile target-ləri
- Prometheus / Grafana

## Praktik nəticə
- Test mesajları üçünlü olsun: FUNKSİYA(input) = actual, want expected — gələcəkdə debug edən sən olmaya bilərsən.
- Table-driven test — default seçim; anonymous struct ilə case-lər 1 sətirdə əlavə olunur.
- Mock üçün ÖZ interfeysini yaz (1 metod kifayət) — kitabxana implicit implement edir; mock-da çağırış arqumentlərini də QEYD ET.
- HTTP handler-ları serverə ehtiyac olmadan Recorder ilə test et.
- Coverage → HTML report → qırmızı sətirləri hədəflə → 100% üçün edge case-lər.
- Benchmark-da setup-i ResetTimer-dan əvvələ qoy; kiçik kod benchmarklarında %-i YOX, &m-1 istifadə et (2x dəqiqlik).
- -race-i hər testdə, hətta go run -race ilə İNTEQRASİYA testlərində də işlət — race sadece go test-də deyil, production path-də də görünür.
- HSTS/CSP header-larını wrapper handler ilə MƏRKƏZİ şəkildə qoy.
- CI: Makefile (lint/build/test target-ləri) + .travis.yml (multi-version) — Go Report Card modeli.
- Kitab Monitoring/Optimization bölmələri yarımçıqdır — Prometheus/Grafana yalnız xatırladılır.

## Mənbə
Pages: 88-145 (PDF), book pages 82-139
