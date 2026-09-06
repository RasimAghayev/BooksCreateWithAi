# Chapter 3 — Memory Allocations (Yaddaş ayrılmaları)

## Bu chapter nədən bəhs edir?
Go runtime-da yaddaş bloklarının necə ayrıldığına, stack vs heap fərqinə, hansı əməliyyatların allocation doğurduğuna, size class israflarına və allocation azaltma strategiyalarına (əvvəlcədən yer ayırma, kopyasız filtrləmə, blok birləşdirmə, cache pool).

## Əsas fikirlər

### 1. Memory blocks (yaddaş blokları)
**Nədir:** Yaddaş ayrmanın əsas vahidi — kəsintisiz (continuous) yaddaş seqmenti. Hər value part bir blokda daşıyır; bir blok birdən çox value part daşıya bilər.

**Necə işləyir:** Blok tapmaq CPU resursu tələb edir → nə qədər az blok yaradılsa, o qədər az CPU xərclənir.

**Nəyə lazımdır:** "Hansı kod allocation yaradır?" sualına cavab — çünki allocation = xərc.

### 2. Stack vs Heap
**Nədir:** İki yaddaş zonası — goroutine-in stack-i və proqramın ümumi heap-i.

**Necə işləyir:**
- **Stack:** tapmaq ucuz, yığmaq lazım deyil — goroutine çıxanda bütöv yığılır
- **Heap:** tapmaq bahalı + istifadə olunmadıqda GC cycle lazımır — blok sayı artdıqca GC təzyiqi böyüyür
- Benchmark `allocs/op` metriyası **yalnız heap** allocation-larını sayır

**Nəyə lazımdır:** Escape analysis (Ch 4) hansı dəyərin hara düşəcəyini müəyyən edir — bu fərqi bilmək onun önşərtidir.

### 3. Allocation yaradan əməliyyatlar
Hər biri ən azı bir allocation doğurur:
- dəyişən elan etmək (bəzi hallarda)
- `new` / `make` çağırışı
- slice/map-ə composite literal ilə yazma
- int → string çevirmə
- `+` ilə string birləşdirmə
- string ↔ []byte çevirmə
- string → []rune çevirmə
- dəyəri interfeysə boxing etmə
- capacity çatmayanda `append`
- underlying array dolanda map-ə entry qoymaq

**Kitabın qeydi:** Compiler müəyyən hallarda bu əməliyyatları optimizasiya edib allocation-sız icra edir (Ch 9-da detalları).

### 4. Size class israfları
**Nədir:** Runtime əvvəlcədən müəyyən ölçü sinifləri istifadə edir — 32KB-a qədər diskret class-lar (8, 16, 24, 32, 48, 64, 80, 96...), daha böyüklər 8KB səhifələrdən.

**Necə işləyir:** 33 baytlıq dəyər → 48 baytlıq blok (15 bayt israf); 32769 baytlıq slice → 40960 bayt (5 səhifə, 8191 bayt israf).

**Kitabdan kod nümunəsi:**
```go
var s = []byte{32: 'b'} // len = 33
r = string(s) + string(s)
// 3 allocations: 48 + 48 + 80 = 176 bayt, 44 bayt israf
```
**Sub-kod izahı:**
- `string(s)` ×2 → hər biri 48 baytlıq blok (33→48 class)
- birləşmənin nəticəsi (66 bayt) → 80 baytlıq blok
- `AllocedBytesPerOp()` bunu ölçür: `testing.Benchmark(f).AllocedBytesPerOp()`

**Nəyə lazımdır:** Kiçik dəyərləri toplu saxlamaq (bax: blok birləşdirmə) individual allocation israfını aradan qaldırır.

### 5. Əvvəlcədən kifayət qədər yer ayır (pre-allocate)
**Nədir:** `append` zəncirinə başlamazdan əvvər nəticə slice-ın ölçüsünü hesablayıb `make(0, n)` ilə bir dəfə ayırmaq.

**Kitabdan kod nümunəsi:**
```go
func MergeWithOneLoop(data ...[]int) []int {      // 636 ns, 352 B, 4 allocs
    var r []int
    for _, s := range data { r = append(r, s...) }
    return r
}
func MergeWithTwoLoops(data ...[]int) []int {     // 268 ns, 144 B, 1 alloc — 2.4× sürətli
    n := 0
    for _, s := range data { n += len(s) }
    r := make([]int, 0, n)
    for _, s := range data { r = append(r, s...) }
    return r
}
```
**Sub-kod izahı:**
- OneLoop → capacity çatmayanda 4 dəfə böyümə: 0→2→6→12→24 (hər böyümə = yeni blok + köhnə datanın kopyası)
- TwoLoops → ümumi uzunluğu əvvəlcədən hesabla, bir allocation, kopyasız böyümə
- **Overflow qorunması:** `if k := n + len(s); k < n { panic(...) }` — production kod üçün zəruri

