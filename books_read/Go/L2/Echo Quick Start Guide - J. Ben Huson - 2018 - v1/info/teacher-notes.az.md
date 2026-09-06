# Echo Quick Start Guide — Müəllim Qeydləri (Teacher Notes)

> Bu sənəd kitabı öyrədərkən istifadə üçün metodiki qeydləri, çətin anlaşılan yerləri, müzakirə suallarını və praktik tapşırıqları birləşdirir.

---

## 1. Auditoriya və ön şərtlər

- **Hədəf auditoriya:** Go əsaslarını bilən, veb-development-ə yeni keçən developer-lər.
- **Ön şərtlər:** Go sintaksisi (struct, interface, funksiya tipləri, pointer), HTTP-nin əsas anlayışı, komputer xəttində iş (Git, curl).
- **Kitabın versiya konteksti:** Go 1.9.4 + Echo v3.3.5 (2018). **Müasir dərslərdə diqqət:** Echo v4-də handler imzası eynidir (`func(echo.Context) error`), amma import yolu `github.com/labstack/echo/v4`, `e.Start`/`e.Logger` API-lərində fərqlər var; `dep` əvəzinə Go modules gəlir. Konseptual biliklərin hamısı tətbiq olunandır.

---

## 2. Tədris axını (8 module üzrə plan)

| Module | Chapter | Konsentrasiya | Süret (təxmini) |
|--------|---------|---------------|------------------|
| 1 | Ch1 | HTTP + net/http primitivləri | 10% |
| 2 | Ch2 | Layihə strukturu + ilk handler-lər | 15% |
| 3 | Ch3 | Routing (Radix tree, group, benchmark) | 12% |
| 4 | Ch4 | Middleware (3 qoşulma yolu, custom) | 15% |
| 5 | Ch5 | Context + Bind + Validate + Render | 18% |
| 6 | Ch6 | Logging + Error + Recover | 10% |
| 7 | Ch7 | Test (unit/benchmark/integration) | 12% |
| 8 | Ch8 | Static + Templates | 8% |

**Tək layihə prinsipi:** Kitab boyu bir layihə (health-check → login → reminder API) genişlənir — dərslərdə hər module-da eyni layihə üzərində işləyin, tələbələr kodun artımını görür.

---

## 3. Çətin anlaşılan anlar və izah üsulları

### a) `func(Context) error` vs `ServeHTTP(w, r)` (Ch1-2)
**Çətinlik:** Nə üçün Echo handler imzasını dəyişib?
**İzah:** İki parametrli imza context saxlamağa imkan vermir (Go 1.7-dən əvvəl); tek context parametri həm request, həm response, həm state-i bir obyektdə birləşdirir. Yan təsir — vendor lock-in (Ch5-də müzakirə olunur).

### b) Middleware zəncirinin sırası (Ch4)
**Çətinlik:** `Use` çağrılarının sırası niyə icra sırasıdır?
**İzah:** `Use(Logger)` sonra `Use(Recover)` → `Logger(Recover(h))` — Logger XARİCDƏDIR. Logger-in içində panic olsa, Recover (içəridə) onu tutur. Əks sıra olarsa, Recover-in öz panic-i Logger-dən sonra baş verən Logger panic-ini tutan bilməz. Taxtaxtaxta funksiya kompozisiyasını çəkin: `m1(m2(m3(h)))`.

### c) 404 vs 405 (Ch3)
**Çətinlik:** Niyə method-a görə ayrı ağac düzgün deyildi?
**İzah:** `/reminder/123` mövcuddur amma POST icazəli deyil → düzgün cavab 405 Method Not Allowed. Ayrı ağaclarda POST ağacı boş → 404 Not Found (səhv: resurs "yoxdur"). Həll: tək ağac, node daxilində method→handler xəritəsi.

### d) `sync.Pool` və 0 allocation (Ch3, Ch5)
**Çətinlik:** Hər request üçün context yaranırsa niyə allocation yoxdur?
**İzah:** Pool istifadə olunmuş context-ləri saxlayır; yeni request köhnə instansiyasını "icarəyə götürür", iş bitincə geri qaytarır. GC təzyiqi azalır — benchmark-da 0 allocs/op-un səbəbi.

