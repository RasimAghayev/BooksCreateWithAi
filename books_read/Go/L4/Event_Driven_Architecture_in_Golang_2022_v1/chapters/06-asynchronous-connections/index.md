# Chapter 6 — Asynchronous Connections (səh. 160-191)

## Bu chapter nədən bəhs edir?

Asinxron kommunikasiya: mesaj anlayışı və növləri, integration pattern-lər
(notifications vs event-carried state transfer), delivery guarantee-lər
(at-most/at-least/exactly-once), ordered delivery problemləri, NATS JetStream
implementasiyası (am/jetstream paketləri, eventStream, proto marshalling) və
Store Management → Shopping Baskets asinxron axını.

## Əsas fikirlər

### 1. Mesaj nədir?
**Event = mesaj, amma mesaj həmişə event deyil.** Mesaj — payload + əlavə
informasiya daşıyan konteynerdir. Növlər:
| Növ | Ömür | Sərhəd | Versiya | Emal |
|---|---|---|---|---|
| Domain event | ən qısa | app daxili | lazım deyil | sinxron |
| Event-sourced event | ən uzun | servis daxili | TƏLƏB OLUNUR | sinxron (replay) |
| Integration event | uzun | modullararası | tələb olunur | asinxron |

### 2. Asinxron inteqrasiya pattern-ləri
**Notifications ( bildiriş):** Producer "dəyişdi" deyir → Consumer geriyə
callback edib datanı soruşur. Problemlər: Producer scalable olmalı (bilinməyən
sayda callback), Consumer Producer-in API-sini bilməli → **tam decoupling YOX**.
Multi-notification latency-də eyni data üçün artıq callback-lər (debounce/
serializasiya həlləri çaşdırıcı olur).

**Event-carried state transfer (stateful events):** Event ÖZÜNÜNÜZÜ daşıyır →
Consumer Producer-ə müraciət etmir. **Əsas üstünlük:** Consumer Producer-in
availability-dən müstəqildir (temporal decoupling). Tələb: Consumer-ın lokal
kopiyası saxlaması.

**Eventual consistency tələsi (read-after-write):** yazışdan DƏRHAL sonra oxu →
replikaya düşəndə stale data qayıda bilər. UI bunu nəzərə almalıdır.

### 3. Delivery guarantee-lər
| Model | Nə təmin edir | Qiymət |
|---|---|---|
| At-most-once | Heç bir ack gözlənilmir — mesaj İTƏ BİLƏR | dedup lazım deyil, amma itki riski |
| At-least-once | Broker təkrar çatdırır — DUPLİKAT mümkün | consumer dedup/idempotency TƏLƏB OLUNUR |
| Exactly-once | Tək dəfə çatdırma | dedup + transaction + ack-sırası koordinasiyası |

**Dedup (transaction ilə):** DB transactionunda (1) mesaj identity yoxla/yaz,
(2) mesajı emal et → xətanda rollback (identity də gedir), uğurda commit
SONRA broker-ə ack. Təkrar gələndə identity artıq var → skip.

### 4. Ordered delivery
- **FIFO single consumer:** istehsal ≤ istehlak olduqda kifayət
- **Competing consumers:** yükü bölür, amma SIRA pozulur — 2-ci mesaj 1-cidən
  əvvəl emal oluna bilər. Həllər: requeue gözləməsi (queue deadlock riski!),
  zaman pəncərəsi (latency), **partitioned queue** (eyni aggregate = eyni
  partition = sıra qorunur)
- **Goroutine tələsi:** consumer daxilində goroutine-lə paylarsan competing
  consumers problemi geri qayıdır

### 5. NATS JetStream
Core NATS üstünə: **durable streams** (subscription-dan əvvəl də yayınlanmış
mesajlar), consumer cursor-lar, **message dedup**, exactly-once semantikası.

**Paket arxitekturası:**
- `/internal/am` — ümumi interfeyslər: `MessagePublisher`, `MessageSubscriber`
  (generic `MessageHandler[T]`), `MessageStream`
- `/internal/jetstream` — NATS-spesifik implementation; `rawMessage`
  aralıq tipi (id, name, data) — event seriyalaşdırma NATS-dan ayrı saxlanılır

