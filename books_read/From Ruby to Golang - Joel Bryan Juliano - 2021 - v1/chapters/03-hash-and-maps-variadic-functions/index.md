# Chapter 3 — Hash və Maps + Variadic Funksiyalar (book səh. 29-54)

## Bu chapter nədən bəhs edir?

Ruby Hash-in Go map-inə tərcüməsi: declaration/assignment ilkinləşdirmə üsulları, make + hint, literal type assignment, struct dəyərli maplər, array dəyərlər (`map[K][]V`), `interface{}` dinamik açarlar, `delete`, mövcud olmayan açarın oxunması (zero value / comma-ok) və Ruby double-splat (`**kwargs`) → Go variadic `...interface{}` emulyasiyası.

---

## Əsas fikirlər

### 1. Hash → Map

Map = açarla dəyərə çıxışın qeyd-qoruyucusu — **adla (indekslə YOX) çıxış**. Key-value DB-lərin (Redis, DynamoDB, Riak, Kyoto Cabinet) konseptual qohumu.

**Ruby:**

```ruby
basket = {}
basket[:fruits] = %w[apple mango avocado]
basket[:veggies] = %w[carrot cucumber kale]
basket.each { |key, val| puts "Key #{key} -- Value #{val.join(' ')}" }
```

### 2. İki ilkinləşdirmə ailəsi

**A. Declaration (sonradan allocate):**

```go
// 1. make ilə (hint size ilə — performans üçün; avtomatik da idarə olunur):
var city map[string]string
city = make(map[string]string)
city["Netherlands"] = "Amsterdam"

// 2. Literal type assignment (bir sətirdə):
var car = map[string]string{"Tesla": "Model 3"}

// 3. Boş literal — sonra doldur:
car := map[string]string{}
car["Tesla"] = "Model 3"
```

**B. Assignment (bir ifadədə yarat+ilkinləşdir):**

```go
variable := map[string]string{keyName: "Value"}   // key/value ilə
variable := map[string]string{}                    // boş — sonra
```

### 3. Struct dəyərli maplər

```go
type produce struct {
    flavour string
    kind    string
}

// Declaration üsulu:
var basket map[string]produce
basket = make(map[string]produce)
basket["apple"] = produce{flavour: `...`, kind: "fruit"}

// Assignment üsulu (bir sətirdə):
basket := map[string]produce{
    "apple": produce{flavour: `...`, kind: "fruit"},
    "kale":  produce{flavour: `...`, kind: "veggies"},
}
for key, value := range basket { /* çap */ }
```

### 4. Array dəyərlər — `map[K][]V`

Bir açara ÇOX dəyər: valueType-in qarşısında `[]`:

```go
basket := map[string][]string{
    "fruits":  []string{"apple", "mango", "avocado"},
    "veggies": []string{"carrot", "cucumber", "kale"},
}
```

`[]valueType` forması map-lərlə yanaşı adi dəyişən və funksiya arqumentlərində də array təyin edir.

### 5. Dinamik tiplər — `interface{}` açarlar

`interface{}` = **Ruby duck-typing davranışı** — int, string, UUID qarışıq açarlar:

```go
variable := map[interface{}]string{}
variable[1] = "from an integer key"
variable["a"] = "from string key"
uuid := uuid.New()
variable[uuid] = "from a UUID key"
```

**Xəbərdarlıq:** Bu rahatlıq PUL TƏLƏB EDİR (tip itkisi) — ucuz/sürətli alternativ: **anonymous struct** (əvvəlki fəsl).

### 6. delete

```ruby
# Ruby: basket.delete(:veggies) — kopyada: basket.dup.tap { ... }
```

```go
delete(basket, "kale")   // mövcud olmasa belə təhlükəsiz
```

### 7. Mövcud olmayan açarın oxunması — zero value

Yalnız FALSE xəta YOX — **tipin default dəyəri qayıdır:**

| Tip | Zero value |
|-----|-----------|
| bool | false |
| int/float | 0 / 0.0 |
| string | "" |
| func, interface, slice, pointer, channel, map | **nil** |

```go
// 1. Zero value müqayisəsi:
if exists := basket["cabbage"] != 0; exists == false { ... }

// 2. Comma-ok (dəqiq üsul — Ruby's fetch-with-default analogiya mənada):
produce, exists := basket["cabbage"]
if !exists { fmt.Println("Item does not exists") }

// 3. Zero-dan istifadə — tip yoxlamadan təhlükəsiz artırma:
basket["cabbage"]++
```

