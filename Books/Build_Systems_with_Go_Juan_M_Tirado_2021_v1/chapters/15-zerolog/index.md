# Chapter 15 — Logging with Zerolog (səh. 326-348)

## Bu fəsil nədən bəhs edir?

Go-da logging: standart `log` paketi (Fatal/Panic/Println, prefix, flags,
fayl logger-i) və Zerolog üçüncü tərəf kitabxanası — zero-allocation JSON
logging: səviyyələr, kontekst sahələri, stack trace, ConsoleWriter,
MultiLevelWriter, sub-logger, Hook, Sampler və hlog ilə HTTP inteqrasiyası.

## Əsas fikirlər

### 1. Standart log paketi
**Nədir:** 3 mesaj növü: `Fatal` (Print + os.Exit), `Panic` (Print + panic),
`Println` (sadə çap). Prefix və flag-lərlə xüsusi format.

**Kitabdan kod nümunəsi:**
```go
log.Println("This is a log message")
log.SetPrefix("prefix -> ")
log.Println("This is a log message")   // prefix -> 2021/02/08 ... mesaj
log.SetFlags(log.Lshortfile)            // fayl:sətir əlavə olunur
```

**Xüsusi logger (fayla yazma):**
```go
tmpFile, err := ioutil.TempFile(os.TempDir(), "logger.out")
if err != nil {
    log.Panic(err)
}
logger := log.New(tmpFile, "prefix -> ", log.Ldate)  // writer, prefix, flag
logger.Println("This is a log message")
```

### 2. Zerolog əsasları
**Nədir:** JSON-istiqamətli, zero-allocation logging həlli. Hər mesaj =
səviyyə + vaxt + mesaj JSON obyekti.

**7 səviyyə (Example 15.3):**
```go
import "github.com/rs/zerolog/log"

// log.Panic().Msg("...")   — dayandırır
// log.Fatal().Msg("...")   — dayandırır
log.Error().Msg("This is an error message")
log.Warn().Msg("This is a warning message")
log.Info().Msg("This is an information message")
log.Debug().Msg("This is a debug message")
log.Trace().Msg("This is a trace message")
```
Çıxış: `{"level":"error","time":"2021-02-08T19:00:21+01:00","message":"This is an error message"}`

**Qlobal səviyyə idarəetməsi:**
```go
zerolog.SetGlobalLevel(zerolog.DebugLevel)
log.Debug().Msg("Debug message is displayed")
log.Info().Msg("Info Message is displayed")

zerolog.SetGlobalLevel(zerolog.InfoLevel)
log.Debug().Msg("Debug message is no longer displayed")   // buraxılır
log.Info().Msg("Info message is displayed")
```

### 3. Kontekst sahələri (typed fields)
```go
log.Info().Str("mystr", "this is a string").Msg("")
log.Info().Int("myint", 1234).Msg("")
log.Info().Int("myint", 1234).Str("str", "some string").Msg("And a regular message")
```
- `Str`, `Int`, `Err`, `Interface`, `RawJSON` və s. — typed context
 metodları; JSON-a birbaşa sahə kimi düşür

**Struct + JSON tag:**
```go
type AStruct struct {
    FieldA string
    FieldB int
    fieldC bool     // unexported — çıxmır
}
type AJSONStruct struct {
    FieldA string `json:"fieldA,omitempty"`
    FieldB int    `json:"fieldB,omitempty"`
    fieldC bool
}

log.Info().Interface("a", a).Msg("AStruct")   // FieldA/FieldB adları ilə
log.Info().Interface("b", b).Msg("AJSONStruct") // fieldA/fieldB (tag-lər!)

encoded, _ := json.Marshal(b)
log.Info().RawJSON("encoded", encoded).Msg("Encoded JSON")
```

**Error logging:**
```go
err := errors.New("there is an error")
log.Error().Err(err).Msg("this is the way to log errors")
// {"level":"error","error":"there is an error",...}
```

**Stack trace (pkg/errors ilə):**
```go
import (
    "github.com/pkg/errors"
    "github.com/rs/zerolog/log"
    "github.com/rs/zerolog"
    "github.com/rs/zerolog/pkgerrors"
)

func failA() error { return failB() }
func failB() error { return failC() }
func failC() error { return errors.New("C failed") }

func main() {
    zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

    err := failA()
    log.Error().Stack().Err(err).Msg("")
}
// "stack":[{"func":"failC","line":"19","source":"main.go"},
//           {"func":"failB",...},{"func":"failA",...},{"func":"main",...}]
```
- `zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack` → Stack()
  tam zənciri verir

### 4. Zerolog tənzimləmələri

**Fayl logger (New + io.Writer):**
```go
tempFile, err := ioutil.TempFile(os.TempDir(), "deleteme")
if err != nil {
    log.Error().Err(err).Msg("there was an error creating a temporary file")
}
defer tempFile.Close()
fileLogger := zerolog.New(tempFile).With().Logger()
fileLogger.Info().Msg("This is an entry from my log")
```

**ConsoleWriter (insan oxunaqlı format):**
```go
output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
output.FormatLevel = func(i interface{}) string {
    return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
}
output.FormatMessage = func(i interface{}) string {
    return fmt.Sprintf(">>>%s<<<", i)
}
output.FormatFieldName = func(i interface{}) string {
    return fmt.Sprintf("[%s]:", i)
}
output.FormatFieldValue = func(i interface{}) string {
    return strings.ToUpper(fmt.Sprintf("[%s]", i))
}

log := zerolog.New(output).With().Timestamp().Logger()
log.Info().Str("foo", "bar").Msg("Save the world with Go!!!")
// 2021-02-08T19:16:09+01:00 | INFO  | >>>Save the world with Go!!!<<< [foo]:[BAR]
```

