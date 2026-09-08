# Chapter 14 — gRPC (səh. 282-325)

## Bu fəsil nədən bəhs edir?

gRPC — dil/platf-müneyyən RPC həlli: .proto-da service tərifi, protoc +
gRPC plugin ilə stub generasiyası, server və client qurulması, streaming
(server/client/bidirectional), HTTP/JSON transcoding (grpc-gateway) və
interceptor-lar (gRPC middleware-i).

## Əsas fikirlər

### 1. Service tərifi (.proto)
**Nədir:** RPC = uzaq funksiya çağırışı. gRPC-də service = RPC-lər qrupu;
IDL rolunu .proto oynayır (PB mesajları + service bloku).

**Kitabdan kod nümunəsi:**
```protobuf
syntax = "proto3";
package user;

option go_package="github.com/juanmanuel-tirado/savetheworldwithgo/13_grpc/example_01/user";

message User {
    string user_id = 1;
    string email = 2;
}

message UserRequest {
    string user_id = 1;
}

service UserService {
    rpc GetUser (UserRequest) returns (User);
}
```

**Stub generasiyası:**
```bash
>>> go get google.golang.org/protobuf/cmd/protoc-gen-go \
         google.golang.org/grpc/cmd/protoc-gen-go-grpc
>>> export PATH="$PATH:$(go env GOPATH)/bin"
>>> protoc -I=. --go_out=$GOPATH/src --go-grpc_out=$GOPATH/src *.proto
```
- `--go_out` → mesaj kodu (user.pb.go); `--go-grpc_out` → server/client
  kodu (user_grpc.pb.go)

**Yaradılan stub (görünüş):**
```go
type UserServiceClient interface {
    GetUser(ctx context.Context, in *UserRequest, opts ...grpc.CallOption) (*User, error)
}
func NewUserServiceClient(cc grpc.ClientConnInterface) UserServiceClient

type UserServiceServer interface {
    GetUser(context.Context, *UserRequest) (*User, error)
    mustEmbedUnimplementedUserServiceServer()
}
```

### 2. Server qurulması (3 addım)
**Kitabdan kod nümunəsi:**
```go
// 1) Implementasiya
type UserServer struct {
    pb.UnimplementedUserServiceServer   // forward-uyğunluq üçün embed
}

func (u *UserServer) GetUser(ctx context.Context, req *pb.UserRequest) (*pb.User, error) {
    fmt.Println("Server received:", req.String())
    return &pb.User{UserId: "John", Email: "john@gmail.com"}, nil
}

// 2)+3) Register + dinləmə
func main() {
    lis, err := net.Listen("tcp", "localhost:50051")
    if err != nil {
        panic(err)
    }
    s := grpc.NewServer()
    pb.RegisterUserServiceServer(s, &UserServer{})

    if err := s.Serve(lis); err != nil {   // bloklanır — sonsuz gözləyir
        panic(err)
    }
}
```
- RPC imzası: `func RPCName(ctx, req) (*Response, error)` — bütün
  metodlar implement olunmalı, yoxsa tip server olmur
- `pb.Unimplemented...Server` embed-i gələcək RPC-lər üçün sıçramaya
  qarşı sığorta

### 3. Client qurulması
**Kitabdan kod nümunəsi:**
```go
func main() {
    conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
    if err != nil {
        panic(err)
    }
    defer conn.Close()

    c := pb.NewUserServiceClient(conn)

    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()

    r, err := c.GetUser(ctx, &pb.UserRequest{UserId: "John"})
    if err != nil {
        panic(err)
    }
    fmt.Println("Client received:", r.String())
}
```

**Sub-kod izahı:**
- `grpc.Dial(adres, opts...)` → bağlantı; `WithInsecure` = şifrələməsiz
  (test üçün), `WithBlock` = bağlana qədər gözlə
