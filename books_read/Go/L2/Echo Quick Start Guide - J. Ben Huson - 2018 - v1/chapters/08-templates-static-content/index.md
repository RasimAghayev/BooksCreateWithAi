# Chapter 8 — Providing Templates and Static Content

## Bu chapter nədən bəhs edir?

Bu chapter Echo-nun **statik məzmun** (CSS, JS, fayl) xidməti və **dinamik şablonlar** (server-side rendering) imkanlarını izah edir: `e.Static`/`e.File` metodları, Go standart kitabxanasının `text/template`/`html/template` sintaksisi (pipeline, if/range/with/block), Echo-ya `echo.Renderer` interfeysi ilə şablon inteqrasiyası və şablonlardan Echo funksionallığına (Reverse URL) geri çağırış.

---

## Əsas fikirlər

### 1. Serving static files (Statik faylların xidməti)

**Nədir?** Statik fayllar — brauzerin işlək UI render etməsi üçün lazım olan assetlər: CSS, JavaScript və yüklənməli digər fayllar. Echo bunları route-lar təqdim etdiyiniz kimi təqdim edir.

**Direktoriya səviyyəsində statik xidmət:**

```go
e.Static("/static", "static")
```

**Necə işləyir?**
- 1-ci arqument: route path prefiksi (`/static`) — Echo-nun wildcard matching imkanına əsaslanır (Ch3): `/static/` prefiksi ilə gələn HƏR request statik fayl handler-inə yönəlir.
- 2-ci arqument: filesystem-dəki qovluq (`static`) — binary-nin **current working directory**-sinə bağlıdır.
- Root boş string verilsə, təhlükəsizlik üçün default `.` (CWD) istifadə olunur.

**Tək fayl xidməti:**

```go
e.File("/", "static/index.html")
```

`/` URL-ninə gələn hər sorğu `./static/index.html` məzmunu ilə cavablanır.

**Echo-nun daxili `static` implementasiyası:**

```go
// Static registers a new route with path prefix to serve static files from the
// provided root directory.
func (e *Echo) Static(prefix, root string) *Route {
    if root == "" {
        root = "." // For security we want to restrict to CWD.
    }
    return static(e, prefix, root)
}

func static(i i, prefix, root string) *Route {
    h := func(c Context) error {
        p, err := url.PathUnescape(c.Param("*"))   // wildcard parametrin URL decode
        if err != nil {
            return err
        }
        name := filepath.Join(root, path.Clean("/"+p)) // "/"+ for security
        return c.File(name)                             // faylı serve et
    }
    i.GET(prefix, h)
    if prefix == "/" {
        return i.GET(prefix+"*", h)
    }
    return i.GET(prefix+"/*", h)
}
```

**Vacib detal:** `path.Clean("/"+p)` — path traversal hücumlarına (məs., `../../etc/passwd`) qarşı qoruyur; `url.PathUnescape` isə %-encoded fayl adlarını düzgün dekod edir.

**Nəyə lazımdır?** HTML/JS/CSS məzmununu birbaşa **Go deployment artifact**-i içində paketləyib ayrı statik CDN server qurmadan xidmət etmək — deploy scope-unu kiçildir.

---

### 2. Go template əsasları

`text/template` (yalın mətn) və ya `html/template` (HTML — **code injection qorumalı**, təhlükəsiz) paketi. Nümunə — Reminders cədvəli:

```go
const tmpl = `
<html>
<head> <title>{{.Title}}</title> </head>
<body>
<table>
<th>
    <td>Reminder Name</td> <td>Reminder ID</td> <td>Reminder Due</td>
</th>
{{ range .Reminders }}
<tr>
    <td>{{ .Name }}</td> <td>{{ .ID }}</td> <td>{{ .Due }}</td>
</tr>
{{ else }}
<tr>
    <td colspan=3>No rows!</td>
</tr>
{{ end }}
</body>
</html>
`
```

**Go kodu tərəfi:**

```go
type Reminder struct {
    ID   uuid.UUID
    Name string
    Due  time.Time
}

reminders := []Reminder{
    {ID: uuid.NewV4(), Name: "Oil Change", Due: time.Now().Add(13 * Week)},
    {ID: uuid.NewV4(), Name: "Birthday Party", Due: mustTime(time.Parse("2006-01-02", "2020-01-01"))},
}

// 1. Template-i PARSE et (yalnız bir dəfə!)
t, err := template.New("reminders").Parse(tmpl)
if err != nil {
    log.Fatalf("failed to parse template: %s\n", err.Error())
}

// 2. Data strukturunu hazırla (struct və ya map)
tmplData := struct {
    Reminders []Reminder
    Title     string
}{reminders, "Reminders Page"}

// 3. İORAZ et — io.Writer + data
err = t.Execute(os.Stdout, tmplData)
if err != nil {
    log.Fatalf("failed to render template: %s\n", err.Error())
}
```

