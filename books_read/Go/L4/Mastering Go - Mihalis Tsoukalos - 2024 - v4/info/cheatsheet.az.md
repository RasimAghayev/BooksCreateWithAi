# Mastering Go, Fourth Edition — Cheat Sheet (Azərbaycanca)

Mihalis Tsoukalos, Packt 2024 — 15 fəsil + GC əlavəsi, L4 Advanced

## 1. Əsas Sintaksis (Ch1)
```go
var x int; var x = 10; var x int   // var formaları
x := 10                              // funksiya daxilində
if err != nil { }                    // əsas pattern
switch { case v > 0: }               // ifadəsiz switch
for i := range 5 { }                 // Go 1.22: integer range
```
- `{` eyni sətirdə; implicit tip çevirməsi QADAĞA; unused var/import = compile xətası

## 2. Slice-lər (Ch2)
```go
s := make([]int, 5)          // len=cap=5
s := make([]int, 0, 10)      // pre-allocate — ən yaxşı
s = append(s, x); append(s, o...)   // nəticəni TƏYİN ET
t := s[2:5:10]               // len=3, cap=8
copy(dst, src)               // dst böyümür
slices.Delete(s, i, j)       // Go 1.22: silinən yerləri sıfırlayır
```

## 3. String/Rune (Ch2)
```go
r := '€'                     // rune = int32 = 8364
string(100)                  // "d" (code point!) — İtoi/FormatInt istifadə et
utf8.ValidString(s)           // validasiya ƏVVƏL; []rune(s) ilə manipulyasiya
```

## 4. Generics (Ch4)
```go
func Print[T any](s []T) { }
func Same[T comparable](a, b T) bool { return a == b }
type Numeric interface { int | int64 | float64 }   // union
type AllInts interface { ~int }                    // supertype — alias-ları da qəbul
type Stack[T any] struct { vals []T }
// slices.Concat, slices.Sort, maps.Copy, cmp.Compare (Go 1.21+)
```

## 5. İnterfeyslər (Ch5)
```go
type I interface { M() }
// implicit satisfaction — implements YOXDUR
v, ok := i.(T)              // type assertion (güvənli forma)
switch T := x.(type) { }    // type switch
_, ok := interface{}(a).(Shape2D)   // satisfy yoxlaması
```
- error = `Error() string` interfeysi; io.EOF xəta DEYİL

## 6. Funksiyalar (Ch6)
```go
func f() (min, max int) { return }        // adlı return
func f(vals ...int) {}                    // variadic; f(slice...)
f2 := func(x int) int { return x * x }    // anonim
sort.Slice(data, func(i, j int) bool { return data[i].G < data[j].G })
// defer LIFO; closure-də dəyər PARAMETR kimi:
defer func(n int) { fmt.Print(n) }(i)     // DÜZGÜN; func(){fmt.Print(i)}() = TƏLƏ (1.21-)
```

## 7. Fayl I/O (Ch7)
```go
f, _ := os.Open(file); defer f.Close()
r := bufio.NewReader(f)
line, err := r.ReadString('\n')          // sətir-sətir; io.EOF-da son sətri emal et
w := bufio.NewWriter(f); w.WriteString(s); w.Flush()   // Flush MÜTLƏQ
f, _ = os.OpenFile(p, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
os.ReadFile(p); os.WriteFile(p, data, 0644)   // ioutil ƏVƏZİNİ (1.16+)
//go:embed static/file
var content string                        // embed
```

## 8. JSON (Ch7)
```go
json.Marshal(&v)                 // struct → JSON
json.Unmarshal([]byte(s), &v)     // JSON → struct (SAHƏLƏR EXPORTED!)
`json:"name,omitempty"`          // boş → düşmür
`json:"-"`                        // tam gizli
json.MarshalIndent(v, "", "\t")   // pretty
json.NewEncoder(w).Encode(v)     // stream; encoder BİR DƏFƏ yarat
```

## 9. Cobra/Viper (Ch7,9,11)
```go
rootCmd.PersistentFlags().StringP("server", "S", "localhost", "...")
cmd.Flags().StringP("user", "u", "", "...")       // lokal
viper.BindPFlag("server", rootCmd.PersistentFlags().Lookup("server"))
viper.GetString("server")
viper.SetConfigType("json"); viper.ReadInConfig()  // fayl yoxdursa SƏSSİZ keçir!
viper.Unmarshal(&t)                                 // mapstructure tag-lərlə
```

## 10. Paralellik (Ch8)
```go
var wg sync.WaitGroup
for ... {
    wg.Add(1)                        // go-dan ƏVVƏL
    go func(x int) { defer wg.Done(); ... }(x)   // dəyər PARAMETR
}
wg.Wait()

ch := make(chan int); make(chan int, 5)   // unbuffered/buffered
close(ch)                                  // YALNIZ SENDER
for v := range ch { }                      // close-dan sonra çıxır
select {
case v := <-ch1:
case <-time.After(d):                      // timeout
case <-ctx.Done():
default:                                   // non-blocking
}
```
| Əməliyyat | Nəticə |
|---|---|
| bağlı kanala YAZ | panic |
| bağlıdan OXU | zero value |
| nil kanal | blok |
| nil close | panic |

