# Chapter 9 — Basic Debugging (Əsas Debuginq)

## Bu fəsil nədən bəhs edir?

Bug-lərin səbəbləri və profilaktika (incremental test, unit test, error handling,
logging), fmt formatlaşdırma (Println, Printf, verb-lər, genişlik/dəqiqlik/align),
debug metodları (kod markerləri, %T, %#v), log paketi (timestamp-lar, SetFlags,
Ldate/Lmicroseconds/Llongfile), Fatal vs Panic və SSN validasiyası activity.

## Əsas fikirlər

### 1. Bug Səbəbləri (Production-da)
1. **Test yalnız sonda** — incremental test YOX; funksiya bitən kimi sına
2. **Enhancement/dəyişikliklər** — bir sahənin dəyişməsi digərinə təsir; unit test
   qoruma verir
3. **Real olmayan timeframe** — qısa müddət → shortcut, az test, qeyri-aydın tələb
4. **Error-ların işlənməməsi** — fayl tapılmadı, /0, bağlantı qurulmadı → handle ET

**Ciddilik:** heparin dərman dozası maşını bug-u həyat üçün təhlükəlidir —
bug-free kod kritikdir.

### 2. Profilaktika — 4 Metod
1. **Kod incremental + çox test:** hər parça bitəndə test → bug İZLƏMƏ sadələşir
2. **Unit test:** input→nəticə yoxlaması; dəyişiklikdən SONRA fail = dəyişiklik
   bug gətirdi; push-dan əvvəl KEÇMƏLİ
3. **BÜTÜN error-ları handle et** (Ch6)
4. **Logging:** debug/info/warn/error/fatal/trace; dəyişən dəyərləri, icra yeri,
   arqumentlər, nəticələr; performance təsirini nəzərə al

### 3. fmt.Printf — Verb Formatlaşdırma
**Println vs Printf:** Println = avtomatik space + newline, default format;
Printf = verb placeholder-lər, newline AVTOMATİK DEYİL (\n əl ilə).

**Əsas verb-lər:**
| Verb | Tip |
|---|---|
| %s | string |
| %d | int (base-10) |
| %f | float (default 6 onluq) |
| %t | bool |
| %v | default value format |
| %T | tipin ADI (BÖYÜK T!) |
| %#v | Go sintaksis repr |

**Genişlik/dəqiqlik/align:**
```go
fmt.Printf("%.2f\n", 3.75)        // 3.75 — dəqiqlik
fmt.Printf("%10.2f\n", 1234.67)   // "     1234.67" — width 10, SAĞA align
fmt.Printf("%-10.2f\n", v)         // SOLA align (- flag)
fmt.Printf("%8.b %2.x\n", i, i)   // binary/hex verb-ləri də formatlanır
```

### 4. Debug Metodları — 4 Alət
**1) Kod markerləri** — "haradayam" ifadələri:
```go
fmt.Println("We are in function calculateGPA")
```
Yerləşdir → çatmırsa bug daha YUXARIDA; çatırsa markeri başqa yerə köçür —
skopu DARALT.

**Kitabın random nümunəsi:** eyni error iki funksiyadan gəlir (a: i<10, b: i>=10) —
hər funksiyada öz markeri → error HARADAN gəldiyini dərhal göstərir.

**2) Dəyişən TİPİ:**
```go
fmt.Printf("fname is of type %T\n", fname)     // string
fmt.Printf("grades is of type %T\n", grades)   // []int
fmt.Printf("states is of type %T\n", states)  // map[string]string
fmt.Printf("p is of type %T\n", p)             // main.person
```
Böyük %T = tip, kiçik %t = bool — case-sensitive!

**3) Dəyişən dəyəri — %#v (Go repr):**
```go
fmt.Printf("fname value %#v\n", fname)   // "Joe" — kodu kopyalaya biləcəyin formada
fmt.Printf("p value %#v\n", p)            // main.person{lname:"Lincoln", age:210, ...}
```
Go-un datanı GÖRDÜYÜ şəkildə — struct/slice/map-in dəqiq strukturu.

**4) Debug log:**
```go
log.Printf("fname value %#v\n", fname)
```

### 5. log Paketi
**fmt-dən fərqi:** hər sətirə AVTOMATİK tarix+vaxt:
```
2019/11/10 23:00:00 Demo app
2019/11/10 23:00:00 Thanos is here!
```

**SetFlags — daha çox detal:**
```go
log.SetFlags(log.Ldate | log.Lmicroseconds | log.Llongfile)
log.Println("Demo app")
// 2019/11/10 23:00:00.000001 /path/main.go:14: Demo app
```
- Ldate — tarix; Lmicroseconds — mikro-saniyə; Llongfile — TAM fayl yolu + sətir
- Başqa flag-lər: Ltime, Lshortfile, LUTC... (log paketinin sənədi)

