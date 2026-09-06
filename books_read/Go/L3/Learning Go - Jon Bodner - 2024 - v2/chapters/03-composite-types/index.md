# Chapter 3 — Composite Types (Kompozit Tiplər)

## Bu chapter nədən bəhs edir?

Go-nun kompozit tipləri: array-lər (və niyə birbaşa az istifadə olunur), slice-lər
(len/cap/append/make, slicing, memory paylaşımı), string-byte-rune UTF-8 münasibəti,
map (hash map), comma-ok idiomu, map-i set kimi istifadə və struct-lar.

## Əsas fikirlər

### 1. Array — Birbaşa İstifadə Üçün Çox Sərt
**Nədir:** Müəyyən sayda, eyni tipli elementin sabit uzunluqlu konteyneri.

**Necə işləyir:** `[3]int` — ölçü tipin bir hissəsidir. `[3]int` ilə `[4]int` **fərqli
tiplərdir** — bir-birinə təyin/müqayisə/convert oluna bilmir. Ölçü dəyişənlə verilə
bilmir (tiplər compile vaxtı çözülməlidir).

**Kitabdan kod nümunəsi:**
```go
var x [3]int                      // [0 0 0]
var x = [3]int{10, 20, 30}
var x = [12]int{1, 5: 4, 6, 10: 100, 15}  // sparse: index-dən asılı dəyərlər
var x = [...]int{10, 20, 30}      // ölçü avtomatik
var x [2][3]int                   // "multidimensional" — əslində array-in array-i
fmt.Println(len(x))               // ölçü
```

**Sub-kod izahı:**
- `[12]int{1, 5: 4, ...}` → 5-ci indeksə 4 qoyulur, boşluqlar zero value qalır
- `[...]` → literal-dakı element sayına görə ölçü təyin olunur
- `==`/`!=` — array-lər müqayisə oluna bilir (eyni tip olduqda)

**Niyə az istifadə olunur:** Ölçü tipin hissəsi olduğu üçün hər ölçü üçün ayrı funksiya
yazmaq lazım gəlir. İstisna: ölçü alqoritmlə tanınan hallar (kriptoqrafik checksum-lar).
**Array-lərin əsas səbəbi slice-lər üçün backing store olmaqdır.**

### 2. Slice — Go-nun Ən Faydalı Konteyneri
**Nədir:** Dinamik uzunluqlu ardıcıllıq; ölçü tipin hissəsi DEYİL.

**Necə işləyir:** `[]int` — ölçüsüz bəyan. Zero value — `nil`. Nil slice boşdur; `len`
0 qaytarır. **Slice-lar müqayisə olunmur** (`==` compile xətası) — yalnız `== nil`
mümkündür (müqayisə üçün `reflect.DeepEqual` — test üçün).

**Kitabdan kod nümunəsi:**
```go
var x = []int{10, 20, 30}   // literal
var x []int                  // nil slice
x = append(x, 10)            // nil-ə append edilə bilər
x = append(x, 5, 6, 7)      // çoxlu dəyər
y := []int{20, 30, 40}
x = append(x, y...)          // slice-in slice-a açılması
```

**Sub-kod izahı:**
- `append` **dəyər qaytarır və təyini mütləqdir** — Go call-by-value olduğundan funksiyaya
  slice-in kopyası gedir; dəyişmiş kopya geri qaytarılır
- `y...` → `...` operatoru slice-ı ayrı-ayrı arqumentlərə açır

**len vs cap:** hər slice-in capacity-si (ayrılmış ardıcıl yaddaş) var. Uzunluq
capacity-yə çatdıqda `append` daha böyük yeni slice ayırır və köhnə dəyərləri kopyalayır.
Böyümə qaydası (Go 1.14+): capacity < 1024 → 2x; sonra ≥ 25%. Nümunə ardıcıllığı:
`[] → [10] → [10 20] → [10 20 30](cap 4) → [10 20 30 40] → [10 20 30 40 50](cap 8)`.

**make:** Ölçünü əvvəlcədən bilmək — ən efficient yol.
```go
x := make([]int, 5)      // len 5, cap 5 — [0 0 0 0 0]
x = append(x, 10)        // DIQQƏT: [0 0 0 0 0 10] — append həmişə sona əlavə edir!
x := make([]int, 0, 10)  // len 0, cap 10 — append üçün ən təhlükəsiz forma
```

**Hansı bəyanat forması (qayda):**
1. Slice boş qala bilərsə → `var data []int` (nil)
2. Başlanğıc dəyərlər varsa/dəyişməyəcəksə → literal `[]int{...}`
3. Ölçünü bilirsənsə → `make`; buffer üçün nonzero len, dəqiq bilirsən nonzero len,
   digər hallarda **zero len + nonzero cap + append** (ən az bug riski; müəllifin
   tövsiyəsi). `cap < len` təyin etmək — compile xətası (konstantdırsa) / panic (variabledırsa).

