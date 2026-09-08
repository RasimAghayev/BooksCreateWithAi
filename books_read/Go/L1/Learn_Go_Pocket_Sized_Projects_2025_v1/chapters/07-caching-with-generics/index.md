# Chapter 7 — Caching with Generics (səh. 282-327)

## Bu fəsil nədən bəhs edir?

Generic cache layihəsi: generics (type parametr, constraint, inference),
goroutines, channels, sync.WaitGroup / errgroup, race condition aşkarı
(`go test -race`), sync.Mutex/RWMutex, TTL (time to live), LRU-mənalı maxSize.

## Əsas fikirlər

### 1. Generics — əsaslar
```go
func prettyPrint[T any](t T) { fmt.Printf("> %v", t) }
prettyPrint(0.25)     // T = float64 (inference — YALNIZ input-lardan)
prettyPrint("pockets") // T = string
```
- `[T any]` — kvadrat mötərizə, funksiya adından SONRA; `any` = istənilən tip
- **Generic type:**
```go
type Group[T any] []T
func (g Group[T]) PrettyPrint() { ... }
var g Group[Cloud]  // Cloud üçün xüsusi versiya
```
- **Custom constraint:** `comparable` (müqayisəli — map açarı üçün) və ya öz
  interfeysin: `type Number interface { int | int64 | float64 }`

### 2. Cache strukturu
```go
type entryWithTimeout[V any] struct {
    value   V
    expires time.Time
}

type Cache[K comparable, V any] struct {
    ttl time.Duration
    mu  sync.Mutex
    data map[K]entryWithTimeout[V]
    maxSize int
    chronologicalKeys []K  // LRU üçün sıra
}

func New[K comparable, V any](maxSize int, ttl time.Duration) Cache[K, V]
```
- `K comparable` — map açarı; `V any` — dəyər
- Upsert (insert+update), Delete (idempotent — yoxdursa da xəta yox),
  Read (value, ok)

### 3. Goroutines + parallel vs concurrency
```go
go printEverySecond("Hello")  // go — yeni goroutine (KB-lərcə yüngül)
go printEverySecond("World")
```
- **Concurrency:** struktur (aynı anda idarə edilmə); **Parallelism:** eyni anda
  icra (çox nüvə) — ROB PIKE fərqi
- 2 nüvə = 2 dəfə sürət DEYİL (context switch, resurs mübarizəsi)

### 4. Channels — kommunikasiya
```go
c := make(chan int)
c <- 4        // göndər (blocking)
v := <-c      // oxu (blocking)
close(c)      // yazmağı bağla — oxu davam edir (zero value + false)
```
- İstiqamətli imza: `func read(c <-chan string)` / `func write(c chan<- string)`
- **Go proverb:** "Don't communicate by sharing memory; share memory by
  communicating."

### 5. WaitGroup + errgroup
```go
wg := &sync.WaitGroup{}
wg.Add(2)
go cookRice(wg)   // daxilində defer wg.Done()
go cookCurry(wg)
wg.Wait()

// errgroup (golang.org/x/sync) — xətaları toplayır:
g, ctx := errgroup.WithContext(ctx)
g.Go(func() error { return cookRice(ctx) })
err := g.Wait()   // İLK xəta qaytarılır
```
- `wg *sync.WaitGroup` — POINTER məcburidir (kopyası Done saymaz!)

### 6. Race condition — aşkar və həll
```go
func TestCache_Parallel(t *testing.T) {
    c := cache.New[int, string]()
    t.Run("write six", func(t *testing.T) { t.Parallel(); c.Upsert(6, "six") })
    t.Run("write kuus", func(t *testing.T) { t.Parallel(); c.Upsert(6, "kuus") })
}
```
```bash
go test --trimpath -race .   # → WARNING: DATA RACE (runtime.mapassign)
```
- "Testing shows presence of bugs, never their absence" — Dijkstra

**Mutex həlli:**
```go
func (c *Cache[K, V]) Upsert(key K, value V) {
    c.mu.Lock()
    defer c.mu.Unlock()   // HƏMİŞƏ defer — erkən return unutmur
    // map əməliyyatı
}
```
- **RWMutex:** Read → RLock/RUnlock (paralel oxu), Upsert/Delete → Lock/Unlock
- **Go proverb:** "Channels orchestrate; mutexes serialize."

### 7. TTL — vaxt keçmiş dəyər
```go
func (c *Cache[K, V]) Read(key K) (V, bool) {
    c.mu.Lock()
    defer c.mu.Unlock()
    var zeroV V
    e, ok := c.data[key]
    switch {
    case !ok:
        return zeroV, false
    case e.expires.Before(time.Now()):
        delete(c.data, key)          // lazy eviction
        return zeroV, false
    default:
        return e.value, true
    }
}
```
- Upsert yeniləmə = TTL də yenilənir

### 8. maxSize — ən köhnəni at
- `chronologicalKeys []K` — yazılış sırası; yeni element yer yoxdursa →
  ən qədim açarı sil (slice ilə O(1) axtarış mövqeyi)
- Mövcud açar yenilənərsə → köhnə mövqeyindən çıxarıb SONA qoyur (yenilənmiş
  = "ən gənc")

### 9. Kanal QADAĞALARI (Go proverbs)
- Ehtiyac yoxdursa kanal İŞLƏTMƏ — learning curve var
- Goroutine sızdırma: hamısı bitmədən proqram çıxmamalı; yazış bitəndə
  close() et
- `--trimpath` — race reportda modul-kökə nisbi yollar

## Əsas terminlər

- Generic (tip parametri)
- Type Constraint (any, comparable)
- Type Inference (çıxarış — yalnız input-lardan)
- Goroutine (yüngül thread)
- Channel (kanal) / close
- sync.WaitGroup (bitmə gözləmə)
- errgroup (xətalı gözləmə)
- Data Race (məlumat yarışı)
- sync.Mutex / RWMutex (qarşılıqlı istisna)
- TTL (yaşam müddəti)
- LRU (ən az istifadə)

## Praktik nəticə

- Map + paralel yazış = mütləq mutex; `go test -race` MƏCBURİ
- Oxu çoxdursa RWMutex; kiçik critical section üçün Mutex sadədir
- TTL: hər entry-ə expires; Read-də lazımsızlaşanı sil (lazy)
- Kanal ehtiyac olmadıqça sadə mutex üstündür

## Mənbə

Pages: 282-327 (Chapter 7, Learn Go with Pocket-Sized Projects)
