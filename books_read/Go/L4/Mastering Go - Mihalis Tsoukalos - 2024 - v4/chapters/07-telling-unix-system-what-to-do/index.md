# Chapter 7 — Telling a UNIX System What to Do (UNIX Sistemə Əmr Vermək)

## Bu chapter nədən bəhs edir?

stdin/stdout/stderr, UNIX prosesləri, io.Reader/io.Writer interfeyslərinin öz tiplərində
implementasiyası, buffered/unbuffered I/O (bufio), mətn faylının sətir/söz/simvol oxunuşu,
/dev/random, yazma üsulları, JSON (Marshal/Unmarshal, tag-lər, stream, pretty print),
viper (flag + konfiq faylı), cobra (komandalar, subkomandalar, alias-lar), go:embed,
os.ReadDir/DirEntry, io/fs, slog və statistika tətbiqinin cobra+JSON+ slog versiyası.

## Əsas fikirlər

### 1. stdin, stdout, stderr və UNIX Prosesləri
**Əsaslar:** UNIX hər şeyi fayl hesab edir; 3 standart deskriptor: 0=stdin, 1=stdout,
2=stderr (`/dev/stdin` və s.). Go-da: `os.Stdin`, `os.Stdout`, `os.Stderr` — portable və
təhlükəsiz yol.

**Proses kateqoriyaları:**
- User proseslər — user space, xüsusi hüquq yox
- Daemon proseslər — fondda, terminal tələb etmir
- Kernel prosesləri — kernel space, tam giriş

Go fork(2) dəstəkləmir — goroutine-lər runtime idarəli thread-lər üzərində qurulur.

### 2. io.Reader və io.Writer
```go
type Reader interface {
    Read(p []byte) (n int, err error)   // p buferini doldurur; n = oxunan bayt
}
type Writer interface {
    Write(p []byte) (n int, err error)   // p məzmununu yazır
}
```
**Kitabdan kod nümunəsi (öz tipində implementasiya):**
```go
// S2 üçün ənənəvi io.Reader — text sahəsi bufer kimi:
func (s *S2) Read(p []byte) (n int, err error) {
    if s.eof() {
        err = io.EOF              // gözlənilən axın sonu — XƏTA DEYİL
        return 0, err
    }
    l := len(p)
    if l > 0 {
        for n < l {
            p[n] = s.readByte()
            n++
            if s.eof() {
                s.text = s.text[0:0]
                break
            }
        }
    }
    return n, nil
}

// İstifadə: bufio ilə klassik şəkildə oxunur:
r := bufio.NewReader(&s2var)
for {
    n, err := r.Read(buf)
    if err == io.EOF {
        break
    } else if err != nil { /* xəta */ }
    fmt.Println(string(buf[:n]))
}
```
**Dərs:** bu interfeysləri öz tiplərində implement etmək sənə hər hansı writer-ı (log
servis, ehtiyat nüsxə, /dev/null) qəbul edən funksiyalara öz tipini ötürmək imkanı verir.

### 3. Buffered vs Unbuffered I/O
- **Unbuffered:** sistem çağırısı birbaşa — kritik data üçün daha etibarlı (güc kəsintisində
  buferdəki data itir)
- **Buffered (bufio):** sistem çağırışlarının sayını azaldır → performans; amma yazma
  gecikir (Flush tələb olunur)

### 4. Mətn Faylının Oxunması — 3 Səviyyə
**Sətir-sətir (ən vacib pattern):**
```go
func lineByLine(file string) error {
    f, err := os.Open(file)
    if err != nil { return err }
    defer f.Close()
    r := bufio.NewReader(f)
    for {
        line, err := r.ReadString('\n')     // '\n-ə qədər oxu
        if err == io.EOF {
            if len(line) != 0 {             // SON SƏTİR \n-siz ola bilər!
                fmt.Println(line)
            }
            break
        }
        if err != nil { return err }
        fmt.Print(line)                     // line özündə \n daşıyır
    }
    return nil
}
```

**Söz-söz:** regex ilə: `re := regexp.MustCompile("[^\\s]+")` → `re.FindAllString(line, -1)`.

**Simvol-simvol:** `for _, x := range line { fmt.Println(string(x)) }` — range rune qaytarır,
string(x) ilə simvola çevrilir.

