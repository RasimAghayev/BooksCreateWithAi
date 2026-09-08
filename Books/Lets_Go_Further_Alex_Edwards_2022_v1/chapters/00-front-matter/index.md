# Chapters 0+1 — Front Matter & Introduction

## Bu hissə nədən bəhs edir?

Kitabın giriş hissəsi: Greenlight JSON API layihəsinin məqsədi, kimə uyğun
olduğu, prerequisite-lər və endpoint xəritəsi.

## Əsas fikirlər

### 1. Greenlight layihəsi
**Nədir:** Open Movie Database-ə bənzər JSON API — filmlərin idarəsi
(CRUD), istifadəçi qeydiyyatı, aktivasiya, authentication (kimlik
təsdiqi), authorization (səlahiyyət), CORS, metrics və deployment.

**Son məqsəd endpoint-ləri:**
| Method | URL | Action |
|---|---|---|
| GET | /v1/healthcheck | Status + versiya |
| GET/POST | /v1/movies | List / Create |
| GET/PATCH/DELETE | /v1/movies/:id | Show / Update / Delete |
| POST | /v1/users | Qeydiyyat |
| PUT | /v1/users/activated | Aktivasiya |
| PUT | /v1/users/password | Parol yenilə |
| POST | /v1/tokens/authentication | Token al |

### 2. Prerequisites (tələblər)
- "Let's Go" kitabının bilikləri (əsaslar yenidən izah olunmur; testing yoxdur)
- Go 1.20 (`go version` ilə yoxla)
- curl (HTTP sorğular), hey (load testing), PostgreSQL, golang-migrate, make

### 3. Kitabın strukturu
23 fəsil: JSON → DB → CRUD → filtering → logging → rate limiting → graceful
shutdown → users → email → activation → auth → permissions → CORS → metrics →
build/QC → deployment → appendices (JWT, password reset, context timeouts).

## Əsas terminlər
- JSON API — JSON formatında sorğu/cavab ilə işləyən API
- Authentication (kimlik təsdiqi) — istifadəçinin kimliyinin yoxlanması
- Authorization (səlahiyyətvermə) — istifadəçinin nəyə icazəsi olduğu

## Praktik nəticə
Kitab "follow-along" tiplidir — hər fəsil əvvəlkinin üzərinə qurulur; kodu
bərpa etmək üçün greenlight.alexedwards.net repo-dan istifadə olunur.

## Mənbə
Pages: 1-10 (raw 001-010)
