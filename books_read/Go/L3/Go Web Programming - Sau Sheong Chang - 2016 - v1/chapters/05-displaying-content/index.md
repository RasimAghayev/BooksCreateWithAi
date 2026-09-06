# Chapter 5 — Displaying content

## Bu chapter nədən bəhs edir?

Go template engine-ləri (`text/template` + `html/template`): parse/execute dövrü, 4 əsas action (conditional, iterator, set, include), arqument/ləyiş (variable)/pipeline, custom funksiyalar (FuncMap), kontekst-fəqliliyi (context awareness) və XSS müdafiəsi, nested template-lər və layout-lar, block action (default template).

## Əsas fikirlər

### 1. Template engine spektri
**İki ideal uç:**
1. **Logic-less:** yalnız placeholder əvəzetməsi, məntiq handler-də. Təmiz presentation/logic ayrılığı. Nümunə: Mustache (amma "logic-less" iddiasına baxmayaraq şərti/loop tag-ləri var).
2. **Embedded logic:** template-də proqram kodu — çox güclü, amma məntiq handler-lər arasında səpələnir, saxlanması çətin. Ekstrem nümunə: PHP (özü template engine kimi başlayıb, indi ona görə Smarty/Blade template engine-ləri qurulur!).

**Go-nun mövqeyi:** hibrid — logic-less kimi işlədilə bilər, amma kifayət qədər embedded funksionallıq var.

### 2. Template engine işə salma — 2 addım
1. **Parse:** template mənbəyini (fayl/string) parse edib Template struct yarat
2. **Execute:** parse olunmuş template-i ResponseWriter + data ilə icra et

**Kitabdan kod nümunəsi:**
```go
package main

import (
    "net/http"
    "html/template"
)

func process(w http.ResponseWriter, r *http.Request) {
    t, _ := template.ParseFiles("tmpl.html")
    t.Execute(w, "Hello World!")
}

func main() {
    server := http.Server{
        Addr: "127.0.0.1:8080",
    }
    http.HandleFunc("/process", process)
    server.ListenAndServe()
}
```

**Template faylı (tmpl.html):**
```html
<!DOCTYPE html>
<html>
  <head>
    <title>Go Web Programming</title>
  </head>
  <body>
    {{ . }}
  </body>
</html>
```

**Sub-kod izahı:**
- `{{ . }}` → **nöqtə = ən vacib action**: handler-dən ötürülən data-nı təmsil edir
- `template.ParseFiles("tmpl.html")` → faylı parse edir; qısayoldur, əslində:
```go
t := template.New("tmpl.html")
t, _ := t.ParseFiles("tmpl.html")
```
- `t.Execute(w, "Hello World!")` → "Hello World!" data kimi nöqtənin yerinə düşür

### 3. Parse variasiyaları
- `ParseFiles("f1", "f2", ...)` → variadic; **bir Template qaytarır**: ilk fayl = əsas template (adı fayl adı), qalanları template **set**-i kimi saxlanılır
- `ParseGlob("*.html")` → pattern (glob) ilə parse
- `t.Parse(string)` → string-dən; bütün digər yollar sonunda bunu çağırır
- `template.Must(...)` → parse xətasında panic — error yoxlama boilerplate-i azaldır

### 4. Execute — ExecuteTemplate
```go
t, _ := template.ParseFiles("t1.html", "t2.html")
t.Execute(w, "Hello World!")                        // t1.html icra olunur (birinci)
t.ExecuteTemplate(w, "t2.html", "Hello World!")     // t2.html icra olunur (ad ilə)
```
Template setində **Execute hər zaman birincini** götürür — başqasını seçmək üçün `ExecuteTemplate(w, ad, data)`.

### 5. 4 əsas action

**a) Conditional:**
```go
{{ if arg }}
  some content
{{ else }}
  other content
{{ end }}
```
Handler-dən `rand.Intn(10) > 5` bool ötürüləndə `{{ if . }}` ilə şərti göstərmə.

**b) Iterator (range):**
```go
{{ range array }}
  Dot is set to the element {{ . }}
{{ end }}
```
- array/slice/map/channel üzrə iterasiya; loop daxilində **nöqtə = cari element**
- Fallback variantı: `{{ else }}` → boş kolleksiyada göstərilir ("Nothing to show")

**c) Set (with):**
```go
{{ with arg }}
  Dot is set to arg
{{ end }}
```
- `{{ with "world" }}` → bu blok daxilində nöqtə = "world"; blokdan sonra köhnə dəyərə qayıdır
- `{{ else }}` → arg boşsa fallback