### 5. /dev/random və Bayt Sırası (Endianness)
```go
f, _ := os.Open("/dev/random")
defer f.Close()
var seed int64
binary.Read(f, binary.LittleEndian, &seed)   // baytları int64-ə yığır
```
- Big endian: baytlar soldan-sağa (01|23|45|67); little endian: əksinə (67|45|23|01)
- /dev/random — random data mənbəyi (seed üçün)

### 6. Fayla Yazma — 4 Üsul + Append
**Kitabdan kod nümunəsi:**
```go
// 1) fmt.Fprintf — formatlı:
f1, _ := os.Create("/tmp/f1.txt")     // mövcuddursa TRUNCATE edir!
defer f1.Close()
fmt.Fprintf(f1, string(buffer))

// 2) WriteString:
n, _ := f2.WriteString(string(buffer))

// 3) bufio.Writer (Flush MÜTLƏQ!):
w := bufio.NewWriter(f3)
w.WriteString(string(buffer))
w.Flush()

// 4) io.WriteString:
io.WriteString(f4, string(buffer))

// Append:
f4, _ = os.OpenFile(f, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
n, _ = f4.Write([]byte("Put some more data at the end.\n"))   // Write = bayt slice
```
**Sub-kod izahı:** `os.Create` — truncate təhlükəsi var; append üçün OpenFile + O_APPEND.
Write() bayt slice tələb edir — binary data üçün düzgün yol; string convenience verir,
amma GC təzyiqi artırır.

### 7. JSON — encoding/json
**Marshal/Unmarshal:**
```go
type UseAll struct {
    Name    string `json:"username"`
    Surname string `json:"surname"`
    Year    int    `json:"created"`
}

// Struct → JSON:
t, err := json.Marshal(&useall)
// JSON → struct (string-i []byte-ə çevirmək MÜTLƏQ):
err = json.Unmarshal([]byte(str), &temp)
```
**Tag texnikaları:**
- `omitempty` — boş sahə JSON-a DÜŞMÜR
- `json:"-"` — sahə tamamilə GİZLƏDİLİR (parol kimi sensitiv data)

**1 nömrəli bug (kitabdan):** struct sahələri EXPORTED (böyük hərf) olmalıdır —
marshal/unmarshal problemində debug-a buradan başla!

**Stream emalı (çoxlu yazma):**
```go
func Serialize(e *json.Encoder, slice interface{}) error { return e.Encode(slice) }
func DeSerialize(e *json.Decoder, slice interface{}) error { return e.Decode(slice) }
// Encoder/Decoder birdən yarat, təkrar istifadə et — hər dəfə yaratmaq performans itirir
```

**Pretty print:**
```go
b, _ := json.MarshalIndent(v, "", "\t")        // tək record
encoder := json.NewEncoder(buffer)               // stream
encoder.SetIndent("", "\t")
encoder.Encode(data)
```

### 8. viper — Flag və Konfiqurasiya
**Nədir:** flag paketindən güclü alternativ; pflag üzərində; JSON/YAML/TOML/HCL konfiq
fayllarını ÖZÜ parse edir; env dəyişənləri də oxuya bilir.

**Kitabdan kod nümunəsi (flag):**
```go
pflag.StringP("name", "n", "Mike", "Name parameter")   // ad, qısa, default, izah
pflag.StringP("password", "p", "hardToGuess", "Password")
pflag.CommandLine.SetNormalizeFunc(aliasNormalizeFunc)  // --pass/--ps → --password alias
pflag.Parse()
viper.BindPFlags(pflag.CommandLine)                      // flag-ləri viper-ə bağla

name := viper.GetString("name")
viper.BindEnv("GOMAXPROCS")                               // env dəyişəni
val := viper.Get("GOMAXPROCS")
viper.Set("GOMAXPROCS", 16)                               // env DƏYİŞDİRİLİR
```

**JSON konfiq faylı:**
```go
viper.SetConfigType("json")
viper.SetConfigFile(CONFIG)
viper.ReadInConfig()          // DIQQƏT: fayl yoxdursa sessizcə boş kimi davam edir!

if viper.IsSet("macos") { fmt.Println(viper.Get("macos")) }
value := viper.GetBool("active")

// Struct-a Unmarshal — mapstructure tag ilə (json tag YOX!):
type ConfigStructure struct {
    MacPass string `mapstructure:"macos"`
}
var t ConfigStructure
viper.Unmarshal(&t)
```

