# Chapter 13 — Functions (Funksiyalar)

## Bu chapter nədən bəhs edir?
Funksiya inlining-in bütün aspektlərinə (hansı funksiyalar inline olur, inline cost, hot-path ayırma, manual inlining), pointer vs value parametrlərə, named result-lara, defer-in xərcinə, arqument qiymətləndirməsinə və hot path-də escape azaltmağa.

## Əsas fikirlər

### 1. Inlining nədir və nə üçün sürət verir
**Nədir:** Kiçik funksiya çağırışının çağırılan kodla əvəzlənməsi — stack əməliyyatları (arqument yazma, return oxuma) qənaət olunur.

**Tool:** `go build -gcflags="-m"` (sadə), `-gcflags="-m -m"` (səbəb + cost):
```
can inline bar with cost 14
cannot inline foo: cost 96 exceeds budget 80
```
**Inline cost:** hər statement-in xərci var; cəm threshold (v1.19-da 80) aşarsa inline olmurmə. Binary böyüməsinin qarşısını almaq üçün yalnız kiçik funksiyalar.

### 2. Hansı funksiyalar HEÇ VAXT inline olmur
- Rekursiv funksiyalar
- `recover()` çağırışı olanlar
- Tip elanı olanlar (`type _ int`)
- `defer` və `go` çağırışı olanlar (v1.19)
- **Funksiya-dəyişən çağırışı** — dəyəri compile-zamanı bilinmirsə:
```go
var addFunc = add               // qlobal — inline OLMUR (compile-time dəyəri bilinmir)
func main() {
    println(addFunc(11, 22))   // not inlined
    var addFunc = add            // LOKAL — inline OLUR (compiler lokalı təyin edir)
    println(addFunc(11, 22))   // inlined!
}
```

**Versiya keçidləri:** plain-for (v1.16+), for-range (v1.18+), tək-case select (v1.19+), closure (v1.17+) — getdi-gedə daha çox şey inline olur.

### 3. Az-inline-costlu yazım üslubları
**Kitabın sübutları (cost fərqləri):**
- Lokal elanlar cost artırır: `var r = 0` (25) vs named result `r int` (20) vs bare `return` (19)
- Bare return ucuzdur: `return r` (18) vs `return` (14)
- `||` zənciri ayrı if-lərdən ucuz: `if a||b||c||d` (12) vs 4 if (18)
- **for-range plain-for-dan ucuzdur** (v1.18+): 10 vs 18
- Funksiya-dəyişən = bahalı: switch+çağırış (30) vs funksiya seçimi+çağırış (84 — inline olmurmə!)
- **Manual inlining compiler-in costunu da azaldır:** cost 96 → 76 — eyni kod əllə açılanda!

### 4. Hot path-i inline-edilə bilən saxla — soyuq hissəni ayır
**Kitabdan kod nümunəsi:**
```go
// Orijinal: cost 85 → inline OLNMUR
func concat(bss ...[]byte) []byte {
    if len(bss) == 0 { return nil }
    else if len(bss) == 1 { return bss[0] }
    else if len(bss) == 2 { return append(bss[0], bss[1]...) }
    // ... soyuq general yol
}

// Yenidənqurulmuş: cost 74 → İNLİNE OLUR
func concat(bss ...[]byte) []byte {
    if len(bss) == 2 {                 // HOT yol (ən çox bu case)
        return append(bss[0], bss[1]...)
    }
    return concatSlow(bss...)          // soyuq yol — ayrı funksiya
}
//go:noinline                        // əks halda düzülüb cost-u geri artırar!
func concatSlow(bss ...[]byte) []byte { /* orijinal məntiq */ }
```
**Sub-kod izahı:**
- Əgər çağırışların 90%-i 2-slice birləşmədisə — hot yol həmişə inline icra olunur
- `//go:noinline` concatSlow-da MƏCBURİDİR — yoxsa compiler onu geri düzər (inline edib cost-u qaldırar)
- **Cost riyaziyyatı:** inline edilməmiş çağırış = 59 cost → 2+ belə çağırış = inline pozulur

### 5. Manual inline > auto inline (bəzən)
```go
Benchmark_ManualInline:  5.13 ns   // r[i&127] = *(*[N]byte)(buf) — birbaşa
Benchmark_AutoInline:   11.49 ns  // Slice2Array(buf) — inline olsa belə 2× yavaş!
```
Amma: qlobal array üzərində manual+auto inline NOT-inline-dən YAVAŞ ola bilər (127.9 → 196.4 ns — v1.19 compiler boşluğu).

