# Chapter 19 — Kafka (səh. 433-462)

## Bu fəsil nədən bəhs edir?

Apache Kafka event streaming platforması: anlayışlar (broker, event, topic,
partition, offset, commit), Confluent klienti (librdkafka wrapper — sync/async
producer/consumer), Segmentio klienti (tam Go — Connection, Writer/Reader) və
Kafka REST Proxy (adi HTTP ilə produce/consume).

## Əsas fikirlər

### 1. Kafka anlayışları
**Nədir:** Paylanmış publish/subscribe platforması.

| Anlayış | Təsvir |
|---|---|
| Broker | klientlərin bağlandığı server |
| Event (hadisə) | key + value + timestamp (+metadata) yazısı |
| Topic | hadisələrin məntiqi qrupu (qovluq kimi); retention siyasəti ilə saxlanır |
| Partition | topic-in brokerlər arası paylanması; replikasiya (adətən 3) fault-tolerance verir |
| Offset | partition daxilində hadisənin ardıcıl indeksi |
| Commit | consumerın son istehlak etdiyi offset-in brokerə bildirilməsi |

**Qaydalar:** eyni key-li hadisələr eyni partition-da; yazılış sırası ilə
istehlak zəmanəti; auto-commit default, manual mümkün.

### 2. Confluent klienti — sync producer
**Nədir:** librdkafka C kitabxanasının Go wrapper-i.

**Message tipi:**
```go
type Message struct {
    TopicPartition TopicPartition
    Value          []byte     // mesaj məzmunu
    Key            []byte     // partition idarəsi üçün
    Timestamp      time.Time
    Opaque         interface{}
    Headers        []Header
}
```

**Kitabdan kod nümunəsi:**
```go
cfg := kafka.ConfigMap{"bootstrap.servers": "localhost:9092"}
p, err := kafka.NewProducer(&cfg)
if err != nil {
    panic(err)
}
defer p.Close()

deliveryChan := make(chan kafka.Event, 10)
defer close(deliveryChan)
topic := "helloTopic"
msgs := []string{"Save", "the", "world", "with", "Go!!!"}
for _, word := range msgs {
    err = p.Produce(&kafka.Message{
        TopicPartition: kafka.TopicPartition{
            Topic: &topic, Partition: kafka.PartitionAny},
        Value: []byte(word),
    }, deliveryChan)
    if err != nil {
        panic(err)
    }
    e := <-deliveryChan          // hər mesajın çatdırılmasını gözlə (sync)
    m := e.(*kafka.Message)
    fmt.Printf("Sent %v\n", m)   // Sent helloTopic[0]@30
}
// p.Flush(1000) — qalan mesajları gözlə (ms)
```
- `PartitionAny` → Kafka partition-u özü seçir
- Sync gözləmə throughput-u azaldır → async variant aşağıda

### 3. Confluent klienti — consumer
**Kitabdan kod nümunəsi:**
```go
c, err := kafka.NewConsumer(&kafka.ConfigMap{
    "bootstrap.servers": "localhost:9092",
    "group.id":          "helloGroup",     // consumer qrupu
    "auto.offset.reset": "earliest",       // əvvəldən oxu
})
if err != nil {
    panic(err)
}
defer c.Close()

c.Subscribe("helloTopic", nil)
for {
    ev := c.Poll(1000)                      // 1s gözlə → Event və ya nil
    switch e := ev.(type) {
    case *kafka.Message:
        c.Commit()                           // manual commit
        fmt.Printf("Msg on %s: %s\n", e.TopicPartition, string(e.Value))
    case kafka.PartitionEOF:
        fmt.Printf("%v\n", e)
    case kafka.Error:
        fmt.Printf("Error: %v\n", e)
        break
    default:
        fmt.Printf("Ignored: %v\n", e)
    }
}
```
- `enable.auto.commit=false` → manual `Commit()` zərurəti

