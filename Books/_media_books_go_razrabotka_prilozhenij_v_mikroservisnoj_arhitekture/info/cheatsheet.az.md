# Go mikroservisləri (Popova 2026) — Cheat Sheet (Azərbaycanca)

## 1. Layihə strukturu (hər Go mikroservisi üçün şablon)

```
service/
├── cmd/main.go            ← YALNIZ start: config+logger+app.New+app.Run
├── internal/              ← XARİCDƏN İMPORT OLUNA BİLMƏZ
│   ├── app/app.go         ← DI kompozitoru (əl ilə)
│   ├── config/            ← env parsing (caarlos0/env + godotenv)
│   ├── logger/            ← zerolog setup
│   ├── model/             ← biznes modelləri
│   ├── mapper/            ← qatlararası çeviricilər
│   ├── repository/        ← GORM + öz model və mapper
│   │   ├── model/  mapper/ migrations/
│   ├── service/           ← biznes məntiqi (interfeyslərlə!)
│   ├── server/            ← gRPC handler-lər
│   └── kafka/ account/ auth/  ← xarici inteqrasiyalar (wrap)
├── pkg/                   ← paylana bilən kod
├── Dockerfile  .dockerignore
└── Makefile
```

## 2. Konfiqurasiya

```go
type Config struct {
    ServiceName string `env:"SERVICE_NAME" required:"true" default:"auth-service"`
    Port        int    `env:"GRPC_PORT" default:"50052"`
    DbDsn       string `env:"DB_DSN" required:"true"`
    JwtSecret   string `env:"JWT_SECRET" required:"true"`
    KafkaBrokers []string `env:"KAFKA_BROKER_HOST" envSeparator:","`
}
func Load() (*Config, error) {
    godotenv.Load()                       // .env (gitignore-da!)
    cfg := &Config{}
    env.Parse(cfg)
}
```
- `.env.example` paylaşılır; `.env` hər kəsdə öz.

## 3. GORM əsasları

```go
db, _ := gorm.Open(postgres.Open(cfg.DbDsn), &gorm.Config{})
// Upsert:
db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Create(&user)
// Axtarış:
db.WithContext(ctx).Model(&User{}).Where("id = ? AND is_deleted = ?", id, false).First(&user)
// Sorğu qurucu + filtrlər + pagination:
q := db.WithContext(ctx).Model(&Transaction{})
if params.UserID != nil { q = q.Where("user_id = ?", *params.UserID) }
q.Offset(offset).Limit(limit).Order("created_at DESC").Find(&items)
// Tranzaksiya:
err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error { ... })
```

## 4. goose miqrasiyası

```
goose -dir internal/migrations create initdb
```
```go
func init() { goose.AddNamedMigrationContext("2025..._initdb.go", up, down) }
func up(ctx context.Context, tx *sql.Tx) error {
    _, err := tx.Exec(`CREATE TABLE IF NOT EXISTS users (...)`)
}
// Rollback: TƏRS sıra — indekslər → PK → cədvəllər
```

## 5. gRPC kontrakt + generasiya

```protobuf
service Account {
  rpc GetUser(GetUserRequest) returns (GetUserResponse) {
    option (google.api.http) = { get: "/api/v1/users/{user_id}" };
  }
}
```
- Ayrı repo (contracts/account); servis+model faylları AYRI (cyclic!);
  Docker proto-builder + Makefile gen; `git tag v1.1.0` → go get.

## 6. JWT Auth servisi

```go
hash, _ := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)  // Register
bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(pw))    // Login
claims := jwt.MapClaims{"sub": userID, "exp": accessExp.Unix()}         // issue
access := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// Refresh: hex(sha256(userID:unixnano:secret)) → DB (expires_at, revoked_at)
// Logout: Update("revoked_at", gorm.Expr("NOW()"))
```

## 7. Gateway JWT interceptor

```go
var publicMethods = map[string]bool{"/gateway.Gateway/Login": true, ...}
func (i *JWTinterceptor) UnaryInterceptor() grpc.UnaryServerInterceptor {
    return func(ctx, req, info, handler) (interface{}, error) {
        if publicMethods[info.FullMethod] { return handler(ctx, req) }
        userID := i.extractUserIDFromJWT(ctx)      // metadata → Bearer → parse
        return handler(context.WithValue(ctx, UserIDKey, userID), req)
    }
}
// Handler-də: interceptor.GetUserIDFromContext(ctx)
// Alg yoxlaması ZƏRURİ:
if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok { ... }
```

