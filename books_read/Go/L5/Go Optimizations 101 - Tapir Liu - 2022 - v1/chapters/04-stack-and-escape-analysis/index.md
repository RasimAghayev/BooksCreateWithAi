# Chapter 4 — Stack and Escape Analysis (Stek və Escape Analysis)

## Bu chapter nədən bəhs edir?
Goroutine stack-lərinin işləməsinə, escape analysis-in nə vaxt dəyəri stack-ə, nə vaxt heap-ə qoyduğuna, funksiya stack frame-lərinə, stack böyümə/kiçilməsinə, dəyərin heap-ə "qaçmasına" səbəb olan hallara və allocation yerini idarə etməyin yollarına. Kitabın ən texniki fəsli.

## Əsas fikirlər

### 1. Goroutine stack-ləri
**Nədir:** Hər goroutine üçün runtime-un yaratdığı fasiləsiz yaddaş zonası (contiguous stack).

**Necə işləyir:** Stack allocation heap-dən **3 cəhətdən üstündür**: (1) tapmaq daha sürətli; (2) yığmaq lazım deyil — goroutine çıxanda bütöv yığılır; (3) CPU cache üçün daha dostdur. Ona görə compiler mümkün olduqca stack-i üstün tutur.

**Nəyə lazımdır:** Niyə `allocs/op` optimize etməyin əsas yolu stack-də qalmaqdır.

### 2. Escape analysis (kaçış analizi)
**Nədir:** Compiler-in hansı value part-ların yalnız bir goroutine-də istifadə olunduğunu təyin edib stack-ə yerləşdirmə modulu.

**Necə işləyir:** 
- **Əsas şərt:** dəyər yalnız cari goroutine-də istifadə olunmalıdır
- Dəyər birdən çox goroutine-də istifadə olunursa və ya compiler əmin ola bilmirsə → **heap**
- Bütün package-level dəyişənlər mütləq heap-dədir
- Heap-dəki dəyər tərəfindən istinad olunan dəyər də heap-ə düşür
- "Stack-ə bilər" ≠ "stack-də olacaq" — ölçü həddindən böyükdə heap-ə düşür (aşağıda)

**Tool:** `go run -gcflags=-m file.go` — escape nəticələrini göstərir:
```
./escape.go:6:3: moved to heap: a
./escape.go:7:3: b does not escape
```

**Kitabdan kod nümunəsi:**
```go
func main() {
    var (
        a = 1        // moved to heap: a — dəyişdirilir, 2 goroutine istifadə edir
        b = false    // does NOT escape — yalnız oxunur, gizli kopya yaradılır
        c = make(chan struct{})
    )
    go func() {
        if b { a++ }
        close(c)
    }()
    <-c
    println(a, b)
}
```
**Sub-kod izahı:**
- `a++` → `a` hər iki goroutine-də live → heap
- `b` heç dəyişdirilmir → compiler closure-a gizli kopya verir → stack-də qalır
- `b = !b` sətri əlavə ediləndə → `b` də heap-ə düşür

### 3. Stack frame-lər və niyə stack ucuzdur
**Nədir:** Funksiyanın stack-də işğal edəcəyi seqment; ölçüsü **compile zamanı** hesablanır (bütün kod budaqları nəzərə alınmaqla).

**Necə işləyir:** `go run -gcflags=-S` frame ölçülərini göstərir. Nümunə: `bar` funksiyası 5024 bayt (5000-baytlıq array + servis), `foo` 10056 bayt. Frame işarələnən kimi dəyərlərin ünvanları **artıq məlumdur** — runtime blok axtarmır → stack allocation-un sürətinin səbəbi.

### 4. Stack böyümə və kiçilmə
**Nədir:** Stack tələbat artdıqca böyüyür (yeni zona + köhnə istifadə hissənin kopyası + pointer-lərin yenilənməsi), GC dövründə lazımsız böyükdə yarıya kiçilir.