**Vacib qaydalar:**
- `Parse` — string şablondan; `ParseFiles` — fayl(lar)dan (fayl ilə işləmək daha rahatdır).
- **Template parsing-i application entry point-də bir dəfə edin** — hər request-də parse etmək israfdır və qaçınılmalıdır.
- `t.Execute(w io.Writer, data interface{})` — data struct və ya `map` ola bilər.

---

### 3. Template sintaksisi — əsas konstruksiyalar

| Konstruksiya | İzah |
|--------------|------|
| `{{/* comment */}}` | Şərh — render olunmur; çoxsətirli ola bilər, iç-içə yasaq |
| `{{pipeline}}` | Dəyişənin default reprizentasiyası ilə çapı (məs., `{{.Title}}`) |
| `{{if pipeline}} T1 {{end}}` | Pipeline boş deyilsə bloku render et |
| `{{if pipeline}} T1 {{else}} T0 {{end}}` | Şərti else |
| `{{if pipeline}} T1 {{else if pipeline}} T0 {{end}}` | Şərti else-if |
| `{{range pipeline}} T1 {{end}}` | Pipeline üzərində iterasiya: array, slice, map, channel |
| `{{range pipeline}} T1 {{else}} T0 {{end}}` | Boş kolleksiya olduqda else bloku render olunur |
| `{{template "name"}}` | Adlandırılmış (nested) template-i icra et |
| `{{template "name" pipeline}}` | Nested template-i verilmiş kontekst ilə icra et |
| `{{block "name" pipeline}} T1 {{end}}` | Subtemplate yaradıb yerində icra et |
| `{{with pipeline}} T1 {{end}}` | Pipeline boş deyilsə bloku render et (if bənzəri) |
| `{{with pipeline}} T1 {{else}} T0 {{end}}` | with + else |

**"Boş" (empty) dəyərlər:** `false`, `0`, `nil`, boş array/slice/map/string.

**Nəyə lazımdır?** Template engine yalnız web render üçün deyil — **email mətni** və istifadəçiyə göndərilən digər mesajların renderində də istifadə olunur.

**Təhlükəsizlik:** HTML çıxışı üçün **həmişə `html/template`** istifadə edin — code injection-a qarşı qoruyur; `text/template` qorumur.

---

### 4. Templates within Echo — echo.Renderer interfeysi

**Addım 1 — Renderer-i custom tip ilə implement et** (`handlers/reminder.go`):

```go
type CustomTemplate struct {
    *template.Template
}

func (ct *CustomTemplate) Render(w io.Writer, name string, data interface{},
    ctx echo.Context) error {
    return ct.ExecuteTemplate(w, name, data)
}
```

**Addım 2 — başlanğıcda parse edib Echo-ya qoş:**

```go
t, err := template.New("reminders").Parse(handlers.RemindersTmpl)
if err != nil {
    panic(err.Error())
}
e.Renderer = &handlers.CustomTemplate{t}
```

**Addım 3 — handler-də `c.Render` ilə işlət:**

```go
func RenderReminders(c echo.Context) error {
    reminders := []Reminder{
        {ID: uuid.NewV4(), Name: "Oil Change", Due: time.Now().Add(30 * 3 * 24 * time.Hour)},
        {ID: uuid.NewV4(), Name: "Birthday Party", Due: mustTime(time.Parse("2006-01-02", "2020-01-01"))},
    }
    tmplData := struct {
        Reminders []Reminder
        Title     string
    }{reminders, "Reminders Page"}
    return c.Render(http.StatusOK, "reminders", tmplData)
}
```

`c.Render(status, "template-name", data)` — `echo.Renderer` interfeysindəki `Render` metodunu çağırır. Yalnız **bir interfeys** implement etməklə istənilən template kitabxanası (Go standard, pongo2 və s.) Echo-ya inteqrasiya olunur.

---

### 5. Calling Echo from templates — Reverse URL

**Problem:** Şablon içində səhifə linkləri düz yazılmamalıdır — rotalar dəyişsə linklər qırılar. Echo-nun `Reverse(name)` funksiyası adlandırılmış route-un URL-ini qaytarır.

**Həll — funksiyanı template data-ya daxil et.** Go şablonları data strukturu üzərindən **metod çağıra bilir**:

```go
type TmplData struct {
    Reminders []Reminder
    Title     string
    rev       func(name string, params ...interface{}) string
}

func (td TmplData) Reverse(name string, params ...interface{}) string {
    return td.rev(name, params...)
}
```

Handler-də `c.Echo().Reverse` funksiyasını daxil et:

