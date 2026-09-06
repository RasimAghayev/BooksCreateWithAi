# Web Programming with Go — Terminologiya (Azərbaycanca)

> Texniki terminlər `English (Azərbaycanca qarşılığı)` formatında — kitabın ardıcıllığı ilə.

## Layihə və modullar (Ch 1-2)
- **Go Modules (Go modulları)** — v1.11+ dependency idarəsi; go.mod manifest
- **go.sum** — kriptoqrafik checksum reyestri
- **Module proxy (modul proksisi)** — keshləyici mediator
- **Vendoring (vendinq)** — dependency-lərin /vendor qovluğuna kopyalanması
- **GOPATH** — köhnə workspace modeli
- **MVC (Model-View-Controller)** — arxitektura pattern-i
- **Dependency Injection (asılılıq inyeksiyası)** — xaricdən dependency təminatı
- **Inversion of Control** — asılılıq idarəsinin tərsinə çevrilməsi
- **Health check endpoint** — sistem sağlamlıq yoxlaması URL-i

## HTTP və Routing (Ch 1, 3)
- **HTTP verb (HTTP fel)** — GET/POST/PUT/DELETE/PATCH/HEAD/OPTIONS
- **Idempotence (eynilik)** — təkrar sorğu = eyni nəticə
- **Request-Response cycle** — sorğu-cavab dövrü
- **Handler (idarəedici)** — ServeHTTP interfeysli obyekt
- **HandlerFunc** — funksiyadan handler-ə adapter
- **ServeMux** — standart multiplexer
- **Routing (yönləndirmə)** — URI+metod → handler
- **Path parameter / Query parameter** — yol/dot parametri
- **Subrouter (alt router)** — prefiksli route qrupu
- **Middleware (arayataq)** — handler-ı saran emal qatı
- **CRUD** — Create/Read/Update/Delete

## Rate limiting (Ch 3)
- **Rate limiting (sürət limiti)** — vaxt vahidində maksimum sorğu
- **Request throttling (tempo tənzimi)** — sorğu axınının nizamlanması
- **Token bucket (jeton çəni)** — rate limit alqoritmi
- **DDoS** — paylanmış xidmət inkarı hücumu

## Templating (Ch 4)
- **Dynamic rendering (dinamik renderinq)** — server+client hibrid render
- **SPA (Single Page Application)** — tək səhifəli app
- **Nested template (iç-içe şablon)** — define/template quruluşu
- **Template caching (şablon keşi)** — parse bir dəfə, icra çox
- **XSS (Cross-Site Scripting)** — skript inyeksiya hücumu
- **Contextual encoding (kontekstual kodlaşdırma)** — yerə görə escape
- **template.HTML/JS/CSS/URL** — escape-siz etibar tipləri
- **Delims** — `{{ }}` ayracı

## Verilənlər bazası (Ch 5)
- **database/sql** — generic DB interfeysi
- **Driver (sürücü)** — lib/pq (PostgreSQL), go-sql-driver (MySQL)
- **Blank import (`_ "..."`)** — yalnız register məqsədli import
- **Foreign key (xarici açar)** — REFERENCES əlaqəsi
- **Constraint (məhdudiyyət)** — CHECK/UNIQUE/DEFAULT
- **SERIAL PRIMARY KEY** — avtoartan unikal ID
- **DECIMAL(10,2)** — dəqiq pul tipi (float YOX)
- **ORM (Object-Relational Mapping)** — GORM
- **Primary/Secondary/Composite/Covering/Partial index** — indeks növləri
- **Connection pooling (bağlantı hovuzu)** — MaxOpen/MaxIdle/ConnMaxLifetime
- **INNER/LEFT JOIN** — cədvəl birləşdirmə
- **Subquery (alt sorğu)** — iç-içə SELECT
- **GROUP BY / HAVING** — qruplaşdırma / sonrakı filtr (WHERE-dən fərqli!)

## Paralellik (Ch 6)
- **Goroutine** — yüngül runtime thread
- **Channel (kanal)** — tipəlyi data borusu
- **Buffered channel (buferli kanal)** — tutumlu `make(chan T, n)`
- **Channel direction** — `chan<-` (göndər) / `<-chan` (qəbul)
- **select** — çoxkanal multiplexing
- **sync.WaitGroup** — bitmə gözləyici sayğac
- **sync.Mutex** — qarşılıqlı istisna kilidi
- **Race condition (yariş şəraiti)** — eynizamanı oxu/yazı zərəri
- **Worker pool** — məhdud goroutine hovuzu
- **Concurrency vs Parallelism** — idarəetmə vs eynizamanlılıq
- **b.RunParallel** — paralel benchmark
- **Eviction (sıxışdırma)** — keşdən element silmə (LRU)

## Sessiya və təhlükəsizlik (Ch 7)
- **Session (sessiya)** — istifadəçi vəziyyət davamlılığı
- **Cookie attributes:** **HTTPOnly** (JS girişi yox), **Secure** (yalnız HTTPS), **SameSite** (CSRF qoruması)
- **bcrypt** — parol hash alqoritması
- **Hashing** — biryollu transformasiya
- **Salting (duzlama)** — unikal təsadüfi əlavə → rainbow table qarşısı
- **Pepper (istiot)** — DB-dən ayrı saxlanan gizli əlavə
- **Rainbow table** — hazır hash cədvəli hücumu
- **CSRF (Cross-Site Request Forgery)** — saxta sorğu hücumu
- **Man-in-the-middle** — aradinsan hücumu
- **OAuth (Open Authorization)** — token-əsaslı 3-cü tərəf auth
- **Access / Refresh token** — qısa/uzunömürlü token cütü
- **AuthCodeURL / Exchange** — kod-token mübadiləsi

## API və JSON (Ch 8)
- **API (Application Programming Interface)**
- **REST** — Representational State Transfer (Roy Fielding, 2000)
- **Statelessness (vəziyyətsizlik)** — sorğunun özünü kifayət etməsi
- **Uniform interface** — vahid interaksiya qaydası
- **Layered system (qatlı sistem)** — qat-qat şəffaflıq
- **JSON (JavaScript Object Notation)** — data mübadilə formatı
- **Struct tag** — `json:"field_name"` sahə xəritələnməsi
- **Marshal / Unmarshal** — encode / decode
- **JWT (JSON Web Token)** — HS256 imzalı token
- **Payment gateway** — Stripe ödəniş inteqrasiyası
- **Rate limit / Quota** — xarici API məhdudiyyətləri

## Test və Debug (Ch 9)
- **testing.T / testing.B** — test / benchmark tipləri
- **Table-driven test** — cədvəl əsaslı ssenari testi
- **Subtest (t.Run)** — adlı alt-test
- **t.Parallel** — paralel test icrası
- **Errorf / Fatalf / Skipf** — test axın idarəsi
- **Mock** — asılılıq simulyasiyası (interface əsaslı)
- **Setup/Teardown** — test öncəsi/sonrası kontekst
- **Structured logging** — JSON formatlı log (logrus/zap)
- **OpenTracing / Span** — sorğu izi vahidi
- **pprof** — CPU/block/memory profiler
- **Delve** — Go debugger
- **N+1 query** — list + element-sorğu antipattern-i
- **Memory leak (yaddaş sızmısı)**
- **Request-ID** — sorğu izləmə identifikatoru
