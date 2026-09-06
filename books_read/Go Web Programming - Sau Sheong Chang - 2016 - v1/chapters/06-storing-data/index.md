# Chapter 6 — Storing data

## Bu chapter nədən bəhs edir?

Data saxlamanın 3 səviyyəsi: yaddaşda (struct + map), fayllarda (CSV, gob binary), verilənlər bazasında (database/sql + PostgreSQL — tam CRUD və one-to-many əlaqələr), həmçinin relational mapper-lər (Sqlx, Gorm).

## Əsas fikirlər

### 1. In-memory storage — map + pointer
**Nədir:** Proqram icrasında saxlanan data — struct-lar konteynerlərdə (array, slice, map; hətta stack/tree/queue).

**Kitabdan kod nümunəsi:**
```go
type Post struct {
    Id      int
    Content string
    Author  string
}

var PostById map[int]*Post
var PostsByAuthor map[string][]*Post

func store(post Post) {
    PostById[post.Id] = &post
    PostsByAuthor[post.Author] = append(PostsByAuthor[post.Author], &post)
}

func main() {
    PostById = make(map[int]*Post)
    PostsByAuthor = make(map[string][]*Post)

    post1 := Post{Id: 1, Content: "Hello World!", Author: "Sau Sheong"}
    post2 := Post{Id: 2, Content: "Bonjour Monde!", Author: "Pierre"}
    post3 := Post{Id: 3, Content: "Hola Mundo!", Author: "Pedro"}
    post4 := Post{Id: 4, Content: "Greetings Earthlings!", Author: "Sau Sheong"}

    store(post1)
    store(post2)
    store(post3)
    store(post4)

    fmt.Println(PostById[1])            // &{1 Hello World! Sau Sheong}
    for _, post := range PostsByAuthor["Sau Sheong"] {
        fmt.Println(post)
    }
}
```

**Sub-kod izahı:**
- 2 **indeks** = 2 map: `PostById` (ID → post pointer) + `PostsByAuthor` (müəllif → post pointer slice)
- `*Post` (pointer) — ID ilə də, author ilə də oxuyanda **eyni** post gəlsin deyə (kopya yox)
- İstifadə sahəsi: DB-dən gələn datanı cache-ləmək (Redis əvəzinə yaddaşda), tez cavab

### 2. Fayl əsaslı yazı/oxu — 2 səviyyə
**Sadə (ioutil):**
```go
data := []byte("Hello World!\n")
err := ioutil.WriteFile("data1", data, 0644)
read1, _ := ioutil.ReadFile("data1")
```

**File struct (daha çevik):**
```go
file1, _ := os.Create("data2")
defer file1.Close()
bytes, _ := file1.Write(data)

file2, _ := os.Open("data2")
defer file2.Close()
read2 := make([]byte, len(data))
bytes, _ = file2.Read(read2)
```
- `defer file.Close()` → funksiya qayıdanda çağrılır (LIFO stack)
- File struct: seek və digər mövqe əməliyyatları üçün əlavə metodlar

### 3. CSV (encoding/csv)
**Nədir:** Cədvəl datasının mətn formatında saxlanması; Excel/Numbers dəstəyi. İstifadəçi böyük data yükləməli olanda praktikidir.

**Kitabdan kod nümunəsi:**
```go
csvFile, err := os.Create("posts.csv")
if err != nil {
    panic(err)
}
defer csvFile.Close()

allPosts := []Post{
    Post{Id: 1, Content: "Hello World!", Author: "Sau Sheong"},
    // ...
}

// Yazma:
writer := csv.NewWriter(csvFile)
for _, post := range allPosts {
    line := []string{strconv.Itoa(post.Id), post.Content, post.Author}
    err := writer.Write(line)
    if err != nil {
        panic(err)
    }
}
writer.Flush()

// Oxuma:
file, err := os.Open("posts.csv")
defer file.Close()

reader := csv.NewReader(file)
reader.FieldsPerRecord = -1
record, err := reader.ReadAll()
if err != nil {
    panic(err)
}

var posts []Post
for _, item := range record {
    id, _ := strconv.ParseInt(item[0], 0, 0)
    post := Post{Id: int(id), Content: item[1], Author: item[2]}
    posts = append(posts, post)
}
```

**Sub-kod izahı:**
- `csv.NewWriter(f)` / `csv.NewReader(f)` → writer/reader fayl üzərində
- `strconv.Itoa(post.Id)` → int-i CSV sətri üçün string-ə çevir; oxuyanda `ParseInt` geri çevirir
- **`writer.Flush()` → buffer-i fayla yazmaq MÜTLƏQDIR** (yoxsa buffered data itə bilər)
- `reader.FieldsPerRecord = -1` → sahə sayı məcburi deyil; müsbət rəqəm = gözlənilən say (az olsa xəta); 0 = ilk sətirdən götür
- `ReadAll()` → bütün sətirlər ([]([]string)); böyük faylda sətir-sətir Read
- CSV nəticəsi: `1,Hello World!,Sau Sheong`

