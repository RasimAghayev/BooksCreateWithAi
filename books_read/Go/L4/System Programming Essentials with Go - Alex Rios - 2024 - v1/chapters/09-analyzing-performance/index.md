# Chapter 9 — Analyzing Performance (Performansın Təhlili)

## Bu chapter nədən bəhs edir?
Escape analysis-in praktikasına (pointer/stack/heap interaksiyası), benchmark yazımına (b.N, sub-benchmark, benchmem, benchstat), 6 klassik benchmark tələsinə, CPU profiling-ə (pprof + file monitor nümunəsi), memory profiling-ə (WriteHeapProfile) və cache trade-off təcrübəsinə.

## Əsas fikirlər

### 1. Escape analysis — compiler landlord
**Nədir:** Compiler-in hər dəyişən üçün "stack icarəsinə layiqdirmi?" qərarı — funksiyadan qaçanlar heap-ə.

**Pointer və stack interaktivliyi:**
```go
a := 42
b := &a          // a-nın ünvanı → compiler "qaça bilər" deyə ehtiyat edir
*b = 21          // a artıq 21
```

**Klassik escape nümunəsi:**
```go
//go:noinline
func createPerson() *person {
    p := person{name: "Alex Rios", age: 99}
    return &p        // İSTİNAD qaytarılır → p escapes to heap!
}
```

**Analiz aləti:**
```bash
go build -gcflags "-m -m" .
```
**Çıxışın oxunması:**
```
./main.go:17:2: p escapes to heap:
  flow: ~r0 = &p:
    from &p (address-of) at ./main.go:18:9
    from return &p (return) at ./main.go:18:2
./main.go:17:2: moved to heap: p
./main.go:10:6: cannot inline main: function too complex: cost 141 exceeds budget 80
```
| Sətir | Məna |
|---|---|
| escapes to heap + flow | qaçış zənciri (haradan-hara) |
| moved to heap | heap allocation edildi |
| cost 141 exceeds budget 80 | inline mürəkkəblik büdcəsi (80) aşıldı |
| go:noinline | manual inline qadağası |

**Goroutine izolyasiyası:** Goroutine stack-i MÜTLƏQ özündür — başqa goroutine onun stack-inə pointer saxlaya BİLMƏZ (stack resize problemlərini aradan qaldırır).

### 2. Benchmark əsasları
```go
func BenchmarkSum(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Sum(1, 2)
    }
}
```
```bash
go test -bench=.
# BenchmarkSum-8  1000000000  0.277 ns/op
# ad-GOMAXPROCS | iterasiyalar | vaxt/əməliyyat
```
**b.N:** framework özü kalibrləyir — etibarlı ölçmə üçün iterasiya sayını tənzimləyir.

**Sub-benchmark (ssenari müqayisəsi):**
```go
func BenchmarkSumSub(b *testing.B) {
    cases := []struct{ name string; a, b int }{
        {"small", 1, 2}, {"large", 1000, 2000},
    }
    for _, c := range cases {
        b.Run(c.name, func(b *testing.B) {
            for i := 0; i < b.N; i++ { Sum(c.a, c.b) }
        })
    }
}
# BenchmarkSumSub/small-8 ... 0.3070 ns/op
# BenchmarkSumSub/large-8 ... 0.2970 ns/op
```

**Yaddaş ölçməsi:**
```bash
go test -bench=. -benchmem
# BenchmarkSum-8 ... 0.277 ns/op  16 B/op  2 allocs/op
```

### 3. benchstat — statistik müqayisə
```bash
go install golang.org/x/perf/cmd/benchstat@latest

go test -bench=. > old.txt       # dəyişiklikdən əvvəl
# ...kod dəyiş...
go test -bench=. > new.txt
benchstat old.txt new.txt
```
```
name           old time/op   new time/op   delta
BenchmarkSum-8 200ns ± 1%    150ns ± 2%   -25.00% (p=0.008 n=5+5)
```
**Interpretasiya:** delta mənfi = sürət artımı; p<0.05 = statistik əhəmiyyət; ±% = xəta hüdüdi; n = nümunə sayı.

### 4. Benchmark flag kombinasiyaları
```bash
go test -bench=BenchmarkMultiply -benchtime=3s -count=5
# seçici regex | minimum müddət | təkrar sayı
```
- Broad: `-bench=.` (hamısı); Sub: `-bench=BenchmarkMultiply/large`
- Diqqət: `-bench=Multiply` → `ComplexMultiply`-i də tutar!

### 5. 6 benchmark tələsi + həlləri
| Tələ | Həll |
|---|---|
| **1. Yanlış hədəf** (sort edilmiş slice yenidən sort) | hər iterasiyada state yenidən qur + b.ResetTimer() |
| **2. Compiler optimizasiyası** (nəticə istifadə olunmur → çağırış silinir) | nəticəni package-level dəyişənə yaz / runtime.KeepAlive |
| **3. Warmup yoxdur** (cache hələ istilənməyib) | qızdırma dövrü + sonra ResetTimer |
| **4. Mühit fərqi** (prod hardware-dan fərqli) | prod-a yaxın şərait: eyni Go runtime, realistic load |
| **5. GC/runtime xərcləri nəzərə alınmır** | uzun benchmark-lar; runtime metric izlə |
| **6. b.N istifadəsi** — 3 alt-tələ: rekurdiv funksiyada arqument kimi; loop-daxili setup üçün; iterasiya-şərtli kod üçün | b.N YALNIZ loop həddi; setup loop-dan XARİC; iterasiyadan asılı məntiq YOX |

