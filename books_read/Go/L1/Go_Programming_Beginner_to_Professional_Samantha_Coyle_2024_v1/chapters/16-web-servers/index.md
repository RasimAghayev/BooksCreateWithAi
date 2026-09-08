# Chapter 16 — Web Servers (səh. 504-539)

## Bu fəsil nədən bəhs edir?

HTTP serverlər: Handler interfeysi (ServeHTTP), Handle/HandleFunc
routing, handler vs handler funksiyası, middleware (wrapper pattern),
querystring parametrləri, html/template ({{.Field}}, {{if}}), statik
fayllar (ServeFile, FileServer, StripPrefix), xarici template faylları
və //go:embed ilə binary-yə daxil etmə.

## Əsas fikirlər

### 1. Handler interfeysi + Hello World
```go
// Handler = ServeHTTP metodu olan struct:
type hello struct{}

func (h hello) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    msg := "<h1>Hello World</h1>"
    w.Write([]byte(msg))
}

func main() {
    log.Fatal(http.ListenAndServe(":8080", hello{}))
}
```
- w.Write — cavaba yaz; r — sorğu (parametrlər, metod...)
- log.Fatal — server xətasında proqramı dayandır
- localhost:8080 + İSTƏNİLƏN yol → Hello World

### 2. Routing (HandleFunc / Handle)
```go
func main() {
    // HandleFunc — FUNKSIYA ilə:
    http.HandleFunc("/chapter1", func(w http.ResponseWriter, r *http.Request) {
        msg := "<h1>Chapter 1</h1>"
        w.Write([]byte(msg))
    })

    // Handle — HANDLER ilə:
    http.Handle("/", hello{})

    log.Fatal(http.ListenAndServe(":8080", nil))   // nil = DefaultServeMux
}
```
- hello{} baş handler olsa /chapter1-iƏƏ üstələyir (override!) —
  Handle ilə "/" ayrıca qeyd edilməlidir
- Naməlum yollar "/" fallback-inə düşür

### 3. Handler vs HandlerFunc
| http.Handle | http.HandleFunc |
|---|---|
| Handler (struct + ServeHTTP) qəbul edir | sadə funksiya qəbul edir |
| vəziyyət saxlaya bilir (sayğac və s.) | sadə cavablar üçün ideal |

**Qayda:** statik sadə səhifələr → HandleFunc; sayğac/dəyər saxlama
lazımdır → Handler struct.

**Sayğac activity (Activity 16.01):**
```go
type PageWithCounter struct {
    counter int
    heading string
    content string
}

func (h *PageWithCounter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    h.counter++                    // pointer receiver — sayğacı dəyiş!
    w.Write([]byte(fmt.Sprintf(
        "<h1>%s</h1><p>%s</p><p>Views: %d</p>",
        h.heading, h.content, h.counter)))
}

func main() {
    http.Handle("/", &PageWithCounter{heading: "Hello World", content: "..."})
    http.Handle("/chapter1", &PageWithCounter{heading: "Chapter 1", ...})
    http.Handle("/chapter2", &PageWithCounter{heading: "Chapter 2", ...})
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```
- Hər səhifənin MÜSTƏQİL sayğacı

### 4. Middleware
**Nədir:** ortaq davranışı çıxarıb zəncirə salan wrapper funksiyası.

```go
func Hello(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        msg := "Hello there,"
        w.Write([]byte(msg))
        next.ServeHTTP(w, r)        // növbəti funksiyanı çağır
    }
}

func Function1(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte(" this is function 1"))
}
func Function2(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte(" and now we are in function 2"))
}

func main() {
    http.HandleFunc("/hello1", Hello(Function1))   // WRAP!
    http.HandleFunc("/hello2", Hello(Function2))
    log.Fatal(http.ListenAndServe(":8085", nil))
}
// /hello1 → "Hello there, this is function 1"

// Zəncir:
Hello(Middleware2(Middleware3(Function2)))
```
- Middleware = "aradakı adam": sorğunu tutur → iş görür → next-i çağırır

### 5. Querystring (GET parametrləri)
```go
func Hello(w http.ResponseWriter, r *http.Request) {
    vl := r.URL.Query()              // map[string][]string!

    name, ok := vl["name"]
    if !ok {
        w.WriteHeader(400)            // Bad Request
        w.Write([]byte("Missing name"))
        return
    }
    w.Write([]byte(fmt.Sprintf("Hello %s", strings.Join(name, ","))))
}

http.HandleFunc("/", Hello)
// localhost:8080?name=john → Hello john
// localhost:8080 → 400 Missing name
```
- Query() → map[string][]string — hər açar SLİCE-dir (çox dəyər ola
  bilər) — Join ilə birləşdir

### 6. Templating (html/template)
**Nədir:** boşluqlu mətn iskeleti; motor dəyərlərlə doldurur.

