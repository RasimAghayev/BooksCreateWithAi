# Chapter 12 — Communicating with Non-Go Code (səh. 183-196)

## Bu fəsil nədən bəhs edir?

Go-dan kənar kodla işləmə: os/exec ilə xarici əmrlər (ping medianı, bc
kalkulyator prototipi — stdin/stdout pipe), cgo ilə C funksiyaları
(ioctl terminal eni, snowball stemmer) və "Cgo is not Go" qiyməti.

## Əsas fikirlər

### Recipe 61 — os/exec ilə ping
**Tapşırıq:** median ping vaxtı (average yox!) — ən sürətli server seçimi.

```go
func median(values []float64) float64 {
    nums := make([]float64, len(values))   // orijinalı dəyişmə!
    copy(nums, values)
    sort.Float64s(nums)
    i := len(nums) / 2
    if len(nums)%2 == 0 {
        return (nums[i-1] + nums[i]) / 2.0
    }
    return nums[i]
}

// findTime: "64 bytes ... time=142 ms" → 142.0
func findTime(line []byte) (float64, bool, error) {
    var prefix = []byte("time=")
    start := bytes.Index(line, prefix)
    if start == -1 {
        return 0, false, nil                  // bu sətirdə yoxdur
    }
    start += len(prefix)
    end := bytes.IndexByte(line[start:], ' ')
    if end == -1 {
        return 0, false, fmt.Errorf("can't find end")
    }
    end += start
    val, err := strconv.ParseFloat(string(line[start:end]), 64)
    if err != nil {
        return 0, false, err
    }
    return val, true, nil
}

// medianPing returns the median time of <count> pings to <host>.
func medianPingTime(host string, count int) (float64, error) {
    sw := "-c"
    if runtime.GOOS == "windows" {
        sw = "-n"                    // Windows ping fərqli flag!
    }
    cmd := exec.Command("ping", sw, fmt.Sprintf("%d", count), host)
    data, err := cmd.Output()        // stdout-u topla (prosesi gözlə)
    if err != nil {
        return 0, err
    }

    values := make([]float64, 0, count)
    s := bufio.NewScanner(bytes.NewReader(data))
    for s.Scan() {
        val, found, err := findTime(s.Bytes())
        if err != nil {
            return 0, err
        }
        if !found {
            continue
        }
        values = append(values, val)
    }
    if err := s.Err(); err != nil {
        return 0, err
    }
    return median(values), nil
}
```
- `cmd.Output()` → []byte; parse SİZİN işinizdir
- Timeout üçün `exec.CommandContext` (Recipe 50 ilə birgə)
- ICMP özün yazmaq mümkün — amma ping hər OS-də var

### Recipe 62 — cgo: terminal eni (ioctl)
**Tapşırıq:** çıxışı mərkəzləşdir — terminal eni lazımdır (Linux-only!).

```go
/*
#include <sys/ioctl.h>
#include <stdio.h>
#include <unistd.h>
int term_width() {
    struct winsize w;
    ioctl(STDOUT_FILENO, TIOCGWINSZ, &w);
    return w.ws_col;
}
*/
import "C"                       // C kodu şərh daxilində + import "C"

// İstifadə:
width := int(C.term_width())     // C.int → int çevrilməsi ŞƏRT

textLen := utf8.RuneCountInString(text)   // bayt YOX — rune sayı!
lpad := (width - textLen) / 2
if lpad <= 0 {
    return text
}
rpad := width - len(text) - lpad
return fmt.Sprintf("%*s%s%*s", lpad, " ", text, rpad, " ")
```
- `ldd` binary-i libc kimi shared library-lərə bağlayır → target maşında
  eyni library-lər olmalı
- C compiler tələbi → cross-compilation çətinləşir

### Recipe 63 — bc ilə kalkulyator (stdin/stdout pipe)
**Tapşırıq:** konfiqurasiyada riyazi ifadələr: `maxSize = 5 * (2^20)`.

