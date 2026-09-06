# Fəsil 7 — Go web services

## Bu fəsil nədən bəhs edir?

Web servislərin nəzəriyyəsi (SOAP vs REST), XML-in parse/create (Unmarshal/Decoder, Marshal/Encoder, struct tag qaydaları), JSON-un parse/create və tam RESTful web servisinin qurulması (Ch 6-dakı CRUD-un HTTP üzərində wrap edilməsi).

## Əsas fikirlər

### 1. Web service nədir?
**Tərif:** Son istifadəçisi insan yox, **proqram** olan, HTTP üzərindən ünsiyyət quran proqram. W3C tərifinə görə "machine-to-machine interaction over a network".

**Növlər:** SOAP-based (enterprise), REST-based (public API-lərin 54%-i), XML-RPC-based. 

**SOAP vs REST müqayisəsi:**
| Xüsusiyyət | SOAP | REST |
|---|---|---|
| Nədir? | Protokol (XML, envelope) | Dizayn fəlsəfəsi |
| Data format | Məcburi XML | Sərbəst (ən çox JSON) |
| Tərif dili | WSDL (məcburi kontrakt) | WADL (az istifadə, standart deyil; Swagger/RAML alternativ) |
| Üslub | Function-driven (RPC tipli) | Data-driven (resurs + verb) |
| Güc | Standartlaşmış, WS-* uzantıları (WS-Security, WS-Addressing), UDDI discovery | Sürətli, sadə, çevik |
| Zəif tərəf | Verbose, troubleshooting çətin, WSDL dəyişikliyi client regenerasiya tələb edir (version lock-in) | Formal kontrakt yoxdur |

**Praktik strategiya:** enterprise daxili inteqrasiya üçün SOAP + xarici developer-lər üçün REST — hər ikisinin gücü öz yerində.

### 2. SOAP — quruluş
**Envelope modeli:** SOAP mesajı = envelope (göndərmə konteyneri) + body (request/response). `Content-Type: application/soap+xml`. Adətən HTTP POST (SOAP 1.2 GET də icazə verir).

**Request nümunəsi:**
```xml
<?xml version="1.0"?>
<soap:Envelope xmlns:soap="http://www.w3.org/2001/12/soap-envelope"
  soap:encodingStyle="http://www.w3.org/2001/12/soap-encoding">
  <soap:Body xmlns:m="http://www.chitchat.com/forum">
    <m:GetCommentRequest>
      <m:CommentId>123</m:CommentId>
    </m:GetCommentRequest>
  </soap:Body>
</soap:Envelope>
```

**WSDL:** service → port → binding → operation → input/output message → part (type). Client-lar WSDL-dən **generasiya olunur** — dəyişiklik = client regenerasiya. "GetCommentService" nümunəsində: message GetCommentRequest (CommentId: string) → portType GetCommentPortType → operation GetComment → binding (rpc style, HTTP transport) → service (location: localhost:8080/GetComment).

### 3. REST — fəlsəfə
**Əsas:** resurs-lar (URL) + verb-lər (HTTP metodları).

**CRUD ↔ HTTP metodları (1-dən-1-ə mapping DEYİL):**
| HTTP | İstifadə | Nümunə |
|---|---|---|
| POST | Resurs yaradır (yeni URL!) | POST /users |
| GET | Resurs oxuyur | GET /users/1 |
| PUT | Müəyyən URL-dəki resursu yeniləyir/yaradır | PUT /users/1 |
| DELETE | Resursu silir | DELETE /users/1 |
| PATCH | Qismən yeniləmə | PATCH /users/1 |

**POST vs PUT:** PUT — hansı URL-i əvəz edəcəyini bilirsən; POST — yeni URL yaradır. PUT **idempotent**dir (təkrar çağırış dəyişiklik yaratmır), POST yox.

**"Aha!" anları:**
1. HTTP metodları ↔ CRUD mapping-i
2. 4 metoddan çoxu var: PATCH — qismən update

