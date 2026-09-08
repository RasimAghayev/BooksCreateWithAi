# Chapters 17-18 — Authentication, AES (səh. 278-329)

## Bu fəsillər nədən bəhs edir?

MAC/HMAC (mesaj autentifikasiyası), açar mübadiləsi problemi (Diffie-Hellman),
public-key crypto (RSA), CA/etimad zənciri. Və nəhayət: AES — tarixi, daxili
struktur (SubBytes/ShiftRows/MixColumns/AddRoundKey), crypto/aes ilə
implementasiya, AES-GCM (şifrə + bütövlük bir paketdə).

## Əsas fikirlər

### 1. Niyə yalnız hash azdır? (Ch17)
Alice hash-i ayrıca ötürür — amma Mallory BOTH-u (mesaj + hash) dəyişə bilər.
Həll: **MAC** — hash(key + message). Mallory açarı bilmir → düzgün hash
qura bilmirz.

**HMAC ( Nested variant):**
```
inner = H(key + message)
outer = H(key + inner)   // göndərilən: ciphertext + outer
```
- Bir addımlı `H(key + msg)` length-extension hücumuna açıqdır — nested
  quruluş bunu kəsir
```go
mac := hmac.New(sha256.New, key)  // crypto/hmac
mac.Write(message)
tag := mac.Sum(nil)
```

### 2. Açar mübadiləsi problemi
Simmetrik crypto üçün Alice-Bob eyni açarı paylaşmalı — amma təhlükəsiz kanal
YOXDUR (olsaydı, crypto lazım olmazdı!).

**Diffie-Hellman-Merkle (DHM):**
- İctai razılaşdırılmış: əsas `g`, modul `p` (public)
- Alice: `A = g^a mod p` göndərir; Bob: `B = g^b mod p`
- Alice hesablayır: `B^a mod p` = `g^(ba) mod p`
- Bob hesablayır: `A^b mod p` = `g^(ab) mod p` → **EYNİ session açarı!**
- Eve görür: `g^a, g^b` → amma `g^ab` tapmaq üçün **discrete logarithm**
  həll etməli — one-way funksiya (brute-force qədər çətin)
- Zəif tərəfi: TAM mesaj replay hələ mümkündür; authentication (kimlik) yoxdur

### 3. Public-key crypto (RSA)
**İki açar:** public (hamıya) + private (yalnız sənə).
- Encrypt(public) → yalnız private ilə açılır → gizli mesaj
- Encrypt(private) = imza → yalnız sən edə bilərdin → authenticity

**Riyaziyyat:** public `n = p × q` (iki BÖYÜK primal); private = p, q.
Multiply asan, **faktorizasiya çətin** — one-way. (DHM: eksponensiya asan,
logaritm çətin.)

**Etimad problemi (Trent / CA):** bankofbob.com-un public açarı həqiqətən
bankınkımıdır? → **Certificate Authority** Trent imza edir; CA-lər öz-özlərini
imzalayırlar (root CA-lər OS/browser-da "pre-installed"). TLS bunun üzərində
qurulub.

### 4. AES — Advanced Encryption Standard (Ch18)
**Tarix:** DES (1970s, 56-bit — NSA-nın o dövrkü crack gücü haqqında ipucu!) →
açık müsabiqə (NIST) → **Rijndael** qalib (Rijmen + Daemen) → AES: blok 128
bit (16 bayt), açar 128/192/256 bit.

**Struktur (4×4 grid, column-major):**
1. **SubBytes** — hər baytı S-box cədvəlindən əvəz et (confusion / qarışdırma)
2. **ShiftRows** — sətirləri rotasiya et (diffusion / yayılma)
3. **MixColumns** — hər sütuna matris transformasiyası (diffusion)
4. **AddRoundKey** — round açarı XOR (açarın təsiri)

- 10/12/14 round (açar ölçüsünə görə); son round MixColumns-suz
- **Konfusion + Difuziya:** SubBytes birtərəfli S-box qarışdırır, ShiftRows/
  MixColumns statistik struktur dağıdır — ikisi BİRLİKDƏ güclü
