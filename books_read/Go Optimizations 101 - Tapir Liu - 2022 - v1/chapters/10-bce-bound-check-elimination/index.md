# Chapter 10 — BCE (Bound Check Elimination) (Sərhəd yoxlamalarının aradan qaldırılması)

## Bu chapter nədən bəhs edir?
Bounds checking-in xərcinə və compiler-in BCE optimizasiyasına: hansı hallarda yoxlamalar avtomatik silinir, hansı hallarda kod-aiddə "hint" lazımdır, `-d=ssa/check_bce` diaqnostika aləti və BCE-friendly yazım üslubları.

## Əsas fikirlər

### 1. Bounds checking nədir və niyə xərclidir
**Nədir:** Memory safety üçün hər index/subslice əməliyyatında index-in aralıqdan çıxması yoxlanılır (`IsInBounds`/`IsSliceInBounds`) — aşmanda panic.

**Necə işləyir:** v1.7-dən BCE optimizasiyası zəruri olmayan yoxlamaları silir. **Tool:** `go run -gcflags="-d=ssa/check_bce" file.go` — qalan yoxlamaları göstərir.

### 2. Nə avtomatik BCE olur
- Array-lərdə sabit indekslər: `a[0], a[4]` (müddət tipin içindədir) — yoxlama yoxdur
- `range` və `len(s)` şərtli loop-larda bütün indeks/subslice-lar
- `for i := range s { _ = s[i] }` — hamısı BCE-olunur
- `if len(s) > 2 { s[0], s[1], s[2] }` — şərt yoxlamanı əvəz edir

**Kiçik boşluqlar (v1.19):** geri sayılan loop-da `s[:i+1]` subslice-i (`s[i]` olmadan) yoxlama saxlayır; `len(s) % 2` kimi modullo indekslər v1.19-a qədər alınmırdı.

### 3. Ən böyük indeksi ƏNVVƏL yaz — yoxlamalar azalır
**Kitabdan kod nümunəsi:**
```go
func f3a(s []int32) int32 {          // 4 bound check!
    return s[0] | s[1] | s[2] | s[3]
}
func f3b(s []int32) int32 {          // 1 bound check — s[3] keçsə, hamısı OK
    return s[3] | s[0] | s[1] | s[2]
}
```
**Prinsip:** Ən yüksək indeksli erişim öncə icra olunarsa, qalanları məntiqən təhlükəsizdir.

### 4. Compiler-ə hint ver — slice-ı yenidən təyin et
**Nədir:** Loop-daxili indekslərin aralığından əmin olmaq üçün slice-ı yenidən təyin (REASSIGN) etmək — `_ = is[:256]` kimi "atılan" ifadələr işləmir!

**Kitabdan kod nümunəsi:**
```go
func f4b(is []int, bs []byte) {
    if len(is) >= 256 {
        is = is[:256]              // UĞURLU hint — slice YENİDƏN təyin edildi
        for _, n := range bs {
            _ = is[n]              // BCEed!
        }
    }
}
func f4c(is []int, bs []byte) {     // işləməyən hintlər:
    if len(is) >= 256 {
        _ = is[:256]                // nəticə istifadə olunmur → kömək ETMİR
        for _, n := range bs { _ = is[n] }  // Found IsInBounds
    }
}
```
**Redundant if hint-i (maraqlı nümunə):**
```go
func NumSameBytes_2(x, y T) int {
    if len(x) > len(y) { x, y = y, x }
    if len(x) > len(y) { panic("unreachable") }   // əlbəttə çatmır — amma BCE-ə isbat!
    for i := 0; i < len(x); i++ {
        if x[i] != y[i] { return i }               // BCEed!
    }
    return len(x)
}
```

### 5. BCE-friendly yazım üslubları

**Subslice-la daralt (3-indeks üstünlük):**
```go
func f7c(s []byte, i int) {   // 1 yoxlama + 3477 ns — ƏN SÜRƏTLİ
    s = s[i:i+4:i+4]
    _ = s[3]; _ = s[2]; _ = s[1]; _ = s[0]
}
```
Diqqət: f7b (2-indeks) məntiqən 1 yoxlamalıdır amma praktikada ən YAVAŞ (4223ns) — üç-indeks həm BCE, həm performans qazandırır.

**İterasiyanı slice-kəsmə ilə sürətləndir:**
```go
func f8z(s []byte) {                  // ən sürətli pattern
    for i := 0; len(s) >= 4; i += 4 { // şərt: len(s) ≥ 4 — s[3] də təhlükəsiz
        _ = s[3]; _ = s[2]; _ = s[1]; _ = s[0]
        s = s[4:]                     // irəlilə
    }
}
```

**Loop şərtini len(buf) ilə yaz:**
```go
buf := make([]int, n+1)
for i := 0; i < len(buf); i++ { buf[i] = k }   // BCE!
for i := 0; i <= n; i++     { buf[i] = k }     // Found IsInBounds
```

**Sabit açılış > döngülü indeks (kiçik array-lər):**
```go
r[0] = x[3] ^ x[2] ^ x[1] ^ x[0]   // açılış — yoxlamasız, compiler sabit bilir
```

### 6. Qlobal slice-lar BCE düşmənidir — lokala köçür
```go
var s = make([]int, 5)
func fa0() { for i := range s { s[i] = i } }   // Found IsInBounds — qlobal
func fa1() { s := s; for i := range s { s[i] = i } }  // BCE — LOKAL kopya!
```
**Səbəb:** Qlobala başqa goroutine yaza bilər — compiler uzunluğu sabit saymır; lokal kopya sabitdir.

### 7. Array slice-dən BCE-dostdur
```go
var s = make([]int, 256)  // slice
var a = [256]int{}         // array
func fb1() int { return s[100] }  // Found IsInBounds
func fb2() int { return a[100] }  // BCE — ölçü TIPTƏ məlum!
func fc2(n byte) int { return a[n] }  // hətta dəyişən indeks — BCE (n < 256 tipin zəmanəti)
```
**Müqayisə dərinliyi:** byte-indeksli array sorğusu heç yoxlama istəmir — çünki byte ≤ 255 < 256.

### 8. BCE hələ bacarmadıqları (v1.19)
`copy`-dən sonra `s[n:]`; `i+N-1` formalı modullo addımlı looplar; funksiya-şərtli filter kimi `data[:k]` — gələcək versiyalarda düzələcək.

## Əsas terminlər
- Bounds checking (sərhəd yoxlaması) — IsInBounds / IsSliceInBounds
- BCE (Bound Check Elimination) — yoxlamaların silinməsi
- Compiler hint (kompilyator işarəsi) — yenidən təyin/tip zəmanəti
- check_bce — diaqnostika flag-i
- Index table (indeks cədvəli) — yoxlamasız array sorğusu

## Praktik nəticə
1. Hot loop-ları `-d=ssa/check_bce` ilə auditoriya et; qalan `Found IsInBounds`-ləri bərpa et.
2. Şərti indeksləmədə ən böyük indeksi birinci yaz.
3. Loop daxili slice-ları əvvəlcədən kəs: `is = is[:256]` (yenidən təyin!), `s = s[4:]` iterasiya.
4. Qlobal slice-ları funksiya daxilində lokal dəyişənə köçür.
5. Ölçüsü məlum konteynerlərdə slice əvəzinə array işlət.
6. Sabit ölçülü blok emalında 3-indeksli subslice (`s[i:i+4:i+4]`).

## Mənbə
Pages: 106-117 (PDF səh. 106-117)
