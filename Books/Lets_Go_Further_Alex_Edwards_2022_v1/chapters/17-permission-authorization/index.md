# Chapter 17 — Permission-based Authorization

## Bu fəsil nədən bəhs edir?

Authorization qatının qurulması: requireActivatedUser middleware (401/403
ayrımı), permissions + users_permissions cədvəlləri (many-to-many + composite
PK), PermissionModel (GetAllForUser/AddForUser), requirePermission
middleware zənciri və qeydiyyat zamanı default permission.

## Əsas fikirlər

### 1. 401 vs 403 status kodları
- **401 Unauthorized** → autentifikasiya yoxdur/pis (anonim istifadəçi)
- **403 Forbidden** → autentifikasiya OK, amma icazə yoxdur (aktiv olmayan
  hesab və ya permission çatışmır)
- **Qayda:** 401 həmişə auth mərhələsi üçün; 403 ondan SONRA gələn
  səlahiyyət yoxlamaları üçün

### 2. requireActivatedUser middleware
**Kitabdan kod nümunəsi:**
```go
// errors.go
func (app *application) authenticationRequiredResponse(w http.ResponseWriter, r *http.Request) {
    app.errorResponse(w, r, http.StatusUnauthorized, "you must be authenticated to access this resource")
}
func (app *application) inactiveAccountResponse(w http.ResponseWriter, r *http.Request) {
    app.errorResponse(w, r, http.StatusForbidden, "your user account must be activated to access this resource")
}

// middleware.go
func (app *application) requireAuthenticatedUser(next http.HandlerFunc) http.HandlerFunc {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        user := app.contextGetUser(r)
        if user.IsAnonymous() {
            app.authenticationRequiredResponse(w, r)
            return
        }
        next.ServeHTTP(w, r)
    })
}

func (app *application) requireActivatedUser(next http.HandlerFunc) http.HandlerFunc {
    fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        user := app.contextGetUser(r)
        if !user.Activated {
            app.inactiveAccountResponse(w, r)
            return
        }
        next.ServeHTTP(w, r)
    })
    // Mövcud olanı avtomatik requireAuthenticatedUser ilə wrap et
    return app.requireAuthenticatedUser(fn)
}
```

**Sub-kod izahı:**
- İmza `http.HandlerFunc → http.HandlerFunc` (Handler YOX) — handler
  funksiyalarını birbaşa wrap etmək üçün konversiyasız
- **Middleware kompozisiyası:** requireActivatedUser daxildən
  requireActivatedUser → requireAuthenticatedUser çağırır — "kim olduğunu
  bilmədən aktivliyi yoxlamaq mənasızdır"
- Zəncir: anonymous → 401; authenticated+inactive → 403; hər ikisi OK → next

**Alternativ (Additional Info):** Endpoint azdıqsa handler DAXİLİNDƏ yoxlama
da normal yanaşmadır (`user := app.contextGetUser(r); if user.IsAnonymous()
{...}`).

### 3. permissions cədvəlləri — many-to-many
**Kitabdan kod nümunəsi (migration 000006):**
```sql
CREATE TABLE IF NOT EXISTS permissions (
    id bigserial PRIMARY KEY,
    code text NOT NULL
);

CREATE TABLE IF NOT EXISTS users_permissions (
    user_id bigint NOT NULL REFERENCES users ON DELETE CASCADE,
    permission_id bigint NOT NULL REFERENCES permissions ON DELETE CASCADE,
    PRIMARY KEY (user_id, permission_id)  -- composite PK!
);

INSERT INTO permissions (code) VALUES ('movies:read'), ('movies:write');
```

**Sub-kod izahı:**
- Many-to-many → **joining table** (users_permissions) klassik həlli
- **Composite PK (user_id, permission_id)** → eyni cüt dublikat ola bilməz
- Hər iki FK `ON DELETE CASCADE` → user/permission silinəndə əlaqələr
  avtomatik təmizlənir
- Permission kodları `resource:action` formatında (movies:read,
  movies:write) — endpoint-lərlə xəritələnir

**Endpoint → permission xəritəsi:**
| Endpoint | Permission |
|---|---|
| GET /v1/movies, GET /v1/movies/:id | movies:read |
| POST/PATCH/DELETE /v1/movies... | movies:write |
| healthcheck, users, tokens | — (ictimai) |

### 4. PermissionModel — GetAllForUser
**Kitabdan kod nümunəsi:**
```go
type Permissions []string

func (p Permissions) Include(code string) bool {
    for i := range p {
        if code == p[i] {
            return true
        }
    }
    return false
}

type PermissionModel struct {
    DB *sql.DB
}

func (m PermissionModel) GetAllForUser(userID int64) (Permissions, error) {
    query := `
        SELECT permissions.code
        FROM permissions
        INNER JOIN users_permissions ON users_permissions.permission_id = permissions.id
        INNER JOIN users ON users_permissions.user_id = users.id
        WHERE users.id = $1`
    // QueryContext 3s → rows.Next() → append → rows.Err()
}
```

**Sub-kod izahı:**
- `Permissions []string` custom slice tipi + `Include()` helper — handler-də
  deklarativ `permissions.Include("movies:write")`
