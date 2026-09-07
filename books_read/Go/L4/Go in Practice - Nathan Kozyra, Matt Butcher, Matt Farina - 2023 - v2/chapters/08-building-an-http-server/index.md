# Chapter 8 — Building an HTTP server (HTTP Server Qurulması)

## Bu chapter nədən bəhs edir?

Metod əsaslı routing (Go 1.22 "GET /path" pattern-i), server konfiqurasiyası (timeouts,
TimeoutHandler), context ilə state (middleware), struct handler-lər (DB), path dəyişənləri
(PathValue), nested mux, query/form data, cookie, JWT (golang-jwt) və basic auth.

## Əsas fikirlər

### 1. Metod Routing (Go 1.22)
**Kitabdan kod nümunəsi:**
```go
http.HandleFunc("GET /comments", getComments)     // metodu AYRI handler
http.HandleFunc("POST /comments", postComments)   // 405 avtomatik digər metodlara
```
- Standart olmayan metod da leqaldır ("FOOBAR /comments") — maraqdan başqa
- Manual upsert pattern: `r.Method == http.MethodPost` + ID 0-ı create/1+ update-ə
  şaxələndir; PUT/PATCH birbaşa update
- REST = sadəcə verb deyil (Fielding dissertasiyası); OPTIONS self-documenting API üçün

### 2. Server Konfiqurasiyası
```go
muxer := http.NewServeMux()
muxer.HandleFunc("GET /timeout", timeoutHandler)
server := http.Server{
    Addr:         ":8000",
    ReadTimeout:  1 * time.Second,    // 503 qaytarır
    WriteTimeout: 2 * time.Second,    // context deadline — 200 amma boş cavab
    Handler:      muxer,
}
server.ListenAndServe()
```
`http.TimeoutHandler(muxer, 2*time.Second, "request took too long")` — 503 ilə global
wrap. Handler səviyyəsində `r.Context()` ilə fərdi timeout.

### 3. Middleware və Context State
**Kitabdan kod nümunəsi:**
```go
var validAgent = regexp.MustCompile(`(?i)(chrome|firefox)`)
func uaMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        userAgent := r.UserAgent()
        if !validAgent.MatchString(userAgent) {
            w.WriteHeader(http.StatusBadRequest)
            return
        }
        ctx := context.WithValue(r.Context(), "agent", userAgent)
        r = r.WithContext(ctx)
        next(w, r)
    }
}
// handler-də: r.Context().Value("agent").(string)
```
İstifadə: `http.HandleFunc("GET /withcontext", uaMiddleware(uaStatusHandler))` —
auth üçün eyni wrap istənilən handler-ə tətbiq olunur.

### 4. Struct Handler — DB Kənara Yayılması
```go
type serverControl struct { db *pgx.Conn }
func main() {
    var sc serverControl
    sc.db, _ = pgx.Connect(context.Background(), "postgres://localhost:5432")
    http.HandleFunc("GET /database", sc.databaseHandler)   // metod = db çıxışı
}
func (sc serverControl) databaseHandler(w http.ResponseWriter, r *http.Request) {
    rows, err := sc.db.Query(`select user_id, comment from comments limit $1`, 5)
    defer rows.Close()
    for rows.Next() {
        var c comment
        rows.Scan(&c.UserID, &c.Comment)
        comments = append(comments, c)
    }
    w.Header().Set("Content-type", "application/json")
    w.Write(output)
}
```
Üstünlük: constructor-da bağlanan resurs, metodda hazır; risk — handler bağı bağlaya bilər.

### 5. Path Dəyişənləri və Nested Mux
```go
http.HandleFunc("GET /comments/{id}", getComment)
commentID, err := strconv.Atoi(r.PathValue("id"))   // {id} dəyəri
if commentID == 0 || len(comments) < commentID { w.WriteHeader(404); return }
```
Precedence: path dəyişənləri XÜSUSİLİK üzrə — `/comments/{id}` > `/comments/uzunpath`.
Sub-router-lar:
```go
mainRouter.Handle("/users/", http.StripPrefix("/users", usersRouter))
mainRouter.Handle("/comments/", http.StripPrefix("/comments", commentsRouter))
```
Framework-lər: chi (qrup middleware), httprouter (sürət, exact match), Gin (httprouter
+ Sinatra). Sadə halda ServeMux birinci seçim.

