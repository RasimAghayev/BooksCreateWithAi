# Chapter 2 — Command and Control (səh. 82-109)

## Bu fəsil nədən bəhs edir?

Axın idarəetməsi: if / if-else / else if, initial if statements (scope
təmizliyi), expression switch (çoxlu case dəyərləri, expressionless
switch, fallthrough), for loop-un bütün formaları (klassik, şərtli,
sonsuz, range), break/continue, goto və activity-lər (FizzBuzz, bubble
sort).

## Əsas fikirlər

### 1. if / if else / else if
```go
// Sadə if — modul ilə tək/cüt:
input := 5
if input%2 == 0 {
    fmt.Println(input, "is even")
}
if input%2 == 1 {
    fmt.Println(input, "is odd")
}

// if else — dedduksiya (tək deyilsə, cütdür):
if input%2 == 0 {
    fmt.Println(input, "is even")
} else {
    fmt.Println(input, "is odd")
}

// else if — çoxşaxəli (mənfilər üçün əvvəl yoxla!):
input := -10
if input < 0 {
    fmt.Println("input can't be a negative number")
} else if input%2 == 0 {
    fmt.Println(input, "is even")
} else {
    fmt.Println(input, "is odd")
}
```
- if yalnız funksiya scope-un-da işlədilir
- Şərtlər yuxarıdan aşağı qiymətlənir; true tapılanda qalanlar ATLANIR

### 2. initial if statements (scope təmizliyi)
**Nədir:** `if <initial>; <boolean>` — dəyişəni if-in ÖZ scope-unda saxla;
istifadəsiz dəyişənlər ətrafda qalmasın.

```go
func validate(input int) error {
    if input < 0 {
        return errors.New("input can't be a negative number")
    } else if input > 100 {
        return errors.New("input can't be over 100")
    } else if input%7 == 0 {
        return errors.New("input can't be divisible by 7")
    } else {
        return nil
    }
}

input := 21
if err := validate(input); err != nil {   // err YALNIZ burada yaşayır
    fmt.Println(err)
    return
} else if input%2 == 0 {
    fmt.Println(input, "is even")
} else {
    fmt.Println(input, "is odd")
}
```
- İcazəli simple statements: `i := 0`, `i = (j*10) == 40`, `i++`,
  kanal göndərməsi; **var QADAĞANDIR** — qısa təyinat işlədin
- error + initial if = Go-nun ən Klassik pattern-i

### 3. Expression switch
```go
switch <initial>; <expression> {
case <expr>, <expr>:
    <statements>
default:
    <statements>
}
```

**Dəyər əsaslı (gün doğum):**
```go
dayBorn := time.Monday
switch dayBorn {
case time.Monday:
    fmt.Println("Monday's child is fair of face")
case time.Tuesday:
    fmt.Println("Tuesday's child is full of grace")
// ... bütün günlər
default:
    fmt.Println("Error, day born not valid")
}
```

**Çoxlu case dəyərləri (, ilə):**
```go
dayBorn := time.Sunday
switch dayBorn {
case time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday:
    fmt.Println("Born on a weekday")
case time.Saturday, time.Sunday:
    fmt.Println("Born on the weekend")
default:
    fmt.Println("Error, day born not valid")
}
```

**Expressionless switch (yalnız initial):**
```go
switch dayBorn := time.Sunday; {
case dayBorn == time.Sunday || dayBorn == time.Saturday:
    fmt.Println("Born on the weekend")
default:
    fmt.Println("Born some other day")
}
```

**Vacib qaydalar:**
- Uyğun case TAPILANDA yalnız ONUN statements-i icra olunur —
  başqa dillərdəki kimi fall-through avtomatik YOXDUR
- `fallthrough` açar sözü ilə süni keçid mümkündür
- default istənilən yerdə ola bilər — adətən sonda (else rolunu oynayır)
- Case-lər yuxarıdan-aşağı, soldan-sağa yoxlanılır
- Hər üç hissə (initial, expression, case) mümkün qədər sərbəst buraxıla
  bilər: `switch {` = `switch true {`

