# Chapter 1 — Variables and Operators (Dəyişənlər və Operatorlar)

## Bu fəsil nədən bəhs edir?

Go-ya giriş (filosofiya, funksionallıq), ilk Go proqramının sətir-sətir təhlili, dəyişən
bəyan etmənin bütün yolları (var, :=, çoxsaylı), tip çıxarılması və onun səhvləri,
operatorlar (arifmetik, müqayisə, məntiq, qısa formalar), zero value-lar, value vs
pointer (stack/heap, escape analysis), sabitlər, iota ilə enum-lar və scope qaydaları
(shadowing daxil).

## Əsas fikirlər

### 1. Go Nədir və Niyə Populyardır
**Nədir:** Google-dan yaranan, statik tipli, GC-li, channel əsaslı paralelliyi olan dil.
**Üstünlükləri:**
- Dinamik dil hissi (JS/PHP) + güclü tip dillərin performansı (C++/Java)
- Tək, müstəqil binary — deploy asan; cross-compile (Windows/Linux/macOS/Android)
- İldırım sürətli compiler — layihə böyüsə belə
- Statik tip + təhlükəsiz yaddaş modeli + GC → çox bug və təhlükəsizlik deşiyindən qoruyur
- Multicore CPU-lar üçün dizayn — paralel/konkurent kod SADƏ və TƏHLÜKƏSİZ

### 2. İlk Go Proqramı — Struktur
**Kitabdan kod nümunəsi:**
```go
package main                        // paket bəyannaməsi — hər faylda MÜTLƏQ
                                    // main = bilavasitə icra olunur
import (
    "errors"
    "fmt"
    "log"
    "math/rand"
    "strconv"
    "time"
)                                   // importlar FAYL-a aiddir (paketyox)
var helloList = []string{           // qlobal slice — UTF-8 dəstəyi
    "Hello, world",
    "Καλημέρα κόσμε",
}

func main() {                       // ENTRY POINT — Go avtomatik çağırır
    rand.Seed(time.Now().UnixNano())
    index := rand.Intn(len(helloList))   // [0, len) random
    msg, err := hello(index)
    if err != nil {
        log.Fatal(err)              // xəta + proqramı ÖLDÜRÜR
    }
    fmt.Println(msg)
}

func hello(index int) (string, error) {
    if index < 0 || index > len(helloList)-1 {
        return "", errors.New("out of range: " + strconv.Itoa(index))
    }
    return helloList[index], nil
}
```
**Sub-kod izahı:**
- Sonuncu qaytarma dəyəri error — Go konvensiyası
- `len(list)-1` = sonuncu indeks; `list[len(list)]` = runtime panic
- Standart kitabxana yüksək keyfiyyətlidir — maksimal istifadə et; xarici paket URL-ə
  bənzər görünür (github.com/fatih/color)

### 3. Dəyişən Bəyannaməsi — Bütün Yollar
**Tam var (bünövrə):**
```go
var foo string = "bar"    // var + ad + tip + ilkin dəyər
```
**Çoxsaylı var bloku (paket səviyyəsində common):**
```go
var (
    Debug      bool   = false
    LogLevel   string = "info"
    startUpTime time.Time = time.Now()
)
```
**Tip və ya dəyəri atla:**
```go
var Debug bool        // dəyər YOX → ZERO VALUE
var LogLevel = "info" // tip YOX → inference
```
**Short declaration `:=` (funksiya daxilində YALNIZ):**
```go
Debug := false
LogLevel := "info"
Debug, LogLevel, startUpTime := false, "info", time.Now()  // çoxsaylı
// Funksiyadan çox dəyər (ƏN COMMON istifadə):
Debug, LogLevel, startUpTime := getConfig()
```
**Tip inference xətasi:**
```go
var seed = 1234456789     // int CHIxarılır
rand.Seed(seed)           // XƏTA: rand.Seed int64 istəyir!
var seed int64 = 1234456789   // HƏLL
```

**UTF-8 adlar:** `デバッグ := false`, `日志级别 := "info"` — birinci simvol hərf/_,
qalanları hərf/rəqəm/_. Asiyada populyarlıq səbəblərindən biri.

### 4. Dəyər Dəyişmə
```go
offset := 5
offset = 10                       // tək
query, limit, offset = "ball", offset, 20   // çoxsaylı — paralel assign
// DİQQƏT: oxunaqlılıq üçün zərurətsə çoxsətirdən istifadə et
```

