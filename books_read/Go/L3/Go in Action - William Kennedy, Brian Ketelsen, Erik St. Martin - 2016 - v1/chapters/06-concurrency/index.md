# Chapter 6 — Concurrency

## Bu chapter nədən bəhs edir?

Go-nun concurrency modeli: OS thread-ləri, logical processor-lər və Go runtime scheduler arasındakı əlaqə; goroutine yaradılması və idarəsi; race condition-lərin yaranması, aşkarlanması (`go build -race`) və aradan qaldırılması (atomic funksiyalar, mutex); unbuffered və buffered channel-lərin sinxronizasiya mexanikası.

## Əsas fikirlər

### 1. Concurrency vs Parallelism (eyni-vaxtlılıq vs paralellik)
**Nədir:** Concurrency — çox işin **idarə edilməsi**; parallelism — çox işin **eyni anda fiziki icrası**.

**Necə işləyir (arxitektura zənciri):**
- **Process (proses)** — tətbiqin resurs konteyneri: yaddaş ünvan sahəsi, fayl/device handle-ləri, thread-lər
- **Thread (axın)** — icra yolu; OS tərəfindən shedulə olunur; hər process-də ən az 1 thread (main thread) var; main thread bitəndə proqram bitir
- **Logical Processor (məntiqi prosessor)** — Go runtime-in uniti; hər biri **bir OS thread-ə** bağlıdır. Go 1.5-dən default: hər fiziki core üçün 1 logical processor
- **Goroutine** — scheduler tərəfindən logical processor-lərdə icra olunan müstəqil iş vahidi

**Scheduler işi:** goroutine yaradılır → global run queue → logical processor-ə təyin → **local run queue** → növbəsi gələndə icra.

**Blocking syscall (fayl açma kimi):** goroutine + thread logical processor-dən ayrılır → thread syscall gözləyir → scheduler **yeni thread yaradıb** processor-ə bağlayır → başqa goroutine icra olunur. Syscall qayıdanda goroutine local run queue-ya qaytarılır.

**Network I/O:** goroutine logical processor-dən ayrılıb **network poller**-ə keçirilir; oxuma/yazma hazır olanda geri qaytarılır.

**Limit:** hər proqram default olaraq maksimum **10,000 OS thread** (`runtime/debug.SetMaxThreads` ilə dəyişilə bilər); aşsa proqram crash edir.

**"Less is more" fəlsəfəsi:** concurrency çox vaxt parallelism-dən daha effektivdir — OS və hardware yükü az olur. Amma真正的 parallelism üçün: 1-dən çox logical processor + çoxfiziki core lazımdır.

### 2. Goroutine yaradılması və WaitGroup
**Nədir:** `go` açar sözü ilə funksiya (anonymous də olar) goroutine kimi icraya verilir.

**Kitabdan kod nümunəsi (alphabet nümunəsi):**
```go
package main

import (
    "fmt"
    "runtime"
    "sync"
)

func main() {
    // Allocate 1 logical processor for the scheduler to use.
    runtime.GOMAXPROCS(1)

    var wg sync.WaitGroup
    wg.Add(2)

    fmt.Println("Start Goroutines")

    // Declare an anonymous function and create a goroutine.
    go func() {
        defer wg.Done()

        for count := 0; count < 3; count++ {
            for char := 'a'; char < 'a'+26; char++ {
                fmt.Printf("%c ", char)
            }
        }
    }()

    go func() {
        defer wg.Done()

        for count := 0; count < 3; count++ {
            for char := 'A'; char < 'A'+26; char++ {
                fmt.Printf("%c ", char)
            }
        }
    }()

    fmt.Println("Waiting To Finish")
    wg.Wait()

    fmt.Println("\nTerminating Program")
}
```

**Sub-kod izahı:**
- `runtime.GOMAXPROCS(1)` → scheduler-ə 1 logical processor ayır; env variable GOMAXPROCS da var
- `wg.Add(2)` → 2 goroutine gözləniləcək
- `go func() { ... }()` → anonymous funksiya goroutine kimi
- `defer wg.Done()` → funksiya bitəndə sayğacı azaldır — defer bunu qarantiya edir
- `wg.Wait()` → main goroutine-lər bitməyincə gözləyir (main qayıtsa proqram bitir!)
- `runtime.GOMAXPROCS(runtime.NumCPU())` → hər core üçün 1 logical processor = paralel icra imkanı

