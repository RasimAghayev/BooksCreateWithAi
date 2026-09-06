# Chapter 1 — Setting Up Your Go Environment (Go Mühitinin Qurulması)

## Bu chapter nədən bəhs edir?

Go proqramlaşdırma mühitinin qurulması: Go alətlərinin installasiyası, workspace
anlayışı, `go run` və `go build` komandaları, üçüncü tərəf alətlərin quraşdırılması,
kod formatlaması (`go fmt`), lint və vet alətləri, IDE seçimləri və Makefile ilə
avtomatik build prosesinin təşkili.

## Əsas fikirlər

### 1. Go Alətlərinin Installasiyası
**Nədir:** Go dilini kompilyasiya etmək üçün lazım olan rəsmi toolchain-in kompüterə qurulması.

**Necə işləyir:** Go saytından platformaya uyğun installer endirilir. Mac üçün `.pkg`,
Windows üçün `.msi` installer Go-nu avtomatik doğru yerləşdirir, köhnə versiyanı silir
və `go` binary-ni default PATH-ə əlavə edir. Linux/FreeBSD üçün gzipped tar faylı
`/usr/local`-ə açılır və `/usr/local/go/bin` PATH-ə əlavə olunur.

**Nəyə lazımdır:** Go kodu yazmağa başlamazdan əvvəl mühitin düzgün qurulması üçün.

**Kitabdan kod nümunəsi:**
```bash
# Linux installasiyası
$ tar -C /usr/local -xzf go1.15.2.linux-amd64.tar.gz
$ echo 'export PATH=$PATH:/usr/local/go/bin' >> $HOME/.profile
$ source $HOME/.profile
```

**Sub-kod izahı:**
- `tar -C /usr/local -xzf` → arxivi birbaşa `/usr/local` qovluğuna açır
- `export PATH=$PATH:...` → `go` komandasını istənilən terminaldan işlədilən edir
- `source $HOME/.profile` → dəyişikliyi cari terminal sessiyasına tətbiq edir

**Quraşdırmanın yoxlanması:**
```bash
$ go version
# go version go1.15.2 darwin/amd64
```

**İzah:** Versiya məlumatı çıxırsa mühit hazırdır. Xəta alınsa, `which go` ilə PATH-də
hansı `go`-nun olduğu yoxlanılır (başqa `go` adlı proqram mane ola bilər).

### 2. Go Workspace və GOPATH
**Nədir:** Üçüncü tərəf Go alətlərinin (`go install` ilə qurulan) saxlandığı mərkəzi qovluq.

**Necə işləyir:** Default olaraq `$HOME/go`-dur: source `$HOME/go/src`-də, compile
olunmuş binary-lər `$HOME/go/bin`-də. `$GOPATH` dəyişəni ilə fərqli yer təyin etmək olar.

**Nəyə lazımdır:** `go install` ilə qurulan alətlərin avtomatik tapılması və PATH üzərindən
işlədilməsi üçün. Layihələrin təşkilində artıq məcburi deyil (müasir Go modulları ilə
layihəni istənilən yerdə saxlamaq olar).

**Kitabdan kod nümunəsi:**
```bash
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin
```

**Sub-kod izahı:**
- `GOPATH=$HOME/go` → workspace-i açıq şəkildə təyin edir
- `PATH=$PATH:$GOPATH/bin` → qurulan alətləri komanda sətirindən işlədilən edir

**Vacib qeyd:** `GOROOT` dəyişəni artıq təyin etmək lazım deyil — `go` tool özü bunu bilir.
Bütün go mühit dəyişənlərini `go env` komandası ilə görmək olar.

### 3. `go run` və `go build`
**Nədir:** Go kodunu işə salmanın iki əsas yolu.

**Necə işləyir:**
- `go run hello.go` → kodu müvəqqəti qovluqda binary-yə compile edir, işə salır və
  binary-ni **silir**. Diskdə iz qalmır.
- `go build hello.go` → cari qovluqda `hello` (Windows-da `hello.exe`) binary yaradır.

**Nəyə lazımdır:**
- `go run` → kiçik test proqramları, skript kimi istifadə
- `go build` → paylamaq üçün binary yaratmaq

**Kitabdan kod nümunəsi:**
```go
package main
import "fmt"
func main() {
    fmt.Println("Hello, world!")
}
```

```bash
go run hello.go        # dərhal işə salır, binary saxlamır
go build hello.go      # hello binary yaradır
go build -o hello_world hello.go   # fərqli ad ilə binary
```