### e) Bind-in arxasında nə olur (Ch5)
**İzah:** `Content-Type` header-ə baxılır → JSON/XML/form deserializatoru seçilir → `json.Unmarshal`-ə bənzər avtomatika. Struct tag-lər (`json:"username"`) sahə xəritələnməsini idarə edir. XML testi üçün curl nümunəsini canlı göstərin — eyni handler iki format!

### f) Test-də server yoxdur — nəticə necə yaranır (Ch7)
**Çətinlik:** `e.ServeHTTP(w, r)` — portu dinləmədən necə işləyir?
**İzah:** Echo server = `net/http` handler; `ServeHTTP` onu birbaşa çağırır — TCP qatı lazımsızdır. `ResponseRecorder` — `ResponseWriter`-in yadda saxlayan saxtası. Diaqram: request → Echo → handler → recorder → assertion.

### g) Parse bir dəfə, Execute çox (Ch8)
**İzah:** Template parse = string-in ağac strukturu parse edilməsi (baha); Execute = data ilə render (ucuz). Hər request-də parse etmək — klassik səhv; entry point-də bir dəfə, `e.Renderer`-da saxla.

---

## 4. Müzakirə sualları (dərs üçün)

**Ch1-2:**
1. Standart kitabxana minimalizmi — design qərarı mı, qüsur mu? Hansı hallarda framework-süz qalmaq daha yaxşıdır?
2. Niyə bindings/renderings/models AYRI qovluqlar? Hamısını bir `types/` qovluğunda saxlasaq nə itiririk?

**Ch3:**
3. Radix tree vs hash map — hansı API formasında hansı qalib gəlir? (Uzun ortaq prefikslər → Radix; çox qısa unikal path-lər → map)
4. `/v:version/...` parametr pattern-i vs `/v1` qrup — hansı vaxt hansı? (Bir handler çox versiya idarə edəndə vs middleware fərqli olanda)

**Ch4:**
5. Logger Recover-dən əvvəl YOXSA sonra Use olunmalıdır? (Logger xaricdə = request logları həmişə yazılır; Recover xaricdə = Logger-in öz panic-i də tutulur — prioritetə qərar verin)
6. RequestID nə üçün Pre-də olmalıdır? (Routing-dən əvvəl gələn hər request — hətta 404-lər — izlənə bilsin)

**Ch5:**
7. Gorilla/context yanaşmasının memory leak-i nədən yaranır? (Request pointer-ləri map-də yığılır, təmizlənməzsə)
8. Validation strukturun YANINDA niyə daha yaxşıdır? (Yenidən istifadə, test oluna bilərlik, handler təmizliyi)

**Ch6:**
9. Production-da səviyyə niyə DEBUG yox INFO? (Verbosluq, performans, gündəlik analiz üçün siqnal/shum nisbəti)
10. Secret-lərin log-a düşməməsi üçün hansı prosedur lazımdır? (Sanitasiya, səviyyə intizamı, code review)

**Ch7:**
11. 100% coverage yaxşı hədəf dərəcəsidirmi? (Yox — assertion keyfiyyəti; over-mock tələsi)
12. MockDB pattern-ində funksiya sahələri niyə lazımdır? (Default passiv davranış + test-də xüsusi override)

**Ch8:**
13. `html/template` vs `text/template` — email şablonu hansı? (Text — HTML context yoxdursa; XSS riski yoxdur)

---

## 5. Praktik tapşırıqlar (təkmilləşdirmə sırası ilə)

**Tapşırıq 1 (Ch1):** Saf `net/http` ilə `/health` endpoint yazın; sonra əlinizlə JSON marshal edin. Bu "əziyyəti" hiss etmək framework dəyərini göstərir.

**Tapşırıq 2 (Ch2):** Kitabın strukturunu təkrarlayın: bindings/renderings/handlers üçlüsü ilə `/ping` API. `dep init` (və ya `go mod init`) icra edin.

