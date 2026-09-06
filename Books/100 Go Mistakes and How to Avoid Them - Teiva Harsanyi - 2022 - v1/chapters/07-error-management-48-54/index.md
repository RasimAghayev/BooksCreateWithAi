# Chapter 7 — Error management (#48-#54)

## Bu chapter nədən bəhs edir?

Bu chapter xəta idarəsinin 7 səhvini əhatə edir: panic-in düzgün/halal istifadəsi, error wrapping qərarları (%w vs %v), error TİPİNİN yoxlanması (errors.As), error DƏYƏRİNİN yoxlanması (errors.Is, sentinel errors), xətanın iki dəfə idarəsi, xətanın bilərəkdən iqnoru (`_`) və defer xətalarının idarəsi (named result parameters ilə propagasiya).

---

## Əsas fikirlər

### #48: Panicking (Panic-ləmək)

**Panic:** adi axını dayandıran built-in funksiya; call stack boyu yuxarı qalxır — goroutine bitənə və ya `recover()` tutana qədər:

```go
func main() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("recover", r)
        }
    }()
    f()
}
func f() {
    fmt.Println("a")
    panic("foo")
    fmt.Println("b")   // icra olunmur
}
// a
// recover foo
```

`recover()` yalnız **defer funksiya daxilində** işləyir (defer-lər panic zamanı da icra olunur).

**Panic HALLAL olan 2 hal:**

1. **Programmer error (proqramçı xətası):**

```go
// net/http: yanlış status kodu = programmer error
func checkWriteHeaderCode(code int) {
    if code < 100 || code > 999 {
        panic(fmt.Sprintf("invalid WriteHeader code %v", code))
    }
}

// database/sql: nil/dublikat driver qeydiyyatı
func Register(name string, driver driver.Driver) {
    if driver == nil {
        panic("sql: Register driver is nil")
    }
    if _, dup := drivers[name]; dup {
        panic("sql: Register called twice for driver " + name)
    }
    drivers[name] = driver
}
```

2. **Məcburi asılılığın ilkinləşməməsi:**

```go
// Email validasiyası üçün regex — MƏCBURİ asılılıq:
// regexp.Compile → (re, err); MustCompile → re, PANIC on error
// Compile alınmırsa heç bir input validasiya oluna bilməz → MustCompile + panic haqlıdır
```

**Prinsip:** Panicə güzəştli istisna (exceptional) hallarda; əks halda error qaytaran funksiya.

---

### #49: Ignoring when to wrap an error (Error wrapping qərarları)

**Wrapping:** source error-u wrapper konteynerə qoyub AÇIQ saxlamaq. **2 use case:**
1. **Context əlavə etmək** — "user X resource Y-yə girəndə permission denied"
2. **Xətanı markalamaq** — Forbidden tipli etiket (handler 403 qaytarsın)

**4 qaytarma seçimi:**

```go
if err != nil {
    // 1. Birbaşa qaytar — marka yox, context lazımsız:
    return err

    // 2. Custom error tipi (Go 1.13-əvvəli wrapping):
    return BarError{Err: err}

    // 3. %w (Go 1.13+) — source error AÇIQ qalır:
    return fmt.Errorf("bar failed: %w", err)

    // 4. %v — transform, source error BAĞLANIR:
    return fmt.Errorf("bar failed: %v", err)
}
```

| Seçim | Context | Marka | Source açıq |
|-------|---------|-------|-------------|
| Birbaşa return | Yox | Yox | Bəli |
| Custom type | Mümkün | Bəli | Mümkün |
| `%w` | Bəli | Yox | **Bəli** |
| `%v` | Bəli | Yox | **Yox** |

**Niyə %v hələ də lazımdır?** Wrapping = **coupling (qoşulma) riski**: caller source error-ün tipinə/yə baxa bilir; implementasiya dəyişsə (başqa funksiya → başqa error tipi) caller-in yoxlaması QIRILAR. Source error implementation detail-dirsə → `%v` ilə transform et, coupling yaratma.

---

### #50: Checking an error type inaccurately (Error TİPİ yoxlaması — errors.As)

**Ssenari:** Handler transient error-a 503, digərlərinə 400 qaytarır:

```go
type transientError struct{ err error }
func (t transientError) Error() string {
    return fmt.Sprintf("transient error: %v", t.err)
}

func handler(w http.ResponseWriter, r *http.Request) {
    amount, err := getTransactionAmount(transactionID)
    if err != nil {
        switch err := err.(type) {
        case transientError:
            http.Error(w, err.Error(), http.StatusServiceUnavailable)
        default:
            http.Error(w, err.Error(), http.StatusBadRequest)
        }
        return
    }
    // ...
}
```

**Tələ:** Refactoringdə `transientError` bir qat dərinləşir və üst qat `%w` ilə wrap edirsə:

```go
// getTransactionAmount artıra DİREKT transientError YOX, WRAPPED error qaytarır:
return 0, fmt.Errorf("failed to get transaction %s: %w", transactionID, err)
// switch-in case transientError → HƏMİŞƏ false → hər xəta 400!
```

