# Chapter 16 — Concurrent Work (Paralel İş)

## Bu fəsil nədən bəhs edir?

Goroutine-lər (go açar sözü, main rutin ilə birlikdə), WaitGroup (Add/Done/Wait),
race condition-lər (paylaşılan dəyişən təhlükəsi, `-race` flagi), atomic əməliyyatlar
(sync/atomic, AddInt32), Mutex (Lock/Unlock), kanallar (buffered/unbuffered, deadlock,
close, range, yönlü kanallar, worker pool, done channel), HTTP server-də state,
metod-goroutine-lər (struct Worker) və context paketi (WithCancel, Done).

## Əsas fikirlər

### 1. Concurrency vs Parallelism
- **Concurrent (paralel strukturlu):** tapşırıqlar KİÇİK HİSSƏLƏRƏ bölünür, maşın hər
  tapşırıqdan bir hissə icra edir — hamısı "eyni anda başlayır" (hammer nümunəsi:
  bir çəkic, hamı növbə ilə istifadə edir)
- **Parallel:** çoxnüvəli maşında həqiqətən eyni anda
- Proqramçı baxımından fərq VACİB DEYİL — tapşırıqlar müstəqil yazılır

### 2. Goroutine — go Açar Sözü
```go
func hello() { fmt.Println("hello world") }
go hello()        // asan bu qədər — main onu GÖZLƏMİR

func main() {
    fmt.Println("Start")
    go hello()      // İcra buraya TULLANIR, davam edir
    fmt.Println("End")   // End əvvəl çap oluna bilər!
}
```
**Hər proqram ən azı 1 goroutine-da işləyir (main).**

**Problem:** main çıxsa, goroutine-lər ÖLÜR:
```go
go func() { s1 = sum(1, 100) }()
fmt.Println(s1)     // 0! — hesablamadan əvvəl çap edir
```
**Kobud həll:** `time.Sleep(time.Second)` — production-da YOX.

### 3. WaitGroup — Add/Done/Wait
**Kitabdan kod nümunəsi:**
```go
func sum(from, to int, wg *sync.WaitGroup, res *int) {
    for i := from; i <= to; i++ {
        *res += i
    }
    wg.Done()                     // BITDİM bildirişi (goroutine İÇİNDƏ)
}

wg := &sync.WaitGroup{}
wg.Add(1)                         // 1 goroutine gözləniləcək (goroutine-DAN ƏVVƏL)
go sum(1, 100, wg, &s1)
wg.Wait()                         // sayğaç 0 olana qədər BLOKLA
log.Println(s1)                    // 5050 — dəqiq
```

### 4. Race Condition — Paylaşılan Dəyişən Təhlükəsi
```go
func next(v *int) {
    c := *v        // hər 3 goroutine *v=0 oxuya bilər!
    *v = c + 1     // hamısı 1 yazır → gözlənilən 3 əvəzinə 1!
}
go next(&a); go next(&a); go next(&a)   // a = 1, 2 və ya 3 — RANDOM
```
**Səbəb:** read-modify-write atomik DEYİL; goroutine-lər bir-birini MİRAR bildirmir.
**Tespit:** `go test -race` — race detector xətaları aşkarlayır (10000 iterasiya test
etmedən görsənməyə bilər).

### 5. Atomic Operations — sync/atomic
```go
func sum(from, to int, wg *sync.WaitGroup, res *int32) {
    for i := from; i <= to; i++ {
        atomic.AddInt32(res, int32(i))     // ATOMİK artırma — race YOX
    }
    wg.Done()
}

s1 := int32(0)                            // int32/uint32/int64 TƏLƏB olunur
wg.Add(4)                                  // 4 goroutine
go sum(1, 25, wg, &s1)
go sum(26, 50, wg, &s1)
go sum(51, 75, wg, &s1)
go sum(76, 100, wg, &s1)
wg.Wait()                                 // həmişə 5050, `-race` təmiz
```
Sadə riyazi əməliyyatlar üçün; mürəkkəb strukturlar üçün Mutex.

### 6. Mutex — Kilidləmə
```go
mtx := &sync.Mutex{}         // pointer — hamı EYNI mutex-u istifadə etsin

mtx.Lock()                   // başqa goroutine-lər DAYANIR
s = s + 5                     // TƏK goroutine dəyişir
mtx.Unlock()                   // burax — növbəti keçir
```
**Qaydalar:**
- Lock ilə Unlock arası MİNİMAL saxla — nə qədər çox kod, o qədər AZ concurrency
- Kilidlənən yalnız təhlükəsizlik tələb edən PAYLAŞILAN hissə
- Goroutine sırası RANDOM — sıralama lazımdırsa sonra sortla

