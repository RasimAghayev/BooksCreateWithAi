# Chapter 11 — Rate Limiting

## Bu fəsil nədən bəhs edir?

Token-bucket rate limiter (x/time/rate) ilə 429 Too Many Requests
middleware-i: qlobal limiter → IP-əsaslı per-client limiter → runtime
konfiqurasiyası + deaktivasiya.

## Əsas fikirlər

### 1. Token Bucket alqoritmi
**Nədir:** Bucket ölçüsü `b`, başlanğıc dolu, `r` token/saniyə doldurulur.
Hər HTTP sorğusu 1 token götürür; bucket boşdursa → 429.

**Davranış:** burst = `b` sorğu ardıcıllıqla; uzunmüddətdə orta = `r`
sorğu/saniyə. Nümunə: `rate.NewLimiter(2, 4)` → 4-lük burst + 2 rps orta.

```go
limiter := rate.NewLimiter(2, 4) // Limit(float64) alias
```
- `Allow()` → 1 token istehlak edir; mutex-lə qorunur, concurrency-safe

**Middleware pattern məqamı:** "initialization code" (limiter yaratmaq)
yalnız bir dəfə — wrap zamanında işləyir; closure limiter-i "yaxında saxlayır".

### 2. Global rate limiter
**Kitabdan kod nümunəsi:**
```go
func (app *application) rateLimit(next http.Handler) http.Handler {
    limiter := rate.NewLimiter(2, 4)
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if !limiter.Allow() {
            app.rateLimitExceededResponse(w, r)
            return
        }
        next.ServeHTTP(w, r)
    })
}

// errors.go
func (app *application) rateLimitExceededResponse(w http.ResponseWriter, r *http.Request) {
    app.errorResponse(w, r, http.StatusTooManyRequests, "rate limit exceeded")
}
```

**Zəncir sırası:**
```go
return app.recoverPanic(app.rateLimit(router))
```
- recoverPanic ƏVVƏLDƏ (rateLimit-in öz panic-i də tutulsun), rateLimit isə
  router-dən əvvəl — bahalı işə (JSON decode, DB sorğu) getməzdən əvvəl
  reject etsin

### 3. IP-based per-client rate limiting
**Nədir:** Hər client üçün ayrı limiter — map[ip]*rate.Limiter; mutex ilə
qorunan map.

**Kitabdan kod nümunəsi (cleanup ilə tam versiya):**
```go
func (app *application) rateLimit(next http.Handler) http.Handler {
    type client struct {
        limiter  *rate.Limiter
        lastSeen time.Time
    }
    var (
        mu      sync.Mutex
        clients = make(map[string]*client)
    )

    // Hər dəqiqə 3+ dəqiqə görünməyən client-ləri sil
    go func() {
        for {
            time.Sleep(time.Minute)
            mu.Lock()
            for ip, client := range clients {
                if time.Since(client.lastSeen) > 3*time.Minute {
                    delete(clients, ip)
                }
            }
            mu.Unlock()
        }
    }()

    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ip, _, err := net.SplitHostPort(r.RemoteAddr)
        if err != nil {
            app.serverErrorResponse(w, r, err)
            return
        }
        mu.Lock()
        if _, found := clients[ip]; !found {
            clients[ip] = &client{limiter: rate.NewLimiter(2, 4)}
        }
        clients[ip].lastSeen = time.Now()
        if !clients[ip].limiter.Allow() {
            mu.Unlock()
            app.rateLimitExceededResponse(w, r)
            return
        }
        mu.Unlock() // defer YOX! — yoxsa bütün downstream handler-lər bitənə
                     // qədər kilid saxlanılar
        next.ServeHTTP(w, r)
    })
}
```

**Sub-kod izahı:**
- `net.SplitHostPort(r.RemoteAddr)` → "IP:port" → IP ayrılır
- **Mutex-lərin critical bölgəsi**: map oxu/yazısı + Allow() çağırışı — bütün
  bölünməz (atomic) olmalıdır
- **`mu.Unlock()` defer-lə DEYİL** — defer işləsəydi kilid handler zənciri
  boyunca saxlanılardı (bütün下游 handler-lər bitənədək) — deadlock
  yox, amma ciddi bottleneck
- `lastSeen` + background cleanup goroutine → map-in sonsuz böyüməsinin
  (memory leak) qarşısı
- Cleanup goroutine-i middleware initialization-də (bir dəfə) start olunur

**Limit (Additional Info):** Bu pattern yalnız tək-maşın üçündür.
Distributed sistemdə (multi-server + LB): HAProxy/Nginx built-in
rate-limiting, yaxud mərkəzi Redis-də request count.

### 4. Runtime konfiqurasiya
**Kitabdan kod nümunəsi:**
```go
type config struct {
    port int
    env  string
    db   struct { /* ... */ }
    limiter struct {
        rps     float64
        burst   int
        enabled bool
    }
}

flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "Rate limiter maximum requests per second")
flag.IntVar(&cfg.limiter.burst, "limiter-burst", 4, "Rate limiter maximum burst")
flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")
```

**Middleware-də:**
```go
if app.config.limiter.enabled {
    // ... yoxla
}
next.ServeHTTP(w, r)
```

**Sub-kod izahı:**
- `-limiter-burst=2` → burst-i azalt; `-limiter-enabled=false` → load
  test/benchmark üçün tam deaktiv ( bütün sorğular keçir)
- `rate.Limit(app.config.limiter.rps)` → float64 → Limit tip konversiyası

## Bu fəsildə layihəyə nə əlavə olundu? (project continuity)

**Dəyişdirilən:**
- middleware.go → rateLimit (IP-əsaslı + cleanup goroutine + enabled yoxlaması)
- errors.go → rateLimitExceededResponse (429)
- routes.go → `recoverPanic(rateLimit(router))` zənciri
- main.go → config.limiter struct + 3 flag

**Yeni asılılıq:** `golang.org/x/time/rate`

## Əsas terminlər

- Token Bucket (jeton vedrəsi) — burst + refill rate modeli
- Burst (partlayış) — ardıcıllıqla icazə verilən maksimum sorğu
- 429 Too Many Requests — rate limit status kodu
- Closure (qapanış) — xarici dəyişəni "yaxınlaşdıran" funksiya (limiter)
- Critical Section (kritik bölgə) — mutex altında atomik icra olunmalı kod
- Memory Leak (yaddaş sızması) — sonsuz böyüyən map (cleanup qarşısını alır)
- Reverse Proxy (tərs proksi) — Nginx/HAProxy kimi ön qat; built-in
  rate-limit imkanı

## Praktik nəticə

Rate limiting API-nin ön qapı mühafizəsidir: token bucket rəqəmləri
(r=2, b=4 kimi) domain-ə görə kalibrlənir; IP-əsaslı ayırma "bir pis
client hamını çəkib düşürməsin" deyə. Ən vacib incəlik: mutex-i `defer`
ilə YOX, açılış bağlamadan əvvəl açmaq — əks halda bütün下游 handler-lər
kilid gözləyər. Distributed deploy-da Redis/LB həllərə keçilməlidir.

## Mənbə
Pages: 248-262 (raw 248-262)
