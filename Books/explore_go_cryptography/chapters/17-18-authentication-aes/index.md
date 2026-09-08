# Chapters 17-18 — Authentication, AES (səh. 278-329)

## Bu fəsillər nədən bəhs edir?

Mesaj bütövlüyü → autentifikasiya pilləsi: hash açıq göndərilirsə brute-force
preimage MÜMKÜNDÜR; hash+mesaj birlikdə şifrələnirsə — **chosen ciphertext
attack** (Mallory Bob-u istədiyi ciphertext-ləri decrypt etməyə MƏCBUR edir →
açar haqqında məlumat; **Cryptographic Doom Principle** — Bob heç vaxt
YANLIŞ hash-li mesajı decrypt ETMƏMƏLİ → hash zərfdən XARİCDƏ oxunmalı);
MAC (hash(ciphertext + key) — yalnız açarı olan doğrulaya bilər), length
extension attack (SHA-256 ailəsi: Mallory mesaja MƏTN ƏLAVƏ edib YENİ düzgün
hash hesablaya bilir — "…goodbye forever"), **HMAC** (H(key + H(key + message))
— inner hash gizli → extension qarşısı; autentiklik GİZLİLİKSİZ də mümkün),
açar mübadiləsi problemi (email vasitəsilə açar → Eve oxuyur/Mallory əvəz
edir; key splitting — yarısı email, yarısı telefon), **asymmetric/public-key**
(şəxsi hissə GİZLİ + açıq hissə İSTƏDİYİN YERƏ), Diffie-Hellman-Merkle
(exponentiation ASAN vs logaritm ÇƏTİN — one-way; (g^a)^b = (g^b)^a; p=13,
g=6 nümunəsi: Alice a=5 → 2; Bob b=4 → 9; 9^5 mod 13 = 2^4 mod 13 = 3 =
SESSİYA AÇARI; forward secrecy — hər söhbətə YENİ açar), RSA (Cocks →
Rivest/Shamir/Adleman; 2 böyük sadə → public/private cütü; faktorlaşdırma
ÇƏTİN; çoxlu alıcılara encrypt; öz açarına da), **imza** (hash-i PRIVATE key-lə
encrypt → hər kəs PUBLIC ilə doğrulayır; yalnız Alice yarada bilər),
autentifikasiya = bütövlük (HMAC) + imza (RSA); TLS (handshake-da asimmetrik
→ sessiya açarı → simmetrik), **chain of trust** (Trent — CA; kök CA-dan
sənədə qədər zəncir; pizza date problemi həll olunmur → CA-lar), Ch18: AES
tarixi (DES 1973 — NSA 56-bit MƏHDUDİYYƏTİ; Triple DES 80-ci illər; 90-ların
sonu bir neçə günə düşür → müsabiqə), **Rijndael** (Rijmen+Daemen; blok 128
bit = 4×4 grid, column-major; açar 128/192/256 — AES-256; Ferrari-birinci-
sürət metaforası), 14 round (hər round üçün ayrıca round key — açardan
derived), round resepti (1: AddRoundKey; 2-13: SubBytes [S-box 256 girişli
lookup] + ShiftRows [sıraları 1/2/3 sola] + MixColumns [matris] + AddRoundKey;
14: SubBytes+ShiftRows+AddRoundKey), confusion (DƏYİŞMƏ) vs diffusion
(YER DƏYİŞMƏ — rail fence nümunəsi), Go crypto/aes (encryptBlockGo —
uint32 sütunları; ^ XOR; nr := len(xk)/4 - 2; te0-te3 cədvəlləri BÜTÜN
middleri birləşdirir; 45 sətir!), AES-CBC (shift → aes KEÇİDİ: aes.NewCipher
— interfeys eynidir!; Pad/Unpad köçürülür), **AES-GCM** (Seal/Open; nonce
prefix + auth tag suffix; öz padding-i VAR — Pad LAZIM DEYİL; CBC+HMAC
birlikdə), zəif tətbiqlər (Zoom ECB; IRS all-zero IV; Hyundai tutorial
nümunəsi açar; Philips Hue side-channel — şəhər lampalar botnet), kvant
(qubit superpozisiya → quantum parallelism — BÜTÜN açarları eyni anda;
RSA/DHM faktorlaşdırma asanlaşar; DECOHERENCE — stray photon ölçür =
superpozisiya çökür; maye helium; Q-Day), post-quantum (daha böyük açarlar;
lattice-based crypto; quantum-resistant hashing).

