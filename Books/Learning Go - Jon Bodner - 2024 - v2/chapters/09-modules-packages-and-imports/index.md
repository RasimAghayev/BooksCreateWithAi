# Chapter 9 — Modules, Packages, and Imports (Modullar, Paketlər və Importlar)

## Bu chapter nədən bəhs edir?

Repository/modul/paket iyerarxiyasi, go.mod, import/export (böyük-kiçik hərf qaydası),
paket adlandırma və təşkili (cmd/pkg pattern), internal paketi, init funksiyası, dövrəvi
asılılıqlar, type alias, üçüncü tərəf modullar, SemVer, minimum version selection, v2+
versiyalama, vendoring, proxy serverlər və go.sum təhlükəsizliyi.

## Əsas fikirlər

### 1. Repository → Module → Package
**Anlayışlar:** Repository (VCS-də kod) ⊃ modul (`go.mod`-lu kök; birlikdə versiyalanır)
⊃ paketlər (təşkilati vahid). Bir repo-da çox modul məsləhətli DEYİL (versiya izləməsi
çətinləşir). Modul yolu = repo yolu (`github.com/jonbodner/proteus`), case-sensitive,
böyük hərfsiz.

**go.mod:** `go mod init MODULE_PATH` ilə yaradılır. Struktur:
```
module github.com/learning-go-book/money
go 1.15
require (
    github.com/learning-go-book/formatter v0.0.0-...
    github.com/shopspring/decimal v1.2.0
)
```
Əlavə bölmələr: `replace` (asılılığın yerini dəyiş), `exclude` (versiyanı qadağan).

### 2. Import və Export — Böyük Hərf Qaydası
**Nədir:** Paketdən kənar görünmə = adın **böyük hərflə** başlaması. Heç bir keyword
yoxdur (`public`/`export` YOX). Kiçik hərf/underscore — yalnız paket daxili.

**Qaydalar:**
- Export hər şey API-yə çevrilir → sənədlə (godoc), backward-compatible saxla.
- İstifadə olunmayan import — compile xətası (binary-də yalnız istifadə olunan kod).
- Eyni qovluqdakı bütün fayllarda eyni package clause; paketin adı package clause-dan
  götürülür, **qovluq adından YOX** (`formatter/` qovluğunda `package print` → `print.Format`).
- Qayda: qovluq adı = paket adı; istisnalar: `main` (import olunmur), keçərsiz
  identifier-li qovluq adları, versiya qovluqları (/v2).
- Relative import İŞLƏTMƏ — mütləq yol refactor dostudur.
- Import fayl block-undadır — hər faylda təkrar import lazımdır.

### 3. Paket Adlandırma Sənəti
**Prinsip:** Ad identifikatorların prefiksi kimi oxunur.

- `util` paketi YOX — funksional ad ver (`extract`, `format`).
- Eyni funksiya adı fərqli paketlərdə OK: `extract.Names` vs `format.Names`.
- Paket adını identifikatorda təkrar etmə: `extract.ExtractNames` YOX; `sort.Sort` /
  `context.Context` istisnadır (ad = paket adı olduqda təbii).
- `util.ExtractNames` heç nə deməir; `extract.Names` hər istifadədə mənanı daşıyır.

### 4. Modul Təşkili (cmd/pkg Pattern)
**Kiçik modul:** hər şey bir paketdə. **Böyüyəndə:**
```
module-root/
├── cmd/        ← hər binary üçün bir qovluq (package main)
├── pkg/        ← əsas Go kodu (test/deploy faylları çox olanda)
```
pkg daxilində funksionallıq "dilimləri" üzrə (customer paketi, inventory paketi) —
asılılıqları azaldır, gələcəkdə mikro-servisə parçalama asanlaşır. (Kat Zien,
GopherCon 2018 "How Do You Structure Your Go Apps".)

### 5. Paket Adının Override Edilməsi
**Nə vaxt:** Ad toqquşması (crypto/rand + math/rand hər ikisi `rand`):

**Kitabdan kod nümunəsi:**
```go
import (
    crand "crypto/rand"   // alias
    "math/rand"           // normal — rand kimi istifadə
)
func seedRand() *rand.Rand {
    var b [8]byte
    _, err := crand.Read(b[:])   // kripto-qrafik seed
    ...
    return rand.New(rand.NewSource(int64(binary.LittleEndian.Uint64(b[:]))))
}
```
`.` import (namespace hoplanması) — QADAĞAYAXIN (mənşəyi oxunmaz edir). `_` import
(blank) — yalnız init-i işə salır (aşağıda). Paket adları kölgələnə bilər — zərurətdə alias.

