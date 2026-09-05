# Chapter 4 — Handling errors and panics (Technique 16-21)

## Bu chapter nədən bəhs edir?

Go error idiomları (error həmişə sonuncu qayıdış, nil minimizasiyası, custom error tipləri, package-scoped error dəyişənləri) və panic sistemi (error vs panic fərqi, panic-ə error ötürmək, defer + recover, deferred closure scope, goroutine panic tələsi, safely.Go pattern).

## Əsas fikirlər

### 1. Error əsası — sonuncu return
**Kitabdan kod nümunəsi:**
```go
func Concat(parts ...string) (string, error) {
    if len(parts) == 0 {
        return "", errors.New("No strings supplied")
    }
    return strings.Join(parts, " "), nil
}

func main() {
    args := os.Args[1:]
    if result, err := Concat(args...); err != nil {     // if GET READY; EVALUATE
        fmt.Printf("Error: %s\n", err)
    } else {
        fmt.Printf("Concatenated string: '%s'\n", result)
    }
}
```
- **İdiom:** error HƏMİŞƏ sonuncu return value
- `if result, err := ...; err != nil` — assignment + yoxlama bir sətirdə; dəyişənlər else-də də scope-dadır — `if a = b` (assignment = müqayisə qarışması) bug-unun qarşısını alır
- **Variadic:** `parts ...string` → daxildə `[]string`; çağırışda `args...` → slice açılır

### TECHNIQUE 16: Minimize the nils
**Problem:** error + nil qaytarmaq istifadəçiyə əlavə iş yükləyir.

**Həll:** xəta olanda da İŞLƏNƏ BİLƏN dəyər qaytar:
```go
// Concat xətasında "" qaytarır → istifadəçi error-u ignorə edə bilər:
result, _ := Concat(args...)
fmt.Printf("Concatenated string: '%s'\n", result)   // error olmadan istifadə!
```
**Qayda:** faydalı nəticə qaytarıla bilirsə — qaytar; heç nə yoxdursa — nil.

**Error yaratma alətləri:**
- `errors.New("mesaj")` — sadə
- `fmt.Errorf("format %s", arg)` — formatlı

**Sənədləşdirmə konvensiyası:**
```go
// Concat concatenates a bunch of strings, separated by spaces.
// It returns an empty string and an error if no strings were passed in.
func Concat(parts ...string) (string, error) { ... }
```

### TECHNIQUE 17: Custom error tipləri
**error interfeysi:**
```go
type error interface {
    Error() string
}
```
**Nə vaxt:** xəta ƏLAVƏ MƏLUMAT daşımalıdır.

**Kitabdan kod nümunəsi:**
```go
type ParseError struct {
    Message    string
    Line, Char int
}

func (p *ParseError) Error() string {
    format := "%s on Line %d, Char %d"
    return fmt.Sprintf(format, p.Message, p.Line, p.Char)
}
```
→ Line/Char sahələri ilə çağırıcı xəta sətrini bərpa edib göstərə bilər (sadəcə mətndən daha güclü).

**Qeyd:** Go developer-ləri nadirən custom error yaradır — core kitabxanalar generic error istifadə edir; əksər xətanın xüsusi atributları yoxdur.

### TECHNIQUE 18: Error dəyişənləri (package-scoped)
**Problem:** bir funksiya bir neçə MƏNALI xəta qaytara bilər, amma heç biri əlavə məlumat daşımır.

**Həll:** paket səviyyəsində error DƏYİŞƏNLƏRİ — standart kitabxanada `io.EOF` kimi.

**Kitabdan kod nümunəsi:**
```go
var ErrTimeout = errors.New("The request timed out")
var ErrRejected = errors.New("The request was rejected")
var random = rand.New(rand.NewSource(35))    // fixed seed — dəqiq nümunə üçün

func SendRequest(req string) (string, error) {
    switch random.Int() % 3 {
    case 0:
        return "Success", nil
    case 1:
        return "", ErrRejected
    default:
        return "", ErrTimeout
    }
}

func main() {
    response, err := SendRequest("Hello")
    for err == ErrTimeout {                  // DƏYİŞƏN MÜQAYİSƏSİ!
        fmt.Println("Timeout. Retrying.")
        response, err = SendRequest("Hello")
    }
    if err != nil {
        fmt.Println(err)
    } else {
        fmt.Println(response)
    }
}
```
- `err == ErrTimeout` — müqayisə sadə equality
- Bir dəfə instantiate → səmərəli; konseptual olaraq sadə
- Java/Python `throw new Instance()` + try/catch YOX — Go yolu: dəyişən + ==

