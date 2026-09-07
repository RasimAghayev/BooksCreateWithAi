# Chapter 1 — Getting started with Go (Go ilə Başlanğıc)

## Bu chapter nədən bəhs edir?

Go dilinin yaranma tarixi, dizayn fəlsəfəsi, əsas xüsusiyyətləri (çoxlu qaytarma,
standart kitabxana, goroutine/channel, toolchain), digər dillərlə müqayisə (C, Rust, Zig,
Nim, Java, Python/PHP/JS), quraşdırma, workspace, environment dəyişənləri və ilk web
server tətbiqi.

## Əsas fikirlər

### 1. Go Nədir?
**Nədir:** 2007-də Google-də Robert Griesemer, Rob Pike, Ken Thompson tərəfindən
yaradılmış, statik tipli, kompilyasiya olunan açıq mənbəli dil. 2009-da elan edilib.

**Fəlsəfə:** Nəzəri təmizlik YOX — real praktik situasiyalar. C, Pascal, Java, Python-un
ən yaxşı cəhətləri; 3 yaradıcının HAMISININ razılığı olmadan heç bir feature daxil olmur
→ sadə, amma güclü dil. Generics 10+ il müzakirədən sonra yalnız backward-compatibility
pozulmadan əlavə edildi.

**3 qat (layers of Go):**
1. Proqramlaşdırma dili (müasir hardware üçün)
2. Toolchain — test, sənədləndirmə, formatlaşdırma, paket idarəetməsi built-in
3. Ekosistem — Git əsaslı paket sistemi üzərində böyüyən kitabxanalar

### 2. Multiple Return Values (Çoxlu qaytarma)
**Nədir:** Funksiya birdən çox dəyər qaytara bilir — tuple/hash-ə hopdurmadan.

**Kitabdan kod nümunəsi:**
```go
func getStrings() (string, string) {
    return "Foo", "Bar"
}
func main() {
    n1, n2 := getStrings()   // hər dəyərə dəyişən
    n3, _ := getStrings()     // ikincisini iqnor et
}
```

**Named return (naked return):**
```go
func getStrings() (first string, second string) {
    first = "Foo"
    second = "Bar"
    return            // adlı dəyərlər avtomatik qayıdır
}
```
Qeyd: naked return qısa funksiyalarda oxunaqlıdır; uzun funksiyalarda return-ün izlənməsi
çətinləşdiyindən getdikcə az istifadə olunur. Konvensiya: `error` həmişə SON qaytarma;
`a, b, err := f()` pattern-i bütün Go-da standartdır.

### 3. Modern Standart Kitabxana
**Necə işləyir:** Şəbəkə, kriptoqrafiya, serializasiya, riyaziyyat out-of-the-box.

**Şəbəkə nümunələri:**
```go
// TCP birbaşa:
conn, err := net.Dial("tcp", "golang.org:80")
fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
status, err := bufio.NewReader(conn).ReadString('\n')

// HTTP client:
resp, err := http.Get("http://example.com/")
defer resp.Body.Close()          // defer — blok sonunda təmizləmə
body, err := io.ReadAll(resp.Body)
```

**Sahələr:** html/template (təhlükəsiz HTML şablonları); crypto (MD5/SHA, TLS, DES, AES,
HMAC, crypto/rand); encoding (JSON↔struct; Go daxildə hər şey UTF-8 — UTF-8-in
yaradıcıları Go-nu da yaradıb!).

### 4. Goroutine və Channel
**Goroutine:** `go` keyword-ilə başladılan funksiya; Go runtime tərəfindən OS thread-lərə
map edilir; block olarsa (I/O) başqa thread-ə keçirilir; çox nüvəli paralel icra.

**Kitabdan kod nümunəsi:**
```go
func count() {
    for i := 0; i < 5; i++ {
        fmt.Println(i)
        time.Sleep(time.Millisecond * 5)
    }
}
func main() {
    go count()          // count və main paralel işə düşür
    time.Sleep(20 * time.Millisecond)
    fmt.Println("Hello World")
}
```

**Channel:** Goroutine-lərarası ünsiyyət; default bloklayıcı (sinxronizasiya):
```go
c := make(chan int)
go printCount(c)        // printCount <-c ilə gözləyir
for _, v := range []int{8, 6, 7, 5, 3, 0, 9, -1} {
    c <- v              // hər yaz printCount-dək bloklanır
}
```
Channel = typed API üzərindən danışan yüngül daxili servis. Kitab boyu serverlərdə,
mesaj keçirmədə, WaitGroup-larla nəzarətdə istifadə olunacaq.