```go
var m sync.Mutex
m.Lock(); defer m.Unlock()            // kritik bölmə
rl.RLock(); defer rl.RUnlock()        // RWMutex oxu — paralel
atomic.AddInt64(&v, 1)                // sayğac
ctx, cancel := context.WithTimeout(ctx, d); defer cancel()
sem := semaphore.NewWeighted(n)       // x/sync/semaphore
sem.Acquire(ctx, 1); defer sem.Release(1)
```
- `go run -race` HƏMİŞƏ; Add/Done balansı; `go test -shuffle=on`

## 11. HTTP Server (Ch9,11)
```go
mux := http.NewServeMux()                   // default mux-dan QAÇIN
mux.Handle("/x", http.HandlerFunc(h))       // mux.HandleFunc qısa yol
s := &http.Server{Addr: PORT, Handler: mux,
    ReadTimeout: 3*time.Second, WriteTimeout: 3*time.Second}  // DoS qorunması
s.ListenAndServe(); s.Shutdown(nil)          // graceful
// gorilla/mux (Go 1.21-): r.Methods(GET).Subrouter(); {id:[0-9]+}
// Go 1.22+: mux.HandleFunc("GET /users/{id}", h) — standart!
r.PathValue("id")                            // 1.22 wildcard dəyəri
```

## 12. HTTP Client (Ch9,11)
```go
req, _ := http.NewRequest(http.MethodGet, url, bytes.NewReader(body))
req.Header.Set("Content-Type", "application/json")
c := &http.Client{Timeout: 15 * time.Second}
resp, err := c.Do(req); defer resp.Body.Close()
if resp.StatusCode != http.StatusOK { }       // STATUS YOXLA
// context timeout:
req, _ = http.NewRequestWithContext(ctx, GET, url, nil)
```

## 13. TCP/UDP/WS (Ch10)
```go
c, _ := net.Dial("tcp", "host:port")         // client; c = Reader+Writer
l, _ := net.Listen("tcp", ":1234")
for { c, _ := l.Accept(); go handleConn(c) } // CONCURRENT pattern
n, addr, _ := conn.ReadFromUDP(buf); conn.WriteToUDP(data, addr)
// WebSocket:
upgrader.Upgrade(w, r, nil) → ws.ReadMessage()/ws.WriteMessage() (YALNIZ bu API!)
websocket.DefaultDialer.Dial(url, nil)
c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(...))  // graceful
```

## 14. Test (Ch12-13)
```go
func TestX(t *testing.T) { t.Error / t.Errorf / t.Fatal(f) }
func BenchmarkX(b *testing.B) { for i := 0; i < b.N; i++ { fxed() } }
func FuzzX(f *testing.F) { f.Add(seed); f.Fuzz(func(t, ...){ }) }
func ExampleX() { /* Output:\n 7 */ }
b.Run(name, func(b *testing.B) { })         // sub-benchmark
t.TempDir(); t.Cleanup(fn)                  // avtomatik təmizlik
```
```bash
go test -v -count=1 -timeout 30s -shuffle=on
go test -bench=. -benchmem
go test -coverprofile=c.out && go tool cover -html=c.out
go test -fuzz=FuzzX -fuzztime 10s
go vet ./...; govulncheck ./...
go test -race
```
- httptest: `rr := httptest.NewRecorder(); handler.ServeHTTP(rr, req)`
- mux.SetURLVars(req, map) — path dəyişənləri testdə

## 15. Profiling/Monitorinq (Ch12-13)
```go
pprof.StartCPUProfile(f); defer pprof.StopCPUProfile()
pprof.WriteHeapProfile(f)
import _ "net/http/pprof"                   // /debug/pprof
trace.Start(f); defer trace.Stop()
metrics.Read([]metrics.Sample{{Name: "/sched/goroutines:goroutines"}})
prometheus.MustRegister(prometheus.NewGauge(prometheus.GaugeOpts{Namespace: "ns", ...}))
http.Handle("/metrics", promhttp.Handler())
```
```bash
go tool pprof cpu.out        # top / top10 -cum / list F / pdf / -http=:1234
go tool trace trace.out
GODEBUG=gctrace=1 go run x.go
```

## 16. Performans (Ch14)
```go
buf := new(bytes.Buffer)         // BİR DƏFƏ yarat
for { w(buf); buf.Reset() }      // Reset = 0 alloc (150ns vs 1056ns)
mySlice := make([]int, 0, n)     // pre-allocate
// map leak: delete azaltmır → m = nil
// subslice leak: copy(make([]T, k), s) ilə kopla
```
```bash
go run -gcflags '-m' x.go       # escape analizi
goos=linux goarch=amd64 go build  # cross-compile
```

## 17. Go 1.21/1.22 Yenilikləri (Ch15)
```go
sync.OnceFunc(f)()      // tək icra
clear(m)                 // map sil / slice sıfırla
for x := range 5 { }     // integer range
rand.N(100)              // math/rand/v2 — generic
```
- Loop dəyişəni artıq paylaşılır DEYİL (goroutine-safe)
- slices.Delete family silinənləri sıfırlayır

## Universal Qaydalar
1. Sadəliyi qoru — abstraksiya haqq qazansın
2. Error nil qədər yoxla; xətanı bağlama (wrap) və yuxarı ötür
3. Body/Close defer; context Cancel defer
4. Concurrency-də -race; benchmark-da -benchmem
5. Data strukturu = performans: slice > map; dəyər > pointer
6. Əvvəl işləsin, sonra gözəl, yalnız lazımsa sürətli
