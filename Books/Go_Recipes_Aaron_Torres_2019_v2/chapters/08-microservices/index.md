# Chapter 8 — Микросервисы для приложений в Go (Go mikroservisləri)

## Bu chapter nədən bəhs edir?

Handler/HandlerFunc/ResponseWriter əsasları, state-li handler-lər (struct +
closure), validasiya, kontent negotiasiya (unrolled/render), middleware,
reverse proxy və gRPC-ın JSON API kimi ixracı.

## Əsas fikirlər

### 1. Handler, Request və ResponseWriter
**Nədir:** Go vebinin əsas vahidi — `func(ResponseWriter, *http.Request)`.

**Kitabdan kod nümunəsi:**
```go
type HandlerFunc func(http.ResponseWriter, *http.Request)
type Handler interface {
    ServeHTTP(http.ResponseWriter, *http.Request)
}

// GET handler:
func HelloHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/plain")     // header ƏVVƏL
    if r.Method != http.MethodGet {
        w.WriteHeader(http.StatusMethodNotAllowed)     // metod yoxlaması
        return
    }
    name := r.URL.Query().Get("name")                 // query param
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(fmt.Sprintf("Hello %s!", name)))
}

// POST handler — JSON cavab:
type GreetingResponse struct {
    Payload struct {
        Greeting string `json:"greeting,omitempty"`
        Name     string `json:"name,omitempty"`
        Error    string `json:"error,omitempty"`
    } `json:"payload"`
    Successful bool `json:"successful"`
}
func GreetingHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    if r.Method != http.MethodPost {
        w.WriteHeader(http.StatusMethodNotAllowed)
        return
    }
    if err := r.ParseForm(); err != nil { /* 400/500 */ }
    name := r.FormValue("name")          // form data
    w.WriteHeader(http.StatusOK)
    gr.Successful = true
    payload, _ := json.Marshal(gr)
    w.Write(payload)
}

// Qoşulma:
http.HandleFunc("/name", handlers.HelloHandler)
http.ListenAndServe(":3333", nil)
```

**Sub-kod izahı:**
- Header-lər WriteHeader-dan ƏVVƏL set olunmalı
- `r.URL.Query().Get` → GET param; `r.ParseForm()` + `r.FormValue` → POST form
- JSON üçün `json.NewDecoder(r.Body).Decode(&p)` alternativdir

### 2. State-li handler-lər (struct və closure)
**Nədir:** Handler funksiyasının imzası state saxlamır — struct (DI) və ya
closure ilə state ötürülür.

**Kitabdan kod nümunəsi:**
```go
// CONTROLLER — struct ilə state (Storage interfeysi injeksiyası):
type Controller struct {
    storage Storage      // interfeys — test/mock üçün
}
func New(storage Storage) *Controller { return &Controller{storage: storage} }

type Storage interface {
    Get() string
    Put(string)
}
type MemStorage struct{ value string }
func (m *MemStorage) Get() string { return m.value }
func (m *MemStorage) Put(s string) { m.value = s }

// Sadə metod-handler:
func (c *Controller) SetValue(w http.ResponseWriter, r *http.Request) {
    // ... c.storage.Put(value)
}

// CLOSURE-handler — arqumentlə qaytarılan HandlerFunc:
func (c *Controller) GetValue(UseDefault bool) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        value := "default"
        if !UseDefault {
            value = c.storage.Get()
        }
        // ...
    }
}

// İki route — eyni handler, fərqli konfiqurasiya:
http.HandleFunc("/get", c.GetValue(false))
http.HandleFunc("/get/default", c.GetValue(true))
```

**Sub-kod izahı:**
- Struct → bir neçə handler arasında paylaşılan state (DB, logger, storage)
- Closure → eyni handler-in parametrlə fərqli variantları
- Metod dəyəri `c.SetValue` HandlerFunc imzasını qane edir

### 3. Input validasiyası (closure ilə mock-able)
**Nədir:** Validasiya funksiyası struct sahəsi kimi — test vaxtı əvaz edilə
bilir.

**Kitabdan kod nümunəsi:**
```go
type Controller struct {
    ValidatePayload func(p *Payload) error   // funksiya sahəsi!
}
func New() *Controller {
    return &Controller{ValidatePayload: ValidatePayload}  // default impl
}

type Verror struct{ error }      // typed error — istifadəçiyə göstərilə bilən
func ValidatePayload(p *Payload) error {
    if p.Name == "" {
        return Verror{errors.New("name is required")}
    }
    if p.Age <= 0 || p.Age >= 120 {
        return Verror{errors.New("age is required and must be >0 and <120")}
    }
    return nil
}

// Handler — JSON decode + validasiya + typed error switch:
func (c *Controller) Process(w http.ResponseWriter, r *http.Request) {
    decoder := json.NewDecoder(r.Body)
    defer r.Body.Close()
    var p Payload
    if err := decoder.Decode(&p); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return
    }
    if err := c.ValidatePayload(&p); err != nil {
        switch err.(type) {
        case Verror:                              // validasiya xətası → 400
            w.WriteHeader(http.StatusBadRequest)
            w.Write([]byte(err.Error()))          // mesaj istifadəçiyə
            return
        default:                                  // daxili xəta → 500
            w.WriteHeader(http.StatusInternalServerError)
            return
        }
    }
    w.WriteHeader(http.StatusOK)
}
```

