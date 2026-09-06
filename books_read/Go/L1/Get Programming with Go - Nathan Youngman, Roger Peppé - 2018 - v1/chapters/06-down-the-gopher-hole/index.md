# Unit 6 — Down the Gopher Hole (Lesson 26-29)

## Bu unit nədən bəhs edir?

Pointer-lər (& ve *, tip elanı, struct/array ilə ergonomiya, mutation, interior pointer, pointer-in arxasındakıları — map/slice), nil (nil pointer panic, guard clause, nil funksiya/slice/map/interface, nil alternativi), error handling (çoxqayıdış konvensiyası, defer, safeWriter pattern, errors.New, error dəyişənləri, custom error tipləri, type assertion, panic/recover). Capstone: Sudoku qaydaları.

**PDF səhifələr:** 216-266 (L26: 216-229, L27: 235-244, L28: 245-262, L29: 263-266)

## Əsas fikirlər

### 1. Pointer əsasları (L26)
**Nədir:** Başqa dəyişənin ÜNVANINI saxlayan dəyişən. "Köçürük" nişanı metaforası — indirection.

**2 simvol:**
- `&` (ampersand) — **address operator**: dəyişənin yaddaş ünvanı
- `*` (asterisk) — **dereference**: ünvanın göstərdiyi DƏYƏR (tək dəyişən prefiksində) YAXUD tip prefiksində pointer TİPİ

```go
answer := 42
fmt.Println(&answer)        // 0x1040c108 — ünvan
address := &answer          // address: *int tipli
fmt.Println(*address)       // 42 — dereference
fmt.Println(*&answer)      // 42
```
- `&42` → compile xətası (literal-ların ünvanı alınmaz; composite literal-ların ALINA BİLƏR!)
- Go-da **pointer aritmetikası YOXDUR** (`address++` qadağan — C-dən fərqli, təhlükəsizlik)

### 2. Pointer tipi (L26)
```go
canada := "Canada"
var home *string            // tip: *string
fmt.Printf("%T\n", home)    // *string
home = &canada              // yalnız string dəyişənlərə işarə edə bilər
fmt.Println(*home)          // Canada
```

### 3. Pointer-lar dəyişənliyi (L26) — NASA nümunəsi
```go
var administrator *string
scolese := "Christopher J. Scolese"
administrator = &scolese
fmt.Println(*administrator)     // Scolese

bolden := "Charles F. Bolden"
administrator = &bolden           // yenidən yönəlt
fmt.Println(*administrator)     // Bolden

bolden = "Charles Frank Bolden Jr."
fmt.Println(*administrator)     // yenilənmiş — pointer canlı bağlıdır!

*administrator = "Maj. Gen. ..."   // TƏRS istiqamət: pointer ilə bolden-i dəyiş
major := administrator            // pointer kopyası — EYNİ ünvan
*major = "Major General ..."       // hər ikisi təsirlənir
fmt.Println(administrator == major)   // true — eyni ünvan

lightfoot := "Robert M. Lightfoot Jr."
administrator = &lightfoot
fmt.Println(administrator == major)    // false — artıq fərqli ünvanlar

charles := *major                  // DEREFERENCE ilə kopya — müstəqil!
fmt.Println(charles == bolden)     // true (dəyər eynidir)
fmt.Println(&charles == &bolden)   // false (ünvanlar fərqli)
```

### 4. Struct pointer-ları — ergonomiya (L26)
```go
type person struct {
    name, superpower string
    age              int
}

timmy := &person{              // composite literal + & — ICAZƏLİDİR!
    name: "Timothy",
    age:  10,
}
timmy.superpower = "flying"     // AVTOMATİK dereference — (*timmy).superpower lazım DEYİL
fmt.Printf("%+v\n", timmy)     // &{name:Timothy superpower:flying age:10}
```

**Array pointer-ları:**
```go
superpowers := &[3]string{"flight", "invisibility", "super strength"}
fmt.Println(superpowers[0])       // flight — avtomatik dereference
fmt.Println(superpowers[1:2])     // [invisibility] — slicing də OK
```
Amma slice/map-ların pointer-ına avtomatik dereference YOXDUR.