### 3. Slicing — Slice-dan Slice (Təhlükə!)
**Nədir:** `x[a:b]` ifadəsi — alt ardıcıllıq.

**Necə işləyir:** Nüsxə YOX — **hər iki slice eyni yaddaşı paylaşır**. Birindəki element
dəyişikliyi digərində də görünür. Subslice-ın capacity-si = valideynin capacity-si minus
offset — yəni istifadə edilməmiş capacity də paylaşılır: `y := x[:2]; y = append(y, 30)`
→ `x`-in 3-cü elementi 30-la əvəzlənir!

**Kitabdan kod nümunəsi:**
```go
x := []int{1, 2, 3, 4}
y := x[:2]
z := x[1:]
x[1] = 20; y[0] = 10; z[1] = 30
// x: [10 20 30 4], y: [10 20], z: [20 30 4] — hamısı təsirlənir
```

**Full slice expression (müdafiə):** `x[start:end:capLimit]` — üçüncü hissə subslice-ın
maksimum capacity-sini məhdudlaşdırır, append-in valideynin yaddaşını tapdalmasını
əngəlləyir:
```go
y := x[:2:2]   // cap(y) = 2 → append yeni slice yaradır, x toxunulmaz
z := x[2:4:4]
```

**Qayda:** Slice-lanmış/slice-lanacaq slice-ləri dəyişməkdən çəkin; ya da üçhissəli
ifadə işlət. Array-dən slice almaq da eyni paylaşım xüsusiyyətinə malikdir.

### 4. copy — Müstəqil Nüsxə
**Kitabdan kod nümunəsi:**
```go
x := []int{1, 2, 3, 4}
y := make([]int, 4)
num := copy(y, x)        // 4 — kopyalanan element sayı
num = copy(x[:3], x[1:])  // üst-üstə düşən sahələr daxilində: [2 3 4 4] 3
copy(y, d[:])            // array ↔ slice kopyasiya mümkün
```

**Sub-kod izahı:** `copy(dst, src)` — daha kiçik olanın uzunluğu qədər kopyalayır;
qaytarılan dəyər kopyalanan saydır (istifadəsizsə təyin etmək olmaz).

### 5. String-lər, Rune-lər, Baytlar və UTF-8
**Nədir:** Go string-i bayt ardıcıllığıdır — rune YOX. UTF-8 gözlənilir.

**Necə işləyir:**
- `s[i]` → i-ci **bayt** (rune deyil); `len(s)` → **bayt** sayı (`"Hello ☀"` → 10, 7 deyil)
- UTF-8-də code point 1-4 bayt ola bilər → çoxbaytlı simvolların ortasından slicing
  kor simvol yaradır. Slice/index yalnız 1-baytlıq məzmun (ASCII) olduqda təhlükəsizdir.
- Conversion-lar: `[]byte(s)`, `[]rune(s)` (az işlənir), `string(a)` (rune/byte-dan),
  `string(x)` int-dən — **TƏLƏ**: `string(65)` = `"A"` (Unicode code point!), `"65"` yox.
  Go 1.15+ `go vet` bunu yalnız rune/byte xaricində bloklayır.
- Substring/simvol əməliyyatları üçün `strings` və `unicode/utf8` paketləri.

**UTF-8 arxitekturası:** <128 dəyərlər üçün 1 bayt, maksimum 4 bayt; byte-order problemi
yoxdur; bayt axınında simvolun başı/ortası tanınır; təkcə random access mümkün deyil
(başdan saymaq lazım). UTF-8-i 1992-də Ken Thompson və Rob Pike (Go-nun yaradıcıları)
ixtira edib.

**Kitabdan kod nümunəsi:**
```go
var s string = "Hello, ☀"
var bs []byte = []byte(s)  // [72 101 108 108 111 44 32 240 159 140 158]
var rs []rune = []rune(s)  // [72 101 108 108 111 44 32 127774]
```

### 6. Map — Hash Map
**Nədir:** `map[keyType]valueType` — açar-dəyər assosiasiyası; daxili implementasiya hash
map-dır (açar → hash → bucket; toqquşma (collision) halında bucket daxilində axtarış).

**Necə işləyir:**
- Zero value — `nil`; **nil map-dən oxuma zero value qaytarır, yazmaq PANIC edir**.
- `map[string]int{}` — empty literal: nil DEYİL, oxu-yazı açıq.
- `make(map[int][]string, 10)` — başlanğıc ölçü ilə; uzunluq 0-dan başlayır, böyüyə bilər.
- Açar — yalnız müqayisə edilə bilən (comparable) tip ola bilər (slice/map açar OLMAZ).
- Map == ilə müqayisə OLUNMUR (yalnız `== nil`); `len` — cüt sayı.

