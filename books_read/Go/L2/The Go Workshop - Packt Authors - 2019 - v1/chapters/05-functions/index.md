# Chapter 5 — Functions (Funksiyalar)

## Bu fəsil nədən bəhs edir?

Funksiya anatomiyası (func, identifier, parameter list, return types, body, signature),
dizayn prinsipləri (single responsibility, ~25 sətir), parameter/argument fərqi, shorthand
notation, çoxlu qaytarma, blank identifier, naked returns və onun shadowing təhlükəsi,
variadic funksiyalar (pack/unpack), anonim funksiyalar, closures, funksiya tipləri
(parametr/qaytarma kimi), defer (FILO, dəyər bağlanması) və salary/pay hesablama
activity-ləri.

## Əsas fikirlər

### 1. Funksiyanın Anatomiyası
```go
func calculateTax(amount float64, rate float64) (float64, error) {
    // func | identifier | parameter list | return types | body
}
```
- **camelCase** adlar; ilk hərf kiçik = private (paketdən kənara çıxmır), BÖYÜK = export
- **Parameter vs Argument:** TƏYİN zamanı parameter; ÇAĞIRIŞ zamanı argument
- **Shorthand:** `(firstName, lastName string)` — ardıcıl eyni tip bir dəfə
- **Signature** = parametrlər + qaytarma tipləri — dəyişməməyə çalış (istifadəçiləri
  asılıdır)
- `{` funksiya bəyannaməsi ilə EYNİ sətirdə — kompulsiv qayda

### 2. Dizayn Prinsipləri
- **Single responsibility:** 1 funksiya = 1 tapşırıq (məsafə hesabla VƏ ya vaxt təxmin
  et — İKİ funksiya)
- **~25 sətir həddi:** böyük funksiya = refaktorinq siqnalı; testi çətinləşdirir
- **Parametr şişkinliyi:** 7-8 parametr = çoxşaxılılıq qoxusu — refaktor et

### 3. Çoxlu Qaytarma — Go-nun Fərqi
```go
func fizzBuzz(i int) (int, string) {
    switch {
    case i%15 == 0: return i, "FizzBuzz"
    case i%3 == 0:  return i, "Fizz"
    case i%5 == 0:  return i, "Buzz"
    }
    return i, ""
}
n, s := fizzBuzz(i)     // sıra mütləq uyğun
_, s := fizzBuzz(i)     // _ = blank identifier — int-i İGNORLA
```
**İdiomatik qayda:** 2-ci qaytarma tez-tez `error` tipidir; error-u İGNOR ETMƏ.

### 4. Naked Returns — Adlı Qaytarma
```go
func greeting() (name string, age int) {   // adlı return = LOKAL dəyişənlər
    name = "John"
    age = 21
    return                                  // NAKED — adlıları qaytarır
}
```
**TƏHLÜKƏ — shadowing:**
```go
func message() (message string, err error) {
    message = "hi"
    if message == "hi" {
        err := fmt.Errorf("say bye\n")   // SHADOW! blok-lokal; return err = NIL
        return
    }
    return
}
```
**Tövsiyə:** naked return-lardan QAÇIN — oxunaqlığı və shadowing riskini artırır;
anlaşılmazlıq yaradır.

### 5. Variadic Funksiyalar
```go
func sum(nums ...int) int {          // pack operatoru ...T
    total := 0
    for _, num := range nums {       // DAXİLDƏ SLİCE-dır!
        total += num
    }
    return total
}
sum(5, 4)          // çoxsaylı arqument
sum()              // sıfır arqument
sum(i...)          // SLİCE ötürülərkən UNPACK ... MÜTLƏQ
// sum(i)          // XƏTA — slice birbaşa variadic-ə verilə BİLMƏZ
```
**Qaydalar:** ...T parametrlərin SONUNDA; funksiyada TƏK variadic; daxildə tipi = `[]T`.

### 6. Anonim Funksiyalar
```go
func() {                              // ad YOX
    fmt.Println("Greeting")
}()                                   // İCRA mötərizələri — dərhal çağırış

func(str string) {                     // parametrli
    fmt.Println(str)
}(message)                             // arqument icra mötərizəsində

f := func() { fmt.Println("later") }   // DƏYİŞƏNDƏ — func() tipi
f()                                    // sonradan çağır
```
**İstifadə:** closure-lar, defer, goroutine blokları, birdefəlik kod.

### 7. Closures — Xarici Dəyişən Bağlanması
**Sadə hal (qorumasız):**
```go
i := 0
incrementor := func() int {
    i += 1        // main-in i-sini GÖRÜR və DƏYİŞİR
    return i
}
incrementor()  // 1
incrementor()  // 2
i += 10        // problem: hər kəs i-yə çata bilir → 12
```

