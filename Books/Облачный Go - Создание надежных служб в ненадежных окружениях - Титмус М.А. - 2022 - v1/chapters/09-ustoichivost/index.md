# Chapter 9 — Устойчивость

## Bu chapter nədən bəhs edir?

Sistemin nöqsanlar (failures) vəziyyətində də funksiyasını yerinə yetirmə qabiliyyəti — uyğun gəlmə (resilience) — izah edir. Kaskad xətalar (cascading failures), təkrar cəhd fəlakətləri (retry storms), eksponensial artım və qarışıq (jitter) algoritmləri, rate limiting, load shedding, bulkhead şablonu, xidmət redundansı və avtomatik miqyaslama, sağlamlıq yoxlamaları (health checks) müzakirə edilir.

## Əsas fikirlər

### 1. Uyğun gəlmə (Resilience) və xəta növləri
**Nədir:** Sistemin əsas funksiyasını itirmədən xətalar vəziyyətində davranma qabiliyyəti.

**Necə işləyir:**
- **Kaskad xətalar (cascading failures):** Bir komponent düşdükdə ondan asılı digər komponentlər də düşür — domino effekti
- **Təkrar cəhd fəlakəti (retry storm):** Bir neçə xidmət eyni vaxtda uğursuz asılılığı təkrar cəhd edərkə yükü daha da artırır — Amazon DynamoDB hadisəsi misal olaraq göstərilir
- **Qorxulu döngü (death spiral):** Xidmət yenidən başlananda yük artır, yenidən düşür — bu döngü təkrarlanır

### 2. Təkrar cəhd strategiyası (Retry with Backoff and Jitter)
**Nədir:** Uğursuz əməliyyatları təkrar cəhd edərkən yükü azaldan və təkrarlanma effektini minimuma endirən strategiya.

**Necə işləyir:**
- **Sabit gecikmə (fixed delay):** Hər cəhddən sonra sabit müddət gözlə — sadə, lakin bütün client-lər eyni anda təkrar cəhd edərsə, kütləvi tələbat (thundering herd) yaranır
- **Eksponensial artım (exponential backoff):** Hər cəhd arasında gecikmə iki qat artır — 1s, 2s, 4s, 8s...
- **Qarışıq (Jitter):** Eksponensial artıma təsadüfi miqdar əlavə edir — client-lərin eyni anda təkrar cəhd etməsini qarşısı alır

**Alqoritm nümunəsi:**
```
delay = base_delay * 2^attempt + random_jitter
```

### 3. Rate Limiting və Load Shedding
**Nədir:** Xidməti aşırı yükdən qoruyan mexanizmlər.

**Necə işləyir:**
- **Rate Limiting:** Müəyyən vaxt intervalında qəbul ediləcək sorğu sayını məhdudlaşdırır — token bucket və ya leaky bucket alqoritmləri
- **Load Shedding:** Xidmət öz-özünü qoruyur — yük həddini aşdıqda bəzi sorğuları rədd edir ( rejecting requests ), əsas funksiyanı saxlayır
- Bu, "yarmadan çəkmək" (fail-fast) prinsipi ilə birlikdə işləyir

### 4. Bulkhead (Gertməli Bölmə)
**Nədir:** Gəmi gertmələri (bulkhead) kimi, xidmətin müxtəlif resursunu izolyasiya edən şablon.

**Necə işləyir:**
- Hər asılı xidmət üçün ayrı resurs kütləsi (pool) təyin edilir — məsələn, DB connection pool
- Bir xidmət düşərsə, digərlərinin resursları təsir olmaz
- Bu, kaskad xətaların yayılmasını məhdudlaşdırır

### 5. Xidmət redundansı və avtomatik miqyaslama
**Nədir:** Bir neçə instance işlətmək və yükü avtomatik olaraq bölüşdürmək.

**Necə işləyir:**
- **Redundans:** Birdən çox instance işlədir — biri düşsə, digərləri davam edir
- **Avtomatik miqyaslama:** Yüksək yükdə yeni instance-lar əlavə edilir, aşağı yükdə azaldılır
- Nöqsanların ehtimalı: 2 instance 99.9% yüksəkliyi təmin edir, 3 instance 99.999% (5-nine)

### 6. Sağlamlıq yoxlamaları (Health Checks)
**Nədir:** Xidmətin vəziyyətini izləmək və idarəetmə sistemlərinə (orchestrator) məlumat vermək üçün endpoint-lər.

**Necə işləyir:**
- **Liveness Probe (canliliq yoxlaması):** Proses hələ də yaşayır? Əgər deyilsə, yenidən başlat.
- **Readiness Probe (hazırlıq yoxlaması):** Sorğuları qəbul edə bilirmi? Əgər deyilsə, yükü çək.
- **Deep Probe (dərin yoxlama):** Sadəcə proses yox, həm də asılılıqlar (DB, cache) çalışır?

**Kitabdan kod nümunəsi (health endpoint):**
```go
func healthHandler(w http.ResponseWriter, r *http.Request) {
    // DB və ya digər asılılıqları yoxla
    if err := db.Ping(); err != nil {
        http.Error(w, "db unavailable", http.StatusServiceUnavailable)
        return
    }
    w.WriteHeader(http.StatusOK)
    fmt.Fprintln(w, "OK")
}
```

## Əsas terminlər
- Resilience (uyğun gəlmə / dayanıqlılıq) — xəta zamanında funksiyanı saxlayış
- Cascading failure (kaskad xəta) — domino effekti
- Retry storm (təkrar cəhd fəlakəti) — kütləvi təkrar cəhdlər
- Exponential backoff (eksponensial artım) — artan gecikmə ilə təkrar cəhd
- Jitter (qarışıq) — təsadüfi gecikmə, thundering herd qarşısı
- Rate limiting (sürət məhdudlaşdırma) — sorğu sayını məhdudlaşdırma
- Load shedding (yük atma) — həddi aşan sorğuları rədd etmə
- Bulkhead (gertməli bölmə) — izolyasiya pool
- Redundancy (redundans) — artıq nümunələr
- Health checks (sağlamlıq yoxlamaları) — liveness, readiness, deep
- Fail-fast (zərərli döngüdən çıxma) — tez uğursuzluq

## Praktik nəticə
Uyğun gəlmə (resilience) distribüed sistemlərin ən kritik xüsusiyyətidir. Kaskad xətalar və təkrar cəhd fəlakətləri real sistemlərdə tez-tez baş verir. Eksponensial artım + jitter ilə təkrar cəhdlər, rate limiting və load shedding ilə yükün idarə edilməsi, bulkhead ilə izolyasiya — bunların hamısı birgə işləyərək sistemin dayanıqlılığını artırır.

## Mənbə
Pages: 279-314 (PDF səh. 279-314)
