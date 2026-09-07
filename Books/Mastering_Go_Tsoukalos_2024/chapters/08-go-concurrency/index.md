# Chapter 8 — Go Concurrency (Go-da Paralellik)

## Bu chapter nədən bəhs edir?

Proses/thread/goroutine ayrımı, Go scheduler (m:n, work-stealing, GOMAXPROCS), concurrency
vs parallelism, goroutine yaradılışı, sync.WaitGroup, channel-lər (buffered/nil/signal,
yönümlü parametrlər, close semantikası), race condition + `-race` detektora, select,
timeout (time.After), worker pool, goroutine icra sırası, UNIX siqnal tutma, sync.Mutex/
RWMutex/atomic, monitor goroutine, closure tələsi, context paketi, semaphore paketi və
statistika tətbiqinin concurrency versiyası.

## Əsas fikirlər

### 1. Proses, Thread, Goroutine
- **Proses** — işləyən proqramın OS təmsili (yaddaş, deskriptor, resurslarla)
- **Thread** — prosesin içindəki daha yüngül icra vahidi; öz stack + idarə axını
- **Goroutine** — ƏN KİÇİ Go icra vahidi; thread-dən də yüngül (böyüyə bilən kiçik stack,
  sürətli start); minlərləsi bir maşında problem deyil. Goroutine-lər thread-lərdə, thread-lər
  proseslərdə yaşayır.

**Goroutine-lər birbaşa kommunikasiya edə BİLMİR** — channels / socket / shared memory.

### 2. Go Scheduler və GOMAXPROCS
**Necə işləyir:** m:n scheduling — m goroutine n OS thread üzərində multiplex. 3 entities:
**M** (OS threads), **G** (goroutines), **P** (logical processors). GOMAXPROCS = eyni anda
istifadə edilə bilən P sayı (Go 1.5+: default = logical core sayı).

**Work-stealing:** boş qalan P digərlərinin local queue-sundan iş OĞURLAYIR; Go
continuation stealing istifadə edir (task stealing yerinə) — daha tez-tez rast gəlinən
vahid oğurlanır. Global run queue + per-P local queue.

**Goroutine vəziyyətləri:** executing (M-də icra), runnable (M gözləyir), waiting (bloklanır).

```go
runtime.GOMAXPROCS(0)     // 0 = dəyişmə, cari dəyəri QAYTARIR; ≥1 = təyin edir
fmt.Println(runtime.Compiler, runtime.GOARCH, runtime.Version())
```
Core-dan böyük GOMAXPROCS context-switchingə görə SÜRƏT VERMİR.

### 3. Concurrency ≠ Parallelism
- **Concurrency (paralellik strukturu)** — komponentləri müstəqil icra oluna biləcək şəkildə
  STRUKTURLAŞDIRMAQ — dizayn qərarı
- **Parallelism (paralellik icrası)** — eyni anda fiziki icra — hardware/OS imkanı
- Düzgün concurrent dizayn paralelliyi avtomatik gətirir (imkan olsa); parallel imkan
  olmasa belə concurrent dizayn data flow + maintainability YAXŞILAŞDIRIR. 
  **"Concurrency is better than parallelism."**

### 4. Goroutine Yaratma və sync.WaitGroup
**Kitabdan kod nümunəsi:**
```go
go func(x int) {
    fmt.Printf("%d ", x)
}(10)                    // dəyər PARAMETR kimi — closure-dan daha oxunaqlı

go printme(15)           // adlı funksiya

// Düzgün gözləmə (Sleep YOX!):
var waitGroup sync.WaitGroup
for i := 0; i < count; i++ {
    waitGroup.Add(1)                  // go-dan ƏVVƏL — race qorunması
    go func(x int) {
        defer waitGroup.Done()        // Done = Add(-1)
        fmt.Printf("%d ", x)
    }(i)
}
waitGroup.Wait()                      // sayğac 0 olana qədər blokla
```
**Add/Done uyğunsuzluğu:**
- Done > Add → `panic: sync: negative WaitGroup counter`
- Add > Done → `fatal error: all goroutines are asleep - deadlock!` (Wait əbədi gözləyir)
- Add(n) bir çağırışla mümkün; Done yalnız 1 azaldır

