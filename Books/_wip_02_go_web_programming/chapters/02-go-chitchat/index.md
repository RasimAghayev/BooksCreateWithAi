# Chapter 2 — Go ChitChat

## Bu chapter nədən bəhs edir?

Tam funksional internet forum tətbiqinin (ChitChat) qurulması: multiplexer, handler funksiyaları, static fayl servisi, cookie ilə autentifikasiya, html/template ilə HTML generasiyası, PostgreSQL quraşdırılması, database/sql ilə data əlaqəsi. Kitabın arxitektur sxemi burada formalaşır: multiplexer → handler → model → template engine.

**Kod reposu:** https://github.com/sausheong/gwp (ChitChat-in tam kodu)

## Əsas fikirlər

### 1. Tətbiq dizaynı
**Nədir:** İnternet forumu — istifadəçilər qeydiyyatdan keçib login olur, thread (mövzu) yaradır, post yazır; anonim istifadəçi yalnız oxuyur.

**Request formatı:** `http://<servername>/<handler-name>?<parameters>`
- Handler adı hierarxikdir: `/thread/read` → thread modulu, read handler-i
- Parametrlər URL query-dir: `id=123`
- Nümunə: `http://chitchat/thread/read?id=123`

**Server axını:** client request → **multiplexer** URL-i yoxlayıb düzgün handler-a yönləndirir → handler emal edir → template engine-ə data verir → HTML response.

### 2. Data model
4 struktur, PostgreSQL-də 4 cədvələ map olunur:
- **User** — forum istifadəçisi
- **Session** — cari login sessiyası
- **Thread** — forum mövzusu (söhbət)
- **Post** — thread daxilindəki mesaj

### 3. Multiplexer (main.go)
**Kitabdan kod nümunəsi:**
```go
package main

import (
    "net/http"
)

func main() {
    mux := http.NewServeMux()
    files := http.FileServer(http.Dir("/public"))
    mux.Handle("/static/", http.StripPrefix("/static/", files))

    mux.HandleFunc("/", index)
    mux.HandleFunc("/err", err)
    mux.HandleFunc("/login", login)
    mux.HandleFunc("/logout", logout)
    mux.HandleFunc("/signup", signup)
    mux.HandleFunc("/signup_account", signupAccount)
    mux.HandleFunc("/authenticate", authenticate)
    mux.HandleFunc("/thread/new", newThread)
    mux.HandleFunc("/thread/create", createThread)
    mux.HandleFunc("/thread/post", postThread)
    mux.HandleFunc("/thread/read", readThread)
    server := &http.Server{
        Addr:    "0.0.0.0:8080",
        Handler: mux,
    }
    server.ListenAndServe()
}
```

**Sub-kod izahı:**
- `mux := http.NewServeMux()` → multiplexer (router) yaradır — URL-i yoxlayıb handler-a yönləndirən kod
- `mux.HandleFunc("/", index)` → root URL-i `index` handler funksiyasına bağlayır; HandleFunc URL + funksiya adı qəbul edir
- `server := &http.Server{Addr: ..., Handler: mux}` → server struct-u mux ilə konfiqurasiya olunur; `ListenAndServe()` başladır

**Handler funksiyası:** sadəcə `(ResponseWriter, *http.Request)` imzalı Go funksiyasıdır. Eyni qovluqdakı bütün fayllar `package main` olduqda avtomatik birləşir (PHP/Python-dakı include YOXDUR).

### 4. Static fayl servisi
```go
files := http.FileServer(http.Dir("/public"))
mux.Handle("/static/", http.StripPrefix("/static/", files))
```

**Sub-kod izahı:**
- `http.FileServer(http.Dir("/public"))` → qovluqdan fayl verən handler
- `http.StripPrefix("/static/", files)` → URL-dən `/static/` prefiksini kəsir
- `/static/css/bootstrap.min.css` sorğusu → `<root>/css/bootstrap.min.css` faylı birbaşa qaytarılır (emal YOXDUR)

