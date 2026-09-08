# Chapter 6 — Understanding Pipes in Inter-Process Communication (IPC-də Borular)

## Bu chapter nədən bəhs edir?
Anonymous pipe-ların mexanikasına (parent-child IPC), exec.Command ilə pipeline qurulumuna, named pipe-lara (Mkfifo — task mailbox), chunking/compression best practice-lərinə, timeout/secure pipe istifadəsinə və real-time log processing alətinə.

## Əsas fikirlər

### 1. Pipe nədir
**Nədir:** Yaddaşda data nəql edən kəsik — producer-consumer modeli; **unidirectional** (bir istiqamət: write-end → read-end).

**İstifadə:** CLI zəncirləri (`cat file.txt | grep "flower"`), data streaming, proseslərarası mübadilə.

### 2. Pipe vs Go channel
**Oxşarlıqlar:** kommunikasiya mexanizmi, data ötürmə, sinxronizasiya (dolu pipe/kanal yazanı bloklayır), buffering.
**Fərqlər:**

| | Pipe | Channel |
|---|---|---|
| İstiqamət | unidirectional | bidirectional (default) |
| Səviyyə | OS/syscall (çoxdilli, proseslərarası) | dil-native (goroutine-lərarası) |

**Seçim qaydası:**
- **Pipe:** ayrı proseslər (fərqli dillər ola bilər), ayrı executable-lər, Unix IPC
- **Channel:** tək Go proqramı daxilində goroutine sinxronu; fan-in/out, worker pool patternləri

### 3. Anonymous pipes — exec ilə
**Go API-lər:** `io.Pipe()` (in-memory, syscall-sız) və `os.Pipe()` (SYS_PIPE2 əsaslı).

**cat | grep replikası (manual pipe):**
```go
echoCmd := exec.Command("echo", "Hello, world!")
grepCmd := exec.Command("grep", "Hello")

pipe, err := echoCmd.StdoutPipe()   // echo → pipe
grepCmd.Stdin = pipe                 // pipe → grep
grepOut, err := grepCmd.StdoutPipe() // grep → nəticə

grepCmd.Start()      // gözləmədən başlat
echoCmd.Run()        // echo icra (output pipe-a)
pipe.Close()          // EOF siqnalı
scanner := bufio.NewScanner(grepOut)
for scanner.Scan() { fmt.Println(scanner.Text()) }
grepCmd.Wait()
```
**Axın:** echo stdout → pipe → grep stdin → grepOut → scanner.

**Sadə alternativ (Output() metodu):**
```go
echoOutput, _ := echoCmd.Output()              // echo-nu icra + output tut
grepCmd.Stdin = strings.NewReader(string(echoOutput))  // output-u input kimi
grepOutput, _ := grepCmd.Output()
```

### 4. Named pipes (Mkfifo) — task mailbox
**Anonymous məhdudiyyətləri:** yalnız yaradan proses + onun övladları yaşayarkən; unidirectional.
**Named həll:** filesystem-də mövcuddur — İSTƏNILƏN proseslər arası, prosesdən asılı deyil.

**Office mailbox analogiyası:** İşçi tapşırığı poçt qutusuna atır → digər işçi götürüb icra edir — birbaşa əlaqə olmadan.

**Yaratma + oxuma:**
```go
func namedPipeExists(pipePath string) bool {
    _, err := os.Stat(pipePath)
    return err == nil || !os.IsNotExist(err)
}

if !namedPipeExists(mailboxPath) {
    unix.Mkfifo(mailboxPath, 0666)          // FIFO yarat
}
mailbox, err := os.OpenFile(mailboxPath, os.O_RDWR, os.ModeNamedPipe)
defer mailbox.Close()
```

**Writer + Reader (WaitGroup ilə):**
```go
// writer.go:
func SendTask(pipe *os.File, data string) error {
    _, err := pipe.WriteString(data)
    return err
}

// reader.go:
func ReadTask(pipe *os.File) error {
    scanner := bufio.NewScanner(pipe)
    for scanner.Scan() {
        task := scanner.Text()
        fmt.Printf("Processing task: %s\n", task)
        if task == "EOD" { break }          // end-of-day sentinel
    }
    return scanner.Err()
}

// main:
wg := &sync.WaitGroup{}
wg.Add(2)
go func() { defer wg.Done(); ReadTask(mailbox) }()
go func() {
    defer wg.Done()
    for i := 0; i < 10; i++ {
        SendTask(mailbox, fmt.Sprintf("Task %d\n", i))
    }
    SendTask(mailbox, "EOD\n")              // sonlanma siqnalı
}()
wg.Wait()
```
**İki istiqamət üçün:** 2 named pipe (hərəsi tək istiqamət).

### 5. Best practice-lər