### 6. Mümkünsə allocation-ı tamamilə yox et
**Kitabdan kod nümunəsi (in-place filtrləmə):**
```go
func FilterOneAllocation(data []int) []int {   // 7263 ns, 8192 B, 1 alloc
    var r = make([]int, 0, len(data))
    for _, v := range data {
        if check(v) { r = append(r, v) }
    }
    return r
}
func FilterNoAllocations(data []int) []int {   // 903 ns, 0 B, 0 allocs — 8× sürətli!
    var k = 0
    for i, v := range data {
        if check(v) {
            data[i] = data[k]
            data[k] = v
            k++
        }
    }
    return data[:k]
}
```
**Sub-kod izahı:**
- Yeni slice YOX — mövcud `data` içində saxlanan elementləri özələyib `data[:k]` qaytarır
- `data[i] = data[k]` → `data[k] = v` — swap texnikası ilə dolu hissə sıxılır
- Şərt: input-un dəyişməsinə icazə olmalıdır (yoxsa kopya mütləqdir)

### 7. Blokları birləşdir — bir böyük blok çox kiçik blokdan yaxşıdır
**Nədir:** 100 ayrı `new(Book)` əvəzinə bir `make([]Book, 100)` + pointer slice.

**Kitabdan kod nümunəsi:**
```go
func CreateBooksOnOneLargeBlock(n int) []*Book {  // 4372 ns, 4992 B, 2 allocs
    books := make([]Book, n)
    pbooks := make([]*Book, n)
    for i := range pbooks { pbooks[i] = &books[i] }
    return pbooks
}
func CreateBooksOnManySmallBlocks(n int) []*Book { // 18017 ns, 5696 B, 101 allocs
    books := make([]*Book, n)
    for i := range books { books[i] = new(Book) }
    return books
}
```
**Sub-kod izahı:**
- 2 vs 101 allocation → 4× sürət fərqi; hər 40-baytlıq `Book` ayrı-ayrı ayrılanda 48-bayt class-a düşür (8 bayt israf × 100)
- **Tələ:** kiçik hissələrin ömrü təxminən eyni olmalıdır (birlikdə yaranıb birlikdə ölür → fragmentation olmur)
- **Diqqət:** N=820 hallında (820×40=32800 > 32768 class) səhifə israfı yaranır — böyük blok əksinə 3.5% çox yaddaş istifadə edə bilər

### 8. Value cache pool — allocation-ları təkrar istifadə et
**Nədir:** Tez-tez yaradılıb məhv edilən dəyərlər (oyun NPC-ləri, buffer-lər) üçün pool saxlamaq.

**Kitabdan kod nümunəsi (manual pool):**
```go
var npcPool = struct {
    sync.Mutex
    *list.List
}{ List: list.New() }

func newNPC() *NPC {
    npcPool.Lock()
    defer npcPool.Unlock()
    if npcPool.Len() == 0 { return &NPC{} }
    return npcPool.Remove(npcPool.Front()).(*NPC)
}
```
**Sub-kod izahı:**
- Boş pool → yeni allocation; dolu pool → mövcud NPC-ni təkrar istifadə (list-in önündən götür)
- `sync.Mutex` — çox-goroutine təhlükəsizliyi
- Müasir alternativ: `sync.Pool` (std kitabxana) — GC tərəfindən təmizlənə bilər, amma zero-code tələb edir
- Ssenari: RTS oyununda minlərlə NPC spawn/despawn — GC təzyiqini FPS qorumaq üçün azaltmaq

## Əsas terminlər
- Memory block (yaddaş bloku) — ayrılmanın əsas vahidi
- Size class (ölçü sinifi) — əvvəlcədən təyin edilmiş blok ölçüləri
- Memory page (yaddaş səhifəsi) — 8192 bayt, 32KB-dan böyük blokların vahidi
- Pre-allocation (əvvəlcədən ayırma) — `make(0, n)` ilə köhnəlmiş
- In-place (yerində) — yeni allocation olmadan
- Allocation pressure (ayırmaca təzyiqi) — GC-yə yük
- Value cache pool (dəyər keş hovuzu) — təkrar istifadə üçün
- allocs/op, B/op — benchmark allocation metrikları

## Praktik nəticə
1. `append` zəncirlərindən əvvəl nəticə ölçüsünü hesabla: `make([]T, 0, n)` — 2-4× sürət + 3× az yaddaş.
2. Input-u dəyişmək olursa, filtrləmə/sıxma əməliyyatlarını **in-place** icra et — allocation tam sıfırlanır.
3. Eyni ömürlü kiçik obyektləri bir slice-da topla, pointer-lərlə istinad et — 100 allocation 2-yə düşür.
4. Spawn/destroy dövrlü sistemlərdə pool işlət (`sync.Pool` və ya manual list) — FPS kritik ssenarilərdə.
5. `AllocedBytesPerOp()`/`AllocsPerOp()` ilə real israfı ölç — fərziyyə yox, fakt.
6. 33-48 baytlıq strukturlar 48 bayt class-a düşür — bir strukturu 32 bayta sığdırmaq 33% qənaətdir.

## Mənbə
Pages: 22-32 (PDF səh. 22-32)
