# Chapter 21 — Go in the Cloud (səh. 646-680)

## Bu fəsil nədən bəhs edir?

Production-a hazırlıq: Prometheus monitoring (counter/gauge/histogram/
summary, /metrics endpoint), OpenTelemetry (tracing + structured logs),
Docker containerization (multi-stage, scratch, CGO_ENABLED=0) və
Kubernetes orchestration konseptləri.

## Əsas fikirlər

### 1. Monitoring vs observability
- **Monitoring** — əvvəlcədən təyin edilmiş metrik/həddlərlə sistem
  sağlamlığının izlənməsi + alerting
- **Observability** — daha araşdırıcı: sistemin daxili davranışının
  comprehensive anlaşılması (debugging üçün)

### 2. Prometheus
**Nədir:** pull-based monitoring toolkit — tətbiqləri müəyyən
interval-la scrape edir, time-series DB-də saxlayır; query/vizualizasiya/
alerting (Grafana ilə birləşir).

**Metrik tipləri:**
| Tip | Xassə | Nümunə |
|---|---|---|
| Counter | yalnız ARTIR; restart-da sıfırlanır | sorğu/xəta sayı |
| Gauge | anlıq dəyər; artır/azalır | CPU, memory, aktiv bağlantı |
| Histogram | distribusiya — bucket-lər; percentile/median | response time |
| Summary | dinamik quantile hesablanması | yüksək-car dinality latensiya |

**Kitabdan kod nümunəsi (/healthz counter):**
```go
package main

import (
    "fmt"
    "net/http"
    "time"
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
    healthzCounter = prometheus.NewCounter(prometheus.CounterOpts{
        Name: "healthz_calls_total",
        Help: "Total number of calls to the healthz endpoint.",
    })
)

func init() {
    prometheus.MustRegister(healthzCounter)   // qeydiyyat
}

func main() {
    http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
        healthzCounter.Inc()                   // hər çağırışda +1
        w.WriteHeader(http.StatusOK)
        fmt.Println("Monitoring endpoint invoked! Counter was incremented!")
    })

    http.Handle("/metrics", promhttp.Handler())   // metrics endpoint

    server := &http.Server{
        Addr:         ":8080",
        ReadTimeout:  10 * time.Second,
        WriteTimeout: 10 * time.Second,
    }
    fmt.Println("Server listening on port 8080...")
    if err := server.ListenAndServe(); err != nil {
        fmt.Printf("Error starting server: %s\n", err)
    }
}
```

**Çıxış (/metrics):**
```
# HELP healthz_calls_total Total number of calls to the healthz endpoint.
# TYPE healthz_calls_total counter
healthz_calls_total 3
# ... Go runtime metrikaları (memory, GC, goroutines) avtomatik
```
- Scrape config — Prometheus tərəfində hədəf + interval təyini
- Name konvensiyası: `_total` suffix counter-lar üçün

### 3. OpenTelemetry (OTel)
**Üç sütun:**
- **Tracing** — sorğunun servis-lər arası AXININI izlə; latency,
  asılılıqlar, xəta yayılımı
- **Metrics** — sağlamlığın kəmiyyət görüntüsü
- **Logs** — hadisələrin narrativi (strukturlu, query-lənən)

**Kitabdan kod nümunəsi (tracing + structured logging):**
```go
import (
    "context"
    "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/exporters/otlp/otlptrace"
    "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
    "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
    sdktrace "go.opentelemetry.io/otel/sdk/trace"
    "go.opentelemetry.io/otel/trace"
    "go.uber.org/zap"
)

// Trace exporter (OTLP gRPC):
func initTraceExporter(ctx context.Context) *otlptrace.Exporter {
    traceExporter, err := otlptracegrpc.New(
        ctx,
        otlptracegrpc.WithEndpoint("http://localhost:4317"),
    )
    if err != nil {
        log.Fatalf("failed to create trace exporter: %v", err)
    }
    return traceExporter
}

// Tracer provider:
func initTracerProvider(traceExporter *otlptrace.Exporter) *sdktrace.TracerProvider {
    exp, err := stdouttrace.New(stdouttrace.WithPrettyPrint())  // stdout demo
    if err != nil {
        log.Println("failed to initialize stdouttrace exporter:", err)
    }
    bsp := sdktrace.NewBatchSpanProcessor(exp)
    tp := sdktrace.NewTracerProvider(
        sdktrace.WithBatcher(traceExporter),
        sdktrace.WithSpanProcessor(bsp),
    )
    return tp
}

// Handler — strukturlu log + span:
func handler(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    span := trace.SpanFromContext(ctx)
    defer span.End()

    logger := zap.NewExample().Sugar()
    logger.Infow("Received request",
        "service", "exercise22.02",
        "httpMethod", r.Method,
        "httpURL", r.URL.String(),
        "remoteAddr", r.RemoteAddr,
    )
    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "Monitoring endpoint invoked!")
}

func main() {
    ctx := context.Background()
    traceExporter := initTraceExporter(ctx)
    defer traceExporter.Shutdown(context.Background())

    tp := initTracerProvider(traceExporter)
    otel.SetTracerProvider(tp)         // qlobal provider

    logger := initLogger()
    defer logger.Sync()

    // Handler-ı OTel instrumentation ilə BÜK:
    httpHandler := otelhttp.NewHandler(http.HandlerFunc(handler), "HTTPServer")
    http.Handle("/", httpHandler)

    server := &http.Server{Addr: ":8080",
        ReadTimeout: 10 * time.Second, WriteTimeout: 10 * time.Second}
    if err := server.ListenAndServe(); err != nil {
        fmt.Printf("Error starting server: %s\n", err)
    }
}
```

