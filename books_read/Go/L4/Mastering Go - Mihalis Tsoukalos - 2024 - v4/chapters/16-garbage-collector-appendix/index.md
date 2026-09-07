# Appendix — The Go Garbage Collector (Go Zibil Toplayıcısı)

## Bu chapter nədən bəhs edir?

Garbage collection anlayışı, Go GC-nin xarakteristikaları (concurrent/parallel, tri-color
mark-and-sweep, write barrier, trigger mexanizmləri, manual nəzarət), runtime.MemStats
ilə GC statistikası, GODEBUG=gctrace=1, tri-color alqoritminin detallı işləməsi, write
barrier invariantı, runtime.GC, map/slice GC perspektivi və 4 data strukturu variantının
performans müqayisəsi.

## Əsas fikirlər

### 1. Garbage Collection Nədir
**Nədir:** İstifadə olunmayan (referans edilə bilməyən) obyektlərin yaddaşını azad
etmə prosesi. Go-da proqram icrası ilə PARALEL (concurrent) gedir — proqramdan əvvəl
və ya sonra YOX.

**Rəsmi təsvir (kitabdan):** "The GC runs concurrently with mutator threads, is type
accurate (also known as precise), allows multiple GC threads to run in parallel. It is a
concurrent mark and sweep that uses a write barrier. It is non-generational and
non-compacting. Allocation is done using size segregated per P allocation areas to
minimize fragmentation while eliminating locks in the common case."

### 2. Go GC-nin Xarakteristikaları
- **Concurrent və parallel:** tətbiq thread-ləri ilə eyni anda; müəyyən fazalar çoxnüvəli
  paralelləşir — tətbiqi DURDURMUR
- **Tri-color mark-and-sweep:** white/gray/black rəngləri ilə işarələmə (aşağıda ətraflı)
- **Write barrier:** heap-də pointer dəyişəndə GC-nin xəbərdarlığı — izləmə dəqiqliyi
- **Trigger:** allocation həddi/heap ölçüsü aşılınca avtomatik başlayır
- **Manual nəzarət:** runtime.GC() ilə açıq sorğu (bloklayıcı!)

### 3. GC Statistikası — runtime.MemStats
**Kitabdan kod nümunəsi:**
```go
func printStats(mem runtime.MemStats) {
    runtime.ReadMemStats(&mem)
    fmt.Println("mem.Alloc:", mem.Alloc)          // GC azad ETMƏYİŞİ obyekt baytları
    fmt.Println("mem.TotalAlloc:", mem.TotalAlloc) // KÜMÜLATİV — heç azalmır
    fmt.Println("mem.HeapAlloc:", mem.HeapAlloc)   // ≈ Alloc
    fmt.Println("mem.NumGC:", mem.NumGC)            // tamamlanmış dövr sayı
}

// Böyük allocation-larla GC tetikləmə:
for i := 0; i < 10; i++ {
    s := make([]byte, 50000000)
    _ = s
}
```
**Səhv oxunuş:** böyük NumGC = yaddaş istifadəsini YENİDƏN NÖZƏRƏT AL — allocation
optimizasiyası lazımdır.

### 4. GODEBUG=gctrace=1
```bash
GODEBUG=gctrace=1 go run gColl.go
# gc 20 @45.111s 0%: 0.095+0.26+0.009 ms clock, ..., 95->95->0 MB, 95 MB goal, 10 P
```
**Üçlüyün mənası (95->95->0 MB):** GC başlanğıc heap → GC sonu heap → LIVE heap.
P = istifadə olunan logical processor sayı.

### 5. Tri-Color Alqoritmi — Əsas Mexanizm
**Üç dəst:**
| Rəng | Mənası |
|---|---|
| **Black** | Tərkib edilmiş (scan olunub) — white-a pointer-i YOXDUR (invariant!) |
| **Gray** | İşlənmək üzrə — white-a pointer-i OLA BİLƏR |
| **White** | Hələ görünməyib — GC-nin NAMİZƏD-ləri |

**Alqoritm gedişatı:**
1. Başlanğıcda BÜTÜN obyektlər white
2. Root obyektlər (qlobal dəyişənlər, stack) GRAY rənglənir
3. Gray obyekt götürülür → scan → BLACK olur
4. Scan white-a pointer tapırsa → həmin white GRAY olur
5. Gray boşalanadək təkrar
6. Qalan WHITE obyektlər = çatıcılmaz = TOPLANIR

**Vacib invariant:** Black-dən White-a pointer OLA BİLMƏZ — pozulsa GC səhv işləyib
proqramı partladır.

