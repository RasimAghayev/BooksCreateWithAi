# Chapter 4 — Structs, Methods, and Interfaces (səh. 64-82)

## Bu fəsil nədən bəhs edir?

Go-da öz data strukturlarının qurulması: struct-lar (adlı sahələr ardıcıllığı),
konstruktor funksiyaları, anonim/nested/embedded struct-lar, metodlar (value və
pointer receiver), embedded metodlar (irsiyyət əvəzi) və interfeyslər (implicit
implementation, empty interface, type switch).

## Əsas fikirlər

### 1. Struct-lar
**Nədir:** Adlı sahələrdən (fields) ibarət strukturdur; hər sahənin adı və tipi var
(C/C++ struct analogu).

**Kitabdan kod nümunəsi:**
```go
type Rectangle struct {
    Height int
    Width  int
}

func main() {
    a := Rectangle{}                    // {0 0} — zero value-larla
    b := Rectangle{4, 4}                // sıra ilə
    c := Rectangle{Width: 10, Height: 3} // sahə adı ilə
    d := Rectangle{Width: 7}            // {0 7} — qalanlar zero
}
```

**Sub-kod izahı:**
- `Rectangle{}` → bütün sahələr zero value alır
- `Rectangle{4,4}` → dəyərlər elan sırası ilə
- `Rectangle{Width: 7}` → ad ilə; göstərilməyən sahələr zero qalır

### 2. Konstruktor funksiyaları
**Nəyə lazımdır:** Go-da klassik konstruktor anlayışı YOXDUR. Zero value
problemli olduqda (məs. eni 0 olan düzbucaqlı mənasızdır) `New...`
funksiyası yazmaq yaxşı praktikadır — tipik olaraq pointer qaytarır.

**Kitabdan kod nümunəsi:**
```go
func NewRectangle(height int, width int) *Rectangle {
    return &Rectangle{height, width}
}

// Validasiya ilə gücləndirilmiş variant:
func NewRectangle(height int, width int) (*Rectangle, error) {
    if height <= 0 || width <= 0 {
        return nil, errors.New("params must be greater than zero")
    }
    return &Rectangle{height, width}, nil
}

r, err := NewRectangle(2, 0)
if err != nil {
    // ...
}
```

### 3. Anonim struct-lar
**Nədir:** Adı olmayan, birbaşa istifadəyə təyin edilən struct. Eyni sahələri
olan adlı struct-la müqayisə oluna bilər.

```go
ac := struct{ x int; y int; radius int }{1, 2, 3}
c := Circle{10, 10, 3}   // type Circle struct {...}

fmt.Printf("%+v\n", ac)             // {x:1 y:2 radius:3}
fmt.Println(reflect.TypeOf(ac))     // struct { x int; y int; radius int }
ac.x = 3                             // sahələr dəyişilir
ac = c                               // eyni sahələr → müqayisə mümkün
```

### 4. Nested və Embedded struct-lar
**Nested (adlı daxiletmə):**
```go
type Coordinates struct {
    x int
    y int
}
type Circle struct {
    center Coordinates   // adlı sahə
    radius int
}
c := Circle{Coordinates{1, 2}, 3}   // {center:{x:1 y:2} radius:3}
```

**Embedded (adsız daxiletmə) — irsiyyətə bənzər:**
```go
type Circle struct {
    Coordinates      // adsız sahə → embedding
    radius int
}
c := Circle{Coordinates{1, 2}, 3}
fmt.Printf("%+v\n", c.Coordinates)  // {x:1 y:2}
fmt.Println(c.x, c.y)               // x, y birbaşa əlçatandır → 1 2
```

**Dvoqluqluq (ambiguity) qaydası:** hər iki embedded struct eyni adlı sahə
saxlayırsa birbaşa çıxış compile xətası verir:
```go
type A struct { fieldA int }
type B struct { fieldA int }
type C struct { A; B }
// c.fieldA      → XƏTA: ambiguous
c.A.fieldA, c.B.fieldA   // sahibi göstərilməlidir
```

### 5. Metodlar
**Nədir:** Go-da class YOXDUR; metod = receiver-i (qəbul edən tipi) olan xüsusi
funksiyadır.

**Kitabdan kod nümunəsi:**
```go
func (r Rectangle) Surface() int {
    return r.Height * r.Width
}

r := Rectangle{2, 3}
fmt.Printf("rectangle %v has surface %d", r, r.Surface())  // 6
```

**Value vs Pointer receiver — əsas fərq:**
```go
func (r Rectangle) Enlarge(factor int) {   // VALUE → kopya üzərində işləyir
    r.Height = r.Height * factor
    r.Width = r.Width * factor
}   // orijinal DƏYİŞMİR

func (r *Rectangle) EnlargeP(factor int) { // POINTER → orijinalı dəyişir
    r.Height = r.Height * factor
    r.Width = r.Width * factor
}

rect := Rectangle{2, 2}
rect.Enlarge(2)    // {2 2} — dəyişməz
rect.EnlargeP(2)   // {4 4} — dəyişdi
```

**Vacib qeydlər:**
- `rect.EnlargeP(2)` → Go bunu avtomatik `(&rect).EnlargeP(2)` kimi tərcümə edir
- Pointer receiver daha effektivdir (kopya azaldır); amma hər halda ardıcıl
  olun — value/pointer receiver-ları qarışdırmayın

