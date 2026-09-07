# Chapter 7 — Interfaces (İnterfeyslər)

## Bu fəsil nədən bəhs edir?

İnterfeys anlayışı (metod dəsti, davranış), implicit implementasiya (Java-dan fərqi),
duck typing, polimorfizm (variadic interfeys parametri, Shape nümunəsi), "accept
interfaces, return structs" prinsipi (io.Reader ilə JSON decode), empty interface{}
(type assertion, comma-ok, type switch) və map[string]interface{} analizi +
Payer activity.

## Əsas fikirlər

### 1. İnterfeys Nədir
**Nədir:** metod İMZALARI toplusu — tipin DAVRANIŞINI təsvir edir; implementation
detalları YOXDUR. 4 təsvir: metod kolleksiyası, metod blueprint-u, davranış,
implementation-sız.

**Kitabın giriş problemi:** 3 eyni funksiya — loadEmployee(s string),
loadEmployeeFromFile(f *os.File), loadEmployeeFromHTTP(r *Request) — hamısı eyni
davranış (data OXU) → `loadEmployee(r io.Reader)` İLKİ funksiya KİMİ əvəz edir.

### 2. Təyini
```go
type Speaker interface {      // type + ad + interface
    Speak() string            // metod dəsti
}
```
**Adlandırma:** -er şəkilçisi (Reader, Speaker, Stringer); tək-metod interfeysi = metod
adı + er.

**Standart nümunələr:**
```go
type Reader interface {          // io — 1 metod
    Read(p []byte) (n int, err error)
}
type FileInfo interface {         // os — çoxsaylı
    Name() string
    Size() int64
    Mode() FileMode
    ...
}
```

### 3. Implicit Implementasiya
**Java:** `class Dog implements Pet` — AÇIQ bəyan MÜTLƏQ.
**Go:** implements açar sözü YOXDUR — tip interfeysin BÜTÜN metodlarını daşıyırısa,
AVTOMATİK satisfy edir.

**Kitabdan kod nümunəsi:**
```go
type Speaker interface { Speak() string }
type cat struct { name string; age int }

func (c cat) Speak() string { return "Purr Meow" }   // interfeys satisfied — SÖZSÜZ
func (c cat) String() string {                        // ƏLAVƏ metod — maneə DEYİL
    return fmt.Sprintf("%v (%v years old)", c.name, c.age)
}
// cat HƏM Speaker, HƏM fmt.Stringer (çapda String() çağırılır) implement edir
```
**Üstünlüklər:** interfeysi dəyişsən heç bir tipdə implements sözlərini dəyişmək
lazım deyil; BAŞQA paketin tiplərinə öz interfeysini yaza bilərsən (decoupling).

### 4. Duck Typing
"If it looks like a duck... then it must be a duck" — tip TIPLƏ yox, METODLARLA
müqayisə olunur. Speaker interfeysi hər hansı Speak() metodlu tip tərəfindən
satisfied → həmin tip hər yerdə Speaker kimi istifadə oluna bilər.

```go
func chatter(s Speaker) {          // interfeys PARAMETR
    fmt.Println(s.Speak())
}
chatter(cat{})                     // cat satisfy edir → KEÇİR
```

### 5. Polimorfizm — Müxtəlif Formalar
Go subclassing ETMİR (class YOX) — polimorfizm interfeyslərlə.

**Kitabdan refaktorinq nümunəsi (3 redundan funksiyadan 1-ə):**
```go
// ƏVVƏL: catSpeak(c), dogSpeak(d), personSpeak(p) — 3 eyni funksiya
// SONRA — variadic interfeys:
func saySomething(say ...Speaker) {
    for _, s := range say {
        fmt.Println(s.Speak())     // hər tipin ÖZ implementasiyası çağırılır
    }
}
saySomething(c, d, p)             // cat, dog, person — hamısı
```
**Gəlir:** bir dəfə yazılmış+test edilmiş kod hər satisfied edən tip üçün işləyir.

**Shape nümunəsi (kitabdan):**
```go
type Shape interface {
    Area() float64
    Name() string
}
type triangle struct{ base, height float64 }
func (t triangle) Area() float64 { return (t.base * t.height) / 2 }
func (t triangle) Name() string  { return "triangle" }
// rectangle{length,width}, square{side} — eyni interfeys

func printShapeDetails(shapes ...Shape) {
    for _, item := range shapes {
        fmt.Printf("The area of %s is: %.2f\n", item.Name(), item.Area())
    }
}
printShapeDetails(t, r, s)     // triangle 155.78, rectangle 200.00, square 100.00
```

### 6. "Accept Interfaces, Return Structs"
**Go atalar sözü + Postel's Law** ("what you accept" — liberal ol).

**Kitabdan müqayisə (JSON decode):**
```go
// MƏHDUD — yalnız string:
func loadPerson2(s string) (Person, error) {
    err := json.NewDecoder(strings.NewReader(s)).Decode(&p)
    ...
}
// FLEKSİBEL — istənilən mənbə:
func loadPerson(r io.Reader) (Person, error) {
    err := json.NewDecoder(r).Decode(&p)
    ...
}
```
loadPerson-i 3 mənbə DƏSTƏKLƏYİR:
- string → `strings.NewReader(s)`
- fayl → `os.Open("data.json")` (dönüş *os.File = io.Reader satisfied)
- HTTP → `request.Body`

