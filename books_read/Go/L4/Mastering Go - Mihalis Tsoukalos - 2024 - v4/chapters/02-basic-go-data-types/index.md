# Chapter 2 — Basic Go Data Types (Go-nun Əsas Data Tipləri)

## Bu chapter nədən bəhs edir?

error data tipi, numeric tiplər (int/uint/float/complex), overflow, string/rune/byte,
strings və unicode paketləri, time (parse, time zone), constants və iota, arrays vs slices
(length/capacity, slice header, hissə seçimi, silmə, copy, sort), pointer-lər, unsafe paketi,
random number/string generasiyası (math/rand + crypto/rand) və statistika tətbiqinin
z-normalizasiya ilə yenilənməsi.

## Əsas fikirlər

### 1. error Data Tipi — Xətalar Dəyər Kimi
**Nədir:** error — xəta şəraitini təmsil edən xüsusi interfeys tipi (interfeyslər Ch5-də).
Go xətaları dəyər kimi qaytarır: nil = xəta yoxdur.

**Kitabdan kod nümunəsi:**
```go
func check(a, b int) error {
    if a == 0 && b == 0 {
        return errors.New("this is a custom error message")
    }
    return nil
}

func formattedError(a, b int) error {
    if a == 0 && b == 0 {
        return fmt.Errorf("a %d and b %d. UserID: %d", a, b, os.Getuid())
    }
    return nil
}

i, err := strconv.Atoi("-123")   // (int, error) qaytarır
if err == nil {
    fmt.Println("Int value is", i)
}
```
**Sub-kod izahı:**
- `errors.New()` → sadə sabit xəta mesajı; `fmt.Errorf()` → formatlı xəta
- Xəta mesajının mətnini müqayisə etmək (`err.Error() == "..."`) BAD PRACTICE-dir — sadəcə
  nil yoxlaması + çap etmək kifayətdir
- **Qlobal xəta strategiyası:** bütün xətalar eyni səviyyədə emal olunmalı (ya qaytar, ya
  yerində həll et); kritik xətaların davranışı dokumentasiya olunmalı; cloud-native app-lərdə
  stderr-ə yönəltmək log servisinə yox

### 2. Numeric Tiplər və Overflow
**Cədvəl:** int8/int16/int32/int64; uint8-64; int/uint (platforma 32/64 bit — CPU register
ölçüsündən asılı, ən çox istifadə olunan); float32/float64; complex64/complex128.

**Vacib qaydalar:**
- int / int → int nəticə: `12 / 5 = 2` (tam bölünmə)
- float nəticə üçün açıq cast: `float64(x) / float64(k)`
- `math.MaxInt`/`math.MinInt` — sərhəd dəyərlər; MaxInt-i bir artırsan MinInt alınır
  (silindir kimi wrap-around: 9223372036854775807 + 1 → -9223372036854775808)
- `complex(5, 7)` və ya `12 + 1i` ilə kompleks ədədlər

### 3. String, Rune və Byte
**Nədir:** string = bayt kolleksiyası; rune = int32 — tək Unicode code point (məs. '€' =
8364); byte = uint8 — tək ASCII simvol.

**Kitabdan kod nümunəsi:**
```go
aString := "Hello World! €"
fmt.Println("First byte", string(aString[0]))   // string array kimi

r := '€'
fmt.Println("As an int32 value:", r)              // 8364
fmt.Printf("As a character: %c\n", r)            // €

for _, v := range aString {                      // rune üzrə iterasiya
    fmt.Printf("%c", v)
}
```
**Sub-kod izahı:**
- rune ilə int32 AYNI tip SAYILMIR — müqayisə edilmir (nominal typing)
- `int→string` fərqi: `strconv.Itoa(100)` → "100"; `string(100)` → "d" (Unicode code point!)
- Byte slice uzunluğu ≠ simvol sayı: "Byte slice €" = 12 simvol, 14 bayt
- Raw string literal = back quote `...` (escape emal olunmur); interpreted = "..."

