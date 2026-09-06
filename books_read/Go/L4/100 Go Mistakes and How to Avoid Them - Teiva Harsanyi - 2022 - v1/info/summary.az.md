# 100 Go Mistakes and How to Avoid Them — Xülasə (Azərbaycanca)

> **Kitab:** 100 Go Mistakes and How to Avoid Them — Teiva Harsanyi, Manning, 2022 (ISBN 9781617299599, 364 səh., 12 chapter)
> **Səviyyə:** 🚀 Advanced (4/5) · **Dil:** Go
> **Struktur:** 100 müstəqil səhv — hər biri ssenari + yanlış kod + həll + müzakirə formatında. Tək layihə YOXDUR.

---

## Kitabın məqsədi və yanaşması

Müəllif 2019-cu ildəki "Top 10 Go mistakes" blog yazısından yola çıxaraq 2 ildə 100 səhvi real-world nümunələrlə topladı. Metodologiya: neyroelm sübut edir ki, beyin səhvlər zamanı ən çox böyüyür (Moser 2011); Metcalfe (2017) — səhvlərin **facilitative effect**-i: xətanı yalnız özünü deyil, konteksti də xatırlayırıq. Hər səhv concrete misal + çoxhallı həll + "context is key" müzakirəsi ilə təqdim olunur.

## Chapter-by-chapter xülasə

### Ch1 — Go: Simple to learn but hard to master
Go 4 cəhətdən güclüdür: stability, expressivity, sürətli compilation, safety. Amma **simple ≠ easy** — 2019 ASPLOS tədqiqatı (Docker/gRPC/K8s): blocking bug-ların əksəriyyəti channel istifadəsindən. 7 səhv kateqoriyası: bugs, needless complexity, weaker readability, suboptimal organization, lack of API convenience, under-optimized code, lack of productivity.

### Ch2 — Code and project organization (#1-#16)
Variable shadowing (`:=` daxili blokda yeni dəyişən) — 2 həll: müvəqqəti var və ya `=`+əvvəlcedən elan. Nested code — happy path solda, else burax. init yalnız uğura bilən statik konfiq üçün (DB → adi `createX() error`). Interface-lər **kəşf olunur** — consumer tərəfdə, minimal; qaytar struct, qəbul et interface. `any` ifadəsizdir. Generics 3 halda (data strukturlar, slice/map/chan, davranış factor-out). Embed-də 2 qadağa (şəkər/gizlətmə). Functional options pattern (`WithPort`). project-layout + meaningful paket adları (`stringset` YOX `util`). Hər exported element dokumentli. Linter-lər (vet/errcheck/golangci-lint) avtomatik.

### Ch3 — Data types (#17-#29)
`010`=8 (oktal → `0o644` yaz). Overflow silent — MaxInt yoxlamaları. Float yaxınlaşmadır: == YOX, vurma əvvəl, böyüklük-qruplu toplama. Slice = ptr+len+cap — append ortaq array yazır; `make(len/cap)` ilkinləşdirmə 400%; nil vs empty (allocation + JSON null/[]); `len(s) != 0` yoxlama; `copy` min(len); full slice expression `s[:2:2]`; cap leak → kopya; pointerli elementləri nil-lə. Map yalnız böyüyür (pik yaddaş qalır) → re-create / `map[int]*[128]byte` (38MB vs 293MB). `==` comparable-lərdə; `errors.As/Is` + DeepEqual(test)/custom equal(run-time).

### Ch4 — Control structures (#30-#35)
Range value = KOPYA → indekslə mutasiya. Range ifadəsi 1 dəfə qiymətlənir (channel dəyişmək təsirsiz; `range &a` kopyasız). Loop dəyişəni tək sabit ünvanlıdır — pointer saxlama (`current := customer`). Map sırası unspecified; insert-iterasiya → kopya üzərində yaz. `break` innermost → label (`break loop` — idiomatik). Loop-da defer yığılır → köməkçi funksiya.

### Ch5 — Strings (#36-#41)
Rune = int32 code point; `len`=bayt; `utf8.RuneCountInString`=rune. `for i, r := range s` — r rune-un özü (i bayt başlanğıcı). TrimRight/Left=SET, TrimSuffix/Prefix=affiks. 5+ birləşdirmə → `strings.Builder`+`Grow` (99% qazanc). I/O-da `bytes` paketi (konversiya=kopya). Substring ortaq array yaşadır → `string([]byte(...))` / `strings.Clone`.

### Ch6 — Functions and methods (#42-#47)
Receiver qərarı: mutasiya/sync sahəsi → pointer (mütləq); map/func/chan → value (mütləq); şübhə → pointer. Named resultlar oxunaqlıq (`lat, lng`) + zero-value tələsi (`return 0,0,err` err hələ nil!). **Nil pointer → interface = non-nil** (MultiError tələsi) — sonda `return nil`. `io.Reader` YOX filename (test + reuse). Defer arqumentləri DƏRHAL → pointer/closure.

### Ch7 — Error management (#48-#54)
Panic yalnız programmer error + məcburi asılılıq. `%w`=wrap (source açıq — coupling), `%v`=transform. `errors.As(err, &T{})` / `errors.Is(err, sentinel)` — ==/type-switch wrap-lə işləmir. Expected → sentinel, unexpected → custom tip. **Bir xəta 1 dəfə idarə olunar** — log YAXUD return (context-i wrap-la). İqnor = `_ =` + səbəb kommenti. Defer Close xətaları: log / named err + closeErr prioritizasiyası.

