# Chapter 6 — Fuzzy thinking (Fuuzz Düşüncəsi)

## Bu fəsil nədən bəhs edir?

Random inputların generasiyası (rand.Intn determinizmi, seed məsələsi),
rand.Perm ilə permutasiya, property-based testing (invariant-lar — Square
non-negative), fuzz testing (Go 1.18+ builtin): Fuzz* funksiyası, fuzz target,
f.Fuzz/f.Add, corpus, -fuzz flag, testdata/fuzz faylları = avtomatik regression
testlər, minimizasiya (27 bayt → 0 bayt), rune vs byte bug nümunəsi (FirstRune),
oracle pattern, explicit argument index (%[1]x), interesting input konsepti, -fuzztime.

## Əsas fikirlər

### 1. Dijkstra — Test Heç Vaxt Tamamlanmır
"Testing shows the presence, not the absence of bugs." — testlər xətaların
YOXLUĞUNU sübut edə BİLMƏZ; yalnız inamı ARTIRIR.

### 2. Random Input + Determinizm
```go
input := rand.Intn(10)          // DEFAULT SEED = 1!
encoded := codec.Encode(input)
got := codec.Decode(encoded)
if input != got { ... }         // roundtrip test
```
**Gözəl balans:** math/rand default seed=1 → ardıcıllıq HƏMİŞƏ EYNİDİR → test
DETERMİNİSTİK (flaky YOX). Random seed istermise: `rand.Seed(time.Now().UnixNano())`
→ hər run-da fərqli, amma FLAKY risk.

**Strategy:** adətən FIXED seed (reproducible); input space-i genişləndirmək
üçün ÇOX dəyər generasiya et.

### 3. rand.Perm — Random Permutasiya
```go
inputs := rand.Perm(100)        // 0-99 HƏR BİRİNİ BİR DƏFƏ, random SIRRA
for _, n := range inputs { ... }
```
Random SEÇİM deyil (təkrarlar olur) — random SIRA. 0-ın İTKİNSİZ əhatəsi +
random sıra = order-bug + tam sərhəd örtüyü.

### 4. Property-Based Testing — İnvariantlar
**Problem:** random input-un DƏQİQ nəticəsini bilmirik (want nə olmalı?).

**Səhv yanaşma:** nəticəni testin ÖZÜ hesabla — eyni alqoritm istifadə etsən
HEÇ NƏ test etmirsən!

**Həll — property (invariant):** Square nə qaytarsa qaytarsın — MƏNFİ
OLMAMALIDIR:
```go
func TestSquareResultIsAlwaysNonNegative(t *testing.T) {
    t.Parallel()
    inputs := rand.Perm(100)
    for _, n := range inputs {
        t.Run(strconv.Itoa(n), func(t *testing.T) {
            got := square.Square(n)
            if got < 0 {
                t.Errorf("Square(%d) is negative: %d", n, got)
            }
        })
    }
}
```
**Property-based = davranışı dəqiq dəyərlərlə YOX, INVARIANT-larla ifadə.**
Random input xəta tapsa → o, sadə example-based test kimi statik korpusa keçir.

### 5. Fuzz Testing — Go 1.18+ Built-in
**Nədir:** random generasiya edilmiş inputlarla avtomatik hücum; "fuzz" =
"randomly generated". Məqsəd — düşÜNÜLMƏYƏN input-ları tapmaq (slice overrun
klasikası: `results[0].Name` — boş nəticədə panic, illər sonra aşkara düşə bilər).

**Kitabdan kod nümunəsi:**
```go
// Deliberately buggy:
func Guess(n int) {
    if n == 21 {
        panic("blackjack!")        // gizli trigger
    }
}

// Fuzz test:
func FuzzGuess(f *testing.F) {        // Fuzz prefiks + *testing.F
    f.Fuzz(func(t *testing.T, input int) {   // FUZZ TARGET
        guess.Guess(input)
    })
}
```

**Fuzz target imzası:** `func(t *testing.T, input T)` — birinci *testing.T,
qalanları — fuzzerin random generasiya edəcəyi İNPUT-lar. Sadə target = yalnız
çağırır → yalnız PANİK və ya TIMEOUT sınayır.

