# Chapter 14 — Efficiency and Performance (Səmərəlilik və Performans)

## Bu chapter nədən bəhs edir?

Benchmark funksiyaları (testing.B, b.N, b.Run sub-benchmark-lar, -benchmem), memory
allocation optimizasiyası (buffer paylaşımı, Reset ilə yenidən istifadə), buffered vs
unbuffered I/O ölçümü, benchstat, səhv benchmark tərifləri, Go yaddaş modeli (heap/stack,
escape analysis, -gcflags -m), slice/map memory leak-ləri və həlləri, memory
pre-allocation və eBPF ilə kernel səviyyəli tracing (gobpf).

## Əsas fikirlər

### 1. Benchmarking — testing.B
**Konvensiyalar:** `Benchmark` prefiksi + böyük hərf/_; `*_test.go` fayllarında;
`b *testing.B` parametri; `go test -bench=.` (regex; parametr YOXSA benchmark İCRA
OLMUR).

**b.N mexanizmi:** hər benchmark ən azı 1 saniyə işləyir; vaxt az olsa b.N artır:
1→2→5→10→20→50... Sürətli funksiya daha ÇOX dəfə icra olunur → dəqiq nəticə.

**Kitabdan kod nümunəsi (make ilə pre-allocation vs append):**
```go
func InitSliceNew(n int) []int {
    s := make([]int, n)          // BİR DƏFƏ allocation
    for i := 0; i < n; i++ { s[i] = i }
    return s
}
func InitSliceAppend(n int) []int {
    s := make([]int, 0)          // hər append-də YENİ allocation riski
    for i := 0; i < n; i++ { s = append(s, i) }
    return s
}

var t []int   // QLOBAL — compiler-ın nəticəni "istifadəsiz" optimizasiyasını əngəlləyir

func BenchmarkNew(b *testing.B) {
    for i := 0; i < b.N; i++ { t = InitSliceNew(i) }
}
```
Nəticə: New 79712 ns/op vs Append 143459 ns/op — append underlying array dolduqca
yeni array + kopyalama.

**Şübhəsiz qayda:** məşğul maşında benchmark ETMƏ; loop xaricində hazırlıq (dəyişən,
bağlantı) mümkündür. `-shuffle=on` — test/benchmark sırasını qarışdır (sıra nəticəyə
təsir edə bilər).

### 2. Memory Allocation Benchmarkı — -benchmem
**Problem versiya:** hər çağırışda YENİ buffer:
```go
func writeMessage(msg []byte) {
    b := new(bytes.Buffer)     // 50 allocs/op!
    b.Write(msg)
}
```
**3 yaxşılaşdırma (kitabdan):**
```go
func writeMessageBuffer(msg []byte, b bytes.Buffer)          // value — KOPYA, faydasız
func writeMessageBufferPointer(msg []byte, b *bytes.Buffer)  // pointer — kopya YOX
func writeMessageBufferWriter(msg []byte, b io.Writer)       // interface — pointer kimi

// Benchmark pattern: buffer BİR DƏFƏ yaradılır, dövrdə PAYLAŞILIR:
func BenchmarkWBufPointerReset(b *testing.B) {
    msg := []byte("Mastering Go!")
    buffer := new(bytes.Buffer)
    for i := 0; i < b.N; i++ {
        for k := 0; k < 50; k++ {
            writeMessageBufferPointer(msg, buffer)
            buffer.Reset()      // ← MƏHƏMMƏM QƏDƏM
        }
    }
}
```
**Nəticələr:**
| Benchmark | ns/op | B/op | allocs/op |
|---|---|---|---|
| WBuf (value) | 1056 | 3200 | 50 |
| WBufPointerNoReset | 337 | 2120 | 0 |
| **WBufPointerReset** | **150.7** | **0** | **0** |
| WBufWriterReset | 151.8 | 0 | 0 |

**Reset() niyə sürətləndirir:**
1. Underlying slice DEALLOC olunmur — len=0 edib memory-ni YENİDƏN istifadə edir
2. Allocation overhead-i (memory manager struktur güncəllənməsi) yox olur
3. Qısaömürlü obyekt sayı azalır → GC təzyiqi düşür

**Compiler optimizasiyası əleyhinə:** nəticəni QLOBAL dəyişənə yaz (ERR = err) — yoxsa
istifadə olunmayan nəticə silinə bilər.

### 3. Buffered vs Unbuffered I/O — b.Run Sub-benchmark-lar
**Kitabdan kod nümunəsi:**
```go
func BenchmarkRead(b *testing.B) {
    buffers := []int{1, 16, 96}                     // 3 bufer ölçüsü
    files := []string{"10.txt", "1000.txt", "5k.txt"} // 3 fayl
    for _, filename := range files {
        for _, bufSize := range buffers {
            name := fmt.Sprintf("%s-%d", filename, bufSize)   // AD MÜTLƏQ!
            b.Run(name, func(b *testing.B) {
                for i := 0; i < b.N; i++ {
                    countChars = CountChars("./testdata/"+filename, bufSize)
                }
            })
        }
    }
}
```
9 sub-benchmark = 3×3 — cədvəl testi benchmark-da da işləyir.

