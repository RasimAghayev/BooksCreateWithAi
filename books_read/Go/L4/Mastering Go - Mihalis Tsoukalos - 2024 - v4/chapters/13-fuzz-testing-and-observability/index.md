# Chapter 13 — Fuzz Testing and Observability (Fuzz Test və Müşahidəlilik)

## Bu chapter nədən bəhs edir?

Fuzz testing (testing.F, seed corpus, testdata, mənfi ədəd və UTF-8 bug nümunələri),
observability konsepti (logs, metrics, traces, monitoring, alerting, distributed tracing),
runtime/metrics paketi, funksiya icra vaxtının ölçülməsi, expvar, cpuid, Prometheus
client (counter/gauge/histogram/summary, /metrics, Docker image), prometheus.yml,
docker-compose monitoring şəbəkəsi və Grafana vizuallaşdırması.

## Əsas fikirlər

### 1. Fuzz Testing Nədir
**Nədir:** Proqrama RANDOM/qeyri-valid/gözlənilməz input generasiya edərək test etmək.
`testing.F` tipi (Test/B kimi); funksiyalar `Fuzz` ilə başlayır.

**Nəyə lazımdır:** input parse edən kodun (buffer overflow, SQL injection) təhlükəsizlik
yoxlaması; unit testin ƏVƏZİ DEYİL — tamamlayıcısı.

**Güclü tərəfləri:** invalid inputa davamlılıq; ciddi təhlükəsizlik bug-ları; attacker-lər
də fuzzing işlədir — hazırlıqlı ol.

**Mexanizm:** uğursuz input `testdata/fuzz/FuzzX/` qovluğuna YAZILIR → bundan sonra
adi `go test` də həmin data ilə FAIL edir (bug düzələnə və ya testdata silinənə qədər).
Corpus = seed corpus (f.Add + testdata/fuzz) + generated corpus (machine-generated).

### 2. Sadə Fuzz Nümunəsi — Mənfi Ədəd Bug
**Kitabdan kod nümunəsi:**
```go
func AddInt(x, y int) int {
    for i := 0; i < x; i++ {   // x<0 olanda loop İŞLƏMİR!
        y = y + 1
    }
    return y
}

func FuzzAddInt(f *testing.F) {
    testCases := []struct{ x, y int }{
        {0, 1}, {0, 100},
    }
    for _, tc := range testCases {
        f.Add(tc.x, tc.y)          // SEED CORPUS — Add imza ilə uyğun
    }
    f.Fuzz(func(t *testing.T, x, y int) {   // parametrlər Add ilə EYNİ
        result := AddInt(x, y)
        if result != x+y {
            t.Errorf("X: %d, Y: %d, Result %d, want %d", x, y, result, x+y)
        }
    })
}
```
```bash
go test -fuzz=FuzzAddInt *.go
# FAIL: X: -63, Y: 32, Result 32, want -31   ← REGULAR TEST GÖRMƏDİ, FUZZ GÖRDÜ!
# Failing input written to testdata/fuzz/FuzzAddInt/b403d5353f8afe03
# To re-run: go test -run=FuzzAddInt/b403d5353f8afe03
```
Müntəzəm testlər {1,2,3},{100,10,110} kimi YALNIZ müsbət ədədlərlə keçirdi — fuzz
mənfi ədədi tapdı.

### 3. İnkişaf Etmiş Nümunə — UTF-8 Reverse Bug
**Problem:** bayt üzrə reverse multi-byte Unicode simvolları pozur.
**Fuzz yoxlaması (2 xassə):**
```go
func FuzzR1(f *testing.F) {
    for _, tc := range []string{"Hello, world", " ", "!12345"} {
        f.Add(tc)
    }
    f.Fuzz(func(t *testing.T, orig string) {
        rev := R1(orig)
        doubleRev := R1(string(rev))
        if orig != string(doubleRev) {          // XASSƏ 1: 2x reverse = original
            t.Errorf("Before: %q, after: %q", orig, doubleRev)
        }
        if utf8.ValidString(orig) && !utf8.ValidString(string(rev)) {  // XASSƏ 2: validlik
            t.Errorf("Reverse: invalid UTF-8 string %q", rev)
        }
    })
}
```
Fuzz tapdı: `string("Ԕ")` (2-baytlı simvol) pozulur → `invalid UTF-8 string "\x94\xd4"`.

**Düzəliş (kitabdan):**
```go
func R2(s string) (string, error) {
    if !utf8.ValidString(s) {                    // 1) VALIDASİYA ƏVVƏL
        return s, errors.New("Invalid UTF-8")
    }
    r := []rune(s)                               // 2) RUNE üzrə (bayt YOX)
    for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
        r[i], r[j] = r[j], r[i]
    }
    return string(r), nil
}
```
```bash
go test -fuzz=FuzzR1 *.go -fuzztime 10s    # 4.5M exec, PASS — bug DÜZƏLDİ
```
`-fuzztime` — vaxt limiti; bitənə və ya ilk xətaya qədər işləyir.

