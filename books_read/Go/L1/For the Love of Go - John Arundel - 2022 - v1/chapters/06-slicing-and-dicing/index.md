# Chapter 6 — Slicing & dicing (Dilimləmə və Doğrama)

## Bu fəsil nədən bəhs edir?

Slice-lər — eyni tipli dəyərlər toplusu: []Book. Slice dəyişənləri, slice
literal-ları (struct-larda tip adı təkrarı lazım DEYİL), index (books[0],
0-based offset), len() built-in, element/sahə modifikasiyası (books[0].Title),
append(). İkinci user story: GetAllBooks — catalog parametri, "setting up the
world", slice müqayisəsi == İLƏ MÜMKÜN DEYİL → go-cmp paketi (go get -t,
cmp.Equal, cmp.Diff — Unix diff tərzli, "-" want / "+" got), null impl fail
çıxışının təhlili (nil vs slice), GetAllBooks = return catalog. Üçüncü story:
GetBook — unikal identifikator problemi (title/author unikal deyil, ISBN hər
kitabda YOX), ID int sahəsi, TestGetBook, zero-value fail ("" / 0 sahələr),
for-range ilə axtarış, "Crime doesn't pay" — `return catalog[0]` yalançı
implementasiya testdən KEÇİRDİ → 2 kitablı test kəsifi.

## Əsas fikirlər

### 1. Slice — Eyni Tiplər Toplusu
Struct = MÜXTƏLİF tiplər bir vahiddə; slice = EYNİ tip dəyərlər toplusu.
```go
var books []Book      // "slice of Book" — boş mötərizələr + element tipi
books = []Book{}      // boş slice literal

books = []Book{
    {Title: "Delightfully Uneventful Trip on the Orient Express"},
    {Title: "One Hundred Years of Good Company"},
}
```
**Gözəllik:** slice elementlərində tip adı YAZILMIR ({Title: ...}, Book{...}
DEYİL) — []Book yalnız Book qəbul edə bilər, compiler İNFİRA edir. Hər slice
tİPİ ayrı tipdir ([]Book ≠ []string).

### 2. Index + len + Modifikasiya + append
```go
first := books[0]              // index: 0-based (başdan OFFSET kimi düşün)
fmt.Println(len(books))        // 2 — built-in length
books[0] = Book{Title: "Heart of Kindness"}  // element əvəzi
books[0].Title = "Heart of Kindness"         // elementin SAHƏSİ
b := Book{Title: "The Grapes of Mild Irritation"}
books = append(books, b)       // SONA əlavə — NƏTİCƏNİ TƏYİN ET!
```
append: slice + element(?) → YENİ slice (yeni element sonuncu). Digər tiplərdə
olmayan xüsusiyyət: slice BÖYÜYÜR.

### 3. GetAllBooks — İkinci Core Story
Davranış: catalog ([]Book) → GetAllBooks onu qaytarsın. Test strukturu:
```go
catalog := []bookstore.Book{
    {Title: "For the Love of Go"},
    {Title: "The Power of Go: Tools"},
}
want := []bookstore.Book{...aynı...}
got := bookstore.GetAllBooks(catalog)
// müqayisə...
```
**Adlandırma qaydası:** eyni dəyişən hər yerdə EYNİ adla (catalog).

### 4. Slice Müqayisəsi — go-cmp Paketi
```go
if want != got {
// invalid operation: want != got (slice can only be compared to nil)
```
**== slice-larda TƏYİN OLUNMAMIŞDIR!** (yalnız nil ilə müqayisə olunar).
Köhnə yol: reflect.DeepEqual. Daha yaxşı: **go-cmp**:
```bash
go get -t     # testlər üçün lazımi paketləri YÜKLƏ (added go-cmp v0.5.6)
```
```go
import "github.com/google/go-cmp/cmp"

if !cmp.Equal(want, got) {
    t.Error(cmp.Diff(want, got))
}
```
**cmp.Diff nümunəsi:**
```go
want := []string{"same", "same", "same"}
got := []string{"same", "different", "same"}
```
```
[]string{
    "same",
-   "same",        ← want (gözlənilən) - prefiksi
+   "different",   ← got (alınan) + prefiksi
    "same",
}
```
Element-element, Unix diff kimi. Sadə tiplər üçün ==; struct/slice/mürəkkəb
 üçün cmp.Equal/cmp.Diff — İNDISPENSABLE.

**Null impl fail çıxışı:**
```
- {{Title: "For the Love of Go"}, {Title: "The Power of Go: Tools"}},
+ nil,
```
→ want = 2-kitablı slice, got = nil. Implementasiya:
```go
func GetAllBooks(catalog []Book) []Book {
    return catalog
}
```

### 5. GetBook — Unikal ID Problemi
Üçüncü story: konkret kitabın detalları. Sual: kitabı NƏ İLƏ tapırıq?
- Title? — eyni adlı çox kitab VAR
- Author/price/hər hansı sahə? — heç biri UNİKAL deyil
- ISBN? — var, amma HƏR KİTABDA YOX (bu kitabda yoxdur!)
**Nəticə:** unikal identifikator (ID) lazım. Nə olduğu ƏHƏMİYYƏTSİZ — yalnız
HƏR kitabı FƏRQLİ identifikasiya etsin. Seçim: `ID int` (kataloq nömrəsi).

