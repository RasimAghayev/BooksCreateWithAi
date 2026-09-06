# Chapter 5 — Garbage Collection (Zibil toplanması)

## Bu chapter nədən bəhs edir?
Go GC-nin iş prinsipinə (scan+mark/sweep), GC pacer stragetiyalarına, GOGC formula-sına, GC təzyiqini azaltma yollarına, memory ballast və GOMEMLIMIT (memory limit strategy) texnikalarına.

## Əsas fikirlər

### 1. GC cycle necə işləyir
**Nədir:** İki fazalı avtomatik yığma — scan+mark (işarələmə) və sweep (təmizləmə).

**Necə işləyir:**
- **Scan+mark:** canlı bilinən value part-lardakı BÜTÜN pointer-lər izlənilir → istinad olunanlar da canlı kimi işarələnir → artıq yeni canlı tapılmayana qədər
- **Sweep:** işarələnməyən heap blokları zibil sayılıb yığılır
- **Korollar:** nə qədər çox pointer, o qədər çox scan işi → GC təzyiqi

**Nəyə lazımdır:** "Pointer az = GC xoşbəxt" prinsipinin mənşəyi.

### 2. GC pacer — 3 avtomatik başlanğıc strategiyası
1. **Heap percentage (GOGC):** scan+mark bitəndə hədəf heap ölçüsü hesablanır; heap hədəfi aşanda yeni cycle
2. **2 dəqiqəlik taymer:** finalizer-lərin işləməsi və stack-lərin kiçilməsi üçün
3. **Memory limit (v1.19+, GOMEMLIMIT):** runtime-un istifadə etdiyi yaddaş limiti aşınca yeni cycle

### 3. Pointer istifadəsinin tarazlığı
**Kitabın balanslı mövqeyi:** Pointer-lər həm GC yükü, həm də **kopya xərcini azaltma alətidir** — mütləq qaçmaq YOX, yalnız zərəri aydın olan yerlərdən qaç:
- Heap allocation azaldıqsa pointer də azalır (hər heap-dəki lokalın stack-də gizli pointeri var)
- `Go C deyil` — bəzi lazımsız allocation-lar qaçılmazdır, qəbul et

### 4. Paylaşılan blok tələsi — bir canlı element bütün bloku saxlayır
**Nədir:** Bir blokda bir neçə value part daşıyanda, blok yalnız hamısı öləndə yığılır.

**Kitabdan kod nümunəsi:**
```go
var s = make([]int, 1000)   // hamısı EYNİ blokda
var p = &s[999]             // 1 element canlı saxlanılır
runtime.GC()                // → BÜTÜN 1000-elementlik blok yığıla BİLMİR!
```
**Həlli:**
```go
var v = s[999]   // kiçik dəyəri KOPYALA
var p = &v       // böyük blok artıq azaddır
```
**Praktik ssenari:** `strings.Fields/Split/Trim` nəticələri arqumentin **substring-i**dir — böyük qısaömürlü string-dən kiçik uzunömürlü substring saxlayırsansa, **kopyala** (`strings.Clone`/`[]byte(...)`), yoxsa 10MB-lıq bloq 1 bayt üçün yığılmayacaq.

### 5. GOGC formula-sı (hədəf heap)
**v1.18-dən:** `Target = Live heap + (Live heap + GC roots) × GOGC/100`
**v1.18-dən əvvəl:** `Target = Live heap + Live heap × GOGC/100`

- GOGC default = 100 (live heap-in 100% böyüməsinə icazə)
- `GOGC=off` / `SetGCPercent(-1)` → percentage + 2min strategiyaları söndürülür
- Yalnız percentage söndürmək: `SetGCPercent(math.MaxInt64)`
- Minimum hədəf: `GOGC × 4 / 100` MB

**Gözlem aləti:** `GODEBUG=gctrace=1` — hər cycle üçün sətir:
```
gc # @#s #%: ..., #->#-># MB, # MB goal, # MB stacks, # MB globals
```
- `#->#->#` → GC başlanğıc / scan sonu / sweep sonu heap ölçüsü
- `# MB goal` → hədəf heap

**Kitabın nümunəsi:** GOGC=100 ilə GC vaxtı 8-12% (pis!); `GOGC=1000` ilə → 1-2% (cycle-lar arası 10× böyüdü) — amma heap 4× böyüdü.

### 6. GC roots böyükdürsə cycle-lar arası böyüyür (v1.18+)
**Nədir:** Roots = stack-dəki + qlobal dəyişənlərdəki scannable value part-lar; v1.18-dən formulyaya daxil edilib.

**Nəticə:** Böyük stack (150MiB) və ya böyük pointerli qlobal → hədəf heap avtomatik böyüyür → cycle-lar seyrəlir, stagger azalır.

### 7. Memory ballast (v1.18-ə qədər üçün klassik hiylə)
**Nədir:** Proqramın maksimal live heap-indən böyük "köhnə balast" slice — live heap-i süni şəkildə yuxarıda saxlayır.

**Kitabdan kod nümunəsi:**
```go
func main() {
    const ballastSize = 150 << 20 // 150 MiB
    ballast := make([]byte, ballastSize)  // virtual allocation — fiziki RAM almır!
    garbageProducer()
    runtime.KeepAlive(&ballast)  // ballast-ı canlı saxla
}
```
**Sub-kod izahı:**
- Elementlər heç toxunulmadığından Linux-da sadəcə **virtual** ayrılır — fiziki yaddaş istehlakı ~0
- Live heap ~150MB sabitlənir → GOGC=100-də hədəf ~300MB → cycle-lar böyük və stabil
- `runtime.KeepAlive` — optimizatorun ballast-ı erkən "öldürməsin" deyə
- **v1.19-dan:** GOMEMLIMIT gəldi, ballast artıq köhnə üsuldur

### 8. GOMEMLIMIT — memory limit strategy (v1.19+)
**Nədir:** `GOMEMLIMIT` env və ya `debug.SetMemoryLimit(limit)` — runtime yaddaşının ümumi miqdarı limitə yaxınlaşınca GC cycle başlayır.

**Nəyə lazımdır:** Container (K8s/Docker) limit-lərində proqramı OOM-dan qorumaq + ballast-ın müasir alternativi.

## Əsas terminlər
- GC cycle (GC dövrü) — scan+mark + sweep
- GC pacer (GC tempoluğu) — cycle başlanğıc qərarları
- GOGC — hədəf heap faizi (default 100)
- Live heap (canlı yığın) — scan+mark sonundakı qeyri-zibil heap
- GC roots — stack + qlobal scannable hissələr
- Memory ballast (yaddaş balastı) — süni live heap saxlayan böyük blok
- GOMEMLIMIT — v1.19 yaddaş limit strategiyası
- Finalizer — obyekt yığılanda çağırılan hook
- `runtime.SetFinalizer` / `runtime.KeepAlive` — instrumentlər

## Praktik nəticə
1. GC vaxt faizini `GODEBUG=gctrace=1` ilə ölç; >5% olsa GOGC yoxla və ballast/limit düşün.
2. Uzunömürlü kiçik dəyəri böyük blokdan saxlamaq istəyirsənsə — kopyala (substring/subslice tələsi!).
3. Qısaömürlü allocation qaynaqlarını (string concat, string↔[]byte, boxing) hot path-də azalt.
4. Container-da işləyirsənsə `GOMEMLIMIT` təyin et — OOM qorunması + stabil GC.
5. Pointer-lərdən qorxma — kopya xərci vs scan yükünü tarazla; yalnız zərəri aydın olanda qaç.

## Mənbə
Pages: 56-67 (PDF səh. 56-67)
