# Chapter 7 — Testing Applications

## Bu chapter nədən bəhs edir?

Bu chapter Go və Echo veb-tətbiqetmələrinin **test edilməsini** əhatə edir: testin 5 əsas kateqoriyası (Unit, Benchmark, Behavior, Integration, Security), standart kitabxana ilə handler/middleware unit testləri (`httptest.ResponseRecorder`), `sql.DB` mock-lanması (interface + funksiya sahəli `MockDB`), benchmark testləri, və `-coverpkg` + instrumented server ilə **real integration coverage** ölçülməsi izah olunur.

> Testin məqsədi kodun dizayn edildiyi kimi işlədiyini **empirik sübut etməkdir**. Testlər olmadan kodunınızın nə etdiyindən heç bir xəbəriniz olmur.

---

## Əsas fikirlər

### 1. Testin 5 əsas kateqoriyası

| Tip | Sual cavablandırır | Xarakteristika |
|-----|-------------------|----------------|
| **Unit testing** | "Bu kod vahidi düşündüyüm işi edir?" | Kodun yanında yaşayır, ən aşağı funksionallıq vahidinə qarşı; ilk səviyyə test |
| **Benchmark testing** | "Bu kod/sistem düşündüyüm qədər sürətli işləyir?" | Sürat və throughput ölçmə; unit səviyyəsində və inteqrasiya səviyyəsində; distributed sistemlər üçün kritik |
| **Behavior testing** | "Sistem müştərinin istifadə edəcəyi şəkildə işləyir?" | Product owner + müştəri + developer körpüsü — hər kəs anlayır, rəy asılılığı yoxdur |
| **Integration testing** | "Tətbiqetmə ətraf sistemlərlə düzgün inteqrasiya olunur?" | Mock edilməmiş real data ilə (məs., real DB) — "full-flow" test də deyilir |
| **Security testing** | "Tətbiqetmə düşündüyüm qədər təhlükəsizdir?" | Penetration testing: XSS, SQL/command injection, CORS və s. |

---

### 2. Go unit test konvensiyaları

**Adlandırma qaydaları** (standart kitabxana `go test` üçün məcburi):
1. Fayl adı **`_test.go`** ilə bitməlidir.
2. Test funksiyası **`Test`** prefiksi ilə başlamalıdır: `func TestXxx(t *testing.T)`.
3. Test **eyni package**-də (eyni qovluqda) olmalıdır — xüsusi `test/` direktoriyası yalnız exported funksionallıq üçün mümkündür, praktik deyil (non-exported test edilə bilməz).
4. `testing.T` strukturu `Fail`, `Log`, `Error` kimi helper-lər təqdim edir.

Go-da unit test üçün **üçüncü tərəf kitabxana lazım deyil** — hamısı standart kitabxanadadır.

---

### 3. Handler-in unit testi — HealthCheck nümunəsi

`handlers/health_check_test.go`:

```go
func TestHealthCheck(t *testing.T) {
    // 1. Echo instance qur — main-dəki kimi
    e := echo.New()
    e.Pre(middlewares.RequestIDMiddleware)
    e.GET("/health-check", HealthCheck)

    // 2. Saxta request + cavab yazıyıcı yarat
    w := httptest.NewRecorder()
    r, _ := http.NewRequest("GET", "/health-check", nil)

    // 3. Request-i birbaşa Echo-ya ver (Şəbəkə serverinə ehtiyac YOX)
    e.ServeHTTP(w, r)
    resp := w.Result()
    if resp.StatusCode != http.StatusOK {
        t.Error("unexpected status code: ", resp.Status)
    }

    // 4. JSON cavabı dekodla və yoxla
    healthCheckResponse := new(renderings.HealthCheckResponse)
    dec := json.NewDecoder(resp.Body)
    err := dec.Decode(healthCheckResponse)
    if err != nil {
        t.Error("error decoding", err)
    }
    if healthCheckResponse.Message != "Everything is good!" {
        t.Error("invalid response message: ", healthCheckResponse.Message)
    }
}
```

**`httptest.ResponseRecorder` nədir?** `http.ResponseWriter` interfeysinin test implementasiyasıdır — status kodu, header-lər və body-ni yadda saxlayır; `Result()` metodu ilə `*http.Response` qaytarır. Real HTTP server və şəbəkə lazım olmadan handler-i tam işlədə bilirsiniz (`e.ServeHTTP(w, r)`).

**Test axını:** Echo qur → middleware + route əlavə et → Request + Recorder yarat → `ServeHTTP` → status kodu, JSON dekodluğu, data dəyəri — 3 assertion.

İşə salma və nəticə:

```bash
go test -v ./handlers -cover
=== RUN   TestHealthCheck
--- PASS: TestHealthCheck (0.00s)
PASS
coverage: 7.9% of statements
```

---

### 4. Mocking — sql.DB nümunəsi

**Problem:** Unit test-də real SQL server (schema + fixture data) qoşulmaq integration testdir. `sql.DB` konkret tip olduğundan onu əvəz etmək olmur.

