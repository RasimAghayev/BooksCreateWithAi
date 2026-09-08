# Chapter 14 — Улучшения производительности (Performans iyileşmələri, məsləhətlər)

## Bu chapter nədən bəhs edir?

pprof profiling aləti, benchmark testləri (RWMutex vs atomic), yaddaş
allocasiyası (concat vs strings.Join) və fasthttp/fasthttprouter.

## Əsas fikirlər

### 1. pprof profiling aləti
**Nədir:** Runtime profiling datası toplayan/ixrac edən Go aləti — CPU,
yaddaş, bloklanma profilləri; veb endpoint-ləri var.

**Necə işləyir:** `_ "net/http/pprof"` importu + portu dinləmək kifayətdir —
/debug/pprof/* endpoint-ləri avtomatik yaranır.

**Kitabdan kod nümunəsi:**
```go
// Ağır handler (bcrypt) — profiling üçün nümunə:
func GuessHandler(w http.ResponseWriter, r *http.Request) {
    if err := r.ParseForm(); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return
    }
    msg := r.FormValue("message")
    real := []byte("$2a$10$2ovnPWuIjMx2S0HvCxP/mutzdsGhyt8rq/...")
    if err := bcrypt.CompareHashAndPassword(real, []byte(msg)); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("try again"))
        return
    }
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("you got it"))
}

// main — yalnız 2 sətrlik əlavə:
import _ "net/http/pprof"     // side-effect import → endpoint-lər yaranır
http.HandleFunc("/guess", crypto.GuessHandler)
http.ListenAndServe("localhost:8080", nil)
```
```bash
# 30 saniyəlik CPU profili (bu müddətdə trafik yarat!):
go tool pprof http://localhost:8080/debug/pprof/profile

# pprof konsolunda:
(pprof) top 10      # ən çox CPU yeyən 10 funksiya
# 870ms 93.55% blowfish.encryptBlock    ← İŞIQLI DÜYÜN
# 30ms  3.23%  blowfish.ExpandKey
(pprof) web         # Graphviz vizuallaşdırma (qırmızı zolaqlar = hot path)

# Yaddaş profili:
go tool pprof http://localhost:8080/debug/pprof/heap
```

**Sub-kod izahı:**
- `top10` → flat (öz müddəti) və cum (uşaqlarla birgə) vaxtlar
- Bütün vaxtın ~94%-i bir funksiyada → dar boğaz aydındır
- Profili TopLADIĞIN müddətdə real trafik yaratmaq şərtdir

### 2. Benchmark — RWMutex vs atomic
**Nədir:** `go test -bench` ilə funksiya sürəti; `RunParallel` ilə
konkurrent yükdə müqayisə.

**Kitabdan kod nümunəsi:**
```go
// Həll 1 — RWMutex sayğacı:
type Counter struct {
    value int64
    mu    *sync.RWMutex
}
func (c *Counter) Add(amount int64) {
    c.mu.Lock()
    c.value += amount
    c.mu.Unlock()
}
func (c *Counter) Read() int64 {
    c.mu.RLock()
    defer c.mu.RUnlock()
    return c.value
}

// Həll 2 — atomic sayğacı:
type AtomicCounter struct{ value int64 }
func (c *AtomicCounter) Add(amount int64) {
    atomic.AddInt64(&c.value, amount)    // kilidsiz atomik
}
func (c *AtomicCounter) Read() int64 {
    return atomic.LoadInt64(&c.value)
}

// Benchmark testləri:
func BenchmarkCounterAdd(b *testing.B) {
    c := Counter{0, &sync.RWMutex{}}
    for n := 0; n < b.N; n++ {      // b.N — framework təyin edir
        c.Add(1)
    }
}
// PARALEL benchmark — lock rəqabətini göstərir:
func BenchmarkCounterAddRead(b *testing.B) {
    c := Counter{0, &sync.RWMutex{}}
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {              // hər goroutine üçün sayğac
            c.Add(1)
            c.Read()
        }
    })
}
```
```bash
go test -bench .
# BenchmarkAtomicCounterAdd-4       8.38 ns/op    ← atomic ~4x daha sürətli
# BenchmarkAtomicCounterRead-4      2.09 ns/op
# BenchmarkAtomicCounterAddRead-4  24.5 ns/op     ← paraleldə fərq DAHA BÖYÜK
# BenchmarkCounterAdd-4           34.8 ns/op
# BenchmarkCounterRead-4           66.0 ns/op
# BenchmarkCounterAddRead-4        146 ns/op     ← RWMutex rəqabətdə 6x yavaş
```

**Sub-kod izahı:**
- `b.N` döngəsi → framework statistik dəqiqlik üçün təkrarlayır
- `b.RunParallel` → çox goroutine yaranır; lock contention ölçülür
- Nəticə: sadə sayğac üçün atomic RWMutex-dən əhəmiyyətli sürətlidir

### 3. Yaddaş allocasiyası — concat vs strings.Join
**Nədir:** `-benchmem` flag-i allocasiya sayı/həcmi göstərir; string
birləşdirmədə iki üsulun müqayisəsi.

**Kitabdan kod nümunəsi:**
```go
// Üsul 1 — tsikl + konkatennasiya (hər += yeni string YARADIR):
func concat(vals ...string) string {
    finalVal := ""
    for i := 0; i < len(vals); i++ {
        finalVal += vals[i]
        if i != len(vals)-1 {
            finalVal += " "
        }
    }
    return finalVal
}

// Üsul 2 — strings.Join (bir dəfə buffer):
func join(vals ...string) string {
    return strings.Join(vals, " ")
}

// Sub-benchmark-lar — fərqli giriş ölçüləri:
func Benchmark_join(b *testing.B) {
    b.Run("one", func(b *testing.B) { join("1") })
    b.Run("five", func(b *testing.B) { join("1","2","3","4","5") })
    b.Run("ten", func(b *testing.B) { join("1",...,"10") })
}
```
```bash
GOMAXPROCS=1 go test -bench=. -benchmem -benchtime=1s
# concat/one    13.6 ns/op    0 B/op   0 allocs
# concat/five  386 ns/op     48 B/op   8 allocs    ← hər += = yeni alloc!
# concat/ten   992 ns/op    256 B/op  18 allocs
# join/one      6.30 ns/op    0 B/op   0 allocs
# join/five   124 ns/op     32 B/op   2 allocs     ← 5x daha az alloc
# join/ten    183 ns/op     64 B/op   2 allocs     ← 5x sürətli, 9x az alloc
```

**Sub-kod izahı:**
- `-benchmem` → B/op (ayrılmış bayt) + allocs/op (ayırma sayı)
- Giriş ölçüsü artdıqca fərq BÖYÜYÜR — sub-benchmark-larla müxtəlif
  ölçüləri sına
- Strings dəyişməzdir: `+=` hər addımda yeni string allocasiyası edir
- Müəllif tövsiyəsi: standart kitabxana funksiyası (Join) adətən optimaldır

### 4. fasthttp + fasthttprouter
**Nədir:** net/http alternativi — hot path optimallaşdırması; router
metod-spesifik route-larla. MƏHDUDİYYƏT: HTTP/2 dəstəyi YOXDUR.

**Kitabdan kod nümunəsi:**
```go
// Thread-safe items (RWMutex):
var items []string
var mu *sync.RWMutex
func AddItem(item string) {
    mu.Lock()
    items = append(items, item)
    mu.Unlock()
}
func ReadItems() []string {
    mu.RLock()
    defer mu.RUnlock()
    return items
}

// fasthttp handler — RequestCtx (ResponseWriter+Request birgə!):
func GetItems(ctx *fasthttp.RequestCtx) {
    enc := json.NewEncoder(ctx)        // ctx = io.Writer
    items := ReadItems()
    enc.Encode(&items)
    ctx.SetStatusCode(fasthttp.StatusOK)
}
func AddItems(ctx *fasthttp.RequestCtx) {
    item, ok := ctx.UserValue("item").(string)   // router parametri!
    if !ok {
        ctx.SetStatusCode(fasthttp.StatusBadRequest)
        return
    }
    AddItem(item)
    ctx.SetStatusCode(fasthttp.StatusOK)
}

// Metod-spesifik route-lar + path parametri:
router := fasthttprouter.New()
router.GET("/item", GetItems)
router.POST("/item/:item", AddItems)      // :item → UserValue("item")
log.Fatal(fasthttp.ListenAndServe("localhost:8080", router.Handler))
```

**Sub-kod izahı:**
- `*fasthttp.RequestCtx` → həm sorğu həm cavab; UserValue ilə route params
- `router.POST("/item/:item", ...)` → :item avtomatik çıxarılır (net/http
  DefaultServeMux-də YOXDUR)
- Handler-lər köçürülməlidir (RequestWriter→RequestCtx) — fasthttp-nin
  qiyməti; HTTP/2 lazımdırsa net/http (1.6+) istifadə et

## Əsas terminlər

- pprof / net/http/pprof (blank import)
- CPU/Heap Profile / top10 / web (Graphviz)
- Hot Path (isti yol — ən çox işlənən kod)
- Benchmark / b.N / RunParallel
- Lock Contention (kilid rəqabəti)
- RWMutex vs atomic (performans müqayisəsi)
- -benchmem / B/op / allocs/op
- String Concatenation vs strings.Join
- fasthttp / fasthttprouter / RequestCtx / UserValue
- Path Parameter (`:item`)
- HTTP/2 (fasthttp-da YOXDUR)

## Praktik nəticə

- pprof: 1 import + port; 30s CPU profili zamanı trafik yarat; top10 →
  dar boğazı tap
- Sərbəst sayğaclarda atomic seç; b.RunParallel rəqabəti real göstərir
- Yaddaş: `+=` konkatennasiyadan qaç; strings.Join/bytes.Buffer istifadə et;
  fərqli giriş ölçülərini sub-benchmark-la ölç
- fasthttp: RequestCtx + metod-spesifik route + :param — sürət qazanmaq
  üçün HTTP/2-dən imtina; şübhə halda net/http kifayətdir

## Mənbə

Pages: 446-464 (Chapter 14, Go Programming Cookbook 2nd ed)