### 4. unicode və strings Paketləri
**unicode.IsPrint()** — simvolun çap oluna biləcəyini yoxlayır (input filter kimi).

**strings paketinin əsas funksiyaları (kitabdan):**
```go
s.ToUpper("Hello THERE"); s.ToLower("...")
s.EqualFold("Mihalis", "MIHAlis")   // case-insensitive müqayisə
s.Index("Mihalis", "ha")            // tapılmayan yerde -1
s.Count("Mihalis", "i")
s.Repeat("ab", 5)                   // "ababababab"
s.TrimSpace(" \ttext \n")           // hər iki tərəfdən whitespace
s.TrimLeft/TrimRight(s, "\n\t ")
s.HasPrefix("Mihalis", "Mi"); s.HasSuffix("Mihalis", "is")
s.Fields("This is a string!")       // whitespace-lə böl → []string
s.Split("abcd efg", "")             // "" = simvol-simvol
s.SplitAfter("123++432++", "++")    // separator saxlanılır
s.Replace("abcd efg", "", "_", -1)  // -1 = limitsiz
s.TrimFunc("123 abc", func(c rune) bool { return !unicode.IsLetter(c) })
```

### 5. time Paketi — Parse və Time Zone
**Nədir:** time.Time — nanosaniyə dəqiqliyi ilə an; hər dəyər bir Location (time zone) ilə
assosiasiya olunur.

**Parse qaydası (vacib!):** Go tarix şablonunu **mənalı referans vaxtı** ilə göstərir:
`01/02 03:04:05PM '06 -0700` → ay/gün, saat/dəqiqə/saniyə, il, timezone.

**Kitabdan kod nümunəsi:**
```go
// "30 January 2023" parse:
t, err := time.Parse("02 January 2006", "30 January 2023")
// "15 August 2023 10:00" parse:
t, _ = time.Parse("02 January 2006 15:04", "15 August 2023 10:00")

// Time zone çevirmə:
loc, _ = time.LoadLocation("America/New_York")
fmt.Printf("New York Time: %s\n", now.In(loc))

// UNIX epoch:
now := time.Now().Unix()   // 1970-01-01-dən saniyələr
```
**Şablon elementləri:** `03` 12-saat, `15` 24-saat, `04` dəqiqə, `05` saniyə, `Mon/Monday`
gün, `02` ayın günü, `2006/06` il, `Jan/January` ay, `MST` timezone.
**Xəta nümunəsi:** "25:00" parse → `hour out of range`.
`time.Since(t)` → time.Duration (int64 underlying, amma implicit çevirmə yoxdur).

### 6. Constants və iota
**Nədir:** const — compile-time dəyər; Boolean/string/numeric tiplərində ola bilər.
iota — artan ardıcıllıq generatoru; const blokunda hər sətirdə avtomatik artır.

**Kitabdan kod nümunəsi:**
```go
const (
    Zero Digit = iota     // 0
    One                   // 1
    Two                   // 2
    Three                 // 3
)

const (
    p2_0 Power2 = 1 << iota   // iota=0 → 1
    _                          // iota=1 SKIP
    p2_2                       // iota=2 → 1<<2 = 4
    _
    p2_4                       // 16
    _
    p2_6                       // 64
)
```
**Sub-kod izahı:**
- `_` — istənməyən iota dəyərini keçir
- `type Digit int` — yeni adlı tip (nominal tip); `1 << iota` — iota ifadələrdə istifadə
  oluna bilir

**Typed vs untyped constant:**
```go
const (
    typedConstant   = int16(100)   // int16 — ancaq int16 ilə işləyir
    untypedConstant = 100           // tipsiz — hər numeric kontekstdə işləyir
)
i := int(1)
i * typedConstant      // COMPILE XƏTASI: mismatched types int and int16
i * untypedConstant     // OK
```

