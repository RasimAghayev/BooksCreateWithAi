# Chapter 15 — Shipping Your Code (səh. 234-252)

## Bu fəsil nədən bəhs edir?

Production-a çıxarış: konfiqurasiya (ardanlabs/conf, validasiya), asılılıq
patch-i (replace direktivi), Docker multistage build, signal-larla graceful
shutdown, səmərəli logging (zap Check), expvar metrics (middleware) və
Delve ilə canlı servis debug-u.

## Əsas fikirlər

### Recipe 78 — konfiqurasiya
**Tələb:** env (APP_ prefiks) və ya CLI: `APP_ADDR=localhost:9999` /
`--web-addr=...`.

```go
import "github.com/ardanlabs/conf/v3"    // /v3 — modul versiyası!

var cfg struct {
    Web struct {
        Addr string `conf:"default::8080,env:ADDR"`
    }
}

help, err := conf.Parse("APP", &cfg)     // prefiks APP → APP_ADDR
if err != nil {
    if errors.Is(err, conf.ErrHelpWanted) {
        fmt.Println(help)
        os.Exit(0)
    }
    log.Fatalf("error: bad config - %s", err)
}

// VALIDASİYA — başlanğıcdan ƏVVƏL:
if err := validateAddr(cfg.Web.Addr); err != nil {
    log.Fatalf("error: invalid config - %s", err)
}

func validateAddr(addr string) error {
    i := strings.Index(addr, ":")
    if i == -1 {
        return fmt.Errorf("%q: missing : in address", addr)
    }
    port, err := strconv.Atoi(addr[i+1:])
    if err != nil {
        return fmt.Errorf("%q: invalid port - %w", addr, err)
    }
    const maxPort = 65_535
    if port < 0 || port > maxPort {
        return fmt.Errorf("%q: invalid port number", addr)
    }
    return nil
}
```
- Konfiqurasiya iyerarxiyası: defaults → fayl → env → CLI
- Çox variant = çox test: 10 "bəli/xeyr" = 2^10 = 1024 kombinasiya!
- **İnsident hesabatlarının əksəriyyəti pis konfiqurasiyadan** — validate
  edin

### Recipe 79 — asılılıq patch-i (replace)
**Bug:** geo.Euclidean overflow — developer 3 ay yoldadır.

```go
// 1) Reproduksiya testi:
func TestEuclideanBug(t *testing.T) {
    d := geo.Euclidean(0, 0, 95e200, 168e200)
    if math.IsInf(d, 1) {
        t.Fatal(d)
    }
}

// 2) Paketi _patch/ qovluğuna kopyala + replace:
// go.mod:
require github.com/353solutions/geo v1.2.3
replace github.com/353solutions/geo => ./_patch/geo

// 3) Lokal kopiyada fix:
func Euclidean(x1, y1, x2, y2 float64) float64 {
    dx := x1 - x2
    dy := y1 - y2
    return math.Hypot(dx, dy)     // Hypot overflow-etmir!
}
```
- **Yanaşmalar:** Fork (import path dəyişir — pis), Vendor (böyük diff),
  Replace (minimal) — qısamüddətli üçün replace, uzunmüddətli vendor

### Recipe 80 — Docker multistage
```dockerfile
# Mərhələ 1 — BUILD:
FROM golang:1.20-bookworm AS build
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download          # ƏVVƏL — cache layer!
COPY . .
ENV CGO_ENABLED=0
RUN go build -o ./dbq

# Mərhələ 2 — DEPLOY:
FROM debian:bookworm-slim
COPY --from=build /build/dbq /usr/local/bin
# Don't run as root user
RUN groupadd -r app && useradd --no-log-init -r -g app app
USER app
CMD dbq
```
- golang 845MB → debian-slim 74.8MB (~11x kiçik)
- go.mod əvvəl kopyalanır → mod download cache-lənir
- Go binary → SDK lazım deyil (JVM-dən fərqli)
- root İSTİFADƏ ETMƏ — app user yarat

### Recipe 81 — graceful shutdown (signal)
```go
mux := http.NewServeMux()
mux.HandleFunc("/health", healthHandler)
addr := ":8080"
srv := &http.Server{Addr: addr, Handler: mux}

// Server goroutine-də:
go func() {
    log.Printf("server starting on %s", addr)
    err := srv.ListenAndServe()
    if err != nil && err != http.ErrServerClosed {
        log.Printf("error: listen - %s", err)
        os.Exit(1)
    }
}()

// Signal tutma:
ch := make(chan os.Signal, 1)                       // BUFFERED — şərt!
signal.Notify(ch, unix.SIGTERM, unix.SIGINT)

<-ch                                                // gözlə
log.Printf("shutting down")
ctx, cancel := context.WithTimeout(context.Background(), time.Second)
defer cancel()
if err := srv.Shutdown(ctx); err != nil {          // graceful!
    log.Printf("error: shutdown - %s", err)
}
```
- SIGINT (Ctrl-C) və SIGTERM (kill default) — ən çox işlədilənlər
- `http.ErrServerClosed` — graceful shutdown-da gözlənilən xəta YOX sayılır
- Shutdown ctx-i — bitməmiş sorğular üçün vaxt pəncərəsi

