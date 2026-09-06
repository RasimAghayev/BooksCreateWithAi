# Chapter 4 — Implementing Middleware

## Bu chapter nədən bəhs edir?

Bu chapter **middleware konsepsiyasını** tam açır: handler-lərdə təkrarlanan kodun middleware-ə çıxarılması, Echo-nun 3 middleware qoşulma mexanizmi (`Use`, `Pre`, route-spesifik), zəncirin daxili işləməsi (nested funksiya çağrıları), 4 məntiqi sərhəd (global pre-routing, global post-routing, group, route), custom middleware yazılması (`RequestIDMiddleware`) və JWT auth middleware-in işləməsi.

---

## Əsas fikirlər

### 1. Middleware nəyə lazımdır?

**Problem:** Handler-lərdə təkrarlanan ortaq işlər — request body parse/binding, request logging, auth/session yoxlaması, metadata yaratma, response render, panic/error idarəsi.

**Alternativ həll — helper funksiyaları** (qeyri-optimal):
- Hər handler-də yenidən çağırmaq lazım gəlir — handler helper çağrıları ilə dolu olur.
- Helper imzası dəyişəndə İSTİFADƏ EDƏN BÜTÜN handler-lərə toxunmalı.

**Middleware həlli:** Ortaq funksionallığı handler-dən tam kənarlaşdır — handler-lər bu komplekslikdən xəbərsiz qalır. Echo `Context.Bind` buna yaxşı nümunədir (framework parsing-i öz üzərinə götürür).

**RequestID nümunəsi — əvvəl (duplikasiya):**

```go
// HealthCheck handler
func HealthCheck(c echo.Context) error {
    requestID := uuid.NewV4()
    c.Logger().Infof("RequestID: %s", requestID)
    resp := renderings.HealthCheckResponse{Message: "Everything is good!"}
    return c.JSON(http.StatusOK, resp)
}

// Login handler — EYNİ 2 sətir təkrarlanır
func Login(c echo.Context) error {
    requestID := uuid.NewV4()
    c.Logger().Infof("RequestID: %s", requestID)
    // ...
}
```

**Sonra (middleware-də konsolidasiya):**

```go
const (
    RequestIDContextKey = "request_id_context_key"
)

func RequestIDMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
    return echo.HandlerFunc(func(c echo.Context) error {
        requestID := uuid.NewV4()
        c.Logger().Infof("RequestID: %s", requestID)
        c.Set(RequestIDContextKey, requestID)
        return next(c)
    })
}
```

Handler artıq təmizdir — RequestID-dən xəbərsizdir:

```go
func HealthCheck(c echo.Context) error {
    resp := renderings.HealthCheckResponse{
        Message: "Everything is good!",
    }
    return c.JSON(http.StatusOK, resp)
}
```

**Middleware-in gücü:** Request handler-a çatmamış intercept edə bilər — auth/authorization üçün ideal (yetkisiz istifadəçiyə handler-ə çatmamıdan cavab qaytarmaq).

---

### 2. Middleware qoşulmasının 3 yolu + sıralanma

#### Yol 1 — `Use` (ən çox istifadə olunan, post-routing)

```go
e.Use(middleware.Logger())   // logger middleware "wrap" edir
e.Use(middleware.Recover())  // recovery-i
```

Echo instance və Group-larda mövcuddur. **Sıra vacibdir** — Use çağrılarının sırası zəncirin icra sırasını müəyyən edir. Misal: Logger Recover-dən əvvəl Use olunubsa, Logger-in ÖZÜ panic etsə, Recover onu tutan bilməz (Recover Logger-in İÇİNDƏDİR, xaricində yox).

#### Yol 2 — `Pre` (routing-dən əvvəl)

```go
e.Pre(middlewares.RequestIDMiddleware)
```

Global, route-spesifik olmayan middleware-lər üçün. Use case-lər:
- Routing-dən əvvəl request-i manipulyasiya etmək (trailing slash silmə, redirect)
- Panic recover — stack trace-in response-a düşməməsi üçün
- RequestID kimi "hər request mütləq alsın" logic-i

