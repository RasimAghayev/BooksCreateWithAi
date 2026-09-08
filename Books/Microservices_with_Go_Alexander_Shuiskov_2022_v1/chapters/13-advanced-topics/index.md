# Chapter 13 — Advanced Topics (səh. 267-290)

## Bu fəsil nədən bəhs edir?

pprof profiling, Grafana dashboard-lar, framework-lər (IoC), ownership
metadata və JWT ilə təhlükəsizlik — kitabın bağlayıcı fəsli.

## Əsas fikirlər

### 1. Profiling (pprof)
- Real-time performans datası: CPU istifadəsi, heap allocation, call graph
- `net/http/pprof` handler-ı :6060-da dinlənir:
```go
go http.ListenAndServe("localhost:6060", nil)  // debug/pprof
// CPU:  go tool pprof http://localhost:6060/debug/pprof/profile?seconds=5
// Heap: go tool pprof http://localhost:6060/debug/pprof/heap
// pprof> web  (Graphviz ilə qraf açılır)
```
- **CPU profile:** node = paket/funksiya + keçən vaxt; heavyOperation 0.01s
  özü, amma daxil çağırışlar 4.39s (md5.Write 2.78s + rand.Read 1.59s) —
  böyük düzbucaqlılar = ən çox CPU yeyən
- **Heap profile:** funksiyanın heap-dəki payı; api.serviceRegister,
  zap.NewProduction, trace.init kimi başlanğıc alloc-ları görünür
- Əlavə: `go test -cpuprofile/-memprofile`; pprof `top`/`top10`; **goroutine**
  profili (stack trace-lərlə)

### 2. Dashboards (Grafana)
- Dashboard = metric qrafikləri toplusu; debugging + data korrelyasiyası
- Hər servis üçün + sistem-global dashboard-lər (aktiv instans, throughput)
- Quraşdırma: `docker run -d -p 3000:3000 grafana/grafana-oss` (admin/admin)
  → Data source: Prometheus, URL `http://host.docker.internal:9090` → New
  dashboard → Add panel → metric (process_open_fds, go_gc_duration_seconds)
  → Run queries → Apply
- **Panel metrikləri (Golden Signals + ):** client/server error rate, API
  throughput, latency **p90/p95/p99 percentile**-lərlə, CPU%, memory,
  network throughput; servisə xas əlavələr (Kafka, cache, DB)

### 3. Framework-lər və IoC
- Konvensiyalar (naming/yerləşdirmə) ↔ framework-lər (məcburi struktur)
- **Inversion of Control (IoC):** net/http ListenAndServe çağırılanda
  idarəni paket alır — handler funksiyalarınızı avtomatik çağırır
- İstifadə: web serverlər (Thrift/gRPC), async emal (Kafka handler-lər)
- **Çatışmazlıqlar:** debug çətin (gizli arxa plan işi), dik öyrənmə
  əyrisi, `reflect` dynamic çağırışları static yoxlamaları zəiflədir
- Tövsiyə: sadə variantdan başla; framework yalnız fayda > ziyan olanda

### 4. Service Ownership (servis sahibliyi)
Minlərlə servisdə "kim məsuldur?" sualı — vulnerability paylandıqda həlledici.

3 aspekt:
- **Accountability (məsuliyyət):** engineering manager başına domain
  (shared/team-based məsuliyyət bulanıq — team anlayışı dəyişkəndir)
- **Support:** Slack kanalı / Jira biletləri URL
- **On-call:** PagerDuty rotation ID xəritələnməsi

```yaml
ownership:
  rating-service:
    accountable: example@somecompany.com
    support: {slack: rating-service-support-group}
    oncall: {pagerduty_rotation: SOME_ROTATION_ID}
```
- API ilə sorğulana bilən sistem + ownership MECBURİ (service creation
  prosesində tələb et)

### 5. JWT (JSON Web Token)
- **Authentication** (kimliy yoxlama — login) → token; **Authorization**
  (hüquq yoxlama — role)
- JWT = header.payload.signature (Base64url):
  - **Payload (claims):** `{"name":"Alexander","role":"admin","iat":...}`
  - **Header:** `{"alg":"HS256","typ":"JWT"}`
  - **Signature:** `HMACSHA256(base64(header) + "." + base64(payload), secret)`
    — secret yalnız autentifikasiya serverində; dəyişiklik = imza pozulur

**Server (autentifikasiya) — golang-jwt/jwt:**
```go
token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
    "username": username, "iat": time.Now().Unix()})
tokenString, _ := token.SignedString(secret)   // 200 + token / 401
```

**Client (authorization):**
```go
req.Header.Set("Authorization", "Bearer "+token)
```

**Server (yoxlama):**
```go
token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
    if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
        return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
    }
    return secret, nil
})
claims := token.Claims.(jwt.MapClaims)  // valid → etibarlı payload
```
- Alqoritm yoxlaması MÜTLƏQ (alg dəyişdirilərək hücumların qarşısını alır)

**Best practices:**
1. **exp** (bitmə vaxtı) — admin üçün bir neçə saat
2. **iat** (verilmə vaxtı) — breach olsa, həmin andan əvvəlki tokenləri etibarsızlaşdır
3. HTTPS üzərindən — MITM attack header-ləri oğurlaya bilər
4. Standard claim sahələrini üstün tut (jwt.io siyahısı)

## Termindirmə (AZ)
- Profiling — Profilləşdirmə (performans ölçməsi)
- Heap — Yığın (dinamik yaddaş)
- Call Graph — Çağırış Qrafı
- Inversion of Control (IoC) — İdarənin Tərsinə Çevrilməsi
- Ownership — Sahiblik (məsuliyyət)
- Accountability — Məsuliyyət (hesabatlılıq)
- Authentication — Kimliyin Təsdiqi
- Authorization — Səlahiyyətləndirmə
- Claim — İddia (JWT payload sahəsi)
- Bearer Token — Daşıyıcı Tokeni
- Percentile (p90/p95/p99) — Persentil

## Kviz sualları
1. JWT-nin 3 hissəsi hansılardır? (Header — alqoritm; Payload — claims;
   Signature — HMAC(secret))
2. Niyə JWT serverdə alqoritm yoxlanılır? (alg sahəsi dəyişdirilib zəif
   imza ilə hücum oluna bilər — HS256 gözlənilir)
3. IoC nə deməkdir? (Framework icra axınını özünə götürür — sizin
   handler-lərinizi çağırır)
4. Ownership metadata-da accountable kimi nə saxlanılır? (Engineering
   manager — insan, team YOX)
