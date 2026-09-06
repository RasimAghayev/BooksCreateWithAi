# Chapter 5 — Debugging and testing (Technique 22-31)

## Bu chapter nədən bəhs edir?

Logging (log paketi, Logger, network logging, back pressure, UDP vs TCP, syslog), stack trace-lər (debug.PrintStack, runtime.Stack, Caller), unit testing (mock/stub interfeyslərlə, canary testlər, generative testing/quick), benchmarking (testing.B, b.N, RunParallel, -cpu, --race).

## Əsas fikirlər

### 1. Debugger vəziyyəti
Go-da rəsmi debugger YOX (GDB plugin — etibarsız). Alternativlər:
- **Delve** (github.com/derekparker/delve) — tam funksiyalı, goroutine tracing
- **Godebug** (github.com/mailgun/godebug) — kod içi breakpoint
- Amma: log + testlər çox vaxt kifayət edir

### 2. log paketi əsasları
```go
log.Println("This is a regular message.")    // stderr-ə
log.Fatalln("This is a fatal error.")         // stderr + os.Exit(1) — sonrakı sətir İCRA OLMAZ
log.Printf("The value of i is %s", i)         // printf-variant
log.Panicln(...)                              // log + panic
```
- Default hədəf: **Stderr**
- `Fatal*` → `os.Exit(1)` → **defer-lər çağırılmır!** (Ch 4 xatırla)

### TECHNIQUE 22: Logger — arbitrary io.Writer
```go
logfile, _ := os.Create("./log.txt")     // diqqət: hər dəfə üstünə yazır
defer logfile.Close()
logger := log.New(logfile, "example ", log.LstdFlags|log.Lshortfile)
logger.Println("This is a regular message.")
```

**Log mesajının 3 hissəsi:**
```
example 2015/05/12 08:42:51 outfile.go:16: This is a regular message.
│      │                   │                └ Message
│      └───────────────────┴ Generated data
└ Prefix
```

**Flag-lar (bitmask, | ilə birləşir):**
| Vaxt | Yer |
|---|---|
| `Ldate` — tarix | `Llongfile` — tam yol + sətir |
| `Ltime` — vaxt | `Lshortfile` — fayl adı + sətir |
| `Lmicroseconds` — µs dəqiqliyi | (birlikdə işlək deyil!) |
| `LstdFlags` = Ldate \| Ltime | |

Prefix sonunda boşluq qoy — logger boşluq əlavə etmir.

### TECHNIQUE 23: Network logging
**Sınaq serveri:** `nc -lk 1902` (Netcat, TCP listener).

```go
conn, err := net.Dial("tcp", "localhost:1902")
if err != nil {
    panic("Failed to connect to localhost:1902")
}
defer conn.Close()                    // panic belə olsa buffer flush olunur!
f := log.Ldate | log.Lshortfile
logger := log.New(conn, "example ", f)
logger.Println("This is a regular message.")
logger.Panicln("This is a panic.")     // Fatal YOX! (defer işləsin deyə)
```
- `net.Conn` — io.Writer-dir → logger birbaşa şəbəkəyə yaza bilər
- **Panicln vs Fatalln:** Fatal defer-i çağırmır (os.Exit) → bağlantı düzgün bağlanmır → **network logging-də həmişə Panicln**
- Sərbəst timestamp: host vaxtı loglanır → log server gecikəndə belə hadisə xronologiyası qurula bilər

### TECHNIQUE 24: Back pressure
**Problem:** TCP-də hər mesaj üçün ACK gözlənilir. Log server yavaşlayanda client bloklanır = **back pressure**.

**Həll A — UDP:**
```go
timeout := 30 * time.Second
conn, err := net.DialTimeout("udp", "localhost:1902", timeout)
defer conn.Close()
logger := log.New(conn, "example ", log.Ldate|log.Lshortfile)
logger.Println("This is a regular message.")
```
Server: `nc -luk 1902`

| TCP | UDP |
|---|---|
| ACK gözləmə → back pressure | ACK yoxdur → appsə təzyiq YOX |
| Zəmanətli çatdırılma | Mesaj İTƏ BİLƏR |
| Sıralı | Qarışıq sıra mümkün |
| Yavaş | Sürətli |

**Seçim GIF/PNG vs JPEG dilemması kimi:** dəqiqlik vs sürət. Log server həcmi proqnozlaşdırıla bilirsə UDP OK; mesaj itkisi qəbuledilməzdirsə TCP + böyük bufer.

