# Chapter 14 — Interfaces (İnterfeyslər)

## Bu chapter nədən bəhs edir?
Value boxing-un xərc təsvilinə (hansı tiplər pulsuz, hansılar 50× bahalı), boxing optimizasiyalarının tam siyahısına, interface→interface assignment qənaətinə, virtual table çağırış xərcinə və sıx çağırılan funksiyalarda interface parametrlərdən yayınmağa. Kitabın son fəsli.

## Əsas fikirlər

### 1. Boxing nədir və nə qədər bahalıdır
**Nədir:** Interface dəyəri = qutudur; non-interface dəyərin interface-ə təyinatı = boxing (kopya qutulanır), type assertion = unboxing.

**Ümumi qayda:** Hər boxing ~1 allocation + dəyərin ölçüsü qədər kopya. Amma ölçüyə görə fərq İLLƏTİVARİDİR:

| Boxing | ns/op | alloc |
|--------|------|-------|
| int16 | 18.4 | 2 B |
| int64 | 55.9 | 8 B |
| string | 104 | 16 B |
| slice | 115 | 24 B |
| [100]int array | **592** | 896 B! |

### 2. PULSUZ boxing-lər (compiler optimizasiyaları)
Bu hallarda boxing allocation-sızdır və pointer qədər sürətlidir (~1.2 ns):
- **Pointer dəyərlər** — və buna görə map/channel/function da (daxildə pointerdir!)
- **Sabit dəyərlər** (constant int64, float64, string)
- **Zero-size** (`struct{}`, `[0]T`)
- **bool və 8-bit int** (birbaşa qutuda saxlanılır)
- **[0, 255] aralığında kiçik int-lər** (~3.4 ns — runtime staticbytes cache)
- **Zero dəyərlər** — 0 float, "" string, nil slice (~3.4 ns)
- **Tək sahəli struct/tək elementli array** — sahənin öz qədər (sahə pulsuzdırsa, struct da)

**Bahalı olanlar:** kənar [0,255] int / non-zero float → ~20× pointer-dən yavaş; non-nil slice / non-constant string → **~50× yavaş**.

### 3. Lookup table ilə bahalı boxing-i ucuzlaşdır
**Nədir:** Dəyərlərin kiçik çoxluğu məlumdirsə, qlobal array-dan pointer götür — pointer boxing pulsuzdur.

**Kitabdan kod nümunəsi:**
```go
var values [65536]uint16      // bütün uint16 dəyərləri saxlayır
func init() { for i := range values { values[i] = uint16(i) } }

r = uint16(i)            // Box_Normal: 22.7 ns
r = &values[uint16(i)]   // Box_Lookup: 1.14 ns — 20× sürətli, 0 allocation!
```
**Tətbiq:** Enum tipli dəyərləri interface ilə daşıyanda (API cavab kodları, statuslar).

### 4. İki interface-ə box etmək lazımdırsa — BİR box-la, assign et
**Kitabdan kod nümunəsi:**
```go
// Pis: 2 allocation
x = v   // box #1
y = v   // box #2
Benchmark_BoxBox:     130.5 ns, 2 allocs

// Yaxşı: 1 allocation
x = v   // box
y = x   // interface→interface assign — ƏMƏLİYYATSIZ!
Benchmark_BoxAssign:   68.3 ns, 1 alloc — 2× sürətli
```

**Praktik tətbiq — fmt:**
```go
var x = "aaa"
fmt.Fprint(io.Discard, x, x, x)              // 3 allocation — hər arqument ayrıca box!
var i interface{} = x                          // 1 allocation
fmt.Fprint(io.Discard, i, i, i)               // 0 əlavə — EYNİ dəyər bir qutuda
```
**Qayda:** Eyni dəyər bir fmt çağırışında neçə dəfə keçirsə — əvvəlcədən bir dəfə interface dəyişəninə box et.

### 5. Interface metod çağırışı = vtable + inline pozulması
**Kitabdan benchmark:**
```go
type BinaryOp interface { Do(x, y float64) float64 }
type Add struct{}
func (a Add) Do(x, y float64) float64 { return x + y }

Benchmark_Add_Inline:       0.63 ns   // konkret tip — inline
Benchmark_Add_NotInlined:   2.34 ns   // konkret tip — inline-siz
Benchmark_Add_Interface:    4.94 ns   // interface — vtable + inline-siz: ~8× yavaş
```
**Dekompozisiya:** vtable lookup ≈ 2.5ns; inline itkisi ≈ 1.7ns.

**Kitabın balansı:** Interface-dən qorxma — təmiz dizayn kiçik perf-dən vacibdir; compiler de-virtualizasiya edə bilir (Ch 4.6.2-dəki kimi). Amma sıx hot path-də konkret tip çağır.

### 6. Sıx çağırılan kiçik funksiyalarda interface parametr/result YOX
**Nümunə — image paketinin dərsi:**
```go
// PIS dizayn (image.Image — Go 1.16-dək):
type Image interface {
    At(x, y int) color.Color        // interface RESULT — hər piksel üçün boxing!
    Set(x, y int, c color.Color)    // interface PARAMETR
}
// Gözəl dizayn (image/draw.RGBA64Image, Go 1.17+):
type RGBA64Image interface {
    image.Image
    RGBA64At(x, y int) color.RGBA64       // KONKRETT tip result — boxing yox
    SetRGBA64(x, y int, c color.RGBA64)   // KONKRETT tip parametr
}
```
**Sub-kod izahı:** Milyonlarla piksel emalında hər At/Set boxing allocation-I doğururdu; RGBA64 konkret tipi hər ikisini aradan qaldırdı — standart kitabxananın öz dərsi: **interface sərhədlərində konkret tiplər tərəfindən keç**.

## Əsas terminlər
- Boxing / Unboxing (qutulama/açma) — interface dəyər konversiyası
- Virtual table (vtable) — dinamik metod cədvəli
- De-virtualization — compile-time konkret tipə qayıdış
- Staticbytes cache — [0,255] kiçik int boxing buferi
- Lookup table (axtarış cədvəli) — dəyər→pointer xəritəsi
- Interface-to-interface assignment — qutunu paylaşma (allocation-sız)
- `image/draw.RGBA64Image` — konkret-tip API dizayn nümunəsi

## Praktik nəticə
1. Hot path-də interface-ə box olunacaq dəyərlərdə pointer işlət (metodları `*T` üzərində elan et — T üzərində yox).
2. Eyni dəyər çox arqumentdə keçirsə — bir dəfə box et (`var i interface{} = x`), sonra `i` ötür.
3. Kiçik dəyər çoxluqları üçün qlobal lookup array — `&values[n]` boxing-i 20× ucuzlaşdırır.
4. Piksel/bayt səviyyəli API-lərdə konkret tip parametr/result (RGBA64 modeli); interface yalnız sərhəddə.
5. 592ns-lik array boxing — böyük struct-ları heç vaxt birbaşa box etmə, pointer box et.
6. Interface metod çağırışı hot loop-da ~8× — de-virtualizasiya şansı olsa da, kritik yerlərdə konkret saxla.

## Mənbə
Pages: 150-161 (PDF səh. 150-161)
