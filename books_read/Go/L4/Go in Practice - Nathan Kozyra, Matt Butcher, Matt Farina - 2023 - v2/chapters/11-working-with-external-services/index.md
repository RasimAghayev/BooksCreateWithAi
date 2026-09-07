# Chapter 11 — Working with external services (Xarici servismlə iş)

## Bu chapter nədən bəhs edir?

HTTP client (custom client, timeout detect, Range ilə davamlı download), xətaların
HTTP üzərindən ötürülməsi (http.Error, JSON error), status kodu emalı, arbitrary JSON
(interface{} + rekursiv gəzinti), REST API versiyalama (URL və content-type) və gRPC
(proto3, codegen, Go client).

## Əsas fikirlər

### 1. HTTP Client
```go
cc := &http.Client{Timeout: time.Second}    // custom client — timeout MÜTLƏQ
res, err := cc.Get("http://www.manning.com")
b, _ := io.ReadAll(res.Body)
defer res.Body.Close()

// Hər hansı metod — DefaultClient + NewRequest:
req, _ := http.NewRequest("DELETE", "http://example.com/foo/bar", nil)
res, _ := http.DefaultClient.Do(req)
```
Custom client: transport, cookie, redirect davranışı konfiqurasiyası.

### 2. Timeout Aşkarlama
**Problem:** net-in error-lərində Timeout() həmişə true vermir; url paketindən gələnlər
metod daşımır. Type switch ilə 3 hal + string fallback:
```go
func hasTimedOut(err error) bool {
    switch err := err.(type) {
    case *url.Error:
        if err, ok := err.Err.(net.Error); ok && err.Timeout() { return true }
    case net.Error:
        if err.Timeout() { return true }
    case *net.OpError:
        if err.Timeout() { return true }
    }
    return err != nil && strings.Contains(err.Error(), "use of closed network connection")
}
```
(errors.Is ilə də mümkün — Ch4 patterni.)

### 3. Davamlı Download (Range Header)
**Kitabdan kod nümunəsi (çoxnaməlli retry):**
```go
func download(location string, file *os.File, retries int64) error {
    req, _ := http.NewRequest("GET", location, nil)
    fi, _ := file.Stat()
    current := fi.Size()
    if current > 0 {
        req.Header.Set("Range", "bytes="+strconv.FormatInt(current, 10)+"-")  // davam nöqtəsi
    }
    cc := &http.Client{Timeout: 5 * time.Minute}
    res, err := cc.Do(req)
    if err != nil && hasTimedOut(err) {
        if retries > 0 { return download(location, file, retries-1) }   // rekursiv retry
        return err
    }
    if res.StatusCode < 200 || res.StatusCode > 300 {
        return fmt.Errorf("Unsuccess HTTP request. Status: %s", res.Status)
    }
    if res.Header.Get("Accept-Ranges") != "bytes" {
        retries = 0                       // server range dəstəkləmir — retry faydasız
    }
    io.Copy(file, res.Body)               // append (file append rejimində açılmalı)
    ...
}
```
HTTP 1.1-dən (1999) Range standartdır; hash mümkünsə integrity check əlavə et.

### 4. Xətaların Ötürülməsi
**Sadə:** `http.Error(w, "An Error Occurred", http.StatusForbidden)` — text/plain +
X-Content-Type-Options: nosniff.

**JSON error (API-lər üçün):**
```go
type Error struct {
    HTTPCode int    `json:"-"`
    Code     int    `json:"code,omitempty"`
    Message  string `json:"message"`
}
func JSONError(w http.ResponseWriter, e Error) {
    data := struct{ Err Error `json:"error"` }{e}
    b, _ := json.Marshal(data)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(e.HTTPCode)
    fmt.Fprint(w, string(b))
}
// {"error": {"code": 123, "message": "An Error Occurred"}}
```
SDK-ya app-specific code + retry göstərişi kimi metadata verir.

**Client tərəfdə oxuma:** status 2xx-dən kənardsa + Content-Type application/json →
body-ni Error struct-a unmarshal, error interfeysini implement et (`func (e Error) Error() string`)
→ standart Go error axınında istifadə. 429/500 üçün incremental backoff retry strategiyası.
Status klassları: 30x redirect (client 10-20 avtomatik izləyir), 40x client, 50x server.

### 5. Arbitrary JSON (struct bəlli deyilsə)
```go
var f interface{}
json.Unmarshal(ks, &f)
m := f.(map[string]interface{})     // top-level obyekt
fmt.Println(m["firstName"])
```
**JSON→Go tip xəritəsi:** bool / float64 / []interface{} / map[string]interface{} / nil /
string.

