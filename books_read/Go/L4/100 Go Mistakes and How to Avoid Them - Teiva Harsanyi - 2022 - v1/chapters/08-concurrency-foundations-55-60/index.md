# Chapter 8 — Concurrency: Foundations (#55-#60)

## Bu chapter nədən bəhs edir?

Bu chapter konkurentliyin fundamental anlayışlarını əhatə edir: concurrency vs parallelism (kofe dükanı metaforası), Go scheduler-in daxili mexanizmi (G/M/P, work stealing, preemptive), konkurrentliyin həmişə sürətli OLMAMASI (parallel merge sort), channels vs mutexes seçimi, data races vs race conditions, Go memory model-in zəmanətləri, workload tiplərinin (CPU/IO-bound) worker pool-a təsiri və `context.Context`-in tam analizi.

---

## Əsas fikirlər

### #55: Mixing up concurrency and parallelism (Konkurrentlik və paralellik qarışdırmaq)

**Kofe dükanı metaforası:**

| Dizayn | Nə edir | Nədir |
|--------|---------|-------|
| 1 ofisiant + 1 kofe maşını | Sadə (sequential) model | — |
| 2 ofisiant + 2 maşın (tam dublikat) | Eyni işi eyni anda çoxdəfə | **PARALELLİK** |
| Rolların bölünməsi: sifariş qəbul / dənə üyüdücü / maşın + gözləmə növbəsi | Strukturun dəyişməsi | **KONKURRENTLİK** |
| 2-ci üyüdücü ofisiant | Bir addımın paralelləşməsi | Struktur EYNİ, paralellik ARTIR |
| 3 kofe maşını | Contention azalır | Paralellik daha da artır |

**Rob Pike:** *"Concurrency is about dealing with lots of things at once. Parallelism is about doing lots of things at once."*

- **Concurrency = STRUKTUR** — addımlara bölüb konkurrent thread-lərə paylamaq
- **Parallelism = İCRA** — bir addımda çoxlu thread
- **Concurrency ENABLES parallelism** — konkurrent struktur paralelləşdirilə bilən hissələr yaradır

---

### #56: Thinking concurrency is always faster (Konkurrentlik həmişə sürətlidir?)

**Go scheduler-in əsasları:**

**Thread vs Goroutine:**

| Xüsusiyyət | OS thread | Goroutine |
|-----------|-----------|-----------|
| Context switch edən | OS | Go runtime |
| Yaddaş ölçüsü | ~2 MB (Linux x86-32) | **2 KB** (Go 1.4+) |
| Switch sürəti | baz | 80-90% DAHA SÜRƏTLİ |

**Context switch:** Thread-in icra state-i (registerlər) yaddaşa yazılır → baha əməliyyat.

**G/M/P terminologiyası:**
- **G** = Goroutine · **M** = OS thread (machine) · **P** = CPU core (processor)
- Hər M bir P-yə təyin olunur; hər G bir M üzərində işləyir
- **GOMAXPROCS** = paralel user-kod icra edən M limiti (Go 1.5+: default = CPU core sayı); sistem çağırısında bloklanan M-lərin əvəzinə runtime DAHA ÇOX M aça bilər

**Goroutine həyat dövrü:** Executing (M-də icra) → Runnable (növbədə) → Waiting (syscall/mutex/channel gözləmə).

**Növbələr:** Hər P üçün **local queue** + bütün P-lər üçün **global queue**.

**Schedule psevdo-kodu:**

```
runtime.schedule() {
    // Hər 61-ci dəfə: global queue-də G varmı? yoxdursa → local queue
    // Yoxdursa → başqa P-lərdən OĞURLA (work stealing)
    // Yoxdursa → global queue → poll network
}
```

**Work stealing** — istifadəsiz P digər P-lərin goroutine-lərini oğurlaya bilər.

**Preemption (Go 1.14+):** Əvvəl kooperativ idi (yalnız bloklayan hallarda switch); indi 10 ms işləyən goroutine preemtive marker alır — uzun iş CPU paylaşmağa məcbur olunur.

**Parallel merge sort təcrübəsi:**

```go
// V1 — hər yarım ayrı goroutine-də:
func parallelMergesortV1(s []int) {
    if len(s) <= 1 { return }
    middle := len(s) / 2
    var wg sync.WaitGroup
    wg.Add(2)
    go func() { defer wg.Done(); parallelMergesortV1(s[:middle]) }()
    go func() { defer wg.Done(); parallelMergesortV1(s[middle:]) }()
    wg.Wait()
    merge(s, middle)
}
```

