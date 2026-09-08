# Chapter 9 — HTTP (səh. 183-207)

## Bu fəsil nədən bəhs edir?

`net/http` paketi ilə HTTP klient (Get/Post/PostForm, custom request,
headers, cookies, CookieJar) və server (ServerMux, handler-lər) tərəfləri,
plus middleware pattern-inin qurulması.

## Əsas fikirlər

### 1. GET request
**Nədir:** `http.Get(url)` — ən sadə sorğu; `*Response` + error qaytarır.

**Kitabdan kod nümunəsi:**
```go
resp, err := http.Get("https://httpbin.org/get")
if err != nil {
    panic(err)
}

fmt.Println(resp.Status)                  // 200 OK
fmt.Println(resp.Header["Content-Type"])  // [application/json]

defer resp.Body.Close()                    // body Reader-dir — bağlanmalıdır!
buf := bufio.NewScanner(resp.Body)
for buf.Scan() {
    fmt.Println(buf.Text())
}
```

**Sub-kod izahı:**
- `resp.Status` → status sətri (RFC7231: "200 OK")
- `resp.Header` → map; header adı ilə çıxılır
- `resp.Body` → `io.Reader` — Chapter 7 bilikləri burada tətbiq olunur
- `defer resp.Body.Close()` → deskriptor itkisinin qarşısını alır

### 2. POST request
**Nədir:** `http.Post(url, contentType, body)` — body `io.Reader` tələb edir.

```go
bodyRequest := []byte(`{"user": "john","email":"john@gmail.com"}`)
bufferBody := bytes.NewBuffer(bodyRequest)

url := "https://httpbin.org/post"
resp, err := http.Post(url, "application/json", bufferBody)
if err != nil {
    panic(err)
}
defer resp.Body.Close()
// body oxunması — GET ilə eyni
```

**PostForm — url-encoded formalar:**
```go
resp, err := http.PostForm("https://httpbin.org/post",
    url.Values{"user": {"john"}, "email": {"john@gmail.com"}})
```
- Content-type avtomatik `application/x-www-form-urlencoded`
- Server tərəfdə data `form` sahəsinə düşür (raw POST-da `data` sahəsinə)

### 3. Custom request + Client (PUT/DELETE və s.)
**Nəyə lazımdır:** GET/POST-dan başqa metodlar, xüsusi header-lər, timeout.

**Kitabdan kod nümunəsi:**
```go
bodyRequest := []byte(`{"user": "john","email":"john@gmail.com"}`)
bufferBody := bytes.NewBuffer(bodyRequest)
url := "https://httpbin.org/put"

header := http.Header{}
header.Add("Content-type", "application/json")
header.Add("X-Custom-Header", "somevalue")
header.Add("User-Agent", "safe-the-world-with-go")

req, err := http.NewRequest(http.MethodPut, url, bufferBody)
if err != nil {
    panic(err)
}
req.Header = header

client := http.Client{
    Timeout: time.Second * 5,
}
resp, err := client.Do(req)   // request-in göndərilməsi
if err != nil {
    panic(err)
}
defer resp.Body.Close()
```

**Sub-kod izahı:**
- `http.NewRequest(metod, url, body)` → `*Request` qurur
- `req.Header = header` → header-lər request yarandıqdan SONRA təyin edilir
- `http.Client{Timeout: ...}` → konfiqurasiya oluna bilən klient (timeout,
  redirect policy, cookie jar); `client.Do(req)` göndərir
- `http.MethodPut` → "PUT" konstantası (həmçinin DELETE, PATCH...)

### 4. HTTP server — ServerMux + handler
**Nədir:** `ServerMux` — URL pattern-lərini handler-lərlə əlaqələndirən
 multiplexer. Handler imzası: `func(ResponseWriter, *Request)`.

**Kitabdan kod nümunəsi:**
```go
func info(w http.ResponseWriter, r *http.Request) {
    for name, headers := range r.Header {
        fmt.Println(name, headers)
    }
    w.Write([]byte("Perfect!!!"))
    return
}

func main() {
    http.HandleFunc("/info", info)
    panic(http.ListenAndServe(":8090", nil))  // bloklanır — sonsuz gözləyir
}
```

**Test:**
```bash
>>> curl -H "Header1: Value1" :8090/info
Perfect!!!
```

**Handler interfeysi ilə (struct əsaslı):**
```go
type MyHandler struct {}

func (c *MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    switch r.RequestURI {
    case "/hello":
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("goodbye\n"))
    case "/goodbye":
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("hello\n"))
    default:
        w.WriteHeader(http.StatusBadRequest)
    }
}

func main() {
    handler := MyHandler{}
    http.ListenAndServe(":8090", &handler)   // bütün URI-lər bu handler-ə
}
```
- `ServeHTTP` → `Handler` interfeysinin tələbi; URL seçimi öz əl ilə

### 5. Cookies
**Nədir:** `http.Cookie` tipi; request-ə əlavə, response-dən oxunur.

**Kitabdan kod nümunəsi (sayğac cookie-si):**
```go
// SERVER tərəfi:
func cookieSetter(w http.ResponseWriter, r *http.Request) {
    counter, err := r.Cookie("counter")     // cookie oxu
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
    value, err := strconv.Atoi(counter.Value)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
    value = value + 1
    newCookie := http.Cookie{
        Name:  "counter",
        Value: strconv.Itoa(value),
    }
    http.SetCookie(w, &newCookie)           // response-a yaz
    w.WriteHeader(http.StatusOK)
}

// KLİENT tərəfi:
func main() {
    http.HandleFunc("/cookie", cookieSetter)
    go http.ListenAndServe(":8090", nil)

    url := "http://localhost:8090/cookie"
    req, _ := http.NewRequest("GET", url, nil)

    client := http.Client{}
    c := http.Cookie{
        Name: "counter", Value: "1", Domain: "127.0.0.1",
        Path: "/", Expires: time.Now().AddDate(1, 0, 0)}
    req.AddCookie(&c)                       // request-ə cookie əlavə et

    fmt.Println("—>", req.Header)           // map[Cookie:[counter=1]]
    resp, err := client.Do(req)
    fmt.Println("<—", resp.Header)           // Set-Cookie:[counter=2]
}
```

