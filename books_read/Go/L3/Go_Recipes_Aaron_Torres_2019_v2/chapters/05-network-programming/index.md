# Chapter 5 — Сетевое программирование (Şəbəkə proqramlaşdırması)

## Bu chapter nədən bəhs edir?

TCP/IP echo server, UDP broadcast, DNS lookup (net paketi), WebSocket
(gorilla/websocket), net/rpc ilə uzaq metod çağırışı və net/mail ilə e-poçt
parse.

## Əsas fikirlər

### 1. TCP/IP echo server və client
**Nədir:** net paketi ilə TCP bağlantı — HTTP-nin altında yatan protokol.

**Necə işləyir:** Server `net.Listen` + `Accept` tsikli; hər bağlantı ayrı
goroutine-da emal olunur; client `net.Dial` ilə qoşulur.

**Kitabdan kod nümunəsi:**
```go
// SERVER:
const addr = "localhost:8888"
func echoBackCapitalized(conn net.Conn) {
    reader := bufio.NewReader(conn)          // conn = io.Reader!
    data, err := reader.ReadString('\n')     // ilk sətiri oxu
    if err != nil {
        fmt.Printf("error reading data: %s\n", err.Error())
        return
    }
    conn.Write([]byte(strings.ToUpper(data)))  // böyük hərflə geri yaz
    conn.Close()
}
func main() {
    ln, err := net.Listen("tcp", addr)
    if err != nil { panic(err) }
    defer ln.Close()
    for {
        conn, err := ln.Accept()      // bağlantı gözlə
        if err != nil { continue }    // xətada davam et
        go echoBackCapitalized(conn)  // hər client ayrı goroutine
    }
}

// CLIENT:
conn, err := net.Dial("tcp", addr)    // qoşul
fmt.Fprintf(conn, data)               // yaz
status, err := bufio.NewReader(conn).ReadString('\n')  // cavab oxu
conn.Close()
```

**Sub-kod izahı:**
- `net.Conn` → həm Reader həm Writer — Chapter 1 interfeysləri burada da işləyir
- `go echoBackCapitalized(conn)` → server accept tsiklini bloklamır
- `Accept` xətasında continue → server ölmür

### 2. UDP server və client (broadcast)
**Nədir:** Connectionless protokol — sürət etibarlılıqdan önəmlidir (oyunlar
üçün). Connect YOXDUR — ünvanlara WriteToUDP.

**Kitabdan kod nümunəsi:**
```go
// SERVER — client ünvanlarını toplayıb hamısına broadcast:
type connections struct {
    addrs map[string]*net.UDPAddr
    mu    sync.Mutex                // map paralel dəyişilir — mutex ŞƏRT
}
func broadcast(conn *net.UDPConn, conns *connections) {
    count := 0
    for {
        count++
        conns.mu.Lock()
        for _, retAddr := range conns.addrs {
            msg := fmt.Sprintf("Sent %d", count)
            conn.WriteToUDP([]byte(msg), retAddr)   // ünvana yaz
        }
        conns.mu.Unlock()
        time.Sleep(1 * time.Second)
    }
}
func main() {
    addr, _ := net.ResolveUDPAddr("udp", addr)
    conn, err := net.ListenUDP("udp", addr)
    defer conn.Close()
    go broadcast(conn, conns)
    msg := make([]byte, 1024)
    for {
        _, retAddr, err := conn.ReadFromUDP(msg)   // gələn mesajın ÜNVANI
        if err != nil { continue }
        conns.mu.Lock()
        conns.addrs[retAddr.String()] = retAddr    // ünvanı qeyd et
        conns.mu.Unlock()
    }
}

// CLIENT:
conn, err := net.DialUDP("udp", nil, addr)
conn.Write([]byte("connected"))    // "qeydiyyat" mesajı
for {
    n, err := conn.Read(msg)        // broadcast-ləri dinlə
    fmt.Printf("%s\n", string(msg[:n]))
}
```

