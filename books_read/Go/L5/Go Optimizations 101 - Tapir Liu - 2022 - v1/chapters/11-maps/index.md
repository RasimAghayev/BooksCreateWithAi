# Chapter 11 — Maps (Xəritələr)

## Bu chapter nədən bəhs edir?
Map-in daxili strukturu (hashtable + backing array), entry təmizləmə, `m[k]++` optimizasiyası, GC-scan azaldan açar/element tip seçimi, pointer-element və index-table patternlərinə və kiçik dəyərli açarlar üçün map əvəzinə array cədvəli istifadəsinə.

## Əsas fikirlər

### 1. Map daxili: hashtable, cap yoxdur, backing array azalmır
**Nədir:** Map daxildə hashtable-dır — entry-lər backing array-da; dolanda yeni böyük array ayrılır, entry-lər köçürülür.
**Vacib:** Backing array heç vaxt KİÇİLMİR — bütün entry-lər silinsə belə! (yaddaş israfı forması)

### 2. Entry təmizləmə və backing array azad etmə
**Kitabdan kod nümunəsi:**
```go
// təmizləmə (xüsusi optimizasiya — çox sürətli):
for key := range aMap { delete(aMap, key) }

// backing array-i azad etmək:
aMap = nil              // təkrar istifadə önəmli deyilsə (yeni entry az gələcəksə)
aMap = make(map[K]V)    // yeni boş map
```
**Sub-kod izahı:** Loop sürətlidir amma array yaddaşda qalır; `nil`/`make` köhnə array-i GC-ə buraxır. Təmizlənəndən sonra ÇOX entry gələcəksə — loop (böyümədən qaçınmaq üçün) saxla; gəlməyəcəksə — release et.

### 3. `m[k]++` — bir hashləmə; `m[k] = m[k] + 1` — İKİ hashləmə
```go
m[99]++              // 11.31 ns
m[99] += 1           // 11.21 ns
m[99] = m[99] + 1    // 16.10 ns — 43% yavaş — açar 2 dəfə hashlanır!
```

### 4. Açar və element tiplərində pointer olmasın — GC scan olunmaz
**Nədir:** Açar+element tiplərində pointer YOXDURSA, GC scan fazasında map-in entry-ləri ÜMUMİYYƏTLƏ yoxlanılmır.

**Tətbiq:** Slice/array/channel-lərə də aiddir. Ən böyük qazanc — milyonlu map-lərdə.

### 5. Qısa string açarlar → `[N]byte` açar
**Nədir:** String daxildə pointerdir → string açarlı map hər GC-də scan olunur. Maksimal uzunluq məlum və kiçikdirsə, `[N]byte` açar pointer-sizdir.
```go
var mapA = make(map[string]int, 1 << 16)     // scan olunur — string açar pointerlidir
var mapB = make(map[[32]byte]int, 1 << 16)   // scan YOX — entries pointersiz!
```
**Xərc:** `[32]byte` açar qoşulanda 32 bayt kopyalanır — 65K entry-də GC qənaəti adətən üstündür.

### 6. Sıx modifikasiya → pointer element (və ya index table)
**Problem:** `m[string(w)]++` L-value modifikasiyadır — hər dəfə `string(w)` çevirməsi allocate edir (Ch 9.1.3-dən fərqli olaraq!).

**Kitabdan kod nümunəsi (3 variant):**
```go
// A: birbaşa — 11600 ns, 62 allocs/op
func IncA(w []byte) { wordCounterA[string(w)]++ }

// B: pointer element — 1543 ns, 0 allocs
func IncB(w []byte) {
    p := wordCounterB[string(w)]      // oxuma — allocation-sız!
    if p == nil {
        p = new(int)                    // yalnız ilk dəfə
        wordCounterB[string(w)] = p     // yalnız ilk dəfə allocate
    }
    *p++
}

// C: index table — 1609 ns, 0 allocs (uzunmüddət ƏN YAXŞI — pointer az!)
var wordIndexes = make(map[string]int)   // söz → indeks
var wordCounters []int                   // sayğac cədvəli
func IncC(w []byte) {
    if i, ok := wordIndexes[string(w)]; ok {
        wordCounters[i]++
    } else {
        wordIndexes[string(w)] = len(wordCounters)
        wordCounters = append(wordCounters, 1)
    }
}
```
**Sub-kod izahı:**
- B: pointer element yarandıqdan sonra dəyişmir → yalnız oxuma qalır (allocation-sız)
- C: map-in özü heç modifikasiya olunmur; sayğaclar slice-dadır — B-dən az pointer, GC scan daha yüngül
- Benchmark "0 allocs" → ortalama <1 allocation-un truncation-u (dizayn qərarı)

### 7. Map-i bir addımda böyüt
Maksimum entry sayı məlumdursa: `make(map[K]V, n)` — çoxaddımlı böyümə və rehash-lardan qaç.

### 8. Bool/kiçik-dəyər açarlı map əvəzinə index table
**Problem:** `map[bool]func(){true: f, false: g}` oxunaqlıdır, amma 11× yavaşdır:
```go
Benchmark_IfElse:     4.16 ns    // if-else
Benchmark_MapSwitch: 47.46 ns   // map[bool] — hashing + lookup!
```
**Həlli — array + bool→int:**
```go
func b2i(b bool) (r int) { if b { r = 1 }; return }
var a = [2]func(){g, f}         // index table
func IndexTable(x bool) func() { return a[b2i(x)] }  // if-else sürətində!
```
**Ümumi qayda:** Açar dəyərlərinin çoxluğu kiçikdirsə (bool, az saylı enum) — array cədvəli map-i əvəz edir.

## Əsas terminlər
- Backing array (dəstəkləyici massiv) — entry saxlayan hashtable massivi
- Rehash (yenidən hashləmə) — böyümə zamanı entry köçürməsi
- Pointer-free container (pointersiz konteyner) — GC-scan-dan azad
- Index table (indeks cədvəli) — map əvəzinə array sorğusu
- L-value map index — `m[k] = ...` (allocate edir) vs R-value (`v := m[k]` — etmir)
- Word counter pattern — sayğac üçün pointer-element/index-table həlli

## Praktik nəticə
1. `m[k]++` / `m[k] += v` işlət — `m[k] = m[k] + v` 43% yavaşdır.
2. Böyük map-lərdə açar/elementdə pointer olmasın: `[N]byte` açar, int element — GC scan sıradan çıxır.
3. Söz-sayğac patterni: `map[string]*int` (pointer element) və ya `map[string]int` + `[]int` cədvəli (ən yaxşı) — 7× sürət, 0 allocation.
4. Bool/enum "map-switch" — gözəldir amma 10× yavaşdır; `[2]T` index table eyni gözəlliyi if-else sürətində verir.
5. Map ölçüsünü əvvəlcədən bilirsənsə `make(map, n)` ilə bir addımda ayır.
6. Təmizlənən map-in backing array-i yaddaşda qalır — tam azad etmək üçün `nil`/`make`.

## Mənbə
Pages: 118-124 (PDF səh. 118-124)
