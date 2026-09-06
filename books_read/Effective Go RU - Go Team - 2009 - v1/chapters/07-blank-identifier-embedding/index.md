# Chapter 7 — Boş identifikator və Embedding (səh. 32-36)

## Bu chapter nədən bəhs edir?

Boş identifikatorun (`_`) 5 istifadəsi (istifadəsiz dəyərlərin söndürülməsi, istifadəsiz import/var, side-effect import, interfeys kompayl-time yoxlaması) və Embedding — interfeys və strukturların daxilə tiplər yerləşdirməsi (metod promote, irs-dən fərqi, ad toqquşması qaydaları).

---

## Əsas fikirlər

### 1. Boş identifikator `_` — yazıla bilən /dev/null

Hər hansı dəyəri/Tipi təhlükəsiz İQNOR etmək üçün universal vasitə.

**a) Çoxlu mənimsətmədə lazımsız dəyər:**

```go
if _, err := os.Stat(path); os.IsNotExist(err) {
    fmt.Printf("%s не существует\n", path)
}
// PİS — xətanı iqnor: fi, _ := os.Stat(path) — mövcud olmayan path → panic!
```

**b) İstifadə olunmayan import/var (in development):**

```go
var _ = fmt.Printf    // debug; sonradan sil
var _ io.Reader       // debug; sonradan sil
_ = fd                 // TODO: istifadə et
```

Konvensiya: import-dan dərhal sonra + komment — təmizlənmə xatırlatması.

**c) Side-effect üçün import:**

```go
import _ "net/http/pprof"   // init HTTP handler-lər qeyd edir — API istifadə olunmur
```

**d) Interfeysin kompayl-time yoxlaması (ən vacib pattern):**

```go
var _ json.Marshaler = (*CustomData)(nil)
```

Bu bəyan: `*CustomData` → `Marshaler` təyin edilməsi KOMPİLYASİYA ZAMANI yoxlanılır. İnterfeys dəyişsə paket KOMPİLE OLMAZ — yeniləmə siqnalı. (json.RawMessage nümunəsi: xüsusi marshal var, amma statik çevirmə yoxdursa bu bəyan zəmanət verir.) Yalnız statik yoxlama olmayan yerlərdə istifadə.

**e) Runtime yoxlaması:**

```go
if _, ok := val.(json.Marshaler); ok { /* təmin edir */ }
```

### 2. Embedding — borrow edilmiş implementasiya

Go tipli MİRS yoxdur, amma **hissələri yerləşdirmək (embed)** olar.

**İnterfeyslərdə:**

```go
type Reader interface{ Read(p []byte) (n int, err error) }
type Writer interface{ Write(p []byte) (n int, err error) }

type ReadWriter interface {
    Reader    // embed — metodların BİRLƏŞMƏSİ
    Writer
}
// Sadə, aydın: ReadWriter = Read + Write
```

**Strukturlarda — metod promote:**

```go
type ReadWriter struct {
    *Reader    // *bufio.Reader embed
    *Writer    // *bufio.Writer embed
}
// Forwarding-metodsuz: bufio.ReadWriter AVTOMATİK hər üç interfeysi şərtlənir:
// io.Reader, io.Writer, io.ReadWriter
```

Adlı sahə ilə yazılsaydı, hər metod üçün forwarding yazmalı olardıq (`rw.reader.Read(p)`).

**İrsdən FƏRQİ — receiver:** embed metodun receiver-i DAXİLİ tipdir, xarici YOX. `bufio.ReadWriter.Read` çağrılanda `rw.reader` receiver-dir — forwarding-in dəqiq ekvivalentidir.

**Job nümunəsi — embed + adlı sahə birgə:**

```go
type Job struct {
    Command string
    *log.Logger    // embed — Print/Printf/Println promote
}

job.Println("начинаю выполнение...")           // log.Logger metodu!

func NewJob(command string, logger *log.Logger) *Job {
    return &Job{command, logger}
}

// Daxili çıxış sahəsi = tip adı (packaqet prefiksi siz):
job.Logger.Printf(...)    // birbaşa Logger-ə çıxış
func (job *Job) Printf(format string, args ...interface{}) {
    job.Logger.Printf("%q: %s", job.Command, fmt.Sprintf(format, args))
}
```

### 3. Ad toqquşması qaydaları

1. Xarici sahə/metod adı DİQQƏTLİ eyni adı GİZLƏDİR (Job.Command > log.Logger-in Command-i olsaydı).
2. Eyni səviyyədə təkrar ad = adətən XƏTA; amma proqramda heç vaxt istifadə olunmursa — qəbul ediləndir. Xarici tiplərin dəyişməsindən qoruyur.

---

## Əsas terminlər

| Termin | İzah |
|--------|------|
| `_` | Yazıla bilən /dev/null — dəyər/-tip iqnoru |
| `fi, _ := ...` | XƏTA iqnoru — pis praktika (panic riski) |
| `import _ "pkg"` | Side-effect import (pprof patterni) |
| `var _ Iface = (*T)(nil)` | Kompayl-time interfeys şərtlənmə yoxlaması |
| Embedding (interfeys) | İnterfeyslərin birləşməsi: ReadWriter |
| Embedding (struct) | Adsız sahə — metod promote; forwarding yox |
| Receiver = daxili tip | İrsdən fərq: embed metod daxilidə işləyir |
| Ad gizlətmə qaydası | Xarici ad daxilini örtür; eyni səviyyə təkrarı = xəta |

---

## Praktik nəticə

1. **`_` yalnız niyyətli iqnor üçün:** dəyər lazım deyilsə (key/value, err ilə yoxlama); xətanı ASLA `_` ilə atma.
2. **Side-effect importlar `_ "..."`** — pprof kimi init-qeydiyyat paketləri.
3. **`var _ Iface = (*T)(nil)` — public tiplər üçün sığorta:** interfeys dəyişəndə kompayl xətası dərhal xəbər verir; hər tip üçün YOX (statik yoxlama olmayan hallarda).
4. **İnterfeys birləşməsi embed ilə:** ReadWriter = Reader+Writer — metod sadalamağından təmiz.
5. **Struct embed = forwarding avtomatikası:** bufio.ReadWriter 3 interfeysi bir anda; adlı sahə istəyirsənsə forwarding-lər öz yaz.
6. **Receiver-i unutma:** embed metodu çağıran tipin ÖZ daxili sahəsidir — this-siz dünya.
7. **Toqquşmaq qaydaları sadədir:** xarici örtür; eyni səviyyə təkrarı istifadə olunmursa qanuni.

---

## Mənbə

- Sənəd: *Effective Go* (rus tərcüməsi), 2009
- PDF səhifələri: 32-36
