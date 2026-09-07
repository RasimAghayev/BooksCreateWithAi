# Chapter 10 — Working with TCP/IP and WebSocket (TCP/IP və WebSocket ilə İş)

## Bu chapter nədən bəhs edir?

TCP/IP protokolları (TCP, IP, UDP, IPv4/IPv6), nc(1) aləti, net paketi ilə TCP client/
server-lər (Dial vs DialTCP, Listen vs ListenTCP), UDP client/server-lər, concurrent TCP
server patterni, WebSocket protokolu (gorilla/websocket ilə server + client), websocat,
RabbitMQ message broker (AMQP, producer/consumer) və go.mod-dan module silmə.

## Əsas fikirlər

### 1. TCP/IP Protokolları
- **TCP (Transmission Control Protocol)** — ETİBARLI: hər paketin çatdırılmasını təsdiqlər,
  itirsə yenidən göndərir; full-duplex virtual dövrə (telefon zəngi kimi). Source/destination
  port + IP = unikal bağlantı.
- **IP (Internet Protocol)** — ÇATDIRICI, amma özü ETİBARLI DEYİL; paketləri IP ünvanına
  görə yönləndirir (routing).
- **UDP** — IP üzərində, sadə və ETİBARSIZ: mesaj itə/dublikat/sırasız gələ bilər; sürət
  vacib olanda istifadə olunur.
- **IPv4** — 32 bit (2³² ünvan, tükənir); **IPv6** — 128 bit. Format: `10.20.32.245` /
  `3fce:1706:...`.

**Vacib:** portlar 0-1024 yalnız root-a — root-la işləmək təhlükəsizlik riskidir.

### 2. nc(1) — Şəbəkə Test Aləti
```bash
nc 10.10.1.123 1234        # TCP client
nc -l 1234                 # TCP server (listen)
nc -u                      # UDP rejimi
nc -v / -vv                # verbose (troubleshooting)
```

### 3. net Paketi və TCP Client
**net.Conn — həm io.Reader, həm io.Writer!** Şəbəkə bağlantısı = fayl I/O kodu.

**Dial üsulu (generik, tövsiyə):**
```go
c, err := net.Dial("tcp", connect)    // protokollar: tcp/tcp4/tcp6/udp/unix...
reader := bufio.NewReader(os.Stdin)
for {
    fmt.Print(">> ")
    text, _ := reader.ReadString('\n')
    fmt.Fprintf(c, "%s\n", text)                    // WRITE
    message, _ := bufio.NewReader(c).ReadString('\n') // READ
    fmt.Print("->: " + message)
    if strings.TrimSpace(string(text)) == "STOP" {
        return
    }
}
```

**DialTCP üsulu (TCP-spesifik):**
```go
tcpAddr, err := net.ResolveTCPAddr("tcp4", connect)
conn, err := net.DialTCP("tcp4", nil, tcpAddr)
// qalan məntiq eyni; STOP-da conn.Close()
```
Resolve+dial — eyni nəticə, fərqli API səviyyəsi. Müəllif generik `net.Dial`-i üstün
tutur.

### 4. TCP Server
**Listen üsulu:**
```go
PORT := ":" + arguments[1]
l, err := net.Listen("tcp", PORT)   // hostname-siz = bütün IP-lər
defer l.Close()
c, err := l.Accept()                // BLOKLAYIR — client gözləyir
for {
    netData, err := bufio.NewReader(c).ReadString('\n')
    if strings.TrimSpace(string(netData)) == "STOP" { return }
    t := time.Now()
    c.Write([]byte(t.Format(time.RFC3339) + "\n"))   // vaxt xidməti
}
```

**ListenTCP üsulu (echo server):**
```go
s, err := net.ResolveTCPAddr("tcp", SERVER)
l, err := net.ListenTCP("tcp", s)
buffer := make([]byte, 1024)
conn, err := l.Accept()
for {
    n, err := conn.Read(buffer)               // raw bayt oxu
    if strings.TrimSpace(string(buffer[0:n])) == "STOP" {
        conn.Close(); return
    }
    conn.Write(buffer)                        // ECHO — eyni bufferi qaytar
}
```
**Diqqət:** bu tək-client-li nümunələrdir — Accept loopun İÇİNDƏ deyil.

### 5. UDP Client və Server
**Client:**
```go
s, err := net.ResolveUDPAddr("udp4", CONNECT)
c, err := net.DialUDP("udp4", nil, s)
defer c.Close()
fmt.Printf("The UDP server is %s\n", c.RemoteAddr().String())

data := []byte(text + "\n")
c.Write(data)
buffer := make([]byte, 1024)
n, _, err := c.ReadFromUDP(buffer)     // UDP-spesifik oxu
```

