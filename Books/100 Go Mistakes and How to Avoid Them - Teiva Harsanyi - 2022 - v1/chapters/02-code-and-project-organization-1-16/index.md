# Chapter 2 — Code and project organization (#1-#16)

## Bu chapter nədən bəhs edir?

Bu chapter kodun və layihənin **təşkilati səhvlərini** əhatə edir: variable shadowing, nested code, init funksiyalarının sui-istifadəsi, getter/setter, interface pollution, interface-in yerləşməsi (producer/consumer), `any` tipləri, generics, type embedding, functional options pattern, layihə strukturu, utility paketlər, paket adı toqquşmaları, dokumentasiya və linter-lər.

---

## Əsas fikirlər

### #1: Unintended variable shadowing (Qeyri-ixtiyari dəyişən kölgələnməsi)

**Nədir?** Daxili blokda eyni adlı dəyişənin yenidən elan edilməsi (`:=` ilə) — xarici dəyişəni YOX, yeni lokal dəyişən yaradır.

```go
var client *http.Client
if tracing {
    client, err := createClientWithTracing()  // client KÖLGƏLƏNİR (yeni dəyişən)
    if err != nil { return err }
    log.Println(client)
} else {
    client, err := createDefaultClient()      // yenidən kölgələnir
    if err != nil { return err }
    log.Println(client)
}
// Use client  →  XARİCİ client HƏMİŞƏ nil!
```

**Niyə təhlükəlidir?** Kod kompilyasiya olunur (kölgə dəyişənləri log-da istifadə olunur), amma funksiya nəticəsi gözlənilən dəyişənə düşmür.

**Həll 1 — müvəqqəti dəyişən:**

```go
if tracing {
    c, err := createClientWithTracing()
    if err != nil { return err }
    client = c            // xarici dəyişənə mənimsədilir
}
```

**Həll 2 — assignment operatoru (`=`) + err dəyişəni:**

```go
var client *http.Client
var err error
if tracing {
    client, err = createClientWithTracing()   // = kölgə YARATMIR
} else {
    client, err = createDefaultClient()
}
if err != nil {
    // Common error handling — error idarəsi birləşdirilir
}
```

**Sub-kod izahı:**
- `:=` → yeni dəyişən elan edir (kölgə yaradır)
- `=` → mövcud dəyişənə mənimsədir (yalnız əvvəlcədən elan olunmuşsa işləyir)
- `var err error` → `=` istifadəsi üçün err-in əvvəlcədən elanı