### 4. for loop formaları
```go
// 1) Klassik (initial; condition; post):
for i := 0; i < 5; i++ {
    fmt.Println(i)     // 0..4
}

// 2) Kollesiya üzrə (len ilə):
names := []string{"Jim", "Jane", "Joe", "June"}
for i := 0; i < len(names); i++ {
    fmt.Println(names[i])
}

// 3) Şərtli (while kimi — data mənbəyi Boolean qaytaranda):
for <condition> { ... }

// 4) Sonsuz (break ilə dayandırılır):
for { ... }

// 5) range — map/slice üzrə:
for key, value := range config { ... }
```

**range ilə map:**
```go
config := map[string]string{
    "debug":    "1",
    "logLevel": "warn",
    "version":  "1.2.1",
}
for key, value := range config {
    fmt.Println(key, "=", value)
}
```
- **Map sırası RANDOM-dır** — sıraya etibar etməyin (pseudo-randomizasiya
  kimi istifadə oluna bilər)
- Lazımsız dəyişən: `_` (for _, v := range s)

### 5. break və continue
```go
for {
    r := rand.Intn(8)
    if r%3 == 0 {
        fmt.Println("Skip")
        continue          // bu iterasiyanı ötür → post + condition
    } else if r%2 == 0 {
        fmt.Println("Stop")
        break            // loop-u tam dayandır
    }
    fmt.Println(r)
}
```
- continue: cari iterasiyanın QALANINI ötür; break: loop bitir
- continue → bir element xarab, qalanı emal et; break → xəta varsa,
  hamısını dayandır
- Sonsuz loop break-siz → prosesi əl ilə öldür (terminalı bağla)

### 6. goto
```go
for {
    r := rand.Intn(8)
    if r%3 == 0 {
        fmt.Println("Skip")
        continue
    } else if r%2 == 0 {
        fmt.Println("Stop")
        goto STOP          // label-ə atla
    }
    fmt.Println(r)
}
STOP:
    fmt.Println("Goto label reached")
```
- Label funksiya daxilində olmalıdır — xaricə goto compile xətası
- Standart kitabxanada da var (math) — ləqəbsiz dəyişənləri azaldır
- **Ehtiyatlı:** oxunaqlılığı azaldır; sadə hallarda saxla

### Activity-lər (praktika)

**2.01 — ən populyar söz (map + range):**
```go
words := map[string]int{
    "Gonna": 3, "You": 3, "Give": 2, "Never": 1, "Up": 4,
}
// range ilə gəz; max count saxla; "Up 4" çap et
```

**2.02 — FizzBuzz (müsahibə klassiki):**
```go
// 1..100: 3-ə bölünən → "Fizz"; 5-ə → "Buzz"; hər ikisinə → "FizzBuzz"
// strconv.Itoa(i) ilə rəqəmi string-ə çevir
```

**2.03 — Bubble sort (swap ilə):**
```go
// In-place swap:
nums[i], nums[i-1] = nums[i-1], nums[i]
// Yeni slice + append:
nums2 = append(nums2, 1)
```

## Əsas terminlər
- Boolean expression — true/false nəticələn ifadə
- Modulus (%) — bölmə qalığı
- Deduktiv məntiq — "xeyrsə, deməli" qısaltması
- Initial statement — if/switch/for-un `;`-dən əvvəlki hissəsi
- Simple statement — ifadə/təyinat/increment (var yox!)
- Expression switch — dəyər müqayisəli switch
- fallthrough — case-dən case-ə süni keçid
- default — switch-in "else"-i
- Infinite loop (`for {}`) — break tələb edən dövr
- break / continue — dayandır / ötür
- range — kolleksiya iterasiyası (key, value verir)
- Map iteration order — random (qəsdən!)
- goto + label — funksiya daxili sıçrayış
- FizzBuzz / bubble sort — klassik məşqlər

## Praktik nəticə
Məntiq qurarkən: iki variant → if/else; çoxşaxə → əvvəl xüsusi hallar
(mənfi kimi), sonra ümumi; dəyər müqayisəsi çoxdursa → switch (, ilə
qruplaşdır); mürəkkəb şərtlər → expressionless switch. Hər error
yoxlamasında `if err := f(); err != nil` — dəyişən scope-da qalır. Loop-lar:
indeksli kolleksiyalar → klassik for; map → range (sıra random!);
sonsuz + break/continue → axın nəzarəti. goto yalnız funksiya daxilində və
ehtiyatla.

## Mənbə
Pages: 82-109 (PDF 82-109)
