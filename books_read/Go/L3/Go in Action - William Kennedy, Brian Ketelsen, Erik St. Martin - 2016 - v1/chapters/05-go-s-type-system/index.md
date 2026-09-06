# Chapter 5 — Go's type system

## Bu chapter nədən bəhs edir?

Go-nun statik tip sisteminin dərinləşməsi: user-defined tiplər (struct və mövcud tip əsasında), metodlar (value/pointer receiver), tiplərin təbiəti (primitive/nonprimitive), interfeyslərin implementasiya mexanikası (iTable, method sets), polimorfizm, type embedding (kompozisiya) və identifier-lərin export/unexport qaydaları.

## Əsas fikirlər

### 1. User-defined tiplər
**Nədir:** Öz tiplərinizi elan etmək üçün 2 yol: `struct` (kompozit tip) və ya mövcud tipi əsas götürərək yeni tip.

**Struct elanı:**
```go
// user defines a user in the program.
type user struct {
    name      string
    email     string
    ext       int
    privileged bool
}
```

**Struct literal formaları:**
```go
// Zero value üçün — var istifadə et (idiomatik)
var bill user

// Sahə adları ilə — sıra əhəmiyyətsiz
lisa := user{
    name:      "Lisa",
    email:     "lisa@email.com",
    ext:       123,
    privileged: true,
}

// Sahə adsız — sıra mütləq struct elanı ilə eyni olmalıdır
lisa := user{"Lisa", "lisa@email.com", 123, true}
```

**Nested struct-lar:**
```go
type admin struct {
    person user
    level  string
}
fred := admin{
    person: user{
        name: "Lisa", email: "lisa@email.com", ext: 123, privileged: true,
    },
    level: "super",
}
```

**Mövcud tip əsasında yeni tip:**
```go
type Duration int64
```

**Sub-kod izahı:** `Duration`-ın base type-ı `int64`-dir, amma bunlar **iki fərqli tipdir** — kompilyator implicit conversion etmir:
```go
var dur Duration
dur = int64(1000)   // COMPILE XƏTASI:
// cannot use int64(1000) (type int64) as type Duration in assignment
```

**Qayda:** zero value üçün `var`, ilkin dəyər varsa `:=` + struct literal.

### 2. Metodlar — value vs pointer receiver
**Nədir:** Metod = funksiya + `func` açar sözü ilə adı arasında elan olunan **receiver** parametri. Receiver funksiyanı tipə bağlayır.

**Kitabdan kod nümunəsi:**
```go
type user struct {
    name  string
    email string
}

// notify implements a method with a value receiver.
func (u user) notify() {
    fmt.Printf("Sending User Email To %s<%s>\n", u.name, u.email)
}

// changeEmail implements a method with a pointer receiver.
func (u *user) changeEmail(email string) {
    u.email = email
}

func main() {
    bill := user{"Bill", "bill@email.com"}
    bill.notify()                    // value → value receiver OK
    bill.changeEmail("bill@newdomain.com")  // value → pointer receiver OK (Go &bill edir)

    lisa := &user{"Lisa", "lisa@email.com"}
    lisa.notify()                    // pointer → value receiver OK (Go *lisa edir)
    lisa.changeEmail("lisa@comcast.com")    // pointer → pointer receiver OK
}
```

**Sub-kod izahı:**
- `func (u user) notify()` → **value receiver**: metod dəyərin **kopyası** üzərində işləyir
- `func (u *user) changeEmail(email)` → **pointer receiver**: metod real dəyərlə paylaşılır — dəyişiklik çağırıcıda görünür
- `bill.changeEmail(...)` → Go arxa planda `(&bill).changeEmail(...)` edir — convenience
- `lisa.notify()` → Go arxa planda `(*lisa).notify()` edir

### 3. Tiplərin təbiəti (nature of types)
**Əsas sual:** "Bu tipdən nəsə əlavə/çıxarmaq yeni dəyər yaradır, yoxsa mövcudu dəyişir?"
- Yeni dəyər → **value receiver** (primitive nature)
- Dəyişmək → **pointer receiver** (nonprimitive nature)

