# Chapter 14 — Effective Coding Practices (Effektiv Kodlaşdırma Praktikaları)

## Bu chapter nədən bəhs edir?
Resursların effektiv idarəsinə: sync.Pool (obyekt yenidən istifadəsi), sync.Once + OnceFunc/OnceValue/OnceValues (tək icra), singleflight (paralel çağırışların deduplikasiyası), mmap (fayl↔yaddaş xəritələmə) və klassik performans tələləri (time.After leak-i, defer loop-da, map-in kiçilməməsi, açıq resurslar, channel idarəsizliyi).

## Əsas fikirlər

### 1. sync.Pool — obyekt hovuzu
Müvəqqəti, tez-tez yaradılan obyektləri yenidən işlət: alloc/GC overhead-i azalır. Hovuzdakı obyektlər GC tərəfindən istənilən an silinə bilər — pool cache-dir, qarantiya deyil.

**BufferPool wrapper:**
```go
type BufferPool struct {
    pool sync.Pool
}
func NewBufferPool() *BufferPool {
    return &BufferPool{
        pool: sync.Pool{
            New: func() interface{} { return new(bytes.Buffer) },  // pool boşsa
        },
    }
}
func (bp *BufferPool) Get() *bytes.Buffer { return bp.pool.Get().(*bytes.Buffer) }
func (bp *BufferPool) Put(buf *bytes.Buffer) {
    buf.Reset()                    // təmizlə — data qalıntısı yox
    bp.pool.Put(buf)
}
func ProcessData(data []byte, bp *BufferPool) {
    buf := bp.Get()
    defer bp.Put(buf)             // mütləq geri qaytar
    buf.Write(data)
    fmt.Println(buf.String())
}
```
**Sub-kod izahı:** `New` boş pool-da çağırılır; `Put`-dan əvvəl `Reset()` — köhnə məzmun növbəti istifadəçiyə keçmir; `defer` ilə qaytarma zəmanətlənir.

**3 ssenari:** bytes.Buffer (yaddaş bufferi), şəbəkə serverində 1KB oxu bufferi (min 1024 bayt `make([]byte, 1024)` ilə — hər bağlantı üçün pool-dan), JSON marshaling (encoder bufferə yazılır → nəticəni `copy` ilə yeni slice-a köçür → buffer pool-a qaytar).

**Qaydalar:** yaratmaı bahalı obyektlər üçün; qısa ömürlü obyektlər üçün (long-lived üçün YOX); GC memory pressure-də pool-u təmizləyə bilər.
**Tələlər:** data qalıntısı (Reset lazımdır), yaddaş şişməsi, sync overhead — benchmark/profile olmadan işlətmə.

### 2. sync.Once — dəfəyə icra
İçində bool + mutex saxlayır: ilk `Do(f)` f-i icra edir, sonrakı bütün çağırışlar (başqa goroutine-lərdən də) yalnız gözləyir — f bir də icra olunmur.

```go
var once sync.Once
func setup() { fmt.Println("Initializing...") }
once.Do(setup)
once.Do(setup)     // icra olunmur
```

**Go 1.21 qısayolları:**
```go
// Klassik:
var once sync.Once
var config *Config
func getConfig() *Config {
    once.Do(func() { config = loadConfig() })
    return config
}
// OnceValue (generic — dəyər də saxlayır):
var getConfig = sync.OnceValue(func() *Config {
    return loadConfig()
})
```
| Funksiya | İmza | Nə üçün |
|---|---|---|
| OnceFunc | func() → func() | dəyərsiz, tək icra |
| OnceValue | func() T → func() T | tək dəyər qaytarır |
| OnceValues | func() (T1,T2) → func() (T1,T2) | iki dəyər (məs. error ilə) |

**Mövqe:** sync.Once — sinxronizasiya alətidir, nəticə-paylaşdırma mexanizmi deyil. Nəticə paylaşmaq + paralel zəngləri deduplikasiya etmək lazımdırsa → singleflight.

### 3. singleflight (golang.org/x/sync/singleflight)
Eyni anda gələn eyni açarlı çağırışlardan yalnız BİRİ icra olunur; qalanları gözləyib eyni nəticəni alır.

