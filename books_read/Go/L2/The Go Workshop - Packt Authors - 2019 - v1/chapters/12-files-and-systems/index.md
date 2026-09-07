# Chapter 12 — Files and Systems (Fayllar və Sistemlər)

## Bu fəsil nədən bəhs edir?

Fayl sistemi (metadata, iyerarxiya), UNIX fayl icazələri (rwx, oktal/simvolik, owner/
group/others), flag paketi (CLI arqumentlər, default dəyərlər, məcburi flag-lər),
siqnallar (SIGINT/SIGTSTP tutma, signal.Notify, graceful cleanup), fayl yaratma/yazma
(os.Create, Write/WriteString, ioutil.WriteFile), mövcudluq yoxlaması (os.Stat +
IsNotExist, FileInfo), tam oxuma (ReadFile/ReadAll), os.OpenFile flag-ləri (O_APPEND
və s.), fayl backup nümunəsi, CSV (encoding/csv, header keçmə, sahə parse) və bank
transaction activity.

## Əsas fikirlər

### 1. Fayl Sistemi və İcazələr
**3 icazə tipi:** Read (r=4), Write (w=2), Execute (x=1).
**3 subyekt qrupu:** Owner, Group, Others.

**Oktal/simvolik:** `rw-` = 6 (4+2); `rwx` = 7; `-rw-r--r--` = 0644.
`ls -l` ilk simvol: `-` fayl, `d` qovluq.

**Cədvəl:** 0=---, 1=--x, 2=-w-, 3=-wx, 4=r--, 5=r-x, 6=rw-, 7=rwx.

### 2. flag Paketi — CLI Arqumentlər
```go
func Int(name string, value int, usage string) *int   // Bool, String, Float64...
```
- `name` → `-name` flag; `value` → default; `usage` → -h çıxışında görünür
- **Qaytaran POINTER** — `*v` ilə oxu
- **flag.Parse()** — çağırmaq MÜTLƏQ; bundan sonra dəyərlər dolur

**Kitabdan kod nümunəsi:**
```go
i := flag.Int("age", -1, "your age")
n := flag.String("name", "", "your first name")
b := flag.Bool("married", false, "are you married?")
flag.Parse()
fmt.Println(*n, *i, *b)
// ./exFlag -h → usage avtomatik çap olunur!
// Sıra əhəmiyyətsiz; -name=John və ya -name John hər ikisi
```

**Məcburi flag patterni:**
```go
if *n == "" {                       // default = boş → istifadəçi DAXIL ETMƏYİB
    fmt.Println("Name is required.")
    flag.PrintDefaults()            // usage çapı
    os.Exit(1)
}
```
Default dəyər "guard" kimi seçilir (misal dəyərlər: "", -1).

### 3. Siqnallar — Graceful Shutdown
**Problem:** defer-lər os.Exit / Ctrl+C / OS interrupt-da İŞLƏMİR — resurslar açıq
qalır (məs. DepositCheck icra olunmur).

**Kitabdan kod nümunəsi:**
```go
sigs := make(chan os.Signal, 1)          // signal kanalı
done := make(chan bool)                  // exit koordinatoru
signal.Notify(sigs, syscall.SIGINT, syscall.SIGTSTP)   // TUT

go func() {
    for {
        s := <-sigs
        switch s {
        case syscall.SIGINT:              // Ctrl+C
            fmt.Println("My process has been interrupted...")
            cleanUp()                      // TƏMİZLİK — siqnal TUTULDUĞU üçün MÜMKÜNDÜR!
            done <- true
        case syscall.SIGTSTP:            // Ctrl+Z
            fmt.Println("Someone pressed CTRL-Z")
            cleanUp()
            done <- true
        }
    }
}()

fmt.Println("Program is blocked until a signal is caught")
<-done                                // BLOK — siqnal gözlə
fmt.Println("Out of here")
```
**Güc:** interrupt-i TUTUB təmizliyi ÖZÜMÜZ icra edirik — log yazma, fayl silmə,
bağlantı kapatma. signal.Notify imza: `Notify(c chan<- os.Signal, sig ...os.Signal)`
— variadic siqnal siyahısı.

### 4. Fayl Yaratma və Yazma
```go
f, err := os.Create("test.txt")     // touch kimi; mövcuddursa TRUNCATE!
if err != nil { panic(err) }        // panic > os.Exit: defer-lər İŞLƏYİR
defer f.Close()

f.Write([]byte("Using Write function.\n"))     // []byte
f.WriteString("Using WriteString function.\n") // string
```
os.File io.Writer/Reader-ı satisfy edir.

**Tək əmrlə (ioutil):**
```go
err := ioutil.WriteFile("test.txt", []byte("Look!"), 0644)
// Yoxdursa YARADIR (0644); varsa TRUNCATE edir
```

### 5. Mövcudluq Yoxlaması — os.Stat
```go
file, err := os.Stat("junk.txt")
if err != nil {
    if os.IsNotExist(err) {          // xətanı TİPİNƏ görə ayır!
        fmt.Println("File does not exist!")
    }
}
// Mövcuddursa — FileInfo:
file.Name()      // ad
file.Size()      // bayt
file.IsDir()     // qovluqdurmu
file.ModTime()   // dəyişdirilmə vaxtı
file.Mode()      // FileMode
```
Stat birdən çox xəta qaytara bilər — IsNotExist olmadığını yoxla.

### 6. Tam Oxuma — ReadFile/ReadAll
```go
content, err := ioutil.ReadFile("test.txt")     // FAYL adı ilə
content, err := ioutil.ReadAll(f)               // io.Reader ilə — FLEKSİBEL!
r := strings.NewReader("No file here.")          // string belə Reader-dır
ioutil.ReadAll(r)                                // faylsız test mümkün
```
**Diqqət:** kiçik fayllar üçün — böyük fayl = memory exhaustion.
EOF bu funksiyalarda XƏTA SAYILMIR.

