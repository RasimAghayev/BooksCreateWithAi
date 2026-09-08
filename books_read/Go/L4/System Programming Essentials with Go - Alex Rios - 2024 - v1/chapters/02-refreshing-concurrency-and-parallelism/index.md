# Chapter 2 — Refreshing Concurrency and Parallelism (Paralellik və Paralelizmə Təzə Baxış)

## Bu chapter nədən bəhs edir?
Goroutine-lərə (WaitGroup ilə idarə), data race-lərə (warehouse analogiyası + race detector), atomic/mutex sinxronizasiyasına, kanallara (unbuffered/buffered, deadlock halları, range, close), delivery guarantee + latency seçiminə və kanal 3 vəziyyətinə (nil/open/closed).

## Əsas fikirlər

### 1. Goroutine-lər
**Nədir:** Go scheduler tərəfindən müstəqil icra üçün yaradılan funksiyalar. `go` keyword — kompleks scheduler alqoritmi sadə interfeyslə.

**async/await-dən fərqi:** Funksiya imzası DƏYİŞMİR — gözlənilən (awaitable) elan etmək, xüsusi çağırış notasyonu lazım DEYİL.

**Say hello nümunələri:**
```go
func say(s string) {
    for i := 1; i < 5; i++ {
        time.Sleep(500 * time.Millisecond)
        fmt.Println(s)
    }
}
func main() {
    go say("hello")    // goroutine — paralel
    say("world")        // main kontekstində — seqvensial
}
// Nəticə: hello/world alternativ
```
**Tələ:** `say("hello"); go say("world")` — main-in sonu goroutine-in bitməsini gözləmir → "world" çap olunmur! Görünür yalnız 4 dəfə "hello".

### 2. WaitGroup — bitmə gözləyici
```go
func main() {
    wg := sync.WaitGroup{}      // zero-value istifadəyə hazırdır
    wg.Add(2)                    // 2 goroutine gözləyəcəyik
    go say("world", &wg)
    go say("hello", &wg)
    wg.Wait()                    // hamısı bitənə qədər blokla
}
func say(s string, wg *sync.WaitGroup) {
    defer wg.Done()              // bitəndə sayğacı azalt
    // ...
}
```
**API:** Add(n) (qeydiyyat) → Done() (bitmə, defer ilə) → Wait() (bloklama).

### 3. Data race — warehouse analogiyası
**Ssenari:** 2 işçi (worker), hərəsi 1000 element qablaşdırır; totalItems ümumi sayğac — sinxronizasiyasız!

**Race-li kod (PIT):**
```go
func PackItems(totalItems int) int {
    const workers = 2
    const itemsPerWorker = 1000
    var wg sync.WaitGroup
    itemsPacked := 0
    for i := 0; i < workers; i++ {
        wg.Add(1)
        go func(workerID int) {
            defer wg.Done()
            for j := 0; j < itemsPerWorker; j++ {
                itemsPacked = totalItems    // KOYNU oxu (race!)
                itemsPacked++
                totalItems = itemsPacked    // KOYNU yaz (race!)
            }
        }(i)
    }
    wg.Wait()
    return totalItems
}
```
**Nondeterminizm:** İcra tez-tez "düzgün" (2000) verir — müəllifin testində 16421-ci cəhddə 1170 tapıldı! **runtime.Gosched()** ("pausa mənə yaxşı olar" hiyesi) ilə "noise" emulyasiyası → 4-cü cəhddə xəta.

**Race detector:**
```bash
go test -race
# WARNING: DATA RACE
# Read at 0x... by goroutine 9 / Previous write at 0x... by goroutine 8
```

### 4. Atomic əməliyyatlar — sync/atomic
**Nədir:** Tək əməliyyat səviyyəsində sinxronizasiya — load/store/add/CAS. Yalnız int32/int64/uint32/uint64/uintptr/float32/float64.

