# Chapter 4 — Arrays və Slices + Navigasiya (book səh. 55-81)

## Bu chapter nədən bəhs edir?

Ruby Array-in Go-da İKİ klasifikasiyası: Fixed Array (dəyər tipi, sabit ölçü) və Sliced Array (referans tipi, dinamik). Big O konteksti, `[...]` avto-ölçü, tip uyğunsuzluğu, value-copy vs referans davranışı, make-in 3-cu parametri (capacity), deep copy, append, variadic ötürmə, `[]interface{}` qarışıq massivlər və for range-in 5 semantik forması.

---

## əsas fikirlər

### 1. Array nə vaxt? — Big O konteksti

Array = **O(n)** — ardıcıl emal üçün optimal (cheap, fast, index-çaluşur). Map = **O(1)** — adla çıxış. Kütləvi/batch emal → array.

Go-da İKİ növ: **Fixed Array** (sabit ölçü, value type) və **Sliced Array** (dəyişən ölçü, reference type) — fərqini bilmək "doğru iş üçün doğru struktur" deməkdir.

### 2. Fixed Array

```go
const MAX_ARRAY_SIZE = 2
var array [MAX_ARRAY_SIZE]string    // ölçü TİPİN HİSSƏSİDİR
array[0] = "foo"
array[1] = "bar"
```

**[...] avtomatik ölçü** — literal sayından hesablanır:

```go
var flavours [4]string
fruits := [...]string{"Apple", "Mango", "Orange", "Banana"}   // [4]string olur
flavours = fruits    // OK — ölçülər üst-üstə düşür
```

**Ölçü tipin hissəsidir** — `[6]string = [4]string` KOMPİLYASİYA XƏTASI (statik tip!). Hədəf ölçüsü avtomatik hesablanMİR — yalnız mənbə literalı.

### 3. Fixed Array = VALUE tipi (copy assignment)

```go
var oldArray [2]string
var newArray [2]string
oldArray[0] = "foo"
newArray = oldArray         // BÜTÜM elementlər KOYALANIR

newArray[1] = "baz"
oldArray[1] = "bar"

fmt.Println(oldArray)  // [foo bar] — MÜSTƏQİL!
fmt.Println(newArray)  // [foo baz]
```

Ruby-da (referans dili) hər iki dəyişən EYNİ array-i göstərərdi. Go-da fixed array assign = **tam kopya** — böyük datas ETİKLİDİR → slice istifadə et.

### 4. Sliced Array — referans tipi

```go
// 4 ilkinləşdirmə:
slicedArray1 := make([]string, 2)
var slicedArray2 []string
slicedArray3 := []string{}
slicedArray4 := []string{"foo", "bar"}
```

**Assign = referans paylaşımı — 4 dəyişən eyni array-i göstərir:**

```go
slicedArray1 := make([]string, 4)
slicedArray1[0] = "foo"
slicedArray2 = slicedArray1       // kopya YOX — istinad
slicedArray2[1] = "bar"
slicedArray3 = slicedArray2
slicedArray3[2] = "baz"
slicedArray4 = slicedArray3
slicedArray4[3] = "quux"
// Hamısı çap olunanda: [foo bar baz quux] — EYNİ yaddaş!
```

Bütün standart kitabxana və public API slice istifadə edir.

### 5. Capacity — make-in 3-cü parametri

```go
make(type, length, capacity)
// array := make([]string, 2) ≡ make([]string, 2, 2)
```

**Capacity = slicing-dən sonra yeni array-in QALAN tutumu.** Tələ: `make([]string, 2, 3)` — len 2 (yaza bilərik: [0], [1]), [2]-ü DİREKT yazmaq XƏTADIR. Amma slicing ilə açılır:

```go
array := make([]string, 2, 3)
array[0] = "foo"; array[1] = "bar"
array2 := array[:3]    // capacity-ni SƏMƏRƏLİNDİR
array2[2] = "baz"      // indi OK
// array2[3] = "..."  → yenidən error (cap aşımı)
```

**Ortaq yaddaş tələsi:**

```go
basket1 := make([]string, 1, 3)
basket1[0] = "apple"
basket2 := basket1[:2]
basket2[1] = "mango"
basket3 := basket2[:3]
basket3[2] = "banana"

basket3[0] = "banana"; basket3[1] = "banana"; basket3[2] = "banana"
fmt.Println(basket1)  // [banana] — basket3-ün dəyişikliyi basket1-ə DƏ DÜŞDÜ!
```

Slicing YENİ array YOX — eyni yaddaş ünvanlarına yeni pəncərə.

### 6. Deep Copy — make + copy

Orijinaldan MÜSTƏQİL nüsxə:

```go
basket1 := make([]string, 1, 3)
basket1[0] = "apple"
basket2 := basket1[:2]
basket2[1] = "mango"

basket3 := make([]string, cap(basket1))    // yeni ayrılış
copy(basket3, basket2)                     // dərin kopya

basket3[1] = "pineapple"
// basket2 = [apple mango] — toxunulmaz; basket3 = [apple pineapple]
```

