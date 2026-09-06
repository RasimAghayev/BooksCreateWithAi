# Chapter 2 — Developing Echo Projects

## Bu chapter nədən bəhs edir?

Bu chapter Echo layihəsinin **strukturunu və əsas inkişaf primitivlərini** təqdim edir: sənaye-tipli qovluq strukturu (bindings/cmd/handlers/middlewares/models/renderings/static), `dep` dependency management, routing + handler imzaları, Login handler üzərində tam auth axını (Bind → Validate → bcrypt → JWT), middleware anlayışı və rendering imkanları. Kitabın boyu davam edən **health-check + login/logout API** layihəsi bu chapter-də başlayır.

---

## Əsas fikirlər

### 1. Layihə strukturu — Echo-nun opinionated olmayan təşkilatı

Django (MVC) kimi framework-lər sərt struktur tələb edir; Echo **flexibility** verir. Amma güclü struktur: siklik import-ları azaldır, təmiz dependency tree yaradır, komanda işi üçün intuitivdir.

```
chapter2/
    bindings/       # request input tipləri (form/query/JSON)
    cmd/
        service/    # entrypoint (main)
    handlers/       # handler kodu + business logic
    middlewares/    # birdən çox handler-ə aid logic
    models/         # DB-nə persist olunan tiplər
    renderings/     # response serializasiya tipləri
    static/         # statik assetlər
    vendor/         # dep idarə etdiyi dependency-lər
```

**Hər qovluq nəyə görə:**

| Qovluq | Məqsəd | Müəllif şərhi |
|--------|--------|---------------|
| `bindings` | Protokol-spesifik input strukturları | Generic model istifadəsi API versiyalama problem yaradır |
| `cmd/service` | Entrypoint | `cmd/scheduler` kimi yeni xidmətlər eyni repodan ayrıla bilər — microservice hazırlığı |
| `handlers` | Echo handler + business logic | Fayl adları funksiya ilə eyniadlı olmalı (health_check.go → HealthCheck) |
| `middlewares` | Çox-handler üçün ortaq logic | Məs., "Verify Token" |
| `models` | DB persist tiplləri | bindings/renderings-ə ümumi base ola bilər |
| `renderings` | ResponseWriter-a serializə olunan tiplər | |
| `static` | Statik məzmun | HTML template / JS / image |

**Fayl adlandırmada simmetriya:** `handlers/health_check.go` + `renderings/health_check.go` — strukturu intuitiv saxlayır.

**Entry point** (`cmd/service/main.go`):

```go
package main

import (
    "github.com/PacktPublishing/Echo-Essentials/chapter2/handlers"
    "github.com/labstack/echo"
)

func main() {
    // create a new echo instance
    e := echo.New()
    // Route / to handler function
    e.GET("/health-check", handlers.HealthCheck)
    // Authentication routes
    e.POST("/login", handlers.Login)
    e.POST("/logout", handlers.Logout)
    // start the server, and log if it fails
    e.Logger.Fatal(e.Start(":8080"))
}
```

**HealthCheck handler + rendering:**

```go
// handlers/health_check.go
func HealthCheck(c echo.Context) error {
    resp := renderings.HealthCheckResponse{
        Message: "Everything is good!",
    }
    return c.JSON(http.StatusOK, resp)
}

// renderings/health_check.go
type HealthCheckResponse struct {
    Message string `json:"message"`
}
```

---

### 2. Dependency management — dep tool

**Problem:** `go get` versiyaları pinləmir — dependency internetdən silinərsə build qırılır (NPM left-pad hadisəsi kimi).

**Həll — dep:**

```bash
# tool-u yüklə
go get -u github.com/golang/dep/cmd/dep

# layihədə vendor qovluğunu initializə et — import-ları oxuyub kodu vendor/ içəri köçürür
dep init
```

**Necə işləyir:** dep source import-ları oxuyur, lazım olan paketləri `$GOPATH`-dan (və ya `go get`/`git` ilə) alıb layihənin `vendor/` qovluğuna kopyalayır və versiyaları pinləyir.

