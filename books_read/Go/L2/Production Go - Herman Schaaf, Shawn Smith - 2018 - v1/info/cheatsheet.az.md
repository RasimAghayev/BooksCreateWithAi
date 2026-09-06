# Production Go — Cheat Sheet (bütün kitabdan toplanmış)

Herman Schaaf, Shawn Smith — Leanpub, 2018

## Tooling

### Quraşdırma / editor
```bash
go version                              # versiya yoxla
go get golang.org/x/tools/cmd/goimports   # format + import idarəsi
# vim: let g:go_fmt_command="goimports"
go get -u github.com/alecthomas/gometalinter && gometalinter --install
gometalinter --deadline=180s --exclude=vendor --enable=misspell ./...
```

### go vet
```bash
go vet main.go
# arg "test" for printf verb %d of wrong type: string
# possible formatting directive in Println call
```

### gofmt
```bash
gofmt -w -s file.go     # -s: simplify ([]Animal{Animal{...}} → {{...}})
```

### golint
```bash
golint ./...
# exported type Entry should have comment or be unexported
# if block ends with a return statement, so drop this else
```

## Sintaksis

### Dəyişənlər
```go
var x int        // zero value
var x = 1        // tip çıxarılır
x := 1           // short — funksiya daxilində
b, c := 1, 2     // OK — c yeni; var ilə REDECLARE XƏTADIR
```

### Tiplər
```go
int64(i32) + i64                  // explicit cast mütləq
var i int8 = 128                  // constant overflows int8
var p, s uint32 = 0, 1
p - s                              // 4294967295 — UNDERFLOW!
```

### Struct + JSON
```go
type Person struct {
    Name  string `json:"name"`      // kiçik hərf; sahələr EXPORTED olmalı (yoxsa {})
}
p := struct{ Name string `json:"name"` }{...}     // anonymous
```

### Kanallar
```go
ch := make(chan int, 5)
select {
case ch <- 6:          // dolu YOXDURSA göndər
default:               // doludursa ötür
}
```

### Nil interface
```go
var s *string              // s == nil: true
var i interface{} = s      // i == nil: FALSE — (tip, dəyər) hər ikisi nil olmalı
```

## Strings / Unicode

```go
strings.Split("a,b,c", ",")          // ["a" "b" "c"]
strings.Count("banana", "ana")        // 1 — NON-OVERLAPPING
strings.Index("banana", "an")        // 1; yoxdursa -1
strings.Contains / HasPrefix / HasSuffix
strings.FieldsFunc(s, isNotLetter)     // funksiya İLƏ split → sözlər
strings.EqualFold(a, b)               // case-insensitive ==
```

### Printf verb-ləri
```go
%s   "ABC 你好"          %x   41424320e4bda0e5a5bd
%q   "ABC 你好"          %+q  "ABC \u4f60\u597d" (ASCII-only!)
%v / %+v (sahə adları) / %#v / %T / %%
Flag: + - # ' ' 0       Genişlik RUNE ilə (C-də bayt!)
```

### Range + rune
```go
for i, r := range "ABC你好" { }
// 'A'(0) 'B'(1) 'C'(2) '你'(3) '好'(6) — index BAYT, dəyər RUNE
```

## Concurrency

### WaitGroup
```go
var wg sync.WaitGroup
wg.Add(1)                   // Add sayı = goroutine sayı (çoxu deadlock!)
go func() {
    defer wg.Done()
    ...
}()
wg.Wait()
```

### Handler-da fon iş + poller
```go
http.HandleFunc("/", func(w, r) {
    go func() { ... }()     // handler qayıdır, goroutine DAVAM
    return
})

func poll() {
    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()
    for {
        select {
        case <-ticker.C:
            resp, _ := http.Get(endpoint)
            defer resp.Body.Close()
            ...
        }
    }
}
go poll()    // sonsuz loop → goroutine ŞƏRT
```

### Race + həllər
```go
// RACE: 2 goroutine eyni map-ə yazır
// 1) sync.Map — YALNIZ: write-once read-many VEYA disjoint keys
// 2) TÖVSİYƏ: safeMap:
type safeMap struct {
    sync.Mutex
    m map[string]string
}
func (sm *safeMap) Store(key, val string) {
    sm.Lock()
    defer sm.Unlock()
    sm.m[key] = val
}
// + Load, Delete
```

