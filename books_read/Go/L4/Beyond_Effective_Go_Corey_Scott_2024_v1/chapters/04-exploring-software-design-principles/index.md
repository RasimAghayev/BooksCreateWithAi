# Chapter 4 — Exploring Software Design Principles (səh. 18-68)

## Bu chapter nədən bəhs edir?

Go dilinin ideologiyasına uyğun software design prinsipləri (Unix Philosophy,
DRY vs KISS, Delegation, Composition over inheritance, Accept interfaces /
return structs, Singles Principle, ISP, DIP) və 4 klassik OO design pattern-in
(Singleton, Factory Method, Observer, Adapter) Go-da tətbiqi. Müəllifin əsas
motto-su: **"Make it work, make it clean, then (maybe) make it fast"** — əvvəl
işlə, sonra təmizlə, optimallaşdırma yalnız lazım olanda.

## Əsas fikirlər

### 1. Minimalist və kompozisiya (Unix Philosophy)
**Nədir:** Hər proqram/paket/struct/funksiya yalnız **bir** işi görsün və asan
birləşdirilə (composable) olsun.

**Necə işləyir:** `encoding/json` paketi bunun nümunəsidir — `Marshaler` və
`Unmarshaler` **iki ayrı minimal interfeys**dir (tək metodlu), bir `JSONObject`
"boş" interfeysi yox. `crypto/rand` paketi tək məqsəd daşıyır: kriptoqrafik
güclü random ədədlər (3 publik funksiya + 1 dəyişən).

**Nəyə lazımdır:** Paketin məqsədi aydın olsun, yenidən istifadə asanlaşsın,
`io.Reader` kimi ubiquit interfeyslər sayəsində kod type-cast-sız birləşsin.

**Kitabdan kod nümunəsi:**
```go
// Minimal interfeys (encoding/json-dən):
type Marshaler interface {
    MarshalJSON() ([]byte, error)
}
```
**Sub-kod izahı:**
- Tək metod → istifadəçi üçün minimal tələb, asan mock, asan compose.

