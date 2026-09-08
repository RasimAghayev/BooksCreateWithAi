# Chapter 5 — Tracking Changes with Event Sourcing (səh. 118-159)

## Bu chapter nədən bəhs edir?

Event sourcing-in tam implementasiyası: append-only event stream, Event interfeysinin
genişləndirilməsi (Metadata, EventOption), event-sourced Aggregate (versiya
izləmə), Registry + serde, AggregateRepository/AggregateStore, PostgreSQL event
store (optimistic concurrency), CQRS read model-lər, AggregateStoreMiddleware
(EventPublisher), snapshot-lar.

## Əsas fikirlər

### 1. Event sourcing nədir?
**Nədir:** Hər aggregate dəyişikliyini **append-only stream**-ə yazmaq; aggregate-in
son vəziyyətini bərpa etmək üçün bütün event-ləri ardıcıl tətbiq etmək (replay).

**Tələblər:**
- Event store **strong consistency** təmin etməli
- **Optimistic concurrency:** eyni stream-ə paralel yazılışda yalnız İlKİ
  uğur qazanır; qalanlar retry və ya fail

**Vacib fərq:** Event streaming ≠ event sourcing — amma bir-birini gücləndirir;
domain/aggregate interaksiyaları pozulanda həm sourcing, həm distribution üçün
eyni event-lərdən istifadə olunur.

### 2. Yenilənmiş Event modeli
Chapter 4-ün sadə event-inin üzərinə: **EventOption** (variadic), **Metadata**
(aggregate adı, versiya), private event struct (Entity embed):
```go
evt := newEvent(name, payload, options...)  // EventOption-lar tətbiq olunur
```
- `...EventOption` — variadic SON parametr olmalı; event qaytarılmazdan əvvəl
  dəyişilir (məs. aggregate məlumatı Metadata ilə əlavə olunur)
- Go privacy = PAKET səviyyəsi: kiçik hərf = yalnız paket daxili. Private
  sahələri xaricdən yalnız EventOption ilə dəyişmək olur

### 3. Aggregate + AggregateEvent
```go
func (a *Aggregate) AddEvent(name string, payload EventPayload, options ...EventOption) {
    options = append(options, Metadata{
        AggregateNameKey: a.name,   // aggregate məlumatı avtomatik qeyd olunur
    })
    // aggregateEvent{event: newEvent(name, payload, options...)}
}
```
- Event adı artıq payload-da deyil, AddEvent-in parametri
- **Generics (Go 1.18+) həlli:** EventDispatcher generic olur → yeni
  AggregateEventDispatcher yazmadan type-safety qorunur (əvvəl: ya
  duplicate dispatcher, ya type-cast)

### 4. es.Aggregate — event-sourced aggregate
ddd.Aggregate üzərinə versiya qatı:
```go
func (a *es.Aggregate) AddEvent(...) {
    options = append(options, ddd.Metadata{
        ddd.AggregateVersionKey: a.PendingVersion() + 1,  // versiya izləmə
    })
    a.Aggregate.AddEvent(...)
}
```
Hər event-sourced aggregate **EventApplier** implement edir:
```go
type EventApplier interface {
    ApplyEvent(ddd.Event) error
}

func (p *Product) ApplyEvent(event ddd.Event) error {
    switch payload := event.Payload().(type) {   // concrete type-a görə
    case *ProductCreated: ...
    case *ProductRemoved: ...
    }
    return nil
}
```
**Events = source of truth:** aggregate dəyərləri birbaşa təyinatla YOX,
event-lərin tətbiqi ilə dəyişir.

### 5. Event payload dəqiqləşməsi
Domain event-də bütövlükdə aggregate ötürülə bilirdi; sourcing-də event YALNIZ
dəyişən sahələri daşıyır (`StoreCreated{Name, Location}`).

**Event-lər versiyalanmalıdır:** event-in məzmunu dəyişdirilə BİLMƏZ (DB-də
qalıb!) — yeni event yaradılır; ApplyEvent köhnə event-ləri də işləməyə
davam etməlidir. App unutmaq hüququ没有.

### 6. Registry + serde
```go
agg := registry.Build(registry.ID(id), registry.Name(name))  // concrete instance
```
- Tip qeydiyyatı: `RegisterEvent(name, zeroValue, serde)` — zero value-lar
- **Serde** = serializer/deserializer; fərqli obyekt qrupları fərqli serde ala bilər
- Load zamanı interface YOX — concrete tip qaytarılır (type-safety)

