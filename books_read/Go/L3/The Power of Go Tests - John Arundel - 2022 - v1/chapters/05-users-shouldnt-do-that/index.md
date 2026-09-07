# Chapter 5 — Users shouldn't do that (İstifadəçilər bunu ETMƏMƏLİDİR)

## Bu fəsil nədən bəhs edir?

Test inputlarının dizaynı (equivalence class, boundary value, adversarial thinking),
uint vs int (type sistem = testlərin bir hissəsi), istifadəçi testi (Margaret
Hamilton/Apollo hekayəsi, hallway usability test), nil/boş/böyük dəyərlərlə
bespoke bug detektorları, table test + subtest (map + t.Run, ad fəzaları), dummy
adlar ("Fake User", "dummy token") ilə RELEVANT/İRRELEVANT data siqnalı, test
data-nın variable/funksiya/fayl bölgəsi (reference type təhlükəsi, makeX()
pattern), testdata/ qovluğu, io.Reader/Writer vs *os.File, fstest.MapFS,
t.TempDir, golden file idarəsi (-update TEHLÜKƏSİ, .gitattributes) və random
inputlara keçid.

## Əsas fikirlər

### 1. Test = Nümunə götürmə (Sampling)
**Problem:** HƏR input-u test etmək mümkünsüz → HANSILARI seçməli?

**Equivalence Class (bərabərlik sinifi):** eyni sinifdəki inputlar test baxımından
EKVİVALENT (sqrt üçün bütün müsbətlər) — bir neçənə baxmaq kifayət. Fərqli
siniflər FƏRQLİ davranış deməkdir (mənfilər — kök TƏYİN EDİLMƏMİŞ).

**Boundary Value (sərhəd dəyəri):** siniflərin KƏNARları — sqrt üçün 1, 0, -1.
Off-by-one xətaları BURADA aşkara çıxır.

**Adversarial thinking:** "A QA engineer walks into a bar. Orders a beer. Orders
0 beers. Orders 99999999999 beers. Orders a lizard. Orders -1 beers."

### 2. Mənfi Dəyər Məsələsi — uint vs int
```go
// Validasiya ÇATIŞMIR:
func DoRequests(reqs int) error {
    if reqs == 0 {                        // -1 KEÇİR!
        return errors.New("number of requests must be nonzero")
    }
}

// Həll 1 — uint (type sistemində):
func DoRequests(reqs uint) error
// DoRequests(-1) → COMPILE XƏTASI: cannot use -1 as uint
```
**Michael Feathers:** "The type system is the set of tests that the language
designer could think up without ever seeing your program."

**Amma uint-in ÇATIŞMazlığı:** -1 case-i artıq REGRESSION TEST kupon YAZILA
BİLMİR (compile olunmur!) → daha yaxşısı if-i gücləndir:
```go
if reqs < 1 {   // 0 VƏ mənfilər — hamısı
    return errors.New("number of requests must be > 0")
}
```

### 3. İstifadəçi Testi — "Users shouldn't do that" YALAN
**Margaret Hamilton (Apollo):** 6 yaşlı qızı P01-i uçuş zamanı seçdi → crash.
Management: "Astronavtlar yaxşı təlim görüb, səhv ETMƏZ." → Apollo 8-də ELƏ
OLDU.

**Praktikalar:**
- Öz proqramını İSTİFADƏ ET: arqumentsiz işə sal, düymələri bas → bug tap →
  test əlavə et → düzəlt → davam et
- İstifadəçilər bug report YAZMIR — sadəcə başqa proqrama keçirlər
- "Usability issues" AYRICA DEYİL — bunlar BUG-dır, testlə düzəlt
- **Hallway usability test** (Spolsky): keçən ilk 5 nəfərə istifadə etdir →
  problemlərin 95%-i aşkar; onların ekran+səsini yaz; SUS — kömək ETMƏ
- Prod-a çıxdıqdan sonra: support biletlərinin HAMISINI filtrləNMƏMİŞ gör

### 4. Özəl Bug Detektorları — Kodu Oxu
Sistem İMPLEMENTASİYASINI bildiyindən spesifik hücumlar:
- slice parametr → NIL slice və YA BOŞ ötür
- pointer → nil ötür
- fayl yazır → disk dolu/read-only halını MƏNTİQİ təhlil et
- Hər "nə səhv ola bilər?" sualına TEST yaz — bug-indi OLMASA BELƏ
  (gələcəkdə qoruma)

### 5. Table Test + Subtest — map + t.Run
**Kitabdan kod nümunəsi:**
```go
func TestListItems(t *testing.T) {
    t.Parallel()
    cases := map[string]struct {           // MAP — ad = açar!
        input []string
        want  string
    }{
        "no items":  {input: []string{}, want: ""},
        "one item":  {input: []string{"a battery"},
                      want: "You can see a battery here."},
        // ...
    }
    for name, tc := range cases {
        t.Run(name, func(t *testing.T) {      // SUBTEST
            got := game.ListItems(tc.input)
            if tc.want != got {
                t.Error(cmp.Diff(tc.want, got))
            }
        })
    }
}
```
**Subtest çıxışı:**
```
--- FAIL: TestListItems/two_items (0.00s)     ← ad avtomatik prefiks!
    game_test.go:46: ...cmp.Diff...
```
- Boşluqlar → `_` ilə əvəz olunur (`two items` → `two_items`)
- Subtest adı = input təsviri, bug ID, problem adı — İSTƏDİYİN etiket
- **Map-in bonusu:** iterasiya SIRASI RANDOM → sıraya asılı gizli bug-lar
  (order-dependent tests) aşkara çıxır!

