# Chapter 11 — Working with REST APIs (REST API-lərlə İş)

## Bu chapter nədən bəhs edir?

REST arxitekturası (prinsiplər, HTTP metodları, status kodları), sadə RESTful server/
client (net/http ilə), gorilla/mux ilə funksional RESTful server (subrouters, NotFound/
MethodNotAllowed handler-lər, SQLite backend, login/auth sistemi), Gin vs Gorilla
müqayisəsi, curl ilə testlər, cobra ilə RESTful client və API versiyalama üsulları.

## Əsas fikirlər

### 1. REST Nədir
**REpresentational State Transfer** — web servislərin dizayn arxitekturası; protokol
DEYİL — HTTP(S) üzərində konvensiya. Adətən JSON over HTTP.

**REST prinsipləri:** client-server dizayn, stateless (hər interaksiya müstəqil),
cacheable, uniform interface, layered system.

**HTTP metodları:**
| Metod | Məqsəd |
|---|---|
| POST | yeni resurs YARAT |
| GET | mövcud resursu OXU |
| PUT | tam updated versiya ilə YENİLƏ |
| PATCH | yalnız dəyişikliklərlə YENİLƏ |
| DELETE | SİL |

Go-da konstantlar: http.MethodGet, MethodPost, MethodPut, MethodPatch, MethodDelete...

**Status kodları:** 200 OK, 201 Created, 202 Accepted (uzun əməliyyat), 301 Moved
(rare — versiyalama üstünlük verilir), 400 Bad Request, 401 Unauthorized, 403 Forbidden
(haqq yoxdur), 404 Not Found, 405 Method Not Allowed, 500 Internal Server Error.

### 2. Sadə RESTful Server — net/http ilə
**Kitabdan kod nümunəsi (POST handler):**
```go
type User struct {
    Username string `json:"user"`
    Password string `json:"password"`
}
var DATA = make(map[string]string)

func addHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {              // METOD YOXLAMASI MÜTLƏQ
        http.Error(w, "Error:", http.StatusMethodNotAllowed)
        return
    }
    d, err := io.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "Error:", http.StatusBadRequest)
        return
    }
    err = json.Unmarshal(d, &user)
    if err != nil { /* 400 */ return }
    if user.Username == "" { /* 400 */ return }
    DATA[user.Username] = user.Password
    w.WriteHeader(http.StatusCreated)             // 201 — resurs yaradıldı
}
```
**http.Error() qaydası:** xəta mesajı + status kodu göndərir; SONRA heç nə yazmaq
olmaz — ən sonda çağır.

**DELETE təhlükəsizliyi:** username DƏ, password DƏ uyğun gəlməli — ancaq sonra sil.
curl testi: `-H 'Content-Type: application/json' -d '{...}' -X DELETE`.

### 3. RESTful Client — Test Proqramı
**Kitabdan kod nümunəsi (client pattern):**
```go
func deleteEndpoint(server string, user User) int {
    userMarshall, err := json.Marshal(user)
    u := bytes.NewReader(userMarshall)          // Marshal → Reader
    req, err := http.NewRequest(http.MethodDelete, server+deleteEndPoint, u)
    req.Header.Set("Content-Type", "application/json")   // JSON TÖVSİYƏ EDİLİR
    c := &http.Client{ Timeout: 15 * time.Second }        // timeout MÜTLƏQ
    resp, err := c.Do(req)
    defer resp.Body.Close()
    data, err := io.ReadAll(resp.Body)
    return resp.StatusCode                       // status kodu = uğur/fail
}
```
Client YALNIZ endpoint + metod + data göndərir — serverin implementasiyasını BİLMİR;
eyni endpoint-ləri dəstəkləyən hər hansı serverlə işləyir.

### 4. Funksional Server — gorilla/mux
**Daha professional API variantı (tövsiyə edilən üslub):**
- `GET /users/` — hamısını al; `GET /users/:id` — birini al
- `POST /users/` — yarat; `DELETE /users/:id` — sil
- `PATCH/PUT /users/:id` — yenilə

**gorilla/mux imkanları (default router-dan fərqlər):**
```go
r.HandleFunc("/url", handler).Methods(http.MethodPut)   // metod-metod match!
mux.NotFoundHandler = http.HandlerFunc(DefaultHandler)  // catch-all
mux.MethodNotAllowedHandler = notAllowed                // 405 avtomatik
r.HandleFunc("/users/{id:[0-9]+}", handler)              // PATH DƏYİŞƏNİ + regex!
```

