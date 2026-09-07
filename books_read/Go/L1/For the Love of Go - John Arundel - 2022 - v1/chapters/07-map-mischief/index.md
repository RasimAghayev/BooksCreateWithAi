# Chapter 7 — Map mischief (Xəritə Məzələti)

## Bu fəsil nədən bəhs edir?

Map tipi — slice-ın axtarış problemi həlli: `map[int]Book` (key → element,
loop YOXDUR, birbaşa müraciət). Map literal, catalog-un slice→map transformu
(test dəyişməz qalır — decoupling dərsi), `catalog[ID]` lookup, element əlavə/
overwrite (`catalog[3] = ...`, dublikat KEY yoxdur, dublikat ELEMENT ola
bilər), sahə oxunuşu `catalog[1].Title` (İŞLƏYİR) vs sahə YAZIŞI
(`catalog[1].Title = ...` — COMPILE XƏTASI! "cannot assign to struct field in
map" — çıxar → dəyiş → geri yaz), mövcud olmayan açar = zero value (panic
YOX!) + fail nümunəsi, `b, ok := catalog[ID]` comma-ok idiomu, GetBook
error-lu versiyası (fmt.Errorf ilə ID interpolasiyası), TestGetBookBadIDReturnsError,
GetAllBooks-un map variantı (loop + append ilə slice qur), **map sırası
RANDOM-dır** — yarı testlər fail olur! → sort.Slice + function literal
(anonymous func, i/j parametrləri, müqayisə funksiyası), flaky testin
sistematik həlli.

## Əsas fikirlər

### 1. Problem: Loop Axtarışı Yavaşdır
Böyük katalogda ID-ləri bir-bir müqayisə etmək SÜRƏTLİ DEYİL. Lazım olan:
ID → Book BİRBAAŞA mapping. Bu data strukturun adı: **map**.

### 2. Map — Key → Element
Slice: elementlər RƏQƏM (index) ilə identifikasiya. Map: İSTƏDİYİN TİP key
olun (bizimki int — book ID-lər):
```go
catalog := map[int]bookstore.Book{
    1: {ID: 1, Title: "For the Love of Go"},
    2: {ID: 2, Title: "The Power of Go: Tools"},
}
```
- Tip oxunuşu: "map of int to bookstore.Book"
- Literal: `key: element` cütləri + hər elementdən sonra VİRGÜL (slice/map/struct
  literal ümumi qaydası)
- Book literal-ları dəyişməz — yalnız önlərinə `1:` açarı gəlir

### 3. Test Dəyişmədi — Decoupling Dərsi
TestGetBook yalnız catalog TİPİ-ni dəyişdi, qalanı EYNİ: "tests shouldn't be
too tightly coupled to implementation details" — əgər sıx bağlı olsaydı, public
API da bağlı olardı. Faydalı abstraksiya istifadəçinin maraqlı olmadığı
detalları GİZLƏDİR.

### 4. GetBook Map ilə — Loop Artıq Lazım Deyil
```go
func GetBook(catalog map[int]Book, ID int) Book {
    return catalog[ID]         // birbaşa lookup!
}
```
Kvadrat mötərizə slice index-i ilə EYNİ sintaksis.

### 5. Map Element Əməliyyatları
```go
// ƏLAVƏ / OVERWRITE:
catalog[3] = Book{ID: 3, Title: "Spark Joy"}
// key varsa → KÖHNƏ dəyir əvəz olunur; yoxdursa → YENİ element
// dublikat KEY MÜMKÜN DEYİL; dublikat ELEMENT (fərqli key) MÜMKÜNDÜR

// SAHƏ OXU (işləyir):
fmt.Println(catalog[1].Title)    // For the Love of Go
// 2 hissə: map index (catalog[1]) + field selector (.Title)

// SAHƏ YAZ (İŞLƏMƏZ!):
catalog[1].Title = "..."
// cannot assign to struct field catalog[1].Title in map
```
**Sahə dəyişmə resepti** — 3 addım:
```go
b := catalog[1]          // 1. çıxar
b.Title = "For the Love of Go"   // 2. dəyişdir
catalog[1] = b           // 3. geri yaz
```

