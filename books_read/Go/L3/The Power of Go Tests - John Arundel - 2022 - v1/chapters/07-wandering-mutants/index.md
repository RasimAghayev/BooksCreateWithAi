# Chapter 7 — Wandering mutants (Gəzən Mutantlar)

## Bu fəsil nədən bəhs edir?

Test coverage (go test -cover, -coverprofile, go tool cover -html), coverage
instrumentasiyası (compiler əlavə instruksiyalar — fuzz "interesting" ilə eyni
mexanizm), Goodhart qanunu / Cobra Effect (raket effecti), coverage = siqnal,
hədəf YOX; TestAdd(2,2) + Add→42 = 100% coverage, 0% test; bebugging (bilərəkdən
bug əkmək — != → ==), feeble test aşkarı, lazımsız/unreachable kodun tapılması;
mutation testing (avtomatik bebugging): go-mutesting, mutation score, IsEven
nümunəsi — memoizasiya + sync.Mutex, 4 unsurvivable mutation, təkrar input
testdə YOX idi; diff oxunuşu, score 1.0→0.2→düzəliş; gremlins aləti.

## Əsas fikirlər

### 1. Coverage Nədir
"Executing the test suite causes 80% of the package's source statements to
run → coverage 80%." (Rob Pike, "The cover story")

**Amma:** goal statement-ləri İCRA ETMƏK deyil — davranışı TAPMAQ. Yəni:
```bash
go test -cover
# coverage: 85.7% of statements

go test -coverprofile=coverage.out
go tool cover -html=coverage.out   # browser: YAŞIL=covered, QIRMIZI=uncovered, BOZ=irrelevant
```
**Qırmızı sətirlərə bax və qərar ver:** test lazımdır, yoxsa yox?
`if err != nil { return err }` — özünə görə testə ehtiyac YOXUDUR: "We don't
need to test Go's return statement: we know it works."

**Coverage mexanizmi:** coverage rejimində compiler hər statement-dən əvvəl/sonra
əlavə instruksiyalar yerləşdirir — icra edilib-edilmədiyini qeyd edir. Bu məlumat
həm coverage profile, həm də fuzz-in "interesting path" detektoru üçün istifadə
olunur (eyni alət!).

**IDE:** VS Code — "Go: Toggle Test Coverage in Current Package"; GoLand —
editor daxilində highlight.

