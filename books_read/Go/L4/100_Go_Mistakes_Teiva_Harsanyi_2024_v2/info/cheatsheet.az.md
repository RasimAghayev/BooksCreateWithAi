# 100 Go Mistakes — Cheat Sheet (AZ)

## Kod təşkilatı

```go
// #1 shadowing: xaricdə elan, daxildə assign
var client *http.Client
var err error
if tracing {
    client, err = createClientWithTracing()   // = YOX :=
}

// #2 happy path solda
if s1 == "" { return "", errors.New("s1 is empty") }
return concatenate(s1, s2)

// #11 functional options
type Option func(*Server)
func WithPort(p int) Option { return func(s *Server) { s.port = p } }
func NewServer(addr string, opts ...Option) *Server
```

## Slice / map

```go
// #21 init
s := make([]T, 0, n)           // n məlumsa

// #22 boşluq yoxlaması — hər ikisini tutur
if len(s) == 0 { ... }

// #25 append yan təsiri
sCopy := make([]int, len(s)); copy(sCopy, s)   // f(sCopy) təhlükəsiz

// #26 capacity leak
uuid := strings.Clone(log[:36])                // və ya kopya

// #28 map kiçilmir — recreating
m = make(map[K]V, len(m))  // əvvəlcədən kopyala-yeni

// #30 range kopya
for i := range accounts { accounts[i].balance += 100 }

// #32 pointer tələsi
for _, c := range customers { c := c; m[c.ID] = &c }  // və ya &customers[i]

// #33 map sırasız — sort-la
```

## String

```go
// #37 rune iterasiya
for i, r := range s { ... }         // r rune, i bayt indexi
// #39 birləşdirmə
var sb strings.Builder
sb.Grow(total)
// #41 substring
u := strings.Clone(bigStr[:36])
```

## Funksiya / xəta

```go
// #46 API
func count(r io.Reader) (int, error)     // fayl adı YOX
// #47 defer closure
defer func() { notify(status) }()         // status DƏYİŞƏNDƏN SONRA oxunur
// #49-51 wrap/is/as
return fmt.Errorf("x: %w", err)
errors.Is(err, ErrX); errors.As(err, &target)
// #54 defer close xətası
defer func() {
    if ferr := f.Close(); ferr != nil && err == nil { err = ferr }
}()
```

## Paralellik

```go
// #58 memory model: send → receive zəmanəti
// #59 CPU-bound: GOMAXPROCS; IO-bound: daha çox
// #60 bloklama select ilə
select {
case <-ctx.Done(): return ctx.Err()
case v := <-ch:
}
// #61 detach
publish(context.Background(), resp)   // və ya cancel-noop ctx
// #62 lifecycle
go func() {
    for {
        select {
        case <-w.quit: return          // exit point MÜTLƏQ
        case v := <-w.ch: ...
        }
    }
}()
// #64 prioritet — nested select
for {
    select {
    case v := <-messageCh: handle(v)
    case <-disconnectCh:
        for { select {
            case v := <-messageCh: handle(v)
            default: return
        }}
    }
}
// #66 nil kanal merge
ch1Closed → ch1 = nil  // select onu gözləmir
// #71 WaitGroup
wg.Add(n)  // ƏVVƏLDƏN (spin-updan öncə)
// #73 errgroup
g, ctx := errgroup.WithContext(ctx)
g.Go(func() error { ... })
err := g.Wait()
// #74 sync kopyası YOX — pointer receiver / go vet
```

## Standart kitabxana

```go
// #75
time.Second      // 100 YOX (ns!)
// #76 loop timer
timer := time.NewTimer(d)
defer timer.Stop()
for { timer.Reset(d); select {...} }
// #77 JSON embedded → adlandır; map[string]any → float64
// #79
resp, err := http.Get(url)
if resp != nil { defer resp.Body.Close() }   // err!=nil olsa da
// #81
client := &http.Client{Timeout: 5 * time.Second}
```

## Test

```bash
go test -race ./...          # CI-də məcburi
go test -short ./...         # ağır testlər skip
go test -shuffle=on ./...
go test -bench=. -count=10 | tee stats.txt; benchstat stats.txt
```
```go
// #85 TDT
tests := map[string]struct{ in, want string }{ ... }
// #87 time injeksiya
type Cache struct{ now func() time.Time }
// #88
httptest.NewRequest("GET", "/x", nil)
httptest.NewRecorder()
httptest.NewServer(handler)
// #89 sink
var sink Result
for i := 0; i < b.N; i++ { sink = compute() }
```

## Optimallaşdırma

```go
// #92 false sharing
type Result struct {
    sumA int64
    _    [56]byte   // padding — ayrı cache line
    sumB int64
}
// #94 sahə sırası böyükdən kiçiyə
type T struct {
    a int64    // 8
    b int32    // 4
    c byte     // 1  → az padding
}
// #96 sync.Pool
var pool = sync.Pool{New: func() any { return make([]byte, 1024) }}
buf := pool.Get().([]byte); defer pool.Put(buf)
```
```bash
# #95 escape
go build -gcflags="-m"
# #98 profil
import _ "net/http/pprof"   # /debug/pprof/
go tool pprof http://host/debug/pprof/profile?seconds=30
# #100
GOMAXPROCS=container_limit   # və ya uber-go/automaxprocs
```

## Qadağan siyahısı (yaddaş cədvəli)

- `:=` daxildə xarici dəyişən üçün → kölgələmə (#1)
- `&customer` loop daxilində (#32), `time.After` loop-da (#76)
- DefaultClient / DefaultServer (#81)
- `fmt.Sprintf("%v", mutexliStruct)` — deadlock (#68)
- map sırasına güvənmə (#33), select determinizm (#64)
- sync tipini value kopyalamaq (#74)
- `100` duration kimi (#75); `resp != nil`dən əvvəl close (#79)
- util paketi (#13); producer-side interfeys (#6); hər yerdə any (#8)
