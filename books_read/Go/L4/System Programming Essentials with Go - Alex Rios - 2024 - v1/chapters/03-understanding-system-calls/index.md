# Chapter 3 — Understanding System Calls (Sistem Çağırışlarını Anlamaq)

## Bu chapter nədən bəhs edir?
Sistem çağırışlarının təbiətinə (user/kernel mode, pasport analogiyası), Go-nun syscall/os/x/sys paketlərinə, strace ilə syscall izləməyə, standard stream-lərə (stdin/stdout/stderr + file descriptor-lər) və testable CLI proqramın qurulmasına (Functional Options pattern-i ilə).

## Əsas fikirlər

### 1. Sistem çağırışı (syscall) nədir
**Nədir:** Kernel-in user-proseslərə təqdim etdiyi aşağı səviyyəli funksiyalar — proseslərin kernel-dən xidmət istəməsinin yolu.

**Səyahət analogiyası:**
- **Kernel = sərhəd məntəqəsi** (güclü qorunan keçid)
- **Syscall = pasport** — user space-dən kernel space-ə keçid üçün icazə
- **Xidmət kataloqu = bələdçi kitabı** — syscall API (proses yaratma, I/O...)
- **Syscall nömrəsi = pasport №** — kernel daxili cədvəldə rəqəmlə təmsil olunur (məs. open = №5), amma biz adı ilə işlədirik
- **Arqumentlər = sərhəd söhbəti** — data user↔kernel arasında mübadilə (write() üçün: fd + buffer)

### 2. User mode vs Kernel mode
| | User mode | Kernel mode (supervisor/privileged) |
|---|---|---|
| Giriş | MƏHDUD — kritik resurslara birbaşa YOX | TAM — bütün resurslar |
| Səbəb | Təhlükəsizlik izolyasiyası | Kernel öz işi üçün |

CPU-nun iki rejimi — syscall user-dən kernel-ə "sərhəd keçidi" ilə müraciət edir.

**Syscall cədvəli:** OS+arxitekturadan asılı; https://filippo.io/linux-syscall-table/ — maraqlı üçün.

### 3. syscall paketinin süqutu (5 səbəb)
| Problem | Təsvir |
|---|---|
| **Bloat** | Hər syscall+konstant üçün müdafiə — "dolab həddindən artıq dolu" |
| **Testing** | Açıq testlər YOX; cross-platform test imkansız |
| **Curation** | "Vəhşi qərb" — istənilən dəyişiklik qəbul olunurdu → ən azı saxlanılan/test edilən paket rekoru |
| **Dokumentasiya** | Sistem-özəl varyasiyalar; godoc "yalnız treyler göstərir" |
| **Compatibility** | OS dəyişir (FreeBSD nümunəsi) — "tək buynuzu qovalamaq" |

**Go komandasının qərarı (Rob Pike, Go 1.3):**
1. **Freeze** — syscall paketi donduruldu
2. **x/sys** — yeni, saxlanıla bilən paket
3. **Deprecation** — yeni inkişaf x/sys-ə keçdi

### 4. x/sys — aşağı səviyyəli alternativ
```bash
go get -u golang.org/x/sys
```
**API (unix alt-paketi):**
- `unix.Syscall() / unix.Syscall6()` — birbaşa syscall çağırışı
- `unix.SYS_READ, unix.SYS_WRITE` — SYS_* konstantları
- Fayl: Create/Unlink/Mkdir/Rmdir/Link/Getdents
- Siqnal: Kill, SIGINT
- User/Group: Setuid/Setgid/Setgroups
- Sysinfo; FcntlInt, Dup2; Mmap

**fmt vs x/sys müqayisəsi (eyni nəticə):**
```go
fmt.Println("Hello World!")
// vs:
unix.Syscall(unix.SYS_WRITE, 1,
    uintptr(unsafe.Pointer(&[]byte("Hello, World!")[0])),
    uintptr(len("Hello, World!")),
)
```
**Dərs:** Aşağı səviyyə qısa məsafədə ÇƏTİNLƏŞİR — əksər hallarda yüksək abstraksiya seç.

### 5. os paketi — portativ üst qat
**Fayl/kataloq:** Create, Mkdir/MkdirAll, Remove/RemoveAll, Stat, Open, Rename, Truncate, Getwd, Chdir, IsExist/IsNotExist/IsPermission.
**Proses/siqnal:** Getpid, Getppid, Getuid/Getgid (effective: Geteuid/Getegid), StartProcess, Exit, Signal, os/signal.Notify.
**Mühit:** os.Args, Getenv, Setenv.

**Proses yaratma nümunəsi:**
```go
cmd := exec.Command("ls", "-l")
cmd.Stdout = os.Stdout
cmd.Stderr = os.Stderr
err := cmd.Run()
pid := os.Getpid()
```

### 6. Portativlik — "imza tetrisi"
Windows: `SetEnvironmentVariable(name *uint16, value *uint16) error`
Unix: `Setenv(key, value string) error`
os paketi: `Setenv(key, value string) error` — HƏR İKİ platformada eyni!

**Qayda (rəsmi tövsiyə):** "The primary use of x/sys is inside other packages" — istifadəçi paketləri: **os, time, net**. x/sys yalnız müstəsna hallarda.

