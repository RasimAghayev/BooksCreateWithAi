# Chapter 7 — Веб-клиенты и API (Veb klientləri və API)

## Bu chapter nədən bəhs edir?

http.Client strukturu və ötürülməsi, REST API client, parallel/async
sorğular, OAuth2 client + token saxlancı, RoundTripper decorator middleware,
gRPC client/server, twirp RPC.

## Əsas fikirlər

### 1. http.Client: init, saxlama, ötürmə
**Nədir:** net/http-in çevik client strukturu — Transport dəyişdirilərək
istənilən səviyyədə fərdiləşdirilir.

**Kitabdan kod nümunəsi:**
```go
func Setup(isSecure, nop bool) *http.Client {
    c := http.DefaultClient
    if !isSecure {
        c.Transport = &http.Transport{
            TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
        }
    }
    if nop {
        c.Transport = &NopTransport{}     // test üçün saxta transport
    }
    http.DefaultClient = c                // qlobal default-u da dəyişir
    return c
}

// RoundTripper interfeysi — test double üçün:
type NopTransport struct{}
func (n *NopTransport) RoundTrip(*http.Request) (*http.Response, error) {
    return &http.Response{StatusCode: http.StatusTeapot}, nil  // 418!
}

// Client-i funksiyaya ötür:
func DoOps(c *http.Client) error {
    resp, err := c.Get("http://www.google.com")
    fmt.Println(resp.StatusCode)
    return nil
}
// Struct-a EMBED et:
type Controller struct{ *http.Client }
func (c *Controller) DoOps() error { c.Client.Get(...) }
```

**Sub-kod izahı:**
- `http.Get()` → `http.DefaultClient` istifadə edir; fərdi client ötürmək
  daha yaxşı test imkanı verir
- Transport dəyişəndə bütün sorğular ondan keçir → mock/test üçün ideal
- 418 (Teapot) — testdə tanınan status kodu kimi istifadə olunur

### 2. REST API client (Transport ilə auth)
**Nədir:** API çağırışlarını metodlara gizlədən client; auth Transport
qatında həll olunur.

**Kitabdan kod nümunəsi:**
```go
type APIClient struct{ *http.Client }

func NewAPIClient(username, password string) *APIClient {
    t := http.Transport{}
    return &APIClient{
        Client: &http.Client{
            Transport: &APITransport{
                Transport: &t,
                username: username,
                password: password,
            },
        },
    }
}

// Transport hər sorğuya basic auth əlavə edir:
type APITransport struct {
    *http.Transport
    username, password string
}
func (t *APITransport) RoundTrip(req *http.Request) (*http.Response, error) {
    req.SetBasicAuth(t.username, t.password)   // HƏR sorğu üçün avtomatik
    return t.Transport.RoundTrip(req)
}

// REST detalları metodlarda gizlidir:
func (c *APIClient) GetGoogle() (int, error) {
    resp, err := c.Get("http://www.google.com")
    return resp.StatusCode, err
}
```

**Sub-kod izahı:**
- `RoundTrip` → hər sorğunun keçdiyi vahid nöqtə; auth/token refresh bura
- API səthi metodlardan ibarətdir: `GetUsers()`, `CreateUser(u)` və s.
- İstifadəçi yalnız metod çağırır — HTTP detalı görünmür

### 3. Parallel və async client sorğuları
**Nədir:** Buffered kanallarla paralel GET; cavablar və xətalar ayrı
kanallara toplanır.

**Kitabdan kod nümunəsi:**
```go
type Client struct {
    *http.Client
    Resp chan *http.Response    // buffered cavab kanalı
    Err  chan error              // buffered xəta kanalı
}
func NewClient(client *http.Client, bufferSize int) *Client {
    return &Client{
        Client: client,
        Resp: make(chan *http.Response, bufferSize),
        Err:  make(chan error, bufferSize),
    }
}
func (c *Client) AsyncGet(url string) {
    resp, err := c.Get(url)
    if err != nil { c.Err <- err; return }   // xəta kanala
    c.Resp <- resp                            // cavab kanala
}
func FetchAll(urls []string, c *Client) {
    for _, url := range urls {
        go c.AsyncGet(url)        // hər URL ayrı goroutine
    }
}

// İstifadə — select ilə hər iki kanalı dinlə:
c := async.NewClient(http.DefaultClient, len(urls))
async.FetchAll(urls, c)
for i := 0; i < len(urls); i++ {       // sayı bəlli — o qədər oxu
    select {
    case resp := <-c.Resp:
        fmt.Printf("Status received for %s: %d\n", resp.Request.URL, resp.StatusCode)
    case err := <-c.Err:
        fmt.Printf("Error received: %s\n", err)
    }
}
```

**Sub-kod izahı:**
- Buffered kanal → goroutine-lər bloklanmır (bufferSize = URL sayı)
- `select` → iki kanaldan hansı gəlibsə emal; main sayğac ilə bitməsini bilir
- Alternativ: worker pool (limit üçün), context timeout (ləğv üçün)

