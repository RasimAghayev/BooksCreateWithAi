# Chapter 4 — Complex Types (səh. 134-191)

## Bu fəsil nədən bəhs edir?

Kolleksiyalar: array (fixed, müqayisəli, açar ilə init), slice (gizli
massiv, len/cap, make, append, 5 kopyalama ssenarisi), map (hash, ok
idiomu, delete). Custom type-lar, struct (init formaları, müqayisə,
embedding/promotion), tip çevrilmələri (lossy), interface{} / any, type
assertion və type switch.

## Əsas fikirlər

### 1. Array-lər
**Nədir:** `[<size>]<type>` — ölçü TİPİN bir hissəsidir; compile-da sabit.

```go
var arr [10]int                            // zero-larla dolu
arr2 := [5]int{0}                          // {0,0,0,0,0}
arr3 := [...]int{0, 0, 0, 0, 0}           // ... = ölçünü öz çıxarır (5)
arr4 := [10]int{1, 9: 10, 4: 5}            // AÇAR ilə: 0:1, 4:5, 9:10, qalanları 0
```

**Müqayisə (yalnız eyni tip = eyni ölçü!):**
```go
arr1 == arr2   // [5]int == [5]int{0}      → true
arr1 == arr4   // [5]int vs [9]int{...}    → COMPILE XƏTASI — fərqli TİP
```
- Slice/map müqayisə OLMAZ — array müqayisəsi onun ÜSTÜNLÜYÜDÜR

**Oxu/yazma/loop:**
```go
arr := [4]string{"ready", "Get", "Go", "to"}
fmt.Sprintln(arr[1], arr[0], arr[3], arr[2])   // "Get ready to Go"
arr[1] = "It's"
for i := 0; i < len(arr); i++ { ... }    // for i loop
arr[i] = arr[i] * arr[i]                 // loop-da dəyişmə
```
- len runtime-da SAYMIR — Go uzunluğu izləyir (sürətli)
- Hardcode ölçü = gizli bug; həmişə len

### 2. Slice-lər
**Nədir:** array-in üzərində nazit qat; ölçü dinamik; append ilə böyüyür.

```go
// Klassik istifadə (os.Args nümunəsi):
func getPassedArgs(minArgs int) []string {
    if len(os.Args) < minArgs {
        fmt.Printf("At least %v arguments are needed\n", minArgs)
        os.Exit(1)
    }
    var args []string
    for i := 1; i < len(os.Args); i++ {    // [0] = proqramın öz adı!
        args = append(args, os.Args[i])
    }
    return args
}

func findLongest(args []string) string {
    var longest string
    for i := 0; i < len(args); i++ {
        if len(args[i]) > len(longest) {
            longest = args[i]
        }
    }
    return longest
}
```

**Çoxlu append (variadic):**
```go
locales = append(locales, "en_US", "fr_FR")   // bir neçə dəyər
locales = append(locales, extraLocales...)   // slice-i AÇ (... explode)
```

**Slice-dan slice (range notasiyası):**
```go
s := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
s[0:1]      // [1]     — low və high
s[:1]       // [1]     — low default 0
s[len(s)-1:]          // [9]   — high default len
s[:5]       // ilk 5
s[5:]       // son 4
s[2:7]      // ortadan 5
s[:]        // hamısı (array→slice çevirmək üçün)
```
- **KOPYA YOXDUR** — yeni slice HƏMİŞƏ gizli massivə BAXIŞdır ("view")!

**Slice daxili mexanika (len/cap/make):**
- Gizli strukt: pointer (gizli massivə) + uzunluq (len) + başlanğıc
 yeri; cap = gizli massivin ölçüsü
- append: cap boşdursa → len artır; doludursa → yeni daha BÖYÜK massiv
 + kopyalama + pointer dəyişir
```go
var s1 []int                // len=0 cap=0
s2 := make([]int, 10)       // len=10 cap=10
s3 := make([]int, 10, 50)   // len=10 cap=50 — genişlənmə üçün hazır!
```

