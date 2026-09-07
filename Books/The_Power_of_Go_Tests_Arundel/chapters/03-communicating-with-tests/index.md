# Chapter 3 — Communicating with tests (Testlərlə Kommunikasiya)

## Bu fəsil nədən bəhs edir?

Testlər KOMMUNİKASİYA aləti kimi (niyyətin sənədləşdirilməsi), test adlarının
cÜMLƏ olması (ACE: Action-Condition-Expectation), "should" YOX "is/does",
gotestdox aləti (test adlarından sənəd generasiyası), happy+sad path kombinasiyası,
deep vs shallow abstraction, informativ failure mesajları və onların süni fail ilə
sınanması, executable examples (Output comment, Unordered output, ad konvensiyaları).

## Əsas fikirlər

### 1. Testlər Niyə Yazılır — Kommunikasiya
**"Tests are stories we tell the next generation of programmers"** (Osherove).
Test = yoldaş proqramçılarla, XƏLƏFLƏRLƏ və GƏLƏCƏK ÖZÜNƏ kommunikasiya.

**Tests capture INTENT:** sistem nə etməli, hansı şərtlərdə, hansı input-la.
İlk ünsiyyət qarşısı — ÖZÜN: davranışı SÖZLƏ təsvir et ("When the user does X,
then Y happens") — bu, boşluqları doldurur. "No amount of elegant programming
will solve a problem if it is improperly specified" (Bryce).

**İlk sual:** "What are we really testing here?" — cavab sözə/sətərə ifadə
edilməyənə qədər test YAZMA.

### 2. Test Adları CÜMLƏ Olmalıdır
**Fikir:** test FAIL olanda əvvəl çap olunan şey — ADIDIR. Ad = mesaj.

**Kitabdan transformasiya:**
```go
// ZƏİF (test kodunu OXUMALI ki, anla):
func TestValid(t *testing.T) { ... }
// --- FAIL: TestValid
//     valid_test.go:12: want true, got false     ← TAM OPAK!

// GÜCLÜ — cümlə = ad:
func TestValidIsTrueForValidInput(t *testing.T) {
    t.Parallel()
    if !valid.Valid("valid input") {
        t.Error(false)
    }
}
// --- FAIL: TestValidIsTrueForValidInput
//     main_test.go:20: false
```
Adın ÖZÜ konteksti daşıyır: "Valid is true for valid input fail edir → Valid
valid input üçün false qaytarır" — test kodunu OXUMADAN.

### 3. ACE — Adın 3 Komponenti
| Komponent | Valid nümunəsində |
|---|---|
| **A**ction | Valid çağırma |
| **C**ondition | valid input ilə |
| **E**xpectation | true qaytarır |

Ad CƏMLƏ SADƏLƏŞDİRİR: input parametrdirsə → nəticəyə TƏSİR etməlidir; yeni
şərtlər (invalid input?) → YENİ TEST-lər üçün fikir verir.

### 4. Happy + Sad Path Kombinasiyası
```go
// YALNIZ happy path ("return true" hələ keçir!):
func TestValidIsTrueForValidInput(t *testing.T) { ... }

// Sad path ƏLAVƏ ET:
func TestValidIsFalseForInvalidInput(t *testing.T) {
    t.Parallel()
    if valid.Valid("invalid input") {
        t.Error(true)
    }
}
```
**Səhv implementasiyaların yoxlanışı:**
```go
func Valid(input string) bool { return true }          // sad yoxdursa keçir!
func Valid(input string) bool {
    return strings.Contains(input, "valid input")      // hər ikisi varsa keçir!
}
```
**Substring bug:** "valid input" daxildir amma input == "invalid input" → Contains
TRUE — yanlış! Daha bir test: substring ehtiva edən AMMA tam bərabər OLMAYAN input
üçün false. **"Test behaviours, not functions"** — hər DAVRANIŞ ayrıca test.

**Əsas prinsip:** "make the easiest way also the correct one" — testlər düzgün
implementasiyanı ƏN SÜRƏTLİ yol edir (lazy implementer üçün).

### 5. gotestdox — Test Adlarından Sənəd
```bash
go install github.com/bitfield/bitfield/gotestdox/cmd/gotestdox@latest
gotestdox
# valid:
#  ✔ Valid is true for valid input (0.00s)
#  ✔ Valid is false for invalid input (0.00s)
```
Test adlarından `Test` prefiksil silib boşluqları bərpa edir → OXUNARLI SƏNƏD.
Fail → ✘ işarəsi ilə. **Code review ikinci addımı:** tələbə layihəsinə gotestdox
işə sal → adların doluluğu dərhal görünür (Dan North-un BDD-dəki agiledox-dan
ilhamlanıb).

**İstifadə dəyəri:** std/encoding/csv test-ləri: "Read simple / Read bare quotes /
Read huge lines" — hər biri DAVRANIŞ cümləsi; fail = "nə işləmir" dərhal aydın.

### 6. "Does", not "Should"
```go
// YOX: Valid should be true for valid input
// BƏLİ: Valid IS true for valid input
```
**Arqument:** hər test cümləsi görünməz şəkildə "...when the code is correct" ilə
bitir — "should" ASPIRASİYADIR, test isə TƏYİNATDIR (definition). "Should"=
qeyri-müəyyənlik; "is" = spesifikasiya.

### 7. Cümlə = Unit-un Ölçüsü
- 1 cümlədə ifadə olunmursa → ÇOX davranış bir testdə → BÖL
- 1 cümlədə ifadə olunmursa → ÇOX davranış bir UNIT-də → unit-i BÖL (dizayn siqnalı)
- İstifadəçi-görünən davranışa ÇEVİRİLƏ BİLMİYƏN test → LAZIMSIZ, yazma

**Dan North:** "How much to test becomes moot: you can only describe so much
behaviour in a single sentence."

### 8. Deep vs Shallow Abstraction
**Ousterhout (A Philosophy of Software Design):** yaxşı modul = DEEP — sadə
interfeys arxasında GÜCLÜ maşın; shallow = kompleks API amma az fayda →
kağız-işləri.

**Test cümləsi bunu AŞKARLAYIR:** cümlə "bu, az iş görür" hissi verirsə →
REDESIGN siqnalı. API-komplekslik / fayda nisbəti ARTIR.

### 9. Uzun Test Adları — Problem DEYİL
- Bu funksiyalar heç vaxt ÇAĞIRILMIR
- Public API-da YOX
- Yeganə istifadə = FAIL mesajında GÖRÜNMƏ

**→ Qısaltma: mənalı ol, qısa olmaq məcburiyyəti YOX.**

### 10. Informativ Failure Mesajları
**Təkmilləşmə nərdivanı (kitabdan):**
```go
t.Error("fail")                            // ƏN PİS — heç nə demir
t.Error("want not equal to got")           // infer olunur — faydasız
t.Errorf("want %v, but didn't get it", want)  // gözlənilən VAR, alınan YOX
t.Errorf("want %v, got %v", want, got)      // İKİSİ də var — ƏSAS standart
t.Error(cmp.Diff(want, got))               // uzun/mürəkkəb dəyərlər üçün
t.Error("wrong account balance after concurrent deposits", cmp.Diff(want, got))
// ↑ İDEAL: USER-GÖRÜNƏN PROBLEM + dəyərlər
t.Error("input %d caused unexpected error: %v", input, err)
// ↑ problem yaradan İNPUT daxil et
```
**MƏQSƏD:** developer mesajdan TAM MƏLUMAT almalı — test KODUNU OXUMAMALI.

### 11. Failure Mesajlarını SÜNİ FAIL ilə Sına
```go
// Test:
func TestDouble2Returns4(t *testing.T) {
    want := 4
    got := double.Double(2)
    if want != got {
        t.Errorf("Double(2): want %d, got %d", want, got)
    }
}

// Mesajı sınamaq üçün SİSTEM koduna müvəqqəti BUG SAL:
func Double(n int) int {
    return n*2 + 1        // qəsdən
}
// FAIL çıxışını qiymətləndir → BUG-u GERİ QAYTAR
```
**QEYDİYYAT:** want-u dəyişmə (5-ə) — SƏHV yoldur: want-un DÜZGÜNLÜYÜ də
yoxlanışın bir hissəsidir. Bug FUNKSIYAYA salınır, testa YOX.

### 12. Executable Examples — Sənəd + Test
```go
func ExampleDouble() {
    fmt.Println(double.Double(2))
    // Output:
    // 4
}
```
**Xüsusiyyətlər:**
- `Example` prefiks; parametr/qaytarma YOX
- `// Output:` comment = gözlənilən stdout — TEST MEXANİZMİ YOXLAYIR
- pkg.go.dev sənədinə AVTOMATİK daxil olur + brauzerdə RUN düyməsi
- go test zamanı da İCRA OLUNUR — Output != actual → FAIL

```go
// --- FAIL: ExampleDouble
// got: 4
// want: 5
```

### 13. Example Ad Konvensiyaları
| Ad | Sənəddə yeri |
|---|---|
| `Example()` | bütün paket |
| `ExampleDouble()` | Double funksiyası |
| `ExampleDouble_with2()` | 2-ci variant (suffix kiçik hərflə) |
| `ExampleUser()` | User tipi |
| `ExampleUser_NameString()` | User.NameString metodu |

```go
func ExampleUser_NameString() {
    u := user.User{Name: "Gopher"}
    fmt.Println(u.NameString())
    // Output:
    // Gopher
}
```

### 14. Unordered Output + Playground
```go
for k, v := range m {
    fmt.Println(k, v)
}
// Unordered output:        ← sıra ƏHƏMİYYƏTSİZ (map iterasiyası random!)
// 1 true
// 2 false
// 3 true
```
**Playground məhdudiyyətləri:** şəbəkə YOX; CPU/time limiti; time.Now həmişə
"2009-11-10 23:00:00 UTC" (deterministik). Lokal go test-də məhdudiyyət YOX.

**Dəyər:** sənəddəki kod NUNMAYIB — compile+run OLUNUR, köhnələ bilməz. Std
kitabxana geniş istifadə edir; 3-cü tərəf paketlərə PROFESSIONAL görüntü verir.

## Əsas terminlələr
- Tests capture intent — niyyətin kodlaşdırılmış sənədlənməsi
- "What are we really testing here?" — hər testin ilk sualı
- Test adı = davranış CÜMLƏSİ
- ACE (Action-Condition-Expectation) — test adının 3 komponenti
- Happy/Sad Path — müsbət/mənfilərin İKİSİ test olunmalı
- Test behaviours, not functions — davranış vahidi, funksiya YOX
- Substring Bug — Contains-in == yerinə istifadəsi klassikası
- gotestdox — test adlarından oxunaqlı sənəd generatoru
- BDD/agiledox — cümlə-adlar konsepsiyasının mənşəyi (Dan North)
- "Is", not "should" — təyinat dilində spesifikasiya
- Deep/Shallow Abstraction — Ousterhout; interfeys/funksional nisbəti
- Failure Message — testin ƏN vacib çıxışı; istifadəçi-görünən problem
- Deliberate Bug — mesajı yoxlamaq üçün funksiyaya müvəqqəti xəta
- Executable Example — Example* + Output comment; sənəd+test
- `// Unordered output:` — sıra-müstəqil yoxlama
- pkg.go.dev / Go Playground — brauzer icra mühiti
- Zero Maintenance — sənəd kodu heç vaxt köhnəlmür

## Praktik nətidə

(1) Test yazmadan əvvəl "What are we really testing here?" — cavabı bir CÜMLƏ
olmalı. (2) Cümləni boşluqsuz funksiya adına çevir — TestValidIsTrueForValidInput.
(3) ACE-yə uyğunlaş: Action-Condition-Expectation eksiksiz olsun. (4) Hər happy
path üçün sad path DÜŞÜN — "hansı səhv kodlar hələ keçir?" sualını ver. (5) Adlarda
"should/must" YOX — "is/does" (test = təyinat, arzu YOX). (6) 1 cümləyə
sığmırsa → testi BÖL; hələ sadədirsə → unit-i BÖL (shallow abstraction siqnalı).
(7) Uzun ad = normal — fail mesajının YARISI addır. (8) Failure mesajı:
dəyərlər + İSTİFADƏÇİ-GÖRÜNƏN problem; "fail" yazmaq = ən pis hal. (9) Mesajları
deliberate bug ilə SINA — want-u dəyişmə, FUNKSIYAYA bug sal. (10) Hər public
funksiyaya Example* yaz — sənəd köhnəlmir, pkg.go.dev-də görünür, test kimi icra
olunur. (11) Random sıralı çıxış üçün Unordered output. (12) gotestdox-u
code-review-da işə sal — dolu/boş adlar dərhal üzə çıxır.

## Mənbə
Pages: 61-94 (PDF 72-105)
