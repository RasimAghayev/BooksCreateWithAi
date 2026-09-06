# 100 Go Mistakes and How to Avoid Them — Müəllim Qeydləri (Teacher Notes)

> Bu sənəd kitabı öyrədərkən istifadə üçün metodiki qeydləri, çətin anlaşılan yerləri, müzakirə suallarını və praktik tapşırıqları birləşdirir.

---

## 1. Auditoriya və ön şərtlər

- **Hədəf auditoriya:** 6+ ay Go təcrübəsi olan, prodakşnda işləyən developer-lər. Junior-lar üçün ch1-2 uyğun; ch8-12 üçün konkurrentlik əsasları tələb olunur.
- **Ön şərtlər:** Go sintaksisi, goroutine/channel əsasları, HTTP əsasları, Git.
- **Kitabın özəlliyi:** 100 müstəqil səhv — **modul sistemi elastikdir**: istənilən səhv ayrıca 15-30 dəq dərs kimi istifadə oluna bilər. Ardıcıllıq tövsiyə olunur, amma məcburi deyil.

---

## 2. Tədris axını (12 module üzrə plan)

| Module | Chapter | Səhvlər | Konsentrasiya | Süret |
|--------|---------|---------|---------------|--------|
| 1 | Ch1 | — | Fəlsəfə + 7 kateqoriya | 5% |
| 2 | Ch2 | #1-#16 | Kod təşkilatı + abstraksiya | 13% |
| 3 | Ch3 | #17-#29 | Yaddaş modeli (slice/map) | 13% |
| 4 | Ch4-5 | #30-#41 | Loop + string məkanikası | 10% |
| 5 | Ch6-7 | #42-#54 | Funksiya imzaları + xətalar | 12% |
| 6 | Ch8 | #55-#60 | Konkurrentlik əsasları | 12% |
| 7 | Ch9 | #61-#74 | Konkurrentlik praktikası | 15% |
| 8 | Ch10 | #75-#81 | Standart kitabxana | 8% |
| 9 | Ch11 | #82-#90 | Test intizamı | 7% |
| 10 | Ch12 | #91-#100 | Aşağı-səviyyə optimizasiya | 5% |

**Tək layihə YOXDUR** — hər səhv mikto-ssenaridir. Dərsdə hər səhv üçün: (1) başlıqı gizlədib "bu kod nə çap edər?" sualı verin; (2) tələbi müzakirə edin; (3) həlli açın; (4) çoxhallı həllərin trade-off-unu müzakirə edin.

---

## 3. Çətin anlaşılan anlar və izah üsulları

### a) Slice-in pointer təbiəti (Ch3 — 8 səhvin kökü)
**İzah aracı:** Dialoqda whiteboard-da slice-in 3 hissəsini çəkin (ptr+len+cap). Hər səhv üçün ƏVVELƏ memory diagram, SONRA kod:
- #20: append len<cap → eyni array
- #25: s2 append s1-in elementini dəyişir
- #26: msg[:5] 1MB array-i yaşadır
- #28: map bucket-ləri silinsə də qalır

### b) Nil pointer → non-nil interface (Ch6 #45)
**Sual:** "Bu kod niyə `customer is invalid: <nil>` çap edir?" — interface = (tip, dəyər) cütü; (T, nil) ≠ (nil, nil). Wrapper/wrappee diaqramı şərtdir. `var foo *Foo; foo.Bar()` nümunəsi ilə nil receiver-in qanuniliyini göstərin.

### c) Memory model — unbuffered receive < send (Ch8 #58)
**Əks-intuitivdir.** Oxu yazıdan əvvəl bərabərdir?! Sıralama ZƏMANƏTİNİ (deyil causality-ni) göstərdiyini vurğulayın: send tamamlanması receive-in bitməsini gözləyir → yazı mütləq əvvəl baş vermiş olmalı. Buffered-da send dərhal tamamlanır → zəmanət yoxdur. 2 diaqram yan-yanə.

### d) Select random seçimi (Ch9 #64)
**Sual:** "case sırası əhəmiyyətli deyilsə, niyə source-da yazılı sıra var?" — starvation qorunması. Tələbi canlı göstərin: 10 mesajın 5-i itir. Inner select + default həllini addım-addım icra edin.

