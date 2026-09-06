# Echo Quick Start Guide — Cheatsheet (Azərbaycanca)

> **Kitab:** Echo Quick Start Guide — J. Ben Huson, Packt 2018 · **Echo v3.3.5** · Go və Echo framework-unun əsas komandalarının və pattern-lərinin sürətli istinadı.

---

## 1. Mühit və layihə qurulması

```bash
# Go quraşdır (kitabda 1.9.4; hazırda daha yeni)
brew install go                      # macOS
# Windows: MSI installer (c:\Go\)

# Echo (versiya pinləmə ilə)
go get github.com/labstack/echo
cd $GOPATH/src/github.com/labstack/echo && git checkout tags/3.3.5

# dep — dependency management
go get -u github.com/golang/dep/cmd/dep
dep init                             # vendor/ yaradır, import-ları köçürür

# Layihə kodu
git clone https://github.com/PacktPublishing/Echo-Essentials
```

**Layihə strukturu:** `bindings/` (input) · `cmd/service/` (main) · `handlers/` · `middlewares/` · `models/` (DB) · `renderings/` (output) · `static/` · `vendor/`.

---

## 2. Minimal Echo tətbiqi

```go
e := echo.New()
e.GET("/", handler)
e.Logger.Fatal(e.Start(":8080"))     // :8080-da dinlə

func handler(c echo.Context) error {
    return c.String(http.StatusOK, "Hello World")
}
```

---

## 3. Routing

```go
e.GET/POST/PUT/PATCH/DELETE/HEAD/OPTIONS(path, handler, mws...)  // metod helper-ləri
e.Add(method, path, handler, mws...)                              // ümumi
e.Any(path, handler)                                              // bütün metodlar
e.Static("/static", "static")     // qovluqdan statik fayl
e.File("/", "static/index.html")  // tək fayl route-a
e.GET("/users/:id", h)            // URL parametri → c.Param("id")
e.GET("/v:version/x/:id", h)      // versiya parametr kimi
e.GET("/assets/*", h)             // wildcard (statik suffix)
```

**Group-lar:**

```go
v1 := e.Group("/v1")
v1.POST("/login", handlers.Login)
v1Rem := v1.Group("/reminder", middleware.JWT(signingKey))  // iç-içə + middleware
v1Rem.POST("", handlers.CreateReminder)
v1Rem.GET("/:id", handlers.GetReminder)
```

- Radix tree axtarışı: ortaq prefikslər performance artırır.
- `/login/` (trailing slash) 404 verir → `TrailingSlash` middleware istifadə et.
- Benchmark: Echo 38,662 ns/op, **0 allocs** (context pool reuse).

---

## 4. Middleware

```go
type MiddlewareFunc func(HandlerFunc) HandlerFunc

e.Pre(mw)      // routing-DƏN ƏVVƏL (RequestID, slash düzəlişi, early recover)
e.Use(mw)      // routing-dan SONRA (global)
grp.Use(mw)    // qrup-spesifik
e.GET("/", h, mw1, mw2)  // route-spesifik (variadic)
// Sıra = icra sırası: m1(m2(h)) — Logger-in panic-i üçün Recover xaricdə olsun
```

**Custom middleware şablonu:**

```go
func RequestIDMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        requestID := uuid.NewV4()
        c.Logger().Infof("RequestID: %s", requestID)
        c.Set(RequestIDContextKey, requestID)
        return next(c)      // çağırmayanda zəncir kəsilir
    }
}
e.Pre(middlewares.RequestIDMiddleware)
```

**Handler-də oxu:**

```go
if reqID, ok := c.Get(middlewares.RequestIDContextKey).(uuid.UUID); ok {
    c.Logger().Debugf("RequestID: %s", reqID.String())
}
```

**Contributed:** `middleware.Logger()`, `middleware.Recover()`, `middleware.JWT(signingKey)`, `TrailingSlash`...

---

## 5. Context (`echo.Context`)

```go
c.Param("id")              // URL parametri
c.QueryParam("q")          // ?q=... query string
c.FormValue("field")       // form dəyəri
c.FormFile("upload")       // multipart fayl
c.Cookie("name")           // request cookie
c.SetCookie(&http.Cookie{...})  // Set-Cookie (response)
c.Get(key) / c.Set(key, v)      // middleware↔handler state
c.Logger()                 // Echo logger
c.Echo()                   // Echo instance → .Reverse("route-name")
c.Handler()                // router-in tapdığı handler
```

