# Chapter 1 — Testing times (Test Vaxtları)

## Bu fəsil nədən bəhs edir?

İşçi günü simulyasiyası: Texio Instronics-də Go kalkulyator layihəsi. Go mühitinin
qurulması: go mod init (module, go.mod), calculator.go + calculator_test.go,
go test (PASS çıxışının oxunması), gofmt (-d diff, -w fix, Windows diff problemi,
gofumpt), funksiya sintaksisi (func, ad, parametrlər, nəticə, body), FAIL
çıxışının təhlili (--- FAIL: TestSubtract, fayl:sətir, want/got mesajı),
Subtract bug-ı (b - a əvəzinə a - b), testing paketi (standart kitabxana,
3-cü tərəf lazım deyil), test funksiyasının imzası (_test.go, Test prefiks,
*testing.T, nəticəsiz), test bodisi sətir-sətir (t.Parallel, want dəyişəni,
got :=, müqayisə), want-and-got pattern, if statement (şərt, bool, George Boole,
conditional expression, indentasiya = icra yolu).

## Əsas fikirlər

### 1. Yeni Layihə — go mod init
Hər Go layihəsi: qovluq + go.mod faylı (module adını müəyyən edir).
```bash
cd ~/go/calculator
go mod init calculator
# go: creating new go.mod: module calculator
```
go.mod məzmunu:
```
module calculator
go 1.19
```
**Module qaydaları:** modul = kod vahidi; hər modul ayrı qovluqda; modul
daxilində modul OLMAZ. Bütün kitab boyu eyni addımlar: qovluq yarat →
`go mod init MODULE_NAME`.

