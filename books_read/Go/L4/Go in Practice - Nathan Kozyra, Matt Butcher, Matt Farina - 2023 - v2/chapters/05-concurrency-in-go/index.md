# Chapter 5 — Concurrency in Go (Go-da Paralellik)

## Bu chapter nədən bəhs edir?

CSP modeli, goroutine-lərin başladılması, sync.WaitGroup (paralel gzip nümunəsi), loop
dəyişəni tələsi, sync.Mutex/RWMutex (word counter race nümunəsi), channel-lər (socket
analogiyası, select, time.After), channel bağlama qaydaları (done channel pattern), buffered
channel ilə kilidləmə.

## Əsas fikirlər

### 1. Goroutine Əsasları
**Nədir:** `go` keyword-ilə çağırılan funksiya — müstəqil işləyən yüngül "thread".
Java/Python OS/green thread modelindən fərqli olaraq Go CSP (Communicating Sequential
Processes, Tony Hoare) modelini izləyir.

**Kitabdan kod nümunəsi (echo, 30 saniyəlik):**
```go
func main() {
    fmt.Println("Type anything below for up to 30 seconds")
    go echo(os.Stdin, os.Stdout)     // fondda işə düşür
    time.Sleep(30 * time.Second)
    fmt.Println("Timed out.")
    os.Exit(0)
}
func echo(in io.Reader, out io.Writer) { io.Copy(out, in) }
```
**Əsas qayda:** main çıxanda bütün proqram QURTARULUR — goroutine işdə olsa belə. time.Sleep
burada həyat dövrünü saxlayır (production-da context.WithDeadline bunu avtomatiklaşdırır).

### 2. sync.WaitGroup — Paralel İşlərin Gözlənilməsi
**Pattern:** N işçi + 1 gözləyən.

**Kitabdan kod nümunəsi (paralel gzip):**
```go
var wg sync.WaitGroup            // init tələb etməz (zero value)
for i, file := range os.Args[1:] {
    wg.Add(1)                     // hər goroutine ƏVVƏL sayğaca daxil et
    go func(filename string) {
        compress(filename)
        wg.Done()
    }(file)
}
wg.Wait()                          // hamı bitənə qədər blokla
fmt.Printf("Compressed %d files\n", len(os.Args[1:]))
```

**Qaydalar:**
- `Add`-i loop İÇİNDƏ hər goroutine üçün çağır (bütün sayı əvvəlcədən Add etmək —
  loop ortasında xəta olsa Wait əbədi asılır)
- `Done`-u goroutine daxilində (idealdır `defer wg.Done()`) — Add/Done yaxın durur
- **Loop dəyişəni tələsi (mühüm!):** closure `file`-i capture etsə, goroutine-lər
  işə düşəndə hamısı SON `file` dəyərini görər → **parametr kimi ötür** (`(file)`)
- Worker funksiyası (compress) WaitGroup-u BİLMİRLƏR — əmək bölgüsü: fayl sıralı
  istifadə üçün də eyni funksiya işləyir

### 3. Race Condition və sync.Mutex
**Problem:** İki+ goroutine eyni dataya yaza bilər → "fatal error: concurrent map writes"
və ya səssiz data korlanması. `--race` flagi data race-i dəqiq göstərir.

**Kitabdan kod nümunəsi (word counter):**
```go
type words struct {
    sync.Mutex                     // EMBED — Lock/Unlock metodları promote olunur
    found map[string]int
}
func (w *words) add(word string) {
    w.Lock()
    defer w.Unlock()
    if count, ok := w.found[word]; ok {
        w.found[word] = count + 1
    } else {
        w.found[word] = 1
    }
}
// Oxuyanda da: w.Lock(); defer w.Unlock()
```
**Vacib:** Kilid YALNIZ bütün çıxışlar eyni kiliddən keçəndə işləyir — bir yer kilidsiz
qalsa race qalır. `sync.RWMutex`: RLock/RUnlock — çoxlu paralel oxu, tək yazma.

### 4. Channel-lər — Socket Analogiyası
**Nədir:** Goroutine-lərarası TİPƏMƏN mesaj kanalı (socket kimi, amma marshal YOX).

**Kitabdan kod nümunəsi (echo + select):**
```go
echo := make(chan []byte)           // make ilə; buffer-siz = 1 mesaj bloklanır
go readStdin(echo)
for {
    select {
    case buf := <-echo:
        os.Stdout.Write(buf)
    case <-time.After(30 * time.Second):
        return
    }
}
func readStdin(out chan<- []byte) {   // YALNIZ-yaz kanalı — compile-time qoruma
    for {
        data := make([]byte, 1024)
        l, _ := os.Stdin.Read(data)
        if l > 0 { out <- data }
    }
}
```

