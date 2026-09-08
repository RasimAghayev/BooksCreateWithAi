# Chapter 11 — Working with Sockets (səh. 170-182)

## Bu fəsil nədən bəhs edir?

TCP/UDP socket proqramlaşdırması: TCP üzərində fayl qəbul/göndərmə
(öz protokolu ilə), Unix domain socket üzərində mini-RPC (JSON serverless),
NTP protokolu ilə UDP-dən vaxt oxunması (binary.BigEndian) və "Fallacies
of distributed computing".

## Əsas fikirlər

### Recipe 57 — TCP-də fayl qəbulu (server)
**Protokol:** `2020-01-01-httpd.log 137383 ` header (ad + ölçü) + məzmun;
cavab: status sətri.

```go
addr := ":8765"
ln, err := net.Listen("tcp", addr)
if err != nil {
    log.Fatalf("error: %s", err)
}
log.Printf("server ready on %s", addr)

for {                                  // sonsuz qəbul döngüsü
    c, err := ln.Accept()
    if err != nil {
        log.Fatalf("error: %s", err)
    }
    log.Printf("new connection from %s", c.RemoteAddr())
    go fileHandler(c)                  // hər bağlantı ayrıca goroutine
}

func fileHandler(c net.Conn) {
    defer c.Close()

    // Header oxu (socket = Reader!):
    var fileName string
    var size int64
    if _, err := fmt.Fscanf(c, "%s %d", &fileName, &size); err != nil {
        fmt.Fprintf(c, "error: %s\n", err)
        return
    }
    if size > maxSize {                 // ölçü həddi
        fmt.Fprintf(c, "error: size %d > %d\n", size, maxSize)
        return
    }

    // Məzmunu kopyala (dəqiq size bayt):
    fileName = path.Join("logs", strings.TrimSpace(fileName))
    file, err := os.Create(fileName)
    if err != nil {
        fmt.Fprintf(c, "error: %s\n", err)
        return
    }
    defer file.Close()
    n, err := io.CopyN(file, c, size)   // CopyN — məlum say!
    if err != nil {
        fmt.Fprintf(c, "error: %s\n", err)
        return
    }
    if n != size {
        fmt.Fprintf(c, "error: %d of %d bytes", n, size)
        return
    }
    fmt.Fprintf(c, "ok: %d bytes written to %s\n", n, fileName)
}
```

**Test:**
```bash
$ echo 'pi.txt 4 3.14' | nc localhost 8765
ok: 4 bytes written to logs/pi.txt
```
- TCP-də hər şeyi ÖZÜN edirsən: protokol, xəta yoxlaması...
- Socket və fayl hər ikisi Reader/Writer → fmt.Fscanf, io.CopyN işləyir

### Recipe 58 — TCP-də fayl göndərilməsi (client)
```go
func sendFile(addr, fileName string) error {
    c, err := net.Dial("tcp", addr)     // "tcp" / "udp" / "unix"
    if err != nil {
        return err
    }
    defer c.Close()

    file, err := os.Open(fileName)
    if err != nil {
        return err
    }
    defer file.Close()
    info, err := file.Stat()
    if err != nil {
        return err
    }

    // Header göndər:
    _, err = fmt.Fprintf(c, "%s %d ", fileName, info.Size())
    if err != nil {
        return err
    }

    // Məzmunu axıt:
    _, err = io.Copy(c, file)
    if err != nil {
        return err
    }

    // Cavabı oxu (limitli!):
    const maxReply = 1 << 10                    // 1KB
    data, err := io.ReadAll(io.LimitReader(c, maxReply))
    if err != nil {
        return err
    }
    reply := string(data)
    if strings.HasPrefix(reply, "error") {
        return fmt.Errorf(reply)
    }
    return nil
}
```
- `net.Conn` → ümumi interfeys; spesifik funksiya üçün assertion:
```go
tc, ok := c.(*net.TCPConn)
if ok {
    tc.SetKeepAlive(true)     // yalnız TCPConn-də var
}
```

### Recipe 59 — serverless platform (Unix socket + JSON RPC)
**Dizayn:** RPC = serializasiya formatı + transport; JSON (çoxdilli) +
Unix domain socket (lokal).

```go
func processMessages(conn io.ReadWriteCloser, ch <-chan Message) error {
    dec := json.NewDecoder(conn)     // həm kod, həm göndər
    enc := json.NewEncoder(conn)

    for msg := range ch {
        if err := enc.Encode(msg); err != nil {    // göndər
            return err
        }
        var reply struct {
            Output any
        }
        if err := dec.Decode(&reply); err != nil { // cavab al
            return err
        }
        log.Printf("%#v -> %#v", msg, reply.Output)
    }
    return nil
}

// Server-i başlat (exec.Command):
cmd := exec.Command("go", "run", "server/main.go", "-socket", socketFile)
cmd.Stdout = os.Stdout
cmd.Stderr = os.Stderr
if err := cmd.Start(); err != nil {
    log.Fatalf("error: %s", err)
}
defer cmd.Process.Kill()
time.Sleep(time.Second)              // FIXME: düzgün gözləmə yaz

// Unix socket-ə qoşul:
sock, err := net.Dial("unix", socketFile)
defer sock.Close()

ch := make(chan Message)
go func() {
    ch <- Message{Text: "Was it a cat I saw?"}
    close(ch)
}()
processMessages(sock, ch)
// main.Message{Text:"Was it a cat I saw?"} -> "?was I tac a ti saW"
```
- JSON Decoder/Encoder socket üzərində birbaşa (sətir-sətir)
- Yarımçıq: xəta idarəetmə, timeout, crash detection, streaming —
  yetkin RPC framework hamısını həll edir
