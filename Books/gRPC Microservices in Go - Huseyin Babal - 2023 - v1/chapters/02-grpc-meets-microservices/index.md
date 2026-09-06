# Chapter 2 — gRPC meets microservices (gRPC mikroservislərlə görüşür)

## Bu chapter nədən bəhs edir?
Monolith vs mikroservis arxitekturaları müqayisəsinə, scale cube modelinə, Saga pattern-inə (choreography/orchestrator), service discovery-ə və gRPC+Protobuf ilə servis-kommunikasiyasının ilk addımlarına (proto → stub → dial).

## Əsas fikirlər

### 1. Monolithic arxitektura: müşahidələr
**Nədir:** UI + server + DB bir vahid paketdə (JAR/WAR/Go binary). İlk versiya üçün ƏLA — biznes domenini nonfunctional çətinliklərsiz öyrənmək.

**Problemlər (böyüdükcə):**
- **Development:** kod böyüyəndə IDE yüklənir; izolyasiyasız testlər hər dəyişiklikdə hamısı işə düşür; compile vaxtı artır
- **Deployment:** newsletter komponentində 1 kiçik dəyişiklik = TAM app deploy + HAMISININ testi; çox-komandalı monolithdə flaky test bütün deploy-u pozur
- **Scaling:** 16GB monolith × 2 = 32GB — amma lazım olan yalnız 2GB customer modulu idi! Texnologiya stack-ə uzunmüddədli öhdəlik

### 2. Scale cube — 3 ölçülü scalability modeli
| Ox | Model | Nə edir | Çatışmazlıq |
|---|---|---|---|
| **X-axis** | Eyni app-in N nüsxəsi + load balancer | hər instans 1/N yük | yaddaş israfı, kod mürəkkəbliyi azalmır |
| **Z-axis** | Nüsxələr DATA SUBSET-lərinə görə | xərc qənaəti + fault isolation | repartition mürəkkəbliyi |
| **Y-axis** | FUNKSIONAL parçalanma | mikroservislər! | — |

**Mikroservis = Y-axis scaling-in tətbiqi.**

### 3. Mikroservis arxitekturasının xarakteristikaları
Gevşək bağlı (loosely coupled); müstəqil deploy/scale; biznes-kabiliyyət fokusu; komanda-özəl ownership; texnologiya azadlığı; bir servis çöksə digərləri işləyir.

**Tövsiyə:** Monolith-dən başla → scale/test/release problemlərində yenidən qiymətləndir → funksional decompose.

### 4. Data consistency: Transaction → Saga
**Monolith:** Transaction (begin → biznes → commit/rollback) — `Order:create()` → `Payment:create()` + `Shipping:start()` — biri düşsə hamısı rollback.

**Mikroservisdə problem:** Hər servis ÖZ data store-unu saxlayır → bir 2PC (two-phase commit) seqvensial koordinasiya bahalı və köhnəmişdir → **Saga pattern**.

**Saga nədir:** Lokal transaksiyalar zənciri — hər servis lokal TX icra edir → növbəti servisi tetikleyən mesaj yayınlayır.

### 5. Choreography-based saga (rəqs-saga)
**Nədir:** Mərkəzi idarəçi YOX — servislər hadisələrlə (event) bir-birini tetikləyir.

**Order flow:** Order (PENDING) → `order_created` event → Payment charge → `payment_created/payment_failed` → Shipping → `shipping_created/failed` → Order statusu SUCCESS/FAILED.

**Saga tamamlanma üsulları:** (1) saga bitincə client-ə nəticə; (2) client polling (saga ID ilə); (3) WebSocket push.

**Kanal növləri:**
- **Command channel:** birbaşa növbəti servsə + `replyTo` parametri — publisher növbəti servisin ÜNVANINI bilməlidir
- **Pub/sub:** domain event → maraqlananlar abunə olur — amma broker single point of failure

### 6. Orchestrator-based saga (dirijor-saga)
**Nədir:** Orchestrator + iştirakçılar; orchestrator iştirakçılara NƏ edəcəyini DEYİR (command/req-resp).

**Order flow:** Create Order Saga → Payment:process() → uğursuzsa **compensation transaction** (refund) → Shipping:start() → uğursuzsa yuxarıdan aşağı rollback: Shipping:cancel() → Payment:refund() → Order:cancel().

**Kitabın seçimi:** request/response style — bu, gRPC-nin doğma rejimidir.

