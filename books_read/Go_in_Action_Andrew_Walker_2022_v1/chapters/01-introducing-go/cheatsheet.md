# Cheat Sheet — Chapter 1: Introducing Go

## `go run main.go`

**Nə edir:** Go faylını kompilyə edib birbaşa icra edir.

**Sub-komanda/flag izahı:**
- `go` → Go command-line (komanda sətri) aləti
- `run` → Faylı kompilyə edib icra et
- `main.go` → Giriş faylı

**Mənbə:** Chapter 2, page 25

---

## `package main`

**Nə edir:** Go proqramının icra edilə bilən package (paket) olduğunu təyin edir.

**Sub-kod izahı:**
- `package` → Package declaration (paket elanı)
- `main` → Executable program (icra edilən proqram) üçün xüsusi ad

**Mənbə:** Chapter 1, page 6

---

## `import "fmt"`

**Nə edir:** Formatlama və çıxış üçün standart kütubxəni daxil edir.

**Sub-kod izahı:**
- `import` → Kitabxanə daxil etmə
- `"fmt"` → Format (formatlaşdırma) paketi

**Mənbə:** Chapter 1, page 6

---

## `func main()`

**Nə edir:** Program işə düşəndə ilk çağırılan funksiya.

**Sub-kod izahı:**
- `func` → Funksiya elanı
- `main` → Giriş nöqtəsi funksiyasının adı

**Mənbə:** Chapter 1, page 6

---

## `fmt.Println(...)`

**Nə edir:** Konsola mətn və ya dəyər çıxışı verir, avtomatik olaraq yeni sətir əlavə edir.

**Sub-kod izahı:**
- `fmt` → Format paketi
- `Println` → Print line (sətir çap et) — yeni sətir ilə birlikdə çıxış verir

**Mənbə:** Chapter 1, page 6

---

## `go makeMeConcurrent()`

**Nə edir:** Funksiyonu goroutine (gözaltı proses) kimi işə salır, paralel icra olunur.

**Sub-komanda/flag izahı:**
- `go` → Goroutine yaratmaq üçün açar sözcük

**Mənbə:** Chapter 1, pages 6-24

---

## `time.Sleep(1 * time.Second)`

**Nə edir:** Proqramı verilən müddət üçün dayandırır.

**Sub-komanda/flag izahı:**
- `time` → Zaman və tarix əməliyyatları paketi
- `Sleep` → Müddət üçün dayandır
- `1 * time.Second` → 1 saniyə

**Mənbə:** Chapter 1, pages 6-24

---

## `make(chan string)`

**Nə edir:** String (mətn) tipli kanal (channel) yaradır.

**Sub-komanda/flag izahı:**
- `make` → Yeni map, slice və ya channel yaradır
- `chan` → Kanal tipi
- `string` → Kanalda ötürüləcək məlumat tipi

**Mənbə:** Chapter 1, pages 6-24

---

## `var number int`

**Nə edir:** Tam ədəd tipində (int) dəyişən elan edir, sıfır dəyəri təyin edir.

**Sub-kod izahı:**
- `var` → Dəyişən elanı
- `number` → Dəyişənin adı
- `int` → Tam ədəd tipi

**Mənbə:** Chapter 1, pages 6-24

---

## `type Person struct { ... }`

**Nə edir:** Yeni struct (strukt) tipi yaradır, müxtəlif sahələri birləşdirir.

**Sub-kod izahı:**
- `type` → Yeni tip elanı
- `Person` → Tipin adı
- `struct` → Strukt tipləri qrupu
- `FirstName string` → Mətn sahəsi
- `Age int` → Tam ədəd sahəsi

**Mənbə:** Chapter 1, pages 6-24

---

## `func (p Person) Name() string`

**Nə edir:** Person tipinə metod əlavə edir, instance (nüsxə) üzərində işləyir.

**Sub-kod izahı:**
- `func` → Funksiya elanı
- `(p Person)` → Receiver (qəbul edən) — metodun hansı tipə aid olduğu
- `Name()` → Metodun adı
- `string` → Geri qaytardığı tip

**Mənbə:** Chapter 1, pages 6-24

---

## `fmt.Sprintf("%s is %d years old", p.Name(), p.Age)`

**Nə edir:** Formatlı mətn yaradır, dəyərləri stringə çevirir.

**Sub-kod izahı:**
- `%s` → String (mətn) placeholder (yer tutucu)
- `%d` → Integer (tam ədəd) placeholder
- `p.Name()` → Ad sahəsi
- `p.Age` → Yaş sahəsi

**Mənbə:** Chapter 1, pages 6-24

---

## `func SumNumbers[N int | float64](numberSlice ...N) N`

**Nə edir:** Generics (generik) funksiya — int və ya float64 tipləri üçün işləyir.

**Sub-kod izahı:**
- `[N int | float64]` → Type parameter (tip parametri) N, constraint (məhdudiyyət) daxilində
- `numberSlice ...N` → Variadic (dəyişən saylı) parametr
- `N` → Geri qaytarmak tipi

**Mənbə:** Chapter 1, pages 6-24
