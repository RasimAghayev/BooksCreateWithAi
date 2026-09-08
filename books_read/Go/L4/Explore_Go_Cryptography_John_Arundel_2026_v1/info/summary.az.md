# Explore Go: Cryptography — Xülasə (Azərbaycanca)

**Müəllif:** John Arundel | **Nəşriyyat:** Bitfield Consulting | **İl:** 2026 | **Səviyyə:** L4 (Advanced)

## Kitabın ümumi məqsədi

Bu kitab kriptoqrafiyanı SIFIRDAN, proqramçı üçün qurur: Caesar şifrəsindən
başlayaraq oxucu öz shift cipher paketini TDD ilə yazır, öz Crack alətini
quraraq onu SINDIRIR, açarın vacibliyini (256-bit = kainatdakı atomlar)
kəşf edir, blok-cipher interfeyslərini (cipher.Block/BlockMode) implementasiya
edir, ECB-nin vizual iflasını ÖZ GÖZÜ ilə görür (şəkil konturu ciphertext-də
QALIR!), CBC+IV-yə keçir, hash funksiyalarının niyə mövcud olduğunu
(LenHash/SumHash-in trivial preimage hücumlarından) anlayır, MAC/HMAC/digital
imzaları qaldırır, Diffie-Hellman açar mübadiləsini ədədlərlə izləyir və
finalda AES-GCM (Seal/Open) ilə peşəkar kriptoqrafiyaya çatır. Nəticə:
kriptoqrafiyanın bir dərin intuisiyası + "just use AES" praktiki ağılı.

## Fəsil qruplarının xülasəsi

**Ch1-4 (Cipher-in doğuşu):** codes vs ciphers; behaviour-cümlədən-testə
("Encipher transforms HAL to IBM"); null-implementation bug detektoru sınağı;
table test + subtest; filter CLI (stdin|stdout); Eve frequency analysis
(E→ən çox hərf); brute-force + crib ("GO", PNG magic 0x89); Decipher =
Encipher(-key) refaktoru; ortaq test cədvəli (plaintext/ciphertext).

**Ch5-7 (Açarlar və parollar):** keyspace eksponensialı (hər bayt = ×256);
supernova/qara dəlik fizikası — 256-bit praktik limit; multi-byte açar =
key[i%len(key)]; hex açar flag (DEADBEEF); crib ≥ açar uzunluğu; dictionary
attack (40% hesab), obfuscations, passphrases; parol POLİTİKASI = zərər
(simvol tələbi keyspace-i 2x azaldır!); minimum 14 simvol + paste icazəsi;
K^N; random = prediktabelliyin OLMAMASI (sequences!); bit = məlumat vahidi.

