# Chapter 3 — Errors and expectations (Xətalar və Gözləntilər)

## Bu fəsil nədən bəhs edir?

Divide implementasiyası: `(float64, error)` imza dəyişikliyi (mötərizə tələbi),
compiler xətalarının təhlili ("assignment mismatch: 2 variables but ... returns
1 values", "not enough arguments to return: have/want"), `return a / b, nil`,
invalid-input testi (TestDivideInvalid, blank identifier `_`), b == 0
detection, errors.New, zero value konvensiyası, test strukturunun ümumiləşdirilməsi,
"red, green, refactor" workflowunun 7 addımı, floating-point dəqiqliyi
(1/3 = 0.333... problemi, closeEnough helper, math.Abs, tolerance, kosmik gəmi
vs robot cərrahiyyə), Sqrt məşqi (negative input error), go run / go build
(cmd/calculator/main.go, package main, main funksiyası, `-o add`).

## Əsas fikirlər

### 1. Divide — Compiler Xətaları Arxılı İnkişaf
Ilkin draft (1 nəticə):
```go
func Divide(a, b float64) float64 {
    return a / b
}
```
Xəta 1 — test 2 dəyər gözləyir:
```
./calculator_test.go:80:12: assignment mismatch: 2 variables
but calculator.Divide returns 1 values
```
Tərcümə: `:=` sol tərəfdə 2 dəyişən var (`got, err`), funksiya 1 qaytarır.
Həll — imzada 2 nəticə ELAN ET:
```go
func Divide(a, b float64) (float64, error) {
```
**Vacib:** 2+ nəticə MÖTƏRİZƏ tələb edir `(float64, error)`; tək nəticədə
mötərizə optional idi.

Xəta 2 — return az dəyər verir:
```
./calculator.go:30:2: not enough arguments to return
have (float64)
want (float64, error)
```
Həll: `return a / b, nil` — nil = "xəta yoxdur". → PASS (valid davranış hazır).

**Dərs:** Go compiler xətaları "intimidating" görünür, amma EXTREMELY accurate:
fayl:sətir + have/want — hər fix üçün lazımlı BÜTÜN informasiya var.

### 2. Invalid Input Testi — TestDivideInvalid
"One behaviour, one test" — yeni davranış, yeni test:
```go
func TestDivideInvalid(t *testing.T) {
    t.Parallel()
    _, err := calculator.Divide(1, 0)
    if err == nil {
        t.Error("want error for invalid input, got nil")
    }
}
```
- **`_` blank identifier:** got lazım DEYİL (error halında data etibarsızdır),
  amma sintaksis sol tərəf İSTƏYİR → `_` = "bu dəyərə ehtiyacım yoxdur"
- **Compiler qaydası:** istifadə olunmayan DƏYİŞƏN = compile xətası (düzgün
  təxmin: səhvdir) — `_` bunun qaçış yolu
- **t.Error (f-siz):** format arqumentləri YOXDUR → Errorf lazımsız
- Bu test DİQQƏTİ yalnız err-ə verir: non-nil olmalı; nə olduğu ƏHƏMİYYƏTSİZ

### 3. Invalid Input Detection + errors.New
```go
func Divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("division by zero not allowed")
    }
    return a / b, nil
}
```
- **b == 0 yoxlaması:** a-nın 0 olması XƏTA DEYİL (nəticə 0-dır); b == 0 =
  qeyri-mümkün bölmə → error davranışı
- **errors.New("mesaj")** — string→error konstruktoru; fmt.Printf ilə çap
  olunanda mesaj görünür
- **Zero value konvensiyası:** error halında data = tipin sıfır dəyəri
  (float64 → 0) — "something and error" semantikası: non-nil error = data
  İGNORE edilməlidir

### 4. Testlərin Ümumi Strukturu
1. Input-ları və gözlənilən nəticələri HAZIRLA
2. Test edilən funksiyanı UYĞUN input-la ÇAĞIR
3. Nəticəni gözlənti ilə MÜQAYİSƏ et; uyğunsuzluqda İZAHLI mesajla FAIL et

### 5. "Red, Green, Refactor" — 7 Addımlıq Workflow
1. Funksiya üçün test yaz (hələ mövcud olmasa da)
2. Compile xətası gör (funksiya yoxdur) — bu, "red"-ə gedən yol
3. Compile OLsun, amma test FAIL etsin — minimum kod
4. Fail-i GÖZLƏNİLƏN səbəbdən olduğunu yoxla; mesajın dəqiq/informativ olduğunu
   yoxla
5. Testi KEÇDİRƏCƏK minimum kodu yaz
6. İstəyə görə: kodu təkmilləşdir (testi QIRMADAN)
7. **Commit!**

"Red" = fail eden test, "Green" = keçən test, "Refactor" = gözəlləşdirmə.
Mürəkkəb funksiya = parçala: ayrı funksiyalar YAXUD sub-behaviour-lar + ayrı
testlər; sonra addım-addım yığ.

### 6. Floating-Point — Dəqiqlik Problemi
float64 = 64-bit floating-point. 1/3 = 0.333333... — sonlanmır; sonlu bitdə
ITİLƏN dəqiqlik. Nümunə:
```go
{a: 1, b: 3, want: 0.333333},
// FAIL: Divide(1.000000, 3.000000): want 0.333333, got 0.333333
```
GÖRÜNÜŞ eyni, == fərqli görür (fərq print-də görünməyəcək qədər kiçik).
**Həll — "close enough" müqayisəsi:**
```go
func closeEnough(a, b, tolerance float64) bool {
    return math.Abs(a-b) <= tolerance
}

if !closeEnough(tc.want, got, 0.001) {
    ...
}
```
**Tolerance = kontekstə bağlı:** 0.001 kifayət ola bilər; kosmik gəmi naviqasiyası
robot beyin cərrahiyyəsindən AZ onluq tələb edə bilər. Tune olunur.

### 7. Sqrt Məşqi — Sərbəst Praktika
VP of Sales: Enterprise kalkulyatora Sqrt lazım! Tapşırıq:
- `Sqrt(float64) (float64, error)` — kök + negative input-da error
  (heç bir real ədədin kvadratı mənfi deyil)
- Testlər Divide modelində: TestSqrt (valid, closeEnough ilə — tolerance bəlkə
  0.1 qədər böyük) + TestSqrtInvalid (err != nil)
- Standart paketlərdən istifadə icazəli (math.Sqrt)

### 8. Proqramın İşə Düşürülməsi — go run / go build
Testlər go test ilə idi; İCRA OLUNAN proqram üçün — executable binary lazımdır
(machine code, OS formatında). Compiler: mənbə kodu → binary.

**package main + main() tələbi:**
- Executable üçün XÜSUSI paket: `package main`
- Bir qovluqda bir paket (test paketlərindən başqa) → main alt-qovluqda
- Konvensiya: `cmd/` (command) qovluğu → `cmd/calculator/main.go`
- İcra `main()` funksiyasından BAŞLAYIR, sətir-sətir; sonuna çatanda proqram
  DAYANIR

```go
package main

import (
    "calculator"
    "fmt"
)

func main() {
    result := calculator.Add(2, 2)
    fmt.Println(result)
}
```
Struktur:
```
calculator/
    calculator.go
    calculator_test.go
    cmd/
        calculator/
            main.go
    go.mod
```

**go run** — compile + icra bir komandada (sürətli yoxlama):
```bash
go run cmd/calculator/main.go
# 4
```

**go build** — paylaşılan BINARY yaradır:
```bash
go build -o add ./cmd/calculator        # Unix
go build -o add.exe ./cmd/calculator    # Windows
./add                                    # → 4
```
Xəta YOXDURSA output YOXDUR; fayl yaranır. (build haqqında daha çox: "Building
blocks" fəsli.)

Kalkulyator layihəsi BAŞA ÇATDI — növbəti gig: ONLAYN KİTAB MAĞAZASI.

## Əsas terminlələr
- (float64, error) — çoxnəticəli imza; mötərizə MÜTLƏQ
- assignment mismatch — sol/sağ dəyər sayı uyğunsuzluğu
- not enough arguments to return — return azlıq; have/want formatı
- nil error — "xəta yoxdur" dəyəri
- Blank Identifier (_) — istifadəolunmayan dəyərin yer tutucusu
- errors.New — string-dən error dəyəri yaratma
- Zero Value Konvensiyası — error halında data = 0/""/nil tipə görə
- t.Error vs t.Errorf — formatsız / formatlı fail
- Red, Green, Refactor — TDD tsikli adları
- 7-Addımlıq Workflow — test→compile-error→fail→düzgün-fail→pass→təkmil→commit
- Floating-Point — sonlu bitdə kəsir dəqiqliyi
- closeEnough(a, b, tolerance) — math.Abs fərq ≤ tolerance
- Tolerance — kontekstə bağlı "kifayət qədər yaxın" həddi
- package main / main() — executable başlanğıc paketi/funksiyası
- cmd/ Konvensiyası — executable-lar üçün qovluq
- go run — compile + birbaşa icra
- go build -o — binary faylı yarat
- Compiler — mənbə → machine code tərcüməçisi

## Praktik nəticə
(1) Çoxnəticəli funksiya: nəticələr MÖTƏRİZƏDƏ `(float64, error)`; return-də
hamısı verilməlidir. (2) Compiler xətalarını QORXMA oxu: have/want sənin
düzəliş planındır. (3) Valid davranış: `return a / b, nil`. (4) Invalid
testdə data lazım deyilsə `_` işlət; yalnız err != nil yoxla. (5) Error
qaytararkən data üçün zero value qaytar. (6) Workflow: red (fail) → green
(pass) → refactor → COMMIT; mürəkkəbliyi sub-behaviour-lara parçala. (7)
Float müqayisəsində == ETMƏ: closeEnough + tolerance (tətbiqə görə tune). (8)
Sqrt tərzü məşq: valid + invalid testləri, standart paket istifadəsi. (9)
Executable = package main + main() + cmd/ strukturu. (10) go run = dərhal
yoxlama; go build -o = paylanan binary.

## Mənbə
Pages: 40-49 (PDF 41-50)
