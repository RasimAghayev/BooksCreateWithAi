# Chapter 11 — Suite smells (Suite Qoxuları)

## Bu fəsil nədən bəhs edir?

Bütün test SUITE-inin sağlamlıq yoxlaması: (1) No tests — "guerilla testing",
icazə lazım deyil (Robert C. Martin, Tarlinder, Mike Bland/Goto Fail), testing
dini deyil (Tim Bray), "don't be the person who fights for an idea"; (2) Legacy
code — Feathers-in adalar metaforu, test-what-you-touch, Edit and Pray vs Cover
and Modify, untestable üçün rescue proseduru (yeni funksiya + test + köhnəni
sil), Ousterhout complexity sitatı; (3) Insufficient tests — "donate blood"
(Leopold), Feathers-in team experiment-i, tests-first/last mübahisəsi (dogma
yox, nəticə hesabına); (4) Code review — 10-addımlıq TEST-FOKUSLU checklist;
(5) Optimistic tests — precondition yoxluğu (Create mirage, Delete loophole,
"nuke the whole DB" WHERE-omission bug-ı, Bob-un sağqalması, createUserOrFail
helper); (6) Persnickety tests — "halt and catch fire" (Weinberg), Poincaré,
seçici assert (GOOS, Donovan&Kernighan), fingerprint property testi (MD5→SHA256
refactor), impossible-ı test etmə; (7) Over-precise — float epsilon,
cmpopts.EquateApprox, epsilon dəqiqləşdirmə metodu; (8) Too many tests —
"lines spent" (Dijkstra), bağça budaması, scaffolding, delta coverage (Herb
Derby, Kent Beck), Fowler-in balansı; (9) Test frameworks — testify/ginkgo,
assert.Equal "premature abstraction", import "testing"; (10) Flaky tests —
timing, map sırası, cmp.Equal + cmpopts.SortSlices, "thou shalt not suffer a
flaky test to live", brittle ≠ flaky; (11) Shared worlds — private world,
take-nothing-but-pictures, test-adlı DB, transaction rollback, "tests must be
safe to run against prod, because one day they will"; (12) Failing tests —
zero defects, broken windows, infinite defects (Word death march), bug
bankruptcy/amnesty (Weinberg); (13) Slow tests — Beck time (10 dəq), 7 sürət
resepti; Kent Beck-in son sitatı (shower).

## Əsas fikirlər

### 1. No Tests — Guerilla Testing
Startup fazasında speed > maintainability → nəticə: kod toxunulmaz olur, hamı
bəhanə axtarır. Həll: **icazə İSTƏMİ, sadəcə yaz.** "You don't need to convince
the customer... You just need to convince yourself." (Uncle Bob)

- Reviewer xəbərdarlığı: "guerilla testing" işdən qovura bilər — müəllif:
  yaxşı olsa, yaxşı kompaniyaya yönləndirər (öz missiyası!)
- Testing dini deyil — convert lazım deyil: "be the person who has ideas worth
  fighting for". Yaxşı işləyən yolu göstər, hamı özü gələcək (Emily Bache:
  freedom to change).
- Mike Bland (Goto Fail/Heartbleed): kompaniyalar code quality üçün ödəmir —
  "onların problemi, mənim marağım SƏNİN işin yaxşı getməsidir".

### 2. Legacy Code — Adalar Strategiyası
"Test what you touch": 1337-ci sətir üçün monster-funksiyanı test etmə — YOX.
Kiçik bloku FUNGSIYAYA çıxart → kiçik test yaz → dəyişikliyi et → qalanına
toxunma. Təkrar-təkrar → "tested islands rise out of the ocean... eventually
continents of test-covered code" (Feathers). Səy avtomatik DÜZ GEDƏN yerə düşür
— dəyişdirilən hissələrə. Kent Beck: tez-toxunulan hissələr qaya kimi sərt,
periferiya ləkəli ola bilər.

**Legacy tərifi (Tarlinder):** testsiz + testability-ni nəzərə almadan
dizayn olunmuş — temporal coupling, indirect input, monster methods,
"bizarre web of dependencies".

**Edit and Pray vs Cover and Modify** (Feathers) — birincisi industry
standarddır, amma təhlükəli.

