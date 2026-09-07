# Chapter 4 — Happy Fun Books (Xoş Əyləncəli Kitablar)

## Bu fəsil nədən bəhs edir?

Yeni layihə: Happy Fun Books onlayn kitab mağazası. Data anlayışı: string
(mətn ardıcıllığı), numeric, boolean. Value vs variable vs type fərqi; var
declaration (`var title string`), assignment (=, string literal "It's literally
this string!"), type checking (`cannot use copies (type int) as type string`),
parametr ötürməsində type error; bool/float64 tipləri; zero values (0, "", false)
+ var vs `:=` stil qaydası; struct — composite type (Customer: Name, Email),
type definition (`type Customer struct`), documentation comments (pkg.go.dev),
exported/unexported identifiers (Böyük hərf = public, kiçik = private); core
package bookstore; compile-only test (`_ = bookstore.Book{...}` struct literal,
hər sahədən sonra VİRGÜL), "unfailable test" amma compile xətası = undefined:
bookstore.Book; Book struct (Title, Author, Copies).

## Əsas fikirlər

### 1. Data Növləri
- **String** — "a string of pearls" kimi ardıcıllıq; simvollar zənciri (kitab adı)
- **Numeric** — rəqəmsal (qiymət); **int** — tam ədədlər, **float64** — kəsr
  ("just read it as fraction")
- **Boolean** — true/false, on/off məlumat

### 2. Value / Variable / Type
- **Value** — data parçası ("Hello, world"); value-un TİPİ var
- **Variable** — yaddaşda yer tutan AD; tipi NƏ saxlaya biləcəyini müəyyən edir
- **Type** = "kind of thing" — AVTOMATİK düzgünlük yoxlaması üçün

### 3. var + Assignment
```go
var title string       // "Please create a variable called title of type string"
var copies int
title = "For the Love of Go"   // sol: ad, sağ: dəyər (string LITERAL)
copies = 99                    // int literal
fmt.Println(title, copies)
```
**Literal** — kodda birbaşa yazılmış dəyər: "It's literally this string!"

### 4. Type Checking — Compiler Yoldaşımız
```go
title = copies
// cannot use copies (type int) as type string in assignment
```
Başqa yer: funksiya parametri:
```go
func printTitle(title string) { fmt.Println("Title: ", title) }
printTitle(copies)
// cannot use copies (type int) as type string in argument to printTitle
```
Go "types in a twist" olsa DƏRHAL xəbər verir — mürəkkəb proqramlarda xilaskar.

### 5. Zero Values + var vs :=
```go
var x int
fmt.Println(x)     // 0 — DEFAULT dəyər
```
| Tip | Zero value |
|-----|-----------|
| int, float64 (numeric) | 0 |
| string | "" (boş string) |
| bool | false |

**Stil qaydası:**
- Zero dəyər LAZIMDIRSA → `var x int`
- Başqa dəyərlə başlayırsa → `x := 1` (declare + assign bir sətirdə)

### 6. Struct — Composite Type
Kitabın datası: title, author, price, ISBN... Hər birini AYRI tutmaq əvəzinə —
kitab TƏK vahid kimi: funksiyaya ötür, DB-də saxla.
**Composite type** — bir neçə dəyərdən təşkil olunmuş vahid.
**Struct** ("structured record") — ən faydalı composite:
```go
// Customer represents information about a customer.
type Customer struct {
    Name  string
    Email string
}
```
- **Field** — struct-ın daxili dəyəri (istənilən tip, o cümlədən digər struct!)
- testCase (kalkulyator testlərindən) də struct idi

### 7. type Definition + Doc Comments
```go
type Customer struct { ... }
```
- `type` — yeni tip yarat (var-ın tip versiyası): "Please create a type
  called... and here are the details"
- Ardınca AD + type literal (struct + field list)
- **// şərhi** — documentation comment; hər tip üçün YAZ (good practice);
  GitHub-a publish → pkg.go.dev-də avtomatik dokumentasiya
  (nümunə: github.com/bitfield/script)

### 8. Exported vs Unexported Identifiers
- **Böyük hərflə başlayır** (Customer, Book, Add) → **EXPORTED** = public —
  paketdən kənarda görünür (testdə də!)
- **Kiçik hərflə başlayır** (customer) → **UNEXPORTED** = private — yalnız paket
  daxilində
- Funksiya, dəyişən, tip — hamısı üçün eyni qayda
- Paketin niyyəti: böyük hərf = açıq API; kiçik = daxili detallar

### 9. Core Package — bookstore
Paket = əlaqəli kodun vahidi. Core paket — problemin ƏSAS funksionallığı.
Ad seçimi: TƏK söz, ideal QISA, tam təsviv edən → **bookstore**.

### 10. Compile-Only Test — Struct-u Testlə Yaratmaq
Dilemma: Book tipi lazımdır, amma non-test kod YAZMAQ olmaz (TDD)! Nəysə ki
struct-un "davranışı" yoxdur... **Həll — compile-only test:**
```go
func TestBook(t *testing.T) {
    t.Parallel()
    _ = bookstore.Book{
        Title:  "Spark Joy",
        Author: "Marie Kondo",
        Copies: 2,
    }
}
```
- **Struct literal:** `Book{Title: "...", Author: "...", Copies: 2}` — tip adı
  + mötərizə + field:dəyər cütləri; HƏR sahədən sonra VİRGÜL (sonuncu daxil!) —
  struct DEFINİSİYASINDA virgül YOXDUR, literalda VAR
- **Niyə `_`?** Literal tək başına statement DEYİL; adlı dəyişənə assign etsək
  "unused variable" xətası; `_` = "sadəcə qeyd etmək istəyirəm, istifadə
  etməyəcəyəm"
- Bu test FAIL EDE BİLMƏZ ("unfailable") — amma hələ COMPILE OLUMUR:
```
./bookstore_test.go:10:6: undefined: bookstore.Book
FAIL bookstore [build failed]
```
Məqsəd məhz budur: `undefined: bookstore.Book` → indi Book-u YAZ.

### 11. Book Struct — Nəticə
```go
// Book represents information about a book.
type Book struct {
    Title  string
    Author string
    Copies int
}
```
**Kapital hərf tələbi:** bookstore_test pakətindən istifadə → Book EXPORTED
olmalı; `book` yazsan — unexported → testdə GÖRÜNMƏZ.

**TDD prinsipi:** "The test is the definition of what's okay!" — müəllifin
versiyası ilə fərqli olsan, amma test keçir — DÜZGÜNDÜR.

## Əsas terminlələr
- String / int / float64 / bool — əsas data tipləri
- Value / Variable — data parçası / adlı yaddaş yeri
- Type — "kind of thing"; avtomatik yoxlama vahidi
- var — dəyişən yarat (+ tip bildir)
- = Assignment — sol: ad, sağ: dəyər
- Literal — koddan birbaşa dəyər (string/int/struct literal)
- Type Checking — compiler-in tip uyğunsuzluğu yoxlaması
- Zero Value — tipin defaultu: 0 / "" / false
- var vs := — zero dəyər istəyirsən var; dəyər ilə başlayırsa :=
- Composite Type — bir neçə parçadan təşkil vahid
- Struct / Field — strukturlaşmış qeyd / sahələri
- type Definition — type Ad struct {...}
- Documentation Comment — pkg.go.dev-də çıxan // şərhi
- Exported / Unexported — Böyük hərf public / kiçik hərf private
- Core Package — layihənin əsas paketi (bookstore)
- Struct Literal — Book{Title: "...", ...} — sahələrdə VİRGÜL
- Compile-Only Test — `_ = T{...}` — yalnız compile tələb edən test
- "The test is the definition of what's okay" — test = spesifikasiya

## Praktik nəticə
(1) Tip seçimi: mətn→string, tam→int, kəsr→float64, bayraq→bool. (2) var =
zero dəyər istəyirsən; `:=` = başlanğıc dəyər ilə — bu, stil QAYDADIR. (3)
Tip xətalarını compiler-dən qəbul et: have/want mesajları sənin tərcüməçindir.
(4) Əlaqəli data → struct: `type Ad struct` + field list; hər tipə // doc
şərhi. (5) Paketdən kənar istifadə → BÖYÜK hərf; yalnız daxili → kiçik. (6)
Core paketə qisa, tək, tam ad ver (bookstore). (7) Strukt-ları TDD ilə yarat:
compile-only test (`_ = bookstore.Book{...}`) → undefined xətası → struct yaz →
PASS. (8) Struct literalda hər sahədən sonra VİRGÜL; struct definisiyasında
YOX. (9) Literal istifadə olunmayacaqsa `_`-ə assign et (unused xətasından
qaçın). (10) "Test = okay tərifi": müəllif həlli ilə fərqlənsən amma test keçirsə
— düzgündür.

## Mənbə
Pages: 50-60 (PDF 51-61)
