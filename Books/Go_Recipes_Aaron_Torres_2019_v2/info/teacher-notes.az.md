# Go Recipes (Torres) — Müəllim qeydləri

📖 **Kitab deyir:**

Kitab 85+ reseptdən ibarət praktik mənbədir. Əsas mövzuları: standart
kitabxananın interfeys-əsaslı istifadəsi, xəta emalı protokolu (pkg/errors),
mikroservis pattern-ləri (Storage interfeysi, Transport, middleware),
concurrency (channel, WaitGroup, pool, pipeline), distributiv sistemlər
(Consul, Raft, Docker), test (mock, table-driven, BDD), serverless və
performans alətləri (pprof, benchmark).

👨‍🏫 **Müəllim qeydi:**

Məncə bunu praktikada belə tətbiq etmək daha yaxşı olar:

1. **2019 tarixini nəzərə alın.** Kitabın əksər pattern-ləri hələ də
   aktualdır (io, concurrency, testing), amma bəziləri superesildi:
   `ioutil` → `os`/`io` paketlərinə köçdü (Go 1.16+); `net/rpc` → demək
   olar öldü, gRPC standartdır; `docker-compose` v2 sintaksisi dəyişdi;
   Sarama `-u` fork (IBM/sarama) istifadə edin; fasthttp yalnız HTTP/2
   lazım olmayan ekstremal hallarda.

2. **Mock üçün gomock-un yerinə artıq mocks/mockery və ya sadə əl ilə
   yazılmış closure mock-lar.** Kitabın Patch/Restore texnikası maraqlı
   refleksiya nümunəsidir, amma gündəlik işdə DI + interface mock
   (kitabın öz Storage pattern-i) daha təmiz həlldir.

3. **SQL NullTypes üçün müasir həll:** `*int`/`*string` pointer sahələri
   və ya `sql.Null*` tipləri hələ də işləyir, amma `pgx` + `database/sql`
   interfeysi kombinasiyası produksiyada standartdır.

4. **Worker pool üçün `errgroup.WithContext` (golang.org/x/sync)** —
   kitabın əl ilə yazdığı WaitGroup+error toplama pattern-ini bir paketə
   sığışdırır; öyrənəndə əvvəl kitabın əl ilə versiyasını yazın, sonra
   errgroup-a keçin.

## Ən vacib 5 fikir

1. **Funksiyalar interfeys qəbul etsin** (io.Reader, Storage, Client) —
   test, mock və backend dəyişməsi bir sətirdə həll olunur (ch 1, 6, 7).
2. **Xəta zənciri: Wrap ilə kontekst, Cause ilə orijinal; log yalnız son
   nöqtədə** (ch 4) — servis loqlarında duplikasiya və itirilmiş kök səbəb
   xətalarının qarşısını alır.
3. **Concurrency primitiv seçimi məqsədə görə:** sayğac → atomic; map →
   RWMutex; bir dəfəlik init → sync.Once; tapşırıq gözləmə → WaitGroup;
   mərhələli emal → pipeline (ch 10).
4. **Worker pool CPU-ağır işlərdə mütləqdir** — hər sorğu üçün yeni
   goroutine crypto/bcrypt kimi işlərdə prosesi boğar (ch 10).
5. **Optimallaşdırma ölçmədən əvvəl yox:** pprof dar boğazı göstərir,
   benchmark (RunParallel + benchmem) mikro-qərarları rəqəmlə doğrulayır
   (ch 14).

## Kitabın ən dəyərli hissəsi

**Chapter 9 (Test) + Chapter 10 (Concurrency)**, səh. 297-360 — çünki bu
iki fəsil kitabın bütün arxitektur pattern-lərinin (interfeys, DI, mock,
pool, pipeline) sınaq mühitinə çevrilmiş halını göstərir. Test yazmağı
öyrənən developer eyni anda SOL prinsiplərini də mənimsəyir. Həmçinin
ch 9-dakı table-driven + httptest kombinasiyası Go-da veb testlərinin
90%-ni əhatə edir.

## İşlənmə təqvimi (tövsiyə)

| Həftə | Mövzu |
|---|---|
| 1 | Ch 1-3: I/O, CLI, data — standart kitabxana əsasları |
| 2 | Ch 4-5: error handling + şəbəkə |
| 3 | Ch 6-7: bazalar + veb klientlər |
| 4 | Ch 8-9: mikroservis + testlər |
| 5 | Ch 10: concurrency (ən vacib fəsil — 2 dəfə oxuyun) |
| 6 | Ch 11-12: distributiv + reaktiv |
| 7 | Ch 13-14: serverless + performans |