**Həll — `errors.As` (Go 1.13):** Error zəncirini REKURSİV unwrap edib tipə uyğunlaşdırır — birbaşa YOXSA wrap olunmuş fərq etməz:

```go
if err != nil {
    if errors.As(err, &transientError{}) {
        http.Error(w, err.Error(), http.StatusServiceUnavailable)
    } else {
        http.Error(w, err.Error(), http.StatusBadRequest)
    }
    return
}
```

**Vacib:** errors.As-in 2-ci arqumenti POINTER olmalıdır (`&transientError{}`) — əks halda kompilyasiya olunur amma run-time-da panic edir.

---

### #51: Checking an error value inaccurately (Error DƏYƏRİ yoxlaması — errors.Is)

**Sentinel error:** Qlobal error dəyişəni — gözlənilən (expected) xətanı daşıyır:

```go
var ErrFoo = errors.New("foo")   // konvensiya: Err prefiks
```

**Expected vs Unexpected:**
- **Expected error** → sentinel error DƏYƏRİ: `sql.ErrNoRows` (sorğu 0 sətir qaytardı — normal hal), `io.EOF` (oxu bitdi)
- **Unexpected error** → custom error TİPİ: network xətaları və s.

```go
// PİS — == operatoru wrap olunmuş sentinel-i TUTA BİLMƏZ:
if err == sql.ErrNoRows { /* wrap olsa həmişə false! */ }

// DÜZGÜN — errors.Is: zənciri rekursiv unwrap edərək müqayisə edir:
if errors.Is(err, sql.ErrNoRows) { /* wrap olunmuş olsa da true */ }
```

**Cədvəl qaydası:** `%w` wrapping istifadə edirsənsə → tip üçün `errors.As`, dəyər üçün `errors.Is` — == və type switch KİFAYƏT ETMİR.

---

### #52: Handling an error twice (Xətanı iki dəfə idarə etmək)

**Qayda:** Xəta YALNIZ BİR DƏFƏ idarə olunmalıdır. **Logging = idarəetmə. Return = idarəetmə.** Hər ikisi birdən — YOX.

**Pis nümunə — hər səviyyədə log + return:**

```go
func GetRoute(...) (Route, error) {
    err := validateCoordinates(srcLat, srcLng)
    if err != nil {
        log.Println("failed to validate source coordinates")   // 2-ci log
        return Route{}, err
    }
    // ...
}
func validateCoordinates(lat, lng float32) error {
    if lat > 90.0 || lat < -90.0 {
        log.Printf("invalid latitude: %f", lat)    // 1-ci log
        return fmt.Errorf("invalid latitude: %f", lat)
    }
    // ...
}
// Nəticə — 2 log sətri:
// 2021/06/01 invalid latitude: 200.000000
// 2021/06/01 failed to validate source coordinates
```

**Problemlər:** Konkurrent çağrılarda 2 sətir ardıcıllığı pozula bilər → debugging çətinləşir; mesajlar dublikasiya olunur.

**Həll — log-ları sil, context-i wrap ilə ver:**

```go
func GetRoute(srcLat, srcLng, dstLat, dstLng float32) (Route, error) {
    err := validateCoordinates(srcLat, srcLng)
    if err != nil {
        return Route{},
            fmt.Errorf("failed to validate source coordinates: %w", err)
    }
    // ... target eyni şəkildə
    return getRoute(srcLat, srcLng, dstLat, dstLng)
}

func validateCoordinates(lat, lng float32) error {
    if lat > 90.0 || lat < -90.0 {
        return fmt.Errorf("invalid latitude: %f", lat)   // YALNIZ return
    }
    // ...
}
// Tək log (caller-da): failed to validate source coordinates: invalid latitude: 200.000000
```

Tək log + itirilmiş məlumat yox (source/target fərqi wrap-də saxlanılır).

---

### #53: Not handling an error (Xətanı iqnor etmək)

```go
// PİS — niyyət bilinmir (unudulubmu? bilərəkdənmi?):
func f() {
    notify()   // error nəsə oldusa da YOX SAYILIR — oxucu şübhələnir
}

// DÜZGÜN — bilərəkdən iqnor AŞKAR bildirilir:
_ = notify()
```

**Rationale commenti (kodun təkrarı deyil — SƏBƏB):**

```go
// At-most once delivery.
// Hence, it's accepted to miss some of them in case of errors.
_ = notify()

// PİS comment: // Ignore the error  → kodun dublikatı, lazımsız
```

**Prinsip:** İqnor istisnadır; mümkünsə low log level-də logla. İqnor mütləqdirsə — blank identifier ilə aşkarla.

---

### #54: Not handling defer errors (Defer xətaları)

```go
func getBalance(db *sql.DB, clientID string) (float32, error) {
    rows, err := db.Query(query, clientID)
    if err != nil {
        return 0, err
    }
    defer rows.Close()   // Close() error qaytarır — İQNOR! (#53 pozuntusu)
    // Use rows
}
```

