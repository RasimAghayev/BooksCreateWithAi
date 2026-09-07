# Chapter 11 — Opening statements (Açılış Statementləri)

## Bu fəsil nədən bəhs edir?

Control flow-ya giriş: statement anlayışı (instruction), assignment
(dəyər→"qutu", overwrite, ən son dəyər qalır), declaration (var — runtime
effekti YOX, "yer ayır"; declaration ≠ statement; funksiya xarici kod =
yalnız declaration), short declaration (`:=` — ":" = "yeni dəyəşən gəlir",
tip inference), tuple assignment (`a, b, c := 1, 2, 3` — əlaqəli qruplar üçün,
"lat, long := getPosition()", oxunaqlıq qaydası), blank identifier (`_ = "at",
"cannot use _ as value"; `_, ok := menu["eggs"]` — left-to-right sıra
məcburiyyəti), increment/decrement (STATEMENT ifadə DEYİL — `x := y++` syntax
error!), if statement (şərt → blok; indent = icra yolu), happy path konsepsiyası
və REFACTORING (iç-içe if-lərdən çıxış — şərtləri FLIP et: `x <= 0 → return
false` əvvəl; "left-align the happy path" qaydası), else branch (if/else —
"bir və ya digəri, heç vaxt hər ikisi"), early return (Go-da else AZ — if
return-la bitirsə else lazım deyil), müqayisə operatorları (==, !=, <, >, <=,
>=), logical operatorlar (&& and, || inclusive or — "Go öyrən və ya Python"
vs "səhər meditasiya YOXUMA yuxu" eksklüziv nümunələr; ! not), bool
dəyişənləri (if ok — sadə şərt), compound if (`if STATEMENT; CONDITION` —
comma-ok və err pattern-ləri; SCOPE yalnız blok daxilində; shadowing təhlükəsi;
uzun compound = qır).

## Əsas fikirlər

### 1. Statement — Əsas Kərpic
Statement = "et" əmri: `b = Book{...}` — assignment. Statementlər funksiya
DAXİLİNDƏ olur; paket səviyyəsində YALNIZ declaration-lar.

### 2. Declaration vs Assignment
```go
var b Book      // DECLARATION: runtime-da heç nə ETMİR; compiler-ə "yer ayır"
b = Book{...}   // ASSIGNMENT: dəyəri qutuya qoy; köhnəni OVERWRITE edir
name = "Kim"; name = "Jo"   // → "Jo" (sonuncu qalır, istənilən sayda)
```
**Dəyişən = qutu:** assignment = qutuya dəyər qoymaq; əvvəlki varsa əvəz
olunur. Declaration texniki olaraq statement DEYİL (runtime effekti yoxdur).

### 3. Short Declaration — `:=`
```go
b := Book{Title: "The Making of a Butterfly"}
```
":" = "yeni dəyişən gəlir — elə bil declared olub". Tip inference: sağ tərəfin
tipi = dəyişənin tipi (Book literal → b Book).

### 4. Tuple Assignment
```go
a, b, c := 1, 2, 3           // 3 ayrı statement əvəzinə 1
lat, long := getPosition()    // çoxdəyərli funksiya — ƏSAS istifadə
price, ok := menu["eggs"]     // map lookup da çoxdəyərlidir
```
**Adab qaydası:** tuple = ƏLAQƏLİ dəyərlər qrupu ("bunlar eyni şeydir"
mesajı); əlaqəsiz dəyişənləri birleşdirmə. Uzun/mürəkkəb tuple = oxunuşu
zəifleşdirir → parçala. Aydın olmırsa — ayrı-ayrı yaz.

### 5. Blank Identifier — `_`
```go
_ = "Hello"
fmt.Println(_)     // compile error: cannot use _ as value
```
"İçinə qoyulur, çıxarılmır." Niyə lazım? Çoxdəyərli təyinatda BAZI dəyərlər
maraqsızdır:
```go
price, ok := menu["eggs"]
// istifadə olunmayan price → "declared but not used" compile xətası!

ok := menu["eggs"]    // YANLIŞ HƏLL: soldan-sağa — ok PRICE-i alacaq!

_, ok := menu["eggs"] // DOĞRU: "bu dəyəri AT"
_, _, z := getPosition()  // təkrar istifadə MÜMKÜNDÜR
```

### 6. ++ / -- — Statement, İfadə DEYİL
```go
x++
x := y++
// syntax error: unexpected ++ at end of statement
```
Dəyəri YOXDUR — təyinata qoyulmaz. (C-dən fərqli, tələ unsurudur.)

### 7. if Statement + Happy Path Refaktoru
```go
if x > 0 {
    fmt.Println("x is positive")
}
```
Şərt true → blok; false →növbəti statement. **Happy path** = "hər şey
düzgündürsə" normal yol. Validasiya nümunəsi — İÇ-İÇƏ pis versiya:
```go
if x > 0 {                    // 2 səviyyə indentasiya...
    if x%2 == 0 {             // "positive AND even" sağa çox uzaq
        return true
    }
    return false
}
return false
```
**REFACTORING (davranışı dəyişmədən təkmilləşdirmə): şərtləri FLIP et:**
```go
if x <= 0 {           // mümkünsüz hallar ƏVVƏL, early return
    return false
}
if x%2 != 0 {
    return false
}
return true           // happy path — düz xətt sola
```
**Qayda: "Left-align the happy path"** — minimum indent, vizual gözlənilən
axın; xətalar indented bloklarda ayrı görünür. (% — remainder operatoru;
x%2 == 0 = cüt.)

### 8. else — "Bir və ya Digəri"
```go
if x <= 0 {
    fmt.Println("x is zero or negative")
} else {
    fmt.Println("x is positive")
}
```
if/else = YA biri YA digəri (heç vaxt hər ikisi). **Amma Go-da else AZ
istifadə olunur** — early return onu əvəz edir (aşağıda).

### 9. Early Return — else-siz Go Tərzi
```go
if x <= 0 {
    return false      // blok return-la bitirsə...
}
return true           // ...buraya çatdıqsa şərt FALSE idi — aydındır!
```
"If bloku returnla bitirsə, else LAZIM DEYİL" — funksiya daha qısa, daha
oxunaqlı.

### 10. Conditional + Logical Operatorlar
Müqayisə: `== != < > <= >=`. Birləşdirmə:
```go
if x > 0 && x%2 == 0 { ... }     // && = AND — hər ikisi
if x > 0 || x%2 == 0 { ... }     // || = İNCLUSİVE OR
if !ok { ... }                   // ! = NOT
```
**|| dəqiqləşdirmə (inklüzivdir!):** true, əgər a YAXUD b YAXUD HƏR İKİSİ;
false YALNIZ hər ikisi false. İngilis "or" amphibolikdir: "Go və ya Python
öyrən" (inklüziv — hər ikisini öyrən!), "səhər dur yoga YOXUMSA yataqda qal"
(eksklüziv — astral səyahət ustası deyilsən!). Go-da || HƏMİŞƏ inklüzivdir.
**XOR məşqi:** `xor(a, b)` — dəqiq BİRİ true → true; hər ikisi → false.

### 11. bool Dəyişənləri Sadə Şərtlərdir
```go
_, ok := menu["eggs"]
if ok {                  // if ok == true YAZMA — artıq!
    fmt.Println("Eggs are on the menu!")
}
```
`if true {...}` / `if false {...}` — mənasız: həmişə icra / heç vaxt → çıxart.

### 12. Compound If — `if STATEMENT; CONDITION`
```go
if _, ok := menu["eggs"]; ok {
    fmt.Println("Eggs are on the menu!")
}
```
= assignment + condition bir if-də. **Sirkulyasiya:** ok YALNIZ bu blokda
mövcuddur (scope):
```go
fmt.Println(ok)    // error: undefined: ok — blokdan çıxıb YOXDUR
```
Ən məşhur forma:
```go
if err := doStuff(); err != nil {
    fmt.Println("oh no")
}
```
err scope-bloka bağlı — paket boyu err "sızıntısı" YOX.

**Shadowing təhlükəsi:** xarici err varsa, compound-if içindəki YENİ err onu
ÖRTÜR (shadow) — blokda doStuff-un dəyəri, blokdan sonra xarici err ORİJİNAL
qalır. Go çaşmır (aydın scope qaydaları), BİZ çaşırıq — "xarici dəyişəni
yenilədiyimi düşünüb shadow etmək" = klassik bug. **Qayda: shadowing-dən
QAÇIN.**

**Uzun compound = oxunmaz:**
```go
if err := apply(func(x int) error {...uzun...}, -1); err != nil { ... }
// }, -1); err != nil { → oxuyan çaşır
```
**Qayda:** compound QISA və tək sətirdirsə — OK; uzun/dəqiq deyilsə — PARÇALA.

## Əsas terminlələr
- Statement — "et" əmri; funksiya daxilində
- Declaration — var; runtime effekti yox; yalnız paket səviyyəsində icazə
- Assignment — dəyəri qutuya qoy; overwrite = sonuncu qalır
- Qutu Metaforası — dəyişən = tipli qutu
- Short Declaration (:=) — declare+assign; tip inference sağdan
- Tuple Assignment — a, b, c := 1, 2, 3 — əlaqəli qrup
- Left-to-Right Sıralama — çoxdəyərli təyinatın icra qaydası
- Blank Identifier (_) — "at"; təkrar istifadə OK; oxunmaz
- Increment/Decrement — statement; ifadə DEYİL
- Control Flow — icra axını; sadə/şərti/tsiklik
- Happy Path — hər şey normal olandakı axın
- Refactoring — davranışı saxlayaraq təkmilləşdirmə
- Flip Şərtlər — x > 0 → x <= 0; mümkünsüzlər əvvəl
- "Left-Align the Happy Path" — minimum indent prinsipi
- else — bir və ya digəri (hər ikisi heç vaxt)
- Early Return — if + return = else-sizlik
- && / || / ! — and / inklüziv or / not
- Inclusive vs Exclusive Or — hər ikisi mümkün / mümkün deyil
- Remainder Operator (%) — x%2 == 0 = cüt
- bool Şərt Sadəliyi — if ok; if ok == true = artıq
- Compound If — if stmt; cond — scope bloka bağlı
- Shadowing — daxili dəyişənin xaricini örtməsi — QAÇIN
- Compound Uzunluq Qaydası — qısa/tək sətir OK, uzun = parçala

## Praktik nəticə
(1) `:=` yeni dəyişən üçün; `=` mövcuduna. (2) Tuple yalnız ƏLAQƏLİ
dəyərlərə; uzunsa parçala. (3) Maraqsız dəyərlər `_`-ə; left-to-right sıranı
unutma (ok price-i almasın!). (4) `x := y++` MÜMKÜNSÜZ — ++ statementdir.
(5) Validasiya/şərtlilikdə: mümkünsüz halları FLIP edib ƏVVƏLƏ qoy, early
return ilə qayıt — happy path SOLA hizalanmış düz xətt olsun. (6) else-dən
əvvəl early return düşün — Go-da else nadirdir. (7) `||` həmişə inklüziv —
"hər ikisi" halını unutma; XOR lazımdırsa öz funksiyan. (8) `if ok` yaz,
`if ok == true` YOX. (9) `if err := f(); err != nil` — əsas pattern; err scope
bloka bağlı, lakin shadowing-i izlə — eyni adlı xarici dəyişən risklidir. (10)
Compound if yalnız QISA ikən; uzunsa ayrı statement-lərə böl — oxunaqlılıq
hər şeydən üstündür.

## Mənbə
Pages: 138-149 (PDF 139-150)
