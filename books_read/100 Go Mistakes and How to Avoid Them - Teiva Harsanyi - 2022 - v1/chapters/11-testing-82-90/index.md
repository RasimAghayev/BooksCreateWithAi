# Chapter 11 — Testing (#82-#90)

## Bu chapter nədən bəhs edir?

Bu chapter test prosesinin 9 səhvini əhatə edir: testlərin kateqoriyalaşdırılması (build tags, env dəyişənləri, short mode), `-race` flaqı, icra rejimləri (parallel, shuffle), table-driven tests, sleep əsaslı flaky testlər, time API ilə bağlı testlər, `httptest`/`iotest` paketləri, dəqiq olmayan benchmark-lar (timer, micro-benchmark, compiler optimization, observer effect) və test imkanları (coverage, _test paketi, utility funksiyaları, setup/teardown).

---

## Əsas fikirlər

### #82: Not categorizing tests (Testlərin kateqoriyalaşdırılmaması)

**Testing pyramid:** Unit (pəyə — çoxsaylı, ucuz, sürətli, deterministik) → Integration → E2E (zirvə — mürəkkəb, yavaş).

**3 kateqoriyalaşdırma üsulu:**

**1. Build tags** (fayl səviyyəsində):

```go
// db_test.go:
//go:build integration     // (Go 1.17+: //go:build; köhnə: // +build)
package db

func TestInsert(t *testing.T) { /* ... */ }
```

```bash
go test .                      # yalnız tagsiz fayllar (unit)
go test --tags=integration .   # tagsiz + integration
```

Yalnız integration üçün unit fayllarına `//go:build !integration` (neqasiya).

**Çatışmazlıq (Peter Bourgon):** Skip olunan testlər haqqında siqnal YOXDUR — testlər yaddan çıxa bilər.

**2. Environment variables** (test səviyyəsində — skip aşkar görünür):

```go
func TestInsert(t *testing.T) {
    if os.Getenv("INTEGRATION") != "true" {
        t.Skip("skipping integration test")   // --- SKIP görünür!
    }
    // ...
}
```

**3. Short mode** (sürət kateqoriyası):

```go
func TestLongRunning(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping long-running test")
    }
    // ...
}
```

```bash
go test -short -v .
```

Üsulları birləşdirmək olar (unit/integration üçün tag/env + uzun testlər üçün short).

---

### #83: Not enabling the -race flag

**Race detector:** Kompilyasiya-time statik analiz YOX — **runtime aləti**. Kompilyator kodu instrumentasiya edir (bütün yaddaş çıxışlarını izləyir) → icra zamanı data race axtarır.

```bash
go test -race ./...
```

**Overhead:** Yaddaş 5-10x, icra vaxtı 2-20x → yalnız local test/CI (production: YOX, canary istisna).

**Warning-in oxunması:**

```
WARNING: DATA RACE
Write at 0x00c000026078 by goroutine 7:   ← yazan goroutine + kod sətri
Previous read ... by main goroutine:      ← oxuyan + sətir
Goroutine 7 (running) created at:          ← yaradılma nöqtəsi
```

**Daxili mexanizm:** Vector clock-lar (distributed sistemlərdə də istifadə olunan qismi sıralama strukturu) — hər goroutine yaradılışında clock, hər yaddaş/sinxronizasiya hadisəsində update, müqayisə ilə race aşkarı.

**Zəmanətlər:**
- False positive YOXDUR → warning = mütləq race.
- False negative MÜMKÜNDÜR → race yoxlama testini loop-da təkrarla (100 dəfə).

**İstisna tag:** `//go:build !race` — race-deteksiyadan faylı çıxarır.

---

### #84: Not using test execution modes (İcra rejimləri)

**Parallel (`t.Parallel`):**

```go
func TestFoo(t *testing.T) {
    t.Parallel()
    // ...
}
```

İcra modeli: sequensial testlər əvvəl bir-bir → sonra paralel testlər (`=== PAUSE` → `=== CONT`). Maksimum paralellik default GOMAXPROCS-a bərabər:

```bash
go test -parallel 16 .    # I/O-ağır uzun testlər üçün artır
```

**Shuffle (Go 1.17+, `-shuffle=on`):** Test sırasını randomlaşdırır — gizli asılılıqları (sıra asılı testlər, shared state) ifşa edir.

```bash
go test -shuffle=on -v .
# -test.shuffle 1636399552801504000  ← SEED çap olunur
go test -shuffle=1636399552801504000 -v .   # eyni sıranı reproduce et
```