### 6. Embedded metodlar
**Nədir:** Embedded struct-ın metodları xarici struct-a avtomatik əlçatandır
 olur — Go-dakı irsiyyət (inheritance) əvəzi.

```go
type Rectangle struct { Height int; Width int }
func (r Rectangle) Surface() int { return r.Height * r.Width }

type Box struct {
    Rectangle    // embedding
    depth int
}

func (b Box) Volume() int {
    return b.Surface() * b.depth   // Rectangle metodunu birbaşa çağırır
}

b := Box{Rectangle{3, 3}, 3}
fmt.Println("Volume", b.Volume())   // Volume 27
```

Dvoqluqluq embedded metodlarda da keçərlidir — sahibi açıq göstərin:
```go
func (g Greeter) Speak() string {
    // return g.Hi() → ambiguous, xəta
    return g.A.Hi() + g.B.Hi()   // düzgün
}
```

### 7. Interfeyslər
**Nədir:** Metod toplusu; məntiq və ya dəyər saxlamır. **Implicit
implementation** — tip interfeysi yalnız bütün metodları tətbiq edəndə
həyata keçirir (extends/implements açar sözü lazım deyil).

**Kitabdan kod nümunəsi:**
```go
type Animal interface {
    Roar() string
    Run() string
}

type Dog struct {}
func (d Dog) Roar() string { return "woof" }
func (d Dog) Run() string  { return "run like a dog" }

type Cat struct {}
func (c *Cat) Roar() string { return "meow" }
func (c *Cat) Run() string  { return "run like a cat" }

func RoarAndRun(a Animal) {
    fmt.Printf("%s and %s\n", a.Roar(), a.Run())
}

func main() {
    myDog := Dog{}
    myCat := Cat{}
    RoarAndRun(myDog)    // işləyir — Dog value receiver ilə
    RoarAndRun(&myCat)   // işləyir — Cat pointer receiver ilə
}
```

**Receiver/interfeys qaydaları (çox vacib):**
| Çağırış | Nəticə | Səbəb |
|---|---|---|
| `RoarAndRun(myDog)` | işləyir | Dog metodları value receiver |
| `RoarAndRun(&myDog)` | işləmir | *Dog üçün metod dəsti fərqlidir |
| `RoarAndRun(myCat)` | işləmir | Cat metodları pointer receiver |
| `RoarAndRun(&myCat)` | işləyir | *Cat metod dəsti tam uyğundur |

Pointer receiver metod yalnız pointer tipinə aiddir:
```go
type Greeter interface { SayHello() string }
type Person struct{ name string }
func (p *Person) SayHello() string { return "Hi! This is me " + p.name }

var g Greeter
p := Person{"John"}
// g = p      → XƏTA: Person *deyil*, SayHello *Person-ündür
g = &p        // işləyir
```

**Method overloading QADAĞANDIR** — eyni metod adı ilə həm pointer, həm value
receiver yaza bilməzsiniz.

### 8. Empty interface və type switch
**Nədir:** `interface{}` — metodu olmayan interfeys; hər tip tərəfindən
həyata keçirilir. Tip əvvəlcədən bilinməyəndə istifadə olunur.

```go
var aux interface{}
fmt.Println(aux)   // <nil>
aux = 10           // int təyin oluna bilər
aux = "hello"      // string də
```

**Type switch — runtime tip müəyyənləşdirmə:**
```go
aux := []interface{}{42, "hello", true}

for _, i := range aux {
    switch t := i.(type) {
    default:
        fmt.Printf("%T —> %s\n", t, i)
    case int:
        fmt.Printf("%T —> %d\n", t, i)
    case string:
        fmt.Printf("%T —> %s\n", t, i)
    case bool:
        fmt.Printf("%T —> %v\n", t, i)
    }
}
```

**Sub-kod izahı:**
- `i.(type)` → type switch konstruksiyası; `t` cari tipi alır
- `%T` → dəyərin tip adını çap edir
- Nəticə: `int —> 42`, `string —> hello`, `bool —> true`

## Əsas terminlər
- Struct — adlı sahələr strukturu
- Constructor pattern — `New...` funksiyası ilə (dildə yox, konvensiya)
- Anonymous struct (anonim strukturdur) — adsız struct literal
- Embedding (daxiletmə) — adsız sahə; irsiyyət əvəzi
- Receiver (qəbul edən) — metodu bağlıyan tip `(r Rectangle)` / `(r *Rectangle)`
- Interface — metod imzaları toplusu
- Implicit implementation (implisit tətbiq) — implements açar sözü olmadan
- Empty interface — `interface{}`; hər tipi qəbul edir
- Type switch — runtime tip yoxlaması `i.(type)`

## Praktik nəticə
Go OOP dili deyil — class/hierarchy yoxdur; struct + method + interface
kombinasiyası oxşar funksionallıq verir. Konstruktor lazımdırsa `New...`
funksiyası yazın (pointer + error qaytarmaqla). Pointer receiver istifadə
edən tip interfeysə yalnız pointer kimi uyğun gəlir — bunu dizayn zamanı
nəzərə alın. Embedding zamanı ad toqquşmasına diqqət edin: `c.A.fieldA`
formasında sahibi göstərin.

## Mənbə
Pages: 64-82 (PDF 64-82)
