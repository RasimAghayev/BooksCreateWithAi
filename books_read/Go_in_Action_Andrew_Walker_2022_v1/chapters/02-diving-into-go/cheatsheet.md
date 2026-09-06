# Cheat Sheet — Chapter 2: Diving Into Go

## `go mod init gobook/wordcount`

**Nə edir:** Yeni Go module yaradır və `go.mod` faylı yaradır.

**Sub-komanda/flag izahı:**
- `go` → Go command-line aləti
- `mod` → Module (modul) əməliyyatları
- `init` → Yeni module initialize (başlat) et
- `gobook/wordcount` → Module adı

**Mənbə:** Chapter 2, page 25

---

## `go build`

**Nə edir:** Go fayllarını kompilyə edir, executable binary yaradır.

**Sub-komanda/flag izahı:**
- `build` → Kompilyasiya və binary yarat

**Mənbə:** Chapter 2, page 25

---

## `go run main.go`

**Nə edir:** Go faylını kompilyə edib birbaşa icra edir, binary saxlamır.

**Sub-komanda/flag izahı:**
- `run` → Kompilyə edib icra et
- `main.go` → Giriş faylı

**Mənbə:** Chapter 2, page 25

---

## `gofmt -d main.go`

**Nə edir:** Faylın formatlaşdırma fərqlərini göstərir (diff).

**Sub-komanda/flag izahı:**
- `-d` → Diff (fərq) rejimi — dəyişiklikləri göstər, faylı dəyişdirmə

**Mənbə:** Chapter 2, page 25

---

## `gofmt -w main.go`

**Nə edir:** Faylı yerində dəyişdirərək standart formata salır.

**Sub-komanda/flag izahı:**
- `-w` → Write (yaz) — faylı diskdə dəyişdir

**Mənbə:** Chapter 2, page 25

---

## `go doc -all strings`

**Nə edir:** Package (paket) dokumentasiyasını konsola çıxarır.

**Sub-komanda/flag izahı:**
- `doc` → Dokumentasiya göstəricisi
- `-all` → Ətraflı məlumat (description, methods, examples)
- `strings` → Package adı

**Mənbə:** Chapter 2, page 25

---

## `go doc -ex bufio.NewScanner`

**Nə edir:** Package-dakı example (nümunə) funksiyaları ilə birlikdə dokumentasiyanı göstərir.

**Sub-komanda/flag izahı:**
- `-ex` → Example-ləri göstər

**Mənbə:** Chapter 2, page 25

---

## `package main`

**Nə edir:** Bu faylın executable (icra edilə bilən) proqram olduğunu təyin edir.

**Sub-kod izahı:**
- `package` → Package elanı
- `main` → Executable proqram üçün xüsusi ad

**Mənbə:** Chapter 2, page 25

---

## `import "fmt"` / `import ("fmt" "strings")`

**Nə edir:** Xarici package-ləri proyektə daxil edir.

**Sub-kod izahı:**
- `import` → Package daxil etmə
- Block syntax (`()`) — bir neçə package eyni anda daxil edilir

**Mənbə:** Chapter 2, page 25

---

## `var numSpaces int`

**Nə edir:** Tam ədəd tipində dəyişən elan edir, sıfır dəyəri (0) təyin edir.

**Sub-kod izahı:**
- `var` → Dəyişən elanı
- Zero value (sıfır dəyəri) — numeric tiplər üçün 0

**Mənbə:** Chapter 2, page 25

---

## `text := "let's count some words!"`

**Nə edir:** Short declaration (qısa elan) operatoru ilə string dəyişəni yaradır.

**Sub-kod izahı:**
- `:=` → Elan və təyinat eyni anda
- Type inference (tip çıxarma) — dəyərdən tip avtomatik müəyyən olunur

**Mənbə:** Chapter 2, page 25

---

## `for i := 0; i < len(text); i++ { }`

**Nə edir:** Klassik 3-hissəli for döngüsü — indeks əsaslı iterasiya.

**Sub-komanda/flag izahı:**
- Init → `i := 0` bir dəfə işləyir
- Condition → `i < len(text)` hər iterasiyadan əvvəl yoxlanır
- Post → `i++` hər iterasiyadan sonra işləyir

