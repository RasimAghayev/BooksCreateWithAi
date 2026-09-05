# Go in Action — Müəllim qeydləri (Teacher Mode)

📖 Kitab deyir:
Go-nun 4 sütunu — sürətli kompilyasiya, built-in concurrency, hierarşiyasız tip sistemi, GC. Concurrency üçün goroutine + channel (CSP modeli). Tip sisteminde composition > inheritance; interfeyslər kiçik və davranış-yönümlü. Receiver seçimi tipin təbiətinə görə. Race condition-lər atomic/mutex/channel ilə həll olunur. Standart kitabxana geniş istifadə üçün nəzərdə tutulub.

👨‍🏫 Müəllim qeydi:
Məncə bunu praktikada belə tətbiq etmək daha yaxşı olar:

1. **Kitab GOPATH dövründə yazılıb — bugünkü Go-dan çox fərqli sabitliyin var.** Chapter 3-dəki `GOPATH`, vendoring, godep, gb mövzuları tarixi maraqdir. Bugün: `go mod init` + `go.mod` + `go.sum` + `go mod vendor` + `go mod tidy`. GOPATH qurmağa heç vaxt ehtiyacın yoxdur. Amma paket = qovluq, import konvensiyaları, `go vet`/`go fmt`/`go doc` — bunlar hamısı hələ də canlı qaydalardır.

2. **Concurrency-də kitabın atomic/mutex nümunələrini real layihədə son seçim kimi gör.** İdiomatik Go: (a) "share memory by communicating" — data-nı bir goroutine sahiblənsin, digərləri channel ilə soruşsun; (b) sadə sayğac üçün `atomic` OK; (c) sadə əlaqələndirilmə üçün `sync.WaitGroup`, `errgroup.Group` (kitabdan sonrakı əlavə) daha praktikdir. Kitabın Runner pattern-i `context.Context`-dən əvvəl dövrdür — bugünkü proqramlarda timeout/cancel üçün `context.WithTimeout` standartdır.

3. **Kitabın log package yanaşması hələ də işləsə də, real sistemlərdə struktur edilmiş logging seç.** `log.New` + MultiWriter + Discard pattern-i əsas olaraq dərsləri ilə qiymətlidir, amma productionda səviyyə-filtirlənən (zap, zerolog, slog) logger-lər — JSON çıxış, dinamik səviyyə, sahə əsaslı kontekst — daha yaxşı observablelik verir. Go 1.21+ isə `log/slog` standart kitabxanadadır — başlanğıc üçün mükəmməl seçimdir.

4. **Chapter 9-un test konvensiyaları möhkəmdir, amma 2020-ci illərdə əlavələr var:** kitabın "Should ..." mesaj üslubu gözəldir, amma assert funksiyaları üçün `github.com/stretchr/testify` (assert/require) və ya Go 1.9+ üçün `t.Helper()` + xüsusi assert-lər daha qısa testlər yazmağa imkan verir. `testify`-də `require` = kitabdakı `t.Fatal`, `assert` = `t.Error` fərqi ilə eyni məntiqdir.

5. **`ioutil` paketi (kitabda geniş istifadə olunur) deprecated-dır.** `ioutil.Discard` → `io.Discard`; `ioutil.ReadFile` → `os.ReadFile`; `ioutil.WriteFile` → `os.WriteFile`; `ioutil.ReadAll` → `io.ReadAll`. Köhnə kodu oxuyarkən bunu bilmək lazımdır.

6. **`err == ErrTimeout` müqayisəsi (Runner pattern) müasir Go-da `errors.Is`-lə əvəz olunmalıdır.** `switch err { case runner.ErrTimeout: ... }` paketdəki error dəyişəni ilə eynilik müqayisəsidir — wrapped error-lər üçün `errors.Is(err, ErrTimeout)` lazımdır. Chapter 2-dəki `log.Fatal` istifadəsi də diqqət tələb edir — main-də OK, library-də qadağan.