### 5. Operatorlar
**Qruplar:** arifmetik, müqayisə, məntiq, ünvan (&/`*`), qəbul (`<-` — kanallar).

**Restoran hesabı simulyasiyası (kitabdan):**
```go
var total float64 = 2 * 13             // 2 məhsul × 13$
total = total + (4 * 2.25)             // içkilər əlavə
total = total - 5                      // endirim
tip := total * 0.1                     // 10% bəxşiş
total = total + tip
split := total / 2                     // 2 nəfərə böl

visitCount := 24
visitCount = visitCount + 1
remainder := visitCount % 5            // hər 5-ci ziyarət
if remainder == 0 {
    fmt.Println("With this visit, you've earned a reward.")
}
```
**String birləşdirmə:** `fullName := givenName + " " + familyName`

**Qısa operatorlar:** `+=`, `-=`, `++`, `--` (loop-larda); string üçün `+=`.

**Müqayisə:** `==`, `!=`, `<`, `<=`, `>`, `>=`
**Məntiq:** `&&` (və), `||` (və ya), `!` (deyil) — YALNIZ bool ilə.

**Membership nümunəsi:**
```go
fmt.Println("Silver member:", visits >= 10 && visits < 21)
fmt.Println("Gold member  :", visits > 20 && visits <= 30)
fmt.Println("Platinum     :", visits > 30)
```

### 6. Zero Value-lar
İlkin dəyər verilməyən dəyişən tipin DEFAULT dəyərini alır:
| Tip | Zero value |
|---|---|
| int/float | 0 |
| bool | false |
| string | "" |
| slice/map | nil |
| struct | bütün sahələr zero |
| pointer | nil |

**fmt.Printf və formatlar:** `%#v` (tip+dəyər), `%v`, `%T`, `%s`, `%d` — `\n` əl ilə.

### 7. Value vs Pointer — Əsas Fərq
**Value (kopya):** funksiyaya kopya ötürülür → dəyişiklik xarici dəyəri TƏSİR ETMİR;
daha az bug; stack istifadəsi (sadə yaddaş idarəetməsi); amma çox kopya = artıq yaddaş.

**Pointer:** dəyərin ÜNVANI; kopya YOX; dəyişiklik XARİCİDƏN GÖRÜNÜR; heap-da yaşayır
(GC idarəsində); nil ola bilər.

**Stack vs heap:**
- Stack — sadə scope logic ilə avtomatik təmizlənir
- Heap — pointer-lər üçün; GC (garbage collection) fondda periodic toplayır
- Escape analysis hansı dəyərin heap-a düşdüyünü müəyyən edir — birbaşa nəzarət YOXDUR;
  implementasiya detalıdır, dəyişə bilər
- CPU tradeoff: pointer kopya xərcini azaldır, amma GC yükü artır — **premature
  optimization etmə; ölç, sonra dəyiş**

**Pointer yaratmaq — 3 yol:**
```go
var count1 *int           // 1) var — nil pointer
count2 := new(int)        // 2) new() — zero value-lu yaddaş + pointer
countTemp := 5
count3 := &countTemp      // 3) & — mövcud dəyişəndən
t := &time.Time{}         // struct-literal-dan birbaşa
```

**Dereference:** `*count` ilə dəyər alınır; **nil pointer-i dereference = runtime
panic** — compiler xəbərdarlıq ETMİR; həmişə `if ptr != nil` yoxla. Struct metodu
çağırırkən dereference LAZIM DEYİL (`t.String()`).

**Funksiya dizaynı (kitabdan):**
```go
func add5Value(count int) {      // kopya dəyişir
    count += 5
}
func add5Point(count *int) {     // OrijİNAL dəyişir
    *count += 5
}
var count int
add5Value(count)     // count hələ 0
add5Point(&count)    // count = 5 — funksiya İÇİNDƏ dəyişdi
```
**Swap (activity):**
```go
func swap(a *int, b *int) {
    *a, *b = *b, *a
}
```
**Tövsiyə:** beginner kimi pointer-lardan QAÇ — performans problemi və ya dizayn
təmizliyi tələb edənə qədər.

### 8. Constants (Sabitlər)
```go
const GlobalLimit = 100
const MaxCacheSize int = 10 * GlobalLimit   // ifadə ilə də ola bilər
const (
    CacheKeyBook = "book_"
    CacheKeyCD   = "cd_"
)
```
Dəyişməz; runtime-da dəyişmir; gələcək dəyişikliklər üçün bir nöqtədən idarə —
hardcode-dan üstün.

