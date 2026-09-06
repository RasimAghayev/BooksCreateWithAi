# Chapter 2 — Diving Into Go

## Bu chapter nədən bəhs edir?

Bu chapter, praktik tətbiq hazırlayaraq Go-nun əsaslarına giriş edir: ilk Go proqramını yazmaq, Go ecosystem-ünü (compiler, tools, modules) tanımaq, fayl oxumaq, error handling (xəta idarəetməsi) və buferli oxuma (buffered I/O) konsepsiyalarını əhatə edir.

## Əsas fikirlər

### 1. Go Ecosystem və Toolchain
**Nədir:** Go yalnız proqramlaşdırma dili deyil, tam ecosystem — kompilyator, package manager (paket menecer), formatter (formatlaşdırıcı) və dokumentasiya alətləri daxildir.

**Necə işləyir:** `go` komandası Swiss Army Knife (şveys çəkici) kimi işləyir: `go build`, `go run`, `go mod`, `go doc` və s. alt əmrləri ilə hər şeyi idarə edir.

**Nəyə lazımdır:** Developer təcrübəsini yüksəltmək, kod formatlaşdırmasını standartlaşdırmaq, asılılıqları idarəetmək.

**Üstünlükləri:**
- `go fmt` ilə kod avtomatik formatlanır — debate (mübahisə) olunmur
- `go mod` ilə asılılıqlar versiyalanır və təkrarlanır
- `go doc` ilə yerli dokumentasiya oxunur

**Çatışmamazlıqları:**
- Başlangıçda öyrəniləsi çox alət var

**Kitabdan kod nümunəsi:**
```go
$ go mod init gobook/wordcount
go: creating new go.mod: module gobook/wordcount
```
**Sub-kod izahı:**
- `go mod init` → Yeni Go module (modul) yaradır
- `gobook/wordcount` → Module adı

**Mənbə:** Chapter 2, pages 25-66

### 2. go build və go run
**Nədir:** `go build` — kodu kompilyə edib binary (icra edilən fayl) yaradır; `go run` — kompilyə edib birbaşa icra edir, binary saxlamır.

**Necə işləyir:** `go build` package-ları və asılılıqları compile edərək executable (icra edilə bilən) fayl yaradır. `go run` eyni prosesi edir, lakin nəticəni diskdə saxlamadan icra edir.

**Nəyə lazımdır:** Tətbiqi test etmək və paylaşmaq.

**Üstünlükləri:**
- `go run` — rapid testing (sürətli test) üçün idealdır
- `go build` — distributable (paylaşıla bilən) binary yaradır

**Çatışmamazlıqları:**
- `go run` böyük tətbiqlər üçün uyğun deyil

**Kitabdan kod nümunəsi:**
```go
$ go build
$ ./wordcount words.txt
Found 4 words
```
**Mənbə:** Chapter 2, pages 25-66

### 3. gofmt — Standart Formatlaşdırma
**Nədir:** Bütün Go kodunu eyni stilə gətirən avtomatik formatlaşdırıcı alət.

**Necə işləyir:** `gofmt` kodun strukturunu (indent, boşluq, brackets) standart Go style-ə uyğunlaşdırır. `gofmt -d` diff çıxarır, `gofmt -w` faylı dəyişdirir.

**Nəyə lazımdır:** Code review (kod nəzəri) vaxtını azaltmaq, ekibin oxunaqlılığını artırmaq.

**Üstünlükləri:**
- Bütün Go kodunu eyni formata salır
- Code review-da stil mübahisələri aradan götürülür

**Çatışmamazlıqları:**
- Fərdi format seçimi qeyri-mümkün

**Kitabdan kod nümunəsi:**
```bash
$ gofmt -d main.go
```
**Mənbə:** Chapter 2, pages 25-66

### 4. Package (Paket) və Import (Daxil Etmə)
**Nədir:** Package — Go kodunun təşkilat vahidi; import — başqa package-ləri daxil etmək.

**Necə işləyir:** `package main` executable program göstərir. `import "fmt"` kimi başqa paketləri daxil edirik. Export (ixrac) olunan identifier-lər böyük hərflə başlayır.

**Nəyə lazımdır:** Kodun təkrar istifadə edilməsi, modulluq, encapsulation (kapsulyasiya).

