# Chapter 8 — Advanced CRUD Operations

## Bu fəsil nədən bəhs edir?

Partial update (PATCH + pointer input), optimistic concurrency control
(version-əsaslı race qoruması) və context timeout ilə uzun çəkən SQL
sorğularının ləğvi.

## Əsas fikirlər

### 1. Partial Updates — pointer input struct
**Problem:** `{"year": 1985}` göndərildikdə decode olunmayan sahələr
zero-value alır (title="" və s.) → validation xətası. `{"title": ""}`
(göndərilib, boşdur) ilə `"title"` (göndərilməyib) fərqini necə ayırd edək?

**Həll:** Input struct sahələrini **pointer** et — pointer-in zero-value-si
`nil`-dir:

**Kitabdan kod nümunəsi:**
```go
var input struct {
    Title   *string       `json:"title"`   // JSON-da yoxdursa nil
    Year    *int32        `json:"year"`
    Runtime *data.Runtime `json:"runtime"`
    Genres  []string      `json:"genres"` // slice-ın zero-value-si onsuz da nil
}

if input.Title != nil {
    movie.Title = *input.Title  // dereference lazımdır
}
if input.Year != nil {
    movie.Year = *input.Year
}
if input.Runtime != nil {
    movie.Runtime = *input.Runtime
}
if input.Genres != nil {
    movie.Genres = input.Genres // slice dereference tələb etmir
}
```

**Sub-kod izahı:**
- `*string` → JSON-da `"title"` key-i VARSA pointer dolur (dəyər "" olsa
  belə), YOXDURSA nil qalır
- `*input.Title` → dereference ilə əsl dəyər
- Slice onsuz da nil-able olduğundan pointer lazım deyil
- Client boş string göndərərsə (`{"title": ""}`) → `*input.Title != nil` →
  boş dəyər MÜTLƏQ → ValidateMovie 422 qaytarır (düzgün davranış)

**PATCH vs PUT:** Partial update = **PATCH**; PUT = tam əvəzetmə. Route
dəyişdirildi: `router.HandlerFunc(http.MethodPatch, "/v1/movies/:id", ...)`

**JSON null xüsusi halı (Additional Info):** `{"title": null}` → pointer
YİNE nil qalır → field dəyişməz qalır (yalnız version artır). Bu, "key
yoxdur" ilə eyni görünür — encoding/json ilə fərqi ayırd etmək MÜMKÜN
DEYİL. Praktik həll: sənədləşdirmədə "null dəyərlər ignore olunur" yaz.

### 2. Optimistic Concurrency Control (OCC)
**Nədir:** İki client eyni record-u eyni anda update etsə **data race**
yaranır: Alice Get(version=N) → Bob Get(version=N) → Alice Update → Bob
Update (Alice-in dəyişikliyini SƏSSİCƏ üstələyir — "lost update").

**Həll:** Update yalnız gözlənilən versiya hələ DB-də olsa icra olunsun:

```sql
UPDATE movies
SET title = $1, year = $2, runtime = $3, genres = $4, version = version + 1
WHERE id = $5 AND version = $6
RETURNING version
```

**Kitabdan kod nümunəsi:**
```go
// internal/data/models.go
var (
    ErrRecordNotFound = errors.New("record not found")
    ErrEditConflict   = errors.New("edit conflict")
)

func (m MovieModel) Update(movie *Movie) error {
    query := `
        UPDATE movies
        SET title = $1, year = $2, runtime = $3, genres = $4, version = version + 1
        WHERE id = $5 AND version = $6
        RETURNING version`
    args := []any{movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres),
                  movie.ID, movie.Version}

    err := m.DB.QueryRow(query, args...).Scan(&movie.Version)
    if err != nil {
        switch {
        case errors.Is(err, sql.ErrNoRows):
            return ErrEditConflict // versiya dəyişib və ya record silinib
        default:
            return err
        }
    }
    return nil
}
```

**Handler + 409 Conflict:**
```go
// errors.go
func (app *application) editConflictResponse(w http.ResponseWriter, r *http.Request) {
    message := "unable to update the record due to an edit conflict, please try again"
    app.errorResponse(w, r, http.StatusConflict, message)
}

// updateMovieHandler daxilində:
case errors.Is(err, data.ErrEditConflict):
    app.editConflictResponse(w, r)
```

**Necə işləyir:** Birinci Update uğur qazanır (version N→N+1); ikinci Update
`WHERE version = N` tapa bilmir → `sql.ErrNoRows` → `ErrEditConflict` →
409. Test: `xargs -I % -P8 curl -X PATCH ...` 8 paralel sorğu → 1-2 uğur,
qalanları 409.