### 2. 100%-ə Yaxınlaşdıqca Xəzər Azalır
85% → 100% üçün lazım olan testlər adətən BUNA DƏYMƏZ. Testlər PULSUE DEYİL:
zaman, enerji, texniki borc — illər boyu oxuyub saxlayacaqsan. İnvestisiyanın
qarşılığını verməli dir. (Donovan & Kernighan: "trade-off between the cost of
writing tests and the cost of failures")

### 3. Coverage = Siqnal, Hədəf YOX (Goodhart + Cobra Effect)
**Goodhart qanunu:** ölçü hədəfə çevriləndə yaxşı ölçü olmaqdan çıxır — adamlar
rəqəmi şişirtmək üçün mənasız işlər görər.

**Cobra Effect:** əks-məhsuldar stimullar → istənməyən nəticələr. Zəhərli
ilanlara mükafat qoyubsan — adamlar İLAN YETİŞDİRMƏYƏ başlayar. "Bug bounty"
sxemləri də bəzən belə.

**Manager tələbi "minimum X% coverage"** = əks-məhsuldardır: 100 sətirlik
heç-nə-etməyən funksiya yazıb testdən çağırmaq coverage-i ARTIRIR, amma sistemə
heç nə qatmir.

### 4. Add Nümunəsi — 100% Coverage, 0% Test
```go
func TestAdd(t *testing.T) {
    t.Parallel()
    add.Add(2, 2)              // nəticə YOXLANMIR
}

func Add(x, y int) int {
    return 42                  // WRONGER THAN WRONG — amma test PASS
}
```
`go test -cover` → `PASS, coverage: 100.0%`. **Test coverage ≠ test quality.**
Bu test YALNIZ "Add panic etmir" deyir (bəlkə "better than nothing, but not
by much").

**Coverage-in əsl qiyməti:**
1. Bəzən vacib davranışların testdən kənar qaldığını göstərir
2. Kodun SƏNİN FİKİRLƏŞDİYİNDƏN fərqli davrandığını öyrədir (Tim Bray:
   "Every time I do this I get surprises — huge gaps in the coverage")
3. **Coverage ratchet:** "no check-in is allowed to make the coverage numbers
   go down" — aşağı salan commit qadağan. Ya untested kod əlavə etdin, ya lazımsız
   test sildin. Aşağı coverage-li layihədə ən azı pisləşməni dayandırır.

**Donovan & Kernighan sitatı:** "Just because a statement is executed does not
mean it is bug-free." — İCRA ≠ DÜZGÜN. Birinci: bir şey ETMƏLİ. İkinci (daha
vacib): DOĞRU şeyi etməli.

### 5. Bebugging — Feeble Test Ovu (Weinberg)
**Feeble test** — kodu COVER edir, amma kifayət qədər TEST ETMİR (çağırır, amma
davranışı yoxlamır).

**Bebugging (seeding):** bilərəkdən, testçilərə DEMƏDƏN məlum bug-lar ək →
tapılan % = qalan naməlum bug-ların təxmini. (Gerald Weinberg)

**Manual bebugging:**
```go
if err != nil { return err }    // → flip: !=  →  ==
if err == nil { return err }    // SINTAKSİSƏN DÜZGÜN, MƏNTİQƏN XƏTA
```
Testlər FAIL edirsə — coverage yaxşıdır. FAIL ETMİRSƏ → ya feeble test (DÜZƏLT!),
ya da lazımsız kod (SİL!). Hər iki halda qələbə.

### 6. Bebugging = Unreachable Kod Axtarıcı
`if false { return 1 }` — aşkar unreachable. Amma gizli variantlar da var:
```go
if s.TLSEnabled {
    if err := s.initiateTLS(); err != nil {
        return err
    }
}
```
Bebug → test FAIL ETMİR → `s.TLSEnabled` heç vaxt true OLMUR (test istifadəsində
əlçatmaz kod). Nəticə: (a) vacibdirsə TEST YAZ (flip FAIL etməlidir); (b)
istifadəçilər maraqlanmırsa — SETTİNGLƏ BİRLİKDƏ SİL.

**Bebuggin ikili xəbəri:** feeble test YAXUD lazımsız behaviour — hər ikisi
dəyərli tapıntıdır.

### 7. Mutation Testing = Avtomatik Bebugging
**Necə işləyir:** (1) proqramı syntax tree-yə parse edir; (2) bir kod sətrini
mutasiya edir (məs. == → !=, statement-i blank assign ilə əvəz edir); (3)
testləri İŞLƏDİR. FAIL → adekvat covered, növbəti. PASS (heç bir test tutmadı) →
PROBLEM: ya feeble test, ya lazımsız kod — DIFF kimi report edilir.

**Nəticə:** "changes that don't fail any test" diff-lər toplusu — şərhi SƏN
verirsən: testi gücləndir / kodu sil / false positive kimi yox say.

**False positive nümunəsi:**
```go
fmt.Println("Hello, world")            // → mutant:
_, _ = fmt.Println("Hello, world")      // Println nəsə qaytarır (who knew?)
```
Davranış dəyişməyib → problem YOX.

### 8. go-mutesting Aləti
```bash
go install github.com/avito-tech/go-mutesting/cmd/go-mutesting@latest

# !!! ƏVVƏLCƏ BÜTÜN LAYİHƏNİN BACKUP-I (.git daxil) !!!
# alət source faylları DƏYİŞİR — zərər riski

go-mutesting .
# The mutation score is 1.000000 (1 passed, 0 failed, ...)
```
- **Score 1.0** = bütün mutantlar testlər tərəfindən öldürüldü (survivor YOX)
- **--do-not-remove-tmp-folder** — mutasiya olunmuş faylları saxla (baxmaq üçün)
- Başqa alət: **gremlins** — github.com/go-gremlins/gremlins

### 9. IsEven Nümunəsi — Növbəti Addımlar
**Adım 1 — iki həqiqi test:**
```go
func TestIsEven_IsTrueForEvenNumbers(t *testing.T) {
    t.Parallel()
    for i := 0; i < 100; i += 2 {
        t.Run(strconv.Itoa(i), func(t *testing.T) {
            if !even.IsEven(i) { t.Error(false) }
        })
    }
}

func TestIsEven_IsFalseForOddNumbers(t *testing.T) {
    t.Parallel()
    for i := 1; i < 100; i += 2 {
        t.Run(strconv.Itoa(i), func(t *testing.T) {
            if even.IsEven(i) { t.Error(true) }
        })
    }
}
```
Subtest adı = input dəyərinin ÖZÜ (i) — başqa nə demək lazım deyil, i bu test
case-i TAM təsvir edir. `return n%2 == 0` sadəliyi "mutation tester üçün kifayət
qədər material verməz" — uzun versiya saxlanılır.

**Adım 2 — memoizasiya əlavə et (bilərəkdən bug üçün):**
```go
var cache = map[int]bool{}
func IsEven(n int) (even bool) {
    even = n%2 == 0
    if _, ok := cache[n]; !ok {
        cache[n] = n%2 == 0
        even = n%2 == 0
    }
    return even
}
```
`go test` → **"fatal error: concurrent map read and map write"** — t.Parallel()
goroutinlərindən map-ə yazmaq = DATA RACE. Bir does not simply...
**Düzəliş: sync.Mutex** (m.Lock / defer m.Unlock) → testlər PASS.

**Adım 3 — go-mutesting:** score 0.2 (1 passed, 4 FAILED, total 5) — 4 mutant
SAĞ QALDI. DİQQƏT: burada "passed" = mutant öldürüldü (test fail etdi), "failed"
= mutant sağ qaldı (testlər HEÇ NƏDƏRMƏDİ) — terminologiya tərsinə!

### 10. 4 Mutantın Təhlili — Dərin Dərslər
**Mutant A (cache yazısı söndürüldü):**
```go
-   cache[n] = n%2 == 0
-   even = n%2 == 0
+   _, _, _, _, _ = cache, n, n, even, n
```
Cache-ə YAZILMIR = memoization tamamilə söndürülüb — amma testlər PASS!
→ **Testlər cache DAVRANIŞINI yoxlamır.**

**Mutant B (ilk hesablama söndürüldü):**
```go
-   even = n%2 == 0       # if blokundan ƏVVƏL
+   _, _ = even, n
```
Testlər PASS qalır. Niyə? Cache-də OLMAYAN input: if bloku hesablayıb qoyur —
nəticə eyni. **Amma cache-də OLAN input:** if FALSE → `return even` → bool
default = false → **cache-dən oxunan hər nəticə YALNIS!** Problem: **heç bir
test eyni dəyərlə İKİ DƏFƏ çağırmır!** Geniş input yaratdıq (0..98), amma
TƏKRAR input növünü verməyi unutduq. → Düzəliş: "same input → same result"
testi yaz (= property-based test, ötən fəsil!)

**Mutant C (ikinci hesablama söndürüldü):** eyni işi edən ikinci assignment
silinir — davranış DƏYİŞMİR → Bu mutant HƏQİQƏTƏN də silinə bilər: birinci
hesablama lazımdır (cache hit case), ikincisi ARTIQDİR. Kod qısalır, sadələşir.

**Mutant D (qeyd):** — müvafiq diff-lər (1 passed) təhlil barədə (yuxarıdakı
kitab mətnində göstərilib).

### 11. Düzəldilmiş Versiya + Final Dərslər
```go
func IsEven(n int) (even bool) {
    m.Lock()
    defer m.Unlock()
    even, ok := cache[n]
    if !ok {
        even = n%2 == 0
        cache[n] = even
    }
    return even
}
```
Amma **əsl dərs:** IsEven lazımsızdır! `n%2 == 0` birbaşa yaz, hətta
memoizasiya əbəsdir — cüt/tək yoxlama ən aşağı bitə baxmaqdır, hardware-də
çox sürətli. Sadə kodla mutation testing öyrənmək yaxşıdır, sonra REAL
layihəyə keç.

### 12. Nə Vaxt / Necə İstifadə
- Hər gün, hər həftə YOX — "regular health check": bir neçə ayda bir və ya
  BÖYÜK dəyişikliklərdən sonra
- Böyük layihədə ÇOX nəticə verir, çoxu YARARSIZ — amma "bir-iki REAL bug,
  feeble test və ya lazımsız kod path" demək olar QARANTİYA
- **Score-u qeyd et:** dəyişiklikdən sonra score DÜŞÜBSƏ → untested / lazımsız
  kod əlavə etmisən — review et

## Əsas terminlələr
- Test Coverage — testlərin işə saldığı statement %-i
- -cover / -coverprofile=coverage.out — faiz / profil faylı
- go tool cover -html — yaşıl/qırmızı/boz highlight
- Coverage Instrumentation — compiler-in statement counter-ləri (fuzz ilə ortaq)
- Goodhart's Law — ölçü hədəf olsa, yaxşı ölçü olmağı dayandırır
- Cobra Effect — pərvər edici stimul → əks nəticə (ilan yetişdirmə)
- Coverage Ratchet — coverage-i aşağı salan check-in qadağası
- Bebugging (Seeding) — məlum bug əkib tapılma %-i ölçmə (Weinberg)
- Feeble Test — çağırır, amma davranışı yoxlamır
- Mutation Testing — avtomatik bebugging (syntax tree mutasiyası)
- go-mutesting — avito-tech aləti; gremlins — yeni alternativ
- Mutation Score — 0..1, öldürülən mutant nisbəti (1 passed = mutant ölü)
- Survivable/Survived Mutant — testlərin TUTMADIĞI dəyişiklik
- Memoization — nəticələrin cache-də saxlanması
- Data Race — paralel map read/write = fatal error
- Unreachable Code — heç vaxt icra olunmayan, silinməli kod

## Praktik nəticə
(1) `go test -cover` + `-coverprofile` + `cover -html` — qırmızıları GÖZLƏ
KEÇİR, hər biri üçün qərar ver: test yaz / görüntü ilə düzgündür / sil.
(2) Coverage hədəf YOX, siqnaldır — 85→100 testləri adətən dəyməz; yalnız
AŞAĞI DÜŞÜŞ ratchet ilə dayandır. (3) 100% coverage + 0 test mümkündür
(Add→42) — icra ≠ düzgünlük. (4) Coverage texnologiyası fuzz "interesting"
 detection ilə eyni compiler instruksiyalarıdır — bir alət, iki istifadə.
(5) Yad kodu qiymətləndirirsənsə: testləri oxu, "What are we really testing
here?" soruş — bu ilk həftələrin ƏN DƏYƏRLİ fəaliyyətidir. (6) Bebugging:
!= → == flip et, test FAIL etmirsə — feeble test düzəlt YAXUD kodu sil (hər
ikisi qələbə). (7) Go-mutesting işə salmazdan ƏVVƏL tam backup (.git daxil).
(8) "passed/failed" terminologiyası TƏRSİ nəzərdə tutulur — "1 passed" =
mutant öldürüldü. (9) Cache/memoizasiya testi: eyni input TƏKRAR çağır və
nəticənin EYNİLİKLİYİNİ yoxla — təkrar input unutma. (10) Mutation testing =
aylıq sağlamlıq yoxlaması; score qeyd et, düşüş = warning. (11) Sadə iş üçün
öz helper yazma (IsEven) — `n%2 == 0` birbaşa; mutation tester sadə kodda da
ARTIQ hesablamaları (ikinci assignment) tapır.

## Mənbə
Pages: 186-216 (PDF 198-224)
