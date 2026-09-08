# Chapter 11 — Telemetry (Telemetriya)

## Bu chapter nədən bəhs edir?
Tətbiqin müşahidə olunması (observability) üçün 3 sütuna: Logs (struktur log-lar — slog vs zap), Traces (runtime/trace + distributed tracing) və Metrics (Prometheus metric tipləri). Sonda hamısı vendor-neutral OTel (OpenTelemetry) çətiri ilə birləşdirilir.

## Əsas fikirlər

### 1. Logs — struktur loglama
Plain-text yerinə struktur (adətən JSON) format: axtarış, filtr, analiz asanlaşır.

**Standart kitabxana — slog (Go 1.21+, log/slog):**
```go
handler := slog.NewJSONHandler(os.Stdout)
logger := slog.New(handler)
logger.Info("A group of walrus emerges from the ocean",
    slog.Attr("animal", "walrus"), slog.Attr("size", 10))
```
**Sub-kod izahı:** handler JSON formatlayır, os.Stdout-a yazır; slog.Attr key-value cütlüyüdür — log mesajı maşın-oxunarlı struktur data ilə zənginləşir. Xarici asılılıq tələb etmir.

### 2. slog vs zap
| | slog | zap (uber) |
|---|---|---|
| Mənbə | standart kitabxana | xarici |
| Performans | yaxşı | ən sürətlilərdən |
| Konfiqurasiya | sadə | çox detallı (encoder/core/sink) |
| context inteqrasiyası | təbii | — |

**zap nümunəsi (konfiqurasiya səviyyəsi):**
```go
encoderConfig := zapcore.EncoderConfig{
    MessageKey: "message", LevelKey: "level",
    EncodeLevel: zapcore.CapitalLevelEncoder,
    TimeKey: "time", EncodeTime: zapcore.ISO8601TimeEncoder,
    CallerKey: "caller", EncodeCaller: zapcore.ShortCallerEncoder,
}
consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
core := zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zap.InfoLevel)
logger := zap.New(core)
sugar := logger.Sugar()                          // rahat, amma yavaş variant
sugar.Infow("A group of walrus emerges", "animal", "walrus", "size", 10)
```
**Sub-kod izahı:** encoder formatı, sink (çixış yeri), level filtri `NewCore`-da birləşir; `Sugar()` — key-value arqumentlərlə rahat loglama API-si (non-sugared versiyadan yavaşdır).

### 3. Debug vs Monitoring logları
| | Debug | Monitoring |
|---|---|---|
| Məqsəd | xəta diaqnostikası | production sağlamlıq müşahidəsi |
| Detallıq | çox verbose (stack trace, dəyişənlər) | az detal, balanslı |
| Ömür | müvəqqəti | daimi |
| Auditoriya | developer-lər | dev + ops + admin |

**Nə log etmək:** xətalar (stack trace ilə), system state dəyişiklikləri, istifadəçi əməlləri, performans metrikaları (metrics server yoxdursa), security event-lər, kritik biznes transaksiyaları (audit yoxdursa).
**Nə log ETMƏMƏK:** parollar/PII/kredit kartı/token-lər, production-da verbose debug, artıq məlumat, iri binary data, sanitasiya edilməmiş istifadəçi girişi (injection riski).
**JSON vs structured text seçimi:** ELK/Splunk kimi alətlər var / nested data → JSON; insan oxuyur / performans kritik → structured text.

### 4. Traces — runtime/trace
İcra yolunu (goroutine-lər, heap, GC event-ləri) qeyd edir — log-dan fərqli olaraq proqramın fasiləsiz, detallı icra hesabatı.

**Baza nümunə:**
```go
f, _ := os.Create("trace.out")
defer f.Close()
err := trace.Start(f)          // qeyd başlayır
if err != nil { panic(err) }
defer trace.Stop()             // qeyd dayanır
```
**Analiz:**
```bash
go tool trace trace.out        # web UI (http://localhost:xxxx)
```
**Sub-kod izahı:** `trace.Start` bütün runtime event-lərini fayla yazır; `go tool trace` veb interfeys açır — goroutine analizi, heap, GC pause-lar görünür.

