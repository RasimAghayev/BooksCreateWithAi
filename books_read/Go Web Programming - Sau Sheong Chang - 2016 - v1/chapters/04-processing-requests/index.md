# Chapter 4 — Processing requests

## Bu chapter nədən bəhs edir?

Request-in emalı və response göndərilməsi: Request struct (URL, Header, Body, Form/PostForm/MultipartForm), HTML form enctype növləri, fayl yükləmə, JSON body, ResponseWriter-in 3 metodu (Write/WriteHeader/Header), cookie-lər (struct, göndərmə/oxuma, flash mesajlar).

## Əsas fikirlər

### 1. Request struct
**Nədir:** Client-dən gələn HTTP request-in parse edilmiş təmsili — text deyil, strukturlaşdırılmış Go tipi.

**Əsas sahələr:** URL, Header, Body, Form, PostForm, MultipartForm. Request həm serverə gələn, həm clientdən gedən sorğunu təmsil edə bilər.

**url.URL struct:**
```go
type URL struct {
    Scheme   string
    Opaque   string
    User     *Userinfo
    Host     string
    Path     string
    RawQuery string
    Fragment string
}
```
- Ümumi forma: `scheme://[userinfo@]host/path[?query][#fragment]`
- `RawQuery` → xam query string (`id=123&thread_id=456`); parse etmək üçün Form sahələri istifadə olunur
- **Fragment boş qalır** — brauzer fragment-i serverə göndərməzdən əvvəl kəsir (Go-nun problemi deyil!); amma HTTP client kitabxanaları/JS framework-ləri fragment göndərə bilər

### 2. Header — map[string][]string
**Nədir:** Header = map; açar string, dəyər **string slice**-ı (eyni adlı bir neçə header dəyəri ola bilər).

**4 metod:** Add, Delete, Get, Set.
- **Set** → dəyəri sıfırdan yazır (köhnə sliceı əvəz edir)
- **Add** → mövcud slice-a element əlavə edir

**Oxuma fərqi:**
```go
h := r.Header["Accept-Encoding"]   // [gzip, deflate] — map çıxışı (slice)
h := r.Header.Get("Accept-Encoding") // "gzip, deflate" — bir string (comma-delimited)
```

### 3. Request body — io.ReadCloser
**Nədir:** Body = Reader + Closer interfeyslərinin birləşməsi (io.ReadCloser).

**Kitabdan kod nümunəsi:**
```go
func body(w http.ResponseWriter, r *http.Request) {
    len := r.ContentLength
    body := make([]byte, len)
    r.Body.Read(body)
    fmt.Fprintln(w, string(body))
}
```
- `r.ContentLength` → body-nin ölçüsü; `make([]byte, len)` → buffer; `Read` → oxu
- Test (cURL): `curl -id "first_name=sausheong&last_name=chang" 127.0.0.1:8080/body`
- Adətən xam body oxunmur — FormValue/FormFile hazırdır

### 4. HTML form və enctype
**Form HTML:**
```html
<form action="/process" method="post">
  <input type="text" name="first_name"/>
  <input type="text" name="last_name"/>
  <input type="submit"/>
</form>
```

**2 əsas enctype:**

**a) `application/x-www-form-urlencoded` (default):** body = URL-encoded query string:
```
first_name=sau%20sheong&last_name=chang
```
Sadə mətn üçün — daha sadə və effektiv.

**b) `multipart/form-data`:** hər cüt ayrı MIME hissə olur:
```
------WebKitFormBoundaryMPNjKpeO9cLiocMw
Content-Disposition: form-data; name="first_name"
sau sheong
------WebKitFormBoundaryMPNjKpeO9cLiocMw
Content-Disposition: form-data; name="last_name"
chang
------WebKitFormBoundaryMPNjKpeO9cLiocMw--
```
Böyük data/fayl yükləmə üçün (Base64 ilə binary də mümkün).

**GET forması:** body yoxdur — data URL-də name-value cütləri kimi.

### 5. Form / PostForm / MultipartForm — 3 data sahəsi
**Ümumi alqoritm:**
1. `ParseForm` yaxud `ParseMultipartForm` çağır
2. Müvafiq sahəyə çıxış et

**Kitabdan kod nümunəsi (Form):**
```go
func process(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    fmt.Fprintln(w, r.Form)
}
```
Formda `hello=sau sheong` (POST) + URL-də `?hello=world&thread=123` olduqda:
```
map[thread:[123] hello:[sau sheong world] post:[456]]
```
- `r.Form` → **URL + body** dəyərlərinin birləşməsi; dəyərlər URL-decode olunur
- Eyni açarlı dəyərlər bir slice-da; **form dəyəri URL dəyərindən əvvəl**

**r.PostForm:** yalnız body (URL query YOX):
```
map[post:[456] hello:[sau sheong]]
```

