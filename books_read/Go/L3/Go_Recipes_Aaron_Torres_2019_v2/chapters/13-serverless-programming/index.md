# Chapter 13 — Бессерверное программирование (Serverless proqramlaşdırma)

## Bu chapter nədən bəhs edir?

AWS Lambda (Apex ilə deploy), Apex logging/metrics, Google App Engine
(Datastore) və Firebase (Firestore) — serverless platformalarda Go.

## Əsas fikirlər

### 1. AWS Lambda + Apex
**Nədir:** Server runtime idarə etmədən funksiya işlətmə; Apex — Lambda
funksiyalarının build/deploy/idarə aləti.

**Kitabdan kod nümunəsi:**
```go
// Lambda funksiyası — aws-lambda-go:
type Message struct {
    Name string `json:"name"`
}
type Response struct {
    Greeting string `json:"greeting"`
}

// HandleRequest — çağırılan funksiya:
func HandleRequest(ctx context.Context, m Message) (Response, error) {
    return Response{Greeting: fmt.Sprintf("Hello, %s", m.Name)}, nil
}
func main() {
    lambda.Start(HandleRequest)     // handler-i qeydiyyatdan keçir
}
```
```bash
apex init                    # project + IAM rol/policy yaradır
apex deploy                  # funksiyaları Lambda-ya yükləyir (v1, v2...)
echo '{"name": "Reader"}' | apex invoke greeter1
# {"greeting":"Hello, Reader"}
apex logs greeter2           # CloudWatch logları
apex rollback greeter        # versiya geri qaytarma
```

**Handler imza qaydaları (Apex/aws-lambda-go):**
- 0-2 arqument; 2-liyin birincisi `context.Context` olmalıdır
- 0-2 qaytarma; 2-liyin ikincisi `error`; 1-li return yalnız `error`

**Sub-kod izahı:**
- `lambda.Start(handler)` → main-in yeganə vəzifəsi
- Apex versiyalama: funksiya dəyişsə YENİ versiya (v2), yeni ad = yeni funksiya
- JSON input/output avtomatik marshal/unmarshal

### 2. Apex logging və metrics
**Nədir:** Lambda-da strukturlaşdırılmış logging (apex/log) + `apex metrics`
ilə xərc/çağırış statistikası.

**Kitabdan kod nümunəsi:**
```go
func HandleRequest(ctx context.Context, input Input) (string, error) {
    log.SetHandler(text.New(os.Stderr))     // stderr → CloudWatch
    log.WithField("secret", input.Secret).Info("secret guessed")
    if input.Secret == "klaatu barada nikto" {
        return "secret guessed!", nil
    }
    return "try again", nil
}
```
```bash
apex logs secret
# INFO[0000] secret guessed secret=open sesame
# REPORT Duration: 52.23 ms Billed: 100 ms Memory: 128 MB Max Used: 19 MB
apex metrics secret
# total cost / invocations / duration / throttles / errors / memory
apex delete       # təmizləmə — PUL ÖDƏMƏMƏK ÜÇÜN MÜTLƏQ
```

**Sub-kod izahı:**
- Logger stderr-ə yazır → CloudWatch-da görünür
- REPORT sətri: real müddət, hesablanan müddət (100ms minimum), yaddaş
- Apex metrics xərc də göstərir; AWS Console-da da mövcuddur

### 3. Google App Engine
**Nədir:** Yüklə avtomatik scale olunan veb app platforması; Datastore
(MongoDB-ə bənzər NoSQL) ilə birlikdə.

**Kitabdan kod nümunəsi:**
```yaml
# app.yaml:
runtime: go112
manual_scaling:
  instances: 1
env_variables:
  GCLOUD_DATASET_ID: go-cookbook
```
```go
// Datastore model + CRUD:
type Message struct {
    Timestamp time.Time
    Message   string
}
func (c *Controller) storeMessage(ctx context.Context, message string) error {
    m := &Message{Timestamp: time.Now(), Message: message}
    k := datastore.IncompleteKey("Message", nil)   // avtomatik ID
    _, err := c.store.Put(ctx, k, m)
    return err
}
func (c *Controller) queryMessages(ctx context.Context, limit int) ([]*Message, error) {
    q := datastore.NewQuery("Message").
        Order("-Timestamp").      // ən yeni birincil
        Limit(limit)
    messages := make([]*Message, 0)
    _, err := c.store.GetAll(ctx, q, &messages)
    return messages, err
}

// Handler — guestbook kimi:
func (c *Controller) handle(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "invalid method", http.StatusMethodNotAllowed)
        return
    }
    ctx := context.Background()
    r.ParseForm()
    if message := r.FormValue("message"); message != "" {
        c.storeMessage(ctx, message)
    }
    fmt.Fprintln(w, "Messages:")
    messages, _ := c.queryMessages(ctx, 10)
    for _, message := range messages {
        fmt.Fprintln(w, message.Message)
    }
}

// main — lokal və App Engine-də eyni:
func main() {
    ctx := context.Background()
    projectID := os.Getenv("GCLOUD_DATASET_ID")
    datastoreClient, err := datastore.NewClient(ctx, projectID)
    c := Controller{datastoreClient}
    http.HandleFunc("/", c.handle)
    port := os.Getenv("PORT")           # App Engine PORT verir
    if port == "" { port = "8080" }
    http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
```
```bash
gcloud config set project go-cookbook
gcloud auth application-default login
gcloud app deploy      # produksiyaya
gcloud app browse
```

