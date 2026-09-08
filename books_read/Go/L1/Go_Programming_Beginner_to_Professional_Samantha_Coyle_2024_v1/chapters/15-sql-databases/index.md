# Chapter 15 — SQL and Databases (səh. 470-503)

## Bu fəsil nədən bəhs edir?

PostgreSQL + database/sql: bağlantı (Open/Ping/Close), cədvəl yaratma
(CREATE TABLE), SQL injection qorunması (Prepare + $1/$2), data əlavə,
oxu (Query/QueryRow + Scan), UPDATE/DELETE (Exec + RowsAffected),
TRUNCATE/DROP və GORM ORM (AutoMigrate, Create, First/Find, tx.Error).

## Əsas fikirlər

### 1. Bağlantı
**Tələblər:** host, port, dbname, user, password (+ səlahiyyətlər).

```go
package main

import "database/sql"
import _ "github.com/lib/pq"     // yalnız side effect (driver qeydiyyatı)!

db, err := sql.Open("postgres",
    "user=postgres password=Start!123 host=127.0.0.1 port=5432 dbname=postgres sslmode=disable")
if err != nil {
    panic(err)
}

// Bağlantı Sınaması (Open bağlantı AÇMIR — yalnız konfiqurasiya!):
if err := db.Ping(); err != nil {
    panic(err)
}

defer db.Close()                  // idiomatik bağlanma
```
- API + driver arxitekturası: database/sql = vahid interfeys; driver
  dəyişilsə kod DƏYİŞMİR (MySQL→Postgres keçidi sərfəli)
- Ping — long-running servislərdə periodic yoxlama (goroutine-də)

### 2. Cədvəl yaratma
```sql
CREATE TABLE table_name (
    column1 datatype constrain,
    column2 datatype constrain
);
```
- Tiplər: INT, DOUBLE, FLOAT, VARCHAR
- Məhdudiyyətlər: NOT NULL, PRIMARY KEY, adlı funksiya
- ACL vs İrs+Rollar (Postgres: irs+rollar yanaşması)

```go
DBCreate := `
CREATE TABLE public.test (
    id integer,
    name character varying COLLATE pg_catalog."default"
)
WITH (
    OIDS = FALSE
)
`
_, err = db.Exec(DBCreate)
if err != nil {
    panic(err)
}
```
- Təkrar icra → "table already exists" xətası (gözlənilən)
- COLLATE — sıralama qaydası; oid — Postgres avtomatik sətir ID-si
  (biz id-i əl ilə idarə edirik)

### 3. SQL injection + Prepare
**Hücum:** `WHERE username=<input> OR '1'='1'` — həmişə true!

**Qoruma:** Prepare + parametr avazlığı:
- Sorğularda: `WHERE col = $1`
- Əlavə/dəyişmədə: `VALUES($1, $2)`

```go
insert, err := db.Prepare("INSERT INTO test VALUES ($1, $2)")
if err != nil {
    panic(err)
}
defer insert.Close()

_, err = insert.Exec(2, "second")   // $1=2, $2="second"
```

### 4. Data oxuma
**Query — çoxsətir; QueryRow — maksimum 1 sətir.**

```go
// Çoxsətir:
var id int
var name string

rows, err := db.Query("SELECT * FROM test")
if err != nil {
    panic(err)
}
for rows.Next() {
    err := rows.Scan(&id, &name)          // sütunları dəyişənlərə köçür
    if err != nil {
        panic(err)
    }
    fmt.Printf("Retrieved data from db: %d %s\n", id, name)
}
err = rows.Err()                            // iterasiya xətaları
if err != nil {
    panic(err)
}
rows.Close()
db.Close()
```

**QueryRow + Prepare:**
```go
var name string
id := 2

qryrow, err := db.Prepare("SELECT name FROM test WHERE id=$1")
if err != nil {
    panic(err)
}
err = qryrow.QueryRow(id).Scan(&name)
if err != nil {
    panic(err)
}
fmt.Printf("The name with id %d is %s", id, name)
qryrow.Close()
```
- `SELECT *` production-da adətkan deyil (performans/təhlükəsizlik)

