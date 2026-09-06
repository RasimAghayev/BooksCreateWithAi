# Chapter 9 — String and Byte Slices (String və bayt dilimləri)

## Bu chapter nədən bəhs edir?
String↔[]byte çevirmələrinin allocation qaydalarına, compiler-in 5 xüsusi çevirmə optimizasiyasına, string birləşdirmə üsullarının müqayisəsinə (+ vs strings.Builder vs []byte), strings.Compare-in tələlərinə və allocation azaldan praktik nümunələrə.

## Əsas fikirlər

### 1. Niyə çevirmə allocation doğurur?
**Nədir:** String = dəyişməz byte slice. Dəyişməzlik zəmanəti üçün string↔[]byte çevirməsində baytlar **dublikasiya olunmalıdır** (yeni blok + kopya).

**Amma:** Compiler 5 xüsusi halda dublikasiyanı atır:

### 2. Optimizasiya 1 — `range` daxilində []byte(s) allocation-sızdır
```go
func f(s string) {
    for i, b := range []byte(s) {  // []byte(s) allocate ETMİR
        _, _ = i, b
    }
}
```

### 3. Optimizasiya 2 — müqayisə operandı kimi string(slice) allocation-sızdır
```go
func Equal(a, b []byte) bool {
    return string(a) == string(b)   // bytes.Equal-in daxili implementasiyası!
}
```
**Tələ:** müqayisə **ifadə** olmalıdır — `switch string(x) { case string(y): }` formalı "təmiz" switch DƏYİŞDİRİCİ yaradır (3 allocation!), verbose `switch { case string(x) == string(y): }` isə yox (0):
```go
func verbose(x, y, z []byte) {  // 0 allocation — bəli, verbose daha sürətli!
    switch {
    case string(x) == string(y): // ...
    case string(x) == string(z): // ...
    }
}
```

### 4. Optimizasiya 3 — map OXUMA açarı kimi string(key) allocation-sızdır
```go
n = m[string(key)]           // get — 0 allocation
m[string(key)] = 123          // set — 1 allocation (modifikasiya!)
m[string(key)]++              // inc — 1 allocation (oxu+yazı birləşik)
```
**Maraq:** `m[string(key)]++` mantıqca `m[k] = m[k] + 1` kimi modifikasiya sayılır. Pointer elementli map isə oxuma kimi qalır:
```go
var m2 = map[string]*int{"key": new(int)}
*m2[string(key)]++            // 0 allocation! — p := m2[k]; *p++
```
**Tövsiyə:** Silinmə az, modifikasiya çox olan map-lərdə element tipini pointer et.

### 5. Optimizasiya 4 — sabit string-lə birləşmədə çevirmə allocation-sız
```go
// 1 allocation (yalnız nəticə):
return (" " + string(s) + string(s))[1:]
// 3 allocation:
return string(s) + string(s)
```
**Mexanizm:** Ən azı bir boş olmayan string sabiti olan concat ifadəsində çevirmələr yaddaş tutmur. Yalnız len > 32 slice-lərdə fərq görünür (≤32 artıq stack-dədir). 1 bayt israfını azaltmaq üçün bilinən baytı birinci qoy: `"$" + string(s[1:]) + ...`

### 6. Optimizasiya 5 — >32 baytlıq nəticələr həmişə heap-də
>32 nəticələndirən çevirmə/concat heap allocation demədir. 32-yə bölünmüş versiya 3→1 allocation endirir:
```go
str = string(s37[:32]) + string(s37[32:]) + string(s37[:32]) + string(s37[32:])
// 1 alloc, 80B — normal string(s37)+string(s37) isə 3 alloc, 176B
```
**Qeyd:** Bu "bölmə" hiyləsi künclüdür, gələcək versiyalar dəyişə bilər.

### 7. String birləşdirmə: 3 yol
** yol — bir ifadədə çoxlu +:**
- ≤32 bayt + escape etmirsə stack-də nəticə — ən yaxşı ssenari
- Bir statement-də mümkün olan hər yerdə + işlət (192 ns vs Builder 197 ns — bərabər, amma sadə)

** strings.Builder — kod-zamanı bilinməyən sayda:**
```go
var b strings.Builder
b.Grow(n)              // uzunluğu əvvəlcədən bilirsə — mütləq böyüt!
for _, s := range ss { b.WriteString(s) }
return b.String()
```
- Grow olmadan kapasite artıq böyüyə bilər → nəticə string blokun israfını miras alır
- Nəticə həmişə heap-də; make-in sıfırlaması Builder-in daxilində də var (amma unsafe ilə son allocation yoxdur)

