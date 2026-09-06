# Chapter 3 — Primitive Types And Operators

## Bu chapter nədən bəhs edir?

Bu chapter, Go-nun əsas (primitive) tiplərini və operatorlarını əhatə edir: integer, floating-point, complex, bool, struct, pointer, string, rune və üzəri şriftli (bitwise) operatorlar. Bu tiplər bütün Go proqramlarının əsasını təşkil edir.

## Əsas fikirlər

### 1. Integer Types (Tam Ədəd Tipləri)
**Nədir:** Go-da signed (işarətli) və unsigned (işarətsiz) tam ədəd tipləri, müxtəlif ölçülərdə (8, 16, 32, 64 bit).

**Necə işləyir:** `int`, `int8`, `int16`, `int32`, `int64` — signed. `uint`, `uint8`, `uint16`, `uint32`, `uint64` — unsigned. `byte` (`uint8` alias), `rune` (`int32` alias) xüsusi adlarla. `int` architecture-dependent (müasir platformadan asılı) — 32 və ya 64 bit.

**Nəyə lazımdır:** Əsas hesablamalar, indeksləmə, konfiqurasiya dəyərləri.

**Üstünlükləri:**
- `int` ən çox istifadə edilən, default seçim
- Bitwise (bit) əməliyyatları üçün unsigned tiplər uyğundur

**Çatışmamazlıqları:**
- Unsigned integer "wrap around" (geri dönüş) edə bilər — `0 - 1` maksimum dəyər olur
- Negative dəyər `uint`-ə convert edilərkən qəribə nəticələr

**Kitabdan kod nümunəsi:**
```go
var u uint64 // 0
u = u - 1
fmt.Println(u) // output: 18446744073709551615 (wrap around)
```
**Mənbə:** Chapter 3, pages 67-103

### 2. Floating-Point Types (Həqiqi Ədəd Tipləri)
**Nədir:** Kəsr (onluq) hissəsi olan ədədlər — `float32` (single precision) və `float64` (double precision). IEEE-754 standartına uyğun.

**Necə işləyir:** Default olaraq `float64` təyin olunur. `12.1`, `.1`, `1e3` kimi literal (sabit dəyər) sintaksisi var. `math.IsInf`, `math.IsNaN` ilə xüsusi dəyərlər (`+Inf`, `-Inf`, `NaN`) yoxlanılır.

**Nəyə lazımdır:** Elmi hesablamalar, faiz dəyərləri, statistik məlumatlar.

**Üstünlükləri:**
- `float64` daha yüksək precision (dəqiqlik) — əsas seçim
- `float32` yalnız yaddaş/performans tələbi varsa

**Çatışmamazlıqları:**
- Dəqiqlik itkisi (precision loss) ola bilər
- `+Inf`, `-Inf`, `NaN` xüsusi halları yoxlanmalıdır

**Kitabdan kod nümunəsi:**
```go
f := 2.0
posInf := math.Pow(f, 10_000) // +Inf
fmt.Println(math.IsInf(posInf, 0)) // true
```
**Mənbə:** Chapter 3, pages 67-103

### 3. Complex Number Types (Kompleks Ədəd Tipləri)
**Nədir:** Real və imagery (xəyali) hissələrindən ibarət tiplər — `complex64` (float32 əsasında) və `complex128` (float64 əsasında).

**Necə işləyir:** `complex(1.1, 2.2)` ilə yaradılır. `1.1 + 2.2i` sintaksisi də mümkündür. `real()` və `imag()` funksiyaları ilə komponentlər çıxarılır.

**Nəyə lazımdır:** Fizika, signal processing (signal emalı), müxtəlif elmi sahələr.

**Üstünlükləri:**
- Daxili tip — xarici kitabxana tələb etmir
- `complex128` dəqiqliyə görə üstündür

**Çatışmamazlıqları:**
- Nadir istifadə — əksər proqramçılar ilk dəfə görəcək

**Mənbə:** Chapter 3, pages 67-103