### e) Escape analysis (Ch12 #95)
**Demo:** `-gcflags "-m=2"` canlı icra — tələbələr kompilyator qərarını ÖZ GÖZLƏRİ ilə görür. sumValue vs sumPtr benchmark (1.26 vs 14.84 ns) — pointer qaytarmağın qiyməti.

### f) False sharing (Ch12 #92)
**Ən çətin mövzu.** MESI-ni qısaca təqdim edin: cache coherency üçün yazma = digər core-lərdə line invalidasiyası. "Sharing memory is an illusion" sitatını yazın. Padding həllünü benchmark nəticəsi ilə (40%) təqdim edin.

### g) CFS throttling (Ch12 #100)
**Hesabladırtdırın:** 8 thread × 50ms = 400ms quota → 100ms periodda 50ms THROTTLE. Tələbələr 300% latency cəzasını öz hesablasınlar. automaxprocs = 1 sətirlik həll — amma NİYƏ işlədiyini anlamak şərtdir.

---

## 4. Müzakirə sualları

**Ch2-3:**
1. `wg.Add` goroutine DAXİLİNDƏ olsa nə baş verir? (71 — nondeterminizm sübutu)
2. Niyə `stringset.New` `util.NewStringSet`-dən yaxşıdır? Amma nə vaxt utility qəbulediləndir?
3. `[]string{}` ilə `var s []string` — JSON klienti strict-dirsə hansı?

