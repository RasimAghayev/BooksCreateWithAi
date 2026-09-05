# Müəllim Qeydləri — Explore Go: Cryptography

**Müəllif:** John Arundel  
**İl:** 2026

## 📖 Kitab deyir:

Kitab, kriptoqrafiyanı öyrənmək üçün test-driven development (TDD) metodikasından istifadə edir. Hər chapter-da yeni bir kriptoqrafik konsept öyrədilir və həmin konsept üçün Go kodları yazılır.

1. **Sadəlikdən başlamaq:** Shift cipher (sürüşmə şifrəsi) kimi ən sadə şifrələrdən başlayaraq mürəkkəb alqoritmlərə qədər.
2. **Test yazmaq qanunudur:** Kod yazmazdan əvvəl test yazmaq, testi keçmək, sonra refactor etmək.
3. **Təcrübəyə əsaslanmaq:** Hər konsept üçün praktik Go kodları təqdim edilir.
4. **Təhlükəsizlik prinsipləri:** Öz şifrəni yaratmaqdan (roll your own crypto) çəkinmək.

## 👨‍🏫 Müəllim qeydi:

Məncə bu kitab təxminən 2-3 ay ərzində öyrənilməsi çətin deyil. Hər chapter 15-25 səhifədir və praktiki tətbiqlərdən ibarətdir. TDD metodikasına üstünlük verirsənsə, bu kitab sənin üçün idealdır.

### Ən vacib 5 fikir

1. **TDD nə deyilsə, test yaz:** Hər funksiyadan əvvəl test yazmaq alış-verişinizi dəyişəcək. Sadəcə "işləyir" demək kifayət deyil — hər edge case üçün test yazmaq öyrənməyi dərinləşdirir.
2. **Açar uzunluğu təhlükəsizliyi müəyyən edir:** 1-byte açar (256 mümkün dəyər) asan kırılır. 256-bit açar (2²⁵⁶ mümkün dəyər) praktiki olaraq qırılmazdır.
3. **ECB rejimi təhlükəsiz deyil:** Eyni plaintext həmişə eyni ciphertext yaradır. Hər zaman CBC, GCM və ya AEAD kimi təhlükəsiz rejimlərdən istifadə edin.
4. **Öz şifrəni yaratma:** Heç vaxt özünüz şifrə alqoritmi icad etməyin. `crypto/aes`, `crypto/hmac`, `crypto/rsa` kimi standart Go paketlərindən istifadə edin.
5. **Entropiya vacibdir:** Təhlükəsiz şifrələr üçün yüksək entropiyalı (təsadüfiy) açarlar və nonce-lər lazımdır.

### Kitabın ən dəyərli hissəsi

**Chapter 5 — Keys (səhifələr 85-103)** — çünki açar məkanı (keyspace) anlayışı kriptoqrafiyanın əsasını təşkil edir. Açar uzunluğunun təhlükəsizlik üzərində təsiri və brute force hücumunun limitləri bu chapter-da əla izah edilir.

### Tövsiyələr

- **Praktika:** Hər chapter-da yazılan testləri özünüz yazın və `go test` ilə işlədiyinizə əmin olun.
- **Genişləndirmək:** Kitabda AES-CBC öyrənilir, lakin real layihələrdə GCM daha təhlükəsizdir. `crypto/cipher` paketinin GCM dəstəyini öyrənin.
- **Təhlükəsizlik auditləri:** Kriptoqrafik sistemlərin qarşısını almaq üçün professional auditlər keçirmək lazımdır. Bu kitab təməl verir, lakin real dünya sistemləri daha mürəkkəbdir.
- **Daha çox oxu:** Kitab "Real-World Cryptography" (David Wong) ilə birlikdə oxumaq tövsiyə olunur.