**select semantikası:** case-lərdən HAZIR olanı icra; BİRDƏN ÇOX hazırdırsa RANDOM seç;
heç biri hazır deyil + default YOX → blokla. time.After(d) — d sonra mesaj verən kanal
qaytarır (time.Sleep-in kanal forması). Tək kanal dinləyirsənsə `for range ch` sadədir
(amma select daha ümumi).

### 5. Channel Bağlama — Səhv və Düzgün Yollar
**Leak problemi:** Açıq kanal + bitməyən goroutine GC tərəfindən yığıla BİLMƏZ.

**Səhv 1 — receiver close edir:**
```go
case <-until:
    close(msg)      // send hələ yazır → panic: send on closed channel!
```
**Səhv 2 — sender close edir, receiver bilmir:**
```go
// send: ch <- true; close(ch)
// select: bağlı kanaldan oxu HƏMİŞƏ "hazır"dır — zero value qaytarır (bool üçün false)
// Nəticə: minlərlə "Got message" — sonsuz zibil fırtınası
```
(comma-ok ilə `m, ok := <-ch; if !ok` bağlılıq yoxlanılır — amma daha yaxşısı aşağıdakı.)

**Düzgün — done channel pattern:**
```go
msg := make(chan string)
done := make(chan bool)
go send(msg, done)
for {
    select {
    case m := <-msg: log.Println(m)
    case <-time.After(5 * time.Second):
        done <- true       // receiver SIQNAL verir
        return
    }
}
func send(ch chan<- string, done <-chan bool) {
    for {
        select {
        case <-done:
            close(ch)      // SENDER bağlayır — məcburi qayda
            return
        default:
            ch <- "hello"
            time.Sleep(500 * time.Millisecond)
        }
    }
}
```
**Prinsiplər:** close-u YALNIZ sender edir; "bitdim" siqnalı ayrıca done kanalı ilə;
hər iki tərəf təmiz şəkildə çıxır.

### 6. Buffered Channel ilə Kilid
**İdeya:** sync.Mutex əvəzinə kanal-orient kodda stil baxımından; buffer = "kilid slotu".

**Kitabdan kod nümunəsi:**
```go
lock := make(chan bool, 1)      // 1 slotlu buffer
for i := 1; i < 7; i++ {
    go worker(i, lock)
}
func worker(id int, lock chan bool) {
    log.Printf("%d wants the lock\n", id)
    lock <- true                 // KİLİD AL — slot doludursa BLOKLANIR
    log.Printf("%d has the lock\n", id)
    <-lock                       // KİLİD BURAX — slot boşalır
}
```
- Unbuffered bunu BACARMIR (send bloklanır) — buffer dolmayınca yaz sərbəst
- Buffer-2 = eyni anda 2 goroutine; message queue/pipeline də bu prinsiplə
- Amma mutex əvəzi DEYİL — "oxşar sahə"; race detector yenə vacibdir

**Channel istismarı xəbərdarlığı:** channel-lər overhead + mürəkkəblik daşıyır; Go-da
yaddaş idarəetmə problemlərinin ƏN BÖYÜK mənbəyi — ehtiyac yarananda işlət, hər yerə çəkib-getirmə.

## Əsas terminlələr
- CSP — Communicating Sequential Processes modeli (Hoare)
- WaitGroup — Add/Done/Wait sayğaclı gözləmə strukturu
- Race condition — paralel çıxışın data yarışı
- sync.Mutex / RWMutex — eksklüziv / oxu-paralel kilidlər
- select — çoxkanal hazır-olma gözətçisi (random seçim)
- time.After — müddət sonunda siqnal verən kanal
- Done channel — "bit" siqnalı kanalı (close-u sender-ə həvalə edir)
- Buffered channel — N-slotlu; buffer dolu olmayana qədər yaz bloklanmır

## Praktik nətidə

Paralellik qərarları: (1) fon iş — `go f()`; amma main-in həyatı = proqramın həyatı;
(2) N paralel iş + gözləmə — WaitGroup (Add loop-da, Done defer ilə); (3) loop dəyişəni
goroutine-a PARAMETR; (4) paylaşılan dəyişən — Mutex (bütün çıxışlar!); çox oxu/tək yazma —
RWMutex; (5) kanallar mesaj ötürür — select ilə çox kanal idarəsi; (6) close-u ancaq
sender; "bit" xəbəri ayrıca done kanalı ilə; (7) kanal-kilid — buffer 1 + send/read
cütü; (8) `--race` ilə test HƏMİŞƏ; (9) kanalı hər problemə çəkib aparma.

## Mənbə
Pages: 115-137 (PDF 136-158)
