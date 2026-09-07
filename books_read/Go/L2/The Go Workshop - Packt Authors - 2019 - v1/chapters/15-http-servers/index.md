# Chapter 15 — HTTP Servers (HTTP Serverlər)

## Bu fəsil nədən bəhs edir?

HTTP handler interfeysi (ServeHTTP), ListenAndServe, sadə routing (Handle vs
HandleFunc), Handler vs HandlerFunc seçimi, səhifə sayğacı (state), JSON cavabı,
querystring parametrləri (r.URL.Query), html/template (placeholder, if/else, .Field),
statik fayllar (ServeFile, FileServer, StripPrefix), xarici template (ParseFiles),
HTML form (POST, ParseForm, r.Form.Get), JSON request (NewDecoder) və REST anlayışı.

## Əsas fikirlər

### 1. HTTP Handler — Əsas Blok
```go
type MyHandler struct{}
func (h MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("HI"))
}

log.Fatal(http.ListenAndServe(":8080", MyHandler{}))
```
- Handler interfeysi = TƏK ServeHTTP metodu
- `w` — cavaba yazma; `r` — sorğudan oxuma
- ListenAndServe xətası → log.Fatal (proqram dayansın)

### 2. Sadə Routing — Handle vs HandleFunc
**Problem:** Handler birbaşa verilsə BÜTÜN yolları o tutur — HandleFunc işləmir.

**Kitabdan kod nümunəsi:**
```go
func main() {
    http.HandleFunc("/chapter1", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("<h1>Chapter 1</h1>"))
    })
    http.Handle("/", hello{})                   // DEFAULT handler "/" yolunda
    log.Fatal(http.ListenAndServe(":8080", nil))  // nil = DefaultServeMux
}
```

**Handle vs HandleFunc:**
| Handle | HandleFunc |
|---|---|
| http.Handler (struct) gözləyir | Funksiya gözləyir |
| State/durum saxlamaq üçün (sayğac) | Sadə statik cavablar üçün |

### 3. Səhifə Sayğacı — Handler-in Gücü
```go
type PageWithCounter struct {
    counter  int
    heading  string
    content  string
}

func (h *PageWithCounter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    h.counter++                                    // HƏR istifadəçinin ÖZ sayğacı
    w.Write([]byte(fmt.Sprintf(
        "<h1>%s</h1><p>%s</p><p>Views: %d</p>",
        h.heading, h.content, h.counter)))
}

// 3 MÜSTƏQİL instance — 3 müstəqil sayğac:
hello := &PageWithCounter{heading: "Hello World", ...}
chapter1 := &PageWithCounter{heading: "Chapter 1", ...}
http.Handle("/", hello)
http.Handle("/chapter1", chapter1)
```
Handler-i seçmək = STATE saxlamaq olar; sadə funksiya bunu edə BİLMƏZ.
(Qeyd: paralel istifadə üçün konkurensiya lazımdır — sonrakı fəsil.)

### 4. JSON Cavabı
Handler-in İÇİNDE json.Marshal → w.Write — hər hansı struct JSON kimi xidmət edilə
bilər (microservice əsası).

### 5. Querystring — Dynamic Content
```go
func Hello(w http.ResponseWriter, r *http.Request) {
    vl := r.URL.Query()                     // map[string][]string!
    name, ok := vl["name"]                  // comma-ok açar mövcudluğu
    if !ok {
        w.WriteHeader(400)                   // Bad Request
        w.Write([]byte("Missing name"))
        return
    }
    w.Write([]byte(fmt.Sprintf("Hello %s", strings.Join(name, ","))))
}
```
- `?name=john` — sorğu dəyərləri GET ilə ötürülür
- Dəyər []string — çoxsaylı eyni adlı parametr ola bilər → Join

### 6. Templating — html/template
**Nədir:** boşluqlu mətn iskeleti; engine dəyərlərlə doldurur.

**Kitabdan kod nümunəsi:**
```go
var tplStr = `
<html>
  <h1>Customer {{.ID}}
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
</html>`

type Customer struct {
    ID      int
    Name    string
    Surname string
    Age     int
}

// Query-dən doldur:
cust := Customer{}
if id, ok := vl["id"]; ok {
    cust.ID, _ = strconv.Atoi(strings.Join(id, ","))
}

tmpl, _ := template.New("test").Parse(tplStr)
tmpl.Execute(w, cust)         // doldurub birbaşa ResponseWriter-a yaz
```
**Sintaksis:** `{{.Sahə}}` — struct sahəsi; `{{if}}...{{else}}...{{end}}` — şərt;
hər if-in ÖZ end-i var.

### 7. Statik Fayllar
**Tək fayl:**
```go
http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "./index.html")
})
```
Server işləyərkən faylı dəyiş → refresh → YENİ məzmun (recompile YOX!) — resurs
koddan ayrılır.

