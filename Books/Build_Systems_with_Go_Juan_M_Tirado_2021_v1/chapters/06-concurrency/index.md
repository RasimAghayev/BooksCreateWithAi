# Chapter 6 — Concurrency (səh. 102-142)

## Bu fəsil nədən bəhs edir?

Go-nun konkurensiya (paralel işləmə) alətləri: goroutine-lər, kanallar
(buffered/unbuffered, close, direction), select, WaitGroup, Timer/Ticker,
Context (WithCancel/WithTimeout/WithDeadline/WithValue), sync.Once, Mutex
və atomic əməliyyatları.

## Əsas fikirlər

### 1. Goroutine-lər
**Nədir:** Go runtime tərəfindən idarə olunan yüngül thread (~2KB yaddaş).
`go` açar sözü ilə funksiya paralel işə salınır.

**Kitabdan kod nümunəsi:**
```go
func ShowIt() {
    for {
        time.Sleep(time.Millisecond * 100)
        fmt.Println("Here it is!!!")
    }
}

func main() {
    go ShowIt()          // paralel işə salınır
    for i := 0; i < 5; i++ {
        time.Sleep(time.Millisecond * 50)
        fmt.Println("Where is it?")
    }
}
```

**Vacib qaydalar:**
- `main` bitəndə proqram dayanır — digər goroutine-lər işləməyə davam
  etmir (sonsuz loop olan ShowIt belə)
- Birdən çox goroutine asanlıqla yaradılır: `go ShowIt(t1, 100)`,
  `go ShowIt(t2, 10)`...
- Çoxlu goroutine eyni çıxışa (console) yazanda nəticə qeyri-deterministikdir

### 2. Kanallar (Channels)
**Nədir:** Paralel işləyən funksiyalar arasında kommunikasiya mexanizmi.
`make(chan T)` ilə yaradılır; göndər `ch <- v`, qəbul `v := <- ch`.

**Kitabdan kod nümunəsi (əsas kanal axını):**
```go
func generator(ch chan int) {
    sum := 0
    for i := 0; i < 5; i++ {
        time.Sleep(time.Millisecond * 500)
        sum = sum + i
    }
    ch <- sum          // nəticəni göndər
}

func main() {
    ch := make(chan int)
    go generator(ch)
    fmt.Println("main waits for result…")
    result := <- ch    // BLOKLANIR — göndəriş gələnə qədər
    fmt.Println(result)
}
```

**Oxuma/yazma sintaksisi:**
```go
ch <- 5        // 5-i kanala göndər
n := <- ch     // kanaldan oxu
<- ch          // göndəriş gələnə qədər gözlə
```

### 3. Buffered kanallar
**Nədir:** `make(chan int)` — unbuffered (hər iki tərəf hazır olmalıdır);
`make(chan int, 10)` — 10 elementlik buffer; göndərici qəbuledici
hazır olmasa da bufferə yaza bilər.

Unbuffered kanala hazır olmayan tərəf yazsa: `fatal error: all goroutines
are asleep - deadlock!`.

```go
func main() {
    ch := make(chan string, 1)   // 1-element buffer
    ch <- "This is main"          // buffer sayəsində deadlock yoxdur
    go MrA(ch)
    go MrB(ch)
    fmt.Println(<-ch)
    fmt.Println(<-ch)
    fmt.Println(<-ch)
}
```

**len/cap qeydi:** `len(ch)==cap(ch)` yoxlaması etibarsızdır — yoxlama
sonrası vəziyyət dəyişə bilər; `select` istifadə edin.

### 4. close və "comma ok"
**close(ch)** → kanalı bağlayır; bağlı kanala yazma → panic.
Qəbuledici bağlanmanı belə tutur:

```go
func sender(out chan int) {
    for i := 0; i < 5; i++ {
        time.Sleep(time.Millisecond * 500)
        out <- i
    }
    close(out)
}

func main() {
    ch := make(chan int)
    go sender(ch)
    for {
        num, found := <- ch     // found=false → bağlanıb
        if found {
            fmt.Println(num)
        } else {
            fmt.Println("finished")
            break
        }
    }
}
```

### 5. range ilə kanal istehlakı
Kanalın bağlanmasına qədər dəyərləri avtomatik gözləyir — element sayını
 əvvəlcədən bilmək lazım deyil:

```go
func generator(ch chan int) {
    for i := 0; i < 5; i++ {
        ch <- i
    }
    close(ch)
}

func main() {
    ch := make(chan int)
    go generator(ch)
    for x := range ch {    // close olunana qədər çap edir
        fmt.Println(x)
    }
    fmt.Println("Done")
}
```

### 6. Kanal istiqaməti (direction)
**Nədir:** Kanalın tipində axın istiqaməti göstərilir — type-safety artırır;
 səhv istiqamət **compile xətası** verir.

```go
ch := make(chan int)       // hər ikisi
ch := make(<-chan int)     // yalnız qəbuledici (receive-only)
ch := make(chan<- int)     // yalnız göndərici (send-only)

func receiver(input <-chan int) {   // yalnız oxuya bilər
    for i := range input {
        fmt.Println(i)
    }
}

func sender(output chan<- int, n int) {  // yalnız yaza bilər
    for i := 0; i < n; i++ {
        time.Sleep(time.Millisecond * 500)
        output <- i * i
    }
    close(output)
}
```

### 7. select
**Nədir:** Switch-ə bənzəyir; bir neçə kanal əməliyyatından **birincisi
hazır olana qədər bloklanır**. Bir neçə hazır olsa pseudo-random seçilir.
select bir dəfə icra olunur — davamlı gözləmək üçün loop-a alınır.

```go
for i := 0; i < 10; i++ {
    select {
    case num := <-numbers:
        fmt.Printf("number %d\n", num)
    case msg := <-msgs:
        fmt.Printf("msg %s\n", msg)
    }
}
```

**"Comma ok" ilə bağlanma izləmə:**
```go
closedNums, closedMsgs := false, false
for !closedNums || !closedMsgs {
    select {
    case num, ok := <-numbers:
        if ok { fmt.Printf("number %d\n", num) } else { closedNums = true }
    case msg, ok := <-msgs:
        if ok { fmt.Printf("msg %s\n", msg) } else { closedMsgs = true }
    }
}
```

**default ilə non-blocking select:**
```go
select {
case i := <-ch:
    fmt.Println("Received", i)
default:
    fmt.Println("Nothing received")   // kanal hazır deyilsə panicsiz keç
}
```

### 8. WaitGroup (sync paketi)
**Nədir:** Goroutine-lərin tamamlanmasını gözləmək üçün sayğac: `Add(n)`
artırır, `Done()` azaldır, `Wait()` sıfırlanana qədər bloklanır.

**Kitabdan kod nümunəsi (producer/consumer):**
```go
func generator(ch chan int, wg *sync.WaitGroup) {
    defer wg.Done()               // zəmanətli bitmə bildirişi
    for i := 0; i < 5; i++ {
        time.Sleep(time.Millisecond * 200)
        ch <- rand.Int()
    }
    close(ch)
    fmt.Println("Generator done")
}

func consumer(id int, sleep time.Duration, ch chan int, wg *sync.WaitGroup) {
    defer wg.Done()
    for task := range ch {
        time.Sleep(time.Millisecond * sleep)
        fmt.Printf("%d - task[%d]\n", id, task)
    }
    fmt.Printf("Consumer %d done\n", id)
}

func main() {
    rand.Seed(42)
    ch := make(chan int, 10)
    var wg sync.WaitGroup
    wg.Add(3)                      // 3 goroutine gözlənilir

    go generator(ch, &wg)          // wg referansla ötürülür
    go consumer(1, 400, ch, &wg)
    go consumer(2, 100, ch, &wg)

    wg.Wait()                       // hamısı bitənə qədər gözlə
}
```

### 9. Timer və Ticker
- **Timer** → gələcəkdə bir dəfə siqnal verir (`time.NewTimer(5*time.Second)`)
- **Ticker** → müəyyən periodla təkrar siqnal (`time.NewTicker(time.Second)`)
- Hər ikisinin `C` kanalı var; `select` ilə gözlənilir

```go
timer := time.NewTimer(time.Second * 5)
ticker := time.NewTicker(time.Second)
x := 0
go worker(&x)

for {
    select {
    case <-timer.C:
        fmt.Printf("timer -> %d\n", x)
    case <-ticker.C:
        fmt.Printf("ticker -> %d\n", x)
    }
    if x >= 10 { break }
}
```