### 6. godoc Sənədləşdirməsi
**Qaydalar:** Şərh sənədlənən elementin DƏRHAL üstündə, boş sətir yoxdur; `//` + elementin
adı ilə başlayır; boş şərh sətri = yeni paraqraf; indent = preformatted. Paket şərhi
package clause-dan əvvəl; uzun şərhlər `doc.go` faylında.

**Kitabdan kod nümunəsi:**
```go
// Package money provides various utilities to make it easy to manage money.
package money

// Money represents the combination of an amount of money and the currency.
type Money struct { ... }

// Convert converts the value of one currency to another.
//
// It has two parameters: ... (parametr və qaytarma izahı, error davranışı)
// Supported currencies are:
//        USD - US Dollar
//        ...
func Convert(from Money, to string) (Money, error) { ... }
```
`go doc PAKET` / `go doc PAKET.IDENTIFIER` ilə baxılır. Export olunan hər identifikator
şərhsiz qalmasın — golint/golangci-lint yoxlayır.

### 7. internal Paketi
**Nədir:** `internal` adlı paketin export-ları yalnız **birbaşa valideyn + sibling**
paketlər üçün görünür; modul kənarından və uzaq qohum paketlərdən — compile xətası
("use of internal package ... not allowed").

**Nəyə lazımdır:** API-yə salmadan modul daxili paylaşım.

### 8. init Funksiyası — Mümkünsə QAÇIN
**Nədir:** Parametrsiz/qaytarmasız `init` funksiyası paket ilk istinad ediləndə avtomatik
icra olunur; yalnız side-effect ilə işləyir. Bir paketdə (hətta bir faylda!) bir neçə init
olabilir — icra sırası sənədlidir, amma əzbərləmək əvəzinə çəkin.

**Qaydalar:**
- Bugünkü əsas istifadə: tək təyinatla qurulmayan package-level dəyişənlər — amma
  onlar da effektiv immutable olmalıdır.
- Package-də birdən çox init YAZMA. Fayl/şəbəkə I/O edirsə — sənədləşdir.
- **Blank import pattern:** `_ "github.com/lib/pq"` — identifikator istifadəsi olmadan
  init (driver qeydiyyatı). Obsolete sayılır (gizli qeydiyyat), amma standart kitabxana
  (database drivers, image formats) compatibility zəmanətinə görə davam edir. Öz kodunda
  registry-ləri açıq qeyd et.

### 9. Dövrəvi Asılılıqlar (Import Cycle)
**Qadağan:** A→B import edirsə, B→A (birbaşa və ya dolayı) OLMAZ — "import cycle not
allowed". Sürətli kompilyator + başa düşülən kod üçün.

**Həllər:** (1) paketləri birləşdir (hər ikisi bir-birinə asılıdırsa, ehtimal ki, bir
paketedilər); (2) dövr yaradan elementləri bir paketə/yeni paketə köçür.

### 10. Type Alias — API-nin Zarafatı Pozulmadan Yenidən Təşkili
**Nədir:** `type Bar = Foo` — mövcud tipin yeni adı (yeni tip YOX!).

**Kitabdan kod nümunəsi:**
```go
type Foo struct { x int; S string }
type Bar = Foo        // alias — eyni tip
func MakeBar() Bar {
    bar := Bar{a: 20, B: "Hello"}  // Foo-nun unexported sahəsi belə görünür (eyni paket)
    var f Foo = bar                // convertsuz təyin — eyni tipdir
    return bar
}
```
Yeni metod/sahə alias-a YOX, orijinala əlavə olunur. Digər paketdən alias mümkündür,
amma unexported hissələrə çıxış vermir. Alias-edilməyənlər: package-level dəyişənlər,
struct sahə adları (bir dəfə ad seçilibsə, dəyişməzdir).

### 11. Üçüncü Tərəf Kod və Versiyalar
**Import:** eyni sistem — repo yolu (`github.com/shopspring/decimal`). `go build` etdiyində
getməmiş asılılıqları avtomatik endirir, `go.mod`-a əlavə edir; `go.sum`-da modulun
özünün və go.mod faylının hash-i yazılır. **go.mod + go.sum həmişə VCS-də saxla.**

**SemVer:** `vMAJOR.MINOR.PATCH` — patch = bug fix; minor = yeni backward-compatible
xüsusiyyət; major = sinan dəyişiklik.

**Əməliyyatlar:**
```bash
go list -m -versions MOD          # mövcud versiyalar
go get MOD@v1.0.0                 # konkret versiya (downgrade də)
go get -u=patch MOD               # cari minor-un ən son patch-i
go get -u MOD                     # ən yeni (minor da)
go mod tidy                       # istifadə olunmayan versiyaları təmizlə
```