### 4. gob (encoding/gob) — Go binary format
**Nədir:** Go-ya xas binary serializasiya; encoder/decoder writer/reader üzərində. Sürətli backup, orderly shutdown, session/checkout-cart kimi müvəqqəti saxlama.

**Kitabdan kod nümunəsi:**
```go
func store(data interface{}, filename string) {
    buffer := new(bytes.Buffer)
    encoder := gob.NewEncoder(buffer)
    err := encoder.Encode(data)
    if err != nil {
        panic(err)
    }
    err = ioutil.WriteFile(filename, buffer.Bytes(), 0600)
    if err != nil {
        panic(err)
    }
}

func load(data interface{}, filename string) {
    raw, err := ioutil.ReadFile(filename)
    if err != nil {
        panic(err)
    }
    buffer := bytes.NewBuffer(raw)
    dec := gob.NewDecoder(buffer)
    err = dec.Decode(data)
    if err != nil {
        panic(err)
    }
}

func main() {
    post := Post{Id: 1, Content: "Hello World!", Author: "Sau Sheong"}
    store(post, "post1")
    var postRead Post
    load(&postRead, "post1")
    fmt.Println(postRead)   // {1 Hello World! Sau Sheong}
}
```

**Sub-kod izahı:**
- `bytes.Buffer` → həm Reader, həm Writer (Read/Write metodları var)
- `gob.NewEncoder(buffer).Encode(data)` → struct → binary → buffer → fayl
- `load` → fayl → buffer → `gob.NewDecoder(buffer).Decode(&postRead)` → struct
- `data interface{}` → empty interface: hər tip serializə oluna bilər

### 5. SQL — setup və qoşulma
**Setup addımları:**
```bash
createuser -P -d gwp      # -P: parol soruşsun; -d: DB yarada bilsin
createdb gwp              # user adı ilə eyni adlı DB
psql -U gwp -f setup.sql -d gwp
```

**setup.sql:**
```sql
create table posts (
  id     serial primary key,
  content text,
  author  varchar(255)
);
```

**Qoşulma:**
```go
import (
    "database/sql"
    "fmt"
    _ "github.com/lib/pq"
)

var Db *sql.DB

func init() {
    var err error
    Db, err = sql.Open("postgres", "user=gwp dbname=gwp password=gwp sslmode=disable")
    if err != nil {
        panic(err)
    }
}
```

**Sub-kod izahı:**
- `sql.Open("postgres", dsn)` → **connection pool** handle (`sql.DB`); DSN = data source name (driver-a xas string)
- **Lazily** qoşulur — Open yalnız strukturları hazırlayır, connection-i sonra yaradır
- `sql.DB` bağılmalı deyil — handle-dir, öz pool-u idarə edir
- `_ "github.com/lib/pq"` → blank import: driver paketinin `init`-i `sql.Register("postgres", &drv{})` edir; birbaşa istifadə ETMƏ — yalnız database/sql interfeysi istifadə olunur → driver dəyişsə kod dəyişmir
- Quraşdırma: `go get "github.com/lib/pq"`

### 6. CRUD — Create (prepared statement)
```go
func (post *Post) Create() (err error) {
    statement := "insert into posts (content, author) values ($1, $2) returning id"
    stmt, err := Db.Prepare(statement)
    if err != nil {
        return
    }
    defer stmt.Close()
    err = stmt.QueryRow(post.Content, post.Author).Scan(&post.Id)
    return
}
```

**Sub-kod izahı:**
- `Db.Prepare(statement)` → **prepared statement** — $1/$2 placeholder-li şablon, təkrar icra üçün; `defer stmt.Close()`
- `returning id` → PostgreSQL auto-increment id-ni qaytarır
- `stmt.QueryRow(...).Scan(&post.Id)` → tək sətir + Scan: dəyərləri parametrlərə kopyalayır → receiver-in Id sahəsi dolur
- Nəticə: `Post{0, ...}` → `Create()` → `Post{1, ...}` (struct artıq tam DB record-a uyğundur)

### 7. CRUD — Retrieve (QueryRow)
```go
func GetPost(id int) (post Post, err error) {
    post = Post{}
    err = Db.QueryRow("select id, content, author from posts where id = $1", id).Scan(&post.Id, &post.Content, &post.Author)
    return
}
```
- `Db.QueryRow` → tək `sql.Row` (çoxsətir nəticədə yalnız 1-ci götürülür); Scan xətaları qaytarır
- Boş struct yarat → chain QueryRow().Scan() → doldur

