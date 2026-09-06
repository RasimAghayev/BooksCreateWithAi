# Go Web Programming — Müəllim qeydləri (Teacher Mode)

📖 Kitab deyir:
Go ilə web proqramlaşdırma standard library (net/http + html/template) üzərindən öyrənilməlidir — framework-lər arxasındakı anlayışları başa düşmədən istifadə etmək cargo cult proqramlaşdırmadır. Web app = handler + template engine. Handler-ları zəncirlə, request-i Form/PostForm/MultipartForm ilə emal et, ResponseWriter-ə yaz, cookie ilə sessiya idarə et, database/sql ilə CRUD qur, REST üçün method dispatch et, httptest + dependency injection ilə test et, concurrency ilə perfomansı artır və 4 fərqli üsulla deploy et.

👨‍🏫 Müəllim qeydi:
Məncə bunu praktikada belə tətbiq etmək daha yaxşı olar:

1. **Kitab Go 1.4-1.6 dövründə yazılıb — bir çox detallar artıq dəyişib.** Ən vacib dəyişikliklər:
   - GOPATH/godep dövrü keçib — bu gün `go mod init` + `go.mod`. Heroku deploy üçün Go buildpack modules-u avtomatik işlədir, Procfile hələ də gərəklidir.
   - HTTP/2 indiki Go-da tam avtomatikdir (HTTPS = HTTP/2, heç bir ConfigureServer lazım deyil).
   - `http.Server` üçün `ReadHeaderTimeout` (Go 1.8) əlavə olundu — slowloris hücumlarına qarşı bunu da qoy.
   - `ioutil.ReadAll/ReadFile/WriteFile/Discard` → `io.ReadAll`, `os.ReadFile`, `os.WriteFile`, `io.Discard` (ioutil deprecated).
   - `io/ioutil` importlarını yeni kodda işlətmə.

2. **Sinxron handler-ların bir problematı var: kontekst.** Kitabın handler-ları `*http.Request`-dən başqa heç nə qəbul etmir — bu gün `http.Handler`-ları context.Context ilə yazmaq standartdır (qrafik: `func handler(w, r)` → `r.Context()` ilə cancel/deadline alın). Ch 9-da ki "konkurrent goroutine-lərin idarəsi" mövzusunda bu gün əvvəl `r.Context().Done()` kanalını select-ə əlavə etmək lazımdır.

3. **Şablonlar üçün müasir cərəyan:** kitabın hər handler-də `template.ParseFiles` çağırması productionda hər request üçün disk I/O deməkdir. Bu gün: (a) `embed.FS` (Go 1.16+) ilə template-ləri binary-yə daxil et + startda bir dəfə parse; (b) yaxud html/template-in `Clone()` funksiyası. `template.Must` başlanğıcda, bir dəfə — hər requestdə YOX.

4. **Databasenin ölü bağlantısı üçün:** kitabdakı `Db.SetMaxOpenConns`/`SetMaxIdleConns`/`SetConnMaxLifetime` qeyd olunmur — productionda bunlar mütləq konfiqurasiya olunmalıdır, yoxsur pool açıq bağlantılarla şişir.

5. **Handler-ların imzaları:** `handleRequest(t Text) http.HandlerFunc` kitabın DI nümunəsi düzgündür, amma bu gün hansısa DI framework (wire, fx) YOX — sadə konstruktor funksiyaları (`NewServer(deps)`) Go-da idiomatik DI-dır. FakePost nümunəsi də bu gün testify mock ilə əvəz oluna bilər — amma interfeysin rolu dəyişməyib.

6. **Testing:** gocheck əsasən tarixi əhəmiyyət daşıyır; testify (assert/require) bu gün de-facto standartdır. Ginkgo/Gomega hələ də canlı, amma çox layihə üçün sadə testing + testify kifayətdir. Kitabın TestMain nümunəsi Go 1.15-dən bəri `t.Cleanup()` ilə daha sadə tətbiq olunur.

