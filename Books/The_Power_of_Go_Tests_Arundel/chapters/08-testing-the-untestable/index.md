# Chapter 8 — Testing the untestable (Test Edilməzın Testi)

## Bu fəsil nədən bəhs edir?

"Untestable" görünenlərin testi: walking skeleton (Freeman & Pryce), Hello
Production (Pete Hodgson), test-in-prod (Charity Majors), Kent Beck "make the
problem easy", API dizaynı (funksional options — LoadTester nümunəsi, Alan Kay
sitatı), external vs internal testlər (package_X_test.go, _internal_test.go —
Mat Ryer), sinxron API üstünlüyü (Dave Cheney ListDirectory 3 variantı, callback
/fs.WalkDir), async test pattern-ləri: randomLocalAddr (port 0 fəndi), eager
test uğursuzluğu, sleep flakiness, wait-for-success + time.Timer + select,
waitForServer helper, t.Fatal yalnız test goroutine-ində; concurrency safety:
global map data race, smoke test, -race detector, Store+mutex refaktoru,
Rich Hickey; context.WithTimeout ilə uzun işlərin testi (RunBatchJob,
ctx.Done(), DeadlineExceeded); user interaction: Greet fmt.Print→os.Stdout
hardwire, Fprint/Fscanln more-general API-lər, io.Reader/io.Writer param-lər,
strings.NewReader, "Mary Jo" boşluq bug-ı (Fscan space-delimited), bufio.Scanner;
CLI testi: main→timer.Main(os.Args), ParseArgs/Sleep ayrımı, Timer struct,
NewTimerFromArgs.

## Əsas fikirlər

### 1. Walking Skeleton — İlk Addım
"Thinnest possible slice of real functionality that we can build, deploy and
test end-to-end... functionality so simple it's obvious and uninteresting."
(Freeman & Pryce, GOOS)

Software engineering komponentlərin İNTERNINDƏN çox, komponentlər arası
QOŞULUŞLAR haqqındadır — bağlantıları dizayn etmək/debug etmək komponentin
özündən uzun çəkir. Dummy komponentləri "düzgün görünən şəkildə" qoş, İŞLƏD,
iterasiya et — funksiya az olduğu üçün çox sürətli iterate olunur.

Sonra: bütün vacib davranışları testli package-ə köçür. **Heç vaxt runnable
proqramdan uzaq düşmə** — dəyişikliklər kiçik və inkremental.

### 2. Hello, Production — Ən Tez Deploy
"Get the simplest version into prod as soon as you can." (Pete Hodgson)

İzolyasiyada test yalnız izolyasiyadakı davranışı göstərir. Prod = users + code +
environment + infrastructure + point in time — emergent, deterministik test
OLMAZ. "I test in prod" — false dichotomy YOX: məsul komandalar HƏR İKİNİ edir.
(Charity Majors)

Full pipeline işləyir → istənilən an ship edə bilərsən. Ən pis vaxt
deadline-dan ƏVVƏL böyük dizayn issue tapmaqdır — "half-finished as soon as
possible" prinsipi.

### 3. İlk Test = Ən Çətin → Sadələşdir
"Make the problem easy, then solve the easy problem (warning: making the
problem easy may be hard)." (Kent Beck)

**"Buy product" endpoint:** form data + payment + DB order = horrendous →
ÖDƏNİŞİ KƏS, order-i AT: handler form-u İGNORE edib "hello, world" qaytarsın.
Faydalı? YOX. Test olunurmu? BƏLİ. İnkremental davam?

