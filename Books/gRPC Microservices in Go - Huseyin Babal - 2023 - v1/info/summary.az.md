# gRPC Microservices in Go — Xülasə (Azərbaycanca)

> **Kitab:** gRPC Microservices in Go — Hüseyin Babal, Manning, 2023 (ISBN 978-1-63343-920-7) · 200 səh.
> 🚀 Advanced (4/5) · Praktik e-commerce layihəsi: Order/Payment/Shipping servisləri Go+gRPC+K8s ilə qurulur.

---

## Kitabın yanaşması
3 hissəli strukturl: (1) arxitektura nəzəriyyəsi — monolith vs mikro servis, saga, gRPC giriş; (2) tam praktiki inkişaf dövrü — proto→hexagonal→kommunikasiya→resiliency→test→K8s deploy; (3) observability. Kitab boyu tək e-commerce layihəsi böyüyür; kod nümunələri GitHub Actions-da icra olunur.

## Fəsl-fəsil xülasə

### Ch 1 — Giriş: gRPC nədir, nə vaxt
gRPC-nin 5 faydası: performans (Protobuf+HTTP/2), kod generasiyası (çoxdilli stub-lar), fault tolerance (idempotency əsaslı), security (TLS/ALTS), streaming (unary/server/client/bidi). REST vs RPC cədvəli: REST = brauzer + sərbəst dəyişiklik; gRPC = low-latency + polyglot + SDK. Hibrid həll: gRPC-web/Twirp. Production e-commerce arxitekturası: 5 servis (Product/Cart/Checkout/Payment/Shipping) + K8s + CI/CD + observability + API gateway.

### Ch 2 — Mikro servis arxitekturası
Monolith pro/con (deploy/scaling problemi); **scale cube** (X nüsxə / Z data-partition / Y funksional = mikro servis); data consistency: monolith TX → paylanmış: 2PC YOX, **Saga** — choreography (event zənciri) vs orchestrator (dirijor + compensation transaction yuxarıdan-aşağı); service discovery (client/server-side, K8s built-in); ilk proto nümunəsi + protoc + grpc.Dial.

### Ch 3 — Protobuf dərinliyi
Mesaj strukturu (rules/types/numbers); field № 1-15 = 1 bayt (tez sahələrə!); reserved qaydası; **binary encoding** (varint, MSB, 7-bit blok, little-endian); protoc flagləri; ayrı proto repo (dil qovluqları + git tag konvensiyası) + GitHub Actions matrix avtomatizasiyası (tag push → generate → go mod init → yeni tag); backward/forward compatibility ssenariləri (yeni sahə, server/client upgrade asimmetriyası, oneof tələləri).

### Ch 4 — Hexagonal Order servisi
Ports (APIPort/DBPort interfeysləri) + adapters (GORM DB, gRPC server) + application core + domain; DI main-da (db→app→grpc); 12-factor env config (fail-fast); GORM AutoMigrate + 1:N relation; gRPC Run (Listen/Register/serve + reflection yalnız dev-də); grpcurl ilə Create testi.

### Ch 5 — İnterservis kommunikasiyası
Server/client LB müqayisəsi (K8s = server-side); PaymentPort+Adapter (stub → Charge → Create RPC); PAYMENT_SERVICE_URL env; core-da PlaceOrder → Save + Charge; **gRPC status kodları**; status.Errorf (xam err = Unknown!); errdetails.BadRequest_FieldViolation ilə zəngin xətalar; client tərəfdə Convert + Details emalı; iki servisin lokal işə salınması.

### Ch 6 — Resiliency
Cascading failure problemi; **timeout** (context.WithTimeout + proqnozlaşdırma); **retry** (grpc_retry interceptor: WithCodes(Unavailable/ResourceExhausted), WithMax, Backoff Linear/Exponential/Jitter); **circuit breaker** (gobreaker: Closed→Open→Half-Open; Settings: MaxRequests/Interval/Timeout/ReadyToTrip/OnStateChange) — CB interceptor kimi mərkəzləşdirilir; gRPC error modeli (Code/Status/Message/ErrorDetails) + validasiya nümunəsi; TLS/mTLS (handshake, openssl sertifikat generasiyası, credentials.NewTLS server+client konfiqurasiyası).

### Ch 7 — Testing
Test piramidi; SUT ssenari konsepti; testify mock (On/Return) + mockery avtomatik generasiya; unit testlər (happy path, DB xəta, payment xəta status-assert ilə); **integration**: testify suite + Testcontainers MySQL (wait.ForSQL, SetupSuite/TearDownSuite); **E2E**: docker-compose stack (mysql+payment+order) + LocalDockerCompose up/down + gRPC client tam axın assert; coverage + CI threshold.

### Ch 8 — Deployment
Multistage Dockerfile (builder + scratch, CGO_ENABLED=0); K8s resursları (Pod/Deployment/Service); LB strategiyaları: LoadBalancer = servis başına LB (bahalı!) → **NGINX Ingress controller** = 1 LB + Ingress resursları (backend-protocol: GRPC); **cert-manager** (CRD-lər, ClusterIssuer self-signed, Let's Encrypt/Vault inteqrasiya, Ingress annotation); minikube tunnel + grpcurl TLS testi; deployment strategiyaları: RollingUpdate (default; backward-compat tələbi!), Blue-Green (LB switch + dublikat), Canary (eyni selector, tədricən); GitOps (ArgoCD/FluxCD) qeydi.

### Ch 9 — Observability
Traces (trace ID + span zənciri) / metrics (SLA/SLI/SLO; average tələsi → **percentile p95**) / logs (node vs cluster-level); **OpenTelemetry** interceptor-ları (otelgrpc) + GORM instrumentasiyası; metrics backend: Jaeger All in One + OTel Collector + Prometheus (Helm); Jaeger SPM (p95/rate/error); trace-inyeksiyalı loglar (logrus custom formatter: SpanFromContext → trace_id/span_id JSON-da); log stack: Fluent Bit DaemonSet → Elasticsearch (ECK operator) → Kibana (TraceID filtri ilə tam axın).

## Kitabın əsas mesajları
1. **gRPC = mikro servis kommunikasiyasının standartı:** stub generasiyası, HTTP/2, strukturlaşdırılmış xətalar, built-in resiliency qatları.
2. **Hexagonal + portlar:** hər asılılıq interfeysə — bu, mock-test və adapter dəyişməzliyinin əsasıdır.
3. **Resiliency patterni qat-qat:** timeout → retry (kod-seçici) → circuit breaker (fail-ratio) — hamısı interceptor-larla mərkəzləşdirilir.
4. **Compatibility protobuf-un qanunudur:** field number qorunması, oneof ehtiyatı, semver — yoxsa RollingUpdate zamanı v1/v2 qarışığı data itirir.
5. **K8s standart dəstəklə gəlir:** service discovery, LB, sertifikat (cert-manager), deployment strategiyaları.
6. **Observability üçlüyü korrelyasiya tələb edir:** trace_id hər metriyin və logun içində — anomal tap → trace bax → log filtrlə.

## Kitabdan sonra
- gRPC streaming dərinləşməsi (kitab əsasən unary ilə işləyir)
- Service mesh (Istio/Linkerd) — resiliencyni application-dan alır
- Distributed tracing standartları (OpenTelemetry Collector pipelines)
- Event-driven arxitektura (Kafka) — choreography saga-nın real implementasiyası
- ArqoCD/FluxCD ilə tam GitOps workflow