### 6. -fuzz Flag və Çıxış
```bash
go test -fuzz .          # '.' = REGEX — bütün fuzz testlər
go test -fuzz FuzzGuess  # konkret
```
Çıxış (müvəffəqiyyətli tapıntı):
```
warning: starting with empty corpus
--- FAIL: FuzzGuess (0.05s)
    testing.go:1356: panic: blackjack!
Failing input written to testdata/fuzz/FuzzGuess/xxx
    To re-run: go test -run=FuzzGuess/xxx
```
- Fuzzer gizli 21-i SANİYƏLƏRDƏ tapdı
- FAIL İNPUT avtomatik testdata/fuzz/ faylına YAZILDI

### 7. Faylın Məzmunu — Versioned Format
```
go test fuzz v1        ← format identifier (geri-uyğunluq üçün)
int(21)                ← FAIL input
```
**Dərs:** maşın-oxunan fayl generasiya edirsənsə VERSİYA saxla — format
dəyişərsə köhnə oxuyucular sınmaz.

### 8. Fuzz Testlər = Regression Test Generator
**Fuzz olmadan adi go test nə edir:** random generasiya ETMİR — corpus-da
(f.Add + testdata/fuzz faylları) olan HƏR İNPUT-u çağırır. Yəni:
- Fuzzer xəta input tapır → testdata/fuzz faylı yazılır
- `go test` (adi) → HƏMİN input HƏMİŞƏ yoxlanılır
- **"A failing test case blocks deployment"** — false positive (yalan xəta)
  false negative-dən (buraxılan xəta) HƏMİŞƏ yaxşıdır

### 9. f.Add — Corpus / Seed Data
```go
func FuzzFirstRune(f *testing.F) {
    f.Add("Hello")          // TRAINING DATA
    f.Add("world")
    f.Fuzz(func(t *testing.T, s string) {
        got := runes.FirstRune(s)
        want, _ := utf8.DecodeRuneInString(s)    // ORACLE
        if want == utf8.RuneError {
            t.Skip()            // invalid UTF-8 → SKIP (pass/fail DEYİL)
        }
        if want != got {
            t.Errorf("given %q (0x%[1]x): want '%c' (0x%[2]x)", s, want)
            t.Errorf("got '%c' (0x%[1]x)", got)
        }
    })
}
```
**Oracle pattern:** gözlənilən nəticəni BİLĞİLİ etibarlı mənbədən al
(utf8.DecodeRuneInString) — özün hesabla YOX.

**t.Skip:** keçilməli input (invalid rune) — amma skip-dən ƏVVƏL funksiya
çağırılır → PANİK yoxlaması QORUNUR!

### 10. Explicit Argument Index — %[N]verb
```go
t.Errorf("given %q (0x%[1]x)", s)
//          ↑ %q = s;  %[1]x = YENİDƏN 1-ci arqument, hex format
fmt.Printf("%s %[1]s %[1]s %[1]s", "hello")
// hello hello hello hello
```
Eyni dəyər bir neçə formatta lazımdırsa — arqumenti TƏKRARLAMA, indeks ilə işarətlə:
`%q` (quoted) + `%[1]x` (hex) = s-nin hər iki görünüşü.

### 11. Minimizasiya — 27 baytdan 0-a
Fuzzer fail input tapanda onu KİÇİLDIRİR (shrink):
```
fuzz: minimizing 32-byte failing input file
→ given "ʭ" (0xcaad)          # 2 bayta endi!
```
Manuel bug report üçün də eyni qayda: minimal reproducer hazırla — bu,
yalnız tetikleyen hissəni saxlayır.

### 12. FirstRune — 3 Iterativ Bug
**Bug 0 (null):** `return 0` → seed "world" FAIL: `want 'w' got ''`.
**Bug 1 (bayt düşünə):** `return rune(s[0])` → boş string PANIC:
`index out of range [0] with length 0` → fuzzer 0 bayta endirdi →
`if s == "" { return utf8.RuneError }` düzəlişi.
**Bug 2 (multi-byte):** input "ʭ" (0xcaad 2-baytlık rune) → funksiya İLK BAYTI
(0xca = 'Ê') rune saydı → oracle-umuxalif: `want 'ʭ' (0x2ad), got 'Ê' (0xca)`.
BU, rune vs byte fərqini BİLMƏYƏN proqramçının klassik bug-ı — fuzz-suz tapmaq ÇƏTİN.

**Doğru implementasiya:**
```go
func FirstRune(s string) rune {
    for _, r := range s {   // range = RUNE üzrə iterasiya!
        return r           // ilk iterasiyada qaytar
    }
    return utf8.RuneError  // boş string → loop 0 dəfə
}
```

