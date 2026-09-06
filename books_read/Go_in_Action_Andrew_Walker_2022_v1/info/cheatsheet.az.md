# Ümumi Cheat Sheet — Go in Action, Second Edition

Bu cheat sheet kitabın chapter 1-4-dən bütün kod nümunələrini, komandalarını və konfiqurasiyalarını əhatə edir.

---

## Chapter 1 — Introducing Go

### `go run main.go`

**Nə edir:** Go faylını kompilyə edib birbaşa icra edir.

**Sub-komanda/flag izahı:**
- `go` → Go command-line aləti
- `run` → Faylı kompilyə edib icra et
- `main.go` → Giriş faylı

**Mənbə:** Chapter 1, page 6

---

### `package main`

**Nə edir:** Go proqramının icra edilə bilən package (paket) olduğunu təyin edir.

**Sub-kod izahı:**
- `package` → Package declaration (paket elanı)
- `main` → Executable program (icra edilən proqram) üçün xüsusi ad

**Mənbə:** Chapter 1, page 6

---

### `import "fmt"`

**Nə edir:** Formatlama və çıxış üçün standart kütubxəni daxil edir.

**Mənbə:** Chapter 1, page 6

---

### `func main()`

**Nə edir:** Program işə düşəndə ilk çağırılan funksiya.

**Mənbə:** Chapter 1, page 6

---

### `fmt.Println(...)`

**Nə edir:** Konsola mətn və ya dəyər çıxışı verir, avtomatik olaraq yeni sətir əlavə edir.

**Sub-kod izahı:**
- `fmt` → Format paketi
- `Println` → Print line (sətir çap et)

**Mənbə:** Chapter 1, page 6

---

### `go makeMeConcurrent()`

**Nə edir:** Funksiyonu goroutine (gözaltı proses) kimi işə salır, paralel icra olunur.

**Sub-komanda/flag izahı:**
- `go` → Goroutine yaratmaq üçün açar sözcük

**Mənbə:** Chapter 1, page 6

---

### `time.Sleep(1 * time.Second)`

**Nə edir:** Proqramı verilən müddət üçün dayandırır.

**Sub-komanda/flag izahı:**
- `time` → Zaman əməliyyatları paketi
- `Sleep` → Müddət üçün dayandır
- `1 * time.Second` → 1 saniyə

**Mənbə:** Chapter 1, page 6

---

### `make(chan string)`

**Nə edir:** String (mətn) tipli kanal (channel) yaradır.

**Sub-komanda/flag izahı:**
- `make` → Yeni map, slice və ya channel yaradır
- `chan` → Kanal tipi
- `string` → Kanalda ötürüləcək məlumat tipi

**Mənbə:** Chapter 1, page 6

---

### `var number int`

**Nə edir:** Tam ədəd tipində (int) dəyişən elan edir, sıfır dəyəri (0) təyin edir.

**Mənbə:** Chapter 1, page 6

---

### `type Person struct { ... }`

**Nə edir:** Yeni struct (strukt) tipi yaradır, müxtəlif sahələri birləşdirir.

**Sub-kod izahı:**
- `type` → Yeni tip elanı
- `struct` → Strukt tipləri qrupu

**Mənbə:** Chapter 1, page 6

---

### `func SumNumbers[N int | float64](numberSlice ...N) N`

**Nə edir:** Generics (generik) funksiya — int və ya float64 tipləri üçün işləyir.

**Sub-kod izahı:**
- `[N int | float64]` → Type parameter (tip parametri) N, constraint daxilində
- `numberSlice ...N` → Variadic (dəyişən saylı) parametr

**Mənbə:** Chapter 1, page 6

---

## Chapter 2 — Diving Into Go

### `go mod init gobook/wordcount`

**Nə edir:** Yeni Go module yaradır və `go.mod` faylı yaradır.

**Mənbə:** Chapter 2, page 25

---

### `go build`

**Nə edir:** Go fayllarını kompilyə edir, executable binary yaradır.

**Mənbə:** Chapter 2, page 25

---

### `gofmt -d main.go`

**Nə edir:** Faylın formatlaşdırma fərqlərini göstərir (diff).

**Sub-komanda/flag izahı:**
- `-d` → Diff (fərq) rejimi

**Mənbə:** Chapter 2, page 25

---

### `go doc -all strings`

**Nə edir:** Package dokumentasiyasını konsola çıxarır.

**Sub-komanda/flag izahı:**
- `-all` → Ətraflı məlumat

**Mənbə:** Chapter 2, page 25

---

### `var numSpaces int`

**Nə edir:** Tam ədəd tipində dəyişən elan edir, sıfır dəyəri (0) təyin edir.

