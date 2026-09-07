# Chapter 9 — HTML and email template patterns (HTML və email şablon patternləri)

## Bu chapter nədən bəhs edir?

html/template (context-aware escaping, range/pipelines, custom FuncMap, template keşləmə,
buffered render, nested templates, define/block inheritance, template.HTML), whitespace
trim və text/template + net/smtp ilə email.

## Əsas fikirlər

### 1. Context-Aware Escaping (Təhlükəsizlik Nüvəsi)
**Nədir:** html/template dəyişəni çıxış kontekstinə görə (HTML/JS/URL) avtomatik escape
edir — `<script>alert('xss')</script>` → `&lt;script&gt;...`. Təhlükəsizlik modeli:
şablon yazıcısı ETİBARLI, istifadəçi datası ETİBARSIZ.

`text/template` escape ETMİR — funksiyaları özün əlavə et. Escape-siz çıxış üçün:
`template.HTML` tipi — amma məsuliyyət SƏNİN (user input-u HEÇ VAXT safe sayma).

### 2. Məlumat Əlaqələndirmə və range
```go
type comment struct{ Username, Text string }
type Page struct {
    Title, Content string
    Comments       []comment
}
t.ExecuteTemplate(w, "list.html", p)   // {{.Title}}, {{.Comments}} sahələri
```
Şablon daxilində:
```html
{{ range .Comments }}
  <div class="comment">{{.Text}} <div>by {{.Username}}</div></div>
{{ end }}
```
- range BLOK-SCOPE-ludur: içəridə `{{.}}` = elementin özü; xaricə çıxmaq üçün `{{$.Title}}`
- Şablon faylları: `t.ParseGlob("templates/*.html")` (init-də, log.Fatal ilə)

### 3. Pipelines və Custom Funksiyalar
**Pipe:** `{{"output" | printf "%q"}}` — UNIX CLI kimi; funksiya zənciri `|` ilə.

**FuncMap (öz funksiyaların):**
```go
var funcMap = template.FuncMap{
    "dateFormat": dateFormat,
}
func dateFormat(layout string, d time.Time) string { return d.Format(layout) }

t := template.New("date")
t.Funcs(funcMap)                  // parse-dən ƏVVƏL
t.Parse(tpl)                       // {{.Date | dateFormat "Jan 2, 2006"}}
t.Execute(res, data)
```
Qaydalar: (1) FuncMap Parse-dən əvvəl verilir; (2) pipe bir funksiyadan digərinə SON
arqument kimi ötürür; (3) funksiya adları map açarıdır (Go adı ilə eyni olmaya bilər);
(4) çoxlu şablon üçün helper: `template.Must(t.Parse(tpl))` + Funcs bir yerdə.

### 4. Parse Keşləmə (Performans)
**Qayda:** Şablonu handler-də YOX — package init-də BİR DƏFƏ parse et:
```go
var t = template.Must(template.ParseFiles("templates/simple.html"))
// handler-da yalnız: t.Execute(w, p)
```
Parse ağır işdir (node modeli qurur); hər request-də parse = dublikat iş. `Must` —
səsiz uğursuzluğa yol verməz.

### 5. Buffered Render (Yarım Səhifə Problemi)
**Problem:** Execute error verəndə əvvəlki hissə ARTIQ istifadəçiyə axıdılıb.

**Həll:**
```go
var b bytes.Buffer
err := t.Execute(&b, p)          // əvvəlcə buffera
if err != nil {
    fmt.Fprint(w, "A error occured.")
    return
}
b.WriteTo(w)                     // təmizdirsə onda göndər
```
Trade-off: streaming daha sürətli UX; buffer qəza və test/fuzz üçün təhlükəsiz.
Şablonlar "ağılsız" olmalı — error-lar data mərhələsində həll olunsun.

### 6. Nested Templates
```html
<!-- index.html -->
{{template "head.html" .}}      <!-- daxil et; . = bütün dataset -->
```
```go
t = template.Must(template.ParseFiles("index.html", "head.html"))
t.ExecuteTemplate(w, "index.html", p)   // AD seçimi — Execute İLK faylı işlədər
```
`{{template "head.html" .Foo}}` — subtemplat-a yalnız .Foo dataset-i keçir.

