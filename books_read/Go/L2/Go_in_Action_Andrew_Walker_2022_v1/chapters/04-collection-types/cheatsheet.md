# Cheat Sheet — Chapter 4: Collection Types

## Array Declarations

```go
var intArray [5]int                          // var + zero values
intArrayShort := [5]int{1, 2, 3, 4, 5}       // Short declaration + literal
intArraySparse := [5]int{1: 2, 3: 4}         // Sparse: specific indices
intArrayAuto := [...]int{5, 4, 3, 2, 1}      // Automatic length
intArrayAutoSparse := [...]int{0: 5, 3: 2}   // Auto + sparse
```

**Nə edir:** Müxtəlif formatlarda array elan edir.

**Sub-kod izahı:**
- `[5]int` — 5 elementli int massiv
- `[...]int` — Compiler ölçünü hesablayır
- `{1: 2, 3: 4}` — Yalnız 1-ci və 3-cü indeksi initialize edir

**Mənbə:** Chapter 4, page 104

---

## Array Type And Length

```go
var (
    fourArray [4]int
    fiveArray [5]int
)
fiveArray = fourArray // error: cannot use [4]int as [5]int
fmt.Println(len(intArray)) // 3 (full size)
```

**Nə edir:** Array ölçüsü tipin tərkib hissəsidir — fərqli ölçülər assign edilə bilməz.

**Mənbə:** Chapter 4, page 104

---

## Array Element Access

```go
intArray := [5]int{1, 2, 3, 4, 5}
firstValue := intArray[0]        // 1
index := 0
integerValue := intArray[index]  // 1
intArray[3] = 6                  // Modify
```

**Nə edir:** Square bracket (kvadrat mötərizə) ilə elementlərə müraciət və dəyişiklik.

**Mənbə:** Chapter 4, page 104

---

## Array of Pointers

```go
ptrArray := [5]*int{0: new(int), 1: new(int)}
*ptrArray[0] = 10
*ptrArray[1] = 20
valOfPtr := *ptrArray[0] // 10
```

**Nə edir:** Pointerlar massivi — hər element bir pointer saxlayır.

**Mənbə:** Chapter 4, page 104

---

## Iterating Arrays

```go
array := [...]int{1, 2, 3, 4}

// Range: indices only
for i := range array { fmt.Println(i) }

// Range: indices + values
for i, item := range array { fmt.Println(i, item) }

// Range: values only (blank identifier)
for _, item := range array { fmt.Println(item) }

// Traditional for loop
for i := 0; i < len(array); i++ {
    fmt.Println(array[i])
}
```

**Nə edir:** Array üzərində iterasiya nümunələri.

**Sub-kod izahı:**
- `range` — Avtomatik iterasiya
- `_` — Blank identifier, indeksı ignore et

**Mənbə:** Chapter 4, page 104

---

## Array Assignment (By Value)

```go
var colors [5]string
favoriteColors := [5]string{"Red", "Blue", "Green", "Yellow", "Pink"}
colors = favoriteColors // Full copy!
```

**Nə edir:** Array dəyər tipidir — kopyalanır.

**Mənbə:** Chapter 4, page 104

---

## Multidimensional Arrays

```go
var array [4][2]int
array = [4][2]int{{10, 11}, {20, 21}, {30, 31}, {40, 41}}
array = [4][2]int{1: {20, 21}, 3: {40, 41}}

// Access inner elements
array2D := [2][2]int{{10, 11}, {20, 21}}
array1D := array2D[1]     // [20 21]
intValue := array1D[0]     // 20
```

**Nə edir:** Array of arrays — 2D massiv elanı və müraciət.

**Mənbə:** Chapter 4, page 104

---

## Passing Arrays to Functions

```go
type image4K [829400]int // 3840 * 2160

// By value (expensive — full copy)
func processPixels(pixels image4K) image4K { ... }

// By pointer (efficient)
func processPixels(pixels *image4K) { ... }
```

**Nə edir:** Böyük massivləri funksiyaya ötürmək üçün pointer istifadə edin.

**Mənbə:** Chapter 4, page 104

---

## Slice Declarations

```go
slice := []int{10, 20, 30}              // Literal
slice2 := make([]int, 0, 10)            // make: len=0, cap=10
slice3 := make([]int, 5)                // make: len=5, cap=5
```

**Nə edir:** Slice elanı — literal və ya `make` ilə.

**Sub-kod izahı:**
- `make([]T, len, cap)` — Len və cap təyin edərək yaradır

**Mənbə:** Chapter 4, page 104

---

## Slice Append

