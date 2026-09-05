# Unit 0-1 — Getting Started + Imperative Programming (Lesson 1-5)

## Bu unit nədən bəhs edir?

Go-ya ilk addım: Go Playground, package/import/func, Print/Println/Printf, dəyişənlər və konstantalar, arifmetika, if/else/switch, for loop, short declaration (`:=`) və scope qaydaları. Capstone: Mars biletləri generatoru.

**PDF səhifələr:** 18-57 (Lesson 1: 18-26, L2: 28-37, L3: 38-48, L4: 49-55, L5: 56-57)

## Əsas fikirlər

### 1. Go nədir? (L1)
**Nədir:** Kompilyasiya olunan dil — icradan əvvəl bütün kod maşın dilinə çevrilir, tək executable hasil olur; kompilyator tipləri/typoları başlanğıcda yaxalayır.

**Müqayisə:** Python/Ruby interpreter-li — sətir-sətir tərcümə, səhvlər test edilməmiş yollarda gizlənir; amma yazmaq sürətli/dinamikdir. **Go ikisini birləşdirir**: "safety and performance of statically compiled languages + lightness and fun of dynamically typed interpreted languages" (Rob Pike).

**Rəqəmlərlə:** Iron.io 30 Ruby server → 2 Go server; Bitly Python→Go performans qazancı, sonra C-ni də əvəz etdi. Şirkətlər: Google, Docker, CloudFlare, Let's Encrypt, Twitch...

**Motto:** "Go is an open source programming language that enables the production of simple, efficient, and reliable software at scale."

**Axtarış üçün:** "golang" açar sözü istifadə et (golang.org).

### 2. İlk proqram (L1)
**Kitabdan kod nümunəsi:**
```go
package main

import (
    "fmt"
)

func main() {
    fmt.Println("Hello, playground")
}
```

**Sub-kod izahı:**
- `package main` → kodun paketi; bütün Go kodu paketlərdə təşkil olunur
- `import "fmt"` → formatlaşdırılmış giriş/çıxış paketi; hər paket bir ideyanı təmsil edir (math: Sin/Cos/Sqrt kimi)
- `func main()` → icra **main paketindəki main funksiyasından** başlayır; yoxdursa kompilyator xəta verir
- `fmt.Println(...)` → funksiyalar `paketAdı.` prefiksi ilə çağırılır — funksiyanın hansı paketdən olduğu dərhal görünür
- Go Playground (play.golang.org): Run → Google serverlərində kompilyasiya+icra; Share → paylaşım linki

**3 açar söz istifadə olundu** (Go-nun 25 açar sözündən): `package`, `import`, `func`.

### 3. The one true brace style (L1)
**Nədir:** `{` funksiya/keyword ilə EYNİ sətirdə, `}` ayrıca sətirdə — başqa yol YOXDUR:
```go
func main()
{                    // SYNTAX ERROR!
}
```
**Səbəb:** kompilyator nöqtəli vergülləri (semicolon) gizlincin əlavə edir — sətir sonlarında avtomatik; yanlış yerdə `;` → "unexpected semicolon" xətası. Nəticə: brace üslubu üzərində mübahisə mümkün deyil — hamı eyni formatda yazır.

### 4. Riyaziyyat və şərh (L2)
**Kitabdan kod nümunəsi:**
```go
// My weight loss program.
package main

import "fmt"

func main() {
    fmt.Print("My weight on the surface of Mars is ")
    fmt.Print(149.0 * 0.3783)
    fmt.Print(" lbs, and I would be ")
    fmt.Print(41 * 365 / 687)
    fmt.Print(" years old.")
}
```
- Operatorlar: `+ - * / %` (% = modulus — qalıq; 42 % 10 = 2)
- `//` → şərh — kompilyatorun üçün görünməz; `/* ... */` çoxsətərli də var
- Mars: çəki 37.83%, il 687 gün

### 5. Printf — format verb-lər (L2)
```go
fmt.Printf("My weight on the surface of Mars is %v lbs,", 149.0*0.3783)
fmt.Printf(" and I would be %v years old.\n", 41*365/687)
fmt.Printf("My weight on the surface of %v is %v lbs.\n", "Earth", 149.0)
```
- `%v` → növbəti arqumentin dəyəri ilə əvəz olunur (istənilən tip!)
- `\n` → yeni sətir (Printf/Print avtomatik keçmir; Println keçir)
- **Genişlik:** `%4v` → 4 simvola sol tərəfdən boşluqla doldur; `%-15v` → solda yaz, sağdan doldur:
```go
fmt.Printf("%-15v $%4v\n", "SpaceX", 94)
fmt.Printf("%-15v $%4v\n", "Virgin Galactic", 100)
// SpaceX          $  94
// Virgin Galactic $ 100
```

