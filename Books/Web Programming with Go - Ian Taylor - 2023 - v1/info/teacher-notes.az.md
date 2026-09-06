# Web Programming with Go — Müəllim Qeydləri (Teacher Notes)

> Bu sənəd kitabı öyrədərkən istifadə üçün metodiki qeydlər, çətin anlar, müzakirə sualları və praktik tapşırıqlar.

---

## 1. Auditoriya və ön şərtlər

- **Hədəf auditoriya:** Digər dillərdə təcrübəli proqramçılar + Go başlanğıc; veb inkişaf əsasları (HTTP, HTML) fərz olunur.
- **Ön şərtlər:** Go sintaksisi (kitab öyrətmir!), HTML form-lar, əsas SQL.
- **Həcm:** 264 səh., 9 fəsil — semestr kursu və ya 4-5 günlük intensiv.

## 2. Tədris axını — layihə-əsaslı model

Kitab təbii kurs kürəsidir: hər dərs bookstore-ə yeni qat qoyur:

| Həftə | Chapter | Çıxış məhsulu |
|-------|---------|---------------|
| 1 | 1-2 | İşləyən server + modul + qovluq strukturu + App obyekti |
| 2 | 3 | Mux routing + middleware + CRUD + rate limit |
| 3 | 4 | Şablonlar + form + kesh + XSS qoruma |
| 4 | 5 | DB sxemi + CRUD + indekslər + pool |
| 5 | 6 | Paralel axtarış + WaitGroup/Mutex + keş |
| 6 | 7 | Sessiya + bcrypt + OAuth |
| 7 | 8 | JSON API + JWT + xarici inteqrasiya |
| 8 | 9 | Test piramidi + mock + pprof + xəta təbabəti |

**Praktik struktur:** Hər dərs = konsept izahı (15 dəq) → kitabın kodunun birgə yazılması (30 dəq) → tələbə öz variantını genişləndirir (45 dəq).

## 3. Çətin anlar və izah üsulları

### a) Niyə interface? (Ch 2-3 DI mock)
Tələbələr "niyə BookService birbaşa SQLBookRepository istifadə etmir?" soruşur. Canlı demo: MockDatabase ilə test yaz → real DB olmadan 0.01 s test. "Interface = dəyişən nöqtə" aforizmi.

### b) Middleware-in "soğan" modeli (Ch 3)
Birgə diaqram çəkin: sorğu → [Logging → RateLimit → Auth] → handler → geri. Sonra kodda handler-ı bir neçə dəfə "sarıb" göstərin — LoggingMiddleware(RateLimit(AuthMiddleware(handler))). İç-içe funksiya çağırışının vizualizasiyası.

### c) Mutex niyə lazım, kanal varkı? (Ch 6)
SAYĞAC demosu: 1000 goroutine `counter++` — nəticə hər dəfə FƏRQLİ (race!). Mutex əlavə et → stabil. Kanal vs Mutex seçim qaydası: data AXINI → kanal; DURUM qoruma → mutex.

### d) Template caching fərqi (Ch 4)
Benchmark canlı: parse-hər-sorğu vs keshli — ekranda görünən fərq təsirli olur. Sonra "dev-də dəyişiklik görünmür" problemini live-da göstər → fresh aləti.

### e) SQL injection (Ch 5)
Tələbələrə string-concat sorğunu "hacker" rolunda sındırın: `"'; DROP TABLE books; --"`. Sonra parameterized ($1) versiya — eyni hücum zərərsiz. Heç bir izah kəşfiyyət qədər təsirli deyil.

### f) Cookie üçlüyü niyə? (Ch 7)
Hər atributun hücumunu ayrıca göstər: HTTPOnly-siz → JS konsolunda `document.cookie` (XSS); Secure-siz → Wireshark HTTP-də cookie görünür; SameSite-siz → 3-cü saytdan sorğu cookie daşıyır. Hər atribut = konkret hücumun qarşısı.

### g) N+1 query (Ch 9)
Kiçik DB-də 100 kitab + 100 müəllif sorğusu: sekvensial ~3s, JOIN ~0.1s. Pprof/log ilə sübut. "Bookshelf" gerçək dataset böyüdükcə eksponensial pisləşməni müzakirə.

## 4. Müzakirə sualları

1. Bookstore üçün REST ya GraphQL? Nə vaxt GraphQL üstün olur?
2. `internal/` qovluğu compiler tərəfindən necə qorunur? Niymə design qərarıdır?
3. Sessiya Redis-də saxlanırsa, Redis çöksə nə olur? Fallback strategiyaları?
4. Rate limit PER-USER etmək üçün limiter harada saxlanmalı? (qlobal map + mutex? Redis?)
5. JWT stateless-dir — amma logout necə işləyir? (qara siyahı paradoksu → token expiry trade-off)
6. GORM vs xam SQL: 1000-objektlik layihədə hansı? Nə vaxt ORM-dən qaçmalı?
7. Worker pool ölçüsü necə seçilir? (CPU core? I/O gözləmə? benchmark cavabı)
8. Mock testlər interface-i 100% əks etdirmir — hansı bug-lar yalnız integration test-də üzə çıxır?