**Həll:** Öz **`MockableDB` interfeysinizi** yaradın — `sql.DB`-nin bütün metod imzaları ilə (`models/mocks.go`):

```go
type MockableDB interface {
    Begin() (*sql.Tx, error)
    BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
    Close() error
    Conn(ctx context.Context) (*sql.Conn, error)
    Driver() driver.Driver
    Exec(query string, args ...interface{}) (sql.Result, error)
    ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
    Ping() error
    PingContext(ctx context.Context) error
    Prepare(query string) (*sql.Stmt, error)
    PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
    Query(query string, args ...interface{}) (*sql.Rows, error)
    QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
    QueryRow(query string, args ...interface{}) *sql.Row
    QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
    SetConnMaxLifetime(d time.Duration)
    SetMaxIdleConns(n int)
    SetMaxOpenConns(n int)
    Stats() sql.DBStats
}
```

**`MockDB` strukturu** — hər metod üçün funksiya sahəsi daşıyır (override edilə bilən davranış):

```go
type MockDB struct {
    mockBegin    func() (*sql.Tx, error)
    mockExec     func(query string, args ...interface{}) (sql.Result, error)
    mockQuery    func(query string, args ...interface{}) (*sql.Rows, error)
    mockQueryRow func(query string, args ...interface{}) *sql.Row
    // ... qalan metodlar üçün də eyni pattern
}

func (db *MockDB) Begin() (*sql.Tx, error) {
    if db.mockBegin != nil {
        return db.mockBegin()
    }
    return nil, nil // default: passiv davranış
}
```

Test-də istifadə — xəta halını süni yarat (`models/user_test.go`):

```go
func TestGetUserByUsername(t *testing.T) {
    db := &MockDB{
        mockQuery: func(query string, args ...interface{}) (*sql.Rows, error) {
            return nil, errors.New("test query failure!")
        },
    }
    _, err := GetUserByUsername(db, "test")
    if err != nil {
        if errors.Cause(err).Error() != "test query failure!" {
            t.Errorf("incorrect failure expected: %s", err.Error())
        }
    }
}
```

**Vacib:** `GetUserByUsername(db, ...)` funksiyası `MockableDB` interfeysi qəbul etməlidir — beləcə real `*sql.DB` və `*MockDB` hər ikisi keçərilidir.

**Mocking-in riskləri:**
- "Tragically brittle" (qrırcılığa qədər kövrək) — mock sizin kodun gözlədiyini kodun özünə geri verir: 100% unit coverage integration-in işləyəcəyinə zəmanət deyil.
- **Over-mock etməyin** — mock coverage artırır, amma real inteqrasiyanı yoxlamır.

### Monkey patching