**Built-in tiplər (rəqəm, string, bool):** həmişə primitive — kopya ötür. Standart kitabxana nümunəsi (`strings.Trim`):
```go
func Trim(s string, cutset string) string {
    if s == "" || cutset == "" {
        return s
    }
    return TrimFunc(s, makeCutsetFunc(cutset))
}
```
Kopya qəbul edir, **yeni string** qaytarır.

**Reference tiplər (slice, map, channel, interface, function):** header value-dirlər (alt struktura pointer + idarə sahələri). Header **kopyalanmaq üçün tasarlanıb** — kopya alt struktur ilə oxşarı paylaşır. Nümunə — `net.IP`:
```go
type IP []byte

func (ip IP) MarshalText() ([]byte, error) { ... }   // value receiver — normal!
```

**Struct tiplər:** hər ikisi ola bilər.
- **Primitive nümunə** — `time.Time`: zaman nöqtəsi dəyişməzdir. `Add` metodu value receiver ilə **yeni** Time qaytarır:
```go
func (t Time) Add(d Duration) Time {
    t.sec += int64(d / 1e9)
    // ...
    return t
}
```
- **Nonprimitive nümunə** — `os.File`: kopyalanması təhlükəlidir. `Open` **pointer** qaytarır (factory pointer qaytarırsa — nonprimitive siqnalıdır). `Chdir` heç nə dəyişməsə belə pointer receiver istifadə edir — dəyişməyən halda da paylaşım qaydası saxlanılır:
```go
func (f *File) Chdir() error { ... }
```

**Qayda:** Receiver seçimi "metod nə edir"-ə deyil, "tipin təbiəti nədir"-ə əsaslanır. Yeganə istisna — interfeyslərlə işləyərkən value receiver lazım ola bilər (method sets-ə görə).

### 4. Interfeyslər — standart kitabxana gücü
**Nədir:** Davranış elan edən tiplər; implementasiya user-defined tiplər tərəfindən metodlarla edilir.

**Kitabdan kod nümunəsi — sadə curl (io.Reader/io.Writer):**
```go
package main

import (
    "fmt"
    "io"
    "net/http"
    "os"
)

func main() {
    r, err := http.Get(os.Args[1])
    if err != nil {
        fmt.Println(err)
        return
    }

    // Copies from the Body to Stdout.
    io.Copy(os.Stdout, r.Body)
    if err := r.Body.Close(); err != nil {
        fmt.Println(err)
    }
}
```

**Sub-kod izahı:**
- `http.Get(...)` → cavabın `Body` sahəsi `io.ReadCloser` (io.Reader + io.Closer) interfeysidir
- `io.Copy(os.Stdout, r.Body)` → 2 interfeys: dest `io.Writer` (Stdout), source `io.Reader` (Body) — web server → terminal streaming
- `io.Copy` hər hansı Reader/Writer cütü ilə işləyir — bu, polimorfizmin gücü

**bytes.Buffer nümunəsi:**
```go
var b bytes.Buffer
b.Write([]byte("Hello"))
fmt.Fprintf(&b, "World!")
io.Copy(os.Stdout, &b)   // Buffer həm Reader, həm Writer-dir
```

### 5. İnterfeys implementasiyası — internals
**Necə işləyir:** İnterfeys dəyəri **2 sözlük** strukturdur:
1. **iTable-a pointer** — saxlanan dəyərin tip məlumatı + metod siyahısı
2. **Saxlanan dəyərə pointer**

Konkret tip dəyəri mənimsədiləndə (`n = user{"Bill"}`) — iTable `Type(user)`, dəyər özü saxlanılır. Pointer mənimsədiləndə (`n = &user{...}`) — iTable `Type(*user)`, ünvan saxlanılır. İnterfeysdən metod çağırışında iTable-dakı konkret metod icra olunur → **polimorfizm**.

