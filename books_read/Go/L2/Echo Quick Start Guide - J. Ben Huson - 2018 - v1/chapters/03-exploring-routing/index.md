# Chapter 3 — Exploring Routing Capabilities

## Bu chapter nədən bəhs edir?

Bu chapter **routing-in daxili mexanizmlərini** açır: Echo-nun Radix tree axtarış alqoritmi, HTTP method-ların rolu (404 vs 405 düzəlişi — PR #205), group routing və versiyalanmış API-lər, trailing slash problemi, 27 framework-in benchmark müqayisəsi (Echo: 38,662 ns/op, 0 allocs), URL parametrləri (`:id`) və wildcard-ların (`*`) daxili implementasiyası.

---

## Əsas fikirlər

### 1. Routing nədir?

URI path komponentini handler koduna xəritələmə prosesi. Request-də həm target (resurs), həm də **method (niyyət)** var — method + path birlikdə hansı kodun icra olunacağını müəyyən edir.

**Add və metod-helper imzaları:**

```go
func (e *Echo) Add(method, path string, handler HandlerFunc,
    middleware ...MiddlewareFunc) *Route

// Hər metod üçün hazır helper:
func (e *Echo) GET(path string, h HandlerFunc, m ...MiddlewareFunc) *Route
func (e *Echo) POST(path string, h HandlerFunc, m ...MiddlewareFunc) *Route
func (e *Echo) PUT(path string, h HandlerFunc, m ...MiddlewareFunc) *Route
func (e *Echo) PATCH(path string, h HandlerFunc, m ...MiddlewareFunc) *Route
func (e *Echo) DELETE(path string, h HandlerFunc, m ...MiddlewareFunc) *Route
func (e *Echo) HEAD/OPTIONS/CONNECT/TRACE(...) *Route

// Bütün metodlara birdən:
func (e *Echo) Any(path string, handler HandlerFunc, m ...MiddlewareFunc) []*Route
```

Metod konstantları: `MethodGet`, `MethodPost`, `MethodPut`, `MethodPatch` (RFC 5789), `MethodDelete`, `MethodHead`, `MethodConnect`, `MethodOptions`, `MethodTrace`.

---

### 2. Radix tree — Echo routing necə işləyir

**Nədir?** Radix tree — prefiks ağacıdır: hər valideyn node-u öz uşaqlarının prefiksini təmsil edir. String match üçün ideal alət.

Kitab layihəsinin rotaları (`chapter3/cmd/service/main.go`):

```go
e.Static("/static", "static")                          // statik assetlər
reminderGroup := e.Group("/reminder")                   // JWT-qorumalı qrup
reminderGroup.Use(middleware.JWT(signingKey))
reminderGroup.POST("", handlers.CreateReminder)
e.GET("/health-check", handlers.HealthCheck)
e.POST("/login", handlers.Login)
e.POST("/logout", handlers.Logout)
```

Ağaq strukturu (sadələşdirilmiş):

```
/                      ← kök (bütün path-lər / ilə başlayır)
├── log                ← handler-siz prefiks node
│   ├── in  → Login
│   └── out → Logout
├── health-check → HealthCheck
└── reminder → ...
```

`/logout`-u tapmaq üçün: `/` → `log` → `out` — cəmi 3 səviyyə. **Böyük API-lərdə gücü görünür** — GitHub API-nin tam route set-i ağacda uzun ortaq prefikslər (`repos/.../events`) üzərindən effektiv axtarılır.

**Metod problemi — RFC uyğunsuzluğu (PR #205):**
- **Köhnə dizayn:** Hər HTTP method üçün ayrı Radix tree. `/reminder/123`-ə POST çatanda POST ağacında node tapılmırdı → **404 Not Found** qaytarılırdı.
- **Düzgün davranış:** Resurs mövcuddur, sadəcə POST icazəli deyil → **405 Method Not Allowed** olmalıdır.
- **Həll:** Node strukturasına `methodHandler` əlavə olundu — hər node öz daxilində method→handler xəritəsini daşıyır. Artıq **tək Radix tree** bütün metodlar üçün; node tapılan kimi ya handler icra olunur, ya 405 qaytarılır.

---

### 3. Group routing — versiyalanmış API-lər

**Nəyə lazımdır?** Eyni prefiksli route dəstlərini məntiqi qruplaşdırmaq — xüsusən API versiyalama (`/v1/login`, `/v2/login`). Request/response sxeması dəyişəndə yeni versiya yaratmaq best practice-dir — inteqrasiya edənlər üçün breaking change siqnalıdır.

```go
func (e *Echo) Group(prefix string, m ...MiddlewareFunc) (g *Group)
```

Group eyni route funksiyalarını daşıyır (Add, Any, GET, POST, PUT...). **İç-içə qruplar** mümkündür:

```go
// V1 Routes
v1 := e.Group("/v1")
// V1 Authentication routes
v1.POST("/login", handlers.Login)
v1.POST("/logout", handlers.Logout)
// V1 Reminder Routes
v1Reminders := v1.Group("/reminder", middleware.JWT(signingKey))
v1Reminders.POST("", handlers.CreateReminder)
```

**Daxili mexanizm:** Group sadəcə Echo-nun `Add` funksiyasını wrap edir — qrup prefiksi path-ə əlavə olunub router ağacına daxil edilir.

**Trailing slash problem:**

```bash
# Route: e.POST("/login", handlers.Login)
curl -XPOST http://localhost:8080/login/ -D -
HTTP/1.1 404 Not Found
{"message":"Not Found"}
```

Router **exact string match** edir — `/login/` ≠ `/login`. Hər endpoint üçün 2 route yazmaq cumbersome-dır; həll — `TrailingSlash` middleware (Ch4-də ətraflı), routing-dən əvvəl slash-ı tənzimləyir.

---

### 4. Router benchmark müqayisəsi

Mənbə: https://github.com/julienschmidt/go-http-routing-benchmark (GitHub API-nin tam route set-i, 2015):

| Router | ns/op | allocs/op |
|--------|-------|-----------|
| **Echo** | **38,662** | **0** |
| Gin | 43,467 | 0 |
| Denco | 83,114 | 167 |
| Ace | 93,675 | 167 |
| HttpRouter | 51,192 | 167 |
| Bear | 264,194 | 943 |
| GocraftWeb | 386,829 | 1,889 |
| Goji | 561,131 | 334 |
| Beego | 1,109,160 | 2,092 |
| **GorillaMux** | **7,431,130** | 1,791 |
| Martini | 10,261,331 | 2,686 |
| GoRestful | 15,569,513 | 7,725 |

**Nəticələr:**
- **Echo + Gin** liderlərdir; Echo 38,662 ns/op ilə ən sürətlisi və **0 allocation** təklif edir.
- 0 alloc səbəbi: context-in **context pool**-da reuse edilməsi (Ch5-də izah olunur).
- **Gorilla Mux (populyar amma yavaş):** Regex-əsaslı routing (Django kimi) — hər request-də BÜTÜN regex-lər üzərində iterasiya. Echo-dan **175 dəfə yavaş** (7.43 ms vs 38.7 µs).
- 7 ms əhəmiyyətsiz görünə bilər, amma bu, **request başına** hesablanır — çox request-də vaxt yığılır. Routing hər request-də icra olunduğundan performance kritikdir.

**Nüans:** Radix tree hər zaman ən yaxşı seçim deyil — axtarış O(k) complexity-si var (k = path uzunluğu). Çox uzun endpoint path-lərində Hash Map (standart kitabxanada builtin) daha yaxşı nəticə verə bilər.

---

### 5. URL parametrləri və wildcard-lar

#### Parametr (`:name`)

```go
v1Reminders.GET("/:id", handlers.GetReminder)   // /v1/reminder/:id

func GetReminder(ctx echo.Context) error {
    ctx.Logger().Info("Reminder id is: ", ctx.Param("id"))
    return nil
}
```

#### Versiya parametr kimi — maraqlı pattern

```go
e.GET("/v:version/reminder/:id", GetReminder)

func GetReminder(c echo.Context) error {
    switch c.Param("version") {
    case "1":
        return GetReminderV1(c)
    case "2":
        return GetReminderV2(c)
    }
}
```

Hər versiyanı ayrıca qeyd etmək əvəzinə parametr kimi yakalayıb yönləndirmə — routing performance-unu artırır.

#### Daxili implementasiya

Route `e.GET("/reminders/:id", ...)` qeydiyyatdan keçəndə:
- Axtarış stringi: `/reminders/:`
- `id` adı node-un metadata-sında **indeksli ad siyahısına** yazılır: `["verb", "adjective", "noun"]` kimi

Request gələndə path-in qalan hissəsi **dəyər siyahısına** yazılır. `c.Param("id")` çağrılanda: adın indeksi tapılır → eyni indeksli dəyər qaytarılır. Bu parallel siyahı cütlüyü çoxparametrli route-ları sadələşdirir.

Madlib nümunəsi:

```go
e.GET("/i/:verb/with/a/:adjective/:noun", ParameterMadlibHandler)

func ParameterMadlibHandler(c echo.Context) error {
    return c.String(http.StatusOK,
        fmt.Sprintf("I %s with a %s %s!",
            c.Param("verb"), c.Param("adjective"), c.Param("noun")))
}
```

```bash
curl http://localhost:8080/i/run/with/a/tall/woman -D -
# HTTP/1.1 200 OK — "I run with a tall woman!"
```

Node metadata: `["verb", "adjective", "noun"]` — dəyərlər eyni sıra ilə qarşılanır.

#### Wildcard (`*`)

Prefiksə uyğun HƏR suffix-i yakalar: `e.GET("/something/*", ...)` — `/something/` prefiksli bütün GET-lər. Ən yaxşı use case: **statik fayl xidməti** (`e.Static` daxildən wildcard istifadə edir).

**Məhdudiyyət:** Adlı wildcard yoxdur — `/i/*/from/a/*/*` kimi çox wildcard-lı routa dəyərlərə asan çıxış vermir. Odur ki, `*` adətən route-un sonu üçün variable suffix kimi istifadə olunur.

---

## Əsas terminlər

| Termin | İzah |
|--------|------|
| Routing | URI path + HTTP method → handler xəritələnməsi |
| Radix tree | Prefiks-ağacı — Echo-nun routing strukturu, O(k) axtarış |
| Prefiks node | Öz uşaqlarının ortaq prefiksini daşıyan ağac düyünü |
| `methodHandler` | Node daxilində method→handler xəritəsi — 405 düzəlişinin açarı |
| 404 vs 405 | Resurs yoxdur (404) vs metod icazəli deyil (405) — RFC uyğunluğu |
| `e.Add(method, path, h)` | Ümumi route qeydiyyatı — helper-lərin əsası |
| `e.Any` | Bütün HTTP metodlarına birdən route qeyd et |
| Group routing | Prefiks + middleware ilə route dəstləri (`e.Group("/v1")`) |
| API versiyalama | `/v1`, `/v2` prefiksləri — breaking change idarəsi |
| Trailing slash | `/login` vs `/login/` — exact match problem, TrailingSlash middleware həlli |
| `ns/op` / `allocs/op` | Benchmark vahidləri: nanosaniyə/əməliyyat, allocation/əməliyyat |
| Zero allocation | Echo-nun 0 allocs/op performansı — context pool reuse |
| Regex-based routing | Gorilla Mux/Django yanaşması — bütün regex-lərin ardıcıl yoxlanması |
| URL parametri (`:id`) | Path daxilində adlandırılmış dəyişən — `c.Param("id")` |
| İndeksli ad/dəyər siyahıları | Parametr adları və dəyərlərinin paralel siyahıları node metadata-sında |
| Wildcard (`*`) | Prefiksə uyğun bütün suffix-i yakalayan route elementi |

---

## Praktik nəticə

1. **Ortaq prefikslərdən istifadə edin:** Radix tree prefiks-ağacı olduğundan oxşar path-lər (məs., `/api/v1/users/...`) routing performance-unu artırır — API URL dizaynında bunu nəzərə alın.
2. **Route-ları qruplaşdırın:** `e.Group("/v1")` + iç-içə alt-qruplar + qrup-spesifik middleware (JWT) — versiyalama və funksional bölgü üçün.
3. **Trailing slash üçün middleware quraşdırın:** İki route yazmaq əvəzinə `TrailingSlash` middleware — `/login` və `/login/` hər ikisi işləyər.
4. **405-i gözləyin:** Modern Echo tək ağac + methodHandler sayəsində düzgün status kodları qaytarır — öz router-inizdə RFC uyğunluğunu unutmayın.
5. **Versiyanı parametr kimi yakalamaq:** `e.GET("/v:version/...")` — hər versiya üçün ayrı route siyahısından qaçış.
6. **Wildcard-ı statik məzmun üçün saxlayın:** `e.Static` daxildən `/*` istifadə edir; adlı dəyərlər lazımdırsa `:param` seçin.
7. **Benchmark-lərə baxın, amma kor-koranə inanmayın:** Echo/Gin router-ləri 0-alloc; amma çox uzun path-lərdə hash map nəzərə alın.

---

## Mənbə

- Kitab: *Echo Quick Start Guide* — J. Ben Huson, Packt Publishing, 2018 (ISBN 9781789340849)
- Chapter 3: "Exploring Routing Capabilities", book səh. 49–68
- PDF səhifələri: 70–86
- Kod: https://github.com/PacktPublishing/Echo-Essentials/tree/master/chapter3
- Video: https://goo.gl/SPqBes
- Router benchmark: https://github.com/julienschmidt/go-http-routing-benchmark
- 405 düzəlişi: https://github.com/labstack/echo/pull/205
