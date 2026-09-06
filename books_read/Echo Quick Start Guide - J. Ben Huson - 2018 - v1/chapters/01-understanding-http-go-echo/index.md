# Chapter 1 — Understanding HTTP, Go, and Echo

## Bu chapter nədən bəhs edir?

Bu chapter veb-tətbiqetmənin **protokol əsaslarına** qayıdış edir: HTTP request/response mesaj strukturu (RFC-7230), Go-nun `net/http` standart kitabxanasının təqdim etdiyi primitivlər (Handler interfeysi, web server, goroutine modeli), standart kitabxananın minimalizmi nədəni ilə framework ehtiyacı və Go + Echo mühitinin qurulması izah olunur.

---

## Əsas fikirlər

### 1. HTTP basics — request/response strukturu

HTTP — **stateless (dəyişməz hal saxlamayan)**, mətn əsaslı protokoldur: hər request bir response doğurur. Veb-tətbiqetmənin girişi HTTP request, çıxışı HTTP response-dur.

#### HTTP Request (RFC-7230 Section 3)

Məcburi elementlər (start-line / request-line hissəsi):

| Element | RFC | İzah |
|---------|-----|------|
| **Request method** | RFC-7231 4.3; RFC-5789 2 | Klientin resurs üzərində icra etmək istədiyi əməl (finite fel siyahısı — GET, POST, PUT...). Go-da `net/http` konstantları |
| **Request target** | RFC-7230 5.3 | URI-dən törəmiş resurs identifikatoru: host, path, query, fragment hissələri. Go-da `net/url` paketi parse edir |
| **HTTP version** | RFC-7230 2.6 | Klientin dəstəklədiyi protokol versiyası (məs., HTTP/1.1) — Go `http.Request`-ə doldurulur |

Opsional elementlər:
- **Header sahələri (RFC-7230 3.2):** `Açar: Dəyər` formatında metadata (Content-Length, auth credentials və s.) — Go-da `http.Request.Header` map-i. Qeydiyyatdan keçmiş header siyahısı: http://www.iana.org/assignments/message-headers
- **Message body (RFC-7230 3.3):** Request-in payload-ı — body mövcudluğunu `Content-Length` və ya `Transfer-Encoding` header-ləri bildirir. Go-da `Request.Body` — `io.ReadCloser` interfeysi (`Read` + `Close` metodları)

Tam request nümunəsi:

```http
POST /file.txt HTTP/1.1
Host: localhost
Accept: */*
Content-Length: 8
Content-Type: application/x-www-form-urlencoded

hi=there
```

#### HTTP Response

Struktur request-ə bənzəyir; status-line ibarətdir:
- **HTTP-Version** — serverin cavab verdiyi protokol versiyası
- **Status-Code (RFC-7231 6)** — 3 rəqəmli kod, 5 sinif: informational (1xx), successful (2xx), redirection (3xx), client error (4xx), server error (5xx)
- **Reason-Phrase** — klient tərəfindən nəzərə alınmaması lazım olan mətn izahı

Go-da cavab adətən **`http.ResponseWriter` interfeysi** üzərindən idarə olunur:
- `Header()` — header map-ini doldurmaq
- `WriteHeader(code)` — status-line yazmaq
- `Write([]byte)` — body göndərmək

> Framework qiymətləndirərkən protokolun sadəliyini yadda saxlayın: mürəkkəb sistemlər belə sadə mətn protokolu üzərində qurulur.

---

### 2. Go HTTP handlers

`net/http`-da **`http.Handler` interfeysi** — server-in sizin kodunuzu çağıracağı vahid giriş nöqtəsi:

```go
type Handler interface {
    ServeHTTP(ResponseWriter, *Request)
}
```

Minimal nümunə (`SimpleHandler.go`):

```go
package main

import "net/http"

func main() {
    http.Handle("/", new(myHandler))
    http.ListenAndServe(":8080", nil)
}

type myHandler struct{}

func (mh *myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(200)
    w.Write([]byte("hello!"))
}
```

İşə salma: `go run SimpleHandler.go` + `curl localhost:8080/`.

**Bu sadə nümunədə yaranan suallar** (framework-lərin mövcudluq səbəbi):
1. `/resource/$ID` kimi URL path dəyişənlərini necə istifadə edim?
2. Handler içində panic baş versə nə olur?
3. Response body-ni hər dəfə `[]byte`-a özüm encodeliyirəm?
4. Eyni prosesləri birdən çox handler-də necə təkrar istifadə edim?

---

### 3. Go HTTP web server — daxili model

Go `net/http` serveri **bare bones (çıx qab)** serverdir:

**Request emal axını:**
1. TCP listen + accept həyat dövrü
2. Hər yeni request üçün **yeni goroutine** başladılır
3. Developer-in təyin etdiyi handler həmin goroutine içində icra olunur
4. Application səviyyəli panic-lər izlənir, handler-in yaratdığı response klientə göndərilir

**Goroutine üstünlüyü:** Traditional thread/subprocess-lərdən qat-qat ucuz başlatma və işlətmə — minlərlə goroutine eyni anda işləyə bilər. TLS dəstəyi daxilidir (istəyə bağlı).

---

### 4. Framework nəyə lazımdır?

`net/http`-nın minimalizmi **qüsursuzluq deyil, design qərarıdır** — developer-ə xüsusi use case-lərə uyğun həll qurma azadlığı verir. Amma bu, iş yaratır:

| Çatışmazlıq | Nəticə |
|-------------|--------|
| Middleware / request pipeline yoxdur | Hər layihə öz pipeline-ini yazmalı |
| URL routing funksionallığı zəif | `/resource/$ID` kimi pattern-lər üçün əl ilə regex |
| Render/binding helper-ləri yoxdur | Encoding və request-in dəyişənlərə bağlanması əl ilə |
| Request context yox idi (son vaxtlara qədər) | Handler-lər arası state ötürülməsi problemli |

**Həll:** Echo — performance-focused (performansa yönəlmiş), open source Go veb framework-u ki, `net/http`-ın buraxdığı boşluqları doldurur.

---

### 5. Mühitin qurulması

#### Go quraşdırması (kitabda Go 1.9.4)

Linux/macOS:

```bash
# Linux
curl -O https://dl.google.com/go/go1.10.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.10.linux-amd64.tar.gz

# macOS (Homebrew)
brew install go

# Environment dəyişənləri (bash)
mkdir -p ~/go
echo "export GOPATH=$HOME/go" >> ~/.bashrc
echo "export PATH=$PATH:$HOME/go/bin:/usr/local/go/bin" >> ~/.bashrc
source ~/.bashrc
```

Windows: MSI installer (`go1.10.windows-amd64.msi`) — `c:\Go\` qovluğuna, environment dəyişənləri avtomatik.

#### Echo quraşdırması (v3.3.5 — kitabın versiyası)

```bash
go get github.com/labstack/echo
cd $GOPATH/src/github.com/labstack/echo
git checkout tags/3.3.5   # versiya pinləmə — go get versiyalı download dəstəkləmir
cd -
```

> Sonrakı chapter-lərdə **dep** tool istifadə olunacaq — layihə-bazalı versiya pinləməni avtomatlaşdırır.

#### Yoxlama proqramı (`environment_setup.go`)

```go
package main

import (
    "net/http"
    "github.com/labstack/echo"
)

func main() {
    // create a new echo instance
    e := echo.New()
    // Route / to handler function
    e.GET("/", handler)
    // start the server, and log if it fails
    e.Logger.Fatal(e.Start(":8080"))
}

// handler - Simple handler to make sure environment is setup
func handler(c echo.Context) error {
    // return the string "Hello World" as the response body
    // with an http.StatusOK (200) status
    return c.String(http.StatusOK, "Hello World")
}
```

`go run environment_setup.go` — Echo banner çap olunsa (v3.3.5 + `http server started on [::]:8080`), mühit hazırdır. Nümunə kod: `git clone https://github.com/PacktPublishing/Echo-Essentials` (dependency-lər `./chapter1/vendor` içində).

---

## Əsas terminlər

| Termin | İzah |
|--------|------|
| Stateless protocol | Hər request müstəqildir — server hallar arasında state saxlamır |
| Start-line / Request-line / Status-line | HTTP mesajının ilk sətri (request-də method+target+version, response-da version+code+phrase) |
| Request target | URI-dən törəmiş resurs identifikatoru (host/path/query/fragment) |
| Header fields | `Açar: Dəyər` metadata cütləri (Content-Length, Content-Type, auth...) |
| Message body | Mesajın payload-ı — `io.ReadCloser` kimi modellənir |
| `http.Handler` interfeysi | `ServeHTTP(ResponseWriter, *Request)` — server-in kodunuzu çağırdığı vahid nöqtə |
| `http.ResponseWriter` | Cavab yazma interfeysi: `Header()` / `WriteHeader()` / `Write()` |
| Goroutine | Go-nun yüngül thread-i — hər request üçün yeni goroutine |
| Bare bones server | Yalnız əsas funksionallıq daşıyan minimalist server |
| Echo | Performance-focused, minimalist Go veb framework-u (labstack/echo) |
| `GOPATH` | Go workspace dəyişəni — paketlərin yükləndiyi qovluq |
| `go get` | Paket yükləmə əmri (versiyalı deyil) |
| Vendor directory | Layihə daxilində dependency kodunun saxlandığı xüsusi qovluq |
| dep | Go dependency management tool — versiya pinləmə |

---

## Praktik nəticə

1. **Protokolu anlayın:** HTTP sadəcə mətn mesajları mübadiləsidir — request (start-line + header + body) → response (status-line + header + body). Bütün framework maşınları bu sadəliyin üzərindədir.
2. **`net/http` primitivləri əsasdır:** `Handler` interfeysi, `ResponseWriter` və request-başına-goroutine modeli bütün Go veb framework-lərinin (Echo daxil) altında işləyir.
3. **Framework ehtiyacı realdır:** URL dəyişənləri, panic idarəsi, render helper-ləri, middleware — bunları sıfırdan yazmaq əvəzinə Echo təqdim edir.
4. **Versiyanı pinləyin:** `go get` + `git checkout tags/...` əl yolu, `dep` isə avtomatik yol — kitab boyu dep istifadə olunacaq.
5. **İlk Echo tətbiqi 6 sətirlik:** `echo.New()` → `e.GET` → `e.Logger.Fatal(e.Start(":8080"))` — handler-da `c.String(http.StatusOK, "Hello World")`.

---

## Mənbə

- Kitab: *Echo Quick Start Guide* — J. Ben Huson, Packt Publishing, 2018 (ISBN 9781789340849)
- Chapter 1: "Understanding HTTP, Go, and Echo", book səh. 1–22
- PDF səhifələri: 30–47
- Kod: https://github.com/PacktPublishing/Echo-Essentials/tree/master/chapter1
- Video: https://goo.gl/QjupHX
- RFC-lər: RFC-7230 (HTTP/1.1 message syntax), RFC-7231 (methods/status codes), RFC-5789 (PATCH)
