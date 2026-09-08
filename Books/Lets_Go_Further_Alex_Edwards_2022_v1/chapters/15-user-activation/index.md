# Chapter 15 — User Activation

## Bu fəsil nədən bəhs edir?

Aktivasiya workflow-unun tam dövrü: tokens cədvəli (hash + FK + scope),
CSPRNG token generasiyası (crypto/rand + base32), TokenModel,
registerUserHandler-da token yaradılıb email-ə göndərilməsi, PUT
/v1/users/activated ilə INNER JOIN əsaslı aktivasiya.

## Əsas fikirlər

### 1. Aktivasiya workflow-u (6 addım)
1. Qeydiyyatda cryptographic-random aktivasiya token-i yaradılır
2. Token-in SHA-256 hash-i tokens cədvəlində saxlanılır (userID + expiry ilə)
3. Plaintext token email-də göndərilir
4. User token-i PUT /v1/users/activated-ə təqdim edir
5. Hash tapılır və vaxtı keçməyibsə activated=true
6. Token silinir (bir dəfəlik istifadə)

**Faydası:** botlara əlavə maneə + saxta email qeydiyyatının qarşısı.

### 2. tokens cədvəli — migration 000005
**Kitabdan kod nümunəsi:**
```sql
CREATE TABLE IF NOT EXISTS tokens (
    hash bytea PRIMARY KEY,
    user_id bigint NOT NULL REFERENCES users ON DELETE CASCADE,
    expiry timestamp(0) with time zone NOT NULL,
    scope text NOT NULL
);
```

**Sub-kod izahı:**
- `hash bytea PRIMARY KEY` → yalnız SHA-256 hash saxlanılır — plaintext
  YOX (DB leak halında qoruma; paroldan fərqli olaraq token yüksək
  entropy-li olduğu üçün SÜRƏTLİ SHA-256 kifayətdir, bcrypt lazım deyil)
- `REFERENCES users` → foreign key: user_id həmişə mövcud user-ə işarə edir
- `ON DELETE CASCADE` → user silinəndə tokenləri avtomatik silinir;
  alternativ `ON DELETE RESTRICT` (user-i tokenləri silmədən silmək olmaz)
- `expiry` → 3 gün; qısa müddət brute-force pencərəsini və sonrakı email
  hesabının təhlükəsini azaldır
- `scope` → tokenin məqsədi ("activation"; sonra "authentication" gələcək) —
  bir cədvəl, çox token tipi

### 3. Token generasiyası — crypto/rand + base32
**Kitabdan kod nümunəsi:**
```go
const ScopeActivation = "activation"

type Token struct {
    Plaintext string
    Hash      []byte
    UserID    int64
    Expiry    time.Time
    Scope     string
}

func generateToken(userID int64, ttl time.Duration, scope string) (*Token, error) {
    token := &Token{
        UserID: userID,
        Expiry: time.Now().Add(ttl),
        Scope:  scope,
    }

    randomBytes := make([]byte, 16)             // 128 bit entropy
    _, err := rand.Read(randomBytes)            // OS CSPRNG
    if err != nil {
        return nil, err
    }

    token.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).
        EncodeToString(randomBytes)             // 26 simvol

    hash := sha256.Sum256([]byte(token.Plaintext))
    token.Hash = hash[:]                        // [32]byte array → slice
    return token, nil
}
```

**Sub-kod izahı:**
- `rand.Read` (crypto/rand) → OS-in CSPRNG-dən 16 random bayt; math/rand
  (deterministic PRNG) token/secret üçün HEÇ VAXT işlədilməz
- base32 (no padding) → 16 bayt → 26 simvol string
  (`Y3QMGX3PJ3WLRL2YRTQGQ6KRHU`); hex olsaydı 32 simvol olardı
- **Entropy 16 baytdır, string uzunluğu yox** — kodlaşdırma uzunluğu dəyişir,
  təhlükəsizlik ölçüsü entropy-dir