Handler imzası: `func(echo.Context) error`. Context `sync.Pool`-dan reuse olunur → 0 allocation.

---

## 6. Request Binding + Validation

```go
// Binding strukturu (bindings/ paketində)
type LoginRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

// Handler-də:
lr := new(bindings.LoginRequest)
if err := c.Bind(lr); err != nil { /* 400 */ }      // Content-Type-a görə: JSON/XML/form
if err := c.Validate(lr); err != nil { /* 400 */ }  // Validatable interface
```

**Validatable + Validator qurulması:**

```go
// bindings/common.go
type Validatable interface{ Validate() error }
type Validator struct{}
func (v *Validator) Validate(i interface{}) error {
    if val, ok := i.(Validatable); ok { return val.Validate() }
    return ErrNotValidatable
}

// bindings/login.go — qaydalar strukturun yanında
func (lr *LoginRequest) Validate() error {
    errs := new(RequestErrors)
    if lr.Username == "" { errs.Append(ErrUsernameEmpty) }
    if lr.Password == "" { errs.Append(ErrPasswordEmpty) }
    if errs.Len() == 0 { return nil }
    return errs
}

// main.go — bir dəfə:
e.Validator = new(bindings.Validator)
```

---

## 7. Response Rendering

```go
c.JSON(code, i)                  // ən çox istifadə (struct tag-lərlə)
c.JSONPretty(code, i, "  ")  c.JSONBlob(code, b)  c.JSONP(code, cb, i)
c.String(code, s)                c.HTML(code, html)     c.HTMLBlob(code, b)
c.XML(code, i)                   c.XMLPretty / XMLBlob
c.Blob(code, contentType, b)     c.Stream(code, ct, r)   c.File(path)
c.Attachment(file, name)         c.Inline(file, name)
c.NoContent(code)                c.Redirect(code, url)
```

**Aşağı səviyyə — `echo.Response`:**

```go
c.Response().Header().Add("X-My-Header", "value")
c.Response().Status / .Size / .Committed        // Status/Size/Committed sahələri
c.Response().Before(fn) / .After(fn)            // yazıdan əvvəl/sonra hook-lar
// WriteHeader: Committed qoruması (double-write → warn)
// Write: !Committed → implicit WriteHeader(200)
```

---

## 8. Logging və Error Handling

```go
// Səviyyə tənzimi
e.Logger.SetLevel(log.DEBUG)   // production: log.INFO

// Context-dən:
c.Logger().Infof(...)  .Debugf(...)  .Errorf(...)  .Warnf(...)
// JSON variant: .Infoj(...) .Debugj(...) ...

// Request log middleware:
e.Use(middleware.Logger())
// → {"time","id","remote_ip","host","method","uri","status","latency","bytes_in/out"}

// Handler-də sadəcə error qaytar:
return errors.New("failure!")               // → 500 {"message":"Internal Server Error"}
return echo.NewHTTPError(404, "not found")  // → 404 status+message

// Custom error handler:
func myHTTPErrorHandler(err error, c echo.Context) {
    code := http.StatusInternalServerError
    if httpErr, ok := err.(*echo.HTTPError); ok { code = httpErr.Code }
    c.String(code, ...)
}
e.HTTPErrorHandler = myHTTPErrorHandler

// Panic qoruması — MÜTLƏQ:
e.Use(middleware.Recover())   // recover() ilə panic-i tutub 500 cavab verir
```

**Log səviyyələri:** DEBUG (troubleshoot) → INFO (production default, əməliyyat məlumatı) → WARN (bərpa olunmuş problem) → ERROR (təhqiqat lazım) → FATAL/PANIC (yalnız başlanğıc xətaları). Secret-ləri log-lamayın!

---

## 9. Testing

