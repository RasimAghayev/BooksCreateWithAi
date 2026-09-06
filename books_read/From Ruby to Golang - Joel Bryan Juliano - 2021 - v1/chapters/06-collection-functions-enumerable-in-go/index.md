# Chapter 6 — Collection Functions: Ruby Enumerable → Go (book səh. 93-100)

## Bu chapter nədən bəhs edir?

Ruby-nin Enumerable metod ailəsinin (all?, any?, collect, cycle, detect, drop, drop_while) Go-da manual emulyasiyası — Ruby üçün builtin olan kolleksiya əməliyyatlarının Go loop-ları ilə necə qurulduğunu öyrədir.

---

## Əsas fikirlər

**Kontekst:** Konfransda "hamının telefonu qapatılıb?" (all?) və "kofe fasiləsi istəyən var?" (any?) sualları kimi — kolleksiyaya ÜMUMİ sorğu. Ruby-də builtin; Go-da ÖZ LOOP-Umuzu yazırıq.

### 1. `all?` — bütün elementlər şərti qarşılayır?

```ruby
["Apple", "Orange"].all? { |fruit| fruit.length > 4 }  # true
["Apple", "Orange"].all? { |fruit| fruit.length > 5 }  # false
```

**Go (kitabın yanaşması — 2 mərhələ):**

```go
// 1. Şərt nəticələrini bool massivinə yığ:
var all = make([]bool, len(fruits))
for i := 0; i < len(fruits); i++ {
    all[i] = len(fruits[i]) > length
}

// 2. Hamısını yoxla, false-ta break:
var cond bool
for v := range all {
    if cond = all[v] == true; cond == false {
        break
    }
}
return cond
```

Tam funksiya: `is_lengthy(length int, fruits []string) bool` — `is_lengthy(4, fruits)` → true; `is_lengthy(5, ...)` → false.

**İdiomatik Göstəriş (müasir praktika):** bool massivi YOX — bir loop-da early-exit:

```go
for _, f := range fruits {
    if len(f) <= length {
        return false
    }
}
return true
```

### 2. `any?` — hər hansı element şərti qarşılayır?

```ruby
fruits.any? { |word| word.length < 4 }  # false
fruits.any? { |word| word.length > 5 }  # true
```

Go: bool massivi doldur → true tapanda **break**.

### 3. `collect` / `map` — nəticələri yeni massivə yığ

```ruby
(1..4).collect { |i| i*i }   # [1, 4, 9, 16]
```

```go
var collect [4]int
integers := [4]int{1, 2, 3, 4}
for i := 0; i < len(integers); i++ {
    collect[i] = integers[i] * integers[i]
}
```

### 4. `cycle` — n dəfə təkrarla

```ruby
a.cycle(2) { |x| puts x }   # a, b, c, a, b, c
```

```go
count := 2
a := [3]string{"a", "b", "c"}
for x := 0; x < count; x++ {
    for i := 0; i < len(a); i++ {
        fmt.Println(a[i])
    }
}
```

### 5. `detect` / `find` — ilk uyğun element

```ruby
(1..100).detect { |i| i % 5 == 0 and i % 7 == 0 }   # 35
```

```go
for i := 1; i < 100; i++ {
    if (i%5 == 0) && (i%7 == 0) {
        fmt.Println(i)   // 35
        break            // İLK uyğun — dayan
    }
}
```

### 6. `drop` — ilk n elementi at

```ruby
a.drop(3)   # [1,2,3,4,5,0] → [4, 5, 0]
```

```go
a := [...]int{1, 2, 3, 4, 5, 0}
drop := a[3:]       // SLICING — bir sətirdə!
fmt.Println(drop)
```

### 7. `drop_while` — şərt true olduqca at

```ruby
a.drop_while { |i| i < 3 }   # → [3, 4, 5, 0]
```

Kitabın Go versiyası: range + şərt + slicing + break (sətir 3-də kəsir).

**Nüans (ehtiyat):** Kitabın `for _, i := range a` + `a[i] < 3` kodu range-dən gələn dəyəri indeks kimi qarışdırır — düzgün filtr `v < 3` dəyər üzərində olmalıdır. Konsept: iterasiya + ilk uyğunda `a[k:]` + `break`.

---

## Əsas terminlər

| Termin | İzah |
|--------|------|
| all? / any? | Predicate metodlar — hamısı / hər hansı biri |
| collect (map) | Transformasiya nəticələrini yeni massivə yığ |
| cycle(n) | Massivi n dəfə təkrarla |
| detect (find) | İlk uyğun element — break ilə |
| drop(n) | İlk n-i at → slicing `a[n:]` |
| drop_while | Şərt true olduqca at |
| Predicate method | Boolean qaytaran metod (Ruby termini) |
| Enumerable | İterasiya oluna bilən kolleksiya + builtin metodlar (Ruby) |

---

## Praktik nəticə

1. **Ruby-nin builtin-ləri Go-da YOXDUR** — loop patternləri ilə öz "koleksiya funksiyalarını" yaz.
2. **all?/any? = bool massivi + hamısı/hər hansı yoxlaması** (müasir: bir loop-da early return).
3. **collect = nəticə massivi + indeks yazması.**
4. **detect = şərt + break.**
5. **drop = slicing bir sətirdə** — ən xəlitə analoq (Go-da natural).
6. **Cəmiyyət həlli:** gobyexample.com/collection-functions — hazır patternlər.
7. Ruby düşüncəsindən qalAN dərs: "massivə ümumi sual ver" abstraksiyası — Go-da funksiyaya çıxarıla bilər (müasir: generics/callback).

---

## Mənbə

- Kitab: *From Ruby to Golang* — Joel Bryan Juliano, Leanpub, 2021
- Chapter 6: Collection Functions in Go, book səh. 93-100
- PDF səhifələri: 99-107