**Stop idarəetməsi:**
```go
quick := time.NewTicker(time.Second)
slow := time.NewTimer(time.Second * 5)
stopper := time.NewTimer(time.Second * 4)
go reaction(quick)
go slowReaction(slow)

<-stopper.C
quick.Stop()                     // ticker dayandırılır (dəyər qaytarmır)
stopped := slow.Stop()           // true → hadisə baş verməmiş dayandırıldı
fmt.Println("Stopped before the event?", stopped)
```

### 10. Context
**Nədir:** Ləğv etmə (cancellation), timeout, deadline və request-scoped
dəyərləri ötürmək üçün standart mexanizm (API çağırışları üçün əsas).

**Context interfeysi:**
```go
type Context interface {
    Deadline() (deadline time.Time, ok bool)
    Done() <-chan struct{}    // bağlananda siqnal verir
    Err() error               // səbəb (deadline exceeded, canceled...)
    Value(key interface{}) interface{}
}
```

**WithCancel — əl ilə ləğv:**
```go
func setter(id int, c *int32, ctx context.Context) {
    t := time.NewTicker(time.Millisecond * 300)
    for {
        select {
        case <-ctx.Done():
            fmt.Println("Done", id)
            return
        case <-t.C:
            atomic.AddInt32(c, 1)
        }
    }
}

func main() {
    ctx, cancel := context.WithCancel(context.Background())
    var c int32 = 0
    for i := 0; i < 5; i++ {
        go setter(i, &c, ctx)
    }
    time.Sleep(time.Second)
    cancel()          // bütün goroutine-lər Done siqnalı alır
    time.Sleep(time.Second)
}
```

**WithTimeout — müddət bitincə ləğv:**
```go
d := time.Millisecond * 300
ch := make(chan int)
for i := 0; ; i++ {
    ctx, cancel := context.WithTimeout(context.Background(), d)
    go work(i, ch)
    select {
    case x := <-ch:
        fmt.Println("Received", x)
    case <-ctx.Done():
        fmt.Println("Done!!")
    }
    if ctx.Err() != nil {
        fmt.Println(ctx.Err())   // "context deadline exceeded"
        return
    }
    cancel()   // loop daxilində defer EDİLMƏZ!
    _ = i
}
```
**Vacib:** loop daxilində yaranan kontekstlərin `cancel`-i defer
edilməz — funksiya qayıdana qədər resurslar saxlanılır; hər iterasiyada
 açıq çağırın.

**WithDeadline — mütləq vaxt:**
```go
d := time.Now().Add(time.Second * 3)
ctx, cancel := context.WithDeadline(context.Background(), d)
defer cancel()
// ...
<-ctx.Done()   // 3 saniyə sonra avtomatik bağlanır
```

**WithValue — kontekstə dəyər daşıma:**
```go
f := func(ctx context.Context, a int, b int) (int, error) {
    switch ctx.Value("action") {
    case "+": return a + b, nil
    case "-": return a - b, nil
    default:  return 0, errors.New("unknown action")
    }
}
ctx := context.WithValue(context.Background(), "action", "+")
r, err := f(ctx, 22, 20)   // 42 <nil>
```

**Parent contexts — məhdudiyyətlərin birgə istifadəsi:**
```go
t := time.Millisecond * 300
ctx, cancel := context.WithTimeout(context.Background(), t)
qCtx := context.WithValue(ctx, "action", "quick")  // timeout + value
defer cancel()
go calc(qCtx)
<-qCtx.Done()
```
- Value konteksti timeout kontekstindən törəyir — hər iki məhdudiyyət
  eyni anda işləyir ("slow" əməliyyatı timeout-a düşür, çap olunmur).

### 11. sync.Once
**Nədir:** Əməliyyatın yalnız bir dəfə icrasını zəmanətləndirir —
initializasiya üçün idealdır. `Do(fn)` funksiyanı ilk çağırışda icra
edir, qalanları sakitcə keçir.

```go
var first int

func setter(i int, ch chan bool, once *sync.Once) {
    t := rand.Uint32() % 300
    time.Sleep(time.Duration(t) * time.Millisecond)
    once.Do(func() {     // yalnız 1 goroutine bunu icra edir
        first = i
    })
    ch <- true
    fmt.Println(i, "Done")
}

func main() {
    rand.Seed(time.Now().UnixNano())
    var once sync.Once
    ch := make(chan bool)
    for i := 0; i < 10; i++ {
        go setter(i, ch, &once)
    }
    for i := 0; i < 10; i++ {
        <-ch
    }
    fmt.Println("The first was", first)
}
```