**Sub-kod izahı:**
- `ValidatePayload` funksiya sahəsi → test-də sadə stub-la əvəz olunur
- Verror tipi → validasiya xətaları istifadəçiyə göstərilə bilən xətalardır
- Çoxsaylı validasiya xətası → Verror-a map/slice əlavə etmək olar

### 4. Render və Content Negotiation
**Nədir:** Content-Type header-inə görə cavab formatı seçilir (JSON/XML).

**Kitabdan kod nümunəsi:**
```go
// unrolled/render wrapper:
type Negotiator struct {
    ContentType string
    *render.Render
}
func GetNegotiator(r *http.Request) *Negotiator {
    contentType := r.Header.Get("Content-Type")
    return &Negotiator{ContentType: contentType, Render: render.New()}
}

func (n *Negotiator) Respond(w io.Writer, status int, v interface{}) {
    switch n.ContentType {
    case render.ContentJSON:
        n.Render.JSON(w, status, v)
    case render.ContentXML:
        n.Render.XML(w, status, v)
    default:
        n.Render.JSON(w, status, v)
    }
}

// Hər iki formatda eyni model:
type Payload struct {
    XMLName xml.Name `xml:"payload" json:"-"`   // XML üçün kök element
    Status  string   `xml:"status" json:"status"`
}
// text/xml → <payload><status>Successful!</status></payload>
// application/json → {"status":"Successful!"}
```

**Sub-kod izahı:**
- `XMLName xml.Name` → XML kök elementin adı; json:"-" → JSON-da gizli
- unrolled/render: JSON, XML, HTML şablonları, streaming dəstəyi

### 5. Middleware (ApplyMiddleware + context)
**Nədir:** Handler-ləri saran funksiyalar — logging, ID generasiyası,
kontekst dəyərləri.

**Kitabdan kod nümunəsi:**
```go
// Middleware = HandlerFunc → HandlerFunc:
type Middleware func(http.HandlerFunc) http.HandlerFunc

func ApplyMiddleware(h http.HandlerFunc, middleware ...Middleware) http.HandlerFunc {
    applied := h
    for _, m := range middleware {
        applied = m(applied)       // ard-arda sar
    }
    return applied
}

// Logger middleware:
func Logger(l *log.Logger) Middleware {
    return func(next http.HandlerFunc) http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()
            l.Printf("started request to %s with id %s", r.URL, GetID(r.Context()))
            next(w, r)               // növbəti handler-i çağır
            l.Printf("completed request to %s with id %s in %s",
                r.URL, GetID(r.Context()), time.Since(start))
        }
    }
}

// Context-də ID saxlama:
type ContextID int
const ID ContextID = 0
func SetID(start int64) Middleware {
    return func(next http.HandlerFunc) http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {
            ctx := context.WithValue(r.Context(), ID,
                strconv.FormatInt(start, 10))
            start++
            r = r.WithContext(ctx)       // YENİ request qaytar
            next(w, r)
        }
    }
}
func GetID(ctx context.Context) string {
    if val, ok := ctx.Value(ID).(string); ok {
        return val
    }
    return ""
}

// Qoşulma (aşağıdan yuxarı tətbiq olunur):
h := middleware.ApplyMiddleware(
    middleware.Handler,
    middleware.Logger(log.New(os.Stdout, "", 0)),
    middleware.SetID(100),       // Logger-dən SONRA → ID artıq context-də
)
http.HandleFunc("/", h)
```

**Sub-kod izahı:**
- `r.WithContext(ctx)` → request dəyişməzdir; yeni nüsxə qaytarılır
- Typed context key (`ContextID`) → toqquşma qorunması
- Middleware-lərin SIRASI vacibdir — Logger ID-dən asılı olduğu üçün
  SetID əvvəl tətbiq olunmalı (ən yaxşısı asılılıqdan qaçmaq: UUID Logger-də)

### 6. Reverse proxy
**Nədir:** Gələn bütün sorğuları başqa host-a yönləndirən və cavabı geri
qaytaran handler.