**Qovluq + maskalama:**
```go
http.Handle("/statics/",
    http.StripPrefix("/statics/",
        http.FileServer(http.Dir("./public"))))
```
- FileServer(Dir) — qovluqdakı bütün faylları xidmət edir
- StripPrefix — URL yolundan "/statics/" silir → public/ daxili yolu
- /statics/body.css → ./public/body.css

**HTML-də CSS:** `<link rel="stylesheet" href="/statics/body.css">`

### 8. Xarici Template — ParseFiles
```go
func NewHello(tplPath string) (*Hello, error) {
    tmpl, err := template.ParseFiles(tplPath)     // FAYLDAN template
    if err != nil { return nil, err }
    return &Hello{tpl}, nil
}
type Hello struct {
    tpl *template.Template                         // handler template saxlayır
}
```
Startup-da yüklə → performans (fayl əməliyyatı ən yavaşdır); runtime reload üçün
konkurensiya lazımdır.

### 9. HTML Form + POST
**form.html:**
```html
<form method="post" action="/">
  <li>Name: <input type="text" name="name"></li>
  <li><input type="submit" value="send"></li>
</form>
```

**Handler (kitabdan):**
```go
func (h Hello) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    vst := Visitor{}
    if r.Method == http.MethodPost {              // METOD yoxla!
        err := r.ParseForm()                       // formu parse et
        if err != nil {
            w.WriteHeader(400)
            return
        }
        vst.Name = r.Form.Get("name")             // sahə dəyərləri
        vst.Surname = r.Form.Get("surname")
        vst.Age = r.Form.Get("age")
    }
    h.tpl.Execute(w, vst)                          // GET = boş; POST = dolu
}
```

### 10. JSON Request Qəbulu
**Kitabdan kod nümunəsi:**
```go
http.HandleFunc("/", func(wr http.ResponseWriter, req *http.Request) {
    decoder := json.NewDecoder(req.Body)            // stream decoder
    var data Request
    err := decoder.Decode(&data)
    if err != nil {
        wr.WriteHeader(400)
        return
    }
    rsp := Response{Greeting: fmt.Sprintf(
        "Hello %s %s", data.Name, data.Surname)}
    bts, _ := json.Marshal(rsp)
    wr.Write(bts)                                   // JSON → JSON
})
```
Test alətləri: Insomnia/Postman (brauzer JSON body göndərə BİLMƏZ).

### 11. REST Metodları
| Metod | Məqsəd |
|---|---|
| GET | data AL |
| POST | resurs ƏLAVƏ et (bilinməyən yerdə) |
| PUT | məlum yerdə dəyər QOY (update) |
| DELETE | sil |

REST API = yollar + metodlar toplusu. POST/PUT fərqi: POST = axtarış tələb edən
dəyişiklik; PUT = məlum lokasiyaya yazı. Browser yalnız GET/POST edir — qalanları
client kodu (Insomnia, curl, Go client).

## Əsas terminlələr
- HTTP Server — sorğu qəbul edib cavab verən proqram
- Handler İnterfeysi — ServeHTTP(ResponseWriter, *Request)
- Handle/HandleFunc — Handler struct / funksiya qeydiyyatı
- DefaultServeMux — default router (ListenAndServe nil)
- Query Parameter/Querystring — ?name=dəyər
- r.URL.Query() — map[string][]string
- html/template — HTML template engine
- Placeholder {{.Field}} — struct sahə istinadı
- ParseFiles — fayldan template yaratma
- tmpl.Execute(w, data) — doldur + cavaba yaz
- http.ServeFile — tək statik fayl xidməti
- http.FileServer(http.Dir) — qovluq fayl serveri
- http.StripPrefix — URL prefiks maskeleme
- r.ParseForm / r.Form.Get — POST form məlumatı
- json.NewDecoder(req.Body) — stream JSON decode
- REST — yol+metod API modeli

## Praktik nətidə

(1) Handler = state lazımdırsa (sayğac); HandleFunc = sadə cavab. (2) Handler-i
ListenAndServe-a birbaşa veribsən HandleFunc-lar İŞLƏMİR — handler-ı "/" yoluna
qeyd et, ListenAndServe-a nil ver. (3) Query map[string][]string-dir — Get/Join
ilə string al. (4) Parametr yoxdursa 400 + return — istifadəçiyə aydın xəta.
(5) Template: {{.Field}} struct-dan; {{if}}...{{end}} hər şərtin ÖZ end-i; Parse
bir dəfə, Execute hər sorğuda. (6) Statik fayllar koddan AYRI — dəyişiklik restart
TƏLƏB ETMİR; template-ləri belə fayla çıxar. (7) FileServer + StripPrefix —
qovluq adını gizlət və yolu yönləndir. (8) POST form: r.Method yoxla → ParseForm →
r.Form.Get. (9) JSON API: NewDecoder(req.Body).Decode — stream; cavab Marshal →
Write. (10) Template-i startup-da yüklə (performance); handler struct template
saxlasın. (11) Sayğac kimi state paralel sorğularda təhlükəsiz DEYİL — konkurensiya
fəsili gəlir.

## Mənbə
Pages: 515-557 (PDF 548-591)
