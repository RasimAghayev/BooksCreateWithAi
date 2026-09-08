# Chapter 4 — Event Foundations (səh. 94-117)

## Bu chapter nədən bəhs edir?

MallBots modular monolith-inin strukturu (internal packages, composition root,
gRPC/REST), modullararası sinxron inteqrasiya nümunələri və side-effect-lərin
domain event-lərə refactor edilməsi: AggregateBase, Event interfeysi,
EventDispatcher, ignoreUnimplementedDomainEvents pattern-i.

## Əsas fikirlər

### 1. Modular monolith strukturu
**Qovluq qaydaları (internal visibility):**
- `/root/internal` → yalnız `/root` və onun ağacı import edə bilər
- `/root/pkg-b/internal` → yalnız `/root/pkg-b` ağacı — `/root` və `pkg-a`
 onda YOX. Modullararası sərhəd Go-nun öz visible mexanizmi ilə qorunur.

**Generic layer adları QADAĞANDIR** (`controllers`, `services`, `models`) —
qovluqlar domain modulları adını daşıyır (baskets, stores, ordering...).

**Hexagonal (ports & adapters) ilə interface yerləşməsi fərqi:**
- "Accept interfaces, return structs" (consumer-side kiçik interfeys) — modul
  DAXİLİNDƏ
- Ports/contracts üçün — daha BÖYÜK interfeyslər mərkəzi yerdə (application/
  domain qovluğu), implementation-lardan sonra yazılır

### 2. Composition root
**Nədir:** Infrastructure + konfiqurasiya + app komponentlərinin birləşdirildiyi
yer — dependency injection burada baş verir. Sıra (sadə, proqnozlaşdırılan):
1. **Driven adapters** (DB, domain dispatcher) — yalnız infrastructure lazımdır
2. **Application** — Driven adapterlər lazımdır, Driver-lər YOX
3. **Driver adapters** (gRPC server, REST gateway) — infrastructure + application

Bu səviyyədə konkret dəyərlər üstünlük təşkil edir — abstraksiyaya ehtiyac azdır.

### 3. gRPC + REST
- Tək gRPC server bir neçə service daşıya bilər; namespace conflict olmasın
  deyə protobuf-lar parent qovluq prefiksi ilə compile olunur:
  `basketspb.Item` vs `orderingpb.Item`
- REST: `grpc-gateway` ilə — module-lar gRPC API-lərini REST-ə açır;
  Swagger UI http://localhost:8080/ ünvanında

### 4. Sinxron inteqrasiya nümunələri (indiki vəziyyət)
**AddItem (səbətə məhsul):** Baskets.AddItem → Stores.GetProduct (məlumat
anında çəkilir — data duplication yoxdur, amma hər sorğu cross-module çağırış).

**CheckoutBasket → CreateOrder (zəncirvari side-effect-lər):**
```
Baskets.CheckoutBasket
Ordering.CreateOrder
Customers.AuthorizeCustomer
Payments.CreateInvoice / ...
Stores... 
```
Bir handler-də bir neçə modula çağırış → birləşmə (coupling) artır.

### 5. Event tipləri (əvvəlki fəsillərdən qısa xatırlatma)
- **Domain event** — bounded context daxilində state dəyişikliyi haqqında
  məlumat; ən çox eyni prosesdə sinxron işlənir
- **Event sourcing event** — aggregate-in state-nin bərpası üçün davam
  (append-only) jurnal yazısı

### 6. Side-effect → Domain event refactoru
**Problemlərin yaranması:** `CreateOrder`-da notification birbaşa çağırılırsa:
```go
if err = h.orders.Save(ctx, order); err != nil { ... }
if err = h.notifications.NotifyOrderCreated(ctx, order.ID, order.CustomerID); err != nil { ... }
```
Bir qayda üçün yaxşıdır, real app-lərdə isə bir çox side-effect olur →
handler şişir, birləşmə artır.

