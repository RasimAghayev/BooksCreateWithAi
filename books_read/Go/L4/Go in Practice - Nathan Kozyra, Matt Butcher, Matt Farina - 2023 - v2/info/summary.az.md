# Go in Practice, Second Edition — Xülasə (Azərbaycanca)

**Müəlliflər:** Nathan Kozyra, Matt Butcher, Matt Farina | **Nəşriyyat:** Manning | **İl:** 2023 | **Səviyyə:** L4 (Advanced)

## Kitabın ümumi məqsədi

"Go in Practice" — praktikaya yönəlmiş, 13 fəsillik kitabdır: hər fəsil real dünya
problemini (Problem/Solution/Discussion formatında) həll edir. Fundamental kitablardan
fərqli olaraq burada məqsəd — Go ilə production-hazır sistemlər (CLI alətləri, serverlər,
mikroservislər, bulud tətbiqləri) qurmağın İDİOMATİK yollarını göstərməkdir. Kitab
4 hissəyə bölünür: fundamentallar → möhkəm tətbiqlər → web tətbiqləri → bulud və
qabaqcıl mövzular.

## Fəsil-fəsil xülasə

1. **Getting started:** Go-nun tarixi/fəlsəfəsi (3 qat: dil+toolchain+ekosistem), çoxlu
   qaytarma (error sonda), standart kitabxana icmalı, goroutine/channel ilk baxış,
   go mod/get, testing ilk baxış, C/Rust/Java/Python müqayisələri, ilk web server.
2. **CLI:** flag paketi (Plan 9 stili), enum çatışmazlığı, Cobra subcommand-ları,
   JSON/YAML/INI konfiqurasiya, env dəyişənləri (12-factor), init daemonlar, graceful
   shutdown (signal.Notify + server.Shutdown), routing — çoxlu handler-lər, Go 1.22
   metod+PathValue, regex router, üçüncü tərəf routerlər.
3. **Struct/interfeys/generics:** anonim struct, funksiya sahəsi vs metod, pointer
   receiver, anonim sahələr, tag-lər (reflection ilə), JSON encode, OOP ilə müqayisə
   (implicit implementasiya), alias tiplər, any/interface{}, union tiplər, generic
   filter/map, ~ təxmini, constraints paketi.
4. **Xətalar:** konvensiya (son qaytarma + zero value), custom tiplər, sentinel
   dəyişənlər, %w wrap + errors.Is, panic vs error fəlsəfəsi, defer/recover, adlı
   return ilə cleanup (OpenCSV), goroutine panic-ləri (handler idiom, safely.Go).
5. **Paralellik:** goroutine, WaitGroup (paralel gzip), loop dəyişəni tələsi, Mutex
   (word counter race), channel-lər (select, time.After), close qaydaları (done
   pattern), buffered channel ilə kilid.
6. **Keyfiyyət:** gofmt/goimports, go vet (context leak, json tag), go mod versiyalar,
   log paketi (flags, custom Writer), slog (1.21, JSONS), stack trace-lər, table testlər,
   fuzzing (1.18), adlı subtestlər, coverage, Delve, benchmark (Bubble Sort).
7. **Fayl/şəbəkə:** ReadFile/Open/Stat/bufio.Scanner, io.Copy, TCP log client (Fatal
   YOX — Panicln!), back pressure, UDP logging trade-off, websocket chat, SSE + fsnotify.
8. **HTTP server:** metod routing (1.22), server timeouts, TimeoutHandler, middleware +
   WithValue, struct handler (DB), PathValue, nested mux, query/form, cookie, JWT, basic
   auth.
9. **Template:** context-aware escaping, range/pipe, FuncMap, parse keşləmə, buffer
   render, nested/inheritance (define/block), template.HTML, email (text/template+smtp).
10. **Data:** FileServer/ServeFile, custom 404, go:embed, CDN konfiq, form/multipart,
    çoxfayllı upload, MIME sniffing (3 səviyyə), MultipartReader stream upload.
11. **Xarici servis:** custom client, timeout detect (type switch), Range ilə davamlı
    download, JSON error, arbitrary JSON, API versiyalama (URL/content-type), gRPC
    (proto3, codegen, client).
12. **Bulud:** mikroservislər (Kafka), REST sürəti (keep-alive transport, dərhal
    Body.Close), sürətli JSON (codecgen/sonic), runtime aşkarlama, cross-compile
    (GOOS/GOARCH/gox/build tag), runtime monitorinqi.
13. **Refleksiya/codegen:** value/type/kind, kind switch, interfeys yoxlaması (nil
    pointer hiyləsi), struct walk, öz ini: tag-ləri (Marshal/Unmarshal), go:generate
    (queue generator).

## Ən vacib 5 fikir

1. **Explicit hər yerdə:** error-lar sonuncu qaytarma, `_, ok :=` yoxlamaları, defer ilə
   resurs zəmanəti — Go-nun "sıxdırıcı"lığı onun etibarlılığıdır.
2. **1.22+ standart router:** metod prefiksli route-lər + PathValue ilə üçüncü tərəf
   mux kitabxanalarının 80%-ni standard əvəz edir.
3. **Concurrency ədəbiyyatı:** loop dəyişəni parametr, close-u sender edir, done kanalı
   ilə bit siqnalı, `-race` mütləq — bunlar race/leak-free kodun açarıdır.
4. **Performans asdıqları:** parse/template init-də; custom Transport-da KeepAlive;
   serial request-lərdə Body dərhal bağlanır; JSON bottleneck — codegen/JIT.
5. **Codegen > reflection:** go:generate dev-time alətidir; generasiya olunmuş kod
   reflection-dan sürətli və sadədir; VCS-ə commit edilir.

## Kitabın ən dəyərli hissəsi

Chapter 8 (HTTP server) + Chapter 12 (bulud) — production Go servisinin tam həyat
dövrünü əhatə edən praktik ikili; Chapter 4 (errors) isə Go-nun ən çox səhv başa
düşülən sahəsinin (panic/recover/defer interaksiyası) dərin təhlilidir.