- `sha256.Sum256` → [32]byte ARRAY qaytarır; `hash[:]` ilə slice-a çevrilir
  (pq driver array qəbul etmir)

### 4. TokenModel + ValidateTokenPlaintext
**Kitabdan kod nümunəsi:**
```go
func ValidateTokenPlaintext(v *validator.Validator, tokenPlaintext string) {
    v.Check(tokenPlaintext != "", "token", "must be provided")
    v.Check(len(tokenPlaintext) == 26, "token", "must be 26 bytes long")
}

type TokenModel struct {
    DB *sql.DB
}

// Shortcut: yarat + daxil et
func (m TokenModel) New(userID int64, ttl time.Duration, scope string) (*Token, error) {
    token, err := generateToken(userID, ttl, scope)
    if err != nil {
        return nil, err
    }
    err = m.Insert(token)
    return token, err
}

func (m TokenModel) Insert(token *Token) error {
    query := `INSERT INTO tokens (hash, user_id, expiry, scope) VALUES ($1, $2, $3, $4)`
    // ExecContext 3s ...
}

func (m TokenModel) DeleteAllForUser(scope string, userID int64) error {
    query := `DELETE FROM tokens WHERE scope = $1 AND user_id = $2`
    // ...
}
```

### 5. registerUserHandler — token + email
**Kitabdan kod nümunəsi:**
```go
token, err := app.models.Tokens.New(user.ID, 3*24*time.Hour, data.ScopeActivation)
if err != nil {
    app.serverErrorResponse(w, r, err)
    return
}

app.background(func() {
    data := map[string]any{
        "activationToken": token.Plaintext,
        "userID":          user.ID,
    }
    err = app.mailer.Send(user.Email, "user_welcome.tmpl", data)
    if err != nil {
        app.logger.PrintError(err, nil)
    }
})
```

**Sub-kod izahı:**
- `3*24*time.Hour` → 3 günlük TTL
- Template data map → `{{.activationToken}}`, `{{.userID}}` render olunur
- Email PUT endpoint-inə yönləndirir — LINK YOX (aşağıda səbəbi)

**Niyə GET-link deyil, PUT?**
1. GET "safe" metod olmalı — state dəyişməməli (HTTP prinsipi)
2. Brauzer/antivirus prefetch → hesab istəmədən aktivləşər (Eve Alice-in
   email-i ilə qeydiyyatdan keçərsə, Alice-in brauzeri linki özü açıb
   aktivləşdirər)
- **Qayda:** state dəyişən hərəkətlər YALNIZ POST/PUT/PATCH/DELETE

### 6. activateUserHandler — PUT /v1/users/activated
**Kitabdan kod nümunəsi:**
```go
func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
    var input struct {
        TokenPlaintext string `json:"token"`
    }
    err := app.readJSON(w, r, &input)
    if err != nil {
        app.badRequestResponse(w, r, err)
        return
    }

    v := validator.New()
    if data.ValidateTokenPlaintext(v, input.TokenPlaintext); !v.Valid() {
        app.failedValidationResponse(w, r, v.Errors)
        return
    }

    user, err := app.models.Users.GetForToken(data.ScopeActivation, input.TokenPlaintext)
    if err != nil {
        switch {
        case errors.Is(err, data.ErrRecordNotFound):
            v.AddError("token", "invalid or expired activation token")
            app.failedValidationResponse(w, r, v.Errors)
        default:
            app.serverErrorResponse(w, r, err)
        }
        return
    }

    user.Activated = true
    err = app.models.Users.Update(user)
    // ErrEditConflict → 409 triage ...

    err = app.models.Tokens.DeleteAllForUser(data.ScopeActivation, user.ID)
    // ...

    err = app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
    // ...
}
```

**Niyə PUT (POST yox)?** Endpoint idempotentdir — təkrar sorğularda ilkin
uğurdan sonra xəta qayıtsa da, application state artıq dəyişmir.

