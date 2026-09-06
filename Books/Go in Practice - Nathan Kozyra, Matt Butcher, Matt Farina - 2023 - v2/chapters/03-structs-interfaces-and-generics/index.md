# Chapter 3 — Structs, interfaces, and generics (Strukturlar, interfeyslər və generiklər)

## Bu chapter nədən bəhs edir?

Struct-lar (anonim, funksiya sahələri, metodlar, pointer receiver, anonim sahələr, tag-lər,
JSON encode), Go-nun OOP/funksional dillərlə müqayisəsi, interfeyslərlə tip məhdudiyyəti,
alias tiplər, empty interface/any, generics (union tiplər, generic funksiyalar, filter/map
patterns, tip təxmini ~, constraints paketi).

## Əsas fikirlər

### 1. Struct Əsasları
**Nədir:** İstənilən şeyi təmsil edən istifadəçi məlumat tipi; C-vari + OOP analogu (böyük
hərf = public, kiçik = paket-daxili).

**Anonim struct:** tək metodda lazım olan data üçün; test table-larda (Ch6) və sürətli
marshal üçün:
```go
animal := struct {
    name  string
    speak func() string      // funksiya sahəsi — first-class
}{
    name:  "cat",
    speak: func() string { return "meow" },
}
```
İnisializasiya: `value := struct{...}{...}` (dəyərlər bilinirsə təmiz) və ya `new(T)`
(zeroed pointer) — üstünlük zövq məsələsidir.

### 2. Funksiya Sahəsi vs Metod
**Funksiya sahəsi** struct tərifindədir, amma parent struct-a çıxışı YOXDUR. **Metod** —
receiver ilə bağlı funksiya; parent-in sahələrinə çatır:

**Kitabdan kod nümunəsi:**
```go
type Animal struct { name string }
func (a Animal) speak() string {          // value receiver — dəyişmir
    switch a.name {
    case "cat":  return "meow"
    case "dog":  return "woof"
    default:     return "nondescript animal noise?"
    }
}
// Pointer receiver — dəyişdirir:
type character struct{ name string }
func (ch *character) fixName() { ch.name = "Inigo Montoya" }
```
**Metod vs funksiya seçimi:** çox parametr keçirirsənsə receiver data axını azaldır;
mutasiya üçün pointer receiver məcburidir; struct-ın bütün sahələrini funksiyaya ötürmək
yerinə metoddan istifadə access-i idarə edir.

### 3. Anonim (Adsız) Sahələr
Hər tipdən YALNIZ BİR olmaqla sahə adı buraxıla bilər — tip özü ad kimi:
```go
type Animal struct{ string }   // a.string kimi çıxış
```
Adlı sahə tip ilə çıxış edilməz (ad varsa). Literal-da son vergül məcburidir.

### 4. Struct Tag-lər
**Nədir:** Sahələrə metadata; backtick içində `key:"value"` cütləri; əsasən encoding.

```go
type Animal struct {
    name string `help:"the name or type of any animal, as long as it is a cat or dog"`
}
// Tag-in runtime oxunması — reflection:
member, ok := reflect.TypeOf(a).FieldByName("name")
member.Tag.Get("help")
```
Reflection runtime + mürəkkəblik xərcli — ilk seçim olmasın; dəyişkən formatlı JSON
(string VƏ YA number sahə) təhlilində görünür.

**JSON encode:**
```go
type Animal struct {
    Name           string  `json:"animal_name"`
    ScientificName string  `json:"scientific_name"`
    Weight         float32 `json:"animal_average_weight"`
}
output, _ := json.Marshal(a)   // {"animal_name":"cat",...}
```
Böyük hərf məcburidir (encoding/json paket xaricindən baxır); tagsiz sahə adı olduğu
kimi çıxar. XML və digər formatlar eyni pattern.

### 5. Go vs OOP vs Funksional
**Go multiparadigm-dir:** Java-nın `class Cat extends Animal` iyerarxiyası YOX; interfeys
ilə eyni məntiq:
```go
type Animal interface { speak() }
type Cat struct{}
func (c Cat) speak() { fmt.Println("meow") }
type Llama struct{}   // speak YOXDUR →
// *Llama does not implement Animal (missing speak method) — COMPILE xətası!
```
Implicit implementasiya: `extends/implements` YOX — ya uyğundur, ya xəta. Bu, dizayn
sualı da daxil edir: "Llama danışmırsa, onu Animal hesab etmək olarmı?" (duck typing).
Constructor konvensiyası: `NewCat() *Cat`. Superclass fallback-u yoxdur — metod hamıda
olmalı və ya tip interfeysə düşməməlidir. Ekstra metodlar (woof, meow) sərbəstdir.
Funksional baxışdan: funksiyalar first-class-dir, amma generics-dən əvvəl tipəmən iterator
yazmaq mümkün deyildi; idiomatik Go funksional estetikaya uyğun deyil — state ötürmək
üçün dizayn edilib.