**Qeyd:** Shadowing-i aşkarlamaq üçün `vet` + `shadow` linter-i mövcuddur (#16).

---

### #2: Unnecessary nested code (Lazımsız iç-içə kod)

**Problem:** Mental model (kodun daxili təmsili) qurmaq üçün koqnitiv səy nested səviyyələrin sayı ilə artır.

```go
// PİS — 5 nested səviyyə:
func join(s1, s2 string, max int) (string, error) {
    if s1 == "" {
        return "", errors.New("s1 is empty")
    } else {
        if s2 == "" {
            return "", errors.New("s2 is empty")
        } else {
            concat, err := concatenate(s1, s2)
            if err != nil {
                return "", err
            } else {
                if len(concat) > max {
                    return concat[:max], nil
                } else {
                    return concat, nil
                }
            }
        }
    }
}
```

```go
// YAXŞI — 2 səviyyə, happy path solda:
func join(s1, s2 string, max int) (string, error) {
    if s1 == "" {
        return "", errors.New("s1 is empty")
    }
    if s2 == "" {
        return "", errors.New("s2 is empty")
    }
    concat, err := concatenate(s1, s2)
    if err != nil {
        return "", err
    }
    if len(concat) > max {
        return concat[:max], nil
    }
    return concat, nil
}
```

**Qaydalar:**
1. `if` bloku return edirsə → `else`-i burax.
2. Non-happy path-i yoxlayıb **erkən return** et; şərti çevir:

```go
// PİS: if s != "" { ... } else { return errors.New("empty string") }
// YAXŞI:
if s == "" {
    return errors.New("empty string")
}
// ...
```

**Mat Ryer (Go Time podcast):** *"Align the happy path to the left; you should quickly be able to scan down one column to see the expected execution flow."* — Happy path sol sütündə, edge case-lər ikinci sütunda oxunur.

---

### #3: Misusing init functions (init funksiyalarının sui-istifadəsi)

**Init funksiyası nədir?** Arqumentsiz, nəticəsiz `func()` — paket state-inin ilkinləşdirilməsi üçün.

**İcra sırası:**
1. Paketin const/var elanları qiymətləndirilir
2. İmport olunan paketlərin init-ləri (asılılıq sırasıyla — `redis` paketi `main`-dən əvvəl)
3. Bu paketin init funksiyaları (multi-init: faylların əlifba sırası — buna GÜVƏNMƏ)
4. `main()`

Aynı faylda birdən çox init mümkündür (mənbə sırası ilə). `_ "foo"` importu — yalnız side-effect üçün init-i işə salır.

**Məhdudiyyət:** init birbaşa çağırıla bilməz: `undefined: init`.

**Pis nümunə — DB connection pool init-də:**

```go
var db *sql.DB
func init() {
    d, err := sql.Open("mysql", os.Getenv("MYSQL_DATA_SOURCE_NAME"))
    if err != nil { log.Panic(err) }   // 1. yalnız panic mümkün!
    if err = d.Ping(); err != nil { log.Panic(err) }
    db = d                            // 3. qlobal dəyişən məcburiyyəti
}
```

**3 çatışmazlıq:**
1. **Error idarəsi məhduddur** — error qaytara bilmir, yalnız panic; çağıran paket retry/fallback edə bilmir.
2. **Testləri çətinləşdirir** — init testlərdən əvvəl icra olunur (unit test DB bağlantısı istəmirsə belə).
3. **Qlobal dəyişən tələb edir** — hər funksiya dəyişəni dəyişə bilər; unit test izolyasiyası pozulur.

**Düzgün həll — adi funksiya:**

```go
func createClient(dsn string) (*sql.DB, error) {
    db, err := sql.Open("mysql", dsn)
    if err != nil { return nil, err }
    if err = db.Ping(); err != nil { return nil, err }
    return db, nil
}
```

- Error idarəsi çağırana qalır · integration test mümkündür · bağlantı encapsulate olunur.

**Init-in YAXŞI olduğu hal:** Uğura bilən statik konfiqurasiya (Go rəsmi blogu misalı — HTTP handler qeydiyyatı: error mümkün deyil, qlobal state yoxdur, testlərə təsir yoxdur).

---

### #4: Overusing getters and setters

**Go-da getter/setter məcburi deyil** — standart kitabxana belə birbaşa sahə çıxışı verir: `timer := time.NewTimer(time.Second); <-timer.C`.

**Getter/setter-in düzgün olduğu hallar:**
- Sahəyə davranış bağlamaq (validasiya, hesablanmış dəyər, mutex wrap)
- Daxili təmsilçiliyi gizlətmək
- Run-time dəyişikliyi üçün debugging nöqtəsi

**Adlandırma konvensiyası:**

```go
currentBalance := customer.Balance()   // getter: Balance (GetBalance YOX)
if currentBalance < 0 {
    customer.SetBalance(0)              // setter: SetBalance
}
```

**Prinsip:** Dəyər gətirməyən getter/setter-lərlə kodu doldurma — praqmatik ol.

---

### #5: Interface pollution (Interface çirklənməsi)

**Interface nədir?** Obyektin davranışını təyin edən abstraksiya. Go-da **implicit (gizli)** təmin olunur — `implements` keyword yoxdur.

**io.Reader/io.Writer nümunəsi:**

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}
type Writer interface {
    Write(p []byte) (n int, err error)
}
```

Abstraksiyanın gücü — generic funksiya:

```go
func copySourceToDest(source io.Reader, dest io.Writer) error { /* ... */ }
```

`*os.File` hər ikisini implement edir; testdə fayl YOX — `strings.NewReader` + `bytes.NewBuffer`:

```go
func TestCopySourceToDest(t *testing.T) {
    const input = "foo"
    source := strings.NewReader(input)          // io.Reader
    dest := bytes.NewBuffer(make([]byte, 0))    // io.Writer
    err := copySourceToDest(source, dest)
    // assertion-lar...
}
```

**Rob Pike:** *"The bigger the interface, the weaker the abstraction."* — io.ReadWriter kimi kombinasiya mümkündür:

```go
type ReadWriter interface {
    Reader
    Writer
}
```

**Interface-in 3 düzgün istifadə yeri:**

1. **Common behavior (Ümumi davranış)** — `sort.Interface{Len, Less, Swap}` — hər hansı index-based kolleksiyanı sort edən abstraksiya; `sort.IsSorted` kimi utility-lər də eyni abstraksiyadan istifadə edir.
2. **Decoupling (Qoparma)** — konkret `mysql.Store` asılılığı əvəzinə `customerStorer` interface: unit test mock, integration test konkret tip — hər ikisi mümkün. (Liskov Substitution Principle — SOLID-in L-i.)
3. **Restricting behavior (Davranışı məhdudlaşdırma)** — `IntConfig` həm `Get` həm `Set` daşıyır, amma kod yalnız oxumaq istəyirsə: `intConfigGetter{Get() int}` interface inject et — semantik read-only təmin olunur.

**Interface pollution:** C#/Java fonlu developer-lərin konkret tipdən əvvəl interface yaratmaq vərdişi. Go-da **"abstractions should be discovered, not created"** — interface gələcək ehtiyac üçün deyil, mövcud ehtiyac üçün yaradılır. Lazımsız interface = lazımsız indirection → kod axını mürəkkəbləşir (+ hash table lookup overhead, əksər kontekstdə minimal). **Rob Pike: "Don't design with interfaces, discover them."**

---

### #6: Interface on the producer side (Producer tərəfində interface)

**Terminologiya:**
- **Producer side** — interface konkret implementasiya ilə EYNİ paketdə
- **Consumer side** — interface istifadə olunduğu XARİCİ paketdə

**Pis yanaşma:** `store` paketi öz `CustomerStorage` interface-ini (6 metodlu) export edir — bütün client-lərə tək abstraksiya MƏCBURİ edilir.

**Go-nun gücü — implicit satisfaction:** Client ÖZ paketində lazım olan minimal interface yaradır:

```go
// client paketində (unexported — yalnız burada istifadə olunur):
type customersGetter interface {
    GetAllCustomers() ([]store.Customer, error)
}
```

**Niyə bu mümkündür?** Interface gizli təmin olunduğundan `store` paketinin `client`-ə asılılığı YOXDUR (görüntüdə circular dependency bənzərir, amma deyil). Bu — Interface Segregation Principle (SOLID-in I-si): client istifadə etmədiyi metodlara asılı olmaq məcburiyyətində deyil.

**İstisna:** Standart kitabxanada `encoding` paketi producer-side interface-lər müəyyən edir — dil dizaynerləri abstraksiyanın qabaqcadan dəyərli olacağını BİLİRDİLƏR (sonradan yox). Belə hallar istisnadır; interface producer tərəfində olacaqsa minimal saxlanmalıdır.

---

### #7: Returning interfaces (Interface qaytarma)

**Problem:** `store` paketi `NewInMemoryStore()` funksiyasından `client.Store` interface qaytarırsa → **store paketi client paketinə asılı olur** → client `NewInMemoryStore` çağıra bilməz (cyclic dependency) → digər paketdən inject etmək məcburiyyəti → **code smell**.

**Postel's law (TCP, RFC 761):** *"Be conservative in what you do, be liberal in what you accept from others."* Go-ya tətbiqi:
- **Qaytar: struct (konkret tip)** · **Qəbul et: interface (mümkünsə)**

**İstisnalar:** `error` (özü interface); `io.LimitReader` → `io.Reader` qaytarır — çünki `io.Reader` up-front abstraksiyadır, dili dizaynerləri reusability/composability üçün qabaqcadan biliblər.

---

### #8: any says nothing (any heç nə demir)

**Go 1.18:** `any` = `interface{}` alias. `any` hər tip saxlaya bilər, amma **bütün tip məlumatı itir** — type assertion tələb olunur.

```go
// Pis — ifadəlilik sıfır, compile-time qoruma yoxdur:
func (s *Store) Get(id string) (any, error) {}
func (s *Store) Set(id string, v any) error {}
s := store.Store{}
s.Set("foo", 42)   // int göndərilir — heç nə bunu dayandırmır!

