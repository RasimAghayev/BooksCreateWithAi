# Learning Go, Second Edition — Cheat Sheet (Azərbaycanca)

Jon Bodner, O'Reilly 2024 — 15 fəsil, L3 Intermediate

## 1. Alətlər
```bash
go version                          # quraşdırma yoxlaması
go run hello.go                     # script kimi — binary saxlamır
go build -o hello_world hello.go    # paylanacaq binary
go install github.com/rakyll/hey@latest   # alət qur (GOPATH/bin)
goimports -l -w .                   # import-ları təmizlə + formatla
go vet ./...                        # şübhəli kod
golint ./... / golangci-lint run    # stil + lint dəsti
shadow ./...                        # kölgələmə yoxlaması
go doc fmt                          # sənəd oxu
```

## 2. Tiqlər və Bəyanatlar
```go
var x int = 10; var x = 10; var x int   // 3 var forması
x := 10                                 // funksiya daxilində
x, y := 10, "hi"                        // ən azı biri yeni olmalı
const x = 10                            // untyped (çevik) — üstünlük
type Score int; type HighScore Score    // fərqli tiplər! conversion lazım
rune == int32; byte == uint8             // alias-lar
float64(x) + y                          // açıq conversion MÜTLƏQ
```

## 3. Slice-lar
```go
x := []int{1,2,3}; var x []int          // nil slice — append OK
x = append(x, 4); x = append(x, y...)
make([]int, 5)                          // [0 0 0 0 0] — append SONA yazır!
make([]int, 0, 10)                      // cap 10 + append — ən təhlükəsiz
y := x[:2:2]                            // 3-hissəli — append paylaşımını kəs
copy(dst, src)                           // müstəqil nüsxə
```

## 4. Kontrol Strukturları
```go
if n := f(); n == 0 { } else if n > 5 { }   // scoped bəyanat
for i := 0; i < 10; i++ { }                   // tam
for cond { }                                   // while
for { }                                        // sonsuz (break ilə)
for k, v := range m { }                        // map — TƏSADÜFİ sıra!
for i, r := range s { }                        // string — RUNE üzrə, i = bayt
switch x := f(); x { case 1,2: ; default: }    // fall-through YOX
switch { case a < 5: }                          // blank switch
outer: for { continue outer }                  // label
```

## 5. Funksiyalar
```go
func f(a, b int) (int, error) { }
func f(vals ...int) { }          // variadic; slice: f(xs...)
v, ok := m[k]                    // comma-ok: map, channel, type assertion
defer f.Close()                  // LIFO; adlı return + defer = commit/rollback
defer func() { err = wrap(err) }()   // ortaq wrap pattern
```
**QADAĞAN:** blank return; `var err CustomType` + return (həmişə non-nil interface!).

## 6. Metodlar və İnterfeyslər
```go
func (c *Counter) Inc()      // dəyişir/nil idarə edir → pointer MÜTLƏQ
func (c Counter) String()   // dəyişmir → value
type I interface { M() }     // client tərəfdə, kiçik təyin et
f := myAdder.AddTo            // method value
f2 := Adder.AddTo             // method expression — func(Adder, int) int
```
- Embed = kompozisiya, inheritance YOX (dynamic dispatch YOXDUR)
- Accept interfaces, return structs
- `i = typedNil` → `i == nil` FALSE (tip pointer-i non-nil)

## 7. Xətalar
```go
errors.New("msg"); fmt.Errorf("%d err: %w", n, err)   // %w = wrap
if errors.Is(err, os.ErrNotExist) { }                   // sentinel zəncirdə
var se StatusErr
if errors.As(err, &se) { }                              // tip zəncirdə
```
Sentinel: `Err` prefiksi, == ilə test, NADİR. Custom: həmişə `error` qaytar.
panic: yalnız library sərhəddində recover + error-a çevir.