**Həll — aggregate event-lər:**
```go
type Order struct {
    ddd.AggregateBase        // ID + event idarəsi kompozisiya ilə
    CustomerID string
    PaymentID  string
    // ...
}
```
- `AggregateBase` — ID sahəsi, `AddEvent(eventName, payload)`, `Events()` 
- **Aggregate event** — aggregate-in özündə yaranan, həyat dövrü ilə bağlı
  event (OrderCreated, OrderCanceled)

**Event interfeysi:**
```go
type Event interface {
    EventName() string   // bounded context daxilində unikal
}
type OrderCreated struct { ... }
func (e OrderCreated) EventName() string { return "events.order.created" }
```

### 7. ignoreUnimplementedDomainEvents pattern
DomainEventHandlers interfeysi böyüyəndə (yeni event əlavə olunanda) hər
handler-i dərhal yazmaq məcburiyyətindən qurtarmaq üçün:
```go
type ignoreUnimplementedDomainEvents struct{ storeNotificationHandlers }

func (ignoreUnimplementedDomainEvents) OnOrderCreated(...) error { return nil }
func (ignoreUnimplementedDomainEvents) OnOrderReadied(...) error { return nil }
// ...
```
- Embed etdiyi struct tərəfindən realləşdirilməyən metodlar boş qaytarır
- **Interfeys yoxlaması sayəsində** yeni event interfeysə əlavə olunanda
  compile xətası — update etməyi xatırladır (susmaq YOX)

### 8. EventDispatcher
```go
type EventDispatcher struct {
    mu       sync.Mutex
    handlers map[string][]handlerFunc
}

func (h *EventDispatcher) Subscribe(eventName string, handler handlerFunc) {
    h.mu.Lock()
    defer h.mu.Unlock()
    h.handlers[eventName] = append(h.handlers[eventName], handler)
}

func (h *EventDispatcher) Publish(ctx context.Context, events ...Event) error {
    for _, event := range events {
        for _, handler := range h.handlers[event.EventName()] {
            if err := handler(ctx, event); err != nil {
                return err
            }
        }
    }
    return nil
}
```
Sub-kod izahı: mutex ilə qorunan map — event adı → handler sırası; Publish
təbii olaraq sinxrondur (eyni proses, eyni goroutine).

### 9. Driver adapter kimi qeydiyyat
```go
func RegisterNotificationHandlers(
    subscriber EventSubscriber, handlers DomainEventHandlers,
) {
    subscriber.Subscribe(
        domain.OrderCreated{}, handlers.OnOrderCreated,
    )
    // ...
}
```
- NotificationHandlers artıq CreateOrder handler-inin parametri DEYİL —
  `domainDispatcher` (Driven) qurulur, handler-lər subscribe olunur
- Application constructor-u dəyişir: notifications asılılığı silinir →
  side-effect məntiqi handler-dən AYRILIR

## Əsas terminlər

- Modular Monolith (modul monolit)
- Composition Root (tərkib kökü)
- Internal Package Visibility (daxili paket görünməzliyi)
- Ports & Adapters (hexagonal)
- Domain Event / Aggregate Event
- EventDispatcher / Subscribe / Publish
- AggregateBase
- grpc-gateway
- Namespace Conflict (protobuf prefiksləri)

## Praktik nəticə

- Modul sərhədlərini Go internal qovluqları ilə tənzimlə — compiler qoruyur
- Composition root-daqı qurulma sırası: Driven → Application → Driver
- Side-effect-ləri aggregate event-lərə çıxar; handler yalnız business
  məntiqi saxlasın
- ignoreUnimplementedDomainEvents — interfeys böyüyəndə boş implementasiyalarla
  compile təhlükəsizliyi
- Domain event adları context-daxili unikal; diskrimininator kimi EventName()

## Mənbə

Pages: 94-117 (Chapter 4, Event-Driven Architecture in Golang)
