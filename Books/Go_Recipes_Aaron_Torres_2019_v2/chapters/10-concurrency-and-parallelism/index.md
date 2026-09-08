# Chapter 10 — Параллелизм и конкурентность (Paralellik və konkurensiya)

## Bu chapter nədən bəhs edir?

Channel + select, sync.WaitGroup, mutex/atomic, context (cancel/timeout/
value), state-li channel-lər, worker pool pattern və worker pipeline.

## Əsas fikirlər

### 1. Channel-lər və select operatoru
**Nədir:** Goroutine-lər arası kommunikasiya; select bir neçə kanalı
eyni anda dinləyir.

**Kitabdan kod nümunəsi:**
```go
// Sender — done gələnə qədər "tick" göndərir:
func Sender(ch chan string, done chan bool) {
    t := time.Tick(100 * time.Millisecond)
    for {
        select {
        case <-done:               // bitirmə siqnalı
            ch <- "sender done."
            return
        case <-t:                   // 100ms keçdi
            ch <- "tick"
        }
    }
}

// Printer — context ilə idarə olunur:
func Printer(ctx context.Context, ch chan string) {
    t := time.Tick(200 * time.Millisecond)
    for {
        select {
        case <-ctx.Done():          // cancel/timeout
            fmt.Println("printer done.")
            return
        case res := <-ch:           // gələn mesaj
            fmt.Println(res)
        case <-t:                   // 200ms "tock"
            fmt.Println("tock")
        }
    }
}

// main — hər iki goroutine + idarə:
ch := make(chan string)
done := make(chan bool)
ctx, cancel := context.WithCancel(context.Background())
defer cancel()
go channels.Printer(ctx, ch)
go channels.Sender(ch, done)
time.Sleep(2 * time.Second)
done <- true      // sender-i dayandır
cancel()          // printer-i dayandır
```

**Sub-kod izahı:**
- `select` → hansı case hazır DISA onu icra edir (random seçim)
- `time.Tick` select daxilində BAŞQA case seçiləndə SIFIRLANIR — bu, yaygın
  tələdir (prioritet yoxdur)
- Bitirmə yolları: done kanalı YA context cancel — ikisi də göstərilib

### 2. sync.WaitGroup ilə async əməliyyatlar
**Nədir:** Müstəqil paralel tapşırıqları gözləmək — sayğac mexanizmi:
Add(1) artır, Done() azaldır, Wait() sıfıra düşməyi gözləyir.

**Kitabdan kod nümunəsi:**
```go
// Xəta aqreqasiyası — custom error tipi:
type CrawlError struct {
    Errors []string
}
func (c *CrawlError) Add(err error) {
    c.Errors = append(c.Errors, err.Error())
}
func (c *CrawlError) Error() string {
    return fmt.Sprintf("All Errors: %s", strings.Join(c.Errors, ","))
}
func (c *CrawlError) Present() bool { return len(c.Errors) != 0 }

// Crawl — bütün URL-ləri paralel yüklə:
func Crawl(sites []string) ([]int, error) {
    start := time.Now()
    wg := &sync.WaitGroup{}
    var resps []int
    cerr := &CrawlError{}
    for _, v := range sites {
        wg.Add(1)                        // sayğacı artır
        go func(v string) {              // v ARQUMENT kimi — loop dəyişəni deyil!
            defer wg.Done()              // goroutine bitəndə azalt
            resp, err := GetURL(v)
            if err != nil {
                cerr.Add(err)            // xətaları topla (qeyd: mutex istəyir)
                return
            }
            resps = append(resps, resp.StatusCode)
        }(v)
    }
    wg.Wait()                             // hamısı bitənə qədər blokla
    if cerr.Present() {
        return resps, cerr               // aqreqasiya xətası
    }
    return resps, nil
}
```

**Sub-kod izahı:**
- `wg.Add(1)` goroutine BAŞLAMAZDAN ƏVVƏL (döngü içində)
- Loop dəyişənini goroutine-ə ARQUMENT kimi ötür — əks halda shared vəziyyət
- Tövsiyə: praktikada xəta/cavab üçün kanal istifadə et — slice append
  race yarada bilər (kitabın sadələşdirilmiş nümunəsi)