### 2. Error vs Panic — konseptual fərq
| | Error | Panic |
|---|---|---|
| Məna | gözlənilən pozuntu | sistem davam edə BİLMƏZ |
| Gözlənmə | sənədləşdirilmiş, gözlənilən | gözlənilməz |
| Kim idarə edir | proqramçı (ignore → heç nə olmur) | runtime — stack unwind, handler axtarır, tapmasa PROQRAM ÖLÜR |
| Nümunə | fayl açılmadı, timeout | divide by zero, array out of bounds |

**Kitabdan kod nümunəsi:**
```go
var ErrDivideByZero = errors.New("Can't divide by zero")

func precheckDivide(a, b int) (int, error) {
    if b == 0 {
        return 0, ErrDivideByZero        // ERROR — yoxlanmış
    }
    return divide(a, b), nil
}

func divide(a, b int) int {
    return a / b                         // b=0 → PANIC!
}

// main-də:
// precheckDivide(1, 0) → "Error: Can't divide by zero" ✓
// divide(2, 0) → panic: runtime error: integer divide by zero → CRASH
```

**Qayda:** "don't panic unless there's no clear way to handle the condition within the present context. When possible, return errors instead."

### TECHNIQUE 19: Panic-ə NƏ ötürmək
`panic(interface{})` — hər şey qəbul edir. Amma:
- `panic(nil)` → "panic: nil" — işlənənməz
- `panic("string")` → bir az yaxşı
- **`panic(errors.New("..."))` → İDIOMATİK ✓**

**Səbəblər:**
1. Intuitiv — panikə səbəb olan şey xətadır
2. **recover-də asanlıqla error sistemə qaytarıla bilər**

```go
panic(errors.New("Something bad happened."))
```

### TECHNIQUE 20: Recover — defer ilə
**defer əsasları:**
```go
func main() {
    defer goodbye()         // main qayıdanda icra
    fmt.Println("Hello world.")   // ƏVVƏL bunu çap et
}
// Output: Hello world. / Goodbye
```
Fayl/socket qapatmaq, resource azad etmək, panic tutmaq — üçün ideal.

**Sadə recover:**
```go
func main() {
    defer func() {
        if err := recover(); err != nil {        // panic VARSA dəyər qaytarır
            fmt.Printf("Trapped panic: %s (%T)\n", err, err)
        }
    }()
    yikes()       // panic edir
}
// Trapped panic: Something bad happened. (*errors.errorString)
```
- `recover()` → panic yoxdursa nil; varsa panic-ə ötürülən dəyər

**Deferred closure scope qaydaları:**
```go
// DOĞRU — msg ƏVVƏLCƏDƏN elan edilib:
var msg string
defer func() { fmt.Println(msg) }()
msg = "Hello world"     // dəyişiklik closure GÖRÜR → "Hello world"

// SƏHV — msg SONRADAN elan:
defer func() { fmt.Println(msg) }()    // COMPILE ERROR!
msg := "Hello world"
```
Closure ELAN ZAMANINDA scope-u görür; icra sonra olsa da, elandan sonrakı dəyişənlər görünmür.

**Tam nümunə — panic-i error-a çevir + təmizlə:**
```go
func OpenCSV(filename string) (file *os.File, err error) {   // ADLI return!
    defer func() {
        if r := recover(); r != nil {
            file.Close()                    // təmizlə
            err = r.(error)                 // panic-i error-a type assert → qaytış!
        }
    }()
    file, err = os.Open(filename)
    if err != nil {
        fmt.Printf("Failed to open file\n")
        return file, err
    }
    RemoveEmptyLines(file)       // həmişə panic edir (demo)
    return file, err
}
```
- **Named return values** — closure içindən file/err manipulyasiyası mümkün
- `r.(error)` — type assertion: panic dəyərini error-a çevirir
- Pattern: panic → recover → təmizlə → normal error qaytışı

**Defer qaydaları:**
- Funksiyanın BAŞINA yaxın qoy
- Sadə elanlar (x:=1) defer-dən əvvəl; mürəkkəb var-lar elan əvvəl + init sonra
- Çox defer bir funksiyada — eybəcər (frowned upon)
- Fayl/şəbəkə/resource qapatma → həmişə defer

### TECHNIQUE 21: Goroutine panic tələsi
**Problem:** hər goroutine ÖZ function stack-inə malikdir. **Panic stack-lərarası KEÇƏ BİLMƏZ** — handle-dəki panic listen-dəki recover-a çatmır!

```
main() → start() → listen() ──go handle()──→ handle() → response() → panic!
                                              ↑ BU stack-də recover YOXDURSA
                                              → PROQRAM ÇÖKÜR
```

