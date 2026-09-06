# Chapter 11 — Наблюдаемость

## Bu chapter nədən bəhs edir?

Sistemin daxili vəziyyətini xarici müşahidələrdən anlama imkanı — gözlənilirlik (observability) — izah edir. Üç sütun: jurnal qeydləri (logs), metrikalar (metrics), izlər (traces). OpenTelemetry standartı, strukturlaşdırılmış jurnalçılıq (Zap), Prometheus metrikaları, avtomatik instrumentasiya və izləmə (tracing) mexanizmləri əhatə edilir.

## Əsas fikirlər

### 1. Gözlənilirlik (Observability) və Monitorinq arasındakı fərq
**Nədir:** Gözlənilirlik — sistemin daxili vəziyyətini xarici müşahidələrdən anlama imkanı.

**Monitorinq vs Gözlənilirlik:**
- **Monitorinq:** Gözlədiyiniz şeyin (bəlli mətrdə doğru getməməsi) olduğunu sizə deyir
- **Gözlənilirlik:** Gözləmədiyiniz, bilmədiyiniz suallara cavab verə bilmək — "why is it doing that?" sualına cavab

**Nəyə lazımdır:** Bulud sistemləri çox mürəkkəb və qeyri-müəyyən olduğu üçün, ancaq monitorinq ilə kifayətlənmək olmaz. Gözlənilirlik, sistemin gələcəkdə bilinməyən problemlərini aşkar etmək imkanı verir.

### 2. Üç sütun: Logs, Metrics, Traces
**Nədir:** Gözlənilirliyin üç əsas elementi.

**Necə işləyir:**
- **Logs (Jurnal qeydləri):** Hadisələr (events) — nə oldu, nə vaxt, harada, hansı sorğu. Məsələn: "HTTP 500, path=/api, latency=250ms"
- **Metrics (Metrikalar):** Yığmış ölçmələr (aggregated measurements) — sayğaclar (counters), göstəricilər (gauges), histogramlar (distributions)
- **Traces (İzlər):** Sorğunun xidmətlər arasında səyahətini izləyir — distributed tracing

### 3. Strukturlaşdırılmış Jurnalçılıq (Structured Logging)
**Nədir:** Mətn qeydləri əvəzinə açar/dəyər (key/value) cütləri ilə log yazmaq.

**Necə işləyir:**
```go
import "go.uber.org/zap"

logger, _ := zap.NewProduction()
defer logger.Sync()

logger.Info("Testing sampling",
    zap.Int("index", i),
    zap.String("path", "/api/v1/resource"),
)
```

**Üstünlükləri:**
- Filtreləmə və axtarış asan — açar üzrə sorğu
- Maşın tərəfindən emal (machine-parseable)
- Statistik və dashboard yaratmaq asan

**Zap paketi:**
- `zap.NewProduction()` — JSON format, yüksək performans
- `zap.NewDevelopment()` — insan oxunabilən format
- `zap.SugaredLogger` — rahat, amma bir qədər yavaş
- `zap.Logger` — maksimum performans

### 4. OpenTelemetry standartı
**Nədir:** Vendor-neytral (satıcıya bağlı olmayan) gözlənilirlik standartı — izləmə (traces), metrikalar və loglar üçün.

**Necə işləyir:**
- `go.opentelemetry.io/otel` — əsas API
- `go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp` — `net/http` üçün avtomatik instrumentasiya
- `go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc` — gRPC üçün avtomatik instrumentasiya
- `go.opentelemetry.io/otel/exporters/metric/prometheus` — Prometheus üçün metrika ixracatçısı

### 5. İzlər (Traces) və Spans
**Nədir:** Sorğunun paylanmış sistemdəki yolunu izləmək.

