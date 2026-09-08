# Chapter 12 — Modules and Documentation (səh. 252-259)

## Bu fəsil nədən bəhs edir?

Go modules — asılılıqların idarəolunması (`go.mod`, `go mod init`,
avtomatik versiya çözümü) və Go-nun minimalizminə söykənən kod
sənədləşdirməsi (comment konvensiyaları, `go doc`, `godoc` server,
işlədilə bilən nümunələr).

## Əsas fikirlər

### 1. Go modules — go.mod
**Nədir:** Asılılıq idarəetməsinin rəsmi standartı; layihə tələblərini
`go.mod` faylında versiya ilə birlikdə saxlayır (`go get`-in yerini tutur —
versioning problemi həll olunur).

**Kitabdan kod nümunəsi (üçüncü tərəf logger):**
```go
package main

import "github.com/rs/zerolog/log"

func main() {
    log.Info().Msg("Save the world with Go!!!")
}
```

**Module initialization:**
```bash
>>> go mod init
# Yaranan go.mod:
module github.com/juanmanuel-tirado/SaveTheWorldWithGo/11_modules/modules/example_01
go 1.15
```

**Avtomatik asılılıq çözümü:**
```bash
>>> go build main.go
go: finding module for package github.com/rs/zerolog/log
go: found github.com/rs/zerolog/log in github.com/rs/zerolog v1.20.0
```

**Nəticə go.mod:**
```
module github.com/juanmanuel-tirado/SaveTheWorldWithGo/11_modules/modules/example_01

go 1.15
require github.com/rs/zerolog v1.20.0 // indirect
```

**Sub-izahlar:**
- `go mod init` → modul adı + Go versiyası ilə go.mod yaradır (yalnız kök
  qovluqda)
- `go build/run/test` → asılılıqları avtomatik analiz edib go.mod-a yazır
- `require paket versiya` → asılılıq sətri (sətirə bir paket)
- `// indirect` → dolaylı asılılıq — zerolog-un öz go.mod-u var, onun
  ehtiyacları üçün əlavə paketlər endirilib
- Paketlər adi halda GOPATH-a düşür

**Vendor qovluğu (tövsiyə olunmur):**
```bash
>>> go mod vendor   # bütün paketləri vendor/ qovluğünə kopyalayır
```
- Yalnız xüsusi ehtiyac varsa; repo-nu paylaşarkan vendor İSTİFADƏ ETMƏYİN —
  modul sisteminin bütün üstünlükləri məhv olur, kod dublikası yaranır
- Köhnə alətlər (go dep) ilə qarışıq layihələrdə uyğunluğu yoxlayın

### 2. Sənədləşdirmə konvensiyaları
**Qaydalar:**
1. `package` elanından əvvəlki şərh → paket şərhi
2. Hər paketin şərhi olmalıdır (çoxfayllı paketdə bir faylda kifayətdir)
3. Hər **exported** adın (böyük hərflə başlayan) şərhi olmalıdır
4. Şərh **təsvir olunan elementin adı ilə başlayır**

**Kitabdan kod nümunəsi:**
```go
// Package example_01 contains a documentation example.
package example_01

import "fmt"

// Msg represents a message.
type Msg struct {
    // Note is the note to be sent with the message.
    Note string
}

// Send a message to a target destination.
func (m *Msg) Send(target string) {
    fmt.Printf("Send %s to %s\n", m.Note, target)
}

// Receive a message from a certain origin.
func (m *Msg) Receive(origin string) {
    fmt.Printf("Received %s from %s\n", m.Note, origin)
}
```

### 3. go doc aləti
```bash
>>> go doc -all          # cari layihənin tam sənədi
>>> go doc fmt           # istənilən GOPATH paketinin sənədi
>>> go doc json.decode   # Decode metodunun sənədi
```

**Çıxış nümunəsi:**
```
package example_01 // import "github.com/.../example_01"
Package example_01 contains a documentation example.

TYPES
type Msg struct {
    // Note is the note to be sent with the message.
    Note string
}
    Msg represents a message.
func (m *Msg) Receive(origin string)
    Receive a message from a certain origin.
func (m *Msg) Send(target string)
    Send a message to a target destination.
```

### 4. godoc server + işlədilə bilən nümunələr
```bash
>>> godoc -http=:8080          # localhost:8080-da sənəd serveri
>>> godoc -http=:8080 -play    # nümunələr runtime box-da İŞLƏDİLƏ BİLƏR
```

**Nümunələr `_test` paketində yerləşir (Chapter 11.2 qaydaları):**
```go
package example_01_test

import "github.com/juanmanuel-tirado/savetheworldwithgo/11_modules/godoc/example_01"

func ExampleMsg_Send() {
    m := example_01.Msg{"Hello"}
    m.Send("John")
    // Output:
    // Send Hello to John
}

func ExampleMsg_Receive() {
    m := example_01.Msg{"Hello"}
    m.Receive("John")
    // Output:
    // Received Hello from John
}
```
- `ExampleTip_Metod` adlandırması → godoc-da konkret metodun yanında göstərilir
- `-play` rejimində istifadəçi nümunəni brauzerdə icra edə bilər

## Əsas terminlər
- Module (modul) — asılılıqları təsvir edən vahid
- go.mod — modul tərif faylı (module, go, require)
- Indirect dependency — dolaylı asılılıq
- Vendor — asılılıqların layihə daxilinə kopyalanması (deprecated)
- go doc — terminal sənəd aləti
- godoc — web sənəd serveri (-play ilə interaktiv)
- Exported name — başqa paketlərdən görünən ad (böyük hərf)

## Praktik nəticə
Yeni layihədə ilk əmr `go mod init` olmalıdır; asılılıqlar build/run/test
zamanı avtomatik həll olunub go.mod-a yazılır — `go get`-i yalnız əl ilə
əlavə etmək istəyəndə işlədin. Sənədləşdirmə minimal formaldır: hər
exported ad + paket üçün "ad ilə başlayan" şərh. `go doc` tez baxış üçün,
`godoc -http=:8080 -play` isə komandaya sənəd göstərmək üçün idealdır.

## Mənbə
Pages: 252-259 (PDF 252-259)
