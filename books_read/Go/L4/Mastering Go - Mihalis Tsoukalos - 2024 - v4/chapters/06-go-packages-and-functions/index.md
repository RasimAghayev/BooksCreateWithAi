# Chapter 6 — Go Packages and Functions (Go Paketləri və Funksiyaları)

## Bu chapter nədən bəhs edir?

Paket anlayışı və görünürlük qaydası (public/private), go get vs go install, funksiyalar
(anonim, çoxlu qaytarma, adlı return, funksiya-parametr, funksiya qaytaran, variadic),
defer (LIFO + closure tələsi), Big O, init() və icra sırası, GitHub-da paket saxlama,
SQLite3 paketi (database/sql + go-sqlite3), modullar/SemVer, paket keyfiyyət qaydaları,
sənədləşdirmə (go doc, BUG), Workspaces (go.work) və ldflags ilə versiyalama.

## Əsas fikirlər

### 1. Paketlər və Görünürlük Qaydası
**Nədir:** Hər Go faylı `package AD` ilə başlayır; main istisna olmaqla paketlər müstəqil
proqram DEYİL (`go run` işləmir: "cannot run non-main package").

**Görünürlük qanunu:** BÖYÜK hərflə başlayan hər şey = public (package xaricindən çatılır);
kiçik hərf = private. Strukt sahələrinə də şamil olunur. Yeganə istisna: paket adları
(adətən kiçik hərflə).

**Dizayn prinsipi:** paketi "yaza bildiyin üçün" YOX — eyni kod başqa proqramlarda lazım
olanda və ya funksiyalar qruplanmalı olanda yarat. Səbəblər: encapsulation, test, security.

### 2. go get vs go install
```bash
go install github.com/mattn/go-sqlite3@latest   # Go 1.16+ — TÖVSİYƏ OLUNAN
go get github.com/spf13/cobra                    # deprecated (modul daxilində əvvəlki rol)
go get -u github.com/spf13/viper                 # upgrade
go mod init / go mod tidy                        # öz layihən üçün əsl yol
```
**~/go strukturu:** `bin/` (alət binary-ları), `pkg/` (compile olunmuş/modul paketlər),
`src/` (köhnə qaydada paket mənbələri — modullar `~/go/pkg/mod`-a düşür).

### 3. Funksiyalar — Birinci Sinif Vətəndaşlar
**Qayda:** funksiya BİR işi yaxşı görsün; bir neçə iş görürsə — bölün.

**Kitabdan kod nümunələri:**
```go
// Çoxlu qaytarma (>3 dəyərdə struct/slice düşün):
func doubleSquare(x int) (int, int) {
    return x * 2, x * x
}

// Adlı return — boş return adlı dəyərləri qaytarır:
func minMax(x, y int) (min, max int) {
    if x > y { min = y; max = x; return min, max }
    min = x; max = y
    return                       // = return min, max
}

// Anonim funksiya dəyişəndə:
anF := func(param int) int { return param * param }

// Funksiya parametr kimi (sort.Slice):
sort.Slice(data, func(i, j int) bool { return data[i].Grade < data[j].Grade })
isSorted := sort.SliceIsSorted(data, func(i, j int) bool { ... })

// Funksiya QAYTARAN funksiya:
func funRet(i int) func(int) int {
    if i < 0 {
        return func(k int) int { k = -k; return k + k }
    }
    return func(k int) int { return k * k }
}
```

### 4. Variadic Funksiyalar
**Qaydalar:** `...T` (pack operatoru) — funksiyada BİR dəfə, həmişə SONDA; daxildə slice
kimi istifadə olunur; çağırışda siyahı və ya `slice...` (unpack).

