# Get Programming with Go — Cheat Sheet (bütün kitabdan toplanmış)

Nathan Youngman, Roger Peppé — Manning, 2018

## Sintaksis əsasları

### Proqram strukturu
```go
package main
import "fmt"
func main() {                      // { eyni sətirdə — one true brace style
    fmt.Println("Hello, playground")
}
```

### Çap funksiyaları
```go
fmt.Print("a", "b")              // aralarında boşluq, sətir keçməz
fmt.Println("hello")              // + newline
fmt.Printf("My weight on %v is %v lbs.\n", "Mars", 149.0*0.3783)
fmt.Printf("%-15v $%4v\n", "SpaceX", 94)   // align
fmt.Printf("Type %T for %[1]v\n", days)     // %T = tip, [1] = arqument təkrarı
fmt.Printf("%+v\n", curiosity)               // struct sahə adları ilə
fmt.Printf("%#v\n", v)                       // Go representasiyası
```

### Dəyişənlər
```go
const lightSpeed = 299792          // dəyişməz
var distance = 56000000           // dəyişən
distance := 56000000               // short declaration (yalnız funksiya daxili!)
var (                             // qrup
    distance = 56000000
    speed    = 100800
)
var a, b = 1, 2                    // çoxlu
weight *= 0.3783; age++; count--; price /= 2   // qısayollar (prefix ++ YOX)
```

### Random
```go
var num = rand.Intn(10) + 1       // 1-10; Intn(10) = 0-9 (off-by-one!)
var distance = rand.Intn(345000001) + 56000000 // range random
```

### Şərtlər və looplar
```go
if command == "go east" { ... } else if ... { ... } else { ... }
if num := rand.Intn(3); num == 0 { ... }        // scoped declaration

switch command {
case "go east", "go inside":     // vergüllü çoxlu dəyər
    ...
default:
    ...
}
switch {                          // şərtsiz switch (if zənciri)
case room == "cave": ...
}

for count > 0 { ... }            // while
for count := 10; count > 0; count-- { ... }     // klassik
for { ... break ... }            // sonsuz
for i, c := range "¿Cómo?" { ... }              // rune-lara açır
for _, v := range slice { ... }                  // _ = ignore
```

### Məntiq
```go
true / false                      // yalnız bunlar — "" / 0 false DEYİL
== != < > <= >=                   // müqayisə (yalnız eyni tip)
&& || !                           // məntiq (short-circuit)
leap := year%400 == 0 || (year%4 == 0 && year%100 != 0)
```

## Tiplər

### Rəqəmlər
| Tip | İstifadə | Qeyd |
|---|---|---|
| float64 | default kəsr | %f, %.3f precision, %4.2f width, %05.2f zero-pad |
| float32 | yaddaş qənaəti | single precision |
| int / uint | default tam | platforma 32/64 |
| int8-64 / uint8-64 | spesifik | uint8 = CSS rəng; int64 = Unix time |
| big.Int/Float/Rat | nəhəng | `big.NewInt(86400)`, `SetString("24e18", 10)`, `.Div()` |
- Float müqayisə: `math.Abs(a-b) < 0.0001` — `==` YOX
- Vurma bölmədən ƏVVƏL: `(c * 9 / 5) + 32`
- Wrap: 255+1 (uint8) = 0; yoxla: `math.MaxInt16` konstantları
- Exponent: `41.3e12`; hex: `0x00, 0x8d` + `%02x`; bits: `%08b`
- **Untyped constants:** `const d = 24e18` — int64-i aşsa da OK (compile-time big); funksiyaya ötürüləndə tip istəyir

### Mətn
```go
s := "literal"                    // escape işləyir
s := `raw string`                  // \n mətn kimi; çoxsətir; C:\go üçün ideal
rune = int32                       // 1 code point; 'A' = 65
byte = uint8                       // 1 bayt; ASCII
c := 'A'                           // rune literal
s[i]                               // BYTE (rune YOX!)
len(s)                             // bayt sayı
utf8.RuneCountInString(s)          // simvol sayı
for i, c := range s                // rune-ara açır
string(pi)                          // rune → string
strconv.Itoa(10)                   // int → string
fmt.Sprintf("%v", v)                // hər şey → string
countdown, err := strconv.Atoi("10")   // string → int + error
```
- String immutable: `message[5] = 'd'` → compile xətası