### 6. Mövcud Olmayan Açar — Panic Yox, Zero Value!
`GetBook(catalog, 3)` (ID 3 map-də YOX) → panic GÖZLƏNİR (slice out of range
kimi)... amma XƏTA YOXDUR: boş Book qaytarılır! Fail output:
```
-  Title: "The Power of Go: Tools",     ← want
+  Title: "",                            ← got: ZERO VALUE Book
   Author: "", Copies: 0, ID: 0,
```
**Map qaydası:** olmayan key → element tipinin ZERO VALUE. (Null impl-in
`Book{}` qaytarması ilə EYNİ nəticə — aydın izah.)

### 7. Comma-Ok Idiomu
Mövcudluğu necə BİLƏK? Book{} müqayisəsi — "hacky". Go-nun NEAT uzantısı:
```go
b, ok := catalog[ID]
```
- `ok` (konvensional ad) — bool: true = key TAPILDI, false = yoxdur
- Bu, "value-in-map" yoxlamasının STANDART yolu

### 8. GetBook — Error-lu Final Versiya
Design qərarı: ok qaytarmaq əvəzinə STANDART error mexanizmi (calculator
testiləri kimi):
```go
func TestGetBookBadIDReturnsError(t *testing.T) {
    t.Parallel()
    catalog := map[int]bookstore.Book{}
    _, err := bookstore.GetBook(catalog, 999)
    if err == nil {
        t.Fatal("want error for non-existent ID, got nil")
    }
}

func GetBook(catalog map[int]Book, ID int) (Book, error) {
    b, ok := catalog[ID]
    if !ok {
        return Book{}, fmt.Errorf("ID %d doesn't exist", ID)
    }
    return b, nil
}
```
**fmt.Errorf dərsi:** invalid-input mesajına SƏBƏB DATA-nı sal ("ID 3 doesn't
exist") — errors.New yerinə fmt.Errorf + %d interpolasiya. Valid testdə:
`got, err :=` + `t.Fatal(err)` (error = fatal, data etibarsız).

### 9. GetAllBooks — Map-dən Slice Qurma
```go
func TestGetAllBooks(t *testing.T) {
    t.Parallel()
    catalog := map[int]bookstore.Book{   // MAP versiyası
        1: {ID: 1, Title: "For the Love of Go"},
        2: {ID: 2, Title: "The Power of Go: Tools"},
    }
    want := []bookstore.Book{...}        // want SLICE qalır
    got := bookstore.GetAllBooks(catalog)
    ...
}

func GetAllBooks(catalog map[int]Book) []Book {
    result := []Book{}
    for _, b := range catalog {
        result = append(result, b)
    }
    return result
}
```
range map üzrə: İKİ dəyər (key + element); `_` = key (lazım deyil).
Əvvəlcə compiler xətası: `cannot use catalog (type map[int]Book) as type
[]Book in return argument` — köhnə `return catalog` artıq TİPƏ uyĞUN DEYİL.

### 10. FLAKY TEST — Map Sırası Random-dır!
Test "keçdi" → sonra FAIL → sonra KEÇDİ... "Wait, what?"

Simptom: yarı vaxt kitablar TƏRS sıradadır (ID 2 birinci). Kiçik təcrübə:
```go
for ID := range catalog { fmt.Println(ID) }
// Run 1: 1, 2
// Run 2: 2, 1     (≈ hər ikinci dəfə!)
```
**İzah — DİZAYN BƏYANATIDIR:** Slice = SIRALI (spesifik ardıcıllıq); map =
"just a bunch of key-value pairs" — HEÇ BİR TƏYİN OLUNMUŞ SIRA YOXDUR. range
hər run-da FƏRLİ sıra verə bilər. 2 kitab = 2 mümkün sıra = ~50% fail.

### 11. Həll — Testi Yox, SORĞUNU Düzəlt
Sual: "What are we really testing?" → Kitabların SIRASI? MARAQSIZDIR!
Yalnız HAMISININ olması. Yəni GetAllBooks DOĞRU işləyir; PROBLEM TESTDƏ —
sıralı slice ilə müqayisə edir.

**sort.Slice:**
```go
sort.Slice(got, func(i, j int) bool {
    return got[i].ID < got[j].ID
})
```
- Standart kitabxana; arqumentlər: slice + müqayisə funksiyası
- Müqayisə funksiyası: (i, j) → bool — "i, j-dən ƏVVƏL gəlməlidirmi?"
- Bizim hal: ID artan sıra → `got[i].ID < got[j].ID`
- **İn-place** sort — got dəyişir
- want onsuz da ID sıralıdıq → müqayisə HƏMİŞƏ bərabər → test RELIABLE

Nəticə: PASS, PASS, PASS — "No more mysterious intermittent failures."

### 12. Function Literal (Anonymous Function)
```go
func(i, j int) bool {
    return got[i].ID < got[j].ID
}
```
- `func` + parametrlər + nəticə + body — ad YOXDUR (adsız)
- Niyə? Bir dəfə istifadə olunur, ad VERMƏK israf — anonymous
- Literal ailəsi: string, struct, slice, map, FUNKSİYA — hamısı literaldir
- (Closure effekti: `got` funksiya literalından görünür — daha sonra)

## Əsas terminlələr
- map[K]V — key-dən elementə birbaşa xəritə (loop YOX)
- Map Literal — {key: element, ...} — hər cütdən sonra VİRGÜL
- catalog[ID] — slice index-i ilə eyni sintaksisli lookup
- Element Overwrite — mövcud key-ə yazma köhnəni əvəz edir
- Dublikat Key Qadağası — eyni map-də eyni key 2 dəfə OLMAZ
- "Cannot Assign to Struct Field in Map" — sahə yazışı qadağan → 3 addım
- Çıxar-Dəyiş-Geri Yaz — b := m[k]; b.F = ...; m[k] = b
- Zero Value Lookup — olmayan key → element tipinin sıfır dəyəri (panic YOX)
- Comma-Ok Idiomu — b, ok := m[k]; ok = mövcudluq
- fmt.Errorf — interpolasiyalı error ("ID %d doesn't exist")
- range map — (key, element) cütü; `_` ilə key ötürülür
- Random Map Order — map sırası TƏYİN OLUNMAMIŞDIR — DİZAYNDIR
- Flaky Test — yarı keçən test; kökü: sıralı müqayisə + sırasız mənbə
- sort.Slice(slice, cmpFunc) — in-place sıralama
- Comparison Function — (i, j) → "i əvvəldirmi?" bool
- Function Literal / Anonymous Function — adsız, bir istifadəlik func

## Praktik nəticə
(1) ID → obyekt birbaşa müraciət lazımdırsa — map (loop-lu axtarışı əvəz et).
(2) Lookup sintaksisi slice ilə eynidir; yazışda mövcud key = overwrite.
(3) Map elementinin SAHƏSİNƏ birbaşa yazma MÜMKÜN DEYİL — çıxar/dəyiş/geri yaz.
(4) Olmayan key panic DEYİL — zero value qayıdır; ayırd etmək üçün comma-ok.
(5) "Tapılmadı" halı üçün error qaytar; mesajda DATA-nı əks etdir (fmt.Errorf).
(6) Map-dən slice: `result := []T{}; for _, v := range m { result = append(...) }`.
(7) Map üzrə range SIRASI RANDOM-dur — 50% flaky testlərin mənbəyi; order
vacib deyilsə testi DÜZƏLT (sort.Slice + ID müqayisəsi), bir dəyəri də qoru.
(8) sort.Slice müqayisə funksiyası: `got[i].ID < got[j].ID` = artan ID;
in-place dəyişir. (9) Function literal: bir istifadəlik anonymous funksiya —
sort-un 2-ci arqumenti kimi təbii. (10) "What are we really testing?" —
sıra yox, TAMLIQDUR; testi real gözləntiyə uyğunlaşdır.

## Mənbə
Pages: 85-97 (PDF 86-98)