### 5. Channel-lər — Əsas Semantika
```go
c := make(chan int)        // unbuffered — sender receiver gözləyir
c := make(chan int, 5)     // buffered — 5 dəyərə qədər receiver-siz
c <- x                     // yaz
v := <-c                   // oxu
_, ok := <-c               // ok=false → bağlı kanal
for v := range ch { }      // yalnız close(olanda çıxır!
close(c)                   // YALNIZ SENDER bağlayır
```

**Qaydalar cədvəli:**
| Əməliyyat | Nəticə |
|---|---|
| Bağlı kanala YAZ | **panic: send on closed channel** |
| Bağlı kanaldan OXU | zero value (təhlükəsiz) |
| nil kanala yaz / nil-dən ox | **BLOKLANIR** (əbədi) |
| nil kanalı close | **panic: close of nil channel** |

**Yönümlü kanal parametrləri (compile-time qoruma!):**
```go
func printer(ch chan<- bool) { ch <- true }        // yalnız YAZ
func f2(out <-chan int, in chan<- int) { }         // oxu / yaz ayrılığı
```
Səhv istiqamətdə əməliyyat → compile xətası (istifadə olunmasa belə yoxlanılır).

**Race nümunəsi (kitabdan):** 5 goroutine eyni kanala yazır + receiver tərəfdən
close → `go run -race` DATA RACE + `panic: send on closed channel` aşkarlayır.
**Həll:** TƏK goroutine yazır və özü bağlayır:
```go
func printer(ch chan<- bool, times int) {
    for i := 0; i < times; i++ {
        ch <- true
    }
    close(ch)          // sender bağlayır — sıralı, racesiz
}
go printer(ch, 5)
for val := range ch { }   // close-dan sonra avtomatik çıxır
```

### 6. select — Çoxkanal Gözlənti
**Semantika:** bütün case-lər EYNİ ANDA yoxlanılır; hazır çoxsasa RANDOM seçilir; heç
biri hazır deyilsə bloklanır (default varsa default icra olunur). `select{}` əbədi gözləyir.

**Kitabdan kod nümunəsi:**
```go
for {
    select {
    case createNumber <- rand.Intn(max-min) + min:   // göndər
    case <-end:                                       // bit siqnalı
        fmt.Println("Ended!")
        // return YADDA SAXLA! — yoxsa funksiya bitmir
    case <-time.After(4 * time.Second):              // "clever default"
        fmt.Println("time.After()!")
        return
    }
}
```
**time.After(d)** — d vaxt sonra mesaj verən kanal; exit strategiyası kimi select-ə
timeout qatır.

### 7. Goroutine Timeout Texnikası
```go
// main daxilində:
select {
case res := <-c1:
    fmt.Println(res)
case <-time.After(time.Second):     // 1 san — Sleep 3 san olduğundan timeout
    fmt.Println("timeout c1")
}

// Ayrı funksiyada, parametrik müddət:
func timeout(t time.Duration) {
    temp := make(chan int)
    go func() {
        time.Sleep(5 * time.Second)
        defer close(temp)
    }()
    select {
    case <-temp:
        result <- false               // vaxtında bitdi
    case <-time.After(t):
        result <- true                // timeout
    }
}
```

### 8. Buffered Channel — Növbə/Semaphore
**Kitabdan kod nümunəsi (select+default ilə doluluq yoxlaması):**
```go
numbers := make(chan int, 5)
for i := 0; i < 10; i++ {
    select {
    case numbers <- i * i:
        fmt.Println("About to process", i)
    default:
        fmt.Print("No space for ", i, " ")   // doludur — drop et
    }
}
for {
    select {
    case num := <-numbers:
        fmt.Print("*", num, " ")
    default:
        fmt.Println("Nothing left to read!")
        return
    }
}
```

### 9. nil Channel — Select Branch Söndürmə
**Nədir:** nil kanal hər oxu/yazı BLOKLAYIR → select-də həmin case-i faktiki SÖNDÜRÜR.

