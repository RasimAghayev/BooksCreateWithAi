# Chapter 6 — Все о базах данных и хранилищах (Bazalar və saxlanc)

## Bu chapter nədən bəhs edir?

database/sql + MySQL, transaksiya interfeysi, bağlantı pulu + context
timeout, Redis, MongoDB (BSON), saxlanc interfeysləri (portability və mock
üçün).

## Əsas fikirlər

### 1. database/sql + MySQL
**Nədir:** Standart SQL paketi — driver-lər daxil edilərək plaqin modeli
işləyir; bağlantı pulu avtomatik idarə olunur.

**Kitabdan kod nümunəsi:**
```go
import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"   // yalnız side-effect üçün import!
)

func Setup() (*sql.DB, error) {
    db, err := sql.Open("mysql",
        fmt.Sprintf("%s:%s@/gocookbook?parseTime=true",
            os.Getenv("MYSQLUSERNAME"), os.Getenv("MYSQLPASSWORD")))
    return db, err
}

// DDL + insert — Exec (nəticə satırı YOX):
db.Exec("CREATE TABLE example (name VARCHAR(20), created DATETIME)")
db.Exec(`INSERT INTO example (name, created) values ("Aaron", NOW())`)

// Sorğu — Query + rows.Scan:
rows, err := db.Query("SELECT name, created FROM example where name=?", name)
defer rows.Close()                       // rows connection-u azad edir
for rows.Next() {
    var e Example
    if err := rows.Scan(&e.Name, &e.Created); err != nil {  // sütun→sahə
        return err
    }
}
return rows.Err()                        // tsikldən sonrakı xəta
```

**Sub-kod izahı:**
- `_ "driver"` → blank import: driver özünü database/sql-ə Qeyd edir
- `sql.Open` bağlantı AÇMIR — sadəcə validation; bağlantı ilk sorğuda alınır
- `?parseTime=true` → DATETIME Go time.Time-ya çevrilir
- `?` placeholder → SQL injection qoruması
- `rows.Close()` vacibdir — pull-dan bağlantı qaytarmır

### 2. Transaksiya interfeysi (DB və Tx üçün ümumi)
**Nədir:** Funksiyaların həm sql.DB, həm sql.Tx ilə işləməsi üçün DB
interfeysi — mock üçün də istifadə olunur.

**Kitabdan kod nümunəsi:**
```go
// DB — həm sql.DB, həm sql.Tx təmin edir:
type DB interface {
    Exec(query string, args ...interface{}) (sql.Result, error)
    Prepare(query string) (*sql.Stmt, error)
    Query(query string, args ...interface{}) (*sql.Rows, error)
    QueryRow(query string, args ...interface{}) *sql.Row
}
type Transaction interface {
    DB                                   // DB-ni EMBED edir
    Commit() error
    Rollback() error
}

// Funksiya interfeys qəbul edir — DB və ya Tx fərqi yoxdur:
func Create(db DB) error { db.Exec(...) }

// İstifadə — tx-i ötür:
tx, err := db.Begin()
defer tx.Rollback()          // commit uğurlu olsa rollback no-op-dur
if err := dbinterface.Exec(tx); err != nil {
    panic(err)
}
tx.Commit()
```

**Sub-kod izahı:**
- `defer tx.Rollback()` → xəta halında avtomatik geri alma; commit-dən
  sonra çağırışı heç nə etmir
- Eyni Create/Query funksiyaları tək bağlantıda və transaksiya daxilində
  İSTİFADƏ OLUNUR — çox cədvəl atomik update üçün
- Interfeys mock-lanması asan olur (Chapter 9)

### 3. Bağlantı pulu, rate limit və timeout
**Nədir:** MaxOpen/MaxIdle Conns + context ilə timeout — mikroservis
miqyasında DB-ni qorumaq.

**Kitabdan kod nümunəsi:**
```go
func Setup() (*sql.DB, error) {
    db, err := sql.Open("mysql", dsn)
    if err != nil { return nil, err }
    db.SetMaxOpenConns(24)     // maksimum AÇIQ bağlantı
    db.SetMaxIdleConns(24)     // idle pul dərinliyi (MaxOpen-dan çox ola bilməz)
    return db, nil
}

// Context ilə deadline:
func ExecWithTimeout() error {
    db, _ := Setup()
    ctx := context.Background()
    ctx, cancel := context.WithDeadline(ctx, time.Now())  // dərhal timeout!
    defer cancel()                                          // resurs sızıntısının qarşısı
    _, err := db.BeginTx(ctx, nil)     // context-aware transaksiya
    return err                         // "context deadline exceeded"
}
```

**Sub-kod izahı:**
- `SetMaxOpenConns` → paralel bağlantı limiti; mikroshservis horizontal
  scale olunanda DB-nin boğulmaması üçün
- `context.WithDeadline` → sorğu (bağlantı gözləmə də daxil) vaxtı bitir
- `defer cancel()` MÜTLƏQ — context resursları təmizlənsin
- Mülayim: məhdud pul + sərt timeout → yüklü sistemdə tez-tez timeout xətaları

### 4. Redis
**Nədir:** Key-value saxlanc — sessiya, müvəqqəti data, TTL dəstəyi.

