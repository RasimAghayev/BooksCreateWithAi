# Unit 3 — Building Blocks (Lesson 12-15)

## Bu unit nədən bəhs edir?

Funksiya anatomiya (parametr/arqument/nəticə, export, çoxqayıdış, variadic + empty interface), öz tiplərinizi yaratmaq (`type celsius float64`) və metodlar (receiver), first-class funksiyalar (dəyişənə mənimsətmə, parametr kimi, closure/anonymous funksiyalar, funksiya tipləri). Capstone: temperatur cədvəlləri.

**PDF səhifələr:** 108-134 (L12: 108-115, L13: 116-122, L14: 123-131, L15: 132-134)

## Əsas fikirlər

### 1. Funksiya deklarasiyası — anatomiya (L12)
```
func Intn(n int) int
 │    │    │      └ result type
 │    │    └ parameter (ad + tip)
 │    └ funksiya adı (böyük hərf = exported)
 └ açar söz
```

**Vacib anlayışlar:**
- **Parameter** — funksiyanın qəbul etdiyi (elan); **arqument** — çağırışda ötürülən. (`func Contains(s, substr string) bool` → `strings.Contains("walk", "alk")`)
- **Böyük hərf = exported** — digər paketlərdən istifadə oluna bilər; kiçik hərf = yalnız paket daxili
- Eyni tip ard-arda parametrlərdə tip bir dəfə: `func Unix(sec, nsec int64) Time`

### 2. Çoxlu nəticə (L12)
```go
countdown, err := strconv.Atoi("10")
// Deklarasiya:
func Atoi(s string) (i int, err error)   // adlı nəticələr
func Atoi(s string) (int, error)         // adsız — hər ikisi leqal
```
- `error` — daxili tip; L28-də dərinləşir

### 3. Println-in sirri — variadic + empty interface (L12)
```go
func Println(a ...interface{}) (n int, err error)
```
- `...` → **variadic**: dəyişən sayda arqument (a = arqumentlər toplusu)
- `interface{}` → **empty interface**: hər tipi qəbul edir — Println-ə istənilən tip verməyin sirri
- L18-də variadic funksiyaların yazılması göstərilir

### 4. Funksiya yazma + pass by value (L12)
**Kitabdan kod nümunəsi:**
```go
// kelvinToCelsius converts ºK to ºC
func kelvinToCelsius(k float64) float64 {
    k -= 273.15
    return k
}

func main() {
    kelvin := 294.0
    celsius := kelvinToCelsius(kelvin)
    fmt.Print(kelvin, "º K is ", celsius, "º C")
    // 294º K is 20.850000000000023º C — kelvin DƏYİŞMƏDİ!
}
```

**Sub-kod izahı:**
- Şərh konvensiyası: funksiya adı ilə başlayır
- **Pass by value:** `k` parametri `kelvin` arqumentinin KOPYASı — funksiya daxilində `k`-nın dəyişməsi `kelvin`-ə təsir etmir → funksiyalar izolyasiya olunur
- Eyni adlı dəyişənlər müxtəlif funksiyalarda müstəqildir (scope)
- Side-effect-free funksiyalar: yalnız parametr qəbul edir, yalnız nəticə qaytarır — ən asan test/reuse olunanlar

### 5. Yeni tip yaratmaq (L13)
```go
type celsius float64

var temperature celsius = 20     // literal untyped → OK
temperature += 10               // float64 kimi davranır
```

**Alias DEYİL — müstəqil tip:**
```go
var warmUp float64 = 10
temperature += warmUp                    // ERROR: mismatched types
temperature += celsius(warmUp)            // OK — konvertasiya

type fahrenheit float64
var c celsius = 20
var f fahrenheit = 20
if c == f { }        // ERROR — müqayisə də qadağan!
c += f               // ERROR
```
**Fayda:** celsius + fahrenheit qarışdırma səhvi kompilyasiyada yaxalanır — readability + reliability.

### 6. Metodlar (L13)
**Kitabdan kod nümunəsi:**
```go
type kelvin float64

// funksiya versiyası:
func kelvinToCelsius(k kelvin) celsius {
    return celsius(k - 273.15)
}

// metod versiyası — EYNİ daxili məntiq:
func (k kelvin) celsius() celsius {
    return celsius(k - 273.15)
}

type fahrenheit float64
func (f fahrenheit) celsius() celsius {
    return celsius((f - 32.0) * 5.0 / 9.0)
}

func main() {
    var k kelvin = 294.0
    c := k.celsius()          // dot notation
}
```

**Sub-kod izahı:**
- `func (k kelvin) celsius() celsius` → **receiver** `(k kelvin)` adından ƏVVƏL; metodda məcburi 1 receiver — daxildə adi parametr kimi davranır
- Çağırış: `k.celsius()` — variable.metod()
- **Simmetriya faydası:** hər temperatur tipinin öz `celsius()` metodu ola bilər; funksiyalarda `celsius` adı 1 dənə ola bilər, metodlarda hər tipdə 1 dənə!
- Metodlar yalnız **eyni paketdə elan olunan tiplərə** bağlana bilər (int/float64 predeclared-lara YOX)

