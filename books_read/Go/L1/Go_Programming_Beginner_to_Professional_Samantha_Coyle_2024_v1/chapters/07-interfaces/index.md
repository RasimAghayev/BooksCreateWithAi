# Chapter 7 — Interfaces (səh. 272-309)

## Bu fəsil nədən bəhs edir?

İnterfeyslər: tərifi (method set), implisit implementasiya (Java-nın
`implements`-indən fərqi), duck typing, polimorfizm (Speaker/Shape
nümunələri), "Accept interfaces, return structs" qəbulü, empty
interface (hər tip uyğun), type assertion (v, ok) və type switch, həmçinin
Go 1.18 `any` alias.

## Əsas fikirlər

### 1. İnterfeys nədir?
**Nədir:** metod dəstsinin (method set) təsviri — davranışın blueprintfiliz
imzası; İCRA DETALLARI YOXDUR (tətbiqi implement edən tipin öhdəliyində).

```go
type Speaker interface {
    Speak() string
}

// Standart kitabxanadan:
type Reader interface {
    Read(p []byte) (n int, err error)
}
```

**Motivasiya nümunəsi (loadEmployee):** string-dən, fayldan, HTTP-dən
oxuyan ÜÇ funksiya əvəzinə — bir dənə:
```go
func loadEmployee(r io.Reader)       // hamısı Reader implement edir!
```

- Adlandırma: `-er` suffiksi (Speaker, Reader, Stringer); bir metodlu
  interfeys = metod adı + er

### 2. İmlisit implementasiya
```go
// Java (EKSPLİSİT):
class Dog implements Pet { ... }

// Go (İMLİSİT):
type cat struct{}

func (c cat) Speak() string {
    return "Purr Meow"
}
// cat Speaker-i AVTOMATİK implement edir — implements açar sözü YOXDUR
```

**Üstünlüklər:**
1. Metod dəsti dəyişəndə bütün tiplərdə elan düzəltmək lazım deyil
2. Digər paketlərin tipləri SİZİN interfeysi implement edə bilər
   (decoupling) — məs. `Stringer` fmt paketində, cat onu öz paketində

**Stringer nümunəsi (çoxlu interfeys bir tipdə):**
```go
func (c cat) String() string {         // Stringer implement olundu
    return fmt.Sprintf("%v (%v years old)", c.name, c.age)
}
fmt.Println(c)   // fmt String()-i ÇAĞIRIR → "Oreo (9 years old)"
// cat: həm Speaker, həm Stringer
```

### 3. Duck typing
"Ördək kimi görünür, ördək kimi üzür, ördək kimi vurur → ördəkdir":
tipin interfeysə uyğun gəlməsi İSİ (inheritance) yox, METODLARLA
müəyyən olunur.

```go
func chatter(s Speaker) {
    fmt.Println(s.Speak())
}

c := cat{}
chatter(c)     // cat Speaker-in metod dəstinə uyğundur → keçir
```

### 4. Polimorfizm
**Nədir:** müxtəlif formalarda görünmə qabiliyyəti (Go-da subclassing
YOXDUR → interfeyslərlə polimorfizm).

**Əvvəl (pis — təkrarlanan funksiyalar):**
```go
catSpeak(c); dogSpeak(d); personSpeak(p)   // hər tip üçün ayrıca!
```

**Sonra (yaxşı — interfeys parametri):**
```go
type Speaker interface { Speak() string }

type cat struct{}
func (c cat) Speak() string   { return "Purr Meow" }
type dog struct{}
func (d dog) Speak() string   { return "Woof Woof" }
type person struct{ name string }
func (p person) Speak() string { return "Hi my name is " + p.name + "." }

// VARIADİK interfeys parametri:
func saySomething(say ...Speaker) {
    for _, s := range say {
        fmt.Println(s.Speak())       // hər tipin ÖZ implementasiyası
    }
}

c := cat{}; d := dog{}; p := person{name: "Heather"}
saySomething(c, d, p)
```

**Shape nümunəsi (Exercise 7.02):**
```go
type Shape interface {
    Area() float64
    Name() string
}

type triangle struct{ base, height float64 }
func (t triangle) Area() float64 { return (t.base * t.height) / 2 }
func (t triangle) Name() string  { return "triangle" }

type rectangle struct{ length, width float64 }
func (r rectangle) Area() float64 { return r.length * r.width }
func (r rectangle) Name() string  { return "rectangle" }

type square struct{ side float64 }
func (s square) Area() float64 { return s.side * s.side }
func (s square) Name() string  { return "square" }

func printShapeDetails(shapes ...Shape) {
    for _, item := range shapes {
        fmt.Printf("The area of %s is: %.2f\n", item.Name(), item.Area())
    }
}

t := triangle{base: 15.5, height: 20.1}
r := rectangle{length: 20, width: 10}
s := square{side: 10}
printShapeDetails(t, r, s)
// The area of triangle is: 155.78 ...
```

### 5. "Accept interfaces, return structs"
**Go proverb + Postel's Law:** "qəbul etdiyində liberal, etdiyində
konservativ ol".

