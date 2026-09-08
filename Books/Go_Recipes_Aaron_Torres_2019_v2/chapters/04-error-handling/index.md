# Chapter 4 — Обработка ошибок в Go (Xəta emalı)

## Bu chapter nədən bəhs edir?

Error interfeysi və xəta yaratma üsulları, pkg/errors ilə wrap/unwrap,
log paketi, strukturlaşdırılmış logging (logrus, apex), context ilə log
sahələri, sync.Once ilə qlobal logger, uzunmüddətli proseslərdə panic
tutma (recover).

## Əsas fikirlər

### 1. Error interfeysi və xəta yaratma
**Nədir:** Go xətaları exception DEYİL — sadə interfeys: `Error() string`.
Hər səviyyədə explicit emal tələb olunur.

**Kitabdan kod nümunəsi:**
```go
type Error interface {
    Error() string
}

// 4 yaratma üsulu:
err := errors.New("quick and easy way")        // 1) sadə
err = fmt.Errorf("an error occurred: %s", x)   // 2) formatlı

// 3) Package-level error VALUE (müqayisə üçün):
var ErrorValue = errors.New("this is a typed error")
if err == ErrorValue { /* bu xətanı xüsusi emal et */ }

// 4) Error TYPE (type switch üçün):
type TypedError struct{ error }
err = TypedError{errors.New("typed error")}

// 5) Custom struct error (əlavə sahələrlə):
type CustomError struct{ Result string }
func (c CustomError) Error() string {
    return fmt.Sprintf("there was an error; %s was the result", c.Result)
}
```

**Sub-kod izahı:**
- `errors.New` → sabit mətnli xəta; package-level dəyişən kimi saxlanılır
- `fmt.Errorf` → formatlaşdırılmış xəta
- Struct error → əlavə sahələr + metodlar; çağıranda type assertion ilə
  dindirilə bilir

### 2. pkg/errors — wrap və unwrap
**Nədir:** Standart errors paketinin əvəzi; xətaları kontekstlə örtmək
(annotation) üçün Wrap/Cause funksiyaları.

**Nəyə lazımdır:** Adi wrap (`fmt.Errorf("custom: %s", err)`) tipi DƏYİŞİR —
type assertion pozulur. pkg/errors orijinal tipi saxlayır.

**Kitabdan kod nümunəsi:**
```go
// Wrap — nil-i də düzgün idarə edir (nil qaytarır):
func WrappedError(e error) error {
    return errors.Wrap(e, "An error occurred in WrappedError")
}
// errors.Wrap(nil, "msg") → nil!

// Unwrap — Cause ilə orijinal tipə qayıdış:
func Unwrap() {
    err := error(ErrorTyped{errors.New("an error occurred")})
    err = errors.Wrap(err, "wrapped")
    switch errors.Cause(err).(type) {   // Cause orijinal tipi verir
    case ErrorTyped:
        fmt.Println("a typed error occurred: ", err)
    default:
        fmt.Println("an unknown error occurred")
    }
}

// Stack trace — %+v formatı ilə:
err = errors.Wrap(err, "wrapped")
fmt.Printf("%+v\n", err)   // tam stack trace çap olunur
```

**Sub-kod izahı:**
- `errors.Wrap(err, "ctx")` → "ctx: original message"; stack trace saxlayır
- `errors.Cause(err)` → wrap zəncirinin ən altdakı orijinal xətası
- `%+v` → xətanın tam stack trace-i
- Nil-check ehtiyacını aradan qaldırır: `return errors.Wrap(err, "...")`

### 3. log paketi və xətanın HARADA loglanması
**Nədir:** Standart log paketi; xəta yalnız SON təyinatında loglanmalıdır.

**Kitabdan kod nümunəsi:**
```go
// Logger konfiqurasiyası:
logger := log.New(&buf, "logger: ", log.Lshortfile|log.LDate)
// Lshortfile → fayl:sətir; LDate → tarix; | (OR) ilə birləşdirilir
logger.Println("test")
logger.SetPrefix("new logger: ")   // prefiks dəyişdirilir
logger.Printf("you can also add args(%v)", true)

// Pass-through pattern — yolda log YOX:
func PassThroughError() error {
    err := OriginalError()
    return errors.Wrap(err, "in passthrougherror")  // wrap et, ötür
}

// Final destination — burada log:
func FinalDestination() {
    err := PassThroughError()
    if err != nil {
        log.Printf("an error occurred: %s\n", err.Error())
        return
    }
}
```

**Sub-kod izahı:**
- Xəta zənciri boyu yalnız wrap; LOG yalnız son nöqtədə — yoxsa eyni
  hadisə üçün bir neçə dəfə log yazılır
- `log.New(out, prefix, flag)` → öz logger obyektin; standart package-level
  logger də var (`log.Printf`)

### 4. Strukturlaşdırılmış logging (logrus, apex)
**Nədir:** Mesaj deyil, SAHƏLƏR (key-value) ilə logging — mikroservis
loglarının axtarışı/indexing üçün (JSON formatına çevrilir).

**Kitabdan kod nümunəsi:**
```go
// LOGRUS:
logrus.SetFormatter(&logrus.TextFormatter{})   // JSONFormatter də olar
logrus.SetLevel(logrus.InfoLevel)
logrus.AddHook(&Hook{"123"})                    // hər entry-ə id əlavə et

type Hook struct{ id string }
func (hook *Hook) Fire(entry *logrus.Entry) error {
    entry.Data["id"] = hook.id                   // hook = entry modifikasiyası
    return nil
}
func (hook *Hook) Levels() []logrus.Level { return logrus.AllLevels }

fields := logrus.Fields{"success": true, "complex_struct": struct{...}{...}}
x := logrus.WithFields(fields)
x.Warn("warning!")
x.Error("error!")

// APEX:
log.WithField("id", "123").Trace("ThrowError").Stop(&err)  // trace + duration
log.SetHandler(&CustomHandler{"123", text.New(os.Stdout)})  // handler modeli
log.WithError(err).Error("an error occurred")               // error field
```

