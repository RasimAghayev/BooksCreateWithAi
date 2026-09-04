# Müəllim qeydləri

Bu kitabı oxuduqdan sonra əsas nəticələr və praktik tövsiyələr.

## Ən vacib 5 fikir

1. **Mikroservis seçimi həmişə məqsədə əsaslanmalıdır.** Monolit
   sadəlik, mikroservislər isə miqyaslana bilmə və müstəqillik təklif edir.
   Layihənin hədəflərinə və komanda ölçüsünə əsasən qərar verin.

2. **Go dili yüksək performanslı, distribut edilmiş sistemlər üçün
   nəzərdə tutulmuşdur.** Goroutines və channels vasitəsilə sadə
   paralellik, güclü standard library ilə az asılılıq — bu onu
   mikroservislər üçün ideal edir.

3. **Go Modules standartdır.** Köhnə $GOPATH və paket menecerləri
   (godep, glide) artıq relevant deyil. go.mod və go.sum ilə
   asılılıqları idarə edin.

4. **Paylanmış tranzaksiyalar ən çətin hissədir.** Saga pattern-i
   vaxtilə 2PC-də daha yaxşı seçim ola bilər. Distributed systems
  -də data consistency həmişə kompromis tələb edir.

5. **Docker və Kubernetes production üçün deyil, inkişaf və
   yerləşdirmə üçün standartdır.** CI/CD, ArgoCD, Vault kimi
   alətlərin inteqrasiyasını əvvəlcədən nəzərdə tutun.

## Praktik yanaşma

### Tövsiyə 1: Lokal mühiti tez quraşdırın
Kitabda göstərilən addımları izləyin:
1. Visual Studio Code + Go plaqin
2. Go quraşdırma (https://go.dev/dl/)
3. Docker Desktop
4. PostgreSQL və ya Docker ilə bazanı işə salın

### Tövsiyə 2: Go Modules ilə işləyin
Köhnə $GOPATH strukturundan uzaq durun. `go mod init` ilə
başlayın, `go mod tidy` ilə asılılıqları təmizləyin.

### Tövsiyə 3: Təmiz arxitektura tətbiq edin
Chapter 1-də göstərilən strukturu izləyin:
- `internal/` — daxili paketlər
- `pkg/` — xaricdən istifadə edilə bilən paketlər
- `cmd/` — tətbiq giriş nöqtələri
- `api/` — handler və routerlar

### Tövsiyə 4: Testləməni əvvəlcədən qurun
- Unit test-lər: `go test`
- Integration test-lər: test bazası və ya Docker container
- API test-lər: Postman və ya curl

### Tövsiyə 5: Docker-dan istifadə edərək development qurun
```bash
docker-compose up -d
```
Bu yanaşma hər developer üçün eyni mühiti təmin edir.

## Kitabın ən dəyərli hissəsi

**Chapter 1, pages 13-86** — Bu chapter tam bir mikroservis
hazırlanmasını əhatə edir: kod redaktordan, Go quraşdırmasından,
verilənlər bazası layihələndirilməsindən və testləməyə qədər.
Bu bölməni 2-3 dəfə təkrar edin, hər addımı özünüz təkrar edin.

## Diqqət edilməli nöqtələr

1. **Versiya idarəsi:** Go versiyasını dəyişərkən poetik addımlar
   istifadə edin. Kənardan bir neçə versiya keçmək risklidir.
2. **Təhlükəsizlik:** Kitabda authentication və authorization
   bölməsi var. JWT, refresh token, RBAC konseptlərini yaxşı
   başa düşün.
3. **Paylanmış tranzaksiyalar:** Chapter 4-də 2PC və Saga
   pattern-ləri izah edilir. Real layihələrdə Saga tez-tez daha
   yaxşı seçimdir.
4. **Deployment:** Chapter 5-də Docker və Kubernetes əhatə edilir.
   Productiona çıxış üçün bu biliklər vacibdir.
5. **Kod keyfiyyəti:** Kitabda codestyle, logging, error handling
   kimi mövzular var. Bunları öz layihənizdə tətbiq edin.

## Kitabın hansı sferaya aid olduğu

Bu kitab **Backend Development (Backend development / server tərəfi
inkişafı)** və **Microservices Architecture (Mikroservis
 Arxitekturası)** sferalarına aiddir.

Əsas texnologiyalar: Go, PostgreSQL, Docker, Kubernetes, gRPC,
Kafka, Redis.