**Tövsiyə:** `vendor/` qovluğunu birbaşa repoya commit edin — internetdən asılılığı aradan qaldırır. Çatışmazlığı: böyük dependency dəyişikliklərində PR-lər böyüyür (amma bu, workflow problemi əlamətidir).

---

### 3. Routing və handler-lər

**Routing nədir?** Target path → handler funksiyası xəritələnməsi.

**Sorğu növləri:**
1. **Static mapping:** `/reminder` → `CreateReminder` — hər dəfə eyni funksiya
2. **Regex yanaşması:** `/reminder/\d+` — `/reminder/123`-ü yakalayır, amma HƏR request-də bütün route-ları sırayla yoxlamaq lazım gəlir (O(n))
3. **Echo yanaşması — Radix tree:** Prefiks-ağacı ilə effektiv axtarış (Ch3-də dərinləşəcək)

**Handler tiplərinin müqayisəsi:**

```go
// Standart kitabxana:
type Handler interface {
    ServeHTTP(w ResponseWriter, r *Request)
}
type HandlerFunc func(ResponseWriter, *Request)

// Echo:
type HandlerFunc func(Context) error
```

**Login handler — tam auth axını** (`handlers/login.go`):

```go
// Login - Login Handler will take a username and password from the request
// hash the password, verify it matches in the database and respond with a token
func Login(c echo.Context) error {
    resp := renderings.LoginResponse{}
    lr := new(bindings.LoginRequest)

    // 1. BIND — request-in avtomatik deserializasiyası
    if err := c.Bind(lr); err != nil {
        resp.Success = false
        resp.Message = "Unable to bind request for login"
        return c.JSON(http.StatusBadRequest, resp)
    }

    // 2. VALIDATE — input qaydaları
    if err := lr.Validate(c); err != nil {
        resp.Success = false
        resp.Message = err.Error()
        return c.JSON(http.StatusBadRequest, resp)
    }

    // 3. USER LOOKUP — DB context-dən
    db := c.Get(models.DBContextKey).(*sql.DB)
    user, err := models.GetUserByUsername(db, lr.Username)
    if err != nil {
        resp.Success = false
        resp.Message = "Username or Password incorrect"
        return c.JSON(http.StatusUnauthorized, resp)
    }

    // 4. BCRYPT — parol yoxlaması
    if err := bcrypt.CompareHashAndPassword(
        user.PasswordHash, []byte(lr.Password)); err != nil {
        resp.Success = false
        resp.Message = "Username or Password incorrect"
        return c.JSON(http.StatusUnauthorized, resp)
    }

    // 5. JWT TOKEN — uğurlu login
    signingKey := c.Get(models.SigningContextKey).([]byte)
    claims := &jwt.StandardClaims{
        ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
        Issuer:    "service",
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    ss, err := token.SignedString(signingKey)
    if err != nil {
        resp.Success = false
        resp.Message = "Server Error"
        return c.JSON(http.StatusInternalServerError, resp)
    }
    resp.Token = ss
    return c.JSON(http.StatusOK, resp)
}
```

**Köməkçi axın:** Bind → Validate → DB lookup → bcrypt müqayisə → JWT generasiyası → JSON response. Handler-lər business logic-in olduğu yerdir.

---

### 4. Middleware — ilk baxış

**Nədir?** Handler-lərin **wrapper-i** — eyni logikanı çoxlu resursa tətbiq etmək üçün. Echo middleware tərifi: növbəti çağrılacaq funksiyanı parametr kimi qəbul edib Echo handler qaytaran funksiya → **zəncir qurulur**.

```go
// Signing Key for our auth middleware
var signingKey = []byte("superdupersecret!")

e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        c.Set(models.SigningContextKey, signingKey)
        return next(c)
    }
})

// Qorumalı route qrupu:
reminderGroup := e.Group("/reminder")
reminderGroup.Use(middleware.JWT(signingKey))
reminderGroup.POST("", handlers.CreateReminder)
```

