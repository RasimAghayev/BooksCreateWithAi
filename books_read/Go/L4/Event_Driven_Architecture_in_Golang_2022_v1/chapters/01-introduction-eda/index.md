# Chapter 1 — Introduction to Event-Driven Architectures (səh. 22-39)

## Bu chapter nədən bəhs edir?

EDA-nın təməli: hadisə (event) anlayışı, 4 hadisə ünsiyyət növü, producer/consumer
rolları, broker-lər (queue vs stream), MallBots nümunə tətbiqi, EDA-nın faydaları
(loose coupling, agility, scalability) və çətinlikləri (eventual consistency,
dual writes, debug).

## Əsas fikirlər

### 1. Event nədir?
**Hadisə = keçmişdə baş vermiş, dəyişməz fakt.** "PaymentReceived" — ödəniş
QƏBUL EDİLDİ (indi deyil, bitib). Adlar past tense olur.

### 2. Hadisə ünsiyyətinin 4 forması

**(a) Event Notification (hadisə bildirişi):**
```go
type PaymentReceived struct {
    PaymentID string // yalnız kimli — detal YOX
}
```
Consumer detal üçün geri çağırıq etməli → bağlılıq qalır.

**(b) Event-Carried State Transfer (state daşıyan hadisə):**
```go
type PaymentReceived struct {
    PaymentID string
    CustomerID string  // tam vəziyyət daxil edilir
    Amount     int
    // ...
}
```
Consumer geri müraciət etmədən öz lokal nüsxəsini saxlaya bilir.

**(c) Event Sourcing (hadisə mənbəyi):** Hadisələr ünsiyyət üçün deyil,
**event store-da saxlanılır** — cari vəziyyət = hadisələrin replay-i.

**(d) Event Streaming (hadisə axını):** Retention + replay imkanlı (Kafka,
NATS JetStream); system-lər sonra qoşula bilir və tarixə baxır.

### 3. Queue vs Stream
| Broker növü | Təyinat | Xüsusiyyət |
|---|---|---|
| Message Queue | Sadə pub/sub | Mesaj götürüldükdə silinir |
| Event Stream | Retention + replay | Yeni consumer keçmişi də oxuyur |

### 4. Komponentlər
- **Producer** — vəziyyət dəyişəndə hadisəni broker-ə dərc edir
- **Consumer** — abunə olur, hadisələri emal edir
- **Broker** — arada; producer/consumer bir-birini BİLMİR (loose coupling)

### 5. MallBots layihəsi
Kitab boyu qurulan demo: mall-da robotlarla alış-veriş idarəetməsi. Modullar:
Customers, Stores, Products, Baskets, Orders, Payments, Depot, Notifications.
Hexagonal arxitektura ilə (sync API + async hadisə qarışığı).

### 6. Faydalar
- **Loose coupling:** Analytics komandası hadisə axınından istədiyini götürür —
  heç kimlə razılaşma lazım deyil
- **Agility:** Yeni istifadəçi (consumer) mövcud hadisələrə qoşulur
- **Scalability, resiliency, UX:** Komponentlər ayrıca miqyaslanır; bir komponent
  çöksə, hadisələr broker-də gözləyir
- **Analytics/auditing:** Hadisələr təbii audit izidir

### 7. Çətinliklər
| Problem | Təsvir | Həll istiqaməti |
|---|---|---|
| **Eventual consistency** | Vəziyyət dərhal yox, "sonda" uyğunlaşır | UX-də bunu dəstəklə; task-based UI |
| **Dual writes** | DB yazma + hadisə göndərmə 2 ayrı əməliyyat — biri çatışmır | Hadisəni DB-yə atomik yaz (outbox, ch9) |
| **Dist. async workflows** | Uzun zəncirlər idarəçiliyi | Saga pattern (ch8) |
| **Debugging** | Mənşəyi izləmək çətin | Correlation ID, observability (ch12) |

## Əsas terminlər

- Event (hadisə) / Fact (fakt)
- Event Notification vs State Transfer
- Event Sourcing / Event Stream
- Producer / Consumer / Broker
- Message Queue
- Loose Coupling (boş bağlılıq)
- Eventual Consistency (sonrakı uyğunluq)
- Dual Write (ikiqat yazma problemi)

## Praktik nəticə

- Hadisə adları past tense; məzmun "nə lazımdır" sualına cavab versin
- Notification (yüngül, amma callback lazım) vs state transfer (kopya saxla,
  amma böyük payload) tradeoff-unu şüurlu seç
- Dual write EDA-nın 1 nömrəli tələsidir — outbox pattern ilə həll olunacaq
- EDA "free" deyil: eventual consistency mədəniyyəti tələb edir

## Mənbə

Pages: 22-39 (Chapter 1, Event-Driven Architecture in Golang)
