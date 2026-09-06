# Chapter 1 — Introducing Go

## Bu chapter nədən bəhs edir?

Go proqramlaşdırma dilinin niyə yaradıldığını, hansı problemləri həll etdiyini və əsas fərqlərini (sürətli kompilyasiya, concurrency, sadə type sistemi, garbage collection) izah edir. Chapter müəllifin "kitabı intermediate developer üçün yazıldığı" qeydi ilə başlayır və Hello World + Go Playground ilə bitir.

## Əsas fikirlər

### 1. Go hansı problemi həll edir?
**Nədir:** Proqramçılar arasında "sürətli inkişaf" (Ruby, Python) və "sürətli icra" (C, C++) arasında seçim problemi var. Go bu iki dünyanı birləşdirir.

**Necə işləyir:** Ağıllı kompilyator və sadələşdirilmiş asılılıq resolüsiyası — kompilyator yalnız birbaşa import edilən kitabxanaları yoxlayır, bütün asılılıq zəncirini gəzmir (Java/C/C++-dan fərqli). Nəticə: böyük Go proqramları 1 saniyədən az, bütün Go source tree 20 saniyədən az müddətdə kompilyasiya olunur.

**Nəyə lazımdır:** Nəhəng layihələrdə build vaxtını qısaldmaq, developer məhsuldarlığını artırmaq.

**Üstünlükləri:**
- Compile time (kompilyasiya vaxtı) saniyə səviyyəsindədir
- Static typing (statik tip) — runtime type bug-larını kompilyasiya vaxtı yaxalayır (dynamic dillərdən fərqli olaraq `ID` sahəsinin int/string/UUID olduğunu soruşmağa ehtiyac yoxdur)

**Çatışmamazlıqları:**
- Dynamic dillərdəki "yaz və dərhal işə sal" rahatlığı tam yoxdur (compile addımı var)

### 2. Goroutines (qorutinlər)
**Nədir:** Başqa gorutinlər ilə paralel/eyni vaxtda işləyən funksiyalar. Thread-lərə bənzəyir, amma daha az yaddaş istifadə edir.

**Necə işləyir:** Bir çox goroutine tək bir OS thread üzərində icra oluna bilir. Go runtime onları konfiqurasiya olunmuş logical processor-lərə (məntiqi prosessorlar) avtomatik schedulе edir, hər logical processor bir OS thread-ə bağlanır. `net/http` kitabxanasında hər inbound request avtomatik öz goroutine-da işə düşür.

**Nəyə lazımdır:** Web server-də eyni anda minlərlə request-i işləmək, fon işləri (logging, processing) əsas axından ayırmaq.

**Üstünlükləri:**
- Minimal overhead — on minlərlə goroutine yaratmaq normaldır
- Thread-lərdən fərqli olaraq sinxronizasiya kodu tələb etmir

**Çatışmamazlıqları:**
- Paylaşılan data oxu-yazması olduqda yenə də sinxronizasiya (channel/lock) tələb olunur

**Kitabdan kod nümunəsi:**
```go
func log(msg string){
        // ... some logging code here
}
// Elsewhere in our code after we've discovered an error.
go log("something dire happened")
```

**Sub-kod izahı:**
- `go log(...)` → `go` açar sözü log funksiyasını goroutine kimi planlaşdırır, kodun qalan hissəsi bloklanmadan davam edir

### 3. Channels (kanallar)
**Nədir:** Gorutinlər arasında tipləşdirilmiş mesajlaşmanı təmin edən data strukturları. Sinxronizasiya daxildir (built-in).

**Necə işləyir:** Unbuffered channel-dan data göndərən goroutine, alan goroutine hazır olana qədər gözləyir — hand-off sinxronlaşdırılır. Birinci goroutine data-nı ikinciyə ötürür, ikinci işi bitirib üçüncüyə ötürür. Heç bir lock tələb olunmur.

**Nəyə lazımdır:** Data-nın eyni anda yalnız bir goroutine tərəfindən dəyişdirilməsini təmin etmək — "shared memory" problemlərindən qaçmaq.

**Üstünlükləri:**
- Data-safe mübadilə, lock-lara ehtiyac yoxdur
- "Data-nı goroutine-lər arası göndər, döyüşdürmə" modeli

**Çatışmamazlıqları:**
- Data copy send olunursa hər goroutine öz nüsxəsində təhlükəsiz dəyişiklik edə bilir; amma **pointer** mübadilə olunursa, yenə sinxronizasiya lazımdır — channel özü pointer-ə qarşı access protection vermir

### 4. Composition over inheritance (kompozisiya irsiyyət əvəzinə)
**Nədir:** Go-nun hierarşiyasız (hierarchy-free) type sistemi. Klassik OOP irsiyyəti (inheritance) yox, kiçik tiplərin böyüklərə daxil edilməsi (embedding) var.

**Necə işləyir:** `Client extends User extends Entity` zənciri əvəzinə kiçik tiplər (`Customer`, `Admin`) yaradılır və böyük tiplərə embed olunur. "Truck: Drive + Carry cargo + Carry passengers" kompozisiya nümunəsi — `Vehicle`-dən irsiyyət almaq əvəzinə lazımi xüsusiyyətlər birbaşa embed edilir.