**Untestable rescue:** köhnə funksiyaya test yaza BİLMƏZSƏN (Catch-22) →
1) İSTƏDİYİN funksiya üçün test yaz; 2) scratch-dən yeni funksiya yaz, test
keçsin; 3) köhnəni SİL. (Müəllifin uğurla işlədiyi emergency prosedur.)

**Ousterhout:** complexity minlərlə kiçik səhvlərin cəmidir — baş verərkən
görmüzsən, sonra düzəltmək mümkünsüz. **Legacy haqqında müsbət:** işləyir,
buglar məlum, workarounds var. Əsas qayda: YENİ legacy YAZMA.

### 3. Insufficient Tests
Sürət bahasına keyfiyyət: "You don't save time... It's like: 'I need to lose a
couple of pounds' — 'Yeah. Donate blood.'" (Leopold)

**Feathers-in eksperimenti:** bir iterasiya boyu HƏR dəyişikliyin testi olsun;
çağra bilməyən komanda iclası çağırsın. Əvvəl dəhşətli görünür → sonra
dəyişikliklər asanlaşır. Yarısı özünə inamdır: HƏR ŞEYİ test edə biləcəyinə.

**Tests-first vs tests-last:** dogma YOX — "If your engineers are producing
code with effective tests, don't be giving them any static about how it got
that way." (Tim Bray). Müəllif tests-first öyrədir, amma "moral issue deyil" —
nəticəyə bax. Amma kod pis + testlər feeble/optimistic/missing → tests-last
səbəb ola bilər → test-first keçidi nəticəni yaxşılaşdıra bilər.

### 4. Code Review — 10 Addımlıq Checklist (TEST-FOKUS)
1. Dəyişiklik NİYƏ? Biznes tələbi/bug bilmədən düzgünlüyü HANSI ölçü ilə?
2. Hər tələb üçün TEST varmı? Həqiqətən test edirmi, yoxsa iddia?
3. Testlər FRESH CHECKOUT-dan keçirmi? (gizli fayl/binary/env asılılığı!)
4. Testlər AÇIQ bug-ları TUTURMU? — bug ək (mutation testing!), tutmursa
   approve ETMƏ: "at least the bugs you thought of"
5. Fail mesajları dəqiq/köməkçi/tamdır? — hər fail path-ini İCRA ET, test
   kodunu oxumadan diaqnoz qoyula bilərmı?
6. Arbitrary input qəbul edirsə — FUZZ edilibmi? Boş slice/map/string?
   İnanclı input nümunəsi və diapazonu?
7. Coverage tam? Qalan pathlər gözlə dəqiqləşir? Yoxsa lazımsızdır — sil?
8. Kod testdən ÇOX şey edirmi? (over-engineering) Sadələşdirilsə keçərmi?
9. Testlər testdən çox YOXLAYIRMI? ("while we're here...") Başqa kodun
   düzgünlüyünə çox bağlanıb — brittle?
10. Test kodundan davranış AYDINDIRMI? Müştəriyə göstərsən başa düşərmi?

Ümumi prinsip: review əsasən TESTLƏRƏ fokuslanmalı — line-by-line nitpick
faydasızdır; kod artıq mövcuddur → müdafiə psixologiyası.

### 5. Optimistic Tests — Preconditions + Implicit Postconditions
**Create mirage:**
```go
func TestCreate(t *testing.T) {
    t.Parallel()
    user.Create("Alice")
    if !user.Exists("Alice") { t.Error("Alice not created") }
}
```
Create HEÇ NƏ ETMƏSƏ? Yalnız Alice ƏVVƏLCDƏN varsa keçər! Create-a `func
Create(name string) {}` + seed map-də Alice yazsaq → həmişə PASS. **"Mirage
test": test var kimi görünür, amma YOXDUR.**

**Düzəliş — precondition yoxla:**
```go
if user.Exists("Alice") {
    t.Fatal("Alice unexpectedly exists")     // ƏVVƏLCƏDƏN yoxdur!
}
user.Create("Alice")
if !user.Exists("Alice") { t.Error("Alice not created") }
```
Create-un vacibliyi: dünyanı "Alice yoxdur"dan "Alice var"a DƏYİŞDİRMƏSİDİR.