## Testing

### Table-driven
```go
func TestX(t *testing.T) {
    cases := []struct{ give int; want bool }{
        {19, true}, {21, false}, {0, false}, {-1, false},
    }
    for _, c := range cases {
        if got := F(c.give); got != c.want {
            t.Errorf("F(%d) = %t, want %t", c.give, got, c.want)  // FUNKSİYA(p) = actual, want expected
        }
    }
}
```

### HTTP handler
```go
req, _ := http.NewRequest("GET", "/hello", nil)
r := httptest.NewRecorder()
handler := http.HandlerFunc(helloHandler)
handler.ServeHTTP(r, req)      // r.Code / r.Body.String() yoxla
```

### Mock (öz interfeys)
```go
type randIntGenerator interface{ Intn(int) int }
type fixedRandIntGenerator struct{ randomNum, calledWithN int }
func (g *fixedRandIntGenerator) Intn(n int) int {
    g.calledWithN = n          // çağırışı QAYD ET
    return g.randomNum          // istənilən "random"
}
```

### Coverage / Examples
```bash
go test -cover                          # 62.5%
go test -coverprofile=coverage.out
go tool cover -html=coverage.out        # sətir-sətir brauzerdə
```
```go
func ExampleUsername() {
    ...
    // Output:        ← go test yoxlayır
    // "gopher": true
}
// Ad: ExampleF / ExampleT_M / Example_suffix
```

## Benchmarks

```go
func BenchmarkFastF(b *testing.B) {
    m := len(numbers) - 1
    b.ResetTimer()             // setup-dan SONRA
    b.SetBytes(int64(n))       // bayt/s hesabı
    for n := 0; n < b.N; n++ {
        FastF(nums[n&m])       // %m YOX &m (2^n) — 2x dəqiqlik!
    }
}
```
```bash
go test -bench=. -benchmem
# ns/op    B/op    allocs/op
go test -bench . > old.txt; ... > new.txt; benchcmp old.txt new.txt
go test -benchmem -gcflags=-m          # escape analysis
go tool compile -S file.go | grep ...  # ASM
```

**Fibonacci dərsi:** recursive 1,255,534 ns/op → iterativ (a, b = b, a+b) 20.3 ns/op (~60,000x).

**n % m == n & (m-1)** — m 2-nin qüvvəti olanda; benchmark overhead azaldır:
```
Modulo:      16.7 ns/op
BitwiseAnd:   7.40 ns/op
```

## Race detector

```bash
go test -race       # test
go run -race        # run — production ssenarisi
go build -race
```
```
WARNING: DATA RACE
Write at 0x... by goroutine 7:   SetNoise()  cat.go:14   ← DƏQİQ SƏTİR
Previous read ... goroutine 6:   updateCat() cat.go:19
```
Həll: struct-a Mutex embed / handler-da WaitGroup ilə fon işi gözlə.

## Security

```go
// HSTS wrapper — BÜTÜN handlerlara tətbiq et:
func headerWrap(h http.HandlerFunc) http.HandlerFunc {
    return http.HandlerFunc(func(w, r) {
        w.Header().Add("Strict-Transport-Security",
            "max-age=31536000; includeSubDomains")
        h.ServeHTTP(w, r)
    })
}
```
- CSRF: per-session token → nosurf / gorilla/csrf
- CSP: mənbə məhdudiyyəti (XSS) → bluemonday, unrolled/secure
- golang-announce `[security]` elanlarını izlə

## CI

### .travis.yml
```yaml
language: go
go:
  - 1.11.x
  - tip
install:
  - make install
script:
  - make lint
  - make test
```

### Makefile
```makefile
all: lint build test
build:
	go build ./...
lint:
	gometalinter --exclude=vendor --disable-all --enable=golint --enable=vet --enable=gofmt ./...
	find . -name '*.go' | xargs gofmt -w -s
test:
	go test -cover ./check ./handlers
```
