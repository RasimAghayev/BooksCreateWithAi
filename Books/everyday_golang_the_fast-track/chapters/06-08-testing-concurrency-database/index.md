# Chapters 6-8 — Testlər, Konkurrentlik, Verilənlər bazası (səh. 34-79)

## Bu fəsillər nədən bəhs edir?

**Ch6 (Unit-testlər):** Go-nun assertion-siz test fəlsəfəsi (want/got + t.Fatalf),
test-table pattern, coverage (`-cover`, HTML report), benchmark (`-bench`,
b.N), testlərin binariyada OLMAMASI (GoFiles vs TestGoFiles), `-test.count=100`
stress, `t.Parallel()` + `-test.parallel N` (5s → 1.3s), `go test -c` test
binariləri, asılılıqların izolyasiyası (implied interfaces, FakeReader,
testing/iotest), GetWebRequest interfeysi ilə HTTP-in fake-lənməsi,
httptest (NewRecorder/NewRequest) ilə handler testləri (basic-auth 200/401).

**Ch7 (Konkurrentlik):** goroutine + WaitGroup (Add/Done/Wait), loop-variable
capture tələsi (i → 6!, 2 həll: parametr və ya `j := i`), errgroup (xəta
toplama), singleflight (eyni açarlı paralel çağırışların birləşdirilməsi —
"thundering herd" qarşısı; OpenFaaS gateway-də real tətbiq: 330 sorğu → 33
upstream), sync.Mutex vs RWMutex (oxu-yazı ayrılığı, in-memory cache nümunəsi),
channels + select (titleCh/errorCh), context (Background/TODO, WithCancel/
WithDeadline/WithTimeout, NewRequestWithContext), worker pool pattern-i
(workQueue channel + N worker + `-w` flag; 6.2s → 2.3s).

**Ch8 (DB):** PostgreSQL + lib/pq ("pure Go" — CGO yox), database/sql
interfeysi, DigitalOcean managed DB, DSN flag-lərlə qurulur, connect+Ping,
todo cədvəli (GENERATED ALWAYS AS IDENTITY), create/list əməliyyatları.

## Əsas fikirlər

### 1. Go Test Fəlsəfəsi (Ch6)
```go
func TestMultiply_ForANegativeNumber(t *testing.T) {
    want := -10
    got := c.Multiply(10, -1)
    if want != got {
        t.Fatalf("want %d, but got %d", want, got)
    }
}
```
- **Assertion YOX:** BDD/assert kitabxanaları (expect, RSpec) — xətalar kriptik
  ("assert: 0 == 1"), stack trace daşqını, ilk fail-də dayanır (nümunə
  maskalanır). Go: adi kod, aydın mesaj, hamısı İŞLƏYİR
- **Səbəb:** Go-nun ən güclü xüsusiyyəti KONSİSTENSLİK — hər Go developer
  başqa codebase-i kitsiz oxuya bilir
- **Test adı:** Test + söz/qarışıq (Test_sum, TestMultiply_ForANegative);
  fayl: *_test.go; testlər EYNİ paketdə (private metodlar da test olunur!)
- **İlk testİ SİNDIR:** köhnə testi qırmızı gör → yalnız sonra yaşıl

### 2. Test-Table Pattern
```go
cases := []struct{ a, b, want int; name string }{
    {-10, 1, -10, "negative case"},
    {10, 1, 10, "positive case"},
}
for _, tc := range cases {
    t.Run(tc.name, func(t *testing.T) {
        got := c.Multiply(tc.a, tc.b)
        if got != tc.want { t.Errorf(...) }
    })
}
```
- NÜMUNƏ: t.Run subtest adı çapda görünür; **2-3 dublikatdan sonra** cədvələ
  keç; PARALEL versiyada `tc := c` kopyası lazımdır (closure capture!)

### 3. go test Bayraqları
| Bayraq | Nə edir |
|---|---|
| -v | hər testin PASS/FAIL-i + t.Log |
| -cover, -coverprofile + go tool cover -html | statement coverage (50%...) — amma İFA-DAN silinmiş if 50% saxlayır → rəqəm DAVRANIŞI əvəz etmir |
| -bench=., b.N | benchmark: ns/op |
| -test.count=100 | stress: race/shared-state ifşa; cache-i qır (count=1) |
| -test.parallel N + t.Parallel() | paralel: 5 test 5.3s → 1.3s |
| -c | testləri BİNARİYAYA compile → paylanan validator aləti |

- **Go testləri ship etmir:** `go list -f={{.GoFiles}} / {{.TestGoFiles}}` —
  binaridə test kodu OLMUR

