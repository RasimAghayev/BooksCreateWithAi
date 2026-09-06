# Chapter 8 — Frontend and Backend Communication (Frontend və Backend Ünsiyyəti)

## Bu chapter nədən bəhs edir?
Frontend-backend əməkdaşlığına, REST prinsiplərinin təkrarına, ilk API-nin qurulmasına, JSON encode/decode-a, DB-dən data endpoint-inə, paralel data fetch-ə, API token autentifikasiyasına (JWT) və xarici API inteqrasiyasına (Stripe, book API-ləri).

## Əsas fikirlər

### 1. Frontend-backend bölgüsü
**Backend:** biznes məntiqi, DB, security, scalability (Go net/http). **Frontend:** istifadəçi görünüşü (HTML/CSS/JS), rendering.
**Körpü:** API-lər — sorğu-nəticə mübadiləsinin standartlaşdırılmış yolu. SPA-lar (React/Vue) API-lərlə klassik bölgünü bulandırıb.

### 2. REST prinsipləri (təkrar, API kontekstində)
- **Statelessness:** hər sorğu özündə tam məlumat daşıyır
- **Client-server:** UI və data concerns ayrılığı
- **Uniform interface:** platformadan asılı olmayaraq vahid interaksiya (mobile/web eyni GET)
- **Layered system:** hər qat yalnız qonşu qatı görür
- **Representations:** JSON/XML resurs təsvirləri

### 3. İlk API endpoint (addım-addım)
```go
package main

import (
    "net/http"
    "encoding/json"
)

type Book struct {
    ID     int     `json:"id"`      // struct tag — JSON sahə adı
    Title  string  `json:"title"`
    Author string  `json:"author"`
    Price  float64 `json:"price"`
}

var books = []Book{
    {ID: 1, Title: "Go Basics", Author: "John Doe", Price: 10.99},
    {ID: 2, Title: "Advanced Go", Author: "Jane Smith", Price: 15.49},
}

func fetchBookDetails(w http.ResponseWriter, r *http.Request) {
    bookID := r.URL.Query().Get("id")           // query parametr
    for _, book := range books {
        if strconv.Itoa(book.ID) == bookID {
            json.NewEncoder(w).Encode(book)     // JSON cavab
            return
        }
    }
    w.WriteHeader(http.StatusNotFound)
    json.NewEncoder(w).Encode("Book not found")
}

func main() {
    http.HandleFunc("/book", fetchBookDetails)
    http.ListenAndServe(":8080", nil)
}
```
**Test:** `http://localhost:8080/book?id=1` → JSON.
**Production qeydləri:** routing paketi, CORS, logging middleware, real DB, mənalı xəta cavabları.

### 4. JSON — data mübadilə formatı
**Struktur:** `{}` obyekt (key-value), `[]` massiv, string/number/true/false/null.

**Encoding (struct → JSON):**
```go
book := Book{ID: 1, Title: "Go Basics", Author: "John Doe", Price: 10.99}
jsonData, err := json.Marshal(book)
fmt.Println(string(jsonData))
```
**Decoding (JSON → struct):**
```go
jsonString := `{"id":1,"title":"Go Basics","author":"John Doe","price":10.99}`
var book Book
err := json.Unmarshal([]byte(jsonString), &book)
fmt.Printf("Book Title: %s\n", book.Title)
```
**Struct tags:** `json:"id"` — sahə adlarını idarə edir; nested strukturlar/slice-lar mürəkkəb iyerarxiyaları təmsil edir.

### 5. DB-backed endpoint
```go
func getBooks(w http.ResponseWriter, r *http.Request) {
    rows, err := db.Query("SELECT id, title, author, price FROM books")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var books []Book
    for rows.Next() {
        var b Book
        if err := rows.Scan(&b.ID, &b.Title, &b.Author, &b.Price); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        books = append(books, b)
    }

    jsonData, err := json.Marshal(books)
    if err != nil { /* 500 */ }
    w.Header().Set("Content-Type", "application/json")
    w.Write(jsonData)
}
```
**Zəncir:** Query → Scan loop → Marshal → Content-Type + Write.
**Production:** prepared statements, connection reuse, error handling, concurrency.

### 6. Paralel data fetch (goroutine + channel)
**Niyə:** Sekvensial fetch = bir API bitənə qədər gözlə; paralel = hamısı eyni anda.

