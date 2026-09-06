# From Ruby to Golang — Cheatsheet (Azərbaycanca)

> **Kitab:** From Ruby to Golang — Joel Bryan Juliano, Leanpub 2021 · Ruby→Go keçidi üçün sürətli istinad.

---

## Əsas dönüşüm xəritəsi (Ruby → Go)

| Ruby | Go |
|------|-----|
| `@instance_var` | struct sahəsi |
| `attr_accessor` | pointer receiver (setter+getter) |
| class + method | struct + `func (p *T) Method()` |
| inheritance (`<`) | embedding (`*Parent` daxildə) |
| mixin (`include`) | interface (implicit) |
| module | package |
| `require` | `import` |
| Hash `{}` | `map[K]V` |
| Array `[]` | slice `[]T` (fixed `[N]T` nadirən) |
| `.each {}` | `for i, v := range` |
| `*args` (splat) | `...T` variadic |
| `**kwargs` (double-splat) | `...interface{}` + type assertion |
| `is_a?` / typeof | `fmt.Sprintf("%T", v)` |
| `all?/any?/map/detect` | manual loop patternlər |
| RubyGems + Bundler + Gemfile | go get + Go Modules + go.mod |
| `h.delete(k)` | `delete(m, k)` |
| `h[k]` (default) | `m[k]` → zero value |
| `fetch` + default | `v, ok := m[k]` (comma-ok) |
| duck typing | `interface{}` / anonymous struct |

---

## Struct

```go
type dog struct {                    // PRIVATE (kiçik hərf — paket daxili)
    name  string
    breed string
    age   int
}

type Dog struct { Breed string }      // PUBLIC (böyük hərf — export)

pet := dog{name: "Maximus", breed: "Rottweiler", age: 5}
fmt.Printf("%+v", pet)               // {name:Maximus breed:Rottweiler age:5}
```

**Metod (receiver):**

```go
func (p *basket) add_item(entry produce) {   // POINTER — attr_accessor kimi
    *p = append(*p, entry)                   // orijinalı mutasiya
}
func (p basket) items() { ... }               // VALUE — kopya, thread-safe
```

**Embedding (irs):**

```go
type Dog struct {
    *Animal        // pointer embed — sahələr birbaşa
    *Owner         // ÇOXLU irs OK!
    Name  string
}
dog.Animal = &Animal{Kind: "Dog", Diet: "Omnivorous"}
fmt.Println(dog.Animal.Diet)
```

**Anonymous struct:**

```go
Animal := struct {
    Kind string
    Diet string
}{"Dog", "Omnivorous"}
// sahə adı YOX: animals.string, animals.int (hər tip 1 dəfə)
```

---

## Map

```go
// Declaration + make (hint size — performans):
var city map[string]string
city = make(map[string]string)
city["Netherlands"] = "Amsterdam"

// Literal (bir sətir):
var car = map[string]string{"Tesla": "Model 3"}

// Assignment:
basket := map[string]produce{
    "apple": produce{flavour: "...", kind: "fruit"},
}

// Array dəyərlər:
basket := map[string][]string{"fruits": {"apple", "mango"}}

// Dinamik açar (duck-typing — qiyməti var!):
variable := map[interface{}]string{}
variable[1] = "int key"; variable["a"] = "string key"

// Silmə — mövcud olmasa belə təhlükəsiz:
delete(basket, "kale")

// Mövcud olmayan açar — ZERO VALUE (0/""/false/nil):
basket["cabbage"]++                    // təhlükəsiz artırma

// Comma-ok — mövcudluq ayırıcı:
v, ok := basket["cabbage"]
if !ok { /* yoxdur */ }

// Pointer dəyər — nil qaytarır:
var basket map[string]*produce
if exists := basket["cabbage"] != nil; !exists { ... }
```

**Zero value cədvəli:** bool→false · int→0 · float→0.0 · string→"" · func/interface/slice/pointer/channel/map→**nil**

---

## Array & Slice

