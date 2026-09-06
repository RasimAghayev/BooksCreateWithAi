# Chapter 5 — Interaction with Databases (Verilənlər Bazası ilə Qarşılıqlı Əlaqə)

## Bu chapter nədən bəhs edir?
Persistent storage anlayışına, database/sql paketinə (driver, bağlantı, sorğu emalı), tam Bookstore sxeminə (6 cədvəl), ORM (GORM) alternativinə, indeksləmə strategiyalarına, connection pool tuning-inə və advanced SQL-ə (JOIN, subquery, aggregation).

## Əsas fikirlər

### 1. Persistent storage nədir
**Nədir:** Server instansiyasından/sessiyadan asılı olmayan uzunmüddətli data anbarı — personalizasiya, sifariş tarixçəsi, tövsiyələrin əsası.

**Seçim:** Relational (PostgreSQL/MySQL — strukturlaşdırılmış, mürəkkəb sorğular) vs NoSQL (MongoDB). Bookstore üçün relational: kitab/müəllif/janr/review münasibətləri.

**Scale problemləri:** Milyonlarla kitab + minlərlə eynizamanı user → distributed DB, caching, sharding.

### 2. database/sql — Go-nun DB interfeysi
**Nədir:** Generic interfeys — ORM-in tam abstraksiyası ilə xam SQL tam nəzarəti arasında balans.

**Driver qurulumu (PostgreSQL nümunəsi):**
```bash
go get -u github.com/lib/pq
```
```go
import (
    "database/sql"
    _ "github.com/lib/pq"   // _ = yalnız register üçün, ad sahəsinə bulaşdırmır
)

connStr := "user=username dbname=gitforgits_bookstore password=password host=localhost sslmode=disable"
db, err := sql.Open("postgres", connStr)
if err != nil { log.Fatal(err) }
defer db.Close()
```

**Oxu (Query + Scan):**
```go
rows, err := db.Query("SELECT title, author FROM books")
if err != nil { log.Fatal(err) }
defer rows.Close()                          // sızıntı qarşısı!
for rows.Next() {
    var title, author string
    if err := rows.Scan(&title, &author); err != nil { log.Fatal(err) }
    fmt.Printf("Title: %s, Author: %s\n", title, author)
}
```

**Yazı (Exec — INSERT/UPDATE/DELETE):**
```go
result, err := db.Exec("INSERT INTO books (title, author) VALUES ($1, $2)",
    "Go Programming", "Jane Doe")
```
**Vacib:** `$1, $2` — parameterized query = SQL injection qoruması.

**Transaksiyalar:** `Begin` → əməliyyatlar → `Commit`/`Rollback` — atomiklik.

### 3. Bookstore sxemi (tam DDL)
```sql
CREATE TABLE books (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    author VARCHAR(255) NOT NULL,
    isbn VARCHAR(13),
    description TEXT,
    publication_date DATE,
    genre_id INT,                      -- FK → genres
    price DECIMAL(10,2)
);

CREATE TABLE genres (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL
);

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash CHAR(64) NOT NULL,   -- hash saxlanır, AÇIQ PAROL YOX!
    signup_date DATE DEFAULT CURRENT_DATE
);

CREATE TABLE reviews (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    book_id INT REFERENCES books(id),
    rating INT CHECK (rating >= 1 AND rating <= 5),   -- 1-5 məhdudiyyət
    comment TEXT,
    review_date DATE DEFAULT CURRENT_DATE
);

CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    book_id INT REFERENCES books(id),
    purchase_date DATE DEFAULT CURRENT_DATE,
    amount DECIMAL(10,2)
);

CREATE TABLE inventory (
    book_id INT REFERENCES books(id),
    stock_quantity INT NOT NULL,
    restock_threshold INT,             -- alt hədd xəbərdarlığı
    last_restock_date DATE
);
```
**Dizayn qeydləri:** SERIAL PK, FK ilə əlaqələr (1:N), CHECK constraint (rating), UNIQUE (username/email), DEFAULT (tarixlər), DECIMAL (pulluq üçün float YOX).

### 4. ORM (GORM) nə vaxt?
**Üstünlüklər:** obyekt paradigması (Book/User obyektləri), avtomatik migrasiyalar, assosiasiyalar (1:N, N:M), hooks (pre-save validation), sürətli inkişaf, raw SQL ehtiyatı.
**Trade-off:** xam SQL-in tam nəzarəti itir; mürəkkəb sorğularda öz SQL-ini yazmalı ola bilərsən (GORM icazə verir).

