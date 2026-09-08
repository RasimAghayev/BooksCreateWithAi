# Chapter 5 — Gordle: Play a Word Game in Your Terminal (səh. 154-202)

## Bu fəsil nədən bəhs edir?

Terminal Wordle oyunu: input oxu (bufio.Reader), runes vs baytlar, xəta
wrap (%w) və sentinel error, fmt.Stringer (emoji feedback), strings.Builder,
feedback alqoritmi, korpus faylı, random söz seçimi, multi-attempt loop.

## Əsas fikirlər

### 1. Game strukturu + pointer receiver
```go
type Game struct {
    solution    []rune
    maxAttempts int
    reader      *bufio.Reader
}

func New(reader io.Reader, corpus []string, maxAttempts int) (*Game, error) {
    if len(corpus) == 0 {
        return nil, ErrCorpusIsEmpty
    }
    return &Game{
        reader:      bufio.NewReader(reader),
        solution:    []rune(strings.ToUpper(pickWord(corpus))),
        maxAttempts: maxAttempts,
    }, nil
}
```
- Pointer receiver-lər — state dəyişən bütün metodlarda; bir tipdə pointer +
  value receiver QARIŞDIRMA (konsistensiya qaydası)

### 2. Runes vs baytlar (Unicode!)
```go
guess := []rune(string(playerInput))
if len(guess) != solutionLength { ... }  // RUNE sayı — bayt yox
```
- `len("həllo")` = BAYT sayı; `len([]rune("həllo"))` = simvol sayı
- Ərəb/Yunan dəstəyi üçün bütün müqayisələr rune səviyysində

### 3. Xəta idarəetməsi
```go
// sentinel error — öz tipi ilə:
type corpusError string
func (e corpusError) Error() string { return string(e) }
const ErrCorpusIsEmpty = corpusError("corpus is empty")

// wrap:
return nil, fmt.Errorf("unable to open %q for reading: %w", path, err)
```
- `errors.Is(err, ErrCorpusIsEmpty)` — sentinel yoxlaması
- `%w` — xəta zənciri qorunur

### 4. Feedback — Stringer + emoji
```go
type hint int

const (
    absentCharacter   hint = iota
    wrongPosition
    correctPosition
)

func (h hint) String() string {  // fmt.Stringer — çapda avtomatik
    switch h {
    case absentCharacter:   return "🩶"
    case wrongPosition:     return "🟡"
    case correctPosition:   return "💚"
    }
    return "💔"
}
```

### 5. strings.Builder vs +=
```go
var sb strings.Builder
for _, h := range fb { sb.WriteString(h.String()) }
return sb.String()
```
- Benchmark ilə sübut: += hər addımda yeni string allokasiyası; Builder —
  amortizə olunmuş

### 6. computeFeedback — alqoritm
```go
func computeFeedback(guess, solution []rune) feedback {
    result := make(feedback, len(guess))
    used := make([]bool, len(solution))
    // 1) hamısı absent başlayır
    // 2) eyni mövqe = correctPosition; used[i]=true
    // 3) qalan hərflər üçün yalnız İSTİFADƏSİZ mövqelərdə wrongPosition
    return result
}
```
- "double character" halları: təkrar hərf yalnız bir dəfə "wrong" işarələnir —
  pseudo-kod + kağız üzərində dizayn TƏLKİM OLUNUR

### 7. Korpus oxunuşu
```go
func ReadCorpus(path string) ([]string, error) {
    data, err := os.ReadFile(path)  // tam fayl — []bayt
    if err != nil {
        return nil, fmt.Errorf("unable to open %q for reading: %w", path, err)
    }
    words := strings.Fields(string(data))
    if len(words) == 0 { return nil, ErrCorpusIsEmpty }
    return words, nil
}
```

### 8. Play loop
```go
func (g *Game) Play() {
    fmt.Println("Welcome to Gordle!")
    for currentAttempt := 1; currentAttempt <= g.maxAttempts; currentAttempt++ {
        guess := g.ask()
        fb := computeFeedback(guess, g.solution)
        fmt.Println(fb.String())
        if slices.Equal(guess, g.solution) {
            fmt.Printf(" 🎉 You won! The word was: %s.\n", string(g.solution))
            return
        }
    }
    fmt.Printf(" 😞 You've lost! The solution was %q.\n", string(g.solution))
}
```

## Əsas terminlər

- Pointer Receiver
- Rune (Unicode kod nöqtəsi)
- Sentinel Error
- fmt.Stringer
- strings.Builder
- Corpus (sözlük korpusu)
- slices.Equal (slice müqayisəsi)

## Praktik nəticə

- İstənilən mətn emalı üçün []rune əsaslı düşün
- Kiçik alqoritmləri pseudo-kodla kağızda QUR, sonra kodlaşdır
- Çap formatını Stringer-ə həvalə et — test və UX bir yerdə
- `strings.Fields` — boşluqlarla ayırma (Split-in yumşaq forması)

## Mənbə

Pages: 154-202 (Chapter 5, Learn Go with Pocket-Sized Projects)
