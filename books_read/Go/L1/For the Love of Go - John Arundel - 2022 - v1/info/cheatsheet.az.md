# For the Love of Go — Cheat Sheet (Azərbaycanca)

John Arundel, Bitfield Consulting 2022 — 15 fəsil, L1 Beginner

## 1. Layihə Qur (Ch1)
```bash
mkdir calculator && cd calculator
go mod init calculator
go test                    # testləri işə sal
go run cmd/calculator/main.go   # proqramı işə sal
go build -o add ./cmd/calculator    # binary yarat
gofmt -w calculator.go     # formatı düzəlt
```

## 2. Test Struktur (Ch1-2)
```go
package calculator_test

import (
    "calculator"
    "testing"
)

func TestAdd(t *testing.T) {
    t.Parallel()
    var want float64 = 4
    got := calculator.Add(2, 2)
    if want != got {
        t.Errorf("want %f, got %f", want, got)
    }
}
```
- _test.go suffiks + Test prefiks + *testing.T + nəticəsiz
- TDD: test → `undefined: X` (compile err) → null impl (`return 0`) →
  FAIL → real impl → PASS

## 3. Table Test (Ch2)
```go
func TestDivide(t *testing.T) {
    t.Parallel()
    type testCase struct {
        a, b float64
        want float64
    }
    testCases := []testCase{
        {a: 2, b: 2, want: 1},
        {a: -1, b: -1, want: 1},
        {a: 10, b: 2, want: 5},
    }
    for _, tc := range testCases {
        got, err := calculator.Divide(tc.a, tc.b)
        if err != nil {
            t.Fatalf("want no error, got %v", err)   // data etibarsız → dayan
        }
        if tc.want != got {
            t.Errorf("Divide(%f, %f): want %f, got %f", tc.a, tc.b, tc.want, got)
        }
    }
}

func TestDivideInvalid(t *testing.T) {
    t.Parallel()
    _, err := calculator.Divide(1, 0)
    if err == nil {
        t.Error("want error for invalid input, got nil")
    }
}
```
- One behaviour, one test! Valid/invalid AYRI

## 4. Float Müqayisə (Ch3)
```go
func closeEnough(a, b, tolerance float64) bool {
    return math.Abs(a-b) <= tolerance
}
if !closeEnough(want, got, 0.001) { ... }
```

## 5. Struct + Validasiya (Ch4, 10)
```go
type Book struct {
    Title           string
    Author          string
    Copies          int
    ID              int
    PriceCents      int     // float YOX — sent-lə int!
    DiscountPercent int
    category        Category // UNEXPORTED — xarici yazış QADAĞA
}

func (b *Book) SetPriceCents(price int) error {   // POINTER mütləq!
    if price < 0 {
        return fmt.Errorf("negative price %d", price)
    }
    b.PriceCents = price
    return nil
}

func (b Book) Category() Category { return b.category }  // getter value
func (b *Book) SetCategory(c Category) error {            // setter pointer
    if !validCategory[c] {
        return fmt.Errorf("unknown category %v", c)
    }
    b.category = c
    return nil
}
```

## 6. iota + Map-as-Set (Ch10)
```go
type Category int

const (
    CategoryAutobiography     Category = iota   // 0, 1, 2...
    CategoryLargePrintRomance
    CategoryParticlePhysics
)

var validCategory = map[Category]bool{
    CategoryAutobiography:     true,
    CategoryLargePrintRomance: true,
    CategoryParticlePhysics:   true,
}   // if validCategory[c] — missing → false!
```

## 7. Kolleksiyalar (Ch6-7)
```go
// SLICE
books := []Book{{Title: "A"}, {Title: "B"}}
first := books[0]; n := len(books)
books = append(books, Book{Title: "C"})    // TƏYİN ET!
books[0].Title = "Z"

// MAP
catalog := map[int]Book{
    1: {ID: 1, Title: "A"},
    2: {ID: 2, Title: "B"},
}
b := catalog[1]                    // olmayan key → ZERO VALUE (panic YOX)
b, ok := catalog[3]                // ok = mövcudluq
catalog[3] = Book{ID: 3}          // əlavə/overwrite
// sahə yazışı QADAĞA: çıxar → dəyiş → geri yaz
tmp := catalog[1]; tmp.Title = "X"; catalog[1] = tmp
```

