# Chapter 6 — Asynchronous Communication (səh. 115-132)

## Bu fəsil nədən bəhs edir?

Asinxron kommunikasiya modeli, message broker/publisher-subscriber naxışları,
Apache Kafka əsasları və rating servisinə Kafka ingestion-un əlavə edilməsi.

## Əsas fikirlər

### 1. Asinxron kommunikasiya nədir?
- Göndərən dərhal cavab GÖZLƏMİR (cavab gecikə bilər, ya heç gəlməyə bilər)
- Bənzətmə: telefon = sinxron, poçt = asinxron
- Asinxron yavaş deyil — əksinə, kontekst dəyişmədiyi üçün tez-tez DAHA SÜRƏTLİDİR

**Faydaları:**
1. Streamlined emal — cavab gözləmədən növbəti mesajı göndər
2. Göndərmə/emal DE-COUPLE — video emal kimi uzun tasks (ACK → sonra nəticə)
3. Load balancing — spike-lərdə server boğulmur, növbə ilə emal edir

**Uyğun məsələlər:** long-running tasks; publish-once/consume-many (status
hesabatları); sequential/batch emal (HDD ardıcıl yazma)

**Çətinliklər:** mürəkkəb xəta idarəsi (çatmadı mı, itdi mi?); əlavə
infrastruktur (broker); intuitiv olmayan data flow.

### 2. Message Broker (Mesaj Vasitəçisi)
Rollar: delivery + transformation + aggregation + routing
Poçt ofisi kimi — göndərən çatdırımla maraqlanmir.

**Delivery guarantees (Çatdırılmə Zəmanətləri):**
| Zəmanət | Mənası | Xarakter |
|---|---|---|
| At-most-once | 0 və ya 1 dəfə | lossy (itirə bilər), sürətli |
| At-least-once | ≥1 dəfə (retry) | lossless, dubl mümkün |
| Exactly-once | dəqiq 1 dəfə | lossless, ən çətini (metadata/dedupe lazım) |

### 3. Publisher-Subscriber (Pub/Sub) modeli
- Hər komponent publish edə və abunə ola bilər (Twitter bənzətməsi)
- Nümunə: profil silinməsi event-i → BÜTÜN servislər ayrı-ayrı xəbərdar
  edilmək əvəzinə bir event-i consume edir və arxivləşdirir
- Aşağı mesaj tezliyində də faydalı (düzəldilmiş abunəlik)

### 4. Apache Kafka
- LinkedIn mənşəli, ən populyar open-source broker
- **Data model:** producer → **topic** (ardıcıl mesajlar, unikal **offset**) →
  consumer; topic-lər **partition**-lana bilər (paralel emal)
- Güclü tərəfləri: yüksək throughput (ardıcıl I/O), partitioning ilə
  scalability, konfiqurə olunan durability (retention: 7 gün / sonsuz)
- Mürəkkəb infrastruktur — bu fəsildə Docker versiyası

### 5. Rating ingestion layihəsi
Ssenari: IMDB kimi data provider rating-ləri dərc edir, rating servisi
consume edir (pub/sub).

**Model:**
```go
type RatingEvent struct {
    UserID, RecordID, RecordType, Value, EventType  // "put" | "delete"
}
```

**Producer (cmd/ratingingester):**
```go
producer, _ := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost"})
// fayldan JSON oxu → hər event-i JSON encode → p.Produce(&kafka.Message{
//   TopicPartition: {Topic: "ratings", Partition: kafka.PartitionAny}})
producer.Flush(10000)  // hamısinin göndərilməsini gözlə
```
- Kitabxana: `confluent-kafka-go` (alternativ: Shopify/sarama)

**Consumer (rating/internal/ingester/kafka):**
```go
consumer, _ := kafka.NewConsumer(&kafka.ConfigMap{
    "bootstrap.servers": addr, "group.id": groupID,
    "auto.offset.reset": "earliest"})
// Ingest(): SubscribeTopics + goroutine → ReadMessage(-1) → JSON unmarshal → channel
```
- `ReadMessage(-1)`: -1 = öz mövqedən yox, həmişə mövzunun BAŞINDAN oxu
- Channel + `for e := range ch` — connection bağlı yekunlaşır

**Controller:** `StartIngestion` — channel-dən gələn hər event-i `PutRating`
ilə DB-yə yazır; servis həm sinxron API, həm asinxron ingester təqdim edir.

### 6. Best practices
**Versioning (Versiyalama):** event formatına `version` sahəsi əlavə et —
consumer köhnə/yeni formatları ayırd edə bilsin (versiya-specific validation,
bəzi versiyaları ignore etmək).

**Partitioning (Bölmələmə):** partition-ı manual seçmək = **data locality**
(məs. userID ilə partition → hər istifadəçinin datası tək partitionda,
axtarış sadələşir).

## Termindirmə (AZ)
- Asynchronous Communication — Asinxron Rabitə
- Message Broker — Mesaj Vasitəçisi
- Publisher-Subscriber — Nəşr edən-Abunə olan modeli
- Delivery Guarantee — Çatdırılmə Zəmanəti
- Topic — Mövzu (mesajların ardıcıl saxlandığı kanal)
- Offset — Ofset (mesajın mövzu daxilindəki mövqeyi)
- Partition — Bölmə (paralel emal üçün mövzunun hissəsi)
- Data Locality — Data Yerliliyi

## Kviz sualları
1. Exactly-once zəmanəti niyə çətindir? (Heç vaxt təkrar göndərilməməsi üçün
   əlavə yoxlama/metadata saxlanması tələb olunur)
2. ReadMessage(-1) nə edir? (Mövzunun ən başından — bütün mövcud mesajları oxuyur)
3. Versioning nə üçün lazımdır? (Consumer format dəyişikliyini ayırd etsin —
   sahənin olmaması "yeni formatdırmı?" sualını qarışdırmır)
4. Partition-ı userID ilə seçməyin faydası? (Data locality — bir istifadəçinin
   bütün mesajları eyni partitionda saxlanılır)
