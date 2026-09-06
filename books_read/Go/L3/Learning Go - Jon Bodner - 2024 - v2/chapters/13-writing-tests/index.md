# Chapter 13 — Writing Tests (Testlərin Yazılması)

## Bu chapter nədən bəhs edir?

testing paketi və go test alətləri, test uğursuzluqlarının hesabatı, TestMain,
t.Cleanup, testdata, packagename_test (black-box), go-cmp, table testlər, code
coverage, benchmark-lər, stub/mock patternləri, httptest, integration testlər (build
tags) və race checker.

## Əsas fikirlər

### 1. Test Əsasları
**Struktur:** Test kodu production kodu ilə **eyni paketdə** (`adder_test.go`) —
unexported funksiyaları da test etmək mümkündür. Funksiyalar `Test` prefiksi +
`*testing.T` parametri; qaytarma YOX; ad testin sənədidir (`Test_addNumbers` —
unexported üçün alt xətt konvensiyası).

**Kitabdan kod nümunəsi:**
```go
func addNumbers(x, y int) int { return x + x }   // bug: x+x!

func Test_addNumbers(t *testing.T) {
    result := addNumbers(2, 3)
    if result != 5 {
        t.Error("incorrect result: expected 5, got", result)
    }
}
```
```bash
go test            # FAIL → düzəlt → PASS
go test ./...      # bütün alt qovluqlar
go test -v         # detallı çıxış
go test -count=1   # cache-i ötür (nəticə keşi keçərli testlər üçün saxlanılır)
```

### 2. Uğursuzluq Hesabatı — Error/Fatal
- `t.Error` / `t.Errorf` — uğursuzluq qeyd et, test DAVAM edir (müstəqil yoxlamalar
  üçün — struct sahələri kimi; bütün problemləri bir dəfə görmək).
- `t.Fatal` / `t.Fatalf` — qeyd et + dərhal ÇIX (sonrakı yoxlamalar həmişə düşəcəksə
  və ya panic törədəcəksə). Yalnız cari test funksiyasını bitirir — digər testlər işə düşür.

### 3. TestMain və t.Cleanup
**TestMain** — paket səviyyəli setup/teardown (yalnız 1 dənə; hər test üçün YOX):
```go
func TestMain(m *testing.M) {
    // setup: xarici repo (DB), package-level dəyişənlərin ilkinləşdirilməsi
    exitVal := m.Run()
    // teardown
    os.Exit(exitVal)
}
```
Package-level dəyişən üçün TestMain lazımdırsa — kodu refactor etmək haqqında düşün.

**t.Cleanup** — test-ə məxsus təmizləmə; helper funksiyalarda ideal (defer testin
özündə deyil, yaradıcıdadır); çoxdəfə çağırıla bilər — LIFO sıra:
```go
func createFile(t *testing.T) (string, error) {
    f, err := os.Create("tempFile")
    if err != nil { return "", err }
    t.Cleanup(func() { os.Remove(f.Name()) })
    return f.Name(), nil
}
```

### 4. testdata
Nümunə data üçün rezervə edilmiş qovluq adı; go test cari iş qovluğunu paket qovluğuna
dəyişir → **mütləq nisbi yol** istifadə et: `"testdata/data.txt"`.

### 5. Public API Testi — packagename_test
Eyni qovluqda fərqli paket adı: `package adder_test` — yalnız export olunanlarla işləmək
(mütləq import tələb edir). Black-box test: həqiqi istifadəçi perspektivi. İki paket adı
eyni qovluqda qarışığa bilər.

### 6. go-cmp ilə Müqayisə
`reflect.DeepEqual`-dən yaxşısı — Google-un go-cmp modulu; fərqi dəqiq göstərir:
```go
if diff := cmp.Diff(expected, result); diff != "" {
    t.Error(diff)
}
```
Custom sahələri iqnor etmək üçün comparator (simmetrik, deterministik, təmiz funksiya
olmalıdır):
```go
comparer := cmp.Comparer(func(x, y Person) bool {
    return x.Name == y.Name && x.Age == y.Age   // DateAdded iqnor
})
cmp.Diff(expected, result, comparer)
```

