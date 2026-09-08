# Chapter 9 — Goroutines, Channels, and Context (səh. 138-155)

## Bu fəsil nədən bəhs edir?

CSP əsaslı konkurensiya: sequential kodun paralelə çevrilməsi (fan-out/fan-in),
buffered channel ilə goroutine limiti, worker pool nümunəsi, Context ilə
timeout-lar (select + Done) və kontekstdə logger/request ID daşıma.

## Əsas fikirlər

### Recipe 47 — sequential → parallel
**Tapşırıq:** URL-lərin gecikməsini ölç — sequential 3.1s → paralel.

```go
// Info — nəticə İÇİNƏ url və err daxil edildi (hansı goroutine bitirəcək bilinmir):
type Info struct {
    url        string
    statusCode int
    delay      time.Duration
    err        error
}

// siteInfo DƏYİŞMİRDİ — sadə, test edilə bilən qalır

func sitesInfo(urls []string) (map[string]Info, error) {
    // Fan-out: hər URL üçün goroutine:
    ch := make(chan Info)
    for _, url := range urls {
        go func(u string) {              // loop dəyişənini parametr kimi ötür!
            info, err := siteInfo(u)
            info.err = err
            info.url = u
            ch <- info
        }(url)
    }

    // Fan-in: dəqiq len(urls) nəticə topla:
    out := make(map[string]Info)
    for range urls {
        info := <-ch
        if info.err != nil {
            return nil, info.err
        }
        out[info.url] = info
    }
    return out, nil
}
// Nəticə: 3.1s → 1.18s (≈2.6x)
```
- Amdahl qanunu: paralellik sehr deyil — teorik speedup formulunu oxuyun
- Concurrency = optimizasiya növü; sequential kifayətdirsə — dəyişməyin

### Recipe 48 — buffered channel ilə goroutine limiti
**Tapşırıq:** PNG-ləri resize et — CPU-intensive, core sayından çox
goroutine mənasızdır.

```go
func scaleDir(srcDir, destDir string, size image.Rectangle) error {
    var mu sync.Mutex
    var errs error

    pool := make(chan bool, runtime.GOMAXPROCS(0))  // = CPU core sayı
    var wg sync.WaitGroup
    srcFiles, err := filepath.Glob(filepath.Join(srcDir, "*.png"))
    if err != nil {
        return err
    }
    wg.Add(len(srcFiles))
    for _, src := range srcFiles {
        go func(srcFile string) {
            defer wg.Done()
            pool <- true                  // slot al (doluysa bloklanır)
            defer func() { <-pool }()     // slotu QAYTAR — defer ilə zəmanət
            destFile := path.Join(destDir, filepath.Base(srcFile))
            if err := scaleFile(srcFile, destFile, size); err != nil {
                mu.Lock()
                errs = errors.Join(errs, err)   // xətaları birləşdir
                mu.Unlock()
            }
        }(src)
    }
    wg.Wait()
    return errs
}
```
- `runtime.GOMAXPROCS(0)` → core sayı; >0 → limit qoyar
- Semaphore pattern-i: buffered channel + defer ilə slot azad etmə
- `errors.Join` — bir çox goroutine xətasını bir error-da birləşdirir

### Recipe 49 — worker pool
**Tapşırıq:** çoxsaylı saytları izlə; konkurensiya həddi istifadəçi təyin edir.

```go
type workerInfo struct {
    Info
    err error
}

type infoReq struct {
    url string              // URL to query
    ch  chan<- workerInfo   // return channel
}

func infoWorker(ch <-chan infoReq) {
    for req := range ch {                 // queue bağlanana qədər
        info, err := siteInfo(req.url)
        winfo := workerInfo{Info: info, err: err}
        req.ch <- winfo                    // nəticəni qaytar
    }
}

// Pool is a fixed pool of workers.
type Pool struct {
    queue chan infoReq
}

func NewPool(n int) (*Pool, error) {
    if n <= 0 {
        return nil, fmt.Errorf("n must be > 0 (got %d)", n)
    }
    queue := make(chan infoReq)
    for i := 0; i < n; i++ {
        go infoWorker(queue)               // n işçi başlat
    }
    return &Pool{queue: queue}, nil
}

// Close signals the worker goroutines to terminate.
func (p *Pool) Close() error {
    if p.queue != nil {
        close(p.queue)                     // range döngüləri bitir
        p.queue = nil
    }
    return nil
}

// SiteInfo returns a channel with info.
func (p *Pool) SiteInfo(url string) (Info, error) {
    ch := make(chan workerInfo, 1)         // BUFFERED(1) — leak qarşısı!
    p.queue <- infoReq{url, ch}
    info := <-ch
    return info.Info, info.err
}
```
- API sadə görünür (parametr + nəticə), arxada pool işləyir
- **Goroutine leak:** bloklanmış goroutine → kanal GC-edilmir + stack
  yaddaş; Dave Cheney: **"Never start a goroutine without knowing how it
  will stop."**
- Kanal tiplərində İSTİQAMƏT göstərin (`chan<-`, `<-chan`) — compiler
  səhvləri yaxalayır

### Recipe 50 — Context ilə timeout
**Tapşırıq:** SiteInfo üçün timeout — hədd aşarsa xəta.