**Tapşırıq 3 (Ch3):** `/users/:id`, `/users/:id/posts/:postId`, `/v1/...` və `/v2/...` qrupları yaradın; handler-də hər iki parametri çap edin. TrailingSlash middleware əlavə edib `/v1/users/1/` sorğusunu test edin.

**Tapşırıq 4 (Ch4):** `RequestIDMiddleware` yazın (UUID + `c.Set`) — Pre ilə qoşun; handler-də `c.Get` ilə oxuyun. Sonra middleware-in `next`-i ÇAĞIRMAYAN variantını yazın (early return) — 401 cavab observer olsun.

**Tapşırıq 5 (Ch5):** `LoginRequest` üçün Bind + Validate + `echo.HTTPError` ilə dəqiq statuslar. XML curl testi ilə content negotiation göstərin.

**Tapşırıq 6 (Ch6):** `middleware.Logger()` + `middleware.Recover()` qoşun; panic edən handler yazıb Recover-in stack trace-i gizlətdiyini görün. Custom `HTTPErrorHandler` yazın.

**Tapşırıq 7 (Ch7):** HealthCheck üçün unit test + benchmark (`b.N`); `MockDB` ilə xəta halı testi. `-cover` çıxışını şərh edin.

**Tapşırıq 8 (Ch8):** `e.Static` qoşun; `CustomTemplate` + `e.Renderer` ilə cədvəl səhifəsi; route adı + `TmplData.Reverse` ilə login linki.

---

## 6. Sınav sualları (nümunə)

**Asan:**
1. Echo handler-inin imzası nədir? (`func(echo.Context) error`)
2. `e.Pre` ilə `e.Use` arasında fərq? (routing-dən əvvəl / sonra)
3. `c.Bind` hansı content-type-ları tanıyır? (JSON, XML, form-urlencoded)

**Orta:**
4. Radix tree nə üçün regex-dən effektiv olur? (Prefiks paylaşımı — bütün regex-lərin ardıcıl yoxlanışı yoxdur; ortaq prefikslərin tək yolu)
5. `echo.Response`-də `Committed` nəyi qoruyur? (İkiqat WriteHeader/Write — cavab artıq göndərilibsə warn)
6. `MockableDB` interface yaratmağın məqsədi? (sql.DB-ni interface arxasınca gizlədib unit test-də mock keçirmək)

**Çətin:**
7. Echo-nun 0 allocs/op benchmark nəticəsinin arxasında nə dayanır? (sync.Pool context reuse)
8. Instrumented server testinə (TestRun + /stop-test-server) nə üçün ehtiyac var? (Real xarici testlər üçün canlı server + dəqiq coverage toplamaq)
9. Middleware sırası `Logger` → `Recover` olduqda hansı panic tutulmur? (Recover-in özündən SONRA gələn — yəni Logger xaricdədirsə, Recover-in işindən sonra gələn heç nə; sıralamaya diqqət)

---

## 7. Kollektiv layihə ideyası (modul sonu)

**"Bookshelf API"** — Echo üzərində kitab kolleksiyası idarəsi:
- Auth: login/logout (JWT) + `auth` middleware qrupu
- Resource: `/v1/books` CRUD — `:id` parametrli, MockDB ilə testlənən `models`
- Keyfiyyət: `Logger` + `Recover` + `RequestID` (custom, Pre)
- Sınaq: unit + benchmark + `-coverpkg` integration
- UI: static CSS + reminders cədvəli template-i + `.Reverse` ilə linklər

Təqdimat meyarları: handler-lərdə yalnız business logic; bütün input bindings-də + Validate; error-lar `echo.HTTPError` ilə; coverage ≥ 60%.

---

## 8.Əlavə resurslar

- Echo rəsmi sənədləri: https://echo.labstack.com
- Kitabın kod reposu: https://github.com/PacktPublishing/Echo-Essentials
- Router benchmark: https://github.com/julienschmidt/go-http-routing-benchmark
- Go blogları: "Error handling and Go", "Defer, Panic and Recover" — xəta/panic dərslərində
- Monkey patching: https://bou.ke/blog/monkey-patching-in-go/
- Echo 405 PR: https://github.com/labstack/echo/pull/205
