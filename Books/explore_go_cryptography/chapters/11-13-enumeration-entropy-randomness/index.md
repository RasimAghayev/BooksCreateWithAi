# Chapters 11-13 — Enumeration, Entropy, Randomness (səh. 178-225)

## Bu fəsillər nədən bəhs edir?

Crack-in blok-cipher versiyası (bayt-bayt "içəri baxış" QADAĞA — block = black
box; yalnız tam blokları sına), 32-bayt açarların sanılması (uint64-da sığmır
→ Next funksiyası — bayt-bayt increment, carry, overflow error; "car odometer"),
endianness (big/little — sanma üçün VACİB DEYİL, sadəcə ardıcıllıq), Crack-in
sonsuz döngəsi (Next → NewCipher → Decrypt → müqayisə), benchmarking (b.N,
ns/op; 2-bayt: 24µs/257 = ~90ns/cəhd; 3-bayt: 6ms [×253 təsdiq]; 4-bayt:
1.6s; 16-bayt = 256^12 × 1.6s = 584 TRİLYON il) → test AÇARI SEÇİMİ (testKey
01 01 01 00... = asan; "test suite saniyələr içində"), paralelləşdirmə
analizi (enerji limiti; işıq sürəti bandwith → böyük = yavaş; sıx = qara
dəlik; "linear knife vs exponential gunfight" — GÜCLÜ açarlarda parallel
faydasız, ZƏİF-lərdə lazımsız), entropiya (log2 N bit — N bərabər-ehtimallı
mənalardan hansı; redundant bitlər; "F U CN RD THIS" bilbordu — ingilis
dili REDUNDANTdır; 250 redundant bitli 256-bit açar = 6 bit!), kompression
(lossless = entropi saxlayır; lossy = azaldır — insan gözü/qulağı domeni;
ŞİFRƏMƏTN sıxıla BİLMƏZ — sıxılan = ciphertext deyil!), Kolmogorov kompleksliyi
(K(s) — s-i generasiya edən ƏN QISA proqram; cüt ədədlər/π/Mandelbrot = AŞAĞI
K — "pi" = 32 bit, milyonları təsvir edir; ulduz kütlələri = YÜKSƏK K = random),
key generation (1M bit amma 128-bit proqramla generasiya olunur = 128 bit
entropi; random seçim = yüksək K), chaos (coin toss deterministik amma
praktiki hesablanmaz — başlanğıc şərtlərinə həssas; lava lamp Wall of Entropy
— Cloudflare), poker/rock-paper-scissors (insan beyni random DEYİL — Lecter
sitatı), PRNG-lər (Fibonacci 12 əmrlə — Elite/BBC Micro; pseudo = yalançı;
math/rand eyni ideya), periodicity (5,5 → 0,5,5,0 DÖVRÜ; HƏR deterministik
sequence təkrarlanır — kriptoqrafiya üçün FATAL), distribution/uniformity
(2,4 toxumları → YALNIZ cüt rəqəmlər!; 50/50 gözləntisi), sistem entropiya
hövzəsi (environmental noise — klaviatura/fayl gecikmələri; /dev/random; OS
resurs kimi təqdim edir; crypto/rand buradan), security is relative (bear
metaforası — ayıdan deyil, DİGƏRNDƏN tez qaç; "Evelyn from IT" sosial
mühəndislik; $5 wrench), TRNG (thermal noise — CPU daxili; USB), kvant
random (superpozisiya; ölçmə = nəzəriən ƏVVƏLCƏDƏN BİLİNMƏZ; God Emperor Eve
belə 50%-dən yaxşı edə BİLMƏZ; ANU qrand API — amma şəbəkədən keçən bitlər
INTERCEPT oluna bilər → LOKAL hovzə çox vaxt DAHA təhlükəsizdir).

## Əsas fikirlər

### 1. Block = Black Box (Ch11)
- **Köhnə Crack-in günahı:** N-ci bayt + N-ci açar → N-ci plaintext —
  interfeysi BYPASS edir; blok İÇİNDƏ qarışıq (shuffling) olan cipher-lər
  üçün MƏNASIZDIR
