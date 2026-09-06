# Chapter 4 — Blocks, Shadows, and Control Structures (Bloklar, Kölgələmə və İdarə Strukturları)

## Bu chapter nədən bəhs edir?

Bloklar və identifier görünüşü, dəyişən kölgələməsi (shadowing), `if`, `for`-un dörd
forması, `for-range` (string/map/slice üzrə), label-lər, `switch` (blank switch daxil)
və `goto`-nun qaydaları.

## Əsas fikirlər

### 1. Bloklar
**Nədir:** Bəyanatların görünmə sahəsini müəyyən edən ərazi.

**Necə işləyir:** Blok iyerarxiyası: **universe block** (ən üst — `int`, `string`, `true`,
`nil`, `make` kimi predeclared identifier-lər; 25 keyword-dən fərqli olaraq bunlar
keyword DEYİL, shadow oluna bilirlər!) → **package block** (funksiyadan kənar bəyanatlar)
→ **file block** (import adları) → **funksiya bloku** (parametrlər daxil) → hər `{}` daxili blok.
Daxili blok xarici blokdan gələn hər identifier-ə çıxış edir.

### 2. Shadowing (Kölgələmə) — Səssiz Bug Mənbəyi
**Nədir:** Daxili blokda xarici blokla eyni adlı dəyişən bəyan etmək.

**Necə işləyir:** Kölgələyən dəyişən mövcud olduqca kölgələnənə çıxış YOXDUR; blok
başqa çatanda kölgələnən dəyişən dəyəri ilə qayıdır.

**Kitabdan kod nümunəsi:**
```go
func main() {
    x := 10
    if x > 5 {
        fmt.Println(x)  // 10
        x := 5          // yeni x — kölgələyir!
        fmt.Println(x)  // 5
    }
    fmt.Println(x)        // 10 — x heç yerə getməyib
}
```

**`:=` ilə tələ:** `x, y := 5, 20` — sol tərəfdə ən azı bir yeni dəyişən olduqda
leqaldır, amma **yalnız cari blokda bəyan olunanları yenidən istifadə edir** — xarici
blokdakı `x` yenə də kölgələnir (təyinat YOX, yeni dəyişən!). Package adlarını da
kölgələmək olar: `fmt := "oops"` → `fmt.Println` artıq `string` tipində çağrılır →
compile xətası.

**Müdafiə:** `shadow` linter-i (`go install golang.org/x/tools/go/analysis/passes/shadow/cmd/shadow@latest`)
+ Makefile vet target-inə əlavə. **Diqqət:** universe block identifier-lərinin kölgələnməsini
heç bir linter tutmur — `true := 10` leqal compile olunur!

### 3. `if` və scoped bəyanatlar
**Nədir:** `if`-in Go spesifikası — şərt sahəsinə məxsus dəyişən bəyanatı.

**Necə işləyir:** `if n := rand.Intn(10); n == 0 { ... }` — `n` if/else zəncirinin
hansı sahəsindədirsə yalnız orada görünür; zəncirdən sonra `undefined: n`. Bu,
dəyişənləri lazım olduqları yerdə yaratmağa imkan verir. Şərtdə sadə statement qoymaq
texniki olaraq mümkündür, amma yalnız yeni dəyişən bəyanı üçün işlət (digərləri qəlizdir).
Qısa if/else gövdələri solda düzülür; nested kod çətin oxunur.

**Kitabdan kod nümunəsi:**
```go
if n := rand.Intn(10); n == 0 {
    fmt.Println("That's too low")
} else if n > 5 {
    fmt.Println("That's too big:", n)
} else {
    fmt.Println("That's a good number:", n)
}
```

### 4. `for` — Dörd Forma
**Nədir:** Go-da tək loop keyword — `for`.

**Kitabdan kod nümunəsi:**
```go
// 1. Tam forma (C-style)
for i := 0; i < 10; i++ { fmt.Println(i) }

// 2. Yalnız şərt (while əvəzi)
i := 1
for i < 100 { fmt.Println(i); i = i * 2 }

// 3. Sonsuz
for { fmt.Println("Hello") }

// 4. for-range
for i, v := range evenVals { fmt.Println(i, v) }
```

**Sub-kod izahı:**
- İnisializasiyada `var` QADAĞANDIR, yalnız `:=`
- `do/while` yoxdur → `for { ...; if !COND { break } }` pattern-i
- `break` — dərhal çıxış; `continue` — təkrarın qalanını ötür (nested if zəncirlərini
  yox, düz if-lər sırasını mümkün edir — FizzBuzz nümunəsində göstərilib)

### 5. `for-range` — Kompozit Tiplər üzrə
**Nədir:** Built-in kompozit tiplərin iterasiyası; 2 dəyişən: mövqe/açar + dəyər.

**Qaydalar:**
- İstifadə edilməyən dəyişən üçün `_` (unused qaydası loop dəyişənlərinə də aiddir);
  dəyər lazım deyilsə sadəcə tərk et: `for k := range m`
- Adlandırma: array/slice/string-də `i`, map-də `k`, dəyər `v`
- **Map iterasiyası təsadüfi sırada gedir** — təhlükəsizlik xüsusiyyətidir (Hash DoS
  hücumlarının qarşısını alır; hər map yarandığında random hash seed). `fmt.Println`
  istisnadır — debug üçün açarları sortlu çap edir.
- **String üzərində for-range bayt YOX, rune üzrə iterasiya edir**: `apple_π!` —
  `6 960 π` sonra `8 33 !` (7-ci bayt `π`-nin ikinci baytıdır). Mövqe = bayt ofseti,
  dəyərin tipi = rune. Yanlış UTF-8 baytı → `0xfffd` replacement simvolu.
