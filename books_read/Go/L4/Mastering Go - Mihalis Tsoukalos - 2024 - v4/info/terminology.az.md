# Mastering Go, Fourth Edition — Terminologiya (Azərbaycanca)

Format: **English Term (Azərbaycanca qarşılıq)** — qısa izah

## Giriş və Əsas Tiqlər (Ch1-2)
- **Goroutine** — yüngül icra vahidi, `go` açar sözü ilə
- **Zero Value (sıfır dəyəri)** — təyin edilməmiş dəyişənin default dəyəri
- **Short Assignment Statement (qısa təyinat)** — `:=` ilə tip çıxarılması
- **Rune (simvol kodu)** — int32; tək Unicode code point
- **Slice Header (slice başlığı)** — pointer + len + cap üçlüyü
- **Capacity (tutum)** — allocation-sız böyümə həddi; cap = 2x artım
- **Escape Analysis (çıxış analizi)** — dəyişənin stack/heap qərarı
- **Dereference (dəyərə keçid)** — `*p` ilə pointer hədəfini oxumaq
- **Z-normalization** — (val-mean)/stdDev standartlaşdırması

## Kompozit Tiqlər və Generics (Ch3-4)
- **Comparable (müqayisəolunan)** — == dəstəyi; map açarı üçün şərt
- **Nominal Typing** — myInt və int eyni underlying-lə fərqli tiplər
- **Constraint (məhdudiyyət)** — generic funksiyanın tip siyahısı-interfeysi
- **Union (union)** — `int | float64` tip birləşməsi
- **Supertype (~)** — underlying tipə görə qəbul (`~int`)
- **iota** — const ardıcıllıq generatoru
- **Shallow/Deep Copy (dayaz/dərin nüsxə)** — assignment / rekursiv kopya

## İnterfeyslər və Refleksiya (Ch5)
- **Implicit Interface Satisfaction (bəyansız satisfaksiya)** — duck typing
- **Method Set (metod dəsti)** — tipin metodları toplusu
- **Value/Pointer Receiver** — kopya üzərində / birbaşa dəyər üzərində iş
- **Type Assertion (tip iddiası)** — `i.(T)` runtime yoxlaması
- **Type Switch (tip açarı)** — tip üzrə switch
- **Embedding (daxil etmə)** — anonim sahə; sahə/metod promote
- **reflect.Value/Type/Kind** — dəyər / tip / tip növü üçlüyü

## Paketlər və Sistem Proqramlaşdırma (Ch6-7)
- **Public/Private (böyük/kiçik hərf)** — görünürlük qaydası
- **Blank Import** — `_ \"pkg\"` — yalnız init/side-effect
- **init()** — paket inicializasiyası; main-dən əvvəl, bir dəfə
- **Closure (bağlanma)** — leksik scope dəyişənini capture edən funksiya
- **LIFO defer** — son daxil, ilk icra
- **SemVer** — vMAJOR.MINOR.PATCH modul versiyalması
- **io.Reader/io.Writer** — oxuma/yazma universal interfeysləri
- **Buffered I/O (bufiolu I/O)** — buferli giriş-çıxış; Flush tələb edir
- **Struct Tag (strukt etiketi)** — `json:\"ad\"` metadata; omitempty; `json:\"-\"`
- **go:embed** — binary-ə fayl daxil etmə direktivi
- **slog** — səviyyəli, handler-li strukturlaşdırılmış logging

## Paralellik (Ch8)
- **m:n Scheduling** — m goroutine n thread-də multiplex
- **Work-Stealing (iş oğurlama)** — boş P-nin iş oğurlaması
- **GOMAXPROCS** — aktiv logical P həddi
- **WaitGroup** — Add/Done/Wait sayğaclı sinxronizasiya
- **Race Condition (data yarışı)** — kilidsiz paralel yaddaş çıxışı
- **Critical Section (kritik bölmə)** — mutex-lə qorunan kod
- **Mutex/RWMutex** — eksklüziv / oxu-paralel kilid
- **select** — çoxkanal eyni-anda gözlənti; random seçim
- **Signal Channel (siqnal kanalı)** — yalnız işarə üçün struct{} kanalı
- **Context** — ləğv oluna bilən əməliyyat konteksti; WithCancel/Timeout/Cause
- **Monitor Goroutine** — datanı sahiblənən tək goroutine

## Şəbəkə və Web (Ch9-11)
- **ServeMux / HandlerFunc** — router / funksiya-Handler adapteri
- **PersistentFlags vs Flags** — qlobal vs lokal cobra flag-i
- **ReadTimeout/WriteTimeout/IdleTimeout** — server tərəfli hədlər
- **Multi-Stage Build** — kiçik image üçün ikimərhələli Dockerfile
- **Full-Duplex WebSocket** — tək TCP üzərində ikiyönlü kanal
- **AMQP** — Advanced Message Queuing Protocol (RabbitMQ)
- **REST** — REpresentational State Transfer; HTTP üzərində konvensiya
- **Subrouter** — ümumiləşdirilmiş şərtli nested route
- **Path Variable (yol dəyişəni)** — {id:[0-9]+} kimi route parametri
- **Graceful Shutdown** — siqnal ilə təmiz dayanma

## Test, Fuzz və Observability (Ch12-13)
- **testing.T/F/B** — unit test / fuzz / benchmark tipləri
- **Table-Driven Test** — struct slice üzərində çoxsaylı ssenari
- **Code Coverage (kod örtüyü)** — testlərin toxunduğu kod nisbəti
- **httptest.NewRecorder** — server-siz handler cavab qeydedicisi
- **Seed/Generated Corpus** — fuzz inputlarının mənbələri
- **Profiling (pprof)** — icra ölçmələri; CPU/heap profilləri
- **Observability** — xarici siqnallardan daxili vəziyyət anlayışı
- **Counter/Gauge/Histogram/Summary** — Prometheus metrik tipləri
- **scrape_interval** — Prometheus pull tezliyi

## Performans və GC (Ch14-15, Appendix)
- **b.N** — avtomatik artan benchmark təkrar sayğacı
- **Escape Analysis (`-gcflags -m`)** — heap/stack axınının görünməsi
- **Leaking Parameter** — funksiya qayıtdıqdan sonra parametri yaşadan referans
- **Pre-allocation** — `make([]T, 0, n)` ilə tutum əvvəlcədən
- **Tri-Color Mark-and-Sweep** — white/gray/black GC alqoritmi
- **Write Barrier** — pointer dəyişikliyini GC-yə bildirən kod
- **Mutator** — GC ilə paralel işləyən tətbiq
- **sync.OnceFunc** — tək icra garantili funksiya wrapper-i (Go 1.21)
- **clear** — map silmə / slice sıfırlama built-in (Go 1.21)
- **math/rand/v2** — generic rand.N; rand.Read deprecated (Go 1.22)