### 7. strace — syscall izləyici
```bash
apt-get install strace -y    # Debian
strace ls                     # bütün syscall-lar
strace -e execve ls          # yalnız execve

# Go proqramda:
go build -o app main.go
strace -e write ./app 2>&1
# write(1, "Hello, World!", 13Hello, World!) = 13
```

### 8. Standard stream-lər
**Unix fəlsəfəsi:** "everything is a file".

| FD | Ad | Vəzifə |
|----|----|--------|
| 0 | stdin | giriş (klaviatura/pipe/fayl) |
| 1 | stdout | nəticə çıxışı |
| 2 | stderr | xəta çıxışı |

**Niyə vacibdir:** pipeline inteqrasiyası (`ls -l | xargs app | grep even`), input/output çevikliyi, xəta ayrılması, cross-platform, test/debug rahatlığı, logging konvensiyası.

**File descriptor tipləri:** regular fayllar, kataloqlar, character device-lar (klaviatura), block device-lar (disk), socket-lər, pipe-lar.

### 9. CLI proqram nümunəsi
```go
func main() {
    words := os.Args[1:]           // program adı istisna
    if len(words) == 0 {
        fmt.Fprintln(os.Stderr, "No words provided.")
        os.Exit(1)                  // nonzero = xəta
    }
    for _, w := range words {
        if len(w)%2 == 0 {
            fmt.Fprintf(os.Stdout, "word %s is even\n", w)
        } else {
            fmt.Fprintf(os.Stderr, "word %s is odd\n", w)
        }
    }
}
```
**Redirection testi:**
```bash
go run main.go word1 word2 word3 > stdout.txt 2> stderr.txt
# >  = fd 1 (1 öz-özündən görünür), 2> = fd 2
```

### 10. Testable refaktor — 3 addım
**Addım 1 — məntiqi ayır:**
```go
func app(words []string, cfg CliConfig) { ... }
```

**Addım 2 — konfiqurasiya strukturu:**
```go
type CliConfig struct {
    ErrStream, OutStream io.Writer
}
```

**Addım 3 — Functional Options pattern-i:**
```go
type Option func(*CliConfig) error

func WithErrStream(errStream io.Writer) Option {
    return func(c *CliConfig) error {
        c.ErrStream = errStream
        return nil
    }
}
func WithOutStream(outStream io.Writer) Option { /* eyni pattern */ }

func NewCliConfig(opts ...Option) (CliConfig, error) {
    c := CliConfig{
        ErrStream: os.Stderr,    // DEFAULT dəyərlər
        OutStream: os.Stdout,
    }
    for _, opt := range opts {
        if err := opt(&c); err != nil { return CliConfig{}, err }
    }
    return c, nil
}

// İstifadə — istədiyin qədər option:
NewCliConfig(WithOutStream(&var1), WithErrStream(&var2))
```
**Patternin 3 faydası:** oxunaqlılıq (parametr sırası yox), genişlənəbilənlik (imza dəyişmir), təhlükəsizlik (default + həmişə valid vəziyyət).

### 11. Test — bytes.Buffer ilə stream tutma
```go
func TestMainProgram(t *testing.T) {
    var stdoutBuf, stderrBuf bytes.Buffer
    config, err := NewCliConfig(
        WithOutStream(&stdoutBuf),
        WithErrStream(&stderrBuf),   // real stdout/stderr əvəzinə buffer!
    )
    if err != nil { t.Fatal("Error creating config:", err) }

    app([]string{"main", "alex", "golang", "error"}, config)

    output := stdoutBuf.String()
    if !strings.Contains(output, "word alex is even") {
        t.Fatal("Expected output does not contain 'word alex is even'")
    }
    errors := stderrBuf.String()
    if !strings.Contains(errors, "word error is odd") {
        t.Fatal("Expected errors does not contain 'word error is odd'")
    }
}
```
**Mexanizm:** io.Writer interfeysi həm fayl, həm buffer — testdə buffer qoyub çıxışı yoxla.

## Əsas terminlər
- System call (syscall) — kernel xidmət interfeysi
- User mode / Kernel mode (supervisor mode) — CPU icra rejimləri
- Syscall table — nömrə↔ad xəritəsi
- syscall package (frozen, Go 1.3) → x/sys paketi
- unix.Syscall / SYS_* — birbaşa çağırış qatı
- os package — portativ üst qat
- strace / -e flag — syscall tracer
- Standard streams — stdin(0)/stdout(1)/stderr(2)
- File descriptor — açıq resursun rəqəmsal ID-si
- Redirection — `>` (stdout) / `2>` (stderr)
- Functional Options — Option funksiyaları ilə konfiqurasiya
- bytes.Buffer — testdə stream qəbulçusu

## Praktik nəticə
1. syscall paketinə YENİ kod YAZMA — x/sys (müstəsna hallar) və ya os (standart).
2. CLI-lərdə konvensiya: nəticə → stdout, xəta → stderr, exit code 0/1.
3. Testable CLI = məntiqi funksiyaya köçür + stream-ləri io.Writer parametr et (Functional Options).
4. strace ilə "gomlum nə edir" sualına cavab: `strace -e <syscall> ./app`.
5. Pipeline dostu ol: `> output` redirect oluna bilən çıxış ver — istifadəçi onu grep/sort ilə birləşdirəcək.

## Mənbə
Pages: 39-59 (PDF səh. 60-81)