- RPC yalnız proseslər ARASI / çoxdilli hallarda; eyni prosesdə lazım deyil

### Recipe 60 — UDP-də NTP vaxtı
**Protokol:** RFC5905 — binary, UDP. NTP epoch = 1900; Unix = 1970!

```go
// NTPMessage is an NTP message.
type NTPMessage struct {
    VNMode             uint8
    Stratum            uint8
    Poll               uint8
    Precision          uint8
    RootDelay          uint32
    RootDispersion     uint32
    RefID              uint32
    RefTimeSec         uint32
    RefTimeFrac        uint32
    OrigTimeSec        uint32
    OrigTimeFrac       uint32
    ReceivedTimeSec    uint32
    ReceivedTimeFrac   uint32
    TransmitTimeSec    uint32
    TransmitTimeFrac   uint32
}

func ntpDelta() time.Duration {           // 1900 ↔ 1970 fərqi
    unixEpoch := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
    ntpEpoch := time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)
    return ntpEpoch.Sub(unixEpoch)
}
var NTPDelta = ntpDelta()

// TransmitTime returns the transmit time.
func (n *NTPMessage) TransmitTime() time.Time {
    secs := int64(n.TransmitTimeSec)
    nanos := (int64(n.TransmitTimeFrac) * 1e9) >> 32   // fraksiya → ns
    return time.Unix(secs, nanos).Add(NTPDelta)
}

// CurrentTime returns the current time from NTP host.
func CurrentTime(addr string) (time.Time, error) {
    conn, err := net.Dial("udp", addr)
    if err != nil {
        return zeroTime, err
    }
    defer conn.Close()

    msg := NTPMessage{
        VNMode: 0b00011011,   // li=0, vn=3, mode=3 (client)
    }

    // Binary yaz/oxu — BigEndian (şəbəkə bayt sırası):
    if err := binary.Write(conn, binary.BigEndian, msg); err != nil {
        return zeroTime, err
    }
    if err := binary.Read(conn, binary.BigEndian, &msg); err != nil {
        return zeroTime, err
    }
    return msg.TransmitTime(), nil
}
```
- UDP: zəmanət YOXDUR (sıra/çatdırma) — sürətli; NTP standartı olduğundan
  UDP-dədir; HTTP/3 də UDP-yə keçir
- binary.Write/Read + BigEndian — struct ↔ baytlar
- Real həyatda: ntpd servisi (bu resept UDP nümunəsidir)

## Final Thoughts-dən

Socket işi HTTP/gRPC-dən AŞAĞI səviyyədir — mümkünsə hazır protokol seçin.
Gözəl cəhət: bağlantı qurulandan sonra hər şey Reader/Writer dünyasına
keçir.

**Fallacies of distributed computing (L. Peter Deutsch):**
1. Şəbəkə etibarlıdır  2. Gecikmə sıfırdır  3. Zolaq sonsuzdur
4. Şəbəkə təhlükəsizdir  5. Topologiya dəyişmir  6. Bir admin var
7. Nəqliyyət pulsuzdur  8. Şəbəkə homojendir

Bu yanlış inanclara görə `io.LimitReader` kimi qorunmalar yazılır.

## Əsas terminlər
- net.Listen / net.Dial — server/klient bağlantısı
- net.Conn — ümumi socket interfeysi (Reader/Writer!)
- ln.Accept — bağlantı qəbulu
- io.CopyN — dəqiq N bayt kopyalama
- Protokol dizaynı — header formatı (ad + ölçü)
- Unix domain socket — lokal IPC
- json.NewDecoder/Encoder — socket üzərində axın
- exec.Command — xarici proses başlatma
- UDP vs TCP — zəmanətsiz/sürətli ↔ etibarlı
- binary.Write/Read + BigEndian — binary serializasiya
- NTP epoch (1900) — Unix epoch-dan 70 il əvvəl
- Fallacies of distributed computing — 8 yanlış inanc

## Praktik nəticə
TCP-də öz protokolu: header (fmt.Fprintf) + məzmun (io.Copy) + status
cavabı; server = Listen + Accept loop + hər bağlantı goroutine. UDP-də:
Dial + binary.Write/Read; zəmanət yoxdur — öz yoxlamanı yaz. Socket qoşulan
kimi Reader/Writer kimi davranır — bütün standart alətlər işləyir. RPC
qurmaq istəsəniz: serializasiya + transport seçimi; amma gRPC kimi hazır
framework-un həll etdiklərini (timeout, retry, streaming) özünüz yazmalı
qalacaqsınız. LimitReader ilə girişi hədləyin — "şəbəkə təhlükəsizdir"
fallacy-sinə düşməyin.

## Mənbə
Pages: 170-182 (PDF 170-182)
