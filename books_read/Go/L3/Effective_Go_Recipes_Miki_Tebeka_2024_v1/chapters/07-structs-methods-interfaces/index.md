# Chapter 7 — Structs, Methods, and Interfaces (səh. 112-124)

## Bu fəsil nədən bəhs edir?

Go-nun OO-ya yanaşması (data-oriented programming): ad hoc interfeyslər
(syncer), http.ResponseWriter sargısı (errWriter), generics — tip məhdudiyyəti
(Number), generic Ring buffer, generic UnmarshalJSON ilə pointer təhlükəsizliyi.

## Əsas fikirlər

### Recipe 37 — ad hoc interfeys (syncer)
**Tapşırıq:** Logger flush etsin — amma yalnız writer Sync() dəstəkləyirsə.

```go
// Böyük interfeysdən qaç — yalnız lazım olan metod:
type syncer interface {
    Sync() error
}

// NOP syncer — heç nə edən standart:
type nopSyncer struct{}
func (nopSyncer) Sync() error { return nil }

type Logger struct {
    level Level
    w     io.Writer
    s     syncer          // həmişə Sync çağırıla bilər
}

func NewLogger(level Level, out io.Writer) Logger {
    log := Logger{level, out, nopSyncer{}}   // default: nop
    if s, ok := out.(syncer); ok {           // comma-ok assertion
        log.s = s                             // dəstəkləyirsə — həqiqisi
    }
    return log
}

func (l Logger) log(level Level, format string, args ...any) {
    if l.level > level {
        return
    }
    msg := fmt.Sprintf(format, args...)
    ts := time.Now().UTC().Format(time.RFC3339)
    fmt.Fprintf(l.w, "[%s] - %s - %s\n", ts, level, msg)
    l.s.Sync()                                // hər yazıdan sonra
}
```
- Alternativ (İMTİNA edilən): `WriteSyncer{io.Writer; Sync() error}` — amma
  bu, istifadə oluna bilən tipləri MƏHDUDlaşdırır
- **Go Proverb (Rob Pike): "The bigger the interface, the weaker the
  abstraction."**

### Recipe 38 — http.ResponseWriter sargısı
**Tapşırıq:** bütün >=400 status-ları logla (monitorinq standartı).

```go
type errWriter struct {
    http.ResponseWriter        // embed → interfeysin hamısı əlçatan
    statusCode int
}

func (ew *errWriter) WriteHeader(statusCode int) {
    ew.statusCode = statusCode                  // YADDA SAXLA
    ew.ResponseWriter.WriteHeader(statusCode)    // orijinala ötür
}

func logErr(w http.ResponseWriter, r *http.Request) {
    ew := errWriter{ResponseWriter: w}           // w-ni əvəz et
    http.DefaultServeMux.ServeHTTP(&ew, r)      // mux-a sarğılı ver
    if ew.statusCode >= http.StatusBadRequest {
        log.Printf("error: %s %s <%d>", r.Method, r.URL.Path, ew.statusCode)
    }
}

func main() {
    http.HandleFunc(offsetPrefix, offsetHandler)
    handler := http.HandlerFunc(logErr)     // qlobal handler
    if err := http.ListenAndServe(":8080", handler); err != nil {
        log.Fatalf("error: %s", err)
    }
}
```
- Embed ilə interfeysi "hissəvi" implement edirsiniz — yalnız
  dəyişdirmək istədiyiniz metodu yazın
- Bu klassik irsiyyət DEYİL: digər metodların receiver-i embed olunan
  tipdir, sizin tip yox

### Recipe 39 — generics ilə kod həcmi
```go
// Number is set of possible numbers.
type Number interface {
    ~float64 | ~int          // ~ → int ƏSASINDA tiplər də (type Level int)
}

// Max returns the maximal value in values.
func Max[T Number](values []T) (T, error) {
    if len(values) == 0 {
        var zero T           // tipin zero value-su — 0 yazmaq olmaz!
        return zero, fmt.Errorf("Max of empty slice")
    }
    max := values[0]
    for _, v := range values[1:] {
        if v > max {
            max = v
        }
    }
    return max, nil
}

// İstifadə — kompilyator tipi çıxarır:
iVals := []int{15, 42, 16, 8, 23, 4}
fmt.Println(Max(iVals))              // 42 <nil>
fVals := []float64{3.14, 2.718, 6.283, 1.618}
fmt.Println(Max(fVals))              // 6.283 <nil>
_, err := Max[int](nil)              // nil tipsizdir → AÇIQ tip lazımdır
```
- `~int` — underlying type int olanlar (Level int kimi)
- `var zero T` — generic zero value idiomu
- Mənbələr: Ian Lance Taylor "When to Use Generics";
  golang.org/x/exp/{maps,slices,constraints}

