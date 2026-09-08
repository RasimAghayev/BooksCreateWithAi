# Chapter 6 — SQL Migrations

## Bu fəsil nədən bəhs edir?

`golang-migrate` tool ilə version-lanmış SQL migration-lar: up/down fayl
cütlüyü, sequential nömrələmə, `schema_migrations` cədvəli, dirty state
bərpası və migration-ların proqramatik icrası.

## Əsas fikirlər

### 1. Migration konsepsiyası
**Nədir:** Hər schema dəyişikliyi (cədvəl yarat, column əlavə et...) üçün
**up** (dəyişikliyi tətbiq et) + **down** (geri qaytar) fayl cütlüyü;
sıralı nömrələnir (0001, 0002... və ya Unix timestamp).

**Üstünlükləri:**
- Schema + onun təkamülü tam SQL faylları ilə təsvir olunur → version
  control-da kodla birlikdə saxlanılır
- Başqa maşında dəqiq eyni schema bərpa olunur (dev/test/prod sinxronu)
- Geri qaytarma (roll-back) mümkündür

### 2. migrate create — fayl cütlüyünün yaradılması
**Kitabdan komanda:**
```bash
migrate create -seq -ext=.sql -dir=./migrations create_movies_table
# Nəticə:
# migrations/000001_create_movies_table.up.sql
# migrations/000001_create_movies_table.down.sql
```

**Sub-flag izahı:**
- `-seq` → sequential nömrələmə (default Unix timestamp-dır)
- `-ext=.sql` → fayl genişlənməsi
- `-dir=./migrations` → qovluq (avtomatik yaranır)
- `create_movies_table` → təsviri ad

### 3. movies cədvəli — up/down migration
**Kitabdan kod nümunəsi (up):**
```sql
-- migrations/000001_create_movies_table.up.sql
CREATE TABLE IF NOT EXISTS movies (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    title text NOT NULL,
    year integer NOT NULL,
    runtime integer NOT NULL,
    genres text[] NOT NULL,
    version integer NOT NULL DEFAULT 1
);
```

**Sub-kod izahı:**
- `bigserial` → 64-bit avtomatik artan ID (1-dən başlayır)
- `timestamp(0) with time zone` → saniyə dəqiqliyi, timezone-aware
- `text[]` → Postgres massiv tipi — **queryable və indexable** (Janrlar
  filtri üçün sonra istifadə olunacaq)
- `text` (varchar YOX) → Postgres-də varchar(n) ilə text arasında performans
  fərqi yoxdur; text daha sadədir
- `NOT NULL` + `DEFAULT` → Go-da NULL idarəsi çətin olduğundan bütün
  sütunlar NOT NULL-dur

**Down (geri qaytarma):**
```sql
DROP TABLE IF EXISTS movies;
```
- DROP TABLE index və constraint-ləri də avtomatik silir — tək statement
  kifayətdir

### 4. Business rule-lar DB səviyyəsində — CHECK constraint-lər
**Kitabdan kod nümunəsi:**
```sql
-- 000002_add_movies_check_constraints.up.sql
ALTER TABLE movies ADD CONSTRAINT movies_runtime_check CHECK (runtime >= 0);
ALTER TABLE movies ADD CONSTRAINT movies_year_check CHECK (year BETWEEN 1888 AND date_part('year', now()));
ALTER TABLE movies ADD CONSTRAINT genres_length_check CHECK (array_length(genres, 1) BETWEEN 1 AND 5);

-- down.sql
ALTER TABLE movies DROP CONSTRAINT IF EXISTS movies_runtime_check;
ALTER TABLE movies DROP CONSTRAINT IF EXISTS movies_year_check;
ALTER TABLE movies DROP CONSTRAINT IF EXISTS genres_length_check;
```

**Sub-kod izahı:**
- `CHECK (...)` → dəyər şərti; `date_part('year', now())` → cari il dinamik
  yoxlanır
- `array_length(genres, 1)` → 1-ci ölçü üzrə massiv uzunluğu
- Pozulsa driver xətası: `pq: new row for relation "movies" violates check
  constraint "movies_year_check"`
- Bir migration faylında çoxlu SQL ola bilər

**Arxitektura qərarı:** Application-level validation (422) +
DB-level CHECK = ikiqat müdafiə — birbaşa SQL/paralel yazıcı olsa belə
data pozula bilməz.

### 5. migrate up — icra və schema_migrations
**Kitabdan komanda:**
```bash
migrate -path=./migrations -database=$GREENLIGHT_DB_DSN up
# 1/u create_movies_table (38.19761ms)
# 2/u add_movies_check_constraints (63.284269ms)
```

