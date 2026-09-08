# Chapter 8 — Message Workflows (səh. 216-247)

## Bu chapter nədən bəhs edir?

Distributed transaction-lar: 2PC vs Saga (choreography vs orchestration), SEC
(Saga Execution Coordinator), Command/Reply mesaj tipləri, Create Order sagasının
implementasiyası və saga repozitoriyası.

## Əsas fikirlər

### 1. Distributed transaction problemi
Bir əməliyyat 3 servisi əhatə edirsə (order + payment + shopping list), onlardan
biri çatarsa nə olur? Lokal ACID (atomicity, consistency, isolation, durability)
kənar komponentlərə uzanmır — **ACID-in paylanmış ekvivalenti axtarılır**.

### 2. 2PC (Two-Phase Commit)
| Fayda | Problem |
|---|---|
| Yaxşı bilinən protokol | Koordinator + iştirakçılar **bloklanır** (lock) |
| | Uzun müddətli kilidlər |
| | Yavaş iştirakçı hamını saxlayır |

### 3. Saga pattern
Uzunmüddətli əməliyyat = lokal transaction-lar zənciri + **compensating
actions** (geri qaytaran addımlar). Bir addım çatarsa, əvvəlkilərin kompensasiyası
işə düşür.

**2 növ:**
| Növ | Mexanizm | Mübadilə |
|---|---|---|
| **Choreography** | İştirakçılar hadisə dərc edir, digərləri reaksiya verir | Mərkəzi nəzarət YOX, amma axın görünməz olur |
| **Orchestration** | SEC koordinator addım-addım command göndərir | Mərkəzi axın, amma SEC tək nöqtədir |

### 4. SEC (Saga Execution Coordinator)
Komponentlər: Orchestrator (saga məntiqi + message stream), saga definitions,
reply topic dinləyicisi.
- İki rejim: manual start + cavab-lara reaktiv
- Hər cavab (reply) gələndə saga növbəti addıma keçir və ya kompensasiya başladır

### 5. Command/Reply tipləri
```go
// ddd paketində Event-in analoqu:
type Command interface { CommandName() string }
type Reply    interface { ReplyName() string }

// CommandMessage:
//   - xüsusi header: cavablar HARAYA göndərilsin (reply destination)
//   - gözlənti: hər command cavab (reply) qaytarır
//   - CommandHandler: HandleCommand(ctx, cmd) (Reply, error)
```

### 6. Create Order saga addımları
```go
saga := sec.NewSaga[CreateOrderData]("create-order", replyChannel)
saga.AddStep().
    Action(saga.authorizeCustomer)      // 1. Customers → AuthorizeCustomer
saga.AddStep().
    Action(saga.createOrder)            // 2. Order Processing → CreateOrder
    CompensateAction(saga.rejectOrder)  //    çatarsa → RejectOrder
saga.AddStep().
    Action(saga.confirmOrder)           // 3. Payments → ConfirmPayment
    CompensateAction(saga.cancelOrder)
saga.AddStep().
    Action(saga.initiateShopping)       // 4. Depot → CreateShoppingList
    CompensateAction(saga.cancelShoppingList)
saga.AddStep().
    Action(saga.approveOrder)           // 5. Order Processing → ApproveOrder
```
Hər iştirakçı modulda:
1. Protobuf-da command/reply mesajları təyin olunur
2. CommandHandler funksiyası (`doAuthorizeCustomer` — payload cast → domain
   çağırışı → Reply qur)
3. `Registrations()` — serializer qeydiyyatı
4. Subscriber abunəliyi (GroupName ilə)

### 7. Saga davamı (resumability)
- Saga instance-ları DB-də saxlanılır (saga stream cədvəli: saga, step, status)
- `Start()` → repozitoriyaya yaz; `HandleReply()` → status update + növbəti
  addım
- Proses çöksə saga **qaldığı addımdan davam edir** — at-least-once + idempotent
  addımlar

### 8. Nəticə
CreateOrder artıq heç bir xarici servisi birbaşa çağırmır — saga command-ları
göndərir, cavabları gözləyir; uğursuzluqda kompensasiya zənciri işləyir. Sistem
resilient olur: iştirakçı müvəqqəti yoxdursa saga gözləyir.

## Əsas terminlər

- Distributed Transaction (paylanmış tranzaksiya)
- ACID
- 2PC (Two-Phase Commit)
- Saga / Compensating Action
- Choreography vs Orchestration
- SEC (Saga Execution Coordinator)
- Command / Reply mesajları
- Saga Resumability (davam etdiriləbilənlik)

## Praktik nəticə

- 2PC-dən qaç — bloklanan kilidlər paylanmış sistemə ölümcüldür
- Mərkəzi görünüş lazımdırsa orchestration (SEC); sadə axınlar üçün choreography
- Hər addım üçün kompensasiya yaz; kompensasiyalar da idempotent olmalıdır
- Saga statusunu DB-də saxla — çökmədən davam etmək üçün
- Command-ın reply header-i saga-nın cavab yolunu müəyyənləşdirir

## Mənbə

Pages: 216-247 (Chapter 8, Event-Driven Architecture in Golang)
