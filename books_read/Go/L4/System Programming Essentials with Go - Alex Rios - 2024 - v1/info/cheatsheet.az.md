# System Programming Essentials with Go — Cheat Sheet (Azərbaycanca)

> **Kitab:** System Programming Essentials with Go — Alex Rios, Packt, 2024 (ISBN 978-1-83512-772-8) · 🚀 Advanced (4/5)
> Sistem çağırışlarından paylanmış keşə qədər bütün əsas kod/konfiqurasiya — bir baxışda.

---

## 1. Sistem çağırışları (Ch 3)
```go
// Səviyyələr: syscall (donaq) → golang.org/x/sys (müasir) → os (portativ)
syscall.Syscall(syscall.SYS_...)          // aşağı səviyyə
unix.EpollCreate1(0)                     // x/sys/unix
os.Open("fayl")                          // yüksək səviyyə abstraksiya
```
**Tracing:** `strace -c ./app` (statistika), `strace -e trace=openat ./app` (filtrlə).
**Standart stream-lər:** stdin=0, stdout=1, stderr=2; `> fayl 2>&1` = hamısı bir yerdə.

## 2. Fayl əməliyyatları (Ch 4)
```go
info, _ := os.Stat(path)                 // mode: info.Mode().String()
filepath.WalkDir(root, func(p string, d os.DirEntry, err error) error {...})
os.Symlink("hədəf", "link")              // symlink
os.Remove(path)                          // unlink
```
**Təhlükəsiz permission-lər:** dir 0777/0775, fayl 0666/0600; `0o777 &^ umask`.
**Dublikat tapmaq:** ölçü → hash (SHA-256) qısayolu.

## 3. Signal + scheduling + fayl izləmə (Ch 5)
```go
sigs := make(chan os.Signal, 1)
signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)
go func() { <-sigs; cleanup(); os.Exit(0) }()

exec.CommandContext(ctx, "cmd")          // timeout idarəli prosses
```
**fsnotify (inotify wrapper):**
```go
watcher, _ := fsnotify.NewWatcher()
watcher.Add("/var/log")
for { select { case ev := <-watcher.Events: ...; case err := <-watcher.Errors: ... } }
```

## 4. Pipes (Ch 6)
```go
r, w, _ := os.Pipe()                     // anonymous pipe
cmd.Stdout = w; cmd2.Stdin = r           // prosesləri birləşdir
// Named pipe (FIFO): unix.Mkfifo("/tmp/pipe", 0o666) — byte-stream, bloklanır
```

## 5. Unix domain socket (Ch 7)
```go
ln, _ := net.Listen("unix", "/tmp/app.sock")
conn, _ := net.Dial("unix", "/tmp/app.sock")
// HTTP over UDS:
u, _ := url.Parse("http://localhost")
u.Host = ""                              // host skip
tr := &http.Transport{ DialContext: func(...) { net.Dial("unix", sock) } }
// Inspeksiya: lsof -U | grep app.sock
```

## 6. GC tuning (Ch 8)
```bash
GOGC=200            # heap artım limiti (100 = default, 2x)
GOMEMLIMIT=4GiB     # yumşaq yaddaş limiti (Go 1.19+)
GODEBUG=gctrace=1    # GC dövrü trace
```
**Ballast:** `ballast := make([]byte, 1<<30)` — süni heap ilə GC tezliyini azalt.
**Arena (eksperimental):** `GOEXPERIMENT=arenas` — manual lifecycle, `arena.New()`.

## 7. Escape analysis + benchmark + pprof (Ch 9)
```bash
go build -gcflags="-m" ./...     # "escapes to heap" səbəbləri
go test -bench=. -benchmem      # ns/op + allocs/op
go test -cpuprofile=cpu.out -memprofile=mem.out .
go tool pprof cpu.out            # top, list, web
```
**Qısayollar:**指针→heap ehtimalı; `&x` funksiyadan çıxırsa escape.

