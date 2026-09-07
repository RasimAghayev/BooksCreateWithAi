# database/sql (Ultimate Guide) — Cheat Sheet (Azərbaycanca)

## 1. Qoşulma + konfiqurasiya

```go
import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"   // side-effect: init() → Register
)
db, err := sql.Open("mysql", "root:@tcp(:3306)/test")  // LAZY — açmır!
defer db.Close()
if err := db.Ping(); err != nil { ... }                  // doğrulama

db.SetMaxOpenConns(n)       // aktiv+idle limit; 0=limitsiz; dolu → BLOCK
db.SetMaxIdleConns(n)       // default 0 = thrash təhlükəsi!
db.SetConnMaxLifetime(d)    // köhnəlmiş bağlantıların ömrü
```

## 2. Sorğu matrisası — bağlantı sahibliyi

| Metod | Nəticə | Bağlantı nə vaxt qayıdır | İstifadə |
|---|---|---|---|
| db.Exec | sql.Result | DƏRHAL | INSERT/UPDATE/DELETE/DDL |
| db.Query | *sql.Rows | rows.Close() / tam iterasiya | SELECT (Close MÜTLƏQ!) |
| db.QueryRow | *sql.Row | .Scan() | tək sətir |
| db.Begin | *sql.Tx | Commit()/Rollback() | tranzaksiya |

## 3. rows üçlüyü (kanonik şablon)

```go
rows, err := db.Query("SELECT ...")
if err != nil { ... }
defer rows.Close()              // uzun funksiyada DEYİL — erkən Close!
for rows.Next() {
    err = rows.Scan(&a, &b)
    if err != nil { ... }
}
if err := rows.Err(); err != nil { ... }   // döngüdən SONRA MÜTLƏQ
```

## 4. Tək sətir

```go
var s string
err := db.QueryRow("SELECT ... LIMIT 1").Scan(&s)
if errors.Is(err, sql.ErrNoRows) { /* yoxdur */ }
```

## 5. NULL + avtomatik çevrilmə

```go
var s sql.NullString          // {String string; Valid bool}
rows.Scan(&s)
if s.Valid { use(s.String) }

var f float64
rows.Scan(&f)                  // VARCHAR "12345.6789" → avtomatik ParseFloat
```

## 6. Naməlum sütunlar

```go
cols, _ := rows.Columns()
dest := []interface{}{new(string), new(int), ...}
rows.Scan(dest[:len(cols)])                    // tipləri bilinir
vals[i] = new(sql.RawBytes)                    // tam naməlum (QueryRow YOX)
// ColumnType: Name/Length/Nullable/ScanType/DatabaseTypeName
rows.NextResultSet()                           // stored procedures (1.8+)
```

## 7. Prepared statements

```go
stmt, err := db.Prepare("INSERT INTO t(v) VALUES(?)")  // DÖNGÜ XARİCİNDƏ!
defer stmt.Close()
for _, v := range values {
    stmt.Exec(v)
}
// Tək-istifadə: db.Exec("INSERT ... VALUES(?)", v) gizli 3-raund-trip edir
// Plaintext bypass: driver Execer/Queryer + ÖZÜN validasiya (injection!)
```
- MySQL `?` / PostgreSQL `$1,$2` (təkrar istifadə OK) / Oracle `:name` / SQLite hər ikisi
- sql.NamedArg (1.8+) — adlı parametrlər

## 8. Tranzaksiyalar

```go
tx, err := db.Begin()                    // VƏ YA db.BeginTx(ctx, opts)
_, err = tx.Exec("UPDATE ...")           // db.Exec YOX!!! (başqa bağlantı)
if err != nil { tx.Rollback(); return }
err = tx.Commit()                        // VƏ YA Rollback; sonra tx ölü
```
- **QADAĞALAR:** tx daxilində paralel sorğu (tək bağlantı!); tx-dən kənara
  db istifadə; gizli retry YOX → deadlock özün idarə et
- tx.Stmt(dbStmt) — db stmt-unun tx-scope "klonu" (yeniden prepare)
- Isolation: BeginTx + sql.IsolationLevel option

## 9. sql.Conn (tək bağlantı, tranzaksiyasız)

```go
conn, err := db.Conn(ctx)
defer conn.Close()
conn.QueryContext(ctx, "...")    // YALNIZ Context-vari metoddar
```
- Temp cədvəllər, USE, lock-lar, bağlantı-state.

## 10. Xəta idarəetməsi

```go
// PİS: strings.Contains(err.Error(), "Deadlock")  ← dil asılı!
// YAXŞI:
if driverErr, ok := err.(*mysql.MySQLError); ok {
    if driverErr.Number == mysqlerr.ER_LOCK_DEADLOCK { ... }  // VividCortex/mysqlerr
}
```

## 11. Valuer / Scanner (şəffaf transformasiya)

```go
type GzippedText []byte
func (g GzippedText) Value() (driver.Value, error) { /* gzip */ }
func (g *GzippedText) Scan(src interface{}) error  { /* gunzip */ }
// Insert-də və Scan-də adi []byte kimi istifadə — arxada sıxılır/açılır
```

## 12. Monitoring

```go
st := db.Stats()  // MaxOpen/Open/InUse/Idle/WaitCount/WaitDuration/
                  // MaxIdleClosed/MaxLifetimeClosed — metriyaya göndər
```
- WaitDuration yüksək = sorğu deyil, POOL gözləyir (limitləri artır/düzəlt)
- OpenCensus database/sql wrapper — distributed tracing; VividCortex —
  şəbəkə səviyyəsində sorğu təhlili

## 13. Pitfall-siyahısı (istirəsiz)

| # | Pitfall | Həll |
|---|---|---|
| 1 | Döngüdə defer | döngüdən əvvəl Close/dərhal |
| 2 | Hər sorğuda yeni sql.Open | QLOBAL db, uzun ömür |
| 3 | rows.Close unut | leak; erkən+təkrar Close zərərsiz |
| 4 | QueryRow-Scan ayır | zəncir bir sətirdə |
| 5 | Tək-istifadə prepare | Sprintf+validasiya ilə plaintext |
| 6 | Yüksək konkurents prepare | bloat; plaintext düşün |
| 7 | strconv/kast kodu | Scan avtomatik çevirir |
| 8 | Öz retry/pool kodu | database/sql 10x retry var |
| 9 | rows.Err unut | hər döngüdən sonra |
| 10 | Non-SELECT-də Query | Exec (yoxsa leak) |
| 11 | Eyni-bağlantı gümanı | tx VƏ YA sql.Conn |
| 12 | tx içində db | db tranzaksiyada DEYİL |
| 13 | NULL scan | NullXXX; schema audit |
| 14 | uint64 MSB | fmt.Sprint() → string |

## 14. Driverlər

- MySQL: github.com/go-sql-driver/mysql (context timeout/cancel DƏSTƏKLİ)
- PostgreSQL: github.com/lib/pq (context dəstəyi əslində YOX; NullBool bonusu)
- siyahı: Go wiki drivers; `sql.Drivers()` — yüklü sürücülər
