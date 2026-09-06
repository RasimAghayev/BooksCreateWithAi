# Chapter 7 — Types, Methods, and Interfaces (Tiplər, Metodlar və İnterfeyslər)

## Bu chapter nədən bəhs edir?

İstifadəçi təyinli tiplər, metodlar (pointer/value receiver), nil receiver-lar, method
value/expression, iota enum-ları, embedding (kompozisiya), implicit interface-lər, interface
+ nil semantikası, boş interface, type assertion/switch və dependency injection.

## Əsas fikirlər

### 1. İstifadəçi Təyinli Tiplər
**Nədir:** `type Ad T` — hər hansı tipin üzərində yeni ad.

**Kitabdan kod nümunəsi:**
```go
type Score int
type Converter func(string) Score
type TeamScores map[string]Score
type HighScore Score      // user-defined → user-defined da mümkün
```

**Vacib:** Bu, **inheritance DEYİL** — `HighScore` ilə `Score` eyni underlying type-a
sahib fərqli tiplərdir; bir-birinə təyin compile xətasıdır (yalnız `Score(i)` conversion
keçir). Metodlar da keçmir. **Tiplər = executable documentation**: `Percentage`
parametri `int`-dən aydındır.

### 2. Metodlar və Receiver Qaydaları
**Nədir:** `func (p Person) String() string` — `func` ilə ad arasında receiver.

**Qaydalar (mühüm sıralama):**
1. Metod receiver-i dəyişirsə → **pointer receiver MÜTLƏQ**
2. Metod nil instance işləməlidirsə → **pointer receiver MÜTLƏQ**
3. Dəyişmirsə → value receiver ola bilər
4. Tipdə hər hansı pointer receiver varsa → hamısına pointer (tutarlılıq)

Receiver adı: tipin qısaltması (ilk hərf); `this`/`self` QADAĞAN. Eyni tipdə metod adı
overload OLUNMUR. Metod tip ilə eyni paketdə bəyan olunmalı (böyümək hüququmuz olmayan
tiplərə metod əlavə etmək mümkün deyil).

**Kitabdan kod nümunəsi:**
```go
type Counter struct {
    total       int
    lastUpdated time.Time
}
func (c *Counter) Increment() {   // pointer — dəyişir
    c.total++
    c.lastUpdated = time.Now()
}
func (c Counter) String() string { ... }  // value
```

**Avtomatik address:** `c.Increment()` → `(&c).Increment()`. Amma funksiyaya value
ötürüləndə kopya üzərində çağırılır: `doUpdateWrong(c Counter)` main-in `c`-sini
dəyişməz saxlayır; `doUpdateRight(c *Counter)` — dəyişir.

**Method set:** pointer instance → pointer + value metodlar; value instance → yalnız
value metodlar. (Interface tələblərində kritik.)

**Getter/setter YOX** — sahəyə birbaşa çıxış idiomatikdir; metodlar business logic
üçün. İstisna: çoxsahəli atomik yeniləmə və ya qeyri-düz təyinat (məs. `Increment`).

### 3. Nil Receiver-lər üçün Kodlaşdırma
**Nədir:** Go nil pointer metod çağırışını icra etməyə çalışır — value receiver-də panic,
pointer receiver-də metod nil-i idarə edə bilirsə işləyir.

**Kitabdan kod nümunəsi (binary tree — nil uşaq düyünü normaldır):**
```go
func (it *IntTree) Insert(val int) *IntTree {
    if it == nil {
        return &IntTree{val: val}
    }
    if val < it.val {
        it.left = it.left.Insert(val)
    } else if val > it.val {
        it.right = it.right.Insert(val)
    }
    return it
}
func (it *IntTree) Contains(val int) bool {
    switch {
    case it == nil:
        return false
    ...
}
```

**Məhdudiyyət:** nil pointer-i metoddaxilində non-nil etmək mümkün deyil (pointer-in
kopyası dəyişir). Nil-i qəbul edə bilməyən pointer metod → nil yoxla + error qaytar.

