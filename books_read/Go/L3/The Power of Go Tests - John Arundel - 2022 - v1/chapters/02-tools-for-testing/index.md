# Chapter 2 — Tools for testing (Test Alətləri)

## Bu fəsil nədən bəhs edir?

testing paketinin əsasları (_test.go, Test* funksiya imzası, package_test), test
çıxışının interpretasiyası, "magic package" yanaşması (mövcud olmayan funksiyaları
çağıraraq dizayn), red-green-refactor tsikli, t.Parallel, t.Error/Errorf, t.Fatal/
Fatalf, t.Log, -v/-run/-count flagləri, t.Helper, t.TempDir/t.Cleanup, feeble test
anlayışı ("hansı səhv implementasiyalar hələ keçir?") və cmp.Equal/Diff.

## Əsas fikirlər

### 1. Test Faylının Anatomiyası
```go
// service_test.go:
package service_test          // _test paketi — test kodu OLDUĞU bilinir

import "testing"

func TestAlwaysPasses(t *testing.T) { }   // Test prefiks + *testing.T MÜTLƏQ
```
**Konvensiyalar:**
- Fayl adı `_test.go` ilə bitməlidir — Go YALNIZ bu fayllarda test AXTARIR
- Funksiya adı `Test` prefiksi ilə başlamalıdır
- Parametr `(t *testing.T)` — imza yanlışsa compile xətası:
  `wrong signature for TestX, must be: func TestX(t *testing.T)`
- Build üçün test kodu İGNORLANIR — binary-a 1 bayt belə əlavə olunmur

**Boş test = həmişə PASS:** "No news is good news" — fail deməzsən, keçir.

### 2. Test Çıxışının Oxunması
```bash
go test
--- FAIL: TestAlwaysFails (0.00s)      # test adı + müddət
    service_test.go:8: oh no            # fayl:sətir + MESAJ (IDE-də clickable)
FAIL
exit status 1
FAIL    service 0.292s                    # package + module FAILED
```
Hər hansı test fail → package fail → module fail. Sətir nömrəsi = t.Error çağırışının
YERİ.

### 3. "Magic Package" — Mövcud Olmayanları Çağırmaq
**İdeya:** testi yazarkən funksiyaları VAR SAY — çağırma kodunu İDEAL şəkildə yaz,
sonra paketi implementasiya et. Dizayn problemi → kodiş problemə çevrilir.

```go
func TestRunningIsTrueWhenServiceIsRunning(t *testing.T) {
    t.Parallel()
    service.Start()          // HƏLƏ YAZILMAYIB!
    if !service.Running() {
        t.Error(false)
    }
}
```
**IDE psixoloji təzyiqi:** red squiggle = "səhv etdin" refleksi — amma test-first-də
bu NORMALDIR. "Hey, I'm working test-first over here! Quit red-lining my stuff."

**Bill Kennedy:** "Half the test is to validate the API is intuitive and simple to
use." — öz API-ni istifadəçi kimi sına.

### 4. Red-Green-Refactor
1. **RED:** davranış üçün test yaz → compile edənə qədər MİNİMUM stub
   (`func Running() bool { return false }` — "shameless red") → fail GÖZLƏ
2. **GREEN:** shameless üsulla pass et
3. **REFACTOR:** pass-in təhlükəsiz toyuğunda gözəlləşdir

**"Being the test" thought experiment:** run etməzdən ƏVVƏL proqnoz ver — hansı
sətirdə, hansı mesajla fail edəcək? Proqnoz ≠ nəticə → modelində SƏHV var, DAYAN
və düşün.

### 5. t.Parallel — Hamısı Paralel
```go
func TestX(t *testing.T) {
    t.Parallel()        // hər testdə — "All tests should be parallel"
    ...
}
```
Çoxnüvəli maşında test sürəti ARTIR; feedback = testin MAQSADI olduğundan mütləq
tətbiq et (yaxşı səbəb olmadıqda).

### 6. t.Error vs t.Errorf
```go
t.Error("want", want, "got", got)      // istənilən sayda istənilən tip
t.Errorf("want %q, got %q", want, got) // format string (%q boş string-i görünür edir!)
```
**Error FAIL edir amma DAVAM edir** — bir neçə uğursuzluq ŞABLONU özü ipucu verir:
"not x / not z" → hansı kombinasiya barədə məlumat.