```go
// PİS — yalnız string:
func loadPerson2(s string) (Person, error) {
    var p Person
    err := json.NewDecoder(strings.NewReader(s)).Decode(&p)
    return p, err
}

// YAXŞI — io.Reader qəbul:
func loadPerson(r io.Reader) (Person, error) {
    var p Person
    err := json.NewDecoder(r).Decode(&p)
    return p, err
}

// İndi string, fayl, HTTP body — hamısı işləyir:
loadPerson(strings.NewReader(s))    // string
f, _ := os.Open("data.json")
loadPerson(f)                        // *os.File
loadPerson(r.Body)                   // *http.Request.Body
```

**Qaytarmada interfeys YOX** — istifadəçi metod dəstini gəzməli
olardı; struct qaytar → lazımdırsa istifadəçi ÖZÜ interfeys qurur.

### 6. Empty interface (interface{})
```go
type interface{}    // metodsuz — HƏR TİP uyğun gəlir

func emptyDetails(s interface{}) {
    fmt.Printf("(%v, %T)\n", s, s)   // dəyər + konkret tip
}

emptyDetails(cat{name: "oreo"})   // ({oreo}, main.cat)
emptyDetails(99)                   // (99, int)
emptyDetails(false)                // (false, bool)
```

### 7. Type assertion
```go
// Sintaksis: v := s.(T)
var str interface{} = "some string"
v := str.(string)                  // uğurlu — v artıq string
strings.Title(v)                  // "Some String"

// UĞURSUZ → PANIC:
var i interface{} = 49
v := i.(string)                   // PANİK!

// Təhlükəsiz forma — comma, ok:
v, isValid := str.(int)          // uğursuzsa: v=0, isValid=false (panic YOX)
```
- Type CONVERSION (strconv.Atoi(i)) interfeys dəyərinə işləməz —
  ASSERTION lazımdır

### 8. Type switch
```go
func typeExample(i []interface{}) {
    for _, x := range i {
        switch v := x.(type) {
        case int:
            fmt.Printf("%v is int\n", v)
        case string:
            fmt.Printf("%v is a string\n", v)
        case bool:
            fmt.Printf("a bool %v\n", v)
        default:
            fmt.Printf("Unknown type %T\n", v)
        }
    }
}

i := []interface{}{42, "The book club", true, c}
typeExample(i)
// 42 is int / The book club is string / a bool true / Unknown type main.cat
```

**map[string]interface{} analizi (Exercise 7.03):**
```go
type record struct {
    key       string
    valueType string
    data      interface{}
}

func newRecord(key string, i interface{}) record {
    r := record{key: key}
    switch v := i.(type) {
    case int:
        r.valueType = "int"; r.data = v
    case bool:
        r.valueType = "bool"; r.data = v
    case string:
        r.valueType = "string"; r.data = v
    case person:
        r.valueType = "person"; r.data = v
    default:
        r.valueType = "unknown"; r.data = v
    }
    return r
}

m := make(map[string]interface{})
m["person"] = p; m["animal"] = a; m["age"] = 54
m["isMarried"] = true; m["lastName"] = "Smith"

rs := []record{}
for k, v := range m {
    rs = append(rs, newRecord(k, v))
}
```
- İstifadə: naməlum schema-lı JSON (map[string]interface{}), tipə
  görə transformasiya

### 9. any (Go 1.18+)
```go
// any = interface{} alias — tam mübadiləlidir:
func f(x any) { ... }   // func f(x interface{}) { ... } ilə eyni
```

## Activity icmalı
- **7.01:** Payer interfeysi + Developer (hourly×hours + review rating)
  + Manager (salary + commission) + payDetails(p Payer) — polimorfizm +
  map[string]interface{} rating type switch (Excellent=5...Poor=2 və
  int rating-lər birlikdə!)

## Əsas terminlər
- Interface / method set — davranış blueprintfiliz
- İmlisit implementasiya — implements YOXDUR
- Duck typing — metodlarla uyğunluq
- Polimorfizm — bir interfeys, çox forma
- Stringer — String() metodu (fmt tərəfindən çağırılır)
- Accept interfaces, return structs — API dizayn proverbü
- io.Reader/Writer — ən çox istifadə olunan interfeyslər
- Postel's Law — liberal qəbul, konservativ icra
- Empty interface / any — metodsuz; hər tip
- Type assertion — v := s.(T); v, ok := s.(T)
- Type switch — switch v := i.(type)
- map[string]interface{} — naməlum JSON strukturu

## Praktik nəticə
API funksiyaları interfeys qəbul etsin (io.Reader qızıl standart),
struct qaytarsın. Bir çox konkret tip üçün təkrar funksiyalar yerinə
tək interfeys-parametrli funksiya (variadic `...Speaker` kimi). Öz
interfeyslərini -er adlandır; bir-metodlu olsalar metoddan ad al.
Naməlum tipli data üçün interface{} (yaxud any) + type switch — amma
comma-ok assertindən istifadə edib paniki ötür. Tiplər bir neçə
interfeysi eyni anda implement edir; bəzən xəbərsizcə (Stringer).

## Mənbə
Pages: 272-309 (PDF 272-309)
