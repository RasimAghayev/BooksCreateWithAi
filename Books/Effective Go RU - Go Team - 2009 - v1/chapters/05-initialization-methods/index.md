# Chapter 5 — İlkinləşdirmə və metodlar (səh. 21-24)

## Bu chapter nədən bəhs edir?

İlkinləşdirmə quruluşu (const + iota, var ilə runtime ifadələri, init funksiyası — determinist sıra) və metodlar (hər adlandırılmış tipə metod, pointer vs value receiver qaydası, io.Writer implementasiyası nümunəsi).

---

## Əsas fikirlər

### 1. Konstantlar — yalnız compile-time

Konstantlar **kompilyasiya zamanı** yaradılır (lokal olsalar belə): rəqəm, rune, string, bool. İfadə **konstant ifadə** olmalıdır:

- `1<<3` — konstant ifadə ✓
- `math.Sin(math.Pi/4)` — XƏYR (runtime çağırış)

**iota enumerator** — ifadədə iştirak edir, ifadələr implicit təkrarlanır:

```go
type ByteSize float64
const (
    _  = iota            // 0-a görə _ — SKIP
    KB ByteSize = 1 << (10 * iota)   // iota=1 → 1<<10
    MB                   // implicit təkrar: 1<<20
    GB; TB; PB; EB; ZB; YB
)
```

### 2. Konstant tipinə String metodu

```go
func (b ByteSize) String() string {
    switch {
    case b >= YB: return fmt.Sprintf("%.2fYB", b/YB)
    // ... ZB EB PB TB GB MB KB
    }
    return fmt.Sprintf("%.2fB", b)
}
// YB → "1.00YB";  ByteSize(1e13) → "9.09TB"
```

**Rekursiya təhlükəsi YOXDUR burada:** `%f` float formatıdır — String çağırmır (yalnız string formatları String-i işə salır).

### 3. Variabllər — runtime ifadələr

```go
var (
    home   = os.Getenv("HOME")     // runtime — VAR üçün OK (const üçün YOX)
    user   = os.Getenv("USER")
    gopath = os.Getenv("GOPATH")
)
```

### 4. init funksiyası — determinist sıra

Hər fayl (birdən çox da) init ola bilər. **Sıra:** import olunan paketlər → paket var-ları → init. İstifadə: bəyan edilə bilməyən setup + **proqramın state düzgünlüyünün yoxlanması:**

```go
func init() {
    if user == "" {
        log.Fatal("$USER not set")
    }
    if home == "" {
        home = "/home/" + user
    }
    if gopath == "" {
        gopath = home + "/go"
    }
    flag.StringVar(&gopath, "gopath", gopath, "переопределить стандартный GOPATH")
}
```

### 5. Metodlar — hər adlandırılmış tipə

Receiver struct məcburiyyətində DEYİL — hər adlandırılmış tip (pointer/interfeys xaric) metod daşıya bilər:

```go
type ByteSlice []byte

func (slice ByteSlice) Append(data []byte) []byte { /* ... */ }
// Amma yenilənmiş slice qaytarmalı — jqənaətbəxş!
```

**Pointer receiver ilə caller-ın slice-unu dəyiş:**

```go
func (p *ByteSlice) Append(data []byte) {
    slice := *p
    // ... məntiq
    *p = slice     // qayıtış YOXDUR — dəyişiklik birbaşa caller-da
}
```

**io.Writer implementasiyası:**

```go
func (p *ByteSlice) Write(data []byte) (n int, err error) {
    slice := *p
    // ...
    *p = slice
    return len(data), nil
}
// *ByteSlice artıq io.Writer ŞƏRTLƏNİR:
var b ByteSlice
fmt.Fprintf(&b, "This hour has %d days\n", 7)   // format-yazma bufferə!
```

### 6. Pointer vs value qaydası

| Receiver | Çağırıla bilər |
|----------|----------------|
| **Value metodu** | Value ÜZƏRİNDƏ də, pointer ÜZƏRİNDƏ də |
| **Pointer metodu** | YALNIZ pointer üzərində |

**İstisna:** value ünvanlana biləndə (`b` kimi lokal), `b.Write(data)` → kompilyator `(&b).Write`-a çevirir.

**Səbəb:** pointer metodu value üzərində çağrılsa kopya dəyişər, dəyişiklik İTƏRDİ — dil bunu qadağan edir. Value metodu hər ikisində təhlükəsizdir (pointer avtomatik dereference).

Bu ideya `bytes.Buffer` implementasiyasının əsasıdır.

---

## əsas terminlər

| Termin | İzah |
|--------|------|
| Konstant ifadə | Kompilyasiya-time hesablanan — funksiya çağırışı YOX |
| iota | Enumerator — ifadədə təkrar, implicit davam |
| init funksiyası | Paket ilkinləşdirməsi: import→var→init sırası |
| Adlandırılmış tip receiver | Metod hər adlı tipə (struct şərt deyil) |
| Pointer receiver | Caller-ın dəyərini birbaşa mutasiya |
| Value metodu hər yerdə | Value+pointer üzərində çağrılır |
| Adresləmə avtomatikası | `b.Write` → `(&b).Write` |
| io.Writer şərtlənmə | `*ByteSlice` Write ilə — Fprintf hədəfi olur |

---

## Praktik nəticə

1. **Const = compile-time:** runtime hesabı lazımdırsa var; `1<<3` OK, `math.Sin(...)` YOX.
2. **iota + `_` skip:** KB-YB kimi qüvvət cərləri bir ifadə ilə.
3. **Scalar tipə String:** `%f` kimi qeyri-string formatlar recursion təhlükəsizdir.
4. **init = yoxlama+toxunuş:** dəyişənlərdən sonra işləyir — bərpaedilməz xətalar üçün Fatal.
5. **Slice mutasiya edən metod = pointer receiver:** `func (p *ByteSlice) Append(...)` — qayıtışı aradan qaldırır.
6. **Pointer metodu = yalnız pointer:** value üzərində çağrılış kompilyasiya xətası (kopya itkisi qarşısı); lokal value-da istisna — avtomatik `&`.
7. **Std interfeysləri öz tipinlə təmin et:** Write → Fprintf-in hədəfi — `bytes.Buffer`-in də yolu.

---

## Mənbə

- Sənəd: *Effective Go* (rus tərcüməsi), 2009
- PDF səhifələri: 21-24
