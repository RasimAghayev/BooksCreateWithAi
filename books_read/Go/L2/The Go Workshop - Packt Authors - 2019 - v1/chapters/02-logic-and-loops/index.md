# Chapter 2 — Logic and Loops (Məntiq və Döngürlər)

## Bu fəsil nədən bəhs edir?

if/else if/else strukturları, initial if statement, if-ın switch-ə çevrilməsi
(expression switch, çoxlu case, ifadəsiz switch), for loop-un 3 forması (klassik,
şərtli, sonsuz), range ilə map/slice iterasiyası, break/continue, FizzBuzz activity
və Bubble Sort activity.

## Əsas fikirlər

### 1. if Statement — Əsas Məntiq
```go
if <bool ifadə> {
    // kod bloku — bool true olanda icra olunur
}
```
**Kitabdan kod nümunəsi (tək/cüt yoxlama):**
```go
input := 5
if input%2 == 0 {
    fmt.Println(input, "is even")
} else {
    fmt.Println(input, "is odd")     // dedduksiya: even deyilsə — odd
}
```
- `%` (modul) — bölgüdən qalan qalıq; tək/cüt: `%2 == 0`
- Yalnız funksiya scope-da istifadə olunur

### 2. else if — Çoxşaxlı Məntiq
```go
if input < 0 {
    fmt.Println("input can't be a negative number")
} else if input%2 == 0 {
    fmt.Println(input, "is even")
} else {
    fmt.Println(input, "is odd")
}
```
Yuxarıdan aşağı qiymətləndirilir; ilk true case İCRA OLUR və qalanları ATLANIR;
heç biri true deyilsə və else yoxdursa — heç nə icra olunmur.

### 3. Initial if Statement — Ən Vacib Go İdioması
**Nədir:** if-ın şərt hissəsindən ƏVVƏL kiçik bəyanat — `if <initial>; <şərt> { }`.
Dəyişən yalnız if bloku daxilində yaşayır → scope xaricinə "sizmir".

**İcazə verilən sadə bəyanatlar:** `:=`, assignment, ifadələr (`i++`), kanal göndərmə.
**QADAĞA:** `var` ilə bəyan — YALNIZ `:=` işləyir.

**Kitabdan kod nümunəsi (error yoxlama idioması):**
```go
func validate(input int) error {
    if input < 0 {
        return errors.New("input can't be a negative number")
    } else if input > 100 {
        return errors.New("input can't be over 100")
    } else if input%7 == 0 {
        return errors.New("input can't be divisible by 7")
    }
    return nil
}

func main() {
    input := 21
    if err := validate(input); err != nil {    // err YALNIZ bu blokda yaşayır
        fmt.Println(err)
        return
    }
    // burada err ARTIQ MÖVCUD DEYİL — scope təmizliyi
}
```
Bu, Go-nun ən məşhur pattern-idir: funksiya çağır + error yoxla + dəyişən avtomatik ölür.

### 4. FizzBuzz (Activity)
Qaydalar: 1-100; 3-ə bölünən → "Fizz"; 5-ə → "Buzz"; hər ikisinə → "FizzBuzz".
```go
for i := 1; i <= 100; i++ {
    if i%15 == 0 {          // 3 VƏ 5 → 15 (İLK yoxla!)
        fmt.Println("FizzBuzz")
    } else if i%3 == 0 {
        fmt.Println("Fizz")
    } else if i%5 == 0 {
        fmt.Println("Buzz")
    } else {
        fmt.Println(strconv.Itoa(i))   // int → string
    }
}
```
**Dərs:** kombinə halı əvvəlcədən yoxlanmalı — yoxsa heç vaxt tutulmur.

### 5. Expression switch
**Niyə:** böyüyən else-if zəncirlərini sıxışdırmaq üçün.

**Sintaksis:**
```go
switch <initial>; <ifadə> {    // hər ikisi OPTIONAL
case <dəyər>:
    // kod bloku — {} LAZIM DEYİL
case <dəyər1>, <dəyər2>:        // çoxlu dəyər vergüllə
    // ...
default:                        // = else; ən sonda saxlamaq good practice
}
```

**Kitabdan kod nümunələri:**

Dəyər match (həftə günü):
```go
switch dayBorn := time.Monday; dayBorn {
case time.Monday:
    fmt.Println("Monday's child is fair of face")
case time.Tuesday:
    fmt.Println("Tuesday's child is full of grace")
default:
    fmt.Println("Error, day born not valid")
}
```

Çoxlu case dəyəri (weekday vs weekend):
```go
switch dayBorn {
case time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday:
    fmt.Println("Born on a weekday")
case time.Saturday, time.Sunday:
    fmt.Println("Born on the weekend")
default:
    fmt.Println("Error, day born not valid")
}
```

İfadəsiz switch (mürəkkəb məntiq case-də):
```go
switch dayBorn := time.Sunday; {      // ifadə YOX → true kimi davranır
case dayBorn == time.Sunday || dayBorn == time.Saturday:
    fmt.Println("Born on the weekend")
default:
    fmt.Println("Born some other day")
}
```

**Vacib qaydalar:**
- Case-lər yuxarıdan-aşağı, soldan-sağa yoxlanılır; ilk match icra olunur
- **Fallthrough AVTOMATİK DEYİL** (C-dən fərqli!) — istəyirsənsə `fallthrough`
  açar sözü yaz
- `default` istənilən yerdə ola bilər, amma sonda saxla
- Case ifadələri full bool məntiqi daşıya bilər (ifadəsiz switch-də)
- Type switch ayrı formadır (sonrakı fəsildə)

