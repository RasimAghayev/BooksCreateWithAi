# Chapter 8 — Standard library

## Bu chapter nədən bəhs edir?

Go standart kitabxanasının strukturu və zəmanətləri, sonra ən çox istifadə olunan 3 paketin dərinləşməsi: **log** (standart və custom logger-lər), **encoding/json** (decode/encode), **io** (Writer/Reader interfeysləri və paketlərin bir-biri ilə işləməsi).

## Əsas fikirlər

### 1. Standart kitabxana — niyə etibarlı?
100+ paket, 38 kateqoriya: `archive, bufio, bytes, compress, container, crypto, database, encoding, errors, flag, fmt, hash, html, image, io, log, math, mime, net, os, path, reflect, regexp, runtime, sort, strconv, strings, sync, testing, text, time, unicode...`

**Zəmanətlər:** hər minor releazedə mövcuddur; backward-compatibility promise (geriuyğunluq vədi); Go-nun dev/build/release prosesinə daxildir; contributor-lər tərəfindən review olunur; hər releazedə test+benchmark edilir. 

**Mənbələr:** sənədlər http://golang.org/pkg/; interaktiv axtarış https://sourcegraph.com/; source kod `$GOROOT/src/pkg`-də; precompiled arxivlər (.a faylları) `$GOROOT/pkg`-də.

### 2. log paketi — əsas konfiqurasiya
**Nədir:** UNIX logging ənənəsinə əsaslanan paket: çıxış stdout-a, log stderr-ə (və ya əksinə — yalnız log yazan proqramlarda log stdout, error stderr).

**Kitabdan kod nümunəsi:**
```go
func init() {
    log.SetPrefix("TRACE: ")
    log.SetFlags(log.Ldate | log.Lmicroseconds | log.Llongfile)
}

func main() {
    // Println writes to the standard logger.
    log.Println("message")

    // Fatalln is Println() followed by a call to os.Exit(1).
    log.Fatalln("fatal message")

    // Panicln is Println() followed by a call to panic().
    log.Panicln("panic message")
}
```

**Sub-kod izahı:**
- `log.SetPrefix("TRACE: ")` → hər sətirə prefiks (identifikasiya üçün)
- `log.SetFlags(...)` → flag-lərin bit-müqayisəli birləşdirilməsi
- `log.Println` → normal yazı; `log.Fatalln` → yazıb `os.Exit(1)`; `log.Panicln` → yazıb `panic()`
- Hər funksiyanın `f` versiyası da var: `Printf`, `Fatalf`, `Panicf`
- **Logger-lər multigoroutine-safe-dir** — eyni anda bir çox goroutine yazıla bilər

**Flag-lər və iota — bit şift konstantları:**
```go
const (
    Ldate         = 1 << iota  // 1  = tarix: 2009/01/23
    Ltime                       // 2  = vaxt: 01:23:23
    Lmicroseconds               // 4  = mikrosaniyə
    Llongfile                   // 8  = tam fayl yolu + sətir
    Lshortfile                  // 16 = yalnız fayl adı + sətir (Llongfile-i override edir)
    LstdFlags  = Ldate | Ltime  // 3  = standart
)
```

**Sub-kod izahı (iota mexanizmi):**
- `1 << iota` → **bitwise left shift**: hər konstant öz unikal bit mövqeyini alır (1, 2, 4, 8, 16...)
- `iota` → konstant blokunda ifadəni təkrarlayır, hər addımda iota +1 (ilkin 0)
- `|` (pipe/OR) → bit-ləri birləşdirir: `Ldate | Ltime` = 3
- Nəticə: bir `int`-də çox seçeneği bit-bit saxlamaq — flag pattern-inin Go implementasiyası

### 3. Custom logger-lər (log.New)
**Nədir:** Fərqli log səviyyələri (Trace/Info/Warning/Error) üçün fərqli destinasiya + prefiks + flag.

