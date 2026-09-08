# Chapters 5-7 — Strings, Functions, Errors (#39-#54) (səh. 133-183)

## Bu fəsillər nədən bəhs edir?

(5) Rune vs bayt, string iterasiya, strings.Builder, lazımsız konversiyalar,
substring memory leak. (6) Value vs pointer receiver, named result parametr,
nil receiver tələsi, io.Reader əvəzi fayl adı, defer arqument/receiver
qiymətləndirilməsi. (7) Panic etmək, error wrap, error tip/dəyər yoxlama, iki
dəfə handle, defer xətaları.

## Əsas səhvlər və həllər

### #36-#37: Rune anlayışı
- String = bayt ardıcıllığı (UTF-8); rune = code point
- `len("hêllo")` = BAYT sayı (6), rune sayı 5
- `for i, r := range s` — r RUNE, i bayt indexi; `s[i]` — BAYT verir ("Ã" kimi pozuntular)
- Konkret rune: `[]rune(s)[4]` (allokasiya) və ya `utf8.DecodeRuneInString`

### #38: Trim müqayisəsi
- `strings.TrimRight(s, "123")` — CUTSET (hər hansı simvol ardıcılla)
- `strings.TrimSuffix(s, "123")` — yalnız bütün suffix

### #39: String birləşdirmə
```go
// PIS: += hər addımda yeni string
// YAXŞI:
sb := strings.Builder{}
sb.Grow(total)          // əvvəlcədən yer — daha da sürətli
for _, v := range values { sb.WriteString(v) }
return sb.String()
```

### #40: Lazımsız konversiya
- I/O dünyası []byte-dır: `io.ReadAll` → `string(b)` → trim → `[]byte(s)` =
  2 əlavə kopya. `bytes.TrimSpace(b)` ilə birbaşa

### #41: Substring leak
```go
uuid := string(log[:36])  // log 1KB olsa da backing paylaşır — hamısı yaddaşda
// həll: kopya və ya strings.Clone (Go 1.18+)
```

### #42: Receiver seçimi (value vs pointer)
**Pointer:** struct mutasiya olunursa; böyük struct (kopya bahalı); metod
dəstində homojenlik (biri pointer → hamısı). **Value:** kiçik dəyişməz;
map/chan/func/pointer sahələri ZATƏN referansdır; concurrently istifadə
(nəticə təcrid olunmalı). Qeyd: value receiver də pointer sahəsini dəyişə
bilər (`c.data.balance += x` — data pointerdir).

### #43: Named result parametrlər
- `(a int, err error)` — zero-value ilə başlayır; naked return mümkün
- **Tələ:** defer/callback nəticəni dəyişir və naked return onu götürür
- Dəyər: oxunaqlılıq; risk: yan təsir → yalnız hallarda işlət

### #44: (asılı müqayisə üçün named parameter qəlibləri)

### #45: nil receiver qaytarmaq
```go
type MultiError struct{ errs []error }
func (m *MultiError) Error() string { return strings.Join(...) }
var m *MultiError        // m == nil
return m                 // interface NON-nil olur (tip=nil, dəyər nil)!
// caller: err != nil TRUE — səhv baş verir
// həll: konkret tipi yoxla: var m *MultiError; errors.As / return m.errs[..]
// və ya: if m != nil { return m }; return nil
```
- Nil pointer VALID receiverdir — `foo.Bar()` işləyir (bar çap olunur)

### #46: Fayl adı YOX, io.Reader
```go
// PIS: func countEmptyLines(filename string)
// YAXŞI: func countEmptyLines(r io.Reader) (int, error)
// test: strings.NewReader("foo\nbaz\n") — faylsız, self-contained
```

### #47: defer arqument/receiver qiymətləndirməsi
```go
var status string
defer notify(status)         // status ÖZÜ HƏMİN ANDA kopyalanır → həmişə ""
// həll 1: defer func() { notify(status) }()  // closure — icra anında oxuyur
// həll 2: defer notifyPointer(&status)
// receiver da dərhal qiymətlənir (value → kopya)
```

### #48: Panicking
- panic YALNız: proqram davam edə BİLMƏZ hallarda (init xətası, invariant
  pozuntusu); xidmət kodunda error QAYTAR

### #49: Error wrap
- `fmt.Errorf("...: %w", err)` — errors.Is/As zənciri işləsin deyə
- **Yalnız:** xətanı İSTEHLAK edən yerdə wrap et (marker mesaj əlavə etmək
  üçün yox)

### #50: Xəta TİPİ yoxlaması
- `errors.As(err, &target)` — type switch `err.(type)` əvəzinin müasir formu

### #51: Xəta DƏYƏRİ yoxlaması
- `errors.Is(err, target)` — `==` müqayisəsi əvəzinin müasir formu

### #52: İki dəfə handle
- Xəta ya LOG et, ya QAYTAR — HƏR İKİSİ YOX (duplicate logging zənciri)

### #53: Xətanı handle ETMƏMƏK
- `_ = err` QADAĞAN (yəni hətta bilərəkdən ignore — izah şərhi olmadan);
  ioutil-dən fərqli olaraq: yoxlamayan kod = kod baxışında tutulmalı

### #54: defer xətaları
```go
// f.Close() xətası yazmada itir bilər!
defer func() {
    if ferr := f.Close(); ferr != nil && err == nil {
        err = ferr   // named result + closure — yalnız əvvəlki xəta yoxdursa
    }
}()
```

## Əsas terminlər

- Rune / Code Point
- strings.Builder / Grow
- strings.Clone
- Named Result Parameter
- nil interface vs nil pointer (tip+dəyər cütü)
- errors.Is / errors.As
- Wrap (%w)
- Defer closure pattern

## Praktik nəticə

- Bayt/rune ayrımını bil: len = bayt; range = rune
- Birləşdirmə = Builder.Grow; I/O = []byte-də qal
- API-lər fayl adı deyil io.Reader qəbul etsin
- nil tip qaytarma — errors üçün struct-pointer tələsi; xəta handle =
  log XATTA return, ikisi YOX

## Mənbə

Pages: 133-183 (Chapters 5-7, 100 Go Mistakes 2nd ed.)
