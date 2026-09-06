# Chapter 2 — Go ChitChat Cheatsheet

## `mux := http.NewServeMux()`

**Nə edir:** Yeni çoxləşdirici (ServeMux) yaradır. Sorğuları URL üzrə handler-lərə yönləndirir.

**Mənbə:** Chapter 2, page 48

---

## `mux.HandleFunc("/", index)`

**Nə edir:** `/` URL path-i üçün `index` handler funksiyasını qeyd edir.

**Parametrlər:**
- `"/"` → URL pattern
- `index` → `func(http.ResponseWriter, *http.Request)` imzalı handler

**Mənbə:** Chapter 2, page 48

---

## `mux.Handle("/static/", http.StripPrefix("/static/", files))`

**Nə edir:** `/static/` ilə başlayan sorğuları statik fayl serverinə yönləndirir, URL-dən `/static/` prefix-ni çıxarır.

**Alt parametrlər:**
- `http.FileServer(http.Dir(...))` → fayl xidmət edən handler
- `http.StripPrefix("/static/", files)` → prefix-i silib asıl path-ı qalan handler-a ötürür

**Nəcə işləyir:** Sorğu `/static/css/bootstrap.min.css` gələrsə, `/static/`-i çıxıb `css/bootstrap.min.css` faylını `public/` kataloğunda axtarır.

**Mənbə:** Chapter 2, page 49

---

## `func index(w http.ResponseWriter, r *http.Request)`

**Nə edir:** HTTP handler funksiyası. `ResponseWriter`-ə cavab yazır, `Request`-dən sorğu məlumatlarını oxur.

**Parametrlər:**
- `w` → `http.ResponseWriter` — cavabı client-ə göndərmək üçün
- `r` → `*http.Request` — sorğunun method, URL, header, body məlumatları

**Mənbə:** Chapter 2, page 50

---

## `http.SetCookie(w, &cookie)` / `r.Cookie("_cookie")`

**Nə edir:** Response-a cookie yazaır / REQUEST-dən cookie oxur.

**Cookie strukturu:**
```go
cookie := http.Cookie{
    Name: "_cookie",
    Value: session.Uuid,
    HttpOnly: true,
}
```
- `HttpOnly: true` → cookie JavaScript-dən (XSS ilə) qorunur
- `Value`-i təyin etmək, `Expiry` təyin etməmək → session cookie olur

**Mənbə:** Chapter 2, page 54

---

## `template.Must(template.ParseFiles(files...))`

**Nə edir:** Template faylları parse edir. `Must` xəta başınsa panic atır.

**Alt funksiyalar:**
- `template.ParseFiles(files...)` → `*template.Template` və `error` qaytarır
- `template.Must(t, err)` → `err != nil`-dən panic atır, `err == nil`-da `t`-i qaytarır

**Mənbə:** Chapter 2, page 51

---

## `templates.ExecuteTemplate(w, "layout", data)`

**Nə edir:** `"layout"` adlı şablonu `data`-nı birləşdirərək `ResponseWriter`-ə HTML yazır.

**Parametrlər:**
- `w` → `http.ResponseWriter`
- `"layout"` → icra ediləcək template adı (define ilə təyin edilir)
- `data` → template-ə ötürülən məlumat (`. ` ilə template-də çıxar)

**Mənbə:** Chapter 2, page 52

---

## `func generateHTML(w http.ResponseWriter, data interface{}, fn ...string)`

**Nə edir:** Generic HTML generasiya funksiyası. Template fayl adlarından və data-dan HTML yaradır.

**Parametrlər:**
- `data interface{}` → boş interfeys, istənilən tip qəbul edir
- `fn ...string` → variadic parametr (sıfır və ya daha çox template fayl adı)

**Mənbə:** Chapter 2, page 55

---

## `func (thread *Thread) NumReplies() (count int)`

**Nə edir:** `Thread` struct-ının metodu. Database-də o thread-a aid post sayını qaytarır.

**Necə işləyir:**
```go
func (thread *Thread) NumReplies() (count int) {
    rows, err := Db.Query("SELECT count(*) FROM posts where thread_id = $1", thread.Id)
    ...
}
```
- `(thread *Thread)` → receiver (alıcı), pointer-receiver metod
- Template-də `{{ .NumReplies }}` ilə çıphır

**Mənbə:** Chapter 2, page 61

---

## `var Db *sql.DB` + `func init()`

**Nə edir:** Global database connection pool-u və proqram başladıqda DB-ə qoşulma.

```go
var Db *sql.DB

func init() {
    var err error
    Db, err = sql.Open("postgres", "dbname=chitchat sslmode=disable")
    if err != nil {
        log.Fatal(err)
    }
}
```

**Alt parametrlər:**
- `*sql.DB` → connection pool-un təmsiliyyatı (birbaşına connection deyil)
- `"postgres"` → driver adı
- `"dbname=chitchat sslmode=disable"` → connection string

**Mənbə:** Chapter 2, page 60

---

## `func Threads() (threads []Thread, err error)`

**Nə edir:** Bütün thread-ləri database-dən oxuyub slice-a toplayır.

```go
func Threads() (threads []Thread, err error){
    rows, err := Db.Query("SELECT id, uuid, topic, user_id, created_at FROM threads ORDER BY created_at DESC")
    if err != nil { return }
    for rows.Next() {
        th := Thread{}
        if err = rows.Scan(&th.Id, &th.Uuid, &th.Topic, &th.UserId, &th.CreatedAt); err != nil { return }
        threads = append(threads, th)
    }
    rows.Close()
    return
}
```

**Database query addımları:**
1. `Db.Query(...)` ilə sorğu göndər
2. `rows.Next()` ilə iterasiya et
3. `rows.Scan(...)` ilə field-ləri struct-a yığıl
4. `rows.Close()` ilə bağlantı serbest burax

**Mənbə:** Chapter 2, page 61

---

## `go build` + `./chitchat`

**Nə edir:** ChitChat serverini kompilyasiya edir və işə salır.

**Əmri:**
```bash
go build
./chitchat
```
- `go build` → `.go` faylları `chitchat` binary-ına kompilyasiya edir
- `./chitchat` → 8080 portunda serveri `0.0.0.0`-da işə salır

**Mənbə:** Chapter 2, page 64

---

## `server := &http.Server{Addr: "0.0.0.0:8080", Handler: mux}`

**Nə edir:** Server konfiqurasiyasını qurur və `ListenAndServe` ilə dinləyərək sorğuları `mux`-a yönləndirir.

**Field-lər:**
- `Addr` → istidlak olunacaq ünvan və port (`0.0.0.0` = hamı interface-lər)
- `Handler` → istifadə olunacaq çoxləşdirici (mux)

**Mənbə:** Chapter 2, page 64