**Kitabdan kod nümunəsi:**
```go
var (
    Trace   *log.Logger // Just about anything
    Info    *log.Logger // Important information
    Warning *log.Logger // Be concerned
    Error   *log.Logger // Critical problem
)

func init() {
    file, err := os.OpenFile("errors.txt",
        os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        log.Fatalln("Failed to open error log file:", err)
    }

    Trace = log.New(ioutil.Discard,
        "TRACE: ",
        log.Ldate|log.Ltime|log.Lshortfile)

    Info = log.New(os.Stdout,
        "INFO: ",
        log.Ldate|log.Ltime|log.Lshortfile)

    Warning = log.New(os.Stdout,
        "WARNING: ",
        log.Ldate|log.Ltime|log.Lshortfile)

    Error = log.New(io.MultiWriter(file, os.Stderr),
        "ERROR: ",
        log.Ldate|log.Ltime|log.Lshortfile)
}

func main() {
    Trace.Println("I have something standard to say")
    Info.Println("Special Information")
    Warning.Println("There is something you need to know about")
    Error.Println("Something has failed")
}
```

**Sub-kod izahı:**
- `log.New(out io.Writer, prefix string, flag int) *Logger` → destinasiya **io.Writer interfeysi** qəbul edir — fayl, stdout, hər hansı Writer
- `ioutil.Discard` → **hər şeyi udan Writer** (Write heç nə etmir) — Trace səviyyəsini "söndürmək" üçün texnika:
```go
// Discard is an io.Writer on which all Write calls succeed without doing anything.
var Discard io.Writer = devNull(0)
func (devNull) Write(p []byte) (int, error) { return len(p), nil }
```
- `io.MultiWriter(file, os.Stderr)` → variadic — bir yazını **eyni anda bir neçə** Writer-a yönləndirir
- `os.O_CREATE|os.O_WRONLY|os.O_APPEND` → fayl flag-lərinin birləşdirilməsi (bit OR)
- Logger metodları paket funksiyaları ilə eynidir: `Print/Printf/Println`, `Fatal*`, `Panic*` + `SetFlags/SetPrefix/Flags/Prefix`

**stdout/stderr unvanı:** `os.Stdout` və `os.Stderr` — `*File` tipləridir, `io.Writer` implement edir:
```go
var (
    Stdin  = NewFile(uintptr(syscall.Stdin), "/dev/stdin")
    Stdout = NewFile(uintptr(syscall.Stdout), "/dev/stdout")
    Stderr = NewFile(uintptr(syscall.Stderr), "/dev/stderr")
)
```

### 4. JSON Decoding — NewDecoder + Decode
**Nədir:** Stream (HTTP cavabı, fayl) üzərindən JSON-u struct-a çevirmək.

**Kitabdan kod nümunəsi:**
```go
type (
    // gResult maps to the result document received from the search.
    gResult struct {
        GsearchResultClass string `json:"GsearchResultClass"`
        UnescapedURL       string `json:"unescapedUrl"`
        URL                string `json:"url"`
        VisibleURL         string `json:"visibleUrl"`
        CacheURL           string `json:"cacheUrl"`
        Title              string `json:"title"`
        TitleNoFormatting  string `json:"titleNoFormatting"`
        Content            string `json:"content"`
    }

    // gResponse contains the top level document.
    gResponse struct {
        ResponseData struct {
            Results []gResult `json:"results"`
        } `json:"responseData"`
    }
)

func main() {
    uri := "http://ajax.googleapis.com/ajax/services/search/web?v=1.0&rsz=8&q=golang"

    resp, err := http.Get(uri)
    if err != nil {
        log.Println("ERROR:", err)
        return
    }
    defer resp.Body.Close()

    // Decode the JSON response into our struct type.
    var gr gResponse
    err = json.NewDecoder(resp.Body).Decode(&gr)
    if err != nil {
        log.Println("ERROR:", err)
        return
    }

    fmt.Println(gr)
}
```

**Sub-kod izahı:**
- `json.NewDecoder(r io.Reader) *Decoder` → Reader qəbul edir (resp.Body artıq Reader-dir!)
- `.Decode(v interface{}) error` → **empty interface** qəbul edir, reflection ilə tipi müəyyənləşdirib doldurur
- **Tag-lər** `` `json:"url"` `` → JSON açarını sahəyə map edir; tag yoxdursa case-insensitive ad uyğunluğu; tapılmasa sahə zero value qalır
- Anonim nested struct (`ResponseData struct {...}`) → bir dəfəlik daxili tiplər ayrıca elan etmədən
- Pointer-in ünvanını vermək (`Decode(&gr)`), hətta `var gr *gResponse; Decode(&gr)` — Decode nil pointer-in dəyərini özü yaradır

