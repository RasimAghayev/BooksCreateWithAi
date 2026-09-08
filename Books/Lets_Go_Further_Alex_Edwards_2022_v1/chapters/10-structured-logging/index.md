# Chapter 10 — Structured Logging and Error Handling

## Bu fəsil nədən bəhs edir?

Custom JSON logger (level-based, mutex-li, io.Writer uyğun), http.Server
ErrorLog inteqrasiyası və panic recovery middleware-i.

## Əsas fikirlər

### 1. Niyə structured logging?
Sada `log.Logger` sətirləri axtırma/filtrləmə/monitorinq üçün struktur vermir.
JSON loglar: level, time, message, properties, trace — 3rd-party analiz
sistemlərinə (ELK, Loki və s.) birbaşa uyğun.

**Hədəf format:**
```json
{"level":"INFO","time":"2020-12-16T10:53:35Z","message":"starting server","properties":{"addr":":4000","env":"development"}}
```

### 2. jsonlog.Logger — custom logger
**Kitabdan kod nümunəsi (tam əsas hissələr):**
```go
// internal/jsonlog/jsonlog.go
type Level int8

const (
    LevelInfo  Level = iota // 0
    LevelError              // 1
    LevelFatal              // 2
    LevelOff                // 3
)

func (l Level) String() string { /* "INFO"/"ERROR"/"FATAL" */ }

type Logger struct {
    out      io.Writer
    minLevel Level
    mu       sync.Mutex
}

func New(out io.Writer, minLevel Level) *Logger { ... }

func (l *Logger) PrintInfo(message string, properties map[string]string) {
    l.print(LevelInfo, message, properties)
}
func (l *Logger) PrintError(err error, properties map[string]string) {
    l.print(LevelError, err.Error(), properties)
}
func (l *Logger) PrintFatal(err error, properties map[string]string) {
    l.print(LevelFatal, err.Error(), properties)
    os.Exit(1)
}

func (l *Logger) print(level Level, message string, properties map[string]string) (int, error) {
    if level < l.minLevel {
        return 0, nil // minimum severity altında → yazmır
    }
    aux := struct {
        Level      string            `json:"level"`
        Time       string            `json:"time"`
        Message    string            `json:"message"`
        Properties map[string]string `json:"properties,omitempty"`
        Trace      string            `json:"trace,omitempty"`
    }{
        Level:      level.String(),
        Time:       time.Now().UTC().Format(time.RFC3339),
        Message:    message,
        Properties: properties,
    }
    if level >= LevelError {
        aux.Trace = string(debug.Stack()) // ERROR/FATAL-da stack trace
    }
    line, err := json.Marshal(aux)
    if err != nil {
        line = []byte(LevelError.String() + ": unable to marshal log message: " + err.Error())
    }
    l.mu.Lock()
    defer l.mu.Unlock()
    return l.out.Write(append(line, '\n'))
}

// io.Writer interfeysinin təminatı — http.Server ErrorLog üçün
func (l *Logger) Write(message []byte) (n int, err error) {
    return l.print(LevelError, string(message), nil)
}
```

**Sub-kod izahı:**
- `iota` → ardıcıl konstant dəyərləri (0,1,2...); daha şiddətli səviyyə =
  daha böyük rəqəm → `level < l.minLevel` müqayisəsi işləyir
- `sync.Mutex` → paralel yazmaları seriyalaşdırır — mutexsiz iki goroutine-in
  yazması bir sətirdə qarışıq (intermingled) çıxardı
- `debug.Stack()` → ERROR+ səviyyələrində trace avtomatik əlavə olunur
- `os.Exit(1)` → FATAL semantic: logla + prosesi bitir
- `Write()` → Logger-i `io.Writer` edir — stdlib `log.Logger`-in hədəfi
  ola bilər

**İstifadə:**
```go
logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)
logger.PrintInfo("database connection pool established", nil)
logger.PrintInfo("starting server", map[string]string{
    "addr": srv.Addr,
    "env":  cfg.env,
})
logger.PrintFatal(err, nil)
```

**logError helper-inin təkmilləşdirilməsi:**
```go
func (app *application) logError(r *http.Request, err error) {
    app.logger.PrintError(err, map[string]string{
        "request_method": r.Method,
        "request_url":    r.URL.String(),
    })
}
```
- Artıq hər xəta request konteksti ilə loglanır

### 3. http.Server ErrorLog inteqrasiyası
**Problem:** `http.Server` öz mesajlarını (unrecovered panics, connection
xətaları) stderr-ə plain-text yazır — JSON formatından kənar düşür.

