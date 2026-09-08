# Chapter 8 — Gordle as a service (səh. 328-392)

## Bu chapter nədən bəhs edir?

Gordle oyunu REST servisinə çevrilir: `net/http` server, ServeMux routing (method +
path pattern), status kodları, GET/POST/PUT/DELETE endpoint-lər, hexagonal
architecture (api/handlers/internal), in-memory repository, stub testlər və
təhlükəsizlik (random ID, rate limiting).

## Əsas fikirlər

### 1. Server, service, endpoint anlayışları
- **Server** — portu dinləyən kompüter/proses; **service** — tək məsuliyyətli
  funksionallıq; **web service** — HTTP üzərindən istifadə olunan service
- **Endpoint** = yol + metod; **HTTP handler** — sorğunu qarşılayan funksiya

İlk server:
```go
func main() {
    err := http.ListenAndServe(":8080", handlers.NewRouter())
    if err != nil { panic(err) }
}
```
`ListenAndServe` **heç vaxt qayıtmır** (bloklanır) — Ctrl-C ilə dayandırılır.

### 2. REST və idempotentlik
| Metod | Mənası | İdempotent? |
|---|---|---|
| GET | oxu | Bəli |
| POST | yarat | Xeyr (hər çağırış yeni resurs) |
| PUT | yenilə | Bəli |
| DELETE | sil | Bəli |

### 3. Routing — ServeMux
```go
func NewRouter(db *repository.GameRepository) *http.ServeMux {
    r := http.NewServeMux()
    r.HandleFunc(http.MethodPost+" "+api.NewGameRoute,  newgame.Handler(db))
    r.HandleFunc(http.MethodGet+" "+api.GetStatusRoute, getstatus.Handler(db))
    r.HandleFunc(http.MethodPut+" "+api.GuessRoute,      guess.Handler(db))
    return r
}
```
**Sub-kod izahı:**
- `"POST /games"` formatı → metod + yol pattern-i birlikdə; metod yoxlaması
  avtomatik olur
- `http.HandlerFunc` imzası: `func(w http.ResponseWriter, r *http.Request)`

### 4. Status kodları
```go
w.WriteHeader(http.StatusCreated) // 201 — Write-dan ƏVVƏL çağırılmalı
_, _ = w.Write(encoded)
```
| Kod | Məna |
|---|---|
| 200 | OK |
| 201 | Created (POST uğurlu) |
| 400 | Bad Request (malformed sorğu) |
| 404 | Not Found |
| 405 | Method Not Allowed |
| 500 | Internal Server Error |

Default 200-dür — müxtəlif cavab üçün Write-dan əvvəl yaz.

### 5. Endpoint parametrlərinin 4 növü
1. **Path** — `/game/{id}` — resursu hədəfləyir
2. **Query** — `?key=value` — filtrlər, optional
3. **Body** — POST/PUT məzmunu (JSON)
4. **Meta** — header-lər (auth, content-type)

### 6. Layihə strukturu (hexagonal-ə yaxın)
```
httpgordle/
├── api/                  → route sabitləri (publik dəyərlər)
├── internal/
│   ├── handlers/         → HTTP adapteri (newgame, getstatus, guess)
│   │   ├── newgame/
│   │   ├── getstatus/
│   │   └── guess/
│   ├── session/          → domain (GameID, Status, GameResponse)
│   ├── repository/       → in-memory DB
│   └── gordle/           → oyun məntiqi (kitabxana)
└── main.go
```
- `internal/` → Go kompilyatoru xarici modullara import-u QADAĞAN edir
- Domain (session) HTTP haqqında heç nə bilmir — adapterlər çevirir
- `pkg/` qovluğu (görüldüyü kimi) — xarici istifadəyə açıq kitabxanalar üçün

### 7. Handler → domain adapteri
```go
func Handle(w http.ResponseWriter, req *http.Request) {
    game, err := createGame(req.Context())
    if err != nil {
        log.Printf("unable to create game: %v", err)
        http.Error(w, "unable to create game", http.StatusInternalServerError)
        return
    }
    // session.Game → GameResponse (API modelinə çevir)
}
```
Tərcümə funksiyaları (`GameResponse`) — daxili tipi API formatından ayırır
(anti-corruption layer).

### 8. In-memory repository
```go
type GameRepository struct {
    db   *repository.MemoryGameRepository // əsl implementasiya
    // və ya:
    games map[session.GameID]session.Game
    mu    sync.Mutex                      // paralel sorğulara qarşı
}

func (r *GameRepository) Add(game session.Game) error
func (r *GameRepository) Find(id session.GameID) (session.Game, error) // ErrNotFound
func (r *GameRepository) Update(game session.Game) error
```
- Servis tək insan üçün deyil — **eyni anda çox sorğu** → mutex şərt
- `ErrNotFound` sentinel xətası → handler-də `errors.Is` ilə 404-ə çevrilir:
```go
if errors.Is(err, repository.ErrNotFound) {
    http.Error(w, "game not found", http.StatusNotFound)
    return
}
```

### 9. Stub testlər
```go
type gameAdderStub struct{ err error }
func (s gameAdderStub) Add(game session.Game) error { return s.err }

handleFunc := Handler(gameAdderStub{}) // stub inject olunur
handleFunc(recorder, req)              // httptest.ResponseRecorder
```
- Hər handler kiçik interfeysə (yalnız ehtiyacı olan metodlara) bağlıdır →
  stub-lama asan
- `testify` (require/assert) + `assert.JSONEq` → JSON müqayisəsi

### 10. Random ID — təhlükəsizlik
İncremental ID-lər **təhlükəsizlik deşiyidir** (hər kəs növbəti ID-ni
təxmin edir). `crypto/rand` ilə random UUID benzeri identifikator.

### 11. Təhlükəsizlik tədbirləri (sadalanan)
- Request sayının məhdudlaşdırılması (rate limiting)
- İstifadəçi autentifikasiyası — ayrıca auth servisindən signed token, header
  ilə validasiya
- Loglama + xəta formatı uniform
- Query param decode (`r.URL.Query()` → `url.Values`)

## Əsas terminlər

- ServeMux (marşrutlaşdırıcı)
- Endpoint / Handler
- Idempotency (idempotentlik)
- Hexagonal Architecture (altıbucaqlı arxitektura)
- Anti-Corruption Layer (korupsiya əleyhinə qat)
- In-Memory Repository
- Stub Test
- Rate Limiting (sorğu limiti)

## Praktik nəticə

- `internal/` + kiçik paketlər → sərhədlər məcburi; domain HTTP-dən təcrid
- HandleFunc pattern-də metod+yol → metod yoxlaması avtomatik
- Status kodu Write-dan ƏVVƏL; xətaları sentinel + errors.Is ilə HTTP-ə xəritələ
- Handler asılılıqlarını mini-interfeyslərə bağla → stub test sadələşir
- İD-lər random (crypto/rand) — sequential ID security flaw-dur

## Mənbə

Pages: 328-392 (Chapter 8, Learn Go with Pocket-Sized Projects)
