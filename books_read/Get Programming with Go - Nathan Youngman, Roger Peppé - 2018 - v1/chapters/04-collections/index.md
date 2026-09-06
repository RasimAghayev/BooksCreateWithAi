# Unit 4 — Collections (Lesson 16-20)

## Bu unit nədən bəhs edir?

Komposit tiplər: array-lər (sabit uzunluq, kopyalanır, bounds), slice-lar (array-ə pəncərə, append, len/cap, 3-indeks slicing, make, variadic), map-lər (key-value, comma-ok, delete, make, set kimi istifadə, qruplama). Capstone: Conway's Game of Life.

**PDF səhifələr:** 136-173 (L16: 136-144, L17: 145-152, L18: 153-160, L19: 161-169, L20: 170-173)

## Əsas fikirlər

### 1. Array-lər (L16)
**Nədir:** Sabit uzunluqlu, sıralı kolleksiya. Uzunluq tipin HİSSƏSİDİR.

```go
var planets [8]string          // 8 element, hamısı "" (zero value)
planets[0] = "Mercury"
earth := planets[2]             // indeks 0-dan

dwarfs := [5]string{"Ceres", "Pluto", "Haumea", "Makemake", "Eris"}   // composite literal
planets := [...]string{           // ... → kompilyator sayır
    "Mercury", "Venus", "Earth", "Mars",
    "Jupiter", "Saturn", "Uranus", "Neptune",   // sondaki vergül MÜTLƏQ
}
fmt.Println(len(planets))        // 8
```

**Bounds:**
```go
planets[8] = "Pluto"    // compile-time xəta: invalid array index 8
i := 8
planets[i] = "Pluto"     // runtime PANIC: index out of range
```
Panic — C-dəki kimi yaddaş korlanmasından (unspecified behavior) YAXŞIDIR.

**İterasiya:**
```go
for i := 0; i < len(dwarfs); i++ { ... }      // klassik
for i, dwarf := range dwarfs { ... }           // range — daha az səhv
for _, dwarf := range dwarfs { ... }            // indeks lazım deyilsə
```

### 2. Array-lər KOPYALANIR (L16)
```go
planetsMarkII := planets       // TAM KOPYA
planets[2] = "whoops"
fmt.Println(planetsMarkII)      // Earth hələ də yaşayır!
```

**Funksiya parametri kimi — terraformƏSƏRSİZDİR:**
```go
func terraform(planets [8]string) {     // kopya üzərində işləyir
    for i := range planets {
        planets[i] = "New " + planets[i]
    }
}
terraform(planets)                       // original DƏYİŞMİR
terraform(dwarfs)                        // XƏTA: [5]string ≠ [8]string!
```
→ **Array-lər funksiya parametri kimi nadir istifadə olunur** — slice işlə.

**Çoxölçülü:**
```go
var board [8][8]string     // 8 array-dan array
board[0][0] = "r"
```

### 3. Slice — array-ə pəncərə (L17)
**Nədir:** Altta yatan array-ə görünüş; kopya deyil — paylaşılan görünüş.

```go
planets := [...]string{"Mercury", "Venus", "Earth", "Mars",
    "Jupiter", "Saturn", "Uranus", "Neptune"}

terrestrial := planets[0:4]     // [Mercury Venus Earth Mars] — half-open: 0 daxil, 4 XAİR
gasGiants := planets[4:6]
iceGiants := planets[6:8]

terrestrial := planets[:4]       // default başlanğıc = 0
iceGiants := planets[6:]          // default son = len
allPlanets := planets[:]           // hamısı
```

**Paylaşma — dəyişiklik hamıya görünür:**
```go
iceGiantsMarkII := iceGiants          // slice kopyası — ama EYNİ array-ə işarə!
iceGiants[1] = "Poseidon"
fmt.Println(planets)                   // ...Uranus Poseidon — ARRAY dəyişdi
fmt.Println(iceGiants, iceGiantsMarkII, ice)   // hamısında Poseidon
```

