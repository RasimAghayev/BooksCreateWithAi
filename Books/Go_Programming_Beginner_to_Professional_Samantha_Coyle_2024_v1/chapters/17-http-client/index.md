# Chapter 17 — Using the Go HTTP Client (səh. 540-561)

## Bu fəsil nədən bəhs edir?

Go HTTP klienti: GET (http.Get + io.ReadAll), JSON parse
(json.Unmarshal), POST (http.Post + json.Marshal), fayl upload (multipart
form), custom client (http.Client{Timeout}) + authorization header və
NewRequest ilə tam nəzarət.

## Əsas fikirlər

### 1. GET request
**URL quruluşu:** `https://example.com/downloads?filter=latest&os=windows`
- Protocol (https) + Hostname + URI (/downloads) + Query parametrləri
  (? ilə ayrılır, & ilə birləşir)

**Sadə GET (Exercise 17.01):**
```go
func getDataAndReturnResponse() string {
    r, err := http.Get("https://www.google.com")   // DEFAULT client
    if err != nil {
        log.Fatal(err)
    }
    defer r.Body.Close()                            // body bağla!
    data, err := io.ReadAll(r.Body)                  // hamısını oxu
    if err != nil {
        log.Fatal(err)
    }
    return string(data)
}
```
- Qaytarılan data — Google home page HTML-i (browser-in arxasındakı
  mexanizmin eynisi!)

### 2. Struktur data (JSON) + Unmarshal
**Server (JSON qaytaran):**
```go
type server struct{}

func (srv server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    msg := "{\"message\": \"hello world\"}"
    w.Write([]byte(msg))
}

func main() {
    log.Fatal(http.ListenAndServe(":8080", server{}))
}
```

**Klient:**
```go
// 1) Struct + JSON tag (EXPORTED sahə şərt!):
type messageData struct {
    Message string `json:"message"`
}

func getDataAndReturnResponse() messageData {
    // 2) GET:
    r, err := http.Get("http://localhost:8080")
    if err != nil {
        log.Fatal(err)
    }
    defer r.Body.Close()
    data, err := io.ReadAll(r.Body)
    if err != nil {
        log.Fatal(err)
    }

    // 3) PARSE:
    message := messageData{}
    err = json.Unmarshal(data, &message)     // & — pointer!
    if err != nil {
        log.Fatal(err)
    }
    return message
}

func main() {
    data := getDataAndReturnResponse()
    fmt.Println(data.Message)                   // hello world
}
```
- JSON tag-lər best practice — hər zaman açıq yaz

### 3. POST request
**Server (JSON decode edib geri qaytaran):**
```go
func (srv server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    jsonDecoder := json.NewDecoder(r.Body)
    messageData := messageData{}
    err := jsonDecoder.Decode(&messageData)
    if err != nil {
        log.Fatal(err)
    }
    jsonBytes, _ := json.Marshal(messageData)
    w.Write(jsonBytes)                       // eynini geri göndər
}
```

**Klient:**
```go
func postDataAndReturnResponse(msg messageData) messageData {
    // 1) Struct → bytes:
    jsonBytes, _ := json.Marshal(msg)

    // 2) POST (url, contentType, body):
    r, err := http.Post("http://localhost:8080", "application/json",
        bytes.NewBuffer(jsonBytes))
    if err != nil {
        log.Fatal(err)
    }

    // 3) Cavabı oxu + parse (GET ilə eyni):
    defer r.Body.Close()
    data, err := io.ReadAll(r.Body)
    if err != nil {
        log.Fatal(err)
    }
    message := messageData{}
    err = json.Unmarshal(data, &message)
    if err != nil {
        log.Fatal(err)
    }
    return message
}

func main() {
    msg := messageData{Message: "Hi Server!"}
    data := postDataAndReturnResponse(msg)
    fmt.Println(data.Message)     // Hi Server!
}
```

### 4. Fayl upload (multipart form)
**Server (faylı qəbul edən):**
```go
func (srv server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    uploadedFile, uploadedFileHeader, err := r.FormFile("myFile")
    if err != nil {
        log.Fatal(err)
    }
    defer uploadedFile.Close()
    fileContent, err := io.ReadAll(uploadedFile)
    if err != nil {
        log.Fatal(err)
    }
    // ... saxla (uploadedFileHeader.Filename ilə)
}
```

