# Chapter 8 — Packages (Paketlər)

## Bu fəsil nədən bəhs edir?

Paket anlayışı və strukturu (qovluq + .go faylları), adlandırma konvensiyası, package
declaration, exported/unexported (böyük/kiçik hərf), GOROOT/GOPATH (bin/pkg/src),
import yolları, paket aliasları, main paketi (executable vs non-executable), init()
funksiyası (icra sırası, çoxsaylı init, istifadə halları) və shape/payroll paketləşdirmə
activity-ləri.

## Əsas fikirlər

### 1. Paket Nəyə Lazımdır — 3 Sütun
- **Maintainable (saxlanıla bilən):** dəyişiklik asan + aşağı risk; rəqabət üçün
  vacib
- **Reusable (yenidən istifadə):** gələcək layihə xərci ↓, vaxt ↓, keyfiyyət ↑
  (çox test+istifadə), innovasiyaya vaxt ↑
- **Modular:** hər diskret tapşırıq öz YERİNDƏ — böyük kodda axtarış mümkünlüyü

**DRY prinsipi:** eyni kodu 2 dəfə YAZMA. İnkişaf zənciri: funksiyalar → fayllar →
PAKETLƏR. Paketlər olmadan hər proqrama kodu KOPYALAMAQ lazım gələrdi — bug düzəlişi
N proqramda.

### 2. Paket Strukturu
**Paket = qovluq + 1+ .go faylı (hamısı EYNİ paket declaration ilə).** Kod paket
daxilində fayllararası TAM GÖRÜNÜR; xaricdən yalnız EXPORTED.

**strings nümunəsi (standart kitabxana):**
```
strings/
├── builder.go    ← fayl adları FUNKSİYAYA görə
├── compare.go
├── reader.go
├── replace.go
├── search.go
└── strings.go
```
Struktur: paket → fayl → funksiya — hər səviyyə məntiqi qruplaşdırma.

### 3. Paket Adlandırması
**Qaydalar:** qısa, konseptual, KİÇİK hərflər, alt xətt YOX, camelCase YOX, tək formada.

| Pis | Yaxşı |
|---|---|
| stringconversion | strconv |
| synchronizationprimitives | sync |
| measuringtime | time |
| StringConversion | (stil pozuntusu) |

Qısaltmalar TÖVSİYƏ olunur (strconv, regexp, os) — icma üçün tanışdırsa.
**QAÇIN:** misc, util, common, data — məqsəd aydın DEYİL.

### 4. Package Declaration
Hər .go faylıın İLK sətri: `package <ad>`. Eyni qovluqdakı faylların declaration-u
EYNİ olmalıdır — fərqlisə compile xətası ("found packages shape and notright").

### 5. Exported vs Unexported
**Yeganə qayda:** BÖYÜK hərf = xaricə görünür; kiçik hərf = yalnız paket daxili.
Access modifier YOXDUR (Java-dakı private/public kimi).

```go
// strings paketindən:
func Contains(...)      // EXPORTED — main-dən çağırılır
func explode(...)       // unexported — strings.explode() XƏTA verir
```
**Praktika:** yalnız LAZIM OLANI export et — qalanını gizlət (interfeysin iç
mexanizmi kimi). Ch7-dəki Shape nümunəsində: PrintShapeDetails EXPORTED, area()/name()
metodları UNEXPORTED saxlanılır — xarici istifadəçi yalnız hazır API çağırır.

### 6. GOROOT və GOPATH
**Compiler paketləri haradan tapır:**
- **$GOROOT** — standart kitabxananın yeri (Go quraşdırması)
- **$GOPATH** — SƏNİN + 3-cü tərəf paketlər

**$GOPATH strukturu:**
```
$GOPATH/
├── bin/    ← go install binary-ləri
├── pkg/    ← object faylları (compile sürəti üçün)
└── src/    ← MƏNBƏ kodlar (bizim maraq dairəsi)
```

**Import yolu = $GOPATH/src-dən nisbi yol:**
```
$GOPATH/src/person/address/       → import "person/address"
$GOPATH/src/github.com/X/Y/Z     → import "github.com/X/Y/Z"  (repo pattern)
```
Paket ADI = yolun SON qovluğu.

### 7. İlk Paket (kitabdan)
```go
// $GOPATH/msg/msg.go:
package msg
import "fmt"

func Greeting(str string) {           // BÖYÜK G = exported
    fmt.Printf("Greeting %s\n", str)
}

// $GOPATH/demoimport/demoimport.go:
package main
import (
    "fmt"
    "msg"                              // paket importu
)
func main() {
    msg.Greeting("George")            // paketAdı.Funksiya
}
```

### 8. Paket Aliasları
```go
import f "fmt"        // f alias
func main() {
    f.Println("Hello, Gophers")   // f ilə çağırış
}
```
**Niyə:** ad uzun/qrtuluqdursa qısaldma; İKİ eyni adlı paketi ayırmaq.

