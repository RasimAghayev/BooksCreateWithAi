# Chapter 13 — Fun with functions (Funksiyalarla Əyləncə)

## Bu fəsil nədən bəhs edir?

Funksiya anatomiyası: declaration (func + ad + parameter list + result list +
body), parameter adları/çoxlu eyni tip (`x, y, z float64`), result list
`(float64, float64, error)`, SIGNATURE anlayışı (parametr+result tip kombinasiyası),
çağırış = control flow (jump → body → geri qayıt; "get orada bunu et, sonra
buraya qayıt"), arqumentlər, nəticə ifadə kimi (`3 * (double(2.5) + 7)` — tək
nəticədə; çox nəticə = YALNIZ tuple), return (nəticəsiz implicit bitiş; nəticəli =
say uyğunluğu məcburi), withdraw məşqi (balance/amount → newBalance/error,
2 test); **functions are values** (assign, parametr ötür, qaytar), TestCase
struct-da `function func(float64, float64) float64` sahəsi (Add/Subtract/Multiply
eyni signature — bir testlə hamısı!), sort.Slice-in less parametri; function
literal (adsız, on-the-fly, math.Pow nümunəsi, sort üçün `nums[i] < nums[j]`),
apply məşqi; **closures** (nums haradan gəlir? — scope görünürlüyü; "closure
over nums", bubble analogy — dəyişənlər içinə həbs olunur); **loop variable
tələsi** (funcs slice → hamısı 3 çap edir! — closure dəyəri ÇAĞRILIŞDA alır,
definisiyada YOX); defer (resource leak, f.Close hər exit-də unudulur, "funksiya
çıxanda HƏR HANSİ yolla icra et"), stacking (LIFO — son defer əvvəl icra),
named results (latitude/longitude — documentation; funksiya daxilində dəyişən
kimi), **naked returns considered harmful** (qadağa kimi qəbul et), deferred
closure ilə nəticənin ÇIXIŞDAN SONRA dəyişməsi (closeErr → err = closeErr;
yalnız fail-da overwrite — `err = f.Close()` tək sətir YANLIŞ: uğurda nil-i
silərdi!), variadic (`...float64` → daxildə SLICE, fmt.Println, AddMany/
DivideMany(12,4,3) məşqləri).

## Əsas fikirlər

### 1. Declaration — func + Ad + Parametrlər + Nəticələr + Body
```go
func double(x float64) float64 {
    return x * 2
}
```
- Yalnız PAKET səviyyəsində (funksiyalar daxilində YOX — literal istisna)
- Parametrlər: ad + tip; adsız da OLAR (istifadə edilmirsə nəyə lazım?)
- Eyni tip qısaltması: `add(x, y, z float64)` (təkrar-təkrar tip YOX)
- Result list: `(float64, float64, error)` — adlar OPTIONAL (bax: named results)
- **Signature** = parametr + result TİPLƏRİNİN kombinasiyası — funksiyanın
  kimliyi
- Boş body legaldır (amma faydasız); çağırılmayan funksiya icra olunmur

### 2. Çağırış = Control Flow
`hello()` → icra hello-nun BAŞINA atılır → bitəndə ÇAĞIRILAN YERİN BİRDƏN
SONRAKISINA qayıdış. "Get orada bunu et, sonra buraya qayıt."
- Arqumentlər = parametrlərlə SAY/ TİP uyğun
- `fmt.Println(1, 2, 3)` — standart kitabxana nümunəsi

### 3. Nəticə İfadələri
```go
answer := double(2.5)              // → 5
answer := 3 * (double(2.5) + 7)    // tək nəticə = ifadə daxilində sərbəst
lat, long, err := location()       // çox nəticə = YALNIZ tuple assignment
```

### 4. return Qaydaları
- Nəticəsiz funksiya: body sonuna çatma = implicit return; erkən `return`
- Nəticəli: `return 0, errors.New(...)` — SAY imza ilə eyni olmalı (dəyərlər
  ixtiyari, SAY mütləq)

**withdraw məşqi:**
```go
func Withdraw(balance, amount int) (int, error) {
    if amount > balance {
        return 0, errors.New("insufficient funds")
    }
    return balance - amount, nil
}
// TestWithdrawValid: 100, 20 → 80, nil
// TestWithdrawInvalid: 20, 100 → _, err != nil
```

### 5. Functions Are Values!
Funksiya = DƏYƏRDIR: assign et / parametr ötür / qaytar. Konsekvensiya:
**TestCase struct-da funksiya sahəsi:**
```go
type TestCase struct {
    a, b     float64
    function func(float64, float64) float64   // ← TİP! imza
    want     float64
}
// literal:
{a: 2, b: 2, function: calculator.Add, want: 4}
```
`func(float64, float64) float64` = "2 float64 parametr + 1 float64 nəticə"
imzası — Add/Subtract/Multiply HAMISI uyğundur → 3 funksiya 1 TABLE TEST ilə!

**Funksiya parametri** — standart kitabxanadan:
```go
func Slice(x any, less func(i, j int) bool)
```
less = müqayisə funksiyası; sort.Slice onu təkrar-təkrar çağırır.

### 6. Function Literal — Adsız Dəyər
Hər dəyər kimi funksiyanın da LITERAL-I var:
```go
function: func(a, b float64) float64 {
    return math.Pow(a, b)
},
```
Adlı funksiya elan etmədən ON THE FLY yarat (math.Pow nümunəsi — məzmun vacib
deyil, ŞƏKİL-dir). sort.Slice klasikası:
```go
nums := []int{3, 1, 2}
sort.Slice(nums, func(i, j int) bool {
    return nums[i] < nums[j]     // i/j-nin ÖZÜ deyil, ELEMENTLƏRİ müqayisə
})
```
(apply məşqi: `apply(1, func(x int) int { return x * 2 })` → 2.)

### 7. Closures — "Bubble" Effect
Sual: `nums` literalının İÇİNDƏ haradan gəlir? Parametr deyil! Cavab: literal
definisiya olunduğu SKOPla görür — closure!
- Package-level funksiya paket dəyişənlərini görür — eyni məntiq
- Function literal = "closure over" əhatə dəyişənləri
- sort.Slice üçün VACİB: less imzası sabit (2 int) — slice neylə ötürlsün?
  CLASURE! "Bubble" analogy: dəyişənlər baloncuxa həbs → balon istənilən
  yerə ötürülür → dəyişənlər içeridə AVAILABLE

### 8. Loop Variable Tələsi — Klassik Gözləntilərdən İmtina
```go
funcs := []func(){}
for _, v := range []int{1, 2, 3} {
    funcs = append(funcs, func() {
        fmt.Println(v)        // closure over v
    })
}
for _, f := range funcs {
    f()
}
// Gözlənti: 1, 2, 3
// REALİK: 3, 3, 3  !!
```
**Sual:** closure dəyəri DEFINİSİYA anında yoxsa ÇAĞIRIŞ anındakını görür?
**CAVAB: ÇAĞIRIŞ anındakını!** Loop bitəndə v = 3; 3 funksiya hamısı ÇAĞIRIŞDA
v-yə baxır → 3, 3, 3. Bu, DÜZGÜN davranışdır — sadəcə gözlənilməzdi.
**Dərs:** closure-un gördüyü dəyişəni YENİLƏYİRSƏNSƏ — çağırış vaxtına diqqət
(ən çox loop-larda baş verir).

### 9. defer — "Çıxanda Hər Hansı Yolla"
Problem: f.Close() body sonunda; amma returnlər ORTA YERDƏ → f BAĞLANMIR →
resource LEAK (OS/runtime yaddaş saxlayır; long-running serverdə böyüyüb
CRASH edir). "Hər return-da bağlamağı xatırla" — heç kim qüsursuz deyil; kodu
dəyişən başqası bilməz də.

```go
f, err := os.Open("testdata/somefile.txt")
if err != nil {
    return err            // açılmayıbsa BAGLAMAQ YOXDUR
}
defer f.Close()           // "çıxanda bağla" — İNDİ QEYD ET
...                       // istənilən return — f avtomatik bağlanır
```
defer funksiya ÇAĞIRIŞINI təxirə salır: indi YOX, funksiya çıxarkən — "no
matter where or how".

### 10. Stacking Defers — LIFO
```go
defer cleanup1()   // exit-də SONUNCU icra
defer cleanup2()   // exit-də İLK icra
```
İstənilən sayda; TƏRS SIRA (last deferred, first run).

### 11. Named Result Parameters
```go
func location() (latitude float64, longitude float64, err error)
```
- **Əsas dəyər: DOKUMENTASİYA** — kod insanlar üçündür; (float64, float64)
  → nədir? latitude/longitude AYDIN edir
- Funksiya daxilində adi dəyişən kimi görünür; return adi qalır:
```go
return 50.5897, -4.6036, nil       // hələ də İFADİLİ
latitude = 50.5897; longitude = -4.6036
return latitude, longitude, err     // yaxud belə
```

### 12. Naked Returns Considered Harmful
Adlı nəticədə `return` (bos) — LEGAL: cari latitude/longitude dəyərləri
qaytarılar. **Amma PRAKTİKA DEYİL:** aydınlıq üçün dəyərləri AÇIQ yaz —
"heç bir faydası yoxdur, özünü qurtarmaq üçün yazmayacaqsan". Adlı nəticə
OLMASI naked return MƏCBURİYYƏTİ DEMƏK DEYİL — "can, and should, make explicit".

### 13. Deferred Closure — Nəticəni Çıxışdan Sonra Dəyişmək
Ssenari: fayla yazırıq, `return nil` (uğur); amma deferred f.Close() FAIL edə
bilər → istifadəçi datası İTİB — error QAYTARILMALI! Problem: returnnil artıq
"baked in"... Həll zənciri:

**Addım 1 — function literal defer:**
```go
defer func() {
    closeErr = f.Close()
    if closeErr != nil {
        fmt.Println("oh no")     // heç olmasa ağlayaq...
    }
}()
```
(Boş mötərizə axırda! defer FUNGSIYA ÇAĞIRIŞINI gecikdirir — literal sonrası
`()` = çağırış.)

**Addım 2 — closure + named result = MÜKƏMMƏL həll:**
```go
func WriteData(...) (err error) {    // ADLI nəticə!
    ...
    defer func() {
        closeErr = f.Close()
        if closeErr != nil {
            err = closeErr          // closure err-i GÖRÜR və DƏYİŞİR
        }
    }()
    ...
    return nil
}
```
"Nəticəni funksiya çıxdıqdan SONRA (qayıtmazdan ƏVVƏL) dəyişmək" — named
results + defer + closure ÜÇLÜYÜ ilə!

**Vacib incəlik:** `err = f.Close()` (şərtsiz) YANLIŞDIR — err artıq return-ın
dəyərini daşıyır (məs. nil); close UĞURLU olsa da SİLİB nil yazardı. Yalnız
FAIL halda overwrite!

### 14. Variadic Functions — ...
```go
fmt.Println(1, 2, 3)            // istənilən sayda arqument!

func AddMany(inputs ...float64) float64 {
    total := 0.0
    for _, input := range inputs {   // inputs DAXİLDƏ SLİCE KİMİ!
        total += input
    }
    return total
}
```
- `...float64` = 0, 1, 2... istənilən sayda float64
- Daxildə inputs = SLICE (`[]float64` kimi range/append/index)
- Məşqlər: AddMany, SubtractMany, MultiplyMany (10 input → hasil),
  DivideMany(12, 4, 3) → 12/4/3 = 1 (soldan-sağa zəncir)

## Əsas terminlələr
- Signature — parametr + result tipləri; funksiyanın tipi
- Parameter / Result List — ad+tip / tiplər (adlar optional)
- Çoxparametr Qısaltması — x, y, z float64
- Arqument — çağırışda ötürülən dəyər
- Çağırış = Control Flow — jump + qayıdış
- Tək/Çox Nəticə — ifadədə sərbəst / yalnız tuple
- return Say Qaydası — nəticə sayı = imza sayı
- Functions Are Values — assign/öter/qaytar
- func(float64, float64) float64 — funksiya TİPİ yazılışı
- TestCase.function — table testdə funksiya sahəsi
- sort.Slice(x, less) — funksiya parametrli standart funksiya
- Function Literal — adsız, on-the-fly funksiya
- Closure — definisiya skopunu görən literal ("closure over X")
- Bubble Analogy — dəyişənlər həbs olunub gəzir
- Loop Variable Tələsi — closure ÇAĞIRIŞ vaxtını görür → 3,3,3
- Resource Leak — bağlanmamış fayl; GC təmizləyə bilmir
- defer — "çıxanda icra et, necə çıxırsa"
- defer f.Close() — uğurlu açılışdan DƏRHAL sonra
- Stacking Defers (LIFO) — son defer əvvəl icra
- Named Result Parameters — latitude/longitude; documentation dəyəri
- Naked Return — bos return; LEGAL amma HARMFUL — explicit yaz
- Deferred Closure — err = closeErr; nəticəni çıxışdan sonra dəyiş
- Şərtsiz err = f.Close() YANLIŞ — uğurlu halda nil-i silər
- Variadic (...T) — istənilən sayda; daxildə slice
- AddMany / DivideMany — variadic məşqləri

## Praktik nəticə
(1) Eyni imzalı funksiyalar → table testdə `function` sahəsi ilə bir testə yığ.
(2) Funksiya parametrli API (sort.Slice modeli): less func(i, j) — closure
data-nı ÖZÜ daşıyır, ötürməyə ehtiyac YOX. (3) Closure dəyişəni YENİLƏYİRSƏNSƏ
çağırış vaxtını xatırla — loop içində literal yığırsansa dəyərlər SON versiya
olacaq. (4) Resurs açdın — DƏRHAL defer Close (uğursuz açılış branch-ından
sonra!); LIFO sırasını nəzərə al. (5) Nəticələri ADLANDIR — documentation
(hansı float64 hansı koordinatdır?); amma return HƏMİŞƏ explicit — naked
return yazma. (6) Close-ın error-u return-ə düşməli olsa: named err + deferred
closure + YALNIZ fail halda overwrite. (7) Çox nəticəli çağırışı ifadə daxilinə
qoyma — yalnız tuple. (8) İstənilən sayda arqument lazımdırsa `...T`; daxildə
slice kimi işlət. (9) Parametrlərdə eyni tipləri qısald (x, y, z float64).
(10) return-un sayı imza ilə dəqiq uyğun — compiler qəbul etməz əks halda.

## Mənbə
Pages: 161-177 (PDF 162-178)