### 6. Table Test Cüd olunması
**SMELL — eyni testdə 2 FƏRQLİ növ case:**
```go
tcs := []struct{
    input int
    wantError bool       // ← MÜXTƏLİF MƏNTİQ = komplikasiya
}
// subtest body:
if tc.wantError { ... } else { ... }    // ← BU VARSƏ BÖL!
```
**Həll:** valid və invalid input-u İKİ AYRI test-ə böl — subtest body QISA və
HOMOGEN olur. **Qayda:** eyni şeyi yoxlayan input-lar → table test; fərqli
davranış növləri → ayrı test-lər.

**VS Code "Generate Unit Tests" — istifadə ETMƏ:** funksiya-orient (davranış
YOX) testlərə təzyiq + slice (map YOX) generasiya edir.

### 7. Dummy Adlar — Relevantiyyət Siqnalları
**2 növ test data:** RELEVANT (davranışa təsir edir) vs İRRELEVANT (etmir).

**Kitabdan nümunə:**
```go
u := user.New("Fake User")        // İRRELEVANT — AŞKAR saxta ad
u.Language = "Chinese"            // RELEVANT — REAL dəyər
want := "你好"

got, err := CallAPI(req, "dummy token")   // İRRELEVANT — dummy
```
**Konvensiya:** "Fake User", "dummy token", "Bogus Language" = AŞKAR saxta →
oxucu DİQQƏTİ relevant dəyərlərə yönəlir. Real görünən dəyər = "buna bax".

**Kent Beck:** "If there is a difference in the data, then it should be
meaningful."

### 8. Test Data-nın Yerləşdirilməsi
**3 səviyyə:**
```go
// 1) Package-level variable — məzmunu VACİB DEYİLsə:
var validInput = `TGlzdGVuaW5n...`      // ad İZAHEDİR

// 2) Generasiya kodu:
var longInput = strings.Repeat("a", 10_000)
var complicatedInput = func() string { ... }()    // IIFE — 1 dəfə hesabla

// 3) FAYLLAR — testdata/ qovluğunda (Go tools IGNORE edir):
data, err := os.Open("testdata/hugefile.txt")
if err != nil { t.Fatal(err) }
defer data.Close()
```
**Sıra qaydası:** test FUNKSİYALAR əvvəl, sonra variable/helper — oxucu üçün
vacib şey əvvəl.

### 9. Reference Type Təhlükəsi — makeX() Pattern
**Problem:** map/slice GLOBAL variable → hamısı EYNİ underlying memory →
biri dəyişsə hamı POZULUR (paralel testlərdə KONFUZ nəticələr).

```go
// TƏHLÜKƏLİ:
var ageData = map[string]int{ "sam": 18, "ashley": 72 }

// TƏHLÜKƏSİZ — hər çağırış YENİ map:
func makeAgeData() map[string]int {
    return map[string]int{
        "sam": 18, "ashley": 72, "chandra": 38,
    }
}
// Testdə:
got := ages.Total(makeAgeData())     // MÜSTƏQIL nüsxə
```
**Qayda:** map/slice (və belə sahəli struct) test data-ları həmişə FUNKSİYA
ilə — direkt paylaş.

### 10. Reader/Writer vs Fayl
**Prinsip:** funksiya sadəcə bayt axını istəyirsə → io.Reader/Writer (fayl
YOX). Çoxlu şey onları implement edir (test üçün *bytes.Buffer HƏR İKİSİ!).
Funksiya *os.File-spesifik (Name, Stat) istəyirsə → fayl MÜTLƏQ.

```go
// Reader YARATMA (fayl YOX):
input := strings.NewReader("hello world")
parse.ParseReader(input)

// Writer — məzmun maraqsızsa:
tps.WriteReportTo(io.Discard)          // heç yerə yaz

// Writer — məzmun YOXLANACAQSA:
buf := &bytes.Buffer{}
tps.WriteReportTo(buf)
got := buf.String()
```

### 11. fs.FS + fstest.MapFS
```go
// Real fayl sistemi:
fsys := os.DirFS("/home/john/go")
results := find.GoFiles(fsys)

// IN-MEMORY (test üçün — disk I/O YOX, SÜRƏTLİ):
fsys := fstest.MapFS{
    "file.go":                {},
    "subfolder/subfolder.go": {},
    "subfolder2/file.go":     {},
}
results := find.GoFiles(fsys)   // funksiya FƏRQİ BİLMİR
```

