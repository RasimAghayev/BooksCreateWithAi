# Chapter 9 — Filtering, Sorting, and Pagination

## Bu fəsil nədən bəhs edir?

GET /v1/movies list endpoint-inin mərhələli qurulması: query string parse
helper-ləri, Filters struct + validation, `QueryContext` ilə çoxsətirli oxu,
optional filtering (fixed SQL + `@>`), PostgreSQL full-text search + GIN
index-lər, safelist-əsaslı sort, LIMIT/OFFSET pagination və `count(*) OVER()`
window function ilə metadata.

## Əsas fikirlər

### 1. Query string parse helper-ləri
**Nədir:** `r.URL.Query()` → `url.Values` map; dəyərlərin çıxarılması +
default + konversiya üçün 3 helper.

**Kitabdan kod nümunəsi:**
```go
func (app *application) readString(qs url.Values, key string, defaultValue string) string {
    s := qs.Get(key)
    if s == "" {
        return defaultValue
    }
    return s
}

func (app *application) readCSV(qs url.Values, key string, defaultValue []string) []string {
    csv := qs.Get(key)
    if csv == "" {
        return defaultValue
    }
    return strings.Split(csv, ",")
}

func (app *application) readInt(qs url.Values, key string, defaultValue int, v *validator.Validator) int {
    s := qs.Get(key)
    if s == "" {
        return defaultValue
    }
    i, err := strconv.Atoi(s)
    if err != nil {
        v.AddError(key, "must be an integer value")
        return defaultValue
    }
    return i
}
```

**Sub-kod izahı:**
- `readString` → string + default fallback
- `readCSV` → vergüllə ayrılmış dəyərlər slice-a (`genres=crime,drama`)
- `readInt` → int konversiya; xətada validator-a error yazılır (422 mexanizmi
  ilə birləşir) + default qaytarılır

**Hədəf query nümunəsi:**
`/v1/movies?title=godfather&genres=crime,drama&page=1&page_size=5&sort=-year`
- `-` prefiksi → descending (`sort=-year` = il üzrə azalan)

### 2. Filters struct + validation
**Nədir:** page/page_size/sort parametrləri başqa endpoint-lərdə də işləyəcək
— yenidən istifadə olunan struct.

**Kitabdan kod nümunəsi:**
```go
// internal/data/filters.go
type Filters struct {
    Page         int
    PageSize     int
    Sort         string
    SortSafelist []string
}

func ValidateFilters(v *validator.Validator, f Filters) {
    v.Check(f.Page > 0, "page", "must be greater than zero")
    v.Check(f.Page <= 10_000_000, "page", "must be a maximum of 10 million")
    v.Check(f.PageSize > 0, "page_size", "must be greater than zero")
    v.Check(f.PageSize <= 100, "page_size", "must be a maximum of 100")
    v.Check(validator.PermittedValue(f.Sort, f.SortSafelist...), "sort", "invalid sort value")
}
```

**Handler-də:**
```go
var input struct {
    Title  string
    Genres []string
    data.Filters  // embed
}
input.Filters.SortSafelist = []string{"id", "title", "year", "runtime",
                                      "-id", "-title", "-year", "-runtime"}
```

**Sub-kod izahı:**
- `SortSafelist` → endpoint-in hansı sort dəyərlərini dəstəklədiyini
  təyin edir — SQL injection qorumasının əsası
- `PermittedValue` generic helper (Chapter 4-dən) safelist yoxlaması edir
- Default-lar: page=1, page_size=20, sort="id"

### 3. GetAll — QueryContext + rows iterasiyası
**Nədir:** Çoxsətirli nəticə üçün `Query()`/`QueryContext()` + `rows.Next()`
döngüsü.

