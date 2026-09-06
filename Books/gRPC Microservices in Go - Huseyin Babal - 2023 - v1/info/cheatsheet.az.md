# gRPC Microservices in Go — Cheat Sheet (Azərbaycanca)

> **Kitab:** gRPC Microservices in Go — Hüseyin Babal, Manning, 2023 (ISBN 978-1-63343-920-7) · 🚀 Advanced (4/5)
> E-commerce layihəsi (Order/Payment/Shipping) boyunca bütün əsas kod/konfiqurasiya — bir baxışda.

---

## 1. Protobuf (Ch 2-3)
```protobuf
syntax = "proto3";
option go_package="github.com/user/microservices-proto/golang/order";

message CreateOrderRequest {
    int64 user_id = 1;              // 1-15 = 1 bayt metadata (tez istifadə!)
    repeated Item items = 2;
    float amount = 3;
}
service Order {
    rpc Create(CreateOrderRequest) returns (CreateOrderResponse) {}
    // streaming: req → (stream Resp) | (stream Req) → resp | (stream) → (stream)
}
```
**Qaydalar:** silinən sahələr `reserved 1, 3 to 7; reserved "ad";`; oneof dəyişikliyi = BREAKING (semver major).

## 2. protoc + stub generasiyası (Ch 3)
```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

protoc -I ./proto \
   --go_out ./golang --go_opt paths=source_relative \
   --go-grpc_out ./golang --go-grpc_opt paths=source_relative \
   ./proto/order.proto
# → order.pb.go (mesajlar) + order_grpc.pb.go (servis)
```
**Ayrı proto repo:** `golang/order/` modulları + `git tag golang/order/v1.2.3` → `go get .../golang/order@v1.2.3`

## 3. Hexagonal layihə strukturu (Ch 4)
```
internal/
  ports/            # APIPort, DBPort, PaymentPort (interfeyslər)
  application/core/ # api (biznes) + domain (model)
  adapters/         # db (GORM), grpc, payment
cmd/main.go         # DI: db→app→grpc→Run
config/            # env (fail-fast)
```

## 4. gRPC server adapter (Ch 4)
```go
type Adapter struct {
    api  ports.APIPort
    order.UnimplementedOrderServer      // forward compat şərti
}
func (a Adapter) Run() {
    listen, _ := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
    grpcServer := grpc.NewServer()
    order.RegisterOrderServer(grpcServer, a)
    if config.GetEnv() == "development" { reflection.Register(grpcServer) }
    grpcServer.Serve(listen)
}
```

## 5. Client dial + status kodları (Ch 5)
```go
conn, err := grpc.Dial(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
paymentClient := payment.NewPaymentClient(conn)
paymentClient.Create(ctx, &payment.CreatePaymentRequest{...})
```
| Kod | Nə vaxt |
|---|---|
| INVALID_ARGUMENT | client input xətası |
| DEADLINE_EXCEEDED | timeout bitdi |
| NOT_FOUND / ALREADY_EXISTS | resurs yox/duplikat |
| RESOURCE_EXHAUSTED | quota/limit |
| UNAVAILABLE | servis yox (retry!) |

## 6. Error details (Ch 5-6)
```go
// Server: kod + mesaj + detal:
err = status.Errorf(codes.InvalidArgument, "failed to charge user: %d", id)
st := status.New(codes.InvalidArgument, "order creation failed")
badReq := &errdetails.BadRequest{}
badReq.FieldViolations = append(..., &errdetails.BadRequest_FieldViolation{
    Field: "payment", Description: msg })
statusWithDetails, _ := st.WithDetails(badReq)
return nil, statusWithDetails.Err()

// Client: detalları aç:
st := status.Convert(err)
for _, d := range st.Details() {
    if br, ok := d.(*errdetails.BadRequest); ok {
        for _, v := range br.GetFieldViolations() { ... }
    }
}
```

