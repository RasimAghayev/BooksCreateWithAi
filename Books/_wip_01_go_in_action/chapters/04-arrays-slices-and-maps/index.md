# Chapter 4 — Arrays, slices, and maps

## Bu chapter nədən bəhs edir?

Go-nun üç daxili kolleksiya strukturunun daxili quruluşu və istifadəsi: **arrays** (sabit uzunluqlu massivlər), **slices** (dinamik hissələr) və **maps** (açar/dəyər cütləri). Hər birinin internal representation, elan/initializasiya, böyümə, iterasiya və funksiyalararası ötürülmə davranışı izah olunur.

## Əsas fikirlər

### 1. Arrays (massivlər)
**Nədir:** Sabit uzunluqlu, eyni tip elementlərdən ibarət ardıcıl (contiguous) yaddaş bloğu.

**Necə işləyir:** Yaddaş ardıcıl ayrılır → CPU cache-də daha uzun qalır, index arifmetikası ilə sürətli iterasiya. Elan edildikdən sonra nə tip, nə uzunluq dəyişə bilər.

**Kitabdan kod nümunələri (elan formaları):**
```go
// 5 elementli int array — bütün elementlər zero value (0)
var array [5]int

// Array literal ilə
array := [5]int{10, 20, 30, 40, 50}

// Uzunluğu Go hesablasın
array := [...]int{10, 20, 30, 40, 50}

// Yalnız index 1 və 2-ni initialize et, qalanları zero
array := [5]int{1: 10, 2: 20}

// Pointer array
array := [5]*int{0: new(int), 1: new(int)}
*array[0] = 10
```

**Sub-kod izahı:**
- `[5]int` → tipə uzunluq da daxildir! `[4]string` və `[5]string` **fərqli tiplərdir**
- `array[2] = 35` → `[]` operatoru ilə elementə çıxış
- `*array[0] = 10` → pointer array-də dəyəri `*` ilə ver
- `new(int)` → zero-initialized int üçün pointer qaytarır

**Array value-dır (kopyalanır):**
```go
array1 := [5]string{"Red", "Blue", "Green", "Yellow", "Pink"}
var array2 [5]string
array2 = array1          // BÜTÜN dəyərlər kopyalanır
```
- `[4]string`-ə `[5]string` mənimsətmə → compile xətası: `cannot use array2 (type [5]string) as type [4]string`
- Pointer array kopyalananda **pointer dəyərləri** kopyalanır — iki array eyni string-lərə işarə edir

**Böyük array-in funksiyaya ötürülməsi:**
```go
// 8 MB kopyalanır — bahalı:
var array [1e6]int
foo(array)
func foo(array [1e6]int) { ... }

// Yalnız 8 byte (pointer) kopyalanır — səmərəli:
foo(&array)
func foo(array *[1e6]int) { ... }
```

**Multidimensional:**
```go
array := [4][2]int{{10, 11}, {20, 21}, {30, 31}, {40, 41}}
array := [4][2]int{1: {20, 21}, 3: {40, 41}}
array[0][0] = 10
var array3 [2]int = array1[1]   // bir dimension-u kopyala
```

### 2. Slices (hissələr)
**Nədir:** Dinamik böyüyüb kiçilə bilən kolleksiya — altta yatan array-i (underlying array) abstrakt edən 3 sahəli kiçik obyekt.

**Necə işləyir (internals):** Slice 3 sahədən ibarətdir:
1. Underlying array-ə pointer
2. **Length** (uzunluq) — slice-ın çıxışı olan element sayı
3. **Capacity** (tutum) — underlying array-də böyümə üçün mövcud element sayı

64-bit arxitekturada slice cəmi **24 byte**-dır (8+8+8).

**Yaratma formaları:**
```go
// make — length 5, capacity 5
slice := make([]string, 5)

// make — length 3, capacity 5
slice := make([]int, 3, 5)

// make([]int, 5, 3) → COMPILE XƏTASI: len larger than cap in make([]int)

// Slice literal — uzunluq və tutum = element sayı
slice := []int{10, 20, 30}

// 100 elementli slice, yalnız 100-cü initialize
slice := []string{99: ""}

// NIL slice — initialization yoxdur (ən çox istifadə olunan)
var slice []int

// EMPTY slice
slice := make([]int, 0)
slice := []int{}
```

