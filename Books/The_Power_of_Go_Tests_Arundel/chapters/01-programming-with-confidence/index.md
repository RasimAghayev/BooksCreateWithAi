# Chapter 1 — Programming with confidence (Konfidansla Proqramlaşdırma)

## Bu fəsil nədən bəhs edir?

Self-testing kod konsepti, test-first fəlsəfəsi (niyə əvvəl test), ListItems
funksiyasının tam TDD tsikli (test → yalançı implementasiya → fail yoxlaması →
Shameless Green → table test → panic halları → refaktorinq), "New behaviour, new
test" qaydası, cmp.Diff ilə fərq çapı və 10 dəqiqəlik prodaktivlik müqayisəsi.

## Əsas fikirlər

### 1. Niyə Test Əvvəl?
**Test-first-in əsas səbəbi:** kod yazmazdan ƏVVƏL istifadəçi nöqtei-nəzərindən
davranışı AYDIN düşünmək. Kodu davranışı bilmədən yazmaq = vaxt itkisi +
implementator üçün convenient, istifadəçi üçün YARARSIZ dizayn.

**Gələn faydalar:**
- Kiçik inkrementlər → yanlış yolda uzağa getməmək
- Öz API-ni İSTİFADƏ etmə təcrübəsi → yaxşı dizayn
- Böyük proqramı kiçik, müstəqil, dəqiqliyi göstərilən modullara bölmə
- **Testability ↔ good design sinergiyası** (dərin əlaqə)

### 2. Self-Testing Code (Martin Fowler)
**Tərif:** kod bazasına seriya avtomatik testlər işlədə bilirsənsə və testlər
KEÇİRSƏ — kodun əhəmiyyətli defeckt-lərdən azad olduğuna inana bilirsənsə.
**Kent Beck:** "Tests are the Programmer's Stone, transmuting fear into boredom" —
qorxunu can sıkıntısına çevirir; "No, I didn't break anything. The tests are all
still green."

### 3. İlk Test — ListItems Nümunəsi
**Ssenari:** Zork tipli oyun — "You can see here a battery, a key, and a tourist
map." cümləsi qurulmalı.

**Kitabdan kod nümunəsi:**
```go
func TestListItems(t *testing.T) {
    t.Parallel()
    input := []string{
        "a battery",
        "a key",
        "a tourist map",
    }
    want := "You can see here a battery, a key, and a tourist map."
    got := game.ListItems(input)      // funksiya HƏLƏ YOXDUR!
    if want != got {
        t.Errorf("want %q, got %q", want, got)
    }
}
```
**Vacib müşahidə:** funksiya mövcud olmasa da, test yazmaqla bir sıra DİZAYN
QEQRARLARI verildi — ad (ListItems), paket (game), parametr ([]string), nəticə
(string), DAVRANIŞ (input üçün expect). Bu, komponent haqqında ən çətin
qərarlardır — kod yazmadan qəbul edildi!

**Pure funksiya:** test yazmaq ListItems-i terminaldan AYIRDI — nəticəni PRINT
etmək əvəzinə QAYTARAN funksiya (deterministik, side-effect-siz) oldu. Testability
dizaynı dəyişdirir.

### 4. Testi Yoxlamaq — "Until you've seen the test fail..."
**Problem:** testlər Go-da DEFAULT olaraq KEÇİR (boş test = həmişə pass):
```go
func TestAlwaysPasses(t *testing.T) {}   // HEÇ NƏ yoxlamır!
if want != want { ... }                    // özünlə müqayisə = həmişə false
```
**Həll — YALANÇI implementasiya ilə fail yoxla:**
```go
func ListItems(items []string) string {
    return ""        // qəsdən SƏHV nəticə
}
```
```bash
go test
# --- FAIL: TestListItems
#     want "You can see here a battery...", got ""
```
Fail gözlənildiyi kimi → TEST DÜZGDÜR. **"Until you've seen the test fail as
expected, you don't really have a test."** — testin ÖZÜ də bug-lu ola bilər!

### 5. cmp.Diff — Fərqin Dəqiq Göstərilməsi
```go
import "github.com/google/go-cmp/cmp"

if want != got {
    t.Error(cmp.Diff(want, got))     // yalnız FƏRQİ çap et
}
```
Çıxış formatı:
```
strings.Join({
    "You can see here",
-   " ",                    ← YALNIZ want-da var (çatışmır)
    "a battery, a key,",
-   " and",                  ← YALNIZ want-da (unuduldu!)
    " a tourist map.",
  }, "")
```
`-` = want-da amma got-da YOX; `+` = got-da amma want-da YOX. Uzun nəticələrdə
hansı hissənin səhv olduğunu dərhal göstərir.

### 6. Shameless Green
**Tərif:** testləri ən SADƏ, ən SÜRƏTLİ, ən BAŞA DÜŞÜLƏN yolla keçən kod.
"Correct first, elegant later" — kod hələ işləmirsə onu gözəlləşdirmək mənasızdır:
"The goal right now is not to get the perfect answer but to pass the test."

### 7. "New Behaviour, New Test" — Qızıl Qayda
**Tələ:** ListItems 3 item üçün işləyir; 2/1/0 item istəyirsənsə — ƏVVƏL test
əlavə et, SONRA kod. Yeni DAVRANIŞ = yeni TEST.
**Tələbi aradan qaldırır:** ideyaları kodlamadan ƏVVƏL dəqiqləşməyə məcbur edir;
"nə qədər kod kifayətdir" sualına cavab verir (test pass = bitti).

