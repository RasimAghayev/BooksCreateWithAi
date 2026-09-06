# Chapter 2 — Primitive Types and Declarations (Primitiv Tiplər və Bəyanatlar)

## Bu chapter nədən bəhs edir?

Go-nun daxili (built-in) tipləri — boolean, tam ədədlər, üzən nöqtə ədədləri, kompleks
ədədlər, string və rune; literal formaları; açıq tip çevrilməsi (explicit type conversion);
`var` və `:=` bəyanat üsulları; `const` (typed/untyped); istifadə olunmayan dəyişənlər
qaydası və adlandırma konvensiyaları.

## Əsas fikirlər

### 1. Zero Value (Sıfır dəyəri)
**Nədir:** Go-da dəyər təyin edilmədən bəyan olunan hər dəyişənin avtomatik default dəyəri.

**Necə işləyir:** Hər tipin öz zero value-su var: rəqəmsal tiplər üçün `0`, `bool` üçün
`false`, string üçün `""`. C/C++-dakı kimi təyin edilməmiş yaddaş (uninitialized memory)
problemi yoxdur.

**Nəyə lazımdır:** Kodu oxunaqlı edir və uninitialized dəyişən bug-larını aradan qaldırır.

**Kitabdan kod nümunəsi:**
```go
var x int     // 0
var flag bool // false
var s string  // ""
```

### 2. Literals (Literal-lar)
**Nədir:** Koddada birbaşa yazılan ədəd, simvol və ya sətir dəyərləri.

**Necə işləyir:**
- Integer: `0b` (binari), `0o` (oktal), `0x` (heksadecimal), `0` (köhnə oktal —
  QARIŞDIRICI sayılır, istifadə ETMƏ). Oxunaqlıq üçün `_` ayırıcı: `1_234`.
- Float: onluq nöqtə + `e` eksponent (`6.03e23`), `0x` + `p` heksadecimal forma.
- Rune: tək dırnaq — `'a'`, `'\141'`, `'\x61'`, `'\u0061'`, `'\U00000061'`; escape-lər:
  `'\n'`, `'\t'`, `'\''`, `'\"'`, `'\\'`.
- String: ikili dırnaq = interpreted literal (`"...\n..."`); backquote = raw literal
  (`` `...` ``) — backquote xaric hər simvolu literal saxlayır, çoxsətirli mətn üçün ideal.
- **Literal-lar typedır** — uyğun gələn istənilən tiplə dəyişənə təyin edilə bilər;
  amma string literal-i rəqəmə (və ya əksinə) təyin etmək compile xətasıdır.

**Nəyə lazımdır:** Oktal ən çox POSIX permission-lar üçün (`0o777`), binar/heksadecimal
bit filtrləri və network/infrastruktur üçün.

### 3. Rəqəmsal tiplər və seçim qaydaları
**Nədir:** 12 rəqəmsal tip 3 qrupda — int8/16/32/64, uint8/16/32/64, float32/64,
complex64/128.

**Xüsusi adlar:**
- `byte` = `uint8` alias — Go kodunda `uint8` yox, `byte` yazılır.
- `int` = platforma dəyişənli (32-bit CPU-da int32, 64-bit CPU-da adətən int64);
  platformalar arası fərqliliyinə görə `int` ↔ `int32`/`int64` arasında təyin/müqayisə
  compile xətasıdır — açıq conversion tələb olunur. Integer literal-ın default tipi `int`-dir.
- `uint` — `int` kimi, amma unsigned.
- `rune` = `int32` alias — simvol (code point) üçün `int32` yox, `rune` yazılır (niyyət
  aydın olsun).
- `uintptr` — Chapter 14-də.

**Hansı integer-i seçmək (3 qayda):**
1. Binary format/network protokolu konkret ölçü tələb edirsə → həmin tip (`int32` və s.).
2. Hamıya uyğun library funksiyası yazırsansa → bir cüt funksiya: `int64` və `uint64`
   parametrli (məs. `strconv`-dakı `FormatInt`/`FormatUint` kimi).
3. Digər bütün hallarda → `int`. (Başqa tip "premature optimization"dır — profiler sübut
   edənə qədər.)

