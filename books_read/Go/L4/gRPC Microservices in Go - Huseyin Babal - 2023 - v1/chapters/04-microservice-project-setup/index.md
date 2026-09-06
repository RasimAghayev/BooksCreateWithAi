# Chapter 4 — Microservice project setup (Mikroservis layihə qurulumu)

## Bu chapter nədən bəhs edir?
Hexagonal (ports & adapters) arxitekturasının Order servisinin tam implementasiyası ilə: domain modeli, API core, portlar (interfeyslər), DB adapter (GORM+MySQL), gRPC adapter, 12-factor konfiqurasiya (env vars), dependency injection və grpcurl ilə endpoint testi.

## Əsas fikirlər

### 1. Hexagonal arxitekturanın qovluq strukturu
```
order/
├── cmd/main.go                          # entry point + DI
├── config/config.go                     # env var oxunuşu
├── internal/
│   ├── ports/                           # MÜQAVİLƏLƏR (interfeyslər)
│   │   ├── api.go                       #   APIPort
│   │   └── db.go                        #   DBPort
│   ├── application/core/
│   │   ├── api/api.go                   # biznes məntiqi
│   │   └── domain/order.go             # domain modeli
│   └── adapters/
│       ├── db/db.go                     # GORM implementasiyası
│       └── grpc/{grpc.go, server.go}   # gRPC implementasiyası
```
**Qayda:** Application/Port/Adapter `internal/` içində; core-dan xaricə asılılıq (xarici hexaqon daxilə asılıdır).

### 2. Domain modeli (core)
```go
type OrderItem struct {
    ProductCode string  `json:"product_code"`
    UnitPrice   float32 `json:"unit_price"`
    Quantity    int32   `json:"quantity"`
}

type Order struct {
    ID         int64       `json:"id"`
    CustomerID int64       `json:"customer_id"`
    Status     string      `json:"status"`
    OrderItems []OrderItem `json:"order_items"`
    CreatedAt  int64       `json:"created_at"`
}

func NewOrder(customerId int64, orderItems []OrderItem) Order {
    return Order{
        CreatedAt:  time.Now().Unix(),
        Status:     "Pending",
        CustomerID: customerId,
        OrderItems: orderItems,
    }
}
```

### 3. Application core (API)
```go
type Application struct {
    db ports.DBPort            // INTERFACE-ə asılı, konkret DB-yə YOX
}

func NewApplication(db ports.DBPort) *Application {
    return &Application{db: db}
}

func (a Application) PlaceOrder(order domain.Order) (domain.Order, error) {
    err := a.db.Save(&order)
    if err != nil { return domain.Order{}, err }
    return order, nil
}
```

### 4. Portlar — müqavilə interfeysləri
```go
type APIPort interface {
    PlaceOrder(order domain.Order) (domain.Order, error)
}

type DBPort interface {
    Get(id string) (domain.Order, error)
    Save(*domain.Order) error
}
```

### 5. DB adapter (GORM + MySQL)
**Qurulum:** `go get -u gorm.io/gorm` + `go get -u gorm.io/driver/mysql`

**DB modelləri (1:N relation):**
```go
type Order struct {
    gorm.Model                 // ID, CreatedAt, UpdatedAt, DeletedAt
    CustomerID int64
    Status     string
    OrderItems []OrderItem     // 1:N — GORM avtomatik join
}
type OrderItem struct {
    gorm.Model
    ProductCode string
    UnitPrice   float32
    Quantity    int32
    OrderID     uint           // back-reference
}
```
**Adapter:**
```go
type Adapter struct{ db *gorm.DB }

func NewAdapter(dataSourceUrl string) (*Adapter, error) {
    db, openErr := gorm.Open(mysql.Open(dataSourceUrl), &gorm.Config{})
    if openErr != nil {
        return nil, fmt.Errorf("db connection error: %v", openErr)
    }
    if err := db.AutoMigrate(&Order{}, &OrderItem{}); err != nil {
        return nil, fmt.Errorf("db migration error: %v", err)
    }
    return &Adapter{db: db}, nil
}
```
**Get/Save:** DB entity ↔ domain model dönüşümü ilə:
```go
func (a Adapter) Get(id string) (domain.Order, error) {
    var orderEntity Order
    res := a.db.First(&orderEntity, id)
    // entity → domain çevirmə (orderItems loop ilə)
    return order, res.Error
}

func (a Adapter) Save(order *domain.Order) error {
    // domain → entity çevir; a.db.Create(&orderModel)
    // uğurda order.ID = int64(orderModel.ID)
}
```

