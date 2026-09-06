# Chapter 1 — Getting into Go

## Bu chapter nədən bəhs edir?

Go-nun dizayn fəlsəfəsi, əsas xüsusiyyətləri (çoxqayıdış, modern standart kitabxana, goroutine/channel, toolchain), digər dillərlə müqayisə (C, Java, Python/PHP, JS/Node.js), quraşdırma + workspace (GOPATH) və Hello World web server.

## Əsas fikirlər

### 1. Go nədir?
**Tərif:** Statik tipli, kompilyasiya olunan, açıq mənbəli dil — Robert Griesemer, Rob Pike, Ken Thompson (Google) tərəfindən **real dünya problemləri** üçün (nəzəri təmizlik üçün yox). İlham: C, Pascal, Smalltalk, Newsqueak, C#, JavaScript, Python, Java.

**3 qat (layers):**
1. **Dil:** müasir aparat arxitekturası üçün
2. **Toolchain:** testing, documentation, formatting, package management — hamısı daxili
3. **Ekosistem:** paket sistemi + Git inteqrasiyası → paket/kitabxana ekosistemi

**Simplicity qanunu:** hər özəllik 3 müəllifin RAZILIĞI ilə daxil olub. "Simple enough to keep in your head yet powerful enough to write a wide variety of software."

**Simplicity nümunəsi — dəyişən sintaksisi:**
```go
var i int = 2     // tam forma
var i = 2          // inference
i := 2              // short declaration — yarısı qədər uzunluq!
```

**Qeyd:** Go-da ternary operator (`?:`) və (kitabın yazıldığı dövrdə) generics YOXDUR — hər problem üçün başqa yaxşı həll var.

### 2. Multiple return values
**Kitabdan kod nümunəsi:**
```go
func Names() (string, string) {
    return "Foo", "Bar"
}

func main() {
    n1, n2 := Names()
    fmt.Println(n1, n2)     // Foo Bar
    n3, _ := Names()
    fmt.Println(n3)          // Foo — ikincisi ignorə
}
```

**Named return values:**
```go
func Names() (first string, second string) {
    first = "Foo"
    second = "Bar"
    return                  // dəyərsiz return — adlı dəyərlər qayıdır!
}
```
Bu pattern bütün Go kitabxanalarında (value, err kimi) standartdır.

### 3. Modern standart kitabxana
**Networking (TCP):**
```go
conn, _ := net.Dial("tcp", "golang.org:80")
fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
status, _ := bufio.NewReader(conn).ReadString('\n')
fmt.Println(status)
```

**HTTP client:**
```go
resp, _ := http.Get("http://example.com/")
body, _ := ioutil.ReadAll(resp.Body)
fmt.Println(string(body))
resp.Body.Close()
```

**Sahələr:**
- **Networking/HTTP:** TCP/UDP dial/listen; http client (proxy, TLS, cookies, transport dəyişimi) + server
- **HTML:** html (escape/unescape) + html/template (reusable, secure templates)
- **Cryptography:** MD5, SHA-*, TLS, DES, AES, HMAC + crypto-secure random
- **Data encoding:** UTF-8 daxili (UTF-8-in yaradıcıları Go-nu da yaratdı!); JSON/XML → obyekt çevirmə; interfeyslərlə genişlənən

### 4. Goroutine + channel — əsas konsept
**Kitabdan kod nümunəsi:**
```go
func count() {
    for i := 0; i < 5; i++ {
        fmt.Println(i)
        time.Sleep(time.Millisecond * 1)
    }
}

func main() {
    go count()                                   // goroutine başlat
    time.Sleep(time.Millisecond * 2)
    fmt.Println("Hello World")
    time.Sleep(time.Millisecond * 5)
}
// Output qarışıq: 0 1 Hello World 2 3 4
```

**Channel nümunəsi:**
```go
func printCount(c chan int) {
    num := 0
    for num >= 0 {
        num = <-c                 // kanaldan gözlə
        fmt.Print(num, " ")
    }
}

func main() {
    c := make(chan int)          // int kanalı yarat
    a := []int{8, 6, 7, 5, 3, 0, 9, -1}
    go printCount(c)
    for _, v := range a {
        c <- v                    // kanala göndər
    }
    time.Sleep(time.Millisecond * 1)
    fmt.Println("End of main")
}
// 8 6 7 5 3 0 9 -1 End of main
```
- `-1` → printCount loop-dan çıxır; main go routine bitənə qədər gözləyir
- Goroutine-lər = "lightweight threads" — runtime idarə edir, multi-core-də paralel
- Kanal → "internal micro-services communicating over a type-defined API"

### 5. Toolchain — go açar aləti
**Package management:**
```go
import (
    "golang.org/x/net/html"    // URL = xarici paket!
    "fmt"
    "net/http"
)
```
```bash
$ go get ./...       # bütün xarici asılılıqları çək
```
- Git, Mercurial, SVN, Bazaar dəstəyi (lokalda quraşdırılmış)
- Xarici paket standart paketlə EYNİ üsulla istifadə olunur

**Testing:**
```go
// hello_test.go — _test.go suffiksi MÜTLƏQ
func TestName(t *testing.T) {
    name := getName()
    if name != "World!" {
        t.Error("Respone from getName is unexpected value")
    }
}
```
```bash
$ go test            # PASS / FAIL
$ go test ./...      # subdirectory-lər də
```