### 7. os.OpenFile — Flag Kombinasiyaları
```go
func OpenFile(name string, flag int, perm FileMode) (*File, error)
```
| Flag | Məna |
|---|---|
| O_RDONLY/O_WRONLY/O_RDWR | oxu/yaz/ikisi |
| O_APPEND | SONA əlavə |
| O_CREATE | yoxdursa yarat |
| O_EXCL | O_CREATE ilə: mövcuddursa xəta |
| O_TRUNC | açanda kəs |
| O_SYNC | sinxron I/O |

**Kombinasiyalar:**
```go
os.OpenFile(f, os.O_CREATE, 0644)                          // sadəcə yarat
os.OpenFile(f, os.O_CREATE|os.O_WRONLY, 0644)              // üstünə yaz (overwrite)
os.OpenFile(f, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)  // SONA ƏLAVƏ — ən common
```

### 8. Backup Patterni (kitabdan)
```go
var ErrWorkingFileNotFound = errors.New("The working file is not found.")

func createBackup(working, backup string) error {
    _, err := os.Stat(working)               // 1) mövcudluq
    if err != nil && os.IsNotExist(err) {
        return ErrWorkingFileNotFound         // custom error
    }
    workFile, _ := os.Open(working)           // 2) oxu
    content, _ := ioutil.ReadAll(workFile)    // 3) tam oxu
    return ioutil.WriteFile(backup, content, 0644)  // 4) backup-a yaz
}

func addNotes(workingFile, notes string) error {
    notes += "\n"
    f, err := os.OpenFile(workingFile,
        os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil { return err }
    defer f.Close()
    _, err = f.Write([]byte(notes))          // append
    return err
}
// main: backup et → 10 qeyd əlavə et
```

### 9. CSV — encoding/csv
```go
in := `firstName,lastName,age
Celina,Jones,18
Cailyn,Henderson,13
`
r := csv.NewReader(strings.NewReader(in))     // io.Reader qəbul edir!
for {
    record, err := r.Read()                  // SƏTİR-sətir; []string
    if err == io.EOF { break }                // EOF = bitdi
    if err != nil { log.Fatal(err) }
    fmt.Println(record)                       // [Celina Jones 18]
}
```
**Header keçmə + sahələrə çıxış:**
```go
header := true
for {
    record, err := r.Read()
    if err == io.EOF { break }
    if !header {                              // ilk sətir = başlıqlar — SKIP
        for idx, value := range record {
            switch idx {                       // indeks = sütun
            case 0: fmt.Println("First Name:", value)
            case 1: fmt.Println("Last Name:", value)
            case 2: fmt.Println("Age:", value)
            }
        }
    }
    header = false                            // ilk keçiddən sonra
}
```

### 10. Bank Transaction Activity
CLI: 2 məcburi flag (-c fayl, -l log); restart-da köhnə logu sil; budget kateqoriya
enum-ları (autoFuel, food, mortgage, repairs, insurance, utilities, retirement);
bank category → bizim category mapping funksiyası + naməlum kateqoriyada custom error;
parse → error olsa LOGA YAZ + DAVAM ET; nəticə []transaction.

## Əsas terminlələr
- File Permissions (rwx) — 3 tip × 3 subyekt; oktal (6=rw) / simvolik
- flag.Parse — CLI arqumentlərin işlənməsi
- flag.PrintDefaults — usage çapı
- os.Signal — siqnal tipi; kanal üzərindən
- signal.Notify — siqnalları kanala TUTMAQ
- SIGINT/SIGTSTP — Ctrl+C / Ctrl+Z
- Graceful Shutdown — siqnal tutub təmizliklə çıxma
- os.Create — touch; TRUNCATE davranışı
- Write/WriteString — []byte / string yazma
- ioutil.WriteFile — yarat+yaz bir əmrlə
- os.Stat / IsNotExist — mövcudluq yoxlaması
- FileInfo — Name/Size/IsDir/ModTime/Mode interfeysi
- ReadFile/ReadAll — tam oxuma (kiçik fayl!)
- O_APPEND|O_CREATE|O_WRONLY — append üçün əs üçlük
- csv.Reader — sətir-sətir CSV oxunuşu
- Header Skip — ilk sətiri keçmə patterni

## Praktik nətidə

(1) İcazə oktalla hesabla: rwx sırası 4-2-1; 0644 = rw-r--r--. (2) Flag: Parse
ÇAĞIRMADAN dəyərlər BOŞDUR; pointer-dən * ilə oxu. (3) Məcburi flag: default = guard
dəyər; yoxlamadan sonra PrintDefaults + exit. (4) defer OS interrupt-da ölüdür —
kritik təmizlik üçün siqnal TUT (Notify + kanal loop + done). (5) panic > os.Exit:
defer-lər panic-də YAŞAYIR. (6) os.Create faylı TRUNCATE edir — append üçün
OpenFile(O_APPEND|O_CREATE|O_WRONLY). (7) Stat xətasını IsNotExist ilə TİPLƏ
yoxla — digər xətalar (permission və s.) fərqli idarə tələb edir. (8) ReadAll
io.Reader qəbul edir — strings.NewReader ilə faylsız test; amma memory həddini unutma.
(9) CSV: io.EOF = loop sonu; ilk sətir header — skip patterni ilə keç. (10) Parse
app-larında: xəta → LOGA yaz + davam; hard fail yalnız struktur xətalarında.

## Mənbə
Pages: 421-462 (PDF 454-497)