### 7. Channels — Message Passing
```go
ch := make(chan int)        // unbuffered — 0 tutum
ch := make(chan int, 10)    // buffered — 10 slot

ch <- 2                     // GÖNDƏR (buffer doludursa BLOKLA)
i := <-ch                    // QƏBUL ET (boşdursa BLOKLA)
close(ch)                   // bağla — range-lər bitər
```

**Unbuffered tələsi:**
```go
ch := make(chan int)        // TUTUM 0
ch <- 1                     // TƏK rutində → DEADLOCK (qəbul edən YOX!)
// fatal error: all goroutines are asleep - deadlock!
```
Buffer-li (2+) eyni kodda işləyir; unbuffered ƏN AZ 2 rutin tələb edir.

**İki yönlü mesaj:**
```go
func greet(ch chan string) {
    msg := <-ch                          // QƏBUL
    ch <- fmt.Sprintf("Thanks for %s", msg)  // CAVAB
    ch <- "Hello David"
}
ch := make(chan string)
go greet(ch)
ch <- "Hello John"                       // göndər
log.Println(<-ch)                        // "Thanks for Hello John"
log.Println(<-ch)                        // "Hello David"
```

### 8. close + range — Kanalın Sonu
```go
// Göndərən tərəf:
for _, s := range strs { in <- s }
close(in)            // "DAHA MESAJ YOXDUR" siqnalı

// Qəbul edən:
for i := range in {   // close-dan sonra avtomatik DAYANIR
    fmt.Println(i)
}
```
**Worker pool (kitabdan):**
```go
func worker(in chan int, out chan int) {
    sum := 0
    for i := range in {     // kanal bağlanana qədər
        sum += i
    }
    out <- sum               // qismi cəmi qaytar
}
func sum(workers, from, to int) int {
    out := make(chan int, workers)
    in := make(chan int, 4)
    for i := 0; i < workers; i++ {
        go worker(in, out)    // N worker — EYNİ in kanalından oxuyur
    }
    for i := from; i <= to; i++ {
        in <- i                // iş bölünür
    }
    close(in)                  // hamıya BİT bildir
    sum := 0
    for i := 0; i < workers; i++ {
        sum += <-out           // qismi cəmləri topla
    }
    close(out)
    return sum                  // sum(100, 1, 100) = 5050
}
```
**Fan-out/fan-in:** bir mənbədən N worker-ə (fan-out), nəticələr tək yerdə (fan-in);
source → workers → sink pipeline pattern-i.

### 9. Yönlü Kanallar — Funksiya İmzalarında
```go
func push(from, to int, in chan bool, out chan int) {
    for i := from; i <= to; i++ {
        <-in          // İSTƏK gözlə (trigger)
        out <- i       // soruşulanda göndər
    }
}
// main:
in <- true            // sorğu göndər
i := <-out            // cavabı qəbul et

// Yönlü tiplər — özünüsənədləşdirmə:
<-chan int             // yalnız OXU
chan<- int             // yalnız YAZ
```

### 10. Done Channel — Bitmə Bildirişi
```go
func readThem(in, out chan string) {
    for i := range in { log.Println(i) }
    out <- "done"                     // İŞ BİTDİ siqnalı
}
// main: bütün mesajları göndər → close(in) → <-out (done gözlə)
```
WaitGroup-suz sinxronizasiya — channel-lə özü.

### 11. Kanal Qaytaran Funksiya
```go
func doSomething() chan int {
    ch := make(chan int)
    go func() { ... }()      // ÖZÜ goroutine başladır
    return ch
}
ch := doSomething()         // go YAZMAGA EHTİYAC YOX
ch <- 1
```

### 12. HTTP Server + Concurrency
**Kəşf:** hər HTTP request ƏLAVƏ GOROUTINE-da işlənir (avtomatik!). Request-lər
müstəqildir — kanal paylaşımı lazım deyil; amma STATE (sayğac) race yaradır →
Mutex ilə qoru. Stateful server-lər (chat, oyun) bu fəslin texnikalarını tələb edir.

