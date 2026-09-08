# Build systems with Go — Cheat Sheet (AZ)

## Go əsasları (ch1-2)

### `go build main.go`
Kompilyasiya edib icra olunan fayl yaradır.
- nəticə: `main` binary; `./main` ilə icra

### `go test [-v] [-run regex] [-bench .] [-cover]`
Test/benchmark icra aləti.
- `-v` → hər test ayrıca; `-run /Encoding=JSON` → subtest regex
- `-bench .` → benchmark-lar; `-cpu 1,2,4` → goroutine matrisi
- `-cover` → örtük faizi; `-coverprofile=f` + `go tool cover -html=f`

### `go mod init` / `go get paket` / `go mod vendor`
Asılılıq idarəetməsi.
- build/run/test asılılıqları avtomatik go.mod-a yazır
- vendor YALNIZ xüsusi ehtiyacda; repo-ya commit etməyin

### `var x T` / `x := v` / `const` / `iota`
Dəyişən/konstant elanları.
- `:=` tip çıxarır; `iota` enum-lara ardıcıl dəyər verir

### `func f(a, b int) (int, error)` / `nums ...int` / closure
Funksiya formaları.
- çoxdəyərli qaytarma; variadic; `func() int` — closure qaytaran funksiya

### `defer` / `panic` / `recover`
- defer → LIFO təmizlik; recover → yalnız deferred daxilində panic tutur

## Slice/Map (ch3)

### `make([]T, len, cap)` / `append(s, v...)`
Slice yaratma/böyütmə.
- cap avtomatik artır; `len`-dən böyük index yazmaq runtime xətadır

### `v, ok := m[k]` / `delete(m, k)`
Map əməliyyatları — comma-ok idiomu.

### `for i, v := range s`
İterasiya — `v` KOPYADIR; orijinalı dəyişmək üçün `s[i]`.

## Struct/Interface (ch4)

### `func NewX(...) (*X, error)`
Konstruktor konvensiyası — pointer + validasiya error-u.

### `func (r *T) M()` vs `func (r T) M()`
Pointer receiver orijinalı dəyişir; interfeysə yalnız *T uyğun gəlir.

### `type I interface { M() }`
Implicit implementation — implements açar sözü YOXDUR.

### `switch t := i.(type)`
Empty interface tip yoxlaması.

## Concurrency (ch6)

### `go f(x)`
Goroutine işə salma — main bitəndə hamısı dayanır.

### `ch := make(chan T)` / `make(chan T, n)`
Unbuffered (hər iki tərəf hazır) / buffered (n element bufer).

### `ch <- v` / `v := <-ch` / `close(ch)`
Göndər / qəbul / bağla; `v, ok := <-ch` → bağlanma yoxlaması.

### `for x := range ch`
Kanal bağlanana qədər istehlak.

### `select { case <-a: ... case <-b: ... default: ... }`
Çoxkanallı gözləmə; default → non-blocking.

### `var wg sync.WaitGroup` → `wg.Add(n)` / `defer wg.Done()` / `wg.Wait()`
Goroutine tamamlanma sayğacı.

### `ctx, cancel := context.WithTimeout(parent, d)` → `defer cancel()`
Vaxt məhdud kontekst; `<-ctx.Done()` ləğv siqnalı.
- loop daxilində cancel-i defer ETMƏYİN

### `sync.Once` → `once.Do(fn)`
Bir dəfəlik icra.

### `m.Lock()` / `m.Unlock()`
Mutex — kritik sahəni kiçik saxlayın.

### `atomic.AddInt32(&c, 1)` / `atomic.LoadInt32(&c)`
Atomik yazı/oxu.

## I/O (ch7)

### `ioutil.ReadFile(p)` / `ioutil.WriteFile(p, b, 0644)`
Fayl oxu/yaz (müasir Go: `os.ReadFile`/`os.WriteFile`).

### `f, _ := os.Create(p)` → `defer f.Close()` → `f.WriteString(s)` / `f.Seek(off, 0)`
Aşağı səviyyə fayl əməliyyatları; Seek(0,0) → əvvələ qayıt.

### `bufio.NewReader(r).ReadString('\n')` / `bufio.NewScanner(r).Scan()`
Buferli giriş; Scanner sətir-sətir.

### `bufio.NewWriter(w).WriteByte(b)` + `w.Flush()`
Buferli çıxış — Flush-unutmayın.

## Encoding (ch8)

### `json.Marshal(v)` / `json.Unmarshal(b, &v)`
- tag: `json:"ad,omitempty"`; unexported sahələr düşmür
- XML: map marshal xətasız; custom `MarshalXML` yazın
- YAML: `go get gopkg.in/yaml.v2` → `yaml.Marshal/Unmarshal`