```go
type infoReq struct {
    ctx context.Context          // sorğu konteksti
    url string
    ch  chan<- workerInfo
}

func infoWorker(ch <-chan infoReq) {
    for req := range ch {
        outCh := make(chan workerInfo, 1)
        go func() {                          // işçi goroutine (ləğv olunmur!)
            info, err := siteInfo(req.ctx, req.url)
            outCh <- workerInfo{Info: info, err: err}
        }()

        select {                              // yarış: nəticə YOXSA timeout
        case info := <-outCh:
            req.ch <- info
        case <-req.ctx.Done():                // timeout/deadline/cancel
            req.ch <- workerInfo{err: req.ctx.Err()}
        }
    }
}

func (p *Pool) SiteInfo(ctx context.Context, url string) (Info, error) {
    ch := make(chan workerInfo, 1)

    // Göndərmə timeout-u:
    select {
    case p.queue <- infoReq{ctx, url, ch}:
        // növbəyə düşdü
    case <-ctx.Done():
        return Info{}, ctx.Err()
    }

    // Qəbul timeout-u:
    select {
    case info := <-ch:
        return info.Info, info.err
    case <-ctx.Done():
        return Info{}, ctx.Err()
    }
}

// Aşağı səviyyə — kontekst HTTP-yə ötürülür:
func siteInfo(ctx context.Context, url string) (Info, error) {
    req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
    // ...
    resp, err := http.DefaultClient.Do(req)   // ləğv oluna bilən!
    // ...
}
```

**Qızıl qaydalar:**
- `ctx` — İLK parametr; struct-da SAXLAMA (rəsmi tövsiyə)
- Timeout YOX, Context qəbul et → istifadəçi deadline/timeout/cancel
  seçir
- select + `<-ctx.Done()` — yarış pattern-i
- Goroutine-i xaricdən ləğv etmək mümkün DEYİL — kontekst bitəndə də
  daxili goroutine işi bitirəcək (kanal buferli olduğundan leak yoxdur)
- net/http köhnə API — `NewRequestWithContext` birləşdirici

### Recipe 51 — kontekstdə logger + request ID
**Tapşırıq:** paralel handler logları qarışır — hər sorğuya ID.

```go
// XID — mərkəzi kilidsiz unikal ID:
func newID() string {
    return xid.New().String()
}

// Toqquşmasız açar — öZ tip:
type ctxKey string
const valuesKey ctxKey = "ctxArgs"

type Values struct {
    RequestID string
    Logger    *log.Logger
}

func idLogger(id string) *log.Logger {
    prefix := fmt.Sprintf("<%s> %s", id, log.Prefix())
    return log.New(log.Writer(), prefix, log.Flags())
}

// Kontekstdən çıxar (comma-ok!):
func ctxLogger(ctx context.Context) *log.Logger {
    vals, ok := ctx.Value(valuesKey).(*Values)
    if !ok {
        return stdLogger()          // standart logger — fallback
    }
    if vals.Logger == nil {
        panic(fmt.Sprintf("no logger in %#v", vals))
    }
    return vals.Logger
}

// Handler-də:
timeout := 100 * time.Millisecond
ctx, cancel := context.WithTimeout(context.Background(), timeout)
defer cancel()

id := newID()
logger := idLogger(id)
logger.Printf("info: usersHandler: id=%s", id)

values := Values{RequestID: id, Logger: logger}
ctx = context.WithValue(ctx, valuesKey, &values)   // pointer — yuxarı da
users, err := getAllUsers(ctx)
```
- Go-da thread-local storage YOXDUR — kontekst daşıyıcıdır
- Öz ctxKey tipin → pakitlər arası toqquşma mümkün deyil
- Bir çox dəyər üçün: tək Values struct + tək kontekst girişi
- Əlavə məlumatı (IP, user-agent) BİR dəfə ID ilə logla; sonra yalnız ID

## Əsas terminlər
- CSP (Communicating Sequential Processes) — Tony Hoare, 1978
- "Share memory by communicating" — Go proverbü
- Fan-out / fan-in — goroutine yayımı / nəticə yığılması
- Amdahl's law — paralel speedup limiti
- Semaphore (buffered channel) — goroutine həddi
- Worker pool — sabit işçi dəsti + iş növbəsi
- Goroutine leak — əbədi bloklanmış goroutine
- errors.Join — xətaların birləşdirilməsi
- ctx.Done() — ləğv siqnal kanalı
- NewRequestWithContext — ləğv olunan HTTP
- ctxKey idiomu — toqquşmasız kontekst açarı
- "Never start a goroutine without knowing how it will stop"

## Praktik nəticə
Sequential-dən başlayın; baseline ölçün; yalnız lazım olanda paralelləşdirin.
Strukturlu pattern-lər: fan-out/fan-in (nəticəyə url/err daxil edin),
semaphore (buffered channel + defer slot qaytarma), worker pool (queue +
n işçi + Close). Hər əməliyyat üçün ctx ilk parametr; select+Done ilə
yarış; HTTP-yə NewRequestWithContext. Logger/ID kimi dəyərlər üçün öz
ctxKey tipi + Values struct + WithValue. Loop dəyişənini goroutine closure-a
parametr kimi ötürün.

## Mənbə
Pages: 138-155 (PDF 138-155)