```go
// FIXED — ölçü TİPİN hissəsi, VALUE (kopya assign):
var array [2]string
fruits := [...]string{"Apple", "Mango", "Orange", "Banana"}   // avto [4]
// [4]string ≠ [6]string — assign XƏTA!

// SLICE — referans, dinamik:
s1 := make([]string, 2)
var s2 []string
s3 := []string{}
s4 := []string{"foo", "bar"}

// make(type, len, CAP) — capacity slicing-lə açılır:
array := make([]string, 2, 3)
array[0] = "foo"; array[1] = "bar"
// array[2] = "baz"        → XƏTA (len 2)
array2 := array[:3]         // capacity aktiv
array2[2] = "baz"           // OK
// array2[3] = ...          → XƏTA (cap 3)

// Referans paylaşımı — hamısı eyni yaddaş:
b2 := b1[:2]; b2[1] = "mango"      // b1 də dəyişir!

// DEEP COPY — müstəqil:
b3 := make([]string, cap(b1))
copy(b3, b2)
b3[1] = "pineapple"    // b2 toxunulmaz

// append — cap aşımında böyüdür:
s = append(s, "baz")
s = append(s, otherSlice...)    // ... açılışı ŞƏRT

// Qarışıq tiplər:
arr := []interface{}{"apple", 1, uuid.New()}
```

---

## For range — 5 forma

```go
for i := 0; i < len(fruits); i++ { }       // C-style
for i, fruit := range fruits { }            // index + value
for _, fruit := range fruits { }            // muted index
for i := range fruits { fruits[i] }         // index-only
for i := range &fruits { }                  // pointer — BAD PRACTİCE
```

---

## Variadic (splat/double-splat)

```go
func basket(fruits ...string) {              // Ruby *args
    for _, fruit := range fruits { }
}
basket("Apple", "Mango")
basket(slice...)                             // array açılışı

func f(args ...interface{}) {                 // Ruby **kwargs emulyasiyası
    item := args[0].(map[string]string)      // type assertion
    kind := args[1].(string)
}

// Tip yoxlama (is_a? / typeof):
if fmt.Sprintf("%T", item) == "map[string]string" { ... }
```

---

## Enumerable → Go patternlər

```go
// all? — hamısı:
all := true
for _, f := range fruits {
    if len(f) <= n { all = false; break }
}

// any? — hər hansı:
any := false
for _, f := range fruits {
    if len(f) > n { any = true; break }
}

// collect/map:
out := make([]int, len(in))
for i, v := range in { out[i] = v * v }

// detect/find — ilk uyğun:
for _, v := range in {
    if cond(v) { found = v; break }
}

// drop(n) — slicing bir sətirdə:
rest := a[n:]

// cycle(n):
for x := 0; x < n; x++ {
    for _, v := range a { }
}
```

---

## Interface

```go
type ProduceBasket interface {
    AddItem(entry Produce)
    RemoveItem(entry Produce)
    Items()
}

type Basket []Produce
type Produce string

func (p *Basket) AddItem(e Produce) { *p = append(*p, e) }
func (p *Basket) RemoveItem(e Produce) {
    for i, v := range *p {
        if v == e { *p = append((*p)[:i], (*p)[i+1:]...); break }
    }
}

var basket ProduceBasket = &Basket{}   // interfeysə yüklə

// Heterojen kolleksiya — tip kontraktı:
var items []ProduceBasket
items = append(items, fruits, veggies)   // hər iki Basket uyuşur
for _, v := range items { v.Items() }

// Return contract:
type Produce interface {
    Flavour() string
    Kind() string
}
var basket []Produce
basket = append(basket, Item{Name: "Apple"})
```

---

## Package Management

```bash
# RubyGems ≈ go get:
go get -u github.com/foo/coffee

# Go Modules (1.11+):  GO111MODULE=on|off|auto
go mod init github.com/foo/menu        # go.mod yarat
# go.mod:
#   module github.com/foo/menu
#   require github.com/foo/coffee v0.0.1
#   require github.com/quux/foo v0.0.1 // indirect

go mod -sync        # import-lardan generasiya
go run .            # asılılıqları avtomatik çək

# legacy dep:
dep init            # Gopkg.toml + Gopkg.lock + vendor/
dep ensure
```

**TOML:**

```toml
[[constraint]]
name = "github.com/foo/bar"
version = "0.0.1"          # və ya branch = "baz"
source = "github.com/quux/baz"   # fork yönləndirməsi
```
