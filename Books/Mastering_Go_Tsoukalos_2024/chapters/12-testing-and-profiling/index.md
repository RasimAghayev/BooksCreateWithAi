# Chapter 12 — Code Testing and Profiling (Kod Testi və Profilinq)

## Bu chapter nədən bəhs edir?

Kod optimallaşdırma fəlsəfəsi, main()-in test üçün run()-ə çevrilməsi, profiling
(runtime/pprof CLI + net/http/pprof HTTP + web interfeys), go tool trace, httptrace,
mux.Walk ilə route yoxlaması, testing paketi (test funksiyaları, UNIX siqnal testi,
keş söndürmə, TempDir/Cleanup, testing/quick, timeout, Fatal, table-driven testlər,
code coverage), go vet, httptest ilə HTTP+DB testi, govulncheck, cross-compilation,
go:generate və example funksiyaları.

## Əsas fikirlər

### 1. Optimallaşdırma Fəlsəfəsi
- **"Make it work, then make it beautiful, then if you really have to, make it fast"**
  (Joe Armstrong); "Premature optimization is the root of all evil" (Knuth)
- Əvvəl DÜZGÜNLÜK (correctness, readability, simplicity, maintainability), sonra sürət
- Performans testi üçün production-dan biraz YAVAŞ maşında işə sal

### 2. main() → run() Test Üçün
**Problem:** main() test kodundan çağırıla BİLMİR.
**Həll:**
```go
func main() {
    err := run(os.Args, os.Stdout)
    if err != nil { fmt.Printf("%s\n", err) }
}

func run(args []string, stdout io.Writer) error {   // test-dən çağırıla bilir!
    if len(args) == 1 { return errors.New("No input!") }
    return nil
}
```
Parametrlər eksplisit — testlərdə öz dəyərlərini ötürmək olar.

### 3. Profiling — runtime/pprof (CLI app)
**Kitabdan kod nümunəsi:**
```go
// CPU profili:
cpuFile, _ := os.Create(path.Join(os.TempDir(), "cpuProfileCla.out"))
pprof.StartCPUProfile(cpuFile)
defer pprof.StopCPUProfile()            // main çıxanda yazılır
// ... CPU-intensiv kod ...

// Yaddaş profili:
memory, _ := os.Create(path.Join(os.TempDir(), "memoryProfileCla.out"))
defer memory.Close()
// ... yaddaş-intensiv kod (böyük slice-lar) ...
err = pprof.WriteHeapProfile(memory)
```
**Analiz:**
```bash
go tool pprof cpuProfileCla.out
(pprof) top              # top 10 flat
(pprof) top10 -cum       # kumulativ
(pprof) list main.N1     # funksiya səviyyəli analiz — hansı SƏTİR vaxt aparır
(pprof) pdf             # PDF qraf — müəllifin ilk seçimi
go tool pprof -http=127.0.0.1:1234 cpuProfileCla.out   # WEB interfeys (Go 1.10+)
```
`list main.N1` göstərdi: N1-in bütün vaxtı `if (n % i) == 0` sətirində — sətir
səviyyəsində bottleneck aşkarlanır.

### 4. HTTP Server Profili — net/http/pprof
```go
import _ "net/http/pprof"    // BLANK import — /debug/pprof/ handler-ları qeyd olunur

// Server işləyərkən İKİNCİ termimalda:
go tool pprof http://localhost:8001/debug/pprof/profile
# 30 san data toplayır → interaktiv pprof shell (eyni komandalar)
```
Uzunmüddətli (long-running) proqramlar üçün ideal.

### 5. go tool trace — runtime/trace
**Nə göstərir:** GC əməliyyatı, goroutine həyatı, hər P-nin aktivliyi, OS thread sayı.
```go
f, _ := os.Create(path.Join(os.TempDir(), "traceCLA.out"))
trace.Start(f)
defer trace.Stop()
// ... GC-un tetiklənməsi üçün böyük allocation-lar ...
```
```bash
go tool trace traceCLA.out    # avtomatik browser açılır — View trace = goroutine+GC
```
Üç mənbə: runtime/trace, net/http/pprof, `go test -trace`.

