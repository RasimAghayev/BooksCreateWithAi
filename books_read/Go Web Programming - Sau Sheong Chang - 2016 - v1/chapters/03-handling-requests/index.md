# Chapter 3 — Handling requests

## Bu chapter nədən bəhs edir?

net/http paketinin server tərəfi: Server konfiqurasiyası, HTTPS (öz SSL sertifikatını yaratmaq daxil), handler vs handler function fərqi, HandlerFunc adapteri, handler zəncirlənməsi (chaining — middleware), ServeMux/DefaultServeMux routing qaydaları, üçüncü tərəf router HttpRouter və HTTP/2 aktivləşdirmə.

## Əsas fikirlər

### 1. Niyə framework yox, standard library?
**Cargo cult programming təhlükəsi:** framework konvensiyalarını anlamadan izləyən proqramçı (kod kopyalayan, StackOverflow-dan yapışdıran) tətbiqi framework-a bağlayır — dəyişmək/taşıməq üçün framework-in dərin biliyi lazım olur. Cookie/session NİYƏ var? Çünki HTTP connection-lessDIR — bunu bilmədən framework interfeysi "maqik" görünür. Kitab: standard library (`net/http` + `html/template`) üzərindən anlama → sonra istəsən framework.

**net/http 2 hissəsi:**
- **Client:** Client, Response, Header, Request, Cookie
- **Server:** Server, ServeMux, Handler/HandleFunc, ResponseWriter, Header, Request, Cookie

### 2. Web server yaratma
**Ən sadə server:**
```go
package main

import (
    "net/http"
)

func main() {
    http.ListenAndServe("", nil)
}
```
**Sub-kod izahı:**
- `ListenAndServe(addr, handler)` → addr boşsa bütün interfeyslər port 80; handler `nil` isə **DefaultServeMux** istifadə olunur

**Konfiqurasiyalı server — Server struct:**
```go
server := http.Server{
    Addr:    "127.0.0.1:8080",
    Handler: nil,
}
server.ListenAndServe()
```

**Server struct sahələri:**
```go
type Server struct {
    Addr           string
    Handler        Handler
    ReadTimeout    time.Duration
    WriteTimeout   time.Duration
    MaxHeaderBytes int
    TLSConfig      *tls.Config
    TLSNextProto   map[string]func(*Server, *tls.Conn, Handler)
    ConnState      func(net.Conn, ConnState)
    ErrorLog       *log.Logger
}
```
- ReadTimeout/WriteTimeout → yavaş client müdafiəsi; ErrorLog → server daxili xətalar üçün ayrıca logger

### 3. HTTPS servisi
**Nədir:** HTTP + SSL/TLS qatı. Login olan sayt üçün mütləq; kredit kartı qəbul edənlər üçün PCI DSS tələbi.

**Server tərəfi:**
```go
server := http.Server{
    Addr:    "127.0.0.1:8080",
    Handler: nil,
}
server.ListenAndServeTLS("cert.pem", "key.pem")
```
- `cert.pem` → SSL sertifikat; `key.pem` → server private key
- Production-da sertifikat CA-dən (VeriSign, Thawte, Comodo); development-da özün yarada bilərsən

**Öz sertifikatını yaratmaq (crypto paketləri ilə):**
```go
package main

import (
    "crypto/rand"
    "crypto/rsa"
    "crypto/x509"
    "crypto/x509/pkix"
    "encoding/pem"
    "math/big"
    "net"
    "os"
    "time"
)

func main() {
    max := new(big.Int).Lsh(big.NewInt(1), 128)
    serialNumber, _ := rand.Int(rand.Reader, max)
    subject := pkix.Name{
        Organization:       []string{"Manning Publications Co."},
        CommonName:         "Go Web Programming",
    }
    template := x509.Certificate{
        SerialNumber: serialNumber,
        Subject:      subject,
        NotBefore:    time.Now(),
        NotAfter:     time.Now().Add(365 * 24 * time.Hour),
        KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
        ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
        IPAddresses:  []net.IP{net.ParseIP("127.0.0.1")},
    }
    pk, _ := rsa.GenerateKey(rand.Reader, 2048)
    derBytes, _ := x509.CreateCertificate(rand.Reader, &template, &template, &pk.PublicKey, pk)
    certOut, _ := os.Create("cert.pem")
    pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
    certOut.Close()
    keyOut, _ := os.Create("key.pem")
    pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(pk)})
    keyOut.Close()
}
```