**Kitabdan kod nümunəsi:**
```go
type Proxy struct {
    Client *http.Client
    BaseURL string
}
// Handler interfeysini realizə edir:
func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    if err := p.ProcessRequest(r); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return
    }
    resp, err := p.Client.Do(r)      // dəyişdirilmiş sorğunu GÖNDƏR
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()
    CopyResponse(w, resp)
}

// Sorğunu yönləndirməyə hazırla:
func (p *Proxy) ProcessRequest(r *http.Request) error {
    proxyURLRaw := p.BaseURL + r.URL.String()
    proxyURL, err := url.Parse(proxyURLRaw)
    if err != nil { return err }
    r.URL = proxyURL
    r.Host = proxyURL.Host
    r.RequestURI = ""               // client sorğusunda BOŞ olmalıdır!
    return nil
}

// Cavabı köçür (header + status + body):
func CopyResponse(w http.ResponseWriter, resp *http.Response) {
    var out bytes.Buffer
    out.ReadFrom(resp.Body)
    for key, values := range resp.Header {
        for _, value := range values {
            w.Header().Add(key, value)
        }
    }
    w.WriteHeader(resp.StatusCode)
    w.Write(out.Bytes())
}

p := &proxy.Proxy{Client: http.DefaultClient, BaseURL: "https://www.golang.org"}
http.Handle("/", p)      // Handle (Handler), HandleFunc DEYİL
```

**Sub-kod izahı:**
- `r.RequestURI = ""` → əks halda client xəta verir (yalnız server oxuyur)
- Request/Response obyektləri həm handler həm client tərəfində istifadə olunur
- Alternativ: `net/http/httputil.ReverseProxy` (Director + ModifyResponse)

### 7. gRPC-ın JSON API kimi ixracı
**Nədir:** Eyni internal RPC kodu həm gRPC, həm JSON REST server kimi serve
edilir — kod dublikasiyası yoxdur.

**Kitabdan kod nümunəsi:**
```protobuf
service KeyValue {
    rpc Set(SetKeyValueRequest) returns (KeyValueResponse) {}
    rpc Get(GetKeyValueRequest) returns (KeyValueResponse) {}
}
```
```go
// INTERNAL — hər iki serverin istifadə etdiyi ümumi impl:
type KeyValue struct {
    mutex sync.RWMutex          // paralel map qoruması
    m     map[string]string
}
func (k *KeyValue) Set(ctx context.Context, r *keyvalue.SetKeyValueRequest) (*keyvalue.KeyValueResponse, error) {
    k.mutex.Lock()
    k.m[r.GetKey()] = r.GetValue()
    k.mutex.Unlock()
    return &keyvalue.KeyValueResponse{Value: r.GetValue()}, nil
}
func (k *KeyValue) Get(ctx context.Context, r *keyvalue.GetKeyValueRequest) (*keyvalue.KeyValueResponse, error) {
    k.mutex.RLock()
    defer k.mutex.RUnlock()
    val, ok := k.m[r.GetKey()]
    if !ok {
        return nil, grpc.Errorf(codes.NotFound, "key not set")  // gRPC error kodu
    }
    return &keyvalue.KeyValueResponse{Value: val}, nil
}

// gRPC server:
grpcServer := grpc.NewServer()
keyvalue.RegisterKeyValueServer(grpcServer, internal.NewKeyValue())
grpcServer.Serve(lis)

// HTTP server — EYNİ internal obyekti wrap edir:
func (c *Controller) GetHandler(w http.ResponseWriter, r *http.Request) {
    key := r.URL.Query().Get("key")
    kv := keyvalue.GetKeyValueRequest{Key: key}
    gresp, err := c.Get(r.Context(), &kv)       // birbaşa RPC metodu!
    if err != nil {
        if grpc.Code(err) == codes.NotFound {   // gRPC kodu → HTTP status
            w.WriteHeader(http.StatusNotFound)
            return
        }
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
    resp, _ := json.Marshal(gresp)              // protobuf → JSON
    w.WriteHeader(http.StatusOK)
    w.Write(resp)
}
```

**Sub-kod izahı:**
- `internal/` paketi → hər iki binari eyni kodu istifadə edir
- gRPC error kodları → HTTP statuslarına xəritələnir (NotFound → 404)
- Protobuf struct-ları JSON tag-ləri daşıyır → birbaşa json.Marshal işləyir

## Əsas terminlər

- Handler / HandlerFunc / ServeHTTP
- ResponseWriter / Request
- Dependency Injection (asılılıq inyeksiyası)
- Closure Handler
- Validation / Typed Error (Verror)
- Content Negotiation (məzmun danışığı)
- Middleware / ApplyMiddleware
- context.WithValue / r.WithContext
- Reverse Proxy / Director
- gRPC error codes / gRPC-to-JSON wrapper

## Praktik nəticə

- Header set → WriteHeader → Write sırası; metod yoxlaması hər handler-də
- State: struct (çox handler) və ya closure (fərdi konfiq) — Storage
  interfeysi ilə test olunan DI
- Validasiyanı funksiya sahəsi kimi saxla — runtime əvəzetmə imkanı
- Middleware sırası vacibdir; context dəyərləri typed key ilə
- r.RequestURI = "" reverse proxy-də unudulmaz
- RPC kodu internal paketdə → gRPC + JSON API birgə serverlər

## Mənbə

Pages: 261-296 (Chapter 8, Go Programming Cookbook 2nd ed)
