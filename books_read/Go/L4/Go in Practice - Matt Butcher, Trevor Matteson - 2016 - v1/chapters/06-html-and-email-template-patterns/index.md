# Chapter 6 — HTML and email template patterns (Technique 32-38)

## Bu chapter nədən bəhs edir?

html/template paketi: context-aware escaping, custom funksiyalar (FuncMap), parse cache (performans), buffer ilə error handling, nested templates, template inheritance (define/block), data obyekt → HTML mapping (template.HTML), email template-ləri (text/template + smtp).

## Əsas fikirlər

### 1. html/template əsasları
**Nədir:** text/template üzərində qurulmuş, **context-aware** HTML şablon mühərriki.

**Kitabdan kod nümunəsi:**
```go
// templates/simple.html:
<!DOCTYPE HTML>
<html>
  <head>
    <title>{{.Title}}</title>
  </head>
  <body>
    <h1>{{.Title}}</h1>
    <p>{{.Content}}</p>
  </body>
</html>

// Go kodu:
type Page struct {
    Title, Content string
}

func displayPage(w http.ResponseWriter, r *http.Request) {
    p := &Page{
        Title:   "An Example",
        Content: "Have fun stormin' da castle.",
    }
    t := template.Must(template.ParseFiles("templates/simple.html"))
    t.Execute(w, p)
}
```

**Context-aware escaping — avtomatik:**
```html
<!-- Yazılan: -->
<a href="/user?id={{.Id}}">{{.Content}}</a>
<!-- Mühərrrik avtomatik genişləndirir: -->
<a href="/user?id={{.Id | urlquery}}">{{.Content | html}}</a>
```
- **Security model:** template kodu ETİBARLI, user data (dəyişənlər) ETİBARSIZ → escape olunur
- `<script>alert('xss')</script>` input → `&lt;script&gt;...` kimi render → XSS qarşısı alınır

### TECHNIQUE 32: Custom template funksiyaları
**Kitabdan kod nümunəsi (dateFormat):**
```go
var tpl = `<!DOCTYPE HTML>
<html>
  <body>
    <p>{{.Date | dateFormat "Jan 2, 2006"}}</p>   <!-- pipeline! -->
  </body>
</html>`

var funcMap = template.FuncMap{
    "dateFormat": dateFormat,          // ad → funksiya map-i
}

func dateFormat(layout string, d time.Time) string {
    return d.Format(layout)
}

func serveTemplate(res http.ResponseWriter, req *http.Request) {
    t := template.New("date")
    t.Funcs(funcMap)                    // SIRASI VACIB: Funcs → Parse!
    t.Parse(tpl)
    data := struct{ Date time.Time }{Date: time.Now()}
    t.Execute(res, data)
}
```

**Sub-kod izahı:**
- Pipeline `|` — UNIX boru kimi; ötürülən dəyər SONUNCU arqument kimi funksiyaya girir
- `dateFormat "Jan 2, 2006"` → layout birinci, Date (piped) sonuncu parametr
- **FuncMap adları fərqli ola bilər** (template adı ≠ Go funksiya adı)
- Funcs Parse-dan əvvəl çağırılmalı

** Çoxlu template üçün helper:**
```go
func parseTemplateString(name, tpl string) *template.Template {
    t := template.New(name)
    t.Funcs(funcMap)
    t = template.Must(t.Parse(tpl))
    return t
}
```

### TECHNIQUE 33: Parse cache — performans
```go
// PİS: hər request-də parse!
func displayPage(w http.ResponseWriter, r *http.Request) {
    t := template.Must(template.ParseFiles("templates/simple.html"))  // ← İÇƏRIDƏ
    t.Execute(w, p)
}

// YAXŞI: paket init-də BİR DƏFƏ parse:
var t = template.Must(template.ParseFiles("templates/simple.html"))

func displayPage(w http.ResponseWriter, r *http.Request) {
    t.Execute(w, p)         // sadəcə icra
}
```
Ch 5 benchmark dərslə: parse+execute = 10167 ns vs execute-only = 1318 ns → **~10x**. Template pəncərələri: parse ağır işdir (AST node-ları), icra ucuz.

### TECHNIQUE 34: Execution failure — buffer pattern
**Problem:** Execute yarımçıq xəta versə, YARIMÇIQ səhifə istifadəçiyə gedir.