```go
// Calc is a calculator.
type Calc struct {
    p *os.Process
    w io.Writer
    r *bufio.Reader
}

// NewCalc creates a new calculator
func NewCalc() (*Calc, error) {
    // -l adds mathlib, -q mean no banner
    cmd := exec.Command("bc", "-lq")

    w, err := cmd.StdinPipe()       // prosesin stdin-i
    if err != nil {
        return nil, err
    }
    r, err := cmd.StdoutPipe()       // stdout
    if err != nil {
        return nil, err
    }
    cmd.Stderr = cmd.Stdout          // stderr → stdout-a yönəlt

    if err := cmd.Start(); err != nil {
        return nil, err
    }
    return &Calc{p: cmd.Process, w: w, r: bufio.NewReader(r)}, nil
}

// Eval evaluates a math expression such as "3 / 7".
func (c *Calc) Eval(expr string) (float64, error) {
    if _, err := fmt.Fprintf(c.w, "%s\n", expr); err != nil {
        return 0, err
    }
    line, err := c.r.ReadString('\n')
    if err != nil {
        return 0, err
    }
    line = line[:len(line)-1]
    return strconv.ParseFloat(line, 64)
}

// Close closes the calculator.
func (c *Calc) Close() error {
    return c.p.Kill()                // prosesi ÖLDÜR — şərt!
}
```
- Pipe-lar Reader/Writer → standart funksiyalar işləyir
- Bir proses dəfələrlə istifadə olunur (hər sorğu üçün yenisi YOX)
- Proses həddi var → Close mütləq (defer ilə)

### Recipe 64 — cgo: snowball stemmer
**Tapşırıq:** works/working/worked → work. Pure Go stemmer yoxdur.

```go
/*
#include <libstemmer.h>
#include <stdlib.h>
#cgo LDFLAGS: -lstemmer        // link əmri!
*/
import "C"

// Stemmer stems words for a specific language.
type Stemmer struct {
    st *C.struct_sb_stemmer      // C struct → C.struct_*
}

// NewStemmer creates a new stemmer for a language.
func NewStemmer(lang string) (*Stemmer, error) {
    cLang := C.CString(lang)     // Go → C string (MƏMORİ AYRILIR!)
    st := C.sb_stemmer_new(cLang, nil)
    C.free(unsafe.Pointer(cLang))   // azad et — ŞƏRT
    if st == nil {
        return nil, fmt.Errorf("can't create stemmer for %q", lang)
    }
    return &Stemmer{st}, nil
}

// Close closes the stemmer, freeing allocated memory.
func (s *Stemmer) Close() {
    if s.st != nil {
        C.sb_stemmer_delete(s.st)    // C tərəfindəki yaddaş
        s.st = nil
    }
}

// Stem will stem a word.
func (s *Stemmer) Stem(word string) string {
    cWord := C.CBytes([]byte(word))       // Go bayt → C buffer
    size := C.int(len(word))
    sym := C.sb_stemmer_stem(s.st, (*C.uchar)(cWord), size)
    if sym == nil {
        return ""
    }
    i := C.sb_stemmer_length(s.st)
    data := C.GoBytes(unsafe.Pointer(sym), i)   // C → Go bayt
    return string(data)                           // stemmer idarə edir — free yox
}
```
- Çevrilmələr: `C.CString` (free lazım!), `C.CBytes`, `C.GoBytes`,
  `C.int` ↔ `int`
- **Kim ayırırsa, o azad edir** — C yaddaşını əl ilə idarə edirsiniz
- `#cgo LDFLAGS/pkg-config` — build inteqrasiyası

## Final Thoughts-dən

**"Cgo is not Go" (Rob Pike):** itkilər — cross-platform, sürətli build,
avtomatik memory management; + C compiler və C biliyi. Mümkünsə "pure
Go" paket axtarın. os/exec təhlükəsizdir, amma executable hədəf
maşında olmalıdır. C yaddaşsında şübhən varsa — KOPYALA. (`i int` vs
`int i` xatirələri :) )

## Əsas terminlər
- os/exec — xarici əmr icrası
- exec.Command / cmd.Output / cmd.Start
- CommandContext — kontekstli (ləğv olunan) əmr
- StdinPipe / StdoutPipe — proses I/O körpüləri
- cgo — C inteqrasiyası (`import "C"`)
- #cgo LDFLAGS / pkg-config — link direktivləri
- C.CString / C.CBytes / C.GoBytes — tip çevrilmələri
- unsafe.Pointer — C yaddaşına düz çıxış
- Kim ayırır, azad edir — C yaddaş qaydası
- Stemming — sözün kök forması
- "Cgo is not Go" — texniki borc xəbərdarlığı
- Median vs average — median outlier-lərə davamlıdır

## Praktik nəticə
Xarici kod seçimi: (1) pure Go paket; (2) os/exec — əmrlə çıxışı parse
et (Output + Scanner); uzunömürlü proses üçün pipe-lar + Kill; (3) cgo
— son çarə: struct pointer, CString/CBytes çevrilmələri, free zəmanəti,
LDFLAGS. Hər halda risk/deyer balansını yoxlayın — cgo operational
mürəkkəblik gətirir.

## Mənbə
Pages: 183-196 (PDF 183-196)
