# Chapter 3 — Packaging and tooling

## Bu chapter nədən bəhs edir?

Go kodunun paketlərə təşkili, import mexanizmi, `init` funksiyaları, go alətləri (`build`, `run`, `vet`, `fmt`, `doc`), kod paylaşması qaydaları və asılılıq idarəsi (vendoring, godep, gb). Bu chapter Go ekosisteminin infrastrukturunu izah edir.

## Əsas fikirlər

### 1. Paketlər (Packages)
**Nədir:** Bütün Go proqramları paket adlanan fayl qruplarına təşkil olunur — kiçik, təkrar istifadə olunan kod vahidləri.

**Necə işləyir:** Paket bir qovluqda yerləşir; bir qovluqda yalnız bir paket ola bilər, paket bir neçə qovluğa bölünə bilməz. Hər `.go` faylının ilk sətri (şərh/boşluqdan başqa) paket elanı olmalıdır.

**Nümunə — `net/http` paketinin daxili strukturu:**
```
net/http/
    cgi/
    cookiejar/
    fcgi/
    httptest/
    httputil/
    pprof/
```
Hər alt qovluq ayrıca import oluna bilər — HTTP client qurursansa yalnız `http` paketi bəsdir.

**Adlandırma konvensiyası:** qovluq adı = paket adı. Qısa, kiçik hərflə,冗edaktə zamanı çox yazıldığı üçün concise (`cgi`, `httputil`, `pprof` kimi). Unikal ad məcburi deyil — import tam path ilə olur.

### 2. `main` paketi — command (executable)
**Nədir:** `main` paketi xüsusi mənanı daşıyır: bu paket binary executable-a kompilyasiya olunacaq.

**Necə işləyir:** Kompilyator `main` paketinə rast gəlincə `main()` funksiyasını da axtarır — yoxsa executable yaranmır. Binary-nin adı `main` paketinin qovluq adını alır.

**Terminologiya:** Go sənədlərində **command** = executable proqram (məs. CLI aləti); **package** = import olunan semantik vahid. Yeni developer-lər üçün qarışdırılır.

**Kitabdan kod nümunəsi (yalnış vs düzgün):**
```go
// DÜZGÜN — executable yaranır:
package main

import "fmt"

func main(){
    fmt.Println("Hello World!")
}
```
```go
// YALNıŞ — main funksiyası var, amma paket adı main deyil:
package hello

import "fmt"

func main(){
    fmt.Println("Hello, World!")
}
```
İkinci fayldan binary yaranmır — sadəcə importable paket kimi qəbul edilir.

### 3. Import mexanizmi
**Nədir:** `import` kompilyatora paketin diskdə harada axtarılacağını bildirir.

**Necə işləyir:** Axtarış sırası:
1. Go quraşdırma qovluğu (`GOROOT/src/pkg/...`)
2. GOPATH-dəki hər qovluq sırayla (`$GOPATH1/src/...`, `$GOPATH2/src/...`)

Nümunə: Go `/usr/local/go`-da, `GOPATH=/home/myproject:/home/mylibraries` isə, `net/http` üçün axtarış:
```
/usr/local/go/src/pkg/net/http
/home/myproject/src/net/http
/home/mylibraries/src/net/http
```
İlk tapılan yer dayandırır; heç biri tapılmasa build xətası.

**İdiomatik multi-import:**
```go
import (
    "fmt"
    "strings"
)
```

### 4. Remote imports və `go get`
**Nədir:** Import path URL olduqda Go tooling paketi şəbəkədən çəkə bilər.

**Kitabdan kod nümunəsi:**
```go
import "github.com/spf13/viper"
```

**Necə işləyir:** `go build` bu path-i GOPATH-də axtarır; `go get` isə GitHub/Launchpad/Bitbucket-dən paketi çəkib GOPATH-də URL-ə uyğun yerə qoyur. `go get` **rekursivdir** — paketin öz asılılıqlarını da gəzib çəkir.

### 5. Named imports (adlandırılmış import)
**Nədir:** Eyni adda iki paket import edilməli olanda biri yenidən adlandırılır.

**Kitabdan kod nümunəsi:**
```go
package main

import (
    "fmt"
    myfmt "mylib/fmt"
)

func main() {
    fmt.Println("Standard Library")
    myfmt.Println("mylib/fmt")
}
```

**Sub-kod izahı:**
- `myfmt "mylib/fmt"` → `myfmt` alias (təxəllüs) verir; artıq `myfmt.Println` ilə çağırılır
- İstifadə olunmayan import **build xətasıdır** — Go komandası bunu kod şişməsinin qarşısını almaq üçün qəsdən edir
- İstifadə etmədən import etmək lazımdırsa: blank identifier `_`