### Type conversion
```go
float64(age); int(earthDays)       // truncate (round YOX!)
int16(bh)                          // Arianne-5: range yoxla → wrap!
if bh < math.MinInt16 || bh > math.MaxInt16 { ... }
uint8(red); celsius(warmUp)
string(false) → XƏTA              // bool ↔ numeric/str yalnız manual if
```
- Koercion YOX: `"10" - 1` → mismatched types xətası

## Funksiyalar və metodlar

### Funksiya
```go
func kelvinToCelsius(k float64) float64 {
    k -= 273.15
    return k
}
func Unix(sec, nsec int64) Time                    // eyni tip — bir yazılış
func Atoi(s string) (i int, err error)             // adlı nəticələr
func Println(a ...interface{}) (n int, err error)  // variadic + empty interface
```

### Yeni tip + metod
```go
type celsius float64               // ALIAS DEYİL — müstəqil tip
type kelvin float64
// celsius + fahrenheit qarışmır — compile xətası!

func (k kelvin) celsius() celsius {   // receiver
    return celsius(k - 273.15)
}
c := k.celsius()                    // dot notation
```

### First-class functions
```go
sensor := fakeSensor                // funksiya ÖZÜ
sensor()                            // çağırış
var sensor func() kelvin            // tip elanı
type sensor func() kelvin           // adlı funksiya tipi

func measureTemperature(samples int, s func() kelvin) { ... }
func calibrate(s sensor, offset kelvin) sensor {   // funksiya QAYTAR
    return func() kelvin {
        return s() + offset          // closure — s, offset yaşayır
    }
}
```

## Kolleksiyalar

### Array
```go
var planets [8]string
planets[0] = "Mercury"
dwarfs := [5]string{...}
planets := [...]string{...}          // kompilyator sayır; sondakı vergül məcburi
// KOPYALANIR; [5]string ≠ [8]string (fərqli TİPLƏR!)
```

### Slice
```go
terrestrial := planets[0:4]         // half-open; [0:4) daxil
planets[:4]; planets[4:]; planets[:]
dwarfs := []string{...}              // [] boş = slice tipi
len(s); cap(s)
dwarfs = append(dwarfs, "Orcus", "Salacia")  // variadic
make([]string, 0, 10)                // preallokasiya
planets[0:4:4]                      // 3-indeks: cap məhdud → append təcrid
planets...                           // slice → arqumentlər (variadic çağırış)
func terraform(prefix string, worlds ...string) []string
```

### Map
```go
temperature := map[string]int{"Earth": 15, "Mars": -65}
temp := temperature["Moon"]         // yoxdursa ZERO VALUE
if moon, ok := temperature["Moon"]; ok { }   // comma-ok
delete(temperature, "Earth")
make(map[float64]int, 8)            // preallokasiya
// KOPYALANMIR — paylaşılır!
frequency[t]++                       // sayğac
groups[g] = append(groups[g], t)     // map of slices — qruplama
set := make(map[float64]bool); set[t] = true   // SET improvizasiyası
```

## Struct, metod, embedding

```go
type location struct {
    name string
    lat, long float64
}
loc := location{lat: -4.5, long: 135.9}   // field-value (davamlı)
loc := location{-4.5, 135.9}              // values-only (kövrək)
loc.name                              // dot notation
// KOPYALANIR

type Grid [9][9]int8
g := &Grid{...}                        // pointer literal

// Constructor konvensiyası:
func newLocation(lat, long coordinate) location { ... }

// Embedding — sahə adı YAZILMIR:
type report struct {
    sol         int
    temperature       // metodlar + sahələr promote
    location
}
report.average()              // forwarding avtomatik
report.high                    // = report.temperature.high
// Name collision → "ambiguous selector" → öz metod yaz
```

