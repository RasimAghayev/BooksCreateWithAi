# Chapter 11 — Deploying Applications to the Cloud (səh. 316-343)

## Bu chapter nədən bəhs edir?

Modular monolith-in mikro-servislərə parçalanması (compose profile-lar), nginx
reverse proxy, DevOps alətləri (kubectl, Helm, K9s, Terraform) və AWS-ə
deploy (EKS, RDS, 2 AZ).

## Əsas fikirlər

### 1. Monolith → microservices (eyni kod, 2 rejim)
- `/cmd/mallbots` → bütün modullar bir prosesdə
- `/cmd/service` → **hər modul ayrıca servis** (seçilənlər env dəyişəni ilə)

**Addımlar:**
1. `system` paketi — ümumi interfeyslərin duplikatı (modullar bir-birinə
   gRPC ilə baxır)
2. Composition root-daşınması: modul başlanğıcı `startModule()` funksiyasına →
   istənilən rejimdə istifadə olunur
3. Hər servis üçün ayrı Dockerfile target-ları:
```dockerfile
FROM golang AS builder
RUN go build -ldflags="-s -w" -v -o service ./${svc}/cmd/service
FROM alpine:3 AS runtime
COPY --from=builder /mallbots/service .
```
4. Docker Compose **profile-lar**: `monolith` (tək app) / `microservices`
   (customers, baskets, orders, ... ayrı-ayrı + DB + NATS)

### 2. Reverse proxy (nginx)
Mikro-servislərdə Swagger UI konflikt edir (hamısı eyni portda). Həll:
```nginx
listen 8080;
location /api/baskets   { proxy_pass http://docker-baskets; }
location /baskets-spec/ { proxy_pass http://docker-baskets/; }
```
- Xarici ünvan vahid (8080), daxili yönləndirmə modullara
- Reverse proxy də `microservices` profile-da

### 3. Environment parametrləri
Hər servis öz env dəyişənlərini daşıyır (PG_CONN, NATS_URL...); servis siyahısı
`SERVICES` dəyişəni ilə pars olunur (`key=value` cütləri).

### 4. DevOps alətləri
| Alət | Məqsəd |
|---|---|
| kubectl | K8s CLI |
| **K9s** | TUI — :deployment, :ingress, pod-ları izləmək (çıxış: :quit) |
| **Helm** | K8s paket meneceri (chart-lar) |
| **Terraform** | İnfrastruktur kodu (IaC) |
| eksctl (tərtibat) | EKS klaster yaratma |

Alətlər Docker konteynerə də yerləşdirilə bilər (lokal sistemin təmizliyi) və ya
birbaşa qurulur (AWS user: AdministratorAccess ilə test).

### 5. Terraform + AWS
- **State:** S3-də saxlanılan state faylı — plan/apply/destroy-u mümkün edir
- Fayllar: `vpc.tf`, `eks.tf` (module istifadəsi), `rds.tf`, `kms.tf`...
- **2 Availability Zone** — tipik kiçik mühit
- `terraform init` → `terraform plan` (nə yaradılacağını gör) → `terraform apply`
- Loglar resurs yaratma prosesini izah edir

### 6. Deploy + teardown
- Helm chart-lar vasitəsilə tətbiq EKS-ə (ingress-lər, deployment-lər)
- K9s-də :deployment → pod-ların onlayn gəlişi görünür
- **Teardown:** tətbiqi sil, sonra `terraform destroy` (state sayəsində bütün
  resurslar tapılır) — AWS xərclərini dayandırmaq üçün məcburi addım

## Əsas terminlər

- Compose Profile
- Reverse Proxy (tərs proksi)
- K9s (Kubernetes TUI)
- Helm Chart
- Terraform State / Plan / Apply / Destroy
- EKS / RDS / VPC / AZ
- IaC (Infrastructure as Code)

## Praktik nəticə

- Modular monolith-in gücü: eyni composition root həm tək proses, həm ayrı
  servis-lər kimi işləyir
- Swagger/path konfliktlərini reverse proxy ilə həll et
- İnfrastruktur kodu Terraform-da; state S3-də — destroy üçün vacibdir
- Test hesabı AdministratorAccess ver, amma real layihədə least-privilege

## Mənbə

Pages: 316-343 (Chapter 11, Event-Driven Architecture in Golang)
