# Chapter 3 — Data types (#17-#29)

## Bu chapter nədən bəhs edir?

Bu chapter əsas tiplər (octal literals, integer overflow, floating points) və slice/map-lərin daxili mexanizmləri (length/capacity, initialization, nil vs empty, copy, append side effects, memory leaks, müqayisə) üzrə 13 səhvi əhatə edir.

---

## Əsas fikirlər

### #17: Creating confusion with octal literals (Octal literallar ilə qarışıqlıq)

```go
sum := 100 + 010
fmt.Println(sum)   // 108 — 110 YOX!
```

**Səbəb:** 0 ilə başlayan integer literal **octal (8-li) say sistemidir** — `010` = 8 (base 10). 

**Octalın faydalı olduğu yer:** Linux fayl icazələri:

```go
file, err := os.OpenFile("foo", os.O_RDONLY, 0644)   // oxunaqlı amma qeyri-aşkar
file, err := os.OpenFile("foo", os.O_RDONLY, 0o644)  // AŞKAR — tövsiyə olunan
```

`0o` prefiksi (və ya `0O` — amma 0/O görünüşcə qarışır) eyni mənanı verir, oxunaqlılığı artırır.

**Digər literal formaları:**
- Binary: `0b`/`0B` (məs., `0b100` = 4)
- Hexadecimal: `0x`/`0X` (məs., `0xF` = 15)
- Imaginary: `i` suffix (məs., `3i`)
- Oxunaqlılıq üçün `_` separator: `1_000_000_000`, `0b00_00_01`

---

### #18: Neglecting integer overflows (Integer overflow-ları nəzərə almamaq)

**10 integer tipi:** int8/16/32/64, uint8/16/32/64 + `int`/`uint` (sistemə bağlı — 32/64 bit).

**Overflow nədir?** Riyazi əməliyyat nəticəsi tipin təmsil edə bildiyi aralıqdan çıxır:

```go
var counter int32 = math.MaxInt32
counter++
fmt.Println(counter)   // -2147483648 — SİLİNT overflow!
```

01111111111111111111111111111111 (MaxInt32) → increment → 10000000000000000000000000000000 (MinInt32). Sign bit 1 olur → mənfi. Two's complement əməliyyatı ilə x + (-x) = 0 təmin olunur.

**Davranış qaydaları:**
- **Compile-time** aşkarlanan overflow → kompilyasiya xətası: `var counter int32 = math.MaxInt32 + 1` → `constant 2147483648 overflows int32`
- **Run-time** overflow → SIRLI, panic YOX → sneaky bug-lar (müsbət ədədlərin cəmi mənfi çıxır)

**Nə vaxt narahat olmalı:** Request counter kimi adi hallarda yox; memory-məhdud kiçik tiplər, böyük ədədlər, konversiyalarda BƏLİ. (Ariane 5 uçuşunun uğursuzluğu (1996) — float64→int16 konversiya overflow-unun nəticəsi.)

**Deteksiya funksiyaları:**

```go
// İncrement:
func IncInt(counter int) int {
    if counter == math.MaxInt { panic("int overflow") }
    return counter + 1
}
// Addition:
func AddInt(a, b int) int {
    if a > math.MaxInt-b { panic("int overflow") }
    return a + b
}
// Multiplication (çoxaddımlı):
func MultiplyInt(a, b int) int {
    if a == 0 || b == 0 { return 0 }
    result := a * b
    if a == 1 || b == 1 { return result }
    if a == math.MinInt || b == math.MinInt { panic("integer overflow") }
    if result/b != a { panic("integer overflow") }   // geri bölmə ilə yoxlama
    return result
}
```

**Qeydlər:** `math.MaxInt`, `math.MinInt`, `math.MaxUint` Go 1.17+-da math paketindədir. Böyük ədədlər üçün `math/big` paketi.

---

### #19: Not understanding floating points (Floating point-ları anlamamaq)

**Əsas həqiqət:** float32/float64 — real arifmetikanın **yaxınlaşmasıdır (approximation)**. Sonsuz real dəyər → məhdud bit (64): dəqiqlik itkisi qaçılmazdır.

```go
var n float32 = 1.0001
fmt.Println(n * n)   // 1.0002 — 1.00020001 gözlənilirdi!
```

**IEEE-754 strukturu:**

