# Chapter 1 — Reading and Writing (I/O) (səh. 12-30)

## Bu fəsil nədən bəhs edir?

`io.Reader`/`io.Writer` abstraksiyası ətrafında 7 resept: []byte↔Reader
çevrilməsi (bytes.NewReader), gzip ilə log sıxıştırma (janitor utility),
bytes.Buffer ilə SQL generasiyası, şərtli dekompressiya, öz io.Writer
implementasiyası (Benford qanunu), os.Pipe ilə dinamik data və memory
mapped fayllarda (mmap) axtarış.

## Əsas fikirlər

### Recipe 1 — []byte → io.Reader (bytes.NewReader)
**Tapşırıq:** io.Reader gözləyən API-yə (encoding/gob, json.NewDecoder)
[]byte ötürmək.

**Kitabdan kod nümunəsi:**
```go
// UnmarshalRide returns a Ride from serialized data.
func UnmarshalRide(data []byte, ride *Ride) error {
    r := bytes.NewReader(data)      // []byte → io.Reader
    return NewDecoder(r).DecodeRide(ride)
}
```
- `bytes.NewReader`, `strings.NewReader`, `bytes.Buffer` — yaddaş içi
  Reader/Writer-lər; test üçün də çox əlverişlidir

### Recipe 2 — gzip ilə köhnə log fayllarının sıxıştırılması
**Tapşırıq:** 30 gündən köhnə .log fayllarını gzip-lə, SHA1 müqayisəsi ilə
təsdiq et, orijinalı sil.

```go
func gzCompress(src, dest string) error {
    file, err := os.Open(src)
    if err != nil {
        return err
    }
    defer file.Close()

    out, err := os.Create(dest)
    if err != nil {
        return err
    }
    defer out.Close()

    w := gzip.NewWriter(out)
    defer w.Close()

    // Metadata — io.Copy-dan ƏVVƏL olmalıdır
    w.Name = src
    info, err := file.Stat()
    if err == nil {
        w.ModTime = info.ModTime()
    }

    if _, err := io.Copy(w, file); err != nil {
        os.Remove(dest)
        return err
    }
    return nil
}
```

**Fayl seçimi (FS abstraksiyası ilə):**
```go
func shouldCompress(path string, maxAge time.Duration) bool {
    info, err := os.Stat(path)
    if err != nil {
        log.Printf("warning: %q: can't get info: %s", path, err)
        return false
    }
    if info.IsDir() {
        return false
    }
    return time.Since(info.ModTime()) >= maxAge
}

func filesToCompress(dir string, maxAge time.Duration) ([]string, error) {
    root := os.DirFS(dir)                    // qovluq → fs.FS
    logFiles, err := fs.Glob(root, "*.log")  // pattern ilə seç
    // ...
}
```

**SHA1 imzası (hash = io.Writer!):**
```go
func fileSHA1(fileName string) (string, error) {
    file, err := os.Open(fileName)
    if err != nil {
        return "", nil
    }
    defer file.Close()
    var r io.Reader = file
    if path.Ext(fileName) == ".gz" {         // .gz → gzip-dən oxu
        var err error
        r, err = gzip.NewReader(r)
        if err != nil {
            return "", err
        }
    }
    w := sha1.New()                          // sha1 io.Writer-dir
    if _, err := io.Copy(w, r); err != nil { // fayl → hash
        return "", err
    }
    sig := fmt.Sprintf("%x", w.Sum(nil))
    return sig, nil
}
```
- Axın: sıxıştır → müqayisə et → eynidirsə orijinalı sil
- `*os.File` io.Reader/Writer/Closer-in hamısını implement edir

### Recipe 3 — bytes.Buffer ilə SQL generasiyası
**Tapşırıq:** SQL injection riskinə görə istifadəçi SQL yazmağa izni
yoxdur — cədvəl+sütunlardan SELECT qur.

