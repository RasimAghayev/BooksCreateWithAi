# Chapter 8 — Errors (Xətalar)

## Bu chapter nədən bəhs edir?

Go-nun error yanaşması: error interface, string xətalar, sentinel error-lar, custom error
tipləri (və nil tələsi), error wrapping (%w, errors.Is/As), defer ilə wrapping, panic və
recover, stack trace üçün üçüncü tərəf kitabxanalar.

## Əsas fikirlər

### 1. Error Əsasları — Dəyər Qaytarma
**Nədir:** `error` — tək metodlu interface: `type error interface { Error() string }`.

**Necə işləyir:** Error həmişə **sonuncu** qaytarma dəyəridir; uğurda `nil`. Çağıran
`if err != nil` ilə yoxlayır. Exception-dan fərqi (2 səbəb):
1. Exception gizli kod yolları yaradır — çağırış zəncirində nəyin partlayacağı bəlli deyil.
2. Unused-variable qaydası error-u ya yoxlamağa, ya `_` ilə açıq iqnor etməyə MƏCBUR
   edir — compiler dəstəyi ilə şəffaflıq.

**Tərzi:** Error handling `if` daxilində indent olunur, business logic solda "golden
path" kimi qalır — vizual ayırdetmə. Mesajlar: kiçik hərflə, sonunda punktuasiya/newline YOX.

**Kitabdan kod nümunəsi:**
```go
func calcRemainderAndMod(numerator, denominator int) (int, int, error) {
    if denominator == 0 {
        return 0, 0, errors.New("denominator is 0")
    }
    return numerator / denominator, numerator % denominator, nil
}
```

### 2. Sadə Xətalar — errors.New və fmt.Errorf
```go
errors.New("only even numbers are processed")
fmt.Errorf("%d isn't an even number", i)   // format verb-ləri
```

### 3. Sentinel Error-lar
**Nədir:** Davamın mümkünsüzlüyünü bildirən xüsusi package-level dəyərlər (Dave Cheney
termini). Adlandırma: `Err` prefiksi (`zip.ErrFormat`, `rsa.ErrMessageTooLong`,
`io.EOF`, `context.Canceled`). Test: `err == zip.ErrFormat`.

**Qaydalar:**
- Nadir hallarda təyin et — public API-yə çevrilir, geri götürülməz.
- Standart kitabxanadakını reuse et, yoxsa kontekst daşıyan custom tip yarat.
- Yalnız: konkret vəziyyət + əlavə kontekst lazım deyil.
- Konstant sentinel-lər (Dave Cheney təklifi: `const ErrFoo = consterr.Sentinel("foo")`)
  — unidiomatik: eyni string-li error-lar bərabər çıxar (paketlərarası bərabərlik
  istənməz).

### 4. Error-lar Dəyərdir — Custom Tiplər (və nil Tələsi!)
**Nədir:** Status kodu kimi strukturallaşdırılmış məlumat daşıyan error tipləri — string
müqayisəsindən xilas.

**Kitabdan kod nümunəsi:**
```go
type Status int
const (
    InvalidLogin Status = iota + 1
    NotFound
)
type StatusErr struct {
    Status  Status
    Message string
}
func (se StatusErr) Error() string { return se.Message }

if err != nil {
    return nil, StatusErr{Status: InvalidLogin, Message: "invalid credentials for " + uid}
}
```

**Qaydalar:** Qaytarma tipi həmişə `error` (custom tip YOX) — müxtəlif error tipləri +
çağıranın asılılıqdan azadlığı. Sahələrə type assertion ilə YOX, `errors.As` ilə çıx.

**MÜHÜM TƏLƏ (interface nil semantikası):**
```go
func GenerateError(flag bool) error {
    var genErr StatusErr        // YANLIŞ yanaşma!
    if flag { genErr = StatusErr{...} }
    return genErr               // flag=false olsa da non-nil qayıdır!
}
```
Səbəb: interface nil üçün həm tip, həm dəyər nil olmalı — `StatusErr` tipi non-nil.
Həll: (1) uğurda açıq `return nil`; (2) lokal dəyişəni `error` tipli elan et.

### 5. Error Wrapping — Zəncir
**Nədir:** Error-u kontekstlə qoruyaraq qaytarmaq; zəncir = wrap olunmuş errorlar silsiləsi.

**Kitabdan kod nümunəsi:**
```go
func fileChecker(name string) error {
    f, err := os.Open(name)
    if err != nil {
        return fmt.Errorf("in fileChecker: %w", err)  // %w = wrap
    }
    f.Close()
    return nil
}
// wrap OLMADAN mesaj saxlamaq: %v (zəncir qırılır — bəzən düzgündür)
return fmt.Errorf("internal failure: %v", err)
```

**Custom tipdə wrapping:** `Unwrap() error` metodu implement et:
```go
type StatusErr struct {
    Status  Status
    Message string
    err     error
}
func (se StatusErr) Unwrap() error { return se.err }
```