### 7. Service discovery (servis kəşfiyyatı)
| Növ | Mexanizm |
|---|---|
| **Client-side** | registry-ə self-register; client registry-dən servis adı ilə ünvan soruşur |
| **Server-side** | load balancer registry ilə inteqrasiya; client LB-ə qoşulur |
| **Kubernetes built-in** | servis ADI = ünvan (Chapter 8-də detal) |

### 8. gRPC ilə kommunikasiya — 5 addım
1. **.proto faylları** təyin et (mesajlar + servislər) — ayrıca repo və ya layihə daxilində
2. **Stub-lar generasiya et** (protoc ilə)
3. **Server tərəfi** biznes məntiqi (istənilən dillə)
4. **Client tərəfi** stub ilə bağlanır
5. İşə sal

**Proto nümunəsi (Listing 2.1):**
```protobuf
syntax = "proto3";
option go_package="github.com/huseyinbabal/grpc-microservices-in-go/listing_2.1/payment";

message CreatePaymentRequest {
    float price = 1;       // field NUMBERS — serialization açarı!
}
message CreatePaymentResponse {
    int64 bill_id = 1;
}
service Payment {
    RPC Create(CreatePaymentRequest) returns (CreatePaymentResponse) {}
}
```

**Stub generasiyası (protoc):**
```bash
protoc \
    --go_out=. \
    --go_opt=paths=source_relative \
    --go-grpc_out=. \
    --go-grpc_opt=paths=source_relative \
    payment.proto
```

### 9. 4 RPC üslubu
| Üslub | İmza | İstifadə |
|---|---|---|
| **Unary** | `Req → Resp` | klassik sorğu/cavab |
| **Server streaming** | `Req → stream Resp` | bulk data part-part |
| **Client streaming** | `stream Req → Resp` | çoxlu sorğu, bir yekun |
| **Bidirectional** | `stream Req → stream Resp` | asinxron davamlı axın (unique ID ilə xəta izləmə!) |

```protobuf
service Payment {
    Create(CreatePaymentRequest) returns (stream CreatePaymentResponse) {}  // server streaming
    Create(stream CreatePaymentRequest) returns (CreatePaymentResponse) {} // client streaming
    Create(stream CreatePaymentRequest) returns (stream CreatePaymentResponse) {} // bidi
}
```

### 10. Client tərəfi bağlanma (Listing 2.2)
```go
var opts []grpc.DialOption
conn, err := grpc.Dial(*serverAddr, opts...)   // Payment ünvanına dial
if err != nil { /* handle */ }
defer conn.Close()

payment := pb.NewPaymentClient(conn)                    // client stub yarat
result, err := payment.Create(ctx, &CreatePaymentRequest{})  // RPC çağırışı
if err != nil { /* handle */ }
```
**Addımlar:** stub import (Go dependency) → Dial → client → method çağırışı.

## Əsas terminlər
- Monolithic architecture — vahid paket
- Scale cube — X/Y/Z oxları
- Y-axis scaling — funksional decompose
- Transaction (begin/commit/rollback)
- 2PC (two-phase commit) — paylanmış atomic TX
- Saga — lokal TX zənciri
- Choreography saga — event-idarəli
- Orchestrator saga — mərkəzi dirijor + compensation transaction
- Command channel vs Pub/sub
- Correlation ID — hadisə zəncirinin izi
- Service discovery (client/server-side)
- protoc — protobuf compiler
- Unary/Server-stream/Client-stream/Bidirectional RPC
- grpc.Dial / pb.NewClient

## Praktik nəticə
1. Scale cube-dən istifadə et: problemin Y-axis-di (funksional) yoxsa X (surət) — mikroservis yalnız Y həllidir.
2. Paylanmış consistency üçün 2PC YOX — Saga: choreography (mərkəzsiz) vs orchestrator (mərkəzli + compensation).
3. Command channel-də replyTo + correlation ID mütləq — nəticə və izləmə üçün.
4. Proto field nömrələri serialization açarıdır — yenidən istifadə ETMƏ (Ch 3-də kompatibillik).
5. Streaming növünü data xarakterinə görə seç: bulk → server stream; asinxron davamlı → bidi.
6. K8s-də service discovery problemi yoxdur — servis adı kifayətdir.

## Mənbə
Pages: 14-27 (PDF səh. 31-44)