Run-time-da funksiya təyinatını əvəz etmə texnikası (https://bou.ke/blog/monkey-patching-in-go/). Məs., `sql.Open` xətasını simulyasiya etmək üçün test-də onu öz funksiyasınız ilə əvəz edirsiniz. **Çatışmazlıq:** paralel test icrasını itirirsiniz (funksiya qlobal şəkildə əvəz olunur).

---

### 5. Benchmark testi

`Benchmark` prefiksi + `*testing.B` parametri; ölçülən hissə `b.N` dövründə:

```go
func BenchmarkHealthCheck(b *testing.B) {
    e := echo.New()
    e.Pre(middlewares.RequestIDMiddleware)
    e.GET("/health-check", HealthCheck)
    w := httptest.NewRecorder()
    r, _ := http.NewRequest("GET", "/health-check", nil)
    for i := 0; i < b.N; i++ {
        e.ServeHTTP(w, r)
    }
}
```

```bash
go test -v ./... -bench Benchmark
BenchmarkHealthCheck-8  300000  5095 ns/op
```

Nəticə: handler 5,095 ns/op — 1.6 saniyədə 300K çağırış. Kodu dəyişəndə performance təsirini dəqiq bilirsiniz.

---

### 6. External/Integration testing — instrumented server

**`-coverpkg` flaqı:** test main package-də işə düşsə də coverage informasiyasını başqa package-dən (məs., `./handlers`) toplayır:

```bash
go test -coverprofile=cov.txt -coverpkg ./handlers -run TestRunMain ./cmd/service/
```

**Test rejimli server** (`cmd/service/main_test.go`):

```go
package main

import "testing"

func TestRunMain(t *testing.T) {
    TestRun = true
    go main()          // server-i ayrı goroutine-da başlat
    <-StopTestServer   // stop siqnalını gözlə
    TestRun = false
}
```

`main.go` modifikasiyası — test üçün dayandırma kanalı:

```go
var (
    StopTestServer = make(chan bool)
    TestRun        = false
)

func main() {
    //...
    if TestRun {
        e.POST("/stop-test-server", func(ctx echo.Context) error {
            StopTestServer <- true
            return nil
        })
    }
    //…
}
```

**Axın:**
1. `go test` server-i `go main()` ilə başladır (Echo banner görünür: `⇨ http server started on [::]:8080`).
2. Server live-dır — curl ilə xarici testlər aparılır:

```bash
curl http://localhost:8080/health-check
# {"message":"Everything is good!"}
curl -XPOST http://localhost:8080/stop-test-server   # testi dayandır
```

3. Coverage real rəqəmləri: `coverage: 7.9% of statements in ./handlers`.

**Funksiya-səviyyəli coverage:**

```bash
go tool cover -func=cov.txt
# handlers/err.go:10:          Error ...
# handlers/health-check.go:13: HealthCheck ...
# total: (statements) 7.9%
```

**Nəyə lazımdır?**
- **Cucumber** kimi Go-dan olmayan davranış framework-lərində test yaza bilən QE komanda üzvlərinə real coverage rəqəmləri təqdim etmək.
- CI/CD pipeline-a qoşulmaq — instrumented server real DB yanında işləyib gerçek integration coverage verir.

---

### 7. Coverage haqqında dürüst qeydlər

- Coverage **yaxşıdır**: irəli-geri hərəkət edən keyfiyyət göstəricisi.
- Coverage **pisdir**: developer-lər düzgün assertion olmadan rəqəmi şişirdə bilər — "bunk" testlər coverage artırır, amma keyfiyyət yox.
- **100% coverage tələbi** developer-i künc qısaltmalarına sövq edir; funksionallıq tələbi artdıqca keyfiyyət düşür — **riskə əsaslanan tarazlıq** vacibdir.

---

## Əsas terminlər

| Termin | İzah |
|--------|------|
| Unit testing | Ən kiçik kod vahidinin (funksiya/handler) izolyə şəkildə yoxlanması |
| Benchmark testing | Kodun sürəti/throughput-unun ölçülməsi (`ns/op`) |
| Behavior testing | İstifadəçi perspektivindən sistem davranışının yoxlanması (məs., cucumber) |
| Integration testing | Real xarici sistemlərlə (DB və s.) birlikdə yoxlama |
| Security testing | XSS, SQL injection, CORS kimi təhlükəsizlik deşiklərinin axtarışı |
| `_test.go` konvensiyası | Go test fayl adı şərti; `Test` prefiksi + `*testing.T` imzası |
| `httptest.ResponseRecorder` | `http.ResponseWriter`-in yadda saxlayan test implementasiyası |
| `e.ServeHTTP(w, r)` | Echo-nun real HTTP server olmadan çağrılması — testlərin açarı |
| Mocking | Interface arxılına saxta asılılıq davranışının qoyulması |
| `MockableDB` interface | `sql.DB`-nin üzlə bənzər interfeysi — mock üçün qapı |
| Monkey patching | Run-time-da funksiyanın əvəz edilməsi (paralelliyi itirir) |
| `b.N` | Benchmark dövrü sayı — framework tərəfindən kalibrlənir |
| `-coverpkg` | Testin coverage-u fərqli package-dən toplaması |
| `-coverprofile` | Coverage nəticələrinin fayla yazılması |
| `go tool cover -func` | Funksiya-səviyyəli coverage hesabatı |
| Instrumented build | Coverage ölçən kodu içində daşıyan icra edilə bilən build |
| `testing.T` | Test status idarəçisi: `Fail`, `Log`, `Error` |

---

## Praktik nəticə

1. **Handler testləri üçün şablon:** `echo.New()` → middleware + route → `httptest.NewRecorder()` + `http.NewRequest` → `e.ServeHTTP(w, r)` → status/body assertion. HTTP serverə ehtiyac yoxdur.
2. **DB asılılıqlarını interface-ə bağlayın:** `GetUserByUsername(db MockableDB, ...)` — beləcə unit test-də `MockDB`, production-da `*sql.DB` keçir; funksiya sahəli mock-lar xəta/ugur hallarını dəqiq simulyasiya edir.
3. **Mock-u ehtiyatla istifadə edin:** Over-mock coverage şişirdir, integration problemlərini gizlədir — mock-un yanında mütləq integration test olmalıdır.
4. **Benchmark yazmaq ucuzdur:** `BenchmarkXxx` + `for i := 0; i < b.N; i++` — dəyişikliklərin performance təsirini ədədlə görün.
5. **Real coverage üçün instrumented server:** `TestRun` flag + `/stop-test-server` endpoint + `-coverpkg` — curl və ya cucumber testləri üçün canlı, ölçən server.
6. **Coverage rəqəminə kor-koranə inanmayın:** assertion-lərin keyfiyyəti rəqəmdən vacibdir; keyfiyyət-risk tarazlığını şirkətlə müqavilə kimi müəyyənləşdirin.
7. **`go test -v ./handlers -cover`** — gündəlik inkişaf dövrəsində əsas əmr.

---

## Mənbə

- Kitab: *Echo Quick Start Guide* — J. Ben Huson, Packt Publishing, 2018 (ISBN 9781789340849)
- Chapter 7: "Testing Applications", book səh. 134–156
- PDF səhifələri: 142–161
- Kod: https://github.com/PacktPublishing/Echo-Essentials/tree/master/chapter7
- Video: https://goo.gl/PtZgpP
- Əlavə: https://bou.ke/blog/monkey-patching-in-go/