**Delete loophole:** Create(heçnə) + Delete(heçnə) → test KEÇİR (Alice heç
olmamışdı, Exists false idi bəs). Precondition: Create-dən sonra Alice VAR
yoxla. Müəllif özü bu xətanı edib: stub yerləşdirib unudub — test keçdiyi üçün.

**Implicit postcondition — "nuke the whole DB":** Delete WHERE-siz SQL ilə
HAMINI silə bilər! Test yalnız "Alice silindi" yoxlayır. Həll — İKİ istifadəçi,
birini sil, DİGƏRİ SAĞ QALSIN:
```go
user.Delete("Alice")
if user.Exists("Alice") { t.Error("Alice still exists after delete") }
if !user.Exists("Bob") { t.Error("Bob was unexpectedly deleted") }
```
Paperwork → `createUserOrFail(t, name)` helper (t.Helper()). "It's arguably even
more important to test that it doesn't do what it shouldn't do." (Mike Bland)
Nəticə: skeptik ol — "congratulations: you're thinking like a tester."

### 6. Persnickety Tests — Həddindən Artıq Test
Optimizmə nisbətən NADİR amma var. **"Halt and catch fire" (Weinberg):**
kompyuter partlamasın, OS çökməsin, istifadəçi xəsarət almasın... — praktikada
gözden kaçmır, test ETMƏ.

**Qaydalar:**
- "Specify precisely what should happen and no more" (GOOS) — scenario-ya aid
  OLMAYAN dəyərləri assert etmə; başqa testlərin davranışını yenidən assert etmə
- İMKANSIZI test etmə (true və false eyni anda) — skeptisizmin dəliliyi.
  Poincaré: "To doubt everything or to believe everything are two equally
  convenient solutions; both dispense with the necessity of reflection."
- ÇOX geniş muqayisə: bütün strukturu/faylı yoxlamaq, halbuki 1 sahə vacibdir
  → brittle + məqsədi gizlədir. "Check only the properties you care about...
  relevant substrings." (Donovan & Kernighan)
- Exact error dəyəri yox, "error VAR/YOX" (errors fəsli)
- Golden fayl tam-miqyaslı yoxlaması — lazım olan hissəvi
- Property-based: dəqiq dəyər YOX, invariant

**Fingerprint nümunəsi (real review):**
```go
// İLK (brittle): MD5-in ÖZÜNÜ yoxlayır
want := md5.Sum(data)
got := fingerprint.Hash(data)
cmp.Equal(want, got)     // → "MD5 insecure, SHA256 edək" → TEST QIRILDI, amma Hash İŞLƏYİR!
```
"What are we really testing here?" — MD5 dəyəri deyil: EYNİ data → eyni hash;
FƏRQLİ data → fərqli hash:
```go
orig := fingerprint.Hash(data)
same := fingerprint.Hash(data)                    // determinizm
different := fingerprint.Hash([]byte("Hello, Newman"))  // fərqlilik
if !cmp.Equal(orig, same) { t.Error("same data produced different hash") }
if cmp.Equal(orig, different) { t.Error("different data produced same hash") }
```
Alqoritm dəyişəndə QIRILMIR (SHA256-də keçir!), unstable/random hash-u TUTUR,
"always return fixed value" tənbəlliyini TUTUR. Daha az brittle + daha az
feeble + az əlavə kod — yalnız əlavə DÜŞÜNMƏ.

### 7. Over-Precise — Float Epsilon
Float64 = lossy compression; dəqiq == demək olar mümkünsüz (64-bitə "neat"
yerləşən rəqəm şansı). DOOM/RP2040 sitatı: demo desync 1/65536 xətada belə —
bit-exact lazım olduqda başqa hadisədir (regression test kimi demos!).

Həll — epsilon müqayisəsi:
```go
if !cmp.Equal(want, got, cmpopts.EquateApprox(0, 0.00001)) {
    t.Errorf("not close enough: want %.5f, got %.5f", want, got)
}
```
**Epsilon dəqiqləşdirmə metodu:** test KEÇƏRKƏN epsilon-u kiçilt → FAIL olanda
bir order-of-magnitude böyüt. (Keçən testin sərhədini tap.)

