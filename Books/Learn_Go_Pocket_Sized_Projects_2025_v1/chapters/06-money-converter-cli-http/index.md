# Chapter 6 — Money Converter: CLI around an HTTP Call (səh. 203-281)

## Bu fəsil nədən bəhs edir?

Pul konvertoru CLI: float dəqiqlik problemləri və Decimal tipli həll, ISO-4217
currency, XML parsinq (ECB API), HTTP status kodları, dependency injection
(interfeys YOX — funksiya), httptest ilə mock, timeout (http.Client), errors.As
ilə tip yoxlaması, Stringer.

## Əsas fikirlər

### 1. Float dəqiqliyi — niyə Decimal?
```go
fmt.Printf("%.f", float32(123_456_789))  // → 123456792  (7-dən sonra itir!)
fmt.Println(math.Sin(math.Pi))            // → 1.22e-16  (0 olmalı idi)
```
- float32 ~7 rəqəm dəqiqlik, float64 ~15 — PULA YARAMIR
- **Həll:** fixed precision integer arifmetikası:
```go
type Decimal struct {
    subunits  int64  // "1.52" → 152
    precision byte   // → 2
}
```
- `ParseDecimal("1.52")` → {subunits:152, precision:2}; `simplify()` — sondakı
  sıfırları atır (32.0 → 32, 320 → 320 qalır)
- `pow10(power byte) int64` — 10^p tez hesabla (switch 0-3, math.Pow qalan)

### 2. Currency + Amount (value object)
```go
type Currency struct {
    code       string  // "USD" (ISO-4217)
    precision  byte    // dollar = 2, BHD = 3
}
type Amount struct {
    quantity Decimal
    currency Currency
}

func NewAmount(quantity Decimal, currency Currency) (Amount, error) {
    switch {
    case quantity.precision > currency.precision:
        return Amount{}, ErrTooPrecise   // 0.0001 cent yaratmaq olmaz
    case quantity.subunits > maxDecimal:
        return Amount{}, ErrTooLarge
    }
    // ...
}
```

### 3. Konversiya — applyExchangeRate
```go
func multiply(d Decimal, r ExchangeRate) Decimal {
    dec := Decimal{
        subunits:  d.subunits * r.subunits,   // int arifmetika — DƏQİQ
        precision: d.precision + r.precision,
    }
    dec.simplify()
    return dec
}
```
- Xəta yox, dəqiqlik itirmə yox — kəsr hissələr int-də yaşayır

### 4. CLI: flag + argument
```go
from := flag.String("from", "", "source currency, required")
to   := flag.String("to", "EUR", "target currency")
flag.Parse()
value := flag.Arg(0)  // pozisional arqument
if value == "" {
    flag.Usage()
    os.Exit(1)
}
```
- **Flag:** adlı (-from), default-lu, ixtiyari; **Argument:** anonim, sıralı,
  adətən məcburi

### 5. Stringer — oxunaqlı çap
```go
func (d Decimal) String() string {
    if d.precision == 0 { return fmt.Sprintf("%d", d.subunits) }
    centsPerUnit := pow10(d.precision)
    format := "%d.%0" + strconv.Itoa(int(d.precision)) + "d"  // "%d.%02d"
    return fmt.Sprintf(format, d.subunits/centsPerUnit, d.subunits%centsPerUnit)
}
```

### 6. Dependency injection — "interfaces are discovered, not designed"
**Variant A (Java-tipli):** interfeys + struct metodu.
**Variant B (Go-tipli, seçilən):** FUNKSİYA parametri:
```go
func Convert(amount Amount, to Currency, rates exchangeRates) (Amount, error)
// exchangeRates — KİÇİK interfeys: FetchExchangeRate(from, to) (ExchangeRate, error)
```
- Konvertasiya paketi bankı TANIMIR; ecb paketi ayrıca təmin edir
- İnterfeys istehlakçı tərəfdə təyin olunur ("discovered") — Go fəlsəfəsi

### 7. ECB XML API çağırışı
```go
type Client struct{ path string }  // path — mock üçün INJEKTƏ olunur

func (c Client) FetchExchangeRate(source, target money.Currency) (money.ExchangeRate, error) {
    resp, err := http.Get(c.path)  // təhlükəli — hamıya eyni DefaultClient!
    if err != nil {
        return 0., fmt.Errorf("%w: %s", ErrCallingBank, err)
    }
    defer resp.Body.Close()
    if err := checkStatusCode(resp.StatusCode); err != nil { return 0., err }
    return readRateFromResponse(source.Code(), target.Code(), resp.Body)
}
```
- `xml.NewDecoder(resp.Body)` + struct tag-lər:
  \t\t`xml:"Cube Cube Currency Code Rate"`
- Status kodları: 1xx info / 2xx uğur / 3xx redirect / 4xx client / 5xx server;
  4xx → ErrClientResponse, 5xx → ErrServerResponse

### 8. httptest ilə mock server
```go
ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, `<?xml ...test data...?>`)
}))
defer ts.Close()
ecb := Client{path: ts.URL}  // real URL əvəzinə test server
```

### 9. Timeout — DefaultClient təhlükəsi
```go
// http.Get = http.DefaultClient.Get — QLOBAL, hər kəsə təsir edir
// DÜZGÜN:
rates := ecbank.NewBank(30 * time.Second)  // öz Client-i + Timeout
```
- `url.Error` yoxlama: `errors.As(err, &urlErr)` + `urlErr.Timeout()`

## Əsas terminlər

- Fixed-Point Arithmetic (sabit nöqtəli arifmetika)
- Value Object (dəyər obyekti)
- ISO-4217 (valyuta kodu standartı)
- Dependency Injection (asılılıq inyeksiyası)
- httptest (HTTP mock server)
- Sentinel Error + errors.Is / errors.As
- struct tag (xml:"...")

## Praktik nətitə

- PUL üçün float QADAĞAN — int subunits + precision
- Xarici API cavabı etibarsızdır: status yoxla + sentinel error-la wrap
- http.Get yox, öz http.Client + Timeout
- Mock: httptest.NewServer + injektə olunan URL

## Mənbə

Pages: 203-281 (Chapter 6, Learn Go with Pocket-Sized Projects)
