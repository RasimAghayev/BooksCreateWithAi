# Go Web Programming — Cheat Sheet (bütün kitabdan toplanmış)

Sau Sheong Chang — Manning, 2016

## Server əsasları

### `http.ListenAndServe("", nil)`
**Nə edir:** Ən sadə web server. Boş addr = bütün interfeyslər, port 80; handler nil = DefaultServeMux.

### Server struct (konfiqurasiya)
```go
server := http.Server{
    Addr:         "127.0.0.1:8080",
    Handler:      mux,             // nil = DefaultServeMux
    ReadTimeout:  ..., WriteTimeout: ...,
}
server.ListenAndServe()               // HTTP
server.ListenAndServeTLS("cert.pem", "key.pem")   // HTTPS
```

### HTTPS sertifikat yaratmaq (dev üçün)
```go
template := x509.Certificate{SerialNumber: ..., Subject: ..., NotAfter: time.Now().Add(365*24*time.Hour), ...}
pk, _ := rsa.GenerateKey(rand.Reader, 2048)
derBytes, _ := x509.CreateCertificate(rand.Reader, &template, &template, &pk.PublicKey, pk)
pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
```

## Handler-lər

### Handler (interfeys)
```go
func (h *MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello World!")
}
http.Handle("/hello", &handler)       // DefaultServeMux-a qeydiyyat
```

### Handler funksiyası + adapter
```go
func hello(w http.ResponseWriter, r *http.Request) { ... }
http.HandleFunc("/hello", hello)      // = mux.Handle(pattern, HandlerFunc(handler))
```

### Middleware (zəncirləmə)
```go
func log(h http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        fmt.Printf("Handler called - %T\n", h)
        h.ServeHTTP(w, r)
    })
}
http.Handle("/hello", protect(log(hello)))     // pipeline
```

### ServeMux qaydaları
- `/hello` (sonu slashsız) → yalnız DƏQİQ match
- `/hello/` (sonu slashlı) → PREFİKS match (`/hello/there` düşür)
- `/` → fallback (heç nə tutulmayanda)
- `http.NewServeMux()` → öz mux (qlobal ad məkanından təcrid)

### HttpRouter (named params)
```go
mux := httprouter.New()
mux.GET("/hello/:name", hello)
func hello(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    fmt.Fprintf(w, "hello, %s!\n", p.ByName("name"))
}
```

## Request emalı

### Form data cədvəli
| Sahə | Metod | URL | urlencoded form | multipart |
|---|---|---|---|---|
| `r.Form` | ParseForm | ✅ | ✅ | ✅ |
| `r.PostForm` | ParseForm | ❌ | ✅ | ❌ |
| `r.MultipartForm` | ParseMultipartForm | ❌ | ✅ | ✅ |
| `r.FormValue(k)` | avto | ✅ | ✅ | ❌ (GOTCHA!) |
| `r.PostFormValue(k)` | avto | ❌ | ✅ | ❌ |

### Fayl yükləmə
```go
file, _, err := r.FormFile("uploaded")     // ən qısa yol
data, _ := ioutil.ReadAll(file)
// Uzun yol: ParseMultipartForm → MultipartForm.File["uploaded"][0].Open()
```

### ResponseWriter
```go
w.Write([]byte(str))                // body; Content-Type yoxdursa 512 byte sniff
w.WriteHeader(501)                  // status; sonra header dəyişməz!
w.Header().Set("Location", url)     // WriteHeader-dən ƏVVƏL
w.Header().Set("Content-Type", "application/json")
```

### Cookie-lər
```go
c := http.Cookie{Name: "flash", Value: base64.URLEncoding.EncodeToString(msg), HttpOnly: true}
http.SetCookie(w, &c)               // Set-Cookie header
c1, err := r.Cookie("first_cookie") // adlı oxu (ErrNoCookie)
cs := r.Cookies()                    // hamısı
// Silmə: eyni ad + MaxAge: -1 + Expires: time.Unix(1, 0)
```

## Template-lər

