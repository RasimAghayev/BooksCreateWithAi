# gRPC Microservices in Go — Müəllim Qeydləri (Teacher Notes)

> Bu sənəd kitabı öyrədərkən istifadə üçün metodiki qeydlər, çətin anlar, müzakirə sualları və praktik tapşırıqlar.

---

## 1. Auditoriya və ön şərtlər

- **Hədəf auditoriya:** Go bilən intermediate developer-lər; mikro servis anlayışı faydalıdır amma kitab izah edir.
- **Ön şərtlər:** Go sintaksisi, Docker əsasları, əsas networking (port/TLS), protoc qurulmuş.
- **Həcm:** 200 səh., 9 fəsil — 3 aylıq kurs (həftədə 1 fəsil + lab) və ya 5 günlük intensiv.

## 2. Tədris axını — e-commerce layihəsi paralelində

| Həftə | Chapter | Layihə mərhələsi |
|-------|---------|------------------|
| 1 | 1-2 | Arxitektura müzakirəsi; mono→micro keçid planı |
| 2 | 3 | Proto repo + CI generasiya + compatibility testləri |
| 3 | 4 | Order servisi tam qurulur (hexagonal) |
| 4 | 5 | Payment inteqrasiyası + status/error modeli |
| 5 | 6 | Timeout/retry/CB + mTLS sertifikat axını |
| 6 | 7 | Unit→integration→E2E piramidi |
| 7 | 8 | Docker + minikube + Ingress + cert-manager |
| 8 | 9 | OTel + Jaeger + Kibana observability |
| 9 | — | Final: tam stack demo + strategy təqdimatı |

**Unikal imkan:** Kitabın GitHub repo-su Actions ilə listing icra edir — şagirdlər lokal mühit qurmadan kodu "işlədə" bilər (listing_execute.yaml).

## 3. Çətin anlar və izah üsulları

### a) Field number əhəmiyyəti (Ch 3)
Canlı demo: iki versiyalı proto → köhnə client + yeni server → sahə nömrəsi dəyişilsə DEKODİROVKA pozulur. "Field number = sahənin pasportu" metaforası.

### b) Binary encoding (Ch 3)
Whiteboard-da 65 → `0001000 01000001` bit parçasını birgə əldə et (varint + MSB). Sonra 21567 üçün 3-blok split — "protocoğrafiya" o qədər də sehrli deyil.

### c) Hexagonal asılılıq istiqaməti (Ch 4)
Diaqram çəkin: core ← ports ← adapters. Xəbərdarlıq: adapter PORTU tanıyır, əksinə YOX. Mock test bu istiqamətin mükafatıdır —birgə mock yazıb göstərin.

### d) Saga iki variantı (Ch 2)
Rol oyunu: şagirdlər servis rolunda — choreography = "növbətinə event at, özünü idarə et"; orchestrator = bir şagird dirijor, əmrləri verir, compensation-ları geri çağırır.

### e) CB state machine (Ch 6)
Fanla 3 vəziyyət kartı: CLOSED (işləyir) → OPEN (dərhal fail) → HALF-OPEN (sınaq). ReadyToTrip riyaziyyatını birgə hesabla (6 fail / 10 sorğu = 0.6).

### f) Percentile vs average (Ch 9)
Klassik tələ: 1000×1ms + 10×4s → orta 5ms. Müştəri kimi düşün: "5ms orta" sözü yalandır — 10 müştəri 4 saniyə gözləyir! p95 hesabla birgə.

### g) mTLS sertifikat zənciri (Ch 6)
CA → server → client imza zəncirini çək; sonra openssl əmrlərini canlı icra — faylların yaranmasını gör. "Evdən çıxarkən iki pasport" (server + client sert) metaforası.

## 4. Müzakirə sualları

1. Monolith-dən mikro servisə keçid hansı SIQNALDA başlamalı? (scale/velocity/team)
2. Choreography vs orchestrator saga: 3 servisli flow üçün hansı? 30 servisli üçün?
3. Retry InvalidArgument üçün niyə zərərdir? Hansı kodlar retry "layiqdir"?
4. RollingUpdate zamanı v1+v2 paralel işləyir — protobuf oneof dəyişikliyi bu ssenaridə nə edir?
5. Hər servisə ayrı LoadBalancer niyə mikro servisdə antipattern? Ingress nə vaxt KIFAYƏT DEYİL?
6. Canary vs Blue-Green: hansı halda dublikat infra xərcini doğruldur?
7. Mock test DB-ni tutmadığı hallar hansılardır? (integration testin mövcudluq səbəbi)
8. Trace_id olmadan Kibana-da bir sifarişin 5 servislik logunu necə taparsan?