### 7. Template İneritance (define/block)
**Faylı yox, SEÇİRLƏRİ şablon kimi düşün:**
```html
{{define "base"}}<!DOCTYPE HTML>
<html>
  <title>{{template "title" .}}</title>
  {{ block "styles" . }}<style>...</style>{{ end }}   <!-- default + override oluna bilər -->
  <body>
    {{template "content" .}}
    {{block "scripts" .}}{{end}}                       <!-- boş default -->
  </body>
</html>{{end}}
```
Uzadan şablon məcburi hissələri doldurur, opsiyonalları istəsə override:
```html
{{define "title"}}User: {{.Username}}{{end}}
{{define "content"}}<ul>...</ul>{{end}}
{{define "styles"}}<style>h1{color:#800080}</style>{{end}}
```
```go
t := make(map[string]*template.Template)
t["user.html"] = template.Must(template.ParseFiles("base.html", "user.html"))
t["user.html"].ExecuteTemplate(w, "base", u)   // KÖK "base" invoke
```
`block` (Go 1.6+) = define + dərhal invoke; 1.6-dan əvvəl məzmunlu template-i redefine
etmək olmazdı. **Whitespace trim:** `{{ define "head" -}}` / `{{- end }}` — çıxışda
artıq boşluqları kəs.

### 8. Rendered-Obyekt Pattern (template.HTML)
**Nədir:** Obyekti AYRI render et → HTML-i keşlə → yuxarı şablona daxil et (çoxlu
render addımı — cache + separation of concerns).

```go
type Page struct {
    Title   string
    Content template.HTML        // safe marker — escape OLUNMAZ
}
q := &Quote{Quote: "...", Person: "Inigo Montoya"}
var b bytes.Buffer
t.ExecuteTemplate(&b, "quote.html", q)
qc = template.HTML(b.String())    // artıq escape-edilmiş məzmun → safe elan
// sonra: p := &Page{Title: "A User", Content: qc}
```
Escape edilmiş mənbədən gələn HTML safe-dir; **AMMA istifadəçi input-u heç vaxt.**

### 9. Email (text/template + net/smtp)
```go
const emailTemplate = `From: {{.From}}
To: {{.To}}
Subject {{.Subject}}
{{.Body}}
`
type EmailMessage struct {
    From, Subject, Body string
    To                  []string
}
type EmailCredentials struct {
    Username, Password, Server string
    Port                       int
}
t = template.New("email")
t.Parse(emailTemplate)

message := &EmailMessage{From: "me@example.com", To: []string{"you@example.com"}, ...}
var body bytes.Buffer
t.Execute(&body, message)                        // buffera render

auth := smtp.PlainAuth("", creds.Username, creds.Password, creds.Server)
smtp.SendMail(creds.Server+":"+strconv.Itoa(creds.Port), auth,
    message.From, message.To, body.Bytes())       // göndər
```
text/template = escape YOX (html/template onun üzərində qurulub); HTML email üçün
html/template işlət. Şablon + buffer + SMTP — hər istənilən formatlı email.

## Əsas terminlələr
- Context-aware escaping — çıxış kontekstinə görə avtomatik escape
- template.HTML — "safe HTML" marker tipi
- Pipeline (`|`) — UNIX-vari funksiya zənciri
- FuncMap — şablon daxili funksiya qeydiyyatı
- template.Must — səhv halında panic edən parse wrapper
- {{template "name" .}} — subtemplate daxil etmə
- define / block — adlı bölmə / default-məzmunlu dərhal-icra bölməsi
- {{- / -}} — whitespace trim
- net/smtp + PlainAuth — email göndərmə

## Praktik nətidə

Şablon qərarları: (1) user datası → həmişə html/template (XSS avtomatik); (2) format
funksiyaları FuncMap ilə şablona daşı; (3) parse — init + Must (keş); (4) error ehtimalı
varsa buffer-ə render; (5) ortak hissə → {{template}}; səhifə skeleti → define/block
inheritance (map üzrə adlandırılmış template dəstləri); (6) render keşi — template.HTML
(dadlı amma təhlükəli: user HTML heç vaxt); (7) email — eyni mühərrik: şablon → buffer →
smtp.SendMail; (8) şablonlar nəcib-kod deyil — məntiq orada YOX.

## Mənbə
Pages: 220-243 (PDF 241-264)
