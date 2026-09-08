# Chapter 1 — Variables and Operators (səh. 32-81)

## Bu fəsil nədən bəhs edir?

Go-ya giriş: ilk proqram (package, import, func main), dəyişən elanının
bütün formaları (var bloku, tip/value atlama, `:=` qısa elan), dəyər
dəyişməsi, operatorlar (arifmetik, müqayisə, məntiqi, shorthand), zero
value-lar, pointer-lər (var/new/&, dereference, funksiya dizaynı),
konstantlar, iota ilə enum-lar və scope qaydaları (shadowing).

## Əsas fikirlər

### 1. İlk Go proqramı
**Kitabdan kod nümunəsi (random salam):**
```go
package main

import (
    "errors"
    "fmt"
    "log"
    "math/rand"
    "strconv"
    "time"
)

var helloList = []string{
    "Hello, world",
    "Καλημέρα κόσμε",
    "こんにちは世界",
    "Привет, мир",
}

func main() {
    // Seed random number generator using the current time
    rand.NewSource(time.Now().UnixNano())
    // Generate a random number in the range of out list
    index := rand.Intn(len(helloList))
    // Call a function and receive multiple return values
    msg, err := hello(index)
    // Handle any errors
    if err != nil {
        log.Fatal(err)
    }
    // Print our message to the console
    fmt.Println(msg)
}

func hello(index int) (string, error) {
    if index < 0 || index > len(helloList)-1 {
        // Create an error, convert the int type to a string
        return "", errors.New("out of range: " + strconv.Itoa(index))
    }
    return helloList[index], nil
}
```

**Sub-kod izahı:**
- `package main` → icra olunan proqram; kitabxana paketləri istənilən
  ad ala bilər; eyni qovluqdakı fayllar eyni paketdə
- `import (...)` → standart kitabxana yüksək keyfiyyətlidir; xarici
  paketlər URL kimi görünür (github.com/...); import fayl-lokaldır
- `func main()` → yeganə giriş nöqtəsi (main paketində bir dənə)
- `msg, err := hello(index)` → çoxdəyərli qaytarma; err-sonuncu
  konvensiyası
- `if err != nil` → error yoxlaması (nil = dəyərsizlik)
- `rand.Intn(len(helloList))` → 0..len-1; slice-da son index len-1

**Go haqqında:** statik tipli + GC; sürətli kompilyator; tək fayl
executable; multi-core CPU-lar üçün konkurensiya (kanallar).

### 2. Dəyişən elan formaları
```go
// Tam var:
var foo string = "bar"

// var bloku (paket səviyyəsində çox işlədilir):
var (
    Debug       bool   = false
    LogLevel    string = "info"
    startUpTime time.Time = time.Now()
)

// Tip və ya dəyəri atla:
var (
    Debug       bool              // zero value: false
    LogLevel         = "info"     // tip çıxarılır
    startUpTime      = time.Now()
)

// Tip çıxarışı xəta verəndə tam elan lazımdır:
var seed int64 = 1234456789     // rand.NewSource int64 tələb edir!

// Qısa elan (YALNIZ funksiya daxilində):
Debug := false
LogLevel := "info"

// Çoxsaylı qısa elan (bir sətir, sayı bərabər):
Debug, LogLevel, startUpTime := false, "info", time.Now()

// Funksiyadan çoxdəyər:
Debug, LogLevel, startUpTime := getConfig()
```
- `:=` real Go kodunda ən çox görülən üsuldur — kompakt, dinamik hiss
- Eyni sətirdə var: `var start, middle, end float32` (eyni tip);
  dəyərlə müxtəlif tiplər ola bilər

### 3. Non-English adlar (UTF-8)
```go
デバッグ := false
日志级别 := "info"
_A1_Μείγμα := ""
```
- Go UTF-8 uyğundur; ilk simvol hərf və ya `_` — Asiyada populyarlıq
  səbəblərindən biri

### 4. Dəyər dəyişməsi
```go
offset := 5
offset = 10                   // tək
query, limit, offset = "ball", offset, 20   // çoxsaylı; PARALEL təyinat
```
- Çoxsaylı təyinatda sağ tərəf ƏVVƏL hesablanır (swap üçün təhlükəsiz)

### 5. Operatorlar
**Qruplar:** arifmetik (+ - * / %), müqayisə (== != < <= > >=),
məntiqi (&& || !), address (& *), receive (<-).

