# Go Recipes (Torres) — Cheat Sheet

## I/O (Ch 1)

```go
io.Copy(dst, src)                        // axını köçür
io.CopyBuffer(dst, src, buf)             // buffer ilə (böyük data)
io.MultiWriter(w1, w2)                   // iki yerə eyni anda
r, w := io.Pipe()                        // yaddaş konveyeri
bytes.NewBuffer(b) / bytes.NewBufferString(s)
strings.NewReader(s)                     // string → io.Reader
buf := make([]byte, 64)
ioutil.ReadAll(r)                        // Go 1.16+: io.ReadAll
// CSV:
csvr := csv.NewReader(r); csvr.Comma = ';'
rows, _ := csvr.Read()                    // record-record
csvw := csv.NewWriter(w); csvw.Flush(); csvw.Error()
// Temp:
t, _ := ioutil.TempDir("", "tmp"); defer os.RemoveAll(t)   // os.MkdirTemp
// Template:
t := template.New("x").Funcs(template.FuncMap{"split": strings.Split})
template.ParseGlob(dir + "/*.tmpl")       // bütün fayllar bir şablon
// html/template avtomatik escape edir (XSS qoruması)
```

## CLI (Ch 2)

```go
flag.StringVar(&s, "subject", "", "desc")
flag.Var(&customVal, "c", "custom tip: String()+Set()")
flag.Parse()                              // yalnız main-də
// Subcommand:
fs := flag.NewFlagSet("cmd", flag.ExitOnError)
fs.Usage = func() { fmt.Println(usage); fs.PrintDefaults() }
fs.Parse(os.Args[2:])
os.Getenv("X") / os.Setenv("X", "y")
// envconfig:
type C struct {
    V string `required:"true"`
    B bool   `default:"true"`
}
envconfig.Process("PREFIX", &c)
// Siqnal:
ch := make(chan os.Signal, 1)
signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
<-ch                                     // bloklanmış gözlə
```

## Data (Ch 3)

```go
i, err := strconv.ParseInt("123", 10, 64)
val, ok := i.(string)                    // comma-ok
switch v := x.(type) { case string: ... } // type switch
big.NewInt(0).Add(a, b)                  // dəyişən uzunluq ədədlər
// Pul — float YOX, string→sent:
strings.Split("15.93", ".")               // ["15","93"] → int64(1593)
// gob:
gob.NewEncoder(buf).Encode(&p)
gob.NewDecoder(buf).Decode(&p2)
// base64:
base64.URLEncoding.EncodeToString(b)
// JSON null-able:
Age *int `json:"age,omitempty"`           // nil vs 0 fərqi
type NI sql.NullInt64                     // + MarshalJSON/UnmarshalJSON
```

## Error handling (Ch 4)

```go
err := errors.Wrap(err, "kontekst")       // pkg/errors — nil-i də idarə edir
root := errors.Cause(err)                 // orijinal xəta
fmt.Printf("%+v", err)                   // stack trace
switch errors.Cause(err).(type) { ... }   // typed xəta yoxlaması
// Log yalnız SON təyinatda; yolda wrap et, loglama
var once sync.Once
once.Do(func() { init() })               // bir dəfəlik init
defer func() { if r := recover(); r != nil { log(r) } }()
```

## Şəbəkə (Ch 5)

```go
ln, _ := net.Listen("tcp", ":8888")
conn, _ := ln.Accept(); go handle(conn)   // goroutine per conn
net.Dial("tcp", addr)                     // client
// UDP:
addr, _ := net.ResolveUDPAddr("udp", addr)
conn, _ := net.ListenUDP("udp", addr)
conn.WriteToUDP(b, retAddr) / conn.ReadFromUDP(buf)
// DNS:
net.LookupCNAME(host); net.LookupHost(host)
// WebSocket (gorilla):
conn, _ := upgrader.Upgrade(w, r, nil)
conn.ReadMessage() / conn.WriteMessage(t, p)
c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(1000, ""))
// Mail:
m, _ := mail.ReadMessage(r); m.Header.Get("To"); header.Date()
```

## Baza (Ch 6)

```go
import _ "github.com/go-sql-driver/mysql"  // driver qeydiyyatı
db, _ := sql.Open("mysql", dsn)
db.SetMaxOpenConns(24); db.SetMaxIdleConns(24)
rows, _ := db.Query("SELECT ... WHERE name=?", name)
defer rows.Close()
for rows.Next() { rows.Scan(&a, &b) } rows.Err()
// Tx + ümumi DB interfeysi:
type DB interface { Exec(...) ; Query(...) ; QueryRow(...) ; Prepare(...) }
tx, _ := db.Begin()
defer tx.Rollback()
dbinterface.Exec(tx)                       // həm DB həm Tx qəbul edir
tx.Commit()
ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
db.BeginTx(ctx, nil)
// Redis:
client.Set("k", "v", 5*time.Second)
client.Get("k").Scan(&s)                   // redis.Nil = açar yoxdur
// Mongo:
coll.FindOne(ctx, bson.M{"name": x}).Decode(&s)
coll.InsertMany(ctx, vals)
// Storage interfeysi → backend dəyişməsi/mock bir sətirlə
```

## Veb klient (Ch 7)