### Parse/Execute
```go
t := template.ParseFiles("tmpl.html")          // fayl(lar)
t := template.ParseGlob("*.html")                // pattern
t := template.Must(template.ParseFiles(...))    // panic-on-error
t.Execute(w, data)                               // birinci template
t.ExecuteTemplate(w, "layout", data)            // ad verilmiş template
```

### Action-lər
```html
{{ . }}                          <!-- data -->
{{ if . }} ... {{ else }} ... {{ end }}
{{ range . }} {{ . }} {{ else }} Boş {{ end }}
{{ with "world" }} {{ . }} {{ end }}
{{ template "name" . }}          <!-- include + data ötür -->
{{ define "layout" }} ... {{ end }}
{{ block "content" . }} default {{ end }}   <!-- Go 1.6+ -->
```

### Variable/Pipeline/FuncMap
```html
{{ range $key, $value := . }} ... {{ end }}
{{ 12.3456 | printf "%.2f" }}
{{ fdate . }}
```
```go
funcMap := template.FuncMap{"fdate": formatDate}
t := template.New("tmpl.html").Funcs(funcMap)   // parse-dən ƏVVƏL!
```

### Context awareness (html/template)
`{{ . }}` yerləşməsinə görə avtomatik escape: div → HTML escape; href="/..." → URL path escape; onclick → JS escape → XSS müdafiəsi. Unescape: `template.HTML(r.FormValue(...))` (riskli!).

## Data saxlama

### Yaddaş (cache)
```go
var PostById map[int]*Post          // 2 indeks = 2 map
var PostsByAuthor map[string][]*Post
func store(post Post) { PostById[post.Id] = &post; ... }
```

### CSV
```go
writer := csv.NewWriter(csvFile)
writer.Write([]string{strconv.Itoa(id), content, author})
writer.Flush()                        // MÜTLƏQ!
reader := csv.NewReader(file); reader.FieldsPerRecord = -1
records, _ := reader.ReadAll()
```

### gob (Go binary)
```go
buffer := new(bytes.Buffer)
gob.NewEncoder(buffer).Encode(data)
ioutil.WriteFile(filename, buffer.Bytes(), 0600)
// load: ReadFile → bytes.NewBuffer → gob.NewDecoder(buffer).Decode(&v)
```

### database/sql
```go
var Db *sql.DB
Db, _ = sql.Open("postgres", "user=gwp dbname=gwp password=gwp sslmode=disable")  // lazy pool
_ "github.com/lib/pq"              // blank import: driver init qeydiyyatı

// CRUD:
stmt, _ := Db.Prepare("insert ... values ($1,$2) returning id")    // prepared
stmt.QueryRow(a, b).Scan(&post.Id)                                  // tək sətir
Db.QueryRow("select ... where id=$1", id).Scan(&a, &b, &c)          // tək oxu
rows, _ := Db.Query("select ... limit $1", n)                       // çoxlu
for rows.Next() { rows.Scan(&...) }                                  // iterator
Db.Exec("update/delete ...", args)                                   // nəticəsiz
```

### Gorm (ORM)
```go
Db.AutoMigrate(&Post{}, &Comment{})
Db.Create(&post)
Db.Model(&post).Association("Comments").Append(comment)
Db.Where("author = $1", "x").First(&readPost)
Db.Model(&readPost).Related(&comments)
```

## XML / JSON

### XML struct tag qaydaları
```go
type Post struct {
    XMLName xml.Name `xml:"post"`          // element adı
    Id      string   `xml:"id,attr"`        // atribut
    Content string   `xml:"content"`        // child element
    Author  Author   `xml:"author"`          // nested struct
    Xml     string   `xml:",innerxml"`       // raw daxili XML
    Comments []Comment `xml:"comments>comment"`  // leap-frog!
}
type Author struct {
    Id   string `xml:"id,attr"`
    Name string `xml:",chardata"`           // character data
}
```

