# Fəsil 9 — Go concurrency-dən istifadə

## Bu fəsil nədən bəhs edir?

Concurrency vs parallelism fərqi, goroutine-lər (istifadə, performans benchmark-ləri, WaitGroup), channel-lər (sinxronizasiya, mesaj ötürmə, buffered, select, close), və real web tətbiqi: foto-mozaika generatorunun concurrent versiyası (fan-out/fan-in pattern, mutex, race condition).

## Əsas fikirlər

### 1. Concurrency ≠ Parallelism
**Təriflər:**
- **Concurrency:** 2+ tapşırıq eyni dövrdə başlayır/işləyir/bitir; icraları üst-üstə düşür; planlaşdırılır və bir-biri ilə əlaqələnə bilər
- **Parallelism:** tapşırıqlar həqiqətən EYNİ ANDA icra olunur; böyük problem hissələrə bölünüb paralel emal olunur; müstəqil resurslar (CPU) tələb edir

> "Concurrency is about dealing with lots of things at once. Parallelism is about doing lots of things at once." — Rob Pike

**Market analogiyası:**
- Concurrency = 2 növbə, 1 kassa (növbələr növbə ilə xidmət olunur)
- Parallelism = 2 növbə, 2 kassa (eyni anda 2 müştəri)

**Go-da:** GOMAXPROCS > 1 = paralellik mümkün (Go 1.5-dən default = CPU sayı; əvvəl 1 idi). Concurrent proqram tək CPU-də də işləyir. **Go concurrency üçün yaradılıb, parallelizm üçün yox.**

### 2. Goroutine-lər
**Nədir:** Müstəqil işləyən funksiyalar; thread-lər üzərində multiplex olunur, amma thread deyillər. Kiçik stack (Go 1.4-də 8KB), böyüyüb/kiçilə bilir. Bloklandıqda runtime eyni thread-dəki digər goroutine-ləri başqa thread-ə daşıyır.

**İstifadə:** funksiyanın qarşısına `go` açar sözü (adlı və ya anonim).

**Kitabdan kod nümunəsi:**
```go
func printNumbers1() {
    for i := 0; i < 10; i++ {
        fmt.Printf("%d ", i)
    }
}

func printLetters1() {
    for i := 'A'; i < 'A'+10; i++ {
        fmt.Printf("%c ", i)
    }
}

func print1() {
    printNumbers1()     // ardıcıl
    printLetters1()
}

func goPrint1() {
    go printNumbers1()  // goroutine
    go printLetters1()
}
```

**Vacib müşahidə:** goroutine-lər başladıqdan sonra main (və ya test) bitərsə — **onlar icra olunmur!** `time.Sleep(1 * time.Millisecond)` ilə gözləmə hack-dir (düzgün həll: WaitGroup/channel).

**İş yükü ilə fərq görünür:** hər iterasiyada 1µs Sleep olduqda output qarışır:
```
A 0 B 1 C D 2 E 3 F 4 G H 5 I 6 J 7 8 9
```
Hər run-da fərqli nəticə — scheduler qeyri-deterministikdir. `-cpu 1` → təkrarlanan nəticələr.

### 3. Goroutine performansı — benchmark dərsləri
**Tək CPU (`-cpu 1`):**
```
BenchmarkPrint1      100000000      13.9 ns/op    ← ardıcıl
BenchmarkGoPrint1     1000000    1090 ns/op       ← goroutine: ~78x YAVAŞ!
```
**Dərs 1:** Goroutine başlatmaq PULSUZ DEYİLDİR — kiçik/trivial işlərdə overhead üstünlük təşkil edir.

**1µs iş yükü ilə:**
```
BenchmarkPrint2        10000   121384 ns/op       ← ardıcıl
BenchmarkGoPrint2    1000000    17206 ns/op       ← goroutine: ~7x SÜRƏTLİ
```

**100 iterasiya (10x iş):**
```
BenchmarkPrint2      2000    1184572 ns/op        ← ~10x yavaşladı (proqnoz)
BenchmarkGoPrint2  1000000      17564 ns/op       ← demək olar dəyişmədi!
```
**Dərs 2:** İş böyüdükcə goroutine overhead-i əhəmiyyətsizləşir.