### 4. Metodlar = Funksiyalar (Method Value / Expression)
**Kitabdan kod nümunəsi:**
```go
myAdder := Adder{start: 10}
f1 := myAdder.AddTo     // method value — closure kimi instance sahələrini saxlayır
fmt.Println(f1(10))     // 20
f2 := Adder.AddTo       // method expression — func(Adder, int) int
fmt.Println(f2(myAdder, 15))  // 25
```

**Funksiya yoxsa metod:** məntiq başqa datadan asılıdırsa (config/state) → struct +
metod; yalnız input-lardan asılıdırsa → funksiya.

### 5. iota — Enum-lar (bəzən)
**Nədir:** const blokunda artan dəyərlər verən identifikator (APL-dan gəlib).

**Kitabdan kod nümunəsi:**
```go
type MailCategory int
const (
    Uncategorized MailCategory = iota  // 0
    Personal                           // 1
    Spam                               // 2
    Social                             // 3
    Advertisements                     // 4
)
```

**Qaydalar (Danny van Heumen məsləhəti):** Yalnız **daxili** konstantlar üçün (ad ilə
müraciət olunursa); xarici sistemlərlə/spec ilə dəyər əlaqədardırsa — açıq dəyər yaz
(ortada əlavə konstant sonrakıları renumber edir — səssiz break!). Zero value mənalıdırsa
`iota` 0-dan başlaması faydalı; mənasızdırsa ilk dəyəri `_` və ya "invalid"a ver. Bit
maska pattern (`1 << iota`) ağıllıdır amma sənədləndir.

### 6. Embedding — Kompozisiya (inheritance YOX)
**Nədir:** Sahə adı olmayan embedded struct — onun sahə/metodları xarici struct-a
"promote" olunur.

**Kitabdan kod nümunəsi:**
```go
type Employee struct {
    Name string
    ID   string
}
func (e Employee) Description() string { ... }
type Manager struct {
    Employee         // embedded
    Reports []Employee
}
m := Manager{Employee: Employee{Name: "Bob", ID: "12345"}, ...}
fmt.Println(m.ID)            // promote: 12345
fmt.Println(m.Description()) // promote
```

**Ad toqquşması:** `Outer`-in öz `X`-i varsa `o.X` xaricinə, `o.Inner.X` daxilinə işarə
edir.

**Embedding ≠ inheritance (3 sübut):**
1. `var e Employee = m` — compile xətası (Manager ≠ Employee; yalnız `m.Employee` açıq)
2. **Dynamic dispatch YOXDUR concrete tiplər üçün**: `Outer.Double()` `Inner.Double`-i
   çağırır, o da `Inner.IntPrinter`-i çağırır — `Outer.IntPrinter` YOX (nəticə:
   "Inner: 20"). Embed edilən metod öz embed olunduğundan xəbərsizdir.
3. Amma embed edilən metodlar method set-ə daxildir → interface implementasiyasını
   təmin edirlər.

### 7. Interface-lər — Type-Safe Duck Typing
**Nədir:** Go-nun yeganə abstract tipi; metod dəsti tələbi.

**Deklarasiya:**
```go
type Stringer interface {
    String() string
}
```
Adlandırma: "-er" konvensiyası (`io.Reader`, `http.Handler`, `json.Marshaler`).

