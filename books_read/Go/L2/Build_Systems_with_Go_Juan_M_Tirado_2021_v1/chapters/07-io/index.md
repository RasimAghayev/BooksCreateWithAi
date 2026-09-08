# Chapter 7 — Input/Output (səh. 143-156)

## Bu fəsil nədən bəhs edir?

Go-da I/O əməliyyatlarının əsası: `io.Reader`/`io.Writer` interfeysləri,
öz Reader/Writer implementasiyalarının yazılması, fayl oxuma/yazma
(`ioutil`, `os`), standart I/O (`os.Stdin/Stdout/Stderr`) və buferlənmiş
I/O (`bufio` — Reader, Scanner, Writer).

## Əsas fikirlər

### 1. Reader və Writer interfeysləri
**Nədir:** Go-da bütün I/O-nun təməli. Reader — oxunan datanın yazılacağı
byte massivi qəbul edir; Writer — yazılacaq massivi qəbul edir. Hər ikisi
(n say, err error) qaytarır.

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}
type Writer interface {
    Write(p []byte) (n int, err error)
}
```

### 2. Öz Reader implementasiyası
**Kitabdan kod nümunəsi (string-dən oxuyan Reader):**
```go
type MyReader struct {
    data string
    from int     // son oxu mövqeyi
}

func (r *MyReader) Read(p []byte) (int, error) {
    if p == nil {
        return -1, errors.New("nil target array")
    }
    if len(r.data) <= 0 || r.from == len(r.data) {
        return 0, io.EOF        // bütün simvollar bitib
    }
    n := len(r.data) - r.from
    if len(p) < n {
        n = len(p)              // massivə sığan qədər
    }
    for i := 0; i < n; i++ {
        b := byte(r.data[r.from])
        p[i] = b
        r.from++
    }
    if r.from == len(r.data) {
        return n, io.EOF        // sonuncu oxu — EOF ilə bildirilir
    }
    return n, nil
}

func main() {
    target := make([]byte, 5)
    mr := MyReader{"Save the world with Go!!!", 0}
    n, err := mr.Read(target)
    for err == nil {
        fmt.Printf("Read %d: Error: %v -> %s\n", n, err, target)
        n, err = mr.Read(target)
    }
}
```

**Sub-kod izahı:**
- `p` massivi **təkrar istifadə olunur** — qısa oxularıda köhnə simvollar
  qala bilər
- `io.EOF` → "daha data yoxdur" xüsusi error dəyəri; sonuncu oxu ilə birlikdə
  də qayıda bilər
- `from` sahəsi axın mövqeyini (stream position) izləyir

### 3. Öz Writer implementasiyası
**Kitabdan kod nümunəsi (batch-limited Writer):**
```go
type MyWriter struct {
    data string
    size int    // çağırış başına maksimum bayt
}

func (mw *MyWriter) Write(p []byte) (int, error) {
    if len(p) == 0 {
        return 0, io.EOF
    }
    n := mw.size
    var err error = nil
    if len(p) < mw.size {
        n = len(p)
    } else {
        err = errors.New("p larger than size")  // yazılan < len(p) → error
    }
    mw.data = mw.data + string(p[0:n])
    return n, err
}

func main() {
    msg := []byte("the world with Go!!!")
    mw := MyWriter{"Save ", 6}
    i := 0
    var err error
    for err == nil && i < len(msg) {
        n, err := mw.Write(msg[i:])
        fmt.Printf("Written %d error %v —> %s\n", n, err, mw.data)
        i = i + n
    }
}
```
- **Go qaydası:** yazılan bayt sayı `len(p)`-dən kiçikdirsə error dolu
  qaytarılmalıdır.

**Qeyd:** io paketində hazır interfeyslər çoxdur (ByteReader, ReadSeeker,
PipeReader...) — özünüzü yazmazdan əvvəl standart kitabxananı yoxlayın.

### 4. Fayl oxuma/yazma — ioutil (yüksək səviyyə)
**Kitabdan kod nümunəsi:**
```go
msg := "Save the world with Go!!!"
filePath := "/tmp/msg"

err := ioutil.WriteFile(filePath, []byte(msg), 0644)  // yaz
if err != nil {
    panic(err)
}

read, err := ioutil.ReadFile(filePath)                 // oxu → []byte
if err != nil {
    panic(err)
}
fmt.Printf("%s\n", read)
```

**Sub-kod izahı:**
- `ioutil.WriteFile(path, []byte, perm)` → bir çağırışda yazır; string →
  `[]byte` cast trivialdır
- `ioutil.ReadFile(path)` → faylın bütün məzmununu `[]byte` kimi qaytarır
- 0644 → Unix fayl icazələri
- Hər iki funksiya error qaytara bilər — nəzarət mütləqdir

### 5. Fayl əməliyyatları — os paketi (aşağı səviyyə)
**Nədir:** `os` — OS-dən asılı olmayan, amma daha aşağı səviyyəli interfeys;
 fayl deskriptorları, Seek və s.

**WriteString ilə sətir-sətir yazma:**
```go
filePath := "/tmp/msg"
msg := []string{"Rule", "the", "world", "with", "Go!!!"}

