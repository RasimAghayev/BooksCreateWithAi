# Chapter 1 — Go and web applications

## Bu chapter nədən bəhs edir?

Web tətbiqlərinin (web applications) tərifi, Go-nun web proqramlaşdırması üçün üstünlükləri, HTTP protokolunun tam əsasları (request/response strukturu, metodlar, header-lər, status kodları, URI), web tətbiqinin 2 əsas komponenti (handler + template engine) və ilk Go web tətbiqinin yazılması.

## Əsas fikirlər

### 1. Go niyə web üçün uyğundur?
**Nədir:** Kiçik və böyük miqyaslı web tətbiqləri üçün Go-nun 4 əsas üstünlüyü.

**Necə işləyir:**
1. **Scalable (miqyaslanan):** 
   - Vertical scaling (bir maşınada CPU artırma) — Go-nun concurrency dəstəyi: tək OS thread üzərində yüz minlərlə goroutine
   - Horizontal scaling (maşın sayını artırma) — Go **statik binary**-yə kompilyasiya olunur, dynamic dependency yoxdur → hər hansı sistemə bir fayl kopyalayıb işə salmaq olar
2. **Modular:** interfeys mexanizmi davranışı təsvir edir → yeni kod köhnə funksiyalarla işləyir; empty interface hər tipi qəbul edir; funksiya tipləri + closure-lar funksional modulliqliq verir; mikro-servis arxitekturası üçün idealdır
3. **Maintainable:** təmiz sintaksis, `gofmt` vahid format, `godoc` kod+şərhdən sənəd, `gotest` daxili testing
4. **High-performance:** native code kompilyasiyası (interpreted dillərdən sürətli), goroutine-lər sayəsində eyni anda çox request

**İstifadə edən şirkətlər (kitab dövrü):** Dropbox, SendGrid, Square, Hailo, BBC, The New York Times.

### 2. Web application nədir?
**Tərif (kitabın dəqiq çərçivəsi):** Proqram bu 2 meyara cavab verirsə web app-dir:
1. Çağıran client-ə HTML qaytarır (client render edib istifadəçiyə göstərir)
2. Data HTTP ilə client-ə ötürülür

**Web server vs web app:** Web server (Apache) sadəcə qovluqdakı (docroot) faylları qaytarır; web app isə request-i **emal edir**, proqramlaşdırılmış əməliyyatları icra edir. Web server = fayl qaytaran ixtisaslaşmış web app.

**Web service:** HTML render etməyən, başqa proqrama data qaytaran tətbiq (Ch 7-də ətraflı).

### 3. HTTP — əsas protokol
**Tərif:** "HTTP is a stateless, text-based, request-response protocol that uses the client-server computing model."

**Komponentlərin izahı:**
- **Stateless (vəziyyətsiz):** hər request müstəqildir — server əvvəlki request-ləri xatırlamır (FTP/Telnet kimi persistent kanal yoxdur; HTTP 1.1 performans üçün connection-i saxlayır, amma state yoxdur)
- **Text-based:** plain text ötürülür → telnet ilə debug mümkün
- **Request-response:** client sorğu göndərir, server cavab verir
- **Client-server:** client (user-agent, adətən brauzer) həmişə söhbəti başladır

**Tarixi:** HTTP 0.9 (yalnız GET, yalnız HTML) → 1.0 (1996, POST/HEAD) → 1.1 (1999) → HTTP/2 (draft). Kitab HTTP 1.1-ə fokuslanır.

**Web app-in təkamülü:** CGI (1993, NCSA — xarici proqram env variable ilə input, stdout ilə output) → SSI (HTML-də directive: `<!--#include file="navbar.shtml" -->`) → PHP/ASP/JSP/ColdFusion → Mustache/ERB/Velocity template engine-lər.

### 4. HTTP Request strukturu
```
1. Request-line
2. Zero or more request headers
3. An empty line
4. The message body (optional)
```

