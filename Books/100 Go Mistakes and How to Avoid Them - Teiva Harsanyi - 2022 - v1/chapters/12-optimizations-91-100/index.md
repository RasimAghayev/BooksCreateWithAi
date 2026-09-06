# Chapter 12 — Optimizations (#91-#100)

## Bu chapter nədən bəhs edir?

Bu son chapter optimallaşdırmanın 10 səhvini əhatə edir: CPU cache-ləri (cache line, spatial locality, striding, cache placement policy), false sharing, instruction-level parallelism, data alignment, stack vs heap + escape analysis, allocation azaldılması (API, compiler optimizasiyaları, sync.Pool), inlining (fast-path), diaqnostika alətləri (pprof profiling, execution tracer), GC-nin işləməsi (GOGC) və Docker/Kubernetes-də Go (CFS throttling, automaxprocs).

> **Disclaimer (kitabın özü):** *"Make it correct, make it clear, make it concise, make it fast — in that order."* (Wes Dyer) — Oxunaqlıq optimallıqdan öncə gəlir; bu chapter yalnız ehtiyac olan kod yolları üçündür.

---

## Əsas fikirlər

### #91: Not understanding CPU caches (CPU cache-ləri)

**Mechanical sympathy** (Jackie Stewart, F1): sistemin dizaynını anlayıb ona uyğun işləmək → optimal performans.

**Cache iyerarxiyası (Intel i5-7300 nümunəsi):**

| Səviyyə | Ölçü | Gecikmə |
|---------|------|---------|
| L1 (D=32KB + I=32KB, per core) | 64 KB | ~1 ns |
| L2 (per core) | 256 KB | ~4x L1 |
| L3 (shared, off-die) | 4 MB | ~10x L1 |
| Main memory (RAM) | GB-lər | **50-100x L1** |

L1-də 100 dəyişən oxumaq = RAM-da 1 oxumaq.

**Cache line:** CPU tək dəyişən YOX — **64 baytlıq kontiqental blok** kopyalayır (8 × int64). Səbəb: locality of reference — temporal (təkrar çıxış) + spatial (qonşu çıxış).

**sum2 vs sum8 (hər 2-ci/hər 8-ci element):** Intuisiya 4x fərqi deyirdi → real fərq ~10%. Cache line-dakı hit-lər dominating.

**Slice of structs vs struct of slices:**

```go
// sumFoo: [a b a b a b ...] — hər cache line-da YARISI lazımsız
func sumFoo(foos []Foo) int64 { /* Foo{a, b int64} */ }
// sumBar: [a a a a ... b b b b ...] — cache line-lərin hamısı DOLU
func sumBar(bar Bar) int64  { /* Bar{a, b []int64} */ }
```

sumBar ~20% sürətli — daha az cache line fetch.

**Striding növləri:**

| Stride | Nümunə | Predictable? |
|--------|--------|--------------|
| **Unit** | Slice kontigu | Bəli — ən effektiv |
| **Constant** | Hər 2-ci element | Bəli — amma daha çox cache line |
| **Non-unit** | Linked list, pointer slice | XEYR — CPU cache line fetch ETMİR |

Linked list (kontigu alloc olunsa belə!) slice-dən ~70% yavaş — CPU stride-i proqnozlaşdıra bilmir.

**Cache placement policy (512 vs 513 matrix sirri):**

- **Set-associative cache:** cache set-lərə bölünür; blok → set = adresin set index bitləri ilə.
- Bütün istifadə olunan bloklar eyni set-ə düşərsə → set dolur → **conflict miss** (dəyişdirilmiş bloklar itir).
- **Critical stride:** set index-i təkrarlanan addım — 32KB 8-way L1D-də: 64 set × 64 bayt = **4 KB = 512 int64** → 512 kolonlu matrix critical stride-ə DÜŞÜR, 513 düşMİR (yerdəyməş 50% fərq).
- Müalicə: critical stride-lərdən qaç; micro-benchmark nəticələrini production cache arxitekturasına ehtiyatla köçür.

---

### #92: Writing concurrent code that leads to false sharing

**Ssenari:** İki goroutine — biri `result.sumA`, digəri `result.sumB` artırır (data race YOXDUR — ayrı dəyişənlər!).

