# Unit 7 — Concurrent Programming (Lesson 30-32)

## Bu unit nədən bəhs edir?

Goroutine-lər (go açar sözü, time sharing, ixtiyari sıra), channel-lər (send/receive, blocking, deadlock, nil channel, select + time.After timeout), pipeline-lar (source → filter → print, sentinel → close → range), concurrent state (race condition, mutex, worker pattern, command channel), Capstone: Mars rover-lərin paralel axtarışı.

**PDF səhifələr:** 268-299 (L30: 268-282, L31: 284-296, L32: 297-299)

## Əsas fikirlər

### 1. Goroutine başlatma (L30)
**Nədir:** Müstəqil işləyən tapşırıq — "gopher factory" metaforası: hər gopher öz işini görür, ünsiyyət kanallarla.

```go
func main() {
    go sleepyGopher()          // go açar sözü — sadəcə bu!
    time.Sleep(4 * time.Second) // main bitəndə BÜTÜN goroutine-lər dayandırılır
}
func sleepyGopher() {
    time.Sleep(3 * time.Second)
    fmt.Println("... snore ...")
}
```

**Bir neçə goroutine:**
```go
for i := 0; i < 5; i++ {
    go sleepyGopher(i)         // arqumentlər KOPYA ilə ötürülür
}
```
- **Time sharing:** məhdud processor-lar goroutine-lər arasında bölüşdürülür
- **Sıra QARANTIYA OLUNMUR** — hər run fərqli; həmişə ixtiyari sıra ehtimal et

### 2. Channel əsasları (L30)
**Metafora:** pneumatic tube (pochta borusu) — qoyursan → o biri tərəfdən çıxır.

```go
c := make(chan int)     // tipli kanal yarat

c <- 99                  // SEND: arrow kanala İSTİQAMƏTLİNƏ
r := <-c                  // RECEIVE: arrow kanaldan UZAQ
```
- **Send bloklanır** — alan goroutine yaranana qədər
- **Receive bloklanır** — göndərən yaranana qədər
- Kanal adi Go tipidir: dəyişən, parametr, struct sahəsi...

**Sleepy gophers + channel:**
```go
func main() {
    c := make(chan int)
    for i := 0; i < 5; i++ {
        go sleepyGopher(i, c)
    }
    for i := 0; i < 5; i++ {        // 5 mesaj gözlə — Sleep YOX!
        gopherID := <-c
        fmt.Println("gopher ", gopherID, " has finished sleeping")
    }
}
func sleepyGopher(id int, c chan int) {
    time.Sleep(3 * time.Second)
    fmt.Println("... ", id, " snore ...")
    c <- id                          // öz ID-sini göndər
}
```

### 3. select — çox kanal gözləmə (L30)
**Vəzifə:** müxtəlif kanallarda EYNİ ANDA gözlə — hansı hazır olsa onu icra et (switch-ə bənzər, amma kanal əməliyyatları).

```go
timeout := time.After(2 * time.Second)    // vaxt bitən kanal!
for i := 0; i < 5; i++ {
    select {
    case gopherID := <-c:
        fmt.Println("gopher ", gopherID, " has finished sleeping")
    case <-timeout:
        fmt.Println("my patience ran out")
        return
    }
}
```
- `time.After(d)` → d-dən sonra dəyər göndərən kanal (Go runtime goroutine-u)
- **Timeout pattern:** hər hansı əməliyyatı goroutine-ə qoy + kanal + select → nə isə vaxtla məhdudlaşdır
- Boş `select {}` → sonsuz gözlə (main-i sağ saxlamaq üçün)
- Nil kanal case-də → sadəcə heç vaxt hazır olmayacaq (seçim disabler!)

### 4. Blocking, deadlock, nil channel (L30)
**Blocked goroutine:** gözləyən goroutine **resurs istifadə etmir** (busy loop-dan fərqli — fan dönmür!). Sadəcə park olunub.

**Deadlock:**
```go
func main() {
    c := make(chan int)
    <-c      // kimse göndərməyəcək — deadlock → program asılır/crash
}
```

**Nil channel:** zero value = nil; nil kanala send/receive → əbədi bloklanır (panic YOX!); nil kanalı close → panic.

### 5. Pipeline — gopher assembly line (L30)
**Nədir:** Mərhədəli emal zənciri — hər worker yuxarıdan alır, işləyir, aşağıya göndərir. Böyük data stream-lərini az yaddaşla emal etmək üçün.

**Versiya 1 — sentinel dəyərlə (köhnə üsul):**
```go
func sourceGopher(downstream chan string) {
    for _, v := range []string{"hello world", "a bad apple", "goodbye all"} {
        downstream <- v
    }
    downstream <- ""                  // sentinel: bitiş!
}

func filterGopher(upstream, downstream chan string) {
    for {
        item := <-upstream
        if item == "" {
            downstream <- ""          // sentinel-i keçir
            return
        }
        if !strings.Contains(item, "bad") {
            downstream <- item
        }
    }
}

func printGopher(upstream chan string) {
    for {
        v := <-upstream
        if v == "" {
            return
        }
        fmt.Println(v)
    }
}

func main() {
    c0 := make(chan string)
    c1 := make(chan string)
    go sourceGopher(c0)
    go filterGopher(c0, c1)
    printGopher(c1)          // sonuncu eyni goroutine-da — proqram onun bitməsini gözləyir
}
```
**Problem:** boş string legal dəyərdirsə? → sentinel işləmir.