**REST request vs SOAP request:**
```
GET /comment/123 HTTP/1.1          ← REST: body YOXDUR, URL-də hər şey
POST /GetComment HTTP/1.1          ← SOAP: body-da envelope
```

### 4. Restriktiv olmayan həllər (action → resurs)
REST arbitrary action-lara icazə vermir (`ACTIVATE /user/456` YOX). 2 həll:

**a) Action-u resursa çevir (reify):**
```
POST /user/456/activation HTTP/1.1
{ "date": "2015-05-15T13:05:05Z" }
```
→ "activation" özəlkisi olan resurs; əlavə parametr qazanılır.

**b) Action-u resursun xüsusiyyətinə çevir:**
```
PATCH /user/456 HTTP/1.1
{ "active" : "true" }
```
→ Qismən update (PATCH) ilə status dəyişmək.

### 5. XML parse — Unmarshal
**Alqoritm:** struct yarat → `xml.Unmarshal(data, &struct)`.

**post.xml:**
```xml
<?xml version="1.0" encoding="utf-8"?>
<post id="1">
  <content>Hello World!</content>
  <author id="2">Sau Sheong</author>
</post>
```

**Kitabdan kod nümunəsi:**
```go
type Post struct {
    XMLName xml.Name `xml:"post"`
    Id      string   `xml:"id,attr"`
    Content string   `xml:"content"`
    Author  Author   `xml:"author"`
    Xml     string   `xml:",innerxml"`
}

type Author struct {
    Id   string `xml:"id,attr"`
    Name string `xml:",chardata"`
}

func main() {
    xmlFile, err := os.Open("post.xml")
    defer xmlFile.Close()
    xmlData, err := ioutil.ReadAll(xmlFile)

    var post Post
    xml.Unmarshal(xmlData, &post)
    fmt.Println(post)
    // {{ post} 1 Hello World! {2 Sau Sheong} <content>...</content>...}
}
```

**Struct tag qaydaları (6 qayda):**
1. Element adı üçün: `XMLName xml.Name` sahəsi + `` `xml:"post"` ``
2. Atribut: `` `xml:"<ad>,attr"` ``
3. Character data: `` `xml:",chardata"` ``
4. Raw daxili XML: `` `xml:",innerxml"` ``
5. Mode flag yoxdursa → eyni adlı child element-ə map olunur
6. **Leap-frog (aşma):** `` `xml:"a>b>c"` `` → ara element-ləri keçib c-yə birbaşa çıxış (məs. `comments>comment` — Comments strukturuna ehtiyac yoxdur)

**Tag qaydaları:** key=xml, value cüt dırnaqda, backtick içində (backtick = escape-siz string). Struct və sahələr **exported** (böyük hərf) olmalıdır.

### 6. XML parse — Decoder (streaming)
**Nəyə lazım:** Böyük/streaming XML faylları — Unmarshal bütün məlumatı yaddaşa yükləyir; Decoder element-element oxuyur.

**Kitabdan kod nümunəsi:**
```go
decoder := xml.NewDecoder(xmlFile)
for {
    t, err := decoder.Token()
    if err == io.EOF {
        break
    }
    if err != nil {
        fmt.Println("Error decoding XML into tokens:", err)
        return
    }
    switch se := t.(type) {
    case xml.StartElement:
        if se.Name.Local == "comment" {
            var comment Comment
            decoder.DecodeElement(&comment, &se)
        }
    }
}
```

**Sub-kod izahı:**
- `xml.NewDecoder(io.Reader)` → decoder
- `decoder.Token()` → növbəti XML token (interface — XML element təmsilçisi); `io.EOF` → bitdi
- `switch se := t.(type)` → tokenin tipini yoxla: `xml.StartElement` = açılış tag-i
- `se.Name.Local == "comment"` → element adı yoxla
- `decoder.DecodeElement(&comment, &se)` → həmin elementi struct-a decode et