### 4. Mathematical Operators (Riyazi Operatorlar)
**Nədir:** Əsas arifmetik (`+`, `-`, `*`, `/`, `%`) və bitwise (`&`, `|`, `^`, `&^`, `<<`, `>>`) operatorlar.

**Necə işləyir:** Eyni tip dəyərlər üzərində işləyir. Bölmə `0` ilə integer-də panic, float-də `+Inf` verir. Type conversion (tip çevrilməsi) fərqli tiplər üçün lazımdır.

**Nəyə lazımdır:** Hər hansı hesablama, bit manipulyasiyası, şifrləmə.

**Üstünlükləri:**
- Tam tiplər üçün bütün standart operatorlar mövcuddur
- `&^` (bit clear) Go-ya xasdır — 2-ci operandın 1 bitlərini 0 edir

**Çatışmamazlıqları:**
- Fərqli tiplər arasında birbaşa əməliyyat mümkün deyil — convert lazımdır

**Kitabdan kod nümunəsi:**
```go
u := uint64(1)
i := int(2)
u = u + uint64(i) // Type conversion required
```
**Mənbə:** Chapter 3, pages 67-103

### 5. bool Type
**Nədir:** `true` və ya `false` dəyərlərini saxlayan tip. Conditional (şərtli) məntiq üçün əsasdır.

**Necə işləyir:** `if`, `for` şərtlərində, funksiya geri qaytarması (return) kimi yerlərdə işləyir. Boolean expression (ifadə) `bool` tip qaytarır. Zero value (sıfır dəyəri) `false`-dir.

**Nəyə lazımdır:** Şərtli idarəetmə, doğrulama yoxlamaları, filter-ləşmə.

**Üstünlükləri:**
- Sadə və aydın — `true`/`false`
- Helper funksiyalar ilə mürəkkəb yoxlamaları sadələşdirmək

**Çatışmamazlıqları:**
- Bitwise operatorlər `bool` üçün mövcud deyil

**Kitabdan kod nümunəsi:**
```go
func IsEven(i int) bool {
    return i%2 == 0
}

if IsEven(num) {
    fmt.Printf("%d is even.\n", num)
}
```
**Mənbə:** Chapter 3, pages 67-103

### 6. Struct Types
**Nədir:** Müxtəlif tipli sahələri bir araya gətirən composite (birləşik) tip.

**Necə işləyir:** `type Person struct { name string; age int }` ilə yaradılır. Sahələr dot-notation (`p.name`) ilə oxunur/yazılır. Struct literal (`Person{name: "Andy", age: 42}`) ilə initialize olunur. Trailing comma (son vergül) təlab edilir.

**Nəyə lazımdır:** Əlaqəli məlumatları qruplaşdırmaq, data modeling (məlumat modeli).

**Üstünlükləri:**
- Təkrar istifadə edilən, mə'nalı tiplər yaratmaq
- Composition (birləşmə) üçün əsas

**Çatışmamazlıqları:**
- Trailing comma unutmaq asan xəta
- Inheritance (miras) yoxdur — composition istifadə edin

**Mənbə:** Chapter 3, pages 67-103

### 7. Pointer Types
**Nədir:** Yaddaş adresini saxlayan tiplər — `*int`, `*Person` və s. Dəyəri kopyalamadan bölüşmək üçün istifadə olunur.

**Necə işləyir:** `&variable` — adres operatoru, `*pointer` — dereference (göstəricini izləmə). `new(Type)` ilə pointer yaradılıb sıfır dəyəri verilir. `nil` — işarətsiz pointer.

**Nəyə lazımdır:** Böyük data strukturlarını kopyalamadan funksiyalara ötürmək, shared state (bölüşülmüş vəziyyət) idarəetmək.

**Üstünlükləri:**
- Memory efficiency (yaddaş səmərəliliyi) — kopya yox
- Slice və map-lər daxilində pointer əsasında işləyir

**Çatışmamazlıqları:**
- Nil pointer dereference — panic
- Dangling pointer riski (azmiş göstərici)

**Kitabdan kod nümunəsi:**
```go
var intPtr *int
intValue := 0
intPtr = &intValue
fmt.Println(*intPtr) // 0
*intPtr = 1
fmt.Println(intValue) // 1
```
**Mənbə:** Chapter 3, pages 67-103