### 7. GetForToken — INNER JOIN
**Kitabdan kod nümunəsi:**
```go
func (m UserModel) GetForToken(tokenScope, tokenPlaintext string) (*User, error) {
    tokenHash := sha256.Sum256([]byte(tokenPlaintext))

    query := `
        SELECT users.id, users.created_at, users.name, users.email,
               users.password_hash, users.activated, users.version
        FROM users
        INNER JOIN tokens
        ON users.id = tokens.user_id
        WHERE tokens.hash = $1
        AND tokens.scope = $2
        AND tokens.expiry > $3`

    args := []any{tokenHash[:], tokenScope, time.Now()}
    // QueryRowContext → Scan; ErrNoRows → ErrRecordNotFound
}
```

**Sub-kod izahı:**
- Client-in plaintext-ini hash-ləyib axtarırır — plaintext heç vaxt
  WHERE-ə düşmür
- `INNER JOIN ... ON users.id = tokens.user_id` → müvəqqəti birləşmiş
  cədvəl; hash PK olduğundan nəticə dəqiq 1 sətir (və ya 0)
- `expiry > time.Now()` → vaxtı keçmişlər tapılmır
- Model mübahisəsi: UserModel.GetForToken user qaytarır (model = entity),
  TokenModel.GetAllForUser tokenlər qaytarardı — hər model öz entity-sini
  qaytarır

### 8. Security qeydləri (Additional Info)
- **HTTPS məcburidir** — plaintext token yalnız şifrəli kanalla qəbul
  olunmalıdır
- **Web workflow:** form-a copy-paste (ən robust) və ya link + JS (Referrer
  Policy: `Origin` lazımdır — token referrer-də sızmasın)
- **Host header injection:** email-dəki link-i `r.Host`-dan QURMAQ OLMAZ —
  domain hard-code və ya flag olmalıdır
- **Timing attack nəzəri riski:** `tokens.hash = $1` Postgres-də
  constant-time müqayisə deyil; amma uğursuz belə olsa yalnız HASH
  sızlar — plaintext-i tapmaq üçün SHA-256 preimage brute-force lazımdır
  (texnoloji olaraq qeyri-mümkün)

## Bu fəsildə layihəyə nə əlavə olundu? (project continuity)

**Yeni fayllar:**
- `migrations/000005_create_tokens_table.{up,down}.sql`
- `internal/data/tokens.go` — Token, generateToken, ValidateTokenPlaintext,
  TokenModel (New/Insert/DeleteAllForUser), ScopeActivation
- users.go: GetForToken, activateUserHandler

**Dəyişdirilən:** registerUserHandler (token + map data); routes.go → PUT
/v1/users/activated; user_welcome.tmpl (activationToken); Models + Tokens.

## Əsas terminlər

- CSPRNG — cryptographic secure pseudo-random generator (OS səviyyəli)
- Entropy (entropiya) — təsadüfüyin ölçüsü (128 bit = 2^128 ehtimal)
- Base32 — 5 bit/simvol kodlaşdırma (token 26 simvol)
- SHA-256 — sürətli cryptographic hash (yüksək entropy üçün kifayət)
- Foreign Key (xarici açar) — cədvəllərarası istinad məhdudiyyəti
- ON DELETE CASCADE — parent silinəndə uşaqların avtomatik silinməsi
- INNER JOIN — uyğun sətirlərin kəsişim birləşməsi
- One-to-Many — bir user → çox token münasibəti
- Idempotency (eynilik) — təkrar tətbiq eyni nəticə (PUT seçiminin səbəbi)
- Timing Attack — müqayisə vaxtından məlumat çıxarma hücumu

## Praktik nəticə

Token pattern-inin əsas qaydası: **DB-də hash, client-də plaintext** —
parol pattern-inin tokenlər üçünvariantı, amma yüksək entropy sayəsində
bcrypt yerinə SHA-256 (sürətli) kifayətdir. Scope sahəsi bir cədvələ çox
token tipi (activation/auth) verir. INNER JOIN + `expiry > now()` + tək
istifadə (delete) — aktivasiyanın üç təhlükəsizlik halqasıdır.

## Mənbə
Pages: 323-347 (raw 323-347)
