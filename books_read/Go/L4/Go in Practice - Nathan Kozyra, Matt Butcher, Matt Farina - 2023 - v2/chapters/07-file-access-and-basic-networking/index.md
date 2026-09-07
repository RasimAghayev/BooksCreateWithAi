# Chapter 7 — File access and basic networking (Fayl girişi və əsas şəbəkələşmə)

## Bu chapter nədən bəhs edir?

Fayl oxuma (ReadFile, Open+Stat, bayt-chunk JSONS, bufio.Scanner), fayl yazma (Create,
io.Copy), TCP log client (net.Dial + log.New(conn)), back pressure problemi, UDP ilə
logging (DialTimeout), websocket chat app (golang.org/x/net/websocket) və server-sent
events (fsnotify + Flusher).

## Əsas fikirlər

### 1. Fayl Oxuma — 3 Yol
**os.ReadFile:** Bütün fayl → []byte; ən sadə; gözlənilməz EOF-u tutmur; böyük faylda
prealloc yaddaş riski.

**os.Open + Stat:** Fayl deskriptoru + metadata (Name/Mode/Size/IsDir) → ölçüyə görə
strategiya seç; böyük faylı stream et.

**Bayt-chunk JSONS oxunuşu (le dərsi):** 2-bayt artımlarla Unmarshal yoxla — kövrəkdir
(malformed JSON, 4-bayt unicode); real parser lazım olardı. Qeyd: Go 1.21 `clear()` —
map-də entry-ləri silir, slice-da zeroValue; []byte{} REASSIGN (clear zibil baytları buraxar).

**bufio.Scanner (idiomatik):**
```go
file, err := os.Open("structured.log")
defer file.Close()
scan := bufio.NewScanner(file)
scan.Split(bufio.ScanLines)            // CSV/JSONS kimi sətir formatları üçün
lineJSON := make(map[string]interface{})
for scan.Scan() {
    if err := json.Unmarshal([]byte(scan.Text()), &lineJSON); err != nil {
        log.Println(err)                // tək səhv sətir iqnor — davam
    } else {
        log.Println(lineJSON["level"])
    }
}
```
Böyük faylları parçalayıb proseslər arasında paylamaq olar.

### 2. Fayl Yazma
```go
file, err := os.Create("test.txt")     // yeni/mövcud fayl
defer file.Close()
file.WriteString("test")
```
**Stream copy:** `io.Copy(dest, src)` — Reader-dan Writer-a buffer ilə; boilerplatesiz
tam köçürmə.

### 3. TCP ilə Şəbəkə Logging
**Test server:** `nc -lk 1902` (Netcat listener).

**Kitabdan kod nümunəsi:**
```go
conn, err := net.Dial("tcp", "localhost:1902")
defer conn.Close()                       // panic-də belə buffer FLUSH olunur
f := log.Ldate | log.Lshortfile
logger := log.New(conn, "example ", f)   // Writer = socket!
logger.Println("This is a regular message.")
logger.Panicln("This is a panic.")
```
- log Writer dəyişmək = konfiq dəyişmə (fayl ↔ socket eyni interfeys)
- Host timestamp log serverdən asılı olmayaraq hadisə rekonstruksiyası üçün vacib
- **log.Fatal YOX, Panicln:** Fatal os.Exit çağırır → defer-lər İCRA OLMAZ → socket
  bağlanmamış/qalan mesajlar İTİR. Panic defer-ləri işə salır.

### 4. Back Pressure
**Nədir:** Server emal girişdən geri qalır; TCP Ack gecikir; client Ack gözləyib
bloklanır → resurslar zəncirlə bağlı sıxılır (damba arxasında su).

### 5. UDP ilə Logging
```go
timeout := 30 * time.Second
conn, err := net.DialTimeout("udp", "localhost:1902", timeout)
```
Server: `nc -luk 1902`. UDP: bağlantı YOX, Ack YOX — client back pressure-dan azaddır.

**Üstünlüklər:** back pressure + outage müqaviməti; sürətli; sadə kod.
**Çatışmazlıqlar:** mesaj İTƏ bilər; sıra pozula bilər (timestamp qismən kömək);
log server özü köklənə bilər (client azaddır, server yox).
**Seçim fəlsəfəsi:** GIF/PNG (dəqiq, böyük) vs JPEG (kiçik, itkili) — TCP (zəmanətli,
back pressure) vs UDP (sürətli, itkili). Middle yol: böyük buffer + müvəqqəti anbar.

### 6. Websocket Chat (golang.org/x/net/websocket)
**Nədir:** Tək TCP bağlantıda davamlı İKİYÖNLÜ əlaqə; HTTP polling-in təkrar/overhead
problemini həll edir; REST-ə TAM əvəz YOX — real-time, yüngül yayım üçün.

