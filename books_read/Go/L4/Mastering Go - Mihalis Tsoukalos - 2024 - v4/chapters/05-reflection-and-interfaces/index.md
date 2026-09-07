# Chapter 5 — Reflection and Interfaces (Refleksiya və İnterfeyslər)

## Bu chapter nədən bəhs edir?

reflect paketi (ValueOf/TypeOf/Kind, struct sahələrinin kəşfi və dəyişdirilməsi), refleksiyanın
3 çatışmazlığı, type methods (value vs pointer receiver), interfeyslər (implicit satisfaction,
sort.Interface, empty interface, type assertion/switch, map[string]interface{}), error interfeysi,
OOP mimikriyası, interfaces vs generics vs reflection müqayisəsi və statistika tətbiqinin
çoxlu fayl + sort.Interface versiyası.

## Əsas fikirlər

### 1. Reflection — Nədir və Nə Vaxt
**Nədir:** Runtime-da ixtiyari obyektin tipini və strukturunu dinamik öyrənmək. reflect
paketi: `reflect.Value` (dəyər), `reflect.Type` (tip), `reflect.Kind` (xüsusi tip növü).

**Nəyə lazımdır:**
- Yazılanda mövcud olmayan, gələcəkdə gələcək tiplərlə iş (JSON decode kimi)
- Ümumi interfeysi olmayan qeyri-adi tiplərlə iş
- fmt.Println arxasında da reflection işləyir

**Kitabdan kod nümunəsi (struct kəşfi):**
```go
A := Record{"String value", -12.123, Secret{"Mihalis", "Tsoukalos"}}
r := reflect.ValueOf(A)
iType := r.Type()

fmt.Printf("The %d fields of %s are\n", r.NumField(), iType)
for i := 0; i < r.NumField(); i++ {
    fmt.Printf("\t%s ", iType.Field(i).Name)               // sahə adı
    fmt.Printf("\twith type: %s ", r.Field(i).Type())      // sahə tipi
    fmt.Printf("\tand value _%v_\n", r.Field(i).Interface()) // dəyər

    k := reflect.TypeOf(r.Field(i).Interface()).Kind()
    if k == reflect.Struct {          // Kind yoxlaması — embedded struct tapır
        fmt.Println(r.Field(i).Type())
    }
}
```
**Sub-kod izahı:** `main.Record` — pakat adı + struct adı unikal identifikator;
`NumField()`/`Field(i)` — struct sahələri üzrə iterasiya.

**Sahə dəyərini dəyişmək (pointer + Elem mütləq):**
```go
r := reflect.ValueOf(&A).Elem()    // &A — dəyişmək üçün pointer
for i := 0; i < r.NumField(); i++ {
    k := reflect.TypeOf(r.Field(i).Interface()).Kind()
    if k == reflect.Int {
        r.Field(i).SetInt(-100)         // int sahəni dəyiş
    } else if k == reflect.String {
        r.Field(i).SetString("Changed!") // string sahəni dəyiş
    }
}
```

**Refleksiyanın 3 çatışmazlığı:**
1. Oxunması çətin kod — maintenance problemləri
2. Yavaş icra — statik tipə görə kod həmişə daha sürətli
3. Xətalar compile-deyil RUNTIME-da panic kimi — aylar sonra belə partlaya bilər

**Sahə olduğu hallar:** JSON/XML serializasiya, dynamic code generation, struct→DB table
mapping.

### 2. Type Methods — Tipə Bağlı Funksiyalar
**Nədir:** `func (a ar2x2) MetodAdı(...)` — funksiyanı tipə bağlayan receiver.
Compiler metodları arxada `MetodAdı(a, ...)` funksiya çağırışına çevirir.

**Kitabdan kod nümunəsi:**
```go
type ar2x2 [2][2]int

// Ənənəvi funksiya — nəticəni QAYTARIR:
func Add(a, b ar2x2) ar2x2 {
    c := ar2x2{}
    for i := 0; i < 2; i++ {
        for j := 0; j < 2; j++ {
            c[i][j] = a[i][j] + b[i][j]
        }
    }
    return c
}

// Type method — nəticəni ÇAĞIRANA yazır (pointer receiver):
func (a *ar2x2) Add(b ar2x2) {
    for i := 0; i < 2; i++ {
        for j := 0; j < 2; j++ {
            a[i][j] = a[i][j] + b[i][j]
        }
    }
}

a.Add(b)     // nəticə a-da; Add(a,b) isə yeni array qaytarır
```

**Value vs Pointer receiver seçimi:**
- **Value receiver:** dəyişiklik etmir; kiçik/immutable tiplər; məntiqən dəyərə aiddir
- **Pointer receiver:** receiver-in vəziyyətini DƏYİŞİR; böyük struct (kopya xərci);
  məntiqən konkret instansiyaya aiddir
- Metod adı struct sahə adı ilə TOQQUŞMAMALI — compiler rədd edir
- Sabit ölçülü, performance-kritik data üçün array (slice yox) — 2x2 matris nümunəsi