### 5. Mutation — pointer parametr (L26)
```go
func birthday(p *person) {
    p.age++                     // dolayısı ilə mutasiya
}

rebecca := person{name: "Rebecca", superpower: "imagination", age: 14}
birthday(&rebecca)
fmt.Printf("%+v\n", rebecca)    // age:15 — DƏYİŞDİ!
```
**Pass by value + pointer:** funksiya ünvanın KOPYASINI alır — amma kopya da eyni yaddaş yerinə işarə edir.

### 6. Pointer receiver (L26)
```go
func (p *person) birthday() {
    p.age++
}

terry := &person{name: "Terry", age: 15}
terry.birthday()     // pointer ilə çağırış ✓

nathan := person{name: "Nathan", age: 17}
nathan.birthday()     // pointer SİZ value — Go &nathan edir! ✓
```
**Vacib:** receiver *person olmalı — value receiver olsaydı yaşı 15-dən artırmazdı (kopya üstündə işləyərdi).

**time.Time dərsi:** `day.Add(24 * time.Hour)` — pointer receiver YOX! time.Time **dəyişməz**dir — Add YENİ dəyər qaytarır. "Tomorrow is a new day."

**Qayda:** Pointer receiver-ları TUTUMLI istifadə et — 1 metod pointer istifadə edirsə, hamısı etsin.

### 7. Interior pointer-lar (L26)
```go
type stats struct {
    level             int
    endurance, health int
}
func levelUp(s *stats) {
    s.level++
    s.endurance = 42 + (14 * s.level)
    s.health = 5 * s.endurance
}

type character struct {
    name  string
    stats stats
}
player := character{name: "Matthias"}
levelUp(&player.stats)          // STRUKTURUN SAHƏSİNƏ pointer!
fmt.Printf("%+v\n", player.stats)   // {level:1 endurance:56 health:280}
```
Sahələr pointer elan olunmamış — amma lazım olanda `&player.stats` interior pointer alınır.

### 8. Array mutasiyası (L26)
```go
func reset(board *[8][8]rune) {
    board[0][0] = 'r'
}
var board [8][8]rune
reset(&board)          // pointer → mutasiya mümkün
```

### 9. Pointer-in arxasında — map/slice (L26)
**Map = pointer-in gizli forması:**
```go
func demolish(planets *map[string]string)   // QƏBUL EDİLMƏZ! map artıq pointerdir
```
**Slice = daxili 3 element:** array-ə pointer + length + capacity. Elementlər birbaşa mutasiya olunur; pointer-a ehtiyac yalnız slice-in ÖZÜNÜ (len/cap/offset) dəyişmək olanda:
```go
func reclassify(planets *[]string) {
    *planets = (*planets)[0:8]     // uzunluğu dəyişir — pointer LAZIM
}
```
(Daha təmiz: yeni slice qaytar.)

### 10. Pointer + interface (L26)
```go
// Value receiver: hər ikisi interfeysi təmin edir
func (m martian) talk() string { ... }
shout(martian{})      // ✓
shout(&martian{})     // ✓ — pointer hər şeyi təmin edir

// Pointer receiver: YALNIZ pointer!
func (l *laser) talk() string { ... }
pew := laser(2)
shout(&pew)     // ✓
shout(pew)      // XƏTA!
```

### 11. Nil (L27)
**Tony Hoare:** "Null References: The Billion Dollar Mistake" (1965 ixtirası; 2009 etiraf). Amma o, CSP (1978) da icad etdi — Go-nun concurrency əsası!

**Nil pointer dereference:**
```go
var nowhere *int
fmt.Println(nowhere)    // <nil>
fmt.Println(*nowhere)   // PANIC: nil pointer dereference

if nowhere != nil {     // guard
    fmt.Println(*nowhere)
}
```

