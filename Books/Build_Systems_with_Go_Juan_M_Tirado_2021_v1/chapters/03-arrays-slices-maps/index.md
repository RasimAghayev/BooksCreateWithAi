# Chapter 3 — Arrays, Slices, and Maps (səh. 52-63)

## Bu fəsil nədən bəhs edir?

Go-nun üç əsas data strukturu: array-lər (sabit ölçülü indeksli ardıcıllıq),
slice-lər (array-ə dinamik görünüş) və map-lər (key-value cütləri).
Deklarasiya, len/cap, append/make, range ilə iterasiya və map əməliyyatları.

## Əsas fikirlər

### 1. Array-lər
**Nədir:** Verilmiş uzunluğu olan indeksli element ardıcıllığı. Tipə bağlıdır,
ölçüsü sabitdir, initialləşdirilməyəndə sıfırlarla dolur.

**Kitabdan kod nümunəsi:**
```go
var a [5]int              // [0 0 0 0 0] — sıfırlarla dolur
b := [5]int{0,1,2,3,4}    // bir sətirdə dəyər mənimsətmə
c := [5]int{0,1,2}        // [0 1 2 0 0] — qalanı zero value
fmt.Println(len(a))       // 5 — len funksiyası uzunluğu qaytarır
```

### 2. Slice-lər
**Nədir:** "Alt layan array-in fasizəsiz seqmentinin deskriptoru" — array-ə
istiqamətlənmiş görünüş (view). Slice özü data saxlamır; elementlərə
çıxış verir. Go-da əksər vaxt array yox, slice ilə işlənir.

**Slicing seçimləri (Table 3.1):**
| İndeks | Seçilən elementlər |
|---|---|
| `a[0]` | 0-cı pozisiyadakı element |
| `a[3:5]` | 3-dən 4-ə qədər (5 daxil deyil) |
| `a[3:]` | 3-dən sona qədər |
| `a[:3]` | əvvəldən 2-yə qədər |
| `a[:]` | bütün elementlər |

**Tip fərqi (reflect.TypeOf):**
```go
a := [5]string{"a","b","c","d","e"}
reflect.TypeOf(a)        // [5]string — array (sabit ölçü)
reflect.TypeOf(a[0:3])   // []string  — slice (ölçüsüz)
reflect.TypeOf(a[0])     // string   — tək element
```

### 3. Length və Capacity (len/cap)
**Nədir:** Slice üçün uzunluq (dolu hissə) və tutum (ayrılmış yaddaş)
fərqli anlayışlardır. Array-də hər ikisi eynidir.

**Kitabdan kod nümunəsi:**
```go
a := []int{0,1,2,3,4}
fmt.Println(a, len(a), cap(a))   // [0 1 2 3 4] 5 5

b := append(a, 5)
fmt.Println(b, len(b), cap(b))   // [0 1 2 3 4 5] 6 10
b = append(b, 6)
fmt.Println(b, len(b), cap(b))   // [0 1 2 3 4 5 6] 7 10

c := b[1:4]
fmt.Println(c, len(c), cap(c))   // [1 2 3] 3 9

d := make([]int, 5, 10)
fmt.Println(d, len(d), cap(d))   // [0 0 0 0 0] 5 10
// d[6]=5  → runtime xəta: length-dən böyük index yazmaq olmaz
```

**Sub-kod izahı:**
- `append(slice, dəyər)` → slice-a element əlavə edir; lazım olsa tutumu
  avtomatik artırır (5→10 kimi)
- `make([]T, len, cap)` → müəyyən uzunluq/tutumlu slice yaradır; cap
  verilməzsə len ilə eyni olur
- `b[1:4]` → yeni slice-in capacity-si orijinalın cap-ından başlanğıc
  index çıxılması ilə hesablanır (10-1=9)
- **Vacib:** length ≥ olan indexə yazmaq capacity-dən asılı olmayaraq
  runtime xətadır

