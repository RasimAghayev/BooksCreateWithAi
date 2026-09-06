# Cheat Sheet — Chapter 3: Primitive Types And Operators

## Integer Declarations

```go
var myInt int                    // Default: 0
var largeInt int64
i := 3                           // Type: int (inference)
u := uint64(4)                   // Explicit type conversion
decInt := 1000                   // Base 10
hexInt := 0x3E8                  // Hexadecimal
octInt := 01750                  // Octal
binInt := 0b1111101000           // Binary
withSep := 1_000                 // Separator for readability
negativeInt := -10
```

**Nə edir:** Müxtəlif ölçülü və formatlı tam ədəd dəyişənləri elan edir.

**Sub-kod izahı:**
- `int` — Architecture-dependent (32/64 bit), default seçim
- `int64` — 64-bit, timestamp, file offset üçün
- `uint64(4)` — Explicit type conversion (açıq tip çevrilməsi)
- `0x` — Hex, `0` — Octal, `0b` — Binary
- `_` — Readability (oxunaqlılıq) üçün separator (ayırıcı)

**Mənbə:** Chapter 3, page 67

---

## Unsigned Integer Wraparound

```go
var u uint64 // 0
u = u - 1
fmt.Println(u) // output: 18446744073709551615
```

**Nə edir:** Unsigned integer-ın minimum dəyərdən aşağı düşdükdə maksimuma qayıtdığını göstərir.

**Mənbə:** Chapter 3, page 67

---

## Float Declarations

```go
var doubleFloat float64           // Default: 0
var singleFloat float32
f := 12.1                         // Type: float64 (inference)
g := float32(12.1)                // Explicit type conversion
floatVal = 12.                    // Decimal point optional
floatVal = 12e0                   // Scientific notation
floatVal = .12e+2                 // Integer part optional
negativeFloat := -12.0
```

**Nə edir:** Müxtəlif formatlarda həqiqi ədəd dəyişənləri yaradır.

**Sub-kod izahı:**
- `float64` — Default, double precision
- `float32` — Single precision, yaddaş/performans tələbi varsa
- `12e0` — `12 × 10^0 = 12`
- `.12e+2` — `0.12 × 10^2 = 12`

**Mənbə:** Chapter 3, page 67

---

## Special Floating-Point Values

```go
f := 2.0
posInf := math.Pow(f, 10_000)    // +Inf
negInf := f - math.Pow(f, 10_000) // -Inf
notANumber := math.Log(-f)        // NaN

fmt.Println(math.IsInf(posInf, 0))  // true
fmt.Println(math.IsNaN(notANumber)) // true
```

**Nə edir:** `+Inf`, `-Inf`, `NaN` xüsusi dəyərlərini yaradır və yoxlayır.

**Mənbə:** Chapter 3, page 67

---

## Complex Number Declarations

```go
cmplxValue1 := complex(1.1, 2.2)    // complex128
cmplxValue2 := complex(float32(1.1), float32(2.2)) // complex64
cmplxValue3 := 1 + 2i               // complex128
cmplxValue4 := -3i                  // equivalent to complex(0, -3)

fmt.Println(real(cmplxValue1), imag(cmplxValue1)) // "1.1 2.2"
```

**Nə edir:** Kompleks ədədləri yaradır və real/imaginary hissələrini ayırır.

**Sub-kod izahı:**
- `complex(real, imag)` — Built-in funksiya
- `1 + 2i` — Complex expression
- `real()`, `imag()` — Komponentləri çıxarmaq

**Mənbə:** Chapter 3, page 67

---

## Arithmetic Operators

```go
a := 7
b := 3
i = a + b    // 10 (sum)
i = a - b    // 4  (difference)
i = a * b    // 21 (product)
i = a / b    // 2  (quotient)
i = a % b    // 1  (remainder/modulo)

u := uint64(1)
u = u + uint64(i) // Type conversion required
```

**Nə edir:** Əsas arifmetik əməliyyatlar — müxtəlif tiplər üçün.

**Mənbə:** Chapter 3, page 67

---

## Bitwise Operators

```go
// AND, OR, XOR, Bit Clear, Shift
a & b    // AND
a | b    // OR
a ^ b    // XOR
a &^ b   // Bit clear (Go-unikal)
a << 1   // Shift left
a >> 1   // Shift right
```

**Nə edir:** Bit səviyyəsində əməliyyatlar.

**Sub-komanda/flag izahı:**
- `&^` — Bit clear: 2-ci operandın 1 bitlərini 0 edir (Go-ya xas)

**Mənbə:** Chapter 3, page 67

---

## Boolean Operations

```go
var boolVar bool // Default: false
boolVar = true
x := 0
y := 1
v := x < y // Type: bool, value: true

func IsEven(i int) bool {
    return i%2 == 0
}
```