### 4. Observability — Komponentlər
**Nədir:** XARİCİ siqnallara əsaslanan sistemin DAXİLİ vəziyyətinin anlaşılması.
**Birinci qayda:** NƏ axtardığını BIL — yoxsa yanlış data toplayıb vacib metrikaları
buraxırsan.

| Komponent | Vəzifə |
|---|---|
| Logs | hadisə tarixçəsi; debug/audit |
| Metrics | kəmiyyət ölçümləri (response time, error rate) |
| Traces | sorğunun axını; latency+asılılıq |
| Monitoring | fasiləsiz izləmə + anomaliya |
| Alerting | threshold aşımında bildiriş |
| Distributed Tracing | mikroservis sorğu izlənməsi |

Metrikları saxlamadan (visualization) oxumaq EFEKTİVSİZDIR.

### 5. runtime/metrics — Go Runtime Metrikaları
**Kitabdan kod nümunəsi:**
```go
const nGo = "/sched/goroutines:goroutines"     // PATH formatı

getMetric := make([]metrics.Sample, 1)        // Sample{Name, Value}
getMetric[0].Name = nGo

metrics.Read(getMetric)                         // OXU
if getMetric[0].Value.Kind() == metrics.KindBad {
    // metric artıq dəstəklənmir
}
mVal := getMetric[0].Value.Uint64()             // Uint64-ə çevir
```
Bütün siyahı: `metrics.All()`. Nümunə: 3 goroutine + main → 4; wg.Wait-dən sonra → 1.

### 6. Funksiya İcra Vaxtının Ölçülməsi
```go
now := time.Now()
myFunction()
elapsed := time.Since(now)
logger.Info("Observability", slog.Int64("time_taken", int64(elapsed)))
// time.Duration = int64 (nanosaniyə) — cast mütləqdir
```
runtime/metrics bunu dəstəkləmir — öz əlinlə time.Now/Since sarğısı + slog ilə
strukturlaşdırılmış qeyd.

### 7. expvar — Dəyişənlərin HTTP-dən İfşası
```go
intVar := expvar.NewInt("intVar")
intVar.Set(1234)
intVar.Add(10)                                   // atomik artırma

expvar.Publish("customFunction", expvar.Func(func() interface{} {
    return "Hi from Mastering Go!"
}))

http.Handle("/debug/expvars", expvar.Handler())  // standart YERİ: /debug/vars
```
JSON formatında ifşa olunur: `{"intVar": 1244, "customFunction": "..."}`.

### 8. cpuid — CPU Xüsusiyyətləri
```go
import . "github.com/klauspost/cpuid/v2"   // DOT import — CPU birbaşa

CPU.BrandName; CPU.PhysicalCores; CPU.LogicalCores; CPU.ThreadsPerCore
CPU.FeatureSet(); CPU.CacheLine; CPU.Cache.L1D/L1I/L2/L3; CPU.Hz
```
Runtime-da CPU cache/freq/feature məlumatı — observability metrikaları üçün kontekst.

### 9. Prometheus Metrikaları — 4 Tip
| Tip | Davranış | İstifadə |
|---|---|---|
| Counter | yalnız ARTIR (və ya reset) | cəmi sorğu sayı, cəmi xəta |
| Gauge | artır/azalır | cari goroutine sayı, temp |
| Histogram | bucket-lərə sample | request duration |
| Summary | histogram + quantile (sliding window) | statistika |

Adətən counter/gauge kifayətdir.

### 10. Metrikların İfşası — client_golang
**Kitabdan kod nümunəsi:**
```go
var counter = prometheus.NewCounter(prometheus.CounterOpts{
    Namespace: "mtsouk",           // QRUPŞƏKİLİLİK üçün vacib!
    Name:      "my_counter",
    Help:      "This is my counter",
})
var gauge = prometheus.NewGauge(...)

prometheus.MustRegister(counter)    // QEYDİYYAT MÜTLƏQ (define kifayət deyil)
prometheus.MustRegister(gauge)

// Dəyərlərin yenilənməsi (fon goroutine-ində):
counter.Add(rand.Float64() * 5)      // counter
gauge.Add(rand.Float64()*15 - 5)    // gauge (mənfi də ola bilər)
histogram.Observe(rand.Float64() * 10)   // observe
summary.Observe(...)

// HTTP ifşa — promhttp hazır handler:
http.Handle("/metrics", promhttp.Handler())
http.ListenAndServe(PORT, nil)
```
```bash
curl localhost:1234/metrics --silent | grep mtsouk
# mtsouk_my_counter 19.948
# mtsouk_my_gauge 29.335
# mtsouk_my_histogram_bucket{le="5"} 4 ...
```

