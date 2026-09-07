# Chapter 9 — Building Web Services (Web Servislərin Qurulması)

## Bu chapter nədən bəhs edir?

net/http paketinin 3 əsas strukturu (Response, Request, Transport), web server quruluşu
(handler-lər, DefaultServeMux vs öz ServeMux, http.Server timeout-ları), statistika
tətbiqinin web servisinə çevrilməsi (API dizaynı, handler implementasiyası), Docker
image qurulması, HTTP client-lər (http.Get, http.NewRequest, errgroup), cobra ilə CLI
client və HTTP timeout texnikaları (SetDeadline, context, server-side).

## Əsas fikirlər

### 1. net/http-nin 3 Strukturu
- **http.Response** — serverin cavabı: Status, StatusCode, Body (io.ReadCloser),
  ContentLength, Header
- **http.Request** — sorğu: Method, URL, Header, Body, Host, RemoteAddr
- **http.Transport** — aşağı səviyyəli bağlantı nəzarəti: Dial/TLS funksiyaları, proxy,
  keep-alive (DisableKeepAlives, MaxIdleConns), müxtəlif timeout-lar (TLSHandshakeTimeout,
  IdleConnTimeout, ResponseHeaderTimeout). http.Client-on Transport sahəsi var — nil
  olsa DefaultTransport.

### 2. Web Server — Əsas Pattern
**Kitabdan kod nümunəsi:**
```go
func myHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Serving: %s\n", r.URL.Path)   // w = io.Writer → client-ə
    fmt.Printf("Served: %s\n", r.Host)            // terminala
}

PORT := ":8001"           // ":" vacibdir — hostname-siz = bütün interface-lər
http.HandleFunc("/time", timeHandler)
http.HandleFunc("/", myHandler)     // "/" = default handler — digərlərinin tutmadığı hər şey
err := http.ListenAndServe(PORT, nil)   // nil = DefaultServeMux
```
**Sub-kod izahı:**
- `http.HandleFunc` DefaultServeMux-a qeydiyyat; `ListenAndServe`-in 2-ci parametri
  nil = default router
- Port 0 → RANDOM boş port (test üçün)
- HTTP (HTTPS YOX) — mikroservis praktikası: Go server Docker arxasında, TLS-i
  Nginx/Caddy ön tərəf həll edir

### 3. Statistika Web Servisi — API Dizaynı
**Endpointlər:** `/list`, `/insert/name/d1/d2/.../`, `/delete/name/`, `/search/name/`,
`/status` (entry sayı). Data mübadilə formatı: plain text (JSON Ch11-də, HTML parse
tələb etdiyindən imtina).

**Default router YOX — öz ServeMux:**
```go
mux := http.NewServeMux()
s := &http.Server{
    Addr:         PORT,
    Handler:      mux,
    IdleTimeout:  10 * time.Second,
    ReadTimeout:  time.Second,       // sorğunun oxunma həddi
    WriteTimeout: time.Second,       // cavabın yazma həddi
}

// mux.Handle — http.HandlerFunc adapteri MÜTLƏQ (default-da avtomatik idi):
mux.Handle("/list", http.HandlerFunc(listHandler))
mux.Handle("/insert/", http.HandlerFunc(insertHandler))
mux.Handle("/delete/", http.HandlerFunc(deleteHandler))
mux.Handle("/", http.HandlerFunc(defaultHandler))

err = s.ListenAndServe()
```
`mux.HandleFunc("/list", listHandler)` — eyni şeyin qısa yolu (implicit conversion).

### 4. Handler Implementasiyaları
**Kitabdan kod nümunələri:**

DELETE — URL split + status kodları:
```go
func deleteHandler(w http.ResponseWriter, r *http.Request) {
    paramStr := strings.Split(r.URL.Path, "/")
    if len(paramStr) < 3 {
        w.WriteHeader(http.StatusNotFound)     // cavabdan ƏVVƏL status kodu!
        fmt.Fprintln(w, "Not found:", r.URL.Path)
        return
    }
    dataset := paramStr[2]
    err := deleteEntry(dataset)
    if err != nil {
        w.WriteHeader(http.StatusNotFound)
        fmt.Fprintf(w, "%s", err.Error()+"\n")
        return
    }
    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "%s deleted!\n", dataset)
}
```

