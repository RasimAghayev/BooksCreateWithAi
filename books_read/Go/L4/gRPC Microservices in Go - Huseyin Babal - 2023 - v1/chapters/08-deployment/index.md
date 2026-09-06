# Chapter 8 — Deployment (Deployment)

## Bu chapter nədən bəhs edir?
Docker image qurulumuna, Kubernetes arxitekturasına (Pod, Deployment, Service, Ingress), cert-manager ilə sertifikat avtomatlaşdırmasına və deployment strategiyalarına (RollingUpdate, Blue-Green, Canary).

## Əsas fikirlər

### 1. Docker — cloud-native artifact
**Multistage Dockerfile** (Ch 7-də tanış):
```dockerfile
FROM golang:1.18 AS builder
WORKDIR /usr/src/app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o order ./cmd/main.go
FROM scratch
COPY --from=builder /usr/src/app/order ./order
CMD ["./order"]
```
- `builder` alias → scratch (boş) runtime — minimal image
- `CGO_ENABLED=0` — statik binary (scratch-da libc yoxdur!)

### 2. Kubernetes resursları
| Resurs | Vəzifə |
|---|---|
| **Pod** | ən kiçik iş vahidi — konteyner qrupu |
| **Deployment** | replica idarəsi + image update (ReplicaSet-lər) |
| **Service** | sabit DNS adı + LB (ClusterIP/NodePort/LoadBalancer) |
| **Ingress** | HTTP(S) routing qaydaları |

**Service tipləri:**
- **ClusterIP** — daxili (default)
- **LoadBalancer** — hər servis üçün ayrı cloud LB → mikroservisdə BAHALI (N servis = N LB!)
- **NodePort** — worker node portunu açır (IP bilmək lazımdır)

### 3. NGINX Ingress controller — vahid giriş nöqtəsi
**Həll:** 1 LB (controller) + Ingress resursları routing edir → N servisə 1 LB.

```bash
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm repo update
helm install nginx-ingress ingress-nginx/ingress-nginx
```
**Ingress resursu (Order servisi, gRPC + TLS):**
```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  annotations:
    kubernetes.io/ingress.class: nginx
    nginx.ingress.kubernetes.io/backend-protocol: GRPC     # gRPC backend
    nginx.ingress.kubernetes.io/ssl-redirect: "true"
    cert-manager.io/cluster-issuer: selfsigned-issuer       # avto-sertifikat
  name: order
spec:
  rules:
    - http:
        paths:
          - path: /Order
            pathType: Prefix
            backend:
              service:
                name: order
                port: { number: 8080 }
  tls:
    - hosts: [ ingress.local ]
```
**Mexanizm:** Ingress dəyişəndə controller nginx.conf-u yenidən render edir.

### 4. Certificate management — cert-manager
**Qurulum:**
```bash
helm repo add jetstack https://charts.jetstack.io
helm repo update
helm install cert-manager jetstack/cert-manager \
  --namespace cert-manager --create-namespace \
  --version v1.10.0 --set installCRDs=true
```
**CRD-lər:** CertificateRequests, Certificates, Challenges, **ClusterIssuers**, Issuers, Orders.

**ClusterIssuer (self-signed, lokal dev):**
```yaml
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: selfsigned-issuer
spec:
  selfSigned: {}
```
```bash
kubectl apply -f cluster-issuer.yaml
kubectl get clusterissuers -o wide selfsigned-issuer
```
**İnteqrasiyalar:** Let's Encrypt (production), Vault, Cloudflare.

**Client sertifikatı:** `minikube tunnel` (127.0.0.1 → ingress proxy) + browser-dan sertifikat export → sistem keychain.
**grpcurl TLS testi:**
```bash
grpcurl -import-path /path/to/order -proto order.proto ingress.local:443 Order.Create
```
- `-import-path`/`-proto` — reflection olmadan metod kəşfi; proqrammatik client (Go) üçün proto dependency kifayətdir.

### 5. Deployment strategiyaları

#### a) RollingUpdate (default)
**Mexanizm:** yeni ReplicaSet yarat → köhnəsini tədricən azalt.
```bash
kubectl set image deployment/order order=order:1.1.0
```
**Fazalar:** v1 (2 pod) → v2-yə 1 pod əlavə → v1-dən 1 sil → ... → v2 tam.
**Tələb:** BACKWARD COMPATIBILITY — keçid müddətində v1+v2 paralel işləyir!

#### b) Blue-Green
**Mexanizm:** v2-ni tam hazır saxla (blue) → test et → DNS/LB record-u switch → v1 deprecate.
- İkinci Ingress controller/LB blue üçün
- Blue daxili şəbəkədən test olunur
- **Artı:** sıfır-downtime switch; **Eksi:** tam dublikat sistemi = xərc

#### c) Canary
**Mexanizm:** eyni SELEKTORLU ikinci Deployment → Service hər ikisinə paylayır → tədricən v2 çəkir, v1-i zero-ya endirir.
- Eksperimental feature-u user subset-inə göstərmək
- Tələb: `app: order` selektoru HƏR İKİ deploymentdə eyni

**Seçim rəhbərliyi:**
| Strategiya | Nə vaxt |
|---|---|
| RollingUpdate | standart, backward-compatible dəyişikliklər |
| Blue-Green | sıfır-risk switch, əlavə infra ödənilə bilərsə |
| Canary | feedback tələb edən eksperimental çıxışlar |

### 6. GitOps qeydi
`kubectl` əvəzinə: **ArgoCD**, **FluxCD** — repo push → avtomatik deploy. ExternalSecrets — gizli data üçün.

## Əsas terminlər
- Multistage build — builder + scratch
- CGO_ENABLED=0 — statik kompilyasiya
- Pod / Deployment / ReplicaSet
- ClusterIP / NodePort / LoadBalancer
- Ingress / Ingress controller / annotations
- cert-manager / ClusterIssuer / CRD
- minikube tunnel
- RollingUpdate / Blue-Green / Canary
- Backward compatibility — paralel versiya tələbi
- GitOps — ArgoCD / FluxCD

## Praktik nəticə
1. Public ekspozisiyada hər servisə LoadBalancer YOX — 1 Ingress controller + path-əsaslı routing.
2. gRPC backend üçün `backend-protocol: GRPC` annotation-i mütləq (NGINX HTTP/2 biləndir).
3. Sertifikatları əllə YOX — cert-manager + ClusterIssuer (dev: self-signed, prod: Let's Encrypt).
4. Default strategiya RollingUpdate — amma protobuf mesajları dəyişəndə backward compatibility yoxdursa, bu strategiya təhlükəlidir (v1-v2 qarışıq trafik!).
5. Canary üçün eyni selector şərti; user-facing eksperimentlər üçün idealdır.
6. Production-de GitOps alətinə keç — kubectl əl əməliyyatı tez texniki borc olur.

## Mənbə
Pages: 127-148 (PDF səh. 155-176)
