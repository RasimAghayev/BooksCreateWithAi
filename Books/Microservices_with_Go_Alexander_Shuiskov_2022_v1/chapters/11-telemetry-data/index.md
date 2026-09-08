# Chapter 11 — Collecting Service Telemetry Data (səh. 215-250)

## Bu fəsil nədən bəhs edir?

Observability-nin üç sütunu: logs, metrics, traces. zap ilə strukturlaşdırılmış
logging, tally/Prometheus ilə metrics, OpenTelemetry + Jaeger ilə distributed
tracing.

## Əsas fikirlər

### 1. Telemetry (Telemetriya) anlayışı
- **Logs** — hadisə mesajları; **Metrics** — kəmiyyət ölçmələri (RPS, latency);
  **Traces** — request-in servislər arası səyahəti
- İmmutabledir; birlikdə güclü informasiya mənbəyidir
- **White-box** (daxili data ilə) vs **Black-box** (yalnız xarici göstəricilər,
  məs. health API) monitoring
- İstifadə: trend analizi, semantik qraf, anomaliya aşkarı, hadisə korrelyasiyası
- Çətinlik: böyük datasetlər, xüsusi tooling, mürəkkəb setup

### 2. Logging
Built-in `log` paketi: sadə, amma structured logging YOX, level YOX, Errorf YOX.

**Structured logging (Strukturlaşdırılmış log):** JSON + level + sahələr:
```json
{"level":"info","time":"...","message":"Service started","service":"metadata"}
```
**Level-lər:** Info / Error / Warning / Fatal / Debug (default bağlı — disk
qənaəti; troubleshooting-də açılır). "Connection terminated" — levelsiz şərh
mümkün deyil!

**Kitabxanalar:** zap (ən sürətli — Uber; Logger + SugaredLogger), zerolog,
go-kit/log, apex/log (Elasticsearch/Greylog built-in), log15 (yavaş).
Kitab **zap** seçir.

**zap praktikası:**
```go
logger, _ := zap.NewProduction()   // JSON, debug off, stacktrace
logger.Info("Started", zap.String("serviceName", "metadata"))
logger.Info("Timeout", zap.Duration("timeout", 10*time.Second))
// struct üçün String() metodu + zap.Stringer
logger = logger.With(zap.String("endpoint", "PutRating"))  // sub-logger
logger.With("component": "ratingController")  // komponent başına logger
```

**Saxlama (Elastic Stack/ELK):** Logstash (toplama) → Elasticsearch
(indeksləmə, 10ms-də TB-lər) → Kibana (UI). Retention policy mütləq.

**Best practices:**
1. **İnterpolasiyalı string YOX** — `zap.String("userId", id)` sahə kimi
   (parse/search asanlaşır)
2. **Standartlaşdır** — ortaq field konstantları (`logging.FieldService`);
   root logger-ə service adı; handler başında endpoint sahəsi
3. **Periodik review** — PII (ad, SSN) logda OLMASIN; debug data çoxalmır?
4. **Retention** — disk dolmasın
5. **Mənbəyi göstər** — service/component/endpoint/file sahələri

### 3. Metrics
- Logs request başına 1 hadisə = 1M RPS-də 1M event/s — qeyri-mümkün!
- **Metric:** aqreqasiyalı dəyər; **time series** = timestamp + value + tags

**Növlər:**
| Tip | Nə ölçür | Nümunə |
|---|---|---|
| **Counter** | kumulyativ dəyər/dəyişim | request sayı, error sayı |
| **Gauge** | tək skalyar dəyər | goroutine sayı, aktiv connection, boş disk |
| **Histogram** | bucket-lər üzrə paylanma | latency (0-100ms, 100-200ms...), yaş qrupları |

**Alətlər:** **Prometheus** (Go-da yazılıb; counter/gauge/histogram;
labels; **PromQL**: `http_requests_total{environment="production",method!="GET"}`),
Graphite (Carbon + Whisper + Graphite-web; Grafana inteqrasiyası).

