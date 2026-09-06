# Go in Practice — Cheat Sheet (bütün kitabdan toplanmış — 70 texnika)

Matt Butcher, Trevor Matteson — Manning, 2016

## CLI (T1-2)

### flag paketi
```go
var name = flag.String("name", "World", "help")      // pointer qaytarır
var spanish bool
flag.BoolVar(&spanish, "spanish", false, "help")     // mövcud dəyişənə
flag.Parse()                                          // dəyərlər düşür
*name                                                 // dereference
flag.PrintDefaults() / flag.VisitAll(fn)              // help
flag.Args() / flag.Arg(i)                            // flag-olmayan arqumentlər
```

### go-flags (struct tag əsaslı)
```go
var opts struct {
    Name string `short:"n" long:"name" default:"World" description:"..."`
}
flags.Parse(&opts)
```

### cli.go (urfave)
```go
app := cli.NewApp()
app.Flags = []cli.Flag{cli.StringFlag{Name: "name, n", Value: "World", Usage: "..."}}
app.Action = func(c *cli.Context) error {
    name := c.GlobalString("name")     // komanda flag-ı: c.String
    return nil
}
app.Commands = []cli.Command{{Name: "up", ShortName: "u", Action: func...}}
app.Run(os.Args)
```

## Konfiq (T3-4)
```go
// JSON:
decoder := json.NewDecoder(file); decoder.Decode(&conf)
// YAML (go-gypsy):
config, _ := yaml.ReadFile("conf.yaml"); config.Get("path") / GetBool(...)
// INI (gcfg):
gcfg.ReadFileInto(&config, "conf.ini")     // config.Section.Enabled
// ENV (12-factor):
os.Getenv("PORT")                           // yoxdursa ""
```

## Web server (T5-9)

### Graceful shutdown (manners)
```go
ch := make(chan os.Signal)
signal.Notify(ch, os.Interrupt, os.Kill)
go func() { <-ch; manners.Close() }()
manners.ListenAndServe(":8080", handler)
```

### Router variantları
```go
// path wildcard (pathResolver):
pr.Add("GET /hello", hello); pr.Add("* /goodbye/*", goodbye)  // path.Match
// Regex:
rr.Add("(GET|HEAD) /goodbye(/?[A-Za-z0-9]*)?", goodbye)       // compile cache!
// 3rd-party: httprouter / gorilla/mux / pat
```

## Concurrency (T10-15)

### Goroutine + closure + WaitGroup
```go
var wg sync.WaitGroup
for i, file := range files {
    wg.Add(1)
    go func(filename string) {          // LOOP DƏYİŞƏNİNİ PARAMETR ET!
        compress(filename)
        wg.Done()
    }(file)
}
wg.Wait()
```

### Mutex (struct-a embed)
```go
type words struct {
    sync.Mutex                        // anonim → words.Lock() açılır
    found map[string]int
}
func (w *words) add(word string, n int) {
    w.Lock()
    defer w.Unlock()
    ...
}
```

### Kanallar
```go
c := make(chan []byte)                  // unbuffered
out chan<- []byte                       // send-only
done <- true                            // done-kanal patternı
select { case m := <-msg: ...; case <-time.After(30s): os.Exit(0) }   // timeout
lock := make(chan bool, 1)              // buffered(1) = KİLİD
lock <- true; ...; <-lock               // acquire / release
// close YALNIZ sender; bağlı kanal → zero value (sonsuz loop təhlükəsi!)
```

## Error/Panic (T16-21)
```go
// Error dəyişənləri (== müqayisəsi):
var ErrTimeout = errors.New("The request timed out")
if err == ErrTimeout { /* retry */ }

// Custom error:
type ParseError struct{ Message string; Line, Char int }
func (p *ParseError) Error() string { ... }

// Panic → error idiomu:
panic(errors.New("Something bad happened."))

// Recover:
defer func() {
    if err := recover(); err != nil {
        fmt.Printf("Trapped panic: %s (%T)\n", err, err)
    }
}()

// Panic-i error-a çevir:
defer func() {
    if r := recover(); r != nil {
        file.Close()
        err = r.(error)              // named return!
    }
}()

// safely.Go — goroutine panic trap:
func Go(todo GoDoer) {
    go func() {
        defer func() { if err := recover(); err != nil { log.Printf(...) } }()
        todo()
    }()
}
```

