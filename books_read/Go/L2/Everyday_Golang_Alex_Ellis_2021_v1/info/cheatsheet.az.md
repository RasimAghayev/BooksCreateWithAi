# Everyday Golang — Cheat Sheet (Azərbaycanca)

## 1. Layihə başlanğıcı

```bash
go mod init github.com/username/project   # GOPATH xaricində ad mütləq
go mod vendor; go build -mod vendor       # private deps / CI
go install path@v1.3.0                    # binari (go get 1.17-dən DEPRECATE)
```

## 2. HTTP + JSON klient (kanonik)

```go
spaceClient := http.Client{Timeout: time.Second * 2}
req, _ := http.NewRequest(http.MethodGet, url, nil)
req.Header.Set("User-Agent", "app-name")
res, err := spaceClient.Do(req)
defer res.Body.Close()
body, _ := io.ReadAll(res.Body)
err = json.Unmarshal(body, &p)   // &p UNUTMA; sahə BÖYÜK hərf + `json:"tag"`
```

## 3. CLI: flag → Cobra

```go
flag.StringVar(&inputVar, "message", "", "desc")
flag.Parse()
if len(strings.TrimSpace(secretVar)) == 0 { panic("required") }
```
- Verb-noun: `faas-cli deploy --image=...` (Cobra subcommands)
- 5 prinsip: Go statik binari; sadə başlanğıc; avtomatlaşdır; paket
  menecerləri; geri bağlantı

## 4. Cross-compile + release

```bash
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o app-arm64 \
    -ldflags "-s -w -X main.Version=$(git describe --tags) -X main.GitCommit=$(git rev-parse HEAD)"
```
- `-s -w` strip (~30% kiçik); CGO=0 portativlik; `file`/`du -h` ilə yoxla
- Makefile dist 5 platform; GitHub Actions tag → upload-assets

## 5. Unit-testlər

```go
func Test_sum(t *testing.T) {
    tables := []struct{ x, y, want int }{{1, 1, 2}, {5, 2, 7}}
    for _, tc := range tables {
        if got := sum(tc.x, tc.y); got != tc.want {
            t.Errorf("sum(%d+%d): want %d, got %d", tc.x, tc.y, tc.want, got)
        }
    }
}
```
| Bayraq | İş |
|---|---|
| -v | detallı |
| -cover -coverprofile=c.out + go tool cover -html | örtük + vizual |
| -bench=. | b.N dövrü, ns/op |
| -test.count=100 | stress (race ifşa) |
| -test.parallel N + t.Parallel() | 5s→1.3s |
| -c | test binariləri |

- **İlk test SİNDİRilməlidir**; assertion YOX; logu azalt

## 6. Asılılıq izolyasiyası + handler test

```go
type GetWebRequest interface { FetchBytes(url string) ([]byte, error) }
// fake: testWebRequest{} → hazır JSON qaytarır
w := httptest.NewRecorder()
r := httptest.NewRequest(http.MethodGet, "http://localhost:8080", nil)
handler(w, r)
if w.Code != wantCode { t.Fatalf(...) }
```

## 7. Concurrency pilləliyi

```go
// 1) WaitGroup
wg.Add(n); go func() { defer wg.Done(); work() }(); wg.Wait()
// 2) loop tələsi — 2 həll:
go func(j int) { ... }(i)   // VƏ YA  j := i; go func() { ... }()
// 3) errgroup
g.Go(func() error { return err })   // g.Wait() → ilk xəta
// 4) singleflight — eyni key-ləri birləşdir
res, err, shared := s.Do(key, func() (any, error) { return expensive() })
// 5) RWMutex — N oxu + 1 yazı
lock.RLock(); defer lock.RUnlock()   // oxu
lock.Lock(); defer lock.Unlock()     // yazı
// 6) worker pool
workQueue := make(chan string)
for i := 0; i < workers; i++ { go func() { for u := range workQueue { ... } }() }
go func() { for _, u := range uris { workQueue <- u }; close(workQueue) }()
// 7) context
ctx, cancel := context.WithTimeout(ctx, 5*time.Second); defer cancel()
http.NewRequestWithContext(ctx, ...)
```

## 8. Konfiq

```go
yaml.Unmarshal(bytesOut, &spec)             // `yaml:"name,omitempty"`
mergo.Merge(&merged, override, mergo.WithOverride)   // base+prod
```

## 9. embed (1.16+)

```go
//go:embed schema.sql
var schema string
//go:embed public/*
var frontend embed.FS
http.FileServer(http.FS(frontend))
```

## 10. HTTP server + middleware

```go
s := &http.Server{Addr: ":8080", ReadTimeout: 10 * time.Second,
    WriteTimeout: 10 * time.Second, MaxHeaderBytes: 1 << 20}
r := mux.NewRouter()
r.HandleFunc("/customer/{id:[-a-zA-Z_0-9.]+}", h).Methods("GET")
func makeAuth(token string, next http.HandlerFunc) http.HandlerFunc {
    return func(w, r) { if !valid { http.Error(...); return }; next(w, r) }
}
```

## 11. Prometheus

```go
// RED avtomatik:
InstrumentHandlerCounter(counter, InstrumentHandlerDuration(hist, next))
// custom inflight gauge:
prometheus.MustRegister(g)
g.Inc(); defer g.Dec()
// dinamik label-lar (K8s replicas):
Collect: Gauge.Reset() → WithLabelValues(name).Set(v) → Gauge.Collect(ch)
// sorğu: rate(http_requests_total{code="200"}[1m])
```

## 12. Docker multi-arch

```dockerfile
FROM --platform=$BUILDPLATFORM golang:1.17 AS builder
ARG TARGETOS; ARG TARGETARCH
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -ldflags "..." -o /app
FROM gcr.io/distroless/static:nonroot
COPY --from=builder /app /
USER nonroot:nonroot
```
```bash
docker buildx build --platform linux/amd64 --platform linux/arm64 --push -t u/a:v .
```

## 13. Tələsiklər cədvəli

| Tələsik | Həll |
|---|---|
| go run-da goroutine çıxışı görunmür | wg.Wait əskik |
| Loop-da hamısı son dəyər | parametr/j := i |
| Testdə race | -test.count=100 |
| Kanal deadlock | buffered/close-sonra-range |
| DB parol kodda | flags/env/secrets faylı |
| Köhnə binari istifadə olunur | -version + ldflags |
| Böyük image | multi-stage + distroless |
| Thundering herd | singleflight |
| Goroutine leak | rate(go_goroutines[1m]) monitorinq |
