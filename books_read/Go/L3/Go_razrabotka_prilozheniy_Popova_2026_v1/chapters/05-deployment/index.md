# Глава 5 — Mikroservislərin yayımlanması (səh. 267-310)

## Bu fəsil nədən bəhs edir?

5 yayım üsulu (server+proseslər → çoxserver+nginx balanslaşdırma →
konteynerlər → orkestrator (Kubernetes) → serverless), nginx quraşdırma və
proxy konfiqurasiyası, CI/CD konsepsiyası (GitLab Runner / GitHub Actions),
GitHub Actions workflow (YAML strukturu: name/on/jobs/runs-on/steps,
workflow_dispatch vs push trigger, Marketplace actions/checkout), gomock +
testify unit-testlər (mock-lar repo/kafka üçün, table-driven testlər,
84.4% coverage), multi-stage Dockerfile (golang:alpine build → alpine:3.18
runtime, layer cache optimallaşdırması), .dockerignore, tam docker-compose
(4 servis + 3 DB + Kafka + kafka-ui, xidmət adları host kimi), Docker Hub
publishing (tag/push), kompose ilə Docker-Compose → Kubernetes miqrasiyası,
Deployment manifesti, kubectl əsasları (run/describe/apply/get/logs), nəticə.

## Əsas fikirlər

### 1. Yayım Üsullarının 5-liyi
| # | Üsul | Xarakteristika |
|---|------|----------------|
| 1 | 1 server, çox proses | lokal dev kimi; production-da instans + balanslaşdırıcı standartdır; SSH ilə köçür + işə sal — miqyaslandırma və dayanıqlıq ÇƏTİN |
| 2 | Çox server, çox proses | nginx/Apache balanslaşdırıcı; konfiqurasiya əl ilə, çox manual iş |
| 3 | Konteynerlər | Docker paketi; orkestrasiyaya ilk addım |
| 4 | Orkestrator | Kubernetes/Nomad — tam platforma |
| 5 | Serverless | buludda birbaşa kod; platforma təlimatı ilə |

- **Yol 1:** proseslər → konteynerlər → Kubernetes; **Yol 2:** serverless.
- 1-ci üsulda ev serveri: public IP + port açmaq; fayl redaktəsi vim/nano
  ilə; çökmə = manual yenidən quruluş.

### 2. Nginx Proxy
```bash
sudo apt install nginx
sudo systemctl status nginx        # active (running) yoxlaması
sudo nano /etc/nginx/sites-available/default
```
```nginx
location / {
    proxy_pass http://localhost:50053;   # 80 → Gateway
}
```
- İkinci üsulda nginx balanslaşdırıcı kimi (çox instans arasında paylama);
  manual iş çox → xidmətlər (services) + CI/CD lazım.

### 3. CI/CD — GitHub Actions
- **Anlayış:** commit = hadisə → avtomatik test/lint/build/deploy;
  kommersiya: GitLab+Runner; alternativ: **GitHub Actions** (pulsuz plan
  kifayətdir)
- **Workflow (YAML):**
```yaml
name: Transaction
on: [workflow_dispatch, push]     # manual + hər push
jobs:
  stage_dev:
    runs-on: ubuntu-latest        # Linux/macOS/Windows
    steps:
    - name: get code
      uses: actions/checkout@v3.5.3   # Marketplace hazır aksiyası
    - name: install deps
      run: go mod download
    - name: Run tests
      run: go test -v -race -coverprofile=coverage.out ./...
```
- `name` (ixtiyari), `on` (hadisə abunəliyi — massiv olar), `jobs` →
  `runs-on` (OS) → `steps` → `uses` (hazır) / `run` (öz komanda)
- Bütün proseslər təmizlənir; commit-dən push-a → test keçibsə kod yaşayır

### 4. Unit-Testlər: gomock + testify
```bash
go get go.uber.org/mock/gomock github.com/stretchr/testify/assert
# Makefile: mockgen -source=internal/service/service.go -destination=.../mocks.go
```
- **Table-driven pattern:** `tests := []struct{ name; setupMocks; expectedError; expectedResult }`
  — uğur + repository xətası + kafka publish xətası ssenariləri