**Qorunmuş closure (fabrika):**
```go
func incrementor() func() int {
    i := 0                        // YALNIZ bu funksiyanın scope-unda
    return func() int {
        i += 1                    // i ölçüsünə bağlı — xaricdən ÇATILMIR
        return i
    }
}
increment := incrementor()        // incrementor BİR dəfə icra
increment()  // 1                // hər çağırış anonim funksiyanı işə salır
increment()  // 2                 // i YAŞAYIR (closure saxlayır)
```
**Parametrli closure (kitabdan):**
```go
func decrement(i int) func() int {
    return func() int { i--; return i }   // parametr də bağlanır
}
x := decrement(4); x()  // 3, 2, 1, 0
```

### 8. Funksiya Tipləri — First-Class Citizens
```go
type calc func(int, int) string        // ÖZ funksiya tipin!
func add(i, j int) string { ... }        // eyni signature = calc-a uyğun

func calculator(f calc, i, j int) {       // funksiya PARAMETR kimi
    fmt.Println(f(i, j))
}
calculator(add, 5, 6)

// Alternativ — adlı tip olmadan:
func calculator(f func(int, int) int, i, j int) { }
calculator(add, 5, 6)
calculator(subtract, 10, 5)             // signature uyğun olduqca hər funksiya
```
**Funksiya QAYTARMA:**
```go
func square(x int) func() int {
    return func() int { return x * x }    // closure da eyni vaxtda
}
v := square(9)   // Type of v: func() int
v()              // 81
```
**Salary strategy (kitabdan):**
```go
func salary(x, y int, f func(int, int) int) int {
    pay := f(x, y)
    return pay
}
func developerSalary(hourlyRate, hoursWorked int) int { return hourlyRate * hoursWorked }
func managerSalary(baseSalary, bonus int) int { return baseSalary + bonus }

salary(50, 2080, developerSalary)     // 104000
salary(150000, 25000, managerSalary) // 175000
// YENİ rol = YENİ uyğun funksiya; salary funksiyası DƏYİŞMİR — açıq扩展
```

### 9. defer — Təxirə Salınmış İcra
**Nədir:** funksiyanı əhatə edən funksiya QAYTARANA qədər gecikdirir. Təmizlik üçün
ideal: fayl bağlama, DB bağlantısı, temp silmə, panic-dən qurtarma.

```go
defer done()                       // SONDA icra olunacaq
defer func() { ... }()             // anonim variant
```

**FILO (LIFO) sırası:**
```go
defer fmt.Println("birinci")   // ƏN SON icra
defer fmt.Println("ikinci")
defer fmt.Println("üçüncü")    // İLK icra
// Main: Start → Main: End → üçüncü → ikinci → birinci
```

**Dəyər bağlanması tələsi (vacib!):**
```go
age := 25
defer personAge(name, age)    // 25 BAĞLANIR — bu anda kopyalanır
age *= 2                      // 50 olsa da...
personAge → "John is 25."      // defer 25-i GÖRÜR!
```
Dəyişənin defer anındakı dəyəri istifadə olunur; sonrakı dəyişikliklər ƏKS OLUNMUR.

## Əsas terminlələr
- Parameter/Argument — təyin edilən / ötürülən
- Shorthand Parameter Notation — `(a, b T)`
- Function Signature — parametrlər + qaytarma; sabit saxla
- Named Return — return siyahısında ad = lokal dəyişən
- Naked Return — dəyərsiz `return`; shadowing riski
- Blank Identifier `_` — qaytarılan dəyəri ignore et
- Pack `...T` — variadic parametr; daxildə slice
- Unpack `slice...` — slice-i variadic arqumentə aç
- Anonymous Function — adsız funksiya literalı; dəyişənə mənimsədilir
- Closure — xarici scope dəyişənlərini bağlayan anonim funksiya
- Function Type — `type calc func(...)`; imza uyğunluğu
- Higher-Order Function — funksiya qəbul edən funksiya
- defer — gecikdirilmiş icra; FILO sırası
- Value Binding — defer anında arqumentlərin kopyalanması

## Praktik nətidə

(1) 1 funksiya = 1 məsuliyyət, ~25 sətir; böyüyəndə refaktor et. (2) Parametr siyahısı
şişibsə — funksiya çoxşaxılıdır, böl. (3) Çoxlu qaytarma Go-nun gücüdür; error sonuncu;
blank identifier lazımsız dəyəri udur. (4) Naked return-dan qaç — shadowing gizli
nil qaytara bilir. (5) Variadic = son parametr, tək; daxildə slice; slice ötürmək üçün
`...`. (6) Closure dəyişəni qorumaq üçün onu fabrika funksiyanın İÇİNDƏ yarat — xaricdən
çatılmaz. (7) Funksiya = tip: strategy pattern-i imzas uyğun funksiyalarla həll et —
salary nümunəsi. (8) defer: resurs açılışının yanında yaz; FILO sırasını planlaşdır.
(9) defer arqumentləri defer ANINDA bağlanır — sonrakı dəyişikliklər görünmür. (10)
Exported funksiya = Böyük hərf; camelCase private.

## Mənbə
Pages: 169-214 (PDF 202-249)