**Üstünlükləri:**
- Açıq export qaydası — `Println` export edilir, `println` deyil
- Package scope (paket görüş sahəsi) daxili idarəetmə

**Çatışmamazlıqları:**
- Import cycle (döngə) qarşısını almaq lazımdır

**Kitabdan kod nümunəsi:**
```go
package main

import (
    "fmt"
    "strings"
)
```
**Sub-kod izahı:**
- `package main` → Executable program
- `import` → Asılılıqları daxil et
- `"fmt"`, `"strings"` → Standart kitabxanalar

**Mənbə:** Chapter 2, pages 25-66

### 5. Variable Declaration (Dəyişən Elanı)
**Nədir:** Dəyişənlərin yaradılması — `var` açar sözü ilə ya da qısaca `:=` operatoru ilə.

**Necə işləyir:** `var numSpaces int` — tip və sıfır dəyəri ilə yaradır. `text := "hello"` — tip inference (tip çıxarma) ilə yaradır. Short declaration `:=` yalnız funksiya daxilində işləyir.

**Nəyə lazımdır:** Məlumatları saxlamaq və idarəetmək.

**Üstünlükləri:**
- Short declaration daha az kod tələb edir
- Type inference tipi avtomatik müəyyən edir

**Çatışmamazlıqları:**
- Short declaration shadowing (kölgələmə) riski — eyni adlı dəyişən scope (görüş sahəsi) daxilində yenidən yaradıla bilər

**Kitabdan kod nümunəsi:**
```go
var numSpaces int          // Zero value: 0
text := "let's count some words!"  // Type inference: string
```
**Mənbə:** Chapter 2, pages 25-66

### 6. For Loops (Döngülər)
**Nədir:** Go-da yeganə döngü konstruksiyası — `for`. Üç Variety (növ) var: klassik 3-hissəli, while kimi, infinite (sonsuz).

**Necə işləyir:** `for init; condition; post { }` — init bir dəfə işləyir, condition hər dəfə yoxlanılır, post hər iterasiyadan sonra işləyir. `for condition { }` while kimi. `for { }` sonsuz döngü.

**Nəyə lazımdır:** Təkrarlanan əməliyyatlar, kolleksiyalar üzərində iterasiya.

**Üstünlükləri:**
- Tək döngü konstruksiyası — öyrənmə asanlığı
- `range` ilə kolleksiyaları asan iterasiya etmək

**Çatışmamazlıqları:**
- While/do-while yoxdur — adaptasiya lazım ola bilər

**Kitabdan kod nümunəsi:**
```go
for i := 0; i < len(text); i++ {
    if text[i] == ' ' {
        numSpaces++
    }
}
```
**Mənbə:** Chapter 2, pages 25-66

### 7. Error Handling (Xəta İdarəetməsi)
**Nədir:** Go-da xətalar `error` tipi ilə təmsil olunur. `nil` xəta yoxdur, non-nil xəta var deməkdir.

**Necə işləyir:** Funksiya `(result, error)` qaytarır. Caller (çağıran) `if err != nil` ilə yoxlayır. `log.Println` ilə log (qeyd) edib `os.Exit(1)` ilə dayandırır.

**Nəyə lazımdır:** Runtime xətalarını idarəetmək, proqramın düzgün şəkildə bağlanmasını təmin etmək.

**Üstünlükləri:**
- Explicit (açıq) error handling — implicit (daxili) exception yox
- Xətalar kodda görünür, debug asan

**Çatışmamazlıqları:**
- verbose (çox sözlü) ola bilər — `if err != nil` təkrarlanır

**Kitabdan kod nümunəsi:**
```go
fileContents, err := os.ReadFile(filename)
if err != nil {
    log.Println(err)
    os.Exit(1)
}
```
**Sub-kod izahı:**
- `os.ReadFile` → Faylı oxuyur, `([]byte, error)` qaytarır
- `if err != nil` → Xəta yoxlanışı
- `log.Println` → Standard error-a timestamp ilə yazır
- `os.Exit(1)` — Proqramı 1 exit kodu ilə dayandırır

**Mənbə:** Chapter 2, pages 25-66

