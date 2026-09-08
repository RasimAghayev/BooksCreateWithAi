# Chapter 18 — Cross Origin Requests (CORS)

## Bu fəsil nədən bəhs edir?

Same-origin policy, origin anlayışı, simple/preflight CORS sorğuları,
Access-Control header ailəsi, trusted origin safelist + reflection,
preflight OPTIONS interceptəsi və CORS təhlükəsizlik qaydaları.

## Əsas fikirlər

### 1. Origin nədir?
**Tərif:** scheme + host + port (göstərilibsə) eynidirsə — eyni origin.

| URL A | URL B | Eyni origin? | Səbəb |
|---|---|---|---|
| `https://foo.com/a` | `http://foo.com/a` | Xeyr | scheme fərqli |
| `http://foo.com/a` | `http://www.foo.com/a` | Xeyr | host fərqli |
| `http://foo.com/a` | `http://foo.com:443/a` | Xeyr | port fərqli |
| `http://foo.com/a` | `http://foo.com/b` | Bəli | yalnız path |
| `http://foo.com/a` | `http://foo.com/a?b=c` | Bəli | yalnız query |

### 2. Same-Origin Policy (brauzer təhlükəsizlik modeli)
**Nədir:** Brauzerin default qadağası:
- Başqa origin-dən **embed** (img, CSS, JS) → icazəli
- Başqa origin-ə **data göndərmək** (form submit) → icazəli (buna görə CSRF
  mümkündür — SameSite cookies/CSRF token lazımdır!)
- Başqa origin-dən **data OXUMAQ** → QADAĞA

**Vacib:** Sorğu gedir, server işləyir, 200 qayıdırır — brauzer yalnız
JavaScript-in RESPONSE-u GÖRMƏSİNİ bloklayır. Policy yalnız brauzer
məhdudiyyətidir — curl/wget kimi client-lər tam sərbəstdir.

**Brauzer avtomatik `Origin: http://localhost:9000` header-i qoyur** —
server üçün etibar siqnalı.

### 3. Simple CORS — Access-Control-Allow-Origin
**Nədir:** Cavaba `Access-Control-Allow-Origin` qoymaq brauzerə "oxumağa
icazə ver" deməkdir.

**Wildcard (hamıya açıq):**
```go
func (app *application) enableCORS(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        next.ServeHTTP(w, r)
    })
}
```

**Zəncir mövqeyi (vacib!):** enableCORS rateLimit-dən ƏVVƏL — yoxsa 429
cavabı Allow-Origin almayacaq → client 429-u yox, sadəcə bloku görəcək:
```go
return app.recoverPanic(app.enableCORS(app.rateLimit(app.authenticate(router))))
```

**Trusted origins (reflection pattern):**
**Problem:** Header-də YALNIZ BİR origin ola bilər (siyahı qadağandır;
browsers space-separated dəstəyi yoxdur). **Həll:** Origin request
header-ini safelist ilə müqayisə et → uyğun gələrsə EKO et.

**Kitabdan kod nümunəsi (konfiqurasiya):**
```go
// config
cors struct {
    trustedOrigins []string
}

// flag: space-separated origins
flag.Func("cors-trusted-origins", "Trusted CORS origins (space separated)", func(val string) error {
    cfg.cors.trustedOrigins = strings.Fields(val)
    return nil
})
// -cors-trusted-origins="https://www.example.com https://staging.example.com"
```

**Kitabdan kod nümunəsi (middleware):**
```go
func (app *application) enableCORS(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Add("Vary", "Origin")
        origin := r.Header.Get("Origin")
        if origin != "" {
            for i := range app.config.cors.trustedOrigins {
                if origin == app.config.cors.trustedOrigins[i] { // EXACT, case-sensitive
                    w.Header().Set("Access-Control-Allow-Origin", origin) // reflection
                    // preflight (aşağıda) ...
                    break
                }
            }
        }
        next.ServeHTTP(w, r)
    })
}
```

**Vary: Origin (mütləq!):** Cavab Origin-dən asılıdır (header dəyəri fərqli
və ya yoxdur) → cache-lər xəbərdar olmalıdır. **Qayda:** cavab hər hansı
request header-in məzmunundan asılıdırsa, o header adı Vary-da olmalıdır —
header sorğuda OLMASA BELƏ.

### 4. Preflight CORS sorğuları
**Nədir:** "Simple" olmayan cross-origin sorğulardan ƏVVƏL brauzerin
göndərdiyi OPTIONS yoxlaması.

**Simple sorğu şərtləri (hamısı):**
1. Metod: HEAD, GET və ya POST
2. Header-lər yalnız forbidden və ya 4 CORS-safe: Accept,
   Accept-Language, Content-Language, Content-Type
3. Content-Type (varsa): `application/x-www-form-urlencoded`,
   `multipart/form-data`, `text/plain`

**Nümunə:** `Content-Type: application/json` → simple DEYİL → preflight
tetiklenir!