### 7. XML create — Marshal / Encoder
**Marshal:**
```go
post := Post{Id: "1", Content: "Hello World!", Author: Author{Id: "2", Name: "Sau Sheong"}}

output, err := xml.Marshal(&post)                     // bir sətir XML
output, err := xml.MarshalIndent(&post, "", "\t")     // girintili
err = ioutil.WriteFile("post.xml", []byte(xml.Header+string(output)), 0644)
```
- `xml.MarshalIndent(&post, prefix, indent)` → gözəl format
- **`xml.Header`** → XML declaration (`<?xml version="1.0" encoding="UTF-8"?>`) avtomatik yazılmır — manual əlavə et!

**Encoder (streaming/fayl):**
```go
xmlFile, _ := os.Create("post.xml")
encoder := xml.NewEncoder(xmlFile)
encoder.Indent("", "\t")
err = encoder.Encode(&post)
```

### 8. JSON parse — Unmarshal / Decoder
**Tag qaydası (yalnız 1):** `` `json:"<key>"` `` — sahəni JSON açarına map et. Qalan hər şey XML ilə eyni.

**Kitabdan kod nümunəsi:**
```go
type Post struct {
    Id       int       `json:"id"`
    Content  string    `json:"content"`
    Author   Author    `json:"author"`
    Comments []Comment `json:"comments"`
}

var post Post
json.Unmarshal(jsonData, &post)
```

**Decoder (streaming):**
```go
decoder := json.NewDecoder(jsonFile)
for {
    var post Post
    err := decoder.Decode(&post)
    if err == io.EOF {
        break
    }
    if err != nil {
        fmt.Println("Error decoding JSON:", err)
        return
    }
    fmt.Println(post)
}
```

**Decoder vs Unmarshal seçimi:**
- Data `io.Reader` stream-dən gəlir (məs. `http.Request.Body`) → **Decoder**
- Data string/yaddaşda → **Unmarshal**

### 9. JSON create — Marshal / Encoder
```go
output, err := json.MarshalIndent(&post, "", "\t\t")
ioutil.WriteFile("post.json", output, 0644)

// Encoder:
jsonFile, _ := os.Create("post.json")
encoder := json.NewEncoder(jsonFile)
err = encoder.Encode(&post)
```

### 10. Tam RESTful web service (server.go + data.go)
**data.go — DB qatı (Ch 6-dan sadələşdirilmiş):**
```go
package main

import (
    "database/sql"
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

func retrieve(id int) (post Post, err error) {
    post = Post{}
    err = Db.QueryRow("select id, content, author from posts where id = $1", id).Scan(&post.Id, &post.Content, &post.Author)
    return
}

func (post *Post) create() (err error) {
    statement := "insert into posts (content, author) values ($1, $2) returning id"
    stmt, err := Db.Prepare(statement)
    if err != nil {
        return
    }
    defer stmt.Close()
    err = stmt.QueryRow(post.Content, post.Author).Scan(&post.Id)
    return
}

func (post *Post) update() (err error) {
    _, err = Db.Exec("update posts set content = $2, author = $3 where id = $1", post.Id, post.Content, post.Author)
    return
}

func (post *Post) delete() (err error) {
    _, err = Db.Exec("delete from posts where id = $1", post.Id)
    return
}
```

**server.go — HTTP routing (method-based dispatch):**
```go
package main

import (
    "encoding/json"
    "net/http"
    "path"
    "strconv"
)

type Post struct {
    Id      int    `json:"id"`
    Content string `json:"content"`
    Author  string `json:"author"`
}

func main() {
    server := http.Server{
        Addr: "127.0.0.1:8080",
    }
    http.HandleFunc("/post/", handleRequest)
    server.ListenAndServe()
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
    var err error
    switch r.Method {
    case "GET":
        err = handleGet(w, r)
    case "POST":
        err = handlePost(w, r)
    case "PUT":
        err = handlePut(w, r)
    case "DELETE":
        err = handleDelete(w, r)
    }
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}
```

**Sub-kod izahı:**
- `switch r.Method` → REST-in method dispatch-i — hamısı eyni URL-də (`/post/`), fərqlər HTTP metodunda
- `http.Error(w, err.Error(), 500)` → hər handler-dən qalxan xəta → 500 + error mətni

