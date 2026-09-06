# Chapter 5 — Strings (#36-#41)

## Bu chapter nədən bəhs edir?

Bu chapter string-lərin daxili quruluşu (rune konsepti), dəqiqsiz iterasiya, trim funksiyalarının qarışdırılması, səmərəsiz birləşdirmə (`strings.Builder`), lazımsız string/[]byte konversiyaları və substring-lərdən yaranan memory leak-ləri əhatə edir.

---

## Əsas fikirlər

### #36: Not understanding the concept of a rune (Rune konsepti)

**Əsas anlayışlar:**

| Anlayış | Tərif |
|---------|-------|
| **Charset** | Simvollar toplusu (Unicode: 2^21 simvol) |
| **Encoding** | Simvolların ikili tərcüməsi (UTF-8: 1-4 bayt) |
| **Code point** | Unicode-da tək dəyərlə təmsil olunan vahid (汉 = U+6C49) |
| **Rune** | Go-da Unicode code point — `type rune = int32` alias |

汉 UTF-8-də 3 baytdır: `0xE6, 0xB1, 0x89` — code point 1-4 bayta encode oluna bildiyindən rune 32-bitdir.

**Vacib nüanslar:**
1. **Go string-ləri həmişə UTF-8 deyil.** Source code UTF-8-dir → string literal-lər UTF-8-dir. Amma string özü **istənilən bayt ardıcıllığıdır** — fayl sistemindən oxunan data UTF-8 olmaya bilər. (UTF-16/32 üçün golang.org/x paketləri.)
2. **`len(s)` bayt sayını qaytarır, simvol sayını YOX:**

```go
fmt.Println(len("hello")) // 5
fmt.Println(len("汉"))    // 3 — bir simvol, üç bayt!
```

3. **Baytlardan string yaratmaq:**

```go
s := string([]byte{0xE6, 0xB1, 0x89})
fmt.Printf("%s\n", s)    // 汉 — 3 bayt → 1 simvol
```

---

### #37: Inaccurate string iteration (Dəqiqsiz string iterasiyası)

```go
s := "hêllo"
for i := range s {
    fmt.Printf("position %d: %c\n", i, s[i])   // PİS!
}
fmt.Printf("len=%d\n", len(s))
```

Çıxış:

```
position 0: h
position 1: Ã       ← ê YOX!
position 3: l       ← position 2 hara getdi?
position 4: l
position 5: o
len=6               ← 5 rune var amma 6 bayt!
```

**Səbəblər:**
- `ê` 2 baytlıdır (0xC3 0xAA) → `len` = 6 bayt.
- Range string üzərində **hər rune-un BAŞLANĞIC İNDEKSİ** üzərində iterasiya edir — `i` rune indeksi deyil!
- `s[i]` — i-ci rune-u YOX, i-ci baytın UTF-8 təmsilini çap edir.

**Həll 1 — value dəyişənini istifadə et (iterasiya üçün ən effektiv):**

```go
for i, r := range s {
    fmt.Printf("position %d: %c\n", i, r)   // r — rune-un ÖZÜ
}
// position 0: h
// position 1: ê
// position 3: l   ← i hələ bayt başlanğıcıdır
// ...
```

**Həll 2 — []rune konversiyası (rune İNDEKSİ lazım olanda):**

```go
runes := []rune(s)
for i, r := range runes {
    fmt.Printf("position %d: %c\n", i, r)
}
// position 0: h ... position 4: o — i artıq rune indeksidir
```

Overhead: O(n) kopya + allocation → yalnız rune indeksi lazımdırsa istifadə et.

**İ-xinci rune-un çıxarılması:**

```go
r := []rune(s)[4]           // ümumi hal
fmt.Printf("%c\n", r)       // o

// Optimizasiya — string YALNIZ 1-baytlıq rune-lardan ibarətdirsə (A-Z, a-z):
fmt.Printf("%c\n", rune(s[4]))   // bayt = rune, kopyasız!
```

**Rune sayı:** `utf8.RuneCountInString(s)` → 5.

---

### #38: Misusing trim functions (Trim funksiyalarının qarışdırılması)

```go
fmt.Println(strings.TrimRight("123oxo", "xo"))   // 123 (set əməliyyatı!)
fmt.Println(strings.TrimSuffix("123oxo", "xo")) // 123o (suffix əməliyyatı)
```

| Funksiya | Nə silir | Təkrarlanır? |
|----------|-----------|---------------|
| `TrimRight(s, set)` | Sondakı bütün rune-lar SET-dən olduqca (geriyə iterasiya, set-dən olmayan rune-a qədər) | Bəli |
| `TrimLeft(s, set)` | Başdakı rune-lar SET-dən | Bəli |
| `TrimSuffix(s, suf)` | Verilmiş TƏK suffix | Xeyr |
| `TrimPrefix(s, pre)` | Verilmiş tək prefiks | Xeyr |
| `Trim(s, set)` | Hər iki tərəfdən set rune-ları | Bəli |