Sıralama eyni qaydadır: nə qədər erkən Pre, o qədər erkən başlayır.

#### Yol 3 — Route qeydiyyatında variadic parametr

```go
e.GET("/", HandlerFunction, Middleware1, Middleware2, Middleware3)
// RouteHandler = Middleware1(Middleware2(Middleware3(HandlerFunction)))
```

Yalnız həmin konkret route-a tətbiq olunur.

#### 4 məntiqi sərhəd, 2 daxili yer

API dizayn baxımından 4 sərhəd var:
1. Global middleware (routing-dən əvvəl) — `Pre`
2. Global middleware (routing-dən sonra) — `Use`
3. Qrup route middleware — `Group.Use`
4. Route-spesifik middleware — variadic parametr

Echo implementasiyasında isə yalnız **2 daxili yer** var: routing-dən əvvəl və routing-dən sonra. Qrup/route middleware-ləri route insert zamanı handler-in wrap-inə çevrilir.

---

### 3. Zəncirin daxili işləməsi

**Axın:**
1. Klient request-i → server qəbul edir
2. **Pre-middleware zənciri** sırayla icra olunur — hər middleware `next(c)` çağırmaqla növbətinə keçir
3. Pre zənciri bitəndə **router** handler-i axtarır
4. Router tapılan route-un wrap-olunmuş handler zəncirini (route middleware-ləri + handler) çağırır
5. İstənilən nöqtədə middleware `next` ÇAĞIRMADAN return etsə — zəncir geri sarınır, cavab klientə gedir

**Sır:** Middleware verya sadədir — "heç bir sehr yoxdur": middleware = sıralanmış **nested funksiya çağırıları dəsti**, handler ən son nested çağırışdır.

```go
// Konseptual olaraq:
finalHandler = Middleware1(Middleware2(HandlerFunction))
// Middleware1-ə daxil olanda "next" = Middleware2(Handler)-in nəticəsidir
```

---

### 4. Custom middleware yazmaq

**Go-da funksiyalar birinci sinif vətəndaşdır** — dəyişən kimi istifadə oluna bilər, parametr kimi ötürülə bilər. Echo middleware tipi:

```go
type MiddlewareFunc func(HandlerFunc) HandlerFunc
```

**"Növbətini qəbul edən, handler qaytaran funksiya."**

Tam custom nümunə (`middlewares/request_id.go`):

```go
package middlewares

import (
    "github.com/labstack/echo"
    uuid "github.com/satori/go.uuid"
)

const (
    requestIDContextKey = "request_id_context_key"
)

func RequestIDMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
    return echo.HandlerFunc(func(c echo.Context) error {
        requestID := uuid.NewV4()               // 1. UUID yarat
        c.Logger().Infof("RequestID: %s", requestID) // 2. Loqla
        c.Set(RequestIDContextKey, requestID)    // 3. Context-ə qoy
        return next(c)                           // 4. Zənciri davam etdir
    })
}
```

**Qoşulma — Pre kimi** (hər request routing-dən əvvəl ID almalıdır):

```go
e.Pre(middlewares.RequestIDMiddleware)
```

**Handler-də oxunması:**

```go
func HealthCheck(c echo.Context) error {
    if requestID, ok := c.Get(middlewares.RequestIDContextKey).(uuid.UUID); ok {
        c.Logger().Infof("RequestID: %s", requestID)
    }
    resp := renderings.HealthCheckResponse{
        Message: "Everything is good!",
    }
    return c.JSON(http.StatusOK, resp)
}
```

`c.Get` type assertion ilə oxunur (`.(uuid.UUID)`), `ok` ilə mövcudluq yoxlanılır. Echo context — middleware-dən handler-ə state keçirməyin yeri (Ch5-də dərinləşir).

---

### 5. Contributed middleware — JWT auth

Echo contributing middleware dəsti ilə gəlir. JWT middleware-in daxili mahiyyəti:

```go
// JWT token parse + validasiya (konseptual):
token, err = jwt.ParseWithClaims(auth, claims, config.keyFunc)
if err == nil && token.Valid {
    // Store user information from token into context.
    c.Set(config.ContextKey, token)
    return next(c)
}
```

