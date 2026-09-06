# Chapter 9 — Concurrency: Practice (#61-#74)

## Bu chapter nədən bəhs edir?

Bu chapter konkurrentlik praktikasının 14 səhvini əhatə edir: yanlış context propagasiyası, goroutine leak-lər, loop dəyişəni tələləri, select-in random seçimi, notification channel-lər, nil channel gücü, channel ölçüsü, string formatting yan təsirləri (etcd data race, deadlock), append data race-ləri, mutex/slice/map sərhədləri, sync.WaitGroup, sync.Cond, errgroup və sync tiplərinin kopyalanması.

---

## Əsas fikirlər

### #61: Propagating an inappropriate context (Uyunsuz context propagasiyası)

**Ssenari:** HTTP handler cavabı yazır + Kafka publish-i asinxron goroutine-də edir:

```go
func handler(w http.ResponseWriter, r *http.Request) {
    response, err := doSomeTask(r.Context(), r)
    if err != nil { /* 500 */ return }
    go func() {
        err := publish(r.Context(), response)   // HTTP context-i götürür!
        // ...
    }()
    writeResponse(response)
}
```

**Problem:** HTTP request context-i 3 halda ləğv olunur:
1. Klient bağlantısı qapanır
2. HTTP/2 request ləğv olunur
3. **Response yazıldıqdan sonra** ← bura təhlükədir!

Race condition: response Kafka publish-dən əvvəl/daha yazılsa → context ləğv olunur → publish mesajı itirir.

**Həllər:**

```go
// 1. Boş context (dəyərlər itir):
err := publish(context.Background(), response)

// 2. Custom detach context — cancel siqnalı YOX, dəyərlər SAXLI:
type detach struct{ ctx context.Context }
func (d detach) Deadline() (time.Time, bool) { return time.Time{}, false }
func (d detach) Done() <-chan struct{}       { return nil }
func (d detach) Err() error                  { return nil }
func (d detach) Value(key any) any           { return d.ctx.Value(key) }

err := publish(detach{ctx: r.Context()}, response)
```

`context.Context` interfeysi: `Deadline / Done / Err / Value` — Done+Err-ni söndürərək ləğvsiz, Value-ni parentdən miras alan context.

---

### #62: Starting a goroutine without knowing when to stop it (Dayanma planı olmayan goroutine)

**Goroutine leak nə deməkdir:** Minimum 2 KB stack (böyüyə bilər: 1 GB-a qədər 64-bit) + heap-də saxladığı variable istinadları + tutduğu resurslar (HTTP/DB bağlantıları, fayllar, socketlər).

```go
ch := foo()
go func() {
    for v := range ch { /* ... */ }   // ch bağlananda çıxır — amma nə vaxt??
}()
```

`foo`-nun `ch`-ni heç vaxt bağlamaması = leak. **Hər goroutine-in exit nöqtəsi bəlli olmalıdır.**

**Watcher nümunəsi — 3 versiya:**

```go
// PİS — dayanma yolu yox:
func newWatcher() { w := watcher{}; go w.watch() }

// YARIMÇIÇA — context siqnalı, amma bağlanacağına ZƏMANƏT YOX:
ctx, cancel := context.WithCancel(context.Background())
defer cancel()
newWatcher(ctx)   // watch ctx-ləğvini görsə də main çıxışı gözləmir

// DÜZGÜN — blocking close + defer:
func main() {
    w := newWatcher()
    defer w.close()          // main çıxmağında resurslar QARANTİYALI bağlanır
    // Run the application
}
func newWatcher() watcher {
    w := watcher{}
    go w.watch()
    return w
}
func (w watcher) close() { /* Close the resources */ }
```

**Prinsip:** Goroutine resursdur; siqnal göndərmək YOX — bloklanaraq bağlanmanı GÖZLƏ. Ömrü aplikasiyaya bağlı goroutine-lərdə çıxışdan əvvəl Wait.

---

### #63: Not being careful with goroutines and loop variables (Loop dəyişəni tələsi)

```go
s := []int{1, 2, 3}
for _, i := range s {
    go func() {
        fmt.Print(i)   // closure — i-yi referans edir
    }()
}
// Çıxış qeyri-deterministik: "233", "333"... "123" gözlənilirdi!
```

**Səbəb:** Bütün closure goroutine-ləri EYNİ `i` dəyişəninə istinad edir — goroutine İCRASI zamanı cari dəyəri oxuyur (yarananda YOX).

**Həllər:**