- **Yeni Crack:** yalnız tam blokları Decrypt ilə sına — interfeysə hörmət

### 2. Next — Açar Sanma Alqoritmi (Ch11)
```go
func Next(key []byte) ([]byte, error) {
    for i := range key {
        if key[i] < 255 {
            key[i]++              // bu baytı artır
            return key, nil
        }
        key[i] = 0                // carry: 255→0, növbəti bayta keç
    }
    return nil, errors.New("overflow")   // hamısı 255 = BİTMİŞ
}
```
- **Car odometer metaforası:** bayt dolanda sıfırlanır, 1 növbəti bayta
  DAŞINIR; testlər {0,0,0}→{1,0,0}; {255,0,0}→{0,1,0}; {255,255,0}→{0,0,1};
  overflow {255,255,255}→error
- **Endianness QEYRİ-VACİB:** sıra əhəmiyyətsizdir — YALNIZ hamısının
  çıxması lazımdır (sadəcə Go slice başından sona getmək rahatdır)
- **uint64 YETƏRSMİZ:** 256 bit > 64 bit; math/big VAR amma bayt manipulyasiyası
  daha sadə yol

### 3. Crack-in Yeni Döngəsi (Ch11)
```go
func Crack(ciphertext, crib []byte) (key []byte, err error) {
    plaintext := make([]byte, len(crib))
    key = make([]byte, BlockSize)              // sıfır açarla başla
    for {
        block, err := NewCipher(key)           // hər sınaqda YENİ cipher
        block.Decrypt(plaintext, ciphertext[:len(crib)])
        if bytes.Equal(crib, plaintext) { return key, nil }
        key, err = Next(key)
        if err != nil { return nil, errors.New("no key found") }
    }
}
```
- **Test açarı SEÇİMİ:** 01 01 01 00... — sanma sırasının BAŞINA yaxın;
  TestCrack başlanğıcda 10 DƏQİQƏ timeout verdi (010101... açarı MİLYONLARCA
  il uzaqdır!) → benchmark-et → asan açar seç → 6µs

### 4. Benchmark Metodologiyası (Ch11)
```go
func BenchmarkCrack(b *testing.B) {
    ciphertext := []byte("Uiis message is exactly 32 bytes")  // 01 01 00 açarı
    b.ResetTimer()
    for range b.N { _, _ = shift.Crack(ciphertext, plaintext) }
}
```
- **Extrapolasiya:** hər bayt ×256 → 2B: 24µs; 3B: 6ms (×253 təsdiqləndi);
  4B: 1.6s (TƏXMİNƏN DƏQİQ); 16B: ×256^12 = 584 trilyon il; 32B: kainatın
  yaşından qat-qat çox
- **Çıxarış:** güclü açarlar = brute-force VAXT ITKİSİDİR; Eve üçün daha
  yaxşı: alqoritm zəiflikləri TAPMAQ

### 5. Paralelləşdirmə Fizikası (Ch11)
- **N core:** core-lar arası orta məsafə artır → işıq sürəti limiti →
  kommunikasiya yavaşlayır → BÜYÜK kompüterƏ lazımsız core əlavə = YAVAŞLAMA
- **Sıx kompüter:** kritik sıxlıq → QARA DƏLİK — nəticə ƏLDƏ EDİLMƏZ
  ("Eve cavabı görə bilər — qara dəliyə düşərək; amma kainata qayıda bilməz")
- **Praktika:** paralel brute-force = "linear knife vs exponential gunfight";
  güclü açar → bölmək lazımsız (100 trilyon il / 8 core = hələ də trilyonlar);
  zəif açar → onsuz da sürətli; **NƏTİCƏ: crack-i paralelləşdirməyə dəyməz**

### 6. Entropiya — log2 N (Ch12)
- **Tərif (Schneier):** N bərabər-ehtimallı mənadan hansıının seçildiyini
  yazmaq üçün MİNİMUM bit sayı = log2 N; coin toss sayı kimi də
