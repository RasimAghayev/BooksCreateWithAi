# Chapter 3 — Composite Data Types (Kompozit Data Tipləri)

## Bu chapter nədən bəhs edir?

Maps (yaradma, nil map, iterasiya), structures (type keyword, new, pointer vs value, slice of
struct), regular expressions (regexp paketi, Compile vs MustCompile, pattern cədvəli),
raw vs interpreted string literal, CSV fayl emalı (encoding/csv) və statistika tətbiqinin
CSV-dən data oxuyan funksional versiyası.

## Əsas fikirlər

### 1. Maps — Dinamik Açarlı Struktur
**Nədir:** Açar-dəyər cütlüyü; açar tipi comparable olmalıdır. Array/slice-dən fərqi:
müsbət tam ədəd indexlər və boşluqsuz ardıcıllıq tələbi YOXDUR.

**Nəyə lazımdır:** Sabit vaxtda (constant time) element axtarışı/access — index 99-a yazmaq
üçün 100 elementlik slice ayrılmır.

**Üstünlükləri:**
- Versatile — database index kimi belə istifadə oluna bilər
- Insert/retrieve — constant time
- Aydın dizayn

**Çatışmamazlıqları:**
- Açar kimi bool mənasızdır (yalnız 2 dəyər); float açar dəqiqlik bug-ları riskidir
- Data locality yoxdur (arrays/slices ilə müqayisədə)

**Kitabdan kod nümunəsi:**
```go
// make ilə:
aMap := make(map[string]int)
// literal ilə (yaratma anında data qoşmaq üçün daha sürətli):
m := map[string]int{"key1": -1, "key2": 123}

len(aMap)                  // açar sayı
delete(aMap, "key1")       // cütü sil

// Açarın mövcudluğu:
v, ok := aMap[k]           // ok=true → k var, v dəyəri; ok=false → v zero value
```
**Vacib qaydalar:**
- Olmayan açardan oxu → xəta YOX, zero value qaytarır
- **Iterasiya sırası RANDOM-dır** — dilin şüurlu dizayn qərarı; sıra üzərində pressumpt YOX

### 2. nil Map — Xüsusi Hal
**Nədir:** `aMap = nil` → map dəyişəni heç yerə işarə etmir; yazma cəhdi PANIC.

**Kitabdan kod nümunəsi:**
```go
aMap := map[string]int{}
aMap["test"] = 1      // OK — mövcud map

aMap = nil
if aMap == nil {      // yoxlama GOOD PRACTICE
    aMap = map[string]int{}
}
aMap["test"] = 1

aMap = nil
aMap["test"] = 1      // panic: assignment to entry in nil map!
```
**Praktik qayda:** map qəbul edən funksiya daxilə girişdə nil yoxlamalıdır.

### 3. Map İterasiyası
```go
for key, v := range aMap { }   // həm açar, həm dəyər
for _, v := range aMap { }     // yalnız dəyərlər
```

### 4. Structures — Qruplaşdırılmış Data
**Nədir:** Müxtəlif tipli sahələri tək ad altında birləşdirən kompozit tip; metodlar
(functions attached) daşıya bilər.

**type keyword:** `type myInt int` — yeni adlı tip; myInt və int Go üçün TAM FƏRQLİ
tiplərdir, birbaşa müqayisə olunmur. Sahə sırası tip identikliyinin hissəsidir — eyni
sahəli 2 struct, fərqli sıra ilə = fərqli tiplər.

**Kitabdan kod nümunəsi:**
```go
type Entry struct {
    Name    string
    Surname string
    Year    int
}

// 4 yaradma üsulu:
func zeroS() Entry { return Entry{} }              // Go zero-value ilə init
func initS(N, S string, Y int) Entry {              // istifadəçi ilə init
    if Y < 2000 { Y = 2000 }                        // input validation burada!
    return Entry{Name: N, Surname: S, Year: Y}
}
func zeroPtoS() *Entry { return &Entry{} }          // pointer variantı
func initPtoS(N, S string, Y int) *Entry {
    if len(S) == 0 { S = "Unknown" }
    return &Entry{Name: N, Surname: S, Year: Y}
}

pS := new(Entry)   // yaddaş ayır + sıfırla + POINTER qaytar
```
**Sub-kod izahı:**
- Zero value struct → bütün sahələr öz zero value-sunda ("" , 0)
- `new(T)` — həmişə pointer qaytarır; channel və map ÜÇÜN işləmir
- Böyük sayda struct init edirsənsə funksiya yaratmaq good practice — daha az xətaya meylli
- Struct definition funksiya XARİCİNDƏ (qlobal scope) olur
- Böyük hərflə başlayan sahə = package xaricinə görünür (Ch6-də ətraflı)
- Struct-i struct-in İÇİNƏ EMBED etmək (definition) BAD PRACTICE; hazır struct tip kimi
  saxlamaq isə normal

