# Chapter 4 — Шаблоны программирования облачных приложений

## Bu chapter nədən bəhs edir?

Bulud tətbiqləri üçün əsas idman (idiomatic) proqramlaşdırma şablonlarını izah edir: distribüed hesablamanın səhv fərzləri (fallacies), context paketi, dayanıqlılıq şablonları (circuit breaker, debounce, retry), rate limiting, massivləşdirmə və sinxronizasiya nümunələri. Bu chapter konseptual və praktiki kod nümunələri ilə bəzədilib.

## Əsas fikirlər

### 1. Distribüed hesablamanın səhv fərzləri (Fallacies of Distributed Computing)
**Nədir:** 1991-ci ildə L Peter Deutsch tərəfindən formulə edilmiş, paylanmış sistemlərlə işləyən proqramçıların tez-tez etdiyi səhv fərzlər siyahısı.

**Səhv fərzlər:**
- Şəbəkə etibarlıdır (network is reliable) — kommutatorlar və marsrutizatorlar sınağa bilir
- Gecikmə sıfırdır (latency is zero) — şəbəkə üzərindən data ötürmə vaxt tələb edir
- Kanal tutumu sonsuzdur (bandwidth is infinite) — şəbəkə məhdud miqdarda data emal edə bilir
- Şəbəkə təhlükəsizdir (network is secure) —敏感 data açıq mətnlə ötürülməməlidir
- Topologiya dəyişməzdir (topology doesn't change) — serverlər və xidmətlər yaranır və yox olur
- Yalnız bir administrator var (single admin) — bir neçə administrator ziddiyyətli qərarlar qəbul edə bilər
- Daşıma xərcləri sıfırdır (transport cost is zero) — data ötürmə vaxt və pul tələb edir
- Şəbəkə homojendir (network is homogeneous) — hər şəbəkə birindən fərqlənir

**Nəyə lazımdır:** Bulud sistemlərini düzgün dizayn etmək üçün bu real məhdudiyyətləri başa düşmək.

### 2. Context paketi (Go 1.7+)
**Nədir:** Proseslər arasında ləğv edilmə sinxalları, son tarixlər və sorğu kontekstində dəyərlər ötürmək üçün idman (idiomatic) vasitə.

**Necə işləyir:**
- `context.Background()` — əsas kontekst, heç vaxt ləğv edilmir, `main` və testlərdə istifadə olunur
- `context.TODO()` — keçici olaraq istifadə edilən boş kontekst
- `context.WithTimeout(ctx, duration)` — vaxt aşımı ilə ləğv edilən kontekst
- `context.WithDeadline(ctx, time)` — konkret vaxtda ləğv edilən kontekst
- `context.WithCancel(ctx)` — manual ləğv edilən kontekst
- `context.WithValue(ctx, key, val)` — kontekstə açar/dəyər cütü əlavə etmək (spesifik sorğu dəyərləri üçün)

**Məqam:** Kontekst ləğv edildikdə, bütün törəmə kontekstlər də ləğv edilir, amma valideyn kontekstlər deyil. Kontekstlər thread-safe-dir, bir neçə eyni vaxtda işləyən gorutindərdə təhlükəsiz istifadə edilə bilər.

**Kitabdan kod nümunəsi:**
```go
func Stream(ctx context.Context, out chan<- Value) error {
    dctx, cancel := context.WithTimeout(ctx, time.Second*10)
    defer cancel()
    res, err := SlowOperation(dctx)
    if err != nil {
        return err
    }
    for {
        select {
        case out <- res:
        case <-ctx.Done():
            return ctx.Err()
        }
    }
}
```

### 3. Circuit Breaker (Devir Açıcı)
**Nədir:** Asılı xidmətlərdə xətaları avtomatik aşkar edib "devri açaraq" — müvəqqəti olaraq sorğuları dayandıraraq kaskad xətaları qarşısı alan şablon.

**Necə işləyir:** Elektrik devir açıcısı kimi iki vəziyyət var: bağlı (closed) və açıq (open). Bağlı vəziyyətdə sorğular adi şəkildə ötürülür. Açıq vəziyyətdə xidmətə sorğu göndərilmir, dərhal "service unreachable" xətası qaytarılır. Bir-birinin ardınca gələn xətalar sayı həddi aşdıqda devir açılır. Bəzi realizasiyalarda avtomatik bağlanma mexanizmi də var — üstünlük mərhələsi ilə artan gecikmə ilə.

**Kitabdan kod nümunəsi:**
```go
type Circuit func(context.Context) (string, error)

func Breaker(circuit Circuit, failureThreshold uint) Circuit {
    var consecutiveFailures int = 0
    var lastAttempt = time.Now()
    var m sync.RWMutex
    return func(ctx context.Context) (string, error) {
        m.RLock()
        d := consecutiveFailures - int(failureThreshold)
        if d >= 0 {
            shouldRetryAt := lastAttempt.Add(time.Second * 2 << d)
            if !time.Now().After(shouldRetryAt) {
                m.RUnlock()
                return "", errors.New("service unreachable")
            }
        }
        m.RUnlock()
        response, err := circuit(ctx)
        m.Lock()
        defer m.Unlock()
        lastAttempt = time.Now()
        if err != nil {
            consecutiveFailures++
            return response, err
        }
        consecutiveFailures = 0
        return response, nil
    }
}
```

### 4. Debounce (Anti-flash)
**Nədir:** Funksiyanın çağırma tezliyini məhdudlaşdıran şablon — vaxt intervalı daxilində gələn çağırışlardan yalnız birincini (DebounceFirst) və ya sonuncunu (DebounceLast) icra edir.

**Necə işləyir:**
- **DebounceFirst:** Mütləq məxaric (mutex) ilə qorunur. Hər çağırışda threshold (son çağırma vaxtı + interval) yenilənir. Əgər növbəti çağırış threshold-dan əvvəl gəlsə, kəşflənmiş cavab qaytarılır. İlk çağırış həmişə icra olunur.
- **DebounceLast:** `time.Ticker` və `sync.Once` istifadə edir. Son çağırışdan sonra interval qeder gözləyir, sonra funksiyanı icra edir.

**Kitabdan kod nümunəsi (DebounceFirst):**
```go
func DebounceFirst(circuit Circuit, d time.Duration) Circuit {
    var threshold time.Time
    var result string
    var err error
    var m sync.Mutex
    return func(ctx context.Context) (string, error) {
        m.Lock()
        defer func() {
            threshold = time.Now().Add(d)
            m.Unlock()
        }()
        if time.Now().Before(threshold) {
            return result, err
        }
        result, err = circuit(ctx)
        return result, err
    }
}
```

### 5. Retry (Təkrar Cəhd)
**Nədir:** Paylanmış sistemlərdə müvəqqəti xətaların təbii olaraq öz-özlüyündə yox olması ehtimalı nəzərə alınaraq, uğursuz əməliyyatları mərhələli olaraq təkrar edən şablon.

**Necə işləyir:** Funksiya (Effector) üçün təkrar cəhdləri idarə edən bağlamanı (closure) qaytarır. Parametrlər: maksimum təkrar sayı (retries) və gecikmə (delay). Hər cəhd arasında delay qədər gözləyir. Normal şəraitdə eksponensial artım (backoff) algoritmi ilə birlikdə istifadə olunur.

**Kitabdan kod nümunəsi:**
```go
type Effector func(context.Context) (string, error)

func Retry(effector Effector, retries int, delay time.Duration) Effector {
    return func(ctx context.Context) (string, error) {
        for i := 0; i < retries; i++ {
            response, err := effector(ctx)
            if err == nil {
                return response, nil
            }
            time.Sleep(delay)
        }
        return effector(ctx)
    }
}
```

### 6. Rate Limiting (Sürət Məhdudlaşdırma) və Throttle
**Nədir:** Xidmətə sorğu axınını nəzarət edən mexanizmlər. Debounce'dan fərqli olaraq throttle sadəcə vaxta görə sıxlığı məhdudlaşdırır.

**Necə işləyir:** Token bucket və ya leaky bucket alqoritmləri ilə hər bir müştəridən qəbul ediləcək sorğu sayını məhdudlaşdırır. Bu, həddindən artıq yük və təkrar cəhd fəlakətlərini (retry storms) qarşısı alır.

### 7. Massivləşdirmə və Sinxronizasiya nümunələri
**Nədir:** Birdən çox gorutinin eyni anda paylaşılan resurslara (map, counter və s.) daxil olması vəziyyətində yaradılan yarış vəziyyətlərini (race conditions) və kilidləşməni (lock contention) azaldan texnikalar.

**Necə işləyir:**
- **Vertical Sharding:** Assosiativ massivi açarın hash-i ilə N alt-massivə bölmək — hər gorutin öz alt-massivinə yazır, ümumi mütləq məxarici (mutex) ehtiyacı azalır.
- `sync.RWMutex`: Çox oxu, az yazı senaryoları üçün — oxuma üçün `RLock()`, yazma üçün `Lock()`.
- `sync.Once`: Funksiyanı mütləq bir dəfə icra etmək üçün istifadə olunur.
- `time.Ticker`: Periyodik yoxlamalar üçün, bitirildikdə `ticker.Stop()` çağırılmalıdır, yoxsa yaddaş sızıntısı olur.

## Əsas terminlər
- Fallacies of Distributed Computing (paylanmış hesablamanın səhv fərzləri)
- Context — ləğv, timeout və dəyər ötürmə mexanizmi
- Circuit Breaker (devir açıcı) — xətaları aşkar edib dayandırma şablonu
- Debounce (anti-flash) — çağırma tezliyini məhdudlaşdırma
- Retry (təkrar cəhd) — uğursuz əməliyyatların yenidən cəhdi
- Rate Limiting / Throttle — sorğu axınını nəzarət mexanizmi
- Vertical Sharding — hash əsaslı massiv bölmə
- sync.RWMutex — oxu/yazı kilidi
- sync.Once — tək icra garantiyası
- time.Ticker — periyodik taymer

## Praktik nəticə
Bu chapter paylanmış bulud sistemlərində ən çox rast gəlinən problemlərin — şəbəkə etibarsızlığı, xətalar, kaskad nasazlıqlar — həlli üçün şablonlar təqdim edir. Context paketi Go-da sorğu həyat dövrünü idarə etmək üçün standart vasitədir. Circuit Breaker, Debounce və Retry şablonları closure və sync paketi ilə həyata keçirilir — bu şablonlar real tətbiqlərdə xidmət asılılıqlarını idarə etmək üçün əsasdir.

## Mənbə
Pages: 88-124 (PDF səh. 88-124)