**Server (random ədəd xidməti):**
```go
s, err := net.ResolveUDPAddr("udp4", PORT)
connection, err := net.ListenUDP("udp4", s)
defer connection.Close()
buffer := make([]byte, 1024)
for {
    n, addr, err := connection.ReadFromUDP(buffer)   // GÖNDERƏNİN ünvanı qayıdır
    if strings.TrimSpace(string(buffer[0:n])) == "STOP" { return }
    data := []byte(strconv.Itoa(random(1, 1001)))
    connection.WriteToUDP(data, addr)                 // həmin ünvana yaz
}
```
**UDP fərqi:** bağlantı YOXDUR — hər mesaj müstəqil; addr ilə cavab verilir; eyni server
bir neçə client-ə xidmət edir (Accept lazım deyil).

### 6. Concurrent TCP Server — Əsas Pattern
**Kitabdan kod nümunəsi:**
```go
var count = 0

func handleConnection(c net.Conn, myCount int) {
    netData, err := bufio.NewReader(c).ReadString('\n')
    if err != nil { return }
    for {
        temp := strings.TrimSpace(string(netData))
        if temp == "STOP" { break }        // yalnız BU client-in goroutine-i ölür
        c.Write([]byte("Client number: " + strconv.Itoa(myCount) + "\n"))
    }
    defer c.Close()
}

func main() {
    l, err := net.Listen("tcp4", PORT)
    defer l.Close()
    for {
        c, err := l.Accept()               // ACCEPT LOOPDA!
        if err != nil { return }
        go handleConnection(c, count)      // HƏR CLIENT = AYRI GOROUTINE
        count++
    }
}
```
**Dərs:** Accept sonsuz loopda; hər bağlantı öz goroutine-ində — server yeni client-ləri
qəbul etməkdə azad qalır. Bir client STOP dedikdə yalnız öz goroutine-i bitir, server
yaşayır. Production serverlərin esası.

### 7. WebSocket Protokolu
**Nədir:** RFC 6455; TƏK TCP bağlantı üzərində full-duplex kanal; ws:// / wss:// URL-lər.
HTTP-dən upgrade olunur — amma əksi mümkün deyil.

**Üstünlükləri:**
- Full-duplex — server client sorğusu gözləmədən data göndərə bilər
- Raw TCP — HTTP handshake overhead-i yoxdur
- Bağlantı ölənə qədər yaşayır — reopen yoxdur
- Real-time web app-lər üçün; HTML5 standardı — bütün müasir brauzerlər

### 8. WebSocket Server — gorilla/websocket
**Kitabdan kod nümunəsi:**
```go
var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool { return true },
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
    ws, err := upgrader.Upgrade(w, r, nil)   // HTTP → WebSocket UPGRADE
    if err != nil { return }
    defer ws.Close()
    for {
        mt, message, err := ws.ReadMessage()    // YALNIZ bu API!
        if err != nil { break }
        err = ws.WriteMessage(mt, message)      // echo — eyni mesajı qaytar
        if err != nil { break }
    }
}

// HTTP server: mux.Handle("/ws", http.HandlerFunc(wsHandler)) — WS endpoint
```
**QAYDA:** upgrade-dən sonra fmt.Fprintf YAZMAQ OLMAZ — yalnız ReadMessage/WriteMessage.
(golang.org/x/net/websocket alternativi var, amma funksiya cətinliklərinə görə gorilla
üstün.)

**Test aləti — websocat:**
```bash
websocat ws://localhost:1234/ws     # interaktiv echo testi
websocat -v ws://...                # verbose
```

### 9. WebSocket Client — Professional Pattern
**Kitabdan kod nümunəsi:**
```go
interrupt := make(chan os.Signal, 1)
signal.Notify(interrupt, os.Interrupt)          // SIGINT → kanal

input := make(chan string, 1)
go getInput(input)                              // stdin → input kanalı

URL := url.URL{Scheme: "ws", Host: SERVER, Path: PATH}
c, _, err := websocket.DefaultDialer.Dial(URL.String(), nil)
defer c.Close()

done := make(chan struct{})
go func() {                                     // oxu goroutine-i
    defer close(done)
    for {
        _, message, err := c.ReadMessage()
        if err != nil { return }
        log.Printf("Received: %s", message)
    }
}()

for {
    select {
    case <-time.After(4 * time.Second):         // idle sayğacı
        TIMESWAIT++
        if TIMESWAIT > TIMESWAITMAX {
            syscall.Kill(syscall.Getpid(), syscall.SIGINT)  // özünə siqnal!
        }
    case <-done:
        return
    case t := <-input:
        c.WriteMessage(websocket.TextMessage, []byte(t))
        TIMESWAIT = 0
        go getInput(input)                      // növbəti input üçün
    case <-interrupt:
        // GRACEFUL CLOSE — server-ə CloseMessage göndər:
        c.WriteMessage(websocket.CloseMessage,
            websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
        select {
        case <-done:
        case <-time.After(2 * time.Second):
        }
        return
    }
}
```
**Sub-kod izahı:** 4 kanal (input, done, interrupt, time.After) bir select-də birləşir —
professional CLI alət arxitekturası. Idle timeout öz prosesinə SIGINT göndərir →
graceful close yolu işə düşür.

