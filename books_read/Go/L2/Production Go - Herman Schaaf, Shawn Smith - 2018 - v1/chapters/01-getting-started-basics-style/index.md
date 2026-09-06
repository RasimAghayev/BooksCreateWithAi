# Chapters 1-4 — Introduction, Getting Started, Basics, Style & Error Handling

## Bu bölmələr nədən bəhs edir?

Kitabın məqsədi (production-ready Go), quraşdırma + editor/linter tooling, tam sintaksis icmalı (data tipləri, operatorlar, struct, slice, map, loop, funksiya, pointer, goroutine, kanal, interfeys, error) və idiomatik üslub (gofmt, goimports, golint, qısa dəyişən adları, error handling fəlsəfəsi).

**PDF səhifələr:** 6-57 (Intro: 6, Getting Started: 7-12, Basics: 13-51, Style & Error Handling: 52-57)

## Əsas fikirlər

### 1. Introduction (Ch 1)
**Niyə production-da Go?**
- Müəlliflərin PHP→Go təcrübəsi: **10x azaldılmış cavab vaxtı**, yüksək user retention, azalmış server xərcləri
- Developer xoşbəxtliyi — təhlükəsizlik zəmanətləri production-qırıcı bug-ları azaldır
- "Works → just works" — sadə dizayn + oxunaqlıq > clever constructs

**Auditoriya:** beginner proqramçılar YOX — CS əsaslarını bilən, ilk dəfə Go yazan mühəndislər.

### 2. Getting Started (Ch 2)

**Quraşdırma:** binary release (golang.org/dl) / Homebrew (`brew install go`) / yoxlama: `go version`.

**GOPATH:** Go 1.8+ default `$HOME/go` (Unix) / `%USERPROFILE%/go` (Windows).
```bash
go get github.com/gorilla/mux     # $HOME/go/src/github.com/gorilla/mux-a endirir
```

**gorilla/mux nümunəsi:**
```go
package main

import (
    "fmt"
    "net/http"

    "github.com/gorilla/mux"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello, world")
}

func main() {
    r := mux.NewRouter()
    r.HandleFunc("/", HomeHandler)
    http.Handle("/", r)
    http.ListenAndServe(":8080", r)
}
```

**Editor tooling:**
- **goimports on save** — hər editorda tövsiyə: format + import-ları avtomatik əlavə/sil
  ```bash
  go get golang.org/x/tools/cmd/goimports
  ```
- GoLand (File Watchers → goimports) / Sublime (GoSublime) / vim (`let g:go_fmt_command="goimports"`)

**Linter suite — gometalinter:**
```bash
$ go get -u github.com/alecthomas/gometalinter
$ gometalinter --install
$ gometalinter --deadline=180s --exclude=vendor --disable-all --enable=misspell ./...
```
Faydalılar: deadcode, ineffassign, misspell.

**go vet (MÜTLƏQ):** correctness yoxlaması:
```bash
$ go vet main.go
main.go:8: arg "test" for printf verb %d of wrong type: string
main.go:8: possible formatting directive in Println call
```
- `fmt.Printf("%d", "test")` → runtime-da `%!d(string=test)` — vet compile-dan ƏVVƏL tutur
- Println içində printf feyli → `Println("%d", x)` çap etmir — vet xəbərdarlıq edir

### 3. Basics (Ch 3)

**Program structure:** package main + func main; `go run` / `go build` + `go install` → $GOPATH/bin.
Paketə çevirmə: `package dinner` + `func Choose() string` (exported!) → başqa yerdən `dinner.Choose()`.

**Variables — 3 yol + 2 incəlik:**
```go
var x int        // zero value (0)
var x = 1        // tip çıxarılır
x := 1           // short — ƏN ÇOX istifadə
```
**İncəlik 1:** `:=` funksiya XARİCİNDƏ YOX — global scope-da yalnız `var`.
**İncəlik 2 (redeclare):**
```go
b, c := 1, 2    // OK — c YENİDİR (b mövcuddursa belə)
var b, c = 1, 2  // XƏTA — b artıq var (redeclared in this block)
```

**Data tipləri:**
- bool (yalnız true/false — digər tiplərə truthy/falsey YOX)
- string (çoxsətirli müzakirə Strings chapterda)
- int8-64 / uint8-64 / int / uint (alias int32/int64 — architecture)
- float32/64 (IEEE-754); complex64/128 (`complex(1,2)`, real/imag, `1.3i`)
- byte = uint8, rune = int32

**Tiplər qarışmır:** `i32 + i64` → `mismatched types` — **explicit cast mütləq**: `int64(i32) + i64`.
**Konstant overflow-i kompilyasiya saxlayır:** `var i int8 = 128` → `constant 128 overflows int8`.

