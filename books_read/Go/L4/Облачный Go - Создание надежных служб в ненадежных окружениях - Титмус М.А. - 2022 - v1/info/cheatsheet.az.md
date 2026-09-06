# Облачный Go / Создание надежных служб в ненадежных окружениях — Cheatsheet (Azərbaycanca)

> **Kitab:** Облачный Go / Создание надежных служб в ненадежных окружениях — Мэтью А. Титмус, ДМК Пресс 2022 · 12 chapter · Cloud Native + Go Patterns · Sürətli istinad.

---

## Chapter 4 — Cloud Application Programming Patterns

```go
// Context — sorğu həyat dövrü idarəetməsi
func Stream(ctx context.Context, out chan<- Value) error {
    dctx, cancel := context.WithTimeout(ctx, time.Second*10)
    defer cancel()
    res, err := SlowOperation(dctx)
    if err != nil { return err }
    for {
        select {
        case out <- res:
        case <-ctx.Done():
            return ctx.Err()
        }
    }
}

// Circuit Breaker — xətaları aşkar edib dayandırma
type Circuit func(context.Context) (string, error)

func Breaker(circuit Circuit, threshold uint) Circuit {
    var fails int; var last = time.Now(); var m sync.RWMutex
    return func(ctx context.Context) (string, error) {
        m.RLock(); d := fails - int(threshold)
        if d >= 0 {
            if !time.Now().After(last.Add(time.Second * 2 << d)) {
                m.RUnlock(); return "", errors.New("unreachable")
            }
        }
        m.RUnlock()
        r, err := circuit(ctx)
        m.Lock(); defer m.Unlock()
        last = time.Now()
        if err != nil { fails++; return r, err }
        fails = 0; return r, nil
    }
}

// DebounceFirst — ilk çağırışı icra et, sonrakıları interval ərzində ignore et
func DebounceFirst(circuit Circuit, d time.Duration) Circuit {
    var threshold time.Time; var result string; var err error; var m sync.Mutex
    return func(ctx context.Context) (string, error) {
        m.Lock(); defer func() { threshold = time.Now().Add(d); m.Unlock() }()
        if time.Now().Before(threshold) { return result, err }
        return circuit(ctx)
    }
}

// Retry — uğursuz əməliyyatları təkrar cəhd et
type Effector func(context.Context) (string, error)

func Retry(e Effector, retries int, delay time.Duration) Effector {
    return func(ctx context.Context) (string, error) {
        for i := 0; i < retries; i++ {
            if r, err := e(ctx); err == nil { return r, nil } { time.Sleep(delay) }
        }
        return e(ctx)
    }
}
```

---

## Chapter 5 — Building a Cloud Service

```go
// gorilla/mux — RESTful marşrutlaşdırma
r := mux.NewRouter()
r.HandleFunc("/keys/{key}", getHandler).Methods("GET")
r.HandleFunc("/keys/{category}/{id:[0-9]+}", getByCategory).Methods("GET")

// Transaction Logger interfeysi
type TransactionLogger interface {
    WritePut(key, value string)
    WriteDelete(key string)
    Run()
    Err() <-chan error
}

// Fayl əsaslı qeydiyyatçı — append-only, tab-separated
fmt.Fprintf(file, "%d\t%s\t%s\t%s\n", seq, eventType, key, value)

// Docker multi-stage build (Dockerfile)
// Stage 1: build
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o kvstore ./cmd/kvstore
// Stage 2: minimal image
FROM scratch
COPY --from=builder /app/kvstore /kvstore
ENTRYPOINT ["/kvstore"]
```

---

## Chapter 7 — Scalability

```go
// Stateless dizayn — sorğu öz data ilə gəlir, serverə bağlı vəziyyət yoxdur

// Gorutin sızıntısını qarşısı al
func pollWithGoroutine(ctx context.Context, resources []*Resource) {
    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()
    for {
        select {
        case <-ticker.C:
            // polling məntiqi
        case <-ctx.Done():
            return
        }
    }
}

// Goroutine exit planı şərt — hər gorutin çıxış rejasına malik olmalıdır
// Sync tiplər (Mutex, WaitGroup) HEÇ VAXT kopyalanmaz — pointer istifadə et
```

---

## Chapter 8 — Loose Coupling

```go
// REST — uniform interfeys, HTTP metodları, JSON/XML
// gRPC — contract-first, .proto → Go stub, HTTP/2, binary serialization

// gRPC interfeysi (.proto'dan avtomatik yaradılır)
type KeyValueStoreClient interface {
    Get(ctx context.Context, in *GetRequest) (*GetResponse, error)
    Put(ctx context.Context, in *PutRequest) (*PutResponse, error)
}

// Heksagonal Arxitektura — core logic, ports, adapters
type KeyValueStore interface {
    Get(key string) (string, error)
    Put(key, value string) error
}

type FrontEnd interface {
    Start(kv *core.KeyValueStore) error
}

// REST adapter
type restFrontEnd struct { router *mux.Router }

// gRPC adapter
type grpcFrontEnd struct { server *grpc.Server }
```