**Kitabdan kod nümunəsi:**
```go
func add(c chan int) {
    sum := 0
    t := time.NewTimer(time.Second)
    for {
        select {
        case input := <-c:
            sum = sum + input
        case <-t.C:
            c = nil                 // vaxt bitdi → <-c case-i ölür
            fmt.Println(sum)
            wg.Done()
        }
    }
}
```

### 10. Worker Pool
**Nədir:** Məhdud sayda işçi goroutine-in növbədən iş götürməsi (Apache modeli).
Buffered channel-lə vasitəsilə.

**Kitabdan kod nümunəsi:**
```go
type Client struct { id, integer int }
type Result struct { job Client; square int }

var size = runtime.GOMAXPROCS(0)
var clients = make(chan Client, size)     // iş növbəsi
var data = make(chan Result, size)        // nəticə növbəsi

func worker(wg *sync.WaitGroup) {
    for c := range clients {              // işləri götür
        data <- Result{c, c.integer * c.integer}
    }
    wg.Done()
}

func create(n int) {
    for i := 0; i < n; i++ {
        clients <- Client{i, i}
    }
    close(clients)                        // işlərin sonu — worker-lər çıxır
}

// main: go create(nJobs); nWorkers ədəd worker goroutine;
// wg.Wait() sonra close(data); nəticə consumer for range data ilə oxuyur
```

### 11. Signal Channel ilə İcra Sırası
**Kitabdan kod nümunəsi (A→B→C→D zənciri):**
```go
func A(a, b chan struct{}) {
    <-a                 // öncəkini gözlə
    fmt.Println("A()!")
    close(b)            // növbətinə işarə
}
// main: x,y,z,w kanalları; close(x) prosesi BAŞLADIR
close(x)    // A açılır → A y-ni bağlayır → B açılır → ... → D
```
Sonuncu funksiya (D) kanal bağlamır → bir neçə dəfə çağırıla bilər; close yalnız BİR dəfə.

### 12. UNIX Siqnallar
**Kitabdan kod nümunəsi:**
```go
sigs := make(chan os.Signal, 1)     // tutum 1 — 1 siqnal anda kifayətdir
signal.Notify(sigs)                  // BÜTÜN tutula bilən siqnallar
// signal.Notify(sigs, syscall.SIGINT, syscall.SIGUSR1)  // konkretlər

go func() {
    for {
        sig := <-sigs
        switch sig {
        case syscall.SIGINT:
            fmt.Println("Execution time:", time.Since(start))
        case syscall.SIGUSR1:
            os.Exit(0)
        default:
            fmt.Println("Caught:", sig)
        }
    }
}()
for { time.Sleep(10 * time.Second); fmt.Print("+") }  // iş simulyasiyası
```
- **SIGKILL və SIGSTOP tutula BİLMƏZ** — kernel/privileged imtiyaz
- Linux-da SIGINFO yoxdur → SIGUSR1/SIGUSR2 istifadə et

### 13. sync.Mutex — Kritik Bölmə Qorunması
**Critical Section (kritik bölmə)** — paralel icra oluna bilməyən kod; mutex ilə qorunur.

**Kitabdan kod nümunəsi:**
```go
var m sync.Mutex
var v1 int

func change() {
    m.Lock()
    defer m.Unlock()          // həmişə defer ilə — unutma deadlock!
    time.Sleep(time.Second)
    v1 = v1 + 1
}

func read() int {
    m.Lock()                  // OXU da qorunmalı — yazma ilə yarış!
    a := v1
    defer m.Unlock()
    return a
}
```
**Unlock unutdulsa:** `fatal error: all goroutines are asleep - deadlock!`
**Qaydalar:** kritik bölmə eyni mutex-li başqa kritik bölməyə EMBED oluna bilməz; mutex-ləri
funksiyalar arasına SÖRÜŞDÜRMƏ — embed olub-olmadığını görməyi çətinləşdirir.

### 14. sync.RWMutex — Oxu Paralel
- **Lock/Unlock** — yazma: TƏK sahib
- **RLock/RUnlock** — oxu: ÇOXLU paralel oxucu
- Bütün oxucular bitənə qədər yazma kilidi gözləyir — qiymət/premium balansı
- RLock blokunda DƏYİŞİKLİK QADAĞDIR

