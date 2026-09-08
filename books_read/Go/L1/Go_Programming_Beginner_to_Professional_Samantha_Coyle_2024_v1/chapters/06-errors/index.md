# Chapter 6 — Don't Panic! Handle Your Errors (səh. 238-271)

## Bu fəsil nədən bəhs edir?

Xəta növləri (syntax/runtime/semantic), Go-nun exception modelindən
fərqi, error interfeysi, errors.New ilə xəta yaratma (Err adlandırma),
panic (defer-lərin icrası, stack yuxarı), recover (yalnız deferred
daxilində), qaydalar toplusu və error wrapping (%w).

## Əsas fikirlər

### 1. Xəta növləri
| Növ | Nə vaxt | Nümunə | Aşkarlanma |
|---|---|---|---|
| Syntax | yazılışda | `fmt.println` (kiçik p) | IDE/golint dərhal |
| Runtime | icrada | slice index out of range | çökmə/panic |
| Semantic (logic) | icrada, SƏSSİZ | `>` əvəzinə `>=` lazımi | səhv NƏTİCƏ |

**Tarixi dərslər:** Northeast Blackout 2003 (race condition — 50M adam
14 gün electricsiz); Mars Climate Orbiter 1998 ($235M — imperial/metric
vahid səhvi) = klassik semantic error.

**Runtime xəta nümunəsi:**
```go
nums := []int{2, 4, 6, 8}
for i := 0; i <= 10; i++ {      // 4 element, 11 iterasiya!
    total += nums[i]             // PANIC: index out of range
}
// Düzgün: for i := range nums
```

**Semantic xəta nümunəsi:**
```go
km := 2
if km > 2 {                      // >= olmalı idi — km=2 "walk" çap olundu
    fmt.Println("Take the car")
}
```

### 2. Exceptions vs Go errors
```java
// Java/Python — implicit, axını kəsir:
try { /* code */ } catch (Exception e) { /* handle */ }
```
```go
// Go — EXPLICIT, axın görünür:
val, err := someFunc()
if err != nil {
    return err
}
```
- Go-da istənilən funksiya exception "ata bilər" və növü bəlli deyil —
  Go isə error DƏYƏRİ qaytarır; yoxlamayan AÇIQ görünür.

### 3. Error interfeysi
```go
type error interface {
    Error() string
}
```
- **"Errors are values" (Rob Pike)** — proqramlaşdırıla bilən dəyərlərdir
- Standart kitabxananın daxili tətbiqi:
```go
type errorString struct {     // unexported!
    s string
}
func (e *errorString) Error() string {
    return e.s
}
```
- errors.errorString-ə birbaşa çıxış YOXDUR — yalnız Error() metodu ilə

### 4. Xəta dəyərlərinin yaradılması
```go
func New(text string) error {
    return &errorString{text}
}

// İdiomatik adlandırma — Err prefiks + camelCase:
var (
    ErrHourlyRate  = errors.New("invalid hourly rate")
    ErrHoursWorked = errors.New("invalid hours worked per week")
)
```
- Standart kitabxana nümunəsi (http): ErrBodyNotAllowed, ErrHijacked...
- Mesaj: kiçik hərflə, nöqtəsiz (concatenation üçün)

**PayDay nümunəsi (Exercise 6.03):**
```go
func payDay(hoursWorked, hourlyRate int) (int, error) {
    if hourlyRate < 10 || hourlyRate > 75 {
        return 0, ErrHourlyRate
    }
    if hoursWorked < 0 || hoursWorked > 80 {
        return 0, ErrHoursWorked
    }
    if hoursWorked > 40 {                 // overtime 2x
        hoursOver := hoursWorked - 40
        overTime := hoursOver * 2
        regularPay := hoursWorked * hourlyRate
        return regularPay + overTime, nil
    }
    return hoursWorked * hourlyRate, nil
}

pay, err := payDay(81, 50)
if err != nil {
    fmt.Println(err)                      // invalid hours worked per week
}
```

### 5. Panic
**Nədir:** proqramı çökdürən built-in funksiya; normal icranı dayandırır,
defer-ləri çağırır, stack boyu yuxarı qalxır, main-də çökir.

**Panic addımları:**
1. İcra dayanır → 2. Panicking funksiyanın defer-ləri → 3. Stack-dəki
bütün defer-lər → 4. main() → 5. Sonrakı sətirlər İCRA OL MUR → 6. Crash

```go
func message(msg string) {
    if msg == "good-bye" {
        panic(errors.New("something went wrong"))   // error ötür — idiomatik!
    }
}

func main() {
    message("good-bye")
    fmt.Println("This line will not get printed")   // ÇAP OLUNMAZ
}
```

**defer + panic (defer-lər yenə də işləyir):**
```go
func test() {
    defer n()                    // "Defer in test" çap olunur
    message("good-bye")          // panic
}
// Çap sırası: "Defer in message func" → "Defer in test" → crash
```