### 11. runtime/metrics → Prometheus Körpüsü
**Kitabdan kod nümunəsi:**
```go
var nGoroutines = prometheus.NewGauge(prometheus.GaugeOpts{
    Namespace: "packt", Name: "n_goroutines", ...})
var nMemory = prometheus.NewGauge(...)

const nGo = "/sched/goroutines:goroutines"
const nMem = "/memory/classes/heap/free:bytes"

getMetric := make([]metrics.Sample, 2)
getMetric[0].Name = nGo
getMetric[1].Name = nMem

http.Handle("/metrics", promhttp.Handler())

go func() {
    for {
        // ... goroutine yarat (metrika dəyişsin) ...
        runtime.GC()                   // memory metrikasını dəyişmək üçün GC çağır
        metrics.Read(getMetric)
        nGoroutines.Set(float64(getMetric[0].Value.Uint64()))   // SET ilə yaz
        nMemory.Set(float64(getMetric[1].Value.Uint64()))
        time.Sleep(...)
    }
}()
```
**Arxitektura:** HTTP server main goroutine-də, metric toplama ayrı goroutine-də.

### 12. Docker Image + Monitoring Stack
**Go app Dockerfile (multi-stage):** builder (golang:alpine + git + go mod) → alpine +
binary; EXPOSE 1234.

**prometheus.yml:**
```yaml
scrape_configs:
  - job_name: GoServer
    scrape_interval: 5s              # hər 5 san PULL
    static_configs:
      - targets: ['goapp:1234']     # container_name ilə hostname!
```

**docker-compose.yml (3 konteyner, monitoring şəbəkəsi):**
```yaml
goapp:       {image: goapp, ports: 1234:1234, networks: [monitoring]}
prometheus:  {image: prom/prometheus, volumes: [./prometheus/:/etc/prometheus/,
               ./prometheus_data/:/prometheus/],
              command: ['--config.file=/etc/prometheus/prometheus.yml', ...],
              ports: 9090:9090}
grafana:     {image: grafana/grafana, depends_on: [prometheus],
              environment: [GF_SECURITY_ADMIN_PASSWORD=helloThere, ...],
              ports: 3000:3000, volumes: [./grafana_data/:/var/lib/grafana/]}
networks: {monitoring: {driver: bridge}}
```
**Vacib:** volumes ilə data konteyner restartında İTMİR; şəbəkə daxilində hostnamelər
container_name-lərdir (goapp, prometheus, grafana).

### 13. Grafana Qoşulması
1. http://localhost:3000 → admin + GF_SECURITY_ADMIN_PASSWORD
2. Add data source → Prometheus → URL: `http://prometheus:9090` (ŞƏBƏKƏ hostname-i!)
→ Save & Test
3. Dashboard → Panel → data source Prometheus → Metrics drop-down → istənilən metrika
→ Save

## Əsas terminlər
- Fuzz Testing — random/qeyri-valid input generasiyası ilə test
- testing.F — fuzz test tipi; f.Add (seed) + f.Fuzz (icra)
- Seed Corpus — başlanğıc inputlar; f.Add/testdata/fuzz
- Generated Corpus — machine-generated inputlar
- testdata/fuzz — uğursuz inputun saxlanma yeri
- -fuzztime — fuzz icra müddəti
- Observability — xarici siqnallardan daxili vəziyyət anlayışı
- runtime/metrics — Go runtime metrikaları (path-formatlı)
- metrics.Sample/Read/All — metrika oxu quruluşu
- expvar — /debug/vars JSON ifşa paketi
- Counter/Gauge/Histogram/Summary — Prometheus metrik tipləri
- prometheus.MustRegister — metrik qeydiyyatı
- promhttp.Handler — /metrics HTTP handler
- Namespace — metrik qrup prefiksi
- scrape_interval — Prometheus pull tezliyi
- Grafana Data Source — vizuallaşdırma üçün backend bağlantısı

## Praktik nətidə

(1) Unit testlər "gözlənilən" halları, fuzz "gözlənilməyən"i tapır — hər ikisi birlikdə.
(2) f.Add imzası = f.Fuzz parametrləri = test edilən funksiyanın arqumentləri. (3) Uğursuz
fuzz inputu testdata-da qalır — düzəltmədən sonra da adi test FAIL edir (regression
qorunması kimi DÜZGÜNDÜR). (4) Property-based fuzz invariantları: double-reverse =
original, valid UTF-8 → valid UTF-8. (5) String manipulyasiyasında []rune + utf8.ValidString
ilə başla — bayt əməliyyatları multi-byte pozur. (6) Observability-ə NƏYİ ölçməli sualı
ilə başla. (7) runtime/metrics + Prometheus Set — runtime data-nın monitorinqə körpüsü.
(8) Prometheus PULL modeli: app /metrics ifşa edir, server gəlir götürür. (9) Metrik
adlarında Namespace — qruplaşdırma. (10) Monitoring stack-i Docker şəbəkəsində saxla —
container_name = hostname; volumes = data davamlılığı. (11) Grafana datasource URL
şəbəkə hostname-idir (localhost YOX). (12) Counter yalnız artır — azalan dəyər üçün gauge.

## Mənbə
Pages: 581-624 (PDF 612-657)
