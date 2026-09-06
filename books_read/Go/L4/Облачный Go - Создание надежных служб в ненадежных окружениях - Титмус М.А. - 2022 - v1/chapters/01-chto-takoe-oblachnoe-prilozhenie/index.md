# Chapter 1 — Что такое «облачное» приложение?

## Bu chapter nədən bəhs edir?

Bulud arxitekturasının əsas atributlarını — miqyaslana bilənlik (scalability), zəif əlaqəlilik (loose coupling), dayanıqlıq (resilience), idarəetmə (manageability) və gözlənilirlik (observability) — izah edir. Bu chapter konseptualdir, kod nümunəsi yoxdur.

## Əsas fikirlər

### 1. Bulud tətbiqinin tarixi və inkişafı
**Nədir:** 1950-ci illərdəki mainframe-lardan başlayaraq 1980-ci illərdəki şəxsi kompyuterlərə, 2006-cı ildə AWS-işə salınmasına qədər olan arxitekturaların təkamülü.

**Necə işləyir:** Hər yeni mərhələ məhdudiyyətləri aradan qaldırır və miqyaslama imkanı yaradır. Mainframe → PC (multi-tier) → Cloud (IaaS/PaaS/SaaS).

**Nəyə lazımdır:** Bugünkü "bulud" termininin mənşəyini başa düşmək, onun yeni bir şey deyil, təkamül nəticəsi olduğunu bilmək.

### 2. CNCF tərəfindən bulud tətbiqinin tərifi
**Nədir:** Cloud Native Computing Foundation tərəfindən verilən rəsmi tərif: bulud texnologiyaları təşkilatların miqyaslana bilən tətbiqlərini müasir dinamik mühitlərdə (ictimai, xüsusi, hibrid bulud) yaradıb işə sala bilməsinə imkan verir.

**Necə işləyir:** Zəif bağlı sistemlər (loose coupling), dayanıqlılıq (resilience), idarəetmə (manageability) və gözlənilirlik (observability) birləşdilir. Avtomatlaşdırma ilə birlikdə engineer-lər tez və təxmin edilə bilən dəyişikliklər edə bilir.

### 3. Miqyaslana bilənlik (Scalability)
**Nədir:** Sistemin artan yükə uyğunlaşaraq performansını saxlayabilməsi.

**Necə işləyir:** iki strategiya var:
- **Vertikal miqyaslama (Vertical scaling):** Eyni serverdə RAM/CPU əlavə etmək. Sadədir, amma limitli.
- **Horizontal miqyaslama (Horizontal scaling):** Yeni server instance-ları əlavə etmək (load balancer arxasında). Limitsizdir, amma idarəsi daha çətindir.

### 4. Zəif əlaqəlilik (Loose Coupling)
**Nədir:** Komponentlərin bir-birindən minimal asılılıqla qurulması; bir komponent dəyişsə, digəri təsir olmadan işləməyə davam edir.

**Necə işləyir:** Standart protokollar (HTTP, gRPC) və müəyyən müqavilələr (contracts) vasitəsilə əlaqə. Nümunə: veb-brauzer və veb-server standart HTTP protokolu üzərində zəif bağlıdır — server yenilənsə brauzer yenilənməyə ehtiyac yoxdur.

### 5. Dayanıqlılıq (Resilience)
**Nədir:** Sistemin xətalar (failures) vəziyyətində də funksiyasını yerinə yetirməsi vəziyyəti.

**Necə işləyir:** Hər komponentin xəta verə biləcəyini qəbul edərək dizayn edilir. Xəta alt sistemdə baş verərsə, uyğun olaraq izolyasiya edilir və kaskad xətalar (cascading failures) qarşısı alınır.

### 6. İdarəetmə (Manageability)
**Nədir:** Sistemin davranışını kod dəyişdirmədən dəyişdirmə imkanı.

**Necə işləyir:** Konfiqurasiya xarici mənbələrdə (environment variable, config file, flag) saxlanır. Məsələn, URL hardcode edilmir, config faylından oxunur.

### 7. Gözlənilirlik (Observability)
**Nədir:** Sistemin daxili vəziyyətini xarici müşahidələrdən (logs, metrics, traces) anlama imkanı.

**Necə işləyir:** Sadəcə logging və dashboard əlavə etmək kifayət deyil — sistem real şəraitə hazırlanmalıdır, gələcəkdə bilinməyən suallara cavab verə bilməlidir.

## Əsas terminlər
- Scalability (miqyaslana bilənlik)
- Loose Coupling (zəif əlaqəlilik)
- Resilience (dayanıqlılıq / xətalara davamlılıq)
- Manageability (idarəetmə)
- Observability (gözlənilirlik)
- Vertical Scaling (vertikal miqyaslama)
- Horizontal Scaling (horizontal miqyaslama)
- Cascading Failures (kaskad xətalar)
- CNCF (Cloud Native Computing Foundation)

## Praktik nəticə
Bulud tətbiqi yalnız "server buluda" demək deyil — bu miqyaslana bilən, zəif bağlı, dayanıqlı, idarə edilə bilən və gözlənilən sistem dizayn məntiqidir.

## Mənbə
Pages: 25-35 (PDF səh. 25-35)
