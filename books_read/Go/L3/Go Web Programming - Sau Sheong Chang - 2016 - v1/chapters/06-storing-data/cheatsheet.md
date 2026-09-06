# Chapter 6 — Storing data Cheatsheet

## `ioutil.WriteFile("file", data, 0644)` / `ioutil.ReadFile("file")`

**Nə edir:** Asıl fayl yazı Rox / oxur. Bayt massiv ilə işləyir.

```go
data := []byte("Hello World!\n")
err := ioutil.WriteFile("data1", data, 0644)  // 0644 = owner rw, group/others r
read1, _ := ioutil.ReadFile("data1")
fmt.Print(string(read1))
```

**Parametrlər:** `perm` — fayl icazəsi.

**Mənbə:** Chapter 6, page 129

---

## `os.Create` + `file.Write` + `defer file.Close()`

**Nə edir:** `File` struct-ı ilə fayl yaradır, `Write`/`Read` metodları ilə işləyir. Daha çox üçün `ioutil`-ə görə.

```go
file1, _ := os.Create("data2")
defer file1.Close()           // defer — funksiya bitdikdən sonra bağla
bytes, _ := file1.Write(data)
file2, _ := os.Open("data2")
defer file2.Close()
read2 := make([]byte, len(data))
file2.Read(read2)
```

**`defer`** — çağrı `defer`-dən sonra stack-ə qoyulur; funksiya qaytdıqda FIFO ilə icra olunur.

**Mənbə:** Chapter 6, page 130

---

## CSV: `csv.NewWriter(csvFile)` + `writer.Write(line)` + `writer.Flush()`

**Nə edir:** CSV fayla yazır və oxur.

```go
writer := csv.NewWriter(csvFile)
line := []string{strconv.Itoa(post.Id), post.Content, post.Author}
writer.Write(line)
writer.Flush()  // bufferı boşalt — əlaqəsiz Flush edilə bsa data itir

reader := csv.NewReader(file)
reader.FieldsPerRecord = -1  // -1: dəqiqlik yox; 0: ilk setirdən; müsbət: gözlənilən
record, _ := reader.ReadAll()
for _, item := range record {
    id, _ := strconv.ParseInt(item[0], 0, 0)
    post := Post{Id: int(id), Content: item[1], Author: item[2]}
}
```

**Mənbə:** Chapter 6, page 131

---

## Gob: `gob.NewEncoder(buffer)` + `encoder.Encode(data)` / `gob.NewDecoder` + `dec.Decode`

**Nə edir:** Go binary serializasiyası — `gob` formatı. `bytes.Buffer` Reader+Writer iki üsüldə.

```go
func store(data interface{}, filename string) {
    buffer := new(bytes.Buffer)
    encoder := gob.NewEncoder(buffer)
    encoder.Encode(data)
    ioutil.WriteFile(filename, buffer.Bytes(), 0600)
}
func load(data interface{}, filename string) {
    raw, _ := ioutil.ReadFile(filename)
    buffer := bytes.NewBuffer(raw)
    dec := gob.NewDecoder(buffer)
    dec.Decode(data)
}
```