```go
// 1. Lokal kopya:
for _, i := range s {
    val := i
    go func() { fmt.Print(val) }()
}

// 2. Parametrli funksiya (closure deyil):
for _, i := range s {
    go func(val int) { fmt.Print(val) }(i)
}
```

**Qayda:** Goroutine closure loop dəyişəninə bədəndən KƏNARDAN istinad edirsə — problem. (Go 1.22-dən əvvəl; yeni versiyalarda loop var semantikası dəyişib olsa da, köhnə kod bazalarında aktual.)

---

### #64: Expecting deterministic behavior using select and channels

```go
// messageCh prioritetli təsəvvür olunur (ilk case):
for {
    select {
    case v := <-messageCh:
        fmt.Println(v)
    case <-disconnectCh:
        fmt.Println("disconnection, return")
        return
    }
}
```

**Spesifikasiya:** *"If one or more of the communications can proceed, a single one that can proceed is chosen via a uniform pseudo-random selection."*

**Select ≠ switch:** İlk uyğun case YOX — multiple ready olduqda RANDOM seçilir. (Səbəb: starvation qorunması — sürətli sender bir channel-i dominant etməsin.) Nəticə: 10 mesajın 5-i çap olundu, sonra disconnect.

**Həllər:**

**Tək producer:**
1. `messageCh`-ni unbuffered et — sender receiver hazır olana qədər bloklanır → bütün mesajlar disconnect-dən əvvəl çatır.
2. Tək channel istifadə et — struct (mesaj YOXSA disconnect) tipli; channel FİFO qarantisi verir → disconnect ən sonda çatır.

**Çoxlu producer:** Inner select + default pattern:

```go
for {
    select {
    case v := <-messageCh:
        fmt.Println(v)
    case <-disconnectCh:
        for {
            select {
            case v := <-messageCh:     // qalan mesajlar
                fmt.Println(v)
            default:                    // mesaj qalmadıqda
                fmt.Println("disconnection, return")
                return
            }
        }
    }
}
```

`default` yalnız heç bir case hazır deyilsə seçilir → disconnect gələndə messageCh-dəkilər hamısı emal olunana qədər dövr.

---

### #65: Not using notification channels (Notification channel-lər)

```go
// PİS — false nə deməkdir?? yenidən qoşulma? bağlanmama? tezliyi?:
disconnectCh := make(chan bool)

// DÜZGÜN — datasız siqnal:
disconnectCh := make(chan struct{})
```

**Empty struct = 0 bayt:** `unsafe.Sizeof(struct{}{})` → 0. (`interface{}` isə 8/16 bayt — istifadə ETMƏ.) Digər use case: hash set — `map[K]struct{}`.

**Prinsip:** Mesaj məzmunu deyil, yalnız "mesaj gəldi" faktı daşıyan channel → `chan struct{}` — receiver-dən dəyər gözləməməsi gözlənilir.

---

### #66: Not using nil channels (Nil channel gücü)

**Nil channel davranışı:** Göndərmə/qəbul = ƏBƏDİ bloklama (panic YOX).

**Merge nümunəsi — inkişaf yolu:**

```go
// V1 — sequensial: ch1 bağlanana qədər ch2 oxunmur — uyğunsuz
// V2 — for/select: ch1 bağlansa receive 0 qaytarır → SONSUZ "received: 0"!
```

**Qapalı channel qəbulu = non-blocking, zero value:** `v, open := <-ch1` — `open=false` bağlanma siqnalını ayırır.

```go
// V3 — Boolean state machine: busy-waiting CPU yandırır (qapalı case HƏMİŞƏ hazır olur)
// V4 — NİL CHANNEL həlli:
func merge(ch1, ch2 <-chan int) <-chan int {
    ch := make(chan int, 1)
    go func() {
        for ch1 != nil || ch2 != nil {
            select {
            case v, open := <-ch1:
                if !open { ch1 = nil; break }   // case-i select-dən SÖNDÜR
                ch <- v
            case v, open := <-ch2:
                if !open { ch2 = nil; break }
                ch <- v
            }
        }
        close(ch)
    }()
    return ch
}
```

**Mahiyyət:** Qapalı channel-i nil-lə → select o case-i görmür (nil = əbədi blok) → yalnız açıq channel-lər gözlənilir; busy-wait YOX.

---

### #67: Being puzzled about channel size (Channel ölçüsü)

**Unbuffered (synchronous channel):** sender receiver qəbul edənə qədər bloklar — **GÜCLÜ SİNXRONİZASİYA** (iki goroutine məlum vəziyyətdədir).

