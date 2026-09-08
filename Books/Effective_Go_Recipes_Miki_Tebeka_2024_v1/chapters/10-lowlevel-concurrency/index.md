# Chapter 10 — Lower-Level Concurrency (səh. 156-169)

## Bu fəsil nədən bəhs edir?

sync və sync/atomic paketləri: sync.Once (idempotent konfiqurasiya),
sync.WaitGroup (paralel broadcast), sync.RWMutex (çox oxuyan/tək yazan),
race detector (`-race`), atomic.AddInt64 (yarış şəraiti fix-i) və
atomic.Value ilə 21x sürətli Now().

## Əsas fikirlər

### Recipe 52 — sync.Once (idempotent əməliyyatlar)
**Tapşırıq:** env-dən konfiqurasiya bir dəfə yüklənsin (idempotent).

```go
var (
    cfgOnce sync.Once
    Config  struct {
        ListenAddress string
        Verbose       bool
    }
)

func loadConfig(envPrefix string) {
    getEnv := func(name string) string {       // closure — prefiks daxili
        key := fmt.Sprintf("%s_%s", envPrefix, name)
        return os.Getenv(key)
    }
    addr := getEnv("ADDRESS")
    if len(addr) == 0 {
        Config.ListenAddress = ":8080"          // default
    } else {
        Config.ListenAddress = addr
    }
    verbose := getEnv("VERBOSE")
    if verbose == "1" || verbose == "yes" || verbose == "on" {
        Config.Verbose = true
    }
}

// LoadConfig loads configuration once from environment.
func LoadConfig(envPrefix string) {
    cfgOnce.Do(func() {                        // YALNIZ bir dəfə icra
        loadConfig(envPrefix)
    })
}
```
- **Idempotent:** bir neçə dəfə tətbiq = bir dəfənin nəticəsi (On düyməsi)
- Do funksiyasız arqument gözləyir → prefiks anonymous funksiya ilə ötürülür
- İstifadə sahələri: konfiqurasiya, ödəniş ikinci dəfə işləməməsi,
  birdəfəlik bildiriş

### Recipe 53 — sync.WaitGroup (paralel broadcast)
**Tapşırıq:** chat otağında bütün client-lərə paralel bildiriş — hamısı
bitəndə qayıt.

```go
// Room is a chat room.
type Room struct {
    clients []io.Writer
}

// Notify sends msg to all clients in parallel.
// It will return after all messages are sent.
func (r *Room) Notify(msg string) {
    var wg sync.WaitGroup
    wg.Add(len(r.clients))               // ƏVVƏLCƏDƏN hamısı üçün
    for _, c := range r.clients {
        go func(w io.Writer) {
            defer wg.Done()               // defer — panic-də belə zəmanət
            w.Write([]byte(msg))
        }(c)
    }
    wg.Wait()                             // hamısı bitənə qədər
}
```

**Test (Sink + mutex):**
```go
type Sink struct {
    mu       sync.Mutex
    messages []string
}
func (s *Sink) Add(message string) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.messages = append(s.messages, message)
}
```
- Add → Done → Wait sayğacı; sıfıra çatanda Wait azad olur
- `defer wg.Done()` — deadlock qorunması
- Xətalar da lazımdırsa: `golang.org/x/sync/errgroup`

### Recipe 54 — sync.RWMutex (çox oxuyan / tək yazan)
**Tapşırıq:** dinamik konfiqurasiya — oxuma sürətini pozmadan.

```go
var (
    cfgLock sync.RWMutex
    config  = make(map[string]string)
)

// Oxu — RLock: bir çox oxuyan eyni anda!
func GetConfig(key string) string {
    cfgLock.RLock()
    defer cfgLock.RUnlock()
    return config[key]
}

// Yazma — Lock: tək yazan, oxuyan yoxdur
func ReloadConfig() {
    cfgLock.Lock()
    defer cfgLock.Unlock()
    config["updated"] = time.Now().String()
}
```
- RWMutex: YA bir yazan (oxuyan yoxdur), YA çox oxuyan (yazan yoxdur)
- Adi Mutex-də oxuyanlar da növbələnir — konkurensiya israfı
- Map/slice goroutine-lər üçün təhlükəsizdir — kilid şərtdir
- **Tövsiyə:** konfiqurasiya dəyişəndə restart etmək daha təhlükəsizdir —
  yarımçıq konfiqurasiya riski

