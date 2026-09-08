# Chapter 18 — NoSQL Databases (səh. 414-432)

## Bu fəsil nədən bəhs edir?

Apache Cassandra NoSQL bazası və GoCQL driver-i: anlayışlar (keyspace,
partition, CQL), bağlantı/sessiya, Cassandra-ya xas data modelling (has-one,
has-many, many-to-many — əlaqəli DB-dən fərqli), CRUD, batch-lər, sorğular
(Scan/Iter/Scanner) və UDT struct-larla işləmə.

## Əsas fikirlər

### 1. Cassandra anlayışları
**Nədir:** Yüksək performans/arabbütlük tələb edən paylanmış NoSQL.
SQL-ə bənzər **CQL** dili var.

| Anlayış | Təsvir |
|---|---|
| Keyspace | cədvəlləri saxlayır; replikasiya siyasəti kimi xüsusiyyətlər təyin edir |
| Table | yalnız sxemanı təyin edir; bir neçə node üzrə partitionlarla instansiya olunur |
| Partition | primary key-in məcburi hissəsi — bütün sətirlərdə olmalıdır |
| Row | unikal primary key (+optional clustering key) ilə sütun toplusu |
| Column | tipli verilən; sxemaya bağlı |

### 2. GoCQL bağlantısı (Session)
**Kitabdan kod nümunəsi:**
```go
import "github.com/gocql/gocql"

const CREATE_TABLE = `CREATE TABLE example.user (
 id int PRIMARY KEY,
 name text,
 email text
)
`

func main() {
    cluster := gocql.NewCluster("127.0.0.1:9042")
    cluster.Keyspace = "example"
    cluster.Consistency = gocql.Quorum     // consistency səviyyəsi
    cluster.Timeout = time.Minute          // bağlantı timeout-u
    session, err := cluster.CreateSession()
    if err != nil {
        panic(err)
    }
    defer session.Close()                  // resurs buraxılır

    ctx := context.Background()
    err = session.Query(CREATE_TABLE).WithContext(ctx).Exec()
    if err != nil {
        panic(err)
    }
}
```

**Sub-kod izahı:**
- `gocql.NewCluster(host)` → ClusterConfig; bir host versək, qalan node-lar
  avtomatik kəşf olunur
- `CreateSession()` → sorğu bağlantıları üçün Session
- `session.Query(CQL).WithContext(ctx).Exec()` → dəyişdirici əməliyyat

### 3. Cassandra data modelling (dəyişik düşüncə!)
**Qızıl qaydalar:** 1) datanı klasterdə bərabər yay; 2) partition sayını
minimallaşdır. Yazma ucuzdur — əsas məqsəd sorğuların səmərəliliyidir.
Modelləşdirmə **sorğu patterninə** görə edilir, əlaqəvi normallaşdırmaya görə YOX.

**Has-one — 3 variant:**
```sql
-- a) əlavə əlaqə cədvəli (laptop müstəqil varlıqdırsa):
CREATE TABLE user_laptop (
   user_id int PRIMARY KEY,   -- user bir dəfə
   sn int);                    -- amma laptop bir neçə user-ə bağlana bilər

-- b) laptop TYPE kimi (laptop həmişə user-ə bağlıdirsə):
CREATE TYPE laptop (
   sn int,
   model text,
   memory int);
CREATE TABLE users (
   user_id int PRIMARY KEY,
   name text,
   email text,
   laptop frozen<laptop>);   -- frozen → yalnız bütöv yenidən yazılır

-- c) laptop cədvəlində user_id (tək sahiblik zəmanəti):
CREATE TABLE laptops (
   sn int PRIMARY KEY,
   model text,
   memory int,
   user_id int);
```

**Has-many — collection ilə:**
```sql
CREATE TABLE users (
   user_id int PRIMARY KEY,
   name text,
   email text,
   laptops set<frozen<laptop>>);   -- set → təkrar laptop yoxdur
-- Amma eyni laptop 2 user-də ola bilər; zəmanət üçün laptops.user_id
```

