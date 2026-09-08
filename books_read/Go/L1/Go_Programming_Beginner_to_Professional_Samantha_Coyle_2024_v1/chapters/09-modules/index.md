# Chapter 9 — Using Go Modules to Define a Project (səh. 326-341)

## Bu fəsil nədən bəhs edir?

Go modulları: go.mod/go.sum faylları, modulun yaradılması (`go mod init`),
xarici asılılıqlar (`go get`), çoxmodullu layihə strukturu və Go 1.18
workspaces (`go.work` — replace əvəzinə).

## Əsas fikirlər

### 1. Modul nədir?
**Nədir:** versiyalanmış paket dəsti — layihənin asılılıq və versiya
idarəetməsinin təməl vahidi (Go 1.11-dən). Modul = self-contained
qutu; kod paylaşımı/kollaborasiya asanlaşır.

**Komponentlər:**

**go.mod — blueprint:**
```
module mymodule
require (
    github.com/some/dependency v1.2.3
    github.com/another/dependency v2.0.0
)
replace (
    github.com/dependency/v3 => github.com/dependency/v4
)
exclude (
    github.com/some/dependency v2.0.0
)
```
- Module path + asılılıqlar + versiyalar
- replace — dəyişdirilmiş asılılıq (lokal test/uyğunluq üçün)
- exclude — problemli versiyanı kənarlaşdır

**go.sum — təhlükəsizlik:**
```
github.com/some/dependency v1.2.3 h1:abcdefg...
github.com/some/dependency v1.2.3/go.mod h1:hijklm...
```
- SHA-256 checksumları — paketlərin zədələnmədiyini/təhrif olunmadığını
  yoxlayır; avtomatik generasiya olunur

### 2. Modulun yaradılması
```bash
mkdir bookutil
cd bookutil
go mod init bookutil        # → go.mod yaranır (go.sum hələ YOX)
```

**Paket əlavəsi (author/author.go):**
```go
package author

import "fmt"

// Author represents an author of a book.
type Author struct {
    Name    string
    Contact string
}

func NewAuthor(name, contact string) *Author {
    return &Author{Name: name, Contact: contact}
}

func (a *Author) WriteChapter(chapterTitle string, content string) {
    fmt.Printf("Author %s is writing a chapter titled '%s'\n", a.Name, chapterTitle)
    fmt.Println(content)
}

func (a *Author) ReviewChapter(chapterTitle string, content string) {
    fmt.Printf("Author %s is reviewing a chapter titled '%s'\n", a.Name, chapterTitle)
    fmt.Println(content)
}

func (a *Author) FinalizeChapter(chapterTitle string) {
    fmt.Printf("Author %s has finalized the chapter titled '%s'.\n", a.Name, chapterTitle)
}
```

**İstifadə (main.go):**
```go
package main

import "bookutil/author"      // modul-yolu/paket

func main() {
    authorInstance := author.NewAuthor("Jane Doe", "jane@example.com")
    chapterTitle := "Introduction to Go Modules"
    authorInstance.WriteChapter(chapterTitle, "Go modules provide...")
    authorInstance.ReviewChapter(chapterTitle, "This chapter looks great...")
    authorInstance.FinalizeChapter(chapterTitle)
}
```
- Modul adı ≠ paket adı (bir modulda çox paket ola bilər)
- Adlandırma: layihənin ƏSAS məqsədinə görə; github.com/<project>/
  praktikası

### 3. Xarici modullar
**Nə üçün:** kod təkrarından qaçınma, funksionallıq genişlənməsi,
sınaqdan keçmiş etibarlı kod (Google-un uuid paketi kimi).

```go
// main.go:
import (
    "fmt"
    "github.com/google/uuid"
)

func main() {
    id := uuid.New()
    fmt.Printf("Generated UUID: %s\n", id)
}
```
```bash
go get github.com/google/uuid    # endir + go.mod-a yaz
# require github.com/google/uuid v1.3.1
# go.sum-a checksum-lar əlavə olundu
go run main.go
# Generated UUID: 7a533339-58b6-4396-b7f7-d0a50216bf88
```
- Çoxlu xarici modul: uuid + rsc.io/quote (random sitat) birlikdə —
  hər biri öz require sətri

### 4. Çoxmodullu layihə
```
myproject/
├── mainmodule/      (main.go + go.mod + go.sum)
├── secondmodule/    (own go.mod)
└── thirdmodule/     (own go.mod)
```

**Submodul yaratmağın məqsədli halları:**
1. Komponentlərin FƏRQLİ asılılıqları var
2. Eyni asılılığın FƏRQLİ versiyaları lazımdır
3. Komponent başqa layihələrdə yenidən istifadə olunacaq
4. Ayrıca idarəetmə/test sadəliyi

Adi halda — adətkan DEYİL; aydın ehtiyac olmadan etməyin.

### 5. Go workspaces (Go 1.18+)
**Problem:** iki lokal modul — othermodule printer modulunu istifadə
etmək istəyir, amma printer GitHub-da YOXDUR (yalnız lokal).

**Köhnə yol — replace direktivi:**
```bash
go mod edit -replace github.com/sicoyle/printer=../printer
# othermodule/go.mod:
#   replace github.com/sicoyle/printer => ../printer
go mod tidy
go run main.go          # indi işləyir
```

**Yeni yol — go.work:**
```bash
# printer qovluğunda:
go work init            # go.work yaradır
go work use ./printer   # lokal modulu işə sal

go run othermodule/main.go   # replace YOXDAN işləyir!
```
- go.work bir neçə go.mod-u koordinasiya edir — hər modulu əl ilə
  redaktə etmək lazım deyil
- İri layihələrdə/multi-repo işdə əvəzsiz

## Əsas terminlər
- Module — versiyalanmış paket kolleksiyası
- go.mod — modul blueprintfiliz (path/require/replace/exclude)
- go.sum — SHA-256 checksum anbarı (təhrif qoruması)
- Semantic versioning — v1.2.3 (major.minor.patch)
- `go mod init` — modulun yaradılması
- `go get` — asılılığın endirilməsi
- `go mod tidy` — lazımsızları sil, çatışmayanları əlavə et
- replace direktivi — asılılığın dəyişdirilməsi
- Submodule — ayrıca go.mod-lu komponent
- UUID — universally unique identifier
- Go workspace / go.work — çoxmodullu lokal koordinasiya
- `go work init` / `go work use` — workspace idarəetməsi

## Praktik nəticə
Hər yeni layihə: `go mod init <ad>` → paketlər qovluqlarda → main.go-da
`<modul>/<paket>` importu ilə istifadə. Xarici paket: `go get` —
go.mod/go.sum avtomatik dolur; əl ilə asılılıq izləmə yoxdur. Lokal
ikili-modul işində köhnə `replace` hələ işləyir, amma `go work init`
+ `go work use` daha təmizdir. Submodullar yalnız real ehtiyacda:
fərqli asılılıq dəstləri, versiya münaqişəsi, yenidən-istifadə.

## Mənbə
Pages: 326-341 (PDF 326-341)