**Buffered:** sender buffer dolana qədər bloklamır — yalnız causality qarantisi var ("qəbul göndərmədən əvvəl ola bilməz"). Obscure deadlock-lara səbəb ola bilər (unbuffered-da dərhal görünənlər).

**Qərar ağacı:**
1. Sinxronizasiya lazımdır → unbuffered (məcbur).
2. Buffered lazımdırsa → **default ölçü 1.**
3. Başqa dəyər yalnız səbəblə:
   - Worker pooling → ölçü = goroutine sayı
   - Rate limiting → ölçü = limit
4. Magic number-lar (`make(chan int, 40)`) — səbəb + komment (benchmark/testlərlə doğrulanmış) şərtsiz qadağadır.

**LMAX Disruptor kağızından:** *"Queues are typically always close to full or close to empty..."* — istehlak/istehsal temp fərqi üzündən sabit dəqiq ölçü tapmaq nadir mümkündür. CPU contention vs memory trade-off.

---

### #68: Forgetting about possible side effects with string formatting (String formatting yan təsirləri)

**Hal 1 — etcd data race (real bug, PR 7816):**

```go
func (w *watcher) Watch(ctx context.Context, ...) WatchChan {
    ctxKey := fmt.Sprintf("%v", ctx)   // context-i formatlaşdırır!
    // ...
}
```

`fmt.Sprintf("%v", ctx)` — context-in BÜTÜN dəyərlərini oxuyur (Value zənciri). Context mutable dəyər daşıyırsa (pointer sahəli struct) → bir goroutine dəyəri YAZIR, Watch isə OXUYUR → **data race**. Düzəliş: `fmt.Sprintf` əvəzinə immutable dəyərdən custom `streamKeyFromCtx` funksiyası.

**Hal 2 — Stringer deadlock:**

```go
type Customer struct {
    mutex sync.RWMutex
    id    string
    age   int
}
func (c *Customer) UpdateAge(age int) error {
    c.mutex.Lock()
    defer c.mutex.Unlock()
    if age < 0 {
        return fmt.Errorf("age should be positive for customer %v", c)  // %v → String() → RLock!
    }
    c.age = age
    return nil
}
func (c *Customer) String() string {
    c.mutex.RLock()                       // UpdateAge Lock saxlayır → DEADLOCK!
    defer c.mutex.RUnlock()
    return fmt.Sprintf("id %s, age %d", c.id, c.age)
}
// fatal error: all goroutines are asleep - deadlock!
```

**Həllər:**
1. **Mutex scope-unu daralt** — yoxlama LOCK-DAN ƏVVƏL:

```go
func (c *Customer) UpdateAge(age int) error {
    if age < 0 {
        return fmt.Errorf("age should be positive for customer %v", c)  // lock YOXDUR → String() işləyir
    }
    c.mutex.Lock()
    defer c.mutex.Unlock()
    c.age = age
    return nil
}
```

2. Lock saxlarkən String()-i çağırmayan format: `fmt.Errorf("... customer id %s", c.id)` — birbaşa sahə.

**Moral:** Concurrent kodda formatlaşdırma gizli funksiya çağırışlarıdır (`Stringer`) — data race və deadlock mənbəyi. Negative-age testi yazılmalı idi!

---

### #69: Creating data races with append (Append data race)

**Append race-free DOLU slice-də:** Yeni backing array yaradılır → ortaq array yazılmır:

```go
s := make([]int, 1)          // len=1, cap=1 → DOLU
go func() { s1 := append(s, 1); fmt.Println(s1) }()  // race YOXDUR
go func() { s2 := append(s, 1); fmt.Println(s2) }()
```

**Race YARADAN hal — DOLU OLMAYAN slice:** Hər ikisi eyni indeksə (1) yazır:

```go
s := make([]int, 0, 1)      // len=0, cap=1 → BOŞ YER VAR
go func() { s1 := append(s, 1) }()   // WARNING: DATA RACE!
go func() { s2 := append(s, 1) }()
```

**Həll — kopya üzərində işlə:**

```go
sCopy := make([]int, len(s), cap(s))
copy(sCopy, s)
s1 := append(sCopy, 1)
```

**Prinsip:** Doluluqdan asılı iki müxtəlif davranış YOX — ortaq slice üzərində append = data race riski; konkurunt kodda qaçınılmalı.

---

### #70: Using mutexes inaccurately with slices and maps

**Data race qaydaları (slice vs map):**

