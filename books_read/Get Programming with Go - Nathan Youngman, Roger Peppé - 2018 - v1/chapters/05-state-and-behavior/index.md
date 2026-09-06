# Unit 5 — State and Behavior (Lesson 21-25)

## Bu unit nədən bəhs edir?

Struct tipləri (composite literal, kopya semantikası, JSON + struct tag), metodlarla struct-lar (constructor konvensiyası, "class alternative"), kompozisiya və embedding (metod forwarding, name collision, inheritance deyil!), interfeyslər (implicit satisfaction, Stringer, polimorfizm). Capstone: Mars heyvan sığınacağı simulyasiyası.

**PDF səhifələr:** 176-214 (L21: 176-184, L22: 185-191, L23: 192-200, L24: 201-210, L25: 211-214)

## Əsas fikirlər

### 1. Struct elanı (L21)
**Problem:** `func distance(lat1, long1, lat2, long2 float64)` — 4 müstəqil float xəta-yönümlüdür.

**Həll:** bağlı dəyərləri bir vahidə qruplaşdır:
```go
var curiosity struct {
    lat  float64
    long float64
}
curiosity.lat = -4.5895          // dot notation
curiosity.long = 137.4417
fmt.Println(curiosity.lat, curiosity.long)   // -4.5895 137.4417
fmt.Println(curiosity)                        // {-4.5895 137.4417}
```

**Tip kimi yenidən istifadə:**
```go
type location struct {
    lat, long float64              // eyni tipli sahələr bir sətirdə
}
var spirit location
spirit.lat = -14.5684
```

### 2. Composite literal — 2 forma (L21)
```go
// 1. Field-value (DƏYƏRMƏTLİ!):
opportunity := location{lat: -1.9462, long: 354.4734}
insight := location{lat: 4.5, long: 135.9}
// + istənilən sıra; + göstərilməyənlər zero value; + sahə əlavəsi kodu pozmur

// 2. Yalnız dəyərlər (qısa, amma kövrək):
spirit := location{-14.5684, 175.472636}
// - sıra mütləq elanla eyni; - yeni sahə → compile xətası; - lat/long qarışsa səssiz səhv
```

**%v vs %+v:**
```go
fmt.Printf("%v\n", curiosity)    // {-4.5895 137.4417}
fmt.Printf("%+v\n", curiosity)   // {lat:-4.5895 long:137.4417} — sahə adları!
```

### 3. Struct-lar KOPYALANIR (L21)
```go
bradbury := location{-4.5895, 137.4417}
curiosity := bradbury               // TAM KOPYA
curiosity.long += 0.0106             // Yellowknife Bay-ə şərqə
fmt.Println(bradbury, curiosity)     // bradbury dəyişməz!
```
Array kimi — funksiyaya ötürülsə kopya gedir (dəyişiklik çağırıcıya görünmür).

### 4. Slice of structs (L21)
**Antipattern:** 2 paralel slice (`lats`, `longs`) — misalign təhlükəsi!
```go
// DÜZGÜN:
type location struct {
    name string
    lat  float64
    long float64
}
locations := []location{
    {name: "Bradbury Landing", lat: -4.5895, long: 137.4417},
    {name: "Columbia Memorial Station", lat: -14.5684, long: 175.472636},
    {name: "Challenger Memorial Station", lat: -1.9462, long: 354.4734},
}
```

### 5. JSON + struct tag (L21)
```go
type location struct {
    Lat  float64 `json:"latitude"`
    Long float64 `json:"longitude"`
}
curiosity := location{-4.5895, 137.4417}
bytes, err := json.Marshal(curiosity)
fmt.Println(string(bytes))
// {"latitude":-4.5895,"longitude":137.4417}
```
- **Sahələr EXPORTED olmalı** (böyük hərf) — yoxsa output `{}` olur
- Struct tag: raw string (backtick) daxilində `key:"value"`; çoxlu tag: `` `json:"lat" xml:"lat"` ``
- JSON keys = sahə adları default; tag ilə override (snake_case üçün)

### 6. Struct + metodlar (L22)
**Kitabdan kod nümunəsi — coordinate (DMS → decimal):**
```go
// coordinate in degrees, minutes, seconds in a N/S/E/W hemisphere.
type coordinate struct {
    d, m, s float64
    h       rune
}

// decimal converts a d/m/s coordinate to decimal degrees.
func (c coordinate) decimal() float64 {
    sign := 1.0
    switch c.h {
    case 'S', 'W', 's', 'w':
        sign = -1
    }
    return sign * (c.d + c.m/60 + c.s/3600)
}

lat := coordinate{4, 35, 22.2, 'S'}
long := coordinate{137, 26, 30.12, 'E'}
fmt.Println(lat.decimal(), long.decimal())   // -4.5895 137.4417
```

