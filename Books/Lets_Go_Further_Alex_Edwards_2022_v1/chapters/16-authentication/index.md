# Chapter 16 — Authentication

## Bu fəsil nədən bəhs edir?

5 authentication yanaşmasının müqayisəsi, stateful token pattern-inin seçilməsi,
POST /v1/tokens/authentication ilə credential→token mübadiləsi və
authenticate() middleware (Bearer header + request context).

## Əsas fikirlər

### 1. Authentication vs Authorization
- **Authentication (kimlik təsdiqi):** sorğu HANSI user-dən gəlir
- **Authorization (səlahiyyət):** o user nəyə icazəlidir (növbəti fəsil)
- Bütün metodlar HTTPS tələb edir

### 2. 5 yanaşmanın müqayisəsi

| Yanaşma | Mexanizm | Üstün | Çatışmazlıq | Nə zaman |
|---|---|---|---|---|
| **HTTP Basic** | `Authorization: Basic base64(user:pass)` hər sorğuda | Sadə; universal dəstək; dillər/curl/brauzer built-in | Slow hash hər sorğuda yoxlanılır → server yüku + latency | Real user hesabı olmayan, özünü qoruyan API-lər |
| **Stateful token** (bearer) | Credential→token mübadilə (endpoint), sonra `Authorization: Bearer <token>`; token (və ya fast hash) DB-də user_id + expiry ilə | Password check yalnız token aləmi/olanda; revoke asan (DB-dən sil); sadə və robust konsept | DB lookup hər sorğuda (amma activation-status üçün onsuz da lazım) | Website/SPA backend — login moment var |
| **Stateless token (JWT/PASETO/Branca)** | User ID + expiry token ÖZÜNDƏ; imzalanır (+şifrələnə bilər) | DB lookup yox; in-memory encode/decode | **Revoke edilə bilmir** (secret rotate hamısı + blocklist stateless-i öldürür); JWT konfiqurasiya xətalarına açıq; token daxilində köhnə data authorization üçün TƏHLÜKƏLİDİR (activation status/permission token-a YAZMA) | Delegated auth — auth servisi istehlakçı servisdən FƏRQLİDİR (microservices) |
| **API key** | `Authorization: Key <key>`; daimi, vaxtı keçmir; fast hash DB-də | Client üçün ən sadə — token idarəsi yoxdur | 2 uzunömürlü secret (pass + key); rotate/multi-key support lazım | General-purpose API-lər (script/inteqrasiya) |
| **OAuth 2.0 / OpenID Connect** | 3rd-party identity provider (Google və s.) | Parolları özün saxlamırsan | OAuth2 auth protokolu DEYİL — auth üçün OpenID Connect lazımdır; mürəkkəb (go-oidc kömək edir); brauzer interaction tələb edir | Bütün user-lərin provider hesabı var + API website backend-i |

**Seçim qaydaları (rough):**
1. Real hesab yoxdur + slow hash yoxdur → **Basic**
2. Parol saxlamaq istəmirsən + OIDC provider + website backend → **OpenID Connect**
3. Delegated auth (microservices) → **Stateless (JWT)**
4. Əks halda → **Stateful token** (login momenti varsa — bizim hal) və ya
   **API key** (general-purpose)

**Kitabın qərarı:** Stateful authentication token — activation token
infrastrukturunun (tokens cədvəli, generateToken, GetForToken) birbaşa
yenidən istifadəsi. JWT appendix-də ayrıca göstərilir.

### 3. Token struct tag-ləri + authentication scope
**Kitabdan kod nümunəsi:**
```go
const (
    ScopeActivation     = "activation"
    ScopeAuthentication = "authentication"
)

type Token struct {
    Plaintext string    `json:"token"`  // client üçün "plaintext" → "token"
    Hash      []byte    `json:"-"`
    UserID    int64     `json:"-"`
    Expiry    time.Time `json:"expiry"`
    Scope     string    `json:"-"`
}
```
- JSON cavabda yalnız `{"token": "...", "expiry": "..."}` görünür