## Əsas fikirlər

### 1. Chosen Ciphertext + Doom Principle (Ch17)
- **Sıra 1:** hash(mesaj) açıq → Mallory eyni-hash saxta mesaj brute-force
- **Sıra 2:** mesaj + hash birlikdə ŞİFRƏLƏ → Bob bilmediyi ciphertext-ləri
  decrypt etməli olur → chosen ciphertext attack: açar haqqında MƏLUMAT
  SIZIR
- **Cryptographic Doom Principle:** Bob YALNIZ düzgün-hash-li mesajı decrypt
  etsin → hash ŞİFRƏSİZ oxunmalı (zərfün XARİSİNDƏ)

### 2. MAC → HMAC (Ch17)
```go
// MAC: hash(ciphertext + key) — doğrulama açar tələb edir
// ZƏİFLİK (length extension): SHA-256 ailəsində Mallory mövcud hash-dən
// YENİ MƏTNƏ düzgün hash çıxarır ("...goodbye forever" əlavəsi!)
// HMAC: H(key + H(key + message))
//   inner hash GİZLİ → extension mümkünsüz
//   outer açar tələb edir → preimage brute-force-yə qapalı
```
- **Autentiklik gizlilikdən ASILDIR:** açıq mesaj + HMAC — yalnız Alice
  yarada, yalnız Bob doğrulaya bilər

### 3. Açar Mübadiləsi → DHM (Ch17)
- **Dilemma:** təhlükəsiz rabitə üçün açar lazım; açar üçün təhlükəsiz
  rabitə → **açarı PAYLAŞMA**
- **Key splitting:** hissə-hissə fərqli kanallarla (email+telefon) — kömək,
  həll YOX
- **DHM one-way funksiyası:** exp ASAN (2^10=1024) vs log ÇƏTİN (cəhdlə-
  səhvə — brute-force kimi)
```
p=13, g=6 (public)    Alice: a=5 (gizli)    Bob: b=4 (gizli)
Alice → Bob: 6^5 mod 13 = 2      Bob → Alice: 6^4 mod 13 = 9
Alice: 9^5 mod 13 = 3             Bob: 2^4 mod 13 = 3   ← EYNİ sessiya açarı!
Sirr: (g^a)^b = (g^b)^a
```
- **Forward secrecy:** hər SESSİYAYA yeni açar — bir açar sınırsa yalnız O
  mesaj itir
- DHM tam public-key DEYİL — simmetrik açar mübadiləsinin TƏHLÜKƏSİZ yolu

### 4. RSA + İmza (Ch17)
- **Açar cütü:** 2 böyük sadə rəqəm → public (açıq) + private (gizli);
  vahidlik: faktorlaşdırma ÇƏTİN (one-way); çoxlu alıcıya (özünə daxil)
  encrypt
- **İmza = ƏKS istiqamət:** hash-i PRIVATE ilə encrypt → hər kəs PUBLIC
  ilə doğrulayar: "yalnız Alice yarada bilər"
- **İki məqsəd:** (1) yalnız SƏN oxuya bilərsən; (2) yalnız SƏN göndərə
  bilərsən
- **Autentifikasiya = bütövlük + imza:** HMAC (simmetrik, sürətli) +
  RSA (asimmetrik, YAVAŞ) → praktikada: RSA/DHM yalnız SESSİYA AÇARI üçün
  (TLS handshake → sonra simmetrik)

### 5. Chain of Trust (Ch17)
- **Problem:** "Bob-un public açarı həqiqətən Bob-unkudurmu?" — Eve
  öz açarını Bob adına qeyd edə bilər
- **Trent (CA):** hər ikisinə ETİBAR edilən 3-cü tərəf; Bob-un sertifikatını
  İMZALAYIR; CA-lara DELEGASİYA; **kök CA-ya qədər zəncir** — bir kəs
  doğruladıqda hamısı doğrudur
- TLS: bankofbob.com-ə kim zəmanət verir? — zəncir

