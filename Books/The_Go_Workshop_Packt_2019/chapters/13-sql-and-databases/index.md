# Chapter 13 — SQL and Databases (SQL və Məlumat Bazaları)

## Bu fəsil nədən bəhs edir?

PostgreSQL qurulumu, database/sql API + driver arxitekturası (unified interface),
bağlantı (sql.Open, connection string, Ping, defer Close), cədvəl yaratma (CREATE
TABLE, datatype/constraint), data əlavəsi (Prepare + $1/$2 — SQL injection qorunması),
sorğu (Query + Scan, QueryRow), update/delete (Exec + RowsAffected), TRUNCATE/DROP
və prime numbers + user/message activity-ləri.

## Əsas fikirlər

### 1. API + Driver Arxitekturası
**database/sql** = unified API; driver-lər = konkret DB implementasiyaları (lib/pq
Postgres üçün, MySQL, Oracle, ODBC...).

**Niyə API-yə yazmalı:** MySQL-ə xüsusi yazılmış script-lər DB dəyişəndə (məs. AWS
Athena) YENİDƏN yazılmalı; API-yə yazılanlar yalnız DRİVER dəyişməklə işləyir.

**DB növləri:** Relational (SQL), NoSQL, Search/Analytic.

### 2. Bağlantı
```go
import "database/sql"
import _ "github.com/lib/pq"     // BLANK import — yalnız side-effect
                                   // (driver-ı database/sql-a QEYDİYYATDAN keçirir)

db, err := sql.Open("postgres",
    "user=postgres password=Start!123 host=127.0.0.1 port=5432 dbname=postgres sslmode=disable")
if err != nil { panic(err) }

connectivity := db.Ping()         // real bağlantı YOXLAMASI (Open yalnız init edir!)
if connectivity != nil { panic(err) }

defer db.Close()                  // idiomatik bağlanma; resurs azadlığı
```
**Open vs Ping:** Open connection string-i VALIDATE edir, bağlantı AÇMIR; Ping real
şəbəkə yoxlamasıdır — long-running app-lərdə kritik.

### 3. Cədvəl Yaratma
```sql
CREATE TABLE table_name (
    column1 datatype constrain,
    column2 datatype constrain
);
```
**Datatype-lər:** INT, DOUBLE, FLOAT, VARCHAR.
**Constraint-lər:** NOT NULL, PRIMARY KEY, adlı funksiya (insert/update-də işə düşür).
**ACL vs Rollar:** Postgres rol-inheritance modeli istifadə edir.

**Kitabdan kod nümunəsi:**
```go
TableCreate := `
CREATE TABLE Number
(
    Number integer NOT NULL,
    Property text COLLATE pg_catalog."default" NOT NULL
)
WITH (OIDS = FALSE)
TABLESPACE pg_default;`

_, err = db.Exec(TableCreate)          // cədvəl yaratma Exec ilə
if err != nil { panic(err) }
```
Təkrar işə salma → "table already exists" xətası (gözlənilən).

### 4. SQL Injection və Prepare
**Hücum nümunəsi:**
```sql
SELECT password FROM Auth WHERE username=<input> OR '1'='1'
-- OR '1'='1' həmişə TRUE → bütün parol hash-ləri sızır!
```
**Qorunma — Prepare + placeholder:**
```go
insert, err := db.Prepare("INSERT INTO test VALUES ($1, $2)")
if err != nil { panic(err) }
_, err = insert.Exec(2, "second")       // dəyərlər AYRI ötürülür — təhlükəsiz
```
**Placeholder formaları:** WHERE-də `$1`; INSERT/UPDATE-də `VALUES($1,$2)`.
**Prinsip:** istifadəçi input-u HƏMİŞƏ placeholder ilə — string concat ilə YOX.

### 5. Sorğu — Query/QueryRow
**Query (çoxsətir):**
```go
var id int
var name string
rows, err := db.Query("SELECT * FROM test")
if err != nil { panic(err) }
for rows.Next() {                       // sətir-sətir iterasiya
    err := rows.Scan(&id, &name)        // & ilə dəyişənlərə DOLDUR
    if err != nil { panic(err) }
    fmt.Println(id, name)
}
err = rows.Err()                        // iterasiya SONRAKI xətalar
rows.Close()
```