**Necə işləyir:**
- İlkin ölçü 2KiB (v1.19-dan adaptiv), ölçülər həmişə 2ⁿ
- Böyümə bahalıdır: kopyalama + bütün stack pointer-lərin yenilənməsi
- Kiçilmə yalnız GC dövründə, goroutine boş olanda; hər dəfə 2× azalır; tələbatın 4 qatından aşağı düşmür
- Maksimum: 64-bit-də 1GB default (`runtime/debug.SetMaxStack` ilə dəyişilir) — aşarsa proqram çökir

**Kitabdan kod nümunəsi:**
```go
func f(i int) byte {
    var a [1<<13]byte // 8KiB stack-də — stack-i böyüdür
    return a[i]
}
```

### 5. Dəyərin heap-ə qaçmasına səbəb olan klassik hallar (bir goroutine-də olsa belə!)

#### a) Loop daxilində elan edilən, loop-dan kənara istinad olunan dəyişən
```go
var x *int
for {
    var n = 1 // moved to heap: n
    x = &n
    break
}
```
**Səbəb:** n-in bir çox eyni-zamanlı instansı ola bilər; frame ölçüsü compile-da sabit olmalıdır.

#### b) Interface metod çağırışına ötürülən arqument
```go
func main() {
    var x int      // does not escape — konkret tip
    t.M(&x)        // T məlumdur → stack
    var y int      // moved to heap: y
    i.M(&y)        // interface — konkret metod bəlli deyil →保守 heap
}
```
**İstisna:** De-virtualizasiya — compiler `i`-nin dinamik tipini bilirsə (məs. yeni `var i I = t` sonra heç dəyişmirsə), stack-də saxlaya bilər.

#### c) `reflect.ValueOf` çağırışı
```go
var n = 1        // moved to heap: n
_ = reflect.ValueOf(&n)
```
**Vacib praktik nəticə:** `fmt.Println` daxildə reflect işlədir → argumentləri heap-ə qaçırır; builtin `println` — yox!
```go
var x = 1; fmt.Println(&x)  // x heap-də
var y = 2; println(&y)      // y stack-də — hot path-də fərq böyükdür
```

#### d) Funksiya nəticəsi tərəfindən istinad olunan dəyər
```go
//go:noinline
func f(x *int) *int {
    var n = *x + 1 // moved to heap: n — return &n
    return &n
}
```
**Səbəb:** Nəticənin caller-da necə istifadə olunmasını izləmək bahalıdır → hamısı heap.

### 6. Inlining escape analysis-i yüngülləşdirir
**Nədir:** Funksiya inline olunanda onun dəyişənləri caller-in lokal dəyişənlərinə çevrilir → escape analizi asanlaşır.

**Nümunə:** Yukarıdakı `f` funksiyasından `//go:noinline` çıxarılarsa, inline səbəbi ilə `n` **stack-də** qalır (çap olunan ünvanlar yanyana durur — sübut).

**Lakin:** v1.19-da constant-lar inline-də yaxşı yayılmır:
```go
func createSlice(n int) []byte { return make([]byte, n) }
var x = createSlice(32) // make([]byte, n) escapes to heap — baxmayara ki, 32 sabitdir
var y = make([]byte, 32) // does not escape
```

### 7. Allocation yerinə nəzarət — threshold-lər (v1.19)

| Hadisə | Hədd | Nəticə |
|--------|------|--------|
| string↔[]byte çevirmə | >32 bayt | heap; **sabit string üçün 64KB** |
| `new(T)` / `&T{}` | >64KB tip | heap |
| `make([]T, N)` — N **sabit** | >64KB backing array | heap |
| `make([]T, n)` — n **dəyişən** | n > 0 | **həmişə** heap |
| Dəyişən elanı (direct part) | >10MB | heap |

**Kitabdan sübut nümunəsi:**
```go
const N = 64 * 1024
func makeSlice65535() bool {          // 0 allocation
    s := make([]bool, N-1)            // does not escape
    ...
}
func makeSliceVarSize() bool {        // 1 allocation — n dəyişəndi
    s := make([]bool, n)              // escapes — n runtime bilinir
    ...
}
```
**Dərs:** Sabit ölçülü slice-lar stack bonusu ala bilər; runtime ölçü = mütləq heap.

