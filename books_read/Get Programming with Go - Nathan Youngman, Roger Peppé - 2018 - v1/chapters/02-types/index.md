# Unit 2 — Types (Lesson 6-11)

## Bu unit nədən bəhs edir?

Go-nun 15 numeric tipi: float32/64 (dəqiqlik problemi, %f formatı, müqayisə tolerance), 10 integer tipi (wrap-around, hex/binar, uint8 CSS rəngləri, int64 ilə 2038 problemi), big package (big.Int/Float/Rat), untyped constants, mətn (string/rune/byte, UTF-8, raw string, Caesar/ROT13 şifrələri), tiplərarası konvertasiya (strconv, Arianne 5 dərsi). Capstone: Vigenère cipher.

**PDF səhifələr:** 60-106 (L6: 60-67, L7: 68-74, L8: 77-82, L9: 83-93, L10: 94-102, L11: 103-106)

## Əsas fikirlər

### 1. Floating-point əsasları (L6)
**IEEE-754:** real ədədlər ikilik kasa/siklon modeli ilə (bucket + offset bit-ləri). 

**İki tip:**
- `float64` — **default** (kompilyator kəsr ədəd üçün bunu çıxarır; 8 bayt; "double precision")
- `float32` — 4 bayt, dəqiqliyi az; "single precision"; 3D oyun minilliklə vertex-də yaddaş qənaəti üçün

```go
days := 365.2425                    // float64 çıxarır
var pi32 float32 = math.Pi           // 3.1415927 (float32)
var pi64 = math.Pi                   // 3.141592653589793
var answer float64 = 42              // tam ədəd → tip MÜTLƏQ göstər!
```
- `math` paketi float64 üzərində işləyir → **float64-ü üstün tut**
- golint: "should omit type float64" — type inference-i çirkləşdirmə

**Zero value:** `var price float64` → `0` (qiymət "hələ gələcək" mənası; `price := 0.0` = "pulsuz" mənası).

### 2. %f formatı (L6)
```go
third := 1.0 / 3
fmt.Printf("%f\n", third)      // 0.333333 (default 6 rəqəm)
fmt.Printf("%.3f\n", third)    // 0.333   (precision: 3)
fmt.Printf("%4.2f\n", third)   // 0.33     (width: 4)
fmt.Printf("%05.2f\n", third)  // 00.33    (zero-padded)
```
`%4.2f` → width=4 (minimum simvol, sol boşluq), precision=2 (noktadan sonra).

### 3. Floating-point dəqiqliyi (L6)
**Rounding errors — binary təbiət:**
```go
third := 1.0 / 3.0
fmt.Println(third + third + third)   // 1 ✓ (binary ⅓ dəqiqdir!)

piggyBank := 0.1
piggyBank += 0.2
fmt.Println(piggyBank)               // 0.30000000000000004 ✗
// 11 dimes: 1.0999999999999999
```
- **Pul üçün floating-point YARAMAZ** → sent-ləri int-də saxla
- Müqayisə üçün: `piggyBank == 0.3` → **false**! Əvəzinə tolerance:
```go
fmt.Println(math.Abs(piggyBank-0.3) < 0.0001)   // true
```
- Machine epsilon: 2⁻⁵² (float64) / 2⁻²³ (float32) — amma xətalar CƏMLƏNİR → app-a xas tolerance seç
- **Qayda:** vurma bölmədən ƏVVƏL → `(celsius * 9.0 / 5.0) + 32.0` = 69.8 ✓ vs `(9.0/5.0*celsius)+32` = 69.80000000000001 ✗

### 4. Integer tipləri (L7)
**10 tip — signed (5) + unsigned (5):**

| Tip | Range | Yaddaş |
|---|---|---|
| int8 | -128..127 | 1B |
| uint8 | 0..255 | 1B |
| int16 | -32,768..32,767 | 2B |
| uint16 | 0..65,535 | 2B |
| int32 | ±2.1 milyard | 4B |
| uint32 | 0..4.29 milyard | 4B |
| int64 | ±9.2 kvintillion | 8B |
| uint64 | 0..18.4 kvintillion | 8B |
| int / uint | 32 və ya 64-bit (platforma) | — |

