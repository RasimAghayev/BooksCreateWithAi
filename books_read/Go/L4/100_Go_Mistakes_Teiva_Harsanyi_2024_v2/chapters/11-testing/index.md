# Chapters 11-12 — Testing, Optimizations (#82-#100) (səh. 282-375)

## Bu fəsillər nədən bəhs edir?

(11) Test kategorizasiyası, -race, parallel/shuffle, TDT, sleep-siz testlər, time
injeksiyası, httptest/iotest, dəqiq benchmarklar, Go test xüsusiyyətləri.
(12) CPU cache-lər, false sharing, instruction-level parallelism, data
alignment, stack vs heap, escape analysis, allocation azaltma, inlining,
pprof, GC, Docker/K8s CFS quota.

## Əsas səhvlər və həllər

### #82: Testlərin kategoriyalaşdırılmaması
- Test piramidi: unit çox → integration → E2E az
- Ayırma üsulları: build tag-lər (`//go:build integration`), env dəyişənləri,
  `-short` + `testing.Short()` ilə skip
- `go test -short -v .` → long-running skip olunur

### #83: -race bayrağı
- `go test -race` — data race-ləri yaxalayır (yaddaş 5-15×, CPU 2-20×
  overheadi ilə) → lokal + CI-də MƏCBURİ, production-da YOX
- Yalnız testlər yaxşıdırsa işləyir (false negatives mümkün)

### #84: Test icra modları
- `t.Parallel()` — paralel testlər (`-parallel 16` ilə məhdudlaşdır);
  `//go:build` fayllarının icrası
- `-shuffle=on` (Go 1.17+) — sıra asılılıqlarını ifşa edir

### #85: TDT (table-driven tests)
```go
tests := map[string]struct{ input, expected string }{
    "suffix with multiple new lines": {input: "a\n\n\n", expected: "a"},
}
for name, tt := range tests {
    t.Run(name, func(t *testing.T) { ... })
}
```
- map istənilən sıra (determinizm üçün də yaxşı); yeni hal = 1 sətir

### #86: Unit testlərdə sleep
- Sleep = FLAKINESS siqnalı. Alternativlər: kanallarla mock (publisherMock
  notify-i await edir), `testify` `Eventually` (retry + timeout)
- Gözləmə logikası: async emali SYNC-ə çevir (channel ilə "hazır" siqnalı)

### #87: Time API testlərdə
```go
type Cache struct {
    now func() time.Time   // unexported funksiya sahəsi — injeksiya nöqtəsi
}
// istehsalda: now: time.Now; testdə: sabit vaxt
```
- Daha da yaxşısı: `TrimOlderThan(t time.Time)` — "indi"ni parametr et

### #88: httptest / iotest
- `httptest.NewRequest("GET", "/x", nil)` + `httptest.NewRecorder()` — handler
  testləri server-siz
- `httptest.NewServer(handler)` — real URL lazım olanda
- `iotest.TestReader` / `iotest.OneByteReader` — custom Reader-lərin düzgünlüyü

### #89: Dəqiq benchmarklar
- `b.ResetTimer()` setup-dan sonra (setup vaxtı hesaba girməsin)
- `b.StopTimer/StartTimer` — hər iterasiyada dynamically
- **Mikro-benchmark yanlışlığı:** -count=10 + benchstat ilə statistik müqayisə
- **Compiler optimizasiyası:** `Sink = result` global sink (ölçmə silinməsin)
- **Observer effekti:** eyni matris üzrə iterasiya = cache-dən oxuyur → HƏR
  loop-da yeni data (512 vs 513 sütun fərqi: 50% — critical stride!)

### #90: Test xüsusiyyətləri
- Code coverage `-cover`; **external package test** (`foo_test`) — public
  API-dən baxış
- Utility funksiyalar üçün `testing.AlperPerks` yox — sadəcə helper-lər
- Setup/teardown: `t.Cleanup(func(){...})` — defer-in test-mərkəzli forması

### #91: CPU cache-lər (mechanical sympathy)
- L1/L2/L3; **cache line = 64 bayt**; spatial + temporal locality
- **Unit stride** (slice) = prefetch mümkün; **non-unit** (linked list,
  pointer slice) = hər element üçün miss
- SumFoo vs SumBar (struktur daxilində slice): əlaqəli elementləri BİRLİKDƏ
  saxla → 1 fetch çox element gətirir