**Many-to-many — sorğu əsasında İKİ cədvəl:**
```sql
CREATE TABLE user_laptops (
   user_id int,
   sn int,
   PRIMARY KEY (user_id, sn));      -- compound PK

CREATE TABLE laptops_user (
   sn int,
   user_id int,
   PRIMARY KEY (sn, user_id));
-- "user-in laptopları" → user_laptops (tək partition-hit!)
-- "laptop-un userləri" → laptops_user
```
- Compound PK (user_id, sn) → user-ə görə partition, daxilə sn sıralı

### 4. Data manipulyasiyası (Exec + Batch)
```go
const INSERT_QUERY = `INSERT INTO users
 (id,name,email) VALUES(?,?,?)`

err = session.Query(INSERT_QUERY, 1, "John",
    "john@gmail.com").WithContext(ctx).Exec()
```

**Batch (atomiklik + performans):**
```go
b := session.NewBatch(gocql.UnloggedBatch).WithContext(ctx)
b.Entries = append(b.Entries,
    gocql.BatchEntry{
        Stmt: INSERT_QUERY,
        Args: []interface{}{1, "John", "john@gmail.com"},
    },
    gocql.BatchEntry{
        Stmt: UPDATE_QUERY,   // UPDATE users SET email=? WHERE id=?
        Args: []interface{}{"otheremail@email.com", 1},
    })
err = session.ExecuteBatch(b)
```
- Batch ən effektiv **tək partition** hədəfində; çox-partition batch-i
  performansı aşağı salır (yalnız atomiklik lazımdırsa)

### 5. Sorğular
**Tək nəticə — Scan:**
```go
name := ""
err = session.Query("SELECT name FROM users WHERE id=1").
    WithContext(ctx).Scan(&name)      // dəyişənlər tip-uyğun olmalıdır
fmt.Println("Retrieved name", name)
```

**Çoxlu nəticə — Iter + Scanner (pagination daxili):**
```go
scanner := session.Query("SELECT * FROM users").
    WithContext(ctx).Iter().Scanner()
for scanner.Next() {
    var id int
    var name, email string
    err = scanner.Scan(&id, &name, &email)
    if err != nil {
        panic(err)
    }
    fmt.Println("Found:", id, name, email)
}
```
- **Qeyd:** sətirlər daxilərdə primary key-ə görə sıralanır —
  daxilət sırası ilə əlaqəsi yoxdur

### 6. UDT — struct + cql tag
**Nədir:** User Defined Type — Go struct-ı ilə birbaşa xəritələnir.

```go
type Laptop struct {
    Sn     int    `cql:"sn"`
    Model  string `cql:"model"`
    Memory int    `cql:"memory"`
}

// UDT yarat + cədvələ daxil et:
err = session.Query(LAPTOP_TYPE).WithContext(ctx).Exec()   // CREATE TYPE
err = session.Query(USERS_TABLE).WithContext(ctx).Exec()   // Laptop frozen<Laptop>
err = session.Query(INSERT, 1, "John", "john@gmail.com",
    &Laptop{100, "Lenovo", 10}).Exec()                     // struct birbaşa

// Oxu — Scan struct-a birbaşa dolur:
var retrieved Laptop
err = session.Query("select laptop from users where user_id=1").
    Scan(&retrieved)
fmt.Println("Retrieved", retrieved)   // {100 Lenovo 10}
```
- `cql:"sahə"` tag-i → UDT sahə adı ilə birləşdirir

## Əsas terminlər
- NoSQL — əlaqəvi olmayan data modeli
- Keyspace — Cassandra "database" anlayışı
- Partition — primary key-in paylanma hissəsi
- CQL — Cassandra Query Language
- GoCQL — rəsmi Go driver
- Session — sorğu bağlantıları vahidi
- Batch — atomik sorğu dəsti
- frozen — bütöv yazılan type modifier
- UDT (User Defined Type) — istifadəçi tipi (struct + cql tag)

## Praktik nəticə
Cassandra-ya SQL kimi yanaşmayın: model sorğuya görə qurulur (queries-first),
normallaşdırma deyil. İkiistiqamətli axtarış üçün iki cədvəl (denormalizasiya)
normaldır. Batch-lər tək partition üçün sürət, çox partition üçün yalnız
atomiklik verir. GoCQL iş axını: NewCluster → CreateSession → defer Close →
Query().Exec()/Scan()/Iter().Scanner(). UDT-lər cql tag-li struct-larla
birbaşa oxunub yazılır.

## Mənbə
Pages: 414-432 (PDF 414-432)
