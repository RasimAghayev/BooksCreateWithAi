# Chapter 3 — Handling HTTP Requests and Routing (HTTP Sorğularının İdarəsi və Routing)

## Bu chapter nədən bəhs edir?
HTTP protokolu, request-response dövrünün, handler/handlerFunc anlayışının, dinamik route-ların, URL parametrlərinin, route qruplaşdırmasının, Gorilla/Mux routerinin, xəta idarəsinin (404/500), rate limiting-in və CRUD endpoint-lərinin praktik tətbiqinə.

## Əsas fikirlər

### 1. HTTP və routing əsasları
**Routing (yönləndirmə):** URI + HTTP metod kombinasiyasına uyğun handler funksiyasının seçilməsi. `/books` GET ≠ `/books` POST — fərqli məntiq.

**HTTP metodları:**
| Metod | Vəzifə | Bookstore nümunəsi |
|-------|--------|---------------------|
| GET | data oxu (dəyişməz) | kitab detalına bax |
| POST | yeni data yarat | review əlavə et, qeydiyyat |
| PUT | resource yenilə (idempotent) | review redaktə et |
| DELETE | resource sil | wishlist-dən sil |
| PATCH | qismən dəyişiklik | yalnız email-i yenilə |
| HEAD | yalnız header-lər | meta-məlumat |
| OPTIONS | kommunikasiya imkanları | capability sorğusu |

### 2. Request-Response dövrü (Alice ssenarisi)
**Request fazası:** brauzer GET sorğusu → header-lər (cookies, user-agent) → server router-i URL+metodla təhlil edir → handler çağırılır.

**Response fazası:** handler DB-dən data alır → status kodu (200/404) + header-lər (Content-Type) + body (JSON/HTML) → brauzer render edir.

### 3. Handler və HandlerFunc
- **http.Handler** interfeysi: `ServeHTTP(ResponseWriter, *Request)` — hər hansı tip tətbiq etsə handler olur
- **http.HandlerFunc** — `func(ResponseWriter, *Request)` tipinin adapteri — ayrıca tip yaratmadan adi funksiya handler olur

```go
func BookListHandler(w http.ResponseWriter, r *http.Request) {
    books := getBooksFromDatabase()
    renderBooks(w, books)
}
http.HandleFunc("/books", BookListHandler)
```

### 4. Middleware pattern-i
**Nədir:** Handler-i handler-in İÇİNƏ saran funksiya — sorğu əsas handler-ə çatmazdan əvvəl emal olunur.

**Kitabdan kod nümunəsi (logging):**
```go
func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Printf("Request received: %s %s", r.Method, r.URL.Path)
        next.ServeHTTP(w, r)   // növbəti handler-ə ötür
    })
}
loggedHandler := LoggingMiddleware(BookListHandler)   // sarğı
```
**Auth middleware:**
```go
func AuthenticationMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if !userIsAuthenticated(r) {
            http.Redirect(w, r, "/login", 302)
            return
        }
        next.ServeHTTP(w, r)
    })
}
bookRouter.Use(AuthenticationMiddleware)   // qrup səviyyəsində tətbiq
```

### 5. Dinamik route-lar və Gorilla/Mux
```go
r := mux.NewRouter()
r.HandleFunc("/book/{id:[0-9]+}", BookDetailHandler)  // regex-li path parametr
r.HandleFunc("/book/add", AddBookHandler).Methods("POST")  // metod məhdudiyyəti
```
**Sub-kod izahı:**
- `{id:[0-9]+}` — path parametr + regex: yalnız rəqəmli ID-lər
- `.Methods("POST")` — bu route yalnız POST-a cavab verir
- `mux.Vars(r)` — handler daxilində parametri çıxarır

### 6. URL parametr növləri
**Query parametrləri** — `?`-dən sonra `key=value&key2=value2`: filtr/sıralama üçün (`/books?genre=fantasy`):
```go
genre := r.URL.Query().Get("genre")
```
**Path parametrləri** — URL yolunun içində: identifikasiya üçün (`/book/5678`):
```go
vars := mux.Vars(r)
bookID := vars["id"]
```
**Birləşmiş istifadə:** `/book/5678/reviews?after=2023-01-01` — hər ikisi eyni sorğuda.

### 7. Route qruplaşdırması (Subrouter)
```go
// Kitab qrupu
bookRouter := r.PathPrefix("/books").Subrouter()
bookRouter.HandleFunc("/", AllBooksHandler)
bookRouter.HandleFunc("/{id:[0-9]+}", BookDetailHandler)
bookRouter.HandleFunc("/{id:[0-9]+}/reviews", BookReviewsHandler)

// İstifadəçi qrupu
userRouter := r.PathPrefix("/users").Subrouter()
userRouter.HandleFunc("/register", RegisterHandler)
userRouter.HandleFunc("/login", LoginHandler)
userRouter.HandleFunc("/profile", ProfileHandler)
userRouter.Use(AuthenticationMiddleware)    // yalnız bu qrupa middleware!

// Sifariş qrupu
orderRouter := r.PathPrefix("/orders").Subrouter()
orderRouter.HandleFunc("/cart", CartHandler)
orderRouter.HandleFunc("/checkout", CheckoutHandler)
orderRouter.HandleFunc("/history", OrderHistoryHandler)
```
**Üstünlüklər:** kod təşkilatı, qrupa-özəl middleware, ardıcıl URL pattern-ləri, scalability.

