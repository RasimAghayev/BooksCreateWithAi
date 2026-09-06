# Chapter 7 — Testing microservices (Mikroservislərin testi)

## Bu chapter nədən bəhs edir?
Test piramidinə (unit/integration/E2E), SUT anlayışına, testify mock-larına (mockery avtomatik generasiyası), Testcontainers ilə integration testlərə, Docker Compose stack üzərində end-to-end testlərə, coverage hesabına və CI pipeline inteqrasiyasına.

## Əsas fikirlər

### 1. Testing piramidi
```
     /E2E\        az, yavaş, baha (tam stack)
    /Integr\      orta (2 komponent + real DB)
   /Unit tests\   çox, sürətli, ucuz
```
Unit əsas payı daşıyır — sürətli feedback + test infra xərci az.

### 2. Unit test — SUT (System Under Test)
**SUT:** test edilən dəqiq yer — fayl YOX, SSENARİ. Bir funksiya şərt bloklarına görə çox ssenari daşıyır.

### 3. Mock-lar (testify/mock)
**MockedPayment (Ch 5-dən davam):**
```go
type mockedPayment struct {
    mock.Mock        // çağırış izləmə daxil olur
}
func (m *mockedPayment) Charge(order *domain.Order) error {
    args := m.Called(order)          // çağırışı qeyd et
    return args.Error(0)
}
```
**Happy path testi:**
```go
func Test_Should_Place_Order(t *testing.T) {
    payment := new(mockedPayment)
    db := new(mockedDb)
    payment.On("Charge", mock.Anything).Return(nil)   // davranış təyin et
    db.On("Save", mock.Anything).Return(nil)
    application := NewApplication(db, payment)
    _, err := application.PlaceOrder(domain.Order{
        CustomerID: 123,
        OrderItems: []domain.OrderItem{{ProductCode: "camera", UnitPrice: 12.3, Quantity: 3}},
    })
    assert.Nil(t, err)
}
```
**Xəta ssenariləri:**
```go
// DB xətası:
db.On("Save", mock.Anything).Return(errors.New("connection error"))
assert.EqualError(t, err, "connection error")

// Payment xətası — status+details assert:
payment.On("Charge", mock.Anything).Return(errors.New("insufficient balance"))
st, _ := status.FromError(err)
assert.Equal(t, st.Message(), "order creation failed")
assert.Equal(t, st.Details()[0].(*errdetails.BadRequest).FieldViolations[0].Description, "insufficient balance")
assert.Equal(t, st.Code(), codes.InvalidArgument)
```

### 4. Mockery — avtomatik mock generasiyası
```bash
mockery --all --keeptree     # bütün interfeyslər → *_mock.go
```
Interfeys dəyişəndə yenidən generate — əl ilə mock saxlamaq YOX.

### 5. Integration test — Testcontainers
**Nədir:** İki modulun birgə işi (DB adapter + REAL MySQL); Docker-da MySQL spin up → URL-i adapterə ötür.

**Test suite (testify/suite):**
```go
type OrderDatabaseTestSuite struct {
    suite.Suite                // assert-lər receiver-dən gəlir (o.Nil, o.Equal)
    DataSourceUrl string
}
```
| Hook | Vaxt | İstifadə |
|------|------|----------|
| SetupSuite | bütün testlərdən əvvəl | konteyner başlat |
| SetupTest | hər testdən əvvəl | mock hazırla |
| TearDownTest | hər testdən sonra | fayl sil |
| TearDownSuite | hamısından sonra | konteyner söndür |

**SetupSuite — MySQL konteyner:**
```go
func (o *OrderDatabaseTestSuite) SetupSuite() {
    ctx := context.Background()
    port := "3306/tcp"
    dbURL := func(port nat.Port) string {
        return fmt.Sprintf("root:s3cr3t@tcp(localhost:%s)/orders?charset=utf8mb4&parseTime=True&loc=Local", port.Port())
    }
    req := testcontainers.ContainerRequest{
        Image:        "docker.io/mysql:8.0.30",
        ExposedPorts: []string{port},
        Env: map[string]string{
            "MYSQL_ROOT_PASSWORD": "s3cr3t",
            "MYSQL_DATABASE":      "orders",
        },
        WaitingFor: wait.ForSQL(nat.Port(port), "mysql", dbURL).Timeout(time.Second * 30),
    }
    mysqlContainer, err := testcontainers.GenericContainer(ctx,
        testcontainers.GenericContainerRequest{ContainerRequest: req, Started: true})
    endpoint, _ := mysqlContainer.Endpoint(ctx, "")
    o.DataSourceUrl = fmt.Sprintf("root:s3cr3t@tcp(%s)/orders?...", endpoint)
}
```
**Testlər + runner:**
```go
func (o *OrderDatabaseTestSuite) Test_Should_Save_Order() {
    adapter, err := NewAdapter(o.DataSourceUrl)
    o.Nil(err)
    o.Nil(adapter.Save(&domain.Order{}))
}

func (o *OrderDatabaseTestSuite) Test_Should_Get_Order() {
    adapter, _ := NewAdapter(o.DataSourceUrl)
    order := domain.NewOrder(2, []domain.OrderItem{{ProductCode: "CAM", Quantity: 5, UnitPrice: 1.32}})
    adapter.Save(&order)
    ord, _ := adapter.Get(order.ID)
    o.Equal(int64(2), ord.CustomerID)
}

func TestOrderDatabaseTestSuite(t *testing.T) {
    suite.Run(t, new(OrderDatabaseTestSuite))
}
```

### 6. E2E test — Docker Compose stack
**Struktur:** `e2e/` ayrıca Go modul + `resources/docker-compose.yml` + `init.sql`.