### Recipe 82 — logging performansı (zap Check)
**Bug:** Error level-də info log-lar görünmür, amma users.Load YİNE
çağırılır — parametrlər funksiya çağırışından ƏVVƏL qiymətlənir!

```go
// Yavaş (hər çağırışda users.Load!):
s.logger.Info("logged in",
    zap.Any("user", users.Load(uid)),
)

// DÜZGÜN — səviyyə yoxlaması:
if info := s.logger.Check(zap.InfoLevel, "logged in"); info != nil {
    info.Write(
        zap.Any("user", users.Load(uid)),   // yalnız enabled-ə
    )
}
```
- Bill Kennedy: "Log development-də kömək etmirsə, production-da əmin ol
  kömək etməyəcək"
- Dəbli log-larda bahalı çağırışları Check qorumasına alın

### Recipe 83 — expvar metrics (middleware)
```go
// Status kodu yadda saxlayan Writer (Recipe 38-in əsası):
type statusWriter struct {
    http.ResponseWriter
    statusCode int
}

func (s *statusWriter) WriteHeader(status int) {
    s.statusCode = status
    s.ResponseWriter.WriteHeader(status)
}

// Middleware:
func addMetrics(name string, h http.Handler) http.Handler {
    calls := expvar.NewInt(fmt.Sprintf("%s.calls", name))
    errors := expvar.NewInt(fmt.Sprintf("%s.errors", name))
    oks := expvar.NewInt(fmt.Sprintf("%s.oks", name))
    fn := func(w http.ResponseWriter, r *http.Request) {
        calls.Add(1)                            // atomic daxilən
        sw := statusWriter{w, http.StatusOK}
        h.ServeHTTP(&sw, r)
        if sw.statusCode >= http.StatusBadRequest {
            errors.Add(1)
        } else {
            oks.Add(1)
        }
    }
    return http.HandlerFunc(fn)
}

h := addMetrics("lookup", http.HandlerFunc(lookupHandler))
http.Handle("/lookup", h)
```
```bash
$ curl 'http://localhost:8080/debug/vars'
# "lookup.calls": 2, "lookup.errors": 1, "lookup.oks": 1
```
- expvar: JSON format, pull model (/debug/vars); Prometheus pull edir
- Metrics = production gözü + alerting (4am çağırışları üçün :)

### Recipe 84 — Delve ilə canlı debug
**Bug:** boş mesajlar DB-də — production nüsxəsi lokalda yaradıla bilmir.

```bash
# Prosesi tap:
$ ps aux | grep messag
# miki 608230 ... ./messaging     ← PID

# Qoşul (sudo — elevated tələb):
$ sudo dlv attach 608230
(dlv) b main.addHandler          # breakpoint
(dlv) c                           # davam et

# Başqa shell-dətet ver:
$ curl -d'BAD JSON' http://localhost:8080/messages
# → dlv breakpoint-də dayanır

(dlv) n    # next — addım-addım
# → BUG: http.Error-dən sonra RETURN YOXDUR!
(dlv) q
# Would you like to kill the process? [Y/n] Y
```
- Quraşdırma: `go install github.com/go-delve/delve/cmd/dlv@latest`
- IDE-də asandır; SSH/məsafə üçün TUI — help ilə öyrən
- Tapılan bug: `http.Error(w, ...)` sonrası `return` unudulub —
  uğursuz JSON-da da mesaj əlavə olunurdu!
- Fix proseduru: əvvəl FAILING TEST, sonra return

## Final Thoughts-dən

Production acıq yerdir: disk dolur, bağlantı timeout olur, istifadəçilər
əyləncəli şeylər edir. Metrics = göz; loglar = səbəb; təcrübə = daha az
4am pager çağırışı.

## Əsas terminlər
- Production readiness — istismara hazırlıq siyahısı
- conf.Parse — env/CLI/defaults birlikdə
- Konfiqurasiya validasiyası — başlanğıcdan əVVƏL
- replace direktivi — lokal asılılıq əvəzi
- math.Hypot — overflow-süz məsafə
- Multistage Docker — build + slim deploy
- Layer cache — go.mod əvvəl
- SIGINT/SIGTERM — dayandırma siqnalları
- signal.Notify — buffered kanala siqnal
- srv.Shutdown(ctx) — graceful shutdown
- zap Check — səviyyə yoxlaması (bahalı parametr qoruması)
- expvar — /debug/vars metrics
- Delve (dlv) — Go debugger; attach — canlı prosesə
- Breakpoint/TUI (b, c, n, q) — debug əmrləri

## Praktik nəticə
Ship-etməzdən əvvəl: konfiqurasiya (prefiksli env + CLI + VALIDASİYA),
Docker multistage (SDK-dən ayrı slim image, root-suz, CGO_ENABLED=0),
signal handler (buffered chan + Notify + Shutdown ctx). İstismarda:
metrics middleware (expvar/Prometheus) + Check-qli logs + Delve
attach üçün sudo ilə çıxış. Xarici asılılıq bug-ı: test yaz → _patch +
replace → upstream düzələndə qaldır. Bug tapanda: əvvəl failing test,
sonra fix.

## Mənbə
Pages: 234-252 (PDF 234-252)