**5 kopyalama ssenarisi (link həqiqəti!):**
```go
// 1) LINKED — eyni gizli massiv:
s1 := []int{1, 2, 3, 4, 5}
s2 := s1
s3 := s1[:]
s1[3] = 99                  // s2[3] və s3[3] DA 99!

// 2) NO LINK — append cap-i aşdı → yeni massiv:
s2 := s1
s1 = append(s1, 6)          // gizli massiv dəyişdi
s1[3] = 99                  // s2[3] = 4 (əski massivdə qalır)

// 3) CAP LINK — artıq cap var, append yeni massiv YARATMADI:
s1 := make([]int, 5, 10)
s2 := s1
s1 = append(s1, 6)          // cap içində — eyni massiv
s1[3] = 99                  // s2[3] = 99 — HƏLƏ LİNK!

// 4) CAP NO LINK — append cap-dən çox:
s1 = append(s1, []int{10: 11}...)  // massiv dəyişdi → link qırıldı

// 5) COPY/APPEND — müstəqil kopya:
s2 := make([]int, len(s1))
copied := copy(s2, s1)      // copy uzunluğu DƏYİŞMİR — əvvəlcədən boy ver!
s2 := append([]int{}, s1...)            // ən çox görülən üsul
s2 := append(s1[:0:0], s1...)           // ƏN effektiv (3-indeksli range)
```

### 3. Map-lər
**Nədir:** hashmap — açar MƏLUMATDIR (placeholder yox); O(1) axtarış.

```go
users := map[string]string{
    "305": "Sue",
    "204": "Bob",
    "631": "Jake",
}
users["073"] = "Tracy"          // əlavə
// var ilə tanımlanan map-a yazmaq → RUNTIME PANIC! (make/literal şərt)
```

**Ok idiomu (açar mövcudluğu):**
```go
user, exists := users[id]        // exists=false → user = zero value ""
if !exists {
    // tapılmadı
}
```
- Açar yoxdursa → zero value qaytarır (xəta YOX — dizayn qərarı)
- range sırası QƏSDƏN random

**Delete:**
```go
delete(users, "305")     // yoxdursa — heç nə olmur
```

### 4. Custom type-lar
```go
type id string            // string əsaslı yeni tip

var id1 id                // zero: ""
var id2 id = "1234-5678"
id3 := id("4567")         // çevrilmə ilə

string(id2) == "1234-5678"   // əsas tipə QAYITMA mümkün
```

### 5. Struct-lar
```go
type user struct {
    name    string
    age     int
    balance float64
    member  bool
}

// Init formaları:
u1 := user{name: "Tracy", age: 51, balance: 98.43, member: true} // açar-dəyər (ən çox)
u2 := user{age: 19, name: "Nick"}       // qalanlar zero
u3 := user{"Bob", 25, 0, false}         // sıralı — HAMISI olmalı
var u4 user                             // hamısı zero
u4.name = "Sue"                         // . ilə yaz

// Anonim struct (bir dəfəlik):
point1 := struct{ x, y int }{10, 10}
```

**Müqayisə:** bütün sahələr comparable → struct da comparable; anonim
struct eyni strukturla müqayisə OLUR (go-nun çevikliyi).

**Embedding (irsiyyət əvəzi):**
```go
type name string
type location struct { x, y int }
type size struct { width, height int }

type dot struct {
    name          // adsız sahələr = EMBED
    location
    size
}

// PROMOTION — daxili sahələr birbaşa əlçatan:
dot2.x = 5
dot2.width = 10
dot2.name = "A"

// Init-də promotion YOX — tip adı ilə:
dot3 := dot{
    name: "B",
    location: location{x: 13, y: 27},
    size:    size{width: 5, height: 7},
}
// Və ya tam yol:
dot4.location.x = 101
```
- Ad toqquşması olsa promotion BAŞ VERMİR (tip yolu ilə çat)
- Real koddə embedding az-çox — adlı sahə üstünlük verilir