### 9. cobra — Professional CLI
**Nədir:** docker/kubectl/hugo-nun əsası — komanda, subkomanda, alias dəstəyi.

**Workflow:**
```bash
go install github.com/spf13/cobra-cli@latest
cobra init                    # layihə strukturu (cmd/root.go + main.go)
cobra add one                 # cmd/one.go
cobra add list -p 'threeCmd'  # SUBKOMANDA - parent internal adı ilə
```

**Kitabdan kod nümunəsi:**
```go
// Global (persistent) flag — root.go init():
rootCmd.PersistentFlags().StringP("directory", "d", "/tmp", "Path")
rootCmd.PersistentFlags().Uint("depth", 2, "Depth of search")
viper.BindPFlag("directory", rootCmd.PersistentFlags().Lookup("directory"))

// Lokal flag — yalnız two komandası:
twoCmd.Flags().StringP("username", "u", "Mike", "Username")

// Alias:
var oneCmd = &cobra.Command{
    Use:     "one",
    Aliases: []string{"cmd1"},     // `cmd1` də işləyir
    Short:   "Command one",
}
```
**Sub-kod izahı:** PersistentFlags — bütün komandalara; Flags — təkcə öz komandasına.
`-p threeCmd` — subcommand-ı üç komandasına bağlayır (`two list` işləməz, `three list`
işləyəcək). Eyni alias 2 komandada olsa — yalnız İLKİ icra olunur (compiler xəbərdarlıq
etmir).

### 10. go:embed — Binary Daxilində Fayl (Go 1.16+)
**Kitabdan kod nümunəsi:**
```go
import _ "embed"                     // blank import

//go:embed static/image.png
var f1 []byte                         // binary fayl → []byte
//go:embed static/textfile
var f2 string                         // mətn fayl → string

// Self-printing proqram:
//go:embed printSource.go
var src string
func main() { fmt.Print(src) }        // öz mənbə kodunu çap edir!
```
**Qayda:** `//go:embed` sətri dərhal sonra gələn `var`-a tətbiq olunur; binary-lar üçün
[]byte, mətnlər üçün string.

### 11. os.ReadDir/DirEntry və io/fs (Go 1.16+)
**ioutil_DEPRECATED:** `os.ReadFile` ← ReadFile; `os.WriteFile` ← WriteFile;
`os.MkdirTemp` ← TempDir; `os.CreateTemp` ← TempFile; `os.ReadDir` ← ReadDir ([]DirEntry
qaytarır — []FileInfo YOX, daha sürətli).

**Kitabdan kod nümunəsi (rekursiv qovluq ölçüsü):**
```go
func GetSize(path string) (int64, error) {
    contents, err := os.ReadDir(path)
    if err != nil { return -1, err }
    var total int64
    for _, entry := range contents {
        if entry.IsDir() {
            temp, err := GetSize(filepath.Join(path, entry.Name()))   // REKURSİYA
            if err != nil { return -1, err }
            total += temp
        } else {
            info, err := entry.Info()     // FileInfo lazımdır → Info()
            if err != nil { return -1, err }
            total += info.Size()
        }
    }
    return total, nil
}
```

**io/fs — virtual fayl sistemi:** `fs.FS` read-only interface; `embed.FS` onu implement
edir → embedded fayllar fs funksiyaları ilə gəzilir:
```go
func list(f embed.FS) error {
    return fs.WalkDir(f, ".", walkFunction)   // hər entry üçün walkFunction çağırılır
}
s, err := fs.ReadFile(f, filepath)           // embedded faylı oxu
fileInfo, err := fs.Stat(f, path)
```
`filepath.Walk()` → `filepath.WalkDir()` keçidi performans verir (DirEntry).

