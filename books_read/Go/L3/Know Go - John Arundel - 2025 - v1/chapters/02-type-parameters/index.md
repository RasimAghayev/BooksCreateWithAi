# Chapter 2 — Type parameters (Tip Parametrləri)

## Bu fəsil nədən bəhs edir?

Generics sintaksisinin əsası: `func PrintAnything[T any](v T)` — T = tip
parametri, "PrintAnything of T", parameterised function. Instantiation —
çağırış yerində Go T-ni çıxarır (x int → int versiyası compile); "cannot
infer T" halı (nəticəsi T olan çağırış) → `Something[int]()` explicit
instantiation; stencilling (spray-paint metaforası — hər tip üçün ayrıca
maşın kodu, indirection YOX, type assertion YOX). Identity funksiyası
(any-versiya: parametr/nəticə ƏLAQƏSİZ tiplər ola bilərdi → generic versiya:
HƏR İKİSİ T!). PrintAnythingTo məşqi (io.Writer + T). Composite generic
tiplər: `Len[E any](s []E)`, `Drain[E any](ch <-chan E)`, `Merge[E any](chs
...<-chan E)`; konvensional adlar (T, E, BÖYÜK hərf). Generic type-lər:
`type Bunch[E any] []E` — instantiate: `Bunch[int]{1,2,3}`; Bunch[T]
elementləri HAMISI T (mixed-type YOX!); "runtime-da generic tip YOXDUR"
(interfeys — qismən istisna). Group məşqi. Generic funksiya TIPLƏRI:
`type idFunc func[T any](T) T` — SYNTAX ERROR; `type idFunc[T any] func(T) T`
— DOĞRU (generic TİP); `Identity[string]` dəyər kimi, `%T` → func(string)
string; "no generic functions" paradoksu; PrintBunch — generic funksiya
generic tip parametri qəbul edir; Len məşqi; any-ə + məhdudiyyəti generics-də
DƏ davam edir: "operator + not defined on x (variable of type T constrained
by any)" → constraint lazımdır (ch3).

## Əsas fikirlər

### 1. T — İstənilən Tip üçün Yer Tutucu
```go
func PrintAnything[T any](v T) {
    // ...
}
```
- **T = type parameter** — yeni növ parametr; v isə adi (tipi T olan) parametr
- Oxunuş: "For any type T, PrintAnything[T] takes a T parameter"
- **Parameterised function** = generic function = "PrintAnything of T"

### 2. Instantiation — Şablon → Konkret Funksiya
```go
var x int = 5
PrintAnything(x)     // Go anlayır: T = int → int-versiyası compile/call
PrintAnything("hi")  // başqa yerdə: T = string → string-versiyası
```
- **Stencilling:** spray-paint — hər tip üçün "rənglənmiş" ayrıca versiya
- Kodu biz ƏVVƏLCƏDƏN generator-la yazırdıq (go:generate) — indi Go ÖZÜ edir
- **Effektiv maşın kodu:** interfeys indirection-u YOX; assertion YOX — hər
  versiya konkret tipi BİLİR

**Type inference limitation:**
```go
x := Something()      // nəticə T → T bilinmir
// cannot infer T
x := Something[int]() // EXPLICIT instantiation: kvadrat mötərizə!
```
Explicit instantiation HƏMİŞƏ icazəli; yalnız inference alınmayanda LAZIM.

### 3. Identity — Parametr/Nəticə Əlaqəsi
```go
func Identity(v any) any {    // PİS: int ver, string ala bilərsən!
    return v
}

func Identity[T any](v T) T { // DOĞRU: hər ikisi EYNİ T
    return v
}

fmt.Println(Identity("Hello"))   // T = string çıxarılır
```
Generic olmayan həll ƏLAQƏNİ itirirdi — generics onu SÖZLƏŞDİRİR.

### 4. Məşq: PrintAnythingTo
```go
func TestPrintAnythingTo_PrintsToGivenWriter(t *testing.T) {
    t.Parallel()
    buf := new(bytes.Buffer)
    print.PrintAnythingTo(buf, "Hello, world")
    want := "Hello, world\n"
    got := buf.String()
    ...
}

// həll:
func PrintAnythingTo[T any](w io.Writer, p T) {
    fmt.Fprintln(w, p)      // adi + T parametr qarışığı — tam normal
}
```
Go 1.18 tələbi; 1.17-də: `type string is not an expression` (çaşdırıcı
xəta!) → Go-nu YENİLƏ.

### 5. Composite Generic Parametrlər
```go
func Len[E any](s []E) int              // slice of E
func Drain[E any](ch <-chan E)          // channel of E
func Merge[E any](chs ...<-chan E) <-chan E   // variadic channels!
```
T yalnız T-in ÖZÜ deyil — T-dən QURULMUŞ composite tip də! Konvensiya:
BÖYÜK hərf; T = istənilən, E = element tipi.

