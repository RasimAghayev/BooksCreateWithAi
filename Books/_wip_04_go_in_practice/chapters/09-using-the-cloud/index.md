# Chapter 9 — Using the cloud (Technique 56-61)

## Bu chapter nədən bəhs edir?

Cloud computing əsasları (IaaS/PaaS/SaaS, container-lər, cloud-native), vendor lock-in-dən qaçınmaq (interfeys pattern), divergent error handling (package error-ləri), host detection, dependency check, cross-compiling (GOOS/GOARCH, gox, build tags), runtime monitoring.

## Əsas fikirlər

### 1. Cloud computing növləri
**İdarəetmə spektri:**
| Model | Sən idarə edirsən | Servis idarə edir | Nümunə |
|---|---|---|---|
| Traditional | HƏR ŞEY (rack, server, OS...) | — | Öz server otağın |
| **IaaS** | OS, runtime, app, data | server, storage, network | AWS, Azure, GCE |
| **PaaS** | app, data | OS, runtime + infra | Heroku, Cloud Foundry, Deis |
| **SaaS** | — | hər şey | Salesforce, Office 365, Stripe |

- **IaaS fərqi:** server aylarla YOX — API ilə yarat/məhv et; REST API + CLI
- **PaaS:** kod push → platform build+run; horizontal scale API ilə
- **SaaS:** DB-as-a-service kimi — qurma/konfiqurasiya/scale vendor-da

### 2. Container-lər vs Virtual Machines
| | VM | Container |
|---|---|---|
| İzolyasiya | Hypervisor + **Guest OS (öz kernel)** | Host kernel paylaşır |
| Başlanğıc | Yavaş (kernel boot) | **Dərhal** |
| İçindəkilər | Bütün guest OS replikasiya olunur | Yalnız bins/libs |
| Sıxlıq | Ağır | Daha sıx yerləşmə |
| Binary fərqi | — | App A: Debian libs, App B: CentOS libs eyni hostda! |

**Cloud-native:** proqramlaşdırıla bilən cloud-dan istifadə — avtomatik scale, **remediation** (avtomatik problem düzəltmə), **microservices** (kiçik müstəqil proseslər, API kommunikasiyası — Ch 10).

### TECHNIQUE 56: Multi-provider — interfeys pattern
**Problem:** hər vendor-un öz SDK/API — birinə bağlanan kod LOCK-IN.

**Həll (printer driver modeli):** interfeys təyin et → hər vendor üçün implementasiya.

```go
// 1. İNTERFEYS — vendor neytral:
type File interface {
    Load(string) (io.ReadCloser, error)
    Save(string, io.ReadSeeker) error
}

// 2. LOKAL İMPLEMENTASİYA (dev/test üçün ən asan):
type LocalFile struct {
    Base string
}

func (l LocalFile) Load(path string) (io.ReadCloser, error) {
    p := filepath.Join(l.Base, path)
    return os.Open(p)                       // *os.File = io.ReadCloser
}

func (l LocalFile) Save(path string, body io.ReadSeeker) error {
    p := filepath.Join(l.Base, path)
    d := filepath.Dir(p)
    err := os.MkdirAll(d, os.ModeDir|os.ModePerm)   // qovluq yarat
    if err != nil {
        return err
    }
    f, err := os.Create(p)
    if err != nil {
        return err
    }
    defer f.Close()
    _, err = io.Copy(f, body)
    return err
}

// 3. FABRİKA — konfiqurasiyadan seç:
func fileStore() (File, error) {
    return &LocalFile{Base: "."}, nil        // AWS/GCE/Azure implementasiyası dəyişə bilər
}

// 4. İSTİFADƏ — interfeysə bağlı, vendor-a YOX:
body := bytes.NewReader([]byte(content))
store, _ := fileStore()
store.Save("foo/bar", body)
c, _ := store.Load("foo/bar")
o, _ := ioutil.ReadAll(c)
```
- AWS S3 / GCS / Azure Blob — hər biri File interfeysini implement edir
- Dev/test LOKAL filesystem → cloud-a keçid = factory dəyişmək
- Credentials dev/test/prod FƏRLİ olmalı