- Default inference: tam ədəd → `int`
- **int ≠ int32 alias** — müstəqil tip!
- `%T` → tipi göstər; `%[1]v` → 1-ci arqumenti təkrar istifadə et

### 5. uint8 + hex (L7)
**CSS rəngləri:** RGB 0-255 → mükəmməl uint8:
```go
var red, green, blue uint8 = 0, 141, 213
var red, green, blue uint8 = 0x00, 0x8d, 0xd5   // hex = eyni
fmt.Printf("color: #%02x%02x%02x;", red, green, blue)   // #008dd5;
```
- `0x` prefiksi → hex (A=10...F=15; 2 hex rəqəm = dəqiq 1 bayt)
- `%02x` → 2 rəqəmli, sıfırla doldur

### 6. Wrap-around (L7)
```go
var red uint8 = 255
red++
fmt.Println(red)          // 0 — wrap!

var number int8 = 127
number++
fmt.Println(number)       // -128 — əks tərəfə wrap
```
- `%b` + `%08b` → bit görünüşü: 255 = `11111111` → `00000000` (carry itir)
- `math.MaxUint16` kimi min/max konstantlar mövcuddur

**2038 problemi dərsi:** Unix time (1970-dən saniyə) 2038-də int32-ni aşır → `time.Unix(12622780800, 0)` yalnız **int64** ilə: 2370-cu il ✓. 2 milyarddan böyük + köhnə 32-bit hardware → int64/uint64.

### 7. big package (L8)
**Exponent sintaksisi:** `var distance int64 = 41.3e12` (41.3 trilion; Alpha Centauri).
**uint64 limiti:** 18.4 kvintillion — Andromeda 24 kvintillion km → **overflow**:
```go
var distance uint64 = 24e18    // ERROR: overflows uint64
```

**3 big tipi:** `big.Int` (böyük tam), `big.Float` (arbitrary precision), `big.Rat` (kəsr: ⅓).

**Kitabdan kod nümunəsi (Andromeda):**
```go
package main

import (
    "fmt"
    "math/big"
)

func main() {
    lightSpeed := big.NewInt(299792)      // int64 → big.Int
    secondsPerDay := big.NewInt(86400)
    distance := new(big.Int)
    distance.SetString("24000000000000000000", 10)   // string-dən, base 10
    fmt.Println("Andromeda Galaxy is", distance, "km away.")
    seconds := new(big.Int)
    seconds.Div(distance, lightSpeed)      // Div metodu
    days := new(big.Int)
    days.Div(seconds, secondsPerDay)
    fmt.Println("That is", days, "days of travel at light speed.")
    // 926568346 days
}
```
- `big.NewInt(int64)` → normal ölçüdən; `SetString("...", 10)` → int64-ə sığmayandan
- Aritmetika METODLARLA (`Div`), operatorlarla yox — daha çoxşirik və YAVAŞ

### 8. Untyped constants (L8)
**Sehir:** tipli const daşır, amma **tip verməmiş const** heç bir tipə məhdud deyil:
```go
const distance uint64 = 24000000000000000000   // ERROR: overflows uint64
const distance = 24000000000000000000           // OK! untyped
const lightSpeed = 299792
const secondsPerDay = 86400
const days = distance / lightSpeed / secondsPerDay   // compile-time hesablanır!
fmt.Println("Andromeda Galaxy is", days, "light days away.")   // 926568346
```
- Bütün literal-lar untyped constant-dır; kompilyator arxada **big package** istifadə edir
- Hesablama **compile-time** baş verir
- Nəticə int-ə sığırsa → `km := distance` OK; sığmırsa → overflow
- Amma `days` konstantı funksiyaya ötürüləndə tipə çevrilir; `distance`-i Println-a vermək olmaz (int-ə sığmır)

### 9. String əsasları (L9)
```go
peace := "peace"                    // inference: string
var blank string                    // zero value: ""
fmt.Println(`peace be upon you
upon you be peace`)                 // RAW string: \n işləmir, çoxsətirli!
```
- Raw string literal (backtick) → escape-lər mətn kimi qalır, çoxsətir mümkün — `C:\go` üçün ideal
- String **immutable**: `message[5] = 'd'` → compile xətası

