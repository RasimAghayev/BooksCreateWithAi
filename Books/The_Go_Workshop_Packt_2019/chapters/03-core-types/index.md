# Chapter 3 — Core Types (Əsas Tiplər)

## Bu fəsil nədən bəhs edir?

Tip sistemi anlayışı (tipin 4 təyinedicisi), bool, rəqəmlər (int/uint sinifləri, float,
overflow/wraparound, math/big), byte, string literal-lar (raw vs interpreted), rune və
multi-byte təhlükəsizliyi, nil və onun yoxlanması; password complexity, satış vergisi
kalkulyatoru və kredit kalkulyatoru activity-ləri.

## Əsas fikirlər

### 1. Tip Nədir — 4 Təyinedici
Tipi təyin edən: saxladığı data növü, icazə verilən əməliyyatlar, əməliyyatların təsiri,
istifadə etdiyi yaddaş. Go GÜCLÜ TİPLİ dildir — tip sabitdir, dəyişə BİLMƏZ.

**Qeyri-intuitiv nümunə:** ən çox istifadə olunan `int` — kompilyasiya hədəfindən asılı
olaraq 32 VƏ YA 64 bit.

### 2. bool
Yalnız `true`/`false`; zero value = `false`. Müqayisə operatorlarının (==, >, <)
nəticəsi bool-dur.

**Password complexity (kitabdan):**
```go
func passwordChecker(pw string) bool {
    pwR := []rune(pw)                     // multi-byte təhlükəsiz
    if len(pwR) < 8 { return false }      // uzunluq rune SAYI ilə
    hasUpper, hasLower, hasNumber, hasSymbol := false, false, false, false
    for _, v := range pwR {
        if unicode.IsUpper(v) { hasUpper = true }
        if unicode.IsLower(v) { hasLower = true }
        if unicode.IsNumber(v) { hasNumber = true }
        if unicode.IsPunct(v) || unicode.IsSymbol(v) { hasSymbol = true }
    }
    return hasUpper && hasLower && hasNumber && hasSymbol   // && zənciri
}
if passwordChecker("This!I5A") { }        // bool → birbaşa şərt
```

### 3. Rəqəmlər — İki Sinif
**Tam ədədlər (signed/unsigned):**
- `int8..int64` (mənfi də ola bilər), `uint8..uint64` (yalnız müsbət)
- `int`/`uint` — platforma asılı (32/64 bit) — MÜTLƏQİYYƏT int64 DEYİL ayrı tipdir!
- BÜTÜN integer tipləri bir-biri ilə birbaşa istifadə OLUNA BİLMƏZ — implicit çevirmə YOX
- **Qayda: həmişə `int` ilə başla** — yalnız memory problemi olanda kiçik tipə keç

**Yaddaş təcrübəsi (kitabdan):**
```go
// 10M element:
var list []int      // TotalAlloc = 403 MiB
var list []int8     // TotalAlloc = 54 MiB  ← ~8x qənaət
```

**Float (float32/float64):** dəqiqlik vs yaddaş balansı; float64 DEFAULT seçim.

**Dəqiqlik təcrübəsi (kitabdan):**
```go
var a int = 100
var b float32 = 100
var c float64 = 100
fmt.Println(a / 3)   // 33                    ← kəsr hissə İTİRİLİR (xətasız!)
fmt.Println(b / 3)   // 33.333332
fmt.Println(c / 3)   // 33.333333333333336   ← ən dəqiq
fmt.Println((a/3)*3) // 99 (!)  (b/3)*3 → 100, (c/3)*3 → 100
```

### 4. Overflow və Wraparound
**Compile-time:** `var a int8 = 128` → overflow XƏTASI (görünür, düzəltilir).
**Runtime wraparound (gizli təhlükə!):**
```go
var a int8 = 125     // max = 127
var b uint8 = 253    // max = 255
for i := 0; i < 5; i++ {
    a++    // 126, 127, -128 (!), -127, -126 — MƏNFİYƏ SARĞIR
    b++    // 254, 255, 0 (!), 1, 2 — SIFIRA SARĞIR
}
```
**Dərs:** dəyişənin maksimal mümkün dəyərini DÜŞÜN — tip ona uyğun seç.

### 5. Böyük Rəqəmlər — math/big
```go
import "math/big"

intA := math.MaxInt64
intA = intA + 1                       // WRAPAROUND — xəta YOX, səhv dəyər!

bigA := big.NewInt(math.MaxInt64)
bigA.Add(bigA, big.NewInt(1))         // DÜZGÜN toplama — API ilə
fmt.Println(bigA.String())            // 9223372036854775808
```
API "qəlizdir" (metod əsaslı) — amma int hədlərini aşanda YEGANƏ yol.

### 6. byte
`byte = uint8` alias. 256 mümkün dəyər (0-255) — istənilən 8-bit pattern. Şəbəkə/fayl
I/O-da hər yerdə.