- **Mock gözləntiləri:**
```go
mockRepo.EXPECT().Deposit(gomock.Any(), expectedParams).Return(expectedResult, nil)
mockKafka.EXPECT().Publish(gomock.Any(), "transaction_data", "1", gomock.Any()).Return(nil)
mockRepo.EXPECT().UpdateTransactionStatus(gomock.Any(), uint64(1), model.TransactionStatusFailed)
```
- **Sınaq dövrü:** ctrl := gomock.NewController(t); defer ctrl.Finish();
  mocks.NewMockRepository(ctrl); zerolog.Nop() — logger səssizləşdirmə
- **Nəticə:** coverage 84.4%; xəta halları da mütləq test edilməlidir —
  bütün kod dalları + gözlənilməz cavab reaksiyaları

### 5. Multi-Stage Dockerfile
```dockerfile
FROM golang:1.24.5-alpine AS build   # 1) build mərhələsi
WORKDIR /app
COPY go.mod go.sum ./                # 2) ƏVVƏLCƏ — layer cache üçün
RUN go mod download                 #    asılılıq dəyişməyibsə cache-dən
COPY . .                            # 3) sonra ağır kod
RUN go build -o app ./cmd

FROM alpine:3.18                     # 4) təmiz runtime (~5 MB, kompilyator YOX)
WORKDIR /app
COPY --from=build /app/app
COPY --from=build /app/internal/migrations ./internal/migrations
CMD ["./app"]
```
- **Version pin-lə:** golang:latest YOX — gözlənilməz uyğunsuzluqlar;
  alpine — ölçüsü kiçik
- **CMD vs RUN:** RUN = image yaradılarkən; CMD = konteyner Başlayanda
- **.dockerignore:** (.idea, vendor) — .gitignore kimi; ölçünü azaldır
- **Miqrasiya qovluğu KOPYALANMALI:** servis startda miqrasiya yoxlayır —
  yoxdursa başlamır

### 6. Tam Docker-Compose (bütün sistem)
```yaml
services:
  account_db:   {image: postgres:15-alpine, ports: "5432:5432", ...}
  auth_db:      {ports: "5433:5432"}       # portlar FƏRLİ
  transaction_db: {ports: "5434:5432"}
  kafka:        {confluentinc/cp-kafka:7.5.0, ...}     # + kafka-ui
  account:      {image: account:latest, build: ./account, ports: "50051:8080",
                 environment: {DB_DSN: postgres://...@account_db:5432/...},
                 depends_on: [account_db]}
  auth / transaction / gateway: {eyni pattern}
```
- **KONTEYNER DAXİLİNDE host = SERVİS ADI:** localhost YOX — account_db,
  kafka, account:8080; bütün env-lər compose-da
- **depends_on:** DB hazır olmayan servisi işə salmaz
- docker-compose.yaml — bütün layihələrdə BİR SƏVİYYƏ YUXARI (build yolunu
  sadələşdirmək üçün)

### 7. Docker Hub Publishing
```bash
docker login                                   # hub.docker.com hesabı
docker image ls                                # 4 image: transaction/account/auth/gateway
docker tag book_all-auth yuliapopova/book_all-auth:latest
docker push yuliapopova/book_all-auth:latest
```
- Nəticə: hər yerdən çatılan nəşrlər; compose-da `image:
  yuliapopova/account:latest` — build ƏVİZİNƏ buluddan yükləmə
- VM-ə çatdırılma: GitHub-dan compose + env klon → docker-compose up

### 8. Kubernetes Əsasları
- **Konsepsiyalar:**
  - **Node (düyün)** — fiziki/VM; konteynerlərin işlədiyi yer
  - **Pod** — 1+ konteyner, bir node-da zəmanətli; KLASTER-DAXİLİ unikal IP;
    sabit port istifadə edə bilər (konflikt riski YOX)
  - **Volume** — pod daxili konteynerlər arasında paylaşılan saxlanc
  - **Replication controller** — nüsxə sayı idarəetməsi; çökmüş pod-un
    əvəzi
  - **Service** — podlar dəsti + giriş siyasəti abstraksiyası
  - **kubectl** — API server ilə CLI