### 12. Nil receiver guard (L27)
**Go xüsusiyyəti:** nil receiver-də metod çağırmaq JAVA-dan fərqli olaraq dərhal crash ETMİR — panic yalnız SAHƏYƏ ÇIXIŞDA baş verir:
```go
func (p *person) birthday() {
    p.age++               // ← p nil olsa BURADA panic
}

var nobody *person
nobody.birthday()          // metod çağrılır — p.age++-da panic

// Guard clause — metodun İÇINDƏ:
func (p *person) birthday() {
    if p == nil {
        return               // səssizcə keç
    }
    p.age++
}
```
Sənin qərarın: zero value qaytar / error qaytar / crash et — Go seçimi sənə buraxır.

### 13. Nil funksiya dəyərləri (L27)
```go
var fn func(a, b int) int
fmt.Println(fn == nil)    // true
// fn(1, 2) → panic

// Default funksiya patterni:
func sortStrings(s []string, less func(i, j int) bool) {
    if less == nil {
        less = func(i, j int) bool { return s[i] < s[j] }   // default!
    }
    sort.Slice(s, less)
}
sortStrings(food, nil)      // nił ötür → default istifadə olunur
```

### 14. Nil slice / nil map (L27)
**Nil slice — demək olar hər şey LEGAL:**
```go
var soup []string
fmt.Println(soup == nil)     // true
for _, ingredient := range soup { ... }   // OK — 0 iterasiya
fmt.Println(len(soup))       // 0
soup = append(soup, "onion", "carrot", "celery")   // append işləyir!

// nil slice = empty slice qədər istifadə oluna bilər:
soup := mirepoix(nil)         // make edəəək ehtiyac yoxdur
```

**Nil map — OXU legal, YAZI panic:**
```go
var soup map[string]int
fmt.Println(soup == nil)      // true
measurement, ok := soup["onion"]    // OK — oxu
for ingredient, measurement := range soup { ... }   // OK — 0 iterasiya
soup["onion"] = 1             // PANIC: assignment to entry in nil map
```

### 15. Nil interface — sürpriz! (L27)
```go
var v interface{}
fmt.Printf("%T %v %v\n", v, v, v == nil)
// <nil> <nil> true — hər ikisi nil

var p *int
v = p
fmt.Printf("%T %v %v\n", v, v, v == nil)
// *int <nil> FALSE!!! — interface tip + dəyər SAXLAYIR
fmt.Printf("%#v\n", v)    // (*int)(nil) — tip dolu, dəyər nil
```
**Dərs:** interface == nil yalnız HƏR İKİSİ nil olanda true. Nil pointer-i interface-ə qoyursansa — `v == nil` FALSE olacaq. Müqayisələrdə `nil` identifikatorunu açıq yaz.

### 16. Nil-ə alternativ (L27)
```go
type number struct {
    value int
    valid bool            // nil əvəzinə flag!
}
func newNumber(v int) number {
    return number{value: v, valid: true}
}
func (n number) String() string {
    if !n.valid {
        return "not set"
    }
    return fmt.Sprintf("%d", n.value)
}
```
Valid bool = niyyət açıq; nil pointer = məna qeyri-müəyyən. Pointer yalnız pointing üçün!

### 17. Error handling — əsaslar (L28)
**Go fəlsəfəsi:** "Errors aren't exceptional" — xətalar istisna deyil, gözləniləndir. Çoxqayıdışlı dillərdən fərq: error həmişə **sonuncu return value**, caller dərhal yoxlayır:
```go
files, err := ioutil.ReadDir(".")
if err != nil {
    fmt.Println(err)
    os.Exit(1)
}
for _, file := range files {
    fmt.Println(file.Name())
}
```
- Xəta olanda digər dəyərlərə ETİBAR ETMƏ (zero yaxud qismən data ola bilər)
- Go Proverbs: "Errors are values. Don't just check errors, handle them gracefully. Don't panic."

### 18. defer — təmizləmə (L28)
```go
func proverbs(name string) error {
    f, err := os.Create(name)
    if err != nil {
        return err
    }
    defer f.Close()          // HƏR return-dan əvvəl icra olunacaq
    _, err = fmt.Fprintln(f, "Errors are values.")
    if err != nil {
        return err            // defer Close çağırır
    }
    ...
    return err
}
```
- defer → funksiya qayıtmazdan əvvəl; panic belə olsa; fayl açılan yerdən dərhal sonra yazılır — unutma ehtimalı 0

