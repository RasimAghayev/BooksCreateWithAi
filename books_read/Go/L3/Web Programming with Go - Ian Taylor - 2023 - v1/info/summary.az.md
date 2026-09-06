# Web Programming with Go — Xülasə (Azərbaycanca)

> **Kitab:** Web Programming with Go: Building and Scaling Interactive Web Applications with Go's Robust Ecosystem — Ian Taylor, GitforGits, 2023 (ISBN 978-93-91526-99-9) · 264 səh.
> 🎯 Intermediate (3/5) · Praktik layihə-kitab: **GitforGits Bookstore** app-i fəsil-fəsil qurulur.

---

## Kitabın yanaşması
Tək kitab boyu davam edən bookstore layihəsi vasitəsilə öyrənmə: hər fəsil app-ə yeni qat əlavə edir — server → struktura → routing → templating → DB → concurrency → auth → API → testing.

## Fəsl-fəsl xülasə

### Ch 1 — Giriş: Go vebdə niyə?
Go-nun üstünlükləri (concurrency, tək binary, standard library, sadəlik), toolchain (go build/test/mod/get), Linux qurulumu, net/http icmalı, REST prinsipləri, **ilk server**: HandleFunc + ListenAndServe + handler.

### Ch 2 — Struktur: Go Modules + arxitektura
go.mod/go.sum/proxy; standart qovluq strukturu (cmd/internal/pkg/web); vendoring; **MVC Go-da**; App obyekti (Router+DB+Config struct); konfiqurasiya env/.env (godotenv); **Dependency Injection** interface-lərlə (BookRepository → BookService); health check endpoint-ləri.

### Ch 3 — Routing: HTTP sorğuları
Metodlar (GET/POST/PUT/PATCH/DELETE); request-response dövrü; Handler/HandlerFunc; **middleware pattern** (logging/auth); Gorilla/Mux: regex-li path parametr (`{id:[0-9]+}`), metod məhdudiyyəti, **subrouter qrupları**; custom 404/500 (NotFoundHandler, panic-recover middleware); **rate limiting** (token bucket 429) + throttling; CRUD handler-ləri (201/204/400/404).

### Ch 4 — Templating: dinamik məzmun
Dinamik renderinq (SEO+interaktivlik); variables/loops/conditions (`{{range}}/{{if}}`); **nested şablonlar** (base+header+content+footer); form emalı (FormValue, HTML5 validation); **template caching** (ParseGlob başlanğıcda — sorğu başına parse YOX); **XSS qoruması** (html/template contextual encoding; template.HTML ehtiyatla); debug ssenariləri (parse/runtime/məntiq xətaları).

### Ch 5 — Database: data qatı
database/sql + driver (lib/pq); Query+Scan / Exec / transaksiyalar (Begin/Commit/Rollback); **tam bookstore sxemi** (6 cədvəl: books, genres, users(password_hash), reviews, transactions, inventory); GORM alternativi; **indeksləmə** (primary/secondary/composite/covering/partial); **connection pool tuning** (MaxOpen/MaxIdle/Lifetime + Ping); advanced SQL (JOIN, subquery, GROUP BY/HAVING).

### Ch 6 — Concurrency: paralellik
Goroutine + channel əsasları; **paralel janr axtarışı** nümunəsi; yönlü/buffered/locked kanallar; **WaitGroup** (bitmə gözləmə) + **Mutex** (paylaşılan resurs); birgə pattern (bulk sifariş emalı); **ConcurrentCache** implementasiyası (map+mutex+eviction); concurrency vs parallelism fərqi; do's/don'ts (race, panic, worker pool, graceful shutdown, test).

### Ch 7 — Auth: sessiya və identitasiya
Sessiya storage seçimləri (cookie/memory/file/DB/Redis); **təhlükəsiz cookie üçlüyü** (HTTPOnly+Secure+SameSite); cookie encrypt + expiration; **bcrypt qeydiyyat/login**; parol qorunması 5 qat (hash+salt+pepper+iterations+modern alqoritm); **OAuth 2.0** tam axını (Config, AuthCodeURL, callback, Exchange, refresh token).

### Ch 8 — API: frontend-backend körpüsü
REST təkrarı; **ilk API endpoint** (struct tag-lı Book + query parametr + JSON); **JSON Marshal/Unmarshal**; DB-backed endpoint (Query→Scan→Marshal); **paralel data fetch** (goroutine+buffered channel+**Result{Data,Err}** pattern); **JWT auth** (GenerateToken/ValidateToken/middleware 401); xarici API inteqrasiyası (http.Get, **Stripe** ödəniş, rate limit idarəsi).

### Ch 9 — Test: keyfiyyət təminatı
testing.T əsasları (Errorf/Fatalf/Skipf); **table-driven testlər** + subtest (t.Run); helper setup/teardown; t.Parallel; **interface əsaslı mock** (Database interface → MockDatabase); logging (standart log → **logrus** structured JSON → ELK); **pprof üçlüyü** (CPU/block/memory); **10 klassik veb xətası** və həlləri (N+1, race, leak, auth, validation, error handling, dependency, config, concurrency, pagination).

### Ch 10 — Epiloq
Yekun mesajlar: ömürboyu öyrənmə, problem-həllərçi mentalitet, Go cəmiyyətinə töhfə.

## Kitabın əsas mesajları
1. **Layihə-əsaslı öyrənmə:** Bookstore app-i böyüyən real senari — hər texnika kontekstdə öyrənilir.
2. **Struktur əvvəl:** modullar → qovluq strukturu → App obyekti → sonra funksionallıq.
3. **Go idiomatikaları:** interface əsaslı DI, middleware sarğı pattern-i, mux subrouter qrupları.
4. **Təhlükəsizlik birinci sinif vətəndaş:** bcrypt+salt+pepper, cookie üçlüyü, XSS/CSRF qoruma, parameterized SQL.
5. **Concurrency Go-nun üstünlüyüdür:** goroutine+channel+WaitGroup+Mutex tam arsenalı real ssenarilərdə.
6. **Performans əməkdir:** keshlənmiş şablonlar, connection pool, indekslər, pprof, rate limit.

## Kitabdan sonra
- Deployment mövzuları (Docker, K8s) — kitab toxunmur, təbii davam
- Advanced goroutine patternləri (errgroup, context cancellation)
- GraphQL API-lər
- Real-world Go veb frameworkləri (Echo, Gin, Chi) müqayisəsi
