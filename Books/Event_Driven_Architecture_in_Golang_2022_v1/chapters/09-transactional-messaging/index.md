# Chapter 9 — Transactional Messaging (səh. 248-277)

## Bu chapter nədən bəhs edir?

Dual write probleminin həlli: tranzaksiya sərhədləri (DI scoped containers),
event store + mesaj publish-in eyni DB tranzaksiyasında birləşdirilməsi və
Inbox/Outbox pattern-ləri.

## Əsas fikirlər

### 1. Dual write problemi
DB-yə yazma + broker-ə publish = 2 AYRI əməliyyat. Biri uğursuz olsa:
- DB yazıldı, mesaj getmədi → sistem korlanmış vəziyyətdə
- Və ya əksinə → ghost mesajlar

Hər state dəyişikliyi üçün bu problem mövcuddur.

### 2. Həll prinsipi
**State dəyişikliyi + mesaj EYNİ tranzaksiyada.** Event store-a yazma artıq
tranzaksiya daxilindədir; publisher-in də eyni tx-i görməsi lazımdır.

### 3. DI scoped containers
Paketlər iki qrupa bölünür:
- **Singleton** — tətbiq ömrü boyu (DB pool, NATS bağlantısı)
- **Scoped** — hər sorğu/tranzaksiya üçün yeni (tx, handler-lər)

```go
container.Provide("tx", func(c di.Container) interface{} {
    db := c.Get("db").(*sql.DB)
    tx, _ := db.BeginTx(ctx, nil)
    ctx = di.With(ctx, "tx", tx)
    return tx
})
// hər handler: di.Get(ctx, "tx").(*sql.Tx)
```
- `container.Scoped(ctx)` → sorğu başına uşaq konteyner
- Sorğu bitəndə rollback/commit + tx bağlanışı

### 4. Inbox pattern (gələn tərəf)
Məqsəd: mesajın yalnız 1 dəfə emal olunması (at-least-once → effectively-once):
```sql
-- inbox cədvəli:
-- (message_id PRIMARY KEY, received_at, ...)
INSERT INTO inbox (id) VALUES ($1)
-- unikal pozuntu → bu mesaj artıq EMAL OLUNUB → skip
```
Middleware: RawMessageHandler-a tx daxilində inbox yazısı əlavə olunur.

### 5. Outbox pattern (gedən tərəf)
Publish-i 2 mərhələyə böl:
1. **Eyni DB tranzaksiyasında** outbox cədvəlinə yaz (`published_at IS NULL`)
2. Ayrıca **processor goroutine** blok-halında çəkir, broker-ə göndərir,
   `published_at = now()` yazır

```sql
-- periodic olaraq:
SELECT ... FROM depot.outbox WHERE published_at IS NULL;
-- publish → UPDATE depot.outbox SET published_at = now()
```
Partial index: `ON depot.outbox (published_at) WHERE published_at IS NULL`.

### 6. Nəticə
- Event store + outbox eyni tx-də → dual write YOX
- Inbox + idempotent emal → duplicate təsirsiz
- Processor çöksə — yayımlanmamış mesajlar outbox-da gözləyir

## Əsas terminlər

- Dual Write (ikiqat yazma)
- Transactional Boundary (tranzaksiya sərhədi)
- Scoped Container (məhdud konteyner)
- Inbox / Outbox pattern
- Partial Index (qismi indeks)
- Message Processor (mesaj emalçısı)
- Effectively-once delivery

## Praktik nəticə

- Hər broker publish-i DB tx-inə bağla (outbox) — bu EDA-nın mütləq tələbidir
- Gələn mesajları inbox ilə dedupe et
- DI-da singleton/scoped ayrımı: tx hər sorğu üçün təzə
- Outbox processor-i goroutine-də, blok-halında işlət

## Mənbə

Pages: 248-277 (Chapter 9, Event-Driven Architecture in Golang)
