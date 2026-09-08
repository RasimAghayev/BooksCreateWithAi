# Chapter 13 — User Model Setup and Registration

## Bu fəsil nədən bəhs edir?

users cədvəli (citext + bytea), password custom type (bcrypt Set/Matches),
ValidateUser, UserModel (Insert/GetByEmail/Update + duplicate email triage)
və POST /v1/users registration endpoint-i.

## Əsas fikirlər

### 1. users cədvəli — migration 000004
**Kitabdan kod nümunəsi:**
```sql
CREATE TABLE IF NOT EXISTS users (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    name text NOT NULL,
    email citext UNIQUE NOT NULL,
    password_hash bytea NOT NULL,
    activated bool NOT NULL,
    version integer NOT NULL DEFAULT 1
);
```

**Sub-kod izahı:**
- `citext` → case-insensitive text: dəyər orijinal casing-də saxlanılır,
  amma müqayisələr (UNIQUE daxil) həmişə case-insensitive → `Alice@X.com`
  və `alice@x.com` duplicate sayılır
- `UNIQUE` → users_email_key avtomatik adlı constraint — duplicate insert
  xətası budur
- `bytea` → binary sütun; bcrypt hash burada saxlanılır (psql hex göstərir:
  `\x243261...`)
- `activated` → default false; email təsdiqindən sonra true
- `version` → optimistic locking üçün (movie-dəki pattern-in təkrarı)

### 2. password custom type — bcrypt
**Kitabdan kod nümunəsi:**
```go
type User struct {
    ID        int64     `json:"id"`
    CreatedAt time.Time `json:"created_at"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    Password  password  `json:"-"`     // JSON-a heç vaxt düşmür
    Activated bool      `json:"activated"`
    Version   int       `json:"-"`
}

type password struct {
    plaintext *string // pointer — "" (boş) ilə "yoxdur" fərqi
    hash      []byte
}

func (p *password) Set(plaintextPassword string) error {
    hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
    if err != nil {
        return err
    }
    p.plaintext = &plaintextPassword
    p.hash = hash
    return nil
}

func (p *password) Matches(plaintextPassword string) (bool, error) {
    err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))
    if err != nil {
        switch {
        case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
            return false, nil // normal mismatch — xəta deyil
        default:
            return false, err
        }
    }
    return true, nil
}
```

**Sub-kod izahı:**
- `bcrypt.GenerateFromPassword(input, 12)` → cost=12; nəticə formatı:
  `$2b$[cost]$[22-char salt][31-char hash]`
- Cost qəbulü: attacker üçün bahalı, user üçün dözülən olmalıdır
- `CompareHashAndPassword` → eyni salt+cost ilə yenidən hash edib
  `subtle.ConstantTimeCompare` (constant-time — timing attack qoruması)
  ilə müqayisə edir
- `*string` plaintext → nil ("mövcud deyil") vs "" ("boş verilib") fərqi —
  PATCH pattern-inin təkrarı
- `Matches` mismatch-i xəta QAYTARMIR (false) — yalnız texniki xətalar error

**bcrypt 72 bayt limiti:** input 72 baytdan kəsilir → validation-da
maksimum 72 qoyulub (uzun parolların "gizli hissəsi" ignore olunmasın deyə).

### 3. ValidateUser + standalone helper-lər
**Kitabdan kod nümunəsi:**
```go
func ValidateEmail(v *validator.Validator, email string) {
    v.Check(email != "", "email", "must be provided")
    v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}

func ValidatePasswordPlaintext(v *validator.Validator, password string) {
    v.Check(password != "", "password", "must be provided")
    v.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
    v.Check(len(password) <= 72, "password", "must not be more than 72 bytes long")
}

func ValidateUser(v *validator.Validator, user *User) {
    v.Check(user.Name != "", "name", "must be provided")
    v.Check(len(user.Name) <= 500, "name", "must not be more than 500 bytes long")
    ValidateEmail(v, user.Email)
    if user.Password.plaintext != nil {
        ValidatePasswordPlaintext(v, *user.Password.plaintext)
    }
    if user.Password.hash == nil {
        panic("missing password hash for user") // developer xətası — fail fast
    }
}
```

**Sub-kod izahı:**
- Standalone ValidateEmail/ValidatePasswordPlaintext → sonra login/activation
  workflow-larında təkrar istifadə olunacaq
- `hash == nil` → client xətası deyil, MANTIQ xətasıdır → validation map
  yerinə panic (fail-fast prinsipi)

### 4. UserModel — Insert/GetByEmail/Update
**Kitabdan kod nümunəsi:**
```go
var ErrDuplicateEmail = errors.New("duplicate email")

type UserModel struct {
    DB *sql.DB
}

func (m UserModel) Insert(user *User) error {
    query := `
        INSERT INTO users (name, email, password_hash, activated)
        VALUES ($1, $2, $3, $4)
        RETURNING id, created_at, version`
    args := []any{user.Name, user.Email, user.Password.hash, user.Activated}
    // ctx 3s + QueryRowContext ...
    err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.CreatedAt, &user.Version)
    if err != nil {
        switch {
        case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
            return ErrDuplicateEmail
        default:
            return err
        }
    }
    return nil
}

