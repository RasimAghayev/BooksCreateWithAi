# Chapter 5 — Database Setup and Configuration

## Bu fəsil nədən bəhs edir?

PostgreSQL quraşdırılması, psql ilə DB/user/extension yaradılması, `lib/pq`
driver ilə connection pool qurulması və pool parametrlərinin
(MaxOpenConns/MaxIdleConns/ConnMaxLifetime/ConnMaxIdleTime) nəzəriyyəsi +
konfiqurasiyası.

## Əsas fikirlər

### 1. PostgreSQL setup (psql)
**Nədir:** Lokal Postgres quraşdırması və ilkin DB/user/extension hazırlığı.

**Kitabdan komandalar (CLI):**
```bash
# macOS / Linux / Windows
brew install postgresql
sudo apt install postgresql
choco install postgresql

# Superuser kimi qoşul (peer authentication)
sudo -u postgres psql
```

**Sub-komanda izahı:**
- `sudo -u postgres psql` → OS `postgres` user-i adına psql; peer
  authentication = OS user adı Postgres user adı ilə üst-üstə düşərsə parolsuz
  qoşulur

**SQL daxilində:**
```sql
CREATE DATABASE greenlight;
\c greenlight
CREATE ROLE greenlight WITH LOGIN PASSWORD 'pa55word';
CREATE EXTENSION IF NOT EXISTS citext;
```

**İzah:**
- `\c` → database dəyiş; digər meta-komandalar: `\l` (DB-lər), `\dt`
  (cədvəllər), `\du` (user-lər), `\?` (kömək)
- `CREATE ROLE ... WITH LOGIN PASSWORD` → superuser OLMAYAN, parollu
  istifadəçi (migration və app üçün)
- `citext` extension → case-insensitive text tipi (email üçün); yalnız
  superuser DB-a əlavə edə bilər

**Qoşulma testi:**
```bash
psql --host=localhost --dbname=greenlight --username=greenlight
```

### 2. Go ilə qoşulma — pq driver + DSN
**Nədir:** `database/sql` + `github.com/lib/pq` ilə connection pool.

**DSN (Data Source Name):**
```
postgres://greenlight:pa55word@localhost/greenlight
# SSL xətası olarsa:
postgres://greenlight:pa55word@localhost/greenlight?sslmode=disable
```

**Kitabdan kod nümunəsi:**
```go
import (
    "database/sql"
    _ "github.com/lib/pq" // driver özünü database/sql-a qeyd edir
)

func openDB(cfg config) (*sql.DB, error) {
    db, err := sql.Open("postgres", cfg.db.dsn)
    if err != nil {
        return nil, err
    }
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    err = db.PingContext(ctx)
    if err != nil {
        return nil, err
    }
    return db, nil
}
```

**Sub-kod izahı:**
- `_ "github.com/lib/pq"` → blank import: paket istifadə olunmur kimi
  görünür, amma `init()`-də driver-i qeydiyyata salır
- `sql.Open()` → boş pool yaradır; bağlantılar LAZY yaranır (ilk sorğuda)
- `db.PingContext(ctx)` → real bağlantını İNDİ məcbur edir + 5s deadline;
  konfiqurasiya xətalarını başlanğıcda yoxlamaq üçün
- `defer db.Close()` → main çıxanda pool qapanır

**DSN-in koddan ayrılması:**
```go
flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("GREENLIGHT_DB_DSN"), "PostgreSQL DSN")
```
- Default dəyər environment variable-dan → parol koda/common repo-a düşmür;
  `psql $GREENLIGHT_DB_DSN` kimi təkrar istifadə da mümkün olur

### 3. Connection pool nəzəriyyəsi
**Nədir:** `sql.DB` — 2 növ bağlantı: **in-use** (aktiv task) və **idle**
(boş). Task gələndə idle varsa təkrar istifadə olunur, yoxsa yeni yaradılır.
Pis bağlantılar 2 dəfə retry edilir, sonra pool-dan silinib yenisi yaradılır.

**Dörd konfiqurasiya metodu:**

