# Chapter 17 — Relational Databases (səh. 381-413)

## Bu fəsil nədən bəhs edir?

Go ilə SQL: standart `database/sql` paketi (driver-lər, bağlantı pool-u,
Exec/Query, transaksiyalar, prepared statements) və GORM ORM (modellər,
tag-lər, ilişkilər — belongs-to/has-one/has-many/many2many, CRUD, Where,
Preload, transaksiyalar, savepoint-lər).

## Əsas fikirlər

### 1. database/sql — bağlantı
**Nədir:** Standart paket ortaq SQL interfeysi verir; konkret driver-lər
üçüncü tərəfdən import olunur (eyni kod fərqli DB-lərlə işləyir).

**Kitabdan kod nümunəsi (SQLite):**
```go
import (
    "database/sql"
    _ "github.com/mattn/go-sqlite3"   // yalnız side effect — driver qeydiyyatı
)

func main() {
    db, err := sql.Open("sqlite3", "/tmp/example.db")
    if err != nil {
        panic(err)
    }
    defer db.Close()
    db.SetConnMaxLifetime(time.Minute * 3)   // bağlantı ömrü
    db.SetMaxOpenConns(10)                    // maks açıq bağlantı
    db.SetMaxIdleConns(10)                    // maks boş bağlantı

    ctx := context.Background()
    if err := db.PingContext(ctx); err == nil {
        fmt.Println("Database responds")
    } else {
        panic(err)
    }
}
```

**Sub-kod izahı:**
- `_ "driver"` → blank import; driver-i `sql`-ə qeyd edir (Chapter 2.9)
- `sql.Open(driver, dataSource)` → BAĞLANTI YARATMIR — yalnız *DB (pool)
  obyekti; real bağlantı `Ping`/`PingContext` ilə yoxlanılır
- dataSource: sqlite → fayl yolu; MySQL/PG → host/user/db məlumatı
- Pool tənzimləmələri: ömür və bağlantı limitləri

### 2. Dəyişdirici sorğular (Exec)
**Nədir:** CREATE/INSERT/UPDATE/DELETE — sətir qaytarmayan sorğular.
`Exec`/`ExecContext` → `Result` (son ID + təsir edilən sətir sayı).

```go
const (
    USERS_TABLE = `CREATE TABLE users(
    name varchar(250) PRIMARY KEY,
    email varchar(250)
)`
    USERS_INSERT = "INSERT INTO users (name, email) VALUES(?,?)"
)

func create_table(db *sql.DB) {
    ctx := context.Background()
    _, err := db.ExecContext(ctx, USERS_TABLE)
    if err != nil {
        panic(err)
    }
}

func insert_rows(db *sql.DB) {
    ctx := context.Background()
    result, err := db.ExecContext(ctx, USERS_INSERT, "John", "john@gmail.com")
    if err != nil {
        panic(err)
    }
    lastUserId, _ := result.LastInsertId()
    numRows, _ := result.RowsAffected()
    fmt.Printf("Row ID: %d, Rows: %d\n", lastUserId, numRows)
}
```

**Placeholder fərqləri (DB-dən asılı):**
| DB | Sintaksis |
|---|---|
| MySQL | `VALUES(?,?)` |
| PostgreSQL | `VALUES($1,$2)` |
| Oracle | `VALUES(:val1,:val2)` |

### 3. Oxuyan sorğular (Query/Scan)
**Nədir:** SELECT üçün `Query`/`QueryContext` → `*Rows` iterator; hər
sətir `Scan` ilə struct-a ötürülür.

```go
type User struct {
    Name  string
    Email string
}

func get_rows(db *sql.DB) {
    ctx := context.Background()
    rows, err := db.QueryContext(ctx, "SELECT * from users")
    if err != nil {
        panic(err)
    }
    defer rows.Close()          // rows deskriptoru bağlanmalıdır

    for rows.Next() {            // növbəti sətir varmı
        u := User{}
        rows.Scan(&u.Name, &u.Email)   // sıra SELECT sütunları ilə eyni olmalıdır
        fmt.Printf("%v\n", u)
    }
}
```
- Tək sətir üçün `QueryRow` + `Scan` birbaşa işləyir

### 4. Transaksiyalar (standart paket)
**Nədir:** Əməliyyatlar dəsti — hamısı uğurlu olmalı, yoxsa rollback.
`BeginTx` → `Tx` (Exec/Query metodları onun üzərindən).