### 8. Too Many Tests — Bağça Budaması
"Lines spent" (Dijkstra) — hər sətir texniki borc. Sistem inkişaf edir → köhnə
davranışlar ölür → testləri implementasiya ilə BİRLİKDƏ sil.

**Scaffolding:** inşaat boyu dəstəkləyir, bina qalxanda SİLİNƏR. Məs.: 10
davranış × 10 test → komponent bitəndə 1 tam-komponent test əvəz edə bilər.

**Həddi:** 1 test bütün sistem üçün? Bug-u tapar, amma HARDA tapmaz —
"murder suspect is somewhere in North America" (location = detection qədər
vacib). Modul sistemdirsə, test strukturu sistem strukturu əks etir; komponent
testləri ilə birlikdə çıxarıla bilən olmalıdır — yoxsa "big ball of mud".
Bəzi dublikasiya TƏBİİDİR — kritik davranışı müxtəlif bucaqlardan test etmək
faydalıdır; "eyni şeyi eyni yolla 2 dəfə" YOX.

**Delta coverage (Herb Derby/Kent Beck):** testlərin HƏR BİRİNİN unikal
əlavə etdiyi coverage. Delta-sız testlər silinsin (kommunikasiya məqsədi
yoxdursa). Alət: deltacoverage (experimental).

**Fowler-in balansı:** "If you can't confidently change the code, you don't
have enough tests. But if you change the code and you have to fix a bunch of
tests, then you have too many tests."

### 9. Test Frameworklər — import "testing"
testify (require.NoError, assert.Equal), ginkgo/gomega (Describe/Context/It,
Ω Should) — "weird flex, but okay". Pisləri YOXDUR — amma LAZIMSIZDIRLAR.

Go standart kitabxanası HƏR ŞEYİ verir. 3-cü tərəf əlavə etmək: (1) asılılıq;
(2) hər yeni developer üçün öyrənmə əyrisi; (3) testləri anlamağı ÇƏTİNLƏŞDİRİR
— "The last thing we want to do is discourage people from writing tests."

**assert.Equal problemi:** "premature abstraction" — "1 does not equal 0"
deyir, bunu onsuz da bilirdik. Köməkli test FƏRHİ İZAH EDİR + potensial bug-u
təklif edir: `t.Error("want error for empty slice; missing len() check?")`.
Assertion API bunu TƏŞVİQ ETMİR — sonrakı assert-e keçməyi təşviq edir.
Testlər REAL Go koduna oxşamalıdır; real istifadəçilər assert YAZMIR. "Stick
to the standard library unless there's a really compelling reason."

### 10. Flaky Tests
"Sometimes fail, sometimes pass, regardless of correctness."

**Səbəblər + həllər:**
- Fixed sleeps → wait-for-success (ən qısa lazımı intervalı gözlə)
- Vaxt özü test olunursa → ƏN QISA interval (1ms 1s əvəzinə)
- Günün saatı → time-into-data + fake Now() (ch10)
- Map sırası: `for range map` müəyyənsizdir! cmp.Equal isə map-ları SIRADAN
  ASILI OLMAYAN müqayisə edir. Slice-lar isə SIRALIDIR (fail: 1,2,3 vs 3,2,1)
  → sıra vacib deyilsə: cmpopts.SortSlices(func(a, b int) bool { return a < b })
  ilə sort option-u.

**Brittle ≠ flaky:** brittle — ƏLAQƏSİZ dəyişiklikdə fail; flaky — NƏ VAXT
İSTƏSƏ. Brittle → decouple/scope-azalt. Flaky → kök səbəb tap (dəyər
artırdıqsa), yoxsa SİL: "Thou shalt not suffer a flaky test to live. Bad
tests are worse than no tests." ("Oh yeah, that test just fails sometimes" =
test ÖLÜB, hamı ignore edir — dəyər eroziyası əbədi.)

### 11. Shared Worlds
Yavaş/baha setup → testlər arasında DÜNYA paylaşılır → test dünyanın vəziyyətini
BİLMİR → flaky. **Ən yaxşı: hər testə ÖZ private dünyası** (başlanğıc vəziyyət
dəqiq, müdaxilə yoxdur).