**⚠️ uint UNDERFLOW tələsi:**
```go
var p, s uint32 = 0, 1
fmt.Println(p - s)     // 4294967295 — -1 YOX, wrap-around!
```
Uint seçməzdən ƏVVƏL: len/cap int qaytarır — cast-lar əlavə olunacaq.

**Struct + JSON tag:**
```go
type Person struct {
    Name  string `json:"name"`     // kiçik hərflə çıxar
    Email string `json:"email"`
}
b, _ := json.Marshal(p)             // {"name":...,"email":...}
```
- Sahələr EXPORTED olmalı — yoxsa **boş JSON {}** qalır!
- Anonymous struct: `p := struct{ Name string `json:"name"` }{...}` — table-driven testlərdə məşhur.

**Operatorlar (spec-dən):**
- Arithmetic: `+ - * / %` + `& | ^ &^ << >>` (yalnız int-lər); `+` string-lərə də
- Comparison: tip UYQUNLUĞU şərt — `int == int32` compile xətası
- Logical: `&& || !` — sağ operand conditional (short-circuit)
- Address: `&x` (pointer), `*p` (deref, nil → panic)
- Receive: `<-` (kanallar)

**Şərtlər:**
```go
if got := string(buf); got != want { ... }   // scoped declaration
switch { case '0' <= c && c <= '9': ... }    // şərtsiz = switch true
switch verb { case 't', 'v': ... default: }   // şərtli
switch err := err.(type) { case NotFoundError: ... }   // type switch
```

**Arrays/Slices:**
```go
x := []int{1, 2, 3, 4, 5}
y := x[:3]          // [1 2 3] — slicing
var x []int          // nil (zero value)
x = append(x, 1, 2, 3)
x = append(x, y...)  // slice aç
copy(y, x)           // y-nin CAPACITY olmalı! (yoxsa kopyalanmaz)
sort.Slice(c, func(i, j int) bool { return c[i].Name < c[j].Name })   // comparator
```

**Maps:**
```go
m := make(map[string]string)
m["en"] = "Hello"
v := m["zh"]                    // "" — zero value
if _, ok := m["ja"]; ok { ... } // mövcudluq
delete(m, "en")
```
- **Concurrency-safe DEYİL** → sync.Map
- İterasiya sırası QARANTIYASIZ
- `var m map[string]string` (nil) → yazmaq PANIC

**Loops:** yalnız `for`:
```go
for x < 4 { ... }          // while
for i := 0; i < 3; i++ { }  // klassik
for { }                      // sonsuz
for i, x := range nums { }  // index + value
for _, x := range nums { }  // index ignore
```

**Pointers:** `var p *int` (nil); `p = &x`; `*p` = 100. **Pointer aritmetikası YOXDUR** (GC sadəliyi + safety).

**Goroutines:**
```go
go hello()                    // lightweight
go func() { ... }()           // anonymous
```
**sync.WaitGroup:**
```go
var wg sync.WaitGroup
wg.Add(3)
go func() { defer wg.Done(); ... }()
...
wg.Wait()
```
**Vacib:** Add(2) → vaxtından əvvəl qayıt; Add(4) → `fatal error: all goroutines are asleep - deadlock!` — SAYI DƏQİQ bil!

**Channels:**
```go
ch := make(chan int)          // unbuffered
ch := make(chan int, 5)       // buffered
ch <- 5; v := <-ch
```
- Buferə 6-cı göndəriş → deadlock; 1 oxu → OK
- **len(ch)** → bufer doluluğu; **select + default** → doludursa ötür:
```go
select {
case ch <- 6:
default:
    fmt.Println("channel is full, ignoring send")
}
```

**Interfaces:**
```go
type Entry interface { Title() string }
type Book struct{ ... }
func (b Book) Title() string { ... }       // implicit — implements!
func Display(e Entry) string { return e.Title() }   // hər iki tip keçər
```
- Empty interface: `var i interface{}` → hər tip
- Type assertion: `s, ok := i.(string)` — ok=false → zero value

**⚠️ Nil interface tələsi (kitabın vacib dərsi):**
```go
var s *string                 // s == nil: true
var i interface{} = s         // i == nil: FALSE!!!
// interface = (tip, dəyər) CÜTÜ — hər ikisi nil olmalı
fmt.Printf("%T, %v\n", i, i)  // *string, <nil>
```

**Error handling:**
```go
func ParseBool(str string) (bool, error) {
    switch str {
    case "1", "t", "T", "true", "TRUE", "True":
        return true, nil
    case "0", "f", "F", "false", "FALSE", "False":
        return false, nil
    }
    return false, fmt.Errorf("invalid input %q", str)
}
```

