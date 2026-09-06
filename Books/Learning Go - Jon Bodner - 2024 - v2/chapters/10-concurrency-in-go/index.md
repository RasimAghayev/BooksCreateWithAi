# Chapter 10 — Concurrency in Go (Go-da Paralellik)

## Bu chapter nədən bəhs edir?

CSP modeli, goroutine-lər, channel-lər (buffered/unbuffered, close, comma-ok, davranış
cədvəli), `select` (random seçim, deadlock qarşısı, default), done channel pattern,
backpressure, timeout idiom-u, WaitGroup, sync.Once, mutex vs channel qərarı və tam
pipeline nümunəsi.

## Əsas fikirlər

### 1. Nə vaxt Paralellik (Concurrency)?
**Fəlsəfə:** Concurrency ≠ parallelism — concurrency problemi struktur etmək alətidir;
paralel icra hardware-dən asılıdır (Amdahl's Law). Yeni developer-in 5 mərhələli
xəyal qırığı: "hər şeyi goroutine-ə qoy → buffer əlavə et → böyük buffer → mutex →
imtina".

**Qaydalar:**
- Addımlar bir-birindən asılı deyilsə və I/O-idisə → concurrency faydalıdır (I/O
  yaddaşdan minlərlə dəfə yavaşdır).
- In-memory sürətli alqoritmlərdə concurrency overhead-i qazancı yeyə bilər.
- Şübhə varsan: serial yaz → benchmark → müqayisə.

### 2. Goroutine-lər
**Nədir:** Go runtime-ın idarə etdiyi yüngül proses; OS thread-ləri üzərində M:N scheduler.

**Üstünlüklər:** yaradılması sürətli (OS resursu yox), kiçik başlanğıc stack (böyüyə
bilir), keçid sürətli (proses daxili), network poller + GC ilə inteqrasiya. Minlərlə
goroutine problem deyil.

**Sintaksis:** `go f()` — qaytarma dəyərləri İQNOR olunur.

**Konvensiya:** goroutine biznes məntiqindən ayrı closure ilə başladılır — closure
channel oxuyub məntiqə ötürür, nəticəni channel-ə yazır. Məntiq concurrency-dən xəbərsiz
→ modul, test edilə bilən, API-də concurrency yox:

**Kitabdan kod nümunəsi:**
```go
func process(val int) int { /* biznes məntiqi */ }

func runThingConcurrently(in <-chan int, out chan<- int) {
    go func() {
        for val := range in {
            result := process(val)
            out <- result
        }
    }()
}
```

### 3. Channel-lər
**Nədir:** Goroutine-lərin ünsiyyət vasitəsi; `make(chan int)`; reference tipi (pointer
kimi ötürülür); zero value — nil.

**Əməliyyatlar:**
```go
ch := make(chan int)      // unbuffered
ch := make(chan int, 10)  // buffered, cap=10
a := <-ch                 // oxu
ch <- b                   // yaz
len(ch); cap(ch)          // buffer doluluğu / ölçüsü
```
- **Yön tipləri:** `<-chan int` (yalnız oxu), `chan<- int` (yalnız yaz) — compiler
  yoxlaması üçün.
- **Unbuffered:** yaz oxuyana qədər, oxu yazana qədər bloklanır — relay estafeti.
- **Buffered:** buffer dolana qədər yaz sərbəst; dolu yaz / boş oxu bloklanır.
- **for-range channel:** tək dəyişən (dəyər); channel bağlanana qədər.
- **close:** yazan goroutine bağlayır; bağlı channel-ə yaz/close → **PANIC**; bağlıdan
  oxu → sıralı qalan dəyərlər, sonra zero value + `v, ok := <-ch` (ok=false = bağlı).
  Channel lazımsız olanda GC onu yığır — close yalnız gözləyən varsa məcburidir.

**Channel davranış cədvəli (əzbərləniləcək):**

| Vəziyyət | Oxu | Yaz | Close |
|---|---|---|---|
| Unbuffered, açıq | yazılana qədər pause | oxunana qədər pause | işləyir |
| Unbuffered, bağlı | zero value (comma-ok) | **PANIC** | **PANIC** |
| Buffered, açıq | boşdursa pause | doludursa pause | işləyir (qalanlar saxlanır) |
| Buffered, bağlı | qalanlar, sonra zero | **PANIC** | **PANIC** |
| Nil | **həmişə asılır** | **həmişə asılır** | **PANIC** |

### 4. `select` — Paralellik üçün Switch
**Nədir:** Çox channel-dən bir əməliyyat seçən kontrol strukturu.

**Qaydalar:**
- Hər case = channel oxu/yazı + gövdə bloku.
- Çoxlu hazır case → **RANDOM seçim** (switch-in "ilk true"-dan fərqli) → starvation
  yoxdur, inconsistent lock sıralaması deadlock-u əngəlləyir.
- **for-select loop:** `for { select { case <-done: return; case v := <-ch: ... } }` —
  exit yolu mütləq.
- **default:** heç bir case hazır deyilsə → nonblocking read/write üçün. Amma for-select
  daxilində default demək olar ki, həmişə YANLIŞDIR (loop CPU-nu yeyir).
- Deadlock: bütün goroutine-lər yatıbsa runtime öldürür: "fatal error: all goroutines
  are asleep - deadlock!". `select` bunu break edir.

**Kitabdan kod nümunəsi (deadlockdan qurtuluş):**
```go
select {
case ch2 <- v:
case v2 = <-ch1:
}
```

### 5. Goroutine + Loop Dəyişəni Tələsi
**Problem:** for-range dəyişənləri hər iterasiyada YENİDƏN istifadə olunur — closure
son dəyəri görür:
```go
for _, v := range a {
    go func() { ch <- v * 2 }()   // hamısı 20 yazacaq!
}
```
**Həll 1:** kölgələmə `v := v`; **Həll 2 (daha aydın):** parametr ötürmə:
```go
for _, v := range a {
    go func(val int) { ch <- val * 2 }(v)
}
```
**Qayda:** Goroutine dəyişə bilən dəyişən istifadə edirsə — cari dəyəri parametr kimi
ötür.

### 6. Goroutine Leak — Hər Zaman Təmizlə
**Nədir:** Çıxa bilməyən goroutine — scheduler ona boş vaxt verir, program yavaşıyır.
Runtime goroutine-in artıq lazımsız olduğunu BİLMİR (variablardan fərqli).

**Generator tələsi:** `countTo` generatorunu for-range yarımçıxdırsa (`break`), goroutine
`ch <- i` yazısında əbədi bloklanır.

**Done Channel Pattern:**
```go
func searchData(s string, searchers []func(string) []string) []string {
    done := make(chan struct{})        // struct{} — dəyər əhəmiyyətsiz, yalnız close
    result := make(chan []string)
    for _, searcher := range searchers {
        go func(searcher func(string) []string) {
            select {
            case result <- searcher(s):
            case <-done:               // bağlı done → sıfır dəyər → çıxış
            }
        }(searcher)
    }
    r := <-result                      // ilk (ən sürətli) nəticə
    close(done)                         // digərlərini yatırt
    return r
}
```

**Cancel funksiyası qaytarma (Ch5 pattern-inin tətbiqi):**
```go
func countTo(max int) (<-chan int, func()) {
    ch := make(chan int)
    done := make(chan struct{})
    cancel := func() { close(done) }
    go func() {
        for i := 0; i < max; i++ {
            select {
            case <-done: return
            case ch <- i:
            }
        }
        close(ch)
    }()
    return ch, cancel
}
// istifadə: ch, cancel := countTo(10); defer/sonra cancel()
```

### 7. Buffered Channel Nə Vaxt?
**Tək cümlə:** "Neçə goroutine başlatdığını bilirsənsə, goroutine sayını məhdudlaşdırmaq
istəyirsənsə, yaxud növbələnən iş miqdarını məhdudlaşdırmaq istəyirsənsə."

**Nümunə 1 — n goroutine nəticəsi:** `results := make(chan int, conc)` — hər goroutine
bloklanmadan yazıb çıxır; leaksiz.

**Nümunə 2 — Backpressure:** paradoks — sistemi məhdudlaşdırmaq onu yaxşılaşdırır.
```go
type PressureGauge struct { ch chan struct{} }
func New(limit int) *PressureGauge {
    ch := make(chan struct{}, limit)
    for i := 0; i < limit; i++ { ch <- struct{}{} }  // "token"lər
    return &PressureGauge{ch: ch}
}
func (pg *PressureGauge) Process(f func()) error {
    select {
    case <-pg.ch:          // token götür
        f()
        pg.ch <- struct{}{}  // qaytar
        return nil
    default:
        return errors.New("no more capacity")  // HTTP 429 kimi
    }
}
```

### 8. select-də case-i Söndürmə — nil Channel Hiyəsi
**Problem:** select-də bağlı channel case-i həmişə "hazır"dır (zero value) — zibil
dəyərlər oxunur. **Həll:** bağlanan channel-i nil-ə mənimsət — nil channel heç vaxt
hazır olmur:
```go
case v, ok := <-in:
    if !ok {
        in = nil    // case bir daha seçilməz
        continue
    }
```

### 9. Timeout İdiomu
```go
func timeLimit() (int, error) {
    var result int
    var err error
    done := make(chan struct{})
    go func() {
        result, err = doSomeWork()
        close(done)
    }()
    select {
    case <-done:
        return result, err
    case <-time.After(2 * time.Second):
        return 0, errors.New("work timed out")
    }
}
```
Qeyd: timeout-dan sonra goroutine DAVAM EDİR (nəticə iqnor) — tam dayandırma üçün
context cancellation (Ch12). Context varsa time.After yerinə context timer işlət
(yuxarı çağırışların timeout-unu da hörmək olur).

### 10. sync.WaitGroup
**Nədir:** Bir goroutine-in bir neçəsini gözləməsi. Zero value istifadəyə hazırdır.

```go
var wg sync.WaitGroup
wg.Add(3)                 // sayğacı qabaqcadan artır
go func() { defer wg.Done(); doThing1() }()
go func() { defer wg.Done(); doThing2() }()
go func() { defer wg.Done(); doThing3() }()
wg.Wait()                  // sayğac 0 olana qədər
```
- `defer wg.Done()` panic halında da zəmanət verir.
- WaitGroup **kopyalana BİLMƏZ** — pointer və ya closure ilə capture (kopyanın Done-u
  orijinalı azaltmaz).
- Ən real istifadə: çoxlu yazan goroutine olduqda channel-i bir dəfə bağlamaq:
```go
go func() {
    wg.Wait()
    close(out)   // bütün işçilər bitəndən sonra
}()
```
- **Birinci seçim DEYİL** — yalnız "hamı bitdikdən sonra təmizlənəcək şey" olduqda.
  Error-toxunan hallar üçün golang.org/x/sync/ErrGroup (bir error hamını dayandırır).

### 11. sync.Once — Dəqiq Bir Dəfə
```go
var parser SlowComplicatedParser
var once sync.Once
func Parse(dataToParse string) string {
    once.Do(func() { parser = initParser() })  // lazy init, yalnız ilk çağırışda
    return parser.Parse(dataToParse)
}
```
Zero value faydalı; kopyalama QADAĞAN; funksiya daxilində elan etmək YANLIŞDIR (hər
çağırış yeni instance). init-in faydalı alternativi: package mutable state yoxdur.

### 12. Tam Pipeline Nümunəsi (A+B → C, 50ms limit)
**Struktur:** context 50ms timeout + buffered channel-lər (hər goroutine bloklanmadan
çıxa bilsin) + select ilə gözləmələr:
```go
func GatherAndProcess(ctx context.Context, data Input) (COut, error) {
    ctx, cancel := context.WithTimeout(ctx, 50*time.Millisecond)
    defer cancel()                      // cancel mütləq çağırılmalı (leak!)
    p := processor{
        outA: make(chan AOut, 1), outB: make(chan BOut, 1),
        inC: make(chan CIn, 1), outC: make(chan COut, 1),
        errs: make(chan error, 2),       // 2 potensial error
    }
    p.launch(ctx, data)                  // 3 goroutine: A, B, C-select
    inputC, err := p.waitForAB(ctx)     // for-select: outA/outB/errs/ctx.Done
    if err != nil { return COut{}, err }
    p.inC <- inputC
    return p.waitForC(ctx)               // select: outC/errs/ctx.Done
}
```
**waitForAB pattern:** count-un 2-yə çatmasını izləyən for-select — hər uğurlu oxu sayğacı
artırır; error və ya ctx.Done hər hansı mərhələdə çıxışı verir. **Sadələşdirmə prinsipi:**
getResultC ctx-i hörməyə etibar etsən — inC/outC channel-ləri və 3-cü goroutine lazımsız:
"düzgün olmaq üçün lazım olan ən az concurrency".

### 13. Mutex-lər — Channel-lərin Qardaşı
**Fəlsəfə:** "Share memory by communicating; do not communicate by sharing memory."
Channel data axınını göstərir; mutex sahibliyi gizlədir. Amma **sahəni qorumaq** üçün
mutex aydındır.

**Katherine Cox-Buday qərar ağacı:**
1. Goroutine-ləri əlaqələndirirsən / dəyər transformasiya zənciri → **channels**
2. Struct sahəsinə çıxışı paylaşmaq → **mutex**
3. Channel benchmark-də kritik performans problemi → **mutex**

**Kitabdan kod nümunəsi:**
```go
type MutexScoreboardManager struct {
    l          sync.RWMutex
    scoreboard map[string]int
}
func (msm *MutexScoreboardManager) Update(name string, val int) {
    msm.l.Lock()
    defer msm.l.Unlock()          // həmişə defer ilə aç!
    msm.scoreboard[name] = val
}
func (msm *MutexScoreboardManager) Read(name string) (int, bool) {
    msm.l.RLock()
    defer msm.l.RUnlock()         // oxu kilidi — paralel oxuyucular OK
    val, ok := msm.scoreboard[name]
    return val, ok
}
```
- `Mutex`: Lock/Unlock — tək sahib. `RWMutex`: yazma tək, oxma çoxuşağalı (RLock/RUnlock).
- **Reentrant DEYİL** — eyni goroutine eyni kilidi 2 dəfə alsa öz-özünə deadlock.
  Rekursiv funksiyada kilidi çağırışdan ƏVVƏL burax; kilid saxlayarkən funksiya çağırmaqdan
  ehtiyatlı ol (o funksiya eyni kilidi istəyə bilər).
- Kopyalama QADAĞAN (WaitGroup/Once kimi) — pointer.
- sync.Map — yalnız "bir dəfə yaz, çox oxu" və ya "goroutine-lər ayrı açarlarda"
  hallarında; əks halda `map + sync.RWMutex` (generics-ə görə sync.Map interface{}-dir).
- sync/atomic — CAS/load/store maşın əməliyyatları; eksperlər üçün — hamı üçün YOX.
- Channel implementasiyası (funksiya göndərən channel + done gözləmə) müqayisədə
  qəlizdir və tək oxucuya imkan verir — scoreboard üçün mutex düzgün seçim.

## Əsas terminlər
- CSP (Communicating Sequential Processes) — Tony Hoare 1978, Go-nun paralellik modeli
- Goroutine — runtime idarəli yüngül icra vahidi
- Channel (kanal) — goroutine-lərarası tipəmən ünsiyyət
- Starvation (aclıq) — həmişə eyni case-in üstün tutulub digərlərinin ac qalması
- Done channel (bitiş kanalı) — close ilə çıxış siqnalı verən struct{} kanalı
- Backpressure (geri təzyiq) — sistemə daxil olan işin məhdudlaşdırılması
- Goroutine leak (goroutine itkisi) — əbədi çıxa bilməyən goroutine
- Critical section (kritik bölmə) — mutex ilə qorunan kod/dəyər
- Reentrant lock (yenidən giriş kilidi) — Go-da OLMAYAN xüsusiyyət

## Praktik nətidə

Paralellik playbook: (1) lazımdırmı? — asılı olmayan I/O addımları + vaxt limiti; (2) API-də
channel/mutex export ETMƏ — closure ilə wrap; (3) dəyişən dəyişəni goroutine-a parametr
ötür; (4) hər goroutine-in çıxış yolunu əvvəlcədən qur: done/cancel/select; (5) buffered
yalnız: sayı bilinən nəticə/limit/növbə idarəsi; (6) timeout: select + time.After (yaxud
ctx); (7) çoxyazanlı channel-də close-u WaitGroup monitoruna həvalə et; (8) struct sahəsi
→ RWMutex; data axını → channel; (9) kilid = defer + kopyasız + funksiya çağırışında
ehtiyat.

## Mənbə
Pages: 289-316