**Goroutine davranışı:** scheduler hər hansı goroutine-in processor-u "əsir almasına" imkan vermir — işləyən goroutine dayandırılıb run queue-ya qaytarıla bilər (preemption). Uzun işlənən tapşırıqda (5000-ə qədər prime number tapma kimi) bu **time-slicing** açıq görünür: output-da `A:...`/`B:...` nümunələri qarışır.

**Vacib:** goroutine-dən return dəyəri ALMAQ olmur — nəticə channel ilə qaytarılır.

### 3. Race condition (yarış şəraiti)
**Nədir:** 2+ goroutine sinxronizasiyasız şəkildə eyni shared resursa eyni anda oxu/yazma etdikdə yaranan bug.

**Kitabdan kod nümunəsi (səhv kod):**
```go
var (
    counter int
    wg sync.WaitGroup
)

func incCounter(id int) {
    defer wg.Done()

    for count := 0; count < 2; count++ {
        // Capture the value of Counter.
        value := counter

        // Yield the thread and be placed back in queue.
        runtime.Gosched()

        // Increment our local value of Counter.
        value++

        // Store the value back into Counter.
        counter = value
    }
}
// main: go incCounter(1); go incCounter(2)
// Gözlənilən: 4, Alınan: Final Counter: 2
```

**Necə baş verir:** hər goroutine `counter`-in kopyasını götürür (`value := counter`), sonra `runtime.Gosched()` ilə thread-i təslim edir — swap zamanı digər goroutine eyni köhnə dəyəri oxuyur; hər ikisi eyni dəyəri artırıb yazır — biri digərinin işini **overwrite** edir. Nəticə: 4 əvəzinə 2.

**Aşkarlama — race detector:**
```bash
go build -race    # race detector flag ilə build
./example         # işə sal
```
Output:
```
==================
WARNING: DATA RACE
Write by goroutine 5:
  main.incCounter()  /example/main.go:49
Previous read by goroutine 6:
  main.incCounter()  /example/main.go:40
...
Found 1 data race(s)
```

### 4. Atomic funksiyalar
**Nədir:** `sync/atomic` paketinin aşağı-səviyyəli lock mexanizmi — int və pointer-lər üçün.

**Kitabdan kod nümunəsi (AddInt64):**
```go
var (
    counter int64
    wg sync.WaitGroup
)

func incCounter(id int) {
    defer wg.Done()

    for count := 0; count < 2; count++ {
        // Safely Add One To Counter.
        atomic.AddInt64(&counter, 1)

        runtime.Gosched()
    }
}
// Final Counter: 4 — düzgün!
```

**Sub-kod izahı:**
- `atomic.AddInt64(&counter, 1)` → əlavə əməliyyatı atomikdir — yalnız bir goroutine eyni anda icra edə bilər
- Digər faydalılar: `atomic.LoadInt64` (təhlükəsiz oxu), `atomic.StoreInt64` (təhlükəsiz yaz)

**Shutdown flag nümunəsi:**
```go
var shutdown int64

// worker goroutine:
for {
    fmt.Printf("Doing %s Work\n", name)
    time.Sleep(250 * time.Millisecond)
    if atomic.LoadInt64(&shutdown) == 1 {
        fmt.Printf("Shutting %s Down\n", name)
        break
    }
}

// main:
time.Sleep(1 * time.Second)
atomic.StoreInt64(&shutdown, 1)   // bütün worker-lərə siqnal
```
- `Load` təhlükəsiz kopya qaytarır, `Store` təhlükəsiz yazır — eyni vaxtda çağrılsa belə sinxronlaşdırılır

### 5. Mutex (mutual exclusion)
**Nədir:** `sync.Mutex` — kritik bölmə (critical section) yaradan lock; bir vaxtda yalnız bir goroutine daxil ola bilər.

**Kitabdan kod nümunəsi:**
```go
var (
    counter int
    wg sync.WaitGroup
    mutex sync.Mutex
)

func incCounter(id int) {
    defer wg.Done()

    for count := 0; count < 2; count++ {
        // Only allow one goroutine through this critical section at a time.
        mutex.Lock()
        {
            value := counter
            runtime.Gosched()
            value++
            counter = value
        }
        mutex.Unlock()
    }
}
// Final Counter: 4 — düzgün!
```

