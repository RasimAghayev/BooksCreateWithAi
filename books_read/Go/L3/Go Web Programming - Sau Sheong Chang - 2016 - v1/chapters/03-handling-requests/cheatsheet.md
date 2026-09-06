# Chapter 3 — Handling requests Cheatsheet

## `http.ListenAndServe("", nil)`

**Nə edir:** Boş `Addr` ilə serveri işə salır — hamı network interface-lər port 80. `nil` handler → `DefaultServeMux` istifadə olunur.

**Parametrlər:**
- `""` (boş) → hamı interface-lər, default port 80
- `nil` → `DefaultServeMux` multiplexer

**Mənbə:** Chapter 3, page 51

---

## `type Server struct {...}`

**Nə edir:** HTTP server konfiqurasiyasını təyin edən struct.

**Field-lər:**
```go
type Server struct {
    Addr string
    Handler Handler
    ReadTimeout time.Duration
    WriteTimeout time.Duration
    MaxHeaderBytes int
    TLSConfig *tls.Config
    ConnState func(net.Conn, ConnState)
    ErrorLog *log.Logger
}
```

**Mənbə:** Chapter 3, page 52

---

## `server.ListenAndServeTLS("cert.pem", "key.pem")`

**Nə edir:** HTTPS serveri işə salır — cert.pem (sertifikat), key.pem (_private açar).

**Nəyə lazımdır:** Giriş (login) səhifələri üçün təhlükəsiz əlaqə, PCI DSS compliance.

**Mənbə:** Chapter 3, page 53

---

## `func (h *MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request)`

**Nə edir:** Handler interfeysünü (ServeHTTP metodu) implement edən struct. `http.Server`-də `Handler`-ə ötürülür.

```go
type MyHandler struct{}
func (h *MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello World!")
}
// server := http.Server{Handler: &handler}
```

**Mənbə:** Chapter 3, page 56

---

## `http.HandleFunc("/hello", hello)` / `http.Handle("/hello", &hello)`

**Nə edir:** Handler funksiyası/Handler-ı URL-ə bağlayır (`DefaultServeMux`-a qeyd olunur).

**Fərqlər:**
- `HandleFunc(pattern, handler func)` → daxili `HandlerFunc(handler)` ilə handler-ə çevirir
- `Handle(pattern, Handler)` → özündə `Handler` interfeysi tələb edir

**HandlerFunc adaptasiyası:**
```go
func HandleFunc(pattern string, handler func(ResponseWriter, *Request)) {
    DefaultServeMux.HandleFunc(pattern, handler)
}
func (mux *ServeMux) HandleFunc(pattern string, handler func(ResponseWriter, *Request)) {
    mux.Handle(pattern, HandlerFunc(handler))
}
```

**Mənbə:** Chapter 3, page 59

---

## Handler chaining: `log(hello)` / `protect(log(hello))`

**Nə edir:** Handler funksiyaları zəncirləyərək cross-cutting concerns (log, security) əlavə edir.

```go
func log(h http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        name := runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name()
        fmt.Println("Handler function called - " + name)
        h(w, r)
    }
}
http.HandleFunc("/hello", protect(log(hello)))
```

**Pipeline:** `protect` (f3) → `log` (f2) → `hello` (f1) — sorğu tərsinə gedir, cavab axtı-ma ilə.

**Mənbə:** Chapter 3, page 61

---

## `mux := httprouter.New()` + `mux.GET("/hello/:name", hello)`

**Nə edir:** HttpRouter ilə dinamik URL routing. `:name` — named parameter.

```go
mux := httprouter.New()
mux.GET("/hello/:name", hello)
// handler: func hello(w http.ResponseWriter, r *http.Request, p httprouter.Params)
fmt.Fprintf(w, "hello, %s!\n", p.ByName("name"))
```

**Quraşdırma:** `go get github.com/julienschmidt/httprouter`

**Mənbə:** Chapter 3, page 67

---

## `ServeMux` pattern matching qanunu

**Nə edir:** `/` ilə bitən pattern — prefix match; `/` ilə bitməyən pattern — tam eşleşmə.

**Nümunə:**
- `/hello` qeyd etsən, `/hello/there` eşleşməz → fallback `/`-a düşür
- `/hello/` qeyd etsən, `/hello/there` eşleşir (`helloHandler` çağırılır)

**Mənbə:** Chapter 3, page 66

---

## `go get <package>` / `go build`

**Nə edir:** Üçüncü şəxs kitabxananı `$GOPATH/src`-ə endirir.

**Sıra:**
1. `go get github.com/julienschmidt/httprouter` — install
2. `go build` — import ilə kodu compile et

**Mənbə:** Chapter 3, page 68

---

## `http2.ConfigureServer(&server, &http2.Server{})`

**Nə edir:** Serveri HTTP/2-ə güncər. Go 1.6+ HTTPS ilə default HTTP/2.

**Yoxlama:**
```bash
curl -I --http2 --insecure https://localhost:8080/
# HTTP/2.0 200
```

**Mənbə:** Chapter 3, page 69
