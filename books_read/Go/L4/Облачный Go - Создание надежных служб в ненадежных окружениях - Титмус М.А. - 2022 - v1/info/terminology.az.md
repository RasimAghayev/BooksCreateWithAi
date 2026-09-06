# Облачный Go / Создание надежных служб в ненадежных окружениях — Terminologiya (Azərbaycanca)

> Kitab boyu rast gəlinən terminlərin Azərbaycanca izahları. Format: `English Term (Azərbaycanca qarşılıq)`.

## A

**Adapter (Adaptor)** — Heksagonal Arxitekturada portun xarici interfeysdə (REST, gRPC, CLI) həyata keçirilməsi.

**Append-only log (Yalnız əlavə qeydiyyat)** — məlumatların dəyişdirilməməsi, yalnız sonuna əlavə edilməsi — dayanıqlılıq (durability) üçün.

**Application state (Tətbiq vəziyyəti)** — müəyyən server instance-na bağlı olan sorğu konteksti, sessiya məlumatları. Stateless dizaynda istifadə edilmir.

## B

**Backoff (Geri çəkilmə)** — uğursuz cəhddən sonra gözləmə müddətinin artması; eksponensial artım + jitter.

**Bulkhead (Gertməli bölmə)** — xidmətin müxtəlif resursunu izolyasiya edən şablon — bir xidmət düşərsə, digərlərinin resursları təsir olmaz.

## C

**Cascading failure (Kaskad xəta)** — bir komponent düşdükdə ondan asılı digər komponentlər də düşür — domino effekti.

**Circuit Breaker (Devir açıcı)** — xidmət xətasını aşkar edərək müvəqqəti olaraq sorğuları dayandıran şablon; bağlı/ açıq vəziyyətləri var.

**Cloud Native (Bulud Native)** — CNCF tərifi: bulud texnologiyaları ilə miqyaslana bilən, dinamik tətbiqlər yaradmaq və işə salmaq imkanı.

**CNCF (Cloud Native Computing Foundation)** — bulud native texnologiyalarını inkişaf etdirən və təşviq edən təşkilat.

**Composition over Inheritance (Kompozisiya irsiyyət əvəzinə)** — Go-da klassik OOP inheritance (miras) yoxdur; əvəzinə embedding və interface-lər ilə kompozisiya.

**Context** — Go 1.7+ paketi; sorğu həyat dövrünü idarə etmək: deadline, timeout, cancel, value ötürmə.

**Context.Background()** — əsas kontekst, heç vaxt ləğv edilmir; `main`, testlər, inisializasiya funksiyaları üçün.

**Context.WithCancel** — manual ləğv edilən kontekst yaradan funksiya.

**Context.WithTimeout** — vaxt aşımı ilə avtomatik ləğv edilən kontekst.

**Context.WithValue** — kontekstə açar/dəyər cütü əlavə edən funksiya (spor pratikada istifadə edilir).

**CSP (Communicating Sequential Processes)** — Tony Hoare tərəfindən təklif olunmuş konkurensiya modeli; Go goroutine + channel əsasında.

## D

**Dapr (Distributed Application Runtime)** — bulud native tətbiqlər üçün runtime, sidecar arxitekturası.

**Debounce (Anti-flash)** — funksiyanın çağırma tezliyini məhdudlaşdıran şablon; DebounceFirst (ilk çağırış) və DebounceLast (son çağırış) variantları.

**Docker** — konteyner texnologiyası; Go statik binary-ləri ilə mümkün olan minimal image-lər.

## F

**Fallacies of Distributed Computing (Paylanmış hesablamanın səhv fərzləri)** — Deutsch tərəfindən 1991-ci ildə formulə edilmiş 8 səhv fərz: şəbəkə etibarlıdır, gecikmə sıfırdır, və s.

**Feature flags (Xüsusiyyət flaqları)** — funksionallığı deploy etmədən açıb-bağlamaq imkanı verən mexanizm; statik və ya dinamik (per-request) ola bilər.

**Functional options pattern** — `Option func(*options) error` + `WithPort` closure-ları ilə optional konfiqurasiya.

## G

**Goroutine (Qorutin)** — Go dilində asinxron işləyən yüngül funksiya; OS thread-dən daha az yaddaş istehlak edir (~2 KB stack).

**gRPC** — Google tərəfindən hazırlanmış RPC framework; `.proto` kontraktları, HTTP/2, binary serialization (Protocol Buffers).

## H

**Health check (Sağlamlıq yoxlaması)** — xidmətin vəziyyətini izləmək üçün endpoint: liveness (canliliq), readiness (hazırlıq), deep (asılılıqlar).

**Horizontal scaling (Horizontal miqyaslama)** — yeni server instance-ları əlavə etmək, load balancer arxasında; limitsiz, amma idarəsi daha çətindir.

## I

**Immutable infrastructure (Dəyişməz infrastruktur)** — serverlərə "mal" (cattle) kimi yanaşma; problem olduqda yenidən qurub əvəz etmək, yamaq etmək deyil.

**Instrumentation (Instrumentasiya)** — kodda izləmə (tracing), metrikalar, loglar əlavə etmək mexanizmi.

## J

**Jitter (Qarışıq)** — eksponensial artıma təsadüfi miqdar əlavə edərək kütləvi təkrar cəhdləri (thundering herd) qarşısı alan texnika.

## L

**Loose coupling (Zəif əlaqəlilik)** — komponentlərin minimal asılılıqla qurulması; bir komponent dəyişsə, digəri təsir olmadan işləməyə davam edir.

