# Chapter 4 — Parsing JSON Requests

## Bu fəsil nədən bəhs edir?

Client sorğularının oxunması: `json.Decoder` ilə decode, xəta triage-ı
(errors.Is/errors.As), input məhdudlaşdırılmaları (unknown fields, body size,
single value), `json.Unmarshaler` custom decode və Validator paketi ilə
business-rule yoxlamaları.

## Əsas fikirlər

### 1. json.Decoder ilə decode
**Nədir:** Sorğu body-sini stream kimi oxuyub Go obyektinə çevirmək.
HTTP body üçün `json.Unmarshal()`-dan üstündür (daha az yaddaş ~80% az B/op,
daha az kod).

**Kitabdan kod nümunəsi:**
```go
var input struct {
    Title   string       `json:"title"`
    Year    int32        `json:"year"`
    Runtime data.Runtime `json:"runtime"`
    Genres  []string     `json:"genres"`
}
err := json.NewDecoder(r.Body).Decode(&input)
```

**Sub-kod izahı:**
- `json.NewDecoder(r.Body)` → body-dən oxuyan Decoder
- `Decode(&input)` → hədəf **non-nil pointer** olmalıdır — yoxsa
  `json.InvalidUnmarshalError`; struct field-ləri export olunmalıdır
- Tag uyğunluğu: dəqiq match > case-insensitive match; uyğun gəlməyən
  key-lər **səssəcə ignore** olunur
- `r.Body` bağlamaq lazım deyil — `http.Server` avtomatik qapatır

**JSON → Go destination cədvəli:**
| JSON type | Supported Go types |
|---|---|
| boolean | bool |
| string | string |
| number | int*, uint*, float*, rune |
| array | array, slice |
| object | struct, map |

**Zero value problemi:** `{"year":0}` və year-un olmaması eyni nəticə verir —
fərqi sonra pointer/özlük tiplərlə ayırd edəcəyik (GET partial update
fəslində).

### 2. Decode xətalarının triage-ı — readJSON() helper
**Nədir:** Decode-in 5 xəta tipini tutub aydın mesajlarla əvəz etmək.

**Xəta tipləri:**
| Error | Səbəb |
|---|---|
| `json.SyntaxError`, `io.ErrUnexpectedEOF` | JSON sintaksis xətası |
| `json.UnmarshalTypeError` | dəyər tipi uyğunsuz |
| `json.InvalidUnmarshalError` | hədəf pointer deyil (proqram xətası) |
| `io.EOF` | body boşdur |

**Kitabdan kod nümunəsi (tam readJSON):**
```go
func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst any) error {
    maxBytes := 1_048_576
    r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

    dec := json.NewDecoder(r.Body)
    dec.DisallowUnknownFields()

    err := dec.Decode(dst)
    if err != nil {
        var syntaxError *json.SyntaxError
        var unmarshalTypeError *json.UnmarshalTypeError
        var invalidUnmarshalError *json.InvalidUnmarshalError
        var maxBytesError *http.MaxBytesError

        switch {
        case errors.As(err, &syntaxError):
            return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)
        case errors.Is(err, io.ErrUnexpectedEOF):
            return errors.New("body contains badly-formed JSON")
        case errors.As(err, &unmarshalTypeError):
            if unmarshalTypeError.Field != "" {
                return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
            }
            return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)
        case errors.Is(err, io.EOF):
            return errors.New("body must not be empty")
        case strings.HasPrefix(err.Error(), "json: unknown field "):
            fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
            return fmt.Errorf("body contains unknown key %s", fieldName)
        case errors.As(err, &maxBytesError):
            return fmt.Errorf("body must not be larger than %d bytes", maxBytesError.Limit)
        case errors.As(err, &invalidUnmarshalError):
            panic(err)
        default:
            return err
        }
    }

    err = dec.Decode(&struct{}{})
    if err != io.EOF {
        return errors.New("body must only contain a single JSON value")
    }
    return nil
}
```

**Sub-kod izahı:**
- `http.MaxBytesReader(w, r.Body, 1MB)` → DoS-a qarşı body limiti;
  aşarsa `*http.MaxBytesError`