```go
var g singleflight.Group

fetchData := func(key string) (interface{}, error) {
    time.Sleep(2 * time.Second)               // bahalı əməliyyat simulyasiyası
    return fmt.Sprintf("Data for key %s", key), nil
}
for i := 0; i < 5; i++ {
    go func(i int) {
        result, err, shared := g.Do("my_key", func() (interface{}, error) {
            return fetchData("my_key")        // YALNIZ 1 dəfə icra olunur
        })
        fmt.Printf("Goroutine %d got result: %v (shared: %v)\n", i, result, shared)
    }(i)
}
```
**Sub-kod izahı:** 5 paralel goroutine eyni açarı soruşur; `g.Do` birincinin icrasını gözləyir, hamısı eyni cavabı alır — `shared: true` = nəticə paylaşıldı. Açarbaşına (per-key) deduplikasiya: "alpha" və "beta" ayrı-ayrı icra olunur.

**Ssenarilər:** eyni resurs üçün paralel HTTP request-lərin deduplikasiyası (cache stampede qarşısı), bahalı hesablamaların kəşlənməsi, rate-limited API throttling, tək nümunəli background task.

### 4. mmap — faylı yaddaşa xəritələ
Fayl baytları birbaşa prosesin adres fəzəsində görünür — böyük fayllarda traditional I/O-dan qat-qat sürətli. `golang.org/x/exp/mmap` (cross-platform, syscall-ları abstraksiya edir).

```go
reader, err := mmap.Open(filename)      // ReaderAt qaytarır
defer reader.Close()
fileSize := reader.Len()
data := make([]byte, fileSize)
reader.ReadAt(data, 0)                  // fayl məzmunu slice-a
lastByte := data[fileSize-1]
```
**Değişiklikləri diska məcburi yazmaq:**
```go
data[fileSize-1] = 'A'
syscall.Msync(data, syscall.MS_SYNC)    // синхron flush
// MS_ASYNC — növbəyə qoyur, OS sonradan yazır (crash-də belə OS yazar, özümüz crash etsək yox)
```
**Protection flag-lər:** PROT_READ / PROT_WRITE / PROT_EXEC / kombinasiya (READ|WRITE).
**Mapping flag-lər:** MAP_SHARED (dəyişikliklər faylla paylaşılır — IPC!) / MAP_PRIVATE (dəyişikliklər prosesə şəxsi, diskə yazılmır).

**Nə vaxt:** nəhəng faylların axtarışı/anallizi; proseslərarası shared memory. **Risk:** sinxronizasiya, xəta yoxlaması, yaddaş korlanması — hamısı sənin üzərinə.

### 5. Tələ 1 — time.After yaddaş leak-i
`time.After` hər çağırışda yeni channel + timer yaradır; timer işə düşənə qədər GC-dən azad OLMAZ və dayandırmaq mümkün deyil. Uzun timeout + erkən bitən əməliyyat = yaddaşda yığılır.

```go
// TƏLƏ: timeout := time.After(duration) — select bitəndən sonra da timer yaşayır

// DÜZ: time.NewTimer + Stop
func processWithManualTimer(duration time.Duration) {
    timer := time.NewTimer(duration)
    defer timer.Stop()                    // resurs dərhal azad
    done := make(chan bool)
    go func() { time.Sleep(duration / 2); done <- true }()
    select {
    case <-done:     fmt.Println("Finished processing")
    case <-timer.C:  fmt.Println("Timed out")
    }
}
```

### 6. Tələ 2 — loop içində defer
`defer` funksiya ÇIXANDA icra olunur — loop-da hər iterasiyada yığılır, minlərlə açıq fayl yaddaşda qalır (OOM-a qədər).

```go
// TƏLƏ:
for _, filename := range filenames {
    f, _ := os.Open(filename)
    defer f.Close()                  // hamısı funksiya sonuna yığılır!
}

// DÜZ: iterasiya daxilində açlıq bağla
for _, filename := range filenames {
    f, err := os.Open(filename)
    if err != nil { return err }
    // ... əməliyyatlar ...
    f.Close()                        // dərhal, iterasiya sonunda
}
```

### 7. Tələ 3 — map kiçilmir
`delete()` yaddaşını GERİ QAYTARMIR — runtime sürət üçün bucket strukturasını saxlayır (yenidən yerləşdirmə tezləşsin). Session map-i minlərlə dövrlə şişir → leak.

```go
// Həll: az qalmışsa yeni map-ə köçür
if len(sessions) < len(deletedSessions) {
    newSessions := make(map[string]Session, len(sessions))
    for k, v := range sessions {
        newSessions[k] = v
    }
    sessions = newSessions            // köhnə bucket-lər GC-ə gedir
}
```