```go
func genSelect(table string, columns []string) (string, error) {
    var buf bytes.Buffer              // buf = io.Writer
    if len(columns) == 0 {
        return "", fmt.Errorf("empty select")
    }
    fmt.Fprintln(&buf, "SELECT")
    for i, col := range columns {
        suffix := ","
        if i == len(columns)-1 {
            suffix = "" // son sütunda vergül yoxdur
        }
        fmt.Fprintf(&buf, "\t%s%s\n", col, suffix)
    }
    fmt.Fprintf(&buf, "FROM %s;", table)
    return buf.String(), nil
}
```
- Üstünlük: fmt paketi + generasiya alqoritmi sadələşir; data gələ-gələ
  yazılır, yığılmır; text/template daha mürəkkəb hallar üçün

### Recipe 4 — şərtli dekompressiya
**Tapşırıq:** qarışıq (.log + .log.gz) fayllarında HTTP redirect (3XX)
sətirlərini say.

```go
// numRedirects io.Reader qəbul edir — mənbədən asılı DEYİL
func numRedirects(r io.Reader) (int, int, error) {
    s := bufio.NewScanner(r)
    nLines, nRedirects := 0, 0
    for s.Scan() {
        nLines++
        // 203.252.212.44 - - [...] "GET /ksc.html HTTP/1.0" 200 7280
        fields := strings.Fields(s.Text())
        code := fields[len(fields)-2]       // status sondan bir əvvəl
        if code[0] == '3' {                  // 3XX = redirect
            nRedirects++
        }
    }
    if err := s.Err(); err != nil {
        return -1, -1, err
    }
    return nLines, nRedirects, nil
}

// İstifadə:
matches, _ := filepath.Glob("logs/http-*.log*")
for _, fileName := range matches {
    file, err := os.Open(fileName)
    // ...
    var r io.Reader = file          // = ilə! := yeni dəyişən yaradar
    if strings.HasSuffix(fileName, ".gz") {
        r, err = gzip.NewReader(r)  // mövcud r-ni yenilə
        // ...
    }
    nl, nr, err := numRedirects(r)  // hər iki halda eyni funksiya
    // ...
}
```
- **Vacib incəlik:** `var r io.Reader = file` → `=` işarəsi; `:=` if
  blokunda YENİ r yaradır, xaricdəki dəyişməz

### Recipe 5 — öz io.Writer implementasiyası (Benford qanunu)
**Tapşırıq:** rəqəmlərin ilk rəqəminin tezliyini hesabla (maliyə fırıldağının
yoxlanması).

```go
type DigitsFreq struct {
    Freqs map[rune]int // ilk rəqəm tezliyi
    inNum bool         // hazırda rəqəm daxilindəyikmi
}

// Write implements io.Writer.
func (d *DigitsFreq) Write(data []byte) (int, error) {
    if d.Freqs == nil {
        d.Freqs = make(map[rune]int)
    }
    for _, b := range data {
        if r := rune(b); unicode.IsDigit(r) {
            if !d.inNum {       // rəqəm başlayır → İLK rəqəm
                d.Freqs[r]++
                d.inNum = true
            }
            continue
        }
        if d.inNum {            // rəqəm bitdi
            d.inNum = false
        }
    }
    return len(data), nil
}

// İstifadə — io.Copy ilə istənilən Reader-dən:
var df DigitsFreq
io.Copy(&df, strings.NewReader(data))
for r, c := range df.Freqs {
    fmt.Printf("%c →%d\n", r, c)
}
```
- **"Make the zero value useful"** (Rob Pike): `var df DigitsFreq` — sıfır
  dəyərlə istifadə oluna bilər (New lazım deyil); Write içində map yaradılır
- crypto paketindəki bütün hash tipləri (sha1 və s.) eyni yanaşma

### Recipe 6 — os.Pipe ilə dinamik data generasiyası
**Tapşırıq:** DB sorğusu channel qaytarır; HTTP handler bu datanı JSON kimi
axıtmalıdır.