### 5. Cookie ilə access control
**Kitabdan kod nümunəsi (authenticate handler-i):**
```go
func authenticate(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    user, _ := data.UserByEmail(r.PostFormValue("email"))
    if user.Password == data.Encrypt(r.PostFormValue("password")) {
        session := user.CreateSession()
        cookie := http.Cookie{
            Name:     "_cookie",
            Value:    session.Uuid,
            HttpOnly: true,
        }
        http.SetCookie(w, &cookie)
        http.Redirect(w, r, "/", 302)
    } else {
        http.Redirect(w, r, "/login", 302)
    }
}
```

**Sub-kod izahı:**
- `r.ParseForm()` → form datanı parse edir
- `data.UserByEmail(...)` / `data.Encrypt(...)` → data paketinin funksiyaları (aşağıda)
- `user.CreateSession()` → DB-də Session yaradır (Uuid random unikal ID)
- `http.Cookie{Name, Value, HttpOnly}` → cookie struct-u; expiry yoxdur → **session cookie** (brauzer bağlananda silinir); `HttpOnly: true` → yalnız HTTP/HTTPS çıxışı (JavaScript yox)
- `http.SetCookie(w, &cookie)` → Set-Cookie header-i ilə response-a əlavə
- `http.Redirect(w, r, "/", 302)` → 302 redirect

**Session yoxlama utiliti (util.go):**
```go
func session(w http.ResponseWriter, r *http.Request) (sess data.Session, err error) {
    cookie, err := r.Cookie("_cookie")
    if err == nil {
        sess = data.Session{Uuid: cookie.Value}
        if ok, _ := sess.Check(); !ok {
            err = errors.New("Invalid session")
        }
    }
    return
}
```
- `r.Cookie("_cookie")` → request-dən cookie oxu
- `sess.Check()` → DB-də sessiya mövcudluğunu yoxla
- Dönüş: sessiya + error — handler-lar bunu istifadə edib public/private navbar seçir

**index handler-i (session ilə):**
```go
func index(writer http.ResponseWriter, request *http.Request) {
    threads, err := data.Threads()
    if err == nil {
        _, err := session(writer, request)
        if err != nil {
            generateHTML(writer, threads, "layout", "public.navbar", "index")
        } else {
            generateHTML(writer, threads, "layout", "private.navbar", "index")
        }
    }
}
```
- Session xətası = login olmayıb → public navbar; uğur = login → private navbar

### 6. Template-lər (html/template)
**3 fayl quruluşu:** `layout.html` + `public.navbar.html`/`private.navbar.html` + `index.html`

**Kitabdan kod nümunəsi (layout.html):**
```html
{{ define "layout" }}
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <title>ChitChat</title>
    <link href="/static/css/bootstrap.min.css" rel="stylesheet">
  </head>
  <body>
    {{ template "navbar" . }}
    <div class="container">
      {{ template "content" . }}
    </div>
    <script src="/static/js/jquery-2.1.1.min.js"></script>
  </body>
</html>
{{ end }}
```

**Sub-kod izahı:**
- `{{ define "layout" }} ... {{ end }}` → bu hissənin "layout" adlı template olduğunu bildirir
- `{{ template "navbar" . }}` → başqa template-i bu mövqeyə daxil et; `.` (nöqtə) — data-nı alt template-ə ötürür
- Actions = `{{ }}` içindəki annotasiyalar (Mustache/CTemplate-ə bənzər)

**Parse və execute:**
```go
templates := template.Must(template.ParseFiles(private_tmpl_files...))
templates.ExecuteTemplate(writer, "layout", threads)
```
- `template.ParseFiles(files...)` → faylları parse edib template dəsti yaradır
- `template.Must(...)` → xəta olsa panic edir (error yoxlama boilerplate-i azaldır)
- `ExecuteTemplate(w, "layout", threads)` → **layout** template-i icra olunur (o navbar + content-i daxil etdiyi üçün yalnız onu çağırmaq bəsdir — digərlərini çağırmaq yarım HTML verərdi)

