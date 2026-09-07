# The Go Workshop — Cheat Sheet (Azərbaycanca)

Packt Workshop Series 2019 — 19 fəsil, L2 Elementary

## 1. Dəyişənlər (Ch1)
```go
var x int = 10; var x = 10; x := 10      // := YALNIZ funksiyada
var ( A bool = false; B string = "x" )   // qrup bəyanı
const Pi = 3.14                          // dəyişməz
const ( Sun = iota; Mon; Tue )            // 0,1,2 avtomatik
// Pointer-lər:
var p *int          // nil
q := new(int)      // zero-lu yaddaş + pointer
r := &x            // mövcud dəyişəndən
*q                 // dereference — nil YOXSA PANİK
```

## 2. Kontrol (Ch2)
```go
if err := f(); err != nil { }         // initial if — GO-nun ƏSAS idiomu
switch x { case 1, 2: default: }      // fallthrough YOX
switch { case v > 0: }                 // ifadəsiz
for i := 0; i < n; i++ { }            // klassik
for k, v := range m { }                // map/slice; _ ilə udmaca
for { break }                          // sonsuz
// FizzBuzz: 15-i ƏVVƏL yoxla!
```

## 3. Tiqlər (Ch3)
```go
float64(i) / float64(k)               // implicit cast YOX
math.MaxInt64 + 1                      // WRAPAROUND! big.NewInt lazım
string(100)                            // "d" — İtoi istifadə et!
s[:10]                                 // MULTI-BYTE POZULA BİLƏR → []rune(s)[:10]
len(s)                                 // BAYT; len([]rune(s)) — SİMVOL
// Raw literal: `C:\Users` — escape-siz
```

## 4. Slice (Ch4)
```go
s := make([]int, 0, 10)               // pre-allocate
s = append(s, v); append(s, o...)     // nəticəni TƏYİN ET
t := s[a:b]                            // underlying PAYLAŞIR
t2 := append(s1[:0:0], s1...)         // MÜSTƏQİL kopya (ƏN effektiv)
v, ok := m[k]                          // map mövcudluq
delete(m, k)
// interface{} + type switch:
switch v := x.(type) { case int: case string: default: }
```

## 5. Funksiyalar (Ch5)
```go
func f() (a, b int) { return }          // naked return — SHADOWING RİSKİ
func f(vals ...int) {}                  // variadic; f(slice...)
inc := func() int { i++; return i }     // closure — i qorunmuş
defer f.Close()                         // FILO; dəyərlər DEFER ANINDA
type calc func(int, int) int            // funksiya TİPİ
```

## 6. Xətalar (Ch6)
```go
var ErrX = errors.New("x")              // Err prefiks, kiçik hərf
v, err := f(); if err != nil { }
defer func() { if r := recover(); r != nil { } }()   // recover YALNIZ defer-də
panic(ErrX)                             // error TİPİ ötür
```

## 7. İnterfeyslər (Ch7)
```go
type Speaker interface { Speak() string }
// implements YOXDUR — metodlar AVTOMATİK satisfy edir
func load(r io.Reader) (Person, error)   // ACCEPT INTERFACES
func String() string                      // Stringer — fmt çapı
i, ok := v.(int)                         // assertion (panic-siz)
```

## 8. JSON (Ch11)
```go
json.Unmarshal(data, &v)                // JSON → struct (POINTER!)
b, _ := json.Marshal(v)                  // struct → JSON
json.MarshalIndent(v, "", "    ")        // pretty
`json:"name,omitempty"`                  // boş → düşmür
`json:"-"`                               // gizli
// Naməlum struktur:
var v interface{}; json.Unmarshal(d, &v)
m := v.(map[string]interface{})          // rəqəmlər = float64!
```

## 9. Fayllar (Ch12)
```go
flag.String("n", "", "usage"); flag.Parse()
if *n == "" { flag.PrintDefaults(); os.Exit(1) }   // məcburi flag
f, _ := os.Create("f.txt")               // TRUNCATE!
defer f.Close()
f.WriteString("...")                     // və ya ioutil/os.WriteFile
os.OpenFile(f, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)   // APPEND
_, err := os.Stat(f); os.IsNotExist(err) // mövcudluq
// Signal:
sigs := make(chan os.Signal, 1)
signal.Notify(sigs, syscall.SIGINT)
go func() { s := <-sigs; cleanUp(); done <- true }()
```

## 10. SQL (Ch13)
```go
import _ "github.com/lib/pq"             // driver side-effect
db, _ := sql.Open("postgres", "user=... password=...")
db.Ping()                                 // real connectivity
defer db.Close()
stmt, _ := db.Prepare("INSERT ... VALUES ($1, $2)")   // INJECTION QORUNMASI
stmt.Exec(v1, v2)
rows, _ := db.Query("SELECT ...")
for rows.Next() { rows.Scan(&id, &name) }
db.Exec("DELETE ..."); res.RowsAffected()
```