- `NewUserServiceClient(conn)` → stub-dan hazır klient
- Sorğu kontekstlə (timeout) birlikdə göndərilir — Chapter 6.6

### 4. Server streaming
**Nədir:** Client bir sorğu göndərir, server cavabı **axınla** (stream) qaytarır.
HTTP/2 full-duplex-dən istifadə edir.

```protobuf
service NumService {
    rpc Rnd (NumRequest) returns (stream NumResponse);
}
```

**Server tərəfi:**
```go
func (n *NumServer) Rnd(req *pb.NumRequest, stream pb.NumService_RndServer) error {
    if req.N <= 0 {
        return errors.New("N must be greater than zero")
    }
    done := make(chan bool)
    go func() {
        for counter := 0; counter < int(req.N); counter++ {
            i := rand.Intn(int(req.To)-int(req.From)+1) + int(req.To)
            resp := pb.NumResponse{I: int64(i), Remaining: req.N - int64(counter)}
            stream.Send(&resp)          // hər element axına göndərilir
            time.Sleep(time.Second)
        }
        done <- true
    }()
    <- done
    return nil                            // nil → EOF → kanal bağlanır
}
```

**Client tərəfi (Recv döngüsü):**
```go
ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
defer cancel()

stream, err := c.Rnd(ctx, &pb.NumRequest{N: 5, From: 0, To: 100})
if err != nil {
    panic(err)
}

done := make(chan bool)
go func() {
    for {
        resp, err := stream.Recv()
        if err == io.EOF {
            done <- true
            return
        }
        if err != nil {
            panic(err)
        }
        fmt.Println("Received:", resp.String())
    }
}()
<- done
```
- `stream.Recv()` → EOF-na qədər oxu
- Kontekst axının bütün ömrünü idarə edir: timeout bitərsə axın bağlanır

### 5. Client streaming
**Nədir:** Client axınla data göndərir, server sonda bir cavab verir.

```protobuf
service NumService {
    rpc Sum (stream NumRequest) returns (NumResponse);
}
```

**Client:**
```go
stream, err := c.Sum(ctx)     // sorğu yox, yalnız kontekst

from, to := 1, 100
for i := from; i <= to; i++ {
    err = stream.Send(&pb.NumRequest{X: int64(i)})
    if err != nil {
        panic(err)
    }
}
result, err := stream.CloseAndRecv()   // bağla + cavabı gözlə
fmt.Printf("The sum from %d to %d is %d\n", from, to, result.Total)
```

**Server:**
```go
func (n *NumServer) Sum(stream pb.NumService_SumServer) error {
    var total int64 = 0
    counter := 0
    for {
        next, err := stream.Recv()
        if err == io.EOF {
            fmt.Printf("Received %d numbers sum: %d\n", counter, total)
            stream.SendAndClose(&pb.NumResponse{Total: total})
            return nil
        }
        if err != nil {
            return err
        }
        total = total + next.X
        counter++
    }
}
```
- `CloseAndRecv` (client) ↔ `SendAndClose` (server) cütü
- `Close` = cavab gözləmədən bağla

### 6. Bidirectional streaming
**Nədir:** Hər iki tərəf asinxron göndərə/ala bilər — chat protokolları üçün.

```protobuf
service ChatService {
    rpc SendTxt (stream ChatRequest) returns (stream StatsResponse);
}
```

**Client (iki goroutine — göndər/izlə):**
```go
func Chat(stream pb.ChatService_SendTxtClient, done chan bool) {
    t := time.NewTicker(time.Millisecond * 500)
    for {
        select {
        case <- done:
            return
        case <- t.C:
            stream.Send(&pb.ChatRequest{Txt: "Hello", Id: 1, To: 2})
        }
    }
}

func Stats(stream pb.ChatService_SendTxtClient, done chan bool) {
    for {
        stats, err := stream.Recv()
        if err != nil {
            panic(err)
        }
        fmt.Println(stats.String())
        if stats.TotalChar > 35 {
            done <- true
            stream.CloseSend()     // client bağlayır
            return
        }
    }
}

stream, err := c.SendTxt(context.Background())
done := make(chan bool)
go Stats(stream, done)
go Chat(stream, done)
<- done
```