func (m UserModel) GetByEmail(email string) (*User, error) {
    // SELECT ... WHERE email = $1 → Scan (hash birbaşa &user.Password.hash)
    // sql.ErrNoRows → ErrRecordNotFound
}

func (m UserModel) Update(user *User) error {
    // UPDATE ... WHERE id = $5 AND version = $6 RETURNING version
    // duplicate key → ErrDuplicateEmail; ErrNoRows → ErrEditConflict
}
```

**Sub-kod izahı:**
- Duplicate email string müqayisəsi ilə tutulur (pq hələ typed error
  vermir — readJSON-dakı "json: unknown field" kimi string prefix yanaşması)
- GetByEmail: UNIQUE sayəsində max 1 sətir; hash birbaşa password struct-ına
  scan olunur
- Update: movie pattern-inin mükəmməl təkrarı — version optimistic locking +
  ErrEditConflict
- Models struct-a `Users UserModel` əlavə olundu

### 5. registerUserHandler — POST /v1/users
**Kitabdan kod nümunəsi:**
```go
func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
    var input struct {
        Name     string `json:"name"`
        Email    string `json:"email"`
        Password string `json:"password"`
    }
    err := app.readJSON(w, r, &input)
    if err != nil {
        app.badRequestResponse(w, r, err)
        return
    }

    user := &data.User{
        Name:     input.Name,
        Email:    input.Email,
        Activated: false,
    }

    err = user.Password.Set(input.Password)
    if err != nil {
        app.serverErrorResponse(w, r, err)
        return
    }

    v := validator.New()
    if data.ValidateUser(v, user); !v.Valid() {
        app.failedValidationResponse(w, r, v.Errors)
        return
    }

    err = app.models.Users.Insert(user)
    if err != nil {
        switch {
        case errors.Is(err, data.ErrDuplicateEmail):
            v.AddError("email", "a user with this email address already exists")
            app.failedValidationResponse(w, r, v.Errors)
        default:
            app.serverErrorResponse(w, r, err)
        }
        return
    }

    err = app.writeJSON(w, http.StatusCreated, envelope{"user": user}, nil)
    if err != nil {
        app.serverErrorResponse(w, r, err)
    }
}
```

**Sub-kod izahı:**
- `Password.Set()` → hash yaradılır; plaintext saxlanılır (validation üçün)
- Duplicate email → 422 + email field mesajı (500 YOX) — sentinel xətanın
  handler-də field xətasına çevrilməsi
- 201 Created + envelope{"user": user} — Password/Version `json:"-"`
  sayəsində JSON-a DÜŞMÜR

### 6. Email case-sensitivity (Additional Info)
- RFC 2821: domain hissə həmişə case-insensitive; username hissəsi
  provider-dən asılı (əksəriyyət insensitive, zəmanət yoxdur)
- **Saxlama:** orijinal casing ilə (emails yanlış adama gedə bilər)
- **Müqayisə:** case-insensitive (login/activation-də forgiving UX)
- citext bu siyasətin DB-də avtomatik təcəssümüdür

### 7. User enumeration riski (Additional Info)
**Problem:** "a user with this email address already exists" cavabı
attacker-a email-in mövcudluğunu Deyirən edir.

**Risklər:** privacy + social engineering + leaked password cəhdləri.

**Qarşısı (2 tələb):**
1. Cavab həmişə EYNİ olmalıdır (ambiguous wording + email side-channel)
2. Response TIME həmişə eyni (background goroutine ilə)

**Trade-off:** Mitigation UX friksiyası artırır; Twitter/GitHub/Amazon belə
qarşısını almır — "is it worth the trade-off?" layihədən-layihəyə dəyişən
qərardır.

## Bu fəsildə layihəyə nə əlavə olundu? (project continuity)

**Yeni fayllar:**
- `migrations/000004_create_users_table.{up,down}.sql`
- `internal/data/users.go` — User, password, Set/Matches, Validate*,
  UserModel (Insert/GetByEmail/Update), ErrDuplicateEmail
- `cmd/api/users.go` — registerUserHandler

**Dəyişdirilən:** Models + NewModels; routes.go → POST /v1/users.

**Yeni asılılıq:** `golang.org/x/crypto/bcrypt`

## Əsas terminlər

- bcrypt — adaptive parol hash alqoritmi (cost + salt)
- Cost Parameter (xərc parametri) — hashın hesablama bahası
- Salt (duz) — hash-ə əlavə olunan təsadüfi dəyər (rainbow table qoruması)
- Constant-Time Comparison — müqayisə vaxtından məlumat sızmaması
- citext — case-insensitive Postgres tipi
- bytea — binary data sütun tipi
- User Enumeration (istifadəçi sayımı hücumu) — mövcudluğu müəyyən etmə
  hücumu

## Praktik nəticə

Password custom type — domain-in öz içində təhlükəsizlik inkapsulyasiyasının
nümunəsidir: hash/plaintext lifecycle, validation qaydaları və comparison
məntiqi bir type-da cəmlənib. citext + UNIQUE + ErrDuplicateEmail zənciri
DB-dən HTTP-ə qədər tutarlı duplicate email idarəsidir. User enumeration
müzakirəsi isə production qərarlarının "security vs UX" tarazlığını
göstərir.

## Mənbə
Pages: 276-293 (raw 276-293)