### 4. OAuth2 client (golang.org/x/oauth2)
**Nədir:** GitHub, Google, Facebook üçün hazır endpoint-ləri olan OAuth2
kitabxanası; kod↔token mübadiləsi.

**Kitabdan kod nümunəsi:**
```go
func Setup() *oauth2.Config {
    return &oauth2.Config{
        ClientID:     os.Getenv("GITHUB_CLIENT"),
        ClientSecret: os.Getenv("GITHUB_SECRET"),
        Scopes:       []string{"repo", "user"},
        Endpoint:     github.Endpoint,     // hazır provider endpoint-i
    }
}

// Token əldə etmə axını:
func GetToken(ctx context.Context, conf *oauth2.Config) (*oauth2.Token, error) {
    url := conf.AuthCodeURL("state")     // istifadəçini yönləndirəcək URL
    fmt.Printf("Type the following url into your browser: %v\n", url)
    var code string
    fmt.Scan(&code)                       // code-u əl ilə daxil et
    return conf.Exchange(ctx, code)       // code → token mübadiləsi
}

// Token-li client avtomatik refresh edir:
client := conf.Client(ctx, tok)
resp, err := client.Get("https://api.github.com/user")
defer resp.Body.Close()
io.Copy(os.Stdout, resp.Body)
```

**Sub-kod izahı:**
- `AuthCodeURL("state")` → istifadəçinin brauzerdə təsdiqləyəcəyi URL
- `Exchange(ctx, code)` → authorization code-u tokenə çevirir
- `conf.Client(ctx, tok)` → token-i avtomatik yeniləyən http.Client
- Token saxlanmır → crash-də yenidən exchange lazımdır (növbəti resept həlli)

### 5. OAuth2 token saxlancı (Storage interfeysi)
**Nədir:** Token-i faylda/DB-də saxlayan TokenSource — hər run-da yenidən
auth-a ehtiyac yoxdur.

**Kitabdan kod nümunəsi:**
```go
// Storage interfeysi:
type Storage interface {
    GetToken() (*oauth2.Token, error)
    SetToken(*oauth2.Token) error
}

// Config — Exchange-də saxlayır:
type Config struct {
    *oauth2.Config
    Storage
}
func (c *Config) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
    token, err := c.Config.Exchange(ctx, code)
    if err != nil { return nil, err }
    c.Storage.SetToken(token)              // hər alınanda SAXLA
    return token, nil
}

// TokenSource — storage-dan oxu, yoxdursa al:
func (s *storageTokenSource) Token() (*oauth2.Token, error) {
    if token, err := s.Config.Storage.GetToken(); err == nil && token.Valid() {
        return token, err                    // saxlanmış etibarlı token
    }
    token, err := s.TokenSource.Token()     // yenisi al / refresh
    if err != nil { return token, err }
    s.Config.Storage.SetToken(token)        // yenini saxla
    return token, nil
}

// Fayl implementasiyası (RWMutex ilə):
type FileStorage struct {
    Path string
    mu   sync.RWMutex      // paralel oxu/yazma qoruması
}
func (f *FileStorage) SetToken(t *oauth2.Token) error {
    if t == nil || !t.Valid() { return errors.New("bad token") }
    f.mu.Lock()
    defer f.mu.Unlock()
    out, err := os.OpenFile(f.Path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
    defer out.Close()
    data, _ := json.Marshal(&t)
    _, err = out.Write(data)
    return err
}
```

**Sub-kod izahı:**
- İlk run: tam exchange; sonrakı run-lar: fayldan oxu + avtomatik refresh
- Storage interfeysi → fayl əvəzinə DB/Redis qoşula bilər
- `sync.RWMutex` → çox istifadəçili serverdə eyni anda oxuma/qeyd
- İstifadəçi fərqi → cookie/DB açarı ilə fayl adı fərqləndirilə bilər

### 6. Client decorator middleware (funksional kompozisiya)
**Nədir:** RoundTripper funsiya-tiplə birləşdirərək middleware zənciri
(logging, auth) — Thomas Senarin ideyası.

**Kitabdan kod nümunəsi:**
```go
// Funksiya interfeysi realizə edir:
type TransportFunc func(*http.Request) (*http.Response, error)
func (tf TransportFunc) RoundTrip(r *http.Request) (*http.Response, error) {
    return tf(r)
}

// Decorator = RoundTripper → RoundTripper:
type Decorator func(http.RoundTripper) http.RoundTripper

// Zəncir hər bir middleware-i tətbiq edir:
func Decorate(t http.RoundTripper, rts ...Decorator) http.RoundTripper {
    decorated := t
    for _, rt := range rts {
        decorated = rt(decorated)      // iç-içə sarır
    }
    return decorated
}

// Logging middleware:
func Logger(l *log.Logger) Decorator {
    return func(c http.RoundTripper) http.RoundTripper {
        return TransportFunc(func(r *http.Request) (*http.Response, error) {
            start := time.Now()
            l.Printf("started request to %s at %s", r.URL, start.Format("2006-01-02 15:04:05"))
            resp, err := c.RoundTrip(r)          // növbəti qata ötür
            l.Printf("completed request to %s in %s", r.URL, time.Since(start))
            return resp, err
        })
    }
}

// BasicAuth middleware:
func BasicAuth(username, password string) Decorator {
    return func(c http.RoundTripper) http.RoundTripper {
        return TransportFunc(func(r *http.Request) (*http.Response, error) {
            r.SetBasicAuth(username, password)
            return c.RoundTrip(r)
        })
    }
}

// İstifadə — bir sətirlə zəncir:
c.Transport = Decorate(&http.Transport{},
    Logger(log.New(os.Stdout, "", 0)),
    BasicAuth("username", "password"),
)
```