**d) Include:**
```go
{{ template "name" }}       // daxilə qəbul et, data ÖTÜRÜLMÜR
{{ template "name" arg }}   // daxilə qəbul et + data ötür
```
- Hər iki fayl parse edilməlidir, yoxsa nəticə boş olar!

### 6. Arguments, variables, pipelines
**Argument:** dəyər (bool, string, struct sahəsi, açar, metod — 1 dəyər yaxud dəyər+error qaytarmalı, funksiya, nöqtə).

**Variable (`$`):**
```go
{{ range $key, $value := . }}
  The key is {{ $key }} and the value is {{ $value }}
{{ end }}
```
Map iterasiyasında açar/dəyər cütlərini tutmaq üçün.

**Pipeline (Unix bənzəri):**
```go
{{ 12.3456 | printf "%.2f" }}   → 12.35
```
`|` → solun çıxışı sağın girişinə ötürülür; ardıcıllıq mümkündür: `{{ p1 | p2 | p3 }}`.

### 7. Custom funksiyalar (FuncMap)
**Məhdudiyyət:** template funksiyası istənilən sayda giriş qəbul edir; **yalnız 1 dəyər** (yaxud 2 — 2-ci error) qaytara bilər.

**Kitabdan kod nümunəsi:**
```go
func formatDate(t time.Time) string {
    layout := "2006-01-02"
    return t.Format(layout)
}

func process(w http.ResponseWriter, r *http.Request) {
    funcMap := template.FuncMap{ "fdate": formatDate }
    t := template.New("tmpl.html").Funcs(funcMap)
    t, _ = t.ParseFiles("tmpl.html")
    t.Execute(w, time.Now())
}
```

Template-də istifadə (2 bərabər yol):
```html
<div>The date/time is {{ . | fdate }}</div>
<div>The date/time is {{ fdate . }}</div>
```

**Sub-kod izahı:**
- `template.FuncMap{"fdate": formatDate}` → ad → funksiya map-i
- `.Funcs(funcMap)` → FuncMap template-ə bağlanır
- **GOTCHA 1:** FuncMap parse-dən ƏVVEL attach olunmalıdır (parse zamanı funksiyalar artıq tanınmalıdır)
- **GOTCHA 2:** `template.New("tmpl.html")` adı ilə `ParseFiles`-in fayl adından törən ad EYNİ olmalıdır, yoxsa xəta
- Pipeline versiyası daha güclüdür: funksiyaları zəncirləmək olar

### 8. Context awareness (kontekst-fəqliliyi)
**Nədir:** Eyni data template-də **harada yerləşdiyinə görə** fərqli şəkildə escape olunur — bu html/template-un ən maraqlı xüsusiyyətidir.

**Kitabdan təcrübə:** data = `` I asked: <i>"What's up?"</i> ``

| Yerləşmə | Nəticə |
|---|---|
| `<div>{{ . }}</div>` | `I asked: &lt;i&gt;&#34;What&#39;s up?&#34;&lt;/i&gt;` (HTML escape) |
| `<a href="/{{ . }}">` | `I%20asked:%20%3ci%3e...` (URL path escape) |
| `<a href="/?q={{ . }}">` | `I%20asked%3a%20...` (URL query escape) |
| `<a onclick="f('{{ . }}')">` | `I asked: \x3ci\x3e\x22What\x27s up?\x22\x3c\/i\x3e` (JS escape) |

Hər kontekst üçün düzgün escape — avtomatik defensive programming.

### 9. XSS müdafiəsi
**Persistent XSS:** attacker şərhə `<script>alert('Pwnd!');</script>` yazır → server saxlayır → başqa istifadəçiyə xam göstərilir → kod icra olunur.

**Go-nun müdafiəsi:** `html/template` input-u xam yazsa belə escape edir:
```html
<div>&lt;script&gt;alert(&#39;Pwnd!&#39;);&lt;/script&gt;</div>
```
`text/template` YOX — o, context-aware deyil, HTML üçün heç vaxt istifadə etmə!

**Escape-i söndürmək (öz riskinlə):**
```go
t.Execute(w, template.HTML(r.FormValue("comment")))
```
`template.HTML` typecast → escape olmur. Brauzerlərin öz XSS qoruması (IE/Chrome/Safari) da var; `X-XSS-Protection: 0` header-i onu da söndürür:
```go
w.Header().Set("X-XSS-Protection", "0")
```