**Nümunə:**
```
GET /Protocols/rfc2616/rfc2616.html HTTP/1.1
Host: www.w3.org
User-Agent: Mozilla/5.0
(boş sətir)
```

**Request-line:** `METOD URI VERSIYA` — məs. `GET /path HTTP/1.1`.

### 5. Request metodları (tam siyahı)
| Metod | Nə edir |
|---|---|
| GET | Resursu qaytarır |
| HEAD | GET kimidir, amma body YOX — yalnız header-lər lazımda |
| POST | Body-dəki datanı URI-dəki resursa ötürür; server nə edəcəyini özü müəyyən edir |
| PUT | Body-dəki data URI-dəki resurs OLMALIDIR; varsa əvəz et, yoxsa yarat |
| DELETE | URI-dəki resursu sil |
| TRACE | Request-i geri qaytarır — aradakı serverlərin nə etdiyini görmək üçün |
| OPTIONS | Serverin dəstəklədiyi metodların siyahısını qaytarır |
| CONNECT | Şəbəkə bağlantısı qur (əsasən SSL tunneling/HTTPS üçün) |
| PATCH | Body-dəki data resursu **dəyişdirir** (modify) |

**Vacib:** HTTP 1.1-ə görə GET və HEAD implement olunmağı MÜTLƏQDIR; qalanları (hətta POST) optional.

**Safe (təhlükəsiz) metodlar:** serverin vəziyyətini dəyişməyen — GET, HEAD, OPTIONS, TRACE.

**Idempotent (eyniləşdirilə bilən) metodlar:** eyni data ilə 2-ci çağırış vəziyyəti dəyişməyən:
- Safe metodlar → avtomatik idempotent
- PUT, DELETE → idempotent amma safe deyil (PUT 2-ci dəfə eyni nəticəni verir; DELETE 2-ci dəfə error verə bilər, amma state dəyişmir)
- **POST → nə safe, nə idempotent** (hər çağırış yeni state yarada bilər — web service-lərdə kritik fərq, Ch 7-də qayıdılır)

**Brauzer dəstəyi:** HTML form yalnız `get` və `post` qəbul edir; PUT/DELETE üçün **XMLHttpRequest (XHR)** API istifadə olunur (XML ilə məhdud deyil — JSON/text də olur).

### 6. Request header-ləri
**Format:** colon-ayrılmış ad-dəyər cütləri, CRLF ilə bitən. HTTP 1.1-də **yalnız Host mütləqdir**; body varsa Content-Length yaxud Transfer-Encoding də lazımdır.

**Kitabdan cədvəl (əsas header-lər):**
| Header | Təsvir |
|---|---|
| `Accept` | Client-in qəbul etdiyi content type-lar: `Accept: text/html` |
| `Accept-Charset` | Charset: `utf-8` |
| `Authorization` | Basic Authentication credential-ləri |
| `Cookie` | Serverin əvvəl qoyduğu cookie-lər: `Cookie: a=hello; b=world` |
| `Content-Length` | Request body-nin uzunluğu (oktet) |
| `Content-Type` | Body type: default `x-www-form-urlencoded`; fayl yükləmə üçün `multipart/form-data` |
| `Host` | Server adı + port (yoxdursa 80) |
| `Referrer` | Əvvəlki səhifənin ünvanı |
| `User-Agent` | Client təsviri |

### 7. HTTP Response strukturu
```
1. Status line (status code + reason phrase)
2. Zero or more response headers
3. An empty line
4. The message body (optional)
```

**Nümunə:**
```
200 OK
Date: Sat, 22 Nov 2014 12:58:58 GMT
Server: Apache/2
Content-Length: 33115
Content-Type: text/html; charset=iso-8859-1
<!DOCTYPE html ...>...
```