```go
slice := []int{10, 20, 30}
slice = append(slice, 40, 50)
fmt.Println(slice) // [10 20 30 40 50]

slice2 := []int{60, 70}
slice = append(slice, slice2...) // Spread operator
fmt.Println(slice) // [10 20 30 40 50 60 70]
```

**Nə edir:** Slice-ə elementlər əlavə edir, lazımsa capacity böyüyür.

**Sub-komanda/flag izahı:**
- `append` — Həmişə yeni slice header qaytarır
- `slice2...` — Spread operator (slice-i ayrı elementlərə açır)

**Mənbə:** Chapter 4, page 104

---

## Slice Expressions

```go
colors := []string{"Red", "Blue", "Green", "Yellow", "Pink"}
sub := colors[1:3]    // ["Blue" "Green"] — index 1 inclusive, 3 exclusive
tail := colors[2:]    // ["Green" "Yellow" "Pink"] — 2-dən sonuna
head := colors[:3]    // ["Red" "Blue" "Green"] — 0-dan 3-ə qədər
full := colors[:]     // Full copy view
```

**Nə edir:** Slice-dən alt kolleksiya (view) yaradır.

**Sub-komanda/flag izahı:**
- `[low:high]` — low inclusive, high exclusive
- Capacity — View-in capacity orijinaldən qalan hissədir

**Mənbə:** Chapter 4, page 104

---

## Nil vs Empty Slice

```go
var nilSlice []int              // nil
emptySlice := []int{}           // empty, not nil
fmt.Println(nilSlice == nil)    // true
fmt.Println(emptySlice == nil)  // false
fmt.Println(len(nilSlice))      // 0
fmt.Println(len(emptySlice))    // 0
```

**Nə edir:** Nil slice (`nil`) və empty slice (`[]`) fərqini göstərir.

**Mənbə:** Chapter 4, page 104

---

## Map Declarations

```go
groupNouns := map[string]string{
    "eagle": "convocation",
    "cat":   "clowder",
    "dog":   "pack",
}

// Access
group := groupNouns["cat"]   // "clowder"
group, found := groupNouns["cat"] // found: true

// Initialize with make
m := make(map[string]int)
m["key"] = 42

// Delete
delete(groupNouns, "cat")
```

**Nə edir:** Map elanı, daxil etmə, silmə, lookup.

**Sub-komanda/flag izahı:**
- `map[KeyType]ValueType` — Map tipi
- `make(map[K]V)` — Initialize (başlat)
- `delete(map, key)` — Key-value cütünü sil
- `value, found := map[key]` — Two-value form, found mövcudluq

**Mənbə:** Chapter 4, page 104

---

## Map Iteration

```go
for animal, group := range groupNouns {
    fmt.Printf("A group of %ss is called a %s.\n", animal, group)
}

// Ignore key
for _, group := range groupNouns {
    fmt.Println(group)
}
```

**Nə edir:** Range ilə map üzərində iterasiya.

**Mənbə:** Chapter 4, page 104

---

## Map Filtering

```go
familyAges := map[string]int{
    "Grandpa": 83,
    "Homer":   36,
    "Bart":    10,
    "Lisa":    8,
}

for person, age := range familyAges {
    if age < 18 {
        delete(familyAges, person)
    }
}
fmt.Println(familyAges) // map[Grandpa:83 Homer:36]
```

**Nə edir:** Map-dən şərtə uyğun gəlməyən elementləri silir.

**Mənbə:** Chapter 4, page 104

---

## Map with Struct Keys

```go
type person struct {
    firstName string
    lastName  string
    age       int
}

favoriteFoods := map[person][]string{}
andy := person{firstName: "Andy", lastName: "Walker", age: 42}
favoriteFoods[andy] = []string{"jambalaya", "bulgogi"}
```

**Nə edir:** Struct-tipli açar ilə map yaradır.

**Mənbə:** Chapter 4, page 104

---

## len() and cap() for Slices

```go
slice := make([]int, 3, 5)
fmt.Println(len(slice)) // 3
fmt.Println(cap(slice)) // 5
```

**Nə edir:** Slice-in uzunluğu və tutumunu göstərir.

**Sub-komanda/flag izahı:**
- `len` — Mövcud element sayı
- `cap` — Azad edilə biləcək maksimum element sayı (backing array ölçüsü)

**Mənbə:** Chapter 4, page 104

---

## %T Print Verb

```go
var array3D [4][3][2]string
fmt.Printf("array3D is a %T\n", array3D)
// output: array3D is a [4][3][2]string
```

**Nə edir:** Dəyişənin tipini string olaraq göstərir.

**Mənbə:** Chapter 4, page 104