**Sub-kod izahı:**
- UDP-də client "qoşulmur" — ilk mesaj serverə ünvanı ötürür
- `WriteToUDP(data, addr)` → hədəf ünvanı açıq verilir
- Server client sayından asılı olmayaraq hər saniyə broadcast edir

### 3. DNS lookup
**Nədir:** Domain → CNAME + IP ünvanları; `dig` komandasının sadələşdirilmiş
analogsu.

**Kitabdan kod nümunəsi:**
```go
type Lookup struct {
    cname string
    hosts []string
}
func (d *Lookup) String() string {     // Stringer → çapda avtomatik
    result := ""
    for _, host := range d.hosts {
        result += fmt.Sprintf("%s IN A %s\n", d.cname, host)
    }
    return result
}
func LookupAddress(address string) (*Lookup, error) {
    cname, err := net.LookupCNAME(address)   // kanonik ad
    if err != nil {
        return nil, errors.Wrap(err, "error looking up CNAME")
    }
    hosts, err := net.LookupHost(address)    // IPv4 + IPv6 ünvanları
    if err != nil {
        return nil, errors.Wrap(err, "error looking up HOST")
    }
    return &Lookup{cname: cname, hosts: hosts}, nil
}
// GODEBUG=netdns=go və ya cgo → resolver seçimi
```

**Sub-kod izahı:**
- `net.LookupCNAME` → "golang.org." kimi kanonik ad
- `net.LookupHost` → bütün A/AAAA qeydləri
- Default: saf Go DNS resolver (cgo alternativdir)

### 4. WebSocket (gorilla/websocket)
**Nədir:** HTTP üzərində iki tərəfli (bidirectional) bağlantı — chat, real-time.

**Necə işləyir:** HTTP handler `upgrader.Upgrade` ilə WS bağlantısına
keçir; sonra mesaj tsikli oxu/yaz.

**Kitabdan kod nümunəsi:**
```go
// SERVER — HTTP handler-dən WS upgrade:
var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}
func wsHandler(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)   // HTTP → WebSocket
    if err != nil {
        log.Println("failed to upgrade connection: ", err)
        return
    }
    for {
        messageType, p, err := conn.ReadMessage()
        if err != nil { return }
        conn.WriteMessage(messageType, p)      // echo back
    }
}
log.Panic(http.ListenAndServe("localhost:8000", http.HandlerFunc(wsHandler)))

// CLIENT:
u := "ws://localhost:8000/"                    // ws:// scheme
c, _, err := websocket.DefaultDialer.Dial(u, nil)
defer c.Close()
// Graceful close — siqnal gələndə:
func catchSig(ch chan os.Signal, c *websocket.Conn) {
    <-ch
    c.WriteMessage(websocket.CloseMessage,
        websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
}
interrupt := make(chan os.Signal, 1)     // buffered kanal — 1 yer!
signal.Notify(interrupt, os.Interrupt)
go catchSig(interrupt, c)
// Mesaj tsikli:
err = c.WriteMessage(websocket.TextMessage, []byte(data))
_, message, err := c.ReadMessage()
```

**Sub-kod izahı:**
- `upgrader.Upgrade(w, r, nil)` → HTTP sorğusunu WS-ə "yüksəldir"
- `CloseMessage` göndərmədən bağlamaq = qeyri-düzgün bağlanma
- `signal.Notify`-ə buffered kanal verilməlidir (siqnal itmir)
- Handler async olduğu üçün min bir bağlantı paralel işləyir

### 5. net/rpc ilə uzaq metod çağırışı
**Nədir:** Go-nun sadə daxili RPC paketi (gRPC-dən yüngül, amma məhdud).

**Məhdudiyyətlər:** Metod imzası ciddi formada: 2 arqument (ixrac olunan),
ikincisi pointer, return error; tip və metod ixrac olunmalı.

