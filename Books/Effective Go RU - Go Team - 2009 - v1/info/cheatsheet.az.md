# Effective Go (RU) — Cheatsheet (Azərbaycanca)

> **Sənəd:** Effective Go (rus tərcüməsi) — Go Team, 2009 · 45 səh. · İdiomatik Go-nun kanonik bələdçisi.

---

## Formatlaşdırma & Adlar

```bash
gofmt                  # bütün format müzakirələrini bitirir
```

- İndent: TAB; sətir limiti YOX
- if/for/switch: mötərizə YOX, `{` eyni sətirdə (lekserin `;` qaydası!)
- Paket: tək lowercase söz; export təkrarsız (`ring.New`, `bufio.Reader`)
- Getter: `Owner()` (Get YOX); Setter: `SetOwner()`
- 1-metodlu interfeys: `-er` (Reader/Writer/Stringer) — kanonik imzaya uyğunluq ŞƏRT
- Çoxsözlü: MixedCaps (alt xətt YOX)

---

## İdarə konstruksiyaları

```go
// init ifadəsi + else buraxma:
if err := file.Chmod(0664); err != nil { return err }

// := təkrar mənimsətmə (eyni scope + yeni var da olmalı):
d, err := f.Stat()      // err YENİDƏN elan olunmur

// 3 for forması:
for i := 0; i < n; i++ {}
for cond {}
for {}

// Range string = rune parse (yanlış → U+FFFD):
for pos, char := range "日本\x80語" { ... }

// Paralel mənimsətmə (vergül operatoru YOX):
for i, j := 0, len(a)-1; i < j; i, j = i+1, j-1 { a[i], a[j] = a[j], a[i] }

// Qiymətsiz switch = if-else-if zənciri:
switch {
case '0' <= c && c <= '9': return c - '0'
case 'a' <= c && c <= 'f': return c - 'a' + 10
}

// Fallthrough YOX; çoxlu case: case ' ', '?', '&':
// Loop qırma: break Loop (label)

// Type switch — hər case öz tipi ilə eyni ad:
switch t := t.(type) {
case bool:  ... // t bool
case *bool: ... // t *bool
}
```

---

## Funksiyalar

```go
// Çoxlu qayıdış (C pointer-out əvəzi):
func nextInt(b []byte, i int) (int, int) { ... return x, i }

// Named nəticə + naked return:
func ReadFull(r Reader, buf []byte) (n int, err error) {
    for len(buf) > 0 && err == nil { ... }
    return
}

// Defer: açılışın yanında bağla — hər return qorunur:
f, err := os.Open(filename)
defer f.Close()

// Defer arqumentləri DƏRHAL; icra LIFO:
for i := 0; i < 5; i++ { defer fmt.Printf("%d ", i) }   // 4 3 2 1 0

// Trace pattern:
func a() { defer un(trace("a")); ... }
```

---

## new / make / massiv / slice / map

```go
// new(T) → *T (zero); make(T, args) → T (slice/map/channel ONLY)
p := new(SyncedBuffer)     // zero value DƏRHAL İŞLƏK (Buffer/Mutex patterni)
v := make([]int, 100)      // idiomatik
// new([]int) → nil-slice pointer — istifadə ETMƏ

// Composite literal — konstruktor əvəzi (lokalin adresi qanuni!):
return &File{fd: fd, name: name}
// new(File) == &File{}

// Massiv = VALUE: assign kopya; ölçü tipin hissəsi; C-lazım → pointer; idioma → SLICE

// Slice = ptr+len+cap value; append QAYITMALI:
s = append(s, elems...)
s = append(s, otherSlice...)    // ... ŞƏRT

// 2D: tək array + slicing:
picture := make([][]uint8, YSize)
pixels := make([]uint8, XSize*YSize)
for i := range picture {
    picture[i], pixels = pixels[:XSize], pixels[XSize:]
}

// Map comma-ok:
v, ok = m[k]          // mövcudluq
_, present := m[k]    // dəyərsiz yoxlama
delete(m, "PDT")      // təhlükəsiz silmə
```

---

## Çap formatları

```
%v      universal (array/slice/struct/map; map-lər SIRALANIR)
%+v     struct sahə adları ilə
%#v     tam Go sintaksisi
%T      tip
%q      string/rune literal
%d %x   tipdən oxuyur — C flag-ları YOX
```

```go
// Custom format + rekursiya tələsi:
func (t *T) String() string { return fmt.Sprintf("%d/%g/%q", t.a, t.b, t.c) }
// PİS: Sprintf("...%s", m)  (String→Sprintf→String...)
// YAXŞI: Sprintf("...%s", string(m))  /  %f (string deyil → təhlükəsiz)

// Variadic ötürmə: f(v...)
func Min(a ...int) int
```

---

## Konstantlar / init

```go
const (
    _  = iota                    // SKIP 0
    KB ByteSize = 1 << (10 * iota)
    MB; GB; TB; PB; EB; ZB; YB  // implicit təkrar
)

// var runtime ifadələr; init: import→var→init sırası
func init() { if user == "" { log.Fatal(...) } }
```

---

## Metodlar: pointer vs value

