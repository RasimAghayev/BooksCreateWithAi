# Chapter 6 — Formatting, testing, debugging, and benchmarking (Formatlaşdırma, testləmə, debug və benchmark)

## Bu chapter nədən bəhs edir?

gofmt/goimports (IDE inteqrasiyası), go vet (context leak, JSON tag yoxlamaları), go mod
versiya idarəetməsi, log paketi (flags, custom Writer, fayl çıxışı), slog struktur log
(1.21), stack trace (debug.PrintStack, runtime.Stack), table-driven testlər, fuzzing
(Go 1.18+), subtest adları, coverage, Delve debugger və benchmarklar (Bubble Sort).

## Əsas fikirlər

### 1. gofmt və goimports
**gofmt:** Kanonik stil; tab indent, artıq whitespace/semicolon təmizləmə. `go fmt` =
gofmt alias. `gofmt -d` — dəyişikliyi diff kimi göstər (faylı dəyişmədən). Rob Pike:
"gofmt's style is no one's favorite, yet gofmt is everyone's favorite" — stil mübahisələri
bitir.

**goimports:** gofmt + itkin import-ları əlavə edir, istifadəsizləri silir:
```bash
go install golang.org/x/tools/cmd/goimports@latest
```
**IDE:** GoLand (File Watchers), VS Code (Format on Save); gopls — LSP server hər IDE-yə.

### 2. go vet — Səssiz Səhvlər
Compile KEÇƏN amma yanlış kodu tutur:
- **Context leak:** `ctx, _ := context.WithCancel(...)` — cancel iqnor → vet:
  "the cancel function returned by context.WithCancel should be called, not discarded,
  to avoid a context leak"
- **JSON tag + unexported sahə:** `username string \`json:"username"\`` — sahə export
  olunmadığından marshal NƏTİCƏSİZ; vet: "struct field username has json tag but is not
  exported" (tagsiz unexported sahəyə susur — niyyət bəlli deyil)
- Digərlər: value-vs-pointer marshal, malformed testlər, buffered olmalı kanallar, yanlış
  printf verb-ləri. False positive mümkündür — hər gün yox, build öncəsi sanity check.
- Handler wrapper pattern (chapter-də): `homeHandler(ctx) http.HandlerFunc` — ctx-i
  tutan closure qaytarır; ctx.Value + type assertion ilə konfiq keçirilir.

### 3. Asılılıqların Yenilənməsi
```bash
go get pkg@0.1.1        # konkret versiya (SemVer teqi)
go get -u && go mod tidy # hamını yenilə + asılılıq ağacını təmizlə
go get -u pkg            # tək paket
go mod graph             # versiya ağacını göstər
# -t: test paketlərini də nəzərə al
```

### 4. log Paketi
```go
log.Println("...")   // stdout: 2023/08/29 15:52:24 ...
log.SetFlags(log.Ltime)
log.SetFlags(log.Llongfile)              // tam yol + sətir
log.SetFlags(log.LUTC | log.Lshortfile)  # birləşdir (OR)
```
Flag-lar: Ldate, Ltime, Lmicroseconds, LstdFlags, Llongfile, Lshortfile (uzun/qısa
uyğunsuzdur). **Custom Writer:** `log.SetOutput(myWriter)` — Write([]byte) (int, error)
implement edən hər şey; runtime.CallersFrames ilə çağıran funksiyanı, rəng kodları və
öz timestamp formatı əlavə etmək olar (rəng — opt-in flag ilə; fayla escape zibil!).

**Fayla log:**
```go
file, _ := os.OpenFile("logging.log", os.O_RDWR|os.O_CREATE, 0755)
log.SetOutput(file)
```
12-factor: stdout stream + shell redirection (`go run x.go &> file.log`) daha yaxşı.

### 5. slog — Struktur Log (Go 1.21)
```go
slog.Info("this is default logging")
slog.Warn(...); slog.Error(...); slog.Debug(...)   // Debug default gizlidir!
```
JSON handler + səviyyə idarəsi:
```go
logger := slog.New(slog.NewJSONHandler(file, &slog.HandlerOptions{
    Level: slog.LevelDebug,
}))
slog.SetDefault(logger)
slog.Info("complex message",
    slog.String("accepted_values", "..."),   // typed key/value
    slog.Int("an int:", 30),
    slog.Group("grouped_info", slog.String("you_can", "do this too")),  // iyerarxiya
)
```
Çıxış — JSONS (newline-ayrı JSON mesaj axını); analitik platformalar dəstəkləyir; `jq`
ilə sorğulanır. Struct + tag əvəzinə slog tipli sahələrlə sürətli JSON.

### 6. Stack Trace
```go
import "runtime/debug"
debug.PrintStack()          // stdout-a tam trace (main → foo → bar → runtime)

import "runtime"
buf := make([]byte, 1024)   // ölçünü əvvəlcədən təyin etməlisən
runtime.Stack(buf, false)   // false = yalnız cari goroutine; true = BÜTÜN goroutine-lər
```
Çox-goroutine trace paralellik debug-ında dəyərli (amma böyükdür). runtime.Callers /
CallersFrames — proqramatik frame girişi.