// Yaxşı — tipə xas metodlar:
func (s *Store) GetContract(id string) (Contract, error) {}
func (s *Store) SetContract(id string, contract Contract) error {}
func (s *Store) GetCustomer(id string) (Customer, error) {}
func (s *Store) SetCustomer(id string, customer Customer) error {}
```

Client öz abstraction-u yarada bilər: `ContractStorer` interface — yalnız Contract metodları.

**any-nin düzgün olduğu hallar (standart kitabxana):**
- `json.Marshal(v any)` — hər hansı tip marshall oluna bilər
- `sql.Conn.QueryContext(ctx, query, args ...any)` — parametrlər hər tip ola bilər

**Prinsip:** `any` yalnız real "hər tip" ehtiyacında. Dublikasiya kodu bəzən expressivlik üçün daha yaxşıdır.

---

### #9: Being confused about when to use generics (Generics nə vaxt?)

**Problem:** `getKeys` funksiyası `map[string]int` üçün yazılıb; `map[int]string` üçün nə etməli? Generics-dən əvvəl: code generation, reflection, dublikasiya.

**Reflection/any yanaşmasının problemləri:** boilerplate artır, tip yoxlaması run-time-a keçir (error qaytarmaq məcburiyyəti), `[]any` qaytarmağa məcbur edir.

**Type parameters (Go 1.18):**

```go
func getKeys[K comparable, V any](m map[K]V) []K {
    var keys []K
    for k := range m {
        keys = append(keys, k)
    }
    return keys
}
```

- `K comparable` — map key `==`/`!=` tələb edir (slice key ola bilməz: `invalid map key type []byte`)
- `V any` — dəyər hər tip
- Instantiation compile-time — type safety qorunur, run-time overhead yoxdur

**Custom constraint:**

```go
type customConstraint interface {
    ~int | ~string     // union operator
}
func getKeys[K customConstraint, V any](m map[K]V) []K { /* ... */ }
```

**~int vs int:** `int` yalnız həmin tipi; `~int` — underlying type-i int olan bütün tiplər (`type customInt int` kimi). `~int` + `String() string` constraint — customInt qəbul edir, `int` + `String()` isə YOX (int String implement etmir).

**Generics data strukturda:**

```go
type Node[T any] struct {
    Val  T
    next *Node[T]
}
func (n *Node[T]) Add(next *Node[T]) { n.next = next }   // receiver instantiate olunur
```

**Məhdudiyyət:** Metodlar type parameter QƏBUL EDƏ BİLMƏZ (yalnız funksiya arqumenti və ya receiver): `func (Foo) bar[T any](t T)` → `methods cannot have type parameters`.

**Generics DÜZGÜN olduğu yerlər:**
1. **Data strukturlar** — binary tree, linked list, heap — element tipini factor out
2. **Slices/maps/channels üzərində funksiyalar** — `func merge[T any](ch1, ch2 <-chan T) <-chan T`
3. **Davranışı factor out** — `SliceFn[T any]{S []T, Compare func(T,T) bool}` implement `sort.Interface` → hər tip üçün ayrı funksiya YOX:

```go
type SliceFn[T any] struct {
    S       []T
    Compare func(T, T) bool
}
func (s SliceFn[T]) Len() int           { return len(s.S) }
func (s SliceFn[T]) Less(i, j int) bool { return s.Compare(s.S[i], s.S[j]) }
func (s SliceFn[T]) Swap(i, j int)     { s.S[i], s.S[j] = s.S[j], s.S[i] }

