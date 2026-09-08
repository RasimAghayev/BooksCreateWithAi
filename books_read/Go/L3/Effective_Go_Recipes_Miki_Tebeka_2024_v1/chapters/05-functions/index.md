# Chapter 5 — Working with Functions (səh. 73-87)

## Bu fəsil nədən bəhs edir?

Funksiyalar birinci sinif obyekt kimi: function registry (dispatch),
functional options (With... funksiyaları), closure ilə arqumentli
optionlar, watcher/notification mexanizmi və go:linkname ilə unexported
funksiyalara çıxış (dirty hack).

## Əsas fikirlər

### Recipe 26 — function registry (dispatch)
**Tapşırıq:** Metrik server — Accept header-inə görə JSON/CSV/... format;
yeni format əlavəsi kod dəyişmədən olsun.

```go
// Encoder function type:
type Encoder func(w io.Writer, metrics []Metric) error

// mime type → encoder map:
var registry = make(map[string]Encoder)

// Register registers an encoder for a mime type.
func Register(mimeType string, enc Encoder) error {
    if _, ok := registry[mimeType]; ok {
        return fmt.Errorf("%q already registered", mimeType)
    }
    registry[mimeType] = enc
    return nil
}

// Konkret encoderlər:
func EncodeJSON(w io.Writer, metrics []Metric) error {
    return json.NewEncoder(w).Encode(metrics)
}

func EncodeCSV(w io.Writer, metrics []Metric) error {
    wtr := csv.NewWriter(w)
    if err := wtr.Write([]string{"time", "name", "value"}); err != nil { // header
        return err
    }
    r := make([]string, 3)
    for _, s := range metrics {
        r[0] = s.Time.Format(time.RFC3339)
        r[1] = s.Name
        r[2] = fmt.Sprintf("%f", s.Value)
        if err := wtr.Write(r); err != nil {
            return err
        }
    }
    wtr.Flush()
    return nil
}

func init() {
    Register(csvMimeType, EncodeCSV)
    Register(jsonMimeType, EncodeJSON)
}

// HTTP handler — registry-dən oxu:
func queryHandler(w http.ResponseWriter, r *http.Request) {
    query := r.URL.Query().Get("query")
    if query == "" {
        http.Error(w, "missing query", http.StatusBadRequest)
        return
    }
    mimeType := requestMimeType(r)      // Accept header → default JSON
    enc, ok := registry[mimeType]       // dispatch!
    if !ok {
        http.Error(w, fmt.Sprintf("unsupported mime type - %q", mimeType),
            http.StatusBadRequest)
        return
    }
    metrics, err := queryDB(query)
    // ...
    w.Header().Add("Content-Type", mimeType)   // mütləq əvvəlcədən!
    if err := enc(w, metrics); err != nil {
        log.Printf("can't encode %d metrics with %q - %s", len(metrics), mimeType, err)
    }
}
```

**Test:**
```bash
$ curl 'http://localhost:8080/metrics?query=CPU'                        # JSON
$ curl -H 'Accept: text/csv' 'http://localhost:8080/metrics?query=CPU' # CSV
```
- Yeni format = yeni funksiya + Register — handler dəyişmir
- Registry pattern-in digər misalları: HTTP router (mux),
  `database/sql` driver qeydiyyatı (blank import `_ "github.com/mattn/go-sqlite3"`)

### Recipe 27 — functional options
**Tapşırıq:** Server kitabxanasına seçimlər əlavə et — API-ni pozmadan.

```go
// Server is an HTTP server.
type Server struct {
    verbose bool     // unexported — implementation azadlığı
    port    int
}

// NewServer returns a Server with options.
func NewServer(options ...func(*Server) error) (*Server, error) {
    srv := &Server{
        port: 8080,   // default dəyərlər
    }
    for _, opt := range options {      // hər option-u tətbiq et
        if err := opt(srv); err != nil {
            return nil, err
        }
    }
    return srv, nil
}

// WithVerbose sets the verbose option on s.
func WithVerbose(s *Server) error {
    s.verbose = true
    return nil
}

// İstifadə:
srv, err := NewServer(WithVerbose)
```
- `...func(*Server) error` — variadic option funksiyaları; sıfır option =
  default-lar
- Unexported sahələr → daxili dəyişiklik istifadəçiyə toxunmur
- Dave Cheney-nin "Functional options for friendly APIs" variasiyası

### Recipe 28 — closure ilə arqumentli optionlar
```go
const portErrFmt = "port must be between 0 and %d, got %d"

func WithPort(port int) func(*Server) error {
    const maxPort = 0xFFFF
    return func(s *Server) error {      // closure — port-u yadında saxlayır
        if port <= 0 || port > maxPort {
            return fmt.Errorf(portErrFmt, maxPort, port)   // validasiya!
        }
        s.port = port
        return nil
    }
}

// İstifadə:
srv, err := NewServer(WithPort(9999), WithVerbose)
```
- `WithPort(9999)` ƏVVƏL hesablanır → dəyəri `func(*Server) error` tipinə
  uyğun gəlir