**Kitabdan kod nümunəsi:**
```go
func Setup() (*redis.Client, error) {
    client := redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: os.Getenv("REDISPASSWORD"),
        DB:       0,
    })
    _, err := client.Ping().Result()   // bağlantı yoxlaması
    return client, err
}

// Set/Get + TTL:
conn.Set("key", "value", 5*time.Second)    // 5 saniyəlik ömür
var result string
if err := conn.Get("key").Scan(&result); err != nil {
    switch err {
    case redis.Nil:          // açar TAPILMADI — xəta deyil!
        return nil
    default:
        return err
    }
}

// List + sort:
conn.LPush(listkey, 1); conn.LPush(listkey, 3); conn.LPush(listkey, 2)
defer conn.Del(listkey)     // təmizləmə
res, err := conn.Sort(listkey, redis.Sort{Order: "ASC"}).Result()  // [1 2 3]
```

**Sub-kod izahı:**
- `redis.Nil` → açar yoxdur xətası — adi err-dən AYRI emal edilməlidir
- `Set(key, value, ttl)` → avtomatik expires
- LPush → listin başına; Sort → server tərəfində sıralama

### 5. MongoDB (BSON)
**Nədir:** Document NoSQL — struct-lar demək olar sxemsiz saxlanılır;
BSON marshal JSON-a bənzər.

**Kitabdan kod nümunəsi:**
```go
func Setup(ctx context.Context, address string) (*mongo.Client, error) {
    ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()
    client, err := mongo.NewClient(options.Client().ApplyURI(address))
    if err != nil { return nil, err }
    if err := client.Connect(ctx); err != nil {   // bağlantı
        return nil, err
    }
    return client, nil
}

type State struct {
    Name       string `bson:"name"`
    Population int    `bson:"pop"`
}

coll := db.Database("gocookbook").Collection("example")
// Çoxlu insert bir çağırışla:
vals := []interface{}{&State{"Washington", 7062000}, &State{"Oregon", 3970000}}
coll.InsertMany(ctx, vals)

// Sorğu + decode:
var s State
coll.FindOne(ctx, bson.M{"name": "Washington"}).Decode(&s)
coll.Drop(ctx)     // təmizləmə
```

**Sub-kod izahı:**
- `bson.M` → map şəklində sorğu filtri (JSON-dakı obyektə bənzər)
- `bson:"..."` struct tag-ləri sahə adlarını xəritələşdirir
- `FindOne(...).Decode(&s)` → sənədi struct-a çevirir
- Context əsasparametrdir — timeout/cancel imkanı

### 6. Saxlanc interfeysləri (portability)
**Nədir:** Storage arxitekturasını interface arxasında gizlətmək — backend
dəyişə bilər, mock asan olur.

**Kitabdan kod nümunəsi:**
```go
type Item struct {
    Name  string
    Price int64
}

// İnterfeys — konkrut backend-dən asılı deyil:
type Storage interface {
    GetByName(context.Context, string) (*Item, error)
    Put(context.Context, *Item) error
}

// MongoDB realizasiyası:
type MongoStorage struct {
    *mongo.Client
    DB         string
    Collection string
}
func (m *MongoStorage) GetByName(ctx context.Context, name string) (*Item, error) {
    c := m.Client.Database(m.DB).Collection(m.Collection)
    var i Item
    if err := c.FindOne(ctx, bson.M{"name": name}).Decode(&i); err != nil {
        return nil, err
    }
    return &i, nil
}
func (m *MongoStorage) Put(ctx context.Context, i *Item) error {
    c := m.Client.Database(m.DB).Collection(m.Collection)
    _, err := c.InsertOne(ctx, i)
    return err
}

// ƏSAS: PerformOperations YALNIZ interfeys bilir:
func PerformOperations(s Storage) error {
    ctx := context.Background()
    i := Item{Name: "candles", Price: 100}
    s.Put(ctx, &i)                       // mongo? fayl? API? — fərqi yoxdur
    candles, err := s.GetByName(ctx, "candles")
    fmt.Printf("Result: %#v\n", candles)
    return nil
}
```

**Sub-kod izahı:**
- `PerformOperations(s Storage)` → backend-i bir sətirlə əvəz etmək olar
- Context parametrləri → timeout idarəsi interfeys səviyyəsində
- Çətinlik: bir neçə əməliyyatın EYNİ transaksiyada olması — bəzən
  komposit əməliyyatlar və ya əlavə arqumentlər tələb edir

## Əsas terminlər

- database/sql / Driver (blank import)
- Connection Pool (bağlantı pulu)
- MaxOpenConns / MaxIdleConns
- Query / Exec / Scan
- Transaction (transaksiya) / Commit / Rollback
- Context-aware əməliyyatlar (BeginTx)
- Redis / TTL (yaşam müddəti) / redis.Nil
- MongoDB / BSON / Document (sənəd)
- Storage Interface (saxlanc interfeysi)
- Mock (yamsılama) / Portability (daşına bilənlik)

## Praktik nəticə

- Driver-lər blank import ilə qoşulur; `?` placeholder həmişə parametrləşdir
- rows.Close() və defer tx.Rollback() — hər zaman; rows.Err() yoxla
- Pul limitləri: mikroservis nüsxə sayı × MaxOpenConns ≤ DB max_connections
- DB/Tx üçün ümumi interfeys → funksiyalar hər iki kontekstdə işləyir
- Storage interfeysi: test mock + backend dəyişməsi; sərhətləri düzgün çək
- redis.Nil fərqli emal olunur; TTL sessiya data üçün ideal

## Mənbə

Pages: 191-217 (Chapter 6, Go Programming Cookbook 2nd ed)
