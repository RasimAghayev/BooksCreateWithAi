# Build systems with Go — Müəllim qeydləri (AZ)

## Ümumi qiymətləndirmə

Kitab Juan M. Tirado tərəfindən yazılıb və Go-nun demək olar bütün
əsas səthlərini (dil nüvəsi + sistem qurulması) kiçik, işlək nümunələrlə
əhatə edir. Kod nümunələri müstəqildir — hər biri `main.go` kimi icra
oluna bilir; real layihə strukturu göstərilmir, amma bu, dərs kitabı
kimi qəbul ediləndə güc deyil, zəiflikdir.

## Kitabın deyil, praktikanın tələb etdiyi əlavələr

1. **Project layout** — kitabda `cmd/` strukturu yalnız Cobra fəslində
   görünür. Real layihədə: `cmd/app/main.go` + `internal/` bölgüsü
   standartdır.
2. **errgroup** — çoxsaylı paralel işlərin xətasını bir yerdə toplamaq üçün
   (`golang.org/x/sync/errgroup`) — kitabda yoxdur, amma hər hansı
   fan-out işində faydalıdır.
3. **Graceful shutdown** — `signal.NotifyContext` + `http.Server.Shutdown`
   birləşməsi production HTTP serverləri üçün vacibdir.
4. **Viper** — Cobra ilə birlikdə konfiqurasiya idarəetməsi (kitab yalnız
   flag-ləri göstərir).
5. **Migration alətləri** — GORM AutoMigrate dev üçündür; production-da
   versiyalı migrasiya alətləri (goose, golang-migrate) üstünlük təşkil edir.

## Diqqət çilən nöqtələr (tələbələr üçün)

- **Chapter 5 (Reflection)** — üç qanun mütləş yadda saxlanmalıdır;
  xüsusilə "settable" qanunu: `ValueOf(&x).Elem()` olmadan yazmaq
  həmişə panic verir.
- **Chapter 6 (Context)** — loop daxilində yaranan kontekstlərin
  `cancel()`-i defer EDİLMƏZ — resurs bitməz. Kitab bunu vurğulayır —
  nümunələrdə də eynilə tətbiq olunub.
- **Chapter 9 (HTTP)** — middleware zənciri `ApplyMiddleware` ilə
  proqramatikləşir; tələbələr bunu öz layihələrində routerə bağlamalıdır.
- **Chapter 14 (gRPC)** — `WithInsecure` yalnız test üçündür;
  production-da TLS məcburidir. Auth interceptor nümunəsi tədris
  səviyyəsindədir.
- **Chapter 17 (GORM)** — hər yazma default transaksiyadadır;
  `SkipDefaultTransaction: true` performans qazandırır, amma çoxsaylı
  yazmaların atomikliyi itir.

## Kitabın ən dəyərli hissəsi

Chapter 6 (Concurrency, səh. 102-142) və Chapter 14 (gRPC, səh. 282-325) —
birincisi Go-nun sinonimi olan sahəni tam əhatə edir, ikincisi müasir
mikroservis rabitəsinin (PB + gRPC + REST transcoding + interceptor) bütöv
mənzərəsini verir.

## Nə oxumaq davam etmək üçün

- **Concurrency in Go** (Katherine Cox-Buday) — ch6-nın dərinləşdirilməsi
- **100 Go Mistakes** (Teiva Harsanyi) — bu kitabın bütün mövzularının
  "nə etməməli" tərəfi
- **Let's Go Further** (Alex Edwards) — HTTP + layihə strukturu baxımından
  tamamlama
