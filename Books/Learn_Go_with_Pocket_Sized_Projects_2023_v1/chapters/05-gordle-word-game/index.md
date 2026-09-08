# Chapter 5 — Gordle: Play a word game in your terminal (səh. 154-202)

## Bu chapter nədən bəhs edir?

Wordle-in terminal versiyası (Gordle): rune-larla işləmə, pointer receiver,
`io.Reader` asılılığı, xəta sarma (%w), feedback hesablama alqoritmi (absent /
wrong position / correct), `Stringer` interfeysi, corpus oxuma, sentinel xətalar
və `crypto/rand` ilə təsadüfi söz seçmə.

## Əsas fikirlər

### 1. Oyun strukturu və pointer receiver
```go
type Game struct {
    reader      io.Reader // asılılıq — testdə strings.Reader
    solution    []rune    // gizli söz
    maxAttempts int
}

func New(reader io.Reader, corpus []string, maxAttempts int) (*Game, error) {
    if len(corpus) == 0 {
        return nil, ErrCorpusIsEmpty
    }
    return &Game{reader: reader, solution: pickWord(corpus), maxAttempts: maxAttempts}, nil
}
```
**Sub-kod izahı:**
- `New` məcburi parametrləri (reader, corpus) qəbul edir — konstruktorda
  validasiya: boş corpus → sentinel xəta
- Pointer receiver istifadə olunur: (1) metodlar Game state-ni dəyişir, (2)
  metod dəstində tənzimləmək asan olur

### 2. Rune-larla input oxuma
```go
func (g *Game) ask() []rune {
    fmt.Printf("Enter a %d-character guess: ", solutionLength)
    scanner := bufio.NewScanner(g.reader) // testdən gələn reader
    if !scanner.Scan() {
        return ask() // retry loop
    }
    guess := []rune(strings.ToUpper(scanner.Text()))
    if len(guess) != solutionLength {
        // ... → yenidən soruş
    }
    return guess
}
```
**Sub-kod izahı:**
- `[]rune(string)` → simvollar (code point-lər), baytlar DEYİL — "ПРИВЕТ" 6
  rune-dur, amma bayt sayı çoxdur
- `strings.ToUpper` → müqayisə uniform olsun
- `len(guess)` = rune sayı (çünki []rune)

### 3. Xəta wrapping + sentinel
```go
var ErrCorpusIsEmpty = corpusError("corpus is empty") // typed sentinel

type corpusError string
func (e corpusError) Error() string { return string(e) }

err = g.validateGuess(guess)
if err != nil {
    _, _ = fmt.Fprintf(os.Stderr, "...invalid: %s.\n", err)
}
```
- `%w` (`fmt.Errorf`) → xəta zənciri; `errors.Is` ilə müqayisə
- Typed sentinel (`corpusError`) → `errors.Is` üçün sabit, müqayisə edilə bilən
  xəta dəyəri

### 4. Feedback status enum + Stringer
```go
type hint rune
const (
    absentCharacter hint = '_' // yoxdur
    wrongPosition  hint = 'C' // yanlış yer
    correctPosition hint = 'G' // düz yer
)

func (h hint) String() string { // fmt.Stringer implementasiyası
    switch h {
    case absentCharacter: return "⬜"
    case wrongPosition:  return "🟨"
    case correctPosition: return "🟩"
    }
    return " "
}
```
`String() string` metodlu hər tip `fmt.Stringer`-i **örtülü** realləşdirir —
`fmt.Println` avtomatik çağırır. `strings.Builder` ilə səmərəli birləşdirmə
(hər append-da yeni yaddaş ayrılmır).

### 5. computeFeedback alqoritmi
İki mərhələli pseudo-code (karton üzərində düşüncə):
```go
func computeFeedback(guess, solution []rune) feedback {
    fb := make(feedback, len(guess))
    for i := range fb {
        fb[i] = absentCharacter      // 1. hamısı absent kimi işarələ
    }
    for i, char := range guess {
        if solution[i] == char {
            fb[i] = correctPosition // 2a. düz mövqe
        } else {
            for j, target := range solution {
                if char == target && fb[j] != correctPosition {
                    fb[i] = wrongPosition // 2b. sözdə var, yeri yanlış
                    break
                }
            }
        }
    }
    return fb
}
```
**VACİB edge case:** təkrarlanan hərflər — bir hərf iki dəfə `correct` işarə
ala bilməz; pseudo-code + kağız üzərində analiz əvvəlcədən bunu ortaya çıxarır.

### 6. Corpus oxuma + pickWord
```go
func ReadCorpus(path string) ([]string, error) {
    data, err := os.ReadFile(path) // fayl = baytlar
    if err != nil {
        return nil, fmt.Errorf("unable to read corpus: %w", err)
    }
    words := strings.Fields(string(data)) // whitespace-ə görə böl
    if len(words) == 0 {
        return nil, ErrCorpusIsEmpty
    }
    return words, nil
}

func pickWord(corpus []string) []rune {
    index := rand.Intn(len(corpus)) // math/rand — seed məcburi
    return []rune(strings.ToUpper(corpus[index]))
}
```
`math/rand` deterministikdir (seed-lə); `crypto/rand` təhlükəsizdir amma
əsassız burada. Testdə `pickWord`-ün çıxışının corpus-da olmasını yoxlamaq
kifayətdir (random-luğu test etmək olmaz).

### 7. Unicode mürəkkəbliyi
Bəzi dillərdə hərf = çoxlu code point (inkişaf etmiş formalar). `golang.org/x/text`
paketindəki `unicode/norm.Iter` normalizasiya üçün — `golang.org/x` Go-nun
genişləndirilmiş (dilin xaricində) ekosistemdir.

## Əsas terminlər

- Rune (simvol kodu)
- Pointer Receiver (işarəçi alıcısı)
- io.Reader asılılığı
- Sentinel Error (möhür xətası)
- Error Wrapping (%w)
- Stringer (fmt.Stringer)
- Corpus (söz toplusu)
- strings.Builder

## Praktik nəticə

- Simvol sayma = `[]rune` çevirib `len`
- Test olunması üçün bütün I/O asılılıqlarını (reader) parametr kimi ötür
- Feedback tipli tapşırıqlarda pseudo-code + edge case analizi əvvəlcədən
- Statunu dəyişən metodlar → pointer receiver

## Mənbə

Pages: 154-202 (Chapter 5, Learn Go with Pocket-Sized Projects)
