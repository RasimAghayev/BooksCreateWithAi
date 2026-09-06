# Chapter 2 — Value Parts and Value Sizes (Dəyər hissələri və ölçüləri)

## Bu chapter nədən bəhs edir?
Go-da dəyərlərin yaddaşda necə təşkil olunduğuna — direct/indirect hissələr, tip ölçüləri, memory alignment (yaddaş düzləndirilməsi), struct padding (strukturların doldurulması) və dəyər kopyalama xərclərinə. Bu, kitabın bütün optimizasiya mövzularının təməlidir: hansı əməliyyatların "pulsuz", hanslarının "bahalı" olduğunu anlamaq üçün əvvəlcə dəyərin yaddaşdakı görünüşünü bilmək lazımdır.

## Əsas fikirlər

### 1. Direct part və indirect part (birbaşa və dolayı hissələr)
**Nədir:** Go 101 seriyasının terminologiyası — dəyərin bir hissəsi öz yaddaş blokunda birbaşa durur (direct part), digər hissələr isə ona istinad (reference) vasitəsilə bağlanır (indirect parts).

**Necə işləyir:** Bir hissəli tiplər — boolean, numerik, pointer, unsafe pointer, struct, array (elementlər eyni blokda). İndirect hissəli ola bilən tiplər — slice, map, channel, function, interface, string.

**Nəyə lazımdır:** Assign zamanı yalnız **direct part** köçürülür — indirect hissələr paylaşılır. Bu, slice-ın niyə "yarı-kopya" olduğunu, struct-un niyə "tam kopya" olduğunu izah edir.

**Kitabdan kod nümunəsi:**
```go
// Bir hissəli: kopya = tam müstəqil nüsxə
b := a // struct/array: BÜTÜN sahələr/elementlər kopyalanır

// İki hissəli: kopya = yalnız başlıq (direct part)
s2 := s1 // slice: pointer+len+cap kopyalanır, elementlər ŞƏRƏKLİDİR
```

**Sub-kod izahı:**
- `s2 := s1` → yalnız 3 word (ptr, len, cap) kopyalanır; backing array (dəstəkləyici massiv) ümumi qalır — birində element dəyişsən, digəri dəyişir
- Map/channel/function daxildə pointer kimi təmsil olunur — kopyaları ucuzdur (1 word)

### 2. Value/type sizes (dəyər və tip ölçüləri)
**Nədir:** Dəyərin ölçüsü = yalnız direct partın ölçüsü. Indirect hissələr ölçüyə daxil edilmir, çünki onlar bir çox dəyər arasında paylaşılır.

**Necə işləyir (64-bit arxitekturada, 1 word = 8 bayt):**

| Tip | Ölçü |
|-----|------|
| int, uint, uintptr, pointer | 1 word |
| string | 2 words (ptr + len) |
| slice | 3 words (ptr + len + cap) |
| map, channel, function | 1 word (daxildə pointer) |
| interface | 2 words (dinamik tip + dəyər) |
| array | element ölçüsü × uzunluq |
| struct | sahələrin cəmi + padding |

**Nəyə lazımdır:** Map-də key kimi string (2 word) vs kiçik array seçəndə, yaddaş budçası hesablayanda.

### 3. Memory alignments (yaddaş düzləndirilməsi)
**Nədir:** CPU instruksiyalarının səmərəli istifadəsi üçün dəyərin başlanğıc ünvanı N-in qatına düzləndirilir; N = tipin alignment guarantee-si.

**Necə işləyir:** `unsafe.Alignof(t)` ilə ölçülür. bool/int8 → 1, int16 → 2, int32/float32 → 4, digərləri (64-bit) → 8. Struct-un alignment-i ən böyük sahə alignment-i qədərdir; array-in alignment-i element tipinə bərabərdir.

**Nəyə lazımdır:** Tipin ölçüsü həmişə onun alignment-inin qatidadır — bu qayda padding-i doğur.

### 4. Struct padding — sahə sırası ölçüyə təsir edir
**Nədir:** Alignment tələblərini ödəmək üçün compiler sahələrin arasına boş baytlar (padding) əlavə edir; bu baytlar struct ölçüsünə daxildir.

**Necə işləyir:** Eyni 3 sahə, iki fərqli sıra → iki fərqli ölçü (64-bit):

```go
type T1 struct {  // 24 bayt
    a int8       // +7 bayt padding (b üçün 8-ə düzlənmə)
    b int64
    c int16      // +6 bayt padding (ümumi ölçü 24 = 8×3 olsun deyə)
}

type T2 struct {  // 16 bayt
    a int8       // +1 bayt padding
    c int16
    b int64      // +4 bayt padding
}

func main() {
    println(unsafe.Sizeof(T1{})) // 24
    println(unsafe.Sizeof(T2{})) // 16
}
```

**Sub-kod izahı:**
- `T1` → böyük sahələr (int64) arasına kiçik sahə (int8) səpələnib → 8 bayt israf
- `T2` → sahələr ölçüyə görə qruplaşdırılıb → 8 bayt qazanc (33%)
- Qayda: sahələri ölçüyə görə azalan sırayla düzmək padding-i minimuma endirir