#### a) Chunking — bufer şişməsinin qarşısı
```go
func writeInChunks(pipe *os.File, data []byte, chunkSize int) error {
    for i := 0; i < len(data); i += chunkSize {
        end := min(i+chunkSize, len(data))
        if _, err := pipe.Write(data[i:end]); err != nil { return err }
    }
    return nil
}

// Oxu tərəfi (newline-ayrılmış chunk-lar):
for {
    chunk, err := reader.ReadBytes('\n')
    if err == io.EOF { break }
    fmt.Printf("Received chunk: %s\n", string(chunk))
}
```

#### b) Compression — gzip qatı
```go
// Writer:
gzipWriter := gzip.NewWriter(fifo)
gzipWriter.Write(data)
gzipWriter.Flush()

// Reader:
gzipReader, _ := gzip.NewReader(fifo)
defer gzipReader.Close()
io.Copy(&buf, gzipReader)
```
**Trade-off:** sıxılma dərəcəsi vs CPU overhead — hər data yaxşı sıxılmır.

#### c) Timeout ilə oxu (deadlock qarşısı)
```go
timeout := time.After(5 * time.Second)
done := make(chan bool)
go func() {
    _, err := pipe.Read(buffer)
    done <- true
}()
select {
case <-timeout:     // 5s içində oxu gəlmədi
    // pipe bağla / xəta loqla
case <-done:        // uğur
}
```

#### d) Resurs idarəsi
```go
pipeReader, pipeWriter, _ := os.Pipe()
defer pipeReader.Close()
defer pipeWriter.Close()
```
Açıq qalan pipe-lar → fd tükənməsi (file descriptor exhaustion).

#### e) Təhlükəsizlik
- Sensitive data → pipe-dan əvvəl encrypt et; qəbuldan sonra validate et
- Named pipe icazələri: `Mkfifo(path, 0600)` — yalnız owner
- **Name squatting hücumu:** hücumçu gözlənilən adda pipe yaradır → random suffix: `/tmp/my_secure_pipe_ + randomString(10)`

### 6. Real-time log processing aləti
```go
func filterLogs(reader io.Reader, writer io.Writer) {
    scanner := bufio.NewScanner(reader)
    for scanner.Scan() {
        logEntry := scanner.Text()
        if strings.Contains(logEntry, "ERROR") {    // severity filtri
            writer.Write([]byte(logEntry + "\n"))
        }
    }
}

func main() {
    pipePath := "/tmp/my_log_pipe"
    os.RemoveAll(pipePath)
    os.Mkfifo(pipePath, 0600)                     // secure permission
    defer os.RemoveAll(pipePath)

    pipeFile, _ := os.OpenFile(pipePath, os.O_RDONLY|os.O_CREATE, os.ModeNamedPipe)
    defer pipeFile.Close()

    go func() {                                    // log writer simulyasiyası
        writer, _ := os.OpenFile(pipePath, os.O_WRONLY, os.ModeNamedPipe)
        defer writer.Close()
        for {
            writer.WriteString("INFO: All systems operational\n")
            writer.WriteString("ERROR: An error occurred\n")
            time.Sleep(1 * time.Second)
        }
    }()

    filterLogs(pipeFile, os.Stdout)                 // ERROR-ları filtrlə
}
```
**Model:** İstənilən process bu pipe-a log yaza bilər; bizim alət real-zamanlı filtrləyir — monitoring/analitik alətlərin əsası.

## Əsas terminlər
- IPC (Inter-Process Communication) — proseslərarası kommunikasiya
- Pipe — yaddaş data kəsiki; unidirectional
- io.Pipe() vs os.Pipe() — in-memory vs SYS_PIPE2
- exec.StdoutPipe()/Stdin — komanda zənciri
- Output() — icra + output tutma qısayolu
- Named pipe / FIFO / Mkfifo — filesystem-də yaşayan pipe
- Sentinel (EOD) — sonlanma işarəsi
- Chunking — böyük datanın hissə-hissə ötürülməsi
- gzip.NewWriter/Reader — pipe üzərində sıxma
- File descriptor exhaustion — fd tükənməsi
- Name squatting — ad ələkeçirmə hücumu
- O_RDWR/O_RDONLY/O_WRONLY — pipe açılış rejimləri

## Praktik nəticə
1. Proseslərarası (xüsusən çoxdilli) axın üçün pipe; tək Go proqramı daxilində isə channel — sərhədi qarışdırma.
2. Manual pipeline: StdoutPipe + Stdin təyinatı + Start/Run/Wait sırası; sadə hallarda Output() qısayolu.
3. Named pipe-larda sentinel mesaj (EOD) + scanner — axışın təmiz dayandırılması.
4. Böyük data üçün chunking; sıxıla bilən data üçün gzip qatı; blocking əməliyyatlara həmişə timeout.
5. Named pipe yarat: random ad + 0600 icazə — squatting və sızma qarşısı.
6. Defer ilə Close — açıq pipe = fd sızmısı.

## Mənbə
Pages: 105-122 (PDF səh. 126-143)