**İmplicit implementasiya:** concrete tip **bəyan etmir** — method set uyğun gəlirsə,
interface dəyişəninə təyin oluna bilər. Nəticə: client kod interface-i sahiblənir ("interfaces
specify what callers need"), implementator xəbərsizdir. Bu, Java-nın explicit
interface-inin client↔provider bağını qırır — yeni provider interface-i heç bilmədən
əvəz edə bilər.

**Standart interface istifadə et:** `io.Reader`-a yazılan kod fayl, yaddaş, gzip-la
işləyir — **decorator pattern** avtomatik (gzip.NewReader io.Reader qaytarır, eyni
`process(r)` kodu compressed faylı oxuyur).

### 8. "Accept Interfaces, Return Structs"
**Nədir:** Parametrlər interface, qaytarma concrete.

**Səbəblər:**
- Interface qaytarmaq → client-i həmin modülə + onun asılılıqlarına birbaşa bağlayır
  (decoupling itir).
- Interface-ə metod əlavəsi → bütün implementasiyalar sınır (versioning problemi);
  struct-a sahə/metod əlavəsi → backward-compatible.
- Hər concrete tip üçün ayrı factory funksiyası yaz (parametrə görə fərqli tiplər qaytaran
  tək factory-dən qaç).
- **İstisnalar:** `error` (müxtəlif implementasiyalar mümkündür) və ayrı-ayrı token
  tipləri qaytaran parser-lər.
- Performans xərci: interface parametri hər çağırışda heap allocation yaradır — yalnız
  profil göstərəndə concrete-ə qayıt.

### 9. Interface və nil (Tələ!)
**Nədir:** Interface nil-dirsə → həm tip, həm dəyər pointer-i nil olmalıdır.

**Kitabdan kod nümunəsi:**
```go
var s *string
var i interface{}
fmt.Println(i == nil)  // true
i = s                  // tip *string, dəyər nil
fmt.Println(i == nil)  // FALSE!
```

Runtime-da interface = 2 pointer (tip + dəyər); tip non-nil olduqda interface non-nil-dir.
Nil interface-də metod çağırmaq panic; non-nil interface-də nil dəyərlə metod çağırmaq —
metod nil-i idarə edirsə işləyir. Dəyərin nil-olub olmadığını reflection olmadan müəyyən
etmək olmur (Ch14).

### 10. Boş interface (`interface{}`) və Type Assertion/Switch
**Boş interface:** 0 metod tələbi → bütün tiplər uyğun. JSON kimi qeyri-müəyyən sxemli
data və generics-ə qədər data strukturları üçün. Amma **mümkün qədər QAÇIN** — güclü
tip dili ilə mübarizə unidiomatikdir.

**Type assertion:** `i.(MyInt)` — yanlışdırsa **panic**; comma-ok ilə təhlükəsiz:
```go
i2, ok := i.(int)
if !ok {
    return fmt.Errorf("unexpected type for %v", i)
}
```
**Qayda:** 100% əminsən belə comma-ok işlət — gələcəkdə reuse panic yaradar.
**Conversion vs assertion:** conversion compile-daə yoxlanır və dəyişdirir; assertion
runtime-da yoxlanır və aşkar edir ("Conversions change, assertions reveal").

**Type switch:**
```go
switch j := i.(type) {
case nil:          // j: interface{}
case int:          // j: int
case MyInt:        // j: MyInt (underlying int olsa da int case-i AYRI tutulur!)
case io.Reader:    // j: io.Reader — interface də case ola bilər
case bool, rune:   // çoxlu tip → j: interface{}
default:           // unknown → mütləq default yaz (yeni tipləri tut)
}
```
Shadowing burada idiomatikdir: `i := i.(type)` — yeni dəyişən eyni adla.

**Əhəmiyyətli istifadələr:**
- **Optional interface-lər:** `io.Copy`-nun `src.(WriterTo)` yoxlaması (sürətli yol);
  `database/sql`-in `si.(StmtExecContext)` yoxlaması (Go 1.8 context məcburiyyətsiz
  uyğunlaşdırma). Çatışmazazlıq: decorator wrap edəndə optional interface gizlənir
  (`bufio.NewReader` `ReaderFrom` optimizasiyasını məhv edir); wrap olunmuş error-lar
  üçün `errors.Is`/`errors.As` lazımdır.

### 11. Funksiya Tipləri — Interface-ə Körpü
**Nədir:** Funksiya tipinə metod əlavə etmək → funksiyalar interface implement edə bilər.

**Kitabdan kod nümunəsi (http paketindən):**
```go
type HandlerFunc func(http.ResponseWriter, *http.Request)
func (f HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    f(w, r)
}
http.HandleFunc("/hello", c.SayHello)  // hər hansı uyğun imzalı funksiya → Handler
```

**Seçim qaydası:** tək funksiya çox funksiya/state asılıdırsa → interface + funksiya
tipi körpüsü; sadə müqayisə funksiyasıdırsa (sort.Slice less) → funksiya parametri.

### 12. Dependency Injection — Frameworksüz
**Nədir:** Kodun ehtiyac duyduğu funksionallığı açıq şəkildə göstərməsi (Robert Martin,
1996, Dependency Inversion Principle).

**Kitabdan kod nümunəsi (tam axın):**
```go
// Implementatorlar (interface-lərdən xəbərsiz):
func LogOutput(message string) { fmt.Println(message) }
type LoggerAdapter func(message string)
func (lg LoggerAdapter) Log(message string) { lg(message) }

type SimpleDataStore struct { userData map[string]string }
func (sds SimpleDataStore) UserNameForID(userID string) (string, bool) { ... }

// Client tərəfin interface-ləri (nə lazımdırsa o qədər):
type DataStore interface { UserNameForID(userID string) (string, bool) }
type Logger interface { Log(message string) }

// Business logic — interface asılılıqları ilə:
type SimpleLogic struct { l Logger; ds DataStore }
func NewSimpleLogic(l Logger, ds DataStore) SimpleLogic { ... }

// Controller — öz ehtiyac interface-i:
type Logic interface { SayHello(userID string) (string, error) }
type Controller struct { l Logger; logic Logic }
func NewController(l Logger, logic Logic) Controller { ... }

// main — yeganə concrete bilən yer (kompozisiya kökü):
func main() {
    l := LoggerAdapter(LogOutput)
    ds := NewSimpleDataStore()
    logic := NewSimpleLogic(l, ds)
    c := NewController(l, logic)
    http.HandleFunc("/hello", c.SayHello)
    http.ListenAndServe(":8080", nil)
}
```

**Prinsiplər:** interface-lər client tərəfdə təyin olunur; `SayGoodbye` Controller-in
`Logic`-inə düşmür (ehtiyac yoxdur); sahələr unexported (l, ds — paket daxili);
test-də Logger capture edən fake injekt olunur. Wire (Google) — codegen əsaslı DI
helper-i istəyənlər üçün.

## Əsas terminlər
- Receiver (qəbul edici) — metoda bağlılıq təyin edən `(p Person)` hissəsi
- Method set (metod dəsti) — tipin (və ya pointer-in) metodlarının toplusu
- Embedded field (daxil edilmiş sahə) — adsız struct sahəsi, promote ilə
- Implicit interface (implict interfeys) — bəyan edilmədən method set uyğunluğu
- Duck typing (ördək yazısı) — metod mövcudiyyətinə görə tip uyğunluğu
- Type assertion (tip iddiası) — `i.(T)` runtime tip yoxlaması
- Optional interface (seçimə bağlı interfeys) — type assertion ilə aşkarlanan əlavə imkan
- Dependency Injection (asılılığın inyeksiyası) — asılılıqların xarici təminatı

## Praktik nəticə

Go OOP dili deyil — "practical" dildir. Dizayn axını: (1) datanı struct-da saxla, məntiqi
metodla; (2) dəyişən hər şey pointer receiver; (3) inheritance axtarma — embedding
kompozisiyadır, dynamic dispatch yoxdur; (4) interface-ləri istehlakçı tərəfdə kiçik
təyin et (böyük əvvəlcədən interface YOX); (5) accept interfaces, return structs; (6)
`i = typedNil` tələsindən qorx — interface nil yoxlaması tip pointer-inə baxır; (7)
assertion-larda həmişə comma-ok; (8) DI üçün framework lazım deyil — factory funksiyaları +
main-də kompozisiya.

## Mənbə
Pages: 189-230
