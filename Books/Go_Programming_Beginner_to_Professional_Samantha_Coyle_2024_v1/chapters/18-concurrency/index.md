# Chapter 18 — Concurrent Work (səh. 562-607)

## Bu fəsil nədən bəhs edir?

Konkurensiya: goroutine-lər (`go` açar sözü), WaitGroup, race
condition + atomic, Mutex, kanallar (buffer/unbuffered, close, range,
bidirectional), pipeline/fan-out pattern-ləri, worker struct + metod
goroutine kimi, context (WithCancel + Done), sync.Cond və sync.Map.

## Əsas fikirlər

### 1. Goroutine-lər
**Nədir:** eyni prosesdə konkurrent icra olunan yüngül tapşırıqlar
("çəkic paylaşan dəmirçilər" metaforası — hamısı eyni çəkiclə, növbə ilə).

```go
func hello() {
    fmt.Println("hello world")
}

func main() {
    fmt.Println("Start")
    go hello()          // İCRA GÖZLƏMİR — main davam edir!
    fmt.Println("End")
}
```

**Problem:** main bitəndə başqa goroutine-lər MÖVCUD DEYİL:
```go
var s1 int
go func() {
    s1 = sum(1, 100)    // 5050 hesablanır...
}()
fmt.Println(s1)         // 0! — hələ bitməyib (və ya main artıq bitib)

// Kobud həll (test xaricində MƏNASIZDIR):
time.Sleep(time.Second)
fmt.Println(s1, s2)     // 5050 55
```
- time.Sleep konkurensiyanın MƏQSƏDİNİ (sürət) məhv edir — düzgün
  həll aşağıda: WaitGroup

### 2. WaitGroup
```go
func sum(from, to int, wg *sync.WaitGroup, res *int) {
    *res = 0
    for i := from; i <= to; i++ {
        *res += i
    }
    wg.Done()                     // bitdi bildir
}

func main() {
    s1 := 0
    wg := &sync.WaitGroup{}
    wg.Add(1)                     // 1 goroutine gözlənilir
    go sum(1, 100, wg, &s1)
    wg.Wait()                     // hamısı Done olana qədər BLOK
    log.Println(s1)               // 5050
}
```
- Add/Done/Wait üçlüyü; nəticə pointer ilə qaytar

### 3. Race condition + atomic
```go
// YARIS ŞƏRAİTİ:
func next(v *int) {
    c := *v        // hamısı eyni anda 0 oxuya bilər!
    *v = c + 1     // bir-birinin dəyişikliyi İTİR
}
a := 0
go next(&a); go next(&a); go next(&a)
// a = 1, 2 və ya 3 — QƏRARSIZ!
```

**atomic həlli:**
```go
func sum(from, to int, wg *sync.WaitGroup, res *int32) {
    for i := from; i <= to; i++ {
        atomic.AddInt32(res, int32(i))    // BÖLÜNMƏZ artım
    }
    wg.Done()
}

func main() {
    s1 := int32(0)
    wg := &sync.WaitGroup{}
    wg.Add(4)
    go sum(1, 25, wg, &s1)
    go sum(26, 50, wg, &s1)
    go sum(51, 75, wg, &s1)
    go sum(76, 100, wg, &s1)
    wg.Wait()
    log.Println(s1)         // 5050 — həmişə!
}
```
- atomic yalnız int32/int64/uint tipləri üçün
- **-race flag:** `go test -race` — yarışı TUTUR (atomic-sız versiyada
  xəta verir; atomic-lı versiyada PASS)

### 4. Mutex
```go
mtx := &sync.Mutex{}         // pointer — hamısı EYNİ mutex-i paylaşır

mtx.Lock()                    // kritik sahə başlayır
s = s + 5                     // yalnız BİR goroutine daxil ola bilər
mtx.Unlock()                  // azad et
```
- Aradakı kod MİNİMAL olmalıdır — nə qədər az, o qədər konkurrent
- sıra zəmanəti YOXDUR — nəticələri sonradan sırala

### 5. Kanallar
```go
ch := make(chan int)          // unbuffered (buffer = 0)
ch := make(chan int, 10)      // 10-luq buffer

ch <- 2                        // göndər (doludursa — BLOK)
i := <- ch                     // qəbul (boşdursa — BLOK)
<- ch                          // qəbul, at
close(ch)                      // bağla (göndəriş SON, qəbul HƏLƏ MÜMKÜN)
defer close(ch)                 // funksiya sonunda bağla
```

**Deadlock — unbuffered + tək goroutine:**
```go
ch := make(chan int, 1)   // buffered 1 → işləyir
ch <- 1
i := <- ch

ch := make(chan int)      // unbuffered → DEADLOCK!
ch <- 1                  // qəbul edən YOX — əbədi blok
// fatal error: all goroutines are asleep - deadlock!
```
- Unbuffered kanal = iki tərəf hazırsa işləyir