**Sub-flag izahı:**
- `-path` → migration qovluğu; `-database` → DSN; `up` → növbəti tətbiq
- Tool avtomatik `schema_migrations` cədvəli yaradır:
  - `version` → son tətbiq olunan migration nömrəsi
  - `dirty` → xətasız icra statusu (false = clean)

**Digər komandalar:**
```bash
migrate ... version          # cari versiyanı göstər
migrate ... goto 1           # konkret versiyaya keç (up və ya down avtomatik)
migrate ... down 1           # son 1 migration-ı geri qaytar
migrate ... down             # HAMSINI geri qaytar (təsdiq istəyir)
migrate ... force 1          # dirty halında versiyanı məcburi təyin et
migrate -source="s3://<bucket>/<path>" ... up   # uzaq mənbədən
migrate -source="github://owner/repo/path#ref" ... up
```

**`drop` haqqında xəbərdarlıq:** bütün cədvəlləri (schema_migrations daxil)
silir, amma sequence/enum-lar qalır → "messy and unknown state"; `down`
üstünlük təşkil edir.

### 6. Dirty state — xətalı migration bərpası
**Nədir:** Migration faylında sintaksis xətası olsa, xətaya qədərki SQL-lər
TƏTBİQ OLUNUR, sonra `dirty=true` qoyulur; bu halda heç bir yeni migration
(yəne down belə) icra olunmur: `Dirty database version {X}. Fix and force
version.`

**Bərpa proseduru:**
1. Orijinal xətanı araşdır — fayl qismən tətbiq olunubsa manual roll-back et
2. `migrate ... force {N}` → versiyanı düzgün dəyərə məcburi keçir
3. Artıq "clean" — migration-lar yenidən işləyir

### 7. Startup-da migration (anti-pattern müzakirəsi)
**Kitabdan kod nümunəsi (kitab bunu ETMİR, sadəcə göstərir):**
```go
migrationDriver, err := postgres.WithInstance(db, &postgres.Config{})
migrator, err := migrate.NewWithDatabaseInstance("file:///path/to/migrations", "postgres", migrationDriver)
err = migrator.Up()
if err != nil && err != migrate.ErrNoChange {
    logger.PrintFatal(err, nil)
}
```

**Sub-kod izahı:**
- `postgres.WithInstance(db, ...)` → mövcud `*sql.DB` üzərində migration driver
- `migrate.NewWithDatabaseInstance` → file source + DB instance birləşdirir
- `migrate.ErrNoChange` → hər şey tətbiq olunubsa normal hal — xəta deyil

**Niyin anti-pattern ola bilər:** Migration icrasını server start-up-a bağlamaq
uzunmüddətdə problem yaradır — bir neçə app replikası eyni anda start olarsa
race yaranır; deploy tooling-i migration-a cavabdeh olmalıdır ("decoupling
database migrations from server startup").

## Bu fəsildə layihəyə nə əlavə olundu? (project continuity)

**Yeni fayllar:**
- `migrations/000001_create_movies_table.{up,down}.sql`
- `migrations/000002_add_movies_check_constraints.{up,down}.sql`

**DB vəziyyəti:** greenlight DB-da `movies` cədvəli + CHECK
constraint-lər; `schema_migrations` v2/clean.

## Əsas terminlər

- Migration (sxem köçürməsi) — schema dəyişikliyinin version-lanmış SQL
  təsviri
- Rollback (geri qaytarma) — dəyişikliyin əksinə icra
- Dirty State (çirkli vəziyyət) — yarımçıq tətbiq olunmuş migration statusu
- CHECK Constraint (yoxlama məhdudiyyəti) — sütun dəyəri üzərində şərt
- bigserial — 64-bit avtomatik artan serial tip
- Extension points: S3/GitHub source — migration fayllarının uzaq mənbəsi

## Praktik nəticə

golang-migrate + up/down cütlüyü + sequential versiya — standartlaşmış,
bütün dillərə daşınabilən yanaşmadır. Vacib praktikalar: business
rule-ları CHECK-lə DB-yə də verin (422 yalnız ilk müdafiə xəttidir);
xətalı migration-dan sonra `force` ilə bərpa edin; migration-ları app
start-up-ından DEYİL, deploy pipeline-dan icra edin. (Müəllim qeydi:
alternativlər — `goose`, `atlas`, `tern` — oxşar konseptlər; `migrate`
ən geniş yayılmışıdır.)

## Mənbə
Pages: 124-134 (raw 124-134)