### Ch8 — Concurrency: Foundations (#55-#60)
Concurrency=STRUKTUR (kofe dükanı), Parallelism=İCRA. G/M/P + GOMAXPROCS + work stealing + preemptive (1.14+). **Konkurrentlik həmişə sürətli deyil** — merge sort V1 8x yavaş (miniklik işlərdə goroutine xərci); threshold (2048, benchmark ilə) → 40% qazanc. Paralel goroutine → mutex; konkurrent → channel. Data race (yaddaş+yazma; atomic/mutex/channel müdafiəsi) ≠ race condition (sıra). Memory model zəmanətləri: go<icra; unbuffered receive<send-tamamlanma (buffered-da YOX). Worker pool: CPU→GOMAXPROCS, IO→xarici sistem. Context: deadline/cancel/values — `defer cancel()`, unexported key tipi, select+Done ilə blocking qadağası.

### Ch9 — Concurrency: Practice (#61-#74)
HTTP context response yazılanca ölür → **detach context**. Goroutine exit planı şərt (leak = 2KB+resurs). Loop-closure tələsi (`val := i` / parametr). **Select RANDOM seçir** — prioritet üçün inner select+default. `chan struct{}` siqnalı. Nil channel = select case söndürücü (merge pattern). Buffered ölçü=1 default; sinxronizasiya=unbuffered. `%v` gizli Stringer çağırır (etcd race, deadlock). Ortak slice append=rice. Map/slice assign=shallow — tam lock / dərin kopya. WaitGroup Add-valideyn-goroutine-əvvəl. Təkrar broadcast → sync.Cond. errgroup = paralel+error+ctx. **Sync tiplər kopyalanmaz** (pointer receiver / `*sync.Mutex`).

### Ch10 — The standard library (#75-#81)
`time.Second` API (1000 = mikrosaniyə tələsi). Loop-da `time.After` = 1GB leak mümkün → `time.NewTimer`+Reset. Embedded time.Time json.Marshaler-i promote edir (ID itir) → adlandır. Monotonic `m=+` == pozur → `.Equal`/`Truncate(0)`. `map[string]any` → numeric=float64. sql.Open bağlantı açmır → `Ping`; pool 4 parametri; Prepare; NULL → `*string`/`sql.NullString`; `rows.Err()` şərt. `io.Closer` hamısı bağlanır (HTTP body oxunmasa belə — `io.Discard` keep-alive üçün; yazıla bilən fayl close-err propagate/Sync). `http.Error` saxlamır → `return`. **Default client/server qadağan** — 4 client timeout + MaxIdleConnsPerHost(2!); server ReadHeader/Read/TimeoutHandler/Idle.

### Ch11 — Testing (#82-#90)
Kateqoriya: build tags / env+Skip (aşkar) / `-short`. **`-race` məcburi** (vector-clock; false positive yox; loop təkrarı). `-parallel` / `-shuffle=seed`. Table-driven (map+subtests; paraleldə `tt := tt`). **Sleep=flaky** → channel mock (deterministik) > retry assert. `type now func() time.Time` inyeksiyası və ya client ötürməsi. httptest (Recorder+Request / NewServer); iotest (TestReader, TimeoutReader). Benchmark 4 tələsi: timer reset, micro+benchstat, **local→global** (inline qarşı), **yeni data** (observer effect). Coverage `-coverpkg`; `x_test` paketi; utility `t.Fatal`; `t.Cleanup`/`TestMain`.

### Ch12 — Optimizations (#91-#100)
CPU cache: L1 ~1ns, RAM 50-100x; 64B line; **struct-of-slices > slice-of-structs**; linked list non-unit stride (70% yavaş); **critical stride** (512 int64 = 4KB L1D). **False sharing** → padding `_ [56]byte` (40%). **ILP**: `v := s[0]` ilə hazard-parçalama (20%). **Alignment**: sahələr böyükdən kiçiyə (24→16 bayt, 15% iterasiya). **Stack/heap**: sharing up=heap (pointer return 12x yavaş + GC təzyiqi); sharing down=stack; `-gcflags -m=2`. Allocation azaltma: sharing-down API, `m[string(bytes)]` optimizasiyası, **sync.Pool**. **Fast-path inlining** (Mutex.Lock pattern, budget 80). **Diagnostics**: pprof 5 profil + heap diff leak; execution tracer paralellik üçün (merge sort 50% runtime overhead sübutu). **GOGC**: heap 2x → GC; pikdə artır; virtual heap hilesi (`make(1GB)` mmap). **K8s**: GOMAXPROCS host-a görə → CFS throttling (300% latency cəzası) → **automaxprocs**.

---

## Kitabın əsas mesajları

1. **Simple ≠ easy** — Go-nun sadəliyi mastery-ni qarantmirəm; 100 səhv kontekstlə birlikdə yaddaşda qalır.
2. **Abstraksiyalar kəşf olunur** — interface/generics/type-embedding yalnız konkret ehtiyacda; premature abstraksiya = needless complexity.
3. **Semantika birinci, performans sonra** — "Make it correct, make it clear, make it concise, make it fast, in that order."
4. **Konkurrentlik = strukturlaşdırılmış mürəkkəblik** — random select, memory model, false sharing, GOMAXPROCS kimi dərin anlayışlar səhvlərin mənbəyidir.
5. **Yaddaş modelini başa düş** — slice/map-in pointer təbiəti (shallow copy, append race, substring leak) 10+ səhvin kökündədir.
6. **Ölçmə olmadan optimallaşdırma YOX** — benchmark intizamı (timer, benchstat, observer effect) + pprof/tracer şərtdir.
7. **Environment vacibdir** — production client/server, K8s CFS limitləri, GC təzyiqi: default-lar production üçün düzülməyib.

## Kitabdan sonra

- Go 1.22+ dəyişiklikləri (loop var semantikası, safelink, PGO)
- Go issue 33803 (CFS-aware GOMAXPROCS) izləməsi
- benchmark alətləri (benchstat, perflock) praktikası
- runtime/trace user-level task-lər
