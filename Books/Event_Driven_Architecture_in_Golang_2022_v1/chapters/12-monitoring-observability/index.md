# Chapter 12 — Monitoring and Observability (səh. 344-367)

## Bu chapter nədən bəhs edir?

Observability-nin 3 sütunu (logging, metrics, tracing): monitoring vs
observability fərqi, correlation/causation identifier-lər, OpenTelemetry ilə
distributed tracing (gRPC interceptor-lər, message middleware, span event-lər),
Prometheus counter/histogram metric-ləri, Jaeger/Grafana ilə vizuallaşdırma.

## Əsas fikirlər

### 1. Monitoring vs Observability
- **Monitoring:** toplanan datanın analizinə REAKSİYA — "nə oldu?" və "niyə
  oldu?" suallarına cavab; alətlər: logging + metrics (CPU, latency)
- **Observability:** datanı alətlərarası korrelyasiya ETMƏK mümkündür —
  məs. CPU alert → dashboard-larla timeframe tap → loglarda xəta axtarışı
  çətindir; observability bu əlaqələndirməni asanlaşdırır

### 2. Tracing identifikatorları
| ID | Rol |
|---|---|
| Request ID | tək istək |
| **Correlation ID** | bütün istəkləri TƏK mənə istəyə bağlayır — dəyişmir |
| **Causation ID** | növbəti istəyin ƏVVƏLKİ istəyə pointer-i — hər addımda yenilənir |

Trace = span-lər zənciri (icicle / upside-down flame graph dünya). **Trace-lər
log mesajlarından yaradıla BİLMƏZ** — öz tracing implementasiyanı yazmaq
yerinə OpenTelemetry istifadə et.

### 3. OpenTelemetry qurulması
```bash
OTEL_SERVICE_NAME: mallbots
OTEL_EXPORTER_OTLP_ENDPOINT: http://collector:4317   # gRPC collector
```
```go
func initOpenTelemetry(ctx context.Context) error {
    // gRPC connection to collector (env-dən)
    // trace.NewTracerProvider(batch) — batch perf üçün
    // propagation: TraceContext{} + Baggage{}
}
```

**gRPC avtomatik instrumentasiya:**
```go
grpc.DialContext(ctx, endpoint,
    grpc.WithTransportCredentials(insecure.NewCredentials()),
    grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
    // stream üçün: otelgrpc.StreamClientInterceptor()
)
```

**Message publisher/subscriber:** hazır middleware YOX — öz wrapper-lər
(amotel paketi): `amotel.OtelMessageContextExtractor()` constructor-a
əlavə olunur — hər constructor qaplandıqda span data çıxır.

### 4. Span zənginləşdirmə
```go
span := trace.SpanFromContext(ctx)
span.AddEvent("handling domain event", trace.WithAttributes(
    attribute.String("event", event.EventName()),
))
// xəta halında:
span.AddEvent("Error encountered handling domain event",
    trace.WithAttributes(errorsotel.ErrAttrs(err)...))
```
- AddEvent-lər grafikdə vaxt nisbətində düzgün yerləşir
- Alternativ: `span.RecordError(err)` — xətanı birbaşı span-a yaz

### 5. Prometheus metrics
**Endpoint:** `promhttp.Handler()` → `/metrics` route (chi mux-a):
```go
s.mux.Method("GET", "/metrics", promhttp.Handler())
```

**Counter (promauto ilə avtomatik qeydiyyat):**
```go
var receivedMessagesCounter = promauto.NewCounterVec(
    prometheus.CounterOpts{
        Name: "received_messages_count",
    },
    []string{"message", "handled"},   // label-lər
)
receivedMessagesCounter.WithLabelValues("message", "true").Inc()
// promauto OLMASA: prometheus.MustRegister(counter) əl ilə
```

**Histogram (latency + bucket-lər):**
```go
var receivedMessagesLatency = promauto.NewHistogramVec(
    prometheus.HistogramOpts{
        Name: "received_messages_latency_seconds",
        Buckets: []float64{0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5},
    },
    []string{"message", "handled"},
)
```

**Middleware zənciri (ID ilə):**
```go
am.NewMessageSubscriber(
    stream,
    amotel.OtelMessageContextExtractor(),      // tracing
    amprom.ReceivedMessagesCounter("baskets"), // metrics
)
```

### 6. Application wrapper (instrumented app)
Kodu DƏYİŞMƏDƏN instrumentasiya — decorator pattern:
```go
func NewInstrumentedApp(app App, customersRegistered prometheus.Counter) App {
    return instrumentedApp{App: app, customersRegistered: customersRegistered}
}

func (a instrumentedApp) RegisterCustomer(ctx context.Context, register RegisterCustomer) error {
    if err := a.App.RegisterCustomer(ctx, register); err != nil {
        return err
    }
    a.customersRegistered.Inc()   // uğurlu qeydiyyat sayğıcı
    return nil
}
```

### 7. Vizuallaşdırma stack
| Alət | Rol |
|---|---|
| **Jaeger** ( :14250/16686) | Trace-lər: servis seç → Find Traces; dairə ölçüsü = span sayı, hündürlük = müddət; saga trace-ləri (create order) servis əməkdaşlığını göstərir; Tags (BasketID, PaymentID) + Logs |
| **Prometheus UI** | raw metric axtarışı (`cosec_received_messages_count`, `go_gc_duration_seconds`); hər mikroservis üçün scrape job: `- job_name: basket` |
| **Grafana** (:3000) | Hazır dashboardlar; + OpenTelemetry Collector dashboard — collector-un öz yükü |

### 8. Yekun arxitektura nəticəsi
3 sütun tam: logging (əvvəldən var idi) + tracing + metrics.Əlavələrin
tətbiqə ölçülə bilən təsiri yoxdur — amma artıq hər hansı təsiri MONİTOR
etmək mümkündür.

## Əsas terminlər

- Observability vs Monitoring
- Correlation ID / Causation ID
- Distributed Tracing
- Span / Trace
- OpenTelemetry (OTEL) / Collector
- gRPC Interceptor
- promauto / CounterVec / HistogramVec / Labels / Buckets
- Instrumented Application (decorator)
- Jaeger / Grafana
- Three Pillars of Observability

## Praktik nəticə

- Öz tracing yazma — OTel standartı; collector-a batch ilə göndər
- gRPC: interceptor-lar avtomatik span yaradır; custom mesaj kodu üçün
  middleware wrapper
- Metric-lər label-lərlə bölünür (`message`, `handled`); latency üçün
  histogram + bucket-lər
- Application-ı wrapper ilə instrumentasiya et — core kodu dəymədən
- Prometheus scrape config hər servis üçün job; vizuallaşdırma Jaeger+Grafana

## Mənbə

Pages: 344-367 (Chapter 12, Event-Driven Architecture in Golang)
