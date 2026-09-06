# Chapter 3 — Основы языка Go

## Bu chapter nədən bəhs edir?

Go dilinin sintaksis və əsas tiplərini — məlumat tipləri, dəyişənlər, konstantlar, konteynerlər (massivlər, kəsiklər/slice, xəritələr/map), göstəricilər (pointers), idarəetmə strukturları, xətaların emalı, funksiyalar, strukturlar, metodlar, interfeyslər və konkurensiya (gorutinlər, kanallar) — əhatə edir. Bu chapter praktiki kod nümunələri ilə doludur.

## Əsas fikirlər

### 1. Əsas məlumat tipləri (Basic Data Types)
**Nədir:** Go-da üç kateqoriya var: məntiqi tiplər (bool), rəqəm tipləri (int, float, complex) və simli tiplər (string).

**Necə işləyir:**
- `bool`: `true` və ya `false` (1 bit)
- `int8/16/32/64`, `uint8/16/32/64`, `float32/64`, `complex64/128`
- `byte` = `uint8` mnemonik alias
- `rune` = `int32` mnemonik alias (Unicode kod nöqtəsi)
- `string`: dəyişməz (immutable) UTF-8 byte kəsik

**Kitabdan kod nümunəsi:**
```go
var b bool = true
var i int = 42
var s string = "Hello\nworld!"
var r rune = 'A'  // int32, Unicode kodu
```

### 2. Dəyişənlər (Variables)
**Nədir:** `var` açar sözü ilə tip və ya `:=` ilə avtomatik tip müəyyən edərək dəyişən təyin edir.

**Necə işləyir:**
- `var foo int = 42` — aydın tip ilə
- `foo := 42` — qısaltılmış sintaksis, avtomatik tip
- `var foo, bar int = 42, 1302` — birdən çox dəyişən
- İstifadə edilməyən dəyişən compile error vətir — Go "mess"ı sevmir

**Kitabdan kod nümunəsi:**
```go
var a, b, c = true, 2.3, "four"
x, y := 0, 2
```

### 3. Konstantlar (Constants)
**Nədir:** `const` ilə təyin edilən, dəyişməz dəyərlər. Nüve (nil) dəyəri yoxdur, təyin edilir.

**Necə işləyir:** Kompilasiya vaxtında qiymətləndirilir. `const text = "Does %s rule? %t!"` — tip avtomatik təyin edilir.

### 4. Konteynerlər: Massivlər, Kəsiklər, Xəritələr (Arrays, Slices, Maps)
**Nədir:** Topluluq (collection) məlumatları saxlamaq üçün üç struktur.

**Arrays:**
- Sabit uzunluqlu, təyinarı element tipi
- `[3]int{2, 4, 6}` literal
- `len(a)` — uzunluq

**Slices:**
- Dəyişən ölçülü massiv abstraksiyası
- İçində 3 komponent var: pointer, uzunluq (len), tutum (cap)
- `make([]int, 3)` — yaratmaq
- `append(s, 2)` — element əlavə etmək
- `s[i:j]` — kəsik operatoru

**Kitabdan kod nümunəsi:**
```go
s0 := []int{0, 1, 2, 3, 4, 5, 6}
s1 := s0[:4]  // [0 1 2 3]
s2 := s0[3:]  // [3 4 5 6]
s0[3] = 42    // s1 və s2 də dəyişir (bəşkəli backing array)
```

**Maps:**
- Hash table əsasında, key → value map
- `map[string]float32{"a": 1.0, "b": 2.0}`
- `delete(m, "key")` — element silmək
- `val, ok := m["key"]` — mövcudluq yoxlamaq

### 5. Göstəricilər (Pointers)
**Nədir:** Dəyişənin yaddaş ünvanını (address) saxlayan tip. `&` ilə ünvan alınır, `*` ilə dəyər oxunur/deyisdirilir.

**Necə işləyir:**
```go
var a int = 10
var p *int = &a  // p → a-nın ünvanını saxlayır
fmt.Println(*p)  // 10 (dereference)
*p = 20          // a-nı dəyişir
fmt.Println(a)   // 20
```

**Nəyə lazımdır:** Funksiyaya parametr kimi göndərməkdə dəyərin kopyalanması əvəzinə, orijinal dəyəri dəyişmək.

### 6. İdarəetmə strukturları (Control Structures)
**Nədir:** `for`, `if`, `switch` — C-tərkibli dillərdən tanış, lakin sadələşdirilmiş.

**For döngüsü:**
- Yeganə döngü tipi (`while` və `do-while` yoxdur)
- `for i := 0; i < 10; i++ { ... }` — klassik
- `for i < 10 { ... }` — while əvəzi
- `for { ... }` — sonsuz döngü
- `for i, v := range s { ... }` — kolleksiya üzərində iterasiya

**If:** `if _, err := os.Open("foo"); err != nil { ... }` — ilkiniz operatoru daxildir.

**Switch:** Fallthrough avtomatik deyil, `fallthrough` açar sözü ilə aktivləşdirilir. Expression olmadan `switch true` kimi işləyir.

### 7. Xəta emalı (Error Handling)
**Nədir:** Go-da exception (istisna) mexanizmi yoxdur. Xətalar `error` interface vasitəsilə qaytarılır.

**Necə işləyir:** `error` interface yalnız bir metod tələb edir: `Error() string`.

