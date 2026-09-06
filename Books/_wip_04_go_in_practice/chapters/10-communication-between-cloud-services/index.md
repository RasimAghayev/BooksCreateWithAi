# Chapter 10 — Communication between cloud services (Technique 62-65)

## Bu chapter nədən bəhs edir?

Microservice kommunikasiyası: connection reuse (keep-alive, TCP slow-start, body close tələsi), reflection-sız JSON (codecgen), protocol buffers (protobuf + protoc), gRPC (proto3, service, context cancellation).

## Əsas fikirlər

### 1. Microservices əsasları
**Xüsusiyyətlər:** (1) Tək iş görür (UNIX fəlsəfəsi); (2) Elastic — horizontal scale; (3) Resilient — instance-lər ölsə də servis yaşayır.

**Transcoding nümunəsi:**
```
User → UI → API server → File storage + Message queue
                          Message queue → Transcoder → File storage → Notifications
```
- Hər hissə ayrıca deploy/scale olunur (transcoder yükə görə, API fərqli)
- Müxtəlif dillərdə yazıla bilər; storage SaaS kimi xaricdən gələ bilər
- **Kommunikasiya = bottleneck** — Google öz network layer-ini YAZIB (böyükdə performans kritikdir)

### TECHNIQUE 62: Connection reuse
**Problem:** hər request = yeni connection → TLS handshake + **TCP slow-start** (congestion-control ramp) = bir mesaj bir neçə round-trip.

**Həll:** connection-i YENİDƏN İSTİFADƏ ET — 1 dəfə bağla, N request göndər.

**Go-nun daxili dəstəyi:** http.DefaultClient/DefaultTransport keep-alive AÇIQDIR (30s). Pozan hallar:

**a) Custom Transport-də Dial qeyri-mövcuddursa:**
```go
// KÖHNE (Go sənədlərindən!): keep-alive SÖNÜLÜ
tr := &http.Transport{
    TLSClientConfig:    &tls.Config{RootCAs: pool},
    DisableCompression: true,
}
client := &http.Client{Transport: tr}

// DÜZGÜN: Dial keep-alive KÖMÜR:
tr := &http.Transport{
    TLSClientConfig:    &tls.Config{RootCAs: pool},
    DisableCompression: true,
    Dial: (&net.Dialer{
        Timeout:   30 * time.Second,
        KeepAlive: 30 * time.Second,       // ← DefaultTransport konfiqi
    }).Dial,
}
client := &http.Client{Transport: tr}
```
- `DisableKeepAlives: false` ≠ "açıqdır"; "SEÇƏ BİLƏRSƏN" deməkdir
- HTTP keep-alive (server tərəfi, protokol) ≠ TCP keep-alive (OS tərəfi) — DisableKeepAlives hər ikisini öldürür

**b) Body bağlanmayanda (pəncərə bərpası):**
```go
// PİS: defer body-ni SONDRA bağlayır → 2-ci request YENİ connection açır!
r, err := http.Get("http://example.com")
defer r.Body.Close()              // main sonuna qədər AÇIQ qalır!
o, _ := ioutil.ReadAll(r.Body)
r2, err := http.Get("http://example.com/foo")   // köhnə connection MƏŞĞUL → yeni açılır
defer r2.Body.Close()

// DÜZGÜN: oxudugun KİMİ bağla:
r, err := http.Get("http://example.com")
o, err := ioutil.ReadAll(r.Body)
r.Body.Close()                    // ← dərhal! → connection pool-a qaytar
r2, err := http.Get("http://example.com/foo")    // eyni connection reuse!
o, err = ioutil.ReadAll(r2.Body)
r2.Body.Close()
```
HTTP/1.1-də pipelining az istifadə olunurdu → serial serial → body bağlı olmayınca connection boş deyil. **defer burada düzgün deyil!** (HTTP/2 multiplexing bunu dəyişir, amma close yə hər zaman yaxşı vərdişdir.)

### TECHNIQUE 63: Faster JSON — codecgen
**Problem:** encoding/json hər marshal-da **REFLECTION** — eyni strukturda təkrar-təkrar tip hesablanması.

**Həll:** kod GENERASIYASI — tipi 1 dəfə hesabla, runtime-da skip.

```go
//go:generate codecgen -o user_generated.go user.go    // go generate direktivi!
package user

type User struct {
    Name  string `codec:"name"`         // codec tag-ləri (json → codec)
    Email string `codec:",omitempty"`   // boşdursa çıxart
}
```