### JSON
```go
type location struct {
    Lat  float64 `json:"latitude"`    // tag; sahələr EXPORTED olmalı
    Long float64 `json:"longitude"`
}
bytes, err := json.Marshal(curiosity)
json.MarshalIndent(loc, "", "  ")
```

## Interfeyslər

```go
type talker interface {
    talk() string
}
// IMPLICIT — implements sözü YOX
func shout(t talker) { ... }
shout(martian{}); shout(laser(2))

// Embedding interfeysi təmin edir:
type starship struct{ laser }
shout(starship{laser(3)})

// fmt.Stringer:
func (l location) String() string {
    return fmt.Sprintf("%v, %v", l.lat, l.long)
}   // Println avtomatik String() istifadə edir
```

## Pointer və nil

```go
&answer                            // address operator
*address                           // dereference
var home *string                    // pointer tipi
timmy := &person{name: "Timothy"}   // struct literal + &
timmy.superpower                    // avtomatik dereference
func birthday(p *person) { p.age++ }   // mutation
func (p *person) birthday() { ... }    // pointer receiver
levelUp(&player.stats)               // interior pointer
```
- Map = pointer (gizli); slice = ptr+len+cap; pointer-a ehtiyac YOX
- Pointer receiver + interface → yalnız `&pew` satisfy edir
- Nil: dereference → panic; metod guard `if p == nil { return }`
- Nil slice: range/len/append OK; Nil map: oxu OK, yazı PANIC
- Nil interface: tip+dəyər hər ikisi nil olmalı; `(*int)(nil)` tələsi

## Error handling

```go
files, err := ioutil.ReadDir(".")
if err != nil {
    fmt.Println(err)
    os.Exit(1)
}

// defer — təmizləmə:
defer f.Close()

// safeWriter pattern ("errors are values"):
type safeWriter struct {
    w   io.Writer
    err error
}
func (sw *safeWriter) writeln(s string) {
    if sw.err != nil { return }
    _, sw.err = fmt.Fprintln(sw.w, s)
}

// Error dəyişənləri:
var (
    ErrBounds = errors.New("out of bounds")
    ErrDigit  = errors.New("invalid digit")
)
switch err {
case ErrBounds, ErrDigit: ...       // müqayisə ÜNVANLARLA!
}

// Custom error tipi:
type SudokuError []error
func (se SudokuError) Error() string { ... }   // error interfeysi

// Type assertion:
if errs, ok := err.(SudokuError); ok { ... }

// panic / recover:
defer func() {
    if e := recover(); e != nil { ... }
}()
panic("I forgot my towel")
```

## Concurrency

```go
go sleepyGopher()                   // goroutine
c := make(chan int)                  // unbuffered
c <- 99; r := <-c                    // send / receive (hər ikisi bloklanır)
v, ok := <-c                         // + bağlanıb?
close(c)                             // yazma → panic; oxu → zero value
for v := range c { ... }            // bağlanana qədər

// select — çox kanal:
select {
case gopherID := <-c: ...
case <-time.After(2 * time.Second):   // TIMEOUT
    return
}

// Pipeline:
go sourceGopher(c0)                  // close(downstream)
go filterGopher(c0, c1)              // for item := range upstream
printGopher(c1)                      // for v := range upstream

// Mutex:
var mu sync.Mutex
mu.Lock()
defer mu.Unlock()
// sadə, tək state, TƏK mutex

// Worker:
func worker() {
    for {
        select {
        case c := <-r.commandc: ...
        case <-nextMove: ...
        }
    }
}
```

## Performans/keyfiyyət qeydləri (kitabdan)
- Float + pul YOX — sent-ləri int-də saxla
- Arianne 5 (1996): float64→int16 overflow → raket partladı — Go-da wrap olur, yoxla!
- 2038: Unix time int32 aşır → int64 istifadə et
- Iron.io: 30 Ruby server → 2 Go server
- Time sharing: goroutine sırası qeyri-deterministik — həmişə ixtiyari güman et