### 6. Tip çevrilmələri (lossy!)**
```go
var i8 int8 = math.MaxInt8     // implicit: int → int8 (yığılır)
int64(i8)                       // kiçik → böyük: TƏHLÜKƏSİZ
int8(128)                       // int → int8: OVERFLOW → -128!
float64(i8)                     // int → float: dəyişməz
int(3.14)                        // float → int: kəsr ATILIR → 3
uint → int, böyük → kiçik: eyni risklər
```

### 7. interface{} / any + type assertion
**Nədir:** fmt.Print necə istənilən tipi qəbul edir? `a ...interface{}`
— boş interfeys: HƏR TİP uyğun gəlir. Go 1.18+: `any` alias.

```go
// fmt.Print-in öz kodu:
func Print(a ...interface{}) (n int, err error) {
    return Fprint(os.Stdout, a...)
}
```

**Type assertion (kilidi açmaq):**
```go
func doubler(v interface{}) (string, error) {
    if i, ok := v.(int); ok {              // ok olmadan → PANIC riski!
        return fmt.Sprint(i * 2), nil
    }
    if s, ok := v.(string); ok {
        return s + s, nil
    }
    return "", errors.New("unsupported type passed")
}
```

**Type switch:**
```go
switch t := v.(type) {
case string:
    return t + t, nil            // t artıq string — safety check YOX
case bool:
    if t { return "truetrue", nil }
case float32, float64:           // çoxlu tip → assertion lazımdır!
    if f, ok := t.(float64); ok {
        return fmt.Sprint(f * 2), nil
    }
    return fmt.Sprint(t.(float32) * 2), nil
case int, int8, int16, int32, int64:   // hər biri ayrı case də olar
    return fmt.Sprint(t * 2), nil
default:
    return "", errors.New("unsupported type passed")
}
```
- `.(type)` YALNIZ type switch-də
- fallthrough QADAĞANDIR type switch-də
- Çoxlu-tip case-də t interfeys olaraq qalır → assertion şərt

## Əsas terminlər
- Array — sabit ölçülü; ölçü = tipin hissəsi
- [...] — element sayından ölçü çıxarma
- Açar ilə init — `{9: 10}` spesifik indekslər
- Slice — gizli massivə view
- Gizli massiv / pointer / len / cap — slice daxili quruluşu
- make(len, cap) — əvvəlcədən tutum
- append — böyümə (cap dolanda yeni massiv)
- copy — kopyalama (uzunluğu dəyişməz!)
- s[:0:0] — üç-indeksli range (ən effektiv kopya)
- Map — açar-dəyər hash; açar comparable olmalı
- Ok idiomu — `v, ok := m[k]`
- delete — map-dən silmə
- Custom type — `type id string`
- Struct / embedding / promotion
- Anonim struct — adsız bir dəfəlik tip
- Type conversion — `<type>(v)` (lossy ola bilər)
- interface{} / any — hər tipi qəbul edən boş interfeys
- Type assertion — `v.(T)` / `v.(T)` comma-ok
- Type switch — `switch v.(type)`

## Praktik nəticə
Sıralı kolleksiya → slice (array yalnız dəqiq ölçü tələbində); açarlı
axtarış → map (yalnız make/literal ilə — var+write = panic!). Slice kopyası:
append([]int{}, s...) və ya s[:0:0] — birbaşa təyinat/s[:] LINKDIR!
Bağın qırılması append-in cap həddindən asılıdır — əmin olmaq istəyirsənsə
həmişə açıq kopya. Struct modelləşdirmə: açar-dəyər init; embedding
promotion verir amma init-də tip adı tələb edir. Tip keçidlərində
itkiyatları bil; hər tipi qəbul etmək üçün any + type switch (çoxlu-tip
case-də assertion unutma).

## Mənbə
Pages: 134-191 (PDF 134-191)