**QueryRow (tək sətir) + Prepare:**
```go
qryrow, err := db.Prepare("SELECT name FROM test WHERE id=$1")
if err != nil { panic(err) }
err = qryrow.QueryRow(2).Scan(&name)    // tək sətir + dərhal Scan
qryrow.Close()
```

### 6. Update/Delete — Exec
```go
UpdateStatement := `
UPDATE test
SET name = $1
WHERE id = $2
`
UpdateResult, err := db.Exec(UpdateStatement, "well", 2)
UpdatedRecords, _ := UpdateResult.RowsAffected()   // N sətir dəyişdi?
fmt.Println("Number of records updated:", UpdatedRecords)

// DELETE — eyni pattern:
DeleteStatement := `DELETE FROM test WHERE id = $1`
DeleteResult, _ := db.Exec(DeleteStatement, 2)
DeletedRecords, _ := DeleteResult.RowsAffected()
```
**RowsAffected** — əməliyyatın real işə düşdüyünün TƏSDİQİ (0 = heç nə dəyişmədi!).

### 7. TRUNCATE və DROP
```go
db.Exec("TRUNCATE TABLE test")     // bütün məzmunu BOŞALT (cədvəl qalır)
db.Exec("DROP TABLE test")         // cədvəli TAM YOX ET
```

### 8. Prime Numbers Exercise (kitabdan)
```go
import "math/big"

// Bütün rəqəmləri oxu, prime-ləri tap:
for result.Next() {
    result.Scan(&number, &prop)
    if big.NewInt(number).ProbablyPrime(0) {   // primality test
        primeSum += number
    }
}
// Even-ləri sil:
db.Exec("DELETE FROM Number WHERE Property=$1", "Even")
// Qalan tək rəqəmlərə primeSum əlavə et:
db.Exec("UPDATE Number SET Number=$1 WHERE Number=$2 AND Property=$3",
    newNumber, number, prop)
```
Full pipeline: SELECT → analiz → DELETE → UPDATE.

### 9. User/Messages Activity-lər
Users cədvəli (ID PK, Name, Email — NULL YOX); struct + for loop ilə insert; email
update + ikinci istifadəçi delete; hamısı Prepare ilə. Messages cədvəri (UserID,
280-char Message); userID-siz mesajı sil; klaviaturadan username → Prepare-lə sorğu;
boş nəticə → "no such user".

## Əsas terminlələr
- database/sql — unified DB API; driver- müstəqil
- Blank Import — `_ "pkg"` — side-effect (driver qeydiyyatı)
- Connection String — user/password/host/port/dbname/sslmode
- db.Ping — real connectivity yoxlaması
- defer db.Close() — funksiya sonunda bağlantı azadlığı
- SQL — Structured Query Language (standart)
- Constraint — NOT NULL, PRIMARY KEY
- SQL Injection — input-la sorğunun dəyişdirilməsi hücumu
- Prepare + $1/$2 — placeholder qorunması
- Query/QueryRow — çox/tək sətir sorğu
- rows.Next + Scan — nəticə iterasiyası + & doldurma
- Exec — universal icraçı (CREATE/INSERT/UPDATE/DELETE/TRUNCATE/DROP)
- RowsAffected — təsir edilmiş sətir sayı
- TRUNCATE vs DROP — məzmunu boşalt / cədvəli yox et
- ACL vs Roles — icazə modelləri (Postgres = rollar)

## Praktik nətidə

(1) Kod database/sql API-na yaz — DB dəyişsə yalnız driver swap lazımdır. (2) Blank
import driver-ı qeydə alır — görməzsən də MÜTLƏQ. (3) Open = init; Ping = real
yoxlama — connection problemində Ping xətası verir. (4) İstifadəçi input-u HEÇ VAXT
string concat ilə SQL-ə girmez — Prepare + $1/$2. (5) Scan funksiyası POINTER gözləyir
(&id). (6) rows.Err() iterasiyadan SONRA yoxlanır — loop içindəki xətaları da tutur.
(7) RowsAffected = əməliyyat təsirinin təsdiqi; 0 = şərtə uyğun sətir TAPILMADI.
(8) Exec universal: UPDATE/DELETE/TRUNCATE/DROP hamısı eyni qaydada. (9) Loop içində
çoxsaylı insert-də Prepare-i BİR dəfə yarat, Exec-i təkrarla — sonra Close. (10) DB
müddətli baqlantı resursudur — iş bitən kimi (defer) Close et.

## Mənbə
Pages: 465-491 (PDF 498-525)