### 2. DRY vs KISS (praqmatik balans)
**Nədir:** DRY (Don't Repeat Yourself / təkrar etmə) dublikatı azaldır; KISS
(Keep It Simple / sadə saxla) mürəkkəbliyi azaldır. İkisi **tez-tez toqquşur**.

**Necə işləyir:** İki feature oxşar logika istifadə edirsə, DRY onu çıxıb
shared paketə qoymağı tələb edir. Amma çıxarılan kod hər iki tərəfin ehtiyacına
uyğunlaşdırılmalı olduqda mürəkkəbləşir — KISS pozulur.

**Nəyə lazımdır:** Müəllifin tövsiyəsi — DRY tətbiq etməzdən əvvəl bunları
soruş: (a) shared logikanın **sahibi kim olacaq**? (utils/commons paketinə atmaq
QADAĞANDIR); (b) iki istifadə **bir-birinə nə qədər bağlıdır**? Bağlı deyilsə,
bəzən **qəsdən dublikat** (intentional duplication) daha yaxşıdır.

**Vacib:** Məqsəd sətir sayını azaltmaq deyil — **maintainability** (baxım
asanlığı) artırmaqdır.

### 3. Delegation (həvaləetmə)
**Nədir:** Sorğunun emalını daha uyğun obyektə həvalə etmək. 2 məqsəd: dublikat
azaltmaq + kohiziya (cohesion / iç bağlılıq) artırmaq.

**Necə işləyir:** İki `User`-i müqayisə etmək üçün sahə-sahə `if` yazmaq əvəzinə,
müqayisə məntiqi `User`-in öz metoduna keçir: `userA.Equals(userB)`. `IsEqual(a, b)`
kimi ayrı funksiya da işləyir, amma **discoverability** (tapıla bilmə) və kohiziya
metodda daha yaxşıdır.

**Optional delegation (könüllü həvalə):** Type-assertion ilə xüsusi formatı
sınamaq, olmasa default-a düşmək:
```go
func Send(conn net.Conn, data interface{}) error {
    var payload []byte
    if encoder, ok := data.(ByteEncoder); ok { // optional delegation
        payload = encoder.Encode()
    } else {
        payload = []byte(fmt.Sprintf("%s", data)) // fallback
    }
    _, err := conn.Write(payload)
    return err
}
```
**Sub-kod izahı:**
- `data.(ByteEncoder)` → type assertion: data bu interfeysi realləşdirirsə xüsusi yol
- `ok` → realləşdirmirsə default `fmt.Sprintf` fallback
- Standart kitabxanada eyni yanaşma: `encoding/json` (custom marshal), `sql/driver`

### 4. Composition over inheritance (kompozisiya irsə üstünlük)
**Nədir:** Go-da inheritance yoxdur; onun iki məqsədini (kod dublikatını azaltmaq +
eyni şəkildə işlənən obyekt ailəsi yaratmaq) kompozisiya ilə əvəz edirik.

**Necə işləyir:** Shared kodu başqa struct-a köçürüb compose etmək; **anonymous
embedding** (adsız daxiletmə) inheritance görünüşü verir:
```go
type bird struct{}
func (b bird) Fly() { /* ... */ }

type Duck struct {
    bird // anonymous composition + implicit delegation: Duck.Fly() avtomatik işləyir
}
```
**Üstünlükləri:**
- Çoxlu "valideyn" mümkündür (single inheritance-də yalnız 1)
- Daha təbii — komponentlərarası əlaqə **açıq** (inheritance-da magic/indirection var)
- Dəyişməyə daha az davamlı (shotgun surgery riski azalır)
- Runtime-da dinamik — kompozisiya işləmə vaxtında qurula bilər

### 5. Accept interfaces, return structs (interfeys qəbul et, struct qaytar)
**Nədir:** Funksiya **girişdə** interfeys, **çıxışda** konkret struct gözləməli.

**Necə işləyir:** `SendEmail(user *User, ...)` əvəzinə `SendEmail(recipient Recipient, ...)`
— istifadəçi struct qurmaq məcburiyyətində qalmır, istənilən uyğun obyekt verir.

**Return structs səbəbi:**
- Interface qaytarmaq paketin işini artırır (interfeysi daim sync etmək)
- İstifadəçi **öz** interfeysini yazsa, yalnız öz ehtiyacına uyğun minimal
  abstraksiya qurur → coupling azalır, "az təmin etmək = daha çox azadlıq"

### 6. The Singles Principle (təklik prinsipi)
**Nədir:** Müəllifin SRP-dən (Single Responsibility Principle) törətdiyi prinsip:
hər komponentdə **tək məqsəd + tək məsuliyyət + tək abstraksiya səviyyəsi**.

**Necə işləyir:**
- **Tək məqsəd** — "comment test": komponenti bir sətirdə təsvir et; `and, or,
  all, various, several` sözləri görünsə → çox məqsəd var → parçala.
- **Tək məsuliyyət** — Martin-in SRP aydınlaşdırması: komponenti dəyişməyə
  təsir edə biləcək yalnız **bir insan/qrup** olmalıdır (məs. HTTP request formatı
  dəyişən şəxs ≠ DB strukturunu dəyişən şəxs → handler parçalanmalıdır).
- **Tək abstraksiya səviyyəsi** — UserValidator / UserDAO / CreateUserHandler
  piramidası "waterfall"a çevrilir: asılılıqlar minimallaşır, hər komponent
  ayrıca test olunur.

### 7. Interface Segregation Principle (ISP)
**Nədir:** "Client istifadə etmədiyi metodlara asılı olmağa məcbur edilməməlidir"
(Martin). Interfeyslər rolla uyğun minimal olmalıdır.

**Necə işləyir:** `CityModel` (Save/Update/LoadByID/LoadAll) yalnız `LoadByID`
istifadə edən handler-ə verilirse, 3 lazımsız metod da mock-lanmalı olur. Həll:
handler öz paketində tək-metodlu `CityLoader` interfeysi təyin edir.

**Nəyə lazımdır:** Qeyri-ixtiyari coupling-in qarşısını almaq; interfeys eyni
zamanda **tələbi sənədləşdirir**.

### 8. Dependency Inversion Principle (DIP)
**Nədir:** "Yüksək səviyyəli modullar aşağı səviyyəli modullara asılı olmasın;
hər ikisi abstraksiyalara asılı olsun."

**Necə işləyir:** Worker → Car əvəzinə Worker → Transport. "Asılılığı invert
etmək" = konkret deyil, permissive (icazəverici) olmaq. **Go-da fərqlənmə:**
interfeys implementationın yanında deyil, **istifadə olunduğu paketdə** təyin
olunur (consumer-side interface). Bu DRY-nin strict pozuntusudur, amma loose
coupling qazancı çox böyükdür.

### 9. Design Pattern-lər Go-da

**Singleton (tək nümunə):** `sync.Once` ilə təhlükəsiz yaradılır:
```go
var (
    instance   *Cache
    initConfig sync.Once
)

func GetCache() *Cache {
    initConfig.Do(func() {
        instance = &Cache{items: map[string]interface{}{}, createdAt: time.Now()}
    })
    return instance
}
```
**Sub-kod izahı:**
- `sync.Once` → paralel çağırışlarda belə init bir dəfə icra olunur
- Connection pool, cache, logger kimi resurs-intensive obyektlər üçün

**Factory Method (istehsal metodu):** Switch-lə obyekt yaratmaq 2+ yerdə lazımdırsa
və ya yeni tiplər gələcəksə, yaratma məntiqini tək funksiyaya çıxar:
```go
func NewDocumentFormat(format string) DocumentFormat {
    switch format {
    case "md":
        return Markdown{}
    default:
        return HTML{}
    }
}
```

**Observer (müşahidəçi):** Subject (Celebrity) observer-lərə (SuperFan) **channel**
göndərməklə abunə olur; bildiriş `select` + `default` ilə non-blocking:
```go
func (o *Celebrity) Upload(post Post) {
    o.mutex.RLock()
    defer o.mutex.RUnlock()
    for _, fan := range o.fans {
        select {
        case fan <- post: // bildiriş düşdü
        default: // observer hazır deyil — skip, bloklamır
        }
    }
}
```
**Sub-kod izahı:**
- Observer öz channel-ini verir → buffering nəzarəti observer-dədir
- `select/default` → yavaş observer Subject-i bloklamır

**Adapter (uyğunlaşdırıcı):** Köhnə/third-party interfeysi yeni formata uyğunlaşdırır.
Go-nun implicit interface-ləri sayəsində version1 ↔ version2 eyni-metodlu
interfeysləri kiçik adapter struct ilə körpülənir.

## Əsas terminlər

- Composability (birləşdirilə bilənlik)
- Cohesion (iç bağlılıq)
- Loose Coupling (boş əlaqəlilik)
- Implicit Interface (örtülü interfeys realləşdirmə)
- Anonymous Embedding (adsız daxiletmə)
- Type Assertion (tip yoxlaması ilə çevirmə)
- Shotgun Surgery (bir dəyişikliyin çox yerə yayılması)

## Praktik nəticə

- İlk implementasiyada prinsipləri gözləmə — işləyəndən sonra refactor et.
- Interfeysləri **istehlakçı tərəfdə**, tək-metodlu təyin et.
- Go-da inheritance ehtiyacı duyursan → kompozisiya + implicit delegation.
- Singleton → `sync.Once`; Observer → channel + select/default.

## Mənbə

Pages: 18-68 (Chapter 4, Beyond Effective Go Part 2)
