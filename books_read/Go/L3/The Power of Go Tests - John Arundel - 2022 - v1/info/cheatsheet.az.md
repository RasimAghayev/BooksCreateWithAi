# The Power of Go: Tests — Cheat Sheet (Azərbaycanca)

John Arundel, Bitfield Consulting 2022 — 11 fəsil, L3 Intermediate

## 1. Test Strukturunun Esası (Ch1-3)
```go
func TestSquare_IsCorrect(t *testing.T) {
    t.Parallel()
    want := uint64(9)
    got := square.Square(3)
    if want != got {
        t.Error(cmp.Diff(want, got))     // fail mesajı: fərq + kontekst
    }
}
```
- Ad formatı: `Test<Funksiya>_<Davranış>`
- Want əvvəl, got sonra; müqayisə üçün `cmp.Diff`

## 2. Table-Driven Test (Ch3)
```go
for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        t.Parallel()
        // ... input: tt.in, want: tt.want
    })
}
```
- Hər sətir subtest; input tam təsvir edir → name = inputun ÖZÜ
- `go test -run TestSquare/inputs[0-9]+` — konkret subtest

## 3. Error Testləri (Ch4)
```go
wantErr := ""            // "found" / "not found" / ""
got, err := f(x)
if got != 0 { t.Errorf("want 0, got %d", got) }      // zero value qaytar!
if wantErr == "found" && err == nil { ... }
if wantErr == "not found" && err != nil { ... }
```
- errors.Is(err, ErrX) / errors.As(err, &target)

## 4. Invalid Input (Ch5)
```go
tests := []struct{ ... }{
    {name: "empty slice", in: []int{}},
    {name: "nil map", in: nil},
    {name: "huge string", in: strings.Repeat("x", 1000000)},
}
```

## 5. Fuzz (Ch6)
```go
func FuzzGuess(f *testing.F) {
    f.Add("seed1")                    // seed corpus
    f.Fuzz(func(t *testing.T, s string) {
        got := f(s)                    // oracle / property yoxla
        if got < 0 { t.Errorf(...) }
    })
}
```
```bash
go test -fuzz . -fuzztime=10m        # fail inputlar testdata/fuzz-a düşür
go test -run=FuzzGuess               # corpus regression kimi
```
- rand determinizmi: default seed=1 → REPRODUCIBLE (flaky deyil!)
- rand.Perm(100) — tam örtük + random sıra

## 6. Coverage (Ch7)
```bash
go test -cover                            # faiz
go test -coverprofile=cover.out           # profil
go tool cover -html=cover.out             # qırmızı/yaşıl
```
- 85→100% testləri adətən DƏYMƏZ; coverage ratchet: aşağı salan check-in YOX

## 7. Mutation (Ch7)
```bash
go install github.com/avito-tech/go-mutesting/cmd/go-mutesting@latest
# BACKUP ED!!! (.git daxil) sonra: go-mutesting .
```
- Score 1.0 = bütün mutantlar ölü; "passed" = mutant öldürüldü (TERS!)
- Survivable mutant → feeble test YAXUD lazımsız kod

## 8. Async Test (Ch8)
```go
func randomLocalAddr(t *testing.T) string {         // port 0 fəndi
    t.Helper()
    l, err := net.Listen("tcp", "localhost:0")
    if err != nil { t.Fatal(err) }
    defer l.Close()
    return l.Addr().String()
}

func waitForServer(t *testing.T, addr string) {    // wait-for-success
    t.Helper()
    timeout := time.NewTimer(100 * time.Millisecond)
    _, err := net.Dial("tcp", addr)
    for err != nil {
        select {
        case <-timeout.C:
            t.Fatal("timed out")
        default:
            time.Sleep(time.Millisecond)
            _, err = net.Dial("tcp", addr)
        }
    }
}
```
- fixed sleep = FLAKY/YAVAŞ — həmişə wait-for-success
- t.Fatal yalnız test-in ÖZ goroutine-ində işləyir

