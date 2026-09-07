# Chapter 12 — Cloud-ready applications and communications (Bulut-hazır tətbiqlər və rabitə)

## Bu chapter nədən bəhs edir?

Mikroservis arxitekturası (transkodinq nümunəsi, Kafka), high availability, servis rabitə
sürəti (keep-alive bağlantı reuse, response body bağlama, codec/sonic JSON), runtime
aşkarlama (os paketi, exec.LookPath), cross-compile (GOOS/GOARCH, gox, build tag-lər,
fayl suffiksləri) və runtime monitorinqi (runtime.MemStats, NumGoroutine).

## Əsas fikirlər

### 1. Mikroservis Arxitekturası
**Transkodinq nümunəsi:** UI → API server → fayl store + Kafka mesaj queue → transcoder
→ notification → fayl geri. UNIX fəlsəfəsi: "Do one thing and do it well" — hər servis
öz işi. Kafka istehlakçı mikroservisi:
```go
conn, err := kafka.DialLeader(context.Background(), "tcp", kafkaHost, topic, 0)
conn.SetReadDeadline(time.Now().Add(30 * time.Second))
batch := conn.ReadBatch(10e3, 1e6)
for { n, err := batch.Read(message); if err != nil { break } }
```
Servislər: fərqli dillərdə ola bilər; fərqli scale (transcoder yükdən asılı; API fərqli);
"heç vaxt offline olmaz" gözləntisi → hər servis özünə uyğun HA metodu.

### 2. Servislərarası Rabitə — REST Sürətləndirmə
**Bağlantı reuse:** Hər request yeni bağlantı = TLS negotiate + TCP slow-start ramp-up.
Persistent bağlantı = bir dəfə xərc, sonrakılar sürətlidir.

- HTTP keep-alive ≠ TCP keep-alive: birincisi HTTP protokol/server-işi, ikincisi OS/TCP.
  `DisableKeepAlives` hər ikisini öldürür.
- DefaultClient/DefaultTransport keep-alive AÇIQDIR; antipattern — custom Transport-da
  Dial konfiqurasiyasız buraxmaq:

**Kitabdan kod nümunəsi (doğru custom transport):**
```go
tr := &http.Transport{
    TLSClientConfig:    &tls.Config{RootCAs: pool},
    DisableCompression: true,
    Dial: (&net.Dialer{
        Timeout:   30 * time.Second,
        KeepAlive: 30 * time.Second,
    }).Dial,                       // DefaultTransport-un eynisi
}
client := &http.Client{Transport: tr}
```

**Response body bağlama tələsi:** `defer r.Body.Close()` funksiya sonuna qədər gecikir →
ikinci request YENİ bağlantı açır. Çoxsaylı serial requestdə body-ni istifadə edən kimi
BAĞLA:
```go
o, _ := ioutil.ReadAll(r.Body)
r.Body.Close()                    // dərhal — növbəti request eyni bağlantını reuse edir
r2, _ := http.Get("http://example.com/foo")
```
Go 1.6+ HTTP/2 transparent; pipelining (paralel request-lər) HTTP/2 ilə gəlir.

### 3. Sürətli JSON
encoding/json = reflection hər dəfə. Alternativlər:
- **ugorji/go/codec (codecgen):** `//go:generate codecgen -o user_generated.go user.go`
  + `codec:"name"` tag-ləri → generasiya olunmuş CodecEncodeSelf/DecodeSelf metodları
  (reflection-sız); `go generate ./...` ilə. JsonHandle + NewEncoderBytes/NewDecoderBytes.
- **bytedance/sonic:** JIT əsaslı drop-in əvəz (`sonic.ConfigDefault.NewDecoder`);
  böyük JSON həcmlərində görünməz bottleneck aradan qalxır; M-seriya chipində yavaş build.