### 6. httptrace — HTTP Sorğu Fəzalarının İzlənməsi
**Kitabdan kod nümunəsi:**
```go
trace := &httptrace.ClientTrace{
    GotFirstResponseByte: func() { fmt.Println("First response byte!") },
    GotConn: func(connInfo httptrace.GotConnInfo) { /* reuse info */ },
    DNSDone: func(dnsInfo httptrace.DNSDoneInfo) { /* DNS addrları */ },
    ConnectStart: func(network, addr string) { fmt.Println("Dial start") },
    ConnectDone: func(network, addr string, err error) { fmt.Println("Dial done") },
    WroteHeaders: func() { fmt.Println("Wrote headers") },
}
req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
```
**Nəticə:** DNS → Dial → GotConn → WroteHeaders → FirstResponseByte ardıcıllığı görünür —
troubleshooting qızıldır. Yalnız TƏK RoundTrip izləyir.

### 7. mux.Walk — Route Yoxlaması
```go
err := r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
    pathTemplate, _ := route.GetPathTemplate()   // /username/{id:[0-9]+}
    pathRegexp, _ := route.GetPathRegexp()       // ^/username/(?P<v0>[0-9]+)$
    methods, _ := route.GetMethods()             // GET/POST/PUT/DELETE
    return nil
})
```
Yoxlanılır: endpoint mövcuddurmu, düzgün metod, parametr varmı.

### 8. Test Funksiyaları — Əsas Qaydalar
- Ad `Test` ilə başlayır + Böyük hərf/_ → TestFunctionName() (idiomatik)
- `*_test.go` fayllarında; `t *testing.T` parametri; heç nə qaytarmır
- `go test -v` — verbose; test YALNIZ bugin MÖVCUDLUĞUNU göstərir, OLMAMASINI YOX

**Kitabdan kod nümunəsi (regex funksiyası üçün):**
```go
func TestMatchInt(t *testing.T) {
    if matchInt("") { t.Error(`matchInt("") != false`) }        // boş input
    if matchInt("00") == false { t.Error(...) }                  // normal
    if matchInt("-00") == false { t.Error(...) }                  // işarə
    if matchInt("+00") == false { t.Error(...) }                  // işarə
}
func TestWithRandom(t *testing.T) {          // RANDOM input — numeric test üçün üsul
    n := strconv.Itoa(random(-100000, 19999))
    if matchInt(n) == false { t.Error("n = ", n) }
}
```

### 9. UNIX Siqnal Testi
**Kitabdan kod nümunəsi:**
```go
func TestAll(t *testing.T) {
    go Listener()                    // siqnal dinləyicisi — goroutine ŞƏRT
    time.Sleep(time.Second)
    test_SIGUSR1(); time.Sleep(time.Second)
    test_SIGUSR2(); time.Sleep(time.Second)
    test_SIGHUP();    time.Sleep(time.Second)   // son emal üçün vaxt
}
func test_SIGUSR1() {
    syscall.Kill(syscall.Getpid(), syscall.SIGUSR1)   // öz prosesinə siqnal
}
```
Testlər main paketi XARİCİNDƏ olmalı (Listener→main rename etmək asan olsun).

### 10. Keş, TempDir, Cleanup
```bash
go clean -testcache              # bütün keş təmizlə
go test -count=1                 # BU icra üçün keş YOX
go test -timeout 1s -count 2     # vaxt həddi + təkrar
```
```go
t1 := t.TempDir()                 // UNİKAL müvəqqəti qovluq — avtomatik SİLİNİR
t.Cleanup(func() { /* əl ilə yaratdığın resursların təmizliyi */ })
t.Cleanup(myCleanUp())            // adlı funksiya — təkrar istifadə üçün
```
testing.TempDir ≠ os.TempDir (sonuncu sadəcə PATH qaytarır). Təmizlik edilməyən
qovluq növbəti testdə "file exists" xətası verir.

### 11. testing/quick — Property-Based Test
```go
func TestWithItself(t *testing.T) {
    condition := func(a, b Point2D) bool {
        return Add(a, b) == Add(b, a)       // kommutativlik xassəsi
    }
    err := quick.Check(condition, &quick.Config{MaxCount: N})
    if err != nil { t.Errorf("Error: %v", err) }
}
```
quick.Check funksiya imzasından RANDOM dəyərlər GENERATE edir (Haskell QuickCheck
modeli); uğursuz halda FAILED INPUT göstərilir.