**Slice-in slice-i:**
```go
giants := planets[4:8]
gas := giants[0:2]
ice := giants[2:4]
```

**String slicing:** `tune := neptune[3:]` → nəticə STRING; amma indekslər BAYT hesabı ilədir (runa YOX!) — "¿Cómo estás?"[:6] = "¿Cóm" (2-baytlı ¿ + C ó m).

### 4. Slice elanı (composite literal) (L17)
```go
dwarfs := []string{"Ceres", "Pluto", "Haumea", "Makemake", "Eris"}
// arxada: 5-element array + tam görünüş
// %T: []string (slice) vs [5]string (array)
```

### 5. Slice funksiyalarda — güc (L17)
**Kitabdan kod nümunəsi:**
```go
// hyperspace removes the space surrounding worlds
func hyperspace(worlds []string) {
    for i := range worlds {
        worlds[i] = strings.TrimSpace(worlds[i])   // ALT ARRAY-İ DƏYİŞDİRİR
    }
}

planets := []string{" Venus   ", "Earth  ", " Mars"}
hyperspace(planets)
fmt.Println(strings.Join(planets, ""))   // VenusEarthMars — boşluqlar getdi!
```
- Funksiya slice-in **başlangıc/son/istiqamətini dəyişə bilməz**, amma **elementlərini dəyişə bilər** — alt array paylaşılır
- Uzunluq tipin hissəsi DEYİL → istənilən ölçülü slice qəbul olunur

### 6. Slice + metod (L17)
```go
type StringSlice []string
func (p StringSlice) Sort()

// istifadə:
sort.StringSlice(planets).Sort()     // konvertasiya + metod
sort.Strings(planets)                 // qısayol helper
```
**Dərs:** slice üzərində tip yarat → metod bağla → siniflərdən daha çevik.

### 7. append + len/cap (L18)
```go
dwarfs := []string{"Ceres", "Pluto", "Haumea", "Makemake", "Eris"}   // len=cap=5
dwarfs = append(dwarfs, "Orcus")                        // variadic: 1 element
dwarfs = append(dwarfs, "Salacia", "Quaoar", "Sedna")   // bir neçəsi
```

**Capacity araşdırması:**
```go
func dump(label string, slice []string) {
    fmt.Printf("%v: length %v, capacity %v %v\n", label, len(slice), cap(slice), slice)
}

dwarfs1 := []string{"Ceres", "Pluto", "Haumea", "Makemake", "Eris"}  // len 5, cap 5
dwarfs2 := append(dwarfs1, "Orcus")          // len 6, cap 10 — YENİ ARRAY (2x)
dwarfs3 := append(dwarfs2, "Salacia", "Quaoar", "Sedna")   // len 9, cap 10
// dwarfs3[1] = "Pluto!" → dwarfs2 də dəyişir, dwarfs1 YOX (başqa array)
```
- **Capacity dolanda:** append → yeni array ayır (adətən 2x) + köhnəni kopyala + slice-i oraya yönəlt
- dwarfs2/dwarfs3 eyni YENİ array-i görür; dwarfs1 köhnədə qalır

### 8. Üç indeksli slicing (L18)
```go
planets := []string{"Mercury", ..., "Neptune"}

terrestrial := planets[0:4:4]      // len 4, cap 4 — CAP MƏHDUD!
worlds := append(terrestrial, "Ceres")   // cap doludur → YENİ array
fmt.Println(planets)               // Jupiter qorunub saxlanılır ✓

terrestrial = planets[0:4]          // len 4, cap 8 (default!)
worlds = append(terrestrial, "Ceres")   // cap var → alt array-də YAZIR
fmt.Println(planets)               // ...Mars Ceres Saturn... — Jupiter ÖLDÜ ✗
```
**Qayda:** alt array-i korlamaq istəmirsənsə **həmişə 3-indeks slicing** (`s[i:j:j]`) default seçim et.

