# Chapter 3 — Concurrency in Go (Technique 10-15)

## Bu chapter nədən bəhs edir?

Go-nun CSP concurrency modeli: goroutine-lər (closure, scheduler, Gosched, loop-variable tələsi), sync.WaitGroup (paralel gzip), sync.Mutex (race condition + word counter nümunəsi, --race detector), channel-lər (multi-channel + select, direction, time.After), channel closing (səhv/düzgün üsullar, done-channel pattern), buffered channel ilə lock.

## Əsas fikirlər

### 1. CSP modeli
Go thread modeli istifadə etmir — **Communicating Sequential Processes** (Tony Hoare):
- **Goroutine:** müstəqil işləyən funksiya ("öz thread-ində işləyən funksiya kimi")
- **Channel:** data ötürmə kanalı — "proqram daxilində socket"; **tipləşdirilmiş**, mərhələmə lazım deyil

Concurrency Go-da ucuzdur → kitabxanalarda tez-tez istifadə olunur. **Diqqət:** goroutine/channel — Go-da memory leak-in ƏSAS mənbəyi.

### 2. Goroutine əsasları
```go
go echo(os.Stdin, os.Stdout)      // echo funksiyasını background-da işə sal
time.Sleep(30 * time.Second)
fmt.Println("Timed out.")
os.Exit(0)

func echo(in io.Reader, out io.Writer) {
    io.Copy(out, in)
}
```
- `go` açar sözü = funksiyanı cədvələşdir; main çıxanda proqram dayanır
- **Concurrent ≠ parallel** — scheduler icrayı bölüşdürür

### TECHNIQUE 10: Goroutine closure + Gosched
```go
func main() {
    fmt.Println("Outside a goroutine.")
    go func() {
        fmt.Println("Inside a goroutine")
    }()
    fmt.Println("Outside again.")
    runtime.Gosched()      // scheduler-ə yield et!
}
```

**Gosched olmadan:** "Inside a goroutine" ÇAP OLUNMAYA BİLƏR — main qayıdandan əvvəl goroutine icra olunmur!

**runtime.Gosched()** = "bu nöqtədə pauza verə bilərəm" — scheduler başqa goroutine-ləri işə sala bilər. Amma **tamamlanma zəmanəti YOX** (DB query gözləyən goroutine işə düşməyə bilər) → WaitGroup istifadə et.

### TECHNIQUE 11: sync.WaitGroup — paralel gzip
**Serial versiya:** 1 goroutine, 1 core — fayllar sıra ilə.
**Paralel versiya:**
```go
func main() {
    var wg sync.WaitGroup          // init LAZIM DEYİL — zero value işləyir!

    var i int = -1
    var file string
    for i, file = range os.Args[1:] {
        wg.Add(1)                  // "bir iş daha gözləyirəm"
        go func(filename string) {  // ← PARAMETR KİMİ, closure YOX!
            compress(filename)
            wg.Done()              // "bir iş bitdi"
        }(file)
    }
    wg.Wait()                      // hamısı bitənə qədər blokla
    fmt.Printf("Compressed %d files\n", i+1)
}

func compress(filename string) error {
    in, err := os.Open(filename)
    defer in.Close()
    out, err := os.Create(filename + ".gz")
    defer out.Close()
    gzout := gzip.NewWriter(out)
    _, err = io.Copy(gzout, in)
    gzout.Close()
    return err
}
```

**Loop-variable tələsi (vacib!):** `file` loop dəyişənidır — closure istifadə etsəydİ, bütün goroutine-lər SON dəyəri görərdi. **Həll: parametr kimi ötür** — dəyər schedule ediləndə kopyalanır.

WaitGroup 3 metodu: `Add(n)` (gözlənilən iş sayı), `Done()` (bitiş siqnalı), `Wait()` (bloklama).

### TECHNIQUE 12: sync.Mutex — race condition
**Word counter (race-li versiya):**
```go
// SÖZ DƏYİŞMƏZ (səhv):
type words struct {
    found map[string]int
}

func (w *words) add(word string, n int) {
    count, ok := w.found[word]
    if !ok {
        w.found[word] = n
        return
    }
    w.found[word] = count + n
}
```
Nəticə: `fatal error: concurrent map writes` — bir neçə goroutine eyni map-ə yazır.