**Restoran hesabı nümunəsi:**
```go
var total float64 = 2 * 13              // əsas
total = total + (4 * 2.25)              // içkilər
total = total - 5                       // endirim
tip := total * 0.1                      // 10% bəxşiş
total = total + tip
split := total / 2                      // bölüşdürmə

// Hər 5-ci ziyarət mükafatı:
visitCount := 24
visitCount = visitCount + 1
remainder := visitCount % 5
if remainder == 0 {
    fmt.Println("With this visit, you've earned a reward.")
}
```

**String birləşmə:** `fullName := givenName + " " + familyName`

**Shorthand operatorlar:**
```go
count := 5
count += 5    // artır və mənimsət
count++       // +1
count--       // -1
count -= 5
name += " Smith"   // string üçün də işləyir
```

**Müqayisə + məntiq (membership nümunəsi):**
```go
visits := 15
fmt.Println("First visit    :", visits == 1)
fmt.Println("Return visit   :", visits != 1)
fmt.Println("Silver member  :", visits >= 10 && visits < 21)
fmt.Println("Gold member    :", visits > 20 && visits <= 30)
fmt.Println("Platinum member:", visits > 30)
```

### 6. Zero value-lar
Başlanğıc dəyəri verilməyən dəyişənlər tipin default-unu alır:

| Tip | Zero value |
|---|---|
| int | 0 |
| float64 | 0.0 |
| bool | false |
| string | "" (boş) |
| []string | nil |
| time.Time | 0001-01-01... |

```go
var count int
fmt.Printf("Count  : %#v \n", count)   // 0
var debug bool                          // false
var message string                      // ""
var emails []string                     // []string(nil)
```
- `%#v` → dəyər + tip; `\n`-i özün əlavə et

### 7. Value vs Pointer
**Nədir:** Dəyər keçirildikdə Go KOPYA yaradır; pointer = ünvan —
kopya yoxdur, orijinal dəyişir.

**Stack vs Heap:** dəyərlər stack-də (scope bitəndə azad); pointer-li
dəyərlər heap-də (GC tərəfindən idarə) — "escape analysis" qərar verir;
düz nəzarət yoxdur, yalnız ümumi göstərişlər.

**Pointer yaratmağın 3 yolu:**
```go
var count1 *int          // 1) var — nil dəyər
count2 := new(int)       // 2) new — zero value işarətlənir
countTemp := 5
count3 := &countTemp     // 3) & — mövcud dəyişəndən
t := &time.Time{}        // struct-literal-dan birbaşa
```

**Dereference + nil yoxlaması:**
```go
if count1 != nil {
    fmt.Printf("count1: %#v\n", *count1)   // * = dəyəri oxu
}
if t != nil {
    fmt.Printf("time : %#v\n", t.String())  // metod çağırışında * lazım DEYİL
}
```
- Nil pointer-in dereference-i → runtime panic — yoxlama şərtdir
- Pointer-lərin == müqayisəsi YALNIZ eyni ünvan üçün true

**Funksiya dizaynı:**
```go
func add5Value(count int) {     // kopya — xarici dəyişməz
    count += 5
}

func add5Point(count *int) {    // pointer — orijinal dəyişir
    *count += 5
}

var count int
add5Value(count)     // count hələ 0
add5Point(&count)    // count artıq 5
```

**Swap activity (pointer praktikası):**
```go
func swap(a *int, b *int) {
    *a, *b = *b, *a    // paralel təyinat + dereference
}
a, b := 5, 10
swap(&a, &b)           // → 10, 5
```

### 8. Konstantlar
**Nədir:** dəyişməz dəyərlər — hardcode-un idarəli alternativi.

```go
const GlobalLimit = 100
const MaxCacheSize int = 10 * GlobalLimit   // konstant ifadələri bir-birini istifadə edə bilər

const (
    CacheKeyBook = "book_"
    CacheKeyCD   = "cd_"
)

// Paylaşılan keş nümunəsi:
var cache map[string]string

func cacheGet(key string) string { return cache[key] }

func cacheSet(key, val string) {
    if len(cache)+1 >= MaxCacheSize {
        return
    }
    cache[key] = val
}

func GetBook(isbn string) string { return cacheGet(CacheKeyBook + isbn) }
func SetBook(isbn string, name string) { cacheSet(CacheKeyBook+isbn, name) }
func GetCD(sku string) string { return cacheGet(CacheKeyCD + sku) }
func SetCD(sku string, title string) { cacheSet(CacheKeyCD+sku, title) }

func main() {
    cache = make(map[string]string)
    SetBook("1234-5678", "Get Ready To Go")
    SetCD("1234-5678", "Get Ready To Go Audio Book")
    fmt.Println("Book :", GetBook("1234-5678"))
    fmt.Println("CD   :", GetCD("1234-5678"))
}
```

