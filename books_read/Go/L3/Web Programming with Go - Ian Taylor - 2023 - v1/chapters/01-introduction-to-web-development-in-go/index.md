# Chapter 1 — Introduction to Web Development in Go (Go-da Veb İnkişafına Giriş)

## Bu chapter nədən bəhs edir?
Go-nun veb inkişafında yerinə, toolchain-inə, Linux-da qurulumuna, net/http paketinə, GitforGits Bookstore layihəsinə (kitab boyu qurulacaq) və ilk veb server-in yaradılışına. Kitabın praktik əsası.

## Əsas fikirlər

### 1. Go veb inkişafında niyə? (Tarixçə + üstünlüklər)
**Nədir:** 2007-də Google-də Griesemer/Pike/Thompson tərəfindən C++/Java-nın scalability və concurrency çatışmamazlıqlarına cavab olaraq yaradılıb; v1.0 — 2012.

**Veb üçün üstünlükləri:**
- **Concurrency native:** goroutines + channels + select — minlərlə eynizamanlı sorğu
- **Statik tip + kompilyasiya:** C/C++ səviyyəli performans, JVM/CLR overhead-i yox
- **Tək statik binary:** bütün dependency-lər içində — "dependency hell" yox; Docker/K8s-də minimal image
- **Cross-compilation:** bir platformada yaz, başqasına kompilyasiya et
- **Standard library:** HTTP server, DB, crypto — box-dan çıxan
- **Sadəlik + backward compatibility:** uzunömürlü layihələr üçün proqnozlaşdırıla bilənlik

**Real dünya istifadəçiləri:** Docker (konteynerləşdirmə), Kubernetes (orkestrasiya), Twitch (milyonlarla çat mesajı), Uber (mikroservis ekosistem), Dropbox (Python→Go miqrasiyası).

### 2. Go Toolchain
| Alət | Vəzifə |
|------|--------|
| `go build` | kompilyasiya — statik binary (bütün dependency-lərlə) |
| `go install` | binary-nı $GOPATH/bin-ə quraşdırır |
| `go mod` | Go Modules (v1.11+) — go.mod idarəsi, versiyalama, reproducible builds |
| `go get` | paket alma + kompilyasiya |
| `go test` | unit test + coverage |
| `go bench` | benchmark |
| `go tool` | alətlərə birbaşa erişim (compiler, linker, trace) |

### 3. Linux-da qurulum (CLI komandaları)
```bash
# 1. Tarball-ı aç:
sudo tar -C /usr/local -xzf go$VERSION.$OS-$ARCH.tar.gz
# 2. PATH-ə əlavə et:
echo "export PATH=$PATH:/usr/local/go/bin" >> ~/.profile
source ~/.profile
# 3. Workspace (artıq məcburi deyil, amma faydalı):
mkdir -p ~/go_projects/{bin,src,pkg}
echo "export GOPATH=~/go_projects" >> ~/.profile
# 4. Modul init:
go mod init <module-name>
```
**Veb alətləri:**
```bash
go get -u github.com/gorilla/mux     # router
go get -u gorm.io/gorm              # ORM
go get -u github.com/joho/godotenv  # .env faylları
go get github.com/pilu/fresh        # auto-reload dev tool
```
**IDE:** VS Code + Go extension, ya GoLand (Jetbrains).

### 4. net/http paketi — veb inkişafının nüvəsi
**Nədir:** Standard kitabxananın HTTP paketi — server və klienti bir paketdə birləşdirir.

**Komponentlər:**
- **Server:** `http.ListenAndServe(":8080", nil)` — portda dinləmə
- **Klient:** `http.Get`/`http.Post` — outbound sorğular
- **Handler:** `http.Handler` interfeysi — `ServeHTTP(http.ResponseWriter, *http.Request)` metodu
- **HandlerFunc:** adi funksiyanın handler-ə çevrilməsi adapteri
- **Router:** `http.ServeMux` — URL pattern-lərə görə yönləndirmə (mürəkkəb üçün Gorilla Mux)
- **Request/Response:** `*http.Request` — query/form/body oxuma; `http.ResponseWriter` — header/status/cavab yazma

### 5. GitforGits Bookstore — kitab boyu layihə
**Layihənin xüsusiyyətləri:** user authentication, kitab kataloqu, axtarış+filtrlər, səbət, review/rating, admin dashboard.

**Öyrənəcəklərin:** layihə strukturu, DB əməliyyatları (CRUD/migrations/queries), sessiya/cookie idarəsi, templating, input validation, middleware, concurrency, testing, deploy/scaling.

### 6. REST prinsipləri və Go
**REST (Representational State Transfer — nümayəndəlik vəziyyət köçürməsi)** — Dr. Roy Fielding-in arxitektura üslubu:
- **Statelessness (vəziyyətsizlik):** hər sorğu BÜTÜN lazımı məlumatı daşıyır
- **Client-server ayrılığı:** UI vs data/backend
- **Cacheability (keşlənə bilmə):** cavablar açıq şəkildə cache-ə uyğun/markəli
- **Layered system (qatlı sistem):** klient son serveri bilə bilməz (intermediarylər şəffaf)
- **Uniform interface (vahid interfeys):** standartlaşdırılmış interaksiyalar

**Go-nun REST üstünlüyü:** encoding/json dəstəyi, error handling + status kodları, middleware (log/auth), HTTPS/TJWT/CSRF qoruma, input sanitizasiya (SQL injection qarşısı).

### 7. İlk veb server (kitabın kodu)
```go
package main

import (
    "fmt"
    "net/http"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello GitforGits BookStore!")
}

func main() {
    http.HandleFunc("/", homeHandler)       // "/" URL → homeHandler
    http.ListenAndServe(":8080", nil)       // 8080 portunda server
}
```
**Sub-kod izahı:**
- `http.HandleFunc("/", homeHandler)` → root URL-i funksiyaya bağlayır
- `http.ListenAndServe(":8080", nil)` → server başladır; `nil` = default ServeMux istifadə et
- `homeHandler(w, r)` → `w` = cavab yazıcısı (client-ə), `r` = sorğu məlumatları
- `fmt.Fprintf(w, ...)` → cavabı client-ə yazır
- İşə salma: `go run main.go` → http://localhost:8080

**Layihə init (CLI):**
```bash
mkdir gitforgits-bookstore
cd gitforgits-bookstore
go mod init gitforgits-bookstore
```

## Əsas terminlər
- Goroutine (qısa axın) — runtime idarəçili yüngül thread
- Channel (kanal) — goroutine-lər arası kommunikasiya/sinxronizasiya
- Static binary (statik binar) — dependency-ləri daxil kompilyasiya
- Go Modules (Go modulları) — v1.11+ dependency idarəsi
- Handler (idarəedici) — ServeHTTP interfeysli obyekt
- ServeMux — multiplexer/router
- REST — Representational State Transfer
- Statelessness (vəziyyətsizlik) — sorğunun özünü kifayət etməsi

## Praktik nəticə
1. Go veb layihəsi həmişə `go mod init` ilə başlayır — go.mod dependency qapısıdır.
2. Minimal server 3 elementdən ibarətdir: handler funksiyası + HandleFunc qeydiyyatı + ListenAndServe.
3. Default mux (`nil`) sadə layihələr üçün kifayətdir; dinamik route-lar üçün növbəti fəsillərdə Gorilla Mux.
4. REST API dizaynında statelessness → hər sorğuda auth məlumatı (JWT) daşıyacaqsan.
5. Deployment üstünlüyü: tək binary = kiçik Docker image = sürətli CI/CD.

## Mənbə
Pages: 19-46 (PDF səh. 19-46)