**Mənbə:** Chapter 2, page 25

---

### `text := "let's count some words!"`

**Nə edir:** Short declaration operatoru ilə string dəyişəni yaradır.

**Sub-komanda/flag izahı:**
- `:=` → Elan və təyinat eyni anda
- Type inference (tip çıxarma)

**Mənbə:** Chapter 2, page 25

---

### `for i := 0; i < len(text); i++ { }`

**Nə edir:** Klassik 3-hissəli for döngüsü.

**Mənbə:** Chapter 2, page 25

---

### `if err != nil { log.Println(err); os.Exit(1) }`

**Nə edir:** Standart error handling pattern-i.

**Sub-kod izahı:**
- `log.Println` → Timestamp ilə standard error-a yaz
- `os.Exit(1)` → Proqramı uğursuz exit kodu ilə dayandır

**Mənbə:** Chapter 2, page 25

---

### `bufio.NewScanner(file)`

**Nə edir:** Fayl üçün buferli scanner yaradır.

**Mənbə:** Chapter 2, page 25

---

### `scanner.Split(bufio.ScanWords)`

**Nə edir:** Scanner-ın tokenization növünü sözlərə dəyişir.

**Mənbə:** Chapter 2, page 25

---

### `for _, filename := range os.Args[1:] { }`

**Nə edir:** Range döngüsü ilə program argumentlərini iterasiya edir.

**Sub-kod izahı:**
- `os.Args[1:]` — Slice expression, 1-ci indekstdən sonuna

**Mənbə:** Chapter 2, page 25

---

## Chapter 3 — Primitive Types And Operators

### Integer Declarations

```go
var myInt int
i := 3
u := uint64(4)
decInt := 1000
hexInt := 0x3E8
octInt := 01750
binInt := 0b1111101000
withSep := 1_000
```

**Mənbə:** Chapter 3, page 67

---

### Unsigned Integer Wraparound

```go
var u uint64 // 0
u = u - 1
fmt.Println(u) // 18446744073709551615
```

**Mənbə:** Chapter 3, page 67

---

### Float Declarations

```go
var doubleFloat float64
f := 12.1
g := float32(12.1)
floatVal = 12e0
negativeFloat := -12.0
```

**Mənbə:** Chapter 3, page 67

---

### Special Floating-Point Values

```go
f := 2.0
posInf := math.Pow(f, 10_000)
notANumber := math.Log(-f)
fmt.Println(math.IsInf(posInf, 0))  // true
fmt.Println(math.IsNaN(notANumber)) // true
```

**Mənbə:** Chapter 3, page 67

---

### Complex Numbers

```go
cmplxValue1 := complex(1.1, 2.2) // complex128
cmplxValue3 := 1 + 2i            // complex128
fmt.Println(real(cmplxValue1), imag(cmplxValue1)) // "1.1 2.2"
```

**Mənbə:** Chapter 3, page 67

---

### Arithmetic Operators

```go
a := 7
b := 3
i = a + b   // 10
i = a - b   // 4
i = a * b   // 21
i = a / b   // 2
i = a % b   // 1

u := uint64(1)
u = u + uint64(i) // Type conversion required
```

**Mənbə:** Chapter 3, page 67

---

### Bitwise Operators

```go
a & b    // AND
a | b    // OR
a ^ b    // XOR
a &^ b   // Bit clear (Go-unikal)
a << 1   // Shift left
a >> 1   // Shift right
```

**Mənbə:** Chapter 3, page 67

---

### Boolean Helper Function

```go
func IsEven(i int) bool {
    return i%2 == 0
}
```

**Mənbə:** Chapter 3, page 67

---

### Struct Declarations

```go
type person struct {
    name string
    age  int
}

andy := person{
    name: "Andy",
    age:  42,
}
```

**Mənbə:** Chapter 3, page 67

---

### Pointer Operations

```go
var intPtr *int
intValue := 0
intPtr = &intValue
fmt.Println(*intPtr) // 0
*intPtr = 1
fmt.Println(intValue) // 1

newIntPtr := new(int)
newPerson := &person{name: "Andy", age: 42}
```

**Mənbə:** Chapter 3, page 67

---

### String Types

```go
basicStr := "Hello, Gophers!"
rawString := `Hello\n\u5730\u9F20`
rawStrWithNewlines := `I can
span multiple lines`
```

**Mənbə:** Chapter 3, page 67

---

### String Operations

```go
fmt.Println("hello" == "hello") // true
strHello += " Gophers!"         // Concatenation
```

**Mənbə:** Chapter 3, page 67

---

### len() on Strings

