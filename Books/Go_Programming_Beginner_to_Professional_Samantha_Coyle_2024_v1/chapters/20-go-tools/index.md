# Chapter 20 — Using Go Tools (səh. 628-645)

## Bu fəsil nədən bəhs edir?

Go toolchain: go build, go run, gofmt (format), goimports (import
idarəetməsi), go vet (static analysis), race detector (--race), go doc
(sənədləşdirmə) və go get (xarici paketlər).

## Əsas fikirlər

### 1. go build — kompilyasiya
```bash
go build -o hello_world main.go
./hello_world          # Hello World
```
- `-o` — çıxış binary adı; yoxsa paket/qovluq adı ilə

### 2. go run — compile + icra bir addımda
```bash
go run main.go        # Hello Packt
```
- Binary YARATMIR — tez yoxlama üçün ideal

### 3. gofmt — format
```go
// PİS format:
func main(){
  firstVar := 1
    secondVar := 2
  fmt. Println("Hello Packt")
}
```
```bash
gofmt main.go         # NÜMAYİŞ (fayl dəyişmir)
gofmt -w main.go      # YAZ — faylı düzəlt
```
- Boşluq/indentasiya/{} düzəlişləri; bütün Go layihələrində EYNİ
  stil → oxunaqlılıq
- IDE-lərdə save-zamanı avtomatik rejim var

### 4. goimports — import idarəetməsi
```go
// XƏTALI: net/http istifadəsiz, log yoxdur:
import (
    "net/http"
    "fmt"
)
func main() {
    fmt.Println("Hello")
    log.Println("Packt")   // log import edilməyib!
}
```
```bash
goimports main.go      # nümayiş: http silinir, log əlavə olunur
goimports -w main.go   # yaz
```
- Əlavə edir + İSTİFADƏSİZİ SİLİR + əlifba sırası
- İstifadəsiz import = təhlükəsizlik riski

### 5. go vet — static analysis
```go
jointString := fmt.Sprintf("%s", helloString, packtString)
// COMPILE KEÇİR, amma packtString İTİR!
```
```bash
go vet main.go
# main.go:9: Sprintf call needs 1 arg, has 2 args

# Fix: fmt.Sprintf("%s %s", helloString, packtString)
go vet main.go        # (boş) → problem yoxdur
```
- Tutduğu səhvlər: Printf arqument sayı, useless assignments, unmarshal-a
  non-pointer ötürmə, unreachable code
- Build prosesinə daxil etmək tövsiyə olunur

### 6. Race detector
**Nədir:** asinxron yaddaş girişlərini icra vaxtı aşkarlayan alət —
yalnız KOD İCRA OLUNANDA işləyir.

```go
// YARIS ŞƏRAİTLİ kod:
func main() {
    finished := make(chan bool)
    names := []string{"Packt"}
    go func() {
        names = append(names, "Electric")    // YAZ (goroutine)
        names = append(names, "Boogaloo")
        finished <- true
    }()
    for _, name := range names {              // OXU (main) — eyni anda!
        fmt.Println(name)
    }
    <-finished
}
```
```bash
go run --race main.go
# ==================
# WARNING: DATA RACE
# Write at 0x00c0000aa000 by goroutine 6:
#   main.main.func1() main.go:10 +0xe0
# Previous read at 0x00c0000aa000 by main goroutine:
#   main.main() main.go:14 +0x170
# Found 1 data race(s)
```
- Warning: hansı goroutine, hansı sətirlər, oxu/yaz — hamısı göstərilir

**Fix — sırala (kanalla senxronlaşdır):**
```go
<-finished                            // ƏVVƏL bitməsini gözlə
for _, name := range names {          // SONRA oxu
    fmt.Println(name)
}
// go run --race → warning YOXDUR; çıxış: Packt Electric Boogaloo
```
- Professional praktika: həllərin race-sızlığını --race ilə təsdiqlə

### 7. go doc — sənədləşdirmə
```go
// Add returns the total of two integers added together
func Add(a, b int) int {
    return a + b
}

// Multiply returns the total of one integer multiplied by the other
func Multiply(a, b int) int {
    return a * b
}
```
```bash
go doc -all                  # paketin tam sənədi
godoc package/path > out.txt # fayla yaz
```
- Konvensiya: şərh FUNKSİYA ADI ilə başlayır → go doc onu götürür
- Komandalararası əlaqə üçün: sənəd paylaş, import edən istifadə etsin

### 8. go get — xarici paketlər
```go
import "github.com/gorilla/mux"

r := mux.NewRouter()
r.HandleFunc("/", exampleHandler)
log.Fatal(http.ListenAndServe(":8888", r))
```
```bash
go run main.go
# error: no required module provides package github.com/gorilla/mux

go get github.com/gorilla/mux     # ENDİR + go.mod-a yaz
go run main.go                    # işləyir
```

### Activity 20.01 — bütün alətlərin kombinasiyası
Pis formatlı + xətalı faylı düzəlt: gofmt (format) → goimports
(missing/unused importlar) → go vet (return-sonrası log.Println =
unreachable; Func → func) → go get (mux) → işə sal.

## Əsas terminlər
- go build -o / go run
- gofmt (-w) — üslub düzəltməsi
- goimports (-w) — import əlavə/sil/sırala
- go vet — compiler-in buraxdığı səhvlər
- --race flag — data race detektoru (icra vaxtı)
- Race condition — paralel oxu/yazı qarşılıqlı toxunması
- go doc -all / godoc — şərhlərdən sənəd
- Şərh konvensiyası — funksiya adı ilə başlayan doc comment
- go get — xarici asılılığın endirilməsi
- gorilla/mux — router paketi (nümunə)

## Praktik nəticə
Gündəlik iş axını: kod yaz → goimports (save-də avtomatik) → gofmt
(eyni) → go vet (CI-də məcburi) → go test --race (konkurent koddan
əvvəl) → go build. Sənədləşdirmə: hər export funksiyaya "FunksiyaAdı
..." şərhi → go doc -all ilə generasiya. Xarici paket yalnız go get
ilə — əks halda "no required module" xətası. Race warning oxunu:
goroutine ID + fayl:sətir + oxu/yaz yeri → həmin resursu kanalla/
mutexlə senxronlaşdır.

## Mənbə
Pages: 628-645 (PDF 628-645)