```go
type secret struct {
    RWM      sync.RWMutex
    password string
}
func Change(pass string) { Password.RWM.Lock(); Password.password = pass; Password.RWM.Unlock() }
func show() { Password.RWM.RLock(); defer Password.RWM.RUnlock(); fmt.Println(Password.password) }
```

### 15. atomic — Kilidsiz Sadə Hallar
```go
func (c *atomCounter) Value() int64 { return atomic.LoadInt64(&c.val) }
atomic.AddInt64(&counter.val, 1)      // 100 goroutine × 4 = 400 düzgün
```
Sadə sayğaclar üçün mutexdən sadə; amma mutex daha versatildir. BÜTÜN oxu/yazılar
atomic funksiyaları ilə olmalıdır.

### 16. Monitor Goroutine — Communicating by Sharing Əksinə
**İdeya:** dəyişəni TƏK goroutine SAHİBLENİR; digərləri kanallarla sorğu göndərir —
"sharing by communicating" (Go fəlsəfəsi).

**Kitabdan kod nümunəsi:**
```go
var readValue = make(chan int)
var writeValue = make(chan int)

func set(newValue int) { writeValue <- newValue }
func read() int { return <-readValue }

func monitor() {
    var value int
    for {
        select {
        case newValue := <-writeValue:
            value = newValue
        case readValue <- value:
        }
    }
}

go monitor()            // TƏK instansiya — race MÜMKÜNSÜZ
```

### 17. Closure Tələsi — go Statement ilə
**Problem:**
```go
for i := 0; i <= 20; i++ {
    go func() { fmt.Print(i, " ") }()   // i CLOSURE-DİR!
}
// Output: 3 7 21 21 21 21 ... — çoxu SON dəyəri (21) çap edir; -race DATA RACE tapır
```
**Səbəb:** closured dəyişən goroutine İCRA ZAMANI qiymətlənir — scheduler gözləyəndə
loop qurtarır, i=21.

**2 həll:**
```go
i := i                              // shadowing — işləyir amma good practice DEYİL
go func(x int) { fmt.Print(x) }(i) // PARAMETR — TÖVSİYƏ OLUNAN
```
**Go 1.22-də bu problem aradan qalxdı** (hər iterasiya öz dəyişən nüsxəsi).

### 18. context — Ləğv Mexanizmi
**Context interfeysi:** Deadline(), Done(), Err(), Value() — hamısını implement etmək
lazım DEYİL; With* metodları uşaq kontekst + CancelFunc qaytarır.

**Kitabdan kod nümunəsi (3 variant):**
```go
// WithCancel — manual ləğv:
c1, cancel := context.WithCancel(context.Background())
defer cancel()
go func() { time.Sleep(4 * time.Second); cancel() }()
select {
case <-c1.Done():
    fmt.Println("f1() Done:", c1.Err())       // context canceled
case r := <-time.After(time.Duration(t) * time.Second):
    fmt.Println("f1():", r)
}

// WithTimeout — avtomatik vaxt bitməsi:
c2, cancel := context.WithTimeout(c2, time.Duration(t)*time.Second)

// WithDeadline — konkret an:
c3, cancel := context.WithDeadline(c3, time.Now().Add(...))
```
**Vacib:** `defer cancel()` — ləğv parent-in uşağa istinadını kəsir; əks halda
child goroutine-lər yığıla bilmir → memory leak (parent reference saxlamalıdır).

**WithCancelCause (Go 1.21) — səbəbli ləğv:**
```go
ctx, cancel := context.WithCancelCause(context.Background())
cancel(errors.New("Canceled by timeout"))
// select: case <-ctx.Done(): return context.Cause(ctx)   // "Canceled by timeout"
```
Həmçinin: WithTimeoutCause, WithDeadlineCause.

### 19. semaphore — golang.org/x/sync/semaphore
**Nədir:** Çəkili (weighted) resurs limiti — eyni anda maksimum N goroutine.