### 6. AES Tarixi (Ch18)
| Mərhələ | Açar | Tale |
|---|---|---|
| DES (1973) | 56 bit | NSA məhdudiyyəti — "biz SINA BİLƏK, digərləri YOX" |
| Triple DES (80-ci) | 3×56 | can aşağı vurma — "kicking the can" |
| AES (1997 müsabiqə) | 128/192/256 | Rijndael qalib |
- "Attacks always get better; they never get worse" — 90-cı illərdə DES
  bir neçə GÜNƏ ($200k maşın)

### 7. AES Daxili Struktur (Ch18)
- **Grid:** 16 bayt blok = 4×4, COLUMN-major (yuxarıdan aşağı doldur)
- **Round keys:** 256-bit açardan 14× 128-bit round açarı DERIVED (key
  schedule)
- **Resept:**
  - Round 1: AddRoundKey (XOR — ^ operatoru; toplama KİMİ sadə/reversible/
    hardware-dostu)
  - Round 2-13: **SubBytes** (S-box — 256 girişli DƏYİŞDİRİCİ cədvəl;
    confusion) → **ShiftRows** (sıralar 0/1/2/3 sola — diffusion; hər çıxış
    sütunu 4 GİRİŞ sütunundan) → **MixColumns** (sütun-matris; diffusion) →
    AddRoundKey
  - Round 14: SubBytes + ShiftRows + AddRoundKey (MixColumns YOX)
- **Confusion vs Diffusion:** dəyər dəyişmə vs YER dəyişmə — shift cipher
  YALNIZ confusion; rail fence YALNIZ diffusion (zəif); AES HƏR İKİSİ × 14

### 8. Go crypto/aes İmplementasiyası (Ch18)
```go
func encryptBlockGo(xk []uint32, dst, src []byte) {
    _ = src[15]                       // EARLY bounds check — dərhal panic
    s0 := binary.BigEndian.Uint32(src[0:4])   // 4 sütun = 4 uint32
    // ...
    s0 ^= xk[0]  // Round 1: XOR
    nr := len(xk)/4 - 2               // round sayı = açar ölçüsündən
    for r := 0; r < nr; r++ {         // midd rounds: te0-te3 cədvəlləri
        t0 = xk[k] ^ te0[uint8(s0>>24)] ^ te1[...] ^ ...  // HƏR ŞEY BİRLİKDƏ
    }
    // son round: sbox0 + XOR; dst-yə yaz
}
```
- **45 sətir + sabit cədvəllər**; te cədvəlləri SubBytes+ShiftRows+MixColumns
  BİRLİŞDİRİR (sürət üçün, oxunaqlıq QURBAN); hardware AES instructions
  (assembly) əsl seçim

### 9. AES-CBC — Tam Keçid (Ch18)
```go
// shift.NewCipher → aes.NewCipher — HEÇ NƏ DƏYİŞMİR (interfeys!)
block, err := aes.NewCipher(key)               // crypto/aes
iv := make([]byte, aes.BlockSize)              // 16 (32 YOX!)
rand.Read(iv)
os.Stdout.Write(iv)
enc := cipher.NewCBCEncrypter(block, iv)
plaintext = Pad(plaintext, aes.BlockSize)       // Pad/Unpad köçürüldü
enc.CryptBlocks(ciphertext, plaintext)
```
- **Plug-and-play dərsi:** cipher.Block interfeysi sayəsində BEYUTƏ ŞİFT →
  DÜNYA STANDARTI 2 sətirlə

### 10. AES-GCM — Son Həll (Ch18)
```go
block, _ := aes.NewCipher(key)
gcm, _ := cipher.NewGCM(block)
nonce := make([]byte, gcm.NonceSize())
rand.Read(nonce)
ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)   // prefix = nonce!
os.Stdout.Write(ciphertext)

// decipher:
nonce := ciphertext[:gcm.NonceSize()]
ciphertext = ciphertext[gcm.NonceSize():]
plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)  // tag yoxlanır!
```
- **Wire format:** nonce (prefix) + ciphertext + auth tag (suffix — Seal
  avtomatik qoşur)
- **GCM = CBC + HMAC bir əməliyyatda:** şifrələmə + bütövlük; **öz padding-i
  var** — Pad/Unpad SİLİNİR; tampering → Open ERROR qaytarır
- **"Go-da müasir kriptoqrafiya üçün lazım olan BÜTÜN kod budur"**