### 8. CRUD — Update / Delete (Exec)
```go
func (post *Post) Update() (err error) {
    _, err = Db.Exec("update posts set content = $2, author = $3 where id = $1", post.Id, post.Content, post.Author)
    return
}

func (post *Post) Delete() (err error) {
    _, err = Db.Exec("delete from posts where id = $1", post.Id)
    return
}
```
- `Db.Exec` → nəticə tələb etməyən əmrlər; `sql.Result` (təsirlənən sətir sayı + last id) — `_` ilə ignorə

### 9. Çoxsətirli oxu (Query + Rows iterator)
```go
func Posts(limit int) (posts []Post, err error) {
    rows, err := Db.Query("select id, content, author from posts limit $1", limit)
    if err != nil {
        return
    }
    for rows.Next() {
        post := Post{}
        err = rows.Scan(&post.Id, &post.Content, &post.Author)
        if err != nil {
            return
        }
        posts = append(posts, post)
    }
    rows.Close()
    return
}
```
- `Db.Query` → `sql.Rows` — iterator: `rows.Next()` → `io.EOF`-a qədər
- Hər iterasiyada struct yarat + Scan + append

### 10. One-to-many əlaqələr
**4 əlaqə növü:** one-to-one (has one), one-to-many (has many), many-to-one (belongs to), many-to-many.

**setup.sql (2 cədvəl):**
```sql
drop table posts cascade if exists;
drop table comments if exists;
create table posts (
  id     serial primary key,
  content text,
  author  varchar(255)
);
create table comments (
  id     serial primary key,
  content text,
  author  varchar(255),
  post_id integer references posts(id)
);
```
- `post_id ... references posts(id)` → **foreign key**; posts silinməsi üçün `cascade` lazımdır

**Struct-lar:**
```go
type Post struct {
    Id       int
    Content  string
    Author   string
    Comments []Comment      // has many
}

type Comment struct {
    Id      int
    Content string
    Author  string
    Post    *Post           // belongs to (pointer — eyni post!)
}
```
- Hər ikisi əslində pointerdir (slice = array-ə pointer) — kopya istəmirik, **eyni** Post-a işarə.

**Comment.Create:**
```go
func (comment *Comment) Create() (err error) {
    if comment.Post == nil {
        err = errors.New("Post not found")
        return
    }
    err = Db.QueryRow("insert into comments (content, author, post_id) values ($1, $2, $3) returning id",
        comment.Content, comment.Author, comment.Post.Id).Scan(&comment.Id)
    return
}
```
- Post YOXDURSA xəta qaytar (`errors.New("Post not found")`)
- `post_id` daxil etməklə əlaqə qurulur

**GetPost — əlaqəni oxu:**
```go
func GetPost(id int) (post Post, err error) {
    post = Post{}
    post.Comments = []Comment{}
    err = Db.QueryRow("select id, content, author from posts where id = $1", id).Scan(&post.Id, &post.Content, &post.Author)

    rows, err := Db.Query("select id, content, author from comments where post_id = $1", id)
    if err != nil {
        return
    }
    for rows.Next() {
        comment := Comment{Post: &post}
        err = rows.Scan(&comment.Id, &comment.Content, &comment.Author)
        if err != nil {
            return
        }
        post.Comments = append(post.Comments, comment)
    }
    rows.Close()
    return
}
```
- 2 sorğu: post + ona aid comment-lər; hər comment-ə `Post: &post` — geri əlaqə (many-to-one)

### 11. Sqlx — database/sql uzantısı
**Nədir:** database/sql interfeyslərini qoruyan əlavələr: row → struct/map/slice marshalling (struct tag ilə), prepared statement-lər üçün named parameter.

```go
import (
    "github.com/jmoiron/sqlx"
    _ "github.com/lib/pq"
)

type Post struct {
    Id         int
    Content    string
    AuthorName string `db:"author"`    // struct tag: kolon → sahə
}

var Db *sqlx.DB

func init() {
    Db, _ = sqlx.Open("postgres", "user=gwp dbname=gwp password=gwp sslmode=disable")
}

func GetPost(id int) (post Post, err error) {
    post = Post{}
    err = Db.QueryRowx("select id, content, author from posts where id = $1", id).StructScan(&post)
    return
}
```
- `sqlx.DB` + `QueryRowx` → `Rowx` → **`StructScan(&post)`** → kolonları sahələrə avtomatik map edir (tag ilə ad fərqli ola bilər)
- Quraşdırma: `go get "github.com/jmoiron/sqlx"`

### 12. Gorm — tam ORM
**Nədir:** Data-Mapper pattern-li ORM (kitabın əvvəlki nümunələri ActiveRecord pattern-i idi). Cədvəlləri avtomatik yaradır (migration), əlaqələri idarə edir, callback-lər, chain query-lər.

