# Chapter 6 — Concurrency in Go (Go-da Paralellik)

## Bu chapter nədən bəhs edir?
Goroutine-lər, channel-lər (buffered, yönlü, close, select), sinxronizasiya (WaitGroup, Mutex), concurrent cache implementasiyasına, concurrency vs parallelism fərqinə, book recommendation paralelliyinə və do's/don'ts praktikalarına.

## Əsas fikirlər

### 1. Niyə concurrency? (Sürət riyaziyyatı)
Google tədqiqatı: bir neçə yüz millisecond gecikmə user engagement azaldır. Sekvensial icra vs overlap edən tapşırıqlar (aşpaz bir yemək bişərkən digərinə başlayır). Go-nun cavabı: **"share by communicating"** — kilidləmə əvəzinə kanallarla kommunikasiya.

### 2. Goroutine-lər
**Nədir:** Runtime idarəli yüngül thread-lər — minlərlə/milyonlarla yaradıla bilər, traditional thread-lərdən qat-qat ucuz.

```go
go func() {
    // concurrent icra
}()
```

**Kitabdan nümunə — paralel janr axtarışı:**
```go
func searchBooksByGenre(genre string, ch chan<- []Book) {
    books := dbSearch(genre)
    ch <- books
}

genres := []string{"Sci-Fi", "Fantasy", "Mystery"}
resultsChannel := make(chan []Book, len(genres))   // BUFFERED — 3 yer

for _, genre := range genres {
    go searchBooksByGenre(genre, resultsChannel)    // 3 paralel sorğu
}
var allBooks []Book
for i := 0; i < len(genres); i++ {
    books := <-resultsChannel                       // nəticələri yığ
    allBooks = append(allBooks, books...)
}
```
**Sub-kod izahı:**
- `chan<- []Book` — yalnız-GÖNDƏRƏN kanal (yönlü tip)
- Buffer = len(genres) → heç bir goroutine göndərmədə bloklanmır
- Nəticə sırası произвольныйdir, hamısı toplanana qədər loop

**Həddindən artıq istifadədən qaçın:** hər kiçik task üçün goroutine YOX — yalnız I/O-bound, data-intensive, CPU-intensive işlərdə (batch emal, paralel axtarış, background tövsiyə).

### 3. Channel-lər
**Nədir:** Goroutine-lər arası data borusu — race condition-dan təbii qorunmuş sinxron mübadilə.

```go
bookChannel := make(chan string)      // yaratma
bookChannel <- "Golang in Action"     // göndər
title := <-bookChannel                // qəbul
```

**Yönlü kanallar (təhlükəsizlik sərhədi):**
```go
func processOrder(ch chan<- Order) { ch <- newOrder }   // yalnız göndər
func billOrder(ch <-chan Order)   { order := <-ch }     // yalnız qəbul
```

**Buffered:** `make(chan int, 50)` — qəbulçu hazır olmasa da 50-ə qədər göndər.

**Close + range:**
```go
close(bookChannel)                   // daha data gəlməyəcək siqnalı
for title := range bookChannel {      // bağlanana qədər iterate
    // hər title emal et
}
```

**select — çoxkanal multiplexing:**
```go
select {
case order := <-ordersChannel:      // sifariş kanalı hazır
case feedback := <-feedbackChannel: // feedback kanalı hazır
default:                            // heç biri hazır — bloklanma YOX
}
```

### 4. Sinxronizasiya — sync.WaitGroup və sync.Mutex

**WaitGroup — bitmə gözləyici sayğac:**
```go
var wg sync.WaitGroup

func fetchBooks(category string) {
    defer wg.Done()                  // bitəndə sayğacı azalt
    // fetching məntiqi
}

for _, category := range []string{"Thriller", "Sci-Fi", "History"} {
    wg.Add(1)                        // hər goroutine üçün artır
    go fetchBooks(category)
}
wg.Wait()                             // hamısı bitənə qədər blokla
```

**Mutex — paylaşılan resurs kilidi:**
```go
var mutex sync.Mutex
var inventory = make(map[string]int)

func updateInventory(bookName string, sold int) {
    mutex.Lock()                      // kritik bölgəyə gir
    inventory[bookName] -= sold
    mutex.Unlock()                    // burax
}
```
**Niyə lazım:** İki goroutine eyni anda `-=` etsə — double decrement (səhs sayı korlanır). Mutex icraları SIRALIŞDIRIR.