### 6. errors.Is və errors.As
**Nədir:** Wrap olunmuş zəncirlərdə axtarış üçün standart funksiyalar.

**errors.Is (instansiya axtarışı):**
```go
if errors.Is(err, os.ErrNotExist) { ... }  // zəncirdə hər səviyyə == müqayisəsi
```
Custom `Is(target error) bool` metodu ilə: noncomparable tiplər (`reflect.DeepEqual`)
və ya **pattern matching** (filtrlər — `ResourceErr{Resource: "Database"}` bütün kodlarla
uyğun):

**errors.As (tip axtarışı):**
```go
var myErr MyErr
if errors.As(err, &myErr) { fmt.Println(myErr.Code) }

// interface pointer də mümkün:
var coder interface{ Code() int }
if errors.As(err, &coder) { fmt.Println(coder.Code()) }
```
`As` 2-ci parametri error/interface pointer-dən başqa bir şeydirsə **panic**.

**Tövsiyə:** Spesifik dəyər/instans axtarırsan → `Is`; tip axtarırsan → `As`.

### 7. defer ilə Wrap (Ortak Mesaj)
**Nədir:** Eyni mesajla çoxsaylı wrap-i tək yerə toplama.

**Kitabdan kod nümunəsi:**
```go
func DoSomeThings(val1 int, val2 string) (_ string, err error) {
    defer func() {
        if err != nil {
            err = fmt.Errorf("in DoSomeThings: %w", err)
        }
    }()
    val3, err := doThing1(val1)
    if err != nil { return "", err }
    val4, err := doThing2(val2)
    if err != nil { return "", err }
    return doThing3(val3, val4)
}
```
Adlı return mütləqdir (err-i defer daxilində oxu/yenilə); bir dəyər adlanırsa hamısı
adlanmalı → `_ string`. Fərqli mesajlar lazımdırsa — hər yerdə konkret `fmt.Errorf`.

### 8. panic və recover
**Nədir:** panic — runtime-ın davam edə bilməməsi (programmer xətası: slice həddi aşımı;
mühit: yaddaş bitməsi). Funksiya dərhal çıxır, defer-lər yığılır, main-ə qədər — sonra
program stack trace ilə çıxır.

**recover pattern:**
```go
func div60(i int) {
    defer func() {
        if v := recover(); v != nil {
            fmt.Println(v)
        }
    }()
    fmt.Println(60 / i)  // 0 → "runtime error: integer divide by zero" çap, davam
}
```
Yalnız defer daxilində işləyir (panic-dən sonra yalnız defer-lər icra olunur).

**Qaydalar:**
- Bu, exception handling DEYİL. Recover nəyin failed olduğunu göstərmir — idiomatik Go
  mümkün xəta şərtlərini AÇIQ göstərir.
- Panic-dən sonra davam etmək nadir hallarda düzgündür — resource bitkisində: recover →
  monitorinqə log → `os.Exit(1)`. Programmer xətasında eyni xəta təkrarlanacaq.
- **Library qaydası:** public API-dən panic çıxmağa icazə YOX — recover ilə error-a
  çevir, çağırana qaytar.
- Go-nun HTTP server-i handler panic-lərini recover edir — amma Go team bunu indi
  səhv sayır (David Symonds).

### 9. Stack Trace
Default olaraq Go error-da stack trace VERMİR. Wrap zənciri əl ilə "call stack" qurur.
Üçüncü tərəf error kitabxanaları (məs. github.com/pkg/errors) stack-lı wrap təqdim edir;
trace `fmt.Printf("%+v", err)` ilə çap olunur. Qeyd: trace fayl yollarını daşıyır —
gizlətmək üçün build-də `-trimpath` flag-i.

## Əsas terminlər
- Sentinel error (gözcü xətası) — konkret dəyərlə davamsızlıq bildirən error
- Error wrapping (xəta bükülməsi) — %w ilə orijinalı saxlayaraq kontekst əlavəsi
- Error chain (xəta zənciri) — wrap olunmuş errorlar ardıcıllığı
- Golden path (qızıl yol) — xəta olmayan əsas icra yolu
- panic/recover — qəfacə dayanma və yalnız defer-də tutulması
- Pattern matching filter (nümunə filtri) — sahələrinin biri boş keçilən Is müqayisəsi

## Praktik nəticə

Error idarəetmə qərar ağacı: (1) yalnız mesaj → `errors.New`/`fmt.Errorf`; (2) çağıranın
bərabərlik testi lazımdırsa → sentinel (nadir!); (3) struktur məlumat → custom tip +
həmişə `error` qaytar; (4) kontekst əlavə et → `%w` (bərpa lazımdırsa) / `%v` (bərpa
istəmirsən); (5) zəncirdə axtarış → `errors.Is`/`errors.As`, heç vaxt `==`/assertion;
(6) eyni wrap → defer pattern; (7) panic yalnız bərpaolunmaz + library sərhədində
recover; (8) `var err CustomType` + return → həmişə non-nil interface tələsi — `return nil`
və ya `error` tipli lokal.

## Mənbə
Pages: 231-252
