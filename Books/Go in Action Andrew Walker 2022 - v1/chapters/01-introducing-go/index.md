# Chapter 1 — Introducing Go

## Bu chapter nədən bəhs edir?

Bu chapter, Go proqramlaşdırma dilinin əsaslarını və onu nəyə görə seçmək lazım olduğunu izah edir. Müasir proqramlaşdırma çağırışlarını (concurrency, performance, development speed) necə həll etdiyini, Go-nun əsas xüsusiyyətlərini və tip sistemini əhatə edir.

## Əsas fikirlər

### 1. Development Speed (İnkişaf Sürəti)
**Nədir:** Go kompilyatorunun çox sürətli işləməsi və qısa compilation (kompilyasiya) vaxtı.

**Necə işləyir:** Go, sadə asılılıq (dependency) həlli mexanizmi və ağıllı kompilyator vasitəsilə böyük tətbiqləri bir saniyədən qısa vaxtda kompilyə edir. Bütün Go source (mənbə) kodu 20 saniyədən qısa vaxtda kompilyə olunur.

**Nəyə lazımdır:** Developer-lərin kod yazmaqdan icrasına qədərki vaxtı minimala endirmək, dinamik dillərdəki sürətli inkişaf üstünlüyünü static (statik) dillə birləşdirmək.

**Üstünlükləri:**
- Kompilyasiya çox sürətli — böyük layihələr də bir saniyədən qısa vaxtda hazır olur
- Go Modules vasitəsilə asılılıqlar şəffaf və təkrarlanabilir şəkildə idarə olunur
- Dinamik dillər kimi sürətli inkişaf, statik dillər kimi yüksək performans

**Çatışmamazlıqları:**
- Dinamik dillərdən daha az intepretasiya (təfsir) dəstəyi
- Kütubxanaların yenilənməsi üçün kompilyasiya tələb olunur

**Kitabdan kod nümunəsi:**
```go
// Go programının əsas strukturı
package main

import "fmt"

func main() {
    fmt.Println("Hello, Go!")
}
```
**Sub-kod izahı:**
- `package main` → Executable (icra edilə bilən) programın giriş nöqtəsi
- `import "fmt"` → Formatlama və çıxış üçün standart kütubxana
- `func main()` → Program işə düşəndə ilk çağırılan funksiya

**Mənbə:** Chapter 1, pages 6-24

### 2. Type Safety (Tip Təhlükəsizliyi)
**Nədir:** Kompilyatorun tip uyğunsuzluqlarını compile (kompilyasiya) vaxtında yoxlaması.

**Necə işləyir:** Hər dəyər yaradıldıqda tip təyin olunur və kompilyator bunu izləyir. Tip uyğunsuzluğu olarsa, proqram heç vaxt icra edilmir, kompilyasiya xətası verilir.

**Nəyə lazımdır:** Runtime (icra vaxtı) xətalarınıMinimuma endirmək, proqramın daha təhlükəsiz və etibarlı işləməsi.

**Üstünlükləri:**
- Tip xətaları runtime-də deyil, compile vaxtında bərpa olunur
- Debugging (səhv axtarış) asanlaşır
- Dinamik dillərdə yaygın olan "undefined is not a function" kimi xətalar qarşısı alınır

**Çatışmamazlıqları:**
- Static typing bəzi hallarda kod yazmağı daha mürəkkəb edə bilər
- Dinamik dillərdəki kimi flexible (flexible / çevik) tip dəyişiklikləri mümkün deyil

**Kitabdan kod nümunəsi:**
```go
var number int
var str string
number = 1
str = "one"
number = str  // Error: cannot use str (variable of type string) as type int in assignment
```
**Sub-kod izahı:**
- `var number int` → Tam ədəd tipində dəyişən
- `var str string` → Mətn tipində dəyişən
- `number = str` → Tip uyğunsuzluğu — kompilyasiya xətası

**Mənbə:** Chapter 1, pages 6-24

### 3. Concurrency (Paralellik)
**Nədir:** Eyni anda bir neçə tapşırığı icra etmək qabiliyyəti.

**Necə işləyir:** Go runtime (işləmə vaxtı) concurrency-ni daxili olaraq dəstəkləyir. Goroutines (gözaltı proseslər) və Channels (kanallar) vasitəsilə paralel işləmə təmin edilir. Goroutine-lər çox yüngül və yarpaq OS thread (əməliyyat sistemi thread) paylaşırlar.