## Logging/Testing (T22-31)

### log
```go
logger := log.New(writer, "prefix ", log.LstdFlags|log.Lshortfile)
// Ldate/Ltime/Lmicroseconds; Llongfile/Lshortfile
logger.Fatalln(...)      // os.Exit(1) — DEFER-İ KEÇİR!
logger.Panicln(...)      // defer işləyir → network üçün BU
net.Dial("tcp", "localhost:1902")           // TCP (back pressure)
net.DialTimeout("udp", "localhost:1902", 30*time.Second)   // UDP (itki riski)
```

### syslog
```go
logger, _ := syslog.New(syslog.LOG_LOCAL3, "narwhal")
logger.Debug / Notice / Warning / Alert
syslog.Dial(...)        // remote
```

### Test
```go
func TestX(t *testing.T) { t.Error / t.Errorf / t.Fatal / t.Fatalf }
// Mock:
type Messager interface { Send(email, subject string, body []byte) error }
func Alert(m Messager, ...)                    // kod İNTERFEYS qəbul etsin
type MockMessage struct{ ... }                // yadda saxlayan implementasiya
// Canary (compile-time interfeys yoxlaması):
var _ io.Writer = &MyWriter{}
// Generative:
quick.Check(func(s string, max uint8) bool { ... }, &quick.Config{MaxCount: 200})
// Benchmark:
func BenchmarkX(b *testing.B) { for i := 0; i < b.N; i++ { ... } }
go test -bench . -cpu=1,2,4
b.RunParallel(func(pb *testing.PB) { for pb.Next() { ... } })
go test -race             // DATA RACE aşkarla
```

### Stack trace
```go
debug.PrintStack()
buf := make([]byte, 1024); runtime.Stack(buf, true)   // true = bütün goroutine-lər
```

## Templates (T32-38)
```go
// Custom funksiya:
t.Funcs(template.FuncMap{"dateFormat": fn})    // Parse-dən ƏVVƏL!
var t = template.Must(template.ParseFiles(...)) // PAKET səviyyəsində CACHE
// Buffer pattern:
var b bytes.Buffer
err := t.Execute(&b, p)
if err != nil { ...; return }       // yalnız tam səhifə göndər
b.WriteTo(w)
// Nested:
t.ExecuteTemplate(w, "index.html", p)    // {{template "head.html" .}}
// Inheritance:
{{define "base"}}...{{template "title" .}}{{block "styles" .}}default{{end}}...{{end}}
t["user.html"].ExecuteTemplate(w, "base", u)     // HƏMİŞƏ base!
// Safe HTML:
qc = template.HTML(b.String())          // escape-olmaz
// Email:
text/template + buffer + smtp.SendMail(...)
```

## Static/Formlar (T39-48)
```go
// Subdir:
http.StripPrefix("/static/", http.FileServer(http.Dir("./files/")))
// Cache + ServeContent:
mutex.RLock/RUnlock; mutex.Lock/Unlock       // RWMutex!
http.ServeContent(res, req, path, modTime, content)   // MIME/304/Length avtomatik
// Fayl upload:
f, h, err := r.FormFile("file")               // tək fayl
r.ParseMultipartForm(16 << 20)                 // 16 MB
files := r.MultipartForm.File["files"]          // çox fayl
fh.Open(); io.Copy(out, f)
// MIME yoxlaması:
buffer := make([]byte, 512); file.Read(buffer)
http.DetectContentType(buffer)                 // ƏN ETİBARLI
// Streaming (böyük fayl):
mr, _ := r.MultipartReader()
part, err := mr.NextPart()                     // io.EOF → son
part.FileName() == "" → text sahə              // fayl → chunk Read + Write
```