**Server (asinxron sayğaç + periodik send):**
```go
func (c *ChatServer) SendTxt(stream pb.ChatService_SendTxtServer) error {
    var total int64 = 0
    go func() {                       // hər 2 saniyədə stats göndər
        for {
            t := time.NewTicker(time.Second * 2)
            select {
            case <- t.C:
                stream.Send(&pb.StatsResponse{TotalChar: total})
            }
        }
    }()
    for {                              // gələn mesajları oxu
        next, err := stream.Recv()
        if err == io.EOF {
            fmt.Println("Client closed")
            return nil
        }
        if err != nil {
            return err
        }
        fmt.Println("->", next.Txt)
        total = total + int64(len(next.Txt))
    }
}
```

### 7. Transcoding (grpc-gateway) — gRPC → REST/JSON
**Nədir:** HTTP/2-ni dəstəkləməyən client/proxy-lər üçün gRPC-nin
HTTP+JSON REST API kimi təqdim edilməsi. Reverse proxy HTTP sorğunu
gRPC sorğusuna çevirir.

**Annotations ilə .proto:**
```protobuf
import "google/api/annotations.proto";

service UserService {
    rpc Get (UserRequest) returns (User) {
        option(google.api.http) = {
            get: "/v1/user/{user_id}"       // {user_id} → UserRequest sahəsi
        };
    }
    rpc Create (User) returns (User) {
        option(google.api.http) = {
            post: "/v1/user"
            body: "*"                        // bütün body → User mesajı
        };
    }
}
```
- Əlavə bindinglər: `additional_bindings { get: "/v2/user/{user_id}" }`
- URL-i bir tərəfli dəyişməyin; versiyanı URL-də saxlayın (/v1/)

**Generasiya:**
```bash
>>> go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
        github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
        google.golang.org/protobuf/cmd/protoc-gen-go \
        google.golang.org/grpc/cmd/protoc-gen-go-grpc
>>> protoc -I . -I $GOPATH/src/.../googleapis \
    --go_out=plugins=grpc:$GOPATH/src \
    --grpc-gateway_out=logtostderr=true:$GOPATH/src *.proto
# user.pb.gw.go yaranır (gateway kodu)
```

**HTTP server (gateway):**
```go
func (u *UserServer) ServeHttp() {
    mux := runtime.NewServeMux()
    opts := []grpc.DialOption{grpc.WithInsecure()}
    endpoint := u.grpcAddr

    err := pb.RegisterUserServiceHandlerFromEndpoint(context.Background(),
        mux, endpoint, opts)     // gateway-i gRPC endpoint-ə bağla
    if err != nil {
        panic(err)
    }

    httpServer := &http.Server{Addr: u.httpAddr, Handler: mux}
    if err = httpServer.ListenAndServe(); err != nil {
        panic(err)
    }
}

func main() {
    us := UserServer{httpAddr: ":8080", grpcAddr: ":50051"}
    go us.ServeGrpc()    // gRPC :50051
    us.ServeHttp()        // HTTP/JSON :8080
}
```

**Test:**
```bash
>>> curl http://localhost:8080/v1/user/john
{"userId":"John","email":"john@gmail.com"}
>>> curl -d '{"user_id":"john","email":"john@gmail"}' http://localhost:8080/v1/user
{"userId":"john","email":"john@gmail"}
```

### 8. Interceptor-lar (gRPC middleware)
**Nədir:** Client↔server axını arasındaki qat — auth, tracing, validation.
4 növ: client/server × unary/streaming (baxılan: unary).