### 5. JSON Unmarshal — string/[]byte üçün
```go
var JSON = `{
    "name": "Gopher",
    "title": "programmer",
    "contact": {
        "home": "415.333.3333",
        "cell": "415.555.5555"
    }
}`

type Contact struct {
    Name    string `json:"name"`
    Title   string `json:"title"`
    Contact struct {
        Home string `json:"home"`
        Cell string `json:"cell"`
    } `json:"contact"`
}

var c Contact
err := json.Unmarshal([]byte(JSON), &c)
// Output: {Gopher programmer {415.333.3333 415.555.5555}}
```
- `json.Unmarshal([]byte, &v)` → string-i []byte-ə çevirib birbaşa; NewDecoder+Decode fərqli olaraq stream tələb etmir

**Struct bilmiriksə — map:**
```go
var c map[string]interface{}
err := json.Unmarshal([]byte(JSON), &c)

fmt.Println("Name:", c["name"])
// Nested çıxış üçün type assertion lazım olur — qeyri-əlverişli:
fmt.Println("H:", c["contact"].(map[string]interface{})["home"])
```
- `map[string]interface{}` → istənilən JSON strukturu qəbul edir, amma **type assertion** (`.(map[string]interface{})`) ilə manual naviqasiya tələb edir — az manipulyasiya lazım olanda sürətli üsul

### 6. JSON Encoding — MarshalIndent / Marshal
```go
c := make(map[string]interface{})
c["name"] = "Gopher"
c["title"] = "programmer"
c["contact"] = map[string]interface{}{
    "home": "415.333.3333",
    "cell": "415.555.5555",
}

// Marshal the map into a JSON string.
data, err := json.MarshalIndent(c, "", "    ")
if err != nil {
    log.Println("ERROR:", err)
    return
}

fmt.Println(string(data))
```

**Sub-kod izahı:**
- `json.MarshalIndent(v interface{}, prefix, indent string) ([]byte, error)` → pretty-print (formatlanmış, girintili) JSON
- `json.Marshal(v)` → sıxılmış JSON — network/API cavabları üçün
- Reflection ilə map/struct → JSON; []byte qaytarır → `string(data)` ilə çevir

### 7. io paketi — Writer/Reader interfeysləri
**Fəlsəfə:** UNIX-in "bir proqramın çıxışı — digərinin girişi" ideyasının Go tərcüməsi: `io.Writer` + `io.Reader` hər hansı data axını üçün vahid abstraksiya.

**İnterfeys elanları:**
```go
type Writer interface {
    Write(p []byte) (n int, err error)
}

type Reader interface {
    Read(p []byte) (n int, err error)
}
```

**Write qaydaları:** p-nin bütün uzunluğunu yazmağa çalış; mümkün olmasa **n < len(p) + non-nil error** qaytar; slice-ı heç vaxt modifikasiya etmə.

**Read qaydaları (4 qayda):**
1. `len(p)`-ə qədər oxu; az data varsa gözləmədən mövcudunu qaytar
2. EOF ilə birlikdə oxunan baytları qaytara bilər (n>0, err=EOF) və ya EOF-u sonrakı çağırışa saxlaya bilər (n>0, err=nil; növbəti: 0, EOF)
3. Çağırıcı **əvvəl n>0 baytları emal etməli, sonra error-a baxmalı**
4. `n=0, err=nil` qaytarmaq qadağandır — 0 bayt oxunuş hər zaman error ilə gəlməlidir

### 8. Paketlərin birgə işi (Hello World — 3 paket)
**Kitabdan kod nümunəsi:**
```go
package main

import (
    "bytes"
    "fmt"
    "os"
)

func main() {
    // Create a Buffer value and write a string to the buffer.
    var b bytes.Buffer
    b.Write([]byte("Hello "))

    // Use Fprintf to concatenate a string to the Buffer.
    fmt.Fprintf(&b, "World!")

    // Write the content of the Buffer to the stdout device.
    b.WriteTo(os.Stdout)
}
// Output: Hello World!
```

