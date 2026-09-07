# Go in Practice, Second Edition — Cheat Sheet (Azərbaycanca)

Kozyra, Butcher, Farina — Manning 2023 — 13 fəsil, L4 Advanced

## 1. CLI (Ch2)
```bash
# flag paketi (Plan 9 stili: -la BİR flag!)
flag.String("name", "World", "help")   # pointer
flag.StringVar(&v, "name", "def", "help")
flag.Parse(); flag.Args(); flag.PrintDefaults()

# Cobra (subcommand-lı CLI):
&cobra.Command{Use: "add", Run: f, Args: cobra.ExactArgs(2)}
cmd.Flags().StringP("name", "n", "World", "help")
cmd.Flags().GetString("name")
```
Enum YOXDUR → slice + slices.Contains validasiyası.

## 2. Konfiqurasiya
```go
os.Getenv("MYAPP_PORT")                      # 12-factor env (namespace ilə!)
json.NewDecoder(file).Decode(&conf)           # JSON (struct sahələri EXPORT)
yaml.ReadFile → config.Get/GetBool            # go-gypsy
ini.Load → config.Section("S").Key("k")      # go-ini
```

## 3. Struct / Metod / İnterfeys / Generics (Ch3)
```go
func (a Animal) speak()     # value — dəyişmir
func (ch *character) fixName()   # pointer — dəyişir
type Foo string             # alias — builtin-ə metod əlavəsi
type AnimalType interface { Cat | Dog }    # union
type Stack[T any] struct{ vals []T }        # generic
type Numeric interface{ ~int8 | ~int16 }      # ~ = underlying tiplər
constraints.Integer / Signed / Ordered       # hazır dəstlər
var zero T                                   # generic zero value
```

## 4. Errors (Ch4)
```go
return "", errors.New("No strings supplied")   # zero value + error BİRGƏ
var ErrTimeout = errors.New("...")               # sentinel
errors.Is(err, ErrTimeout)                        # == YOX (wrap olsa!)
fmt.Errorf("ctx: %w", err)                        # wrap
fmt.Errorf("ctx: %v", err)                        # wrap YOX
panic(errors.New("..."))                           # panic-ə ERROR ötür
defer func() { if r := recover(); r != nil { err = r.(error) } }()  # adlı return!
```
Nil interface tələsi: `var e StatusErr` + return → həmişə non-nil! `return nil` və ya
`var err error`.

## 5. Paralellik (Ch5)
```go
var wg sync.WaitGroup
wg.Add(1); go func(f string) { defer wg.Done(); compress(f) }(file); wg.Wait()
# LOOP DƏYİŞƏNİ PARAMETR KİMİ!

type words struct{ sync.Mutex; found map[string]int }
w.Lock(); defer w.Unlock()          # bütün çıxışlar eyni kilid!

# Channel:
ch := make(chan int); ch := make(chan bool, 1)   # buferli = KİLİD
select { case v := <-ch: ; case <-time.After(30*time.Second): return }
# close-u YALNIZ sender; done kanalı ilə "bit" siqnalı
```

## 6. Keyfiyyət (Ch6)
```bash
gofmt -d file.go; goimports -w .
go vet                     # context leak, json tag unexported...
go test ./... --race       # RACE MÜTLƏQ
go test -cover -coverprofile=c.out && go tool cover -html=c.out
go test -bench=. -benchmem; go test --fuzz=Fuzz
go get -u && go mod tidy
```
```go
slog.Info("msg", slog.String("k","v"), slog.Group("g", ...))   # 1.21 struktur log
debug.PrintStack(); runtime.Stack(buf, true)
t.Run(test.name, func(t *testing.T) {...})                     # adlı subtest
f.Fuzz(func(t *testing.T, seed string) {...})                   # invariant
```

## 7. Fayl / Şəbəkə (Ch7)
```go
os.ReadFile("f")                     # hamısı
scan := bufio.NewScanner(file); scan.Split(bufio.ScanLines)  # stream
io.Copy(dst, src)
net.Dial("tcp", "host:1902"); log.New(conn, "example ", f)
logger.Panicln(...)                  # Fatal YOX (defer flush!)
//go:embed files                     # binary daxil
var f embed.FS; http.FileServer(http.FS(f))
```
Websocket: `websocket.Handler(ws)` + Message.Send/Receive; SSE: `w.(http.Flusher)` +
`text/event-stream` + `event: X\ndata: Y\n\n`.

