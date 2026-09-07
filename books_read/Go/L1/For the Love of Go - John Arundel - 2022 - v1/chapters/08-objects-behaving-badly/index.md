# Chapter 8 — Objects behaving badly (Obyektlərin Pis Davranışı)

## Bu fəsil nədən bəhs edir?

Data-ya DAVRANIŞ vermək: obyekt anlayışı (Book, Customer, Order, Shape, Card —
real dünya və ya abstrakt), properties → fields; PriceCents (float dəqiqliyindən
qaçınmaq üçün SENT-lə int!) + DiscountPercent sahələri; NetPriceCents funksiyası
(test → null impl → real), METHODS — receiver anlayışı
(`func (b Book) NetPriceCents() int`), funksiyadan metoda çevrilmə (b Book
parametri addan ÖNÇƏ), "dynamic struct field" metaforası, method = tip
definisiyasının hissəsi (paketdə olmalıdır, başqa paketin tipinə method YOX),
non-local tip problemi (`cannot define new methods on non-local type int`),
custom type həlli (`type MyInt int` — LOCAL tip!), type conversion
(`mytypes.MyInt(9)`, int≠MyInt — DISTINCT tiplər, string→int QADAĞA),
MyString.Len, Catalog tipi (`type Catalog map[int]Book` — map[int]Book
receiver OLMAZ: "not a defined type"), metodlaşdırılmış GetAllBooks
(`catalog.GetAllBooks()`), receiver adı qısa (c), GetBook məşqi.

## Əsas fikirlər

### 1. Obyektlər — Data + Davranış
Obyekt = proqramın real dünya (və ya abstrakt) vahidi: Book, Customer, Order;
drawing: Shape, Brush; kart oyunu: Card, Hand, Player; hətta Database!
**Properties** (real dünya faktları) → struct FIELDS.
**Davranış** — obyekt passiv saxlamır, BİR ŞEY EDİR (shape özünü çəkir,
scale edir) → KOD lazımdır.

### 2. Book-a Qiymət + Endirim
Float dəqiqlik problemindən (ch3) qaçınma: qiymət = SENT-lərin TAM ədədi:
```go
type Book struct {
    Title           string
    Author          string
    Copies          int
    ID              int
    PriceCents      int     // 4000 = $40.00
    DiscountPercent int     // 25 = 25% off
}
```
NetPrice (endirimdən sonra): klassik TDD.
```go
func TestNetPriceCents(t *testing.T) {
    t.Parallel()
    b := bookstore.Book{
        Title:           "For the Love of Go",
        PriceCents:      4000,
        DiscountPercent: 25,
    }
    want := 3000
    got := bookstore.NetPriceCents(b)
    if want != got {
        t.Errorf("want %d, got %d", want, got)
    }
}
```
Null impl: `return 0` → FAIL ("want 3000, got 0") → real:
```go
func NetPriceCents(b Book) int {
    saving := b.PriceCents * b.DiscountPercent / 100
    return b.PriceCents - saving
}
```

### 3. Methods — Receiver-lə Qısa Yol
"Book parametrli funksiya" pattern-i O QƏDƏR faydalıdı ki, Go QISA YOL verir:
```go
got := b.NetPriceCents()     // b-ni ÖTÜRMƏ — metodu b ÜZƏRİNDƏN çağır
```
**Metod definisiyası:**
```go
func (b Book) NetPriceCents() int {
    saving := b.PriceCents * b.DiscountPercent / 100
    return b.PriceCents - saving
}
```
Dəyişikliklər: (1) `b Book` funksiya adından ƏVVLƏ (adlar mütləqdən əvvəl
mövhərəzələrə keçir — bu, RECEIVER); (2) adından sonrakı mötərizə ARTIQ BOŞ
(parametrlər yoxdur). Body EYNİ qalır — receiver daxildə NORMAL parametr
kimi görünür.

**Metafora:** metod = "dynamic struct field" — hər sorğuda dəyəri TƏLAİBİLİK
hesablayır.

### 4. Metodların Paket Qaydası
NetPriceCents Book-un tipi ilə Sıx bağlı → EYNİ paketdə definisiya olunmalıdır.
**Nəticə:** başqa paketin tipinə (standart kitabxana daxil) metod YAZMAQ
OLMAZ.

### 5. Non-Local Tip Problemi + Custom Type Həlli
```go
func (i int) Twice() int {   // mytypes pakətində
    return i * 2
}
// cannot define new methods on non-local type int
```
int = non-local (bizim paketdə definisiya olunmayıb). Həll — YENİ LOKAL TİP:
```go
type MyInt int                 // underlying type: int — amma BİZİM tipimiz!

func (i MyInt) Twice() MyInt { // indi METHOD MÜMKÜNDÜR
    return i * 2
}
```
Test — TYPE CONVERSION tələbi:
```go
func TestTwice(t *testing.T) {
    t.Parallel()
    input := mytypes.MyInt(9)      // ← int("Hello") kimi DEYİL — mümkündür
    want := mytypes.MyInt(18)
    got := input.Twice()
    if want != got {
        t.Errorf("twice %d: want %d, got %d", input, want, got)
    }
}
```
**Niyə `input := 9` YOX?** int və MyInt = DISTINCT tiplər: int → MyInt
ötürülməz, Twice int-də ÇAĞIRILMAZ, qaytarılan MyInt int-lə müqayisə OLMAZ.
**Conversion sintaksisi:** `TipName(dəyər)` — `mytypes.MyInt(9)`.
**Hər conversion mümkün deyil:** `int("Hello")` → "cannot convert" (compile
xətası); int↔MyInt isə TƏBİİ.