### 6. for Loop — Tək Loop Konstruksiyası
**3 hissəli klassik forma:**
```go
for <initial>; <şərt>; <post> {
    // kod
}
for i := 0; i < 5; i++ { }        // ənənəvi sayğaclı loop
```

**Şərtli forma** (data mənbəyi bool qaytaranda — DB, fayl, socket):
```go
for <şərt> { }                     // while bənzəri
```

**Sonsuz loop:**
```go
for { }                            // break ilə dayandırmaq MÜTLƏQ
```
İnfinite loop baş verərsə — terminalı bağla; sistem zədələnmir.

**Kitabdan kod nümunəsi (slice üzərində):**
```go
names := []string{"Jim", "Jane", "Joe", "June"}
for i := 0; i < len(names); i++ {    // len() = uzunluq
    fmt.Println(names[i])
}
```

### 7. range Loop — Kolleksiyalar Üçün
**Map iterasiyası (ən vacib istifadə):**
```go
config := map[string]string{
    "debug":    "1",
    "logLevel": "warn",
    "version":  "1.2.1",
}
for key, value := range config {     // hər iterasiyada cüt qaytarır
    fmt.Println(key, "=", value)
}
// Lazımsız dəyişən: _ (blank identifier)
for _, value := range config { }
```
**Map sırası RANDOM-dır** — dilin şüurlu qərarı; sıraya etibar ETMƏ (pseudo-random
data mənbəyi kimi belə istifadə olunur).

**Activity — ən populyar söz (kitabdan):**
```go
words := map[string]int{"Gonna": 3, "You": 3, "Give": 2, "Never": 1, "Up": 4}
maxCount := 0
var topWord string
for word, count := range words {
    if count > maxCount {
        maxCount = count
        topWord = word
    }
}
fmt.Println("Most popular word:", topWord)   // Up, 4
```

### 8. break və continue
- **continue** — cari iterasiyanı YARIDA KƏS → post + şərt icra olunur → növbəti
  iterasiya (tək elementi ötür)
- **break** — loop-u TAM DAYANDIRIR (xəta varsa qalanını emal etməmək)

**Kitabdan kod nümunəsi:**
```go
for {
    r := rand.Intn(8)
    if r%3 == 0 {
        fmt.Println("Skip")
        continue                    // 3-ə bölünənləri ötür
    } else if r%2 == 0 {
        fmt.Println("Stop")
        break                       // 2-yə bölünəndə dayandır
    }
    fmt.Println(r)
}
```
İç-içə if-ləri və nəzarət dəyişənlərini AZALDIR — kodu təmizləyir.

### 9. Bubble Sort (Activity)
**Kitabdan ipucu kodları:**
```go
nums := []int{5, 8, 2, 4, 0, 1, 3, 7, 9, 6}
for i := 0; i < len(nums); i++ {
    for j := 1; j < len(nums)-i; j++ {
        if nums[j] < nums[j-1] {
            // In-place swap — müvəqqəti dəyişən YOX:
            nums[j], nums[j-1] = nums[j-1], nums[j]
        }
    }
}
// Yeni slice + append:
var nums2 []int
nums2 = append(nums2, 1)
```
Komponentlər: nested loop + müqayisə + paralel assign ilə swap.

## Əsas terminlələr
- Boolean Expression — true/false nəticələnən ifadə
- Modulus (%) — bölgü qalığı; tək/cüt/dülünbleme aləti
- Initial Statement — if/switch/for şərtindən əvvəl kiçik bəyanat; scope məhdudlaşdırır
- Expression Switch — dəyər/məntiq əsaslı budaqlanma
- Fallthrough — case-i növbətiyə keçirən açar söz (avtomatik DEYİL)
- Default Case — tutulmayan hallar (else kimi)
- Infinite Loop — `for {}` — break tələb edir
- Post Statement — hər iterasiya sonunda icra olunur (i++)
- range — açar/dəyər cütü ilə kolleksiya iterasiyası
- Map Random Order — map elementlərinin təsadüfi ardıcıllığı
- break/continue — loop-u dayandır / iterasiyanı ötür
- Blank Identifier `_` — lazımsız dəyəri udan ad
- In-place Swap — müvəqqəti dəyişənsiz dəyişdirmə (`a, b = b, a`)
- FizzBuzz — klassik intervyu tapşırığı (3/5/15 bölünmə)
- Bubble Sort — qonşu müqayisə ilə sıralama alqoritmi

## Praktik nətidə

(1) `if err := f(); err != nil` — Go-nun ən işlək idioması; err scope-dan kənar
sizmir. (2) else-if zənciri 3-dən çox olanda switch-ə keç. (3) Kombinə şərtlər
(FizzBuzz-də 15) ƏN ÜSTƏ yaz — yoxsa heç vaxt icra olunmur. (4) Fallthrough yalnız
AÇIQ şəkildə — Go default-da case-i bitirir. (5) Sonsuz `for {}` + break/continue =
nəzarətli axın. (6) Slice/array üçün `for i`; map üçün range. (7) Map sırasına
ETİBAR ETMƏ — randomdur; ehtiyac olsa sortla. (8) `_` ilə lazımsız dəyərləri ud.
(9) Paralel assign swap-ı müvəqqəti dəyişəndən üstün tut. (10) continue = elementi
özür, break = emalı dayandırır — anlayışını seçimini diktə edir.

## Mənbə
Pages: 55-81 (PDF 88-115)