CI-da fail lokalda reproduce üçün seed-i istifadə et.

---

### #85: Not using table-driven tests

**Problem:** Hal-hazırkı 5 test funksiyası — 55 simvollu adlar + struktur dublikasiyası (call → expected → müqayisə → error).

**Həll — subtests (`t.Run`) + map:**

```go
func TestRemoveNewLineSuffix(t *testing.T) {
    tests := map[string]struct {
        input    string
        expected string
    }{
        `empty`:                   {input: "", expected: ""},
        `ending with \r\n`:        {input: "a\r\n", expected: "a"},
        `ending with \n`:          {input: "a\n", expected: "a"},
        `ending with multiple \n`: {input: "a\n\n\n", expected: "a"},
        `ending without newline`:  {input: "a", expected: "a"},
    }
    for name, tt := range tests {
        t.Run(name, func(t *testing.T) {
            got := removeNewLineSuffixes(tt.input)
            if got != tt.expected {
                t.Errorf("got: %s, expected: %s", got, tt.expected)
            }
        })
    }
}
```

**Faydaları:** Oxunaqlı adlar; məntiq 1 dəfə — dəyişiklik/əlavə minimal; subtest-i tək icra: `go test -run=TestFoo/subtest_1`.

**Paralel subtest + loop dəyişəni (#63 tələsi):**

```go
for name, tt := range tests {
    tt := tt                        // KÖLGƏLƏMƏ ŞƏRT!
    t.Run(name, func(t *testing.T) {
        t.Parallel()
        // Use tt — hər closure-un ÖZ tt-si
    })
}
```

---

### #86: Sleeping in unit tests (Sleep-lü flaky testlər)

**Flaky test:** Kod dəyişmədən keçə bilən/fail edə bilən test — debugging baha, etimadsızlıq yaradır.

**Ssenari:** `getBestFoo` cavabı qaytarır + arxa fon goroutine Publish çağırır. Mock + `time.Sleep(10ms)` ilə yoxlama — **10 ms kifayət etməyə bilər** → flaky.

**Həll 1 — retry (assertion周期ik):**

```go
func assert(t *testing.T, assertion func() bool, maxRetry int, waitTime time.Duration) {
    for i := 0; i < maxRetry; i++ {
        if assertion() {
            return
        }
        time.Sleep(waitTime)
    }
    t.Fail()
}

assert(t, func() bool { return len(mock.Get()) == 2 }, 30, time.Millisecond)
```

Passiv 10 ms-dən yaxşıdır — uğur halında 1 ms-dən sonra keçir. (testify `Eventually` hazır təklif edir.)

**Həll 2 — channel sinxronizasiyası (ƏN YAXŞI — tam deterministik):**

```go
type publisherMock struct{ ch chan []Foo }
func (p *publisherMock) Publish(got []Foo) { p.ch <- got }

func TestGetBestFoo(t *testing.T) {
    mock := publisherMock{ch: make(chan []Foo)}
    defer close(mock.ch)
    h := Handler{publisher: &mock, n: 2}
    foo := h.getBestFoo(42)
    // Check foo
    if v := len(<-mock.ch); v != 2 {   // Publish baş verənə qədər blok
        t.Fatalf("expected 2, got %d", v)
    }
}
```

(Tezliklə gözləməmək üçün select + time.After timeout əlavə et.)

**Prioritet:** Sinxronizasiya > retry > passiv sleep. Sinxronizasiya mümkün deyilsə — dizayn yenidən düşün.

---

### #87: Not dealing with the time API efficiently (Time API testləri)

**Problem:** `TrimOlderThan` içində `time.Now()` → test `time.Now()` əsaslı event-lərlə → maşın məşğul olsa flaky.

**Həll 1 — unexported funksiya asılılığı (dependency injection):**

```go
type now func() time.Time

type Cache struct {
    mu     sync.RWMutex
    events []Event
    now    now                    // asılılıq sahəsi
}

func NewCache() *Cache {
    return &Cache{events: make([]Event, 0), now: time.Now}   // production: real
}

// Test — sabit vaxt inyeksiyası:
cache := &Cache{now: func() time.Time {
    return parseTime(t, "2020-01-01T12:00:00.06Z")   // DETERMİNİSTİK
}}
```

Çatışmazlıq: unexported → xarici paket testində əlçatmaz (#90 ilə əlaqə). `time.After` üçün əlavə `after` asılılığı və ya `Now/After` interfeysi.

**Həll 2 — API-nin dəyişdirilməsi (ən sadə):**

```go
// Client cari vaxtı verir:
func (c *Cache) TrimOlderThan(now time.Time, since time.Duration) { /* ... */ }
// Və ya bir daha da sadə — tək arqument:
func (c *Cache) TrimOlderThan(t time.Time) { /* ... */ }
cache.TrimOlderThan(time.Now().Add(time.Second))
```

**Anti-pattern — qlobal dəyişən:** `var now = time.Now` — mutable shared state → testlər izolyasiyası pozulur, paralel icra mümkünsüz. Struct asılılığı seç.

---

### #88: Not using testing utility packages (httptest / iotest)

#### httptest — server testi

```go
func Handler(w http.ResponseWriter, r *http.Request) {
    w.Header().Add("X-API-VERSION", "1.0")
    b, _ := io.ReadAll(r.Body)
    _, _ = w.Write(append([]byte("hello "), b...))
    w.WriteHeader(http.StatusCreated)
}

func TestHandler(t *testing.T) {
    req := httptest.NewRequest(http.MethodGet, "http://localhost",
        strings.NewReader("foo"))       // sorğu
    w := httptest.NewRecorder()          // cavab yazıyıcı
    Handler(w, req)                      // birbaşa çağırış — transport YOX
    // assertion-lar: header, body, status
}
```

Diqqət: transport (HTTP özü) test olunmur — yalnız handler loqikası.

#### httptest — client testi

```go
func TestDurationClientGet(t *testing.T) {
    srv := httptest.NewServer(           // LOKAL HTTP server (bir neçə ms!)
        http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            _, _ = w.Write([]byte(`{"duration": 314}`))
        }),
    )
    defer srv.Close()
    client := NewDurationClient()
    duration, err := client.GetDuration(srv.URL, ...)
    // assertion...
}
```

Variantlar: `NewTLSServer`, `NewUnstartedServer` (lazy start).

#### iotest — reader/writer testləri

```go
// Custom reader-in düzgünlüyü:
err := iotest.TestReader(&LowerCaseReader{reader: strings.NewReader("aBcDeFgHiJ")},
    []byte("acegi"))
```

**Xəta-tolerantlıq alətləri:**

| Alət | Davranış |
|------|----------|
| `iotest.ErrReader` | Verilmiş error qaytarır |
| `iotest.HalfReader` | Sorğulananın yarısını oxuyur |
| `iotest.OneByteReader` | Hər oxumada 1 bayt |
| `iotest.TimeoutReader` | 2-ci oxumada xəta (sonrakılar OK) |
| `iotest.TruncateWriter` | n baytdan sonra səssiz dayanır |

```go
// foo funksiyası io.ReadAll istifadə edirsə — TimeoutReader testi FAIL edəcək
// (io.ReadAll xətaları çatdırır). Custom readAll(r, retries) yaz → test keçər.
```

---

### #89: Writing inaccurate benchmarks (Dəqiqsiz benchmark-lar)

**Benchmark əsasları:**

```go
func BenchmarkFoo(b *testing.B) {
    for i := 0; i < b.N; i++ {   // b.N dəyişkən — benchtime-a (default 1s) uyğunlaşır
        foo()
    }
}
```

#### a) Timer reset/pause unutması

```go
// Setup loop-dan əvvəl:
expensiveSetup()
b.ResetTimer()             // setup vaxtı/allocasiyası sıfırlanır

// Setup hər iterasiyada:
for i := 0; i < b.N; i++ {
    b.StopTimer()          // pauza
    expensiveSetup()
    b.StartTimer()         // davam
    functionUnderTest()
}
```

Qeyd: funksiya setup-dan çox sürətlidirsə benchmark 1 saniyədən ÇOX çəkə bilər (vaxt yalnız funksiya hesablanır) → `-benchtime` azalt.

#### b) Micro-benchmark yanlış pressumpları

`atomic.StoreInt32` vs `StoreInt64` — sıra dəyişdikdə "qalib" dəyişir! Təsir edən: maşın aktivliyi, power management, thermal scaling, cache alignment.

**Həllər:**
- `-benchtime` artır (böyük ədədlər qanunu — gözlənilən dəyərə yaxınlaşma).
- **perflock** — benchmark-ı CPU-nun 70%-i ilə işlət (OS-ə 30%).
- **benchstat** (`golang.org/x`) — statistik müqayisə:

```bash
go test -bench=. -count=10 | tee stats.txt
benchstat stats.txt
# AtomicStoreInt32-4   5.10ns ± 1%
# AtomicStoreInt64-4   5.10ns ± 1%   ← əslində EYNİDİR
```

Həm də production sistemi benchmark maşınından fərqli ola bilər — nəticələri ehtiyatla köçür.

#### c) Compiler optimization tələsi (Go issue 14813)

```go
// PİS — 0.2858 ns/op (~1 saat dövrü — MÜMKÜNSÜZ!):
func BenchmarkPopcnt1(b *testing.B) {
    for i := 0; i < b.N; i++ {
        popcnt(uint64(i))    // inline + side-effect yox → kompilyator BOŞ loop edir
    }
}

// DÜZGÜN — 1.993 ns/op:
var global uint64

func BenchmarkPopcnt2(b *testing.B) {
    var v uint64
    for i := 0; i < b.N; i++ {
        v = popcnt(uint64(i))    // 1) lokal dəyişənə yaz
    }
    global = v                    // 2) sonuncunu QLOBAL-a yaz
}
```

Pattern: hər iterasiyada LOKAL-a, sonda bir dəfə QLOBAL-a — inline+dead-code-optimize qarşı.

#### d) Observer effect

512 vs 513 kolonlu matrisin ilk 8 kolon cəmi — 513 ~50% "sürətli" çıxdı!

**Səbəb:** Eyni matris minlərlə dəfə istifadə olunduğundan CPU cache-də artıq var → cache miss azalır. Ölçülən: yeni matris üzərindəki iş YOX, cache-dəki matris.

**Həll — hər iterasiyada yeni data:**

```go
for i := 0; i < b.N; i++ {
    b.StopTimer()
    s := createMatrix512(rows)    // yeni matris
    b.StartTimer()
    sum = calculateSum512(s)
}
```

Nəticə: 33,547 vs 35,507 ns/op — real fərq yoxdur. CPU-bound micro-benchmark-larda məcburi qayda.

---

### #90: Not exploring all the Go testing features (Test imkanları)

**a) Code coverage:**

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out    # brauzerdə sətir-sətir
# Başqa paketdən test olunmuş kod üçün:
go test -coverpkg=./... -coverprofile=coverage.out ./...
```

Xəbərdarlıq: 100% coverage ≠ bugsuz — testlərin MƏNTİQİ hər hansı statik həddən vacibdir.

**b) `_test` paketi (davranış testi):**

```go
// counter_test.go — xarici paket:
package counter_test

import ("testing"; "myapp/counter")

func TestCount(t *testing.T) {
    if counter.Inc() != 1 {   // YALNIZ exported API
        t.Errorf("expected 1")
    }
}
```

Unexported elementlərə (count dəyişəni) çıxış YOXDUR → test implementation detail-a YOX, açıq davranışa fokuslanır → refactor-da test dəyişmir.

**c) Utility funksiyası + `t` ötürməsi:**

```go
// Error qaytaran variant əvəzinə:
func createCustomer(t *testing.T, someArg string) Customer {
    // Create customer
    if err != nil {
        t.Fatal(err)          // test birbaşa fail
    }
    return customer
}
customer := createCustomer(t, "foo")   // qısa, təmiz
```

**d) Setup/teardown:**

```go
// Test-səviyyə: defer / t.Cleanup:
func TestMySQLIntegration(t *testing.T) {
    setupMySQL()
    defer teardownMySQL()
    db := createConnection(t, "tcp(localhost:3306)/db")
    // ...
}
func createConnection(t *testing.T, dsn string) *sql.DB {
    db, err := sql.Open("mysql", dsn)
    if err != nil { t.FailNow() }
    t.Cleanup(func() { _ = db.Close() })   // avtomatik bağlanma (LIFO sıra)
    return db
}

// Paket-səviyyə: TestMain:
func TestMain(m *testing.M) {
    setupMySQL()
    code := m.Run()
    teardownMySQL()
    os.Exit(code)
}
```

---

## Əsas terminlər

| Termin | İzah |
|--------|------|
| Testing pyramid | Unit (baza) → Integration → E2E — yuxarı getdikcə mürəkkəb/yavaş |
| Build tag (`//go:build`) | Fayl-şərtli kompilyasiya — test kateqoriyası |
| `!integration` neqasiya | Yalnız-integration / yalnız-unit bölgüsü |
| Env dəyişənli skip | `t.Skip` — skip AŞKAR görünür |
| `testing.Short()` | `-short` rejimi — uzun testləri ayır |
| Race detector | Runtime instrumentasiya — vector clock əsaslı |
| False negative (race) | Real race-i qaçırma ehtimalı → loop-icra |
| `t.Parallel()` / `-parallel N` | Paralel testlər — sequensiallardan sonra |
| `-shuffle=on` + seed | Random sıra + reproduce üçün seed |
| Table-driven test | Map + `t.Run` subtests — dublikatsız çoxsaylı case |
| Subtest kölgələmə `tt := tt` | Paralel subtest loop-dəyişən tələsindən qorunma |
| Flaky test | Kod dəyişmədən keçə/fail ola bilən test |
| Retry assertion | `assert(t, cond, maxRetry, wait)` — passiv sleep-dən yaxşı |
| Channel sinxronizasiyası | Mock-də channel → tam deterministik gözləmə |
| `type now func() time.Time` | Vaxt asılılığının funksiya-tip inyeksiyası |
| `httptest.NewRecorder/NewRequest` | Handler testi — transport'suz |
| `httptest.NewServer` | Lokal server — client testi (real HTTP, ms sürət) |
| `iotest.TestReader` | Custom reader davranış yoxlaması |
| `iotest.TimeoutReader` və s. | Xəta-tolerantlıq simulatorları |
| `b.ResetTimer` / `Stop/StartTimer` | Setup-in nəticədən çıxarılması |
| benchstat / perflock | Statistik müqayisə / CPU limitləmə |
| Local→global pattern | Inline/optimize qarşı: lokal yaz + sonda qlobal |
| Observer effect | Təkrar istifadə cache qazandırır → hər iterasiyada yeni data |
| `-coverpkg` | Başqa paketdən test olunan kodun cover-i |
| `package x_test` | Xarici test paketi — davranış fokusu |
| `t.Cleanup` | Test sonu funksiyası (LIFO) |
| `TestMain(m)` | Paket-səviyyə setup/teardown |

---

## Praktik nəticə

1. **Testləri kateqoriyalaşdır:** Build tags (fayl) / env + `t.Skip` (aşkar) / `-short` (temp) — kombinasiya mümkün.
2. **Konkurrent kod = `-race` şərt:** Local/CI-də; warning 100% real; race-testlərini loop-da təkrarla (false negative ehtimalı).
3. **Uzun testlər `t.Parallel` + `-parallel N`:** I/O-ağır testlərdə N > GOMAXPROCS faydalı; `-shuffle` gizli asılılıqları ifşa edir (seed ilə reproduce).
4. **Eyni strukturlu testlər → table-driven:** Map + subtests + `tt := tt` (paralel subtestlərdə MÜTLƏQ).
5. **Sleep = flaky siqnalı:** Channel sinxronizasiyası (ən yaxşı) → retry assert (ikinci) → passiv sleep (qadağan); sinxronizasiya mümkünsə dizaynı dəyiş.
6. **`time.Now()` testlənməsi:** Funksiya-tip asılılıq (unexported `now`) və ya client-ötürmə API — qlobal dəyişən anti-pattern.
7. **`httptest`:** Handler → Recorder+Request; Client → NewServer (Docker mock-larından yüz dəfə sürətli).
8. **`iotest`:** Custom reader → TestReader; xəta-tolerantlıq → TimeoutReader/HalfReader/OneByteReader.
9. **Benchmark intizamı:** Reset/Stop timer; `-count=10` + benchstat (± statistika); local→global pattern (inline qarşı); hər iterasiyada yeni data (cache observer effect).
10. **Coverage və `_test` paketi:** `-coverpkg` ümumi mənzərə; xarici paket davranış-zorlamalı test; `t.Fatal`-lı utility-lər; `t.Cleanup` + `TestMain` setup/teardown.

---

## Mənbə

- Kitab: *100 Go Mistakes and How to Avoid Them* — Teiva Harsanyi, Manning, 2022 (ISBN 9781617299599)
- Chapter 11: "Testing", book səh. 262–298
- PDF səhifələri: 282–318
- İstinadlar: Peter Bourgon test kateqoriyaları (mng.bz/qYlr); httptest (pkg.go.dev/net/http/httptest); iotest (pkg.go.dev/testing/iotest); Go issue 14813; benchstat (golang.org/x); perflock; testify Eventually
