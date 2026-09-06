# Go Optimizations 101 — Terminologiya (Azərbaycanca)

> Texniki terminlər `English (Azərbaycanca qarşılığı)` formatında — kitabın ardıcıllığı ilə.

## Value / Memory (Ch 2-4)
- **Value part (dəyər hissəsi)** — dəyərin yaddaşda hissəsi: direct (birbaşa) / indirect (dolayı)
- **Memory alignment (yaddaş düzləndirilməsi)** — ünvanın N qatına salınması; N = alignment guarantee
- **Struct padding (struktur doldurulması)** — sahələr arası boş baytlar; ölçüyə daxildir
- **Small-size type (kiçik ölçülü tip)** — ≤4 native word; xüsusi kopyalama optimizasiyası
- **Memory block (yaddaş bloku)** — kəsintisiz seqment; ayrılmanın vahidi
- **Size class (ölçü sinfi)** — əvvəlcədən müəyyən blok ölçüləri (8, 16, 24, 32, 48...)
- **Memory page (yaddaş səhifəsi)** — 8192 bayt; 32KB+ blokların vahidi
- **Stack frame (stek kadrı)** — funksiya çağırışının stack seqmenti; ölçüsü compile-da hesablanır
- **Contiguous stack (fasiləsiz stek)** — bütöv seqment kimi goroutine stack-i
- **Escape analysis (kaçış analizi)** — stack/heap qərarı verən compiler modulu
- **Escape to heap (heap-ə qaçma)** — dəyərin heap-də ayrılması
- **De-virtualization (de-virtuallaşdırma)** — interface çağırışının konkret tipə çevrilməsi
- **Inlining (içinə yerləşdirmə)** — çağırışın çağırılan kodla əvəzlənməsi
- **Inline cost / budget** — statement xərcləri cəmi / hədd (80)

## GC (Ch 5)
- **GC cycle (GC dövrü)** — scan+mark (işarələmə) + sweep (təmizləmə)
- **GC pacer (GC tempoluğu)** — cycle başlanğıc qərarları mexanizmi
- **GOGC** — hədəf heap faizi (default 100); env və `SetGCPercent`
- **Live heap (canlı yaddaş)** — scan sonundakı qeyri-zibil heap
- **GC roots** — stack + qlobal scannable hissələr (v1.18+ formulda)
- **Memory ballast (yaddaş balastı)** — süni live heap saxlayan böyük virtual blok
- **GOMEMLIMIT** — v1.19 memory limit strategiyası; soft limit
- **Finalizer** — obyekt yığılanda çağırılan hook (`SetFinalizer`)
- **runtime.KeepAlive** — istinadın canlı qalmasını təmin edən alət

## Pointer / Loop (Ch 6-7)
- **Nil pointer check (nil yoxlaması)** — TESTB instruksiyası; loop-dan çıxarılmalı
- **Pointer dereference (istiqamətləndirmə)** — memory oxu/yaz; register-dən bahalı
- **Aliasing (üst-üstə düşmə)** — iki pointer eyni yaddaşa işarə edir; transformasiya təhlükəsi
- **Register vs memory** — CPU işləmə rejimləri; register həmişə üstün

## Slice / String (Ch 8-9)
- **make+copy optimization** — make-in lazımsız sıfırlamasını atan compiler transformasiyası
- **Full slice expression (tam dilim ifadəsi)** — `s[a:b:c]` üç-indeks; cap kəsmə
- **memclr** — daxili vectorized sıfırlama (for-range sıfırlama pattern-i)
- **Slice growth (dilim böyüməsi)** — append alqoritmi: <256 → 2×; sonrası ~1.25×
- **Backing array (dəstəkləyici massiv)** — slice/map-in arxasındakı massiv
- **Immutable (dəyişməz)** — string-in bayt dəyişməzliyi; çevirmə dublikasiya qaynağı
- **strings.Builder** — artımlı string qurucu; `Grow(n)` ilə pre-allocation
- **strings.EqualFold** — allocation-sız case-insensitive müqayisə
- **Three-way comparison (üçtərəfli müqayisə)** — -1/0/+1 nəticəsi

## BCE (Ch 10)
- **Bounds checking (sərhəd yoxlaması)** — IsInBounds / IsSliceInBounds
- **BCE (Bound Check Elimination)** — zəruri olmayan yoxlamaların silinməsi
- **Compiler hint (kompilyator işarəsi)** — slice yenidən təyini / redundan if
- **check_bce** — `go run -gcflags="-d=ssa/check_bce"` diaqnostikası
- **Index table (indeks cədvəli)** — yoxlamasız array sorğu mexanizmi

## Map / Channel (Ch 11-12)
- **Rehash (yenidən hashləmə)** — map böyüməsində entry köçürməsi
- **Pointer-free container** — açar/element-siz konteyner; GC scan-dan azad
- **L-value map index** — `m[k] = ...` (allocate); R-value oxu — etmir
- **Try-send / try-receive** — `select+default` non-blocking əməliyyat
- **Open-coded defer (v1.14+)** — loop-xarici defer optimizasiyası
- **Short-circuit evaluation** — `&&` ilə arqument qiymətləndirməsinin kəsilməsi

## Interface (Ch 14)
- **Boxing / Unboxing (qutulama/açma)** — interface təyinatı / type assertion
- **Virtual table (vtable)** — dinamik metod cədvəli; çağırış xərci ~2.5ns
- **Staticbytes cache** — [0,255] kiçik int boxing buferi
- **Lookup table (axtarış cədvəli)** — dəyər→pointer qlobal xəritəsi (ucuz boxing)
- **image/draw.RGBA64Image** — konkret-tip API dizayn nümunəsi (Go 1.17+)

## Benchmark alətləri
- **ns/op, B/op, allocs/op** — benchmark vahidləri
- **testing.AllocsPerRun** — funksiya başına allocation sayı
- **AllocedBytesPerOp** — ayrılan bayt ölçüsü
- **b.ResetTimer** — hazırlıq vaxtını xaric etmə
- **-gcflags=-m / -m -m** — escape / inline + cost
- **-gcflags=-S** — assembly çıxışı
- **-d=ssa/check_bce** — BCE qalıqları
- **GODEBUG=gctrace=1** — GC cycle jurnalı