### 6. Query Parametrlər
```go
params := r.URL.Query()
if username := params.Get("username"); username != "" { /* filter */ }
if search := params.Get("search"); search != "" {
    re := regexp.MustCompile(search)   // TƏHLÜKƏ: user regex → ReDoS!
}
```
- `Values.Get` boş string qaytarır; `Has` boş dəyərdə də true — `!= ""` yoxlaması düzgün
- **Regex-i user-dən qəbul etmək QADAĞAYAXIN** (DoS; Go linear-time etsə də)
- OWASP Top 10-u nəzərdən keçir; input sanitize (html/template — Ch9)

### 7. Form Data
```go
func postHandler(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()                          // url.Values strukturuna çevirir
    username := r.Form.Get("username")
    commentText := r.Form.Get("comment")
    comments = append(comments, comment{username, text, time.Now().Format(time.RFC3339)})
    http.Redirect(w, r, "/comments", http.StatusFound)
}
```
HTML form: `<form method="POST" action="/comments">` + name atributlu input/textarea.
Validasiya: `r.Form.Has("username")` — amma hücumkara strukturu ifşa etməmək üçün
sadə 400 kifayət. Göstərilən input sanitize-sizdir — XSS (html/template həlli).

### 8. Cookie
```go
usernameCookie, err := r.Cookie("username")
if err == nil { username = usernameCookie.Value }
http.SetCookie(w, &http.Cookie{
    Name: "username", Value: username,
    Expires: time.Now().Add(24 * time.Hour),
})
```
Stateless HTTP-də state; ephemerik dəyərlər üçün; təkzibedilməz ("You" nümunəsi) — amma
cleartext cookie spoofing-ə açıqdır.

### 9. JWT (JSON Web Token)
**Nədir:** İçində identitet daşıyan İMZALANMIŞ token; cookie/session-ID-dən fərqli olaraq
backend lookup tələb etmir.

**Kitabdan kod nümunəsi (golang-jwt/jwt/v5):**
```go
var SIGNING_KEY = []byte("this-value-should-be-secret")   // .pem-dən oxunur
claims := jwt.RegisteredClaims{
    ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
    IssuedAt:  jwt.NewNumericDate(time.Now()),
    Subject:   "nobody@example.com",
}
token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
ss, _ := token.SignedString(SIGNING_KEY)

// doğrulama:
token, err := jwt.Parse(signed, func(t *jwt.Token) (interface{}, error) {
    return SIGNING_KEY, nil
})
if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
    log.Println(claims["sub"])
}
```
Yanlış SIGNING_KEY → parse xətası (imitasiya mümkünsüz). Çatışmazlıq: invalidation
(revocation) — yalnız hər requestdə external yoxlama ilə (benefit itir). **Security
qaydası:** identitet/auth üçün öz-özünə yazma YOX — vetted kitabxana (jwt.io siyahısı);
Auth0 kimi SSO provayderlərində qarşılaşacaqsan.

### 10. Basic Auth
```go
var validUsers = map[string]string{"bill": "abc123"}
username, password, auth := r.BasicAuth()
if !auth || !login(username, password) {
    w.WriteHeader(http.StatusUnauthorized)
    return
}
```
cURL test: `curl -X POST -u "bill:abc124" ...` → 401. Production-da parollar one-way
encrypt olunmuş xarici mənbədən (Ch9-10). Status kodları konstant adları ilə yaz
(StatusNotFound, 418 Teapot kimi ekzotikləri əzbərləmə).

## Əsas terminlələr
- Method-prefixed route — "GET /path" patterni (Go 1.22)
- TimeoutHandler — 503 qaytaran sarğı
- r.WithContext / context.WithValue — request-də state daşıma
- PathValue — `{id}` path dəyişəninin oxunması
- StripPrefix — sub-router-ın prefiks kəsilməsi
- ParseForm / url.Values — form data parse
- http.SetCookie / r.Cookie — cookie yazı/oxu
- JWT / RegisteredClaims — imzalı identitet tokeni
- BasicAuth — header əsaslı username/password

## Praktik nətidə

Server qərarları: (1) metod+path route-ları ServeMux patternində bir sətirdə; (2) server
struct-da timeout-lar MÜTLƏQ (+ TimeoutHandler global); (3) cross-cutting yoxlama →
middleware + WithValue; (4) DB/resurs → constructor struct + metod handler (qapanma
riskinə diqqət); (5) path dəyişəni → PathValue + validasiya + 404; (6) query/form →
Get + `!= ""` (Has YOX); user regex QADAĞAN; (7) state → cookie (ephemerik) → JWT
(imzalı, lookup-sız; revocation problemi) → basic auth (ən sadə); (8) status kodlarını
konstant adları ilə qaytar; (9) sanitize edilməmiş input-un çapı — XSS; html/template
 növbəti fəslin mövzusudur.

## Mənbə
Pages: 195-219 (PDF 216-240)