**Kitabdan kod nümunəsi:**
```go
func fetchData(url string, ch chan<- string) {
    // ... fetch məntiqi
    ch <- data
}

func main() {
    urls := []string{"https://api1.example.com", "https://api2.example.com"}
    ch := make(chan string, len(urls))     // buffered
    for _, url := range urls {
        go fetchData(url, ch)               // paralel sorğular
    }
    for range urls {
        data := <-ch                       // nəticələr
        fmt.Println(data)
    }
}
```
**Error handling pattern-i — Result strukturu:**
```go
type Result struct {
    Data string
    Err  error
}
ch := make(chan Result)
// hər goroutine data+err birlikdə göndərir
```
**Xəbərdarlıq:** minlərlə nəzarətsiz goroutine sistemi boğar; error-lar ayrıca fail edə bilər — Result pattern məcburidir.

### 7. API autentifikasiyası — JWT token
**Auth vs Authz:** Authentication = kimliy yoxla; Authorization = nə icazəli.

**Token generasiyası (jwt-go):**
```go
func GenerateToken(userID string) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "userID": userID,
    })
    return token.SignedString([]byte("your_secret_key"))
}
```
**Token validasiyası:**
```go
func ValidateToken(encodedToken string) (string, error) {
    token, err := jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
        return []byte("your_secret_key"), nil
    })
    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        userID := claims["userID"].(string)
        return userID, nil
    }
    return "", err
}
```
**Auth middleware:**
```go
func AuthenticationMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")   // header-dən token
        if _, err := ValidateToken(token); err != nil {
            http.Error(w, "Invalid token", http.StatusUnauthorized)  // 401
            return
        }
        next(w, r)
    }
}
http.HandleFunc("/protectedEndpoint",
    AuthenticationMiddleware(HandleProtectedEndpoint))
```
**Token riskləri:** leak olsa sui-istifadə → HTTPS, qısa expire, breach halında etibarsızlaşdırma. **Alternativ:** Auth0/Okta hazır həllər.

### 8. Xarici API inteqrasiyası
**Seçim meyarları:** etibarlılıq, xərc, scalability, uyğunluq, dokumentasiya, community, reputasiya.

**Sadə GET:**
```go
resp, err := http.Get("https://bookreviewsapi.com/book/{bookID}")
```
**Stripe payment gateway:**
```go
stripe.Key = "YOUR_SECRET_KEY"           // ENV VARDAN, hardcoded YOX!
ch, err := charge.New(&stripe.ChargeParams{
    Amount:   stripe.Int64(2000),        // $20.00
    Currency: stripe.String(string(stripe.CurrencyUSD)),
    Source:   &stripe.SourceParams{Token: stripe.String("TOKEN_FROM_CHECKOUT")},
})
```
**Book summary servisi:**
```go
func GetBookSummary(bookID string) (string, error) {
    url := fmt.Sprintf("https://booksummaryapi.com/summary/%s", bookID)
    resp, err := http.Get(url)
    if err != nil { return "", err }
    defer resp.Body.Close()
    body, _ := ioutil.ReadAll(resp.Body)
    return string(body), nil
}
```
**Rate limit idarəsi:** `time.Ticker` və ya `golang.org/x/time/rate` — xarici limitləri qorumaq üçün.
**Monitorinq:** xəta/unexpected response logging; API struktur/xərc dəyişikliklərini periodik audit et.

## Əsas terminlər
- API (Application Programming Interface)
- REST — Roy Fielding, 2000
- Statelessness / Uniform interface / Layered system
- JSON (JavaScript Object Notation)
- Struct tag — `json:"field_name"`
- json.Marshal / Unmarshal / NewEncoder
- JWT (JSON Web Token) — HS256 imzalı
- Auth middleware — Header token yoxlaması
- Payment gateway — Stripe charge
- Rate limit / quota — xarici API məhdudiyyətləri

## Praktik nəticə
1. API struct-larında həmişə JSON tag-lər — sahə adlarını frontend-lə sinxron saxla.
2. Paralel fetch pattern-i: buffered channel + N goroutine + Result{Data, Err} — xətalar itmir.
3. Qorunan endpoint-lər: JWT + auth middleware təbəqəsi; secret-i env-də saxla.
4. Xarici API key-lər heç vaxt koddə olmasın — env var; rate limit alətləri ilə limitləri qoru.
5. Content-Type header-i JSON cavablarda mütləq; 401/500/404 status kodları düzgün.
6. Xarici inteqrasiyalar mü脆弱dır — logging + periodik audit dövrü qur.

## Mənbə
Pages: 206-230 (PDF səh. 206-230)
