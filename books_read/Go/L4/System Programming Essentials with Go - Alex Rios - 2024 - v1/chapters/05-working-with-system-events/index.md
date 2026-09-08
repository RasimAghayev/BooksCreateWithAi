# Chapter 5 — Working with System Events (Sistem Hadisələri ilə İşləmək)

## Bu chapter nədən bəhs edir?
Siqnalların təbiətinə (sinxron/asinxron), os/signal paketinə, task scheduler qurulumuna (Job/Scheduler), Timer/Ticker ilə zaman idarəsinə, fayl monitorinqinə (Inotify + fsnotify), log rotation-a, proses timeout idarəsinə (context+select) və fayl kilidlərinə (fcntl/Flock — distributed lock manager).

## Əsas fikirlər

### 1. Siqnallar nədir
**Nədir:** Prosesə "hadisə baş verdi" bildirişi — software interrupt (proqram axınını kəsən software səbəbi; dəqiq vaxtını proqnozlaşdırmaq mümkün deyil).

**3 kateqoriya:**
| Mənbə | Nümunə |
|---|---|
| Hardware fault | SIGBUS, SIGFPE, SIGSEGV → Go-da runtime panic-ə çevrilir |
| İstifadəçi (terminal) | Ctrl+C → SIGINT; Ctrl+\ → SIGQUIT |
| Software | child prosesin bitməsi və s. |

**Qeyd-tutulmaz siqnallar:** SIGKILL və SIGSTOP — os/signal paketi tərəfindən tutula BİLMƏZ (prosesi heç nə qoruya bilməz).

### 2. os/signal — siqnal tutma
```go
func main() {
    signals := make(chan os.Signal, 1)      // siqnal kanalı (buferli!)
    done := make(chan struct{}, 1)           // çıxış siqnalı
    signal.Notify(signals, os.Interrupt)      // SIGINT qeydiyyatı

    go func() {
        for {
            s := <-signals
            switch s {
            case os.Interrupt:
                fmt.Println("INTERRUPT")
                done <- struct{}{}            // proqramı bitir
            default:
                fmt.Println("OTHER")
            }
        }
    }()
    fmt.Println("awaiting signal")
    <-done                                   // blokla
    fmt.Println("exiting")
}
```
**Pattern:** signal.Notify kanala yönləndirir → goroutine loop-da oxuyur → done kanalı ilə graceful exit.

**Siqnal idarəsinin 4 səbəbi:** graceful shutdown (resurs bağla, state saxla), resource management (SIGUSR1/2: log rotate, config reload), IPC (SIGSTOP/SIGCONT), emergency stop (SIGKILL/SIGABRT).

### 3. Task Scheduler qurulumu
```go
type Job func()                              // iş = funksiya

type Scheduler struct {
    jobQueue chan Job                        // buferli iş kanalı
}

func NewScheduler(size int) *Scheduler {
    return &Scheduler{jobQueue: make(chan Job, size)}
}

func (s *Scheduler) Start() {
    for job := range s.jobQueue {            // kanalı dinlə
        go job()                              // hər iş yeni goroutine-də
    }
}

func (s *Scheduler) Schedule(job Job, delay time.Duration) {
    go func() {
        time.Sleep(delay)                     // gecikmə gözlə
        s.jobQueue <- job                     // sonra növbəyə qoy
    }()
}

func main() {
    scheduler := NewScheduler(10)             // 10-luq bufer
    scheduler.Schedule(func() {
        fmt.Println("Job executed at", time.Now())
    }, 5*time.Second)
    go scheduler.Start()
    fmt.Scanln()                              // çıxış gözlə
}
```
**Niyə schedule:** effektivlik (off-peak icra), etibarlılıq (rotasyon backup), paralellik, proqnozlaşdırıla bilənlik (interval polling/report).

### 4. Timer + Ticker — zaman siqnalları
```go
ticker := time.NewTicker(1 * time.Second)   // hər saniyə
defer ticker.Stop()
timer := time.NewTimer(10 * time.Second)    // 10 saniyədən sonra bir dəfə
defer timer.Stop()

for {
    select {
    case tick := <-ticker.C:                 // periodic hadisə
        fmt.Println("Tick at", tick)
    case <-timer.C:                          // bir dəfəlik son
        fmt.Println("Timer expired")
        return
    }
}
```
**select ilə birgə idarə** — iki müxtəlif zaman mənbəyi vahid loop-da.