**Kitabdan kod nümunəsi:**
```go
func (m MovieModel) GetAll(title string, genres []string, filters Filters) ([]*Movie, Metadata, error) {
    query := fmt.Sprintf(`
        SELECT count(*) OVER(), id, created_at, title, year, runtime, genres, version
        FROM movies
        WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '')
        AND (genres @> $2 OR $2 = '{}')
        ORDER BY %s %s, id ASC
        LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())

    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    args := []any{title, pq.Array(genres), filters.limit(), filters.offset()}
    rows, err := m.DB.QueryContext(ctx, query, args...)
    if err != nil {
        return nil, Metadata{}, err
    }
    defer rows.Close()

    totalRecords := 0
    movies := []*Movie{}
    for rows.Next() {
        var movie Movie
        err := rows.Scan(&totalRecords, &movie.ID, &movie.CreatedAt, &movie.Title,
            &movie.Year, &movie.Runtime, pq.Array(&movie.Genres), &movie.Version)
        if err != nil {
            return nil, Metadata{}, err
        }
        movies = append(movies, &movie)
    }
    if err = rows.Err(); err != nil {
        return nil, Metadata{}, err
    }
    metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)
    return movies, metadata, nil
}
```

**Sub-kod izahı:**
- `defer rows.Close()` → MüTLƏQ — resultset bağlantını tutur; hər iterasiyada
  `rows.Next()` false olanda `rows.Err()` İÇRA XƏTASINI yoxlayır (network
  kəsilməsi və s. — `Scan` xətaları deyil)
- `movies := []*Movie{}` → pointer slice — nil slice-ın `append`-i təhlükəsizdir,
  boş nəticədə JSON `[]` (nil yox) verir

### 4. Optional filtering — fixed SQL pattern
**Nədir:** Dinamik SQL qurmaq (concat/interpolate) yerinə hər filtrin "optional"
olduğu SABİT sorğu:

```sql
WHERE (LOWER(title) = LOWER($1) OR $1 = '')
AND (genres @> $2 OR $2 = '{}')
```

**Sub-kod izahı:**
- `$1 = ''` → title boş string (default) olduqda şərt həmişə true → filtr
  "keçilir"; empty slice default-u `'{}'` ilə eyni məntiq
- `@>` → Postgres array "contains" operatoru: `$2`-dəki HAMSI `genres`
  sütununda varsa true (AND semantikası: `genres=animation,adventure` →
  hər ikisi olmalı)
- Digər array operatorları: `&&` (overlap), `<@` (contained by)

### 5. Full-text search — to_tsvector / plainto_tsquery
**Nədir:** Natural dil axtarışı — `title=panther` "Black Panther"-i tapır
(partial match).

**SQL:**
```sql
WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '')
```

**Komponentlər:**
- `to_tsvector('simple', title)` → title-i lexeme-lərə bölür ("The Breakfast
  Club" → 'breakfast' 'club' 'the'); `simple` konfiqurasiyası yalnız
  lowercase edir; `english` → stopword çıxarma + stemming ('cuckoo' 'flew'
  'nest' 'one')
- `plainto_tsquery('simple', $1)` → axtarış dəyərini sorğu termininə çevirir,
  sözlər arasına `&` (AND) qoyur ("the club" → 'the' & 'club')
- `@@` → matches operatoru

**Alternativlər (Additional Info):**
- `STRPOS(LOWER(title), LOWER($1)) > 0` → substring; amma `the` axtarışı
  "Black Panther"-i də tapır (unintuitive) + index-lənMİR → full-table scan
- `title ILIKE $1` → pattern; client `%` wildcard idarə edir (`the%25`);
  `pg_trgm` extension + GIN index ilə index-lənə bilər

**GIN index migration (000003):**
```sql
CREATE INDEX IF NOT EXISTS movies_title_idx ON movies USING GIN (to_tsvector('simple', title));
CREATE INDEX IF NOT EXISTS movies_genres_idx ON movies USING GIN (genres);
```
- GIN → array və full-text axtarışları üçün index tipi; lexeme-lər hər
  sorğuda yenidən hesablanmaqdan xilas edir

### 6. Sort — safelist + fmt.Sprintf
**Problem:** Placeholder parametrləri column adları və ASC/DESC keyword-ləri
üçün İŞLƏMİR → `fmt.Sprintf` interpolasiyası lazımdır → SQL injection riski!

**Həll: safelist + panic fail-safe:**
```go
func (f Filters) sortColumn() string {
    for _, safeValue := range f.SortSafelist {
        if f.Sort == safeValue {
            return strings.TrimPrefix(f.Sort, "-")
        }
    }
    panic("unsafe sort parameter: " + f.Sort) // fail-safe
}