**os.Exit vs panic:** os.Exit defer-ləri ÇAĞIRMIR; panic çağırır —
təmizlik lazımdırsa panic üstündür.

**Panic-lə payDay (Exercise 6.04):**
```go
func payDay(hoursWorked, hourlyRate int) int {
    report := func() {
        fmt.Printf("HoursWorked: %d\nHourlyRate: %d\n", hoursWorked, hourlyRate)
    }
    defer report()                 // panic-də belə hesabat!

    if hourlyRate < 10 || hourlyRate > 75 {
        panic(ErrHourlyRate)
    }
    if hoursWorked < 0 || hoursWorked > 80 {
        panic(ErrHoursWorked)
    }
    // ... hesablama
}
```

### 6. Recover
**Nədir:** panicking goroutine-in idarəsini geri qaytaran funksiya.
**YALNIZ deferred funksiya daxilində işləyir!**

```go
func recover() interface{}       // panic-ə ötürülən dəyəri qaytarır

func b(msg string) {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("error in func b()", r)
        }
    }()
    if msg == "good-bye" {
        panic(errors.New("something went wrong"))
    }
}
// panic → deferred recover tutur → normal axın a() funksiyasına QAYTARIR
```

**Recover + error yoxlaması (Exercise 6.05):**
```go
func payDay(hoursWorked, hourlyRate int) int {
    defer func() {
        if r := recover(); r != nil {
            if r == ErrHourlyRate {              // xəta MÜQAYİSƏSİ
                fmt.Printf("hourly rate: %d\nerr: %v\n\n", hourlyRate, r)
            }
            if r == ErrHoursWorked {
                fmt.Printf("hours worked: %d\nerr: %v\n\n", hoursWorked, r)
            }
        }
        fmt.Printf("Pay was calculated based on:\nhours worked: %d\nhourly Rate: %d\n",
            hoursWorked, hourlyRate)
    }()
    // ... panic-lər
}
// payDay(100, 200) → recover → "hourly rate: 200 err: invalid hourly rate"
// proqram ÇÖKMÜR, davam edir
```

### 7. Qaydalar toplusu (errors & panics)
1. Xəta dəyişəni `Err` prefiksi ilə: `ErrExampleNotAllowed`
2. Mesaj: kiçik hərf, nöqtəsizz (concat üçün)
3. Qaytarılan hər error QİYMƏTLƏNDİRİLMƏLİDİR (yoxlamamaq = pis kod)
4. `panic()`-ə error tipi ötür (boş yox)
5. Error STRING-ə görə yoxlama — type assertion/interfeys metodları ilə
6. `panic()` AZ işlədin — birinci müdafiə xətti DEYİL
7. Error = gözlənilən vəziyyətlər; panic = istisnai/qeyri-bərpa olunan

### 8. Error wrapping
```go
// Go 1.13+ standart yolu — %w verb:
return fmt.Errorf("failed to read config file: %w", err)

// Üçüncü tərəf (köhnə):
return errors.Wrap(err, "failed to read config file")   // pkg/errors

// Çoxlu xəta:
github.com/hashicorp/go-multierror
```
- Zəncir: orijinal error SAXLANIR + kontekst əlavə olunur → debugging
  asanlaşır
- Qeyd: həddindən artıq kontekst yuxarı ötürməyin — təhlükəsizlik

## Activity icmalları (bank app progresiyası)
- 6.01: ErrInvalidLastName + ErrInvalidRoutingNumber yaratma
- 6.02: directDeposit struct + validateRoutingNumber/validateLastName
  metodları (error qaytaran)
- 6.03: eyni validasiya, amma panic ilə
- 6.04: panic-i defer+recover ilə tutub çap etmə

## Əsas terminlər
- Syntax / runtime / semantic (logic) error
- Race condition — paralel yazı yarışı (Blackout səbəbi)
- Exception handling vs explicit error checking
- error interface — Error() string metodu
- errorString — unexported daxili tip
- errors.New — xəta konstruktoru
- Err adlandırma konvensiyası
- Panic — defer-ləri çağıran çökmə
- Stack unwinding — defer-lərin yuxarı icrası
- os.Exit vs panic — defer fərqi
- recover() — yalnız deferred-də; panic dəyərini qaytarır
- %w verb / error wrapping — kontekstli xəta zənciri
- go-multierror — çoxlu xəta birləşdirmə

## Praktik nəticə
Xətalar dəyər kimi qaytarılır: `(T, error)` + `if err != nil` — hər
çağırışda. Paket səviyyəsində `Err...` konstantaları (validasiya
halları üçün). Panic yalnız bərpaolunmaz/qeyri-iddialı hallarda; onu da
error arqumenti ilə; xarici koda qarşı defer+recover sargısı. Yuxarı
özətərkən `%w` ilə kontekst əlavə et — amma originalı saxla. String
parse etməklə xəta yoxlama — tip müqayisəsi ilə (r == ErrHourlyRate).

## Mənbə
Pages: 238-271 (PDF 238-271)