**Sub-kod izahı:**
- `x509.Certificate{...}` → sertifikat konfiqurasiyası: serial nömrə (random 128-bit), Subject (distinguished name), müddət (1 il), KeyUsage/ExtKeyUsage (server authentication), IP
- `rsa.GenerateKey(rand.Reader, 2048)` → 2048-bit RSA private key; `.PublicKey` ilə public hissə
- `x509.CreateCertificate(...)` → DER formatında byte slice
- `pem.Encode` → DER-i PEM faylına yazır (Base64 + `-----BEGIN CERTIFICATE-----` sərhədləri)
- CA imzalı sertifikatda fayl = server sertifikatı + CA sertifikatı (birləşdirilmiş)

**SSL/TLS/HTTPS anlayışları:**
- SSL (Netscape) → IETF tərəfindən TLS adlandırıldı
- X.509 sertifikatı: ASN.1 formatında, public key + məlumat saxlayır; CA imzası autentikliyi təsdiqləyir
- Növbəti mübadilə: client sertifikatı yoxlayır → random **simmetrik açar** yaradır → public key ilə şifrəlir → o simmetrik açar data şifrələmə üçün istifadə olunur
- PEM = Base64-kodlu DER, BEGIN/END sərhədləri ilə

### 4. Handler nədir?
**Tərif:** Go-da handler — `ServeHTTP(w http.ResponseWriter, r *http.Request)` metoduna sahib **interfeys**dir. Bu imzalı metodu olan İSTƏRİLƏN tip handler-dir.

**Kitabdan kod nümunəsi (tək handler):**
```go
package main

import (
    "fmt"
    "net/http"
)

type MyHandler struct{}

func (h *MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello World!")
}

func main() {
    handler := MyHandler{}
    server := http.Server{
        Addr:    "127.0.0.1:8080",
        Handler: &handler,
    }
    server.ListenAndServe()
}
```

**Sub-kod izahı:**
- Handler qeyd ediləndə mux YOXDUR → **bütün URL-lər** (`/`, `/anything/at/all`) eyni cavabı alır — URL matching yoxdur
- Server-a Handler verildikdə DefaultServeMux sıradan çıxır

### 5. Çox handler — DefaultServeMux ilə
**Kitabdan kod nümunəsi:**
```go
type HelloHandler struct{}

func (h *HelloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello!")
}

type WorldHandler struct{}

func (h *WorldHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "World!")
}

func main() {
    hello := HelloHandler{}
    world := WorldHandler{}
    server := http.Server{
        Addr: "127.0.0.1:8080",
    }
    http.Handle("/hello", &hello)
    http.Handle("/world", &world)
    server.ListenAndServe()
}
```

**Sub-kod izahı:**
- Server-da Handler qeyd edilmir → DefaultServeMux işə düşür
- `http.Handle(url, handler)` → qısayol: əslində `DefaultServeMux.Handle` çağırılır
- `/hello` → HelloHandler, `/world` → WorldHandler

### 6. Handler function və HandlerFunc adapteri
**Fərq:** handler funksiyası — ServeHTTP imzasıyla adi funksiyadır; HandlerFun**c** isə **funksiya tipi**dir ki, adi funksiyanı Handler-ə çevirir:

**Kitabdan kod nümunəsi:**
```go
func hello(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello!")
}

func world(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "World!")
}

func main() {
    server := http.Server{
        Addr: "127.0.0.1:8080",
    }
    http.HandleFunc("/hello", hello)
    http.HandleFunc("/world", world)
    server.ListenAndServe()
}
```