f, err := os.Create(filePath)   // fayl yarat → açıq deskriptor
if err != nil {
    panic(err)
}
defer f.Close()                  // deskriptoru burax

for _, s := range msg {
    f.WriteString(s + "\n")
}
```

**Seek ilə mövqe dəyişmə (faylın içini redaktə etmək):**
```go
tmp := os.TempDir()
file, err := os.Create(tmp + "/myfile")
if err != nil { panic(err) }
defer file.Close()

msg := "Save the world with Go!!!"
file.WriteString(msg)

positions := []int{4, 10, 20}
for _, i := range positions {
    file.Seek(int64(i), 0)     // offset-i mövqeyə keçir
    file.Write([]byte("X"))    // həmin baytı əvəz et
}

file.Seek(0, 0)                // ƏVVƏLƏ QAYT — vacib addım!
result := make([]byte, len(msg))
file.Read(result)
fmt.Printf("%s\n", result)     // SaveXthe wXrld with Xo!!!
```
- `Seek(offset, whence)` → oxu/yazı mövqeyini dəyişir; `0,0` → faylın
  əvvəlinə qaytarır (əks halda oxu son dəyişdilən yerdən davam edər)

### 6. Standart I/O (os.Stdin/Stdout/Stderr)
**Nədir:** Üç standart deskriptor; `*os.File` tipindədirlər — fayl
metodları işləyir.

**Stdout-a yazma:**
```go
msg := []byte("Save the world with Go!!!\n")
n, err := os.Stdout.Write(msg)     // fmt.Print ilə eyni nəticə
if err != nil { panic(err) }
fmt.Printf("Written %d characters\n", n)  // yazılan bayt sayı — 26 ("\n" də simvoldur)
```

**Stdin-dən oxuma (fixed buffer):**
```go
target := make([]byte, 50)
n, err := os.Stdin.Read(target)
if err != nil { panic(err) }
msg := string(target[:n])           // yalnız oxunan hissə!
fmt.Println(n, strings.ToUpper(msg))
```
- Enter basılanda EOF simvolu axını bitirir
- `target[:n]` → massivin yalnız dolu hissəsi; sabit ölçü həddi aşarsa
  qalan data itir → buferləmə lazımdır

### 7. bufio — buferlənmiş I/O
**Nədir:** Data həcmi bilinməyəndə/ifadəli olanda buferlənmiş reader/writer.

**NewReader + ReadString:**
```go
reader := bufio.NewReader(os.Stdin)
fmt.Print(">>> What do you have to say?\n")
fmt.Print("<<< ")
text, err := reader.ReadString('\n')   // delimiter-ə qədər oxu
if err != nil { panic(err) }
fmt.Println(">>> You're right!!!")
fmt.Println(strings.ToUpper(text))
```

**Scanner — axını hissələrə böl:**
```go
scanner := bufio.NewScanner(os.Stdin)
fmt.Println(">>> What do you have to say?\n")
counter := 0
for scanner.Scan() {              // sətir-sətir
    text := scanner.Text()
    counter = counter + len(text)
    if counter > 15 {
        break
    }
}
fmt.Println("that's enough")
```
- `scanner.Scan()` → növbəti element varmı (default: sətir)
- `scanner.Text()` → cari elementi string kimi verir
- Split funksiyası ilə delimiter xüsusi təyin oluna bilər

**NewWriter + Flush:**
```go
writer := bufio.NewWriter(os.Stdout)
msg := "Rule the world with Golang!!!"
for _, letter := range msg {
    time.Sleep(time.Millisecond * 300)
    writer.WriteByte(byte(letter))   // buferə yaz
    writer.Flush()                    // buferi məcburi çap et
}
```
- **Vacib:** buferdəki data avtomatik çıxmır — `Flush()` çağırılmalıdır
  (yazı makinası effekti üçün hər simvoldan sonra flush)

## Əsas terminlər
- Reader — `Read(p []byte) (n, err)` oxuma interfeysi
- Writer — `Write(p []byte) (n, err)` yazma interfeysi
- EOF (End of File) — axının sonu xətası
- File descriptor (fayl deskriptoru) — açıq faylın OS identifikatoru
- Seek — oxu/yazı mövqeyinin dəyişdirilməsi
- Buffered I/O (buferlənmiş I/O) — yaddaş buferi ilə I/O
- Flush — buferin məcburi boşaldılması

## Praktik nəticə
Bütün Go I/O-nu iki interfeys daşıyır: Reader və Writer. Fayl üçün sadə
yol `ioutil.WriteFile/ReadFile`, nəzarətli yol `os.Create` + `defer Close`
+ `Seek`. Konsol I/O `os.Stdin/Stdout` deskriptorlarıdır. Bilinməyən
həcmdə girişlər üçün `bufio.Reader`/`Scanner`, çıxış üçün `bufio.Writer`
+ `Flush` istifadə edin. Standart kitabxanada hazır interfeyslər var —
özlərini yazmazdan əvvəl yoxlayın.

## Mənbə
Pages: 143-156 (PDF 143-156)