- **Dəyər kopyadır**: `for _, v := range evenVals { v *= 2 }` orijinalı dəyişmir.
  (Bunun goroutine təsirləri Chapter 10-da.)
- Standart for ilə string-in başını ötürmək multibyte simvolları korlayır — string-in
  bir hissəsini keçmək üçün mütləq for-range.

**Label-lər:** nested loop-dan xarici loop-un iterasiyasını idarə etmək üçün:
```go
outer:
    for _, sample := range samples {
        for i, r := range sample {
            if r == 'l' { continue outer }
        }
    }
```
Label `go fmt` tərəfindən funksiya səviyyəsinə düzülür (asan görünsün deyə).

**Hansı forma nə vaxt:** Hamısını iterasiya edirsən → for-range. Bir hissəni (orta
elementləri) → standart for (`for i := 1; i < len-1; i++` daha qısa və aydındır).
Hesablanan şərtlə → condition-only. do/while və iterator pattern → infinite for.

### 6. `switch` — Fall-Through yoxdur
**Nədir:** Go-nun switch-i C-dən fərqli olaraq faydalıdır.

**Necə işləyir:**
- Case-lər avtomatik break olunur (fall-through DEFAULT YOXDUR; `fallthrough` keyword-ü
  var, amma istifadə etmə — asılılığı aradan qaldır).
- Eyni məntiq üçün vergüllə çoxlu dəyər: `case 1, 2, 3, 4:`
- Boş case = heç nə etmə (növbəti case-ə keçmir!)
- `if` kimi şərt sahəsinə bəyanat: `switch size := len(word); size { ... }`
- Hər case öz blokudur — case daxilində yeni dəyişən bəyan olunur, yalnız orada görünür
- == ilə müqayisə olunan hər tip switch edilə bilər (slice/map/channel/funksiya YOX)

**Kitabdan kod nümunəsi:**
```go
switch size := len(word); size {
case 1, 2, 3, 4:
    fmt.Println(word, "is a short word!")
case 5:
    wordLen := len(word)   // case-blok daxili dəyişən
    fmt.Println(word, "is exactly the right length:", wordLen)
case 6, 7, 8, 9:           // boş case — heç nə olmur
default:
    fmt.Println(word, "is a long word!")
}
```

**Tələ (break-in case daxilində):** for içində switch-də `break` switch-i break edir,
loop-u YOX → loop-dan çıxmaq üçün label lazımdır: `loop:` + `break loop`.

### 7. Blank Switch (Boş switch)
**Nədir:** Müqayisə dəyəri göstərilməyən switch — hər case istənilən bool ifadə.

**Kitabdan kod nümunəsi:**
```go
switch wordLen := len(word); {
case wordLen < 5:
    fmt.Println(word, "is a short word!")
case wordLen > 10:
    fmt.Println(word, "is a long word!")
default:
    fmt.Println(word, "is exactly the right length.")
}
```

**Qayda:** Əgər bütün case-lər eyni dəyişənin bərabərliyidirsə → adi expression switch
işlət (`switch a { case 2: ... }`). Əlaqəli çoxlu müqayisələr → blank switch (if/else
zəncirindən daha oxunaqlı — müqayisələr sol kənarda düzülür).

### 8. `goto` — Bəli, Go-da var
**Nədir:** Label-li sətirə atlayan statement.

**Qaydalar (məhdudiyyətlər):** Dəyişən bəyanatının üstündən atlamaq QADAĞAN;
daxili/paralel bloka atlamaq QADAĞAN (compile xətası). Yalnız eyni funksiya daxilində.

**Nə vaxt:** Praktikada heç vaxt; istisna — bir sıra şərtdən sonra ortaq "tullantı
təmizləmə" koduna sıçramaq (boolean flag qəlizliyi və ya kod dublikası əvəzinə).
Real nümunə: standart kitabxanadakı `strconv/atof.go`-nun `floatBits` metodu
(`overflow:` və `out:` label-ləri).

**Kitabdan kod nümunəsi:**
```go
a := rand.Intn(10)
for a < 100 {
    if a%5 == 0 {
        goto done
    }
    a = a*2 + 1
}
fmt.Println("loop normal bitəndə işə düşən məntiq")
done:
fmt.Println("hər halda icra olunan son kod")
```

## Əsas terminlər
- Shadowing (kölgələmə) — daxili blokda eyni adlı identifier-in xaricini gizlətməsi
- Universe block (kainat bloku) — predeclared identifier-lərin (int, nil, make) bloku
- for-range (aralıq dövrü) — kompozit tip üzrə iterasiya konstruktoru
- Blank switch (boş switch) — müqayisə dəyəri olmayan, bool case-lərə sahib switch
- Label (etiket) — `break`/`continue`/`goto` üçün adlandırılmış hədəf
- Hash DoS (xəş DoS hücumu) — eyni bucket-ə hash olunan açarlarla slowdown hücumu

## Praktik nəticə

Shadowing Go-nun ən səssiz bug mənbəyidir: `:=` çoxdəyişənli təyinatda xarici dəyişənləri
yenidən istifadə ETMİR — kölgələyir. Müdafiə: shadow linter + universe identifier-lərə
(LC adlarla) diqqət. Kod tərzi: nested if yerinə continue; əlaqəli müqayisələrdə blank
switch; loop-un ortasından çıxış üçün label. Map sırasına heç vaxsoy trust etmə.

## Mənbə
Pages: 101-132