**Net/http source daxili:**
```go
func HandleFunc(pattern string, handler func(ResponseWriter, *Request)) {
    DefaultServeMux.HandleFunc(pattern, handler)
}

func (mux *ServeMux) HandleFunc(pattern string, handler func(ResponseWriter, *Request)) {
    mux.Handle(pattern, HandlerFunc(handler))
}
```
- `HandlerFunc(handler)` → funksiyanı Handler interfeysinə uyğunlaşdırır
- **HandlerFunc tipi özü ServeHTTP metoduna malikdir** — içində saxladığı funksiyanı çağırır (funksiya tipi + metod = adapter pattern)

**Nə vaxt handler, nə vaxt handler funksiya?** Dizayn məsələsi: tipin həm handler, həm başqa şey olmasını istəyirsənsə ServeHTTP metodunu əlavə et — modulyarlıq artır. Sadə hallarda handler funksiyası təmizdir.

### 7. Chaining (zəncirləmə) — middleware
**Problem:** logging, security, error handling — **cross-cutting concern**-lər: hər handler-ə əlavə etmək intruzivdir, kodu dublikasiya edir.

**Həll:** handler funksiyasını funksiyadan keçirib funksiya qaytaran wrapper:

**Handler funksiyaların zənciri:**
```go
func hello(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello!")
}

func log(h http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        name := runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name()
        fmt.Println("Handler function called - " + name)
        h(w, r)
    }
}

func main() {
    server := http.Server{Addr: "127.0.0.1:8080"}
    http.HandleFunc("/hello", log(hello))
    server.ListenAndServe()
}
```

**Sub-kod izahı:**
- `log(h http.HandlerFunc) http.HandlerFunc` → funksiya qəbul edib funksiya qaytarır
- `runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name()` → reflection ilə funksiya adını tapır (`main.hello`)
- `h(w, r)` → öz funksiyasını çağırır — sıra: log → hello
- Çoxluq zəncir: `http.HandleFunc("/hello", protect(log(hello)))` — pipeline processing (Lego kərpicləri)

**Handler-ların zənciri (Handler interfeysi ilə):**
```go
func log(h http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        fmt.Printf("Handler called - %T\n", h)
        h.ServeHTTP(w, r)
    })
}

func protect(h http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // ... authorization yoxlaması
        h.ServeHTTP(w, r)
    })
}

func main() {
    server := http.Server{Addr: "127.0.0.1:8080"}
    hello := HelloHandler{}
    http.Handle("/hello", protect(log(hello)))
    server.ListenAndServe()
}
```
- Fərq: `h.ServeHTTP(w, r)` çağırışında anonim funksiyanı `http.HandlerFunc(...)` ilə Handler-ə çevirir
- Bu idiom bütün Go web framework-lərinin middleware əsasıdır