## 8. HTTP Server (Ch8, Go 1.22)
```go
http.HandleFunc("GET /comments/{id}", h)    # metod + path dəyişəni
r.PathValue("id")
server := http.Server{ReadTimeout: 1s, WriteTimeout: 2s, Handler: mux}
http.TimeoutHandler(mux, 2s, "timeout")
# middleware:
func mw(next http.HandlerFunc) http.HandlerFunc { return func(w,r){ ctx := context.WithValue(r.Context(), k, v); next(w, r.WithContext(ctx)) } }
r.URL.Query().Get("q"); r.ParseForm(); r.Form.Get("name"); r.Form.Has("k")
http.SetCookie(w, &http.Cookie{Name:"u", Value:v, Expires:...})
r.BasicAuth()
# JWT: jwt.NewWithClaims(HS256, claims).SignedString(key); jwt.Parse(...)
```

## 9. Template (Ch9)
```go
var t = template.Must(template.ParseFiles(...))   # İNİC-DƏ parse (keş!)
t.ExecuteTemplate(w, "index.html", data)           # ad seçimi
t.Funcs(template.FuncMap{"dateFormat": fn})        # Parse-dən ƏVVƏL
{{.X | dateFormat "Jan 2, 2006"}}                   # pipe
{{template "head.html" .}}                          # nest
{{define "base"}}...{{block "styles" .}}...{{end}}  # inheritance
template.HTML(b.String())                           # safe (yalnız öz renderindən!)
var b bytes.Buffer; t.Execute(&b, p); if err==nil { b.WriteTo(w) }  # buffer
```

## 10. Statik + Upload (Ch10)
```go
http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("./s"))))
f, h, err := r.FormFile("file"); defer f.Close(); io.Copy(out, f)
r.ParseMultipartForm(16 << 20); r.MultipartForm.File["files"]   # çoxlu
http.DetectContentType(buf512)                                   # MIME (məzmun!)
mr, _ := r.MultipartReader(); mr.NextPart(); part.Read(buf)      # STREAM upload
```

## 11. Xarici Servis (Ch11)
```go
cc := &http.Client{Timeout: time.Second}
# timeout: type switch *url.Error/net.Error/*OpError + closed-conn string
req.Header.Set("Range", "bytes="+start+"-")    # davamlı download
# JSON error: {"error":{code,message}} — hər iki tərəfdən eyni Error tipi
# arbitrary JSON: interface{} + map[string]interface{} + rekursiv type switch
# versiya: /api/v1/ (asan) və ya Accept: application/vnd.X.json; version=2.0
# gRPC: protoc --go-grpc_out=. → grpc.Dial + NewChatServiceClient(conn)
```

## 12. Bulud (Ch12)
```go
tr := &http.Transport{Dial: (&net.Dialer{Timeout: 30s, KeepAlive: 30s}).Dial}  # reuse!
r.Body.Close() DƏRHAL (serial request-lərdə defer YOX)
GOOS=linux GOARCH=arm64 go build; gox -os="..." -arch="..."
filepath.Join/Separator; //go:build !windows; foo_windows.go
go runtime monitorinqi: go monitorRuntime() — NumGoroutine/ReadMemStats (SEYRƏK!)
exec.LookPath("dep")     # asılılıq yoxla
```

## 13. Refleksiya/Codegen (Ch13)
```go
reflect.ValueOf(v).Kind()          # kind switch (MyInt→int ailəsi)
reflect.Indirect(v)                # pointer aç (təhlükəsiz)
t.NumField(); t.Field(i); val.Field(i)   # eyni indeks
field.Tag.Get("ini")                # öz tag-lərin
(*fmt.Stringer)(nil)                # interfeys reflect üçün nil pointer
//go:generate ./queue MyInt         # codegen — dev-time, nəticə VCS-də
```

## Universal Qaydalar
1. Error sonuncu; zero value ilə birgə; iqnor `_` ilə açıq
2. panic → recover yalnız defer-də; goroutine başlan hər yerdə recover
3. HTTP client: custom Client + timeout; DefaultClient yox
4. close-u sender edir; done kanalı ilə siqnal
5. Loop dəyişənini goroutine-a parametr ötür
6. User input: həmişə sanitize (html/template, MIME sniff, regex QADAĞAN)
7. Parse/template bir dəfə — init-də; runtime-da YOX
8. reflection → codegen → generics gücü sırası ilə düşün
9. `-race` `--fuzz` `-cover` CI-da mütləq
