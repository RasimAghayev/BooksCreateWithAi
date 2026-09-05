# Go Web Programming — Xülasə (Azərbaycan dilində)

**Müəllif:** Sau Sheong Chang
**Nəşriyyat:** Manning Publications Co., 2016 · **ISBN:** 9781617292569
**Səhifə:** 280 (+indeks) · **10 chapter, 3 hissə**

## Kitabın ümumi məqsədi

Go ilə real web tətbiqlərinin və servislərinin qurulması — standard library (`net/http` + `html/template`) üzərindən, framework-ə ehtiyac olmadan. Kitabın fəlsəfəsi: framework istifadə etməzdən ƏVVƏL altında yatan anlayışları başa düş — "cargo cult programming"-dən qaç.

## Chapter-by-chapter xülasə

### Ch 1 — Go and web applications (s. 3-21)
Web app tərifi (HTML qaytaran + HTTP ötürən proqram; web service isə proqrama data qaytarır). Go-nun 4 üstünlüyü: scalable (vertical — goroutine-lər; horizontal — statik binary), modular (interfeys + mikro-servis), maintainable (gofmt/godoc/gotest), high-performance (native kod). HTTP tam kurs: request/response strukturu, 9 metod (safe/idempotent təsnifatı ilə), header-lər, status kod sinifləri, URI strukturu, HTTP/2 (binary, multiplexed). Web app = handler + template engine; MVC tarixi və opcionallığı. Hello Go: 10 sətirdə web server.

### Ch 2 — Go ChitChat (s. 22-44)
Tam forum tətbiqi: multiplexer (NewServeMux), 11 route, static fayllar (FileServer+StripPrefix), cookie ilə session auth (CreateSession + Uuid + HttpOnly), session yoxlama utiliti, html/template (define/layout/navbar/content), generateHTML helper (variadic + empty interface), PostgreSQL (setup.sql, sql.Open connection pool, Threads() sorğusu, NumReplies metodu — template-dən metod çağırış). Data layer fəlsəfəsi: handler-lar heç vaxt birbaşa SQL yazmır.

### Ch 3 — Handling requests (s. 47-68)
net/http server tərəfi: Server struct (timeout, ErrorLog), HTTPS (ListenAndServeTLS + crypto/x509 ilə self-signed sertifikat), Handler interfeysi (ServeHTTP metodu), HandlerFun**c** adapteri (funksiya tipi — source kod göstərilir), zəncirləmə/middleware (`func(http.Handler) http.Handler` idiomu — log/protect nümunələri), ServeMux routing (trailing slash prefiks match!), HttpRouter (`:name` named params + metod-əsaslı route), HTTP/2 (Go 1.6+ avtomatik, http2.ConfigureServer + cURL --http2).

### Ch 4 — Processing requests (s. 69-95)
Request struct: URL, Header (map[string][]string — Set vs Add), Body (io.ReadCloser). HTML form enctype-ləri (urlencoded vs multipart MIME). Form/PostForm/MultipartForm cədvəli + ParseForm/ParseMultipartForm + **GOTCHA**: FormValue multipart body-dən oxumur. Fayl yükləmə (FormFile). JSON POST (Angular application/json — ParseForm tutmur!). ResponseWriter: Write (512-byte Content-Type sniffing), WriteHeader (sonra header dəyişməz), Header (redirect: Location+302). Cookie struct (Expires vs MaxAge), SetCookie, oxuma (r.Cookie/Cookies), flash mesajlar (Base64 + eyni adlı MaxAge:-1 cookie ilə silmə).

### Ch 5 — Displaying content (s. 96-124)
Template engine spektri (logic-less vs embedded logic; Go hibriddir). Parse/Execute dövrü (ParseFiles/ParseGlob/Parse/Must; Execute vs ExecuteTemplate). 4 action: conditional (if/else), iterator (range + else fallback), set (with — nöqtəni dəyişir), include (template + arg). Variables (`$key,$value :=`), pipeline (`|` — Unix bənzəri). FuncMap custom funksiyalar (parse-dən əvvəl attach!). **Context awareness**: eyni data HTML/URL/JS kontekstində fərqli escape — XSS müdafiəsi avtomatik; template.HTML ilə unescape + X-XSS-Protection. Nesting: define ilə layout, eyni adlı template fərqli fayllarda (content switch), block action (Go 1.6+ default template).

### Ch 6 — Storing data (s. 125-152)
Yaddaşda cache (map + pointer, 2 indeks = 2 map). Fayllar: CSV (encoding/csv — Flush mütləq, FieldsPerRecord) və gob (binary serializasiya, bytes.Buffer, empty interface). SQL: sql.Open (lazy connection pool) + blank import (driver-in init Register çağırır). Tam CRUD: Create (Prepare + returning id + QueryRow.Scan), GetPost (QueryRow), Update/Delete (Exec), Posts (Query + rows.Next iterator). One-to-many: FK + Comments slice + Post *Post pointer; Comment.Create post yoxdursa error. Sqlx (StructScan + db tag) və Gorm (Data-Mapper, AutoMigrate, Association, Related).

### Ch 7 — Go web services (s. 155-189)
SOAP vs REST: protokol vs fəlsəfə; WSDL kontraktı; envelope strukturu; REST: resurs+verb, CRUD↔HTTP mapping (PUT idempotent/URL-məlum, POST yeni URL), PATCH; action→resurs həlləri (reify: POST /user/456/activation; property: PATCH active). XML parse: 6 struct tag qaydası (XMLName, attr, chardata, innerxml, ad match, `a>b>c` leap-frog), Unmarshal vs Decoder (streaming, Token/StartElement/DecodeElement). XML create: Marshal/MarshalIndent + xml.Header manual; Encoder. JSON: `json:"key"` tag, Decoder (stream) vs Unmarshal (string). Tam RESTful web service: switch r.Method dispatch + 4 handler + http.Error 500 + path.Base id çıxarışı + cURL testləri.

