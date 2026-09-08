# Effective Go Recipes — Cheat Sheet (AZ)

## I/O (ch1)

### `bytes.NewReader(data)` / `strings.NewReader(s)` / `bytes.Buffer`
[]byte/string → io.Reader; Buffer → yaddaş içi Writer.

### `gzip.NewWriter(w)` + `w.Name/ModTime` + `io.Copy`
Sıxışdırma zənciri — metadata Copy-dən ƏVVƏL.

### `os.DirFS(dir)` + `fs.Glob(root, "*.log")`
Fayl sistemi abstraksiyası ilə pattern axtarışı.

### `sha1.New()` + `io.Copy(w, r)`
Hash = io.Writer — faylın imzası bir sətirdə.

### `unix.Mmap(fd, 0, size, PROT_READ, MAP_PRIVATE)` + `bytes.Index`
Böyük fayl → []byte; axtarış yaddaş həddindən asılı deyil.

## Serializasiya (ch2)

### `gob.Register(&T{})` + `enc.Encode(&evt)`
İnterfeys pointeri encode; eyni encoder axın boyu.

### Anonim struct: `var reply struct { Payments []struct { Amount float64 } }`
Yalnız lazımi sahələr — json bilinmişi ignore edir.

### `dec.Decode(&l)` döngüsü + `io.EOF`
Streaming JSON (JSON lines).

### Pointer sahələr + `if l.Time == nil`
İtkin JSON sahələrini yaxala.

### `func (s *T) MarshalJSON()` / `UnmarshalJSON`
Custom tip → ən yaxın JSON tipi; `*s = *node` (pointer update).

### `mapstructure.Decode(obj, &l)`
map[string]any → struct (dinamik "type" mesajları).

## HTTP (ch3)

### `url.Values{}` + `url.PathEscape` + `query.Encode()`
URL-i əllə QURMA — həmişə bu yolla.

### `http.NewRequestWithContext(ctx, "GET", url, nil)` + `io.LimitReader(body, 10MB)`
Timeout + DoS qoruması.

### `srv := http.Server{ReadTimeout, WriteTimeout}` + `srv.ListenAndServe()`
Server timeout-ları.

### `versionRouter` + ayrı `v1Mux()/v2Mux()` + `X-API-ver` header
Paralel API versiyaları.

## Mətn (ch4)

### `%#v` (debug) / `%+v` (sahə adları) / `%v`
Verb seçimi; String() → sonsuz rekursiyaya diqqət.

### `mime.ParseMediaType(ctype)` → `params["charset"]`; yoxsa `charset.DetermineEncoding(data, "text/plain")`

### `strings.EqualFold(a, b)`
Unicode case-insensitive — ToLower YOX.

### `norm.NFKC.String(s)` (golang.org/x/text/unicode/norm)
Giriş kənarlarında normalizasiya — "Kraków" bug-ının dərmanı.

## Funksiyalar (ch5)

### Registry: `map[string]Encoder` + `Register(mime, enc)` + init
Dispatch — handler dəyişmədən format əlavə.

### `NewServer(options ...func(*Server) error)` + `WithVerbose` + `WithPort(port)`
Functional options — closure arqument daşıyır.

### `//go:linkname local remote` + `_ "unsafe"`
Unexported hack — yalnız HACK şərhi + issue ilə.

## Əsas tiplər (ch6)

### `v, ok := m[k]` / `val, ok := <-ch` / `n, ok := i.(int)`
Comma-ok — 3 sahədə.

### `if len(sl) > 1024 && 2*len(sl) < cap(sl)` → kopyala
Stack Pop-da memory leak qoruması.

### `cs := make([]int, 0, len(values))`
Append-dən əvvəl capacity → 19 alloc → 1.

### `time.LoadLocation(name)` + `ParseInLocation(layout, q, loc)` + öz `parseDelta`/`roundTime`
Zone-lu vaxt parse.

## Struct/Interface/Generics (ch7)

### `if s, ok := out.(syncer); ok { log.s = s }` + `nopSyncer{}`
Ad hoc interfeys — "kiçik interfeys, güclü abstraksiya".

### `errWriter{http.ResponseWriter; statusCode}` + WriteHeader override
Mövcud interfeysi embed ilə sargıla.

### `type Number interface { ~int | ~float64 }` + `func Max[T Number]([]T) (T, error)` + `var zero T`
Generic funksiya — tipi kompilyator çıxarır (nil üçün açıq tip).

### `func UnmarshalJSON[T Request](data []byte, d *T) error`
Pointer olmayan arqent → COMPILE xətası.

## Xətalar (ch8)

### `return fmt.Errorf("context: %w", err)`
Hər yerdə wrap + kontekst.

### Named return + `defer func() { if e := recover(); e != nil { err = fmt.Errorf(...) } }()`
Xarici panic → error.

### `safelyGo(func(){...})` — daxilində defer/recover
Goroutine panic qoruması.

