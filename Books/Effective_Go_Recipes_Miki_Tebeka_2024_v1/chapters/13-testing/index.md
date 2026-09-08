# Chapter 13 — Testing Your Code (səh. 197-217)

## Bu fəsil nədən bəhs edir?

Testinq ədəbiyyatı: CI-yalnız testlər (env yoxlama + t.Skip), YAML-dan
test case-lər, fuzz testinq, HTTP mock (RoundTripper), TestMain (qlobal
setup/teardown), real server-in end-to-end testi (os/exec) və xüsusi
linter yazılması (go/analysis).

## Əsas fikirlər

### Recipe 65 — CI-yalnız testlər
**Problem:** test dəsti 20 dəqiqə — developer maşını 5 dəqiqə limiti.

```go
var (
    inCi = os.Getenv("CI") != ""      // CI sistemi CI env qoyur
)

func TestMonthlyReport(t *testing.T) {
    if !inCi {
        t.Skip("not in CI system")     // SKIP statusu
    }
    t.Log("running in CI")
    // uzun test...
}
```
```bash
$ go test -v                            # --- SKIP: TestMonthlyReport
$ CI=true go test -v                    # PASS
$ CI=true go test -run TestMonthlyReport -v   # tək CI test
```
- CI: GitHub Actions, CircleCI, GitLab CI/CD; "works on my machine"-dən
  qoruyur
- Developerləri məmnun saxlayın — yoxsa künc kəsməyə başlayacaqlar

### Recipe 66 — YAML-dan test case-lər
**Tapşırıq:** QA Go yazmadan test əlavə etsin (bin packing, NP-hard!).

```yaml
# packer_cases.yml
- name: simple
  err: false
  box_capacity: 6
  weights: [2.3, 3.7, 5]
  num_boxes: 2
- name: over
  err: true
  box_capacity: 10
  weights: [1.2, 37, 9.5]
```

```go
type TestCase struct {
    Name      string    `yaml:"name"`
    BoxCap    float64   `yaml:"box_capacity"`
    N         int       `yaml:"num_boxes"`
    Weights   []float64 `yaml:"weights"`
    ShouldErr bool      `yaml:"err"`
}

// Helper — error yox, t.Fatal (test üslubu!):
func loadCases(t *testing.T, path string) []TestCase {
    file, err := os.Open(path)
    if err != nil {
        t.Fatal(err)
    }
    defer file.Close()
    var tcs []TestCase
    if err := yaml.NewDecoder(file).Decode(&tcs); err != nil {
        t.Fatal(err)
    }
    return tcs
}

func TestPack(t *testing.T) {
    testCases := loadCases(t, "packer_cases.yml")
    for _, tc := range testCases {
        items := weightsToItems(tc.Weights)
        t.Run(tc.Name, func(t *testing.T) {    // SUBTEST — ad YAML-dan!
            boxes, err := Pack(tc.BoxCap, items)
            if err == nil && tc.ShouldErr {
                t.Fatal("expected error, got nil")
            }
            if err != nil && !tc.ShouldErr {
                t.Fatalf("unexpected error: %s", err)
            }
            if n := len(boxes); n != tc.N {
                t.Fatalf("expected %d boxes, got %d", tc.N, n)
            }
        })
    }
}
```
- Test helper-ləri error qaytarmır — `t.Fatal` (qısa, aydın)
- Subtest adı YAML-dən → xətanı dərhal tanıyırsan

### Recipe 67 — fuzzing
**Tapşırıq:** `NumBytes(n)` — n bit üçün bayt sayı: `(n+7)/8`.

```go
func FuzzNumBytes(f *testing.F) {
    f.Add(0)                      // seed nümunələr

    fn := func(t *testing.T, n int) {
        if n < 0 {
            return                 // mənfilər icarə xarici
        }
        nBytes := NumBytes(n)
        // HEURISTİKA: nə çox, nə az:
        ok := (nBytes*8 >= n) && ((nBytes-1)*8 <= n)
        if !ok {
            t.Fatal(nBytes)
        }
    }
    f.Fuzz(fn)                    // random dəyərlərlə minlərlə icra
}
```
```bash
$ go test -run NONE -fuzz . -fuzztime 10s
# execs: 1848027 ... PASS
```
- Uğursuz nümunə `testdata/fuzz`-a yazılır → sonrakı run-larda təkrar
  yoxlanılır