**Problem:** sumA + sumB eyni 64-baytlıq cache line-da (7/8 ehtimal) → hər iki core həmin line-ın kopyasını saxlayır → **cache coherency (MESI protokolu): bir core yazanda BÜTÜN line digər core-lərdə INVALIDASIYA olunur** — hər iki dəyişən müstəqil olsa belə!

**"Sharing memory is an illusion"** — ən aşağı səviyyədə yoxdur; yazma qranulluğu dəyişən deyil, cache line-dır.

**Həllər:**

```go
// 1. Padding (56 bayt ayrıcı):
type Result struct {
    sumA int64
    _    [56]byte    // sumA-nı ayrı cache line-a məcbur edir
    sumB int64
}
// ~40% sürət artımı!

// 2. Alqoritm dəyişikliyi: hər goroutine LOKAL nəticə + channel ilə ötürmə.
```

---

### #93: Not taking into account instruction-level parallelism (ILP)

**Superscalar processor:** Tək core-da bir neçə instraksiyanı PARALEL icra edir: `total = max(t(I1), t(I2), t(I3))` (sequensial `t(I1)+t(I2)+t(I3)` yerinə).

**Hazard növləri:**
- **Control hazard** (branching) → **branch prediction** həll edir (yanlış proqnoz: 10-20 dövrü cəza ilə pipeline flush).
- **Data hazard** — I2 I1-in nəticəsindən asılıdırsa → paralel icra MÜMKÜNSÜZ; forwarding yalnız yumşaldır.
- Structural hazard — resurs konflikti (bizə aid deyil).

**Nümunə — data hazard azaldılması:**

```go
// V1 — çoxlu data hazard:
s[0]++                    // read s[0] → add → write
if s[0]%2 == 0 {          // read s[0] YENİDƏN (hazard!)
    s[1]++
}

// V2 — müvəqqəti dəyişən (~20% sürətli!):
v := s[0]                 // tək oxu
s[0] = v + 1
if v%2 != 0 {             // v artıq asılı deyil — PARALEL mümkün
    s[1]++
}
```

V2-də "add to s[0]" və "check v" hər ikisi "read s[0]"-dan asılıdır → 3 paralel yol (V1-də 2). Eyni addım sayı, azalan critical path.

**Ehtiyat:** Go kompilyator inkişaf etdikcə generasiya olunan assembly dəyişə bilər — belə micro-optimizasiyalar version bağlıdır.

---

### #94: Not being aware of data alignment

**Data alignment:** dəyişənin adresi ÖZ ölçüsünün multiple-u olmalıdır (int32 → 4-ün multiple-u).

| Tip | Alignment |
|-----|-----------|
| byte/int8/uint8 | 1 |
| int16/uint16 | 2 |
| int32/uint32/float32 | 4 |
| int64/uint64/float64/complex64 | 8 |
| complex128 | 16 |

**Padding problemi:**

```go
type Foo struct {
    b1 byte      // 0x00
                 // 7 bayt PADDING (i 8-ə align olmalı)
    i  int64     // 0x08
    b2 byte      // 0x10
                 // 7 bayt PADDING (struct word-ə tamamlanmalı)
}
// Ümumi: 24 bayt — 10 bayt data + 14 bayt PADDING!
```

**Həll — sahələri ölçüyə görə azalan sırayla:**

```go
type Foo struct {
    i  int64     // 8 bayt
    b1 byte      // + b2 eyni word-də
    b2 byte
}
// Ümumi: 16 bayt — 33% yaddaş qənaəti!
```

**Təsirlər:** Retained strukturlarda artıq yaddaş; tez yaranan heap strukturlarında daha sıx GC; slice iterasiyasında hər cache line 33% daha çox faydalı data → ~15% iterasiya sürəti.

**Qaynaq-qayda:** struct sahələrini ölçü üzrə azalan düz.

---

### #95: Not understanding stack vs. heap

**Stack:** Hər goroutine öz stackına sahib (başlanğıc 2 KB, böyüyüb kiçilir, həmişə kontigu). Funksiya çağrısı → stack frame; return → frame "ölür" (silinmə YOXDUR — yeni variable-lar üstünə yazılır) → **self-cleaning, GC laimsız**.

**Heap:** Bütün goroutine-lərin paylaşılan yaddaş hovuzu → GC tələb edir.