### 6. Method sets (metod dəstləri) — QANUN
**Nədir:** Hansı tiplərin (value/pointer) interfeysi təmin etdiyini müəyyən edən qaydalar.

**Spesifikasiya (receiver perspektivi ilə — əsas cədvəl):**
```
Metod Receiver    | İnterfeysi implement edənlər
(t T)             | T və *T   (hər ikisi)
(t *T)            | yalnız *T (yalnız pointer)
```

**Kitabdan kod nümunəsi (compile xətası):**
```go
type notifier interface {
    notify()
}

type user struct {
    name  string
    email string
}

// notify implements a method with a pointer receiver.
func (u *user) notify() {
    fmt.Printf("Sending user email to %s<%s>\n", u.name, u.email)
}

func main() {
    u := user{"Bill", "bill@email.com"}
    sendNotification(u)      // XƏTA!
    // cannot use u (type user) as type notifier in argument to
    // sendNotification:
    //   user does not implement notifier
    //   (notify method has pointer receiver)
}

func sendNotification(n notifier) {
    n.notify()
}
```

**Həll:** `sendNotification(&u)` — pointer ötür.

**Niyə bu qadağa?** Hər dəyərin ünvanını almaq mümkün deyil:
```go
type duration int

func (d *duration) pretty() string {
    return fmt.Sprintf("Duration: %d", *d)
}

duration(42).pretty()
// cannot call pointer method on duration(42)
// cannot take the address of duration(42)
```
Value-nun method set-ində yalnız value receiver metodları olur; pointer-ın method set-i hər ikisini əhatə edir.

### 7. Polimorfizm nümunəsi
**Kitabdan kod nümunəsi:**
```go
// user və admin hər ikisi notifier implement edir
func (u *user) notify() {
    fmt.Printf("Sending user email to %s<%s>\n", u.name, u.email)
}
func (a *admin) notify() {
    fmt.Printf("Sending admin email to %s<%s>\n", a.name, a.email)
}

bill := user{"Bill", "bill@email.com"}
sendNotification(&bill)   // "Sending user email..."
lisa := admin{"Lisa", "lisa@email.com"}
sendNotification(&lisa)   // "Sending admin email..."
```
Eyni `sendNotification` funksiyası saxladığı konkret tipə uyğun fərqli davranış icra edir.

### 8. Type embedding (tip daxiletməsi)
**Nədir:** Mövcud tipi yeni struct-ın içinə sahə adı YAZMADAN daxil etmək — kompozisiya yolu ilə genişləndirmə.

**Kitabdan kod nümunəsi:**
```go
type user struct {
    name  string
    email string
}

func (u *user) notify() {
    fmt.Printf("Sending user email to %s<%s>\n", u.name, u.email)
}

type admin struct {
    user          // Embedded Type — ad yazılmır!
    level string
}

ad := admin{
    user:  user{name: "john smith", email: "john@yahoo.com"},
    level: "super",
}

// İç tip birbaşa əlçatandır — öz kimliyini itirmir:
ad.user.notify()

// Inner type promotion — iç tipin identifikatorları xarici tipə qalxır:
ad.notify()      // eyni nəticə!
```

**Sub-kod izahı:**
- `user` sahə adı YOXDUR — yalnız tip adı yazılır → embedding
- **Promotion:** iç tipin sahə/metodları xarici tipə aid olur
- İç tip həmişə öz kimliyi ilə mövcuddur: `ad.user.notify()` hələ də işləyir

**İnterfeys + embedding:** iç tipin interfeys implementasiyası da promote olunur:
```go
sendNotification(&ad)   // admin özü notify yazmasa belə — user-in implementasiyası işləyir
```