### Recipe 40 — generic Ring buffer
```go
type Number interface {
    ~int | ~int16 | ~int32 | ~int64 |
    ~uint | ~uint16 | ~uint32 | ~uint64 |
    ~float32 | ~float64
}

// Ring is a circular ring buffer.
type Ring[T Number] struct {
    size   int
    i      int
    values []T
}

func NewRing[T Number](size int) (*Ring[T], error) {
    if size <= 0 {
        return nil, fmt.Errorf("size must be > 0")
    }
    r := Ring[T]{
        size:   size,
        values: make([]T, size),
    }
    return &r, nil
}

// Push pushes item to the ring, possibly overwriting an old value.
func (r *Ring[T]) Push(v T) {
    r.values[r.i] = v
    r.i = (r.i + 1) % r.size       // dairəvi hərəkət
}

// Mean returns the mean of values in the ring.
func (r *Ring[T]) Mean() float64 {
    var s T = 0
    for _, v := range r.values {
        s += v
    }
    return float64(s) / float64(r.size)
}

// İstifadə:
r, _ := NewRing[int](4)
for i := 1; i <= 10; i++ {
    r.Push(i)     // son 4 dəyər: 7, 8, 9, 10
}
fmt.Println(r.Mean())   // 8.5
```
- Generics-dən əvvəl `any` → tip sistemi bypass olunurdu (string-i int
  ringinə push etmək mümkün idi!)
- Məhdudiyyət: generic METOD-lar yalnız T tipindən istifadə edə bilər —
  praktikada problem deyil

### Recipe 41 — generics ilə tip təhlükəsizliyi
**Problem:** 1934, 2013, 2083, 2111 — hamısı `json.Unmarshal`-a non-pointer
göndərilməsindən!

```go
type UserRequest struct {
    Login string
}
type GroupRequest struct {
    ID string
}

// Request is the set of all possible requests.
type Request interface {
    UserRequest | GroupRequest      // tip UNİONU
}

// UnmarshalJSON implements json.Unmarshaler.
func UnmarshalJSON[T Request](data []byte, d *T) error {
    return json.Unmarshal(data, d)
}

// DÜZGÜN — compile olur:
data := []byte(`{"login": "elliot"}`)
var r UserRequest
UnmarshalJSON(data, &r)

// XƏTA — COMPILE OLMAZ ("not a pointer"):
var r UserRequest
UnmarshalJSON(data, r)
```
- encoding/json generics-dən əvvəl yazılıb → `any` qəbul edir → runtime-də
  reflection ilə tip yoxlaması → 4 bug!
- `*T` parametr → pointer olmayan arqument COMPILE vaxtı yıxılır
- Flexibility + type safety birlikdə

## Final Thoughts-dən

**Qayda barmağı: "Accept interfaces, return types."** — `os.Open` →
`*os.File` qaytarır; `io.Copy` → `io.Reader/io.Writer` qəbul edir. Bu,
idiomatik Go-nun ilk addımıdır.

## Əsas terminlər
- Data-oriented programming — Go-nun OO-yoloji yanaşması
- Ad hoc interface — kiçik, məqsədli interfeys (syncer)
- NOP (no operation) — heç nə edən tip/metod
- Embedding ilə interfeys implementasiyası — yalnız fərqli metodlar
- Type constraint — generic tip məhdudiyyəti (Number)
- ~ (tilde) — underlying type uyğunluğu
- Tip unionu — `A | B` interfeys elementləri
- var zero T — generic zero value idiomu
- Ring buffer — dairəvi bufer (% modul)
- Accept interfaces, return types — idiomatik API qaydası

## Praktik nəticə
İnterfeysləri kiçik saxlayın — istifadə yeri təyin edən ad hoc interfeys
(comma-ok assertion) ən çevik həllüddür. Mövcud interfeysi dəyişmək
istəyirsinizsə embed edin, yalnız öz metodunuzu override edin
(errWriter pattern-i — middleware-in əsası). Kolleksiyalar və yardımçı
funksiyalar üçün generics: constraint interface + ~tilde + var zero T;
serializasiya wrapper-lərində `*T` parametri pointer səhvlərini
compile-vaxtı yaxalayır.

## Mənbə
Pages: 112-124 (PDF 112-124)