### TECHNIQUE 25: Syslog
**Log səviyyələri (ümumi konsepsiya):** Trace, Debug, Info, Warn, Error, Critical (+Notice, Alert, Emergency).

**2 üsul:**

**a) Go logger-i syslog-a bağla:**
```go
priority := syslog.LOG_LOCAL3 | syslog.LOG_NOTICE
flags := log.Ldate | log.Lshortfile
logger, err := syslog.NewLogger(priority, flags)
logger.Println("This is a test log message.")
// Cəhd bir dəfə: facility + severity (sonra dəyişməz olur!)
```

**b) Birbaşa syslog funksiyaları (daha güclü):**
```go
logger, err := syslog.New(syslog.LOG_LOCAL3, "narwhal")   // facility + prefix
defer logger.Close()
logger.Debug("Debug message.")     // sistem konfiqi siyahıya salmayə bilər!
logger.Notice("Notice message.")
logger.Warning("Warning message.")
logger.Alert("Alert message.")
```
Output: `Jun 30 08:52:06 host narwhal[PID]: Warning message.`
- **Debug görünmədi** — sistem syslog konfiqi filtrləyir! Bu FAYDALIDIR: developer nəyin görünəcəyinə qərar vermir — operator verir.
- **Remote syslog:** `syslog.Dial(network, raddr, priority, tag)` — cluster log aqreqasiyası.

**3rd-party:** Logrus, Glog.

### TECHNIQUE 26: Stack trace
**Sadə (stdout):**
```go
import "runtime/debug"

func bar() {
    debug.PrintStack()     // main→foo→bar zənciri + runtime
}
```

**Buferə (log/göndərmək üçün):**
```go
import "runtime"

buf := make([]byte, 1024)          // ölçünü əvvəlcədən təyin et — dəqiq hesablama yolu yoxdur
runtime.Stack(buf, false)          // false = yalnız cari goroutine; TRUE = BÜTÜN goroutine-lər!
fmt.Printf("Trace:\n %s\n", buf)
```
- `true` → bütün goroutine stack-ləri (concurrency debug üçün qiymətli, amma çox böyük output)
- Daha incə nəzarət: `runtime.Caller` / `runtime.Callers`

### 3. Unit test əsasları
```go
// hello.go:
package hello
func Hello() string { return "hello" }

// hello_test.go:
package hello                     // EYNİ PAKET — unexported-lara da çıxış!
import "testing"
func TestHello(t *testing.T) {
    if v := Hello(); v != "hello" {
        t.Errorf("Expected 'hello', but got '%s'", v)
    }
}
```

**2 tez-tez səhv:**
1. Test fayllarını ayrı qovluğa qoymaq → **EYNİ qovluqda olmalıdır**
2. Fərqli paketdə (hello_test) → **EYNİ paket** — private kod da test olunsun

**testing.T əsas metodları:**
- `t.Error(args)` / `t.Errorf(format, ...)` — logla + FAIL + davam
- `t.Fatal(args)` / `t.Fatalf(...)` — logla + FAIL + **DAYANDIR** (sonrakı testlər mənasızdırsa)

### TECHNIQUE 27: Mock/stub — interfeyslə
**Problem:** xarici kitabxana (Message.Send) testdə çağırılmasın, amma çağırıldığı YOXLANILSIN.

```go
// Production kod — konkret tip YOX, İNTERFEYS qəbul et:
type Messager interface {
    Send(email, subject string, body []byte) error
}

func Alert(m Messager, problem []byte) error {
    return m.Send("noc@example.com", "Critical Error", problem)
}

// Test — Mock:
type MockMessage struct {
    email, subject string
    body           []byte
}

func (m *MockMessage) Send(email, subject string, body []byte) error {
    m.email = email            // göndərmə — YADDA SAXLA
    m.subject = subject
    m.body = body
    return nil
}

func TestAlert(t *testing.T) {
    msgr := new(MockMessage)
    body := []byte("Critical Error")
    Alert(msgr, body)
    if msgr.subject != "Critical Error" {           // mock-un saxladığını YOXLA
        t.Errorf("Expected 'Critical Error', Got '%s'", msgr.subject)
    }
}
```
+ Bonus: implementasiya dəyişmək asanlaşır (modular proqramlaşma).

### TECHNIQUE 28: Canary testlər — interfeys yoxlaması
**Problem:** interfeys imzası səhv → type assertion runtime-da partlayır; compile-time yox.