### 6. `init` funksiyaları
**Nədir:** Hər paket istənilən sayda `init` funksiyası təyin edə bilər — bunlar `main`-dən əvvəl avtomatik icra olunur.

**Nəyə lazımdır:** Paketin setup-u, dəyişənlərin ilkinləşdirilməsi, bootstrapping — xüsusən **qeydiyyat (registration)** pattern-i.

**Kitabdan kod nümunəsi (database driver pattern):**
```go
// Driver paketi (github.com/goinaction/code/chapter3/dbdriver/postgres):
package postgres

import (
    "database/sql"
)

func init() {
    sql.Register("postgres", new(PostgresDriver))
}
```
```go
// İstifadəçi tərəfi:
package main

import (
    "database/sql"

    _ "github.com/goinaction/code/chapter3/dbdriver/postgres"
)

func main() {
    sql.Open("postgres", "mydb")
}
```

**Sub-kod izahı:**
- Driver paketin `init`-i özünü `sql` paketinə qeydiyyata alır — `sql` paketi kompilyasiya vaxtı hansı driver-lərin mövcud olduğunu bilmir, runtime qeydiyyatından öyrənir
- `_ "…/postgres"` → blank identifier ilə import: identifikator istifadə olunmur, amma `init` çağırılır
- `sql.Open("postgres", "mydb")` → "postgres" adı ilə qeydiyyata düşmüş driver işə düşür

### 7. Go alətləri (`go` command)
**Kitabdan kod nümunəsi (wordcount — listing 3.7):**
```go
package main

import (
    "fmt"
    "io/ioutil"
    "os"

    "github.com/goinaction/code/chapter3/words"
)

// main is the entry point for the application.
func main() {
    filename := os.Args[1]

    contents, err := ioutil.ReadFile(filename)
    if err != nil {
        fmt.Println(err)
        return
    }

    text := string(contents)

    count := words.CountWords(text)
    fmt.Printf("There are %d words in your text. \n", count)
}
```

**Sub-kod izahı:**
- `os.Args[1]` → kompensasiya sətri arqumenti (fayl adı)
- `ioutil.ReadFile(filename)` → bütün faylı oxuyur, `[]byte` + error qaytarır
- `string(contents)` → byte slice-ı string-ə çevirir

### `go` əmrlərinin izahı (cheat sheet üçün əsas):

**`go build <pkg|file|.>` →** paketi/faylı kompilyasiya edib binary yaradır. Fayl adı verməsən cari paket götürülür; `.` = cari qovluq.

**`go build github.com/goinaction/code/chapter3/wordcount` →** müəyyən paketi birbaşa build edir.

**`go build github.com/goinaction/code/chapter3/...` →** `...` wildcard — chapter3 altındakı **bütün** paketləri build edir.

**`go clean <file>` →** build nəticəsində yaranan executable-ı silir (source control-ə təmiz checkin üçün).

**`go run wordcount.go` →** build + execute bir addımda — binary diskdə qalmır, ən çox istifadə olunan əmr.

**`go vet <file>` →** kodu ümumi xətalara göre yoxlayır:
- Printf tipli çağırışlarda yanlış parametrlər
- Ümumi metod təriflərində imza xətaları
- Yanlış struct tag-lər
- Açarsız composite literal-lər

Nümunə: `fmt.Printf("…dogs", 3.14)` format string-də placeholder yoxdursa:
```
go vet main.go
main.go:6: no formatting directive in Printf call
```
**Tövsiyə:** commit-dən əvvəl `go vet` vərdiş et.

**`go fmt <file|pkg>` →** kodu avtomatik formatlayır. Mübahisəli məqzaları (tab/space, `{` yeri) aradan qaldırır:
```
if err != nil { return err }
```
`go fmt`-dan sonra:
```go
if err != nil {
    return err
}
```
**Tövsiyə:** editor-da save-də avtomatik `go fmt` quraşdır.

**`go doc tar` →** paket sənədlərini birbaşa terminalda göstərir — workflow pozulmur.

**`godoc -http=:6060` →** localhost:6060-da sənədləşmə web server-i işə salır (Go saytının özü də bunun variantıdır). Həm standart kitabxana, həm sənin GOPATH-dəki kodlar görünür.

### 8. Kod sənədləşdirməsi (godoc konvensiyası)
**Nədir:** Şərhlər müəyyən konvensiya ilə yazılsa godoc onları avtomatik sənədləşdirməyə çevirir.

**Kitabdan kod nümunəsi:**
```go
// Retrieve connects to the configuration repository and gathers
// various connection settings, usernames, passwords. It returns a
// config struct on success, or an error.
func Retrieve() (config, error) {
    // ... omitted
}
```