s := SliceFn[int]{S: []int{3, 2, 1}, Compare: func(a, b int) bool { return a < b }}
sort.Sort(s)   // [1 2 3]
```

**Generics YANLIŞ olduğu yerlər:**
1. **Type argument-in metodunu çağıranda** — `func foo[T io.Writer](w T)` → `w.Write(b)` — heç bir dəyər yoxdur, birbaşa `io.Writer` qəbul et.
2. **Kodu mürəkkəbləşdirəndə** — generics məcburi deyil; aydınlıq gətirmirsə yenidən düşün.

**Prinsip:** Interface-lə eyni — abstraksiya kimi generics də vaxtından əvvəl istifadə olunmamalı. Boilerplate yazmaq üzrəykən düşün.

---

### #10: Not being aware of the possible problems with type embedding (Type embedding problemləri)

**Embedded field:** Adsız elan edilmiş struct sahəsi — sahə/metodlar promote olunur:

```go
type Foo struct {
    Bar        // embedded
}
type Bar struct {
    Baz int
}
foo := Foo{}
foo.Baz = 42           // promoted path
foo.Bar.Baz = 42       // nominal path — eyni sahə
```

**Pis nümunə — sync.Mutex embed:**

```go
type InMem struct {
    sync.Mutex          // EMBED — təhlükə!
    m map[string]int
}
func (i *InMem) Get(key string) (int, bool) {
    i.Lock()            // promote edilmiş metod
    v, contains := i.m[key]
    i.Unlock()
    return v, contains
}
// XARİCİ client:
m := inmem.New()
m.Lock()   // ?? — Lock/Unlock client-ə görünür!
```

**Həll:** Adlı unexported sahə istifadə et:

```go
type InMem struct {
    mu sync.Mutex      // embed YOX — xaricdən görünmür
    m  map[string]int
}
```

**Yaxşı nümunə — io.WriteCloser embed:**

```go
type Logger struct {
    io.WriteCloser     // embed — forwarding metodları lazım deyil
}
func main() {
    l := Logger{WriteCloser: os.Stdout}
    l.Write([]byte("foo"))   // promote
    l.Close()
}
```

Write/Close promote olunur → Logger `io.WriteCloser`-ı SATISFY edir → forwarding metodları yazmaq lazım deyil.

**Embedding vs OOP subclassing:** Embedding-də metodun receiver-i embedded tip (X) qalır; subclassing-də subclass (Y) olur. **Embedding = kompozisiya, irs deyil.**

**2 məhdudiyyət:**
1. Yalnız sintaktik şəkər kirsə (`Foo.Baz` vs `Foo.Bar.Baz`) — embed ETMƏ.
2. Gizləmək istədiyin data/metodları promote edirsə — embed ETMƏ. (Public struct-larda embed → daxili tip yeni metod əlavə etsə diqqət tələb edir.)

---

### #11: Not using the functional options pattern

**Ssenari:** `NewServer(addr, port)` — client-lər write timeout və s. istəyir; yeni parametr əlavə etmək compatibility qırır. Port loqikası: yoxdursa default, mənfidirsə error, 0-dırsa random, yoxsa verilən.

**Yanaşma 1 — Config struct:** Uyğunluq problemi həll, amma **zero value ambiguity**: `Config{Port: 0}` == `Config{}` — "0 qəsdən verilib" ilə "verilməyib" fərqlənmir. Pointer istifadəsi (`Port *int`) fərqi verir, amma client üçün qeyri-praktik + default üçün boş struct məcburiyyəti: `NewServer("localhost", httplib.Config{})`.

**Yanaşma 2 — Builder pattern (GoF):**

```go
type ConfigBuilder struct{ port *int }
func (b *ConfigBuilder) Port(port int) *ConfigBuilder {  // chain üçün builder qaytarır
    b.port = &port
    return b
}
func (b *ConfigBuilder) Build() (Config, error) {
    if b.port == nil { /* default */ } else if *b.port == 0 { /* random */ }
    // mənfi → error...
}
```

Çatışmazlıqlar: default üçün yenə boş config; chaining error qaytarmağa imkan vermir → validasiya Build-ə gecikir.

**Yanaşma 3 — Functional options pattern (idiomatik):**

```go
type options struct {
    port *int
}
type Option func(options *options) error