**Load shedding (Yük atma)** — xidmət öz-özünü qoruyaraq həddi aşan sorğuları rədd edən mexanizm.

## M

**Manageability (İdarəetmə)** — sistemin davranışını kod dəyişdirmədən dəyişdirmə imkanı.

**Microservices (Mikroservislər)** — kiçik, müstəqil deploy edilə bilən xidmətlər, yüngül mexanizmlər (HTTP, gRPC) ilə əlaqələnir.

**Monolith (Monolit)** — bütün funksionallıq tək prosesdə, tək kod bazasında.

## O

**Observability (Gözlənilirlik)** — sistemin daxili vəziyyətini xarici müşahidələrdən (logs, metrics, traces) anlama imkanı.

**OpenTelemetry** — vendor-neytral gözlənilirlik standartı; traces, metrics, logs üçün birlikdə API.

## P

**Pets vs Cattle** — server idarəetmə fəlsəfəsi: "ev heyvanı" (pet, yamaq tələb edən) əvəzinə "mal" (cattle, istənilən vaxt əvəz edilə bilən).

**Plugin (Plugin)** — tətbiyin funksionallığını yenidən kompilyasiya etmədən dinamik olaraq genişləndirmək; Go `plugin` paketi və ya HashiCorp `go-plugin`.

**Port (Port)** — Heksagonal Arxitekturada core tərəfindən təyin edilən interfeys (müqavilə).

**Prometheus** — açıq mənbəli metrika toplama və izləmə sistemi; time-series DB, PromQL sorğu dili.

**Pub/Sub (Yayın/Abunə)** — asinxron mesajlaşma modeli; istehsalçı (publisher) və istehlakçı (consumer) eyni anda olmasına ehtiyac yoxdur.

## R

**Rate limiting (Sürət məhdudlaşdırma)** — müəyyən vaxt intervalında qəbul ediləcək sorğu sayını məhdudlaşdıran mexanizm; token bucket, leaky bucket.

**Readiness probe (Hazırlıq yoxlaması)** — xidmətin sorğuları qəbul edə bilmə vəziyyəti; orchestrator yükü çəkir və ya götürür.

**Redundancy (Redundans)** — bir neçə xidmət nümunəsi; biri düşsə, digərləri davam edir.

**Reliability (Dayanıqlılıq)** — sistemin gözlənilən vəziyyətdə funksiyasını yerinə yetirməsi, əsas məqsədinə uyğun olaraq.

**Resilience (Uyğun gəlmə)** — xətalar vəziyyətində də funksiyanı yerinə yetirmə qabiliyyəti.

**Retry (Təkrar cəhd)** — uğursuz əməliyyatları təkrar edən şablon; eksponensial artım + jitter ilə birlikdə.

**REST (Representational State Transfer)** — uniform interfeys üzərində HTTP sorğuları; JSON/XML, GET/POST/PUT/DELETE.

## S

**Scalability (Miqyaslana bilənlik)** — sistemin artan yükə uyğunlaşaraq performansını saxlayabilməsi.

**Service discovery (Xidmət kəşfi)** — mikroservis arxitekturasında xidmətlərin bir-birini tapması mexanizmi (Consul, etcd, K8s Service).

**Sharding (Paylama)** — məlumatları hissələrə bölmək; vertical sharding (hash əsaslı), horizontal sharding.

**Sleeper / Thundering herd (Qorxulu döngü / Göy qurd)** — kütləvi təkrar cəhdlər və ya restart-lar nəticəsində yükün kütləvi artması.

**Stateful (Vəziyyətli)** — serverə bağlı vəziyyət saxlayan; miqyaslanma çətinliyi yaradır.

**Stateless (Vəziyyətsiz)** — serverə bağlı vəziyyət saxlamayan; hər sorğu öz data ilə gəlir.

## T

**TCP keep-alive** — boş bağlantıların aradan qaldırılması mexanizmi.

**Termination (Dayandırma)** — xidmətin təmiz şəkildə bağlanması; SIGTERM qəbul edərək resursları təmizləmək.

**Throttle (Drossel)** — funksiyanın çağırma tezliyini sabit sürətlə məhdudlaşdıran şablon; Debounce'dan fərqli olaraq, packet-ləri deyil, sürəti məhdudlaşdırır.

**Throughput (Kəsmə qabiliyyəti)** — vahid vaxtda emal edilən sorğu/element sayı.

**Time coupling (Vaxt əlaqəliliyi)** — komponentlərin eyni anda və ya müəyyən sıra ilə sorğu göndərməsi ehtiyacı; asinxron modellərlə azalır.

**Twelve-Factor App** — SaaS və bulud tətbiqləri üçün 12 ilk metodologiyası; codebase, dependencies, config, backing services, build/run, processes, port-binding, concurrency, disposability, dev/prod parity, logs, admin processes.

## V

**Vertical scaling (Vertikal miqyaslama)** — eyni serverdə RAM/CPU əlavə etmək; sadə, amma limitli.

**Vertical sharding (Vertikal paylama)** — assosiativ massivi açarın hash-i ilə N alt-massivə bölmək — kilidləşməni azaltmaq.

**Viper** — Go konfiqurasiya kitabxanası (spf13/viper); JSON, YAML, TOML, env vars, CLI flags, uzaq konfiqurasiya.

## Z

**Zap** — Go üçün yüksək performanslı strukturlaşdırılmış jurnalçı kitabxanası (uber); `zap.Logger` (maksimum performans) və `zap.SugaredLogger` (rahatlıq).