## 8. Saga (Kafka) axını

```
Transaction.Deposit:
  1. repo.Deposit() → DB tranzaksiya → pending status
  2. kafka.Publish("transaction_data", key=userID, {request_type, user_id,
     amount(qəpik), operation_id, timestamp})
Account.HandleTransaction (abunə "transaction_data"):
  3. switch request_type → Deposit/Withdraw/Transfer
  4. Publish("transaction_response", {operation_id, result})
Transaction.HandleAccountResponse (abunə "transaction_response"):
  5. result ? completed : failed → UpdateTransactionStatus
```
- kafka-go: Writer (LeastBytes, RequiredAcks: RequireOne) / Reader (GroupID,
  MinBytes/MaxBytes, CommitInterval); Subscribe → hər topic üçün unique
  groupID + goroutine.

## 9. Tranzaksiya nəzəriyyəsi cədvəlləri

**Anomaliya → izolyasiya həlli:**
| Anomaliya | Həll edən səviyyə |
|---|---|
| Dirty read | Read Committed |
| Lost update, non-repeatable | Repeatable Read |
| Phantom | Serializable |

**İzolyasiya seçimi:** default Read Committed (PG); maliyyə kritik
transferlər üçün Serializable VƏ YA tranzaksiya + status axını.

## 10. Normal formalar (praktik xatırlatma)

- 1НФ: atomar dəyərlər, PK
- 2НФ: PK-dan tam asılılıq
- 3НФ: transitive asılılıq yox (A→B→C split)
- НФБК: determinant = açar
- 4НФ+: müstəqil çoxdəyərli asılılıqları ayır — praktikada 3-yə qədər.

## 11. Multi-stage Dockerfile

```dockerfile
FROM golang:1.24.5-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./      # ƏVVƏLCƏ — layer cache!
RUN go mod download
COPY . .
RUN go build -o app ./cmd
FROM alpine:3.18
WORKDIR /app
COPY --from=build /app/app
COPY --from=build /app/internal/migrations ./internal/migrations
CMD ["./app"]
```

## 12. Docker-Compose (bütün sistem)

```yaml
services:
  account:
    image: account:latest
    build: ./account
    ports: ["50051:8080"]
    environment:
      DB_DSN: postgres://postgres:postgres@account_db:5432/account?sslmode=disable
      KAFKA_BROKER_HOST: kafka:9092        # host = SERVİS ADI
    depends_on: [account_db, kafka]
    networks: [app-network]
```

## 13. GitHub Actions CI

```yaml
name: Transaction
on: [workflow_dispatch, push]
jobs:
  stage_dev:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3.5.3
    - run: go mod download
    - run: go test -v -race -coverprofile=coverage.out ./...
```

## 14. Unit-test (gomock+testify)

```go
tests := []struct{ name string; setupMocks func(*mocks.MockRepository, ...); expectedError bool }{
    {"successful deposit", ..., false},
    {"repository error", ..., true},
    {"kafka publish error", ..., true},
}
for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        ctrl := gomock.NewController(t); defer ctrl.Finish()
        mockRepo.EXPECT().Deposit(gomock.Any(), params).Return(res, nil)
        logger := zerolog.Nop()
        svc := New(mockRepo, mockAcc, mockKafka, &logger)
    })
}
```

## 15. Kubernetes işə düşmə

```bash
kompose convert -f docker-compose.yaml   # hər servis → deployment+service YAML
kubectl apply -f                          # docker-compose.yaml-ı SİL (manifest sanır!)
kubectl get pods                          # status
kubectl logs <pod>                       # loglar
```

## 16. Səhv və tələsiklər (kitabın dərsləri)

| Tələsik | Düzgünü |
|---|---|
| localhost DB_DSN compose-da | servis adı (account_db) |
| golang:latest image | pin-lənmiş versiya (1.24.5-alpine) |
| proto-da sahə silmək | reserved elan et (backwards compat) |
| parolu Account-a göndərmək | yalnız Auth-da bcrypt |
| float pul | int64 qəpik |
| distributed tranzaksiya | Saga (mümkünsə ümumiyyətlə qaçın) |
| kubectl run | manifest + apply |
| statusu tranzaksiya daxilində təyin et | ayrı UpdateTransactionStatus (uzaq cavabdan sonra) |
