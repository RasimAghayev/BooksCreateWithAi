# Chapter 2 — Arrays and Algorithms for Searching and Sorting (səh. 85-159)

## Bu fəsil nədən bəhs edir?

İlk data strukturu — array-lər: yaddaşda yerləşmə, əməliyyatlar, Go-da
istifadə. Slices (dinamik ölçülü "flexible arrays"), çoxölçülü matrislər,
method və interface anlayışları. İkinci hissə: axtarış (sequential, binary)
və sıralama (insertion, selection, bubble, quick) alqoritmləri Go
implementasiyaları ilə + `sort` paketi.

## Əsas fikirlər

### 1. Array nədir?
- **Linear + homogen** struktur; sonlu sayda element
- **Arranged:** hər elementin mövqeyi məlumdur — **index** mövqeyi təyin edir
- Yaddaşda **sequential** saxlanılır (fasiləsiz; element başına 4 location
  nümunəsi; ünvanlar hex formatında)
- Əməliyyatlar (index vasitəsilə): element dəyərini **seç** (get), **dəyiş** (set)

### 2. Array-lər Go-da
```go
var a [8]int                          // zero value-larla: 8 sıfır
var a = [8]int{1, 18, 5, 27, 25, 8, 21, 9}  // initializasiya
var b = a[3]                          // get: index 3
a[3] = 7                              // set
```
- Uzunluq tipin **inteqral hissəsidir**, dəyişə bilməz — real problemlərdə
  məhdudiyyət → həll: **slice**

### 3. Slice — "flexible array"
- Dinamik ölçü; mahiyyətdə **underlying array-ə pointer**
- Array ilə eyni əməliyyatlar (get/set)

**Go-da slice:**
```go
var s []int = a[1:4]        // a[1], a[2], a[3] — half-open [low:high)
fmt.Println(s[0])           // 18
fmt.Println(len(s))         // 3
```
- **Half-open range:** ilk daxil, son çıxarılır
- Slice real datanı saxlamır — underlying array-ə **istinad** edir: slice-də
  dəyişiklik = array-də dəyişiklik; eyni array-ə istinad edən bütün slice-lərə
  təsir edir

**Bound defaultları:** low=0, high=len(slice):
```go
s1 := s[1:4]  // [18 5 27]
s2 := s[:4]   // [1 18 5 27]
s3 := s[1:]   // [18 5 27 8 25 9 21]
s4 := s[:]    // hamısı
```

**Length vs Capacity:**
| Attribute | Məna |
|---|---|
| len(s) | slice-dəki element sayı |
| cap(s) | underlying array-dəki element sayı (slice-in ilk elementindən sayılır) |
- a[1:4] nümunəsində: len=3, cap=7; capacity çatarsa length uzadıla bilər
- **nil slice:** default dəyər; len=cap=0, underlying array yoxdur

**make() ilə yaratma:**
```go
s := make([]int, 0, 5)   // type, length, capacity
```

**append() — sona əlavə:**
```go
s = append(s, 3, 21, 12, 30)
```
- Nəticə: original + yeni elementlər; underlying array kiçikdirsə daha böyük
  yeni array ayrılır → nəticə eyni slice dəyişəninə təyin edilir

### 4. Çoxölçülü array-lər (matrislər)
```go
var matrix [3][3]int                          // 3x3 zero
var matrix = [3][3]int{{1,18,5},{27,25,8},{21,9,12}}
var n = matrix[1][2]                          // get
a[2][0] = 7                                   // set
var m1 [3][]int                               // sütun ötürülüb
var m2 [][]int                                // hər ikisi
matrix = append(matrix, []int{1, 18, 5}, []int{27, 8, 25})  // sətir əlavəsi
```
- **Matrix (2D array):** ən yaygın; kvadrat = hər iki ölçü eyni (3x3)
- Yaddaşda **sətirlərlə ardıcıl** (row-major); sütunla da ola bilər

### 5. Method və Interface
**Method** — tip üzərində elan olunan funksiya (receiver arqumenti ilə):
```go
type Rectangle struct{ a, b int }

func (r Rectangle) Area() int      { return r.a * r.b }
func (r Rectangle) Perimeter() int { return 2*r.a + 2*r.b }
```
- Receiver `func` sözü ilə metod adı arasında mötərizədə
- **Pointer receiver:** metodu çağıran dəyişəni dəyişir (praktikada çox istifadə
  olunur)