**Server interceptor (auth nümunəsi):**
```go
func AuthServerInterceptor(
    ctx context.Context,
    req interface{},
    info *grpc.UnaryServerInfo,
    handler grpc.UnaryHandler) (interface{}, error) {
    md, found := metadata.FromIncomingContext(ctx)   // gələn metadata
    if !found {
        return nil, status.Errorf(codes.InvalidArgument, "metadata not found")
    }
    password, found := md["password"]
    if !found {
        return nil, status.Errorf(codes.Unauthenticated, "password not found")
    }
    if password[0] != "go" {
        return nil, status.Errorf(codes.Unauthenticated, "password not valid")
    }
    h, err := handler(ctx, req)    // auth uğurlu → işləməyə burax
    return h, err
}

func withAuthServerInterceptor() grpc.ServerOption {
    return grpc.UnaryInterceptor(AuthServerInterceptor)
}

func main() {
    // ...
    s := grpc.NewServer(withAuthServerInterceptor())
    // ...
}
```
- `UnaryServerInterceptor` imzası: (ctx, req, info, handler)
- `status.Errorf(codes..., "...")` → gRPC status kodları ilə xəta

**Client tərəfi (metadata əlavə etmə):**
```go
ctx, cancel := context.WithTimeout(context.Background(), time.Second)
ctx = metadata.AppendToOutgoingContext(ctx, "password", "go")
defer cancel()
r, err := c.GetUser(ctx, &pb.UserRequest{UserId: "John"})
```

**Client interceptor (logging metadata):**
```go
func ClientLoggerInterceptor(
    ctx context.Context,
    method string,
    req, reply interface{},
    cc *grpc.ClientConn,
    invoker grpc.UnaryInvoker,
    opts ...grpc.CallOption) error {

    os := runtime.GOOS
    zone, _ := time.Now().Zone()

    ctx = metadata.AppendToOutgoingContext(ctx, "os", os)
    ctx = metadata.AppendToOutgoingContext(ctx, "zone", zone)

    return invoker(ctx, method, req, reply, cc, opts...)  // axını davam etdir
}

func withUnaryClientLoggerInterceptor() grpc.DialOption {
    return grpc.WithUnaryInterceptor(ClientLoggerInterceptor)
}

conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(),
    grpc.WithBlock(), withUnaryClientLoggerInterceptor())
```

**Server tərəfində metadata-nın oxunması:**
```go
md, found := metadata.FromIncomingContext(ctx)
if found {
    os, _ := md["os"]
    zone, _ := md["zone"]
    fmt.Printf("Request from %s using %s\n", zone, os)  // [CET] [darwin]
}
```

## Əsas terminlər
- gRPC — HTTP/2 üzərində RPC framework
- IDL (Interface Definition Language) — service imzalarının dili (.proto)
- Stub — protoc-un yaratdığı client/server kodu
- Unary RPC — tək sorğu/tək cavab
- Server/client streaming — bir tərəfli axın
- Bidirectional streaming — ikili istiqamətli axın
- Transcoding — gRPC-nin HTTP/JSON REST-ə çevrilməsi
- grpc-gateway — transcoding üçün reverse proxy generatoru
- Interceptor — gRPC middleware qatı
- Metadata — sorğu başlıqlarının gRPC analoqu

## Praktik nəticə
gRPC axını: .proto-da service yaz → protoc + go-grpc plugin → server-də
metodları implement et + `Register...Server` + `Serve`; client-da `Dial` +
`New...Client` + kontekstli çağırış. Streaming-də EOF (`io.EOF`) axının
bağlanma siqnalıdır. Brauzer/REST müştəriləri üçün grpc-gateway annotations
ilə eyni API-nı HTTP/JSON ilə təqdim edin. Auth/logging üçün server
(`UnaryInterceptor`) və client (`WithUnaryInterceptor`) interceptor-ları —
metadata API-si ilə.

## Mənbə
Pages: 282-325 (PDF 282-325)
