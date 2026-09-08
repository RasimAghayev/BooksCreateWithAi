# Chapter 12 — Setting Up Service Alerting (səh. 251-266)

## Bu fəsil nədən bəhs edir?

Alerting-in əsasları, Four Golden Signals, Prometheus + Alertmanager ilə
uçdan-uca alert pipeline qurulması və best practices.

## Əsas fikirlər

### 1. Alerting əsasları
- İnsidentlər qaçılmazdır: resurs limitləri, network conjestion, asılılıq
  xətaları
- **Alerting = insident aşkarı + bildiriş**
- Mexanizm: şərt (query) təyini → periodik qiymətləndirmə → şərt ödənəndə
  action (email/SMS)
- Pseudocode: `active_user_count == 0`

### 2. Four Golden Signals (Google SRE)
| Signal | Nə ölçür | Nümunə |
|---|---|---|
| **Latency** | emal müddəti | API request müddəti |
| **Traffic** | yük | request/s |
| **Errors** | xəta nisbəti | failed/total |
| **Saturation** | resurs doluluğu | CPU/RAM/disk/fd |

**Praktik siqnallar:**
- API: client error rate, server error rate, latency
- Saturation: CPU%, memory, disk, **open file descriptors** (proses limiti var!)
- Digər: **service panics** (toleransiya SIFIR), failed deployments

### 3. Prometheus modeli
- Metric import yolları: **Scraping** (Prometheus çəkir — tövsiyə; servis
  `/metrics` endpoint açır) vs Pushing (Pushgateway vasitəsilə)
- `/metrics` çıxışı: `active_user_count 755` — açar-dəyər formatı
- **PromQL:** `active_user_count{service="rating-ui"}` — matcher filtrləri;
  şərtlər Boolean ifadədir: `active_user_count == 0`;
  `api_request_latency{quantile="0.5"} > 1` (median)
- **Alertmanager:** ayrı komponent — qruplaşdırma, retry, suppression

### 4. Alertmanager konfiqi
```yaml
groups:
- name: Availability alerts
  rules:
  - alert: Rating service down
    expr: service_availability{service="rating"} == 0
    for: 3m              # 3 dəq şərt aramsız ödənəndə
    labels: {severity: page}
    annotations:
      title: Rating service availability down
      description: No available instance of the rating service.
```

### 5. Servis tərəfində (tally + prometheus reporter)
```go
reporter := prometheus.NewReporter(prometheus.Options{})
scope, closer := tally.NewRootScope(tally.ScopeOptions{
    Tags: {"service": "metadata"}, CachedReporter: reporter}, 10*time.Second)
http.Handle("/metrics", reporter.HTTPHandler())
go http.ListenAndServe(":8091", nil)   // metricsPort ayrı portda!
counter := scope.Tagged(...).Counter("service_started"); counter.Inc(1)
```
- `/metrics` cavabına Go runtime metricləri də düşür (GC, goroutines)

### 6. Prometheus konfiq (prometheus.yaml)
```yaml
global: {scrape_interval: 15s, evaluation_interval: 15s}
rule_files: [alerts.rules]
scrape_configs:
- job_name: prometheus
  static_configs:
  - targets: [host.docker.internal:8091]   # metadata
    labels: {service: metadata}
  - targets: [host.docker.internal:8092]   # rating
  - targets: [host.docker.internal:8093]   # movie
```
- `host.docker.internal` — Docker-dan hostdakı servislərə çıxış
- Dinamik mühitdə: **consul_sd_configs** (Consul-dan avtomatik targetlər)

**alerts.rules:** `up{service="metadata"} == 0` — scrape uğursuz = servis down

### 7. İcra və test
1. `docker run -p 9090:9090 -v configs:/etc/prometheus prom/prometheus`
   → UI :9090; `up` query; Alerts tab — inactive
2. alertmanager.yml (email receiver, smtp.gmail.com:587)
3. `docker run -p 9093:9093 prom/alertmanager` → UI :9093
4. Rating + movie servislərini DAYANDIR → alertlər **firing** → email gəlir

### 8. Best practices
1. **Alertlər dərhal actionable olsun** — CPU spike müvəqqəti ola bilər;
   `for: 10m` kimi müddət ver, noise azalt
2. **Runbook linkləri** — hər alert üçün mitigation təlimatı
3. **Konfiqi periodik review** — alert qaydaları kod bazasında (reviewable),
   köhnəlmiş alertləri təmizlə

## Termindirmə (AZ)
- Alerting — Xəbərdarlıq Sistemi
- Four Golden Signals — Dörd Qızıl Siqnal
- Latency — Gecikmə
- Traffic — Trafik (yük)
- Saturation — Doluluq (resurs istifadəsi)
- Scraping — Metrik Çəkilməsi
- Alertmanager — Alert İdarəçisi
- Firing — Alovlanan (aktivləşən alert)

## Kviz sualları
1. Four Golden Signals hansılardır? (Latency, Traffic, Errors, Saturation)
2. `for: 3m` nə edir? (Şərt 3 dəqiqə aramsız ödənəndə alert atılır — noise
   azaldır)
3. Scraping vs Pushing fərqi? (Scrape: Prometheus /metrics-dən çəkir;
   Push: servis Pushgateway-ə göndərir)
4. Dinamik servis siyahısı üçün nə istifadə olunur? (consul_sd_configs —
   Consul registry əsaslı target kəşfi)