**2 CPU (`-cpu 2`):**
```
BenchmarkGoPrint2-2   200000     8607 ns/op      ← 2x sürətləndi (iş paylaşılır)
```

**4 CPU (`-cpu 4`):**
```
BenchmarkGoPrint1-4  3000000       479 ns/op      ← 2 CPU-dan PİS! (overhead > qazanc)
BenchmarkGoPrint2-4   300000      6193 ns/op      ← cəmi 40% yaxşılaşdı
```
**Dərs 3 (ən vacib):** CPU sayını artırmaq avtomatik performans demək DEYİL — çox-CPU planlaşdırma overhead-i var. Benchmark etmədən qərar vermə.

### 4. WaitGroup — goroutine gözləmə
**Addımlar:** WaitGroup elan et → `Add(n)` sayğac qoy → goroutine-lər `Done()` çağırır → `Wait()` sıfırlanana qədər bloklayır.

**Kitabdan kod nümunəsi:**
```go
func printNumbers2(wg *sync.WaitGroup) {
    for i := 0; i < 10; i++ {
        time.Sleep(1 * time.Microsecond)
        fmt.Printf("%d ", i)
    }
    wg.Done()     // sayğacı azalt
}

func main() {
    var wg sync.WaitGroup     // 1. elan
    wg.Add(2)                 // 2. sayğac = 2
    go printNumbers2(&wg)     // 3. goroutine-lər
    go printLetters2(&wg)
    wg.Wait()                // 4. gözlə
}
// 0 A 1 B 2 C 3 D 4 E 5 F 6 G 7 H 8 I 9 J
```

**Deadlock xəbərdarlığı:** `Done()` çağırılanma goroutine olsa:
```
fatal error: all goroutines are asleep - deadlock!
```

### 5. Channel-lər — əsaslar
**Nədir:** Goroutine-lər arasında kommunikasiya üçün tipli dəyərlər. `make` ilə yaradılır, reference dəyərdir. **Qutu metaforası:** bir goroutine qoyur, digəri götürür.

**Yaradılış:**
```go
ch := make(chan int)       // unbuffered (default)
ch := make(chan int, 10)   // buffered, size 10
ch := make(chan <- string) // send-only
ch := make(<-chan string)  // receive-only
```

**Əməliyyatlar:**
```go
ch <- 1     // qoy (send)
i := <-ch   // götür (receive)
```

**Unbuffered = sinxron:** qutu 1 şey saxlayır; goroutine bir şey qoyanda **başqa heç nə qoyula bilməz** — biri götürənə qədər qoyan yuxuya gedir; boş qutudan götürmək də bloklanır.

### 6. Sinxronizasiya channel ilə (WaitGroup əvəzinə)
```go
func printNumbers2(w chan bool) {
    for i := 0; i < 10; i++ { ... }
    w <- true          // bitəndə siqnal
}

func main() {
    w1, w2 := make(chan bool), make(chan bool)
    go printNumbers2(w1)
    go printLetters2(w2)
    <-w1               // bloklanır — w1-ə true gələnə qədər
    <-w2               // bloklanır — w2-ə true gələnə qədər
}
```
Dəyərdən istifadə olunmur — yalnız unblocking məqsədi.

### 7. Mesaj ötürmə — thrower/catcher
```go
func thrower(c chan int) {
    for i := 0; i < 5; i++ {
        c <- i
        fmt.Println("Threw  >>", i)
    }
}

func catcher(c chan int) {
    for i := 0; i < 5; i++ {
        num := <-c
        fmt.Println("Caught <<", num)
    }
}

func main() {
    c := make(chan int)
    go thrower(c)
    go catcher(c)
    time.Sleep(100 * time.Millisecond)
}
```
Output: Threw/Caught sıraları qarışsa da **nömrələr ardıcıllıqdır** — hər atılan nömrə tutulmadan növbətisi atılmır (unbuffered sinxronluq).