## 5. Praktik tapşırıqlar

**Tapşırıq 1 (proto):** 3 servisli (order/payment/shipping) proto repo qur; GitHub Actions matrix generasiya; v1→v2 dəyişikliyində compatibility testi (köhnə client yeni server-ə qoşulur).

**Tapşırıq 2 (hexagonal):** Order servisini kitabdan asılı olmayaraq qur (məs. "delivery" domeni); port/adapter/domain ayrılığı; GORM + AutoMigrate.

**Tapşırıq 3 (kommunikasiya):** Order→Payment inteqrasiya: port + adapter + status.Errorf + errdetails; client tərəfdə Convert ilə 3 xəta ssenari assert.

**Tapşırıq 4 (resiliency):** 1) context timeout-lu server (sleep) + client; 2) retry interceptor (Unavailable kodunda); 3) gobreaker CB — 3 paternin hamısını bir testdə birləşdir.

**Tapşırıq 5 (mTLS):** CA+server+client sertifikatları openssl ilə; Go-da mTLS server+client; TLS-siz client rədd edilsin.

**Tapşırıq 6 (test piramidi):** Unit (mockery+testify), integration (Testcontainers MySQL), E2E (docker-compose stack) — üçünün coverage hesabatı.

**Tapşırıq 7 (K8s):** minikube + Ingress (GRPC annotation) + cert-manager selfsigned; canary deployment təcrübəsi (2 deployment eyni selector).

**Tapşırıq 8 (observability):** OTel interceptor + Jaeger trace baxışı; logrus custom formatter trace_id inyeksiyası; Kibana-da TraceID filtri.

**Final:** Tam e-commerce stack (3 servis + DB + observability) K8s-də + deployment strategy təqdimatı.

## 6. Sınav sualları

**Asan:**
1. protoc generasiya edən 2 fayl? (order.pb.go + order_grpc.pb.go)
2. proto3-da tez sahələr üçün hansı nömrə aralığı? (1-15)
3. gRPC xətalarında Code sahəsi olmadan status.Errorf YOXSA xam err — fərq? (Unknown)

**Orta:**
4. Saga-nın compensation transaction-ı nə vaxt işə düşür? Hansı istiqamətdə?
5. mTLS ilə TLS-in fərqi? Hansı halda mütləq tələb olunur (zero-trust)?
6. Ingress ilə Service LoadBalancer arasındakı xərc fərqi?

**Çətin:**
7. RollingUpdate + oneof silinməsi birləşəndə nə baş verir? Izah et və həll təklif et.
8. CB half-open state-də MaxRequests niyə məhduddur? Limitsiz olsaydı?
9. p95 = 100ms, average = 5ms — müştəri müqaviləsində hansını yazarsan və niyə?

## 7. Kollektiv layihə ideyası

**"gRPC Distributed Doğu":** Qruplar 3-4 servisli domain qurur (e-commerce/food delivery/booking); ümumi proto repo; interservis kommunikasiya müqaviləsi; ortaq K8s klasterində birgə deploy + observability dashboard-un birgə oxunması.

## 8. Əlavə resurslar

- Kitab repo: github.com/huseyinbabal/grpc-microservices-in-go
- gRPC rəsmi: grpc.io · Protobuf: developers.google.com/protocol-buffers
- go-grpc-middleware: github.com/grpc-ecosystem/go-grpc-middleware
- gobreaker: github.com/sony/gobreaker
- Testcontainers: testcontainers.org · mockery: github.com/vektra/mockery
- cert-manager: cert-manager.io · Jaeger: jaegertracing.io · OpenTelemetry: opentelemetry.io
- Google SRE kitabı: sre.google/sre-book
- Effective Go və bu seriyadakı əvvəlki kitablar (kitabın 3 hissəsi)

## 9. Kitabın dəyəri və limitləri (müəllim qiyməti)

**Güclü:** Tək davamlı layihə; hexagonal + port/adapter dərallı DI; resiliency patternlərinin hamısı kod səviyyəsində; compatibility qaydalarının konkret ssenariləri; K8s + observability tam dövrü; GitHub Actions icra modeli unikal.

**Zəif/həssas:** Streaming kodu az (əsasən unary); saga yalnız nəzəri (event broker implementasiyası yoxdur); auth/authz (JWT) geniş əhatə olunmur; bəzi alətlər köhnəlir (dgrijalva/jwt-go kimi library versiyaları); K8s YAML-ların tam siyahısı yoxdur (linklərə istinad).