### TECHNIQUE 57: Divergent error-lər
**Problem:** hər implementasiya öz error tipini qaytarırsa → app kodu implementasiya DETALLARINI bilməli → interfeys pozulur.

**Həll:** paket səviyyəli error DƏYİŞƏNLƏRİ + implementasiya daxilində xəritələmə:

```go
var (
    ErrFileNotFound   = errors.New("File not found")
    ErrCannotLoadFile = errors.New("Unable to load file")
    ErrCannotSaveFile = errors.New("Unable to save file")
)

func (l LocalFile) Load(path string) (io.ReadCloser, error) {
    p := filepath.Join(l.Base, path)
    var oerr error
    o, err := os.Open(p)
    if err != nil && os.IsNotExist(err) {        // os xətasını TANı
        log.Printf("Unable to find %s", path)   // ORİJİNAL xətanı LOG-la (itirmə!)
        oerr = ErrFileNotFound                   // paket error-u qaytar
    } else if err != nil {
        log.Printf("Error loading file %s, err: %s", path, err)
        oerr = ErrCannotLoadFile
    }
    return o, oerr
}

// Caller — implementation bilmədən:
if err == ErrFileNotFound {
    fmt.Println("Cannot find the file")
}
```
**Dərs:** original error LOG-da (monitoring alert üçün), caller-a PAKET error-u (Ch 4 error dəyişənləri pattern).

### TECHNIQUE 58: Host information
```go
// os paketi:
os.Hostname()              // kernel hostname
os.Getpid()                // proses ID
os.Getwd()                 // iş qovluğu
os.PathSeparator           // '/' (Windows '\')
os.PathListSeparator       // ':' (Windows ';')

// IP tapma (interfeyslər uzun siyahı verməsin deyə):
name, err := os.Hostname()
addrs, err := net.LookupHost(name)     // hostname → onun IP-ləri
for _, a := range addrs {
    fmt.Println(a)                    // loopback/IPv6 kirlindən təmiz
}
```
**Prinsip:** runtime-detect > assume — çünki horizontal scale + datacenter müxtəlifliyi konfiqurasiyanı köhnəldir.

### TECHNIQUE 59: Dependency detection
**Problem:** cloud-da minimal Linux distro-lar — `fortune`, `ffmpeg` və s. OLMAYA bilər.

```go
func checkDep(name string) error {
    if _, err := exec.LookPath(name); err != nil {    // PATH-də axtar
        es := "Could not find '%s' in PATH: %s"
        return fmt.Errorf(es, name, err)
    }
    return nil
}

err := checkDep("fortune")
if err != nil {
    log.Fatalln(err)               // yaxud fallback / skip
}
fmt.Println("Time to get your fortune")
```
- `exec.LookPath` — binary PATH-də varmı?
- Addımlar: fail (Fatalln) / fallback / skip — sənin qərarın; ƏSAS: reportable şəkildə bil

### TECHNIQUE 60: Cross-compiling
**Go 1.5+:** daxili cross-compile (CGO qadağası ilə):
```bash
$ GOOS=windows GOARCH=386 go build
# PE32 executable for MS Windows (console) Intel 80386 32-bit
```

**GOOS:** windows, linux, darwin, freebsd... **GOARCH:** amd64, 386, arm...

**gox — paralel çoxlu build:**
```bash
$ go get -u github.com/mitchellh/gox
$ gox \
  -os="linux darwin windows" \
  -arch="amd64 386" \
  -output="dist/{{.OS}}-{{.Arch}}/{{.Dir}}" .
# dist/linux-amd64/, dist/windows-386/ ... hamısı PARALEL
```

**Platform fərqlərinin idarəsi:**
```go
filepath.Separator      // '/'
filepath.ListSeparator  // ':'
filepath.ToSlash(path)  // dönüşdür
filepath.Join(a, b)     // düzgün separator ilə birləşdir
filepath.Split / SplitList
```