## 11. HTTP Klient (Ch14)
```go
r, _ := http.Get(url); defer r.Body.Close()
data, _ := ioutil.ReadAll(r.Body)
// POST:
b, _ := json.Marshal(msg)
http.Post(url, "application/json", bytes.NewBuffer(b))
// Custom:
client := http.Client{Timeout: 11 * time.Second}
req, _ := http.NewRequest("POST", url, body)
req.Header.Set("Authorization", "token")
client.Do(req)
```

## 12. HTTP Server (Ch15)
```go
type H struct{ tpl *template.Template }
func (h H) ServeHTTP(w http.ResponseWriter, r *http.Request) { }
http.HandleFunc("/x", func(w, r) { })
http.Handle("/", hello{})                  // handler "/"-da; ListenAndServe(•, nil)
vl := r.URL.Query()                        // map[string][]string
if v, ok := vl["name"]; ok { }
// Template:
{{.Field}}  {{if .X}}...{{else}}...{{end}}  {{range .Items}}
tmpl, _ := template.ParseFiles("x.html")
tmpl.Execute(w, data)
// Statik:
http.ServeFile(w, r, "f.html")
http.Handle("/statics/", http.StripPrefix("/statics/",
    http.FileServer(http.Dir("./public"))))
// Form:
if r.Method == http.MethodPost { r.ParseForm(); r.Form.Get("name") }
```

## 13. Concurrency (Ch16)
```go
wg := &sync.WaitGroup{}
wg.Add(1)
go func() { defer wg.Done(); work() }()
wg.Wait()

atomic.AddInt32(&n, 1)                   // int32/64 only
mtx.Lock(); x++; mtx.Unlock()             // MINIMAL bölgü
ch := make(chan int); make(chan int, 10) // unbuffered/buffered
ch <- v; v := <-ch; close(ch)
for v := range ch { }                     // close-da dayanır
// Worker pool:
for i := 0; i < N; i++ { go worker(in, out) }
for i := from; i <= to; i++ { in <- i }
close(in)                                  // iş bitdi
for i := 0; i < N; i++ { total += <-out }
// Done channel:
out <- "done"; <-out                       // bitmə siqnalı
// Context:
cl, stop := context.WithCancel(ctx)
select { case <-cl.Done(): return; default: work() }
```

## 14. Alətlər (Ch17)
```bash
go build -o name src.go      # binary
go run main.go               # compile+icra
gofmt -w main.go             # format + YAZ
goimports -w main.go         # import əlavə/sil/çeşidle
go vet main.go               # Printf args, unmarshal pointer
go run --race main.go        # race detector
go doc -all                  # sənəd
go get github.com/gorilla/mux
go test ./...                # recursive test
```

## 15. Təhlükəsizlik (Ch18)
```go
// SQL: Prepare + $1 — concat HEÇ VAXT
// XSS: html/template (YOX text/template)
// Hash:
sha256.Sum256([]byte(s))     // %x ilə çap
// AES-GCM:
block, _ := aes.NewCipher([]byte(key))
gcm, _ := cipher.NewGCM(block)
nonce := make([]byte, gcm.NonceSize()); rand.Read(nonce)
ct := gcm.Seal(nonce, nonce, data, nil)   // [nonce|cipher]
gcm.Open(nil, nonce, ct[gcm.NonceSize():], nil)
// RSA:
priv, _ := rsa.GenerateKey(rand.Reader, 2048)
rsa.EncryptOAEP(sha256.New(), rand.Reader, &priv.PublicKey, msg, nil)
rsa.DecryptOAEP(sha256.New(), rand.Reader, priv, ct, nil)
// Parol:
bcrypt.GenerateFromPassword([]byte(pw), 10)
bcrypt.CompareHashAndPassword(hash, []byte(pw))
```

## 16. Build Constraints (Ch19)
```go
// +build linux            // boş sətir! / yeni: //go:build linux
// +build amd64,darwin 386,!gccgo
// Fayl adı: main_linux_amd64.go
```
```bash
GOOS=linux GOARCH=amd64 go build   # cross-compile
```

## Universal Qaydalar
1. err yoxlanılmadan KOD YOX (`if err != nil`)
2. defer: Close/Unlock həmişə
3. append nəticəsi TƏYİN ET
4. map sırası RANDOM; slice == müqayisə YOX
5. Multi-byte string → []rune
6. SQL input → placeholder; HTML → html/template
7. Concurrency → `-race`; counter → atomic/mutex
8. Security random → crypto/rand
9. Parol → bcrypt YOX plain SHA
10. Fərqli OS → build tag/suffix
