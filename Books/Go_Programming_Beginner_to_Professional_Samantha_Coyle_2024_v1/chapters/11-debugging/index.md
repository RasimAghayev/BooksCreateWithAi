# Chapter 11 — Bug-Busting Debugging Skills (səh. 370-397)

## Bu fəsil nədən bəhs edir?

Debug metodologiyası: bug mənbələri, proaktiv tədbirlər (inkremental
kod+test, unit test, error idarəetməsi, logging), fmt format verbs
(%s/%d/%f/%t/%v, genişlik/dəqiqlik), kod markerləri, %T tip çapı, %#v Go
representasiyası, log paketi (SetFlags, Fatal/Panic fərqi) və məhdud
mühitlərdə debug (Delve, pprof, feature flags).

## Əsas fikirlər

### 1. Bug mənbələri (production-a çıxan)
1. Test sonrakı düşüncədir — funksiya bitən kimi test edilmir
2. Kod dəyişikliyi — tələb dəyişikliyi başqa sahəni pozur
3. Qeyri-real deadline — qısayollar, qısaldılmış dizayn/test
4. İdarə edilməyən xətalar — fayl yoxdur, div/0, bağlantı kəsilib

**Real nümunə:** heparin dozalaşdırma maşını — hesablama bug-u həyatı
təhdid edir.

### 2. Proaktiv metodlar
- **Inkremental kodlaşdırma + tez-tez test** — kiçik hissə bitən kimi
  sına; bug izləmə sahəsi daralır
- **Unit testlər** — giriş→çıxış yoxlaması; dəyişiklik sonrası test
  düşməsi = yeni bug siqnalı
- **BÜTÜN xətaları idarə et** (ch6)
- **Logging** — debug/info/warn/error/fatal/trace səviyyələri;
  performans təsirini nəzərə al (həddindən artıq log = yavaşlama)

### 3. fmt ilə format
```go
// Println — aralarında boşluq, sonda \n:
fmt.Println("Hello:", fname, lname)   // Hello: Edward Scissorhands

// Printf — verbs (C-dən götürülüb):
fmt.Printf("Hello Mr. %s %s", fname, lname)
// Birdən çox verb — ardıcıllıqla əvəz olunur

// Yeni sətir AVTOMATİK YOXDUR:
fmt.Printf("Hello my first name is %s\n", fname)
```

**Əsas verbs:**
| Verb | Təsvir |
|---|---|
| %s | string |
| %d | int (10-luq) |
| %f | float (default 6 kəsr) |
| %t | bool |
| %v | dəyər (default format) |
| %T | TİP adı |
| %#v | Go sintaksis representasiyası |

**Dəqiqlik + genişlik:**
```go
fmt.Printf("%s has a gpa of %.2f.\n", fname, gpa)  // 2 kəsr: 3.75

// Genişlik.dəqiqlik — sağa düzləmə, boşluq doldurma:
fmt.Printf("%10.0f\n", v)   //      1234
fmt.Printf("%10.1f\n", v1)  //     1234.6
fmt.Printf("%10.2f\n", v2)  //    1234.67

// Sola düzləmə — - flaq:
fmt.Printf("%-10.0f\n", v)  // 1234
```

**Sayı sistemləri təlimi (Exercise 11.02):**
```go
for i := 1; i <= 255; i++ {
    fmt.Printf("Decimal: %3.d Base Two: %8.b Hex: %2.x\n", i, i, i)
}
```

### 4. Kod markerləri (bug lokallaşdırma)
**Nədir:** proqramın hansı nöqtəsində olduğumuzu göstərən print-lər —
bug-un YERİNİ daraldır.

```go
func a(i int) error {
    if i < 10 {
        fmt.Println("Error is in func a")      // MARKER
        return errors.New("Incorrect value")
    }
    return nil
}

func b(i int) error {
    if i >= 10 {
        fmt.Println("Error is in func b")       // MARKER
        return errors.New("Incorrect value")
    }
    return nil
}

// Çıxış: "Error is in func a" → bug HARADADIR bilinir
```
- Marker tapılmayanda başqa yerə köçür; iş bitəndə sil

### 5. Tip çapı — %T
```go
type person struct {
    lname string
    age    int
    salary float64
}

fname := "Joe"
grades := []int{100, 87, 67}
states := map[string]string{"KY": "Kentucky"}
p := person{lname: "Lincoln", age: 210, salary: 25000.00}

fmt.Printf("fname is of type %T\n", fname)   // string
fmt.Printf("grades is of type %T\n", grades)   // []int
fmt.Printf("states is of type %T\n", states)  // map[string]string
fmt.Printf("p is of type %T\n", p)            // main.person
```
- %T (böyük) = tip; %t (kiçik) = bool — hərflər həssasdır!