### 10. rune və byte (L9)
- **Unicode code point** → hər simvolun nömrəsi (A=65, 😊=128515)
- `rune` = **int32 alias** — bir Unicode code point
- `byte` = **uint8 alias** — binary data / ASCII (128 simvol)
- Go 1.9+ öz alias-larını yaratmaq: `type byte = uint8`

```go
var pi rune = 960
var alpha rune = 940
var omega rune = 969
var bang byte = 33
fmt.Printf("%c%c%c%c\n", pi, alpha, omega, bang)   // πάω!

grade := 'A'          // character literal — rune inference; dəyəri 65
var star byte = '*'
```

### 11. String-in indekslənməsi + Caesar/ROT13 (L9)
```go
message := "shalom"
c := message[5]           // 'm' — indeks BYTE qaytarır (0-dan)
fmt.Println(len(message)) // 6 — BYTE sayı!
```

**Caesar cipher (shift 3):**
```go
c := 'a'
c = c + 3                       // 'd' — rune riyaziyyat!
if c > 'z' {
    c = c - 26                  // wrap: x→a
}
```

**ROT13 (hər ikisi eyni əməliyyat — şifrələ/dəşifrələ):**
```go
message := "uv vagreangvbany fcnpr fgngvba"
for i := 0; i < len(message); i++ {
    c := message[i]
    if c >= 'a' && c <= 'z' {
        c = c + 13
        if c > 'z' {
            c = c - 26
        }
    }
    fmt.Printf("%c", c)
}
// "hi international space station"
```
- Bu versiya yalnız ASCII üçündür — ispan/rus mətni korlanır!

### 12. UTF-8 decode (L9)
**UTF-8:** variable-length (1-4 bayt/simvol); ASCII ilə üst-üstə düşür; Ken Thompson (Go dizayneri) ixtirası.

```go
question := "¿Cómo estás?"
fmt.Println(len(question), "bytes")                    // 15 bytes
fmt.Println(utf8.RuneCountInString(question), "runes") // 12 runes
c, size := utf8.DecodeRuneInString(question)           // '¿', 2 bytes

for i, c := range question {      // range: UTF-8 decode edir!
    fmt.Printf("%v %c\n", i, c)   // i = bayt indeksi, c = rune
}
for _, c := range question {      // blank identifier: indeks lazım deyil
    fmt.Printf("%c ", c)
}
```
- `len` = bayt; `utf8.RuneCountInString` = simvol sayı
- **range string-i avtomatik rune-lara açır** — çoxdilli ROT13 bununla düzgün işləyir

### 13. Tiplər qarışmır (L10)
```go
countdown := "Launch in T minus " + "10 seconds."    // OK: string+string
countdown := "Launch in T minus " + 10 + " seconds."  // ERROR: mismatched types

age := 41
earthDays := 365.2425
fmt.Println(age*earthDays/...)                        // ERROR: int × float64
```
- JS/PHP: `"10"-1` = 9 (JS/PHP) — Go: compile xətası; `"10"+2` = "102" (JS) / 12 (PHP) — Go: rədd
- **Koercion YOXDUR** — hər konvertasiya EKSPLİSİT

### 14. Numeric konvertasiya + Arianne 5 (L10)
```go
age := 41
marsAge := float64(age)          // int → float64
fmt.Println(int(earthDays))      // 365 — TRUNCATE (yoxlama, round YOX)
```

**Arianne 5 dərsi (1996):** float64 → int16, dəyər 32,767-ni aşdı → Ada-da exception → raket 40 saniyədə partladı. Go-da:
```go
var bh float64 = 32767
var h = int16(bh)    // OK
// bh = 32768 → int16(bh) = -32768 (wrap!) — Go exception VERMİR

if bh < math.MinInt16 || bh > math.MaxInt16 {
    // range yoxlaması — untyped constant müqayisəsi OK
}
```
**Dərs:** hər konvertasiyada nəticənin təhlükəsizliyini DÜŞÜN.

