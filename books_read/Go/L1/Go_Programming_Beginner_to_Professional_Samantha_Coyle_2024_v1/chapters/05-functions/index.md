# Chapter 5 — Functions – Reduce, Reuse, and Recycle (səh. 192-237)

## Bu fəsil nədən bəhs edir?

Funksiya anatomiyası (func, identifier, parameter list, return, body),
parametrlər/argumentlər, scope, return value-lar (çoxdəyərli, naked
returns, blank identifier), variadic funksiyalar (pack/unpack), anonymous
funksiyalar, closure-lar, function type-lar (funksiya arqument/qaytarma
kimi), defer (FILO, value freeze) və kodun fayl/qovluqlara ayrılması.

## Əsas fikirlər

### 1. Funksiya anatomiyası
```go
func checkNumbers() {           // adlı funksiya
    for i := 1; i <= 30; i++ {
        if i%2 == 0 {
            fmt.Println("Even")
        } else {
            fmt.Println("Odd")
        }
    }
}

func main() {
    fmt.Println("Main is in control")
    checkNumbers()               // ÇAĞIRIŞ — () mütləqdir
    fmt.Println("Back to main")
}
```

**Hissələr:**
- `func` — elan açar sözü; `{` HƏMİŞƏ eyni sətirdə (Go məcburiyyəti)
- Identifier — camelCase; **kiçik hərf = private (paketdaxili); böyük
  hərf = exported** (digər paketlərdən çağırıla bilər)
- Parameter list — `(name string, age int)`; eyni tiplər üçün shorthand:
  `(firstName, lastName string)`
- Return types — bir və ya BİR NEÇƏ (Go-nun fərqi!); adətən 2-ci = error
- Body — `{}` arası; return olan funksiya SON ifadəsi return olmalıdır

**Dizayn prinsipləri:** single responsibility (bir funksiya = bir
vəzifə), ~25 sətir həddi, təkrarlanan kod = funksiyalaşdırma siqnalı.

### 2. Parametrlər vs Argumentlər
```go
checkNumbers(10)                 // 10 = ARGUMENT (çağırışda)
func checkNumbers(end int)       // end = PARAMETR (tərifdə)

// Şorthand:
func checkNumbers(start, end int) { ... }
checkNumbers(start, end)
```
- Sıra TİPLƏRƏ uyğun olmalıdır: `greeting(45, "Cayden")` xətadir
- Parametrlər funksiya-lokal dəyişənlərdir
- **Scope:** main-in `m`-i greeting-də görünmür; greeting-in `s`-i
  main-də görünmür — yalnız parametr keçidi ilə

### 3. Return value-lar
```go
// Çoxdəyərli return:
func checkNumbers(i int) (int, string) {
    switch {
    case i%2 == 0:
        return i, "Even"
    default:
        return i, "Odd"
    }
}

num, result := checkNumbers(4)   // sıra ilə təyinat
_, result := checkNumbers(4)     // _ = blank identifier (int ignore)
```
- Fayl oxunuşunda: `_, err := file.Read(bytes)` — sayğacı ignore edirik
- **Error-u ignore etmək tövsiyə olunmur**

**Named returns + naked return:**
```go
func greeting() (name string, age int) {   // ADLI return-lar
    name = "John"                          // lokal dəyişən kimi
    age = 21
    return                                 // NAKED — avtomatik qayıdır
}
```
- Naked return oxunuqluluğu azaldır + shadowing təhlükəsi:
```go
func message() (message string, err error) {
    message = "hi"
    if message == "hi" {
        err := fmt.Errorf("say bye\n")   // := — YENİ err (shadow)!
        return                            // qaytarılan err = nil
    }
    return
}
```

### 4. Variadic funksiyalar
```go
func nums(i ...int) {          // ... = PACK operatoru
    fmt.Println(i)             // daxildə SLİCE-dir!
}

nums(99, 100)                  // [99 100]  — arqumentlər
nums(200)                       // [200]
nums()                          // []

// Variadic SON parametr olmalıdır:
func nums(str string, i ...int) { ... }   // ✓
func nums(i ...int, str string) { ... }   // COMPILE XƏTASI

// Slice ötürmək — UNPACK (...):
i := []int{5, 10, 15}
nums(i...)                     // slice → arqumentlərə açılır
nums(i)                        // XƏTA — []int != ...int

// Cəm nümunəsi:
func sum(nums ...int) int {
    total := 0
    for _, num := range nums {
        total += num
    }
    return total
}
sum(5, 4)         // 9
sum(i...)         // 30
```

### 5. Anonymous funksiyalar
```go
// Dərhal icra (execution parentheses):
func() {
    fmt.Println("Greeting")
}()

// Arqumentlə:
message := "Greeting"
func(str string) {
    fmt.Println(str)
}(message)

// Dəyişənə mənimsət (sonra çağır):
f := func() {
    fmt.Println("Executing an anonymous function using a variable")
}
f()                             // dəyişən = func() tipi

// Parametrli + qaytaran:
x := func(i int) int {
    return i * i
}
fmt.Printf("The square of %d is %d\n", 9, x(9))   // 81
```
İstifadə: closure, defer, goroutine, bir dəfəlik bloklar.