- **Critical stride:** 512 sütunlu matris = hər sətir EYNİ cache set-lərinə
  düşür (conflict miss) → 513 FƏRQLİ 50% sürətli!

### #92: False sharing
- İki goroutine muxtasır dəyişənləri (sumA, sumB) eyni cache line-da → L1D
  per-core REPLİKASIYA → hər yazma digər nüvəni INVALIDATE edir
- Həll: **padding** — `_ [56]byte` sahəsi ilə dəyişənləri ayrı line-lara

### #93: Instruction-level parallelism
- CPU pipeline: data hazard (RAW — növbəti təlimat əvvəlkinin nəticəsini gözləyir)
- `s[0]++; s[1]++` ardıcıl asılılıq → temp dəyişənlə paralleliliyi artır

### #94: Data alignment
- struct sahələri ölçüyə görə ALIGN olunur; compiler padding əlavə edir
- **Sahələri BÖYÜKDƏN KİÇİYƏ sırala** → daha az padding → daha kiçik struct

### #95: Stack vs heap
- Stack: goroutine başına 2KB (böyüyür), self-cleaning, GC YOX
- Heap: GC təmizliyi (25% CPU); escape analysis qərar verir
- Escape halları: funksiyadan çıxan pointer, xarici dəyişənə göndərilən
  pointer, channel-a göndərilən pointer, interface{}-yə boxing, closure yaxalaması
- `go build -gcflags="-m"` — escape səbəblərini göstərir

### #96: Allocation azaltma
- `sync.Pool` — 1024-bayt buferlərin təkrar istifadəsi:
```go
var pool = sync.Pool{New: func() any { return make([]byte, 1024) }}
buffer := pool.Get().([]byte)
defer pool.Put(buffer)  // (düzgün idarə ilə)
```
- Compiler map optimizasiyası; API dəyişikliyi (təkrar istifadə oluna bilən
  bufeyr parametri)

### #97: Inlining
- `-gcflags "-m=2"` — inline qərarları; kiçik funksiyalar inline olur
- **Fast-path inlining:** sync.Mutex.Lock — fast path inline + slow path ayrı
  (mid-stack inlining müasir Go-da)

### #98: Diagnostics tooling
- **pprof:** CPU (100Hz sampling, 30s default), heap (alloc_space/objects),
  goroutine, mutex (contention), block
- Tek profiler vaxtında; `go tool pprof`; execution tracer (GC/scheduler
  detalları — "white spaces" = boş CPU vaxtı)

### #99: GC necə işləyir
- Mark (root-dan traversal) + Sweep; 2 stop-the-world (qısa) + concurrent mark
- **GOGC=100 default:** heap 2× olanda GC tetiklenir; GOGC=off → limitsiz;
  artır (200, 300) → GC tezliyi azalır, yaddaş artır
- Peak-forcing hilesi: böyük bir alokasiya ilə GOGC limitini süni yüksəlt

### #100: Docker/K8s CFS quota
- Goroutine-lər OS thread-lərində; container CPU limiti (məs. 4 core) — amma
  **GOMAXPROCS host nüvə sayını görür** (məs. 8) → 100ms pəncərədə 400ms
  quota 50ms-də bitir → CFS THROTTLE → latency pikləri
- Həll: GOMAXPROCS-u container limitinə bərabərləşdir (runtime.GOMAXPROCS və ya
  automaxprocs lib); Go 1.25+ avtomatik container limitini tanıyır

## Əsas terminlər

- Mechanical Sympathy (maşınla simpatiya)
- Cache Line (64 bayt) / Locality (spatial/temporal)
- Critical Stride / Conflict Miss
- False Sharing / Padding
- Data Hazard (RAW)
- Escape Analysis
- sync.Pool
- Fast-Path Inlining
- pprof / Execution Tracer
- GOGC
- CFS Quota Throttling

## Praktik nəticə

- Ölçmədən optimallaşdırma; pprof + benchmark + benchstat zənciri
- Strukturları böyükdən kiçiyə; əlaqəli datanı bir yerdə
- Pool + slicing + escape biliyi = GC təzyiqi azalır
- Container-da GOMAXPROCS = limit; -race CI-də, prod-da YOX

## Mənbə

Pages: 282-375 (Chapters 11-12, 100 Go Mistakes 2nd ed.)
