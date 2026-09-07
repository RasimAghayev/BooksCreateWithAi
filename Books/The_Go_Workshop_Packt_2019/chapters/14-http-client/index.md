# Chapter 14 — Using the Go HTTP Client (Go HTTP Klienti)

## Bu fəsil nədən bəhs edir?

net/http client (default vs custom), URL anatomiyası, GET sorğusu (http.Get,
Body oxuma), JSON structured data (Unmarshal), POST sorğusu (http.Post, Marshal,
application/json), fayl upload (multipart form, CreateFormFile, custom Request),
custom header-lər (Authorization) və custom client (Timeout) və POST+GET activity-ləri.

## Əsas fikirlər

### 1. HTTP Klient Nədir
**Nədir:** web server-dən data almaq / server-ə data göndərmək aləti. Browser =
ən məşhur HTTP klient.

**2 əsas sorğu:** GET (data AL — səhifə açmaq), POST (data GÖNDƏR — login formu).

**2 istifadə yolu:**
- Default klient — sadə, sürətli başlanğıc
- Custom klient (http.Client{}) — tam nəzarət (timeout və s.)

### 2. URL Anatomiyası
```
https://example.com/downloads?filter=latest&os=windows
  │          │          │              │
Protocol   Hostname     URI        Query Parameters
```
- `?` — query başlanğıcı; `&` — parametrlər arası ayırıcı
- Protocol: HTTP/HTTPS; URI: resurs yolu

### 3. GET Sorğusu
**Kitabdan kod nümunəsi:**
```go
r, err := http.Get("https://www.google.com")    // GET göndər
if err != nil { log.Fatal(err) }
defer r.Body.Close()                             // Body HƏMİŞƏ bağla
data, err := ioutil.ReadAll(r.Body)             // body-ni oxu
return string(data)                              // HTML string kimi
```
Alınan data-nı response.html-ə yazsan → Google home page — brauzerin arxa planı
məhz budur.

### 4. Structured Data — JSON Parse
**Server tərəfi (sadə):**
```go
type server struct{}
func (srv server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    msg := `{"message": "hello world"}`
    w.Write([]byte(msg))
}
http.ListenAndServe(":8080", server{})
```

**Client — Unmarshal (kitabdan):**
```go
type messageData struct {
    Message string `json:"message"`
}

r, _ := http.Get("http://localhost:8080")
defer r.Body.Close()
data, _ := ioutil.ReadAll(r.Body)

message := messageData{}
err = json.Unmarshal(data, &message)      // JSON → struct
fmt.Println(data.Message)                 // "hello world"
```
Pipeline: GET → ReadAll → Unmarshal → struct istifadəsi.

### 5. POST Sorğusu — Data Göndərmə
**Kitabdan kod nümunəsi:**
```go
type messageData struct {
    Message string `json:"message"`
}

func postDataAndReturnResponse(msg messageData) messageData {
    jsonBytes, _ := json.Marshal(msg)          // struct → JSON bytes

    r, err := http.Post("http://localhost:8080",
        "application/json",                     // Content-Type
        bytes.NewBuffer(jsonBytes))            // body = buffer
    if err != nil { log.Fatal(err) }
    defer r.Body.Close()

    data, _ := ioutil.ReadAll(r.Body)           // cavabı oxu
    message := messageData{}
    json.Unmarshal(data, &message)             // cavabı parse et
    return message                              // echo回来了
}
msg := messageData{Message: "Hi Server!"}
```
**Tələlər:** Marshal → bytes.Buffer → http.Post(url, contentType, buffer);
server decode edib eyni struct-u qaytarır (NewDecoder/Decode).