```go
func runTransaction(db *sql.DB) {
    tx, err := db.BeginTx(context.Background(), nil)
    defer tx.Rollback()          // xətasız halda no-op, xətada geri qaytarır

    ctx := context.Background()
    _, err = tx.ExecContext(ctx, "INSERT INTO users(name, email) VALUES(?,?)",
        "Peter", "peter@email.com")
    if err != nil {
        panic(err)
    }
    err = tx.Commit()             // hamısını icra et
    if err != nil {
        panic(err)
    }
}
```
- Transaksiya bağlantını tutur — uzun sürməməlidir

### 5. Prepared statements
**Nədir:** Sorğu bir dəfə parse edilir, təkrar icralarda parse xərci
 yoxdur.

```go
func runStatement(db *sql.DB) {
    stmt, err := db.Prepare("select * from users where name = ?")
    if err != nil {
        panic(err)
    }
    defer stmt.Close()            // bağlantıya bağlı — bağlanmalıdır

    result := stmt.QueryRow("Peter")
    u := User{}
    err = result.Scan(&u.Name, &u.Email)
    switch {
    case err == sql.ErrNoRows:
        panic(err)
    case err != nil:
        panic(err)
    default:
        fmt.Printf("%v\n", u)
    }
}
```
- `sql.ErrNoRows` → nəticə boşdur (xəta deyil, vəziyyət)

### 6. GORM — əsaslar
**Nədir:** ORM — struct əməliyyatlarını SQL-ə çevirir (MySQL, PostgreSQL,
SQLite, SQL Server, Clickhouse).

**Kitabdan kod nümunəsi (tam CRUD):**
```go
import (
    "gorm.io/gorm"
    "gorm.io/driver/sqlite"
)

type User struct {
    ID    uint
    Name  string
    Email string
}

func main() {
    db, err := gorm.Open(sqlite.Open("/tmp/test.db"), &gorm.Config{})
    if err != nil {
        panic("failed to connect database")
    }

    err = db.AutoMigrate(&User{})     // struct → cədvəl (users)
    if err != nil {
        panic(err)
    }

    u := User{Name: "John", Email: "john@gmail.com"}
    db.Create(&u)                       // INSERT — ID avtomatik dolur

    var recovered User
    db.First(&recovered, "name=?", "John")   // SELECT (SQL şərti ilə)
    fmt.Println("Recovered", recovered)

    db.Model(&recovered).Update("Email", "newemail")   // UPDATE
    db.First(&recovered, 1)                            // primary key ilə

    db.Delete(&recovered, 1)                           // DELETE
}
```
- `AutoMigrate(&T{})` → cədvəli qurur/güncəlləyir
- `ID` sahəsi default olaraq primary key-dir

### 7. Modellər və tag-lər
```go
// gorm.Model — hazır metadata struct:
type Model struct {
    ID        uint           `gorm:"primaryKey"`
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt `gorm:"index"`
}
```

**Embedded struct (Scanner/Valuer ilə):**
```go
type User struct {
    Name  string
    Email string
}

func (u *User) Scan(src interface{}) error {     // DB → struct
    input := src.([]byte)
    json.Unmarshal(input, u)
    return nil
}
func (u User) Value() (driver.Value, error) {    // struct → DB
    enc, err := json.Marshal(u)
    return enc, err
}

type Operator struct {
    ID         uint
    User       User  `gorm:"embedded,embeddedPrefix:user_"`
    Platform   string `gorm:"not null"`
    Dedication uint  `gorm:"check:dedication>5"`   // çek məhdudiyyəti
}
```
- `embedded` → struct dəyəri bir sütun dəstinə seriyalaşır (JSON ilə)
- Tag-lər: `not null`, `check:...`, `index`, `primaryKey`...

### 8. İlişkilər (relationships)
**Belongs-to (1:1 — FK sahib tərəfdə):**
```go
type User struct {
    ID      uint
    Name    string
    Email   string
    GroupID uint       // FK
    Group   Group
}
type Group struct {
    ID   uint
    Name string
}
// Xüsusi FK: Group Group `gorm:"foreignKey:Ref"` + Ref uint
```

**Has-one (1:1 — FK digər tərəfdə):**
```go
type User struct {
    ID     uint
    Laptop Laptop      // FK Laptop-da (UserID)
}
type Laptop struct {
    ID           uint
    SerialNumber string
    UserID       uint
}
```

**Has-many (1:N):**
```go
type User struct {
    ID      uint
    Laptops []Laptop    // laptops.user_id FK
}
```

**Many-to-many (join table):**
```go
type User struct {
    ID        uint
    Languages []Language `gorm:"many2many:user_languages"`
    Laptops   []*Laptop   `gorm:"many2many:user_laptops"`
}
type Laptop struct {
    ID   uint
    Users []*User `gorm:"many2many:user_laptops"`
}
```
- `many2many:ad` → join cədvəlinin adı
- Pointer slice-lar (`[]*T`) hər iki tərəfdən əlaqə üçün