**Round-trip locking (Additional Info):** Client gözlədiyi versiyanı
`X-Expected-Version` header-ində göndərə bilər → handler Get-dən sonra
yoxlayır:
```go
if r.Header.Get("X-Expected-Version") != "" {
    if strconv.FormatInt(int64(movie.Version), 32) != r.Header.Get("X-Expected-Version") {
        app.editConflictResponse(w, r)
        return
    }
}
```
- GET-dən PATCH-ə qədər başqa client dəyişiklik edibsə client eləcə də
  başlanğıcda öyrənir

**Alternativ versiya sahələri:** artan integer (ucuz, təhlükəsiz — default
tövsiyə); timestamp (riskli — saat xətaları); UUID (`uuid-ossp`
extension: `version = uuid_generate_v4()`) — guessable olmaması vacibdirsə.

### 3. SQL Query Timeouts — context ilə
**Nədir:** Uzun çəkən sorğuları avtomatik ləğv etmək: `ExecContext()` /
`QueryRowContext()` + `context.WithTimeout()`.

**Kitabdan kod nümunəsi (bütün modellərdə 3s):**
```go
func (m MovieModel) Get(id int64) (*Movie, error) {
    if id < 1 {
        return nil, ErrRecordNotFound
    }
    query := `SELECT id, created_at, title, year, runtime, genres, version
              FROM movies WHERE id = $1`
    var movie Movie

    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    err := m.DB.QueryRowContext(ctx, query, id).Scan(
        &movie.ID, &movie.CreatedAt, &movie.Title, &movie.Year,
        &movie.Runtime, pq.Array(&movie.Genres), &movie.Version,
    )
    // ...
}
```

**Sub-kod izahı:**
- `context.WithTimeout(context.Background(), 3*time.Second)` → 3s deadline-li
  context
- `defer cancel()` → MÜTLƏQ — context resurslarını dərhal azad edir
  (yoxsa timeout və ya parent cancel olana qədər saxlanılır → memory leak)
- Deadline context YARADILANDAN sayılır (QueryRowContext çağrılmazdan əvvəl
  keçən vaxt da hesaba düşür)
- Timeout olsa: pq background goroutine context-in `Done` channel-ini
  dinləyir → Postgres-ə ləğv siqnalı → `pq: canceling statement due to user
  request` xətası → handler 500 qaytarır

**İki fərqli timeout ssenarisi:**
1. Sorğu İŞLƏYƏRKƏN timeout → `pq: canceling statement due to user request`
2. Sorğu hələ POOL-DA GÖZLƏYƏRKƏN timeout (bütün 25 connection məşğuldursa)
   → `context deadline exceeded` (`context.DeadlineExceeded`)
   - Demo: `-db-max-open-conns=1` ilə 2 paralel sorğu → biri işləyir,
     digəri növbədə ikisi də ləğv olunur
3. Scan mərhələsində də deadline vurula bilər → `DeadlineExceeded`

**pg_sleep(10) demo texnikası:** test üçün sorğuya `pg_sleep(10)` əlavə
edib ilk Scan hədəfi `&[]byte{}` etmək; sonra geri qaytarmaq.

**Request context istifadəsi (kitab bunu seçMİR):** `r.Context()` parent
olaraq götürmək mümkündür, amma davranış mürəkkəbliyi (client disconnect
→ bütün uşaqlar ləğv) çox vaxt dəyər — ayrıca appendix-də müzakirə olunur.

## Bu fəsildə layihəyə nə əlavə olundu? (project continuity)

**Dəyişdirilənlər:**
- updateMovieHandler → pointer input + PATCH
- Update() → `AND version = $6` + `ErrEditConflict`
- errors.go → `editConflictResponse()` (409)
- Bütün model methodları → `context.WithTimeout(3s)` +
  `QueryRowContext`/`ExecContext`
- models.go → `ErrEditConflict` sentinel-i
- routes.go → PUT əvəzinə PATCH

## Əsas terminlər

- Partial Update (qismən yeniləmə) — yalnız göndərilən sahələrin dəyişməsi
- Optimistic Locking (optimist kilidləmə) — conflict yoxlaması son anda,
  kilidsiz yazma (pessimistic-in əksi)
- Data Race (data yarışı) — nəticəsi icra sırasından asılı olan paralel
  əməliyyatlar toqquşması
- Lost Update (itirilmiş yeniləmə) — bir yazmanın səssizcə üstələnməsi
- Edit Conflict (409) — versiya uyğunsuzluğu statusu
- Deadline (son məddət) — context-in vaxt limiti
- Done Channel — context ləğv siqnal kanalı (pq bunu dinləyir)

## Praktik nəticə

Bu üç pattern production API-nin "böyüklər" mexanizmidir: PATCH+pointer —
client-ə yumşaq interfeys; `WHERE version = $N` — 2 sətirlik SQL ilə lost
update qoruması; context timeout — resurs sızmasını qarşısını alan təhlükəsizlik
kəməri. Versiya sahəsi artan integer olaraq saxlamağın ucuz və etibarlı
seçimidir; UUID yalnız guessability vacibdirsə.

## Mənbə
Pages: 165-188 (raw 165-188)
