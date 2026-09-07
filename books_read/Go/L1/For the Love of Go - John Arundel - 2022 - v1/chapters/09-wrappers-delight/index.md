# Chapter 9 — Wrapper's delight (Sarmalayıcının Zövqü)

## Bu fəsil nədən bəhs edir?

strings.Builder (effektiv string yığımı: WriteString, String, Len), underlying
tip metodları MİRAS OLUNMUR (`type MyBuilder strings.Builder` → mb.Len compile
xətası: "has no field or method Len"), wrapping texnikası — struct daxilində
`Contents strings.Builder` sahəsi (həm underlying metodları `mb.Contents.X`
ilə, həm OWN metodlar), StringUppercaser məşqi (ToUpper + strings.ToUpper),
pass-by-value kəşfi (Double(input int) → `input *= 2` → x DƏYİŞMİR! want 24,
got 12), pointer yaratma (`&x` — "sharing operator"), pointer parametr
tipi (`*int`, distinct tiplər — *int ≠ *float64), dereferencing (`*input`,
"invalid operation: mismatched types *int and int" xətası), nil pointer
(zero value, "panic: invalid memory address or nil pointer dereference"),
pointer metod receiver (`func (input *MyInt) Double()`) — receiver-i dəyişmək
istəyən metod POINTER OLMALIDIR, p.Double() testi.

## Əsas fikirlər

### 1. strings.Builder — Effektiv String Yığımı
Kiçik string-lərdən BÖYÜK string qurma aləti:
```go
var sb strings.Builder          // xüsusi init GEREK YOX
sb.WriteString("Hello, ")
sb.WriteString("Gophers!")

got := sb.String()              // "Hello, Gophers!" — hamısı bir yerdə
gotLen := sb.Len()             // 15 — cari məzmun uzunluğu
```

### 2. Metodlar MİRAS OLUNMUR!
`type MyInt int` fəndi MyInt-ə metod VERİRDİ. Amma metodLU tip əsaslı yeni tip:
```go
type MyBuilder strings.Builder
```
Sual: MyBuilder strings.Builder-in metodlarını ALIRMI?
```go
var mb mytypes.MyBuilder
mb.Len()
// mb.Len undefined (type mytypes.MyBuilder has no field or method Len)
```
**XƏYİR!** Yeni tip = metodlarsız təməl başlanğıc. "That's a shame. But all
is not lost" — LOCAL tip olduğundan ÖZ metodlarımızı əlavə edə BİLİRİK:
```go
func (mb MyBuilder) Hello() string {
    return "Hello, Gophers!"
}
```

### 3. Wrapping — İkisinin də Faydası
İstək: MyBuilder həm OWN metodlu, həm də strings.Builder metodLARINA malik
olsun. **Həll — struct WRAPPING:**
```go
type MyBuilder struct {
    Contents strings.Builder       // sahə = underlying tip
}
```
İstifadə:
```go
var mb mytypes.MyBuilder
mb.Contents.WriteString("Hello, ")     // underlying METODLARI
mb.Contents.WriteString("Gophers!")
got := mb.Contents.String()            // ← strings.Builder metodları İŞLƏYİR
gotLen := mb.Contents.Len()
// + öz metodlarımıza əlavə edə bilərik!
```
**Sirlər heç nə magik deyil:** `mb.Contents` TİPİ strings.Builder-dür → onun
metodları var; MyBuilder LOCAL → own metodlar əlavə oluna bilər. Bu üsulla
METODLU non-local tiplərin HAMSINI genişləndirə bilərsən.

**StringUppercaser məşqi:**
```go
// StringUppercaser wraps strings.Builder.
type StringUppercaser struct {
    Contents strings.Builder
}

func (su StringUppercaser) ToUpper() string {
    return strings.ToUpper(su.Contents.String())   // wrap → convert → qaytar
}
```
Test: `su.Contents.WriteString("Hello, Gophers!")` → `su.ToUpper()` →
"HELLO, GOPHERS!". Bu challenge-i asanlıqla həll etmək = Go tip/metod
qaydalarını bir çox işçi proqramçıdan DƏRİNDƏN başa düşmək.

### 4. Pass By Value — Double Puzzle
```go
func Double(input int) {
    input *= 2          // input-u 2-ə vurur... BURADA
}

x := 12
mytypes.Double(x)
// x hələ də 12! --- FAIL: want 24, got 12
```
Double TAM İSTƏNİLİNİ edir — amma... **Pass by value:** funksiya yalnız
DƏYƏRİN KOPYASINI alır; `input` və `x` = MÜSTƏQİL dəyişənlər. Double öz
LOKAL kopyasını dəyişir, x-a təsir ETMİR.

