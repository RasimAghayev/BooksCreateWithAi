# Chapter 6 — Working with Basic Types (səh. 88-111)

## Bu fəsil nədən bəhs edir?

Əsas tiplərin daxili mexanikası: comma-ok (map/kanal/type assertion), slice
üzərində Stack (yaddaş itkisi qoruması), append-in böyümə siyasəti və
benchmark optimizasiyası, JSON üçün xüsusi vaxt formatı, composite map
açarları və təbii dil vaxt parse (timezone-larla).

## Əsas fikirlər

### Recipe 31 — comma, ok paradigması
**Bug:** discounts map-ində olmayan meyvələr (lemon, banana) → zero value
0.0 → məbləğ 1.8 çıxır (4.4 yerinə)!

```go
// XƏTALI kod:
discount := discounts[li.Name]           // yoxdursa → 0.0!
total += li.Amount * li.Price * discount // 0-a vurulur

// DÜZGÜN:
discount, ok := discounts[li.Name]
if ok {
    total += li.Amount * li.Price * discount
} else {
    total += li.Amount * li.Price        // endirim yoxdur
}
```

**Comma-ok-un 3 yeri:**
```go
// 1) Map:
v, ok := m[key]

// 2) Kanal (bağlı kanaldan oxu → zero value):
ch := make(chan int)
close(ch)
val, ok := <-ch     // ok=false

// 3) Type assertion (panicsiz):
var i any = "hi"
// n := i.(int)         → PANIC
if n, ok := i.(int); ok {
    fmt.Println("int", n)
} else {
    fmt.Println("not an int")
}
```
- Zero value prediktabellik verir, amma "0 göndərilib" ↔ "yoxdur"
  fərqini yalnız comma-ok ayırır

### Recipe 32 — slice üzərində Stack
```go
type Token struct {
    Loc  int
    Char rune
}

type Stack []Token

func (s Stack) Len() int { return len(s) }

func (s *Stack) Push(tok Token) {
    *s = append(*s, tok)
}

var ErrEmpty = errors.New("empty stack")

func (s *Stack) Pop() (Token, error) {
    size := s.Len()
    if size == 0 {
        return Token{}, ErrEmpty
    }

    sl := *s
    val := sl[size-1]
    sl = sl[:size-1]

    // Yaddaş itkisi qoruması: >1k elementdə cap 2x azıbsa → kopyala:
    if len(sl) > 1024 && 2*len(sl) < cap(sl) {
        sl2 := make([]Token, len(sl))
        copy(sl2, sl)
        sl = sl2
    }

    *s = sl
    return val, nil
}
```

**Slice-in daxili strukturu (runtime/slice.go):**
```go
type slice struct {
    array unsafe.Pointer  // alt massivə pointer
    len   int
    cap   int
}
```
- `s2 := s1[1:3]` — half-open; HƏM s1, HƏM s2 eyni alt massivi göstərir
- Pop-dakı kiçilmə yoxlaması: slice alt massivin kiçik hissəsini göstərə
  bildiyindən GC bütün massivi azad etmir → memory leak; kopyalama ilə
  azad edilir

### Recipe 33 — cumSum benchmark əsaslı optimizasiya
**Problem:** böyük slice-lərdə yavaşlıq + GC aktivliyi.

```go
// Orijinal — boş slice + append:
func cumSum(values []int) []int {
    var cs []int                    // cap=0
    s := 0
    for _, val := range values {
        s += val
        cs = append(cs, s)          // hər dəfə böyümə riski
    }
    return cs
}
```

**Benchmark (yalnış hesabatı da göstərir):**
```
BenchmarkCumsum-12   13896   87945 ns/op   357626 B/op   19 allocs/op
```
- Gözlənilən yaddaş: 8 bayt × 9323 ≈ 74KB; real: 357KB, 19 alloc!
- pprof: `runtime.growslice` və GC işçi funksiyaları vaxtı yeyir

