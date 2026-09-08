# Chapter 8 — Deployment with Kubernetes (səh. 147-162)

## Bu fəsil nədən bəhs edir?

Deployment prosesi, mühit konfiqurasiyası, Docker konteynerləri, Kubernetes
deployment modeli və best practices (rollback, canary, CD).

## Əsas fikirlər

### 1. Deployment əsasları
- **Build** (compile + fayllar) → **Rollout** (uzun serverlərə kopyalama)
- Rollout ARDICILDIR: bir instans yenilənir → health yoxlanır → növbəti —
  yeni versiya bug-ludursa bütün serverlər bir anda düşmür
- **Mühitlər (Environments):**
  - Local/dev — developer maşını, sadələşdirilmiş komponentlər
  - Staging — prod-un güzgüsü, amma ayrı data (prod-a toxunmur)
  - Production — real istifadəçi request-ləri

### 2. Konfiqurasiyanın ayrılması
- Hardcode vs **separate config** — YAML faylları hər mühit üçün ayrıca
- Invalid konfiq dəyişikliyi = prod outage-lərin ƏSAS səbəblərindən —
  pre-commit/receive hooklarla YAML validasiyası tövsiyə olunur
```go
type serviceConfig struct {
    APIConfig apiConfig `yaml:"api"`
}
type apiConfig struct { Port string `yaml:"port"` }
// os.Open("base.yaml") → yaml.NewDecoder(f).Decode(&cfg)
```
- `gopkg.in/yaml.v3` paketi; metadata→8081, rating→8082, movie→8083

### 3. Kubernetes modeli
- Google mənşəli, Linux Foundation dəstəkli orkestrasiya platforması
- **Pod** — ən kiçik deploy vahidi (1+ konteyner) → **Node** (host) →
  **Cluster** (node qrupu)
- Kubernetes öz üstünləri: built-in **service discovery**, **rollback**,
  **automated restart** (crash olmuş pod yenidən qalxır), autoscaling

### 4. Deploy addımları
1. **Konteyner image:** Dockerfile (alpine:latest — bir neçə MB):
```dockerfile
FROM alpine:latest
COPY main .
COPY configs/. .
EXPOSE 8081
CMD ["/main"]
```
2. **Cross-compile:** `GOOS=linux go build -o main cmd/*.go` (Linux üçün!)
3. `docker build -t metadata .` → `docker run -p 8081:8081`
4. Docker Hub-a push: `docker tag`/`docker push <user>/metadata:1.0.0`
5. **Deployment YAML:**
```yaml
apiVersion: apps/v1
kind: Deployment
metadata: {name: metadata}
spec:
  replicas: 2            # pod sayı
  selector: {matchLabels: {app: metadata}}
  template:
    metadata: {labels: {app: metadata}}
    spec:
      containers:
      - name: metadata
        image: <user>/metadata:1.0.0
        imagePullPolicy: IfNotPresent
        ports: [{containerPort: 8081}]
```
6. `minikube start` → `kubectl apply -f kubernetes-deployment.yml`
7. Yoxlama: `kubectl get deployments`, `kubectl get pods` (Running, 2 ədəd)
8. `kubectl logs -f <POD_ID>`; `kubectl port-forward <POD_ID> 8081:8081`

### 5. Deployment best practices
**Automated rollbacks (K8s-də default YOXDUR — özün qur):**
- davamlı health check → problem + son vaxt deployment üst-üstə düşürsə →
  `kubectl rollout undo deployment <NAME>`
- Sıra: canary config → yoxla → prod config

**Canary (Kanari) deployments:**
- Yeni versiya traffic-in 1-3%-i üzərində sınaq (50 pod → 1 canary + 49 prod)
- İki ayrı deployment config: `rating-canary`, `rating-production`
- Kiçik fraqmentdə test → bug təsiri minimal

**Continuous Deployment (CD):**
- Hər kod dəyişikliyində avtomatik deploy (Git hook / commit izləmə)
- Fayda: deployment xətası ERKƏN aşkarlanır
- Yüksək cadence → avtomatik health check tooling tələb edir (Ch11-12)

## Termindirmə (AZ)
- Deployment — Deploy (yayım)
- Rollout — Mərhələli Yayım
- Environment — Mühit
- Canary Deployment — Kanari Yayımı (kiçik hissədə sınaq)
- Continuous Deployment (CD) — Fasiləsiz Yayım
- Pod — Kubernetes-in ən kiçik deploy vahidi
- Node/Cluster — Düyün/Klaster

## Kviz sualları
1. Rollout niyə ardıcıldır, paralel deyil? (Yeni versiya bug-ludursa bütün
   instanslar eyni anda düşməsin)
2. Canary deployment nədir? (Yeni versiyanı prod traffic-in 1-3%-ində sınaq)
3. `GOOS=linux go build` nə üçün lazımdır? (Konteyner Linux əsaslıdır —
   binary Linux üçün compile olunmalıdır)
4. Automated rollback K8s-də varmı? (Yox — health check + deployment vaxtı
   müqayisəsi + rollout undo ilə özün qurulur)