**Cache nümunəsi (kitabdan):**
```go
var cache map[string]string     // PAYLAŞILMIŞ cache

func cacheSet(key, val string) {
    if len(cache)+1 >= MaxCacheSize {   // limit nəzarəti
        return
    }
    cache[key] = val
}
func GetBook(isbn string) string {
    return cacheGet(CacheKeyBook + isbn)  // prefiks ilə unikal açar
}
```

### 9. Enum və iota
Go-nun built-in enum-u YOXDUR — const + iota ilə özünü yaradır:
```go
const (
    Sunday    = iota   // 0 — birinci sətirdə iota = 0
    Monday            // 1 — avtomatik artır
    Tuesday           // 2
    Wednesday
    Thursday
    Friday
    Saturday         // 6
)
```
Ortadan dəyər əlavə etmək asandır — nömrələmə özü düzəlir.

### 10. Scope (Əhatə Dairəsi)
**Quruluş:** package scope (top) → funksiya → blok `{...}` → uşaq bloklar. Axtarış:
caridən başlayır → valideynlərə qalxır → tapanda dayanır → tapmasa error.

**Shadowing (kölgələmə):**
```go
var level = "pkg"        // package scope
func main() {
    fmt.Println(level)   // "pkg"
    level := 42           // SHADOW — yeni, fərqli tip də ola bilər!
    fmt.Println(level)   // 42
}
func funcA() {
    fmt.Println(level)   // "pkg" — funcA main-də HARADA çağırıldığından asılı DEYİL
}
```
**Static scope resolution:** funksiya öz təyin olunduğu yerdəki scope-ları görür —
çağırılış YERİ əhəmiyyətsizdir.

**Uşaq scope-dan çıxış:** daxili blokda təyin edilən dəyişən xaricdən görünməz.

**Message Bug nümunəsi:** if/else hər ikisi öz bloklarında `message` təyin edir →
xaricdən `undefined`. Həll: if-dən ƏVVƏL təyin et, bloklarda SADECE təyin et.

**Bad Count Bug:** `if count < 5 { count := 10; count++ }` — daxili count SHADOW-dur;
xarici count dəyişmir. `:=` yerinə `=` yaz.

## Əsas terminlər
- Package Declaration (`package main`) — icra olunacaq proqram üçün mütləq
- Short Variable Declaration — `:=` — tip çıxarışı ilə qısa bəyan (funksiyada)
- Type Inference — literal-dan tip çıxarma qaydaları
- Zero Value — tipin default dəyəri (0/false/""/nil)
- Pointer — dəyərin yaddaş ünvanı; `*T` tipi
- Dereference — `*ptr` ilə pointer-dən dəyər almaq
- `new(T)` — zero-lu yaddaş + pointer qaytaran built-in
- Stack/Heap — avtomatik/GC-idarəli yaddaş sahələri
- Escape Analysis — dəyərin heap-a "qaçmasının" təhlili
- Shadowing — daxili scope-da eyni adlı YENİ dəyişən (xaricini gizlədir)
- Static Scope Resolution — funksiyanın TƏYİN yerinə görə scope
- Constant — dəyişməz dəyər; runtime-da sabit
- iota — const siyahısında avtomatik artan sayğac (enum aləti)
- Garbage Collection — heap-dəki istifadəsiz dəyərlərin fondda toplanması

## Praktik nəticə

(1) `:=` real kodun 90%-i; qalan 10% tam var — tipə nəzarət (məs. int64) lazım olanda.
(2) Funksiyadan çox dəyər + `:=` bir sətirdə — Go-nun ən məhsuldar kombinasiyası.
(3) Tip çıxarışı literal-in DEFAULT tipinə görədir — API int64 istəyirsə açıq yaz.
(4) Pointer: funksiya XARİCİNƏ dəyişiklik + performans + nil-state; amma error-meillidir —
ehtiyac yaranana qədər value istifadə et. (5) `*ptr`-dən əvvəl nil yoxla — compiler
səni xilas ETMİR. (6) Zero value-ları bil — ilkin dəyərsiz bəyan təhlükəsizdir. (7)
Shadowing: eyni adda daxili dəyişən xaricini GİZLƏDİR — unutma, debug üçün increment
bəyanları bu yolla gedir. (8) `:=` YALNIZ funksiyada; paket səviyyəsində var. (9) Const
+ prefiks pattern — paylaşımlı cache-lərdə unikal açar üçün. (10) iota enum-ın Go
həlli — manual nömrələmədən üstün.

## Mənbə
Pages: 1-53 (PDF 34-87)
