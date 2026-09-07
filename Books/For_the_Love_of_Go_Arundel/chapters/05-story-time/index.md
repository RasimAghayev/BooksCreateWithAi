# Chapter 5 — Story time (Nağıl Vaxtı)

## Bu fəsil nədən bəhs edir?

User stories anlayışı ("As a customer, I want to..." formatı), core stories-in
müəyyənləşdirilməsi (Buy a book / List all books / See details), davranış-əsaslı
dizayn; struct dəyişənləri (var b bookstore.Book, `:=` qısa forması), dot
notation (b.Copies — oxu və yaz), TestBuy testi, informativ fail mesajı sənəti,
null implementation (return Book{}), decrement/increment operatorları (--, ++,
-=, +=), Buy implementasiyası, test coverage (go test -cover,
-coverprofile, go tool cover -html, VS Code Toggle Test Coverage, boz sətirlər
= statement deyil), test-last development təcrübəsi (zero copies bug-ı → -1!),
error qaytarma əlavəsi + signature dəyişmə + test uyğunlaşdırma, coverage 75% →
TestBuyErrorsIfNoCopiesLeft → 100%, "covered vs tested" (TestBuyTrivial — 100%
useless), juking the stats, os.Open error-line nümunəsi, 80-90% real hədəf.

## Əsas fikirlər

### 1. User Stories — İstifadəçi Görüşü
**User story** — istifadəçi və istədiyi əməliyyat haqqında kiçik nağıl:
- "As a customer, I want to see a sample preview of the book in my browser
  so that I can figure out if I'd like to read it."
- "As a bookstore owner, I want to see my total revenue for the month, so
  that I'll know if I can afford pizza."

Feature ≠ user story: feature = "proqram nəyi edə bilər"; bəzi feature-lar
istifadəçiyə HEÇ NƏ vermir (buna baxmayaraq vendorlar əlavə edir). Story =
istifadəçi perspektivi.

### 2. Core Stories — Müzakirəsiz Zəruri Olanlar
Bütün mümkün story-ləri siyahılamaq YOX (uzun + prioritet bilinməz). Sual
tersinə: **"HANSI story-ləri buraxmaq MÜMKÜNSÜZDÜR?"** — productsız məhsul
olmaz:
1. **Buy a book**
2. **List all available books**
3. **See details of a book**

Bunlar olmadan mağaza "mağaza" deyil. Biznes məntiqi: product viable olana qədər
pul qazandırmır → ƏVVƏL critical story-lər.

### 3. "Buy a Book" → Davranış → Test
Story-ni SADƏLƏŞDİR: ödənişləri SONRAYA saxla; yalnız STOCK CONTROL baxımından:
"İstifadəçi kitab alanda mövcud nüsxə sayı AZALIR."

Bu, testin verbal təsviridir. Funksiya dizaynı: `Buy(b Book) Book` — alır Kitabı,
qaytarır dəyişdirilmiş Kitabı.

### 4. Struct Dəyişənləri + Dot Notation
```go
var b bookstore.Book            // var + struct tipi (built-in kimi!)
b := bookstore.Book{            // qısa forma + literal
    Title:  "Nicholas Chuckleby",
    Author: "Charles Dickens",
    Copies: 8,
}
b.Copies                        // dot notation: sahəyə ÇIXIŞ
b.Copies = 7                    // sahəyə YAZ (adi dəyişən kimi)
```
Testdə bookstore pakətindən kənardayıq → fully-qualified `bookstore.Book`.

### 5. TestBuy + İnformativ Mesaj Sənəti
```go
func TestBuy(t *testing.T) {
    t.Parallel()
    b := bookstore.Book{
        Title:  "Spark Joy",
        Author: "Marie Kondo",
        Copies: 2,
    }
    want := 1
    result := bookstore.Buy(b)
    got := result.Copies
    if want != got {
        t.Errorf("want %d copies after buying 1 copy from a stock of 2, got %d", want, got)
    }
}
```
**Fail mesajı səviyyələri:** "failed" — FAYDASIZ; "unexpected number of copies"
— az; "want X, got Y" — hələ çatmır; **DOĞRU:** "2 nüsxədən 1 aldıqdan sonra
1 gözlənirdi, alınan: Y" — KONTEKST + gözlənti + fakt. Səbəb: bir gün izahsız
fail alanda hamı sənə baxır — test DƏRHAL nə düzəltmək lazım olduğunu söyləməlidir.

### 6. Null Implementation + Operatorlar
Compile xətası: `undefined: bookstore.Buy` → minimum kod:
```go
func Buy(b Book) Book {
    return Book{}       // struct ZERO VALUE = boş literal!
}
```
Fail: "want 1 copies ... got 0" — gözlənilən RED. İndi operatorlar:
```go
b.Copies--     // decrement: -1
b.Copies++     // increment: +1
b.Copies -= 5  // -5
b.Copies += 10 // +10
```
Real implementasiya:
```go
func Buy(b Book) Book {
    b.Copies--
    return b
}
```
**DİQQƏTƏ DƏYƏR müşahidə:** funksiya testdən QISAdır — NORMAL. Test bütün
davranışı TƏSVİF edir, funksiya yalnız icra edir. "Kodun ~50%-i testdirsə —
provably right; azı — test çatışır."

### 7. Test Coverage Alətləri
```bash
go test -cover
# PASS, coverage: 100.0% of statements

go test -coverprofile=coverage.out
go tool cover -html=coverage.out    # browser: yaşıl/ QIRMIZI /boz
```
- VS Code: "Go: Toggle Test Coverage in Current Package" — editor içində
- **Boz sətirlər** = type definisiyaları, const-lar — compiler göstərişləri,
  binary-də obyekt kodu YOXDUR → coverage STATİSTİKASINA DAHİL DEYİL