**`--race` flag:** `go run --race race.go *.txt` → WARNING: DATA RACE + dəqiq sətir göstərilir. Development-də istifadə et (yavaşladır amma dəqiq tapır).

**DÜZGÜN (mutex-lu):**
```go
type words struct {
    sync.Mutex                    // ANONIM EMBEDDING → Lock/Unlock metodları!
    found map[string]int
}

func (w *words) add(word string, n int) {
    w.Lock()                     // kilidlə
    defer w.Unlock()             // çıxanda aç (hər return-də!)
    count, ok := w.found[word]
    if !ok {
        w.found[word] = n
        return
    }
    w.found[word] = count + n
}
```
- **sync.Mutex-i struct-a embed et** → `words.Lock()` / `words.Unlock()` açılır — çox rast gəlinən pattern
- `defer w.Unlock()` — unutma ehtimalı sıfır

**Vacib:** lock YALNIZ hamı eyni lock-dan istifadə edəndə işləyir — biri locksuz çıxış etsə race qalır.

**Digər:** `sync.RWLock` — çox oxu + tək yazma; "öz lock-unu yazma" — RUMOR: builtin lock yavaşdır → YALAN; battle-tested.

### 3. Channel əsasları
Socket metaforası:
| Socket | Channel |
|---|---|
| 2 proqram arası | 2 goroutine arası |
| unidirectional/bidirectional | hər ikisi |
| raw bytes | **tipləşdirilmiş data** |
| marshal lazım | mərhələmə YOX |

```go
ch := make(chan []byte)      // unbuffered — 1 dəyər
ch <- v                       // send
v := <-ch                     // receive
out chan<- []byte             // send-only (function parametr)
done <-chan bool              // receive-only
```

### TECHNIQUE 13: Multi-channel + select
**Echo program (channel-lu versiya):**
```go
func main() {
    done := time.After(30 * time.Second)    // <-chan time.Time!
    echo := make(chan []byte)
    go readStdin(echo)

    for {
        select {
        case buf := <-echo:
            os.Stdout.Write(buf)
        case <-done:
            fmt.Println("Timed out")
            os.Exit(0)
        }
    }
}

func readStdin(out chan<- []byte) {   // write-only — oxumaq compile xətası!
    for {
        data := make([]byte, 1024)
        l, _ := os.Stdin.Read(data)
        if l > 0 {
            out <- data
        }
    }
}
```

**time.After(d):** `<-chan time.Time` qaytarır — d-dən sonra dəyər gəlir. `time.Sleep(5s)` ≡ `<-time.After(5s)`.

**select semantikası:**
- 0+ case + optional default
- 1 case hazır → onu icra et
- >1 hazır → RANDOM seç
- 0 hazır + default YOX → **blokla**
- Function signature-də channel **istiqamətini göstər** — yaxşı praktika

### TECHNIQUE 14: Channel closing — düzgün yolu
**Qayda 1:** `close` YALNIZ GÖNDƏRƏN tərəfindən!

**Səhv nümunə (receiver bağlayır):**
```go
case <-until:
    close(msg)        // ← RECEIVER BAĞLAYIR — send hələ yazır!
    return
// → panic: send on closed channel
```

**Səhv nümunə 2 (sender bağlayır, loop davam edir):**
```go
func send(ch chan bool) {
    time.Sleep(120 * time.Millisecond)
    ch <- true
    close(ch)
}
// main-dəki select: bağlı kanaldan oxuma → DƏRHAL zero value (false) qaytarır
// → "Got message." minlərlə dəfə — sonsuz spin!
```
**Bağlı kanal:** hər oxuma = dərhal zero value — döngü yoxlaması olmadan sonsuz dövr.

