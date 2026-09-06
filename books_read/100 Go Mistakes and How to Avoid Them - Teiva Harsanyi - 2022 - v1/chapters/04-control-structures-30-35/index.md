# Chapter 4 — Control structures (#30-#35)

## Bu chapter nədən bəhs edir?

Bu chapter idarə strukturları üzrə 6 səhvi əhatə edir — əsas fokus **range loop**-da: element kopiyası, ifadənin bir dəfə qiymətləndirilməsi, pointer elementləri, map iterasiyası yanlış pressumpları, break/label mexanizmi və defer-in loop daxilində istifadəsi.

---

## Əsas fikirlər

### #30: Ignoring the fact that elements are copied in range loops (Range loop-da elementlərin kopyalanması)

**Range loop nə üzərində işləyir:** string, array, array pointer, slice, map, receiving channel. Klassik for-dan fərqli olaraq indeks/terminasiya əl ilə idarə edilmir — off-by-one xətalarından qoruyur.

```go
s := []string{"a", "b", "c"}
for i, v := range s { /* index + value */ }
for _, v := range s { /* yalnız value */ }
for i := range s   { /* yalnız index */ }
```

**Kopyalama qaydası:** Go-da HƏR assignment kopyadır — struct assign-i structın kopyasını, pointer assign-i yaddaş ünvanının kopyasını yaradır (64-bit sistemdə 64 bit). Range loop hər elementi **value dəyişəninə kopyalayır**.

**Səhv nümunəsi:**

```go
type account struct{ balance float32 }

accounts := []account{{100.}, {200.}, {300.}}
for _, a := range accounts {
    a.balance += 1000     // YALNIZ KOPYANI dəyişir!
}
// accounts: [{100} {200} {300}] — dəyişmədi
```

**Həllər:**

```go
// 1. İndeks ilə (range):
for i := range accounts {
    accounts[i].balance += 1000
}
// 2. Klassik for:
for i := 0; i < len(accounts); i++ {
    accounts[i].balance += 1000
}
// 3. Pointer slice (böyük strukturlar, tez mutasiya):
accounts := []*account{{100.}, {200.}, {300.}}
for _, a := range accounts {
    a.balance += 1000    // a — pointerin kopyasıdır, eyni struct-a işarə edir
}
```

