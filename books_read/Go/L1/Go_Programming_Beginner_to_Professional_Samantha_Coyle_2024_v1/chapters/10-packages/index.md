# Chapter 10 — Packages Keep Projects Manageable (səh. 342-369)

## Bu fəsil nədən bəhs edir?

Paketlər: struktur (qovluq + .go faylları), adlandırma konvensiyaları,
package declaration, exported/unexported (böyük/kiçik hərf), alias-lar,
main paketinin xüsusiyyəti, init() funksiyası (icra sırası, çoxlu init)
və layihənin paketlərə bölünməsi (cmd/ + pkg/ strukturu).

## Əsas fikirlər

### 1. Paket nədir?
**Nədir:** bir qovluqda yaşayan bir və ya bir neçə .go faylından ibarət
kod qrupu — əlaqəli kodu birləşdirir, yalnız lazımi hissələri açır.

**Kod təşkilatının inkişafı:** funksiyalar → fayllar → PAKETLƏR
(strings paketi nümunəsi: builder.go, compare.go, replace.go... hər
fayl öz funksiyası ilə).

**Üç sütun:**
- **Maintainable** — dəyişiklik asan, risk az (SDLC boyu dəyişiklik
  xərci artır)
- **Reusable** — gələcək layihələrdə istifadə; xərc/vaxt azalır, keyfiyyət
  artır
- **Modular** — hər vəzifə öz yerində; axtarış asan