### 7. Append — capacity-ni AŞMAQ

Slicing capacity həddindədir (`array[:3]` → `[3]`-dən sonra error). **append capacity-ni avtomatik böyüdür:**

```go
array := make([]string, 1, 2)
array[0] = "foo"
array2 := array[:2]
array2[1] = "bar"
array3 := array[:2]
array3 = append(array3, "baz")    // cap aşılanda YENİ array ayrılır
```

**Risk:** append dolu slice-da realloсasiya edir — köhnə dəyərlər yeni referansda (müasir: append dəyəri assign etmək MƏCBURİDİR).

### 8. Variadic — Ruby splat (*args) → Go `...T`

```ruby
def basket(*fruits) ... end        # Ruby catch-all
```

```go
func basket(fruits ...string) {    // Go variadic
    for _, fruit := range fruits {
        fmt.Printf("%s is in the basket\n", fruit)
    }
}
basket("Apple", "Mango", "Orange", "Banana")

// Mövcud slice ötürməsi — ... AÇILIŞI:
fruits := []string{"Apple", "Mango", "Orange", "Banana"}
basket(fruits...)
```

### 9. `[]interface{}` — qarışıq məzmun

Əsas tiplər qarışıq YAZILA BİLMƏZ (`[3]string{"string", 1, uuid}` → error). Duck-typing lazımdırsa:

```go
var array []interface{}
array = []interface{}{"apple", 1, uuid.New()}
fmt.Println(array[0])  // apple
fmt.Println(array[1])  // 1
fmt.Println(array[2])  // <uuid>
```

### 10. Naviqasiya — Ruby each → Go for range (5 semantik forma)

**1. C-style:**

```go
for i := 0; i < len(fruits); i++ {
    fmt.Printf("%d - %s\n", i, fruits[i])
}
```

**2. Value semantic (index+value):**

```go
for i, fruit := range fruits { ... }
```

**3. Muted parameter (`_` ilə index sönüşdürmə):**

```go
for _, fruit := range fruits { ... }
```

**4. Index-only:**

```go
for i := range fruits {
    fmt.Println(i, fruits[i])
}
```

**5. Pointer access (`range &fruits`):** Massivin ünvanı üzərində range — **BAD PRACTİCE** kimi qeyd olunur (mutasiya hallarında təhlükəli).

---

## Əsas terminlər

| Termin | İzah |
|--------|------|
| Fixed Array | Sabit ölçülü massiv — ölçü TİPİN HİSSƏSİ; value type (kopya) |
| Sliced Array | Dinamik — referans tipi; standart API-nin hamısı |
| O(n) vs O(1) | Array ardıcıl emal / map açarla çıxış |
| `[...]` | Literal sayından avto-ölçü (yalnız mənbədə) |
| Value copy assign | Fixed array = tam element kopyası — müstəqil |
| Reference assign | Slice = eyni yaddaş — 4 dəyişən 1 massiv |
| make(type, len, cap) | 3-cü parametr capacity — slicing-dən sonrakı tutum |
| Slicing `s[:3]` | Capacity-ni aktivləşdirən pəncərə — eyni yaddaş |
| Deep copy | `make + copy` — müstəqil nüsxə |
| append | Capacity aşımında realloсasiya edən böyümə |
| `...T` parametr | Ruby splat ekvivalenti; `slice...` çağırışda açılış |
| `[]interface{}` | Qarışıq tipli massiv — duck-typing |
| 5 range forması | C-style / index+value / muted / index-only / pointer |

---

## Praktik nəticə

1. **Adi ehtiyac = slice:** Fixed array yalnız dəqiq bilinən sabit ölçüdə (məs., [3]float64 matris).
2. **Ölçü tipdir:** `[4]string ≠ [6]string` — assign üçün hər ikisi üst-üstə düşməlidir.
3. **Slice assign = paylaşım:** Bir dəyişəndən düzəliş digərinə düşür — müstəqillik üçün `make+copy` (deep copy).
4. **Capacity-ni slicing açır:** `make(s, 2, 3)` + `s[:3]` — amma append-dən başqa cap-aşımı mümkünsüz.
5. **Append nəticəsini mənimsət:** realloсasiya yeni array verə bilər.
6. **Splat → `...`:** `func f(args ...T)` + çağırış `f(slice...)`.
7. **Qarışıq massiv → `[]interface{}`** (və ya anonymous struct — tipli alternativ).
8. **each → `for i, v := range`** (5 forma — 2-ci ən idiomatik).

---

## Mənbə

- Kitab: *From Ruby to Golang* — Joel Bryan Juliano, Leanpub, 2021
- Chapter 4: Arrays and Slices + Navigating your Arrays, book səh. 55-81
- PDF səhifələri: 61-87
