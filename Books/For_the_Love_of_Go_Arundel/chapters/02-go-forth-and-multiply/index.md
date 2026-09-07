# Chapter 2 — Go forth and multiply (İrəli və Çoxalt)

## Bu fəsil nədən bəhs edir?

TDD ilə yeni funksiya dizaynı (Multiply): test əvvəlcə — kopyala/dəyiştir
metodu, `undefined: calculator.Multiply` compiler xətasının MƏQİYYƏTLİ olması,
null implementation (`return 0`) — testin BUG TUTABİLMƏSİNİN sübutu, sonra real
implementasiya. Test cases: testCase struct (a, b, want), testCases slice,
`for ... range` loop-u, fail mesajında input-un göstərilməsi. Divide problemi:
6/0 → müəyyən olunmayan nəticə; Go-nun ÇOXdəyərli qaytarma qabiliyyəti;
"something and error" pattern (data + error); error tipi; test behaviours not
functions — "One behaviour, one test" qanunu; TestDivide valid-input testi:
`got, err :=`, `err != nil` → `t.Fatalf` (t.Errorf-dan fərqi — dərhal exit),
ardınca want/got müqayisəsi; maraqlı test case-lər (zero, negative, fractional).

## Əsas fikirlər

### 1. Yeni Funksiya — Testdən Başla
Multiply üçün test yaz (mövcud TestAdd-ı KOPYALA-dəyiştir):
```go
func TestMultiply(t *testing.T) {
    t.Parallel()
    var want float64 = 9
    got := calculator.Multiply(3, 3)
    if want != got {
        t.Errorf("want %f, got %f", want, got)
    }
}
```
Dəyişəcəklər: (1) funksiya adı (eyni paketdə eyni ad OLMAZ); (2) çağırış —
`calculator.Multiply(3, 3)`; (3) want = 9.

İlk `go test` → **compiler xətası:** `undefined: calculator.Multiply` — bu,
TAM İSTƏNİLƏNDİR! Funksiya hələ YOXDUR; xəta implementasiya yazanda YOX OLUR.
TDD addımları: test → compile error → null impl → FAIL → real impl → PASS.

### 2. Null Implementation — Testin Özünün Testi
Sual: TestMultiply Multiply-ın bug-larını TUTA BİLİRMİ? Bilinməyən bir şəkildə
keçə bilər (yanlış want, `==` əvəzinə `!=` kimi səhvlər).

**Sübut üsulu:** BİLƏRƏKDƏN səhv implementation yaz, FAIL OLMALIDIR:
```go
func Multiply(a, b float64) float64 {
    return 0        // null implementation — həmişə yanlış
}
```
```
--- FAIL: TestMultiply
calculator_test.go:31: want 9.000000, got 0.000000
```
Test bu "pretty major bug"-u TUTDU → test işləyir → İNDİ real implementation
yazmaq GÜVƏNLİDİR. Null impl olmadan testin heç nə yoxladığını BİLMƏK MÜMKÜN
DEYİLDİ (Go testləri default olaraq KEÇİR!).

### 3. Table Tests — testCase Struct + Slice + range
Sadə funksiyalar üçün 1 input bəsdir; mürəkkəblərdə ÇOX case lazımdır.
TestAdd1and1, TestAdd2and2... YORUCUDUR. Əvəzinə — TABLE TEST:
```go
func TestAdd(t *testing.T) {
    t.Parallel()
    type testCase struct {
        a, b  float64
        want  float64
    }
    testCases := []testCase{
        {a: 2, b: 2, want: 4},
        {a: 1, b: 1, want: 2},
        {a: 5, b: 0, want: 5},
    }
    for _, tc := range testCases {
        got := calculator.Add(tc.a, tc.b)
        if tc.want != got {
            t.Errorf("Add(%f, %f): want %f, got %f", tc.a, tc.b, tc.want, got)
        }
    }
}
```
- **struct** — strukturlaşmış data tipi: çoxlu sahə, bir qeyd
- **slice** — eyni tipli dəyərlər ardıcıllığı (bütün şeylər toplusu; bəzi
  dillərdə "array" deyilir)
- **`for _, tc := range testCases`** — slice üzrə loop; `_` = index (lazım
  deyil), `tc` = cari element; testCase funksiyaya XASdir (daxildə təyin olunur)
- **Fail mesajı INPUT-u göstərir:** "Add(2, 2): want 4, got 5" — HANSI case
  sındı dərhal görünür

**Case müxtəlifliyi qanunu:** "zero, negative, empty, backwards, forwards,
sideways" — oxşar case-lərin çoxu DƏYƏR vermir; FƏRLİ növlər inamı artırır.
Test eyni zamanda DAVRANIŞ SƏNƏDİDİR.

### 4. Divide Problemi — Invalid Input Mövcuddur
`6 / 0` — "0-la vuraraq 6 verən hansı rəqəm var?" — YOXDUR. Cavab variantları:
0 (unsatisfying), ∞, NaN — hamısı "əsas problemin ətrafından dolanmaq"dır.
**Müqayisə:** Add/Subtract/Multiply üçün invalid input YOX idi; Divide üçün
MÜMKÜNDÜR. Function davranışı input-dan asılı olaraq FƏRLİDİR.

