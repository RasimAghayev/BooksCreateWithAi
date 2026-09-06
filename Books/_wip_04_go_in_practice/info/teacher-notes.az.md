# Go in Practice — Müəllim qeydləri (Teacher Mode)

📖 Kitab deyir:
70 texnika ilə real dünya Go: CLI-dan cloud-a qədər. Concurrency kanal qaydaları, panic/recover idiomları, graceful shutdown, connection reuse, protobuf/gRPC, reflection və code generation — hər biri problem→həll formatında.

👨‍🏫 Müəllim qeydi:
Məncə bunu praktikada belə tətbiq etmək daha yaxşı olar:

1. **Kitab 2016 (Go 1.5-1.6 dövrü) — bir çox texnika arxaikləşib.** Ən vacib dəyişikliklər:
   - **T5 (manners):** `http.Server.Shutdown(ctx)` Go 1.8-dən standartdır — manners paketinə ehtiyac YOX. `signal.NotifyContext` (Go 1.16+) sinyal idarəsini daha da sadələşdirir.
   - **T10 (Gosched):** Go 1.14-dən scheduler asinxron preemption edir — Gosched-in roli praktiki olaraq yox oldu.
   - **Generics (Go 1.18+):** T70 (typed queue generator) artıq `type Queue[T any]` ilə 3 sətirdir. Generator pattern hələ də valid (protobuf, stringer kimi alətlər üçün), amma "collections üçün" arqumenti keçib.
   - **T3/T4 konfiq:** viper (spf13) — JSON/YAML/ENV/flag bir paketdə; CLI üçün cobra standartdır (urfave/cli-nin əsas rəqibi).
   - **T63 (codecgen):** ugorji codec var amma JSON üçün müasir seçim `json-iterator` və ya `sonic`-dur; encoding/json özü də Go 1.21-dən ciddi optimallaşdırılıb.
   - **gRPC:** protoc plugin-ləri → `protoc-gen-go-grpc`; `required` proto3-də YOX; context artıq standart kitabxanada (Go 1.7+, kitabda təxmin edirdi).

2. **T62 (body close) — kitabın ən zamansız texnikası.** `defer r.Body.Close()` vs vaxtında close fərqi bu gündə də DÜZGÜNDÜR — amma HTTP/2 multiplexing-də təsiri azalır. Müasir idiom: defer + ReadAll mümkün qədər erkən; streaming hallarında vaxtında close.

3. **T48 (MultipartReader streaming) — hələ də qızıl bilik.** Böyük fayl upload-da ParseMultipartForm tmp faylları doldurur — bu, cloud-da real problem olaraq qalır. Müasir əlavə: `http.MaxBytesReader` ilə request həcm limiti.

4. **Ch 9 (cloud):** go-cloud (gocloud.dev) kitabın T56 interfeys vizyonunun GOOGLE tərəfindən standartlaşdırılmasıdır — `blob.Bucket` abstraksiyası məhz File interfeysi konseptidir. Tunnsiyə oxucuya çatırmağa dəyər.

5. **Testing:** testify (assert/require/mock) — canary test-lərin 90%-ni əvəz edir; həm də Ch 5-6-dakı mock pattern hələ də DI (dependency injection) üçün doğrudur.

6. **Kitabın güclü tərəfi — "gotcha" payı:** hər texnikada "bunu səhv etsən nə olar" var (Fatal defer-i öldürür; bağlı kanal zero value; defer body close; MyInt type vs kind). Bu, müasir kitablarda az rast gəlinir — kitabı hələ də dəyərli edən elə bu "tələlər toplusu"dur.

## Ən vacib 5 texnika (bugün üçün)

1. **T49+T50 (timeout + Range):** hasTimedOut funksiyası hazır bugünkü kontekstlə də işlək; resumable download HTTP/1.1 Range ilə hələ standartdır.
2. **T52 (custom JSON error):** `{"error":{...}}` envelope — bugünkü API-lərin (RFC 7807 problem+json-ə qədər) əsasını kitab 2016-da göstərirdi.
3. **T56+57 (interfeys + package error):** portability prinsipi — go-cloud-un sübuta çevirdiyi dizayn.
4. **T21 (safely.Go):** goroutine panic izolasiyası — net/http daxilində görülən pattern; müəlliflərin "yaddaşa bel bağlama" mesajı bugün də aktual.
5. **T68+69 (struct walk + tags):** əvvəlcədən JSON marshal şeması görmək — kitabın T69 custom ini codec nümunəsi bu gün sqlc/envconfig kimi tool-ların eyni prinsipidir.

## Kitabın ən dəyərli hissəsi

**Chapter 3 + 4, pages 59-111** — kanalların tam semantikası (close, nil, done-pattern) + panic/recover idiomları + goroutine stack izolasiyası. Bu iki chapter "səhv edənlər üçün Go"nu izah edir: kitabın digar fərqi — yalnız "necə yazmaq" yox, "necə POZULUR".

## Uyğunsuzluq / diqqət qeydləri

⚠️ **Ch 10, gRPC import:** kitabda `"golang.org/x/net/context"` — bugün `"context"` (standart). Kitabın özü də Go 1.7-yə idarə edirdi — bugün birbaşa standart kitabxana.

⚠️ **T15 (buffered channel lock):** texniki olaraq düzgün, amma müasir Go-da kanal-kilid yerinə sync.Mutex/RWMutex üstünlük tövsiyə olunur — kanallar kommunikasiya üçün, kilidlər qoruma üçün (go proverb).

⚠️ **T62 listing 8.1 (Go sənədlərindən köhnə nümunə):** kitabın əsas dərslərindən biri — öz sənədlərindəki köhnə nümunəni tənqid edir (keep-alive-siz Transport). Bu, bugünkü Go sənədlərində artıq düzəldilib — amma köhnə blog/SO cavablarında hələ də rast gəlinir.

⚠️ **Elephant: `redis.Pool` məsələn** — kitab Redis/DB kitabxanalarına toxunmur; müasir oxucu üçün əlavə: database/sql-in own connection pool var; Redis üçün go-redis v9 müasir seçimdir.

⚠️ **Ch 9 (IaaS API-lər):** Go Cloud SDK-ları (aws-sdk-go-v2, google.golang.org/api) bu gün mükəmməl şəkildə müstəqildir — kitabın interfeys abstraksiyası hələ də ən yaxşı praktikadır, sadəcə implementasiya daha asandı.