### 8. Tələ 4 — bağlanmayan resurslar
GC yalnız yaddaşı idarə edir — fayl/şəbəkə/DB bağlantıları AÇIQ bağlanmalıdır.

```go
// TƏLƏ: f, _ := os.Open(path); return io.ReadAll(f)   // Close yoxdur

// DÜZ:
func readFile(path string) ([]byte, error) {
    f, err := os.Open(path)
    if err != nil { return nil, err }
    defer f.Close()
    return io.ReadAll(f)
}
```
**Xəta halında bağlantı (conditional defer):**
```go
conn, err := net.DialTCP("tcp", nil, addr)
if err != nil { return nil, err }
defer func() {
    if err != nil {          // sonrakı əməliyyat xəta veribsə
        conn.Close()
    }
}()
```

### 9. Tələ 5 — HTTP Body bağlanmır
`resp.Body` bağlanmazsa socket-lər açıq qalır — connection tükənir, yaddaş şişir.

```go
func fetchURL(url string) error {
    resp, err := http.Get(url)
    if err != nil { return err }
    defer resp.Body.Close()          // error yoxlamasından DƏRHAL sonra
    body, err := io.ReadAll(resp.Body)
    if err != nil { return err }
    fmt.Println(string(body))
    return nil
}
```

### 10. Tələ 6 — channel idarəsizliyi
- **Unbuffered + qəbulçu yoxdur** → göndərən goroutine sonsuz bloklanır (leak).
- **Buffered + drain olunmur** → bufferdəki data yaddaşda qalır.
- **Həll:** hər channel üçün qəbulçu təmin et; lazımsız channel-ları bağla (qapalıya göndərmə — panic!); select + timeout/default ilə blokdan çıxış.

```go
// TƏLƏ: func produce(ch chan int) { for i := 0; ; i++ { ch <- i } }  // qəbulçu yoxdursa sonsuz blok

// DÜZ: timeout ilə çıxış
func produce(ch chan int) {
    for i := 0; ; i++ {
        select {
        case ch <- i:                  // göndərildi
        case <-time.After(5 * time.Second):
            return                    // sonsuz gözləmədən qurtul
        }
    }
}
```

**Ümumi qaydalar:** resurs yaradıldıqca dərhal defer Close; acquired-sonrakı xətaları yoxla; static analyzer (məs. Staticcheck) bağlanmayan resursları tutur.

## Əsas terminlər
- sync.Pool — müvəqqəti obyekt hovuzu (New/Get/Put)
- bytes.Buffer / Reset() — qalıntı data qarşısı
- sync.Once / Do — bir dəfəlik icra (bool + mutex)
- OnceFunc / OnceValue[T] / OnceValues[T1,T2] — Go 1.21 qısayolları
- singleflight.Group / g.Do(key, fn) / shared — paralel zəng deduplikasiyası
- golang.org/x/sync — Go komandasının eksperimental paketi
- mmap / mmap.Open / ReadAt / Len — memory-mapped fayl
- syscall.Msync / MS_SYNC / MS_ASYNC — diskə yazma sinxronu
- PROT_READ/WRITE/EXEC, MAP_SHARED/MAP_PRIVATE — mmap flag-ləri
- time.After vs time.NewTimer / timer.Stop() — timer leak
- defer loop-da — yığılma tələsi
- map delete() yaddaşı qaytarmır — re-map köçürmə həlli
- resp.Body.Close() — HTTP socket leak-inin qarşısı
- Unbuffered/buffered channel — bloklanma/drain leak-ləri
- select + time.After/default — blokdan çıxış naxarı

## Praktik nəticə
1. sync.Pool-u yalnız benchmark-la təsdiqlənmiş bahalı alloc-lar üçün işlət; hər Put-dan əvvəl Reset.
2. Bir dəfəlik init → sync.Once/OnceValue; paralel eyni sorğuların nəticəsi → singleflight (cache stampede dərmanı).
3. Uzun timeout-larda time.After yerinə time.NewTimer + defer Stop.
4. Loop daxilində defer YOX — resursu iterasiya içinde bağla və ya funksiya çıxara çək.
5. Çox-silinən map-lərdə periodik re-creation — delete() yaddaşı qaytarmır.
6. os.Open / http.Get / net.Dial — hamısında dərhal defer Close; body oxunmasa belə bağlanmalı.
7. Channel göndərənlərini həmişə qəbulçu/tımeout ilə təmin et — bloklanmış goroutine = leak.

## Mənbə
Pages: 305-330 (PDF səh. 326-351)
