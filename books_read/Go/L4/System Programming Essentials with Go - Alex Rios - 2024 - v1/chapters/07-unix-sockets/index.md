# Chapter 7 — Unix Sockets (Unix Socket-lər)

## Bu chapter nədən bəhs edir?
UNIX domain socket-lərə (IPC üçün TCP/IP alternativi), server/client qurulumuna (graceful shutdown ilə), lsof ilə socket inspeksiyasına, çoxşahidlə chat server-in inkişafına (5 mərhələ), HTTP-in unix socket üzərində serve edilməsinə (raw HTTP + textproto) və performans üstünlüklərinə.

## Əsas fikirlər

### 1. UNIX socket nədir
**Nədir:** Eyni maşında proseslərarası sürətli kommunikasiya — TCP/IP-nin LOKAL alternativi (Unix/Linux-özəl).

**Xüsusiyyətlər:**
- Stream (TCP kimi) və ya datagram (UDP kimi)
- **Filesystem node kimi təmsil olunur** (adi fayl deyil — IPC mexanizmi)
- **Effektivlik:** şəbəkə protokolu overhead-i yoxdur
- **Filesystem namespace:** path ilə istinad — asan tapılır, amma silinmədiyi qədər qalır
- **Təhlükəsizlik:** fayl icazələri ilə giriş nəzarəti (user/group ID əsaslı)

### 2. Server qurulumu — 5 addım
```go
// 1. Socket path + cleanup:
socketPath := "/tmp/example.sock"
os.Remove(socketPath)                    // köhnə socket faylını sil

// 2. Listen:
listener, err := net.Listen("unix", socketPath)
defer listener.Close()

// 3. Graceful shutdown:
signals := make(chan os.Signal, 1)
signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
go func() {
    <-signals
    listener.Close()
    os.Remove(socketPath)                 // socket faylını TƏMİZLƏ
    os.Exit(0)
}()

// 4. Accept loop:
for {
    conn, err := listener.Accept()
    if err != nil { continue }
    go handleConnection(conn)
}

// 5. Handler:
func handleConnection(conn net.Conn) {
    defer conn.Close()
    buffer := make([]byte, 1024)
    n, err := conn.Read(buffer)
    if err != nil { return }
    fmt.Println("Received:", string(buffer[:n]))
    conn.Write([]byte("Message received successfully\n"))
}
```

### 3. Arxaplanda nə olur (OS + filesystem baxışı)
- **OS perspektivi:** socket KERNEL-də daxili resurs kimi yaranır — yaddaşda abstraksiya, hələ faylla bağlı deyil
- **Filesystem perspektivi:** bind = socket-i path ilə əlaqələndirir → xüsusi "socket faylı" yaranır (data saxlamır — socketin ADLI təmsilidir)

### 4. Client
```go
conn, err := net.Dial("unix", "/tmp/example.sock")
defer conn.Close()
conn.Write([]byte("Hello UNIX socket!\n"))

buffer := make([]byte, 1024)
n, err := conn.Read(buffer)
fmt.Println("Server response:", string(buffer[:n]))
```

### 5. lsof — socket inspeksiyası
```bash
lsof -Ua /tmp/example.sock
# -U: yalnız UNIX socket-lər; -a: şərtləri AND-la
# PID + server/client qoşulmaları görünür
```

### 6. Chat server — 5 mərhələli inkişaf

**Mərhələ 1 — əsas server:** Listen + Accept + Close.

**Mərhələ 2 — tək client oxu:** handleConnection-də Read + çap.

**Mərhələ 3 — çoxlu client (concurrent):**
```go
var (
    clients []net.Conn
    mutex  sync.Mutex
)
// Accept-də:
mutex.Lock()
clients = append(clients, conn)
mutex.Unlock()
go handleConnection(conn)
```

**Mərhələ 4 — broadcast:**
```go
func handleConnection(conn net.Conn) {
    defer conn.Close()
    buffer := make([]byte, 1024)
    for {
        n, err := conn.Read(buffer)
        if err != nil {
            removeClient(conn)             // disconnect təmizliyi
            break
        }
        broadcastMessage(string(buffer[:n]))
    }
}

func broadcastMessage(message string) {
    mutex.Lock()
    defer mutex.Unlock()
    for _, client := range clients {
        client.Write([]byte(message + "\n"))
    }
}
```

**Mərhələ 5 — mesaj tarixçesi:**
```go
// handleConnection başlanğıcında:
for _, msg := range messageHistory {
    conn.Write([]byte(msg + "\n"))        // yeni clientə KEÇMİŞ göndər
}
```