İstifadə — qrup səviyyəsində:

```go
reminderGroup := e.Group("/reminder")
reminderGroup.Use(middleware.JWT(signingKey))
reminderGroup.POST("", handlers.CreateReminder)
```

Login handler-in JWT hissəsi (Ch2-dən) ilə tam dövrə tamamlayır: login JWT token buraxır → sonrakı request-lərdə JWT middleware token parse edib context-ə user məlumatı qoyur → handler auth-dən xəbərsiz işləyir.

**Sətir-arxası:** handler-lər tam auth-dan xəbərsizdir — route-lar az-kompleks, asan-oxunan olur.

---

## Əsas terminlər

| Termin | İzah |
|--------|------|
| Middleware | Handler-i wrap edən, `next`-i parametr kimi qəbul edən funksiya |
| `MiddlewareFunc` | `func(HandlerFunc) HandlerFunc` — Echo middleware tipi |
| `e.Use` | Post-routing middleware qoşulması (instance/Group səviyyəsində) |
| `e.Pre` | Pre-routing middleware qoşulması (global) |
| Variadic route middleware | `e.GET("/", h, m1, m2)` — yalnız o route üçün |
| Middleware chaining | Nested funksiya çağrıları zənciri — `m1(m2(h))` |
| Zəncirin geri sarılması | Middleware `next` çağırmayanda call-stack geri dönür, cavab yaranır |
| 4 məntiqi sərhəd | Global pre / global post / group / route — 2 daxili yerə xəritələnir |
| `RequestIDMiddleware` | Hər request-ə UUID verən custom middleware nümunəsi |
| `c.Set` / `c.Get` | Context state yazmaq/oxumaq — middleware↔handler körpüsü |
| UUID (`satori/go.uuid`) | Unikal request identifikatoru — tracing üçün |
| Type assertion | `c.Get(key).(uuid.UUID)` — context-dən tipin geri alınması |
| Contributed middleware | Echo ilə gələn hazır dəst (Logger, Recover, JWT, TrailingSlash...) |
| `jwt.ParseWithClaims` | Token parse + imza yoxlaması — JWT middleware-in nüvəsi |
| Wrap (Sarmalama) | Bir funksiyanın digərinin ətrafında qatlanması |

---

## Praktik nəticə

1. **Təkrarı görəndə middleware düşün:** Eyni 2+ sətir birdən çox handler-də varsa (UUID yaratma, logging, auth) — middleware-ə çıxarın; handler-lər xəbərsiz qalsın.
2. **Sıraya diqqət:** `Use`/`Pre` çağrı sırası = icra sırası. Logger-in panic-i Recover tərəfindən tutulmalıdırsa, Recover xaricdə olmalı (əvvəl Use olunmalı). Pre vəziyyətində də eyni məntiq.
3. **"Hər request üçün" logic-i `Pre`-yə qoyun:** RequestID, trailing slash düzəlişi, pre-routing redirect — routing-dən əvvəl icra olunmalılar.
4. **Auth üçün qrup middleware:** `reminderGroup.Use(middleware.JWT(signingKey))` — bütün qrup üzvləri avtomatik qorunur.
5. **Middleware = adi funksiyalar:** Səhr yoxdur — `MiddlewareFunc` imzasına uyğun hər hansı funksiya middleware-dir; first-class funksiyalar sayəsində nested zəncir yaranır.
6. **State-i `c.Set` ilə ötürün:** Handler-də `c.Get(key).(Tip)` + `ok` pattern — təhlükəsiz oxuma.
7. **Helper funksiyalar middleware-i əvəz etmir:** Helper-lər hələ də hər handler-də çağrılmalıdır; middleware çağrını öz üzərinə götürür.

---

## Mənbə

- Kitab: *Echo Quick Start Guide* — J. Ben Huson, Packt Publishing, 2018 (ISBN 9781789340849)
- Chapter 4: "Implementing Middleware", book səh. 69–88
- PDF səhifələri: 87–105
- Kod: https://github.com/PacktPublishing/Echo-Essentials/tree/master/chapter4
- Video: https://goo.gl/u5TBL3
