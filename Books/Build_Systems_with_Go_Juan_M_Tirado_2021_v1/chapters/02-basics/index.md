# Chapter 2 — The basics (səh. 21-51)

## Bu fəsil nədən bəhs edir?

Go dilinin əsasları: paketlər və importlar, dəyişənlər/konstantlar/enum-lar,
funksiyalar (variadic, closure), pointerlər, zero value anlayışı, loops və
branching (if/else, switch, for), error idarəetməsi, defer/panic/recover və
init funksiyaları.

## Əsas fikirlər

### 1. Paketlər və importlar
**Nədir:** Go kodu paketlərə təşkil olunur; paket = bir qrup source fayl.

**Necə işləyir:** `package main` xüsusi paketdir — kompilyatora bu paketin
icra olunan faylın giriş nöqtəsi olduğunu bildirir; `main` funksiyası
tələb olunur. Import zamanı Go əvvəlcə `GOROOT`-u, sonra `GOPATH`-ı yoxlayır.

**Üçüncü tərəf paketlərin yüklənməsi:**
```bash
>>> go get -v go.mongodb.org/mongo-driver
```
- `go get -v` → paketin source kodunu lokala endirir (Go yalnız source
  kodu kompilyasiya edir — binary yox). `-v` verbose (izahlı) rejim.
- Alternativ alias: `import myalias "go.mongodb.org/mongo-driver/mongo"`
- Müasir yanaşma Go modules-dur (Chapter 12-də izah olunur).

### 2. Dəyişənlər, konstantlar, enum-lar
**Go strong statically typed dildir** — tip kompilyasiya vaxtı təyin olunur.

**Dəyişən elan üsulları:**
```go
var a int          // sadə elan (zero value ilə)
a = 42
var aa int = 100   // elan + tip + dəyər
b := -42           // qısa forma; tip çıxarılır (int)
c := "this is a string"
var d, e string    // çoxsaylı elan
d, e = "var d", "var e"
f, g := true, false // çoxsaylı qısa elan
```

**Əsas tiplər (Table 2.1):**
| Tip | Təsvir |
|---|---|
| `bool` | Boolean |
| `string` | Simvol sətri |
| `int, int8..int64` | İşarəli tam ədədlər |
| `uint, uint8..uint64, uintptr` | İşarəsiz tam ədədlər |
| `byte` | Bayt (uint8 analoqu) |
| `rune` | Unicode code point (int32 alias; UTF-8 simvol) |
| `float32, float64` | Kəsr ədədlər |
| `complex64, complex128` | Kompleks ədədlər |

**Rune qeydi:** `rune` = int32 alias; UTF-8 simvollar 32 bit tələb edir.
`'€'` simvolunun dəyəri 8364, `%U` formatında `U+20AC`, `%c` ilə çap olunur.

**Konstantlar:**
```go
const (
    Pi = 3.14                    // tip çıxarılır → float64 (ən böyük tip seçilir)
    Avogadro float32 = 6.022e23  // tip açıq verilib
)
```
- Konstant kompilyasiya vaxtı təyin olunur və dəyişdirilə bilməz.
- Tip verilməzsə Go ən böyük uyğun tipi seçir (Pi → float64).

**Enum və iota:**
```go
type DayOfTheWeek uint8

const(
    Monday DayOfTheWeek = iota
    Tuesday
    Wednesday
    // ...
)
```
- `iota` → enum elementlərinə ardıcıl dəyərlər verir: Monday=0,
  Tuesday=1, Wednesday=2...
- Enum = konstant dəyərlərindən ibarət data tipi (həftə günləri, aylar).

### 3. Funksiyalar
**Nədir:** Kodun təkrar istifadə olunan parçasını saxlayan əsas vahid.

**Çoxdəyərli qaytarma (multiple return):**
```go
func ops(a int, b int) (int, int) {
    return a + b, a - b
}
sum, subs := ops(2, 2)
b, _ := ops(10, 2)   // _ ilə ikinci dəyər ignore edilir
```