**Ch6-7:**
4. Pointer receiver "shallow" mutasiya edə bilməzmi? (`customer.data.balance` nümunəsi — #42)
5. `%w` ilə `%v` arasındakı "coupling" konkret nə deməkdir? İnterfeys dəyişsə nə qırılır?
6. Niyə `logging + return` "ikidəfə idarəetmə" sayılır — amma 2 müxtəlif SİSTM-də log haqqıqdır? (mübahisə mövzusu — server-side vs client-side idarəetmə)

**Ch8-9:**
7. Merge sort threshold 2048 — bu rəqəm başqa maşında eyni qalacaqmı? Magic dəyərlərlə mübarizə strategiyaları?
8. sync.Cond vs channel broadcast: nə vaxt Cond MÜTLƏQ lazımdır? (təkrar broadcast + çoxlu receiver)
9. Detach context (#61) "həmişə düzgün" həlldirmi? Nə vaxt ləğv olunma İSTƏNİLƏN davranışdır?

**Ch10-11:**
10. `resp.Body.Close()` oxumadan vs oxuyub — TCP bağlantıya fərqi nədir? (keep-alive reuse)
11. `errors.As` 2-ci arqumentində pointer olmasa nə olur? (compile OK, run-time panic — #50)
12. Table-driven test-də `tt := tt` yalnız paralel subtestlərdəmi lazımdır? Niyə?

**Ch12:**
13. `make(1GB)` virtual heap hilesi — bu "hack"dirmi və ya qanuni texnika? Risks?
14. Struct sahə sırası dəyişmək davranışı dəyişirmi? (Yox — sadəcə yaddaş/layout)
15. Critical stride 512 = 4KB — bu dəyər HƏR cache-də eynidirmi? (64 set × 64B — cache parametrlərinə bağlı)

---

## 5. Praktik tapşırıqlar

**Tapşırıq 1 (#1-#2):** Kölgələnmə yaradan kod yazın; `vet`+`shadow` ilə aşkarlayın; hər 2 üsulla düzəldin. Happy-path-left refaktorunu edin.

**Tapşırıq 2 (#20-#26):** Slice memory əyləncəsi: (a) `s1`, `s2=s1[1:3]` append → s1 dəyişsin; (b) 1MB slice-dan 5 bayt saxla → `runtime.ReadMemStats` ilə leak göstər; (c) kopya ilə müqayisə.

**Tapşırıq 3 (#42-#47):** `MultiError` nil-receiver tələsini yazın; `-race` ilə yoxlayın; hər 2 həlli tətbiq edin. `io.Reader` qəbul edən `countLines` + `strings.NewReader` testi.

**Tapşırıq 4 (#50-#52):** transientError + %w wrap → type switch-in break olmasını reproduce edin; `errors.As` ilə düzəldin. İkiqat-log variantını tək-log+%w-a refactor edin.

**Tapşırıq 5 (#55-#56):** Sequential merge sort → V1 paralel (8x yavaşlanmanı görün) → threshold-lu V2 (benchmark ilə optimal threshold tapın).

**Tapşırıq 6 (#61-#66):** Kafka publish + HTTP handler detach context; merge2channel nil-channel pattern; notification chan struct{}; inner select + default prioritet.

**Tapşırıq 7 (#71-#73):** WaitGroup səhvini yazın (nondeterminizm) → 2 düzəliş; sync.Cond donation goals; errgroup paralel API calls.

**Tapşırıq 8 (#75-#81):** Loop-da time.After leak → NewTimer+Reset; connection pool 4 parametr konfiqurasiya; production HTTP client (4 timeout); `-race` + body close audit.

**Tapşırıq 9 (#82-#90):** Table-driven + `tt := tt` paralel; `-shuffle=on` ilə gizli asılılıq tapma; httptest.NewServer ilə client test; benchmark-da local→global + yeni-data pattern.

**Tapşırıq 10 (#91-#95):** struct-of-slices vs slice-of-structs benchmark; alignment 24→16 bayt; `-gcflags -m=2` ilə escape analizi; sumValue vs sumPtr.

---

## 6. Sınav sualları (nümunə)

**Asan:**
1. `wg.Add` harada çağrılmalıdır? (Valideyndə, spin-dən əvvəl)
2. `len("hêllo")` = ? (6 — bayt)
3. Loop-da `time.After` nə problem yaradır? (Timer resursları bitənədək yaddaşda)

**Orta:**
4. `s[:2]` vs `s[:2:2]` — append təsirində fərq? (İkincisi cap-ı 2 edir → kənar append yeni array yaradır)
5. Nil pointer qaytaran error funksiyası niyə həmişə non-nil error verir?
6. Buffered vs unbuffered channel sinxronizasiya fərqi? (Unbuffered = receive send-dən əvvəl zəmanət)

**Çətin:**
7. Niyə `for range` loop-da channel dəyişmək range-i etkilemez? (Ifadə 1 dəfə kopyalanır)
8. 512 vs 513 kolon matris fərqinin səbəbi? (Critical stride: 64 set × 64B = 4KB = 512 int64)
9. GOMAXPROCS=8, K8s limit=4 core: nə baş verir? (CFS: 8×50ms=400ms→50ms throttle/period)

---

## 7. Kollektiv layihə ideyası

**"Production-Grade Checklist Generator"** — tələbələr 100 səhvdən praktik audit-checklist qurur:
- HTTP servisi üçün: default client/server qadağası, timeout-lar, body close, return after http.Error
- Konkurrentlik üçün: -race CI-da, goroutine exit planları, sync kopya audit, errgroup
- Yaddaş üçün: preallocation, nil-slice return, clone substring, map lifecycle
- Hər qrup öz domain-i üçün 15-maddəlik checklist + kod nümunələri təqdim edir

Qiymətləndirmə: hər maddə üçün (a) səhvin #nomresi, (b) necə aşkarlanır (linter/test/profiler), (c) necə düzəldilir.

---

## 8. Əlavə resurslar

- Kitabın rəsmi reposu: github.com/teivah/100-go-mistakes (kod nümunələri)
- Race detector: go.dev/doc/articles/race_detector
- Memory model: go.dev/ref/mem
- Effective Go: go.dev/doc/effective_go
- pprof: go.dev/blog/pprof + go.dev/doc/diagnostics
- Diagnostics: go.dev/blog/diagnostics
- automaxprocs: github.com/uber-go/automaxprocs
- Go issue 33803 (CFS-aware GOMAXPROCS)
- benchstat: golang.org/x/perf/cmd/benchstat
- "Understanding Real-World Concurrency Bugs in Go" (ASPLOS 2019)