- **Redundant bitlər:** məzmun ƏLAVƏ ETMİR → effektiv açar uzunluğundan SİL
- **256-bit açar + 250 redundant = 6-bit açar** (görünüş aldadan!)
- **Bilbord dərsi:** "F U CN RD THIS, U CN BCM A SEC & GT A GD JB W HI PA" —
  oxunur! → ingilis dili HƏDDƏN ARTIQ redundant (oxunaqlılıq üçün YAXŞI,
  kripto açarı üçün PİS)

### 7. Kompression Perspektivi (Ch12)
- **Entropiya = lossless sıxmadan SONRA qalan**
- Lossless (run-length və s.): məlumat EYNİ qalır → entropi SAXLANIR
- Lossy (dinamik diapazon, rəng dərinliyi, >20KHz səs): insan qavrayışına
  görə qurban — entropi AZALIR (şəkil/səs üçün OK)
- **ŞİFRƏMƏTN SIRXILA BİLMƏZ:** hər bit ZƏRURİDİR (redundant = çıxarıla
  bilər); **Eve testi:** intercept olunan data-nı sıx — sıxıLMIRSA = ciphertext!

### 8. Kolmogorov Kompleksliyi K(s) (Ch12)
- **Tərif:** s-i generasiya edən ƏN QISA proqramın uzunluğu
- **Aşağı K nümunələri:** cüt ədədlər (984-bit Go proqramı!); π ("pi" = 32
  bit → SONSUZ rəqəm); Mandelbrot (z→z²+c) — GÖRÜNÜŞ aldadır
- **Yüksək K:** qum dənələri diametrləri, ulduz kütlələri — proqramdan
  QISA EDİLMƏZ = random
- **Açar tətbiqi:** 1M-bit açar amma 128-bit proqram yaradır = 128 bit
  entropi → Eve 128-bit proqram fəzasını axtarır; **random seçim = yüksək
  K statistik zəmanəti**

### 9. Chaos və Lava Lampa (Ch12)
- Coin toss = deterministik, amma başlanğıc şərtlərə HƏSSAS (chaos) →
  praktiki hesablanMAZ; "reality has a lot of decimal places"
- **Cloudflare Wall of Entropy:** ~100 lava lamp + kamera → piksel datası =
  entropiya mənbəyi (bəzək DEYİL!)

### 10. PRNG — Elite (Ch13)
```go
// Fibonacci əsaslı — Elite-də 12 MAŞIN ƏMRİ!
func seq() int {
    tmp := seed1 + seed2
    seed1 = seed2
    seed2 = tmp
    return tmp % 10       // mod → artım limitini qır
}
```
- **Seeds:** seçimə görə fərqli sequence; oyun üçün təkrarlanma FAYDALI
  (partikul silinməsi: eyni seedlər = eyni koordinatlar — yaddaş qənaəti)
- **Seeding:** time.Now().Second() (hər iki seed eyni olsa belə əhəmiyyətsiz)
- **math/rand = EYNİ FƏLSƏFƏ:** deterministik + vaxt toxumu

### 11. PRNG-in 2 FATAL ZƏİFLİYİ (Ch13)
1. **Periodicity:** {5,5} → 0,5,5,0 DÖVRÜ; HƏR deterministik sequence
   təkrarlanır → Eve gözləyir, dövrü tapır
2. **Distribution/uniformity:** {2,4} toxumları → YALNIZ CÜT rəqəmlər
   (cüt+cüt=cüt!) — 1000 rəqəmdə 0≈9 gözləntisi POZULUR; güzgü 50/50

### 12. Sistem Entropiya Hövzəsi (Ch13)
- **Environmental noise:** klaviatura intervalları, mouse hərəkəti, disk
  access vaxtları → OS hovzəsi (entropy reservoir) → yüksək keyfiyyətli
  PRNG qidalandırır
- **/dev/random:** `hexdump -n 32 /dev/random` — 32 bayt REAL entropiya;
  **crypto/rand = bu hovzədən** (Go-da təhlükəsiz random-un mənbəyi)
- **Nəticə:** deterministik PRNG + real-noise seed = praktiki təhlükəsiz
  (EVET-in kompüterə girişi YOXDURSA)