**QAYTARMA struct/concrete olmalı:** interfeys qaytarsan istifadəçi metod dəstini
araşdırmalı; concrete qaytarsan — birbaşa istifadə; interfeys lazımdırsa istifadəçi
ÖZÜ yaradır.

### 7. Empty interface{}
0 metod → BÜTÜN tiplər avtomatik satisfied. İstənilən tip qəbul edən funksiya:
```go
func emptyDetails(s interface{}) {
    fmt.Printf("(%v, %T)\n", s, s)   // dəyər + CONCRETE tip
}
emptyDetails(99)        // (99, int)
emptyDetails(cat{name: "oreo"})   // ({oreo}, main.cat)
```

### 8. Type Assertion — `v := s.(T)`
```go
var str interface{} = "some string"
v := str.(string)              // uğur → v = "some string"
strings.Title(v)               // "Some String"

var num interface{} = 49
v := num.(string)              // UĞURSUZ → PANİK!

// TƏHLÜKƏSİZ forma — comma-ok:
v, isValid := str.(int)       // isValid=false → zero value, panic YOX
```
İstifadə: naməlum schema-lı JSON məlumatın emalı — tipə görə fərqli manipulyasiya.

### 9. Type Switch
```go
switch v := x.(type) {          // .(type) YALNIZ switch-də
case int:
    fmt.Printf("%v is int\n", v)
case string:
    fmt.Printf("%v is a string\n", v)
case bool:
    fmt.Printf("a bool %v\n", v)
default:
    fmt.Printf("Unknown type %T\n", v)
}
```
Case-lər DƏYƏR yox TİP yoxlayır; v case-in tipində olur.

### 10. map[string]interface{} Analizi (kitabdan)
```go
type record struct {
    key       string
    valueType string
    data      interface{}
}

func newRecord(key string, i interface{}) record {
    r := record{}
    r.key = key
    switch v := i.(type) {
    case int:
        r.valueType = "int"; r.data = v; return r
    case bool:
        r.valueType = "bool"; r.data = v; return r
    case string:
        r.valueType = "string"; r.data = v; return r
    case person:                        // ÖZ tipin də mümkün!
        r.valueType = "person"; r.data = v; return r
    default:
        r.valueType = "unknown"; r.data = v; return r
    }
}

m := map[string]interface{}{"person": p, "animal": a, "age": 54, ...}
rs := []record{}
for k, v := range m {
    rs = append(rs, newRecord(k, v))    // hər dəyərin tipi AŞKAR OLUNUR
}
```
animal case-də olmadığından "unknown" — naməlum tiplər üçün default tələbi.

### 11. Payer Activity — Polimorfizm + interface{}
Developer (HourlyRate × HoursWorkedInYear) və Manager (Salary + Salary×CommissionRate)
tipləri Payer interfeysini (Pay() (string, float64)) implement edir; payDetails(p Payer)
tək funksiya hər ikisini çap edir. Review map[string]interface{} — qiymətlər string
("Excellent"=5...) VƏYA int ola bilər → type switch ilə float-a normalize → ortalamə.

## Əsas terminlələr
- Interface — metod imzaları dəsti; davranış müqaviləsi
- Method Set — interfeysi satisfied edən metodlar qrupu
- Implicit Implementation — implements açar sözü olmadan satisfied etmə
- Duck Typing — metodlara görə tip uyğunluğu
- Polymorphism — eyni interfeysin müxtəlif formalarda çıxışı
- Stringer — fmt.String() interfeysi; çapın fərdiləşdirilməsi
- io.Reader/io.Writer — Go-nun ən işlək interfeysləri
- Accept Interfaces, Return Structs — API dizayn prinsipi
- Postel's Law — "what you accept"də liberal ol
- Empty Interface interface{} — 0 metod; hər tipi qəbul edir
- Type Assertion — s.(T) — concrete tipi çıxarma; panic/comma-ok
- Type Switch — switch v := x.(type) — çoxtipli yoxlama
- map[string]interface{} — naməlum schema JSON-un Go forması

## Praktik nətidə

(1) Eyni davranışı təkrarlayan funksiyalar görəndə dayan — interfeys çıxart
(io.Reader nümunəsi). (2) Implicit implementasiya: tipin metodları interfeysi
AVTOMATİK satisfy edir; əlavə metodlar maneə deyil. (3) API dizaynı: parametrdə
İNTERFEYS, qaytarmada STRUCT — caller istənilən mənbəni ötürə bilsin. (4) Standart
kitabxana interfeyslərini (io.Reader, fmt.Stringer) əvvəlcədən yoxla — özünü yeniden
icad etmə. (5) String() metodu bütün fmt çaplarını fərdiləşdirir. (6) Variadic
interfeys (...Speaker) — çoxtipli polimorf API. (7) interface{} qəbul edən funksiyada
type switch + default MÜTLƏQ — naməlum tip crash-ə yol verməsin. (8) Assertioni
həmişə comma-ok ilə — tək formal assert uğursuzda panic. (9) JSON naməlum schema →
map[string]interface{} + type switch ilə analiz. (10) Go-da polimorfizmin yeganə yolu
interfeysdir — class/subclass yoxdur, amma bu kifayətdir.

## Mənbə
Pages: 251-288 (PDF 284-323)