### 8. Table Test — testCase Struct + range
```go
func TestListItems(t *testing.T) {
    type testCase struct {
        input []string
        want  string
    }
    cases := []testCase{
        {input: []string{"a battery", "a key", "a tourist map"},
         want: "You can see here a battery, a key, and a tourist map."},
        // ... yeni hallar buraya bir-bir əlavə olunur
    }
    for _, tc := range cases {
        got := game.ListItems(tc.input)
        if tc.want != got {
            t.Error(cmp.Diff(tc.want, got))
        }
    }
}
```
**Refaktorinq addımı:** mövcud test CASES struktuna çevrilir → ƏVVƏL pass-ı yoxla
(refaktorinq heç nə pozmur), SONRA case əlavə et.

### 9. Add One Case at a Time
Hər yeni case: (1) əlavə et → (2) FAIL-ı gör → (3) minimal kodla keçir → (4)
növbəti. Hamısını bir dəfəyə əlavə etmək = itirilmiş fərdi feedback.

### 10. Panic-lərin Söndürülməsi
1-item case əlavəsi → `panic: index out of range [1] with length 1` —
kod items[1]-ə müraciət edir, amma 1 elementli slicedə YOXDUR.
**Həll:** panic-in ÖZÜNÜ case kimi işlə — `if len(items) == 1 { return ... }`.
0-item: `index out of range [0] with length 0` → `if len(items) == 0 { return "" }`.

### 11. Refactoring — Test Təhlükəsiz Toyuğu
**Tərif:** davranışı DƏYİŞMƏDƏN kodu dəyişmək. Test bütün RELEVANT davranışı
müəyyən etdiyindən — kodu TAM AZADLIQLA dəyiş, test anında xəbər verəcək.

**Kitabdan final refaktorinq:**
```go
func ListItems(items []string) string {
    switch len(items) {
    case 0:
        return ""
    case 1:
        return "You can see " + items[0] + " here."
    case 2:
        return "You can see here " + items[0] + " and " + items[1] + "."
    default:
        return "You can see here " +
            strings.Join(items[:len(items)-1], ", ") +
            ", and " + items[len(items)-1] + "."
    }
}
```
4 xüsusi hal → switch — daha oxunaqlı; DAVRANIŞ EYNİ (test hələ pass).

### 12. On Dəqiqəlik Təhlil
Test-first ilə: ~2 dəq test yaz + 3 dəq pass + 2 dəq refaktor = **10 dəq, BİTMİŞ**.
Test-siz: 10 dəq ehtiyatlı yazmaq + düşünmək + throwaway kodla yoxlama + "bir
boşluq unutdum" düzəlişləri = EYNİ vaxt, amma İNAMSIZ nəticə. "When we're done,
we're really done" — görünən yavaşlıq İLLÜZİYADIR.

**Kernighan:** "If you're as clever as you can be when you write it, how will you
ever debug it?" — testlər sadəliyi MƏCBUR edir, bu BÜTÜN proqramçılar üçün yaxşıdır.

## Əsas terminlələr
- Self-Testing Code — testlər pass = kod əhəmiyyətli defekt-siz (Fowler)
- Test-First — davranışı kod yazmazdan əvvəl dəqiqləşdirmə
- Pure Function — deterministik, input-dan asılı, side-effect-siz
- t.Parallel() — paralel test işarəsi
- t.Errorf — test fail ETDİ amma DAVAM edir
- Shameless Green — testi keçən ən sadə kod
- cmp.Diff — want/got fərqini and-ırıq şəkildə göstərən üçüncü tərəf alət
- Table Test — testCase struct + cases slice + range loop
- New Behaviour, New Test — yeni davranış üçün mütləq yeni test
- Red-Green-Refactor — fail göstər → pass et → təkmilləşdir
- Panic Quelling — index-out-of-range xüsusi hallarının kod ilə idarəsi
- Refactoring — davranışı dəyişmədən kod struktununun dəyişdirilməsi

## Praktik nətidə

(1) Test YALNIZ görmək istədiyin DAVRANIŞI ifadə edir — implementasiya detallarını
YOX. (2) Testi yoxlamadan implementasiya ETMƏ: boş test = həmişə pass; `want !=
want` = heç vaxt fail. (3) Funksiyanı YALANÇI (return "") ilə yaz → fail göstər →
ancaq sonra implementasiya. (4) cmp.Diff uzun string-lərdə FƏRQİ bir gözərtmə ilə
göstərir. (5) Shameless Green: sadə, çalış, anlaşılan — gözəllik SONRA. (6) Yeni
davranış fikri gəldi → test əlavə et → fail → minimal həll → təkrarla. (7) Table
test strukturu: refaktordan SONRA pass → sonra case-ləri bir-bir əlavə. (8) Panic =
XÜSUSİ HAL YOXdur: case əlavə et, şərt yaz. (9) Test pass-dan sonra KOD AZADDIR —
switch/refaktor without fear. (10) Test-first YAVAŞ DEYİL — 10 dəqiqə = bitmiş,
inamlı, self-testing kod.

## Mənbə
Pages: 1-28 (PDF 12-40)
