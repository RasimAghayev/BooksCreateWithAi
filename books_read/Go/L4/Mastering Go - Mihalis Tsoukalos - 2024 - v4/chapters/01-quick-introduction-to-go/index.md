# Chapter 1 — A Quick Introduction to Go (Go-ya Sürətli Giriş)

## Bu chapter nədən bəhs edir?

Go-nun tarixi və fəlsəfəsi, üstünlük/çatışmamazlıqları, go doc/godoc alətləri, Hello World,
dəyişənlər, kontrol strukturları (if/switch), for loop və range, istifadəçi input-u (stdin +
command line arguments), goroutine-lərə ilk baxış, which(1) utility-nin Go implementasiyası,
log paketi (syslog, custom log faylları, Fatal/Panic) və kitab boyu inkişaf edəcək statistika
tətbiqinin ilk versiyası.

## Əsas fikirlər

### 1. Go-nun Tarixi və Fəlsəfəsi
**Nədir:** 2009-cu ildə Google-dən çıxan open-source systems programming dili. Müəlliflər:
Robert Griesemer, Ken Thomson, Rob Pike.

**Necə işləyir:** Simplicity (sadəlik) əsas dizayn qərarıdır — hər addım üçün bir yol. Go C
sintaksisinə, Modula-2 paket konseptinə bənzəyir. Rəsmi ad "Go"dur; "Golang" qeyri-rəsmi
adlandırmadır (go.org domain-i alınmadığı üçün golang.org seçildi; hazırda rəsmi sayt
https://go.dev/).

**Nəyə lazımdır:** Böyük komandalarda saxlanıla bilən, oxunaqlı, etibarlı sistem proqramları.

**Üstünlükləri:**
- 25 açar söz — öyrənməsi asan; C/Python/Java bilənlər üçün tanış
- Göstərişi arimetikası YOX (unsafe paketi istisna)
- Goroutine + channel ilə built-in paralellik
- Statik link olunan binary-lar — asılılıqsız paylama
- Backward compatibility (Go 1.x daxilində)
- Compiler sadə səhvləri (istifadə olunmamış import/dəyişən) tutur

**Çatışmamazlıqları:**
- Tam OOP dəstəyi YOX (klass və inheritance yoxdur; interfeyslərlə mimik olunur)
- Goroutine OS thread qədər güclü deyil; fork(2) dəstəyi yoxdur
- Manual memory management mümkün deyil
- High-availability sistemlər üçün Erlang/Elixir daha yaxşıdır

### 2. go doc və godoc Alətləri
**Nədir:** Standart kitabxana sənədlərini offline oxumaq üçün alətlər — man(1)-ın Go
variantı.

**Kitabdan kod nümunəsi:**
```bash
go doc fmt.Printf      # funksiya sənədi
go doc fmt             # bütün paket sənədi
go install golang.org/x/tools/cmd/godoc@latest   # godoc qurulumu
~/go/bin/godoc -http=:8001                          # lokal web server (:6060 default)
```
**Sub-komanda izahı:**
- `go doc fmt.Printf` → Printf funksiyasının imzası və izahı
- `-http=:8001` → brauzerdə http://localhost:8001/ ilə sənədləşmə; port 0-1023 root-a məxsusdur

### 3. Hello World — Paket və Funksiya Strukturu
**Kitabdan kod nümunəsi:**
```go
package main

import "fmt"

func main() {
    fmt.Println("Hello World!")
}
```
**Sub-kod izahı:**
- `package main` → muxtariq (executable) proqram üçün mütləq; entry point buradan başlayır
- `import "fmt"` → formatlı I/O; istifadə olunmayan import compile xətası verir
- `func main()` → proqramın başlanğıc nöqtəsi; **kiçik hərflə başlayan hər şey privatdır**
  (yalnız package main istisna olmaqla — orada main() məcburidir)

### 4. go build və go run
**Nədir:** İki icra üsulu — compile olunmuş binary paylamaq və ya script kimi işə salmaq.

**Kitabdan kod nümunəsi:**
```bash
go build hw.go              # hw binary yaradır
./hw                        # manual icra
go build -o helloWorld hw.go  # fərqli ad ilə binary
go run hw.go                # müvəqqəti binary yaradır, icra edir, silir
```
**Sub-komanda izahı:**
- `build` → paylanacaq binary; `-o` → çıxış adı/yolu dəyişir
- `run` → test üçün üstün; binary diskdə qalmır

### 5. Curly Brace Qaydası (Vacib Sintaksis Qaydası)
**Problem:** `{` yeni sətirdə açılırsa compile xətası:
```
./curly.go:7:6: missing function body
./curly.go:8:1: syntax error: unexpected semicolon or newline before {
```
**Səbəb:** Compiler lazım olan yerlərdə avtomatik semicolon daxil edir; `func main()` sonunda
semicolon daxil edilir və `{` artıq yeni statement başlayır. **Qayda:** `{` həmişə eyni sətirdə.

### 6. Dəyişənlər, Constants və Tip Çevirmələri
**Necə işləyir:** `var ad tip`, `var ad = dəyər` (tip çıxarılır), `ad := dəyər` (short
assignment — yalnız funksiya daxilində; funksiya xaricində hər statement açar sözlə başlamalıdır).

**Vacib qaydalar:**
- İlkin dəyər yoxdursa → tipin **zero value**-si avtomatik təyin olunur
- **Implicit tip çevirməsi QADAĞANDIR** — `math.Abs(float64(x))` kimi açıq cast mütləqdir
- Global dəyişənlər üçün yalnız `var` və `const` işləyir; `:=` funksiya daxilindədir

**Kitabdan kod nümunəsi:**
```go
var Global int = 1234
var AnotherGlobal = -5678

func main() {
    var j int                      // j = 0 (zero value)
    i := Global + AnotherGlobal    // short assignment
    k := math.Abs(float64(AnotherGlobal))  // açıq tip çevirməsi MÜTLƏQ
    fmt.Printf("Global=%d, i=%d, j=%d k=%.2f.\n", Global, i, j, k)
}
```
- `fmt.Println` → avtomatik yeni sətir; `fmt.Printf` → format nəzarəti (%d int, %.2f 2 onluqlu
  float, \n əl ilə)

### 7. if və switch — Kontrol Strukturları
**Necə işləyir:** if-də mötərizə istifadə olunmur. switch-in 2 forması var: ifadəli (expression)
və ifadəsiz (hər case ayrıca bool ifadə). `fallthrough` — növbəti case-i məcburi icra edir
(Go-da break avtomatikdir).

**Ən məşhur Go pattern-i (error yoxlaması):**
```go
err := anyFunctionCall()
if err != nil {
    // xəta var
}
```

**Kitabdan kod nümunəsi (ifadəli + ifadəsiz switch):**
```go
switch argument {
case "0":
    fmt.Println("Zero!")
case "2", "3", "4":
    fmt.Println("2 or 3 or 4")
    fallthrough
default:
    fmt.Println("Value:", argument)
}

switch {
case value == 0:
    fmt.Println("Zero!")
case value > 0:
    fmt.Println("Positive integer")
case value < 0:
    fmt.Println("Negative integer")
}
```
**Sub-kod izahı:**
- case sırası vacibdir — yalnız ilk uyğun gələn icra olunur
- `strconv.Atoi(argument)` → string→int + err; xəta varsa `err != nil`

### 8. for Loop və range — Tək Loop Konstruksiyası
**Nədir:** Go-da yalnız `for` var; while və sonsuz loop onun variantlarıdır.

**Kitabdan kod nümunəsi:**
```go
for i := 0; i < 10; i++ { }        // ənənəvi

for ok := true; ok; ok = (i != 10) { }  // şərtli

for {                              // sonsuz (while true)
    if i == 10 {
        break
    }
    i++
}

aSlice := []int{-1, 2, 1, -1, 2, -2}
for i, v := range aSlice {         // index + value
    fmt.Println("index:", i, "value:", v)
}
for i := range aSlice { }          // yalnız index — ikinci dəyər tələb olunmur
```
**Sub-kod izahı:**
- `range` → array/slice/map üzrə iterasiya; `_` istənməyən dəyəri udur
- `break` → loop-dan çıxış; `continue` → növbəti iterasiya
- Bədən 1 sətir olsa belə `{ }` mütləqdir

### 9. İstifadəçi Input-u — stdin və Command Line Arguments
**fmt.Scanln:**
```go
fmt.Printf("Please give me your name: ")
var name string
fmt.Scanln(&name)      // pointer ötürülür
```

**os.Args strukturu:**
- `os.Args[0]` → həmişə executable-ın öz yolu
- `os.Args[1:]` → istifadəçi arqumentləri; boş slice heç vaxt olmur

**Kitabdan kod nümunəsi (min/max axtarışı, etibarsız input filtri):**
```go
arguments := os.Args
var min, max float64
var initialized = 0
for i := 1; i < len(arguments); i++ {
    n, err := strconv.ParseFloat(arguments[i], 64)
    if err != nil {
        continue              // etibarsız inputu ötür
    }
    if initialized == 0 {      // ilk etibarlı dəyərlə min/max init
        min = n
        max = n
        initialized = 1
        continue
    }
    if n < min { min = n }
    if n > max { max = n }
}
```
**Sub-kod izahı:**
- `strconv.ParseFloat(s, 64)` → (float64, error); xəta varsa input yoxlanılır
- Error dəyişəni ilə **tip fərqləndirmə** texnikası: əvvəl Atoi (spesifik), sonra ParseFloat
  (generik) — hər etibarlı integer eyni zamanda float olduğundan sıra vacibdir
- `invalid := make([]string, 0)` → etibarsız dəyərləri topla; `append(invalid, k)` ilə əlavə

### 10. Go Paralellik Modelinə İlk Baxış
**Nədir:** Goroutine = ən kiçik icra vahidi; `go` açar sözü ilə yaradılır (funksiya və ya
anonim funksiya). Channel = goroutine-lərarası kommunikasiya mexanizmi.

**Kitabdan kod nümunəsi:**
```go
func myPrint(start, finish int) {
    for i := start; i <= finish; i++ {
        fmt.Print(i, " ")
    }
    fmt.Println()
    time.Sleep(100 * time.Microsecond)
}

func main() {
    for i := 0; i < 4; i++ {
        go myPrint(i, 5)      // 4 goroutine paralel
    }
    time.Sleep(time.Second)   // ƏLƏ HƏLL — WaitGroup düzgün üsuldur (Ch8)
}
```
**Vacib müşahidə:** Hər icra-da fərqli output — goroutine-lər random sıra ilə başlayır; Go
scheduler onları idarə edir. main() çıxsa bütün proqram dayanır — buna görə Sleep ilə gözləmə
buradadır (production-da sync.WaitGroup).

### 11. which(1) Utility-nin Go Versiyası
**Nədir:** UNIX utility — PATH dəyişənində verilmiş executable-ı axtarır.

**Kitabdan kod nümunəsi:**
```go
arguments := os.Args
if len(arguments) == 1 {
    fmt.Println("Please provide an argument!")
    return
}
file := arguments[1]
path := os.Getenv("PATH")
pathSplit := filepath.SplitList(path)
for _, directory := range pathSplit {
    fullPath := filepath.Join(directory, file)
    fileInfo, err := os.Stat(fullPath)
    if err != nil {
        continue
    }
    mode := fileInfo.Mode()
    if !mode.IsRegular() {
        continue
    }
    if mode&0111 != 0 {           // executable bit
        fmt.Println(fullPath)
        return
    }
}
```
**Sub-kod izahı:**
- `os.Getenv("PATH")` → environment dəyərini oxuyur
- `filepath.SplitList` → PATH-i portable şəkildə ayırır (Linux `:`, Windows `;`)
- `filepath.Join` → OS-ə uyğun path birləşdirir
- `os.Stat` → fayl məlumatı; `mode.IsRegular()` → adi fayldırmı; `mode&0111` → icra biti
- Tapılmayanda heç nə çap etməmək — UNIX fəlsəfəsi (pipe-ya uyğun)

### 12. log Paketi — Logging Sistemi
**Nədir:** Standart kitabxananın logging qatı; default stderr-ə yazır.

**Səviyyə/Facility haqqında bilik:** UNIX logging **severity level** (debug, info, notice,
warning, err, crit, alert, emerg) və **facility** (auth, cron, daemon, mail, user, local0-7)
konseptlərinə malikdir; Go-nun log paketi səviyyə dəstəkləmir, amma syslog-a göndərə bilir.

**Kitabdan kod nümunələri:**

syslog-a yazış:
```go
sysLog, err := syslog.New(syslog.LOG_SYSLOG, "systemLog.go")
if err != nil {
    log.Println(err)
    return
}
log.SetOutput(sysLog)
log.Print("Everything is fine!")
```

Custom log faylı:
```go
LOGFILE := path.Join(os.TempDir(), "mGo.log")
f, err := os.OpenFile(LOGFILE, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
if err != nil {
    fmt.Println(err)
    return
}
defer f.Close()                       // funksiya qayıtmazdan ƏVVƏL icra
iLog := log.New(f, "iLog ", log.LstdFlags)
iLog.Println("Hello there!")
```

Sətir nömrəsi ilə:
```go
LstdFlags := log.Ldate | log.Lshortfile
iLog := log.New(f, "LNum ", LstdFlags)   // Lshortfile → fayl adı + sətir
iLog.SetFlags(log.Lshortfile | log.LstdFlags)   // icra zamanı format dəyişə bilər
```

Çoxlu log hədəfi:
```go
w := io.MultiWriter(file, os.Stderr)     // eyni anda fayl + stderr
logger := log.New(w, "myApp: ", log.LstdFlags)
logger.Printf("BOOK %d", os.Getpid())
```

**Sub-kod izahı:**
- `log.Fatal()` → log.Print + os.Exit(1) — programı dərhal bitirir, non-zero exit code
- `log.Panic()` → log.Print + panic() — stack trace ilə dayanır; gözlənilməz vəziyyətlər üçün
- `os.O_APPEND|os.O_CREATE|os.O_WRONLY` → mövcuddursa əlavə et, yoxdursa yarad
- `0644` → UNIX permission (owner yazma, hamı oxuma)
- **Müəllim qeydi (kitabdan):** logging tətbiq kodu üçündür, library kodu üçün YOX;
  Docker/cloud dünyasında stdout/stderr-ə loglama 12factor.net/logs tövsiyəsidir

### 13. Statistika Tətbiqi — Kitabın Davamlı Layihəsi
**Nədir:** Kitab boyu fəsil-fəsil inkişaf etdiriləcək CLI utility — N dəyərdən min/max/mean/
standard deviation hesablayır.

**Kitabdan kod nümunəsi:**
```go
var min, max float64
var initialized = 0
nValues := 0
var sum float64
for i := 1; i < len(arguments); i++ {
    n, err := strconv.ParseFloat(arguments[i], 64)
    if err != nil {
        continue
    }
    nValues = nValues + 1
    sum = sum + n
    if initialized == 0 {
        min = n
        max = n
        initialized = 1
        continue
    }
    if n < min { min = n }
    if n > max { max = n }
}
meanValue := sum / float64(nValues)
fmt.Printf("Mean value: %.5f\n", meanValue)

// Standard deviation — mean tələb etdiyi üçün ikinci pass
var squared float64
for i := 1; i < len(arguments); i++ {
    n, err := strconv.ParseFloat(arguments[i], 64)
    if err != nil {
        continue
    }
    squared = squared + math.Pow((n-meanValue), 2)
}
standardDeviation := math.Sqrt(squared / float64(nValues))
```
**Sub-kod izahı:**
- Mean üçün bütün dəyərlər tələb olunduğundan əvvəlcə hamısı toplanır
- Standard deviation = sqrt(Σ(x-mean)²/n) — mean hazır olduqdan sonra ikinci dövr
- `math.Pow`, `math.Sqrt` → float64 tələb edir

## Əsas terminlər
- Goroutine — yüngül icra vahidi, `go` açar sözü ilə
- Channel (kanal) — goroutine-lərarası data mübadiləsi
- Zero Value (sıfır dəyəri) — ilkin dəyərsiz dəyişənin default dəyəri
- Short Assignment Statement (qısa təyinat) — `:=` ilə tip çıxarılması
- Command Line Arguments (komanda sətir arqumentləri) — os.Args slice
- Logging Facility (log kateqoriyası) — syslog-da mesajın mənsub olduğu kateqoriya
- Logging Level (log səviyyəsi) — severity: debug...emerg
- Static Linking (statik bağlama) — binary-nın asılılıq tələb etməməsi
- Backward Compatibility (geriyə uyğunluq) — Go 1.x daxili zəmanət

## Praktik nəticə

(1) Go-nun 25 açar sözü və bir-yol prinsipi — oxunaqlılıq dizayn qərarıdır; (2) implicit
tip çevirməsi heç bir zaman işləmir — hər cast açıq; (3) istifadə olunmamış dəyişən/import
compile xətasıdır, nəticədə ölü kod olmur; (4) error-nun nil olub-olmaması Go-nun əsas
nəzarət axını pattern-idir; (5) `go run` test üçün, `go build` paylama üçün; (6) main
çıxanda goroutine-lər də ölür — SLEEP ilə gözləmə yalnız dərs nümunəsidir, WaitGroup düzgün
yoldur; (7) logging stderr-ə/cloud-native şəkildə; kitabxanalara log QADAĞANDIR; (8) which(1)
nümunəsi: filepath paketi portable path emalı üçün standart yoldur.

## Mənbə
Pages: 1-44 (PDF 32-77)
