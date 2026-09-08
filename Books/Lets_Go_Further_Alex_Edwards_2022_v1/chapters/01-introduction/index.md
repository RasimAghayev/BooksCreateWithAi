# Chapter 1-2 — Introduction & Getting Started (Greenlight API)

## Bu hissə nədən bəhs edir?

Kitabın girişi və layihənin təməlinin qurulması: **Greenlight** — filmlər haqqında
məlumat idarə edən production-istiqamətli JSON API. Bu fəsillərdə skeleton qovluq
strukturu, konfiqurasiya (command-line flags), Dependency Injection (asılılıqların
inyeksiyası) ilə handler-lərin təchizatı, basic HTTP server və httprouter ilə
RESTful routing qurulur.

**Prerequisite-lər:** Kitab "Let's Go"nun davamıdır — əsaslar detallı izah
olunmur; Go 1.20, curl, hey (load test), PostgreSQL, golang-migrate, make
tələb olunur.

## Əsas fikirlər

### 1. Layihənin skeleton strukturu
**Nədir:** Production Go layihəsi üçün standart qovluq iyerarxiyası.

**Necə işləyir:**
```bash
mkdir -p $HOME/Projects/greenlight
cd $HOME/Projects/greenlight
go mod init greenlight.alexedwards.net
mkdir -p bin cmd/api internal migrations remote
touch Makefile
touch cmd/api/main.go
```

**Struktur:**
```
greenlight/
├── bin/          ← kompilyasiya olunmuş binary-lər (deploy üçün)
├── cmd/api/      ← applikasiyaya xas kod (server, handler-lər, auth)
├── internal/     ← köməkçi paketlər (DB, validation, email) — xarici import QADAĞAN
├── migrations/   ← SQL migration faylları
├── remote/       ← production server konfiqurasiyaları
├── go.mod        ← modul yolu + asılılıq versiyaları
└── Makefile      ← adminstrativ task avtomatlaşdırması
```

**Nəyə lazımdır:** Kodun səliqəli bölgüsü; `internal/` Go compiler-i tərəfindən
müdafiə olunur — bu paketləri yalnız parent modul daxilindən import etmək olar.

**Üstünlükləri:**
- Reproducible builds (təkrarlanan build-lər) — go.mod dəqiq versiyaları saxlayır
- Layihə böyüdükcə strukturu dəyişmək lazım deyil

**Çatışmamazlıqları:**
- Yeni başlayanlar üçün başlanğıcda "çox qovluq" təəssüratı yarada bilər

**Kitabdan kod nümunəsi:**
```go
// File: go.mod
module greenlight.alexedwards.net
go 1.20
```
**Sub-kod izahı:**
- `module greenlight.alexedwards.net` → unikal modul identifikatoru; import
  kök yolu kimi istifadə olunur (öz URL-inizə əsaslanmaq yaxşı praktikadır)
- `go 1.20` → minimum Go versiyası

### 2. Konfiqurasiya və Dependency Injection pattern-i
**Nədir:** Bütün konfiqurasiya command-line flag-lərdən struct-a oxunur;
asılılıqlar `application` struct-ında toplanır və handler-lərə **method** kimi
çatdırılır — qlobal dəyişən YOX.

**Necə işləyir:**
```go
type config struct {
    port int
    env  string
}

type application struct {
    config config
    logger *log.Logger
}

func main() {
    var cfg config
    flag.IntVar(&cfg.port, "port", 4000, "API server port")
    flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
    flag.Parse()

    logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

    app := &application{
        config: cfg,
        logger: logger,
    }
    // ...
}
```

**Sub-kod izahı:**
- `flag.IntVar(&cfg.port, "port", 4000, ...)` → `-port` flag-ini `cfg.port`-a
  default 4000 ilə bağlayır
- `flag.Parse()` → flag dəyərlərini faktiki oxuyur
- `log.New(os.Stdout, "", log.Ldate|log.Ltime)` → tarix+vaxt prefiksli logger
- `app := &application{...}` → bütün asılılıqlar bir mərkəzdə; handler-lər
  `app.config.env` kimi çatırlar