func WithPort(port int) Option {
    return func(options *options) error {      // closure — port-u referans edir
        if port < 0 {
            return errors.New("port should be positive")
        }
        options.port = &port
        return nil
    }
}

func NewServer(addr string, opts ...Option) (*http.Server, error) {
    var options options
    for _, opt := range opts {
        if err := opt(&options); err != nil {
            return nil, err
        }
    }
    var port int
    if options.port == nil {
        port = defaultHTTPPort
    } else if *options.port == 0 {
        port = randomPort()
    } else {
        port = *options.port
    }
    // ...
}
```

**Client istifadəsi:**

```go
server, err := httplib.NewServer("localhost",
    httplib.WithPort(8080),
    httplib.WithTimeout(time.Second))

// Default üçün:
server, err := httplib.NewServer("localhost")   // boş arqument lazım deyil!
```

**Elementlər:** unexported `options` struct · `Option func(*options) error` funksiya tipi · hər seçim üçün `With` prefiksli funksiya (closure + validasiya) · variadic `opts ...Option` · iterasiya ilə apply. gRPC kimi kitabxanalarda istifadə olunur.

---

### #12: Project misorganization (Layihənin yanlış təşkili)

**project-layout (github.com/golang-standards/project-layout):**

| Qovluq | Məzmun |
|--------|--------|
| `/cmd` | Əsas source faylları — `/cmd/foo/main.go` |
| `/internal` | Xaricin import etməməli olduğu privat kod |
| `/pkg` | Xaricə açıq publik kod |
| `/test` | Xarici testlər + test datası (unit testlər source ilə eyni paketdə) |
| `/configs` | Konfiqurasiya faylları |
| `/docs` | Dizayn və istifadəçi sənədləri |
| `/examples` | Nümunələr |
| `/api` | API müqavilə faylları (Swagger, Protocol Buffers) |
| `/web` | Veb assetləri |
| `/build` | Paketləmə və CI |
| `/scripts` | Skriptlər |
| `/vendor` | Asılılıqlar |

**Qeydlər:**
- `/src` YOXDUR — çox genericdir.
- **Russ Cox (2021):** bu layout rəsmi standart deyil (golang-standards təşkilatı Go-nun deyil). Məcburi konvensiya yoxdur — **"indecision is the only wrong decision"** — təşkilatıstandartlaşdır, developer-lər repo-lar arasında vaxt itirməsin.

**Paket təşkili:** Go-da subpackage anlayışı yoxdur — `net/http` `net`-dən miras almır, yalnız exported elementləri görür. Subdirektoriyaların faydası: high cohesion.

**Paketlərin təşkili qaydaları:**
1. **Premature packaging-dən qaç** — sadə başla, layihə inkişaf etdikcə böyüt.
2. **Nano packages-dən qaç** (1-2 fayllıqlar) — məntiqi əlaqələr itir; converse: nəhəng paketlər də pis — ad mənasını itirir.
3. **Adlandırma:** paketi təmin etdiyi şeyə görə adlandır (ehtiva etdiyinə görə YOX); qısa, ifadəli, tək lowercase söz.
4. **Export minimallaşdır** — şübhələnirsə export ETMƏ; sonradan açmaq asandır (exception: `encoding/json` unmarshal üçün export sahələr).
5. Kontekstə görə yoxsa layərə görə qruplaşdırma — ikisi də ola bilər, vacibi **consistency**.

---

### #13: Creating utility packages (Utility paketləri)

**Problem:** `util`, `common`, `shared`, `base` — mənasız adlar, paketin nə təmin etdiyini söyləmir.

```go
// PİS:
package util
func NewStringSet(...string) map[string]struct{} {}
func SortStringSet(map[string]struct{}) []string {}