```go
asciiCharStr := "easy, right?"
fmt.Println(len(asciiCharStr))      // 12 (bytes)

unicodeCharStr := "地鼠"
fmt.Println(len(unicodeCharStr))    // 6 (bytes)
```

**Mənbə:** Chapter 3, page 67

---

### Range Loop on String

```go
for i, r := range unicodeCharStr {
    fmt.Printf("%d:%s ", i, string(r))
}
// output: 0:地 3:鼠
```

**Mənbə:** Chapter 3, page 67

---

### []rune Conversion

```go
characters := []rune(unicodeCharStr)
for i, r := range characters {
    fmt.Printf("%d:%s ", i, string(r))
}
// output: 0:地 1:鼠
```

**Mənbə:** Chapter 3, page 67

---

### fmt.Printf %x

```go
fmt.Printf("%x\n", "A")    // 41
fmt.Printf("%x\n", "地")   // e1beb8
```

**Mənbə:** Chapter 3, page 67

---

## Chapter 4 — Collection Types

### Array Declarations

```go
var intArray [5]int
intArrayShort := [5]int{1, 2, 3, 4, 5}
intArraySparse := [5]int{1: 2, 3: 4}
intArrayAuto := [...]int{5, 4, 3, 2, 1}
```

**Mənbə:** Chapter 4, page 104

---

### Array Type And Length

```go
var (
    fourArray [4]int
    fiveArray [5]int
)
fiveArray = fourArray // error: cannot use [4]int as [5]int
fmt.Println(len(intArray)) // 3
```

**Mənbə:** Chapter 4, page 104

---

### Array Element Access

```go
intArray := [5]int{1, 2, 3, 4, 5}
firstValue := intArray[0]   // 1
intArray[3] = 6             // Modify
```

**Mənbə:** Chapter 4, page 104

---

### Array of Pointers

```go
ptrArray := [5]*int{0: new(int), 1: new(int)}
*ptrArray[0] = 10
valOfPtr := *ptrArray[0] // 10
```

**Mənbə:** Chapter 4, page 104

---

### Iterating Arrays

```go
for i := range array { ... }
for i, item := range array { ... }
for _, item := range array { ... }
for i := 0; i < len(array); i++ { ... }
```

**Mənbə:** Chapter 4, page 104

---

### Multidimensional Arrays

```go
var array [4][2]int
array = [4][2]int{{10, 11}, {20, 21}, {30, 31}, {40, 41}}
array1D := array2D[1]    // [20 21]
intValue := array1D[0]    // 20
```

**Mənbə:** Chapter 4, page 104

---

### Slice Declarations

```go
slice := []int{10, 20, 30}
slice2 := make([]int, 0, 10)
slice3 := make([]int, 5)
```

**Mənbə:** Chapter 4, page 104

---

### Slice Append

```go
slice := []int{10, 20, 30}
slice = append(slice, 40, 50)
slice2 := []int{60, 70}
slice = append(slice, slice2...) // Spread operator
```

**Mənbə:** Chapter 4, page 104

---

### Slice Expressions

```go
colors := []string{"Red", "Blue", "Green", "Yellow", "Pink"}
sub := colors[1:3]    // ["Blue" "Green"]
tail := colors[2:]    // ["Green" "Yellow" "Pink"]
head := colors[:3]    // ["Red" "Blue" "Green"]
```

**Mənbə:** Chapter 4, page 104

---

### Nil vs Empty Slice

```go
var nilSlice []int            // nil
emptySlice := []int{}         // empty, not nil
fmt.Println(nilSlice == nil)  // true
fmt.Println(len(nilSlice))    // 0
```

**Mənbə:** Chapter 4, page 104

---

### Map Declarations

```go
groupNouns := map[string]string{
    "eagle": "convocation",
    "cat":   "clowder",
}

group, found := groupNouns["cat"] // found: true
m := make(map[string]int)
m["key"] = 42
delete(groupNouns, "cat")
```

**Mənbə:** Chapter 4, page 104

---

### Map Iteration

```go
for animal, group := range groupNouns {
    fmt.Printf("A group of %ss is called a %s.\n", animal, group)
}
```

**Mənbə:** Chapter 4, page 104

---

### Map Filtering

```go
for person, age := range familyAges {
    if age < 18 {
        delete(familyAges, person)
    }
}
```

**Mənbə:** Chapter 4, page 104

---

### len() and cap() for Slices

```go
slice := make([]int, 3, 5)
fmt.Println(len(slice)) // 3
fmt.Println(cap(slice)) // 5
```

**Mənbə:** Chapter 4, page 104

---

### %T Print Verb

```go
fmt.Printf("array3D is a %T\n", array3D)
// output: array3D is a [4][3][2]string
```

**Mənbə:** Chapter 4, page 104