### `errors.Is(err, os.ErrNotExist)` — yoxsa real xəta
Zəncir yoxlaması.

## Konkurensiya (ch9-10)

### Fan-out/fan-in: `go func(u string){ ch <- info }(url)` + `for range urls { <-ch }`
Nəticəyə url/err daxil et; loop dəyişəni parametr kimi!

### `pool := make(chan bool, runtime.GOMAXPROCS(0))` + `pool <- true` / `defer <-pool`
CPU-bound semafor.

### `wg.Add(n)` / `defer wg.Done()` / `wg.Wait()` + `errors.Join`
Qrup + xəta birləşdirmə.

### `select { case <-outCh: ... case <-ctx.Done(): return ctx.Err() }`
Timeout yarışı; ctx İLK parametr; HTTP-yə NewRequestWithContext.

### `type ctxKey string` + `Values struct` + `context.WithValue(ctx, valuesKey, &values)`
Kontekstdə logger/ID — toqquşmasız açar.

### `cfgOnce.Do(func(){...})`
Idempotent bir dəfəlik.

### `RWMutex`: `RLock/RUnlock` (oxu — paralel) / `Lock/Unlock` (yazı — tək)

### `go test -race` + `atomic.AddInt64(&counter, 1)` / `atomic.Value{}` (Store/Load + cast)
Race tut; atomik düzəlt (benchmark sübutu ilə!).

## Sockets/cgo (ch11-12)

### `net.Listen("tcp", addr)` + `ln.Accept()` + `go handler(c)`
Server; `io.CopyN(file, c, size)` — dəqiq ölçülü kopya.

### `net.Dial("unix", socketFile)` + json.NewDecoder/NewEncoder
Unix socket RPC.

### `binary.Write/Read(conn, binary.BigEndian, msg)`
UDP binary protokolu (NTP).

### `exec.Command("ping", "-c", n, host)` + `cmd.Output()` + Scanner parse
Xarici əmr + çıxış parse; CommandContext ilə timeout.

### `cmd.StdinPipe()/StdoutPipe()` + `cmd.Start()` + `defer c.p.Kill()`
Uzunömürlü proses (bc kalkulyator).

### cgo: `import "C"` + `C.CString` (free!) + `C.CBytes` + `C.GoBytes` + `#cgo LDFLAGS`
C yaddaşını əl ilə idarə et.

## Test (ch13)

### `inCi = os.Getenv("CI") != ""` + `if !inCi { t.Skip(...) }`
CI-yalnız testlər.

### YAML: `loadCases(t, path)` (t.Fatal) + `t.Run(tc.Name, ...)`
Table-driven; helper-lər error yox, t.Fatal qaytarır.

### `f.Add(seed)` + `f.Fuzz(func(t, n){...heuristic...})` + `-fuzztime 10s`
Fuzzing; uğursuz nümunə testdata/fuzz-da qalır.

### `MockTransport{body, err}` → `c.c.Transport = ...`
HTTP klient mock (RoundTripper); httptest.NewRecorder().Result().

### `TestMain(m)` + `runTests(m)` (defer üçün) + `os.Exit`
Qlobal fixtures.

### build + `freePort` + `waitForServer` + `t.Cleanup(Kill)`
End-to-end server test — go run YOX, build!

### `analysis.Analyzer{Name, Run}` + `singlechecker.Main` + `analysistest.Run`
55 sətirlik linter.

## Build/Ship (ch14-15)

### `//go:embed user.sql` + `var userSQL string`; qovluq → `embed.FS`
Asset binary-də.

### `go build -ldflags=-X main.version=$TAG` (git tag yoxsa short commit)
Versiya inyeksiyası.

### `CGO_ENABLED=0 go build` (Alpine/musl üçün MƏCBURİ)
`ldd agent` → "not a dynamic executable".

### `//go:build prof` faylı + `_ "net/http/pprof"` + `go build -tags prof`
Şərti profiling.

### `.goreleaser.yaml` + `goreleaser build`
GOOS/GOARCH matrisi.

### `//go:generate go run _scripts/gen.go in.txt out.go` + `go generate ./...`
Kod generasiyası; generasiya nəticəsi VCS-də.

### `conf.Parse("APP", &cfg)` + `validateAddr` → Fatal
Konfiqurasiya + validasiya.

### `replace github.com/x/y => ./_patch/y`
Asılılıq patch-i.

### Multistage Docker: build (golang) → deploy (slim, root-suz, CGO=0)

### `signal.Notify(ch, SIGTERM, SIGINT)` (buffered!) + `<-ch` + `srv.Shutdown(ctx)`
Graceful shutdown.

### `if info := s.logger.Check(zap.InfoLevel, ...); info != nil { info.Write(...) }`
Bahalı parametrləri qoru.

### `expvar.NewInt("name.calls")` + `statusWriter` middleware
/debug/vars metrics.

### `sudo dlv attach PID` → `b func` / `c` / `n` / `q`
Canlı debug.