**MultipartForm (multipart/form-data üçün):**
```go
r.ParseMultipartForm(1024)
fmt.Fprintln(w, r.MultipartForm)
// &{map[hello:[sau sheong] post:[456]] map[]}
```
- `ParseMultipartForm(1024)` → 1024 byte çıxar (giriş parametri); lazım olanda ParseForm-u da çağırır
- Nəticə: **2 mapli struct** — 1-ci value-lar, 2-ci (boş) fayllar üçün

**Qısayollar:**
- `r.FormValue("hello")` → ParseForm-u özü çağırır, ilk dəyəri verir
- `r.PostFormValue("hello")` → PostFun üçün eyni

**GOTCHA (vacib!):** `FormValue`/`PostFormValue` multipart body-dən dəyər vermir — yalnız Form/PostForm sahələrinə baxır. multipart/form-data + FormValue → URL-dən `world` gəlir, body-dən heç nə!

**Kitabdan cədvəl (Türkçe şərh):**
| Sahə | Çağırılacaq metod | URL-dən | Form-dan (urlencoded) | Multipart |
|---|---|---|---|---|
| Form | ParseForm | ✅ | ✅ | ✅ |
| PostForm | ParseForm | ❌ | ✅ | ❌ |
| MultipartForm | ParseMultipartForm | ❌ | ✅ | ✅ |
| FormValue (metod) | avtomatik | ✅ | ✅ | ❌ |
| PostFormValue (metod) | avtomatik | ❌ | ✅ | ❌ |

### 6. Fayl yükləmə
**Kitabdan kod nümunəsi (uzun yol):**
```go
func process(w http.ResponseWriter, r *http.Request) {
    r.ParseMultipartForm(1024)
    fileHeader := r.MultipartForm.File["uploaded"][0]
    file, err := fileHeader.Open()
    if err == nil {
        data, err := ioutil.ReadAll(file)
        if err == nil {
            fmt.Fprintln(w, string(data))
        }
    }
}
```

**Qısa yol — FormFile:**
```go
func process(w http.ResponseWriter, r *http.Request) {
    file, _, err := r.FormFile("uploaded")
    if err == nil {
        data, err := ioutil.ReadAll(file)
        if err == nil {
            fmt.Fprintln(w, string(data))
        }
    }
}
```
- `r.MultipartForm.File["uploaded"][0]` → FileHeader slice-dan birincisi
- `r.FormFile("uploaded")` → parse + ilk fayl birbaşa (tək fayl üçün ən sürətli)

### 7. POST + JSON body
**Problem:** Angular POST-ları `application/json` Content-Type ilə göndərir → `ParseForm` **heç nə qaytarmır** (yalnız form formatlarını parse edir). JQuery isə `application/x-www-form-urlencoded` işlədir → ParseForm işləyir.

**Dərs:** framework-lərin arxasındakı detalları bilmək vacibdir — "framework seams show at the joints". JSON body üçün: `json.NewDecoder(r.Body).Decode(&v)` (Ch 7-də ətraflı).

### 8. ResponseWriter — 3 metod
**Nədir:** Handler-in response yaratmaq üçün istifadə etdiyi interfeys; arxada nonexported `http.response` struct-u durur (interfeysdən başqa yol yoxdur). Hər 2 ServeHTTP parametri əslində reference ilə ötürülür — ResponseWriter interfeysi pointer-ın üstünə qabıqdır.

**a) Write — body yaz:**
```go
func writeExample(w http.ResponseWriter, r *http.Request) {
    str := `<html>
<head><title>Go Web Programming</title></head>
<body><h1>Hello World</h1></body>
</html>`
    w.Write([]byte(str))
}
```
- Content-Type təyin edilməyibsə ilk **512 byte-a** görə avtomatik aşkarlanır (`text/html; charset=utf-8` çıxır)

**b) WriteHeader — status kodu yaz:**
```go
func writeHeaderExample(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(501)
    fmt.Fprintln(w, "No such service, try next door")
}
```
- Adı aldadıcıdır: header YOX, **status kodu** yazır
- Çağırıldıqdan sonra hələ body yazmaq olar, amma header dəyişmək OLMAZ
- Çağırılmasa avtomatik 200 OK gedir

**c) Header — header map döndür:**
```go
func headerExample(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Location", "http://google.com")
    w.WriteHeader(302)
}
```
- Redirect: `Location` header + 302 status
- **Sıra vacibdir:** Location WriteHeader-dən ƏVVƏL yazılmalıdır

**JSON cavabı:**
```go
func jsonExample(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    post := &Post{
        User:    "Sau Sheong",
        Threads: []string{"first", "second", "third"},
    }
    json, _ := json.Marshal(post)
    w.Write(json)
}
// → {"User":"Sau Sheong","Threads":["first","second","third"]}
```
- Content-Type manual set et, sonra marshalled JSON-u Write ilə göndər

### 9. Cookie struct və göndərmə
**Nədir:** Cookie — client-də saxlanan kiçik məlumat; HTTP-nin stateless-liyini aradan qaldırmanın ən populyar yolu. İki sinif: **session cookie** (brauzer bağlananda silinir) və **persistent cookie** (müddət bitənə qədər qalır).