**HTTP handler üçün trace wrapper:**
```go
func TraceHandler(inner http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        ctx, task := trace.NewTask(r.Context(), r.URL.Path)
        defer task.End()
        trace.Log(ctx, "HTTP Method", r.Method)
        trace.Log(ctx, "URL", r.URL.String())
        inner(w, r.WithContext(ctx))   // ctx downstream-a ötürülür
    }
}
http.HandleFunc("/", TraceHandler(handler))
```
**Sub-kod izahı:** hər request üçün ayrıca trace task yaradılır; context downstream çağırılara ötürülür ki, tam request lifecycle izlənə bilsin.

**Effektiv tracing:** bütün proqramı yox, yalnız kritik hissələri trace et; middleware şəklində; anomaliyalara (uzun bloklanan goroutine, GC pause) diqqət yetir. Tracing-in performans overhead-i log-dan böyükdür — detallıq/xərc balansı lazımdır.

### 5. Distributed tracing — 4 konsept
1. **Unique identifier (trace ID)** — ilkin request-ə təyin olunur.
2. **Propagation** — ID bütün servis-lərə ötürülür (HTTP header, queue mesajı).
3. **Spans** — hər servis öz hissəsini qeyd edir (timestamp, servis adı, funksiyalar, xətalar).
4. **Collection & analysis** — mərkəzi sistem span-ları trace ID-yə görə birləşdirir.

**Faydaları:** bottleneck görünür; root cause analizi (xətanın real mənbəyi tapılır); performans optimizasiyası; microservice debug.
**Alətlər:** Zipkin, Jaeger, Honeycomb, Datadog — amma vendor lock-in riski var (→ OTel həlli).

### 6. Metrics — Prometheus
Sayısal ölçmələr: memory, CPU load, goroutine sayı və s.

**Counter nümunəsi (status koduna görə request sayı):**
```go
var requestsProcessed = prometheus.NewCounterVec(
    prometheus.CounterOpts{
        Name: "http_requests_processed",
        Help: "Total number of processed HTTP requests.",
    },
    []string{"status_code"},        // label — ölçü bölgüsü
)
func init() { prometheus.MustRegister(requestsProcessed) }

http.Handle("/metrics", promhttp.Handler())      // scrape endpoint
http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    code := http.StatusOK
    if time.Now().Unix()%2 == 0 { code = http.StatusInternalServerError }
    requestsProcessed.WithLabelValues(fmt.Sprintf("%d", code)).Inc()
    w.WriteHeader(code)
})
```
**Sub-kod izahı:** CounterVec label-lərlə counter saxlayır; handler hər request-də uyğun status label-i ilə artırır; Prometheus `/metrics` endpoint-dən scrape edir.

**Sorğular:** `http_requests_processed` (ümumi say), `rate(http_requests_processed[1m])` (dəqiqəlik rate).
**İcra:** prometheus.yml (scrape_interval: 15s) → `docker run -p 9090:9090 -v prometheus.yml:/etc/prometheus/prometheus.yml prom/prometheus`.

### 7. 4 metric növü
| Tip | Davranış | Nə vaxt |
|---|---|---|
| **Counter** | yalnız artır, restart-da sıfırlanır | hadisə sayı (request, error, signup), rate ölçmə |
| **Gauge** | istənilən istiqamətə dəyişir | cari səviyyə: memory, aktiv user, queue uzunluğu |
| **Histogram** | müşahidələri bucket-lərə sayır + cəm | distribusiya: latency percentil-ləri (p95) |
| **Summary** | sliding window quantile hesablayır | dəqiq dinamik quantile-lər (son trendlər) |

**Seçim qaydası:** saymaq → Counter; artıb-azalma → Gauge; distribusiya → Histogram; dəqiq dinamik quantile → Summary (Histogram-dan hesablamalı, daha bahalı).

