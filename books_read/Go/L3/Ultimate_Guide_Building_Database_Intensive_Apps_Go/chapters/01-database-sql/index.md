# The Ultimate Guide To Building Database-Intensive Apps with Go (səh. 1-39)

**Müəllif:** Baron Schwartz | **Nəşriyyat:** VividCortex | **İl:** 2019 (Go 1.12-yə qədər) | **Səviyyə:** L3 (Intermediate)

## Bu kitab nədən bəhs edir?

Go-nun standart `database/sql` paketinin KAMM dərsləri: bağlantı pool-u (sql.DB
= bağlantı DEYİL, databazadır; lazy bağlantılar; avtomatik 10 retry), pool
konfiqurasiyası (SetMaxOpenConns/SetMaxIdleConns/SetConnMaxLifetime),
Query/QueryRow/Exec fərqləri (bağlantının nə vaxt pool-a qayıtdığı!),
rows.Next/Scan/Close həyat dövrü (leak → "too many connections"), NULL
(parametrlər — sql.NullString kimi), NULL/scan avtomatik çevrilmələri,
bilmədiyim sütunlar (rows.Columns, sql.RawBytes, sql.ColumnType), çoxlu
nəticə dəstləri (NextResultSet 1.8+), prepared statement-lər (db.Prepare,
gizli avtomatik prepare, təkrar prepare bağlantıdaşınması problem, anti-
patternlər), tranzaksiyalar (db.Begin/Commit/Rollback; tx-də TƏK bağlantı —
paralellik YOX; db-yə müraciət = tranzaksiyadan kənar!), sql.Conn (tək
bağlantı zəmanəti), xəta idarəetməsi (rows.Err, driver-specific error
codes — mysql.MySQLError + VividCortex/mysqlerr), driver.Valuer/sql.Scanner
interfeysləri (şəffaf transformasiyalar: lowercase, gzip, ENKRİPSİYA),
db.Stats() monitoring, context (1.8+: QueryContext), driverlər (go-sql-driver/
mysql, lib/pq), 13 common pitfall siyahısı.

## Əsas fikirlər