- **İmperativ (tövsiyə EDİLMİR):** `kubectl run nginx --image=nginx:latest --port=80`
- **Declarativ manifest:**
```yaml
apiVersion: apps/v1
kind: Deployment
metadata: {name: my-nginx-deployment}
spec:
  replicas: 1
  selector: {matchLabels: {app: nginx}}
  template:
    metadata: {labels: {app: nginx}}
    spec:
      containers:
      - {name: nginx, image: cr.yandex/.../ubuntu-nginx:latest}
```
  - apiVersion (K8s versiyası) / kind (Deployment) / metadata / spec /
    containers
- **Əmrlər:** `kubectl apply -f pod.yaml` (yarat/deyişdir); `kubectl get
  pods` (status siyahısı: Running/ContainerCreating); `kubectl describe pods
  nginx` (tam info + events); `kubectl logs <pod>` (pod logları)

### 9. kompose — Compose → K8s Miqrasiyası
```bash
brew install kompose                       # və ya kompose.io/installation
kompose convert -f docker-compose.yaml     # hər servis üçün deployment + service faylları
kubectl apply -f                           # (docker-compose.yaml-ı sil — kubectl onu da manifest sanır!)
kubectl get pods                           # status yoxla
```
- **Məhdudiyyət:** 1-dən çox env dəyişəni ilə işləmir — amma kitabda
  bütün env-lər onsuz da compose içindədir

### 10. Nəticə
- Kitabın yolu: docker-compose → Kubernetes (kompose) — SADƏLƏŞDİRİLMİŞ;
  real: CI/CD → build → bulud repo → çatdırma → işə salma (dəqiqələr və ya
  daha çox; güc + paralel əməliyyatlardan asılı)
- 5-ci üsul (serverless) — platforma spesifik təlimatla, kitabda işlənmir.

## Əsas terminlələr
- Load balancer (balanslaşdırıcı) — instanslar arası bərabər yükləmə
- CI/CD — fasiləsiz inteqrasiya/çatdırılma avtomatlaşdırması
- Workflow — YAML-da iş axını təsviri (GitHub Actions)
- Trigger/triqqer (on:) — workflow-u işə salan hadisə
- gomock/mock — interfeyslərin saxta implementasiyaları
- Table-driven test — cədvəl variativli test patterni
- Coverage (örtük) — testlərin əhatə etdiyi kod faizi
- Multi-stage build — build + runtime ayrı mərhələlər
- Layer cache — dəyişməyən qatların yenidən istifadəsi
- Docker registry (Docker Hub) — image nəşr anbarı
- Node/Pod/Deployment/Service — Kubernetes obyektləri
- Manifest — K8s üçün deklarativ YAML konfiqurasiya
- Declarative API — istənilən vəziyyəti bildir; sistem özü çatdırır
- kompose — compose→K8s konvertoru

## Praktik Nəticə

1. **Dockerfile sırası VACİBDİR:** go.mod/go.sum → download → COPY . . —
   layer cache asılılıq dəyişməsə kodun hər build-ində paket yüklənmir.
2. **Runtime ayrı edin:** alpine təmiz image ~5 MB — kompilyatorsuz,
   hücum səthi minimal.
3. **Konteyner şəbəkəsində host = servis adı:** DB_DSN, ACCOUNT_GRPC_HOST
   və s. — localhost YOX.
4. **CI-də -race + coverage:** `go test -v -race -coverprofile=...` —
   race detector uzaqdan da aktiv.
5. **Testlərdə xəta yolları:** yalnız uğur YOX — repo/kafka xətaları;
   mock-lar real asılılıqları kənarlaşdırır.
6. **K8s-ə keçid:** kompose convert + apply; imperative `kubectl run`
   YERİNƏ manifest (sürüm/history idarəsi).
7. **Yayım pilləliyi:** 1 üsuldan 4-ə — ehtiyac artdıqca; hər addım
   manual işi azaldır, idarəetmə mürəkkəbliyini artırır.

## Mənbə
Pages: 267-310 (PDF 268-311)