- Closure = funksiyanın təyin olunduğu mühit; compiler `port`-u closure-də
  axtarır
- Anonim funksiya: `var Add = func(a, b int) int {...}` forması

### Recipe 29 — funksiyalarla notification (watchers)
```go
type State byte
const (
    Ready State = iota + 1
    Working
    Done
)

type Task struct {
    ID       uint
    Result   any
    Err      error
    State    State
    Work     func() (any, error)   // işin özü — funksiya!
    Watchers []func(*Task)          // bitəndə xəbər veriləcək funksiyalar
}

func NewTask(id uint, work func() (any, error)) *Task {
    return &Task{ID: id, State: Ready, Work: work}
}

// Subscribe subscribes a watcher to task.
func (t *Task) Subscribe(w func(*Task)) {
    t.Watchers = append(t.Watchers, w)
}

// Execute executes the task.
func (t *Task) Execute() error {
    t.State = Working
    defer func() { t.State = Done }()     // panic olsa belə Done

    t.Result, t.Err = t.Work()
    if t.Err != nil {
        log.Printf("error: Task %d failed - %s", t.ID, t.Err)
    }

    for _, s := range t.Watchers {        // watcher-lərə xəbər ver
        s(t)
    }
    return t.Err
}

// İstifadə — həm anonim, həm bound method:
t := NewTask(7, func() (any, error) { return "done", nil })
t.Subscribe(func(t *Task) {               // anonim watcher
    log.Printf("info: w1: from %d - %#v %v", t.ID, t.Result, t.Err)
})

type Watcher struct{}
func (w *Watcher) Handle(t *Task) { ... }
var w Watcher
t.Subscribe(w.Handle)                     // BOUND METHOD da funksiyadır!

t.Execute()
// w1: from 7 - "done" <nil>
// w2: from 7 - "done" <nil>
```
- `any` generics ilə type-safe edilə bilər (Recipe 39-a bax)
- Bound method funksiya kimi ötürülə bilər — Subscribe hər ikisini qəbul edir

### Recipe 30 — unexported funksiyaya çıxış (go:linkname)
**Tapşırıq:** Client-də 1MB limit var; `setMaxBodySize` unexported; developer
2 ay sonra baxacaq; deadline yaxındır.

```go
import (
    "fmt"
    _ "unsafe"                          // məcburi — xatırlatma üçün!
    "git.corp.com/client"
)

// HACK: Access internal client.setMaxBodySize.
// Remove this once setMaxBodySize becomes exported, see issue #732.
//go:linkname setClientMaxBodySize git.corp.com/client.setMaxBodySize
func setClientMaxBodySize(int64)

// İstifadə:
setClientMaxBodySize(5 * 1_000_000)   // 5 MB
c := client.New()
```
- `//go:linkname` — compiler directive (pragma): iki simvolu birləşdirir;
  encapsulation-ı və hətta type safety-ni pozur
- **Yalnız çox yaxşı səbəblə** — kod oxunmazlaşır, unexported dəyişə bilər
- `_ "unsafe"` importu mütləq — diqqət çəkmək üçün
- HACK şərhi + issue linki — technical debt izlənilsin
- Alternativlər: vendoring (go mod vendor + əl ilə fix), replace directive,
  fork — hamısının öz qiyməti var

**Digər direktivlər:** `//go:build` (şərti kompilyasiya), `//go:noinline`
(benchmark-larda inline qadağası — müəllifin şəxsi istifadəsi).

## Final Thoughts-dən

Funksiyalar birinci sinif obyektdir — registry, options, notification
kimi kompakt həllər mümkündür; hətta database/sql belə işlədir. Öz
kodunuzda belə yerlər axtarın.

## Əsas terminlər
- First-class functions — funksiyaların dəyişən/strukturda saxlanması
- Function registry (dispatch) — açar → funksiya xəritəsi
- Functional options — `...func(*T) error` seçim üslubu
- Closure — təyin olunduğu mühiti yadında saxlayan funksiya
- Anonymous function — adsız funksiya (var Add = func...)
- Bound method — `obj.Method` funksiya kimi ötürülür
- go:linkname — unexported simvola çıxış direktivi
- Compiler directive (pragma) — //go: prefiksli şərhlər
- Blank import (_) — yalnız side-effect (qeydiyyat) üçün

## Praktik nəticə
Funksiya = data: registry map-i ilə format/handler dispatch; functional
options (WithX/WithX(arg)) API-ni sabit saxlayıb konfiqurasiyanı
artırır; closure arqumentli optionların daşıyıcısıdır; watcher slice-i
observer pattern-in ən yüngül formasıdır. go:linkname yalnız açıq
placeholder (HACK + issue) ilə; normal yollar: PR gözlə, vendor, replace,
fork — qiymət/qərar balansını edin.

## Mənbə
Pages: 73-87 (PDF 73-87)
