# Chapter 4 — A log story: Creating a library (səh. 122-153)

## Bu chapter nədən bəhs edir?

"pocketlog" adlı üçsəviyyəli logging library-nin yaradılması: package qaydaları,
iota enum, variadic funksiyalar, New() konstruktor pattern-i, internal vs
external testlər, io.Writer injection, functional options pattern və logging
best practices.

## Layihə məqsədi

Susan ilə 2 saatlıq bug ovu → "nə üçün logger yazmırıq?" → company-wide istifadə
olunacaq STABİL API-li library: `New(level, options...)` + `Debugf/Infof/Errorf`.

## Əsas fikirlər

### 1. Go package qaydaları
- Hər fayl `package NAME` ilə başlayır; qovluqdakı bütün fayllar eyni adı paylaşır
- Ad: qısa, lowercase, tək söz; qovluq adı ilə üst-üstə düşməsi adətdir
- **Export:** böyük hərflə başlayan simvollar (funksiya/type/const/sahə) — xarici
  istifadəyə açıq; kiçik hərflə başlayanlar — package-private
- **Faylları kiçik saxla:** type böyüyəndə metodları fayllara böl (scope = paket,
  fayl deyil — `level.go`, `logger.go`, `options.go`)

### 2. Enum — integer-based new type + iota
```go
type Level int   // yeni tip — təkcə ad deyil, tip təhlükəsizliyi

const (
    LevelDebug Level = iota  // 0
    LevelInfo               // 1
    LevelError              // 2
)
```
**Sub-kod izahı:**
- `Level int` — plain int deyil, ayrıca tip → səhv səviyyə ötürülməsi compile
  vaxtı yaxalanır
- `iota` → hər sətirdə avtomatik artım; eksplisit dəyər yazmağa ehtiyac yoxdur
- Sıra vacibdir: Debug < Info < Error → threshold müqayisəsi `l.threshold <= level`

### 3. API dizaynı — metod sayı seçimi
İki variant müqayisə olunur:
- `l.Log(pocketlog.Info, "message")` — bir metod + səviyyə parametri
- `l.Info("message")` — hər səviyyəyə ayrıca metod ✓ (seçilən)

Səbəb: kod oxunaqlığı — `l.Info(...)` intent-i birbaşa ifadə edir.

### 4. Variadic funksiyalar
```go
func (l Logger) Debugf(format string, args ...any) {
    // args — funksiya daxilində slice kimi: args[0], len(args), range
}
```
- Son parametr `...T` → istənilən sayda (0+) arqument
- `fmt.Printf` məntiqi: format + istənilən sayda dəyər

### 5. New() konstruktoru + zero value faydası
Go-da konstruktor yoxdur → `New()` funksiyası konvensiyadır:
```go
func New(threshold Level, options ...Option) *Logger { ... }
```
- İstifadə: `pocketlog.New(...)` — package prefix ilə, `NewPocketLog` adına
  ehtiyac yoxdur
- **Zero value useful:** `var l Logger` — sıfır dəyəri belə işləməlidir (nil
  output → `os.Stdout` fallback; zero threshold → hamısı loglanır)
- Hər tipin zero dəyəri var: struct/funksiya/channel/interface/pointer/slice/map

### 6. External vs internal testlər
```
pocketlog/logger.go          // package pocketlog
pocketlog/logger_test.go     // package pocketlog_test  ← EXTERNAL
```
- Eyni qovluqda iki paket ola bilər — `foo` + `foo_test` istisnası testlər üçündür
- External test (`package pocketlog_test`): library-ni **istifadəçi gözləntisi
  ilə** yoxlayır — yalnız export edilmiş API
- "Closed-box" test fəlsəfəsi: internal detal yox, public davranış yoxlanılır

### 7. Dokumentasiya
- Export edilmiş hər elementin üstündə `//` şərhi → IDE hover + godoc
- **doc.go** — paket səviyyəli dokumentasiya: `/* Package pocketlog ... */`

### 8. io.Writer injection — çıxışın təyini
```go
type Reader interface { Read(p []byte) (n int, err error) }
type Writer interface { Write(p []byte) (n int, err error) }
```
Logger çıxışını hardcode etmə → `io.Writer` sahəsi ilə injection:
- İstifadəçi istədiyi destination verir: os.Stdout, fayl, network, buffer
- Default: `output == nil` → `os.Stdout`

### 9. logf — ümumi funksiya (DRY)
`Debugf/Infof/Errorf` eyni məntiqi təkrarlayır → ümumi `logf` metodu:
```go
func (l Logger) logf(level Level, format string, args ...any) {
    if l.threshold <= level {
        _, _ = fmt.Fprintf(l.output, format+"\n", args...)
    }
}
```
Runtime-da səviyyə seçəndə `l.Logf(Level, ...)` də ixrac olunur.

### 10. Functional options pattern
```go
type Option func(*Logger)   // konfiqurasiya funksiyası

func WithOutput(output io.Writer) Option {
    return func(lgr *Logger) { lgr.output = output }
}

func New(threshold Level, options ...Option) *Logger {
    lgr := &Logger{threshold: threshold}
    for _, opt := range options {
        opt(lgr)
    }
    return lgr
}
```
**İstifadə:** `pocketlog.New(pocketlog.LevelInfo, pocketlog.WithOutput(os.Stdout))`
- Tələb olunan parametr (threshold) — pozisiya ilə; optional-lar — funksiya ilə
- Yeni Option əlavə etmək mövcud çağırışları POZMUR (backward compatible)
- Default dəyərlər "ehtiyatla" — cognitive load azaldır, amma məcburiyyət yoxdur

### 11. Test helper — io.Writer mock
```go
type testWriter struct{ contents string }

func (tw *testWriter) Write(p []byte) (n int, err error) {
    tw.contents += string(p)   // in-memory yaddaş
    return len(p), nil
}
```
Testlərdə `WithOutput(tw)` → logger çıxışı yoxlanıla bilən string-ə düşür.

### 12. Logging best practices
- **Nə loglamaq:** threshold-un seçimi = signal/noise balansı; storage BAHADIR
- **Uzun mesajlardan çəkin:** minlərlə açarlı map-i yox, ölçüsü/açar varlığı;
  image/audio baytları log-a YOX
- **Həssas data:** parollar, tokendlər, şəxsi məlumat — heç vaxt
- **Logger debug aləti deyil:** kod qeyri-aydındısa → (1) daha aydın kod, (2)
  daha yaxşı dokumentasiya, (3) daha çox test yaz
- Strukturlaşdırılmış format (JSON) — machine-readable; `json.Marshal` yeni
  sətir əlavə ETMİR (`\n` özün əlavə olunmalı)

## Əsas terminlər

- Library (kitabxana)
- Enumeration (siyahı / enum)
- iota (sayğaç sabiti)
- Variadic function (dəyişən saylı funksiya)
- Constructor pattern (`New()`)
- Zero value (sıfır dəyər)
- Closed-box testing (qutu-testi)
- io.Writer injection
- Functional Options (funksional seçənəklər)
- godoc / doc.go

## Praktik nəticə

- Library = stabil export API + internal detallar; `foo_test` external testi
  istifadəçi nöqteyi-nəzərindən yoxlanır
- Enum → yeni tip + iota; sıra semantic (Debug ≤ Info ≤ Error)
- Konstruktor: tələb olunanlar parametr, optional-lar `...Option`
- Çıxış heç vaxt hardcode olunmur — `io.Writer` injection
- Zero value işlək olsun; fayllar kiçik; hər export dokumentasiyalı

## Mənbə

Pages: 122-153 (Chapter 4, Learn Go with Pocket-Sized Projects)