7. **Deploy bölməsi (Ch 10) tarixən ən çox köhnəlib:**
   - Heroku Go buildpack artıq `go.mod` istifadə edir; `godep` arxaikdir.
   - GAE Go 1.12+ üçün tam standart Go runtime dəstəkləyir — artıq `init()` hack-inə, `goapp`-a, `$1`→`?` dəyişikliyinə ehtiyac YOX (standard Go + gcloud app deploy). Kitabdakı köhnə GDE məhdudiyyətləri (60s, read-only FS) 1st-gen sandbox-a idi, 2nd-gen+ freewheelingdir.
   - Upstart yerinə systemd (Ubuntu 15.04+); Dockerfile-da `golang` base yerinə multi-stage build (golang → scratch/alpine) → ~10-15MB final image (kitabın 534MB nümunəsindən fərqli olaraq!) — bu, Go-nun static binary üstünlüyünün əsl ortaya çıxmasıdır.
   - Kubernetes bu gün cloud orkestrasiyasının dominantıdır — kitabda yoxdur (o dövrdə yenidən doğururdu).

8. **HttpRouter qeydi:** httprouterperformanslıdır, amma bu gün standart kitabxanada Go 1.22+ pattern matching VAR: `mux.HandleFunc("GET /post/{id}", ...)` — method + path parametri daxilində. Üçüncü tərəf routerə ehtiyac azaldı.

## Ən vacib 5 fikir

1. **Anlayış > framework:** kitabın əsas mesajı bu gündə də qızıl dəyərindədir — net/http-nin handler/mux/context-mexanikasını bilən şəxs hər hansı router/middleware kitabxanasını 1 gündə öyrənir.
2. **Form/PostForm/MultipartForm cədvəli + FormValue gotcha:** bu nüanslar 2026-da da eynidir — köhnə yox, hələ də aktual bilik.
3. **Context awareness (html/template):** XSS müdafiəsinin avtomatikliyi Go-nun unikal qəhrəmanlığıdır — text/template ilə heç vaxt HTML render etmə.
4. **DI = testability:** global Db-dən interfeys + closure quruluşuna keçid test double-lərin açarıdır — bu pattern dəyişməz qalır.
5. **Concurrency = dizayn qərarı:** 1 CPU-da 3.5x qazanc (fan-out/fan-in) və multi-CPU-da pulsuz 10x — kitabın ən yadda qalan performans dərsidir.

## Kitabın ən dəyərli hissəsi

**Chapter 3 (Handling requests) + Chapter 8 (Testing), pages 47-68 + 190-221** — çünki:
- Ch 3 handler/mux/middleware mexanikasını net/http source kodu ilə göstərir — bütün Go web framework-lərinin ortaqlanguage-ı.
- Ch 8 httptest + DI birləşməsi real layihələrdə ən çox istifadə olunan test pattern-idir (mux.ServeHTTP(recorder, request) + FakePost).
- Bu iki chapter birlikdə "Go web developer"in gündəlik 90%-ni əhatə edir.

## Uyğunsuzluq / diqqət qeydləri

⚠️ **Chapter 4, listing 4.15-4.16 (cookie oxuma):** kitab `getCookie` handler-ini iki dəfə fərqli variantda göstərir — birincisi `r.Header["Cookie"]` (raw string slice), ikincisi `r.Cookie/Cookies` metodları. listing 4.16-da ikinci dəfə `fmt.Fprintln(w, c1)` yazılıb, amma `cs` dəyişşəni (`r.Cookies()`) heç çap olunmur — nəzarət üçün faydalıdır, kitabın özündə kiçik redaksiya səhfi ola bilər.

⚠️ **Chapter 7, listing 7.15 (handleRequest):** default case YOXDUR — tanınmayan metod (məs. PATCH) üçün err = nil olaraq sükutla keçir (405 yox). Productionda default: `http.Error(w, "Method not allowed", 405)` əlavə etmək lazımdır.

⚠️ **Chapter 9, `time.Sleep(100 * time.Millisecond)`** thrower/catcher nümunəsində main-i bitirmək üçün istifadə olunur — kitabın özü də qeyd edir ki, bu "hacky"; WaitGroup/channel ilə düzgün shutdown mümkündür.

⚠️ **Chapter 10 bütün kod nümunələri köhnə tooling ilə** (godep, goapp, Upstart, 534MB golang-base Docker image) — 2026-da bunların hamısının müasir ekvivalenti var (modules, gcloud, systemd, multi-stage build). Kitabı öyrənərkən anlayışları götür, alətləri müasir versiyaları ilə əvəz et.
