# Chapter 14 — Building blocks (Tikinti Blokları)

## Bu fəsil nədən bəhs edir?

Proqramın həyat dövrü: executable binary (machine code, OS formatı), compiler
(mənbə → binary), paylama rahatlığı (TƏK fayl, asılılıqsız — Docker/Kubernetes,
~2 MiB); package main + main() — giriş nöqtəsi (cmd/ konvensiyası); init
funksiyası (main-dən ƏVVƏL, imported paketlərdə də; amma "magical and
non-obvious" — ÜSYƏTƏN MƏSLƏHƏT GÖRÜLMÜR: control flow izləməyi çətinləşdirir),
alternativ — main başlanğıcında et / paket səviyyəsində `var thing =
initialiseThing()` (adi, oxunaqlı Go kodu); go build (binary yaratma, `file`
inspection: Mach-O/PE32+/ELF), cross-compilation (GOOS "goose" + GOARCH
"gorch": `GOOS=windows go build`, `GOOS=linux GOARCH=arm go build`,
`go tool dist list` — hamısı ÖZ maşınından!); exiting (main sonu = avtomatik;
os.Exit(1) — exit status konvensiyası 0=OK/non-zero=xəta; log.Fatal və panic —
dərhal ləğv amma "dərinlərdə basdırılmışsa" oxunuşu çətinləşdirir → "yalnız
main-dən exit et" qaydası — "It brought you into the world... it should be
the one to see you out").

## Əsas fikirlər

### 1. Binary — Proqramın Cismi
Executable binary = machine code + OS-un tələb etdiyi format. Terminaldan adı
yazılaraq YAXUD Finder kimi qrafik interfeysdən başladılır; OS faylı yaddaşa
yükləyir və düzgün nöqtədən icra başlayır.

**Compiler:** "angry messages"-in mənbəyi — Go mənbəsini birbaşa MACHINE
CODE-a çevirir. Payda paylamaq üçün YALNIZ BİR FAYL kifayətdir.

### 2. package main + main()
- Executable üçün: `package main` + `func main()` — giriş nöqtəsi
- Bir qovluqda bir paket → main alt-qovluqda; konvensiya: `cmd/`
- İcra main-in BAŞINDAN başlayır → bildiyimiz control flow → main-in sonuna
  çatanda proqram DAYANIR ("When and if execution reaches the end of main")
- Digər funksiyalar (fmt.Println, öz paketlərin) çağırılır → kontrol geri qayıdır

### 3. init — Magic Funksiya (və NİYƏ YOX)
Fakt: init (varsa) main-dən ƏVVƏL icra olunur; imported paketlərin init-ləri
də hamısı main-dən əvvəl.

**Niyə məsləhət görülmür:** "magical and non-obvious" — paketdə init VARSA
bilməzsən (axtarmandan!); control flow gizlənir; oxunuş çətinləşir. Proqram
axını "as straightforward and obvious as possible" olmalıdır (kitabın leitmotivi).

**Alternativ 1 — main-də başla:** nə lazımdırsa main-in ƏVVƏLİNDƏ et.
**Alternativ 2 — package-level var:**
```go
package stuff

var thing = initialiseThing()    // funksiya çağırışı təyinatda!
```
Paket səviyyəsində var = declaration (ch11-dən: funksiyadan kənarda yalnız
declaration) — amma DƏYƏR də təyin edə bilər, funksiya nəticəsi də! Nəticə:
initialiseThing paketdəki HƏR ŞEYDƏN ƏVVƏL icra olunur — magic YOX, "plain old
Go code that everyone can read".

### 4. go build — Binary Yaratmaq
```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, world!")
}
```
```bash
go build          # hello binary yaranır (output YOX = uğur)
file hello
# hello: Mach-O 64-bit executable x86_64   (macOS nümunəsi)
./hello
# Hello, world!
```
Nə baş verdi: OS binary yüklədi → main-in ilk machine kodundan icra → main
sonu → exit. ~2 MiB tək fayl — Go toolu, kitabxana, HEÇ NƏ tələb etməyən hər
hansı kompüterdə (eyni OS+CPU şərti ilə) işləyər; Docker/Kubernetes — bir fayl.

### 5. Cross-Compilation — "Goose" və "Gorch"
Hər sistem üçün həmin sistemdə compile ETMƏK LAZIM DEYİL! Öz maşınından:
```bash
GOOS=windows go build
# hello: PE32+ executable (console) x86-64 ... for MS Windows

GOOS=linux GOARCH=arm go build
# hello: ELF 32-bit LSB executable, ARM, EABI5 ... statically linked

go tool dist list    # bütün GOOS/GOARCH kombinasiyaları
```
- **GOOS** (tm. "goose" — Go komandasının təsdiyi ilə!) = hədəf OS
- **GOARCH** ("gorch") = hədəf CPU arxitekturası
- Format göstəriciləri: Mach-O (macOS), PE32+ (Windows), ELF (Linux) —
  `file` aləti ilə doğrulanır

### 6. Exiting — Çıxış Yolları
Default: main sonu → avtomatik exit (OS-ə kontrol qayıdır).

**os.Exit(status):**
```go
func main() {
    os.Exit(1)      // exit status 1 → shell görür
}
```
- Exit code = rəqəm; **konvensiya: 0 = hər şey OK; non-zero = xəta**
- Shell: `exit status 1`

**Digər yollar:** log.Fatal(mesaj) və ya panic(...) — dərhal dayandırır
(non-zero status; panic + lokasiya məlumatı).

**Problem:** bunlar kodun DƏRİNLİKLƏRİNDƏ basdırılmışsa — oxuyan hardan/necə
çıxış olacağını GÖRƏ BİLMİR. **Qayda: "Yalnız main-dən exit et."** — "It
brought you into the world, if you like, so it should be the one to see you
out" (dünyaya o gətirdi — çıxarıb da o göndərməlidir!). Buna görə "please don't
panic."

Kitabın finala hazırlığı: son fəsil — kindness, simplicity, humility, not
striving (Tao of Go).

## Əsas terminlələr
- Executable Binary — machine code + OS formatında fayl
- Compiler — mənbə → machine code tərcüməçisi
- Paylama Sadəliyi — tək fayl, asılılıq yox (Docker/K8s uyğun)
- package main / func main() — executable tələbi / giriş nöqtəsi
- cmd/ Konvensiyası — main paketi üçün alt-qovluq
- init — main-dən əvvəl icra olunan magic; POOR STYLE
- Package-Level var — `var thing = initialiseThing()` — init əvəzi
- go build — binary yaratma komandası
- file hello — binary format yoxlaması (Mach-O / PE32+ / ELF)
- ~2 MiB — tipik "Hello world" binary ölçüsü
- GOOS ("goose") — hədəf əməliyyat sistemi dəyişəni
- GOARCH ("gorch") — hədəf CPU arxitekturası
- Cross-Compilation — öz maşından hər platformaya build
- go tool dist list — GOOS/GOARCH cədvəli
- os.Exit(code) — proqramı statusla bitir
- Exit Status Konvensiyası — 0 = OK; non-zero = xəta
- log.Fatal / panic — dərhal ləğv; dərində QADAĞA
- "Yalnız main-dən Exit" — görünən çıxış qaydası
- "It brought you into the world" — main-in rolü metaforası

## Praktik nəticə
(1) CLI proqram: cmd/ qovluğunda package main + main(); library işini ayrı
paketdə saxla. (2) init İSTİFADƏ ETMƏ — eyni məqsəd: main başlanğıcı YAXUD
paket-level `var x = initFn()`. (3) go build → tək asılılıqsız fayl — payda
budur; go run yalnız eksperiment üçün. (4) Başqa platforma lazımdırsa:
GOOS/GOARCH env ilə öz maşında build et — cross-compilation Go-nun super
gücüdür. (5) `go tool dist list` — mövcud hədəflər. (6) Exit status: 0/1
konvensiyasına riayət et; os.Exit YALNIZ main-də. (7) log.Fatal/panic
qatlarında basdırılmış exit-lər oxunuşu öldürür — belə halları refaktorla
main-ə daş. (8) Binary ölçüsü (~2 MiB) — statik linkin qiyməti, amma asılılıq
yoxluğunun faydası. (9) Proqram həyatı: main başlanğıc → control flow → main
sonu — "simple!".

## Mənbə
Pages: 178-184 (PDF 179-185)