**Fix — əvvəlcədən capacity:**
```go
func cumSum(values []int) []int {
    cs := make([]int, 0, len(values))   // len=0, cap=len(values)!
    s := 0
    for _, val := range values {
        s += val
        cs = append(cs, s)              // heç vaxt böyümür
    }
    return cs
}
// Nəticə: 23712 ns/op, 81920 B/op, 1 allocs/op — 3.7x sürət!
```

**append-in böyümə siyasəti (ölçülmüş):**
- 1024-ə qədər: ×2 (1→2→4→...→1024)
- 1024-dən sonra: ~×1.25-1.36 (təxminən üçdə bir artım)
- Versiyalar arası dəyişə bilər, amma ümumi davranış eynidir

### Recipe 34 — JSON-a/JSON-dan vaxt (xüsusi format)
**Tapşırıq:** legacy format `20230421T153217.372` (YYYYMMDDTHHMMSS.MS).

```go
const JSONTimeLayout = "20060102T150405.000"

type JSONTime struct {
    time.Time          // embed — bütün time metodları əlçatandır
}

func (t JSONTime) MarshalJSON() ([]byte, error) {
    s := t.Format(JSONTimeLayout)
    return []byte(`"` + s + `"`), nil
}

func (t *JSONTime) UnmarshalJSON(data []byte) error {
    if len(data) < 2 {
        return fmt.Errorf("data too small: %q", string(data))
    }
    data = data[1 : len(data)-1]     // "" mötərizələrini kəs
    ts, err := time.Parse(JSONTimeLayout, string(data))
    if err != nil {
        return err
    }
    t.Time = ts
    return nil
}

// API modeli — daxili LogRecord-dan ayrı:
type APILogRecord struct {
    Time    JSONTime `json:"time"`
    Level   string   `json:"level"`
    Message string   `json:"message"`
}
```
- JSON-da Time tipi YOXDUR → adətən string (RFC3339) və ya epoch rəqəmi
- Daxili ilə API modelini ayır → daxili dəyişiklik API-yə toxunmur
- time.Time embed → Format/Parse/... birbaşa

### Recipe 35 — composite map açarları
```go
type key struct {          // unexported tip — toqquşma yoxdur
    year   int
    month  time.Month
    day    int
    symbol string
}

type InfoDB struct {
    m map[key]StockInfo    // struct = açar!
}

func (i *InfoDB) Get(symbol string, date time.Time) (StockInfo, bool) {
    k := key{date.Year(), date.Month(), date.Day(), symbol}
    info, ok := i.m[k]
    return info, ok
}
```
- Struct bütün sahələri comparable-dırsa, açar ola bilər
- String düzəltmə ("AAPL:20200302") təhlükəlidir: "joe42" — joe+42 yoxsa
  joe4+2? Öz serializasiya formatını icad etməyin
- Slice-lar comparable DEYİL — açarda istifadə etməyin; immutable tiplər
  (string, int, rune) üstünlük
- time.Time açarda yox — müqayisəsi problemi (saat qurşaqları);

### Recipe 36 — təbii dil vaxt parse
**Tapşırıq:** "2 weeks ago", "today", "5 hours ago", "2020-03-02T11:47",
"[Australia/Sydney]" zone-larla.

