# Chapter 6 — Resilient communication (Davamlı kommunikasiya)

## Bu chapter nədən bəhs edir?
Mikroservis resiliency patternlərinə — Timeout (context), Retry (interceptor), Circuit Breaker (gobreaker) — gRPC error modelinə və TLS/mTLS ilə kommunikasiyanın təhlükəsizliyinə.

## Əsas fikirlər

### 1. Niyə resiliency? (Cascading failure)
Bir servisin DB-si çöksə → Shipping ölür → Order (asılı olduğu üçün) da çökir → **cascading failure** (zəncirvari çökmə). Həll: timeout, retry, circuit breaker.

### 2. Timeout pattern — context
**Nədir:** `context.Context` — deadline, ləğv siqnalı, key-value daşıyıcısı; funksiya deadline-dan sonra avtomatik ləğv olunur.

**Client tərəf (Listing 6.1):**
```go
shippingClient := shipping.NewShippingServiceClient(conn)
ctx, _ := context.WithDeadline(context.Background(), time.Now().Add(1*time.Second))
// və ya:
ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
_, errCreate := shippingClient.Create(ctx, &shipping.CreateShippingRequest{UserId: 23})
```
**Context proqnozlaşdırılması (Listing 6.2):** Client 3s timeout → Order → Product zəncirə eyni context ötürülür; 2s+2s = 4s > 3s → DEADLINE_EXCEEDED bütün zəncirdə.
```go
func (s *server) Create(ctx context.Context, in *order.CreateOrderRequest) (*order.CreateOrderResponse, error) {
    time.Sleep(2 * time.Second)
    productInfo, err := s.productClient.Get(ctx, ...)   // ctx proqnozlaşır
    ...
}
```
**Qayda:** WithTimeout/WithDeadline həm context, həm cancel funksiyası qaytarır — 3 rəqabətli çağırışdan ilk cavab gələndə qalanlarını ləğv et.

### 3. Retry pattern — go-grpc-middleware
**Transient faults (keçici xətalar):** ani network xətası, müvəqqəti əlçatmazlıq, load-dan resurs tükənməsi.

**Listing 6.3:**
```go
import grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"

var opts []grpc.DialOption
opts = append(opts,
    grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(
        grpc_retry.WithCodes(codes.Unavailable, codes.ResourceExhausted),  // HANSI kodlarda
        grpc_retry.WithMax(5),                                            // max say
        grpc_retry.WithBackoff(grpc_retry.BackoffLinear(time.Second)),    // interval
    )))
opts = append(opts, grpc.WithInsecure())
conn, err := grpc.Dial("localhost:8080", opts...)
```
**Parametrlər:**
- **WithCodes:** default Unavailable + ResourceExhausted — validation xətasında (InvalidArgument) retry MAĞAMSIZDIR
- **WithMax:** limit; uğur gələndə dayanır
- **Backoff:** Linear (sabit) / Exponential (2× artan) / +Jitter (random — toqquşma azaldır)

**Interceptor növləri:** WithUnaryInterceptor / WithStreamingInterceptor; client/server tərəflərdə.

### 4. Circuit Breaker — sony/gobreaker
**State machine:** Closed → (failure rate > threshold) → Open → (timeout sonra) → Half-Open → (uğursuzsa Open / uğurlusa Closed).

**Settings:**
| Parametr | Vəzifə |
|---|---|
| **MaxRequests** | half-open-dan keçə bilən sorğu limiti (open = 0) |
| **Interval** | closed halında sayğacın sıfırlanması |
| **Timeout** | open → half-open keçid vaxtı |
| **ReadyToTrip** | dövrü AÇMAQ şərti (threshold yoxlaması) |
| **OnStateChange** | state dəyişikliyi callback (loq üçün) |

**Listing 6.4:**
```go
cb = gobreaker.NewCircuitBreaker(gobreaker.Settings{
    Name:        "demo",
    MaxRequests: 3,
    Timeout:     4,
    ReadyToTrip: func(counts gobreaker.Counts) bool {
        failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
        return failureRatio >= 0.6          // 60% xəta → aç
    },
    OnStateChange: func(name string, from, to gobreaker.State) {
        log.Printf("Circuit Breaker: %s, changed from %v, to %v", name, from, to)
    },
})
cbRes, cbErr := cb.Execute(func() (interface{}, error) {
    // iş məntiqi burada
    return res, nil
})
```
**Listing 6.5 (client dövrü):** `cb.Execute` içində RPC çağırışı + 1s sleep ilə periodik sorğu.

**Listing 6.6 (mərkəzləşdirilmiş interceptor):**
```go
func CircuitBreakerClientInterceptor(cb *gobreaker.CircuitBreaker) grpc.UnaryClientInterceptor {
    return func(ctx context.Context, method string, req, reply interface{},
        cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
        _, cbErr := cb.Execute(func() (interface{}, error) {
            err := invoker(ctx, method, req, reply, cc, opts...)  // əsl RPC
            if err != nil { return nil, err }
            return nil, nil
        })
        return cbErr
    }
}
// İstifadə:
opts = append(opts, grpc.WithUnaryInterceptor(middleware.CircuitBreakerClientInterceptor(cb)))
```
**Prinsip:** CB məntiqini BÜTÜN çağırışlara bir mərkəzdən tətbiq et — hər yerdə əl ilə sarma YOX.

### 5. gRPC error modeli
**Struktur:** Code + Status (insan-dostlu) + Message + ErrorDetails (violation siyahısı).
**Fərq:** strukturlaşdırılmış xəta (kodlu) → qərar məntiqi mümkün: retry? circuit aç? {code: ResourceExhausted, message: Order Limit Exceeded} → sonra retry məntiqlidir.