| Tip | Sign | Exponent | Mantissa | Dəyər |
|-----|------|----------|----------|-------|
| float32 | 1 bit | 8 bit | 23 bit | sign × 2^exponent × mantissa |
| float64 | 1 bit | 11 bit | 52 bit | eyni düstur |

1.0001 float32-də əslində 1.000100016593933 saxlanılır — dəqiqlik itkisi dəyərin dəqiqliyinə təsir edir.

**3 nəticə qaydası:**

1. **== ilə müqayisə ETMƏ** — delta müqayisəsi et: `|a-b| < epsilon`. (testify `InDelta`.) FPU processor-fərqlidir — bir maşında alınan nəticə digərində fərqlənə bilər → delta testləri portativ edir.
2. **Toplama/çıxma sırası:** Böyüklük tərtibinə görə qruplaşdır — `f2` (10,000-i sonda əlavə edir) `f1`-dən (əvvəl 10,000 + dəyişən) dəqiqliyə görə üstündür. Böyük n → böyük fərq:

| n | Dəqiq nəticə | f1 | f2 |
|---|--------------|----|----|
| 10 | 10010.001 | 10010.000999999993 | 10010.001 |
| 1k | 11000.1 | 11000.099999999293 | 11000.099999999982 |

3. **Vurma/bölmə əvvəl:** `a*(b+c)` vs `a*b+a*c` → birincisi 200030.00200030004 (pis dəqiqlik), ikincisi 200030.0020003. Kombinə olmuş əməliyyatlarda vurma/bölməni ƏVVƏL icra et (accuracy vs icra vaxtı seçimi).

**Xüsusi float dəyərləri:**

```go
var a float64
positiveInf := 1 / a     // +Inf
negativeInf := -1 / a   // -Inf
nan := a / a            // NaN
// Yoxlama: math.IsInf, math.IsNaN
```

IEEE-754-ə görə NaN — yalnız `f != f` olan float (özünə bərabər deyil).

---

### #20: Not understanding slice length and capacity (Slice length/capacity başa düşməmək)

**Slice strukturu:** backing array-ə pointer + **length** (slice-da olan element sayı) + **capacity** (backing array-dəki element sayı).

```go
s := make([]int, 3, 6)   // 6 elementlik array yaradılır, ilk 3-ü 0-a initializə
```

- `s[1] = 1` — element yenilənir, len/cap dəyişmir.
- `s[4] = 0` → panic: `index out of range [4] with length 3` — TUTUMDA olsa belə, LEN-dən kənar yazmaq qadağan.
- `s = append(s, 2)` — boş (gray) elementdən istifadə → len 4 olur.
- `s = append(s, 3, 4, 5)` — array dolanda Go **kapasiteti İKİLƏŞDİRİR**, hər şeyi köçürür, 5-i əlavə edir → yeni array (len 7, cap 12).

**Böyümə qaydası:** 1,024 elementə qədər 2x; sonra 25% artım.

**Slicing (`s1[1:3]`):** Yeni slice EYNİ backing array-ə istinad edir, amma başlanğıc indeksi fərqlidir → fərqli len/cap: `s1` (3,6) → `s2 = s1[1:3]` (2,5).

- `s1[1]` və ya `s2[0]` yenilənsə — HƏR İKİSİNDƏ görünür (eyni yaddaş).
- `s2 = append(s2, 2)` — ortaq array-də element yazılır, amma yalnız `s2`-nin len-i artır: `s1=[0 1 0], s2=[1 0 2]`. `s1` hələ (3,6)-dır.
- `s2`-ni dolanadıqca (5 append-dən sonra) — yeni backing array yaradılır; `s1` köhnə array-da qalır, `s2` yeni massivə keçir.