**Misal:** `TrimRight("123oxo", "xo")`: `o` set-dədir → sil; `x` set-dədir → sil; `o` set-dədir → sil; `3` set-də DEYİL → dayan → `123`.

`TrimSuffix("123xoxo", "xo")` → `123xo` (yalnız bir dəfə).

```go
fmt.Println(strings.TrimLeft("oxo123", "ox"))    // 123
fmt.Println(strings.TrimPrefix("oxo123", "ox"))  // o123
fmt.Println(strings.Trim("oxo123oxo", "ox"))     // 123
```

---

### #39: Under-optimized string concatenation (Səmərəsiz string birləşdirmə)

**Problem:** String İMMUTABLE-dır — `+=` hər iterasiyada YENİ allocation yaradır:

```go
// PİS — hər += yeni string allocasiyası:
func concat(values []string) string {
    s := ""
    for _, value := range values {
        s += value
    }
    return s
}
```

**Həll — `strings.Builder`:**

```go
func concat(values []string) string {
    sb := strings.Builder{}
    for _, value := range values {
        _, _ = sb.WriteString(value)   // internal buffere append
    }
    return sb.String()
}
```

- `WriteString` heç vaxt non-nil error qaytarmır; error `io.StringWriter` interfeysinə uyğunluq üçündür.
- Digər metodlar: `Write([]byte)`, `WriteByte(byte)`, `WriteRune(rune)`.
- Daxili byte slice append əsasında işləyir → konkurrent istifadə TƏHLÜKƏLİDİR (race).

**Preallocation — `Grow(n)`:**

```go
func concat(values []string) string {
    total := 0
    for i := 0; i < len(values); i++ {
        total += len(values[i])       // BAYT sayı (rune yox)
    }
    sb := strings.Builder{}
    sb.Grow(total)                     // total bayt üçün yer garantisi
    for _, value := range values {
        _, _ = sb.WriteString(value)
    }
    return sb.String()
}
```

**Benchmark (1,000 string × 1,000 bayt):**

| Variant | ns/op | Fərq |
|---------|-------|------|
| v1 (`+=`) | 72,291,485 | — |
| v2 (Builder) | 878,962 | ~82x sürətli |
| v3 (Builder + Grow) | 190,340 | **v1-dən 99%, v2-dən 78% sürətli** |

**İki dəfə iterasiya niyə sürətli?** #21 prinsipi — slice dolanda böyümək əvəzinə bir dəfə allocasiya.

**Qərar qaydası:** ~5-dən çox string birləşdirilirsə `strings.Builder`; bir neçə string üçün `+=` və ya `fmt.Sprintf` oxunaqlıdır. Gələcək ölçü məlum olsa `Grow` çağır.

---

### #40: Useless string conversions (Lazımsız string konversiyaları)

**Problem:** I/O əməliyyatları (`io.Reader`, `io.Writer`, `io.ReadAll`) `[]byte` ilə işləyir — string seçimi əlavə konversiyalar deməkdir:

```go
// PİS — 2 əlavə allocasiya + 2 kopya:
func sanitize(s string) string { return strings.TrimSpace(s) }
return []byte(sanitize(string(b))), nil   // []byte→string→[]byte!
```

**Həll — `bytes` paketi ilə []byte axını:**

```go
func sanitize(b []byte) []byte { return bytes.TrimSpace(b) }
return sanitize(b), nil    // heç bir konversiya
```

`bytes` paketi `strings`-in bütün əsas əməliyyatlarını təklif edir: `Split`, `Count`, `Contains`, `Index`, `TrimSpace` və s.

**String immutability sübutu:**

```go
b := []byte{'a', 'b', 'c'}
s := string(b)      // KOOPYA yaradılır
b[1] = 'x'
fmt.Println(s)      // abc — axc YOX!
```

`[]byte` → string konversiyası bayt kopyası tələb edir. Bütün workflow-u `bytes` ilə qurmağa çalış — I/O-da və I/O-suz.

---

### #41: Substrings and memory leaks (Substring-lər və memory leak)

**Substring sintaksisi:**

```go
s1 := "Hello, World!"
s2 := s1[:5]                    // Hello — 5 BAYT (rune yox!)

s1 := "Hêllo, World!"
s2 := string([]rune(s1)[:5])    // Hêllo — rune əsaslı kəsim üçün []rune
```

