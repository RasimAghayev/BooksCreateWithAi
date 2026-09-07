# Go: разработка приложений в микросервисной архитектуре с нуля — Mündəricat

**Müəllif:** Попова Ю.Ю. | **Nəşriyyat:** БХВ-Петербург | **İl:** 2026 | **Səviyyə:** L3 (Intermediate) | **Dil:** Rus

Kitabın mövzusu: Go-da mikroservis arxitekturasının SIFIRDAN qurulması — bir
User servisindən başlayaraq Auth (JWT), Gateway, Transaction (Saga/Kafka)
və Docker/Kubernetes yayımınadək tam yol. Praktik layihə yönümlü: hər fəsil
real kod (GORM + goose + zerolog + gRPC + Kafka) ilə gedir.

## Fəsillər

| # | Fəsil | Səhifə | Əsas mövzu |
|---|-------|--------|------------|
| 0 | [Введение](../chapters/00-introduction/index.md) | 7-12 | Monolit vs mikroservis, Go icmalı, texnologiya steki |
| 1 | [User servisi](../chapters/01-user-service/index.md) | 13-86 | Alətlər, layihə strukturu, GORM, miqrasiya, DI, 3-qatlı arxitektura |
| 2 | [Auth servisi](../chapters/02-auth-service/index.md) | 87-110 | Autentifikasiya nəzəriyyəsi, JWT/bcrypt, token lifecycle |
| 3 | [Servislərarası əlaqə](../chapters/03-inter-service-communication/index.md) | 111-156 | OSI/HTTP/gRPC, Kafka/RabbitMQ/Redis, Gateway + interceptor |
| 4 | [Transaction modulu](../chapters/04-transaction-module/index.md) | 157-266 | Normal formalar, ACID, izolyasiya, 2PC vs Saga, kafka-go |
| 5 | [Yayım (Deployment)](../chapters/05-deployment/index.md) | 267-310 | Docker multi-stage, CI/CD (GitHub Actions), Kubernetes |
| 6 | [Nəticə](../chapters/06-conclusion/index.md) | 311-317 | Arxitektura xəritəsi, texnologiya siyahısı, arxiv |

## Oxu ardıcıllığı tövsiyəsi

Ardıcıl 0 → 1 → 2 → 3 → 4 → 5: hər fəsil əvvəlkinin layihəsini davam etdirir
(tək kitab boyu BİR mikro sistem qurulur). Fəsil 4 ən uzun və ən çətindir
(normal formalar + tranzaksiyalar + Saga) — 3-cü fəsldəki Kafka nəzəriyyəsi
ondsuz anlaşılmır.

## Kitabın fayl arxivi

Final layihə kodu: https://zip.bhv.ru/9785977521208.zip (nəşriyyat serveri)