### 8. Buffered channel — FIFO növbə
```go
c := make(chan int, 3)   // 3-lük buffer
```
Output dəyişir:
```
Threw  >> 0
Threw  >> 1
Threw  >> 2      ← 3 atıldı, buffer doldu
Caught << 0      ← catcher boşaldır
...
```
- Dolulana qədər send bloklanmır; boşaldıqdan sonra bloklanır
- **İstifadə:** throughput məhdudlaşdırma (throttling) — proses sayı məhduddursa sorğuları yavaşa t.

### 9. select — çox kanal seçimi
**Nədir:** switch-in channel versiyası — bir çox kanaldan hansı hazırdırsa onu seçir.

```go
select {
case msg := <-a:
    fmt.Printf("%s from A\n", msg)
case msg := <-b:
    fmt.Printf("%s from B\n", msg)
}
```
- Hər ikisi hazırdırsa → **random seçim**
- Heç biri yoxdursa → bloklanır

**Deadlock ssenarisi:** 5 dəfə loop + yalnız 2 mesaj → 3-cü iterasiyada hər iki kanal bloklanır → `fatal error: all goroutines are asleep - deadlock!`

**Həll 1 — default:**
```go
select {
case msg := <-a: ...
case msg := <-b: ...
default:
    fmt.Println("Default")   // hər ikisi blokda → bu işləyir
}
```

**Həll 2 — close (düzgün yol):**
```go
func callerA(c chan string) {
    c <- "Hello World!"
    close(c)          // mesaj bitdi siqnalı
}

func main() {
    a, b := make(chan string), make(chan string)
    go callerA(a)
    go callerB(b)

    var msg string
    ok1, ok2 := true, true
    for ok1 || ok2 {
        select {
        case msg, ok1 = <-a:     // 2 dəyərli receive!
            if ok1 {
                fmt.Printf("%s from A\n", msg)
            }
        case msg, ok2 = <-b:
            if ok2 {
                fmt.Printf("%s from B\n", msg)
            }
        }
    }
}
```

**close qaydaları:**
- Qapalı kanal: heç vaxt bloklanmır, **zero value** qaytarır (`ok=false` ilə)
- Receive-only kanalı bağlamaq OLMAZ; bağlanmış kanalı bağlamaq/bağlanmış kanala send → **panic**
- Bağlamaq MƏCBURİ DEYİL — yalnız "artıq heç nə gəlməyəcək" siqnalıdır

### 10. Web tətbiqi — foto-mozaika (real nümunə)
**Alqoritm (5 addım):** tile DB yarat (qovluq skan + orta rəng) → target şəkli tile ölçüsünə kəs → hər hissənin sol-üst pikselini orta rəng hesabla → DB-dən ən yaxın (Euclid məsafəsi) tap → yerləşdir + DB-dən sil (unikallıq).

**Əsas funksiyalar:**
```go
func averageColor(img image.Image) [3]float64 { ... }    // RGB orta
func resize(in image.Image, newWidth int) image.NRGBA { ... }
func tilesDB() map[string][3]float64 { ... }             // tile bazası
func nearest(target [3]float64, db *map[string][3]float64) string { ... }
func distance(p1, p2 [3]float64) float64 {                // Euclid
    return math.Sqrt(sq(p2[0]-p1[0]) + sq(p2[1]-p1[1]) + sq(p2[2]-p1[2]))
}
var TILESDB map[string][3]float64       // qlobal — 1 dəfə dolur
func cloneTilesDB() map[string][3]float64 { ... }        // hər sorğuda klon
```

**Sinxron versiya handler-i:** nested loop → hər tile üçün nearest → draw; nəticə: **2.25 saniyə** (151KB JPEG).
Şəkil HTML-ə **data URL** (`data:image/jpg;base64,...`) kimi daxil edilir.

### 11. Concurrent versiya — fan-out/fan-in
**Alqoritm:** şəkli 4 kvadranta böl → eyni anda emal et → nəticələri birləşdir.

