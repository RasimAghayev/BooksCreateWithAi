# Chapter 5 — Interservice communication (Servislərarası kommunikasiya)

## Bu chapter nədən bəhs edir?
Load balancing strategiyalarına (server/client-side), Payment servisinin Order servisinə inteqrasiyasına (port/adapter), gRPC status kodlarına, xəta mesajlarının detalları ilə qaytarılmasına (errdetails) və client tərəfində xətaların emalına.

## Əsas fikirlər

### 1. Load balancing strategiyaları
| | Server-side LB | Client-side LB |
|---|---|---|
| **Mexanizm** | client → LB → backend-lər | client birbaşa backend-lərə |
| **Pros** | asan konfiq; backend bilinmir | əlavə hop yox → sürət; SPOF yox |
| **Cons** | latency; throughput limit | konfiq çətin; health/load izləmə; dil-özəl implementasiya |
| **Alqoritmlər** | — | round-robin (default), load-report əsaslı |

**Kitabın seçimi:** Kubernetes istifadə olunacağı üçün server-side LB (K8s Service yerleşdirir).

### 2. Payment modulundan asılılıq
```bash
go get -u github.com/huseyinbabal/microservices-proto/golang/payment@v1.0.38
# sabit versiya — gözlənilməz davranış qarşısı
paymentClient := payment.NewPaymentClient(conn, opts)
paymentClient.Create(ctx, &CreatePaymentRequest{})
```

### 3. Payment port və adapter (driven side)
**Port (müqavilə):**
```go
type PaymentPort interface {
    Charge(*domain.Order) error
}
```
**Adapter (implementasiya):**
```go
type Adapter struct {
    payment payment.PaymentClient   // generasiya olunmuş stub
}

func NewAdapter(paymentServiceUrl string) (*Adapter, error) {
    var opts []grpc.DialOption
    opts = append(opts,
        grpc.WithTransportCredentials(insecure.NewCredentials()))  // TLS-siz ( sadəlik üçün)
    conn, err := grpc.Dial(paymentServiceUrl, opts...)
    if err != nil { return nil, err }
    defer conn.Close()
    client := payment.NewPaymentClient(conn)
    return &Adapter{payment: client}, nil
}

func (a *Adapter) Charge(order *domain.Order) error {
    _, err := a.payment.Create(context.Background(), &payment.CreatePaymentRequest{
        UserId:     order.CustomerID,
        OrderId:    order.ID,
        TotalPrice: order.TotalPrice(),
    })
    return err
}
```
**Domain hesablama:**
```go
func (o *Order) TotalPrice() float32 {
    var totalPrice float32
    for _, orderItem := range o.OrderItems {
        totalPrice += orderItem.UnitPrice * float32(orderItem.Quantity)
    }
    return totalPrice
}
```

### 4. Konfiqurasiya + DI (main.go)
```go
func GetPaymentServiceUrl() string {
    return getEnvironmentValue("PAYMENT_SERVICE_URL")   // fail-fast qaydası qalır
}

func main() {
    dbAdapter, err := db.NewAdapter(config.GetDataSourceURL())
    paymentAdapter, err := payment.NewAdapter(config.GetPaymentServiceUrl())
    application := api.NewApplication(dbAdapter, paymentAdapter)  // 2 adapter inyeksiya
    grpcAdapter := grpc.NewAdapter(application, config.GetApplicationPort())
    grpcAdapter.Run()
}
```

### 5. Core-da payment çağırışı
```go
func (a Application) PlaceOrder(order domain.Order) (domain.Order, error) {
    if err := a.db.Save(&order); err != nil { return domain.Order{}, err }
    if paymentErr := a.payment.Charge(&order); paymentErr != nil {
        return domain.Order{}, paymentErr      // (bax: error handling aşağıda)
    }
    return order, nil
}
```

### 6. gRPC status kodları
| Kod | Məna | Nümunə |
|---|---|---|
| OK | uğur | — |
| CANCELLED | client ləğv etdi | ilk-cavab yarışı |
| INVALID_ARGUMENT | client-in input-u yanlış | boş order ID |
| DEADLINE_EXCEEDED | deadline bitdi | 5s limit, 6s cavab |
| NOT_FOUND | resurs yoxdur | mövcud olmayan order ID |
| ALREADY_EXISTS | duplikat | eyni email-ə user |
| PERMISSION_DENIED | icazə yoxdur | admin resursu |
| RESOURCE_EXHAUSTED | quota/limit bitib | SaaS limiti |
| INTERNAL | server daxili xəta | — |