**Tapıntılar:** böyük fayllarda bufer ölçüsü KRİTİK (5k.txt: 1 bayt → 1.3ms; 96 bayt →
33µs — 40x!). Yazmada isə bufer artması az təsir edir. `-10` sonluğu = GOMAXPROCS.

### 4. benchstat — Nəticə Müqayisəsi
```bash
go install golang.org/x/perf/cmd/benchstat@latest
go test -bench=. > r1.txt        # yalnız Benchmark sətirlərini saxla
~/go/bin/benchstat r1.txt r2.txt
# geomean + "vs base" faizi; ~ = əhəmiyyətli dəyişiklik YOX
```
Keçmiş benchcmp-i əvəz edir.

### 5. Səhv Benchmark Tərifləri
```go
// SƏHV 1: b.N artdıqca iş MIQDARI da artır → converg olmur, heç vaxt bitmir:
func BenchmarkFiboI(b *testing.B) {
    for i := 0; i < b.N; i++ { _ = fibo1(i) }
}

// SƏHV 2: hər iterasiya b.N-ə asılı → eyni problem:
func BenchmarkfiboII(b *testing.B) {
    for i := 0; i < b.N; i++ { _ = fibo1(b.N) }
}

// DÜZGÜN: sabit input:
func BenchmarkFiboIV(b *testing.B) {
    for i := 0; i < b.N; i++ { _ = fibo1(10) }
}
// DÜZGÜN: TƏK çağırış b.N-dən asılıdır (loop YOXDUR):
func BenchmarkFiboIII(b *testing.B) {
    _ = fibo1(b.N)
}
```
**Prinsip:** b.N yalnız TƏKRAR SAYI olmalıdır — iş yükü SABİT qalmalıdır.

### 6. Go Yaddaş Modeli — Heap/Stack
**İki allocation növü:**
- **Automatic** — compiler ömrü bilir (lokal dəyişənlər, arqumentlər) → STACK
- **Dynamic** — funksiya scope-undan kənara çıxan data → HEAP

**Escape analysis:** compiler hər dəyişənin yaddan qərarını verir (Go-da new/make
avtomatik heap DEMƏYƏCƏK — C++-dan fərqli!). Goroutine-lərin öz stack-ı var; heap
bölüşdürülür.

**-gcflags '-m' analizi:**
```bash
go run -gcflags '-m' allocate.go
# t escapes to heap           → funksiyadan kənara çıxır
# &Item{} escapes to heap      → heap-ə keçir
# ... argument does not escape → stack-da qalır
go run -gcflags '-m -m'       # daha detallı (çox qalabalıqlı)
```
**Yaddaş modelinin elementləri:** program code (RO, OS mapped), global data (RO),
uninitialized data (anonymous pages, bir dəfə ayrılır, GC BİLMİR), heap (dinamik, GC
sahəsi), stacks (avtomatik).

**Praktik qayda:** stack heap-dən ucuzdur — amma böyük/uzunömürlü obyektlər üçün heap
məcburidir; qərar compiler-indir. Heap ölçüsünü ölçmək proses memory-nin başa düşülməsi
üçün adətən kifayətdir.

### 7. Slice Memory Leak — Leaking Parametr
**Problem kod:**
```go
func getValue(s []int) []int {
    val := s[:3]      // SUBSLICE — underlying 1M-element array-ə BAĞLI qalır!
    return val        // 3 element 1000000 elementin yaddaşını SAXLAYIR
}
// gcflags çıxışı: leaking param: s to result ~r0
```
**Həll — copy:**
```go
func getValue(s []int) []int {
    returnVal := make([]int, 3)
    copy(returnVal, s)     // MÜSTƏQIL nüsxə — əlaqə kəsilir
    return returnVal
}
// gcflags: s does not escape → leak YOX
```

### 8. Map Memory Leak — Bucket-lər Daralmır
**Kitabdan eksperimenti:**
```go
func printAlloc() {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    fmt.Printf("%d KB\n", m.Alloc/1024)
}

n := 2000000
m := make(map[int][128]byte)
// 111 KB → 927931 KB (populate)
// delete hamısı + runtime.GC() → 600767 KB (!!) — bucket SAYI dəyişmir, sadəcə slotlar sıfırlanır
m = nil          // map referansını ÖLDÜR
runtime.GC()     // → 119 KB — GC artıq yaddaşı azad edir
runtime.KeepAlive(m) // map-in ölümdən ƏVVƏL istifadəsini saxlama aləti
```
**Dərs:** map-i boşaltmaq yaddaşı qaytarmır — `m = nil` (və ya yeni map) tələb olunur.