### 7. AggregateRepository vs AggregateStore
- **AggregateRepository** (generic): Load (registry.Build → store.Load) və
  Save (yeni event-ləri tətbiq et → store-a yaz)
- **AggregateStore** (interfeys): infrastructure bağlantısı (PostgreSQL
  implementation)

### 8. PostgreSQL event store
`events` cədvəli — **compound primary key** (stream id + name + version):
eyni versiyaya iki yazış gələndə ikincisi conflict alır → optimistic
concurrency DB səviyyəsində.

### 9. CQRS mecburiyyəti
Event store yalnız stream əməliyyatları edir; sorğular (GetStores kimi) üçün
**read model** lazımdır → "just enough CQRS":
- `MallRepository` — store event-lərini read model-ə **proyeksiya** edir +
  sorğuları icra edir
- `CatalogRepository` — məhsullar üçün eyni pattern
- Event sourcing CQRS-siz çətin; CQRS sourcing-siz MÜMKÜNDÜR

### 10. EventHandler refactoring
```go
type EventHandler interface { HandleEvent(ctx, Event) error }
type EventHandlerFunc func(ctx context.Context, event ddd.Event) error
```
- `dispatcher.Subscribe(name, ddd.EventHandlerFunc(myHandlerFn))` — funksiya
  adapteri (http.HandlerFunc kimi)
- Unhandled event artıq panic yoxdur → ignoreUnimplementedDomainEvents qalır

### 11. AggregateStoreMiddleware — EventPublisher
Problem: Save uğurlu olsa da publish uğursuz → event-lər itər (və ya əksinə).
**Middleware zənciri** həlli:
```go
func (p EventPublisher) Save(ctx context.Context, agg EventSourcedAggregate) error {
    if err := p.AggregateStore.Save(ctx, agg); err != nil {  // DB-yə yaz
        return err
    }
    return p.publisher.Publish(ctx, agg.Events()...)          // sonra publish
}
// Quruluş: middleware.AggregateStore(eventStore, eventPublisher, ...)
```
Hər middleware `func(store AggregateStore) AggregateStore` qaytarır — zəncir
(Others → EventPublisher → Another → EventStore). Command handler-lərdən
EventPublisher asılılığı silinir.

### 12. Snapshot-lar
**Problem:** Uzun stream-in tam replay-i bahalıdır.
**Həll:** Periodik olaraq aggregate-in serializasiyası + versiyası saxlanılır:
- `ApplySnapshot()` — köhnə snapshot strukturları da işləməli (`switch ss :=
  snapshot.(type)` ilə `*ProductV1` kimi versiyalı tiplər)
- SnapshotStore da AggregateStore interfeysini qane edir (xüsusi interfeys YOX)
- Load: snapshot varsa tətbiq et → event store yalnız snapshot-dan SONRAKI
  versiyaları yükləyir
- Kitabda hər 3 event-dən bir snapshot (demo üçün; real sistemdə daha
  ağıllı strategiya lazımdır)

## Əsas terminlər

- Event Sourcing (hadisə mənbəyi)
- Append-only Stream (yalnız-əlavə axını)
- Replay / EventApplier
- Optimistic Concurrency (optimist uyğunluq)
- Compound Primary Key (mürəkkəb birincil açar)
- Registry / Serde (qeydiyyat / seriyalaşdırıcı)
- Read Model / Projection (oxu modeli / proyeksiya)
- CQRS
- EventHandlerFunc adapteri
- AggregateStoreMiddleware
- Snapshot (anlıq görüntü)

## Praktik nəticə

- Event sourcing = event-lər source of truth; dəyərlər yalnız ApplyEvent ilə
  dəyişir
- Event-in məzmunu dəyişilməz — yeni event + köhnə event-lərin dəstəyi
- Concurrency: DB compound key (id+name+version) ilə optimistic control
- Sorğular üçün read model + proyeksiya (CQRS) — event store sorğu dəstəyi
  vermir
- Uzun stream-lərdə snapshot + snapshot-dan sonrakı event-lərin yüklənməsi

## Mənbə

Pages: 118-159 (Chapter 5, Event-Driven Architecture in Golang)