### 5. Slices of Structures
**Kitabdan kod nümunəsi:**
```go
type record struct {
    Field1 int
    Field2 string
}

s := []record{}
for i := 0; i < 10; i++ {
    text := "text" + strconv.Itoa(i)
    s = append(s, record{Field1: i, Field2: text})
}

fmt.Println(s[0].Field1, s[0].Field2)   // element → .sahə

sum := 0
for _, k := range s {
    sum += k.Field1
}
```

### 6. Regular Expressions — Nəzəriyyə
**Nədir:** Axtarış pattern-i təyin edən simvol ardıcıllığı; compile → finite automaton
(deterministic və ya nondeterministic) → recognizer.

**Compile vs MustCompile:**
```go
exp1, err := regexp.Compile(re)     // (Regexp, error) — xəta idarə olunur
exp2 := regexp.MustCompile(re)      // yalnız Regexp — xəta PANIC
```
- MustCompile başlanğıcda səhvi dərhal üzə çıxarır (init zamanı yoxlanış)
- **Go lookahead `(?=)` və lookbehind DESTƏKLƏMİR** — parse xətası

**Pattern cədvəli:**
| Simvol | Məna |
|---|---|
| `.` | hər hansı simvol |
| `*` | 0 və ya daha çox |
| `+` | 1 və ya daha çox |
| `?` | 0 və ya 1 |
| `^` / `$` | sətir başı / sonu |
| `[A-Z]` | hərf diapazonu |
| `\d` / `\D` | rəqəm / qeyri-rəqəm |
| `\w` / `\W` | söz simvolu [0-9A-Za-z_] / əksinə |
| `\s` / `\S` | whitespace / qeyri-whitespace |

**Kitabdan kod nümunələri:**
```go
// Ad/soyad: Böyük hərflə başla, kiçiklə davam et:
re := regexp.MustCompile(`^[A-Z][a-z]*$`)
re.Match([]byte(s))

// İştirakı imzalı/tam ədəd:
re := regexp.MustCompile(`^[-+]?\d+$`)
```
**Müəllim məsləhəti (kitabdan):** regex bug mənbəyidir — mümkün olan ən sadə pattern;
regex-siz həll varsa (məs. tam ədəd üçün strconv.Atoi) onu seç; regex bilinməyən
strukturda input ayırmaq üçün qiymətlidir.

### 7. Raw vs Interpreted String Literal
- Interpreted `"..."` — escape simvolları (\n və s.) emal olunur
- Raw `` `...` `` — escape emal YOX: çoxsətirli mətn, regex (backslash escape lazımsız),
  struct tag-ları üçün ideal

### 8. CSV Fayl Emalı — encoding/csv
**Nədir:** io.Reader/io.Writer interfeysləri üzərində işləyən CSV oxu/yaz paketi.
CSV-də hər şey STRING-dir — numeric data özün çevirməlisən.