**Nümunə — səhv Write imzası:**
```go
func (m *MyWriter) Write([]byte) error { ... }    // io.Writer DEYİL! (int qaytarmır)
```
Runtime-da `m["w"].(io.Writer)` → panic. Amma compile OK.

**Canary test (compiler-ə yoxlama etdir):**
```go
func TestWriter(t *testing.T) {
    var _ io.Writer = &MyWriter{}     // təyin etməsən də compiler YOXLAYIR
}
```
Compile xətası dəqiq deyir:
```
*MyWriter does not implement io.Writer (wrong type for Write method)
    have Write([]byte) error
    want Write([]byte) (int, error)
```
**4 halda faydalı:** (1) export olunan tip xarici interfeysi implement edir; (2) xarici tipi təsvir edən interfeys yaratmışan; (3) xarici interfeys dəyişibsə; (4) type assertion istifadə olunursa.

### 4. Generative testing — testing/quick
**Nədir:** avtomatik test data generasiyası — bizim test-data bias-ımızı aradan qaldırır.

**Pad funksiyası + bug:**
```go
func Pad(s string, max uint) string {
    ln := uint(len(s))
    if ln > max {
        return s[:max-1]     // BUG! s[:max] olmalıdır — length veririk, index yox
    }
    s += strings.Repeat(" ", int(max-ln))
    return s
}

// Sadə test — KEÇİR (max>len halını əhatə etmir!):
if r := Pad("test", 6); len(r) != 6 { ... }

// Generative test — random string-lərlə BUG-I TAPIR:
func TestPadGenerative(t *testing.T) {
    fn := func(s string, max uint8) bool {       // parametrlərdən data tipi çıxarılır
        p := Pad(s, uint(max))
        return len(p) == int(max)
    }
    if err := quick.Check(fn, &quick.Config{MaxCount: 200}); err != nil {
        t.Error(err)
    }
}
// FAIL: failed on input "\U000305ea...", 0x20 — uzun string truncation bug!
```
- `quick.Check(fn, config)` — funksiyanın parametrlərini introspekt edir, uyğun random data yaradır
- uint8 → qısa stringlər; uint16 → uzun (seçim səni)
- **Random seed:** hər run EYNİ data (seed yoxdur) — repeatable; fərqli istəyirsənsə Config-ə seed ötür (qarşılıq: tapılan bug reproduksiya oluna bilməz!)
- Hətta random STRUCT instansiyaları yarada bilər

### TECHNIQUE 29: Benchmarking — testing.B
```go
func BenchmarkTemplates(b *testing.B) {          // Test→Benchmark prefiksi; T→B
    tpl := "Hello {{.Name}}"
    data := &map[string]string{"Name": "World"}
    var buf bytes.Buffer
    for i := 0; i < b.N; i++ {                  // HƏMİŞƏ b.N-ə qədər loop!
        t, _ := template.New("test").Parse(tpl)
        t.Execute(&buf, data)
        buf.Reset()
    }
}
```
```bash
$ go test -bench .
BenchmarkTemplates  100000  10102 ns/op
# b.N: 1 → 100 → 10000 → 100000 — framework özü kalibrləşir!
```

**Optimizasiya müqayisəsi (compile cache):**
```go
func BenchmarkCompiledTemplates(b *testing.B) {
    t, _ := template.New("test").Parse(tpl)     // LOOP-DAN XARIÇDA — 1 dəfə compile!
    ...
    for i := 0; i < b.N; i++ {
        t.Execute(&buf, data)
    }
}
// 10167 ns/op vs 1318 ns/op — 10x FƏRQ! (compile bir dəfə kənarda)
```

### TECHNIQUE 30: Parallel benchmark — RunParallel
```go
func BenchmarkParallelTemplates(b *testing.B) {
    tpl := "Hello {{.Name}}"
    t, _ := template.New("test").Parse(tpl)
    data := &map[string]string{"Name": "World"}
    b.RunParallel(func(pb *testing.PB) {        // çoxlu goroutine-də icra
        var buf bytes.Buffer                    // HƏR GOROUTINE ÖZ BUFERİ!
        for pb.Next() {                         // davam etmək lazımdır?
            t.Execute(&buf, data)
            buf.Reset()
        }
    })
}
```

**-cpu flag:**
```bash
$ go test -bench . -cpu=1,2,4
BenchmarkTemplates            100000    10019 ns/op   # 1 CPU
BenchmarkParallelTemplates    1000000    1249 ns/op    # 1 CPU — fayda YOX
BenchmarkParallelTemplates-2  2000000     784 ns/op    # 2 CPU — ~1.6x qazanc
BenchmarkParallelTemplates-4  2000000     829 ns/op    # 4 CPU — lock overhead!
```
- `-N` şəkilçisi = N prosessor
- Serial kod çox CPU-dan qazanmır; parallel kod 2-4 CPU-da qazanır (4-də lock overhead görünür)

