# Chapter 8 — Examining Unusual Patterns for Go (səh. 277-316)

## Bu chapter nədən bəhs edir?

Go-nun qeyri-adi pattern-ləri: funksional proqramlaşdırma (FP) konseptləri Go-da,
closures, function-based pattern-lər (abstract method, middleware, functional
options, decorator, function chaining, futures) və struct fndləri (empty struct,
anonymous struct, noCopy).

## Əsas fikirlər

### 1. Functional Programming Go-da
**FP-nin 3 fundamental ideyası:**
1. **Pure functions (təmiz funksiyalar)** — eyni input → həmişə eyni output, yan
   təsir YOX. `TimeRemaining(when time.Time)` `time.Now()`-dan asılıdır → pure
   deyil; `TimeRemaining(when, now time.Time)` — pure (idempotent).
2. **Immutability (dəyişməzlik)** — dəyişən təyin olunduqdan sonra dəyişmir
3. **No implicit state (örtülü state yoxdur)** — debug üçün obyektin bütün
   interaksiyalarını izləmək lazım deyil

**FP-nin faydaları:** idempotent → cache oluna bilən, testi asan, paralel
təhlükəsiz (referential transparency).

### 2. First-class və higher-order funksiyalar
**Nədir:** Funksiyalar dəyişənə mənimsədilə, parametr kimi ötürülə bilər.

```go
func Map(in []string, operation func(string) string) []string {
    out := make([]string, len(in))
    for index, item := range in {
        out[index] = operation(item) // higher-order: operation funksiyası çağrılır
    }
    return out
}
Map(in, toUpper)   // funksiyalar data kimi ötürülür
Map(in, addPeriod)
```

**Recursion (rekursiya) ehtiyatlı:** Hər çağırış stack frame yaradır — 1M elementli
slice-də stack overflow. Pure FP dillərindən fərli olaraq Go-da tail-call
optimizasiyası YOXDUR.

### 3. Currying və partial functions
```go
func Multiply(a int) func(int) int {
    return func(b int) int { // partial function: "b gözləyən" funksiya
        return a * b
    }
}

multiply5 := Multiply(5) // partial: a=5 "yadda saxlanılıb"
result := multiply5(3)   // 15
// Qısa forma: Multiply(5)(3)
```
**Sub-kod izahı:**
- Xarici funksiya qayıtdıqda `a` closure-da "saxlanılır"
- Çoxparametrli funksiyanı addım-addım tətbiq etmək = currying

### 4. Immutability Go-da məhduddur
`const` yalnız təmiz tip-lərə; map/slice-lər mütləq mutable. Qorunma: funksiyaya
verməzdən əvvəl kopyala (ideal deyil). For loop-da `result`-un dəyişməsi FP
pozuntusudur (amma Go-da normal).

### 5. Anonymous functions və closures
**Anonymous function:** adsız funksiya — dərhal çağırılan və ya dəyişənə mənimsədilən:
```go
var GetDB = func() *sql.DB { // monkey patching üçün function var
    dbInit.Do(func() { db, _ = sql.Open("mysql", "user:password@/dbname") })
    return db
}
```
Testdə `GetDB` mock-la əvəz olunur.

**Closure (bağlanma):** Xarici scope-in dəyişənlərinə müraciət edən daxili
funksiya:
```go
func buildGreeting(name string) func() string {
    greeting := fmt.Sprintf("Hello %s. Nice to meet you!", name)
    return func() string { return greeting } // greeting closure-da yaşayır
}
```
**State daşıyıcısı kimi closure:** sayağın cari total-u yalnız o closure üçün
görünür — qlobal dəyişəndən fərqli olaraq **tam izolyasiya**.

### 6. Function-based pattern-lər

**Abstract method (template method):**
```go
type Image struct {
    toBytes func() []byte // "abstract" — gecikməli bənd
}
func (i *Image) Save(filename string) {
    destination := i.openFile(filename)
    data := i.toBytes() // konkret encoder PNG/JPG tərəfindən qurulur
    i.writeToFile(data, destination)
}
myPNG := &Image{toBytes: pngEncoder}
```
OO abstract class-ın funksiya ilə Go ekvivalenti — common logic bir yerdə.

**Middleware:**
```go
func trackRequest(next http.HandlerFunc) http.HandlerFunc {
    return func(resp http.ResponseWriter, req *http.Request) {
        start := time.Now()
        defer func() { log.Printf("took %v", time.Since(start)) }()
        next(resp, req) // handler-i wrap edən closure
    }
}
http.Handle("/", trackRequest(http.HandlerFunc(ListUsersHandler)))
```

**Functional Options (çoxlu konfiqurasiyalı konstruktor):**
```go
type Option func(*StatsClient) // konfiq funksiyası

func WithSampleRate(rate int) Option { return func(c *StatsClient) { c.sampleRate = rate } }

func NewStatsDClient(host string, options ...Option) *StatsClient {
    client := &StatsClient{host: host} // host = tələb olunan (parametr)
    for _, option := range options {   // optional-lar tətbiq olunur
        option(client)
    }
    return client
}
NewStatsDClient("statsd.example.com", WithSampleRate(10))
```
**Qərar qaydası:** tələb olunan konfiq → parametr; optional → Option.

