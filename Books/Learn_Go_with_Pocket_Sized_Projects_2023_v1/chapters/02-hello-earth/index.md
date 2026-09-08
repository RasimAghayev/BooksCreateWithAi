# Chapter 2 — Hello, earth! Extend your hello, world (səh. 46-78)

## Bu chapter nədən bəhs edir?

İlk proqramdan başlayaraq: `fmt.Println`, `greet()` funksiyası, `Example` testləri,
`*testing.T` ilə Test funksiyaları, yeni dil dəstəyi (switch → map), table-driven
testlər və `flag` paketi ilə CLI parametrləri.

## Əsas fikirlər

### 1. İlk proqram və görünüş qaydası
```go
package main

import "fmt"

func main() {
    greeting := greet("en")
    fmt.Println(greeting)
}
```
**Sub-kod izahı:**
- `package main` → icra olunan proqram (executable)
- Böyük hərflə başlayan simvollar (`Println`) **xarici istifadəyə açıqdır**
  (exported); kiçik hərflə başlayanlar yalnız paket daxilində (unexported)
- `:=` → qısa dəyişən elanı (type inference / tip çıxarımı)

### 2. Example testi — sənədləşmə + test birlikdə
```go
func ExampleMain() {
    main()
    // Output: Hello world
}
```
- `// Output:` şərhi `go test` tərəfindən oxunur: stdout ilə müqayisə edilir
- Ad qaydası: `Example<FunksiyaAdı>` → godoc-da funksiya yanında göstərilir

### 3. Test funksiyası — 4 fazalı struktur
```go
func TestGreet_English(t *testing.T) {
    lang := language("en") // 1. Preparation (hazırlıq)
    want := "Hello world"

    got := greet(lang)     // 2. Execution (icra)

    if got != want {       // 3. Assertion (yoxlama)
        t.Errorf("expected: %q, got: %q", want, got) // 4. Failure (uğursuzluq)
    }
}
```
**Sub-kod izahı:**
- `*testing.T` → test konteksti; `t.Errorf` testi FAILED edir (davam edir),
  `t.Fatalf` dərhal dayandırır
- Fayl adı `*_test.go` → build/vaxtı ignore olunur
- `_internal_test.go` → eyni paketdə, unexported funksiyaları yoxlaya bilir

### 4. Custom tip + switch-dən map-ə
```go
type language string // domain mənası olan ad

var phrasebook = map[language]string{
    "en": "Hello world",
    "fr": "Bonjour le monde",
    "el": "Χαίρεται Κόσμε",
    "ur": "ہیلو دنیا",
}

func greet(l language) string {
    greeting, ok := phrasebook[l] // "comma ok" — açar varsa true
    if !ok {
        return "unsupported language: " + string(l)
    }
    return greeting
}
```
**Sub-kod izahı:**
- `type language string` → sadə string deyil, adlandırılmış tip — kompiler
  səhv parametri yaxalayır
- `map` → hash table; `phrasebook[l]` iki dəyər qaytarır `(dəyər, ok)`
- switch-case zənciri əvəzinə data-driven həll; yeni dil əlavə etmək kodu
  DƏYİŞDİRMİРDƏN map-ə sətir əlavə etməkdir

### 5. Table-driven test (TDT)
```go
func TestGreet(t *testing.T) {
    type testCase struct {
        lang language
        want string
    }
    var tests = map[string]testCase{
        "English": {lang: "en", want: "Hello world"},
        "French":  {lang: "fr", want: "Bonjour le monde"},
        "Akkadian": {lang: "akk", want: "unsupported language: akkk"},
    }

    for name, tc := range tests {
        tc := tc // loop dəyişənin tutulması
        t.Run(name, func(t *testing.T) {
            got := greet(tc.lang)
            if got != tc.want {
                t.Errorf("expected: %q, got: %q", tc.want, got)
            }
        })
    }
}
```
**Sub-kod izahı:**
- `testCase` struct → hər ssenari üçün input + gözlənti bir yerdə
- `map[string]testCase` → açar = ssenari adı (test outputunda görünür)
- `t.Run(name, ...)` → subtest — ssenari başına ayrı PASS/FAIL
- Yeni ssenari = yeni map elementi, test funksiyası dəyişmir

### 6. Quotation mark-lar (sitat işarələri)
- `"double"` → literal string
- `` `raw` `` → escape-siz (backslash daxil)
- `'r'` → rune (bütün stringlər UTF-8, rune = code point)

### 7. flag paketi — CLI parametrləri
```go
var lang string
flag.StringVar(&lang, "lang", "en", "The required language, e.g. en, ur...")
flag.Parse()

greeting := greet(language(lang))
```
**Sub-kod izahı:**
- `flag.StringVar(&lang, ad, default, kömək)` → mövcud dəyişənə bağlanır
- `flag.String(...)` → pointer qaytarır (alternativ)
- `flag.Parse()` → OS args parslənir; default `en`
- İstifadə: `go run main.go -lang=fr` və ya `-lang fr`

## Əsas terminlər

- Exported / Unexported (görünən / gizli)
- Example Test (nümunə testi)
- Table-Driven Test (cədvəl testi)
- Subtest (`t.Run`)
- Hash Table / Map (xəritə)
- "Comma ok" idiomu
- Named Type (adlandırılmış tip)
- Flag (CLI parametri)

## Praktik nəticə

- Example testi həm sənədləşdirir, həm yoxlayır — kiçik funksiyalar üçün ideal
- Şərti məntiq (switch) → data (map) çevirmək genişlənməni asanlaşdırır
- Test strukturu: Preparation → Execution → Assertion → Failure
- `flag` ilə default-lu, köməkli CLI parametrləri

## Mənbə

Pages: 46-78 (Chapter 2, Learn Go with Pocket-Sized Projects)
