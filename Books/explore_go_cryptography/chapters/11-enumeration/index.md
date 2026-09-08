# Chapter 11 — Enumeration (səh. 178-198)

## Bu chapter nədən bəhs edir?

Brute-force (qüvvətli) açar sınama: 32-baytlıq (256-bit) açarları sistemli
şəkildə saymaq — `Next` funksiyası, little-endian inkrement, benchmark ilə
crack sürətinin ölçülməsi və fizikanın (Landauer limiti, qara dəlik) Eve-i
məhdudlaşdırması. TDD: test → benchmark → ekstrapolyasiya → realistik test
açarına keçid.

## Əsas fikirlər

### 1. Crack-in yeni problemi — blok şifrəsi qara qutudur
**Nədir:** Əvvəlki `Crack` versiyası N-ci baytı açarla N-ci açıq-mətn baytına
çevirirdi (bayt-bayt məntiq). Blok şifrəsində (AES tipli) bu **mümkün deyil**:
bütün blok bir vahid kimi şifrlənir, daxili qarışdırma ola bilər.

**Nəticə:** Crack `cipher.Block` interfeysinin içinə "girə bilməz" — yalnız
bütün açar fəzamızını (keyspace) tək-tək sınaya bilər: **brute force**.

### 2. Next — açar generatoru
**Nədir:** Verilmiş açardan "növbəti" açarı istehsal edən funksiya. 256-bit
ədəd `uint64`-ə sığmır → **bayt-bayt inkrement**:

```go
func Next(input []byte) ([]byte, error) {
    // hər baytı sağdan sola yoxla:
    // 255 deyilsə → +1, qaytar
    // 255-dirsə → 0 et, "daşıma" (carry) ilə növbəti bayta keç
    // hamısı 255 → error: açarlar bitib
}
```

**Testlər (TDT):**
```go
tcs := []struct{ input, want []byte }{
    {input: []byte{0, 0, 0},   want: []byte{1, 0, 0}},   // sadə inkrement
    {input: []byte{255, 0, 0}, want: []byte{0, 1, 0}},  // carry
}
```

### 3. Little-endian inkrement (kiçik sonluq)
**Nədir:** `{0,0,0} → {1,0,0}` — **birinci** bayt artır. Bu, adətən gözlənilən
"sonuncu bayt"dan fərqlidir.

**Necə işləyir:** Çoxlu mərtəbəli ədədlərdə ən az əhəmiyyətli bayt **birinci**
yazılır (little-endian / balaca sonluq) — normal "123" (böyük sonluq) tərzdən
fərqli. Daşıma soldakı (növbəti) bayta keçir.

**Nəyə lazımdır:**`Next` universal olur — istənilən uzunluqdakı açar üçün
işləyir (yalnız 32 bayt deyil), sadə və şərh olunandır.

### 4. Yeni Crack — sadə loop
**Necə işləyir:** Bütün-sıfır açardan başla → `Decipher` → crib ilə müqayisə →
uyğunsursa `Next` → təkrarla. `Next` error qaytarana qədər (bütün açarlar
sınıb) və ya crib tapılanadək.

### 5. Benchmark ilə realistlik yoxlaması
**Problem:** Tam test (32 baytlıq açar) `panic: test timed out after 10m` —
çünük 256-bit fəza əlçatmazdır.

**Benchmark mexanizmi:**
```go
func BenchmarkCrack(b *testing.B) {
    for i := 0; i < b.N; i++ { // b.N — mexanizm tərəfindən idarə olunur
        _, err := shift.Crack(ciphertext, crib)
        // ...
    }
}
```
**Nəticələrin oxunması:**
```
BenchmarkCrack-8    48979    24456 ns/op
```
- `-8` → 8 nüvə; `48979` → təkrar sayı; `24456 ns/op` → orta vaxt

**Ekstrapolyasiya:** hər əlavə bayt = 256× vaxt. 24µs (1 bayt) → ~6ms (2 bayt,
253×) → 1.6s (4 bayt, ~256×) — proqnoz real ölçmə ilə təsdiqləndi.
16 bayt artıq çox uzun sürür; 32 bayt = **kainatın yaşından kvadrilion dəfə
çox** → test "asan" açarla (`01 01 01 00...`) yenidən yazılır.

### 6. Enerji limiti — fizika Eve-i dayandırır
**Landauer limiti:** Hər bit silmə əməliyyatı minimum enerji tələb edir.
256-bit açarın brute-force-u o qədər enerji istəyir ki, Eve-in "super
kompüuteri" istiliyə çevrilərək **öz qara dəliyinə çökər** — hesablama bitər,
cavab müşahidə oluna bilməz. Paralel core-lar (8×, minlərlə×) gözlənilən vaxtı
bölmür — əsl limit vaxt deyil, **enerji**dir.

## Əsas terminlər

- Brute-force Attack (qüvvəli hücum)
- Keyspace (açar fazası)
- Little-endian (balaca sonluq)
- Carry (daşıma)
- Benchmark ( performans ölçmə)
- Extrapolation (xarici proqnoz)
- Landauer Limit (enerji limiti)

## Praktik nəticə

- Blok şifrəsini qara qutu kimi yoxla — daxili məntiqə güvənmə
- Böyük ədədlərin sayılması üçün bayt inkrement + carry
- Testlərin realistliyini benchmark ilə yoxla; `ns/op` oxumağı bil
- Uzun sürən testləri zəif açarla saxla; 256-bit açarın sınması fiziki olaraq mümkünsüzdür

## Mənbə

Pages: 178-198 (Chapter 11, Explore Go: Cryptography)