### 3. Mutex və atomic əməliyyatlar
**Nədir:** Paralel map yazma təhlükəsizdir — RWMutex; sayğac kimi dəyərlər
üçün sync/atomic; bir dəfəlik init üçün sync.Once.

**Kitabdan kod nümunəsi:**
```go
// SafeMap — RWMutex ilə thread-safe map:
type SafeMap struct {
    m  map[string]string
    mu *sync.RWMutex
}
func (t *SafeMap) Set(key, value string) {
    t.mu.Lock()               // YAZMA kilidi — yalnız 1 yazıçı
    defer t.mu.Unlock()
    t.m[key] = value
}
func (t *SafeMap) Get(key string) (string, error) {
    t.mu.RLock()              // OXU kilidi — çox oxucu paralel
    defer t.mu.RUnlock()
    if v, ok := t.m[key]; ok {
        return v, nil
    }
    return "", errors.New("key not found")
}

// Ordinal — atomic + sync.Once:
type Ordinal struct {
    ordinal uint64
    once    *sync.Once
}
func (o *Ordinal) Init(val uint64) {
    o.once.Do(func() {                        // YALNIZ bir dəfə
        atomic.StoreUint64(&o.ordinal, val)
    })
}
func (o *Ordinal) GetOrdinal() uint64 {
    return atomic.LoadUint64(&o.ordinal)       // atomik oxu
}
func (o *Ordinal) Increment() {
    atomic.AddUint64(&o.ordinal, 1)             // atomik artırma — race yoxdu
}
```

**Sub-kod izahı:**
- RWMutex: çox oxucu PARALEL; yazıçı tək və oxucuları bloklayır
- `atomic.Add/Store/Load` → mutex-dən ucuz, sayğac/flag üçün ideal
- `atomic.CompareAndSwapUint64` → CAS əməliyyatı da mövcuddur
- Kilidi mütləq aç (defer Unlock) — yoxsa deadlock

### 4. context paketi
**Nədir:** Cancel/timeout/value daşıyan sorğu-konteksti — funksiyalar arası
idarəetmə kanalı.

**Kitabdan kod nümunəsi:**
```go
// Value — typed key:
type key string
const (
    timeoutKey  key = "TimeoutKey"
    deadlineKey key = "DeadlineKey"
)
func Setup(ctx context.Context) context.Context {
    ctx = context.WithValue(ctx, timeoutKey, "timeout exceeded")
    ctx = context.WithValue(ctx, deadlineKey, "deadline exceeded")
    return ctx
}
func GetValue(ctx context.Context, k key) string {
    if val, ok := ctx.Value(k).(string); ok {  // type assertion
        return val
    }
    return ""
}

// Timeout vs Deadline — yarışı:
func Exec() {
    ctx := context.Background()
    ctx = Setup(ctx)
    timeoutCtx, cancel := context.WithTimeout(ctx,
        time.Duration(rand.Intn(2))*time.Millisecond)
    defer cancel()
    deadlineCtx, cancel := context.WithDeadline(ctx,
        time.Now().Add(time.Duration(rand.Intn(2))*time.Millisecond))
    defer cancel()
    for {
        select {
        case <-timeoutCtx.Done():            // timeout bitdi
            fmt.Println(GetValue(ctx, timeoutKey))
            return
        case <-deadlineCtx.Done():            // deadline çatdı
            fmt.Println(GetValue(ctx, deadlineKey))
            return
        }
    }
}
```

**Sub-kod izahı:**
- `WithTimeout(ctx, d)` → müddət; `WithDeadline(ctx, t)` → dəqiq vaxt nöqtəsi
- `defer cancel()` → resuş sızıntısının qarşısı — context ağacından çıxarış
- Uşaq contextlər valideyn dəyərlərini İRSƏN alır
- `ctx.Done()` → select daxilində cancel/timeout gözləməsi

### 5. Channel-lərlə state idarəsi
**Nədir:** Kanal mesajı kimi struct göndərmək — bir kanal, çox əməliyyat
növü (op sahəsi ilə).