**MultiLevelWriter (çox çıxış):**
```go
fileWriter := zerolog.New(tempFile).With().Logger()
consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout}

multi := zerolog.MultiLevelWriter(consoleWriter, os.Stdout, fileWriter)
logger := zerolog.New(multi).With().Timestamp().Logger()
logger.Info().Msg("Save the world with Go!!!")
// Həm konsolda (2 formatda), həm faylda
```

**Sub-logger (miras + genişlənmə):**
```go
mainLogger := zerolog.New(os.Stdout).With().Logger()
mainLogger.Info().Msg("This is the main logger")

subLogger := mainLogger.With().Str("component", "componentA").Logger()
subLogger.Info().Msg("This is the sublogger")
// {"level":"info","component":"componentA","message":"This is the sublogger"}
```

### 5. Əlavə qurulumlar

**Hook (hər çağırışda icra):**
```go
type ComponentHook struct {
    component string
}

func (h ComponentHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
    if level == zerolog.DebugLevel {          // yalnız debug mesajlarına
        e.Str("component", h.component)
    }
}

type RandomHook struct{}
func (r RandomHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
    e.Int("random", rand.Int())
}

logger := log.Hook(ComponentHook{"moduleA"})
logger = logger.Hook(RandomHook{})
logger.Info().Msg("Info message")     // random sahəsi
logger.Debug().Msg("Debug message")   // component + random
```

**Sampling (loop spam-inin qarşısı):**
```go
// BasicSampler — hər N-dən bir:
logger := log.Sample(&zerolog.BasicSampler{N: 200})
for i := 0; i < 1000; i++ {
    logger.Info().Int("i", i).Msg("")
}
// i:0, i:200, i:400, i:600, i:800

// BurstSampler — periodda M mesaj, sonra NextSampler:
logger := log.Sample(&zerolog.BurstSampler{
    Burst: 2,
    Period: time.Second * 5,
    NextSampler: &zerolog.BasicSampler{N: 90000000},
})
```

**hlog — HTTP handler inteqrasiyası:**
```go
var log zerolog.Logger = zerolog.New(os.Stdout).With().
    Str("app", "example_04").Logger()

func (c MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    log.Info().Msg("This is not a request contextual logger")
    hlog.FromRequest(r).Info().Msg("")    // sorğu kontekstli logger
    w.Write([]byte("Perfect!!!"))
}

func main() {
    mine := MyHandler{}
    a := hlog.NewHandler(log)                        // logger-i handler-ə keçir
    b := hlog.RemoteAddrHandler("ip")                // IP əlavə edir
    c := hlog.UserAgentHandler("user_agent")         // user-agent
    d := hlog.RequestIDHandler("req_id", "Request-Id") // sorğu ID

    panic(http.ListenAndServe(":8090", a(b(c(d(mine))))))
}
// {"level":"info","app":"example_04","ip":"::1","user_agent":"curl/7.64.1","req_id":"..."}
```

**Wrapper — proqramatik qat zənciri (rekursiya ilə):**
```go
type Wrapper struct {
    layers []func(http.Handler) http.Handler
}

func NewWrapper() *Wrapper {
    layers := []func(http.Handler) http.Handler {
        hlog.NewHandler(log),
        hlog.RemoteAddrHandler("ip"),
        hlog.UserAgentHandler("user_agent"),
        hlog.RequestIDHandler("req_id", "Request-Id"),
        hlog.MethodHandler("method"),
        hlog.RequestHandler("url"),
    }
    return &Wrapper{layers}
}

func (w *Wrapper) GetWrapper(h http.Handler, i int) http.Handler {
    if i >= len(w.layers) {
        return h
    }
    return w.layers[i](w.GetWrapper(h, i+1))   // rekursiv bükme
}

h := wrapper.GetWrapper(mine, 0)
panic(http.ListenAndServe(":8090", h))
```

## Əsas terminlər
- Logging — proqram hadisələrinin qeydiyyatı
- Log level — mesajın vaciblik dərəcəsi (panic..trace)
- Zero-allocation — yaddaş ayrılmadan (sürətli) logging
- Context field — mesaja bağlı typed sahə (Str/Int/Err...)
- ConsoleWriter — insan oxunaqlı konsol çıxışı
- MultiLevelWriter — bir mesaj, çox destinasiya
- Hook — hər log çağırışında işə düşən qat
- Sampler — mesaj sıxlığının azaldılması
- hlog — Zerolog-un HTTP handler paketi

## Praktik nəticə
Sadə ehtiyaclar üçün standart `log` kifayətdir (Fatal/Panic/Println + fayl
writer). Real layihələrdə Zerolog: JSON struktur, səviyyə filtri
(`SetGlobalLevel`), typed kontekst (`Err`, `Str`...), stack trace
(`pkgerrors.MarshalStack`), fayl+konsol paralel çıxış (MultiLevelWriter),
insan-formatı konsol (ConsoleWriter). HTTP serverlərdə hlog middleware
zənciri hər sorğunun IP/agent/req_id kontekstini avtomatik əlavə edir.
Loop-larda Sampler ilə spam qarşısı alınır.

## Mənbə
Pages: 326-348 (PDF 326-348)