**Kitabdan kod nümunəsi:**
```go
func addFloats(message string, s ...float64) float64 {
    sum := float64(0)
    for _, a := range s { sum = sum + a }
    s[0] = -1000          // slice kimi element access mümkün
    return sum
}
addFloats("msg", 1.1, 2.12, 3.14)     // inline
addFloats("msg", s...)                // unpack

// []string → []interface{} CEVRIM MECBURIDIR:
empty := make([]interface{}, len(os.Args[1:]))
for i, v := range os.Args[1:] { empty[i] = v }
everything(empty...)
```
**Tələ:** `os.Args...` → `[]string` ≠ `[]interface{}` — memory repr fərqlidir, compile xətası.
`everything(os.Args)` işləyir, amma BÜTÜN slice-i TƏK element kimi ötürür.

### 5. defer — LIFO və Closure Tələsi
**Nədir:** Funksiyanın icrasını əhatə edən funksiya qayıdana qədər təxirə salır. LIFO
sırası ilə icra olunur. Fayl açılışının yanında Close yazmaq üçün ideal.

**Kitabdan kod nümunəsi (3 variant):**
```go
func d1() {
    for i := 3; i > 0; i-- {
        defer fmt.Print(i, " ")     // DƏYƏRLƏR dərhal bağlanır → çap: 1 2 3
    }
}
func d2() {
    for i := 3; i > 0; i-- {
        defer func() { fmt.Print(i, " ") }()   // CLOSURE — i sonrakı qiymətlə ölçülür → 0 0 0!
    }
}
func d3() {
    for i := 3; i > 0; i-- {
        defer func(n int) { fmt.Print(n, " ") }(i)   // PARAMETR kimi → 1 2 3 — DÜZGÜN YOL
    }
}
```
**Sub-kod izahı:** d2-də closure i-ni referans saxlayır; loop bitəndə i=0 → 3 dəfə 0.
Go 1.22-də loop dəyişəni semantikası dəyişdi (hər iterasiya öz dəyişəni) — amma parametr
ötürmə hələ də ən aydın üsuldur. **Qayda: defer-də dəyəri PARAMETR kimi ötür.**

### 6. Big O Kompleksliyi
| Notasiya | Məna | 100 element üçün əməliyyat |
|---|---|---|
| O(1) | sabit | 1 |
| O(n) | linear — yaxşı hesab olunur | ~100 |
| O(n²) | kvadratik | ~10,000 |
| O(n!) | faktorial | 10^158! |

### 7. init() və İcra Sırası
**init() xüsusiyyətləri:** arqumentsiz, dəyərsiz, optional, implicit çağırılır, main-dən
ƏVVƏL; faylda birdən çox ola bilər (bəyannamə sırası ilə); paket neçə dəfə import olunsa
da BİR dəfə icra olunur; xaricdən çağırıla BİLMƏZ (private by design).

**İcra sırası (main → A → B asılılığı):**
1. main → A → B importlar çözülür
2. B qlobal dəyişənləri → B init()
3. A qlobal dəyişənləri → A init()
4. main qlobal dəyişənləri → main init()
5. main() başlayır

**Məqsədli istifadə halları:** şəbəkə/server bağlantılarının init-i, lazımı fayl/qovluq
yaratma, resurs mövcudluğu yoxlaması. **Ehtiyat:** public paketdə global state dəyişmə.

### 8. SQLite3 ilə İş — database/sql + go-sqlite3
**Kitabdan kod nümunəsi (temel əməliyyatlar):**
```go
import (
    "database/sql"
    _ "github.com/mattn/go-sqlite3"   // BLANK import — yalnız driver qeydiyyatı
)

db, err := sql.Open("sqlite3", "test.db")   // fayl yoxdursa yaradılır
defer db.Close()

// Cədvəl yarat:
_, err = db.Exec(`CREATE TABLE IF NOT EXISTS book (
    id INTEGER NOT NULL PRIMARY KEY, time TEXT NOT NULL, description TEXT);`)

// Prepare + Exec (? parametrli):
stmt, err := db.Prepare("INSERT INTO book VALUES(NULL,?,?);")
_, err = stmt.Exec(cT, dsc)

// Bir sətirlik sorğu:
var version string
err = db.QueryRow("SELECT SQLITE_VERSION()").Scan(&version)

// Çoxsətirlik sorğu:
rows, err := db.Query("SELECT * from book WHERE id > ?", n)
defer rows.Close()
for rows.Next() {
    var id int; var dt, description string
    err = rows.Scan(&id, &dt, &description)
    // ...
}

// Update/Delete də Exec ilə:
db.Exec("UPDATE book SET time = ? WHERE id > ?", cT, 7)
stmt, _ = db.Prepare("DELETE from book where id = ?")
stmt.Exec(8)
```
**Sub-kod izahı:**
- `_ "pkg"` blank import — paket özünü `sql` üçün driver kimi qeyd edir, kodda birbaşa
  istifadə olunmur