**Kitabdan kod nümunəsi:**
```go
type op string
const (
    Add      op = "add"
    Subtract    = "sub"
    Multiply    = "mult"
    Divide      = "div"
)
type WorkRequest struct {
    Operation op
    Value1   int64
    Value2   int64
}
type WorkResponse struct {
    Wr     *WorkRequest
    Result int64
    Err    error
}

// Processor — ümumi iş yönləndirici:
func Processor(ctx context.Context, in chan *WorkRequest, out chan *WorkResponse) {
    for {
        select {
        case <-ctx.Done():
            return
        case wr := <-in:
            out <- Process(wr)        // işi gör, cavabı geri göndər
        }
    }
}

func Process(wr *WorkRequest) *WorkResponse {
    resp := WorkResponse{Wr: wr}
    switch wr.Operation {
    case Add:
        resp.Result = wr.Value1 + wr.Value2
    case Divide:
        if wr.Value2 == 0 {
            resp.Err = errors.New("divide by 0")   // xəta da cavabda
            break
        }
        resp.Result = wr.Value1 / wr.Value2
    default:
        resp.Err = errors.New("unsupported operation")
    }
    return &resp
}

// Buffered kanallar + sayğaclı istifadə:
in := make(chan *state.WorkRequest, 10)
out := make(chan *state.WorkResponse, 10)
go state.Processor(ctx, in, out)
in <- &state.WorkRequest{state.Divide, 8, 0}    // xətalı sorğu
resp := <-out                                    // Err sahəsi dolu
```

**Sub-kod izahı:**
- Sorğu/cavab struct-ları → bir kanal növü ilə çox əməliyyat
- Xəta da cavabın içində (Err sahəsi) — kanal ayrılmağa ehtiyac yoxdu
- Buffered kanallar → göndərici bloklanmır

### 6. Worker pool pattern
**Nədir:** N sabit goroutine eyni in kanalını dinləyir; maksimum paralelliyi
idarə edir (crypto kimi ağır işlər üçün).

**Kitabdan kod nümunəsi:**
```go
// Dispatch — pool yaradıcı:
func Dispatch(numWorker int) (context.CancelFunc, chan WorkRequest, chan WorkResponse) {
    ctx := context.Background()
    ctx, cancel := context.WithCancel(ctx)
    in := make(chan WorkRequest, 10)
    out := make(chan WorkResponse, 10)
    for i := 0; i < numWorker; i++ {
        go Worker(ctx, i, in, out)     // N worker, eyni kanallar
    }
    return cancel, in, out
}

func Worker(ctx context.Context, id int, in chan WorkRequest, out chan WorkResponse) {
    for {
        select {
        case <-ctx.Done():
            return
        case wr := <-in:
            fmt.Printf("worker id: %d, performing %s work\n", id, wr.Op)
            out <- Process(wr)         // hər worker HƏR növ işi edə bilir
        }
    }
}

// İş tipləri (bcrypt misalı):
type WorkRequest struct {
    Op      op        // Hash / Compare
    Text    []byte
    Compare []byte    // Compare üçün
}
func hashWork(wr WorkRequest) WorkResponse {
    val, err := bcrypt.GenerateFromPassword(wr.Text, bcrypt.DefaultCost)
    return WorkResponse{Result: val, Err: err, Wr: wr}
}
func compareWork(wr WorkRequest) WorkResponse {
    var matched bool
    if err := bcrypt.CompareHashAndPassword(wr.Compare, wr.Text); err == nil {
        matched = true
    }
    return WorkResponse{Matched: matched, Err: err, Wr: wr}
}

// İstifadə — hash-la, sonra müqayisə et:
cancel, in, out := pool.Dispatch(10)
defer cancel()
for i := 0; i < 10; i++ {
    in <- pool.WorkRequest{Op: pool.Hash, Text: []byte(fmt.Sprintf("messages %d", i))}
}
for i := 0; i < 10; i++ {
    res := <-out
    in <- pool.WorkRequest{Op: pool.Compare, Text: res.Wr.Text, Compare: res.Result}
}
```