### 5. İndeksləmə strategiyaları
| İndeks növü | Nə edir | Nümunə |
|---|---|---|
| **Primary** | PK üzrə avtomatik | `id` üzrə ani tapıntı |
| **Secondary** | PK-dan başqa sütun | `CREATE INDEX idx_books_title ON books(title)` |
| **Composite** | çoxsütunlu | `ON books(genre_id, publication_date)` |
| **Covering** | sorğunun BÜTÜN sütunları indeksdə | `ON books(title, author)` — cədvələ toxunmur |
| **Partial** | şərtli subset | `ON books(price) WHERE price > 50` |

**Qaydalar:** İndeks = sürətli OXU + yavaş YAZI (hər INSERT index-i də yeniləyir) → sorğu pattern-lərinə görə, ölçülü indekslə.

### 6. Connection pooling
**Nədir:** Əvvəlcədən açıq saxlanan bağlantı növbəsi — taksi parkı metaforası: hər sifariş üçün yeni taksi çağırmaq yox, hazır olanı götür.

**Tuning parametrləri:**
```go
db.SetMaxOpenConns(100)                    // maksimum açıq bağlantı
db.SetMaxIdleConns(50)                     // növbədə saxlanan boş
db.SetConnMaxLifetime(time.Minute * 5)     // maksimum ömür
```
**Health check:** `db.Ping()` — xəta halında: retry, exponential backoff, backup DB-ə yönləndirmə.

### 7. Advanced SQL

**INNER JOIN** — hər iki cədvəldə uyğun gələnlər:
```sql
SELECT books.title, authors.name
FROM books INNER JOIN authors ON books.author_id = authors.id;
```
**LEFT JOIN** — soldan HAMISI, sağdan uyğun gələnlər:
```sql
SELECT genres.name, books.title
FROM genres LEFT JOIN books ON genres.id = books.genre_id;
```
**Subquery** — ən çox kitab yazan müəllif:
```sql
SELECT author_id, COUNT(*) as book_count FROM books
GROUP BY author_id
HAVING COUNT(*) = (
    SELECT MAX(book_count)
    FROM (SELECT author_id, COUNT(*) as book_count FROM books GROUP BY author_id) as counts
);
```
**Aggregations:**
```sql
SELECT COUNT(*) FROM books;                              -- say
SELECT SUM(stock) FROM inventory WHERE book_id = X;      -- cəm
SELECT AVG(price) FROM books;                           -- orta
SELECT MAX(price), MIN(price) FROM books;                -- ekstremum
SELECT author_id, COUNT(*) FROM books GROUP BY author_id;          -- qruplaşdır
SELECT author_id, COUNT(*) FROM books GROUP BY author_id HAVING COUNT(*) > 5;  -- filtr sonrası
```
**KEY fərq:** WHERE qruplaşdırmadan ƏVVƏL, HAVING SONRA filtrləyir.

## Əsas terminlər
- Persistent storage (davamlı saxlama)
- database/sql — generic DB interfeysi
- Driver (sürücü) — lib/pq PostgreSQL üçün
- Connection pooling (bağlantı hovuzu)
- Foreign key (xarici açar) — REFERENCES
- Constraint (məhdudiyyət) — CHECK/UNIQUE/DEFAULT
- ORM (Object-Relational Mapping) — GORM
- Primary/Secondary/Composite/Covering/Partial index
- INNER/LEFT JOIN, Subquery, GROUP BY/HAVING

## Praktik nəticə
1. Driver import-u həmişə `_ "github.com/lib/pq"` — blank import register edir.
2. Parametrli sorğular ($1, $2) — SQL injection-a yeganə peşəkar müdafiə.
3. `rows.Close()` heç vaxt unutma — defer ilə; `password_hash` — heç vaxt açıq parol.
4. İndeksləri sorğu pattern-lərinə görə qur; covering index tez-tez oxunan cütlüklərə; over-indexing yazını yavaşladır.
5. Pool parametrlərini traffic-ə görə tune et + Ping health check + retry/backoff.
6. WHERE vs HAVING fərqini bilmək analytics sorğuların düzgünlüyü deməkdir.

## Mənbə
Pages: 140-162 (PDF səh. 140-162)