### 12. slog — Strukturlaşdırılmış Logging (Go 1.21+)
**Kitabdan kod nümunəsi:**
```go
slog.Error("This is an ERROR message")   // default logger; Debug default-da ÇAP OLUNMUR

logLevel := &slog.LevelVar{}               // runtime dəyişən səviyyə
opts := &slog.HandlerOptions{Level: logLevel}
handler := slog.NewTextHandler(os.Stdout, opts)
logger := slog.New(handler)
logLevel.Set(slog.LevelDebug)              // səviyyəni yüksəlt
logger.Debug("This is a DEBUG message")

logJSON := slog.New(slog.NewJSONHandler(os.Stdout, nil))   // JSON log
logJSON.Error("ERROR message in JSON")
// {"time":"...","level":"ERROR","msg":"ERROR message in JSON"}
```
**JSON log-un dəyəri:** loglar time series DB-də saxlanılıb analiz/vizuallaşdırıla bilər.

**io.Discard hiyləsi (log söndürmə):**
```go
log.SetOutput(io.Discard)     // bütün yazılar uDulur — logger kodu dəyişmədən logging-off
```

### 13. Statistika Tətbiqi — cobra + JSON + slog Versiyası
**Struktur:** cobra init + list/delete/insert/search komandaları; `--log` persistent bool
flag (`BoolVarP(&disableLogging, "log", "l", false, ...)` — dəyişənə birbaşa bağlanır).

**Data modeli:** `index map[string]int` (filename→array index) + `data` slice; filename =
unikal identifikator. Insert əmrində mövcud fayl → əvvəlcə silinir (update emulyasiyası):
```go
_, ok := index[file]
if ok {
    delete(index, file)
    for i, k := range data {
        if k.Filename == file {
            data = slices.Delete(data, i, i+1)   // Go 1.21 slices paketi
            break
        }
    }
}
err := ProcessFile(file)
err = saveJSONFile(JSONFILE)       // Serialize → fayl
```
**Logger pattern:** hər komanda öz loggerini yaradır; `--log` false → io.Discard-a yönləndirir:
```go
logger = slog.New(slog.NewJSONHandler(os.Stderr, nil))
if disableLogging == false {
    logger = slog.New(slog.NewJSONHandler(io.Discard, nil))
}
slog.SetDefault(logger)
```
**List:** `sort.Sort(DFslice(data))` + PrettyPrintJSONstream.

## Əsas terminlər
- File Descriptor (fayl deskriptoru) — açıq faylın tam ədəd identifikatoru (0/1/2)
- Daemon (demon) — fonda işləyən, terminalsız proses
- io.Reader/io.Writer — oxuma/yazma interfeysləri (öz tipində implement olunur)
- bufio — buffered I/O sarğısı; Flush tələb edir
- Endianness (bayt sırası) — little/big endian repr
- Marshaling/Unmarshaling — struct↔JSON çevirməsi
- Struct Tag (strukt etiketi) — `json:"ad"` metadata; omitempty; `json:"-"`
- embed.FS — binary-ə daxil edilmiş virtual fayl sistemi
- DirEntry — ReadDir-in sürətli element tipi
- pflag/viper — müasir flag + konfiq idarəsi (BindPFlags)
- PersistentFlags vs Flags — qlobal vs lokal flag
- Alias (ləqəb) — komandanın alternativ adı
- slog — səviyyəli, handler-li (Text/JSON) strukturlaşdırılmış logging
- io.Discard — yazılan hər şeyi udan writer

## Praktik nətidə

(1) io.Reader/io.Writer hər yerdədir — onları öz tipində implement etmək səni bütün Go I/O
ekosistemə qoşur. (2) bufio ilə oxunun-da Flush-u unutma; buffered write-suz Flush = data
itir. (3) Sətir-sətir oxuma pattern-i: ReadString('\n') + io.EOF-də son \n-siz sətri emal et.
(4) os.Create truncate edir — append üçün OpenFile+O_APPEND. (5) JSON-un 1 bug mənbəyi:
exported sahə adları; tag-lərlə ad dəyiş, omitempty/`-` ilə nəzarət et. (6) Encoder/Decoder
bir dəfə yarat, təkrar istifadə et. (7) viper.ReadInConfig faylın mövcudluğunu SƏSSİCƏ
keçir — özün yoxla. (8) cobra: persistent=qlobal, local=komanda; subcommand -p ilə parent-ə
bağlanır. (9) go:embed binary-ə asset qatır — paylanan tək fayl. (10) slog JSON + io.Discard
= production logging: analiz üçün JSON, səssizlik üçün Discard. (11) filepath.Walk →
WalkDir keçidi performans qazancıdır.

## Mənbə
Pages: 265-318 (PDF 296-351)