Book-ya yeni sahə:
```go
type Book struct {
    Title  string
    Author string
    Copies int
    ID     int
}
```

### 6. TestGetBook + Zero-Value Fail
```go
func TestGetBook(t *testing.T) {
    t.Parallel()
    catalog := []bookstore.Book{
        {ID: 1, Title: "For the Love of Go"},
    }
    want := bookstore.Book{ID: 1, Title: "For the Love of Go"}
    got := bookstore.GetBook(catalog, 1)
    if !cmp.Equal(want, got) {
        t.Error(cmp.Diff(want, got))
    }
}
```
Null impl (`return Book{}`) fail-i — sahə-sahə:
```
-  Title: "For the Love of Go",     ← gözlənilən
+  Title: "",                        ← alınan (zero value!)
   Author: "",
   Copies: 0,
-  ID: 1,
+  ID: 0,
```
Boş Book = bütün sahələr zero value ("" və 0).

### 7. GetBook Implementasiyası — range Axtarışı
Ən sadə işləyən yol: HƏR kitaba BAX (prototype mərhələsi — elegant deyil,
scalable deyil, amma İŞLƏYİR):
```go
func GetBook(catalog []Book, ID int) Book {
    for _, b := range catalog {
        if b.ID == ID {
            return b
        }
    }
    return Book{}
}
```
- `for _, b := range catalog` — test case-lərdəki kimi; b = ardıcıl kitab
- Tapanda DƏRHAL return
- Tapmasa — boş Book (bu fəsildə yeganə mənalı default; error halı sonra)

### 8. "Crime Doesn't Pay" — Yalançı Implementasiya Testi
Sual: "Is there an INCORRECT implementation that would still pass this test?"
Bu test üçün VAR:
```go
func GetBook(catalog []Book, ID int) Book {
    return catalog[0]        // HƏMİŞƏ ilk kitab — AÇIQ-ASKER SƏHV!
}
```
Amma KEÇİR — çünki testdə catalog-da CƏMİ 1 kitab VAR (ilk kitab = istənilən
kitab)!

**Həll:** catalog-a İKİNCİ kitabı əlavə et, ID=2-ni SORUŞ:
```go
catalog := []bookstore.Book{
    {ID: 1, Title: "For the Love of Go"},
    {ID: 2, Title: "The Power of Go: Tools"},
}
want := bookstore.Book{ID: 2, Title: "The Power of Go: Tools"}
got := bookstore.GetBook(catalog, 2)
```
İndi `return catalog[0]` FAIL edəcək. **Metod:** mütləq certainty YOXDUR, amma
"səhv ola biləcək yolları düşünüb testləri gücləndir" = inam artdırıcı.
(Kitabda buna "Crime doesn't pay" — "cinayət ödənilmir" deyilir: yalançı kod
testdən qaça bilmir.)

## Əsas terminlələr
- Slice — eyni tipli elementlərin sıralı toplusu ([]Book)
- []T Notasiyası — boş mötərizə + element tipi
- Slice Literal — []Book{{...}, {...}} — elementlərdə tip adı YOX
- Index (books[0]) — 0-based ofset göstəricisi
- len() — built-in element sayı
- Element Modifikasiyası — books[0].Title = "..."
- append(books, b) — sonu element əlavə; nəticə TƏYİN olunmalı
- go-cmp — google/cmp müqayisə paketi
- cmp.Equal — struktur/slice bərabərliyi (bool)
- cmp.Diff — element-element fərq (- want, + got)
- go get -t — test dependencies yüklə
- Unique Identifier (ID) — dəyərin məzmunu əhəmiyyətsiz, unikallıq vacib
- ID int — Book struct-a kataloq nömrəsi sahəsi
- for-range Axtarışı — bütün slice üzrə müqayisə dövrü
- "İncorrect implementation passes?" — test gücünü yoxlayan sual
- catalog[0] Trap — 1-elementli test yalançı kodu KEÇİRİR → 2 element

## Praktik nəticə
(1) Kolleksiya lazımdırsa — slice: `var books []Book`, literal-da tip adını
təkrarlama. (2) Index 0-dan; len() say verir; append nəticəni qaytarır —
TƏYİN ET (`books = append(books, b)`). (3) Slice-ları == ilə MÜQAYİSƏ ETMƏ —
xəta verir; cmp.Equal + cmp.Diff işlət (go get -t ilə yüklə). (4) cmp.Diff
oxunuşu: `-` = want (gözlənti), `+` = got (fakt). (5) "World setup" pattern:
catalog yarat → want kopyası → funksiya → müqayisə. (6) Unikal seçim lazımdırsa
ID sahəsi əlavə et (int başlanğıc üçün kifayət). (7) Axtarış: for-range + if
b.ID == ID + return; default return Book{}. (8) Test yazdıqdan sonra ÖZÜNƏ
sual ver: "HANSI SƏHV implementasiya bu testdən keçə bilər?" → 1-elementli
katalog = `catalog[0]` tələsi → 2-ci element ilə ID soruş. (9) Null impl
fail-i zero value-ləri göstərir ("" / 0) — sahə-sahə Diff oxu. (10) Adlar
konsistli: eyni data = eyni ad (catalog) hər yerdə.

## Mənbə
Pages: 73-84 (PDF 74-85)
