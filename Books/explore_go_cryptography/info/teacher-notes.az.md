# Explore Go: Cryptography — Müəllim Qeydləri (AZ)

📖 Kitab deyir: shift cipher-dən AES-GCM-ə qədər hər şeyi öz əlinlə TDD ilə
qur, sonra sındır, sonra möhkəmləndir. Nəticə: `aes.NewCipher` +
`cipher.NewGCM` + `gcm.Seal` — "bu, Go-da müasir kriptoqrafiya üçün yazacağımız
əsasən bütün koddur".

👨‍🏫 Müəllim qeydi: bu yanaşma (öyrənmə üçün) əladır, amma production-da
KITABIN ÖZÜNÜ yazma — yalnız son faylı köçür. 4 real praktikada görünən
nuans:

1. **Nonce İDARƏSİ GCM-də vacibdir.** Kitab `rand.Read(nonce)` göstərir —
   doğru başlanğıcdır, amma productionda hər Seal üçün UNIQUE nonce şərtdir;
   counter əsaslı (random başlanğıc + artan) və ya bölünmüş (96-bit)
   quruluşlar mövcuddur. Nonce təkrarı GCM-də KATASTROFİKDIR (açar axını
   pozulur) — kitab bunun üstündə kifayət qədər dayanmır.

2. **Aead III məlumat (associated data):** `gcm.Seal(nonce, nonce, plaintext,
   nil)` — 4-cü parametr (nil) başlıqlar/kontekst üçündür. HTTPS-də TLS
   header-ləri şifrələnmir amma BÜTÖVLÜKDƏ qorunmalıdır — real
   sistemlərdə bu parametri istifadə etmək tez-tez lazım olur.

3. **Açar saxlanması kitabda toxunulmur.** Açarı harada saxlamaq (env var,
   KMS, Vault), rotasiya (dövri dəyişmə), "key wrapping" — production
   kriptoqrafiyasının YARISI məhz budur. "Cryptography is the art of
   transforming information security problems into key management problems"
   epigrafı kitabın özü ilə təsdiqlənir, amma həll productionda axtarılmalıdır.

4. **Post-quantum hazırlığı:** kitabın son fəsli kvant hücumunu
   "hələ uzaq" kimi təqdim edir — əsasən doğru, amma uzunömürlü sirrlər
   üçün ("harvest now, decrypt later") bugündən PQ-hybrid TLS (X25519+
   ML-KEM) düşünmək lazımdır.

## Ən vacib 5 fikir

1. Təhlükəsizlik nisbidir — hücumçu resursu vs data dəyəri
2. Entropiya = açarın gücü; random ≥ ağıllı
3. ECB strukturu göstərir; CBC zənciri; GCM hamısını bağlayır
4. MAC/HMAC bütövlük üçün; GCM-də daxili
5. Standart paketlər (`crypto/*`) hər şeyi verir — özünü yazma

## Kitabın ən dəyərli hissəsi

Chapter 14: ECB-nin şəkil ifşası — "şifrələmişəm" hissinin NƏ ÜÇÜN aldatıcı
olduğunu bir dəfə və həmişəlik göstərir. Bu, təhsildə görə bildiyim ən effektli
kripto dərsi.

⚠️ Uyğunsuzluq/qeyd: kitabın shift cipher kodu tədris məqsədlidir — "bizim
32-bayt blok shift cipher" real qorunma vermir (hər blok müstəqil eyni XOR
məntiqi). Müəllif bunu açıq deyir; oxucunun shift-i "istifadə etməməsi"
tələb olunur.
