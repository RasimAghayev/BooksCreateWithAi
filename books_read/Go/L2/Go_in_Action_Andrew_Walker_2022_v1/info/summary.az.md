# Xülasə — Go in Action, Second Edition

## Kitab Haqqında

**Go in Action, Second Edition** — Andrew Walker və William Kennedy tərəfindən yazılmış, Manning nəşriyyatından çıxmış praktiki Go proqramlaşdırma kitabı. Bu MEAP V03 versiyasındadır və Chapter 1-4 ehtiva edir.

## Əsas Mövzular

### Chapter 1 — Introducing Go
Go-nun nəyə görə yaradıldığını, onun əsas xüsusiyyətlərini izah edir: sürətli kompilyasiya, tip təhlükəsizliyi, daxili paralellik (goroutine + channel), tip sistemi (struct, interface, generics) və zibil toplayıcı (garbage collection).

### Chapter 2 — Diving Into Go
Praktik tətbiq hazırlayaraq Go ecosystem-ünə giriş edir: ilk Go proqramını (word counter) yazmaq, `go build`, `go run`, `gofmt`, `go mod`, error handling, `bufio.Scanner`, `io.Reader` interfeysi və `range` döngüsü.

### Chapter 3 — Primitive Types And Operators
Go-nun əsas tiplərini əhatə edir: integer, floating-point, complex, bool, struct, pointer, string, rune və üzəri şriftli (bitwise) operatorlar. Hər tipin praktik istifadəsi, type inference, zero value və tip seçimi məsləhətləri verilir.

### Chapter 4 — Collection Types
Arrays, slices və maps tiplərini əhatə edir. Slice-lərin daxili işi (pointer, length, capacity), `append` davranışı, slice expressions, map iterasiya və filterləmə, nil/empty slice fərqi.

## Ən Vacib 5 Fikir

1. **Go compiled və statically-typed (kompilyə edilən və statik tip) dildir** — kompilyasiya çox sürətlidir, tip xətaları runtime-də deyil, compile vaxtında bərpa olunur.

2. **Concurrency (paralellik) daxili konsepsiyadır** — Goroutine-lər çox yüngüldür, channel-lər vasitəsilə təhlükəsiz data mübadiləsi. Standart kitabxanada (məsələn, `net/http`) concurrency artıq daxildir.

3. **Error handling explicit (açıq)dir** — `if err != nil` pattern-i istisnasız tətbiq olunur. Exception mexanizmi yoxdur.

4. **Slices praktiki arrays-dir** — Dinamik, reference type, `append` ilə genişlənir. Arrays əsas tipdir, lakin slice daha çox işlənir.

5. **Interface-lər implicit (daxili) implementasiya tələb edir** — `Describer` interfeysi üçün `Describe()` metodu yazmaq kifayətdir, `implements Describer` deyilmir.

## Kitabın ən Dəyərli Hissəsi

Chapter 2, pages 25-66 — çünki burada ilk real Go tətbiqi (word counter) adım-adım qurulur. `go run`, `go build`, `gofmt`, error handling, `bufio.Scanner` və `io.Reader` konsepsiyaları praktik nümunələrlə öyrədilir. Bu chapter oxumadan sonra developer real Go layihəsi qurmağa başlaya bilir.

## Qeydlər

- Bu MEAP V03 versiyası Chapter 1-4 ehtiva edir.
- Kitabın tam nəşri daha çox chapter (5-15) və appendices (əlavələr) ehtiva edir.
- Source code: https://github.com/flowchartsman/go-in-action