### TECHNIQUE 31: --race ilə race detection
**Səhv optimizasiya — paylaşılan bufer:**
```go
func BenchmarkParallelOops(b *testing.B) {
    ...
    var buf bytes.Buffer                     // GOROUTINE-LER ARASI PAYLAŞILIR!
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            t.Execute(&buf, data)            // RACE!
            buf.Reset()
        }
    })
}
```
**--race olmadan:** qeyri-müəyyən panic (`slice bounds out of range`) — səbəb QARANLIQ.
**--race ilə:**
```bash
$ go test -bench Oops -race -cpu=1,2,4
WARNING: DATA RACE
Write by goroutine 20:
  bytes.(*Buffer).Write()     ← DƏQİQ SƏBƏB!
```
→ həll: geriqayıd (goroutine başına bufer) yaxud sync.Mutex.

**Qayda:** `--race` go run / go test / benchmark hamısında — development standartı.

## Alət qısa xülasəsi
| Alət | Komanda | Məqsəd |
|---|---|---|
| Logger | log.New(writer, prefix, flags) | fayl/şəbəkə/syslog |
| Syslog | syslog.New(facility, tag) | sistem səviyyələri |
| Stack | debug.PrintStack / runtime.Stack(buf, all) | trace |
| Test | go test | unit |
| Canary | var _ io.Writer = &MyWriter{} | compile-time interfeys |
| Generative | quick.Check(fn, config) | random data |
| Benchmark | go test -bench . | ns/op |
| Parallel | b.RunParallel + -cpu=1,2,4 | multi-core |
| Race | --race | data race |

## Əsas terminlər
- log.Logger / log.New / Fatal* (defer-i keçir!) / Panicln
- Ldate / Ltime / Lmicroseconds / LstdFlags / Llongfile / Lshortfile (bitmask)
- Back Pressure (ACK gözləməsi)
- UDP vs TCP logging (itki vs zəmanət)
- Syslog: facility / severity / priority / LOG_LOCAL3
- Log Səviyyələri: Trace/Debug/Info/Warn/Error/Critical
- syslog.Dial (remote)
- debug.PrintStack / runtime.Stack (all goroutines) / runtime.Caller
- testing.T: Error / Fatal / Errorf / Fatalf
- Eyni qovluq + eyni paket qaydası
- Mock / Stub (interfeys + yadda saxlayan implementasiya)
- Canary Test (compile-time type yoxlaması; "kömürdəki kanarya")
- Generative Testing / testing/quick.Check / Config{MaxCount}
- testing.B / b.N (avtomatik kalibrləmə)
- Benchmark prefiksi
- Compile-outside-loop optimizasiyası (10x dərsi)
- b.RunParallel / pb.Next / *testing.PB
- -cpu=1,2,4 (-N şəkilçisi)
- --race (go run / go test / bench)

## Praktik nəticə
- Fatal* defer-i çağırmır — network logging-də Panicln istifadə et; ümumiyyətlə production server-də Fatal-ə yaxınlaşma.
- Logger-i io.Writer-ə bağla → fayldan şəbəkəyə keçiş 1 sətir dəyişikliyi.
- UDP logging: mesaj itkisi OK + sürət lazımdırsa; zəmanət lazımdırsa TCP + bufer.
- Syslog severity-lərindən istifadə et — operator konfiqə ilə filtrləmə qərarını versin.
- Testlər eyni paketdə — private funksiyalar da test olunsun; assertion kitabxanası məcburi deyil.
- Mock yalnız 3 sətir: interfeys yarat → koda interfeys qəbul etdir → mock struct yadda saxlasın.
- Export interfeysli tiplər üçün canary test yaz: `var _ io.Writer = &MyWriter{}` — runtime panic-i compile xətasına çevirir.
- quick.Check ilə truncation/kənar hal bug-larını tap — insan test datası bias-lıdır.
- Benchmark-lı müqayisə: parse/compile kimi ağır işləri loop-dan kənara çıxar (10x sübut olundu).
- RunParallel + -cpu multi-core performansı ölçür; --race olmadan parallel bug qeyri-müəyyən görünür.

## Mənbə
Pages: 136-166 (PDF), book pages 113-143