**Subrouters — ümumiləşdirilmiş şərtlər:**
```go
getMux := rMux.Methods(http.MethodGet).Subrouter()
getMux.HandleFunc("/getall", GetAllHandler)
getMux.HandleFunc("/getid/{username}", GetIDHandler)
getMux.HandleFunc("/username/{id:[0-9]+}", GetUserDataHandler)

putMux := rMux.Methods(http.MethodPut).Subrouter()
putMux.HandleFunc("/update", UpdateHandler)

postMux := rMux.Methods(http.MethodPost).Subrouter()
postMux.HandleFunc("/add", AddHandler)

deleteMux := rMux.Methods(http.MethodDelete).Subrouter()
deleteMux.HandleFunc("/username/{id:[0-9]+}", DeleteHandler)
```
Parent subrouter metod şərtini bir dəfə saxlayır — kod təkrarı azalır, match optimallaşır.

### 5. Gin vs Gorilla
- **gorilla/mux** — YALNIZ router (yol + metod + path dəyişənləri)
- **Gin** — router + JSON marshaling + validation + custom response writing = yüksək
  səviyyəli framework (httprouter əsaslı)
**Seçim:** minimalistik/sadə → gorilla/mux; performance-kritik, middleware, tam
framework → Gin. Əvvəl gorilla, bacarmırsa Gin.

### 6. SQLite Backend — restdb.go
**Dizayn qərarı:** handler-lər DB haqqında HİÇ BİLMİR — bütün DB kodu restdb.go-da;
DB dəyişsə handler-lər dəyişməz.

**Cədvəl:**
```sql
CREATE TABLE users (
    UserID INTEGER PRIMARY KEY, username TEXT NOT NULL, password TEXT NOT NULL,
    lastlogin INTEGER, admin INTEGER, active INTEGER
);
```

**DB patterni (kitabdan):**
```go
func OpenConnection() *sql.DB {          // helper — hər funksiya öz bağlantını açır
    db, err := sql.Open("sqlite3", Filename)
    return db
}

func DeleteUser(ID int) bool {
    db := OpenConnection()
    defer db.Close()
    t := FindUserID(ID)                  // mövcudluq yoxlaması əvvəl
    if t.ID == 0 { return false }
    stmt, err := db.Prepare("DELETE FROM users WHERE UserID = $1")
    _, err = stmt.Exec(ID)               // $1 = Exec parametri
    return true
}

func IsUserValid(u User) bool {
    rows, err := db.Query("SELECT * FROM users WHERE username = $1", u.Username)
    // Query + parametr birbaşa; rows.Next/Scan ilə oxu
    if u.Username == temp.Username && u.Password == temp.Password { return true }
    return false
}
```

### 7. Auth Handler-ləri — Admin Yoxlaması
**Kitabdan kod nümunəsi:**
```go
// /add — İKİ User recordu ARRAY kimi: [emr verən, yeni istifadəçi]
users := []User{}
err = json.Unmarshal(d, &users)           // JSON array → slice
if !IsUserAdmin(users[0]) {               // EMR VERƏN admin olmalıdır!
    rw.WriteHeader(http.StatusBadRequest)
    return
}
result := InsertUser(users[1])

// /getall — bir record, admin yoxlaması, nəticə JSON stream:
if !IsUserAdmin(user) { /* 400 */ return }
err = SliceToJSON(ListAllUsers(), rw)     // slice → JSON → ResponseWriter
```

### 8. Graceful Shutdown
```go
go func() {                               // server goroutine-də
    err := s.ListenAndServe()
}()

sigs := make(chan os.Signal, 1)
signal.Notify(sigs, os.Interrupt)
sig := <-sigs                             // siqnal gözlə — main çıxmır
log.Println("Quitting after signal:", sig)
s.Shutdown(nil)                           // GRACEFUL dayanma
```

### 9. curl Testləri — Əsas Komandalar
```bash
curl localhost:1234/time                              # sadə GET
curl -X GET -H 'Content-Type: application/json' \
     -d '{"username": "admin", "password": "x"}' localhost:1234/getall
curl -X POST -H 'Content-Type: application/json' \
     -d '[{"username":"admin",...}, {"username":"packt",...}]' localhost:1234/add
curl -X PUT -H 'Content-Type: application/json' \
     -d '[...]' localhost:1234/update
curl -X DELETE -H 'Content-Type: application/json' \
     -d '{...}' localhost:1234/username/4 -v         # -v status kodu göstərir
```
**Sıra vacibdir:** 405 (metod yanlış) autentikasiyadan ƏV VƏL gəlir — düzgün səbəblə
uğursuzluq.