---

## Chapter 9 — Resilience

```go
// Eksponensial artım + jitter ilə təkrar cəhd
func backoff(attempt int, base time.Duration) time.Duration {
    jitter := time.Duration(rand.Int63n(int64(base)))
    return base * time.Duration(1<<uint(attempt)) + jitter
}

// Health check endpoint-ləri
func healthHandler(w http.ResponseWriter, r *http.Request) {
    if err := db.Ping(); err != nil {
        http.Error(w, "db down", http.StatusServiceUnavailable)
        return
    }
    w.WriteHeader(http.StatusOK)
}

// Liveness Probe — proses canlıdırmı?
// Readiness Probe — sorğuları qəbul edə bilirmi?
// Deep Probe — DB/cache çalışır?
```

---

## Chapter 10 — Manageability

```go
// Viper konfiqurasiya iyerarxiyası
viper.SetDefault("id", "13")           // 1. Defaults
viper.SetConfigFile("config.yaml")      // 2. Config faylı
viper.BindEnv("PORT")                   // 3. Env vars
viper.BindPFlags(rootCmd.PersistentFlags()) // 4. CLI flags

// Dinamik yeniləmə
viper.WatchConfig()
viper.OnConfigChange(func(e fsnotify.Event) {
    // konfiqurasiya dəyişdi, yenilə
})

// Cobra CLI
var rootCmd = &cobra.Command{
    Use: "myapp",
    Run: func(cmd *cobra.Command, args []string) {
        port := viper.GetInt("port")
    },
}
rootCmd.AddCommand(&cobra.Command{
    Use: "serve",
    Run: func(cmd *cobra.Command, args []string) { /* serve logic */ },
})

// Feature flag — dinamik dəyər
type Enabled func(flag string, r *http.Request) (bool, error)
func fromPrivateIP(flag string, r *http.Request) (bool, error) {
    ip, _, _ := net.SplitHostPort(r.RemoteAddr)
    _, cidr, _ := net.ParseCIDR("10.0.0.0/8")
    return cidr.Contains(net.ParseIP(ip)), nil
}
```

---

## Chapter 11 — Observability

```go
// Zap strukturlaşdırılmış jurnalçılıq
logger, _ := zap.NewProduction()
defer logger.Sync()
logger.Info("request processed",
    zap.String("method", "GET"),
    zap.Int("status", 200),
    zap.Duration("latency", 150*time.Millisecond),
)

// OpenTelemetry — izləmə (tracing)
import "go.opentelemetry.io/otel"

tracer := otel.Tracer("myapp")
ctx, span := tracer.Start(ctx, "GetResource")
defer span.End()

span.AddEvent("Acquiring lock",
    label.String("lock.name", "cache"),
    label.Int("lock.id", 42),
)
span.SetAttributes(
    label.String("resource.name", "database"),
    label.Int("http.status_code", 200),
)

// Prometheus metrikaları
exporter, _ := prometheus.NewExportPipeline()
meter := metric.Must(exporter.Meter("myapp"))
requestCounter := metric.Must(meter.NewFloat64Counter("http.requests.total"))
requestCounter.Add(ctx, 1, label.String("method", "GET"))
```

---

## Sürətli yaddaş cədvəli

| Ehtiyac | Həll |
|---------|------|
| Sorğu timeout/ləğv | `context.WithTimeout`, `ctx.Done()` |
| Xəta idarəetmə | Circuit Breaker (`sync.RWMutex` + closure) |
| Çoxlu çağırış məhdudlaşdırma | DebounceFirst (`sync.Mutex` + threshold) |
| Təkrar cəhd | Retry + eksponensial artım + jitter |
| Gorutin sızıntısı qarşısı | `defer ticker.Stop()`, kanal bağlama |
| REST marşrut | `gorilla/mux` — `{key}`, `{id:[0-9]+}` |
| gRPC xidmət | `.proto` → Go stub, `grpc.DialContext` |
| Konteyner | Docker multi-stage build, scratch image |
| Konfiqurasiya | Viper iyerarxiyası + Cobra CLI |
| Feature flag | `func(string, *http.Request) (bool, error)` |
| Log | `go.uber.org/zap` — strukturlaşdırılmış key/value |
| Metrics | OpenTelemetry + Prometheus ixracatçısı |
| Trace | `otel.Tracer`, span events, labels |
| Health check | `/healthz` — 200 OK, `/readyz` — readiness, `/livez` — liveness |