- Dəstəklənən tiplər: string/[]byte, int/uint ailəsi, float, bool;
  mürəkkəblər üçün testing/quick
- **Çətin hissə:** uğur/uğursuzluq heuristikası; "negative numbers!" kimi
  sürprizlər tapılır

### Recipe 68 — HTTP client mock (RoundTripper)
**Tapşırıq:** Users() metodunu server işlətmədən test et.

```go
type MockTransport struct {
    body []byte
    err  error
}

// RoundTrip implements http.RoundTripper
func (t *MockTransport) RoundTrip(r *http.Request) (*http.Response, error) {
    if t.err != nil {
        return nil, t.err                       // şəbəkə xətası emulyasiyası
    }
    w := httptest.NewRecorder()
    if t.body != nil {
        w.Write(t.body)                         // cavab body
    }
    return w.Result(), nil
}

// Uğurlu yol:
func TestUsersOK(t *testing.T) {
    users := []string{"clark", "diana", "bruce"}
    data, _ := json.Marshal(users)
    c := NewAPIClient("http://localhost:8080")
    c.c.Transport = &MockTransport{data, nil}   // MOCK taxılı
    reply, err := c.Users()
    require.NoError(t, err)
    require.Equal(t, users, reply)
}

// Xəta yolları:
c.c.Transport = &MockTransport{nil, fmt.Errorf("network error")}
c.c.Transport = &MockTransport{[]byte(`["clark","diana","bruce"`), nil} // pis JSON
```
- `http.Client.Transport` = `RoundTripper` interfeysi — mock nöqtəsi
- `httptest.NewRecorder().Result()` — Response-ı asan qurmaq
- **Xəbərdarlıq:** mock = aldatma; /users → /api/users dəyişəndə test
  keçəcək, real kod çökəcək. testify/mock daha güclü imkanlar verir

### Recipe 69 — TestMain (qlobal setup/teardown)
```go
func setupTests() error {
    file, err := os.CreateTemp("", "*.yml")
    if err != nil {
        return err
    }
    defer file.Close()
    cfg := map[string]any{
        "verbose": true,
        "dsn":     "postgres://localhost:5432",
    }
    if err := yaml.NewEncoder(file).Encode(cfg); err != nil {
        return err
    }
    os.Setenv(envConfigKey, file.Name())
    return nil
}

func teardownTests() {
    fileName := os.Getenv(envConfigKey)
    if err := os.Remove(fileName); err != nil {
        log.Printf("warning: can't delete %q - %s", fileName, err)
    }
}

func runTests(m *testing.M) int {          // defer üçün ayrı funksiya!
    if err := setupTests(); err != nil {
        return 1
    }
    defer teardownTests()                  // os.Exit defer-i çağırmır!
    return m.Run()
}

func TestMain(m *testing.M) {
    code := runTests(m)
    os.Exit(code)                          // sıfır olmayan kod = xəta
}
```

### Recipe 70 — test-də real server (end-to-end)
```go
func buildServer(t *testing.T) string {
    fileName := path.Join(t.TempDir(), "httpd")
    cmd := exec.Command("go", "build", "-o", fileName, "httpd.go")
    if err := cmd.Run(); err != nil {
        t.Fatal(err)
    }
    return fileName
}

func freePort(t *testing.T) int {          // boş port tap
    conn, err := net.Listen("tcp", "")
    if err != nil {
        t.Fatal(err)
    }
    conn.Close()
    return conn.Addr().(*net.TCPAddr).Port
}

func waitForServer(t *testing.T, addr string) {   // hazır olmasını gözlə
    start := time.Now()
    timeout := 10 * time.Second
    var err error
    var conn net.Conn
    for time.Since(start) < timeout {
        conn, err = net.Dial("tcp", addr)
        if err == nil {
            conn.Close()
            return
        }
        time.Sleep(10 * time.Millisecond)
    }
    t.Fatalf("server not ready after %s (%s)", timeout, err)
}

func runServer(t *testing.T) int {
    exe := buildServer(t)                 // BUILD — go run YOX!
    port := freePort(t)

    env := os.Environ()                   // env + port
    env = append(env, fmt.Sprintf("HTTPD_ADDR=:%d", port))
    cmd := exec.Command(exe)
    cmd.Env = env
    if err := cmd.Start(); err != nil {
        t.Fatal(err)
    }

    addr := fmt.Sprintf("localhost:%d", port)
    waitForServer(t, addr)

    t.Cleanup(func() {                    // test bitəndə öldür
        if err := cmd.Process.Kill(); err != nil {
            t.Logf("warning: can't kill server")
        }
    })
    return port
}

func TestHealth(t *testing.T) {
    port := runServer(t)
    resp, err := http.Get(fmt.Sprintf("http://localhost:%d/health", port))
    if err != nil {
        t.Fatal(err)
    }
    if resp.StatusCode != http.StatusOK {
        t.Fatalf("bad return code: %d", resp.StatusCode)
    }
}
```
- **go run YOX, build + Start:** go run öz övlad prosesini yaradır — Kill
  yalnız go-nu öldürür, serveri yox
- `t.Cleanup` — test sonrası avtomatik təmizlik
- `t.TempDir()` — test bitəndə avtomatik silinən qovluq

### Recipe 71 — xüsusi linter (go/analysis)
**Tapşırıq:** qadağan paketlərin (syscall) istifadəsini yaxala.

```go
var forbidden = []string{"syscall"}

func isForbidden(path string) bool {
    for _, prefix := range forbidden {
        if strings.HasPrefix(path, prefix) {
            return true
        }
    }
    return false
}

Analyzer = &analysis.Analyzer{
    Name:     "forbidden",
    Doc:      "Check for usage of forbidden packages",
    Requires: []*analysis.Analyzer{inspect.Analyzer},
    Run:      run,
}

func run(pass *analysis.Pass) (any, error) {
    filter := []ast.Node{(*ast.ImportSpec)(nil)}     // yalnız import-lar
    inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
    inspect.Preorder(filter, func(node ast.Node) {
        imp, ok := node.(*ast.ImportSpec)
        if !ok {
            return
        }
        path := strings.Trim(imp.Path.Value, `"`)    // `"syscall"` → syscall
        if !isForbidden(path) {
            return
        }
        pass.Reportf(imp.Pos(), "importing forbidden package %q", path)
    })
    return nil, nil
}

