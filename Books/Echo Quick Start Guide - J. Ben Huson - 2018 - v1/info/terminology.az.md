# Echo Quick Start Guide — Terminologiya (Azərbaycanca)

> Kitab boyu rast gəlinən terminlərin Azərbaycanca izahları. Texniki terminlər `English (Azərbaycanca qarşılığı)` formatında.

## A

**API versiyalama (API versioning)** — `/v1`, `/v2` prefiksləri ilə route dəstlərinin ayrılması; breaking change-lərin idarəsi üçün best practice.

**Attachment** — `c.Attachment(file, name)` — faylı klientə yükləmə (download) kimi göndərən cavab metodu.

**Auth/Authentication (Autentifikasiya)** — istifadəçinin kimliyinin yoxlanması; kitabda bcrypt parol + JWT token ilə.

## B

**Benchmark testing (Benchmark testi)** — kod sürətinin/throughput-unun ölçülməsi (`BenchmarkXxx(b *testing.B)`, `ns/op`, `allocs/op`).

**Bare bones server (Yalın server)** — yalnız əsas funksionallıq daşıyan, əlavə abstraksiyasız server (Go `net/http`).

**bcrypt** — parol hash müqayisə alqoritmi; `bcrypt.CompareHashAndPassword(hash, password)`.

**Binding (Bağlama)** — `c.Bind(&x)` — request body-nin `Content-Type`-a görə struktura avtomatik deserializasiyası (JSON/XML/form).