### 4. createAuthenticationTokenHandler
**Kitabdan kod nümunəsi:**
```go
// POST /v1/tokens/authentication
func (app *application) createAuthenticationTokenHandler(w http.ResponseWriter, r *http.Request) {
    var input struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }
    err := app.readJSON(w, r, &input)
    if err != nil {
        app.badRequestResponse(w, r, err)
        return
    }

    v := validator.New()
    data.ValidateEmail(v, input.Email)
    data.ValidatePasswordPlaintext(v, input.Password)
    if !v.Valid() {
        app.failedValidationResponse(w, r, v.Errors)
        return
    }

    user, err := app.models.Users.GetByEmail(input.Email)
    if err != nil {
        switch {
        case errors.Is(err, data.ErrRecordNotFound):
            app.invalidCredentialsResponse(w, r)
        default:
            app.serverErrorResponse(w, r, err)
        }
        return
    }

    match, err := user.Password.Matches(input.Password)
    if err != nil {
        app.serverErrorResponse(w, r, err)
        return
    }
    if !match {
        app.invalidCredentialsResponse(w, r)
        return
    }

    token, err := app.models.Tokens.New(user.ID, 24*time.Hour, data.ScopeAuthentication)
    if err != nil {
        app.serverErrorResponse(w, r, err)
        return
    }

    err = app.writeJSON(w, http.StatusCreated, envelope{"authentication_token": token}, nil)
    // ...
}
```

**Sub-kod izahı:**
- ValidateEmail/ValidatePasswordPlaintext — Chapter 13-dəki standalone
  helper-lərin yenidən istifadəsi (plan işlədi!)
- `invalidCredentialsResponse` → 401 "invalid authentication credentials" —
  email yoxdur və ya parol yanlış EYNİ cavab (user enumeration qarşısı —
  bu endpoint-də)
- `Tokens.New(userID, 24h, ScopeAuthentication)` → 24 saatlıq token
- **201 Created** + `{"authentication_token": {"token", "expiry"}}`
- **Authorization header cavabda göndərmək HTTP spec pozuntusudur** — o
  REQUEST header-dır; body-də göndərilir

### 5. AnonymousUser + request context
**Kitabdan kod nümunəsi:**
```go
// internal/data/users.go
var AnonymousUser = &User{}

func (u *User) IsAnonymous() bool {
    return u == AnonymousUser  // pointer müqayisəsi!
}
```

```go
// cmd/api/context.go
type contextKey string  // custom tip — naming collision qoruması
const userContextKey = contextKey("user")

func (app *application) contextSetUser(r *http.Request, user *data.User) *http.Request {
    ctx := context.WithValue(r.Context(), userContextKey, user)
    return r.WithContext(ctx)
}

func (app *application) contextGetUser(r *http.Request) *data.User {
    user, ok := r.Context().Value(userContextKey).(*data.User)
    if !ok {
        panic("missing user value in request context") // gözlənilməz hal
    }
    return user
}
```

**Sub-kod izahı:**
- `AnonymousUser` — boş User-in singleton pointer-i; `IsAnonymous()`
  pointer identity ilə müqayisə edir (`otherUser := &User{}` → false)
- Context values `any` tipindədir → type assertion ilə geri alınır
- **Custom contextKey tipi** — 3rd-party paketlərin "user" string key-i ilə
  toqquşmasının qarşısı
- `contextGetUser` yalnız user OLMALI olduqda çağrılır → yoxdursa panic
  (fail-fast; developer xətası)

