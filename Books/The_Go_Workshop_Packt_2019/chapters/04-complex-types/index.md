# Chapter 4 — Complex Types (Mürəkkəb Tiplər)

## Bu fəsil nədən bəhs edir?

Kolleksiyalar (array — ölçü, müqayisə, açarlı init; slice — append, range kəsimi,
make/len/cap, hidden array bağlantıları, copy üsulları; map — yaratma/oxuma/silmə,
comma-_ok), sadə custom tiplər, struct-lar (init üsulları, müqayisə, anonim struct,
embedding/promotion), tip çevirmələri (itirici), interface{} + type assertion + type switch.

## Əsas fikirlər

### 1. Array — Sabit Ölçülü Kolleksiya
**Bəyan:** `[<ölçü>]<tip>` — ölçü TİPİN hissəsidir: `[5]int` ≠ `[9]int` (FƏRLİ tiplər!).
```go
var arr [10]int                    // zero-lu, ölçü sabit
arr2 := [5]int{0}                  // literal ilə
arr3 := [...]int{0,0,0,0,0}       // ... = say avtomatik
arr3 := [10]int{1, 9: 10, 4: 5}   // AÇARLI init — sıra əhəmiyyətsiz, boşluqlar OK
```
**Müqayisə:** eyni ölçülü array-lar == ilə müqayisə OLUNAR (slice/map OLMAZ — loop lazım).
**Uzunluq funksiyası:** `len()` — runtime izləyir, hər dəfə SAYMIR — loop şərtində
istifadəsi təhlükəsizdir.

### 2. Slice — Dinamik Qat
**Nədir:** array üstündə nazik qat; ölçü Go tərəfindən idarə olunur. Real kodda DEFAULT.

```go
var s []int                        // nil slice
s = append(s, dəyər)               // NƏTİCƏNİ TƏYİN ET!
s = append(s, "a", "b", "c")      // çoxsaylı
s = append(s, digərSlice...)       // ... = unpack — slice birləşdirmə
```

**Range kəsimi — [low:high]:**
```go
s := []int{1,2,3,4,5,6,7,8,9}
s[0:1]      // [1]      — high DAXİL DEYİL
s[:5]       // ilk 5    — low default 0
s[5:]       // son 4    — high default len
s[:]        // hamısı   — array-i slice-ə çevirir (KOYA YOX!)
```

### 3. Slice Daxili Mexanizmi — 3 Gizli Sahə
Pointer (hidden array-a) + başlanğıc nöqtəsi + uzunluq. **Capacity = hidden array
ölçüsü** (`cap()`).

**Append 2 ssenari:**
1. Gizli arrayda yer VAR → dəyər yazılır, len artır
2. Array DOLU → yeni BÖYÜK array + köhnə dəyərlər köçürülür → pointer YENİYƏ döndərilir

**make ilə nəzarət:**
```go
make([]int, 10)        // len=10, cap=10
make([]int, 10, 50)    // len=10, cap=50 — artıq tutum = append-də köçürmə YOX
```

### 4. Slice Bağlantı Ssenariləri (Ən Vacib Təhlil!)
**Kitabdan 6 ssenari:**
```go
// 1) LINKED — sadə kopya + range kəsimi EYNİ arraya işarə edir:
s1 := []int{1,2,3,4,5}
s2 := s1           // kopya, amma pointer EYNİ
s3 := s1[:]        // eyni hidden array
s1[3] = 99         // s2[3]=99, s3[3]=99 — hamısı DƏYİŞİR

// 2) NO LINK — append cap-i AŞIRSA yeni array → əlaqə QOPUR:
s2 := s1
s1 = append(s1, 6)  // cap dolu → yeni array; s2 köhnədə qalır
s1[3] = 99         // s2[3] = 4 — TƏSİR ETMİR

// 3) CAP LINK — artıq cap VARSA append bağlantını SAXLAYIR:
s1 := make([]int, 5, 10)
s2 := s1
s1 = append(s1, 6)  // hələ cap-dadır → eyni array
s1[3] = 99         // s2[3] = 99 — BAĞLI

// 4) CAP NO LINK — cap-i aşan çox append → qopma (2-dəki kimi)

// 5) COPY NO LINK — built-in copy:
s2 := make([]int, len(s1))
copied := copy(s2, s1)   // nə qaytarır = kopyalanan say; s2 ölçüsü DƏYİŞMİR

// 6) APPEND NO LINK — ən məşhur kopyalama üsulu:
s2 := append([]int{}, s1...)           // yeni boş slice + unpack
// ƏN effektiv variant (0-cap range):
s2 := append(s1[:0:0], s1...)           // 3-lü kəsim [low:high:cap]!
```
**Dərs:** append cap həddini aşanda gizli array ƏVƏZLƏNİR — kopyalara təsir
dayanır. Bu, "səhv çalışan" slice-lərin ən məşhur mənbəyidir.