**index.html — data iterasiyası:**
```html
{{ define "content" }}
<p class="lead">
  <a href="/thread/new">Start a thread</a> or join one below!
</p>
{{ range . }}
  <div class="panel panel-default">
    <div class="panel-heading">
      <span class="lead"> {{ .Topic }}</span>
    </div>
    <div class="panel-body">
      Started by {{ .User.Name }} - {{ .CreatedAtDate }} - {{ .NumReplies }} posts.
      <div class="pull-right">
        <a href="/thread/read?id={{.Uuid }}">Read more</a>
      </div>
    </div>
  </div>
{{ end }}
{{ end }}
```

**Sub-kod izahı:**
- `{{ range . }}` → kök nöqtə (.) = handler-dən ötürülən `threads` slice-ı; iterasiya edir
- `{{ .Topic }}` → range daxilində nöqtə = cari Thread; sahə adı BÖYÜK hərflə başlamalı (exported)
- `{{ .User.Name }}` → **metod** tərəfindən qaytarılan data
- `{{ .NumReplies }}` → Thread-də sahə yoxdur — **metod** çağırılır!

### 7. generateHTML — helper funksiya
**Kitabdan kod nümunəsi:**
```go
func generateHTML(w http.ResponseWriter, data interface{}, fn ...string) {
    var files []string
    for _, file := range fn {
        files = append(files, fmt.Sprintf("templates/%s.html", file))
    }
    templates := template.Must(template.ParseFiles(files...))
    templates.ExecuteTemplate(writer, "layout", data)
}
```

**Sub-kod izahı:**
- `data interface{}` → **empty interface** — hər tipi qəbul edir (interfeyslər metod dəstidir; boş dəst = hər tip uyğun gəlir — statik tipli dildə dinamiklik!)
- `fn ...string` → **variadic** — istənilən sayda template adı; variadic sonuncu parametr olmalıdır

### 8. PostgreSQL quraşdırılması
- **Ubuntu:** `sudo apt-get install postgresql postgresql-contrib` → `sudo su postgres` → `createuser --interactive` → `createdb <AD>`
- **Mac:** Postgres.app drag&drop
- **Windows:** EnterpriseDB installer (pgAdmin III ilə)

**Cədvəllər (setup.sql):**
```sql
create table users (
    id         serial primary key,
    uuid       varchar(64) not null unique,
    name       varchar(255),
    email      varchar(255) not null unique,
    password   varchar(255) not null,
    created_at timestamp not null
);
create table sessions (
    id         serial primary key,
    uuid       varchar(64) not null unique,
    email      varchar(255),
    user_id    integer references users(id),
    created_at timestamp not null
);
create table threads (
    id         serial primary key,
    uuid       varchar(64) not null unique,
    topic      text,
    user_id    integer references users(id),
    created_at timestamp not null
);
create table posts (
    id         serial primary key,
    uuid       varchar(64) not null unique,
    body       text,
    user_id    integer references users(id),
    thread_id  integer references threads(id),
    created_at timestamp not null
);
```
İcra: `psql -f setup.sql -d chitchat`

### 9. Database əlaqəsi (data paketi)
**Paket strukturu:** `data/` qovluğu — `thread.go`, `user.go`, `data.go`; adı `package data`. İstifadə: `import "github.com/sausheong/gwp/Chapter_2_Go_ChitChat/chitchat/data"` → `data.Thread`.

**Kitabdan kod nümunəsi (Thread struct + DB pool):**
```go
// data.go — global connection pool:
var Db *sql.DB

func init() {
    var err error
    Db, err = sql.Open("postgres", "dbname=chitchat sslmode=disable")
    if err != nil {
        log.Fatal(err)
    }
    return
}

// thread.go:
type Thread struct {
    Id        int
    Uuid      string
    Topic     string
    UserId    int
    CreatedAt time.Time
}
```

**Sub-kod izahı:**
- `sql.Open("postgres", "dbname=chitchat sslmode=disable")` → **connection pool** (`*sql.DB`) — hər sorğu üçün yeni bağlantı deyil
- `init()` → proqram başlarkən Db avtomatik yaradılır
- Thread struct-u `threads` cədvəlinə uyğun gəlir

