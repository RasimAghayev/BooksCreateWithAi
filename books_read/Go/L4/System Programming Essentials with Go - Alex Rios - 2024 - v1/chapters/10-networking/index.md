# Chapter 10 — Networking (Şəbəkə Proqramlaşdırması)

## Bu chapter nədən bəhs edir?
net paketinə, TCP socket proqramlaşdırmasına (echo server), HTTP server/client qurulumuna (verb-lər + status kodları), TLS ilə təhlükəsizliyə (sertifikatlar, PEM/CRT, ListenAndServeTLS), UDP-yə (Selective Retransmissions ilə) və WebSocket-ə (gobwas/ws).

## Əsas fikirlər

### 1. net paketi — şəbəkə İsveçrə bıçağı
TCP/UDP bağlantılar, data stream-lər, ünvan parse — hamısı standart kitabxanada.

**HTTP client nümunəsi:**
```go
client := &http.Client{}
resp, err := client.Get("https://pokeapi.co/api/v2/pokemon/ditto")
defer resp.Body.Close()
body, _ := io.ReadAll(resp.Body)
fmt.Println(resp.StatusCode, string(body))
```

### 2. TCP sockets — echo server
**2 abstraksiya:** net.Conn (tək bağlantı) + net.Listener (gözləyici bouncer).

```go
listener, err := net.Listen("tcp", ":8080")
for {
    conn, err := listener.Accept()
    if err != nil { continue }
    go handleConnection(conn)       // hər bağlantı ayrıca goroutine
}

func handleConnection(conn net.Conn) {
    defer conn.Close()
    buf := make([]byte, 1024)
    for {
        n, err := conn.Read(buf)
        if err != nil { break }              // client bağladı
        conn.Write(buf[:n])                  // yalnız OXUNMUŞ hissəni qaytar
    }
}
```
**Tələlər:** bağlantını bağlamaq unudulur → resurs axır; buf[:n] (tam buf YOX) — çirkli data qayıtmır.

### 3. HTTP server — verb + status
**Minimal:**
```go
http.HandleFunc("/", handler)
http.ListenAndServe(":8080", nil)
```

**Verb idarəsi (switch r.Method):**
```go
switch r.Method {
case http.MethodGet:    ...
case http.MethodPost:    ...
case http.MethodPut:     ...
case http.MethodDelete:  ...
default:
    http.Error(w, "Unsupported HTTP method", http.StatusMethodNotAllowed)
}
```

**Status kodları — 5 sinif:**
| Sinif | Məna | Konstantlar |
|---|---|---|
| 1xx | informational | — |
| 2xx | success | StatusOK(200), StatusCreated(201), StatusAccepted(202), StatusNoContent(204) |
| 3xx | redirect | — |
| 4xx | client error | StatusNotFound(404), StatusMethodNotAllowed(405) |
| 5xx | server error | StatusInternalServerError(500) |

**CRUD → status xəritəsi:** GET→200, POST→201, PUT→202, DELETE→204, naməlum→405.

### 4. TLS — bağlantının təhlükəsizliyi
**TLS sertifikatının 2 vəzifəsi:** şifrələmə (eavesdropping qarşısı) + autentifikasiya (impostor yox, real server).

**PEM vs CRT:**
| | .crt | .pem |
|---|---|---|
| Format | DER (binary) və ya PEM | HƏMİŞƏ ASCII (Base64) |
| Marker-lər | — | BEGIN/END CERTIFICATE |
| Çoxsertifikat | — | bir faylda zəncir mümkün |
| Üstünlük | — | universal, insan-oxunarlı |

**Self-signed sertifikat (dev):**
```bash
openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365
```
**CSR axını (CA üçün):**
```bash
openssl genrsa -out mydomain.key 2048                  # 1. private key (GİZLİ!)
openssl req -new -key mydomain.key -out mydomain.csr   # 2. imza istəyi (CN = domain!)
openssl x509 -req -days 365 -in mydomain.csr -signkey mydomain.key -out mydomain.crt  # 3. self-sign
```

**HTTPS server:**
```go
err := http.ListenAndServeTLS(":8443", "cert.pem", "key.pem", nil)
```

**TLS TCP server (crypto/tls):**
```go
cert, _ := tls.LoadX509KeyPair("server.crt", "server.key")
config := &tls.Config{Certificates: []tls.Certificate{cert}}
listener, _ := tls.Listen("tcp", ":8443", config)
```

**6 TLS tələsi:** etibarlılıq (expire yenilə), private key təhlükəsizliyi, production-da CA (self-signed yox), domain uyğunluğu (CN), tam sertifikat zənciri, performans ( effektiv cipher suite).

### 5. UDP vs TCP
| | UDP | TCP |
|---|---|---|
| Protokol | connectionless | connection-oriented |
| Etibarlılıq | zəmanət YOX | sifarişli + xəta düzəltmə |
| Overhead | aşağı | yüksək |
| Sürət | daha sürətli | yavaş |
| Ssenari | real-time (oyun/VoIP/stream), broadcast/multicast, custom reliability | zəmanətli çatdırılma, data bütövlüyü |

**Go API:** net.DialUDP / net.ListenUDP / UDPConn (ReadFromUDP / WriteToUDP).

