# Chapter 6 — Money converter: CLI around an HTTP call (səh. 203-281)

## Bu chapter nədən bəhs edir?

Pul konvertoru CLI: float dəqiqlik problemləri, Decimal tipli öz dəyər obyekti,
Currency value object, Stringer, ECB API-yə HTTP çağırışı, XML decode, dependency
injection (funksiya vs interfeys), httptest mock server, timeout.

## Əsas fikirlər

### 1. Float pul üçün QADAĞANDIR
```go
fmt.Println(math.Sin(math.Pi)) // 1.2246467991473512e-16 — 0 olmalı idi!
fmt.Printf("%.f", float32(123_456_789)) // 123456792 — precision itir
```
Float64 15-16 rəqəm dəqiqliyi, float32 ~7. Pul = **fixed precision**
(subunits: sentlər), float arifmetikası sentləri itirir.

### 2. Decimal value object
```go
type Decimal struct {
    subunits  uint64 // bütün rəqəmlər baytlarla ifadə olunur (10^əsaslı)
    precision byte   // neçə rəqəm onluq hissədə
}
```
**ParseDecimal** string-dən Decimal-ə:
```go
major, minor, found := strings.Cut("32.50", ".") // "32", "50", true
```
- `strings.Cut` → bir ayırıcıda iki hissə; `strings.Split` → slice (lazımsız
  paylar yaradır)
- Precision məhdudiyyəti (maxDecimal 10^12) + custom Error tipi:
  `type Error string; func (e Error) Error() string`

### 3. Currency value object
```go
type Currency struct {
    code      string // ISO-4217 3 hərfli kod (USD, EUR)
    precision byte   // Yaponiya yen: 0; USD: 2
}
```
Valyuta dəqiqliyi fərqlidir — konvertorda hədəf valyutanın dəqiqliyinə
uyğunlaşdırma (`applyExchangeRate`) lazımdır.

### 4. NewAmount + validate
```go
func NewAmount(quantity Decimal, currency Currency) (Amount, error) {
    if quantity.precision > currency.precision {
        return Amount{}, ErrInvalidPrecision // dəqiqlik uyğunsuzluğu
    }
    return Amount{quantity: quantity, currency: currency}, nil
}
```
Value object constructor-u **kənar murakibə buraxmır** — yanlış vəziyyət
yaranmır (invalid state unrepresentable).

### 5. Test helper-lər (mustParse)
```go
func mustParseCurrency(t *testing.T, code string) money.Currency {
    t.Helper() // xəta hesabatında bu sətri göstərmə
    curr, err := money.ParseCurrency(code)
    if err != nil { t.Fatal(err) }
    return curr
}
```
Cədvəl testində giriş qurmağın təkrarını aradan qaldırır — testlər **yalnız
davranışı** göstərir.

### 6. DI: funksiya vs interfeys
```go
// Variant A — funksiya parametri:
func Convert(amount Amount, to Currency,
    rates func(from, to Currency) (ExchangeRate, error)) (Amount, error)

// Variant B — interfeys (seçilən):
type exchangeRates interface {
    FetchExchangeRate(source, target Currency) (money.ExchangeRate, error)
}
func Convert(amount Amount, to Currency, rates exchangeRates) (Amount, error)
```
**Go fəlsəfəsi:** interfeyslər **istehlakçı tərəfdə** təyin olunur (implicit
implementation) — mock edilmə ehtiyacı yoxdursa funksiya bəsdir.

### 7. ECB HTTP çağırışı
```go
func (e EuroCentralBank) FetchExchangeRate(source, target money.Currency) (money.ExchangeRate, error) {
    resp, err := http.Get(ecbDailyRatesURL) // istehsalatda: öz Client + Timeout
    if err != nil { return money.ExchangeRate{}, fmt.Errorf("...: %w", err) }
    defer resp.Body.Close()
    if err := checkStatusCode(resp.StatusCode); err != nil { ... }
    return readRateFromResponse(source.CurrencyCode(), target.CurrencyCode(), resp.Body)
}
```
**Status kod semantikası:**
- 1xx — informational; 2xx — success; 3xx — redirect
- 4xx (`clientErrorClass`) — BİZİN xətamız (soruğu səhf)
- 5xx (`serverErrorClass`) — ONLARIN xətası
- 2XX/3XX xarici hamısı xəta sayılır

**XML decode** (ECB XML qaytarır):
```go
decoder := xml.NewDecoder(respBody)
var envelope envelope
if err := decoder.Decode(&envelope); err != nil { ... } // POINTER mütləq!
```

### 8. httptest mock server
```go
ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, `<?xml version="1.0"...><Cube><Cube time="...">...`) // saxta cavab
}))
defer ts.Close()
// ECB URL-i testdə ts.URL ilə əvəz olunur (konstruktorda path parametri)
```
Real şəbəkə YOX — lokal server; success + timeout + status xətaları üçün ayrı
ssenarilər.

### 9. Timeout + DefaultClient təhlükəsi
```go
var DefaultClient = &Client{} // http.Get bunu işlədir — timeout YOX!
```
Hər kəs `DefaultClient`-ı dəyişə bilər — global mutasiya. Həll:
```go
func NewBank(timeout time.Duration) EuroCentralBank {
    return EuroCentralBank{path: url, client: &http.Client{Timeout: timeout}}
}
```
Timeout çatanda `*url.Error` qayıdır — `errors.As` ilə unwrap.

### 10. Layihə strukturu (alternative tree)
```
moneyconverter/
├── money/    → domain (Amount, Currency, Convert — xalis)
├── ecbank/   → xarici API adapteri (XML, HTTP)
└── main.go   → CLI (flag, Stringer, main)
```
Domain paketi xarici mənbəyi BİLMİR — rates interfeysi ilə təcrid olunur.

## Əsas terminlər

- Floating-Point Precision (üzən nöqtə dəqiqliyi)
- Value Object (dəyər obyekti)
- Subunits (alt vahidlər — sentlər)
- strings.Cut
- ISO-4217 (valyuta kodu standartı)
- Dependency Injection (asılılıq inyeksiyası)
- httptest.NewServer (mock server)
- xml.Decoder
- DefaultClient təhlükəsi

## Praktik nəticə

- Pul üçün float YOX — subunits (uint64) + precision
- Domain (money) və infrastruktur (ecbank) paketləri ayır
- Xarici çağırışları interfeys arxasında saxla → httptest ilə test
- Timeout-lu öz `http.Client` istifadə et — DefaultClient-a etibar etmə
- Səhv status kodlarını siniflər üzrə (4xx/5xx) ayır

## Mənbə

Pages: 203-281 (Chapter 6, Learn Go with Pocket-Sized Projects)
