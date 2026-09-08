# Event-Driven Architecture in Golang — Cheat Sheet (AZ)

## 1. Event qərarı sürətlə

```
Consumer detal üçün geri çağırır?      → Event-Carried State Transfer
Tək fakt bildirmək kifayətdir?         → Event Notification
Cari vəziyyət tarixdən qurulmalıdır?   → Event Sourcing
Keçmişi replay etmək lazımdır?         → Event Stream (JetStream)
```

## 2. Event store sxemi (Ch 5)

```sql
CREATE TABLE events (
  stream_id   TEXT    NOT NULL,
  version     INT     NOT NULL,
  event_id    TEXT    NOT NULL,
  name        TEXT    NOT NULL,
  payload     JSONB   NOT NULL,
  occurred_at TIMESTAMPTZ,
  aggregated  BOOLEAN DEFAULT FALSE,
  PRIMARY KEY (stream_id, version)   -- optimistic concurrency
);
```

## 3. Aggregate + hadisə yığımı (Ch 4-5)

```go
type Aggregate interface {
    ID() string
    Events() []event.Event
    CommitEvents()
}

// dəyişiklik zamanı:
agg.AddEvent("stores.created", payload)  // yalnız yığılır
// handler uğurlu → CommitEvents + dispatch (middleware publish edir)
```

## 4. EventDispatcher (Ch 4)

```go
type EventDispatcher struct {
    mu       sync.Mutex
    handlers map[string][]Handler
}
dispatcher.Subscribe(domain.StoreCreated{}, handler.OnStoreCreated)
// Publish: hər hadisə üçün bütün abunəçilər
```

## 5. Hadisə versiyalanması (Ch 5)

- Köhnə hadisəni DƏYİŞMƏ
- Yeni versiya: `stores.created.v2` yeni payload ilə
- Kompilyator/registry-dən hər ikisi qeyd olunur

## 6. Çatdırılma semantikası (Ch 6)

| Rejim | Tələb |
|---|---|
| at-most-once | itki tolere edilir |
| at-least-once | **dedupe** (inbox) |
| exactly-once | at-least-once + idempotent emal (uydurma deyil) |

```sql
-- inbox dedupe:
INSERT INTO inbox (id) VALUES ($1) ON CONFLICT DO NOTHING;
-- unikal pozuntu = artıq emal olunub → skip
```

## 7. NATS JetStream qurulumu (Ch 6)

```go
nc, _ := nats.Connect(url)
js, _ := nc.JetStream()                      // durable stream
// consumer cursor-lu:Processed mesajlar yerini yadda saxlayır
// Domain → Integration çevirici:
//  domainHandler On(StoreCreated) → pb integration event publish
```

## 8. State transfer keşi (Ch 7)

```go
// hər istehlakçı modulda:
products := postgres.NewProductCacheRepository(
    "search.products_cache", db,
    grpc.NewProductRepository(conn), // sync fallback
)
// upsert-də unikal constraint xətası IGNORE (at-least-once duplicate)
```

## 9. Saga (Ch 8)

```go
saga := sec.NewSaga[CreateOrderData]("create-order", replyTopic)
saga.AddStep().Action(s.authorizeCustomer)
saga.AddStep().
    Action(saga.createOrder).
    CompensateAction(saga.rejectOrder)   // çatarsa geri qaytar
saga.AddStep().Action(saga.confirmOrder).
    CompensateAction(saga.cancelOrder)
```

- Command → Reply (header-də reply destination)
- Saga statusu DB-də: çökmədən davamlı

## 10. Outbox (Ch 9)

```sql
CREATE TABLE outbox (
  id           UUID,
  name         TEXT,
  payload      JSONB,
  published_at TIMESTAMPTZ -- NULL = gözləmədə
);
CREATE INDEX ON outbox (published_at) WHERE published_at IS NULL;
```
```
tx: [domain dəyişikliyi] + [event store yazımı] + [outbox yazımı]  → COMMIT
sonra: processor goroutine → publish → published_at = now()
```

## 11. Test piramidi (Ch 10)

```
go test ./...                       # unit (mock)
go test ./... -tags integration     # DB testləri (testcontainers)
Pact:  consumer gözləntisi → Broker → provider verification
go test ./tests/e2e                # godog (Gherkin)
```

## 12. Deploy (Ch 11)

```bash
docker compose --profile microservices up
terraform init && terraform plan && terraform apply
helm install mallbots ./charts/mallbots
# teardown: helm uninstall → terraform destroy (state-dən bütün resurslar)
```

## 13. Observability (Ch 12)

```go
// OTel avtomatik gRPC interceptor:
grpc.WithUnaryInterceptor(otgrpc.UnaryClientInterceptor())
// span event:
span.AddEvent("before handle", trace.WithAttributes(...))
// Prometheus:
promauto.NewHistogramVec(prometheus.HistogramOpts{
    Name: "received_messages_latency_seconds",
    Buckets: []float64{.01, .025, .05, .1, .25, .5, 1, 2.5, 5},
}, []string{"message"})
```
Bütün mesajlara **correlation ID** header-i ötür.

## 14. Qızıl qaydalar

1. Hadisə adı past tense, məzmun UL əsasında
2. Hər broker publish DB tx-inə bağlı (outbox) — dual write YOX
3. Consumer idempotent yaz (inbox/ON CONFLICT)
4. Sıra lazımdırsa partition açarını aggregate ID-dən götür
5. Kompensasiyasız addım saga-da yazma
6. Hər arxitektura qərarı ADR-də
7. Asinxron API-lər AsyncAPI ilə sənədlənir