### 6. Pointer vs value parametrlər — ölçüyə görə
```go
// 5-sahəli struct (20 bayt) — POINTER QAZANDIRIR:
Benchmark_Add5_TT_T (value): 17.73 ns
Benchmark_Add5_PPP (pointer): 11.95 ns

// 4-sahəli struct (16 bayt, small-size!) — VALUE QAZANDIRIR:
Benchmark_Add4_TT_T (value):  2.72 ns   // xüsusi optimizasiya!
Benchmark_Add4_PPP (pointer): 9.01 ns  // 3.3× yavaş
```
**Qayda:** ≤4 word struct → value param; böyüklər → pointer. (Pointer-in öz xərcləri də var: escape riski, GC scan.)

### 7. Named vs anonymous result — sabit deyil
```go
// Hal 1: named YAVAŞ (inline-lə bağlı boşluq):
ConvertToArray_Named:   472 ns   // (ret [N]byte) ... return
ConvertToArray_Unnamed: 333 ns   // [N]byte ... return *(*[N]byte)(b)
// Hal 2: named SÜRƏTLİ (bare return + sıfırlama qənaəti):
CopyToArray_Named:    408 ns    // (ret [N]byte) { copy(ret[:], b); return }
CopyToArray_Unnamed:  548 ns    // var ret [N]byte; copy; return ret
```
**Dərs:** Universal qayda YOXDUR — hər iki variantı benchmark et; kod üslubu oxunaqlıq seçsin.

### 8. Orta hesabları LOKAL dəyişəndə yığ
```go
func f(s []int) { for _, v := range s { sum += v } }     // qlobala hər iterasiyada — 3293 ns
func g(s []int) { var n = 0; for _, v := range s { n += v }; sum = n } // 654 ns — 5×
```
(Eyni Ch 6/7 prinsipi — qlobal yazışları loop-dan çıxart.)

### 9. defer loop daxilində = 10× xərc
```go
func f(n int) { for i := 0; i < n; i++ { defer inc(); inc() } }        // 61797 ns!
func g(n int) { for i := 0; i < n; i++ { func() { defer inc(); inc() }() } } // 5990 ns
```
**Sub-kod izahı:**
- v1.14+: loop XARİCİNDEKİ defer xüsusi optimizasiya olunur (open-coded defer)
- Loop DAXİLİNDEKİ defer — hər iterasiyada defer zəncirinə qoşulma
- g həlli: daxili anonim funksiya defer-i "loop-dan kənara" çıxarır — amma **məntiq fərqlidir** (defer-lər hər iterasiya sonunda işləyir, f-də isə funksiya sonunda)! Ekvivalensiya lazımdırsa istifadə etmə.
- **Ekstremal perf tələbində defer-dən tamamilə qaç** — deferli funksiya inline olmur.

### 10. Arqumentlər HƏMİŞƏ qiymətləndirilir
```go
debugPrint(h + w)                  // h+w concat BAZA hallarda da icra olunur — 1 allocation!
_ = debugOn && debugPrint(h + w)   // short-circuit: debugOn=false → çağırış YOX → 0 allocation
```
**Pattern:** Debug/trace funksiyalarını bool-and ifadəsinin sağ tərəfində çağır.

### 11. Hot path-də escape edən arqumentləri ayır
```go
func f(x int) string {          // x escapes to heap — çünki g(&x) filialı VAR
    if x >= 0 && x < 10 { return "0123456789"[x:x+1] }   // hot: x heap-e qaçmır amma...
    return g(&x)                // bu sətirə görə BÜTÜN çağırışlarda x heap-də
}
```
Hər filial üçün x-in gələcəyi bəlli olmalıdır — pointer ötürən filial hamını escape edir. Həll: hot filialı ayrı funksiyaya/sabit yoluna çəkmək.

## Əsas terminlər
- Inlining — çağırışın kodla əvəzlənməsi
- Inline cost / budget (80) — statement xərci cəmi / həddi
- Hot path / cold path (isti/soyuq yol) — sıx/nadir icra yolu
- Open-coded defer (v1.14+) — loop-xarici defer optimizasiyası
- Bare return — `return` (adsız nəticədən ucuz)
- Short-circuit evaluation — `&&` ilə qiymətləndirmə kəsilməsi
- `//go:noinline` — inline-i məcburi söndürmə

## Praktik nəticə
1. Hot kiçik funksiyaları `-m -m` ilə yoxla; cost 80 keçibsə — manual aç və ya soyuq hissəni `//go:noinline`-lı ayrı funksiyaya köçür.
2. Funksiya-dəyişən çağırışları (map[bool]func kimi) hot path-də işlətmə — inline pozulur.
3. ≤4 word parametrlər value; böyüklər pointer (benchmark et!).
4. defer-i loopdan çıxar (anonim funksiya ilə) və ya ekstremal halda defer-dən imtina.
5. Debug çağırışlarını `cond && f(args)` ilə qoru — arqument qiymətləndirməsi yalnız lazımda.
6. for-range plain-for-dan həm BCE, həm inline-cost baxımından üstündür (v1.18+).

## Mənbə
Pages: 130-149 (PDF səh. 130-149)
