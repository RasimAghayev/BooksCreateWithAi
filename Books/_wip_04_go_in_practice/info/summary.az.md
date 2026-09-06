# Go in Practice — Xülasə (Azərbaycan dilində)

**Müəlliflər:** Matt Butcher, Trevor Matteson (Masterminds)
**Nəşriyyat:** Manning Publications Co., 2016 · **ISBN:** 9781633430075
**Səhifə:** 280 + indeks · **11 chapter · 70 texnika**

## Kitabın ümumi məqsədi

Go-nun real dünya praktikası: 70 konkret texnika PROBLEM→HƏLL→MÜZAKİRƏ formatında. Səthi tutorial YOX — hər texnika production-da qarşılaşılan spesifik problemin (timeout, race condition, vendor lock-in, böyük fayl upload) dəqiq həlli. Tərtib: foundation (CLI/konfiq/router/concurrency) → dayanıqlıq (error/panic/test) → interfeys (template/form/REST) → cloud (multi-provider/microservice/reflection).

## Chapter-by-chapter xülasə

### Ch 1 — Getting into Go (s. 3-25)
Go 3 qat: dil + toolchain (test/fmt/doc/pkg daxili) + ekosistem. Çoxqayıdış + adlı return. Goroutine/channel icmalı. Dil müqayisələri: C (runtime+GC+compile sürəti), Java (tək binary vs JRE+JIT), Python/PHP (GIL vs daxili server), JS (tək thread vs multi-core). GOPATH workspace (src/bin/pkg). Hello Web Server (HandleFunc + ListenAndServe).

### Ch 2 — A solid foundation (T1-9, s. 27-58)
flag paketi (Plan 9 üslubu!) + gnuflag/go-flags (struct tag) + cli.go (komanda/subkomanda framework — Docker istifadə edir). Konfiq: JSON/YAML(gcfg)/INI + 12-factor ENV. Graceful shutdown (manners + signal.Notify); /shutdown antipatterni. Router 4 səviyyə: multi-handler → path.Match wildcard → regex (compile cache!) → 3rd-party (httprouter/mux/pat).

### Ch 3 — Concurrency (T10-15, s. 59-83)
CSP. Gosched (yield — amma zəmanət YOX). WaitGroup + paralel gzip (loop dəyişəni PARAMETR!). Mutex (struct-a embed + defer Unlock) — word counter race nümunəsi + --race detector. Kanallar: multi-channel + select + time.After; **done-kanal close patternı** (close YALNIZ sender — pánik təhlükəsi); bağlı kanal zero value (sonsuz loop). Buffered(1) kanal = kilid.

### Ch 4 — Errors and panics (T16-21, s. 87-111)
Error = sonuncu return; nil-ləri minimallaşdır (işlənən dəyər + error). Error dəyişənləri (ErrTimeout == müqayisə) və custom error tipləri (ParseError line/char). Error vs panic fərqi (gözlənilən vs davam edilməz). panic(errors.New) idiomu. defer + recover + **named return + r.(error)** (panic→error). Goroutine panic öz stack-də qalır → handler-da recover / net/http daxili qoruyur / safely.Go wrapper patternı.

### Ch 5 — Debugging and testing (T22-31, s. 113-143)
Logger (arbitrary io.Writer; Ldate/Lshortfile bitmask). Network logging: TCP back pressure vs UDP itki; Fatal defer-i öldürür → Panicln. Syslog (facility/severity; LOG_LOCAL3). Stack: debug.PrintStack / runtime.Stack(buf, all). Test: mock (interfeys yarat + kodda istifadə), canary test (`var _ io.Writer = &MyWriter{}` — compile-time), generative (quick.Check — truncation bug tapır!), benchmark (b.N kalibrləmə; compile-loop-dan-xaricə 10x), RunParallel + -cpu, --race ilə paylaşılan bufer tutumu.

### Ch 6 — Templates (T32-38, s. 147-166)
html/template context-aware escape (XSS). FuncMap (Funcs→Parse sırası). **Parse cache** (package-level — 10x). Buffer pattern (atomik cavab). Nested ({{template "ad" .}} + ExecuteTemplate). Inheritance (define + block default-ları; map strukturu; həmişə "base" icra). Obyekt→HTML mapping (template.HTML safe tipli ikili render). Email (text/template + smtp).

### Ch 7 — Static və formlar (T39-48, s. 168-193)
FileServer/StripPrefix/ServeContent (304/MIME avtomatik). Custom 404 (go-fileserver). Memory cache (RWMutex + ServeContent). go.rice (asset-ləri binary-yə). CDN per-mühit konfiq. Form: ParseForm/ParseMultipartForm/FormValue; çoxdəyərli sahələr slice ilə. Upload: FormFile (tək) / MultipartForm.File (çox). MIME yoxlaması: DetectContentType (512 bayt öz məzmunundan — extension/header etibarsız). **MultipartReader + NextPart + chunk** — böyük fayl streaming (yaddaş sabit).