### 4. Asılılıqların İzolyasiyası
- **Go interfeysləri IMPLIED:** konkret tip interfeysi BİLMİR — sonra fake
  qoşmaq olar (C#/Java-dan fərqli: implement demək lazım deyil!)
- **Kiçik interfeyslər:** io.ReadCloser = Read + Close (2 metod)
```go
type GetWebRequest interface {
    FetchBytes(url string) ([]byte, error)
}
// real: LiveGetWebRequest (http.Client...)
// test: testWebRequest{} → return []byte(`{"number": 3}`), nil
```
- **Əvvəl stdlib-a bax:** bytes.NewReader; testing/iotest — yavaş/xətalı
  Reader-lər (resilience testləri)
- **İstehlakçının asılılıq izolyasiyası:** kitabxana interfeys TƏKMİL ET —
  yoxsa istifadəçi üçün böyük iş

### 5. HTTP Handler Testləri — httptest
```go
w := httptest.NewRecorder()                        // cavabı QEYD EDİR
r := httptest.NewRequest(http.MethodGet, "http://localhost:8080", nil)
r.SetBasicAuth("admin", "password")
decorated := DecorateWithBasicAuth(handler, creds)
decorated.ServeHTTP(w, r)                          // real server YOX!
if w.Code != http.StatusOK { t.Fatalf(...) }
gotAuth := w.Header().Get("WWW-Authenticate")
```
- İki dəst ssenari: düzgün parol → 200; səhv → StatusUnauthorized +
  `WWW-Authenticate: Basic realm="Restricted"` başlığı
- **Dərs:** middleware-in (basic-auth) və handler-in testi serverə ehtiyac
  duymadan — recorder + fake request

### 6. Konkurrentlik Əsasları (Ch7)
```go
wg := sync.WaitGroup{}
wg.Add(3)
go func() { defer wg.Done(); printLater(...) }()
wg.Wait()                                  // hamı Done deyənə qədər blokla
```
- Goroutine = runtime idarəli yüngül axın (OS thread-ilə 1:1 DEYİL — ucuz,
  miqyaslı)
- **Concurrency vs parallelism:** eyni vaxtda İDARƏ vs eyni vaxtda İCRA
- Panic İSTİRAKÇI goroutine-dən — WHOLE app çökür

### 7. Loop-Variable Tələsi (ƏSAS PITFALL!)
```go
for i := 1; i <= 5; i++ {
    go func() { fmt.Printf("Hello from %d", i) }()   // hamısı "from 6"!!!
}
// Həll 1 — parametr kimi ötür:
go func(j int) { ... }(i)
// Həll 2 — lokal kopya:
j := i
go func() { ... }()
```
- **Səbəb:** closure döngə dəyişəninin SON dəyərini görür (Go 1.22-dən əvvəl)

### 8. errgroup — Xətalı Paralellik
```go
g := errgroup.Group{}
g.Go(func() error { return fmt.Errorf("world failed") })
err := g.Wait()      // ilk xəta bubble-up
```
- WaitGroup-dan fərqi: funksiyalar ERROR qaytarır → kanal sync lazım deyil;
  serial pipeline də mümkün

### 9. singleflight — Thundering Herd Qarşısı
```go
s := singleflight.Group{}
res, err, shared := s.Do(key, func() (any, error) {
    return callUpstreamAPI()     // 2 saniyə çəkən zəng
})
```
- **Problem:** 1000 paralel istifadəçi → 1000 API zəngi (expensive/rate-limited)
- **Həll:** eyni `key` ilə gələnlər BİR zəngin nəticəsini PAYLAŞIR (shared
  bool); açarlıqsə yeni zəng başlayır, digərləri ona qoşulur
- **Real:** OpenFaaS gateway scale-from-zero: 330 sorğu → 33 upstream zəng
- res `any` qaytarır → type assertion ilə aç

### 10. Mutex vs RWMutex — Data Paylaşımı
```go
type ScrapeRun struct {
    Results map[string]ScrapeResult
    Lock    *sync.Mutex          // VƏ YA *sync.RWMutex
}
sr.Lock.Lock(); defer sr.Lock.Unlock(); sr.Results[url] = ...
```
- **Mutex:** eksklüziv — oxu da yazı da bloklanır → reader+writer bottleneck
- **RWMutex:** N oxu (RLock) PARALEL + 1 yazı (Lock) eksklüziv —
  in-memory cache pattern: timer yazır, handler-lər oxuyur

### 11. Channels + select
```go
titleCh := make(chan string)      // unbuffered: göndəriş QARŞIQIDA alıcı olmalıdır
errorCh := make(chan error)
go titleOf(uri, titleCh, errorCh)
select {
case title := <-titleCh: ...
case err := <-errorCh: os.Exit(1)
}
```
- Buffered: `make(chan T, n)` — alıcısız n qədər qəbul edir
- **Xəbərdarlıq:** natamam kanal bilikləri → goroutine leak, deadlock,
  memory leak — ona görə kitab əvvəl WaitGroup/errgroup öyədir

### 12. Context — Timeout/Ləğv
```go
ctx, cancel := context.WithTimeout(context.Background(), timeout)
defer cancel()
req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
```
- WithCancel / WithDeadline(tarix) / WithTimeout(duration); TODO = hansı
  context olacağı qeyri-müəyyəndə
- Nəticə: `context deadline exceeded` — asılan saytlara qarşı müdafiə

### 13. Worker Pool (Limitli Konkurrentlik)
```go
workQueue := make(chan string)
for i := 0; i < workers; i++ {          // N worker
    go func() {
        for uri := range workQueue {    // iş götür
            title, err := titleOf(uri, 5*time.Second)
        }
        wg.Done()
    }()
}
go func() {
    for _, u := range uris { workQueue <- u }
    close(workQueue)                    // DÖNGÜNÜ BİTİRƏN SİQNAL!
}()
wg.Wait()
```
- `-w 1` → 6.2s; `-w 4` → 2.3s; **məqsəd: yalnız sürət DEYİL — resurs
  LIMITİ** (DB bağlantısı, fayl deskriptoru, API rate limit)
- faas-cli build.go-da real istifadə

### 14. PostgreSQL + database/sql (Ch8)
```go
import _ "github.com/lib/pq"      // pure Go — CGO YOX!
func connect(...) (*sql.DB, error) {
    connStr := "postgres://" + user + ":" + password + "@" + host + ...
    db, err := sql.Open("postgres", connStr)
    err = db.Ping()               // bağlantını doğrula
}
```
- **Schema:** `id INT GENERATED ALWAYS AS IDENTITY`, description NOT NULL,
  completed_date NULL-able (*time.Time pointer!)
- **Təhlükəsizlik:** parol hardcode YOX — flag/env/secrets (OpenFaaS:
  /var/openfaas/secrets/NAME)
- **Kitabxana seçimi:** "pure Go" (lib/pq) VS CGO-lu — portabiliti

## Əsas terminlər
- Assertion-free testing — want/got müqayisəsi ilə adi Go kodu
- Test-table (cədvəl testi) — ssenari slice + t.Run subtest
- Coverage (statement) — icra olunan if %; DAVRANIŞ yox
- t.Parallel() — testi paralel qrupa qoşur
- Implied interface — konkret tip interfeysi bilmir (fake-ə açıq)
- Test-double/Fake — asılılığın saxta implementasiyası
- httptest Recorder/NewRequest — server-siz HTTP testi
- WaitGroup — Add/Done/Wait üçlüyü
- Loop-variable capture — closure-da dəyişən tələsi (i → son dəyər)
- errgroup — xəta toplayan paralel qrup
- singleflight — eyni açarlı zənglərin birləşməsi (thundering herd)
- Mutex/RWMutex — eksklüziv / oxu-yazı ayrıqılı kilidi
- Buffered/unbuffered channel — alıcısız göndərilə bilən / olmayan
- Context — deadline + ləğv siqnalı daşıyıcısı
- Worker pool — N worker + iş kanalı + close() sərhədi
- Pure Go — CGO-suz kitabxana (portabiliti)

## Praktik nəticə

1. **Test yazım ardıcıllığı:** TestFoo(t *testing.T) → want/got → t.Fatalf;
   2-3 dublikatda cədvələ; İLK SİNDIRMƏDƏN başla (yaşıl test = test etmir).
2. **HTTP testləri httptest ilə:** serverə ehtiyac YOX; handler/middleware
   Recorder ilə; asılılıq interfeys + fake (implied olması qolaydır).
3. **Loop-da goroutine:** parametr VƏ YA lokal kopya — əks halda son dəyər!
4. **Paralellik pilləliyi:** WaitGroup (sadə) → errgroup (xəta) → kanallar
   (data axını) → worker pool (limit) → singleflight (dedupe).
5. **RWMutex cache pattern:** yavaş dəyişən, tez oxunan data (aktiv
   müştərilər siyahısı) — reader-lər bloklanmır.
6. **Context hər xarici zəngdə:** WithTimeout + NewRequestWithContext.
7. **Kanala close() Vacibdir:** worker-lərin range döngüsü bitmir onsuz.
8. **DB bağlantısı:** flags/env-dən DSN; Ping ilə doğrula; pure-Go driver.

## Mənbə
Pages: 34-79 (PDF 35-80)