**Docker Compose sahələri:**
| Sahə | Vəzifə |
|---|---|
| build | Dockerfile yolu |
| environment | env vars (port, DB URL) |
| depends_on + condition | asılılıq (service_healthy gözlə) |
| ports | xarici açılan portlar |
| healthcheck | hazırlıq yoxlaması (interval/timeout/retries) |
| volumes | lokal fayl mount |

**docker-compose.yml (3 servis):**
```yaml
version: "3.9"
services:
  mysql:
    image: "mysql:8.0.30"
    environment: { MYSQL_ROOT_PASSWORD: "s3cr3t" }
    volumes: ["./init.sql:/docker-entrypoint-initdb.d/init.sql"]
    healthcheck:
      test: ["CMD", "mysqladmin", "ping", "-h", "localhost", "-uroot", "-ps3cr3t"]
      interval: 5s
      timeout: 5s
      retries: 20
  payment:
    depends_on: { mysql: { condition: service_healthy } }
    build: ../../payment/
    environment:
      APPLICATION_PORT: 8081
      ENV: "development"
      DATA_SOURCE_URL: "root:s3cr3t@tcp(mysql:3306)/payments?..."
  order:
    depends_on: { mysql: { condition: service_healthy } }
    build: ../../order/
    ports: ["8080:8080"]
    environment:
      APPLICATION_PORT: 8080
      ENV: "development"
      DATA_SOURCE_URL: "root:s3cr3t@tcp(mysql:3306)/orders?..."
      PAYMENT_SERVICE_URL: "localhost:8081"
```
**Multistage Dockerfile:**
```dockerfile
FROM golang:1.18 AS builder
WORKDIR /usr/src/app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o payment ./cmd/main.go
FROM scratch
COPY --from=builder /usr/src/app/payment ./payment
CMD ["./payment"]
```
**init.sql:** `CREATE DATABASE IF NOT EXISTS payments; CREATE DATABASE IF NOT EXISTS orders;`

**E2E suite:**
```go
type CreateOrderTestSuite struct {
    suite.Suite
    compose *testcontainers.LocalDockerCompose
}

func (c *CreateOrderTestSuite) SetupSuite() {
    composeFilePaths := []string{"resources/docker-compose.yml"}
    identifier := strings.ToLower(uuid.New().String())      // unikal stack adı
    c.compose = testcontainers.NewLocalDockerCompose(composeFilePaths, identifier)
    execError := c.compose.WithCommand([]string{"up", "-d"}).Invoke()
    if execError.Error != nil { log.Fatalf("Could not run compose stack: %v", execError.Error) }
}

func (c *CreateOrderTestSuite) Test_Should_Create_Order() {
    var opts []grpc.DialOption
    opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
    conn, _ := grpc.Dial("localhost:8080", opts...)
    defer conn.Close()
    orderClient := order.NewOrderClient(conn)

    createOrderResponse, errCreate := orderClient.Create(context.Background(), &order.CreateOrderRequest{
        UserId: 23,
        OrderItems: []*order.OrderItem{{ProductCode: "CAM123", Quantity: 3, UnitPrice: 1.23}},
    })
    c.Nil(errCreate)

    getOrderResponse, errGet := orderClient.Get(context.Background(),
        &order.GetOrderRequest{OrderId: createOrderResponse.OrderId})
    c.Nil(errGet)
    c.Equal(int64(23), getOrderResponse.UserId)
    orderItem := getOrderResponse.OrderItems[0]
    c.Equal(float32(1.23), orderItem.UnitPrice)
    c.Equal(int32(3), orderItem.Quantity)
    c.Equal("CAM123", orderItem.ProductCode)
}

func (c *CreateOrderTestSuite) TearDownSuite() {
    execError := c.compose.WithCommand([]string{"down"}).Invoke()   // stacki söndür
    if execError.Error != nil { log.Fatalf(...) }
}

func TestCreateOrderTestSuite(t *testing.T) { suite.Run(t, new(CreateOrderTestSuite)) }
```
**İcra:** `go test -run "^TestCreateOrderTestSuite$"` — regex ilə spesifik suite.

### 7. Test coverage + CI
```bash
go test -cover ./...                                # paket üzrə %
go test -coverprofile=coverage.out                  # profil faylı
go tool cover -html=coverage.out                     # HTML rapor
```
**CI axını (PR):** test + coverage + static analiz + vulnerability → threshold altına düşsə PR FAIL → merge-ə icazə YOX; uğurlu → Docker image build (tag from Git).

## Əsas terminlər
- Test piramidi — unit (böyük) / integration / E2E (kiçik)
- SUT — System Under Test
- testify/mock — On/Return/Called
- mockery --all — avtomatik mock generasiyası
- suite.Suite + Setup/TearDown (Suite/Test səviyyələri)
- Testcontainers — Docker-da test DB
- wait.ForSQL — hazırlıq zənbili
- Docker Compose — build/environment/depends_on/ports/healthcheck/volumes
- Multistage build — builder + scratch
- LocalDockerCompose — compose up/down proqrammatik
- coverage.out / -html — coverage alətləri
- CI threshold — minimum coverage qaydası

## Praktik nəticə
1. Unit testləri mock On/Return ilə; xəta yollarını da (DB, payment) ayrıca test et.
2. Interfeys çoxdursa mockery ilə generasiya — əl ilə mock qorumaq texniki borc.
3. Integration test DB-yə QARŞI: Testcontainers + MySQL + wait.ForSQL; suite SetupSuite-də spin, TearDown-da down.
4. E2E: docker-compose stack + uuid identifikator (paralel CI-də toqquşmasın) + gRPC client ilə tam axın assert.
5. Coverage profile-u CI-ə ötür — threshold qaydası ilə PR bloklanması; coverage = inam səviyyəsi.

## Mənbə
Pages: 103-124 (PDF səh. 129-150)