```go
// encodeRides kanaldan oxuyub WriteCloser-a JSON yazır
func encodeRides(ch <-chan Ride, w io.WriteCloser) error {
    enc := json.NewEncoder(w)
    defer w.Close() // funksiya çıxanda "data bitdi" siqnalı
    for r := range ch {
        if err := enc.Encode(r); err != nil {
            return err
        }
    }
    return nil
}

func queryHandler(w http.ResponseWriter, r *http.Request) {
    // 1) body oxu (hədd qoyulmuş):
    data, err := io.ReadAll(io.LimitReader(r.Body, maxSize))
    if err != nil {
        http.Error(w, "can't read body", http.StatusBadRequest)
        return
    }
    location := string(data)

    // 2) DB sorğusu → kanal:
    conn, err := Dial(dbDSN)
    if err != nil {
        http.Error(w, "can't connect", http.StatusInternalServerError)
        return
    }
    ch := conn.QueryRidesIn(location)

    // 3) Pipe yarat — Reader/Writer cütü:
    rp, wp, err := os.Pipe()
    if err != nil {
        http.Error(w, "can't create pipe", http.StatusInternalServerError)
        return
    }

    // 4) Goroutine-də kodla, əsas axında kopyala:
    go encodeRides(ch, wp)
    _, err = io.Copy(w, rp)
    if err != nil {
        log.Printf("error: can't encode: %s", err)
    }
}
```
- Nəticə "JSON lines" formatıdır (sətir-sətir JSON) — Decoder həzm edir,
  digər klientlər bilməz
- **Vacib:** yazı tərəfi bağlanmayırsa, oxu tərəfi asılı qalar
- Unix fəlsəfəsi: pipe ilə proqramları birləşdir (Knuth vs. McIlroy hekayəsi)

### Recipe 7 — memory mapped faylda axtarış (mmap)
**Tapşırıq:** böyük qarışıq faylda (binary+text+JSON) `loc:{...}` fraqmentləri
tapmaq — fayl yaddaşa sığmır.

```go
type Location struct {
    Lat float64
    Lng float64
}

file, err := os.Open("data.txt")
if err != nil {
    log.Fatalf("error: %s", err)
}
defer file.Close()

// mmap → fayl məzmunu []byte kimi (Unix-only: golang.org/x/sys/unix):
fi, err := file.Stat()
m, err := unix.Mmap(
    int(file.Fd()), 0, int(fi.Size()),
    unix.PROT_READ, unix.MAP_PRIVATE,
)
defer unix.Munmap(m)   // bitəndə azad et

// []byte üzərində axtarış:
pos := 0
locPrefix := []byte("loc:{")
var loc Location
for {
    i := bytes.Index(m[pos:], locPrefix)      // "loc:" tap
    if i == -1 {
        break
    }
    i += len(locPrefix) - 1                    // prefiksdən keç
    start := pos + i
    size := bytes.IndexByte(m[start:], '}')    // bağlanan } tap
    if size == -1 {
        break
    }
    size++
    if err := json.Unmarshal(m[start:start+size], &loc); err != nil {
        log.Fatalf("error: %s", err)
    }
    fmt.Printf("%+v\n", loc)
    pos = start + size + 1                      // növbəti axtarışa keç
}
```
- mmap: fayl yaddaşa "xəritələnir" — fayl ölçüsündən asılı olmayacaq
  sadə kod; Windows üçün ayrıca kod/build tags lazımdır
- "grep üçün kifayətdirsə, sən üçün də kifayətdir"

## Əsas terminlər
- io.Reader/io.Writer — I/O-nun vahid interfeysləri
- bytes.NewReader / strings.NewReader — []byte/string → Reader
- bytes.Buffer — yaddaş içi Writer (kod generasiyası)
- gzip.Writer/Reader — sıxıştırma qatı
- fs.FS / os.DirFS — fayl sistemi abstraksiyası
- sha1.New — hash tipi Writer kimi
- os.Pipe — Unix pipe (Reader/Writer cütü)
- mmap — memory mapped file (fayl → []byte)
- JSON lines — sətir-sətir JSON formatı
- Zero value useful — sıfır dəyərlə işlək tip (Rob Pike)

## Praktik nəticə
Bayt-istiqamətli I/O-da konkret tip (`*os.File`, `net.Conn`) yerinə
`io.Reader`/`io.Writer` qəbul edin — funksiya hər mənbəyə (fayl, gzip,
socket, string) uyğunlaşar. Reader/Writer-ləri qat-qat bükürəm (gzip →
hash → copy) — hər qat bir iş görür. Öz tiplərinizdə Write/Read
implement edin + zero value işlək olsun. Böyük fayllar üçün mmap, iki
axın arasında körpü üçün os.Pipe (yazı tərəfini bağlamağı unutmayın).

## Mənbə
Pages: 12-30 (PDF 12-30)