### 8. OTel (OpenTelemetry)
CNCF altında vendor-neutral standart: API + SDK + exporter-lər. Go üçün birincil dəstəkli dillərdən: `go.opentelemetry.io/otel/trace`, `otel/metric`, `otel/propagation`.

**Faydaları:** vendor dəyişdirmə kod dəyişikliyi tələb etmir; asan instrumentasiya; vahid data format; aktiv community.
**Qeydi:** OTel-in Go Logs SDK hələ inkişafdadır → log üçün kitabda zap işlədilir. OTel "magic" observable etmir — sensorları (instrumentasiya) kodda özün yerləşdirirsən.

**OTel + otlphttp trace exporter qurulumu:**
```go
traceExporter, err := otlptrace.New(ctx, otlptracehttp.NewClient())
tp := sdktrace.NewTracerProvider(
    sdktrace.WithBatcher(traceExporter),           // topyekün yığım göndərişi
    sdktrace.WithResource(resource.NewWithAttributes(
        semconv.SchemaURL,
        semconv.ServiceNameKey.String("ExampleService"),  // servis kimliyi
    )),
)
otel.SetTracerProvider(tp)
```
**Handler-də span:**
```go
func exampleHandler(w http.ResponseWriter, r *http.Request) {
    _, span := otel.Tracer("example-tracer").Start(r.Context(), "handleRequest")
    defer span.End()
    zap.L().Info("Handling request")
    w.Write([]byte("Hello, World!"))
}
http.Handle("/", otelhttp.NewHandler(http.HandlerFunc(exampleHandler), "Example"))
```
**Sub-kod izahı:** TracerProvider OTLP HTTP exporter-ə batch rejimdə trace göndərir; otelhttp middleware HTTP handler-i avtomatik instrument edir; handler içində manual span yaradılır. Backend (Collector) docker-compose ilə qalxır (ch11/otel/).

## Əsas terminlər
- Observability (müşahidə olunma) — log + trace + metric üçlüyü
- Structured logging (struktur loglama) — JSON/key-value format
- slog / slog.NewJSONHandler / slog.Attr — Go 1.21 standart log
- zap / zapcore.EncoderConfig / NewCore / Sugar() / Infow
- runtime/trace / trace.Start / trace.Stop / go tool trace
- trace.NewTask / trace.Log — task səviyyəli instrumentasiya
- Distributed tracing: trace ID, propagation, span, collection
- Zipkin / Jaeger / Honeycomb / Datadog — tracing backend-lər
- Prometheus: Counter, Gauge, Histogram, Summary
- NewCounterVec / WithLabelValues / promhttp.Handler / /metrics
- rate() — vaxt pəncərəsi üzrə artım sürəti sorğusu
- OTel / OTLP / otlptracehttp / TracerProvider / WithBatcher
- semconv.ServiceNameKey — servis identifikasiyası resurs atributu
- otelhttp.NewHandler — HTTP avtomatik instrumentasiyası
- Vendor lock-in / vendor neutrality

## Praktik nəticə
1. Log formatı alətə görə seç: ELK/Splunk → JSON; konsol oxunuşu/performance → struktur text.
2. Heç vaxt parol/PII/token log etmə; istifadəçi girişini sanitasiya etdikdən sonra logla.
3. Tracing-i yalnız kritik hissələrdə işə sal — trace faylları böyükdür, overhead log-dan yuxarıdır.
4. HTTP server-də trace ctx-i downstream-a ötür; distributed sistemlərdə trace ID-ni header-lərlə propagate et.
5. Metric tipi təbiətə görə: say → Counter, səviyyə → Gauge, latency distribusiya → Histogram.
6. Yeni layihədə OTel ilə başla — backend dəyişərsə (Jaeger → Zipkin) instrumentasiya kodu qalır.

## Mənbə
Pages: 221-245 (PDF səh. 242-266)
