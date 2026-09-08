# Chapter 2 — Инструменты командной строки (CLI alətləri)

## Bu chapter nədən bəhs edir?

Command-line flag-lər (flag paketi), arqumentlər və subcommand-lar, mühit
dəyişənləri (envconfig), TOML/YAML/JSON konfiqurasiyası, Unix pipe-lar,
signal tutma (signal.Notify) və ANSI rəngləmə.

## Əsas fikirlər

### 1. Command-line flag-lər (flag paketi)
**Nədir:** CLI tətbiqinə istifadəçi girişi ötürmək üçün standart kitabxana
paketi: `-flag value` sintaksisi.

**Necə işləyir:** Flag dəyərləri struct-a bağlanır (Var funksiyaları),
`flag.Parse()` main-dən çağırılır; `-h` avtomatik yardım yaradır.

**Kitabdan kod nümunəsi:**
```go
type Config struct {
    subject     string
    isAwesome   bool
    howAwesome  int
    countTheWays CountTheWays   // custom type!
}

func (c *Config) Setup() {
    flag.StringVar(&c.subject, "subject", "", "subject is a string, defaults to empty")
    flag.StringVar(&c.subject, "s", "", "shorthand")  // qısa versiya
    flag.BoolVar(&c.isAwesome, "isawesome", false, "is it awesome or what?")
    flag.IntVar(&c.howAwesome, "howawesome", 10, "how awesome out of 10?")
    flag.Var(&c.countTheWays, "c", "comma separated list of integers")
}

// Custom flag tipi — flag.Value interfeysi:
type CountTheWays []int

func (c *CountTheWays) String() string { /* ... */ }
func (c *CountTheWays) Set(value string) error {
    values := strings.Split(value, ",")   // vergüllə böl
    for _, v := range values {
        i, err := strconv.Atoi(v)
        if err != nil {
            return err
        }
        *c = append(*c, i)
    }
    return nil
}

func main() {
    c := Config{}
    c.Setup()
    flag.Parse()     // yalnız main-dən!
    fmt.Println(c.GetMessage())
}
```

**Sub-kod izahı:**
- `flag.StringVar(&c.subject, "s", ...)` → pointer ilə struct-a birbaşa yazır
- `flag.Var(&c.countTheWays, ...)` → öz tipin üçün `String()` + `Set(value)` realizə et
- Bool flag arqumentsiz çağırılır (`-isawesome`), sıra əhəmiyyətsizdir
- Çatışmamazlıq: qısa flag-lər üçün kod dublikasiyası, yardım əlifba sırası

### 2. Command-line arqumentləri və subcommand-lar
**Nədir:** `app greet name -flag` kimi çoxsəviyyəli komandalar — flag
set-lərin birləşməsi ilə.

**Kitabdan kod nümunəsi:**
```go
func (m *MenuConf) SetupMenu() *flag.FlagSet {
    menu := flag.NewFlagSet("menu", flag.ExitOnError)
    menu.Usage = func() {
        fmt.Printf(usage, os.Args[0])
        menu.PrintDefaults()
    }
    return menu
}

func (m *MenuConf) GetSubMenu() *flag.FlagSet {
    submenu := flag.NewFlagSet("submenu", flag.ExitOnError)
    submenu.BoolVar(&m.Goodbye, "goodbye", false, "Say goodbye instead of hello")
    // ...
    return submenu
}

func main() {
    c := MenuConf{}
    menu := c.SetupMenu()
    if err := menu.Parse(os.Args[1:]); err != nil { return }
    if len(os.Args) > 1 {
        switch strings.ToLower(os.Args[1]) {   // subcommand seçimi
        case "version":
            c.Version()
        case "greet":
            f := c.GetSubMenu()
            if len(os.Args) > 3 {
                if err := f.Parse(os.Args[3:]); err != nil { return } // flag-lər
            }
            c.Greet(os.Args[2])                // positional arqument
        default:
            menu.Usage()
        }
    }
}
```