### 19. safeWriter — kreativ error pattern (L28)
**Rob Pike "Errors are values" (Go blog, 2015):**
```go
type safeWriter struct {
    w   io.Writer
    err error                     // İLK xətanı saxla
}
func (sw *safeWriter) writeln(s string) {
    if sw.err != nil {
        return                     // əvvəlki xəta varsa yazma
    }
    _, sw.err = fmt.Fprintln(sw.w, s)
}

func proverbs(name string) error {
    f, err := os.Create(name)
    if err != nil {
        return err
    }
    defer f.Close()
    sw := safeWriter{w: f}
    sw.writeln("Errors are values.")
    sw.writeln("Don't just check errors, handle them gracefully.")
    // ... 13 sətir — heç bir if err yoxdur!
    return sw.err
}
```
**Böyük ideya:** "errors are values, and the full power of the Go programming language is available for processing them" — hər sətirdə if yazmaq əvəzinə xəta yığıcı tip yarat.

### 20. errors.New + error dəyişənləri (L28)
```go
// Sudoku nümunəsi:
const rows, columns = 9, 9
type Grid [rows][columns]int8

var (
    ErrBounds = errors.New("out of bounds")     // konvensiya: Err prefiksi
    ErrDigit  = errors.New("invalid digit")
)

func (g *Grid) Set(row, column int, digit int8) error {
    if !inBounds(row, column) {
        return ErrBounds               // paket səviyyəli dəyişən
    }
    g[row][column] = digit
    return nil
}

// Caller — dəqiq identifikasiya:
err := g.Set(0, 0, 15)
if err != nil {
    switch err {
    case ErrBounds, ErrDigit:
        fmt.Println("Les erreurs de paramètres hors limites.")
    default:
        fmt.Println(err)
    }
}
```
- **Müqayisə ünvanlar üzərindədir** (errors.New pointer implementasiyası) — mətn yox!
- Parametrləri funksiyanın BAŞINDA validasiya et — qalan kod təmiz qalır

### 21. Custom error tipləri (L28)
**error interfeysi:**
```go
type error interface {
    Error() string
}
```
Hər hansı tip `Error() string` metoduna sahibdirsə — error-dur!

**Çoxlu xəta toplayıcı:**
```go
type SudokuError []error                     // konvensiya: Error şəkilçisi

// Error returns one or more errors separated by commas.
func (se SudokuError) Error() string {
    var s []string
    for _, err := range se {
        s = append(s, err.Error())
    }
    return strings.Join(s, ", ")
}

func (g *Grid) Set(row, column int, digit int8) error {
    var errs SudokuError
    if !inBounds(row, column) {
        errs = append(errs, ErrBounds)
    }
    if !validDigit(digit) {
        errs = append(errs, ErrDigit)
    }
    if len(errs) > 0 {
        return errs              // BİRDƏN-ÇOX xəta!
    }
    g[row][column] = digit
    return nil                    // BOŞ slice YOX — nil! (nil interface dərsi)
}
```
**Vacib:** metod imzası həmişə `error` interfeysi — konkret `SudokuError` YOX. Uğurda boş slice YOX, nil qaytar (yoxsa "nil olmayan nil" tələsi).

### 22. Type assertion (L28)
```go
err := g.Set(10, 0, 15)
if err != nil {
    if errs, ok := err.(SudokuError); ok {     // tip iddiası!
        fmt.Printf("%d error(s) occurred:\n", len(errs))
        for _, e := range errs {
            fmt.Printf("- %v\n", e)
        }
    }
}
// 2 error(s) occurred:
// - out of bounds
// - invalid digit
```
`err.(SudokuError)` → interface-dən konkret tipə çevirmə; `ok` — müvəffəqiyyət.

