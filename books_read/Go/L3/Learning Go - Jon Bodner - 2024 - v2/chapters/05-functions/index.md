# Chapter 5 — Functions (Funksiyalar)

## Bu chapter nədən bəhs edir?

Funksiya bəyanatı və çağırışı, variadic parametrlər, çoxlu qaytarma dəyərləri, named
return values (və blank return təhlükəsi), funksiyaların dəyər olması, anonim
funksiyalar və closure-lar, `defer` və call-by-value semantikası.

## Əsas fikirlər

### 1. Funksiya Bəyanatı
**Nədir:** `func` + ad + input parametrlər + qaytarma tipi + gövdə.

**Kitabdan kod nümunəsi:**
```go
func div(numerator int, denominator int) int {
    if denominator == 0 {
        return 0
    }
    return numerator / denominator
}
// Eyni tipli parametrlər birleşə bilər:
func div(numerator, denominator int) int { ... }
```

**Qeydlər:** Named/optional parametr YOXDUR — emulyasiya üçün struct parametr keçir
(`MyFunc(MyFuncOpts{LastName: "Patel", Age: 50})`); amma funksiya çox parametrli
alırsa — o, çox mürəkkəbdir, refactor et.

### 2. Variadic Parametrlər
**Nədir:** `...T` — son (və ya tək) parametr olmaqla istənilən sayda dəyər qəbul edir.

**Necə işləyir:** Funksiya daxilində parametr `[]T` slice-ı kimi davranır. Slice-ı ötürmək
üçün `...` suffix-i məcburidir (`addTo(3, a...)`), yoxsa compile xətası.

**Kitabdan kod nümunəsi:**
```go
func addTo(base int, vals ...int) []int {
    out := make([]int, 0, len(vals))   // zero-len + nonzero-cap pattern
    for _, v := range vals {
        out = append(out, base+v)
    }
    return out
}
addTo(3)                       // []
addTo(3, 2, 4, 6, 8)          // [5 7 9 11]
addTo(3, []int{1,2,3,4,5}...) // [4 5 6 7 8]
```

### 3. Çoxlu Qaytarma Dəyərləri
**Nədir:** Qaytarma tipləri mötərizədə vergüllə; hamısı qaytarılmalıdır.

**Necə işləyir:** Python tuple destructuring-dən fərqli olaraq Go-da **həqiqi ayrı
dəyərlərdir** — tək dəyişənə mənimsətmək compile xətasıdır. Konvensiya: `error` həmişə
son qaytarma dəyəridir. İstifadə edilməyən dəyərlər `_`-a mənimsədilir; hamısını
susmayaq (drop) olar (məs. `fmt.Println` 2 dəyər qaytarır), amma açıq `_` idiomatikdir.

**Kitabdan kod nümunəsi:**
```go
func divAndRemainder(numerator, denominator int) (int, int, error) {
    if denominator == 0 {
        return 0, 0, errors.New("cannot divide by zero")
    }
    return numerator / denominator, numerator % denominator, nil
}
result, _, err := divAndRemainder(5, 2)
```

### 4. Named Return Values
**Nədir:** Qaytarma dəyərlərinin mümkün adları — funksiya daxilində zero value ilə
əvvəlcədən bəyan olunmuş dəyişənlər.

**Tələlər:**
- Ad yalnız funksiya daxilində etibarlıdır; çağıran tərəf dəyişənlərin adına bağlı deyil.
- Kölgələnə bilər — təyinatı named return-a yox, kölgəyə vermək asandır.
- **Məcburiyyət YOXDUR**: `result, remainder = 20, 30` təyin edib sonra `return 2, 1, nil`
  yazsan — return-dakı dəyərlər qayıdır; compiler return ifadələrini adlı parametrlərə
  mənimsədir. Qarışıqlıq mənbəyidir.
- Bəzi dəyərləri adsız saxlamaq: `_`.

**Kitabdan kod nümunəsi:**
```go
func divAndRemainder(numerator, denominator int) (result int, remainder int, err error) {
    if denominator == 0 {
        err = errors.New("cannot divide by zero")
        return // blank return
    }
    result, remainder = numerator/denominator, numerator%denominator
    return // blank return
}
```

### 5. Blank (Naked) Returns — HEÇ VAXT İSTİFADƏ ETMƏ
**Nədir:** Adlı return-larda dəyərsiz `return` — son mənimsədilmiş dəyərləri qaytarır.

**Çatışmamazlıqları:** Oxucu return-ın nə qaytardığını görmək üçün bütün funksiyanı
scan etməli olur — data flow oxunmaz olur. Kitabın QADAĞASI: dəyər qaytaran funksiya üçün
blank return yazma. Adlı return-in yeganə tələbi: `defer` daxilində return dəyərlərini
oxumaq/dəyişmək (aşağıda).

### 6. Funksiyalar Dəyərdir — Closure-lar
**Nədir:** Funksiyanın tipi = imza (`func(int, int) int`); dəyişənlərə, map-lərə,
parametrə mənimsədilə və qaytarıla bilər. Funksiya daxilindəki funksiyalar closure-dur —
xarici funksiyanın dəyişənlərini görə və dəyişə bilər.

**Kitabdan kod nümunəsi:**
```go
var opMap = map[string]func(int, int) int{
    "+": add, "-": sub, "*": mul, "/": div,
}
opFunc, ok := opMap[op]
result := opFunc(p1, p2)
```

**Funksiya tipi bəyanı:** `type opFuncType func(int, int) int` — sənədləşdirmə və təkrar
istifadə üçün.

