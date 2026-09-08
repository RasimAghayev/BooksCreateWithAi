# Chapters 1-2 — Meet Go, Hello Earth (səh. 33-78)

## Bu fəsillər nədən bəhs edir?

Go-ya giriş: dilin tarixi və fəlsəfəsi (sadəlik — 25 açar söz), müasir
sənaye ehtiyacları (backend, cloud), test/benchmark/fuzz alətləri. Sonra ilk
layihə: hello world-dən çoxdilli salamlama CLI-sına — test, custom tip, switch,
map, table-driven test, flag paketi.

## Əsas fikirlər

### 1. Go nədir? (Ch1)
- **Mənşə:** Google-də böyük miqyaslı problemlər üçün: yavaş build, asılılıq
  idarəsi, mürəkkəblik. **Sadəlik idarə edir** — 25 reserved keyword (2024).
- **Müasir sənaye üçün:** unit test, benchmarking, fuzzing, formatting (gofmt)
  — hamısı daxili alətlər.
- **Generics:** 1.18-dən — type-safe yenidən istifadə, boilerplate azaldır.
- **Harada YARAMIR:** OS yazmaq (GC — yaddaş nəzarəti məhdud), desktop GUI zəif.
- **Dil müqayisəsi:** C++/Python/Java vs Go — statik tip, errors-are-values,
  multiparadigm, daxili test alətləri.
- **Pedaqogika:** John Dewey — "etmək öyrənməyin ən yaxşı yoludur"; kitabın 12
  layihəsi bunun üzərində.

### 2. İlk proqram: hello, world + Example test (Ch2)
```go
package main

import "fmt"

func main() {
    greeting := greet()
    fmt.Println(greeting)
}

// greet returns a greeting to the world.
func greet() string {
    return "Hello world"
}
```
**Example test (stdout yoxlaması):**
```go
func ExampleMain() {
    main()
    // Output:
    // Hello world
}
```
- `// Output:` şərhindən sonrakı sətirlər standart çıxışla müqayisə olunur
- Example həm test, həm DOKUMENTASİYADIR (godoc-da görünür)

**Test faylı konvensiyaları:**
- `*_test.go` — yalnız `go test` zamanı kompayl olunur
- `package main` + `_internal_test.go` → unexported funksiyalara çıxış (internal
  test); `package main_test` → yalnız export olunmuş API (external test)

### 3. Custom tip + switch (çoxdillilik)
```go
// language represents the language's code
type language string  // custom tip — string-in "mənalı" versiyası

func greet(l language) string {
    switch l {
    case "en":
        return "Hello world"
    case "fr":
        return "Bonjour le monde"
    default:
        return ""
    }
}
```
**Sub-kod izahı:**
- `type language string` — nominal tip: `"en"` artıq təkcə string deyil, funksiya
  imzası dəqiq tələb edir
- Switch-də `break` LAZIM DEYİL — implicit

### 4. Test funksiyası (testing.T)
```go
func TestGreet_English(t *testing.T) {
    lang := language("en")      // preparation
    want := "Hello world"
    got := greet(lang)          // execution
    if got != want {            // assertion
        t.Errorf("expected: %q, got: %q", want, got)
    }
}
```
**Testin 4 mərhələsi:** preparation → execution → comparison → error reporting.
`t.Errorf` — test FAIL edir amma davam edir; `%q` — stringi sitatla çap edir.

### 5. Map (hash table) — phrasebook
```go
var phrasebook = map[language]string{
    "el": "Χαίρετε Κόσμε", // Greek
    "en": "Hello world",   // English
    "fr": "Bonjour le monde", // French
}

func greet(l language) string {
    greeting, ok := phrasebook[l]  // "comma ok" idiom-u
    if !ok {
        return fmt.Sprintf("unsupported language: %q", l)
    }
    return greeting
}
```
**Sub-kod izahı:**
- `v, ok := map[k]` — açar yoxdansa `ok=false`, `v` = zero value ("" string üçün)
- Switch əvəzinə map: O(1) axtarış + yeni dil = 1 sətir (kod dəyişmir)

### 6. Table-Driven Test (TDT)
```go
func TestGreet(t *testing.T) {
    tests := map[string]struct {   // ad → test case
        lang language
        want string
    }{
        "English":            {lang: "en", want: "Hello world"},
        "French":            {lang: "fr", want: "Bonjour le monde"},
        "Akkadian, unsupported": {lang: "akk", want: "unsupported language: \"akk\""},
    }
    for name, tc := range tests {
        tc := tc
        t.Run(name, func(t *testing.T) {  // subtest — ad görünən
            got := greet(tc.lang)
            if got != tc.want {
                t.Errorf("expected: %q, got: %q", tc.want, got)
            }
        })
    }
}
```
- Yeni dil əlavə etmək = 1 sətir test cədvəlində, test KODU dəyişmir
- `t.Run` — subtestlər ayrıca run/parallel ola bilər

### 7. Flag paketi (CLI parametrləri)
```go
var lang string
flag.StringVar(&lang, "lang", "en", "The required language, e.g. en, ur...")
flag.Parse()
greeting := greet(language(lang))
```
- `flag.StringVar(&var, name, default, help)` — pointer ilə (və ya
  `lang := flag.String(...)` — pointer qaytarır)
- `flag.Parse()` — OS args parslənir; default "en"
- `os.Args` ilə müqayisə: flag `--lang=fr` / `-lang fr` / `--lang fr` — hamısını
  avtomatik emal edir

## Əsas terminlər

- Reserved Keywords (açar sözlər — 25)
- Example Test (nümunə testi — stdout müqayisəsi)
- Internal/External Test (daxili/xarici test paketi)
- Custom Type / Nominal Typing (xüsusi tip)
- Comma Ok Idiom (v, ok := map[k])
- Table-Driven Test (cədvəl testi)
- Flag (komanda sətri parametri)

## Praktik nəticə

- Hər funksiya yaranan kimi test yaz — refactoring təhlükəsiz olur
- Mənalı tiplər (`type language string`) — imzada dokumentasiya
- Switch yox, map + comma-ok — çoxdilli cədvəllər üçün
- TDT: test data vs test kod ayrılır; yeni hal = 1 sətir

## Mənbə

Pages: 33-78 (Chapters 1-2, Learn Go with Pocket-Sized Projects)