### 9. CRUD detalları

**Create + CreateInBatches:**
```go
u := User{Name: "John"}
res := db.Create(&u)                      // ID avtomatik set olunur
fmt.Printf("User ID: %d, rows: %d\n", u.ID, res.RowsAffected)

users := []User{{Name: "Peter"}, {Name: "Mary"}}
db.CreateInBatches(users, 5)              // 5-lük batch-lərlə
```

**Query — First/Take/Last/Find/Where:**
```go
var u User
db.First(&u)              // ilk (PK sırası ilə)
u = User{}
db.Take(&u)               // sırasız tək
u = User{}
db.Last(&u)              // son
u = User{}
db.First(&u, 2)          // ID=2

var retrievedUsers []User
db.Find(&retrievedUsers)                  // hamısı
db.Find(&retrievedUsers, []int{2, 4})     // PK siyahısı ilə

db.Where("name = ?", "Jeremy").First(&u)
db.Where("name LIKE ?", "%J%").Find(&retrievedUsers)
db.Where("name LIKE ?", "%J%").Or("name LIKE ?", "%y").Find(&retrievedUsers)
db.Where(&User{Name: "Mary"}).First(&u)   // struct arqument kimi!
db.Order("name asc").Find(&retrievedUsers)
```
- **Tələ:** hər sorğudan əvvəl `u = User{}` sıfırlama — yoxsa köhnə
  sahələr sorğu şərti kimi oxunur
- `Where` → SQL şərti + arqumentlər; `Or` → və ya

### 10. Eager loading (Preload)
**Nədir:** Əlaqəli cədvəllərin avtomatik yüklənməsi — başqa cədvəl üçün
əlavə sorğu.

```go
db.Preload("Laptops").First(&recovered)   // çoxlu: plural ad!
// Recovered with preload {1 John john@gmail.com [{1 sn0000001 1} ...]}
// preload-sız Laptops boş [] olardı

db.Preload("Group").First(&recovered, 1)  // tək: singular (belongs-to)
```
- `Preload("Laptops")` vs `Preload("Group")` — sahə adının CƏM/tək
  olmasına diqqət

### 11. GORM transaksiyaları

**Funksiya əsaslı (avtomatik rollback):**
```go
db.Transaction(func(tx *gorm.DB) error {
    if err := tx.Model(&u).Update("Email", "newemail").Error; err != nil {
        return err          // return ≠ nil → rollback
    }
    var inside User
    tx.First(&inside)
    if err := tx.Create(&u).Error; err != nil {
        return err          // UNIQUE constraint → rollback
    }
    return nil              // commit
})
```
- Default: hər yazma əməliyyatı öz transaksiyasında; `SkipDefaultTransaction`
  ilə söndürülə bilər (performans)

**Əl ilə idarə + SavePoint:**
```go
func RunTransaction(u *User, db *gorm.DB) error {
    tx := db.Begin()
    if tx.Error != nil {
        return tx.Error
    }
    if err := tx.Model(u).Update("Email", "newemail").Error; err != nil {
        return err
    }
    tx.SavePoint("savepoint")            // arxaq nöqtə
    if err := tx.Create(u).Error; err != nil {
        tx.RollbackTo("savepoint")       // yalnız buraya qədər geri
    }
    return tx.Commit().Error
}
// Nəticə: e-poçt yenilənməsi QORUNUR, uğursuz create isə rollback olunur
```

## Əsas terminlər
- database/sql — standart SQL interfeys paketi
- Driver — konkret DB adapteri (blank import ilə qeydiyyat)
- Connection pool — bağlantıların idarə olunan dəsti
- Exec vs Query — dəyişdirici vs oxuyan sorğu
- Prepared statement — bir dəfə parse edilən sorğu şablonu
- ORM — obyekt-relasiya xəritələmə
- AutoMigrate — struct-dan cədvəl qurma
- Preload (eager loading) — əlaqəli datanın ilkin yüklənməsi
- Savepoint — transaksiya daxilində qismən rollback nöqtəsi

## Praktik nəticə
İki qat seçimi: tam nəzarət üçün `database/sql` (driver blank import, pool
tənzimi, Exec/Query + Scan, defer Rollback + Commit); sürətli inkişaf üçün
GORM (`AutoMigrate` + struct tag-lər + `Preload`). GORM-da sorğu dəyişənlərini
istifadədən əvvəl sıfırlayın; əlaqələrdə FK yerini və many2many join adını
tag-lərlə idarə edin; xətalı ardıcıllıqda savepoint qismən bərpa verir.

## Mənbə
Pages: 381-413 (PDF 381-413)