**Float qaydaları:** Default `float64` (literal-lar da float64 default-dur). `float32`
yalnız 6-7 onluq dəqiqəlikdir. **PUL üçün float İSTİFADƏ ETMƏ** — IEEE 754 dəqiq onluq
təmsil edə bilmir (`-3.1415` yaddaşda `-3.14150000000000018118839761883` kimi saxlanılır).
Float-u `==` ilə müqayisə ETMƏ — epsilon (icazəli fərq) yanaşması işlət.
`0.0-a bölmə` → `+Inf`/`-Inf`; `0.0/0.0` → `NaN`.

**Kompleks ədədlər:** `complex64`/`complex128`, `complex(re, im)` built-in funksiyası;
`real()`/`imag()` hissələri çıxarır; `math/cmplx` paketi. Praktikada az istifadə olunur.

**Kitabdan kod nümunəsi:**
```go
func main() {
    x := complex(2.5, 3.1)
    y := complex(10.2, 2)
    fmt.Println(x + y)            // (12.7+5.1i)
    fmt.Println(real(x))          // 2.5
    fmt.Println(imag(x))           // 3.1
    fmt.Println(cmplx.Abs(x))     // 3.982461550347975
}
```

### 4. String və Rune
**Nədir:** `string` — immutable (dəyişməz) built-in tip; `rune` — tək code point.

**Necə işləyir:** String-lər Unicode dəstəkləyir; `+` ilə birləşdirilir; `==`/`!=`/`<`
ilə müqayisə olunur. String-in daxili bayt dəyərini dəyişmək olmaz — yalnız dəyişənə
yeni string təyin etmək olar. Detallar (bytes/runes/encoding) Chapter 3-də.

### 5. Explicit Type Conversion (Açıq tip çevrilməsi)
**Nədir:** Go-da avtomatik tip yüksəltməsi (automatic type promotion) YOXDUR — bütün
çevrilmələr açıq yazılmalıdır.

**Necə işləyir:** Hətta fərqli ölçülü int-lər və float-lar arasında belə. Truthiness də
yoxdur: heç bir tip `bool`-a nə açıq, nə dolayı çevrilə bilmir — müqayisə operatorları
(`x == 0`, `s == ""`) işlət.

**Nəyə lazımdır:** Go "clarity over conciseness" (qısalıqdan aydınlıq üstünlük) seçir —
conversion qaydalarını əzbərləmək lazım deyil, niyyət kodda görünür.

**Kitabdan kod nümunəsi:**
```go
var x int = 10
var y float64 = 30.2
var z float64 = float64(x) + y   // 40.2
var d int = x + int(y)           // 40
```

**Sub-kod izahı:**
- `float64(x)` → x-i float64-ə çevirir (y-ilə toplanma üçün)
- `int(y)` → y-ni kəsərərək int-ə çevirir (30.2 → 30)

### 6. `var` Versus `:=`
**Nədir:** Dəyişən bəyanatının iki əsas forması.

**Necə işləyir:**
- `var x int = 10` — ən uzun forma (keyword + tip + dəyər)
- `var x = 10` — tip sağ tərəfdən çıxarılır
- `var x int` — dəyər verilmir, zero value alır
- `var x, y = 10, "hello"` — müxtəlif tiplər bir sətirdə
- Declaration list: `var ( ... )` — qrup şəklində
- `x := 10` — yalnız funksiya daxilində; `:=` bir fəndəliyi var: sol tərəfdə **ən azı
  bir yeni** dəyişən olduqda mövcud dəyişənlərə də təyin edə bilir (`x, y := 30, "hello"`
  — x mövcud, y yeni).

**Nə vaxt hansı:**
- Funksiya daxilində default → `:=`
- Zero value üçün → `var x int` (niyyət: zero)
- Untyped literal-ın default tipi istənilməyən tipdirsə → `var x byte = 20` (idiomatik;
  `x := byte(20)` hüquqi, amma az rəsmi)
- `:=` shadowing yaratmaq riski olduqda → yeni dəyişənləri `var` ilə ayır, sonra `=`
  işlət (bax: Chapter 4, Shadowing)
- Çoxdəyişənli təyin yalnız funksiyanın çoxlu qaytarma dəyərləri və ya comma-ok idiomu
  üçün
- Package səviyyəsində `:=` QADAĞANDIR, `var` işlənir; package-level dəyişənlər ümumiyyətlə
  azaldılmalıdır — dəyişən data axınını izlənməz edir.

**Kitabdan kod nümunəsi:**
```go
// := reuse (x mövcud, y yeni)
x := 10
x, y := 30, "hello"
```

### 7. `const` — Literal-lara ad vermək
**Nədir:** Compile vaxtı hesablana bilən dəyərlərə ad verən bəyanat.