```go
// 1) Zone çıxarma (regex):
var tzRe = regexp.MustCompile(`\[.+\]`)

func extractLocation(query string) (*time.Location, string, error) {
    loc := tzRe.FindStringIndex(query)
    if loc == nil {
        return time.UTC, query, nil     // zone yoxdursa UTC
    }
    locName := query[loc[0]+1 : loc[1]-1]   // []-i at
    tz, err := time.LoadLocation(locName)
    if err != nil {
        return nil, query, err
    }
    query = strings.TrimSpace(query[:loc[0]])
    return tz, query, nil
}

// 2) "today" → günün başlanğıcı:
func today(loc *time.Location) time.Time {
    t := time.Now().In(loc)                  // zone-a çevir
    return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// 3) Delta parse ("3 days" → -3*24h):
var unitNames = map[string]time.Duration{
    "minute": time.Minute,
    "hour":   time.Hour,
    "day":    24 * time.Hour,
    "week":   7 * 24 * time.Hour,
}

func parseDelta(query string) (time.Duration, time.Duration, error) {
    var amount time.Duration
    var unit string
    _, err := fmt.Sscanf(query, "%d %s", &amount, &unit)
    if err != nil {
        return 0, 0, err
    }
    unit = strings.TrimSuffix(unit, "s")     // weeks → week
    d, ok := unitNames[unit]
    if !ok {
        return 0, 0, fmt.Errorf("unknown duration: %q", unit)
    }
    return -(amount * d), d, nil             // keçmişə = mənfi
}

// 4) Yuvarlaqlaşdırma:
func roundTime(t time.Time, delta time.Duration) time.Time {
    year, month, day := t.Year(), t.Month(), t.Day()
    hour, minute := t.Hour(), t.Minute()
    switch {
    case delta >= time.Hour:
        minute = 0
        fallthrough
    case delta >= 24*time.Hour:
        minute, hour = 0, 0
    }
    return time.Date(year, month, day, hour, minute, 0, 0, t.Location())
}

// 5) Əsas:
func parseTime(query string) (time.Time, error) {
    loc, query, err := extractLocation(query)
    if err != nil {
        return time.Time{}, err
    }
    if query == "today" {
        return today(loc), nil
    }
    // Əvvəlcə ISO format cəhdi:
    t, err := time.ParseInLocation("2006-01-02T15:04", query, loc)
    if err == nil {
        return t, nil
    }
    delta, round, err := parseDelta(query)
    if err != nil {
        return time.Time{}, err
    }
    t = time.Now().In(loc).Add(delta)
    return roundTime(t, round), nil
}
```

**Vaxt haqqında vacib qeydlər:**
- `time.ParseDuration`-da day/week YOXDUR — ona görə öz parseDelta
- `Truncate` absolute duration-la işləyir, təqdimetmə formu ilə YOX —
  öz roundTime yazıldı
- Zone DB kompüterdədir — Docker-da scratch/busybox images-lərində yoxdur;
  Go 1.15+: `time/tzdata` embed modulu
- Epoch (1 yanvar 1970), Y2038 problemi (32-bit overflow), DST keçidləri,
  NTP sıçrayışları — "Falsehoods programmers believe about time" oxuyun

## Əsas terminlər
- Comma, ok — mövcudluq siqnalı (map/kanal/assertion)
- Zero value — tipin default dəyəri (0.0, "", nil)
- Slice header (array+len+cap) — runtime strukturu
- Memory leak (slice-də) — alt massivin azad edilməməsi
- append growth policy — 1024-ə qədər ×2, sonra ~1.3x
- make(len, cap) — əvvəlcədən ayrılmış tutum
- Benchmark/profile — -benchmem, -cpuprofile, pprof
- json.Marshaler/Unmarshaler — xüsusi serializasiya
- Composite key — struct map açarı
- Comparable type — == ilə müqayisə oluna bilən
- time.Location/LoadLocation — zone idarəetməsi
- Epoch/unix time — 1970-dən saniyələr
- Truncate vs round — daxili və təqdimə formaları

## Praktik nəticə
Map oxuma həmişə comma-ok ilə — zero value 0.0 kimi tələyə çevrilir.
Slice ilə strukturlar quranda: Pop-da cap yoxlaması (memory leak),
doldurulan slice-lərdə `make(0, n)` (19 alloc → 1). JSON-a vaxt üçün
embed-li wrapper tip + Marshal/Unmarshal; daxili/API modelləri ayrı.
Composite açarlar üçün struct (string düzəltmə YOX); vaxt parse-da
zone-ları regex-lə ayır, delta üçün öz cədvəlin, Truncate-ə inanma —
yuvarlaqlaşdırmanı özün yaz.

## Mənbə
Pages: 88-111 (PDF 88-111)