```go
var b bytes.Buffer
err := t.Execute(&b, p)          // ƏVVƏLCƏ buferə yaz
if err != nil {
    fmt.Fprint(w, "A error occured.")   // xəta → təmiz cavab
    return
}
b.WriteTo(w)                      // xəta yoxdursa → istifadəçiyə
```
- Trade-off: streaming (tez hiss) vs buffer (atomiklik). Template-lər "ağılsız" olmalı — xətalar data-da yox, icrada baş verməli deyil.

### TECHNIQUE 35: Nested templates
```html
<!-- index.html: -->
<!DOCTYPE HTML>
<html>
  {{template "head.html" .}}     <!-- başqa template-i DAXİL ET, dataset-i ÖTÜR -->
  <body>
    <h1>{{.Title}}</h1>
  </body>
</html>

<!-- head.html: -->
<head>
  <meta charset="utf-8">
  <title>{{.Title}}</title>       <!-- eyni dataset! -->
</head>
```
```go
t = template.Must(template.ParseFiles("index.html", "head.html"))  // İKİSİ bir obyektə
t.ExecuteTemplate(w, "index.html", p)     // AD ilə icra!
```
- `{{template "ad" .}}` — 3 hissə: directive, ad, dataset
- `.Foo` ötürsən — yalnız Foo-nun sahələri görünər
- **ExecuteTemplate(w, ad, data)** — çox template varsa hansının icra olunacağını seçir (Execute → birinci parse olunanı)

### TECHNIQUE 36: Template inheritance — define + block
**Base (base.html):**
```html
{{define "base"}}<!DOCTYPE HTML>
<html>
  <head>
    <title>{{template "title" .}}</title>              <!-- HANSISA yerdə doldurulacaq -->
    {{ block "styles" . }}<style>                      <!-- DEFAULT + dərhal icra -->
      h1 { color: #400080 }
    </style>{{ end }}
  </head>
  <body>
    <h1>{{template "title" .}}</h1>
    {{template "content" .}}                            <!-- MƏCBURİ — doldurulmalı -->
    {{block "scripts" .}}{{end}}                       <!-- optional, boş default -->
  </body>
</html>{{end}}
```

**Extending user.html:**
```html
{{define "title"}}User: {{.Username}}{{end}}
{{define "content"}}
<ul>
  <li>Username: {{.Username}}</li>
  <li>Name: {{.Name}}</li>
</ul>
{{end}}
```

**Optional override (page.html):**
```html
{{define "title"}}{{.Title}}{{end}}
{{define "content"}}<p>{{.Content}}</p>{{end}}
{{define "styles"}}<style>h1 { color: #800080 }</style>{{end}}   <!-- DEFAULT-U ƏVƏZ EDİR -->
```

```go
t = make(map[string]*template.Template)
temp := template.Must(template.ParseFiles("base.html", "user.html"))
t["user.html"] = temp
temp = template.Must(template.ParseFiles("base.html", "page.html"))
t["page.html"] = temp

// İcra: HƏMİŞƏ base-dən başla!
t["user.html"].ExecuteTemplate(w, "base", u)
t["page.html"].ExecuteTemplate(w, "base", p)
```

**Semantika:**
- `define` — adlı template təyin et (icra YOX)
- `block` (Go 1.6+) — define + dərhal icra + DEFAULT məzmun (override edilə bilər)
- Məcburi bölmələr (default-suz `{{template "content"}}`) extending tərəfindən DOLDURULMALI
- Optional block-lar default qəbul edir, istəyənə override
- **Map strukturu:** hər səhifə öz base+ext cütü ilə parse → `t[ad]` → həmişə "base" icra

### TECHNIQUE 37: Obyekt → HTML mapping (template.HTML)
**Konsepsiya:** hissələri AYRI render et → yüksək səviyyəli template-ə hazır HTML kimi ötür.

**Səbəblər:**
1. **Cache:** bahalı dataset + render → HTML snapshot cache-də saxla → hər səhifədə ötür
2. **Modularity:** hər template BİR şey render etsin

**Kitabdan kod nümunəsi:**
```go
// quote.html — yalnız Quote obyekti:
<blockquote>
&ldquo;{{.Quote}}&rdquo;
&mdash; {{.Person}}
</blockquote>

type Page struct {
    Title   string
    Content template.HTML        // SAFE HTML tipi — escape OLMAZ!
}

type Quote struct {
    Quote, Name string
}

func main() {
    q := &Quote{Quote: `You keep using that word...`, Person: "Inigo Montoya"}
    var b bytes.Buffer
    t.ExecuteTemplate(&b, "quote.html", q)
    qc = template.HTML(b.String())      // bufer string-i → SAFE HTML-ə çevir

    http.HandleFunc("/", displayPage)
    ...
}

func displayPage(w http.ResponseWriter, r *http.Request) {
    p := &Page{
        Title:   "A User",
        Content: qc,                     // hazırlanmış HTML birbaşa daxil olur
    }
    t.ExecuteTemplate(w, "index.html", p)
}
```