### Ch 8 — Testing your application (s. 190-221)
testing.T: Error/Fatal/Skip; -short + testing.Short(); t.Parallel + -parallel N. Benchmark: b.N loop, `-bench . -run x`, Decode vs Unmarshal (25% fərq). HTTP testing: httptest.NewRecorder + http.NewRequest + mux.ServeHTTP — server qaldırmadan handler testi. TestMain (mərkəzi setUp/tearDown). **Dependency injection**: Text interfeysi + Post.Db sahəsi + handleRequest(t Text) http.HandlerFunc closure + FakePost test double → DB-siz müstəqil test. gocheck: Suite + SetUpTest/TearDownTest fixture-ləri + c.Check/Assert. Ginkgo+Gomega (BDD): user story → Describe/Context/It + BeforeEach + Expect().To(Equal()).

### Ch 9 — Leveraging Go concurrency (s. 223-254)
Concurrency (üst-üstə düşən icra, 1 CPU-da mümkün) ≠ parallelism (eyni anda, çox CPU). Goroutine: 8KB stack, multiplexed; **performans benchmark-ləri**: trivial işdə 78x yavaş (başlatma overhead), ağır işdə 7x sürətli, multi-CPU 2x daha da artır, 4 CPU isə bəzən 2 CPU-dan PİS. WaitGroup (Add/Done/Wait; Done yoxdursa deadlock). Channel-lər: unbuffered (sinxron qutu) vs buffered (FIFO, throughput throttle); thrower/catcher; select + default (deadlock həlli); close + `v, ok := <-ch` (düzgün həll). Foto-mozaika web app: sinxron 2.25s → concurrent fan-out/fan-in (cut x4 + combine + select + WaitGroup) 1 CPU-da 646ms (paralellik YOXDUR!), multi-CPU-da 216ms — 10x. Race condition → DB struct + mutex (axtarış+silmə bir kritik bölmədə).

### Ch 10 — Deploying Go (s. 256-280)
4 üsul: (1) **Standalone**: nohup (crash-ə qarşı kor) → Upstart init daemon (respawn, respawn limit, setuid, exec — kill-ə baxmayaraq avtomatik restart). (2) **Heroku**: os.Getenv("PORT") + godep save (Godeps.json asılılıqlar) + Procfile → git push. (3) **GAE**: sandbox (read-only FS, 60s limit, network yox) → main→init, MySQL driver + cloudsql DSN, `$1`→`?`, app.yaml, goapp serve/deploy. (4) **Docker**: VM vs container (OS-level virtualizasiya); Dockerfile (FROM golang, ADD, RUN go get/install, ENTRYPOINT, EXPOSE) → build → run (--publish 80:8080) → Docker Machine ilə cloud host (digitalocean driver, env eval). Müqayisə cədvəli: kod dəyişikliyi, sistem işi, scalability, platform bağlanqlılığı.

## Kitabın əsas mesajları

1. **Standard library kifayətdir** — net/http + html/template ilə tam funksional web app; framework anlayışını anlamadan istifadə etmə.
2. **Handler + Template Engine** — web app-in onurğa sütunu; middleware üçün `func(Handler) Handler` idiomu.
3. **Data layer ayrıdır** — handler-lar SQL yazmaz; struct + metod/funksiya = ORM-siz data qatı.
4. **Test üçün dizayn** — dependency injection (interfeys + closure) test double-lərə imkan verir.
5. **Concurrency dilin genindədir** — amma benchmark olmadan qəhrəmanlıq etmə; correct dizayn → paralellik sonra pulsuz gəlir.
6. **Deploy seçimi trade-off-dur** — sadəlik (Heroku) vs güclü məhdudiyyət (GAE) vs çeviklik (Docker) vs tam nəzarət (standalone).

## Ən dəyərli 5 fikir (az oxucu üçün)

1. **ServeMux trailing slash qaydası** — `/hello` dəqiq, `/hello/` prefiks; bunu bilmək mütləqdir, yoxsuroute-lar gözlənilməz handler-lara düşür.
2. **FormValue multipart-dən oxumur** — fayl yükləmədə FormFile/ParseMultipartForm istifadə et; bu gotcha real layihələrdə rast gəlinir.
3. **html/template context awareness** — XSS-ə qarşı avtomatik escape; text/template heç vaxt HTML üçün istifadə etmə.
4. **Dependency injection = test oluna bilən kod** — global sql.DB asılılığını interfeysə çevirmək FakePost testlərinin açarıdır.
5. **Goroutine başlatmaq pulsuz deyil** — trivial işdə 78x itki; concurrency yalnız kifayət qədər ağır işə layiqdir.

## Kitabın ən dəyərli hissəsi

**Chapter 3 + Chapter 4 (Handling/Processing requests), pages 47-95** — çünki:
- Handler vs HandlerFunc vs ServeMux münasibətlərini net/http source kodundan göstərərək açır — başqa mənbələrdə adətən qaranlıq qalan fərq.
- Middleware idiomu (`func(http.Handler) http.Handler`) bütün Go web framework-lərinin (gin, echo, chi) arxitektur əsasıdır — bu 2 chapter onları anlaməğın açarıdır.
- Form/PostForm/MultipartForm cədvəli + ResponseWriter Write/WriteHeader/Header sırası — gündəlik işin 80%-ni əhatə edən praktik bilikdir.
