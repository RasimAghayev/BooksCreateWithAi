# Chapter 7 — CRUD Operations

## Bu fəsil nədən bəhs edir?

MovieModel data-access qatının qurulması və dörd CRUD endpoint-inin tam
icrası: Insert (POST), Get (GET :id), Update (PUT :id), Delete (DELETE :id) —
database/sql + pq ilə.

## Əsas fikirlər

### 1. MovieModel — data access layer
**Nədir:** Bütün SQL məntiqini handler-lərdən təcrid edən model qatı.

**Kitabdan kod nümunəsi:**
```go
// internal/data/movies.go
type MovieModel struct {
    DB *sql.DB
}

func (m MovieModel) Insert(movie *Movie) error { return nil } // placeholder
func (m MovieModel) Get(id int64) (*Movie, error) { return nil, nil }
func (m MovieModel) Update(movie *Movie) error { return nil }
func (m MovieModel) Delete(id int64) error { return nil }
```

**Models container (internal/data/models.go):**
```go
var ErrRecordNotFound = errors.New("record not found")

type Models struct {
    Movies MovieModel
}

func NewModels(db *sql.DB) Models {
    return Models{Movies: MovieModel{DB: db}}
}
```

**Sub-kod izahı:**
- `MovieModel` → `*sql.DB` pool-u wrap edir; handler-lər
  `app.models.Movies.Insert(...)` kimi oxunaqlı çağırır
- `Models` struct-ı → bütün modellər üçün vahid konteyner (sonra UserModel,
  PermissionModel gələcək); `NewModels(db)` ilə main()-də bir dəfə yaradılır,
  `application` struct-ına dependency kimi ötürülür
- `ErrRecordNotFound` → paket-səviyyəli sentinel xəta; handler `errors.Is`
  ilə 404-ə çevirir

**Mocking (Additional Information):** `Models.Movies` sahəsini konkret tip
yerinə **interface** elan etməklə (`interface { Insert(...) error; Get(...) }`)
+ `MockMovieModel` + `NewMockModels()` → unit test-lərdə real DB olmadan model
mock-lanır. (Interface yalnız test lazım olanda daxil edin.)

### 2. Insert — RETURNING ilə sistem dəyərlərinin qaytarılması
**SQL:**
```sql
INSERT INTO movies (title, year, runtime, genres)
VALUES ($1, $2, $3, $4)
RETURNING id, created_at, version
```

**Kitabdan kod nümunəsi:**
```go
func (m MovieModel) Insert(movie *Movie) error {
    query := `
        INSERT INTO movies (title, year, runtime, genres)
        VALUES ($1, $2, $3, $4)
        RETURNING id, created_at, version`

    args := []any{movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres)}

    return m.DB.QueryRow(query, args...).Scan(&movie.ID, &movie.CreatedAt, &movie.Version)
}
```

**Sub-kod izahı:**
- `$N` placeholder-ləri → SQL injection-a qarşı mütləq (untrusted input)
- `RETURNING` → Postgres-ə xas (standart deyil) — INSERT/UPDATE/DELETE-in
  toxunduğu sətirdən dəyər qaytarır; id/created_at/version DB tərəfindən
  yaradılır
- Cavab tək sətir olduğu üçün `Exec()` YOX, `QueryRow()` + `Scan()`
- `*Movie` pointer → method struct-ı **mutasiya** edir (id, created_at,
  version yazılır)
- `pq.Array(movie.Genres)` → `[]string` → `pq.StringArray` (`driver.Valuer` +
  `sql.Scanner` implement edir) → `text[]` sütununa məncərə; `[]bool`,
  `[]byte`, `[]int32`... üçün də işləyir
- `args` slice pattern-i → 3+ parametrli sorğularda aydınlıq üçün

**Handler (createMovieHandler) — 201 + Location:**
```go
err = app.models.Movies.Insert(movie)
if err != nil { app.serverErrorResponse(w, r, err); return }

headers := make(http.Header)
headers.Set("Location", fmt.Sprintf("/v1/movies/%d", movie.ID))
err = app.writeJSON(w, http.StatusCreated, envelope{"movie": movie}, headers)
```
- `201 Created` + `Location` header → REST ənənəsi: yeni resursun URL-i
- writeJSON-ın `headers` parametri məhz bu hal üçün lazım idi

