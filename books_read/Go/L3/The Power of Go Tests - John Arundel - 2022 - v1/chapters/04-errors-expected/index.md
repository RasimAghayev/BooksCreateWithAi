# Chapter 4 — Errors expected (Gözlənilən Xətalar)

## Bu fəsil nədən bəhs edir?

Testlərdə error nəticələrinin yoxlanması (blank identifier tələsi, t.Fatal
zərurəti), "want error"/"want no error" test cütləri, error simulyasiyası
(errReader, iotest paketi), error-ların == ilə MÜQAYİSƏ OLUNMAMASI (pointer
semantikası), string matching-in kövrəkliyi, sentinel error-lar, errors.Is + %w
wrapping, custom error tipləri + errors.As və 4 əsas nəticə.

## Əsas fikirlər

### 1. Blank Identifier Tələsi
**PIS test:**
```go
got, _ := CreateUser("some valid user")     // error _ İLƏ UDULUB!
if want != got { ... }
```
**2 problem:**
1. GÖZLƏNMƏZ error (funksiya işləməyəndə) testdən YOX — happy-path fixation
2. Testin ÖZÜ səhv olarsa (yanlış input) → "want user ID 1, got 0" — SEBƏB
   GÖRÜNMÜR: CreateUser (0, error) qaytarmış, yəni got ETİBARSIZDIR amma biz
   onu müqayisə edirik!

Go-da exception YOX — error yoxlanışı CALLER-in (yəni TESTin) borcudur. Test
kodunu "defensive" yaz: həmişə error yoxla — səhv olsa belə failure İPUCU verəcək.
**Mövcud testləri təkmilləşdirmə:** errcheck kimi statik analizlə `_`-li yerləri
tap → err = təyin et → yoxla → t.Fatal.

### 2. Gözlənməz Error = t.Fatal (Error YOX!)
```go
// YANLIŞ:
s, err := store.Open("testdata/store.bin")
if err != nil {
    t.Error(err)             // DAVAM EDIR...
}
for _, v := range s.Data {    // s == nil → PANIC!
```
**Səbəb:** error halında konvensiya = nil pointer qaytar; nil-i dereference →
KONFUZ panic (test fail olur, amma YANLIŞ SƏBƏBDƏN — problem Open-da YOX, testdədir).

```go
// DÜZGÜN:
s, err := store.Open("testdata/store.bin")
if err != nil {
    t.Fatalf("unexpected error opening test store: %v", err)   // DƏRHAL DAYAN
}
// indi s-i təhlükəsiz istifadə et
```
**Fixture xətaları da Fatal:** golden file adı səhv ("godlen.txt") →
t.Fatalf("unexpected error reading golden file: %v", err).

### 3. Error Davranışı = API-in Hissəsi
**Fikir:** xəta mesajları API-in ƏN VACİB hissəsidir — işləyəndə heç kim
düşünmür, İŞLƏMƏYƏNDƏ istifadəçi köməyə ehtiyacı var: nə işləmədi, niyə, nə
etməli.

**Minimum 2 test:**
```go
// 1) "WANT ERROR":
func TestFormatData_ErrorsOnInvalidInput(t *testing.T) {
    t.Parallel()
    _, err := format.Data(invalidInput)      // _ — dəyər MARAQSIZ
    if err == nil {
        t.Error("want error for invalid input")
    }
}

// 2) "WANT NO ERROR" (+ nəticə):
func TestFormatData_IsCorrectForValidInput(t *testing.T) {
    t.Parallel()
    want := validInputFormatted
    got, err := format.Data(validInput)
    if err != nil {
        t.Fatal(err)                         // bail out
    }
    if want != got {
        t.Error(cmp.Diff(want, got))
    }
}
```

### 4. Error Simulyasiyası — Fake Reader
**Problem:** io.Reader-dən gələn xətanı İNDİSƏ ETMƏK lazımdır; strings.Reader
HƏMİŞƏ müvəffəqiyyətlidir.

**Həll — xəta qaytaran öz Reader:**
```go
type errReader struct{}
func (errReader) Read([]byte) (int, error) {
    return 0, io.ErrUnexpectedEOF          // HƏMİŞƏ xəta
}

func TestReadAll_ReturnsAnyReadError(t *testing.T) {
    input := errReader{}
    _, err := reader.ReadAll(input)
    if err == nil {
        t.Error("want error for broken reader, got nil")
    }
}
```
**Standart alternativ:** `iotest.ErrReader` (+ TimeoutReader, HalfReader). Çünki
io.Reader KİÇİK interfeysdir — istənilən davranışı saxlayan fake yazmaq asandır.

