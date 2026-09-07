# Everyday Golang — Xülasə (Azərbaycanca)

**Müəllif:** Alex Ellis | **İl:** 2021 (Go 1.17-yə uyğun) | **Səviyyə:** L2 (Elementary)

## Kitabın ümumi məqsədi

OpenFaaS qurucusu Alex Ellis-in 5-6 illik Go təcrübəsindən yığılmış
GÜNDƏLİK pattern kolleksiyası: kitab nəzəriyyə deyil, PRODUKSİYA alətlərindən
(OpenFaaS CLI, arkade, k3sup, inlets) çıxarılmış real nümunələrlə gedir.
18 fəsil Go lifecycle-ni əhatə edir: quraşdırma → modules → ilk proqram →
CLI → testlər → konkurrentlik → DB → YAML konfiq → version injeksiyası →
embed → templates → HTTP server → Prometheus → release pipeline →
Docker → OpenFaaS serverless.

## Mövzu bloklarının xülasəsi

**Təməllər (1-5):** GOPATH/modules/vendoring keçidi; HMAC webhook
doğrulaması; flag parsing (TrimSpace validasiyası); çoxpaketli layihələr
(export = böyük hərf); cross-compile (GOOS/GOARCH, -ldflags strip, CGO=0);
HTTP+JSON (timeout, User-Agent etikası, funksiyaya refaktor); CLI-də
flags→Cobra pilləliyi və 5 usability prinsipi.

**Testlər+Konkurrentlik+DB (6-8):** assertion-siz want/got fəlsəfəsi;
test-table; coverage/benchmark/parallel/stress party tricks; testlərin
binariya düşməməsi; implied interfaces ilə asılılıq izolyasiyası; httptest
Recorder/NewRequest; goroutine loop-capture tələsi; errgroup; singleflight
(thundering herd: 330→33); Mutex/RWMutex; channels+select; context timeout;
worker pool (6.2s→2.3s, resurs limiti); lib/pq + database/sql todo app.

**Konfiq+HTTP+Metrik+Release (9-17):** YAML struct tag-ləri; mergo
birləşdirmə; -ldflags -X version injeksiyası; embed (string/[]byte/FS +
http.FileServer); text/template (custom metodlar, Must, test məcburiyyəti);
http.Server timeouts; gorilla/mux (regex path, Methods); middleware zənciri;
Prometheus RED (CounterVec/HistogramVec code+method), custom Gauge
(inflight), Collector interfeysi (dinamik label Reset+doldur); Makefile dist
5 platform; GitHub Actions tag→publish; multi-stage+multi-arch Docker
(distroless/nonroot, buildx+QEMU, GHCR).

**OpenFaaS (18):** funksiya = ixtisaslaşmış mikro servis; lock-in azadlığı
vs cloud functions; golang-middleware şablonu; faas-cli up axını; bcrypt-fn;
handler testləri build-də icra; --shrinkwrap lokal iterasiya; secrets
(fayl) vs env; init()-də DB bağlantısı (sorqular arası açıq); REST routing
POST/GET+path.

## Kitabın əsas mesajları

1. **Ən sadə işləyən şeylə başla:** standart flag → Cobra; kanallar
   İŞDÜZGÜN anlaşılanda; assertion YOX (konsistenslik Go-nun dəyəridir).
2. **Produksiya adətləri tez öyrənilir:** timeout, User-Agent, defer Close,
   version injeksiyası, non-root image — hamısı "gündəlik" detal deyil,
   MÜTLƏQ detallardır.
3. **Concurrency pilləli:** WaitGroup → errgroup → channels → worker pool →
   singleflight — hər alətin öz problemi.
4. **Pipeline avtomatlaşdır:** tag push → test+build+upload; buildx
   multi-arch — paylanma artıq ƏL işi deyil.

## Kim üçündür

Go-ya yeni başlayan və ya başqa dildən gələn developer-lər üçün ideal
"fast track" — hər nümunə əyləncəli real alət (astronaut sayacı, ISS
pozisiyası, todo, bcrypt) üzərində. L2: dərin nəzəriyyə YOX, amma hər
pattern-in "niyə"-si izah olunur. Arundel seriyasından fərqli: daha
praktik, less fəlsəfi.