**Sub-kod izahı:**
- `flag.NewFlagSet(ad, flag.ExitOnError)` → müstəqil flag dəsti
- `os.Args[1]` → komanda, `os.Args[2]` → positional arqument, `os.Args[3:]` → flag-lər
- `menu.Usage = func(){...}` → öz yardım mətnini təyin etmək olar

### 3. Mühit dəyişənləri + envconfig
**Nədir:** `os.Getenv`/`os.Setenv` + struct tag-lərlə avtomatik env oxuma
(kelseyhightower/envconfig).

**Kitabdan kod nümunəsi:**
```go
// Fayldan (JSON) + env-dən birləşdirilmiş konfiq yükləmə:
func LoadConfig(path, envPrefix string, config interface{}) error {
    if path != "" {
        if err := LoadFile(path, config); err != nil {
            return errors.Wrap(err, "error loading config from file")
        }
    }
    err := envconfig.Process(envPrefix, config)   // env DƏYƏRLƏRİ FAYLI OVERRIDE EDİR
    return errors.Wrap(err, "error loading config from env")
}

// Tag-lərlə tələblər:
type Config struct {
    Version string `json:"version" required:"true"`
    IsSafe  bool   `json:"is_safe" default:"true"`
    Secret  string `json:"secret"`
}
// EXAMPLE_VERSION=1.0.0, EXAMPLE_ISSAFE=false env-ləri fayldakıları əvəz edir
```

**Sub-kod izahı:**
- `errors.Wrap(err, "...")` → xəta annotasiyası — orijinal xəta itmir (Chapter 4)
- `envconfig.Process("EXAMPLE", &c)` → EXAMPLE_ prefiksli env-ləri struct-a yazır
- `required:"true"` → yoxdursa xəta; `default:"true"` → default dəyər

### 4. TOML / YAML / JSON konfiqurasiyası
**Nədir:** Üç populyar data formatı — struct tag-lərlə Marshal/Unmarshal.

**Kitabdan kod nümunəsi:**
```go
// TOML (BurntSushi/toml):
type TOMLData struct {
    Name string `toml:"name"`
    Age  int    `toml:"age"`
}
func (t *TOMLData) ToTOML() (*bytes.Buffer, error) {
    b := &bytes.Buffer{}
    if err := toml.NewEncoder(b).Encode(t); err != nil { return nil, err }
    return b, nil
}
func (t *TOMLData) Decode(data []byte) (toml.MetaData, error) {
    return toml.Decode(string(data), t)
}

// YAML (go-yaml/yaml):      yaml.Marshal(t) / yaml.Unmarshal(data, t)
// JSON (encoding/json):     json.Marshal(t) / json.Unmarshal(data, t)

// JSON-in stream variantı:
res := make(map[string]string)
b := bytes.NewReader([]byte(`{"key2": "value2"}`))
decoder := json.NewDecoder(b)
decoder.Decode(&res)      // struct YOX, map + decoder ilə stream
```

**Sub-kod izahı:**
- Hər formatın öz struct tag-i: `toml:"..."`, `yaml:"..."`, `json:"..."`
- JSON ən tam imkanlıdır (encoder/decoder, stream dəstəyi)
- []byte ↔ string ↔ bytes.Buffer çevirmələri sürətlə edilir

### 5. Unix pipe-lar (os.Stdin)
**Nədir:** `echo "test" | ./app` — sol tərəfin çıxışı proqrama stdin ilə
daxil olur.

**Kitabdan kod nümunəsi:**
```go
func WordCount(f io.Reader) map[string]int {
    result := make(map[string]int)
    scanner := bufio.NewScanner(f)      // f = os.Stdin
    scanner.Split(bufio.ScanWords)
    for scanner.Scan() {
        result[scanner.Text()]++       // söz sayğıcı
    }
    if err := scanner.Err(); err != nil {
        fmt.Fprintln(os.Stderr, "reading input:", err)
    }
    return result
}
func main() {
    for key, value := range WordCount(os.Stdin) {
        fmt.Printf("%s: %d\n", key, value)
    }
}
```

