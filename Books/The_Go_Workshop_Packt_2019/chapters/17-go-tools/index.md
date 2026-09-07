# Chapter 17 — Using Go Tools (Go Alətləri)

## Bu fəsil nədən bəhs edir?

Go alətlər dəsti: go build (binary yaratma), go run (compile+icra bir addımda),
gofmt (kod formatı), goimports (import idarəsi), go vet (statik analiz — compiler-in
buraxdığını tutma), race detector (`--race`), go doc (dokumentasiya generasiyası),
go get (3-cü tərəf paket endirmə) və kombinə activity.

## Əsas fikirlər

### 1. go build vs go run
```bash
go build -o hello_world main.go    # BINARY yaradır → ./hello_world
go run main.go                      # compile + İCRA, binary YOX — sürətli test
```
build: paylanacaq executable; run: development sürət yolu.

### 2. gofmt — Kod Formatı
```bash
gofmt main.go          # necə GÖRMƏLİ olduğunu göstərir (dəyişməz)
gofmt -w main.go       # FAYLI DƏYİŞTİRİR + saxlayır
```
- Spacing, indentation, alignment avtomatik düzəlir
- `func` / `main(){` split → bir sətirdə birləşdirilir
- BÜTÜN Go layihələrində EYNİ stil — oxunaqlılıq standartı
- IDE-lərdə save-də avtomatik işə salına bilər

### 3. goimports — Import İdarəsi
```bash
goimports main.go     # preview
goimports -w main.go  # yaz
```
- LAZIMSIZ importu SİLİR (security risk!)
- ÇATIŞMAYANI ƏLAVƏ EDİR (net/http istifadə olunmayıb → sil; log çatışmır → əlavə)
- Importları ƏLİFBA sırasına düzür
- Kod yazarkən import haqqında DÜŞÜNMƏ — alət həll edir

### 4. go vet — Statik Analiz
**Nədir:** compiler-in BURAXDIĞI xətaları tutan analiz aləti — kod İŞLƏYİRSA BELƏ.

**Tutduqları:**
- **Printf arqument sayı:** `fmt.Sprintf("%s", a, b)` — 1 verb, 2 arqument → xəta!
  (compiler KEÇİR, runtime səhv nəticə verir)
- Lazımsız assignment-lər
- **Unmarshal-a non-pointer ötürmə** — valid kod, amma data YAZILA BİLMİR → debug
  əzabı; vet ERKƏN tutur
- Unreachable kod

```bash
go vet main.go
# main.go:9:2: Sprintf call needs 1 arg but has 2
```
**Praktika:** build prosesinə DAXİL ET — production-a çatmamış tut.

### 5. Race Detector — `--race`
**Nədir:** runtime-da asinxron yaddaş erişimini aşkarlayan alət; kod İCRA
OLUNMADAN tapa bilmir → go run/test ilə birləşdirilib.

**Kitabdan kod nümunəsi:**
```go
finished := make(chan bool)
names := []string{"Packt"}
go func() {
    names = append(names, "Electric")     // YAZ (goroutine)
    names = append(names, "Boogaloo")
    finished <- true
}()
for _, name := range names {               // OXU (main) — EYNİ ANDA!
    fmt.Println(name)
}
<-finished
```
```bash
go run --race main.go
# WARNING: DATA RACE — Read at main.go:15 / Write at main.go:10
```
**Həll — əməliyyat SIRASINI dəyiş:**
```go
<-finished                                 // ƏVVƏL goroutine bitməsini gözlə
for _, name := range names { ... }         // SONRA oxu — race YOX
```
Race-lər TRANSIENT-dirlər — deploy-dan ÇOX SONRA üzə çıxır; detector professional
standart yoxlama alətidir.

### 6. go doc — Dokumentasiya
**Konvensiya:** funksiya/paket üstündə // komment, FUNKSIYA ADI ilə başlayır:
```go
// Add returns the total of two integers added together
func Add(a, b int) int { return a + b }

// Multiply returns the total of one integers multiplied the other
func Multiply(a, b int) int { return a * b }
```
```bash
go doc -all          // paket + bütün funksiyaların sənədi
```
Komment → avtomatik dokumentasiya; komandalararası paylaşım üçün.

### 7. go get — 3-cü Tərəf Paketlər
```bash
go get github.com/gorilla/mux     # paketi LOKAL endirir
```
Importda `github.com/gorilla/mux` görülən kimi paket lokalda OLMAZSA → compile
xətası; go get həll edir. Standart kitabxanadan çox güclü ekosistem.

### 8. Activity — Alətlər Zənciri
Pis format + çatışmayan import + unreachable kod + 3-cü tərəf paket olan faylı:
1. `gofmt -w` — format
2. `goimports -w` — importları düzəlt (http çatışmır!)
3. `go vet` — return-dən SONRAKI log.Println (unreachable) tapılır
4. `go get github.com/gorilla/mux` — paket
5. `go run` → işləyən server

## Əsas terminlələr
- go build / go run — binary / biraddımlıq icra
- gofmt (-w) — kod formatını standartlaşdırma
- goimports (-w) — import əlavə/sil/çeşidləmə
- go vet — statik analiz; Printf arqumentləri, unmarshal pointer
- Race Detector (--race) — runtime data race aşkarlanması
- DATA RACE Warning — oxu+yaz toqquşması raportu
- go doc -all — kommentlərdən sənəd
- Sənəd Konvensiyası — "// FuncName ..." formatı
- go get — 3-cü tərəf paketin lokal endirilməsi
- Unreachable Code — ölü kod; vet tapır

## Praktik nətidə

(1) Sürətli test → go run; paylama → go build -o ad. (2) gofmt/goimports -w — save
avtomatikası qur; import barədə düşünməyi DAYANDIR. (3) go vet-i build-ə daxil et —
Printf arqument xətaları runtime-da SIRLI itir. (4) Unmarshal çağırışında pointer
YOXDURSA vet səni xilas edir — ən çox vaxt aparan bug-lardan. (5) Paralel kodda
`go run --race` MÜTLƏQ — race müvəqqətidir, adi testlərdə gizli qalır. (6) Race
həlli: sadə hallarda əməliyyat sırası (bitməyi gözlə, sonra oxu); mürəkkəbdə
WaitGroup/mutex/channel. (7) go doc üçün kommentlər FunksiyaAdı ilə başlasın —
sənəd avtomatik. (8) go get ilə paketi əVVƏLCƏDƏN endir — compile xətasının qarşısı.
(9) Alətlər zənciri: format → import → vet → get → run — pis kodun əsas təmizlənmə
ardıcıllığı. (10) Go-nun populyarlığının səbəblərindən biri — bütün bu alətlər
DİLİN TƏRKİBİ kimi gəlir.

## Mənbə
Pages: 605-623 (PDF 638-657)