INSERT — dəyişən uzunluqlu data:
```go
paramStr := strings.Split(r.URL.Path, "/")
dataset := paramStr[2]
dataStr := paramStr[3:]                    // URL-dəki qalan hissələr
data := make([]float64, 0)
for _, v := range dataStr {
    if val, err := strconv.ParseFloat(v, 64); err == nil {
        data = append(data, val)
    }
}
entry := process(dataset, data)
if err := insert(&entry); err != nil {
    w.WriteHeader(http.StatusNotModified)   // 304 — dəyişiklik OLMADI
} else {
    w.WriteHeader(http.StatusOK)
}
```

**Vacib detal:** `/delete` (slash-siz) → Go router `301 Moved Permanently` ilə
`/delete/`-a yönləndirir — hər iki variantı qeyd etmək bunu aradan qaldırır.
`WriteHeader` bodidən ƏVVƏL çağırılmalı. Server düzgün status kodu göndərməlidir;
yoxlama isə client-in vəzifəsidir.

### 5. Docker Image — Multi-Stage Build
**Kitabdan kod nümunəsi (buildDocker faylı):**
```dockerfile
FROM golang:alpine AS builder
RUN apk update && apk add --no-cache git     # asılılıqlar üçün
RUN mkdir $GOPATH/src/server
ADD ./stats.go $GOPATH/src/server
ADD ./handlers.go $GOPATH/src/server
WORKDIR $GOPATH/src/server
RUN go mod init && go mod tidy
RUN mkdir /pro
RUN go build -o /pro/server stats.go handlers.go

FROM alpine:latest                           # İKİNCİ mərhələ — yalnız binary
RUN mkdir /pro
COPY --from=builder /pro/server /pro/server  # builder-dən yalnız nəticə
EXPOSE 1234
WORKDIR /pro
CMD ["/pro/server"]
```
```bash
docker build -f buildDocker -t goapp .
# docker-compose.yml: image goapp, ports 1234:1234, restart always
docker-compose up
```
**Multi-stage build dərsi:** builder mərhələsində Go toolchain + mənbə; son mərhələ
yalnız minimal alpine + binary → kiçik image. Diqqət: konteyner dayansa data.json
İTİR — xarici DB və ya volume mount lazımdır.

### 6. HTTP Client — Sadə və Nəzarətli
**Sadə:**
```go
data, err := http.Get(URL)
if err != nil { return }
io.Copy(os.Stdout, data.Body)   // body → stdout stream
data.Body.Close()                // GC üçün MÜTLƏQ
```

**Nəzarətli (kitabdan):**
```go
URL, _ := url.Parse(os.Args[1])           // URL validasiyası
c := &http.Client{
    Timeout: 15 * time.Second,             // client-övəzi timeout
}
request, err := http.NewRequest(http.MethodGet, URL.String(), nil)
httpData, err := c.Do(request)            // sorğunu GÖNDƏR

fmt.Println("Status code:", httpData.Status)
header, _ := httputil.DumpResponse(httpData, false)   // debug üçün header dump
contentType := httpData.Header.Get("Content-Type")
// charset = strings.SplitAfter(contentType, "charset=")

// Manual ölçü hesabı:
var buffer [1024]byte
r := httpData.Body
length := 0
for {
    n, err := r.Read(buffer[0:])
    if err != nil { break }               // io.EOF da burada
    length += n
}
```
`http.NewRequest` + `c.Do` = http.Get-in detallı forması — header, method, body
nəzarəti verir.

### 7. errgroup — Paralel Sorğular
**Kitabdan kod nümunəsi:**
```go
import "golang.org/x/sync/errgroup"

g := new(errgroup.Group)
for _, url := range os.Args[1:] {
    url := url                       // closure nüsxəsi (Go 1.21 idiomu)
    g.Go(func() error {
        resp, err := http.Get(url)
        if err != nil { return err }
        defer resp.Body.Close()
        fmt.Println(url, "is OK.")
        return nil
    })
}
if err := g.Wait(); err != nil {     // İLK xəta qayıdır
    fmt.Println("Error:", err)
}
```
errgroup = sinxronizasiya + xəta ötürülməsi + context ləğvi — bir paketdə.

### 8. CLI Client — cobra ilə
**Qlobal parametrlər (root.go init):**
```go
rootCmd.PersistentFlags().StringP("server", "S", "localhost", "Server")
rootCmd.PersistentFlags().StringP("port", "P", "1234", "Port number")
viper.BindPFlag("server", ...)
```