### 6. const və var (L2)
**Kitabdan kod nümunəsi:**
```go
const lightSpeed = 299792 // km/s
var distance = 56000000   // km
fmt.Println(distance/lightSpeed, "seconds")     // 186
distance = 401000000
fmt.Println(distance/lightSpeed, "seconds")     // 1337
```
- `const` → dəyişməz; yeni dəyər mənimsətmək → "cannot assign to lightSpeed" xətası
- `var` → proqram boyu dəyişə bilər; elan edilmiş dəyişən olmadan mənimsətmə (`speed = 16`) → xəta (typo-ları yaxalayır: `distence` vs `distance`)
- "Magic numbers" (mənasız literal-lar) → adlandırılmış const/var ilə əvəz et

**Elan formaları:**
```go
var distance = 56000000          // ayrı-ayrı
var (                           // qrup
    distance = 56000000
    speed    = 100800
)
var distance, speed = 56000000, 100800   // bir sətir
```

### 7. Qısayol operatorlar (L2)
```go
weight = weight * 0.3783   ⇔   weight *= 0.3783
age = age + 1               ⇔   age += 1  ⇔  age++
count--                     // azalt
price /= 2                  // böl
```
**Qeyd:** Go **prefix** `++count` dəstəkləmir (C/Java-dan fərqli).

### 8. Pseudorandom rəqəmlər (L2)
```go
package main

import (
    "fmt"
    "math/rand"
)

func main() {
    var num = rand.Intn(10) + 1   // 1-10
    fmt.Println(num)
}
```
- `rand.Intn(10)` → 0-9; `+ 1` olmadan **off-by-one error** (klassik səhv!)
- Import path `math/rand`, paket adı `rand`
- **GOTCHA:** `const` funksiya çağırışının nəticəsini qəbul edə BİLMƏZ — yalnız `var`
- Playground-da time dondurulub — hər run eyni "random" rəqəmlər

### 9. bool və müqayisə (L3)
```go
var walkOutside = true
var takeTheBluePill = false
```
- **Go-da yalnız `true` true-dur, yalnız `false` false-dur** — Python/JS-dəki kimi "" və ya 0 false DEYİL (Ruby-də true!)
- Operatorlar: `== != < > <= >=`
- JavaScript-in `===` (strict) lazımsızdır — Go-da tək `==` var və text ilə rəqəmi müqayisə etmək olmaz

### 10. if / else if / else (L3)
```go
var command = "go east"
if command == "go east" {
    fmt.Println("You head further up the mountain.")
} else if command == "go inside" {
    fmt.Println("You enter the cave where you live out the rest of your life.")
} else {
    fmt.Println("Didn't quite get that.")
}
```
- `=` (mənimsətmə) `==` (bərabərlik) əvəzinə yazilsa → kompilyator xətası yaxalayır

### 11. Məntiq operatorları + short-circuit (L3)
```go
// Leap year:
var year = 2100
var leap = year%400 == 0 || (year%4 == 0 && year%100 != 0)
```
- `||` (or), `&&` (and), `!` (not — bool-u tərsinə çevirir)
- **Short-circuit:** `||` birinci true-dan sonra qalanı qiymətləndirmir; `&&` birinci false-dan sonra dayanır
- 2100: `false || (true && false)` → false → "Keep your feet on the ground."

### 12. switch (L3)
**Forma 1 — dəyər müqayisəsi:**
```go
var command = "go inside"
switch command {
case "go east":
    fmt.Println("You head further up the mountain.")
case "enter cave", "go inside":        // vergüllə çoxlu dəyər!
    fmt.Println("You find yourself in a dimly lit cavern.")
case "read sign":
    fmt.Println("The sign reads 'No Minors'.")
default:
    fmt.Println("Didn't quite get that.")
}
```

**Forma 2 — şərtsiz switch (hər case if kimi):**
```go
var room = "lake"
switch {
case room == "cave":
    fmt.Println("You find yourself in a dimly lit cavern.")
case room == "lake":
    fmt.Println("The ice seems solid enough.")
    fallthrough                        // növbəti case-i də icra et!
case room == "underwater":
    fmt.Println("The water is freezing cold.")
}
```
- **Fallthrough default DEYİL** (C/Java-dan fərqli!) — yalnız `fallthrough` açar sözü ilə
- `default` — heç biri tutulmadıqda

### 13. for loop (L3)
```go
// Klassik form:
var count = 10
for count > 0 {
    fmt.Println(count)
    time.Sleep(time.Second)
    count--
}
fmt.Println("Liftoff!")

// Sonsuz loop + break:
var degrees = 0
for {
    fmt.Println(degrees)
    degrees++
    if degrees >= 360 {
        degrees = 0
        if rand.Intn(2) == 0 {
            break
        }
    }
}
```
- **Go-da yalnız `for` var** (while yoxdur!)
- Şərtsiz `for {}` → sonsuz; `break` → çıxış

### 14. Scope qaydaları (L4)
**Nədir:** dəyişən elan edildikdə **scope-a daxil olur** (görünür); scope bitəndə görünməz olur — `undefined` xətası.