**Threads funksiyası — bütün thread-ləri oxu:**
```go
func Threads() (threads []Thread, err error) {
    rows, err := Db.Query("SELECT id, uuid, topic, user_id, created_at FROM threads ORDER BY created_at DESC")
    if err != nil {
        return
    }
    for rows.Next() {
        th := Thread{}
        if err = rows.Scan(&th.Id, &th.Uuid, &th.Topic, &th.UserId, &th.CreatedAt); err != nil {
            return
        }
        threads = append(threads, th)
    }
    rows.Close()
    return
}
```

**4 addımlıq DB sorğu paterni:**
1. Connection pool ilə qoşul
2. SQL query göndər → rows qayıdır
3. Struct yarat
4. Rows üzrə iterasiya edib Scan ilə struct-a doldur

**Metod nümunəsi — NumReplies (template-dən çağırılır):**
```go
func (thread *Thread) NumReplies() (count int) {
    rows, err := Db.Query("SELECT count(*) FROM posts where thread_id = $1", thread.Id)
    if err != nil {
        return
    }
    for rows.Next() {
        if err = rows.Scan(&count); err != nil {
            return
        }
    }
    rows.Close()
    return
}
```
- `func (thread *Thread) NumReplies()` → Thread-ə bağlı metod; receiver = Thread
- `$1` → PostgreSQL placeholder (SQL injection qorunması)
- Template-də `{{ .NumReplies }}` bu metodu çağırır — sahə də, metod da eyni sintaksislə!

**Data layer fəlsəfəsi:** handler-lar heç vaxt birbaşa SQL yazmır — struct + funksiya/metod kombinasiyası data qatını təşkil edir. ORM kitabxanası yoxdur — "no magic, just simple, straightforward code".

### 10. Server başlatma
```go
server := &http.Server{
    Addr:    "0.0.0.0:8080",
    Handler: mux,
}
server.ListenAndServe()
```
CLI: `go build` → `./chitchat` → http://localhost:8080

## Tam axın (recap)
```
1. Client request göndərir
2. Multiplexer düzgün handler-a yönləndirir
3. Handler request-i emal edir
4. Data lazım olanda data struct-larından istifadə edir
5. Model funksiya/metodları ilə DB-yə qoşulur
6. Handler template engine-i işə salır (data ilə)
7. Template faylları parse olunub data ilə birləşərək HTML yaradır
8. HTML response ilə client-ə qayıdır
```

## Əsas terminlər
- Multiplexer / ServeMux (çoxlayıcı / router)
- Handler Function (emalçı funksiya)
- FileServer (fayl serveri)
- StripPrefix (prefiks kəsm)
- Session Cookie (sessiya kökəyi)
- HttpOnly (yalnız-HTTP bayrağı)
- Template Engine (şablon mühərriki)
- Action (`{{ }}` — şablon əmri)
- define / template / range actions
- Empty Interface (`interface{}`)
- Variadic Parameter (`...T`)
- Connection Pool (`sql.DB`)
- DDL (Data Definition Language)
- Data Layer (data qatı)
- Receiver (metodun bağlı olduğu tip)
- Placeholder (`$1` — parametr yer tutucusu)

## Praktik nəticə
- Router-i `http.NewServeMux()` ilə öz mux-unu yarat — default mux-da bütün handler-lar qlobal ad məkanını paylaşır.
- Static fayllar üçün `FileServer + StripPrefix` kombinasiyası standart qəndaşlıqdır.
- Login vəziyyətini session cookie (Uuid) + DB cədvəli ilə yoxla; HttpOnly sayəsində JS cookie-ə çıxa bilmir.
- `template.Must` parse xətalarını başlanğıcda partladır — runtime-da yarımçıq template yox.
- Handler-ları drum etmə: data əməliyyatları struct metodlarına, HTML yaratmağı `generateHTML` helper-inə ver.
- Template-də sahə ilə metodu ayırd etmək lazım deyil — `{{ .X }}` hər ikisini də çağırır; sahələr exported (böyük hərf) olmalıdır.

## Mənbə
Pages: 43-65 (PDF), book pages 22-44