```go
data := TmplData{reminders, "Reminders Page", c.Echo().Reverse}
return c.Render(http.StatusOK, "reminders", data)
```

**Route-u adlandır** (`cmd/service/main.go`):

```go
e.POST("/login", handlers.Login).Name = "login"
```

**Şablonda istifadə:**

```html
<a href={{ .Reverse "login" }}>Login</a>
```

Render nəticəsi:

```html
<a href=/login>Login</a>
```

`TmplData.Reverse("login")` → `c.Echo().Reverse("login")` → "login" adlı route-un axtarılıb tam URL-inin qaytarılması.

**Nəyə lazımdır?** Bu pattern ilə TmplData strukturu üzərindən tətbiqetmənin **istənilən aspektini** şablonlara açıq edirsiniz — server-side rendering dinamik və çevik olur.

---

## Əsas terminlər

| Termin | İzah |
|--------|------|
| Static content (Statik məzmun) | CSS, JS, şəkil kimi dəyişmən asset faylları |
| `e.Static(prefix, root)` | Prefiks altında qovluqdan statik fayl xidməti |
| `e.File(path, file)` | Tək faylı konkret route-a bağlama |
| `c.File(name)` | Context metod — fayl məzmununu serve edir |
| `path.Clean` | Path traversal hücumlarına qarşı fayl yolunun təmizlənməsi |
| `text/template` | Yalın mətn üçün şablon paketi (injection qoruması YOX) |
| `html/template` | HTML üçün təhlükəsiz şablon paketi (code injection qoruması VAR) |
| `template.New(...).Parse` / `ParseFiles` | Şablon mətninin/faylının parse edilməsi — entry point-də bir dəfə |
| `t.Execute(w, data)` | Parse edilmiş şablonun data ilə writer-a yazılması |
| Pipeline | Şablon kontekstindəki dəyər/ifadə (`.Title`, `.Reminders`) |
| `{{range}}`/`{{if}}`/`{{with}}`/`{{block}}` | Şablon iterasiya/şərt/blok konstruksiyaları |
| Empty value | `false`, `0`, `nil`, boş string/array/slice/map — if/with üçün "false" |
| `echo.Renderer` | Echo-nun şablon interfeysi — `Render(w, name, data, ctx)` metodu |
| `c.Render(code, name, data)` | Handler-dən qeydiyyatlı renderer ilə render |
| Reverse URL (Əks URL) | Adlandırılmış route-un URL-inin generasiyası — `e.Reverse(name)` |
| `TmplData.Reverse` | Şablon içindən Echo funksiyasını çağırmaq üçün metod-wrapper |
| Server-side rendering (SSR) | Dinamik məzmunun serverdə şablonla render edilməsi |
| Deployment artifact | Deploy edilən icra paketi — statik fayllar daxilində daşıya bilər |

---

## Praktik nəticə

1. **Statik assetlər bir sətirlə:** `e.Static("/static", "static")` + `e.File("/", "static/index.html")` — ayrı CDN xidməti lazım olmadan Go binary-si içindən xidmət; daxili `path.Clean` qoruması path traversal-a qarşı mövcuddur.
2. **Şablon parse-ı bir dəfə, entry point-də:** `template.New(...).Parse(...)` main-də çağrılır, `e.Renderer`-a qoşulur; hər request-də parse etmək böyük israfdır.
3. **HTML üçün `html/template`:** `text/template` yalın mətn (email və s.) üçündür — web səhifələrdə injection riski yaradır.
4. **Custom renderer yalnız 1 interfeys:** `CustomTemplate.Render(w, name, data, ctx)` metodunu implement edib `e.Renderer`-a verin — sonra handler-lərdə `c.Render(200, "name", data)` hazır.
5. **Linkləri hardcode etməyin:** Route-lara `.Name = "login"` verin, şablonlarda `{{ .Reverse "login" }}` istifadə edin — URL dəyişəndə şablonlar avtomatik uyğunlaşır.
6. **TmplData pattern genişlənə biləndir:** İstənilən funksiyanı (auth status, konfiq dəyəri və s.) `rev` kimi struktura daxil edib şablonlara aça bilərsiniz.
7. **Şablon engine-dən yalnız web-də deyil:** Email/mesaj renderində də eyni `Parse`+`Execute` axını işləyir.

---

## Mənbə

- Kitab: *Echo Quick Start Guide* — J. Ben Huson, Packt Publishing, 2018 (ISBN 9781789340849)
- Chapter 8: "Providing Templates and Static Content", book səh. 157–173
- PDF səhifələri: 162–177
- Kod: https://github.com/PacktPublishing/Echo-Essentials/tree/master/chapter8
- Video: https://goo.gl/BAHuF3
