# Chapter 2 — Structuring Go Web Application (Go Veb Tətbiqinin Strukturu)

## Bu chapter nədən bəhs edir?
Go Modules-in dərinliyinə, standart layihə qovluq strukturuna, package management + vendoring-ə, MVC pattern-inin Go tətbiqinə, əsas Application obyektinin proqramlaşdırılmasına, konfiqurasiya/env dəyişənlərinə, Dependency Injection-a və health check endpoint-lərinə.

## Əsas fikirlər

### 1. Go Modules — dependency idarəsinin təkamülü
**Nədir:** go.mod manifestli paket toplusu — GOPATH-dən asılılıq aradan qaldırıldı (v1.11+).

**Üstünlüklər:**
- **Versioning:** dəqiq versiya sabitlənir — build-lər ardıcıl
- **İzolyasiya:** hər layihə öz dependency dəsti — konflikt yox
- **go.sum:** kriptoqrafik checksum-lər — təhlükəsizlik + tutarlılıq
- **Lokasiya azadlığı:** fayl sistemi, şəbəkə, proxy — istənilən yerdən

**Əmrlər:**
```bash
go mod init <module-name>        # yeni modul (go.mod yaradır)
go get <module>@<version>        # dəqiq versiya
go mod tidy                      # istifadə edilməyənləri təmizlə
go get -u <package>              # update
go mod vendor                    # /vendor qovluğunu doldur (opt-in)
```

### 2. Module Proxy (modul proksisi)
**Nədir:** Layihə ilə source repo arasındaki keshləyici aracı (GitHub-dan alır, saxlayır, təkrar sorğuları cachedən verir).
- Sürət: eyni modul təkrar endirilmir
- Etibarlılıq: orijinal source silinsə belə cached versiya qalır

### 3. Standart Go layihə strukturu
```
gitforgits-bookstore/
├── cmd/server/          # main funksiya — giriş nöqtəsi
├── pkg/models/         # Book, User strukturları
├── pkg/utils/          # auth.go, helpers.go
├── api/definitions/     # OpenAPI/gRPC müqavilələri
├── web/static/         # CSS, JS, şəkillər
├── web/templates/      # HTML şablonlar
├── internal/handlers/  # HTTP handler-lər (kitab, istifadəçi)
├── internal/middleware/# logging, CORS, auth
├── internal/config/    # konfiqurasiya yükləməsi
├── scripts/             # deploy.sh, test.sh
├── database/migrations/# SQL migrasiya skriptləri
├── database/seeds/      # ilkin data
├── docs/                # setup.md, api.md
├── go.mod + go.sum + README + LICENSE
```
**Qaydalar:** `cmd/` — hər app öz qovluğunda 1 main; `internal/` — yalnız bu app üçün (compiler xarici import-a qadağa qoyur); `pkg/` — paylaşılabilən kitabxanalar.

### 4. Package Management + Vendoring
**Niyə lazımdır:** mühitlərarası tutarlılıq, qlobal vəziyyətdən izolyasiya, reproducible build-lər (CI/CD etibarı).

**Vendoring:** dependency-lərin layihə daxilinə `/vendor` qovluğuna kopyalanması:
- İzolyasiya — qlobal dəyişikliklər təsir etmir
- Offline build mümkündür
- `go mod vendor` — Go Modules dövründə opt-in rejim

### 5. MVC pattern-i Go-da
**Model** — data strukturu + biznes məntiqi (`pkg/models/book.go`); **View** — təqdimat (`web/templates`); **Controller** — input→model→view körpüsü (`internal/handlers/bookHandler.go`).

**Faydaları:** concern ayrılığı (paralel inkişaf + unit test), modulluq, yenidən istifadə, asan maintenance.

**Sorğu axını:** request → controller → model (data) → controller → view (render) → user.

**Qeyd:** Go MVC framework-i deyil (Rails/Django kimi), amma net/http + üçüncü paketlərlə MVC arxitekturası qurulur.

### 6. Əsas Application obyekti (orchestrator)
**Kitabdan kod nümunəsi:**
```go
package main

import (
    "database/sql"
    "net/http"
    "github.com/gorilla/mux"
)

type App struct {
    Router *mux.Router
    DB     *sql.DB
    Config map[string]string
}

func (a *App) Initialize(config map[string]string) {
    connectionString := config["database"]
    var err error
    a.DB, err = sql.Open("mysql", connectionString)
    if err != nil { panic(err) }
    a.Router = mux.NewRouter()
    a.initializeRoutes()
}

func (a *App) initializeRoutes() {
    a.Router.HandleFunc("/books", a.getBooks).Methods("GET")
}

func (a *App) getBooks(w http.ResponseWriter, r *http.Request) { /* ... */ }

func (a *App) Run(addr string) { http.ListenAndServe(addr, a.Router) }

func main() {
    config := map[string]string{"database": "user:password@/dbname"}
    app := &App{}
    app.Initialize(config)
    app.Run(":8080")
}
```
**Sub-kod izahı:**
- `App struct` — Router + DB + Config bir mərkəzdə
- `Initialize` — DB bağlantısı + route qurulumu
- `a.getBooks` — metodu handler kimi ötürmək (pointer receiver → mux-yə method value)
- `Run` — router-i ListenAndServe-ə ötürür (default mux YOX, gorilla mux)