### 7. t.Fatal — Testi Tərk Et
```go
f, err := os.Open("testdata/input")
if err != nil {
    t.Fatal(err)          # fail + DƏRHAL DAYAN
}
```
**Seçim qaydası:** "Davam etməyin faydası varmı?" — fixture qurulmadısa → Fatal
(f == nil ilə davam = konfuz xətalar). Məntiqi yoxlamalar → Error (davam).

### 8. t.Log — Debug Çıxışı
```go
got := StageOne()
t.Log("StageOne result", got)     # yalnız FAIL halında görünür!
got = StageTwo(got)
t.Log("StageTwo result", got)
```
**Niyə fmt.Println DEYİL:**
1. Pass-də də çap olunur → "Tests should pass silently and fail loudly"
2. Paralel testlərdə SIRA qarışır; t.Log avtomatik test nəticəsinə bağlanır
3. Yalnız FAIL edən testlərin logları görünür — istədiyin qədər Log yaz, spam YOX

### 9. go test Flag-ləri
```bash
go test ./...                      # recursive bütün paketlər
go test -v                         # hər testin adı + nəticəsi + LOGLAR
go test -count=1 .                 # cache-i ötür (dəyişiklik olmasa da icra)
go test -run TestRunningIsTrue...  # TƏK test (refaktorinq zamanı)
go test -run TestDatabase          # regex — qrup testlər
```
- `-v` HƏMİŞƏ istifadə etmə — "verbose tests are tedious company"; ekran dolu =
  önəmli xətanı görməmək riski
- `ok service (cached)` — dəyişməyən paketlər keşlənir (sürət)

### 10. t.Helper — Yardımçı Funksiyalar
```go
func createTestUser(t *testing.T, user, pass string) {
    t.Helper()                    # helper kimi işarələ
    ... // setup
    if err != nil {
        t.Fatal(err)              # xəta HELPER-in yox, TEST-in sətirində göstərilir
    }
}

func TestUserCanLogin(t *testing.T) {
    createTestUser(t, "Jo Schmo", "dummy password")   # setup GİZLƏDİLİR
    ... // əSAS MƏNTİQ görünür
}
```
**Fayda:** (1) setup mexanikası gizlənir, testin HEKAYƏSI görünür; (2) xəta
testdə göstərilir (helperdə YOX); (3) `t` ötürməklə helper ÖZÜ fail edə bilər —
error return + yoxlama clutter-i aradan qalxır.

### 11. t.TempDir və t.Cleanup
```go
f, err := os.Create(t.TempDir()+"/result.txt")
# - hər test üçün UNİKAL qovluq (bir-birinə mane YOX)
# - test bitəndə AVTOMATİK silinir

res := someResource()
t.Cleanup(func() {          # test SONUNDA çağırılır
    res.GracefulShutdown()
    res.Close()
})
```
**Cleanup > defer:** helper-də `defer db.Close()` — helper QAYTARANDA bağlanar (db
hələ lazımdır!); t.Cleanup — test bitəndə. Helper-dən resurs qaytararkən MÜTLƏQ
Cleanup.

### 12. Feeble Test — Kifayətsiz Testlər
**Diaqnostik sual: "What incorrect implementations would still pass this test?"**

**Kitabdan nümunə — 3 səviyyəli təkmilləşmə:**
```go
// SƏVİYYƏ 0 — FEEBLE:
func TestNewThing(t *testing.T) {
    _, err := thing.NewThing(1, 2, 3)    // nəticə _ İLƏ UDULUB!
    if err != nil { t.Fatal(err) }
}
// "NewThing doesn't return an error → it works" — ETİBARSAZ DEDUKSİYA
```
```go
// Səhv implementasiya 1: return nil, nil — hələ KEÇİR!

// SƏVİYYƏ 1: nəticəni yoxla:
got, err := thing.NewThing(1, 2, 3)
if got == nil { t.Error("want non-nil *Thing, got nil") }
// Səhv implementasiya 2: return &Thing{}, nil — hələ KEÇİR!

// SƏVİYYƏ 2 — INPUT-UN TƏSİRİNİ yoxla (DÜZGÜN):
x, y, z := 1, 2, 3
got, err := thing.NewThing(x, y, z)
if got.X != x { t.Errorf("want X: %v, got X: %v", x, got.X) }
if got.Y != y { ... }
if got.Z != z { ... }
// Artıq null implementasiya FAIL edir: "want X: 1, got X: 0"
```
**Müqayisə (test qəbulu):** parametrlər nəticəyə TƏSİR ETMƏLİDİR — yoxsa niyə
parametrdir? "Tests that are designed to confirm a prevailing theory tend not to
reveal new information" (Michael Bolton).