### 5. Pointer — "Sharing Operator" &
Orijinalı dəyişmək üçün funksiyaya REFERANS lazımdır:
```go
mytypes.Double(&x)         // & = SHARE ET — x-ın pointer-ını ötür
```
Compiler dərhal type mismatch:
```
cannot use &x (type *int) as type int in argument to mytypes.Double
```
İmza güncəllənməli:
```go
func Double(input *int) {    // *int = "pointer to int"
```

### 6. *int ≠ int — və Hətta *int ≠ *float64
Pointer tiplər DISTINCT-dir: *int və *float64 hər ikisi pointerdir, amma
MÜXTƏLİF base tiplərə işarə edirlər → fərqli tiplər, qarışdırmaq OLMAZ.

### 7. Dereferencing — * Operatoru
Pointer-in ÖZÜ ilə riyaziyyat YOX:
```go
input *= 2
// invalid operation: input *= 2 (mismatched types *int and int)
```
2 = int, input = *int → mismatch. Göstərilən DƏYƏR lazımdır — DEREFERENCE:
```go
func Double(input *int) {
    *input *= 2        // *input = "star-input" — pointer-in göstərdiyi dəyər
}
```

### 8. Nil Pointer — "Nil Desperandum"
Pointer tipinin ZERO VALUE = **nil** ("heç nəyə işarə etmir") — error-un
nil-i kimi ("no error" ↔ "no pointee").
**Nil dereference = PANIC:**
```
panic: runtime error: invalid memory address or nil pointer dereference
```
Normal şəraitdə BAŞ VERMƏMƏLİDİR — "panic" ≈ "unrecoverable internal program
error", proqram DAYANIR.

### 9. Pointer Metodlar — Pointer Receiver
Funksiya parametri pointer ola bildiyi kimi, RECEIVER də ola bilər:
```go
func (input *MyInt) Double() {
    *input *= 2
}
```
**QAYDA:** receiver-i DƏYİŞDİRƏCƏK metod = POINTER RECEIVER mütləq!
(Value receiver kopya üzərində işləyərdi — pass-by-value dərsi receiver-ə də
tətbiq olunur.)

Test:
```go
func TestDouble(t *testing.T) {
    t.Parallel()
    x := mytypes.MyInt(12)
    want := mytypes.MyInt(24)
    p := &x          // pointer yarat
    p.Double()       // pointer üzərindən metod çağır
    if want != x {
        t.Errorf("want %d, got %d", want, x)
    }
}
```
Zəncir: &x → p.Double() → *input *= 2 → x = 24. PASS!

## Əsas terminlələr
- strings.Builder — WriteString/String/Len; effektiv string yığımı
- Metod MİRASI YOXDUR — `type X Y` underlying metodları GƏTİRMİR
- Wrapping — struct daxilində underlying tip sahəsi (Contents)
- mb.Contents.X — underlying metodlara çıxış
- StringUppercaser — wrap + ToUpper genişlənmə nümunəsi
- Pass by Value — funksiya DƏYƏRİN KOPYASINI alır; original dəyişmir
- Pointer — dəyişənə REFERANS
- & (Sharing Operator) — &x = x-a pointer yarat
- *int — "pointer to int" tipi
- Distinct Pointer Tipləri — *int ≠ *float64
- Dereferencing — *p = pointer-in göstərdiyi dəyər
- * (Star) Operatoru — pointer → dəyər
- Nil Pointer — pointer zero value; "heç nəyə işarə etmir"
- Nil Pointer Dereference Panic — "invalid memory address"
- Pointer Receiver — `func (p *T) M()` — receiver-i dəyişən metod
- p.Double() — pointer üzərindən metod çağırışı

## Praktik nəticə
(1) Metodlu tipi "genişləndirmək" istəsən: `type X Y` MİRAS VERMİR — struct
WRAP et (`Contents Y`) → hər iki dünyanın metodları. (2) Pass-by-value:
funksiya kopya alır — originalı dəyişmək istəyirsənsə POINTER ötür (&x),
imzada *T qəbul et. (3) *p ilə dereference et — pointer-ilə riyaziyyat
OLMAZ (type mismatch). (4) Pointer zero value = nil; nil dereference = panic
— normal kodda baş vermeməlidir. (5) Receiver-i dəyişən metod = POINTER
RECEIVER (`func (p *MyInt) Double()`); value receiver yalnız OXUYUR. (6)
Pointer tipləri base-ə görə fərqlidir — *int-u *float64 gözləyən yerə ötürmə.
(7) `x *= 2` = `x = x * 2` qısalması (compound assignment). (8) Testlə
düşün, sonra implementasiya "crank the handle" — TDD-nin gücü: ağır düşüncə
testdə, implementasiya mexanikidir.

## Mənbə
Pages: 109-119 (PDF 110-120)