### 6. Selective Retransmissions (SACK)
**Go-Back-N (TCP klassik):** 1 paket düşsə → ondan SONRAKILAR da yenidən göndərilir (israf).
**SACK (UDP-də custom):** qəbulçu itirilmiş paketlərin SIYAHISINI deyir → yalnız onlar retransmit.

**Server nümunəsi:**
```go
const (
    maxDatagramSize = 1024
    packetLossRate  = 0.2                     // 20% itki simulyasiyası
)

type Packet struct {
    SeqNum  uint32
    Payload []byte
}

addr, _ := net.ResolveUDPAddr("udp", ":5000")
conn, _ := net.ListenUDP("udp", addr)
defer conn.Close()

// Oxu goroutine:
go func() {
    buf := make([]byte, maxDatagramSize)
    for {
        n, clientAddr, _ := conn.ReadFromUDP(buf)
        receivedSeq, _ := unpackUint32(buf[:4])   // ilk 4 bayt = seq
        sendAck(conn, clientAddr, receivedSeq)     // ACK qaytar
    }
}()

// Göndərmə (itki simulyasiyalı):
func sendPacket(conn *net.UDPConn, addr *net.UDPAddr, packet *Packet) {
    buf := make([]byte, 4+len(packet.Payload))
    binary.BigEndian.PutUint32(buf[:4], packet.SeqNum)   // network byte order!
    copy(buf[4:], packet.Payload)
    if rand.Float32() > packetLossRate {                  // 80% göndər
        conn.WriteToUDP(buf, addr)
    } else {
        fmt.Printf("Simulated packet loss, seq: %d\n", packet.SeqNum)
    }
}

func sendAck(conn *net.UDPConn, addr *net.UDPAddr, seqNum uint32) {
    ackPacket := make([]byte, 4)
    binary.BigEndian.PutUint32(ackPacket, seqNum)
    conn.WriteToUDP(ackPacket, addr)
}
```
**Big-endian:** şəbəkə bayt sırası — MSB birinci (binary.BigEndian.PutUint32/Uint32).

**UDP/TCP seçimi:** zəmanət → TCP; sürət + itki tolerantlığı → UDP (+ SACK komplexliyi).

### 7. WebSocket — real-time ikiistiqamətli
**Model:** HTTP handshake → upgrade → uzunmüddətli TCP; hər iki tərəf özbaşına göndərir; minimal framing overhead.

**Server (gobwas/ws):**
```go
go install github.com/gobwas/ws@latest

http.ListenAndServe(":8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    conn, _, _, err := ws.UpgradeHTTP(r, w)     // HTTP → WebSocket
    go func() {
        defer conn.Close()
        for {
            msg, op, err := wsutil.ReadClientData(conn)
            if err != nil { break }
            wsutil.WriteServerMessage(conn, op, msg)   // echo
        }
    }()
}))
```

**Go client:**
```go
conn, _, _, err := ws.DefaultDialer.Dial(ctx, "ws://localhost:8080")
defer conn.Close()
wsutil.WriteClientMessage(conn, ws.OpText, []byte("Hello, server!"))
response, _, _ := wsutil.ReadServerData(conn)
```

**JS client (brauzer):**
```javascript
var ws = new WebSocket('ws://localhost:8080');
ws.onopen = () => ws.send('Hello, server!');
ws.onmessage = (e) => console.log('Message from server:', e.data);
ws.onerror = (e) => console.log('Error:', e);
ws.onclose = (e) => console.log('Closed:', e);
```

**Arxitektura qeydi:** Bu bloklar üzərində REST / mesajlaşma / gRPC (RPC) qurulur — bazanı bilmək troubleshooting-də aydın mental xəritə verir.

## Əsas terminlər
- net.Conn / net.Listener — TCP abstraksiyaları
- http.Client / client.Get / io.ReadAll
- HTTP verbs: GET/POST/PUT/DELETE/PATCH
- Status kodları: 200/201/202/204/404/405/500
- TLS: şifrələmə + autentifikasiya
- PEM (ASCII/Base64) vs CRT (DER/binary)
- CSR (Certificate Signing Request) + CN (Common Name)
- openssl req -x509 / genrsa / x509 -req
- ListenAndServeTLS / tls.Listen / LoadX509KeyPair
- UDP: connectionless, DialUDP/ListenUDP/ReadFromUDP/WriteToUDP
- Go-Back-N vs Selective Retransmissions (SACK)
- binary.BigEndian — network byte order
- WebSocket: upgrade, gobwas/ws, wsutil
- ws.UpgradeHTTP / DefaultDialer.Dial / OpText

## Praktik nəticə
1. TCP serverlərdə hər Accept üçün goroutine + defer Close; buf[:n] yaz (tam buf YOX).
2. HTTP cavablarında verb-uyğun status kodu: POST→201, DELETE→204 — API istehlakçısı bunlara etibar edir.
3. Dev-də self-signed PEM; production-da CA imzalı sertifikat + tam zəncir + CN domain uyğunluğu.
4. UDP seçəndə SACK-dəstəkli seq/ACK protokolu öz üzərinə götürürsən — TCP-nin pulsuz xidmətlərindən imtina edirsən.
5. Şəbəkə formatında həmişə BigEndian — binary.BigEndian alətləri.
6. Real-time üçün WebSocket (HTTP upgrade + long-lived TCP); browser JS WebSocket API hazırdır.

## Mənbə
Pages: 191-219 (PDF səh. 212-241)