### 7. Konfiqurasiya və environment dəyişənləri
**Nədir:** Port, DB creds, API key-lər koddan XARİCDƏ — environment-variable və ya .env.

**Kitabdan kod nümunəsi:**
```go
type Configuration struct {
    ServerAddress, DbUser, DbPassword, DbHost, DbName string
}

func loadConfiguration() Configuration {
    err := godotenv.Load(".env")   // lokal dev üçün
    if err != nil {
        log.Println("no .env — production env vars gözlənilir")
    }
    return Configuration{
        ServerAddress: os.Getenv("SERVER_ADDRESS"),
        DbUser:        os.Getenv("DB_USER"),
        // ...
    }
}
```
**Praktik qaydalar (müəllifin tövsiyələri):**
1. `.env`-i heç vaxt git-ə commit ETMƏ (sensitive data)
2. Hierarxik konfiqurasiya: default → env var → CLI flag
3. Aydın adlar: `GITFORGITS_DB_USER` > `GFG_DB_U`

### 8. Dependency Injection (asılılıqların inyeksiyası)
**Nədir:** Komponent dependency-lərini ÖZÜ yaratmır — xaricdən "inyeksiya" olunur (constructor/method) → Inversion of Control (idarəetmənin tərsinə çevrilməsi).

**Faydaları:** decoupling (interface-lərə asılılıq), testability (mock/stub əvəzetməsi), flexibility, lifecycle idarəsi, scalability.

**Kitabdan kod nümunəsi (framework-süz DI):**
```go
type BookRepository interface {
    GetByID(id int) (*Book, error)
}

type SQLBookRepository struct{ db *sql.DB }
func (s *SQLBookRepository) GetByID(id int) (*Book, error) { /* SQL */ }

type BookService struct{ repo BookRepository }   // INTERFACE-ə asılı

func NewBookService(r BookRepository) *BookService {
    return &BookService{repo: r}               // constructor inyeksiyası
}

func main() {
    sqlRepo := &SQLBookRepository{db: db}
    bookService := NewBookService(sqlRepo)      // konkret → interface-ə
}
```
**Sub-kod izahı:** `BookService` konkret SQL reposundan deyil, `BookRepository` interface-indən asılıdır → test-də `MockBookRepository` inyeksiya et, production-da `SQLBookRepository`.

### 9. Health Check endpoint-ləri
**Sadə (alive göstəricisi):**
```go
func (a *App) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("OK"))
}
a.Router.HandleFunc("/healthcheck", a.healthCheckHandler).Methods("GET")
```
**Detallı (dependency yoxlaması):**
```go
func (a *App) detailedHealthCheckHandler(w http.ResponseWriter, r *http.Request) {
    if err := a.DB.Ping(); err != nil {          // DB bağlantısını sına
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte("Database Unreachable"))
        return
    }
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("OK"))
}
```
**Təhlükəsizlik:** Detallı health check sistem haqqında məlumat verir — yalnız autorizasiyalı sistemlərə açıq saxla (K8s liveness probe üçün sadə versiya kifayətdir).

## Əsas terminlər
- Go Module (Go modulu) — go.mod manifestli paket dəsti
- go.sum — kriptoqrafik checksum reyestri
- Module proxy (modul proksisi) — keshləyici mediator
- Vendoring — dependency-lərin layihəyə kopyalanması
- MVC — Model-View-Controller arxitekturası
- Dependency Injection (asılılıq inyeksiyası) — xarici dependency təminatı
- Inversion of Control — asılılıq yaradılışının tərsinə çevrilməsi
- Health check endpoint — sistem sağlamlıq yoxlama URL-i
- godotenv — .env faylından env var yükləyici

## Praktik nəticə
1. Yeni layihə: `go mod init` + `cmd/`, `internal/`, `pkg/` strukturu — internal xarici istifadədən compiler səviyyəsində qorunur.
2. App obyekti pattern: Router+DB+Config struct-da → Initialize → Run; handler-lər metod kimi.
3. Konfiqurasiya heç vaxt koddə olmasın: env var + .env (lokal) + git-də YOX.
4. Service-lər həmişə interface asılılığı ilə qur (repository pattern) — DI framework-süz, Go-nun təbii yolu.
5. Health check-lər: sadə versiya (liveness) + detallı versiya (dependency-lərlə, qorunan) — K8s/Docker üçün vacibdir.

## Mənbə
Pages: 47-79 (PDF səh. 47-79)