**Kitabdan kod nümunəsi:**
```go
// PAYLAŞILAN TİP (server və client istifadə edir):
type StringTweaker struct{}
type Args struct {
    String  string
    ToUpper bool
    Reverse bool
}
// RPC qaydalarına uyğun imza:
func (s StringTweaker) Tweak(args *Args, resp *string) error {
    result := args.String
    if args.ToUpper { result = strings.ToUpper(result) }
    if args.Reverse {
        runes := []rune(result)         // rune slice — unicode-safe reverse
        for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
            runes[i], runes[j] = runes[j], runes[i]
        }
        result = string(runes)
    }
    *resp = result                       // cavab pointer-ə yazılır
    return nil
}

// SERVER:
s := new(tweak.StringTweaker)
rpc.Register(s)                 // metodu qeyd et
rpc.HandleHTTP()                // HTTP üzərindən serve et
l, _ := net.Listen("tcp", ":1234")
http.Serve(l, nil)

// CLIENT:
client, err := rpc.DialHTTP("tcp", "localhost:1234")
args := tweak.Args{String: "test", ToUpper: true, Reverse: true}
var result string
err = client.Call("StringTweaker.Tweak", args, &result)   // "Tip.Metod" adı
```

**Sub-kod izahı:**
- `client.Call("StringTweaker.Tweak", args, &result)` → string ilə metodu çağır
- Rune slice ilə reverse → bayt deyil, simvol səviyyəsində (UTF-8 düzgün)

### 6. net/mail ilə e-poçt parse
**Nədir:** Xam e-poçt mətnindən header/body ayrılması.

**Kitabdan kod nümunəsi:**
```go
r := strings.NewReader(msg)          // xam e-poçt stringi
m, err := mail.ReadMessage(r)        // m.Header + m.Body

func printHeaderInfo(header mail.Header) {
    toAddress, err := mail.ParseAddress(header.Get("To"))
    if err == nil {
        fmt.Printf("To: %s <%s>\n", toAddress.Name, toAddress.Address)
    }
    fromAddress, err := mail.ParseAddress(header.Get("From"))
    if err == nil {
        fmt.Printf("From: %s <%s>\n", fromAddress.Name, fromAddress.Address)
    }
    fmt.Println("Subject:", header.Get("Subject"))
    if date, err := header.Date(); err == nil {     // RFC5322 tarix parse
        fmt.Println("Date:", date)
    }
}
// Body — sadə io.Reader-dir:
io.Copy(os.Stdout, m.Body)
```

**Sub-kod izahı:**
- `mail.ReadMessage(r)` → Reader qəbul edir (string/file/şəbəkə — istənilən)
- `m.Header.Get("X")` → istənilən header; `.Date()` xüsusi RFC5322 parse edir
- `mail.ParseAddressList` → çoxlu ünvan üçün (tək ünvanda ParseAddress)
- `m.Body` io.Reader → böyük məktublar da stream ilə oxunur

## Əsas terminlər

- TCP/IP (etibarlı bağlantı protokolu)
- UDP (bağlantısız protokol) / Broadcast (yayım)
- net.Listen / net.Dial / net.Accept
- net.Conn (şəbəkə bağlantısı — Reader/Writer)
- DNS / CNAME / A Record ( domen adı sistemi)
- WebSocket (iki tərəfli veb bağlantısı) / Upgrade
- RPC (Remote Procedure Call — uzaq prosedur çağırışı)
- net/rpc imza qaydaları
- net/mail / RFC5322 (e-poçt standartı)
- Stringer interfeysi (String() metodu)

## Praktik nəticə

- net.Conn həm Reader həm Writer — bütün I/O pattern-ləri şəbəkədə də işləyir
- TCP server: Accept tsikli + goroutine per connection (worker pool ilə
  təkmilləşdirilə bilər)
- UDP-də ünvanlar map-də toplanır — mutex ilə qorunmalı
- WS close mesajı göndərmək qonaqqartasiya normal bağlıqdır
- net/rpc sadə hallarda gRPC-yə qədər kifayətdir; imza qaydalarına diqqət
- net/mail xam e-poçtları DevOps script-lərində parse etmək üçün idealdır

## Mənbə

Pages: 163-190 (Chapter 5, Go Programming Cookbook 2nd ed)