### 10. RabbitMQ — Message Broker
**Nədir:** Asenkron mesaj mübadiləsi üçün açıq mənbə broker; AMQP protokolu (binary,
frame-lərlə); mesajlar FIFO növbələrdə saxlanılır — consumer oxuyana qədər EHTİYATDA
qalır. Go modulu: github.com/rabbitmq/amqp091-go.

**RabbitMQ vs Kafka:** Kafka sürətli, pull-model, batching var, payload limiti var;
RabbitMQ push-model, payload limitsiz, sadə. Sürət → Kafka; sadəlik+payload → RabbitMQ.

**İşə salma (docker-compose):**
```yaml
rabbitmq:
  image: 'rabbitmq:3.12-management'
  ports: ['5672:5672', '15672:15672']   # AMQP + management UI
  environment:
    RABBITMQ_DEFAULT_USER: "guest"
    RABBITMQ_DEFAULT_PASS: "guest"
```

**Producer (yazma):**
```go
conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
ch, err := conn.Channel()
defer ch.Close()

// QueueDeclare(name, durable, autoDelete, exclusive, noWait, args)
q, err := ch.QueueDeclare("Go", false, false, false, false, nil)

err = ch.PublishWithContext(nil, "", "Go", false, false,
    amqp.Publishing{ContentType: "text/plain", Body: []byte(message)},
)
```
**Diqqət:** növbə mövcud deyilsə YARADILIR — typo yaxalanmır!

**Consumer (oxuma):**
```go
// Consume(queue, consumer, autoAck, exclusive, noLocal, noWait, args)
msgs, err := ch.Consume("Go", "", true, false, false, false, nil)

forever := make(chan bool)
go func() {
    for d := range msgs {
        fmt.Printf("Received: %s\n", d.Body)
    }
}()
<-forever                    // proqramı saxla — yeni mesaj gözlə
```

### 11. go.mod-dan Module Silmə
```bash
go get github.com/rabbitmq/amqp091-go@none   # asılılığı SİL
go mod tidy                                   # geri gətirmək lazımsa
```

## Əsas terminlər
- TCP — etibarlı, təsdiqli, full-duplex nəqliyyat protokolu
- UDP — etibarsız, sadə, sürətli datagram protokolu
- IPv4/IPv6 — 32/128-bit ünvanlaşdırma
- net.Conn — Reader+Writer interfeysli şəbəkə bağlantısı
- net.Dial/Listen — client/server başlanğıcı (generik)
- ResolveTCPAddr/DialTCP — TCP-spesifik adres + bağlantı
- ReadFromUDP/WriteToUDP — UDP-spesifik (göndərən ünvanı ilə)
- Full-Duplex (ikiyönlü) — eyni anda iki istiqamət
- WebSocket Upgrade — HTTP bağlantının WS-ə keçirilməsi
- ReadMessage/WriteMessage — gorilla WS-in YEGANƏ data API-si
- AMQP — Advanced Message Queuing Protocol
- Queue (növbə) — FIFO mesaj anbarı; QueueDeclare/Consume
- Message Broker (mesaj brokeri) — asenkron mesajların arayıcısı
- Pull vs Push Model — consumer çəkir / broker itələyir

## Praktik nətidə

(1) net.Conn = Reader+Writer — şəbəkə kodu fayl I/O kodudur; bufio burada da işləyir.
(2) TCP serverinə concurrency: Accept loopda + hər client öz goroutine-də (concTCP pattern).
(3) UDP-də bağlantı yoxdur — addr-dən cavabla; Accept lazım deyil, amma etibarlılıq SƏNİN
vəzifəndir. (4) Port 0-1024-dən uzaq dur — root risk. (5) WebSocket: upgrade-dən sonra
yalnız ReadMessage/WriteMessage; fmt.Fprintf = ölü bağlantı. (6) WS client-də siqnal
kanalı + graceful CloseMessage = professional bağlantı bağlama. (7) RabbitMQ: Dial →
Channel → QueueDeclare → Publish/Consume; növbə adı typo-su SƏSSİZ keçir. (8) Broker
seçimi: sürət → Kafka; payload qüvvəti + sadəlik → RabbitMQ. (9) `go get pkg@none` —
asılılığı go.mod-dan çıxarır. (10) Şübhə varsa TCP/IP ilə başla, ehtiyac olanda WS-ə keç.

## Mənbə
Pages: 427-464 (PDF 458-497)