**Sub-kod izahı:**
- `mutex.Lock()` → kritik bölmənin başlanğıcı; digər goroutine-lər burada bloklanır
- `mutex.Unlock()` → lock azad edilir, gözləyən növbəti goroutine daxil olur
- Mötərizələr `{}` məcburi deyil — yalnız kritik bölməni görsəl dəqiqələşdirir
- `Gosched()` içəridə olsa belə təhlükəsizdir — swap olsa da lock qorunur

### 6. Channels (kanallar) — CSP modeli
**Nədir:** Goroutine-lər arasında tipli mesaj ötürülməsi üçün konduit. CSP (communicating sequential processes) paradigmının əsas tipi: "data-nı locklama, ötür".

**Yaradılış və əməliyyatlar:**
```go
// Unbuffered channel of integers.
unbuffered := make(chan int)

// Buffered channel of strings.
buffered := make(chan string, 10)

// Send a string through the channel.
buffered <- "Gopher"

// Receive a string from the channel.
value := <-buffered
```

**Sub-kod izahı:**
- `make(chan T)` → unbuffered (tutum 0)
- `make(chan T, N)` → buffered (tutum N)
- `ch <- v` → göndər (send)
- `v := <-ch` → qəbul (receive) — `<-` operatoru sol tərəfdə

### 7. Unbuffered channel — qarantiyalı mübadilə
**Nədir:** Heç bir dəyəri saxlaya bilməyən kanal — send və receive **eyni anda hazır** olmalıdır.

**Necə işləyir (tennis topu modeli):** göndərən goroutine kanala "əlini uzadır" — alan hazır olana qədər **kanala kilidlənir**; alan da eyni şəkildə kilidlənir; mübadilə bitəndə hər ikisi azad olur. Sinxronizasiya daxilidir — bir digərsiz baş verə bilməz.

**Kitabdan kod nümunəsi (tennis oyunu):**
```go
// player simulates a person playing the game of tennis.
func player(name string, court chan int) {
    defer wg.Done()

    for {
        // Wait for the ball to be hit back to us.
        ball, ok := <-court
        if !ok {
            // If the channel was closed we won.
            fmt.Printf("Player %s Won\n", name)
            return
        }

        // Pick a random number and see if we miss the ball.
        n := rand.Intn(100)
        if n%13 == 0 {
            fmt.Printf("Player %s Missed\n", name)

            // Close the channel to signal we lost.
            close(court)
            return
        }

        // Display and then increment the hit count by one.
        fmt.Printf("Player %s Hit %d\n", name, ball)
        ball++

        // Hit the ball back to the opposing player.
        court <- ball
    }
}

// main:
court := make(chan int)
wg.Add(2)
go player("Nadal", court)
go player("Djokovic", court)
court <- 1        // oyun başlayır
wg.Wait()
```

**Sub-kod izahı:**
- `ball, ok := <-court` → iki dəyərli receive: dəyər + "kanal açıqdırmı?" flag-i
- `ok == false` → kanal **closed** — oyun bitib (qazandın)
- `close(court)` → uduzduğunu bildirir; qapalı kanal üzərində receive dərhal zero value + `ok=false` qaytarır
- `court <- ball` → topu geri göndər — digər oyunçu receive gözləyir, hər ikisi mübadiləyə kilidlənir

**Estafet (relay race) nümunəsinin məntiqi:** eyni unbuffered sinxronizasiya — `baton := make(chan int)`, hər runner `runner := <-baton` gözləyir, `baton <- newRunner` ilə ötürür.

### 8. Buffered channel — müstəqil iş
**Nədir:** N dəyər saxlaya bilən kanal; send/receive eyni anda olmağı tələb etmir.

**Bloklanma qaydaları:**
- **Receive** yalnız kanal boş olduqda bloklanır
- **Send** yalnız buffer dolu olduqda bloklanır
- Unbuffered-dən fərqli olaraq mübadilənin "eyni anda" qarantiyası YOXDUR