### 5. Toolchain — Paket İdarəetməsi
```bash
go mod init github.com/USER/REPO   # go.mod yaradır
go get ./...                         # import-lardakı xarici paketləri endirir
go mod tidy                          # asılılıq ağacını təmizləyir
go install PKG                       # /bin üçün binary
```
- Import qruplaşdırılır, əlifba sırası; standard kitabxana əvvəl, xarici paketlər sonra
  (goimports avtomatik edir)
- Xarici paketlər URL ilə göstərilir; Git ilə birbaşa işləyir (private repo da)
- Konvensiya: executable-lar `./cmd/NAME/main.go` altında; eyni qovluqda 2 main YOX

### 6. Testing və Coverage
```go
func TestName(t *testing.T) {
    got := reverseNameFixed("William")
    if !strings.EqualFold(got, "mailliW") {
        t.Errorf("got [%s], expected [%s]", got, "mailliW")
    }
}
```
- `_test.go` suffiksi — test yalnız `go test`-də işləyir, build-də YOX
- `go test ./...` — paket + alt qovluqlar
- `go test -cover` — statement səviyyəli coverage (%)
- Kitabın bug dərsi: reverse funksiyası bayt üzrə işləyəndə multibyte simvolları korlayır
  → `[]rune(name)` ilə fix — rune = Unicode code point!

### 7. Dil Landşaftı — Müqayisələr
- **C:** Go runtime+GC verir (memory leak/buffer overflow/race riski azalır); C-də thread
  və yaddaş əl ilə; Go compile sürəti dəygandır; cgo ilə C bağlama (SWIG dəstəyi).
- **Rust/Zig/Nim:** Rust (2015, Mozilla) borrow checker — daha mürəkkəb, daha dəqiq
  lifetime nəzarəti; Linux kernel Rust qəbul edir. Zig C-yə bənzər stil, GC yox; Nim
  Pythonvari sintaksis. Go ən yetkin + backward-compatibility zəmanəti.
- **Java:** Go tək binary (runtime daxil) — Java JRE quraşdırılmış sistem tələb edir;
  Java VM+JIT vs Go native binary — performans müqayisəsi qeyri-müəyyən; GOGC="off"
  ilə GC hətta söndürülə bilər (nadirdir!).
- **Python/PHP/JS:** dinamik tipli, paket menecerləri (pip/composer/npm) sonradan
  əlavə edilib (təhlükəsizlik exploitləri tarixçəsi); Go-nun built-in web server-i
  birbaşa qoşulur (proxy tələb etmir); Node.js tək-thread event loop vs Go çox-thread
  çoxnüvəli goroutine modeli — Go hardware-dən daha çox istifadə edir.

### 8. İşə Başlama
- Quraşdırma: go.dev/doc/install; macOS: `brew install go`; Linux: paket menecerləri
  (köhnə ola bilər — rəsmi endirmə üstün)
- `$GOPATH` — workspace baza qovluğu; modul izolyasiyası layihəni istənilən yerə qoymağa
  imkan verir
- **GOOS/GOARCH** — cross-compile hədəf OS/arxitektura (Ch12 containerlərdə)
- Öyrənmə alətləri: go.dev/tour (brauzerdə icra), Go Playground (paylaşılan linklər),
  Go by Example
- AI alətləri: kitab boyu prompt engineering ilə test data/simulyasiya yaranacaq

### 9. Hello, Go — İlk Web Server
**Kitabdan kod nümunəsi:**
```bash
go mod init hellogo
```
```go
package main
import (
    "fmt"
    "net/http"
)
func hello(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "Hello, my name is Inigo Montoya")
}
func main() {
    http.HandleFunc("/", hello)
    http.ListenAndServe("localhost:4000", nil)
}
```
`go run inigo.go` (temp binary + dərhal icra) və ya `go build inigo.go && ./inigo`
(`-o bin/inigo` ilə çıxış yolu).

## Əsas terminlər
- Goroutine — runtime idarəli paralel funksiya
- Channel (kanal) — goroutine-lərarası tipəmən ünsiyyət
- Toolchain — dil + test/format/paket alətləri dəsti
- Naked return (boş qaytarma) — adlı return dəyişənlərinin dəyərsiz `return`-u
- go mod — modul/paket idarəetmə sistemi
- GOOS/GOARCH — cross-compile hədəf dəyişənləri
- Rune — Unicode code point (multibyte simvol dəstəyi)

## Praktik nəticə

Yeni Go layihəsinin standart axını: `go mod init` → kod (executable üçün ./cmd/NAME/)
→ goimports ilə import təmizliyi → go test -cover → go build. Server yazarkən proxy
qatına ehtiyac yoxdur — net/http birbaşa qoşulur. String emalında multibyte riski varsa
həmişə rune üzrə işlə. Dil seçimində Go-nun yeri: servis/sistem/cloud; embedded üçün
TinyGo istisna olmaqla ilk seçim deyil.

## Mənbə
Pages: 3-28 (PDF 24-49)