### 9. make ilə preallokasiya (L18)
```go
dwarfs := make([]string, 0, 10)    // len 0, cap 10 — 10 append-dək yeni array YOX
dwarfs := make([]string, 10)       // len=cap=10, hamısı "" (zero value)
```
Əvvəldən məlum olan ölçü üçün allokasiya/kopya xərclərindən xilas olur.

### 10. Variadic funksiyalar (L18)
```go
func terraform(prefix string, worlds ...string) []string {
    newWorlds := make([]string, len(worlds))    // kopya yarat — arqumenti dəyişməmək üçün
    for i := range worlds {
        newWorlds[i] = prefix + " " + worlds[i]
    }
    return newWorlds
}

twoWorlds := terraform("New", "Venus", "Mars")     // ayrı-ayrı arqumentlər
planets := []string{"Venus", "Mars", "Jupiter"}
newPlanets := terraform("New", planets...)          // slice-i AÇ (...)
```
**Ellipsis-in 3 istifadəsi:** (1) `[...]` — kompilyator saysın; (2) `worlds ...string` — variadic parametr; (3) `planets...` — slice-i arqumentlərə aç.

### 11. Map əsasları (L19)
**Digər dillərdə:** Python dictionary, Ruby hash, JS object, PHP associative array, Lua table.

```go
temperature := map[string]int{
    "Earth": 15,
    "Mars":  -65,
}
temp := temperature["Earth"]     // 15 — oxu
temperature["Earth"] = 16        // yaz/yenilə
temperature["Venus"] = 464       // əlavə
```

**Comma, ok idiomu (mövcudluq yoxlaması):**
```go
moon := temperature["Moon"]     // 0 — zero value (xəta YOX!)
if moon, ok := temperature["Moon"]; ok {
    fmt.Printf("On average the moon is %vº C.\n", moon)
} else {
    fmt.Println("Where is the moon?")     // ← bu icra olunur
}
// "Moon" varsa və dəyəri 0-dırsa: ok=true — fərqi göstərir
```

### 12. Map-lər KOPYALANMIR (L19)
```go
planets := map[string]string{"Earth": "Sector ZZ9", "Mars": "Sector ZZ9"}
planetsMarkII := planets               // EYNİ data!
planets["Earth"] = "whoops"             // hər ikisində görünür
delete(planets, "Earth")               // builtin — hər ikisindən silinir
fmt.Println(planetsMarkII)             // map[Mars:Sector ZZ9]
```
Slice-ların alt array paylaşması kimi — map funksiyaya ötürüləndə dəyişə bilər!

### 13. make + istifadələr (L19)
```go
temperature := make(map[float64]int, 8)   // preallokasiya; len həmişə 0
```

**Sayğac (frequency):**
```go
temperatures := []float64{-28.0, 32.0, -31.0, -29.0, -23.0, -29.0, -28.0, -33.0}
frequency := make(map[float64]int)
for _, t := range temperatures {
    frequency[t]++                  // hər tapıntıda artır
}
for t, num := range frequency {
    fmt.Printf("%+.2f occurs %d times\n", t, num)
}
```
**Map sırası QARANTIYA OLUNMUR** — hər run-da fərqli ola bilər!

**Qruplama (map of slices):**
```go
groups := make(map[float64][]float64)
for _, t := range temperatures {
    g := math.Trunc(t/10) * 10         // -28 → -20, -31 → -30
    groups[g] = append(groups[g], t)
}
// -20: [-28 -29 -23 -29 -28], -30: [-31 -33], 30: [32]
```

**Set kimi map:**
```go
set := make(map[float64]bool)
for _, t := range temperatures {
    set[t] = true                     // dublikatlar avtomatik yox olur
}
if set[-28.0] { /* set member */ }

// Sıralı unique slice:
unique := make([]float64, 0, len(set))
for t := range set {
    unique = append(unique, t)
}
sort.Float64s(unique)
```