**handleGet:**
```go
func handleGet(w http.ResponseWriter, r *http.Request) (err error) {
    id, err := strconv.Atoi(path.Base(r.URL.Path))
    if err != nil {
        return
    }
    post, err := retrieve(id)
    if err != nil {
        return
    }
    output, err := json.MarshalIndent(&post, "", "\t\t")
    if err != nil {
        return
    }
    w.Header().Set("Content-Type", "application/json")
    w.Write(output)
    return
}
```
- `path.Base(r.URL.Path)` → URL-in son seqmentini al (`/post/1` → "1"); `strconv.Atoi` → int
- Struct → `json.MarshalIndent` → `Content-Type: application/json` → Write

**handlePost:**
```go
func handlePost(w http.ResponseWriter, r *http.Request) (err error) {
    len := r.ContentLength
    body := make([]byte, len)
    r.Body.Read(body)
    var post Post
    json.Unmarshal(body, &post)
    err = post.create()
    if err != nil {
        return
    }
    w.WriteHeader(200)
    return
}
```
- Body → byte slice → `json.Unmarshal` → struct → `create()` → DB-yə yaz → 200

**handlePut:** id al → `retrieve(id)` → body-ni **mövcud struct-a** Unmarshal → `update()`
**handleDelete:** id al → `retrieve(id)` → `delete()` → 200

**cURL test əmrləri:**
```bash
# Create:
curl -i -X POST -H "Content-Type: application/json" \
  -d '{"content":"My first post","author":"Sau Sheong"}' http://127.0.0.1:8080/post/

# Read:
curl -i -X GET http://127.0.0.1:8080/post/1

# Update:
curl -i -X PUT -H "Content-Type: application/json" \
  -d '{"content":"Updated post","author":"Sau Sheong"}' http://127.0.0.1:8080/post/1

# Delete:
curl -i -X DELETE http://127.0.0.1:8080/post/1
```

## Əsas terminlər
- Web Service (web servisi)
- SOAP / Envelope (zarf)
- WSDL (Web Service Definition Language)
- WS-* (WS-Security, WS-Addressing)
- UDDI (kataloq servisi)
- REST (Representational State Transfer)
- Resource / Verb (resurs / fel)
- Idempotency (eynilik — PUT vs POST)
- Reify (abstrakt əməliyyatı resursa çevirmək)
- PATCH (qismən yeniləmə metodu)
- WADL / Swagger / RAML
- struct tag (`xml:"..."` / `json:"..."`)
- `,attr` / `,chardata` / `,innerxml` (XML mode flag-ləri)
- Leap-frog (`a>b>c` — aşma sintaksisi)
- xml.Name / XMLName
- Unmarshal / Marshal / MarshalIndent
- Decoder / Encoder (streaming)
- xml.Header (declaration sabiti)
- io.EOF (axın sonu)
- Method-based Dispatch (metod əsaslı yönəltmə)
- path.Base (yolun son seqmenti)

## Praktik nəticə
- Yeni API qurursansa REST seç — sadəlik + JSON sürəti; SOAP yalnız legacy enterprise inteqrasiyada qalıb.
- Action-ları REST-də modelləşdirməyin 2 yolu: resursa çevir (POST /user/456/activation) yaxud PATCH ilə status dəyiş.
- XML struct tag mode flag-lərini yadda saxla: `attr`, `chardata`, `innerxml`; nested-də `a>b>c` aşma sintaksisi Comments strukturundan xilas edir.
- Stream gəliri varsa (request body, böyük fayl) **Decoder/Encoder**; string varsa **Unmarshal/Marshal**.
- `xml.Header`-i manual əlavə et — Go declaration-ı avtomatik yazmır.
- Method dispatch-i `switch r.Method` ilə tək handler-də həll et — bu sadə REST framework-ün əsasıdır.
- Xətaları mərkəzdən emal et: handler-lar `(err error)` qaytarsın, `handleRequest` 500 qaytarsın.

## Mənbə
Pages: 175-210 (PDF), book pages 155-189