```bash
$ go get -u github.com/ugorji/go/codec/codecgen
$ codecgen -o user_generated.go user.go     # VEYA:
$ go generate ./...                          // go:generate şərhlərini icra et
```

**İstifadə:**
```go
jh := new(codec.JsonHandle)             // JSON handler (Binc/MsgPack/CBOR da var!)
u := &user.User{Name: "Inigo Montoya", Email: "inigo@montoya.example.com"}

var out []byte
err := codec.NewEncoderBytes(&out, jh).Encode(&u)
fmt.Println(string(out))     // {"name":"Inigo Montoya","Email":"..."}

var u2 user.User
err = codec.NewDecoderBytes(out, jh).Decode(&u2)
```
- Generasiya `CodecEncodeSelf`/`CodecDecodeSelf` metodları yaradır → codec bunları tapanda reflection-a EHTİYAC YOX
- **Qeyd:** codecgen main paketlərində işləmir → struct-lar ayrı paketdə
- Handleanın öz tag-i yoxdursa sahə adı BÖYÜK hərflə çıxır (Email); `codec:"name"` kiçik yazır

### TECHNIQUE 64: Protocol buffers
**Nədir:** Google-un binary serializasiyası — JSON/XML-dən **kiçik** + **sürətli**. Transport müstəqildir (HTTP, RPC, queue, fayl).

**user.proto:**
```protobuf
package chapter10;              // ad toqquşması qarşısı (Go paketindən FƏRLİ)
message User {
  required string name = 1;     // = 1: binary tag ( sahə ID)
  required int32 id = 2;        // int64, float32/64, bool...
  optional string email = 3;    // optional sahə
}
```

**Quraşdırma:**
```bash
# 1. protoc compiler: developers.google.com/protocol-buffers/docs/downloads
# 2. Go plugin:
$ go get -u github.com/golang/protobuf/protoc-gen-go
# 3. Generate:
$ protoc -I=. --go_out=. ./user.proto
# -I input dir; --go_out output; .proto faylı
```

**Server:**
```go
pb "github.com/Masterminds/go-in-practice/chapter10/userpb"   // generated
"github.com/golang/protobuf/proto"

func handler(res http.ResponseWriter, req *http.Request) {
    u := &pb.User{
        Name:  proto.String("Inigo Montoya"),    // POINTER-lar! proto.String/Int32
        Id:    proto.Int32(1234),
        Email: proto.String("inigo@montoya.example.com"),
    }
    body, err := proto.Marshal(u)
    if err != nil {
        http.Error(res, err.Error(), http.StatusInternalServerError)
        return
    }
    res.Header().Set("Content-Type", "application/x-protobuf")
    res.Write(body)
}
```

**Client:**
```go
res, err := http.Get("http://localhost:8080")
defer res.Body.Close()
b, _ := ioutil.ReadAll(res.Body)

var u pb.User
err = proto.Unmarshal(b, &u)

fmt.Println(u.GetName())       // Get* metodları pointer-lardan dəyəri alır
fmt.Println(u.GetId())
fmt.Println(u.GetEmail())
```
- Sahələr **pointer** → `proto.String()` / `proto.Int32()` ilə yarat
- Oxu: `u.GetName()` (Get prefiksli metodlar)
- **XƏBƏRDARLIQ:** user məlumatı TLS üzərindən (bu nümunə sadəlik üçün HTTP)

### TECHNIQUE 65: gRPC — RPC + protobuf
**REST vs RPC:** REST = resurs + HTTP verb; RPC = **prosedur çağırışı** ("server-i restart et" kimi əməliyyat semantikası üçün ideal).

**hello.proto (proto3):**
```protobuf
syntax = "proto3";              // MÜTLƏQ ilk sətir! gRPC proto3 tələb edir
package chapter10;

service Hello {                                 // SERVIS təyini
  rpc Say (HelloRequest) returns (HelloResponse) {}
}
message HelloRequest {
  string name = 1;                              // proto3-də required/optional YOX
}
message HelloResponse {
  string message = 1;
}
```

```bash
$ protoc -I=. --go_out=plugins=grpc:. ./hello.proto    # grpc PLUGIN!
# plugins=grpc olmadan service stub yaradılmır
```

