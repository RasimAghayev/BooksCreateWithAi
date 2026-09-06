# Облачный Go / Создание надежных служб в ненадежных окружениях — Xülasə (Azərbaycanca)

> **Kitab:** Облачный Go / Создание надежных служб в ненадежных окружениях — Мэтью А. Титмус (Matthew A. Titmus), ДМК Пресс, 2022 (ISBN 978-5-97060-965-1, 419 səh., 12 chapter + əlavələr)
> **Səviyyə:** 🚀 Advanced (4/5) · **Dil:** Go, Cloud Native
> **Struktur:** 12 konseptual/praktik chapter — cloud native arxitektura, patterns, tətbiq qurulumu, Docker, reliability, scalability, loose coupling, resilience, manageability, observability.

---

## Kitabın məqsədi və yanaşması

Müəllif bulud (cloud native) sistemlərin ən əsas atributlarını — miqyaslana bilənlik, zəif əlaqəlilik, dayanıqlılıq, idarəetmə, gözlənilirlik — izah edir və hər birini Go dilində real şablonlar (patterns) ilə nümayiş etdirir. Kitab "Design Patterns" kitabının strukturu ilə paralel gedir: hər şablon üçün Applicability, Components, Implementation, Code Example bölmələri var. Praktik tərəf əsasdır — chapter 5-də tam bir key/value xidməti qurulur, chapter 11-də isə OpenTelemetry ilə gözlənilirlik (observability) qurulur.

---

## Chapter-by-chapter xülasə

### Ch1 — Что такое «облачное» приложение?
Bulud arxitekturasının 5 əsas atributunu izah edir: scalability, loose coupling, resilience, manageability, observability. Konseptualdir, kod yoxdur. CNCF tərifi, mainframe → cloud təkamülü, hər atributun nəzəri izahı.

### Ch2 — Почему Go правит облачным миром?
Go dilinin bulud hesablamaları üçün uyğunluğunu izah edir. CSP modeli, goroutine, kanallar, kompozisiya irsiyyət əvəzinə, sürətli kompilasiya, yaddaş təhlükəsizliyi, statik bağlantı, sabit API. Konseptualdir, kod yoxdur.

### Ch3 — Основы языка Go
Go sintaksis və əsas konstruksiyalar: məlumat tipləri, dəyişənlər, konstantlar, konteynerlər (massivlər, kəsiklər/slice, xəritələr/map), göstəricilər, idarəetmə strukturları, xətaların emalı, funksiyalar, strukturlar, metodlar, interfeyslər, konkurensiya. 40 səhifəlik praktik chapter, hər anlayışa kod nümunəsi.

### Ch4 — Шаблоны программирования облачных приложений
Distribüed hesablamanın səhv fərzləri (fallacies), context paketi, dayanıqlılıq şablonları: Circuit Breaker, Debounce (First/Last), Retry, Rate Limiting, Throttle. Vertical Sharding və sync.RWMutex ilə kilidləşmə azaltma.

### Ch5 — Конструирование облачной службы
Chapter 4 şablonlarını real tətbiqdə göstərir: `net/http` + `gorilla/mux` ilə RESTful key/value xidməti, əməliyyat qeydiyyatı (transaction logger), fayl əsaslı və PostgreSQL əsaslı realizasiyalar, Docker konteynerizasiyası (multi-stage build).

### Ch6 — Все дело в надежности
Dayanıqlılıq (reliability) anlayışı: Jean-Claude Laprie modeli (fault prevention, tolerance, removal, forecasting), Twelve-Factor App (12 ilk), immutabl infrastruktur (cattle, not pets). Konseptualdir.

### Ch7 — Масштабируемость
Stateless (vəziyyətsiz) dizaynın üstünlükləri, yaddaş sızıntıları (goroutine, ticker, map-lər), gorutin əsaslı polling nümunələri, monolith vs microservices arxitekturaları.

### Ch8 — Слабая связанность
Protokol çeşidləri (REST, gRPC, SOAP, pub/sub), vaxt əlaqəliliyi, plugin sistemləri (Go `plugin`, HashiCorp `go-plugin`), Heksagonal Arxitektura (ports and adapters). `gorilla/mux` marşrutlaşdırma, `.proto` ilə gRPC kontraktları.

### Ch9 — Устойчивость
Kaskad xətalar, təkrar cəhd fəlakətləri (retry storms), eksponensial artım + jitter, rate limiting, load shedding, bulkhead, xidmət redundansı, avtomatik miqyaslama, health checks (liveness, readiness, deep).

### Ch10 — Управляемость
Konfiqurasiya iyerarxiyası (defaults → config → env → flags), dinamik yeniləmə (Viper `WatchConfig`), xüsusiyyət flaqları (feature flags — statik və dinamik), Cobra CLI, Viper konfiqurasiya kitabxanası.

### Ch11 — Наблюдаемость
Gözlənilirlik (observability) üç sütun: logs, metrics, traces. OpenTelemetry standartı, strukturlaşdırılmış jurnalçılık (Zap), Prometheus metrikaları, avtomatik instrumentasiya (`net/http`, gRPC), span events və labels, dinamik nümunə və kardinallıq riski.

---

## Kitabın əsas mesajları

1. **Bulud tətbiqi yalnız "server buluda" deyil** — bu arxitektura məntiqidir: miqyaslana bilən, zəif bağlı, dayanıqlı, idarə edilə bilən və gözlənilən sistem.
2. **Context paketi Go-da sorğu həyat dövrünü idarə etmək üçün standart vasitədir** — ləğv, timeout, dəyər ötürmə.
3. **Shablonlar closure və sync paketi ilə həyata keçirilir** — Circuit Breaker, Debounce, Retry — hamısı `func() Circuit` qaytarır.
4. **Dayanıqlılıq (resilience) təkcə texniki deyil** — təşkilati, prosesləşdirmə və mədəniyyət məsələsidir.
5. **İdarəetmə (manageability) və gözlənilirlik (observability) cloud native-nin iki böyük sütunudur** — konfiqurasiya və loglar/metrikalar/izlər sistemi.
6. **Go statik binary və konteyner dostudur** — ~2 MB "Hello World", heç bir xarici runtime tələb etmir.

---

## Kitabdan sonra

- **Advanced Go Concurrency:** "Concurrency in Go" (Katherine Cox-Buday), "Learning Go" (Jon Bodner)
- **Distributed Systems:** "Designing Distributed Systems" (Brendan Burns), "Microservices Patterns" (Chris Richardson)
- **Observability:** OpenTelemetry in Action, Prometheus Up & Running
- **Cloud Native:** CNCF Landscape, Kubernetes in Action
- **Resilience:** "Release It!" (Michael T. Nygard), "Site Reliability Engineering" (Google)