```go
import "sync/atomic"
for j := 0; j < itemsPerWorker; j++ {
    atomic.AddInt32(&totalItems, int32(itemsPacked))
}
// race detector: təmiz!
```
**Hüdud:** Tək əməliyyat üçün ideal; KOD BLOKU üçün mutex.

### 5. Mutexes — kod bloku qoruması
```go
func PackItems(m *sync.Mutex, totalItems int) int {
    // ...
    go func(workerID int) {
        defer wg.Done()
        for j := 0; j < itemsPerWorker; j++ {
            m.Lock()
            itemsPacked := totalItems
            itemsPacked++
            totalItems = itemsPacked
            m.Unlock()
        }
    }(i)
    // ...
}
```
**Bounce metaforası:** Mutex = rəqs meydançasının nəzarətçisi — bir goroutine kritik bölgədə, qalanları növbədə.

**Lock qranulluğu tələsi:**
```go
// Çox lock (hər sətir üçün ayrıca):
m.Lock(); itemsPacked = totalItems; m.Unlock()
m.Lock(); itemsPacked++;              m.Unlock()
m.Lock(); totalItems = itemsPacked;   m.Unlock()
```
Benchmark: 1 böyük kritik blok 32629 ns/op; çox-kiçik bloklar 91246 ns/op — **~64% yavaş!**
**Dərs:** Qranulluğu incəltmək həmişə yaxşı deyil — lock/unlock overhead-i qazancı yeyir.

### 6. Unbuffered kanallar — "trust fall"
**Qayda:** Göndərən və qəbul edən EYNİ ANDA hazır olmalıdır.

**3 deadlock nümunəsi:**
```go
c := make(chan string)
c <- "message"           // 1) göndərici var, qəbul yox → [chan send] deadlock

c := make(chan string)
fmt.Println(<-c)         // 2) qəbul var, göndərici yox → [chan receive] deadlock

c <- "message"; fmt.Println(<-c)  // 3) seqvensial! send-in özü gözləyir → deadlock
```
**Həll — goroutine ilə sinxron hazırlıq:**
```go
balls := make(chan string)
go throwBalls("red", balls)          // göndərici paraleldə
fmt.Println(<-balls, "received!")    // qəbul hazırdır

func throwBalls(color string, balls chan string) {
    fmt.Printf("throwing the %s ball\n", color)
    balls <- color
}
```

### 7. range + close — kanal üzərində iterasiya
```go
for color := range balls {    // close olunana qədər gözləyir
    fmt.Println(color, "received!")
}
```
**Problemlər:** range kanal BAĞLANMAYINCƏ dayanmır → deadlock; close-u yanlış yerdə çağırmaq → mesajlar itir.

**Tam həll (3 goroutine pattern-i):**
```go
balls := make(chan string)
wg := sync.WaitGroup{}
wg.Add(2)
go func() { defer wg.Done(); throwBalls("red", balls) }()
go func() { defer wg.Done(); throwBalls("green", balls) }()
go func() {
    wg.Wait()      // bütün göndəricilər bitib...
    close(balls)    // ...onda bağla
}()
for color := range balls {
    fmt.Println(color, "received!")
}
```
**Nüans:** throwBalls imzasının "korlanması"nm — anonim funksiyalar funksiyanı concurrency-dən xəbərsiz saxlayır.

### 8. Buffered kanallar — clown car
**Nədir:** `make(chan int, 3)` — 3 yerlik bufer; göndərici qəbulçunu gözləmədən yazır (yer varsa).

**Clown car nümunəsi (select+default ilə):**
```go
clownChannel := make(chan int, 3)   // maşın 3 yer
clowns := 5

// Sürücü (qəbulçu):
go func() {
    defer close(clownChannel)
    for clownID := range clownChannel {
        fmt.Printf("Driver: Drove the car with Balloon %d\n", clownID)
        time.Sleep(500 * time.Millisecond)
    }
}()

// Palçığlar (göndəricilər):
for clown := 1; clown <= clowns; clown++ {
    wg.Add(1)
    go func(clownID int) {
        defer wg.Done()
        select {
        case clownChannel <- clownID:   // yer varsa — min
        default:                          // doludur — gözləMƏ
            fmt.Printf("Clown %d: Oops, the car is full!\n", clownID)
        }
    }(clown)
}
```
**select+default = non-blocking send** — dolu kanalda gözləmə əvəzinə alternativ yol.