**mosaic handler:**
```go
c1 := cut(original, &db, tileSize, bounds.Min.X, bounds.Min.Y, bounds.Max.X/2, bounds.Max.Y/2)
c2 := cut(original, &db, tileSize, bounds.Max.X/2, bounds.Min.Y, bounds.Max.X, bounds.Max.Y/2)
c3 := cut(original, &db, tileSize, bounds.Min.X, bounds.Max.Y/2, bounds.Max.X/2, bounds.Max.Y)
c4 := cut(original, &db, tileSize, bounds.Max.X/2, bounds.Max.Y/2, bounds.Max.X, bounds.Max.Y)
c := combine(bounds, c1, c2, c3, c4)
// ...
"mosaic": <-c,     // kanaldan nəticə gözlə
```

**Race condition problemi:** 4 goroutine eyni tile-ı eyni anda "ən yaxın" tapa bilər (DB-dən silinmədən əvvəl).

**Həll — mutex + DB struct:**
```go
type DB struct {
    mutex *sync.Mutex
    store map[string][3]float64
}

func (db *DB) nearest(target [3]float64) string {
    var filename string
    db.mutex.Lock()          // kritik bölmə: AXTARIŞ + SİLMƏ birlikdə!
    smallest := 1000000.0
    for k, v := range db.store {
        dist := distance(target, v)
        if dist < smallest {
            filename, smallest = k, dist
        }
    }
    delete(db.store, filename)
    db.mutex.Unlock()
    return filename
}
```
**Vacib:** yalnız `delete`-i kilidləmək KİFAYƏT DEYİL — başqa goroutine silinmədən əvvəl eyni tile-ı tapa bilər; axtarış+silmə BİRLİKDƏ kilidlənməlidir.

**cut — fan-out (receive-only channel qaytarır):**
```go
func cut(original image.Image, db *DB, tileSize, x1, y1, x2, y2 int) <-chan image.Image {
    c := make(chan image.Image)
    sp := image.Point{0, 0}
    go func() {
        newimage := image.NewNRGBA(image.Rect(x1, y1, x2, y2))
        for y := y1; y < y2; y = y + tileSize {
            for x := x1; x < x2; x = x + tileSize {
                r, g, b, _ := original.At(x, y).RGBA()
                color := [3]float64{float64(r), float64(g), float64(b)}
                nearest := db.nearest(color)
                file, err := os.Open(nearest)
                if err == nil {
                    img, _, err := image.Decode(file)
                    if err == nil {
                        t := resize(img, tileSize)
                        tile := t.SubImage(t.Bounds())
                        tileBounds := image.Rect(x, y, x+tileSize, y+tileSize)
                        draw.Draw(newimage, tileBounds, tile, sp, draw.Src)
                    }
                }
                file.Close()
            }
        }
        c <- newimage.SubImage(newimage.Rect)
    }()
    return c    // kanal dərhal qayıdır; nəticə hazır olanda gələcək
}
```
- Kanal **bidirectional yaradılır, receive-only qaytarılır** (typecast)
- Anonim goroutine işi bitirəndə kanala şəkil göndərir

**combine — fan-in (select + WaitGroup):**
```go
func combine(r image.Rectangle, c1, c2, c3, c4 <-chan image.Image) <-chan string {
    c := make(chan string)

    go func() {
        var wg sync.WaitGroup
        img := image.NewNRGBA(r)
        copy := func(dst draw.Image, r image.Rectangle, src image.Image, sp image.Point) {
            draw.Draw(dst, r, src, sp, draw.Src)
            wg.Done()
        }
        wg.Add(4)
        var s1, s2, s3, s4 image.Image
        var ok1, ok2, ok3, ok4 bool
        for {
            select {
            case s1, ok1 = <-c1:
                go copy(img, s1.Bounds(), s1, image.Point{r.Min.X, r.Min.Y})
            case s2, ok2 = <-c2:
                go copy(img, s2.Bounds(), s2, image.Point{r.Max.X / 2, r.Min.Y})
            case s3, ok3 = <-c3:
                go copy(img, s3.Bounds(), s3, image.Point{r.Min.X, r.Max.Y/2})
            case s4, ok4 = <-c4:
                go copy(img, s4.Bounds(), s4, image.Point{r.Max.X / 2, r.Max.Y / 2})
            }
            if (ok1 && ok2 && ok3 && ok4) {
                break
            }
        }

        wg.Wait()     // bütün kopyalar bitənə qədər
        buf2 := new(bytes.Buffer)
        jpeg.Encode(buf2, img, nil)
        c <- base64.StdEncoding.EncodeToString(buf2.Bytes())
    }()
    return c
}
```