**Klient:**
```go
func postFileAndReturnResponse(filename string) string {
    // 1) Bufer + multipart writer:
    fileDataBuffer := bytes.Buffer{}
    multipartWriter := multipart.NewWriter(&fileDataBuffer)

    // 2) Lokal faylı aç:
    file, err := os.Open(filename)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    // 3) Form faylı yarat + faylı köçür:
    formFile, err := multipartWriter.CreateFormFile("myFile", file.Name())
    if err != nil {
        log.Fatal(err)
    }
    _, err = io.Copy(formFile, file)          // file → formFile
    if err != nil {
        log.Fatal(err)
    }
    multipartWriter.Close()                    // writer-i BİTİR!

    // 4) Request YARAT (http.Post yetmir — nəzarət lazımdır):
    req, err := http.NewRequest("POST", "http://localhost:8080", &fileDataBuffer)
    if err != nil {
        log.Fatal(err)
    }

    // 5) Content-Type — multipart boundary daxil!
    req.Header.Set("Content-Type", multipartWriter.FormDataContentType())

    // 6) Göndər:
    response, err := http.DefaultClient.Do(req)
    if err != nil {
        log.Fatal(err)
    }
    defer response.Body.Close()
    data, err := io.ReadAll(response.Body)
    if err != nil {
        log.Fatal(err)
    }
    return string(data)
}

postFileAndReturnResponse("./test.txt")
```
- **Multipart form** — fayl + mətn key-value birləşdirən encoding (HTML
  form upload standartı)
- **DefaultClient production-da YOX** — timeout YOXDUR; böyük fayllarda
  yaddaş daşması riski (os.Pipe alternativi)

### 5. Custom client + headers + timeout
**Server (auth yoxlayan + 10s gözləyən):**
```go
func (srv server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    auth := r.Header.Get("Authorization")
    if auth != "superSecretToken" {
        w.WriteHeader(http.StatusUnauthorized)       // 401
        w.Write([]byte("Authorization token not recognized"))
        return
    }
    time.Sleep(10 * time.Second)                       // yavaş cavab
    w.Write([]byte("hello client!"))
}
```

**Klient:**
```go
func getDataWithCustomOptionsAndReturnResponse() string {
    // 1) ÖZ client-i — TIMEOUT ilə:
    client := http.Client{Timeout: 11 * time.Second}

    // 2) Request yarat:
    req, err := http.NewRequest("POST", "http://localhost:8080", nil)
    if err != nil {
        log.Fatal(err)
    }

    // 3) AUTH header:
    req.Header.Set("Authorization", "superSecretToken")

    // 4) İcra:
    resp, err := client.Do(req)
    if err != nil {
        log.Fatal(err)
    }
    defer resp.Body.Close()
    data, err := io.ReadAll(resp.Body)
    if err != nil {
        log.Fatal(err)
    }
    return string(data)
}
```
- Timeout <10s etsən → client xəta ilə qayıdacaq (server hələ yatırdır)
- Token yoxdursa/səlidirsə → 401

## Activity icmalları
- **17.01:** GET + names massivi ({"names":["Electric","Boogaloo",...]})
  → struct-a parse → hər adın sayını hesabla
- **17.02:** POST {"name":"Electric"} → server {"ok":true} → GET ilə
  adları geri al → çap et (write + read dövrü)

## Əsas terminlər
- http.Get / http.Post — default klient qısayolları
- r.Body + defer Close + io.ReadAll — cavab oxu
- json.Unmarshal(&struct) / json.Marshal(struct) — parse/seriya
- JSON struct tag — `json:"message"`; EXPORTED sahələr
- bytes.NewBuffer — []byte → io.Reader
- Multipart form — fayl upload encoding-i
- multipart.NewWriter / CreateFormFile / FormDataContentType
- http.NewRequest — tam nəzarətli request
- http.DefaultClient — timeout-suz (production-da İŞLƏTMƏ!)
- http.Client{Timeout} — öz klientin
- req.Header.Set — Authorization kimi başlıqlar
- client.Do(req) — request-in göndərilməsi
- 401 StatusUnauthorized

## Praktik nəticə
GET: http.Get → Body → ReadAll → (JSON isə Unmarshal). POST:
json.Marshal → bytes.NewBuffer → http.Post(url, "application/json",
buffer). Fayl upload: multipart writer + CreateFormFile + io.Copy +
NewRequest + FormDataContentType. DefaultClient-dan production-da qaçın
— http.Client{Timeout} ilə özünü yarat; header-lər (Authorization)
NewRequest + Header.Set ilə. Strukturlarda həmişə JSON tag + exported
sahələr. Böyük fayllar üçün buffer əvəzinə stream (os.Pipe) nəzərə al.

## Mənbə
Pages: 540-561 (PDF 540-561)
