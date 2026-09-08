# Explore Go: Cryptography — Mündəricat

**Müəllif:** John Arundel | **Nəşriyyat:** Bitfield Consulting | **İl:** 2026 | **Səviyyə:** L4 (Advanced)

Kitabın mövzusu: kriptoqrafiya SIFIRDAN — Caesar cipher-dan AES-GCM-yə qədər
tam TDD səfəri. Hər fəsil GOAL → HINT → SOLUTION formatında: oxucu shift
cipher-in özünü yazır, sındırır (Crack!), blok/mode/padding qurur, entropiya
və random-un fəlsəfəsini öyrənir, ECB-nin vizual iflasını GÖRÜR, CBC-ni
tətbiq edir, hash/MAC/HMAC/digital imzaları qurur və nəhayət AES-GCM ilə
professional kriptoqrafiyaya çatır.

## Bilik faylları (fəsil qrupları üzrə)

| Fayl | Fəsillər | Mövzular | Səhifə |
|---|---|---|---|
| [01-04-ciphers-encipher-decipher-crack](../chapters/01-04-ciphers-encipher-decipher-crack/index.md) | 1-4 | Codes vs ciphers, shift cipher TDD, table tests, filter CLI, frequency analysis, brute-force + cribs, Decipher refaktoru | 20-84 |
| [05-07-keys-cribs-passwords](../chapters/05-07-keys-cribs-passwords/index.md) | 5-7 | Keyspace (256-bit = kainat), multi-byte açarlar (%), hex flag, crib ≥ açar, parol siyasətləri = zərər, K^N, random tərifi | 85-126 |
| [08-10-blocks-modes-padding](../chapters/08-10-blocks-modes-padding/index.md) | 8-10 | cipher.Block/BlockMode interfeyslər, NewCipher + ErrKeySize wrap, blok döngüsü, panic konvensiyası, PKCS#7 Pad/Unpad | 127-177 |
| [11-13-enumeration-entropy-randomness](../chapters/11-13-enumeration-entropy-randomness/index.md) | 11-13 | Next açar sanması, benchmark extrapolasiyası, paralelləşdirmə fizikası, log2 N entropiya, Kolmogorov K(s), PRNG tələləri, /dev/random, TRNG, kvant | 178-225 |
| [14-16-chains-hashing-coins](../chapters/14-16-chains-hashing-coins/index.md) | 14-16 | ECB vizual iflası, Mallory (replay/drop), CTR/nonce, CBC+IV, LenHash/SumHash dərsləri, MD5/SHA-1 ölü, SHA-256, parol hash + salt + bcrypt, Bobcoin, 51%, PoW/PoS | 226-277 |
| [17-18-authentication-aes](../chapters/17-18-authentication-aes/index.md) | 17-18 | Chosen ciphertext + Doom, MAC/HMAC, açar mübadiləsi, DHM (işlənmiş nümunə), RSA + imza, TLS + CA zənciri, AES tarixi/daxili strukturu, AES-CBC keçidi, AES-GCM (Seal/Open), kvant + post-quantum | 278-329 |

## Oxu ardıcıllığı tövsiyəsi

Ardıcıl 1→18: kitab TƏK layihədir — shift paketi hər fəsildə təkamül edir
(Encipher → multi-key → block → CBC → AES). Fəsil 11-12 (entropiya) və 14
(ECB) ayrıca təkrar oxunmağa dəyər — ən dərin nəzəriyyə oradadır.

## Kitabın xarakteri

Arundel-in imza üslubu: hekayə ilə öyrənmə (Alice/Bob/Eve/Mallory), GOAL/
HINT/SOLUTION interaktivliyi, hər texniki qərar üçün "niyə" + "niyə yox"
müzakirəsi. Kod nümunələri minimal, amma hər biri tam testlənmişdir. 329
səhifə, amma sıx emal ilə ~6 bilik faylına sığır.