```go
// handlers/health_check_test.go
func TestHealthCheck(t *testing.T) {
    e := echo.New()
    e.Pre(middlewares.RequestIDMiddleware)
    e.GET("/health-check", HealthCheck)
    w := httptest.NewRecorder()                       // fake ResponseWriter
    r, _ := http.NewRequest("GET", "/health-check", nil)
    e.ServeHTTP(w, r)                                 // server-siz icra!
    resp := w.Result()
    if resp.StatusCode != http.StatusOK { t.Error("unexpected", resp.Status) }
    hc := new(renderings.HealthCheckResponse)
    json.NewDecoder(resp.Body).Decode(hc)             // JSON yoxlaması
    if hc.Message != "Everything is good!" { t.Error("invalid message") }
}

// Benchmark
func BenchmarkHealthCheck(b *testing.B) {
    // setup...
    for i := 0; i < b.N; i++ { e.ServeHTTP(w, r) }
}
```

**Mock pattern (sql.DB):**

```go
type MockableDB interface { /* sql.DB-nin bütün imzaları */ }
type MockDB struct {
    mockQuery func(query string, args ...interface{}) (*sql.Rows, error)
    // ... hər metod üçün funksiya sahəsi
}
func (db *MockDB) Query(q string, a ...interface{}) (*sql.Rows, error) {
    if db.mockQuery != nil { return db.mockQuery(q, a...) }
    return nil, nil
}
// Funksiyalar MockableDB qəbul etməlidir → test-də MockDB, prod-da sql.DB
```

**Integration coverage (instrumented server):**

```go
// cmd/service/main_test.go
func TestRunMain(t *testing.T) {
    TestRun = true; go main(); <-StopTestServer; TestRun = false
}
// main.go: if TestRun { e.POST("/stop-test-server", ...) }

go test -coverprofile=cov.txt -coverpkg ./handlers -run TestRunMain ./cmd/service/
go tool cover -func=cov.txt        // funksiya-səviyyəli hesabat
```

---

## 10. Templates və Statik Məzmun

```go
// Statik:
e.Static("/static", "static")
e.File("/", "static/index.html")

// Renderer (bir dəfə, entry point-də):
type CustomTemplate struct{ *template.Template }
func (ct *CustomTemplate) Render(w io.Writer, name string, data interface{},
    ctx echo.Context) error {
    return ct.ExecuteTemplate(w, name, data)
}
t, _ := template.New("reminders").Parse(handlers.RemindersTmpl)
e.Renderer = &handlers.CustomTemplate{t}

// Handler-də:
return c.Render(http.StatusOK, "reminders", tmplData)
```

**Template sintaksisi (html/template):**

```
{{.Title}} {{.Field}}                     çap
{{if p}}...{{else}}...{{end}}             şərt (boş: false/0/nil/""/boş kolleksiya)
{{range .Items}}...{{else}}No rows{{end}}  iterasiya (array/slice/map/chan)
{{with p}}...{{end}}  {{template "name"}}  {{block "n" p}}...{{end}}
{{/* şərh */}}
```

**Reverse URL (şablondan Echo çağırma):**

```go
e.POST("/login", handlers.Login).Name = "login"    // route adı
type TmplData struct {
    Reminders []Reminder
    Title     string
    rev       func(name string, params ...interface{}) string
}
func (td TmplData) Reverse(name string, p ...interface{}) string { return td.rev(name, p...) }
data := TmplData{reminders, "Title", c.Echo().Reverse}
// Şablonda: <a href={{ .Reverse "login" }}>Login</a>
```

---

## Sürətli yaddaş cədvəli

| Ehtiyac | Əmr/Metod |
|---------|-----------|
| Server başlat | `e.Logger.Fatal(e.Start(":8080"))` |
| Route + middleware | `e.POST("/x", h, mw1, mw2)` |
| Pre-routing mw | `e.Pre(mw)` |
| Body → struct | `c.Bind(&x)` + `c.Validate(&x)` |
| Struct → JSON cavab | `c.JSON(200, resp)` |
| URL param | `c.Param("id")` |
| State keçir | `c.Set(k, v)` / `c.Get(k)` |
| 404/500 error | `echo.NewHTTPError(code, msg)` |
| Panic qoru | `e.Use(middleware.Recover())` |
| Request log | `e.Use(middleware.Logger())` |
| Test handler | `e.ServeHTTP(httptest.NewRecorder(), req)` |
| Statik fayl | `e.Static("/s", "static")` |
| Template render | `e.Renderer = ...` + `c.Render(200, "name", data)` |