## 5. Praktik tapşırıqlar

**Tapşırıq 1 (server + modul):** "Library" layihəsi qur: go mod init, qovluq strukturu, App obyekti, health check, .env konfiqurasiya.

**Tapşırıq 2 (routing):** CRUD API: mux + subrouter (books, users) + metod məhdudiyyəti + custom 404 + logging middleware. Bonus: rate limit.

**Tapşırıq 3 (şablon):** Kitab kataloq səhifəsi: nested templates (base/header/footer), keshlənmiş, kitab kartları loop. Bonus: axtarış formu + validation.

**Tapşırıq 4 (DB):** 6-cədvəllik sxemi qur (SERIAL, FK, CHECK); 5 CRUD endpoint; composite index; pool tuning + Ping.

**Tapşırıq 5 (concurrency):** 3 mənbədən paralel kitab axtarışı (kanal pattern); shared counter-u race-siz et; ConcurrentCache yaz + test.

**Tapşırıq 6 (auth):** bcrypt qeydiyyat/login; sessiya cookie-si (üçlük atributlarla); JWT ilə API qoruma; OAuth (Google) — bonus.

**Tapşırıq 7 (API):** Kitab API: struct tag-lər, pagination, JWT auth middleware; xarisi review API inteqrasiyası; error pattern (Result{Data,Err}).

**Tapşırıq 8 (test):** Bütün keçən tapşırıqlar üçün table-driven testlər; MockDatabase; paralel test; pprof ilə bir bottleneck tap və optimallaşdır.

**Final layihə:** Kitabın bookstore-unu KOMPLEKT qur + 1 orijinal xüsusiyyət (missiya: wishlist, tövsiyə mühərriki, admin dashboard) + test əhatəsi ≥ 60%.

## 6. Sınav sualları

**Asan:**
1. `go mod init` nə yaradır? (go.mod)
2. HTTPOnly atributu nədən qoruyur? (XSS — JS cookie oğurluğu)
3. `rows.Scan` nə edir? (sətir sahələrini dəyişənlərə köçürür)

**Orta:**
4. Subrouter nə üçün lazımdır? (qrup + qrupa-özəl middleware)
5. Parameterized query nəyə görə təhlükəsizdir? (input data kimi, kod kimi YOX icra olunur)
6. WaitGroup.Add(1) niyə goroutine-dən ƏVVƏL çağrılmalı? (yarış: Done Add-dən əvvəl gələ bilər)

**Çətin:**
7. N+1 query-ni JOIN ilə həll et — amma JOIN-in öz hansı riskləri var? (böyük kartez məhsulu, indeks ehtiyacı)
8. JWT logout problemi: token qara siyahı stateless-i POZUR — həllər? (qısa expiry + refresh rotasiyası)
9. 100 user üçün per-user rate limit: harada saxlansın? (in-memory map+mutex vs Redis) — trade-off analizi

## 7. Kollektiv layihə ideyası

**"Mini-Kitabxana müharibələri":** Qruplar bookstore-un fərqli aspektlərini paralel qurur (bir qrup backend API, digəri auth, üçüncüsü test). Git branch-lər + code review ritualı + integration həftəsi. Real komanda işi simulyasiyası.

## 8. Əlavə resurslar

- Go rəsmi tur: go.dev/tour
- Effective Go: go.dev/doc/effective_go
- gorilla/mux docs: gorillatoolkit.org
- database/sql tutorial: go.dev/doc/database
- OWASP Top 10: owasp.org
- jwt.io — JWT debug aləti
- Testify: github.com/stretchr/testify
- Stripe Go docs: stripe.com/docs/api?lang=go

## 9. Kitabın dəyəri və limitləri (müəllim qiyməti)

**Güclü:** Tək davamlı layihə; praxis-biased; bütün standart veb komponentlər təmsil olunur; təhlükəsizlikə real diqqət; kod nümunələri iştirakçıdır.

**Zəif/həssas:** Bəzi bölmələr təkrar xarakterlidir (REST iki dəfə tam izah olunur); jwt-go köhnə lib-dir (golang-jwt modern); kod nümunələrdə kiçik qeyri-tutarlılıqlar (Book vs Name sahələri); deployment/cloud mövzusu yoxdur; dərin ORM (GORM) kodu azdır.