### Ch 8 — Web services (T49-55, s. 194-213)
HTTP client (Request/Client ayrılığı; Timeout body-ni əhatə edir). **hasTimedOut** — 3 tip + string fallback. Range header ilə resumable download (rekursiv retry; Accept-Ranges yoxlaması). Custom JSON error ({"error":{code,message}}; json:"-"). Client tərəfi error interfeysi. Arbitrary JSON (interface{} + type switch; float64/number qaydası). API versiya: URL (/v1/) vs vnd. content-type (semantic; Accept negotiation).

### Ch 9 — Cloud (T56-61, s. 217-234)
IaaS/PaaS/SaaS spektri; container vs VM. **Vendor lock-in qarşısı: interfeys + factory** (File interfeysi; LocalFile dev implementasiyası). Divergent error → package error dəyişənləri (orijinalı logla). Host detect (Hostname+LookupHost). Dependency (exec.LookPath). Cross-compile: GOOS/GOARCH; gox paralel; filepath + build tags + foo_windows.go. Runtime monitoring (NumGoroutine — müəllifin goroutine leak story-si; ReadMemStats bahalıdır).

### Ch 10 — Microservice kommunikasiyası (T62-65, s. 235-252)
**Keep-alive**: custom Transport-də Dial KeepAlive qoy; **body-ni vaxtında bağla** (defer tələsi — 2-ci request yeni connection açır). codecgen (reflection-sız JSON). Protobuf (marshalkiçik+sürətli; proto.String pointer-lər). gRPC (proto3 + service + protoc plugins=grpc; context WithCancel/c.Done; HTTP/2). RPC daxili servis-lər üçün, REST xarici API üçün.

### Ch 11 — Reflection və generasiya (T66-70, s. 253-280)
Value/Type/Kind üçlüyü. **Kind switch** (MyInt type-a yox, kind-a uyğundur; ref.Int() hamısını birləşdirir). İnterfeys yoxlaması: nil pointer tricki + Implements(). Struct gəzintisi (Indirect + NumField/Field(i) — tip və value eyni indekslə). Tag.Get("ini") + custom INI marshal/unmarshal (kind switch + Field(i).Set). **go generate**: direktiv + $GOPACKAGE + text/template ilə typed queue generator — generasiya DEV lifecycle-dadır, nəticə VCS-ə commit. Reflection vs generation trade-off.

## Kitabın əsas mesajları

1. **Problem→həll formatı:** hər texnika real production problem-si — "sadəcə API göstər" yox, "back pressure olanda nə et".
2. **Ənənəvi əməliyyat landşaftı:** graceful shutdown, init daemon, syslog, keep-alive — modern cloud-UNIX ənənələri Go kontekstində.
4. **Concurrency idiomları:** loop-parametr tələsi, done-kanal, mutex-embed — köhnə "Goroutine-ə Start ver, sonra bax" yox.
5. **Vendor lock-in:** interfeys + package error + factory — cloud-da da Go-nun əsas dizayn prinsipi.
6. **Reflection son çarədir:** performance + boilerplate → code generation.

## Ən dəyərli 5 texnika (az oxucu üçün)

1. **T14 (channel closing):** close YALNIZ sender; done-kanal signalı — ən çox panic/leak yaradan səhv.
2. **T21 (safely.Go):** goroutine panic stack-də qalır — net/http-in daxili qorunması dərsləri.
3. **T48 (incremental upload):** MultipartReader + chunk — yaddaş sabit, pass-through mümkün.
4. **T56+57 (multi-provider):** interfeys + package error — cloud-da portable Go kodu.
5. **T62 (body close):** `defer r.Body.Close()` çox vaxt PİSDİR — oxuduqdan dərhal sonra bağla; connection reuse bunun üzərindədir.

## Kitabın ən dəyərli hissəsi

**Chapter 3 (Concurrency) + Chapter 4 (Errors/Panics), pages 59-111** — çünki:
- Ch 3 Go-nun kanal semantikasını (close qaydaları, nil/bağlı davranış, done-pattern) kitablar arasında ən dəqiq izah edir — bu, Ch 5-6-da (Get Programming with Go) yarımçıq qalan "kanal bağlandıqda nə olur" sualına tam cavabdır.
- Ch 4 panic/recover + goroutine stack izolasyonunu — runtime daxilənə enmədən praktik dərinliklə açır; safely.Go real library design dərsidir.
- Birlikdə: "dayanıqlı Go" (robust Go) — bu kitabın ən fərqlənən mövzusu.

## Müəlliflər barədə
Matt Butcher — Masterminds (Helm, Glide); Trevor Matteson — Deis. Kitabın texnikaları real layihələrdən: go-fileserver bu kitab üçün yazılıb; Helm generator pattern-dən ilhamlanıb.
