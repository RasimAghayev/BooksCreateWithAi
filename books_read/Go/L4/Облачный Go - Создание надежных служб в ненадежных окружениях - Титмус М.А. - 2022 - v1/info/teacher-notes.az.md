# Облачный Go / Создание надежных служб в ненадежных окружениях — Müəllim Qeydləri (Teacher Notes)

> Bu sənəd kitabı öyrədərkən istifadə üçün metodiki qeydləri, çətin anlaşılan yerləri, müzakirə suallarını və praktik tapşırıqları birləşdirir.

---

## 1. Auditoriya və ön şərtlər

- **Hədəf auditoriya:** 6+ ay Go təcrübəsi olan, prodakşnda işləyən developer-lər. Chapter 3-dək Go sintaksis əsasları məcburi, chapter 4-dən etibarən cloud native patterns tələb olunur.
- **Ön şərtlər:** Go sintaksisi (ch3), goroutine/channel əsasları, HTTP əsasları, Docker əsasları, Git.
- **Kitabın özəlliyi:** Chapter 3 praktikdir (kod nümunələri var), chapter 4-11 konseptual + praktik şablonlar və real tətbiq qurulumu ilə. Chapter 5-də tam bir key/value xidməti qurulur — bu, chapter 4 şablonlarının birgə tətbiqidir.

---

## 2. Tədris axını (12 chapter üzrə plan)

| Chapter | Mövzu | Konsentrasiya | Süret |
|---------|-------|---------------|--------|
| 1 | Cloud Native atributları | Konseptual | 5% |
| 2 | Go və cloud | Konseptual | 5% |
| 3 | Go əsasları | Praktik | 15% |
| 4 | Cloud patterns + Context | Patterns + kod | 20% |
| 5 | Key/value xidməti qurulumu | Tam tətbiq | 15% |
| 6 | Reliability (12-factor, immutabl) | Konseptual | 5% |
| 7 | Scalability (stateless, leak) | Arxitektura + kod | 5% |
| 8 | Loose coupling (REST, gRPC, plugins) | Protokollar + arxitektura | 10% |
| 9 | Resilience (retry, circuit, bulkhead) | Patterns + kod | 10% |
| 10 | Manageability (config, flags, Cobra) | Praktik + kod | 5% |
| 11 | Observability (logs, metrics, traces) | OpenTelemetry + kod | 5% |

**Tədris strategiyası:**
- Chapter 1-2: Müzakirə əsaslı, kod yox. Cloud native konsepsiyasını qruplara bölüb müzakirə edin.
- Chapter 3: Live coding — hər konstruksiyanı real kodla göstərin. Tələbələr eyni anda kod yazsın.
- Chapter 4: Hər şablon üçün: (1) məqsəd, (2) kod nümunəsi, (3) tələbələrin özü kod yazması. Circuit Breaker və Debounce ən vacibdir.
- Chapter 5: Step-by-step pair programming — key/value xidmətini 2-3 tələbə birgə qursun.
- Chapter 6: Müzakirə əsaslı — Twelve-Factor App hər bir ilki üçün real kompaniya təcrübəsi bölün.
- Chapter 7-9: Araloq (interleave) — teoriya + praktik nümunə.
- Chapter 10-11: Live demos — real xidmətə konfiqurasiya və gözlənilirlik əlavə edin.

---

## 3. Çətin anlaşılan anlar və izah üsulları

### a) Context paketi və ləğv yayılması (Ch4)
**İzah aracı:** Diaqram — istemci → xidmət A → xidmət B. `ctx.Done()` kanalı bağlananda bütün törəmə kontekstlər sinxronlaşdırılmış ləğv sinxalı alır. `defer cancel()` unutmaq resource leak-ə səbəb olur — demo edin.

### b) Circuit Breaker vəziyyət maşını (Ch4)
**Sual:** "Niyə 2 saniyə?" — bu kitabdakı sadə implementasiya, real sistemlərdə (Hystrix, Resilience4j) daha mürəkkəb backoff alqoritmləri var. State diagram çəkin: CLOSED → OPEN (n failures) → HALF-OPEN (timeout) → CLOSED/OPEN.

### c) Gorutin sızıntısı (Ch7)
**Demo:** `time.Ticker` yaradın, `Stop()` çağırmadan gorutini bitirin — goroutine dump ilə gorutin hələ də işlədiyini göstərin. `defer ticker.Stop()` qaydasını vurğulayın.

### d) Heksagonal Arxitektura — nəyə görə? (Ch8)
**İzah aracı:** "Kod test edilə bilməlidir" prinsipi. Core məntiq DB, HTTP, gRPC-dan asılı olmamalıdır. Adapter-lər mock edilə bilir. Diaqram: Core (portlar) → REST Adapter, gRPC Adapter, PostgreSQL Adapter.

### e) Retry storm — real hadisə (Ch9)
**Nümunə:** 2011-ci ildə Amazon DynamoDB hadisəsi — bir uğursuzluq bütün client-lərin təkrar cəhdlərini törətdi, kaskad xəta yaratdı. Jitter olmadan eksponensial artım da kifayət deyil — müəyyən vaxtda hamı eyni anda təkrar cəhd edirsə, problem davam edir.

### f) Health check növləri (Ch9)
**Sual:** "Liveness nə üçün lazımdır? Ready dəyil?" — Liveness prosesin canlılığını yoxlayır (həyat dövrü dəyərləri: deadlock, infinite loop). Ready isə xidmətin sorğuları qəbul edə bilmə vəziyyətini yoxlayır (asılılıqlar hazırdırmı?).

---

## 4. Müzakirə sualları