**Faydaları:** eyni ad yenidən istifadə oluna bilər; oxuyarkən yalnız cari scope-dakıları yadda saxla.

**Sərhədlər = `{ }`:** funksiya → nested for → case blokları (case-lər brace-siz öz scope-u açır!).

```go
func main() {
    var count = 0
    for count < 10 {
        var num = rand.Intn(10) + 1   // num yalnız loop daxilində
        fmt.Println(num)
        count++
    }
    // num burada YOXDUR — undefined
}
```

### 15. Short declaration `:=` (L4)
```go
var count = 10   ⇔   count := 10
```
- `:=` daha qısa və **var-ın gedə bilmədiyi yerlərə gedir**

**for başlanğıcında:**
```go
for count := 10; count > 0; count-- {   // initialize; condition; post
    fmt.Println(count)
}
// count burada artıq scope-dan çıxıb — undefined
```

**if/switch daxilində:**
```go
if num := rand.Intn(3); num == 0 {
    fmt.Println("Space Adventures")     // num bütün branch-larda görünür
} else if num == 1 {
    fmt.Println("SpaceX")
} else {
    fmt.Println("Virgin Galactic")
}

switch num := rand.Intn(10); num {
case 0:
    fmt.Println("Space Adventures")
default:
    fmt.Println("Random spaceline #", num)
}
```

### 16. Package / function / block scope (L4)
```go
var era = "AD"                  // PACKAGE scope — bütün paketdə görünür

func main() {
    year := 2018                 // FUNCTION scope
    switch month := rand.Intn(12) + 1; month {
    case 2:
        day := rand.Intn(28) + 1 // BLOCK scope — yalnız bu case-də
        fmt.Println(era, year, month, day)
    case 4, 6, 9, 11:
        day := rand.Intn(30) + 1 // YENİ day — müstəqil dəyişən!
        fmt.Println(era, year, month, day)
    default:
        day := rand.Intn(31) + 1
        fmt.Println(era, year, month, day)
    }
}
```
**Vacib:** package scope-da `:=` MÜMKÜN DEYİL (yalnız var/const).
**Hər case öz scope-u** — 3 müstəqil `day` dəyişəni.

**Code smell dərsi:** yuxarıdakı kod 3x `Println` dublikasiya edir → refactoring: dəyişənləri geniş scope-a çıxar:
```go
year := 2018
month := rand.Intn(12) + 1
daysInMonth := 31
switch month {
case 2:
    daysInMonth = 28
case 4, 6, 9, 11:
    daysInMonth = 30
}
day := rand.Intn(daysInMonth) + 1
fmt.Println(era, year, month, day)
```
**Qayda:** dar scope mental yükü azaldır, amma həddindən artıq dar → dublikasiya → code smell. Case-by-case.

### 17. Capstone: Ticket to Mars (L5)
Tapşırıq: 10 random bilet cədvəli:
```
Spaceline        Days Trip type  Price
======================================
Virgin Galactic  23   Round-trip $  96
SpaceX           31   One-way    $  41
```
- Spaceline: Space Adventures / SpaceX / Virgin Galactic (random)
- Tarix: 13 Oktyabr 2020, məsafə 62,100,000 km
- Sürət: 16-30 km/s random → müddət = məsafə/sürət/86,400 (gün)
- Qiymət: $36-50M (sürətli = bahalı); Round-trip = 2x
- İstifadə: `%-15v $%4v` align, `rand.Intn`, `switch`, `for`

## Əsas terminlər
- Compiler vs Interpreter (kompilyator vs tərcüməçi)
- Package / import (paket / idxal)
- Format Verb (`%v` — format feyli)
- Escape Sequence (`\n`)
- Padding (`%4v` / `%-15v`)
- Constant / Variable (konstant / dəyişən)
- Magic Number (maqik rəqəm)
- Modulus (`%` — qalıq)
- Off-by-one Error (bir-vahid xətası)
- Pseudorandom (psevdotəsadüfi)
- Boolean / Comparison Operator
- Short-circuit Logic
- Fallthrough
- Infinite Loop (`for {}`)
- Scope (package/function/block)
- Short Declaration (`:=`)
- Code Smell / Refactoring

## Praktik nəticə
- `{` həmişə eyni sətirdə — one true brace style; gofmt bunu dəbi məcbur edir.
- `rand.Intn(N)` → 0-(N-1); 1-N lazımdırsa `+ 1` — off-by-one-dan qorx.
- `if num := ...; num == x` — şərt dəyişənini if-in öz scope-una həbs et.
- Dar scope dublikasiya yaradırsa → genişləndir; amma lazımsız geniş scope-da qlobal çirk yaratma.
- `const` funksiya nəticəsi qəbul edə bilməz; package scope-da `:=` işləmir.
- `switch` + vergüllü case dəyərləri + `fallthrough`-un ekspliçitliyi — Go-nun təhlükəsiz seçimləri.

## Mənbə
Pages: 18-57 (PDF), book pages 3-42 (Lesson 1-5)