### XML əməliyyatları
```go
xml.Unmarshal(data, &post)                    // string → struct
decoder := xml.NewDecoder(r)                  // streaming
for { t, err := decoder.Token(); if err == io.EOF { break }
     switch se := t.(type) { case xml.StartElement: decoder.DecodeElement(&c, &se) } }

xml.Marshal(&post)                            // struct → XML (declaration YOXDUR!)
xml.MarshalIndent(&post, "", "\t")            // girintili
xml.Header                                    // <?xml version="1.0"...> — manual əlavə et
encoder := xml.NewEncoder(file); encoder.Indent("", "\t"); encoder.Encode(&post)
```

### JSON
```go
json.Unmarshal(data, &post)                   // tag: `json:"key"`
decoder := json.NewDecoder(r)                 // stream (məs. http Request.Body)
for { err := decoder.Decode(&post); if err == io.EOF { break } }
json.Marshal(&post)                            // sıxılmış
json.MarshalIndent(&post, "", "\t")            // pretty
json.NewEncoder(w).Encode(&post)              // stream write
```
**Seçim qaydası:** stream var → Decoder/Encoder; string var → Unmarshal/Marshal.

## Web service (REST)

### Method dispatch
```go
func handleRequest(w http.ResponseWriter, r *http.Request) {
    var err error
    switch r.Method {
    case "GET":    err = handleGet(w, r)
    case "POST":   err = handlePost(w, r)
    case "PUT":    err = handlePut(w, r)
    case "DELETE": err = handleDelete(w, r)
    }
    if err != nil { http.Error(w, err.Error(), http.StatusInternalServerError) }
}
```

### Handler daxilində
```go
id, _ := strconv.Atoi(path.Base(r.URL.Path))     // /post/1 → 1
output, _ := json.MarshalIndent(&post, "", "\t\t")
w.Header().Set("Content-Type", "application/json")
w.Write(output)
// POST body: r.ContentLength → make([]byte, len) → r.Body.Read → json.Unmarshal
```

### cURL testləri
```bash
curl -i -X POST -H "Content-Type: application/json" -d '{...}' http://127.0.0.1:8080/post/
curl -i -X GET http://127.0.0.1:8080/post/1
curl -i -X PUT -H "Content-Type: application/json" -d '{...}' http://127.0.0.1:8080/post/1
curl -i -X DELETE http://127.0.0.1:8080/post/1
```

## Testing

### testing.T
```go
func TestXxx(t *testing.T) {
    t.Error(err)          // Log + Fail (davam)
    t.Fatal(err)          // Log + FailNow (dayan)
    t.Skip("səbəb")       // keç
    if testing.Short() { t.Skip(...) }   // -short flag
    t.Parallel()          // -parallel N ilə paralel
}
```
```bash
go test -v -cover -short -parallel 3 -run "TestAdı"
```

### Benchmark
```go
func BenchmarkXxx(b *testing.B) {
    for i := 0; i < b.N; i++ { decode("post.json") }
}
```
```bash
go test -run x -bench . -cpu 1     # funksional testləri kənarlaşdır
```

### HTTP test (httptest)
```go
mux := http.NewServeMux()
mux.HandleFunc("/post/", handleRequest(&FakePost{}))
writer := httptest.NewRecorder()
request, _ := http.NewRequest("GET", "/post/1", nil)
mux.ServeHTTP(writer, request)
// writer.Code, writer.Body.Bytes() yoxla
```

### TestMain + setup
```go
func TestMain(m *testing.M) { setUp(); code := m.Run(); tearDown(); os.Exit(code) }
```

### Dependency injection (test double)
```go
type Text interface { fetch(id int) (err error); create() (err error); ... }
type Post struct { Db *sql.DB; Id int; ... }
func handleRequest(t Text) http.HandlerFunc {   // closure qaytar!
    return func(w http.ResponseWriter, r *http.Request) { ... }
}
http.HandleFunc("/post/", handleRequest(&Post{Db: db}))    // real
mux.HandleFunc("/post/", handleRequest(&FakePost{}))       // test double
```