### 14. Capstone: Conway's Game of Life (L20)
**Qaydalar:** canlı + <2 qonşu → öl; canlı + 2-3 → yaşa; canlı + >3 → öl; ölü + tam 3 → doğul.

**Struktur:**
```go
const (
    width  = 80
    height = 15
)

type Universe [][]bool          // slice — paylaşıla bilər

func NewUniverse() Universe     // make ilə yarat (zero value = false = ölü)
func (u Universe) Show()        // * = canlı, ' ' = ölü
func (u Universe) Seed()        // ~25% random canlı
func (u Universe) Alive(x, y int) bool   // wrap-around: % modulus
func (u Universe) Neighbors(x, y int) int // Alive istifadə et — 0-8
func (u Universe) Next(x, y int) bool      // qaydaları tətbiq et
func Step(a, b Universe)        // A-dan oxu → B-yə yaz
```

**Kritik məqamlar:**
- **Wrap-around:** `(y+height)%height` — grid kənarları birləşir (modulus!)
- **Parallel universe problemi:** növbəti nəsli eyni grid-də yazsan qonşu sayları korlanar → 2 universe + `a, b = b, a` swap
- `"\x0c"` — ekranı təmizlə (ANSI); `time.Sleep` — animasiya

## Array vs Slice vs Map xülasəsi
| | Array | Slice | Map |
|---|---|---|---|
| Uzunluq | Sabit (tipin hissəsi) | Dəyişən (len) | Dəyişən |
| Açar | 0..n-1 | 0..n-1 | istənilən tip |
| Kopyalanır? | ✅ (tam) | ❌ (alt array paylaşır) | ❌ (data paylaşır) |
| Yaradılış | `[N]T{...}` | `[]T{...}` / make / slicing | `map[K]V{...}` / make |
| Böyümə | YOX | append (cap idarə edir) | avtomatik |
| Funksiyaya | nadir (uzunluq fərqləri!) | standart seçim | standart |

## Əsas terminlər
- Array (sabit uzunluq; uzunluq = tip)
- Composite Literal (`{...}`)
- Bounds Check / Panic
- Slice (pəncərə/görünüş)
- Half-open Range (`[0:4)` — 4 daxil DEYİL)
- Default Indices (`s[:4]`, `s[4:]`, `s[:]`)
- Underlying Array
- Shared View (paylaşılan görünüş)
- append / len / cap
- Three-index Slicing (`s[i:j:k]`)
- make (preallokasiya)
- Variadic Function / Ellipsis (3 istifadə)
- Slice Expansion (`slice...`)
- Map (key-value; dict/hash/object)
- Comma-ok Idiom
- delete (builtin)
- Set (improvize — map[T]bool)
- Wrap-around (modulus grid)
- Parallel Buffer (double buffering)

## Praktik nəticə
- Funksiyalara array YOX, **slice** ötür — uzunluq tip fərqi səhvə yol verir; slice uzunluq-agnostikdir.
- Slice-almaqda 3-indeks (`s[i:j:j]`) default et — append alt array-i səthən korlamaya bilər (Jupiter dərsi).
- Böyük miqdarda append olunacaq slice-ı `make([]T, 0, n)` ilə başlat — allokasiya/kopya qənaəti.
- Map mövcudluğu üçün **comma-ok** — `v` zero-value ola bilər, `ok` fərqi göstərir.
- Map sırasına etibar etmə — sıralı lazımdırsa açarları slice-a çıxarıb sort et.
- Set lazımdırsa `map[T]bool` — dublikatsız kolleksiya.
- State-dəyişən simulyasiyalarda (Game of Life) double buffering — oxu A-dan, yaz B-yə, swap et.

## Mənbə
Pages: 136-173 (PDF), book pages 121-158 (Lesson 16-20)