**Birləşmiş pattern — bulk sifariş emalı:**
```go
func processOrder(order Order) {
    defer wg.Done()
    for _, item := range order.Items {
        mutex.Lock()
        inventory[item.BookName] -= item.Quantity
        mutex.Unlock()
    }
}
for _, order := range bulkOrders {
    wg.Add(1)
    go processOrder(order)
}
wg.Wait()                              // WaitGroup orkestr, Mutex qoruma
```

### 5. Concurrent Cache implementasiyası
**Kitabdan kod nümunəsi:**
```go
type ConcurrentCache struct {
    data  map[string]*BookDetails
    mutex sync.Mutex
}

func (c *ConcurrentCache) Store(key string, value *BookDetails) {
    c.mutex.Lock()
    c.data[key] = value
    c.mutex.Unlock()
}

func (c *ConcurrentCache) Fetch(key string) (*BookDetails, bool) {
    c.mutex.Lock()
    value, exists := c.data[key]
    c.mutex.Unlock()
    return value, exists
}

const MAX_CACHE_SIZE = 1000
func (c *ConcurrentCache) CheckEviction() {
    c.mutex.Lock()
    if len(c.data) > MAX_CACHE_SIZE {
        for key := range c.data {
            delete(c.data, key)        // sadələşdirilmiş eviction (LRU alternativ)
            break
        }
    }
    c.mutex.Unlock()
}
```
**İşə düşmə axını:** cache-də var → dərhal qaytar; yoxdur → DB-dən al → cache-ə yaz → qaytar. Populyar kitablar tezliklə cache-də olur.

### 6. Concurrency vs Parallelism
| Concurrency (paralellik idarəsi) | Parallelism (eynizamanlılıq) |
|---|---|
| Çox tapşırığın İDARƏ edilməsi | çox tapşırığın EYNİ ANDA icrası |
| Bir yol — maşınlar gir-çıxır, işıqlarda dayanır | Çoxmərtəbəli yol — hər core bir zolaq |
| Tək core-da mümkün | çox core tələb edir |
| Flash sale — minlərlə sorğunun qayğısız idarəsi | 1000 user-in satınalma analizi — çox core bölünür |

**Kitabın parallelism nümunəsi:**
```go
func BenchmarkRecommendations(b *testing.B) {
    users := fetchSampleUsersForBenchmarking()
    b.RunParallel(func(pb *testing.PB) {     // çox core-da parallel
        for pb.Next() {
            for _, user := range users {
                analyzePurchases(user)
            }
        }
    })
}
```

### 7. Do's and Don'ts
| DO | DON'T |
|---|---|
| Concurrency/parallelism fərqini başa düş | Data race-ləri (oxu+yazı üst-üstə) nəzərə alma — Mutex ilə qoru |
| Goroutine-ləri mötədil işlət (I/O, background) | Hər kiçik task üçün goroutine yarat — resurs tükənməsi |
| Goroutine panic-lərini defer-recover ilə tut | Bir goroutine-in paniki bütün app-i çökdürsün |
| Worker pool/buffered channel ilə sayı məhdudlaşdır | Minlərlə goroutine-ni nəzarətsiz burax |
| Graceful shutdown — işlər bitsin, sonra bağlan | Sorğuları kəskin kəs |
| Concurrency testləri yaz (simulyasiya) | Test-sız burax — race bug-ları gizli qalır |

## Əsas terminlər
- Goroutine — yüngül runtime thread
- Channel (kanal) — tipəlyi data borusu
- Buffered channel — tutumlu kanal
- Channel direction — `chan<-` / `<-chan`
- select — çoxkanal gözləmə
- sync.WaitGroup — bitmə sayğacı
- sync.Mutex — qarşılıqlı istisna kilidi
- Race condition (yariş şəraiti) — eynizamanı oxu/yazı zərəri
- LRU — least recently used eviction
- b.RunParallel — benchmark parallelizm aləti

## Praktik nəticə
1. Paralel DB axtarışları: buffered channel + N goroutine + nəticə toplama loop-u.
2. Paylaşılan map/sayğac HƏMİŞƏ Mutex arxasında — `-=` əməliyyatı atomik DEYİL.
3. WaitGroup.Add(1) goroutine BAŞLAMAZDAN əvvəl, defer wg.Done() funksiyanın İÇİNDƏ.
4. Worker pool ilə goroutine sayını sərhədlə — sonsuz artım resurs tükəndirir.
5. Cache-ə yaz/oxu mutex-lə; populyar data avtomatik "isti" qalır.
6. Flash-sale = concurrency (idare), analitik batch = parallelism (bölüşdür).
7. Hər goroutine-də recover — bir sorğunun xətası bütün serveri yıxmasın.

## Mənbə
Pages: 163-185 (PDF səh. 163-185)