**Variadic funksiyalar** (arqument sayı qeyri-müəyyən):
```go
func sum(nums ...int) int {
    total := 0
    for _, a := range(nums) {
        total = total + a
    }
    return total
}
total := sum(1,2,3,4,5)
```
- `...int` → istənilən sayda int arqument; funksiya daxilində massiv
  kimi davranır.

**Funksiya arqument kimi (first-class functions):**
```go
func doit(operator func(int,int) int, a int, b int) int {
    return operator(a,b)
}
c := doit(sum, 2, 3)
d := doit(multiply, 2, 3)
```

**Closure (kloucer):**
```go
func accumulator(increment int) func() int {
    i := 0
    return func() int {
        i = i + increment
        return i
    }
}
a := accumulator(1)
b := accumulator(2)
```
- Closure → bədənindən kənar dəyişənlərə istinad edən anonim funksiya.
- Hər `accumulator(...)` çağırışı öz `i` dəyişəninə bağlı yeni funksiya
  yaradır — a və b müstəqil sayğaclardır.

### 4. Pointerlər
**Nədir:** Dəyərin özü yox, yaddaş ünvanına istinad. `*T` — T tipinə pointer.

**Kitabdan kod nümunəsi:**
```go
func a(i int) { i = 0 }     // dəyər kimi — kopya dəyişir

func b(i *int) { *i = 0 }   // pointer kimi — orijinal dəyişir

func main() {
    x := 100
    a(x)        // x dəyişmir → 100
    fmt.Println(x)
    b(&x)       // x = 0 olur
    fmt.Println(x)
    fmt.Println(&x)  // yaddaş ünvanı çap olunur
}
```

**Sub-kod izahı:**
- `&x` → x-in ünvanını verir (`*int` tipi)
- `*i = 0` → pointerin göstərdiyi yaddaş yerinə 0 yazır
- Dəyər dəyişənin kopyasını götürür; pointer orijinala çıxış verir.
- Dəyişən kodun müxtəlif yerlərində dəyişməlidirsə pointer istifadə edin.

### 5. Nil və zero value
**Nədir:** Go-da hər tipin initialləşdirilməyəndə avtomatik default dəyəri var.

| Tip | Zero value |
|---|---|
| rəqəmsal (int, float) | `0` |
| bool | `false` |
| string | `""` (boş sətir) |
| pointer, funksiya | `nil` |

- `nil` undefined deyil — özü bir dəyərdir.

### 6. Loops və branching

**if/else** — mötərizə tələb olunmur:
```go
if x < 0.5 {
    fmt.Println("head")
} else {
    fmt.Println("tail")
}
```
- Go-da ternary operator (condition ? a : b) YOXDUR — sadəlik prinsipi.

**switch — dəyər əsaslı:**
```go
switch finger {
case 0: fmt.Println("Thumb")
case 1: fmt.Println("Index")
default: fmt.Println("...")
}
```

**switch — şərt əsaslı (dəyişənsiz):**
```go
switch {
case x < 0.25: fmt.Println("Q1")
case x < 0.5:  fmt.Println("Q2")
case x < 0.75: fmt.Println("Q3")
default:       fmt.Println("Q4")
}
```
- Bir neçə case eyni məntiqi paylaşırsa üst-üstə yığıla bilər.
- Hər switch-də `default` yazmaq yaxşı praktikadır.

**for — 4 forma:**
```go
for counter > 0 { ... }              // while tipli
for i:=0; i < x; i++ { ... }         // klassik C tipli
for { ... break/continue ... }       // şərtsiz + break/continue
for { ... }                           // sonsuz (server loop-ları üçün)
```
- Go-da `while` YOXDUR — bütün while məntiqi `for` ilə ifadə olunur.
- `break` → loop-u dayandırır; `continue` → növbəti iterasiyaya keçir.

### 7. Error idarəetməsi
**Nədir:** Go `try/catch` əvəzinə error dəyərləri ilə işləyir.

**Kitabdan kod nümunəsi:**
```go
func GetMusketeer(id int) (string, error) {
    if id < 0 || id >= len(Musketeers) {
        return "", errors.New(
            fmt.Sprintf("Invalid id [%d]", id))
    }
    return Musketeers[id], nil
}

mosq, err := GetMusketeer(id)
if err == nil {
    fmt.Printf("[%d] %s", id, mosq)
} else {
    fmt.Println(err)
}
```