**$N xüsusiyyəti (Additional Info):** eyni `$1` bir çox yerdə istifadə
oluna bilər (`SET bar = $1 + $2 WHERE bar = $1`). Çoxlu statement bir
`Exec()`-də yerləşir, amma placeholder TƏKİ QADAĞANDIR: `pq: cannot insert
multiple commands into a prepared statement`.

### 3. Get — sql.ErrNoRows → ErrRecordNotFound
**Kitabdan kod nümunəsi:**
```go
func (m MovieModel) Get(id int64) (*Movie, error) {
    if id < 1 {
        return nil, ErrRecordNotFound // DB-ə getmədən shortcut
    }
    query := `
        SELECT id, created_at, title, year, runtime, genres, version
        FROM movies
        WHERE id = $1`

    var movie Movie
    err := m.DB.QueryRow(query, id).Scan(
        &movie.ID, &movie.CreatedAt, &movie.Title, &movie.Year,
        &movie.Runtime, pq.Array(&movie.Genres), &movie.Version,
    )
    if err != nil {
        switch {
        case errors.Is(err, sql.ErrNoRows):
            return nil, ErrRecordNotFound
        default:
            return nil, err
        }
    }
    return &movie, nil
}
```

**Sub-kod izahı:**
- PK üzrə sorğu → maksimum 1 sətir → `QueryRow()`
- `sql.ErrNoRows` (Scan-in boş nəticə xətası) → model-səviyyəli
  `ErrRecordNotFound`-a çevrilir — handler DB detalını BİLMİR
- `pq.Array(&movie.Genres)` → Scan istiqamətində də adapter tələb olunur
  (olmasa: `unsupported Scan, storing driver.Value type []uint8 into type
  *[]string`)

**Handler — xəta triage-ı:**
```go
movie, err := app.models.Movies.Get(id)
if err != nil {
    switch {
    case errors.Is(err, data.ErrRecordNotFound):
        app.notFoundResponse(w, r)
    default:
        app.serverErrorResponse(w, r, err)
    }
    return
}
```

**Niyə uint64 YOX, int64 (Additional Info):**
1. PostgreSQL-də unsigned integer YOXDUR → tipləri düzləşdirmək lazımdır:
   `smallint/smallserial → int16`, `integer/serial → int32`,
   `bigint/bigserial → int64`
2. `database/sql` int64 maksimumundan böyük dəyəri dəstəkləmir: `uint64
   values with high bit set are not supported`

### 4. Update — version artımı ilə
**SQL:**
```sql
UPDATE movies
SET title = $1, year = $2, runtime = $3, genres = $4, version = version + 1
WHERE id = $5
RETURNING version
```

**Kitabdan kod nümunəsi:**
```go
func (m MovieModel) Update(movie *Movie) error {
    query := `...`
    args := []any{movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres), movie.ID}
    return m.DB.QueryRow(query, args...).Scan(&movie.Version)
}
```

**Sub-kod izahı:**
- `version = version + 1` → hər update-də avtomatik artım; `RETURNING
  version` → yeni dəyər struct-a yazılır (optimistic concurrency üçün
  əsasdır — Chapter 8-də istifadə olunacaq)
- **PUT semantikası:** tam əvəzetmə — client BÜTÜN redaktə oluna bilən
  sahələri göndərməlidir (id, created_at, version server tərəflidir və
  dəyişmir)

**updateMovieHandler axını (7 addım):**
1. `readIDParam()` → URL-dən ID
2. `Get(id)` → mövcud record (yoxdursa 404)
3. `readJSON()` → input struct
4. input → movie copy (yalnız editable sahələr)
5. `ValidateMovie()` → 422
6. `Update(movie)` → DB-yə yaz
7. `writeJSON` 200 + envelope