### 12. Fatal vs Error
- `t.Error/t.Errorf` — davam et; asılılıq yoxdursa
- `t.Fatal/t.Fatalf` — DAYANDIR; sonrakı testlər bu uğursuzluqdan MÜTLƏQ düşəcəksə
  (DB bağlantısı, şəbəkə bağlantısı qurulmadıqda)

### 13. Table-Driven Testlər
**Kitabdan kod nümunəsi:**
```go
type myTest struct {
    a, b     int
    resInt   int
    resFloat float64
}
var tests = []myTest{
    {a: 1, b: 2, resInt: 0, resFloat: 0.5},
    {a: 2, b: 2, resInt: 1, resFloat: 1.0},
    // ...
}
func TestAll(t *testing.T) {
    t.Parallel()                    // paralel icra (digər paralel testlərlə)
    for _, test := range tests {
        if intDiv(test.a, test.b) != test.resInt {
            t.Errorf("Expected %d, got %d", test.resInt, ...)
        }
        if floatDiv(test.a, test.b) != test.resFloat {
            t.Errorf("Expected %f, got %f", ...)
        }
    }
}
```
**Üstünlük:** yeni test = struct-a YENİ ENTRİ (yeni funksiya YOX).
**Xəbərdarlıq:** float == müqayisəsi etibarsızdır — production-da epsilonla yoxla.

### 14. Code Coverage
```bash
go test -cover *.go                      # coverage: 50.0% of statements
go test -coverprofile=coverage.out *.go  # hesabat faylı
go tool cover -html=coverage.out          # brauzer: qırmızı = icra olunmadı
```
Coverage məntiqi xətaları üzə çıxarır (həmişə true olan if-else, çatıcılmayan
branch-lar) — unit testi ƏVƏZ ETMİR, tamamlayır. Coverage azsa problem testlərdə də ola
bilər.

### 15. go vet — Çatıcılmayan Kod
```bash
go vet cannotReach.go
# ./cannotReach.go:9:2: unreachable code
```
return-dən sonra gələn kod kimi məntiqi xətaları tapır. Daha güclü alternativ:
**staticcheck** (VS Code, Neovim, Zed inteqrasiyalı). go vet-i workflow-a daxil et.

### 16. httptest — HTTP Handler Testi (DB backend ilə)
**Kitabdan kod nümunəsi:**
```go
func TestLogin(t *testing.T) {
    UserPass := []byte(`{"Username": "admin", "Password": "admin"}`)
    req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(UserPass))
    req.Header.Set("Content-Type", "application/json")
    rr := httptest.NewRecorder()               // cavabı qeyd edən recorder
    handler := http.HandlerFunc(LoginHandler)
    handler.ServeHTTP(rr, req)                 // server İŞLƏDİLMƏDƏN sorğu!
    if rr.Code != http.StatusOK {
        t.Errorf("got %v want %v", rr.Code, http.StatusOK)
    }
}

// Path dəyişəni testi — mux.SetURLVars:
req, _ = http.NewRequest("GET", "/username/1", bytes.NewBuffer(UserPass))
req = mux.SetURLVars(req, map[string]string{"id": "1"})   // URL vars əl ilə
handler.ServeHTTP(rr, req)
expected := `{"id":1,"username":"admin",...}`
if strings.TrimSpace(serverResponse) != expected { t.Errorf(...) }
```
**Dərslər:** server başlamadan handler test edilir; gorilla/mux SetURLVars path
dəyişənlərini süni şəkildə verir; unikal username = timestamp prefiks; test asılılıqları
(Lugin→Logout ardıcıllığı) random order-da problem yaradır.

### 17. govulncheck — Asılılıq Zəiflikləri
```bash
go install golang.org/x/vuln/cmd/govulncheck@latest
govulncheck ./...          # skan
govulncheck -json ./...    # JSON çıxış
# Nümunə tapıntı: GO-2022-1059 DoS in golang.org/x/text/language
#   Found in: v0.3.5, Fixed in: v0.3.8
go get golang.org/x/text@latest    # həll — upgrade
```
Vulnerability DB: vuln.go.dev. Hər layihədə vərdiş kimi işlət.

