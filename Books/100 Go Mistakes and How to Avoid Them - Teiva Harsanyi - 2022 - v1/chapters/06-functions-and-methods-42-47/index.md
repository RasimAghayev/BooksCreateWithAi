# Chapter 6 — Functions and methods (#42-#47)

## Bu chapter nədən bəhs edir?

Bu chapter funksiyalar və metodlar üzrə 6 səhvi əhatə edir: value/pointer receiver seçimi, named result parameters (faydaları və yan təsirləri), nil receiver qaytarma tələsi, filename əvəzinə `io.Reader` qəbul etmək və defer-də arqument/receiver qiymətləndirilməsi.

---

## Əsas fikirlər

### #42: Not knowing which type of receiver to use (Hansı receiver tipi?)

**Value receiver:** Go dəyərin KOPYASINI metoda ötürür — dəyişikliklər lokal qalır:

```go
func (c customer) add(v float64) {
    c.balance += v        // KOPYA dəyişir
}
c := customer{balance: 100.}
c.add(50.)
fmt.Println(c.balance)   // 100.00 — dəyişmədi
```

**Pointer receiver:** Obyektin ünvanı ötürülür (pass-by-reference Go-da YOXDUR — pointer-in kopyası ötürülür) — modifikasiyalar orijinalda:

```go
func (c *customer) add(operation float64) {
    c.balance += operation
}
c := customer{balance: 100.0}
c.add(50.0)
fmt.Println(c.balance)   // 150.00
```

**Qərar cədvəli:**

