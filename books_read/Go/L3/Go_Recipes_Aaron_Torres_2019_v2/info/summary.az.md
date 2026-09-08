# Go. Рецепты программирования (2nd ed) — Xülasə

**Müəllif:** Аарон Торрес (Aaron Torres) | **İl:** 2019 | **Nəşriyyat:** Packt
**ISBN:** 978-1-78980-098-2 | **Səhifə:** 464 | **Fəsil:** 14 (85+ resept)

## Kitabın mahiyyəti

"Go Programming Cookbook" — problem-həlli formatında praktik reseptlər
toplusu. Hər resept: hazırlıq → "Necə edilir" (addım-addım kod) → "Necə
işləyir" (izah). Kitab standart kitabxananın idarəedilməz görünən
genişliyini strukturlaşdırır: I/O interfeyslərindən başlayır, mikroservis,
test, concurrency, distributiv sistemlər, serverless və performance
optimallaşdırmasına qədər aparır.

## Fəsil-fəsil xülasə

| # | Fəsil | Mövzu |
|---|---|---|
| 1 | I/O və fayl sistemləri | io.Reader/Writer, bytes/strings, CSV, temp fayllar, şablonlar (text/html template) |
| 2 | CLI alətləri | flag paketi, subcommand, env + envconfig, TOML/YAML/JSON konfiq, Unix pipe, siqnallar, ANSI |
| 3 | Data transformasiyası | tip çevirmə, math/big, valyuta (float64 problemi), NullTypes, gob/base64, reflection, closure kolleksiyalar |
| 4 | Xəta emalı | error interfeysi, pkg/errors wrap/cause, log, logrus/apex strukturlu logging, context logging, sync.Once, panic recover |
| 5 | Şəbəkə | TCP/UDP socket, DNS lookup, WebSocket (gorilla), net/rpc, net/mail |
| 6 | Bazalar | database/sql, DB/Tx interfeysi, pul + timeout, Redis, MongoDB, Storage interfeysi |
| 7 | Veb klientlər | http.Client/Transport, REST client, async sorğular, OAuth2 + token saxlanc, decorator middleware, gRPC, twirp |
| 8 | Mikroservislər | handler/ResponseWriter, controller DI, validasiya, content negotiation, middleware + context, reverse proxy, gRPC→JSON |
| 9 | Testlər | closure mock + Patch/Restore, gomock/mockgen, table-driven coverage, gocov/goconvey, BDD (godog) |
| 10 | Concurrency | channel/select, WaitGroup, mutex/atomic/Once, context, state kanalları, worker pool, pipeline |
| 11 | Distributiv sistemlər | Consul discovery, Raft konsensus (FSM), Docker/ldflags, Compose, Prometheus, go-metrics |
| 12 | Reaktiv | Goflow dataflow graph, Kafka (sync/async Sarama), Kafka→Goflow, GraphQL server |
| 13 | Serverless | AWS Lambda + Apex, apex logs/metrics, App Engine + Datastore, Firebase Firestore |
| 14 | Performans | pprof profiling, benchmark (RWMutex vs atomic), allocasiya analizi (concat vs Join), fasthttp |

## Kitabın əsas mesajları

1. **İnterfeyslər hər yerdə:** io.Reader/Writer, RoundTripper, Storage,
   Client — funksiyalar konkret tip deyil, davranış qəbul etsin. Bu test
   (mock) və dəyişkənliyi (portability) mümkün edir.

2. **Composition over inheritance:** struct embedding (*http.Client,
   *oauth2.Config) və closure-lar ilə funksional kompozisiya (Decorator,
   Map/Filter).

3. **Xəta ötürmə protokolu:** Wrap ilə kontekst artır, Cause ilə orijinala
   qayıt; log yalnız son təyinatında.

4. **Concurrency primitivləri əl ilə idarə olunur:** kanal, select, WaitGroup,
   mutex/atomic, sync.Once — hər birin öz yeri; worker pool/pipeline ilə
   paralelliyi məhdudlaşdır.

5. **Observability dəsti:** strukturlu logging + Prometheus metric + pprof
   profiling — produksiya Go servisinin 3 ayağı.

6. **Proqramlaşdırma ən yaxşıları sınaqla doğrulanır:** benchmark-lar
   (atomic vs mutex, Join vs concat) mikro-optimallaşdırma qərarlarına
   rəqəmsal əsas verir.

## Auditoriya

- Əsas Go sintaksisini bilən, real dünya pattern-ləri axtaran developerlər
- Standart kitabxananı dərindən öyrənmək istəyənlər
- Go-da mikroservis/backend texnologiya dəsti qurmaq istəyənlər

## Nə öyrənilir

- Standart kitabxana istifadəsi (io, os, net, encoding/*, context, sync)
- 3rd-party ekosistem (gorilla, gomock, apex/log, Sarama, gRPC, graphql-go)
- Deploy və işlətmə (Docker, Compose, Lambda, App Engine, Firebase)
- Performans alətləri (pprof, benchmark, benchmem)