### 5. Map — Açar-Dəyər Hash
**Fərq:** açar MƏNALI data-dır (ID → record), sadecə sayğac DEYİL. Constant-time axtarış.

```go
// Yaratma (var-dan QAÇIN — yazış runtime panic!):
users := map[string]string{"305": "Sue", "204": "Bob"}
m := make(map[string]int, 10)      // cap TÖVSİYƏDİR, cap() İŞLƏMİR

users["073"] = "Tracy"             // set
user, exists := users[id]          // COMMA-OK — açar mövcudluğu!
delete(users, id)                  // tam silmə; olmayan açar = heç nə olmur
```
**Açar tipləri:** yalnız comparable tiplər (slice/map AÇAR OLAMAZ). Sıra RANDOM-dur —
sıralı çıxış üçün slice köməkçisi lazım.
**Olmayan açardan oxu:** zero value qaytarır — zero-value məntiqi mümkünsə comma-ok
əvəzinə istifadə et.

### 6. Sadə Custom Tiplər
```go
type id string                     // string əsaslı yeni tip
var id2 id = "1234-5678"
string(id2) == "1234-5678"          // ƏSAS tipə ÇEVİRMƏ ilə müqayisə
// id2 == "1234-5678"               // XƏTA — id və string UYĞUN DEYİL
```
Eyni davranış (zero value, müqayisə), amma əsas tipdən AYRI — domain modelinin bir vahidi.

### 7. Struct — Sahələr Qrupu
```go
type user struct {
    name    string
    age     int
    balance float64
    member  bool
}

// 3 init üsulu:
u1 := user{name: "Tracy", age: 51, ...}   // AÇARLI — ən məşhur; sıra fərq etmir,
                                           // buraxılan sahələr zero
u3 := user{"Bob", 25, 0, false}           // SIRALI — HAMISI olmalı, sıra dəqiq
var u4 user                                 // zero-lu; sonra .field ilə doldur
u4.name = "Sue"

// Anonim struct (bir dəfəlik):
point1 := struct{ x, y int }{10, 10}
```
**Müqayisə:** bütün sahələr comparable-dırsa struct == müqayisə olunar; anonim struct
eyni strukturlu adlı tip ilə müqayisə OLUNAR (Go-nun çevikliyi).

### 8. Embedding — Kompozisiya
**Nədir:** inheritance YOXDUR; struct-i struct-in İÇİNƏ adsız sahə kimi qoş — sahələr
PROMOTE olunur (birbaşa çatılır).

**Kitabdan kod nümunəsi:**
```go
type name string
type location struct{ x, y int }
type size struct{ width, height int }

type dot struct {
    name          // adsız sahələr = embedding
    location
    size
}

var d dot
d.x = 5                  // PROMOTED — birbaşa
d.width = 10
d.location.x = 13         // TİP ADI ilə də mümkün
d.name = "A"             // sadə tip embedding — tip adı = sahə adı

// İNİSİALİZASİYADA promotion İŞLƏMİR — tip adı ilə:
dot3 := dot{
    name: "B",
    location: location{x: 13, y: 27},
    size:     size{width: 5, height: 7},
}
```
**Qaydalar:** ad toqquşması olsa promote BAŞ VERMİR (tip yolu ilə çatılır); pointer
embedding-də ad `*T` → `T` olur (field pointer olaraq qalır). Real kodda az rast
gəlinir — adlı sahələr daha məşhurdur.

### 9. Tip Çevirmələri — İtirici Ola Bilər
```go
int64(i8)      // kiçik→böyük int — HƏMİŞƏ TƏHLÜKƏSİZ
int8(128)      // overflow! → -128
float64(i8)    // dəyişməz
int(3.14)      // 3 — kəsr TRUNCATE
var i8 int8 = math.MaxInt8   // IMPLICIT (constant int → int8)
```
Data hədləri keçmirsə itirici çevirmə etibarlıdır — real kodda davamlı olur.

