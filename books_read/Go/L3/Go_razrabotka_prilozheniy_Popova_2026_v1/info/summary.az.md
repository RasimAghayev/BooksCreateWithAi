# Go: разработка приложений в микросервисной архитектуре с нуля — Xülasə (Azərbaycanca)

**Müəllif:** Попова Ю.Ю. (Julia Popova) | **Nəşriyyat:** БХВ-Петербург | **İl:** 2026 | **Səviyyə:** L3 (Intermediate)

## Kitabın ümumi məqsədi

Bu kitab Go-da mikroservis arxitekturasını SIFIRDAN, tam praktik yolla öyrədir:
oxucu bir kitab boyu VAHİD layihə qurur — istifadəçi servisi (Account/User)
ilə başlayır, avtorizasiya (Auth), API Gateway, əməliyyatlar (Transaction)
servislərini əlavə edir və sonda bütün sistemi Docker konteynerlərinə
yığıb Kubernetes-də yayımlayır. Hər fəsil real kodla (GORM, goose, zerolog,
gRPC, Protobuf, Kafka) gedir; kontraktlar ayrı repozitoriyada, hər qatın öz
modeli və mapper-ləri var. Kitab rus dilindədir və BHV-nin "С нуля"
(sıfırdan) seriyasındadır.

## Fəsil-fəsil xülasə

0. **Введение:** mikroservisə keçid biznes qərarıdır; monolitin 5 problemi
   (komanda konflikti, təhlükəsizlik, anlaşılmazlıq, gec reliz, tək nöqtə
   fəlakəti); paylanmış monolit təhlükəsi; Go = goroutine+kanal, sürətli
   compile, minimализм.
1. **User (Account) servisi:** Go Modules tarixi ($GOPATH→godep→glide→
   govendor→modules; go.sum hash-ləri Supply Chain Attack-a qarşı), layihə
   kanonu (cmd/internal/pkg; internal xaricə QAPALI), Makefile, .env+
   caarlos0/env, zerolog strukturqlu log, GORM repository pattern, Swagger+Gin,
   PostgreSQL normal formalar (1-3НФ), Protobuf kontraktlar ayrı repo
   (cyclic import qadağası, Docker+Make deterministik generasiya), 3-qatlı
   arxitektura (hər qatın ÖZ modeli — gRPC-dən asılı biznes ANTİPATTERN),
   goose tranzaksiya miqrasiyaları, DI internal/app kompozitoru, Postman gRPC.
2. **Auth servisi:** identifikasiya/autentifikasiya/avtorizasiya ayrımı;
   5 üsul (parol, sertifikat X.509, OTP/2FA, API key, token); HTTP sxemlər
   (Basic/Digest); JWT/SWT/SAML formatları; zəifliklər və müdafiələr (buffer
   overflow→validasiya, XSS/CSRF/clickjacking); users+refresh_tokens cədvəl
   (FK CASCADE, revoked_at); bcrypt (DefaultCost=10, salt); access JWT (HS256,
   sub/exp) + refresh SHA-256 DB-da; Register/Login/Refresh/Validate/Logout.
3. **Servislərarası əlaqə:** sinxron/asinxron; OSI 7 qat; HTTP metodları;
   gRPC (Protobuf, 4 kommunikasiya nümunəsi, enum-0 problemi, google.api.http
   REST annotations); RabbitMQ/AMQP (fanout/direct/topic exchange,
   acknowledgment, At Most/Least/Exactly Once); Kafka (pull, partition/offset,
   compact topic, Connect); Redis (in-memory əlavə, cache); API Gateway
   pattern (tək giriş, cross-cutting); JWT interceptor (publicMethods, Bearer,
   alg yoxlaması, context.WithValue; /users/me üçün GetUserIDFromContext);
   grpc.NewClient + insecure; servislərin interfeys arxasında izolyasiyası;
   kompozit Register (parol YALNIZ auth-da!).
4. **Transaction modulu:** НФБК/4-5-6НФ/ДКНФ; miqrasiya mexanikası (sütun
   rename→çevir→sil; enum 3-addım); indekslər (B-tree, rollback TƏRS sıra);
   ACID; paralel anomaliyalar (lost update, dirty/non-repeatable/phantom,
   serialization anomaly); izolyasiya səviyyələri (PostgreSQL fərqləri);
   double-entry (Transaction + Entry DEBIT/CREDIT); pul = QƏPİK int64;
   distributed tranzaksiya problemi → 2PC (koordinator, blok, yavaş) vs Saga
   (event-based, Kafka transaction_data/transaction_response topic-lər,
   pending→completed/failed, kompensasiya; kafka-go Producer/Consumer/Client
   wrapper, unique groupID per topic); YAGNI (DELETE YOX).
5. **Yayım:** 5 üsul; nginx proxy; CI/CD GitHub Actions (workflow YAML,
   workflow_dispatch→push trigger, actions/checkout); gomock+testify
   table-driven testlər (84.4% coverage, xəta yolları mütləq); multi-stage
   Dockerfile (alpine, layer cache: go.mod əvvəl); .dockerignore; tam compose
   (host=SERVİS ADI, depends_on); Docker Hub tag/push; Kubernetes (Node/Pod/
   Deployment/Service, kubectl apply/get/describe/logs, MANIFEST imperativdən
   üstün); kompose convert miqrasiyası.
6. **Nəticə:** tam arxitektura xəritəsi (5 servis + 3 DB + Kafka), kod arxivi
   (zip.bhv.ru).

## Kitabın əsas mesajları

1. **Mikroservis = qatlı intizam:** transport/biznes/data ayrılığı, hər
   qatın öz modeli, interfeyslərlə asılılıq — dəyişiklik LOKALIZASIYA olunur.
2. **Kontraktlar birinci vətəndaşdır:** ayrı repo, semver tag-lər, Docker
   ilə deterministik protoc generasiyası, cyclic import = modulluq pozuntusu.
3. **Pul və parol heç vaxt açıq saxlanmır:** qəpiklərdə int64; bcrypt+salt;
   parol yalnız Auth servisinə gedir.
4. **Paylanmış tranzaksiya = dizayn problemi:** mümkünsə qaçın; əks halda
   Saga (status + kompensasiya) praktik standartdır.
5. **Avtomatlaşdırma pilləli götürülür:** Makefile → CI/CD → Docker →
   Compose → K8s — hər addım manual xərci azaldır.

## Kitabın güclü və zəif cəhətləri

**Güclü:** tam layihə fokuslu (kod arxivi ilə); rusdilli ekosistemdə nadir
mikroservis Go materialı; real prodakstion idiomlar (DI, interceptor, Saga).
**Zəif:** Kubernetes səthi (kompose keçidi, prodakstion üçün kifayətsiz);
monitorinq/observability (Prometheus/Loki yalnız adı çəkilir); bəzi kod
təkrarları (mapper-lər, app.go) — arxivdə yoxlanılmalı.