**Interface** — metod elanları dəsti; tip bütün metodları implement edərək
implementasiya edir (xüsusi açar sözü YOXDUR):
```go
type Shape interface {
    Area() int
    Perimeter() int
}
// Rectangle və Square hər ikisi Shape-i implement edir
```

### 6. Axtarış alqoritmləri
Axtarış = identifikasiya (**key**) əsasında datasetdə datanı tapmaq.
Nəticə: tapıldı / tapılmadı.

**Sequential search — O(n):**
- Key hər elementlə müqayisə olunur → uğurlu (match) və ya uğursuz (hamısı
  yoxlanıldı)
- Ən az effektiv; lakin **sıralanmamış array-də yeganə seçim**
- Sıralı array-də effektiv artır (key-dən böyük element görünəndə dayanmaq)
```go
func SeqSearch(array []int, key int) int {
    for i, elem := range array {   // for-range: index + element copy
        if key == elem {
            return i
        }
    }
    return -1
}
```

**Binary search — O(log n):**
- Divide and conquer; **yalnız sıralı array-də**
- Sequential-in problemini həll edir: hər addımda element QRUPU atılır
- Hər iterasiyada 3 dəyər: **low, mid=(low+high)/2, high**
  - array[mid] == key → return mid
  - key < array[mid] → high = mid-1 (yuxarı yarım atılır)
  - key > array[mid] → low = mid+1 (aşağı yarım atılır)
- low >= high olana qədər; uğursuzsa -1
```go
func BinSearch(array []int, key int) int {
    low, high := 0, len(array)-1
    for low <= high {
        mid := (low + high) / 2
        if key == array[mid] {
            return mid
        } else if key < array[mid] {
            high = mid - 1
        } else {
            low = mid + 1
        }
    }
    return -1
}
```

### 7. Go `sort` paketi — axtarış
```go
func SearchInts(a []int, x int) int
func SearchFloat64s(a []float64, x float64) int
func SearchStrings(a []string, x string) int
func Search(n int, f func(int) bool) int    // əsas funksiya
```
- Hamısı binary search; match yoxdursa x-in **daxil edilməli olduğu index**
  qaytarılır; slice artan sıralı olmalıdır
- `SearchInts` wrapper-dir: `Search(len(a), func(i int) bool { return a[i] >= x })`
- Azalan sıralama üçün `>=` əvəzinə `<=`

### 8. Sıralama alqoritmləri
Sıralama = dataseti müəyyən sıraya düzme; əksəriyyəti birbaşa müqayisəyə
əsaslanır.

**Insertion sort — O(n²):**
- Array 2 hissəyə bölünür: **sorted + unsorted**; hər iterasiyada unsorted-dan
  1 element sorted-a **uyğun yerə daxil** edilir
- Başlanğıcda yalnız index 0 sorted-a aiddir
```go
func InsertionSort(array []int) []int {
    for i := 1; i < len(array); i++ {
        k := array[i]
        j := i - 1
        for j >= 0 && array[j] > k {
            array[j+1] = array[j]   // sağa sürüşdür
            j--
        }
        array[j+1] = k
    }
    return array
}
```

**Selection sort — O(n²):**
- Hər iterasiyada unsorted hissədən **ən kiçik** tapılır → sorted hissənin
  sonuna qoyulur (ilk elementlə yer dəyişir)
```go
func SelectionSort(array []int) []int {
    for i := 0; i < len(array)-1; i++ {
        min, pos := array[i], i
        for j := i + 1; j < len(array); j++ {
            if array[j] < min {
                min, pos = array[j], j
            }
        }
        array[pos] = array[i]
        array[i] = min
    }
    return array
}
```

**Bubble sort — O(n²):**
- Ən sadə, ən ineffectiv; qonşu elementlər müqayisə olunur, sırada deyilsə
  dəyişdirilir; hər iterasiyada ən böyük element sona "su üstünə çıxır"
- Adı buradan: elementlər baloncuk kimi üzdən yuxarı
```go
func BubbleSort(array []int) []int {
    n := len(array)
    for i := 0; i < n-1; i++ {
        for j := 0; j < n-i-1; j++ {
            if array[j] > array[j+1] {
                array[j], array[j+1] = array[j+1], array[j]
            }
        }
    }
    return array
}
```

**Quick sort — orta O(n log n), ən pis O(n²):**
- Divide and conquer; **pivot** seçilir (adətən ilk element) → 2 partition:
  aşağı (≤ pivot) və yuxarı (> pivot) → hər partition-a eyni proses (rekursiv)