### 6. close + two-value + range (L30)
**close(c):** kanalın bağlanması — "artıq dəyər gəlməyəcək".
- Bağlı kanala yazma → panic
- Bağlı kanaldan oxuma → dərhal zero value (sonsuz loop təhlükəsi — yoxlama!")
- İki dəyərli oxuma: `v, ok := <-c` → ok=false = bağlanıb

**Versiya 2 — close ilə (düzgün):**
```go
func sourceGopher(downstream chan string) {
    for _, v := range []string{"hello world", "a bad apple", "goodbye all"} {
        downstream <- v
    }
    close(downstream)               // BAĞLA
}

func filterGopher(upstream, downstream chan string) {
    // range: kanal bağlanana qədər oxuyur!
    for item := range upstream {
        if !strings.Contains(item, "bad") {
            downstream <- item
        }
    }
    close(downstream)                 // zəncir boyu keçir
}

func printGopher(upstream chan string) {
    for v := range upstream {
        fmt.Println(v)
    }
}
```
**range + channel:** bağlanana qədər iterasiya — ən idiomatik pipeline oxunuşu.

### 7. Race condition — shared phone (L31)
**Metafora:** fabrikdə 1 telefon xətti, çox dəstək — eyni anda danışan gopher-lər sifarişi korlayır.

**Qayda:** 2+ goroutine eyni dəyərə yazırsa → **undefined behavior**. Oxu-oxu təhlükəsiz; yazı iştirak edirsə — qoruyun lazımdır. Race detector: `go run -race` — istifadə et və tapdığın race-i DÜZƏLT.

### 8. Mutex (L31)
**Metafora:** bankadakı metal token — götür → işlət → qaytar; token yoxdursa gözlə.

```go
var mu sync.Mutex
mu.Lock()
defer mu.Unlock()          // defer ilə — hər return-da açılır
// qorunan kod
```

**Web crawler nümunəsi:**
```go
// Visited tracks whether web pages have been visited.
// Its methods may be used concurrently from multiple goroutines.
type Visited struct {
    // mu guards the visited map.
    mu      sync.Mutex
    visited map[string]int
}

func (v *Visited) VisitLink(url string) int {
    v.mu.Lock()
    defer v.mu.Unlock()
    count := v.visited[url]
    count++
    v.visited[url] = count
    return count
}
```
- **Zero value = açıq mutex** — init lazım deyil
- Mutex struct SAHƏSİ kimi + qoruduğu data dərhal altında (şərh ilə assosiasiya açıq)
- **Lock/Unlock paket daxilində gizlədilir** — metodlar arxasında; xarici API-də mutex göstərilmir
- Konvensiya: metod concurrent-safe DİLSƏ — sənət olaraq dokumentasiya et

**Tələlər:**
- Kilid saxlarkən çox iş görmə → bloklama/dreadlock
- Eyni mutex-i kilidli ikən təkrar Lock → **deadlock** (əbədi blok)
- Kilidsiz Unlock → panic
- **Qaydalar:** (1) kiliddə kod sadə olsun; (2) hər shared state üçün TƏK mutex

### 9. Long-lived worker (L31)
**Nədir:** Öz-özünə işləyən, eyni zamanda xarici əmrlərə cavab verən daimi goroutine. Curiosity rover-in modulları belə işləyir (message-passing arxitekturası)!

**Skelet:**
```go
func worker() {
    for {
        select {
        // kanalları burada gözlə
        }
    }
}
go worker()
```

**Timer-lı worker:**
```go
func worker() {
    n := 0
    next := time.After(time.Second)
    for {
        select {
        case <-next:
            n++
            fmt.Println(n)
            next = time.After(time.Second)   // YENİ timer!
        }
    }
}
```

**Go vs event loop:** Node.js-də mərkəzi event loop var; Go-da HƏR worker öz event loop-udur — mərkəzi loop lazımsız.

### 10. RoverDriver — worker + command channel (L31)
**Kitabdan kod nümunəsi:**
```go
type command int
const (
    right = command(0)
    left  = command(1)
)

// RoverDriver drives a rover around the surface of Mars.
type RoverDriver struct {
    commandc chan command
}

func NewRoverDriver() *RoverDriver {
    r := &RoverDriver{
        commandc: make(chan command),
    }
    go r.drive()          // worker-i başlat
    return r
}

// drive is responsible for driving the rover. It
// is expected to be started in a goroutine.
func (r *RoverDriver) drive() {
    pos := image.Point{X: 0, Y: 0}          // image.Point — standart tip!
    direction := image.Point{X: 1, Y: 0}
    updateInterval := 250 * time.Millisecond
    nextMove := time.After(updateInterval)
    for {
        select {
        case c := <-r.commandc:              // xarici əmr
            switch c {
            case right:
                direction = image.Point{X: -direction.Y, Y: direction.X}
            case left:
                direction = image.Point{X: direction.Y, Y: -direction.X}
            }
            log.Printf("new direction %v", direction)
        case <-nextMove:                    // hərəkət taymeri
            pos = pos.Add(direction)
            log.Printf("moved to %v", pos)
            nextMove = time.After(updateInterval)
        }
    }
}

// Left turns the rover left (90º counterclockwise).
func (r *RoverDriver) Left() {
    r.commandc <- left
}

// Right turns the rover right (90º clockwise).
func (r *RoverDriver) Right() {
    r.commandc <- right
}

func main() {
    r := NewRoverDriver()
    time.Sleep(3 * time.Second)
    r.Left()
    time.Sleep(3 * time.Second)
    r.Right()
    time.Sleep(3 * time.Second)
}
```

**Sub-kod izahı:**
- Kanal **implementation detail** kimi metodlar arxasında gizlədilir (Left/Right sadə API)
- Worker direction dəyərinə TƏK sahibdir — metodlar birbaşa dəyişə bilmir → **race YOX, mutex LAZIM DEYİL!**
- `image.Point` — X/Y + Add metodu (standart kitabxanadan)
- select: komanda + taymer eyni anda gözlənilir

### 11. Capstone: Life on Mars (L32)
Tapşırıqlar:
1. **MarsGrid** — mutex ilə concurrent-safe grid; `Occupy(p) *Occupier` → nil (dolu/xaric) yaxud Occupier; `Move(p) bool` — uğursuzsa yerində qal
2. **Rover-lər** — grid-də gəzir; kənya/əngələ çatanda random dönüş
3. **Life axtarışı** — hər cell-də 0-1000 life ehtimalı; >900 → radio mesaj
4. **Buffer goroutine** — relay peyk üfüqdə deyil → mesajları slice-da buferlə
5. **Earth goroutine** — ara-sıra mesajları qəbul edir (koordinat + life dəyəri + rover adı)

Bütün Unit 6-7 bilikləri birləşir: struct, interface, pointer, mutex, worker, channel, select, pipeline.

## Concurrency alətləri xülasəsi
| Alət | Nə üçün | Nümunə |
|---|---|---|
| `go f()` | Tapşırığı paralelə burax | `go sleepyGopher(i, c)` |
| `c <- v` / `v = <-c` | Kanalla send/receive | gopher ID |
| `select` | Çox kanal gözlə | komanda + taymer |
| `time.After` | Timeout / periodik | 2s patience |
| `close(c)` | Stream sonu | source → filter |
| `for v := range c` | Bağlanana qədər oxu | pipeline worker |
| `sync.Mutex` | Shared state qoruması | Visited map |
| Worker (`for` + `select`) | Daimi müstəqil tapşırıq | RoverDriver.drive |

## Əsas terminlər
- Goroutine (coroutine/fiber/thread qohumu)
- Time Sharing (vaxt bölüşdürmə)
- Deterministik olmayan sıra (ixtiyari icra)
- Channel / Send / Receive (`<-`)
- Blocking (bloklanma — resurssız gözləmə)
- Deadlock
- Nil Channel (əbədi blok; close → panic)
- select (multi-kanal gözləmə)
- time.After (timeout/periodik kanal)
- Timeout Pattern
- Sentinel Value (bitiş siqnalı dəyər)
- close / Two-value Receive (`v, ok`)
- range Channel (bağlanana qədər)
- Pipeline (mərhələli emal)
- Race Condition / Undefined Behavior
- Race Detector (`-race`)
- Mutex (mutual exclusion) / Lock / Unlock
- Zero Value Mutex (açıq)
- Long-lived Worker
- Command Channel
- Event Loop vs Goroutine
- Message-passing Architecture (Curiosity rover)

## Praktik nəticə
- `go` + kanal — hər hansı kodu paralelə çevir; arqumentlər kopya gedir, loop dəyişənlərini parametr kimi ötür.
- Sleep ilə gözləmə prototip üçündür — production-da channel/sync ilə bitişi gözlə.
- Timeout lazımdırsa: goroutine + kanal + `select` + `time.After` — hər şeyə tətbiq olunur.
- Pipeline-lərdə sentinel əvəzinə `close` + `range` — hər mərhələ upstream-i bağlanana qədər oxuyub downstream-i bağlayır.
- Shared state + yazı = mutex; amma daha yaxşısı: dəyəri TƏK goroutine-ə ver (worker pattern) — mutex-ə heç ehtiyac qalmır.
- Mutex-dəkilə kodu sadə saxla; `defer Unlock` ilə qaytarmağa zəmanət ver.
- Kanalları API-də gizlət — metodlar (Left/Right) sadə interfeys, kanal implementation detail-dir.
- Worker pattern (`for { select {} }`) — hər hansı daimi tapşırıq: rover, poller, hardware controller.

## Mənbə
Pages: 268-299 (PDF), book pages 253-284 (Lesson 30-32)
