# Chapter 1 — Why Go? (Niyə Go?)

## Bu chapter nədən bəhs edir?
Go-nun sistem proqramlaşdırması üçün niyə uyğun dil olduğuna: seçim arqumentləri, CSP-əsaslı paralellik modeli, OS ilə aşağı səviyyəli qarşılıqlı əlaqəyə, builtin alətlər dəstinə və cross-platform inkişafa.

## Əsas fikirlər

### 1. Sistem proqramlaşdırması nədir
**Nədir:** Aşağı səviyyəli tapşırıqlar mərkəzli proqramlaşdırma: fayl/kataloq yaradılması, proseslərin idarəsi, digər proqramların icrası, thread/proses kommunikasiyası, şəbəkə üzərində IPC.

**Müəllifin mövqeyi:** Sistem proqramlaşdırması cansıxın deyil — sehrbazlıqdır: OS və avadanlığı idarə edirsən, başqa dillərdə mümkün olmayan şeyləri edirsən.

### 2. Go-nu seçmək — rəqiblərə qarşı
**Rəqib spektiri:** C/C++ (köklü), Zig/Rust/Odin (yeni dalğa).

**Rəqib tələləri:** dik öyrənmə əyrisi, yüksək koqnitiv yük, cəmiyyət/destək çatışmazlığı, API qeyri-tutarlılığı, adoptasiya azlığı.

**Go-nun dizayn fəlsəfəsi:** sadəlik, ifadəlilik, möhkəmlik, effektivlik. Təcrübəli proqramçılar ~2 həftədə Go-ya adaptasiya olunur (ekspert olmasa da, orta mürəkkəbli kod yazırlar). Python/Ruby proqramçıları ifadəliliyi itirmədən performans + paralelliyi qazanırlar.

**GC tənqidi və cavab:** GC pauzaları ən pis halda **<100 mikrosaniyə**; Go 1.20+-dan daha granulyar yaddaş idarəsi (arenas — Ch 8).

**Fəlsəfi balans:** Go CPU üçün zero-cost hədəfləmir — proqramçının səyini azaltmaq prioritetdir (bunun yan məhsulu: prosesin dadsından gələn həzz).

### 3. Paralellik (Concurrency)
**Nədir:** Çox tapşırığın eyni anda icrası — real-zamanlı sistemlər üçün kritik (millisaniyələrin əhəmiyyətli olduğu ssenarilər).

**Faydaları:**
- Throughput artımı (vahid vaxtda emal olunan məlumat miqdarı)
- Tapşırıq tamamlanma vaxtı azalır
- Cavabdeqlik/responsivlik artır
- İzolyasiya qabiliyyəti — data bütövlüyü qoruması
- CPU-bound + I/O-bound tapşırıqların orkestrasiyası — CPU işləyərkən I/O gözləyir

### 4. Goroutinlər
**Nədir:** Yüngül icra threadləri — "green threads" (yaşıl axınlar). Yaradılması ucuz, minlərlə eyni anda bir neçə OS thread üzərində işləyə bilir.

### 5. CSP-ilhamlı model
**Nədir:** Communicating Sequential Processes (Hoare, 1978) — riyazi formalizm; kommunikasiya + sinxronizasiya kanallar vasitəsilə.

**Go-nun tətbiqi:**
- Kanallar = CSP kommunikasiyasının birbaşa reallaşdırılması
- Shared memory + lock-lar əvəzinə → azaldılmış mürəkkəblik, azaldılmış risk
- Paralellik DİLİN ÖZÜNƏ daxilidir (xarici kitabxana deyil) → daha az səhvlə meyilli kod
- Callback hell yoxdur — prosedural oxunuş, adi funksiya imzaları

### 6. "Share by communication"
**Məşhur aforizm:** "Don't communicate by sharing memory, share memory by communicating"

**Müəllifin dəqiqləşdirməsi:** "Share by communicating, not by locking" — kanallarla mesaj mübadiləsi, açıq lock-lara ehtiyacı azaldır.

**Funksional proqramlaşdırma bonusu:** Go-da data goroutinlər arasında **implicit paylaşımı YOXDUR — kopyalanır** → data dəyişməzliyi (immutability) problemi avtomatik həll olunur.

### 7. OS ilə qarşılıqlı əlaqə
Go-nun syscall yanaşması: təhlükəsiz və effektiv (paralellik modeli ilə uyğun). Digər dillərlə müqayisədə nisbətən aşağı səviyyəli → dəqiq nəzarət, amma OS API bilik tələbi (öyrənmə əyrisi). Kitabın 2-ci hissəsində dərinləşəcək.

### 8. Tooling — builtin alətlər
| Alət | Vəzifə | Nümunə |
|------|--------|--------|
| `go build` | compile → executable | `go build main.go` → `./main` |
| `go test` | testləri icra et | `go test` (math_test.go avto-tapılır) |
| `go run` | compile-sız birbaşa icra | `go run hello.go` |
| `go vet` | şübhəli konstruktları yoxla | `Printf format %s has arg 1999 of wrong type int` |
| `go fmt` | standart formatlama | `msg:="Hello"` → `msg := "Hello"` |

### 9. Cross-platform inkişaf
**GOOS + GOARCH environment dəyişənləri:**
```bash
GOOS=linux GOARCH=amd64 go build    # Linux üçün
GOOS=darwin GOARCH=amd64 go run     # macOS üçün
```

**Build tags ilə kod seqmentasiyası:**
```go
// go:build windows
package main
import "fmt"
func main() { fmt.Println("This is Windows!") }
```
```go
// go:build linux
package main
import "fmt"
func main() { fmt.Println("This is Linux!") }
```
```bash
GOOS=windows go build -o app.exe   # Windows binary
GOOS=linux go build -o app         # Linux binary
```

**Xəbərdarlıq:** Syscall-larla sıx işə düşən kod OS-özəlliyinə bağlı olur → şərti kompilyasiya (build tags) və ya platforma-özəl düzəlişlər tələb edə bilər (xüsusən golang.org/x/sys ilə işləyəndə).

## Əsas terminlər
- System programming (sistem proqramlaşdırması) — OS/hardware yaxın kod
- Green threads (yaşıl axınlar) — runtime idarəli yüngül threadlər
- CSP (Communicating Sequential Processes) — Hoare-in formalizmi
- Share by communication — kanalla paylaşma fəlsəfəsi
- Callback hell / pyramid of doom — nested callback antipattern-i
- GC pause / stop-the-world — zibil toplayıcının pauzası (<100µs)
- Throughput — vahid vaxtda emal miqdarı
- go build/test/run/vet/fmt — builtin alət dəsti
- GOOS / GOARCH — platforma hədəf dəyişənləri
- Build tags — `// go:build <os>` şərti kompilyasiya

## Praktik nəticə
1. Sistem aləti yazarkən Go seçimi: C/Rust sürəti YOX, amma 2 həftəlik öyrənmə + GC <100µs + native paralellik — praktik balans.
2. Paralellik dizaynında kanallardan başla; lock-lar yalnız performans tələbi olanda.
3. `go vet` + `go fmt` CI-in minimum standartı — heuristikalar səhvləri compiler-dən əvvəl tutur.
4. Platforma-özəl kod üçün build tags + fayl adı konvensiyası (`main_windows.go`).
5. GC narahatlığı Go 1.20+ üçün keçərli deyil — arenas (Ch 8) granulyar idarə verir.

## Mənbə
Pages: 3-12 (PDF səh. 24-33)