### 5. Error == Müqayisəsi NİYƏ İşləmir
```go
want := errors.New("Go home, Go, you're drunk")
got := errors.New("Go home, Go, you're drunk")
if got != want {     // HƏMİŞƏ TRUE — FAIL!
```
**Səbəb:** `errors.New` → `&errorString{text}` — POINTER qaytarır. İki fərli
instance = iki fərli yaddaş ünvanı → pointer-lar YALNIZ eyni obyektə işarə
etsələr bərabərdir. Mesaj eyni olsa belə — FƏRQLİ obyektlər.

**Sual düzgün formada:** "eyni yaddaş yerimi?" YOX — "eyni XƏTƏMİ təmsil edirmi?"

### 6. String Matching — Kövrək Həll
```go
if got.Error() != want.Error() { ... }    // İŞLƏYİR amma BRITTLE
```
**Dave Cheney:** "The Error method exists for humans, not code." — xəta mətni
log/ekran üçündür; müqayisə üçün YOX. Bir vergül əlavəsi → bütün testlər qırılır.

### 7. Sentinel Error — Məhdudiyyətli
```go
var ErrUnopenable = errors.New("can't open store file")

func Open(path string) (*Store, error) {
    f, err := os.Open(path)
    if err != nil {
        return nil, ErrUnopenable      // err ATILIR — məlumat İTİR!
    }
}
```
**Müqayisə İŞLƏYİR** (eyni pointer!): `err != store.ErrUnopenable` — amma
istifadəçi sadəcə "can't open store file" görür: HANSI fayl? NİYƏ? Err-də olan
"context" atıldı.

**Suallar:** (1) Caller həqiqətən fərqli xətalar üçün FƏRQLİ hərəkət edirmi?
(2) Yoxsa sadəcə "xəta var" kifayətdir? — 90% halda `err != nil` YETƏR.

### 8. Sentinel + errors.Is (legitim hal)
**Ratelimit nümunəsi — fərqli hərəkət LAZIM:**
```go
var ErrRateLimit = errors.New("rate limit")

func Request(URL string) error {
    resp, err := http.Get(URL)
    if err != nil { return err }
    defer resp.Body.Close()
    if resp.StatusCode == http.StatusTooManyRequests {
        return ErrRateLimit            // retry MÜMKÜNDÜR — fərqləndirmə VACİB
    }
    return nil
}

// Test — httptest ilə server:
func newRateLimitingServer() *httptest.Server {
    return httptest.NewServer(http.HandlerFunc(
        func(w http.ResponseWriter, r *http.Request) {
            w.WriteHeader(http.StatusTooManyRequests)
        }))
}

func TestRequestReturnsErrRateLimitWhenRatelimited(t *testing.T) {
    ts := newRateLimitingServer()
    defer ts.Close()
    err := req.Request(ts.URL)
    if !errors.Is(err, req.ErrRateLimit) {     // errors.Is!
        t.Errorf("wrong error: %v", err)
    }
}
```
**QAYDA:** sentinel YALNIZ istifadəçinin dəyəri üçün (yoxsa sadəcə test üçün = SMELL).

### 9. %w Wrapping — Sentinel + Dinamik Məlumat
```go
var ErrUserNotFound = errors.New("user not found")

func FindUser(name string) (*User, error) {
    user, ok := userDB[name]
    if !ok {
        return nil, fmt.Errorf("%q: %w", name, ErrUserNotFound)
        //                             ↑ %w = WRAP (%v YOX!)
    }
    return user, nil
}

func TestFindUser_GivesErrUserNotFoundForBogusUser(t *testing.T) {
    _, err := user.FindUser("bogus user")
    if !errors.Is(err, user.ErrUserNotFound) {   // WRAP-i AÇARAQ tapır
        t.Errorf("wrong error: %v", err)
    }
}
```
**%v vs %w:** %v = error-u STRING-ə düzəldir (müqayisə mümkünsüz); %w = wrapped
error — orijinal sentinel "yadında saxlanılır". Çoxqatlı wrap mümkün — errors.Is
HAMISINI açır. Wrapping həmişə TƏHLÜKƏSİZDİR — bilməyən tərəf üçün adi error kimi.