### 2. İlk Kod — calculator.go
```go
// Package calculator does simple calculations.
package calculator

// Add takes two numbers and returns the result of adding
// them together.
func Add(a, b float64) float64 {
    return a + b
}
```
Sənəd şərhleri (//) — doc konvensiyası: "Add takes..." funksiyanın müqaviləsini
izah edir.

### 3. İlk Test — calculator_test.go
```go
package calculator_test

import (
    "calculator"
    "testing"
)

func TestAdd(t *testing.T) {
    t.Parallel()
    var want float64 = 4
    got := calculator.Add(2, 2)
    if want != got {
        t.Errorf("want %f, got %f", want, got)
    }
}
```
**Vacib:** test faylı `package calculator_TEST` — external test paketi (öz
paketindən FƏRLİ ad) — calculator-u İMPORT edir, sanki istifadəçiyi.

### 4. go test Çıxışı
```bash
go test
# PASS
# ok  calculator  0.234s
```
"ok" = paket keçdi; vaxt maşınadan asılıdır. Editor dəstəyi: VS Code-da
hər testin üstündə "Run test" linki, fayl üstündə "Run package tests".

Xəta olsa: Google + https://go.dev/learn/ — hər kəs əvvəlcə bunu yaşayır.

### 5. gofmt — Standart Format
```bash
gofmt -d calculator.go   # diff göstər: - çıxarılır, + əlavə olunur
gofmt -w calculator.go   # faylı YERINDƏ düzəlt
```
Nümunə: `func \n Add(...)` (artıq boşluq) → `func Add(...)`. Təkrar `-d` —
output YOX = format DÜZGÜN.

**Windows xatırlatması:** `gofmt -d` PATH-də `diff` proqramı tələb edir —
yoxdursa "computing diff: exec: diff: executable file not found" — ciddi problem
DEYİL, yalnız -d işləmir. Editoru save-zamanı gofmt-ə təyin et.

**gofumpt** — gofmt + bir az daha çox reformat.

### 6. Funksiya Sintaksisi
```go
func Add(a, b float64) float64 {
    return a + b
}
```
- `func` — "funksiya gəlir!"
- `Add` — ad
- `(a, b float64)` — parametrlər (adlar + ümumi tip)
- `float64` (mötərizədən sonra) — nəticə tipi (adsız)
- `{ ... }` — body: statementlər
- `return a + b` — funksiyanı bitirir, dəyəri qaytarır

### 7. Failing Test — Təhlil Sənəti
Subtract bug-lı versiyası:
```go
func Subtract(a, b float64) float64 {
    return b - a        // BUG: tərs!
}
```
go test çıxışı:
```
--- FAIL: TestSubtract (0.00s)
calculator_test.go:22: want 2.000000, got -2.000000
FAIL
exit status 1
```
**Oxu addımları:** (1) FAIL var; (2) HANSI test: TestSubtract; (3) HARDA:
calculator_test.go:22 (= t.Errorf sətri); (4) NİYƏ: want 2, got -2 →
"4-dən 2-ni çıxmaq əvəzinə, 2-dən 4-ü çıxıb!" → kodda `b - a` görürsən →
düzəlt.

**Dərs:** test düzgün yazılsa, fail mesajı SƏNƏ bug-u İZAH edir. Bu kitabın
məqsədi Subtract düzəltmək deyil — TESTLƏRİN necə kömək etdiyini öyrənməkdir.

### 8. testing Paketi
```go
import (
    "calculator"     // test edilən paket
    "testing"        // standart kitabxana
)
```
Go-da test üçün BUILT-IN dəstək var: framework, 3-cü tərəf paket LAZIM DEYİL.
(var, amma standart testing ilə başla). Sənəd: https://pkg.go.dev/testing

### 9. Test Funksiyasının İmzası — 3 Tələb
1. **Fayl adı `_test.go` ilə bitir** (bir faylda çox test ola bilər)
2. **Ad `Test` ilə başlayır** (yoxsa Go test kimi tanımır)
3. **`*testing.T` parametri qəbul edir** — test icrasına nəzarət "pultu"
4. **Nəticə qaytarmır** (mötərizə ilə { arasında boşluq = nəticə listi yoxdur)

### 10. Test Bodisi — Sətir-Sətir
```go
t.Parallel()                          // 1. paralel icra — vaxt qənaəti
var want float64 = 2                  // 2. gözlənilən dəyər
got := calculator.Subtract(4, 2)      // 3. funksiyanı çağır, nəticəni saxla
if want != got {                      // 4. müqayisə
    t.Errorf("want %f, got %f", want, got)   // 5. fərq varsa FAIL + mesaj
}
```
**Want-and-got pattern:** want = düzgün cavab; got = funksiyanın
qaytardığı; müqayisə — eynidirsə pass, fərqlidirsə fail. Bu, kitabın ən
fundamental test strukturu.

**t.Errorf:** testi mesajla FAIL edir; çağırılmasa test DEFAULT olaraq KEÇİR.
`%f` — float format.

### 11. if Statement — Şərtli İcra
```go
if want != got {
    t.Errorf(...)
}
```
- `if` + şərt ifadəsi (true/false qiymətlənir) + `{ blok }`
- true → blok İCRA OLUNUR (zigzag sağa); false → blok ötürülür, növbəti sətir
- **bool tipi** — George Boole şərəfinə: yalnız `true` / `false` dəyərləri
- `!=` — "fərqlidir" müqayisə operatoru
- **Conditional expression** — if-dən sonrakı ifadə; nəticəsi bool OLMALIDIR
- İndentasiya = icra yolunun vizualizasiyası: düz xətt aşağı, şərt true olsa
  DAXİLƏ (sağa) dönür

## Əsas terminlələr
- Module — Go kod vahidi; go.mod ilə müəyyən olunur
- go mod init MODULE_NAME — yeni modul yaradır
- go test — paket testlərini işə salır
- gofmt (-d / -w) — format yoxla / yerində düzəlt
- gofumpt — daha streng gofmt variantı
- Funksiya declaration — func + ad + parametrlər + nəticə + body
- Parameter / Result list — (a, b float64) float64
- testing paketi — standart kitabxana test dəstəyi
- _test.go suffiksi — test faylını tanıyan pattern
- Test prefiksi — funksiya-adı tələbi
- *testing.T — testə nəzarət obyekti
- t.Parallel — paralel icra
- t.Errorf — fail + formatlı mesaj
- Want and Got pattern — gözlənti/nəticə müqayisə strukturu
- if statement — şərtli icra
- bool / true / false — Boole məntiqi dəyərləri
- Conditional expression — bool qiymətlənən if ifadəsi
- George Boole — bool adının mənşəyi
- External test package — calculator_test (paket + "_test")

## Praktik nəticə
(1) Yeni layihə: qovluq + `go mod init ad` — go.mod yaranır; modul içində
modul yox. (2) Test faylı ayrı paketdə: `package X_test` + import X + testing.
(3) Test imzası: _test.go faylı + Test prefiks + *testing.T + nəticəsiz.
(4) Bütün hallarda want/got strukturu: want dəyişəni → çağırış → müqayisə →
t.Errorf("%f"). (5) Fail çıxışını SİSTEMATİK oxu: test adı → fayl:sətir →
want/got → diaqnoz. (6) t.Parallel() standart prelüddir. (7) gofmt -w ilə
formatı avtomatik saxla; editorda save-zamanı gofmt qoş. (8) Windows-da
gofmt -d diff tələb edir — yoxdursa yalnız -d işləmir, -w işləyir. (9) if =
şərt (bool) + blok; şərt true → blok icra; indentasiya icra yolunu göstərir.
(10) Test keçmək DEFAULT-dur — səhv heç nə yoxlamayan test "yaşıl" görünə
bilər; yoxlama want≠got şərtindədir.

## Mənbə
Pages: 15-27 (PDF 16-28)