**Function-to-interface decorator:**
```go
type MyHandlerFunc func(resp http.ResponseWriter, req *http.Request)

func (m MyHandlerFunc) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
    m(resp, req) // funksiya → http.Handler interfeysi
}
```
Standart `http.HandlerFunc`-un öz dəbdə olan forması.

**Function chaining (pipeline):**
```go
type ChainedFunction func(in <-chan *User) <-chan *User

func EmailIsValid(in <-chan *User) <-chan *User {
    out := make(chan *User)
    go func() {
        defer close(out)
        for user := range in {
            if _, err := mail.ParseAddress(user.Email); err == nil {
                out <- user // filter
            }
        }
    }()
    return out
}
// istifadə: in = chainedFunc(in) — çıxış → növbətinin girişi
```
Goroutine + channel ilə Unix pipeline-ə bənzər data axını.

### 7. Futures (gələcək dəyərlər)
**Nədir:** Funksiya nəticənin ÖZÜNÜ deyil, ona **proxy** qaytarır — paralel
hesablama üçün.

**Closure + mutex versiyası:**
```go
func doWorkFuture(in int64) func() (int64, sync.Mutex) {
    var result int64
    var mux sync.Mutex
    mux.Lock()
    go func() {
        defer mux.Unlock() // oxunuş azad olur
        result = 1000 + in // vaxt aparan hesablama
    }()
    return func() (int64, sync.Mutex) {
        mux.Lock()        // nəticə hazır olana qədər bloklanır
        defer mux.Unlock()
        return result
    }
}
```

**Channel versiyası (idiomatik Go):**
```go
func doWorkFuture(in int64) <-chan int64 {
    result := make(chan int64, 1) // buffer 1: göndərən bloklanmır
    go doTimeConsumingWork(result, in)
    return result
}
future := doWorkFuture(123)
value := <-future // hazır olanda oxu
```

**Ehtiyatlı ol:** (1) channel-də dəyər yalnız BİR dəfə oxunur; (2) funksiya
versiyasında mutex unudulsa deadlock; (3) Go-da bu pattern az populyardır —
əksər hallada birbaşa channel ötürmək kifayətdir.

### 8. Curious Struct Tricks

**Empty struct `struct{}`:** 0 bayt yaddaş. 2 klassik istifadə:
1. **Set tipli map:** `map[string]struct{}` — StringCollection.Match O(n²)
   nested loop-dan O(n) map lookup-a düşür; dəyər yeriynə keys saxlanılır
2. **Semaphore/shutdown signal:** `chan struct{}` — "dəyərin özü əhəmiyyətsiz"
   mənasını daşıyır; `close(shutdownCh)` bütün worker-lərə siqnal verir

**Anonymous struct:** Adsız, inline struct — template data üçün ideal:
```go
data := struct {
    Name     string
    Days     int
    Messages int
}{Name: "Jo", Days: 3, Messages: 7}
```
`html/template` / `text/template` bu struct-ları gözləyir; başqa yerdə lazım
olmadığından ad verməyə dəyməz.

**noCopy (kopyalanma qadağası):** `sync.Locker` implement edən struct-lar dəyərlə
ötürülə bilməz — `go vet` bunu yoxlayır. Öz struct-unu qorumaq üçün:
```go
type MyStruct struct {
    noCopy noCopy // embed — go vet kopyalanmanı yaxalayır
}
type noCopy struct{}
func (noCopy) Lock()   {}
func (noCopy) Unlock() {}
```
Kopyalama yaddaş + GC yükü və gözlənilməz nəticələr (orijinal dəyişsə, kopya
dəyişmir) verir.

## Əsas terminlər

- Pure Function (təmiz funksiya)
- Referential Transparency (istinad şəffaflığı)
- Higher-Order Function (yüksək dərəcəli funksiya)
- Currying (karrinq) / Partial Function (qismən funksiya)
- Closure (bağlanma)
- Middleware (ara qat)
- Functional Options (funksional seçənəklər)
- Future (gələcək dəyər proxy-si)
- Semaphore (semofor)

## Praktik nəticə

- Funksiyalar data kimi işlənə bilər: Map/filter pipeline-lər, middleware,
  DI üçün function var-lar
- Çoxlu optional konfiq → functional options; tələb olunan → parametr
- Set üçün `map[T]struct{}`; siqnal üçün `chan struct{}`
- Futures üçün buffered channel (size 1) — Go-nun idiomatik yolu
- `go vet` + noCopy ilə dəyərlə kopyalanma qadağan et

## Mənbə

Pages: 277-316 (Chapter 8, Beyond Effective Go Part 2)