| Əməliyyat | Data race? |
|-----------|------------|
| Eyni slice indeksi (yazma ilə) | BƏLİ — eyni yaddaş |
| Fərqli slice indeksləri | Xeyr — fərqli yaddaş |
| Eyni map (fərqli key-lər belə!) | BƏLİ — bucket hash-i random initializasiyaya bağlı; race detector həmişə xəbərdarlıq edir |

**Yanlış kopya təsəvvürü:**

```go
// PİS — map assign = DATANI KOYAMIR (runtime.hmap strukturunun kopyası):
func (c *Cache) AverageBalance() float64 {
    c.mu.RLock()
    balances := c.balances     // eyni bucket-lərə pointer!
    c.mu.RUnlock()
    for _, balance := range balances { /* ... */ }   // RACE — AddBalance paralel yazır
}
```

(Slice-da da eyni: `s2 := s1` — eyni backing array.)

**Həllər:**

```go
// 1. Əməliyyat yüngüldürsə — bütün funksiyanı qoru:
func (c *Cache) AverageBalance() float64 {
    c.mu.RLock()
    defer c.mu.RUnlock()
    sum := 0.
    for _, balance := range c.balances {
        sum += balance
    }
    return sum / float64(len(c.balances))
}

// 2. Əməliyyat AĞIRDIRSA — dərin kopya (yalnız kopya kritik seqmentdə):
func (c *Cache) AverageBalance() float64 {
    c.mu.RLock()
    m := make(map[string]float64, len(c.balances))
    for k, v := range c.balances {
        m[k] = v                 // DƏRİN kopya
    }
    c.mu.RUnlock()               // sonra lock'suz emal
    sum := 0.
    for _, balance := range m { sum += balance }
    return sum / float64(len(m))
}
```

Seçim: element sayı, struct ölçüsü, əməliyyatın ağırlığı (məs., xarici DB çağırışı → kopya variantı).

---

### #71: Misusing sync.WaitGroup

**API:** `Add(int)` artır, `Done()` azaldır (= Add(-1)), `Wait()` 0-a qədər bloklar. Counter mənfi ola bilməz (panic).

**Səhv — Add goroutine DAXİLINDƏ:**

```go
wg := sync.WaitGroup{}
var v uint64
for i := 0; i < 3; i++ {
    go func() {
        wg.Add(1)               // PİS — valideyn Wait-dən əvvəl gəlməyə bilər!
        atomic.AddUint64(&v, 1)
        wg.Done()
    }()
}
wg.Wait()                        // v = 0..3 arası QƏRARSIZ + data race!
```

Ssenari: sonuncu goroutine Add çağırmamış Wait keçə bilər → v yarımçıq oxunur. Atomic-in özü WaitGroup sırasını qoruya bilməz — **happens-before Add↔Wait arasıdır.**

**Həllər:**

```go
// 1. Sayı əvvəlcədən məlum → Add(3) loop-dan əvvəl
wg.Add(3)
for i := 0; i < 3; i++ {
    go func() { /* ... */ }()
}

// 2. Add hər iterasiyada, goroutine-dən ƏVVƏL:
for i := 0; i < 3; i++ {
    wg.Add(1)
    go func() { /* ... */ }()
}
```

**Qızıl qayda:** `Add` valideyn goroutine-də, spin-dən əvvəl; `Done` uşaq goroutine-də.

---

### #72: Forgetting about sync.Cond (Şərt dəyişəni)

**Ssenari:** Donation goals — 1 updater balance artırır; listener-lər ($10, $15 hədəfləri) xəbərdarlıq gözləyir.

**Yanaşma 1 — mutex busy loop:** İşləyir amma listener-lər fasiləsiz `for balance < goal` yoxlayır → **CPU israfı dəhşətli.**

**Yanaşma 2 — channel:** Mesaj yalnız BİR goroutine-ə çatır (round-robin distribution!) → $10 hədəfi $11-də xəbər tutdu:

```
$11 goal reached    ← 1-ci goroutine cəmi tək nömrələr aldı (1,3,5...)
$15 goal reached
```

Channel broadcast YALNIZ bağlanma ilə mümkündür (bir dəfəlik!) — təkrar broadcast üçün uyğun deyil.

**Yanaşma 3 — sync.Cond (doğru həll):**