### 5. UPDATE / DELETE (Exec + RowsAffected)
```go
// UPDATE:
UpdateStatement := `
UPDATE test
SET name = $1
WHERE id = $2
`
updateResult, err := db.Exec(UpdateStatement, "well", 2)
if err != nil {
    panic(err)
}
updatedRecords, err := updateResult.RowsAffected()   // neçə sətir?
if err != nil {
    panic(err)
}
fmt.Println("Number of records updated: ", updatedRecords)

// DELETE — eyni pattern:
DeleteStatement := `
DELETE FROM test
WHERE id = $1
`
deleteResult, err := db.Exec(DeleteStatement, 2)
deletedRecords, _ := deleteResult.RowsAffected()
```
- `Exec()` — universal icraçı (SELECT/UPDATE/DELETE hər şeyi)
- RowsAffected — uğurun KVANTİTATİV təsdiqi

### 6. TRUNCATE / DROP
```go
db.Exec("TRUNCATE TABLE test")   // boşalt (cədvəl qalır)
db.Exec("DROP TABLE test")         // TAMAMEN sil
```

### 7. GORM ORM
**Nədir:** cədvəlləri Go struct-ları kimi göstərən obyekt-relyasiya
xəritələyicisi — SQL əvəzinə saf Go kodu.

```go
import (
    "gorm.io/gorm"
    "gorm.io/driver/postgres"
)

type User struct {
    gorm.Model                  // embed: ID (PK), CreatedAt, UpdatedAt...
    FirstName string
    LastName  string
    Email     string
}

// Bağlantı:
db, err := gorm.Open(postgres.Open(connection_string), &gorm.Config{})
if err != nil {
    panic("failed to connect database")
}

// CƏDVƏL YARAT — SQL YAZMADAN:
db.AutoMigrate(&User{})           // users cədvəli avtomatik

// Əlavə:
u := &User{FirstName: "John", LastName: "Smith", Email: "john.smith@gmail.com"}
db.Create(u)

// Oxu — ID ilə:
var user User
db.First(&user, 1)

// Şərtlə (snake_case column adları!):
db.First(&user, "last_name = ?", "Doe")
db.First(&user, "last_name = ? AND first_name = ?", "Smith", "James")

// Struct ilə axtarış (ən sadə):
db.First(&user, &User{FirstName: "James", LastName: "Smith"})

// ÇOX nəticə — Find:
var users []User                       // SLİCE lazımdır!
db.Find(&users, &User{LastName: "Smith"})

// Xətalar — tx.Error:
tx := db.Find(&users, &User{LastName: "Smith"})
if tx.Error != nil {
    fmt.Println(tx.Error)
}
```
- camelCase sahə → snake_case column avtomatik
- tx = transaction obyekti — hər çağırış qaytarır

## Activity icmalları
- **15.01:** Users cədvəli (ID/Name/Email) + struct-dan loop ilə insert
  + email update + user delete — hamısı Prepare ilə
- **15.02:** Messages cədvəli (UserID/Message 280 xarakter) + klaviatura
  girişindən axtarış (yalnış user üçün xəta mesajı)
- **15.02-nin Exercise forması (Exercise 15.01/15.02):** Numbers cədvəli
  (0-99, Odd/Even property) + big.NewInt().ProbablyPrime(0) ilə sadə
  ədədlərin tapılması + cəmi + cüt ədədlərin silinməsi + qalanların
  yenilənməsi

## Əsas terminlər
- SQL / relational database
- database/sql + driver (lib/pq) — vahid API arxitekturası
- Blank import (__) — yalnız side effect üçün
- sql.Open / Ping / defer Close
- Connection string — user/password/host/port/dbname/sslmode
- CREATE TABLE / TRUNCATE / DROP
- SQL injection — OR '1'='1 hücumu
- Prepare + $1/$2 — parametrli sorğu (injection qoruması)
- Query / QueryRow / rows.Next() / Scan
- Exec + RowsAffected — universal icra + təsir sayı
- ORM / GORM — struct↔cədvəl xəritələmə
- gorm.Model — ID + zaman damğaları embed
- AutoMigrate / Create / First / Find
- snake_case transformasiyası — FirstName → first_name
- tx.Error — transaction xətası

## Praktik nəticə
database/sql + driver pattern-i: Open (bağlantı string) → Ping (yoxla)
→ defer Close. Bütün istifadəçi datalı sorğular mütləq Prepare + $N
ilə — SQL injection real təhlükədir. Çoxsətir: Query + Next/Scan/Close;
təksətir: QueryRow. Dəyişikliklər: Exec + RowsAffected (uğur sayını
təsdiqlə). Saf Go yolunu istəyirsənsə GORM: gorm.Model embed +
AutoMigrate + Create/First/Find; xətalar tx.Error ilə.

## Mənbə
Pages: 470-503 (PDF 470-503)