### 7. Xəta kodu + mesaj qaytarma
**Pis (kodsuz):**
```go
return nil, err   // → Code: Unknown (!)
```
```bash
ERROR:
  Code: Unknown
  Message: failed to charge the customer
```
**Yaxşı (status.Errorf):**
```go
err = status.Errorf(codes.InvalidArgument,
    fmt.Sprintf("failed to charge user: %d", request.UserId))
return nil, err
```
```bash
ERROR:
  Code: InvalidArgument
  Message: failed to charge user: 123
```

### 8. Errors with details — errdetails
**Server tərəfində zəngin xəta:**
```go
paymentErr := a.payment.Charge(&order)
if paymentErr != nil {
    st, _ := status.FromError(paymentErr)              // payment xətasından status
    fieldErr := &errdetails.BadRequest_FieldViolation{
        Field:       "payment",                         // HANSI sahənin xətası
        Description: st.Message(),
    }
    badReq := &errdetails.BadRequest{}
    badReq.FieldViolations = append(badReq.FieldViolations, fieldErr)
    orderStatus := status.New(codes.InvalidArgument, "order creation failed")  // kök status
    statusWithDetails, _ := orderStatus.WithDetails(badReq)
    return domain.Order{}, statusWithDetails.Err()
}
```
**Cavab:**
```json
ERROR:
  Code: InvalidArgument
  Message: order creation failed
  Details:
  1) {"@type":"type.googleapis.com/google.rpc.BadRequest",
      "fieldViolations":[{"field":"payment",
      "description":"failed to charge. invalid billing address"}]}
```

### 9. Client tərəfi details emalı
```go
st := status.Convert(paymentErr)                    // FromError-ın qardaşı
var allErrors []string
for _, detail := range st.Details() {
    switch t := detail.(type) {
    case *errdetails.BadRequest:
        for _, violation := range t.GetFieldViolations() {
            allErrors = append(allErrors, violation.Description)
        }
    }
}
fieldErr := &errdetails.BadRequest_FieldViolation{
    Field:       "payment",
    Description: strings.Join(allErrors, "\n"),     // bütün detallar birləşdirilir
}
// ... yuxarıdakı kimi yeni status qaytarılır
```
**Fərq:** `FromError` sadə status; `Convert` hər zaman status qaytarır (err status olmasa belə) + Details slice-ı iterasiya oluna bilir.

### 10. İki servisin lokal işə salınması
```bash
# Payment (port 3001):
DB_DRIVER=mysql DATA_SOURCE_URL=root:pass@tcp(127.0.0.1:3306)/payment \
APPLICATION_PORT=3001 ENV=development go run cmd/main.go

# Order (port 3000 + payment ünvanı):
DB_DRIVER=mysql DATA_SOURCE_URL=root:pass@tcp(127.0.0.1:3306)/order \
APPLICATION_PORT=3000 ENV=development \
PAYMENT_SERVICE_URL=localhost:3001 go run cmd/main.go
```
**Prinsip:** Payment PUBLIC deyil — yalnız Order (şəbəkə içində) çatır; istifadəçi yalnız Order ilə danışır.

## Əsas terminlər
- Server-side / Client-side load balancing
- Round-robin — növbəli paylama
- Driven side — hexaqonun çıxış tərəfi (Order→Payment)
- PaymentPort / PaymentAdapter
- gRPC status kodları — codes.InvalidArgument və b.
- status.Errorf / status.New / status.FromError / status.Convert
- errdetails.BadRequest_FieldViolation — sahə-özəl xəta detalı
- WithDetails / st.Details() — status zənginləşdirmə
- insecure.NewCredentials — TLS-siz dial (dev)
- PAYMENT_SERVICE_URL — servis ünvanı env-i

## Praktik nəticə
1. Servis-arası asılılıqlar port (interface) + adapter (stub) cütü ilə — Order Payment-in İÇİNİ bilmir.
2. Xətaları heç vaxt xam `return nil, err` qaytarma — status.Errorf ilə kodlu; yoxsa "Unknown".
3. Çoxservisli xətalarda errdetails istifadə et — hansı SAHƏnin xətası olduğu görünsün.
4. Client tərəfdə Convert + Details loop-u — SDK-larda istifadəçiyə mənalı mesaj üçün.
5. Dependency versiyasını sabitlə (@v1.0.38) — "latest" gizli breaking riskidir.
6. K8s istifadə edirsənsə server-side LB kifayətdir — client-side yalnız xüsusi ehtiyacda.

## Mənbə
Pages: 65-78 (PDF səh. 83-96)