### 1. sql.DB — Bağlantı DEYİL, Pool-dur
```go
import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"   // side-effect import — init() ilə qeydiyyat
)
db, err := sql.Open("mysql", "root:@tcp(:3306)/test")
defer db.Close()
err = db.Ping()   // Open BAĞLANTI AÇMIR (lazy) — Ping yoxlayır
```
- **sql.DB = database obyekti;** bağlantılar İÇ pool-da; konkurent-safe;
  uzunömürlü — bir dəfə yarat, daimi istifadə et (hər HTTP sorğu üçün YENİ
  db açma — pitfall #2: TIME_WAIT TCP-lər, latency)
- `_` importu: driver yalnız init()-də özünü `sql.Register()` ilə qeyd
  edir — ad sahəsi istifadə olunmur

### 2. Pool Mexanizmi və Konfiqurasiya
- **Əsasənə mənimsəməli qayda:** sorğu → pool-dan boş bağlantı (yoxdursa
  yeni); bitəndə qaytarılır. **Ölü bağlantı → avtomatik retry (max 10) —
  kodda retry yazma!**
- Konfiqurasiya:
  - `SetMaxOpenConns(n)` — aktiv+idle limiti; 0 = limitsiz; dolu → BLOCK
  - `SetMaxIdleConns(n)` — saxlanılan idle sayı; default 0 = bağlantılar
    dərhal QAPADILIR → thrashing (açıb-bağlama təlatümü)
  - `SetConnMaxLifetime(d)` — köhnəlmiş bağlantıların bitmə vaxtı
- **4 problemə diqqət:** thrash; DB-də çox bağlantı; gözləmə bloku; 10+
  ölü → retry limitindən sonra uğursuzluq

### 3. Query vs QueryRow vs Exec — bağlantı sahibliyi FƏRQLİDİR
| Metod | Bağlantı nə vaxt qayıdır |
|---|---|
| db.Ping() | dərhal |
| db.Exec() | dərhal (Result istinad saxlayır) |
| db.Query() | rows tam iterasiya Close() → **qədər tutulur!** |
| db.QueryRow() | Scan() çağırılana qədər |
| db.Begin() | Commit()/Rollback() → |

- **QƏTİ QADAĞA:** `_, err := db.Query("DELETE ...")` — rows IQNORE edilib →
  bağlantı heç vaxt pool-a QAYTIRMIR → LEAK → "too many connections" →
  downtime. Non-SELECT üçün Exec!

### 4. rows Həyat Dövrü
```go
rows, err := db.Query("SELECT ...")
if err != nil { log.Fatal(err) }
defer rows.Close()          // early-return halları üçün
for rows.Next() {           // io.EOF-də false; xətada da false — daxili yaddaşda
    err = rows.Scan(&s1, &s2)
    if err != nil { log.Fatal(err) }
}
if err = rows.Err(); err != nil { log.Fatal(err) }   // DÖNGÜDƏN SONRA MÜTLƏQ!
```
- Normal bitişdə Close avtomatik; break/return → defer YOXDUSA leak
- **Uzun-ömürlü funksiyada defer ETMƏ:** dəyişən gecikir + yaddaş tutulur —
  döngüdən çıxmazdan ƏVVƏL açıq-aşkar Close
- rows.Err() io.EOF-ni özü həll edir (nil qaytarır)

### 5. Scan — Avtomatik Çevrilmələr və NULL
```go
var (s1 string; s2 sql.NullString; i1 int; f1, f2 float64)
// sətir: ["hello", NULL, 12345, "12345.6789", "not-a-float"]
err = rows.Scan(&s1, &s2, &i1, &f1, &f2)
// "12345.6789" avtomatik ParseFloat! "not-a-float" → xəta (qalanları UĞURLU)
// s2: {String: "", Valid: false}
if s2.Valid { use s2.String }
```
- VARCHAR rəqəm → float64 dest-inə: Scan ÖZÜ strconv çağırır — kod təmiz
- **NULL → NullXXX tipləri** (NullString/NullFloat64/...); özəl tip lazımdısa
  — driver-ə bax VƏ YA bir neçə sətir kopyala
- Sıra vacibdir: arqumentlər SIRAYLA işlənir — birincilər uğurlu qalır

### 6. Naməlum Sütunlar + Çoxlu Nəticə Dəstləri
```go
cols, err := rows.Columns()                 // adlar + say
dest := []interface{}{new(uint64), new(string), ...}
err = rows.Scan(dest[:len(cols)])           // tipi bilinəndə
// tam naməlum: vals[i] = new(sql.RawBytes)  (QueryRow ilə İŞLƏMİR!)
// yeni Go: sql.ColumnType — DatabaseTypeName/Length/Nullable/ScanType
rows.NextResultSet()   // 1.8+: stored procedure-lər; false → Err() yoxla
```

### 7. Prepared Statement-lər
```go
stmt, err := db.Prepare("INSERT INTO test.hello(world) VALUES(?)")
defer stmt.Close()
for _, s := range []string{"hello1", "hello2"} {
    res, err := stmt.Exec(s)
}
// gizli: db.Exec("INSERT ... VALUES(?)", "hello") → prepare + exec + close
```
- **Faydaları:** quoting yox; SQL injection müdafiəsi; MySQL binary
  protokolu; təkrar icrada parse xərcləri itir
- **Anti-patternlər:**
  1. Tək-istifadə prepare — 3 raund-trip (prepare/exec/close) = israf
  2. **Döngü DAXİLİNDƏ Prepare** — döngü XARİCİNDƏ!
- **Pool təsiri (GİZLİ):** db.Prepare → hərə hansı bağlantıda prepare →
  pool-a qayıdır; icrada bağlantı məşğuldursa → BAŞQA bağlantıda YENİDƏN
  prepare → statement "çoxalır"; yüksək yükdə server-də statement sayı
  partlaya bilər → "leak"
- **Bypass (plaintext SQL):** driver.Execer/Queryer implement edirsə;
  parametrsiz/təkrarsız/statementsiz DB (Sphinx, MemSQL) — bu halda
  ÖZÜN validasiya (injection!)
- **Parametr sintaksisi (driver-a bağlı):** MySQL `?`; PostgreSQL `$1, $2`
  (təkrar istifadə OLAR); Oracle `:name`; SQLite hər ikisi

### 8. Tranzaksiyalar — sql.Tx
- **YANLIŞ YOL:** db.Exec("BEGIN") / UPDATE / db.Exec("COMMIT") — hər
  Exec FƏRLİ bağlantıda ola bilər!
- **DÜZGÜN:**
```go
tx, err := db.Begin()            // VƏ YA BeginTx(opts{IsolationLevel})
// tx TƏK bağlantı tutur; db ilə ƏLAQƏSİ KƏSİLİR
_, err = tx.Exec("UPDATE ...")   // db.Exec YOX — tx.Exec!
err = tx.Commit()                // VƏ YA tx.Rollback(); sonra tx ETİBARSIZ
```
- **Anti-pattern:** Begin-dən sonra db.Exec çağırmaq — db tranzaksiyada
  DEYİL, FƏRLİ bağlantı!
- **tx daxilində konkurrensiya YOXDUR:** rows iterasiyası zamanı eyni tx-dən
  yeni sorğu = bağlantı məşğul → BOOM; daxili sorğu tranzaksiyaya daxil
  olmayacaqsa — db istifadə et (diqqətli!)
- **tx-də retry qeyri-aktiv:** ölü bağlantıda avtomatik 10 retry YOX —
  deadlock/rollback-ləri ÖZÜN idarə et
- tx.Prepare: yalnız o tx üçün; tx.Stmt(stmt) — db stmt-u tx-a "klonlaşdırır"
  (yenidən prepare)

### 9. sql.Conn (1.9+) — Tranzaksiyasız Tək Bağlantı
- Səbəblər: temp cədvəllər/user dəyişənləri, USE; konkurrentsini məhdudlaşdır;
  explicit lock-lar
- `conn, _ := db.Conn(ctx)` → Close() mütləq; YALNIZ Context-vari imzalar
- Köhnə Go: tx istifadə edib dərhal commit "COMMIT" — tövsiyə EDİLMİR (driver
  xətası mümkün)

### 10. Xəta İdarəetməsi
- **String match YOX:** `strings.Contains(err, "Deadlock")` — dil dəyişəndə
  SİNİR; ANSI kodlar da qeyri-spesifik
- **Düzgün:** driver tipi + xəta kodu:
```go
if driverErr, ok := err.(*mysql.MySQLError); ok {
    if driverErr.Number == mysqlerr.ER_LOCK_DEADLOCK { ... }  // github.com/VividCortex/mysqlerr
}
```
- rows.Err() döngüdən sonra MÜTLƏQ; rows.Close() xətası — logla keç (həll
  yolu YOXDUR)

### 11. Valuer/Scanner — Şəffaf Transformasiya
```go
type LowercaseString string
func (ls LowercaseString) Value() (driver.Value, error) {  // YAZARKƏN
    return driver.Value(strings.ToLower(string(ls))), nil
}
func (ls *LowercaseString) Scan(src interface{}) error {   // OXUYARKƏN
    // string/[]byte → ToLower → *ls
}
```
- Nümunə: "I AM UPPERCASED" NormalString → DB-də olduğu KİMİ; magic tip →
  DB-də KİÇİK; hər ikisi oxunanda KİÇİK
- **Real dünyada:** validasiya; vahid format; gzip (GzippedText []byte —
  Value/Scan ilə sıxılır/açılır); **şəffaf ENKRİPSİYA** (müştəri datası)

### 12. Monitoring + Context + Driverlər
- **db.Stats()** (ucuz, thread-safe): MaxOpenConnections, OpenConnections,
  InUse, Idle, WaitCount, WaitDuration, MaxIdleClosed, MaxLifetimeClosed →
  metrikalara
- **Sorular:** bağlantı açıq/qapadılır YOXSA reuse? prepare təkrar istifadə?
  ən tez/ən yavaş statement-lər? "zibil" ping-lər? sorğu yavaş İCRA olundu
  YOXSA pool gözlədi? (→ VividCortex paket təhlili; OpenCensus wrapper —
  distributed tracing)
- **Context (1.8+):** QueryContext və s. — köhnələr thin wrapper; ləğv/
  timeout/trace; MySQL driver timeout-cancellation DƏSTƏKLƏYİR, lib/pq —
  HƏQİQƏTƏN İCRA ETMİR
- **Driver vəzifələri:** bağlantı (pool YOX — sql özü edir), Rows iterator,
  nəticə interfeysi, prepared statement, tranzaksiya commit/rollback,
  dəyər çevrilmələri; init()-də sql.Register("mysql", &MySQLDriver{})
- **Tövsiyə olunanlar:** MySQL → go-sql-driver/mysql; PostgreSQL → lib/pq

### 13. 13 Common Pitfalls (istirəsiz siyahı)
1. Döngü daxilində defer — yaddaş+bağlantı sınırsız artır
2. Çoxlu db obyekti — global sql.DB; hər sorğu üçün YOX (TIME_WAIT fəlakəti)
3. rows.Close() unutma — leak; mümkün qədər TƏZ Close (təkrarı zərərsizdir)
4. QueryRow().Scan() zənciri — tək sətirdə (leak qarşısı)
5. Tək-istifadə prepare — Sprintf + parametrsiz fikirləş (2 raund-trip qənaət)
6. Prepare bloat — yüksək konkurrentsdə təkrar prepare (çoxaldılma)
7. strconv/kast kodu — Scan avtomatik çevirir; tipi istəyin
8. retry/pool kodu — database/sql özü idarə edir
9. rows.Err() unutma — döngü anormal çıxa bilər
10. db.Query() non-SELECT-də — leak; Exec işlət
11. Eyni bağlantı gümanı — ardıcıl statement-lər FƏRLİ bağlantıda; tək
    bağlantı lazımdırsa tx VƏ YA sql.Conn
12. tx zamanı db istifadə — db tranzaksiyada DEYİL
13. NULL sürprizi — NullXXX-siz NULL scan = xəta; SCHEMANI DİQQƏTLƏ OXU
    (testdə işləyən, prod-da partlayır)
14. uint64 parametr (bonus) — MSB set olanda qəbul EDİLMİR → fmt.Sprint ilə
    string kimi

## Əsas terminlələr
- sql.DB / sql.Rows / sql.Row / sql.Result / sql.Stmt / sql.Tx / sql.Conn — əsas tiplər
- Connection pool (bağlantı hovuzu) — daxili idarəətli bağlantı dəsti
- Lazy connection — ilk real əməliyyatda açılan bağlantı
- Connection leak — pool-a qayıtmayan bağlantı (unudulmuş Close)
- Thrashing — boş idle qadağası səbəbindən açılıb-qapanan bağlantılar
- Side-effect import (`_`) — yalnız init() üçün (driver qeydiyyatı)
- Client-side prepared statement — saxta prepare (Perl DBI); Go HƏQİQƏTƏN server tərəfində
- NullXXX tipləri — NULL-scan üçün {T, Valid} strukturları
- driver.Valuer / sql.Scanner — yazma/oxuma transformasiya interfeysləri
- sql.RawBytes — kopyasız scan (məhdud ömürlü, QueryRow ilə YOX)
- NextResultSet — çoxlu nəticə dəstləri (stored procedures)

## Praktiki nəticə

1. **Bir sql.DB, uzun ömür:** package səviyyəsində; Ping ilə doğrula; sorğu
   başına Open YOX.
2. **Exec vs Query ciddi ayrılıq:** nəticə seti VARSA Query (+Close/Err),
   YOXDURSA Exec — bağlantı leak-inin əsas mənbəyi bu qarışıqlıqdır.
3. **rows trio:** defer/close-early + döngüdən sonra Err() — hər iki uc.
4. **Tranzaksiya qaydaları:** db.Exec("BEGIN") MÜTLƏQ QADAĞA; tx.Exec işlət;
   tx daxilində ardıcllıq; retry-ləri özün yaz.
5. **NULL-ğa hazırlaş:** NullXXX; schemadakı hər NULL sahə bir gün dəyəcək.
6. **Xətaları kodla tut:** driver struct assert + adlandırılmış konstantlar —
   string match YOX.
7. **Valuer/Scanner ilə cross-cutting:** gzip/şifrələmə/validasiya — DB
   qatında ŞƏFFAF, biznes kodu təmiz qalır.

## Mənbə
Pages: 1-39 (PDF 1-40; son 2 səh. VividCortex marketing)