### 12. Golden File İdarəsi
```go
func TestWriteReportFile_ProducesCorrectOutputFile(t *testing.T) {
    t.Parallel()
    output := t.TempDir() + "/" + t.Name()
    tps.WriteReportFile(output)
    want, err := os.ReadFile("testdata/output.golden")
    if err != nil { t.Fatal(err) }
    got, err := os.ReadFile(output)
    if err != nil { t.Fatal(err) }
    if !cmp.Equal(want, got) {
        t.Error(cmp.Diff(want, got))
    }
}
```
**Böyük fayllar:** hamısını yükləmə — hash müqayisəsi və ya io.ReadFull +
bytes.Equal ilə chunk-chunk.

**-update flag TEHLÜKƏSİ (müəllif mövqeyi):** çıxışı AVTOMATİK golden-a
yazmaq = "davranış = indi nə edirsə" tərif etmək → test HƏMİŞƏ PASS olur;
sonra KİMİSİ düzəltsə FAIL! Golden dəyişməyi BÖYÜK İŞ etməli — manual.
**Sıralama:** ən yaxşı = real oracle data; sonra = hand-rolled; ƏN SON =
sistem öz çıxışı (mümkünsə ümumiyyətlə YOX).

### 13. Cross-Platform Line Endings
**Problem:** Git `autocrlf` Windows-da golden faylları dəyişir →
`- "Hello world\r\n"` vs `+ "Hello world\n"` — Linux-də pass, Windows-də FAIL.

**Həll:** repo kökündə `.gitattributes`:
```
* -text
```
= bütün faylları binary say → line-ending çevirmə YOX.

## Əsas terminlələr
- Sampling / Input Space — test = nümunə; bütün input mümkünsüz
- Equivalence Class — eyni test dəyərliyi daşıyan input sinfi
- Boundary Value — sinif kənarları (1, 0, -1); off-by-one ovçusu
- Adversarial Thinking — "nə səhv ola bilər?" məntiqi
- uint tipi — mənfiləri compile-da öldürür; amma regression test-i də
- User Testing — öz proqramını istifadə etmə; usability = bug
- Hallway Usability Test — 5 nəfər → 95% problemlər (Spolsky)
- Bespoke Bug Detector — implementasiya bilənə xüsusi sınaqlar
- t.Run / Subtest — adlı uşaq test; fail-də avtomatik prefiks
- map[string]testCase — adlar + RANDOM sıra → order-bug aşkarı
- Table Test Smell — subtest-də if/else = İKİ test-ə BÖL
- Dummy Names — "Fake User"/"dummy token" = irrelevant siqnalı
- IIFE data — `func() T {...}()` bir dəfə hesablanan mürəkkəb data
- Reference Type Share Təhlükəsi — global map/slice korrupsiyası
- makeX() pattern — hər test üçün MÜSTƏQIL data konstruktoru
- testdata/ — Go-nun İGNOR etdiyi test data qovluğu
- io.Discard — "heç yerə" yazan Writer
- *bytes.Buffer — Reader+Writer hər ikisi (test əsası)
- fs.FS / fstest.MapFS — fayl sistemi abstraksiyası / in-memory versiyası
- Golden File — gözlənilən çıxışın saxlanılan test oraklı faylı
- -update Flag — golden-ı avtomatik yeniləyin — TEHLÜKƏLİ
- .gitattributes `* -text` — CRLF çevirməsinin qarşısı

## Praktik nətidə

(1) Input seçimi: equivalence class-lara bölmə + SƏRHƏD dəyərlərini test et.
(2) QA məntiqi: 0, mənfi, çox böyük, tipə uyğun olmayan (lizard!) — heç biri
"mümkün deyil" sayılmır. (3) Mənfiləri type-sistem ilə yoxla (uint) amma
regression test yazıla bilsin deyə if < 1 daha yaxşıdır. (4) Öz proqramını
İSTİFADƏ et — "users shouldn't do that" ƏSAS YALANDIR (Apollo dərsi). (5)
Hallway test: 5 nəfər, səssiz izlə, müdaxilə ETMƏ. (6) Implementation bilirik →
nil/boş slice/fayl xətalarını KODU oxuyaraq test et. (7) Table test: map + ad
+ t.Run — adlar fail-də görünür + random sıra order-bug tutur. (8) Subtest
body-də if/else = SMELL → test-i böl. (9) İrrelevant data = "dummy/fake" adlandır;
relevant = real dəyər. (10) Böyük data: variable → generasiya → testdata/ faylı.
(11) map/slice data → HƏMİŞƏ makeX() funksiyası — global paylaşma = konfuz
paralel fail. (12) io.Reader istəyənə strings.Reader, Writer istəyənə
io.Discard/buf.Buffer — fayl YARATMA. (13) fs.FS qəbul edənə fstest.MapFS —
disk I/O-suz test. (14) Golden fayl: manual dəyiş; -update = test-ləri öldürür.
(15) .gitattributes `* -text` — CRLF golden pozuntularının dərmanı.

## Mənbə
Pages: 125-154 (PDF 137-166)
