# Chapter 3 — Core Types (səh. 110-133)

## Bu fəsil nədən bəhs edir?

Go-nun əsas tipləri: bool (parol yoxlayıcı nümunəsi), tam ədədlər
(signed/unsigned, int vs int64 fərqi, yaddaş ölçmə), float (dəqiqlik),
overflow/wraparound, math/big, byte, string literal-lar (raw vs
interpreted), rune (multi-byte təhlükəsizliyi) və nil.

## Əsas fikirlər

### 1. bool
**Nədir:** true/false; zero value = false; müqayisələrin nəticə tipi.

```go
fmt.Println(10 > 5)    // true
fmt.Println(10 == 5)   // false
```

**Parol mürəkkəbliyi yoxlayıcısı (Exercise 3.01):**
```go
func passwordChecker(pw string) bool {
    pwR := []rune(pw)              // multi-byte təhlükəsiz!
    if len(pwR) < 8 {
        return false
    }

    hasUpper := false
    hasLower := false
    hasNumber := false
    hasSymbol := false

    for _, v := range pwR {
        if unicode.IsUpper(v) {
            hasUpper = true
        }
        if unicode.IsLower(v) {
            hasLower = true
        }
        if unicode.IsNumber(v) {
            hasNumber = true
        }
        if unicode.IsPunct(v) || unicode.IsSymbol(v) {
            hasSymbol = true
        }
    }
    return hasUpper && hasLower && hasNumber && hasSymbol   // hamısı true olmalıdır
}

func main() {
    if passwordChecker("") {
        fmt.Println("password good")
    } else {
        fmt.Println("password bad")
    }
    if passwordChecker("This!I5A") {
        fmt.Println("password good")   // ✓
    }
}
```
- unicode paketi: IsUpper/IsLower/IsNumber/IsPunct/IsSymbol — hamısı bool

### 2. Tam ədədlər (integers)
**İki ölçü:** işarə (signed/unsigned) + ölçü (8/16/32/64 bit).

| Tip | Diapazon | Ölçü |
|---|---|---|
| int8 | -128..127 | 1 bayt |
| uint8 | 0..255 | 1 bayt |
| int16/uint16 | ±32K / 64K | 2 bayt |
| int32/uint32 | ±2.1K-milyard | 4 bayt |
| int64/uint64 | ±9.2×10¹⁸ / 1.8×10¹⁹ | 8 bayt |
| int/uint | 32 və ya 64 bit — platformdan asılı | — |
| byte | = uint8 alias | 1 bayt |
| rune | = int32 alias | 4 bayt |

**Vacib incəlik:**
- 64-bit sistemdə `int` ≈ `int64` EYNİ ölçüdə, amma FƏRQLİ TİP — birlikdə
  istifadə OLMAZ (32-bit portativliyi üçün)
- Heç bir tam tipi bir-biri ilə qarışdırmaz

**Standard seçim: `int`** — yalnız problem yarananda (yaddaş!) dəyişin.

**Yaddaş ölçmə təcrübəsi (10M element):**
```go
var list []int       // 403 MiB heap
var list []int8      // 54 MiB heap — 8x azalma!
// runtime.ReadMemStats(&m) → m.TotalAlloc
```

### 3. Float nömrələr
```go
var a int = 100
var b float32 = 100
var c float64 = 100
fmt.Println(a / 3)   // 33                  — kəsr ATILIR!
fmt.Println(b / 3)   // 33.333332           — float32 dəqiqliyi
fmt.Println(c / 3)   // 33.333333333333336  — float64 dəqiqliyi

// Geri qaytarma:
fmt.Println((a / 3) * 3)   // 99   — int itki!
fmt.Println((b / 3) * 3)   // 100
fmt.Println((c / 3) * 3)   // 100
```
- Standart seçim: **float64** (yaddaş qənaəti lazımsa float32)
- Maliyyə/bank işlərində xüsusi diqqət tələb olunur

### 4. Overflow və wraparound
```go
var a int8 = 128    // COMPILE XƏTASI — 127 max
```
Compile tutmurssa → RUNTIME wraparound (ən yüksəkdən ən alçağa):
```go
var a int8 = 125      // max 127
var b uint8 = 253     // max 255
for i := 0; i < 5; i++ {
    a++
    b++
    fmt.Println(i, "int8", a, "uint8", b)
}
// a: 126, 127, -128(!), -127, -126   → SİGNEG: max-dan minə
// b: 254, 255, 0(!), 1, 2            → UNSIGNED: 255→0
```
- Diapazonu HƏMİŞƏ nəzərə alın