### 6. eventStream — seriyalaşdırma qatı
```go
func (s eventStream) Publish(ctx context.Context, topicName string, event ddd.Event) error {
    payload, _ := s.reg.Serialize(event.EventName(), event.Payload())
    data, _ := proto.Marshal(&EventMessageData{   // protobuf konteyner
        Payload:    payload,
        OccurredAt: timestamppb.New(event.OccurredAt()),
        Metadata:   metadata,
    })
    return s.stream.Publish(ctx, topicName, rawMessage{
        id: event.ID(), name: event.EventName(), data: data,
    })
}
```
- Publish yalnız `ddd.Event` qəbul edir; Subscribe yalnız
  `MessageHandler[EventMessage]` — konkretlik DRY-dən üstün (universal
  handler YOX)
- Subscribe: anonymous funksiya RawMessage handler kimi — deserializasiya
  orada, alınan EventMessage istehlakçıya ötürülür

### 7. Monolith konfiqurasiyası
```go
type NatsConfig struct { URL, Stream string }
```
JetStream context init:
```go
js.AddStream(&nats.StreamConfig{
    Name:     cfg.Stream,
    Subjects: []string{fmt.Sprintf("%s.>", cfg.Stream)}, // wildcard subscription
})
```
Graceful shutdown: `closed` channel (semaphore) + `errgroup`: bir goroutine
closed gözləyir, digəri `nc.Drain()` — bağlanma təmiz baş verir. `Monolith`
interfeysinə `JS()` metodu əlavə olunur — module-lar JetStream context-i alır.

### 8. Integration event qaydaları
- Event-lər PUBLİK olmalı → `storespb` paketində (protobuf) təyin edilir
- Müstəqil (standalone) — request/response mesajları ilə qarışdırılmır
- Duplicate adlar: domain event `StoreCreated` + integration event `StoreCreated`
  eyni adda ola bilər — yalnız hər ikisini görən modul fərq edir; istəsən
  fərqli ad ver
- `Registrations()` funksiyası + `ProtoSerde` ilə registry-yə qeydiyyat —
  consumer tərəf sadə istifadə edir

### 9. Sender (Store Management)
IntegrationEventHandlers domain event-ləriə qulaq asır → protobuf event
publish edir:
```go
domainSubscriber.Subscribe(eventHandlers,
    domain.StoreCreatedEvent,
    domain.StoreParticipationEnabledEvent,
    // ...
)
// onStoreCreated: domain payload → storespb.StoreCreatedWasm integration event
```

### 10. Receiver (Shopping Baskets)
- `StoreHandlers` — adi event handler (`HandleEvent`), integration event-ləri
  emal edir; logging middleware eyni işləyir
- Subscribe fərqli: eventStream-ə `MessageHandlerFunc` adapteri ilə sarılır
  (HandleEvent != HandleMessage → type-safe wrap)
- Nəticə log-ları: store yarat → `Stores.Catalog.On(stores.StoreCreated)`
  asinxron axın işləyir

## Əsas terminlər

- Message / Event fərqi
- Integration Event (inteqrasiya hadisəsi)
- Event-Carried State Transfer (hadisə ilə state ötürülməsi)
- Temporal Decoupling (zaman ayrılığı)
- Read-after-write inconsistency
- At-most/at-least/exactly-once delivery
- Idempotency / Deduplication
- Competing Consumers (rəqabətli istehlakçılar)
- Partitioned Queue (bölünmiş növbə)
- NATS JetStream / durable stream
- errgroup / Drain

## Praktik nəticə

- Consumer Producer-dən asılı olmasın deyə state-i event-in içində göndər
- At-least-once + idempotent consumer — praktik standart; exactly-once
  transaction+dedup ilə qurulur
- Sıra lazımdırsa: single consumer və ya partition (aggregate ID üzrə)
- Seriyalaşdırma (proto/serde) broker inteqrasiyasından AYRI qatda — am
  interfeys + jetstream implementation
- Integration event-lər protobuf-da, standalone, registry-yə qeydli —
  modul sərhədindən kənar kontrakt

## Mənbə

Pages: 160-191 (Chapter 6, Event-Driven Architecture in Golang)