**Code coverage:**
```bash
$ go test -cover
# Coverage: 33.3% of statements
```
+ HTML report (statement-level — if/else şaxələri ayrı-ayrı görünür).

**Formatting:**
- Effective Go (golang.org/doc/effective_go.html) — idiomatik üslub
- `go fmt` — kanonik formata salır; Vim/Sublime/Eclipse inteqrasiyası

### 6. Dil müqayisələri

**C vs Go:**
| | C | Go |
|---|---|---|
| Kompilyasiya | machine code | machine code |
| Runtime | YOX (özün idarə edirsən) | runtime + GC + thread idarəsi |
| Memory | manual | garbage collector |
| Concurrency | öz thread-lər | goroutine-lər |
| Kompilyasiya sürəti | yavaş (asılılıqlarla) | saniyələr (design goal!) |

- **cgo:** C kitabxanalarını bağlamaq; C-style ↔ Go string keçidi; SWIG dəstəyi; `go doc cgo`

**Java vs Go:**
- Java: runtime sistemdə + JVM/JIT; Go: **tək statik binary** — runtime daxilində
- Deploy: Go = 1 fayl; Java = runtime + app
- JIT vs compiled: aşkar qalib yoxdur (kontekstdən asılı)

**Python/PHP vs Go:**
- Dynamic vs static (dynamic-like type switching ilə)
- Python GIL → 1 thread vaxtında; PHP proses-başına → web server (Nginx/Apache) önə qoyulur
- **Go: daxili web server** — client DİREKT bağlanır; goroutine-lər connection-ları paralel idarə edir
- Python/PHP C-də yazılıb — C-yə keçid performans triki; Go-da fərq YOX (hamısı machine code)

**JavaScript/Node.js vs Go:**
- JS: **tək thread** (async I/O ayrı olsa da main thread bloklanır); Go: multi-thread, multi-core
- V8: JVM + JIT — uzun müddət işləyəndə optimallaşır; Go: statik machine code (JIT ehtiyacı yox)
- npm: mərkəzi repozitoriya (metadata + content); Go: mərkəz YOX — paketlər source-dan çəkilir

### 7. Quraşdırma + workspace
- Tour: tour.golang.org — brauzerdə icra edilə bilən dərslər
- Playground: play.golang.org — kod icra + paylaşım linki
- Windows/OSX: installer; OSX Homebrew: `brew install go`; Linux: apt-get/yum (köhnə) yaxud golang.org-dan son versiya
- Git + Mercurial: paket çəkilməsi üçün lazım

**GOPATH workspace:**
```
$GOPATH/
    src/                    ← source kod (səninki + asılılıqlar)
        github.com/
            Masterminds/
                cookoo/
                glide/
    bin/                    ← go install → executable-lər
        glide
    pkg/                    ← .a arxivləri (kitabxanalar)
        darwin_amd64/
            github.com/Masterminds/cookoo.a
```
```bash
$ mkdir $HOME/go
$ export GOPATH=$HOME/go
$ export PATH=$PATH:$GOPATH/bin     # bin-i PATH-ə əlavə et
# GOBIN — optional, alternativ binary yeri
```

### 8. Hello, Go — ilk web server
**Kitabdan kod nümunəsi:**
```go
package main

import (
    "fmt"
    "net/http"
)

func hello(res http.ResponseWriter, req *http.Request) {
    fmt.Fprint(res, "Hello, my name is Inigo Montoya")
}

func main() {
    http.HandleFunc("/", hello)
    http.ListenAndServe("localhost:4000", nil)
}
```

**Sub-kod izahı:**
- `package main` → executable üçün
- `hello(res, req)` → handler funksiyası — request/response obyektləri qəbul edir
- `http.HandleFunc("/", hello)` → "/" path-ini handler-a bağla
- `http.ListenAndServe("localhost:4000", nil)` → 4000 portunda server başlat
- `nil` handler → DefaultServeMux istifadə olunur

**İşə salma:**
```bash
$ go run inigo.go        # temp-də compile + icra (dev üçün)
$ go build inigo.go      # binary yarat
$ ./inigo                 # icra et
```

## Əsas terminlər
- Toolchain (dil + alətlər + ekosistem)
- Multiple Return Values / Named Returns
- Blank Identifier (`_`)
- Goroutine / Channel (lightweight thread / type-defined kommunikasiya)
- go get / go test / go fmt / go build / go run / go install
- Code Coverage (`-cover`)
- Effective Go (idiom üslubu)
- cgo (C binding) / SWIG
- GOPATH / GOBIN / src-bin-pkg workspace
- Ternary operatorun olmaması (sadəlik qərarı)
- JIT (V8, Java VM) vs Compiled Machine Code
- GIL (Python global interpreter lock)

## Praktik nəticə
- Çoxqayıdışlı funksiyalar + `_` — Go-nun əsas imzası; hər API-də (value, err) görəcəksən.
- `go get ./...` — xarici asılılıqları bir əmrlə çək; import path = URL.
- `go test -cover` — coverage də daxili; əlavə alət quraşdırmağa ehtiyac yoxdur.
- `go fmt` — müzakirəsiz vahid format; editor-a bağla.
- Web app = tək statik binary + daxili server — Java/PHP-nin runtime/nginx qatlarına ehtiyac yoxdur.
- Kanalı "tipləşdirilmiş micro-service API" kimi düşün — bu mental model bütün kitab boyu davam edir.

## Mənbə
Pages: 26-48 (PDF), book pages 3-25