### 7. Table Testlər
**Pattern:** anonim struct slice-u + `t.Run` subtestləri:
```go
data := []struct {
    name     string
    num1, num2 int
    op       string
    expected int
    errMsg   string
}{
    {"addition", 2, 2, "+", 4, ""},
    {"bad_division", 2, 0, "/", 0, `division by zero`},
}
for _, d := range data {
    t.Run(d.name, func(t *testing.T) {
        result, err := DoMath(d.num1, d.num2, d.op)
        if result != d.expected {
            t.Errorf("Expected %d, got %d", d.expected, result)
        }
        // errMsg müqayisəsi...
    })
}
```
Hər subtestin adı çıxışda görünür (`-v`). Error mesaj müqayisəsi kövrəkdir — custom tip
varsa `errors.Is`/`errors.As` işlət.

### 8. Code Coverage
```bash
go test -v -cover -coverprofile=c.out
go tool cover -html=c.out   # brauzerdə: yaşıl=örtülmüş, qırmızı=örtülməmiş
```
100% coverage BUG YOXDUR demək deyil — kitabın öz nümunəsində `*` operatoru `num1+num2`
kimi yazılıb, bütün sətirlər örtülərkən bug qalırdı (copy-paste təhlükəsi!). Coverage
lazımlı, amma kifayət deyil.

### 9. Benchmark-lər
`Benchmark` prefiksi + `*testing.B`; `b.N` loopu MÜTLƏQ; framework N-i artıraraq stabil
nəticə axtarır:
```go
var blackhole int   // kompilyatorun çağırışı optimize-away etməsin deyə

func BenchmarkFileLen(b *testing.B) {
    for _, v := range []int{1, 10, 100, 1000, 10000, 100000} {
        b.Run(fmt.Sprintf("FileLen-%d", v), func(b *testing.B) {
            for i := 0; i < b.N; i++ {
                result, err := FileLen("testdata/data.txt", v)
                if err != nil { b.Fatal(err) }
                blackhole = result
            }
        })
    }
}
```
```bash
go test -bench=. -benchmem
# BenchmarkFileLen/FileLen-10000-12   62992   19000 ns/op   10376 B/op   4 allocs/op
```
Çıxış sütunları: ad-GOMAXPROCS | təkrar sayı | ns/op | B/op | allocs/op.

**Nümunə dərsləri:** buffer böyükdikcə sürət artır (allocation azalır), fayldan böyük
buffer isə artıq yaddaş; bufferi loopdan çıxarıb bir dəfə ayırmaq 65208 alloc/op → 4
alloc/op etdi. **Optimizasiyaya başlamazdan əvvəl:** proqram business tələbləri qarşılayırsa,
vaxtını feature/bug-ə sərf et. Benchmark sürət problemi göstərəndən SONRA profiling
(pprof; Julia Evans məqaləsi).

### 10. Stub-lər (Dependency Abstraction + Test)
**Funksiya tipi / interface asılılıqları test ediləbiləndir.** Tək metodlu interface
stub-u:
```go
type MathSolverStub struct{}
func (ms MathSolverStub) Resolve(ctx context.Context, expr string) (float64, error) {
    switch expr {
    case "2 + 2 * 10": return 22, nil
    ...
}
p := Processor{MathSolverStub{}}   // injekt et
```

**Böyük interface + embed pattern** (yalnız testə lazım olan metodu yaz):
```go
type GetPetNamesStub struct {
    Entities      // embed — bütün metodlar "var", çağırılmayanlar panic verir!
}
func (ps GetPetNamesStub) GetPets(userID string) ([]Pet, error) { ... }
```
Diqqət: çağırılan HƏR metod implement olunmalı, yoxsa panic.

**Funksiya sahəli stub (ən çevik)** — hər test case öz implementasiyasını verir:
```go
type EntitiesStub struct {
    getPets func(userID string) ([]Pet, error)
    // ... hər metod üçün sahə
}
func (es EntitiesStub) GetPets(userID string) ([]Pet, error) {
    return es.getPets(userID)
}
// table test-də:
data := []struct{ ...; getPets func(string) ([]Pet, error); ... }{
    {"case1", func(userID string) ([]Pet, error) { return []Pet{{Name: "Bubbles"}}, nil }, ...},
}
l.Entities = EntitiesStub{getPets: d.getPets}
```