**Vacib fərq:** `[3]int{...}` = **array** (uzunluq yazılıb); `[]int{...}` = **slice** (boş mötərizə).
- **nil slice** → mövcud olmayan kolleksiya (məs. funksiyada xəta halında return)
- **empty slice** → boş kolleksiya (məs. sorğu 0 nəticə qaytaran DB)
- `append`, `len`, `cap` hər ikisində eyni işləyir

**Slicing əməliyyatı:**
```go
slice := []int{10, 20, 30, 40, 50}
newSlice := slice[1:3]    // length 2, capacity 4
```

**Formulalar:**
```
slice[i:j]  (underlying array capacity = k):
  Length:   j - i
  Capacity: k - i

slice[1:3], k=5 → length = 3-1 = 2, capacity = 5-1 = 4
```

**Təhlükə:** iki slice eyni underlying array-i paylaşır!
```go
newSlice[1] = 35   // slice-ın index 2-si də dəyişir!
newSlice[3] = 45   // RUNTIME PANIC: index out of range
                   // (length-dən kənar çıxış yoxdur, capacity yalnız append ilə)
```

### 3. append — slice böyüməsi
**Nədir:** Built-in funksiya; slice-a element əlavə edir, yeni slice qaytarır.

**Necə işləyir:**
1. Capacity varsa → element length-ə daxil edilir, underlying array dəyişmir
2. Capacity yoxdursa → **yeni underlying array** ayrılır, dəyərlər kopyalanır, yeni element yazılır

```go
// Capacity mövcuddur:
slice := []int{10, 20, 30, 40, 50}
newSlice := slice[1:3]          // cap 4, len 2
newSlice = append(newSlice, 60) // len 3; slice-ın index 3-ü (40) → 60 olur!

// Capacity doludur:
slice := []int{10, 20, 30, 40}  // len=cap=4
newSlice := append(slice, 50)   // yeni array! capacity 2x → 8
```

**Böyümə alqoritmi:** capacity < 1000 → **2x**; ≥ 1000 → **1.25x (25%)**. (Dil inkişafı ilə dəyişə bilər.)

**Variadic append:**
```go
s1 := []int{1, 2}
s2 := []int{3, 4}
fmt.Printf("%v\n", append(s1, s2...))   // [1 2 3 4]
```
`s2...` → slice-ın elementlərını ayrı-ayrı arqument kimi açır.

### 4. Üç indeksli slicing — capacity məhdudiyyəti
**Nədir:** `slice[i:j:k]` — üçüncü indeks yeni slice-ın capacity-sini məhdudlaşdırır.

**Formulalar:**
```
slice[i:j:k]:
  Length:   j - i
  Capacity: k - i

source[2:3:4] → length = 1, capacity = 2
```

**Nəyə lazımdır:** Append-in underlying array-i korlamasının qarşısını almaq — detach pattern:
```go
source := []string{"Apple", "Orange", "Plum", "Banana", "Grape"}
slice := source[2:3:3]        // len = cap = 1
slice = append(slice, "Kiwi") // capacity dolu → YENİ array ayrılır
                               // "Banana" (index 3) dəyişməz qalır!
```
Üçüncü indeks olmasaydı, append "Banana"-nı "Kiwi" ilə əvəz edərdi.

**Xəta nümunəsi:** `source[2:3:6]` → mövcud tutumdan böyük → `panic: slice bounds out of range`.

### 5. range ilə iterasiya
**Kitabdan kod nümunəsi:**
```go
slice := []int{10, 20, 30, 40}
for index, value := range slice {
    fmt.Printf("Index: %d Value: %d\n", index, value)
}
```

**Vacib:** `value` — elementin **kopyasıdır**, referans deyil! `&value` hər iterasiyada eyni ünvandır (kopya dəyişəninin ünvanı). Elementin real ünvanı: `&slice[index]`.

```go
// Index lazım deyilsə:
for _, value := range slice { ... }

// Kontrol lazımdırsa — klassik for:
for index := 2; index < len(slice); index++ { ... }
```

**Multidimensional slices:**
```go
slice := [][]int{{10}, {100, 200}}
slice[0] = append(slice[0], 20)   // index 0-a yeni slice təyin edilir
```

**Slice-ın funksiyaya ötürülməsi:** 24 byte kopyalanır, underlying array YOX. Pointer lazım deyil — sadəcə slice-ı qaytar: `slice = foo(slice)`.