**Kitabdan kod nümunəsi:**
```go
totalWins := map[string]int{}
totalWins["Orcas"] = 1
totalWins["Lions"] = 2
totalWins["Kittens"]++      // mövcud olmayan açara ++ — zero value-dan başlayır
fmt.Println(totalWins["Kittens"])  // 1
```

**Oxuma/Yazma:** yazma — `m[k] = v` (`:=` YOX); oxu — `v := m[k]` — açar yoxdursa
value tipinin zero value-su qaytarır.

### 7. Comma-ok İdiomu
**Nədir:** Açarın map-də olub-olmamasını dəyərdən ayıran oxuma forması.

**Necə işləyir:** `v, ok := m[k]` — ok=true açar var, false yox. "Açar var amma dəyəri
zero value" ilə "açar yoxdur" fərqini ayırır.

**Kitabdan kod nümunəsi:**
```go
m := map[string]int{"hello": 5, "world": 0}
v, ok := m["world"]    // 0 true
v, ok = m["goodbye"]   // 0 false
```

Bu idiom daha sonra channel oxumasında (Ch10) və type assertion-da (Ch7) da çıxacaq.

### 8. delete və Map-i Set kimi istifadə
**delete:** `delete(m, "hello")` — açar yoxdursa və ya map nil-dirsə heç nə olmur; dəyər
qaytarmır.

**Set simulyasiyası:** Go-da built-in set YOXDUR → `map[int]bool{}`; `intSet[v] = true`;
üzvlük: `if intSet[100] {...}`. Alternativ: `map[int]struct{}{}` — zero bayt yer (bool
1 bayt), amma kod qəlizləşir və comma-ok tələb edir — böyük set-lər istisna olmaqla
fərq əhəmiyyətsizdir.

### 9. Struct — Əlaqəli Datanın Qruplaşdırılması
**Nədir:** Adlandırılmış sahələri olan tip; class/inheritance YOXDUR (OOP xüsusiyyətləri
Chapter 7-də başqa formada).

**Necə işləyir:**
```go
type person struct {
    name string
    age  int
    pet  string
}
var fred person              // bütün sahələr zero value
bob := person{}              // eynilə — map-dən fərqli olaraq fərq YOXDUR
julia := person{"Julia", 40, "cat"}           // posisional — HƏMİ sahə olmalı, sıra məcburi
beth := person{age: 30, name: "Beth"}          // adlı — istənilən sıra, çatışmayan zero
```
- İki literal stili QARIŞDIRMAQ olmaz; gələcəkdə sahə əlavə olunacaqsa adlı stil təhlükəsizdir
  (posisional compile xətasına keçir — xeyirli).
- Sahə oxu/yazı: `bob.name = "Bob"`.

**Anonymous struct:** adsiz tip — dəyişən bir instansiya ilə kifayətlənirdik: JSON
marshal/unmarshal və table-driven test-lərdə (Ch13) istifadə olunur:
```go
pet := struct {
    name string
    kind string
}{name: "Fido", kind: "dog"}
```

**Müqayisə:** Bütün sahələr comparable idisə struct müqayisə olunur (slice/map/funksiya/channel
sahəsi varsa YOX). Custom equality metodu YOXDUR.

**Conversion (müqayisədən fərqli):** sahə adları, sırası və tipləri **tam eyni** olan
structlar arasında type conversion MÜMKÜNDÜR (hansısa sahə fərqlidirsə — compile xətası).
Anonymous struct bonusu: eyni sahəli adlı + anonim struct-lar arasında həm `=`, həm `==`
convertsuz işləyir.

## Əsas terminlər
- Slice backing store (slice-ın arxa anbarı) — array-in slice üçün yaddaş təmin etməsi
- Capacity (tutum) — ayrılmış yaddaş yerlərinin sayı
- Full slice expression (tam slice ifadəsi) — `x[a:b:c]` üçhissəli form
- Comma-ok idiom (vergül-ok idiomu) — dəyər + mövcudluq göstəricisi oxuma
- Hash map (xəş xəritəsi) — açarı hash edərək bucket-də saxlayan struktur
- Anonymous struct (anonim strukt) — adsız struct tipi
- Marshaling/Unmarshaling (seriyalaşdırma) — struct ↔ xarici data (JSON və s.)

## Praktik nətidə

Ən çox bug yaradan sahə: slice paylaşımı. Qaydalar: (1) slice-dan subslice aldıqda
dəyişikliklərdən çəkin ya üçhissəli `x[:2:2]` işlət; (2) müstəqil nüsxə üçün `copy`;
(3) string-də bayt/rune sayı fərqlidir — `len(s)` bayt sayıdır, simvol əməliyyatı üçün
`[]rune(s)` ya `unicode/utf8`; (4) `make`-də nonzero len + append kombinasiyası
sürpriz zero-lar yaradır — `make(T, 0, n)` + append ən təhlükəsizdir; (5) nil map-ə yazma
panic-dir — map-i əvvəlcədən `make` və ya literal ilə yarat.

## Mənbə
Pages: 65-100