**stdin oxu:**
```go
scanner := bufio.NewScanner(os.Stdin)
for scanner.Scan() { m[scanner.Text()]++ }
if err := scanner.Err(); err != nil { ... }
```

**Fayl yaz:**
```go
f, err := os.Create("langs.txt")     // truncate!
defer f.Close()
f.WriteString(lang + "\n")
// VEYA: ioutil.WriteFile(path, data, perm) — bir çağırışda
```

### 4. Style & Error Handling (Ch 4)

**Gofmt — istisnasız standart:**
- PEP8-dən fərqli: exception YOX; tab YOX (məcburi); 80 simvol limiti YOX
- **Tövsiyənin əsası:** hər kəsin kodu eyni görünür → debug/contribute asanlaşdı → Go community böyüdü
- `gofmt -s` — simplify: `[]Animal{Animal{...}, Animal{...}}` → `[...]{{...}, {...}}` (mantıq dəyişməz, HƏMİŞƏ təhlükəsiz)

**goimports** — gofmt + import idarəsi; on-save hook əvəzinə istifadə et.

**Qısa dəyişən adları (idiom):**
> "Prefer c to lineCount. Prefer i to sliceIndex." — Go CodeReviewComments

**Qayda:** <10 sətir span → 1 simvol; >10 → deskriptiv. Funksiya <15 sətir hədəfi.

**golint:**
```bash
go get -u github.com/golang/lint/golint
```
- gofmt REFORMAT edir; golint XƏBƏRDARLIQ çap edir (mütləq standart deyil, false positive olur)
- `exported type Entry should have comment or be unexported` — godoc üçün
- `if block ends with a return statement, so drop this else` → else-i at, outdent et:
```go
// PİS:
if len(s) < l { return s } else { return s[:l] + suf }
// YAXŞI:
if len(s) < l { return s }
return s[:l] + suf
```

**Error handling fəlsəfəsi:**
```go
func lineCount(filepath string) (int, error) {
    out, err := exec.Command("wc", "-l", filepath).Output()
    if err != nil {
        return 0, fmt.Errorf("could not run wc -l: %s", err)   // spesifikləşdir!
    }
    count, err := strconv.Atoi(strings.Fields(string(out))[0])
    if err != nil {
        return 0, fmt.Errorf("could not convert wc -l output to integer: %s", err)
    }
    return count, nil
}
```
- `if err != nil` təkrarları — İDİOMATİK; try/except YOX
- **Strictlik faydası:** xətanın HARDA və NİYƏ baş verdiyini qaçırmaq mümkün deyil
- **Error mesajı KİÇİK hərflə başla:** `could not run` — log-da axının ortasında BÖYÜK hərf eybəcər görünür:
```
ERROR: lineCount("..."): Could not run wc -l   ← PİS
ERROR: lineCount("..."): could not run wc -l   ← YAXŞI
```

**Sonda:** Go source-u OXU — gofmt-dən oxunaqlı; godoc.org-dan funksiya adına kliklə source-a düş.

## Əsas terminlər
- goimports (format + import idarəsi)
- gometalinter / deadcode / ineffassign / misspell
- go vet (printf verb, formatting directive)
- golint (else at, exported comment)
- gofmt -s (simplify)
- Redeclare (:= əksər yeni olmalı)
- uint Underflow (0-1 = max)
- mismatched types (int32 + int64)
- constant overflows int8
- sync.Map (concurrent map)
- Nil Interface Tələsi ((tip, dəyər) cütü)
- Type Assertion (s, ok := i.(string))
- Anonymous Struct (table-driven testlər)
- CodeReviewComments (qısa adlar)

## Praktik nəticə
- goimports-u editor save hook et; go vet-i CI-da MÜTLƏQ gəzdir — ikisi birlikdə ən ucuz bug yaxalayıcıdır.
- uint seçməzdən əvvəl 2 dəfə düşün — len/cap int-dir, underflow tələsi real.
- Interface nil yoxlamasında tip+dəyər cütünü unutma — nil pointer interface-ə qoyulsa nil OLMAZ.
- WaitGroup Add sayı dəqiq goroutine sayı olmalı — azı erkən qayıdış, çoxu deadlock.
- Error-ları spesifikləşdir (fmt.Errorf kontekst əlavə et), kiçik hərflə başla.
- Else return-dən sonra artıqdır — golint tövsiyəsinə uyğunlaş.
- Go source oxu — ən yaxşı idiom dərsi.

## Mənbə
Pages: 6-57 (PDF), book pages i-51
