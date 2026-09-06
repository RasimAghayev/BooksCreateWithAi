# Chapter 6 — Все дело в надежности

## Bu chapter nədən bəhs edir?

"Dayanıqlılıq" (reliability) anlayışını tərif edir: sistemin əsas funksiyasını həyata keçirməyə davam edə bilməsi, əgər ətraf mühit fundamental olaraq etibarsızdırsa. Jean-Claude Laprie'nin dörd pilləli dayanıqlılıq modeli, Twelve-Factor App metodologiyası və immutabl infrastruktur konsepsiyası izah edilir. Bu chapter konseptualdir, kod nümunəsi yoxdur.

## Əsas fikirlər

### 1. Dayanıqlılığın tərifi (Reliability)
**Nədir:** Sistemin gözlənilən vəziyyətdə fəaliyyət göstərməsi qabiliyyəti, əsas məqsədinə uyğun olaraq.

**Necə işləyir:** Dayanıqlılıq yalnız "çox vaxt işləmək" deyil — "gözlənilməyən hallarda da düzgün işləmək" deməkdir. Bulud sistemləri fundamental olaraq etibarsız komponentlər (şəbəkə, disk, elektrik) üzərində qurulur. Dayanıqlılıq bu komponentlərin nöqsanlarını qəbul edərək dizayn edilməlidir.

### 2. Jean-Claude Laprie'nin dayanıqlılıq modeli
**Nədir:** Dayanıqlılığın dörd əsas mexanizmi:
- **Xətanın qarşısını almaq (Fault prevention):** Test, kod review, dil xüsusiyyətləri ilə xətaları əvvəlcədən aradan qaldırmaq.
- **Xətaya davamlılıq (Fault tolerance):** Redundans (artıq nümunələr), ləüzumi enerji itkisi (graceful degradation) ilə xəta zamanında da funksiyanı saxlayış.
- **Xətanı aradan qaldırmaq (Fault removal):** Debugging, yamaqlar (patch), proqramların yenilənməsi.
- **Xətanı təxmin etmək (Fault forecasting):** Monitorinq, gözlənilirlik (observability) ilə potensial problemləri əvvəlcədən aşkar etmək.

### 3. Twelve-Factor App metodologiyası
**Nədir:** SaaS və bulud tətbiqləri hazırlamaq üçün 12 ilk (factor), James麦凯ney və Adam Wiggins tərəfindən hazırlanmış metodologiya.

**12 ilk:**
1. **Codebase (Kod bazası):** Birlikdə idarə olunan versiya, bir neçə mühit (dev/staging/prod) üçün eyni kod
2. **Dependencies (Asılılıqlar):** Açıq şəkildə elan edilmiş xarici asılılıqlar, sistemdə olanlara əsaslanmayın
3. **Config (Konfiqurasiya):** Konfiqurasiya kodda deyil, mühit dəyişənlərində (env vars) saxlanır
4. **Backing services (Arxa xidmətlər):** Verilənlər bazası, cache, mesaj növü kimi resurslar attachable olmalıdır
5. **Build, release, run (Qurulum, buraxılış, icra):** Müstəqil mərhələlər — build → release (config + build) → run
6. **Processes (Proseslər):** Stateless proseslər — hər sorğu öz data ilə gəlir, diskdə sessiya saxlamayın
7. **Data isolation (Məlumat izolyasiyası):** Port-binding vasitəsilə bir neçə versiya eyni serverdə işlə bilir
8. **Concurrency (Konkurensiya):** Proses nümayəndələri (process scaling) ilə miqyaslama
9. **Disposability (İstifadə edilib atılabilmə):** Proseslər tez başlaya və dayana bilməlidir — SIGTERM ilə təmiz şəkildə bağlanma
10. **Dev/prod parity (Tərtibat/iskəslənmə bərabərliyi):** Tərtibat və iskəslənmə mühitləri minumum fərqlənməlidir
11. **Logs (Jurnal qeydləri):** Loglar stdout/stderr-a axıdılır, xidmət onları idarə etmir
12. **Admin processes (İdarə prosesləri):): Admin əməliyyatları eyni codebase-dən, eyni environment ilə icra olunur

### 4. Immutabl infrastruktur (Dəyişməz infrastruktur)
**Nədir:** Serverlərə "ev heyvanı" (pet) kimi deyil, "mal" (cattle) kimi yanaşma — problem olduqda serverı yamaq yerinə yenidən qurub əvəz etmək.

**Necə işləyir:**
- Konfiqurasiya idarəetmə alətləri (Ansible, Chef, Puppet) ilə serverləri kod şəkildə təyin etmək
- Konteynerlər (Docker) və orchestration alətləri (Kubernetes) ilə immutabl imejlər
- "Pets" (adlı, özəl qulluq tələb edən serverlər) əvəzinə "Cattle" (eyni olan, istənilən vaxt əvəz edilə bilən)
- Yenidən qurmaq və əvəz etmək daha təhlükəsiz və təkrarlanandır

### 5. Dayanıqlılığın təşkili
**Nədir:** Dayanıqlılıq təkcə texniki deyil — təşkilati, prosesləşdirmə və mədəniyyət məsələsidir.

**Necə işləyir:**
- Testləmə, code review, CI/CD boruşu (pipeline)
- Monitoring və alerting
- Post-incident təhlil və runbook-lar
- Blameless postmortem mədəniyyəti

## Əsas terminlər
- Reliability (dayanıqlılıq) — gözlənilən vəziyyətdə funksiyanı yerinə yetirmə
- Fault prevention (xətanın qarşısını almaq)
- Fault tolerance (xətaya davamlılıq)
- Fault removal (xətanı aradan qaldırmaq)
- Fault forecasting (xətanı təxmin etmək)
- Twelve-Factor App — 12 ilkli bulud tətbiqi metodologiyası
- Immutable infrastructure (dəyişməz infrastruktur)
- Pets vs Cattle — server idarəetmə fəlsəfəsi
- Graceful degradation (ləüzumi enerji itkisi) — xəta zamanı əsas funksiyanı saxlayış
- CI/CD — davamlı inteqrasiya və çatdırılma

## Praktik nəticə
Dayanıqlılıq təkcə "kod yazmaq" deyil — sistemin dizaynı, konfiqurasiyası, deployment və idarəetmə proseslərini əhatə edir. Twelve-Factor App bu kontekstdə standart çərçivə təqdim edir. Immutabl infrastruktur isə server idarəetməsini yenidən qurmaq və əvəz etmək prinsipi ilə sadəcə və təhlükəsiz edir.

## Mənbə
Pages: 183-205 (PDF səh. 183-205)