## 8. Müqayisə (Ch6-7)
```go
import "github.com/google/go-cmp/cmp"          // go get -t
import "github.com/google/go-cmp/cmp/cmpopts"

if !cmp.Equal(want, got) { t.Error(cmp.Diff(want, got)) }  // - want, + got
if !cmp.Equal(want, got,
    cmpopts.IgnoreUnexported(Book{})) { ... }    // unexported sahə varsa

sort.Slice(got, func(i, j int) bool {      // RANDOM map sırasına qarşı
    return got[i].ID < got[j].ID
})
```

## 9. Pointer Kalıbları (Ch9-10)
```go
p := &x              // & = share et
x = *p               // * = dereference (pointee)
mytypes.MyInt(9)     // type conversion; MyInt ≠ int

func (input *MyInt) Double() { *input *= 2 }    // pointer receiver

// nil pointer dereference → PANIC! (zero value = nil)
```

## 10. Wrapping (Ch9)
```go
type MyBuilder struct {
    Contents strings.Builder      // underlying metodlar + own metodlar
}
mb.Contents.WriteString("Hello, ")
```
- `type X Y` — Y-nin metodları MİRAS OLUNMUR!

## 11. Always Valid Struct (Ch10)
```go
package creditcard

type card struct { number string }    // UNEXPORTED tip!

func New(number string) (card, error) {
    if number == "" {
        return card{}, errors.New("number must not be empty")
    }
    return card{number}, nil
}
// xaricdən invalid card YARATMAQ MÜMKÜNSÜZDÜR!
```

## 12. Control Flow (Ch11-12)
```go
// HAPPY PATH — sola hizala:
if x <= 0 { return false }     // mümkünsüzlər əvvəl, early return
if x%2 != 0 { return false }
return true                    // happy path düz xətt

if ok { ... }                  // if ok == true YAZMA
if _, ok := menu["eggs"]; ok { ... }   // compound if (qısa olduqda!)

switch x {
case 1, 2, 3:        // çoxdəyərli case
    ...
default:             // HƏMİŞƏ yaz
    ...
}

for i, e := range employees { ... }    // index + element
for _, e := range employees { ... }    // element
for x := 0; x < 10; x++ { ... }        // 3-hissəli (nadir)
for { ... }                            // forever (server pattern)

if !e.IsCurrent { continue }    // element ötür — happy path QORU
if money <= 0 { break }         // loop-dan tam çıx
```

## 13. Closure + defer (Ch13)
```go
sort.Slice(nums, func(i, j int) bool {
    return nums[i] < nums[j]     // closure over nums!
})

// LOOP TƏLƏSİ: closure ÇAĞIRIŞ vaxtını görür:
funcs := []func(){}
for _, v := range []int{1, 2, 3} {
    funcs = append(funcs, func() { fmt.Println(v) })
}
// (müasir Go 1.22+: hər iterasiyada yeni v → 1,2,3; köhnə: 3,3,3)

defer f.Close()                  // çıxanda HƏR HANSİ yolla bağla
defer cleanup2()                 // LIFO: əvvəl bu
defer cleanup1()                 // sonra bu

func WriteData(...) (err error) {    // named result
    defer func() {
        if closeErr := f.Close(); closeErr != nil {
            err = closeErr           // YALNIZ fail halda overwrite!
        }
    }()
    return nil
}
// naked return YAZMA — həmişə explicit!
```

## 14. Variadic (Ch13)
```go
func AddMany(inputs ...float64) float64 {
    total := 0.0
    for _, input := range inputs {    // daxildə SLİCE kimi
        total += input
    }
    return total
}
```

## 15. Binary + Cross-Compile (Ch14)
```bash
go build                                # ./hello yaranır (~2 MiB)
GOOS=windows go build                  # hello.exe (PE32+)
GOOS=linux GOARCH=arm go build         # ARM ELF
go tool dist list                      # bütün hədəflər
```
```go
os.Exit(1)        // 0 = OK, non-zero = xəta; YALNIZ main-də!
// init İSTİFADƏ ETMƏ: var x = initialiseThing() əvəzi
```

## 16. Tao Xülasəsi (Ch15)
| Prinsip | Tətbiq |
|---------|--------|
| Kindness | İstifadəçi/İşlədən/Oxuyan/Özün üçün kod |
| Simplicity | Bir işi yaxşı; defaults; lazımsız extensibility YOX |
| Humility | óbvio > clever; standart kitabxana; pre-engineering YOX |
| Not Striving | Problemi LƏĞV et; "ən yaxşı optimallaşdırma = etməmək" |