### 6. Generic Types — `type Bunch[E any] []E`
```go
type SliceOfInt []int        // konkret — hər tip üçün ayrıca? YOX!

type Bunch[E any] []E        // GENERIC: E üçün "slice of E" adı

b := Bunch[int]{1, 2, 3}     // instantiatE → SIRADI []int
b = append(b, "hello")
// cannot use "hello" (untyped string constant) as int value in argument to append
```
**Ən vacib MIS anlayış:** "generic slice" ≠ "mixed-type slice"! Bunch[T]
-in elementləri HAMISI T olmalıdır. Bunch[int] YAXUD Bunch[string] — qarışiq
YOX. İnstantiationdan sonra SIRADI slice-dır.

**"There are no generic types at run time":** compile-da generic definisiya
olunur; istifadədə konkretləşir; runtime-da generic YOXDUR (yalnız konkret).
Qismən istisna: `[]any` (interfeys) — amma "truly generic" deyil.

**Group məşqi:** `type Group[E any] []E` — slices.Equal ilə test.

### 7. Generic Funksiya Tipləri — İncə Fərq
```go
f := Identity[string]           // instantiatE = DƏYƏR (funksiyalar dəyərdir!)
fmt.Printf("%T\n", f)
// func(string) string           ← adi konkret funksiya tipi!

type Stringulator func(string) string    // konkret funksiya tipinə ad

type idFunc func[T any](T) T
// syntax error: function type must have no type parameters

type idFunc[T any] func(T) T     // DOĞRU: bu, GENERIC TİP-dir!
```
**"There are no generic functions":** smart-alec həqiqət — generic funksiya
"şablondur"; konkret funksiyalar (Identity[string]) compile-da yaranır.
Funksiya tipinə tip parametri QADAĞA; amma generic TİP ola bilər
(idFunc[T] instantiatE olunanda func(int) int kimi konkretləşir).

### 8. Generic Tip = Generic Funksiyanın Parametri
```go
func PrintBunch[E any](v Bunch[E]) {
    fmt.Println(v)          // Bunch[int] ver → çap
}
```
İki sistemin birləşməsi — "pretty exciting stuff".

**Len məşqi:**
```go
func Len[E any](s []E) int {
    return len(s)           // built-in!
}
```

### 9. any + Operatorlar — Limit Davam Edir
```go
func AddAnything[T any](x, y T) T {
    return x + y
}
// invalid operation: operator + not defined on x
// (variable of type T constrained by any)
```
**Eyni köklü problem:** T any = HƏR tip; struct + struct mənasızdır → Go
"guarantee" tələb edir: instantiate OLUNACAQ tiplərin HAMSİ + dəstəkləsin.
Həll: T-i CONSTRAINT-lə (ch3) — icazəli tip dəstinə restrict et.

## Əsas terminlələr
- Type Parameter (T) — tip üçün yer tutucu parametr
- `[T any]` Sintaksisi — tip parametr siyahısı
- Parameterised / Generic Function — "PrintAnything of T"
- Instantiation — konkret tip üçün versiya yaratma
- Type Inference — çağırışdan T-ni çıxarma
- Explicit Instantiation — `Something[int]()`
- "cannot infer T" — nəticə-T çağırışları
- Stencilling — hər tip üçün ayrıca kod (spray-paint)
- No Indirection — interfeysdən fərqli: birbaşa maşın kodu
- Identity[T] — parametr/nəticə eyni-T müqaviləsi
- Composite Generic Parameter — []E, chan E, ...chan E
- Konvensional Adlar — T, E; BÖYÜK hərf
- Generic Type — `type Bunch[E any] []E`
- "No Generic Types at Run Time" — ancaq konkret instantiations
- Mixed-Type MIS — Bunch[T] elementləri homojendir
- `[]any` İstisnası — interfeys; "truly generic" deyil
- Funksiya Dəyəri — `f := Identity[string]`
- %T — tipi çap edən fmt verb
- `type idFunc[T any] func(T) T` — generic funksiya TİPİ (legal)
- "Function Type Must Have No Type Parameters" — funksiya tipinə T QADAĞA
- PrintBunch — generic funksiya × generic tip
- "T Constrained by Any" — operator xətasının formulasiyası
- Constraint (ön baxış) — icazəli tip dəsti (ch3)

## Praktik nəticə
(1) Generic funksiya: `[T any]` + T parametr; inference adətən işləyir,
nəticə-T halında explicit `[int]` yaz. (2) any YERİNE T-in gücü: ƏLAQİLİ
tipləri (parametr↔nəticə) sözleşdirmək — Identity dərsi. (3) Composite:
[]E, chan E, variadic ...<-chan E — T-dən qurulu parametrlər sərbəstdir. (4)
Generic tip: `type Ad[E any] []E`; instantiate — `Ad[int]{...}`; elementlər
HOMOJEN — "mixed slice" İKİSİ DƏ YANLIŞDIR. (5) `f := F[string]` —
instantiatE olunmuş generic funksiya DƏYƏRDIR (%T → func(string) string). (6)
`type F func[T any](T) T` — SYNTAX ERROR; düzgünü: `type F[T any] func(T) T`
(generic TİP). (7) any + operatorlar = compile xətası — "T constrained by
any" işləməz; növbəti addım: constraint-lə tip dəsti məhdudlaşdır (ch3). (8)
Go 1.18+ şərt; "type string is not an expression" = köhnə Go siqnalı.

## Mənbə
Pages: 25-43 (PDF 26-44)