** byte slice yolu — kiçik (≤64) nəticələrdə:**
```go
func Concat_WithBytes(ss ...string) string {
    var n = 0
    for _, s := range ss { n += len(s) }
    var bs []byte
    if n > 64 { bs = make([]byte, 0, n) }     // heap
    else       { bs = make([]byte, 0, 64) }   // STACK (sabit cap!)
    for _, s := range ss { bs = append(bs, s...) }
    return string(bs)
}
// len(s)=16: Bytes 208ns < Plus 235ns; len(s)=17: Bytes 358ns > Plus 236ns
```
**Bonus:** string-ə çevirməzdən əvvəl baytları redaktə etmək mümkündür.

### 8. String + byte slice birləşməsi
```go
// Way 1 (one-line): str BÖYÜKDÜRSƏ daha yaxşı — stringin dublikasiyası zəruri idi
newByteSlice = append([]byte(str), bs...)
// Way 2 (verbose): bs BÖYÜKDÜRSƏ daha yaxşı
newByteSlice = make([]byte, len(str)+len(bs))
copy(newByteSlice, str); copy(newByteSlice[len(str):], bs)
```
Hər ikisinin öz israfı var; ideali yoxdur — nisbətə görə seç.

### 9. strings.Compare yavaşdır (v1.19)
- `strings.Compare` = `if a==b; if a<b` — eyni uzunluq + uzun prefiks halında İKİ müqayisə
- `bytes.Compare` isə `bytealg.Compare` ilə optimize olunub
- **Üçtərəfli müqayisə** üçün `strings.Compare` istənilən halda məntiqlidir; sadə bərabərlik üçün `==`

**Müqayisə sıralaması qaydası (uzunluqlar tez-tez fərqlidirsə):** `x == y` (uzunluq fərqi → anlıq) `x < y`-dən çox sürətlidir → `case x == y`-ni default qoyma:
```go
func f1(x, y string) {        // DOĞRU sıra
    switch {
    case x == y: // ...
    case x < y:  // ...
    default:     // ...
    }
}
```

### 10. Allocation azaldan nümunələr

**Bir birləşmədən subslice:** 2 allocation əvəzinə 1:
```go
func f(a, b, c string) {
    abc := a + b + c
    ab := abc[:len(abc)-len(c)]   // abc-dən kəs — a+b üçün ayrıca concat YOX
}
```

**Composite açar — string concat əvəzinə array/struct:**
```go
var ma = make(map[[2]string]struct{})   // [a, b] açar
func fa(a, b string) { ma[[2]string{a, b}] = struct{}{} }  // 147ns, 0 allocs!
func fs(a, b string) { ms[a+"/"+b] = struct{}{} }           // 508ns, 3 allocs
```

**Case-insensitive müqayisə:**
```go
strings.EqualFold(a, b)              // 1271 ns, 0 allocs
strings.ToLower(a) == strings.ToLower(b)  // 7157 ns, 18 allocs — 5.6× yavaş!
```

**io.Writer-ə string yazmaq — BytesWriter:**
```go
type BytesWriter struct {
    io.Writer
    buf []byte                    // yenidən istifadə olunan buffer
}
func (sw *BytesWriter) WriteString(s string) (int, error) {
    for len(s) > 0 {
        n := copy(sw.buf, s)              // bufferə köçür
        n, err := sw.Write(sw.buf[:n])    // yaz
        s = s[n:]
        if err != nil || n == 0 { break }
    }
}
```
**Prinsip:** `w.Write([]byte(s))` hər çağırışda allocation — bufferli wrapper ilə sıfıra endir.

## Əsas terminlər
- Immutable (dəyişməz) — string-in əsas xassəsi
- Duplication avoidance (dublikasiyadan yayınma) — 5 compiler optimizasiyası
- strings.Builder — artımlı string qurucu (unsafe ilə son kopyasız)
- Grow — kapasiteni əvvəlcədən böyütmə metodu
- Three-way comparison (üçtərəfli müqayisə) — -1/0/+1
- strings.EqualFold — allocation-sız case-insensitive müqayisə
- Composite key (mürəkkəb açar) — `[N]string` / struct açar

## Praktik nəticə
1. Bir statement-də bitən concat üçün `+`; döngülü/qismən üçün `strings.Builder` + `Grow(n)`.
2. Map oxuma açarları `m[string(b)]` pulsuzdur; SET/`++` pulludur — sıx modifikasiyada pointer element düşün.
3. `switch string(x)` deyil, `switch { case string(x) == ... }` — 0 vs 3 allocation.
4. Case-insensitive: həmişə `EqualFold`, heç vaxt `ToLower` müqayisəsi.
5. Çoxhissəli açarlar üçün concat `a+"/"+b` əvəzinə `[2]string{a,b}` array açarı.
6. Writer-yə string axınında bufferli `WriteString` wrapper-i.
7. >32 baytlıq string əməliyyatları həmişə heap-dir — hot path-də sayını azalt.

## Mənbə
Pages: 90-105 (PDF səh. 90-105)