### 7. Table-Driven Testlər
```go
func Test_FizzBuzz(t *testing.T) {
    tests := []struct {
        name     string
        input    int64
        expected string
    }{
        {"fizz buzz test1", 37, "1, 2, Fizz, 4, Buzz, ..."},
        {"fizz buzz test2", 5, "1, 2, Fizz, 4, Buzz"},
    }
    for i := range tests {
        test := tests[i]
        t.Run(test.name, func(t *testing.T) {     // adlı subtest
            if res := fizzbuzz(test.input); res != test.expected {
                t.Fatalf("\ngot \n%s \nexpected \n%s", res, test.expected)
            }
        })
    }
}
```
- `-v` çıxışda hər subtest ayrı göstərilir (IDE inteqrasiyası da bunu istifadə edir)
- Adlar declarativ: "should generate ..." üslubu
- Kitabın bug dərsi: FizzBuzz-da 0-başlangıç + off-by-one + artıq boşluq + vergül
  ayracı — əl testi görməzdi, cədvəl dərhal yaxaladı
- `t.Fatalf` — yalnız növbəti testlər asılı olanda; `t.Errorf` — müstəqil yoxlamalarda

### 8. Fuzzing (Go 1.18+)
**Nədir:** Gözlənilməyən input-ları avtomatik yaradaraq edge-case tapma.

**Kitabdan kod nümunəsi (int16 overflow):**
```go
func FuzzSummedRuneCodes(f *testing.F) {
    tests := []string{"i am trying things", "..."}   // seed-lər
    for t := range tests { f.Add(tests[t]) }
    f.Fuzz(func(t *testing.T, seed string) {
        got := summedRuneCodes(seed)
        if got < 0 {        // invariant: nəticə heç vaxt mənfi olmamalı
            t.Errorf("we got %d from string %s", got, seed)
        }
    })
}
```
```bash
go test --fuzz=Fuzz    # uğursuz seed-i tapır → faylda saxlanılır → corpus-a düşür
```
Gözlənilən cavab YOXDUR — ümumi invariant (mənfilik, overflow) yoxlanılır. Təhlükəsizlik
dəyəri: uzun input, xüsusi simvollar, SQL injectionvari hücum formaları. LLM-lər test
case/artıq edge-case generasiyasında köməkçi (yoxlanış mütləq).

### 9. Coverage
```bash
go test -v -cover
# coverage: 66.7% of statements
go test -v -cover -coverprofile mycover.out
go tool cover -html mycover.out    # brauzerdə yaşıl/qırmızı sətir görünüşü
```

### 10. Debugger
Go-da rəsmi debugger YOXDUR (GDB plugin qeyri-stabil). **Delve** (go-delve/delve) —
cəmiyyət standartı; VS Code Go extension-a built-in.

### 11. Benchmark
```go
func runBenchmark(arr []int, runs int) {
    for i := 0; i < runs; i++ { bubbleSort(arr) }
}
func BenchmarkBubbleSort10(b *testing.B)  { runBenchmark(generateRandoms(10), b.N) }
func BenchmarkBubbleSort100(b *testing.B) { runBenchmark(generateRandoms(100), b.N) }
...
```
```bash
go test -bench=.
# BenchmarkBubbleSort10-8       248036485    4.703 ns/op
# BenchmarkBubbleSort100000-8          1    13268521541 ns/op
```
- `-bench=PATTERN` regex; `-benchtime=10s` minimal vaxt (default 5s? kitabda 5s; 1M-də
  cəmi 1 run — vaxtı artıraraq daha çox data)
- `-benchmem` — allocation məlumatı
- **benchstat** (golang.org/x/perf) — A/B müqayisə + statistik analiz
- Bubble Sort nəticəsi: 10 elementdə 4.7ns, 100k-da 13.2 SANIYƏ — O(n²) realda

## Əsas terminlələr
- gofmt/gofmt -d — kanonik format / diff rejimi
- go vet — compile-keçən səhvlərin statik analizi
- L* flag-ları — log prefix elementləri
- slog — struktur log paketi (1.21); JSONS — newline-JSON axını
- debug.PrintStack / runtime.Stack — trace çapı / buffera
- Table-driven test — input/expected struct slice-u
- Fuzzing — random input ilə invariant yoxlama
- t.Run — adlı subtest
- b.N — benchmark təkrar dəyişəni
- Delve — cəmiyyət debugger-i

## Praktik nəticə

Keyfiyyət axını: (1) save-də goimports; (2) build öncəsi go vet (context cancel və JSON
tag tələlərini özü tutur); (3) `log` fmt-dən üstün; dağıtılacaq proqram üçün slog + JSONS
jq ilə; (4) hər funksiyaya table test (adlı subtestlər); (5) əl ilə tapılmayan hücum
halları üçün fuzz invariant-ları; (6) coverage-i coverprofile+html ilə vizuallaşdır; (7)
performans iddialarını -bench ilə sayallaqlaşdır (-benchmem, benchstat); (8) debug üçün
Delve; (9) `--race`/`--fuzz`/`--cover` flag-ləri CI-ya.

## Mənbə
Pages: 138-168 (PDF 159-189)