**Səviyyələr:**

1. **Aşkar iqnor:**

```go
defer func() { _ = rows.Close() }()
```

2. **Log (DB connection pool azad edilmədi — bilmək lazımdır):**

```go
defer func() {
    err := rows.Close()
    if err != nil {
        log.Printf("failed to close rows: %v", err)
    }
}()
```

3. **Propagasiya — named result parameters ŞƏRTİDİR:**

```go
// Bu KOMPİLİYA OLUNMUR — return anonim funksiyaya aiddir, getBalance-ə YOX:
defer func() {
    err := rows.Close()
    if err != nil {
        return err
    }
}()

// Düzgün — named err:
func getBalance(db *sql.DB, clientID string) (balance float32, err error) {
    rows, err := db.Query(query, clientID)
    if err != nil {
        return 0, err
    }
    defer func() {
        err = rows.Close()    // err named parameterə yazılır → caller-a çatır
    }()
    if rows.Next() {
        err := rows.Scan(&balance)
        if err != nil {
            return 0, err
        }
        return balance, nil
    }
    // ...
}
```

**Sub-tələ:** Yuxarıdakı sadə variantda `rows.Scan` error qaytarsa belə `rows.Close`-un uğuru `err`-i ÜSTÜNƏ YAZIR → scan xətası İTİRİLİR (Close uğurludursa nil qaytarır!).

**Düzgün məntiq cədvəli:**

| rows.Scan | rows.Close | Nə qaytarılır |
|-----------|------------|---------------|
| OK | OK | nil |
| OK | Fail | closeErr |
| Fail | OK | scan err |
| Fail | Fail | scan err (əsas) + close err (log) |

**Final implementasiya:**

```go
defer func() {
    closeErr := rows.Close()
    if err != nil {
        // əsas xəta PRIORİTETDİR; close xətasını logla
        if closeErr != nil {
            log.Printf("failed to close rows: %v", err)
        }
        return
    }
    err = closeErr
}()
```

---

## Əsas terminlər

| Termin | İzah |
|--------|------|
| panic/recover | Axını dayandırma; recover yalnız defer daxilində tutur |
| Programmer error | İstifadəçi deyil, proqramçının xətası — panic üçün haqlı səbəb |
| `MustCompile` | Panikləşən variant — məcburi asılılıq qurulumu üçün |
| Error wrapping | Source error-u wrapper-də AÇIQ saxlamaq (%w) |
| Error transformation | Source-u bağlamaq (%v) — coupling qoruması |
| Coupling risk | Açıq source error caller-i implementasiyaya bağlayır |
| `errors.As(err, &target)` | Zəncir boyu tip yoxlaması — pointer hədəf məcburidir |
| `errors.Is(err, sentinel)` | Zəncir boyu dəyər yoxlaması |
| Sentinel error | Qlobal error dəyəri — expected xətaları daşıyır (ErrFoo, sql.ErrNoRows, io.EOF) |
| Expected error | Gözlənilən hal (0 sətir) → sentinel; unexpected → custom tip |
| Handle once | Xəta bir dəfə idarə olunar: log YAXUD return |
| Blank identifier iqnoru | `_ = notify()` — niyyətli iqnorun aşkar formasi |
| Rationale comment | İqnorun SƏBƏBİNİ izah edən komment |
| `rows.Close()` | Closer interfeysi — error qaytaran təmizləmə |
| Named err + defer | Defer xətasının caller-a propagasiya yolu |

---

## Praktik nəticə

1. **Panic yalnız 2 halda:** programmer error (net/http, database/sql patterni) və ya məcburi asılılıq alınmadıqda (MustCompile). Əks halda error qaytar.
2. **Wrapping qərarı:** context/marka lazımdırsa → `%w` və ya custom tip; source implementation detail-dirsə → `%v` transform (coupling yaratma).
3. **Tip yoxlaması:** type switch YOX → `errors.As(err, &T{})` — wrapping ilə işləyir; hədəf mütləq pointer.
4. **Dəyər yoxlaması:** `==` YOX → `errors.Is(err, sql.ErrNoRows)` — wrapping ilə işləyir.
5. **Expected → sentinel value, unexpected → error type.**
6. **Log YAXUD return, heç vaxt hər ikisi:** İki log = konkurrent debugging xaosu; context-i `%w` wrap ilə ver.
7. **İqnor = `_ =`:** "unudulubmu?" sualına yazı imkanın yoxdur; səbəb kommenti əlavə et.
8. **Defer Close xətaları:** Minimum aşkar iqnor; daha yaxşısı log; propagasiya üçün named err + `closeErr` ayrılması (scan err-i üstə yazma!). Prioritet: əsas xəta > close xətası (log).

---

## Mənbə

- Kitab: *100 Go Mistakes and How to Avoid Them* — Teiva Harsanyi, Manning, 2022 (ISBN 9781617299599)
- Chapter 7: "Error management", book səh. 143–161
- PDF səhifələri: 163–181
