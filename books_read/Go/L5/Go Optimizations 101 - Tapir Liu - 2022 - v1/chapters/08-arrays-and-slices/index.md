# Chapter 8 — Arrays and Slices (Massivlər və dilimlər)

## Bu chapter nədən bəhs edir?
Array/slice əməliyyatlarının performans detallarına: müqayisə, kopyalama (Go 1.17 slice→array-pointer), make+copy optimizasiyası, append-in böyümə alqoritmi, birləşdirmə/daxil etmə idiomları, range-in ikinci dəyişəni tələsi, memclr sıfırlama optimizasiyası və üç-indeksli subslice.

## Əsas fikirlər

### 1. Böyük array literal-ı müqayisə operandı kimi işlətmə
**Kitabdan kod nümunəsi:**
```go
type T [1000]byte
var zero = T{}                    // qlobal — bir dəfə

func CompareWithLiteral(t *T) bool { return *t == T{} }   // 52.2 ns
func CompareWithGlobalVar(t *T) bool { return *t == zero } // 31.0 ns — 1.7× sürətli
```
**İzah:** Literal hər çağırışda yenidən qurulur; qlobal dəyişən compiler üçün hazır. Kiçik array-larda fərq əhəmiyyətsiz.

### 2. Slice→array-pointer çevirmə ilə kopya (Go 1.17+)
**Nədir:** `*(*[N]byte)(d) = *(*[N]byte)(s)` — copy builtin-inin sürətli alternativi kiçik ölçülərdə.

**Kitabdan kod nümunəsi:**
```go
const N = 64
func copy2(d, s []byte) {
    *(*[N]byte)(d) = *(*[N]byte)(s)  // pointer çevirməsi ilə bir COPROKOPYA
}
// N=32: copy 4.58ns → copy2 1.40ns (3.3×)
// N=64: copy 5.19ns → copy2 2.00ns (2.6×)
```
**Sərhəd:** cəmi element ölçüsü ≤64 bayt olduqda istifadə et; [250..2250] və >1MB aralığında copy daha yaxşıdır.

### 3. make+copy optimizasiyası (v1.15+) — clone idiomu
**Nədir:** `y = make([]T, n); copy(y, x)` zəncirində make elementləri sıfırlamır — copy dərhal üstünə yazacaq.

**Şərtlər (hamısı olmalıdır):**
- klonlanan slice təmiz/sadə identifikator olmalı (`s` OK, `a[0]`/`ss.x` YOX)
- make YALNIZ 2 arqument almalı (`make([]T, n, m)` YOX)
- copy ifadə kimi istifadə edilməməli (`_ = copy(y,s)` və `f(copy(y,s))` YOX)

```go
// işləyir:
y = make([]T, len(s)); copy(y, s)
// işlemir:
_ = copy(y, s)               // nəticə istifadə olunmur
y = make([]T, len(s), len(s)) // 3 arqument
copy(y, a[0])                 // qualified identifikator
```

### 4. append böyümə alqoritmi (Go 1.18+)
```go
required := old.len + values.len
if required > old.cap*2 {
    newcap = required
} else if old.cap < 256 {          // v1.17-də hədd 1024 idi
    newcap = old.cap * 2           // kiçiklər: 2×
} else {
    newcap = old.cap               // böyüklər: 1.25×-ə yaxın artım
    for newcap < required { newcap += (newcap + 3*256) / 4 }
}
// sonra memsize → size class-a yuvarlaqlaşdırma
```
- **1.17 drawback:** nəticə kapasitesi monoton olmayan idi (897→2048, 1024→1280!); 1.18 düzəltdi
- Hər böyümə = 1 allocation → **bir addımda böyüt** (öncədən max cap hesabla)

### 5. Allocation edəcək append-də birinci arqumenti klip et
```go
x := make([]byte, 100, 500)
y := make([]byte, 500)
a := append(x, y...)                      // cap(a) = 896 — cap 500-ü miras aldı
b := append(x[:len(x):len(x)], y...)      // cap(b) = 640 — üç-indeks klipi
```
**İzah:** `x[len:cap]` forması (full slice expression) kapasiteni len-ə kəsir → yeni blok artıq köhnə cap-işlədmir. Daha sonra append olunmayacaqsa yaddaş qənaəti.

### 6. Clone: make+copy > append(nil)
```go
sCloned = make([]T, len(s)); copy(sCloned, s)   // optimal
sCloned = append([]T(nil), s...)                // çox vaxt daha yavaş
// x := make([]byte, 1<<15+1); y := append(nil, x...) → cap-y-len(x) = 8191 israf!
```

### 7. İki slice birləşdirmə
```go
// make+copy (sıralı):
merged = make([]T, len(x)+len(y)); copy(merged, x); copy(merged[len(x):], y)
// append (sıralı):    — len(y) >> len(x) isə daha yaxşı ola bilər
merged = append(x[:len(x):len(x)], y...)
```
- Sıra ƏHƏMİYYETSİZDİRSƏ: qısa slice-i birinci ötür — `append(qısa, uzun...)` daha az cap israf edir (768 vs 1360!)
- `x`-in boş cap-i `y`-ni saxlayırsa və paylaşma icazəlidirsə: `append(x, y...)` — allocation-sız, ən sürətli