**Minimum Version Selection:** Eyni modul D fərqli versiyalarla tələb olunursa (A: v1.1.0,
B: v1.2.0, C: v1.2.3) — ən YENİ tələb olunan götürülür (v1.2.3), tək nüsxə. npm-in
çoxversiyalı yanaşmasından fərqli — paket state-i və ölçü problemləri olmur. Uyğunsuzluqda
(import compatibility rule pozulubsa) müəlliflər düzəltməlidir — "community həlli".

**`// indirect`** — modulun öz go.mod-u olmayan/köhnə asılılıqların sizin go.mod-da
görünməsi.

### 12. Major Versiya 2+ — Semver Import Rule
**Qayda:** major > 1 üçün modul yolu `/vN` ilə bitməlidir:
```go
import "github.com/learning-go-book/simpletax/v2"
```
İki fərqli major eyni proqramda paralel import oluna bilər (fərqli paketlər kimi) —
könüllü miqrasiya. Publisher tərəfdə: `/v2` qovluğu (kod kopyası) VƏYA VCS branch
(`v2` yeni kod / `v1` köhnə) + go.mod module yolunun və daxili importların `/vN`
yenilənməsi (Marwan Sulaiman avtomatikləşdirmə aləti) + `v2.0.0` teqi.

### 13. Vendoring
`go mod vendor` → `vendor/` qovluğu bütün asılılıqların nüsxəsi. go.mod dəyişəndə yenidən
işlədilməlidir (yoxsa build imtina edir). Fayda: dəqiq bəlli kod; mənfi: repo şişir.
Proxy serverlər dövründə populyarlıq itirir.

### 14. Mərkəzləşdirilmiş Xidmətlər
- **pkg.go.dev** — modulların godoc + lisenziya + README + asılılıq indeksi.
- **Module proxy** (default: Google) — bütün public modulların versiya kopyaları;
  yoxdursa repo-dan çəkib saxlayır. Modulun repo-dan silinməsindən qoruyur.
- **Sum database** — hər modul versiyasının imzalı hash ağacı; `go build/test/get`
  hər endirmədə hash müqayisə edir — tutuşmazsa modul QURAŞDIRILMIR (könüllü re-tag və
  zərərli kod dəyişikliyinə qarşı).
- **Alternativlər:** `GOPROXY=https://gocenter.io,direct` (JFrog), `GOPROXY=direct`
  (proxy-siz; silinmiş versiya riski), öz proxy-n (Artifactory/Sonatype/Athens) —
  private repo-lar üçün ideal (auth mərkəzləşir, cache sürətlənir).
- **GOPRIVATE** — public proxy istifadə edərkən private repo-ları istisna etmək:
  `GOPRIVATE=*.example.com,company.com/repo`.

### 15. Modul Nəşri
Sadəcə VCS-ə qoy (+ go.mod, go.sum, LICENSE, README, SemVer teqləri). Mərkəzi registry
YOXDUR (Maven/npm kimi upload lazım deyil). **Lisenziya:** öz lisenziyanı YAZMA;
icma permissive lisenziyaları üstünlük tutur (BSD, MIT, Apache) — Go üçüncü tərəf kodu
birbaşa binary-yə kompilyasiya etdiyindən GPL kimi nonpermissive istifadəçiləri məcburi
open-source edir — bir çox təşkilat üçün qəbuledilməz.

## Əsas terminlələr
- Module path (modul yolu) — modulun qlobal unikal identifikatoru
- Package clause (paket bəndi) — faylın ilk sətri: `package NAME`
- Blank import (boş import) — `_ "pkg"` — yalnız init-i işə salır
- Type alias (tip aliası) — `type B = A` — eyni tipin yeni adı
- Semantic versioning (semantik versiyalama) — major.minor.patch qaydaları
- Minimum version selection (minimum versiya seçimi) — ən yeni tələb olunan götürülür
- Import compatibility rule — minor/patch 100% backward-compatible olmalıdır
- Sum database (cəm verilənlər bazası) — imzalı modul hash-lərinin mərkəzi qeydi
- Vendoring (kənarlaşdırma) — asılılıqların vendor/ qovluğunda saxlanması

## Praktik nəticə

Layihə qərarları: (1) paket adları mənalı olsun, `util`/`common` QADAĞAN; (2) cmd/pkg
strukturunu ehtiyac yarananda tətbiq et; (3) modul daxili paylaşım — `internal`; (4)
export hər şeyin sənədi yazılsın (godoc formatında); (5) init-dən qaç, registry-ləri
açıq qeyd et; (6) cycle → paketləri birləşdir; (7) API dəyişikliyi → alias + backward
uyğunluq, major break → /v2 + yeni import yolu; (8) go.mod/go.sum VCS-də; (9) private
repo-lar üçün GOPRIVATE və ya öz proxy.

## Mənbə
Pages: 253-288