set := util.NewStringSet("c", "a", "b")
fmt.Println(util.SortStringSet(set))
```

**Həll — ifadəli paket adı (Go blogundan ilhamlanmış set nümunəsi):**

```go
package stringset
func New(...string) map[string]struct{} { ... }
func Sort(map[string]struct{}) []string { ... }

set := stringset.New("c", "a", "b")
fmt.Println(stringset.Sort(set))
```

**Bir addım irəli — konkret tip:**

```go
type Set map[string]struct{}
func New(...string) Set { ... }
func (s Set) Sort() []string { ... }

set := stringset.New("c", "a", "b")
fmt.Println(set.Sort())    // yalnız 1 dəfə stringset istinadı
```

**Qeyd:** Nano package anlayışı özü pis deyil — high cohesion varsa kiçik paket qəbul olunandır. Ümumi tiplər (client/server arası) üçün Dave Cheney: bəzən hamısını TƏK paketdə birləşdirmək daha yaxşıdır.

---

### #14: Ignoring package name collisions (Paket adı toqquşmaları)

```go
redis := redis.NewClient()   // redis dəyişəni redis PAKETİNİ kölgələyir!
v, err := redis.Get("foo")   // burada OK, amma paket artıq əlçatmazdır
```

**Həllər:**
1. Fərqli dəyişən adı: `redisClient := redis.NewClient()` (ən sadə).
2. Import alias: `import redisapi "mylib/redis"` → dəyişən `redis` qala bilər.

**Əlavə qorunma:** Built-in funksiya adlarından da qaç — `copy := copyFile(src, dst)` → `copy` built-in-i bloklanır. **Dot import** (qualifier-siz çıxış) — qarışıqlığı artırır, adətən qaçınılmalı.

---

### #15: Missing code documentation (Kod dokumentasiyasının çatışmazlığı)

**Qayda 1 — hər exported element dokumentləşdirilməlidir**, comment elementin adı ilə başlayır:

```go
// Customer is a customer representation.
type Customer struct{}