**Buffer dolu + extra göndəriş → deadlock;** buffer artır → düzəlir.

### 6. Kanal kommunikasiya nümunələri
**Bir istiqamətli greeting:**
```go
func greet(ch chan string) {
    ch <- "Hello"
}

func main() {
    ch := make(chan string)
    go greet(ch)
    log.Println(<-ch)      // Hello
}
```

**İkitərəfli:**
```go
func greet(ch chan string) {
    msg := <- ch
    ch <- fmt.Sprintf("Thanks for %s", msg)
    ch <- "Hello David"
}

func main() {
    ch := make(chan string)
    go greet(ch)
    ch <- "Hello John"      // göndər
    log.Println(<-ch)        // Thanks for Hello John
    log.Println(<-ch)        // Hello David
}
```

**Sorğu-cavab (request pattern):**
```go
func push(from, to int, in chan bool, out chan int) {
    for i := from; i <= to; i++ {
        <- in              // İSTƏK gözlə
        out <- i           // sonra göndər
    }
}
// main: in <- true; i := <- out — pull modeli
```

### 7. Worker pool + range + close
```go
func worker(in chan int, out chan int) {
    sum := 0
    for i := range in {      // close-edilənə qədər oxu!
        sum += i
    }
    out <- sum                // qismi cəmi qaytar
}

func sum(workers, from, to int) int {
    out := make(chan int, workers)
    in := make(chan int, 4)
    for i := 0; i < workers; i++ {
        go worker(in, out)            // N worker başlat
    }
    for i := from; i <= to; i++ {
        in <- i                      // işi payla
    }
    close(in)                        // "İŞ BİTDİ" siqnalı

    sum := 0
    for i := 0; i < workers; i++ {
        sum += <-out                  // qismi cəmləri topla
    }
    close(out)
    return sum
}

res := sum(100, 1, 100)     // 5050
```
- **range kanal** — close olana qədər avtomatik iterasiya
- close-dan SONRA da qəbul MÜMKÜN; yalnız GÖNDƏRİŞ qadağan

**Done channel (WaitGroup əvəzi):**
```go
func readThem(in, out chan string) {
    for i := range in {
        log.Println(i)
    }
    out <- "done"                    // bitdi bildir
}

func main() {
    in, out := make(chan string), make(chan string)
    go readThem(in, out)
    for _, s := range []string{"a","b","c","d","e","f"} {
        in <- s
    }
    close(in)
    <-out                             // done gözlə
}
```

### 8. Patternlər
- **Pipeline** — mənbə → mərhələlər → sink
- **Fan-out/Fan-in** — bir kanaldan N oxuyucu; nəticələr bir nöqtədə
  birləşir
- **Funksiya kanal QAYTARIR** (goroutine-i özü başladır):
```go
func doSomething() chan int {
    ch := make(chan int)
    go func() {
        for i := range ch {
            log.Println(i)
        }
    }()
    return ch
}
ch := doSomething()    // go ÇAĞIRMA — daxildə başlayır
ch <- 1
ch <- 4
```
- **İstiqamət annotasiyası:** `<-chan int` (yalnız oxu), `chan<- int`
  (yalnız yaz)

### 9. Metodlar goroutine kimi + Worker struct
```go
type Worker struct {
    in, out chan int
    sbw int                 // subworker sayı
    mtx *sync.Mutex
}

func (w *Worker) readThem() {
    w.sbw++
    go func() {
        partial := 0
        for i := range w.in {
            partial += i
        }
        w.out <- partial

        w.mtx.Lock()          // sayğacı təhlükəsiz azalt
        w.sbw--
        if w.sbw == 0 {
            close(w.out)       // son subworker out-u bağlayır
        }
        w.mtx.Unlock()
    }()
}

func (w *Worker) gatherResult() int {
    total := 0
    wg := &sync.WaitGroup{}
    wg.Add(1)
    go func() {
        for i := range w.out {
            total += i
        }
        wg.Done()
    }()
    wg.Wait()
    return total
}

// main: 10 subworker + 100 rəqəm + close(in) + gatherResult → 5050
```