### 8. ServeMux routing qaydaları
**Nədir:** ServeMux — struct: URL → handler map saxlayır; özü də handlerdir (ServeHTTP metodu var — DefaultServeMux'u Server-a handler kimi vermək məhz buna görə işləyir!).

**Pattern matching qaydaları:**
- `/hello` (slash YOXDUR sonuca) → **yalnız dəqiq match**: `/hello/there` → `/hello`-ya DÜŞMÜR
- `/hello/` (slash VAR) → **prefiks match**: `/hello/there` → `/hello/`-a düşür
- `/` (root) → heç nə tutulmayan URL-lərin hamısı bura düşür (fallback)
- Əgər `/random` sorğulanıbsa və qeydə alınmayıbsa → root-un handler-i icra olunur

**Principle of Least Surprise qutusu:** "ən az təəccüb prinsipi" — dizayn ən gözlənilən davranışı göstərməlidir. Qapı yanında düymə qapıya dair olmalıdır, koridor işığına yox.

### 9. HttpRouter — üçüncü tərəf multiplexer
**ServeMux əsas çatışmamazlığı:** URL-də **dəyişən** yoxdur. `/thread/123`-ü emal etmək üçün handler path-i özü parse etməlidir; `/thread/123/post/456` praktiki olmur.

**Kitabdan kod nümunəsi:**
```go
package main

import (
    "fmt"
    "github.com/julienschmidt/httprouter"
    "net/http"
)

func hello(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    fmt.Fprintf(w, "hello, %s!\n", p.ByName("name"))
}

func main() {
    mux := httprouter.New()
    mux.GET("/hello/:name", hello)
    server := http.Server{
        Addr:    "127.0.0.1:8080",
        Handler: mux,
    }
    server.ListenAndServe()
}
```

**Sub-kod izahı:**
- `mux.GET("/hello/:name", hello)` → **metod-əsaslı** qeydiyyat (GET/POST/...); `:name` → named parameter
- Handler imzasına 3-cü parametr gəlir: `p httprouter.Params`
- `p.ByName("name")` → URL-dəki dəyəri oxu (`/hello/sausheong` → `sausheong`)
- Server-a `Handler: mux` → DefaultServeMux əvəzinə bu mux

**Quraşdırma:**
```bash
$ go get github.com/julienschmidt/httprouter
```
`go get` → repodan endirib `$GOPATH/src`-ə qoyur → `go build` artıq tapır.

Digər alternativlər: Gorilla Toolkit (mux, pat) — www.gorillatoolkit.org.

### 10. HTTP/2 aktivləşdirmə
- Go 1.6+: HTTPS + `ListenAndServeTLS` → HTTP/2 avtomatik
- Əvvəlki versiyalar: `go get "golang.org/x/net/http2"` + `http2.ConfigureServer(&server, &http2.Server{})` — bundan sonra `ListenAndServeTLS` ilə HTTP/2 işləyir.

**Yoxlama (cURL):**
```bash
curl -I --http2 --insecure https://localhost:8080/
```
- `--http2` → HTTP/2 sorğusu (cURL 7.43+, nghttp2 library ilə linkli olmalıdır)
- `--insecure` → self-signed sertifikatı qəbul et
- Nəticə: `HTTP/2.0 200` — uğur

## Handler familyası xülasəsi
| Konsept | Nədir | İmza/Example |
|---|---|---|
| Handler | Interfeys | `ServeHTTP(w, r)` metodu |
| Handler function | Handler kimi davranan funksiya | `func(w, r)` |
| HandlerFunc | Adapter tipi | `http.HandlerFunc(f)` → Handler |
| ServeMux | Multiplexer (özü də handler) | `NewServeMux()` / `DefaultServeMux` |
| Middleware | `func(http.Handler) http.Handler` | `protect(log(hello))` |

## Əsas terminlər
- Cargo Cult Programming (kor-koranə kod kopyalama)
- Cross-cutting Concern (kəsən qayğı — logging/security/error)
- Chaining / Middleware / Pipeline Processing (zəncirləmə)
- Handler / Handler Function / HandlerFunc
- ServeMux / DefaultServeMux
- Exact Match / Prefix Match (dəqiq / prefiks uyğunluq)
- Named Parameter (`:name` — adlandırılmış parametr)
- HTTPS / SSL / TLS / X.509 / PEM / DER / ASN.1
- Certificate Authority (CA)
- RSA Private Key
- PCI DSS (kredit kartı təhlükəsizlik standartı)
- Principle of Least Surprise

## Praktik nəticə
- Production server-də Server struct-un ReadTimeout/WriteTimeout qoy — yavaş client-lər resurs yandırır.
- Self-signed sertifikat yalnız dev üçün; real sayt üçün CA (bu gün LetsEncrypt) istifadə et.
- Middleware üçün `func(http.Handler) http.Handler` idiomunu mənimsə — bütün Go web framework-ləri eyni pattern-i işlədir.
- ServeMux pattern-ləri: trailing slash = prefiks match, slashsiz = dəqiq match — `/hello` və `/hello/` FƏRLİ davranır.
- URL-də dəyişən lazımdırsa (REST path) HttpRouter/Gorilla mux seç, yoxsa handler-də path parse etməli olursan.
- `go get` ilə üçüncü tərəf paketi gətir — $GOPATH-ə yerləşir, build avtomatik tapır.

## Mənbə
Pages: 68-89 (PDF), book pages 47-68
