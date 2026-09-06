# Chapter 8 — Working with web services (Technique 49-55)

## Bu chapter nədən bəhs edir?

REST API istifadəsi: HTTP client (default + custom), timeout detection (net.Error, url.Error), resumable download (Range header), custom JSON error ötürülməsi və oxunması, arbitrary JSON parsing (interface{}), API versioning (URL path vs Accept content type).

## Əsas fikirlər

### 1. HTTP client əsasları
```go
// Helper funksiyalar:
res, _ := http.Get("http://goinpracticebook.com")
b, _ := ioutil.ReadAll(res.Body)
res.Body.Close()

// İstənilən metod (DELETE və s.):
req, _ := http.NewRequest("DELETE", "http://example.com/foo/bar", nil)
res, _ := http.DefaultClient.Do(req)
fmt.Printf("%s", res.Status)

// Custom client:
cc := &http.Client{Timeout: time.Second}    // 1 saniyə timeout!
res, err := cc.Get("http://goinpracticebook.com")
```
- **Ayrılma:** Request (nə) + Client (necə) — hər ikisi customize olunur
- Default client: redirects (10-a qədər), cookies, transport idarəsi
- `Client{Timeout: ...}` — request + BODY OXUMA hamısı bu pəncərədə

### TECHNIQUE 49: Timeout detection
**Problem:** net-in `Timeout()` metodu VAR amma hamı tətbiq etmir; url.Error fərqli tipdir; set olunmamış timeout-da baş verə bilər.

**Kitabdan kod nümunəsi:**
```go
func hasTimedOut(err error) bool {
    switch err := err.(type) {                 // TYPE SWITCH
    case *url.Error:                          // url paketindən
        if err, ok := err.Err.(net.Error); ok && err.Timeout() {
            return true                       // daxili net.Error yoxla
        }
    case net.Error:                           // birbaşa net.Error
        if err.Timeout() {
            return true
        }
    case *net.OpError:                        // əməliyyat xətası
        if err.Timeout() {
            return true
        }
    }
    // Set olunmamış timeout — string yoxlaması (son çarə):
    errTxt := "use of closed network connection"
    if err != nil && strings.Contains(err.Error(), errTxt) {
        return true
    }
    return false
}

// İstifadə:
res, err := http.Get("http://example.com/test.zip")
if err != nil && hasTimedOut(err) {
    fmt.Println("A timeout error occured")
    return
}
```
- 3 tip + string fallback — bütün timeout halları örtülür

### TECHNIQUE 50: Resumable download (Range header)
**HTTP 1.1 (1999) Range standartı:** server faylın HİSSƏSİNİ verə bilər.

**Kitabdan kod nümunəsi:**
```go
func main() {
    file, err := os.Create("file.zip")
    defer file.Close()
    err = download(location, file, 100)       // 100 retry
    fi, _ := file.Stat()
    fmt.Printf("Got it with %v bytes downloaded", fi.Size())
}

func download(location string, file *os.File, retries int64) error {
    req, err := http.NewRequest("GET", location, nil)
    if err != nil {
        return err
    }
    fi, err := file.Stat()
    if err != nil {
        return err
    }
    current := fi.Size()
    if current > 0 {                                    // lokal faylda data varsa
        start := strconv.FormatInt(current, 10)
        req.Header.Set("Range", "bytes="+start+"-")    // "bayt N-dən sona qədər"
    }
    cc := &http.Client{Timeout: 5 * time.Minute}
    res, err := cc.Do(req)
    if err != nil && hasTimedOut(err) {                 // timeout → RETRY
        if retries > 0 {
            return download(location, file, retries-1)  // REKURSİV — Range ilə davam!
        }
        return err
    } else if err != nil {
        return err
    }

    if res.StatusCode < 200 || res.StatusCode >= 300 {
        errFmt := "Unsuccess HTTP request. Status: %s"
        return fmt.Errorf(errFmt, res.Status)
    }
    if res.Header.Get("Accept-Ranges") != "bytes" {     // server Range DƏSTƏKLƏMİR
        retries = 0                                    // → resume mümkün deyil
    }
    _, err = io.Copy(file, res.Body)                   // APPEND davam edir
    if err != nil && hasTimedOut(err) {                // kopya zamanı timeout
        if retries > 0 {
            return download(location, file, retries-1)
        }
        return err
    } else if err != nil {
        return err
    }
    return nil
}
```

**Sub-kod izahı:**
- `fi.Size()` → "hansı baytdan davam" (index 0-based — uzunluq = növbəti index!)
- `Range: bytes=N-` → N-dən EOF-a qədər
- `Accept-Ranges: bytes` header-i — server dəstəyi yoxlaması
- Rekursiv retry — hər dəfə Range yenidən hesablanır (fayl böyüyür)
- Tuning: timeout-u orta fayl müddətinə görə seç; hash yoxlaması əlavə etmək olar