### 10. Context (WithCancel)
```go
func countNumbers(ctx context.Context, r chan int) {
    v := 0
    for {
        select {
        case <-ctx.Done():        // ləğv siqnalı
            r <- v
            return
        default:
            time.Sleep(time.Millisecond * 100)
            v++
        }
    }
}

func main() {
    r := make(chan int)
    ctx := context.TODO()

    cl, stop := context.WithCancel(ctx)   // ləğv edilə bilən kontekst
    go countNumbers(cl, r)

    go func() {
        time.Sleep(time.Millisecond * 100 * 3)
        stop()                           // 300ms sonra DAYANDIR
    }()

    v := <- r
    log.Println(v)     // 3 — sonsuz loop 3 iterasiyadan sonra dayandı
}
```
- ctx zəncir boyu ötürülür; Done() kanalı ləğvdə bağlanır
- Göndərən başqa goroutine-də close etməkdən qaçın (izləmək çətin)

### 11. sync.Cond (şərt dəyişəni)
```go
type WorkQueue struct {
    cond      *sync.Cond
    maxSize   int
    workItems []string
}

func NewWorkQueue(maxSize int) *WorkQueue {
    return &WorkQueue{
        cond:      sync.NewCond(&sync.Mutex{}),
        maxSize:   maxSize,
        workItems: make([]string, 0),
    }
}

func (wq *WorkQueue) enqueue(item string) {
    wq.cond.L.Lock()
    defer wq.cond.L.Unlock()
    for len(wq.workItems) == wq.maxSize {   // DOLUDUR → gözlə
        wq.cond.Wait()
    }
    wq.workItems = append(wq.workItems, item)
    wq.cond.Signal()                        // gözləyənlərə xəbər ver
}

func (wq *WorkQueue) dequeue() string {
    wq.cond.L.Lock()
    defer wq.cond.L.Unlock()
    for len(wq.workItems) == 0 {            // BOŞDUR → gözlə
        wq.cond.Wait()
    }
    item := wq.workItems[0]
    wq.workItems = wq.workItems[1:]
    wq.cond.Signal()
    return item
}
```
- Wait gözləyən goroutine-i yuxuya salır; Signal oyadır
- Mutex ilə birlikdə — WIP-limitli növbə üçün ideal

### 12. sync.Map (thread-safe map)
```go
func generateRandomNumber(max int) (int, error) {
    n, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
    if err != nil {
        return 0, err
    }
    return int(n.Int64()), nil
}

func updateCount(countMap *sync.Map, key int) {
    count, _ := countMap.LoadOrStore(key, 0)   // yoxdursa 0 qoy
    countMap.Store(key, count.(int)+1)
}

func printCounts(countMap *sync.Map) {
    countMap.Range(func(key, value interface{}) bool {
        fmt.Printf("Number %d: Count %d\n", key, value)
        return true                               // davam et
    })
}

func main() {
    var countMap sync.Map             // mutex LAZIM DEYİL!
    var wg sync.WaitGroup
    // 5 goroutine × 1000 random rəqəm sayı
}
```
- LoadOrStore/Store/Range/Load/Delete — daxili kilidlərlə
- Çox oxuyucu + bir yazıçı ssenariləri üçün optimallaşdırılmış

## Əsas terminlər
- Goroutine / `go` açar sözü
- Concurrency vs parallelism
- sync.WaitGroup — Add/Done/Wait
- Race condition — üst-üstə yazma yarışı
- atomic.AddInt32 — bölünməz əməliyyat
- `-race` flag — yarış detektoru
- Mutex — Lock/Unlock kritik sahə
- Channel — tipik kommunikasiya kanalı
- Buffered/unbuffered channel + deadlock
- close(ch) — göndəriş sonu; qəbul hələ mümkün
- range kanal — close-a qədər iterasiya
- Done channel — bitmə bildirişi
- Pipeline / fan-out/fan-in / worker pool
- Kanal istiqaməti — `<-chan` / `chan<-`
- Metod = goroutine (struct daxili kanal inkapsulyasiyası)
- context.WithCancel + ctx.Done() + stop()
- sync.Cond — Wait/Signal şərt primitivi
- sync.Map — LoadOrStore/Range konkurrent map
- "Share by communicating, not communicate by sharing"

## Praktik nəticə
Konkurensiya yalnız MÜSTƏQİL işlər üçün (web server, çoxmənbəli data
yığılımı) — yoxsa faydasızdır. Əsas alət dəsti: go + WaitGroup (bitmə
gözləmə); paylaşılan rəqəm → atomic; mürəkkəb vəziyyət → Mutex (qısa
kritik sahə!); data ötürmə → kanal (unbuffered senxron, buffered
dekuplə; range+close worker pattern). Worker-ləri struct + metod
goroutine-lərinə qur; bitməsi done kanalı/son subworker close ilə.
Ləğvetmə üçün context.WithCancel + select{Done()}. Mutex+Cond =
bloklı növbələr; sync.Map = kilidsiz konkurrent map. Testləri HƏMİŞƏ
`-race` ilə işə sal.

## Mənbə
Pages: 562-607 (PDF 562-607)