### 4. Async producer (goroutine handler)
```go
func Handler(c chan kafka.Event) {
    for {
        e := <-c
        if e == nil {
            return
        }
        m := e.(*kafka.Message)
        if m.TopicPartition.Error != nil {
            fmt.Printf("Partition error %s\n", m.TopicPartition.Error)
        } else {
            fmt.Printf("Sent %v: %s\n", m, string(m.Value))
        }
    }
}

delivery_chan := make(chan kafka.Event, 10)
defer close(delivery_chan)
go Handler(delivery_chan)          // bloklamadan çatdırılma izləməsi

for _, word := range msgs {
    p.Produce(&kafka.Message{...}, delivery_chan)   // gözləmirsiz
}
p.Flush(10000)                      // hamısı göndərilənə qədər blokla
```

### 5. Async consumer — batch commit
Hər mesajı commit etmək overhead-dir; batch commit daha yaxşıdır:

```go
const COMMIT_N = 2

func Committer(c *kafka.Consumer) {
    offsets, err := c.Commit()
    if err != nil {
        fmt.Printf("Error: %s\n", err)
        return
    }
    fmt.Printf("Offset: %#v\n", offsets[0].String())
}

// main loop daxilində:
case *kafka.Message:
    counter += 1
    if counter % COMMIT_N == 0 {
        go Committer(c)      // hər 2 mesajda bir commit
    }
    fmt.Printf("Msg on %s: %s\n", e.TopicPartition, string(e.Value))
```

**Alternativ — commit uğurundan sonra emal:**
```go
case *kafka.Message:
    append(messages, msg)
    counter += 1
    if counter % COMMIT_N == 0 {
        offset, err := c.Commit()
        if err == nil {
            go ProcessMsgs(messages)   // yalnız commit-lə bərabər emal
        }
    }
```

### 6. Segmentio klienti (tam Go) — Connection
**Nədir:** C asılılığı yoxdur; sadə Connection əsaslı API.

**Producer:**
```go
conn, err := kafka.DialLeader(context.Background(),
    "tcp", "localhost:9092", "helloTopic", 0)   // host, topic, partition
if err != nil {
    panic(err)
}
defer conn.Close()

conn.SetWriteDeadline(time.Now().Add(3 * time.Second))
for _, m := range []string{"Save", "the", "world", "with", "Go!!!"} {
    l, err := conn.WriteMessages(kafka.Message{Value: []byte(m)})
    if err != nil {
        panic(err)
    }
    fmt.Printf("Sent %d bytes: %s\n", l, m)
}
```
- `WriteMessages(...Message)` — variadic; hamısını bir çağırışla ötürmək
  throughput-u artırır; `Partition` sahəsi read-only (yazmaq olmaz)

**Consumer (ReadBatch):**
```go
conn.SetReadDeadline(time.Now().Add(time.Second))
batch := conn.ReadBatch(10e3, 10e6)   // min 10KB, max 10MB
defer batch.Close()
for {
    b := make([]byte, 10e3)
    l, err := batch.Read(b)
    if err != nil {
        break
    }
    fmt.Printf("Received %d: %s\n", l, string(b))
}
```

### 7. Segmentio — Writer/Reader (yüksək səviyyə)
```go
// Writer — producer:
w := kafka.Writer{
    Addr:  kafka.TCP("localhost:9092"),
    Topic: "helloTopic",
}
defer w.Close()
err = w.WriteMessages(context.Background(),
    kafka.Message{Value: []byte("Save")},
    kafka.Message{Value: []byte("the")})
w.Stats()   // yazılan bayt statistikası və s.

// Reader — consumer:
r := kafka.NewReader(kafka.ReaderConfig{
    Brokers:  []string{"localhost:9092"},
    Partition: 0,
    Topic:    "helloTopic",
    GroupID:  "testGroup",
    MinBytes: 10e3,
    MaxBytes: 10e6,
})
defer r.Close()

for {
    m, err := r.FetchMessage(context.Background())  // bloklanır
    if err != nil {
        break
    }
    if err := r.CommitMessages(context.Background(), m); err != nil {
        panic(err)          // variadic — batch commit mümkündür
    }
    fmt.Printf("Topic %s msg: %s\n", m.Topic, m.Value)
}
```

### 8. Kafka REST Proxy — producer (adi HTTP)
**Nədir:** TCP protokolu əlçatmaz olan mühitlər (brauzer və s.) üçün
 HTTP API. Content-Type: `application/vnd.kafka.json.v2+json`.