### 6. %#v — Go representasiyası
```go
fmt.Printf("fname value %#v\n", fname)   // "Joe"
fmt.Printf("grades value %#v\n", grades) // []int{100, 87, 67}
fmt.Printf("states value %#v\n", states) // map[KY:Kentucky...]
fmt.Printf("p value %#v\n", p)            // main.person{...}
```
- Kopyala-yapışdır edilə bilən Go sintaksisi — dəyəri Go "gördüyü"
  kimi göstərir

### 7. log paketi
```go
name := "Thanos"
log.Println("Demo app")          // 2019/11/10 23:00:00 Demo app
log.Printf("%s is here!", name)  // vaxt möhürü avtomatik
log.Print("Run")
```

**SetFlags — ətraflı məlumat:**
```go
log.SetFlags(log.Ldate | log.Lmicroseconds | log.Llongfile)
// Ldate — tarix; Lmicroseconds — mikrosoniyələr;
// Llongfile — TAM fayl yolu + sətir nömrəsi
// Digərlər: Ltime, Lshortfile, LUTC...
```
Çıxış: `2024/02/12 07:09:14.015902 /path/main.go:21: mesaj`

### 8. Fatal vs Panic
```go
log.Fatalln(err)    // logla + os.Exit(1) — defer-lər ÇAĞIRMIR
log.Panicln(err)    // logla + panic — defer/recover İŞLƏYİR
```
- **Fatal:** bərpa mümkün deyil (data korlanması təhlükəsi, exit kodu
  lazımdır)
- **Panic:** bərpa mümkündür

```go
func main() {
    log.SetFlags(log.Ldate | log.Lmicroseconds | log.Llongfile)
    log.Println("Start of our app")
    err := errors.New("Application Aborted!")
    if err != nil {
        log.Fatalln(err)             // burada çıxır
    }
    log.Println("End of our app")     // İCRA OLUNMUR
}
```

### 9. Məhdud/canlı mühitlərdə debug
- **Mühiti başa düş** — deployment, şəbəkə, təhlükəsizlik məhdudiyyəti
- **Delve remote debug** — canlı prosesə qoşul, breakpoint, addımlama
- **pprof observability** — profil endpoint-ləri; kodu dəyişmədən
  runtime statistikası; metrics + log aqreqasiyası
- **Log səviyyələri** — amma həssas datanı logda ifşa etmə
- **IDE debuggers** (VS Code, GoLand) — breakpoint/watch; development-də
  effektiv, production-da mümkün deyil
- **Feature flags / canary** — funksionallığı seçici aktiv et; kiçik
  istifadəçi qrupunda müşahidə et
- Klassik dərslər: "works on my machine" — CI-də kiçik inkrementlərlə
  sına; planlı insident məşqləri (game days) komandanı hazırlayır

## Activity icmalı
- **11.01 SSN validator:** ErrInvalidSSNLength/Numbers/Prefix/DigitPlace
  xətaları; 4 validasiya funksiyası; log ilə İZLƏMƏ (vaxt+fayl+sətir);
  xəta baş verəndə davam et (Fatal YOX)

## Əsas terminlər
- Bug/debugging — istənməyən davranış / səbəbin tapılması
- Incremental coding — kiçik hissələrlə kod + tez test
- Unit test — giriş→çıxış müqayisəsi
- Format verb (%s/%d/%f/%t/%v/%T/%#v)
- Width & precision (%10.2f, %-10s)
- Kod markeri — yer göstəricisi print
- log.SetFlags (Ldate|Lmicroseconds|Llongfile)
- Fatal vs Panic — exit(1) vs recoverable
- Delve — Go debugger (remote də dəstəkləyir)
- pprof — runtime profiling
- Feature flag / canary release — seçici aktivləşdirmə

## Praktik nəticə
Proaktiv ol: inkremental yaz + hər hissəni test et + bütün error-ları
idarə et + həmişə logging. Debug üçün alət dəsti: markerlərlə sahəni
daralt → %T tipini gör → %#v dəyəri Go-görməsində çap et → log.Printf
izləmə (SetFlags ilə fayl+sətir). Fatal yalnız bərpaolunmaz halda (defer
çağırılmır!); bərpa lazımdırsa Panic. Production-da: Delve attach,
pprof endpoint, aqreqasiya olunan loglar, feature flag-lərlə tədricən
rollout.

## Mənbə
Pages: 370-397 (PDF 370-397)