**Böyük problemlər:** filli bir diş vur ("if I had this component and that
component, I could...") — "Be a CRM" → "Store customer's name and address" +
"Generate invoices" → təsəvvür ediləbilən vahidlərə qədər parçala. Bu, kağız
oynatmaq deyil — hər parçalama DİZAYNI dəqiqləşdirir.

"NEVER build large apps... assemble testable, bite-sized pieces." (Justin Meyer)

### 4. API Dizaynı — İstifadəçi Görüşündən
Funksiya yazmırsan — istifadəçilərin proqramlarını elegan ifadə edəcəyi DİL
yaradırsan; funksiyalar = o dilin sözləri. (Dijkstra: abstraksiya qeyri-müəyyənlik
DEYİL, yeni dəqiqlik səviyyəsidir.)

**LoadTester nümunəsi — 6-parametrli konstruktor PROBLEM:**
```go
load.NewLoadTester("https://example.com", os.Stdout, os.Stderr,
    http.DefaultClient, 20, "loadtest")   // uzun, qoxulu, variantlı
```
**Funksional options həlli:**
```go
got := load.NewLoadTester("https://example.com")     // sadə hal — qısa

got := load.NewLoadTester("https://example.com",     // customise lazımdırsa
    load.WithOutput(buf),
    load.WithErrOutput(buf),
    load.WithHTTPClient(&http.Client{Timeout: time.Second}),
    load.WithConcurrentRequests(20),
    load.WithUserAgent("loadtest"),
)
```
"Simple things should be simple, complex things should be possible." (Alan Kay)
Boilerplate az-çoxdur, amma məcburi parametrləri YOX etmək dəyər.

### 5. External vs Internal Testlər
**External test** (başqa package): yalnız API görürsün → black box →
implementasiya dəyişsə testlər QIRILMIR. "Your test files should always be in a
different package." (Michael Sorens)

**Unexported funksiyalar:** lowercase = yalnız package daxili. Başqa package-dən
çağırılmır → DİREKT test lazım deyil — onları İSTİFADƏ EDƏN exported
funksiyalar vasitəsilə İNDİRİKT test olunur. "Tests should go through the front
door like everybody else." (Freeman & Pryce)

**İstisna — internal test:** gizli davranış (keçən fəsildə CACHE) ancaq daxildən
yoxlanılır → `*_internal_test.go` suffiksi ilə AYRI fayl (Mat Ryer). Internal
testlər TƏBIİ olaraq daha brittledir — implementasiya detailinə bağlıdırlar.

### 6. Concurrency — Sinxron API Üstünlüyü
**Kanarı rule:** concurrent API test etmək çətindirsə — bu, DİZAYN issuesidir.
"The tests are a canary in a coal mine." (Kent Beck)

Channel API istifadəçiyə də çətindir: sahib kim? lifecycle? kim close edir?
buffered? Sinxron API = asan anlaşılır, asan test.

**Dave Cheney — ListDirectory 3 variantı:**
```go
func ListDirectory(dir string) ([]string, error)      // 1. SINXRON: hamısı bir yerdə
func ListDirectory(dir string) chan string            // 2. ASINXRON: channel + close
func ListDirectory(dir string, fn func(string))       // 3. CALLBACK: istəyə görə
```
Variant 2-nin istifadəçisi nə edir? `for entry := range entriesChan` → slice
toplayır → ...bu, məhz variant 1-in qaytardığıdır! Parallellik YALAN oldu:
istifadəçi blocked qalır. Variant 3 (callback) — fs.WalkDir pattern:
```go
fs.WalkDir(fsys, ".", func(p string, d fs.DirEntry, err error) error {
    if filepath.Ext(p) == ".go" { count++ }
    return nil
})
```
Concurrency QƏRARI istifadəçiyə buraxılır — istəsə öz goroutine-ində çağırır.
**Rule: concurrency-ni API-da DEYİL, caller-da saxla.**

### 7. Async API Testi — 4 İnkişaf Mərhələsi
**Problem:** `serve.ListenAsync(addr)` dərhal qayıdır; server goroutine-da
açılır.

**Port problemi → randomLocalAddr helper:**
```go
func randomLocalAddr(t *testing.T) string {
    t.Helper()
    l, err := net.Listen("tcp", "localhost:0")   // PORT 0 = kernel RANDOM FREE port seçir!
    if err != nil { t.Fatal(err) }
    defer l.Close()
    return l.Addr().String()                     // "127.0.0.1:60606"
}
```
Hər test müstəqil, paralel port alır — toqquşma YOX, "already in use" YOX.

**Mərhələ 1 — EAGER (uğursuz):** Dial dərhal → "connection refused" —
server hələ başlamayıb. Concurrency ≠ parallelism; goroutine başlatmaq sürətlidir,
ANCAQ instant DEYİL.

**Mərhələ 2 — Fixed sleep (FLAKY):**
```go
time.Sleep(10 * time.Millisecond)
```
`go test -count 10` → bəzən FAIL: race — bəzən server 10ms-də çatmır. Sleep
qısa → flake; uzun → suite yavaş. DOĞRU INTERVAL YOXDUR. "Flickering test"
(GBOS) — böyüyən suite-də gərginlik artır.

**Mərhələ 3 — Wait for success:**
```go
_, err := net.Dial("tcp", addr)
for err != nil {
    t.Log("retrying")
    time.Sleep(time.Millisecond)
    _, err = net.Dial("tcp", addr)
}
```
İlk cəhd UĞURLU ola bilər (loop-a girmir); uğursuzsa ~1ms-lik retry-lar —
serverə lazım olan qədər gözlə, amma o qədər çox YOX. "Every tested activity
must affect the system so that its observable state becomes different."
(Freeman & Pryce — observable effect YOXDURSA, wait üçün nə YOXDUR!)

**Mərhələ 4 — Timeout + helper:**
```go
func waitForServer(t *testing.T, addr string) {
    t.Helper()
    timeout := time.NewTimer(100 * time.Millisecond)
    _, err := net.Dial("tcp", addr)
    for err != nil {
        select {
        case <-timeout.C:
            t.Fatal("timed out")
        default:
            t.Log("retrying: ", err)
            time.Sleep(time.Millisecond)
            _, err = net.Dial("tcp", addr)
        }
    }
}
```
Server heç açılmasa 100ms sonra t.Fatal — "never starts" halı sürətli tutulur.
Test:
```go
addr := randomLocalAddr(t)
serve.ListenAsync(addr)
waitForServer(t, addr)
// test whatever we're really testing
```

**Vacib məhdudiyyət — t.Fatal yalnız test-in ÖZ goroutine-ində:**
```go
go func() {
    time.Sleep(100 * time.Millisecond)
    t.Fatal("timeout")     // NO EFFECT! İŞLƏMİR
}()
```
Buna görə time.AfterFunc da YOX — select + timer channel DÜZ həlldir.

### 8. Concurrency Safety — Go-da Həmişə
API özü concurrency istifadə etməsə belə — Bu Go-dur, istifadəçi PARALEL
çağıracaq. **Global mutable state klassikası:**
```go
var data = map[string]string{}   // no! data race!
func Set(k, v string) { data[k] = v }
func Get(k string) string { return data[k] }
```
Smoke test: bir goroutine 1000 Set, test goroutine 1000 Get + runtime.Gosched()
→ "fatal error: concurrent map read and map write". (Gosched = scheduler-ə
nəfəs vermək — race şansını artırır.)

**Amma şans işidir** — fatal error heç törəməyə bilər. **Gizli silah: race
detector:**
```bash
go test -race
# WARNING: DATA RACE — Write at ... by goroutine 8 / Previous read ... by goroutine 7
# --- FAIL: race detected during execution of test
```
Mapassign/mapaccess stack trace-ləri ilə dəqiq göstərir: Set yazır, Get oxuyur,
sinxronizasiya YOX.

**Düzəliş — Store tipi + mutex, global dəyişən ləğvi:**
```go
type Store struct {
    m    *sync.Mutex
    data map[string]string
}
func NewStore() *Store { ... }
func (s *Store) Set(k, v string) { s.m.Lock(); defer s.m.Unlock(); s.data[k] = v }
func (s *Store) Get(k string) string { s.m.Lock(); defer s.m.Unlock(); return s.data[k] }
```
BONUS: bir deyil, ÇOXLU store mümkün oldu. `go test -race` → PASS.

**Race detector hər race-i TUTA BİLMƏZ** — yalnız icra olunan path-ləri.
"Tests are not a substitute for thinking." (Rich Hickey: "Who drives their car
around banging into the guard rails?") Data race-lər ən çox prod-da, PIK
yükdə partlayır — ən istəmədiyin anda.

### 9. Long-Running Tasks — context ilə Test
context = standart ləğv/timeout mexanizmi. Test üçün BUILT-IN timeout var:
```go
func TestRunBatchJob(t *testing.T) {
    t.Parallel()
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
    defer cancel()
    go func() {
        batch.RunBatchJob(ctx)
        cancel()                 // job vaxtında bitdi → öz ləğv et
    }()
    <-ctx.Done()                 // close-only channel: bağlama = mesaj
    if errors.Is(ctx.Err(), context.DeadlineExceeded) {
        t.Fatal("timed out")
    }
}
```
- defer cancel — test hər hansı şəkildə çıxsa resurs boşa axmasın
- `<-ctx.Done()` — ya job bitdi (goroutine cancel çağırdı), ya 10ms doludu
- `ctx.Err()`: context.Cancelled = NÖRMAL (keçdi) / DeadlineExceeded = TIMEOUT
  (FAIL)
- Close-only channel idiomu: heç nə GÖNDƏRİLMİR — bağlanmanın ÖZÜ mesajdır

### 10. User Interaction — Greet Funksiyası
**Ilkin (test olunmaz) versiya:**
```go
func Greet() {
    fmt.Print("Your name? ")
    var name string
    fmt.Scanln(&name)
    fmt.Printf("Hello, %s!\n", name)
}
```
Test çağırır → "Your name? Hello, !" çap edir, PASS — useless + terminali
kirliyir. Nəticə yoxdur, fail yolu yoxdur.

**Sual: "What are we really testing here?"** → Davranış: prompt çap + input oxu
+ greeting çap. Yəni Greet 2 qapıdan işləyir: **io.Writer-a yazır** + **io.Reader-dan
oxuyur**. Bu anlayış testin açarıdır.

**fmt.Print-in siri:** `func Print(a ...any) { return Fprint(os.Stdout, a...) }`
— convenience wrapper, os.Stdout HARDWIRE. Daha ÜMUMİ API = Fprint (hər hansı
writer). Eyni: Scanln = Fscanln(os.Stdin, ...).

**os.Stdout-u əvəz etmək PİS fikir:** (1) paralel testlərdə data race; (2)
exclusive deyilik — Go özü də (test nəticələri) istifadə edir.

**Refaktor — parametrlər kimi axınlar:**
```go
func Greet(in io.Reader, out io.Writer) {
    fmt.Fprint(out, "Your name? ")
    var name string
    fmt.Fscanln(in, &name)
    fmt.Fprintf(out, "Hello, %s!\n", name)
}
```
**Test:**
```go
buf := new(bytes.Buffer)                       // io.Writer + sonra oxunur
input := strings.NewReader("fakename\n")       // io.Reader, "istifadəçi yazıb"
greet.Greet(input, buf)
want := "Your name? Hello, fakename!\n"
got := buf.String()
if want != got { t.Error(cmp.Diff(want, got)) }
```
main: `greet.Greet(os.Stdin, os.Stdout)` — real axınlar yalnız kənarda.

**TDD dərsi:** test ƏVVƏL olsaydı io.Writer parametri ÖZÜ doğardı — "We simply
can't write untestable functions when the test comes first." Legacy üçün: sadə
refaktor kifayət edir.

**"Mary Jo" bug-ı — təxəyyülü test data:**
```go
input := strings.NewReader("Mary Jo\n")
want := "Your name? Hello, Mary Jo!\n"
got  := "Your name? Hello, Mary!\n"     // FAIL!
```
fmt.Fscan sənədləri: "space-separated values" — boşluq = DELIMITER, "Jo" İTKİR.
Boşluqlu adı olan hər kəs üçün bug! Daha təxəyyülü input / müxtəlif komanda /
fuzz — bug-un ovu.

**Düzəliş — sətir oxumaq üçün bufio.Scanner:**
```go
func Greet(in io.Reader, out io.Writer) {
    fmt.Fprint(out, "Your name? ")
    scanner := bufio.NewScanner(in)
    if !scanner.Scan() { return }        // input yoxdursa — heç nə etmə
    fmt.Fprintf(out, "Hello, %s!\n", scanner.Text())
}
```
Scan() false halı: error qaytarmaq İMKANSIZ (imza dəyişməsin), default ad DA
məntiqli deyil → ən yaxşısı sadəcə return.

### 11. CLI Testi — main-dən çıxarış
main test OLUNMUR (binary compile+run tələb edir) — amma main-da çox şey xəta
edə bilər: heç nə etməmək, yanlış çağırış, arqumentləri ignore/yanlış oxumaq,
missing args-da panic.

**Həll 1 — main yalnız ötürür:**
```go
func main() { timer.Main(os.Args) }
```
Test: `timer.Main([]string{"program", "arg1", "arg2"})` — istənilən args.

**Həll 2 — davranışı 2 funksiyaya BÖL (timer 10s):**
```go
func main() {
    interval := timer.ParseArgs(os.Args)   // ARĞUMENT ŞƏRHİ — pure
    timer.Sleep(interval)                   // İCRA — ayrıca test olunur
}
```
ParseArgs pure funksiyadır → sadə test:
```go
want := 10 * time.Millisecond
got := timer.ParseArgs([]string{"timer", "10ms"})
```

**Həll 3 — konfiq çox olanda struct:**
```go
type Timer struct { Interval time.Duration }
func (t Timer) Sleep() { time.Sleep(t.Interval) }
func NewTimerFromArgs(args []string) Timer { ... }

func TestNewTimerFromArgs(t *testing.T) {
    want := timer.Timer{Interval: 10 * time.Second}
    got := timer.NewTimerFromArgs([]string{"program", "10s"})
    if !cmp.Equal(want, got) { t.Error(cmp.Diff(want, got)) }
}
```
CLI-ni proqramın qalanından DECOUPLE et, hər ikisini ayrıca test et. Növbəti
fəsil: bəzən binary-ni birbaşa İŞLƏT MƏK lazımdır — test scripts.

## Əsas terminlələr
- Walking Skeleton — uçtan-uca işləyən ən nazik dilim (GOOS)
- Hello, Production — ən sadə versiyanın dərhal prod-a çıxarılması
- Test in Prod — izolyasiya testi + prod testi — hər ikisi (Majors)
- Make the Problem Easy — Kent Beck sadələşdirmə qanunu
- Functional Options — WithX(...) variadic konfiq (NewLoadTester)
- External Test — başqa package-dən black-box test (package_x_test)
- Internal Test — _internal_test.go, hər şeyə çıxışı var, brittledir
- Front Door — testlər də public API-dan girir
- Close-only Channel — bağlanma özü mesajdır (ctx.Done)
- Port 0 Trick — kernel RANDOM boş port təyin edir (randomLocalAddr)
- Wait for Success — retry loop + timeout (flaky sleep əvəzinə)
- Flickering Test — ara-sıra fail olan, timeout-a bağlı test
- Observable State — async test yalnız GÖRÜNƏN təsirə görə sinxronlaşır
- Race Detector — go test -race; icra olunan path-lərdəki race-lər
- Smoke Test — paralel çağırışla race yaratmaq cəhdi (şans işidir)
- context.WithTimeout / DeadlineExceeded / ctx.Done / ctx.Err
- Fprint/Fscanln — Print/Scanln-in umumi (writer/reader qəbul edən) tərəfi
- strings.NewReader / bytes.Buffer — sahte stdin/stdout
- bufio.Scanner — sətir-sətir oxu (Fscan boşluq delimiterindən fərqli)
- timer.Main(os.Args) / ParseArgs / NewTimerFromArgs — main-i test olunana çevir

## Praktik nəticə
(1) Yeni layihə: walking skeleton — dummy qoşulmuş uçtan-uca dilim, dərhal
deploy, sonra davranışları testli package-ə köçür. (2) Həmişə runnable yaxın
ol — böyük dizayn issue-ni deadline qorxusunda deyil, 1-ci həftədə tap.
(3) Test-ə çətin görünən şey = dizayn qoxusu (canary): async API, global state,
hardwired stdout — hamısı refaktorla test olunandır. (4) API-da concurrency
AÇMA — sinxron qaytar / callback ver; concurrency qərarını CALLER-a burax.
(5) Port 0 fəndi ilə hər test müstəqil random port alır. (6) Async testdə
fixed sleep YOX — wait-for-success + time.Timer + select; t.Fatal yalnız öz
goroutine-dən. (7) Greet nümunəsi: hardwired Print → parametr Fprint; əvvəl
test yazsaydın parametr məcburu ÖZÜ doğardı. (8) Test data təxəyyüllü olsun —
"Mary Jo" boşluqlu ad SİZİN funksiyanın sexist bug-unu aşkar etdi. (9) Fscan
space-separated-dır; tam sətirlər üçün bufio.Scanner. (10) Go-da HƏR funksiya
paralel çağırıla bilər — global mutable state YAZMA; Store+mutex; smoke test +
`go test -race` standart repsiyadır. (11) Uzun işlər: context.WithTimeout,
goroutine + cancel, <-ctx.Done(), ctx.Err() == DeadlineExceeded → FAIL.
(12) main yalnız ötürür: Main(os.Args) / ParseArgs+Sleep ayrılığı / config
struct-u (Timer) — CLI testi pure funksiya testinə çevrilir.

## Mənbə
Pages: 213-258 (PDF 225-270)