Nümunə: login və static resurslar autentifikasiya istəmir; reminder yaratmaq isə istəyir — middleware bu fərqi idarə edir. Handler kodu sadələşir, təkrar istifadə artır (Ch4-də dərinləşir).

---

### 5. Rendering — Echo-nun cavab helper-ləri

`echo.Context` metodları serializasiyanı öz üzərinə götürür:

| Metod | Təyinat |
|-------|---------|
| `HTML` / `HTMLBlob` | HTML render (hazır blob daxil) |
| `JSON` / `JSONBlob` / `JSONPretty` / `JSONP` | JSON render — struct tag-lərə əsaslanır |
| `XML` / `XMLBlob` / `XMLPretty` | XML render |
| `File` / `Attachment` | Fayl / yükləmə cavabı |
| `Blob` | İstənilən content-type ilə bayt blob |
| `String` | Düz mətn |
| `NoContent` | Bodysiz cavab |
| `Redirect` | Yönləndirmə |

API üçün adətən JSON/XML; server-side rendered tətbiqlər üçün HTML + template.

---

## Əsas terminlər

| Termin | İzah |
|--------|------|
| Layihə strukturu | bindings/cmd/handlers/middlewares/models/renderings/static bölgüsü |
| `bindings` | Protokol-spesifik request input strukturları (JSON/form/query) |
| `renderings` | ResponseWriter-a serializə olunan cavab strukturları |
| `cmd/service` | Main entrypoint qovluğu — çox-xidmətli repo hazırlığı |
| `vendor/` | Dependency kodunun layihə daxilində saxlanması |
| `dep init` | Import-ları oxuyub vendor qovluğunu yaradan əmr |
| Static route mapping | Düz path → funksiya xəritəsi |
| Radix tree | Prefiks-ağacı — Echo-nun routing data strukturu |
| `type HandlerFunc func(Context) error` | Echo handler imzası (std `ServeHTTP(w, r)`-dan fərqli) |
| `c.Bind` | Request-in avtomatik struktur deserializasiyası |
| bcrypt | Parol hash müqayisə alqoritmi (`CompareHashAndPassword`) |
| JWT (`StandardClaims`) | Token-based auth — HS256 signing, 72 saat expiry |
| `c.Get`/`c.Set` | Context üzərindən state keçirmə (DB, signing key) |
| `e.Group` + `.Use` | Route qrupu + qrup-spesifik middleware |
| `c.JSON(code, i)` | Struct-ı JSON serializasiya edib yazan helper |

---

## Praktik nəticə

1. **Struktur işin yarsıdır:** bindings (input) / renderings (output) / models (DB) ayrımı API versiyalamasını asanlaşdırır; `cmd/` çox-xidmətli bölgüyə (microservice) imkan verir.
2. **Dependency-ləri pinləyin:** `dep init` + `vendor/`-ı commit-ləyin — left-pad tipli build qırılmalarından sığorta.
3. **Handler = business logic mərkəzi:** Bind → Validate → DB → bcrypt → JWT → JSON — bu ardıcıllıq kitab boyu bütün handler-lərin şablonudur.
4. **Auth-i middleware-ə çıxarın:** `middleware.JWT(signingKey)` qrup səviyyəsində — handler-lər autentifikasiyadan xəbərsiz qalır.
5. **Render helper-lərindən istifadə edin:** `c.JSON/String/XML` — manual `json.Marshal` + `Write` zəncirini əvəz edir.
6. **Fayl adları strukturun güzgüsü olsun:** `handlers/login.go` ↔ `bindings/login.go` ↔ `renderings/login_response.go` simmetriyası.

---

## Mənbə

- Kitab: *Echo Quick Start Guide* — J. Ben Huson, Packt Publishing, 2018 (ISBN 9781789340849)
- Chapter 2: "Developing Echo Projects", book səh. 23–48
- PDF səhifələri: 48–69
- Kod: https://github.com/PacktPublishing/Echo-Essentials/tree/master/chapter2
- Video: https://goo.gl/NudCD9
