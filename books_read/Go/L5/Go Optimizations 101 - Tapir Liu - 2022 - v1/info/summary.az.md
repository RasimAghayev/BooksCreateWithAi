# Go Optimizations 101 — Xülasə (Azərbaycanca)

> **Kitab:** Go Optimizations 101 — Tapir Liu, go101.org, 2022 · 161 səh. · 🧠 Expert (5/5)
> **Səviyyə əsası:** Compiler/runtime daxili mexaniklərini (v1.19) bilən oxucu üçün; hər fənd benchmark sübutu ilə.

---

## Kitabın məqsədi və yanaşması

Kod performansı optimizasiya toplusu — amma "fənd siyahısı" deyil, **mexanizm anlayışı**: hər məsləhət compiler-in/runtime-un NİYƏ belə davrandığı ilə əsaslanır. Əsas tezislər:

1. **Trade-off dünyası:** Oxunaqlıq > performans; kodun 90%-i optimallaşdırma tələb etmir — hot path-i tap.
2. **Ölçməsiz optimallaşdırma yoxdur:** kitab boyu `-m`, `-d=ssa/check_bce`, `gctrace`, benchmark.
3. **Versiya həssaslığı:** v1.19 əsası; compiler boşluqları qeyd olunur və "gelecekte düzələ bilər" xəbərdarlıqları var.

## Fəsl-fəsl xülasə

### Ch 2 — Value parts & sizes
Direct/indirect hissə modeli; tip ölçü cədvəli (string=2 word, slice=3, map/chan/func=1); alignment; **padding: sahə sırası 24B→16B**; ≤4 word xüsusi kopya (9→10 sahə sıçrayışı!); `for _, v` böyük elementdə 4× yavaş.

### Ch 3 — Allocations
Blok/size-class/page modeli (israf nümunələri); pre-allocate `make(0,n)` 2.4×; in-place filter 8×; blok birləşdirmə (101→2 alloc); pool (custom > sync.Pool müəllif üçün).

### Ch 4 — Stack & escape analysis
Frame compile-da hesablanır → stack allocation pulsuz; escape halları (loop dəyişəni, interface arqument, **reflect/fmt**, return istinadı); inlining escape-ni yüngülləşdirir; **threshold-lər: 32B string / 64KB new+make / 10MB elan**; array→slice törət 10MB stack; dummy frame ilə stack-i əvvəlcədən böyüt (9× sürət).

### Ch 5 — GC
Scan+mark/sweep; pacer 3 strategiyası; GOGC formulu (v1.18+ roots daxil); **balast** (150MiB virtual) və **GOMEMLIMIT** (v1.19+); paylaşan blok tələsi (1 canlı element bütün bloku saxlayır → kopyala!).

### Ch 6-7 — Pointers & structs
Loop-daxili nil-check/dereference-dən yayın: `_ = *a` (nil-check loop-dan çıxır), lokal yığ `n := *sum` (5×, amma aliasing yoxdursa); struct sahəsi eyni prinsip.

### Ch 8 — Arrays & slices
Böyük array literal müqayisəsindən qaç; **slice→array-pointer kopya ≤64B (2-3×)**; make+copy optimizasiyasının 4 şərti; append böyümə alqoritmi (1.18 monoton); clone=make+copy; merge/insert idiomları; range 2-ci dəyişən hər yerdə yavaş; **memclr for-range sıfırlama**; 3-indeks subslice.

### Ch 9 — Strings
**5 allocation-sız çevirmə halları** (range, müqayisə, map-oxu, sabitli concat, ≤32B); `m[string(k)]++`Allocate edir → `map[K]*V` və ya index-table; concat: `+` / Builder+Grow / ≤64B byte slice yolu; `EqualFold` 5.6×; `[2]string` açar 3.5×; BytesWriter/bufio.

### Ch 10 — BCE
check_bce diaqnostikası; avtomatik hallar; **ən böyük indeks əvvəl**; hint = slice YENİDƏN təyin; qlobal slice→lokal; array slice-dən BCE-dost; hələ düzəlməyən boşluqlar.

### Ch 11 — Maps
Backing array azalmır; `m[k]++` 43% ; **pointer-siz açar/element → GC scan yox** (`[32]byte`); söz-sayğac 3 variantı (index-table uzunmüddət ən yaxşı); pre-allocate; **bool/enum map 11× yavaş → `[2]T` index table**.

### Ch 12 — Channels
NoSync < atomic < mutex < **channel (27×)**; case-say select xərci; **struct-elementli vahid kanal** (1295→851ns); try-send/receive 5ns (xüsusi).

### Ch 13 — Functions
Inline cost/budget (80); inline POZANLAR: rekursiya, recover, defer, go, funksiya-dəyişən; **hot/cold ayırma + `//go:noinline` soyuq üçün**; manual > auto (2×); ≤4 word value-param, böyük pointer; named/unnamed result sabit deyil — benchmark; **defer loop-da 10×** → anonim funksiya; `cond && debugPrint(args)` short-circuit; hot-filialda escape arqumenti ayır.

### Ch 14 — Interfaces
Boxing xərc xəritəsi: pointer/sabit/bool/int8 ~1.2ns; **[0,255] int ~3.4ns**; kənar int/float 20×; **string/slice 50×**; array 592ns!; lookup-table 20× ucuz; **iki interfaceə: y=x assign (1 box)**; fmt çox-eyni-dəyər: əvvəlcədən box; vtable+inline-itkisi 8×; **RGBA64Image dərsi: sıx API-da konkret tiplər**.

## Kitabın əsas mesajları

1. **Compiler-in dilini öyrən** — escape, inline, BCE, size-class: performans = bu 4 mexanizmi oxuya bilmək.
2. **Yaddaş iyerarxiyası hər şeydir** — stack < heap; register < memory yox; pointer az = GC xoşbəxt.
3. **İdiomlar var, muydlar yox** — hər qaydanın şərti, istisnası və versiya-dependentliyi var.
4. **Benchmark qanundur** — bütün fəndlər təkrarlanan sübutlara əsaslanır; fərziyyə deyil.
5. **Hot path-i tap** — qalan kodu oxunaqlı saxla: kitabın öz fəlsəfəsi.

## Kitabdan sonra

- `pprof` dərinləşməsi (kitab yalnız benchmarkə toxunur)
- `runtime` paketinin oxu (GOMAXPROCS, scheduler interaksiyası)
- Go 1.20-1.23 dəyişiklikləri (PGO — Profile-Guided Optimization!)
- Müəllifin digər kitabı: Go 101 (dil əsasları)

## ⚠️ Müəllifin öz qeydləri / düzəlişlər
- Kitab v1.19 əsasıdır — bəzi "boşluqlar" (compiler flaws) sonrakı versiyalarda düzəlib; hər fəndi öz versiyanda yoxla
- Bəzi fəndlər (composite literal stack, split-concat) **kənar/künklü davranışlardır** — qeyri-rəsmi, gələcəkdə dəyişə bilər