### 15. String konvertasiyası (L10)
```go
// rune → string:
fmt.Print(string(pi), string(alpha), string(omega), string(bang))   // πάω!

// int → string:
str := "Launch in T minus " + strconv.Itoa(countdown) + " seconds."
str := fmt.Sprintf("Launch in T minus %v seconds.", countdown)

// string → int (xəta mümkün!):
countdown, err := strconv.Atoi("10")
if err != nil { /* "10" yoxdursa xəta */ }
```
- `Itoa` = integer to ASCII; `Atoi` = ASCII to integer (err qaytarır — L28-də ətraflı)
- Digits: '0'=48...'9'=57

**Static typing:** dəyişənin tipi elanda BİRLƏŞİR — `countdown = 0.5` (int) → xəta. JS/Python/Ruby dinamikdə tip dəyişir.

### 16. Boolean konvertasiya (L10)
```go
// bool → string:
launchText := fmt.Sprintf("%v", launch)     // "false"
if launch { yesNo = "yes" } else { yesNo = "no" }

// string → bool:
launch := (yesNo == "yes")                  // şərtin nəticəsi birbaşa bool

// bool → int: yalnız if ilə
var oneZero int
if launch { oneZero = 1 } else { oneZero = 0 }
```
- `string(false)`, `int(false)`, `bool(1)` → **compile xətası** — bool-un numeric ekvivalenti YOXDUR

### 17. Capstone: Vigenère cipher (L11)
**Caesar zəifliyi:** frequency analysis (E tez-tez → şifrədə də tez) → asan qırılır.

**Vigenère (16-cı əsr):** hərfə sabit yox, **təkrarlanan keyword** qədər shift (GOLANG → G=6, O=14, L=11...). Caesar = keyword "D"; ROT13 = keyword "N".

Tapşırıq: `CSOITEUIWUIZNSROCNKFD` + `GOLANG` → hər hərf üçün keyword-in müvafiq hərfini çıx ('C'-'G' = 2-6 = -4 → +26 = 22 = 'W'...). 
- `strings.Repeat` ilə keyword-i uzat; modulus ilə wrap (`27 % 26 = 1`; diqqət: `-3 % 26 = -3`!)
- range (rune) vs index (byte) iki variantda yaz

## Əsas terminlər
- IEEE-754 (floating-point standard)
- float64 / float32 (double / single precision)
- Zero Value (sıfır dəyər)
- Machine Epsilon (2⁻⁵² / 2⁻²³)
- Rounding Error (yuvarlama xətası)
- Signed / Unsigned (işarəli / işarəsiz)
- Wrap-around (dövrü taşma)
- Hexadecimal (`0x`), Binary (`%b`), Nibble
- Exponent sintaksisi (`e12`)
- big.Int / big.Float / big.Rat
- Untyped Constant / Literal
- Compile-time Calculation
- String (immutable), Raw String Literal (backtick)
- Unicode / Code Point / UTF-8
- rune (int32 alias) / byte (uint8 alias)
- Character Literal ('A')
- ASCII (128 simvolluq alt çoxluq)
- len (bayt) vs RuneCountInString (simvol)
- Blank Identifier (`_`)
- Type Conversion (explicit)
- Static vs Dynamic Typing
- strconv: Itoa / Atoi
- machine min/max constants (math.MaxInt16...)
- Frequency Analysis
- Vigenère Cipher

## Praktik nəticə
- Pul və digər dəqiq dəyərlər üçün floating-point YOX — sent-ləri int-də saxla.
- Float müqayisəsində `==` QADAĞAN — `math.Abs(a-b) < tolerance` istifadə et.
- Vurma bölmədən əvvəl: `(c * 9 / 5) + 32`.
- 2038 həssaslığı: Unix time / 2 milyarddan böyük dəyərlər → int64.
- int8/int16 konvertasiyasından əvvəl math.Min/Max yoxla — Arianne 5 bu səhvlə partladı; Go exception deyil, wrap verir.
- Çoxdilli mətn üçün range + rune; indeks (`s[i]`) yalnız bayt verir.
- `C:\go` kimi path-lər üçün raw string (backtick).
- Böyük sabitlər üçün tip vermə — untyped const compile-time-da big package ilə hesablanır.

## Mənbə
Pages: 60-106 (PDF), book pages 45-91 (Lesson 6-11)