**Kitabxanalar:** tally (Uber; Prometheus/StatsD/M3; seçilən), rcrowley/go-metrics,
go-kit/metrics.

**tally praktikası:**
```go
scope, _ := tally.NewRootScope(tally.ScopeOptions{
    Tags: map[string]string{"service": "rating"}, Reporter: reporter}, time.Second)
counter := scope.Counter("request_count"); counter.Inc(1)
gauge := scope.Gauge("active_user_count"); gauge.Update(n)
timer := scope.Timer("operation_latency"); sw := timer.Start(); defer sw.Stop()
histogram := scope.Histogram("user_age", tally.MustMakeLinearValueBuckets(0,1,130))
scope.Tagged(map[string]string{"operation":"put"})  // əlavə tag
```
- Scope-lar iyerarxikdir — parent tagləri avtomatik miras alır

**Best practices:**
1. **Tag cardinality** — UUID/object ID tag OLMASIN (indeks partlayır);
   service/endpoint/city kimi aşağı-cardinality OK
2. **Ad standartı** — api_errors vs api_request_errors qarışıqlığı YOX
3. **Retention** — Prometheus default 15 gün; `--storage.tsdb.retention.time=60d`

### 4. Distributed Tracing (Paylanmış İzləmə)
- **Span:** əməliyyat (ad, start, end, tags, logs); **parent/child span**
  iyerarxiyası — GetMovieDetails → GetMetadata (100ms) + GetAggregatedRating (1100ms)
- İstifadə: call analizi, xəta analizi, call path (kod bilmədən sistemi anla),
  əməliyyat parçalanması

**Context propagation (kontekst yayımı):** trace-in SİRİ!
- Servis A → B → C: `ctx.reqId` header-i ötürülür; hər servis eyni request-i
  tanıyır
- Go-da: `context.WithValue(ctx, key, value)` → arqument kimi ötür
- Addımlar: parent span yarat → context ilə uşaqlara ötür → uşaq span parent
  ID saxlayır → hamı report edir

**OpenTelemetry + Jaeger (Go-da yazılıb):**
```go
// config: jaeger.url (base.yaml) → http://localhost:14268/api/traces
tp := tracing.NewJaegerProvider(url, serviceName)  // jaeger exporter + batcher
defer tp.Shutdown(ctx)
otel.SetTracerProvider(tp)
otel.SetTextMapPropagator(propagation.TraceContext{})
// client: grpc.Dial(..., grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()))
// server: grpc.NewServer(grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()))
```
- gRPC interceptorləri avtomatik span yaradır — əlavə kod YOX
- Jaeger: `docker run jaegertracing/all-in-one` → UI :16686
- **Manual instrumentasiya** (DB repo):
```go
const tracerID = "metadata-repository-memory"
_, span := otel.Tracer(tracerID).Start(ctx, "Repository/Get")
defer span.End()
```
- Qayda: >50ms çəkən, şəbəkə/DB/IO əməliyyatları — trace namizədi

## Termindirmə (AZ)
- Telemetry — Telemetriya (servis performans datası)
- Structured Logging — Strukturlaşdırılmış Loglama
- Log Level — Log Səviyyəsi
- Counter/Gauge/Histogram — Saya/Göstərici/Histoqram
- Cardinality — Kardinalite (unikal dəyərlərin sayı)
- Span — Aralıq (bir əməliyyatın izi)
- Context Propagation — Kontekst Yayımı
- Distributed Tracing — Paylanmış İzləmə
- Retention — Saxlama Müddəti

## Kviz sualları
1. Niyə request sayını loglarla deyil, metric-lə ölçürük? (Yüksək RPS-də
   hər request üçün log = data partlayışı; metric aqreqasiya edir)
2. Tag-də UUID olmağın problemi? (Yüksək cardinality — hər dəyər indekslənir,
   TSDB-nin throughput-u düşür)
3. Parent-child span əlaqəsi necə qurulur? (Parent span ID context üzərindən
   uşaq çağırışlara ötürülür)
4. zap kitabxanasının 2 logger növü? (Logger — max performans;
   SugaredLogger — əlavə funksiyalar)