```go
const (
    URL         = "http://localhost:8082/topics/%s/partitions/%d"
    CONTENT_TYPE = "application/vnd.kafka.json.v2+json"
)

func BuildBody(users []User) string {     // {"records":[{"value":{...}},...]}
    values := make([]string, len(users))
    for i, u := range users {
        encoded, _ := json.Marshal(&u)
        values[i] = fmt.Sprintf(`{"value":%s}`, encoded)
    }
    return fmt.Sprintf(`{"records": [%s]}`, strings.Join(values, ","))
}

body := BuildBody(users)
resp, err := http.Post(fmt.Sprintf(URL, "helloTopic", 0), CONTENT_TYPE,
    bytes.NewBuffer([]byte(body)))
defer resp.Body.Close()
// Cavab: {"offsets":[{"partition":0,"offset":165,...},...]}
```

### 9. Kafka REST Proxy — consumer (4 addım)
**Addımlar:** 1) consumer yarat; 2) subscribe; 3) records oxu; 4) instance sil.

```go
const (
    HOST              = "http://localhost:8082"
    GROUP             = "testGroup"
    CONSUMER          = "testConsumer"
    NEW_CONSUMER      = "%s/consumers/%s"
    SUBSCRIBE_CONSUMER = "%s/consumers/%s/instances/%s/subscription"
    FETCH_CONSUMER    = "%s/consumers/%s/instances/%s/records"
    DELETE_CONSUMER   = "%s/consumers/%s/instances/%s"
)

// 1) Yeni consumer (POST):
url := fmt.Sprintf(NEW_CONSUMER, HOST, GROUP)
body := fmt.Sprintf(`{"name":"%s", "format": "json"}`, CONSUMER)
DoHelper(&client, url, []byte(body))
// 200 OK: {"instance_id":"testConsumer","base_uri":"..."}

// 2) Subscribe (POST):
url = fmt.Sprintf(SUBSCRIBE_CONSUMER, HOST, GROUP, CONSUMER)
body = `{"topics":["helloTopic"]}`
DoHelper(&client, url, []byte(body))    // 204 No Content

// 3) Records (GET — Accept header mütləq!):
req, _ := http.NewRequest(http.MethodGet,
    fmt.Sprintf(FETCH_CONSUMER, HOST, GROUP, CONSUMER), nil)
req.Header.Add("Accept", CONTENT_TYPE)
respRecords, _ := client.Do(req)
// 200 OK: [{"topic":"helloTopic","value":{...},"partition":0,"offset":179}, ...]

// 4) Sil (DELETE):
deleteReq, _ := http.NewRequest(http.MethodDelete,
    fmt.Sprintf(DELETE_CONSUMER, HOST, GROUP, CONSUMER), nil)
client.Do(deleteReq)                     // 204
```

**Offset irəliletmə (növbəti batch üçün):**
```json
POST /consumers/testGroup/instances/testConsumer/positions
{
  "offsets": [
    {"topic": "helloTopic", "partition": 0, "offset": 181}
  ]
}
```
- Sorğu parametrləri: `timeout` (fetch müddəti), `max_bytes`

## Əsas terminlər
- Publish/subscribe — yayın/abunə mesajlaşma modeli
- Broker — Kafka server
- Event/Record — key+value+timestamp mesajı
- Topic — hadisələrin kanalı
- Partition — topic hissəsi; replikasiya vahidi
- Offset — partition daxili hadisə indeksi
- Consumer group — birlikdə istehlak edən consumer dəsti
- Commit — offset-in brokerə bildirilməsi
- librdkafka — Confluent-in C nüvəsi
- REST Proxy — Kafka HTTP interfeysi

## Praktik nəticə
Client seçimi: C asılılığı problem deyilsə Confluent (rəsmi, Playground zəngin),
yalnız Go lazımdırsa Segmentio (DialLeader/Connection və ya Writer/Reader
yüksək səviyyə API), HTTP-dən başqa yol yoxdursa REST Proxy. Performans
qaydaları: sync mesaj gözləməkdən qaçın (async handler + Flush), commit-i
batch-lə (hər N mesajda) edin, WriteMessages/FetchMessages-ə kolleksiya ötürün.
Consumer axını: Subscribe → Poll/Fetch loop → type switch → batch commit.

## Mənbə
Pages: 433-462 (PDF 433-462)