### 13. Fuzz İcra Mexanikası
- 8 CORE = 8 worker paralel generasiya (`now fuzzing with 8 workers`)
- **"interesting" input:** YENİ CODE PATH icra edən input → fuzzer onu
  mutasiya üçün BAZA kimi işlədir (coverage instrumentation — coverage-ilə eyni
  mexanizm, növbəti fəsil)
- Fail TAPILMASA → SONSUZ icra! `-fuzztime=10m` ilə məhdudlaşdır:
  `go test -fuzz . -fuzztime=10m`
- CI-da YOX (testdata dəyişiklikləri commit tələb edir) — LOKAL, manual

### 14. Nə Vaxt Fuzz
- Geniş input space + istifadəçi input-u parse edən istənilən funksiya
- Dəqiq nəticəni yoxlaya bilməsən belə — PANİK və HANG həmişə tutulur
- Property (invariant) yoxlama mümkündür (Square kimi)
- Example-based testlərlə yaxşı örtüldükdən SONRA, kritik sistemlərdə

**Richard Dawkins sitatının tərcüməsi:** fuzz = "doğruluğu DISPROVE etməyə
çalış, nə qədər çalışdığını göstər" — müəyyən müddət fail TAPMAMAQ = inam.

## Əsas terminlələr
- Property-Based Testing — invariant-larla davranış testi
- Invariant — input-dan asılı olmayaraq DƏYİŞMƏYƏN xassə
- Example-Based Testing — dəqiq dəyərlə ənənəvi test
- Deterministic Seed — rand seed=1 → təkrarlanan ardıcıllıq
- rand.Perm — random SIRA, tam örtük
- Fuzz Testing — random inputlarla avtomatik bug ovi (Go 1.18+)
- Fuzz* / *testing.F — fuzz test imzası
- Fuzz Target — f.Fuzz-ə ötürülən, input qəbul edən funksiya
- Corpus — seed + kəşf edilmiş inputlar toplusu
- f.Add — corpus-a training data əlavəsi
- f.Fuzz — target qeydiyyatı
- -fuzz [regex] / -fuzztime — fuzzing rejimi / müddət limiti
- testdata/fuzz/ — fail inputların saxlanma yeri
- Regression Test Generator — fail inputlar statik testə çevrilir
- Oracle — nəticənin etibarlı mənbəyi (utf8.DecodeRuneInString)
- t.Skip — keçilməli input; nə pass nə fail
- Shrink/Minimize — fail input-un minimal forması
- Interesting Input — yeni code path tetikleyen input
- Explicit Argument Index — `%[1]x` — eyni arqumentin təkrar istifadəsi
- Rune vs Byte — multi-byte UTF-8 simvollarda 1 bayt ≠ 1 simvol

## Praktik nətidə

(1) Dijstra: test absence sübut ETMƏZ — amma inamı artırır; fuzz = geniş səhra
sübutu. (2) Random inputlardan determinizm gözlənilirsə default seed saxla
(flaky test = düşmən); rand.Seed yalnız istəyən hallarda. (3) rand.Perm =
tam örtük + random sıra — 0-99 hamısı MƏCBURİ VƏ random düzülüşdə. (4) Random
inputda want bilmirsənsə — PROPERTY yaz (Square mənfidir kimi); nəticəni testin
özü hesablaması = testin ləğvi. (5) Fuzz target: adi test funksiyası KİMİ —
sadə halda yalnız çağırış = panik/hang ovu. (6) f.Add ilə seed ver — boş
corpus xəbərdarlığından yaxşı. (7) Oracle pattern — nəticəni etibarlı standart
funksiyadan soruş. (8) Invalid input → t.Skip — amma skip-dən ƏVVƏL çağır ki,
panic qaçmasın. (9) Fail inputlar testdata/fuzz-a AVTOMATİK yazılır — fuzz
testlər adi go test-də REGRESSION kimi işləyir; faylları commit ET. (10) Fuzzer
fail input-u MİNİMİZƏ edir — 27 bayt → 0; manuel reportlarda da minimal
reproducer hazırla. (11) %[1]x — eyni arqumenti 2 formatta çap (təkrar yazma).
(12) Fuzzing sonsuza qədər icra oluna bilər — -fuzztime ilə limit; CI-da YOX,
lokal manuel. (13) İstifadəçi input parse edən HƏR funksiya fuzz kandidatıdır —
rune/byte bug-larını yalnız fuzz tutur. (14) Fuzzer kod path-ləri izləyir —
"interesting" = yeni path; coverage ilə eyni alət.

## Mənbə
Pages: 155-185 (PDF 167-197)
