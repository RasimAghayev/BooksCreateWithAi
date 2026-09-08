# Chapter 6 — Asynchronous Connections (səh. 160-191)

## Bu chapter nədən bəhs edir?

NATS JetStream ilə asinxron inteqrasiya: mesaj növləri, at-most/at-least/exactly-once
çatdırılma, dedupe, FIFO sıralama, JetStream stream-lər, protoserde, Store
Management modulunun asinxronlaşdırılması.

## Əsas fikirlər

### 1. Mesaj növləri (3lük)
| Növ | Ömür | İstifadə |
|---|---|---|
| **Domain event** | Ən qısa — domain daxili, broker-ə çıxmır | daxili yan təsirlər |
| **Integration event** | Modullararası, broker-dən keçir | asinxron inteqrasiya |
| **Command/Reply** | Sorğu-cavab məntiqi | saga iştirakçıları (ch8) |

### 2. Notification vs state transfer (callback problemi)
- Notification → consumer detail üçün producer-i geri çağırır → 2-3 dəfə
  lazımsız callback, producer yükü artır
- **State transfer → temporal decoupling:** consumer producer mövcudluğundan
  asılı deyil, lokal nüsxədən oxuyur

### 3. Eventual consistency tələsi
Uğurlu cavab dərhal oxu modelində görünməyə bilər — "read-your-writes" pozula
bilər; UX bunu dəstəkləməlidir.

### 4. Çatdırılma zəmanətləri
| Rejim | Davranış | Tələb |
|---|---|---|
| **At-most-once** | Ack gözləmir — itərsə itir | itki tolere edilən system |
| **At-least-once** | Təkrar çatdırıla bilər | **dedupe** qabiliyyəti (idempotent emal) |
| **Exactly-once** | Nəzəri — əslində at-least-once + dedupe | broker + alıcı əməkdaşlığı |

Dedupe pəncərəsi: mesaj ID + tranzaksiya daxilində "emal olundu" qeydi.

### 5. Sıralama (ordering)
- FIFO queue → tək consumer ardıcıllıq alır, amma publish rate > emal rate
  olarsa selector queue dolur
- **Requeue strategiyası:** sıradan çıxmış mesaj sona qaytarılır → sıra pozulur
- Partitioned queues → müəyyən açara görə eyni partition-a düşür (sıra qorunur,
  paralellik qazanılır); amma consumer özü də asinxron emal edərsə yenə pozuntu
  mümkündür

### 6. NATS JetStream
- Core NATS: at-most-once (sadə pub/sub)
- **JetStream:** durable streams, cursor-lu consumer-lar (yerini yadda saxlayır),
  replay, ack

**Paket strukturu (/internal/nats):**
```
MessagePublisher interface
MessageSubscriber interface → Message qaytarır
RawMessage (aracı) → Deserialize üçün Serde-ə keçir
EventMessage → event.Event-i wrap edir
```
Serde: `Serialize(event, payload)` / `Deserialize(data)` — JSON/Proto dəstəyi.

### 7. Store Management modulunun asinxronlaşdırılması
- Monolith konfiqurasiyasına NATS URL əlavə; `nc.Close()` defer; JetStream
  context init
- `errgroup` ilə mesaj axını goroutine-i idarə olunur (Ctrl-C → close channel)
- `Monolith.JS()` metodu — modullara JetStreamContext təqdim edir
- **Protobuflarda integration event-lər** təyin olunur (api / proto paketi)
- Ad toqquşması: domain və integration hadisələri eyni adda OLA BİLƏR — fərqli
  paketlərdədirlər; registry-də ayırmaq lazım (alias-lar)
- Domain → integration çevrilməsi: domain handler integration hadisəsini buraxır
  (`On(StoreCreated) → publish(StoreCreated pb event)`)

### 8. Subscriber adapter
```go
// handler DomainEventHandlers HandleEvent imzasına malikdir, MessageHandler YOX
// → RawMessage → Deserialize → HandleEvent körpüsü
```
Hər modul öz `Registrations()` funksiyası ilə abunə olacağı mövzuları qeyd edir.

## Əsas terminlər

- Integration Event
- Temporal Decoupling (zaman qiymət ayrılığı)
- At-most / At-least / Exactly-once
- Deduplication (ikilərin təmizlənməsi)
- FIFO / Partitioned Queue
- NATS JetStream / Durable Stream / Cursor
- Serde (serializasiya/deserializasiya)
- errgroup

## Praktik nəticə

- Integration hadisələri protobuf-da təyin et — kontrakt versiyalanır
- At-least-once + dedupe reyestrini birbaşa kodunuza salın
- Ordering lazımdırsa partition açarı seç (məs. aggregate ID)
- Domain hadisəsi brokerə birbaşa getmir — domain → integration çevirici qatı

## Mənbə

Pages: 160-191 (Chapter 6, Event-Driven Architecture in Golang)