**`data interface{}``** — boş interfeys, istənilən tip qəbul edir.

**Mənbə:** Chapter 6, page 133

---

## `sql.Open("postgres", "user=... dbname=... sslmode=disable")` + `var Db *sql.DB`

**Nə edir:** `sql.DB` (connection pool) yaradır. `init()`-də çağrılır.

```go
var Db *sql.DB
func init() {
    var err error
    Db, err = sql.Open("postgres", "user=gwp dbname=gwp password=gwp sslmode=disable")
    if err != nil { panic(err) }
}
```

**Diqqət:** `Open` connection yox, pool konfiqurasiyasını yaradır (lazy). Real bağlantı `Query`-də.

**Mənbə:** Chapter 6, page 138

---

## `Db.Prepare(stmt)` + `stmt.QueryRow().Scan(&post.Id)` (Create)

**Nə edir:** SQL prepared statement + single-row query + Scan-ın id-ni əldə et.

```go
func (post *Post) Create() (err error) {
    statement := "insert into posts (content, author) values ($1, $2) returning id"
    stmt, err := Db.Prepare(statement)
    defer stmt.Close()
    err = stmt.QueryRow(post.Content, post.Author).Scan(&post.Id)
    return
}
```

**`$1, $2`** — PostgreSQL placeholder (MySQL-də `?`).

**Mənbə:** Chapter 6, page 139

---

## `Db.QueryRow("...$1", id).Scan(...)` (GetPost) / `Db.Query("...limit $1", limit)` + `rows.Next()`

**Nə edir:** Single-row və multi-row sorğular.

```go
// Single
func GetPost(id int) (post Post, err error) {
    err = Db.QueryRow("select id, content, author from posts where id = $1", id).
        Scan(&post.Id, &post.Content, &post.Author)
    return
}

// Multiple
func Posts(limit int) (posts []Post, err error) {
    rows, err := Db.Query("select id, content, author from posts limit $1", limit)
    for rows.Next() {
        post := Post{}
        err = rows.Scan(&post.Id, &post.Content, &post.Author)
        posts = append(posts, post)
    }
    rows.Close()
    return
}
```

**`rows`** — iterator (`io.EOF`-a qədər `Next()`); `rows.Close()` vacib.

**Mənbə:** Chapter 6, page 141

---

## `Db.Exec("update ...")` / `Db.Exec("delete ...")`

**Nə edir:** INSERT/UPDATE/DELETE üçün `Exec` — `sql.Result` + `error`. `sql.Result` — affected rows + last insert id.

```go
func (post *Post) Update() (err error) {
    _, err = Db.Exec("update posts set content = $2, author = $3 where id = $1",
        post.Id, post.Content, post.Author)
    return
}
func (post *Post) Delete() (err error) {
    _, err = Db.Exec("delete from posts where id = $1", post.Id)
    return
}
```

**`_` (underscore)** — `sql.Result`-ı yox edirik, istifadə etmirik.

**Mənbə:** Chapter 6, page 142

---

## Foreign key: `post_id integer references posts(id)`

**Nə edir:** SQL foreign key — `comments.post_id`-ı `posts.id`-yə istinad edir.

```sql
create table comments (
    id serial primary key,
    content text,
    author varchar(255),
    post_id integer references posts(id)
);
drop table posts cascade if exists;  -- FK olan table-lər üçün cascade
```

**`cascade`** — parent table silinəndə uğurlu FK-ları da silir.

**Mənbə:** Chapter 6, page 144

---

## Sqlx: `Db.QueryRowx(...).StructScan(&post)` + `db: author` tag

**Nə edir:** `database/sql`-a extension; `StructScan` ilə row-u struct-a avtomatik map edir, `db:` tag-ı ilə sütun adını override edir.

```go
type Post struct {
    AuthorName string `db: author`
}
var Db *sqlx.DB
err = Db.QueryRowx("select id, content, author from posts where id = $1", id).StructScan(&post)
```

**Mənbə:** Chapter 6, page 148

---

## Gorm: `Db.AutoMigrate(&Post{}, &Comment{})` + `Db.Create(&post)` + `Db.Where("...").First(&readPost)`

**Nə edir:** Data Mapper ORM. `AutoMigrate` — struct-dan table auto yaradır.

```go
type Comment struct {
    Author string `sql:"not null"`
    PostId int `sql:"index"`
    CreatedAt time.Time
}
var Db gorm.DB
Db, _ = gorm.Open("postgres", "user=gwp ...")
Db.AutoMigrate(&Post{}, &Comment{})
Db.Create(&post)
Db.Model(&post).Association("Comments").Append(comment)
Db.Where("author = $1", "Sau Sheong").First(&readPost)
Db.Model(&readPost).Related(&comments)
```

**`CreatedAt`** — avtomatik populat. `Model/Where/First/Related` — chain query.

**Mənbə:** Chapter 6, page 150