### 6. Interfeyslə Builtin Tiplərə Funksionallıq (Alias Pattern)
Builtin tipə metod YAZILA BİLMƏZ → alias tip yarat:
```go
type shuffleString string
func (s *shuffleString) shuffle() { /* strings.Split + rand.Shuffle + Join */ }
func (s *shuffleString) contents() string { return string(*s) }
type shuffleSlice []interface{}
func (sl shuffleSlice) shuffle() { rand.Shuffle(len(sl), ...) }

type Shuffleable interface {
    contents() string
    shuffle()
}
var myShuffle Shuffleable
myShuffle = NewShuffleString("my name is inigo montoya")
myShuffle = &shuffleSlice{1, 2, 3, 4, 5}   // hər ikisi eyni interfeys
```
**Empty interface (`interface{}` = `any`):** 0 metod → hər tip uyğun; `[]interface{}`
slice-da qarışıq tiplər. Amma interface qaytarmaq/qəbul etmək tip müəyyənliyi işi
yaradır — mümkünsə konkret tip qaytar (sadə log istisna).

### 7. Generics — Go 1.18
**Tarix:** 10+ il "mürəkkəblik/xərc qazancı aşmaz" müqayisəsindən sonra 1.18-də gəldi.
Eyni alqoritmi çox tipə yazmaq təkrarını aradan qaldırır.

**Union type + instantiation:**
```go
type AnimalType interface {
    Cat | Dog                     // union — bu tiplər və yalnız bunlar
}
type Animal[T AnimalType] struct {
    value       T
    AnimalNoise func() string
}
func (a Animal[T]) Speak() { fmt.Println(a.AnimalNoise()) }
catAnimal := Animal[Cat]{value: Cat{}, AnimalNoise: func() string { return "meow!" }}
```
`T` konvensiyası C++/Java şablonlarından; `[Cat]` instantiation çox vaxt inference ilə
lazımsızdır (`Speak(catAnimal)` kifayət). **Antipattern xəbərdarlığı:** generic funksiya
daxilində tip-ə görə switch ilə məntiq şaxələndirmək — kodu bulandırır.

**Generic funksiyalar — ən real fayda slice/map üzərində:**
```go
func filter[T any](items []T, fx func(T) bool) []T {
    var filtered []T
    for _, v := range items {
        if fx(v) { filtered = append(filtered, v) }
    }
    return filtered
}
// istifadə — tip göstərilməsi inference ilə azalır:
strings = filter[string](strings, func(s string) bool { return unicode.IsUpper(rune(s[0])) })
ints = filter[int](ints, func(i int) bool { return i%3 == 0 })
```

**Tip constraint (daha dar):**
```go
type Numeric interface {
    int8 | int16 | int32 | int64 | float32 | float64
}
func filterPositive[T Numeric](items []T) []T { ... }
```
`any`-dən dar constraint type-safety artırır.

**Mühüm məhdudiyyət:** Generic funksiya daxilində sahələrə çıxış YOXDUR — hətta bütün
namizəd tiplər eyni sahəni paylaşsa belə. Sahə/metod lazımdırsa generics yanlış seçim ola bilər.

### 8. Tip Təxmini (~) və constraints Paketi
**Problem:** `type Smallint int8` → `doubler[T Numeric]`-ə KEÇMİR ("possibly missing ~
for int8") — union EXACT tip tələb edir, underlying type YOX.

**Həll — tilde:**
```go
type Numeric interface {
    ~int8 | ~int16 | ~int32 | ~int64 | ~float32 | ~float64
}
```
`~int8` = int8 VƏ underlying type-ı int8 olan bütün tiplər (type approximation).

**constraints paketi** (golang.org/x/exp/constraints — experimental): hazır dəstlər:
`constraints.Signed/Unsigned/Integer/Ordered` (Ordered string daxil). Rare hallar
istisna, hamısı underlying-lə birlikdə istənilir:
```go
func doubler[T constraints.Integer](value T) T {
    return value * T(2)
}
```

## Əsas terminlələr
- Receiver — metodun bağlılıq dəyişəni `(a Animal)`
- Anonymous struct/sahə — adsız tip/sahə; sahədə tip adı ilə çıxış
- Struct tag — `json:"name"` metadata; reflection ilə oxunur
- Duck typing — metod uyğunluğuna görə tip münasibəti
- Empty interface / any — 0 metodlu interfeys; hər tip
- Union type — `Cat | Dog` interfeys daxili tip məhdudiyyəti
- Instantiation — `Animal[Cat]` tip arqumenti təyini
- Type approximation — `~int8` underlying tiplər daxil edən constraint
- constraints paketi — hazır Signed/Integer/Ordered dəstləri (experimental)

## Praktik nəticə

Struct/interfeys/generics qərar ağacı: (1) datanı struct-da saxla; (2) funksiya struct-a
daxili state çıxışı lazımdırsa METOD (dəyişmə = pointer receiver); (3) ortak davranış
müqaviləsi = interfeys (client tərəfdə, implicit); (4) builtin tipə metod = alias tip;
(5) eyni alqoritm çox slice/map tipi üçün = generic funksiya (`[T any]` və ya dar constraint);
(6) custom tiplər underlying dəstəyə düşməlidirsə `~`; (7) sahə çıxışı lazım olan
"generics" = interfeys istifadə et; (8) `any`-ni sərhəd kimi saxla — qaytarmada konkret
tip üstün.

## Mənbə
Pages: 61-84 (PDF 82-105)