```go
func QuickSort(array []int, low, high int) []int {
    if low < high {
        array, j := partition(array, low, high)
        QuickSort(array, low, j-1)    // aşağı partition
        QuickSort(array, j+1, high)   // yuxarı partition
    }
    return array
}

func partition(array []int, low, high int) ([]int, int) {
    pivot := array[low]
    i, j := low, high
    for i < j {
        for array[i] <= pivot && i < j { i++ }
        for array[j] > pivot { j-- }
        if i < j {
            array[i], array[j] = array[j], array[i]
        }
    }
    array[low], array[j] = array[j], pivot
    return array, j    // pivotun yekun mövqeyi
}
```
- Effektivlik partition **balansından** asılıdır; ən yaxşı hal — bərabər
  bölünmə; performans array ölçüsündən də asılıdır

**Müqayisə cədvəli (kitabın 2.1):**
| Alqoritm | Komplekslik | Qeyd |
|---|---|---|
| Insertion sort | O(n²) | böyük array-lər üçün yararsız |
| Selection sort | O(n²) | kiçiklər üçün uyğun |
| Bubble sort | O(n²) | tövsiyə olunmur |
| Quick sort | O(n log n) orta | ən pis halda O(n²) |

### 9. Go `sort` paketi — sıralama
```go
func Ints(x []int)
func Float64s(x []float64)
func Strings(x []string)
// hamısı Sort() wrapper-i; artan sıra
sort.Ints(array)   // nümunə
```
- `Sort()` — `sort.Interface` implement edən hər şeyi sıralayır:
```go
type Interface interface {
    Len() int
    Less(i, j int) bool
    Swap(i, j int)
}
type IntSlice []int
func (x IntSlice) Len() int { return len(x) }
func (x IntSlice) Less(i, j int) bool { return x[i] < x[j] }
func (x IntSlice) Swap(i, j int) { x[i], x[j] = x[j], x[i] }
```
- Bərabər elementlərin original sırası lazımdırsa **Stable()** istifadə et
- `sort` funksiyaları **optimallaşdırılmış quick sort**dur (ən pis hal O(n log n))

## Termindirmə (AZ)
- Array — Massiv (linear, homogen, sonlu; indexli)
- Index — İndeks (elementin mövqeyini təyin edən dəyər)
- Slice — Dilim (dinamik ölçülü, underlying array-ə pointer)
- Underlying array — Alt massiv (slice-in istinad etdiyi array)
- Length / Capacity — Uzunluq / Tutum
- Half-open range — Yarı-açıq aralıq (ilk daxil, son çıxarılır)
- Matrix — Matris (2D array; kvadrat = n×n)
- Row-major storage — Sətir-əsaslı yaddaş
- Receiver — Qəbuledici (method-un bağlı olduğu tip)
- Pointer receiver — Pointer qəbuledici (çağıran dəyişəni dəyişir)
- Interface — İnterfeys (metod elanları dəsti)
- Key — Açar (axtarış identifikatoru)
- Sequential search — Ardıcıl axtarış O(n)
- Binary search — İkili axtarış O(log n)
- Pivot — Dayaq (quick sort-da bölünmə nöqtəsi)
- Partition — Bölmə
- Stable sort — Stabil sıralama (bərabərlərin sırası qorunur)

## Kviz sualları
1. `[3]int` və `[]int` fərqi nədir? (array: ölçü tipin daxilində; slice:
   dinamik, underlying array-ə istinad)
2. `a[1:4]` neçə element verir? (3 — index 1,2,3; high çıxarılır)
3. len və cap fərqi nədir? (len = slice-dəki elementlər; cap = underlying
   array-də slice başlanğıcından mövcud yer)
4. `append` nə qaytarır və niyə eyni dəyişənə yazılır? (yeni slice; underlying
   array yenidən ayrıla bilər)
5. Binary search niyə sıralanmamış array-də işləmir? (hansı yarıda olduğunun
   proqnozu sıralılığa əsaslanır)
6. `sort.SearchInts` match tapmasa nə qaytarır? (x daxil edilməli olan index)
7. Hansı sort alqoritmi ən pis halda O(n log n) qarantlıdır? (sort paketinin
   optimallaşdırılmış quick sort-u; bizim sadə QuickSort O(n²)-ə düşə bilər)
8. Selection sort hər iterasiyada nə edir? (unsorted-dən ən kiçiyi tapıb
   sorted-ın sonuna qoyur)