**Sub-kod izahı:**
- `-o hello_world` → çıxış binary-nin adını/yolunu dəyişir

### 4. Üçüncü tərəf alətlərin qurulması — `go install`
**Nədir:** Go alətlərinin mənbəkoddan compile edilib `$GOPATH/bin`-ə qurulması.

**Necə işləyir:** Go mərkəzi registry istifadə etmir (Maven/NPM-dən fərqli olaraq) —
alətlər birbaşa source repository-dən gəlir. `go install repo@latest` repo-nu endirir,
compile edir və binary-ni `$GOPATH/bin`-ə qoyur.

**Nəyə lazımdır:** Məsələn `hey` HTTP load-test alətini qurmaq üçün.

**Kitabdan kod nümunəsi:**
```bash
$ go install github.com/rakyll/hey@latest
$ hey https://www.golang.org
```

**Sub-kod izahı:**
- `@latest` → ən yeni versiyanı götürür; konkret versiya də göstərilə bilər
- Yeniləmək üçün eyni komanda yenidən işə salınır

**Üstünlükləri:** Alət binary-si adi executable-dır, istənilən yerdə saxlanıla bilər;
`go install` digər Go developerlərinə paylamanın ən rahat yoludur.

### 5. Kod formatlaması — `go fmt` və `goimports`
**Nədir:** Go-nun standart kod formatını avtomatik tətbiq edən alətlər.

**Necə işləyir:** Go format müzakirələrini aradan qaldırır — tək standart format var.
`go fmt` whitespace, indentasiya (tab), struct sahələrinin düzülüşünü və operator
aralıqlarını avtomatik düzəldir. `goimports` bundan əlavə import-ları sıralayır,
istifadə olunmayanları silir, çatışanları təxmin edərək əlavə edir.

**Nəyə lazımdır:** Go developer icmasının gözlədiyi vahid kod görünüşü üçün.

**Kitabdan kod nümunəsi:**
```bash
go install golang.org/x/tools/cmd/goimports@latest
goimports -l -w .
```

**Sub-kod izahı:**
- `-l` → yanlış formatlı faylların adlarını çap edir
- `-w` → faylları yerində dəyişir
- `.` → cari qovluq və bütün alt qovluqlar

**Semicolon insertion qaydası (önəmli detal):** Go kompilyatoru növbəti sətirə keçərkən
son token identifier, literal, `break`, `continue`, `fallthrough`, `return`, `++`, `--`,
`)` və ya `}` olduqda avtomatik semicolon əlavə edir. Ona görə `{`-i yeni sətirə yazmaq
(`func main()` + yeni sətirdə `{`) compile xətası yaradır — `func main();` forması
yaranır. Bu qayda həm dilin sadəliyini, həm də vahid brace stilini təmin edir.

### 6. Lint və Vet
**Nədir:** Kod keyfiyyətini yoxlayan alətlər.

**Necə işləyir:**
- `golint` → stil qaydalarını yoxlayır (adlandırma, error mesaj formatı, public
  metodlarda şərh). False positive/negative ola bilər — məcburi deyil.
- `go vet` → sintaktik düzgün, amma məntiqən şübhəli kodu tapır (format funksiyalarına
  yanlış sayda arqument, istifadə olunmayan dəyişən təyinatı və s.). Kompilyasiyanı
  dayandırmayan, amma real bug-ları tutan yoxlamalar.
- `golangci-lint` → golint + go vet + 50-ə yaxın başqa linter-i bir komandada birləşdirir.

**Nəyə lazımdır:** Common bug-ların vaxtından əvvəl tutulması, idiomatik kod.

**Kitabdan kod nümunəsi:**
```bash
go install golang.org/x/lint/golint@latest
golint ./...
go vet ./...
golangci-lint run
```

**Sub-kod izahı:**
- `./...` → cari qovluq və bütün alt qovluqlardakı paketlər
- `.golangci.yml` → hansı linter-lərin aktiv olduğu konfiqurasiya olunur

**Tövsiyə (kitab):** `go vet` build prosesinə məcburi daxil et; `golint` code review
prosesində istifadə et; team razılaşması olmadan golangci-lint qaydalarını məcburi etmə.

### 7. IDE seçimləri: VS Code, GoLand, Go Playground
**Nədir:** Go inkişafı üçün editor/IDE seçimləri.