**Write barrier:** mutator (işləyən proqram) hər heap pointer dəyişikliyində kiçik
funksiya icra edir — dəyişdirilən obyekti GRAY rəngləyir. Bu, concurrent icrada
invariantı qoruyur. **Qiyməti:** write barrier icrası latency — concurrency üçün
ödənilən vergi.

**Qeyd:** gray-də unreachable olan obyekt bu dövrdə YOX, NÖVBƏTİ dövrdə toplanır —
suboptimal amma qəbul olunandır.

**Mənşə:** Dijkstra, Lamport, Martin, Scholten, Steffens — "On-the-Fly Garbage
Collection: An Exercise in Cooperation" məqaləsi.

### 6. runtime.GC() — Manual Dövr
- Bloklayıcıdır — məşğul proqramda bütöv proqramı DURA bilər
- Səbəb: dəyişən yaddaşda rəng dəstlərini dəqiqrənləndirmək üçün "safe point" lazımdır
- Kanal da GC hədəfidir: çatıcılmaz kanal bağlanmamış olsa belə resursları azad olunur

### 7. Maps, Slices və GC — 4 Variant Müqayisəsi
**Ortaq pattern:** N=80M element; `runtime.GC()` + `_ = structure[0]` (erkən
toplamaya qarşı; alternativ runtime.KeepAlive()).

**Variantlar (kitabdan):**
```go
// 1) Slice of structs (2 int sahə):
structure = append(structure, data{value, value})

// 2) Map with pointers:
myMap := make(map[int]*int)
myMap[value] = &value

// 3) Map without pointers:
myMap := make(map[int]int)
myMap[value] = value

// 4) Map SHARDING — 2000 xırda map:
split := make([]map[int]int, 2000)
for i := range split { split[i] = make(map[int]int) }
split[i%2000][value] = value
```

**Performans nəticələri:**
| Variant | Vaxt (user) |
|---|---|
| sliceGC | **0.61s** |
| mapNoStar | 10.01s |
| mapSplit | 11.22s |
| mapStar | 23.86s |

**Analiz:**
- Slice həmişə qalib — CONTIGUOUS data + hash funksiyası YOX
- Map hash funksiyası + bucket-da paylanmış (non-contiguous) data = zəif locality
- Pointer-li map ƏN YAVAŞ — hər element ayrıca heap obyekti → GC üçün minlərlə əlavə
  izlənəcək pointer
- Sharding (2000 map) ümumi sürəti yaxşılaşdırmır — GC işini asanlaşdıra bilər, amma
  hash qiyməti qalır

## Əsas terminlər
- Garbage Collector — çatıcılmaz obyektlərin yaddaşını azad edən runtime komponenti
- Concurrent Mark-and-Sweep — proqramla paralel işarələmə-təmizləmə
- Tri-Color Algorithm — white/gray/black dəstləri ilə işarələmə
- Write Barrier — pointer dəyişikliyini GC-yə bildirən mutator kodu
- Mutator — GC ilə paralel işləyən tətbiq proqramı
- Root Objects — birbaşa çatıla bilən obyektlər (qlobal, stack)
- Stop-the-World — proqramı durduran ənənəvi GC; Go onu aradan qaldırır
- Safe Point — GC-nin dəqiqrənləndirə biləcəyi sabit vəziyyət
- runtime.GC() — bloklayıcı manual GC sorğusu
- runtime.ReadMemStats — Alloc/TotalAlloc/NumGC statistikası
- GODEBUG=gctrace=1 — GC dövr izləməsi
- Sharding — böyük map-in xırda map-lərə bölünməsi
- Generational/Non-generational — cavan/qoca nəsil ayrımı (Go NON-generational-dır)

## Praktik nətidə

(1) GC avtomatik və gizlidir — amma NumGC böyüyürsə allocation strukturunu yenidən
düşün. (2) Tri-color invariantı: black-dən white-a pointer YOX — bu, concurrent
GC-nin mümkünlüyünün şərtidir. (3) Write barrier = concurrency vergisi — pointer
dəyişən kod GC yükü artırır. (4) runtime.GC() bloklayır — yalnız test/ölçmə üçün.
(5) Kanallar bağlanmasa da çatıcılmazdırsa GC təmizləyir — amma bağlama daha aydındır.
(6) Performans sırası: slice of structs > map (values) > shard map > map (pointers).
(7) Pointer-li kolleksiyalar GC-yə ən baha başa gəlir — mümkünsə dəyər saxla. (8)
Contiguous memory (slice) həm cache, həm GC üçün qazandır. (9) gctrace=1 və MemStats
— GC davranışının gözə. (10) GC detalları versiyadan-versiyaya dəyişir — rəsmi
sənədə bax (mgc.go mənbəyi: github.com/golang/go).

## Mənbə
Pages: 673-685 (PDF 704-715)