**Mock vs Stub (Martin Fowler):** stub — verilmiş inputa hazır dəyər qaytarır; mock —
çağırışların gözlənilən ardıcıllıq/inputlarla BAŞ VERDİYİNİ təsdiqləyir. Alətlər:
gomock (Google), testify (Stretchr).

### 11. httptest — HTTP Servis Stub-u
**Nədir:** Random portda real HTTP server — xarici servisi imitasiya edir, integration
ağrısız:
```go
type info struct { expression string; code int; body string }
var io info   // iki closure arasında paylaşılır (test kodunda qəbul olunur)

server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
    expression := req.URL.Query().Get("expression")
    if expression != io.expression {
        rw.WriteHeader(http.StatusBadRequest)
        return
    }
    rw.WriteHeader(io.code)
    rw.Write([]byte(io.body))
}))
defer server.Close()
rs := RemoteSolver{
    MathServerURL: server.URL,    // test server-in URL-i
    Client:        server.Client(), // hazır client
}
```

### 12. Integration Testlər və Build Tags
Integration testlər — real xarici servisə bağlanan testlər (API anlayışının təsdiqi).
Qruplaşdırma — build tag (faylın ilk sətri, package-dən əvvəl boş sətirlə ayrılmış
magic comment):
```go
// +build integration
```
```bash
go test -tags integration -v ./...
```
Tagsiz fayllar həmişə işləyir (unit), tagli fayllar yalnız `-tags` ilə.

**-short flag alternativi:** `if testing.Short() { t.Skip(...) }` — amma kitab tövsiyə
ETMİR: yalnız 2 səviyyə, "asılılıq" (tag) ilə "uzun" (short) fərqli anlayışlardır,
və default-da qısa testlər işləməliyi, `-short`-u "çixartmaq" mənasızdır.

### 13. Race Checker
**Data race** — kilidsiz paralel dəyişən çıxışı (`counter++` 5 goroutine-dən:
gözlənti 5000, real 3673). Axtarış:
```bash
go test -race    # xətamız yoxlamanı göstərir
```
Çıxış: goroutine-lərin oxu/yazma izləri + sətir. **Sleep əlavə etməklə race düzəltmə —
kod əvvəlki kimi səhvdir.** `-race` binary də mümkündür (testsiz kod üçün), amma ~10x
yavaş — ona görə həmişə açıq deyil. Tapılan hər race MÜTLƏQ düzəldilməlidir.

## Əsas terminlər
- Table test (cədvəl testi) — anonim struct slice + t.Run subtestləri
- Coverage (örtük) — testlərin toxunduğu kod nisbəti
- Benchmark — b.N looplu performans ölçmə funksiyası
- Stub — hazır dəyər qaytaran saxta asılılıq
- Mock — çağırış ardıcıllığını təsdiqləyən saxta (Fowler fərqi)
- Build tag — fayl-ərzində kompilyasiya şərti (`// +build`)
- Data race — kilidsiz paralel yaddaş çıxışı
- httptest.NewServer — random portlu test HTTP server

## Praktik nətidə

Test yazma axını: (1) kiçik funksiya + test; (2) uğursuzluqda Error/Fatal seçimi; (3)
müqayisələr go-cmp ilə; (4) çoxsaylı hallar → table test; (5) coverage bax, amma
100%-ə sitayiş etmə; (6) performans şübhəsi → -benchmem benchmark (əvvəl ehtiyac
sübutu!); (7) asılılıqları interface/funksiya tipinə çək → stub-lar (funksiya sahəli
stub ən saxta-dostu); (8) HTTP asılılıqları httptest ilə; (9) integration = build tag;
(10) paralel kod həmişə -race ilə test et; (11) black-box üçün _test paketindən istifadə
et.

## Mənbə
Pages: 367-404