**Sub-kod izahı:**
- Eyni pool muxtəlif iş növlərini işləyir (Op sahəsi ilə dispatch)
- Worker sayı = maksimum paralellik → yaddaş/CPU qoruması
- Crypto kimi CPU-intensiv işlərdə hər sorğu üçün YENİ goroutine YARATMA
  prosesi boğa bilər — pool bunun həllidir

### 7. Worker pipeline
**Nədir:** Pool-ları zəncirləmək — bir pool-un çıxışı digərinin girişi
(konveyer).

**Kitabdan kod nümunəsi:**
```go
// Worker — rolu iş təyinatında müəyyən olunur:
type Worker struct {
    in  chan string
    out chan string
}
type Job string
const (
    Print  Job = "print"
    Encode Job = "encode"
)
func (w *Worker) Work(ctx context.Context, j Job) {
    switch j {
    case Print:
        w.Print(ctx)
    case Encode:
        w.Encode(ctx)
    }
}

func (w *Worker) Encode(ctx context.Context) {
    for {
        select {
        case <-ctx.Done():
            return
        case val := <-w.in:
            w.out <- fmt.Sprintf("%s => %s", val,
                base64.StdEncoding.EncodeToString([]byte(val)))
        }
    }
}
func (w *Worker) Print(ctx context.Context) {
    for {
        select {
        case <-ctx.Done():
            return
        case val := <-w.in:
            fmt.Println(val)
            w.out <- val              // növbəti mərhələyə ötür
        }
    }
}

// NewPipeline — mərhələləri birləşdirir:
func NewPipeline(ctx context.Context, numEncoders, numPrinters int) (chan string, chan string) {
    inEncode := make(chan string, numEncoders)
    inPrint := make(chan string, numPrinters)
    outPrint := make(chan string, numPrinters)
    for i := 0; i < numEncoders; i++ {
        w := Worker{in: inEncode, out: inPrint}   // encode → print GİRİŞİ
        go w.Work(ctx, Encode)
    }
    for i := 0; i < numPrinters; i++ {
        w := Worker{in: inPrint, out: outPrint}
        go w.Work(ctx, Print)
    }
    return inEncode, outPrint        // pipeline girişi və çıxışı
}

// 10 encoder + 2 printer:
in, out := pipeline.NewPipeline(ctx, 10, 2)
go func() {
    for i := 0; i < 20; i++ {
        in <- fmt.Sprint("Message", i)     // mərhələ 1: encode
    }
}()
for i := 0; i < 20; i++ {
    <-out                                   // 20 nəticə gözlə
}
```

**Sub-kod izahı:**
- Encode pool-un çıxışı (inPrint) = Print pool-un girişi — kanallarla birləşmə
- Fərqli ölçülü pool-lar (10 encoder, 2 printer) — baha əməliyyatları
  məhdudlaşdırmaq üçün
- Sayğaclı bitmə: 20 girişə 20 çıxış

## Əsas terminlər

- Goroutine / Channel / Buffered Channel
- select operatoru
- time.Tick (select içində sıfırlanma tələsi)
- sync.WaitGroup (Add/Done/Wait)
- Race Condition (yarış şəraiti)
- Mutex / RWMutex (Lock/RLock)
- sync/atomic (Store/Load/Add/CompareAndSwap)
- sync.Once (bir dəfəlik icra)
- context (Background/WithCancel/WithTimeout/WithDeadline/WithValue/Done)
- Worker Pool (işçi hovuzu)
- Pipeline (konveyer) / Dispatch

## Praktik nəticə

- Bitirmə üçün 2 yol: done kanalı və ya context — mürəkkəb sistemlərdə
  context üstünlük ver
- Loop-dan goroutine başladarkən dəyişəni arqument kimi ötür
- Map paralel yazma → RWMutex; sadə sayğac → atomic; bir dəfəlik init → Once
- context-də typed key işlət; hər With* üçün defer cancel
- CPU-ağır işlər (crypto/bcrypt) üçün worker pool — limitsiz goroutine YOX
- Pool-ları kanallarla zəncirlə → pipeline; hər mərhələnin ölçüsünü ayrıca
  tənzimlə

## Mənbə

Pages: 324-360 (Chapter 10, Go Programming Cookbook 2nd ed)