### 6. Maps (açar/dəyər kolleksiyası)
**Nədir:** Sırasız (unordered) key/value cütləri kolleksiyası — hash table üzərində implementasiya.

**Necə işləyir (internals):** Açar hash funksiyasından keçir → hash-in aşağı bitləri (LOB — low order bits) **bucket** seçir. Bucket daxilində: yuxarı bitlərdən (HOB) ibarət array (elementləri ayırmaq üçün) + açar/dəyərlərin bir yerdə "packed" olduğu byte array. Yaxşı paylanma → 10.000 elementli map-də cəmi ~8 cütl baxmaqla axtarış.

**Yaratma:**
```go
// make
dict := make(map[string]int)

// Map literal — idiomatik üsul
dict := map[string]string{"Red": "#da1337", "Orange": "#e95a22"}

// Slice VALUE ola bilər:
dict := map[int][]string{}

// Slice AÇAR ola BİLMƏZ:
dict := map[[]string]int{}   // COMPILE XƏTASI: invalid map key type []string
```

**Açar qaydası:** Açar `==` ifadəsində istifadə oluna bilən hər hansı built-in/struct tipi ola bilər. Slice və function açar ola bilməz; slice saxlayan struct da ola bilməz.

**Əməliyyatlar:**
```go
colors := map[string]string{}
colors["Red"] = "#da1337"        // əlavə

// NIL map-ə yaz → panic!
var colors map[string]string
colors["Red"] = "#da1337"        // panic: assignment to entry in nil map

// Açar mövcudluğu — 2 dəyərli lookup (TÖVSİYƏ):
value, exists := colors["Blue"]
if exists { fmt.Println(value) }

// Zero-value yoxlaması (yalnız zero value yanlış nəticə deyilsə):
value := colors["Blue"]
if value != "" { fmt.Println(value) }

// Iterasiya — SIRASI QEYRİ-MÜƏYYƏNDİR (hər dəfə fərqli ola bilər!)
for key, value := range colors {
    fmt.Printf("Key: %s Value: %s\n", key, value)
}

// Silmə:
delete(colors, "Coral")
```

**Map funksiyaya ötürülməsi:** Kopyalanmır! Dəyişikliklər bütün referanslarda görünür (slice kimi ucuz, amma reference davranışı):
```go
func removeColor(colors map[string]string, key string) {
    delete(colors, key)
}
// main-də colors["Coral"] artıq yoxdur
```

## len / cap cədvəli
| Funksiya | Array | Slice | Map |
|----------|-------|-------|-----|
| `len()` | ✅ uzunluq | ✅ uzunluq | ✅ cüt sayı |
| `cap()` | ✅ | ✅ yalnız slice-da | ❌ map-də yoxdur |

## Əsas terminlər
- Array (massiv — sabit uzunluqlu ardıcıl struktur)
- Slice (hissə — pointer + length + capacity)
- Underlying Array (altta yatan massiv)
- Length / Capacity (uzunluq / tutum)
- append (böyütmə built-in funksiyası)
- make (ilkinləşdirmə built-in funksiyası)
- nil slice / empty slice (sıfır / boş hissə)
- Hash Table (xəş cədvəli)
- Bucket (xəş çəni)
- LOB / HOB (aşağı / yuxarı bitlər)
- Contiguous Memory (ardıcıl yaddaş)
- Pass by Value (dəyərlə ötürmə)
- Variadic Function (dəyişən saylı arqumentli funksiya)

## Praktik nəticə
- Kolleksiyalar üçün array yox, **slice** istifadə et (idiomatik); sabit ölçü məcburiyyəti yoxdursa array-ə ehtiyac yoxdur.
- Slice paylaşan iki slice bir-birinin datasını görür — təcrid üçün `s[i:j:j]` (len=cap) detach pattern.
- Böyük array-i funksiyaya pointer ilə ötür; slice-lar isə özü-özünü idarə edir (24 byte).
- Map-də açar mövcudluğunu **iki dəyərli lookup** ilə yoxla (`value, exists := m[k]`).
- Map iterasiyasının sırasına etibar etmə — sıralı nəticə lazımdırsa açarları ayrıca sırala.
- Map-i funksiyaya ötürəndə dəyişikliklər orijinalda görünür — bunu ya bilərəkdən istifadə et, ya da kopya yarat.

## Mənbə
Pages: 78-108 (PDF), book pages 57-87