### 7. String və Literal-lar
**İki literal növü:**
- **Raw** `` `mətn` `` — olduğu kimi (escape YOX; real yeni sətir OK)
- **Interpreted** `"mətn"` — escape emalı (\n yeni sətir, \t tab); daxilində real yeni
  sətir QADAĞA

**Kitabdan müqayisə:**
```go
comment1 := `This is the BEST
thing ever!`                        // raw — literal yeni sətir
comment3 := "This is the BEST\nthing ever!"   // interpreted — \n

comment1 := `In "Windows" the user directory is "C:\Users\"`      // təmiz!
comment2 := "In \"Windows\" the user directory is \"C:\\Users\\\"" // escape cəhənnəmi
```
Raw-da ` ` ` (backquote) OLMAZ. Dəyişənə düşdükdən sonra fərq YOXDUR.

### 8. Rune və Multi-Byte Təhlükəsizliyi
**Problem:** string = BAYT kolleksiyası; UTF-8 simvollar 1-4 bayt.
```go
username := "Sir_King_Über"     // 13 simvol, 14 bayt!
len(username)                   // 14 — BAYT sayı (SİMVOL YOX!)
username[13]                    // 'Ü'-nin yarısı — etibarsız data
string(username[:10])           // multi-byte kəsə bilər — korlanmış mətn
```

**Həll 1 — []rune çevirməsi:**
```go
runes := []rune(username)
len(runes)                          // 13 — simvol sayı
string([]rune(username)[:10])     // təhlükəsiz kəsim
```

**Həll 2 — range (dilə gömülüb):**
```go
for index, runeVal := range "デバッグ" {    // rune-RUNE iterasiya
    fmt.Println(index, string(runeVal))     // index = BAYT ofseti!
}
```
range bayt-deyil, RUNE addımı ilə gedir — multi-byte safe doğma davranış.

**len vs len([]rune) fərqi input validasiyasında bug yaradır** — "8 simvol tələbi"
bayt sayı ilə yoxlanarsa multi-byte simvollu istifadəçi 7 simvollu keçər.

**strings paketi** — əvvəlcə ora bax, hazır alət var.

### 9. nil
**Nədir:** tip DEYİL; "dəyər yoxdur" xüsusi dəyər. Pointer, map, slice, interface üçün.
nil ilə işləmək = CRASH.

```go
var message *string
if message == nil {                    // ƏSAS YOXLAMA PATTERNİ
    fmt.Println("error, unexpected nil value")
    return
}
```

### 10. Activity-lər
**Sales Tax Calculator:** item qiyməti × vergi dərəcəsi; float64 ilə cəmləmə.
**Loan Calculator:** kredit reytinqi ≥450 → 15%, yoxsa 20%; aylıq ödəniş həddi
(good: gəlirin 20%-i, əks halda 10%-i); mənfi dəyər/12-yə bölünməyən müddət → error;
faiz = məbləğ × dərəcə × müddət; Approved şərti bool.

## Əsas terminlələr
- Strongly Typed (güclü tipli) — tip sabit; implicit çevirmə YOX
- Signed/Unsigned — mənfi dəstəkli/dəstəksiz tam ədəd
- Platform-Dependent int — 32/64 bitlik hədəfə görə
- Overflow — tip həddini aşan literal → compile xətası
- Wraparound — runtime-da hədd aşımı → min dəyərə "sarğır"
- math/big — int64-dən böyük ədədlər (big.Int API)
- byte — uint8 alias; raw data vahidi
- Raw/Interpreted String Literal — escape-siz / escape emallı
- Rune — multi-byte UTF-8 simvol üçün tip (int32 əsaslı)
- Byte Offset — range-də indeksin bayt mövqeyi olması
- nil — boş dəyər; pointer/map/slice/interface-lərdə crash təhlükəsi

## Praktik nətidə

(1) int ilə başla; kiçik tip yalnız ölçülmüş memory problemində. (2) int/int64 ayrı
tiplərdir — birbaşa qarışmaz; cast lazımdır. (3) Tam ədəd bölgüsü kəsri SƏSSİCƏ itirir —
dəqiqlik lazımdırsa float64. (4) Runtime wraparound GİZLİDİR — sərhəd dəyərlərini nəzərə
al. (5) Raw literal — çoxsətirli/escape-dolu mətnlərdə oxunaqlılıq. (6) `len(s)` BAYT
sayıdır — simvol üçün `len([]rune(s))` və ya range. (7) String kəsməkdən əvvəl
[]rune-a çevir — multi-byte pozulmasın. (8) range string üçün doğma rune iterasiyasıdır.
(9) strings paketini əvvəlcədən yoxla. (10) nil yoxlaması — pointer/map/interface işlətməzdən
ƏVVƏL; yoxsa panic. (11) MaxInt64+1 big.Int olmadan MÜMKÜN DEYİL.

## Mənbə
Pages: 83-107 (PDF 116-141)
