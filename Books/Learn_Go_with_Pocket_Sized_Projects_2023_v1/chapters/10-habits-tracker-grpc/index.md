# Chapter 10 — Habits tracker using gRPC (səh. 452-518)

## Bu chapter nədən bəhs edir?

gRPC + Protobuf ilə vərdiş izləyici servis: proto tərifindən kod generasiyası
(`protoc`, `go generate`), server strukturu, business/api qat ayrımı, typed
xətalar → status kodları, minimock ilə mock testlər, inteqrasiya testləri
(buffers + real listener), `context` dərinləşməsi və ISOWeek ilə tick saxlama.

## Əsas fikirlər

### 1. Protobuf — serializasiya formatı
JSON/XML-dən daha kiçik və sürətli; **sxem əsaslı** (kontrakt birincildir).
```
api/proto/habit.proto:
message Habit {
  string id = 1;          // field NÖMRƏSİ = beynəlxalq identifikator
  string name = 2;
  int32 weekly_frequency = 3;
}
service Habits {
  rpc CreateHabit(CreateHabitRequest) returns (CreateHabitResponse);
  rpc ListHabits(...) returns (...);
}
```
- Field nömrələri dəyişməz qalır → əlavə sahə köhnə klientləri POZMUR
  (backward compatibility)
- `go_package` option → generasiya hədəfi

### 2. protoc + go generate
```bash
protoc -I=api/proto/ --go_out=api/ --go_opt=paths=source_relative \
  --go-grpc_out=api/ --go-grpc_opt=paths=source_relative \
  api/proto/*.proto
```
Kod icrasına bağla:
```go
//go:generate bash -c "protoc ... (yuxarıdakı komanda)"
```
→ `go generate ./...` — bütün generasiya bir əmrlə.

### 3. Server strukturu
```go
type Server struct {
    db Repository // in-memory repo (internal/)
    lgr Logger    // öz Logger interfeysi — "log.Print"-dan təcrid
}

grpcServer := grpc.NewServer()
api.RegisterHabitsServer(grpcServer, &server.Server{...})
lis, _ := net.Listen("tcp", fmt.Sprintf(":%d", port))
grpcServer.Serve(lis)
```

### 4. Business vs API qatları
- **Business (internal/habit):** `type ID string`, `type Habit struct` — protokol
  tərifləri domain-ə SIZMIR ("don't leak protocol definitions into core code")
- **API layer:** proto tiplərini domain-ə çevirir:
```go
func (s *Server) CreateHabit(ctx context.Context, req *api.CreateHabitRequest) (*api.CreateHabitResponse, error) {
    h := habit.Habit{Name: req.GetName(), ...}
    created, err := habit.Create(ctx, h, s.db)
    // xəta → status kodu çevir
}
```

### 5. gRPC status kodları + typed xətalar
| Kod | Məna |
|---|---|
| 0 | OK |
| 3 | InvalidArgument |
| 5 | Not Found |
| 13 | Internal |

```go
type InvalidInputError struct{ field string }
func (e InvalidInputError) Error() string { return "invalid " + e.field }

// API layer-də:
if err := habit.Create(...); err != nil {
    var inv InvalidInputError
    if errors.As(err, &inv) {
        return nil, status.Error(codes.InvalidArgument, err.Error()) // kod 3
    }
    return nil, status.Error(codes.Internal, "internal problem") // 5xx analoqu
}
```

### 6. Mock testlər (minimock)
```go
ctrl := minimock.NewController(t)
defer ctrl.Finish()
db := NewRepositoryMock(ctrl)
db.AddMock.When(ctx, mock.Anything).Then() // gözlənti
// və ya sadə əl ilə:
type MockList struct{ Items []habit.Habit; Err error }
func (m MockList) FindAll(_ context.Context) ([]habit.Habit, error) { return m.Items, m.Err }
```
Mock alətləri: mockgen / mockify / **minimock**. Interface = bir metod
(`habitCreator { Add(...) }`) → mock minimal.

### 7. Integration test — real gRPC
```go
func TestIntegration(t *testing.T) {
    grpcServ := newServer(t)
    listener, err := net.Listen("tcp", "") // random free port
    go grpcServ.Serve(listener)

    conn, err := grpc.Dial(listener.Addr().String(),
        grpc.WithTransportCredentials(insecure.NewCredentials()))
    habitsCli := api.NewHabitsClient(conn)
    // real çağırışlar: add → list → tick → status
}
```
`testing.Short()` → `go test -short` inteqrasiya testlərini skip edir.

### 8. context dərin baxış
```go
ctx, cancel := context.WithTimeout(parent, timeout)
defer cancel()
```
- `Background()` (kök) / `TODO()` (yerdə doldurulacaq)
- `ctx.Done()` ← deadline/cancel; `ctx.Err()` səbəb — std kitabxanada 60+ yerdə
  bu pattern
- **Hər uzaq çağırış öz deadline-li context** almalıdır
- `context.WithValue` — son çarə (yalnız retriv мəlumat üçün interfeys)

### 9. Tick + ISOWeek
```go
type ISOWeek struct {
    Year int
    Week int
}
// time.Time.ISOWeek() → həftəlik hesabat üçün stabili açar
```
AddTick(ctx, habitID, time) — habit yoxdursa inconsistensiya yaratma (corner
case: xəzinəyə yazmazdan əvvəl habit mövcudluğu yoxlanılır).

### 10. grpcurl — CLI klient
```bash
grpcurl -import-path api/proto/ -proto service.proto -plaintext \
  -d '{"habit_id":"96b7..."}' localhost:28710 habits.Habits/GetHabitStatus
```
HTTP-dən fərqli: binary format, `-d` JSON eyni vaxtda protobuf-a çevrilir.

## Əsas terminlər

- Protocol Buffers (protobuf)
- gRPC / rpc Service Definition
- Field Numbering (sahə nömrələnməsi)
- go:generate Direktivi
- Status Codes (gRPC status kodları)
- Typed Error (tipli xəta)
- minimock / Mock Controller
- Integration Test (inteqrasiya testi)
- context.Context / Deadline
- ISO Week (ISO həftəsi)

## Praktik nəticə

- Kontrakt əvvəl: proto faylı → `go generate` → server implementasiyası
- Proto tipləri domain-ə SIZDIRMA — çevirici API layer-də
- Xətalar: typed domain xətası → `status.Error(codes.X)` xəritələnməsi
- Unit mock + `-short`-skip inteqrasiya testləri — CI sürəti qorunur
- Hər RPC çağırışında deadline-li context

## Mənbə

Pages: 452-518 (Chapter 10, Learn Go with Pocket-Sized Projects)