### 6. gRPC adapter
**Adapter:**
```go
type Adapter struct {
    api  ports.APIPort                     // core-a asılılıq
    port int
    order.UnimplementedOrderServer          // forward compatibility şərti
}
```
**Create handler:**
```go
func (a Adapter) Create(ctx context.Context, request *order.CreateOrderRequest) (*order.CreateOrderResponse, error) {
    var orderItems []domain.OrderItem
    for _, orderItem := range request.OrderItems {
        orderItems = append(orderItems, domain.OrderItem{
            ProductCode: orderItem.ProductCode,
            UnitPrice:   orderItem.UnitPrice,
            Quantity:    orderItem.Quantity,
        })
    }
    newOrder := domain.NewOrder(request.UserId, orderItems)
    result, err := a.api.PlaceOrder(newOrder)   // core-a ötür
    if err != nil { return nil, err }
    return &order.CreateOrderResponse{OrderId: result.ID}, nil
}
```
**Server Run:**
```go
func (a Adapter) Run() {
    listen, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
    if err != nil { log.Fatalf(...) }
    grpcServer := grpc.NewServer()
    order.RegisterOrderServer(grpcServer, a)
    if config.GetEnv() == "development" {
        reflection.Register(grpcServer)     // yalnız dev-də (grpcurl üçün)
    }
    if err := grpcServer.Serve(listen); err != nil { log.Fatalf(...) }
}
```

### 7. Konfiqurasiya — 12-factor env vars
**ENV / DATA_SOURCE_URL / APPLICATION_PORT** — çatmayan varsa app FAIL-FAST (boş dəyərlə səssiz işləmək daha təhlükəli!):
```go
func getEnvironmentValue(key string) string {
    if os.Getenv(key) == "" {
        log.Fatalf("%s environment variable is missing.", key)
    }
    return os.Getenv(key)
}
```

### 8. Dependency injection + main
```go
func main() {
    dbAdapter, err := db.NewAdapter(config.GetDataSourceURL())
    if err != nil { log.Fatalf(...) }
    application := api.NewApplication(dbAdapter)        // DB → core
    grpcAdapter := grpc.NewAdapter(application, config.GetApplicationPort())
    grpcAdapter.Run()                                    // core → gRPC
}
```
**Zəncir:** DB adapter → Application → gRPC adapter → Run.

### 9. Local işə salma (Docker MySQL + grpcurl)
```bash
docker run -p 3306:3306 \
    -e MYSQL_ROOT_PASSWORD=verysecretpass \
    -e MYSQL_DATABASE=order mysql

DATA_SOURCE_URL=root:verysecretpass@tcp(127.0.0.1:3306)/order \
APPLICATION_PORT=3000 \
ENV=development \
go run cmd/main.go
```
**grpcurl testi:**
```bash
grpcurl -d '{"user_id": 123, "order_items": [{"product_code": "prod", "quantity": 4, "unit_price": 12}]}' \
  -plaintext localhost:3000 Order/Create
# → {"orderId": "1"}
```
`-plaintext` — TLS-i deaktiv; reflection yalnız dev-də açıq olduğundan localhost-da işləyir.

## Əsas terminlər
- Hexagonal architecture (ports & adapters) — altıbucaqlı arxitektura
- Port — qatlararası müqavilə interfeysi
- Adapter — portun konkret implementasiyası
- Domain model — biznes obyekti
- GORM — Go ORM; gorm.Model (metadata), AutoMigrate
- 1:N relation — Order → OrderItems + OrderID back-ref
- 12-factor app — env var konfiqurasiya metodologiyası
- Fail-fast — çatmayan config = dərhal ölüm
- UnimplementedOrderServer — forward compatibility baza
- grpcurl — gRPC üçün curl
- Reflection — dev rejimində gRPC introspection

## Praktik nəticə
1. Core-dan başla → sonra xarici qatlar; asılılıq HƏMİŞƏ interfeysə (port), konkretə YOX.
2. Domain ↔ DB-entity ↔ gRPC-request dönüşümləri ayrı saxla — 3 fərqli model, bir-birinə qarışmasın.
3. Config fail-fast: çatmayan env var proqramı dərhal dayandırsın.
4. Reflection yalnız development-də — production-da söndür.
5. Local dev: Docker MySQL + env vars + grpcurl üçlüyü — 2 dəqiqədə işlək servis.
6. AutoMigrate prototip üçün yaxşıdır; production-da versioned migration aləti (Ch 8 qeydləri).

## Mənbə
Pages: 47-62 (PDF səh. 67-82)
