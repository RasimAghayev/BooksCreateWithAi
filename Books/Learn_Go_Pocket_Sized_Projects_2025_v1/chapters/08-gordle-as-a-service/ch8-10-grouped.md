# Chapters 8-10 — Gordle Service, Maze Solver, Habits gRPC (səh. 328-518)

## Bu fəsillər nədən bəhs edir?

(8) Gordle oyununu REST web servisə çevirmək: http.ServeMux routing,
handler-lar, domain/api ayırımı, in-memory repository, random ID təhlükəsizliyi.
(9) Paralel maze solver: PNG şəkil yükləmə, linked list path, kanyon
(branch) axiomları, goroutine-lərlə parallel eksplorasiya, select/quit kanalı,
WaitGroup, GIF animasiya. (10) gRPC habits tracker: protobuf, kod generasiyası,
typed errors, minimock, integration test, context.

## Əsas fikirlər

### 1. REST servis — strukturlaşdırma (Ch8)
```
httpgordle/
├── api/           ← API tipləri (JSON strukturları, marşrutlar)
├── internal/
│   ├── handlers/  ← HTTP handler-lər (endpoint başına papka)
│   └── gordle/    ← domain (oyun məntiqi)
└── main.go
```
```go
r := http.NewServeMux()
r.HandleFunc(http.MethodPost+" "+api.NewGameRoute, newgame.Handler(db))
r.HandleFunc(http.MethodGet+" "+api.GetGameRoute, getstatus.Handler(db))
r.HandleFunc(http.MethodPut+" "+api.GuessRoute, guess.Handler(db))
http.ListenAndServe(":8080", r)
```
- **Servis vs server vs endpoint vs handler** ayırımı; hexagonal (ports and
  adapters) mindset — domain ≠ API detal
- Handler imzası: `func(http.ResponseWriter, *http.Request)`; istənilən
  funksiya mux-a qeydiyyat oluna bilir (closure ilə db injekt)

### 2. JSON cavab yazmaq
```go
w.Header().Set("Content-Type", "application/json")
w.WriteHeader(http.StatusCreated)
json.NewEncoder(w).Encode(apiGame)
```
- `WriteHeader` Write-dən ƏVVƏL — Write default 200 yazır

### 3. Random GameID (təhlükəsizlik)
- İncremental ID = təhlükəsizlik deşiyi (başqasının oyununu tapmaq asandır)
- `crypto/rand` ilə UUID-yə bənzər təsadüfi ID

### 4. In-memory repository + mutex
```go
type GameRepository struct {
    storage map[session.GameID]session.Game
    mu      sync.Mutex      // bir neçə request eyni anda yazır
}
```
- Repository xarici modullara QAPALI olmalı (internal/) — yoxsa API-nin
  mənası yoxdur; handler-lər stub ilə test olunur:
```go
type gameAdderStub struct{ err error }
func (g gameAdderStub) Add(_ session.Game) error { return g.err }
```

### 5. HTTP statuslar + xəta idarəsi
```go
game, err := db.Find(session.GameID(id))
if err != nil {
    if errors.Is(err, repository.ErrNotFound) {
        http.Error(w, "this game does not exist", http.StatusNotFound)
        return
    }
    http.Error(w, "...", http.StatusInternalServerError)
}
```
- Handler XƏTANI QAYTARMIR — status kodu + log; daxili xəta istifadəçiyə
  damır (info leak qarşısı)

### 6. Maze solver — paralel eksplorasiya (Ch9)
**Struktur:**
```go
type path struct {
    previousStep *path      // linked list — geri dönüş yoxdur
    at           image.Point
}
type Solver struct {
    maze           *image.RGBA
    pathsToExplore chan *path  // yeni branch-lər buraya daxil olur
    quit           chan struct{} // xəzinə tapıldı siqnalı
    solution       *path
}
```
**Əsas dövrü:**
```go
// listener: hər branch üçün YENİ goroutine
for {
    select {
    case <-s.quit:   // xəzinə tapıldı → dayan
        return
    case p := <-s.pathsToExplore:
        wg.Add(1)
        go func(p *path) {
            defer wg.Done()
            s.explore(p)
        }(p)
    }
}
// explorer daxilində: yeni branch tapanda YA yeni goroutine-a başlanğıc
// verilir, YA select quit-i yoxlayır
```
- **Fəlsəfə:** Theseus (xəzinə axtaran) + Daedalus (dinləyici) goroutine-ləri;
  hər ayrılma nöqtəsi (intersection) = yeni goroutine
- PNG: `image/png.Decode(reader)` → `*image.RGBA` type assertion; sərhəd
  kənarı `RGBAAt` zero value qaytarır — təhlükəsiz

