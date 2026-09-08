# Microservices with Go — Cheatsheet (AZ)

## Servis strukturu
```
service/
├── cmd/main.go              # main: komponentləri yığ, server başlat
├── internal/
│   ├── controller/          # biznes məntiqi
│   ├── handler/grpc|http/   # API qatı (status kodlar/mapping)
│   ├── repository/memory|mysql/  # DB qatı
│   └── gateway/metadata|rating/  # başqa servislərə çağırış
└── pkg/model/               # export tiplər
```

## gRPC sxemi (movie.proto)
```proto
service MetadataService {
    rpc GetMetadata(GetMetadataRequest) returns (GetMetadataResponse);
}
// protoc -I=api --go_out=. --go-grpc_out=. movie.proto
```
- Handler: `gen.UnimplementedMetadataServiceServer` embed + codes.InvalidArgument/NotFound/Internal
- Gateway: registry-dən addr → `grpc.Dial` → generated client → `FromProto` mapping

## Retry + backoff
```go
for i := 0; i < maxRetries; i++ {
    resp, err = client.GetMetadata(ctx, req)
    if err != nil {
        if shouldRetry(err) { continue }  // DeadlineExceeded|ResourceExhausted|Unavailable
        return nil, err
    }
    return resp, nil
}
```

## Rate limiting (token bucket, gRPC interceptor)
```go
l := rate.NewLimiter(rate.Limit(100), 100)   // golang.org/x/time/rate
srv := grpc.NewServer(grpc.UnaryInterceptor(
    ratelimit.UnaryServerInterceptor(&limiter{l})))
// aşım → codes.ResourceExhausted; HTTP → 429 Too Many Requests
```

## Graceful shutdown
```go
ctx, cancel := context.WithCancel(context.Background())
sigChan := make(chan os.Signal, 1)
signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
go func() { <-sigChan; cancel(); srv.GracefulStop() }()
wg.Wait()
```

## Kafka ingestion
```go
consumer := kafka.NewConsumer(&kafka.ConfigMap{
    "bootstrap.servers": addr, "group.id": gid,
    "auto.offset.reset": "earliest"})
// goroutine: ReadMessage(-1) → JSON unmarshal → chan RatingEvent
// controller: for e := range ch { PutRating(...) }
```
- Producer: Produce + Flush(timeout)

## Observability
```go
// zap
logger, _ := zap.NewProduction()
logger.Info("msg", zap.String("k", "v"))
logger = logger.With(zap.String("component", "handler"))
// tally → Prometheus /metrics
reporter := prometheus.NewReporter(prometheus.Options{})
scope, _ := tally.NewRootScope(tally.ScopeOptions{
    Tags: map[string]string{"service": "rating"},
    CachedReporter: reporter}, 10*time.Second)
scope.Counter("request_count").Inc(1)
scope.Gauge("active_users").Update(n)
sw := scope.Timer("op_latency").Start(); defer sw.Stop()
// OpenTelemetry
otel.SetTracerProvider(tracing.NewJaegerProvider(url, name))
otel.SetTextMapPropagator(propagation.TraceContext{})
grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor())   // client
grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor())      // server
```

## Prometheus alerting
```yaml
- alert: Rating service down
  expr: up{service="rating"} == 0
  for: 3m
  labels: {severity: page}
```

## JWT
```go
token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
    "username": u, "iat": time.Now().Unix()})
tokenString, _ := token.SignedString(secret)
// client: req.Header.Set("Authorization", "Bearer "+token)
// verify: jwt.Parse + alqoritm yoxlaması + token.Valid
```

## Testlər
- Table-driven + subtests (`t.Run`, `t.Parallel`, `testing.Short()`)
- gomock: `m.EXPECT().Get(ctx, id).Return(...).Times(1)`
- cmp: `cmp.Diff(want, got, cmpopts.IgnoreUnexported(gen.Metadata{}))`
- İnteqrasiya: setup → əməliyyatlar → teardown (GracefulStop); got/want sırası!

## MySQL repo
```go
import _ "github.com/go-sql-driver/mysql"
db, _ := sql.Open("mysql", "root:pass@/movieexample")
row := db.QueryRowContext(ctx, "SELECT ... WHERE id = ?", id)
// sql.ErrNoRows → ErrNotFound
```