**Override (xarici tip eyni metodu yazsa):**
```go
func (a *admin) notify() {
    fmt.Printf("Sending admin email to %s<%s>\n", a.name, a.email)
}
```
- `sendNotification(&ad)` → **admin-in** versiyası işləyir
- `ad.user.notify()` → iç tipin versiyası hələ də birbaşa çağırıla bilir
- Promotion yalnız xarici tip öz implementasiyasını yazmadıqda işləyir

### 9. Export / Unexport
**Qayda:** İdentifikator böyük hərflə başlayırsa **exported** (paketdən kənarda görünür), kiçik hərflə başlayırsa **unexported**.

**Unexported tipə paketdən kənardan çıxış YOXDUR:**
```go
// counters paketi:
package counters
type alertCounter int   // unexported

// başqa paketdən:
counter := counters.alertCounter(10)
// XƏTA: cannot refer to unexported name counters.alertCounter
```

**Factory function pattern — `New` konvensiyası:**
```go
// counters paketi:
package counters

type alertCounter int   // unexported qalır

// New creates and returns values of the unexported type.
func New(value int) alertCounter {
    return alertCounter(value)
}

// başqa paketdən:
counter := counters.New(10)   // İŞLƏYİR!
```

**Niyə işləyir?**
1. Export/unexport **identifikatorlara** aiddir, **dəyərlərə** yox — funksiya unexported tipin dəyərini qaytara bilər
2. `:=` tipi inferensiya edərək unexported tipin dəyişənini yarada bilir — amma sən bu tipi **açıq şəkildə** elan edə bilməzsən

**Struct sahələri:**
```go
// entities paketi:
type User struct {
    Name  string   // exported
    email string   // unexported
}

u := entities.User{
    Name:  "Bill",
    email: "bill@email.com",   // XƏTA: unknown field 'email' in struct literal
}
```

**Unexported iç tip + exported sahələr:**
```go
// entities paketi:
type user struct {     // unexported
    Name  string       // exported
    Email string
}
type Admin struct {
    user               // embedded, unexported
    Rights int
}

// başqa paketdən:
a := entities.Admin{Rights: 10}
a.Name = "Bill"       // İŞLƏYİR — promotion! İç tipin exported sahələri xaricdən görünür
a.Email = "bill@email.com"
// Amma struct literal-da iç tipi initialize etmək olmaz — user tipi görünməzdir
```

## Əsas terminlər
- User-defined Type (istifadəçi tipi)
- Struct (struktur)
- Method / Receiver (metod / qəbuledici)
- Value Receiver / Pointer Receiver (dəyər / göstərici qəbuledici)
- Primitive / Nonprimitive Nature (ilkin / qeyri-ilkin təbiət)
- Header Value (başlıq dəyəri — reference tipinin daxili strukturu)
- Interface (interfeys)
- iTable (interfeys cədvəli)
- Concrete Type (konkret tip)
- Method Set (metod dəsti)
- Polymorphism (çoxşəkillilik)
- Type Embedding (tip daxiletməsi)
- Inner Type Promotion (iç tipin qaldırılması)
- Exported / Unexported (ixrac olunmuş / olunmamış)
- Factory Function (`New` — istehsal funksiyası)
- Base Type (əsas tip)

## Praktik nəticə
- Receiver seçimi: tipin təbiətini soruş — dəyişməzsə value (time.Time kimi), dəyişənsə pointer (os.File kimi). Konkret metodun nə etdiyindən asılı olmayaraq bütün tip üzrə TUTUMLU ol.
- Pointer receiver ilə implement olunan interfeysə yalnız pointer qəbul edilir — `&u` ötürməyi vərdiş et.
- Kompozisiya üçün embedding istifadə et — inheritance yoxdur, amma override/promotion mexanizmi var.
- API dizaynında unexported tip + `New` factory funksiyası ilə daxili tətbiqi gizlət, amma dəyərlərin istifadəsinə icazə ver.
- Embedding zamanı iç tipin kimliyi qorunur — `ad.user.notify()` həmişə mümkündür.

## Mənbə
Pages: 109-148 (PDF), book pages 88-127
