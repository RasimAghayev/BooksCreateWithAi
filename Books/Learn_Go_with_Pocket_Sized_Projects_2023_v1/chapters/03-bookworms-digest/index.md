# Chapter 3 — A bookworm's digest: Playing with loops and maps (səh. 79-121)

## Bu chapter nədən bəhs edir?

"Bookworms" CLI layihəsi: JSON fayl oxuma, `os.File` + `defer Close()`, struct-lara
decode, `map` ilə sayğac, slice-lar, `sort.Interface` ilə custom sıralama və
`bufio` ilə buffered oxuma.

## Əsas fikirlər

### 1. Layihə və məqsədi
JSON faylında bookworm-ların kitab kolleksiyaları var → ortaq kitabları tap →
stdout-a çap et. Kod 2 fayla bölünür: `main.go` (terminal, çap) + `bookworm.go`
(business logic — test oluna bilən).

### 2. JSON dekodlaşdırma
```go
type Book struct {
    Author string `json:"author"` // struct tag → JSON açarı
    Title  string `json:"title"`
}
type Bookworm struct {
    Name  string `json:"name"`
    Books []Book `json:"books"`
}

func loadBookworms(filePath string) ([]Bookworm, error) {
    f, err := os.Open(filePath) // io.Reader qaytarır
    if err != nil {
        return nil, fmt.Errorf("failed to open file: %w", err) // %w → wrap
    }
    defer f.Close() // funksiya bitəndə resurs azad olunur

    var bookworms []Bookworm
    if err := json.NewDecoder(f).Decode(&bookworms); err != nil {
        return nil, fmt.Errorf("failed to decode JSON: %w", err)
    }
    return bookworms, nil
}
```
**Sub-kod izahı:**
- Struct tag-lər JSON açarlarını field-lərə bağlayır
- `os.Open` → yalnız oxuma; `os.Create` → yaratma/yazma (silir!); `os.OpenFile`
  → qarışıq rejimlər
- `defer f.Close()` → funksiya çıxışında (hətta panic olsa) bağlanır
- `%w` → xəta zənciri qorunur, `errors.Is/As` işləyir

### 3. defer semantikası
```go
func main() {
    fmt.Println("a bookworm") // konsol: "a bookworm" birinci
    defer fmt.Println("you are") // ƏN SON icra olunur
}
```
defer arqumentləri **dərhal** qiymətlənir, icra funksiya sonuna təxirə salınır.
Fayl açma + Close **eyni kod blokunda** görünməlidir.

### 4. Ortaq kitablar — map sayğacı
```go
func booksCount(bookworms []Bookworm) map[Book]uint {
    counts := map[Book]uint{} // Book struct-a görə sayğac
    for _, bookworm := range bookworms {
        for _, book := range bookworm.Books {
            counts[book]++ // zero value = 0 → avtomatik artım
        }
    }
    return counts
}

func findCommonBooks(bookworms []Bookworm) []Book {
    booksOnShelves := booksCount(bookworms)
    var commonBooks []Book
    for book, count := range booksOnShelves {
        if count > 1 {
            commonBooks = append(commonBooks, book)
        }
    }
    return sortBooks(commonBooks)
}
```
**Sub-kod izahı:**
- `map[Book]uint` → struct açar: müqayisə oluna bilən (comparable) tiplər
- Map iteration **təsadüfi sıralıdır** → deterministik çap üçün sort MƏCBURİDİR
- Zero value xüsusiyyəti: `counts[book]++` — açar yoxdursa 0-dan başlayır

### 5. Slice — 3 sahəli struktur
```
slice = pointer (alt array-ə) + len + cap
```
```go
books := make([]Book, 0, 5) // len=0, cap=5 — böyümə üçün ön-ayırma
```
`append` cap dolanda yeni alt array yaradır — performans üçün cap əvvəlcədən ver.

### 6. sort.Interface — custom sıralama
```go
type byAuthor []Book

func (b byAuthor) Len() int           { return len(b) }
func (b byAuthor) Less(i, j int) bool { return b[i].Author < b[j].Author }
func (b byAuthor) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }

sort.Sort(byAuthor(books)) // müəllif → sonra başlıq (Less-də müqayisə)
```
**Sub-kod izahı:**
- 3 metodun hamısı implement olunmalıdır (istifadə olunmasa belə)
- Daha sadə alternativ: `sort.Slice(books, func(i, j int) bool {...})`

### 7. Determinizm və testlər
Xəritə iterasiyası random olduğundan testlərdə gözlənilən nəticə sıralanmış
olmalıdır. `equalBooks` kimi helper-lər `t *testing.T` qəbul edir + `t.Helper()`
ilə test hesabatında düzgün sətir göstərir.

### 8. bufio — buffered oxuma
```go
f, _ := os.Open(filePath)
defer f.Close()
r := bufio.NewReader(f) // bir oxunuşda böyük chunk — sistem çağırışları az
```
Hər baytı bir-bir oxumaq əvəzinə buffer dolana qədər oxu — I/O sürəti artır.

## Əsas terminlər

- Struct Tag (struct etiketi)
- Deferred Function (təxirə salınmış funksiya)
- Zero Value (sıfır dəyər)
- Comparable Type (müqayisəli tip)
- Determinism (determinizm)
- Buffered Reader (buferli oxuyucu)

## Praktik nəticə

- JSON decode: struct tag + `json.NewDecoder(f).Decode(&v)`
- Fayl açılan yerdə defer ilə bağla
- Sayğac üçün `map[T]uint` + zero value artımı
- Map çıxışını sort etmədən istifadə etmə (non-determinizm)
- Böyük fayllar üçün `bufio.NewReader`

## Mənbə

Pages: 79-121 (Chapter 3, Learn Go with Pocket-Sized Projects)