### Recipe 55 — race detector
**Bug:** 1000 hit → 987 metric; 10 goroutine × 1000 = 10000 → 9914 çap!

```go
// XƏTALI:
var callCount int                       // race condition!
func handler(w http.ResponseWriter, r *http.Request) {
    callCount++                          // ++ = oxu+yazı — atomik DEYİL
    // ...
}

// DETERMİNİZM:
$ go run -race counter.go
# WARNING: DATA RACE
# Read at 0x... by goroutine 12: main.handler() counter.go:21
# Previous write at 0x... by goroutine 15: main.handler() counter.go:21
# Found 2 data race(s)

// FIX:
var callCount int64                     // SPESİFİK tip (int yox!)
func handler(w http.ResponseWriter, r *http.Request) {
    atomic.AddInt64(&callCount, 1)       // memory barrier
    // ...
}
# $ go run -race counter.go → 10000 ✓
```

**Səbəb:** 3 səviyyəli cache — hər goroutine fərqli dəyər görə bilər;
`++` = oxu+yazı yarışı. Race çıxışı: hansı goroutine, hansı sətir, oxu/yazı.

**Mutex vs Atomic benchmark:**
```
BenchmarkMutex-12    44517554    25.59 ns/op
BenchmarkAtomic-12  144230883     8.343 ns/op   # ~3x sürətli
```
- sync/atomic adətən "code smell" — yüksək səviyyəli primitivlər üstündür;
  yalnız sıx performans tələbində
- Metrikalar üçün: expvar / Prometheus

### Recipe 56 — atomic.Value ilə sürətli Now()
**Problem:** profiling → log-larda `time.Now()` vaxtı yeyir.

```go
var (
    now atomic.Value                     // kilidsiz təhlükəsiz saxlanc
)

func init() {
    now.Store(time.Now())
    go func() {                          // arxa plan yeniləyicisi
        for {
            time.Sleep(time.Millisecond)
            now.Store(time.Now())        // hər 1ms
        }
    }()
}

// Now return the current time in 1ms granularity
func Now() time.Time {
    return now.Load().(time.Time)        // any → time.Time cast
}
```

**Benchmark:**
```
BenchmarkTimeNow-12   22844953    49.34 ns/op
BenchmarkNow-12     523512082     2.310 ns/op   # ~21x sürətli!
```
- Trade-off: dəqiqlik 1ms — log üçün kifayətdir
- atomic.Value.Store/Load — mutexsiz oxu/yazı; any saxladığından cast
  lazımdır

## Final Thoughts-dən

Aşağı səviyyəli primitivlər çox diqqət tələb edir — testlər tutmur,
production-da debug dəhşətdir. **Yaxşı səbəb olmadan goroutine + kanalda
qalın.**

## Əsas terminlər
- Idempotent — təkrar tətbiq nəticəni dəyişmir
- sync.Once.Do — birdəfəlik icra zəmanəti
- sync.WaitGroup — Add/Done/Wait sayğacı
- RWMutex — oxuyan-yazan kilidi (RLock/Lock)
- Race condition — paralel oxu/yazı yarışı
- Race detector (-race) — data race yoxlayıcısı
- Memory barrier — cache-ləri sinxronlayan əməliyyat
- atomic.AddInt64 — atomik artım
- atomic.Value — istənilən tip üçün atomik saxlanc
- Trade-off (dəqiqlik/sürət) — 1ms qranulluq ↔ 21x sürət
- expvar / Prometheus — metrik sistemləri

## Praktik nəticə
Birdəfəlik işlər üçün sync.Once (bool+mutex özün yazmayın); qrup gözləməsi
üçün WaitGroup (defer Done şərt); oxu-çox-yazı-az data üçün RWMutex;
paylaşılan map/slice həmişə kilidli. `-race` flag-i hər test/inteqrasiya
işində standart olmalıdır — data race gözləniləndən fərqli nəticələr verir.
Sıx performans: atomic.Add* (spesifik int64 tip!) və atomic.Value (cast ilə)
— amma yalnız benchmark sübut edəndə; əks halda kanal/sync üstünlük.

## Mənbə
Pages: 156-169 (PDF 156-169)