**Sub-kod izahı:**
- `select` → hansı kvadrant HAZIRDISA onu emal et (birincini gözləmək concurrency deyil!)
- `go copy(...)` → kopyalamanı da paralelə burax
- 2 dəyərli receive (`s1, ok1 = <-c1`) → kanal bağlanıbmu yoxla; hamısı ok → loop-dan çıx
- `wg.Wait()` → 4 kopya goroutine bitəndə gözlə; sonra encode + string kanala

**Performans nəticələri:**
| Versiya | Vaxt | Qazanc |
|---|---|---|
| Sinxron (1 CPU) | 2250 ms | — |
| Concurrent (1 CPU) | 646 ms | 3.5x — **paralelliksiz!** |
| Concurrent (multi-CPU) | 216 ms | 10x |

**Concurrent (1 CPU) = parallelizm YOXDUR** — goroutine-lər müstəqil işləyir amma bir-birini əvəz edir. Multi-CPU-da **paralellik pulsuz gəlir**.

## Concurrency pattern xülasəsi
| Pattern | Nə edir | Nümunə |
|---|---|---|
| Fan-out | işi parçalayıb paralel goroutine-lərə payla | cut x4 |
| Fan-in | nəticələri bir yerə topla | combine |
| Mutex | paylaşılan resursa eksklüziv çıxış | DB.nearest |
| WaitGroup | goroutine-lərin bitməsini gözlə | copy x4 |
| Channel siqnal | bitiş bildirişi | `w <- true` |
| close + 2 dəyərli receive | mənbək bitmə yoxlaması | `msg, ok = <-a` |

## Əsas terminlər
- Concurrency / Parallelism (eyni-vaxtlılıq / paralellik)
- GOMAXPROCS (CPU sayı konfiqi)
- Goroutine (qorutin)
- Multiplexing (çoxləmə — thread-lər üzərində)
- WaitGroup (gözləmə qrupu): Add/Done/Wait
- Deadlock (qarşılıqlı bloklanma)
- Channel / Unbuffered (Sinxron) / Buffered (FIFO)
- Send-only / Receive-only (`chan <-` / `<-chan`)
- select / default case
- close / two-value receive (`v, ok := <-ch`)
- Zero Value (bağlı kanalın qaytardığı)
- Race Condition (yarış şəraiti)
- Mutual Exclusion / Mutex (qarşılıqlı istisna)
- Critical Section (kritik bölmə)
- Fan-out / Fan-in (paylama / toplama)
- Data URL (base64 şəkil embed)
- Euclidean Distance (Evklid məsafəsi)

## Praktik nəticə
- Goroutine-ləri yalnız **kifayət qədər ağır** işlər üçün başlat — trivial işlərdə overhead 78x sürət itkisi verir.
- Sleep ilə gözləmə hack-dir — WaitGroup və ya channel siqnalından istifadə et.
- Paylaşılan resurs + çox goroutine = race condition; axtarış+dəyişməni eyni kritik bölmədə kilidlə.
- Sonsuz select loop-larında `close` + `v, ok := <-ch` — default sonsuz boşluğa düşür.
- Concurrent dizayn et — paralellik sonra GOMAXPROCS ilə **pulsuz** gəlir (kitabda 10x sübut olundu).
- Benchmark olmadan CPU sayı haqqında qərar vermə — 4 CPU bəzən 2 CPU-dan pisdir.

## Mənbə
Pages: 244-276 (PDF), book pages 223-254