**Pointer dəyərli map:** `map[string]*produce` — mövcud olmayan açar **nil** qaytarır:

```go
var basket map[string]*produce
basket = make(map[string]*produce)
basket["apple"] = &produce{...}
if exists := basket["cabbage"] != nil; exists == false { ... }
```

### 8. Ruby double-splat (`**item`) → Go variadic + interface{}

**Ruby:** `def add_item(kind, **item)` — keyword arqumentləri hash-ə yığır.

**Go emulyasiyası — variadic `...interface{}`:**

```go
func basket(args ...interface{}) {
    item := args[0].(map[string]string)   // TYPE ASSERTION — interface{}-dən geri qaytar
    kind := args[1]
    fmt.Printf("Name: %s\n...", item["name"], item["flavour"], kind)
}

basket(map[string]string{"name": "apple", ...}, "fruit")
```

**Pointer receiver metod variantı:**

```go
func (items *basket) add_item(args ...interface{}) {
    item := args[0].(map[string]string)
    item["kind"] = args[1].(string)
    *items = append(*items, item)
}
```

**Array ötürmə:** variadic funksiyaya mövcud slice ötürərkən `...` açılışı: `basket(produce...)`.

**Tip yoxlaması (Ruby `is_a?` / JS `typeof` analoqu):**

```go
for _, item := range args {
    item_type := fmt.Sprintf("%T", item)          // tipin adını al
    if item_type == "map[string]string" {
        produce := item.(map[string]string)      // assert et
    } else {
        fmt.Println("non-fruit argument", item.(string))
    }
}
```

**...variable vs variable...:** `...interface{}` = funksiyanın variadic PARAMETRİ; `variable...` = çağırışda array AÇILIŞI.

---

## Əsas terminlər

| Termin | İzah |
|--------|------|
| map[keyType]valueType | Go map-in tip forması; Ruby Hash analoqu |
| Declaration vs Assignment | var + make/literal / bir ifadəli `:=` yaratma |
| make hint size | Map-in əvvəlcədən ayrılması — performans; avtomatik idarə olunur |
| map[K][]V | Array dəyərli map — bir açara çox dəyər |
| interface{} açar | Duck-typing emulyasiyası — qarışıq tip açarlar (pulu var!) |
| delete(m, k) | Açar silmə — mövcud olmasa təhlükəsiz |
| Zero value | Mövcud olmayan açarın tip defaultu: 0/""/false/nil |
| Comma-ok | `v, ok = m[k]` — mövcudluğu ayıran forma |
| `map[string]*T` | Pointer dəyər — mövcud olmayan açar nil-dir |
| Variadic `...T` | Ruby splat (*args) analoqu |
| Double-splat emulyasiyası | `...interface{}` + type assertion |
| `%T` | Tipin adı — `is_a?`/`typeof` ekvivalenti yoxlaması |
| Type assertion | `args[0].(map[string]string)` — interface{}-dən konkret tip |

---

## Praktik nəticə

1. **Hash → map seçimləri:** `var m map[K]V` + `make` (hint: böyük map-lərdə performans) və ya `m := map[K]V{...}` (qısa, bir sətirdə).
2. **Qarışıq data → anonymous struct, `interface{}` yalnız son çarə:** Tip itkisi runtime xətası deməkdir.
3. **Mövcudluğucomma-ok ilə yoxla:** yalnız zero value müqayisəsi "0 mövcuddur" halını itirir; `v, ok := m[k]`.
4. **`m[k]++` təhlükəsizdir:** mövcud olmayan açar zero ilə başlayır — sayaclar üçün ideal.
5. **Struct dəyərli map-lərdə pointer seç:** `map[string]*produce` — nil yoxlaması daha ucuz, böyük struct-larda kopya yoxdur.
6. **Ruby kwargs → Go:** `...interface{}` + `args[0].(map[string]string)` — amma real layihədə tipə xas struct arqumenti daha idiomatikdir.
7. **`%T` ilə dinamik tip yoxla:** variadic qarışıq arqumentlərin filtri.

---

## Mənbə

- Kitab: *From Ruby to Golang* — Joel Bryan Juliano, Leanpub, 2021
- Chapter 3: Hash and Maps + Variadic Functions, book səh. 29-54
- PDF səhifələri: 35-60