### 11. Zəif Tətbiq İstatistikası (Ch18)
| Sistem | Səhv |
|---|---|
| Zoom | AES-ECB mode! |
| IRS | ECB + all-zero IV ("password" parolu kimidir) |
| Hyundai | QLOBAL açar = internet tutorial nümunəsi |
| Philips Hue | Side-channel ilə signing key → şəhər botnet |
- **Moral:** AES təhlükəsizdir, İSTİFADƏ YANLIŞDIR; "blame the
  cryptographically-illiterate managers"

### 12. Kuant Gələcəyi (Ch18)
- **Qubit:** 0+1 superpozisiya; N qubit = 2^N eyni anda (quantum
  parallelism — 256 qubit = bütün açarları BİRDƏN sına)
- **RSA/DHM ölümcül:** faktorlaşdırma/logaritm kvantda ASAN (Shor tipli);
  AES-256 daha davamlı (Grover ~yarı-atlamalar)
- **Catch — decoherence:** stray photon = ölçmə = superpozisiya ÇÖKÜR;
  izolyasiya = maye helium, böyük/donmuş/baha; bugünkü qubitlər klassikdən
  YAVAŞ
- **Q-Day:** AES/RSA-nın iflas günü — "on illər, ya da sabah"; **post-
  quantum:** böyük açarlar + lattice-based crypto + quantum-resistant
  hashing (ARAŞDIRMA ARTIQ GEDİR — simulyasiya ilə sına)

## Əsas terminlər
- Chosen ciphertext attack — Bob-u seçilmiş ciphertext decrypt etməyə
  məcbur etmə; Doom Principle: əvvəlcə doğrula
- MAC (Message Authentication Code) — hash + açar
- Length extension attack — mövcud hash-dən uzadılmış mesaja keçid
- HMAC — H(key + H(key + message)); ikili-hash quruluşu
- Key exchange problem — açarı təhlükəsiz ötürmə paradoxu
- Key splitting — açarı kanallara bölüb göndərmə
- Asymmetric/public-key — public (açık) + private (gizli) cüt
- DHM — modular exponentiation one-way; ortaq sessiya açarı
- Forward secrecy — sessiya-başına açar; keçmiş qorunur
- RSA — sadə-cüt faktorlaşdırma əsaslı public-key; imza
- TLS handshake — asimmetrik → simmetrik keçid
- Chain of trust / CA — Trent; kök sertifikatdan zəncir
- DES/Triple DES — 56-bit keçmiş standartlar
- Rijndael/AES — müsabiqə qalibi; 4×4 grid; 128-bit blok
- S-box — SubBytes əvəzetmə cədvəli (confusion)
- ShiftRows/MixColumns — sətir/sütun diffusion
- AddRoundKey — XOR ilə round açarı
- GCM (Galois Counter Mode) — encrypt + authenticate birlikdə; Seal/Open
- Qubit/superpozisiya/decoherence — kvant bit; ikili vəziyyət; çökmə
- Q-Day / post-quantum — kvant iflas günü / müqavimətli kriptoqrafiya

## Praktik nəticə

1. **"Just use AES":** öz cipher İCAT ETMƏ; öz AES implementasiyası YAZMA;
   aes.NewCipher + cipher.NewGCM = Go-da müasir kriptoqrafiyanin 90%-i.
2. **AES-GCM standart seçim:** CBC+HMAC əvəzinə; Seal/Open + NonceSize;
   Pad YOX (öz padding var); açar: crypto/rand-dan 32 bayt.
3. **Doom Principle:** əvvəlcə YOXLAMA (MAC/tag), sonra DECRYPT — heç
   vaxt yad mesajı açma.
4. **DHM/RSA = yalnız açar mübadiləsi/sessiya:** simmetrik sürətli axın;
   forward secrecy üçün hər sessiyaya yeni.
5. **İmza = hash + private key:** authenticity ayrıca məsələdir — HMAC
   (simmetrik, sürətli) və ya RSA (asimmetrik, universal).
6. **İcarə DIY cryptoqrafiya yox:** Zoom/IRS/Hyundai dərsləri — alqoritm
   düzgün, İSTİFADƏ ölüdür; mode seç (ECB YOX), IV random, açar unikal.

## Mənbə
Pages: 278-329 (PDF 279-330; kitab texniki hissəsi ~səh. 326-da bitir —
sonrası afterword/about/credits)
