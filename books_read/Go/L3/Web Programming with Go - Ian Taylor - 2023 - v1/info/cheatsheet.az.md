# Web Programming with Go — Cheat Sheet (Azərbaycanca)

> **Kitab:** Web Programming with Go — Ian Taylor, GitforGits, 2023 · 🎯 Intermediate (3/5)
> Bookstore layihəsi boyunca istifadə olunan bütün kod/konfiqurasiya nümunələri — bir baxışda.

---

## 1. Layihə qurulumu (Ch 1-2)
```bash
mkdir gitforgits-bookstore && cd gitforgits-bookstore
go mod init gitforgits-bookstore
go get -u github.com/gorilla/mux      # router
go get -u gorm.io/gorm                # ORM
go get -u github.com/joho/godotenv    # .env
go get golang.org/x/crypto/bcrypt     # parol hash
go get golang.org/x/time/rate         # rate limit
go get -u github.com/lib/pq           # PostgreSQL driver
```

## 2. Layihə strukturu (Ch 2)
```
cmd/server/          # main
internal/handlers/   # HTTP handler-lər
internal/middleware/ # logging, CORS, auth
internal/config/     # konfiqurasiya
pkg/models/         # Book, User strukturları
web/templates/      # HTML şablonlar
web/static/          # CSS/JS
database/migrations/ # SQL skriptləri
```

## 3. Minimal server + App obyekti (Ch 1-2)
```go
type App struct {
    Router *mux.Router
    DB     *sql.DB
    Config map[string]string
}
func (a *App) Initialize(config map[string]string) {
    a.DB, _ = sql.Open("postgres", config["database"])
    a.Router = mux.NewRouter()
    a.initializeRoutes()
}
http.HandleFunc("/", homeHandler)
http.ListenAndServe(":8080", nil)
```

## 4. Routing (Ch 3)
```go
r := mux.NewRouter()
r.HandleFunc("/book/{id:[0-9]+}", BookDetailHandler).Methods("GET")
bookRouter := r.PathPrefix("/books").Subrouter()   // qrup
bookRouter.HandleFunc("/{id:[0-9]+}/reviews", BookReviewsHandler)
userRouter.Use(AuthenticationMiddleware)            // qrupa middleware
vars := mux.Vars(r); bookID := vars["id"]           // path param
genre := r.URL.Query().Get("genre")                 // query param
r.NotFoundHandler = http.HandlerFunc(Custom404)
```

## 5. Middleware pattern-i (Ch 3)
```go
func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Printf("Request: %s %s", r.Method, r.URL.Path)
        next.ServeHTTP(w, r)
    })
}
func RateLimit(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if !limiter.Allow() {
            http.Error(w, "Too many requests", http.StatusTooManyRequests) // 429
            return
        }
        next.ServeHTTP(w, r)
    })
}
```

## 6. Templating (Ch 4)
```go
// Keshlənmiş şablonlar (başlanğıcda BİR dəfə):
var templateCache = template.New("").Delims("{{", "}}")
templateCache.ParseGlob("./templates/*.gohtml")
templateCache.ExecuteTemplate(w, "booklist.gohtml", books)
```
```html
{{define "base"}}{{template "header" .}}{{template "content" .}}{{end}}
{{range .Books}}<h2>{{.Title}}</h2>{{end}}
{{if .IsRare}}...{{end}}
```

## 7. Form emalı (Ch 4)
```html
<form action="/register" method="post">
  <input type="email" name="email" required>
</form>
```
```go
username := r.FormValue("username")
```

## 8. Database (Ch 5)
```go
db, _ := sql.Open("postgres", connStr)
db.SetMaxOpenConns(100); db.SetMaxIdleConns(50)
db.SetConnMaxLifetime(time.Minute * 5)
rows, _ := db.Query("SELECT title, author FROM books WHERE id = $1", id)  // parametrli!
defer rows.Close()
rows.Scan(&title, &author)
result, _ := db.Exec("INSERT INTO books (title) VALUES ($1)", t)
tx := db.Begin(); tx.Commit() / tx.Rollback()
db.Ping()  // health
```
**Sxem:** books, genres, users (password_hash!), reviews (CHECK rating 1-5), transactions, inventory (restock_threshold).