### 9. Memory Pre-allocation
```go
mySlice := make([]int, 0, 100)        // cap=100 → append-lərdə reallocation YOX
myMap := make(map[string]int, 10)      // map üçün də başlanğıc tutum
```
Nə vaxt: ölçü təxmini BİLİNİRSƏ + çoxsaylı insert gözlənilirsə + böyük data. Panaseya
DEYİL — kiçik strukturlarda fərqsiz.

### 10. eBPF — Kernel Səviyyəli Tracing
**Nədir:** Linux kernel daxilində proqramlaşdırıla bilən virtual maşın (Extended
Berkeley Packet Filter). 3 sahə: networking, security, observability. Kerneli dəyişmədən
kernel davranışını genişləndirir; production-safe; səmərəli.

**Necə:** C kodu string-də → BCC ilə compile → kernel hook-a attach → nəticə table-dan oxu.
Go üçün: github.com/iovisor/gobpf/bcc (libbpf əsaslı).

**Kitabdan kod nümunəsi (getuid(2) izləyici):**
```go
const source string = `
#include <uapi/linux/ptrace.h>
BPF_HASH(counts);                    // kernel-daxili hash table
int count(struct pt_regs *ctx) {
    u64 *pointer; u64 times = 0; u64 uid;
    uid = bpf_get_current_uid_gid() & 0xFFFFFFFF;
    pointer = counts.lookup(&uid);
    if (pointer != 0) times = *pointer;
    times++;
    counts.update(&uid, &times);
    return 0;
}
`

m := bpf.NewModule(source, []string{})     // eBPF modulu
defer m.Close()

Uprobe, err := m.LoadUprobe("count")        // count funksiyasını yüklə
err = m.AttachUprobe("c", "getuid", Uprobe, *pid)   // getuid çağırışına bağla
table := bpf.NewTable(m.TableId("counts"), m)      // nəticə cədvəli

// Ctrl+C gözlə, sonra binary table-ni oxu:
sig := make(chan os.Signal, 1)
signal.Notify(sig, os.Interrupt)
<-sig
for it := table.Iter(); it.Next(); {
    k := binary.LittleEndian.Uint64(it.Key())     // BINARY decode MÜTLƏQ
    v := binary.LittleEndian.Uint64(it.Leaf())
    fmt.Printf("%d\t%d\n", k, v)
}
```
**Sub-kod izahı:** C kodu kernel-də icra olunur (u64 sayğaclı hash); Go tərəfi yalnız
module yükləyir, uprobe attach edir, nəticəni decode edir. `import "C"` tələb olunur;
gcc + BPF kitabxanaları lazımdır; YALNIZ Linux.

**Müəllim məsləhəti:** eBPF-ə system admin KİMİ yanaş — mövcud alətlərlə başla
(bpftrace, perf); yalnız həqiqi ehtiyac olanda öz alətini yaz.

## Əsas terminlər
- Benchmark — `Benchmark*`, testing.B, `go test -bench`
- b.N — avtomatik artan təkrar sayğacı (1 san hədəfi)
- ns/op / B/op / allocs/op — vaxt / bayt / allocation metabirimi
- b.Run — sub-benchmark (cədvəl patterni)
- -benchmem — allocation statistikası
- benchstat — iki benchmark nəticəsinin statistik müqayisəsi
- Escape Analysis — stack/heap qərarı (`-gcflags -m`)
- Automatic/Dynamic Allocation — compiler-bilinən / runtime allocation
- Leaking Parameter — funksiya qayıtdıqdan sonra parametri yaşadan referans
- BPF_HASH — eBPF kernel hash table
- Uprobe — istifadəçi-space funksiya hook-u
- runtime.KeepAlive — GC-nin vaxtından əvvəl toplamasının bloku
- runtime.ReadMemStats — yaddaş statistikası

## Praktik nətidə

(1) b.N yalnız SAY olmalıdır — iş yükü sabit; səhv tərif heç vaxt bitməz. (2) Nəticəni
qlobal dəyişənə yaz — compiler optimizasiyası benchmark-u öldürməsin. (3) Allocation
azaltmaq üçün: pointer ilə paylaş + Reset ilə yenidən istifadə — 0 alloc mümkündür.
(4) Value parametr kopyadır — bytes.Buffer kimi strukturları pointer ötür. (5) Subslice
böyük underlying array-i yaşadır — kiçik hissə lazımdırsa COPY et. (6) Map bucket-ləri
daralmır — tam azad üçün nil təyin et. (7) Ölçü məlumdursa pre-allocate et (make cap /
map size). (8) `-gcflags -m` ilə escape axınını gör — heap istifadəsini azaltmaq üçün
ilk addım. (9) Benchmark məşğul maşında = yanlış nəticə. (10) eBPF kernel hadisələrini
səmərəli izləyir — Go+də kodda tracing; amma əvvəl hazır alətləri sına.

## Mənbə
Pages: 627-659 (PDF 658-691)
