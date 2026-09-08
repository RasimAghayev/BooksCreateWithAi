# Chapter 3 — Sending JSON Responses

## Bu fəsil nədən bəhs edir?

Handler-lərin plain-text əvəzinə JSON cavab qaytarması: fixed-format string,
`json.Marshal()`, struct encoding + struct tags, `json.Marshaler` interface,
envelope pattern və JSON error helper-ləri.

## Əsas fikirlər

### 1. JSON — sadəcə mətndir (Fixed-Format)
**Nədir:** JSON cavabı sadəcə string kimi yazmaq olar; yeganə vacib şey
`Content-Type: application/json` header-idir.

**Kitabdan kod nümunəsi:**
```go
func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
    js := `{"status": "available", "environment": %q, "version": %q}`
    js = fmt.Sprintf(js, app.config.env, version)
    w.Header().Set("Content-Type", "application/json")
    w.Write([]byte(js))
}
```

**Sub-kod izahı:**
- Raw string literal (backtick) → daxilən dırnaqları escape etməyə ehtiyac yoxdur
- `%q` verb-i → dəyərləri dırnağa alır (JSON string tələbi)
- `Content-Type` təyin edilməsə → Go `text/plain` göndərir

**Nəyə lazımdır:** Statik və ya kiçik dinamik JSON cavabları üçün sürətli yol.

**Qeyd (JSON charset):** `charset=utf-8` parametri lazımsızdır — RFC JSON-un
UTF-8 olmasını mütləq tələb edir; `application/json` media type-ın charset
parametri texniki olaraq mövcud deyil.

### 2. json.Marshal() ilə encoding
**Nədir:** Go obyektini (map/struct/slice) `[]byte` JSON-a çevirən standart
funksiya: `func Marshal(v any) ([]byte, error)`.

**Necə işləyir:**
```go
data := map[string]string{
    "status":      "available",
    "environment": app.config.env,
    "version":     version,
}
js, err := json.Marshal(data)
if err != nil {
    app.logger.Print(err)
    http.Error(w, "...", http.StatusInternalServerError)
    return
}
js = append(js, '\n')
```

**Go → JSON type mapping cədvəli:**
| Go type | JSON type |
|---|---|
| bool | boolean |
| string | string |
| int*, uint*, float*, rune | number |
| array, slice | array |
| struct, map | object (map key-ləri **əlifba sırası ilə** sıralanır) |
| nil pointer/interface/slice/map | null |
| chan, func, complex | DƏSTƏKLƏNMİR (`json.UnsupportedTypeError`) |
| time.Time | RFC3339 string (`"2020-11-08T06:27:59+01:00"`) |
| []byte | Base64 string (`"aGVsbG8="`) |

**json.Encoder alternativi:** `json.NewEncoder(w).Encode(data)` bir addımda
encode+yazır, amma encode xətası halında header dəyişmək mümkün olmur
(məs. Cache-Control şərti qoymaq çətinləşir). Benchmark: fərq ~µs səviyyəsində
(Encoder bir az az alloc edir) — HTTP cavabı üçün `json.Marshal()` üstünlük
təşkil edir.

### 3. writeJSON() helper-i
**Nədir:** Bütün JSON cavabları üçün vahid helper — encode, header, status,
body bir yerdə.

**Kitabdan kod nümunəsi:**
```go
type envelope map[string]any

func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
    js, err := json.MarshalIndent(data, "", "\t")
    if err != nil {
        return err
    }
    js = append(js, '\n')
    for key, value := range headers {
        w.Header()[key] = value
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    w.Write(js)
    return nil
}
```

**Sub-kod izahı:**
- `json.MarshalIndent(data, "", "\t")` → tab-indentli, oxunaqlı JSON
  (benchmark: ~65% yavaş, ~30% çox yaddaş — çoxu app üçün qəbuledilər)