### 12. Mutex
**Nədir:** Race condition (yarış şəraiti) — bir neçə goroutine eyni dəyişənə
 paralel çıxanda. `sync.Mutex` + `Lock()`/`Unlock()` kritik sahəni qoruyur.

**Kitabdan kod nümunəsi:**
```go
func writer(x map[int]int, factor int, m *sync.Mutex) {
    i := 1
    for {
        time.Sleep(time.Second)
        m.Lock()                    // kritik sahə başlayır
        x[i] = x[i-1] * factor
        m.Unlock()                   // kritik sahə bitir
        i++
    }
}

func reader(x map[int]int, m *sync.Mutex) {
    for {
        time.Sleep(time.Millisecond * 500)
        m.Lock()
        fmt.Println(x)
        m.Unlock()
    }
}

func main() {
    x := make(map[int]int)
    x[0] = 1
    m := sync.Mutex{}
    go writer(x, 2, &m)
    go reader(x, &m)
    time.Sleep(time.Millisecond * 300)
    go writer(x, 3, &m)
    time.Sleep(time.Second * 4)
}
```
- **Qayda:** exclusion sahəsi mümkün qədər kiçik olmalıdır — böyük sahə
  uzun gözləmə deməkdir.

### 13. Atomics (sync/atomic)
**Nədir:** Aşağı səviyyəli atomik yaddaş primitivləri — düzgün
işlədildikdə mutex-dən sürətlidir. Yalnız native tiplər üçün (Table 6.1,
int32 üçün; uint32/int64/uint64/uintptr üçün də analoqlar var):

| Funksiya | İşlev |
|---|---|
| `AddInt32(addr, delta)` | atomik əlavə et |
| `CompareAndSwapInt32(addr, old, new)` | müqayisə+əvəz |
| `LoadInt32(addr)` | atomik oxu |
| `StoreInt32(addr, val)` | atomik yaz |
| `SwapInt32(addr, new)` | əvəz et, köhnəni qaytar |

**Sadə sayğac:**
```go
func increaser(counter *int32) {
    for {
        atomic.AddInt32(counter, 2)
        time.Sleep(time.Millisecond * 500)
    }
}
func decreaser(counter *int32) {
    for {
        atomic.AddInt32(counter, -1)
        time.Sleep(time.Millisecond * 250)
    }
}
// oxumaq: atomic.LoadInt32(&counter)
```

**atomic.Value — ixtiyari tip üçün:**
```go
type Monitor struct {
    ActiveUsers int
    Requests     int
}

func updater(monitor atomic.Value, m *sync.Mutex) {
    for {
        time.Sleep(time.Millisecond * 500)
        m.Lock()
        current := monitor.Load().(*Monitor)   // interface{} qayıdır → cast
        current.ActiveUsers += 100
        current.Requests += 300
        monitor.Store(current)
        m.Unlock()    // sahə dəyişikliyi atomik DEYİL → mutex lazımdır
    }
}
func observe(monitor atomic.Value) {
    for {
        time.Sleep(time.Second)
        fmt.Printf("%v\n", monitor.Load())
    }
}
```

## Əsas terminlər
- Goroutine — Go runtime idarə etdiyi yüngül thread
- Channel (kanal) — goroutine-lər arası kommunikasiya kanalı
- Unbuffered/Buffered channel — buferli/bufersiz kanal
- Deadlock — bütün goroutine-lərin bloklanması
- select — çoxkanallı gözləmə bloku
- WaitGroup — goroutine tamamlanma sayğacı
- Timer/Ticker — tək/təkrar vaxt siqnalları
- Context — ləğv/timeout/deadline/dəyər daşıyıcısı
- Race condition — paralel yazış yarışı
- Mutex — qarşılıqlı istisna kilidi
- Atomic — bölünməz yaddaş əməliyyatı

## Praktik nəticə
Go-nun konkurensiya fəlsəfəsi: "kanallar vasitəsilə kommunikasiya, yox
 yaddaş paylaşımı". Vaxt məhdudiyyəti üçün `select` + Timer/Ticker/Context;
 bir dəfəlik əməliyyat üçün `sync.Once`; paylaşılan dəyişən üçün Mutex
 (kiçik kritik sahə!) və ya atomik. API serverlərdə hər sorğunun öz
 Context-i olmalıdır (WithTimeout + defer cancel). Loop daxilində cancel-i
 açıq çağırın, defer etməyin.

## Mənbə
Pages: 102-142 (PDF 102-142)