### 13. cmp.Equal/cmp.Diff — Dərin Müqayisə
```go
want := &thing.Thing{X: x, Y: y, Z: z}
got, err := thing.NewThing(x, y, z)
if !cmp.Equal(want, got) {
    t.Error(cmp.Diff(want, got))      // YALNIZ fərqli sahələr
}
```
Sahə-sahə manual yoxlamaya ÜSTÜN: qısa, bütöv (yeni sahə avtomatik əhatə olunur),
diff çıxışı anlayışlı:
```
&thing.Thing{
-   X: 1,
+   X: 0,        # istənilən vs alınan
    ...
}
```

## Əsas terminlələr
- `_test.go` — test faylı konvensiyası; build-də istisna
- `package foo_test` — xarici test paketi
- `Test*(t *testing.T)` — test funksiya imzası
- Red-Green-Refactor — fail göstər → pass et → təkmilləşdir
- Shameless Red — compile üçün minimum stub
- Magic Package — ideal API-ni İMAJİNASİYA ilə çağır, sonra implementasiya et
- Being the Test — run-dan əvvəl nəticə PROQNOZU
- t.Parallel — paralel icra; "all tests should be parallel"
- t.Error/Errorf — fail + DAVAM; çoxlu xəta şablonu
- t.Fatal/Fatalf — fail + DƏRHAL DAYAN; fixture xətaları
- t.Log — yalnız fail-də görünən debug
- `-v` / `-run` regex / `-count=1` — verbose / seçici icra / cache ötürmə
- (cached) — dəyişməyən paketin keşlənmiş nəticəsi
- t.Helper — helperdə xətanı test sətirinə yönləndirir
- t.TempDir — unikal + avtomatik silinən müvəqqəti qovluq
- t.Cleanup — test sonunda icra (defer-dən fərqli: helper üçün ideal)
- Feeble Test — "göründüyünden azını test edən" test
- Null Implementation — `return nil, nil` kimi boş həll
- cmp.Equal/cmp.Diff — dərin bərabərlik + fərq raportu
- "%q" verb — boş string-ləri görünür edən format

## Praktik nətidə

(1) Test yazmaq = 3 qayda: _test.go + Test prefiks + t *testing.T. Heç bir
framework YOX — built-in. (2) Boş test PASS deməkdir — testin dəyəri onu FAIL
edə biləcəyi şərtlərdədir. (3) Testi fail görməmiş ona etibar etmə; "being the
test" proqnozu model səhvlərini ERKƏN tutur. (4) Magic package: IDE red-line =
test-first-in normal vəziyyəti; ideal çağırış kodunu yaz, sonra implementasiya.
(5) hər testə t.Parallel(). (6) Error vs Fatal: davamın faydası sualı. (7) t.Log
- fmt.Println əvəzinə — pass SAKİT, fail SƏS-Lİ. (8) -v həmişə yox — ancaq yeni
sistem öyrənərkən. (9) Setup-i helper-ə: t.Helper + t ötür → test HEKAYƏTİ
açıq qalır. (10) TempDir/Cleanup — fayl və resurs testlərinin təmizlənməsi üçün
avtomatik; helper-də defer YOX, Cleanup. (11) Hər testə "hansı SƏHV kod keçər?"
sualını ver — null implementasiya keçirsə test FEEBLE-dir. (12) Böyük struct-ları
cmp.Equal + Diff ilə — manual sahə yoxlamasından qısa və TAM.

## Mənbə
Pages: 29-60 (PDF 41-71)
