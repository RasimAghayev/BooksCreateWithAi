# Chapter 10 — Very valid values (Çox Etibarlı Dəyərlər)

## Bu fəsil nədən bəhs edir?

Validasiya: sahələrə yanlış dəyər təyin etmə riski (mənfi PriceCents,
DiscountPercent > 100). SetPriceCents metodu (valid/invalid 2 test), value
receiver tələsi (b.Book → "want updated price 3000, got 4000" — dəyişiklik
QALMADI; staticcheck SA4005 "ineffective assignment"), pointer receiver
həlli + **avtomatik dereferencing** (b.PriceCents — Go *Book-dan struct-a
özü keçir; `(*b).PriceCents` yazmaq lazım deyil); **always valid field** —
unexported `category string` + accessor metodlar (Category/SetCategory, getter
value receiver, setter pointer receiver); setter+getter bir test cütü;
**cmp.Equal + unexported sahə PANİK** → cmpopts.IgnoreUnexported; **always
valid struct** — creditcard paketi: unexported `card` + constructor `New`
(number == "" yoxlama), getter Number (ad istifadə olunmur — dəyər kimi işlənir);
**map-as-set** (`validCategory = map[string]bool{...}` — missing key = false
zero value!; ad təklikdə "validCategory" — oxunuş üçün); **constants**
(`const CategoryAutobiography = "..."`, qruplaşdırma, compiler typo yoxlaması,
http.StatusOK nümunəsi); dəyər əhəmiyyətsiz olanda int konstantlar + **type
Category int** + **iota** (0,1,2 avtomatik, ortaya qırmızıqotum, "not one
iota" = Yunan ən kiçik hərfi), table testlə 3 kateqoriya, SetCategory(Category)
validasiyası map ilə.

## Əsas fikirlər

### 1. Problem — Validasiyasız Sahələr
`b.PriceCents = -1` — compile OLUR, amma MƏNASIZDIR. Compiler TİP yoxlayır,
DƏYƏR yoxlamır. Həll: sahəni DİREKT yox, METOD ilə dəyiş — metod YOXLAYA bilər.

### 2. SetPriceCents — Value Receiver Tələsi
```go
func TestSetPriceCents(t *testing.T) {
    t.Parallel()
    b := bookstore.Book{Title: "For the Love of Go", PriceCents: 4000}
    want := 3000
    err := b.SetPriceCents(want)
    if err != nil { t.Fatal(err) }
    got := b.PriceCents
    if want != got { t.Errorf("want updated price %d, got %d", want, got) }
}
```
İlk cəhd — VALUE receiver:
```go
func (b Book) SetPriceCents(price int) error {
    b.PriceCents = price     // nope!
    return nil
}
// --- FAIL: want updated price 3000, got 4000
```
**Səbəb:** value receiver = kopya üzərində iş; dəyişiklik metoddan çıxanda
ATILIR (ch9 pass-by-value dərsi receiver-də!). Həll:
```go
func (b *Book) SetPriceCents(price int) error {
    if price < 0 {
        return fmt.Errorf("negative price %d", price)
    }
    b.PriceCents = price
    return nil
}
```
**staticcheck** (VS Code avtomatik işlədir) xəbərdarlıq edir:
`ineffective assignment to field Book.PriceCents (SA4005)` — bu problem bir çox
yeni Goçu yıxır; diqqətli testlər TUTUR.

### 3. Avtomatik Dereferencing
`b` = *Book (pointer). `b.PriceCents` necə işləyir? Pointer-in SAHƏSİ yoxdur —
struct-un var! Go ağıllıdır: `b.PriceCents` = "dereference b → struct-ın
PriceCents sahəsi". Əks halda `(*b).PriceCents = price` yazmalı olardıq —
awkward. Struct metodlarının ÇOXU pointer olduğundan bu çox faydalıdır.

### 4. Always Valid Field — unexported + Accessor
SetPriceCents olsa da, `b.PriceCents = -1` DİREKT hələ mümkündür (exported
sahə!). Güclü həll — **sahəni UNEXPORTED et:**
```go
type Book struct {
    ...          // existing fields
    category string   // KİÇİK hərf — paketdən kənarda GÖRÜNMƏZ!
}
```
Xarici kod `b.category = "bogus"` yaza BİLMƏZ — compile xətası. Oxu/yazı =
metodlarla:
```go
func (b *Book) SetCategory(category Category) error {
    if !validCategory[category] {
        return fmt.Errorf("unknown category %v", category)
    }
    b.category = category
    return nil
}

func (b Book) Category() Category {   // dəyişmədiyi üçün VALUE receiver
    return b.category
}
```
**Setter/Getter testi — bir cüt:** SetCategory düzgünsə Category göstərir;
Category səhvdirsə set-olunan görünməz; hər ikisi səhvdirsə kateqoriya yoxdur
— hamısı eyni test cütlüyündə aşkar. Ayrı Category testi LAZIM DEYİL (zərəri
də yoxdur).

**Güclü ifadə:** "Making it impossible to compile incorrect programs is the
best kind of validation!" — invalid yolu CONNECTION KƏSİLMƏSİ runtime yoxlamadan
YAXŞIDIR.

### 5. cmp.Equal + Unexported Sahə = PANİK
category sahəsi əlavə olunca köhnə testlər birdən:
```
panic: cannot handle unexported field at {[]bookstore.Book}[1].category
consider using a custom Comparer; ... Exporter, AllowUnexported, or
cmpopts.IgnoreUnexported
```
cmp.Equal paketdən kənarda oldığından unexported sahələrə BAXA BİLMƏZ (hamı
kimi). Panik mesajının ÖZÜ cavabı deyir — IgnoreUnexported opsiyası:
```go
import "github.com/google/go-cmp/cmp/cmpopts"

if !cmp.Equal(want, got,
    cmpopts.IgnoreUnexported(bookstore.Book{})) {
```

### 6. Always Valid Struct — creditcard Paketi
Fikri STRUCT-a qədər genişləndir: ödəniş kartları. Kart məlumatı user
input-undandır (web form) — string (int DEYİL: istifadəçi hərf daxil edə bilər).
```go
package creditcard

type card struct {          // UNEXPORTED TİP!
    number string
}

func New(number string) (card, error) {
    if number == "" {
        return card{}, errors.New("number must not be empty")
    }
    return card{number}, nil
}

func (c *card) Number() string {
    return c.number
}
```
Test:
```go
func TestNew(t *testing.T) {
    t.Parallel()
    want := "1234567890"
    cc, err := creditcard.New(want)
    if err != nil { t.Fatal(err) }
    got := cc.Number()
    if want != got { t.Errorf("want %q, got %q", want, got) }
}

func TestNewInvalidReturnsError(t *testing.T) {
    t.Parallel()
    _, err := creditcard.New("")
    if err == nil { t.Fatal("want error for invalid card number, got nil") }
}
```
**Zərif məqam:** test `creditcard.card` adını HEÇ VAXT İSTİFADƏ ETMİR! New
"bir tipin dəyərini" qaytarır — adını bilmək lazım DEYİL; Number metodu
exported → çağırıla bilər. Nəticə: paketdən kənardan invalid card YARATMAQ
MÜMKÜNSÜZDÜR (literals ad tələb edir, ad görünməz!). SetNumber YOXDUR — nömrə
dəyişmək məntiqi yox, YENİ kart yaradılar.

### 7. Map-as-Set — validCategory
Kateqoriyalar çoxalırsa if-zənciri YOX:
```go
var validCategory = map[string]bool{
    "Autobiography":       true,
    "Large Print Romance": true,
    "Particle Physics":    true,
}

if validCategory[category] {   // var → true; YOX → bool zero = false!
```
**Zərif məntiq:** missing key → zero value → bool-un zero-su = false. Yoxlama
"ideal" — əlavə kod YOX. **Adlandırma incəliyi:** `validCategories` deyil,
`validCategory` — `if validCategory[...]` İNGİLİSCƏ kimi oxunur; kodun zövqlü
oxunuşu üçün kiçik detal.

### 8. Constants — Compiler Yoxlamalı Dəyərlər
String-lər arbitrary-dır: typo = runtime-a qədər gizli. Konstantlar:
```go
const (
    CategoryAutobiography     = "Autobiography"
    CategoryLargePrintRomance = "Large Print Romance"
    CategoryParticlePhysics   = "Particle Physics"
)

err := b.SetCategory(bookstore.CategoryLargePrintRomance)
// typo:
err := b.SetCategory(bookstore.CategoryLargePrintBromance)
// undefined: bookstore.CategoryLargePrintBromance  ← COMPILE XƏTASI!
```
**Standart kitabxana nümunəsi:** net/http status kodları — http.StatusOK,
http.StatusNotFound; http.StatusBogus = compile xətası; "200" literalı mənasızdır,
StatusOkEq İVOKATİVdir.

### 9. Dəyər Əhəmiyyətsiz olanda — type + iota
Dəyərin ÖZÜ önəmsizdisə (yalnız FƏRQLİLİK vacibdir) — string lazım deyil:
```go
type Category int

const (
    CategoryAutobiography     Category = iota   // 0
    CategoryLargePrintRomance                    // 1
    CategoryParticlePhysics                      // 2
)
```
- İlk konstanta `= iota` → 0; ardıcıllar avtomatik 1, 2, 3...
- Tip yalnız BİRİNCİDƏ bildirilir — qalanları eyni tipi alır
- **iota** — Yunan əlifbasının ƏN KİÇİ hərfi; gündəlik məna "ən azacıq"
  ("not one iota")
- **Üstünlük:** ortaya yeni kateqoriya qırmızıqotum → avtomatik renumber —
  "automagically. Thanks, iota!"

### 10. Yekun SetCategory Sistemi
```go
type Category int

const (
    CategoryAutobiography     Category = iota
    CategoryLargePrintRomance
    CategoryParticlePhysics
)

var validCategory = map[Category]bool{
    CategoryAutobiography:     true,
    CategoryLargePrintRomance: true,
    CategoryParticlePhysics:   true,
}

func (b *Book) SetCategory(category Category) error {
    if !validCategory[category] {
        return fmt.Errorf("unknown category %v", category)
    }
    b.category = category
    return nil
}

func (b Book) Category() Category {
    return b.category
}
```
Valid test — table-driven (3 kateqoriya):
```go
cats := []bookstore.Category{
    bookstore.CategoryAutobiography,
    bookstore.CategoryLargePrintRomance,
    bookstore.CategoryParticlePhysics,
}
for _, cat := range cats {
    err := b.SetCategory(cat)
    if err != nil { t.Fatal(err) }
    got := b.Category()
    if cat != got { t.Errorf("want category %q, got %q", cat, got) }
}
```
Invalid test: `b.SetCategory(999)` → error gözlə.

## Əsas terminlələr
- Validating Accessor — dəyəri yoxlayıb təyin edən metod (SetPriceCents)
- Value Receiver Tələsi — kopya dəyişir, original QALIR (SA4005)
- staticcheck / SA4005 — "ineffective assignment" lint xəbərdarlığı
- Automatic Dereferencing — pointer üzərindən b.Field Go özü *b.Field edir
- Always Valid Field — unexported sahə + getter/setter qoruması
- Accessor Methods — getter (Category) + setter (SetCategory)
- Getter Value / Setter Pointer — oxuyan kopya, yazan pointer alır
- Setter+Getter Test Cütü — bir valid/invalid cütü hər ikini yoxlayır
- cmpopts.IgnoreUnexported — unexported sahəli struct müqayisəsi üçün
- Always Valid Struct — unexported TİP + constructor
- Constructor (New) — yalnız VALID dəyərlər yarada bilən funksiya
- "Name istifadə etməmək" — unexported tip dəyəri adi istifadə olunur
- Map-as-Set — map[T]bool; missing key → false (zero value)
- validCategory Adı — natural oxunuş üçün təklik
- const / Konstant Qrupu — compiler yoxlamalı adlar
- http.StatusOK — standart kitabxana konstant nümunəsi
- type Category int — mənalı int tipləri
- iota — 0,1,2... avtomatik enumerator (Yunan kiçik hərfi)
- Table-Driven Category Testi — çox kateqoriya üçün loop

## Praktik nəticə
(1) Sahə dəyişəndə validasiya lazımdırsa — DİREKT təyinatı QADAĞA et:
unexported + accessor cütü. (2) Setter = pointer receiver MÜTLƏQ (kopyada
iş BITƏR); value receiver yazıntısı staticcheck SA4005 tutur. (3) Getter =
value receiver (dəyişmir). (4) cmp.Equal unexported görəndə panik edir —
cmpopts.IgnoreUnexported(T{}) əlavə et. (5) Struktun ÖZÜNÜ qoru: unexported
tip + New constructor; testlər adı istifadə etməz, dəyərlə metodlarla işlənər;
invalid obyekt YARATMAQ MÜMKÜNSÜZ. (6) Çoxlu icazəli dəyər → map[T]bool set;
missing = false magic. (7) Adları KONSTANT et: typo → compile xətası (200 →
http.StatusOK oxunaqlılıq + təhlükəsizlik). (8) Dəyərlər fərqli olsun YETER —
type + iota ilə avtomatik nömrələ; ortaya əlavə = renumber QƏYDİYYATI YOX.
(9) Validasiyalı metod = həmişə 2 test: valid (err==nil + dəyər yoxla) və
invalid (err gözlə). (10) User-input mənbəyi string seç — format fərziyyəsi
etmə (rəqəm gözləniləndə belə).

## Mənbə
Pages: 120-137 (PDF 121-138)