**Nəyə lazımdır:** Yüksək throughput (ötürmə qabiliyyəti), real-time (real-time / canlı) əməliyyatlar, I/O-bound (Giriş/Çıxış əlaqəli) tapşırıqlar.

**Üstünlükləri:**
- Çox yüngül goroutine-lər — minlərlə eyni anda işləyə bilər
- Channel-lər vasitəsilə təhlükəsiz data (məlumat) mübadiləsi
- Mürəkkəb synchronization (sinxronizasiya) mexanizmləri yerine sadə kanal modeli

**Çatışmamazlıqları:**
- Deadlock (sökülmə) riski — kanallar düzgün deyil
- Memory leak (yaddaq sızqanı) ehtimalı — goroutine-lər bitmirsə

**Kitabdan kod nümunəsi:**
```go
func DoSomethingSlow() {
    fmt.Println("SLOW: maybe I'm doing network stuff?")
    time.Sleep(1 * time.Second)
    fmt.Println("SLOW: Okay, finally finished!")
}

func main() {
    fmt.Println("Starting the main task...")
    go DoSomethingSlow()
    fmt.Println("Resuming the main task...")
    time.Sleep(500 * time.Millisecond)
    fmt.Println("Finished the main task!")
    time.Sleep(1 * time.Second)
}
```
**Sub-kod izahı:**
- `go DoSomethingSlow()` → Funksiyonu goroutine kimi işə salır
- `time.Sleep(...)` → Verilən müddət üçün dayandırır
- `fmt.Println(...)` → Konsola mətn çıxışı verir

**Mənbə:** Chapter 1, pages 6-24

### 4. Channels (Kanallar)
**Nədir:** Goroutine-lər arasında data (məlumat) mübadiləsi üçün nəzərdə tutulmuş tip.

**Necə işləyir:** Kanal, göndərən və qəbul edən goroutine-ləri sinxronlaşdırır. Buffer (bufer) kanalı birdən çox mesaj saxlaya bilər, unbuffered (bufer-siz) kanal isə həmişə bir mesaj tutur.

**Nəyə lazımdır:** Concurrent (paralel) proqramlarda data paylaşımı və synchronization.

**Üstünlükləri:**
- Lock (kilid) və mutex (mütəks) istifadə etmədən təhlükəsiz data mübadiləsi
- Göndərən və qəbul edən arasında avtomatik sinxronizasiya

**Çatışmamazlıqları:**
- Unbuffered kanallarda göndərən və qəbul edən eyni anda hazır olmalıdır
- Pointer (göstərici) məlumatları kopyalandıqda, yalnız pointer özü kopyalanır, əsas data yox

**Kitabdan kod nümunəsi:**
```go
func main() {
    channel := make(chan string)
    output := make(chan string)
    
    go func() {
        s := <-channel
        s = s + "World!"
        output <- s
    }()
    
    go func() {
        s := "Hello, "
        channel <- s
    }()
    
    finalString := <-output
    fmt.Println("This string was built concurrently:", finalString)
}
```
**Sub-kod izahı:**
- `make(chan string)` → String (mətn) tipli kanal yaradır
- `<-channel` → Kanaldan məlumat qəbul edir
- `output <- s` → Kanala məlumat göndərir

**Mənbə:** Chapter 1, pages 6-24

### 5. Go Type System (Go Tip Sistemi)
**Nədir:** Go-nun tip sistemi — inheritance (miras) yox, composition (birləşmə) və implicit interface (daxili interfeys) implementasiyası əsasında işləyir.

**Necə işləyir:** Struct (strukt) tipləri sadə tiplərdən birləşdirilir. Interface-lər davranışı (method-ları) təyin edir. Tip interfeysi implement edərsə, avtomatik olaraq o interfeysin nümayəndəsi sayılır.

**Nəyə lazımdır:** Esnek və reusable (təkrar istifadə edilən) kod yazmaq, class hierarchy (sinif ierarxiyası) çatışmamazlıqlarını aradan götürmək.

**Üstünlükləri:**
- No inheritance (miras yox) — daha sadə və aydın kod
- Implicit interface (daxili interfeys) — explicit (açıq) implementasiya tələb olunmur
- Composition over inheritance (miras əvəzinə birləşmə) — daha esnek strukturlar

**Çatışmamazlıqları:**
- OOP (obyekt-yönümlü proqramlaşdırma) pattern-ləri bir qədər fərqli işləyir
- Müəyyən hallarda boilerplate (standart) kod artır