### 10. Nesting + define — layout-lar
**Problem:** include action adı string konstantdır — `{{ template "content.html" }}` yazsan hər səhifə öz layout faylını tələb edir (layout məqsədi puç olur).

**Həll — define:** template faylı içində template-lər AÇIQ adlandırılır:
```html
<!-- layout.html -->
{{ define "layout" }}
<html>
  <head>
    <title>Go Web Programming</title>
  </head>
  <body>
    {{ template "content" }}
  </body>
</html>
{{ end }}
```

**Eyni faylda çox template:**
```html
{{ define "layout" }}
<html>... {{ template "content" }} ...</html>
{{ end }}
{{ define "content" }}
Hello World!
{{ end }}
```

**İcra:**
```go
t, _ := template.ParseFiles("layout.html")
t.ExecuteTemplate(w, "layout", "")
```

**Content dəyişmə texnikası:** eyni adlı template fərqli fayllarda:
```html
<!-- red_hello.html -->
{{ define "content" }}
<h1 style="color: red;">Hello World!</h1>
{{ end }}

<!-- blue_hello.html -->
{{ define "content" }}
<h1 style="color: blue;">Hello World!</h1>
{{ end }}
```
Handler hansı faylı parse edirsə, həmin content işlənir:
```go
if rand.Intn(10) > 5 {
    t, _ = template.ParseFiles("layout.html", "red_hello.html")
} else {
    t, _ = template.ParseFiles("layout.html", "blue_hello.html")
}
t.ExecuteTemplate(w, "layout", "")
```
→ Eyni layout, dəyişən content.

### 11. block action (Go 1.6+) — default template
**Nədir:** define + include birlikdə — template təyin edir və dərhal yerləşdirir; **override edilə bilər**.

```html
{{ define "layout" }}
<html>
  <body>
    {{ block "content" . }}
      <h1 style="color: blue;">Hello World!</h1>
    {{ end }}
  </body>
</html>
{{ end }}
```
- `{{ block "content" . }}` → "content" template-i burada define olunur; parse zamanı başqa fayldan (red_hello.html) gələn eyni adlı template onu **əvəz edir**
- Heç nə gəlmirsə → blokun öz mətni default kimi işlənir → random crash-i aradan qaldırır

## Action xülasəsi
| Action | Sintaksis | Təyinat |
|---|---|---|
| Nöqtə | `{{ . }}` | Data-nın özü |
| Conditional | `{{ if }} {{ else }} {{ end }}` | Şərt |
| Iterator | `{{ range }} {{ else }} {{ end }}` | Loop |
| Set | `{{ with }} {{ else }} {{ end }}` | Nöqtəni dəyiş |
| Include | `{{ template "ad" arg }}` | Template daxil et |
| Define | `{{ define "ad" }} {{ end }}` | Template adlandır |
| Block | `{{ block "ad" . }} {{ end }}` | Define + include + default |

## Əsas terminlər
- Template Engine (şablon mühərriki)
- Logic-less / Embedded Logic (məntiqsiz / daxili məntiqli)
- text/template vs html/template
- Action (`{{ }}`)
- Dot (nöqtə — data)
- Parse / Execute / ExecuteTemplate
- Template Set (şablon dəsti)
- ParseGlob (pattern parse)
- Must (panic-sarmalayan helper)
- Conditional / Iterator / Set / Include / Define / Block actions
- Argument / Variable (`$x`) / Pipeline (`|`)
- FuncMap (funksiya xəritəsi)
- Context Awareness (kontekst-fəqliliyi)
- XSS (Cross-Site Scripting)
- Persistent XSS Vulnerability
- Escaping (HTML/URL/JS escape)
- template.HTML (unescape typecast)
- Layout (səhifə quruluşu)
- Nested Templates (iç-içə şablonlar)

## Praktik nəticə
- HTML üçün **həmişə html/template** — context-aware escape XSS-in qarşısını avtomatik alır; text/template yalnız plain text üçün.
- Layout üçün `define` + `ExecuteTemplate(w, "layout", data)` — fayl adları deyil, template adları ilə işlə.
- Funksiyaları FuncMap ilə əlavə et — amma parse-dən ƏVVƏL; `New("ad").Funcs(fm)` adı fayl adıyla düz girməlidir.
- `{{ range }}` + `{{ else }}` fallback — boş list üçün "Nothing to show" göstər.
- block action ilə default content — override olunmayan halda crash yoxdur.
- User HTML-inin həqiqətən render olunmasını istəyirsənsə `template.HTML` — amma bunun riskini bilərsən.

## Mənbə
Pages: 117-145 (PDF), book pages 96-124