**Status kod sinifləri (1-ci rəqəmə görə):**
| Sinif | Məna |
|---|---|
| 1XX | Informational — request alındı, emal olunur |
| 2XX | Success — əmrə uğurla icra (standart: **200 OK**) |
| 3XX | Redirection — client əlavə addım atmalı (əsasən URL redirect) |
| 4XX | Client Error — request-də problem (məşhur: **404 Not Found**) |
| 5XX | Server Error — server tərəfində problem (**500 Internal Server Error**) |

**Response header-ləri (əsas):**
| Header | Təsvir |
|---|---|
| `Allow` | Serverin dəstəklədiyi metodlar |
| `Content-Length` / `Content-Type` | Body uzunluğu / tipi |
| `Date` | Cari vaxt (GMT) |
| `Location` | Redirect üçün növbəti URL |
| `Server` | Serverin domain adı |
| `Set-Cookie` | Client-də cookie qoyur (bir response-də bir neçə ola bilər) |
| `WWW-Authenticate` | 401 ilə birlikdə — hansı autentifikasiya tələb olunur |

### 8. URI strukturu
**Ümumi forma:** `<scheme>:<hierarchical-part>[?<query>][#<fragment>]`

**Sözbəsöz təhlil nümunəsi:**
```
http://sausheong:password@www.example.com/docs/file?name=sausheong&location=singapore#summary
│     │                 │                  │            │                                     │
scheme user:password@   hierarxik hissə     path         query (? ilə, & ilə ayrılır)          fragment (# ilə)
```
- URN (ad) + URL (yer) = URI (umbrrella); kitabda URL/URI sinonim istifadə olunur
- Yalnız **scheme + hierarchical part mütləqdir**
- Fragment client tərəfində emal olunur — brauzer adətən serverə göndərməzdən əvvəl kəsir (JS/HTTP client ilə göndərmək mümkündür)
- URL-də space və xüsusi simvollar qadağandır → **URL encoding (percent encoding)**: space → `%20` (ASCII 32 = hex 20)

### 9. HTTP/2 qısa icmal
- SPDY/2 (Google) əsasında; **binary protokol** (HTTP 1.x text əsaslı idi) → effektiv parse, kompaktness; amma telnet ilə debug mümkün deyil
- **Fully multiplexed** — bir connection üzərində eyni anda çox request/response
- Header compression, server push
- Semantika (metodlar, status kodları) DƏYİŞMƏYİB — dəyişsəydi web qırılardı
- Go 1.6+: HTTPS işlədilsə HTTP/2 avtomatik aktivdir; əvvəlki versiyalar üçün `golang.org/x/net/http2` paketi

### 10. Web app-in 2 hissəsi
**Web app = Handler + Template Engine.**

**Handler:** HTTP request-i qəbul və emal edir; template engine-i çağırıb HTML generasiya edir; HTML-i HTTP response-a yığıb client-ə qaytarır.

**MVC pattern kontekstində:** handler = controller + model. İdeal MVC-də controller nazikdir (yalnız routing + HTTP unpack/pack), model isə kök (məntiq + data). Service object-lər model-i şişməkdən xilas edir (amma MVC-nin strict hissəsi deyil). **MVC məcburi deyil** — handler hər şeyi edib cavab qaytara bilər.

**MVC tarixçəsi (kitab qutusu):** Smalltalk, 1970-ci illərin sonu, Xerox PARC — Web-dən 10+ il əvvəl. Rails, CodeIgniter, Play, Spring MVC bunun üzərində qurulub. Yeni proqramçılar MVC-ni "web development-in yeganə yolu" sanır — yanlışdır: web app = HTTP üzərində istifadəçi ilə qarşılıqlı tətbiq, istənilən pattern (yaxud heç biri) işləyir.

**Template engine 2 növü:**
1. **Static/logic-less templates:** HTML + placeholder token-lər, template-də məntiq yoxdur — token-lər data ilə əvəz olunur (SSI məntiqi). Nümunə: CTemplate, **Mustache**
2. **Active templates:** placeholder + proqramlaşdırma konstruktları (conditionals, iterators, variables). Nümunə: JSP, ASP, ERB; PHP aktiv template engine kimi başlayıb dilə çevrilib