### 9. main Paketi — Executable
2 paket növü:
- **main** = EXECUTABLE — main() funksiyası MÜTLƏQ (entry point); go build → binary
  (adı = qovluq adı)
- **Hamısı başqaları** = non-executable (kitabxana) — build binary YARATMIR

### 10. init() Funksiyası
**Nədir:** paket inicializasiyası üçün xüsusi funksiya; arqumentsiz, dəyərsiz;
AVTOMATİK çağırılır.

**İstifadə halları:** DB bağlantıları, paket dəyişənlərinin init-i, fayl yaratma,
konfiq yükləmə, state yoxlama/təmir.

**İCRA SIRASI (məcburi):**
1. İmport olunan paketlər initialize olunur (ƏVVƏL)
2. Paket səviyyəli dəyişənlər
3. init() funksiyası
4. main() icra olunur

```go
var name = "Gopher"      // 2-ci: dəyişən
func init() {            // 3-cü: init
    fmt.Println("Hello, ", name)
}
func main() {            // 4-cü: main
    fmt.Println("Hello, main function")
}
// Çıxış: Hello, Gopher → Hello, main function
```
`func init(age int)` — XƏTA: arqument QADAĞA.

### 11. Çoxsaylı init() Funksiyaları
Bir paketdə BİRDƏN ÇOX init() ola bilər — icra KODDAKI YAZILIŞ SIRASI ilə:
```go
func init() { fmt.Println("Hello, ", name) }   // 1-ci
func init() { fmt.Println("Second") }           // 2-ci
func init() { fmt.Println("Third") }            // 3-cü
func main() { fmt.Println("Hello, main function") }  // ən sonda
```
**Məqsəd:** inicializasiyanı modullaşdır (fayl+DB+state ayrı init-lərdə) — debug və
maintenance asanlaşır. **Diqqət:** sıra vacibdir — asılı init-lərdə gözlənilməz
nəticələr.

### 12. Budget Nümunələri (kitabdan)
```go
var budgetCategories = make(map[int]string)
var payeeToCategory = make(map[string]int)

func init() {                              // 1-ci init — kateqoriyaları yüklə
    budgetCategories[1] = "Car Insurance"
    budgetCategories[2] = "Mortgage"
    ...
}
func init() {                              // 2-ci init — payee xəritəsi
    payeeToCategory["Integrity Insurance"] = 1
    ...
}
func main() {
    for k, v := range payeeToCategory {   // cross-map sorğu
        fmt.Printf("Payee: %s, Category: %s\n", k, budgetCategories[v])
    }
}
```

### 13. Payroll Activity — Modul Paketləşdirmə
Ch7-dəki Payer activity-si paketə çevrilir: Developer/Employee/Manager tipləri
`payroll` paketinə köçürülür; tiplər/metodlar düzgün export/unexport edilir; tiplər
MƏNTİQİ FAYLLARA bölünür; main payroll-un alias istifadəçisi olur; 2 init() —
salamlama + dəyişən inicializasiyası.

## Əsas terminlələr
- Package — eyni qovluqdakı .go fayllar qrupu; kod vahidi
- DRY — Don't Repeat Yourself; funksiya→fayl→paket zənciri
- Package Declaration — hər faylın ilk sətri
- Exported/Unexported — BÖYÜK/kiçik hərf görünürlüyü
- GOROOT — standart kitabxana kökü
- GOPATH — istifadəçi+3-cü tərəf kodların src kökü
- pkg/bin/src — GOPATH-in 3 qovluğu
- Import Path — $GOPATH/src-dən nisbi yol
- Package Alias — `import f "fmt"` alternativ ad
- Executable/Non-Executable — main paketi / kitabxana paketi
- init() — avtomatik inicializasiya funksiyası; arqumentsiz
- Init Order — import → variables → init → main
- Package Files — paketi təşkil edən .go faylları

## Praktik nətidə

(1) Eyni kodu kopyalamaq əvəzinə paket yarat — bug düzəlişi BİR yerdə. (2) Paket adı:
qısa, kiçik hərf, tanış qısaltma; misc/util/common QAÇIN. (3) Fayl adları funksiyaya
görə (strings nümunəsi) — paket→fayl→funksiya iyerarxiyası. (4) Export minimumda
saxla — daxili metodlar kiçik hərfdə qalsın; istifadəçiyə hazır API təqdim et.
(5) Eyni qovluq = eyni package declaration — fərqlisə compile dayanır. (6) Import yolu
son qovluq adı ilə biter — paket adı oradan gəlir. (7) Alias: uzun adlar və ad
toqquşmaları üçün. (8) main + main() = executable; digərləri kitabxana. (9) init()
statik data/konfiq yükləmək üçün idealdır; sıra: import → var → init → main. (10)
Çoxsaylı init kod SIRASI ilə icra olunur — asılı inicializasiyaları düz yerləşdir.
(11) Modul paketi + main alias — böyük layihələrin əsas arxitekturası.

## Mənbə
Pages: 291-320 (PDF 324-355)