### 6. CPU profiling — fayl monitor tətbiqi
**Layihə:** Kataloq dəyişikliyi monitoru — scan (WalkDir) → compare (yarat/sil/dəyiş) → alert (goroutine).

```go
type FileInfo struct {
    Name    string
    ModTime time.Time
    Size    int64
}

func scanDirectory(dir string) (map[string]FileInfo, error) {
    results := make(map[string]FileInfo)
    filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
        info, _ := d.Info()
        results[path] = FileInfo{Name: info.Name(), ModTime: info.ModTime(), Size: info.Size()}
        return nil
    })
    return results, nil
}

// main loop:
for {
    newState, _ := scanDirectory(dirToMonitor)
    compareAndEmitEvents(currentState, newState)   // go sendAlert(...) ilə asinxron
    currentState = newState
    time.Sleep(interval)
}
```

**Profiling qurulumu:**
```go
import "runtime/pprof"

f, _ := os.Create("cpuprofile.out")
defer f.Close()
pprof.StartCPUProfile(f)
defer pprof.StopCPUProfile()
```

**Analiz:**
```bash
go tool pprof cpuprofile.out
# Total: 10 samples
#   5  50.0% 50.0%   5  50.0% compareAndEmitEvents   ← HOTSPOT
#   3  30.0% 80.0%   3  30.0% scanDirectory
#   1  10.0% 90.0%   1  10.0% filepath.WalkDir

go tool pprof -web cpuprofile.out     # flame graph
```
**Flame graph:** geniş zolaq = çox CPU; iyerarxiya = call stack; yuxarıdan geniş yolları izlə.

### 7. Memory profiling
```go
f, _ := os.Create("memprofile.out")
defer f.Close()
runtime.GC()                     // zibili təmizlə → yalnız CANLI yaddaş görünür
pprof.WriteHeapProfile(f)
```
```bash
go tool pprof memprofile.out         # alloc hotspotları
go tool pprof -web memprofile.out    # flame (zolaq = alloc həcmi)
```
**Suallar:** leak? hansı funksiya ən çox allocate edir? yüklə altında pattern nə cür? obyekt ölçüləri?

**Zaman boyu:** loop daxilində interval-li heap profile yaz → müqayisə → gözlənilmədən qalan obyektlər = leak şübhəsi.

### 8. Cache trade-off eksperimenti
**Sadə cache (global map):**
```go
var cachedDirectoryState map[string]FileInfo    // qlobal

// scanDirectory başlanğıcında:
if cachedDirectoryState != nil {
    for path, fileInfo := range cachedDirectoryState {
        results[path] = fileInfo               // köhnə datadan başla
    }
}
// WalkDir sonunda: cachedDirectoryState = results
```
**Trade-off təhlili:**
- CPU: scanDirectory vaxtı azalmalı (təkrar walk azalır)
- Memory: cache-in özü YENİ istehlakçıdır — qəbul edilənmi?

**Caveat:** stale data riski (scan-lar arası dəyişikliklər); production keş üçün size limit + LRU eviction + invalidation trigger lazımdır.

**Qızıl qayda:** bir dəfədə YALNIZ BİR dəyişiklik profillə — təmiz müqayisə üçün.

## Əsas terminlər
- Escape analysis — stack/heap qərar mexanizmi
- -gcflags "-m -m" — escape + inline raporu
- Inline cost/budget (80) — mürəkkəblik həddi
- //go:noinline — manual qadağa
- Goroutine stack izolyasiyası
- b.N — kalibrlənən iterasiya sayı
- b.Run — sub-benchmark
- -benchmem — B/op + allocs/op
- benchstat — statistik delta + p-value
- -benchtime / -count — müddət/təkrar
- b.ResetTimer — setup-ı xaric et
- runtime.KeepAlive — optimizasiya qarşısı
- pprof.StartCPUProfile/StopCPUProfile
- runtime.GC() + WriteHeapProfile
- Flame graph — width = CPU/alloc
- LRU eviction — keş sıxışdırma

## Praktik nəticə
1. Escape şübhəsində `-m -m` ilə "flow" zəncirini oxu — return &p, interface, closure əsas qaçış səbəbləridir.
2. Benchmark nəticələri fərqləndirmək üçün benchstat + -count=5+ — yalnız statistik fərq REAL fərqdir.
3. Compiler-ə qurban vermə: nəticəni istifadə et (sink variable) — yoxsa benchmark boş ölçür.
4. Profiling workflow: StartCPUProfile (defer Stop) → workload → go tool pprof (text) → -web (flame) → top few optimallaşdır.
5. Heap profili QARŞIDA GC çağır — dəqiq canlı yaddaş görünümü üçün.
6. Hər optimizasiya addımından sonra hər iki profili YENİDƏN götür (adları dəyişməyi unutma) — trade-off-ları rəqəmlə gör.

## Mənbə
Pages: 161-188 (PDF səh. 182-209)