- `Prepare` + `Exec` — parametrli sorğu; `?` yerləri Exec arqumentləri ilə doldurulur
- `QueryRow` — ≤1 sətir gözlənilən; `Query` — çox sətir; hər ikisi `Scan` ilə oxunur
- go-sqlite3 cgo istifadə edir → gcc tələb olunur
- SQLite lokal tək fayldır — server/port yoxdur

### 9. Öz Paketin — sqlite06 Dizaynı
**Struktur (kitabın nümunəsi):** Users + Userdata cədvəlləri; UserID ilə bağlı.
Paket funksiyaları:
- `openConnection()` — private helper
- `exists(username) int` — private; -1 = yoxdur (indikator dəyəri — sadəlik; error daha
  robust-dur, Go fəlsəfəsinə daha yaxın)
- `AddUser(d Userdata) int` — public; -1 = xəta
- `DeleteUser(id int) error`, `UpdateUser(d Userdata) error`, `ListUsers() ([]Userdata, error)`

**Dizayn dərsləri:**
- İki cədvəl üçün iki struct məcburi DEYİL — birləşdirilmiş Userdata daha praktik oldu
- Username lowercase — dublikatların qarşısı (dizayn qərarı)
- DeleteUser iki addımdır: əvvəl Userdata, sonra Users (xarici açar ardıcıllığı)
- Debug Println-lər final versiyada error dəyərlərinə çevrildi — xəta qərarı istifadəçiyə
  aiddir (library-də log/çap YOX)

### 10. Modullar və SemVer
**Module** = versiyalı paket toplusu (bir paketdən çox ola bilər). vMAJOR.MINOR.PATCH:
- MAJOR — backward-INCOMPATIBLE
- MINOR — yeni feature, uyğun
- PATCH — yalnız bugfix
Modullar Go 1.11-də gəldi, 1.13-də finalize olundu.

### 11. Yaxşı Paket Qaydaları (kitabdan)
1. Elementlər məntiqi bağlı olsun (maşın paketi — amma maşın+təyyarə+pvelosiped bir paketdə YOX)
2. Public elan etməzdən əvvəl özün uzun müddət istifadə et; testlər MÜTLƏQ
3. Aydın, tez məhsuldar API
4. Public API minimal; adlar deskriptiv-amma-qısa
5. Uyğun halda interface/generic parametr tipi
6. Update-lərdə köhnə versiya uyğunluğunu qırma
7. Birləşdirilmiş tapşırıqları ayrı-ayrı fayllara böl
8. Mövcud paketi sıfırdan yazma — töhfə ver/öz versiyanı yarad
9. Paket ekrana log çap etməsin — flag ilə idarə olunsun
10. Tip definition-lərı istifadə yerinə yaxın qoy
11. Test faylı olan paket = peşəkar paket; testsiz paketi istifadə ETMƏ

### 12. Sənədləşdirmə — go doc
**Qaydalar:** komment elementin DƏRHAL üstündə, boş sətir YOX; funksiya kommenti funksiya
adı ilə başlayır; block comment (`/* */`) paket intro üçün; tab-başlayan sətirlər kod kimi
render olunur; `BUG(1): ...` xüsusi açar sözü.

