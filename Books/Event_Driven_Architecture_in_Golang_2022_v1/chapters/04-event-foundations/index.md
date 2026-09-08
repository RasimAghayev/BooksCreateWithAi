# Chapter 4 — Event Foundations (səh. 94-117)

## Bu chapter nədən bəhs edir?

MallBots modular monolith-inin strukturu (internal paketlər, composition root,
driver/driven adapterlər, Docker Compose), modullararası inteqrasiya və domain
hadisələri ilə yan təsirlərin (side effect) refaktoru: Event interface,
aggregate record-ları, EventDispatcher.

## Əsas fikirlər

### 1. Modular monolith struktur
```
/root/internal/          → bütün modullar (customers, orders, ...)
/root/internal/common   → paylaşılan infra (api, di/containers)
/root/cmd/mallbots      → giriş nöqtəsi
```
Hər modul öz hexagonunu daşıyır: Driver adapterlər (gRPC/REST) solda, Driven
adapterlər (Postgres, NATS) sağda.

### 2. Ports & adapters inteqrasiyada
- Hər modul **öz port interfeyslərini** təyin edir (istehlakçı tərəfi)
- Modullararası əlaqələr gRPC üzərində (namespace konfliktinə qarşı protobuflar
  ayrı saxlanılır; tək gRPC server bir çox service qulluq edə bilər)

### 3. Composition root
Hər modul **eyni başlanğıc patterni** ilə qalxır:
```go
// Driven adapters — portları reallaşdırır (Postgres, dispatcher)
grpc.RegisterServer(ctx, app, mono.RPC()) // Driver adapter — sorğuları qəbul edir
```
- **Driver adapterlər** tətbiqə "sürür" (gRPC, REST, cron)
- **Driven adapterlər** tətbiq tərəfindən "idarə olunur" (DB, mesaj broker)
- Docker Compose ilə bütün sistem qalxır (docker compose up); dəyişiklikdən
  sonra down → up

### 4. Problem: yan təsirlər handler-də
Köhnə yanaşma — CheckoutBasket handler birbaşa 3 çağırış edir: order yarat +
payment başlat + notification göndər. Modullar bir-birinə yapışıq (temporal
coupling).

### 5. Domain events ilə həll
**Domain event** — bounded context daxilində baş verən hadisə (inteqrasiya
hadisəsindən fərqli). Yan təsirlər implicit qaydalara çevrilir.

**Aggregate-lərdə hadisə saxlanması:**
```go
type aggregate interface {
    ID() string
    Events() []event.Event    // topladığı hadisələr
    CommitEvents()            // handler uğurla bitəndə təmizlənir
}
```
Hər model dəyişəndə hadisə **yığılır** (dərhal göndərilmir!) — handler-ın
özü müvəffəqiyyətindən sonra dispatch olur (uğursuzluqda hadisələr itmir).

**Event interfeysi:**
```go
type Event interface {
    EventName() string // identifikasiya + router üçün açar
}

type OrderCreated struct {
    OrderID string
    // ...
}
func (OrderCreated) EventName() string { return "orders.created" }
```

**Subscriber interfeysi (hadisə-yə görə metod):**
```go
type EventSubscriber interface {
    OnOrderCreated(evt OrderCreated) error
    OnOrderReadied(evt OrderReadied) error
    // ...
}
```
- Yeni hadisə əlavə olunca interfeys genişlənir → **kompilyator** reallaşdırmayan
  abunəçini yaxalayır
- Köhnə kodun sınmaması üçün "ignore" embed:
```go
type ignoreUnimplementedDomainEvents struct{}
func (ignoreUnimplementedDomainEvents) OnOrderCreated(...) error { return nil }
```

### 6. EventDispatcher
```go
type EventDispatcher struct {
    mu       sync.Mutex
    handlers map[string][]event.Handler
}

func (h *EventDispatcher) Subscribe(event event.Event, handler event.Handler) {
    h.mu.Lock()
    defer h.mu.Unlock()
    h.handlers[event.EventName()] = append(h.handlers[event.EventName()], handler)
}

func (h *EventDispatcher) Publish(ctx context.Context, events ...event.Event) error {
    // hər hadisə üçün abunə olan bütün handler-ləri çağırır
}
```
Composition root-da bağlanış:
```go
domainDispatcher.Subscribe(domain.OrderCanceled{}, notificationHandlers.OnOrderCanceled)
```
Notification artıq order module-ün DAXİLİNƏ yazılmır — hadisə ilə subscribe
olunur. Modul asılılıqları azalır, yan təsirlər sınanabilt.

## Əsas terminlər

- Modular Monolith / Composition Root
- Driver / Driven Adapter
- Domain Event vs Integration Event
- Aggregate Event Collection (hadisə yığımı)
- EventDispatcher / Subscriber
- Temporal Coupling (zaman bağlılığı)
- Side Effect (yan təsir)

## Praktik nəticə

- Yan təsirləri handler-dən hadisəyə köçür — "nə edirəm" yox, "nə baş verdi"
- Aggregate hadisələri yığır; commit zamanı dispatch — transaction təhlükəsizliyi
- Hadisə interfeysləri kompilyator ilə subscriber-ləri dəstəkləyin (yeni
  hadisə = mütləq imzası görünür)
- Dispatcher mutex ilə qorunur; subscribe/publish composition root-da bağlanır

## Mənbə

Pages: 94-117 (Chapter 4, Event-Driven Architecture in Golang)