### 7. Arrays vs Slices
**Array xüsusiyyətləri:**
- Ölçü declaration-da mütləq (və ya `[...]` ilə sayılır); sonra DƏYİŞMİR
- Funksiyaya kopya kimi ötürülür — funksiya daxilindəki dəyişikliklər İTİR

**Slice (dinamik):**
- Başlığa (header) malik: underlying array pointer + len + cap
- `reflect.SliceHeader{Data uintptr; Len int; Cap int}`
- Funksiyaya header-in kopyası ötürülür — underlying array şərikli olur (element dəyişikliyi
  görünür; amma yeni allocation olsa bağlantı qopur)

**Kitabdan kod nümunəsi:**
```go
aSlice := []float64{}                      // boş slice: len=0, cap=0
aSlice = append(aSlice, 1234.56)           // append NƏTİCƏSİ TƏYİN EDİLMƏLİDİR!

t := make([]int, 4)                        // len=4, cap=4, hamısı 0
aSlice = append(aSlice, []int{-1, -2}...)  // ... slice-i aça bilir

twoD := [][]int{{1, 2, 3}, {4, 5, 6}}      // 2D
make2D := make([][]int, 2)                  // 2D make ilə
```

### 8. Length və Capacity
**Nədir:** len = cari element sayı; cap = yeni allocation tələb etmədən böyüyə biləcəyi
maksimum. `make([]T, len, cap)`.

**Artım qaydası:** len cap-ı keçmək üzrə olanda cap İKİLƏNİR (4→8→16).

**Kitabdan kod nümunəsi:**
```go
a := make([]int, 4)          // L:4 C:4
aSlice := make([]int, 4, 4)
aSlice = append(aSlice, 5)   // L:5 C:8 (cap 2x artdı)
aSlice = append(aSlice, []int{-1, -2, -3, -4}...)  // L:9 C:16
```
**Qayda:** cap < len mümkün deyil — `make([]int, 3, 2)` xəta.
**Optimizasiya:** böyük slice üçün əvvəlcədən cap bilsəniz təyin edin — allocation+copy qənaəti.

### 9. Slice Hissəsinin Seçimi (Slicing)
**Sintaksis:** `s[a:b]` (b daxil deyil); üçlü forma `s[a:b:c]` — cap = c-a.

**Kitabdan kod nümunəsi:**
```go
aSlice := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
aSlice[0:5]      // ilk 5 element
aSlice[:5]        // eyni
aSlice[l-2:]      // son 2
t := aSlice[2:5:10]   // len=3, cap=10-2=8
t = aSlice[:5:6]     // len=5, cap=6-0=6
```
**Vacib:** `aSlice[0:5]` YENİ slice yaradır, amma EYNİ underlying array-yə işarə edir.

### 10. Slice-dan Element Silmə
**Go 1.21+:** `slices.Delete()` hazır var. Manual üsullar:

**Kitabdan kod nümunəsi:**
```go
// Texnika 1 — sıra QORUNUR:
aSlice = append(aSlice[:i], aSlice[i+1:]...)
// i-dən əvvəl + i-dən sonra → i-ci element düşür

// Texnika 2 — sıra DƏYİŞİR (son elementi i-yə qoy):
aSlice[i] = aSlice[len(aSlice)-1]
aSlice = aSlice[:len(aSlice)-1]
```

### 11. Slice ↔ Array Bağlantısı
**Nədir:** Slice underlying array-yə pointer saxlayır; slice üzərindəki dəyişikliklər
array-ə təsir edir — CAP DƏYİŞƏNƏ QƏDƏR.

**Kitabdan kod nümunəsi:**
```go
a := [4]string{"Zero", "One", "Two", "Three"}
S0 := a[0:1]
S0[0] = "S0"           // a[0] da dəyişir!
S12 := a[1:3]
S12[0] = "S12_0"       // a[1] dəyişir

change(S12)            // funksiya daxilindən də belə — array dəyişir

S0 = append(S0, "N1")  // ... N4-ə qədər: capacity dəyişəndə
                        // YENİ underlying array → artıq a-ya bağlı DEYİL
```
**Nəticə:** cap dəyişmədiyi müddətdə a ↔ S12 əlaqəsi qalır; S0 artıq müstəqildir.

