# Microservices with Go — Xülasə (AZ)

## Kitab kimin üçündür?
Go-nu bilib, mikroservis arxitekturasını UÇDAN-UCA öyrənmək istəyənlər üçün —
kitab boyu Movie aplikasiyası (metadata, rating, movie servisləri) sıfırdan
qurulur və production-keyfiyyətinə çatdırılır: discovery, gRPC, Kafka, MySQL,
Kubernetes, test, observability, alerting, JWT.

## Layihə xətti (13 fəsil, 290 kitab səh.)

| Hissə | Fəsil | Mövzu | Nə tikilir |
|---|---|---|---|
| 1 | 1 | Mikroservis konsepti | monolit vs mikro, meyarlar |
| 2 | 2-9 | Foundation | Movie app scaffolding, Consul discovery, Protobuf, gRPC, Kafka, MySQL, K8s, testlər |
| 3 | 10-13 | Maintenance | reliability, telemetry, alerting, profiling/dashboard/JWT |

## Ən vacib 5 fikir
1. **Qatlı servis strukturu:** handler → controller → repository/gateway —
   biznes məntiqi API növündən asılı deyil; HTTP-dən gRPC-ə keçiddə controller
   DƏYİŞMİRDİ
2. **Texnologiya-neytral interfeyslər:** discovery Registry, ratingIngester —
   in-memory (test) və Consul/Kafka (prod) implementasiyaları dəyişdirilə bilər
3. **Explicit xətalar + retry məntiqi:** errors.Is/As, sentinel ErrNotFound
   hər komponentdə ayrı; retriable kodlar (DeadlineExceeded,
   ResourceExhausted, Unavailable) + exponential backoff + jitter
4. **Observability 3 sütunu:** zap structured logging, tally/Prometheus
   metrics (cardinality diqqəti!), OpenTelemetry tracing (context
   propagation — interceptorlər avtomatik span yaradır)
5. **Reliability = limitlərin açıq təyini:** 429 rate limit (token bucket),
   graceful shutdown (signal.Notify + GracefulStop), fallback/graceful
   degradation, canary/CD, on-call/runbook/postmortem mədəniyyəti

## Kitabın ən dəyərli hissəsi
Ch2 (scaffolding qatları) + Ch5 (gRPC-ə köçürmə) kombinasiyası: internal vs
generated model ayırımı, mapper funksiyaları, UnimplementedServer embed —
real dünyada protokol dəyişikliyini ağrısız edən naxış.

## Öyrənilən texnologiya steki
Go · gRPC · Protocol Buffers · Apache Kafka (confluent-kafka-go) · Consul ·
MySQL (database/sql) · Docker · Kubernetes (minikube) · zap · tally ·
Prometheus + Alertmanager · Grafana · Jaeger + OpenTelemetry · golang-jwt ·
testify · cmp · gomock