## 7. Resiliency üçlüyü (Ch 6)
```go
// 1) TIMEOUT — context:
ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
client.Create(ctx, req)                       // ctx zəncir boyu yayılır!

// 2) RETRY — interceptor:
grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(
    grpc_retry.WithCodes(codes.Unavailable, codes.ResourceExhausted), // SADECE bunlar!
    grpc_retry.WithMax(5),
    grpc_retry.WithBackoff(grpc_retry.BackoffLinear(time.Second))))

// 3) CIRCUIT BREAKER — gobreaker:
cb = gobreaker.NewCircuitBreaker(gobreaker.Settings{
    MaxRequests: 3, Timeout: 4,
    ReadyToTrip: func(c gobreaker.Counts) bool {
        return float64(c.TotalFailures)/float64(c.Requests) >= 0.6 },
    OnStateChange: func(n string, from, to gobreaker.State) { log... },
})
// Mərkəzləşdirilmiş:
grpc.WithUnaryInterceptor(middleware.CircuitBreakerClientInterceptor(cb))
```

## 8. mTLS (Ch 6)
```bash
# CA → server CSR → CA imza → client CSR → CA imza (openssl req -x509 / req / x509 -req)
```
```go
// Server: ClientAuth: tls.RequireAnyClientCert, ClientCAs: certPool
grpc.Creds(credentials.NewTLS(&tls.Config{...}))
// Client: ServerName: "*.microservices.dev", RootCAs: certPool
grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{...}))
```

## 9. Test piramidi (Ch 7)
```go
// Unit: mock (testify) — On/Return/Called
payment.On("Charge", mock.Anything).Return(errors.New("insufficient balance"))
// mockery --all --keeptree → avtomatik mocklar

// Integration: Testcontainers MySQL + suite.SetupSuite/TearDownSuite
req := testcontainers.ContainerRequest{Image: "mysql:8.0.30",
    Env: ..., WaitingFor: wait.ForSQL(...)}

// E2E: docker-compose stack (mysql+payment+order) + grpc client assert
compose.WithCommand([]string{"up", "-d"}).Invoke()

// Coverage:
go test -cover ./... ; go tool cover -html=coverage.out
```

## 10. K8s deploy (Ch 8)
```yaml
# Ingress — 1 LB, çox servis:
metadata:
  annotations:
    nginx.ingress.kubernetes.io/backend-protocol: GRPC
    cert-manager.io/cluster-issuer: selfsigned-issuer
spec:
  rules: [{ http: { paths: [{ path: /Order, backend: { service: { name: order } } }] } }]
  tls: [{ hosts: [ ingress.local ] }]
```
```bash
helm install nginx-ingress ingress-nginx/ingress-nginx
helm install cert-manager jetstack/cert-manager --set installCRDs=true
kubectl set image deployment/order order=order:1.1.0   # RollingUpdate
```
**Strategiyalar:** RollingUpdate (default, backward-compat tələb!) | Blue-Green (switch + dublikat xərc) | Canary (eyni selector, tədricən keçid).

## 11. Observability (Ch 9)
```go
// OTel interceptor:
grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor())
// Trace-inyeksiyalı log (logrus):
span := trace.SpanFromContext(entry.Context)
entry.Data["trace_id"] = span.SpanContext().TraceID().String()
entry.Data["span_id"] = span.SpanContext().SpanID().String()
```
**Stack:** OTel Collector → Prometheus (time-series) + Jaeger (trace/SPM p95); Fluent Bit DaemonSet → Elasticsearch → Kibana (TraceID filtri).

---

## Qızıl qaydalar
1. Field №1-15 tez sahələrə; silinənlər reserved.
2. Xətalar HƏMİŞƏ status.Errorf ilə — xam err "Unknown" qaytarır.
3. Retry yalnız Unavailable/ResourceExhausted; InvalidArgument retry-si zərərdir.
4. CB-i interceptor kimi mərkəzləşdir (ReadyToTrip = fail ratio).
5. Interfeysə (port) asılı ol — konkret adapter-ə YOX (mock/test mümkünü).
6. p95 > average — SLA vədlərində percentile işlət.
7. gRPC Ingress-də `backend-protocol: GRPC` şərtdir.
8. Hər loga trace_id — Kibana-da tam axın axtarışı.

## grpcurl tez istinadlar
```bash
grpcurl -d '{"user_id": 123}' -plaintext localhost:3000 Order/Create
grpcurl -import-path ./proto -proto order.proto ingress.local:443 Order.Create  # TLS ilə
```