**DÜZGÜN — done-channel pattern:**
```go
func main() {
    msg := make(chan string)
    done := make(chan bool)              // ayrı siqnal kanalı!
    until := time.After(5 * time.Second)
    go send(msg, done)

    for {
        select {
        case m := <-msg:
            fmt.Println(m)
        case <-until:
            done <- true                  // receiver SİQNAL verir
            time.Sleep(500 * time.Millisecond)
            return
        }
    }
}

func send(ch chan<- string, done <-chan bool) {
    for {
        select {
        case <-done:                     // siqnal gəldi
            println("Done")
            close(ch)                    // GÖNDƏRƏN bağlayır ✓
            return
        default:
            ch <- "hello"
            time.Sleep(500 * time.Millisecond)
        }
    }
}
```
**Pattern:** receiver bitmə şərtini bilir → `done` kanalından xəbər verir → **sender** öz kanalını bağlayıb qayıdır. Hər iki tərəf təmiz çıxır — leak YOX.

### TECHNIQUE 15: Buffered channel ilə lock
**Unbuffered:** send bloklanır (alan gözləməsə) — lock üçün yararsız.
**Buffered (size 1):** buferdə yer varsa send bloklanmır!

**Lock mexanizmi:**
```go
lock := make(chan bool, 1)     // 1 yerlik bufer

func worker(id int, lock chan bool) {
    fmt.Printf("%d wants the lock\n", id)
    lock <- true                // ACQUIRE: buferə yaz — yer yoxdursa BLOK
    fmt.Printf("%d has the lock\n", id)
    time.Sleep(500 * time.Millisecond)
    fmt.Printf("%d is releasing the lock\n", id)
    <-lock                      // RELEASE: buferdən oxu → yer açılır
}

func main() {
    lock := make(chan bool, 1)
    for i := 1; i < 7; i++ {
        go worker(i, lock)
    }
    time.Sleep(10 * time.Second)
}
// 2 wants... 2 has... 2 releasing → 1 has... → ardıcıllıq
```
- `lock <- true` = kilidi götür; `<-lock` = burax
- İlk gələn 1 yeri tutur; qalanları blok
- **Size N = N paralel kilid** mümkün
- Buferli kanallar həm də message queue/pipeline üçün

## Concurrency qaydaları xülasəsi
1. Loop daxilində goroutine → dəyişəni PARAMETR kimi ötür
2. Bitmə gözlə → WaitGroup; Gosched zəmanət vermir
3. Shared data + yazı → mutex (struct-a embed et, defer Unlock) yaxud kanal
4. `go run --race` — development-də həmişə
5. close YALNIZ sender tərəfdən + done-kanalı ilə siqnal
6. Bağlı kanal zero value qaytarır — sonsuz loop təhlükəsi
7. Kanalı həddindən artıq istifadə etmə — overhead + memory leak mənbəyi

## Əsas terminlər
- CSP (Communicating Sequential Processes — Tony Hoare)
- Goroutine (müstəqil funksiya; concurrent ≠ parallel)
- runtime.Gosched (yield)
- sync.WaitGroup: Add / Done / Wait
- Loop Variable Closure Tələsi (parametr kimi ötür!)
- Race Condition / `--race` Detector
- sync.Mutex / sync.RWLock / sync.Locker
- Mutex Embedding (struct-da anonim sahə)
- defer Unlock
- Channel (tipli daxili socket)
- Unbuffered / Buffered Channel (make 2-ci arqument)
- Direction: `chan<-` / `<-chan`
- time.After (`<-chan time.Time`)
- select: random choice / default / blok
- Done-channel Pattern (close siqnalı)
- Zero Value on Closed Channel (sonsuz loop tələsi)
- Channel-based Lock (bufer=1)
- Channel/Goroutine Leak

## Praktik nəticə
- Gosched/Sleep ilə "gözləmə" yox — WaitGroup və ya done-kanal düzgün senxronizasiyadır.
- Mutex-i struct-a embed et + defer Unlock — unutma mümkün deyil; amma hamı eyni lock-dan istifadə etməlidir.
- close çağırışını YALNIZ sender edir; receiver done-kanalı ilə xəbər verir — bu, kitabın ən vacib kanal qaydasıdır.
- Bağlı kanaldan oxuma zero value verir — range/ok-yoxlaması olmadan loop-da spin.
- Buferli kanal (size 1) = kanal-əsaslı kilid; size N = N-lik semafor.
- --race development cycle-da standart; performansa təsir edir, amma race-ləri dəqiq göstərir.
- Kanallar güclüdür amma hammer-deyil — hər yerde istifadə etmə; leak mənbəyidir.

## Mənbə
Pages: 82-106 (PDF), book pages 59-83
