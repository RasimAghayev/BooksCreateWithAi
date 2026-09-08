# Chapter 12 — Реактивное программирование и потоки данных (Reaktiv proqramlaşdırma və data axınları)

## Bu chapter nədən bəhs edir?

Goflow ilə dataflow proqramlaşdırması, Kafka + Sarama (sync/async producer),
Kafka-nın Goflow-a qoşulması və Go-da GraphQL server.

## Əsas fikirlər

### 1. Goflow — dataflow proqramlaşdırma
**Nədir:** Komponentləri kanallarla birləşdirən graph framework — data axını
qara qutu kimi modelləşir.

**Necə işləyir:** Komponentlər input/output kanal sahələri daşıyır;
`Add` + `Connect` + `MapInPort` ilə graph qurulur.

**Kitabdan kod nümunəsi:**
```go
// Komponent — kanal sahələri:
type Upper struct {
    Val <-chan string     // input (yalnız oxu)
    Res chan<- string    // output (yalnız yaz)
}
func (e *Upper) Process() {
    for val := range e.Val {          // input bağlanana qədər
        e.Res <- strings.ToUpper(val)
    }
}

type Printer struct {
    flow.Component       // embed — goflow base
    Line <-chan string
}
func (p *Printer) Process() {
    for line := range p.Line {
        fmt.Println(line)
    }
}

// Graph — komponentlərin birləşdirilməsi (string ID-lərlə):
func NewUpperApp() *goflow.Graph {
    u := goflow.NewGraph()
    u.Add("upper", new(Upper))
    u.Add("printer", new(Printer))
    u.Connect("upper", "Res", "printer", "Line")   // Res → Line
    u.MapInPort("In", "upper", "Val")               // xarici giriş portu
    return u
}

// İşə salma:
net := kafkaflow.NewUpperApp()
in := make(chan string)
net.SetInPort("In", in)
wait := flow.Run(net)                 // wait channel qaytarır
defer func() {
    close(in)                          // input-u bağla → Process bitir
    <-wait                              # graph-ın bitməsini gözlə
}()
in <- "mesaj"                          # axın başlayır
```

**Sub-kod izahı:**
- `u.Connect("ad1", "outChan", "ad2", "inChan")` → string-lərlə bağlantı —
  səhv olsa çox erkən runtime xətası verir
- `close(in)` → range tsiklləri təbii bitir; `<-wait` → təmiz shutdown
- Worker sayı kimi qənnadlar YOXDUR — goflow özü idarə edir

### 2. Kafka + Sarama — sync producer
**Nədir:** Distributiv mesaj növbəsi; sync producer mesajın çatdırılmasını
GÖZLƏYİR (partition/offset qaytarır).

**Kitabdan kod nümunəsi:**
```go
// PRODUCER (sync):
producer, err := sarama.NewSyncProducer([]string{"localhost:9092"}, nil)
defer producer.Close()

func sendMessage(producer sarama.SyncProducer, value string) {
    msg := &sarama.ProducerMessage{
        Topic: "example",
        Value: sarama.StringEncoder(value),
    }
    partition, offset, err := producer.SendMessage(msg)
    if err != nil {
        log.Printf("FAILED to send message: %s\n", err)
        return
    }
    log.Printf("> message sent to partition %d at offset %d\n", partition, offset)
}

// CONSUMER — partition səviyyəsində:
consumer, err := sarama.NewConsumer([]string{"localhost:9092"}, nil)
defer consumer.Close()
partitionConsumer, err := consumer.ConsumePartition("example", 0,
    sarama.OffsetNewest)          // yalnız YENİ mesajlar
defer partitionConsumer.Close()
for {
    msg := <-partitionConsumer.Messages()    // kanaldan mesaj
    log.Printf("Consumed message: \"%s\" at offset: %d\n", msg.Value, msg.Offset)
}
```

**Sub-kod izahı:**
- `SendMessage` → (partition, offset, err) — çatdırılma TƏSDİQ olunur
- `ConsumePartition(topic, 0, OffsetNewest)` → 0-cı partition, yeni mesajlardan
- Sync producer çatışmazlığı: Kafka ölərsə bloklanır → veb handler-lərdə
  timeout riski (sərt asılılıq)

### 3. Kafka — async producer
**Nədir:** Mesajları kanala yazır və dərhal qayıdır; success/error ayrıca
kanallardan emal olunur.

**Kitabdan kod nümunəsi:**
```go
// Config — success/error qaytarışını AKTİVLƏŞDİR:
config := sarama.NewConfig()
config.Producer.Return.Successes = true
config.Producer.Return.Errors = true
producer, err := sarama.NewAsyncProducer([]string{"localhost:9092"}, config)
defer producer.AsyncClose()

// Nəticə emalı — AYRI goroutine (mütləq, yoxsa bloklanır!):
func ProcessResponse(producer sarama.AsyncProducer) {
    for {
        select {
        case result := <-producer.Successes():
            log.Printf("> message: \"%s\" sent to partition %d at offset %d\n",
                result.Value, result.Partition, result.Offset)
        case err := <-producer.Errors():
            log.Println("Failed to produce message", err)
        }
    }
}

// Handler — producer-i struct-da saxla, mesajı kanala yaz:
type KafkaController struct {
    producer sarama.AsyncProducer
}
func (c *KafkaController) Handler(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    msg := r.FormValue("msg")
    if msg == "" {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("msg must be set"))
        return
    }
    c.producer.Input() <- &sarama.ProducerMessage{
        Topic: "example",
        Value: sarama.StringEncoder(msg),
    }
    w.WriteHeader(http.StatusOK)     // gözləmədən dərhal cavab
}

go ProcessResponse(producer)        // arxa plan emalı
c := KafkaController{producer}
http.HandleFunc("/", c.Handler)
```

