# Глава 1 — İlk mikroservisin hazırlanması (User) (səh. 13-86)

## Bu fəsil nədən bəhs edir?

Tam mikroservis User-in (Account) sıfırdan qurulması: lokal mühit (VS Code +
Go extension, GVM, protoc + protoc-gen-go/protoc-gen-go-grpc, Git/git flow,
gofmt), Go-nun dependency tarixi ($GOPATH → godep → glide → govendor → Go
Modules — go.mod/go.sum, hash müdafiəsi Supply Chain Attack-a qarşı),
layihə strukturu (cmd/internal: service/model/repository/api/migrations/
config/pkg), Makefile automation, .env konfiqurasiyası, zerolog structured
logging + spew, GORM ORM (repository pattern), Swagger + Gin (Radix Tree),
PostgreSQL dizaynı (normal formalar 1НФ/2НФ/3НФ), Protobuf kontraktlar
(ayrı repo, cyclic import qadağası, Docker+Make ilə deterministik
generasiya), 3-qatlı arxitektura (transport → service → repository, hər
qatın ÖZ modeli + mapper-lər), goose miqrasiyaları, DI (internal/app
kompozitoru, əl ilə yazılmış DI konteyneri), Postman gRPC testi.

## Əsas fikirlər

### 1. Alət Zənciri və Go Modules
- **VS Code + Go Team extension:** IntelliSense, sintaksis diaqnostikası
- **GVM (Go Version Manager):** çoxversiyalı layihələr; versiya upgrade
  MƏRHƏLƏLİ (bir neçə versiyaya "tullanma" QADAĞA — asılılıqlar sınmalı)
- **Go Modules tarixi:** $GOPATH (tək kataloq, versiya fiksliyi YOX) →
  godep (Godeps.json) → glide (glide.yaml + glide.lock) → govendor
  (vendor.json) → **Go Modules (1.11+, 1.16-dan default)**:
  - `go.mod` — modul adı + Go versiyası + require(dependency@versiya)
  - `go.sum` — SHA-256 hash-lər (h1: prefiksi, base64); dəyişdirilmiş
    paket build zamanı aşkarlanır → **Supply Chain Attack** müdafiəsi
- **protoc:** `protoc --go_out=. --go-grpc_out=. file.proto` → struct +
  gRPC servis kodu avtomatik

### 2. Layihə Strukturunun Kanonu
```
account/
├── cmd/            ← main.go (yalnız START; package main + func main)
├── internal/       ← XARİCDƏN İMPORT OLUNA BİLMƏZ (Go kompilyator qəddiyyəti)
│   ├── app/       ← DI kompozitoru (app.go)
│   ├── service/   ← biznes məntiqi (infrastruktur haqqında HEÇ NƏ BİLMİR)
│   ├── model/     ← biznes modelləri
│   ├── repository/← DB çıxışı (+ öz model və mapper qatı)
│   ├── server/    ← transport (gRPC handler-lər, validasiya)
│   ├── migrations/← goose miqrasiyaları
│   ├── config/    ← env parsing
│   └── logger/
├── pkg/            ← yenidən istifadə edilə bilən kod
└── Makefile
```
- **internal xüsusiyyəti:** kənar layihələr import edə BİLMƏZ — daxili
  implementasiya qorunur

### 3. Makefile — Rutin Avtomatlaşdırması
```makefile
APP_NAME ?= account
build:
	go build $(LDFLAGS) -o $(BINARY_NAME) $(MAIN_PATH)
fmt:        ## gofmt -s -w, gofumpt, goimports
lint:       ## golangci-lint run
test:       ## go test -v ./...
test-coverage:  ## go test -coverprofile + cover -html
dev: fmt lint test build
ci:  deps fmt lint test build
```
- Versiya/build-time `-ldflags "-X main.Version=..."` ilə binariyaya

### 4. Konfiqurasiya: .env + env strukturu
```go
type Config struct {
    ServiceName string `env:"SERVICE_NAME" required:"true" default:"account-service"`
    AppEnv      string `env:"APP_ENV" required:"true" default:"development"`
    Host        string `env:"HTTP_HOST" default:"localhost"`
    Port        int    `env:"HTTP_PORT" default:"9000"`
}
// config.Load(): godotenv.Load() + env.Parse(cfg)
```
- `.env` → `.gitignore`-da; komandaya `.env.example` (şablon) paylaşılır
- Paketlər: `caarlos0/env/v10` (struktur parsing), `joho/godotenv` (.env faylı)
- DB: `DB_DSN=postgres://postgres:postgres@localhost:5432/account?sslmode=disable`