- `headers` map-i (nil ola bilər — nil map range etmək təhlükəsizdir) →
  uğurlu cavaba əlavə header-lər (məs. `Location`) əlavə etməyə imkan verir
- `WriteHeader(status)` → body-dən ÖNCƏ status yazılmalıdır

**Təhlükəsizlik qaydası:** Header-lər yalnız encode uğurlu olandan SONRA
yazılır — xəta halında status dəyişmək mümkün qalır.

### 4. Struct encoding + struct tag-lər
**Nədir:** Custom struct JSON-a encode olunur; field adları/tag-lər nəzarəti
sizdədir.

**Kitabdan kod nümunəsi:**
```go
type Movie struct {
    ID        int64     `json:"id"`
    CreatedAt time.Time `json:"-"`            // heç vaxt görünmür
    Title     string    `json:"title"`
    Year      int32     `json:"year,omitempty"`
    Runtime   Runtime   `json:"runtime,omitempty"`
    Genres    []string  `json:"genres,omitempty"`
    Version   int32     `json:"version"`
}
```

**Tag direktivləri:**
- `json:"id"` → key adını dəyişir (snake_case)
- `json:"-"` → field JSON-da heç vaxt görünmür (parol hash-i kimi)
- `json:",omitempty"` → dəyər "boş"dursa gizlədilir (false, 0, "", boş
  slice/map, nil pointer/interface)
- `json:"runtime,omitempty,string"` → rəqəmi JSON **string** kimi məcburi
  edir (yalnız int*/uint*/float*/bool işləyir)

**Vacib:** Field-lər export olunmalıdır (böyük hərf) — yoxsa encoding/json
onları görmür. `json:"-"` unexported-dan yaxşıdır: niyyət açıqdır.

### 5. Envelope pattern
**Nədir:** Cavab datasını `{"movie": {...}}` kimi parent obyektə sarmaq.

```go
err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
```

**Üstünlükləri:**
1. Self-documenting (özünü izah edən) cavablar
2. Client tərəfində yanlış işləmə riski azalır
3. Köhnə brauzerlərdə top-level JSON array qaytarmağın yaratdığı təhlükəsizlik
   boşluğu aradan qalxır

**Qeyd:** JSON:API / jsend kimi populyar formatlar var, amma vacib olan
**tutarlı strukturdur**.

### 6. json.Marshaler interface — dərin customization
**Nədir:** Type `MarshalJSON() ([]byte, error)` method-u implement edirsə, Go
encoding zamanı məhz onu çağırır.

```go
type Runtime int32

func (r Runtime) MarshalJSON() ([]byte, error) {
    jsonValue := fmt.Sprintf("%d mins", r)
    quotedJSONValue := strconv.Quote(jsonValue)
    return []byte(quotedJSONValue), nil
}
```

**Sub-kod izahı:**
- `Runtime int32` → underlying type int32 olan custom type
- `strconv.Quote()` → JSON string tələbi üçün dırnağa salır (getməsə:
  `invalid character 'm' after top-level value` runtime xətası)
