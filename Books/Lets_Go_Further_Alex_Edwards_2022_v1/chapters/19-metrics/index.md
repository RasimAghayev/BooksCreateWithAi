# Chapter 19 — Metrics

## Bu fəsil nədən bəhs edir?

`expvar` paketi ilə runtime metrics: default (cmdline, memstats), custom
(NewString/NewInt/NewMap, Publish+Func), `db.Stats()` pool vəziyyəti,
request-level metrics middleware-i və custom ResponseWriter ilə status
kodlarının tutulması.

## Əsas fikirlər

### 1. expvar.Handler — GET /debug/vars
**Nədir:** Stdlib-in hazır metrics endpoint-i — JSON formatında.

```go
router.Handler(http.MethodGet, "/debug/vars", expvar.Handler())
```

**Default output:**
- `cmdline` → os.Args (proqram adı + CLI flag-lər) — hansı non-default
  parametrlərlə start olunduğunu göstərir
- `memstats` → runtime.MemStats snapshot:
  - `TotalAlloc` — heap-də kümülativ ayrılan bayt (AZALMIR)
  - `HeapAlloc` — cari heap baytı
  - `HeapObjects` — cari heap obyekt sayı
  - `Sys` — OS-dən alınmış ÜMUMİ yaddaş (heap, stack, internal)
  - `NumGC` — tamamlanmış GC dövrü sayı
  - `NextGC` — növbəti GC hədəf heap ölçüsü (HeapAlloc ≤ NextGC)

### 2. Custom metrics — New* + Publish
**Kitabdan kod nümunəsi:**
```go
expvar.NewString("version").Set(version)

expvar.Publish("goroutines", expvar.Func(func() any {
    return runtime.NumGoroutine()
}))

expvar.Publish("database", expvar.Func(func() any {
    return db.Stats() // sql.DBStats: in-use, idle, WaitCount...
}))

expvar.Publish("timestamp", expvar.Func(func() any {
    return time.Now().Unix()
}))
```

**Sub-kod izahı:**
- `expvar.NewString/NewFloat/NewInt/NewMap` → register + publish + pointer;
  concurrency-safe — runtime-da dəyişmək olar; eyni adla ikinci
  register → PANIC
- `expvar.Publish(name, expvar.Func(...))` → dinamik metric — hər sorğuda
  funksiya çağrılır (lazımi hesablama)
- Func `any` qaytarır — JSON-a encode OLUNMALI; olmasa item səssicə
  düşür və /debug/vars cavabı MALFORMED olur

**db.Stats() interpretasiyası (load test nəticəsi):**
- 118 goroutine, 11 in-use + 14 idle bağlantı
- `WaitCount` / `WaitDuration` → pool-un dolu olduğuna görə gözləmə sayı/vaxtı
  (normal production-da ~0 olmalıdır; nümunədə ~98ms orta gözləmə →
  MaxOpenConns artırılmalı)
- `MaxIdleTimeClosed` → ConnMaxIdleTime (15m) üzrə qapanmış bağlantılar

**Təhlükəsizlik (vacib!):** /debug/vars DoS üçün informasiya verir;
`cmdline` DSN kimi sensitive datanı AÇA bilər → production-da məhdudlaşdırın:
- `metrics:view` permission, YOXSA
- HTTP Basic auth, YOXSA
- Reverse proxy (Caddy) localhost-only (kitabın seçimi — Chapter 21-də)
- Default cmdline/memstats-i SİL MƏK mümkün deyil (open issue)

### 3. metrics() middleware — request-level
**Kitabdan kod nümunəsi:**
```go
func (app *application) metrics(next http.Handler) http.Handler {
    var (
        totalRequestsReceived           = expvar.NewInt("total_requests_received")
        totalResponsesSent              = expvar.NewInt("total_responses_sent")
        totalProcessingTimeMicroseconds = expvar.NewInt("total_processing_time_μs")
        totalResponsesSentByStatus      = expvar.NewMap("total_responses_sent_by_status")
    )
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        totalRequestsReceived.Add(1)

        mw := &metricsResponseWriter{wrapped: w}
        next.ServeHTTP(mw, r)

        totalResponsesSent.Add(1)
        totalResponsesSentByStatus.Add(strconv.Itoa(mw.statusCode), 1)

        duration := time.Since(start).Microseconds()
        totalProcessingTimeMicroseconds.Add(duration)
    })
}

// Zəncirin ƏN BAŞINDA:
return app.metrics(app.recoverPanic(app.enableCORS(app.rateLimit(app.authenticate(router)))))
```

**Sub-kod izahı:**
- expvar dəyişənləri middleware İNİSİALİZASİYASINDA bir dəfə (wrap zamanı)
  yaradılır — closure ilə hər request-də yenidən istifadə