### 7. Chat client (tam)
```go
conn, err := net.Dial("unix", socketPath)
defer conn.Close()

var wg sync.WaitGroup
wg.Add(1)

go func() {                                 // oxu: server → console
    defer wg.Done()
    scanner := bufio.NewScanner(conn)
    for scanner.Scan() {
        fmt.Println("Message from server:", scanner.Text())
    }
}()

scanner := bufio.NewScanner(os.Stdin)      // yaz: stdin → server
for scanner.Scan() {
    conn.Write([]byte(scanner.Text()))
}
wg.Wait()
```

### 8. HTTP over UNIX socket — server
**Niyə:** IP/port idarəsi yox; TCP port tükənməsi yox; kiçik gecikmə; filesystem icazələri ilə API nəzarəti; maşın-özəl servislər üçün ideal.

```go
http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Hello, world!"))
})

listener, err := net.Listen("unix", "/tmp/go-server.sock")

// graceful shutdown (eyni pattern)...
err = http.Serve(listener, nil)             // listener-i http.Serve-ə ötür!
if err != nil && err != http.ErrServerClosed {
    log.Fatal("HTTP server error:", err)
}
```

### 9. Raw HTTP client + textproto
```go
conn, err := net.Dial("unix", socketPath)
defer conn.Close()

// 1. Manual HTTP sorğusu:
request := "GET / HTTP/1.1\r\n" +
    "Host: localhost\r\n" +
    "\r\n"                                  // boş sətir = header sonu!
conn.Write([]byte(request))

// 2. Cavab parse:
reader := bufio.NewReader(conn)
tp := textproto.NewReader(reader)

statusLine, _ := tp.ReadLine()               // "HTTP/1.1 200 OK"
headers, _ := tp.ReadMIMEHeader()           // header-lər map kimi
for key, values := range headers {
    for _, value := range values {
        fmt.Printf("%s: %s\n", key, value)
    }
}

for {                                        // body sətir-sətir
    line, err := reader.ReadString('\n')
    if err != nil { break }
    fmt.Print(line)
}
```
**HTTP request anatomiyası:**
- `GET / HTTP/1.1\r\n` — metod + path + protokol + CRLF
- `Host: localhost\r\n` — HTTP/1.1-də MƏCBURİ (virtual hosting)
- `\r\n` — boş sətir: header bitir, body başlayır (GET-də request bitir)

**textproto niyə:** CRLF nüansları avtomatik; MIME header parse; bufio inteqrasiyası; HTTP-dən əlavə digər text protokollar üçün də.

### 10. Performans + istifadə sahələri
**Loopback-dən (localhost TCP) daha sürətlidir:** loopback belə TCP/IP stack-dən keçir (TCP seqmentləmə + IP paketləmə); unix socket kernel DAXİLİNDE birbaşa; bəzi hallarda zero-copy mümkün.

**Real istifadəçilər:** System V IPC, X11 (windowing), D-Bus (Linux mesaj bus), systemd (init), MySQL/PostgreSQL (lokal client-server), Docker socket.

## Əsas terminlər
- UNIX domain socket — filesystem-path-ə bağlı IPC
- net.Listen("unix", path) / net.Dial("unix", path)
- Socket faylı — socketin filesystem təmsili (data saxlamır)
- Graceful shutdown — signal.Notify + Close + Remove
- lsof -Ua — socket inspeksiyası
- clients slice + mutex — qoşuluşların qorunması
- Broadcast pattern — bütün clientlərə yazı
- Message history — yeni qoşulana keçmiş
- http.Serve(listener, nil) — HTTP over unix socket
- Raw HTTP: request line + headers + boş sətir (\r\n)
- net/textproto — text protokol parser (ReadLine/ReadMIMEHeader)
- Zero-copy — kernel-user köçürməsiz ötürünmə
- D-Bus / systemd / X11 — real istifadəçilər

## Praktik nəticə
1. Unix socket başlanğıcında KÖHNƏ faylı sil (os.Remove) — bind xətasının qarşısı; shutdown-da da təmizlə.
2. Hər unix-socket serverdə signal.Notify + graceful shutdown — socket faylı kirlənməsin.
3. Çoxlu client: slice + mutex + hər qoşuluş üçün ayrıca goroutine; disconnect-də client sil.
4. HTTP server unix socketdə: net.Listen + http.Serve(listener, nil) — portlardan azad ol.
5. Raw protokol clienti üçün textproto — CRLF/MIME manual parse etmə.
6. Lıkal IPC üçün localhost TCP əvəzinə unix socket — loopback belə stack overhead daşıyır.

## Mənbə
Pages: 123-144 (PDF səh. 144-165)