**Server:**
```go
pb "github.com/Masterminds/go-in-practice/chapter10/hellopb"
"golang.org/x/net/context"                // Go 1.7-dən standart kitabxanada
"google.golang.org/grpc"

type server struct{}

func (s *server) Say(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
    msg := "Hello " + in.Name + "!"
    return &pb.HelloResponse{Message: msg}, nil
}

func main() {
    l, err := net.Listen("tcp", ":55555")
    ...
    s := grpc.NewServer()
    pb.RegisterHelloServer(s, &server{})   // interfeys qeydiyyatı
    s.Serve(l)
}
```

**Client:**
```go
conn, err := grpc.Dial("localhost:55555", grpc.WithInsecure())   // productionda TLS!
defer conn.Close()
c := pb.NewHelloClient(conn)               // generated client

hr := &pb.HelloRequest{Name: "Inigo Montoya"}
r, err := c.Say(context.Background(), hr)   // Say imzası serverdəki kimi!
fmt.Println(r.Message)                      // "Hello Inigo Montoya!"
```

**Context — cancellation:**
```go
// Client:
ctx, cancel := context.WithCancel(context.Background())
defer cancel()                   // funksiya bitəndə SIQNAL

// Server (uzun RPC daxilində):
select {
case <-ctx.Done():               // client ləğv etdisə
    return nil, ctx.Err()        // "canceled" error-u
}
```
- Context deadline/cancel/values daşıyır — xətt boyu
- HTTP/2 (gRPC default) → connection reuse + multiplexing avtomatik

**gRPC PRO/CON:**
| Pro | Con |
|---|---|
| Protobuf: kiçik + sürətli | Binary — insan oxuya bilmir |
| Context: cancel/timeout | Interfeys bilmək lazım (yalnız payload kifayət etmir) |
| Prosedur semantikası | Dərin inteqrasiya — security həssaslığı |
| Transport semantikası azad | İctimai API üçün uyğun DEYİL (daxili microservice-lər üçün ideal) |

## Kommunikasiya performans pilləsi
```
JSON + yeni connection (ən yavaş)
  → JSON + keep-alive (body-ni vaxtında bağla!)
    → codecgen JSON (reflection-sız)
      → Protobuf + HTTP
        → gRPC (protobuf + HTTP/2 + context) (ən sürətli)
```

## Əsas terminlər
- Microservices (single action, elastic, resilient)
- Transcoding pipeline (UI → API → queue → transcoder → storage → notifications)
- Connection Reuse / Keep-Alive
- TCP Slow-Start (congestion control ramp)
- HTTP keep-alive vs TCP keep-alive
- DisableKeepAlives
- net.Dialer{KeepAlive}
- Body Close tələsi (defer-in zərəri)
- Pipelining vs Serial (HTTP/1) vs Multiplexing (HTTP/2)
- Reflection Overhead
- codecgen / ugorji codec / JsonHandle
- `//go:generate` directive
- Protocol Buffers (.proto / message / required/optional / field tag)
- protoc / protoc-gen-go plugin
- proto.String/Int32 (pointer wrapper-lər)
- proto.Marshal / Unmarshal
- application/x-protobuf
- gRPC / proto3 (syntax ilk sətir)
- service / rpc təyini
- plugins=grpc
- RegisterHelloServer / NewHelloClient
- context.Background / WithCancel / ctx.Done() / ctx.Err()
- RPC vs REST semantikası

## Praktik nəticə
- Default client/transport keep-alive-a HAZIRDIR; custom Transport yaradırsansa Dial-ə KeepAlive qoy — sənəddəki köhnə nümunələr bu addımı qaçırır.
- HTTP body-ni `defer`-lə yox, OXUDUQDAN SONRA bağla — 2-ci request eyni connection-dan gedir; defer yeni connection məcburiyyəti yaradır.
- Təkrar eyni strukturu marshal edirsənsə codecgen (və ya oxşarı) reflection-un qarşısını alır.
- JSON hələ də ictimai API standartıdır; protobuf — daxili microservice-lərin ölçüsü və sürəti üçün.
- proto sahə tag-ləri (1, 2, 3) sabit saxlanmalı — binary uyğunluq asılıdır.
- gRPC = proto3 + service + protoc plugins=grpc; daxili servis üçün RPC, xarici API üçün REST — kitabın aydın tövsiyəsi.
- Context WithCancel → uzun RPC-lərdə client-in getməsini server-ə bildir; server ctx.Done() seçir.

## Mənbə
Pages: 258-275 (PDF), book pages 235-252