- TDD ilə getsən, uncovered sistem kodu OLMAZ (hər sətir testdən doğub)

### 8. Test-Last Development — Zero Copies Fəlakəti
Ssenari: təcrübəsiz developer feature əlavə edir test YAZMADAN. Bug: Spark Joy
satılıb bitib (Copies=0) → Buy -1 qaytarır! Qəzəbli müştərilər: "mənfi sayda
kitab aldıq, bizə pul borclusuz!"

**Düzəliş — error əlavə et (testdən SONRA!):**
```go
func Buy(b Book) (Book, error) {
    if b.Copies == 0 {
        return Book{}, errors.New("no copies left")
    }
    b.Copies--
    return b, nil
}
```
Zəncir: signature dəyişdi → testdə `result, err := bookstore.Buy(b)` → err
yoxla (`t.Fatal(err)`) → PASS. Amma...

### 9. Coverage 75% → Yeni Davranış Testi
```
coverage: 75.0% of statements
```
HTML: QIRMIZI sətir = `return Book{}, errors.New("no copies left")` — error
qaytaran return HEÇ BİR test tərəfindən İCRA OLUNMUR!

**Risk:** kod bu gün düzgündür, amma sabah kimsə dəyişər, testlər KEÇƏR, xəta
gizli qalar. Həll — yeni DAVRANIŞ testi:
```go
func TestBuyErrorsIfNoCopiesLeft(t *testing.T) {
    t.Parallel()
    b := bookstore.Book{
        Title:  "Spark Joy",
        Author: "Marie Kondo",
        Copies: 0,
    }
    _, err := bookstore.Buy(b)
    if err == nil {
        t.Error("want error buying from zero copies, got nil")
    }
}
```
→ 100% coverage qayıtdı. **Qayda:** `if` = 2 mümkün path = 2 DAVRANIŞ = 2 test.

### 10. "Covered" ≠ "Tested" + Juking the Stats
```go
func TestBuyTrivial(t *testing.T) {
    t.Parallel()
    bookstore.Buy(bookstore.Book{})
    bookstore.Buy(bookstore.Book{Copies: 1})
}
```
100% coverage — amma HEÇ NƏ yoxlamır: "Covered" ≠ "tested". TestBuyTrivial
coverage-ı "juking" (aldatma) edir. "Test behaviours, not functions" — tələ.

**Real dünyada 100% NADİRDIR:**
```go
f, err := os.Open("some file")
if err != nil {
    return err      // bu sətir adətən QIRMIZI qalır
}
```
Test etmək olar (mövcud olmayan faylla) — amma nə YOXLAYIRIQ? Yalnız if-i.
Real davranış yoxdursa — test DƏYMƏZ; göz ilə düzgünlük görürük.

**Hədəf:** xüsusi RƏQƏM yox — "100% of important behaviours" (100% of lines
DEYİL). Sadə `if err != nil` blokları üçün 100%-dən aşaqlıq NORMAL. 80-90%
düzgün rəqəmdir.

## Əsas terminlələr
- User Story — istifadəçi perspektivli qısa nağıl
- Core Stories — məhsulun varlığı üçün zəruri minimum story dəsti
- Feature vs Story — proqram bacarığı vs istifadəçi dəyəri
- Behaviour-First — "how does the program need to behave?" sualı
- Struct Variable — var b bookstore.Book / b := Book{...}
- Dot Notation — b.Copies (oxu + yaz)
- ++ / -- / += / -= — increment/decrement/compound operatorlar
- Zero Value of Struct — Book{} (boş literal)
- Null Implementation (struct) — return Book{} → FAIL sübutu
- İnformativ Fail Mesajı — kontekst + gözlənti + fakt bir mesajda
- Test Coverage — sətirlərin testlər tərəfindən icra %-i
- go test -cover / -coverprofile / go tool cover -html
- Boz Sətirlər — type/const — coverage-a daxil deyil
- Test-Last Development — kod əvvəl, test sonra → uncovered risk
- TestBuyErrorsIfNoCopiesLeft — error davranışı testi
- "Covered" vs "Tested" — icra ≠ yoxlama
- TestBuyTrivial — juking the stats nümunəsi
- 80-90% — real coverage hədəfi; 100% = yalnız vacib davranışlar

## Praktik nəticə
(1) Paketə başlayanda: user story-lər yaz; "buraxmaq mümkünsüz" olanları CORE
kimi seç, onlardan başla. (2) Story → SADƏLƏŞDİR (ödənişsiz stock-control) →
DAVRANIŞ ifadəsi → test. (3) Struct dəyişəni: var (zero) yaxud `:=` literal;
sahələrə dot notation. (4) Fail mesajına KONTEKST: "want X after buying from
stock of Y, got Z". (5) Null impl: `return Book{}` → fail → sonra `b.Copies--`.
(6) Kodun yarısı testdirsə — sağlam; funksiyanın testdən qısalığı NORMAL. (7)
Coverage alətləri: -cover, -coverprofile, -html, editor highlight. (8) if
görsən — 2 path = 2 test; error qaytaran path-i YOXSA coverage göstərəcək. (9)
Signature dəyişəndə testi də dəyiş: `result, err :=` + err yoxla. (10) Rəqəm
OVU YOX: "100% of important behaviours"; sadə err-nil bloklarını gözlə görmək
kifayətdir.

## Mənbə
Pages: 61-72 (PDF 62-73)