**MyString nümunəsi:**
```go
// MyString is a custom version of the `string` type.
type MyString string

// Len returns the length of the string.
func (s MyString) Len() int {
    return len(s)
}
```
Test: `input := mytypes.MyString("Hello, Gophers!")` → want 15.

### 6. Catalog Tipi — map-ə Metod
İstək: GetAllBooks KATALOQUN metodu olsun (kataloqu soruşan/yeniləyən HƏR
funksiya metoddur). İlk cəhd:
```go
func (c map[int]Book) GetAllBooks() []Book {
// invalid receiver type map[int]Book (map[int]Book is not a defined type)
```
map[int]Book — built-in YOX, amma BİZİM definisiya etdiyimiz tip də DEYİL
(kompozit tip). Həll — eyni fənd:
```go
type Catalog map[int]Book
```
İndi Catalog = LOCAL defined type → metod OLAR!

### 7. Metodlaşmış GetAllBooks
Test yenilənməsi — Catalog LITERAL (map literal ilə EYNİ görünüş, yalnız ad!):
```go
catalog := bookstore.Catalog{
    1: {ID: 1, Title: "For the Love of Go"},
    2: {ID: 2, Title: "The Power of Go: Tools"},
}
got := catalog.GetAllBooks()      // ← funksiya çağırışı DEYİL, METOD
```
Implementasiya — yalnız İMZA dəyişir:
```go
func (c Catalog) GetAllBooks() []Book {
    result := []Book{}
    for _, b := range c {
        result = append(result, b)
    }
    return result
}
```
**Receiver adlandırma konvensiyası:** QISA ad (c = catalog) — body-də çox
istifadə olunacaq. Məntiq EYNİ (slice yığıcı loop).

### 8. GetBook Metodlaşdırma Məşqi
Eyni dönüşüm: `catalog.GetBook(2)` — testlər (valid + invalid ID error) dəyiş,
compile xətalarını al, imzanı dəyiş:
```go
func (c Catalog) GetBook(ID int) (Book, error) {
    b, ok := c[ID]
    if !ok {
        return Book{}, fmt.Errorf("ID %d doesn't exist", ID)
    }
    return b, nil
}
```

## Əsas terminlələr
- Obyekt / Entity — real dünya və ya abstrakt vahidin proqram təmsili
- Properties → Fields — faktlar struct sahələrinə
- Davranış — obyektin EDİYİ şey; kod tələb edir
- PriceCents (int) — float dəqiqliyindən qaçış: sent-lərlə tam ədəd
- Method — receiver-lü funksiya
- Receiver — `(b Book)` — funksiya adından ƏVVƏL xüsusi parametr
- "Dynamic Struct Field" — hər çağırışda hesablanan metod metaforası
- Dot Notation Metod Çağırışı — b.NetPriceCents()
- Non-Local Type — pakət xaricində definisiya olunan tip (int, kitabxana)
- "Cannot Define New Methods on Non-Local Type" — başqa paketin tipinə metod QADAĞA
- type MyInt int — underlying tipli YENİ LOCAL tip
- Underlying Type — yeni tipin əsaslandığı tip
- Type Conversion — mytypes.MyInt(9); TipName(dəyər)
- Distinct Types — MyInt ≠ int; bir-birinə ötürülməz
- MyString.Len — string əsaslı custom tip metodu
- type Catalog map[int]Book — defined type; map receiver-i YOX
- Receiver Ad Konvensiyası — qısa ad (c, b, s)
- Metodlaşdırma — funksiya → metod dönüşümü (yalnız imza dəyişir)

## Praktik nəticə
(1) Obyektə hesablama lazımdırsa — METOD yaz: `func (b Book) Name() T`;
parametr addan əvvəl mötərizədə, çağırış b.Name() formasında. (2) Funksiyanı
metoda çevirmək ucuzdur: imza dəyiş, body eynidir; testdə çağırışı dəyiş. (3)
Başqa paketin tipinə metod YAZMAQ OLMAZ — `type MyX ExistingType` ilə ÖZ
local tipini yarat və metod ona. (4) Yeni tip underlying-dən AYRIDIR:
dəyərlər CONVERT olunmalı (`MyInt(9)`); string→int kimi mənasız conversion
compile xətasıdır. (5) Qiymətləri FLOAT-la YOX, sent-lərlə İNT saxla — dəqiqlik
 problemi yoxdur. (6) Composite tipə metod lazımdırsa (`map[int]Book`) —
adlı defined type yarat (`type Catalog ...`) və metodu ona yaz. (7) Receiver
adını qısa tut (c/b/s). (8) "Kataloqla işləyən hər funksiya kataloqun metodu
olmalıdır" — API dizayn prinsipi. (9) Metodlar tipin PAKETİNDƏ olur — paket
sərhədini pozmaq mümkün deyil.

## Mənbə
Pages: 98-108 (PDF 99-109)