**Build tags:**
```go
// +build !windows          // Windows-da BU FAYLI SKIP et
// +build !linux,!darwin    // Linux və OS X-də skip
```

**Fayl adı konvensiyası:** `foo_windows.go` → Windows build-də; `foo_386.go` → 386 build-də.

**XƏBƏRDARLIQ:** cgo istifadəsində cross-compile problemləri — hər platformada test et!

### TECHNIQUE 61: Runtime monitoring
**Real story (müəlliflərdən):** kitabxana bug-u goroutine-ləri öldürmürdü → milyonlarla goroutine yığılırdı (yüzlər olmalıydı) — monitoring bunu GÖSTƏRDİ.

```go
func monitorRuntime() {
    log.Println("Number of CPUs:", runtime.NumCPU())
    m := &runtime.MemStats{}
    for {
        r := runtime.NumGoroutine()
        log.Println("Number of goroutines", r)
        runtime.ReadMemStats(m)
        log.Println("Allocated memory", m.Alloc)
        time.Sleep(10 * time.Second)
    }
}

func main() {
    go monitorRuntime()       // fon goroutine — app ilə paralel
    ...
}
```
- **runtime.NumCPU / NumGoroutine / ReadMemStats**
- **XƏBƏRDARLIQ:** ReadMemStats runtime-u MOMENTAL DAYANDIRIR — performans təsiri; debug rejimində istifadə et / seyrək çağır
- MemStats: GC məlumatları (son pass, növbəti heap həddi, pass müddəti), heap statları, cgo call sayı
- Log əvəzinə New Relic kimi monitoring servisinə göndərmək olar

## Cloud qaydaları xülasəsi
1. Vendor lock-in → interfeys + factory (T56)
2. Divergent error → package error-ləri + original-ı logla (T57)
3. Environment → runtime-detect (T58)
4. Xarici binary → LookPath ilə yoxla (T59)
5. Multi-OS → GOOS/GOARCH + filepath + build tags (T60)
6. Runtime health → NumGoroutine/ReadMemStats monitorinqi (T61)

## Əsas terminlər
- IaaS / PaaS / SaaS (idarəetmə spektri)
- Vendor Lock-In
- Remediation (avtomatik düzəltmə)
- Cloud-Native / Microservices (ön baxış)
- Hypervisor / Guest OS vs Host Kernel (VM vs Container)
- Interface-based Provider Abstraction (printer driver modeli)
- LocalFile (dev/test üçün lokal implementasiya)
- Factory Function (konfiqurasiya əsaslı seçim)
- io.ReadCloser / io.ReadSeeker (interfeys parametrləri)
- Package-level Error Variables (divergent error həlli)
- os.IsNotExist / exec.LookPath
- os.Hostname / net.LookupHost
- GOOS / GOARCH
- gox (paralel cross-compile)
- filepath.Separator / Join / ToSlash
- Build Tags (`// +build`)
- Filename Suffix Convention (foo_windows.go)
- runtime.NumCPU / NumGoroutine / ReadMemStats
- MemStats (GC, heap, cgo)
- New Relic (monitoring servisi)

## Praktik nəticə
- Cloud əməliyyatlarını (storage, compute, DB) HƏMİŞƏ interfeys arxasına — 1 lokal implementasiya dev-i asanlaşdırır, cloud-a keçid 1 factory dəyişməsi.
- Implementasiya xətaları paket error-unə çevrilsin; original yalnız log-da — monitoring alert-ləri üçün.
- Host məlumatlarını kodda assumed YOX, runtime-da detect et — scale/datacenter müxtəlifliyi.
- Xarici binary asılılığı olan hər yerdə LookPath + aydın error.
- Cgo-suz layihələrdə cross-compile 2 env dəyişəni; çoxlu target üçün gox; fərqlər üçün filepath paketi + build tags + fayl adları.
- Runtime monitorinqi (goroutine sayı!) goroutine leak və memory bug-larını əvvələn yaxalayır; ReadMemStats bahalıdır — debug-only/seyrək.

## Mənbə
Pages: 240-257 (PDF), book pages 217-234
