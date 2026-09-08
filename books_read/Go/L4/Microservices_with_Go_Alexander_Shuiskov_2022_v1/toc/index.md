# Microservices with Go — Mündəricat (TOC)

**Müəllif:** Alexander Shuiskov (Uber staff engineer, observability) · **İl:** 2022 · **Packt Publishing** · 329 PDF səh. (290 kitab səh.)

**Səhifə ofseti:** kitab səh. + 23 = PDF səh.

## Struktur

| Hissə | Fəsil | Səh. (kitab) | PDF | Mövzu |
|---|---|---|---|---|
| 1 | 1 | 1–12 | 24–35 | Mikroservislərə giriş: nədir, nə üçün, Go-nun rolu |
| 2 | 2 | 13–54 | 36–77 | Movie aplikasiyasının qurulması (scaffolding), metadata/rating/movie servisləri |
| 2 | 3 | 55–76 | 78–99 | Service Discovery — Consul, registry, health monitoring |
| 2 | 4 | 77–92 | 100–115 | Serialization — JSON, Protocol Buffers, msgpack |
| 2 | 5 | 93–114 | 116–137 | Sinxron kommunikasiya — gRPC gateway/client |
| 2 | 6 | 115–132 | 138–155 | Asinxron kommunikasiya — Apache Kafka, partitioning, versioning |
| 2 | 7 | 133–146 | 156–169 | Data saxlanması — MySQL, database intro |
| 2 | 8 | 147–162 | 170–185 | Kubernetes deployment, canary, rollback, CD |
| 2 | 9 | 163–192 | 186–215 | Unit/Integration testlər — mocking, cmp |
| 3 | 10 | 193–214 | 216–237 | Reliability — graceful shutdown, on-call, incident management |
| 3 | 11 | 215–250 | 238–273 | Telemetry — logging, metrics, tracing, OpenTelemetry |
| 3 | 12 | 251–266 | 274–289 | Alerting — Prometheus alert rules |
| 3 | 13 | 267–290 | 290–313 | Advanced — profiling, dashboards, frameworks, JWT |

## Layihə davamlılığı

Kitab boyu **Movie application** layihəsi: movie metadata, rating və movie aqreqasiya
mikroservisləri sıfırdan qurulur, sonra discovery/serialization/Kafka/MySQL/K8s/observability
ilə production-hazır hala gətirilir.