**Niyə vacib:** sətir nömrəsi + mikro-saniyə — hadisə ardıcıllığını və YERİNİ dəqiq
müəyyən edir; production-da debug-in əsas aləti.

### 6. Fatal vs Panic (log-da)
| | log.Fatal(n) | log.Panic(n) |
|---|---|---|
| Çıxış | os.Exit(1) DƏRHAL | panic başlanır |
| Defer-lər | İŞLƏMİR (os.Exit onları öldürür) | İŞLƏYİR (recover mümkün) |
| Bərpa | MÜMKÜNSÜZ | recover() ilə mümkün |

**Fatal nə vaxt:** data korlanması riski / bərpa mənasız / exit code vermək lazımdır
(CLI utility-lər).

**Kitabdan nümunə:**
```go
log.Println("Start of our app")
err := errors.New("Application Aborted!")
if err != nil {
    log.Fatalln(err)          // log + os.Exit(1)
}
log.Println("End of our app")   // İCRA OLUNMUR
```

### 7. Logging Fəlsəfəsi (kitabdan)
- Logging İNFRASTRUKTURDUR — bug olanda YOX, həmişə olmalıdır
- OS-lər daim loglayır (health, resource access) — eyni yanaşma
- Production şərtləri dev-dən FƏRQLİDİR (yüksək yük, malform data) — logsuz
  səbəbi tapmaq SAATSIZ işdir
- Performance: pik yükdə extensive log proqramı yavaşladır — balans saxla

### 8. SSN Activity — 4 Custom Error + İz Trace
```go
var (
    ErrInvalidSSNLength   = errors.New("invalid SSN length")
    ErrInvalidSSNNumbers  = errors.New("SSN must be all numbers")
    ErrInvalidSSNPrefix   = errors.New("SSN prefix cannot be 000")
    ErrInvalidDigitPlace  = errors.New("SSN starting with 9 requires 7 or 9 in 4th place")
)
// 4 validasiya funksiyası — hər biri SSN + error qaytarır
// main: slice üzrə loop → xəta olsa LOG ET + DAVAM ET (Fatal YOX!)
for _, ssn := range validateSSN {
    if err := validateSSN(ssn); err != nil {
        log.Printf("%v: %v\n", ssn, err)    // log + continue
    }
}
```
**Dərs:** validasiya pipeline-da bir element xətası BÜTÜN proqramı dayandırmamalı —
log et, növbətinə keç.

## Əsas terminlələr
- Bug/Debugging — gözlənilməz davranış / səbəbin tapılması
- Incremental Testing — hər funksiya bitəndə test
- Unit Test — input→nəticə müqaviləsi; regressiya qoruyucusu
- Format Verb — %s/%d/%f/%t/%v/%T/%#v placeholder-lər
- Width/Precision — %10.2f genişlik.dəqiqlik
- Left/Right Align — %-10.2f sol / default sağ
- Code Marker — "haradayam" çap ifadəsi
- %#v — Go sintaksis təmsili (kopyalanabilir)
- %T — concrete tipin adı
- log.Printf — timestamp-li çap
- log.SetFlags — Ldate/Lmicroseconds/Llongfile detal qatı
- log.Fatal — çap + os.Exit(1); defer-lər işləmir
- log.Panic — çap + panic; recover mümkün
- Debug Logging — proqram vəziyyətinin izlənməsi

## Praktik nətidə

(1) Test SONDAN YOX, incremental — hər funksiya bitən kimi yoxla. (2) Unit test
dəyişiklikdən sonra fail edirsə — dəyişikliyin bug gətirdiyini dərhal bilirssən.
(3) Bug axtarışına MARKER ilə başla — skopu addım-addım daralt. (4) %T tip
qarışıqlığını, %#v data strukturunu göstərir — hər ikisi debug arsenalıdır.
(5) %f default 6 onluq verir — %.2f ilə qısalt. (6) log > fmt debug üçün:
timestamp + sətir nömrəsi AVTOMATİK. (7) SetFlags(Llongfile) — log-un HANSİ
sətirdən gəldiyini göstərir. (8) Fatal: bərpa mənasızsa; amma defer-lər ölür —
resurs təmizliyi lazımdsa PANIC+recover seç. (9) Logging həmişə aktiv — production
fərqli şərtlərdə işləyir; logsuz orada baş verəni tapmaq mümkünsüzdür. (10)
Validasiya zəncirində element xətası — log + continue, proqramı dayandırma.

## Mənbə
Pages: 323-350 (PDF 356-385)
