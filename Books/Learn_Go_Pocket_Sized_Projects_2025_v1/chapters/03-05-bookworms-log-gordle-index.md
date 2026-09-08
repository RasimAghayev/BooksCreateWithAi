# Chapters 3-5 — Bookworm's Digest, Log Library, Gordle (səh. 79-202)

## Bu fəsillər nədən bəhs edir?

Üç layihə: (3) JSON faylından bookworm siyahısı oxumaq və ortaq kitabları
tapmaq — loops, map, sort.Interface, bufio; (4) pocketlog — üçsəviyyəli log
kitabxanası — enum, io.Writer, functional options, external/internal test,
doc.go; (5) Gordle — terminal Wordle oyunu — runes, error wrapping, sentinel
errors, strings.Builder, korpus oxunuşu, crypto/rand.

## Əsas fikirlər

### 1. JSON faylın dekodlanması (Ch3)
```go
type Bookworm struct {
    Name  string
    Books []Book
}
type Book struct {
    Author string
    Title  string
}

// loadBookworms reads the file and returns the list of bookworms.
func loadBookworms(filePath string) ([]Bookworm, error) {
    f, err := os.Open(filePath)      // io.Reader qaytarır
    if err != nil {
        return nil, fmt.Errorf("unable to open %q: %w", filePath, err)
    }
    defer f.Close()                   // fayl açıldıqda dərhal defer

    var bookworms []Bookworm
    err = json.NewDecoder(f).Decode(&bookworms)
    if err != nil {
        return nil, fmt.Errorf("unable to decode: %w", err)
    }
    return bookworms, nil
}
```
**Sub-kod izahı:**
- `os.Open` vs `os.Create` vs `os.OpenFile` — oxu / yarat / full nəzarət
- `defer f.Close()` — funksiya bitəndə MÜTLƏQ bağlanır (resource leak qarşısı);
  defer LIFO — son açılan ilk bağlanır
- `json.NewDecoder(reader)` — stream oxuma (Unmarshal-dan üstün: bütün faylı
  yaddaşa yükləmədən)
- JSON sahə adları == struct sahə adları (PascalCase) → avtomatik match

### 2. defer zamanlaması
```go
func main() {
    defer fmt.Println("a bookworm")  // ƏVVƏL "you are", SONRA "a bookworm"
    fmt.Println("you are")
}
```
- defer arqumentləri DƏRHAL hesablanır, icra funksiyanın RETURN-undandır
- Fayl/dəstə bağlanması, mutex unlock — klassik istifadə

### 3. Loop-lar (hamısı `for`)
```go
for i := 0; i < n; i++ { }  // klassik
for cond { }                  // while
for { }                       // sonsuz — daxilində break/return olmalı
for k, v := range m { }       // map/slice/array/channel üzrə
```
- `range` map-də **qeyri-determinist** sıra verir → test üçün sort lazımdır

### 4. map-in set kimi istifadəsi (kitab sayğacı)
```go
func booksCount(bookworms []Bookworm) map[Book]uint {
    count := make(map[Book]uint)  // Book struct-u = açar (müqayisə oluna bilən)
    for _, bookworm := range bookworms {
        for _, book := range bookworm.Books {
            count[book]++
        }
    }
    return count
}

func findCommonBooks(bookworms []Bookworm) []Book {
    booksOnShelves := booksCount(bookworms)
    var commonBooks []Book
    for book, count := range booksOnShelves {
        if count > 1 {
            commonBooks = append(commonBooks, book)
        }
    }
    return commonBooks
}
```
- Açarlı map = counter; `count > 1` filter = "ortaq kitablar"
- Map açarı olmaq üçün struct **müqayisə oluna bilən** olmalıdır (pointer/slice/
  map/channel sahələr OLMAZ)

### 5. sort.Interface (custom sıralama)
```go
type byAuthor []Book

func (b byAuthor) Len() int           { return len(b) }
func (b byAuthor) Less(i, j int) bool { return b[i].Author < b[j].Author }
func (b byAuthor) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }

func sortBooks(books []Book) []Book {
    sort.Sort(byAuthor(books))
    return books
}
```
- 3 metod daxili `sort.Sort` üçün kifayətdir; tip konversiyası `byAuthor(books)`
  — slice-ı "başqa görmək"
- Daha qısa alternativ: `sort.Slice(books, func(i, j int) bool {...})`