**Echo server (goroutine + panic təhlükəsi):**
```go
func listen() {
    listener, err := net.Listen("tcp", ":1026")
    if err != nil {
        fmt.Println("Failed to open port on 1026")
        return
    }
    for {
        conn, err := listener.Accept()
        if err != nil {
            fmt.Println("Error accepting connection")
            continue
        }
        go handle(conn)              // hər connection öz goroutine
    }
}

func handle(conn net.Conn) {
    reader := bufio.NewReader(conn)
    data, err := reader.ReadBytes('\n')
    if err != nil {
        fmt.Println("Failed to read from socket.")
        conn.Close()
    }
    response(data, conn)             // panic etsə → SERVER ÖLÜR!
}
```

**Həll 1 — handler-da recover:**
```go
func handle(conn net.Conn) {
    defer func() {
        if err := recover(); err != nil {
            fmt.Printf("Fatal error: %s", err)
        }
        conn.Close()                 // hər halda bağla
    }()
    reader := bufio.NewReader(conn)
    ...
    response(data, conn)             // panic → recover tutur → server YAŞAYIR
}
```

**Həll 2 — net/http dərsi:** handler panic etsə server DAVAM EDİR — kitabxana ÖZÜ recover edir:
```
2015/04/08 07:57:31 http: panic serving [::1]:51178: Fake panic!
... (server işləməyə davam)
```
→ **Server kitabxanaları panic-handling-i ÖZ İÇİNDƏ saxlamalıdır.**

**Həll 3 — safely.Go pattern:**
```go
package safely

import (
    "log"
)

type GoDoer func()

func Go(todo GoDoer) {
    go func() {
        defer func() {
            if err := recover(); err != nil {
                log.Printf("Panic in safely.Go: %s", err)
            }
        }()
        todo()                   // istənilən funksiya — recover-li mühitdə!
    }()
}

// İstifadə:
safely.Go(message)    // go message ƏVƏZİNƏ — panic avtomatik tutulur
```
- `GoDoer` — parametrsiz funksiya tipi
- `safely.Go(fn)` → goroutine + daxili recover → panic loglanır, proqram yaşayır
- **Ders:** developer-in yaddaşına bel bağlama — preventativ library yaz

## Error/Panic qərar cədvəli
| Vəziyyət | Error | Panic |
|---|---|---|
| Gözlənilən kənar (fayl, şəbəkə, istifadəçi inputu) | ✓ | |
| Bərpa mümkündür | ✓ | |
| Sistem davam edə bilmir | | ✓ |
| Muəyyən edilməz invariant pozuntusu | | ✓ |
| Library-dəki goroutine | ✓ (recover daxildə) | |

## Əsas terminlər
- Error as Last Return Value (idiom)
- if assignment clause (`if v, err := ...; err != nil`)
- Variadic (`...T` / `args...`)
- errors.New / fmt.Errorf
- Minimize the Nils (işlənən dəyər + error)
- Custom Error Type (struct + Error() metodu)
- Package-scoped Error Variable (ErrTimeout, io.EOF)
- Error Equality Check (`err == ErrTimeout`)
- Panic vs Error (gözlənilməz vs gözlənilən)
- Stack Unwinding
- panic(errors.New(...)) — idiomatik ötürmə
- defer / Deferred Closure
- recover (yalnız defer daxilində)
- Deferred Closure Scope (elandan SONRAKI dəyişənlər görünmür)
- Named Return Values + recover
- Type Assertion (`r.(error)`)
- Goroutine Function Stack (panic keçə bilməz)
- Server-side Panic Handling (net/http dərsi)
- GoDoer / safely.Go (panic-safe goroutine wrapper)

## Praktik nəticə
- Xəta olanda da işlənə bilən dəyər qaytar ("" və s.) — istifadəçi `_` ilə keçə bilsin.
- Fərqli xəta halları üçün `ErrX` dəyişənləri + `==` müqayisəsi — try/catch-siz Go yolu.
- Custom error yalnız əlavə məlumat lazım olanda — əks halda generic error kifayətdir.
- Panic-ə error ötür — recover-də `r.(error)` ilə normal error axınına qaytar.
- Recover YALNIZ defer daxilində işləyir; closure elanından sonrakı dəyişənləri görmür — named return-larla birləşdir.
- Goroutine başlatdıqda panic-in ÖZ stack-də qaldığını bil — handler-da recover YOXSA server çökür; net/http bunu daxilə alır — sənin library-lərin də etsin.
- safely.Go pattern — proqramçı yaddaşına deyil, library quruluşuna güvən.

## Mənbə
Pages: 110-134 (PDF), book pages 87-111
