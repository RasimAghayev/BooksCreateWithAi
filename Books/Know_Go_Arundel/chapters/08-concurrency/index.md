# Chapter 8 — Concurrency (Konkurentlik)

## Bu fəsil nədən bəhs edir?

Set tipinə concurrency safety əlavəsi: data race təhlükəsi, race qarşısının
alınmasının 3 yolu (mutable shared data-nun olmaması > guard goroutine +
channel > mutex), mutex çətinlikləri (contention, deadlock, opt-in xətaları),
SetC tipi — `struct{ mutex *sync.RWMutex; data map[E]struct{} }`,
write lock (Add) vs read lock (All/Contains) fərqi, smoke test (1000
oxu/yazı), mutex olmadan "fatal error: concurrent map iteration and map
write", race detector (`go run -race`, `go test -race`), Intersection-un
optimizasiyası (bir RLock ilə bütün döngü vs granulyar lock trade-off),
Channel məşqi — instrumented channel, sync/atomic paketi (atomic.Uint64 —
mutex paperwork-siz sayaçlar), Set genişləndirmə ideyaları (Len, Delete,
Equals, IsSubset, immutable variant, NewSet/NewUnsafeSet konstruktorları).

## Əsas fikirlər

### 1. Data Race və Onun 3 Qarşısı
```go
// goroutine 1: s.Add(1) — yaz
// goroutine 2: s.All() — oxu  → DATA RACE
```
- **Nədir:** iki goroutine eyni map-ə konkurent yazı/oxu edir — nəticə:
  crash VƏ YA gizli məlumat pozuntusu
- **Qarşının alınması (ən yaxşıdan ən pisə):**
  1. Mutable shared data OLMASIN (ən yaxşı)
  2. Data-nı paylaşınlama: guard goroutine + channel (channel özü
     concurrency-safe)
  3. Mutex (general-purpose set paketi üçün yeganə real yol)
- **sync.Map qeydi:** standart kitabxanada var, amma generics-dən ƏVVƏL
  olduğundan yalnız any saxlayır → az istifadə olunur

### 2. Mutex Problemləri
- **Contention:** goroutine-lər işə yararlı iş əvəzinə lock gözləyirlər
- **Deadlock:** goroutine A mutex-1 tutub mutex-2 gözləyir; B mutex-2
  tutub mutex-1 gözləyir → hər ikisi əbədi bloklanır (yol qovşağındakı
  iki sürücü metaforası). Ən pis hal: proqram "işləyir" amma heç nə
  etmir — aşkarlanması çətin
- **Opt-in xətaları:** lock almağı unutmaq (fəlakət), buraxmağı unutmaq
  (israf) — paket müəllifi accessor metodlarla MƏCBUR edə bilər

### 3. SetC — Konkurrent Təhlükəsiz Set
```go
type SetC[E comparable] struct {
    mutex *sync.RWMutex
    data  map[E]struct{}
}

func NewSetC[E comparable](vals ...E) *SetC[E] {
    s := SetC[E]{
        mutex: &sync.RWMutex{},
        data:  map[E]struct{}{},
    }
    for _, v := range vals {
        s.data[v] = struct{}{}
    }
    return &s
}
```
- **Dizayn:** mutex + data eyni struct-da; mutex POINTER kimi — kopya
  olunarsa iki ayrı mutex yaranar → qoruma ölür ("mutex kopyalanmaz"
  qaydası)
- **RWMutex:** oxu üçün RLock (çoxlu eyni anda mümkün), yazı üçün Lock
  (yalnız tək, oxuyanları da bloklayır)
- **Konstruktor POINTER qaytarır** — dəyişən metodlar pointer receiver-də

### 4. Write Lock və Read Lock
```go
func (s *SetC[E]) Add(vals ...E) {          // YAZI — tam lock
    s.mutex.Lock()
    defer s.mutex.Unlock()
    for _, v := range vals {
        s.data[v] = struct{}{}
    }
}

func (s SetC[E]) All() []E {                // OXU — read lock
    s.mutex.RLock()
    defer s.mutex.RUnlock()
    result := make([]E, 0, len(s.data))
    for v := range s.data {
        result = append(result, v)
    }
    return result
}

func (s SetC[E]) Contains(v E) bool {
    s.mutex.RLock()
    defer s.mutex.RUnlock()
    _, ok := s.data[v]
    return ok
}
```
- **Qayda:** dəyişir → Lock; yalnız oxuyur → RLock
- **defer Unlock** — panic olsa belə lock buraxılır (idiomatik)
- String/Union dəyişməz qalır — yalnız All/Add çağırır, onlar artıq
  lock-aware-dır (kompozisiya faydası)

### 5. Test: Smoke Test + Race Detector
```go
s := NewSetC(1, 2, 3)
var wg sync.WaitGroup
wg.Add(1)
go func() {
    for i := range 1000 {
        s.Add(i)
    }
    wg.Done()
}()
for i := range 1000 {
    _ = s.All()
}
wg.Wait()
fmt.Println("We made it!")
```
- **Smoke test:** tək oxu/yazı kifayət etməz — 1000 dəfə ehtimalı artırır
- **Lock silinərsə:**
  `fatal error: concurrent map iteration and map write` — Go runtime
  xətanı aşkar edib proqramı dərhal dayandırır
