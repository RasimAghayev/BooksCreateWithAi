# Chapter 6 — Errors (Xətalar)

## Bu fəsil nədən bəhs edir?

Xəta növləri (syntax/runtime/semantic), exception handling ilə müqayisə, error
interfeysi (errors.New, errorString, Err adlandırmasi), custom xəta dəyərləri, panic
(mexanizmi, defer ilə interaksiyası, os.Exit-dən fərqi), recover (defer daxilində,
xəta tipinə görə fərqləndirmə) və bank direct-deposit activity zənciri.

## Əsas fikirlər

### 1. Xətaların 3 Növü
| Növ | Səbəb | Tapma |
|---|---|---|
| **Syntax** | dil qaydalarının pozulması (fmt.println, } itmiş) | ƏN ASAN — IDE/compile |
| **Runtime** | icra zamanı mümkünsüz əməliyyat (index out of range, olmayan fayl, /0) | crash/panic |
| **Semantic (logic)** | kod İŞLƏYİR amma NƏTİCƏ SƏHV (>= vs >) | ƏN ÇƏTİN — heç nə xəta vermir |

**Real nümunələr (kitabdan):** Northeast Blackout 2003 (race condition bug → 50M
insan 14 gün elektriksiz); Mars Climate Orbiter $235M (imperial vs metric — semantic).

**Runtime nümunəsi:**
```go
nums := []int{2, 4, 6, 8}
for i := 0; i <= 10; i++ {     // 4 element, 11 iterasiya!
    total += nums[i]           // panic: index out of range
}
// Həll: for i := range nums
```

**Semantic nümunəsi:** `if km > 2 { car } else { walk }` — km=2-də "walk" (səhv);
`>=` düzəlişi.

### 2. Go vs Exception Handling
**Java/Python (try/catch/finally):** implicit; hər funksiya exception ata bilər —
hansını? bilmirsən; axın kəsilir.

**Go (error return):** EXPLICIT — funksiya `error` qaytarır; yoxlanmayıbsa GÖRÜNÜR:
```go
val, err := someFunc()
if err != nil {
    return err        // yuxarıya ötür
}
return nil
```
Xətalar funksiyalara ötürülən/qaytarılan ADİ DƏYƏRLƏRdir (Rob Pike: "Errors are values").

### 3. error İnterfeysi
```go
type error interface {
    Error() string        // yeganə şərt: Error() metodu → string
}
```
errors paketinin daxili (standart kitabxana):
```go
type errorString struct { s string }        // unexported
func (e *errorString) Error() string { return e.s }   // interfeysi satisfied edir
func New(text string) error { return &errorString{text} }
```
errorString UNEXPORTED-dır — `errors.errorString{}` birbaşa YAZILA BİLMƏZ; yalnız
New() ilə.

### 4. Custom Xəta Dəyərləri — Err Konvensiyası
**Kitabdan kod nümunəsi (payDay validasiyası):**
```go
var (
    ErrHourlyRate  = errors.New("invalid hourly rate")     // Err prefiks + camelCase
    ErrHoursWorked = errors.New("invalid hours worked per week")  // kiçik hərf, NO punctuation
)

func payDay(hoursWorked, hourlyRate int) (int, error) {
    if hourlyRate < 10 || hourlyRate > 75 {
        return 0, ErrHourlyRate
    }
    if hoursWorked < 0 || hoursWorked > 80 {
        return 0, ErrHoursWorked
    }
    if hoursWorked > 40 {
        hoursOver := hoursWorked - 40
        overTime := hoursOver * 2
        regularPay := hoursWorked * hourlyRate
        return regularPay + overTime, nil
    }
    return hoursWorked * hourlyRate, nil
}

pay, err := payDay(81, 50)
if err != nil { fmt.Println(err) }    // HƏR çağırışda err YOXLA
```
**Standart kitabxanadan nümunə:** http paketi — ErrBodyNotAllowed, ErrHijacked...
paket səviyyəsində qruplaşdırılmış Err dəyişənləri.

### 5. Panic — Nə Vaxt və Necə
**Nədir:** built-in funksiya; goroutine-un normal icrasını DAYANDIRIR; proqramın
BÜTÜNLÜYÜNÜ qorumaq üçün. Exception-dan fərqi: Go-da panic NORMA DEYİL — anormallıq.

**Panic addımları:**
1. İcra dayanır
2. Panicking funksiyanın defer-ləri çağırılır
3. Stack-də yuxarı bütün defer-lər çağırılır (main-ə qədər)
4. Handler yoxdursa → PROQRAM ÇÖKÜR

**Kitabdan kod nümunəsi (defer + panic birlikdə):**
```go
func test() {
    n := func() { fmt.Println("Defer in test") }
    defer n()
    message("good-bye")
}
func message(msg string) {
    f := func() { fmt.Println("Defer in message func") }
    defer f()
    if msg == "good-bye" {
        panic(errors.New("something went wrong"))   // error TİPİ ötür — idiomatik
    }
}
// Çıxış sırası: Defer in message → Defer in test → CRASH
```
**os.Exit vs panic:** os.Exit defer-ləri İŞLƏTMİR; panic İŞLƏDİR — panic üstündür.