### 8. bufio.Scanner — Buffered Oxuma
**Nədir:** Faylları buferli (buffered) şəkildə sətir-sətir və ya token-token oxumaq üçün nəzərdə tutulmuş tip.

**Necə işləyir:** `bufio.NewScanner(file)` Scanner yaradır. `scanner.Scan()` növbəti token-i oxuyur, `scanner.Text()` oxunmuş mətn qaytarır, `scanner.Err()` xəta yoxlayır. SplitFunc ilə tokenization (tokenlərə ayırma) növü dəyişdirilə bilər.

**Nəyə lazımdır:** Böyük faylları yaddaşa yükləmədən emal etmək, line-by-line (sətir-sətir) və ya word-by-word (söz-söz) oxumaq.

**Üstünlükləri:**
- Faylın hamısını yaddaşa yükləməyə ehtiyac yox
- `ScanWords`, `ScanLines` kimi hazır split funksiyaları var

**Çatışmamazlıqları:**
- Token çox böyük olarsa panic (çöküş) baş verə bilər

**Kitabdan kod nümunəsi:**
```go
scanner := bufio.NewScanner(file)
scanner.Split(bufio.ScanWords)
var wordCount int
for scanner.Scan() {
    wordCount++
}
if scanner.Err() != nil {
    log.Println(scanner.Err())
}
```
**Mənbə:** Chapter 2, pages 25-66

### 9. io.Reader Interface
**Nədir:** Oxunan məlumat mənbələri üçün standart interfeys — yeganə metod `Read(p []byte) (n int, err error)`.

**Necə işləyir:** Hər hansı tip `Read` metodunu implement edərsə, `io.Reader` interfeysini ötənə bilər. `os.File`, `bufio.Scanner`, `strings.Reader` və s. bu interfeysi implement edir.

**Nəyə lazımdır:** Generic (ümumi) I/O kodu yazmaq — funksiya `io.Reader` qəbul edərsə, fayl, şəbəkə və ya buffer ilə işləyə bilər.

**Üstünlükləri:**
- Decoupling (ayrılıq) — I/O mənbəyindən asılı olmayan kod
- Standart library-də geniş istifadə

**Çatışmamazlıqları:**
- Oxunmalı byte sayını özü təyin edir — buffer doludursa qayıda bilər

**Mənbə:** Chapter 2, pages 25-66

### 10. Range Loop (Aralıq Döngüsü)
**Nədir:** Kolleksiyalar (array, slice, map, string) üzərində iterasiya üçün xüsusi `range` operatoru.

**Necə işləyir:** `for i, v := range collection` — i indeks, v dəyər (kopya). `for i := range collection` — yalnız indeks. `for _, v := range collection` — indeksi ignore et, yalnız dəyər.

**Nəyə lazımdır:** Kolleksiyaları təkrar etmək, indeksləmək.

**Üstünlükləri:**
- Avtomatik dayanma — bounds (sərhəd) yoxlanışı lazım deyil
- İndeks və dəyər eyni anda əlçatan

**Çatışmamazlıqları:**
- Dəyər kopyalanır — dəyişdirmək üçün indeks istifadə edin

**Mənbə:** Chapter 2, pages 25-66

## Əsas terminlər
- Module (modul) — Go layihəsinin asılılıqları və adı ilə təmsil olunan vahid
- Package (paket) — Funksiya və tiplərin təşkilat vahidi
- go mod — Go modules sistemi
- gofmt — Standart kod formatlaşdırıcı
- io.Reader — Oxuma interfeysi
- bufio.Scanner — Buferli oxuyucu
- Error handling (xəta idarəetməsi) — `error` tipi ilə xətaları idarəetmə
- Range loop (aralıq döngüsü) — Kolleksiya iterasiyası üçün operator
- Zero value (sıfır dəyəri) — Dəyişən yaradıldıqda avtomatik verilən dəyər

## Praktik nəticə
Bu chapter ilk hands-on (praktiki) Go tətbiqini yazdıq: word counter. Go toolchain (kompilyator, `go run`, `go build`, `gofmt`), error handling, `bufio.Scanner` və `io.Reader` konsepsiyalarını praktik nümunələrlə öyrəndik.

## Mənbə
Pages: 25-66
