# Everyday Golang — Terminoloji Lüğət (Azərbaycanca)

## B

**Benchmark** — go test -bench; b.N dövründə ns/op ölçüsü.

**Buffered channel** — alıcısız qəbul edən kanal (make(chan T, n)).

**Buildx** — çox-platformalı Docker build aləti (QEMU ilə cross-compile).

## C

**CGO** — Go↔C; cross-compile üçün CGO_ENABLED=0 ilə söndür.

**Collector (Prometheus)** — Describe/Collect; dinamik label dəsti ixracı.

**Coverage (statement)** — icra olunan if %; rəqəm davranışı ƏVƏZ ETMİR.

**Cross-compile** — GOOS/GOARCH env-i ilə başqa platformaya build.

## D

**Distroless** — SDK-sız, minimal, non-root baza image.

## E

**embed.FS** — build-ə qoşulmuş fayl sistemi (tələbə görə yüklənir).

**errgroup** — paralel funksiyalardan xəta toplayan qrup.

## F

**FaaS/Serverless** — idarə olunan funksiyalar; funksiya = ixtisaslaşmış
mikro servis.

**faas-cli up** — build + push + deploy bir əmrdə.

## G

**Goroutine** — runtime idarəli yüngül axın (OS thread 1:1 DEYİL).

**GOPATH** — pkg (module cache) + bin (install binariləri) + src (köhnə).

## H

**httptest** — NewRecorder (cavab qeydi) + NewRequest (fake sorğu);
server-siz handler testi.

**HMAC** — simmetrik açarla mesaj doğrulaması (webhook digest).

## I

**Implied interface** — konkret tip interfeysi bilmir → fake qoşmaq asan.

**init()** — paket/funksiya yüklənməsindən ƏVVƏL icra (OpenFaaS-də DB
bağlantısı burada).

## L

**Loop-variable capture** — closure dəyişənin SON dəyərini görür; həll:
parametr və ya j := i.

**ldflags -X** — link zamanı string dəyişən injeksiyası (version/commit).

## M

**mergo** — strukturları birləşdirən kitabxana (WithOverride).

**Middleware** — handler-ı bürüyen funksiya; zəncirlənir.

**mux (gorilla)** — regex path parametr + Methods() verb ayrımı.

## O

**OpenFaaS** — portable open-source serverless platforma; faasd = tək-host
variant.

## P

**Pseudo-version** — tag-siz modul: v0.0.0-tarix-SHA.

## R

**RED metrikaları** — Rate (counter) / Errors (code) / Duration (histogram).

**RWMutex** — N paralel oxu + 1 eksklüziv yazı.

## S

**Secret** — konfidansial data; OpenFaaS-də /var/openfaas/secrets/NAME faylı.

**Shrinkwrap** — Docker-sız lokal build variantı (sürətli iterasiya).

**singleflight** — eyni açarlı paralel zənglərin dedupe-i (shared=true).

**Stack.yml** — OpenFaaS funksiya tərif faylı (lang/handler/image/secrets/env).

**Statically-linked** — runtime tələbsiz tək binari.

## T

**Template (text/template)** — {{.Field}} tokenləri; custom funksiya = tip
metodu; Must; statik yoxlanma YOX → test.

**Test-table** — ssenari slice + t.Run; paraleldə tc := c kopyası.

**Test-double/Fake** — asılılığın saxta implementasiyası.

## U

**Unbuffered channel** — göndəriş ancaq hazır alıcı ilə.

**User-Agent** — HTTP identifikasiya başlığı (etiket mütləqdir).

## V

**Vendor** — asılılıqların layihə daxilində kopyası; go mod vendor.

## W

**WaitGroup** — Add/Done/Wait; goroutine-lərin bitməsini gözləmə.

**Worker pool** — N worker + iş kanalı + close() — konkurrentliyi LIMITLƏ.

## X

**X-Call-Id / X-Duration-Seconds** — OpenFaaS çağırış izləmə/ölçü başlıqları.