**Nəyə lazımdır:** Kod təkrarından qurtulmaq, refactoring overhead-ini azaltmaq.

**Üstünlükləri:**
- Refactoring zamanı class hierarxiyasını yenidən planlamağa ehtiyac yoxdur
- Java/C++-dəki "bir həftə abstract class planlama" məşğuliyyəti aradan qalxır

**Çatışmamazlıqları:**
- Klassik OOP-ya öyrəşmiş developer-ə adaptasiya vaxtı lazımdır

### 5. Go interfeysləri — duck typing
**Nədir:** Go interfeysləri **davranışı** (behavior) model edir, tipi yox. İnterfeysi implement etdiyini elan etmək lazım deyil.

**Necə işləyir:** "Quacks like a duck → duck" prinsipi. Tip interfeysin metodlarını implement edirsə, o tipli dəyər həmin interfeys tipində saxlana bilər. Xüsusi declaration (elan) tələb olunmur.

**Kitabdan kod nümunəsi (Java müqayisəsi):**
```java
interface User {
    public void login();
    public void logout();
}
```
Java-da interfeysi implement edən class bütün vədləri yerinə yetirib **açıq şəkildə** `implements` elan etməlidir.

Go-da isə ən çox istifadə olunan interfeyslərdən biri `io.Reader`:
```go
type Reader interface {
    Read(p []byte) (n int, err error)
}
```

**Sub-kod izahı:**
- `type Reader interface` → Reader adlı interfeys tipi təyin edir
- `Read(p []byte) (n int, err error)` → tək metod: byte slice qəbul edir, oxunan byte sayı və mümkün xəta qaytarır
- Bu interfeysi implement etmək üçün yalnız `Read` metodunu yazmaq kifayətdir — `implements` açar sözü YOXDUR

**Nəyə lazımdır:** Kod reuse və composability — `io.Reader` implement edən istənilən tip (file, buffer, socket, network connection) `io.Reader` qəbul edən istənilən Go funksiyasına ötürülə bilər.

**Üstünlükləri:**
- Standart kitabxanadakı interfeyslər çox kiçikdir (1-2 metod) — "small behaviors"
- Bütün networking kitabxanası `io.Reader` üzərində qurulub: şəbəkə implementasiyası tətbiq məntiqindən ayrılır

### 6. Memory Management (yaddaşın idarəsi)
**Nədir:** Go-nun müasir garbage collector-u (zibil toplayıcısı) yaddaş idarəsini avtomatik aparır.

**Necə işləyir:** C/C++-da yaddaş istifadədən əvvəl allocate, iş bitincə free edilir — səhv edilsə crash və ya memory leak. Go-da GC bu işi özü görür.

**Üstünlükləri:**
- Proqram icra vaxtına əlavə overhead kiçikdir, development səyi isə xeyli azalır

## Hello, Go — ilk proqram

**Kitabdan kod nümunəsi:**
```go
package main

import "fmt"

func main(){
        fmt.Println("Hello World!")
}
```

**Sub-kod izahı:**
- `package main` → Go proqramları paketlər şəklində təşkil olunur; `main` paketi icra olunan proqramı bildirir
- `import "fmt"` → xarici kod istifadəsini təmin edir; `fmt` standart kitabxananın format/çıxış paketidir
- `func main()` → proqram işə düşəndə icra olunan giriş nöqtəsi (C-dəki kimi)
- `fmt.Println(...)` → ekrana sətir çap edir

### Go Playground
**Nədir:** Brauzerdən Go kodu yazıb işə sala biləcəyin onlayn mühit — http://play.golang.org.

**Necə işləyir:** Kod brauzer pəncərəsində redaktə olunur, **Run** düyməsi ilə işə salınır, **Share** ilə paylaşılabilir URL alınır (nümunə: http://play.golang.org/p/EWIXicJdmz).

**Nəyə lazımdır:** Go quraşdırmadan ideya paylaşmaq, nümunə göstərmək, kodu debug etmək. IRC kanallarında, Slack qruplarında, mailing list-lərdə Go developer-ləri Playground linkləri ilə ünsiyyət qurur.

## Əsas terminlər
- Goroutines (qorutinlər — yüngül paralel funksiyalar)
- Channels (kanallar — tipli mesajlaşma strukturları)
- Concurrency (paralellik / eyni-vaxtlılıq)
- Composition (kompozisiya — tiplərin birləşdirilməsi)
- Duck typing (ördək tipi — davranışa əsaslanan tip uyğunluğu)
- Garbage Collector GC (zibil toplayıcı)
- Logical Processor (məntiqi prosessor)
- Static Typing (statik tip sistemi)
- Interface (interfeys — davranış müqaviləsi)

## Praktik nəticə
- Go seçimi "sürətli yaz" vs "sürətli işlə" dilemma-sını aradan qaldırır.
- Server tətbiqlərində concurrency üçün xüsusi thread kitabxanasına ehtiyac yoxdur — `go` açar sözü + channels bəsdir.
- Interfeysləri kiçik saxla (tək metod), böyük inheritance zəncirlərindən qaç — `io.Reader` bunun canlı nümunəsidir.
- Go Playground ilə quraşdırma olmadan Go-nu sına.

## Mənbə
Pages: 23-30 (PDF), book pages 1-8