### 6. authenticate() middleware
**Kitabdan kod nümunəsi:**
```go
func (app *application) authenticate(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Cavabın cache-lənməsinin Authorization-dan asılı olduğunu bildir
        w.Header().Add("Vary", "Authorization")

        authorizationHeader := r.Header.Get("Authorization")

        // Header yoxdur → AnonymousUser → davam
        if authorizationHeader == "" {
            r = app.contextSetUser(r, data.AnonymousUser)
            next.ServeHTTP(w, r)
            return
        }

        // "Bearer <token>" parçala
        headerParts := strings.Split(authorizationHeader, " ")
        if len(headerParts) != 2 || headerParts[0] != "Bearer" {
            app.invalidAuthenticationTokenResponse(w, r)
            return
        }
        token := headerParts[1]

        // Format validasiyası
        v := validator.New()
        if data.ValidateTokenPlaintext(v, token); !v.Valid() {
            app.invalidAuthenticationTokenResponse(w, r)
            return
        }

        // DB lookup — authentication scope!
        user, err := app.models.Users.GetForToken(data.ScopeAuthentication, token)
        if err != nil {
            switch {
            case errors.Is(err, data.ErrRecordNotFound):
                app.invalidAuthenticationTokenResponse(w, r)
            default:
                app.serverErrorResponse(w, r, err)
            }
            return
        }

        r = app.contextSetUser(r, user)
        next.ServeHTTP(w, r)
    })
}

// errors.go
func (app *application) invalidAuthenticationTokenResponse(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("WWW-Authenticate", "Bearer") // 401 ilə birlikdə gözlənilən scheme
    app.errorResponse(w, r, http.StatusUnauthorized, "invalid or missing authentication token")
}
```

**Sub-kod izahı:**
- **`Vary: Authorization`** → cache-lərə: cavab Authorization header-inin
  dəyərindən asılıdır — fərqli user-lərin cavabları bir-birini əvəz edə bilməz
- Header yoxdur → ANONİM davam (401 YOX!) — auth məcburiyyəti endpoint
  səviyyəsində (requireAuthenticatedUser — növbəti fəslin require-ı)
- `Bearer` böyük hərflə müqayisə — RFC 6750 scheme
- `GetForToken(ScopeAuthentication, ...)` → activation token ilə auth
  token bir-birini əvəz edə BİLMƏZ (scope izolyasiyası)
- `WWW-Authenticate: Bearer` → 401 cavabının standart hissəsi (client-a
  gözlənilən scheme-i bildirir)

**Zəncir sırası:**
```go
return app.recoverPanic(app.rateLimit(app.authenticate(router)))
```
- authenticate router-dən ƏVVƏL — hər endpoint user konteksti alsın

**3 nəticə ssenarisi:**
1. Valid token → user context-də → davam
2. Header yoxdur → AnonymousUser → davam
3. Malformed/invalid token → 401 + Vary + WWW-Authenticate

## Bu fəsildə layihəyə nə əlavə olundu? (project continuity)

**Yeni fayllar:**
- `cmd/api/tokens.go` — createAuthenticationTokenHandler
- `cmd/api/context.go` — contextKey, contextSetUser/contextGetUser

**Dəyişdirilən:**
- tokens.go → ScopeAuthentication + struct tags
- users.go → AnonymousUser, IsAnonymous
- middleware.go → authenticate
- errors.go → invalidCredentialsResponse, invalidAuthenticationTokenResponse
- routes.go → POST /v1/tokens/authentication + authenticate zəncirə

## Əsas terminlər

- Bearer Token — "gətirən şəxsə aid" token (RFC 6750)
- Stateful vs Stateless — server tərəfdə vəziyyət saxlanılır ya token özündə daşıyır
- Revocation (ləğv) — tokenin qüvvədən salınması (stateful asan, stateless çətin)
- Delegated Authentication (nümayəndəlik autentifikasiyası) — tokeni bir servis
  yaradır, digəri istehlak edir
- Vary Header — cache-lərə cavab dəyişənlərini bildirən header
- WWW-Authenticate — 401 ilə gözlənilən auth scheme-i bildirən header
- Request Context — sorğunun ömrü boyu mövcud key/value anbarı
- Type Assertion — `any` dəyərin konkret tipə çevrilməsi
- Singleton — tək nüsxə (AnonymousUser pointer identity)

## Praktik nəticə

Yanaşma seçimi kontekstlədir: bu fəslin ən dəyərli hissəsi 5 üsulun
trade-off cədvəlidir — "JWT hər yerdə yaxşıdır" mümkün deyil: revoke
problemi + stale-data təhlükəsi + yalnız delegated auth üçün optimal.
Stateful pattern bu layihədə activation infrastructure-ın 95%-inin
yenidən istifadəsi deməkdir. Vary/WWW-Authenticate header-ləri isə
professional HTTP semantikasının incə detallarıdır.

## Mənbə
Pages: 348-370 (raw 348-370)