### 10. Cobra REST Client
**Qlobal flag-lər (root.go):**
```go
rootCmd.PersistentFlags().StringVarP(&username, "username", "u", "username", "...")
rootCmd.PersistentFlags().StringVarP(&password, "password", "p", "admin", "...")
rootCmd.PersistentFlags().StringVarP(&data, "data", "d", "{}", "JSON Record")
rootCmd.PersistentFlags().StringVarP(&SERVER, "server", "s", "http://localhost", "...")
rootCmd.PersistentFlags().StringVarP(&PORT, "port", "P", ":1234", "...")
```
StringVarP — flag dəyəri birbaşa DƏYİŞƏNƏ yazılır (viper BindPFlag-a ehtiyac yoxdur).

**Komanda patterni (list):**
```go
user := User{Username: username, Password: password}
buf := new(bytes.Buffer)                  // Buffer = Reader + Writer!
user.ToJSON(buf)
req, _ := http.NewRequest(http.MethodGet, SERVER+PORT+endpoint, buf)
req.Header.Set("Content-Type", "application/json")
c := &http.Client{Timeout: 15 * time.Second}
resp, err := c.Do(req)
if resp.StatusCode != http.StatusOK { return }
users := []User{}
SliceFromJSON(&users, resp.Body)          // stream decode
data, _ := PrettyJSON(users)
```

**İstifadə nümunələri:**
```bash
./rest-cli list -u admin -p admin
./rest-cli add -u admin -p admin --data '{"Username":"newUser", "Password":"aPass"}'
./rest-cli getid -u admin -p admin --data '{"Username":"newUser"}'   # ID tap
./rest-cli delete -u admin -p admin --data '{"ID":4}'
./rest-cli logged -u admin -p notPass    # 400 Bad Request
```
CLI client curl-dən üstün: status kodlarını INTERPRET edir, datanı preprocess edir;
bahası — inkişaf vaxtı.

### 11. REST API Versiyalama
| Üsul | Nümunə |
|---|---|
| Custom header | version-used: v1 |
| Subdomain | v1.servername |
| Content negotiation | Accept + Content-Type |
| Path | /v1/endpoint |
| Query param | ?version=v1 |

Düzgün cavab YOXDUR — tutarlılıq vacibdir. Müəllif: `/v1/...`, `/v2/...` yol
prefikslərini üstün tutur.

## Əsas terminlər
- REST — REpresentational State Transfer arxitekturası
- Stateless (vəziyyətsiz) — hər sorğu müstəqil
- Uniform Interface (vahid interfeys) — standartlaşdırılmış sorğu/cavab
- JSON Tags — struct↔JSON sahə adları
- http.Error — mesaj + status kodu cavabı
- Subrouter — ümumiləşdirilmiş şərtli nested route (mux)
- Path Variable — {id:[0-9]+} kimi yol dəyişəni
- NotFoundHandler / MethodNotAllowedHandler — mux-in catch-all/405 qatları
- Graceful Shutdown — siqnal ilə təmiz dayanma
- Content Negotiation — Accept header ilə versiya/format seçimi
- API Versioning — REST API-nin versiyalı idarəsi

## Praktik nətidə

(1) REST = HTTP üzərində konvensiya; metod+status kodu semantikası SÖZLEŞMƏDİR —
dokumentasiya edin və tutarlı olun. (2) Handler-də ilk iş: metod yoxlaması (405),
sonra body parse (400), sonra biznes məntiqi. (3) http.Error ən SON çağrılır —
sonrasında yazmaq olmaz. (4) gorilla/mux: metodları route səviyyəsində match et —
handler-də manual yoxlamalar aradan qalxır; subrouter-lər kod təkrarını kəsir.
(5) {id:[0-9]+} — path dəyişənləri + regex validasiyası router-də. (6) DB-ni ayrı
faylda saxla — handler-lər implementasiyadan asılı olmasın. (7) /add üçün İKİ
recordu ARRAY-də göndər: [əmr verən, hədəf] — server əmr verənin hüququnu yoxlayır.
(8) 405 avtorizasiyadan öndə gəlir — düzgün səbəblə reddir. (9) Graceful shutdown:
server-i goroutine-də işə sal + siqnal kanalı + s.Shutdown. (10) Client-də
Content-Type header + Timeout + StatusCode yoxlaması — üçlük HƏMİŞƏ. (11) Versiyalamanı
bir üsulla, hər yerdə eyni tətbiq et.

## Mənbə
Pages: 467-517 (PDF 498-549)