### 13. Security is Relative — Bear (Ch13)
- **"You don't need to outrun the bear":** açarlar həqiqətən qırılmaz
  olmalı deyil — YALNIZ düşmənin ƏN YAXŞI SEÇİMİ OLmayacaq qədər çətin
- **Eve-in daha ucuz yolları:** sosial mühəndislik ("Evelyn from IT:
  parolu ADMIN edin"); $5 wrench; meme link ilə kompüter ələ keçirmə; Bobun
  evinə girib DEŞİFRƏ olunmuş mesajları oxumaq
- **Hövzə təhlükəsizliyi:** Eve Alice-in kompüterində syscall hijack edə
  bilərsə — random bitləri oxuyur; amma O ZAMAN onsuz da hər şeyi oxuyur

### 14. TRNG və Kvant (Ch13)
- **TRNG:** thermal noise (istilik interferensiyası) — CPU-lara DAXİLİ;
  USB xarici variantlar
- **Kvant:** superpozisiya (0 VƏ 1 eyni anda); ölçmə = həqiqətən qeyri-
  müəyyən — NƏZƏRİƏN ƏVVƏLCƏDƏN BİLİNMƏZ (gizli switch YOXDUR); God
  Emperor Eve (pan-dimensional, qara dəlik enerjisi) belə 50%-dən yaxşı
  EDƏ BİLMƏZ
- **ANU qrand API:** `rand.New(qrand.NewSource(qrand.NewReader(apiKey)))` —
  internetdən kvant bitləri; **AMMA:** şəbəkə kabelindən oxuna bilər →
  LOKAL /dev/random çox vaxt DAHA yaxşı (intercept imkanları az)

## Əsas terminlər
- Key enumeration — bütün açarların ardıcıl sınanması
- Carry/overflow — bayt dolanda növbətiyə keçid / hamısı dolu
- Big/little-endian — bayt sırası (sanmada QEYRİ-VACİB)
- ns/op — benchmark ortalama vaxt (nanosaniyə/əməliyyat)
- Black box — daxili mexanikası görünməyən komponent
- Entropy (log2 N) — bilməmizlik ölçüsü; minimum bit
- Redundant bit — məzmun daşımayan bit
- Lossless/lossy compression — məlumat qoruyan/azaldan sıxma
- Kolmogorov complexity K(s) — ən qısa generasiya proqramı
- Chaos — başlanğıc şərtlərə həssaslıq (coin toss)
- Wall of Entropy — Cloudflare lava lamp qurğusu
- PRNG (pseudo-random) — deterministik yalançı-random
- Periodicity — sequence təkrarı (PRNG zəifliyi #1)
- Uniformity/distribution — bərabər yayılma (#2)
- Entropy pool — OS-in entropiya anbarı (/dev/random)
- TRNG — thermal-noise hardware generator
- Quantum superposition — 0+1 eyni anda; ölçmə = qeyri-müəyyən
- "Outrun the bear" — mükəmməl deyil, kifayət qədər təhlükəsiz

## Praktik nəticə

1. **Açar sanma idiomu:** Next(key) — increment + carry + overflow error;
   u64/büyük.Int yerinə bayt manipulyasiyası.
2. **Uzun testlərdə açar SEÇ:** benchmark ilə ölç; test açarı = sanma
   sırasının başına yaxın; "test suite 5 saniyədən çox olmasın".
3. **Kripto açarı ÜÇÜN YALNIZ crypto/rand:** math/rand = PRNG (periodic +
   distribution tələləri); /dev/random hovzəsindən gəlir.
4. **Ciphertext sıxılmır:** data sıxıla bilirsə = şifrələnməYİB (Eve testi
   kimi də işlək).
5. **Aşağı K = zəif açar:** π-rəqsəmləri, seq funksiyaları, qısa proqramla
   generasiya olunan hər şey; random (crypto/rand) = yüksək K.
6. **Paralel crack mənasızdır:** enerji sabit; güclü açar → onsuz da
   qeyri-mümkün; zəif → onsuz da sürətli.
7. **"Kifayət qədər" təhlükəsizlik:** düşmənin ən ucuz hücumundan baha
   olmasın kifayətdir — mükəmməllik axtarma.

## Mənbə
Pages: 178-225 (PDF 179-226)
