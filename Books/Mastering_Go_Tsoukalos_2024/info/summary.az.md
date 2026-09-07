# Mastering Go, Fourth Edition — Xülasə (Azərbaycanca)

**Müəllif:** Mihalis Tsoukalos | **Nəşriyyat:** Packt | **İl:** 2024 | **Səviyyə:** L4 (Advanced)

## Kitabın ümumi məqsədi

"Mastering Go" 4-cü nəşri — amatör/orta səviyyəli Go proqramçılarını ADVANCED səviyyəyə
çəkən, 15 fəsil + GC əlavəsindən ibarət hərtərəfli kitabdır. Xüsusiyyəti: kitab boyu
TƏK layihə — statistika tətbiqi — fəsil-fəsil təkamül edir (CLI → CSV → JSON → web
service → REST → cobra CLI → concurrent). Hər fəsil Go-nun bir qatını açır: tiplərdən
başlayır, generics/reflection/interfeyslərlə abstraksiya qurur, package arxitekturası
və sistem proqramlaşdırmasından keçərək concurrency, web, TCP/IP, REST, test/profiling,
fuzz/observability, performans və Go 1.21-1.22 yeniliklərinə qədər aparır.

## Fəsil-fəsil xülasə

1. **Giriş:** Go tarixi/fəlsəfəsi, which(1), log, statistika app v1.
2. **Əsas tiplər:** error, numeric, string/rune, slices (header/cap), pointers, unsafe,
   crypto/rand.
3. **Kompozit tiplər:** maps, structs, regexp, CSV.
4. **Generics:** constraints, ~supertype, cmp/slices/maps paketləri.
5. **Refleksiya+interfeyslər:** reflect, type methods, implicit satisfaction, sort.Interface,
   empty interface, error interfeysi, OOP mimikriyası, generics vs reflection.
6. **Paketlər:** görünürlük, defer LIFO+closure tələsi, init sırası, SQLite paketi (GitHub),
   go doc, workspaces, ldflags versiyalama.
7. **Sistem proqramlaşdırma:** io.Reader/Writer, buffered I/O, JSON tag-lər, viper, cobra,
   go:embed, ReadDir/DirEntry, io/fs, slog.
8. **Paralellik:** scheduler (m:n, work-stealing), WaitGroup, channel-lər, race+`-race`,
   select, worker pool, signal, Mutex/RWMutex/atomic, context, semaphore, kanalsız
   concurrent stats.
9. **Web servis:** net/http, ServeMux, handler-lər, Docker multi-stage, errgroup, timeout-lar.
10. **TCP/UDP/WebSocket:** net paketi, concurrent TCP patterni, gorilla websocket, RabbitMQ.
11. **REST:** status kodları, gorilla/mux subrouters, SQLite auth serveri, cobra REST
    client, API versiyalama.
12. **Test/profiling:** run() patterni, pprof, trace, httptest, table-driven, coverage,
    go vet, govulncheck, cross-compilation, go:generate, example funksiyaları.
13. **Fuzz+observability:** testing.F, corpus, UTF-8 bug, runtime/metrics, expvar,
    Prometheus+Grafana stacki.
14. **Performans:** benchmark (b.N, benchmem), buffer Reset optimizasiyası, escape
    analysis, slice/map leak-lər, eBPF.
15. **Go 1.21/1.22:** sync.OnceFunc, clear, loop var semantikası, math/rand/v2, ServeMux
    metod+wildcard.
16. **Əlavə (GC):** tri-color alqoritmi, write barrier, MemStats, map vs slice GC yükü.

## Ən vacib 5 fikir

1. **Go-nun sadəliyi qərarlı dizayndır:** 25 açar söz, bir-yol prinsipləri, implicit
   interfeyslər — böyük komandalarda saxlanıla bilənlik üçün.
2. **Concurrency dizayn qərarıdır, parallelism mükafatdır:** WaitGroup/unikal-index kimi
   kanalsız həllər deadlock-suz sadə paralellik verir; kanal hər problemin cavabı deyil.
3. **Yaddaş idarəetməsi GC-yə buraxılıb, amma strukturu SƏN seçirsən:** slice-of-structs
   pointer-li map-dən 40x sürətli; pre-allocation; leak (subslice, map bucket) — Go-da
   mümkündür.
4. **Etibarlılıq üçün toxum testlər:** unit + table-driven + fuzz + `-race` + coverage +
   govulncheck — hamısı bir komplementer dəstədir.
5. **Production Go = observability ilə:** slog/prometheus/pprof — işləyən sistemin
   görünməz tərəfini görünür edən alətlər toplusu.

## Kitabın ən dəyərli hissəsi

Chapter 8 (Concurrency) + Appendix (GC) ikilisi — Go-nun runtime kimliyinin (scheduler,
channel, tri-color GC) ən dərin və praktik izahı; Chapter 12-14 üçlüyü isə test-performans
mədəniyyətini tam dövrə ilə öyrədir.
