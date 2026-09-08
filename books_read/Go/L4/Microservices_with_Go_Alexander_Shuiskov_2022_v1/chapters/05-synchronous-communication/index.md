# Chapter 5 — Synchronous Communication (səh. 93-114)

## Bu fəsil nədən bəhs edir?

Sinxron (request-response) kommunikasiya, RPC framework-ləri (Thrift, gRPC),
Protocol Buffers ilə servis API təyini və Movie aplikasiyasının HTTP+JSON-dan
gRPC-yə köçürülməsi.

## Əsas fikirlər

### 1. Sinxron kommunikasiya əsasları
- Request-response modeli; ən populyar protokol HTTP
- HTTP ötürmə kanalları: URL parametrləri, headers, body
- Xəta növləri: **Client error** (yanlış arqument, unauthorized, not found) vs
  **Server error** (bug, upstream komponent xətası)

### 2. RPC framework-lərin üstünlükləri
Manuel HTTP+JSON əvəzinə RPC framework:
- **Client/server kod generasiyası** (çoхdilli)
- Autentifikasiya (TLS, token)
- **Context propagation** (trace-lər — Ch11-də)
- Sənədləşdirmə generasiyası

| | Apache Thrift | gRPC |
|---|---|---|
| Mənşə | Facebook | Google |
| Transport | öz protokolu | **HTTP/2** |
| Serializasiya | Thrift format | **Protocol Buffers** |
| Populyarlıq | azalır | yüksək |

Kitab gRPC seçir.

### 3. Servis API təyini (movie.proto)
```proto
service MetadataService {
    rpc GetMetadata(GetMetadataRequest) returns (GetMetadataResponse);
}
message GetMetadataRequest  { string movie_id = 1; }
message GetMetadataResponse { Metadata metadata = 1; }
```
- Hər endpoint üçün AYRI Request/Response strukturləri — adlar funksiya
  adı ilə prefikslənir
- RatingService: GetAggregatedRating, PutRating; MovieService: GetMovieDetails
- Generasiya: `protoc -I=api --go_out=. --go-grpc_out=. movie.proto`
- Nəticə: `MetadataServiceClient` interfeysi (Invoke ilə) və
  `MetadataServiceServer` interfeysi + RegisterMetadataServiceServer

### 4. Internal vs Generated model
İKI model saxlanılır:
- **Internal model** (pkg/model): repository/controller/biznes məntiqində
- **Generated model** (gen): YALNIZ serializasiya üçün (şəbəkə, saxlama)

Niyə generated-i hər yerdə istifadə etmək olmaz:
1. Serializasiya formatına coupling (format dəyişsə bütün aplikasiya)
2. Generasiya alətinin versiyaları arası dəyişikliklər break edə bilər
3. Bütün sahələr optional → hər yerdə nil-check → panic riski

Həll: **mapper funksiyaları:**
```go
func MetadataToProto(m *Metadata) *gen.Metadata {...}
func MetadataFromProto(m *gen.Metadata) *Metadata {...}
```

### 5. gRPC handler implementasiyası
```go
type Handler struct {
    gen.UnimplementedMetadataServiceServer  // forward-compat üçün MÜTLƏQ
    svc *controller.MetadataService
}
func (h *Handler) GetMetadata(ctx, req) (*gen.GetMetadataResponse, error) {
    if req == nil || req.MovieId == "" {
        return nil, status.Errorf(codes.InvalidArgument, "...")
    }
    m, err := h.svc.Get(ctx, req.MovieId)
    if err != nil && errors.Is(err, controller.ErrNotFound) {
        return nil, status.Errorf(codes.NotFound, err.Error())
    }
    ...
    return &gen.GetMetadataResponse{Metadata: model.MetadataToProto(m)}, nil
}
```
- Xəta kodları: `codes.InvalidArgument`, `codes.NotFound`, `codes.Internal`
- main: `net.Listen` → `grpc.NewServer()` → `gen.Register...Server(srv, h)` → `srv.Serve(lis)`

### 6. Gateway (client tərəfi)
```go
// grpcutil: discovery + gRPC birləşməsi
func ServiceConnection(ctx, serviceName, registry) (*grpc.ClientConn, error) {
    addrs, _ := registry.ServiceAddresses(ctx, serviceName)
    return grpc.Dial(addrs[rand.Intn(len(addrs))],
        grpc.WithTransportCredentials(insecure.NewCredentials()))
}
```
- Gateway: connection → generated client → çağırış → `FromProto` mapping
- Movie main: static registry (metadata:8081, rating:8082, movie:8083)
- HTTP handler → gRPC handler dəyişdiriləndə CONTROLLER DƏYİŞMİRDİ — qatların
  ayırılması (Ch2) öz bəhrəsini verir

## Termindirmə (AZ)
- Synchronous Communication — Sinxron Rabitə (request-response)
- RPC (Remote Procedure Call) — Uzaq Prosedur Çağırışı
- gRPC — Google RPC framework-u (HTTP/2 + Protobuf)
- Gateway — Keçid (başqa servisi çağıran client kodu)
- Status Codes — Vəziyyət Kodları (gRPC xəta kodları)
- Code Generation — Kod Generasiyası

## Kviz sualları
1. Niyə internal və generated modellər ayrı saxlanılır? (Format dəyişikliyi,
   tooling versiyaları, nil-check yükü — generated yalnız serializasiyada)
2. `UnimplementedMetadataServiceServer` niyə embed olunur? (Gələcək protokol
   dəyişikliklərində kompatibilik — yeni RPC-lər break etməsin)
3. grpcutil.ServiceConnection nə edir? (Registry-dən təsadüfi instans seçib
   gRPC connection qaytarır — client-side LB)
4. gRPC Thrift-dən nə ilə fərqlənir? (HTTP/2 transport + Protobuf + geniş
   qəbul)
