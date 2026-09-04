# Chapter 5 — Развертывание микросервисов

Pages: 267-310

## Bu chapter nədən bəhs edir?
Bu chapter hazırlanmış mikroservislərin necə konteynerləşdirildiyini və
Kubernetes orkestratoru ilə necə miqyaslandırıldığını izah edir.

## Əsas fikirlər
### 1. Docker konteynerləşdirmə
- Dockerfile yazmaq
- Multi-stage build
- Docker Compose ilə çoxservisli tətbiq
- Volume və port yönləndirməsi

### 2. Kubernetes orkestrasiyası
- Pod, Deployment, Service, Ingress anlayışları
- ConfigMap və Secret
- Health checks (liveness, readiness)
- Horizontal Pod Autoscaler (HPA)

### 3. CI/CD (CI/CD bunun üçün deyil, amma yaxşıdır)
- Image build və push
- Deployment strategiyaları

## Əsas terminlər
- Containerization (konteynerləşdirmə)
- Docker (konteyner platforması)
- Docker Compose (çoklu konteyner orkestrasiyası)
- Kubernetes (K8s, konteyner orkestrasiya platforması)
- Pod (Kubernetes-ın ən kiçik deployable vahidi)
- Deployment (Kubernetes deploy obyekti)
- Service (Kubernetes servis abstraksiyası)
- Ingress (HTTP/HTTPS trafik idarəcisi)
- ConfigMap (konfiqurasiya məlumatları)
- Secret (məxfi məlumatlar)
- Health Check (sağlamlıq yoxlaması)
- Autoscaling (avtomatik miqyaslama)
- Load Balancer (yük bərabərləndirici)
- Rolling Update (tərpənən yeniləmə)

## Praktik nəticə
Bütün mikroservislər Docker konteynerlərə qablaşdırılır və Kubernetes
klasterə yerləşdirilir.

## Mənbə
Pages: 267-310