**Necə işləyir:** Yalnız bu dəyərlər ola bilər: rəqəmsal literal-lar, `true`/`false`,
string-lər, rune-lər, `complex`/`real`/`imag`/`len`/`cap` built-in-ləri və bunlardan
operatorlarla qurulan ifadələr. Runtime-da hesablanan dəyəri immutable elan etmək
MÜMKÜN DEYİL — immutable array/slice/map/struct YOXDUR.

**Typed vs untyped const:**
- `const x = 10` (untyped) — literal kimi davranır: `var y int = x`, `var z float64 = x`,
  `var d byte = x` — hamısı leqal. Default tipi var (`int`), amma çıxarış olmadıqda.
- `const typedX int = 10` (typed) — yalnız `int`-ə təyin oluna bilər; `float64`-ə təyini
  compile xətası: `cannot use typedX (type int) as type float64 in assignment`.

**Tövsiyə:** Əksər halda untyped saxla (çox çeviklik); enum (iota) hallarında typed işlət.

**İstifadə olunmayan const-lara icazə var:** compile vaxtı hesablanır, side-effect-i
yoxdur — binary-ə sadəcə daxil edilmir.

**Kitabdan kod nümunəsi:**
```go
const x int64 = 10
const (
    idKey   = "id"
    nameKey = "name"
)
const z = 20 * 10
```

### 8. Unused Variables (İstifadə olunmayan dəyişənlər)
**Nədir:** Go-nun unikal qaydası — lokal dəyişən oxunmadıqda compile xətası.

**Necə işləyir:** Dəyişən ən azı bir dəfə oxunmalıdır. Amma yoxlama tam deyil —
oxunmadan qalan təyinatları (`x := 10; x = 20; fmt.Println(x); x = 30`) yalnız
`golangci-lint`-in `ineffassign` linter-i tutur. Package-level unread dəyişənləri isə
kompilyatoru maraqlandırmır — bu da package-level dəyişənlərdən qaçmağın səbəblərindən.

**Kitabdan kod nümunəsi:**
```go
$ golangci-lint run
unused.go:6:2: ineffectual assignment to `x` (ineffassign)
```

### 9. Adlandırma konvensiyaları
**Necə işləyir:**
- İdiomatik Go **camelCase** işlədir; snake_case (`index_counter`) YOX.
- Konstantlar böyük hərflərlə yazılmır (`INDEX_COUNTER` YOX) — çünki böyük hərflə
  başlayan ad paketdən kənar üçün public deməkdir (Chapter 9).
- Funksiya daxilində qısa adlar: scope nə qədər kiçikdirsə, ad bir o qədər qısa; `k`/`v`
  (for-range), `i`/`j` (index) standartdır. Tipin ilk hərfi də adətən dəyişən adı olur
  (`i` int, `f` float, `b` bool).
- Qısa adları izləmək çətindirsə — kod çox iş görür deməkdir; refactor et.
- Package blokunda uzun, deskriptiv adlar.
- Unicode hərfləri leqaldır (`π := 3`), amma look-alike simvollar təhlükəlidir: `ａ`
  (U+FF41) ilə `a` (U+0061) tamamilə fərqli dəyişənlərdir.

## Əsas terminlər
- Zero Value (sıfır dəyəri) — bəyan edilmiş, təyin edilməmiş dəyişənin default dəyəri
- Literal (literal) — kodda birbaşa yazılan dəyər
- Rune (simvol kodu) — int32 alias, tək Unicode code point
- Explicit Type Conversion (açıq tip çevrilməsi) — `T(v)` formasında manual çevrilmə
- Untyped Constant (tipsiz konstant) — öz tipi olmayan, default tipi olan const
- Ineffassign (təsirsiz təyinat) — oxunmayan təyinatı tutan linter
- camelCase (dəvə qoşası) — sözOrtasıBöyükHərf adlandırma stili

## Praktik nəticə

Yeni Go developerinin ən tez düşdüyü tələlər: (1) float-la pul hesablamaq, (2) float-u
`==` ilə müqayisə etmək, (3) `int`-i `int64`-ə açıq conversion olmadan qatmaq, (4)
package-level dəyişənə sahib olmaq. Düzgün default-lar: pul → integer cents; müqayisə →
epsilon; universal seçim → `int` və `float64`; package dəyişəni → yox, const və ya
funksiya daxili dəyişən.

## Mənbə
Pages: 37-64