- İkili INNER JOIN: permissions ↔ users_permissions ↔ users — 3 cədvəlli
  sorğu; model öz entity-sini (permission kodları) qaytarır

### 5. requirePermission middleware
**Kitabdan kod nümunəsi:**
```go
func (app *application) requirePermission(code string, next http.HandlerFunc) http.HandlerFunc {
    fn := func(w http.ResponseWriter, r *http.Request) {
        user := app.contextGetUser(r)

        permissions, err := app.models.Permissions.GetAllForUser(user.ID)
        if err != nil {
            app.serverErrorResponse(w, r, err)
            return
        }

        if !permissions.Include(code) {
            app.notPermittedResponse(w, r) // 403
            return
        }
        next.ServeHTTP(w, r)
    }
    return app.requireActivatedUser(fn) // zəncir: authenticated → activated → permission
}

func (app *application) notPermittedResponse(w http.ResponseWriter, r *http.Request) {
    app.errorResponse(w, r, http.StatusForbidden,
        "your user account doesn't have the necessary permissions to access this resource")
}
```

**Route tətbiqi:**
```go
router.HandlerFunc(http.MethodGet, "/v1/movies",
    app.requirePermission("movies:read", app.listMoviesHandler))
router.HandlerFunc(http.MethodPost, "/v1/movies",
    app.requirePermission("movies:write", app.createMovieHandler))
// ...
```

**Sub-kod izahı:**
- Middleware-in BİRİNCİ parametri permission kodu — per-endpoint konfiqurasiya
- Zəncir kompozisiyası: `requirePermission` → `requireActivatedUser` →
  `requireAuthenticatedUser` → handler — üç yoxlama bir deklarasiyada
- Tam zəncir: recoverPanic → rateLimit → authenticate → router →
  requirePermission → handler

### 6. AddForUser — variadic + ANY($2)
**Kitabdan kod nümunəsi:**
```go
func (m PermissionModel) AddForUser(userID int64, codes ...string) error {
    query := `
        INSERT INTO users_permissions
        SELECT $1, permissions.id FROM permissions WHERE permissions.code = ANY($2)`
    // ExecContext 3s
    _, err := m.DB.ExecContext(ctx, query, userID, pq.Array(codes))
    return err
}
```

**Sub-kod izahı:**
- `codes ...string` → variadic: `AddForUser(2, "movies:read", "movies:write")`
- `ANY($2)` → Postgres array üzvlük yoxlaması; `pq.Array(codes)` adapteri
  slice → array
- INSERT...SELECT → müvəqqəti cədvəl: (userID, hər kodun ID-si) sətirləri
  bir sorğuda daxil edilir

### 7. Default permission qeydiyyatda
registerUserHandler-a əlavə (Insert uğurundan sonra):
```go
err = app.models.Permissions.AddForUser(user.ID, "movies:read")
```
- Yeni user-ə avtomatik movies:read; movies:write admin tərəfindən əl ilə
  verilir (demo: `INSERT INTO users_permissions VALUES ((SELECT id FROM
  users WHERE email='...'), (SELECT id FROM permissions WHERE code='...'))`)

## Bu fəsildə layihəyə nə əlavə olundu? (project continuity)

**Yeni fayllar:**
- `migrations/000006_add_permissions.{up,down}.sql` — permissions +
  users_permissions + seed
- `internal/data/permissions.go` — Permissions type + Include,
  PermissionModel (GetAllForUser, AddForUser)

**Dəyişdirilən:**
- middleware.go → requireAuthenticatedUser, requireActivatedUser,
  requirePermission (kompozit zəncir)
- errors.go → authenticationRequiredResponse, inactiveAccountResponse,
  notPermittedResponse
- routes.go → 5 movies endpoint-i requirePermission ilə
- users.go → registerUserHandler-da AddForUser
- Models → PermissionModel

## Əsas terminlər

- Authorization (səlahiyyətvermə) — kimin nəyə icazəsi var
- 401 vs 403 — auth yoxluğu vs icazə yoxluğu
- Many-to-Many Relationship — çox-çox əlaqə (user ↔ permission)
- Joining Table (birləşdirmə cədvəli) — many-to-many üçün aralıq cədvəl
- Composite Primary Key (mürəkkəb əsas açar) — birden çox sütundan PK
- Middleware Composition — middleware-lərin bir-birini wrap etməsi
- Variadic Parameter — `...T` — dəyişən sayda arqument
- Permission Code — `resource:action` formatlı icazə identifikatoru

## Praktik nəticə

Middleware zənciri kompozisiyası bu fəslin memarlıq dərsidir:
requirePermission daxilən requireActivatedUser, o da requireAuthenticatedUser
çağırır — hər yoxlama müstəqil, lakin istifadə yerdə bir sətirdə. 401→403→handler
ardıcıllığı HTTP semantikasının düzgün təcəssümüdür. Many-to-many +
composite PK + cascade FK isə relational DB-də RBAC-ın minimal, təmiz
formasını verir.

## Mənbə
Pages: 371-394 (raw 371-394)
