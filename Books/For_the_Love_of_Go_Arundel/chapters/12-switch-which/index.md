# Chapter 12 — Switch which? (Hansına Keçək?)

## Bu fəsil nədən bəhs edir?

switch statement: N mümkün yol; case-lər ardıcıl yoxlanır, İLK match icra
olunur (yalnız bir!), fallthrough (nadir), default case (həmişə yaz — clarity),
case-lər ƏLAQƏLİ olmalıdır (unrelated if-ləri gizlətmək üçün YOX), switch
expression (`switch x { case 1: ... }` — "switches on x", tip uyğunluğu),
comma-separated çoxdəyərli case (`case 1, 2, 3:`), break (case-dən tez çıxış —
"məhbus cəzasını bitirmək istəməyən kimi"), for loops: conditional (`for x <
10`), forever (`for {}` — istehsal robotları, HTTP serverlər), range
(collection üzrə: `for range employees`, `for i, e := range`, `for _, e :=`,
`for i := range`), 3-hissəli for (`for x := 0; x < 10; x++` — init;cond;post,
Go-da NADİR — range və ya forever daha çox), continue (cari elementi ötür —
happy path sola hizalanır, `if !e.IsCurrent { continue }`), break (loop-dən
tam çıxış, funksiya işi varsa return-dən fərqli), labels (nested loop-larda
`break outer` / `continue outer` — NADİR, refaktor siqnalı), goto (AVOID! —
oxunaqlıq düşməni).

## Əsas fikirlər

### 1. switch — N Yoldan Biri
if/else zənciri böyüyəndə verbose olur. switch:
```go
switch {
case x < 0:
    fmt.Println("negative")
case x > 0:
    fmt.Println("positive")
default:
    fmt.Println("zero")
}
```
- `case ŞƏRT:` — şərtli formalarda hər case öz ifadəsini gətirir
- **İCRA QAYDASI:** case-lər SIRAYLA qiymətlənir; İLK true icra olunur və
  switch BİTİR — digər match-lər YOXLANMIR (yalnız bir case!)
- `default:` — heç biri tutmayanda; MÜMKÜNSA HƏMİŞƏ YAZ (daha aydın, etibarlı)
- **Dizayn qaydası:** case-lər ƏLAQƏLİ olmalıdır — unrelated if-ləri "switch-də
  gizlətmək" confusing koddur

### 2. fallthrough — Nadir İstisna
Default: ilk match → bit. `fallthrough` case-in sonunda → NÖVBƏTİ case də
icra olunar. Rəngarif/ifadə — varlığını bil, amma gümanən heç vaxt işlətməyəcəksən.

### 3. Switch Expression — Dəyərə Görə
```go
switch x {
case 1:
    fmt.Println("one")
case 2:
    fmt.Println("two")
case 3:
    fmt.Println("three")
}
```
- "switches on x" — case-lər = ifadənin MÜMKÜN DƏYƏRLƏRİ
- İstənilən ifadə ola bilər — case dəyərləri ilə EYNİ TİP olmalıdır (1,2,3 → int)
- **Çoxdəyərli case:**
```go
case 1, 2, 3:
    fmt.Println("one, two, or three")
```
- greet məşqi: Alice → "Hey, Alice.", Bob → "What's up, Bob?", default →
  "Hello, stranger."

### 4. break — Case-dən Tez Çıxış
```go
switch x {
case 1:
    if SomethingWentWrong() {
        break          // case-i bitir, switch-dən SONRAKİ koda keç
    }
    ...               // davam etmə
}
```
Metafora: məhbus cəzasının qalanını çəkmək istəmir — "breaks out".

### 5. for — Loop Açar Sözü
Təkrar = loop; Go-da YALNIZ `for` (while YOXDUR — dil sadəliyi):
```go
for x < 10 {          // conditional: şərt true olsa təkrar
    ...
}
```
**Mexanika:** şərt yoxla → true: body → SONDA YENİDƏN şərt → ... şərt false
olana qədər.

### 6. Forever Loop — `for {}`
Şərtsiz for = `for true` = SONSUZ. Faydalıdır!:
- İstehsal prosesləri: qaynəmə robotu — konveyer işlədiyi müddətəcə eyni
  qaynəkləri təkrarlayır
- Şəbəkə serverləri (HTTP): sonsuz "gözlə → serve et → gözlə" dövrü
Dayandırma mexanizmləri VAR (daha sonra), amma loop şərti kimi yazılmır.

### 7. range — Collection üzrə Loop
```go
for range employees {              // hər element üçün 1 dəfə (data YOX)
    fmt.Println("Found another employee!")
}

for i, e := range employees {      // index + element
    fmt.Println("Employee number %d: %v", i, e)
}

for _, e := range employees {     // yalnız element (index maraqsız)
for i := range employees {        // yalnız INDEX (element qəbul edilmir)
```
- **Qısa qeyd:** 3 element → 3 çap. Payroll: hər işçiyə payslip.
- String üzrə — RUNE-lar; map üzrə — key-value cütləri (hamısı range ilə!)
- total məşqi: `[]int{1,2,3}` → 6 (sum loop)

