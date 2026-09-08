# Chapter 3 — A bookworm's digest: Playing with loops and maps (səh. 79-121)

## Bu chapter nədən bəhs edir?

"Bookworms" CLI layihəsi: JSON faylından kitab kolleksiyaları oxunur, ümumi
kitablar tapılır və çap olunur. Mövzular: loops, map, struct-ın map key-i
kimi istifadəsi, JSON decode, defer ilə fayl bağlanması, sort.Interface,
bufio, set üçün `map[X]struct{}`.

## Layihə məqsədi

1. Kitab qurdularının (bookworm) kolleksiyalarını JSON faylından oxu
2. Ortaq kitabları tap (2+ kitabqurdunun rəfində olan)
3. Standart çıxışa çap et
4. (Bonus) Oxu tövsiyələri — "Other Readers Bought" üslubunda

## Əsas fikirlər

### 1. JSON faylından decode — loadBookworms()
**Struktur JSON strukturu ilə üst-üstə düşməlidir** (PascalCase exposure —
sahələr EXPORTED olmalıdır, yoxsa decoder "görmür" və boş qalır):
```go
type Bookworm struct {
    Name  string
    Books []Book
}
type Book struct {
    Author string
    Title  string
}
```

**Fayl açma + decode (klassik pattern):**
```go
func loadBookworms(filePath string) ([]Bookworm, error) {
    f, err := os.Open(filePath)      // read-only file descriptor
    if err != nil {
        return nil, err
    }
    defer f.Close()                   // resurs sızması olmasın deyə

    var bookworms []Bookworm
    err = json.NewDecoder(f).Decode(&bookworms)  // pointer ötürülür!
    return bookworms, err
}
```
**Sub-kod izahı:**
- `os.Open` → `*os.File` (io.Reader) + error; `os.Create` yazma üçün (truncate edir!),
  `os.OpenFile` — hüquqlarla açma
- `defer f.Close()` → funksiya çıxanda resurs azad olunur; **defer-lərin LIFO
  sırası** — son yazılan birinci icra olunur
- `Decode(&bookworms)` → dəyər yox, POİNTER ötürülür — decoder slice-i doldurur
- `json.Decoder` vs `json.Unmarshal`: Decoder `io.Reader` oxuyur (stream),
  Unmarshal tam `[]byte` tələb edir → Decoder üstünlük təşkil edir

### 2. Test datası — testdata/
`testdata` qovluğu Go tooling üçün xüsusiyyətlidir — compile ediləndə skip
olunur. Test JSON-ları və expected dəyərləri orada saxlanılır.

### 3. Loops — hamısı `for`
Go-da bütün loop-lar `for` açar sözü ilə:
```go
for i := 0; i < 10; i++ { }   // klassik
for i < 10 { }                 // while ekvivalenti (şərt + post yoxdur)
for { }                        // sonsuz — içində break/return olmalı
for i, bw := range bookworms { } // range: slice/map/channel üzrə iterasiya
```
- `range` slice-da: `i` = index (0..len-1), ikinci dəyər = **copy** (referans deyil)
- İndeks lazım deyilsə: `for _, bookworm := range ...`

### 4. map — unikal açar kolleksiyası
**Nədir:** Go-da `map` — unordered assosiativ massiv. Unikal açarlar üçün
idiomatik yoldur.

**Açar tələbi:** açar **comparable** olmalıdır — `key1 == key2` yazıla bilirsə.
Struct da açar ola bilər (daxilində pointer/slice/map/channel/func/interface
OLMAZ — bunlar comparable deyil).

**Zero-value axtarışı:**
```go
v, ok := mapped[3]        // ok=false → açar yoxdur, v = tipin zero dəyəri
if v, ok := mapped[3]; ok { /* ... */ }  // scope məhdud qısa forma
```

### 5. Kitab sayğacı — map ilə counter
```go
func booksCount(bookworms []Bookworm) map[Book]uint {
    count := make(map[Book]uint)          // kitab → sayı
    for _, bookworm := range bookworms {
        for _, book := range bookworm.Books {
            count[book]++                 // mövcud olmasa da: zero=0 → 1
        }
    }
    return count
}
```
**Sub-kod izahı:**
- `map[Book]uint` — struct açar! Go compiler `Hash` funksiyası tələb etmir
  (Java-dan fərqli), hashable struct-ı özü hesablayır