### 4. Slice iterasiyası (range)
**Kitabdan kod nümunəsi:**
```go
names := []string{"Jeremy", "John", "Joseph"}

for i := 0; i < len(names); i++ {   // klassik yanaşma
    fmt.Println(i, names[i])
}

for position, name := range names {  // range ilə — daha qısa
    fmt.Println(position, name)
}
```

**Vacib tələ:** `range` elementin **kopyasını** qaytarır — dəyişəni loop
daxilində dəyişmək orijinal slice-i dəyişməz:
```go
for _, name := range names {
    name = name + "_changed"        // kopya dəyişir — təsirsiz
}
for position, name := range names {
    names[position] = name + "_changed"  // düzgün yanaşma — index ilə
}
```

### 5. Map-lər
**Nədir:** `map[K]V` — açarı dəyərlə əlaqələndirən strukturdur. Açarlar
unikaldır; `==`/`!=` operatorlarını dəstəkləyən hər tip açar ola bilər.
Initialləşdirilməmiş map `nil`-dir — ona yazmaq runtime xəta verir.

**Yaratma üsulları:**
```go
var ages map[string]int
fmt.Println(ages)              // map[] (nil)
// ages["Jesus"] = 33          → XƏTA: init olunmayıb

ages = make(map[string]int, 5) // make ilə init (başlanğıc ölçü optional)
ages["Jesus"] = 33             // indi işləyir

ages = map[string]int{         // literal ilə init
    "Jesus": 33,
    "Mathusalem": 969,
}
```

**Oxuma, silmə, "comma ok" idiomu:**
```go
birthdays := map[string]string{
    "Jesus": "12-25-0000",
    "Budha": "563 BEC",
}
fmt.Println(birthdays, len(birthdays))  // map[...] 2

xmas, found := birthdays["Jesus"]   // (dəyər, tapıldımı) — comma ok
fmt.Println(xmas, found)            // 12-25-0000 true

delete(birthdays, "Jesus")          // delete built-in silir
fmt.Println(birthdays, len(birthdays))

_, found = birthdays["Jesus"]
fmt.Println("Did we find when its Xmas?", found)  // false

birthdays["Jesus"] = "12-25-0000"   // yenidən əlavə/overwrite
```

**Sub-kod izahı:**
- `v, ok := m[k]` → tapılmasa `ok=false`, dəyər zero value olur
- `delete(m, k)` → açarı silir
- `m[k] = v` → yeni açar əlavə edir və ya mövcudu yeniləyir

### 6. Map iterasiyası
```go
sales := map[string]int{
    "Jan": 34345, "Feb": 11823, "Mar": 8838, "Apr": 33,
}
fmt.Println("Month\tSales")
for month, sale := range sales {
    fmt.Printf("%s\t\t%d\n", month, sale)
}
```
**Vacib qeyd:** `range` map üzərində **ardıcıllıq zəmanəti vermir** —
 ardıcıllıqlar arası sıra dəyişə bilər. Sıralı çıxış lazımdırsa açarları
 ayrıca sort etmək lazımdır.

## Əsas terminlər
- Array (massiv) — sabit ölçülü indeksli strukturdur
- Slice (dilim) — array-ə dinamik görünüş; data saxlamır
- Length (uzunluq) — dolu element sayı (`len`)
- Capacity (tutum) — ayrılmış yaddaş həcmi (`cap`)
- Map (xəritə) — key-value cütlüyü strukturu
- Comma ok idiomu — `v, ok := m[k]` tapılma yoxlaması
- range (arıalıq operatoru) — kolleksiyalar üzrə iterasiya klauzası

## Praktik nəticə
Gündəlik Go kodunda array yerinə slice istifadə olunur: `make`/`append`
ilə dinamik böyümə. Slice yazarkən `len` sərhədini aşmaq runtime xəta verir —
`cap`-dən asılı deyil. Map işlətməzdən əvvəl mütləq init edin (`make` və ya
literal); açar mövcudluğunu `comma ok` idiomu ilə yoxlayın. `range` slice
elementlərinin kopyasını verir — orijinalı dəyişmək üçün index istifadə edin.

## Mənbə
Pages: 52-63 (PDF 52-63)
