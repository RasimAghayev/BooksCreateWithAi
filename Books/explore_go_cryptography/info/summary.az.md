# Explore Go: Cryptography — Xülasə (AZ)

## Kitab kimin üçündür?

Go biliyi olan, kriptoqrafiyanı **praktik TDD yolu ilə** öyrənmək istəyən
developer-lər üçün. Riyaziyyat dərinliyi əvəzinə: hər konsept sadə kodla
qurulur, sonra OXŞAR HÜCUM İLƏ SINDIRILIR, sonra möhkəmləndirilir — və bu
dövr kitab boyu təkrarlanır.

## Əsas xətt (18 fəsil, 329 səh.)

1. **Sadə başlanğıc (Ch1-4):** Caesar/shift cipher — encipher/decipher/crack.
   Cribs (bilinən plaintext) + brute-force. Testable kod, mod 256 bayt dünyası.
2. **Açarlar (Ch5-7):** qısa açarın zəifliyi → 32-bayt (256-bit) açar;
   açar yaratmaq üçün parolların LIMIT-ləri; SHA-256 ilə "uzatma" entropi
   artırmır.
3. **Blok dünyası (Ch8-10):** cipher.Block interfeysi, blok-döngüsü,
   PKCS#7 padding (həmişə pad — "illegal pixel" problemi).
4. **Açarların mənası (Ch11-13):** keyspace enumeration (Next funksiyası),
   benchmarking, entropiya (Kolmogorov), PRNG vs TRNG, kvant.
5. **Rejimlər (Ch14):** ECB-nin şəkil ifşası → CBC + IV zənciri.
6. **Bütövlük (Ch15-17):** hashing (LenHash → SHA-256), salt, MAC/HMAC,
   açar mübadiləsi (DHM), public-key (RSA), CA etimadı.
7. **Real crypto (Ch18):** AES strukturu + AES-GCM — "yazacağımız əsasən
   bütün müasir kriptoqrafiya kodu".

## Ən vacib 5 fikir

1. **Təhlükəsizlik nisbidir** — "secure/not secure" yox; hücumçunun resursları
   vs məlumatın dəyəri. 256-bit açar = enerji/fizika limitinə qədər güclü.
2. **Açar = entropiya** — random bitlər (`crypto/rand`); "ağıllı görünən"
   qısa qaydalı açarlar (pi, seq, parollar) zəifdir.
3. **İnterfeyslər dəyişdirilə bilənliyi verir** — `cipher.Block` / `BlockMode` /
   `AEAD` sayəsində öz shift cipher → AES keçidi 2 sətir.
4. **ECB QADAĞANDIR** — eyni input → eyni output; replay, codebook, şəkil
   ifşası. CBC zənciri + IV; amma bütövlük üçün MAC/AEAD lazımdır.
5. **Yekun kod:** `aes.NewCipher` + `cipher.NewGCM` + `gcm.Seal/Open` +
   random nonce — padding, MAC, IV hamısı daxildir.

## Kitabın ən dəyərli hissəsi

Chapter 14 (ECB → CBC) və Chapter 18 (AES-GCM): birinci "şifrələyirəm, amma
məntiqi strukturu qoruyuram" səhvini VİZUAL göstərir (şəkil), ikinci isə
müasir standartın nə qədər sadə olduğunu göstərərək kitabı bağlayır.

## Metodoloji qeydlər

- Kitab tam TDD ilə yazılıb: hər funksiya əvvəl test, sonra testdən çıxan
  GOAL/HINT/SOLUTION ritorikası
- Benchmarking dərsi: benchmarking ilə brute-force sürətini ölçmək →
  eksponensial böyümənin intuisiyası
- Hər chapter epigrafla açılır (Schneier, Knuth, Orwell, Usual Suspects)
