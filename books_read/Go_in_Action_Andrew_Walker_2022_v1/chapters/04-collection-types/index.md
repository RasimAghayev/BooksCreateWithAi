# Chapter 4 — Collection Types

## Bu chapter nədən bəhs edir?

Bu chapter, Go-nun kolleksiya tiplərini əhatə edir: arrays (sabit ölçülü massivlər), slices (dinamik massivlər) və maps (açar-dəyər xüritələri). Bu tiplər data strukturlaşdırmanın əsasını təşkil edir və proqramlamada ən çox işlənilən strukturlardır.

## Əsas fikirlər

### 1. Arrays (Massivlər)
**Nədir:** Sabit ölçülü, eyni tip elementlərdən ibarət kolleksiya. Ölçü tipin tərkib hissəsidir — `[5]int` və `[4]int` fərqli tiplərdir.

**Necə işləyir:** `var intArray [5]int` — 5 elementli int massiv, sıfır dəyərlərlə. `[...]int{1,2,3}` — avtomatik ölçü. `[5]int{1: 20, 3: 40}` — sparse (səliqəli) elan. Dəyişməz ölçü — sonradan böyütmək mümkün deyil.

**Nəyə lazımdır:** Sabit sayda element lazım olduqda, CPU cache (tez yaddaş) səmərəliliyi.

**Üstünlükləri:**
- Sequential memory (ardıcıl yaddaş) — cache-friendly
- Dəyər (value) tipi — kopyalama mümkün

**Çatışmamazlıqları:**
- Ölçü sabitdir — böyütmək üçün yeni massiv yaradılmalı
- Funksiyaya ötürülərkən kopyalanır — böyük massivlər üçün pointer istifadə edin

**Kitabdan kod nümunəsi:**
```go
var intArray [5]int
intArray[1] = 20
fmt.Println(len(intArray)) // 5
```
**Mənbə:** Chapter 4, pages 104-178

### 2. Slices (Kəsiklər)
**Nədir:** Dinamik ölçülü, massiv əsasında işləyən kolleksiya tipi. Array-dan daha flexible (çevik) və praktik.

**Necə işləyir:** `[]int{1,2,3}` — literal ilə yaradılır. `make([]int, 0, 10)` — `make` ilə. Slice 3 hissədən ibarətdir: pointer (göstərici), length (uzunluq), capacity ( tutum). `append()` ilə böyütmək mümkündür.

**Nəyə lazımdır:** Dinamik məlumatlar, siyahılar, API response-ları.

**Üstünlükləri:**
- Dinamik ölçü — `append` ilə genişlənir
- Reference type (referans tipi) — kopya xərcləri az

**Çatışmamazlıqları:**
- `append` yeni slice qaytarır — `slice = append(slice, x)` yazmaq lazımdır
- Nil slice vs empty slice fərqi var

**Kitabdan kod nümunəsi:**
```go
slice := []int{10, 20, 30}
slice = append(slice, 40, 50)
fmt.Println(slice) // [10 20 30 40 50]

slice2 := []int{60, 70}
slice = append(slice, slice2...)
fmt.Println(slice) // [10 20 30 40 50 60 70]
```
**Mənbə:** Chapter 4, pages 104-178

### 3. Slice Expressions (Kəsmi İfadələr)
**Nədir:** Massiv və ya slice-dən alt kolleksiya yaratmaq üçün sintaksis — `array[low:high]`, `array[low:]`, `array[:high]`, `array[:]`.

**Necə işləyir:** `slice[1:3]` — 1-ci indekstdən 3-cü indeksin qabağına qədər (1 və 2). Capacity orijinal massivin capacity-dan başlanır. 3-cü indeks daxil deyil.

**Nəyə lazımdır:** Subset (alt kolleksiya) almaq, view (baxış) yaratmaq.

**Üstünlükləri:**
- Yeni kopya yox, orijinal data üçün görünüş
- Performanslı — yaddaş ayırma yox

**Çatışmamazlıqları:**
- Alt slice orijinalı dəyişdirsə, view də dəyişir

**Mənbə:** Chapter 4, pages 104-178

### 4. Maps (Xüritələr)
**Nədir:** Açar (key) və dəyər (value) cütlüklərini saxlayan kolleksiya tipi — hash table (hash cədvəli) əsasında.

**Necə işləyir:** `map[string]int{"a": 1, "b": 2}` — literal ilə. `make(map[string]int)` — `make` ilə. `map[key]` ilə dəyərə müraciət, `map[key] = value` ilə təyinat, `delete(map, key)` ilə silmə. Nil map oxumaq/ yazmaq mümkün deyil — `make` və ya literal istifadə edin.

**Nəyə lazımdır:** Key-value məlumatlar, dictionary lookup (lüğət axtarışı), grouping (qruplaşdırma).