**Escape qaydaları:**

| Hal | Nəticə |
|-----|--------|
| **Sharing down** (funksiya pointer QƏBUL edir) | Stack-da qalır |
| **Sharing up** (funksiya pointer QAYTARIR) | **Heap-ə escape** |
| Qlobal dəyişənlər | Heap |
| Channel-a pointer/value göndərmək | Heap |
| Çox böyük lokal variable | Heap |
| Dəyişən ölçülü (`make([]int, n)`) | Heap (`make([]int, 10)` ola bilmez) |
| Append ilə backing array realloсasiyası | Heap |

**Benchmark fərqi:**

| Funksiya | ns/op | allocs |
|----------|-------|---------|
| `sumValue` (int qaytarır) | 1.26 | 0 |
| `sumPtr` (*int qaytarır) | **14.84** — ~12x yavaş | 1 alloc, 8 B/op |

Heap idarəsi data-intensive tətbiqlərdə CPU vaxtının 20-30%-ni yeyə bilər.

**Moral:** "Kopyadan qaçmaq üçün pointer qaytarma" — premature optimization; semantika birinci (pointer = paylaşım niyyəti).

**Yoxlama:** `go build -gcflags "-m=2"` → `./main.go:12:2: z escapes to heap`.

---

### #96: Not knowing how to reduce allocations