### 12. Out-of-Bounds Tutma Hiyləsi
```go
func foo(s []int) int { return s[0]+s[1]+s[2]+s[3] }   // runtime panika riski

func bar(slice []int) int {
    a := (*[3]int)(slice)                               // slice → array pointer (Go 1.17)
    return a[0]+a[1]+a[2]+a[3]                          // COMPILE XƏTASI: index 3 out of bounds [0:3]
}
```
**Sub-kod izahı:** slice-ı sabit ölçülü array pointerinə çevirmək compile-time bound
yoxlaması AÇIR — runtime xətasını compile mərhələsinə keçirir.

### 13. copy() Funksiyası
```go
copy(destination, source)
```
**Qaydalar:**
- Destination AVTOMATİK böyümür — kopyalanan element sayı = min(len(dst), len(src))
- Destination böyükdürsə qalan elementlər DOKUNULMAZ qalır
- `copy(a1, a2)` — a2 (source) heç vaxt dəyişmir

### 14. sort Paketi
```go
sort.Ints(sInts)                                 // artan
sort.Float64s(sFloats); sort.Strings(sStrings)
sort.Sort(sort.Reverse(sort.IntSlice(sInts)))    // azalan
```
Go 1.21+ slices paketi (generics) artıq mövcuddur — Ch4-də.

### 15. Pointer-lər
**Nədir:** Pointer = dəyişənin yaddaş ünvanı. `&x` — ünvan al; `*p` — dereference (dəyər al).
Pointer arifmetikası YOXDUR (unsafe istisna).

**Nəyə lazımdır:**
- Funksiyaya by-reference ötürmə — dəyişikliklər qayıdanda itmir
- Zero value ilə "təyin edilməyib" (nil) fərqləndirməsi (struct-larda vacib)
- `Next *Node` kimi öz-özünə istinad edən data strukturları (linked list, binary tree)

**Kitabdan kod nümunəsi:**
```go
func processPointer(x *float64) {
    *x = *x * *x            // x-in hədəfi dəyişir — main-də görünür
}
func returnPointer(x float64) *float64 {
    temp := 2 * x
    return &temp           // escape → heap allocation (stack frame öləndə yaşayır)
}

fP := &f                   // f-in ünvanı
processPointer(fP)         // f dəyişir

var k *aStructure          // nil pointer — heç yerə işarə etmir
if k == nil {
    k = new(aStructure)    // yaddaş ayır, zero value ilə doldur
}
```
**Xəbərdarlıq:** nil pointer-i dereference → crash; slice-lar onsuz da underlying array
pointeri daşıyır — slice üçün pointer ötürmək lazım deyil.

**Slice → Array çevirmə:**
```go
arrayPtr := (*[3]byte)(slice)   // Go 1.17: array pointer
array := [3]int(slice2)          // Go 1.20: array KOPYASI
```

### 16. unsafe Paketi (SON ÇARƏ)
**Nədir:** Tip təhlükəsizliyini pozan əməliyyatlar — StringData, String, Slice, SliceData.
Böyük string/slice çevirmələrində kopyasız sürət qazancı.

**Kitabdan kod nümunəsi:**
```go
func byteToString(bStr []byte) string {
    if len(bStr) == 0 { return "" }
    return unsafe.String(unsafe.SliceData(bStr), len(bStr))  // KOYPASIZ
}
func stringToByte(str string) []byte {
    if str == "" { return nil }
    return unsafe.Slice(unsafe.StringData(str), len(str))    // KOYPASIZ
}
```
**Sub-kod izahı:** `unsafe.StringData(s)` → string-in underlying byte pointer; `unsafe.Slice(p, n)`
→ p-dən n elementli slice. **QADAĞA:** unsafe.StringData qaytırdığı baytları dəyişmək —
Go string-ləri immutable-dir; hər hansı dəyişiklik = tanınmayan davranış.

