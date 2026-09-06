# Chapter 1 — Introduction to Go gRPC microservices (Go gRPC mikroservislərinə giriş)

## Bu chapter nədən bəhs edir?
gRPC-nin faydalarına (performans, kod generasiyası, fault tolerance, təhlükəsizlik, streaming), REST vs RPC müqayisəsinə, gRPC-nin nə vaxt seçilməli olduğuna və kitab boyu qurulacaq production-grade e-commerce arxitekturasına.

## Əsas fikirlər

### 1. gRPC nədir
**Nədir:** Google-un 2015-də yaratdığı open source RPC framework — load balancing, tracing, fault tolerance və security daxili dəstəklə; çoxdilli server/client stub generasiyası.

**Mikroservislərlə əlaqə:** Monolithdə "payment çağırmaq" = class method; mikroservisdə = ŞƏBƏKƏ çağırışı (TCP/HTTP/queue) — daha çətin idarə. gRPC bu çağırışların çətinliklərini (network failure, TLS, connection pooling) həll edir.

### 2. gRPC-nin 5 faydası

#### a) Performance (performans)
- **Protobuf (protocol buffers):** dil/platforma-müstəqil serialization — kiçik, kompakt mesajlar, sürətli encode/decode
- **HTTP/2:** server push, multiplexing, header compression

#### b) Code generation və interoperabilitet
- IDL (Interface Definition Language) ilə mesajları təyin et → istənilən dil üçün stub generasiya et
- Çoxdilli dəstək: Go, Java, Python, Ruby, JavaScript, C#...
- **Sitat:** "A little duplication is far cheaper than wrong abstraction" — birgə model kitabxanası əvəzinə, hər servisə generasiya olunan stub
- Server dili ≠ client dili — Python client Java server-ə asanlıqla bağlanır

#### c) Fault tolerance (xətalara davamlılıq)
- **Idempotency (eynilik):** eyni əməliyyatın təkrarı eyni nəticə — retry üçün açar şərt
- Idempotent olmayan əməliyyatlarda proper validation error-lar → retry-ın nə vaxt dayanacağını bil
- Rate limiting, circuit breaker, fault injection (Ch 6-da dərinləşəcək)

#### d) Security
- HTTP/2 over SSL/TLS standart
- Alternativlər: ALTS (Application Layer Transport Security), token-based auth

#### e) Streaming
- Pagination əvəzinə axın: server-side, client-side, və ya **bidirectional streaming**
- Bağlantı bir dəfə açılır → data axır

### 3. REST vs RPC müqayisəsi
| Xüsusiyyət | REST | gRPC |
|---|---|---|
| Protokol | HTTP/1.0 (2.0 mümkün) | HTTP/2 built-in |
| Format | JSON/XML (insan-oxunarlı) | Protobuf binary |
| Stub generasiyası | Swagger Codegen kimi framework lazım | built-in |
| Streaming | yox (klasik) | unary + bidirectional |
| Browser dəstəyi | TAM | minimal → gRPC-Web proxy lazım |
| Dəyişiklik çevikliyi | sahələr sərbəst dəyişilə bilər | qaydalara uyğun evolyusiya tələb edir |

### 4. Nə vaxt gRPC (və nə vaxt YOX)
**gRPC seç:**
- Low latency tələbi, polyglot (çoxdilli) mühit
- Interservice (daxili) kommunikasiya
- Çoxlu SDK saxlamaq lazım olduqda

**REST seç:**
- Browser dəstəyi kritikdırsa (gRPC üçün çevrici qat lazım olar)
- Sadə layihələr (1-2 servis — proto fayl saxınmaq əziyyətdir)
- İstehlakçıya açıq API SDK-sız (gRPC interfeysini müştəriyə açmaq çətin)

**Hibrid həll:** daxildə gRPC + public-da REST: gRPC load balancer (Envoy) və ya **Twirp** (Protobuf üzərində REST qatı):
```bash
curl -X POST -H "Content-Type: application/json" \
  -d '{"name": "dev-cluster"}' \
  http://localhost:8080/twirp/github.com/huseyinbabal/microservices-proto/cluster/Create
```

### 5. Production-grade e-commerce arxitekturası (kitabın layihəsi)
**5 biznes-kabiliyyət servisi:**
| Servis | Vəzifə |
|---|---|
| **Product** | məhsul axtarış/siyahı/detallar |
| **Cart** | virtual səbətə məhsul əlavə |
| **Checkout** | Cart-dan item-lər, Payment-lə ödə, Shipping-lə göndər |
| **Payment** | ödəniş qapısı (Stripe inteqrasiyası) |
| **Shipping** | müştəri ünvanına çatdırılma (shipping şirkəti) |

**Kommunikasiya:** Hər servisin gRPC stub-u avto-generasiya olunur; Checkout Cart/Shipping/Payment stublarını Go dependency kimi əlavə edib çağırır.

**İnfrastruktur komponentləri:**
1. **Kubernetes (container runtime):** servis kəşfiyyatı = servis adı (domain adı kimi); horizontal scaling servis-özəl (Product daha çox scale etsə, hamısı YOX); resource requests/limits ayrıca
2. **CI/CD:** stub generasiyası push-da avtomatik; unit/integration/contract test-lər; static analiz; vulnerability yoxlaması; UAT → production; deployment strategiyaları: rolling upgrade, canary, blue-green
3. **Observability:** metrics + logs + tracing; trace ID ilə sorğu həyat dövrü izləmə; Prometheus alarmları ("memory > 80%"); Elastic Stack log analizi
4. **Public access:** API gateway — auth/authz, rate limiting; K8s-də NGINX Ingress controller konfiqurasiyası

## Əsas terminlər
- gRPC (gRPC Remote Procedure Call)
- Protobuf (protocol buffers) — binary serialization
- IDL (Interface Definition Language)
- Stub — client tərəfdə generasiya olunan service method proxy-si
- Idempotency (eynilik) — təkrar təhlükəsizliyi
- HTTP/2 multiplexing — paralel axın bir TCP bağlantıda
- ALTS — Google-un application-layer təhlükəsizliyi
- gRPC-Web — brauzer proxy qatı
- Twirp — Protobuf üzərində REST framework
- Polyglot — çoxdilli mühit
- CI/CD — davamlı inteqrasiya/delivery
- UAT — User Acceptance Testing

## Praktik nəticə
1. Daxili servis-kommunikasiyası üçün gRPC; public API üçün REST (və ya gRPC-Web/Twirp proxy).
2. Müştəri SDK-ları gRPC stub-larından generasiya et — əl ilə model dublikatı saxlamadan.
3. Retry məntiqindən əvvəl idempotency təhlil et — eyniliksiz retry data korlanır.
4. Servis parçalanmasını BİZNES-kabiliyyətlərinə görə et (Product/Cart/Checkout/...), texniki qatlara görə YOX.
5. Kubernetes + CI/CD + observability üçlüyünü başdan nəzərdə tut — mikroservis "işləyir" statusu bunlarsız ölçülməz.

## Mənbə
Pages: 3-12 (PDF səh. 20-29)