**Sub-kod izahı:**
- `Decorator` tipi → middleware funksiyasının imzası standartlaşdırılır
- `Decorate(t, mw1, mw2)` → sırayla sarır; massiv variadic
- Funksiyanın interfeysi realizə etməsi → closure gücü + interface uyğunluğu

### 7. gRPC client
**Nədir:** Protobuf + HTTP/2 üzərində yüksək performanslı RPC.

**Kitabdan kod nümunəsi:**
```protobuf
syntax = "proto3";
package greeter;
service GreeterService {
    rpc Greet(GreetRequest) returns (GreetResponse) {}
}
message GreetRequest { string greeting = 1; string name = 2; }
message GreetResponse { string response = 1; }
```
```go
// SERVER — protoc-un yaratdığı interfeysi doldur:
type Greeter struct{ Exclaim bool }
func (g *Greeter) Greet(ctx context.Context, r *greeter.GreetRequest) (*greeter.GreetResponse, error) {
    msg := fmt.Sprintf("%s %s", r.GetGreeting(), r.GetName())
    if g.Exclaim { msg += "!" } else { msg += "." }
    return &greeter.GreetResponse{Response: msg}, nil
}
grpcServer := grpc.NewServer()
greeter.RegisterGreeterServiceServer(grpcServer, &Greeter{Exclaim: true})
lis, _ := net.Listen("tcp", ":4444")
grpcServer.Serve(lis)

// CLIENT:
conn, err := grpc.Dial(":4444", grpc.WithInsecure())   // test üçün insecure
defer conn.Close()
client := greeter.NewGreeterServiceClient(conn)
resp, err := client.Greet(ctx, &greeter.GreetRequest{Greeting: "Hello", Name: "Reader"})
// protoc --go_out=plugins=grpc:. greeter/greeter.proto  → kod generasiyası
```

**Sub-kod izahı:**
- `.proto` → `protoc` server stub + client doldurur; JSON tag-ləri də var
- `grpc.WithInsecure()` → yalnız test; produksiyada TLS şərtdir
- Streaming bu reseptdə yoxdur (sadə unary çağırış)

### 8. twirp RPC (twitchtv/twirp)
**Nədir:** gRPC-nin faydaları (protobuf modellər) + HTTP/1.1 + JSON dəstəyi —
curl ilə test oluna bilən RPC.

**Kitabdan kod nümunəsi:**
```go
// SERVER — sadə HTTP handler kimi:
server := &Greeter{}
twirpHandler := greeter.NewGreeterServiceServer(server, nil)
http.ListenAndServe(":4444", twirpHandler)

// CLIENT — custom http.Client qəbul edir:
client := greeter.NewGreeterServiceProtobufClient("http://localhost:4444", &http.Client{})
resp, err := client.Greet(ctx, &req)

// CURL ilə sına:
// curl --request POST --location \
//   "http://localhost:4444/twirp/greeter.GreeterService/Greet" \
//   --header "Content-Type:application/json" \
//   --data '{"greeting": "Greetings to", "name":"you"}'
```

**Sub-kod izahı:**
- HTTP/1.1 üzərində → gRPC-yə infrastruktur dəstəyi olmayan muhitlərdə işləyir
- JSON content-type → protobuf xaricində curl/browser də qoşula bilir
- Eyni `.proto` → gRPC və twirp arasında keçid asandır

## Əsas terminlər

- http.Client / http.DefaultClient
- Transport / RoundTripper interfeysi
- REST Client pattern
- Async / Parallel sorğular
- Buffered Channel / select
- OAuth2 / Authorization Code / Token Exchange / Refresh Token
- TokenSource / Token Storage
- Decorator / Middleware (funksional kompozisiya)
- gRPC / Protobuf / protoc / unary RPC
- twirp RPC / HTTP/1.1 RPC

## Praktik nəticə

- Client-i həmişə parametr/struct kimi ötür — DefaultClient-dan asılılıq testi
  çətinləşdirir
- Auth/token refresh Transport-da (RoundTrip) həll olunmalı — API metodu təmiz qalır
- Async üçün: buffered kanallar + select + sayğac (və ya worker pool + context)
- OAuth2 token-i mütləq saxla (Storage interfeysi) — hər run-da auth olma
- Decorate pattern: bir sətirlə middleware zənciri — logging/auth/metrics
- twirp: gRPC modelləri + HTTP/1.1 sadəliyi (curl debug mümkün)

## Mənbə

Pages: 218-260 (Chapter 7, Go Programming Cookbook 2nd ed)