### 23. Panic və recover (L28)
**Exceptions ilə müqayisə:** Go-da exceptions YOX — xətalar dəyərlərdir; exceptions opt-out çətin, Go-da ignore etmək AÇIQ görünür.

```go
// panic — proqram çökür:
panic("I forgot my towel")

// Runtime panic-lər: 42 / zero → "integer divide by zero"

// recover — YALNIZ defer daxilində:
defer func() {
    if e := recover(); e != nil {
        fmt.Println(e)          // panic dayandırıldı
    }
}()
panic("I forgot my towel")
// İcazə verilir: panic > os.Exit (defer-lər icra olunur!)
```
**Qayda:** panic NADİR olsun — error qaytar.

### 24. Capstone: Sudoku rules (L29)
- 9x9 array (fixed) + pointer mutasiya
- `Set(r, c, digit)` — row/column/3x3 subregion qaydaları + error qaytar
- `Clear(r, c)` — qaydasız
- `NewSudoku(...)` constructor + composite literal
- Fixed (başlanğıc) vs penciled rəqəmlər — fixed-lər dəyişilməz (error)

## Error handling müqayisəsi
| | Exceptions (Java/Python) | Go |
|---|---|---|
| Mexanizm | try/catch/throw xüsusi sözlər | Çoxqayıdış + error dəyəri |
| Default | Ignore (opt-in tutma) | Açıq yoxlama (ignore AÇIQ görünür) |
| Yayılmа | Call stack üzrə bubble | Dəyər kimi ötürülür |
| Xüsusi kod | try/catch/finally blokları | adi if/switch/for |
| Crash | Unhandled exception | panic (nadir) |

## Əsas terminlər
- Address Operator (`&`)
- Dereference (`*`)
- Pointer Type (`*int`)
- Pointer Arithmetic (Go-da YOX)
- Indirection
- Automatic Dereference (struct sahə/array)
- Interior Pointer (`&player.stats`)
- Pass by Value + Pointer Copy
- Pointer Receiver (tutum qaydası)
- Immutable Type (time.Time nümunəsi)
- nil (zero value: pointer/slice/map/interface)
- Nil Pointer Dereference / Panic
- Guard Clause
- Billion Dollar Mistake (Tony Hoare)
- Nil Function Value (default pattern)
- Nil Slice / Nil Map (oxu-yazı fərqi)
- Interface = type + value (nil tələsi)
- Sentinel Error (`ErrBounds`)
- error Interface (`Error() string`)
- Custom Error Type (`SudokuError`)
- Multiple Errors
- defer (təmizləmə zəmanəti)
- safeWriter Pattern ("errors are values")
- Type Assertion (`err.(SudokuError)`)
- panic / recover (yalnız defer-də)
- Go Proverbs

## Praktik nəticə
- Pointer-ları struct və array mutasiyası üçün istifadə et; map/slice üçün YOX (onlar artıq pointerdirlər).
- Pointer receiver 1 metodda lazımdırsa — hamıda istifadə et; dəyişməz tiplərdə (time.Time) value receiver.
- Nil receiver-də metodu guard et — Go çağırışı crash etmədən icazə verir.
- Nil map-ə yazma panic; nil slice append/reflection/len ilə problemsiz — funksiyalarını nil-friendly yaz.
- Interface-i nil ilə müqayisə edəndə HƏR İKİ sahənin nil olmasını yoxla; `(*int)(nil)` tələsinə diqqət.
- Xətaları saxla/değiştir — errors.New dəyərləri `Err` prefiksi ilə paket səviyyəsində; caller switch ilə dəqiq identifikasiya edir.
- Təkrarlanan `if err != nil`-lərdən yorulmusan — safeWriter kimi xəta yığıcı tip yarat ("full power of Go").
- Uğur halında custom error funksiyası nil qaytarsın — boş slice "nil olmayan nil" xətası verər.
- panic nadir; crash-dan qaçmaq lazımdırsa defer + recover (yalnız defer daxilində işləyir).

## Mənbə
Pages: 216-266 (PDF), book pages 201-251 (Lesson 26-29)