**Sub-kod izahı:**
- `errors.New(...)` → mesajı olan error yaradır
- Konvensiya: `(dəyər, error)` qaytar; xəta halında dəyərə zero value
  (`""`), error-a mesaj ver
- `err == nil` → xəta yoxdur deməkdir
- Go-da try/catch/except idiomu QSDIĞIRDIR (FAQ-a əsasən kodu sadə saxlamaq üçün)

### 8. Defer, panic, recover
**defer:** Funksiyanı "sonradan icra" siyahısına əlavə edir; əhatə edən
funksiya bitəndə icra olunur — resurs təmizliyi üçün nəzərdə tutulub.
Deferred funksiyalar **ters sıra ilə** icra olunur (LIFO).

```go
func main() {
    defer CloseMsg()
    fmt.Println("Doing something…")
    defer fmt.Println("Certainly closed!!!")
    fmt.Println("Doing something else…")
}
// Çıxış: Doing something… / Doing something else… / Certainly closed!!! / Closed!!!
```

**panic:** İcrai axını dayandırır, deferred funksiyaları icra edir və
proqram çökənə qədər nəzarəti yuxarı çağıran tərəfə qaytarır. Proqramın
nəzarətindən kənar vəziyyətlər üçün:
```go
if i > 2 {
    panic("Panic was called")
}
```

**recover:** Yalnız deferred funksiya daxilində işləyir; panic-i tutub
normal icraiı bərpa edir:
```go
defer func() {
    r := recover()
    if r != nil {
        fmt.Println("No need to panic if i=", r)
    }
}()
panic(i)  // recover bu dəyəri tutur
```

### 9. Init funksiyaları və runtime başlanğıc sırası
**Nədir:** Paket başına bir dəfə icra olunan konfiqurasiya funksiyaları.
Kitabxanaların `main`-siz başlanğıc tələbləri üçün.

**Go runtime başlanğıc sırası:**
1. Import olunmuş paketlər rekursiv initialləşdir
2. Dəyişənlərə dəyər mənimsət (var x = xSetter())
3. `init()` funksiyaları icra olunur

```go
var x = xSetter()  // 1-ci: xSetter çap olunur

func init() {      // 2-ci: Init function
    fmt.Println("Init function")
}

func main() {       // 3-cü: This is the main
    fmt.Println("This is the main")
}
```

- init funksiyası arqumentsiz və qaytarmasızdır; paketdə bir neçə init ola bilər
- İstifadə olunmayan paketi import etmək QADAĞANDIR; yalnız init
  lazımdırsa blank import istifadə olunur: `_ "a"` → yalnız side effect-lər
  (init-lər) işə düşür, paket kodda istifadə edilmir.

## Əsas terminlər
- Package (paket) — kodun təşkilat vahidi
- GOPATH/GOROOT — Go workspace və quraşdırma yolları
- Static typing (statik tiplilik) — tip kompilyasiya vaxtı təyin olunur
- iota — enum konstantlarına ardıcıl dəyər verən identifikator
- Variadic function (variadik funksiya) — dəyişən sayda arqument qəbul edən
- Closure (kloucer) — xarici dəyişənləri "yadda saxlayan" funksiya
- Zero value (sıfır dəyər) — initialləşdirilməmiş dəyişənin default dəyəri
- defer/panic/recover — təmizlik/qəza/bərpa mexanizmi
- Blank import (boş import) — `_ "pkg"` yalnız init-i işə salır

## Praktik nəticə
Go-nun əsas quruluşu: main paket + main() giriş; xətalar try/catch əvəzinə
`(dəyər, error)` konvensiyası ilə; resurs təmizliyi defer ilə; while yoxdur —
for hamısını əvəz edir; ternary yoxdur — if/else yazılır. Init funksiyaları
paket konfiqurasiyası üçün, `_ "pkg"` blank importu isə yalnız side-effect
 üçün istifadə olunur.

## Mənbə
Pages: 21-51 (PDF 21-51)