### 7. Loop dəyişəni tələsi (Go ≤1.22!)
```go
for p := range s.pathsToExplore {
    p := p             // KÖLGƏ — closure-a safe copy
    wg.Add(1)
    go func() { s.explore(p) }()
}
```
- Go 1.22-dən əvvəl loop dəyişəni YENİDƏN İSTİFADƏ olunur — paralel
  goroutine-lərdə hamısı SON dəyəri görür. `p := p` klassik müdafiə.

### 8. GIF animasiya
- Explored pixel-lər kanaldan async toplanır → hər ~N pikselə 1 frame
  (30 frame hədəfi) → `gif.GIF` + `image/color/palette` (paket adı toqquşması:
  alias `plt "image/color/palette"`)

### 9. gRPC + Protobuf (Ch10)
**Proto faylı:**
```protobuf
syntax = "proto3";
package habits;
option go_package = "learngo-pockets/habits/api";

message CreateHabitRequest {
  string name = 1;
  optional uint32 weekly_frequency = 2;
}
service Habits {
  rpc CreateHabit(CreateHabitRequest) returns (CreateHabitResponse);
}
```
**Generasiya:**
```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
protoc -I=api/proto/ --go_out=api/ --go_opt=paths=source_relative \
       --go-grpc_out=api/ --go-grpc_opt=paths=source_relative api/proto/*.proto
```
- `//go:generate` direktivi — `go generate ./...` ilə protoc çağırışı
- Request/Response hər endpoint üçün AYRI (version compatibility: sahə
  əlavə etmək mövcud istifadəçiləri pozmur)

### 10. gRPC server implementasiyası
```go
type Server struct {
    db  Repository   // bizim interfeysimiz
    lgr Logger
}

func (s *Server) CreateHabit(ctx context.Context, req *api.CreateHabitRequest) (*api.CreateHabitResponse, error) {
    h := habit.Habit{Name: habit.Name(req.GetName())}
    h, err := habit.Create(ctx, s.db, h)
    if err != nil {
        var invalidErr habit.InvalidInputError
        if errors.As(err, &invalidErr) {
            return nil, status.Errorf(codes.InvalidArgument, "%s", err)
        }
        return nil, status.Errorf(codes.Internal, "%s", err)
    }
    return &api.CreateHabitResponse{Habit: habitToAPI(h)}, nil
}
```
- **Typed error → status kodu:** `errors.As(err, &invalidErr)` → codes.InvalidArgument
  (3); qalan → codes.Internal (13)
- **Domain ≠ proto:** biznes tipləri (`habit.Habit`) ayrı, api adapterləri
  çevirir — protokol dəyişəndə domain toxunulmaz

### 11. Mock + integration test
```go
ctrl := minimock.NewController(t)
defer ctrl.Finish()
db := tt.db(ctrl)         // minimock — typed parameter mock
got, err := habit.Create(ctx, db, h)
assert.ErrorIs(t, err, tt.expectedErr)
```
- Integration test: real gRPC server (random port) + real client — scenario
  list (create → list → tick → status)
- `go test -short` → `testing.Short()` ilə AĞIR testlər skip (CI sürəti)

### 12. Context
```go
ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
defer cancel()
```
- 4 metod: Deadline, Done, Err, Value; `Background()`/`TODO()` — kök;
  hər uzaq çağırış ÖZ deadline-i; `ctx.Done()` kanalı — standard kitabxanada
  60+ yerdə `ctx.Err()` ilə birlikdə
- **Value QADAĞAN-dır** parameter ötürülməsi üçün (yalnız son çarə) — reqvest-scoped
  dəyərlər protokol öz Context-i ilə gəlir

## Əsas terminlər

- http.ServeMux (router)
- Hexagonal Architecture (ports/adapters)
- Repository pattern
- Linked List (yol strukturu)
- select (çoxkanal gözləmə)
- Loop variable shadowing (p := p)
- Protobuf / gRPC
- go:generate
- Typed Error + errors.As
- status codes (gRPC)
- context.Context
- Integration test / testing.Short

## Praktik nəticə

- Domain-paketdə heç vaxt proto/http tipləri — adapter çevirir
- Handler: xəta → status kod; daxili detal istifadəçiyə YOX
- Paralel eksplorasiyada: select + quit kanalı + WaitGroup — goroutine sızması
- Go <1.22-də loop dəyişənini kölgələ (p := p)
- gRPC: Request/Response ayrı; typed error → status code

## Mənbə

Pages: 328-518 (Chapters 8-10, Learn Go with Pocket-Sized Projects)