**Sub-kod izahı:**
- `producer.Input()` → yazma kanalı; handler bloklanmır
- Successes()/Errors() emal edilməzsə → kanallar dolur → BLOKLANMA
- Sync vs async seçimi: təsdiq lazımdırsa sync; throughput üçün async

### 4. Kafka → Goflow birləşməsi
**Nədir:** Kafka consumer-dən gələn mesajlar goflow input kanalına yönəldilir —
növbə + pipeline evliliyi.

**Kitabdan kod nümunəsi:**
```go
// Consumer + goflow birlikdə:
net := kafkaflow.NewUpperApp()          // Upper → Printer graph-ı
in := make(chan string)
net.SetInPort("In", in)
wait := flow.Run(net)
defer func() {
    close(in)
    <-wait
}()
for {
    msg := <-partitionConsumer.Messages()   // Kafka-dan mesaj
    in <- string(msg.Value)                  // goflow-a ötür
}
// Producer "Message 0" göndərir → consumer "MESSAGE 0" çap edir
```

**Sub-kod izahı:**
- Kafka axını = goflow üçün təbii mənbə — komponentlər modul/veyenidən
  istifadə olunur (eyni Printer müxtəlif graph-larda)
- İstənilən producer eyni topic-a yazırsa — hamısı eyni axına birləşir

### 5. GraphQL server (graphql-go)
**Nədir:** REST alternativi — client schema üzərindən istədiyi sahələri soruşur.

**Necə işləyir:** Type → Resolve funksiyası → Schema → graphql.Do.

**Kitabdan kod nümunəsi:**
```go
// Model:
type Card struct {
    Value string
    Suit  string
}

// GraphQL Object tipi — sahə + resolve:
func CardType() *graphql.Object {
    return graphql.NewObject(graphql.ObjectConfig{
        Name: "Card",
        Fields: graphql.Fields{
            "value": &graphql.Field{
                Type: graphql.String,
                Resolve: func(p graphql.ResolveParams) (interface{}, error) {
                    if card, ok := p.Source.(Card); ok {   // type assertion
                        return card.Value, nil
                    }
                    return nil, nil
                },
            },
            "suit": &graphql.Field{ /* analoog */ },
        },
    })
}

// Root query resolve — filtrləmə:
func Resolve(p graphql.ResolveParams) (interface{}, error) {
    finalCards := []Card{}
    suit, suitOK := p.Args["suit"].(string)     // arqumentlər
    value, valueOK := p.Args["value"].(string)
    for _, card := range cards {
        if suitOK && !strings.EqualFold(card.Suit, suit) { continue }
        if valueOK && !strings.EqualFold(card.Value, value) { continue }
        finalCards = append(finalCards, card)
    }
    return finalCards, nil
}

// Schema — root query + arqumentlər:
func Setup() (graphql.Schema, error) {
    fields := graphql.Fields{
        "cards": &graphql.Field{
            Type: graphql.NewList(CardType()),
            Args: graphql.FieldConfigArgument{
                "suit":  &graphql.ArgumentConfig{Type: graphql.String},
                "value": &graphql.ArgumentConfig{Type: graphql.String},
            },
            Resolve: Resolve,
        },
    }
    rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
    return graphql.NewSchema(graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)})
}

// Sorğu icrası:
query := `
{
    cards(value: "A"){
        value
        suit
    }
}`
params := graphql.Params{Schema: schema, RequestString: query}
r := graphql.Do(params)
if len(r.Errors) > 0 {
    log.Fatalf("failed: %+v", r.Errors)
}
rJSON, _ := json.MarshalIndent(r, "", "  ")
// {"data": {"cards": [{"suit":"Spades","value":"A"}, ...Hearts/Clubs/Diamonds]}}
```

**Sub-kod izahı:**
- `p.Source.(Card)` → cari obyekt; `p.Args["suit"]` → sorğu arqumenti
- Client YALNIZ lazımi sahələri soruşur (value/suit) — overfetching yox
- Resolve funksiyasını DB sorğusuna dəyişmək kifayətdir — schema eyni qalır
- Bu nümunə query-ni birbaşa icra edir; REST endpoint-ə qoşmaq kiçik işdir

## Əsas terminlər

- Reactive Programming (reaktiv proqramlaşdırma)
- Dataflow / Graph (data axını / qraf)
- Goflow: Component / Connect / MapInPort / Run
- Kafka: Topic / Partition / Offset
- Sarama: SyncProducer / AsyncProducer / ConsumePartition
- Input/Successes/Errors kanalları
- GraphQL: Schema / Object / Field / Resolve / Args / RootQuery
- Overfetching (artıq data çəkilməsi)

## Praktik nəticə

- Goflow: komponentlər kiçik və yenidən istifadə oluna bilən; close(in) +
  <-wait ilə təmiz shutdown
- Sync producer: çatdırılma təsdiqi lazım olanda; Kafka ölərsə bloklanma riski
- Async producer: handler bloklanmır; Successes/Errors ayrıca goroutine-də
  emal etmək MÜTLƏQDİR
- Kafka + pipeline: növbəni mənbə kimi istifadə et — komponentləri müxtəlif
  graph-larda paylaş
- GraphQL: type + resolve + schema üçlüyü; Resolve-i dəyişərək DB-yə keç

## Mənbə

Pages: 396-424 (Chapter 12, Go Programming Cookbook 2nd ed)