### 8. String Type (Mətn Tipi)
**Nədir:** Read-only (yalnız oxunan) byte array (bayt massivi) — UTF-8 mətni saxlamaq üçün nəzərdə tutulmuşdur.

**Necə işləyir:** Interpreted string (tərcümə edilmiş mətn) `"Hello"` — escape sequence-ləri (`\n`, `\u5730`) izah edir. Raw string (işarətli mətn) `` `Hello\n` `` — heç bir interpretasiya yox, multilinedir. Strings are immutable (dəyişməz) — dəyişiklik yeni string yaradır.

**Nəyə lazımdır:** Mətn emalı, istifadəçi mesajları, API responses (cavablar).

**Üstünlükləri:**
- UTF-8 dəstəyi — bütün dillər
- `+` ilə concatenation (birləşmə)
- `strings` paketi ilə zəngin funksiyalar

**Çatışmamazlıqları:**
- `len()` byte qaytarır, character (simvol) sayı deyil
- Unicode üçün `range` və ya `[]rune` istifadə edin

**Mənbə:** Chapter 3, pages 67-103

### 9. Rune Type (Simvol Tipi)
**Nədir:** Unicode character (simvol) təmsilçisi — `int32` alias. `byte` (`uint8`) ilə fərqli olaraq, rune Unicode code point (kod nöqtəsi) saxlayır.

**Necə işləyir:** `for i, r := range str` — `r` rune tipindədir. `unicode.IsLetter(r)`, `unicode.ToLower(r)` ilə simvol mə'lumatları alınır. `[]rune(string)` ilə string → rune slice çevrilir.

**Nəyə lazımdır:** Unicode emalı, karakter sayma, case conversion (böyük/kiçik hərf).

**Üstünlükləri:**
- Çoxdilli mətnlərlə asan işləmə
- `utf8` paketi ilə byte uzunluğunu ölçmək

**Çatışmamazlıqları:**
- `len()` byte qaytarır, `len([]rune)` character sayı

**Mənbə:** Chapter 3, pages 67-103

### 10. String Concatenation və Manipulyasiya
**Nədir:** Mətnləri birləşdirmək və emal etmək üsulları.

**Necə işləyir:** `+` operatoru ilə birləşmə. `strings.ToLower`, `strings.Replace`, `strings.Fields` kimi standart funksiyalar. `fmt.Sprintf` formatlı birləşmə.

**Nəyə lazımdır:** Mətn formatlaşdırma, axtarış, çıxış hazırlamaq.

**Üstünlükləri:**
- `strings` paketi geniş funksiya dəsti təklif edir
- Immutability (dəyişməzlik) — thread-safe (thread təhlükəsiz)

**Çatışmamazlıqları:**
- `+` ilə çoxlu birləşmə performansı aza bilər — `strings.Builder` daha yaxşıdır

**Mənbə:** Chapter 3, pages 67-103

## Əsas terminlər
- Integer (tam ədəd) — signed/unsigned, architecture-dependent/independent
- Floating-point (həqiqi ədəd) — `float32`, `float64`
- Complex (kompleks ədəd) — `complex64`, `complex128`
- bool (məntiqi tip) — `true`/`false`
- Struct (strukt) — Composite tip, sahələr qrupu
- Pointer (göstərici) — `*Type`, yaddaş adresi
- String (mətn) — Read-only byte array, UTF-8
- Rune (simvol) — `int32`, Unicode code point
- Bitwise operator (bit əməliyyatı) — `&`, `|`, `^`, `&^`, `<<`, `>>`
- Zero value (sıfır dəyəri) — Tipin default dəyəri

## Praktik nəticə
Primitive tiplər hər Go proqramının əsasını təşkil edir. `int` seçmək ilk addım, `float64` kəsr üçün, struct və pointer data strukturları üçün əsasdır. String-lər UTF-8 əsasında işləyir, `len()` byte qaytarır — character sayı üçün `range` və ya `[]rune` istifadə edin.

## Mənbə
Pages: 67-103