**Sub-kod izahı:**
- `datastore.IncompleteKey` → datastore avtomatik key yaradır
- `Order("-Timestamp")` → desc sıralama; Limit ilə birlikdə
- Port env-dən — App Engine inject edir; lokal default 8080
- Yeni mesaj dərhal görünməyə bilər (eventual consistency)

### 4. Firebase (Firestore)
**Nədir:** Mobil-üstün realtime NoSQL; service account faylı ilə auth;
Get/Set interfeysi ilə wrap.

**Kitabdan kod nümunəsi:**
```go
// Client interfeysi — mock üçün:
type Client interface {
    Get(ctx context.Context, key string) (interface{}, error)
    Set(ctx context.Context, key string, value interface{}) error
    Close() error
}

// Firestore realizasiyası:
type firebaseClient struct {
    *firestore.Client
    collection string
}
func (f *firebaseClient) Get(ctx context.Context, key string) (interface{}, error) {
    data, err := f.Collection(f.collection).Doc(key).Get(ctx)
    if err != nil {
        return nil, errors.Wrap(err, "get failed")
    }
    return data.Data(), nil
}
func (f *firebaseClient) Set(ctx context.Context, key string, value interface{}) error {
    set := make(map[string]interface{})
    set[key] = value
    _, err := f.Collection(f.collection).Doc(key).Set(ctx, set)
    return errors.Wrap(err, "set failed")
}

// Auth — service_account.json ilə:
func Authenticate(ctx context.Context, collection string) (Client, error) {
    opt := option.WithCredentialsFile("/tmp/service_account.json")
    app, err := firebase.NewApp(ctx, nil, opt)
    if err != nil {
        return nil, errors.Wrap(err, "error initializing app")
    }
    client, err := app.Firestore(ctx)
    return &firebaseClient{Client: client, collection: collection}, nil
}

// İstifadə:
c, _ := firebase.Authenticate(ctx, "collection")
defer c.Close()
c.Set(ctx, "key", []string{"val1", "val2"})
res, _ := c.Get(ctx, "key")      // [val1 val2]
```

**Sub-kod izahı:**
- Collection → Doc (key) → Get/Set — document modeli
- Service account tokeni adminsdk konsolundan yüklənir
- Client interfeysi → test mock; Close defer ilə bağlanır
- map[string]interface{} saxlanır — veb/mobil client-lər oxuyur

## Əsas terminlər

- Serverless (serversiz arxitektura)
- AWS Lambda / Handler imza qaydaları
- Apex (deploy/invoke/logs/metrics/rollback/delete)
- IAM Role/Policy
- CloudWatch (Lambda logları)
- Billed Duration (hesablanan müddət)
- Google App Engine / manual_scaling
- Cloud Datastore (IncompleteKey/Query/Order/Limit)
- Firebase / Firestore / Collection / Document
- Service Account (service_account.json)
- Eventual Consistency (mütləq ardıcıllıq olmayan)

## Praktik nəticə

- Lambda handler: ctx + custom type → custom type/error; JSON avtomatik
- Serverless PUL ödəyir → iş bitəndə apex delete / console təmizləməsi ŞƏRT
- App Engine: PORT və DATASET_ID env-dən; lokal və cloud eyni kod
- Firebase: interfeys ilə wrap → mock-test; service account faylına ehtiyat
- Datastore sorğuları Order+Limit ilə; yeni yazı dərhal görünməyə bilər

## Mənbə

Pages: 425-445 (Chapter 13, Go Programming Cookbook 2nd ed)