### 11. Hello Go — ilk web app
**Kitabdan kod nümunəsi (listing 1.1):**
```go
package main

import (
    "fmt"
    "net/http"
)

func handler(writer http.ResponseWriter, request *http.Request) {
    fmt.Fprintf(writer, "Hello World, %s!", request.URL.Path[1:])
}

func main() {
    http.HandleFunc("/", handler)
    http.ListenAndServe(":8080", nil)
}
```

**Sub-kod izahı:**
- `package main` → executable proqram (web app də daxil) main paketində olmalıdır
- `"net/http"` → Go-da web server **application server tələb etmir** (Ruby/Python/Javadan fərqli) — mühiti net/http təmin edir və kodla birlikdə kompilyasiya olunur → standalone deployable binary
- `func handler(writer http.ResponseWriter, request *http.Request)` → handler funksiyası (texniki olaraq "handler function" — Ch 3-də fərq açılıcaq):
  - `http.ResponseWriter` interface → HTTP response yazmaq üçün
  - `*http.Request` struct → sorğunun bütün məlumatları
- `fmt.Fprintf(writer, "Hello World, %s!", request.URL.Path[1:])` → ResponseWriter-a formatlanmış yazı; `request.URL.Path[1:]` → URL path-in ilk `/`-dan sonrakı hissəsi
- `http.HandleFunc("/", handler)` → root URL-i handler-a bağlayır
- `http.ListenAndServe(":8080", nil)` → 8080 portunda server başladır (dayandırmaq: Ctrl-C)

**Build və run (CLI):**
```bash
$ go install first_webapp
$ first_webapp          # $GOPATH/bin PATH-dədirsə
```
- `go install first_webapp` → `$GOPATH/bin`-də binary yaradır
- `http://localhost:8080/sausheong/was/here` → ekranda "Hello World, sausheong/was/here!"

**Kitabın kod reposu:** https://github.com/sausheong/gwp

## Əsas terminlər
- Web Application (web tətbiqi)
- Web Service (web servisi — proqramlara data qaytaran tətbiq)
- HTTP (Hypertext Transfer Protocol)
- Stateless Protocol (vəziyyətsiz protokol)
- Request-Response Model (sorğu-cavab modeli)
- Request Line / Status Line (sorğu vəziyyəti sətri)
- Request Header (sorğu başlığı)
- Safe Method (təhlükəsiz metod)
- Idempotent Method (eyniləşdirilən metod)
- Status Code (vəziyyət kodu)
- URI / URL / URN (ümumi resurs identifikatoru / lokator / adı)
- URL Encoding / Percent Encoding (faiz kodlaşdırması)
- Handler (emalçı)
- Template Engine (şablon mühərriki)
- Static/Logic-less Template (statik şablon)
- Active Template (aktiv şablon)
- MVC (Model-View-Controller)
- CGI (Common Gateway Interface)
- SSI (Server-Side Includes)
- HTTP/2 (SPDY əsaslı binary protokol)
- XMLHttpRequest / XHR (brauzer HTTP API-si)

## Praktik nəticə
- Web app = handler + template engine — bu ikili kitab boyu bütün arxitekturanın onurğasıdır.
- REST API dizaynında safe/idempotent fərqi kritikdir: GET/PUT/DELETE retry edilə bilər, POST yox.
- Response status kodlarını 1XX-5XX siniflərinə görə oxu; debugging-də 4XX client-in, 5XX serverin günahıdır.
- Go web app tək statik binarydir — DockerHeroku-dan əvvəl bu deploy sadəliyi Go-nun killer feature-idir.
- HTTP/2 avtomatik gəlir (HTTPS + Go 1.6+) — semantika dəyişmədiyi üçün kod eyni qalır.

## Mənbə
Pages: 24-42 (PDF), book pages 3-21