Məcburən paylaşma → "take nothing but pictures, leave nothing but footprints":
başlanğıc vəziyyətə etibar ETMƏ; heç nə dəyişdir ki, başqalarına təsir etsin;
dünyanı tapdığın kimi QOYUB ÇIX.

**Pilləkən (DB nümunəsi):** (1) hər test ÖZ bazası (test adıyla!); (2) mümkün
deyilsə — öz CƏDVƏLİ (lock konfliktləri yox); (3) onda da yox — TRANSACTION
daxilində hər şey + sonda ROLLBACK (izolyasiya: hər test öz snapshot-ını görür).

**Qırmızı xətlər:** yaratmadığın heç nəyi SİLMƏ; "drop & recreate database"
avtomasiyası YAZMA — bir gün prod-da başlayacaq. **"It should always be safe
for the tests to run against the production DB, because one day they will."**

### 12. Failing Tests — Zero Defects
Həmişə fail eden test = insanlar inamı itirir ("oh yeah, that test always
fails") → testlər value-siz olur. **Bir test belə fail-ə icazə YOX** —
"as we can never have any bugs" (Maguire: "Don't fix bugs later; fix them now").

Fail başlayanda: hər kəsin TOP PRIORITY-si; bug-fixdən başqa HEÇ NƏ deploy
edilmir. Bir test buraxsan — hamısı değersizləşir. **Broken windows** (Pragmatic
Programmer): "All the rest of this code is crap, I'll just follow suit."

**Infinite defects methodology** (Spolsky): Word 1.0 death march — bug-fix
cədvəldə deyildi → cədvəl = "features waiting to be turned into bugs". Zero
defects radikal deyil — ALTERNATİV NƏDİR?

**Bug bankruptcy/amnesty** (Weinberg): böyük backlog + biznes yaşayır →
köhnə bug-ları HAMISINI bağla / fail testləri SİL → vacib olanlar tezliklə
YENİDƏN açılar (insanlar maraqlıdırsa). Backlog dərhal azalır.

### 13. Slow Tests — Beck Time
Hədd: testlərdən ~5 dəqiqə uzaqda olma. Kent Beck: **10 dəqiqə = "9.8 m/s²"**
psixoloji limit; uzun suite-lər ya kəsilir ya tune olunur → 10 dəqiqəyə qayıdır.
10-dan yuxarı hamı görür + səy qoyur; ALTI idarə olunmur amma şikayət olur —
Beck time ADLANDIRIR. 10 dəqiqə "okay" demək DEYİL.

**7 Sürət Resepti:**
1. **Paralel testlər** — paralellik qeyri-mümkünlüyü = design smell; hər test
   öz dünyası + qlobal state yox → ~10X sürət
2. **I/O-nu ləğv et** — "off the chip" = yavaş; fstest.MapFS (in-memory FS),
   memory-backed Reader/Writer (bytes.Buffer, strings.NewReader)
3. **Remote network çağırışı YOX** — local fake; local networking kernel-dədir,
   wire-dan SÜRƏTLİDİR
4. **Fixture paylaş** — baha setup bir testdə + SUBTEST-lər; amma diqqət:
   fixture-paylaşma flaky yaradarsa — "a flaky test is worse than a slow test"
5. **Fixed sleep YOX** — wait-for-success; sleep = maksimum itilmiş vaxt
6. **Hardware at** — CPU-bound suite → 256-core cloud maşın; "CPU time costs a
   lot less than programmer time" (ucuz proqramçı işə götürmək baha-coverdir)
7. **Slow suite-i gecəlik** — son option; ayrı "slow suite", gündəlik schedule;
   "better than not running tests at all"

### 14. Final — Kent Beck-in Sözü
"If you know what to type, then type. If you don't know what to type, then
take a shower, and stay in the shower until you know what to type. Many teams
would be happier, more productive, and smell a whole lot better if they took
this advice." — kitabın "qoxu" temasına uyğun son.

## Əsas terminlələr
- Test Suite Smells — bütün suite səviyyəsində pozuntu
- Guerilla Testing — icazəsiz, toxunduğun yerdə test yazmaq
- Test What You Touch — Feathers adalar strategiyası
- Legacy Code — testsiz + testability-nəzərəsiz kod (Tarlinder)
- Edit and Pray / Cover and Modify — Feathers dikotomiyası
- Rescue Proseduru — istənilən funksiya üçün test + yeni impl + köhnəni sil
- Zero Defects / Broken Windows / Infinite Defects
- Bug Bankruptcy / Amnesty — backlog-un tam ləğvi (Weinberg)
- Preconditions / Postconditions — test müqaviləsinin iki qolu
- Mirage Test — keçən amma heç nə test etməyən (Create {})
- WHERE-omission testi — Bob sağqalmalı (implicit postcondition)
- "Halt and Catch Fire" — test edilməli olmayan aşkar şeylər (Weinberg)
- Seçici Assert — yalnız scenario-ya aid dəyərlər (GOOS)
- Fingerprint Property — eyni input→eyni output; fərqli→fərqli
- cmpopts.EquateApprox — float epsilon müqayisəsi
- Epsilon tuning — kiçilt → fail → 10X böyüt
- Lines Spent — Dijkstra: hər sətir xərc/borc
- Scaffolding Tests — müvəqqəti, sonra silinən testlər
- Delta Coverage — testin unikal əlavə etdiyi örtük (Derby/Beck)
- Beck Time — 10 dəqiqə psixoloji limit
- Beck's Balance — dəyişə bilmirsən→azdır; çox test qırılırsa→çoxdur
- import "testing" — 3-cü tərəf framework-lərdən imtina
- Premature Abstraction (assert) — fərhi kontekstsiz bildirir
- Flaky vs Brittle — nə vaxt istəsə fail / əlaqəsiz dəyişiklikdə fail
- cmpopts.SortSlices — sıradan asılı olmayan slice müqayisəsi
- Shared World / Private World — test environment paylaşımı
- Transaction + Rollback — DB test izolyasiyası
- Prod-safe Tests — "bir gün prod-a qarşı işləyəcək"

## Praktik nəticə
(1) Test yoxdursa — İCAZƏ GÖZLƏMƏ: toxunduğun hər şeyə test yaz; nəticə =
adalar → qitələr. (2) Untestable legacy: istədiyin funksiyanın testi + yeni
impl + köhnənin silinməsi — Catch-22-dan çıxış. (3) Review-da STİLƏ deyil,
TESTLƏRƏ bax: 10-addımlıq checklist (niyə? test var? fresh checkout? bug
tutur? fail mesajı? fuzz/boş input? coverage? over-engineering? kənar işlər?
aydınlıq?). (4) Optimistic testləri ovla: precondition yoxla (Alice ƏVVƏL
yoxdur!); implicit postcondition yoxla (Bob sağ qalır!); stub-unutma
klassikası. (5) Həddi aşma: halt-and-catch-fire yox; seçici assert; property
(MD5→SHA256 refactor-proof); exact-string/qızıl-fayl tam yoxlama. (6) Float ==
YOX — EquateApprox; epsilon-u sərhəddə gətirib 10X böyüt. (7) Testsiz sıra
YOX, amma "eyni şeyi eyni yolla 2 dəfə" də yox: delta coverage ilə buday,
scaffolding-i sök, ölə davranışın testini+kodunu sil. Fowler balansını yoxla.
(8) Framework YOX: import "testing" + go-cmp; assert.Equal "1≠0" deyir, test
FƏRHİ izah etməlidir. (9) Flaky: kökü tap / SİL — "bad tests are worse than
no tests"; map-ları cmp.Equal ilə, sırasız slice-ları SortSlices ilə. (10)
Shared world pilləkəni: öz DB > özn cədvəl > transaction+rollback; prod-a
qarşı da TƏHLÜKƏSİZ olmalıdır. (11) Fail = top priority, zero defects;
backlog böyükdürsə — bug bankruptcy. (12) Beck time 10 dəq; 7 resept: paralel,
I/O-suz, lokal fake, fixture paylaş (flaky etmədən!), sleepsiz, hardware,
gecelik slow suite. (13) Son söz: nə yazacağını bilirsənsə yaz; bilmirsənsə —
duş al, bİLƏNƏ qədər duşda qal.

## Mənbə
Pages: 346-385 (PDF 358-397)
