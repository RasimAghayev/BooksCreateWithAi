# Chapters 9-17 — YAML, Version İnjeksiyası, Embed, Şablonlar, HTTP, Prometheus, Release (səh. 80-125)

## Bu fəsillər nədən bəhs edir?

**Ch9 (YAML):** 12-factor env dəyişənləri vs config faylı (uzun/nested
məhdudiyyətlər); gopkg.in/yaml.v2 — struct tag-lər (`yaml:"name,omitempty"`),
Unmarshal/Marshal; mergo ilə obyekt birləşdirmə (base + overrides → bir Spec;
OpenFaaS stack.yml real istifadəsi).

**Ch10 (Version injeksiyası):** `var GitCommit string` + `go build -ldflags
"-X main.GitCommit=$GIT_COMMIT"`; git rev-list ilə capture; Dockerfile-da
COPY .git; docker version nümunəsi; -version flag.

**Ch11 (embed, 1.16+):** `//go:embed schema.sql` → string; []byte (logo);
embed.FS (tələb üzəri yükləmə — böyük fayllar); çoxlu qovluq + http.FS ilə
statik frontend server.

**Ch12 (Templates):** text/template — `{{.Name}}`, `{{if .Attended}}`,
`{{- else}}` whitespace-trim; template.Must; custom metodlar (Updates.
YearBreak — tipə metod bağlamaqla şablonda funksiya); arkade BinaryTemplate
(OS/Arch-a görə fayl adı; statik yoxlanma MÜMKÜNSÜZ → hər yolu unit-test).

**Ch13 (HTTP server):** http.Server{ReadTimeout/WriteTimeout/MaxHeaderBytes};
sorğu oxuma (RawQuery/Path/Method/Host/Header/Body); gorilla/mux — regex
path parametrlər (`/customer/{id:[-a-zA-Z_0-9.]+}`, mux.Vars), Methods()
zənciri; middleware (makeAuth Bearer — next-i qaytaran funksiya; zəncir
mümkün).

**Ch14 (Prometheus):** promhttp.Handler() /metrics; scrape_interval/target;
default metrikalar (go_goroutines, memstats); rate() sorğuları; RED —
InstrumentHandlerCounter/Duration (CounterVec/HistogramVec {code,method}
label-ləri); custom Gauge (MustRegister, Inc/defer Dec — inflight); Collector
interfeysi (Describe/Collect; dinamik label dəsti: Reset() + doldur —
Kubernetes replicas nümunəsi).

**Ch15-16 (Release):** Makefile dist (5 platform, LDFLAGS, -mod vendor,
-installsuffix cgo); GitHub Actions (tag push → build → upload-assets);
Docker multi-stage + multi-arch (BUILDPLATFORM/TARGETOS/ARCH ARG-lar,
distroless/static:nonroot, buildx --platform --push, QEMU); GHCR + secrets.

**Ch17:** İrəli yol — OpenFaaS/arkade/k3sup/derek/inlets-operator layihələrinə
töhfə.

## Əsas fikirlər

### 1. YAML Konfiqurasiyası (Ch9)
```go
type Spec struct {
    Name    string            `yaml:"name"`
    Image   string            `yaml:"image"`
    Environment map[string]string `yaml:"environment,omitempty"`
}
bytesOut, _ := os.ReadFile("config.yaml")
yaml.Unmarshal(bytesOut, &spec)          // JSON ilə eyni pattern
yaml.Marshal(spec)                        // yazmaq da simmetrik
```
- **12-factor vs fayl:** env = konteyner-dost; amma uzun/nested qiymətlər
  çətin → stack.yml kimi YAML config
- **mergo birləşdirmə:** base + production overrides →
  `mergo.Merge(&merged, override, mergo.WithOverride)` — standart fayl +
  istifadəçi overrides; strukturlarla işləyir (YAML-spesifik DEYİL)

### 2. Build-Zamanı Version (Ch10)
```go
// GitCommit is set at build-time
var GitCommit string
```
```bash
export GIT_COMMIT=$(git rev-list -1 HEAD)
go build -ldflags "-X main.GitCommit=$GIT_COMMIT"
```
- **Niyə:** bug reportlarda version; köhnəlmiş versiya aşkarı; docker version
  çıxışı nümunəsi
- Dockerfile: COPY .git → RUN git rev-list → build; amma multi-stage + nonroot
  şərtləri (Ch16-da tam)

### 3. embed — Məlumat Daxilə Qoşma (Ch11)
```go
import _ "embed"
//go:embed schema.sql
var tableCreate string        // build zamanı daxil olur
//go:embed logo.png
var logo []byte
import "embed"
//go:embed image/* template/*
var frontend embed.FS          // TALEBƏ GÖRE yüklənir — böyük fayllar
http.Handle("/public/", http.StripPrefix("/public/",
    http.FileServer(http.FS(frontend))))   // React frontend + Go API
```

### 4. text/template (Ch12)
```go
const letter = `Dear {{.Name}},
{{if .Attended}}...{{- else}}...{{- end}}`
t := template.Must(template.New("letter").Parse(letter))
t.Execute(&buffer, person)
```
- `{{- ... }}` — whitespace trim; Must — parse xətasında panic
- **Custom funksiyalar = tipə metod:** `{{ if ($.Updates.YearBreak $in) }}`
  → func (u Updates) YearBreak(i int) bool — ilkin elementi keçərək il
  sərhədlərini tapır (qruplaşmış siyahı)
- **arkade BinaryTemplate:** OS/arch-a görə düzgün binari adı (ming→exe,
  darwin, armhf, arm64); **statik validasiya YOXDUR → hər kod yolunu
  unit-testlə ört** (coverage ilə yoxla)