### 5. Big numbers (math/big)
```go
import ("math"; "math/big")

intA := math.MaxInt64
intA = intA + 1                  // wraparound!

bigA := big.NewInt(math.MaxInt64)
bigA.Add(bigA, big.NewInt(1))    // düzgün artım — big API (biraz künc)
fmt.Println(bigA.String())       // 9223372036854775808
```

### 6. byte
- `byte = uint8` alias; 256 kombinasiya (0..255)
- Şəbəkə/fayl I/O-da hamısı bayt kimi gəlir

### 7. String literal-lar
```go
// RAW (`) — nə yazsan o:
comment1 := `This is the BEST
thing ever!`
comment2 := `This is the BEST\nthing ever!`    // \n HƏRFİMİ hərf kimi

// INTERPRETED (") — transformasiya:
comment3 := "This is the BEST\nthing ever!"    // \n = yeni sətir
```
- Raw-dan istifadə: çox sətirli/"\/" dolu mətnlərdə oxunaqlılıq:
```go
comment1 := `In "Windows" the user directory is "C:\Users\"`
comment2 := "In \"Windows\" the user directory is \"C:\\Users\\\""   // eyni, amma dəhşət
```
- Raw-da ` olmur — interpreted lazımdır
- Literal fərqi YALNIZ yazılışda; dəyişəndə fərq yoxdur

### 8. rune — multi-byte təhlükəsizliyi
**Problem:** string = bayt kolleksiyası; UTF-8 simvollar 4 bayta qədər!

```go
username := "Sir_King_Über"     // 13 hərf, 14 BAYT!

// BAYT üzrə dövr — XƏTALI emal:
for i := 0; i < len(username); i++ {
    fmt.Print(string(username[i]), " ")
    // Ü baytlara parçalanır → "Ã," kimi zibil
}

// RUNE üzrə — düzgün:
runes := []rune(username)
for i := 0; i < len(runes); i++ {
    fmt.Print(string(runes[i]), " ")   // S i r _ K i n g _ Ü b e r ✓
}

// range — dil səviyyəsində rune-addımlı:
logLevel := "デバッグ"
for index, runeVal := range logLevel {
    fmt.Println(index, string(runeVal))
    // index = BAYT ofseti (0, 3, 6), runeVal = simvol
}
```

**len tələsi:**
```go
fmt.Println("Bytes:", len(username))             // 14 — yanlış "uzunluq"
fmt.Println("Runes:", len([]rune(username)))     // 13 — düzgün!

fmt.Println(string(username[:10]))               // Ü yarımçıq kəsilir!
fmt.Println(string([]rune(username)[:10]))       // təmiz kəsim
```
- İstifadəçi girişi yoxlarkən `len(s)` BAYT sayır — rune sayı deyil!
- strings paketini əvvəlcədən yoxlayın

### 9. nil
**Nədir:** "dəyərsizlik" — tipi yoxdur; pointer/map/interface/slice-lar
başlanğıcsız olanda nil olurlar; nil-ə toxunma → crash.

```go
var message []string          // nil
if message == nil {
    fmt.Println("error, unexpected nil value")
    return
}
```

### Activity-lər
- **3.01 Satış vergisi:** məbləğ × vergi dərəcəsi; ümumi cəm
- **3.02 Kredit kalkulyatoru:** creditScore ≥ 450 → 15% (yoxsa 20%),
  ödəniş həddi, validasiyalar (mənfilər, 12-yə bölünmə), hesabat çapı

## Əsas terminlər
- Core types — bool, ədədlər, string — mürəkkəb tiplərin əsası
- Signed/unsigned — mənfi saxlayır / saxlamır
- int vs int64 — eyni ölçü, FƏRQLİ tip
- Zero value — başlanğıcsız default (false/0/"")
- Overflow — tip həddindən aşma (compile tutur)
- Wraparound — max→min sıçrayışı (runtime, gizli bug!)
- math/big — həddsiz ədədlər API
- byte (uint8) / rune (int32) — xammal / simvol
- Raw (`) vs interpreted (") literal
- UTF-8 — 4 bayta qədər simvol kodlaması
- len(s) — bayt sayı (runе yox!)
- range string — bayt-ofset + rune dəyəri
- nil — dəyərsizlik iştirakçısı

## Praktik nəticə
Standartlar: int, float64, interpreted string. Kiçik tip yalnız yaddaş
sübutu ilə (10M int→int8 = 8x). Diapazonları bil — wraparound gizli
çökmə yaradır; həddi aşma ehtimalı varsa math/big. Mətn emalında
həmişə []rune (yaxud range) — bayt-indexli string emalı Ü kimi simvolları
pozur; len(string) bayt sayır! Başlanğıcsız pointer/map/slice = nil —
toxunmazdan əvvəl yoxla.

## Mənbə
Pages: 110-133 (PDF 110-133)