### 5. Delete — Exec + RowsAffected
**Kitabdan kod nümunəsi:**
```go
func (m MovieModel) Delete(id int64) error {
    if id < 1 {
        return ErrRecordNotFound
    }
    query := `DELETE FROM movies WHERE id = $1`
    result, err := m.DB.Exec(query, id)
    if err != nil {
        return err
    }
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }
    if rowsAffected == 0 {
        return ErrRecordNotFound
    }
    return nil
}
```

**Sub-kod izahı:**
- Sətir qaytarmayan sorğu → `Exec()` (QueryRow YOX)
- `sql.Result.RowsAffected()` → 1 = silindi; 0 = belə ID yox idi →
  `ErrRecordNotFound` → handler-də 404
- Əvvəlcədən Get ilə yoxlamağa ehtiyac YOXDUR — RowsAffected bu məlumatı
  özü verir (race-free: GET-then-DELETE arasında record silinə bilərdi)

**Handler cavabı:** 200 + `{"message": "movie successfully deleted"}` —
və ya 204 No Content (maşın-clientlər üçün; insan-clientlər üçün mesaj
daha yaxşı UX — kitabın qeydi).

### 6. Final routes
```go
router.HandlerFunc(http.MethodGet,    "/v1/healthcheck", app.healthcheckHandler)
router.HandlerFunc(http.MethodPost,   "/v1/movies",      app.createMovieHandler)
router.HandlerFunc(http.MethodGet,    "/v1/movies/:id",  app.showMovieHandler)
router.HandlerFunc(http.MethodPut,    "/v1/movies/:id",  app.updateMovieHandler)
router.HandlerFunc(http.MethodDelete, "/v1/movies/:id",  app.deleteMovieHandler)
```

## ⚠️ Uyğunsuzluq (Bölmə 21.4)
Chapter 7.4-də (səh. 159) update route-u `http.MethodPut` kimi qeydiyyata
keçirilir, amma Chapter 7.5-in yekun routes siyahısında (səh. 164) eyni
route `http.MethodPatch` göstərilir. Kitabın gövdə mətni PUT (tam
əvəzetmə) semantikası izah edir → **PUT düzgündür**, səh. 164 çap
səhfi kimi görünür. (Nəticə: öz layihənizdə PUT istifadə edin; PATCH
yalnız Chapter 8-də partial update üçün gələcək.)

## Bu fəsildə layihəyə nə əlavə olundu? (project continuity)

**Yeni fayllar:**
- `internal/data/models.go` — Models konteyneri, NewModels, ErrRecordNotFound
- movies.go-da MovieModel (Insert/Get/Update/Delete tam implementasiya)
- movies.go: updateMovieHandler, deleteMovieHandler
- routes.go: PUT/DELETE route-ları

**Application struct:** `models data.Models` sahəsi əlavə olundu (DI).

**DB state:** 4 test filmi (Moana, Black Panther, Deadpool, The Breakfast
Club); Deadpool silindi.

## Əsas terminlər

- Data Access Layer (məlumat çıxışı qatı) — SQL məntiqinin təcridi
- Placeholder Parameter (yer tutucu parametr) — `$N` — SQL injection qoruması
- RETURNING clause — dəyişdirilən sətirdən dəyər qaytaran Postgres imkanı
- driver.Valuer / sql.Scanner — Go↔DB tip çevirici interfeyslər (pq.Array)
- Sentinel Error — paket səviyyəsində müqayisə xətası (ErrRecordNotFound)
- Optimistic Concurrency Control — versiya yoxlaması ilə yazma konflikti
  idarəsi (hazırlıq)

## Praktik nəticə

Model pattern-inin gücü handler-lərin qısalığındadır: Get → 404 triage,
Insert → 201+Location, Delete → RowsAffected → 404 — bütün DB detalları
modeldə, HTTP semantikası handler-də. `ErrRecordNotFound` sentinel-i bu
ayırıcının açarıdır. `version = version + 1` isə növbəti fəslin optimistic
locking-inin təməlidir.

## Mənbə
Pages: 135-164 (raw 135-164)