**Kitabdan kod nümunəsi (worker pool):**
```go
const (
    numberGoroutines = 4  // Number of goroutines to use.
    taskLoad          = 10 // Amount of work to process.
)

func main() {
    // Create a buffered channel to manage the task load.
    tasks := make(chan string, taskLoad)

    // Launch goroutines to handle the work.
    wg.Add(numberGoroutines)
    for gr := 1; gr <= numberGoroutines; gr++ {
        go worker(tasks, gr)
    }

    // Add a bunch of work to get done.
    for post := 1; post <= taskLoad; post++ {
        tasks <- fmt.Sprintf("Task : %d", post)
    }

    // Close the channel so the goroutines will quit
    // when all the work is done.
    close(tasks)

    wg.Wait()
}

// worker is launched as a goroutine to process work from the buffered channel.
func worker(tasks chan string, worker int) {
    defer wg.Done()

    for {
        // Wait for work to be assigned.
        task, ok := <-tasks
        if !ok {
            // This means the channel is empty and closed.
            fmt.Printf("Worker: %d : Shutting Down\n", worker)
            return
        }

        fmt.Printf("Worker: %d : Started %s\n", worker, task)

        // Randomly wait to simulate work time.
        sleep := rand.Int63n(100)
        time.Sleep(time.Duration(sleep) * time.Millisecond)

        fmt.Printf("Worker: %d : Completed %s\n", worker, task)
    }
}
```

**Sub-kod izahı:**
- `make(chan string, 10)` → 10 tapşırıq buffer-i
- 4 worker goroutine `task, ok := <-tasks` gözləyir
- `close(tasks)` → bütün tapşırıqlar göndərildikdən sonra bağlanır
- **Qapalı kanal qaydaları:** send olmaz (panic), receive **olur** — kanal boşalana qədər dəyərlər itmir; boş + qapalı kanaldan receive dərhal zero value + `ok=false` qaytarır → worker özünü söndürür
- `time.Duration(sleep) * time.Millisecond` → int-i Duration-a çevirir

## Unbuffered vs Buffered müqayisəsi
| Xüsusiyyət | Unbuffered | Buffered |
|---|---|---|
| Tutum | 0 | N |
| Send/Receive sinxronluğu | Eyni anda — qarantiya var | Müstəqil — qarantiya yoxdur |
| Send bloklanması | alan hazır deyilsə | buffer doludursa |
| Receive bloklanması | göndərən yoxdursa | kanal boşdursa |
| İstifadə | handshake, notification, dəqiq mübadilə | worker queue, tapşırıq paylayıcı |

## Əsas terminlər
- Concurrency / Parallelism (eyni-vaxtlılıq / paralellik)
- Process / Thread (proses / axın)
- Logical Processor (məntiqi prosessor)
- Scheduler (vaxt planlayıcısı)
- Run Queue (icra növbəsi)
- Preemption / Time-slicing (zorla dayandırma / vaxt dilimlənməsi)
- Race Condition (yarış şəraiti)
- Race Detector (`-race` flag)
- Atomic Functions (atomik funksiyalar)
- Mutex / Critical Section (mutual exclusion / kritik bölmə)
- CSP (communicating sequential processes — mesaj ötürmə modeli)
- Unbuffered / Buffered Channel (bufersiz / bufferli kanal)
- Blocking / Non-blocking Syscall (bloklayan sistem çağırışı)
- Network Poller (şəbəkə sorğulayıcısı)

## Praktik nəticə
- Kod linear yaza bilirsənsə — concurrency-dən imtina et (sadəlik).
- GOMAXPROCS-u kor-koranə dəyişmə; dəyişməzdən əvvəl benchmark et.
- Concurrent kodu həmişə `go build -race` ilə yoxla — race condition-lər görünməz olur.
- Yalnız int/pointer qoruyursansa atomic; kod bloku qoruyursansa mutex; **məntiqləri bölüşmək lazımdırsa channel** (idiomatik seçim).
- Worker pool pattern: buffered channel + N goroutine + `close()` — shutdown siqnalı üçün 2 dəyərli receive (`task, ok := <-ch`).
- Unbuffered kanal handshake qarantiyası verir — dəqiq mübadilə lazım olduqda onu seç.

## Mənbə
Pages: 149-178 (PDF), book pages 128-157
