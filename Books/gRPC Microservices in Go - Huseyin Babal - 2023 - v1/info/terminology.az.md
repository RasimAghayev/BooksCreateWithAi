# gRPC Microservices in Go — Terminologiya (Azərbaycanca)

> Texniki terminlər `English (Azərbaycanca qarşılığı)` formatında — kitabın ardıcıllığı ilə.

## Arxitektura (Ch 1-2)
- **gRPC (gRPC Remote Procedure Call)** — Google-un RPC frameworkü (2015)
- **Protobuf (protocol buffers)** — dil/platforma-müstəqil binary serialization
- **IDL (Interface Definition Language)** — interfeys təyinat dili
- **Stub** — client tərəfli generasiya olunan service-method proxy
- **Idempotency (eynilik)** — eyni əməliyyatın təkrarı eyni nəticə
- **Polyglot** — çoxdilli mühit
- **Monolithic architecture** — vahid paket
- **Scale cube** — X (nüsxə) / Y (funksional) / Z (data-partition) oxları
- **Y-axis scaling** — funksional decompose = mikro servislər
- **2PC (two-phase commit)** — paylanmış atomic transaksiya
- **Saga** — lokal transaksiyalar zənciri
- **Choreography saga** — event-idarəli zəncir
- **Orchestrator saga** — mərkəzi dirijor + compensation transaction
- **Command channel / Pub-sub** — növbəti servsə birbaşa / hadisə yayımı
- **replyTo / correlation ID** — cavab kanalı / zəncir izi
- **Service discovery** — client-side / server-side / K8s built-in
- **gRPC-Web** — brauzer proxy qatı
- **Twirp** — Protobuf üzərində REST framework

## Protobuf (Ch 3)
- **proto3** — sadələşdirilmiş syntax versiyası
- **Field rule** — singular / repeated
- **Field number** — binary identifikator (1-15 = 1 bayt)
- **reserved** — silinmiş sahə nömrə/adının qorunması
- **Varint** — int wire type (0)
- **MSB (Most Significant Bit)** — davam biti
- **Little-endian** — tərs bayt sırası
- **protoc / protoc-gen-go / protoc-gen-go-grpc** — compiler + Go pluginlər
- **go_package** — generasiya modulunun yolu
- **oneof** — ya-o-ya sahə qrupu
- **Backward / Forward compatibility**
- **Matrix strategy** — GitHub Actions çox-job dəyişəni
- **semver** — breaking change = major bump

## Hexagonal (Ch 4)
- **Hexagonal architecture (ports & adapters)** — altıbucaqlı arxitektura
- **Port** — qatlararası müqavilə interfeysi
- **Adapter** — portun konkret implementasiyası
- **Driver / Driven side** — giriş / çıxış tərəfləri
- **Domain model** — biznes obyekti
- **GORM** — Go ORM (gorm.Model, AutoMigrate)
- **12-factor app** — env var konfiqurasiya
- **Fail-fast** — çatmayan config = dərhal dayanma
- **UnimplementedXServer** — forward compatibility baza
- **grpcurl** — gRPC curl; **reflection** — dev-də introspection

## Kommunikasiya (Ch 5)
- **Server-side / Client-side load balancing**
- **Round-robin** — növbəli paylama
- **insecure.NewCredentials** — TLS-siz dial (dev)
- **gRPC status kodları** — OK/CANCELLED/INVALID_ARGUMENT/DEADLINE_EXCEEDED/NOT_FOUND/ALREADY_EXISTS/PERMISSION_DENIED/RESOURCE_EXHAUSTED/INTERNAL
- **status.Errorf / New / FromError / Convert**
- **errdetails.BadRequest_FieldViolation** — sahə-özəl xəta detalı
- **WithDetails / Details()** — status zənginləşdirmə/açma

## Resiliency (Ch 6)
- **Cascading failure** — zəncirvari çökmə
- **context.WithTimeout / WithDeadline** — deadline + cancel funksiyası
- **Deadline propagation** — context-in zəncirdə yayılması
- **Transient fault** — keçici xəta
- **grpc_retry interceptor** — WithCodes/WithMax/WithBackoff
- **Backoff: Linear / Exponential / Jitter**
- **Circuit breaker: Closed / Open / Half-Open**
- **gobreaker.Settings** — MaxRequests/Interval/Timeout/ReadyToTrip/OnStateChange
- **UnaryClientInterceptor / UnaryServerInterceptor**
- **TLS handshake / mTLS (mutual TLS)** — biryönlü/ikiyönlü sert təsdiqi
- **CA / CSR / self-signed / subjectAltName**
- **credentials.NewTLS** — TransportCredentials

## Testing (Ch 7)
- **Test piramidi** — unit/integration/E2E nisbəti
- **SUT (System Under Test)**
- **testify/mock** — On/Return/Called
- **mockery** — avtomatik mock generator
- **suite.Suite + SetupSuite/SetupTest/TearDownTest/TearDownSuite**
- **Testcontainers** — Docker-da test DB
- **wait.ForSQL** — DB hazırlıq zənbili
- **Docker Compose** — build/depends_on/healthcheck/volumes
- **Multistage build** — builder + scratch
- **LocalDockerCompose** — compose up/down proqrammatik
- **coverage.out / -coverprofile / -html**
- **CI threshold** — minimum coverage PR qaydası

## Deployment (Ch 8)
- **Pod / Deployment / ReplicaSet / Service / Ingress**
- **ClusterIP / NodePort / LoadBalancer**
- **NGINX Ingress controller** — vahid LB + routing
- **backend-protocol: GRPC** — HTTP/2 annotation
- **cert-manager / ClusterIssuer / CRD** — sertifikat avtomatlaşdırması
- **minikube tunnel** — lokal ingress proxy
- **RollingUpdate / Blue-Green / Canary**
- **GitOps** — ArgoCD / FluxCD
- **ExternalSecrets** — xarici gizli data

## Observability (Ch 9)
- **Trace / Span** — sorğu səyahəti / əməliyyat vahidi
- **rpc.service / rpc.method / span.kind / rpc.system** — span taqları
- **SLA / SLI / SLO**
- **Percentile (p95)** — sıralanmış faiz nöqtəsi
- **OpenTelemetry (OTel)** — telemetriya SDK/API
- **otelgrpc.UnaryClientInterceptor** — gRPC instrumentasiyası
- **OTel Collector** — telemetriya qəbul nöqtəsi
- **Jaeger / Jaeger SPM** — trace UI + performans monitorinqi
- **Prometheus** — time-series metric anbarı
- **trace.SpanFromContext** — context-dən span
- **logrus JSONFormatter / FieldMap** — struktur log
- **DaemonSet** — node başına 1 pod
- **Fluent Bit** — log collector/shipper
- **Elasticsearch (ECK) / Kibana** — log backend + dashboard
- **Logstash_Format / Retry_Limit** — Fluent Bit output