**Ch1-3:**
1. Cloud native və cloud hosted arasındakı fərq nədir? (hosting = xidmət etdirmə, native = dizayn prinsipi)
2. Go-nu seçməyin ən güclü arqumenti nədir? (kompilasiya sürəti + yaddaş təhlükəsizliyi + konkurensiya + sabit API)
3. `:=` və `=` arasındakı fərq nədir? Kölgələnmə (shadowing) nəyə səbəb olur?

**Ch4-5:**
4. Context `Value` metodu nəyə görə mübahisəlidir? (compile-time type təhlükəsizliyini pozur)
5. DebounceFirst və DebounceLast nə vaxt istifadə edilir? (UI axtarışı vs server-side rate limiting)
6. Key/value xidmətində əməliyyat qeydiyyatı nə üçün lazımdır? (dayanıqlılıq — restart sonrası bərpa)

**Ch6-8:**
7. Immutabl infrastruktur və CD/CI arasındakı əlaqə nədir?
8. Stateful və stateless arasındakı fərq real tətbiqdə necə görünür? (sessiya idarəetməsi, cache, DB)
9. Heksagonal Arxitektura REST və gRPC üçün eyni core kodunu necə istifadə edir? (port interfeysi, adapter injection)

**Ch9-11:**
10. Retry storm nəyə səbəb olur və jitter onu necə qarşısı alır?
11. Bulkhead nə vaxt lazımdır? (Bəzi xidmətlərin resurslarını digərlərindən izolyasiya etmək)
12. OpenTelemetry nəyə görə Prometheus və Zipkin-dən daha yaxşıdır? (vendor-neytral, birlikdə traces + metrics + logs)

---

## 5. Praktik tapşırıqlar

**Tapşırıq 1 (Ch4):** Circuit Breaker və Debounce şablonlarını real xidmət üçün yazın. `Breaker(Debounce(myFunc))` şəklində birləşdirin. Test üçün artificial latency (time.Sleep) əlavə edin.

**Tapşırıq 2 (Ch5):** Key/value xidmətini Docker konteynerə yerləşdirin. Multi-stage build ilə ~2 MB image yaradın. `docker run` ilə işə salın və API test edin.

**Tapşırıq 3 (Ch7):** Stateless REST API qurun. Hər sorğu öz data ilə gəlsin, heç bir server-side sessiya saxlanmasın. Nümunə: `/api/v1/users/{id}` — hər instance eyni cavabı qaytarsın.

**Tapşırıq 4 (Ch8):** Heksagonal arxitektura ilə sadə key/value xidməti qurun: Core (interface), PostgreSQL adapter, REST adapter. Hər adapterı ayrıca test edin.

**Tapşırıq 5 (Ch9):** Retry + Circuit Breaker + Bulkhead birlikdə işləyən xidmət qurun. `errgroup` ilə paralel sorğular göndərin, kaskad xətaları simulyasiya edin.

**Tapşırıq 6 (Ch10):** Cobra + Viper ilə CLI xidməti qurun. Konfiqurasiya faylı, env vars, CLI flags iyerarxiyası. Feature flag middleware əlavə edin.

**Tapşırıq 7 (Ch11):** OpenTelemetry + Zap ilə xidmətinizi gözlənilir edin. `/metrics` (Prometheus) və `/trace` endpoint-ləri əlavə edin. Goroutine leak-i `pprof` ilə tapın.

---

## 6. Sınav sualları (nümunə)

**Asan:**
1. Go context paketinin 4 əsas funksiyası nədir? (Background, WithTimeout, WithCancel, WithValue)
2. Twelve-Factor App-in 1 və 6-cı illəri nə deməkdir? (Codebase, Processes)
3. Vertical və horizontal scaling arasındakı fərq?

**Orta:**
4. Circuit Breaker 2 vəziyyəti nədir? (Closed, Open) HALF-OPEN nə vaxtdir?
5. Heksagonal Arxitekturanın 3 əsas komponenti? (Core, Ports, Adapters)
6. Goroutine sızıntısı nədir? Necə qarşısı alınır?

**Çətin:**
7. Context.Value nəyə görə mübahisəlidir? Nəyə görə istifadə edilməməsi tövsiyə olunur?
8. Eksponensial artım + jitter olmadan nə baş verir? Retry storm nədir?
9. OpenTelemetry-in 3 sütunu nədir? Nəyə görə vendor-neytral standart vacibdir?

---

## 7. Kollektiv layihə ideyası

**"Bulud Native Go Key/Value Xidməti"** — tələbələr kitabın chapter 4-5-11 şablonlarını birlikdə tətbiq edərək tam xidmət qururlar:
- Chapter 4: Circuit Breaker, Retry, Debounce
- Chapter 5: REST API, əməliyyat qeydiyyatı, Docker
- Chapter 11: Structured logging (Zap) + health checks + Prometheus metrics
- Qiymətləndirmə: funksionallıq + testləmə + Docker image ölçüsü + observability

---

## 8. Əlavə resurslar

- Kitabın rəsmi reposu: github.com/cloud-native-go (nümunə kodlar)
- CNCF Landscape: cncf.io
- OpenTelemetry Go: opentelemetry.io/docs/instrumentation/go
- Prometheus: prometheus.io/docs/guides/go-application
- Zap: go.uber.org/zap
- Cobra/Viper: spf13.com/pages/cobra, spf13.com/pages/viper
- HashiCorp go-plugin: github.com/hashicorp/go-plugin
- "Designing Distributed Systems" (Brendan Burns, O'Reilly)
- "Release It!" (Michael T. Nygard)
