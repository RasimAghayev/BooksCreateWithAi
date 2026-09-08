# Chapters 11-13 — Enumeration, Entropy, Randomness (səh. 178-225)

## Bu fəsillər nədən bəhs edir?

Keyspace-in brute-force ilə gəzilməsi (Next funksiyası, benchmarking), entropiya
anlayışı (mesajın "həqiqi" informasiya miqdarı — Kolmogorov kompleksliyi,
lossy compression, qısa təsvir = zəif açar), və təsadüfülük (PRNG vs TRNG,
Fibonacci seq, lava lamp "Wall of Entropy", kvant randomness).

## Əsas fikirlər

### 1. Enumeration — 256-bit keyspace gəzintisi (Ch11)

**Crack-in yeni versiyası:** Blok cipher "black box"dir — N-ci baytı ayrıca
sınamaq (əvvəlki Crack-in yanaşması) MÜMKÜNSÜZDÜR. Yeganə yol: bütün açarları
ardıcılla sınamaq.

**Next funksiyası (bayt-bayt increment):**
```go
// {0,0,0} → {1,0,0}; {255,0,0} → {0,1,0} (carry); {255,255,255} → error
func Next(input []byte) ([]byte, error) {
    // hər baytı 0-dan 255-ə qədər artır; overflow → 0 yaz, carry növbətiyə
}
```
**Sub-kod izahı:**
- Little-endian ardıcıllıq: İLK bayt ən az əhəmiyyətli (ƏRƏB ədədi tərsinə!)
  — test gözləntisi: input `[0,0,0]` → want `[1,0,0]`
- Hamısı 255 → "no more keys" error
- `math/big` əvəzinə birbaşa bayt manipulyasiyası — sadə və başa düşülən

**Benchmarking (keyspace sürətini ölçmək):**
```go
func BenchmarkCrack(b *testing.B) {
    for i := 0; i < b.N; i++ { // b.N — benchmark maşını təyin edir
        Crack(ciphertext, crib)
    }
}
// BenchmarkCrack-8   48979   24456 ns/op
```
- 1 baytlıq açar: ~24 miko/saniyə; 2 bayt: ×256 (6.2ms); 3 bayt: ×256 (1.59s)
- 4 bayt = ~1.6 milyard ns ≈ 1.5 saniyə. 16 bayt = artıq çox uzun.
- **32 bayt (256-bit):** 256^16 dəfə daha çox = kainatın yaşı × trilyon
  quadrillion — test üçün "asan" açar (0101...) seçilir

**Fizika limiti (müəllifin əyləncəli dərs):** Eve hesablama gücünü artırsa da,
enerji xərci eynidir; çox böyük kompüuter öz kütləvi ilə **qara dəlik**
yaradır — cavabı görə bilməz. Strong keys (256-bit+) crack-lənməz: "hər yerli
kompüteri ələ keçirsə belə, 100 trilyon il → heç bir sirrin o qədər dəyəri
yoxdur."

### 2. Entropy — informasiyanın həqiqi miqdarı (Ch12)

**Nədir:** N mümkün mesaj varsa, seçimin entropiyası = onu unikal ifadə etmək
üçün lazım olan bit sayı (= qəpik atış sayı). Mesajın bit ölçüsündən FƏRQLİ
olabilir — redundancy (artıqlıq) ola bilər.

**Nümunələr:**
- Kamera freymi: 100MiB raw, amma görüntüdə qonşu piksellər oxşar → entropiya
  az; lossy compression (JPEG/MP3) məhz bunu istismar edir — insan gözü/qulağı
  fərqi görmür
- **Kolmogorov kompleksliyi (K):** sistem = ən qısa təsviri. "1 milyondan kiçik
  bütün cüt ədədlər" = 984 bit proqram (milyonluq siyahı əvəzinə!). π rəqəmləri
  təsadüfi GÖRÜNÜR, amma "pi" sözü kifayətdir → aşağı K = zəif açar
- Mandelbrot şəkli: mürəkkəb görünür, `z → z² + c` → aşağı K

**Açar üçün dərs:** görünüşdə mürəkkəb, amma qısa qayda ilə yaradıla bilən
açar (seq, formula) Eve üçün HƏDİYƏDİR. Əksinə: **random seçilmiş bitlər =
yüksək entropiya** (statistik olaraq asan təsviri yoxdur).

**Wall of Entropy:** Cloudflare-in SF ofisində ~100 lava lamp + kamera — piksel
datası entropiya mənbəyi. İnsan beyni təsadüfülük üçün YARARSIZDIR (pattern
düşür) — Dağıstan şahmatı misalı: proqnozlaşdırıla bilməmək üstünlükdür.

### 3. Randomness — PRNG vs TRNG (Ch13)

**Sadə PRNG (Fibonacci-vari):** iki seed ədədi → cəm → mod 10:
```go
seq: 1, 2, 3, 5, 8, 3, 1, 4...   // 2+3=5, 3+5=8, 5+8=13 mod 10=3...
```
- **Deterministik:** eyni seed → eyni silsilə (Elite oyunu seed-ləri
  yenidən işlədərək qalaktikanı "bərpa edir" — yaddaşa qənaət)
- **Statistik zəiflik:** seed 2 və 4 → həmişə cüt ədədlər; periodik təkrar
- `math/rand` belədir: vaxt seed-li deterministik → kriptoqrafiyaya YARAMAZ

**Təhlükəsizlik nisbidir:** "secure/not secure" yox — Eve-in resursları
karşısında dəyər. $5 wrench hücumu (fiziki zor), phishing, "ADMIN" parol
məcburiyyəti — kripto güclü olsa da, insan zəif həlqədir.

**Environmental noise:** real dünya hadisələri (klaviatura vaxtı, disk
latency, mouse hərəkəti) + deterministik "yuma" (mixing) → OS entropy pool
(`/dev/random`). CPU-lara daxili **hardware RNG (TRNG)** quraşdırılır; USB
qurğu belə var.

**Kvant randomness:** fotonun polyarizasiyası — superpozisiya ölçüləndə 0 və ya
1; **prinsipcə** əvvəlcədən bilinə bilməz (gizli dəyişən YOXDUR). ANU QRN
servisi internetdən kvant bitləri verir. Açar yaratmaq asandır:
```go
key := make([]byte, 32)
rand.Read(key) // crypto/rand — OS entropy pool / TRNG
```
Amma ANU-dan gələn bitlər şəbəkədən keçir — real iş üçün **lokal mənbə**
(entropy pool) tövsiyə olunur.

## Əsas terminlər

- Keyspace (açarlar fəzası)
- Brute-force attack (qüvvə hücumu)
- Entropy (entropiya)
- Kolmogorov Complexity (Kolmogorov kompleksliyi)
- Redundancy (artıqlıq)
- PRNG (psevdotəsadüfi generator)
- TRNG (həqiqi təsadüfi generator)
- Quantum Superposition (kvant superpozisiyası)

## Praktik nəticə

- Brute-force sürəti benchmark ilə ölçülür; hər +1 bayt açar = ×256 vaxt
- Açar üçün: random bitlər (crypto/rand); "qısa qayda" açarları ZƏİF
- math/rand oyunlara, crypto/rand açarlara
- İnsan tərəfindən "random" seçim proqnozlaşdırıla bilər

## Mənbə

Pages: 178-225 (Chapters 11-13, Explore Go: Cryptography)
