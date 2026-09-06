# Chapter 4 — Templating and Rendering Content (Şablonlar və Məzmun Renderi)

## Bu chapter nədən bəhs edir?
Go html/template paketinə: dinamik renderinq, dəyişənlər/loop/şərtlər, nested şablonlar, əsas səhifə layout-u, form emalı, template caching, XSS-dən qorunma (safe rendering) və şablon debug ssenarilərinə.

## Əsas fikirlər

### 1. Dinamik renderinq nədir
**Nədir:** Server tərəfdə render + client tərəfdə interaktivliyin hibridi — botlara (SEO) hazır statik versiya, istifadəçiyə zəngin JS versiyası.

**Tarixçə:** Statik səhifələr → SPA (React/Vue — interaktiv amma yavaş ilk render) → dinamik renderinq (ikisinin də yaxşısı).

**Bookstore üçün:** Dəyişən kataloq (yeni kitablar, qiymətlər) real vaxtda render olunur; axtarış motorları pre-render versiyanı indeksləyir.

### 2. Şablon konstruksiyaları
**Dəyişənlər (variables):** `{{.Title}}`, `{{.Author}}` — Go struct sahələrinin yerini tutur.

**Loop-lar:**
```html
{{range .Books}}
    <h2>{{.Title}}</h2>
    <p>By {{.Author}}</p>
{{end}}
```

**Şərtlər (birləşmiş nümunə — rating filtri):**
```html
{{range .TrendingBooks}}
    {{if .Rating > 4.5}}
        <h2>{{.Title}}</h2>
        <p>By {{.Author}}</p>
    {{end}}
{{end}}
```

### 3. Nested şablonlar (modul quruluşu)
**Base template:**
```html
{{define "base"}}
<html>
<head><title>GitforGits Bookstore</title></head>
<body>
    {{template "header" .}}
    {{template "content" .}}
    {{template "footer" .}}
</body>
</html>
{{end}}
```
**Header / Footer sabit, content dəyişən:**
```html
{{define "header"}}
<div class="header">
    <img src="logo.png" alt="GitforGits Logo">
    <ul class="navigation">
        <li>Home</li><li>Bestsellers</li><li>Genres</li>
    </ul>
</div>
{{end}}

{{define "booklist"}}
<div class="book-list">
    {{range .Books}}
        <div class="book">
            <h2>{{.Title}}</h2>
            <p>By {{.Author}}</p>
        </div>
    {{end}}
</div>
{{end}}
```
**Sub-kod izahı:** `{{template "header" .}}` — adlı şablonu çağırır və CARI data-nı (`.`) ötürür; content hissəsi səhifəyə görə dəyişdirilir (booklist, author-bio, review...).

### 4. Əsas səhifə layout-u
```html
{{define "mainContent"}}
<div class="content">
    <section class="featured">
        <h2>Featured Books</h2>
        {{range .FeaturedBooks}}
            <div class="book">
                <img src="{{.ImageURL}}" alt="{{.Title}}">
                <h3>{{.Title}}</h3>
                <p>By {{.Author}}</p>
            </div>
        {{end}}
    </section>
    <!-- bestsellers / UserRecommended eyni pattern -->
</div>
{{end}}
```
**Prinsip:** Hər bölmə ayrıca template — dəyişiklik bir yerdə, bütün sayta tətbiq olunur.

### 5. Formlar və istifadəçi inputu
**HTML form nümunəsi (qeydiyyat):**
```html
<form action="/register" method="post">
    <label for="username">Username:</label>
    <input type="text" id="username" name="username">
    <input type="email" id="email" name="email" required>
    <input type="password" id="password" name="password">
    <input type="submit" value="Register">
</form>
```
**Seçim elementləri:** `<select>` (janr), `<input type="checkbox">` (müəllif seçimləri), radio.

**Validation:** HTML5 built-in — `required`, `type="email"` pattern yoxlaması.

**Server tərəfi emal (Go):**
```go
username := r.FormValue("username")
email := r.FormValue("email")
password := r.FormValue("password")
```

### 6. Template caching — performans açarı
**Problem:** Hər sorğu üçün `ParseFiles`/`ParseGlob` = təkrar parse xərci + lag.
**Həll:** Başlanğıcda bir dəfə parse → yaddaşda kesh → handler-lər hazır şablonu icra edir.