func main() {
    singlechecker.Main(Analyzer)     // CLI alətinə çevir
}
```

**Test (annotasiya əsaslı):**
```go
// testdata/app.go
import (
    "fmt"
    "syscall" // want `importing forbidden package "syscall"`
)

func TestAnalyzer(t *testing.T) {
    analysistest.Run(t, analysistest.TestData(), Analyzer)
}
```
- ~55 sətirlik linter! Xətanın düzəldilməsi dəyəri zamanla ARTIR —
  erkən tutmaq = qənaət
- staticcheck kimi linter-lər testlərdən ƏVVƏL; make/script avtomatlaşdırır

## Final Thoughts-dən

SQLite: hər source sətrinə 640 test sətri — və yenə də bug var! Test =
"pain vs gain" balansı (nahar seçənə vs təyyarə idarə edən sistem).
NASA belə Marsa bug göndərir (və Marsda düzəldir!).

## Əsas terminlər
- CI (Continuous Integration) — avtomatik build/test
- t.Skip / t.Fatal — test status idarəetməsi
- Test helper üslubu — error yox, t.Fatal
- Table-driven (YAML) — kod-xarici test case-lər
- Fuzzing — random input ilə xəta axtarışı
- Fuzz heuristics — uğur şərti
- testdata/fuzz — tapılan uğursuz nümunələr anbarı
- http.RoundTripper — transport mock nöqtəsi
- httptest.NewRecorder — Response qurucusu
- TestMain / fixtures — qlobal setup/teardown
- t.Cleanup / t.TempDir — avtomatik təmizlik
- go/analysis framework — linter infrastrukturu
- analysistest — "want" annotasiyalı linter testləri

## Praktik nəticə
Uzun testləri CI-yə ayırın (CI env + t.Skip); test case-ləri YAML-a
çəkin — QA kod yazmadan əlavə edir; subtest adları ilə diaqnoz asanlaşır.
Fuzz: mümkün hər funksiya üçün heuristika + f.Add seed-lər. HTTP klient
testi üçün Transport mock-u (3 ssenari: OK, şəbəkə xətası, pis JSON) —
amma "mock = aldatma" unutmayın. Qlobal fixture-lər TestMain+runTests
defer triki ilə; real end-to-end üçün build+Start+freePort+waitForServer+
t.Cleanup. Xüsusi qaydalar üçün go/analysis — 55 sətirdə linter!

## Mənbə
Pages: 197-217 (PDF 197-217)