## REST (T49-55)
```go
// Timeout detection:
switch err := err.(type) {
case *url.Error: err.Err.(net.Error).Timeout()
case net.Error: err.Timeout()
case *net.OpError: err.Timeout()
} // + "use of closed network connection" string yoxlaması

// Resumable download:
req.Header.Set("Range", "bytes="+start+"-")     // fayl ölçüsü = başlanğıc
res.Header.Get("Accept-Ranges") != "bytes"      // server dəstəyi yoxla

// Custom JSON error:
type Error struct {
    HTTPCode int    `json:"-"`
    Code     int    `json:"code,omitempty"`
    Message  string `json:"message"`
}
JSONError(w, e)                                  // {"error":{...}} + status

// Arbitrary JSON:
var f interface{}
json.Unmarshal(data, &f)
m := f.(map[string]interface{})                  // object
switch vv := v.(type) { ... }                    // bool/float64/[]i{}/map/nil/string

// API versiya:
http.HandleFunc("/api/v1/test", h)                // URL-də
req.Header.Set("Accept", "application/vnd.mytodos.json; version=2.0")  // content-type
```

## Cloud (T56-61)
```go
// Vendor lock-in qarşısı:
type File interface { Load(string) (io.ReadCloser, error); Save(string, io.ReadSeeker) error }
// + package error-ləri:
var ErrFileNotFound = errors.New("File not found")
// + host detect:
os.Hostname() / net.LookupHost(name)
// + dependency:
exec.LookPath("name")
// + cross-compile:
GOOS=linux GOARCH=amd64 go build
gox -os="linux darwin windows" -arch="amd64 386" -output="dist/{{.OS}}-{{.Arch}}/{{.Dir}}"
// + runtime monitoring:
runtime.NumGoroutine(); runtime.ReadMemStats(m)   // BAHALI — seyrək!
```

## Microservice kommunikasiya (T62-65)
```go
// Keep-alive (custom Transport):
Dial: (&net.Dialer{Timeout: 30s, KeepAlive: 30s}).Dial
// Body-ni VAXTINDA bağla (defer YOX!):
o, _ := ioutil.ReadAll(r.Body); r.Body.Close()   // → reuse mümkün
// codecgen (reflection-sız JSON):
//go:generate codecgen -o user_generated.go user.go
jh := new(codec.JsonHandle); codec.NewEncoderBytes(&out, jh).Encode(&u)
// Protobuf:
u := &pb.User{Name: proto.String("Inigo"), Id: proto.Int32(1234)}
body, _ := proto.Marshal(u)
// gRPC (proto3):
protoc --go_out=plugins=grpc:. hello.proto
pb.RegisterHelloServer(s, &server{})
c := pb.NewHelloClient(conn)
r, err := c.Say(context.Background(), hr)
ctx, cancel := context.WithCancel(context.Background())   // server: <-ctx.Done()
```

## Reflection (T66-70)
```go
// Kind switch:
ref := reflect.ValueOf(val)
switch ref.Kind() { case reflect.Int, reflect.Int64: ref.Int() ... }
// İnterfeys yoxlaması (runtime seçim):
stringer := (*fmt.Stringer)(nil)                  // nil pointer trick
reflect.TypeOf(target).Elem(); t.Implements(iface)
// Struct gəzinti:
val := reflect.Indirect(reflect.ValueOf(u))
t := val.Type()
t.NumField(); t.Field(i); val.Field(i)            // eyni index!
// Tag:
field.Tag.Get("ini")
// Yazma:
v.Field(i).Set(reflect.ValueOf(converted))
// go generate:
//go:generate ./queue MyInt
tt.Execute(file, map[string]string{"MyType": ..., "Package": os.Getenv("GOPACKAGE")})
```