**Sub-kod izahı:**
- `r.Cookie("ad")` → serverdə cookie axtarır
- `req.AddCookie(&c)` → klientdə request-ə əlavə edir
- `http.SetCookie(w, &c)` → response-un Set-Cookie header-ini yazır

### 6. CookieJar — avtomatik cookie idarəetməsi
**Nədir:** `net/http/cookiejar` — yaddaşda cookie saxlayan anbar; client
hər sorğudan sonra cookiləri avtomatik yeniləyir.

```go
jar, err := cookiejar.New(nil)
if err != nil {
    panic(err)
}
cookies := []*http.Cookie{
    &http.Cookie{Name: "counter", Value: "1"},
}

url := "http://localhost:8090/cookie"
u, _ := url2.Parse(url)
jar.SetCookies(u, cookies)

client := http.Client{Jar: jar}           // jar client-ə bağlanır

for i := 0; i < 5; i++ {
    _, err := client.Get(url)
    if err != nil {
        panic(err)
    }
    fmt.Println("Client cookie", jar.Cookies(u))
}
// Client cookie [counter=2] ... [counter=6] — avtomatik artır
```
- `jar.SetCookies(u, cookies)` → ilkin cookieləri qeyd edir
- `client.Get` sonrası `Set-Cookie` avtomatik jar-a yazılır — əl ilə
  yeniləmək lazım deyil

### 7. Middleware
**Nədir:** Handler-dən əvvəl/sonra icra olunan, təkrar istifadə edilə bilən
qatlar (auth, logging...). net/http hazır alət vermir — handler sarğısı ilə:

```go
func Middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // next işə düşməzdən əvvəl nəsə et
        next.ServeHTTP(w, r)
        // next bitdikdən sonra nəsə et
    })
}
http.ListenAndServe(":8090", Middleware(Middleware(Handler)))  // zəncir
```

**Basic auth middleware nümunəsi:**
```go
func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        header := r.Header.Get("Authorization")
        if header == "" {
            w.WriteHeader(http.StatusUnauthorized)
            return
        }
        authType := strings.Split(header, " ")
        if len(authType) != 2 || authType[0] != "Basic" {
            w.WriteHeader(http.StatusUnauthorized)
            return
        }
        credentials, err := base64.StdEncoding.DecodeString(authType[1])
        if err != nil {
            w.WriteHeader(http.StatusUnauthorized)
            return
        }
        if string(credentials) == "Open Sesame" {
            next.ServeHTTP(w, r)      // yalnız düzgün parolda keçir
        }
    })
}

func main() {
    targetHandler := MyHandler{}
    panic(http.ListenAndServe(":8090", AuthMiddleware(&targetHandler)))
}
```

**Test (base64 kodlama):**
```bash
>>> auth=$(echo -n "Open Sesame" | base64)
>>> curl :8090 -w "%{http_code}"                        # 401
>>> curl :8090 -w "%{http_code}" -H "Authorization: Basic $auth"
Perfect!!!200
```

**Proqramatik middleware zənciri:**
```go
type Middleware func(http.Handler) http.Handler

func ApplyMiddleware(h http.Handler, middleware ...Middleware) http.Handler {
    for _, next := range middleware {
        h = next(h)     // hər middleware əvvəlkini bükür
    }
    return h
}

func SimpleMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        value := w.Header().Get("simple")
        if value == "" {
            value = "X"
        } else {
            value = value + "X"
        }
        w.Header().Set("simple", value)
        next.ServeHTTP(w, r)
    })
}

func main() {
    h := &MyHandler{}
    http.Handle("/three", ApplyMiddleware(
        h, SimpleMiddleware, SimpleMiddleware, SimpleMiddleware))  // Simple: XXX
    http.Handle("/one", ApplyMiddleware(h, SimpleMiddleware))        // Simple: X
    panic(http.ListenAndServe(":8090", nil))
}
```
- `ApplyMiddleware` → middleware-ləri sağdan sola bükür; cavab header-ində
  icra sayı görünür (`Simple: XXX` vs `Simple: X`)

## Əsas terminlər
- Request/Response — HTTP sorğu/cavab strukturları
- ServerMux — URL pattern → handler yönləndiricisi
- Handler — `ServeHTTP(ResponseWriter, *Request)` implementasiyası
- ResponseWriter — cavabın yazıldığı interfeys
- Cookie — client state fraqmenti
- CookieJar — cookielərin avtomatik saxlanması
- Middleware (arasofta qat) — handler ətrafında bükülmüş emal qatı
- Basic auth — base64 credential-lı `Authorization: Basic` header

## Praktik nəticə
Sadə sorğular üçün `http.Get/Post/PostForm`; tam nəzarət üçün
`http.NewRequest` + `http.Client{Timeout}` + `client.Do`. Server üçün
`http.HandleFunc("/yol", handler)` + `ListenAndServe`. Response body həmişə
Reader-dir — `defer Close` şərtdir. Çoxlu cookie ilə CookieJar istifadə edin.
Middleware = `func(http.Handler) http.Handler` funksiyaları zənciri — auth,
log, rate-limit həmin pattern ilə qurulur.

## Mənbə
Pages: 183-207 (PDF 183-207)