**Kitabdan kod nümunəsi:**
```go
type Describer interface {
    Describe() string
}

type Person struct {
    FirstName string
    LastName string
    Age int
}

func (p Person) Describe() string {
    return fmt.Sprintf("%s is a %d year-old human", p.Name(), p.Age)
}
```
**Sub-kod izahı:**
- `type Describer interface` → Davranış təyin edən interfeys
- `func (p Person) Describe()` → Person tipinə metod əlavə edir
- `fmt.Sprintf(...)` → Formatlı mətn yaradır

**Mənbə:** Chapter 1, pages 6-24

### 6. Generics (Generik Proqramlaşdırma)
**Nədir:** Type parameter (tip parametri) ilə bir neçə tip üçün işləyən universal funksiya və data struktur yaratmaq imkanı.

**Necə işləyir:** `[N int | float64]` kimi type constraint (tip məhdudiyyəti) verilir və funksiya həmin constraintə uyğun tiplərlə işləyir.

**Nəyə lazımdır:** Daha generic (ümumi) və reusable kod, type-safe (tip təhlükəsiz) data strukturları.

**Üstünlükləri:**
- Eyni funksiyanı bir neçə tip üçün istifadə etmək
- Type safety saxlanır — compile vaxtında yoxlanış
- Interface-lərdən daha sürətli işləyir

**Çatışmamazlıqları:**
- Yazmağı daha mürəkkəb edə bilər
- Type constraint izahı tələb olunur

**Kitabdan kod nümunəsi:**
```go
func SumNumbers[N int | float64](numberSlice ...N) N {
    var total N
    for i := range numberSlice {
        total += numberSlice[i]
    }
    return total
}

func main() {
    fmt.Println(SumNumbers(1, 2, 3))       // 6
    fmt.Println(SumNumbers(1.1, 2.2, 3.3)) // 6.6
}
```
**Sub-kod izahı:**
- `[N int | float64]` → Tip parametri N, yalnız int və ya float64 qəbul edir
- `numberSlice ...N` → Variadic (dəyişən saylı) parametr
- `range numberSlice` → Slice-i iterasiya edir

**Mənbə:** Chapter 1, pages 6-24

### 7. Garbage Collection (Zibil Toplayıcı)
**Nədir:** Avtomatik yaddaş (memory) idarəetmə sistemi — istifadə edilməyən yaddaş avtomatik olaraq azad edilir.

**Necə işləyir:** Go runtime-i, reference (referans) saymaq mexanizmi ilə hansı obyektlərin artıq istifadə edilmədiyini müəyyən edir və onları collect (toplayır) edir.

**Nəyə lazımdır:** Manual (əl ilə) memory management (yaddaş idarəetmə) çatışmamazlıqlarını aradan götürmək.

**Üstünlükləri:**
- Memory leak (yaddaş sızqanı) riskini azaldır
- Programmer (proqramçı) yaddaş idarəetməsi ilə məşğul olmur
- Fully-concurrent (tam paralel) — proqramı dayandırmadan collect edir

**Çatışmamazlıqları:**
- Sub-millisecond (milyardıma saniyə) pauses (fasilələr) olsa da, performansa təsir edə bilər
- Predictable (təxmin edilə bilən) latency (gecikmə) tələb edən sistemlərdə problem ola bilər

**Mənbə:** Chapter 1, pages 6-24

## Əsas terminlər
- Goroutine (gözaltı proses) — Go dilində paralel işləmə üçün yüngül funksiya
- Channel (kanal) — Goroutine-lər arasında data mübadiləsi üçün nəzərdə tutulmuş tipli strukturlar
- Concurrency (paralellik) — Eyni anda bir neçə tapşırığı idarəetmə qabiliyyəti
- Type Safety (tip təhlükəsizliyi) — Kompilyatorun tip uyğunsuzluqlarını compile vaxtında yoxlaması
- Interface (interfeys) — Davranış təyin edən, tip implementasiyasına imkan verən strukturlar
- Generics (generik proqramlaşdırma) — Universal funksiya və data strukturları yaratmaq imkanı
- Garbage Collection (zibil toplayıcı) — Avtomatik yaddaş idarəetmə sistemi
- Go Modules — Asılılıqları idarəetmək üçün Go ecosystem-ünün standart vasitəsi

## Praktik nəticə
Bu chapter, Go-nu seçmək üçün əsas səbəbləri izah edir: sürətli kompilyasiya, təhlükəsiz tip sistemi, daxili paralellik və avtomatik yaddaş idarəetməsi. Kitabın sonrakı chapter-lərində bu konsepsiyalar praktik nümunələrlə əhatə olunacaq.

## Mənbə
Pages: 6-24