**Benchmark (4 core, 10,000 element) — SÜRPRİZ:**

| Variant | ns/op |
|---------|-------|
| Sequential | 2,278,993,555 |
| ParallelV1 | **17,525,998,709 — ~8x YAVAŞ!** |
| ParallelV2 (threshold) | 1,313,010,260 — 40% sürətli |

**V1 niyə yavaş?** 1024 element → hər yarım üçün goroutine → 512 → 256 → ... → 1 elementə qədər miniklik işlər üçün goroutine yaratmaq + scheduling xərci birbaşa birləşdirmədən BÖYÜKDÜR.

**V2 həlli — threshold:**

```go
const max = 2048    // magic value — BENCHMARK ilə tapılır!

func parallelMergesortV2(s []int) {
    if len(s) <= 1 { return }
    if len(s) <= max {
        sequentialMergesort(s)          // kiçik iş — sequensial
    } else {
        // ... paralel bölgü (V1 kimi)
    }
}
```

**Dərslər:**
1. Konkurrentlik default YOX — sadə sequensial versiyadan başla, profiling (#98) + benchmark (#89) ilə doğrula.
2. Magic threshold dəyərlər production-oxşar mühitdə benchmark ilə təyin olunur. (Java threads üçün optimal ~8192 — goroutine-lərin thread-lərdən nə qədər effektiv olduğunun sübutu.)
3. Modern CPU-lar sequensial + predictable koda çox yaxşıdır (superscalar).

---

### #57: Being puzzled about when to use channels or mutexes (Channel yoxsa mutex?)

**Channel xatırlatma:** kommunikasiya mexanizmi — goroutine-ləri birləşdirən pipe. Unbuffered (göndərən qəbul edənə qədər bloklar) və ya buffered (buffer dolanda qədər).

**Goroutine münasibətləri:**

```
G1 ←paralel→ G2      (eyni addım — məs., 2 eyni HTTP handler)
G1/G2 → konkurrent → G3   (növbəti addım — nəticə toplayıcı)
```

**Qərar qaydası:**

| Goroutine-lər | Ehtiyac | Alət |
|---------------|---------|------|
| **Paralel** (G1↔G2) | **Synchronization** — ortaq resursa eksklüziv çıxış | **MUTEX** (buffered channel da işləməz!) |
| **Konkurrent** (G2→G3) | **Koordinasiya/orkestrasiya** — nəticə siqnalı, ownership transfer | **CHANNEL** |

**Semantika:** Mutex = state/resurs paylaşımı; channel = siqnalizasiya (data ilə `chan T` və ya sizinsiz `chan struct{}`). Birlikdə tamamlayıcıdırlar — birini həmişə məcbur etmək səhvdir.

---

### #58: Not understanding race problems (Race problemləri)

#### Data race vs Race condition

**Data race:** 2+ goroutine eyni yaddaş yerinə eyni anda çıxış edir və ən azı biri YAZIR:

```go
i := 0
go func() { i++ }()
go func() { i++ }()
// -race: WARNING: DATA RACE; nəticə bəzən 1, bəzən 2!
```

**Mexanizm:** `i++` = 3 əməliyyat (oxu → artır → yaz). İnterleaved icrada hər ikisi 0 oxuyur, hər ikisi 1 yazır.

**3 müdafiə yolu:**

```go
// 1. ATOMİK əməliyyat (sync/atomic — yalnız int32/int64/uint32/uint64!):
var i int64
go func() { atomic.AddInt64(&i, 1) }()
go func() { atomic.AddInt64(&i, 1) }()

// 2. MUTEX — kritik seqment:
i := 0
mutex := sync.Mutex{}
go func() { mutex.Lock(); i++; mutex.Unlock() }()
go func() { mutex.Lock(); i++; mutex.Unlock() }()

// 3. KOMMUNİKASİYA — yalnız 1 goroutine yazır:
i := 0
ch := make(chan int)
go func() { ch <- 1 }()
go func() { ch <- 1 }()
i += <-ch
i += <-ch
```

**Race condition (yarış şəraiti):** Data race YOXDUR, amma davranış nəzarət olunmayan hadisə sırasından asılıdır:

```go
i := 0
mutex := sync.Mutex{}
go func() { mutex.Lock(); defer mutex.Unlock(); i = 1 }()   // 1 və ya
go func() { mutex.Lock(); defer mutex.Unlock(); i = 2 }()   // 2 — QƏRARSIZ!
```

Mutex data race-i aradan qaldırır, amma icra SIRASI qeyri-müəyyən qalır. Determinizm üçün koordinasiya lazımdır — məs., channel-lərlə ardıcıllıq qarantisi (hətta mutex-i tam aradan qaldırmaq olar).

#### Go memory model (golang.org/ref/mem)

Spesifikasiya: bir goroutine-dəki yaz digər goroutine-dəki oxu üçün nə vaxt "happens-before" qarantisi verir. Notasiya: A < B (A, B-dən əvvəl baş verir).

**Zəmanətlər:**

1. **Goroutine yaratmaq → onun icrasından əvvəl:** `i := 0; go func() { i++ }()` — RACE YOXDUR (yaratma icradan öndədir).

2. **Goroutine çıxışı HƏR HANSI hadisədən əvvəl olma zəmanəti YOXDUR:**

```go
i := 0
go func() { i++ }()
fmt.Println(i)   // DATA RACE!
```

3. **Channel send → müvafiq receive tamamlanmadan əvvəl:**

```go
i := 0
ch := make(chan struct{})
go func() { <-ch; fmt.Println(i) }()
i++
ch <- struct{}{}
// sıra: increment < send < receive < read → SINKRON, race yoxdur
```

4. **Channel close → həmin bağlanmanın receive-indən əvvəl:** `close(ch)` variantı da race-free.

5. **Unbuffered channel receive → send tamamlanandan əvvəl (ƏKS-İNTUİTİV!):**

```go
// BUFFERED (cap 1) → DATA RACE:
i := 0
ch := make(chan struct{}, 1)
go func() { i = 1; <-ch }()
ch <- struct{}{}     // send dərhal tamamlanır (bufferə düşür)
fmt.Println(i)       // yaz ilə ox eyni anda ola bilər!

// UNBUFFERED → RACE-FREE:
ch := make(chan struct{})
go func() { i = 1; <-ch }()
ch <- struct{}{}     // send QƏBUL bitənə qədər tamamlanMIR
fmt.Println(i)       // yaz < qəbul < send-tamamlanma < oxu → qarantiyalı
```

Ox (yarıqlar) səbəbkarlığı deyil — sıralama ZƏMANƏTİNİ göstərir: unbuffered-da yaz həmişə oxudan əvvəl.

---

### #59: Not understanding the concurrency impacts of a workload type (Workload tipinin təsiri)

**Workload təsnifatı:**

| Tip | Məhdudlaşdıran | Nümunə |
|-----|----------------|--------|
| **CPU-bound** | CPU sürəti | Merge sort alqoritmi |
| **I/O-bound** | I/O sürəti | REST çağırışı, DB sorğusu |
| **Memory-bound** | Yaddaş | (nadir — memory ucuzlaşıb) |

**Worker pooling pattern:**

```go
func read(r io.Reader) (int, error) {
    var count int64
    wg := sync.WaitGroup{}
    var n = 10                       // pool size — workload tipinə görə!
    ch := make(chan []byte, n)       // capacity = pool size (contention azaldır)
    wg.Add(n)
    for i := 0; i < n; i++ {
        go func() {
            defer wg.Done()
            for b := range ch {
                v := task(b)
                atomic.AddInt64(&count, int64(v))
            }
        }()
    }
    for {
        b := make([]byte, 1024)
        // Read from r to b
        ch <- b                      // tapşırıqları yayımla
    }
    close(ch)
    wg.Wait()
    return int(count), nil
}
```

Fiksiləşmiş pool: resurs təsirini məhdudlaşdırır, xarici sistemin flood olunmasının qarşısını alır.

**Pool size qərarı:**

- **I/O-bound:** xarici sistemin dözə biləcəyi konkurrent çıxış sayı (məs., DB connection limiti)
- **CPU-bound:** **GOMAXPROCS-a yaxın** (`n := runtime.GOMAXPROCS(0)`) — 4 core/4 thread-də 4 goroutine; 5-ci goroutine bir thread-i paylaşar → artıq context switch

**Niyə NumCPU() YOX?** GOMAXPROCS dəyişilə bilər və core sayından AZ ola bilər — 4 core amma 3 thread → 3 goroutine düzgün (4-ü context switch artırar).

**İdeal yayılma (4 core):** Work stealing + OS M/P hərəkəti nəticəsində hər M ayrı core-da — amma bu, developer-in tələbi ilə YOX, uygün şərtlərlə (GOMAXPROCS-based pool) MÜMKÜN olur.

**Vacib:** Fikirləri benchmark ilə doğrula — konkurrentlikdə tez-tez yanlış pressumplar.

---

### #60: Misunderstanding Go contexts (Go context-ləri)

**Rəsmi tərif:** *"A Context carries a deadline, a cancellation signal, and other values across API boundaries."*

#### Deadline

```go
// Radar hər 4 saniyə mövqe göndərir → köhnə mövqeyə 4 saniyədən çox ehtiyac YOX:
func (h publishHandler) publishPosition(position flight.Position) error {
    ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
    defer cancel()    // ŞƏRT — yoxsa goroutine 4 saniyə yaddaşda qalır
    return h.pub.Publish(ctx, position)
}
```

- `context.WithTimeout(parent, duration)` → (ctx, cancel)
- `defer cancel()`: WithTimeout daxili goroutine yaradır — cancel çağırılmazsa 4 saniyəlik yaddaş tutulumu.

#### Cancellation signals

```go
func main() {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()    // main qayıdanda file watcher fd bağlansın
    go func() {
        CreateFileWatcher(ctx, "foo.txt")   // context-aware: ctx cancel → descriptor bağlanır
    }()
    // ...
}
```

#### Context values

```go
ctx := context.WithValue(parentCtx, "key", "value")
fmt.Println(ctx.Value("key"))   // value
```

**Key tələsi:** String key → paketlər arası toqquşma (bir dəyər digərini əzər). **Həll — unexported custom tip:**

```go
package provider
type key string
const myCustomKey key = "key"    // unexported tip → toqquşma mümkünsüz

func f(ctx context.Context) {
    ctx = context.WithValue(ctx, myCustomKey, "foo")
}
```

**Use case-lər:** correlation ID (tracing — funksiya imzasını çirkləndirmədən), HTTP middleware kommunikasiyası:

```go
type key string
const isValidHostKey key = "isValidHost"

func checkValid(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        validHost := r.Host == "acme"
        ctx := context.WithValue(r.Context(), isValidHostKey, validHost)
        next.ServeHTTP(w, r.WithContext(ctx))   // növbəti middleware/handler-a ötür
    })
}
```

#### Catching a context cancellation

**`Done()` → `<-chan struct{}`** — context ləğv/deadline olduqda BAĞLANIR. (Bağlanma seçilir, çünki channel-in bağlanması BÜTÜN istehlakçı goroutine-lərin aldığı yeganə channel əməliyyatıdır.)

**`Err()`** — Done bağlı deyilsə nil; bağlıdırsa səbəb: `context.Canceled` və ya `context.DeadlineExceeded`.

```go
func handler(ctx context.Context, ch chan Message) error {
    for {
        select {
        case msg := <-ch:
            // Do something with msg
        case <-ctx.Done():
            return ctx.Err()
        }
    }
}
```

**Context qəbul edən funksiyada BLOCKING channel əməliyyatı YOX:**

```go
// PİS — ctx ləğv olsa belə bloklar:
ch1 <- struct{}{}
v := <-ch2

// DÜZGÜN — hər əməliyyat select ilə:
select {
case <-ctx.Done():
    return ctx.Err()
case ch1 <- struct{}{}:
}
select {
case <-ctx.Done():
    return ctx.Err()
case v := <-ch2:
    // ...
}
```

**Digər qaydalar:**
- Bütün standart context-lər çox-goroutine təhlükəsizdir.
- Gözlənilən funksiyalar context QƏBUL ETMƏLİDİR — upstream caller ləğvetmə qərarı verə bilsin.
- Context hansının bilinmədiyi/əlçatmaz olduğu hallarda `context.TODO()` (`context.Background` əvəzinə) — semantik siqnal: "hələ məlum deyil".

---

## Əsas terminlər

| Termin | İzah |
|--------|------|
| Concurrency | STRUKTUR — addımlara bölmə; "dealing with lots of things at once" |
| Parallelism | İCRA — eyni anda çoxlu; "doing lots of things at once" |
| G / M / P | Goroutine / OS thread (machine) / CPU core (processor) |
| GOMAXPROCS | Paralel user-kod M limiti — default: CPU core sayı |
| Context switch | Thread-in state yaddaşına yazılıb dəyişdirilməsi — baha |
| Goroutine stack | 2 KB (thread ~2 MB) → sürətli switch |
| Work stealing | Boş P-nin digər P-lərdən goroutine oğurlaması |
| Preemptive scheduling | Go 1.14+: 10 ms-dən uzun goroutine məcburi paylaşır |
| Threshold (paralellik) | Magic dəyər — bu dəyərdən kiçik iş sequensial; benchmark ilə tapılır |
| Unbuffered / buffered channel | Göndərən qəbula qədər bloklar / buffer dolana qədər |
| Synchronization vs coordination | Paralel goroutine-lər üçün mutex / konkurrentlər üçün channel |
| Data race | Eyni yaddaş + eyni anda + ən azı bir yaz → qeyri-müəyyən nəticə |
| Race condition | Data race olmasa belə sıra/timing asılılığı |
| `atomic.AddInt64` | Bölünməz artırma — yalnız xüsusi int tipləri |
| Critical section | Mutex ilə qorunmuş, ən çoxu 1 goroutine |
| Go memory model | Yaz/oxu sıralama zəmanətləri (happens-before) |
| CPU-bound / I/O-bound | CPU sürəti / I/O sürəti ilə məhdud iş |
| Worker pooling | Fiksiləşmiş goroutine hovuzu + ortaq channel |
| `context.WithTimeout/WithCancel/WithValue` | Deadline / ləğvetmə / dəyər daşıyan context yaratma |
| `defer cancel()` | WithTimeout-un daxili goroutine-inin dərhal təmizlənməsi |
| `ctx.Done()` / `ctx.Err()` | Ləğv channel-ı / səbəb error-u |
| Unexported context key | Toqquşmasız key üçün custom tip patterni |
| `context.TODO()` | "Hansı context hələ bilinmir" semantik siqnalı |

---

## Praktik nəticə

1. **Struktur (concurrency) ilə icranı (parallelism) ayır:** Rol bölməsi = konkurrentlik; eyni rolun çoxaldılması = paralellik. Konkurrentlik paralelliyi MÜMKÜN edir.
2. **Konkurrentlik həmişə sürətli deyil:** Miniklik işlərdə goroutine yaratmaq xərci qəbuldan böyükdür — sequensial başla, benchmark ilə paralelliyi doğrula, kiçik hissələr üçün threshold əlavə et.
3. **G/M/P modelini başa düş:** GOMAXPROCS M limitidir; work stealing load balans edir; Go 1.14+ preemptivdir.
4. **Paralel goroutine → mutex, konkurrent → channel:** State paylaşımı = mutex; siqnal/koordinasiya/ownership = channel. İkisi tamamlayıcıdır.
5. **Data race ≠ race condition:** `-race` detector-u data race tutur; determinizm üçün isə əlavə koordinasiya (channel ardıcıllığı) lazımdır.
6. **Memory model zəmanətlərini yadda saxla:** go yaratma < icra; goroutine çıxışı zəmanətsiz; unbuffered receive < send-tamamlanma (buffered-da bu yoxdur!).
7. **Worker pool ölçüsü workload-a bağlıdır:** CPU-bound → GOMAXPROCS; I/O-bound → xarici sistemin tutumu. NumCPU() YOX (GOMAXPROCS dəyişilə bilər).
8. **Context 3 şey daşıyır:** deadline + cancel siqnalı + key-value. `defer cancel()` mütləq; unexported key tipi; blocking channel əməliyyatlarını select + Done ilə əvəz et; şübhə halında `context.TODO()`.

---

## Mənbə

- Kitab: *100 Go Mistakes and How to Avoid Them* — Teiva Harsanyi, Manning, 2022 (ISBN 9781617299599)
- Chapter 8: "Concurrency: Foundations", book səh. 162–192
- PDF səhifələri: 182–212
- İstinadlar: Go memory model (golang.org/ref/mem); Go scheduler (mng.bz/N611, mng.bz/lxY8); context (pkg.go.dev/context); Rob Pike "Concurrency is not parallelism"