**Necə işləyir:**
- **Trace:** Tam sorğu səyahəti — istemciden servisə, oradan 3-cü xidmətə, geri qayıdış
- **Span:** Tək əməliyyat (operasiya) — məsələn, DB sorğusu, HTTP call
- **Span Events:** Span daxilindəki hadisələr — məsələn, "mutex əldə edildi"
- **Span Attributes/Labels:** Span haqqında metadata — `label.Int("db.rows", 42)`, `label.String("db.system", "postgresql")`

**Kitabdan kod nümunəsi (span events və labels):**
```go
ctx, span := tracer.Start(ctx, "GetResource")
defer span.End()

span.AddEvent("Acquiring mutex lock",
    label.String("lock.name", "resource_cache"),
    label.Int("lock.id", 42),
)

span.SetAttributes(
    label.String("resource.name", "database"),
    label.Int("resource.size", 1024),
)
```

### 6. Metrikalar (Metrics)
**Nədir:** Sistemin sabit və ya dəyişən göstəriciləri üzərində toplanmış ölçmələr.

**Necə işləyir:**
- **Counter (sayğac):** Yalnız artan dəyər — sorğu sayı, xəta sayı
- **Gauge (ölçü):** Dəyişən dəyər — yaddaş istifadəsi, aktiv bağlantı sayı
- **Histogram (histoqram):** Dəyərlərin paylanması — sorğu gecikməsi, ölçü paylanması

**Prometheus ixracatçısı:**
```go
import "go.opentelemetry.io/otel/exporters/metric/prometheus"

exporter, err := prometheus.NewExportPipeline()
meter := metric.Must(exporter.Meter("github.com/.../myapp"))
requestCount := metric.Must(meter.NewFloat64Counter("http.requests.total"))
```

### 7. Dinamik Nümunə və Log Təmizliyi
**Nədir:** Yüksək yük vəziyyətində logların və izlərin aşırı artmasını əngəlləmək.

**Necə işləyir:**
- **Dinamik nümunə (dynamic sampling):** Bəzi sorğuların 100%-ə, digərlərinin 1%-ə izlənməsi
- **Log təmizliyi (log filtering):** Səviyyə əsaslı filtrləmə — DEBUG yalnız inkişaf mühitində, PRODUCTION-da yalnız ERROR+
- **Həddindən artıq kardinallıq (cardinality):** Hər unikal dəyər üçün yeni seriya (series) yaratmaq — Prometheus-u böyük məhdudiyyətə qoya bilər

## Əsas terminlər
- Observability (gözlənilirlik) — daxili vəziyyəti xarici müşahidələrdən anlama
- Monitoring (monitorinq) — bəlli problemlərin aşkar edilməsi
- Structured logging (strukturlaşdırılmış jurnalçılıq) — key/value qeydləri
- Zap — yüksək performanslı Go jurnalçı kitabxanası
- OpenTelemetry — vendor-neytral gözlənilirlik standartı
- Trace (iz) — distribüed sorğu səyahəti
- Span — tək əməliyyat, trace-in hissəsi
- Span Events — span daxilindəki hadisələr
- Labels/Attributes — metadata
- Counter (sayğac) — yalnız artan metrika
- Gauge (ölçü) — dəyişən metrika
- Histogram (histoqram) — dəyər paylanması
- Prometheus — metrika toplama və izləmə sistemi
- Dynamic sampling (dinamik nümunə) — yükə görə izləmə dərəcəsi
- Cardinality (kardinallıq) — seriya sayının artması riski

## Praktik nəticə
Gözlənilirlik (observability) bulud sistemlərinin "nə baş verir?" sualını cavablandırır. Üç sütun (logs, metrics, traces) birgə işləyir — loglar hadisələri göstərir, metrikalar tendensiyaları izləyir, izlər sorğuların səyahətini izləyir. OpenTelemetry bu üçünü birləşdirən standartdır. Zap ilə strukturlaşdırılmış loglar və Prometheus ilə metrikalar, avtomatik instrumentasiya ilə izlər — bu alətlər modern Go xidmətləri üçün vacibdir.

## Mənbə
Pages: 354-409 (PDF səh. 354-409)