## 9. Concurrency (Ch 6)
```go
ch := make(chan []Book, len(genres))          // buffered
go searchBooksByGenre(genre, ch)               // paralel
books := <-ch                                  // nəticə

var wg sync.WaitGroup
wg.Add(1); go func() { defer wg.Done(); ... }(); wg.Wait()

var mu sync.Mutex
mu.Lock(); inventory[b] -= n; mu.Unlock()

select { case o := <-orders: ...; default: }   // non-blocking
```

## 10. Sessiya və təhlükəsizlik (Ch 7)
```go
http.SetCookie(w, &http.Cookie{
    Name: "session_token", Value: token,
    HTTPOnly: true,                            // XSS qoruması
    Secure: true,                              // yalnız HTTPS
    SameSite: http.SameSiteStrictMode,          // CSRF qoruması
    Expires: time.Now().Add(2 * time.Hour),
})
hashed, _ := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
bcrypt.CompareHashAndPassword(hashed, []byte(pw))  // login yoxlaması
// Salt + Pepper: salt + password + pepper kombinə → hash
```

## 11. JSON API (Ch 8)
```go
type Book struct {
    ID    int     `json:"id"`      // struct tag MÜTLƏQ
    Title string  `json:"title"`
}
json.NewEncoder(w).Encode(book)             // yaz
json.NewDecoder(r.Body).Decode(&newBook)    // oxu
jsonData, _ := json.Marshal(books)          // marshal
json.Unmarshal(data, &book)
w.Header().Set("Content-Type", "application/json")
```

## 12. JWT Auth (Ch 8)
```go
token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"userID": id})
tokenStr, _ := token.SignedString([]byte(secret))
// Header: Authorization: <token>
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc { ... 401 ... }
```

## 13. Xarici API (Ch 8)
```go
resp, err := http.Get("https://api.example.com/data")
defer resp.Body.Close()
stripe.Key = os.Getenv("STRIPE_KEY")        // hardcoded YOX!
```

## 14. Testing (Ch 9)
```go
func TestAdd(t *testing.T) {
    tests := []struct{ name string; a, b, want int }{
        {"pos", 3, 4, 7}, {"neg", -3, -4, -7},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := Add(tt.a, tt.b); got != tt.want {
                t.Errorf("Add(%d,%d)=%d; want %d", tt.a, tt.b, got, tt.want)
            }
        })
    }
}
t.Parallel()        // paralel test
t.Errorf / Fatalf / Skipf
func BenchmarkX(b *testing.B) { for i := 0; i < b.N; i++ { ... } }
```

## 15. Mock pattern (Ch 9)
```go
type Database interface { FetchBookByID(id int) (Book, error) }
type MockDatabase struct{ books []Book }
func (m *MockDatabase) FetchBookByID(id int) (Book, error) { ... }
// test: FetchBookDetails(mockDB, 1)
```

## 16. pprof (Ch 9)
```go
import _ "net/http/pprof"
pprof.StartCPUProfile(f); defer pprof.StopCPUProfile()
runtime.SetBlockProfileRate(1)         // block
pprof.WriteHeapProfile(f)             // memory
// go tool pprof cpu.pprof → top
```

## 17. Konfiqurasiya (Ch 2)
```go
godotenv.Load(".env")                    // lokal dev
os.Getenv("DB_PASSWORD")                 // production env
// .env git-ə commit EDİLMƏZ!
```

---

## Xəta-düzəltmə cədvəli (Ch 9)
| Xəta | Həll |
|------|------|
| N+1 query | JOIN / batch |
| Data race | Mutex / channel |
| Memory leak | pprof + kanal close |
| Input validation | hər input yoxla |
| Error handling | user-ə generic, log-a detal |
| Hardcoded config | env var |
| Unsafe concurrency | mutex ilə qoru |
| Pagination yox | LIMIT/OFFSET |