func (f Filters) sortDirection() string {
    if strings.HasPrefix(f.Sort, "-") {
        return "DESC"
    }
    return "ASC"
}
```

**ORDER BY qaydası:**
```sql
ORDER BY %s %s, id ASC
```
- **İkinci dərəcəli sort həmişə PK (id)** — Postgres ORDER BY-sız qayda
  zəmanət etmir; pagination-da səhifələrarası "tullanma" olmasın deyə
  tam deterministik sıra lazımdır

### 7. Pagination — LIMIT/OFFSET
**Riyaziyyat:**
```
LIMIT  = page_size
OFFSET = (page - 1) * page_size
```

```go
func (f Filters) limit() int  { return f.PageSize }
func (f Filters) offset() int { return (f.Page - 1) * f.PageSize }
```
- Integer overflow riski → ValidateFilters-in max dəyərləri (10M × 100)
  tərəfindən məhdudlaşdırılıb
- LIMIT/OFFSET placeholder-lərlə ($3/$4) — dəyər təbiətinə görə təhlükəsizdir

### 8. Pagination metadata — count(*) OVER()
**Nədir:** Window function — filtrə uyan ÜMUMİ sətir sayını hər sətirdə qaytarır
(ayrıca COUNT sorğusuna ehtiyac yoxdur).

```sql
SELECT count(*) OVER(), id, created_at, ... LIMIT $3 OFFSET $4
```

**İcra ardıcıllığı (simplified):** WHERE filtri → window count → ORDER BY →
LIMIT/OFFSET. Yəni count PAGINATION-DAN ƏVVƏL hesablanır — filtered totali
verir.

**Metadata struct:**
```go
type Metadata struct {
    CurrentPage  int `json:"current_page,omitempty"`
    PageSize     int `json:"page_size,omitempty"`
    FirstPage    int `json:"first_page,omitempty"`
    LastPage     int `json:"last_page,omitempty"`
    TotalRecords int `json:"total_records,omitempty"`
}

func calculateMetadata(totalRecords, page, pageSize int) Metadata {
    if totalRecords == 0 {
        return Metadata{} // boş → omitempty ilə tam gizlənir
    }
    return Metadata{
        CurrentPage:  page,
        PageSize:     pageSize,
        FirstPage:    1,
        LastPage:     int(math.Ceil(float64(totalRecords) / float64(pageSize))),
        TotalRecords: totalRecords,
    }
}
```
- `math.Ceil` → 12 record / 5 page_size = 2.4 → last_page = 3
- `omitempty` → boş metadata `{}` kimi çıxır
- Handler: `envelope{"movies": movies, "metadata": metadata}`

## Bu fəsildə layihəyə nə əlavə olundu? (project continuity)

**Yeni fayllar:**
- `internal/data/filters.go` — Filters, ValidateFilters, sortColumn/sortDirection
  (private), limit/offset (private), Metadata, calculateMetadata (private)
- `migrations/000003_add_movies_indexes.{up,down}.sql` — GIN index-lər

**Dəyişdirilən:** movies.go → GetAll (filter+sort+pagination+metadata);
helpers.go → readString/readCSV/readInt; routes.go → GET /v1/movies;
listMoviesHandler tam funksional.

## Əsas terminlər

- Query String (sorğu sətri) — URL-də `?`-dan sonrakı key=value cütlükləri
- Safelist (təhlükəsiz siyahı) — icazə verilən dəyərlərin ağ siyahısı
- Lexeme (leksim) — full-text search-in token vahidi
- Window Function (pəncərə funksiyası) — sətirlər qrupu üzrə hesablama
  (`count(*) OVER()`)
- GIN Index — array/full-text üçün Generalized Inverted index
- Pagination (səhifələmə) — nəticənin hissə-hissə qaytarılması
- Interpolation (interpolasiya) — dəyərin string içinə yerləşdirilməsi
  (placeholder-in əksi; təhlükəsizlik tələb edir)

## Praktik nəticə

Bu fəsil kitabın ən "real world" hissəsidir: fixed-SQL optional filter pattern,
safelist sort, GIN index, LIMIT/OFFSET + `count(*) OVER()` metadata —
hamısı bir araya yığılmış list endpoint recepi. Filters/Metadata struct-ları
kopyalanabilir utility-dirlər. (Müəllim qeydi: çox böyük datalarda OFFSET
ineffektivdir — deep pagination üçün keyset/seek pagination
(`WHERE id > $last`) alternativini bilin; kitab səviyyəsində LIMIT/OFFSET
kifayətdir.)

## Mənbə
Pages: 189-233 (raw 189-233)