**Strukturlu log çıxışı:**
```json
{"level":"info","msg":"Received request","service":"exercise22.02",
 "httpMethod":"GET","httpURL":"/healthz","remoteAddr":"[::1]:51082"}
```
- JSON log-lar sonra sorğu-lənə bilər: "service=X olan bütün requestlər"
- otelhttp.NewHandler — hər sorğu üçün span avtomatik yaradır
- zap logger — production-grade strukturlu logging

### 4. Containerization best practices
- **Go modules** — asılılıq idarəetməsi konteynerdə də
- **Yüngül image** — multi-stage build + Alpine/scrach baza
- **Dockerfile optimizasiyası** — layer caching: az-dəyişəndən çox-dəyişənə
  sıra (go.mod əvvəl → kod sonra)
- **Təhlükəsizlik** — minimal etibarlı baza, Trivy skan, non-root user,
  sirrlər env/secrets kimi (hardcode YOX), Chainguard images
- **Health checks + logging** — orkestrator üçün /healthz; strukturlu log

**Kitabdan Dockerfile:**
```dockerfile
FROM golang:latest AS builder
ENV CGO_ENABLED=0                    # statik binary
WORKDIR /app
COPY go.mod go.sum ./                 # ƏVVƏL mod faylları (cache!)
COPY main.go ./
RUN go mod download
RUN go build -o monitored_app .

FROM scratch                          # BOŞ minimal baza!
COPY --from=builder /app/monitored_app /
EXPOSE 8080
CMD ["./monitored_app"]
```
```bash
docker build -t monitored-app .
docker run -p 8080:8080 monitored-app
# http://localhost:8080/healthz → işləyir
```
- Multi-stage: builder (SDK) → scratch (yalnız binary) — ~10MB-a yaxın image
- CGO_ENABLED=0 — statik link; scratch-də işləmək üçün MƏCBURİ

### 5. Kubernetes (K8s)
**Nədir:** container orchestration standartı — deploy, scale, management
avtomatizasiyası; unified API + control plane.

**Hazırlıq addımları:**
1. **Containerize** (əvvəlki bölmə)
2. **Deploy** — registry-ə push (Docker Hub/GCR/ECR) + deployment
   manifestləri
3. **K8s resursları** — Deployments, Services, ConfigMaps, Secrets (YAML
   manifestlər — istənilən vəziyyət təsviri)
4. **Lifecycle** — health checks, readiness probes, graceful shutdown,
   logging/metrics — K8s-ə uyğun dizayn
5. **Service discovery + load balancing** — Services daxili/kənar trafik
6. **Monitoring/logging** — Prometheus, Grafana, Fluentd, OTel —
   strukturlu emit

## Əsas terminlər
- Monitoring vs observability
- Pull-based scrape + time-series DB
- Prometheus client_golang — NewCounter/MustRegister/Inc
- promhttp.Handler() — /metrics endpoint
- Counter/Gauge/Histogram/Summary
- OpenTelemetry (OTel) — tracing/metrics/logs üç sütunu
- otlptracegrpc exporter / TracerProvider / BatchSpanProcessor
- otelhttp.NewHandler — HTTP instrumentation sargısı
- zap strukturlu logging (Infow key-value)
- Containerization — portativlik/konsistensiya/scale
- Multi-stage Dockerfile (builder → scratch)
- CGO_ENABLED=0 — statik binary şərti
- Layer caching — COPY sırası
- Trivy / Chainguard / non-root — təhlükəsizlik
- Kubernetes — Deployment/Service/ConfigMap/Secret; probes; registry

## Praktik nəticə
Cloud-a çıxmazdan əvvəl: (1) Prometheus metrikləri əlavə et — counter
(sorğu sayı), gauge (resurs), histogram (latensiya); /metrics-i
promhttp ilə expose et; (2) OTel ilə tracing (otelhttp sargısı +
TracerProvider) və zap ilə strukturlu JSON loglar — log-lar query-lənən
olacaq; (3) multi-stage Dockerfile: builder-də CGO_ENABLED=0 build →
scratch-ə yalnız binary; go.mod-u əvvəl COPY et (cache); (4) K8s üçün:
/healthz + graceful shutdown + readiness probe; manifestlərdə istənilən
vəziyyəti təsvir et. Metrik + trace + log üçlüyü production-da
problemin yerini dərhal göstərir.

## Mənbə
Pages: 646-680 (PDF 646-680)