### 5. Strukturqlu Loglaşdırma — zerolog
```go
writer := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "2006-01-02 15:04:05"}
log := zerolog.New(writer).With().
    Timestamp().Str("source", cfg.ServiceName).Str("env", cfg.AppEnv).Logger()
```
- Səviyyələr rəngli: INFO-yaşıl, WARN-sarı, ERROR-qırmızı, DEBUG-mavi
- `spew` — struktur dump; konfiq başlanğıcda loglanır
- **Qayda:** bir sorğu = bir xəta logu (repository-də logladısa — service-də
  YOX)
- Məqsədlər: audit (təhlükəsizlik), monitoring (Prometheus/Loki), alert
  (504 xətaları normadan artıq → siqnal); JSON format production üçün

### 6. ORM: GORM + Repository Pattern
```go
// qoşulma
db, err := gorm.Open(postgres.Open(cfg.DbDsn), &gorm.Config{})
// repository
type Repository struct {
    db     *gorm.DB
    logger *zerolog.Logger
}
res := r.db.WithContext(ctx).
    Clauses(clause.OnConflict{UpdateAll: true}).
    Create(&userRepo)                     // upsert
r.db.WithContext(ctx).Where("id = ?", userID).First(&user)
r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&users)  // pagination
```
- **ORM trade-off:** sadə əməliyyatlar rahat; mürəkkəb sorğularda
  generasiya olunan SQL suboptimal ola bilər → **qarışıq yanaşma: ORM
  standart üçün, xam SQL optimallaşdırma üçün**
- Repository biznes məntiqini DB detallarından ayırır; bağlantı main-də,
  istifadə *gorm.DB kimi

### 7. Swagger + Gin
```go
// @title Account Service
// @version 1.0
router := gin.Default()
router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
// @Summary ... @Tags health @Success 200 {string} string "pong"
// @Router /ping [get]
swag init -g cmd/main.go
```
- Gin: Radix Tree router — çox marşrutda sürətli; minimalist
- Swagger (swaggo): kod şərhindən (@teg) sənədlər generasiyası — sənəd
  heç vaxt koddan geri qalmır