**Cookie struct:**
```go
type Cookie struct {
    Name       string
    Value      string
    Path       string
    Domain     string
    Expires    time.Time
    RawExpires string
    MaxAge     int
    Secure     bool
    HttpOnly   bool
    Raw        string
    Unparsed   []string
}
```
- Expires boş → session cookie
- Expires vs MaxAge: Expires dəqiq vaxt, MaxAge saniyə müddət; HTTP 1.1-də Expires deprecate olunsa da bütün brauzerlər dəstəkləyir; MaxAge IE 6/7/8-də yoxdur → **praqmatik həll: ikisini də qoy** (yaxud yalnız Expires)

**Göndərmə — 2 üsul:**
```go
c1 := http.Cookie{Name: "first_cookie", Value: "Go Web Programming", HttpOnly: true}
c2 := http.Cookie{Name: "second_cookie", Value: "Manning Publications Co", HttpOnly: true}

// Üsul 1 — manual header:
w.Header().Set("Set-Cookie", c1.String())
w.Header().Add("Set-Cookie", c2.String())

// Üsul 2 — qısayol:
http.SetCookie(w, &c1)
http.SetCookie(w, &c2)
```
- `c1.String()` → cookie-ni Set-Cookie üçün serializə edir
- `Set` birinci, `Add` ikinci cookie üçün (slice-a əlavə)
- `http.SetCookie` cookie-ni **referansla** qəbul edir

### 10. Cookie oxuma
```go
// Xam header (parse etməlisən):
h := r.Header["Cookie"]
// [first_cookie=Go Web Programming; second_cookie=Manning Publications Co]

// Adlı cookie:
c1, err := r.Cookie("first_cookie")   // yoxdursa err (http.ErrNoCookie)

// Bütün cookie-lər:
cs := r.Cookies()   // []*http.Cookie slice
```

### 11. Flash mesajlar (cookie silmə texnikası)
**Nədir:** Müvəqqəti (transient) informativ mesaj — səhifə yenilənəndə bir də göstərilmir.

**Kitabdan kod nümunəsi:**
```go
func setMessage(w http.ResponseWriter, r *http.Request) {
    msg := []byte("Hello World!")
    c := http.Cookie{
        Name:  "flash",
        Value: base64.URLEncoding.EncodeToString(msg),
    }
    http.SetCookie(w, &c)
}

func showMessage(w http.ResponseWriter, r *http.Request) {
    c, err := r.Cookie("flash")
    if err != nil {
        if err == http.ErrNoCookie {
            fmt.Fprintln(w, "No message found")
        }
    } else {
        rc := http.Cookie{
            Name:    "flash",
            MaxAge:  -1,
            Expires: time.Unix(1, 0),   // keçmiş vaxt!
        }
        http.SetCookie(w, &rc)
        val, _ := base64.URLEncoding.DecodeString(c.Value)
        fmt.Fprintln(w, string(val))
    }
}
```

**Sub-kod izahı:**
- `base64.URLEncoding.EncodeToString(msg)` → cookie dəyəri URL-safe olmalıdır (space və % cookie-də problem); decode: `DecodeString`
- **Cookie silmə:** eyni adda + `MaxAge: -1` + `Expires: time.Unix(1, 0)` (1970!) → brauzer köhnə cookie-ni əvəz edib silir
- Axın: `/set_message` → cookie yaranır → `/show_message` → mesaj göstərilir + cookie silinir → yenidən `/show_message` → "No message found"

## Əsas terminlər
- Request / Response struct
- url.URL / RawQuery
- Header (map[string][]string)
- io.ReadCloser (Reader + Closer)
- enctype (form content type)
- application/x-www-form-urlencoded
- multipart/form-data
- Form / PostForm / MultipartForm
- ParseForm / ParseMultipartForm
- FormValue / PostFormValue / FormFile
- FileHeader
- ResponseWriter (Write / WriteHeader / Header)
- Content-Type sniffing (512 byte)
- Session / Persistent Cookie
- Expires / MaxAge
- SetCookie / ErrNoCookie
- Flash Message (müvəqqəti mesaj)
- Base64 URL Encoding

## Praktik nəticə
- Form data oxuyarkən adətən `r.FormValue` bəsdir; fərqli source-ları ayırd etmək lazımdırsa Form/PostForm/MultipartForm cədvəlinə bax.
- multipart/form-data + FormValue = boş — bu gotcha-nı bil; fayl üçün `r.FormFile`.
- Angular/JS client-lərdən gələn JSON POST-ları `ParseForm` tutmur — `json.NewDecoder(r.Body)` istifadə et.
- Header yazma sırası: `Header().Set` → `WriteHeader` → `Write` — WriteHeader-dən sonra header dəyişmir.
- Cookie silmə = eyni ad + MaxAge -1 + keçmiş Expires.
- Bütün brauzerlər üçün cookie müddətini hər iki sahə ilə (Expires + MaxAge) təyin et.

## Mənbə
Pages: 90-116 (PDF), book pages 69-95