### 8. 2+ slice birləşdirmə
```go
func MergeN_MakeAppend(ss ...[]byte) []byte {   // sadə və kifayət
    n := 0
    for _, s := range ss { n += len(s) }
    r := make([]byte, 0, n)
    for _, s := range ss { r = append(r, s...) }
    return r
}
```
- Orderless variant: ən böyük slice-i birinci yerə qoyub make+copy optimizasiyasından tam istifadə — amma verbose; yalnız güclü perf tələbində.

### 9. Slice daxilinə slice yerləşdirmə
```go
func Insert2(s []byte, k int, vs []byte) []byte {
    a := s[:k]
    s2 := make([]byte, len(s)+len(vs))
    copy(s2, a)                          // make+copy optimizasiyası sıfırlamadan
    copy(s2[len(a):], vs)
    copy(s2[len(a)+len(vs):], s[k:])
    return s2
}
```
- `append(x1[:k:k], append(vs, x1[k:]...)...)` klassik one-liner — İKİ dəfə kopya + 2 allocation → YAVAŞ, qaçın
- Boş cap kifayətdirsə: `s = s[:len(s)+len(vs)]; copy(s[i+len(vs):], s[i:]); copy(s[i:], vs)` — allocation-sız
- Tez-tez insert olunan data üçün linked list düşün

### 10. for-range ikinci dəyişən — hətta kiçik elementlərdə belə yavaş
```go
func sum_forrange1(s []int) int {          // 33793 ns — for i := range
    var n = 0; for i := range s { n += s[i] }; return n
}
func sum_forrange2(s []int) int {          // 37819 ns — for _, v := range (12% yavaş)
    var n = 0; for _, v := range s { n += v }; return n
}
func sum_plainfor(s []int) int {           // 33704 ns — klassik for
    var n = 0; for i := 0; i < len(s); i++ { n += s[i] }; return n
}
```
**Ekvivalensiya:** `for i, v := range anArray` ≡ array-in tam kopyası + elementlərin kopyası — array-də hər element İKİ dəfə kopyalanır!

### 11. memclr optimizasiyası — sıfırlamaq üçün for-range
```go
for i := range aSliceOrArray {
    aSliceOrArray[i] = v0     // compiler → daxili memclr — vectorized sıfırlama
}
```
- Klassik `for i := 0; i < len(...); i++` döngüsündən daha sürətlidir (len ≥ 6, byte)
- v1.19-dan array pointer üzərində də işləyir
- Array üçün daha sadə alternativ: `anArray = ArrayType{}`

### 12. Subslice-da kapasiteyi açıq göstər: `s[i:i+4:i+4]`
**Nədir:** İki-indeksli `s[i:i+4]` formasında compiler kapasitenin non-zero olduğunu çıxara bilmir → element pointer-inin blok sonuna işarə etməməsi üçün əlavə instruksiyalar yaradır.

**Kitabdan kod nümunəsi:**
```go
func g(rs, bs []byte) {
    for i, j := 0, 0; i < len(bs)-3; i += 4 {
        s2 := bs[i:i+4:i+4]   // üç-indeks: compiler bilir cap ≥ 4 → yoxlama yoxdur
        rs[j] = s2[3] ^ s2[2] ^ s2[1] ^ s2[0]
        j++
    }
}
```
**Qayda:** hot loop-dakı hər subslice ifadəsinə üçüncü indeksi əlavə et.

### 13. İndeks cədvəli ilə müqayisələri azalt (8.14)
Kiçik dəyər çoxluğu üçün map/switch əvəzinə array lookup — Ch 11.7 ilə eyni pattern.

## Əsas terminlər
- make+copy optimization — make-in sıfırlamasını atan compiler transformasiyası
- Full slice expression (tam dilim ifadəsi) — `s[a:b:c]` üç-indeks
- memclr — daxili vectorized sıfırlama
- Slice growth (dilim böyüməsi) — 2× / 1.25× alqoritm
- Backing array (dəstəkləyici massiv)
- Size class rounding — kapasitenin yaddaş sinfinə yuvarlaqlaşması

## Praktik nəticə
1. Clone: `make([]T, len(s)) + copy` — standart idiom (append(nil) yox).
2. ≤64 baytlıq kopyalar üçün `*(*[N]byte)(d) = *(*[N]byte)(s)` — 2-3× sürət.
3. Sıfırlama: `for i := range s { s[i] = zero }` — memclr-in açarı.
4. Hot subslice-lərdə həmişə üçüncü indeks: `s[i:j:j]`.
5. Birləşdirmədə nəticə uzunluğu əvvəlcədən hesabla, bir allocation ilə böyüt.
6. `for _, v := range`-i yalnız oxunaqlıq üçün istifadə et; perf-kritik loopda `for i := range`.
7. Daha artıq append olunmayacaqsa `append(x[:len(x):len(x)], ...)` klipi ilə yaddaş saxla.

## Mənbə
Pages: 76-89 (PDF səh. 76-89)