### 2. HTTP error əsasları
```go
http.Error(w, "An Error Occurred", http.StatusForbidden)   // text/plain + status 403

// Client tərəfi:
res, _ := http.Get("http://example.com")
fmt.Println(res.Status)       // "200 OK" / "404 Not Found" (mətn)
fmt.Println(res.StatusCode)   // 200 / 404 (int)

switch {
case 300 <= res.StatusCode && res.StatusCode < 400:
    fmt.Println("Redirect message")
case 400 <= res.StatusCode && res.StatusCode < 500:
    fmt.Println("Client error")
case 500 <= res.StatusCode && res.StatusCode < 600:
    fmt.Println("Server error")
}
```
`http.Error` məhdudiyyəti: content type hardcoded text/plain; `X-Content-Type-Options: nosniff` qoyulur.

### TECHNIQUE 51: Custom JSON error (server)
**Kitabdan kod nümunəsi:**
```go
type Error struct {
    HTTPCode int    `json:"-"`               // JSON-a DÜŞMÜR — yalnız status üçün
    Code     int    `json:"code,omitempty"`  // app-apecific kod
    Message  string `json:"message"`
}

func JSONError(w http.ResponseWriter, e Error) {
    data := struct {
        Err Error `json:"error"`             // {"error": {...}} wrapper
    }{e}
    b, err := json.Marshal(data)
    if err != nil {
        http.Error(w, "Internal Server Error", 500)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(e.HTTPCode)                 // 403 və s.
    fmt.Fprint(w, string(b))
}

func displayError(w http.ResponseWriter, r *http.Request) {
    e := Error{
        HTTPCode: http.StatusForbidden,
        Code:     123,
        Message:  "An Error Occurred",
    }
    JSONError(w, e)
}
```
Nəticə:
```json
{
  "error": {
    "code": 123,
    "message": "An Error Occurred"
  }
}
```
+ **`json:"-"` dərsi** — struct sahəsi JSON-a serialize olunmasın.

### TECHNIQUE 52: Custom error OXUNMASI (client)
**Kitabdan kod nümunəsi:**
```go
type Error struct {
    HTTPCode int    `json:"-"`
    Code     int    `json:"code,omitempty"`
    Message  string `json:"message"`
}

func (e Error) Error() string {                // error İNTERFEYSİNİ implement et!
    fs := "HTTP: %d, Code: %d, Message: %s"
    return fmt.Sprintf(fs, e.HTTPCode, e.Code, e.Message)
}

func get(u string) (*http.Response, error) {   // http.Get ƏVƏZİNƏ
    res, err := http.Get(u)
    if err != nil {
        return res, err
    }
    if res.StatusCode < 200 || res.StatusCode >= 300 {
        if res.Header.Get("Content-Type") != "application/json" {
            sm := "Unknown error. HTTP status: %s"
            return res, fmt.Errorf(sm, res.Status)      // JSON deyil → generic
        }
        b, _ := ioutil.ReadAll(res.Body)
        res.Body.Close()
        var data struct {
            Err Error `json:"error"`
        }
        err = json.Unmarshal(b, &data)                    // JSON → Error struct
        if err != nil {
            sm := "Unable to parse json: %s. HTTP status: %s"
            return res, fmt.Errorf(sm, err, res.Status)
        }
        data.Err.HTTPCode = res.StatusCode
        return res, data.Err                             // Error ARTIQ error-dur!
    }
    return res, nil
}
```
- Error struct **error interfeysi** implement edir → normal Go error axınında istifadə
- Client + server EYNİ paketi paylaşa bilər (Error strukturu)
- 401 → re-login; SDK-a specific code ilə dəqiq reaksiya

### TECHNIQUE 53: Arbitrary JSON (schema bilinmir)
**Həll:** `interface{}`-a unmarshal + type walk.

**JSON → Go tip çevrilmə cədvəli:**
| JSON | Go |
|---|---|
| Boolean | `bool` |
| Number | `float64` |
| Array | `[]interface{}` |
| Object | `map[string]interface{}` |
| null | `nil` |
| String | `string` |

**Kitabdan kod nümunəsi:**
```go
var f interface{}
err := json.Unmarshal(ks, &f)          // istənilən JSON → interface{}

m := f.(map[string]interface{})         // type assertion → object
fmt.Println(m["firstName"])             // "Jean"

func printJSON(v interface{}) {         // rekursiv walk
    switch vv := v.(type) {
    case string:
        fmt.Println("is string", vv)
    case float64:
        fmt.Println("is float64", vv)
    case []interface{}:                 // array → elementləri rekursiv
        fmt.Println("is an array:")
        for i, u := range vv {
            fmt.Print(i, " ")
            printJSON(u)
        }
    case map[string]interface{}:       // object → sahələri rekursiv
        fmt.Println("is an object:")
        for i, u := range vv {
            fmt.Print(i, " ")
            printJSON(u)
        }
    default:
        fmt.Println("Unknown type")
    }
}
```
- `f.firstName` → COMPILE xətası (interface{}-də metod yoxdur) — əvvəlcə assertion!
- Type switch `v.(type)` — dinamik JSON naviqasiyasının açarı

### 4. API versioning — 2 üsul
**Semantika:** major (v1, v2 — BREAKING) + point (v1.1 — əlavə).