**Nəyə lazımdır:** Test yazarkən application struct-ını saxta asılılıqlarla
konstruksiya etmək mümkün olur; qlobal state aradan qalxır.

**Üstünlükləri:**
- Handler-lər qlobal dəyişən və ya closure olmadan asılılıqlara çatır
- Test ediləbilənlik (testability) yüksəlir

**Çatışmamazlıqları:**
-struct böyüdükcə çoxsaylı asılılıq sahəsi yığılır (kitab boyunça normal
qəbul edilir)

### 3. Sensible HTTP server timeout-ları ilə
**Nədir:** `http.Server` yaradılırkən mütləq timeout təyin etmək.

**Kitabdan kod nümunəsi:**
```go
srv := &http.Server{
    Addr:         fmt.Sprintf(":%d", cfg.port),
    Handler:      mux,
    IdleTimeout:  time.Minute,
    ReadTimeout:  10 * time.Second,
    WriteTimeout: 30 * time.Second,
}
err := srv.ListenAndServe()
logger.Fatal(err)
```

**Sub-kod izahı:**
- `IdleTimeout: time.Minute` → boş (keep-alive) bağlantı 1 dəqiqədən sonra
  qapatılır — yavaş lori (slowloris) hücumlarına qarşı qoruyur
- `ReadTimeout: 10 * time.Second` → sorğunun tam oxunması üçün maksimum vaxt
- `WriteTimeout: 30 * time.Second` → cavabın yazılması üçün maksimum vaxt
- `logger.Fatal(err)` → server dayandıqda loglayıb prosesi sonlandırır

**Nəyə lazımdır:** Default `http.ListenAndServe()` heç bir timeout təyin etmir —
production-da mütləq lazımdır.

### 4. RESTful endpoint strukturu və HTTP metodları
**Nədir:** Eyni URL pattern-i fərqli metodlarla fərqli handler-lərə gedir.

**Cədvəl (bütün kitab boyunca qurulacaq API):**
| Method | URL Pattern | Handler | Action |
|---|---|---|---|
| GET | /v1/healthcheck | healthcheckHandler | Applikasiya məlumatı |
| GET | /v1/movies | listMoviesHandler | Bütün filmlər |
| POST | /v1/movies | createMovieHandler | Yeni film yarat |
| GET | /v1/movies/:id | showMovieHandler | Xüsusi film |
| PUT | /v1/movies/:id | editMovieHandler | Filmi yenilə |
| PATCH | /v1/movies/:id | — | Partial yenilə |
| DELETE | /v1/movies/:id | deleteMovieHandler | Filmi sil |

**Metod semantikası:**
- `GET` → yalnız məlumat oxuyur, state dəyişmir
- `POST` → non-idempotent yaratma əməliyyatı
- `PUT` → idempotent tam yeniləmə
- `PATCH` → partial yeniləmə (idempotent ola da, olmaya da bilər)
- `DELETE` → resursu silir

### 5. httprouter inteqrasiyası
**Nədir:** `http.ServeMux` (standart router) metod-əsaslı routing və URL
parametrlərini dəstəkləmir — `httprouter` bu boşluğu doldurur.

**Necə işləyir:**
```bash
go get github.com/julienschmidt/httprouter@v1
```

```go
// File: cmd/api/routes.go
func (app *application) routes() *httprouter.Router {
    router := httprouter.New()
    router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
    router.HandlerFunc(http.MethodPost, "/v1/movies", app.createMovieHandler)
    router.HandlerFunc(http.MethodGet, "/v1/movies/:id", app.showMovieHandler)
    return router
}
```

**Sub-kod izahı:**
- `router.HandlerFunc(http.MethodGet, ...)` → metoddan asılı olaraq handler
  qeydiyyatı
- `/v1/movies/:id` → `:id` interpolated parametrdir
- Routing qaydaları `routes()` metodunda Encapsulation (kapsullaşdırma) —
  main() təmiz qalır, test kodu asanlıqla router-ə çatır