### 5. HTTP Server + Sorğu Anatomiyası (Ch13)
```go
s := &http.Server{
    Addr: ":8080",
    ReadTimeout: 10 * time.Second, WriteTimeout: 10 * time.Second,
    MaxHeaderBytes: 1 << 20,
}
```
- **Sorğu mənbələri:** Body / Header / RawQuery / URL.Path / Host
- **gorilla/mux:**
```go
r.HandleFunc("/customer/{id:[-a-zA-Z_0-9.]+}", handler)   // regex parametr
vars := mux.Vars(r); id := vars["id"]
r.HandleFunc("/system/functions", list).Methods(http.MethodGet)  // verb ayrımı
```
- **Middleware pattern:**
```go
func makeAuth(token string, next http.HandlerFunc) http.HandlerFunc {
    return func(w, r) {
        if b := r.Header.Get("Authorization"); len(b) == 0 || b != "Bearer: "+token {
            http.Error(w, "Not authorized", http.StatusUnauthorized)
            return          // next ÇAĞIRILMIR
        }
        next(w, r)
    }
}
```
  Zəncirlənə bilər (auth → logging → ...)

### 6. Prometheus RED Metrikaları (Ch14)
```go
type HttpMetrics struct {
    RequestsTotal          *prometheus.CounterVec
    RequestDurationHistogram *prometheus.HistogramVec
}
// promauto.NewCounterVec(..., []string{"code", "method"})
// promauto.NewHistogramVec(..., Buckets: prometheus.DefBuckets, ...)
InstrumentHandler := promhttp.InstrumentHandlerCounter(counter,
    promhttp.InstrumentHandlerDuration(hist, next))
```
- **RED:** Rate (counter), Errors (code label), Duration (histogram)
- **Default-lar:** go_goroutines (leak izləmə!), go_memstats_*, go_info
- **Custom Gauge:**
```go
g := prometheus.NewGauge(...)
prometheus.MustRegister(g)
hashesInflight.Inc(); defer hashesInflight.Dec()   // inflight işlər
```
- **Collector interfeysi (dinamik label dəsti):**
```go
func (e *Exporter) Describe(ch chan<- *prometheus.Desc)
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
    e.Gauge.Reset()                          // hər scrape-də yenidən qur
    for name, replicas := range e.api.QueryFunctions() {
        e.Gauge.WithLabelValues(name).Set(float64(replicas))
    }
    e.Gauge.Collect(ch)
}
```
  Silinən element label-dən DÜŞÜR (dynamic data — K8s replicas nümunəsi);
  Prometheus özü çağırır (async)

### 7. Release Pipeline (Ch15-16)
**Makefile dist:**
```makefile
LDFLAGS := "-s -w -X main.Version=$(Version) -X main.GitCommit=$(GitCommit)"
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -mod=vendor -ldflags $(LDFLAGS) -o bin/app-arm64
```
**GitHub Actions (tag → publish):**
```yaml
on: {push: {tags: ['*']}}
- run: make all
- uses: alexellis/upload-assets@0.2.2   # → GitHub Releases
```
**Multi-arch Docker:**
```dockerfile
FROM --platform=${BUILDPLATFORM} golang:1.15 as builder
ARG TARGETOS; ARG TARGETARCH
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build ...
FROM gcr.io/distroless/static:nonroot     # SDK YOX, non-root!
COPY --from=builder /usr/bin/release-it /
USER nonroot:nonroot
```
```bash
docker buildx build --platform linux/amd64 --platform linux/arm64 \
    --build-arg Version=0.1.0 --push -t user/app:0.1.0 .
```
- CI-də QEMU + buildx setup action-ları; GHCR (ghcr.io) secrets:
  DOCKER_USERNAME/PASSWORD (VƏ GITHUB_TOKEN)

## Əsas terminlər
- 12-Factor App — env-əsaslı konfiqurasiya fəlsəfəsi
- mergo — struct birləşdirici (WithOverride)
- -ldflags -X — link zamanı dəyişən injeksiyası
- embed.FS — build-ə qoşulmuş virtual fayl sistemi
- template.Must — parse xətasında panic edən konstruktor
- Middleware — handler-ı bürüyen funksiya (zəncir)
- Regex path parametr — mux-da {name:[pattern]}
- RED metrikaları — Rate/Errors/Duration
- CounterVec/GaugeVec/HistogramVec — label-li metrik tipləri
- Collector — dinamik metrik ixrac interfeysi (Describe/Collect)
- scrape_interval/target — Prometheus-in yığma cadence/uç nöqtəsi
- distroless — SDK-sız, minimal baza image
- buildx/QEMU — çox-arxitekturalı build
- GHCR — GitHub Container Registry

## Praktik nəticə

1. **YAML strukturları JSON pattern-i ilə eynidir:** tag + Unmarshal; nested
   üçün mergo (base + override siyasəti).
2. **Hər binaridə version:** -ldflags -X; -version flag; Makefile-da
   git describe --tags.
3. **Statik resurslar embed ilə:** SQL şemaları (startda cədvəl yarat!),
   frontend (embed.FS + http.FileServer).
4. **Şablonlarda custom məntiq = tip metodu:** funksiya registrasiyası YOX;
   statik yoxlanmadığından testlə ört.
5. **HTTP produksiya üçün:** mux (regex + Methods) + middleware zənciri +
   Server timeouts.
6. **RED-ə başla:** promhttp.InstrumentHandlerCounter/Duration — 3 sətirdə;
   custom Gauge inflight üçün; dinamik label-lar → Collector.
7. **Release tam avtomatik:** tag push → CI → binarilər + multi-arch image
   (distroless + nonroot mütləq).

## Mənbə
Pages: 80-125 (PDF 81-126)
