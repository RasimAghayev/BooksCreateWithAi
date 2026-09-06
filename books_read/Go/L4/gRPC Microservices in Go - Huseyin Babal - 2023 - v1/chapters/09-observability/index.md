# Chapter 9 — Observability (Müşahidəlilik)

## Bu chapter nədən bəhs edir?
Observability üçlüyünə (traces, metrics, logs), SLA/SLO/SLI və percentile anlayışlarına, OpenTelemetry instrumentasiyasına, K8s-də Jaeger+Prometheus qurulumuna, logrus ilə trace-inyeksiyalı logla-maya və Fluent Bit+Elasticsearch+Kibana log stackinə.

## Əsas fikirlər

### 1. Observability nədir
Sual cavablandırma qabiliyyəti: "Product servisi niyə dünən çökdü?", "Niyə hər cümə günortası performans düşür?", "Payment indi necə işləyir?" — **traces + metrics + logs** üçlüyü ilə.

### 2. Traces və spans
**Trace:** unikal sorğunun servislər arası səyahəti; **Span:** tək əməliyyat (DB save, RPC çağırışı).

**PlaceOrder trace:** OrderService:saveToDB (span 1) → paymentClient:Create (span 2) → paymentService:charge (span 3).

**Span taqları:** `rpc.service`, `rpc.method`, `span.kind` (client/server), `rpc.system` (grpc).

### 3. Metrics + SLA/SLO/SLI
- **SLA** — müştəri-sağlayıcı müqaviləsi (latency, throughput, uptime)
- **SLI** — xüsusi göstərici ("throughput 1000/ms")
- **SLO** — komanda hədəfi ("99.999% uptime")

**Percentile (faiz nöqtəsi) — average-in tələsi:**
1000×1ms + 10×4000ms = 5000ms / 1010 = **5ms orta** — amma 10 sorğu 4 SANİYƏ gözləyir!
→ **p95 = "95% cavab bu rəqəmdən sürətlidir"** — sırala, 80% nöqtəni tap (nümunədə 11ms).
→ Müştərilərə ORTA YOX, PERCENTILE vəd et.

### 4. OpenTelemetry instrumentasiyası
**Interceptor ilə (retry/circuit-breaker pattern-i kimi):**
```go
var opts []grpc.DialOption
opts = append(opts,
    grpc.WithTransportCredentials(insecure.NewCredentials()),
    grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),   // OTel telemetriyası
)
conn, err := grpc.Dial(paymentServiceUrl, opts...)
```
**İnstrumentasiya nöqtələri:** gRPC client, gRPC server, GORM DB (otelgorm extension).

### 5. Metrics backend — Jaeger + Prometheus
**Arxitektura:** Order/Payment → OpenTelemetry Collector (port 14278) → Prometheus (8889, time-series) + Jaeger UI.

**Helm qurulumu:**
```bash
helm repo add huseyinbabal https://huseyinbabal.github.io/charts
helm install my-jaeger huseyinbabal/jaeger -n jaeger --create-namespace
```
**Komponentlər:** Jaeger All in One + OTel Collector + Prometheus.

**Jaeger SPM (Service Performance Monitoring):** servislərin p95 latency, request rate, error rate — Prometheus query-lərindən aqreqasiya.

**Nümunə trace analizi:** create order 4.2ms → Payment Create 3.36ms → içində DB 2.88ms (GORM instrumentasiyası SQL statement-i belə göstərir).

### 6. Application logging — trace-inyeksiyalı loglar
**logrus + custom formatter:**
```go
type serviceLogger struct {
    formatter log.JSONFormatter
}

func (l serviceLogger) Format(entry *log.Entry) ([]byte, error) {
    span := trace.SpanFromContext(entry.Context)      // context-dən span
    entry.Data["trace_id"] = span.SpanContext().TraceID().String()
    entry.Data["span_id"] = span.SpanContext().SpanID().String()
    return l.formatter.Format(entry)
}

func init() {
    log.SetFormatter(serviceLogger{
        formatter: log.JSONFormatter{FieldMap: log.FieldMap{"msg": "message"}},
    })
    log.SetOutput(os.Stdout)
    log.SetLevel(log.InfoLevel)
}
```
**İstifadə:**
```go
log.WithContext(ctx).Info("Creating order...")
// Çıxış: {"trace_id":"4a55...", "span_id":"1ab6...", "level":"info", "message":"Creating order..."}
```
**Kibana-da TraceID filtri ilə** bütün axın boyu logları birləşdir.
**JSON üstünlüyü:** backend sahə-map edir — plaintext parser lazım deyil.