**Preflight request header-ləri:**
- `Origin` → mənşə
- `Access-Control-Request-Method: POST` → real sorğunun metodu
- `Access-Control-Request-Headers: content-type` → real sorğunun
  non-safe header-ləri (yalnız bunlar; hamısı YOX)

**Preflight identifikasiyası (3 komponent hamısı):**
```
r.Method == OPTIONS && Origin != "" && Access-Control-Request-Method != ""
```

**Cavab header-ləri:**
```
Access-Control-Allow-Origin: <reflected origin>
Access-Control-Allow-Methods: OPTIONS, PUT, PATCH, DELETE
Access-Control-Allow-Headers: Authorization, Content-Type
```
- Safe metodları (HEAD/GET/POST) və safe header-ləri yazmaq lazım deyil

### 5. Preflight interceptəsi (tam enableCORS)
**Kitabdan kod nümunəsi:**
```go
func (app *application) enableCORS(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Add("Vary", "Origin")
        w.Header().Add("Vary", "Access-Control-Request-Method")
        origin := r.Header.Get("Origin")
        if origin != "" {
            for i := range app.config.cors.trustedOrigins {
                if origin == app.config.cors.trustedOrigins[i] {
                    w.Header().Set("Access-Control-Allow-Origin", origin)
                    if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
                        w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, PUT, PATCH, DELETE")
                        w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
                        w.WriteHeader(http.StatusOK) // 204 YOX — köhnə brauzerlər!
                        return
                    }
                    break
                }
            }
        }
        next.ServeHTTP(w, r)
    })
}
```

**Sub-kod izahı:**
- `Vary: Access-Control-Request-Method` → preflight/normal fərqli cavab
- **200 OK (204 YOX)** — bəzi brauzer versiyaları 204-ü dəstəkləmir və
  real sorğunu bloklayır
- Preflight cavabında body yoxdur — WriteHeader(200) + return

### 6. CORS təhlükəsizlik qaydaları (təcrübə qayaları)

1. **Partial match QADAĞAN:** "ends with example.com" → attacker
   `attackerexample.com` qeydiyyatdan keçirər — YALNIZ tam uzunluqlu
   safelist, exact, case-sensitive müqayisə
2. **`Origin: null` safelist-ə HEÇ VAXT** — sandboxed iframe ilə forge
   oluna bilər
3. **Credentials:** `Access-Control-Allow-Credentials: true` cookie/basic
   auth üçün lazımdır; `*` wildcard ilə BİRLİKDƏ QADAĞANDIR; JS tərəfdə
   `credentials: 'include'` / `withCredentials = true` tələb olunur
4. **Authorization header + wildcard:** Authorization icazə verilibsə,
   origin-siz reflection / `*` → distributed brute-force kapısı
   (trusted origin MÜTLƏQDİR)
5. **Preflight cache:** `Access-Control-Max-Age: 60` (Chrome 5s default;
   Firefox cap 24h; Chromium 2h; `-1` = disable; köhnə brauzerlər
   cache-i təmizləyə bilmir — səhv header uzun müddət ilişər)
6. **Preflight wildcards:** `Access-Control-Allow-Methods: *` / `...-Headers: *`
   → yalnız ~74% brauzer dəstəyi; Authorization wildcard OLMAZ
   (`Authorization, *` yaz); credentialed sorğularda `*` literal kimi
   davranır

## Bu fəsildə layihəyə nə əlavə olundu? (project continuity)

**Yeni fayllar:**
- `cmd/examples/cors/simple/main.go` — :9000 test səhifəsi (fetch GET)
- `cmd/examples/cors/preflight/main.go` — POST + JSON Content-Type testi

**Dəyişdirilən:**
- middleware.go → enableCORS (Vary + reflection + preflight intercept)
- routes.go → zəncirə enableCORS (recoverPanic-dən sonra, rateLimit-dən
  əvvəl)
- main.go → cors.trustedOrigins + flag.Func + strings.Fields

## Əsas terminlər

- Origin (mənşə) — scheme+host+port üçlüyü
- Same-Origin Policy — brauzerin cross-origin oxuma qadağası
- CORS (Cross-Origin Resource Sharing) — qadağanı seçici yumşaltma
  mexanizmi
- CORS-safe Method/Header — preflight tələb etməyən sadə elementlər
- Preflight Request — OPTIONS ilə gələn icazə yoxlaması
- Reflection (eko) — Origin header-inin safelist yoxlamasından sonra geri
  qaytarılması
- Vary Header — cache-dəyişənliyi bildirişi
- CSRF (Cross-Site Request Forgery) — same-origin-un göndərməyə icazə
  verməsindən doğan hücum (oxumağa YOX)

## Praktik nəticə

CORS server qərarı deyil, BRAUZER müqaviləsidir: server yalnız "etimadnamə"
(header-lər) verir. Praktik recepi: exact-match safelist + reflection +
Vary:Origin + preflight intercept (200 OK) + credentials üçün ayrıca
diqqət. Wildcard yalnız tamamilə açıq, credentials-sız API-lərdə
qəbulediləndir.

## Mənbə
Pages: 395-424 (raw 395-424)
