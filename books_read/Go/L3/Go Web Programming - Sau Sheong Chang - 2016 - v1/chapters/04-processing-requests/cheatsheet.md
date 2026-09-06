# Chapter 4 — Processing requests Cheatsheet

## `r.ParseForm()` + `r.Form` / `r.PostForm`

**Nə edir:** Sorğunı parse edir, form datanı `Form` (URL+form) və `PostForm` (yalnız form) map-ləri ilə əldə edir.

**Fərqlər:** `Form` — URL query + HTML form birləşimi; `PostForm` — yalnız URL-encoded form (multipart-də boş olur).

```go
r.ParseForm()
fmt.Fprintln(w, r.Form)    // map[thread:[123] hello:[sau sheong world] post:[456]]
fmt.Fprintln(w, r.PostForm) // map[hello:[sau sheong] post:[456]]
```

**Mənbə:** Chapter 4, page 78

---

## `r.FormValue("key")` / `r.PostFormValue("key")`

**Nə edir:** Sorğudan key üzrə dəyər alır. Öz-özlə `ParseForm`/`ParseMultipartForm` çağırarıq (əlavə iş yox).

**Gotyaya dikkat:** `FormValue` yalnız ilk dəyəri qaytarır; `multipart/form-data`-da `PostFormValue` boş qalır (MultipartForm istifadə edin).

**Mənbə:** Chapter 4, page 79

---

## `r.FormFile("uploaded")`

**Nə edir:** Multipart formdan fayl götürür. `ParseMultipartForm` çağırmaya ehtiyac yox (avtomatik).

```go
file, _, err := r.FormFile("uploaded")
data, _ := ioutil.ReadAll(file)
fmt.Fprintln(w, string(data))
```

**Mənbə:** Chapter 4, page 81

---

## `r.ParseMultipartForm(1024)` + `r.MultipartForm.File["name"][0]`

**Nə edir:** Multipart formu parse edir (1024 bayt RAM buffer, istirakdən çox diskdə). `MultipartForm`-dan `FileHeader` alır.

```go
r.ParseMultipartForm(1024)
fileHeader := r.MultipartForm.File["uploaded"][0]
file, _ := fileHeader.Open()
data, _ := ioutil.ReadAll(file)
```

**Mənbə:** Chapter 4, page 80 — `maxMemory` RAM-a yüklənən maksimum bayt sayı

---

## `r.MultipartForm` structure

**Nə edir:** İki map-dan ibarət struct: `{map[string][]string, map[string][]*FileHeader}`.

- Birinci map — form sahələri (key → string slice)
- İkinci map — fayllar (key → FileHeader slice) — üçüncü "hissə" boş (fayl map-i)

**Mənbə:** Chapter 4, page 79

---

## `w.Write([]byte(str))`

**Nə edir:** Cavab body-sinə bayt massivini yazar. İlk 512 baytdan `Content-Type` avtomatik ayrışdırılır.

```go
w.Write([]byte(html))
```

**Mənbə:** Chapter 4, page 83

---

## `w.WriteHeader(501)`

**Nə edir:** HTTP status kodunu təyin edir (default 200). `WriteHeader` çağıldıqdan sonra header dəyişdirilə bilməz.

```go
func writeHeaderExample(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(501)
    fmt.Fprintln(w, "No such service, try next door")
}
```

**Mənbə:** Chapter 4, page 85

---

## `w.Header().Set("Location", "http://google.com")` + `w.WriteHeader(302)`

**Nə edir:** Redirect üçün `Location` header + 302 status. Header **WriteHeader'dan əvvəl** təyin edin.

```go
w.Header().Set("Location", "http://google.com")
w.WriteHeader(302)
```

**Mənbə:** Chapter 4, page 86

---

## `type Cookie struct{...}` + `http.SetCookie(w, &c)`

**Nə edir:** Cookie strukturu və response-a əlavə etmək.

**Field-lər:**
```go
type Cookie struct {
    Name string
    Value string
    Expires time.Time
    MaxAge int
    Secure bool
    HttpOnly bool
}
```
- `HttpOnly: true` → JS-dən (XSS) qorunur
- `Expires` təyin etmək → persistent; təyin etməmək → session cookie

**Mənbə:** Chapter 4, page 88

---

## `r.Cookie("name")` / `r.Cookies()`

**Nə edir:** Sorğudan cookie oxumaq: `Cookie(name)` — ilk uyğun cookie; `Cookies()` — bütün cookie-lər slice-u.

```go
c1, err := r.Cookie("first_cookie")
cs := r.Cookies()
```

**Mənbə:** Chapter 4, page 91

---

## Flash cookie silmə: `MaxAge: -1` + `Expires: time.Unix(1, 0)`

**Nə edir:** Eyni ada sahib cookie ilə əvəz edərək browser-in cookie-ni silməsini söndürür.

```go
rc := http.Cookie{
    Name: "flash",
    MaxAge: -1,
    Expires: time.Unix(1, 0),
}
http.SetCookie(w, &rc)
```

**Mənbə:** Chapter 4, page 93

---

## `json.Marshal(post)` → `w.Write(json)`

**Nə edir:** Struct-ı JSON-ə çevirir və response body-yə yazar.

```go
type Post struct {
    User string
    Threads []string
}
w.Header().Set("Content-Type", "application/json")
json, _ := json.Marshal(post)
w.Write(json)
```

**Mənbə:** Chapter 4, page 87