### 8. 3-Hissəli for — init; cond; post
```go
x := 0
for x < 10 {
    fmt.Println(x)
    x++
}
// eyni, kompakta:
for x := 0; x < 10; x++ {
    fmt.Println(x)
}
```
- **init** (x := 0) — şərtdən ƏVVƏL, bir dəfə
- **cond** (x < 10) — hər iterasiya öncəsi
- **post** (x++) — hər body-dən SONRA
- Nokta-vergüllə ayrılır; init/post scope = for statement
- `for x < 10` = init/post-sız xüsusi hal
- **Go-da NADİRDİR:** adətən ya range (data üzrə) ya da forever lazımdır
- evens məşqi: 0-100 cütlər 3-hissəli for ilə

### 9. continue — Növbəti Element
Problem: yalnız BAZI elementlər üzərində iş:
```go
for _, e := range employees {
    if e.IsCurrent {        // iç-içe if — ƏLAVƏ İNDENT!
        e.PrintCheck()
    }
}
```
Happy path prinsipinə zidd. Həll — FLIP + continue:
```go
for _, e := range employees {
    if !e.IsCurrent {
        continue            // bu elementi AT, növbətinə keç
    }
    e.PrintCheck()          // happy path SOLA düz xətt!
}
```
Bir neçə skip səbəbi ola bilər — hər biri öz continue-u. (nonNegative məşqi:
yalnız mənfidən kənarlərı çap et.)

### 10. break — Loop-dan Tam Çıxış
Funksiya return etsə — loopdan da çıxar, amma loopdan SONRA funksiya işi
qalırsa return YARAMAZ. break:
```go
for _, e := range employees {
    if MoneyLeft() <= 0 {
        fmt.Println("Oops, out of cash!")
        break               // loop BİTİR, for-un } sonrakı koda keç
    }
    ...  // print check
}
```

### 11. Labels — Nested Loop İdarəsi
```go
outer:
for x := 0; x < 10; x++ {
    for y := 0; y < 10; y++ {
        fmt.Println(x, y)
        if y == 5 {
            break outer     // HƏR İKİ loop-dan çıx!
        }
    }
}
```
- Label = kod yerinə AD (`outer:`); `break outer` / `continue outer`
- Tək break = YALNIZ daxili loop (outer davam edər)
- **NADİRDİR:** nested loop-lar onsuz da confusing → qaçınırıq; "jumping"
  axını anlaşılırlığı azaldır → label ehtiyacı = REFAKTOR SİQNALI
- **goto:** label-ə birbaşa eniş — break/continue-dan da pis! Yaxşı səbəb
  olmadan İSTİFADƏ ETMƏ — "clear, readable code" düşməni

## Əsas terminlələr
- switch / case / default — N yollu şərti seçim
- İlk-Match Qaydası — yalnız İLK tutan case icra olunur
- fallthrough — növbəti case-i də icra et (nadir)
- Switch Expression — `switch x` + dəyər case-ləri; "switches on x"
- Comma Case — `case 1, 2, 3:` — hər hansı bir uyğunluq
- break (switch) — cari case-dən tez çıxış
- for — Go-nun yeganə loop açar sözü (while YOX)
- Conditional Loop — `for x < 10`
- Forever Loop — `for {}` / `for true` — serverlər, robotlar
- range — collection elementləri üzrə təkrar
- for i, e := range — index + element cütü
- for _, e := range — element tək (index atılır)
- for i := range — index tək
- Init / Post Statement — for x := 0; cond; x++
- 3-Hissəli for — Go-da nadir (range/forever üstdür)
- continue — cari elementi ötür, növbəti
- break (loop) — loop-dan tam çıxış
- Nested Loop — loop içində loop
- Label (outer:) — kod yerinə ad; break/continue outer
- goto — birbaşa jump — QADAĞAYA YAXIN
- "Label = Refaktor Siqnalı" — nested-loop çətinliyi göstəricisi

## Praktik nəticə
(1) 3+ yollu seçim = if/else zənciri YOX, switch; default HƏMİŞƏ yaz. (2)
Case-lər eyni məsələyə aid olsun; unrelated şərtləri switch-ə yıxma. (3)
Dəyər seçimi: `switch x { case 1, 2, 3: }` — komma ilə qrup. (4) Go-da
loop = for: şərtli (`for x < 10`), sonsuz (`for {}` — server pattern), data
(`for _, v := range xs`). (5) 3-hissəli formanı tanı (`for i := 0; i < n; i++`)
amma ÖZÜNƏ sual ver: range daha yaxşıdır? (6) Element filtri: iç-içe if YOX —
`if !ok { continue }` + happy path SOLA. (7) Loop-dan çıxış: funksiya bitirsə
return; yoxsa break. (8) Nested + label görsən — dayan, REFAKTOR et; goto
isə heç vaxt. (9) range map-də key-value, string-də rune verir — universal
iteratordur. (10) x++ YALNIZ post mövqeyində təbii görünür; statementdir,
ifadə yoxdur.

## Mənbə
Pages: 150-160 (PDF 151-161)