- **Race detector — ƏSAS ALƏT:**
```go
go run -race main.go
go test -race
// WARNING: DATA RACE
// Read at 0x... by main goroutine: ...All()
// Previous write at 0x... by goroutine 7: ...Add()
```
- **Dərs:** konkurent kodda -race BÜTÜN testlərdə işlədil. O, mükəmməl
  deyil, amma data race-lərin çoxunu avtomatik tutur. Race bug-ı bu gün
  olmaya bilər — amma olacaq, ən pis: pik yükdə

### 6. Intersection Optimizasiyası — Lock Trade-off
```go
// PIS: hər element üçün s2.Contains() → min 1000 lock/unlock
func (s SetC[E]) Intersection(s2 SetC[E]) *SetC[E] { ... s2.Contains(v) ... }

// YAXŞI: bir RLock — bütün döngü boyu:
func (s SetC[E]) Intersection(s2 SetC[E]) *SetC[E] {
    result := NewSetC[E]()
    s2.mutex.RLock()
    defer s2.mutex.RUnlock()
    for _, v := range s.All() {
        _, ok := s2.data[v]
        if ok {
            result.Add(v)
        }
    }
    return result
}
```
- **Trade-off:** az lock = sürətli, amma uzun müddət saxlanılan lock =
  digər goroutine-lər gözləyir. Granulyar (çox kiçik lock) vs kobud
  (bir böyük lock) — tətbiqə bağlı qərardır.

### 7. Channel Məşqi — atomic İlə Təmiz Həll
**Problemlər (tullanan yazı 1):**
```go
type Channel[T any] struct {
    lock            *sync.RWMutex
    ch              chan T
    sends, receives int
}
// Send/Receive/Sends/Receives — HƏR BİRİ lock+unlock paperwork...
```
- Çirkin: kodun yarısı lock idarəetməsidir

**Daha yaxşı həll — sync/atomic:**
```go
type Channel[T any] struct {
    ch              chan T
    sends, receives atomic.Uint64
}

func New[T any](length int) *Channel[T] {
    return &Channel[T]{ch: make(chan T, length)}
}

func (c *Channel[T]) Send(v T) {
    c.ch <- v
    c.sends.Add(1)
}

func (c *Channel[T]) Receive() T {
    v := <-c.ch
    c.receives.Add(1)
    return v
}

func (c *Channel[T]) Sends() uint64     { return c.sends.Load() }
func (c *Channel[T]) Receives() uint64  { return c.receives.Load() }
```
- **Nədir:** atomic.Uint64 — mutex-siz, konkurrent-artımlı sayaç;
  Add(1) və Load() atomikdir
- **Nəyə lazımdır:** sadə sayaclarda mutex-in YERİNİ tutur — kod xətləri
  azalır, contention itir; kanalın ÖZÜ (chan T) konkurent təhlükəsizdir
- **Qeyd:** c.ch göstəricisi dəyişmir, ona görə pointer receiver yalnız
  dəyişən sahələr üçün lazımdır

### 8. Genişləndirmə İdeyaları
- **Bir tip + opt-in safety:** NewSet (safe default) + NewUnsafeSet
  (performans üçün açıq istək) — "err on the side of safety"
- **Əlavə əməliyyatlar:** Len, Delete, Equals, IsSubset, Subtract,
  Difference, Make (pre-allocate), Map/Reduce/Filter üzvlər üzərində
- **Immutable set:** Add dəyişdirilmiş YENİ set qaytarır (functional stil)

## Əsas terminlər
- Data race (data yarışı) — konkurent oxu+yazı → təyinolunmaz davranış
- Mutex (mutual exclusion) — qarşılıqlı istisna kilidi
- RWMutex — read lock (paylaşılan) + write lock (eksklüziv)
- Contention (rəqabət) — lock uğrunda gözləmə, işin itirilməsi
- Deadlock (ölüm kilidi) — dövrəvi gözləmə, proqramın donması
- Race detector (-race) — data race avtomatik aşkarlayıcısı
- sync/atomic — mutex-siz atomik əməliyyatlar (Uint64 və s.)
- Guard goroutine — data-nın yeganə sahibi; digərləri channel ilə danışır

## Praktik nəticə

1. **Konkurent kodda İLK addım:** `go test -race` — race detector bütün
   testlərdə aktiv olmalıdır.
2. **Mutex saxlamaq üçün struct-da pointer:** `mutex *sync.RWMutex` —
   kopyalanma qorunur; konstruktorda `&sync.RWMutex{}` ilə init.
3. **Lock seçimi:** yazı → Lock; oxu → RLock; hər ikisi defer Unlock ilə.
4. **Sayac kifayətdirsə atomic seç:** atomic.Uint64.Add/Load — mutex
   paperwork olmadan konkurrent sayaç.
5. **Lock qranulyarlığı trade-off:** döngü başına 1 lock (sürətli, amma
   uzun blok) vs əməliyyat başına lock (cəlil, amma churn) — yükdən asılı.

## Mənbə
Pages: 139-155 (PDF 140-156)