**Həll — adapter pattern:**
```go
srv := &http.Server{
    Addr:     fmt.Sprintf(":%d", cfg.port),
    Handler:  app.routes(),
    ErrorLog: log.New(logger, "", 0), // custom Logger = io.Writer hədəf
    ...
}
```
- `log.New(logger, "", 0)` → stdlib logger, yazma hədəfi bizim jsonlog
  Logger-dir; `""` və `0` → prefix/flags yox
- Nəticə: server-in bütün daxili mesajları ERROR səviyyəsində JSON kimi
  gedir — HƏR ŞEY bir yerdə, bir formatda

**Alternativ (kitabın tövsiyəsi):** custom logger istəməsəniz — `zerolog`
(sürətli, az allocation, geniş customization).

### 4. recoverPanic middleware
**Problem:** Handler-dəki panic-i http.Server avtomatik recover edir, amma
client-ə cavab YOXDUR — bağlantı sadəcə qapanır.

**Kitabdan kod nümunəsi:**
```go
// cmd/api/middleware.go
func (app *application) recoverPanic(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                w.Header().Set("Connection", "close")
                app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
            }
        }()
        next.ServeHTTP(w, r)
    })
}
```

**Sub-kod izahı:**
- `defer func() { recover() }` → panic stack unwind olunarkən deferred funksiya
  işə düşür
- `Connection: close` header → Go server-inə cavabdan sonra bağlantını QAPAT
  siqnalı — panic-dən sonra bağlantı etibarsızdır
- `recover()` `any` qaytarır → `fmt.Errorf("%s", err)` ilə error-a normalize
- `serverErrorResponse` → JSON log (ERROR) + client-ə 500 JSON cavab

**Routes dəyişikliyi:**
```go
// http.Handler qaytarır (router-in özü yox)
func (app *application) routes() http.Handler {
    router := httprouter.New()
    // ... route-lar
    return app.recoverPanic(router) // middleware router-i wrap edir
}
```

**Mühüm limit (Additional Info):** recoverPanic YALNIZ eyni goroutine-dəki
panic-ləri tutur. Handler daxilində spin olunan background goroutine-dəki
panic recover OLUNMUR (middleware də, http.Server də tutmur) → app çökür.
Həll: hər background goroutine-də öz defer-recover olmalıdır (Chapter 14-də
welcome email göstəriləcək).

## Bu fəsildə layihəyə nə əlavə olundu? (project continuity)

**Yeni fayllar:**
- `internal/jsonlog/jsonlog.go` — Logger, Level, Write
- `cmd/api/middleware.go` — recoverPanic

**Dəyişdirilən:**
- main.go → jsonlog.Logger, ErrorLog adapter, PrintInfo/PrintFatal çağırışları
- errors.go → logError properties ilə
- routes.go → `http.Handler` return + recoverPanic wrap
- application.logger tipi: `*jsonlog.Logger`

## Əsas terminlər

- Structured Logging (strukturlaşdırılmış jurnallama) — machine-readable
  formatda log
- Severity Level (ağırlıq səviyyəsi) — INFO/ERROR/FATAL şiddət dərəcələri
- Mutex (qarşılıqlı istisna kilidi) — paralel yazışları seriyalaşdıran lock
- Stack Trace (izləmə dəstəsi) — xətanın yarandığı çağırış zənciri
- Panic Recovery (panik bərpası) — defer+recover ilə çökmənin tutulması
- Middleware (aralıq proqram təminatı) — handler-ə çatmadan əvvəl işləyən
  qat; `func(next http.Handler) http.Handler` imzası
- Adapter Pattern — bir interfeysi digərinə uyğunlaşdırma (stdlib
  log.Logger → jsonlog)

## Praktik nəticə

Logger + ErrorLog + recoverPanic üçlüyü production observability-nin
təməlidir: bütün loglar bir formatda bir axında, panics isə client üçün
dərk edilə bilən 500-ə çevrilir. `Write()` methodunun io.Writer
interfeysinə uyğunluğu — Go-nun composition gücünün kiçik nümunəsidir:
stdlib-ə toxunmadan custom logger-i bütün ekosistemə bağlayır.
(Müəllim qeydi: Go 1.21+ `log/slog` standart structured logging verir —
yeni layihələrdə jsonlog custom implementasiyası yerinə slog + JSON
handler nəzərdən keçirin.)

## Mənbə
Pages: 234-247 (raw 234-247)
