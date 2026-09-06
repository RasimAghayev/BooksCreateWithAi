# Go Web Programming — Terminologiya lüğəti (AZ)

`English (Azərbaycanca)` formatında — SYSTEM_PROMPT qayda 4 üzrə.

## HTTP əsasları (Ch 1)
- Web Application / Web Service (web tətbiqi / servisi)
- HTTP (stateless, text-based, request-response protokolu)
- Request Line / Status Line (sorğu / vəziyyət sətri)
- Request Header (sorğu başlığı): Accept, Cookie, Host, Content-Type, User-Agent...
- Response Header: Allow, Location, Set-Cookie, WWW-Authenticate...
- Status Code sinifləri: 1XX informational, 2XX success, 3XX redirection, 4XX client error, 5XX server error
- Safe Method (təhlükəsiz metod): GET, HEAD, OPTIONS, TRACE
- Idempotent Method (eyniləşdirilən): safe + PUT, DELETE (POST YOX)
- URI / URL / URN; Query / Fragment
- URL Encoding / Percent Encoding (faiz kodlaşdırması)
- HTTP/2 (SPDY əsaslı binary, multiplexed)
- XMLHttpRequest / XHR
- CGI (Common Gateway Interface)
- SSI (Server-Side Includes)
- Handler (emalçı)
- Template Engine (şablon mühərriki)
- Static/Logic-less Template (statik şablon)
- Active Template (aktiv şablon)
- MVC (Model-View-Controller)

## Go web server (Ch 2-3)
- Multiplexer / ServeMux / DefaultServeMux (çoxlayıcı / router)
- Handler / Handler Function (interfeys / imzalı funksiya)
- HandlerFunc (adapter tipi — funksiya → Handler)
- Middleware / Chaining (zəncir / pipeline)
- FileServer / StripPrefix (statik fayl servisi)
- Session Cookie (sessiya kökəyi)
- HttpOnly (yalnız-HTTP bayrağı)
- Cargo Cult Programming (kor-koranə kod kopyalama)
- Cross-cutting Concern (kəsən qayğı)
- Principle of Least Surprise (ən az təəccüb prinsipi)
- Named Parameter (`:name` — HttpRouter)
- SSL / TLS / HTTPS / X.509 / PEM / DER / ASN.1
- Certificate Authority (CA)
- PCI DSS (kredit kartı təhlükəsizlik standartı)

## Request emalı (Ch 4)
- Request struct (URL, Header, Body, Form, PostForm, MultipartForm)
- enctype: application/x-www-form-urlencoded / multipart/form-data
- FormValue / PostFormValue / FormFile (qısayol metodları)
- io.ReadCloser (Body interfeysi)
- ResponseWriter: Write / WriteHeader / Header
- Content-Type Sniffing (ilk 512 bayta görə aşkarlama)
- Session / Persistent Cookie
- Expires / MaxAge
- ErrNoCookie
- Flash Message (müvəqqəti mesaj)
- Base64 URL Encoding

## Template-lər (Ch 5)
- text/template vs html/template
- Action (`{{ }}`)
- Dot (nöqtə — data)
- Conditional / Iterator (range) / Set (with) / Include / Define / Block actions
- Argument / Variable (`$x`) / Pipeline (`|`)
- FuncMap (custom funksiyalar xəritəsi)
- template.Must (panic helper)
- Template Set (şablon dəsti)
- ParseFiles / ParseGlob / Execute / ExecuteTemplate
- Context Awareness (kontekst-fəqliliyi)
- XSS (Cross-Site Scripting) / Persistent XSS
- template.HTML (unescape typecast)
- Layout / Nested Templates
- X-XSS-Protection header