**Kitabdan kod nümunəsi:**
```go
file, err := os.Open("somefile.ext")
if err != nil {
    log.Fatal(err)
    return err
}
```

**Xəta yaratmaq:**
- `errors.New("message")` — sadə xəta
- `fmt.Errorf("format %s", val)` — formatlı xəta

### 8. Funksiyalar (Functions)
**Nədir:** Go funksiyaları birinci sinif dəyərlər (first-class citizens). Birdən çox dəyər qaytara, closure və variadic parametrlər qəbul edə bilər.

**Necə işləyir:**
- `func add(x, y int) int` — tip parametrdən sonra
- `func swap(x, y string) (string, string)` — çoxlu return
- `defer` — funksiyadan çıxışda icra olunur (resource cleanup üçün)
- Closure: daxili funksiya valideynin dəyişənlərinə giriş saxlayır

**Kitabdan kod nümunəsi:**
```go
func incrementer() func() int {
    i := 0
    return func() int {
        i++
        return i
    }
}

inc := incrementer()
fmt.Println(inc()) // 1
fmt.Println(inc()) // 2
```

### 9. Strukturlar, Metodlar və İnterfeyslər (Structs, Methods, Interfaces)
**Strukturlar:**
- Sadə sahələr toplusu: `type Vertex struct { X, Y float64 }`
- Nüve dəyəri `nil` ola bilməz, bütün sahələrin sıfır dəyəri var

**Metodlar:**
- Funksiyalardır, tipə bağlıdır: `func (v Vertex) Area() float64`
- Pointer alıcı: `func (v *Vertex) Scale(f float64)` — orijinal dəyəri dəyişir

**İnterfeyslər:**
- Davranış müqaviləsi, metod siqnaturası toplusu
- Duck typing: "quacks like a duck" — metodları varsa, implement edir
- `io.Reader`, `io.Writer` — ən çox istifadə olunan kiçik interfeyslər

**Kitabdan kod nümunəsi:**
```go
type Shape interface {
    Area() float64
}

type Rectangle struct{ Width, Height float64 }
func (r Rectangle) Area() float64 { return r.Width * r.Height }

type Circle struct{ Radius float64 }
func (c Circle) Area() float64 { return math.Pi * c.Radius * c.Radius }

func PrintArea(s Shape) {
    fmt.Printf("%T's area is %.2f\n", s, s.Area())
}
```

**Kompozisiya (Embedding):**
- `type ReadWriter struct { *Reader; *Writer }` — metodlar və sahələr avtomatik yüksəldilir (promotion)

### 10. Konkurensiya: Gorutinlər və Kanallar (Concurrency)
**Gorutinlər (Goroutines):**
- Go funksiyasını `go` açar sözü ilə asinxron işə salmaq
- Çox yüngül, OS thread-lərdən daha az yaddaş istehlak edir
- Runtime avtomatik schedulə edir

**Kanallar (Channels):**
- Gorutinlər arasında tipləşdirilmiş mesajlaşma
- `ch := make(chan int)` — yaratmaq
- `ch <- val` — göndərmək
- `val := <-ch` — qəbul etmək
- Unbuffered channel: göndərən və alan gorutin bloklanır, hand-off sinxronlaşdırılır

**Kitabdan kod nümunəsi:**
```go
ch := make(chan string)

go func() {
    message := <-ch        // bloklanır, gözləyir
    fmt.Println(message)   // "ping"
    ch <- "pong"           // bloklanır, gözləyir
}()

ch <- "ping"               // bloklanır, gözləyir
fmt.Println(<-ch)          // "pong"
```

**Buffered kanallar:**
- `ch := make(chan string, 2)` — 2 elementlik buffer
- Buffer dolana qədər göndərmə bloklanmaz

**Select:**
- Bir neçə channel üzərində multiplexing
- Default case ilə non-blocking yoxlanış
- Təsadüfi seçim (random) eyni anda hazır case-lərdə

**Kitabdan kod nümunəsi:**
```go
select {
case <-ch1:
    fmt.Println("Got something")
case x := <-ch2:
    fmt.Println(x)
case ch3 <- y:
    fmt.Println(y)
default:
    fmt.Println("None of the above")
}
```

## Əsas terminlər
- Goroutine (qorutin) — yüngül asinxron funksiya
- Channel (kanal) — tipləşdirilmiş mesajlaşma kanalı
- Select — multiplex kanal operatoru
- Slice — dəyişən ölçülü massiv abstraksiyası
- Map — hash table, key/value map
- Pointer (göstərici) — yaddaş ünvanı
- Defer — gecikdirilmiş funksiya çağırışı
- Closure (bağlanma) — daxili funksiya, valideyn scope saxlayır
- Interface (interfeys) — metod müqaviləsi, duck typing
- Error interface — xəta müqaviləsi (Error() string)
- Variadic function (`...interface{}`)
- Fallthrough — switch case-ləri keçmək

## Praktik nəticə
Go sintaksisi sadə, amma güclü: gorutinlər və kanallarla konkurensiya daxili, error interface ilə xəta emalı aydın, interfeyslər ilə kod kompozisiya asan. Bu chapter növbəti chapter-larda (4-11) istifadə ediləcək bütün əsas konstruksiyaları örtür.

## Mənbə
Pages: 48-87 (PDF səh. 48-87)