### 8. Threshold-dən yuxarı stack alma hiylələri
**(a) Array-dən slice törət:**
```go
const N = 10 * 1024 * 1024 // 10M
func f() byte {
    var a [N]byte  // stack-də! (elən 10M həddindədir)
    var s = a[:]   // 10M elementli slice — stack backing array
    return s[n]
}  // 0 allocation
```
**(b) Composite literal (v1.19 bug-üstünlük):**
```go
var s = []byte{N: 0} // 500M element belə stack-də!
```
**(c) `-smallframes` flag-i:** 64K→16K, 10M→128K hədlərini endirir.

### 9. Stack-i az böyütmək — dummy funksiya hiyləsi
**Nədir:** Dərin rekursivada çoxqatlı stack böyüməsini bir dəfəyə reduction.

**Kitabdan kod nümunəsi:**
```go
func demo(n int) byte {
    var a [8192]byte
    if n--; n > 0 { b = demo(n) } // dərin rekursiya — çox böyümə
    return a[n] + b
}
func bar(c chan time.Duration) {
    // 64MiB frame ölçülü dummy anonim funksiya — stack-i əvvəlcədən pik hala gətirir
    func(x *interface{}) {
        type _ int // avoid being inlined
        if x != nil { *x = [1024 * 1024 * 64]byte{} }
    }(nil)
    demo(8192) // artıq böyüməsiz
}
// foo (böyüməsiz): 42.05ms → bar (bir böyümə): 4.74ms — 9× sürətli!
```
**Sub-kod izahı:**
- dummy funksiya heç vaxt işləməsə də, onun frame ölçüsü (67MB) stack-i dərhal 128MiB-ə qaldırır
- `type _ int` → inline olmasın deyə (inline olsa frame yoxa çıxar)

### 10. sync.Pool vs custom pool (Ch 3 davamı)
- `sync.Pool`: 2 GC dövründə istifadə edilməyən obyektlər avtomatik yığılır; ölçüsü dinamik
- Müəllifin fikri: real layihələrdə custom pool (max-size məhdudiyyətli) çox vaxt daha uyğundur

## Əsas terminlər
- Escape analysis (qaçış analizi) — stack/heap qərarı verən compiler modulu
- Stack frame (stek kadrı) — funksiya çağırışının stack seqmenti
- Contiguous stack (fasiləsiz stek) — bütöv seqment kimi stack
- De-virtualization (de-virtuallaşdırma) — interface çağırışının konkret tipe çevrilməsi
- Inlining (funksiyanın içinə yerləşdirilməsi) — çağırının call yerinə kod kopyalanması
- `-gcflags=-m` — escape nəticələrini göstərən flag
- Memory threshold (yaddaş həddi) — 32B / 64KB / 10MB qərar sərhədləri

## Praktik nəticə
1. Hot path-də `fmt.Println` əvəzinə builtin `println` (və ya logging-i hot path-dən kənarlaşdır) — reflect heap-ə qaçırır.
2. Interface metoda pointer ötürməkdən çəkin — mümkünsə konkret tip çağır.
3. `make([]T, n)`-də n runtime dəyəridirsə həmişə heap-dir — sabit ölçü stack şansı verir.
4. Böyük lokal buffer-lar üçün `var a [N]byte` + `a[:]` idiomu 64KB həddini 10MB-a qaldırır.
5. Dərin rekursiya + performans: başlanğıcda dummy funksiya ilə stack-i pik ölçüyə gətir.
6. Escape şübhəndə `-gcflags=-m` ilə yoxla — təxmin yox, compiler-in sözü.
7. Nəticə qaytararkən `&n` (lokalın pointeri) return etmək — həmişə heap allocation deməkdir; value type qaytar (nəticə kopyası ucuzdur, ≤4 word).

## Mənbə
Pages: 33-55 (PDF səh. 33-55)