### gocheck / Ginkgo
```go
// gocheck:
Suite(&PostTestSuite{})
func Test(t *testing.T) { TestingT(t) }
func (s *PostTestSuite) SetUpTest(c *C) { ... }      // hər testdən əvvəl
c.Check(x, Equals, 200)                               // fail+davam; Assert = fail+dayan

// Ginkgo (BDD):
var _ = Describe("Get a post", func() {
    BeforeEach(func() { ... })
    Context("scenario", func() {
        It("behavior", func() { Expect(x).To(Equal(200)) })   // Gomega matcher
    })
})
```
CLI: `ginkgo bootstrap`, `ginkgo generate`, `ginkgo convert .`, `ginkgo -v`

## Concurrency

### Goroutine + WaitGroup
```go
go f()                          // funksiyanı goroutine et
var wg sync.WaitGroup; wg.Add(2)
go func1(&wg); go func2(&wg)    // hər birində wg.Done()
wg.Wait()                       // hamısı bitənə qədər
```

### Channel-lər
```go
ch := make(chan int)            // unbuffered (sinxron)
ch := make(chan int, 10)        // buffered (FIFO)
ch := make(chan <- string)      // send-only
ch := make(<-chan string)       // receive-only
ch <- v; v = <-ch               // send/receive
v, ok := <-ch                   // + "açıqdırmı?"
close(ch)                       // bağlı kanal: zero value + ok=false
```

### select
```go
select {
case v, ok := <-ch1: ...
case v, ok := <-ch2: ...
default: ...                    // heç biri hazırl deyilsə (bloklamır)
}
```

### Mutex (race condition)
```go
type DB struct {
    mutex *sync.Mutex
    store map[string][3]float64
}
func (db *DB) nearest(t [3]float64) string {
    db.mutex.Lock()      // axtarış + silmə BİRLİKDƏ kiliddə!
    defer db.mutex.Unlock()
    ...
}
```

### Fan-out / Fan-in
```go
func cut(...) <-chan image.Image {   // receive-only kanal qaytar
    c := make(chan image.Image)
    go func() { ...; c <- result }()
    return c
}
// fan-in: select + WaitGroup + 2 dəyərli receive
```

## Deploy

### Standalone
```bash
go build; ./ws-s                  # foreground
nohup ./ws-s &                    # HUP-a laqeyd
# Upstart /etc/init/ws.conf: respawn; respawn limit 10 5; setuid; exec ...
sudo start ws
```

### Heroku
```go
Addr: ":" + os.Getenv("PORT")    // port env-dən!
```
```bash
godep save                        # Godeps/ + Godeps.json
# Procfile: web: ws-h
heroku create; git push heroku master
```

### GAE (kod dəyişiklikləri)
1. `package main` → başqa paket
2. `main()` → `init()` (Server/ListenAndServe sil)
3. MySQL driveri + `cloudsql:` DSN
4. `$1` → `?`
5. app.yaml
```bash
goapp serve; goapp deploy
```

### Docker
```dockerfile
FROM golang
ADD . /go/src/...
WORKDIR /go/src/...
RUN go get github.com/lib/pq
RUN go install ...
ENTRYPOINT /go/bin/ws-d
EXPOSE 8080
```
```bash
docker build -t ws-d .
docker run --publish 80:8080 --name svc --rm ws-d
docker ps; docker images
# Docker Machine (cloud):
docker-machine create --driver digitalocean --digitalocean-access-token <t> wsd
eval "$(docker-machine env wsd)"
docker build -t ws-d .; docker run --publish 80:8080 ... ws-d
```

## Performans rəqəmləri (kitabdan)
- Goroutine + trivial iş = **78x yavaş** (13.9 vs 1090 ns/op) — başlatma overhead-i
- Goroutine + ağır iş = **7x sürətli**; multi-CPU ilə daha da artır
- json Decode vs Unmarshal: Decode **25% sürətlidir**
- Mosaic: sinxron 2250ms → concurrent (1 CPU) 646ms → multi-CPU 216ms (**10x**)
