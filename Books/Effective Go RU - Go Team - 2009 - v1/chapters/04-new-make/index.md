# Chapter 4 — Data: new/make, massiv, slice, map, çap, append (səh. 13-21)

## Bu chapter nədən bəhs edir?

Go-nun data təyinatının əsasları: `new` (zero value pointer) vs `make` (initializə olunmuş slice/map/channel), composite literal-lər (konstruktor əvəzi), massivlərin value-semantikası, slice-ların referans təbiəti, 2D slice qurulumu, map-lər (comma-ok idiomu), fmt çap formatları (%v/%+v/%#v/%T) və append builtin-inin dizaynı.

---

## Əsas fikirlər

### 1. new — zero value + pointer

```go
// new(T): yaddaş ayır, SIFIRLA, *T qaytar — "T-nin zero value-sına pointer"
p := new(SyncedBuffer)  // *SyncedBuffer — DƏRHAL İŞLƏK
var v SyncedBuffer       // SyncedBuffer  — DƏRHAL İŞLƏK
```

**Dizayn prinsipi — "zero value istifadəyə hazırdır":** `bytes.Buffer` zero value = boş buffer; `sync.Mutex` zero value = açıq mutex. Bu xassə **transitivdir**:

```go
type SyncedBuffer struct {
    lock   sync.Mutex
    buffer bytes.Buffer
}
// hər iki elan da hazırdır — initializasiya GƏRƏKSİZ
```

### 2. Konstruktorlar + composite literal

Zero value kifayət etmirsə:

```go
// Boilerplate:
func NewFile(fd int, name string) *File {
    if fd < 0 { return nil }
    f := new(File)
    f.fd = fd; f.name = name; f.dirinfo = nil; f.nepipe = 0
    return f
}

// Composite literal ilə:
func NewFile(fd int, name string) *File {
    if fd < 0 { return nil }
    f := File{fd, name, nil, 0}
    return &f
}
// Və ya bir sətirdə — lokalin adresini qaytarmaq Go-da QANUNİDİR:
return &File{fd: fd, name: name}   // sahə:dəyər — sıra əhəmiyyətsiz, yoxlananlar zero
```

**`new(File) == &File{}`** — sahəsiz literal = zero value.

Massiv/slice/map literal-ları (index/key etiketli):

```go
a := [...]string{Enone: "no error", Eio: "Eio", Einval: "invalid argument"}
s := []string{Enone: "no error", Eio: "Eio", Einval: "invalid argument"}
m := map[int]string{Enone: "no error", Eio: "Eio", Einval: "invalid argument"}
```

### 3. make — initializə olunmuş (YOX zero) dəyər

`make(T, args)` YALNIZ slice/map/channel; **T qaytarır (*T YOX)** — bu 3 tip istifadə öncəsi initializasiya tələb edən referans strukturlarıdır:

```go
make([]int, 10, 100)    // 100-lük array + len=10, cap=100 slice strukturu

// new vs make:
var p *[]int = new([]int)       // *p == nil — nadirən lazımlı!
var v []int  = make([]int, 100) // 100-lük işlək slice

// İdiomatik:
v := make([]int, 100)
```

### 4. Massivlər — VALUE semantika

```go
// Go-da massiv = DƏYƏR: assign BÜTÜN elementləri köçürür.
// Funksiyaya ötürülmə = KOPYA (pointer YOX!).
// Ölçü TİPİN hissəsidir: [10]int ≠ [20]int
```

C davranışı lazımdırsa pointer: `func Sum(a *[3]float64)`. Amma **idiomatik deyil — slice istifadə et.**

### 5. Slice-lər — referans + funksiyaya ötürülmə

Slice array-i sarır; assign = EYNİ array-ə istinad; element dəyişikliyi caller-a görünür:

```go
func (f *File) Read(buf []byte) (n int, err error)
// Read slice qəbul edir: uzunluq = oxu limiti (pointer+say əvəzi)
n, err := f.Read(buf[0:32])   // ilk 32 bayt — slicing effektiv
```

**Length/capacity:** len artırıla bilər (cap həddində); `cap` builtin max göstərir. Manual Append:

```go
func Append(slice, data []byte) []byte {
    l := len(slice)
    if l+len(data) > cap(slice) {  // realloсasiya
        newSlice := make([]byte, (l+len(data))*2)
        copy(newSlice, slice)
        slice = newSlice
    }
    slice = slice[0 : l+len(data)]
    copy(slice[l:], data)
    return slice    // QAYTARILMALI — slice strukturu (ptr/len/cap) VALUE kimi ötürülür!
}
```

### 6. 2D slice-lər

```go
type Transform [3][3]float64   // massiv massivləri
type LinesOfText [][]byte        // slice slice-ləri — sətirlər FƏRQLİ uzunluqlu!
```

2 yolla ayrılır: (a) sətir-sətir müstəqil `make` (böyüyə bilən sətirlər); (b) TƏK böyük array + slicing (statik — effektiv):

```go
// (b) tək ayrılma:
picture := make([][]uint8, YSize)
pixels := make([]uint8, XSize*YSize)
for i := range picture {
    picture[i], pixels = pixels[:XSize], pixels[XSize:]
}
```

### 7. Map-lər

Key: equality müəyyən olan hər tip (SLICE YOX!). Map = referans strukturu — funksiyaya ötürülsə dəyişiklik görünür.

```go
var timeZone = map[string]int{
    "UTC": 0 * 60 * 60, "EST": -5 * 60 * 60, /* ... */
}
offset := timeZone["EST"]        // olmayan key → zero value (0)
```

**Mövcudluq yoxlaması — comma-ok idiomu:**

```go
seconds, ok = timeZone[tz]         // ok=false → mövcud deyil
_, present := timeZone[tz]         // dəyər lazım deyilsə _
delete(timeZone, "PDT")           // silmə — olmasa belə təhlükəsiz
```

**Set-in emulyasiyası:** `map[string]bool` + `attended[person]` yoxlaması.

### 8. Çap — fmt

```go
fmt.Printf("Hello %d\n", 23)
fmt.Fprint(os.Stdout, "Hello ", 23, "\n")   // io.Writer-ə!
fmt.Println("Hello", 23)
```

C-dən fərqlər:
- **%d flag TƏLƏB ETMİR** — tipdən oxuyur: `x uint64 = 1<<64-1` → `%d %x` = 18446744073709551615 ffff...; `int64(x)` → -1 -1
- **%v** — universal (array/slice/struct/map də çap edir); map-lər lexikoqrafik SIRALANIR
- **%+v** — struct sahə adları ilə; **%#v** — Go sintaksisi; **%T** — tip; **%q** — string literal (tərs-dirnaq)

**String metodu ilə custom format:**

```go
func (t *T) String() string {
    return fmt.Sprintf("%d/%g/%q", t.a, t.b, t.c)
}
```

**Sonsuz rekursiya tələsi** (`%s` String çağırır → yenə Sprintf → ...):

```go
// PİS:
func (m MyString) String() string {
    return fmt.Sprintf("MyString=%s", m)   // REKURSİYA!
}
// YAXŞI:
return fmt.Sprintf("MyString=%s", string(m))   // əsas tipə çevir
```

**Variadic ötürmə:** `Sprintln(v...)` — `...` slice-i arqument siyahısına aça; `func Min(a ...int)` konkret tipli variadic.

### 9. Append — builtin olma səbəbi

```go
func append(slice []T, elements ...T) []T
```

Go-da T caller tərəfindən təyin olunan funksiya yazmaq mümkün deyil (generics öncəsi!) — ona görə **dilə daxilidir** (kompilyator dəstəyi tələb edir):

```go
x := []int{1, 2, 3}
x = append(x, 4, 5, 6)          // [1 2 3 4 5 6]
x = append(x, y...)             // slice birləşdirmə — ... ŞƏRT
```

---

## Əsas terminlər

| Termin | İzah |
|--------|------|
| `new(T)` | Zero-value T-yə pointer — SIFIRLA, qurMUR |
| Zero value hazırlığı | Buffer/Mutex/SyncedBuffer — dərhal istifadə; transitiv |
| Composite literal | `&File{fd: fd}` — konstruktoru sadələşdirir |
| `make(T, args)` | Slice/map/channel — initializə olunmuş T (YOX *T) |
| Massiv = value | Assign = tam kopya; ölçü tipin hissəsi |
| Slice strukturu | ptr+len+cap — value ötürülür → append qayıtmalı |
| 2D slice | Sətir-sətir YOXSA tək-array+slice |
| Comma-ok | `v, ok = m[k]` — mövcudluq ayırıcı idiom |
| `%v/%+v/%#v/%T` | Universal / sahəli / Go-sintaksis / tip formatları |
| String rekursiya tələsi | %s ilə Sprintf(String) — əsas tipə çevir |
| Builtin append | T tipi caller-dən asılı olduğundan dilə daxil |

---

## Praktik nəticə

1. **Zero value hazırlığı dizayn et:** struct-ı elə qur ki `new`/`var` sonrası dərhal işləsin (Buffer/Mutex patterni) — konstruktor ehtiyacını azaldır.
2. **Konstruktor = composite literal:** `&File{fd: fd}` — sahə adları ilə, sıradan asılı olmayan.
3. **new vs make:** slice/map/channel → YALNIZ make; T qaytarır. `new([]int)` → nil-slice pointer — istifadə etmə.
4. **Massiv parametr KOPYADIR** — C gözləntisi YOX; massiv ehtiyacı → slice.
5. **Append nəticəsini mənimsət:** `s = append(s, x)` — strukturu value-dur, amma element dəyişikliyi caller-da görünür.
6. **Map yoxlaması comma-ok:** yalnız zero value müqayisəsi yanlış (0 = mövcud ola bilər); `_` ilə mövcudluq.
7. **%v ailəsini bil:** %+v struct debug, %#v tam sintaksis, %T tip; String metodunda %s YOX — əsas tip.
8. **`append(x, y...)`** — slice birləşdirmə idiomu; `...` sizsiz kompilyasiya xətası.

---

## Mənbə

- Sənəd: *Effective Go* (rus tərcüməsi), 2009
- PDF səhifələri: 13-21