### 8. DB Dizaynı — Normal Formalar
- **1НФ:** hər cədvəl = entity; PK minimal sahə sayı (yoxdursa BIGSERIAL
  avtoinkrement VƏ YA UUID); sahələr ATOMAR (bir xanada "iPhone 14, 13,
  12" siyahısı QADAĞA); yazı sırası əhəmiyyətsiz (tarix sahəsi istifadə et)
- **2НФ:** 1НФ + hər sütun PK-dan Qeyri-azaldıla bilən asılı (yalnız
  firmanın endirimindən asılı sütun ayrı cədvələ — dekompozisiya)
- **ЗНФ:** 2НФ + transitive asılılıq YOX (Model→Mağaza→Telefon: Telefon
  mağazadan asılıdır, modeldən YOX → 2 cədvələ böl)

### 9. Protobuf Kontraktlar — Ayrı Repo
```protobuf
service Account {
  rpc CreateUser(CreateUserRequest) returns (google.protobuf.Empty) {
    option (google.api.http) = { post: "account/api/v1/create_user" body: "*" };
  }
  rpc GetUser(GetUserRequest) returns (GetUserResponse) { ... get: ".../users/{user_id}" }
  rpc GetUsers ... rpc DeleteUser ... rpc UpdateUser ...
}
message User { uint64 id = 1; string login = 2; ... google.protobuf.Timestamp created_at = 9; }
```
- **Niyə ayrı repo (contracts/account):** bir servisdə protudi paylanma
  ehtiyacı; servis+model proto faylları AYRI — **cyclic import protoc-da
  xəta verir** (modullar bir-birinə bağlıdırsa — modulluq POZULUB)
- **Pagination ayrı entity:** çoxservis istifadəsi
- **Deterministik generasiya:** protoc-u Docker konteynerə (golang:1.24 +
  protoc 27.1 + googleapis + pluginlər) → hər kəsdə eyni nəticə; Makefile:
  `docker run --rm -v ... proto-builder` döngüsü ilə hər .proto üçün
- **Semver ilə tag:** `git tag v1.0.0` → `go get` ilə istifadə; breaking
  olmadan yeni funksiya = v1.1.0

### 10. 3-Qatlı Arxitektura — Hər Qatın ÖZ Modeli
- **Niyə proto struct-larını hər yerdə istifadə ETMƏK OLMAZ:**
  1. gRPC→REST dəyişməsi 1 adapterdə həll olunmalı (yoxsa böyük refactor)
  2. Servis eyni anda gRPC + Kafka + JSON danışa bilər — daxili model
     universaldır
  3. Backwards compatibility: proto-dan sahə SİLİNMƏZ (reserved); daxili
     modeldə azadlıq
  4. Sahə dəyişikliyi cascade yox, LOKALIZASIYA olunmalıdır (kontrakt +
  mapper dəyişir — hamısı YOX)
- **Qatlar:** server (transport: PbToUserCreate → service) → service
  (biznes) → repository (RepoUserToUser) — mapper qatları arası
- **Service interfeyslə işləyir:** `type Repository interface {...}` —
  implementasiya paketinə istinad YOX; unit-test üçün STUB bəs edir,
  real DB lazım deyil

### 11. Goose Miqrasiyaları
```
goose -dir internal/migrations create initdb
goose.AddNamedMigrationContext("2025..._initdb.go", upCreateUsersTable, downCreateUsersTable)
```
```go
func upCreateUsersTable(ctx context.Context, tx *sql.Tx) error {
    _, err := tx.Exec(`CREATE TABLE IF NOT EXISTS users (
        id BIGSERIAL PRIMARY KEY, login TEXT NOT NULL UNIQUE, email TEXT NOT NULL UNIQUE,
        phone TEXT, first_name TEXT, last_name TEXT, middle_name TEXT, age INT,
        created_at TIMESTAMP NOT NULL DEFAULT NOW(), updated_at TIMESTAMP NOT NULL DEFAULT NOW());`)
}
func downCreateUsersTable(...) { tx.Exec(`DROP TABLE IF EXISTS users;`) }
```
- Ad = timestamp + ad (konflikt ehtimalı ~0); up = tətbiq, down = rollback;
  **goose miqrasiyanı tranzaksiyada icra edir** — xəta olarsa avtomatik
  rollback
- `goose.UpContext(ctx, dbGoose, "internal/migrations")` startda

### 12. DI — Əl ilə Yazılmış Konteyner (internal/app)
```go
type App struct { cfg; logger; accountRepository; accountService; accountServer; grpcServer }
func (a *App) Run(ctx context.Context) error {
    // getAccountServer → (lazımsa) getAccountService → getRepository → runMigrations
    // kaskad: yuxarıdan aşağı — transport → biznes → data
    go func() { serveErrCh <- a.grpcServer.Serve(lis) }()
    select {
    case <-ctx.Done():        a.grpcServer.GracefulStop(); return ctx.Err()
    case err := <-serveErrCh: ...
    }
}
```
- **Prinsip:** hər qat asılılıqları XARİCDƏN alır, özü YARATMIR;
  main.go yalnız config + logger + app.New + app.Run
- **Go fəlsəfəsi:** framework DI maqiyası YOX — kod AÇIQ və proqnozlaşdırılan
- GracefulStop + context — mehmanpərvər dayandırılma

### 13. Postman ilə gRPC Test
- Collection (microservices) + Environment (dev: `account =
  localhost:50051` — `{{account}}` şəklində istifadə)
- gRPC sorğu → Service definition → proto fayl (+ import KATALOGU — faylın
  ÖZÜ YOX); "Use Example Message" şablon doldurur
- Status 0-OK = uğur; limit/offset filtrlə sınaq

## Əsas terminlər
- Supply Chain Attack (təchizat zənciri hücumu) — paket dəyişdirilib,
  versiya eyni qalıb; go.sum hash-i tutur
- Surrogate key (sürroqat açar) — BIGSERIAL/UUID süni PK
- Atomicity (atomarlıq) — bir hüceyrə = bir dəyər (1НФ)
- Transitive dependency (keçid asılılığı) — A→B→C (ЗНФ pozuntusu)
- Upsert — insert + conflict update (clause.OnConflict{UpdateAll: true})
- Mapper — qatlararası model çevirici (PbToUser, UserToRepoUser)
- GracefulStop — aktiv sorğuları bitirib dayanma
- Cyclic import (dövrəvi import) — A→B→A; protoc üçün FATAL

## Praktik nəticə

1. **Layihə şablonu:** cmd/internal/pkg + Makefile + .env(.example) +
   .gitignore — ilk gündən; internal BLOKLANMIŞ import.
2. **Model ayrılığı:** hər qatın öz modeli; dəyişikliklər mapper-də
   lokallaşır — gRPC-dən asılı biznes kodu ANTİPATTERNDİR.
3. **Xəta logu bir dəfə:** repository logladısa, service susur.
4. **Miqrasiya tranzaksiyada:** goose up/down cütlüyü + timestamp adları.
5. **DI əl ilə:** internal/app kompozitor — Go-da bu, framework maqiyasından
   üstündür.
6. **Normal formalar:** praktikada 2НФ/ЗНФ kifayətdir — dublikat və keçid
   asılılıqlarını gündəlik dizaynda tanı.

## Mənbə
Pages: 13-86 (PDF 14-87)