```go
/*
The package works on 2 tables on an SQLite database.
...
*/
package document

// AddUser adds a new user to the database
//
// Returns new User ID
// -1 if there was an error
func AddUser(d Userdata) int { }
```
`go doc document.go ListUsers` — yalnız PUBLIC elementlər göstərilir.

### 13. Workspaces — go.work
**Nədir:** Eyni anda bir neçə modul üzərində iş; `use` və `replace` direktivləri go.mod-u
override edir.

**Kitabdan kod nümunəsi:**
```bash
go work init .
go work use ./util
go work use ./sqlite06
# go.work:
# go 1.21.0
# use (./sqlite06, ./util)
replace github.com/mactsouk/sqlite06 => ./sqlite06   # ən vacib sətir — lokal kopyaya yönləndir
go run ./util/sqliteGo.go     # lokal DƏYİŞDİRİLMİŞ versiya icra olunur
```
**İstifadə halı:** stabil versiya sistemdə qalarkən lokal inkişaf nüsxəsi ilə işləmək.

### 14. Versiyalama — ldflags
**Kitabdan kod nümunəsi:**
```go
var VERSION string    // runtime-da linker tərəfindən doldurulacaq

func main() {
    if len(os.Args) == 2 && os.Args[1] == "version" {
        fmt.Println("Version:", VERSION)
    }
}
```
```bash
export VERSION=$(git rev-list -1 HEAD)
go build -ldflags "-X main.VERSION=$VERSION" gitVersion.go
./gitVersion version   # Version: 4dc3d6b5... (commit hash — unikal)
```
**Sub-komanda izahı:** `-ldflags` — cmd/link-ə dəyər ötürür; `-X main.VERSION=...` —
compile vaxtı dəyişənə dəyər yazır. Docker/kubectl də eyni texnikanı işlədir.

## Əsas terminlər
- Package (paket) — `package AD` ilə başlayan təşkilat vahidi
- Public/Private — böyük/kiçik hərflə görünürlük
- Blank Import — `_ "pkg"` — yalnız init/driver qeydiyyatı üçün
- Anonymous Function (anonim funksiya) — adsız; lambda
- Closure (bağlanma) — leksik scope dəyişənlərini capture edən funksiya
- Named Return (adlı qaytarma) — imzada adlandırılmış nəticə; boş return
- Variadic Function — `...T` dəyişən saylı parametr
- Pack/Unpack — `...T` / `slice...`
- defer — LIFO təxirəsalma; resurs təmizliyi
- init() — paket inicializasiyası; main-dən əvvəl, bir dəfə
- Big O — alqoritm artım sırası
- Module (modul) — SemVer ilə versiyalı paket toplusu
- Workspace (iş sahəsi) — go.work ilə çoxmodullu iş
- ldflags -X — build vaxtı dəyər inyeksiyası (versiyalama)

## Praktik nəticə

(1) Böyük hərf = public — API-ni minimal saxla; hər şey öz-özünü izah etsin. (2) Funksiya
BİR iş; >3 qaytarma → struct. (3) defer-də closure YOX — dəyəri parametr kimi ötür (d3
nümunəsi); LIFO sırasını yadda saxla. (4) `[]string` → `[]interface{}` avtomatik YOX —
manual çevrim. (5) init() nüvəsi: main-dən əvvəl, bir dəfə; public paketdə ehtiyatla.
(6) database/sql pattern-i: Open+defer Close, Prepare+Exec, Query+rows.Next+Scan, QueryRow.
(7) Driver paketi blank import. (8) Paket = müqavilə: daxili dəyişikliklər istifadəçini
etməməli; xətalar error kimi qayıtsın, çap YOX. (9) Testsiz paket peşəkar deyil — testlər
opsiya deyil. (10) Versiya = commit hash + ldflags -X; CI/CD uyğun avtomatik unikallıq.

## Mənbə
Pages: 201-263 (PDF 232-295)