### 7. K8s log arxitekturası
| Model | Mexanizm |
|---|---|
| **Node-level** | stdout → fayl + log rotation (böyüyən fayl problemi) |
| **Cluster-level** | daemonset agent node-dakı /var/log/containers fayllarını stream edir → backend |

### 8. Log stack — Fluent Bit + Elasticsearch + Kibana
**Fluent Bit (collector, DaemonSet):**
```bash
helm repo add fluent https://fluent.github.io/helm-charts
helm repo update
helm install fluent-bit fluent/fluent-bit
```
**Elasticsearch (ECK operatoru ilə):**
```bash
kubectl create -f https://download.elastic.co/downloads/eck/2.5.0/crds.yaml
kubectl apply -f https://download.elastic.co/downloads/eck/2.5.0/operator.yaml

cat <<EOF | kubectl apply -f -
apiVersion: elasticsearch.k8s.elastic.co/v1
kind: Elasticsearch
metadata:
  name: quickstart
spec:
  version: 8.5.2
  nodeSets:
  - name: default
    count: 1
    config:
      node.store.allow_mmap: false
EOF

# Parol:
PASSWORD=$(kubectl get secret quickstart-es-elastic-user -o go-template='{{.data.elastic | base64decode}}')
```
**Fluent Bit output konfiqurasiyası (fluent.yaml):**
```yaml
config:
  outputs: |
    [OUTPUT]
        Name es
        Match kube.*
        Host quickstart-es-http
        HTTP_User elastic
        HTTP_Password $PASSWORD
        tls On
        tls.verify Off
        Logstash_Format On
        Retry_Limit False
```
```bash
helm upgrade --install fluent-bit fluent/fluent-bit -f fluent.yaml
```
**Kibana:**
```bash
cat <<EOF | kubectl apply -f -
apiVersion: kibana.k8s.elastic.co/v1
kind: Kibana
metadata:
  name: quickstart
spec:
  version: 8.5.2
  count: 1
  elasticsearchRef:
    name: quickstart
EOF
# port-forward 5601 → http://localhost:5601; index pattern: log*
```
**Axın:** konteyner stdout → Fluent Bit (node daemonset) → Elasticsearch → Kibana dashboard → TraceID filtri.

## Əsas terminlər
- Observability üçlüyü — traces/metrics/logs
- Trace / Span — sorğu səyahəti / əməliyyat vahidi
- SLA / SLI / SLO
- Percentile (p95) — sıralanmış faiz nöqtəsi
- OpenTelemetry — SDK/API kolleksiyası; otelgrpc interceptor
- OTel Collector (14278) → Prometheus (8889)
- Jaeger SPM — servis performans monitorinqi
- trace.SpanFromContext — context-dən span çıxarışı
- logrus JSONFormatter + FieldMap
- DaemonSet — node başına 1 pod
- Fluent Bit — log collector/shipper
- Elasticsearch (ECK CRD) + Kibana dashboard
- Logstash_Format — tarixli index konvensiyası

## Praktik nəticə
1. Performans vədlərində AVERAGE YOX — p95/p99 percentile işlət.
2. OTel instrumentasiyası = interceptor əlavə etmək qədər asan; DB (GORM) də daxil et.
3. Hər loga trace_id + span_id inyeksiya et (custom formatter) — Kibana-da tam axın axtarışı mümkün olur.
4. Cluster-level logging: Fluent Bit DaemonSet + ES + Kibana — stdout JSON formatı parser ehtiyacını aradan qaldırır.
5. Trace-metric-log korrelyasiyası: metrics anomaliyası → trace bax → həmin TraceID-li loglar — kök səbəb 3 addımda.

## Mənbə
Pages: 153-174 (PDF səh. 180-195)