**Mənbə:** Chapter 2, page 25

---

## `for _, filename := range os.Args[1:] { }`

**Nə edir:** Range döngüsü ilə `os.Args`-ın bir hissəsini iterasiya edir.

**Sub-kod izahı:**
- `range` → Kolleksiyanın indeksini və dəyərini qaytarır
- `_` → Blank identifier (boş identifikator) — indeksı ignore et
- `os.Args[1:]` → Slice expression (kəsmi ifadə) — 1-ci indeksten sonuna qədər

**Mənbə:** Chapter 2, page 25

---

## `os.ReadFile(filename)`

**Nə edir:** Faylın tam məzmununu yaddaşə oxuyur.

**Sub-komanda/flag izahı:**
- `os` → Əməliyyat sistemi əməliyyatları paketi
- `ReadFile` → Faylı tam oxuyur, `([]byte, error)` qaytarır

**Mənbə:** Chapter 2, page 25

---

## `os.Open(filename)`

**Nə edir:** Faylı açar və `*os.File` pointer qaytarır.

**Sub-komanda/flag izahı:**
- `Open` → Faylı oxumaq üçün açar
- `*os.File` — Fayl deskriptoru, `io.Reader` interfeysini implement edir

**Mənbə:** Chapter 2, page 25

---

## `bufio.NewScanner(file)`

**Nə edir:** Fayl üçün buferli scanner yaradır, sətir/söz token-lərinə ayırır.

**Sub-komanda/flag izahı:**
- `bufio` — Buffered I/O paketi
- `NewScanner` — Scanner konstruktoru, `io.Reader` qəbul edir

**Mənbə:** Chapter 2, page 25

---

## `scanner.Split(bufio.ScanWords)`

**Nə edir:** Scanner-ın tokenization (tokenlərə ayırma) növünü sözlərə dəyişir.

**Sub-komanda/flag izahı:**
- `Split` — SplitFunc təyin edir
- `ScanWords` — Arada boşluqlarla ayrılmış sözləri qaytarır

**Mənbə:** Chapter 2, page 25

---

## `scanner.Scan()`

**Nə edir:** Növbəti token-i oxuyur, `true` qaytarır (EOF-dək və ya xəta yoxdək).

**Sub-komanda/flag izahı:**
- `Scan` — Növbəti tokenə keç

**Mənbə:** Chapter 2, page 25

---

## `scanner.Text()`

**Nə edir:** Son oxunmuş token-i string olaraq qaytarır.

**Sub-komanda/flag izahı:**
- `Text` — Oxunmuş tokeni stringə çevir

**Mənbə:** Chapter 2, page 25

---

## `scanner.Err()`

**Nə edir:** Scan zamanı baş verən ilk non-EOF xətanı qaytarır.

**Sub-komanda/flag izahı:**
- `Err` — Xəta yoxlanışı

**Mənbə:** Chapter 2, page 25

---

## `if err != nil { log.Println(err); os.Exit(1) }`

**Nə edir:** Standart error handling (xəta idarəetməsi) pattern-i.

**Sub-kod izahı:**
- `err != nil` → Xəta varsa
- `log.Println` → Timestamp ilə standard error-a yaz
- `os.Exit(1)` → Proqramı uğursuz exit kodu ilə dayandır

**Mənbə:** Chapter 2, page 25

---

## `strings.Fields(text)`

**Nə edir:** Mətni boşluqlara görə ayırır, sözlər slice-ni qaytarır.

**Sub-komanda/flag izahı:**
- `strings` — Mətn emalı paketi
- `Fields` — Arada boşluq olan yerlərdən ayırır

**Mənbə:** Chapter 2, page 25

---

## `fmt.Println("Found", numSpaces+1, "words")`

**Nə edir:** Arguementləri boşluqla ayıraraq konsola çıxarır, avtomatik yeni sətir.

**Sub-komanda/flag izahı:**
- `Println` — Print line (sətir çap et)

**Mənbə:** Chapter 2, page 25

---

## `file.Close()`

**Nə edir:** Açıq faylı bağlayır, resource (resurs) sərfiyyatını təmin edir.

**Sub-komanda/flag izahı:**
- `Close` — Fayl deskriptorunu bağla

**Mənbə:** Chapter 2, page 25