| Receiver... | Şərt |
|-------------|------|
| **MÜTLƏQ pointer** | Metod receiver-i MUTASİYA edirsə (slice append də daxil: `func (s *slice) add(e int) { *s = append(*s, e) }`); receiver kopyalana bilməyən sahə daşıyırsa (sync paketi tipləri — #74) |
| **Pointer olmalıdır (should)** | Receiver böyük obyektdir — kopyadan qaçmaq (nə qədər "böyük"? → benchmark) |
| **MÜTLƏQ value** | İmmutability məcburiyyəti; receiver map, funksiya və ya channel (əks halda kompilyasiya xətası) |
| **Value olmalıdır (should)** | Mutasiya olunmayan slice; kiçik array/struct (təbii value tipi — `time.Time`); əsas tiplər (int, float64, string) |

**Xüsusi hal — pointer sahəli struct:**

```go
type customer struct { data *data }
type data struct { balance float64 }

func (c customer) add(operation float64) {
    c.data.balance += operation   // value receiver amma mutasiya İŞLƏYİR!
}
c := customer{data: &data{balance: 100}}
c.add(50.)
fmt.Println(c.data.balance)   // 150.00
```

c-nin kopyası götürülse də `data` pointer-i eyni struct-a işarə edir. Amma **aydınlıq üçün** pointer receiver daha yaxşıdır — customer bütövlükdə mutable olduğu bildirilir.

**Default:** value receiver; yaxşı səbəb yoxdursa. **Şübhə halında: pointer receiver.**

**Receiver tiplərinin qarışması:** Ümumi qadağadır, amma 100% deyil — `time.Time` əksəriyyəti value (`After`, `IsZero`, `UTC`), amma `UnmarshalBinary` (encoding.TextUnmarshaler) pointer receiver olmalıdır (mutasiya edir).

---

### #43: Never using named result parameters (Named result parametrləri)

**Mexanizm:** Nəticə parametrlərinə ad verildikdə funksiya başlanğıcında **zero value** ilə initializə olunur; `return` (naked) cari dəyərləri qaytarır:

```go
func f(a int) (b int) {
    b = a
    return    // naked return — b-nin cari dəyəri
}
```

**Faydalı olduğu hal — oxunaqlıq:**

```go
// Pis — 2 float32 nədir? Hansı birinci?
type locator interface {
    getCoordinates(address string) (float32, float32, error)
}

// Yaxşı — məna imzada görünür:
type locator interface {
    getCoordinates(address string) (lat, lng float32, err error)
}
```

**Faydasız olduğu hal:** `func StoreCustomer(c Customer) (err error)` — tək `err` adı heç nə əlavə etmir → istifadə ETMƏ.

**Convenience halı (Effective Go, io.ReadFull):**

```go
func ReadFull(r io.Reader, buf []byte) (n int, err error) {
    for len(buf) > 0 && err == nil {
        var nr int
        nr, err = r.Read(buf)
        n += nr
        buf = buf[nr:]
    }
    return    // n və err onsuz da zero-dan başlayır
}
```

**Qaydalar:**
- Eyni tipli çoxluq nəticələr → named parameters (və ya adətharə field-lərə struct).
- Naked return yalnız QISA funksiyalarda; uzun funksiyada oxucu çıxışları xatırlamalı olur.
- Bir funksiyada naked və arqumentli return-ları QARIŞDIRMA.
- Named parameters ≠ naked return — yalnız imza aydınlığı üçün də istifadə oluna bilər.
- Şübhə halında istifadə ETMƏ.

---

### #44: Unintended side effects with named result parameters (Named parameters yan təsirləri)

**Tələ:** zero-value initializasiyası səhvləri gizlədir:

```go
func (l loc) getCoordinates(ctx context.Context, address string) (
    lat, lng float32, err error) {

    isValid := l.validateAddress(address)
    if !isValid {
        return 0, 0, errors.New("invalid address")
    }
    if ctx.Err() != nil {
        return 0, 0, err    // err hələ NIL-dır (zero value)! → həmişə nil error qaytarır!
    }
    // ...
}
```

`ctx.Err() != nil` olduqda belə funksiya **nil error qaytarır** — çünki `err`-ə heç nə mənimsədilməyib. Kompilyasiya xətası YOXDUR (err elan olunub!). Adlı olmasaydı: `Unresolved reference 'err'`.

**Həll 1 — kölgələmə ilə assign:**

```go
if err := ctx.Err(); err != nil {
    return 0, 0, err
}
```

**Həll 2 — naked return (lakin qarışdırma qaydasını pozur):**

```go
if err = ctx.Err(); err != nil {
    return
}
```

**Moral:** Named result parameters oxunaqlılıq verir, amma zero-value başlanğıcı səhvləri gizlədə bilər — ehtiyatlı istifadə et.

---

### #45: Returning a nil receiver (Nil receiver qaytarmaq)

**Ən geniş yayılmış Go tələlərindən biri:**

```go
type MultiError struct{ errs []string }
func (m *MultiError) Add(err error) { m.errs = append(m.errs, err.Error()) }
func (m *MultiError) Error() string { return strings.Join(m.errs, ";") }

func (c Customer) Validate() error {
    var m *MultiError          // nil pointer
    if c.Age < 0 {
        m = &MultiError{}
        m.Add(errors.New("age is negative"))
    }
    if c.Name == "" {
        if m == nil { m = &MultiError{} }
        m.Add(errors.New("name is nil"))
    }
    return m                   // m == nil ola bilər → amma error != nil OLACAQ!
}

// Test:
customer := Customer{Age: 33, Name: "John"}   // valid!
if err := customer.Validate(); err != nil {
    log.Fatalf("customer is invalid: %v", err)  // İŞƏ DÜŞÜR! → "<nil>"
}
```

**Səbəb — interface wrapper semantikası:**
1. **Nil pointer VALID receiver-dir:** `var foo *Foo; foo.Bar()` kompilyasiya olunur və işləyir — metod = receiver-i birinci parametr kimi qəbul edən funksiya: `func Bar(foo *Foo) string`.
2. `return m` — `*MultiError` (nil) → `error` interface-inə **konversiya** olunur. Interface = dispatch wrapper: **wrappee (içəridəki pointer) nil, wrapper (interface) isə NON-NIL!**

Nəticə: çağıran həmişə non-nil error alır.

**Həll — nil-i aşkar qaytar:**

```go
func (c Customer) Validate() error {
    var m *MultiError
    // ... yoxlamalar ...
    if m != nil {
        return m
    }
    return nil    // nil İNTERFEYS, nil pointer-in konvertasiyası YOX
}
```

**Prinsip:** Interface qaytaranda nil POINTER yox, birbaşa nil qaytar. Yalnız error-lərə xas deyil — pointer receiver ilə implement olunan HƏR HANSI interfeysdə eyni təhlükə.

---

### #46: Using a filename as a function input (Filename funksiya girişi kimi)

```go
// PİS — filename qəbul edir:
func countEmptyLinesInFile(filename string) (int, error) {
    file, err := os.Open(filename)
    if err != nil { return 0, err }
    // ...
    scanner := bufio.NewScanner(file)
    for scanner.Scan() { /* ... */ }
}

// Testlər üçün hər case-ə FAYL yaratmaq lazım!
// HTTP body versiyası lazım olsa — eyni logic-i DUBLIKASIYALAMAQ:
func countEmptyLinesInHTTPRequest(request http.Request) (int, error) {
    scanner := bufio.NewScanner(request.Body)
    // Copy the same logic
}
```

**İdiomatik həll — `io.Reader` abstraksiyası:**

```go
func countEmptyLines(reader io.Reader) (int, error) {
    scanner := bufio.NewScanner(reader)
    for scanner.Scan() { /* ... */ }
}
```

**Faydaları:**
1. **Data mənbəyi abstraksiyası:** fayl, HTTP request, socket, gRPC — fərq etmir; `*os.File` və `http.Request.Body` hər ikisi `io.Reader` implement edir.
2. **Test asanlığı:** fayl yaratmadan `strings.NewReader` ilə self-contained testlər:

```go
func TestCountEmptyLines(t *testing.T) {
    emptyLines, err := countEmptyLines(strings.NewReader(`foo
bar
baz
`))
    // Test logic
}
```

**Prinsip:** Filename input = code smell (`os.Open` kimi xüsusi funksiyalar istisna). Fayl oxuyan funksiya dizayn edəndə `io.Reader`-dən başla.

---

### #47: Ignoring how defer arguments and receivers are evaluated

**Qayda:** Defer funksiyasının ARQUMENTLƏRİ DƏRHAL (defer çağrılanda) qiymətləndirilir — əhatəedici funksiya qayıdanda YOX.

**Səhv nümunəsi:**

```go
func f() error {
    var status string
    defer notify(status)            // status indi qiymətləndirilir → ""
    defer incrementCounter(status)   // həmçinin ""
    if err := foo(); err != nil {
        status = StatusErrorFoo
        return err
    }
    if err := bar(); err != nil {
        status = StatusErrorBar
        return err
    }
    status = StatusSuccess
    return nil
}
// Hər execution path-də notify/incrementCounter BOŞ string alır!
```

**Həll 1 — pointer ötür:**

```go
defer notify(&status)
defer incrementCounter(status)  →  defer incrementCounter(&status)
// status dəyişir, amma ÜNVAN sabitdir → defer çağrılan zaman cari dəyər oxunur
```

Çatışmazlıq: hər iki funksiyanın imzasını dəyişmək lazımdır.

**Həll 2 — closure (imza dəyişməz):**

```go
defer func() {
    notify(status)            // status closure icrası zamanı qiymətləndirilir
    incrementCounter(status)
}()
```

**Closure qiymətləndirmə qaydası:** Arqument kimi ötürülən dəyişənlər DƏRHAL; bədəndən kənar referans olunan dəyişənlər closure İCRASINDA:

```go
i := 0
j := 0
defer func(i int) { fmt.Println(i, j) }(i)   // i dərhal (0), j icrada (1)
i++
j++
// Çap: 0 1
```

**Metod receiver-ləri — eyni qayda (dərhal qiymətləndirilir):**

```go
// Value receiver — SONRAKİ dəyişiklik GÖRÜNMÜR:
s := Struct{id: "foo"}
defer s.print()     // s-nin kopyası DƏRHAL götürülür
s.id = "bar"
// Çap: foo

// Pointer receiver — dəyişiklik GÖRÜNÜR:
s := &Struct{id: "foo"}
defer s.print()     // pointer dərhal kopyalanır, amma eyni struct-a işarə
s.id = "bar"
// Çap: bar
```

---

## Əsas terminlər

| Termin | İzah |
|--------|------|
| Value receiver | Metoda dəyərin kopyası ötürülür — mutasiya lokal |
| Pointer receiver | Ünvan ötürülür — orijinal mutasiya olunur (pass-by-reference YOXDUR) |
| Receiver tiplərinin qarışması | Ümumi qadağa; time.Time istisna (UnmarshalBinary) |
| Named result parameter | Adlandırılmış nəticə — zero value ilə başlayır |
| Naked return | Arqumentsiz return — cari named dəyərlər |
| Zero-value initializasiya | Named parametrlərin başlanğıc dəyəri — səhv gizlədə bilər |
| Interface wrapper | Interface = dispatch wrapper; nil wrappee ≠ nil wrapper |
| Nil pointer receiver | Valid receiver — metod = receiver-li funksiya |
| `io.Reader` | Mənbəyə-agnostik oxu abstraksiyası — filename əvəzinə |
| `strings.NewReader` | String-dən io.Reader — faylsız testlər |
| Immediate evaluation | Defer arqumentləri defer ZAMANINDA qiymətləndirilir |
| Closure defer | Kənar dəyişənlər closure icrası zamanı qiymətləndirilir |
| Deferred receiver copy | Value receiver defer zamanı dərhal kopyalanır |

---

## Praktik nəticə

1. **Receiver seçim alqoritmi:** Mutasiya/kopyalanmayan sahə → pointer (mütləq); böyük obyekt → pointer; map/func/chan → value (mütləq); immutability/kiçik struct/əsas tip → value. Şübhə → pointer.
2. **Named result parametrləri imzada:** Eyni tipli nəticələrin mənasını açır (`lat, lng float32`); tək `err` üçün lazımsızdır.
3. **Named parameter + `return 0, 0, err` → err-i yoxla:** zero-value tələsi; `if err := ctx.Err(); err != nil` patternini istifadə et.
4. **`return m` (nil pointer) `error` kimi → həmişə non-nil:** Sonda `if m != nil { return m }; return nil`.
5. **Filename yox, `io.Reader`:** Testlər faylsız, funksiyalar mənbəyə-agnostik, dublikasiya sıfır.
6. **Defer arqumentləri DƏRHAL:** Dəyişən dəyəri lazımdırsa — pointer ötür və ya closure istifadə et.
7. **Defer + metod:** Value receiver → sonrakı dəyişikliklər görünmür; pointer receiver → görünür.

---

## Mənbə

- Kitab: *100 Go Mistakes and How to Avoid Them* — Teiva Harsanyi, Manning, 2022 (ISBN 9781617299599)
- Chapter 6: "Functions and methods", book səh. 126–142
- PDF səhifələri: 146–162
- İstinadlar: Effective Go (go.dev/doc/effective_go)