**Kitabdan kod nümunəsi:**
```go
// OXU:
func readCSVFile(filepath string) ([][]string, error) {
    _, err := os.Stat(filepath)      // fayl mövcuddurmu?
    if err != nil { return nil, err }
    f, err := os.Open(filepath)
    if err != nil { return nil, err }
    defer f.Close()
    lines, err := csv.NewReader(f).ReadAll()   // hamısını bir dəfəyə oxu
    return lines, nil
}

// YAZ:
func saveCSVFile(filepath string) error {
    csvfile, err := os.Create(filepath)
    if err != nil { return err }
    defer csvfile.Close()
    csvwriter := csv.NewWriter(csvfile)
    csvwriter.Comma = '\t'               // default ',' dəyişdirilir
    for _, row := range myData {
        temp := []string{row.Name, row.Surname, row.Number, row.LastAccess}
        err = csvwriter.Write(temp)
        if err != nil { return err }
    }
    csvwriter.Flush()                     // buferi diskə boşalt!
    return nil
}
```
**Sub-kod izahı:**
- `ReadAll()` → `[][]string` — sətir × sahə; kiçik fayllar üçün; böyük faylda sətir-sətir
  `Read()` daha yaxşıdır
- `csvwriter.Comma` — sahə ayırıcını dəyişir
- `Flush()` yazma prosesini tamamlayır — unutma = data itir

### 9. Statistika Tətbiqi — CSV Version
**Kitabdan kod nümunəsi:**
```go
func readFile(filepath string) ([]float64, error) {
    f, err := os.Open(filepath)
    if err != nil { return nil, err }
    defer f.Close()
    lines, _ := csv.NewReader(f).ReadAll()
    values := make([]float64, 0)
    for _, line := range lines {
        tmp, err := strconv.ParseFloat(line[0], 64)  // string → float64
        if err != nil {
            log.Println("Error reading:", line[0], err)
            continue                                   // etibarsız sətri ötür
        }
        values = append(values, tmp)
    }
    return values, nil
}

func stdDev(x []float64) (float64, float64) {   // (mean, stdDev) qaytarır
    sum := 0.0
    for _, val := range x { sum += val }
    meanValue := sum / float64(len(x))
    var squared float64
    for i := 0; i < len(x); i++ {
        squared += math.Pow((x[i]-meanValue), 2)
    }
    standardDeviation := math.Sqrt(squared / float64(len(x)))
    return meanValue, standardDeviation
}

// main: oxu → sort.Float64s → min=values[0], max=values[len-1] → stdDev → normalize
```
**Yaxşılıqlar:** funksiyalara bölünmüş main(); diskdən data. **Gələcək inkişaf:**
çoxlu CSV, JSON data, nəticələrin sortlanması.

## Əsas terminlər
- Map (xəritə) — açar-dəyər strukturu, constant time access
- Comparable (müqayisəolunan) — map açarı üçün tələb: == dəstəyi
- Zero Value — olmayan açarın oxunmasında qaytarılan default dəyər
- nil Map — yazılmayan map; assignment → panic
- Structure (struktur) — sahələr qrupu; yeni adlı tip
- Nominal Typing (nominal typing) — myInt və int fərqli tiplər kimi
- Finite Automaton (sonlu avtomat) — regex-in compile nəticəsi
- Lookahead/Lookbehind — Go-da DƏSTƏKLƏNMƏYƏN regex konstruktları
- Raw String Literal (xam sətir literalı) — back quote, escape-siz
- CSV (vergüllə ayrılmış dəyərlər) — encoding/csv paketi
- io.Reader/io.Writer — oxuma/yazma üçün universal interfeyslər

## Praktik nəticə

(1) Slice index tələbindən azad olmaq istəyirsənsə — map; amma sıra lazımdırsa map YOX.
(2) `v, ok := m[k]` — açarın mövcudluğunun yeganə etibarlı yoxlaması. (3) nil map-ə yazma
panic-dir — funksiya girişində yoxla. (4) Map iterasiyası random sıra ilə — sortlanmış
çıxış lazımdırsa açarları ayrıca topla+sortla. (5) Struct zero value = bütün sahələr
default-də; init funksiyası həm validation, həm təhlükəsiz yaradıcıdır. (6) new(T) həmişə
pointer qaytarır, map/channel üçün işləmir. (7) Regex-də MustCompile init vaxtı yoxlanışı
verir; səhv pattern başlanğıcda partlayır — bu faydalıdır. (8) Regex işin yeganə həllidir
deyil — sadə hallarda strconv daha yaxşıdır. (9) CSV oxu/yazı io interfeysləri üzərindən
gedir; numeric data əl ilə parse olunmalı; Flush() yazmanı tamamlayır.

## Mənbə
Pages: 105-130 (PDF 136-163)