### 10. Custom Error Tipləri + errors.As
```go
type ErrUserNotFound struct {
    User string
}
func (e ErrUserNotFound) Error() string {
    return fmt.Sprintf("user %q not found", e.User)
}

// Köhnə üsul — type assertion:
if _, ok := err.(ErrUserNotFound); ok { }

// Müasir — errors.As:
if errors.As(err, &ErrUserNotFound{}) { }
```

### 11. Is vs As — Fərq
| Funksiya | Sual | Nə üçün |
|---|---|---|
| `errors.Is(err, target)` | "IS it this error?" | SENTINEL DƏYƏR (wrap-daxil də) |
| `errors.As(err, &targetType)` | "Is it the same type AS this?" | TİP yoxlaması |

**Custom tipin ÜSTÜNLÜK YOXDUR:** hər xəta üçün struct + Error metodu = uzun;
wrap = sadəcə %v → %w dəyişmək. Çoxqatlı wrap custom tiplərdə İŞLƏMİR.
→ **Köhnə kodu ANLA, amma YENİ kodda %w + errors.Is istifadə et.**

### 12. 4 Əsas Nəticə (kitabın yekunu)
1. Testlərdə error-ları HƏMİŞƏ yoxla — gözlənilən də, gözlənilməyən də
2. Çox vaxt NƏ olduğu vacib deyil — yalnız nil OLMADIĞI
3. Vacib olanda — errors.Is + sentinel
4. Sentinel + dinamik məlumat — fmt.Errorf + %w (wrapped error)

## Əsas terminlələr
- Blank Identifier Trap — error-un _ ilə udulması; ikiqat problem
- Defensive Testing — hər error yoxlanılır; səhv fail İPUCU verir
- errcheck — _-li error-ları tutan statik analiz aləti
- t.Fatal — gözlənməz xətada MÜTLƏQ (nil dereference paniki)
- Want Error / Want No Error — xəta davranışının 2 testi
- errReader — həmişə xəta qaytaran test Reader-i
- iotest.ErrReader/TimeoutReader/HalfReader — standart xəta readerləri
- errorString — errors.New-in unexported arxa tipi (POINTER!)
- Pointer Equality — == yalnız EYNİ ünvanla true
- Brittle Test — string matching; vergül = qırılan testlər
- Sentinel Error — ixrac olunmuş sabit xəta dəyəri (io.EOF)
- httptest.NewServer — test üçün lokal HTTP server
- errors.Is — dəyər yoxlaması (wrap-daxil)
- %w verb — wrap edən fmt.Errorf forması
- Wrapped Error — sentinel + dinamik kontekst birləşdirməsi
- Custom Error Type — struct + Error() metodu (köhnə üsul)
- errors.As — TİP yoxlaması; "same type AS"
- Is vs As — dəyər müqayisəsi / tip müqayisəsi

## Praktik nətidə

(1) `_` ilə error udatma — hərəkəti dayandır; mövcud testlərdə errcheck ilə
təmizlə. (2) Gözlənməz xəta → t.Fatal (Error DAVAM = nil dereference panic =
konfuz fail). (3) Fixture xətaları (golden file) da Fatal. (4) Xəta = API —
"want error" + "want no error" testləri HƏMİŞƏ cüt. (5) İstənilən xəta
davranışını simulyasiya etmək üçün kiçik fake tiplər (errReader) — iotest-də
hazırları var. (6) error == error MÜQAYİSƏ OLMAZ — pointer. (7) Error() string
müqayisəsi = brittle; mümkünsüz et. (8) Sadəcə "xəta varmı" maraqlıdırsa —
err == nil YOX, err != nil BASTA. (9) Fərqli hərəkət tələb edən hallar (retry
on 429) — sentinel + errors.Is. (10) Kontekst lazımdırsa — %w wrap; errors.Is
çoxqatlı wrap-i də açır. (11) Custom error tiplərini ANLA amma YAZMA — %w daha
sadə və güclü. (12) errors.Is = DƏYƏR, errors.As = TİP — qarışdırma.

## Mənbə
Pages: 95-124 (PDF 106-136)