**Necə işləyir:**
- **VS Code** (pulsuz): Go extension qurulur; Delve debugger və `gopls` (Go Language
  Server) avtomatik qurulur. Language server — editərə code completion, lint, istifadə
  axtarışı kimi ağıllı funksiyalar verən standart API spesifikasiyasıdır.
- **GoLand** (JetBrains, pullu): əlavə alət qurulması tələb etmir; refactoring,
  debugger, code coverage daxildir.
- **Go Playground** (onlayn): quraşdırma tələb etmir; kiçik proqramları yazıb paylaşmaq
  üçün. Şərtlər: network bağlantısı yoxdur, uzun/məşəqqətli proseslər dayandırılır,
  saat 10 noyabr 2009 (Go-nun elan tarixi) göstərilir. **Share URL publicdir —
  sensitive məlumat yazmaq QADAĞANDIR.**

### 8. Makefile ilə build avtomatlaşdırması
**Nədir:** Build addımlarını (fmt → vet → build) təkrarlana bilən şəkildə təsvir edən skript.

**Necə işləyir:** Hər əməliyyat target adlanır. `.DEFAULT_GOAL` komandasız `make`
çağırışında işə düşən target-i təyin edir. Target-dən asılılıqlar (`build: vet`) əvvəl
icarə olunur.

**Kitabdan kod nümunəsi:**
```makefile
.DEFAULT_GOAL := build
fmt:
        go fmt ./...
.PHONY:fmt
lint: fmt
        golint ./...
.PHONY:lint
vet: fmt
        go vet ./...
.PHONY:vet
build: vet
        go build hello.go
.PHONY:build
```

**Sub-kod izahı:**
- `.DEFAULT_GOAL := build` → sadə `make` = `make build`
- `build: vet` → build-dən əvvəl vet icra olunur; vet də fmt-ə bağlıdır →
  `make` = fmt + vet + build zənciri
- `.PHONY:fmt` → target adı ilə eyni adlı fayl/qovluq yaranarsa make-in çaşmaması üçün
- Addımlar **mütləq TAB** ilə indent olunmalıdır

**Çatışmamazlıqları:** Makefile sintaksisi dözümlü deyil (tab tələbi); Windows-da
əvvəlcədən yoxdur (`choco install make` ilə qurulur).

### 9. Go versiyalarının paralel saxlanması və yenilənməsi
**Nədir:** Fərqli Go versiyalarının eyni anda saxlanması.

**Necə işləyir:** Go proqramları native binary-dır, ayrıca runtime yoxdur — versiya
yenilənməsi deploy olunmuş köhnə proqramlara təsir etmir. Yeni versiyanı test etmək
üçün əlavə mühit qurulur.

**Kitabdan kod nümunəsi:**
```bash
$ go install golang.org/dl/go.1.15.6@latest
$ go1.15.6 download
$ go1.15.6 build
```

**Sub-kod izahı:**
- `go1.15.6 download` → həmin versiyanın toolchain-ini endirir
- `go1.15.6 build` → layihəni məhz o versiya ilə compile edir

**Go Compatibility Promise:** Go 1.x boyu dil və standart kitabxanada backward-breaking
dəyişiklik YOXDUR (yalnız bug/təhlükəsizlik fix-ləri istisna). `go` komandalarının
flag-lərində dəyişiklik ola bilər. Buraxılışlar təxminən hər 6 ayda bir olur.

## Əsas terminlər
- Workspace (iş sahəsi) — üçüncü tərəf alətlərin saxlandığı `$GOPATH` qovluğu
- Linter (kod təmizləyici) — stil qaydalarını yoxlayan alət
- Language Server (dil serveri) — editərlər üçün ağıllı kod xidmətləri API-si
- Semicolon insertion rule (nöqtəli vergül əlavə qaydası) — Go-nun avtomatik `;` qoyan leksik qaydası
- Compatibility Promise (uyğunluq vədi) — Go 1.x boyu kod-un sınmaması zəmanəti

## Praktik nəticə

Professional Go layihəsinin minimal build zənciri belə qurulur: hər build əvvəli
`go fmt` → `go vet` (məcburi) → `go build`; lint code review mərhələsində. Bunu
Makefile-da `.DEFAULT_GOAL` + asılılıq zənciri ilə avtomatlaşdırmaq "mənim maşınımda
işləyir" problemini aradan qaldırır. `goimports` save-time avtomatikləşdirmək ən
yaxşı praktikadır.

## Mənbə
Pages: 13-36
