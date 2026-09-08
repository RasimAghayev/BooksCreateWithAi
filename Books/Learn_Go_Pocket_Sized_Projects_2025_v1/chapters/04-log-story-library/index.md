# Chapter 4 — A Log Story: Creating a Library (səh. 122-153)

## Bu fəsil nədən bəhs edir?

pocketlog — paylaşıla bilən üçsəviyyəli log kitabxanası: enum (iota), stabil
export API, variadic funksiyalar, New() konstruktoru, functional options,
io.Writer asılılığı, external/internal test, doc.go sənədləşməsi, logging best
practices.

## Əsas fikirlər

### 1. Enum — iota
```go
type Level int

const (
    LevelDebug Level = iota  // 0
    LevelInfo                // 1
    LevelError               // 2
)
```
- Integer-based yeni tip + iota — "enum kimi" oxuna bilən, sətir sayı azdı

### 2. Logger strukturu + threshold
```go
type Logger struct {
    threshold Level
    output    io.Writer
}

func New(threshold Level, options ...Option) *Logger {
    lgr := &Logger{threshold: threshold, output: os.Stderr}  // default
    for _, option := range options { option(lgr) }
    return lgr
}

type Option func(*Logger)  // functional options

func WithOutput(output io.Writer) Option {
    return func(lgr *Logger) { lgr.output = output }
}

func (l *Logger) Infof(format string, args ...any) {
    l.logf(LevelInfo, format, args...)
}
```
**Sub-kod izahı:**
- Məcburi parametrlər (threshold) + optional Option-lar — gələcək dəyişikliklər
  API-ni POZMUR
- Variadic `args ...any` — Printf tərzi istənilən sayda arqument; `any` =
  `interface{}` (1.18 alias)
- Level metodları (Debugf/Infof/Errorf) oxunaqlılıq üçün ayrı-ayrı; daxili
  `logf` — bir implementation

### 3. io.Writer — test edilə bilən çıxış
```go
// istifadəçi: pocketlog.New(LevelInfo, pocketlog.WithOutput(os.Stdout))
// test: öz testWriter-ın ilə
type testWriter struct{ contents string }
func (tw *testWriter) Write(p []byte) (int, error) {
    tw.contents = tw.contents + string(p)
    return len(p), nil
}
```
- Logger stdout-a bağlı deyil — hər hansı io.Writer (fayl, buffer, test
  capture). Bu, kitabxana testinin açarıdır.

### 4. External vs internal test paketi
- `package pocketlog_test` (logger_test.go) — yalnız İXRAC OLUNAN API test
  olunur → istifadəçi perspektivi, refactor azadlığı
- `package pocketlog` (logger_internal_test.go) — unexported-lara da çıxış
- Qayda: eyni qovluqda iki paket — yalnız `_test.go` faylları üçün istisna

### 5. doc.go — paket sənədləşməsi
```go
/*
Package pocketlog exposes an API to log your work.

First, instantiate a logger with pocketlog.New, giving it a threshold level.
Messages of lesser criticality won't be logged.
*/
package pocketlog
```
- Blok şərh paket klauzundan ƏVVƏL → godoc paket səhifəsi
- Export olunan hər şeyin şərhi = dokumentasiya (IDE hover-də görünür)

### 6. Logging best practices
- **Nə log etmək:** funksiya girişləri, xarici çağırışların sonucları, gözlənilməz
  hallar — AXI (value) üçün deyil
- **Qısa mesajlar:** hər log = saxlanma qiyməti; minikli map-i çap etmə
- **Log ≠ debug aləti:** kod qeyri-aydındırsa log əlavə etməkdən öncə kodu
  AYDINLAŞDIR / parçala / test yaz ("trust the code")

## Əsas terminlər

- Enum / iota
- Functional Options Pattern
- Variadic Function
- io.Writer injection
- doc.go
- Internal / External test

## Praktik nəticə

- Kitabxana dizaynı: minimal New + Option-lar; bütün dönmələr stabil API
- Çıxış həmişə io.Writer olmalıdır — test üçün
- Export = müqavilə: şərhlər + external test + doc.go

## Mənbə

Pages: 122-153 (Chapter 4, Learn Go with Pocket-Sized Projects)