### 9. Delivery guarantee və latency seçimi
| | Unbuffered | Buffered |
|---|---|---|
| **Delivery** | HƏMİŞƏ zəmanətli (alıcı hazırdır) | ZƏMANƏT YOX — alıcı oxumaya bilər |
| **Latency** | Yüksək (sinxron bloklama) | Aşağı (dekoruplə) |
| **Nə vaxt** | data bütövlüyü kritik, 1:1, load-balancing | asinxron, contention azaltma, deadlock qarşısı, batch/pipeline |

### 10. Kanalın 3 vəziyyəti — tam cədvəl
| Əməliyyat | nil | open (empty) | open (data var) | closed |
|---|---|---|---|---|
| **Oxu** | sonsuz blok | blok (data gözlə) | data qaytar | zero value + false |
| **Yaz** | sonsuz blok | uğur | (doludursa blok) | **PANIC** |
| **Close** | panic | uğur (drain-ə icazə) | uğur (qalan oxunar) | **PANIC** |

(Write-only kanala oxu / read-only-ə yaz / close(only-read) = compile xətası.)

### 11. Signaling pattern-i
```go
signalChannel := make(chan bool)
// Goroutine 1: <-signalChannel — gözlə
// Goroutine 2: signalChannel <- true — siqnal ver
```
Producer-consumer, fan-out/fan-in pattern-lərinin əsası.

### 12. Mutex yoxsa channel? (praqmatizm)
**Kanal seç:** data sahibliyinin ötürülməsi, iş vahidlərinin paylanması, asinxron nəticə kommunikasiyası.
**Mutex seç:** keşlər (caches), shared state.
**Qızıl qayda:** Praqmatizm — mutex oxunaqlıdırsa, TWO DƏFƏ DÜŞÜNMƏ, mutex-ə get!

## Əsas terminlər
- Goroutine — `go` ilə başladılan yüngül axın
- WaitGroup — Add/Done/Wait üçlüyü
- Data race — sinxronizasiyasız paylaşılan dəyişən erişimi
- runtime.Gosched() — scheduler-ə "pausa ver" hiyesi
- Race detector — `go test -race`
- sync/atomic — Add/Load/Store/CAS (tək əməliyyat)
- sync.Mutex — Lock/Unlock (kod bloku)
- Unbuffered/buffered channel — tutumsuz/tumlu kanal
- Deadlock — qarşılıqlı gözləmə (chan send/receive)
- range + close — kanal iterasiyası + bağlanma
- select+default — non-blocking göndər/qəbul
- Delivery guarantee — çatdırılma zəmanəti (yalnız unbuffered)
- nil/open/closed — kanal 3 vəziyyəti

## Praktik nəticə
1. Goroutine-ləri HƏMİŞƏ WaitGroup (və ya kanal) ilə sinxronlaşdır — main bitməsi goroutine-ləri öldürür.
2. Race şübhəsi = dərhal `go test -race`; nondeterminizm xətanı minlərlə cəhddən sonra göstərir.
3. Tək sayğac əməliyyatı üçün atomic; çoxsətirli kritik blok üçün mutex — qranulluğu incələmə (64% itki!).
4. Unbuffered = zəmanət; buffered = məhsuldarlıq — data itirsem olar-mı sualına cavab ver.
5. Kanal üzərində range = close mütləq; close-u yalnız bütün göndəricilər bitəndən SONRA (WaitGroup+3-cü goroutine pattern).
6. Close-a yazış PANIC-dir; close-un close-u PANIC — bağlanmış kanalı bir də bağlama.

## Mənbə
Pages: 13-36 (PDF səh. 34-57)