### 7. Constructor funksiyası konvensiyası (L22)
**Go-da konstruktor DİL XÜSUSİYYƏTİ DEYİL — konvensiyadır:**
```go
// newLocation from latitude, longitude d/m/s coordinates.
func newLocation(lat, long coordinate) location {
    return location{lat.decimal(), long.decimal()}
}

curiosity := newLocation(
    coordinate{4, 35, 22.2, 'S'},
    coordinate{137, 26, 30.12, 'E'},
)
```
- Ad: `newType` (unexported) / `NewType` (exported)
- Fərqli girişlər üçün fərqli konstruktorlar: `newLocationDMS`, `newLocationDD`
- `errors.New` kimi tək `New` — paket prefiksi ilə `errors.NewError` redundantdır

### 8. Class alternativi (L22)
**Kitabdan kod nümunəsi — world + distance:**
```go
type world struct {
    radius float64
}

var mars = world{radius: 3389.5}

// rad converts degrees to radians.
func rad(deg float64) float64 {
    return deg * math.Pi / 180
}

// distance calculation using the Spherical Law of Cosines.
func (w world) distance(p1, p2 location) float64 {
    s1, c1 := math.Sincos(rad(p1.lat))
    s2, c2 := math.Sincos(rad(p2.lat))
    clong := math.Cos(rad(p1.long - p2.long))
    return w.radius * math.Acos(s1*s2+c1*c2*clong)
}

spirit := location{-14.5684, 175.472636}
opportunity := location{-1.9462, 354.4734}
dist := mars.distance(spirit, opportunity)
fmt.Printf("%.2f km\n", dist)     // 9669.71 km
```
- Radius parametr kimi YOX — world tipinin SAHƏSİ; hər planet üçün eyni metod
- **Squint testi:** struct + metodlar ≈ class (amma miras YOX)

### 9. Kompozisiya (L23)
**Antipattern — bütün sahələr bir structda:**
```go
type report struct {
    sol       int
    high, low float64
    lat, long float64
}
```

**Kompozisiya — kiçik tiplərdən yığ:**
```go
type report struct {
    sol         int
    temperature temperature
    location    location
}
type temperature struct {
    high, low celsius
}
type location struct {
    lat, long float64
}
type celsius float64

bradbury := location{-4.5895, 137.4417}
t := temperature{high: -1.0, low: -78.0}
report := report{sol: 15, temperature: t, location: bradbury}
fmt.Printf("a balmy %vº C\n", report.temperature.high)   // -1
```
+ hər tip müstəqil işləyə bilər + metod asmaq:
```go
func (t temperature) average() celsius {
    return (t.high + t.low) / 2
}
report.temperature.average()     // -39.5
```

**Manual forwarding (boilerplate):**
```go
func (r report) average() celsius {
    return r.temperature.average()
}
report.average()     // dərhal istifadə — amma əl ilə yazıldı
```

### 10. Embedding — avtomatik forwarding (L23)
```go
type report struct {
    sol         int
    temperature               // SAHƏ ADI YOXDUR — embedding!
    location
}

report := report{
    sol:         15,
    location:    location{-4.5895, 137.4417},
    temperature: temperature{high: -1.0, low: -78.0},
}
fmt.Println(report.average())          // forwarding İŞLƏYİR!
fmt.Println(report.temperature.average())   // sahə hələ də mövcuddur
fmt.Println(report.high)                // sahələr DƏ forward olunur!
report.high = 32                        // = report.temperature.high dəyişir
```

**Hər hansı tip embed oluna bilər:**
```go
type sol int
type report struct {
    sol
    location
    temperature
}
func (s sol) days(s2 sol) int { ... }
report.days(1446)        // metod forward
report.sol.days(1446)    // birbaşa
```

### 11. Name collision (L23)
sol və location hər ikisi `days` metoduna sahibdirsə:
- **İstifadə olunmursa** → hər şey OK (kompilyator ağıllıdır)
- `report.days(1446)` çağırılırsa → **"ambiguous selector" xətası**
- **Həll:** report-da öz `days` metodu yaz — o, priority alır:
```go
func (r report) days(s2 sol) int {
    return r.sol.days(s2)          // hansını istəyirsən ona forward et
}
```

### 12. Bu inheritance DEYİL! (L23)
- **Inheritance:** "rover IS-A vehicle" — funksionallığı miras alır
- **Composition:** "rover HAS-A engine, wheels" — hissələrdən yığılır
- Go: **"Favor object composition over class inheritance"** — Gang of Four, 1994
- **Texniki fərq:** `average()` metodu forward edilsə belə **receiver həmişə temperature tipidir** (report YOX) — delegation/inheritance-də receiver dəyişərdi
- "Use of classical inheritance is always optional" — Sandi Metz

### 13. Interfeys nədir? (L24)
**Fəlsəfə:** "what something can DO, not what it IS" — konkret deyil, abstrakt (qələm yox, "yazı aləti").

```go
var t interface {
    talk() string
}

type martian struct{}
func (m martian) talk() string { return "nack nack" }

type laser int
func (l laser) talk() string { return strings.Repeat("pew ", int(l)) }

t = martian{}           // hər ikisi talk() təmin edir
fmt.Println(t.talk())    // nack nack
t = laser(3)
fmt.Println(t.talk())    // pew pew pew
```
- **Polimorfizm** = "many shapes" — t forma dəyişir
- Java-dan fərq: **implements elanı YOXDUR** — implicit