```go
type Donation struct {
    cond    *sync.Cond
    balance int
}
donation := &Donation{cond: sync.NewCond(&sync.Mutex{})}

// Listener-lər:
f := func(goal int) {
    donation.cond.L.Lock()
    for donation.balance < goal {
        donation.cond.Wait()    // şərt dəyişəninə qədər yuxu
    }
    fmt.Printf("%d$ goal reached\n", donation.balance)
    donation.cond.L.Unlock()
}
go f(10)
go f(15)

// Updater:
for {
    time.Sleep(time.Second)
    donation.cond.L.Lock()
    donation.balance++
    donation.cond.L.Unlock()
    donation.cond.Broadcast()   // BÜTÜN gözləyənləri oyat
}
```

**Wait-in daxili addımları:** (1) mutex-i aç → (2) goroutine-i saxla → (3) bildiriş gələndə mutex-i bağla. Buna görə Wait kritik seqment DAXİLİNDƏ çağrılır.

**Broadcast vs Signal:** Broadcast = bütün gözləyənlər; Signal = bir (non-blocking `chan struct{}` send-inə bərabər). **Çatışmazlıq:** Boş gözləyici yoxdansa Broadcast İTİR (channel mesajı kimi buffer-lənmir).

**Nə vaxt sync.Cond:** Təkrar-təkrar çoxlu goroutine-ə broadcast → sync.Cond.

---

### #73: Not using errgroup

**Problem:** N paralel çağırış + error aqreqasiyası — WaitGroup ilə error yığmaq (error slice / mutex / error channel) mürəkkəbdir.

**errgroup (golang.org/x/sync/errgroup):**

```go
func handler(ctx context.Context, circles []Circle) ([]Result, error) {
    results := make([]Result, len(circles))
    g, ctx := errgroup.WithContext(ctx)   // ortaq context!
    for i, circle := range circles {
        i := i
        circle := circle
        g.Go(func() error {
            result, err := foo(ctx, circle)
            if err != nil {
                return err               // error closure-dan qayıdır
            }
            results[i] = result           // indeks-əsaslı yazma = race-free
            return nil
        })
    }
    if err := g.Wait(); err != nil {     // İLK non-nil error
        return nil, err
    }
    return results, nil
}
```

**Faydaları:**
1. `g.Go(f func() error)` — goroutine + error qaytarma vahidə.
2. `g.Wait()` — hamısı bitənə qədər bloklar, ilk error-u çatdırır.
3. **Ortaq context ləğvi:** Bir çağırış 1 ms-də error versə → ctx ləğv olunur → digər 5 saniyəlik çağırışlar dayandırılır → növbəti 4 saniyə gözləmə YOXDUR.

**Şərt:** `g.Go` içi context-aware olmalıdır, yoxsa ləğv təsirsizdir.

---

### #74: Copying a sync type (Sync tipinin kopyalanması)

**Səhv — value receiver mutex-i kopyalayır:**

```go
type Counter struct {
    mu       sync.Mutex
    counters map[string]int
}
func NewCounter() Counter { return Counter{counters: map[string]int{}} }
func (c Counter) Increment(name string) {   // VALUE receiver → struct+mutex kopyası!
    c.mu.Lock()
    defer c.mu.Unlock()
    c.counters[name]++
}
// 2 goroutine paralel → DATA RACE — hər biri ÖZ kopya mutexi ilə işləyir!
```

**Qadağan kopyalanan tiplər:** `sync.Cond`, `sync.Map`, `sync.Mutex`, `sync.RWMutex`, `sync.Once`, `sync.Pool`, `sync.WaitGroup`.

**Həllər:**

```go
// 1. Pointer receiver:
func (c *Counter) Increment(name string) { /* eyni kod */ }

// 2. Pointer sahə (value receiver saxlanılsa):
type Counter struct {
    mu       *sync.Mutex   // pointer kopyası → eyni mutex!
    counters map[string]int
}
func NewCounter() Counter {
    return Counter{
        mu:       &sync.Mutex{},   // nil olarsa Lock() panic!
        counters: map[string]int{},
    }
}
```

**Kopyalanma riski yaradan hallar:** value receiver metod; sync arqumentli funksiya; sync sahəsi olan struct arqumenti. **`go vet` tutur:** `Increment passes lock by value: Counter contains sync.Mutex`.

---

## Əsas terminlər

