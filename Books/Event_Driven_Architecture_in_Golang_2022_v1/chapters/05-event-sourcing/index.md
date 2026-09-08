# Chapter 5 — Tracking Changes with Event Sourcing (səh. 118-159)

## Bu chapter nədən bəhs edir?

Event sourcing-in monolith-ə əlavə olunması: es.Aggregate (versioning), event
store (optimistic concurrency), registry, snapshot-lar, "just enough CQRS" ilə
read modellər (mall, catalog) və event stream ömürləri.

## Əsas fikirlər

### 1. Event sourcing vs event streaming
- **Event sourcing:** Event store = həqiqət mənbəyi; cari vəziyyət = hadisələrin
  replay-i. Güclü konsistensiya + optimistic concurrency tələb edir.
- **Event streaming:** Broker (Kafka/NATS) — ünsiyyət kanalı; replay mümkün,
  amma source-of-truth deyil.

### 2. Event strukturu (private + options)
```go
type event struct {
    entity   Entity         // kim buraxdı
    name     string
    payload  Payload
    // private → yalnız EventOption ilə dəyişilir
}

func NewEvent(name string, payload Payload, options ...EventOption)
```
Aggregate-in AddEvent-i yalnız **aggregateEvent** yaradır — hadisə aggregate-ə
bağlıdır.

### 3. es.Aggregate — versioning qatı
```
Aggregate (ddd) → es.Aggregate (embedding) → Store / Product
```
- Hər dəyişiklikdə aggregate versiyası ARTIR
- Yeni hadisə aggregate-in cari versiyası ilə yazılır
- **Optimistic concurrency:** events cədvəlində compound primary key
  (stream_id + version) → eyni anda 2 yazımdən biri konflikt alır (birincinin
  versiyası qalib)

### 4. Registry — rekonstruksiya üçün
```go
// sıfır dəyərləri qeydiyyatdan keçir:
es.RegisterAggregate(StoreAggregate, Store{}) // Store{}
// Load zamanı: registry-dən yeni nüsxə → hadisələri tətbiq et
```

### 5. Event store sxemi
```
events(stream_id, version, event_id, name, payload, occurred_at, aggregated)
PRIMARY KEY (stream_id, version)
```
- `aggregated` markeri: hər hadisə yalnız bir dəfə read model-ə işlənir
- **Event versiyalanması:** Köhnə hadisəni DƏYİŞMƏK olmaz — yeni versiya hadisə
  (v2) buraxılır; köhnə stream-lər re-play oluna bilməli qalır

### 6. Just enough CQRS
Yazma tərəfi (aggregate repository) dəyişməz; **oxu tərəfi** üçün ayrıca
projection repoları:
- `MallRepository` — store-ların read modeli (plain SQL sorğuları)
- `CatalogRepository` — məhsulların oxu modeli
- Hadisə handler-ləri read model cədvəllərini doldurur (`CatalogHandlers`)

### 7. Middleware zənciri (AggregateStore)
```go
type AggregateStoreMiddleware func(AggregateStore) AggregateStore

// nümunə: EventPublisher middleware-i
eventPublisher := EventPublisher{publisher: publisher}
return func(store AggregateStore) AggregateStore {
    eventPublisher.AggregateStore = store
    return eventPublisher
}
```
Command handler-lər artıq publisher çağırmır — **repository commit edəndə
middleware** hadisələri dərc edir.

### 8. Snapshot-lar
- Uzun stream-lərdə replay bahalıdır → hər N hadisədən sonra aggregate-in **tam
  vəziyyəti** yazılır (kitabda hər 3 hadisədə bir — hardcoded strategiya)
- Load: ən son snapshot → yalnız qalan hadisələr tətbiq olunur
- Snapshot hər şeyi ehtiva etməlidir (dəqiq bərpa üçün)
- Stream ömrü: aggregate məhv olanda (store bağlandı) stream arxivlənə/silinə
  bilər

## Əsas terminlər

- Event Sourcing / Event Stream fərqi
- Optimistic Concurrency Control (optimist paralellik nəzarəti)
- Aggregate Versioning
- Registry (yenidənqurma reyestri)
- Projection / Read Model
- AggregateStoreMiddleware
- Snapshot (əyrilik şəkli)
- Event Versioning (hadisə versiyalanması)

## Praktik nəticə

- Event store: (stream_id, version) PK → paralel yazımlar avtomatik yaxalanır
- Hadisələr immutable; dəyişiklik = yeni versiya hadisəsi
- Oxu ehtiyacları üçün "just enough" CQRS: projection repo + hadisə handler
- Uzun stream-lər üçün snapshot; strategiyanı konfiqurasiya edilən et
- Publisher/dedupe kəsişən işləri middleware qatına qoy

## Mənbə

Pages: 118-159 (Chapter 5, Event-Driven Architecture in Golang)