### 6. Fayl Upload — Multipart Form
**Kitabdan kod nümunəsi (client):**
```go
func postFileAndReturnResponse(filename string) string {
    fileDataBuffer := bytes.Buffer{}
    multipartWriter := multipart.NewWriter(&fileDataBuffer)

    file, err := os.Open(filename)              // lokal faylı aç
    formFile, err := multipartWriter.CreateFormFile(
        "myFile", file.Name())                  // form sahəsi yarat
    io.Copy(formFile, file)                      // fayl bytes-i forma KOPYALA
    multipartWriter.Close()                      // writer-i BİTİR

    // Custom Request — http.Post qısa yolu KİFAYƏT DEYİL:
    req, err := http.NewRequest("POST",
        "http://localhost:8080", &fileDataBuffer)
    req.Header.Set("Content-Type",
        multipartWriter.FormDataContentType())  // MULTIPART başlığı MÜTLƏQ

    response, err := http.DefaultClient.Do(req)  // client.Do icra
    defer response.Body.Close()
    data, _ := ioutil.ReadAll(response.Body)
    return string(data)
}
```
**Server tərəfi:** `r.FormFile("myFile")` → uploadedFile + header → ReadAll →
disk-ə yaz.

**Addımlar:** buffer + multipart writer → CreateFormFile → io.Copy → Close →
NewRequest → FormDataContentType header → Do.

### 7. Custom Headers + Custom Client
**Kitabdan kod nümunəsi:**
```go
client := http.Client{Timeout: 11 * time.Second}    // ÖZ klientin

req, err := http.NewRequest("POST", "http://localhost:8080", nil)
req.Header.Set("Authorization", "superSecretToken")  // TOKEN header

resp, err := client.Do(req)                          // icra
defer resp.Body.Close()
```
**Server yoxlaması:**
```go
auth := r.Header.Get("Authorization")
if auth != "superSecretToken" {
    w.WriteHeader(http.StatusUnauthorized)     // 401
    w.Write([]byte("Authorization token not recognized"))
    return
}
```
**Timeout dərsi:** server 10 saniyə yatırdısa, client timeout-u 11 san → çatır;
<10 san → xəta. Auth token API-lərin ən common tələbidir.

### 8. Activity-lər
- **14.01:** GET → JSON array of names → Electric/Boogaloo SAYĞI
- **14.02:** POST ilə ad əlavəsi → GET ilə yoxlama (data-nın write+read-back
  pattern-i — professional test əsası)

## Əsas terminlələr
- HTTP Client — server ilə kommunikasiya aləti (browser = örnək)
- Default vs Custom Client — http.Get qısa yolu vs http.Client{...}
- GET/POST — data alma / göndərmə
- URL: Protocol/Hostname/URI/Query Parameters
- r.Body — cavabın məzmun axını; defer Close MÜTLƏQ
- Content-Type — məzmun tipi (application/json, multipart/form-data)
- bytes.NewBuffer — POST body wrapper
- Multipart Form — fayl upload formatı
- CreateFormFile — formda fayl sahəsi yaratma
- FormDataContentType — multipart başlığı generatoru
- http.NewRequest — custom request (header/ body nəzarəti)
- client.Do(req) — custom klientdə icra
- Authorization Header — token-based identifikasiya
- http.Client{Timeout} — sorğu vaxt həddi
- StatusUnauthorized — 401

## Praktik nətidə

(1) http.Get/Post qısa yollar — sadə hallar; NewRequest+Do — header/timeout
nəzarəti lazım olanda. (2) r.Body HƏMİŞƏ defer Close ilə — resurs sizintisinin ən
common mənbəyi. (3) JSON cavab: ReadAll → Unmarshal → struct; tag-lərlə mapinqi
dəqiqləşdir. (4) POST: Marshal → buffer → content-type MÜTLƏQ təyin et. (5) Fayl
upload: multipart writer + CreateFormFile + io.Copy; FormDataContentType başlığı
OLMADAN server formu TANIMIR. (6) Custom klientdə timeout qoy — yavaş serverlər
proqramı ASMAMALI. (7) Authorization header — API autentifikasiyasının standart
yolu; token yoxdursa 401. (8) Write+read-back: POST sonra GET — datanın həqiqətən
düşdüyünü YOXLA. (9) URL query: ?param=dəyər&param2=dəyər2. (10) Custom client
DEFAULT-dan fərqli olmalıdırsa yalnız onda yaradılmağa dəyər — amma timeout HƏMİŞƏ
qoyulmalıdır.

## Mənbə
Pages: 493-513 (PDF 526-547)