## 9. Race + Context (Ch8)
```bash
go test -race       # smoke testlərlə birlikdə
```
```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
defer cancel()
go func() { job(ctx); cancel() }()
<-ctx.Done()
if errors.Is(ctx.Err(), context.DeadlineExceeded) { t.Fatal("timed out") }
```

## 10. I/O Parametrləri (Ch8)
```go
func Greet(in io.Reader, out io.Writer) {        // test:
    fmt.Fprint(out, "Your name? ")              // buf := new(bytes.Buffer)
    scanner := bufio.NewScanner(in)             // input := strings.NewReader("John\n")
    if !scanner.Scan() { return }                // Fscan boşluqda KƏSİR! Scanner XEYR
    fmt.Fprintf(out, "Hello, %s!\n", scanner.Text())
}
```

## 11. CLI Test (Ch9)
```go
func main() { timer.Main(os.Args) }              // delegate pattern

func TestMain(m *testing.M) {
    os.Exit(testscript.RunMain(m, map[string]func() int{
        "timer": timer.Main,
    }))
}
```
```
# testdata/script/hello.txtar
exec timer 10ms
stdout 'sleeping for 10ms'
! stderr .
-- golden.txt --
istənilən fayl məzmunu burada
```
- `! exec` = fail gözlə; double quote LITERALDIR ('' single'ı escape edir)

## 12. Dependency Arsenalı (Ch10)
```go
type Store interface { Store(Widget) (string, error) }   // adapter

type mapStore struct { m *sync.Mutex; data map[string]widget.Widget }  // fake
func (ms *mapStore) Store(w widget.Widget) (string, error) { ... }

func Create(s Store, w Widget) (string, error) {       // biznes məntiqi DB-siz
    return s.Store(w)
}
```
```go
var Now = time.Now             // seam: testdə fake funksiya inject et
past.Now = func() time.Time { return testTime }   // PARALEL DEYİL!
```
- Prioritet: ləğv et > scope azalt > chunk > adapter > fake > mock (SON)
- İmtina: mock = indirect outputs → brittle; interface-i YALNIZ mock üçün YARATMA

## 13. Vaxt Müqayisəsi (Ch10-11)
```go
delta := want.Sub(got).Abs()          // == İŞLƏTMƏ
if delta > 10*time.Microsecond { t.Errorf(...) }
```

## 14. Seçici Müqayisələr (Ch11)
```go
cmp.Equal(want, got)                                      // map: sıra önəmsiz
cmpopts.SortSlices(func(a, b int) bool { return a < b })  // slice: sırasız yoxla
cmpopts.EquateApprox(0, 0.00001)                          // float epsilon
```
- Epsilon tuning: kiçilt → fail → 10X böyüt

## 15. Preconditions (Ch11)
```go
if user.Exists("Alice") { t.Fatal("precondition: Alice mövcud OLMAMALI") }
user.Create("Alice")
if !user.Exists("Alice") { t.Error("not created") }
// Delete üçün: İKİ user — Bob SAĞ qalmalı (nuke-bug qorunması)
```

## 16. Suite Sağlamlığı (Ch11)
- Flaky test → düzəlt və ya SİL ("bad tests are worse than no tests")
- Fail-ə icazə YOX → zero defects; backlog böyükdürsə → bug bankruptcy
- DB testləri: öz DB > özn cədvəl > transaction+rollback; prod-a qarşı da təhlükəsiz olmalı
- Beck time: suite ≤ 10 dəq (paralel, I/O-suz, lokal fake, sleepsiz, gecəlik slow)
- Framework YOX: `import "testing"` + go-cmp; assert.Equal fərhi İZAHLAMIR

## 17. Sorğu Çarpanları
| Problemin | Sorğu |
|-----------|-------|
| Test çətindir | "What are we really testing here?" |
| Nə yoxlamaq? | Davranış, implementasiya deyil |
| Random inputda want? | Property (invariant) yaz |
| Mock lazımdır? | Əvvəl ləğv/scope/chunk/adapter sına |
| Test flaky-dirsə | SİL və ya kökü tap |
| Azmı çoxmu test? | Fowler: dəyişə bilmirən → az; çox qırılır → çox |
