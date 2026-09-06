# Chapter 6 — İnterfeyslər və digər tiplər (səh. 24-35)

## Bu chapter nədən bəhs edir?

İnterfeyslərin mahiyyəti (davranış müqaviləsi, implicit təmin), tip çevirmələri (conversion), interfeysdən interfeysə type assertion/switch, "yalnız interfeys export et" prinsipi (crypto/cipher nümunəsi) və interfeyslərin metodlarla əlaqəsi — HTTP Handler-in struct/int/channel/FUNKSİYA ilə təmini.

---

## əsas fikirlər

### 1. İnterfeys = davranış təyini

"Nəsə bu əməli edə bilirsə, burada istifadə oluna bilər." 1-2 metodlu interfeyslər Go-da QANUNİDİR — ad metoddan törəyir: `io.Writer` (Write), `io.Reader` (Read), `fmt.Stringer` (String).

**Bir tip ÇOX interfeys təmin edə bilər.** Sequence nümunəsi — sort + String:

```go
type Sequence []int

// sort.Interface üçün:
func (s Sequence) Len() int           { return len(s) }
func (s Sequence) Less(i, j int) bool { return s[i] < s[j] }
func (s Sequence) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func (s Sequence) Copy() Sequence { /* kopya qaytar */ }

// Stringer üçün — çapdan əvvəl sort edir:
func (s Sequence) String() string {
    s = s.Copy()    // MÜDAFİƏ kopyası — arqumenti üstünə yazmırıq
    sort.Sort(s)
    return fmt.Sprint([]int(s))   // AŞAĞI BAX: çevirmə ilə standart çap
}
```

### 2. Çevirmələr (conversion) — tip dəyişikliyi metod dəsti dəyişir

`[]int(s)` — Sequence → []int: **YENİ dəyər YOX, mövcud dəyərin müvəqqəti yeni tip kimi görünməsi** (eyni alt struktur). Bu, idioma: **tipi çevirərək başqa metod dəstinə çıxış:**

```go
// Sequence-in 3 tipi 3 rolda:
func (s Sequence) String() string {
    s = s.Copy()
    sort.IntSlice(s).Sort()    // IntSlice-in metodu!
    return fmt.Sprint([]int(s))  // []int kimi çap
}
```

### 3. Type switch/assertion — interfeysdən qiymət çıxarma

```go
// Switch versiyası (fmt.Printf daxilisi kimi):
type Stringer interface{ String() string }
var value interface{}
switch str := value.(type) {
case string:
    return str                    // konkret tip
case Stringer:
    return str.String()           // interfeys→interfeys!
}
```

**Tək tip lazımdırsa — type assertion:**

```go
str := value.(string)          // yanlışsa RUNTIME PANIC!

// Təhlükəsiz — comma-ok:
str, ok := value.(string)
if ok {
    fmt.Printf("string: %q\n", str)
} else {
    fmt.Println("not a string")
}
// uğursuzsa str mövcuddur — zero value ("")
```

### 4. Yalnız interfeys export et

Tip YALNIZ interfeys təmin etmək üçün varsa və maraqlı export metodu YOXDURSA — tipin ÖZÜNÜ export ETMƏ. Konstruktor interfeys dəyəri qaytarsın:

```go
// crc32.NewIEEE və adler32.New hər ikisi hash.Hash32 qaytarır —
// CRC-32 → Adler-32 dəyişmək = konstruktor çağırışını dəyişmək ONLY.
```

**crypto/cipher misalı** — Block (blok şifr) + Stream (axın şifr) ayrılığı:

```go
type Block interface {
    BlockSize() int
    Encrypt(dst, src []byte)
    Decrypt(dst, src []byte)
}
type Stream interface {
    XORKeyStream(dst, src []byte)
}
// NewCTR hər hansı Block-dan Stream qurur — detal abstraksiya olunur:
func NewCTR(block Block, iv []byte) Stream
```

### 5. İnterfeys + metod = hər şey Handler ola bilər

**Handler interfeysi:**

```go
type Handler interface {
    ServeHTTP(ResponseWriter, *Request)
}
```

**4 fərqli receiver 4 dəfə təmin edir:**

```go
// 1. STRUCT:
type Counter struct{ n int }
func (ctr *Counter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    ctr.n++
    fmt.Fprintf(w, "counter = %d\n", ctr.n)
}

// 2. PRİMİTİV (int!):
type Counter int
func (ctr *Counter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    *ctr++
    fmt.Fprintf(w, "counter = %d\n", *ctr)
}

// 3. CHANNEL:
type Chan chan *http.Request
func (ch Chan) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    ch <- req
    fmt.Fprint(w, "уведомление отправлено")
}

// 4. FUNKSİYA — HandlerFunc adapter patterni:
type HandlerFunc func(ResponseWriter, *Request)
func (f HandlerFunc) ServeHTTP(w ResponseWriter, req *Request) {
    f(w, req)    // receiver=funksiya f — çağırır!
}

func ArgServer(w http.ResponseWriter, req *http.Request) {
    fmt.Fprintln(w, os.Args)
}
http.Handle("/args", http.HandlerFunc(ArgServer))   // çevirmə → metod çıxışı
```

**Moral:** interfeys = metod dəsti; metodlar (demək olar) hər tipə yaxışır — struct, int, channel, funksiya. `fmt.Fprintf`-in `http.ResponseWriter`-ə yazması eyni prinsipdir (ResponseWriter Write daşıyır → io.Writer).

---

## Əsas terminlər

| Termin | İzah |
|--------|------|
| -er konvensiyası | 1-metodlu interfeys: Reader/Writer/Stringer |
| İmplicit təmin | `implements` YOX — metodlar kifayətdir |
| Conversion | `[]int(s)` — müvəqqəti tip görünüşü; metod dəsti dəyişir |
| Type assertion | `value.(string)` / `v, ok := value.(string)` |
| Type switch | Çoxlu tip budaqlanması; interfeys case-ləri QARİŞDIRILA BİLƏR |
| Yalnız interfeys export | Konstruktor interfeys qaytarır (hash.Hash32) |
| HandlerFunc adapter | Funksiya → Handler: metod receiver-i funksiyadır |
| ResponseWriter | Write daşıyır → Fprintf hədəfi |

---

## Praktik nəticə

1. **Kiçik interfeyslər qur:** 1-2 metod — Reader/Writer patterni; çox interfeys bir tipdə birləşir.
2. **String metodu istehzaçı dən:** Arqumenti dəyişməmək üçün kopya + sort; `fmt.Sprint([]int(s))` çevirməsi O(N²)-dən qurtarır.
3. **Çevirmə ilə metod dəsti:** `sort.IntSlice(s).Sort()` — öz tipin metod yazmadan hazır imkanları al.
4. **Assertion daima comma-ok:** `value.(T)` təkbaşına xəta halında PANIC.
5. **Yalnız-interfeys tipləri gizlət:** Alqoritm dəyişməsi = konstruktor sətri — istifadəçi kodu toxunulmaz.
6. **Adapter pattern (HandlerFunc):** Adi funksiyanı interfeys dünyasına birləşdir — çevirmə kifayətdir.
7. **Metod hər tipə:** int, channel, funksiya — interfeys təminı struct-la məhdud deyil.

---

## Mənbə

- Sənəd: *Effective Go* (rus tərcüməsi), 2009
- PDF səhifələri: 24-35