**Listing 6.7 (server validasiyası):**
```go
var validationErrors []*errdetails.BadRequest_FieldViolation
if in.UserId < 1 {
    validationErrors = append(validationErrors, &errdetails.BadRequest_FieldViolation{
        Field: "user_id", Description: "user id cannot be less than 1",
    })
}
if in.ProductId < 0 { /* eyni pattern */ }
if len(validationErrors) > 0 {
    stat := status.New(400, "invalid order request")
    badRequest := &errdetails.BadRequest{}
    badRequest.FieldViolations = validationErrors
    s, _ := stat.WithDetails(badRequest)
    return nil, s.Err()
}
```
**Client emalı:**
```go
stat := status.Convert(errCreate)
for _, detail := range stat.Details() {
    switch errType := detail.(type) {
    case *errdetails.BadRequest:
        for _, violation := range errType.GetFieldViolations() {
            log.Printf("The field %s has invalid value. desc: %v",
                violation.GetField(), violation.GetDescription())
        }
    }
}
```

### 6. TLS / mTLS
**TLS handshake (server-only):** client qoşulur → server certificate göstərir → client təsdiqləyir → şifrələnmiş kanal.
**mTLS (mutual):** client DƏ certificate göstərir → server də təsdiqləyir → **zero-trust** mühit üçün.

**Sertifikat generasiyası (OpenSSL):**
```bash
# 1. CA (root):
openssl req -x509 -sha256 -newkey rsa:4096 -days 365 \
    -keyout ca-key.pem -out ca-cert.pem \
    -subj "/C=TR/ST=EURASIA/L=ISTANBUL/O=Software/OU=Microservices/CN=*.microservices.dev/emailAddress=..." \
    -nodes

# 2. Server CSR:
openssl req -newkey rsa:4096 -keyout server-key.pem -out server-req.pem \
    -subj ".../OU=PaymentService/CN=*.microservices.dev" -nodes -sha256

# 3. CA ilə imzala:
openssl x509 -req -in server-req.pem -days 60 \
    -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial \
    -out server-cert.pem -extfile server-ext.cnf -sha256
# server-ext.cnf: subjectAltName=DNS:*.microservices.dev,IP:0.0.0.0

# 4. Eyni addımlarla CLIENT sertifikatı (OU=OrderService)
```
**subj parametrləri:** /C ölkə, /ST ştat, /L şəhər, /O təşkilat, /OU departament, /CN domain, /emailAddress.

**gRPC server TLS (Listing 6.8):**
```go
func getTlsCredentials() (credentials.TransportCredentials, error) {
    serverCert, _ := tls.LoadX509KeyPair("cert/server.crt", "cert/server.key")
    certPool := x509.NewCertPool()
    caCert, _ := ioutil.ReadFile("cert/ca.crt")
    certPool.AppendCertsFromPEM(caCert)
    return credentials.NewTLS(&tls.Config{
        ClientAuth:   tls.RequireAnyClientCert,   // mTLS: client sert tələb et
        Certificates: []tls.Certificate{serverCert},
        ClientCAs:   certPool,                     // client sert yoxlama kökü
    }), nil
}

opts = append(opts, grpc.Creds(tlsCredentials))
grpcServer := grpc.NewServer(opts...)
```
**gRPC client TLS:**
```go
func getTlsCredentials() (credentials.TransportCredentials, error) {
    clientCert, _ := tls.LoadX509KeyPair("cert/client.crt", "cert/client.key")
    certPool := x509.NewCertPool()
    certPool.AppendCertsFromPEM(caCert)
    return credentials.NewTLS(&tls.Config{
        ServerName:   "*.microservices.dev",    // sertifikatdakı CN ilə uyğun
        Certificates: []tls.Certificate{clientCert},
        RootCAs:      certPool,                  // server sert yoxlama kökü
    }), nil
}

opts = append(opts, grpc.WithTransportCredentials(tlsCredentials))
conn, err := grpc.Dial("localhost:8080", opts...)
```

## Əsas terminlər
- Cascading failure — zəncirvari çökmə
- context.WithTimeout / WithDeadline / cancel funksiyası
- Deadline propagation — context-in zəncir boyu ötürülməsi
- Transient fault — keçici xəta
- grpc_retry.UnaryClientInterceptor / WithCodes / WithMax / WithBackoff
- Backoff: Linear / Exponential / Jitter
- Circuit breaker: Closed / Open / Half-Open
- gobreaker.Settings — MaxRequests/Interval/Timeout/ReadyToTrip/OnStateChange
- UnaryClientInterceptor / UnaryServerInterceptor
- gRPC error model — Code/Status/Message/ErrorDetails
- TLS handshake / mTLS (mutual TLS)
- CA / CSR / self-signed / subjectAltName
- credentials.NewTLS / ClientCAs / RootCAs / ServerName

## Praktik nəticə
1. Hər interservice çağırışda context timeout qoy — sonsuz gözləmə cascading failure deməkdir.
2. Retry yalnız Unavailable/ResourceExhausted üçün; InvalidArgument retry-si faydasız yükdür.
3. Circuit breaker-ı interceptor kimi mərkəzləşdir — hər çağırışda əl ilə YOX; ReadyToTrip threshold-u fail ratio.
4. Xətalar strukturlaşdırılmış statusla qaytarılsın — kod = qərar (retry/circuit/auth) açarı.
5. Zero-trust mühitdə mTLS: hər servis öz sertifikatını təqdim edir; dev-də self-signed, production-da K8s cert-manager (Ch 8).
6. TLS konfiqurasiyasında ServerName sertifikat CN-ı ilə uyğun olmalıdır — yoxsa hostname verification xətası.

## Mənbə
Pages: 81-98 (PDF səh. 105-122)