- Dəhşətli səmərəli: pre-generated lookup cədvəlləri (te0..te3) ilə bütün
  round-lar bir neçə sətirdə; Go-nun crypto/aes/block.go belə yazılıb

**"Klüzi itələmə" (Rubber-hose):** insider Wikipedia — çox güclü
implementasiyalar belə istifadəçi səhvləri ilə düşür.

### 5. AES-in Go-da istifadəsi
```go
block, err := aes.NewCipher(key)          // key: 32 bayt (AES-256)
encrypter := cipher.NewCBCEncrypter(block, iv)
encrypter.CryptBlocks(dst, src)           // bizim shift cipher ilə EYNİ interfeys!
```
- `cipher.Block` interfeysi sayəsində shift cipher → AES dəyişməsi 2 sətir
- CLI (encipher/decipher): yalnız `shift.NewCipher` → `aes.NewCipher`

### 6. AES-GCM — tövsiyə olunan son həll
CBC yalnız gizlilik verir. **GCM = CBC + MAC BİRLİKDƏ:**
```go
block, _ := aes.NewCipher(key)
gcm, _ := cipher.NewGCM(block)            // AEAD interfeysi
nonce := make([]byte, gcm.NonceSize())
rand.Read(nonce)
ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
// göndərilən: nonce || ciphertext || auth tag — HAMISI Seal-in çıxışında

plaintext, err := gcm.Open(nil, nonce, ciphertext, nil) // err ≠ nil → TAMPARED
```
**Sub-kod izahı:**
- `Seal` öz padding-i edir → Pad/Unpad EHTİYACI YOXDUR
- `Open` authenticity-ni avtomatik yoxlayır — Mallory-in dəyişikliyi SƏSLİ xəta verir
- "Bu, Go-da müasir kriptoqrafiya üçün yazacağımız ƏSASƏN BÜTÜN kod" — kitabın yekun mesajı

### 7. Zəifliklər və kvant (finale)
- AES-in ÖZÜ deyil, İSTİFADƏSİ səhv olarkən təhlükə var: Zoom-un ECB-si,
  IRS-in ECB tələbi, QNap NAS-ın 10 rəqəmli "açarı"
- **Kvant kompüterlər:** qubit superpozisiyası — 256 qubit = bütün açarları
  "paralel" sınamaq (Grover alqoritmi kvadrat sürətləndirir → 256-bit AES →
  effektiv 128-bit, hələ güclü). Şor alqoritmi RSA/DHM-i sındırır.
- Decoherence problemi: ölçmə (və ya stray photon) superpozisiyanı dağıdır —
  böyük kvant kompüuterlər hələ uzaqdır. Amma "post-quantum crypto"
  araşdırmaları davam edir.

## Əsas terminlər

- MAC (Message Authentication Code)
- HMAC (Nested MAC standartı)
- Diffie-Hellman-Merkle (açar mübadiləsi)
- Discrete Logarithm (diskret logaritm)
- Public/Private Key (açiq/xüsusi açar)
- CA (Certificate Authority / sertifikat orqanı)
- S-Box (substitusiya cədvəli)
- Confusion / Diffusion (qarışdırma / yayılma)
- AEAD (Authenticated Encryption with Associated Data)
- AES-GCM (tövsiyə olunan mode)
- Post-Quantum Crypto (post-kvant kriptoqrafiya)

## Praktik nəticə

- Müasir Go kriptoqrafiyası: `aes.NewCipher` + `cipher.NewGCM` + `crypto/rand`
  nonce — başqa heç nə yazma
- MAC/HMAC mesaj bütövlüyü üçün; sadə hash yox
- Açar mübadiləsi: DHM (session) və ya RSA (public-key); TLS = sertifikat + etimad zənciri
- ECB heç vaxt; CBC minimal; **GCM standart**
- Öz cipher-u İXTİRA ETMƏ — "perpetuum mobile ilə məşğul olmaq kimi"

## Mənbə

Pages: 278-329 (Chapters 17-18, Explore Go: Cryptography)