**Sub-kod izahı:**
- Şərh **identifikatorun birbaşa üstündə** durmalıdır — paket, funksiya, tip, global dəyişən üçün işləyir
- Tam cümlələrlə yazılır

**Paket səviyyəli sənəd — `doc.go` faylı:**
```go
/*
Package usb provides types and functions for working with USB
devices.
To connect to a USB device start by creating a new USB
connection with NewConnection
...
*/
package usb
```
- `doc.go` → paketlə eyni adda, paket elanından əvvəl böyük şərh bloku; godoc-də paketin başlanğıc mətni kimi göstərilir

### 9. Kod paylaşma qaydaları
**Paket repozitoriyanın kökündə olmalıdır:** `go get` tam path istifadə edir — public repoda `code/` və ya `src/` qovluğu yaratmaq xətadır; source faylları repozitoriyanın root-unda olmalıdır (import path qısa qalsın).

**Paketlər kiçik ola bilər:** Bir tapşırıq görən kiçik paket Go-da normaldır.

**`go fmt` işlət + sənədləşdir:** İnsanlar koduna baxıb keyfiyyətini ölçür.

### 10. Dependency management (asılılıq idarəsi)
**Problem:** `import` statement-i hansı **reviziyanın** çəkilməli olduğunu bildirmir — `go get` hər dəfə fərqli versiya gətirə bilər → reproducible build (təkrarlana bilən build) problemi.

**Həll 1 — Vendoring + import path rewriting (godep, vendor tools):**
Bütün asılılıqlar layihə repozitoriyasının içindəki qovluğa kopyalanır və import-lar yenidən yazılır.

godep strukturu:
```
$GOPATH/src/github.com/ardanstudios/myproject
|-- Godeps
|   |-- Godeps.json
|   |-- _workspace
|       |-- src
|           |-- bitbucket.org/ww/goautoneg/...
|           |-- github.com/beorn7/perks/...
|-- main.go
```

Import dəyişikliyi (əvvəl → sonra):
```go
// Əvvəl:
import (
    "bitbucket.org/ww/goautoneg"
    "github.com/beorn7/perks"
)
// Sonra (vendored):
import (
    "github.ardanstudios.com/myproject/Godeps/_workspace/src/bitbucket.org/ww/goautoneg"
    "github.ardanstudios.com/myproject/Godeps/_workspace/src/github.com/beorn7/perks"
)
```

**Üstünlük:** bütün source layihə daxilindədir → təkrarlana bilən build; repo go-gettable qalır.
**Çatışmamazlıq:** import path-lər çox uzun və əl ilə yazılmış kimi görünür.

**Həll 2 — gb (project-based build tool):**
GOPATH və Go tooling-i tamamlayır; import path rewriting lazım deyil.

gb layihə strukturu:
```
$PROJECT/
|-- src/           ← sənin kodun
|   |-- cmd/myproject/main.go
|-- vendor/
|   |-- src/       ← vendored asılılıqlar
|       |-- bitbucket.org/ww/goautoneg/
|       |-- github.com/beorn7/perks/
```
- `$PROJECT` — heç bir environment variable tələb olunmur
- Import-lar dəyişmir: gb əvvəl `$PROJECT/src/`-də, sonra `$PROJECT/vendor/src/`-də axtarır
- Build: `gb build all` ($PROJECT qovluğunda)
- **Çatışmamazlıq:** gb layihəsi Go tooling (o cümlədən `go get`) ilə uyğun deyil
- Sayt: getgb.io

## Əsas terminlər
- Package (paket)
- Command (executable proqram)
- GOPATH / GOROOT (workspace / quraşdırma yolu)
- Remote Import (ucqordan import)
- Named Import (adlandırılmış import)
- Blank Identifier (boş identifikator)
- init Function (başlanğıc funksiyası)
- Vendoring (asılılıqların layihə daxilinə kopyalanması)
- Import Path Rewriting (import yolunun yenidən yazılması)
- Reproducible Build (təkrarlana bilən build)
- Struct Tag (struktur teqi)

## Praktik nəticə
- Paket qovluğu = paket adı; bir qovluq bir paket — bu qaydanı pozma.
- Driver/large framework qurarkən init + Register pattern istifadə et (kitabxana istifadəçisi yalnız `_ "driver/path"` yazır).
- `go vet` + `go fmt`-i commit pipeline-ına daxil et.
- Public repoda source root-da saxla, `src/` qovluğu yaratma.
- Müasir Go (v1.11+) üçün göstərilən godep/gb tarixi yanaşmalardır — amma vendoring fəlsəfəsi `go mod vendor`-da yaşayır. (Müəllim qeydi: bugünkü layihələrdə Go modules istifadə et — `go.mod` + `go.sum` bu problemi standart həll edir.)

## Mənbə
Pages: 60-77 (PDF), book pages 39-56