**Sub-kod izahı:**
- logrus: Hook (Fire + Levels) + Formatter ayrı-ayrı; JSON-a çevrilmə asan
- apex: Handler birleşdirir; `WithError(err)`, `Trace().Stop(&err)`
  (müddət ölçmə) convenience funksiyaları
- Hook-dan traceID əlavə etmək — request-in servislər arası izlənməsi

### 5. context ilə logging
**Nədir:** Log sahələrini funksiyalar arası context vasitəsilə daşımaq.

**Kitabdan kod nümunəsi:**
```go
type key int
const logFields key = 0   // untyped key — konflikt qorunması

func getFields(ctx context.Context) *log.Fields {
    if fields, ok := ctx.Value(logFields).(*log.Fields); ok {
        return fields
    }
    f := make(log.Fields)
    return &f
}

func WithField(ctx context.Context, key string, value interface{}) context.Context {
    return WithFields(ctx, log.Fields{key: value})
}
func WithFields(ctx context.Context, fields log.Fielder) context.Context {
    f := getFields(ctx)
    for key, val := range fields.Fields() {
        (*f)[key] = val                   // sahələri topla
    }
    return context.WithValue(ctx, logFields, f)
}

// İstifadə: hər funksiya öz sahəsini əlavə edir:
ctx, e := FromContext(ctx, log.Log)
ctx = WithField(ctx, "id", "123")
e.Info("starting")                 // id=123
gatherName(ctx)                    // name əlavə olunur
e.Info("after gatherName")         // id=123 name=Go Cookbook
```

**Sub-kod izahı:**
- `context.WithValue` ilə sahələr yığılır; son log nöqtəsində hamısı çap olunur
- Bu pattern elə də populyar deyil, amma context-in digər üstünlükləri
  (cancel, timeout) ilə birləşir

### 6. sync.Once ilə qlobal logger
**Nədir:** Package-level logger — yalnız bir dəfə init; ixrac olunan
qısa funksiyalar (WithField, Debug).

**Kitabdan kod nümunəsi:**
```go
var (
    log     *logrus.Logger   // lowercase — ixrac olunmur!
    initLog sync.Once
)

func Init() error {
    err := errors.New("already initialized")
    initLog.Do(func() {          // YALNIZ bir dəfə icra olunur
        err = nil
        log = logrus.New()
        log.Formatter = &logrus.JSONFormatter{}
        log.Out = os.Stdout
        log.Level = logrus.DebugLevel
    })
    return err                   // ikinci çağırışda error
}

func SetLog(l *logrus.Logger) { log = l }       // bypass üçün
func WithField(key string, value interface{}) *logrus.Entry {
    return log.WithField(key, value)            // qlobal logger-ə delegate
}
func Debug(args ...interface{}) { log.Debug(args...) }
```

**Sub-kod izahı:**
- `sync.Once.Do(fn)` → fn yalnız bir dəfə, thread-safe icra olunur
- Global dəyişən gizlidir; istifadəçi yalnız funksiyalarla işləyir
- `Init()` funksiyası `init()`-dən üstündür — parametr ötürmək olar

### 7. Panic tutma (recover)
**Nədir:** Uzun ömürlü proseslərdə panic-i yox, tam crash-i qarşısını almaq —
defer + recover.

**Kitabdan kod nümunəsi:**
```go
func Panic() {
    zero, _ := strconv.ParseInt("0", 10, 64)
    a := 1 / zero          // divide by zero → runtime panic
    fmt.Println("never get here", a)
}

func Catcher() {
    defer func() {
        if r := recover(); r != nil {     // panic-i yakala
            fmt.Println("panic occurred:", r)
        }
    }()
    Panic()      // panic baş versə də Catcher davam edir
}
// main: "before panic" → "panic occurred: ..." → "after panic"
```

**Sub-kod izahı:**
- `defer func(){ recover() }` → panic-i yalnız eyni goroutine-da defer
  zənciri ilə tutmaq mümkündür
- Veb tətbiqlərində: recover + http.StatusInternalServerError qaytarma
  standart pattern-dir (middleware ilə)
- Panic mənbələri: init olunmamış map/pointer, divide-by-zero

## Əsas terminlər

- Error Interface (xəta interfeysi)
- errors.New / fmt.Errorf
- Error Value / Error Type (müqayisə/tip yoxlaması üçün)
- Error Wrapping (xəta örtülməsi) / errors.Cause
- Stack Trace (stek izi)
- Log Levels (səviyyələr: Debug, Info, Warn, Error)
- Structured Logging (strukturlaşdırılmış logging)
- Hook / Handler / Formatter (logrus/apex anlayışları)
- context.WithValue
- sync.Once (bir dəfəlik icra)
- Panic / Recover / Defer

## Praktik nəticə

- Xətanı heç vaxt nəzərə almadan buraxma — hər səviyyədə ya emal et, ya wrap ötür
- pkg/errors: Wrap ilə kontekst əlavə et, Cause ilə orijinal tipi yoxla
- LOG yalnız xətanın son təyinatında — yolda yox (duplicate log problemi)
- Struktur log (JSON) mikroservislərdə axtarışı mümkün edir
- Uzun proseslərdə recover + log; vebdə HTTP 500 middleware
- sync.Once: logger/DB init üçün ideal — bir dəfəlik, thread-safe

## Mənbə

Pages: 133-162 (Chapter 4, Go Programming Cookbook 2nd ed)