- `dec.DisallowUnknownFields()` → tanınmayan key xəta yaradır ("json: unknown
  field X" — hələ ayrı error type yoxdur, string prefix ilə tutulur)
- `errors.As()` → konkret tip; `errors.Is()` → sentinel xəta yoxlaması
- İkinci `dec.Decode(&struct{}{})` → yalnız `io.EOF` qayıdırsa body-də
  tək JSON var; əks halda "single JSON value" xətası
- `panic(err)` (InvalidUnmarshalError) → developer səhvidir, "fail fast"

**Panic vs return (Additional Information):**
- **Gözlənilən xətalar** (DB timeout, pis input) → return + graceful handle
- **Gözlənilməz xətalar** (developer mantıq səhvi) → panic qəbuledilər;
  stdlib də belə edir (out-of-bounds, closed channel). Qayda: "panic means
  something went unexpectedly wrong... fail fast on errors that shouldn't
  occur during normal operation"

### 3. badRequestResponse helper
```go
func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
    app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}
```
- İxtisaslaşmış helper-lər bütün endpoint-lərdə tutarlı xəta formatı saxlayır

### 4. Custom decoding — json.Unmarshaler
**Nədir:** Encoding-in əksi: `UnmarshalJSON([]byte) error` method-u decode-i
əvəz edir. `"107 mins"` string-ini `Runtime`-a çevirmək üçün.

**Kitabdan kod nümunəsi:**
```go
var ErrInvalidRuntimeFormat = errors.New("invalid runtime format")

type Runtime int32

func (r *Runtime) UnmarshalJSON(jsonValue []byte) error {
    unquotedJSONValue, err := strconv.Unquote(string(jsonValue))
    if err != nil {
        return ErrInvalidRuntimeFormat
    }
    parts := strings.Split(unquotedJSONValue, " ")
    if len(parts) != 2 || parts[1] != "mins" {
        return ErrInvalidRuntimeFormat
    }
    i, err := strconv.ParseInt(parts[0], 10, 32)
    if err != nil {
        return ErrInvalidRuntimeFormat
    }
    *r = Runtime(i)
    return nil
}
```

**Sub-kod izahı:**
- **Pointer receiver MÜTLƏQDİR** — receiver-i modifikasiya edirik; value
  receiver yalnız nüsxəni dəyişərdi
- `strconv.Unquote()` → dırnaqları açır (`"107 mins"` → `107 mins`)
- Format yoxlaması: 2 hissə, ikincisi "mins"
- `*r = Runtime(i)` → pointer-deref ilə dəyəri yaz
- Sentinel xəta `ErrInvalidRuntimeFormat` → handler tərəfindən
  `errors.Is` ilə tutula bilər

### 5. Validator paketi — business rule yoxlamaları
**Nədir:** `internal/validator` — sade, lakin güclü validation qatı.

**Kitabdan kod nümunəsi:**
```go
// internal/validator/validator.go
type Validator struct {
    Errors map[string]string
}

func New() *Validator { return &Validator{Errors: make(map[string]string)} }
func (v *Validator) Valid() bool { return len(v.Errors) == 0 }
func (v *Validator) AddError(key, message string) {
    if _, exists := v.Errors[key]; !exists {
        v.Errors[key] = message
    }
}
func (v *Validator) Check(ok bool, key, message string) {
    if !ok { v.AddError(key, message) }
}

func PermittedValue[T comparable](value T, permittedValues ...T) bool { ... }
func Matches(value string, rx *regexp.Regexp) bool { ... }
func Unique[T comparable](values []T) bool { ... }
```

**Sub-kod izahı:**
- `Errors map[string]string` → field adı → mesaj; ilk xəta qalib
  (AddError mövcud key-i üstə yazmır)
- `Check(ok, key, msg)` → koşullu əlavə — handler-də bəyanətvari (declarative)
  görünüş
- Generic `Unique[T comparable]` → slice dəyərləri unikaldırmı (map ilə)
- `EmailRX` → WHATWG HTML spec regex-i

**Handler-də istifadəsi:**
```go
v := validator.New()
// ...v.Check(...) çağırışları...
if !v.Valid() {
    app.failedValidationResponse(w, r, v.Errors)
    return
}
```

**failedValidationResponse:**
```go
func (app *application) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
    app.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}
```
- 422 Unprocessable Entity + `{"error": {"field": "mesaj", ...}}` — bütün
  xətalar bir cavabda

**Movie qaydaları ( ValidateMovie ):**
- title: boş deyil, ≤ 500 bayt
- year: verilib, 1888 ≤ year ≤ cari il
- runtime: verilib, müsbət
- genres: 1-5 unikal janr

### 6. Yenidən istifadə — ValidateMovie() domain-də
**Nədir:** Yoxlamaları handler-dən ayırıb domain type yanına yerləşdirmək.

```go
// internal/data/movies.go
func ValidateMovie(v *validator.Validator, movie *Movie) {
    v.Check(movie.Title != "", "title", "must be provided")
    v.Check(len(movie.Title) <= 500, "title", "must not be more than 500 bytes long")
    v.Check(movie.Year != 0, "year", "must be provided")
    v.Check(movie.Year >= 1888, "year", "must be greater than 1888")
    v.Check(movie.Year <= int32(time.Now().Year()), "year", "must not be in the future")
    v.Check(movie.Runtime != 0, "runtime", "must be provided")
    v.Check(movie.Runtime > 0, "runtime", "must be a positive integer")
    v.Check(movie.Genres != nil, "genres", "must be provided")
    v.Check(len(movie.Genres) >= 1, "genres", "must contain at least 1 genre")
    v.Check(len(movie.Genres) <= 5, "genres", "must not contain more than 5 genres")
    v.Check(validator.Unique(movie.Genres), "genres", "must not contain duplicate values")
}
```

**Memarındandır qərarlar (kitabın izahı):**
1. Validator handler-də yaradılıb funksiyaya ötürülür (return yox) — sonra
   bir neçə validation helper-i zəncirləmək olar
2. Input struct → Movie struct-a kopyalanır; birbaşa Movie-yə decode
   ETMİRLİR — client `id`/`version` göndərib server tərəfli field-ləri
   doldara bilərdi (mass assignment qorunması)

## Bu fəsildə layihəyə nə əlavə olundu? (project continuity)

**Yeni fayllar:**
- `internal/validator/validator.go` — Validator, Check, PermittedValue,
  Matches, Unique, EmailRX
- helpers.go: `readJSON()` (tam versiya)

**Dəyişdirilən:**
- errors.go: `badRequestResponse()`, `failedValidationResponse()`
- movies.go: createMovieHandler (decode → copy → ValidateMovie)
- runtime.go: `UnmarshalJSON()` + `ErrInvalidRuntimeFormat`
- data/movies.go: `ValidateMovie()`

## Əsas terminlər

- Triage (xətaların sinifləndirilməsi) — xətanı tipinə görə ayırıb
  müvafiq cavaba çevirmək
- Sentinel Error (mayak xəta) — müqayisə üçün paket-səviyyəli xəta dəyəri
  (`errors.Is`)
- Pointer Receiver (işarə qəbuledicisi) — method receiver-ini dəyişmək
  üçün pointer
- Mass Assignment Protection — client-in server idarə etdiyi field-ləri
  təyin etməsinin qarşısı
- Unprocessable Entity (422) — semantik olarağ düzgün, amma business
  qaydalarına uyğmayan sorğu
- Denial-of-Service (xidmət inkarı hücumu) — resursları tükətmə hücumu;
  MaxBytesReader qoruyur

## Praktik nəticə

readJSON() helper-i "bir dəfə yaz, hər yerə köçür" tipindəndir: body limiti +
unknown field qadağası + tək-value tələbi + aydın xəta mesajları. Validator
paketi isə 422 + field-əsaslı xəta map-i ilə istənilən API üçün minimal,
genişlənə bilən validation qatıdır — generic funksiyalar sayəsində
(müəllif qeyd: productionda go-playground/validator kimi hazır həll də
müqayisə edin, amma kitabın əl yanaşması izah üçün idealdır).

## Mənbə
Pages: 74-104 (raw 074-104)