**Kitabdan kod nümunəsi (arxitektura <200 sətir):**
```go
var clients map[string]*websocket.Conn     // ID → connection
func ws(ws *websocket.Conn) {
    id := generateId()                       // random 16-simvol ID
    clients[id] = ws
    sendToClients(joinMsg)                   // joinleave + member list
    for {
        var incoming string
        if err := websocket.Message.Receive(ws, &incoming); err != nil {
            disconnectClient(id); break      // receive xətası = disconnect
        }
        if err := sendToClients(msg); err != nil {
            disconnectClient(id); break
        }
    }
}
func sendToClients(msg servermsg) error {
    msgJSON, _ := json.Marshal(msg)
    for k := range clients { websocket.Message.Send(clients[k], string(msgJSON)) }
}
// main: http.Handle("/ws", websocket.Handler(ws))  — handshake/upgrade wrapper
```
- servermsg struct JSON tag-ləri ilə (message_type: message/joinleave)
- **Təhlükəsizlik:** ws:// lokal; production wss:// (TLS); **cookie YOX** → identitet
  JWT (session hijack riskinə qarşı imzalı; invalidation öz riski) ilə qurulmalı
- **Stateful problem:** server bağlansa client bilmir → reconnect logic frontend-də
  mütləq; köhnə brauzer polyfill-i REST-ə deqradasiya edir

### 7. Server-Sent Events (SSE / EventSource)
**Nədir:** Uzunömürlü HTTP bağlantı, YALNIZ server→client; websocket-dən yüngül;
bildiriş/növbə yeniləməsi üçün ideal (sosial media cavabları kimi).

**Kitabdan kod nümunəsi (fayl dəyişikliyi bildirimi — fsnotify):**
```go
func sseHandler(w http.ResponseWriter, r *http.Request) {
    flusher, ok := w.(http.Flusher)          // uyğunluq yoxlaması!
    if !ok {
        http.Error(w, "byte streams not supported", http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "text/event-stream")
    changes := make(chan FileUpdateInfo)
    go fileListener(r.Context(), changes)    // ctx.Done() ilə çıxan watcher
    for change := range changes {
        changeJSON, _ := json.Marshal(change)
        fmt.Fprintf(w, "event: file-update\ndata: %s\n\n", changeJSON)
        flusher.Flush()                        // dərhal göndər
    }
}
func fileListener(ctx context.Context, changes chan<- FileUpdateInfo) {
    for {
        select {
        case <-ctx.Done(): return
        case event, ok := <-watcher.Events:
            if !ok { break }
            changes <- FileUpdateInfo{event.Name, event.Op.String(), size}
        }
    }
}
```
- Format: `event: AD\ndata: JSON\n\n`; frontend: `new EventSource('/sse')` +
  `addEventListener('file-update', ...)`
- **Məhdudiyyət:** brauzer başına domain-ə EventSource limiti (2-10; VƏ bütün
  tab-lar arasında!) — HTTP/2 bu limiti aradan qaldırır; kontingensiya planla.
- Flusher-ə type assertion — dəstəksiz client üçün erkən çıxış.

## Əsas terminlələr
- JSONS — newline-ayrı JSON mesaj axını
- bufio.Scanner/Split — stream sətir oxuyucusu
- Back pressure — emal girişdən geri qalanda axın sıxılması
- Ack — TCP təsdiq paketi (UDP-da yoxdur)
- net.Dial/DialTimeout — TCP/UDP bağlantı (timeout ilə)
- Websocket handshake/upgrade — HTTP→WS keçidi
- Flusher — dərhal göndərmə interfeysi (SSE uyğunluq siqnalı)
- EventSource — SSE üçün JS API
- fsnotify — fayl dəyişikliyi watcher paketi

## Praktik nətidə

Fayl/şəbəkə qərarları: (1) bütün fayl kiçikdirsə ReadFile; (2) böyük/sətir-formatlı —
os.Open + bufio.Scanner (tək səhv sətiri iqnor et); (3) kopyalama — io.Copy; (4) şəbəkə
log: TCP (zəmanət) vs UDP (sürət) — Fatal YOX, Panicln (defer flush!); (5) real-time
İKİYÖNLÜ — websocket (JWT identitet, reconnect planı, wss://); (6) BİRYÖNLÜ bildiriş —
SSE (Flusher + text/event-stream + HTTP/2 limiti unutma); (7) bunların hamısı net/http
alətlərinin üstündə qurulur — Chapter 8-10 davamı.

## Mənbə
Pages: 169-192 (PDF 190-213)