**Sub-kod izahı:**
- `bytes.Buffer` → özü io.Writer implement edir (`func (b *Buffer) Write(p []byte)`), həm də `WriteTo(w io.Writer)` — Reader rolu
- `fmt.Fprintf(w io.Writer, format, ...)` → ilk parametri Writer — Buffer də, File də qəbul edir
- `os.Stdout` → `*File`, onun `Write` metodu io.Writer implementasiyasıdır (`n != len(b)` olsa `io.ErrShortWrite` qaytarır)
- Nəticə: 3 paket bir-birini tanımadan interfeys vasitəsilə əməkdaşlıq edir — **interfeysin gücü**

### 9. Sadə curl — io.Copy + MultiWriter
**Kitabdan kod nümunəsi:**
```go
package main

import (
    "io"
    "log"
    "net/http"
    "os"
)

func main() {
    // r here is a response, and r.Body is an io.Reader.
    r, err := http.Get(os.Args[1])
    if err != nil {
        log.Fatalln(err)
    }

    // Create a file to persist the response.
    file, err := os.Create(os.Args[2])
    if err != nil {
        log.Fatalln(err)
    }
    defer file.Close()

    // Use MultiWriter so we can write to stdout and
    // a file on the same write operation.
    dest := io.MultiWriter(os.Stdout, file)

    // Read the response and write to both destinations.
    io.Copy(dest, r.Body)
    if err := r.Body.Close(); err != nil {
        log.Println(err)
    }
}
```

**Sub-kod izahı:**
- `http.Get(...)` → `r.Body` = `io.Reader`
- `os.Create(os.Args[2])` → fayl yarat (`*File` = Writer)
- `io.MultiWriter(os.Stdout, file)` → eyni yazı 2 destinasiyaya
- `io.Copy(dest, r.Body)` → Reader-dan Writer-a axını kopyalayır — chunk-chunk, bütün low-level iş io paketindədir
- URL → stdout + fayl eyni anda, ~20 sətirdə curl

## Əsas terminlər
- Standard Library (standart kitabxana)
- GOROOT (Go quraşdırma qovluğu)
- Archive File (.a — precompiled statik kitabxana)
- Logger / Log Level (jurnal / səviyyə)
- iota + Bit Shift (`1 << iota` — bit şifti ilə konstant zənciri)
- Flag Pattern (bit bayraqları)
- stdout / stderr (standart çıxış / xəta çıxışı)
- io.Writer / io.Reader (yazıcı / oxuyucu interfeysləri)
- MultiWriter (çəkyazıçı)
- Discard (uducu Writer)
- JSON Tag (JSON teqi)
- Marshal / Unmarshal (seriyalaşdırma / deseriyyalaşdırma)
- Reflection (refleksiya — runtime tip məlumatı)
- Type Assertion (tip müəyyənləşdirmə)
- EOF (fayl sonu)

## Praktik nəticə
- Log konfiqurasiyasını `init()`-də et — proqram başlanğıcında dərhal hazır olsun.
- Fərqli log səviyyələri üçün ayrıca `log.New` logger-ləri yarat; səviyyəni söndürmək üçün `ioutil.Discard` (müasir Go: `io.Discard`) istifadə et.
- HTTP cavabından JSON: `json.NewDecoder(resp.Body).Decode(&v)` — stringdən: `json.Unmarshal([]byte(s), &v)`.
- JSON strukturu məchuldursa `map[string]interface{}`, amma dərin naviqasiya üçün struct + tag-lər həmişə daha təmizdir.
- Öz tiplərinə `Read`/`Write` implement et — io paketindəki bütün hazır funksionallıq (Copy, MultiWriter və s.) pulsuz əlçatkan olur.
- Standart kitabxananın source kodunu oxu — idiomatik Go öyrənməyin ən yaxşı yolu.

## Mənbə
Pages: 205-231 (PDF), book pages 184-210
