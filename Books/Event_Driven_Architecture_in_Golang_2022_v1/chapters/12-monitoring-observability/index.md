# Chapter 12 — Monitoring and Observability (səh. 344-382)

## Bu chapter nədən bəhs edir?

Monitoring vs observability fərqi, correlation/causation identifikatorları,
OpenTelemetry distributed tracing (spans, events), Prometheus metrikaları
(counters, histograms), Jaeger/Grafana ilə vizuallaşdırma.

## Əsas fikirlər

### 1. Monitoring vs Observability
- **Monitoring:** "Nə oldu?" + "Niyə oldu?" — loglar + metrikalar
- **Observability:**monitoring-dən uzağa — sistemlərin daxili vəziyyətini
  **xaricdən çıxarış** etmək; korrelyasiya üçün daha çox siqnal (trace-lər)

### 2. İdentifikatorlar (tracing 3-lüyü)
| ID | Məqsəd |
|---|---|
| **Request ID** | Bir HTTP sorğusunu izləmə |
| **Correlation ID** | Bir iş axını boyu bütün komponentləri birləşdirmək |
| **Causation ID** | Bu mesajı YARADAN səbəbi qeyd etmək |

EDA-da saga bir sorğu üzrə 5 servisi gəzir — correlation ID olmadan izi
itirmək mümkün deyil.

### 3. OpenTelemetry (OTel)
Standart izləmə çərçivəsi — vendor-neytral.

**Quruluş:**
```go
func initOpenTelemetry(ctx context.Context) error {
    conn, err := grpc.DialContext(ctx, endpoint,
        grpc.WithTransportCredentials(insecure.NewCredentials()))
    // ... exporter → trace provider
    otel.SetTextMapPropagator(
        propagation.NewCompositeTextMapPropagator(
            propagation.TraceContext{},
            propagation.Baggage{},
        ))
}
```
- gRPC üzərindən OTLP collector-a yollanır
- **Interceptor-lar:** gRPC client/server çağırışları avtomatik span açır —
  mövcud middleware sıralamasından asılılıq YOX (təmiz yerləşmə)

**Span-lər:**
- Sorğu modullar arası parçalanır → span zənciri
- **Events** (annotasiyalar): hadisə handler-ləri "before/after" işarələri
  qoyur — uğursuzluq səbəbi span daxilində görünür

### 4. Prometheus metrikaları
Endpoint: `promhttp` handler `/metrics`-də.
```go
receivedMessagesLatency = promauto.NewHistogramVec(prometheus.HistogramOpts{
    Name: "received_messages_latency_seconds",
    Buckets: []float64{0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5},
}, []string{"message"})
// labels: message adı → hər hadisə ayrıca sıra
```
- **Counter:** `customers_registered_total` kimi artan sayğaclar
- **Histogram:** latency bölmələri (buckets)
- InstrumentedApp wrapper: App interfeysinə sayğac əlavə et — domain kodu
  təmiz qalır

### 5. Vizuallaşdırma
| Alət | Göstərir |
|---|---|
| **Jaeger** | Trace-lər — servis seç → Find Traces; span detalları, böyük
  dairələr = yavaş yer |
| **Prometheus UI** | Metrik axtarışı (`cosec_received_messages_count`,
  `go_gc_duration_seconds` — GC hər servis üçün) |
| **Grafana** | Dashboard-lar (hazır Go runtime + OTel Collector dashboard-ları) |

Create order prosesi Jaeger-də tam zəncir kimi görünür: baskets → orders saga →
payments → depot → notifications.

## Əsas terminlər

- Observability vs Monitoring
- Correlation / Causation ID
- OpenTelemetry / OTLP / Collector
- Span / Span Event
- Trace Propagation (TraceContext, Baggage)
- Interceptor (avtomatik instrumentasiya)
- Histogram / Counter / Label
- Jaeger / Grafana

## Praktik nəticə

- EDA-da correlation ID mütləqdir — bütün mesajlara header kimi ötür
- gRPC interceptor + otel propagator: kod dəyişmədən tracing
- Hadisə handler-lərinə span event qoy — səbəb izahı trace-də otursun
- Metrikaları wrapper/interceptor qatında saxla; domain-i təmiz burax
- Əvvəl hazır dashboard-lardan başla (Go runtime, OTel collector)

## Mənbə

Pages: 344-382 (Chapter 12, Event-Driven Architecture in Golang)