**Rekursiv gəzinti:**
```go
func printJSON(v interface{}) {
    switch vv := v.(type) {
    case string:  fmt.Println("is string", vv)
    case float64: fmt.Println("is float64", vv)
    case []interface{}:
        for i, u := range vv { printJSON(u) }
    case map[string]interface{}:
        for i, u := range vv { printJSON(u) }
    }
}
```

### 6. REST API Versiyalama
**URL-də:** `/api/v1/todos` — ən populyar (Google, Salesforce, Facebook); cURL/Postman
ilə asan test; amma semver DEYİL (URL obyekt deyil, obyektə çıxışı təmsil edir):
```go
http.HandleFunc("/api/v1/test", displayTest)
```

**Content-type-da (semantik):** `Accept: application/vnd.mytodos.json; version=2.0` —
tək URL, müxtəlif versiyalar:
```go
switch t := r.Header.Get("Accept"); {
case t == "application/vnd.mytodos.json; version=2.0":
    // v2 struktur + cavab Content-Type: ...version=2.0
default:
    // v1 (default + fallthrough)
}
```
Client: `req.Header.Set("Accept", ct)` + cavab Content-Type-ni doğrula. Qeydlər: vnd.
namespace IANA qeydiyyatlı; nondefault versiya üçün client artıq addım atır. Feature
əlavəsi → point versiya (v1.1); breaking → major.

### 7. gRPC
**Nədir:** Google-un 2015 open-source RPC framework-ü; HTTP/2 üzərində protobuf (binary,
tipəmən kontrakt). RPC — köhnə distributiv hesablama konsepti.

**Üstünlüklər:** JSON-dən sürətli serializasiya; strong typing (RawMessage/interface{}
reflection ehtiyacsız); codegen boilerplate-i aradan qaldırır; HTTP/2 sayəsində
persistent + streaming (bir/birə çox/çoxa çox — `stream` keyword).
**Çatışmazlıqlar:** binary — insan-oxunmaz (transcoding istisna); brauzer test çətin
(Postman dəstəyi); kontrakt geri-uyğunluq zəmanəti dəyişməyə imkan vermir; qurma mürəkkəb.
Kiçik layihə üçün overkill — amma istehlakçı kimi tapacaqsan.

**Proto3 faylı:**
```protobuf
syntax = "proto3";
package chat;
import "google/protobuf/timestamp.proto";
option go_package = "./ch11";

message CommentRequest {
    string username = 1;
    string text = 2;
    google.protobuf.Timestamp sent = 3;
}
message CommentResponse {
    int32 commentLength = 1;
    int32 previousCommentCount = 2;
}
service ChatService {
    rpc RouteComments(CommentRequest) returns (CommentResponse) {}
    // streaming: rpc GetComments() returns (stream CommentResponse) {}
}
```
```bash
protoc --go_out=. --go-grpc_out=. protos/chat.proto   # .go faylları (DEYİŞTİRMƏ!)
```

**Go client:**
```go
conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
client := NewChatServiceClient(conn)
meta, err := client.RouteComments(context.Background(),
    &CommentRequest{Username: "Nick", Text: "Hello World!",
        Sent: timestamppb.New(time.Now())})
defer conn.Close()
```
Server Python-da eyni .proto-dan generasiya (grpc.io qurma dəstəyi) — dillərarası
zəmanətli kontrakt. Funksiya çağırışı qədər tip təhlükəsizliyi.

## Əsas terminlələr
- net.Error / url.Error / OpError — timeout yoxlama tipləri
- Range / Accept-Ranges — fayl hissəsi sorğusu / dəstək başlığı
- http.Error — standart text error cavabı
- nosniff — MIME sniff-in qarşısını alan header
- Arbitrary JSON — interface{} + type switch ilə naməlum strukturlu JSON
- vnd. content type — IANA vendor namespace versiyalı media type
- proto3 / protoc — protobuf dili / kompilyatoru
- grpc.Dial / WithInsecure — client bağlantı (test üçün insecure)
- stream (rpc) — gRPC axın tipləri

## Praktik nətidə

Client qərarları: (1) custom http.Client + timeout (DefaultClient production-da YOX);
(2) timeout — type switch hasTimedOut + retry; (3) böyük fayl — Range + fayl ölçüsü +
retry-lə davam; (4) API error-ları — JSON strukturlu (code+message) hər iki tərəfdən
eyni Error tipi ilə; (5) naməlum JSON — interface{} + rekursiv type switch; (6) versiya —
URL-də (asanlıq) və ya content-type-da (semantik) — JSON struktur dəyişəndə mütləq;
(7) cross-language, sürət, streaming lazımdırsa — gRPC (proto3 + protoc + codegen);
(8) hər status klassına (3xx/4xx/5xx) uyğun davranış + backoff retry.

## Mənbə
Pages: 266-290 (PDF 287-311)