### 18. Cross-Compilation
```bash
env GOOS=linux GOARCH=amd64 go build crossCompile.go
# macOS-dan Linux x86-64 üçün binary! Statik linked — asılılıqsız.
go tool dist list     # bütün keçərli GOOS/GOARCH kombinasiyaları
```
runtime.GOOS/GOARCH/Version — cari mühiti çap edir. CI/CD multi-distribution üçün ideal.

### 19. go:generate
```go
//go:generate ./echo.sh
//go:generate echo GOFILE: $GOFILE     // xüsusi dəyişənlər runtime-da expands olunur
//go:generate echo GOOS: $GOOS
//go:generate ls -l
//go:generate ./hello.py
```
```bash
go generate        # icra (avtomatik YOX — eksplisit)
go generate -n     # nə icra olunacaq (dry run)
go generate -x     # icra + çap
```
İstifadə halları: runtime QABAQ data yüklə, versiya generasiya et, sample DB hazırla.
**Müəllif mövqeyi:** developer-dən şeyləri GİZLƏDİR — mümkünsə QAÇ.

### 20. Example Funksiyaları
```go
func ExampleLengthRange() {
    fmt.Println(LengthRange("Mihalis"))
    fmt.Println(LengthRange("Mastering Go, 4th edition!"))
    // Output:
    // 7
    // 26
}
```
`// Output:` şərhindən sonra go test çıxışı MÜQAYİSƏ edir — EXECUTABLE documentation.
Adlar Example ilə başlar; parametr/qaytarma YOX; testing paketi importu lazım DEYİL;
godoc sənədlərində görünür.

## Əsas terminlər
- Profiling — icra ölçmələri ilə davranış analizi
- runtime/pprof / net/http/pprof — CLI / HTTP profiling paketləri
- go tool pprof — profil analiz aləti (top/list/pdf/web)
- Code Tracing — runtime hadisələrinin (GC, goroutine) izi
- httptrace.ClientTrace — HTTP sorğu fazası hook-ları
- Test Function — Test*[A-Z_], *_test.go, *testing.T
- Table-Driven Test — struct slice üzərində çoxsaylı ssenari testi
- t.Parallel — paralel test icrası
- testing.TempDir — avtomatik-silinən unikal müvəqqəti qovluq
- t.Cleanup — test sonu təmizlik callback-i
- testing/quick — property-based random test generator
- t.Fatal vs t.Error — dayandır / davam et
- Code Coverage — testlərin toxunduğu kod nisbəti
- go vet — şübhəli konstruktlar; unreachable code
- httptest.NewRecorder — server-siz handler cavab qeydedicisi
- mux.SetURLVars — test üçün path dəyişənlərinin inyeksiyası
- govulncheck — asılılıq zəiflik skaneri
- Cross-Compilation — GOOS/GOARCH ilə digər platformaya build
- go:generate — build əvvəsi əmr avtomatlaşdırma
- Example Function — sənəd+test ikilisi (// Output:)

## Praktik nətidə

(1) Əvvəl düzgün, sonra gözəl, yalnız lazımsa sürətli. (2) main() minimal shell + run()
test edilə bilən məntiq. (3) pprof: CPU üçün Start/StopCPUProfile, yaddaş üçün
WriteHeapProfile; analizə `pdf`-dən başla, sonra top/list. (4) HTTP server: blank import
net/http/pprof → /debug/pprof. (5) Trace GC/goroutine davranışını, pprof isə HARƏKƏTLİ
VAXTI göstərir — fərqli suallara cavab. (6) Testlərdə: gözlənilən+gözlənilməyən+boş+edge
case + random input. (7) Table-driven — yeni hal = yeni sətir. (8) Fatal yalnız kaskad
uğursuzluqda. (9) TempDir avtomatik silir; əl yaratdıqların üçün Cleanup yaz. (10)
httptest + mux.SetURLVars — serveri işə SALMADAN handler testi. (11) go vet və
govulncheck-i CI-ya daxil et. (12) GOOS/GOARCH — tək maşından bütün platformalara build.
(13) Example funksiyaları həm sənəd həm test — paketlərində istifadə et.

## Mənbə
Pages: 519-578 (PDF 550-611)