**Nə edir:** Məntiqi dəyərlər və boolean funksiyalar.

**Mənbə:** Chapter 3, page 67

---

## Struct Declarations

```go
type person struct {
    name string
    age  int
}

andy := person{
    name: "Andy",
    age:  42,
}
```

**Nə edir:** Composite tip yaradır və instance nümunəsi yaradır.

**Sub-kod izahı:**
- `type person struct` — Yeni tip elanı
- Trailing comma (son vergül) təlab edilir
- Dot-notation ilə sahəyə müraciət: `andy.name`

**Mənbə:** Chapter 3, page 67

---

## Pointer Operations

```go
var intPtr *int              // Pointer elanı, zero value: nil
intValue := 0
intPtr = &intValue           // Address operator
intPtr2 := &intValue         // Short declaration + address

fmt.Println(*intPtr)         // Dereference: 0
*intPtr = 1                  // Modify through pointer
fmt.Println(intValue)        // 1

newIntPtr := new(int)        // new() fundamental types üçün
newPerson := &person{}       // Struct literal + pointer
newPerson = &person{name: "Andy", age: 42}
```

**Nə edir:** Pointer yaradmaq, address (adres) almaq, dereference (izləmə) etmək.

**Sub-komanda/flag izahı:**
- `&` — Address operator (adres operatoru)
- `*` — Dereference (göstəricini izləmə)
- `new(int)` — Fundamental tiplər üçün pointer yaradır

**Mənbə:** Chapter 3, page 67

---

## String Types

```go
basicStr := "Hello, Gophers!"           // Interpreted string
unicodeStr := "你好，地鼠！"              // Unicode
rawString := `Hello\n\u5730\u9F20`       // Raw string (no escapes)
rawStrWithNewlines := `I can
span multiple lines`
```

**Nə edir:** Interpreted və raw string literal-lərini yaradır.

**Sub-komanda/flag izahı:**
- `""` — Interpreted (escape sequence-lər işləyir)
- `` ` `` — Raw (literal, multiline, no escapes)

**Mənbə:** Chapter 3, page 67

---

## String Operations

```go
fmt.Println("hello" == "hello")      // true
fmt.Println("abc" < "cba")           // true (byte comparison)
strHello := "hello"
strHello += " Gophers!"              // Concatenation
```

**Nə edir:** Mətn müqayisəsi və birləşdirmə.

**Mənbə:** Chapter 3, page 67

---

## len() on Strings

```go
asciiCharStr := "easy, right?"
fmt.Println(len(asciiCharStr)) // 12 (bytes)

unicodeCharStr := "地鼠"
fmt.Println(len(unicodeCharStr)) // 6 (bytes, not characters!)
```

**Nə edir:** String-in byte uzunluğunu qaytarır — character sayı deyil.

**Mənbə:** Chapter 3, page 67

---

## Range Loop on String

```go
unicodeCharStr := "地鼠"
for i, r := range unicodeCharStr {
    fmt.Printf("%d:%s ", i, string(r))
}
// output: 0:地 3:鼠
```

**Nə edir:** Unicode sözləri character-character iterasiya edir.

**Sub-komanda/flag izahı:**
- `r` — Rune (Unicode code point)
- `i` — Byte indeksi

**Mənbə:** Chapter 3, page 67

---

## []rune Conversion

```go
unicodeCharStr := "地鼠"
characters := []rune(unicodeCharStr)
for i, r := range characters {
    fmt.Printf("%d:%s ", i, string(r))
}
// output: 0:地 1:鼠
```

**Nə edir:** String-i rune slice-ə çevirir, character indekslərini alır.

**Mənbə:** Chapter 3, page 67

---

## unicode Package

```go
interestingCharacters := "À(cid:4223)"
for _, r := range interestingCharacters {
    fmt.Printf("rune: %s byte length: %d\n", string(r), utf8.RuneLen(r))
    if unicode.IsLetter(r) && unicode.IsUpper(r) {
        fmt.Println("rune is an uppercase letter")
        fmt.Println("lowercase is:", string(unicode.ToLower(r)))
    }
    if unicode.IsSymbol(r) {
        fmt.Println("rune is a symbolic character")
    }
}
```

**Nə edir:** Unicode simvol mə'lumatlarını (tip, böyük/kiçik hərf) yoxlayır.

**Mənbə:** Chapter 3, page 67

---

## fmt.Printf with %x

```go
fmt.Printf("%x\n", "A")       // 41
fmt.Printf("%x\n", "À")       // c380
fmt.Printf("%x\n", "地")      // e1beb8
fmt.Printf("%x\n", "🐰")      // f09f85b0
```

**Nə edir:** String-in byte-lərini hex (onaltılıq) formatda göstərir.

**Mənbə:** Chapter 3, page 67
