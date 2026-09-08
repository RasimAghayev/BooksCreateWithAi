# Learn Go with Pocket-Sized Projects — Cheat Sheet (AZ)

## Test (əsas komandalar)

```bash
go test ./...                 # hamısı
go test -run TestGreet ./...   # ad filtri
go test -race ./...           # data race yoxlaması (MÜTLƏQ paralel kodda)
go test -short ./...           # ağır/integration testləri skip
go test --trimpath -race .    # race reportda təmiz yollar
go test -bench=. -benchmem    # benchmark
```

## Example + TDT

```go
func ExampleMain() {
    main()
    // Output:
    // Hello world
}

func TestGreet(t *testing.T) {
    tests := map[string]struct{ lang language; want string }{
        "English": {lang: "en", want: "Hello world"},
    }
    for name, tc := range tests {
        t.Run(name, func(t *testing.T) {
            if got := greet(tc.lang); got != tc.want {
                t.Errorf("want %q, got %q", tc.want, got)
            }
        })
    }
}
```

## Xəta idarəsi

```go
// wrap:
return nil, fmt.Errorf("unable to open %q: %w", path, err)
// sentinel:
const ErrCorpusIsEmpty = corpusError("corpus is empty")
errors.Is(err, ErrCorpusIsEmpty)
// typed:
type InvalidInputError struct{ field, reason string }
var e InvalidInputError
errors.As(err, &e)
```

## Fayl + JSON

```go
f, err := os.Open(path)
defer f.Close()
json.NewDecoder(f).Decode(&data)   // stream
// tam fayl: data, err := os.ReadFile(path)
// buffer: bufio.NewReaderSize(f, 1024*1024)
```

## Flag + args

```go
from := flag.String("from", "", "source currency")
flag.Parse()
arg := flag.Arg(0)   // pozisional
flag.Usage(); os.Exit(1)
```

## Map idiomları

```go
v, ok := m[key]              // comma-ok
count[book]++                 // counter (zero value)
words := strings.Fields(s)   // boşluqla ayır
sort.Sort(byAuthor(books))   // Len/Less/Swap
sort.Slice(books, func(i, j int) bool { return books[i].Author < books[j].Author })
```

## Kitabxana: options + writer

```go
type Option func(*Logger)
func New(level Level, options ...Option) *Logger {
    l := &Logger{threshold: level, output: os.Stderr}
    for _, o := range options { o(l) }
    return l
}
func WithOutput(w io.Writer) Option { return func(l *Logger) { l.output = w } }
```

## Float YOX — Decimal

```go
type Decimal struct {
    subunits  int64   // 1.52 → 152
    precision byte    // → 2
}
// vur: subunits vurulur, precision TOPLANIR, simplify()
```

## HTTP xəta statusları

```go
if errors.Is(err, repository.ErrNotFound) {
    http.Error(w, "not found", http.StatusNotFound)
    return
}
w.Header().Set("Content-Type", "application/json")
w.WriteHeader(http.StatusCreated)
json.NewEncoder(w).Encode(v)
// Timeout:
client := &http.Client{Timeout: 30 * time.Second}   // http.Get YOX
```

## Mock HTTP

```go
ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, `<xml>...`) }))
defer ts.Close()
c := Client{path: ts.URL}
```

## Generics

```go
type Cache[K comparable, V any] struct { data map[K]V }
func New[K comparable, V any]() Cache[K, V]
c := cache.New[int, string]()
// constraint: type Number interface{ int | int64 | float64 }
```

## Paralellik

```go
var wg sync.WaitGroup
wg.Add(2)
go func() { defer wg.Done(); work() }()
wg.Wait()

// (golang.org/x/sync/errgroup) — xəta toplayır:
g, _ := errgroup.WithContext(ctx)
g.Go(func() error { return work() })
err := g.Wait()

// select + quit:
for {
    select {
    case <-s.quit: return
    case p := <-s.pathsToExplore: go explore(p)
    }
}
```

## Mutex

```go
c.mu.Lock()
defer c.mu.Unlock()   // HƏMİŞƏ defer
// çox-oxu: RWMutex → RLock/RUnlock
```

## Loop kölgəsi (Go ≤1.22)

```go
for p := range ch {
    p := p          // closure-a safe copy
    go func() { use(p) }()
}
```

## gRPC quruluşu

```bash
protoc -I=api/proto/ --go_out=api/ --go_opt=paths=source_relative \
       --go-grpc_out=api/ --go-grpc_opt=paths=source_relative api/proto/*.proto
# və ya: //go:generate bash -c "protoc ..."
go generate ./...
```
```go
// typed error → status:
if errors.As(err, &invalidErr) {
    return nil, status.Errorf(codes.InvalidArgument, "%s", err)
}
```

## Context

```go
ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
defer cancel()
// Value QADAĞAN — parametr kimi ötür; hər uzaq çağırış öz deadline-i
```

## Template

```go
//go:embed index.xhtml
var indexPage string

tpl := template.New("index").Funcs(template.FuncMap{
    "statusCSSClass": statusCSSClass,
}).Parse(indexPage)
tpl.Execute(w, data)
```
```html
{{ if .Habits }}{{ range .Habits }}
<li class="{{ statusCSSClass . }}">{{ .Name }}</li>
{{ end }}{{ end }}
```

## Wasm

```bash
GOOS=js GOARCH=wasm go build -o main.wasm
```
```go
js.Global().Set("fn", js.FuncOf(func(this js.Value, args []js.Value) any { ... }))
<-make(chan struct{})   // main-ı saxla!
```

## TinyGo

```bash
tinygo build --target=arduino-nano33 main.go
tinygo flash -target=arduino-nano33
```
```go
led := machine.LED
led.Configure(machine.PinConfig{Mode: machine.PinOutput})
led.High(); led.Low()
```

## Zero dəyərlər

```go
var i int          // 0
var s string       // ""
var b bool         // false
var sl []int       // nil   (append işləyir!)
var m map[k]v      // nil   (YAZMAQ PANIC; oxu zero qaytarır)
var p *T           // nil
```