**Ch8-10 (Blok quruluşu):** cipher.Block (Encrypt dst-ə YAZIR, açar
receiver-də); 32-bayt fiks açar — genişləndirmə entropi ALMIR ("nəhəng
hərflərlə parol"); NewCipher + sentinel ErrKeySize + %w wrap + errors.Is;
BlockMode/CryptBlocks; blok döngüsü (re-slice progress); panic = proqramçı
səhvi konvensiyası; PKCS#7 (həmişə pad; N qısa → N×N; son bayt = sayğac).

**Ch11-13 (Sanma, entropiya, random):** Next (carry/overflow — odometer);
benchmark extrapolasiya (4 bayt = 1.6s; 16 bayt = 584 trilyon il); paralel
brute-force = fizikaya toxunur (işık sürəti/qara dəlik); entropiya = log2 N
= bilməmizlik; redundant bitlər (256-bit açar + 250 redundant = 6-bit!);
ciphertext sıxıla BİLMƏZ (Eve testi); Kolmogorov K(s) (π = 32 bit!); PRNG
tələləri (periodicity {5,5}→0,5,5; distribution {2,4}→cüt only);
/dev/random hovzəsi; TRNG; kvant (God Emperor Eve belə 50%+ edə bilməz).

**Ch14-16 (Zəncirlər, hash, blockchain):** ECB = eyni blok → eyni ciphertext
(şəkil GÖRÜNÜR); Mallory (replay/drop/modify — Eve-dən FƏRLİ: AKTİV);
CTR/nonce (gizli DEYİL); CBC = əvvəlki ciphertext zənciri + IV (ilk blok,
açıq göndərilir; yalnız ilk bloku qoruyur); hash: LenHash (preimage trivial)
/ SumHash (avalanche YOX; A→B = ±1) dərsləri; MD5 (1-2 s) / SHA-1 ($100k)
MİLLİ; SHA-256 standart; parol hash: salt (rainbow table MƏHV) + YAVAŞ
(bcrypt/scrypt — gündə 1 hash 3ms, Eve milyonlarla); Bobcoin: distributed
ledger + hash zənciri = sıra doğrulaması; consensus (uzun zəncir qalib);
51% attack; PoW (1/10^64; 120 TWh/il) / PoS ("capitalizm artıq var").

**Ch17-18 (Autentifikasiya, AES):** chosen ciphertext + **Doom Principle**
(əvvəlcə doğrula, sonra decrypt); MAC → length extension → **HMAC**
H(key+H(key+msg)); açar mübadiləsi problemi → **DHM** ((g^a)^b=(g^b)^a;
6^5 mod 13 = 2; 9^5 = 2^4 = 3 — SESSİYA AÇARI; forward secrecy); **RSA**
(faktorlaşdırma one-way; public encrypt/private decrypt; İMZA = hash-i
PRIVATE ilə — "yalnız Alice göndərə bilər"); TLS = asimmetrik handshake →
simmetrik axın; CA/zəmanət zənciri; AES: DES(56-bit NSA) → Triple →
Rijndael (4×4 grid, 14 round, S-box/ShiftRows/MixColumns/XOR; confusion +
diffusion); Go crypto/aes = 45 sətir; **AES-GCM = SEÇİM** (Seal/Open; nonce
prefix + tag suffix; öz padding; CBC+HMAX birlikdə); zəif tətbiq kasaları
(Zoom ECB, IRS zero-IV, Hyundai tutorial-açar); kvant (qubit superpozisiya,
decoherence, Q-Day) + post-quantum (lattice-based).

## Kitabın əsas mesajları

1. **TDD kriptoqrafiyada da işləyir:** behaviour cümləsi → test → sindir →
   düzəlt — hər primitiv (Encipher/Crack/Pad/Next) belə doğuldu.
2. **Uzun açar hər şeydir — AMMA entropi ilə:** 256-bit əzabı proqramla
   generasiya olunursa = həqiqətdə 128-bit; random = crypto/rand.
3. **Alqoritm yox, İSTİFADƏ sınır:** ECB/zero-IV/tutorial-açar — Zoom-dan
   IRS-ə qədər hamısı DOĞRU AES-in YANLIŞ istifadəsidir.
4. **Sıra:** encrypt (CBC/GCM) tək başına çatışmır — integrity (hash) +
   authenticity (MAC/imza) ayrıca qatlar.
5. **"Just use AES":** öz cipher YAZMA; aes.NewCipher + cipher.NewGCM =
   müasir Go kriptoqrafiyasının 90%-i; qalan 10% — açarları düzgün götürmək.

## Kim üçündür

L3+ Go developer-ləri: kriptoqrafiyanın İNTUİSİYASINI istəyənlər üçün ən
yaxşı mövcud materiallardan biri; hər konsept öz əlinlə tikilir, sonra
real-dünyada nə baş verdiyi göstərilir (Zoom/IRS kimi). Praktiki gündəlik
istifadə üçün son fəsil (AES-GCM snippet-ləri) kifayət edir; qalanı =
dərin anlayış. Seriya: The Deeper Love of Go → Know Go → Power of Go: Tools
→ Explore Go (bu kitab).