**Nəyə lazımdır:** Minliklərlə elementli massivlərdə struct ölçüsü 24→16 olsa, 25% yaddaş qənaəti deməkdir.

**Kitabın tövsiyəsi:** Oxunaqlıq birincildir — sahələri məntiqi qrupla saxla; yalnız real ehtiyac olduqda (milyonlarla instans) ölçüyə görə sırala.

### 5. Kiçik ölçülü tiplər xüsusi optimizasiya olunur
**Nədir:** Compiler ~4 word-ə qədər struct/array kopyalamanı xüsusi (sürətli) instruksiyalarla icra edir.

**Necə işləyir:** Kitabın benchmark sübutu (Go v1.19, element = 1 word):
- 9-element array kopyası: 3974 ns/op; **10-element: 8896 ns/op** (2.2× sıçrayış!)
- 9-field struct: 2970 ns/op; 10-field: 8471 ns/op
- `Add4` (4×float32 struct): 2.65 ns; `Add5` (5 sahə): 19.15 ns — 7× fərq!

```go
type T4 struct{ a, b, c, d float32 }      // sürətli yol
type T5 struct{ a, b, c, d, e float32 }   // "büyük" sayılır

//go:noinline
func Add4(x, y T4) (z T4) {  // 2.65 ns/op
    z.a = x.a + y.a; z.b = x.b + y.b
    z.c = x.c + y.c; z.d = x.d + y.d
    return
}
```

**Nəyə lazımdır:** Hot path-də struct-u ≤4 word saxlamaq (məs. vektor hesabı 3D-də qalsın, 5D-yə çıxmasın).

### 6. Dəyər kopyası yaradan əməliyyatların siyahısı
Assign-dən əlavə bunlar da kopyalayır:
- value boxing (interfeysə çevirmə)
- funksiya parametr/rezultat ötürməsi
- channel-a send/receive
- map-ə entry qoyma
- slice-a append
- `for-range` ikinci iterasiya dəyişəni (`for _, v := range`) — hər elementi `v`-yə kopyalayır!

### 7. Benchmark dərsləri: range və böyük elementlər
**Kitabdan kod nümunəsi (Example 2, [10]int64 elementli slice):**
```go
func Sum_PlainForLoop(s []Element) (r int64) {      // 911 ns
    for i := 0; i < len(s); i++ { r += s[i][0] }
    return
}
func Sum_OneIterationVar(s []Element) (r int64) {    // 929 ns
    for i := range s { r += s[i][0] }
    return
}
func Sum_UseSecondIterationVar(s []Element) (r int64) { // 3753 ns — 4× YAVAŞ!
    for _, v := range s { r += v[0] }
    return
}
```
**Sub-kod izahı:**
- `for _, v := range` → hər elementin 80 baytlıq kopyası `v`-yə yazılır — böyük elementlərdə qadağan
- `for i := range` + `s[i]` → kopya YOX, indekslə birbaşa oxu
- `v := &s[i]` → pointer ilə oxu — sahələr çox istifadə olunanda ən yaxşı variant

**Array parametr tələsi (Example 1):**
```go
func Sum_RangeArray(a [N]int) int      // 897 ns — kopya 2 dəfə!
func Sum_RangeArrayPtr1(a *[N]int) int // 799 ns — kopya 1 dəfə
func Sum_RangeArrayPtr2(a *[N]int) int // 555 ns — range pointer üzərindən, kopya 0
func Sum_RangeSlice(a []int) int       // 561 ns — ən idiomatik, kopyasız
```
**Dərs:** Böyük array-i funksiyaya **dəyər kimi** ötürmək hər çağırışda tam kopya deməkdir — pointer `*a` və ya slice `a[:]` ötür.

## Əsas terminlər
- Value part (dəyər hissəsi) — direct/indirect bölünmə
- Alignment guarantee (düzləndirmə zəmanəti)
- Struct padding (struktur doldurması) — israf olunan baytlar
- Small-size type (kiçik ölçülü tip) — xüsusi kopyalama optimizasiyası
- Value copy cost (dəyər kopyalama xərci) — ölçüyə mütənasib
- ns/op — benchmark nanosaniyə/əməliyyat vahidi

## Praktik nəticə
1. Struct sahələrini ölçüyə görə sıralamaq 33%-ə qədər yaddaş qənaəti verir — amma yalnız milyonlarla instans olan hot data-larda tətbiq et.
2. Hot path-də struct-ları ≤4 word saxla (9→10 sahə sıçrayışı 2-7× yavaşlama deməkdir).
3. `for _, v := range` yalnız kiçik elementlərdə təhlükəsizdir; böyük elementlərdə `for i := range` + indeks və ya `&s[i]` işlət.
4. Böyük array-ləri heç vaxt dəyər kimi ötürmə — slice və ya pointer işlət.
5. `unsafe.Sizeof`/`unsafe.Alignof` ilə real ölçüləri yoxla — fərziyyə yox, ölçü.

## Mənbə
Pages: 9-21 (PDF səh. 9-21)