**Anonim funksiyalar:** `func(j int) { ... }(i)` — birbaşa çağırılış normal deyil, amma
`defer` və goroutine-lanç üçün əsasdır.

**Closure-ların istifadəsi:**
- `sort.Slice(people, func(i, j int) bool { return people[i].LastName < people[j].LastName })`
  — `people` closure tərəfindən capture olunur; eyni data müxtəlif sahələr üzrə sortlanır.
- Funksiya qaytaran funksiya:
```go
func makeMult(base int) func(int) int {
    return func(factor int) int {
        return base * factor
    }
}
```
- Bu pattern-lər: `sort.Search`, web middleware (sonrakı chapter), resource cleanup (defer).

**Kitabın prinsipi:** "Don't write fragile programs" — calculator nümunəsində 22 sətirlik
loop-un 16-sı error-checkdir; professional fərqi məhz error handling-dədir.

### 7. `defer` — Təmizləmə Zəmanəti
**Nədir:** Əhatə edən funksiya çıxanda icra olunacaq çağırışı qeydiyyata alan keyword.

**Necə işləyir:**
- Funksiyanın neçə çıxış nöqtəsi olmasından asılı olmayaraq icra olunur (return-dan sonra).
- Çoxlu defer-lər **LIFO** sıra ilə: son qeydiyyat — ilk icra.
- Parametrlər defer İCRA ZAMANI qiymətləndirilir, qeydiyyat anında YOX.
- defer üçün closure `()` ilə çağırılmalıdır (`defer f.Close()`, `defer closer()`).

**Kitabdan kod nümunəsi:**
```go
f, err := os.Open(os.Args[1])
if err != nil { log.Fatal(err) }
defer f.Close()   // bütün çıxış yollarında icra olunacaq
```

**Adlı return + defer pattern (yeganə vacib named return istifadəsi):**
```go
func DoSomeInserts(ctx context.Context, db *sql.DB, value1, value2 string) (err error) {
    tx, err := db.BeginTx(ctx, nil)
    if err != nil { return err }
    defer func() {
        if err == nil {
            err = tx.Commit()   // commit xətası err-i yeniləyir
        }
        if err != nil {
            tx.Rollback()
        }
    }()
    _, err = tx.ExecContext(ctx, "INSERT INTO FOO (val) values $1", value1)
    if err != nil { return err }
    return nil
}
```

**Resource + cleanup closure qaytarma pattern:**
```go
func getFile(name string) (*os.File, func(), error) {
    file, err := os.Open(name)
    if err != nil { return nil, nil, err }
    return file, func() { file.Close() }, nil
}
f, closer, err := getFile(os.Args[1])
defer closer()   // closer çağrılmadıqsa unused-var compile xətası — yaddaşköçürücü!
```

**Mədəniyyət qeydi:** try/catch/finally bloklarından fərqli olaraq defer ekstra indent
yaratmır — tədqiqatlar (Antinyan et al. 2017; Miara et al. 1983) nested dərinliyin
mürəkkəbliyin əsas amili olduğunu göstərib.

### 8. Go Is Call By Value (Dəyərlə çağırma)
**Nədir:** Parametr ötürüləndə həmişə dəyərin KOPYASI çatır.

**Necə işləyir:**
- int, string, struct — dəyişikliklər ÇAĞIRANA TƏSİR ETMİR (`modifyFails(i, s, p)` — heç
  nə dəyişmir; Java/JS/Python-dan fərqli olaraq struct sahələri də).
- **Map** — kopya pointer metadata olduğundan dəyişikliklər GÖRÜNÜR (element əlavə/silmə
  daxil).
- **Slice** — elementlərin modifikasiası görünür (elementlər arxa yaddaşda ümumi), amma
  `append` (uzunartma) görünMÜR — kopyanın len/cap metadata-sı dəyişir.

**Kitabdan kod nümunəsi:**
```go
func modMap(m map[int]string) {
    m[2] = "hello"; m[3] = "goodbye"; delete(m, 1)  // hamısı çağıran tərəfdə görünür
}
func modSlice(s []int) {
    for k, v := range s { s[k] = v * 2 }  // görünür
    s = append(s, 10)                      // görünMÜR
}
```

**Praktik fayda:** Funksiyalar input parametrlərini dəyişmədikcə data axını izlənilən
olur — bu, Go-nun const-ın zəifliyini kompensasiya edir. Mutable ötürmək lazımdırsa —
pointer lazımdır (növbəti chapter).

## Əsas terminlər
- Signature (imza) — parametr və qaytarma tiplərinin kombinasiyası
- Variadic (çoxsaylı) — `...T` ilə istənilən sayda arqument
- Named return value (adlı qaytarma) — bəyanatda adlandırılmış qaytarma dəyəri
- Blank/naked return (boş qaytarma) — dəyərsiz `return` (QADAĞAN)
- Closure (bağlanma) — xarici dəyişənləri capture edən daxili funksiya
- Higher-order function (yüksək dərəcəli funksiya) — funksiya qəbul/qaytaran funksiya
- Call by value (dəyərlə çağırma) — həmişə kopya ötürməsi

## Praktik nəticə

Professional Go funksiyası: (1) error həmişə son qaytarma; (2) çox parametr → struct
options; (3) adlı return-lardan yalnız defer əlaqəli hallarda istifadə; (4) resource
açılan kimi defer qeydiyyat; (5) funksiya input-ları dəyişməsin, yeni dəyər qaytarsın
(map/slice "xüsusi hallarını" bilsən belə); (6) closure-ları sort, middleware və
cleanup üçün işlət; (7) hər funksiya kiçik və tək-məqsədli olsun.

## Mənbə
Pages: 133-158