## 8. Şəbəkə (Ch 10)
```go
ln, _ := net.Listen("tcp", ":8080")
conn, _ := ln.Accept()                   // hər bağlantı: go handle(conn)
// TLS: openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365
http.ListenAndServeTLS(":8443", "cert.pem", "key.pem", nil)
```
**Status xəritəsi:** GET→200, POST→201, PUT→202, DELETE→204, naməlum→405.
**Şəbəkə bayt sırası:** `binary.BigEndian.PutUint32(buf[:4], seq)`.

## 9. Telemetriya (Ch 11)
```go
// slog (Go 1.21+):
logger := slog.New(slog.NewJSONHandler(os.Stdout))
logger.Info("mesaj", slog.Attr("key", "val"))
// Prometheus counter:
c := prometheus.NewCounterVec(prometheus.CounterOpts{Name: "http_requests_total"},
    []string{"status_code"})
c.WithLabelValues("200").Inc()
// OTel trace:
tp := sdktrace.NewTracerProvider(sdktrace.WithBatcher(otlptracehttp.NewClient()))
otel.SetTracerProvider(tp)
```
**Metric seçimi:** say→Counter, səviyyə→Gauge, distribusiya→Histogram, dinamik quantile→Summary.

## 10. Modules + CI + Release (Ch 12)
```bash
go mod init github.com/user/proj
go get dep@v1.2.3; go mod tidy
go work init ./proj; go work use ./modA; go work sync
go env -w GOPRIVATE="github.com/<org>/*"
go install github.com/user/tool@v0.5.0
```
**GitHub Actions cache:**
```yaml
- uses: actions/cache@v3
  with:
    path: ~/go/pkg/mod
    key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
```
**GoReleaser:** CGO_ENABLED=0, ldflags `-s -w -extldflags "-static"`, goos/goarch cross-compile, GHCR docker tags; tag push → release workflow.

## 11. Distributed cache capstone (Ch 13)
```go
// LRU: map + container/list
type Cache struct {
    mu sync.RWMutex
    items map[string]*list.Element
    eviction *list.List
    capacity int
}
// Get: MoveToFront; Set: PushFront + evictLRU (Back)
// TTL: time.NewTicker + evictExpiredItems goroutine
// Replikasiya: X-Replication-Request header (sonsuz dövr qarşısı)
// Sharding: hash ring — SHA-1, sorted hashes, binary search, saat əqrəbi ilk node
// Forward: X-Forwarded-For loop qoruması
```

## 12. Effektiv praktikalar (Ch 14)
```go
// sync.Pool:
var pool = sync.Pool{New: func() any { return new(bytes.Buffer) }}
buf := pool.Get().(*bytes.Buffer)
defer func() { buf.Reset(); pool.Put(buf) }()

// sync.OnceValue (Go 1.21):
var getConfig = sync.OnceValue(loadConfig)

// singleflight — paralel eyni sorğu = 1 icra:
var g singleflight.Group
res, err, shared := g.Do(key, fn)

// mmap:
reader, _ := mmap.Open("böyük_fayl")
syscall.Msync(data, syscall.MS_SYNC)     // diskə məcburi yaz
```
**Tələlər:** time.After leak → `time.NewTimer + defer Stop`; defer loop-da → iterasiya içində Close; map şişir → yeni map-ə köçür; `resp.Body` → dərhal `defer Close`; kanal qəbulçusuz → select+timeout.

## 13. Hardware avtomatlaşdırma (Appendix)
```go
// D-Bus USB hadisəsi:
conn, _ := dbus.SystemBus()
ch := make(chan *dbus.Signal); conn.Signal(ch)
conn.BusObject().Call("org.freedesktop.DBus.AddMatch", 0,
    "type='signal',sender='org.freedesktop.UDisks2',interface='org.freedesktop.DBus.ObjectManager',path='/org/freedesktop/UDisks2'")
// Mount: /proc/mounts (6 sahə; \\040 = boşluq) VƏYA Properties.GetAll MountPoints
// Bluetooth RSSI: api.GetDefaultAdapter().GetDevices() → info.RSSI < -70 → xdg-screensaver lock
```