### 3. İnterfeyslər — Davranış Müqaviləsi
**Nədir:** Metod dəsti təyin edən abstrakt tip. **Implicit satisfaction:** tip interfeysin
BÜTÜN metodlarını implement edirsə, interfeysi avtomatik satisfies edir — "implements"
bəyannaməsi YOXDUR (duck typing).

**Dizayn prinsipləri (kitabdan):**
- Proqramı interfeyslərdən başlayaraq dizayn etmə — əvvəl kod, sonra ümumi davranışlar
  üzə çıxanda interfeysə çevir
- Kiçik və dəqiq interfeyslər ən populyardır (1 metod = ideal)
- 2+ konkret tip arasında davranış paylaşılırsa yarat; sadələşdirmirsə — sil

**Kompozisiya:** interfeyslər birləşə bilər — `type IntC interface { IntA; IntB }`.

### 4. sort.Interface — Xüsusi Sıralama
**Kitabdan kod nümunəsi:**
```go
type Personslice []Person

func (a Personslice) Len() int { return len(a) }

// Hansı sahə üzrə müqayisə — burada qərar verilir:
func (a Personslice) Less(i, j int) bool {
    return a[i].F3.F1 < a[j].F3.F1      // embedded struct sahəsi belə mümkün
}

func (a Personslice) Swap(i, j int) {
    a[i], a[j] = a[j], a[i]
}

sort.Sort(Personslice(data))                    // artan
sort.Sort(sort.Reverse(Personslice(data)))      // azalan — avtomatik işləyir
```
**Sub-kod izahı:** 3 metod (Len/Less/Swap) → sort.Interface satisfied; Less müqayisə
məntiqini təyin etdiyindən istənilən meyar (məs. mean dəyəri) mümkündür.

### 5. Empty Interface və Type Assertion/Switch
**Empty interface** (`interface{}` = `any`): 0 metod → BÜTÜN tiplər tərəfindən satisfied.

**Type assertion — `x.(T)`:**
```go
anInt := returnNumber()              // interface{} daxilində int

number, ok := anInt.(int)            // təhlükəsiz forma — ok yoxlaması
if ok { /* istifadə et */ }

i := anInt.(int)                    // tək dəyər — uğursuzsa PANIC
_ = anInt.(bool)                    // PANIC: interface conversion: int, not bool
```

**Type switch:**
```go
switch T := x.(type) {
case Secret:
    fmt.Println("Secret type")
case Entry:
    fmt.Println("Entry type")
default:
    fmt.Printf("Not supported type: %T\n", T)
}
```
- Case sırası vacibdir — yalnız İLK uyğun gələn icra olunur; spesifikdən generikə düz

### 6. map[string]interface{} — Naməlum JSON
**Nədir:** Dəyərləri istənilən tipi daşıyan map — JSON-un schema-sı əvvəlcədən
bilinmədiyi hallar üçün ideal. map[string]string-dən fərqi: orijinal tip MÜQAVİLƏTƏ
QORUNUR.

**Kitabdan kod nümunəsi:**
```go
JSONMap := make(map[string]interface{})
err := json.Unmarshal([]byte(JSONrecord), &JSONMap)

func typeSwitch(m map[string]interface{}) {       // rekursiv kəşfiyyət
    for k, v := range m {
        switch c := v.(type) {
        case string:
            fmt.Println("Is a string!", k, c)
        case float64:                              // JSON rəqəmlər = float64!
            fmt.Println("Is a float64!", k, c)
        case bool:
            fmt.Println("Is a Boolean!", k, c)
        case map[string]interface{}:
            typeSwitch(v.(map[string]interface{}))   // REKURSİYA — nested obyekt
        default:
            fmt.Printf("...Is %v: %T!\n", k, c)
        }
    }
}
```

### 7. error — İnterfeysdir
```go
type error interface {
    Error() string
}
```
**Custom error nə vaxt:** xətaya daha çox kontekst vermək istəyirsən.

**Kitabdan kod nümunəsi (boş fayl fərqləndirməsi):**
```go
type emptyFile struct {
    Ended bool
    Read  int
}

func (e emptyFile) Error() string {          // error interfeysi satisfied
    return fmt.Sprintf("Ended with io.EOF (%t) but read (%d) bytes", e.Ended, e.Read)
}

func isFileEmpty(e error) bool {
    v, ok := e.(emptyFile)                   // type assertion error-dan struct çıxarır
    if ok {
        if v.Read == 0 && v.Ended == true {
            return true
        }
    }
    return false
}
// readFile: io.EOF + n==0 → return emptyFile{true, n}  → kontekstli xəta
```
**io.EOF ciddi xəta deyil — fayl oxunun məntiqi hissəsidir.**

### 8. Öz İnterfeysin + Satisfaction Yoxlaması
```go
type Shape2D interface {
    Perimeter() float64
}
type circle struct{ R float64 }
func (c circle) Perimeter() float64 { return 2 * math.Pi * c.R }

// Tipin interfeysi satisfy edib-etmədiyini yoxla:
_, ok := interface{}(a).(Shape2D)
if ok { fmt.Println("a is a Shape2D!") }
```