## HTTP (ch9)

### `http.Get(url)` → `defer resp.Body.Close()`
Sadə GET — body Reader-dir, bağlanmalıdır.

### `http.Post(url, "application/json", buf)`
POST — body io.Reader.

### `req, _ := http.NewRequest(m, url, body)` + `client.Do(req)`
Custom metod + header; `http.Client{Timeout: ...}`.

### `http.HandleFunc("/yol", h)` + `http.ListenAndServe(":8090", nil)`
Server; h imzası: `func(ResponseWriter, *Request)`.

### Middleware pattern: `func(http.Handler) http.Handler`
Zəncir: `AuthMiddleware(LogMiddleware(handler))`.

## Test (ch11)

### `func TestX(t *testing.T)` + `t.Error(...)`
Test funksiyası; `t.Skip()`, `t.Run("A=1", sub)`.

### `func BenchmarkX(b *testing.B)` → `for i := 0; i < b.N; i++`
+ `b.RunParallel(func(pb *testing.PB){ for pb.Next() {...} })`

### `func ExampleT_M()` + `// Output:`
Sənədləşən, çıxışı yoxlanan nümunə.

### `go tool pprof cpu.out` → `top` / `list Funksiya`
Profil analizi; `-cpuprofile`/`-memprofile` ilə topla.

## gRPC (ch13-14)

### `protoc -I=. --go_out=... --go-grpc_out=... *.proto`
Stub generasiyası (mesajlar + server/client).

### `grpc.NewServer()` + `pb.RegisterXSrv(s, impl)` + `s.Serve(lis)`
Server; RPC imzası: `func(ctx, req) (*Resp, error)`.

### `grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())` + `NewXClient(conn)`
Klient.

### Streaming: `stream.Send/Recv` / `CloseAndRecv` / `SendAndClose` / `CloseSend`
- `io.EOF` → axın sonu

### Annotations + grpc-gateway → REST
`option(google.api.http) = { get: "/v1/user/{id}" };`

### `grpc.UnaryInterceptor(fn)` / `grpc.WithUnaryInterceptor(fn)`
Server/klient interceptor; `metadata.AppendToOutgoingContext`.

## CLI (ch16)

### `var root = &cobra.Command{Use, Short, Run}` + `root.Execute()`
- flag: `root.Flags().StringVar(&v, "msg", def, "help")`
- required: `root.MarkFlagRequired("msg")`
- persistent (qlobal): `root.PersistentFlags()`
- subcommand: `root.AddCommand(cmd)`
- completion: `cmd.Root().GenBashCompletion(os.Stdout)`
- docs: `doc.GenMarkdownTree(root, ".")`

## SQL (ch17)

### `db, _ := sql.Open("sqlite3", dsn)` + `defer db.Close()` + `db.PingContext(ctx)`
Driver blank import ilə; Open bağlantı AÇMIR — Ping yoxlayır.

### `db.ExecContext(ctx, q, args...)` → `Result.LastInsertId/RowsAffected`
Dəyişdirici; `db.QueryContext` → `rows.Next()/Scan/Close`.

### `tx, _ := db.BeginTx(ctx, nil)` → `defer tx.Rollback()` → `tx.Commit()`
Transaksiya.

### `stmt, _ := db.Prepare(q)` → `defer stmt.Close()` → `stmt.QueryRow(args)`
Prepared statement.

### GORM: `gorm.Open(...)` → `AutoMigrate(&T{})` → `Create/First/Find/Where/Preload`
- `db.Where("name = ?", "x").First(&u)` — sorğu dəyişənini sıfırlayın!
- `db.Transaction(func(tx *gorm.DB) error {...})` — return nil=commit

## Kafka (ch19)

### Confluent: `NewProducer(ConfigMap)` → `p.Produce(msg, deliveryChan)` + `p.Flush(ms)`
Consumer: `NewConsumer` → `Subscribe` → `Poll(ms)` → type switch → `Commit()`.

### Segmentio: `kafka.DialLeader(ctx, "tcp", addr, topic, part)`
+ `conn.WriteMessages(kafka.Message{...})` / `conn.ReadBatch(min, max)`.

### Reader/Writer: `kafka.NewReader(ReaderConfig{Brokers, Topic, GroupID})`
→ `r.FetchMessage(ctx)` + `r.CommitMessages(ctx, m...)` (variadic — batch!).

### REST Proxy: POST /topics/t/partitions/n (Content-Type: vnd.kafka.json.v2+json)
Consumer: POST /consumers/qroup → subscription → GET records → DELETE.