**Komanda patterni (status/list/delete/search/insert):**
```go
SERVER := viper.GetString("server")
PORT := viper.GetString("port")
URL := "http://" + SERVER + ":" + PORT + "/status"

data, err := http.Get(URL)
if err != nil { return }
if data.StatusCode != http.StatusOK {     // STATUS KODU YOXLA — həmişə
    fmt.Println("Status code:", data.StatusCode)
    return
}
responseData, err := io.ReadAll(data.Body)
fmt.Print(string(responseData))
```

**INSERT komandasının client tərəfi:**
```go
insertCmd.Flags().StringP("dataset", "d", "", "Dataset name")
insertCmd.Flags().StringP("values", "v", "", "List of values")

VALS := strings.Split(values, ",")
vSend := ""
for _, v := range VALS {
    if _, err := strconv.ParseFloat(v, 64); err == nil {
        vSend = vSend + "/" + v           // client tərəfli validasiya!
    }
}
URL := "http://" + SERVER + ":" + PORT + "/insert/" + "/" + dataset + "/" + vSend + "/"
```

### 9. HTTP Timeout — 3 Texnika
**1) Transport Dial + SetDeadline (aşağı səviyyə):**
```go
func Timeout(network, host string) (net.Conn, error) {
    conn, err := net.DialTimeout(network, host, timeout)
    if err != nil { return nil, err }
    conn.SetDeadline(time.Now().Add(timeout))   // hər oxu/yazıdan ƏVVƏL
    return conn, nil
}
t := http.Transport{ Dial: Timeout }
client := http.Client{ Transport: &t }
```

**2) Client-side context (kitabın tövsiyə üsulu):**
```go
ctx, cncl := context.WithTimeout(context.Background(), time.Second*time.Duration(delay))
defer cncl()
req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
res, err := http.DefaultClient.Do(req.WithContext(ctx))
defer res.Body.Close()
// vaxt bitərsə: "context deadline exceeded"
```

**3) Server-side timeout (ƏN VACİB — DoS qorunması):**
```go
srv := &http.Server{
    ReadTimeout:  3 * time.Second,    // sorğunun tam oxunması həddi
    WriteTimeout: 3 * time.Second,    // cavabın göndərilmə həddi
}
```
Server çox açıq bağlantı saxlayırsa yeni sorğulara xidmət edə bilmir — səbəblər:
bug-lar və DoS hücumu. Timeout bağlantını təmizləyir.

## Əsas terminlələr
- http.ResponseWriter — client-ə yazma interfeysi (io.Writer)
- Handler — `func(http.ResponseWriter, *http.Request)` imzalı funksiya
- ServeMux / DefaultServeMux — HTTP router; HandleFunc avtomatik konversiya
- http.HandlerFunc — funksiyanı Handler-ə çevirən adapter tipi
- http.Server — Addr/Handler/Timeout-larla server konfiqurasiyası
- ReadTimeout/WriteTimeout/IdleTimeout — server tərəfli hədlər
- http.Transport — aşağı səviyyə bağlantı qatı (dial, keep-alive, proxy)
- Multi-Stage Build — kiçik image üçün ikimərhələli Dockerfile
- errgroup — paralel işlərin xəta ötürən qrupu
- SetDeadline — connection səviyyəli i/o həddi
- NewRequestWithContext — context-li sorğu (deadline/ləğv)
- Status Codes (status kodları) — 200 OK, 301, 304 Not Modified, 400, 404

## Praktik nətidə

(1) Go web serverinə mux ver: DefaultServeMux qlobaldır — paylaşılan kitabxana
kodunda zəhərli; öz NewServeMux + HandlerFunc adapteri. (2) "/" handler = catch-all
default; "/"-lu və "/"-siz route-ları hər ikisini qeyd et — 301 redirect qarşısı.
(3) WriteHeader body-dən ƏVVƏL; status kodu sözleşmədir — server göndərir, client
yoxlayır. (4) net/http hər sorğunu ayrıca goroutine-də işlədir — server avtomatik
concurrent-dir; amma paylaşılan state (data slice) racesiz deyil — mutex lazımdır.
(5) Body həmişə bağlanır (defer Close). (6) Multi-stage Docker build — production
image kiçikliyi; data davamlılığı üçün volume/DB. (7) Client-də 3 timeout qatı:
Client.Timeout (ümumi), Transport (dial/header), context (ən çevik — ləğv səbəbi ilə).
(8) Server ReadTimeout/WriteTimeout — DoS-ə qarşı müdafiə xətti. (9) errgroup paralel
sorğularda xəta idarəsini sadələşdirir.

## Mənbə
Pages: 387-424 (PDF 418-457)