7. **Reader/Writer fəlsəfəsi kitabın ən dəyərli hissəsidir.** Bu interfeysləri öz tiplərində implement etmək (`Read`/`Write`) — io ekosisteminə qoşulmağın açarıdır. `io.Copy`-un daxilindəki 32KB buffer + `io.ReaderFrom`/`io.WriterTo` optimallaşdırmaları — məntiq bir bütün kimi davam edir.

## Ən vacib 5 fikir

1. **Race condition-ləri gözləmə, aşkarla:** `go test -race` və `go build -race` — kitabın sadə counter nümunəsi göstərir ki, "işləyir" görünən kod səhv nəticə verə bilər (4 əvəzinə 2).
2. **Slice-in length/capacity ayrılığı Go-nun ən çox bug yaradan incəliyidir:** `s[i:j:k]` detach pattern kəşf edilməlidir — append birbaş digər slice-ın datasını dəyişə bilər.
3. **Method sets qaydası — value vs pointer receiver yalnız metod daxilində dəyişikliklə bağlı deyil:** interfeysə yerləşdirmə qabiliyyətini də müəyyən edir. Bu, chapter 5-in ürəyidir.
4. **Unbuffered channel = qarantiyalı handoff, buffered channel = müstəqil iş.** Runner/Pool/Work pattern-lərində seçim qəsdən idi: hovuz üçün buffered (resurs sıraya düşür), işçi üçün unbuffered (itən iş yoxdur).
5. **Standard library bir dərslikdir:** `strings.Trim`-in value əsaslı imzası, `os.File`-ın pointer receiver qaydası, `log.New`-un Writer interfeysi — bunlar dilin fəlsəfəsinin living kod nümunələridir.

## Kitabın ən dəyərli hissəsi

**Chapter 5 (Go's type system) + Chapter 6 (Concurrency), pages 88-157** — çünki:
- Chapter 5 receiver/təbiət/method sets üçlüyünü standart kitabxananın öz kodundan (`time.Time`, `os.File`, `net.IP`) göstərərək izah edir — bu, başqa heç bir Go kitabında belə dərin deyil.
- Chapter 6 scheduler-in daxili mexanikasını (logical processor, blocking syscall davranışı, network poller) göstərən az sayda mənbədən biridir — "goroutine ucuzdur" iddiasının **nə üçün** doğru olduğunu izah edir.

Bu iki chapter birlikdə Go developer-in "idiomatik düşünməsini" formalaşdırır: tipin təbiətinə görə dizayn + data paylaşmaq əvəzinə ötürmək.

## Uyğunsuzluq / diqqət qeydləri

⚠️ **Chapter 7, work.go (listing 7.28):** kitabın özündə `Pool` struct `work` sahəsi ilə elan olunur, amma `New` funksiyasının constructor literalında `tasks: make(chan Worker)` yazılıb (səhifə 178, listing 7.30: `p := Pool{tasks: make(chan Worker)}` və `close(p.tasks)`). Struct-da sahə `work` adlanır, amma `New` və `Shutdown`-da `tasks` istifadə olunur — kitabda çap səhfi ola bilər; doğru halda hamısı eyni addır (`work`). Diqqətli olun.

⚠️ **Chapter 3, gb/godep:** bugünkü Go (1.21+) üçün bu alətlər arxaikdir; yalnız `go mod` ecosystem-i öyrənin.

⚠️ **Chapter 9, listing 9.26 (ExampleSendJSON):** kitabın özündə `r` ilə request yaradılır, amma decode `w.Body`-dən oxunur (`json.NewDecoder(w.Body)`) — `w` dəyişəni heç yerde elan edilməyib; düzgün versiya `rw.Body`-dir. Çap səhfi.

⚠️ **Go Playground linkləri** (http://play.golang.org/p/...) — köhnə subdomaindir, bugün https://go.dev/play. Kitabdaki http:// URL-lər bugün https://go.dev/p/pkg üslubuna redirect olunur.
