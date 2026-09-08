# Chapter 7 — Caching with generics (səh. 282-327)

## Bu chapter nədən bəhs edir?

Generics əsasları (type parametrlər, constraints, `comparable`/`any`), generic
cache layihəsi, goroutine/channel/sync paketi, paralel testlər, `go test -race`,
Mutex/RWMutex ilə thread-safe cache, TTL və maksimum ölçü.

## Əsas fikirlər

### 1. Generics nədir və nə vaxt LAZIM deyil
Type parametrləri təkrarlanan kodu (hər tip üçün ayrı funksiya) aradan qaldırır.
Amma **lazım olmadıqda istifadə etmə** — `interface{}` (any) kimi tipləri
gizlətmək üçün YOX, konkret tip məcburiyyəti yoxdursa adi interface bəsdir.

### 2. Type parametr sintaksisi
```go
func prettyPrint[T any](t T) {        // [T any] → type parametr bloku
    fmt.Printf("> %v", t)
}
prettyPrint(.25)       // T = float64 (inference / çıxarım)
prettyPrint("pockets") // T = string
```
**Sub-kod izahı:**
- Kvadrat mötərizə funksiya adı ilə parametrlər arasındadır
- `any` = `interface{}` alias — heç bir məhdudiyyət yoxdur
- Call-da tipi yazmaq lazım deyil — inference edir

**Custom constraint:**
```go
type Number interface {
    int | float64 | uint // union — bu tiplərdən biri
}
func Sum[T Number](nums ...T) T
```

### 3. Generic Cache strukturu
```go
type Cache[K comparable, V any] struct { // K map açarı ola bilər
    data map[K]V
}

func New[K comparable, V any]() *Cache[K, V] {
    return &Cache[K, V]{data: make(map[K]V)}
}

func (c *Cache[K, V]) Upsert(key K, value V) { // Upsert = insert + update
    c.data[key] = value
}

func (c *Cache[K, V]) Read(key K) (V, bool) { // "comma ok" semantikası
    v, ok := c.data[key]
    return v, ok
}

func (c *Cache[K, V]) Delete(key K) { // idempotent silmə
    delete(c.data, key)
}
```
- `comparable` → map açarı üçün məcburi şərt (== müqayisəsi)
- Konvensiya: K (key), V (value) təkhərflər
- Delete idempotentdir — olmayan açarı silmək xəta deyil

### 4. Goroutine-lər və kanallar
```go
go printEverySecond("hello") // go → yeni goroutine (yüngül "thread")
```
- **Concurrency (paralellik) vs parallelism (paralelizm):** goroutine-lər
  eyni vaxtda İDARƏ oluna bilən tapşırıqlardır; parallelizm çox nüvədə eyni
  anda İCRA olunmasıdır
- Go atalar sözü: **"Don't communicate by sharing memory; share memory by
  communicating"** (yaddaş paylaşmaqla ünsiyyət yox, ünsiyyətlə yaddaş paylaş)

**Channel əsasları:**
```go
ch := make(chan string)     // unbuffered — bloklanır
ch := make(chan string, 1)  // buffered — 1 mesaja qədər gözləməz
ch <- "msg"                 // yaz
msg := <-ch                 // oxu
close(ch)                   // yazma bağlı → range bitir
func read(c <-chan string)  // yalnız-oxu; send: chan<- (yön göstərir)
```

### 5. sync.WaitGroup — gözləmə nöqtəsi
```go
wg := &sync.WaitGroup{}
wg.Add(2)             // 2 tapşırıq qeydə al
go cookRice(wg)       // hər funksiya daxilində: defer wg.Done()
go cookCurry(wg)
wg.Wait()             // hamısı Done olana qədər blokla
```
**Vacib:** WaitGroup-i pointer ilə ötür (dəyərlə kopyalanarsa Done
itkin düşər).

### 6. errgroup (golang.org/x/sync) — xətalı paralel iş
```go
g := errgroup.Group{}
g.Go(func() error { return cookCurry() }) // xəta qaytarırsa
err := g.Wait()                            // ilk xəta burada
```
WaitGroup-dən fərli: xətaları toplayıb qaytarır.

### 7. Race detection
```go
// paralel test:
t.Run("write six", func(t *testing.T) {
    t.Parallel() // subtest-lər paralel icra
    c.Upsert(6, "six")
})

go test -race . // → WARNING: DATA RACE — map eyni anda yazılır!
```
`--trimpath` → çıxışda modul-relative yollar.

### 8. Mutex və RWMutex
```go
type Cache[K comparable, V any] struct {
    mu   sync.RWMutex     // oxu çox, yaz az → RWMutex
    data map[K]V
}

func (c *Cache[K, V]) Read(key K) (V, bool) {
    c.mu.RLock()          // oxucular paralel girə bilər
    defer c.mu.RUnlock()
    // ...
}
func (c *Cache[K, V]) Upsert(key K, value V) {
    c.mu.Lock()           // yazıçı TİKİLİ təcrid
    defer c.mu.Unlock()
    // ...
}
```
- `defer Unlock()` — early return/panic-də unudulmasın
- Atalar sözü: **"Channels orchestrate; mutexes serialize"** (kanallar
  idarə edir, mutexlər ardıcıllaşdırır) — 90% hallarda kanal düzgün alətdir

### 9. TTL (time to live) + max size
```go
type entryWithTimeout[V any] struct {
    value      V
    bestBefore time.Time // "istifadə tarixi"
}

func (c *Cache[K, V]) Read(key K) (V, bool) {
    // ... entry varsa:
    //   bestBefore keçibsə → zero dəyər + false qaytar
}
```
- Zero value: `var zeroV V` — generic tipin sıfır dəyəri
- Max size: FIFO slice ilə ən köhnə açarı çıxar; update-də mövqe yenilənir
- Alternativ data strukturu: binary search tree; slice əksər ehtiyaclara bəsdir

## Əsas terminlər

- Type Parameter (tip parametri)
- Constraint (məhdudiyyət) / `comparable` / `any`
- Type Inference (tip çıxarımı)
- Goroutine / Channel
- Buffered Channel (buferli kanal)
- sync.WaitGroup
- errgroup
- Data Race / `go test -race`
- Mutex / RWMutex
- TTL (yaşam müddəti) / Idempotent Upsert/Delete

## Praktik nəticə

- Generics: ancaq real tip-müxtəlifliyi varsa; `K comparable, V any` kombinasiyası
  cache/dictionary-lərin standartı
- Paralel test → `t.Parallel()` + `-race` kombinasiyası ilə data race yaxala
- Oxu-çox pattern → RWMutex; yazı dominant → Mutex
- WaitGroup pointer ilə; Done hər goroutine-də defer ilə

## Mənbə

Pages: 282-327 (Chapter 7, Learn Go with Pocket-Sized Projects)