**Seçim:** hər elementi yeniləmək → 1 (qısa); spesifik elementlər (məs., hər ikincisi) → 2. Option 3-ün çatışmazlıqları: slice tipini dəyişmək lazımdır; pointer slice CPU üçün predictability azaldır (#91).

---

### #31: Ignoring how arguments are evaluated in range loops (Range arqumentlərinin qiymətləndirilməsi)

**Qayda:** `for i, v := range exp` — `exp` **YALNIZ BİR DƏFƏ, loop-dan ƏVVƏL** qiymətləndirilir (kopyası müvəqqəti dəyişənə yazılır; range bu dəyişən üzərində iterasiya edir).

**Slice nümunəsi:**

```go
s := []int{0, 1, 2}
for range s {
    s = append(s, 10)   // orijinal s böyüyür
}
// 3 iterasiyadan sonra BİTİR — range kopyası 3-length qalır

// Klassik for:
for i := 0; i < len(s); i++ {   // len(s) HƏR iterasiyada yenidən hesablanır
    s = append(s, 10)
}
// SONSUZ loop!
```

**Channel nümunəsi:**

```go
ch := ch1
for v := range ch {
    fmt.Println(v)
    ch = ch2   // TƏSİRSİZ! range hələ də ch1 üzərindədir (kopya)
}
// Çap: 0 1 2 (ch2 elementləri YOX)
// AMMA: sonrakı close(ch) artıq ch2-ni bağlayar
```

**Array nümunəsi:**

```go
a := [3]int{0, 1, 2}
for i, v := range a {
    a[2] = 10           // orijinal array dəyişir
    if i == 2 {
        fmt.Println(v)  // 2 çap olunur — v kopyadan gəlir!
    }
}
```

Həllər:
1. İndekslə birbaşa: `fmt.Println(a[2])` → 10.
2. Array pointer: `for i, v := range &a` — hər iki pointer eyni array-ə işarə edir → `v` = 10. Üstünlük: böyük array-in tam kopyası YOXDUR.

---

### #32: Ignoring the impact of using pointer elements in range loops (Pointer elementlər range loop-da)

**Pointer istifadəsinin 3 səbəbi:**
1. **Semantika:** elementin paylaşılması (cache-də saxlayarkən caller ilə Store eyni *Foo-nu paylaşır).
2. Artıq pointerlərin mövcudluğu.
3. **Böyük strukt + tez mutasiya:** kopya+geriyazma əvəzinə birbaşa `mapPointer[id].foo = "bar"`.

**Səhv — loop dəyişəninə pointer saxlamaq:**

```go
func (s *Store) storeCustomers(customers []Customer) {
    for _, customer := range customers {
        s.m[customer.ID] = &customer   // BÜTÜN key-lər EYNİ pointeri saxlayır!
    }
}
// Nəticə:
// key=1, value=&Customer{ID:"3", Balance:0}
// key=2, value=&Customer{ID:"3", Balance:0}
// key=3, value=&Customer{ID:"3", Balance:0}
```

**Səbəb:** `customer` dəyişəni loop boyu **tək, sabit ünvanlıdır** (eyni ünvan: `0xc000096020` hər iterasiyada). Son tapşırıq — sonuncu element (Customer 3).

**Həllər:**

```go
// 1. Lokal dəyişən (hər iterasiyada yeni ünvan):
for _, customer := range customers {
    current := customer          // yeni ünvan
    s.m[current.ID] = &current
}
// 2. İndekslə elementə pointer:
for i := range customers {
    customer := &customers[i]    // slice elementinin öz ünvanı
    s.m[customer.ID] = customer
}
```

**Prinsip:** range-də bütün dəyərlər tək ünvanlı dəyişənə assign olunur — hər iterasiyada pointer saxlayırsansa, eyni (sonuncu) pointeri saxlayırsan. Map üçün də eyni problem keçərlidir.

---

### #33: Making wrong assumptions during map iterations (Map iterasiyasında yanlış pressumplar)

**Map-in 4 qeyri-zəmanəti:**
1. Key-lərə görə sıralı deyil (binary tree deyil).
2. Insert sırasını qorumur.
3. **Iterasiya sırası deterministik deyil** — eyni map, iki range: `zdyaec`, sonra `czyade`.
4. Elementlərin cari saxlanma sırası da zəmanət verilmir.

**Niyə belə:** Dil dizaynerlərinin ŞÜURLU seçimi — developer-lər heç vaxt sıra pressumplasına güvənməsin. Rəsmi spesifikasiya: iterasiya **unspecified**-dir (random deyil — distribution uniform deyil).

**İstisna:** `encoding/json` map-i JSON-a marshal edərkən key-ləri əlifba sırasına düzür — bu Go map-in deyil, json paketinin xüsusiyyətidir. Sıralama lazımdırsa: binary heap kimi strukturlar (GoDS: github.com/emirpasic/gods).

**Iterasiya zamanı insert:**

```go
m := map[int]bool{0: true, 1: false, 2: true}
for k, v := range m {
    if v {
        m[10+k] = true    // yeni elementlər bu iterasiyada ÇIXA BİLƏR və ya ÇIXMAYA BİLƏR
    }
}
// 3 fərli icra → 3 fərli nəticə (non-deterministic!)
```

**Go spesifikasiyası:** *"If a map entry is created during iteration, it may be produced during the iteration or skipped. The choice may vary for each entry created and from one iteration to the next."*

**Həll — kopya üzərində yaz:**

```go
m2 := copyMap(m)
for k, v := range m {
    m2[k] = v
    if v {
        m2[10+k] = true   // oxunan map (m) ≠ yazılan map (m2) → predictible
    }
}
```

---

### #34: Ignoring how the break statement works (Break-in işləməsi)

**Qayda:** break **ən daxili** for/switch/select ifadəsini sonlandırır.

```go
for i := 0; i < 5; i++ {
    fmt.Printf("%d ", i)
    switch i {
    default:
    case 2:
        break   // SWITCH-i qırır, loop-u YOX → 0 1 2 3 4 çap olunur
    }
}
```

**Həll — label:**

```go
loop:
for i := 0; i < 5; i++ {
    fmt.Printf("%d ", i)
    switch i {
    default:
    case 2:
        break loop    // FOR loop-u qırır → 0 1 2
    }
}
```

**Select halı (çox rast gəlinən):**

```go
// PİS: break select-i qırır, loop sonsuz:
for {
    select {
    case <-ch:
        // ...
    case <-ctx.Done():
        break
    }
}

// DÜZGÜN:
loop:
for {
    select {
    case <-ch:
        // ...
    case <-ctx.Done():
        break loop
    }
}
```

**Label ≠ goto:** İdiomatik yanaşmadır — standart kitabxanada belə istifadə olunur (net/http: `readlines:` label ilə buffer sətirlərinin oxunması). `continue` də label qəbul edir.

---

### #35: Using defer inside a loop (Loop daxilində defer)

**Qayda:** defer çağırışını əhatə edən FUNKSiya qayıdanda icra edir.

```go
// PİS — file descriptor leak:
func readFiles(ch <-chan string) error {
    for path := range ch {
        file, err := os.Open(path)
        if err != nil {
            return err
        }
        defer file.Close()   // readFiles qayıdanadək YIĞILIR!
        // ...
    }
    return nil
}
// Loop bitməzsə (channel açıqdırsa) fayllar HƏMİŞƏAÇIQ qalır
```

**Həll — hər iterasiyada çağırılan əhatəedici funksiya:**

```go
func readFiles(ch <-chan string) error {
    for path := range ch {
        if err := readFile(path); err != nil {
            return err
        }
    }
    return nil
}

func readFile(path string) error {
    file, err := os.Open(path)
    if err != nil {
        return err
    }
    defer file.Close()   // readFile qayıdanda — hər iterasiyanın sonunda
    // Do something with file
    return nil
}
```

**Alternativ — closure** (eyni məntiq, sadəcə anonim):

```go
for path := range ch {
    err := func() error {
        // ...
        defer file.Close()
        // ...
    }()
    if err != nil {
        return err
    }
}
```

Adi funksiya daha aydındır + unit test yazıla bilər. **Performance-kritik halda:** funksiya çağırış overhead-i istenmirsə — defer-i at, bağlanmanı əl ilə idarə et.

---

## Əsas terminlər

| Termin | İzah |
|--------|------|
| Range loop | String/array/array-ptr/slice/map/channel üzərində idiomatik iterasiya |
| Value copy | Range-də hər elementin value dəyişəninə KOPYALANMASI |
| Blank identifier `_` | İstifadə olunmayan index/value-i söndürmək |
| Expression evaluation | Range ifadəsinin loop-dan əvvəl BİR DƏFƏ kopyalanması |
| `len(s)` re-evaluation | Klassik for-da hər iterasiyada yenidən hesablanır |
| `range &a` | Array pointeri üzərində range — tam kopyadan qaçış |
| Fixed-address loop variable | Range value dəyişəni tək sabit ünvanlıdır — pointer tələsi |
| Map iteration order | Unspecified — sıralı deyil, insert sırası yox, deterministic deyil |
| "may be produced or skipped" | İterasiya zamanı əlavə edilən elementin qeyri-müəyyən görünməsi |
| Label (`loop:`) | Break/continue-in hədəf ifadəsini aşkar təyin etməsi |
| Innermost statement | Break-in default hədəfi — ən daxili for/switch/select |
| defer stacking | Loop daxilində defer-lərin funksiya qayıdışına qədər yığılması |
| File descriptor leak | Bağlanmayan fayl deskriptorlarının yığılması |

---

## Praktik nəticə

1. **Range value = kopya:** struct mutasiyası üçün `accounts[i]` (index) və ya klassik for; pointer slice (`[]*T`) yalnız böyük+tez-dəyişən strukturlarda.
2. **Range ifadəsi BİR DƏFƏ qiymətləndirilir:** loop daxilində slice böyütmək/kanal dəyişmək range-ə təsir etmir; klassik for-da `len(s)` canlıdır — sonsuz loop riski.
3. **Loop dəyişəninə pointer saxlaMA:** `&customer` → hər key sonuncu elementi görəcək; `current := customer` və ya `&customers[i]`.
4. **Map iterasiyasından heç bir sıra gözləmə:** json marshal əlifba düzür (paket xüsusiyyəti); sıralama lazımdırsa binary heap/GoDS.
5. **Iterasiya + map yazma = kopya patterni:** oxunan map ilə yazılan map-i ayır → predictible nəticə.
6. **Switch/select loop daxilində break = label yaz:** `break loop`; label idiomatikdir (net/http istifadə edir), goto deyil.
7. **Loop daxilində defer YOX:** hər iterasiyada qayıdan köməkçi funksiya (və ya closure) içərisinə köçür; performance-kritikdirsə defer-siz əl idarəsi.
8. **`range &a` ilə array kopyasından qaç:** böyük arraylərdə həm düzgünlük həm performans.

---

## Mənbə

- Kitab: *100 Go Mistakes and How to Avoid Them* — Teiva Harsanyi, Manning, 2022 (ISBN 9781617299599)
- Chapter 4: "Control structures", book səh. 95–112
- PDF səhifələri: 115–132
- İstinadlar: Go map iteration randomization müzakirəsi (mng.bz/M2JW); GoDS (github.com/emirpasic/gods)