| Metod | Default | Təsiri |
|---|---|---|
| `SetMaxOpenConns(n)` | limitsiz | open (in-use + idle) üst limit; PostgreSQL default 100 hard limit ("sorry, too many clients already") — ~25 seçmək kifayət; tıxanarsa HTTP request "hang" olur → context timeout MÜTLƏQ |
| `SetMaxIdleConns(n)` | 2 | idle üst limit; yüksək = tez-tez yeni bağlantı yaratmamaq; amma hər idle bağlantı yaddaş aparır + köhnəyə bilər (MySQL 8 saatda qapatır); **həmişə ≤ MaxOpenConns** (Go avtomatik azaldır) |
| `SetConnMaxLifetime(d)` | limitsiz | bağlantının YARADILDIĞI andan ömür müddəti; 1s-də bir background cleanup; çox qısa → dəfərlərə重建 frequency artır (100 conn + 1m lifetime ≈ 1.67 conn/s ölür-yenilənir) |
| `SetConnMaxIdleTime(d)` | limitsiz | SON istifadədən keçən idle vaxtı; yüksək idle limit + idle-timeout kombinasiyası ideal — resources periodik azad olur |

**Kitabın praktik tövsiyələri:**
1. MaxOpenConns-i mütləq təyin et (bu layihə: 25 — kiçik/orta API üçün
   başlanğıc nöqtəsi, benchmark-la tune et)
2. MaxIdleConns = 25 (= MaxOpenConns)
3. ConnMaxIdleTime = 15m (istifadəsiz idlələri təmizlə)
4. ConnMaxLifetime = limitsiz saxla (DB tərəfindən lifetime məhdudiyyəti
   yoxdursa / LB swap lazım deyilsə)

**Kitabdan kod nümunəsi (konfiqurasiya flag-ləri):**
```go
type config struct {
    port int
    env  string
    db   struct {
        dsn          string
        maxOpenConns int
        maxIdleConns int
        maxIdleTime  string
    }
}

flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")

// openDB daxilində:
db.SetMaxOpenConns(cfg.db.maxOpenConns)
db.SetMaxIdleConns(cfg.db.maxIdleConns)
duration, err := time.ParseDuration(cfg.db.maxIdleTime) // "15m" → Duration
db.SetConnMaxIdleTime(duration)
```

**Sub-kod izahı:**
- Duration string formatı (`"5s"`, `"10m"`) + `time.ParseDuration()` →
  integer-saniyə konversiyasından daha oxunaqlı CLI
- `db.Stats()` ilə pool vəziyyəti real-time izlənə bilər (load test
  fəslində istifadə olunacaq)

## Bu fəsildə layihəyə nə əlavə olundu? (project continuity)

**Dəyişdirilən:** main.go — config struct-ına `db` sahəsi (dsn, pool
parametrləri), `openDB()` helper-i, blank pq import, GREENLIGHT_DB_DSN env
default-u; openDB daxilində pool parametrlərinin tətbiqi.

**Yeni xarici tələb:** `github.com/lib/pq@v1`, PostgreSQL server, greenlight
DB + user + citext extension.

## Əsas terminlər

- Connection Pool (bağlantı hovuzu) — bağlantıların yenidən istifadəsi üçün
  idarə olunan dəst
- DSN — Data Source Name (bağlantı parametrlərini daşıyan string)
- Peer Authentication (cütlük autentifikasiyası) — OS user = DB user olduqda
  parolsuz qoşulma
- Lazy Connection (təxirə salınmış bağlantı) — ilk ehtiyac anına qədər
  yaradılmayan bağlantı
- Throttle (sürət məhdudlaşdırıcı) — MaxOpenConns-in şərti rol gücləndirici
  funksiyası
- Extension (genişləndirmə) — Postgres-ə əlavə funksionallıq paketi (citext)

## Praktik nəticə

Pool parametrləri "bir dəfə düzgün təyin et" tipində qərarlardır: MaxOpenConns
DB hard-limitindən aşağa, MaxIdleConns ona bərabər, ConnMaxIdleTime orta
müddətli, ConnMaxLifetime default. Ayrıca: parollar heç vaxt kodda olmasın —
env variable + flag default kombinasiyası bunun ən sadə yolu. (Müəllim
qeydi: productionda DSN secret-manager-dən gəlsin; `pgx` driver-i
müasir layihələrdə `pq`-nu əvəz edir — benchmark edib seçin.)

## Mənbə
Pages: 105-123 (raw 105-123)