### 13. Struct Worker — Metod-goroutine
**Kitabdan kod nümunəsi:**
```go
type Worker struct {
    in, out chan int
    sbw int               // sub-worker sayı
    mtx *sync.Mutex
}

func (w *Worker) readThem() {
    w.sbw++                // sub-worker say
    go func() {
        partial := 0
        for i := range w.in { partial += i }
        w.out <- partial
        w.mtx.Lock()                  // sayğacı MÜDAFİƏ et
        w.sbw--
        if w.sbw == 0 { close(w.out) }  // sonuncu bitəndə out-u BAGLA
        w.mtx.Unlock()
    }()
}
func (w *Worker) gatherResult() int {
    total := 0
    wg := &sync.WaitGroup{}
    wg.Add(1)
    go func() {
        for i := range w.out { total += i }
        wg.Done()
    }()
    wg.Wait()
    return total
}
```
Kanallar struct İÇİNDƏ — parametr ötürmək lazım DEYİL; metod `go`-suz çağrılır,
özü daxildə goroutine yaradır.

### 14. Context — Rutinə Nəzarət
**Nədir:** çağrışlar boyunca ötürülən nəzarət konteyneri; DƏYƏR daşımaz —
İCRA-NI DAYANDIRMAQ üçün. HTTP zənglərində mütləq (server cavab vermirsa timeout).

**Kitabdan kod nümunəsi:**
```go
func countNumbers(c context.Context, r chan int) {
    v := 0
    for {
        select {
        case <-c.Done():          // context LƏĞV edildi?
            r <- v
            return                // SONSUZ loop-dan ÇIX
        default:
            time.Sleep(time.Millisecond * 100)
            v++
        }
    }
}

r := make(chan int)
c := context.TODO()
cl, stop := context.WithCancel(c)   // ləğv edilə bilən variant
go countNumbers(cl, r)
go func() {
    time.Sleep(300 * time.Millisecond)
    stop()                          // 300ms sonra LƏĞV
}()
v := <-r                            // 3 — 3 iterasiyadan sonra dayandı
```

## Əsas terminlələr
- Concurrency/Parallelism — hissə-hissə bölünmüş icra / fiziki eyni anda
- Goroutine — `go` ilə başladılan yüngül rutin; main-i gözləmir
- sync.WaitGroup — Add/Done/Wait sayğaclı sinxronizasiya
- Race Condition — kilidsiz paylaşılan dəyişən üzərində yarış
- `-race` — race detector test flagi
- sync/atomic — AddInt32 və s. atomik əməliyyatlar
- Mutex — Lock/Unlock mutual exclusion
- Channel — tipəmən mesaj kanalı; `<-` göndər/qəbul
- Buffered/Unbuffered — tutumlu (non-blocking yaz) / tutumsuz (sinxron uzlaşma)
- Deadlock — hamı gözləyir, heç kəs işləmir → crash
- close(ch) — mesajların sonu; range-i dayandırır
- for range ch — kanal bağlanana qədər qəbul
- Yönlü kanal — `<-chan T` (oxu) / `chan<- T` (yaz)
- Fan-out/Fan-in — N worker-ə paylama / tək nöqtədə birləşdirmə
- Worker Pool — sabit worker dəsti + paylaşılan in-kanalı
- Done Channel — bitmə siqnalı kanalı
- context.WithCancel / Done() — rutin ləğv mexanizmi
- "Share by communicating" — Go mottosu: paylaşmaq üçün KANAL, mutex yalnız lazımda

## Praktik nətidə

(1) `go f()` main-i gözləMİR — nəticə 0 çap olunması klassik ilk səhv. (2) Sleep
YOX — WaitGroup: Add go-dan əvvəl, Done goroutine içində (defer), Wait sonda.
(3) Paylaşılan dəyişən + çox goroutine = race; `-race` ilə test et — görünməz
oğru nəticələr YALNIZ detector ilə üzə çıxır. (4) Sadə sayğac/ədəd üçün atomic;
mürəkkəb əməliyyat üçün mutex. (5) Mutex bölgəsini MİNİMAL saxla — kilid altında
nə qədər az kod, o qədər çox concurrency. (6) Unbuffered kanal 2 rutin tələb edir —
tək rutində göndəriş = deadlock. (7) close-u GÖNDƏRƏN tərəf edir; range close-da
dayanır — worker pool-un mərkəzi mexanizmi. (8) Worker sayını sabit saxla, işin
bölünməsini kanala HƏVALƏ et (adil paylanma avtomatik). (9) Done channel —
WaitGroup-a alternativ sinxron bitmə siqnalı. (10) Kanalı imzada yönləndir
(`chan<-`) — səhv istiqamət compile xətası verir. (11) HTTP server hər requesti
goroutine-da işlədir — state üçün mutex ŞƏRT; amma "share by communicating" —
mümkünsə kanal seç. (12) Context: sonsuz/həddindən artıq işləri WithCancel+Done
ilə idarə et — HTTP timeout-un standart həlli.

## Mənbə
Pages: 559-602 (PDF 592-637)