```go
type ByteSlice []byte

func (p *ByteSlice) Append(data []byte) { ... }        // caller-ı mutasiya edir
func (p *ByteSlice) Write(data []byte) (n int, err error) { ... }
// *ByteSlice → io.Writer ŞƏRTLƏNİR:
var b ByteSlice
fmt.Fprintf(&b, "This hour has %d days\n", 7)
```

| Receiver metodu | Çağırış |
|-----------------|---------|
| value | value + pointer (hər ikisində) |
| pointer | YALNIZ pointer (value ünvanlanabiləndə avtomatik `&`) |

---

## İnterfeyslər

```go
// Implicit təmin; conversion = metod dəsti dəyişməsi:
sort.IntSlice(s).Sort()
fmt.Sprint([]int(s))

// Type assertion (comma-ok TƏHLÜKƏSİZLİK ŞƏRT):
str, ok := value.(string)
// Yalnız switch-də çoxlu: case string: / case Stringer: (qarışdırma OK)

// Yalnız interfeys export et — konstruktor interfeys qaytarsın:
// crc32.NewIEEE & adler32.New → hər ikisi hash.Hash32

// HandlerFunc adapter — funksiya handler olur:
http.Handle("/args", http.HandlerFunc(ArgServer))

// 4 tip Handler: struct, int, channel, funksiya!
```

---

## `_` (boş identifikator)

```go
if _, err := os.Stat(path); os.IsNotExist(err) { ... }
_, present := timeZone[tz]

var _ = fmt.Printf            // dev-time istifadəsiz import
import _ "net/http/pprof"     // side-effect import

// KOMPAYL-TAİM interfeys yoxlaması (ƏSAS PATTERN):
var _ json.Marshaler = (*CustomData)(nil)

// runtime:
if _, ok := val.(json.Marshaler); ok { ... }
```

---

## Embedding

```go
// İnterfeys birləşməsi:
type ReadWriter interface {
    Reader
    Writer
}

// Struct — metod promote (forwarding YOX):
type Job struct {
    Command string
    *log.Logger        // Print/Printf/Println işlək!
}
job.Println("...")         // promote
job.Logger                 // daxili çıxış (tip adı = sahə adı)
```

- İrsdən FƏRQ: receiver = DAXİLİ tip
- Xarici ad daxilini GİZLƏDİR; eyni səviyyə təkrarı (istifadəsizdirsə) qanuni

---

## Konkurrentlik

```go
// Şüar: yaddaşı kommunikasiya ilə paylaş — CSP
go list.Sort()                     // goroutine (bitiş siqnalı YOX!)

c := make(chan int)                // unbuffered = SİNXRON
go func() { list.Sort(); c <- 1 }()
<-c                                // bitmə gözləməsi

// Semafor:
sem := make(chan int, MaxOutstanding)
sem <- 1; process(r); <-sem

// Worker pool (hər request-ə goroutine YOX):
for i := 0; i < MaxOutstanding; i++ { go handle(queue) }

// Kanal-kanal (mutexsiz RPC):
type Request struct {
    args       []int
    f          func([]int) int
    resultChan chan int    // klientin ÖZ kanalı
}
req.resultChan <- req.f(req.args)

// Paralelləşdirmə:
var numCPU = runtime.GOMAXPROCS(0)
c := make(chan int, numCPU)
for i := 0; i < numCPU; i++ { go v.DoSome(...) }
for i := 0; i < numCPU; i++ { <-c }

// Leaky buffer pool:
select {
case b = <-freeList:        // varsa al
default: b = new(Buffer)    // yoxsa yarat
}
select {
case freeList <- b:         // geri qaytar
default:                    // doludur → GC
}
```

---

## Xətalar / panic / recover

```go
// error interfeysi + zəngin model:
type PathError struct {
    Op, Path string; Err error
}
func (e *PathError) Error() string { return e.Op + " " + e.Path + ": " + e.Err.Error() }
// Prefiks konvensiyası: "image: unknown format"

// Xəta tipi yoxla + bərpa:
if e, ok := err.(*os.PathError); ok && e.Err == syscall.ENOSPC {
    deleteTempFiles(); continue
}

// panic: yalnız davam MÜMKÜNSÜZ olanda (kitabxanada son çarə)
// recover: yalnız defer daxilində — goroutine xilası:
func safelyDo(work *Work) {
    defer func() {
        if err := recover(); err != nil {
            log.Println("work failed:", err)
        }
    }()
    do(work)
}

// Paket-daxili parse patterni (named result + defer mutasiyası):
func Compile(str string) (regexp *Regexp, err error) {
    defer func() {
        if e := recover(); e != nil {
            regexp = nil
            err = e.(Error)    // öz tipi deyilsə re-panic
        }
    }()
    return regexp.doParse(str), nil
}
```

---

## Veb-server minimumu

```go
var templ = template.Must(template.New("qr").Parse(templateStr))   // başlanğıc panic'i OK

http.Handle("/", http.HandlerFunc(QR))
http.ListenAndServe(":1718", nil)

func QR(w http.ResponseWriter, req *http.Request) {
    templ.Execute(w, req.FormValue("s"))    // html/template avtomatik escapinq
}
```