### 9. Enums (iota)
```go
// Əl ilə:
const (
    Sunday    = 0
    Monday    = 1
    // ... Saturday = 6
)

// iota ilə — avtomatik artım:
const (
    Sunday = iota     // 0
    Monday            // 1
    Tuesday           // 2
    Wednesday         // 3
    Thursday          // 4
    Friday            // 5
    Saturday          // 6
)
```
- Go-də built-in enum YOXDUR — iota + konstantlarla qurulur
- `_` ilə dəyər ötürülür; ofset və hesablamalar mümkündür

### 10. Scope və shadowing
**Qaydalar:** hər `{` yeni uşaq scope; dəyişən axtarışı → cari scope →
 valideynlər → paket scope; tapılmasa compile xətası.

```go
// Eyni dəyişən hər yerdə:
var level = "pkg"
func main() {
    fmt.Println("Main start :", level)   // pkg
    if true {
        fmt.Println("Block start :", level)   // pkg
    }
}

// SHADOWING — uşaq scope-də eyni adda yeni dəyişən:
var level = "pkg"
func main() {
    fmt.Println("Main start :", level)   // pkg
    level := 42                           // yeni level (shadow)
    if true {
        fmt.Println("Block start :", level)   // 42
        funcA()
    }
    fmt.Println("Main end :", level)      // 42
}
func funcA() {
    fmt.Println("funcA start :", level)   // pkg — static scope!
}
```
- `funcA` çağırıldığı YERİ yox, TƏYİN olunduğu yeri görür (static scope)
- Uşaq scope-dəki dəyişən valideyndən görünmür:

**Bug-lar (scope dərsləri):**
```go
// Activity 1.03 — if/else daxilindəki message xaricdən görünmür:
if count > 5 {
    message := "Greater than 5"     // YALNIZ if daxilində
} else {
    message := "Not greater than 5"
}
fmt.Println(message)                // compile xətası!
// Fix: message-i if-dən ƏVVƏL elan et

// Activity 1.04 — shadowing:
count := 0
if count < 5 {
    count := 10    // YENİ count (shadow!) — xarici count dəyişmir
    count++
}
fmt.Println(count == 11)   // false — xarici hələ 0-dır
// Fix: := yerinə = istifadə et
```

## Əsas terminlər
- Package (paket) — kod təşkilatı; main icra üçün
- Import — fayl-lokal asılılıq elanı
- Zero value — tipsiz başlanğıcın default dəyəri
- Short variable declaration (:=) — funksiya-daxili qısa elan
- Type inference — dəyərdən tip çıxarması
- Pointer / dereference (*, &) — ünvan / dəyərə çıxış
- new — tip üçün yaddaş + pointer
- Stack vs Heap / escape analysis — yaddaş yerləşməsi
- GC (garbage collector) — heap təmizləyicisi
- Shadowing — uşaq scope-də adın üstünü örtməsi
- Static scope resolution — təyinetmə yerinə görə axtarış
- const / iota — dəyişməz dəyər / enum sayğacı
- Paralel təyinat (a, b = b, a) — eyni vaxtlı dəyişmə

## Praktik nəticə
Dəyişənlərin 90%-ı üçün `:=` (funksiya daxilində); paket səviyyəsində
var bloku; tip qeyri-müəyyənliyində tam var elan. Hər zero value-i bil
— nil slice/map/string ilə işləməzdən əvvəl yoxla. Pointer-lər:
kopyadan yığın yaddaş + funksiyada orijinalı dəyişmək; amma ehtiyatlı —
nil dereference və shadowing klassik bug mənbəyidir. Konstantlar
hardcode-u əvəz edir; əlaqəli sabitlər üçün iota. Scope-a diqqət:
hər `:=` yeni dəyişən yaradır — köhnəsini dəyişmək istəyirsənsə `=` işlət.

## Mənbə
Pages: 32-81 (PDF 32-81)