// ID returns the customer identifier.
func (c Customer) ID() string { return "" }
```

- Comment tam cümlə, durğu işarəsi ilə bitir.
- Funksiya dokumentasiyası NƏ etdiyini yazır (necə yox — o kodun/kommentlərin işi).
- İdeal: klient kodu OXMADAN istifadə edə bilsin.

**Deprecated:** `// Deprecated:` — IDE-lər xəbərdarlıq göstərir:

```go
// ComputePath returns the fastest path between two points.
// Deprecated: This function uses a deprecated way to compute
// the fastest path. Use ComputeFastestPath instead.
func ComputePath() {}
```

**Konstanta/variable:** Məqsəd → code documentation; məzmun → yanlış sətir commenti:

```go
// DefaultPermission is the default permission used by the store engine.
const DefaultPermission = 0o644 // Need read and write accesses.
```

**Paket dokumentasiyası:** `// Package` prefiksi ilə başlayır; ilk sətir qısa olmalı (godoc-da görünür); istənilən fayl (doc.go). Deklarasiyaya bitişik olmayan commentlər (copyright) dokumentasiyaya düşmür — arada boş sətir onu ayırır:

```go
// Copyright 2009 The Go Authors. All rights reserved.  ← DOKUMENTASİYAYA DÜŞMÜR
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Package math provides basic constants and mathematical functions.  ← DÜŞÜR
//
// This package does not guarantee bit-identical results
// across architectures.
package math
```

---

### #16: Not using linters (Linter-lərdən istifadə etməmək)

**Linter** — kodu analiz edib xətaları yakalayan avtomatik alət.