### 8. Xəta idarəsi — custom 404/500
**404 (Not Found):**
```go
func CustomNotFoundHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusNotFound)
    w.Write([]byte("Kitab tapılmadı — kataloqumuzə baxın!"))
}
r.NotFoundHandler = http.HandlerFunc(CustomNotFoundHandler)
```
**500 (Internal Server Error) — panic yaxalama middleware-i:**
```go
func InternalServerErrorHandler(w http.ResponseWriter, r *http.Request) {
    if rec := recover(); rec != nil {
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte("Xəta baş verdi, düzəldirik..."))
    }
}
r.Use(func(h http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer InternalServerErrorHandler(w, r)   // panic olsa yaxala
        h.ServeHTTP(w, r)
    })
})
```
**Dizayn dərsi:** Xəta səhifələri brend səsi ilə — istifadəçi itən yerdə cəlbedici alternativ göstər (axtarış, oxşar kitablar).

### 9. Rate Limiting və Request Throttling
**Rate limiting:** vahid zamanda maksimum sorğu sayı — DDoS qoruması, ədalətli istifadə.
**Throttling:** sorğuların İSTİQAMƏTİNİ tənzimləmək — sabit tempo saxlamaq.

**Token bucket (golang.org/x/time/rate):**
```go
var limiter = rate.NewLimiter(5, 1)   // 5 sorğu/dəqiqə (rate, burst)

func RateLimit(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if !limiter.Allow() {
            http.Error(w, "Too many requests", http.StatusTooManyRequests)  // 429
            return
        }
        next.ServeHTTP(w, r)
    })
}
```
**Sabit gecikmə throttle:**
```go
func RequestThrottle(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        time.Sleep(200 * time.Millisecond)
        next.ServeHTTP(w, r)
    })
}
r.Use(RateLimit, RequestThrottle)   // global tətbiq
```

### 10. CRUD endpoint-ləri (tam dövrə)
**Create (POST):**
```go
func CreateBookHandler(w http.ResponseWriter, r *http.Request) {
    var newBook Book
    if err := json.NewDecoder(r.Body).Decode(&newBook); err != nil {
        http.Error(w, "Invalid book data", http.StatusBadRequest)   // 400
        return
    }
    BookstoreDB.Create(&newBook)
    w.WriteHeader(http.StatusCreated)                                // 201
    json.NewEncoder(w).Encode(newBook)
}
```
**Read (GET + path parametr):**
```go
func GetBookHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    bookID := vars["id"]
    var fetchedBook Book
    if err := BookstoreDB.First(&fetchedBook, bookID).Error; err != nil {
        http.Error(w, "Book not found", http.StatusNotFound)        // 404
        return
    }
    json.NewEncoder(w).Encode(fetchedBook)
}
```
**Update (PUT):** köhnə qeydi tap → body decode → `Save` → 200 + yenilənmiş JSON.
**Delete (DELETE):** tap → `Delete` → `204 No Content`.

**Status kodları:** 201 (Created), 204 (No Content), 400 (Bad Request), 404, 429 (Too Many), 500.

### 11. SkipClean konfiqurasiyası
`mux.NewRouter().SkipClean(true)` — URL təmizləməsini söndürür (performans); diqqət: təmizlənməmiş URL-lər gözlənilməz davranış yarada bilər.

## Əsas terminlər
- Routing (yönləndirmə) — URI+metod → handler xəritəsi
- HTTP verb (HTTP fel) — GET/POST/PUT/DELETE/PATCH
- Idempotence (eynilik) — eyni sorğunun təkrarı eyni nəticə
- Middleware (arayataq) — handler-ı saran emal qatı
- Path parameter / Query parameter — yol/dot parameter
- Subrouter — prefiksli route qrupu
- Token bucket — rate limit alqoritmi
- Request throttling — tempo tənzimi
- CRUD — Create/Read/Update/Delete

## Praktik nəticə
1. REST dizayn: hər resurs üçün qrup subrouter + metod məhdudiyyəti (`.Methods("POST")`).
2. Middleware-i qrup səviyyəsində tətbiq et (`userRouter.Use`) — hər route üçün ayrı-aqrı YOX.
3. 429 + token bucket istehsalı rate limit; hər user üçün ayrı limiter (IP/user-ID map) düşün.
4. Panic-yaxalayan middleware həmişə olsun — bir sorğu bütün serveri çökməsin.
5. CRUD status kodları: POST→201, DELETE→204, xəta→400/404 — müşahidə oluna bilən API.
6. Query = filtrləmə, Path = identifikasiya — heç vaxt əksinə.

## Mənbə
Pages: 81-114 (PDF səh. 81-114)