**Üstünlükləri:**
- O(1) average lookup (ortalama axtarış)
- Hər hansı tip açar və dəyər ola bilər

**Çatışmamazlıqları:**
- Orderless (sırasız) — iterasiya sırası təxminidir
- Nil map panic verir — `make` ilə initialize edin

**Kitabdan kod nümunəsi:**
```go
groupNouns := map[string]string{
    "eagle": "convocation",
    "cat":   "clowder",
}

for animal, group := range groupNouns {
    fmt.Printf("A group of %ss is called a %s.\n", animal, group)
}
```
**Mənbə:** Chapter 4, pages 104-178

### 5. Nil Slices və Empty Slices
**Nədir:** Nil slice (`var slice []int`) — initialize edilməmiş, `nil`-dir. Empty slice (`slice := []int{}`) — initialize edilmiş, uzunluğu 0, `nil` deyil.

**Necə işləyir:** `len(slice) == 0` hər ikində true qaytarır. `slice == nil` yalnız nil slice üçün true. `append` hər ikisinə işləyir.

**Nəyə lazımdır:** Zero value (sıfır dəyəri) ilə başlamaq, optional (istəyə bağlı) dəyər.

**Üstünlükləri:**
- Nil slice `append`-ə qarşı tolerantdır
- Empty slice JSON-da `[]` olaraq serialize olunur

**Çatışmamazlıqları:**
- `if slice == nil` check (yoxlama) bəzi hallarda yanlış ola bilər

**Mənbə:** Chapter 4, pages 104-178

### 6. Slice Append Surprises (append Sürprizləri)
**Nədir:** `append` zamanı baş verən gözlənilməsinə davranışlar — yenidən allocation (ayırma), sharing (bölüşmə), side effects (yan təsirlər).

**Necə işləyir:** `append` capacity dolduqda yeni backing array ayırır. Bu vaxt slice-lər paylaşılan array-i göstərməyə davam edir. Yeni slice köhnəni dəyişdirmirsə, lakin orijinal data dəyişə bilər.

**Nəyə lazımdır:** Slice-lərin daxili işini başa düşmək, unexpected (gözlənilməz) bug-ları qarşısını almaq.

**Üstünlükləri:**
- `append` həmişə yeni slice header qaytarır
- Eyni anda bölüşmə və kopyalama arasında balanced (tarazlılıq) saxlayır

**Çatışmamazlıqları:**
- Paylaşılan array dəyişəndə bütün görünüşlər təsir alır

**Mənbə:** Chapter 4, pages 104-178

### 7. Maps Iterasiya və Filter
**Nədir:** `range` ilə map üzərində iterasiya və filterləmə (süzgəc) nümunələri.

**Necə işləyir:** `for key, value := range myMap` ilə iterasiya. `delete(myMap, key)` ilə silmə. Range sırası təxminidir, hər dəfə fərqli ola bilər.

**Nəyə lazımdır:** Map məlumatlarını emal etmək, filter etmək, çıxış formatına gətirmək.

**Üstünlükləri:**
- `delete` təhlükəsizdir — key yoxdursa heç nə olmur
- Range ilə bütün cütlər təkrar edilir

**Çatışmamazlıqları:**
- Sıra təxminidir — sıralı məlumat üçün slice/massiv istifadə edin
- Concurrent (paralel) oxuma/yazma təhlükəsiz deyil — `sync.Map` və ya mutex istifadə edin

**Kitabdan kod nümunəsi:**
```go
for person, age := range familyAges {
    if age < 18 {
        delete(familyAges, person)
    }
}
```
**Mənbə:** Chapter 4, pages 104-178

## Əsas terminlər
- Array (massiv) — Fixed-size (sabit ölçülü) kolleksiya
- Slice (kəsik) — Dynamic (dinamik) kolleksiya, massiv əsasında
- Map (xüritə) — Key-value (açar-dəyər) kolleksiyası
- Length (uzunluq) — `len()` — Mövcud element sayı
- Capacity (tutum) — `cap()` — Azad edilə biləcək maksimum element sayı
- Append — Slice-ə element əlavə etmək
- Nil slice — Initialize edilməmiş slice
- Empty slice — Initialize edilmiş, 0 uzunluqlu slice
- Range loop (aralıq döngüsü) — Kolleksiya iterasiyası
- Hash table (hash cədvəli) — Map-in daxili implementasiyası

## Praktik nəticə
Arrays, slices və maps Go-da ən çox işlənilən kolleksiya tipləridir. Slice-lər praktiki olaraq arrays əvəzinə işlənir — daha flexible və efficient (səmərəli). Map-lər key-value məlumatları üçün standart seçimdir. Nil slice və empty slice fərqini başa düşmək vacibdir.

## Mənbə
Pages: 104-178