**Üstünlükləri:**
- Radix tree sayəsində çox sürətli URL matching
- Avtomatik `405 Method Not Allowed` + `Allow` header
- Avtomatik `OPTIONS` cavabı

**Çatışmamazlıqları:**
- Conflicting routes (ziddiyyətli marşrutlar) QADAĞANDIR: `GET /foo/new` və
  `GET /foo/:id` eyni vaxtda qeydiyyatdan keçə bilməz. REST strukturunda problem
  deyil; lazım olsa `chi`, `flow`, `pat` routerlərinə baxın.

### 6. URL parametrinin oxunması — helper
**Nədir:** `:id` parametrini context-dən oxuyub int64-ə çevirən yenidən
istifadə olunan helper.

**Kitabdan kod nümunəsi:**
```go
// File: cmd/api/helpers.go
func (app *application) readIDParam(r *http.Request) (int64, error) {
    params := httprouter.ParamsFromContext(r.Context())
    id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
    if err != nil || id < 1 {
        return 0, errors.New("invalid id parameter")
    }
    return id, nil
}
```

**Sub-kod izahı:**
- `httprouter.ParamsFromContext(r.Context())` → request context-dən parametr
  dilimini götürür
- `strconv.ParseInt(..., 10, 64)` → base-10, 64-bit konversiya; uğursuzsa və ya
  `id < 1` olarsa xəta qaytarır → handler `404 Not Found` göndərir
- Asılılıq istifadə etməsə də method kimi yazmaq — struktur tutarlılığı və
  gələcək üçün (kitabın tövsiyəsi)

### 7. API versiyalanması (Additional Information)
**Nədir:** Backwards-incompatible dəyişikliklər üçün API versiya strategiyası.

**İki yanaşma:**
1. URL prefix: `/v1/healthcheck`, `/v2/healthcheck` — istifadəçi dostudur,
   brauzerlə görünür (kitab bunu seçir)
2. Header ilə: `Accept: application/vnd.greenlight-v1` — HTTP semantikası
   baxımından "təmiz", amma praktik deyil

## Bu fəsildə layihəyə nə əlavə olundu? (project continuity)

**Əvvəlki vəziyyət:** Boş layihə.
**Bu fəsildə əlavə olunan:** go mod, skeleton qovluqlar, main.go (config +
flags + server), healthcheck.go, routes.go, movies.go (placeholder
handler-lər), helpers.go (readIDParam).

**Yeni fayllar:**
- `cmd/api/main.go` — giriş nöqtəsi, konfiqurasiya, server
- `cmd/api/healthcheck.go` — status endpoint-i
- `cmd/api/routes.go` — httprouter qeydiyyatları
- `cmd/api/movies.go` — createMovieHandler, showMovieHandler (placeholder)
- `cmd/api/helpers.go` — readIDParam helper-i

## Əsas terminlər

- Dependency Injection (asılılıqların inyeksiyası) — asılılıqların struct
  sahələri kimi çatdırılması
- Encapsulation (kapsullaşdırma) — routing qaydalarının ayrı metoddа saxlanması
- Idempotency (eyniliyin qorunması) — eyni əməliyyatın təkrarı eyni nəticə verir
- Interpolated URL Parameter (URL-ə yerləşdirilmiş parametr) — `:id` kimi
  dinamik URL hissəsi
- Radix Tree (radiks ağacı) — httprouter-in sürətli URL matching strukturu
- Reproducible Build (təkrarlanan build) — dəqiq asılılıq versiyaları ilə build

## Praktik nəticə

Production Go API-si üçün etibarlı başlanğıc arxitekturası: konfiqurasiya
flag-lərdən, asılılıqlar mərkəzi `application` struct-ında, routing təcrid
olunmuş `routes()` metodunda, server isə mütləq timeoutlarla. Bu pattern bütün
kitab boyunca təkrarlanır və istənilən Go API layihəsinə daşına bilər.

## Mənbə
Pages: 1-30 (raw 006-030)