**Panic-lü payDay (kitabdan):**
```go
func payDay(hoursWorked, hourlyRate int) int {
    report := func() {                       // debug məlumatı — həmişə çap olunacaq
        fmt.Printf("HoursWorked: %d\nHourlyRate: %d\n", hoursWorked, hourlyRate)
    }
    defer report()
    if hourlyRate < 10 || hourlyRate > 75 {
        panic(ErrHourlyRate)                // çağırıcıya etibar YOXDURSA
    }
    if hoursWorked < 0 || hoursWorked > 80 {
        panic(ErrHoursWorked)
    }
    ...
}
```

### 6. Recover — Panic-dən Qayıdış
**İmza:** `func recover() interface{}` — panic-ə ötürülən dəyəri qaytarır.
**ŞƏRT:** YALNIZ deferred funksiya daxilində işləyir; yoxsa panic davam edir.

**Kitabdan kod nümunəsi:**
```go
func b(msg string) {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("error in func b()", r)   // r = panic dəyəri (error)
        }
    }()
    if msg == "good-bye" {
        panic(errors.New("something went wrong"))
    }
}
// main → a() → b() → panic → b-nin defer-i recover edir → a() davam edir → main davam
```

**Tip-fərqləndirən recover (kitabdan):**
```go
defer func() {
    if r := recover(); r != nil {
        if r == ErrHourlyRate {                        // == ilə ERR YOXLAMASI
            fmt.Printf("hourly rate: %d\nerr: %v\n", hourlyRate, r)
        }
        if r == ErrHoursWorked {
            fmt.Printf("hours worked: %d\nerr: %v\n", hoursWorked, r)
        }
    }
}()
```
Panic-ə error ötürmək recover-də TİP YOXLAMASI verir — string ötürsən bu imkan YOXDUR.

### 7. Xəta İşləmə Qaydaları (kitabın cəmi)
1. Xəta dəyişəni `Err` prefiksi ilə + camelCase
2. Xəta mətni: kiçik hərflə, durğu işarəsiz (concatenation üçün)
3. Qaytarılan error HƏMİŞƏ yoxlansın — yoxlanmamış err = gözlənilməz davranış
4. panic(error) — boş dəyər YOX, error TİPİ ötür
5. Xətanın STRING dəyərini yoxlama (müqayisə etmə)
6. panic-i AZ işlət

### 8. Activity Zənciri (bank directDeposit)
4 addımlı təkamül: custom Err yaradılışı → validasiya metodları (routingNumber<100 →
ErrInvalidRoutingNum; boş lastName → ErrInvalidLastName) → panic-ə keçid →
defer+recover ilə crash-in qarşısı.

## Əsas terminlələr
- Syntax/Runtime/Semantic Error — 3 xəta sinfi
- Exception Handling — try/catch paradigması; implicit (Go-da YOX)
- error İnterfeysi — `Error() string` metodu; hər xətanın müqaviləsi
- errorString — errors.New-in unexported arxa tipi
- errors.New(text) — xəta dəyəri yaratmaq
- Err Prefiksi — xəta dəyişənlərinin idiomatik adı
- panic — proqramı qorumaq üçün çökdürmə; goroutine-u dayandırır
- Call Stack Unwinding — panic-də defer-lərin yuxarı çağırılması
- recover() — YALNIZ defer daxilində; panic dəyərini qaytarır
- os.Exit — defer-ləri işlətməyən dərhal çıxış
- Race Condition — 2 thread-in eyni yaddaşa yazışı (blackout nümunəsi)

## Praktik nətidə

(1) Semantic xəta ən təhlükəlidir — proqram "işləyir" amma yanlış; sərhəd
hallarını (>-, >=) ikicə dəfə yoxla. (2) Go-da xəta = dəyər: yoxlanmamış err
kod oxunuşunda DƏRHAL görünür — exception fəlsəfəsindən üstünlük. (3) Err-prefiks +
kiçik hərfli, durğusuz mətn — concat üçün. (4) Xətanı string kimi müqayisə ETMƏ —
dəyərlə (==) müqayisə. (5) panic: yalnız bütövlüyü qorumaq üçün; mütləq error
ötür — recover-də tip yoxlaması yaradacaq. (6) recover YALNIZ defer daxilində
işləyir — pattern: `defer func() { if r := recover(); r != nil { ... } }()`.
(7) panic vs os.Exit: defer-lər panic-də İŞLƏYİR — resurs təmizliyi qorunur.
(8) defer-dən report funksiyası — panic halında belə input məlumatı qeyd edir.
(9) İstifadəçi error-ları handle etməyə etibarlı deyilsə — error return əvəzinə
panic. (10) range loop — index xətalarının struktur həlli.

## Mənbə
Pages: 217-249 (PDF 250-283)