### 6. Closure-lar
**Nədir:** xarici dəyişənləri YADINDA SAXLAYAN anonymous funksiya.

```go
// Problematik versiya — i HƏR KƏSƏ açıq:
i := 0
incrementor := func() int {
    i += 1
    return i
}
incrementor()   // 1
incrementor()   // 2
i += 10         // xaricdən dəyişmə MÜMKÜNDÜR — problem!
incrementor()   // 12 (3 deyil!)

// QORUNMUŞ versiya — i funksiya daxilində:
func incrementor() func() int {
    i := 0                        // yalnız buradan əlçatandır
    return func() int {
        i += 1
        return i
    }
}

increment := incrementor()         // BİR dəfə icra
increment()                       // 1
increment()                       // 2
// i artıq xaricdən DƏYİŞMƏZ — inkapsulyasiya!

// Parametrli closure:
func decrement(i int) func() int {
    return func() int {
        i--
        return i
    }
}
x := decrement(4)   // 3, 2, 1, 0
```

### 7. Function type-lar
```go
// Öz tipin yarat:
type calc func(int, int) string

func add(i, j int) string {
    return fmt.Sprintf("Added %d + %d = %d", i, j, i+j)
}

func calculator(f calc, i, j int) {   // funksiya = parametr
    fmt.Println(f(i, j))
}
calculator(add, 5, 6)                // add calc-ın imzasına uyğundur

// İnline imza (tipsiz):
func calculator(f func(int, int) int, i, j int) {
    fmt.Println(f(i, j))
}
calculator(add, 5, 6)       // 11
calculator(subtract, 10, 5) // 5

// Funksiya QAYTARAN funksiya:
func square(x int) func() int {
    f := func() int {
        return x * x
    }
    return f
}
v := square(9)
fmt.Println(v())               // 81
fmt.Printf("Type of v: %T", v) // func() int
```

**Maaş strategiyası nümunəsi:**
```go
func salary(x, y int, f func(int, int) int) int {
    pay := f(x, y)
    return pay
}
func managerSalary(baseSalary, bonus int) int { return baseSalary + bonus }
func developerSalary(hourlyRate, hoursWorked int) int { return hourlyRate * hoursWorked }

devSalary := salary(50, 2080, developerSalary)       // 104000
bossSalary := salary(150000, 25000, managerSalary)   // 175000
```
- Yeni vəzifə = yeni eyni-imzalı funksiya; salary DƏYİŞMİR — genişlənən
  dizayn (strategy pattern funksiyalarla)

### 8. defer
**Nədir:** əhatə funksiyası QAYITMAMAZDAN ƏVVƏL icra olunan çağırış.

```go
defer done()
fmt.Println("Main: Start")     // 1-ci çap
fmt.Println("Main: End")       // 2-ci çap
// done() ƏN SONDA — "Now I am done"
```

**Çoxlu defer = FILO (First In, Last Out):**
```go
defer fmt.Println("I was declared first.")
defer fmt.Println("I was declared second.")
defer fmt.Println("I was declared third.")
// Çap sırası: third → second → first (tərs)
```

**Dəyər dondurma (freeze):**
```go
age := 25
name := "John"
defer personAge(name, age)     // ARQUMENTLƏR İNDİ hesablanır!
age *= 2
fmt.Printf("Age double %d.\n", age)   // 50
// defer işə düşəndə: "John is 25." — 25! 50 YOX
```
- İstifadə: resurs təmizliyi (fayl/bağlantı bağlama), panic recover
  (sonrakı fəsil)

### 9. Kodun ayrılması
- Funksiya → fayl (salary.go, game.go, weather.go) → qovluq → modul
- Faydalar: reusability, readability, testability
- Ətraflı modul fəslində (ch9-10)

## Activity icmalları
- **5.01:** Employee/Developer struct + LogHours/HoursWorked metodları —
  funksiyaların vəzifə ayrılığı
- **5.02:** nonLoggedHours() closure + PayDay (overtime 2x, >40h) +
  PayDetails — closure/çoxreturn/strategy birləşməsi

## Əsas terminlər
- First-class / higher-order functions
- camelCase / exported (böyük hərf) vs private
- Parameter (tərif) vs argument (çağırış)
- Shorthand parameter notation — `(a, b string)`
- Blank identifier (_) — qaytarma ignore
- Named returns / naked return + shadowing riski
- Variadic — `...T` (pack); `slice...` (unpack)
- Anonymous function / function literal
- Closure — xarici dəyişənlərin yaddaş saxlanması
- Function type — `type calc func(int,int) string`
- defer + FILO + arqument freeze
- Single responsibility (~25 sətir)

## Praktik nəticə
Funksiya = tip: dəyişənə mənimsət, parametr keç, qaytar. Strategiya
dəyişənliyi üçün `func(...)` imzalı parametrlər (salary nümunəsi) —
aç-açiq genişlənmə. Dövlətli inkapsulyasiya üçün closure (i yalnız
incrementordan dəyişilir). Çoxdəyərli return + `_` ilə seçici qəbul;
naked return-dan qaçın (shadowing!). Variadic: `...T` qəbul, `s...`
keçir. defer ilə təmizlik — amma unutma: arqumentlər defer VAXTINDA
donur; sıra FILO.

## Mənbə
Pages: 192-237 (PDF 192-237)