**Adlandırma konvensiyası:** tək metodlu interfeys `-er` şəkilçisi: `talker`.
```go
type talker interface {
    talk() string
}

func shout(t talker) {
    louder := strings.ToUpper(t.talk())
    fmt.Println(louder)
}
shout(martian{})      // NACK NACK
shout(laser(2))       // PEW PEW
shout(crater{})        // XƏTA: crater does not implement talker
```

### 14. Embedding interfeysi təmin edir (L24)
```go
type starship struct {
    laser                // laser-i embed et
}
s := starship{laser(3)}
fmt.Println(s.talk())    // pew pew pew — forwarding!
shout(s)                 // starship talker SATISFY edir ✓
```

### 15. Interfeysləri sonra kəşf et (L24)
**Implicit güc:** mövcud koda (yazdığın YOXSA yazmadığın) interfeys tətbiq et.

**Stardate nümunəsi:**
```go
// Əvvəl: konkret time.Time
func stardate(t time.Time) float64 {
    doy := float64(t.YearDay())
    h := float64(t.Hour()) / 24.0
    return 1000 + doy + h
}

// Sonra: interfeysə keç — həm Earth, həm Mars!
type stardater interface {
    YearDay() int
    Hour() int
}
func stardate(t stardater) float64 { ... }

type sol int
func (s sol) YearDay() int { return int(s % 668) }   // Mars ili = 668 sol
func (s sol) Hour() int     { return 0 }

stardate(time.Date(2012, 8, 6, ...))    // 1219.2 — time.Time öz-özünə uyğun!
stardate(sol(1422))                      // 1086.0
```
**Güc:** `java.time` belə interfeysi sonradan "implements" edə bilməzdi — Go-da STANDART KİTABXANA tiplərinə öz interfeyslərini tətbiq edirsən!

### 16. Standart kitabxana interfeysləri (L24)
> "Go encourages composition over inheritance, using simple, often one-method interfaces" — Rob Pike

**fmt.Stringer:**
```go
type Stringer interface {
    String() string
}

// location üçün:
func (l location) String() string {
    return fmt.Sprintf("%v, %v", l.lat, l.long)
}
fmt.Println(curiosity)   // -4.5895, 137.4417 — Println String() çağırır!
```

Digər məşhurlar: `io.Reader`, `io.Writer`, `json.Marshaler`. `io.ReadWriter` — **interfeys embedding-i** (`io.Reader` + `io.Writer` birləşməsi).

### 17. Capstone: Martian Animal Sanctuary (L25)
Tapşırıq: 
- Animal tipləri (name + `String()` — Stringer)
- `move()` və `eat()` metodları (random description)
- 3 sol (72 saat) day/night cycle — gecə hamı yatır, hər saat random heyvan random action
- Struct + interface istifadə et

## Class vs Go müqayisəsi
| Klassik OOP | Go |
|---|---|
| Class | struct + metodlar |
| Constructor (dil xüsusiyyəti) | `NewType` funksiya konvensiyası |
| Inheritance (IS-A) | Embedding (HAS-A) + forwarding |
| implements açar sözü | Implicit satisfaction |
| Interface-ə implement elanı | Yalnız metod yarat — kifayət |
| Polimorfizm class hierarchy ilə | Interfeyslərlə (daha sərbəst) |

## Əsas terminlər
- Structure / Field / Dot Notation
- Composite Literal (field-value vs values-only)
- `%+v` (sahə adları ilə çap)
- Struct kopya semantikası
- Slice of Structs
- Struct Tag (`json:"..."`)
- json.Marshal / MarshalIndent
- Constructor Function (NewType konvensiyası)
- Composition (HAS-A)
- Embedding (sahə adsız elan)
- Method Forwarding / Promotion
- Ambiguous Selector (name collision)
- Inheritance (IS-A) vs Composition
- Interface / Method Set
- Implicit Satisfaction
- Polymorphism ("many shapes")
- `-er` Suffix Convention
- fmt.Stringer / io.Reader / io.Writer / json.Marshaler
- Interface Embedding (io.ReadWriter)
- Gang of Four ("favor composition over inheritance")

## Praktik nəticə
- Composite literal-da field-value formasını seç — kod dəyişikliyinə davamlıdır; values-only yalnız stabil kiçik tiplər üçün.
- Paralel slice-lar YOX — slice of structs; data bir vahid kimi hərəkət etsin.
- Konstruktor: `newLocation` kimi adlandır; mürəkkəb yaratma məntiqini bir yerə topla.
- Embedding ilə boilerplate forwarding-i aradan qaldır; amma name collision-da öz metodunu yaz — o, priority alır.
- İnterfeysləri ƏVVƏLCƏDƏN yox, lazım olanda kəşf et — implicit satisfaction mövcud kodla da işləyir.
- `String()` metodu ilə Println çıxışını beautify et — fmt.Stringer-in ən sadə tətbiqi.
- Embed etdiyin tip interfeysi təmin edirsə, xarici tip də avtomatik təmin edir (starship → talker).

## Mənbə
Pages: 176-214 (PDF), book pages 161-199 (Lesson 21-25)