### 10. interface{} və Type Assertion
**Nədir:** interface = bir tipin SAHİB OLMALI OLDUĞU funksiyalar müqaviləsi.
`interface{}` = 0 funksiya → HƏR TİP uyğun gəlir (fmt.Print-in sirri: `func Print(a ...interface{})`).

**Type assertion — `v.(T)`:**
```go
func doubler(v interface{}) (string, error) {
    if i, ok := v.(int); ok {           // COMMA-OK təhlükəsiz form
        return fmt.Sprint(i * 2), nil
    }
    if s, ok := v.(string); ok {
        return s + s, nil
    }
    return "", errors.New("unsupported type passed")
}
// i := v.(int)                          // tək formal — fail olsa PANİK!
```
**Sub-kod izahı:** interface{} dəyəri bağlayır, amma compiler tip yoxlamasını dayandırır;
assertion runtime-da bu yoxlamanı İCRA EDİR — məsuliyyət sənə keçir.

### 11. Type Switch — Çoxtipli Assertion
```go
switch t := v.(type) {
case string:                     // t = string — birbaşa istifadə
    return t + t, nil
case bool:
    if t { return "truetrue", nil }
    return "falsefalse", nil
case float32, float64:           // ÇOXLU tip — t interface{} qalır!
    if f, ok := t.(float64); ok {
        return fmt.Sprint(f * 2), nil
    }
    return fmt.Sprint(t.(float32) * 2), nil
case int: return fmt.Sprint(t * 2), nil    // ... bütün int/uint tipləri ayrı-ayrı
default:
    return "", errors.New("unsupported type passed")
}
```
**Qaydalar:** `.(type)` YALNIZ switch-də; `fallthrough` İŞLƏMİR; çoxlu-tip case-də
əlavə assertion lazımdır.

## Əsas terminlələr
- Array — `[N]T` ölçü = tipin hissəsi; müqayisə mümkün
- `...` init — element sayından ölçü çıxarma
- Slice — hidden array üzərində pointer+start+len qatı
- Capacity (cap) — hidden array ölçüsü; append həddi
- make([]T, len, cap) — başlanğıc tutum nəzarəti
- Unpack (`...`) — slice-i append arqumentinə açma
- Range Notation [low:high] — high daxil deyil; subslice
- `[a:b:c]` 3-lü — low:high:capacity kəsimi
- copy() — ölçüsü dəyişməyən dəyər kopyası
- Comma-Ok (`v, ok := m[k]`) — açar mövcudluğu yoxlaması
- delete() — map elementinin tam silinməsi
- Custom Type — `type id string`; əsas tipdən ayrı
- Struct — adlı sahələr toplusu; == mümkün (comparable sahələrlə)
- Anonim Struct — adsız birdefəlik strukt
- Embedding/Promotion — adsız sahə; sahələrin yuxarı qalxması
- interface{} — 0-metodlu müqavilə; hər tipi qəbul edir
- Type Assertion — `v.(T)` runtime tip açarı; panic riski
- Type Switch — `switch v.(type)` — çoxtipli assertion

## Praktik nətidə

(1) Array ölçüsü TİPDİR — [5]int və [9]int qarışmaz; funksiyalar slice qəbul etsin.
(2) append nəticəsini HƏMİŞƏ təyin et; `s = append(s, ...)`. (3) Slice kopyası
BAĞLIDIR — append cap həddini aşırsa qopur, yoxsa dəyişikliklər görünür; müstəqil
kopya üçün `append(s1[:0:0], s1...)`. (4) Böyüklüyü bilinirsə make + cap —
reallocation qənaəti. (5) Map: var bəyanından QAÇIN — nil map yazış panic; comma-ok
mövcudluğu, delete tam silir. (6) Map sırası random — sıralı üçün slice saxla.
(7) Custom tip domainə ad verir; əsas tipə çevir. (8) Embedding promotion verir, amma
İNİT-də tip adı tələb edir; real kodda adlı sahə üstünlük təşkil edir. (9) Numeric
çevirmə itiricidir: int8(128)=-128, int(3.14)=3 — hədləri bil. (10) interface{} +
comma-ok assertion = təhlükəsiz genericlik; type switch çoxlu tipləri səliqələyir.
(11) fmt-in mənbəyinə bax — interface{...} variadic patterni öyrənmə qaynağıdır.

## Mənbə
Pages: 109-167 (PDF 142-201)
