# Chapter 11 — Deploying Applications to the Cloud (səh. 316-343)

## Bu chapter nədən bəhs edir?

MallBots modular monolith-inin mikroservislərə parçalanması: Docker Compose
profile-ləri, paylaşılan system paketi, multi-stage Dockerfile, Nginx reverse
proxy, gRPC service discovery (konfiqurasiya ilə), Terraform ilə AWS environment
(EKS, RDS, ECS) və Kubernetes-ə deploy.

## Əsas fikirlər

### 1. Monolith → Mikroservis parçalanması
**Strategiya:** Kodu kopyalamaq YOX — mövcud composition root-u yenidən istifadə
etmək:
1. `app` struct → **`system` paketi** ( yeni `System` strukturu); köhnə
   `Monolith` interfeysi `Service` adlanır
2. `NewSystem(cfg)` konstruktoru: initDB, initJS, initMux, initRpc — monolith
   və hər mikroservis EYNİ initialization-i paylaşır
3. Hər modula üçün `cmd/<module>/service` main paketi: `module.Root(ctx, s)`
   çağırışı + Waiter (WaitForWeb, WaitForRPC, WaitForStream)
4. **Niyə hər ikisi saxlanılır?** Real dünyada monolith rejimini saxlamaq
   seçimi lokal dev üçün asanlıq verir — single process debug

### 2. Docker Compose profile-ləri
```yaml
services:
  mallbots-monolith:
    profiles: [monolith]
  baskets:
    profiles: [microservices]
    expose: ['9000']        # daxili gRPC
    environment:
      PG_CONN: ...
      NATS_URL: nats:4222
    depends_on: [nats, postgres]
    command: ["./wait-for", "postgres:5432", "--", "/mallbots/service"]
```
```bash
docker compose --profile monolith up      # NATS+PG+Pact+monolith
docker compose --profile microservices up # 9 konteyner
```

### 3. Multi-stage Dockerfile (tək, bütün servislər üçün)
```dockerfile
FROM golang:1.18 AS builder
ARG svc                      # build argument — hansı servis
COPY .. ./
RUN go build -ldflags="-s -w" -v -o service ./${svc}/cmd/service

FROM alpine:3 AS runtime
COPY --from=builder /mallbots/docker/wait-for .
COPY --from=builder /mallbots/service /mallbots/service
CMD ["/mallbots/service"]
```
- `-ldflags="-s -w"` — debug simvolları atılır (binary kiçik)
- `wait-for` utility — DB hazır olana qədər gözlə

### 4. Nginx reverse proxy
**Problem:** Swagger UI hər mikroservisdə öz spec-ini gözləyir — UI ayrı-ayrı
portlarda "sınıq" təcrübə verir.
**Həll (nginx.conf):**
```nginx
listen 8080;
location /api/baskets {
    proxy_pass http://docker-baskets;
}
location /baskets-spec/ {
    proxy_pass http://docker-baskets;
}
# hər servis üçün cüt location blokları + swagger-ui üçün son location /
```
- Nginx konteyneri `microservices` profile-də; config mount edilir
- İstənilən URL hər hansı servisə yönəldilir (proxy_redirect off)

### 5. gRPC service discovery ( sadə həll)
Mikroservislərdə gRPC client-lər hardan dial edəcəklərini bilməli. **Konfiq
həlli:** servis cütlükləri env dəyişənləri ilə:
```
RPC_SERVICES=ORDERING=ordering:9000,STORES=stores:9000
```
`Services` tipi üçün custom `Decode()` — `KEY=VALUE` cütlüklərini map-ə
çevirir; `RpcConfig.Service(name)` convenience metodu.

### 6. DevOps alətləri
| Alət | Məqsəd |
|---|---|
| Terraform | IaC — AWS resursları |
| AWS CLI | AWS idarəsi; IAM user (Programmatic access + AdministratorAccess) |
| eksctl / helm | K8s klaster alətləri |
| kubectl | K8s CLI |
| K9s | K8s TUI (:deployments, :ingress; çıxış :quit) |
| Make | çoxaddımlı əmrlərin qısaltması |

Alətləri ya lokal qur, ya da `deploytools` konteynerində işlət (`deploytools aws configure`).

### 7. Terraform ilə AWS infrastrukturu
**2 Availability Zone, planlaşdırılmış resurslar:**
- **ECS Docker repositories** — mikroservis image-ləri buraya push olunur
- **EKS** Kubernetes klaster (eks.tf; ALB üçün IAM policies/roles)
- **RDS** serverless PostgreSQL (klasterdən əlçatan)
- VPC, subnetlər və s.

**Terraform faylları:**
- `terraform init` — modulları endirir; `terraform validate` — yoxlayır
  (`make ready` hər ikisini edir)
- `make deploy` ( `terraform apply`) — resursları qurur; kəsilsə yenidən run
  etmək çox vaxtı düzəldir
- **State:** plan/apply fərqlərini və destroy üçün resurs yerlərini saxlayır

### 8. Kubernetes deploy (kubernetes.tf)
- **ConfigMap** — ümumi env dəyişənləri (ENVIRONMENT, WEB_PORT, NATS_URL)
  bütün deployment-lərə ötürülür
- Hər mikroservis = Deployment resursu (mallbots namespace)
- ALB Ingress Controller → ingress resursları
- `kubectl get deployment -n mallbots` / K9s ilə izləmə

### 9. Təmizləmə (destroy)
```
/deployment/application:   make destroy
/deployment/infrastructure: make destroy
```
Sıra vacibdir — əvvəl application, sonra infrastructure; state məlumatına
əsasən resurslar tapılıb silinir.

## Əsas terminlər

- Microservices vs Modular Monolith
- Composition Root reuse (təkrar istifadə)
- Docker Compose Profiles
- Multi-stage Dockerfile / ARG
- Reverse Proxy (Nginx)
- Service Discovery (env əsaslı)
- Infrastructure as Code (IaC)
- Terraform init/validate/apply/destroy
- EKS / ECS / RDS / ALB
- ConfigMap / Ingress
- K9s (Kubernetes TUI)

## Praktik nətivə

- Monolith-i parçalayan kimi kodu DUBLICƏTLƏŞDİRMƏ — system paketi hər ikisinə
  xidmət edir; monolith rejimini saxlamaq dev üçün üstünlükdür
- Tək parametrli (ARG svc) multi-stage Dockerfile bütün servisləri qurur
- Reverse proxy + spec route-ları UX-i bərpa edir
- Service discovery-ə başlamaq üçün sadə env config kifayətdir — dedike
  system lazım deyil
- Terraform state həm plan, həm destroy üçün addım-addım yazılır; destroy
  sırası: application → infrastructure

## Mənbə

Pages: 316-343 (Chapter 11, Event-Driven Architecture in Golang)
