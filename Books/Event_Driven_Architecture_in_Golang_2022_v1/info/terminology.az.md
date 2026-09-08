# Event-Driven Architecture in Golang — Terminologiya (AZ)

| Termen (EN) | Azərbaycanca qarşılıq | İzah |
|---|---|---|
| Event | Hadisə | Keçmişdə baş vermiş dəyişməz fakt; adı past tense |
| Event Notification | Hadisə bildirişi | Yalnız ID daşıyan yüngül hadisə |
| Event-Carried State Transfer | State daşıyan hadisə | Detallı məlumat daxil; consumer lokal nüsxə saxlayır |
| Event Sourcing | Hadisə mənbəyi | Cari vəziyyət = hadisələrin replay-i |
| Event Stream | Hadisə axını | Broker-də retention + replay (Kafka, JetStream) |
| Producer / Consumer | İstehsalçı / İstehlakçı | Hadisə dərc edən / qəbul edən |
| Message Broker | Mesaj vasitəçisi | Aradakı növbə (NATS, Kafka) |
| Loose Coupling | Boş bağlılıq | Komponentlər bir-birini bilmir |
| Eventual Consistency | Sonrakı uyğunluq | Vəziyyət dərhal yox, axın sonunda uyğunlaşır |
| Dual Write | İkiqat yazma | DB + broker ayrı yazma problemi |
| Domain Event | Doman hadisəsi | Bounded context daxili hadisə |
| Integration Event | İnteqrasiya hadisəsi | Modullararası, broker-dən keçən |
| Bounded Context | Sərhədli kontekst | DDD: modelin yaşadığı sərhəd |
| Ubiquitous Language | Ümumvahid dil | Kontekst daxili ortaq terminologiya |
| Aggregate | Aqreqat | Ardıcıllıq vahidi olan kök model |
| Hexagonal Architecture | Altıbucaqlı arxitektura | Ports & Adapters |
| CQRS | Komanda-sorğu ayrılığı | Yazma və oxu modellərinin bölünməsi |
| Projection / Read Model | Proyeksiya / Oxu modeli | Hadisələrdən qurulan oxu cədvəli |
| Snapshot | Anlık görüntü | Uzun stream üçün periodic tam vəziyyət |
| Optimistic Concurrency | Optimist paralellik | Versiya yoxlaması ilə yazım konflikti |
| At-most / At-least / Exactly-once | Çatdırılma zəmanətləri | Mesaj çatdırılma semantikası |
| Deduplication | İkilərin təmizlənməsi | Eyni mesajın təkrar emalının qarşısı |
| FIFO | İrəli-girəvi sıra | Mesaj ardıcıllığı qorunur |
| Choreography | Xareoqrafiya | Hadisə dərci ilə idarəsiz saga |
| Orchestration | Orkestrasiya | Mərkəzi koordinator (SEC) idarə edir |
| Compensating Action | Kompensasiya addımı | Uğursuz addımın geri qaytarılması |
| Saga | Saqa | Paylanmış tranzaksiya zənciri |
| SEC | Saga icra koordinatoru | Saga addımlarını idarə edən komponent |
| Inbox / Outbox | Gələn / Gedən qutu | Tx-daxili mesaj yazımı pattern-i |
| Correlation ID | Korrelyasiya ID-si | Axını birləşdirən identifikator |
| Causation ID | Səbəb ID-si | Mesajı yaradan səbəbi qeyd edir |
| OpenTelemetry (OTel) | Açıq telemetriya | Standart tracing/metrics çərçivəsi |
| Span | Span | Trace-in bir hissəsi (vaxt intervalı) |
| Pact / CDCT | Müqavilə testi | Consumer-driven contract testing |
| ADR | Arxitektura qərar qeydi | Nygard formatında qərar sənədi |
| EventStorming | Hadisə fırtınası | Sticky-notes workshop dizayn metodu |
| AsyncAPI | Asinxron API spesifikasiyası | OpenAPI-nin async ekvivalenti |
| Modular Monolith | Modul monolit | Tək deploy, modul sərhədli daxili strukturla |

## Acronyms

- **EDA** — Event-Driven Architecture
- **DDD** — Domain-Driven Design
- **CQRS** — Command and Query Responsibility Segregation
- **2PC** — Two-Phase Commit
- **SEC** — Saga Execution Coordinator
- **CDCT** — Consumer-Driven Contract Testing
- **BDD** — Behavior-Driven Development
- **OTLP** — OpenTelemetry Line Protocol
- **EKS** — AWS Elastic Kubernetes Service
- **AZ** — Availability Zone