**Sub-kod izahı:**
- `os.Stdin` → io.Reader kimi fayl deskriptoru
- Oxuma bitəndən SONRA `scanner.Err()` yoxlanılır
- Chapter 1 pipe nümunəsi ilə birləşdirib `tee` proqramı yazmaq olar

### 6. Signal tutma və emalı
**Nədir:** SIGINT (Ctrl+C), SIGTERM (kill) siqnallarını tutub graceful
shutdown etmək.

**Necə işləyir:** `signal.Notify` kanala siqnal göndərir; ayrı goroutine
kanalı gözləyir; `done` kanalı main-i bloklayır.

**Kitabdan kod nümunəsi:**
```go
func CatchSig(ch chan os.Signal, done chan bool) {
    sig := <-ch                          // siqnal gözlə — bloklanır
    switch sig {
    case syscall.SIGINT:
        fmt.Println("handling a SIGINT now!")
    case syscall.SIGTERM:
        fmt.Println("handling a SIGTERM in an entirely different way!")
    default:
        fmt.Println("unexpected signal received")
    }
    done <- true                          // main-i azad et
}

func main() {
    signals := make(chan os.Signal)
    done := make(chan bool)
    signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
    go CatchSig(signals, done)
    fmt.Println("Press ctrl-c to terminate...")
    <-done                                // siqnal gələnə qədər gözlə
    fmt.Println("Done!")
}
```

**Sub-kod izahı:**
- `signal.Notify(signals, SIGINT, SIGTERM)` → yalnız bu siqnallar kanala düşür
- `done <- true` / `<-done` → goroutine sinkronizasiyası
- Praktik: HTTP handler-lərin cari sorğuları bitirməsinə imkan vermək —
  çox goroutine-li vəziyyətli tətbiqlərdə təmizləmə üçün vacibdir

### 7. ANSI rəngləmə
**Nədir:** Terminalda mətn rəngləmə — ANSI escape kodları.

**Kitabdan kod nümunəsi:**
```go
type Color int
const (
    ColorNone = iota
    Red; Green; Yellow; Blue; Magenta; Cyan; White
    Black Color = -1
)
type ColorText struct {
    TextColor Color
    Text      string
}
func (r *ColorText) String() string {
    if r.TextColor == ColorNone {
        return r.Text                       // rəngsiz
    }
    value := 30
    if r.TextColor != Black {
        value += int(r.TextColor)            // ANSI kod hesabı
    }
    return fmt.Sprintf("\033[0;%dm%s\033[0m", value, r.Text)
}
```

**Sub-kod izahı:**
- `\033[0;Nm` → rəng başlanğıcı, `\033[0m` → sıfırlama
- `Stringer` interfeysi → fmt.Println avtomatik String() çağırır
- Tam həll: github.com/agtorre/gocolorize (fmt.Formatter realizə edir)

## Əsas terminlər

- Flag / FlagSet (bayraq dəsti)
- flag.Value interfeysi (custom flag tipi)
- Positional Argument (mövqe arqumenti)
- Subcommand (alt komanda)
- Environment Variable (mühit dəyişəni)
- envconfig struct tags
- Marshal / Unmarshal (seriyalaşdırma)
- Unix Pipe (boru — | operatoru)
- Signal / SIGINT / SIGTERM (siqnallar)
- signal.Notify (siqnal abunəliyi)
- Graceful Shutdown (təmiz bağlanma)
- ANSI Escape Codes (rəng kodları)

## Praktik nəticə

- Flag-ləri struct-a topla — uzunmüddətli baxımdan daha yaxşıdır
- Custom flag üçün String()+Set() kifayətdir
- Konfiq: fayl (JSON) + env override modeli — LoadConfig pattern
- stdin üzərində scanner ilə Unix utilitləri (wc kimi) yazmaq asandır
- Siqnal tutma: uzun ömürlü servislərdə graceful shutdown üçün mütləq
- Rəng kodları log oxunaqlılığını artırır (səhvlər qırmızı və s.)

## Mənbə

Pages: 58-92 (Chapter 2, Go Programming Cookbook 2nd ed)