**Əvvəlki chapter-lərdən:** strings.Builder (#39), []byte işləmə (#40), preallocation (#21, #27), alignment (#94).

**3 yeni yanaşma:**

**1. API dizaynı (sharing-down):**

```go
// PİS dizayn — qaytarılan slice həmişə heap-ə escape:
type Reader interface {
    Read(n int) (p []byte, err error)
}
// YAXŞI (real io.Reader) — caller-in buferi:
type Reader interface {
    Read(p []byte) (n int, err error)
}
```

**2. Kompilyator optimizasiyaları:**

```go
// V1 — açar dəyişənə çevrilir (allocation + kopya):
key := string(bytes)
v, contains = c.m[key]

// V2 — eyni görünür amma KOMPİLYATOR []byte→string konversiyasını ATLAYIR:
v, contains = c.m[string(bytes)]   // map lookup üçün xüsusi optimizasiya!
```

**3. sync.Pool:**

```go
var pool = sync.Pool{
    New: func() any {
        return make([]byte, 1024)   // factory — boş pool üçün
    },
}

func write(w io.Writer) {
    buffer := pool.Get().([]byte)   // pool-dan al və ya yarat
    buffer = buffer[:0]              // RESET — əvvəlki data təmizlənir
    defer pool.Put(buffer)           // geri qaytar
    getResponse(buffer)              // buffer-i doldur
    _, _ = w.Write(buffer)
}
```

**Xüsusiyyətlər:** Cache deyil (fizilmiş ölçü YOXDUR); hər GC-dən sonra pool boşalır; konkurrent-təhlükəsiz; amortize allocation xərcini azaldır.

---

### #97: Not relying on inlining

**Inlining:** funksiya çağırışının funksiya body-si ilə əvəzlənməsi. 2 fayda: çağırış overhead-i yoxdur; əlavə optimizasiyalar açılır (məs., escape analizi dəyişə bilər).

```bash
$ go build -gcflags "-m=2"
./main.go:10:6: can inline sum with cost 4 as: ...     # budget: 80
# və ya: cannot inline foo: function too complex: cost 84 exceeds budget 80
```

**Mid-stack inlining (Go 1.9+):** Yalnız leaf funksiyalar deyil — içində başqa çağırış olan funksiyalar da inline ola bilər.

**Fast-path inlining pattern (sync.Mutex.Lock-dan real nümunə):**

```go
// Əvvəl: bütün loqika bir funksiyada → inline OLAMAZ (budget aşılır)
func (m *Mutex) Lock() {
    if atomic.CompareAndSwapInt32(&m.state, 0, mutexLocked) {
        // fast path: kilitsiz halda
        return
    }
    // ... 30+ sətir slow path ...
}

// Sonra: slow path ayrı funksiyaya → Lock İNLİNE OLUR (~5% sürət):
func (m *Mutex) Lock() {
    if atomic.CompareAndSwapInt32(&m.state, 0, mutexLocked) {
        return                    // fast path — inline daxilində!
    }
    m.lockSlow()                  // slow path — çağırış qalır
}
func (m *Mutex) lockSlow() { /* mürəkkəb loqika */ }
```

**Prinsip:** Fast path-i kiçik saxla (inline budget daxilində), slow path-i ayrı funksiyaya çıxar.

---

### #98: Not using Go diagnostics tooling (pprof + tracer)

**pprof aktivləşdirmə:**

```go
import _ "net/http/pprof"    // http://host/debug/pprof
```

Production-da təhlükəsizdir (CPU profili yalnız aktiv ediləndə işləyir).

**Profil növləri:**

| Profil | Nə ölçür | Endpoint / Qeyd |
|--------|----------|------------------|
| **CPU** | Harada vaxt sərf olunur | `/debug/pprof/profile` — SIGPROF hər 10ms (30s default); `-cpuprofile` benchmark-da |
| **Heap** | Heap allocation | `/debug/pprof/heap` — hər 512KB-da 1 nümunə; alloc_objects/alloc_space/inuse_objects/inuse_space |
| **Goroutine** | Goroutine stack trace-ləri | `/debug/pprof/goroutine` — leak şübhəsində |
| **Block** | Sinxronizasiya gözləmələri | `/debug/pprof/block` — `runtime.SetBlockProfileRate` ilə aktivləşdir |
| **Mutex** | Mutex contention | `/debug/pprof/mutex` — `runtime.SetMutexProfileFraction` |

**Analiz:** `go tool pprof -http=:8080 <file>` — web UI: call graph, Top, Flame Graph, source-line heat.

**İnsaytlar:** `runtime.mallocgc` çoxluğu → kiçik allocation-ları azalt; channel/mutex vaxtı → contention; syscall.Read/Write çoxluğu → I/O buffering.

**Memory leak izləməsi (heap diff):**

```bash
# 1. /debug/pprof/heap?gc=1 (GC + endload)
# 2. bir neçə dəqiqə gözlə
# 3. təkrar endload
go tool pprof -http=:8080 -diff_base <file2> <file1>   # sabit artım = leak
```

**Full goroutine dump:** `/debug/pprof/goroutine/?debug=2` — deadlock analizi ("blocked for 1420 minutes on chan receive" kimi).

**Qaydalar:** Eyni anda YALNIZ 1 profil (CPU+heap paralel → səhv müşahidə); threadcreate profili 2013-dən pozulub (issue 6104).

**Execution tracer:**

```bash
go test -bench=. -trace=trace.out
go tool trace trace.out    # web UI
```

**Profiler vs tracer:**

| | CPU Profiling | Execution Tracer |
|--|---------------|------------------|
| Metod | Sample-based (10ms) | Hadisə-əsaslı |
| Qranulluq | Funksiya | Goroutine |
| Üstün olduğu yer | Hot path tapma | Zəif paralellik, GC davranışı, goroutine orkestrasiyası |

Merge sort V1-də tracer göstərdi: ~50% CPU vaxtı goroutine spin-up/orkestrasiyaya gedir (ağ boşluqlar) → zəif paralelliyin sübutu. **User-level task-lər:** `runtime/trace` — `trace.NewTask`/`WithRegion` ilə funksiya-səviyyəli bölgü.

---

### #99: Not understanding how the GC works

**Go GC = concurrent mark-and-sweep:**
- **Mark:** heap-dəki obyektlərin istifadəsini işarələ
- **Sweep:** root-lardan istinad ağacını gəz, istinadsız blokları dealloсasiya et
- Hər GC dövründə 2 qısa **stop-the-world** + əsas iş konkurrent → "concurrent"

**Scavenger:** İstifadə pikindən sonra boş yaddaşı OS-ə qaytarır (yavaş proses; təcili üçün `debug.FreeOSMemory()`).

**GOGC (default 100):** Son GC-dən bəri heap neçə FAİZ böyüsən növbəti GC tetiklenir.

- Heap 128 MB idi → GOGC=100 → **256 MB-da** növbəti GC.
- Əlavə: 2 dəqiqədən artıq GC olmayıbsa məcburi GC.

**Ssenari analizi (1M istifadəçi):**

| Ssenari | GC təzyiqi | Hərəkət |
|---------|-----------|---------|
| Gün boyunca tədrici artım | Orta → azalan | GOGC=100 KİFAYƏT |
| 1 saatda kəskin artım | Yüksək — sıx STW-lər, latency artımı | GOGC-i ARTIR |
| Bir neçə SƏNİYƏDƏ zirvə | Kritik — tətbiq təlaşkeşən | Virtual heap hilesi |

**Virtual heap hilesi (bilinən pik üçün):**

```go
var min = make([]byte, 1_000_000_000)   // 1 GB VİRTUAL
```

Heap-i 1 GB-dən başlat → GOGC=100 ilə GC hər 2 GB-da tetiklenir → pik anında GC sayı azalır. **Fiziki yaddaş istehlak etmir:** `make` → `mmap()` — virtual address space; yalnız toxunan səhifələr page fault ilə fiziki olur (ps ilə doğrulanır).

**GOGC tuning:** production load profilindən sonra; azaltmaq → sıx GC (çox təzyiq); artırmaq → yavaş böyümə (amma böyük heap = uzun təmizlik — xətti fayda YOXDUR).

---

### #100: Not understanding the impacts of running Go in Docker and Kubernetes

**Problem:** GOMAXPROCS **hostun** logical core sayına görə təyin olunur — konteyner CPU limitinə görə YOX!

```yaml
# 8-core node-da limit 4000m (4 core):
resources:
  limits:
    cpu: 4000m
# GOMAXPROCS = 8 (host!) — limit 4 deyil
```

**CFS (Completely Fair Scheduler) mexanizmi:**
- `cpu.cfs_period_us` = 100 ms (period)
- `cpu.cfs_quota_us` = 400 ms (4 core × 100 ms)
- Hər 100 ms-də max 400 ms CPU vaxtı istehlak oluna bilər.

**Throttling ssenarisi:** GOMAXPROCS=8 → 8 thread paralel işləyir → hər biri 50 ms → 8×50 = 400 ms quota 50 ms-də DOLUR → **qalan 50 ms üçün tətbiq TAM DONDURULUR (throttle)**.

**Latency təsiri:** 50 ms orta gecikməli servis → **150 ms-ə qədər** (300% cəza!).

**Həll — automaxprocs (Uber):**

```go
import _ "go.uber.org/automaxprocs"   // blank import — GOMAXPROCS-u CFS kvotasına uyğunlaşdırır
// 4000m → GOMAXPROCS = 4 → throttling mümkünsüz
```

Uzunmüddətli həll Go issue 33803 (CFS-aware GOMAXPROCS) — izlənilməli.

---

## Əsas terminlər

| Termin | İzah |
|--------|------|
| Mechanical sympathy | Sistem dizaynına uyğun işləyərək optimal performans |
| L1/L2/L3 cache | Core-yəxın yaddaş səviyyələri (1ns / 4x / 10x; RAM 50-100x) |
| Cache line | 64 baytlıq kopyalanan blok (8×int64) |
| Temporal/spatial locality | Eyni / qonşu yaddaşın təkrar istifadəsi |
| Compulsory miss | İlk dəfə çıxış (qaçılmaz) |
| Unit/constant/non-unit stride | Proqnozlaşdırıla bilən kontigu / addımlı / proqnozsuz (linked list) |
| Critical stride | Eyni cache set-inə düşən addım → conflict miss |
| Set-associative cache | Set-lərə bölünmüş cache; yerləşmə adresin set index-i ilə |
| False sharing | Eyni cache line-da ayrı dəyişənlər → MESI invalidasiyası |
| Padding | Ayrıcı baytlar — dəyişənləri ayrı cache line-a məcbur etmə |
| MESI | Cache coherency protokolu (Modified/Exclusive/Shared/Invalid) |
| ILP / superscalar | Tək core-da paralel instraksiya icrası |
| Branch prediction | Control hazard həlli — yanlış proqnoz = pipeline flush |
| Data hazard | Nəticə asılılığı — paralel icra maneəsi |
| Data alignment | Adres = tipin ölçüsünün multiple-u |
| Struct padding | Sahə sırası yanlış olsa 14 bayt israf (24 vs 16) |
| Stack frame | Funksiyaya məxsus yaddaş intervalı — self-cleaning |
| Escape analysis | Kompilyatorun stack/heap qərarı (sharing down/up) |
| `B/op`, `allocs/op` | Benchmark heap allocasiya göstəriciləri |
| `m[string(bytes)]` | Map lookup-da konversiya-atlama optimizasiyası |
| `sync.Pool` | Obyekt reuse hovuzu — GC-də boşalır, cache deyil |
| Inlining budget (80) | Kopya edilə bilən funksiya mürəkkəblik həddi |
| Mid-stack inlining | Non-leaf funksiyaların inline olması (Go 1.9+) |
| Fast-path inlining | Slow path-i ayıraraq fast path-i inline-a açmaq |
| pprof (CPU/heap/goroutine/block/mutex) | Sample-based profil alətləri |
| Heap diff (`-diff_base`) | İki GC-sonrası heap müqayisəsi — leak izi |
| Execution tracer | Hadisə-əsaslı goroutine vizuallaşdırması |
| Mark-and-sweep | İşarələ + təmizlə — Go GC alqoritmi (konkurrent) |
| GOGC | Heap böyümə faizi → GC tetikleyicisi (default 100) |
| `debug.FreeOSMemory` | Təcili yaddaşın OS-ə qaytarılması |
| Virtual heap hilesi | `make(1GB)` + mmap — GC sıxlığını pikdə azaltma |
| CFS throttling | Kvota aşılanda 50 ms tam dondurma |
| automaxprocs | Uber kitabxanası — GOMAXPROCS-u CFS limitinə uyğunlaşdırır |

---

## Praktik nəticə

1. **Cache line düşün:** RAM L1-dən 50-100x yavaş; 64-baytlıq bloklar gəlir — data-nı elə düz ki, hər line maksimal faydalı olsun (struct-of-slices, kontigu slice-lər, linked list-dən qaçın non-unit stride).
2. **Critical stride-dən qorx:** 512 kolon = 4KB = L1D critical stride; cache arxitekturası maşına bağlı — micro-benchmark nəticələri başqa sistemə ehtiyatla köçür.
3. **False sharing yoxla:** Konkurent yazılan dəyişənlər eyni cache line-da olmasın — padding (`_ [56]byte`) və ya channel kommunikasiyası (~40% qazanc).
4. **ILP üçün data hazard azald:** Müvəqqəti dəyişənlə (`v := s[0]`) oxu-dan-asılılığı parçala — paralel icra yolları artır.
5. **Struct sahələrini böyükdən kiçiyə düz:** 24 bayt → 16 bayt — yaddaş + GC sıxlığı + spatial locality.
6. **Pointer qaytarma = heap:** Sharing up escape; sharing down stack. Semantika birinci — "kopyadan qaçma" ön optimallaşdırmadır. `-gcflags "-m=2"` ilə doğrula.
7. **Allocation azalt:** Sharing-down API (`io.Reader` patterni); `m[string(b)]` map lookup optimizasiyası; tez yaranan eyni tipli obyektlərdə `sync.Pool` (reset `buffer[:0]` + `defer Put`).
8. **Fast-path inlining:** Kiçik fast path + ayrı `slow()` funksiya — sync.Mutex.Lock patterni; `-gcflags` ilə budgeti izlə.
9. **Diaqnostika məcburidir:** pprof (CPU hot path, heap leak diff, goroutine dump, block/mutex contention) + execution tracer (paralellik, GC davranışı) — production-da pprof açıq qala bilər.
10. **GC tuning:** GOGC default 100; kəskin piklərdə artır; böyük bilinən zirvə üçün virtual heap hilesi; `GODEBUG=gctrace=1` ilə izlə.
11. **K8s-də mütləq automaxprocs:** Yoxsa GOMAXPROCS=host cores → CFS throttling → 300% latency cəzası.

---

## Mənbə

- Kitab: *100 Go Mistakes and How to Avoid Them* — Teiva Harsanyi, Manning, 2022 (ISBN 9781617299599)
- Chapter 12: "Optimizations", book səh. 299–354
- PDF səhifələri: 319–374
- İstinadlar: go.dev/doc/diagnostics; Go issue 6104 (threadcreate), issue 33803 (CFS-aware GOMAXPROCS); automaxprocs (github.com/uber-go/automaxprocs); Go 2021 developer survey
