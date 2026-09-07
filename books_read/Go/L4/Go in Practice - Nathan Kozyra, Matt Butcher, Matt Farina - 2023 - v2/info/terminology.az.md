# Go in Practice, Second Edition — Terminologiya (Azərbaycanca)

Format: **English Term (Azərbaycanca qarşılıq)** — qısa izah

## Fundamentallar (Ch1-3)
- **Toolchain (alət dəsti)** — go komandasının daxilindəki build/test/fmt/mod alətləri
- **Naked return (boş qaytarma)** — adlı return-lərsiz dəyərsiz `return`
- **Go Playground** — brauzerdə Go icra/paylaşma mühiti
- **GOOS/GOARCH** — hədəf OS/arxitektura dəyişənləri
- **CSP** — Communicating Sequential Processes; Go paralellik modeli
- **Goroutine** — runtime idarəli yüngül icra vahidi
- **Channel (kanal)** — tipli mesaj boru xətti
- **SemVer (semantik versiya)** — major.minor.patch
- **Command-line framework** — subcommand/dəstəkli CLI aləti (Cobra)
- **12-factor app** — config-environment metodologiyası
- **Graceful shutdown (təmiz dayanma)** — mövcud request-ləri bitirərək sönmə
- **ServeMux / multiplexer** — URL→handler router
- **Precedence (üstünlük)** — route uyğunluq qaydası (longest-match-wins)
- **Anonymous struct (anonim strukt)** — adsız inline tip
- **Receiver (qəbul edici)** — metodun tip bağlılığı `(a Animal)`
- **Embedding (daxil etmə)** — adsız sahə; metod promote-ı
- **Struct tag (strukt etiketi)** — `json:"name"` metadata; reflection ilə oxunur
- **Union type (birləşmə tipi)** — `Cat | Dog` interfeys daxili məhdudiyyət
- **Type approximation (tip təxmini)** — `~int8`: underlying tiplər daxil
- **constraints paketi** — hazır Signed/Integer/Ordered constraint dəstləri

## Xətalar (Ch4)
- **Sentinel error (gözcü xətası)** — paket səviyyəli mənalı error instansiyası
- **Error wrapping** — `%w` ilə zəncir qorunması
- **errors.Is** — wrap zəncirində rekursiv instansiya axtarışı
- **Panic** — davam mümkünsüz halın siqnalı; error-əksinə stack unwind
- **Deferred function (gecikmiş funksiya)** — defer ilə funksiya-sonuna planlanan icra
- **Named return (adlı qaytarma)** — defer-in qaytarma dəyərini dəyişməsi üçün
- **Handler server idiom** — istifadəçi handler-i recover-li wrapper-də icra
- **safely.Go** — panic-tutan goroutine starter

## Paralellik (Ch5)
- **WaitGroup** — Add/Done/Wait sayğaclı gözləmə
- **Race condition (data yarışı)** — kilidsiz paralel yaddaş çıxışı
- **sync.Mutex / RWMutex** — eksklüziv / oxu-paralel kilidlər
- **select** — çoxkanal gözləyici; RANDOM hazır-case seçimi
- **time.After** — müddət sonra siqnal verən kanal
- **Done channel** — "bit" siqnalı; close-u sender-ə həvalə edir
- **Back pressure (geri təzyiq)** — istehsal > emal sıxışması

## Keyfiyyət (Ch6)
- **gofmt/gofmt -d** — kanonik format / diff nümayişi
- **goimports** — gofmt + import avtomatikləşdirmə
- **gopls** — Language Server Protocol Go implementasiyası
- **go vet** — compile-keçən səhvlərin statik analizi
- **slog (1.21)** — struktur log paketi
- **JSONS** — newline-ayrı JSON mesaj axını
- **Table-driven test** — input/expected struct cədvəli
- **Fuzzing (1.18)** — random input ilə invariant yoxlama
- **Coverage (örtük)** — testlərin toxunduğu kod faizi
- **Benchmark (b.N)** — təkrar saylaşmalı performans ölçməsi
- **Delve** — cəmiyyət debugger-i
- **Stack trace (iz dəsti)** — runtime.Stack / debug.PrintStack

## Şəbəkə/Web (Ch7-11)
- **TCP slow-start** — bağlantı başında congestion ramp-up
- **Keep-alive** — bağlantı reuse (HTTP protokolu vs TCP OS səviyyəsi)
- **Websocket / wss://** — davamlı ikiyönlü kanal / TLS formaları
- **SSE / EventSource** — server→client biryönlü axın
- **Flusher** — dərhal göndərmə interfeysi
- **fsnotify** — fayl dəyişikliyi watcher-ı
- **PathValue (1.22)** — `{id}` route dəyişəninin oxunması
- **Middleware** — `func(http.HandlerFunc) http.HandlerFunc` sarğı
- **JWT (JSON Web Token)** — imzalı identitet daşıyıcısı
- **BasicAuth** — header əsaslı istifadəçi/parol
- **Context-aware escaping** — html/template-in çıxışa görə avtomatik escape
- **template.HTML** — "safe HTML" marker tipi
- **FuncMap** — şablon daxili funksiya qeydiyyatı
- **define/block** — adlı şablon / default-məzmunlu icra
- **go:embed** — compile-time binary daxil fayl direktivi
- **MultipartReader / NextPart** — stream upload emalı
- **DetectContentType** — 512-bayt MIME sniffing
- **Range header** — bayt aralığı ilə davamlı download
- **vnd. content type** — versiyalı vendor media type
- **proto3 / protoc** — protobuf dili / kod generatoru

## Bulud/Advanced (Ch12-13)
- **Mikroservis** — tək-məsuliyyətli kiçik servis
- **Message queue (Kafka)** — asinxron iş ötürmə növbəsi
- **High availability** — fərdi-servis HA strategiyası
- **exec.LookPath** — PATH-də asılılıq yoxlaması
- **gox** — paralel multi-platform build
- **Build tag / fayl suffiksi** — `//go:build !windows` / `foo_windows.go`
- **runtime.MemStats / NumGoroutine** — runtime sağlamlıq göstəriciləri
- **Sidecar goroutine** — monitorinq üçün fondda işləyən gözləyici
- **reflect.Value/Type/Kind** — refleksiya üçlüsü
- **Kind switch** — quruluş ailəsinə görə birləşmiş switch
- **reflect.Indirect** — pointer-samiiyətə açma (pointer deyilsə no-op)
- **go:generate** — dev-time kod generasiyası direktivi
- **Codegen vs reflection** — compile-time vs runtime tip emalı