**Kitabdan kod nümunəsi:**
```go
var templateCache = template.New("").Delims("{{", "}}")

func LoadTemplates() {
    templates, err := templateCache.ParseGlob("./templates/*.gohtml")
    if err != nil { log.Fatal(err) }
    templateCache = templates      // qlobal kesh dolur
}

// Handler artıq hazır keshdən icra edir:
func BookListHandler(w http.ResponseWriter, r *http.Request) {
    books := fetchBooks()
    err := templateCache.ExecuteTemplate(w, "booklist.gohtml", books)
    if err != nil {
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
    }
}
```
**Sub-kod izahı:**
- `LoadTemplates()` — main-də BİR dəfə çağırılır
- `ExecuteTemplate` — parse-sız, yalnız icra → sorğu başına xərc ~0
- **Dev xəbərdarlığı:** kesh ilkin versiyanı saxlayır — dev zamanı dəyişiklikləri görmək üçün keshi refresh etməli (auto-reload aləti: fresh)

### 7. Safe HTML Rendering — XSS qorunması
**XSS (Cross-Site Scripting — saytlararası skriptləmə):** İstifadəçi review-una zərərli script yerləşdirir → başqa istifadəçilər baxanda icra olunur → data oğurluğu, sessiya ələkeçirmə.

**Go həlli — html/template avtomatik contextual encoding:**
- `{{.UserReview}}` — script-i neytrallaşdırır
- Context hissiyyəti: HTML body, attribute, JS, URL, CSS — HƏR context üçün ayrıca escape qaydası

**Etibarlı mənbədən xam HTML (ehtiyatla!):**
```go
trustedHTML := template.HTML(userGeneratedHTML)   // escape-siz render
trustedJS := template.JS(userGeneratedJS)          // JS üçün
// template.URL, template.CSS — digər tiplər
```
**Qayda:** `template.HTML` yalnız mənbəyə 100% əminsən; əks halda XSS qapısı açılır.

**Əlavə müdafiə:** avtomatik escape-dən əlavə manual validation (uzunluq, format, zərərli məzmun filtri).

### 8. Şablon debug ssenariləri
| Ssenari | Simptom | Həll |
|---------|---------|------|
| **Malformed template** `{{if .IsNewRelease}` (qapanış yoxdur) | parse xətası | error mesajındaki delimiter isaresini yoxla |
| **Itkin data** `{{.UserRevu}}` (typo) | runtime xətası | şablon sahə adlarını struct ilə müqayisə et |
| **Yanlış məntiq** `Discount` string, `> 50` int müqayisə | SƏSSİZ səhv — "Regular Price" göstərilir | data tiplərini yoxla; unit test yaz |

**Alətlər:** Go Playground (izolyasiya), custom error wrapper `template.Execute()` ətrafında (kontekst əlavə edir).

## Əsas terminlər
- Dynamic rendering (dinamik renderinq) — server+client hibrid
- Nested template (iç-içe şablon) — `{{define}}`/`{{template}}` quruluşu
- Template caching (şablon keshlənməsi) — parse bir dəfə, icra çox
- ParseGlob/ParseFiles — şablon fayllarını yükləyən funksiyalar
- XSS (Cross-Site Scripting) — skript inyeksiya hücumu
- Contextual encoding (kontekstual kodlaşdırma) — yerə görə escape
- template.HTML/JS/CSS/URL — escape-siz etibar tipləri
- Delims — `{{` `}}` ayırıcıları

## Praktik nəticə
1. Şablonları main-də bir dəfə parse edib qlobal keshdə saxla — sorğu başına parse xərci sıfır.
2. Böyük səhifələri base+header+content+footer nested strukturu ilə qur — DRY.
3. İstifadəçi məzmunu HƏMİŞƏ `{{.Field}}` kimi escape-olunan formada; `template.HTML` yalnız etibarlı mənbə.
4. Şablon xətaları: parse-time (sintaksis) vs runtime (data) — error növünə görə debug et.
5. Dev-də kesh problemi — auto-reload alətindən (fresh) istifadə et.
6. Form validasiyası iki qat: HTML5 (client) + server tərəfli (FormValue-dan sonra öz yoxlaman).

## Mənbə
Pages: 115-139 (PDF səh. 115-139)
