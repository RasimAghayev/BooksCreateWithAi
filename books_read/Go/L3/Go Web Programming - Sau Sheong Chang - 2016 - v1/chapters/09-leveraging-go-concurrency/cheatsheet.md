# Chapter 9 — Leveraging Go concurrency Cheatsheet

## `go func()` — goroutine başlatma

**Nə edir:** Funksiyanı müstəqil icra vahidi kimi başladır. `main` bitərsə goroutine-lər də ölür.

```go
go printNumbers2()        // adlı funksiya
go func() { /* ... */ }() // anonim funksiya
```

**Mənbə:** Chapter 9, page 226

---

## `sync.WaitGroup` — counter əsaslı sinxronizasiya

**Nə edir:** Bütün goroutine-lərin bitməsini gözləyir.

```go
var wg sync.WaitGroup
wg.Add(2)
go printNumbers2(&wg)
go printLetters2(&wg)
wg.Wait()  // counter 0 olana qədər blok edir
```

**Parametrlər:**
- `Add(n)` — counter +n
- `Done()` — counter -1
- `Wait()` — counter 0 olana qədər blok edir

**Mənbə:** Chapter 9, page 232

---

## `make(chan T)` və `make(chan T, N)` — channel yaratmaq

**Nə edir:** Typed message channel yaradır. Unbuffered sinxron, buffered N ölçüdə FIFO.

```go
ch := make(chan int)         // unbuffered
ch := make(chan int, 10)     // buffered (10 slot)
sendOnly := make(chan<- string)  // yalnız göndərmə
recvOnly := make(<-chan string)  // yalnız alma
```

**Mənbə:** Chapter 9, page 232

---

## `ch <- val` və `val := <-ch` — göndərmə və alma

**Nə edir:** Channel vasitəsilə mesaj ötürür. Unbuffered hər iki tərəf hazır olmayana qədər blok edir.

```go
ch <- 1                  // göndər
i := <-ch                // al
val, ok := <-ch          // ok=false → channel bağlıdır
```

**Mənbə:** Chapter 9, page 233

---

## `select { case ... : ... default: ... }` — çoxsaylı channel seçimi

**Nə edir:** Hazır olan channel case-ini seçir (rastgələ əgər birdən çoxu hazırdırsa). `default` bloklanmır.

```go
select {
case msg := <-a:
    fmt.Println("from A:", msg)
case msg := <-b:
    fmt.Println("from B:", msg)
default:
    fmt.Println("nothing ready")
}
```

**Mənbə:** Chapter 9, page 237

---

## `close(ch)` — channel bağlama

**Nə edir:** Channel-i bağlayır. Almaq zero value qaytarır, artıq göndərmə panic verir.

```go
ch <- "msg"
close(ch)  // bundan sonra almaq mümkündür, göndərmə yox
val, ok := <-ch  // val="", ok=false
```

**Mənbə:** Chapter 9, page 239

---

## `sync.Mutex` — race condition qarşısının alınması

**Nə edir:** Kritik hissəyə yalnız bir goroutine-in girməsini təmin edir.

```go
type DB struct {
    mutex *sync.Mutex
    store map[string][3]float64
}

func (db *DB) nearest(target [3]float64) string {
    db.mutex.Lock()
    defer db.mutex.Unlock()
    // axtarış + delete — bütün kritik hissə
    var filename string
    smallest := 1000000.0
    for k, v := range db.store {
        dist := distance(target, v)
        if dist < smallest { filename, smallest = k, dist }
    }
    delete(db.store, filename)
    return filename
}
```

**Mənbə:** Chapter 9, page 250

---

## Fan-out: işi hissələrə böl, ayrı-ayrı goroutine-lərdə işlət

**Nə edir:** Ağır işi paralel emal etmək üçün hissələrə bölmək.

```go
func cut(original image.Image, db *DB, tileSize, x1, y1, x2, y2 int) <-chan image.Image {
    c := make(chan image.Image)
    go func() {
        newimage := image.NewNRGBA(image.Rect(x1, y1, x2, y2))
        // ... kvadrant emalı ...
        c <- newimage.SubImage(newimage.Rect)
    }()
    return c  // receive-only
}

// 4 kvadrant:
c1 := cut(original, &db, tileSize, x1, y1, x2/2, y2/2)
c2 := cut(original, &db, tileSize, x2/2, y1, x2, y2/2)
c3 := cut(original, &db, tileSize, x1, y2/2, x2/2, y2)
c4 := cut(original, &db, tileSize, x2/2, y2/2, x2, y2)
```

**Mənbə:** Chapter 9, page 250

---

## Fan-in: nəticələri birləşdir (select + WaitGroup)

**Nə edir:** Birdən çox channel-dən gələn nəticələri bir yerdə toplamaq.

```go
func combine(r image.Rectangle, c1, c2, c3, c4 <-chan image.Image) <-chan string {
    c := make(chan string)
    go func() {
        var wg sync.WaitGroup
        wg.Add(4)
        var s1, s2, s3, s4 image.Image
        var ok1, ok2, ok3, ok4 bool
        for {
            select {
            case s1, ok1 = <-c1: go copy(img, s1.Bounds(), s1, pt1)
            case s2, ok2 = <-c2: go copy(img, s2.Bounds(), s2, pt2)
            case s3, ok3 = <-c3: go copy(img, s3.Bounds(), s3, pt3)
            case s4, ok4 = <-c4: go copy(img, s4.Bounds(), s4, pt4)
            }
            if ok1 && ok2 && ok3 && ok4 { break }
        }
        wg.Wait()
        // encode + c <- string
    }()
    return c
}
```

**Mənbə:** Chapter 9, page 251

---

## `GOMAXPROCS=N` — paralel OS thread sayı

**Nə edir:** Go runtime-un istifadə edəcəyi maksimum CPU/thread sayını təyin edir. Go 1.5-dən default = CPU sayı.

```bash
GOMAXPROCS=1 ./mosaic_concurrent   # 1 CPU
./mosaic_concurrent                  # default = bütün CPU-lar
go test -cpu 2                       # test üçün
```

**Mənbə:** Chapter 9, page 254