| Termin | İzah |
|--------|------|
| Detach context | Cancel söndürülmüş, dəyərləri saxlayan custom context |
| HTTP request context ləğvi | Response yazılan andan etibarən! |
| Goroutine leak | Exit nöqtəsi olmayan goroutine — 2KB+ yaddaş və resurs tutur |
| Closure loop variable | Bədəndən kənar i-yə istinad — icra zamanı cari dəyər |
| Pseudo-random select | Multiple ready case-lərdən random seçim (starvation qoruması) |
| Inner select + default | Prioritet pattern — disconnect-də qalan mesajları tüketmə |
| Notification channel | `chan struct{}` — datasız siqnal; struct{} = 0 bayt |
| Nil channel | Göndər/qəbul = əbədi blok; select case söndürmə aləti |
| `v, open := <-ch` | Qapalı channel: open=false, zero value — ayrıseçmə üsulu |
| Busy-waiting | Fasiləsiz loop yoxlama — CPU israfı |
| Synchronous channel | Unbuffered — güclü sinxronizasiya qarantisi |
| LMAX qeydi | Növbələr həmişə az-çox dolu/boş — balanslı ölçü nadirdir |
| `%v` yan təsiri | Stringer çağırışı → etcd data race / RWMutex deadlock |
| Mutex scope daraldılması | Validasiya lock-dan əvvəl — deadlock qarşısı |
| Append race | Dolu slice → yeni array (safe); dolmamış → eyni indeks yazma (race) |
| Shallow copy tələsi | Map/slice assign = eyni data bucket-ləri |
| Round-robin distribution | Çoxlu receiver-li channel-də mesajlar ardıcıl paylanır |
| `sync.Cond` | Şərt dəyişəni — təkrar broadcast üçün; Wait = unlock-suspend-relock |
| Broadcast vs Signal | Hamısını oyat / birini oyat (broadcast boşluqda İTİR) |
| `errgroup.WithContext` | Goroutine qrupu + error + context ləğvi vahiddə |
| Sync kopya qadağası | Mutex/WaitGroup və s. heç vaxt kopyalanmır — pointer istifadə et |

---

## Praktik nəticə

1. **Asinxron iş + HTTP context = təhlükə:** Response yazılınca context ölür — `context.Background()` və ya detach custom context istifadə et.
2. **Hər `go` ifadəsindən əvvəl soruş: "bu nə vaxt dayanacaq?"** — Exit planı + resurs bağlanması üçün blocking close/Wait.
3. **Loop dəyişəni + closure goroutine:** `val := i` və ya parametr `(i)` — yalnız 2 düzgün yol.
4. **Select random seçir:** Case sırası prioritet DEMƏZLİDİR — unbuffered/tək channel (1 producer) və ya inner select+default (çoxlu producer).
5. **Siqnal kanalı = `chan struct{}`:** bool YOX; 0 bayt + aydın semantika.
6. **Nil channel pattern:** Qapalı kanalı nil-lə → select case-i deaktiv et → merge kimi state machine-lər təmiz.
7. **Buffered ölçü = 1 default:** Worker-pool/limit halları istisna; magic rəqəmlər üçün səbəb kommenti şərtdir; sinxronizasiya = unbuffered.
8. **Formatlaşdırma gizli çağırışdır:** `%v` Stringer çağırır (deadlock), `%v` ctx dəyərləri oxuyur (etcd race) — lock saxlarkən birbaşa sahələr çap et.
9. **Ortaq slice-da append = race:** Doluluqdan asılı davranışa etibar etmə — kopya istifadə et.
10. **Map/slice assign = shallow copy:** Kritik seqmentə dərin kopya YOXDUR — ya bütün iterasiya lock içində, ya da dərin kopya sonra lock-suz.
11. **WaitGroup: Add valideyndə, spin-dən əvvəl; Done uşaqda** — əks halda qeyri-determinizm + race.
12. **Təkrar çoxlu broadcast → sync.Cond:** Channel round-robin paylayır; Broadcast bütün gözləyənləri oyadır (amma gözləyən yoxdansa itir).
13. **Paralel + error + context → errgroup:** `g.Go` + `g.Wait` + avtomatik context ləğvi.
14. **Sync tiplər kopyalanmaz:** Pointer receiver və ya `*sync.Mutex` sahəsi; `go vet` yakalayar.

---

## Mənbə

- Kitab: *100 Go Mistakes and How to Avoid Them* — Teiva Harsanyi, Manning, 2022 (ISBN 9781617299599)
- Chapter 9: "Concurrency: Practice", book səh. 193–233
- PDF səhifələri: 213–253
- İstinadlar: Go spec — select (go.dev/ref/spec); etcd PR 7816 (github.com/etcd-io/etcd/pull/7816); LMAX Disruptor white paper; sync (pkg.go.dev/sync); errgroup (golang.org/x/sync/errgroup)