**Kitabdan kod nümunəsi:**
```go
var sem = semaphore.NewWeighted(int64(Workers))   // 4 token

ctx := context.TODO()
for i := range results {
    err = sem.Acquire(ctx, 1)          // token al — boşdursa BLOK
    if err != nil { break }
    go func(i int) {
        defer sem.Release(1)            // token qaytar
        results[i] = worker(i)          // hər goroutine ÖZ elementinə yazır — racesiz
    }(i)
}
// Clever trick — hamının bitməsini gözləmək:
err = sem.Acquire(ctx, int64(Workers))  // BÜTÜN tokenləri al → hamı bitməyib bloklanır
```
Kanal tələb etmir; nəticələr birbaşa slice-a (unikal index-ə) yazılır.

### 20. Statistika Tətbiqi — Kanalsız Concurrent Dizayn
**Kitabdan kod nümunəsi:**
```go
files = make(DFslice, len(os.Args))     // files[0] istifadəsiz — index uyğunluğu
var waitGroup sync.WaitGroup
for i := 1; i < len(os.Args); i++ {
    waitGroup.Add(1)
    go func(x int) {
        process(os.Args[x], x)          // hər goroutine ÖZ files[x] yerinə yazır
        defer waitGroup.Done()
    }(i)
}
waitGroup.Wait()
```
**Dizayn dərsi:** unikal index = race-siz paralel yazı; kanalsız, deadlock-sız, MINIMAL
dəyişikliklə concurrency. Benchmark: 9 faylda 3.27s → 1.24s (~3x).

## Əsas terminlər
- Goroutine — `go` ilə yaradılan yüngül icra vahidi
- m:n Scheduling — m goroutine n thread-də multiplex
- Work-Stealing — boş P-nin iş oğurlaması (continuation stealing)
- GOMAXPROCS — aktiv logical P həddi
- WaitGroup — Add/Done/Wait sayğaclı sinxronizasiya
- Channel — tipəmən kommunikasiya; buffered/unbuffered
- Signal Channel — yalnız siqnal üçün kanal (struct{})
- Race Condition —kilidsiz paralel yaddaş çıxışı
- Critical Section — mutex-lə qorunan kod bölməsi
- Mutex/RWMutex — eksklüziv / oxu-paralel kilid
- atomic — kilidsiz atomik əməliyyatlar
- Monitor Goroutine — datanı sahiblənən tək goroutine
- select — çoxkanal eyni-anda gözlənti; random seçim
- time.After — müddətli kanal (timeout)
- Context — ləğv oluna bilən əməliyyat konteksti
- Semaphore — çəkili goroutine limiti (x/sync)
- Pipeline — goroutine+kanal zənciri

## Praktik nətidə

(1) Goroutine-lər müstəqil deyil — ölmüş main = ölmüş proqram; Sleep yox, WaitGroup.
(2) Add go-dan ƏVVƏL, Done defer ilə; Add/Done balansı pozulsa panic/deadlock.
(3) Close-u YALNIZ sender edir; bağlı kanala yazma panic; nil kanal bloklayır — bunu
select branch söndürmək üçün İSTİFADƏ ET. (4) Kanal parametrlərini `chan<-`/`<-chan`
ilə yönləndir — compile vaxtı səyvlərindən qoru. (5) `go run -race` HƏMİŞƏ — race
bəzən gizlənir, müəyyən sıralarda partlayır. (6) select: hazır case-lər random;
heç biri hazır deyilsə blok; default = non-blocking. (7) time.After — timeout-un kanal
formaası; context isə ləğv + deadline + cause üçün tam həll. (8) defer cancel() —
context leak qarşısı. (9) Loop dəyişənini goroutine-ə PARAMETR ötür (Go 1.21-ə qədər
məcburi; 1.22+ belə daha aydın). (10) Mutex-də kritik bölmələri kiçik və lokal saxla;
oxu-çox halda RWMutex; sadə sayğacda atomic. (11) Kanal hər problemin cavabı deyil —
unikal index + WaitGroup (stats nümunəsi) deadlock-suz sadə paralellik verir. (12)
Concurrency dizayn qərarıdır, parallelism onun mükafatı; sıra-verilən kodu sona qədər
concurrent etmə.

## Mənbə
Pages: 321-385 (PDF 352-417)