### TECHNIQUE 54: Versiya URL-də
```
https://example.com/api/v1/todos
                   └──┬──┘ └─┬─┘
                API versiyası  resurs
```
```go
http.HandleFunc("/api/v1/test", displayTest)     // versiya PATH-də
```
- Google, OpenStack, Salesforce, Twitter, Facebook bu üsulda
- **Pro:** asan test (cURL/Postman), görünür
- **Con:** semantik deyil (URL = versiya, obyekt yox)
- Qeyd: eyni path-də çoxlu metod (POST/PUT/DELETE) üçün Ch 2 router texnikaları

### TECHNIQUE 55: Versiya content type-da (semantic)
**Custom vendor content type:**
```
application/vnd.mytodo.v1.json
application/vnd.mytodos.json; version=1.0
```

**Kitabdan kod nümunəsi:**
```go
func displayTest(w http.ResponseWriter, r *http.Request) {
    t := r.Header.Get("Accept")                  // client NƏ İSTƏYİR
    var err error
    var b []byte
    var ct string
    switch t {
    case "application/vnd.mytodos.json; version=2.0":
        data := testMessageV2{"Version 2"}
        b, err = json.Marshal(data)
        ct = "application/vnd.mytodos.json; version=2.0"
    case "application/vnd.mytodos.json; version=1.0":
        fallthrough                              // v1 = DEFAULT-un üstünə düş
    default:                                     // tanınmayan → v1
        data := testMessageV1{"Version 1"}
        b, err = json.Marshal(data)
        ct = "application/vnd.mytodos.json; version=1.0"
    }
    if err != nil {
        http.Error(w, "Internal Server Error", 500)
        return
    }
    w.Header().Set("Content-Type", ct)
    fmt.Fprint(w, string(b))
}

type testMessageV1 struct {
    Message string `json:"message"`
}
type testMessageV2 struct {
    Info string `json:"info"`                     // FƏRLİ SCHEMA!
}
```

**Client tərəfi:**
```go
ct := "application/vnd.mytodos.json; version=2.0"
req, _ := http.NewRequest("GET", "http://localhost:8080/test", nil)
req.Header.Set("Accept", ct)                      // versiya SORĞUDA
res, _ := http.DefaultClient.Do(req)
if res.Header.Get("Content-Type") != ct {         // gözlənilən versiya gəldimi?
    fmt.Println("Unexpected content type returned")
    return
}
b, _ := ioutil.ReadAll(res.Body)
```

**Nüanslar:**
- `vnd.` namespace → IANA qeydiyyatı tövsiyə olunur
- Nondefault versiya → əlavə addımlar (sadə GET kifayət etmir)
- EYNİ URL həm HTML, həm JSON, həm API versiyaları verir — həqiqi semantic

## Versiya müqayisəsi
| | URL-də (T54) | Content type-da (T55) |
|---|---|---|
| Nümunə | /api/v1/todos | Accept: ...v2.0 |
| Developer asanlığı | ✓ (cURL/browser) | əlavə header |
| Semantik | ✗ | ✓ (URL = resurs) |
| Default davranış | açıq | tanınmayan → default versiya |
| İstifadəçilər | Google, Twitter... | — |

## Əsas terminlər
- http.Get/Head/Post/PostForm (helper)
- http.NewRequest / DefaultClient.Do
- http.Client{Timeout} (body oxuma daxil!)
- net.Error.Timeout() / *url.Error / *net.OpError
- Type Switch (err.(type))
- Range Header (`bytes=N-`) / Accept-Ranges: bytes
- Resumable/Incremental Download
- Recursive Retry
- http.Error / X-Content-Type-Options: nosniff
- Status sinifləri: 2xx/3xx/4xx/5xx
- JSON Error Envelope ({"error":{...}})
- `json:"-"` (serialize etmə)
- error interface implement (Error() string)
- Arbitrary JSON / interface{} unmarshal
- JSON→Go type map (bool/float64/[]interface{}/map[string]interface{}/nil/string)
- API Versioning: major/point
- vnd. Content Type / IANA
- Accept Header negotiation
- fallthrough (default-a düşmə)

## Praktik nəticə
- Timeout-u Client{Timeout} ilə AÇIQ qoy — sonra hasTimedOut ilə bütün halları yoxla (string fallback də daxil).
- Böyük fayl download → Range + retry: fayl ölçüsü =növbəti bayt; Accept-Ranges yoxlaması unutma.
- API error-larını da JSON ver — {"error":{code,message}} envelope; HTTPCode `json:"-"` ilə gizlət.
- Client-da Error struct-a Error() metodunu əlavə et — normal error axınına inteqrasiya.
- Schema naməlum JSON → interface{} + type switch; float64/number, map/object qaydalarını yadda saxla.
- API versiyası: developer rahatlığı → URL-də; semantik təmizlik → Accept content type; default versiya = tanınmayan sorğular üçün təhlükəsiz qayıdış.

## Mənbə
Pages: 217-236 (PDF), book pages 194-213