```go
import (
    "github.com/jinzhu/gorm"
    _ "github.com/lib/pq"
    "time"
)

type Post struct {
    Id        int
    Content   string
    Author    string `sql:"not null"`
    Comments  []Comment
    CreatedAt time.Time
}

type Comment struct {
    Id       int
    Content  string
    Author   string `sql:"not null"`
    PostId   int    `sql:"index"`
    CreatedAt time.Time
}

var Db gorm.DB

func init() {
    Db, _ = gorm.Open("postgres", "user=gwp dbname=gwp password=gwp sslmode=disable")
    Db.AutoMigrate(&Post{}, &Comment{})
}

func main() {
    post := Post{Content: "Hello World!", Author: "Sau Sheong"}
    Db.Create(&post)                                 // create (Data-Mapper!)

    comment := Comment{Content: "Good post!", Author: "Joe"}
    Db.Model(&post).Association("Comments").Append(comment)   // əlaqə

    var readPost Post
    Db.Where("author = $1", "Sau Sheong").First(&readPost)     // sorğu

    var comments []Comment
    Db.Model(&readPost).Related(&comments)                     // əlaqəli data
}
```

**Sub-kod izahı:**
- `Db.AutoMigrate(&Post{}, &Comment{})` → variadic — struct-lardan cədvəlləri yaradir/güncəlləyir (setup.sql LAZIM DEYİL!)
- `` `sql:"not null"` `` / `sql:"index"` → Gorm struct tag-ləri
- `CreatedAt time.Time` → avtomatik doldurulur (created_at kolonu)
- `PostId int` → Gorm avtomatik foreign key kimi tanıyır (Post field yoxdur!)
- `Db.Create(&post)` → **Db** yaradır (Data-Mapper), `post.Create()` deyil (ActiveRecord)
- `Db.Model(&post).Association("Comments").Append(comment)` → əlaqə yaradır, PostId-ni özü idarə edir
- `Db.Where(...).First(&readPost)` + `Db.Model(&readPost).Related(&comments)` → sorğu + əlaqəli data

Digər ORM-lər: Beego ORM, GORP.

## Data saxlama seçimləri cədvəli
| Üsul | Persistence | Sürət | İstifadə sahəsi |
|---|---|---|---|
| Yaddaş (map/struct) | ❌ (proqram bitəndə itir) | Ən sürətli | Cache |
| CSV faylı | ✅ | Orta | İstifadəçi ilə data mübadiləsi (Excel) |
| gob binary | ✅ | Sürətli | Backup, session, müvəqqəti saxlama |
| SQL DB | ✅ + server | DB-yə bağlı | Real web app |
| Sqlx / Gorm | ✅ | ORM overhead | Kod qısaltma |

## Əsas terminlər
- In-memory Storage (yaddaşda saxlama)
- Data Source Name (DSN)
- Connection Pool (`sql.DB`)
- Database Driver (verilənlər bazası drayveri)
- Blank Import (`_ "github.com/lib/pq"`)
- Prepared Statement (hazırlanmış sorğu — `$1` placeholder)
- QueryRow / Query / Exec
- Scan (sətir → parametr kopyası)
- Rows Iterator (sətir təkrarlayıcısı)
- CRUD (Create/Retrieve/Update/Delete)
- One-to-many / Many-to-one (bir-çox / çox-bir)
- Foreign Key (xarici açar)
- CSV (Comma-Separated Values)
- gob (Go binary format)
- bytes.Buffer (oxuyub-yazan bufer)
- Serialisation (seriyalaşdırma)
- ORM (Object-Relational Mapper)
- ActiveRecord vs Data-Mapper pattern
- Struct Tag (`db:"author"` / `sql:"not null"`)
- AutoMigrate (avtomatik miqrasiya)

## Praktik nəticə
- Cache üçün yaddaşda map + pointer saxla — 2 indeks = 2 map, hər biri pointer qaytarsın.
- CSV yazandan sonra `Flush`-u unutma; oxuyanda `FieldsPerRecord = -1` sərbəstlik verir.
- gob Go-dan kənar paylaşım üçün YOX — Go↔Go fayl mübadiləsi/backup üçün ideal.
- Driver-i yalnız blank import et (`_ "lib/pq"`) — kod database/sql üzərindən işləsin, driver dəyişə bilsin.
- Insert + auto-id: `returning id` + `QueryRow().Scan(&post.Id)` — struct DB record-a çevrilir.
- Nəticə lazım deyilsə `Exec`; tək sətir = `QueryRow`; çox sətir = `Query` + Next loop.
- Gorm `AutoMigrate` + `PostId` konvensiyası ilə setup.sql və manual əlaqə idarəsindən qurtulursan — amma arxada nə baş verdiyini bu chapter-dən bilirsən.

## Mənbə
Pages: 146-173 (PDF), book pages 125-152
