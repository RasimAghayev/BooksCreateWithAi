# Go Optimizations 101 — Xarici References

## Kitabın öz resursları
- Kitabın səhifəsi: https://go101.org/optimizations/101.html
- Go 101 (müəllifin əsas kitabı): https://go101.org
- Müəllifin oyunları: https://tapirgames.com
- Go 101 issue tracker (feedback/düzəlişlər): https://github.com/go101/go101
- Go 101 Twitter: https://twitter.com/go100and1

## Rəsmi Go sənədləri (kitabda istifadə olunan)
- Go GC bələdçisi (GOGC/GOMEMLIMIT): https://go.dev/doc/gc-guide
- runtime/debug — SetGCPercent, SetMemoryLimit: https://pkg.go.dev/runtime/debug
- testing — Benchmark, AllocsPerRun: https://pkg.go.dev/testing
- Compiler optimizasiyaları wiki: https://github.com/golang/go/wiki/CompilerOptimizations
- unsafe.Sizeof / Alignof: https://pkg.go.dev/unsafe

## Kitabda adı çəkilən texniki məqalələr
- Go memory ballast (məqalə): https://blog.twitch.tv/en/2019/04/10/go-memory-ballast-how-i-learnt-to-stop-worrying-and-love-the-heap/
- GOMEMLIMIT təqdimatı / memory limit guide (kitabın tövsiyəsi): go.dev/doc/gc-guide daxilində
- Keith Randall, Ian Lance Taylor, Axel Wagner və b. — Go komandasının izahları (kitabın şəxsi dialoqları, müəllifin minnətdarlıq siyahısında)

## Compiler/CLI alətləri (kitabın əsas instrumentləri)
- go build -gcflags=-m / -m -m — escape analysis + inline cost
- go run -gcflags="-d=ssa/check_bce" — BCE qalıntıları
- go run -gcflags=-S — assembly çıxışı
- go test -bench — benchmark
- GODEBUG=gctrace=1 — GC cycle jurnalı
- GOMEMLIMIT / GOGC mühit dəyişənləri

## Əlaqəli oxu (seriya konteksti)
- Go 101 — müəllifin dil əsasları kitabı (go101.org)
- Go specification — runtime davranışların rəsmi tərifi: https://go.dev/ref/spec
- Go runtime source: https://github.com/golang/go/tree/master/src/runtime