**Sub-kod izahı:**
- `template.HTML` — escape-olunmayan tip; çünki quote.html ÖZÜ context-aware escape etdi → nəticə ETİBARLI
- Page struct-un `Content` sahəsi `template.HTML` → index.html icra ediləndə bu sahə escape edilmir
- **XƏBƏRDARLIQ:** user input HEÇ VAXT template.HTML sayılmaz!

### TECHNIQUE 38: Email template
**Kitabdan kod nümunəsi:**
```go
type EmailMessage struct {
    From, Subject, Body string
    To                  []string
}

type EmailCredentials struct {
    Username, Password, Server string
    Port                       int
}

const emailTemplate = `From: {{.From}}
To: {{.To}}
Subject {{.Subject}}
{{.Body}}
`

var t *template.Template

func init() {
    t = template.New("email")
    t.Parse(emailTemplate)
}

func main() {
    message := &EmailMessage{
        From:    "me@example.com",
        To:      []string{"you@example.com"},
        Subject: "A test",
        Body:    "Just saying hi",
    }
    var body bytes.Buffer
    t.Execute(&body, message)          // render → bufer

    authCreds := &EmailCredentials{Username: "myUsername", Password: "myPass",
        Server: "smtp.example.com", Port: 25}
    auth := smtp.PlainAuth("", authCreds.Username, authCreds.Password, authCreds.Server)
    smtp.SendMail(authCreds.Server+":"+strconv.Itoa(authCreds.Port),
        auth, message.From, message.To, body.Bytes())
}
```
- **text/template** istifadə — escape YOX (email mətni üçün düzgün); HTML email üçün html/template
- Render → buffer → `smtp.SendMail(..., body.Bytes())`
- Header formatı (From/To/Subject) template-daxilində

## Template pattern xülasəsi
| Pattern | Həll | İstifadə |
|---|---|---|
| Custom funksiya | FuncMap + Funcs (Parse-dən əvvəl) | dateFormat |
| Parse cache | package-level Must(ParseFiles) | hər server |
| Atomik cavab | buffer → xəta yoxdursa WriteTo | xəta-lı template-lər |
| Nested | `{{template "ad" .}}` + ExecuteTemplate | head/footer paylaşımı |
| Inheritance | define + block + map strukturu | base layout |
| Obyekt mapping | template.HTML + ikili render | cache + modularity |
| Email | text/template + buffer + smtp | bildirişlər |

## Əsas terminlər
- html/template vs text/template (context-aware)
- Context-aware Escaping (urlquery / html)
- XSS müdafiəsi (security model: template = trusted, data = untrusted)
- template.HTML (safe tip)
- Action / Directive (`{{ }}`)
- Pipeline (`|`) — UNIX boru semantikası
- template.FuncMap / t.Funcs (sıra: Funcs → Parse)
- Parse Cache (package-level)
- Buffer Pattern (atomik cavab)
- Streaming vs Buffered cavab
- Nested Templates / ExecuteTemplate
- Template Inheritance
- `define` / `block` (Go 1.6+) / `end`
- Default + Override (block default-u)
- Map of Templates (per-page base+ext)
- İki mərhələli render
- net/smtp / smtp.PlainAuth / smtp.SendMail

## Praktik nəticə
- HTML üçün həmişə html/template — context-aware escaping özbaşına XSS bloklayır; text/template yalnız plain text (email və s.).
- Parse-i handler-dan çıxar → init/package-level: ~10x performans (Ch 5 benchmark sübutu).
- Partial render riski varsa → buffer + error → yalnız tam səhifə göndər.
- Template-ləri "ağılsız" saxla — data xətaları render-dən əvvəl həll olunmalı.
- İrsiyyət üçün: base + block default-ları; hər səhifəni base ilə birgə parse edib map-də saxla; icra həmişə "base"-dən.
- template.HTML yalnız ÖZ şablonun çıxışına — user input-u heç vaxt safe etiketləmə.
- Custom funksiyaları FuncMap ilə qeydiyyatdan keçir; Funcs çağırışı Parse-dən əvvəl — ən çox edilən səhv.

## Mənbə
Pages: 170-189 (PDF), book pages 147-166