**Leak ssenarisi:** Log mesajları — UUID (36 simvol) + mesaj (minlərlə bayt). Son n UUID-ni cache-də saxla:

```go
// PİS — hər uuid BÜTÜN log mesajının backing array-ini yaşadır:
func (s store) handleLog(log string) error {
    if len(log) < 36 {
        return errors.New("log is not correctly formatted")
    }
    uuid := log[:36]     // eyni backing array-ə istinad!
    s.store(uuid)
    // ...
}
```

**Mexanizm:** Standart Go kompilyatoru substring ilə orijinal string-in EYNİ backing array-i paylaşmasına imkan verir (spesifikasiya tələb etməsə də — performans/memory baxımından optimal). Nəticə: 36 baytlıq uuid minlərlə baytlıq mesajı yaşadır → n mesaj = BÖYÜK yaddaş.

**Həll 1 — dərin kopya (manual):**

```go
uuid := string([]byte(log[:36]))   // []byte→string = yeni 36-baytlıq array
```

GoLand kimi IDE-lər "redundant conversion" xəbərdarlığı verə bilər — AMMA bu, real effektli əməliyyatdır; linter xəbərdarlıqları həmişə dəqiq deyil.

**Həll 2 — `strings.Clone` (Go 1.18+):**

```go
uuid := strings.Clone(log[:36])   // təzə allocation
```

**Vacib:** String əslində pointer olduğundan funksiyaya ötürülməsi dərin kopya YARATMIR — kopyalanan strukturlar eyni backing array-ə istinad edir.

---

## Əsas terminlər

| Termin | İzah |
|--------|------|
| Charset | Simvol dəsti (Unicode — 2^21 simvol) |
| Encoding | Simvolluğun ikiliyə tərcüməsi (UTF-8 — 1-4 bayt/simvol) |
| Code point | Tək dəyərlə təmsil olunan Unicode vahidi |
| Rune | Unicode code point — `type rune = int32` |
| `len(s)` | String-də BAYT sayı (rune yox) |
| Rune starting index | Range string-də `i` — rune-un bayt başlanğıcı |
| `utf8.RuneCountInString` | Rune sayı hesablanması |
| TrimRight/TrimLeft | Set-dəki rune-ları təkrar-təkrar silmə |
| TrimSuffix/TrimPrefix | Tək suffix/prefiks silmə (təkrarsız) |
| `strings.Builder` | Birləşdirmə üçün mutable buffer — immutable string-lərdən qurtuluş |
| `Grow(n)` | Builder-in daxili slice-inə n bayt preallocation |
| `io.StringWriter` | `WriteString(s) (n int, err error)` interfeysi |
| `bytes` paketi | `strings`-in []byte analoqları (Split, Contains, TrimSpace...) |
| String immutability | String dəyişməzdir — konversiya/+= həmişə kopyadır |
| Substring backing array | `s[:36]` orijinal string-in array-i ilə paylaşır — leak mənbəyi |
| `strings.Clone` | (Go 1.18+) Yeni allocation ilə string kopyası |

---

## Praktik nəticə

1. **`len(s)` baytdır:** simvol sayı üçün `utf8.RuneCountInString`; "汉" = 1 rune / 3 bayt.
2. **String source-dan gəlmirsə UTF-8 güvəni YOXDUR** — fayl/şəbəkə datasında ehtiyatlı ol.
3. **Range `i` = bayt başlanğıcı, `s[i]` = bayt:** rune üçün `r` (value) istifadə et; i-ci rune indeksi lazımdırsa `[]rune(s)[i]`; tək-baytlı alfabetdə `rune(s[i])` kopyasız optimizasiya.
4. **TrimRight/Left = set, TrimSuffix/Prefix = dəqiq:** "xo" set kimi `123oxo`→`123`, suffix kimi →`123o`.
5. **5+ string birləşdirmə → `strings.Builder`:** `+=` hər addımda allocation; ölçü məlum olsa `Grow(total)` — 99% sürət fərqi.
6. **I/O pipeline-larında `bytes` seç:** `bytes.TrimSpace` və s. — `[]byte→string→[]byte` zəncirindən qaç; hər konversiya = kopya + allocation.
7. **Substring saxlayırsansa kopyala:** `log[:36]` minlərbaytlıq mesajı yaşadır; `string([]byte(...))` və ya `strings.Clone(log[:36])` (Go 1.18+).
8. **IDE/linter xəbərdarlıqlarına kor inanma:** "redundant conversion" real effekt daşıya bilər.

---

## Mənbə

- Kitab: *100 Go Mistakes and How to Avoid Them* — Teiva Harsanyi, Manning, 2022 (ISBN 9781617299599)
- Chapter 5: "Strings", book səh. 113–125
- PDF səhifələri: 133–145