**bindings/** — layihə qovluğu: protokol-spesifik request input strukturları.

**BLOB rendering** — `c.Blob(code, contentType, b)` — istənilən content-type ilə xam baytların cavabı.

## C

**Chaining (Zəncirləmə)** — middleware-lərin nested funksiya çağrıları şəklində birləşməsi: `m1(m2(h))`.

**Committed flag** — `echo.Response`-də cavabın artıq göndərildiyini bildirən bayraq; double-write qoruması.

**Content negotiation (Məzmun danışıqları)** — API-nin eyni endpoint-dən JSON/XML/form formatlarında xidmət etməsi.

**Context (`echo.Context`)** — Echo handler-inin tək parametri; request/response + state + helper-lər. `sync.Pool`-dan reuse olunur.

**Contributed middleware** — Echo ilə gələn hazır middleware dəsti (Logger, Recover, JWT, TrailingSlash, CORS...).

**curl** — command-line HTTP klienti; kitabda bütün endpoint testləri üçün istifadə olunur.

## D

**dep** — Go dependency management tool; `dep init` vendor qovluğu yaradır, versiyaları pinləyir.

**Debug level (DEBUG səviyyəsi)** — troubleshoot üçün; production-da adətən söndürülür.

**Deployment artifact (Deploy artifactı)** — deploy edilən icra paketi; statik fayllar daxilində daşıya bilər.

## E

**Error level (ERROR səviyyəsi)** — bərpa edilən, amma təhqiqat tələb edən tətbiq xətası üçün log səviyyəsi.

**Explicit error check (Açıq xəta yoxlaması)** — Go konvensiyası: exception throw yox, xətanın yaranma yerində yoxlanması.

## F

**Fail with grace (Layiqcə uğursuzluq)** — xəta halında düzgün cavabla (stack trace yox) uğursuz olmaq.

**Fatal/Panic level** — bərpa mümkün olmayan xətalar; adətən başlanğıc zamanı (DB connect, port bind).

**First-class functions (Birinci sinif funksiyalar)** — funksiyaların dəyişən kimi istifadəsi; Go middleware-in əsası.

**FormValue/FormFile** — `c.FormValue(name)` form dəyəri; `c.FormFile(name)` multipart fayl.

## G

**Global context mapping** — request pointer-inin qlobal map key kimi istifadəsi (gorilla/context); memory leak təhlükəsi.

**Goroutine** — Go-nun yüngül thread-i; Go web server hər request üçün yeni goroutine başladır.

**Group routing** — `e.Group("/v1")` — prefiks + middleware ilə route dəstləri.

## H

**Handler (`http.Handler` / `HandlerFunc`)** — request-i emal edən kod vahidi. Std: `ServeHTTP(w, r)`; Echo: `func(echo.Context) error`.

**html/template** — HTML üçün təhlükəsiz Go şablon paketi (code injection qorumalı).

**HTTPErrorHandler** — handler-dən qaytan error-u HTTP cavabına çevirən, əvəz edilə bilən funksiya (`e.HTTPErrorHandler`).

**httptest.ResponseRecorder** — `http.ResponseWriter`-in yadda saxlayan test implementasiyası; `e.ServeHTTP(w, r)` ilə server-siz test.

## I

**Information leakage (Məlumat sızması)** — secret-lərin təsadüfən log-a düşməsi.

**Integration testing (İnteqrasiya testi)** — real xarici sistemlərlə (DB və s.) birlikdə test.

**Instrumented build (İnstrumented build)** — coverage ölçən kodu daşıyan icra build-i (`-coverpkg`).

## J

**JWT (JSON Web Token)** — token-based auth; `jwt.StandardClaims` (ExpiresAt, Issuer), HS256 signing.

## L

**Latency (Gecikmə)** — request emal müddəti; Logger middleware-in JSON sahəsi.

**Lock contention (Lock rəqabəti)** — çoxlu goroutine-in eyni map-ə yazma rəqabəti; qlobal map həllərinin problemi.

**Logger interface** — Echo-nun jurnal abstraksiyası; `SetLevel`, `Xxx/Xxxf/Xxxj` metodları.

**Logging (Jurnallama)** — tətbiq vəziyyətinin iz qeydləri; troubleshooting əsası.

## M

**Memory leak (Yaddaş sızıntısı)** — təmizlənməyən qlobal map-lərdə yığılan yaddaş (gorilla/context riski).

**Message body (Mesaj gövdəsi)** — HTTP mesajının payload-ı; `io.ReadCloser` kimi modellənir.

**Middleware** — handler-i wrap edən, `next`-i parametr kimi qəbul edib handler qaytaran funksiya (`MiddlewareFunc`).

**Mocking** — interface arxılına saxta asılılıq davranışının qoyulması (`MockableDB`/`MockDB`).

**Monkey patching** — run-time-da funksiya təyinatının əvəz edilməsi; paralel test itkisi ilə.

## O

**Over-mock** — həddindən artıq mock istifadəsi; coverage şişirir, integration problemlərini gizlədir.

## P

**Param (`c.Param`)** — URL path parametrinin (`:id`) handler-də oxunması.

**Pipeline (şablonlarda)** — `{{.Field}}` kontekst dəyəri/ifadəsi.

**Pre/Use** — `e.Pre(mw)` routing-dən əvvəl; `e.Use(mw)` routing-dan sonra.

**Prefiks node** — Radix tree-də öz uşaqlarının ortaq prefiksini daşıyan düyün.

**QueryParam** — `c.QueryParam("name")` — query string parametrinin oxunması.

## R

**Radix tree** — prefiks-ağacı; Echo-nun routing data strukturu, O(k) axtarış.

**Range (şablon)** — `{{range .Items}}` — kolleksiya üzərində iterasiya; `{{else}}` boş halda.

**Recover middleware** — `recover()` ilə panic-i tutub generic error-a çevirən qoruyucu middleware; mütləq istifadə edilməli.

**renderings/** — layihə qovluğu: response serializasiya strukturları.

**Rendering** — cavabın hazır helper-lərlə yazılması (`c.JSON/String/XML/HTML/File/...`).

**Request-line/Status-line** — HTTP mesajının ilk sətri (request-də method+target+version; response-da version+code+phrase).

**Resilient (Davamlı)** — xətalara baxmayaraq işləməyə davam edə bilmə qabiliyyəti.

**ResponseWriter** — std cavab interfeysi: `Header()`/`WriteHeader()`/`Write()`.

**Reverse URL** — adlandırılmış route-un URL generasiyası: `e.Reverse("login")`; şablonlarda `TmplData.Reverse`.

**Router** — path→handler xəritələyicisi; Echo Radix tree, Gorilla Mux regex əsaslı.

## S

**Security testing** — XSS, SQL injection, CORS kibi təhlükəsizlik yoxlaması (penetration testing).

**Server-side rendering (SSR)** — dinamik məzmunun serverdə şablonla render edilməsi.

**SigningKey** — JWT imza açarı; middleware-dən context-ə `c.Set` ilə ötürülür.

**Start-line** — HTTP mesajının ilk sətri (ümumi termin).

**Stateless protocol (Stateless protokol)** — hər request müstəqil; server hallar arası state saxlamır.

**Static content (Statik məzmun)** — CSS/JS/şəkil kimi dəyişməyən assetlər; `e.Static`/`e.File`.

**Stack trace (Stek izi)** — panic tutulmasa klientə düşən daxili kod izi — təhlükəli sızma.

**Structured logging (`Xxxj`)** — JSON strukturlu log yazımı.

**sync.Pool** — müvəqqəti instansların reuse hovuzu; Echo context 0-alloc edir.

**Vendor lock-in** — kodun yalnız bir framework-ə bağlılığı; yeni handler imzasının qiyməti.

## T

**testing.T / testing.B** — unit/benchmark test strukturları; `Fail`, `Log`, `Error` helper-ləri.

**Trailing slash** — `/login` vs `/login/` fərqi; exact match 404 verir → TrailingSlash middleware.

**Type assertion** — `c.Get(key).(uuid.UUID)` — interface-dən konkret tipin geri alınması.

## U

**Unit testing** — ən kiçik kod vahidinin izolyə testi; `_test.go` + `TestXxx(t *testing.T)`.

**UUID** — unikal identifikator; `satori/go.uuid` — RequestID üçün.

## V

**Validatable interface** — `Validate() error` daşıyan strukturlar üçün abstraksiya; `e.Validator` ilə qeydiyyat.

**Vendor directory** — layihə daxilində dependency kodunun saxlandığı qovluq.

## W

**Wildcard (`*`)** — prefiksə uyğun bütün suffix-i yakalayan route elementi; statik fayl üçün ideal.

**With (`{{with}}`)** — şablon konstruksiyası: pipeline boş deyilsə bloku render et.
