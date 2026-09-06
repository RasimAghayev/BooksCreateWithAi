# Go in Action — Cheat Sheet (bütün kitabdan toplanmış)

William Kennedy, Brian Ketelsen, Erik St. Martin — Manning, 2016

## Go əmrləri (CLI)

### `go build <pkg|file|.>`
**Nə edir:** Paketi/faylı kompilyasiya edib executable yaradır. Fayl adı verməsən cari paket götürülür; `.` = cari qovluq.

### `go build github.com/.../chapter3/...`
**Nə edir:** `...` wildcard — qovluq altındakı bütün paketləri build edir.

### `go clean <file>`
**Nə edir:** Build nəticəsində yaranan executable-ı silir.

### `go run <file.go>`
**Nə edir:** Build + execute bir addımda — binary diskdə qalmır.

### `go vet <file>`
**Nə edir:** Ümumi xətaları yoxlayır: Printf tipli çağırışlarda yanlış parametrlər, metod imza xətaları, yanlış struct tag-lər, açarsız composite literal-lər.
Nümunə: `main.go:6: no formatting directive in Printf call`

### `go fmt <file|pkg>`
**Nə edir:** Kodu standart formata salır (`if err != nil { return err }` → çoxsətirli).

### `go doc tar`
**Nə edir:** Paket sənədlərini birbaşa terminalda göstərir.

### `godoc -http=:6060`
**Nə edir:** localhost:6060-da browsable sənədləşmə server-i qaldırır.

### `go get <url>`
**Nə edir:** Remote repodan paketi çəkib GOPATH-ə qoyur; **rekursivdir** — asılılıqları da gətirir.

### `go test -v`
**Nə edir:** Testləri verbose output ilə işə salır.

### `go test -run="none" -bench=. -benchmem`
**Nə edir:** Unit testləri regex ilə süzür ("none" = heç biri), bütün benchmark-ları + allokasiya statistikası ilə işə salır.

### `go build -race && ./example`
**Nə edir:** Race detector ilə build — işə salınca data race-ləri WARNING: DATA RACE ilə göstərir.

### `runtime.GOMAXPROCS(runtime.NumCPU())`
**Nə edir:** Hər fiziki core üçün 1 logical processor ayır — paralel icra.

## Sintaktik konstruktlar

### Slice yaratma
```go
slice := make([]int, 5)          // len=cap=5
slice := make([]int, 3, 5)       // len=3, cap=5
slice := []int{10, 20, 30}       // literal
slice := []string{99: ""}        // 100 elementli
var slice []int                   // nil slice
slice := make([]int, 0)           // empty slice
```

### Slicing formulları
```
slice[i:j]   → len = j-i,  cap = k-i
slice[i:j:k] → len = j-i,  cap = k-i  (cap məhdudlaşdırma)
source[2:3:3] → detach pattern: append yeni array yaradır
```

### append
```go
newSlice = append(slice, 60)      // 1 element
combined := append(s1, s2...)     // slice birləşdirmə
```

### Map əməliyyatları
```go
dict := make(map[string]int)
dict := map[string]string{"Red": "#da1337"}   // literal
value, exists := colors["Blue"]                 // 2 dəyərli lookup (TÖVSİYƏ)
value := colors["Blue"]                        // zero-value yoxlaması
delete(colors, "Coral")                        // silmə
```

### Channel əməliyyatları
```go
unbuffered := make(chan int)         // sinxron mübadilə qarantiyası
buffered := make(chan string, 10)    // müstəqil send/receive
ch <- v                             // send
v := <-ch                           // receive
v, ok := <-ch                       // + "açıqdırmı?" flag
close(ch)                           // send olmaz, receive boşalana qədər olur
for v := range ch                   // channel bağlanana qədər
select { case ...: default: }       // çoxşaxəli + nonblocking peek
```

### sync alətləri
```go
var wg sync.WaitGroup
wg.Add(2)          // sayğac qoy
wg.Done()          // azalt (defer ilə!)
wg.Wait()          // sıfırlanana qədər blokla

mutex.Lock()       // kritik bölmə başla
mutex.Unlock()     // bitir

atomic.AddInt64(&counter, 1)   // atomik artırma
atomic.LoadInt64(&shutdown)    // atomik oxu
atomic.StoreInt64(&shutdown, 1) // atomik yaz
```

### OS siqnalları (Runner pattern)
```go
interrupt := make(chan os.Signal, 1)
signal.Notify(interrupt, os.Interrupt)   // Ctrl+C → kanala
signal.Stop(interrupt)                    // abunəliyi dayandır
timeout := time.After(d)                  // vaxt bitən kanal
```

## JSON / XML
```go
// Struct + tag
type Feed struct {
    Name string `json:"site"`
    URI  string `json:"link"`
}
type item struct {
    Title string `xml:"title"`
}

// Decode (stream)
json.NewDecoder(resp.Body).Decode(&gr)
xml.NewDecoder(resp.Body).Decode(&document)

// Unmarshal (string)
json.Unmarshal([]byte(s), &c)

// Encode
data, _ := json.MarshalIndent(v, "", "    ")   // pretty
data, _ := json.Marshal(v)                     // sıxılmış

// Encode → Writer
json.NewEncoder(rw).Encode(&u)
```

## io əməliyyatları
```go
io.Copy(dest, src)                      // Reader → Writer axını
io.MultiWriter(os.Stdout, file)         // eyni yazı → NWriter
fmt.Fprintf(&buffer, "format %s", v)     // Writer-a formatla yaz
buffer.WriteTo(os.Stdout)               // Buffer → Writer
ioutil.Discard                          // hər şeyi udan Writer
os.Open(name) / os.Create(name)         // fayl aç / yarat
```

## Log
```go
log.SetPrefix("TRACE: ")
log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Llongfile | log.Lshortfile)
logger := log.New(out io.Writer, prefix string, flag int) *Logger
log.Fatalln("...")   // yaz + os.Exit(1)
log.Panicln("...")   // yaz + panic()
```

## Test
```go
func TestXxx(t *testing.T) {
    t.Log("...")      // informasiya (-v ilə)
    t.Fatal("...")    // FAIL + dayandır
    t.Error("...")    // FAIL + davam
}
func BenchmarkXxx(b *testing.B) {
    b.ResetTimer()
    for i := 0; i < b.N; i++ { /* ölçülən kod */ }
}
func ExampleXxx() {
    // ...
    // Output:
    // gözlənilən çıxış
}

// Mock server
server := httptest.NewServer(http.HandlerFunc(handlerFunc))
defer server.Close()
resp, _ := http.Get(server.URL)

// Endpoint test
req, _ := http.NewRequest("GET", "/path", nil)
rw := httptest.NewRecorder()
http.DefaultServeMux.ServeHTTP(rw, req)
// rw.Code, rw.Body yoxla
```

## Əsas qaydalar (kitabdan toplanmış)
1. Zero value üçün `var`, ilkin dəyər üçün `:=`
2. Pointer receiver default; value receiver yalnız dəyişməz tiplər üçün
3. Pointer receiver ilə implement interfeysə yalnız pointer uyğun gəlir
4. Loop-dakı dəyişənləri goroutine-ə parametr kimi ötür — closure tələsi!
5. Resurs açanda `defer Close` dərhal yanına
6. Xəta olanda digər return dəyərlərinə etibar etmə
7. İnterfeysləri kiçik saxla (1-2 metod, `-er` şəkilçisi)
8. Unexported tip + `New` factory = gizli implementasiya, açıq API
9. Embedming promotion — iç tipin kimliyi həmişə qalır
10. Channel: data-nı locklama — ötür (CSP)