### 6. pocketlog kitabxanası (Ch4)
**Enum (iota):**
```go
type Level int

const (
    LevelDebug Level = iota  // 0
    LevelInfo                // 1
    LevelError               // 2
)
```
**Logger + io.Writer asılılığı:**
```go
type Logger struct {
    threshold Level
    output    io.Writer
}

func New(threshold Level, options ...Option) *Logger {
    lgr := &Logger{threshold: threshold, output: os.Stderr} // default
    for _, option := range options {
        option(lgr)
    }
    return lgr
}

type Option func(*Logger)  // functional options pattern

func WithOutput(output io.Writer) Option {
    return func(lgr *Logger) { lgr.output = output }
}

func (l *Logger) Debugf(format string, args ...any) {
    if l.threshold <= LevelDebug {  // əvvəlcə yoxla
        _, _ = fmt.Fprintf(l.output, "DEBUG: "+format+"\n", args...)
    }
}
```
**Sub-kod izahı:**
- `io.Writer` — testdə `testWriter` ilə əvəz olunur (stdout yox, captured string)
- `New(LevelInfo, WithOutput(os.Stdout))` — məcburi threshold + optional output
- Variadic `args ...any` — fmt.Printf kimi istənilən sayda arqument
- `any` = `interface{}` (alias, 1.18+)

**External test paketi:** `package pocketlog_test` — yalnız ixrac olunan API
test olunur (istifadəçi nöqteyi-nəzərindən); internal — hər iki.

**doc.go — paket sənədi:**
```go
/*
Package pocketlog exposes an API to log your work.
First, instantiate a logger with pocketlog.New, and giving it a threshold level.
*/
package pocketlog
```
- `/* */` blok şərhi paket adından ƏVVƏL — godoc-da paket səhifəsi olur

### 7. Gordle oyunu (Ch5)
**Runes vs bytes:**
```go
guess := []rune(string(playerInput))  // baytlar → simvollar
if len(guess) != solutionLength { ... } // RUNE sayı, bayt yox!
```
- `len("həllo")` bayt sayıdır; rune sayı üçün `len([]rune(s))`
- Unicode (ərəbc, yunan) düzgün işləsin deyə — rune əsaslı müqayisə

**Sentinel error (öz tip):**
```go
type corpusError string
func (e corpusError) Error() string { return string(e) }

const ErrCorpusIsEmpty = corpusError("corpus is empty")
// istifadə: errors.Is(err, ErrCorpusIsEmpty)
```

**Error wrapping:**
```go
return nil, fmt.Errorf("unable to open %q for reading: %w", path, err)
// %w = wrap — errors.Is/As zənciri işləyir
```

**Stringer interfeysi (emoji feedback):**
```go
type hint int

const (
    absentCharacter   hint = iota
    wrongPosition
    correctPosition
)

func (h hint) String() string {  // fmt.Stringer
    switch h {
    case absentCharacter:
        return "🩶"  // grey
    case wrongPosition:
        return "🟡"  // yellow
    case correctPosition:
        return "💚"  // green
    }
    return "💔"
}
```

**strings.Builder (string birləşdirmə):**
```go
var sb strings.Builder
for _, h := range fb {
    sb.WriteString(h.String())
}
return sb.String()
```
- `+=` əvəzinə — allokasiyaları minimallaşdırır (benchmark ilə sübut olunur)

**computeFeedback alqoritmi (pseudo-kod → kod):**
```go
func computeFeedback(guess, solution []rune) feedback {
    result := make(feedback, len(guess))
    used := make([]bool, len(solution))
    // 1. hamısını absent işarələ
    // 2. düz mövqeləri tap → correctPosition, used[i] = true
    // 3. qalanları mövqe-sız axtar → wrongPosition
    return result
}
```

**Korpus + random seçim:**
```go
func ReadCorpus(path string) ([]string, error) {
    data, err := os.ReadFile(path)   // []byte — tam fayl
    if err != nil { return nil, fmt.Errorf("unable to open %q: %w", path, err) }
    words := strings.Fields(string(data))  // boşluqlarla ayır
    if len(words) == 0 { return nil, ErrCorpusIsEmpty }
    return words, nil
}
// pickWord: crypto/rand ilə təsadüfi (müəllif: crypto/rand tövsiyə —
// deterministik testlər üçün ayrıca yanaşma)
```

## Əsas terminlər

- io.Reader / io.Writer (oxu/yazı interfeysləri)
- defer (təxirə salınmış icra)
- json.NewDecoder (stream JSON)
- iota (enum sayğacı)
- Functional Options Pattern
- Variadic Function (...args)
- Sentinel Error (məlum xəta dəyəri)
- Error Wrapping (%w)
- Rune vs Byte (simvol vs bayt)
- fmt.Stringer interfeysi
- strings.Builder

## Praktik nəticə

- Fayl aç → dərhal defer Close; JSON stream = Decoder
- Counter üçün `map[T]uint`; unikal üçün `map[T]struct{}`
- Kitabxana dizaynı: minimal New + Option-lar; io.Writer injekt et
- Unicode mətn işləyirsən — []rune ilə düşün
- Xətalar: %w ilə wrap, sentinel + errors.Is

## Mənbə

Pages: 79-202 (Chapters 3-5, Learn Go with Pocket-Sized Projects)
