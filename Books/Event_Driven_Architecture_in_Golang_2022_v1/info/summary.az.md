# Event-Driven Architecture in Golang — Xülasə (AZ)

**Müəllif:** Michael Stack · Packt 2022 · 384 səh. · L4 Advanced

## Kitabın bir cümləlik özəyi

MallBots demo tətbiqi üzərindən sinxron modular monolith-in event-driven
mikro-servislərinə çevrilməsi: domain hadisələrindən başlayaraq event sourcing,
NATS JetStream, saga, inbox/outbox, contract testing və AWS deploy-a qədər
tam prod haltedə növbə.

## Hissə-hissə xülasə

### Part 1 — Fundamentals (Ch 1-3)
- **Event** = keçmişdə baş vermiş dəyişməz fakt; 4 forma: notification, state
  transfer, event sourcing, event streaming
- EDA faydaları: loose coupling, agility, scale; çətinliklər: eventual
  consistency, dual writes, debugging
- DDD: bounded context + ubiquitous language; hexagonal (ports/adapters);
  CQRS yazu/oxu ayrılığı; modular monolith → microservices təkamül yolu
- EventStorming: kağız + sticky notes ilə business axınını kəşf; storytelling
  implicit hadisələri üzə çıxarır; nəticə BDD/Gherkin + ADR qeydləri

### Part 2 — Komponentlər (Ch 4-9)
- **Domain events:** aggregate hadisələri yığır, commit zamanı dispatch;
  EventDispatcher + subscriber interfeysləri
- **Event sourcing:** event store = source of truth; (stream_id, version) PK
  ilə optimistic concurrency; registry-dən rekonstruksiya; hər N hadisədə
  snapshot; hadisələr immutable — dəyişiklik yeni versiya
- **NATS JetStream:** durable stream + cursor; integration hadisələri
  protobuf-da; domain→integration çevirici qat
- **State transfer:** hər modul lokal keş + idempotent upsert + sync fallback;
  AsyncAPI sənədləşməsi
- **Saga:** 2PC blokladığı üçün orchestration (SEC koordinator) + kompensasiya
  addımları; Command/Reply mesajları; saga DB-də davamlıdır
- **Transactional messaging:** dual write outbox ilə həll — publish eyni DB
  tranzaksiyasında; inbox dedupe; scoped DI konteynerlər

### Part 3 — Production (Ch 10-12)
- Test piramidi: unit (mock) → integration (Docker konteyner + testify suite +
  build tag) → contract (Pact CDCT + Broker) → e2e (BDD/godog)
- Deploy: eyni composition root monolith/microservices rejimləri; nginx
  reverse proxy; Terraform (state S3-də) + EKS/Helm/K9s
- Observability: correlation/causation ID-lər; OpenTelemetry interceptor-ları;
  span event-lər; Prometheus counter/histogram; Jaeger trace baxışı

## Kitabın ən dəyərli 5 yeri

1. **Ch 9 Inbox/Outbox** — dual write EDA-nın 1 nömrəli xətasıdır; hər publish
   DB tx-inə bağlanmalıdır
2. **Ch 5 Event store sxemi** — (stream_id, version) PK: paralelizm və sıra
   problemi bir zərbədə həll
3. **Ch 8 Saga addım qrafı** — kompensasiya zəncirinin kod görünüşü
4. **Ch 6 çatdırılma semantikası** — at-least-once + dedupe = effectively-once
5. **Ch 10 Pact messaging testi** — event müqavilələrinin real publisher ilə
   doğrulanması

## Kim üçün?

Go bilən, mikro-servis/event-driven dünyasına keçən developersiya üçün. DDD
əsasları tanış olmalıdır. Başlanğıc səviyyəsi üçün deyil.
