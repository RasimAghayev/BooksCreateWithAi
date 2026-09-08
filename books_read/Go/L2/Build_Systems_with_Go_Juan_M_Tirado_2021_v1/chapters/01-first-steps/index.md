# Chapter 1 — First steps with Go (səh. 10-20)

## Bu fəsil nədən bəhs edir?

İlk Go proqramının yazılışı, kompilyasiyası və icrası: `go build` ilə
kompilyasiya, `os.Args` ilə proqrama arqument ötürülməsi, `strconv.Atoi`
ilə string→int çevrilməsi və error idarəolunması. Fəsil praktiki "Save the
world with Go!!!" nümunəsi ətrafında qurulub — kompyuterə iki rəqəm ötürüb
cümlərini verən mini-kalkulyator qurulur.

## Əsas fikirlər

### 1. İlk Go proqramı
**Nədir:** Konsola mesaj çap edən ən sadə Go proqramı (Hello World analoqu).

**Kitabdan kod nümunəsi:**
```go
package main

import "fmt"

func main() {
    fmt.Println("Save the world with Go!!!")
}
```

**Sub-kod izahı:**
- `package main` → kodun daxil olduğu paketin adı; icra olunan proqramlar
  həmişə `main` paketində olur
- `import "fmt"` → `Println` funksiyası üçün lazım olan standart kitabxana
- `func main()` → proqramın giriş nöqtəsi; məntiq bu funksiyanın
  mötərizələri daxilindədir
- `fmt.Println(...)` → mesajı standart çıxışa (stdout) yazır

### 2. Kompilyasiya və icra
**Necə işləyir:** Go kompilyasiya olunan dildir — kod icra olunmazdan əvvəl
`go build` ilə platformaya uyğun icra olunan fayla çevrilir.

```bash
>> go build main.go
>> ls
main main.go
>> ./main
Save the world with Go!!!
```

**Sub-komanda izahı:**
- `go build main.go` → `main.go` faylını kompilyasiya edib `main` adlı
  icra olunan fayl yaradır
- `./main` → `./` prefiksi ilə cari qovluqdakı icra olunan faylı işə salır
- `ls` → nəticə faylının (main) yarandığını yoxlamaq üçün

### 3. Proqrama arqument ötürülməsi (os.Args)
**Nədir:** Command-line argumentlərinə (komanda sətri arqumentlərinə)
proqram daxilindən çıxış yolu.

**Necə işləyir:** `os.Args` — proqramın adı daxil olmaqla bütün arqumentləri
saxlayan massiv (array). Slicing ilə proqram adı çıxardıla bilər.

**Kitabdan kod nümunəsi:**
```go
argsWithProg := os.Args
argsWithoutProg := os.Args[1:]
arg := os.Args[3]
```

**Sub-kod izahı:**
- `os.Args` → bütün arqumentlər + icra olunan faylın adı (index 0)
- `os.Args[1:]` → proqram adı çıxarılmış arqumentlər
- `os.Args[3]` → 3-cü pozisiyadakı arqument (indexing 0-dan başlayır)

### 4. String→int çevrilməsi (strconv.Atoi) və error idarəetməsi
**Nəyə lazımdır:** Arqumentlər həmişə string kimi gəlir; riyazi əməliyyat
üçün integerə çevrilməlidir. `strconv.Atoi` çevirmənin uğursuz olma
ehtimalına görə iki dəyər qaytarır.

**Çevirmə nümunələri:**
- `"42"` → `42` (uğurlu)
- `"-33"` → `-33` (uğurlu)
- `"4.2"` → xəta — float gözlənilmir
- `"thirteen"` → xəta — mətn təsviri rəqəmə çevrilmir

**Kitabdan kod nümunəsi (mini-kalkulyator):**
```go
package main

import (
    "fmt"
    "os"
    "strconv"
)

func main() {
    argsWithProg := os.Args
    numA, err := strconv.Atoi(argsWithProg[1])
    if err != nil {
        fmt.Println(err)
        os.Exit(2)
    }
    numB, err := strconv.Atoi(argsWithProg[2])
    if err != nil {
        fmt.Println(err)
        os.Exit(2)
    }
    result := numA + numB
    fmt.Printf("%d + %d = %d\n", numA, numB, result)
}
```

**Sub-kod izahı:**
- `strconv.Atoi(...)` → string-i integerə çevirir; `(dəyər, error)` cütü
  qaytarır
- `if err != nil` → error boş deyilsə çevirmə alınmadı deməkdir
- `os.Exit(2)` → proqramı xəta kodu (2) ilə dayandırır
- `fmt.Printf("%d + %d = %d\n", ...)` → formatlı çıxış: `%d` yerlərinə
  numA, numB, result dəyərləri yerləşir

**İcra nəticəsi:**
```bash
>>> ./sum 2 2
2 + 2 = 4
>>> ./sum 42 -2
42 + -2 = 40
>>> ./sum 2 two
strconv.Atoi: parsing "two": invalid syntax
```

## Əsas terminlər
- Compilation (kompilyasiya) — mənbə kodun icra olunan fayla çevrilməsi
- CLI arguments (komanda sətri arqumentləri) — proqram işə düşəndə ötürülən dəyərlər
- Error handling (xətaların idarəolunması) — `err != nil` yoxlaması ilə
- Exit code (çıxış kodu) — proqramın bitmə statusu (0 = uğur, 2 = xəta)

## Praktik nəticə
Go proqramı yazmaq üçün minimum: `main` paketi + `main()` funksiyası.
CLI alətləri üçün `os.Args` + `strconv.Atoi` kombinasiyası əsasdır —
hər çevirmə `err` yoxlaması ilə müşayiət olunmalıdır. Xəta halında
`os.Exit(2)` qeyri-sıfır kodla çıxış standart praktikadır.

## Mənbə
Pages: 10-20 (PDF 10-20)