```go
// Custom Transport (auth hər sorğuda):
type APITransport struct { *http.Transport; user, pass string }
func (t *APITransport) RoundTrip(r *http.Request) (*http.Response, error) {
    r.SetBasicAuth(t.user, t.pass)
    return t.Transport.RoundTrip(r)
}
// Async fetch:
respch := make(chan *http.Response, len(urls))
errch := make(chan error, len(urls))
for _, u := range urls { go func(u string) { ... c.AsyncGet(u) }(u) }
for i := 0; i < len(urls); i++ {
    select { case r := <-c.Resp: ...; case e := <-c.Err: ... }
}
// Decorate (middleware zənciri):
c.Transport = Decorate(&http.Transport{}, Logger(l), BasicAuth(u, p))
// gRPC:
conn, _ := grpc.Dial(addr, grpc.WithInsecure())
client.Greet(ctx, &Req{...})
```

## Veb server (Ch 8)

```go
// Controller DI:
type Controller struct { storage Storage }
func New(s Storage) *Controller { return &Controller{s} }
// Closure handler:
func (c *Controller) Get(useDefault bool) http.HandlerFunc { return func(w, r) {...} }
// Middleware:
type Middleware func(http.HandlerFunc) http.HandlerFunc
func Apply(h http.HandlerFunc, mws ...Middleware) http.HandlerFunc {
    applied := h
    for _, m := range mws { applied = m(applied) }
    return applied
}
// context value:
ctx = context.WithValue(r.Context(), ID, "100")
r = r.WithContext(ctx)
// Reverse proxy:
r.URL = parsed; r.Host = parsed.Host; r.RequestURI = ""
// gRPC kodu → JSON API: internal paket + http.Handler wrap
if grpc.Code(err) == codes.NotFound { w.WriteHeader(404) }
```

## Test (Ch 9)

```go
// Closure mock:
type MockX struct{ Do func(...) error }
// Patch (yalnız var A = func(){} üçün):
defer Patch(&fn, fake).Restore()
// gomock:
ctrl := gomock.NewController(t); defer ctrl.Finish()
m.EXPECT().Get("key").AnyTimes().Return(v, err)
m.EXPECT().Get(gomock.Any()).Do(func(k string){}).Return("", nil)
// Table-driven:
tests := []struct{ name string; args args; wantErr bool }{...}
for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) { ... })
}
// go test -bench=. -benchmem -coverprofile=c.out
// httptest:
w := httptest.NewRecorder()
r := httptest.NewRequest("POST", "/x", body)
Handler(w, r)                             // network-süz test
```

## Concurrency (Ch 10)

```go
for { select { case <-ctx.Done(): return; case v := <-ch: ... } }  // əsas pattern
wg.Add(1); go func(v T) { defer wg.Done(); ... }(v)  // loop dəyişəni ARQUMENT!
wg.Wait()
atomic.AddInt64(&n, 1); atomic.LoadInt64(&n)          // sayğac → mutex YOX
mu.Lock()/Unlock(); mu.RLock()/RUnlock()              // map → RWMutex
once.Do(init)                                          // bir dəfəlik
ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
defer cancel()
// Worker pool:
for i := 0; i < num; i++ { go Worker(ctx, i, in, out) }
// Pipeline: pool-A çıxışı = pool-B girişi (kanallarla qoş)
```

## Distributiv (Ch 11)

```go
// Consul:
reg := &api.AgentServiceRegistration{ID: n, Name: n, Port: p, Address: a}
client.Agent().ServiceRegister(reg)
client.Health().Service("name", "tag", false, nil)
// Raft FSM: Apply(r *raft.Log); raft.Apply(data, timeout) — leader-də
// Docker: alpine + statik binary + -ldflags "-X main.version=1.0"
// Compose: servis adı = host adı; go mod vendor → offline build
// Prometheus:
http.Handle("/metrics", promhttp.Handler())
// go-metrics:
c := metrics.GetOrRegisterCounter("ad", nil); c.Inc(1)
t := metrics.GetOrRegisterTimer("ad", nil); t.UpdateSince(start)
```

## Reaktiv (Ch 12)

```go
// Goflow:
u := goflow.NewGraph()
u.Add("a", new(A)); u.Connect("a", "Out", "b", "In")
u.MapInPort("In", "a", "Val")
wait := flow.Run(net); close(in); <-wait   // təmiz shutdown
// Kafka sync:
p, _ := sarama.NewSyncProducer(brokers, nil)
p.SendMessage(&sarama.ProducerMessage{Topic: t, Value: sarama.StringEncoder(v)})
// Kafka async:
cfg := sarama.NewConfig(); cfg.Producer.Return.Successes = true
p.Input() <- msg                          // handler bloklanmır
go func() { for { select { case <-p.Successes(): ...; case <-p.Errors(): ... } } }()
pc, _ := consumer.ConsumePartition("topic", 0, sarama.OffsetNewest)
// GraphQL:
graphql.NewObject / graphql.NewSchema / graphql.Do(graphql.Params{...})
p.Args["suit"].(string); p.Source.(Card)
```

## Performans (Ch 14)

```go
import _ "net/http/pprof"                 // /debug/pprof/* yaranır
go tool pprof http://host/debug/pprof/profile    // 30s CPU profili
go tool pprof http://host/debug/pprof/heap       // yaddaş profili
// Benchmark:
func BenchmarkX(b *testing.B) { for n := 0; n < b.N; n++ { X() } }
b.RunParallel(func(pb *testing.PB) { for pb.Next() { X() } })
go test -bench=. -benchmem                // ns/op, B/op, allocs/op
// Qaydalar: sayğac → atomic; string birləşdirmə → strings.Join
```