**DRY (Don't Repeat Yourself)** — paketlər DRY-nin modul səviyyəsidir.

### 2. Paket adlandırması
| Pis | Yaxşı (standart kitabxana!) |
|---|---|
| stringconversion | strconv |
| synchronizationprimitives | sync |
| measuringtime | time |
| StringConversion, measuring_time | — (stil pozuntusu) |

**Qaydalar:**
- Tam kiçik hərf, alt xətt YOX, camelCase YOX
- Qısa + konkret isim (adətən təkhallı)
- Abreviaturalar icazəli (regexp, strconv) — tanınmışdırsa
- QAÇIN: misc, util, common, data — məqsəd bildirmir
- Ad = self-documentation

### 3. Package declaration
```go
package <packageName>     // HƏR faylın İLK sətri
```
- Eyni paketin bütün faylları eyni deklarasiyaya başlayır
- strings paketinin real nümunəsi: builder.go və compare.go hər ikisi
  `package strings` ilə başlayır
- Paket daxilində bütün kod QARŞILIQLI görünür (fayllararası!)

### 4. Exported vs Unexported
**Qayda:** BÖYÜK hərf = exported (xaricdən görünür); kiçik hərf =
unexported (yalnız paket daxilində). Access modifier YOXDUR!

```go
// strings paketindən:
func Contains(s, substr string) bool { ... }   // C böyük → EXPORTED
func explode(s string, n int) []string { ... }  // e kiçik → unexported

// İstifadə:
strings.Contains(str, "found")    // ✓ işləyir
strings.explode(str, 3)          // ✗ COMPILE XƏTASI
```
- Praktika: yalnız nə lazımdırsa aç; qalanını gizlət

### 5. Package alias
```go
import f "fmt"                  // f = alias

func main() {
    f.Println("Hello, Gophers") // alias ilə çağırış
}
```
Səbəblər: ad məqsədizdirsə, çox uzundursa, iki eyni adlı paket varsa.

### 6. main paketi
- EXECUTABLE paket — binary hasil edir; adı = qovluq adı
- `main()` funksiyası MƏCBURİDİR — giriş nöqtəsi
- main-in məntiqi başqa paketlər tərəfindən istifadə OLUNMUR
- main OLMAYAN paket = non-executable (kitabxana)

### 7. init() funksiyası
**Nədir:** paket başlanğıcı — main-dən ƏVVƏL icra olunan setup.

**İcra sırası:**
1. İmport olunan paketlər (rekursiv)
2. Paket səviyyəli dəyişənlər (`var name = "Gopher"`)
3. init() funksiyaları
4. main()

```go
var name = "Gopher"            // 1) dəyişən

func init() {                   // 2) init
    fmt.Println("Hello,", name)
}

func main() {                   // 3) main
    fmt.Println("Hello, main function")
}
// Çıxış: Hello, Gopher → Hello, main function
```

**init qaydaları:**
- Arqumentsiz, qaytarmasız (`func init(age int)` → compile xətası)
- İstifadə: DB bağlantıları, paket dəyişənləri, konfiq yükləmə, vəziyyət
  yoxlaması

**Çoxlu init():** paketdə bir neçə init ola bilər — KODDAKI SIRAYLA
icra olunur:
```go
func init() { fmt.Println("Hello, Gopher") }   // 1-ci
func init() { fmt.Println("Second") }           // 2-ci
func init() { fmt.Println("Third") }            // 3-cü
func main() { ... }                              // ən sonda
```

**Budget categories nümunəsi (Exercise 10.02-10.03):**
```go
var budgetCategories = make(map[int]string)
var payeeToCategory = make(map[string]int)

func init() {                          // kateqoriyaları yüklə
    fmt.Println("Initializing our budgetCategories")
    budgetCategories[1] = "Car Insurance"
    budgetCategories[2] = "Mortgage"
    budgetCategories[3] = "Electricity"
    // ...
}

func init() {                          // payee-ləri qoş
    payeeToCategory["State Farm"] = 1
    // ...
}

func main() {
    for k, v := range payeeToCategory {
        fmt.Printf("Payee: %s, Category: %s\n", k, budgetCategories[v])
    }
}
```

### 8. Layihə strukturu (Exercise 10.01 — shape paketi)
```
Exercise10.01/
├── go.mod                 (module exercise10.01)
├── cmd/
│   └── main.go            (package main — İCRA)
└── pkg/shape/
    └── shape.go           (package shape — KİTABXANA)
```

**pkg/shape/shape.go:**
```go
package shape

import "fmt"

type Shape interface {
    area() float64          // unexported metod — paket daxili
    name() string
}

type Triangle struct {        // EXPORTED tiplər
    Base   float64            // EXPORTED sahələr
    Height float64
}
type Rectangle struct { Length, Width float64 }
type Square struct { Side float64 }

func (t Triangle) area() float64 { return (t.Base * t.Height) / 2 }
func (t Triangle) name() string { return "Triangle" }
// ... Rectangle, Square eynilə

func PrintShapeDetails(shapes ...Shape) {   // EXPORTED funksiya
    for _, item := range shapes {
        fmt.Printf("The area of %s is: %.2f\n", item.name(), item.area())
    }
}
```

**cmd/main.go:**
```go
package main

import "exercise10.01/pkg/shape"     // lokal paket importu

func main() {
    t := shape.Triangle{Base: 15.5, Height: 20.1}
    r := shape.Rectangle{Length: 20, Width: 10}
    s := shape.Square{Side: 10}
    shape.PrintShapeDetails(t, r, s)
}
```
- Paket adı = yolun SON qovluğu
- go build cmd qovluğunda → `cmd` binary

## Activity icmalı
- **10.01:** Chapter 7 Activity 7.01-i modullaşdır: pkg/payroll paketi
  (Developer/Employee/Manager tipləri + metodlar), tipləri fayllara
  ayır, main-də 2 init() (salamlama + dəyişən setup)

## Əsas terminlər
- Package — qovluqdakı .go faylları qrupu
- DRY — Don't Repeat Yourself
- Package files — paketin .go faylları
- Package declaration — ilk sətir `package X`
- Exported (böyük hərf) / Unexported (kiçik hərf)
- Access modifier — YOXDUR Go-da
- Package alias — `import f "fmt"`
- Executable (main) vs non-executable paket
- init() — paket başlanğıcı; import-lar → var-lar → init → main
- Çoxlu init — kod sırası ilə
- cmd/ + pkg/ strukturu — icra/kitabxana ayrılığı

## Praktik nətiəcə
Paketlər layihə böyüdükcə məcburi alətdir: adı qısa, kiçik hərf,
məqsədli (strconv kimi); misc/util/common YOX. Ailəvi kod bir qovluqda;
`cmd/` (main) + `pkg/` (kitabxana) — standart struktur. API sərhədi =
böyük hərf; daxili implementasiya = kiçik hərf. init() ilə
konfiq/kateqoriya yüklə (sıra: import → var → init kod sırası → main),
amma sadə qalması üçün init-ləri ayrı-ayrı vəzifələrə böl.

## Mənbə
Pages: 342-369 (PDF 342-369)
