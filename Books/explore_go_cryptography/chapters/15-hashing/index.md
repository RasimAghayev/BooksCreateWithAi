# Chapter 15 — Hashing (səh. 248-267)

## Bu chapter nədən bəhs edir?

Mesaj bütövlüyü (integrity): hash funksiyaları, digest, distribution, collision,
preimage attack, naiv hash-lərin (LenHash, SumHash) sınması, MD5/SHA-1 tarixçəsi,
SHA-256 standartı və parol hash-ləmə (dictionary attack, rainbow table, salt).

## Əsas fikirlər

### 1. Mesaj bütövlüyü problemi
**Nədir:** CBC Mallory-in blokları dəyişməsinin qarşını tam almır — blokların
biri dəyişilsə, decrypt pozulur, amma Bob **nəyin** pozulduğunu bilmir.
Real kanallarda da bit xətaları olur (thermal noise, cosmic rays) — TCP/IP
packet-lərin pozulmasını yoxlayıb yenidən göndərir.

**Sınaqlı işarə ("ENDIT") yolverilmezdir:** proqnozlaşdırıla bilən mətn Eve-ə
crib verir — kriptoqrafik həll lazımdır.

### 2. Hash və digest
**Nədir:** İstənilən uzunluqdakı mesajı **sabit uzunluqlu** nümayəndəyə
(digest, məs. 256 bit) çevirmək — bitləri qarışdırmaq (muddle).

**Məqsədlər:** bütövlük yoxlaması (digest-i müqayisə et), data strukturları
(hash table).

**Buckets və distribution:** Hash table-də düzgün paylanma (uniform
distribution) vacibdir — bütün "Z" hərfli fayllar bir bucket-a düşsə, tam skan
geri qaytarır. **Adversarial key set-lərə** (qəsdən toqquşduran düşmən) hazır
olmaq lazımdır — yalnız orta halı deyil, **worst case**-i də qiymətləndir.

### 3. LenHash — naiv sınaq
**Nədir:** Hash = input-un uzunluğu:
```go
func LenHash(input []byte) []byte {
    digest := make([]byte, 8)
    binary.BigEndian.PutUint64(digest, uint64(len(input)))
    return digest
}
```
**Problemlər:**
1. **Collision-lar kütləvi:** eyni uzunluqda bütün mesajlar eyni hash — Mallory
   istənilən alternativ mesajı asanlıqla tapır
2. **Pis distribution:** nəticələr çox kiçik ədədlərdə cluster-lənir (bitlərin
   çoxu sıfır) — hashspace-in kiçik küncündə

### 4. SumHash — bir az yaxşı, amma kifayətsiz
**Nədir:** Hash = bütün baytların cəmi:
```go
func SumHash(input []byte) []byte {
    digest := make([]byte, 8)
    var sum uint64
    for _, b := range input {
        sum += uint64(b)
    }
    binary.BigEndian.PutUint64(digest, sum)
    return digest
}
```
**Preimage attack (sınaq):** Hədəf hash 0505 = 1285 (decimal). Mallory byt-ları
seçir: "zzzzzzzzzz" = 10×122 = 1220, +65 ("A") = 1285 → `zzzzzzzzzzA` hazır.
Məhz bir dəqiqəlik hesabla istənilən hash üçün preimage qurulur.

**Problemlər:**
1. Distribution yenə pis — sum-lar kiçik ədədlərdə cluster
2. **Kiçik input dəyişikliyi = kiçik hash dəyişikliyi** ("A"→"B" = hash +1) —
   Mallory-ə preimage qurmağı asanlaşdırır. Yaxşı hash-da 1 bit dəyişmə
   **avalanche** (çığ) — digest-in yarısı dəyişməlidir

### 5. Real hash alqoritmləri
| Alqoritm | Status | Qeyd |
|---|---|---|
| MD5 | SÖZÜLMÜŞ | Collision-lar 1-2 saniyədə brute-force ilə; kriptoqrafik məqsədlərə tam uyğunsuz |
| SHA-1 | SÖZÜLMÜŞ | Practical collision attacks mövcuddur (2017) |
| SHA-256 | STANDART | Hazırda praktik hücum yoxdur; de facto standart |
| SHA-3 | Alternativ | Fərqli ailə, əlavə ehtiyat üçün |

MD5-in işləmə üsulunun dərsi: input-u bloklara bölür, hər bloku bir-birinə
dolaşır (tortuous bit-twiddling) — beş il state-of-the-art idi, sonra sındı.
**Alqoritmlərin ömrü var** — "sabah da güvənilir" yoxdur.

**SHA-256 istifadəsi:**
```go
digest := sha256.Sum256(data) // [32]byte
```

### 6. Parol hash-ləmə
**Problem:** Alice banka parolu göndərmədən onu bildiyini sübut etməli; bank da
parolu bilməməlidir (kicking the can down the road olmasın).

**Həll:** Qeydiyyatda bank `hash(password)` saxlayır; login-də `hash(input)`-ı
müqayisə edir. Parol heç vaxt şəbəkədən keçmir, serverdə düz saxlanmır.

**Eve-in hücumu — dictionary attack:** Bütün lüğət sözlərini bir dəfə hash-lə
→ **rainbow table** → istənilən sızmış parol siyahısına lookup cədvəli.

**Həll — salt:**
```
store: salt || hash(password || salt)
```
- Salt = istifadəçi başına **təsadüfi** dəyər, hash ilə birlikdə açıq saxlanılır
- Eve lüğəti yalnız bir dəfə yox, **hər parol üçün yenidən** hash-lamalıdır →
  bir-size-fits-all rainbow table ölü
- Slowness də vacibdir: parol hash-ləri üçün qəsdən YAVAŞ alqoritmlər
  (bcrypt, scrypt, Argon2) işlədilir

## Əsas terminlər

- Digest (məhaz / özet)
- Hash Function (hash funksiyası)
- Collision (toqquşma)
- Preimage Attack (ilkin şəkil hücumu)
- Distribution / Uniform (paylanma / bərabər)
- Avalanche Effect (çığ effekti)
- Rainbow Table (göy qurşağı cədvəli)
- Salt (duz)
- Dictionary Attack (lüğət hücumu)

## Praktik nəticə

- Öz hash funksiyanı yazma — LenHash/SumHash sınaqları göstərdi ki, naiv
  yanaşmalar dərhal sınır
- MD5/SHA-1 heç vaxt; SHA-256 (və ya SHA-3) standartdır
- Yaxşı hash: uniform distribution + avalanche + preimage üçün brute-force
- Parol saxlama: `salt + yavaş hash` (bcrypt/scrypt/Argon2); SHA-256 parol
  üçün çox sürətlidir
- Digest müqayisəsi = bütövlük yoxlaması (şifrə AÇMADAN dəyişikliyi tap)

## Mənbə

Pages: 248-267 (Chapter 15, Explore Go: Cryptography)