### 17. Random Ədədlər və Stringlər
**math/rand (pseudo-random):** seed eynidursa → eyni ardıcıllıq (test üçün faydalı,
təhlükəsizlik üçün YOX).

**Kitabdan kod nümunəsi:**
```go
func random(min, max int) int {
    return rand.Intn(max-min) + min     // [min, max)
}

// Random string (ASCII 33='!' .. 126='~', 94 çap olunan simvol):
newChar := string(startChar[0] + byte(myRand))
```

**crypto/rand (cryptographically secure):**
```go
func generateBytes(n int64) ([]byte, error) {
    b := make([]byte, n)
    _, err := rand.Read(b)          // bütün slice-i doldurur; seed LAZIM DEYİL
    return b, nil
}
// sonra: base64.URLEncoding.EncodeToString(b) → çap olunan string
```

### 18. Statistika Tətbiqi — Z-Normalizasiya
**Nədir:** Dəyərlərin müqayisə oluna bilən forma gətirilməsi: (val - mean) / stdDev.
stdDev=0 xüsusi halında data olduğu kimi qaytarılır.

**Kitabdan kod nümunəsi:**
```go
func normalize(data []float64, mean float64, stdDev float64) []float64 {
    if stdDev == 0 {
        return data
    }
    normalized := make([]float64, len(data))
    for i, val := range data {
        normalized[i] = math.Floor((val-mean)/stdDev*10000) / 10000  // 4 dəqiqəlik
    }
    return normalized
}

func randomFloat(min, max float64) float64 {
    return min + rand.Float64()*(max-min)   // [min, max)
}
```
- Etibarsız input-un HAMISI etibarsızdırsa → 10 random dəyərlə doldurulur (test data)

## Əsas terminlər
- error interface — xəta dəyəri (nil = xəta yox)
- Rune — int32, tək Unicode code point
- Zero Value — tipin default dəyəri (0, false, "", nil)
- Slice Header — pointer + len + cap üçlüyü
- Capacity (tutum) — allocation-sız böyümə həddi
- iota — const ardıcıllıq generatoru
- Typed/Untyped Constant — tipli (məhdud) / tipsiz (çevik) konstant
- Dereference (dəyərə keçid) — *p ilə pointerin hədəfini oxu
- Overflow (daşma) — tip sərhədini aşınca wrap-around
- Pseudo-random — seed-dən asılı saxta təsadüfi ardıcıllıq
- Cryptographically Secure Random — proqnozlaşdırıla bilmez generator (crypto/rand)
- Z-normalization — (val-mean)/stdDev standartlaşdırması

## Praktik nəticə

(1) error həmişə dəyər kimi — `if err != nil` Go-nun zərfi; xəta mesajlarını string müqayisəsi
ilə yoxlama. (2) İnteqer bölgə tam bölür — float istəyirsənsə cast et. (3) MaxInt+1 = MinInt —
overflow səssizdir. (4) string(100) ≠ "100" — İtoi/FormatInt istifadə et. (5) append-in
nəticəsini həmişə təyin et; slice cap-i dəyişəndə underlying array bağlantısı QOPUR. (6) cap-ı
əvvəlcədən bilmək = allocation qənaəti. (7) copy() destination-i böyütmür. (8) Pointer =
by-reference + nil fərqi + öz-istinadlı strukturlar. (9) unsafe yalnız ölçülü, riski başa
düşülən hallarda. (10) Təhlükəsizlik lazım olan yerdə math/rand YOX — crypto/rand. (11) Boş
slice = len 0; nil slice = nil pointer — nil slice həmişə boşdur, əksi həmişə doğru deyil.

## Mənbə
Pages: 47-102 (PDF 78-135)