- `Add(1)` — atomik artım (concurrency-safe)
- `time.Since(start).Microseconds()` — request işlənmə müddəti
- expvar.NewMap string-keyed → `strconv.Itoa(statusCode)`

**Törəmə metrikalar (iki /debug/vars snapshot arasında):**
- In-flight: `total_requests_received - total_responses_sent`
- RPS: `Δtotal_requests_received / Δtimestamp`
- Orta request vaxtı: `Δtotal_processing_time_μs / Δtotal_requests_received`

### 4. Status kodun tutulması — custom ResponseWriter
**Problem:** Stdlib-də cavabın status kodunu oxumaq üçün built-in YOL YOXDUR.

**Həll (Go 1.20+ ResponseController/Unwrap era):** ResponseWriter-ı wrap edən
struct:

**Kitabdan kod nümunəsi:**
```go
type metricsResponseWriter struct {
    wrapped       http.ResponseWriter
    statusCode    int
    headerWritten bool
}

func (mw *metricsResponseWriter) Header() http.Header {
    return mw.wrapped.Header()
}

func (mw *metricsResponseWriter) WriteHeader(statusCode int) {
    mw.wrapped.WriteHeader(statusCode) // əvvəlcə pass-through (panic halında
                                       // fərqli status gedə bilər)
    if !mw.headerWritten {
        mw.statusCode = statusCode
        mw.headerWritten = true
    }
}

func (mw *metricsResponseWriter) Write(b []byte) (int, error) {
    if !mw.headerWritten {
        mw.statusCode = http.StatusOK // WriteHeader çağrılmayıbsa default 200
        mw.headerWritten = true
    }
    return mw.wrapped.Write(b)
}

func (mw *metricsResponseWriter) Unwrap() http.ResponseWriter {
    return mw.wrapped // http.ResponseController unwrapping üçün
}
```

**Sub-kod izahı:**
- 3 interface method (Header/WriteHeader/Write) → http.ResponseWriter
  təmin olunur → middleware zəncirində şəffaf
- `headerWritten` → yalnız İLK status qeyd olunur; `Write`-dan əvvəl
  WriteHeader çağrılmayıbsa Go default 200 göndərir → biz də qeyd edirik
- WriteHeader pass-through-dan SONRA qeyd — invalid status panic-i fərqli
  statusa səbəb ola bilər
- **Unwrap()** → ResponseController (və digər wrap-lər) daxili
  ResponseWriter-ə çata bilsin — bunsuz Flush/Hijack kimi imkanlar
  itirilərdi

**Alternativ embedding:**
```go
type metricsResponseWriter struct {
    http.ResponseWriter        // embed — Header() avtomatik promote
    statusCode    int
    headerWritten bool
}
```
- Gain: Header() yazmırlıq; Loss: daha az explicit — dad məsələsidir

**Test nəticəsi (hey ilə):** rate limit aktiv → 4×201 + 196×429 →
`total_responses_sent_by_status: {"201": 4, "429": 196}`

### 5. Metrics ilə nə etməli (Additional Info)
- Spot-check (kiçik layihələr)
- Skript: periodically fetch + analiz + alerting
- **Prometheus** + visualization (real-time graphs)
- expvar platforma verir — inteqrasiya seçimi layihədən asılıdır

## Bu fəsildə layihəyə nə əlavə olundu? (project continuity)

**Dəyişdirilən:**
- routes.go → GET /debug/vars (expvar.Handler) + metrics zəncirin başında
- middleware.go → metrics() + metricsResponseWriter (4 method)
- main.go → version/goroutines/database/timestamp publish

## Əsas terminlər

- expvar — Go-nun standart exposed-variable (metrics) paketi
- MemStats — runtime yaddaş statistikası snapshot-ı
- GC (Garbage Collector) — heap təmizləyicisi (NumGC, NextGC)
- db.Stats() — connection pool statistikası (WaitCount, MaxIdleTimeClosed...)
- Custom ResponseWriter — status kodu qeyd edən wrapper (Unwrap dəstəyi)
- Response Promotion — embedding-lə method-un avtomatik gəlməsi
- In-flight Requests — gözləməkdə olan (received-sent) sorğular

## Praktik nəticə

expvar "sıfır asılılıqla observability" deməkdir: versiya, goroutine sayı,
pool vəziyyəti, request sayı/vaxtı/status payı — hamısı stdlib ilə. İki
vacib incəlik: (1) /debug/vars production-daq QORUNMALI (cmdline DSN
sızdırır); (2) status kodu üçün custom ResponseWriter lazımdır — Unwrap()
metodu isə müasir Go-nun (1.20+) wrapper zənciri ekosisteminə inteqrasiya
açarıdır. (Müəllim qeydi: Prometheus istifadə edəcəksinizsə, /debug/vars
yerinə /metrics + promhttp daha standart həlldir; expvar konsepti eyni
qalır.)

## Mənbə
Pages: 425-453 (raw 425-453)