```go
var tplStr = `
<html>
  <h1>Customer {{.ID}}</h1>
  {{if .ID }}
  <p>Details:</p>
  <ul>
  {{if .Name}}<li>Name: {{.Name}}</li>{{end}}
  {{if .Surname}}<li>Surname: {{.Surname}}</li>{{end}}
  {{if .Age}}<li>Age: {{.Age}}</li>{{end}}
  </ul>
  {{else}}
  <p>Data not available</p>
  {{end}}
</html>
`

type Customer struct {
    ID      int
    Name    string
    Surname string
    Age     int
}

func Hello(w http.ResponseWriter, r *http.Request) {
    vl := r.URL.Query()
    cust := Customer{}

    if id, ok := vl["id"]; ok {
        cust.ID, _ = strconv.Atoi(strings.Join(id, ","))
    }
    if name, ok := vl["name"]; ok {
        cust.Name = strings.Join(name, ",")
    }
    // ... surname, age eynilə

    tmpl, _ := template.New("test").Parse(tplStr)
    tmpl.Execute(w, cust)             // ResponseWriter-a birbaşa!
}
```
- `{{.Sahə}}` — struct sahəsi; `{{if}}...{{else}}...{{end}}` — şərtlər
- Performans qeydi: template hər sorğuda parse EDİLMƏMƏLİ — main-də
  bir dəfə yaradıb handler-ə ötür

### 7. Statik fayllar
```go
// Tək fayl:
http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "./index.html")
})
// Faylı dəyiş → serveri YENİDƏN BAŞLATMADAN refresh kifayətdir!

// Qovluq — FileServer:
http.FileServer(http.Dir("./public"))
// /public/myfile.css kimi avtomatik xidmət edir

// Prefix gizlət (yolu dəyiş):
http.Handle(
    "/statics/",
    http.StripPrefix(
        "/statics/",                    // URL-dən bu hissəni KƏS
        http.FileServer(http.Dir("./public")),
    ),
)
// /statics/body.css → ./public/body.css
```
- HTML-də istinad: `<link rel="stylesheet" href="/statics/body.css">`

### 8. Xarici template faylı
```go
// ParseFiles — faylı yaddaşa yüklə:
t, _ := template.ParseFiles("mytemplate.html")

// Performans üçün STARTUP-da yüklə (fayl əməliyyatları yavaşdır);
// dəyişiklikdən sonra restart tələb olunur
```

### 9. Embedding (//go:embed)
```go
package main

import (
    _ "embed"                        // yalnız side effect
    "html/template"
    "log"
    "net/http"
)

//go:embed mytemplate.html            // faylı binary-yə DAXİL ET
var s string                          // məzmun bu dəyişəndə!

func main() {
    t, _ := template.New("mytemplate").Parse(s)   // STRING-dən parse
    http.HandleFunc("/hello1", func(w http.ResponseWriter, r *http.Request) {
        data := struct{ text string }{text: "Hello there"}
        t.Execute(w, data)
    })
    log.Fatal(http.ListenAndServe(":8085", nil))
}
```
- Binary başqa qovluğa köçürüləndə template YENƏ DƏ işləyir
- (Qeyd: bu nümunədə struct sahəsi lowercase `text` — template
  yalnız EXPORTED sahələri görür; `{{.Text}}` + `Text string` olmalıdır)

## Activity icmalları
- **16.01:** PageWithCounter — 3 səhifə + müstəqil view sayğacları
- **16.02:** Xarici template faylından welcome server — name parametri;
  boş olsa "visitor" ({{if .Name}}{{.Name}}{{else}}visitor{{end}})

## Əsas terminlər
- http.ListenAndServe(addr, handler) — server başlatma
- Handler interfeysi — ServeHTTP(ResponseWriter, *Request)
- http.Handle / http.HandleFunc — routing
- DefaultServeMux — nil keçiləndə standart router
- Middleware — func(HandlerFunc) HandlerFunc wrapper zənciri
- Querystring — ?name=john GET parametrləri
- r.URL.Query() → map[string][]string
- html/template — {{.Field}}, {{if}}/{{else}}/{{end}}
- template.Execute(w, data) — birbaşa ResponseWriter-a
- http.ServeFile — tək statik fayl
- http.FileServer(Dir) — qovluq serveri
- http.StripPrefix — URL prefiks kəsici
- template.ParseFiles — fayldan template
- //go:embed — binary-yə fayl daxil etmə

## Praktik nətiəcə
Server əsası: Handle/HandleFunc + ListenAndServe(nil). Sadə cavablar —
HandleFunc; vəziyyət (sayğac) — Handler struct (pointer receiver!).
Ortaq davranış — middleware wrapper-ləri zəncirlə. Dinamik məzmun —
query parametrlər (Query + ok yoxlaması + Join) + html/template;
template-i main-də parse et, hər sorğuda YOX. Statik resurslar —
ServeFile (tək) / FileServer+StripPrefix (qovluq; fayl dəyişəndə
restart lazım deyil). Deployment üçün template-ləri //go:embed ilə
binary-yə qat — tək fayl paylanması.

## Mənbə
Pages: 504-539 (PDF 504-539)