**Nümunə — shadowing-in aşkarlanması (#1 ilə əlaqə):**

```bash
$ go install golang.org/x/tools/go/analysis/passes/shadow/cmd/shadow
$ go vet -vettool=$(which shadow)
./main.go:8:3: declaration of "i" shadows declaration at line 6
```

**Gündəlik linter-lər:**
- `go vet` (https://golang.org/cmd/vet/) — standart Go analayzerı
- `errcheck` (https://github.com/kisielk/errcheck) — error yoxlamaları
- `gocyclo` (https://github.com/fzipp/gocyclo) — cyclomatic complexity
- `goconst` (https://github.com/jgautheron/goconst) — təkrar string konstantlar

**Formatter-lər:**
- `gofmt` (https://golang.org/cmd/gofmt/) — standart formatter
- `goimports` — import formatter

**golangci-lint** (https://github.com/golangci/golangci-lint) — çoxlu linter/formatter üzərində facade; parallel icra. İcra avtomatlaşdırılsın: CI və ya Git precommit hook.

**Qeyd:** Linter-lər kitabdakı 100 səhvin hamısını yakalamır — oxumağa davam et.

---

## Əsas terminlər

| Termin | İzah |
|--------|------|
| Variable shadowing (Dəyişən kölgələnməsi) | Daxili blokda eyni adlı yeni dəyişənin elanı — xaricini kölgələyir |
| Mental model | Sistemin davranışının daxili təmsili — koqnitiv səylə qorunur |
| Happy path | Gözlənilən icra axını — solda düzülməlidir |
| init function | Paket ilkinləşdirmə funksiyası — error qaytara bilmir, panic məcburiyyəti |
| Zero value ambiguity | `Config{}` == `Config{Port: 0}` — "verilməyib" ilə "0 verilib" fərqlənmir |
| Interface pollution | Lazımsız interface-lərlə kodu çirkləndirmək |
| Implicit satisfaction | Go-da interface `implements` olmadan təmin olunur |
| Producer vs consumer side | Interface-in implementasiya yanında YOX, istifadə olunduğu paketdə olması |
| Liskov Substitution Principle | Implementasiyanın abstraksiya ilə əvəz edilə bilməsi (SOLID-L) |
| Interface Segregation Principle | Client istifadə etmədiyi metodlara asılı olmamalıdır (SOLID-I) |
| Postel's law | "Qəbul etməkdə liberal, etməkdə konservativ ol" — qəbul: interface, qaytarma: struct |
| `any` | `interface{}` alias (Go 1.18) — tip məlumatını itirir |
| Type parameter `[T any]` | Generics — compile-time instantiate olunur, type-safe |
| Constraint | Tip arqumentlərini məhdudlaşdıran interface (`comparable`, `~int\|~string`) |
| `~int` (tilde) | Underlying type-i int olan bütün tiplər |
| Type embedding | Adsız sahə — metod/sahə promote; kompozisiya, irs deyil |
| Promotion | Embedded tipin elementlərinin xarici tərəfə görünməsi |
| Functional options pattern | `Option func(*options) error` + `With` prefiksli closure-lar + variadic |
| Closure | Bədənindən kənar dəyişənləri referans edən anonim funksiya |
| Builder pattern | ConfigBuilder + method chaining + `Build()` |
| project-layout | /cmd /internal /pkg /test /api strukturu (qeyri-rəsmi) |
| `/internal` | Xaricin import edə bilmədiyi privat paketlər |
| Nano package | 1-2 fayllıq mikro paket — cohesion yoxdursa pis |
| Import alias | `import redisapi "mylib/redis"` — ad toqquşması həlli |
| `// Deprecated:` | Exported elementin köhnəlmiş olduğunu bildirən komment konvensiyası |
| Linter / Formatter | Kod xətası analizi / stil düzəlişi alətləri (vet, errcheck, gofmt) |

---

## Praktik nəticə

1. **Kölgədən qorun:** `:=` daxili blokda yeni dəyişən yaradır — xarici dəyişənə təyinat lazımdırsa `=` + əvvəlcədən elan (və ya müvəqqəti `c` + `client = c`); `vet`+`shadow` avtomatik yakalayır.
2. **Happy path solda, else-ləri burax, erkən return:** nested səviyyələr = koqnitiv yük.
3. **init yalnız uğura bilən statik konfiqurasiya üçün:** DB/bağlantılar adi `createX(...) error` funksiyalarında — error idarəsi + test + encapsulation qazanılır.
4. **Abstraksiyalar KƏŞF olunur, yaradılmır:** interface/generics yalnız konkret ehtiyacda; `Don't design with interfaces, discover them`.
5. **Interface consumer tərəfdə:** producer konkret tipi export etsin; client ÖZ minimal (unexported) interface-ni yaratsın — implicit satisfaction bunu mümkün edir.
6. **Funksiya imzası:** qaytar **struct**, qəbul et **interface** (Postel's law).
7. **`any` = ifadəsizlik:** tipə xas metodlar yaz; dublikasiya bəzən `any`-dən yaxşıdır (json.Marshal/sql args istisna).
8. **Generics 3 halda:** data strukturlar, slice/map/channel funksiyaları, davranış factor-out; metodlarda type parametr YOX.
9. **Embed-də 2 qadağa:** yalnız sintaktik şəkər üçün YOX; gizlədiləcək elementləri promote edirsə YOX (sync.Mutex → `mu sync.Mutex`).
10. **Optional konfiqurasiya = functional options:** `WithPort(8080)` variadic closure-lar; default = arqumentsiz çağırış.
11. **Layihə strukturu qərar ver və consistent qal:** indecision is the only wrong decision.
12. **`util/common/base` paketlərini refactor et:** təmin edilən şeyə görə adlandır (`stringset`).
13. **Hər exported element dokumentlə:** `// Name does X` formatı; paket sənədi `// Package name ...` — godoc avtomatik generasiya edir.
14. **Linter-ləri avtomatlaşdır:** vet + errcheck + gofmt + golangci-lint CI/precommit-də.

---

## Mənbə

- Kitab: *100 Go Mistakes and How to Avoid Them* — Teiva Harsanyi, Manning, 2022 (ISBN 9781617299599)
- Chapter 2: "Code and project organization", book səh. 7–55
- PDF səhifələri: 27–75
- Xarimi linklər: Mat Ryer "Line of Sight in Code" (medium.com/@matryer); Rob Pike — Proverbs (youtube.com/watch?v=PAAkCSZUG1c); project-layout (github.com/golang-standards/project-layout); golangci-lint (github.com/golangci-lint)