### 5. Inotify — Linux kernel fayl monitorinqi
**Nədir:** Kernel subsistemi — fayl/kataloq hadisələri (create/modify/delete/move) haqqında bildiriş.

**X/sys/unix ilə aşağı səviyyəli:**
```go
fd, err := unix.InotifyInit()                // inotify nümunəsi
defer unix.Close(fd)

watchDescriptor, err := unix.InotifyAddWatch(fd, watchPath,
    unix.IN_MODIFY|unix.IN_CREATE|unix.IN_DELETE)  // hansı hadisələr
defer unix.InotifyRmWatch(fd, uint32(watchDescriptor))

// Event loop:
const bufferSize = (unix.SizeofInotifyEvent + unix.NAME_MAX + 1)
buf := make([]byte, bufferSize)
for {
    n, err := unix.Read(fd, buf[:])
    // buf-dakı hadisələri parse et:
    var offset uint32
    for offset < uint32(n) {
        event := (*unix.InotifyEvent)(unsafe.Pointer(&buf[offset]))
        nameBytes := buf[offset+unix.SizeofInotifyEvent :
            offset+unix.SizeofInotifyEvent+uint32(event.Len)]
        name := string(nameBytes[:clen(nameBytes)])   // NUL trim
        fmt.Printf("Event: %s/%s\n", watchPath, name)
        offset += unix.SizeofInotifyEvent + uint32(event.Len)
    }
}
```
**Dərinlik:** unsafe.Pointer ilə byte → struct çevirmə; clen() NUL bayt kəsici; bir Read → çox hadisə.

### 6. fsnotify — cross-platform alternativ
**Üstünlükləri:** platforma abstraksiyası (Windows/macOS/Linux), sadə API, corner-case idarəsi, community saxlanması.

```go
import "github.com/fsnotify/fsnotify"

watcher, err := fsnotify.NewWatcher()
defer watcher.Close()
watcher.Add(watchPath)

go func() {
    for {
        select {
        case event := <-watcher.Events:     // hadisə kanalı
            fmt.Printf("Event: %s\n", event.Name)
        case err := <-watcher.Errors:        // xəta kanalı
            log.Println("Error:", err)
        }
    }
}()

signalCh := make(chan os.Signal, 1)
signal.Notify(signalCh, os.Interrupt, syscall.SIGINT)
<-signalCh                                 // Ctrl+C gözlə
```
**Seçim qaydası:** Portativlik/sadəlik → fsnotify; spesifik funksiya/öyrənmə → xam inotify.

### 7. Log rotation — fsnotify əsaslı
**Ssenari:** Log faylı max həddi aşanda adını dəyiş + yenisini yarat.

```go
const (
    logFilePath = "your_log_file.log"
    maxFileSize  = 1024 * 1024 * 10        // 10 MB
)

var mu sync.Mutex                          // rotasiya yarışı qoruması

go func() {
    for {
        select {
        case event, ok := <-watcher.Events:
            if !ok { return }
            if event.Op&fsnotify.Write == fsnotify.Write {
                fi, err := os.Stat(logFilePath)
                if err != nil { continue }
                if fi.Size() >= maxFileSize {
                    mu.Lock()
                    rotateLogFile()         // tək rotasiya!
                    mu.Unlock()
                }
            }
        case err, ok := <-watcher.Errors:
            if !ok { return }
        }
    }
}()

func rotateLogFile() {
    closeLogFile()                                           // 1. Bağla
    timestamp := time.Now().Format("20060102150405")
    newLogFilePath := fmt.Sprintf("your_log_file_%s.log", timestamp)
    os.Rename(logFilePath, newLogFilePath)                   // 2. Adını dəyiş
    createLogFile()                                          // 3. Yenisini yarat
}
```
**İstifadə sahələri:** sistem logları, backuplar, compliance, IoT sensor data, web server access logları.