## Data saxlama (Ch 6)
- In-memory Storage (yaddaş saxlama) — cache
- CSV (Comma-Separated Values)
- gob (Go binary serializasiya)
- bytes.Buffer (həm Reader həm Writer)
- Connection Pool (`sql.DB`)
- Data Source Name (DSN)
- Database Driver / Blank Import (`_ "lib/pq"`)
- Prepared Statement (`$1` placeholder)
- QueryRow / Query / Exec / Scan
- Rows Iterator
- CRUD (Create/Retrieve/Update/Delete)
- One-to-many / Many-to-one / Foreign Key
- ORM (Object-Relational Mapper)
- ActiveRecord vs Data-Mapper pattern
- Struct Tag (`db:"..."` / `sql:"..."`)
- AutoMigrate (Gorm)

## Web services (Ch 7)
- SOAP (envelope + WSDL protokolu)
- WSDL (Web Service Definition Language)
- WS-* (WS-Security, WS-Addressing)
- UDDI (kataloq servisi)
- REST (Representational State Transfer)
- Resource / Verb
- Idempotency (PUT idempotent, POST yox)
- Reify (əməliyyatı resursa çevirmək)
- PATCH (qismən yeniləmə)
- WADL / Swagger / RAML
- XML struct tag mode flag-lər: `attr`, `chardata`, `innerxml`
- Leap-frog (`a>b>c` — aşma sintaksisi)
- xml.Name / XMLName
- xml.Header (declaration sabiti)
- Decoder / Encoder (streaming)
- io.EOF
- Method-based Dispatch (switch r.Method)

## Testing (Ch 8)
- testing.T / testing.B / testing.M
- Error/Errorf (fail + davam) vs Fatal/Fatalf (fail + dayan)
- Skip / testing.Short() (-short)
- Parallel (`t.Parallel()`, `-parallel N`)
- Benchmark (`-bench`, `b.N`, ns/op, `-cpu`)
- Coverage (`-cover`)
- ResponseRecorder (httptest.NewRecorder)
- http.NewRequest (test sorğusu yaratma)
- TestMain (mərkəzi setup/teardown)
- Test Double (test ikiləsi — FakePost)
- Dependency Injection (asılılığın inyeksiyası — Text interfeysi)
- Test Fixture (sınaq şəraiti)
- SetUpTest / TearDownTest / SetUpSuite / TearDownSuite (gocheck)
- Check / Assert (gocheck assertion-ları)
- BDD (Behavior-Driven Development)
- User Story / Scenario / Given-When-Then
- Describe / Context / It / BeforeEach (Ginkgo)
- Matcher (Gomega) / Expect

## Concurrency (Ch 9)
- Concurrency vs Parallelism (eyni-vaxtlılıq vs paralellik)
- GOMAXPROCS (CPU sayı)
- Goroutine (multiplexed, yüngül stack)
- WaitGroup: Add / Done / Wait
- Deadlock ("all goroutines are asleep")
- Unbuffered (sinxron) / Buffered (FIFO) Channel
- Send-only / Receive-only (`chan <-` / `<-chan`)
- select / default case
- close / two-value receive (`v, ok := <-ch`)
- Zero Value (bağlı kanalın qaytardığı)
- Race Condition
- Mutual Exclusion / Mutex / Critical Section
- Fan-out / Fan-in
- Data URL (base64 embed)
- Euclidean Distance

## Deploy (Ch 10)
- IaaS / PaaS / SaaS (NIST cloud modelləri)
- VPS (Virtual Private Server)
- nohup (HUP siqnalına laqeydlik)
- Init Daemon / Upstart / systemd / Stanza / Respawn
- Slug / Dyno (Heroku)
- Procfile / Godep / Godeps.json
- Sandbox (GAE məhdudiyyətləri)
- Google Cloud SQL
- app.yaml / goapp serve / goapp deploy
- `$1` vs `?` (PostgreSQL vs MySQL placeholder)
- Docker: client / daemon / host / image / container
- Dockerfile: FROM / ADD / WORKDIR / RUN / ENTRYPOINT / EXPOSE
- Registry / Docker Hub
- --publish / --rm
- Docker Machine
- Personal Access Token