- `count[book]++` — mövcud olmayan açar zero-value (0) ilə başlayır → increment
  təbii işləyir

**Equal helper (test üçün):**
```go
func equalBooksCount(t *testing.T, got, want map[Book]uint) bool {
    t.Helper()                         // xəta sətri helper-ə yox, testə aid olsun
    if len(got) != len(want) { return false }
    for book, targetCount := range want {
        count, ok := got[book]
        if !ok || targetCount != count { return false }
    }
    return true
}
```

### 6. Ortaq kitablar — filtrləmə
```go
func findCommonBooks(bookworms []Bookworm) []Book {
    booksOnShelves := booksCount(bookworms)
    var commonBooks []Book                 // nil slice — append ilə böyüyür
    for book, count := range booksOnShelves {
        if count > 1 {
            commonBooks = append(commonBooks, book)
        }
    }
    return commonBooks
}
```

### 7. slice — 3 sahə
Slice = **pointer (underlying array) + len + cap**:
- `len` = indiki element sayı; `cap` = reallocated-ə qədər tutum
- `var s []Book` → nil slice (allocation yoxdur); `append` ilə böyüyür
- `make([]Book, 0, 5)` → len=0, cap=5 (son ölçü məlum olanda — 1 realloc qənaəti)

### 8. Determinizm — sort
Map iteration tərtibi **random** olduğundan test yalnız sort edilmiş nəticə ilə
müqayisə oluna bilər → sortBooks. Sorting ayrı funksiyada (əsas alqoritmdən
ayrı) — test asanlaşır.

### 9. sort.Interface — custom comparator
```go
type byAuthor []Book

func (b byAuthor) Len() int           { return len(b) }
func (b byAuthor) Less(i, j int) bool { return b[i].Author < b[j].Author }
func (b byAuthor) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }

sort.Sort(byAuthor(books))    // author sonra title
```
**Sub-kod izahı:** Orta tip (named type) `[]Book` üzərində interfeysi
implement edir; 3 metodun hamısı vacibdir (istifadə olunmasa belə).

### 10. Oxu tövsiyələri (bonus — recommendOtherBooks)
Hər kitabqurdunun rəfindəki hər kitab üçün → həmin rəfdəki DİGƏR kitabları
recommendation xəritəsinə qeyd et → "kim A oxuyursa, B də xoşuna gələr".

### 11. bufio — buferlənmiş oxuma
Disk oxuması kiçik parçalarla bahadır → buferlə:
```go
buffedReader := bufio.NewReaderSize(f, 1024*1024)   // 1 MiB bufer
decoder := json.NewDecoder(buffedReader)
err = decoder.Decode(...)
```
`bufio.Reader` `io.Reader` implement edir, amma `Closer` YOX — `f.Close()`
əvvəlki descriptor üzərindən çağrılır.

### 12. Set pattern — map[X]struct{}
Dublikatları atmaq üçün (məs. CD kolleksiyasında unikal qruplar):
```go
myMap := make(map[Band]struct{})
myMap[band] = struct{}{}    // zero-size dəyər — yalnız açar maraqlıdır
```

## Əsas terminlər

- Table-driven test (cədvəl-əsaslı test)
- File descriptor (fayl deskriptoru)
- defer (gecikməli icra)
- Zero value (sıfır dəyər)
- Comparable type (müqayisə oluna bilən tip)
- Hashable struct (hash-lana bilən struktur)
- len/cap (uzunluq/tutum)
- sort.Interface (sıralama interfeysi)
- Buffered reader (buferlənmiş oxuyucu)
- Set pattern (çoxluq nümunəsi)

## Praktik nəticə

- JSON decode standart axını: `os.Open` → `defer Close` → `json.NewDecoder(f).Decode(&v)`
- Unikal dəyərlər üçün `map[X]struct{}`; sayğac üçün `map[X]uint`
- Strukturu map açarı etmək olar — sadə comparable sahələrlə
- Map nəticəsi həmişə non-deterministikdir → testdə sort et
- Fayl I/O-da defer Close; böyük fayllarda bufio

## Mənbə

Pages: 79-121 (Chapter 3, Learn Go with Pocket-Sized Projects)