**Vacib:** append-in shared array üzərindəki təsirlərini anlamadan istifadə → yanlış pressumplar (#25-də davam).

---

### #21: Inefficient slice initialization (Səmərəsiz slice ilkinləşdirməsi)

```go
// PİS — 0 kapasitetlə başlayır: 1→2→4→8... dəfə böyüyür
func convert(foos []Foo) []Bar {
    bars := make([]Bar, 0)
    for _, foo := range foos {
        bars = append(bars, fooToBar(foo))
    }
    return bars
}
// 1000 elementli input: 10 backing array ayrılır, 1000+ element köçürülür, GC yükü
```

**Həll 1 — kapasitet ver:**

```go
bars := make([]Bar, 0, len(foos))   // n elementlik array bir dəfə ayrılır
```

**Həll 2 — uzunluq ver (bir az daha sürətli):**

```go
bars := make([]Bar, len(foos))
for i, foo := range foos {
    bars[i] = fooToBar(foo)          // append YOX — birbaşa indeks
}
```

**Benchmark (1M element):**

| Variant | ns/op | Fərq |
|---------|-------|------|
| EmptySlice | 49,739,882 | ~400% yavaş |
| GivenCapacity | 13,438,544 | baz |
| GivenLength | 12,800,411 | ~4% daha sürətli (append çağırış overhead-i yoxdur) |

**Pebble nümunəsi (Cockroach Labs):** `keys := make([][]byte, 0, len(tombstones)*2)` — indeks hesablaması (`keys[i*2]`, `keys[i*2+1]`) mürəkkəb olduğundan capacity+append oxunaqlığı seçilib. **Dərs:** performance-kritik deyilsə, oxunaqlıq seçilə bilər. Şərtən əlavə olunan elementlərdə (if something(foo)) dəqiq uzunluq məlum deyil — CPU/memory trade-off qərarıdır.

---

### #22: Being confused about nil vs. empty slices

**Təriflər:**
- **Empty slice:** length == 0
- **Nil slice:** slice == nil

4 inicializasiya variantı:

```go
var s []string        // 1: empty=true,  nil=true   — allocation YOX
s = []string(nil)     // 2: empty=true,  nil=true   — allocation YOX
s = []string{}        // 3: empty=true,  nil=false  — allocation VAR
s = make([]string, 0) // 4: empty=true,  nil=false  — allocation VAR
```

**Nil slice həmçinin empty-dir; empty slice mütləq nil deyil.**

**Qaydalar:**
1. `append` nil slice üzərində işləyir: `append(var_s, "foo")` → `[foo]`.
2. **Funksiya slice qaytarırsa — nil slice qaytar** (müdafiə məqsədli non-nil collection şərti deyil — allocation-a ehtiyac yoxdur).
3. Uzunluq məlum olsa: `make([]string, len(ints))` + `s[i]` pattern (#21).
4. `[]string(nil)` — sintaktik şəkər: `append([]int(nil), 42)` — bir sətirdə nil-slice kopyası (#24).
5. `[]string{}` yalnız İlKİ ELEMENTLƏRLƏ: `[]string{"foo", "bar", "baz"}` — elementsiz forma MƏQULDUR (allocation + semantika dəyişir).

**Kitabxana fərqləri:**

```go
var s1 []float32                            // nil
b, _ := json.Marshal(customer{Operations: s1})   // {"ID":"foo","Operations":null}
s2 := make([]float32, 0)                     // empty
b, _ = json.Marshal(customer{Operations: s2})   // {"ID":"bar","Operations":[]}
```

`encoding/json` — nil slice → `null`; empty → `[]`. Strict JSON klientlər üçün kritik fərq. `reflect.DeepEqual(nilSlice, emptySlice)` → false — unit test kontekstində yadda saxla.

---

### #23: Not properly checking if a slice is empty

```go
// PİS — getOperations heç vaxt nil qaytarmır:
if operations != nil { handle(operations) }   // həmişə true

// DÜZGÜN — length yoxla (nil DAXİL hər halı örtür):
if len(operations) != 0 { handle(operations) }
```

- nil slice → `len != 0` false; empty slice → false. Length yoxlaması hər iki ssenarını örtür.
- Çağırılan funksiyanı dəyişə bilməzsən (xarici kitabxana) — length yeganə etibarlı yol.
- **API dizayn prinsipi (Go wiki):** nil və empty slice-ları AYIRMA — çağıran üçün semantik fərq olmamalıdır. Eyni qayda map-lərə şamil olur.

---

### #24: Not making slice copies correctly

```go
src := []int{0, 1, 2}
var dst []int           // len=0!
copy(dst, src)
fmt.Println(dst)        // [] — 3 element YOX
```

**Qayda:** `copy` **min(len(dst), len(src))** element köçürür. Boş dst → 0 element.

```go
// DÜZGÜN:
dst := make([]int, len(src))
copy(dst, src)          // [0 1 2]
```

**Arqument sırası:** `copy(dst, src)` — təyinat birinci, mənbə ikinci (tərsinə çevirmək adi səhvdir).

**Alternativ — append ilə bir sətirdə:**

```go
dst := append([]int(nil), src...)   // 3-length, 3-capacity kopya
```

Bir sətirdədir, amma `copy` daha idiomatik/oxunaqlıdır.

---

### #25: Unexpected side effects using slice append

```go
s1 := []int{1, 2, 3}
s2 := s1[1:2]           // len=1, cap=2 — EYNİ backing array
s3 := append(s2, 10)    // s2 DOLU DEYİL → array-də element yazılır!

fmt.Println(s1)         // [1 2 10] — s1 DƏYİŞDİ!
// s2=[2], s3=[2 10]
```

**Mexanizm:** append slice dolu deyilsə (len < cap) backing array-də birbaşa yazır — ortaq array-i dəyişir.

**Funksiyaya ötürmə halı:**

```go
s := []int{1, 2, 3}
f(s[:2])          // len=2, cap=3 — 3-cü element risk altında!

func f(s []int) { _ = append(s, 10) }  // s[2]=10 → main-dəki s dəyişir: [1 2 10]
```

**Müdafiə həlləri:**

1. **Slice kopyası ötür:**

```go
sCopy := make([]int, 2)
copy(sCopy, s)
f(sCopy)
result := append(sCopy, s[2])   // orijinal birləşdirilir
```

2. **Full slice expression `s[low:high:max]`:**

```go
f(s[:2:2])   // cap = 2-0 = 2 → append içəridə yeni array yaradar, 3-cü element QORUNUR
```

`s[0:2]` (len 2, cap 3) vs `s[0:2:2]` (len 2, cap 2) — ikincisində append kənar təsir edə BİLMƏZ. Kopya overhead-i yoxdur.

---

### #26: Slices and memory leaks

**Hal 1 — Capacity leak (kapasitet sızması):**

```go
// 1M baytlıq mesajlar; son 1000 mesaj tipini (ilk 5 bayt) saxla:
func getMessageType(msg []byte) []byte {
    return msg[:5]   // len=5, cap=1M — BÜTÜN ARRAY YAŞAYIR!
}
// 1000 mesaj × 1MB = ~1 GB yaddaş (5KB yerinə!)
```

**Həll — kopya:**

```go
func getMessageType(msg []byte) []byte {
    msgType := make([]byte, 5)
    copy(msgType, msg)       // həmişə 5/5
    return msgType
}
```

**Qeyd:** Full slice expression (`msg[:5:5]`) kömək ETMİR — Go spesifikasiyası GC-nin inaccessible hissəni geri qazanacağınına zəmanət vermir; testlər göstərir ki, bütün backing array yaddaşda qalır. (`runtime.ReadMemStats(&m)` ilə `m.Alloc` yoxlanılır.)

**Hal 2 — Pointer-li elementlər:**

```go
type Foo struct{ v []byte }

func keepFirstTwoElementsOnly(foos []Foo) []Foo {
    return foos[:2]   // 998 element GC tərəfindən TOPLANMIR!
}
```

**Səbəb:** Element pointer və ya pointer sahəli struct-dırsa (slice özü pointerdir), backing array-də qalan elementlər GC tərəfindən qazanılmır — slice-in işarə etdiyi massiv yaşadığı müddətcə.

**Nüticə:** 1000 Foo × 1MB → slicingdən sonra GC işləsə belə yaddaş DÜŞMÜR (1024072 KB → 1024072 KB).

**Həllər:**
1. **Kopya** (ilk i element): `res := make([]Foo, 2); copy(res, foos)` → 998 element toplanır.
2. **Qalan elementləri nil-lə:**

```go
for i := 2; i < len(foos); i++ {
    foos[i].v = nil
}
return foos[:2]   // 2-len, 1000-cap — amma GC backing array-ləri toplayır
```

**Seçim meyarı:** i (qalan) n-ə yaxındırsa → option 1 (0→i kopyası); i 0-a yaxındırsa → option 2 (i→n nil-ləmə) — benchmark tövsiyə olunur.

---

### #27: Inefficient map initialization

**Map daxili:** hash table = bucket array-i; hər bucket = 8 elementlik key-value massivinə pointer. `hash(key)` → array indeksi; oxunuşda/update-də bucket içində ardıcıll axtarış → worst-case O(p) (p = bucket-lərdəki element sayı).

**Böyümə şərtləri (grow):**
1. **Load factor** (bucket-də orta element sayı) > 6.5 (daxili konstanta)
2. Həddindən artıq overflow bucket-lər

Böyümə: bucket sayı İKİLƏŞİR + bütün açarlar yenidən paylanır → worst-case insert O(n).

**Həll — size hint ver:**

```go
m := make(map[string]int, 1_000_000)   // "ən azı 1M üçün bucket-lər hazır"
```

**Benchmark (1M element insert):**

| Variant | ns/op |
|---------|-------|
| WithoutSize | 227,413,490 |
| WithSize | 91,117,193 — **~60% sürətli** |

Map-də yalnız size (length anlayışı YOX) verilir; limit deyil — hintdir.

---

### #28: Maps and memory leaks

**Map yalnız BÖYÜYÜR, heç vaxt KİÇİLMİR.**

Ssenari: `m := make(map[int][128]byte)` → 1M element əlavə (461 MB) → hamısını sil + GC (runtime.GC) → **293 MB qalır!**

**Səbəb:** `runtime.hmap.B` sahəsi bucket sayının log2-sidir; 1M elementdən sonra B=18 (2^18=262,144 bucket). Elementlər silinsə də B=18 qalır — bucket sayı dəyişmir, sadəcə slotlar sıfırlanır. Overflow bucket-lər də qalır.

**Real ssenari:** Black Friday — bir saatda milyonlarla müştəri → map pikdəki bucket sayına sahib olur; sonrakı günlərdə də yaddaş düşmür.

**Həllər:**
1. **Müntəzəm map kopyalamaq** — hər saat yeni map yaradıb elementləri köçürmək; çatışmazlıq: kopya anında qısamüddətli 2x yaddaş.
2. **Pointer dəyərləri:** `map[int]*[128]byte` — bucket-da 128 bayt yox, 8 bayt pointer saxlanılır:

| Addım | map[int][128]byte | map[int]*[128]byte |
|-------|-------------------|---------------------|
| Boş map | 0 MB | 0 MB |
| 1M element | 461 MB | 182 MB |
| Sil + GC | 293 MB | **38 MB** |

**Qeyd:** Key/value 128 baytı keçərsə Go bucket-da birbaşa saxlamır, pointer istifadə edir.

---

### #29: Comparing values incorrectly (Dəyərlərin yanlış müqayisəsi)

**== / != comparable tiplərdə işləyir:** booleans, numerics, strings, channels (eyni make çağırışı/nil), interfaces (eyni dinamik tip+dəyər/nil), pointers (eyni yaddaş/nil), comparable-lərdən təşkil olmuş structs və arraylər. **Slice və map MÜQAYİSƏ OLUNA BİLMƏZ.**

```go
type customer struct {
    id         string
    operations []float64   // slice sahəsi
}
cust1 == cust2   // KOMPİLİYA OLUNMUR: struct containing []float64 cannot be compared
```

**`any` tələsi:** `any`-yə qoyulmuş müqayisəli tiplər compile olur, amma run-time-da panic:

```go
var cust1 any = customer{...}   // slice sahəli
var cust2 any = customer{...}
cust1 == cust2   // panic: comparing uncomparable type main.customer
```

**Həll 1 — `reflect.DeepEqual`:** Recursive dərin müqayisə (array, struct, slice, map, pointer, interface, funksiya + əsas tiplər):

```go
reflect.DeepEqual(cust1, cust2)   // true
```

**2 diqqət:**
1. nil vs empty collection-u AYIRIR (nil slice ≠ empty slice → false).
2. **Yavaşdır** — reflection run-time introspection edir; lokal benchmarklərdə `==`-dən ~100x yavaş → testlər üçün idealdır, run-time üçün deyil.

**Həll 2 — custom equal metodu:**

```go
func (a customer) equal(b customer) bool {
    if a.id != b.id { return false }
    if len(a.operations) != len(b.operations) { return false }
    for i := 0; i < len(a.operations); i++ {
        if a.operations[i] != b.operations[i] { return false }
    }
    return true
}
```

100 elementli slice benchmark: custom metot `reflect.DeepEqual`-dən ~96x sürətli.

**Əlavə:** Standart kitabxnanın hazır müqayisələrini yoxla — `bytes.Compare` optimallaşdırılmışdır; custom yazmazdan əvvəl təkəri yenidən icad etmə. Testlər üçün `go-cmp`, `testify`.

---

## Əsas terminlər

| Termin | İzah |
|--------|------|
| Octal literal (0 / 0o prefiks) | 8-li say sistemi — `010`=8; Linux icazələri üçün `0o644` |
| Integer overflow | Nəticə tipin aralığından çıxır — run-time-da SİLENT |
| Two's complement | Mənfi ədədlərin təmsili — x + (-x) = 0 |
| `math.MaxInt`/`MinInt`/`MaxUint` | Overflow yoxlaması üçün sabitlər (Go 1.17+) |
| Mantissa / Exponent | Floating point komponentləri: baza dəyər / vurucu |
| IEEE-754 | Float standardı — float32 (8+23 bit), float64 (11+52 bit) |
| NaN (+Inf, -Inf) | Təyin olunmayan nəticə; `f != f` — `math.IsNaN/IsInf` |
| Backing array | Slice-in altında duran məxsusi array |
| Slice length / capacity | Slice-dəki element sayı / backing array element sayı |
| Slice growth | 1024-ə qədər 2x, sonra +25% |
| Slicing | `s[low:high]` — ortaq backing arrayli yeni slice |
| Full slice expression | `s[low:high:max]` — capacity = max - low; append qoruması |
| Nil vs empty slice | nil == nil (alloc yox); empty len=0 (alloc VAR) |
| `copy(dst, src)` | min(len(dst), len(src)) element köçürür |
| Capacity leak | Kiçik slice böyük array-i yaşadır — kopya həlli |
| Bucket (map) | 8 elementlik key-value massivi; hmap.B = log2(bucket sayı) |
| Load factor | Bucket başına orta element — 6.5-i keçəndə map 2x böyüyür |
| Comparable type | == işləyən tiplər — slice/map YOX |
| `reflect.DeepEqual` | Dərin refleksiyalı müqayisə — nil/empty ayırır, ~100x yavaş |

---

## Praktik nəticə

1. **0 ilə başlayan ədədlər octal-dır** — `0o` prefiksi ilə aşkar yaz; `_` separatoru ilə oxunaqlıq (`1_000_000_000`).
2. **Run-time overflow sirlidir** — kritik yerlərdə MaxInt/MinInt yoxlamaları yaz; böyük ədədlər üçün `math/big`.
3. **Float yaxınlaşmadır:** == YOX (delta), toplama üçün böyüklük-qruplaşdırma, vurma/bölməni əvvəl icra et.
4. **Slice length ≠ capacity:** append len < cap olduqda ortaq array-də yazır — səhv pressumpların əsas mənbəyi.
5. **İlkinləşdirmədə ölçü ver:** `make([]T, 0, n)` və ya `make([]T, n)` — 400% qədər sürət fərq; map üçün `make(map[K]V, n)` — 60%.
6. **Boş slice qaytararkən nil seç:** `var s []string` — allocationsız; `[]string{}` elementsiz MƏQUL.
7. **Boşluq yoxlaması = `len(s) != 0`** — nil/empty hər ikisini örtür; API-də nil/empty ayırmayın.
8. **copy min(len)-dir:** dst-i əvvəlcədən `make(len(src))` ilə yarat; sıra: `copy(dst, src)`.
9. **Ortaq array qoruması:** `s[:2:2]` full slice expression (kopyasız) və ya kopya ötür.
10. **Böyük slice-dan kiçik hissə saxlayırsansa KOPYALA** — `msg[:5]` 1MB array-i yaşadır; pointerli elementlərdə qalanları nil-lə.
11. **Map yalnız böyüyür:** pik yaddaş qalır — periodik re-yaratma və ya `map[K]*V` (38MB vs 293MB).
12. **== yalnız comparable-larda:** slice/map üçün `reflect.DeepEqual` (test) və ya custom `equal()` (run-time, ~96x sürətli).

---

## Mənbə

- Kitab: *100 Go Mistakes and How to Avoid Them* — Teiva Harsanyi, Manning, 2022 (ISBN 9781617299599)
- Chapter 3: "Data types", book səh. 56–94
- PDF səhifələri: 76–114
- İstinadlar: Pebble (github.com/cockroachdb/pebble); testify (github.com/stretchr/testify); go-cmp (github.com/google/go-cmp); Ariane 5 (bugsnag.com/blog/bug-day-ariane-5-disaster)