- **Value receiver** seçilib: value method həm value, həm pointer üzərindən
  çağrıla bilər ("value methods can be invoked on pointers and values, but
  pointer methods can only be invoked on pointers" — Effective Go)
- `omitempty` custom type üzərində də işləyir: underlying 0 → field
  gizlədilir, MarshalJSON çağrılmır

**Alternativ yanaşmalar (kitabın müqayisəsi):**
1. **Bütün struct-a MarshalJSON** — anonymous struct ilə bütün sahələri
   yenidən yazmaq: verbose, amma tam nizam nəzarəti
2. **Alias embedding** — `type MovieAlias Movie` + anonymous struct-a embed:
   qısa, amma "trick" hesab olunur, method varisliyi olmamasına əsaslanır
   (sonsuz loopun qarşısını almaq üçün alias vacibdir!), field sırası
   itirilir (runtime sonda görünür)

### 7. JSON error helper-ləri
**Nədir:** Plain-text `http.Error()` əvəzinə strukturlaşdırılmış JSON xətalar.

**Kitabdan kod nümunəsi:**
```go
// cmd/api/errors.go
func (app *application) logError(r *http.Request, err error) {
    app.logger.Print(err)
}

func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
    env := envelope{"error": message}
    err := app.writeJSON(w, status, env, nil)
    if err != nil {
        app.logError(r, err)
        w.WriteHeader(500)
    }
}

func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
    app.logError(r, err)
    message := "the server encountered a problem and could not process your request"
    app.errorResponse(w, r, http.StatusInternalServerError, message)
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
    app.errorResponse(w, r, http.StatusNotFound, "the requested resource could not be found")
}

func (app *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
    message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
    app.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}
```

**Sub-kod izahı:**
- `logError()` → ayrı helper; sonra structured logging-ə yüksəldiləcək
- `errorResponse()` → mərkəzi `{"error": ...}` formatı; message `any` —
  string və ya map (validasiya xətaları üçün) qəbul edir
- `serverErrorResponse()` → detalı loglayır, clientə GENERIC mesaj (internal
  xəta detalı istemciyə İTİLMİZ)
- writeJSON özü xətalarsa → fallback: status 500 + boş body

### 8. Router xətalarının JSON-a keçirilməsi
**Nədir:** httprouter 404/405-ləri default plain-text göndərir — custom
handler-lər təyin olunur.

```go
func (app *application) routes() *httprouter.Router {
    router := httprouter.New()
    router.NotFound = http.HandlerFunc(app.notFoundResponse)
    router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)
    // ...
}
```
- `router.NotFound` / `router.MethodNotAllowed` → `http.Handler` qəbul edir;
  `http.HandlerFunc(app.notFoundResponse)` adapteri ilə method → Handler
- httprouter `Allow` header-i custom handler ilə də avtomatik qoyur

**Limit:** Go `http.Server` bəzi hallarda (malformed Host header, köhnə
protokol, HTTPS-ə HTTP sorğu və s.) plain-text cavabı özü göndərir — bunlar
stdlib-ə hard-coded-dir, customize edilə bilməz (yalnız吓 yaman clients
görür, problem deyil).

## Bu fəsildə layihəyə nə əlavə olundu? (project continuity)

**Yeni fayllar:**
- `internal/data/movies.go` — Movie struct (struct tags ilə)
- `internal/data/runtime.go` — Runtime custom type + MarshalJSON
- `cmd/api/errors.go` — logError, errorResponse, serverErrorResponse,
  notFoundResponse, methodNotAllowedResponse
- helpers.go: `writeJSON()` (envelope + MarshalIndent) və `envelope` type

**Dəyişdirilən:** healthcheck (envelope + system_info), showMovie (Movie
struct + envelope), routes.go (router.NotFound/MethodNotAllowed).

## Əsas terminlər

- Envelope (zərf pattern-i) — cavabın parent obyektə sarılması
- Struct Tag (struktur etiketi) — `json:"..."` metadata direktivi
- Zero Value (sıfır dəyər) — tipin default dəyəri (0, "", false, nil)
- Interface Satisfaction (interfeysin təminatı) — type method set-i ilə
  kontraktı qarşılayır
- Value Receiver (dəyər qəbuledicisi) — method nüsxə üzərində işləyir
- Base64 Encoding (ikilik-mətn kodlaşdırması) — []byte → string

## Praktik nəticə

`writeJSON()` + `envelope` + error helper trio-su istənilən Go API-nin
skeletidir: bütün cavablar və xətalar vahid JSON formatında olur, header
inzası təhlükəsiz şəkildə yerləşir, router xətaları belə eyni formaya düşür.
Custom type + `MarshalJSON()` isə JSON təqdimatını domain modelindən ayırmağın
idiomatik yoludur.

## Mənbə
Pages: 31-73 (raw 031-073)