## 8. Paralellik
```go
go func() { ... }()                       // closure wrap pattern
ch := make(chan int); ch := make(chan int, 10)
x := <-ch; ch <- x; v, ok := <-ch; close(ch)  // yazan bağlayır!
select { case v := <-ch1: ; case <-time.After(2s): ; default: }  // RANDOM seçim
```
- Done pattern: `done := make(chan struct{}); defer close(done)`
- Loop dəyişəni goroutine-a PARAMETR ötür
- WaitGroup: Add(n) qabaq + defer Done + Wait; kopyalama YOX
- sync.Once: lazy init; `once.Do(func(){...})`
- Mutex: `Lock(); defer Unlock()`; RWMutex oxucular üçün; reentrant DEYİL
- Timeout: select + time.After VƏYA context.WithTimeout + defer cancel

## 9. Modullar
```bash
go mod init path
go get mod@v1.2.3; go get -u=patch mod; go get -u mod
go list -m -versions mod
go mod tidy
go mod vendor
```
v2+: yol `/v2` bitir + yeni qovluq/branch + teq. Minimum version selection.
internal/ → valideyn+sibling üçün. go.mod + go.sum həmişə VCS-də.

## 10. Standart Kitabxana
```go
// io
n, err := r.Read(buf); if err == io.EOF {}   // əvvəl datanı emal et!
json.NewDecoder(r).Decode(&v); json.NewEncoder(w).Encode(v)
io.Copy, io.MultiReader, io.LimitReader, ioutil.NopCloser

// time
d := 2*time.Hour + 30*time.Minute
t1.Equal(t2)   // == YOX (timezone!)
time.Parse("2006-01-02 15:04:05", s)   // reference time şablonu
time.NewTicker(d)   // Tick YOX (GC-olunmaz)

// http
client := &http.Client{Timeout: 30 * time.Second}   // DefaultClient YOX
mux := http.NewServeMux()                          // DefaultServeMux YOX
middleware: func(http.Handler) http.Handler
```

## 11. Context
```go
ctx, cancel := context.WithTimeout(ctx, 50*time.Millisecond)
defer cancel()                          // MÜTLƏQ — leak!
<-ctx.Done(); ctx.Err()                 // Canceled / DeadlineExceeded
ctx = context.WithValue(ctx, key, v)    // unexported key tipi ilə!
req.Context(); req.WithContext(ctx)
```

## 12. Testlər
```go
func TestX(t *testing.T) { t.Error/Fatal(f) }
t.Run(name, func(t *testing.T) {})          // table test
cmp.Diff(expected, got)                      // go-cmp
func BenchmarkX(b *testing.B) { for i := 0; i < b.N; i++ {} }
```
```bash
go test -v ./... -count=1
go test -cover -coverprofile=c.out && go tool cover -html=c.out
go test -bench=. -benchmem
go test -race                                # həmişə paralel kod üçün!
go test -tags integration
```
Stub: böyük interface → embed; fərqli case-lər → funksiya sahəli stub.
httptest.NewServer — HTTP stub. testdata/ — nisbi yollar.

## 13. Reflection / unsafe / cgo (SON ÇARƏ)
- reflect: yalnız sərhəd (marshal/unmarshal); 30-70x yavaş
- unsafe: binary data + KeepAlive; checkptr ilə test
- cgo: yalnız əvəzsiz C kitabxanası (29x yavaş keçid)

## 14. Generics (Go 1.18)
```go
type Stack[T any] struct { vals []T }
type Stack[T comparable] struct { vals []T }   // == üçün
func Map[T1, T2 any](s []T1, f func(T1) T2) []T2
var zero T                                     // zero value hiyləsi
```

## Universal Qaydalar
1. Value > pointer (amma dəyişən receiver → pointer)
2. Interface-lər client tərəfdən, kiçik təyin olunur
3. Error həmişə sonda; xətanı yoxla ya `_` et
4. `:=` kölgələmə tələsi — xarici dəyişəni parametr ötür
5. Package-level mutable state YOX; init-dən qaç
6. Concurrency API-dən kənarda (closure wrap)
7. Hər goroutine-in çıxış yolu olsun (done/cancel)
8. Düzgün Go "sıxdırıcı"dır — aydınlıq qısalıqdan üstündür