### 7. First-class funksiyalar — dəyişənə mənimsətmə (L14)
**Kitabdan kod nümunəsi:**
```go
type kelvin float64

func fakeSensor() kelvin {
    return kelvin(rand.Intn(151) + 150)
}

func realSensor() kelvin {
    return 0
}

func main() {
    sensor := fakeSensor     // funksiya ÖZÜ (ÇAĞIRIŞ YOX — qaranlıq mötərizə yoxdur!)
    fmt.Println(sensor())     // fakeSensor() işə düşür
    sensor = realSensor
    fmt.Println(sensor())     // realSensor()
}
```
- `sensor := fakeSensor` → mənimsətmə; `sensor()` → çağırış
- **Signature uyğunluğu:** parametr/nəticə tipləri eyni olmalıdır — `groundSensor() celsius` → mənimsədilməz (error)
- Tip elanı: `var sensor func() kelvin`

### 8. Funksiyanı parametr kimi ötürmə (L14)
```go
func measureTemperature(samples int, s func() kelvin) {
    for i := 0; i < samples; i++ {
        k := s()
        fmt.Printf("%vº K\n", k)
        time.Sleep(time.Second)
    }
}

measureTemperature(3, fakeSensor)    // funksiya adı birbaşa ötürülür
```
→ measureRealTemperature + measureFakeTemperature dublikasiyasından xilas!

### 9. Funksiya tipi elan etmə (L14)
```go
type sensor func() kelvin          // yeni tip!

func measureTemperature(samples int, s sensor) { ... }
```
Uzun imzaları qısaldır; çox parametrli funksiya tiplərində faydalı.

### 10. Anonymous funksiyalar + closure (L14)
**3 istifadə forması:**
```go
// 1. Dəyişənə mənimsətmə:
f := func(message string) {
    fmt.Println(message)
}
f("Go to the party.")

// 2. Elan + dərhal çağırış:
func() {
    fmt.Println("Functions anonymous")
}()
```

**Closure = kalibrasiya nümunəsi:**
```go
type sensor func() kelvin

func calibrate(s sensor, offset kelvin) sensor {
    return func() kelvin {          // anonim funksiya QAYTAR
        return s() + offset          // s və offset-i "Əhatə edir"!
    }
}

func main() {
    sensor := calibrate(realSensor, 5)
    fmt.Println(sensor())     // 5
}
```
- **Closure:** anonim funksiya əhatə edən scope-un dəyişənlərinə **referans** saxlayır — calibrate qayıtdıqdan sonra belə `s` və `offset` yaşayır
- `sensor()` = realSensor() + 5

**Referans kopya deyil:**
```go
var k kelvin = 294.0
sensor := func() kelvin { return k }
fmt.Println(sensor())   // 294
k++                     // XARİCDƏ dəyişiklik
fmt.Println(sensor())   // 295 — closure GÖRÜR!
```
**Təhlükə:** for loop daxilində closure-lar — dəyişənlər canlı qalır (Ch 2-dəki goroutine tələsini xatırla).

### 11. Capstone: Temperature tables (L15)
```
=======================
| ºC       | ºF       |
=======================
| -40.0    | -40.0    |
| ...      | ...      |
=======================
```
- Cədvəl 1: ºC → ºF (-40-dan +100-ə, 5° addım)
- Cədvəl 2: ºF → ºC (tərsinə)
- `drawTable(rows int, getRow func(row int) (string, string))` — sətir datasını **first-class funksiya** ilə al → eyni drawTable müxtəlif data çəkir

## Əsas terminlər
- Function Declaration / Parameter / Argument
- Exported (böyük hərf) / Unexported
- Multiple Return Values
- Named Result (`(i int, err error)`)
- Variadic Function (`...`)
- Empty Interface (`interface{}`)
- Pass by Value (dəyərlə ötürmə)
- Side-effect-free Function
- `type` Keyword (yeni tip)
- Underlying Type (alt tip)
- Type Alias vs New Type (fərq!)
- Method / Receiver
- Dot Notation
- First-class Function
- Function Signature
- Function Type (`type sensor func() kelvin`)
- Anonymous Function / Function Literal
- Closure (dəyişənlərə referans saxlama)

## Praktik nəticə
- `celsius` + `fahrenheit` kimi tiplər yarat — ölçü vahidi qarışdırma səhvlərini kompilyasiyaya ver; bu, klassik "primitive obsession" antipatterninin dərmanıdır.
- Metod adlarını qısa saxla (`k.celsius()` — receiver-in tipi kontekst verir); funksiya adlarındakı `kelvinTo...` prefiksi lazmsızlaşır.
- Funksiyanı dəyişənə mənimsədəndə mötərizəsiz yaz — `sensor := fakeSensor`; çağırış `sensor()`.
- Funksiya tipləri elan et (`type sensor func() kelvin`) — imzalar təkrarlananda kod qısaldır.
- Closure-lar referans saxlayır — loop dəyişənlərini closure-a bağlama; kalibrasiya/factory pattern-lərin əsası budur.
- Pass by value sayəsində funksiyalar izolyasiya olunur — side-effect-free funksiyalara can at.

## Mənbə
Pages: 108-134 (PDF), book pages 93-119 (Lesson 12-15)