### 4. Runtime Aşkarlama (host barədə)
- `os.Hostname()`, `os.Getpid()`, `os.Getwd()`, `os.PathSeparator/PathListSeparator`
- **IP ünvanı:** hostname-i al → `net.LookupHost(name)` — bütün interfeys ünvanlarını
  yoxlamağın təmiz yolu
- **Asılılıq yoxlaması:**
```go
func checkDep(name string) error {
    if _, err := exec.LookPath(name); err != nil {
        return fmt.Errorf("Could not find '%s' in PATH: %s", name, err)
    }
    return nil
}
```
Bulutda minimal distributivlər (demək olar alətsiz) — çağırılmazdan ƏVVƏL yoxla + logla.

### 5. Cross-Compile
```bash
GOOS=windows GOARCH=386 go build
# çoxsaylı paralel: gox
gox -os="linux darwin windows" -arch="amd64 386" -output="dist/{{.OS}}-{{.Arch}}/{{.Dir}}" .
```
- GOARCH: amd64/386/arm; GOOS: windows/linux/darwin/freebsd
- cgo cross-compile problemlidir — hər platformada test et
- **path/filepath:** Separator/ListSeparator/ToSlash/Split/Join — hardcoded `/` YOX
- **Build tag:** `// +build !windows` — faylı Windows-da yığma; `!linux,!darwin` çoxlu
- **Fayl suffiksləri:** `foo_windows.go`, `foo_386.go` — avtomatik seçim

### 6. Runtime Monitorinqi
**Sidecar goroutine:**
```go
func monitorRuntime() {
    log.Println("Number of CPUs:", runtime.NumCPU())
    m := &runtime.MemStats{}
    for {
        log.Println("Number of goroutines", runtime.NumGoroutine())   // leak aşkarlama
        runtime.ReadMemStats(m)
        log.Println("Allocated memory", m.Alloc)
        time.Sleep(10 * time.Second)
    }
}
func main() { go monitorRuntime(); ... }
```
Mövcud data: GC statistikası (son pass, növbəti trigger heap ölçüsü, müddət), heap
statistikası, goroutine/processor/cgo sayı. **Diqqət:** ReadMemStats runtime-ı anlıq
DAYANDIRIR — tez-tez çağırma; debug rejimə saxla. Goroutine sayı artımı = kitabxana bug-u
siqnalı (milyonlarla goroutine). Log əvəzinə bulud monitoring API-ə göndər.

## Əsas terminlələr
- Mikroservis / message queue — kiçik məsul servis / Kafka-növbəsi
- TCP slow-start — congestion-control ramp-up (persistent bağlantıda 1 dəfə)
- HTTP vs TCP keep-alive — protokol vs OS səviyyəli canlılıq
- codecgen / go:generate — JSON codegen aləti / generasiya komandası
- Sonic JIT — runtime JSON kompilyasiyası
- exec.LookPath — PATH-də asılılıq axtarışı
- GOOS/GOARCH/gox — cross-compile dəyişənləri / paralel multi-build
- Build tag / fayl suffiksi — platforma-filtrli kompilyasiya
- runtime.MemStats / NumGoroutine — runtime sağlamlıq göstəriciləri

## Praktik nətidə

Bulut qərarları: (1) hər mikroservis öz dilində, öz scale-unitsində; rabitə queue ilə
lock-in azaldır; (2) REST sürəti — custom Transport-da KeepAlive dial QUR, serial
request-lərdə body-ni dərhal bağla; (3) JSON bottleneck-i profil göstərəndə codecgen
(codegen) və ya sonic (JIT); (4) host məlumatını assume YOX — runtime-da aşkarla;
asılılıqları LookPath ilə yoxla; (5) cross-compile GOOS/GOARCH + gox; filepath paketi
+ build tag-lər platforma fərqlərini idarə edir; (6) sidecar monitoring goroutine —
goroutine/memory artımı sızntı/bug siqnalı; ReadMemStats-ı seyrək çağır.

## Mənbə
Pages: 293-315 (PDF 314-336)