### 8. Proses icrası + timeout
```go
cmd := exec.Command("sleep", "2")
timeout := 3 * time.Second

ctx, cancel := context.WithTimeout(context.Background(), timeout)
defer cancel()                               // niyyəti self-documenting edir

cmd.Start()
done := make(chan error, 1)
go func() { done <- cmd.Wait() }()           // gözləmə paraleldə

select {
case <-ctx.Done():                          // timeout bitdi
    cmd.Process.Kill()
    fmt.Println("Process killed as timeout reached")
case err := <-done:                         // vaxtında bitdi
    if err != nil { fmt.Println("finished with error") }
}
```
**Niyə timeout:** resurs idarəsi (asılı proses), etibarlılıq (time-sensitive əməliyyatlar), deadlock qarşısı.

### 9. Fayl kilidləri — fcntl/Flock (distributed lock manager)
**Lock növləri:**
| | Advisory | Mandatory |
|---|---|---|
| Tətbiqedici | Proseslər özü əməkdaşlıq edir | OS məcbur edir |
| Əməkdaşlıqsız proses | Kilidi İQNORR edə bilər | Bloklanır |

```go
file, err := os.Open("yourfile.txt")
defer file.Close()

lock := syscall.Flock_t{
    Type:   syscall.F_WRLCK,      // yazma kilidi (F_RDLCK = oxu)
    Whence: io.SeekStart,
    Start:  0,                     // faylın başından
    Len:    0,                     // 0 = bütün fayl
}
syscall.FcntlFlock(file.Fd(), syscall.F_SETLK, &lock)   // kilidi qoy

// Burax:
lock.Type = syscall.F_UNLCK
syscall.FcntlFlock(file.Fd(), syscall.F_SETLK, &lock)
```
**İstifadə ssenariləri:** data korlanmasının qarşısı, DB tək-instans, fayl sinxronizasiyası, resurs təyinatı (cluster: hansı maşın istifadə edir), message queue tək-consumer, keş/shared memory koordinasiyası, editor tək-istifadəçi, backup bütövlüyü.

**Xatırlatma:** Kilidlər advisory-dir — əməkdaşlıq edən proseslər üçün işləyir, OS məcbur etmir.

## Əsas terminlər
- Signal — software interrupt (bildiriş)
- Sinxron (SIGBUS/FPE/SEGV → panic) / asinxron (SIGINT/QUIT) siqnallar
- signal.Notify — siqnal → kanal yönləndirməsi
- Graceful shutdown — təmiz çıxış
- Job/Scheduler — funksiya növbəsi + gecikməli icra
- Timer (bir dəfə) / Ticker (periodik)
- Inotify: InotifyInit/AddWatch/RmWatch, IN_MODIFY/CREATE/DELETE
- unsafe.Pointer — byte → struct bərpa
- fsnotify — cross-platform watcher (Events/Errors kanalları)
- Log rotation — timestamp adı + Rename + yeni fayl
- context.WithTimeout + select + cmd.Process.Kill — timeout öldürücü
- fcntl / Flock_t / F_SETLK / F_WRLCK / F_UNLCK
- Advisory vs mandatory lock

## Praktik nəticə
1. Hər uzun ömürlü Go proqramında signal.Notify + graceful shutdown qur — Ctrl+C-də DB bağlantısı/resurs itmir.
2. Gecikməli işlər üçün minimal scheduler: chan Job + sleep-goroutine; periodik üçün Ticker; select ilə birləşdir.
3. Fayl dəyişikliyi izləmək: fsnotify (portativ) ilə başla; yalnız spesifik ehtiyacda xam inotify.
4. Log rotasiyasında mutex — eyni anda 2 write-event 2 rotasiya başlada bilər.
5. Xarici proseslər üçün HƏMİŞƏ context+timeout+Kill pattern-i — asılı `sleep 9999` serveri bloklamır.
6. Çoxproses fayl erişimində Flock advisory kilidi — amma unutma: əməkdaşlıq etməyən proses üçün müdafiə DEYİL.

## Mənbə
Pages: 83-103 (PDF səh. 104-125)