### 9. Go-da OOP Mimikriyası
Go tam OOP deyil, amma təqlid edir:
- **Struct + type methods** = obyekt + metodları
- **İnterfeyslər** = abstrakt tip / polimorfizm
- **Kiçik hərf** = private (encapsulation — package daxilində)
- **İnterfeys + struct birləşməsi** = composition

**Embedding (anonim sahə):**
```go
type compose struct {
    field1 int
    a                  // anonim — sahələri birbaşa çatılır
}
iComp := compose{123, a{456, 789}}
iComp.XX, iComp.YY     // a-nın sahələri birbaşa
```
Miras DEYİL — sahə/metod promote-un-dır; fərqli tiplər eyni adlı metod daşıya bilər.

### 10. İnterfeys vs Generics vs Reflection
| Meyar | İnterfeys | Generics | Reflection |
|---|---|---|---|
| Tip təhlükəsizliyi | runtime (assertion) | COMPILE vaxtı | runtime (kind) |
| Tətbiq sahəsi | konkret DAVRANIŞ (Reader/Writer) | çoxlu TİP eyni əməliyyat | struktur KƏŞFİ (marshal) |
| Kod sadəliyi | orta | ən sadə | ən mürəkkəb |
| Runtime xərci | type switch | orta | ən yüksək |

**Kitabın qərarı:** çoxlu tip üçün eyni əməliyyat → generics; konkret davranış müqaviləsi →
interfeys; tipin daxili strukturunu bilmədən kəşf → reflection.

**Reflection variantı (kitabdan):**
```go
func PrintReflection(s interface{}) {
    val := reflect.ValueOf(s)
    if val.Kind() != reflect.Slice {   // manual tip yoxlaması MÜTLƏQ
        return
    }
    for i := 0; i < val.Len(); i++ {
        fmt.Print(val.Index(i).Interface(), " ")
    }
}
```
Generics eynini `func PrintSlice[T any](s []T)` ilə 3 sətirdə edir — compiler yoxlaması bonus.

### 11. Statistika Tətbiqi — Çoxlu Fayl + Sort
**Kitabdan kod nümunəsi:**
```go
type DataFile struct {
    Filename string
    Len      int
    Minimum  float64
    Maximum  float64
    Mean     float64
    StdDev   float64
}
type DFslice []DataFile

func (a DFslice) Less(i, j int) bool {
    return a[i].Mean < a[j].Mean      // MEAN üzrə sırala
}

// Hər fayl üçün:
currentFile.Minimum = slices.Min(values)   // sort etmədən min/max!
currentFile.Maximum = slices.Max(values)
meanValue, standardDeviation := stdDev(values)
files = append(files, currentFile)

sort.Sort(files)      // DataFile-lar mean üzrə sıralanır
```

## Əsas terminlər
- Reflection — runtime tip/dəyər kəşfi (reflect paketi)
- reflect.Value / reflect.Type / reflect.Kind — dəyər / tip / tip növü
- Type Method — tipə bağlı funksiya; receiver `(a *T)` / `(a T)`
- Receiver (qəbul edici) — metodun bağlı olduğu dəyişən
- Value/Pointer Receiver — kopya üzərində / birbaşa dəyər üzərində iş
- Interface — metod dəsti; implicit satisfaction (duck typing)
- Empty Interface — interface{} = any; hər tipi qəbul edir
- Type Assertion — x.(T) ilə underlying dəyərin çıxarılması
- Type Switch — tip üzrə switch; ilk uyğun case icra olunur
- io.EOF — axın sonu; xəta deyil, gözlənilən siqnal
- Embedding (daxil etmə) — anonim sahə; sahə/metod promote
- Polymorphism (çoxşekillilik) — ümumi davranışın fərqli tiplərdə implementasiyası

## Praktik nətidə

(1) İnterfeys = DAVRANIŞ müqaviləsi — "nədir" yox "nə edir". (2) Əvvəl konkret kod,
sonra ümumi davranış → interfeys; 1 metodlu interfeys ən gücludür. (3) Reflection-ı yalnız
məcburi hallarda (marshal/unmarshal, dinamik strukturlar) — yavaş, oxunmaz, runtime panic.
(4) Pointer receiver = dəyişən funksiya və böyük struct; value receiver = saf hesab. (5)
sort.Interface: Len/Less/Swap üçlüyü istənilən meyar üzrə sıralama verir. (6) Type assertion
həmişə `, ok` formasında — tək dəyər forması uğursuzda panic. (7) map[string]interface{}
JSON-un naməlum schema-sı üçün; rəqəmlər float64 kimi gəlir. (8) Custom error = kontekst
əlavəsi; io.EOF xəta deyil. (9) Çoxlu tip üçün eyni əməliyyat — generics; konkret davranış —
interfeys. (10) Go OOP deyil — composition + implicit interfaces onun OOP cavabıdır.

## Mənbə
Pages: 153-198 (PDF 184-231)