### 5. "Something and Error" Pattern
Go funksiyaları ÇOX dəyər qaytara bilər (və çox vaxt qaytarır):
```go
// Divide: data dəyəri + error göstəricisi
func Divide(a, b float64) (float64, error)
```
- `error` — Go-nun xüsusi tipi
- Standart pattern: birinci nəticə = DATA, ikinci = ERROR; error data-nın
  etibarlılığını bildirir
- Book terminology: **"something and error"**

### 6. "Test Behaviours, Not Functions" — Bir Davranış, Bir Test
Davranış = proqramın MÜXTƏLİF ŞƏRAİTLƏRDƏ nə etdiyi. Divide-in İKİ davranışı:
1. Valid input → hesablanmış cavab + nil error
2. Invalid input → nəticə + non-nil error (problem bildirən)

**Qanun: "One behaviour, one test."** Testləri kiçik, sadə, oxunaqlı saxlayır.
"One function, one test" YANLIŞ hədəfdir — çoxmərtəbəli mürəkkəb testlər
yaradır. Valid/invalid davranışları AYRI testlərə böl.

### 7. TestDivide — Valid Input Testi
TestAdd-ı kopyala → dəyiş: Divide çağırışı, uyğun case-lər, İKİ dəyər qəbulu:
```go
func TestDivide(t *testing.T) {
    t.Parallel()
    type testCase struct {
        a, b  float64
        want  float64
    }
    testCases := []testCase{
        {a: 2, b: 2, want: 1},
        {a: -1, b: -1, want: 1},
        {a: 10, b: 2, want: 5},
    }
    for _, tc := range testCases {
        got, err := calculator.Divide(tc.a, tc.b)
        if err != nil {
            t.Fatalf("want no error for valid input, got %v", err)
        }
        if tc.want != got {
            t.Errorf("Divide(%f, %f): want %f, got %f",
                tc.a, tc.b, tc.want, got)
        }
    }
}
```
`got, err :=` — çoxdəyərli qəbul. Valid-input testi err-in NIL olmasını gözləyir.

### 8. t.Fatalf vs t.Errorf
- **t.Errorf** — testi FAIL edir, amma İCRA DAVAM EDİR (başqa case-lər yoxlanır)
- **t.Fatalf** — testi FAIL edir və DƏRHAL ÇIXIR ("things are so broken that
  there's no point continuing")

**Niyə error halında Fatalf?** err != nil olsa, `got` İSTİFADƏ OLUNMAZ —
etibarsızdır (error baş verib). Onu müqayisə etmək mənasızdır → dərhal dayan.
Sıra: ƏVVƏL err yoxla (Fatalf), SONRA want/got müqayisəsi (Errorf).

`%v` — error-un default formatı. Invalid-input testi (növbəti fəsil) yalnız
error yoxlayır — data İGNORE edilir, çünki etibarsızdır.

## Əsas terminlələr
- TDD Addımları — test → compile error → null impl → fail → impl → pass
- Copy-Paste-Modify — mövcuz testdən yeni test (ad + çağırış + want dəyiş)
- Null Implementation — `return 0`-luq bilərəkdən səhv kod; testin özünü yoxlama
- Table Test — testCase struct + testCases slice + range loop
- testCase struct — {a, b, want} sahələri; funksiya daxilində lokal tip
- Slice — eyni tipli dəyərlər ardıcıllığı (array analoqu)
- for ... range — slice üzrə iterasiya; `_` = index, tc = element
- Invalid Input — Divide(6, 0) kimi müəyyən olunmayan hal
- Multiple Return Values — Go funksiyaları çox dəyər qaytarır
- error Tipi — xəta göstəricisi; nil = xəta yoxdur
- "Something and Error" — data + error qaytarma patterni
- Behaviour vs Function — şəraitə görə fərqli davranış; funksiya deyil
- "One Behaviour, One Test" — hər test TƏK davranış
- t.Fatalf vs t.Errorf — dərhal exit / davam edən fail
- nil — "xəta yoxdur" xüsusi dəyər
- %v — error-un default formatı

## Praktik nəticə
(1) Yeni funksiya = TESTDƏN başla; kopyala-dəyiştir: ad, çağırış, want.
(2) `undefined: X` compiler xətası = TDD-nin normal 1-ci addımı. (3) Null
implementation (`return 0`) yazmadan testin bug tutduğuna İNANMA — bu, testin
özünün testidir. (4) Table test: testCase {a, b, want} + testCases slice +
`for _, tc := range`; fail mesajına INPUT-u sal ("Add(2, 2): ..."). (5) Case
müxtəlifliyi: zero, negative, fractional, backwards — oxşarlar dəyərsizdir;
testlər DAVRANIŞ SƏNƏDİDİR. (6) Invalid input mümkünsə funksiya "something
and error" qaytarsın: `(float64, error)`. (7) Bir funksiyanın müxtəlif
davranışları = ayrı testlər: valid (err == nil yoxla) / invalid (err != nil
yoxla). (8) err != nil → t.Fatalf (data etibarsızdır, davam mənasız); want≠got
→ t.Errorf. (9) Valid testdə sıra: ƏVVƏL err, SONRA data. (10) İki nəticəli
çağırış: `got, err := f(...)`.

## Mənbə
Pages: 28-39 (PDF 29-40)
