# Chapter 4 — A log story: Creating a library (səh. 122-153)

## Bu chapter nədən bəhs edir?

İlk kitabxana (pocketlog): 3 səviyyəli logger qurma — exported API dizaynı,
iota enum, `New()` konstruktor, variadic funksiyalar, internal/external testlər,
godoc sənədləşməsi, `io.Writer` interfeysi, default dəyərlər mübahisəsi və
functional options pattern.

## Əsas fikirlər

### 1. Library (kitabxana) vs application
Kitabxana = başqalarının kodu üçün təməl. Əsas prinsip: **stabil API** —
evolyusiya olsa da istifadəçi kodu dəyişməməlidir. Fayl təşkilati: `pocketlog/`
paketi + `go.mod` module.

### 2. iota enum — səviyyələr
```go
type Level int // yeni tip = domain mənası

const (
    LevelDebug Level = iota // 0
    LevelInfo               // 1
    LevelError              // 2
)
```
**Sub-kod izahı:**
- `iota` → const blokda avtomatik artım (0, 1, 2...)
- `Level` tipi → adi int deyil; yanlış tip keçməz

### 3. Exported API strukturu
İki yanaşma müqayisə olunur:
- Metodlar `Logger` struct-ında: `Debugf`, `Infof`, `Errorf`
- Strukt daxilində `format` metodu — private implementasiya, public davranış

### 4. Variadic funksiyalar
```go
func (lgr *Logger) Debugf(format string, args ...any) {
    lgr.logf(LevelDebug, format, args...) // ... → slice-i aç
}
```
- `args ...any` → istənilən sayda arqument (slice kimi `args`)
- Çağıranda `args...` → slice-i geri açır

### 5. New() konstruktoru
```go
func New(minLevel Level) *Logger {
    lgr := &Logger{
        minLevel: minLevel,
        output:   os.Stdout, // default
    }
    return lgr
}
```
- Sıfır dəyərli `Logger` nil-pointer verməsin deyə `New` məcburi edir
- Struktur daxili sahələr (output, minLevel) unexported — inkapsulyasiya

### 6. Internal vs external testlər
```
pocketlog/
├── level.go        (package pocketlog)
├── logger.go       (package pocketlog)
└── logger_test.go  (package pocketlog_test — XARİCİ test)
```
- `package pocketlog_test` → yalnız **exported** API yoxlayır (istifadəçi
  gözləntisi kimi) — "closed-box testing"
- `*_internal_test.go` → unexported hissələr üçün
- Hər iki paket bir qovluqda ola bilər — konflikt yoxdur

### 7. godoc sənədləşməsi
```go
/*
Package pocketlog exposes an API to log your work.

First, instantiate a logger with pocketlog.New, giving it a threshold level.
Messages of lesser criticality won't be logged.
*/
package pocketlog
```
- Exported hər şeyin üstündə şərh → `godoc` saytı/doc
- "doc.go" faylı — paket sənədi üçün ayrıca fayl

### 8. io.Writer — çıxışın soyudulması
```go
type Logger struct {
    minLevel Level
    output   io.Writer // hər hansı yazıla bilən hədəf
}
```
Testdə:
```go
type testWriter struct{ contents string }

func (tw *testWriter) Write(p []byte) (int, error) {
    tw.contents += string(p) // yazılanları yoxlumaq üçün saxla
    return len(p), nil
}
```
**Sub-kod izahı:**
- Logger konkret `os.Stdout` deyil, `io.Writer` asılıdır → testdə fake writer
- `fmt.Fprintf(lgr.output, format+"\n", args...) → hər hansı hədəfə yazır

### 9. Default dəyərlər + Functional Options
Mübahisə: default-lar bootstrap-i asanlaşdırır, amma gizli xətalar yaradır.
Həll — Option pattern:
```go
type Option func(*Logger) // konfiq funksiyası

func WithOutput(output io.Writer) Option {
    return func(lgr *Logger) {
        lgr.output = output // yalnız output-u dəyişən closure
    }
}

func New(minLevel Level, opts ...Option) *Logger {
    lgr := &Logger{minLevel: minLevel, output: os.Stdout} // default
    for _, opt := range opts {
        opt(lgr) // istifadəçi istəkləri tətbiq olunur
    }
    return lgr
}

// istifadə:
lgr := pocketlog.New(pocketlog.LevelInfo, pocketlog.WithOutput(os.Stdout))
```
Məcburi parametr (minLevel) — konstruktorda; optional-lar — Option kimi.

### 10. Logging best practices
- Səhv (error) logla, amma çox danışma — pul/log storage birbaşa mütənasibdir
- Uzun mesaj YOX — qısa, strukturlaşdırılmış
- **Logging debugging aləti deyil** — kod niyə işlədiyini bilmirsənsə, kod
  (və ya adlar) səthəd; log-la debug etmək əvəzinə test yaz
- Log-a timestamp əlavə etmək mümkündür (side quest)

## Əsas terminlər

- Library API (kitabxana interfeysi)
- Enumeration (iota) (saylım)
- Variadic Function (dəyişən saylı funksiya)
- Constructor (konstruktor)
- Internal/External Test (daxili/xarici test)
- Closed-Box Testing (qapalı qutu testi)
- Functional Options Pattern
- io.Writer (yazı interfeysi)

## Praktik nəticə

- Kitabxana qurarkən: exported API minimal + stabil; unexported implementasiya
- `New(...)` + functional options = məcburi + optional parametrlərin təmiz qarışığı
- External test paketi istifadəçi gözləntisini yoxlayır
- `io.Writer` asılılığı → testlərdə fake writer, istəsən fayl/buffer

## Mənbə

Pages: 122-153 (Chapter 4, Learn Go with Pocket-Sized Projects)
