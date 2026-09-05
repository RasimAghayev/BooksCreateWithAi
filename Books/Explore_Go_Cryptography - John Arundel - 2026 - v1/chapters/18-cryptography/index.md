# 18. Cryptography

**Səhifələr:** 295-329

## Bu fəsil nədən bəhs edir?

Real dünya kriptoqrafiyası, ümumi təhlükələr, öz şifrəni yaratmaq (roll your own) haqqında məsləhətlər və yanlışlar.

## Əsas fikirlər

### 1. Real-World Cryptography
Kitabda yazılan sadə şifrələri real həyatda işlətmək təhlükəlidir. Profesional kriptoqrafik kitabxanalardan istifadə edilməlidir.

### 2. Common Pitfalls (Ümumi təhlükələr)
- ECB rejimindən istifadə etmək
- Təsadüfi olmayan (predictable) nonce istifadə etmək
- Öz şifrəni icad etmək
- Açarı kod daxilində saxlamaq

### 3. Standard Libraries
Go-nun `crypto/aes`, `crypto/rsa`, `crypto/hmac` kimi paketlər təhlükəsiz və sınanmış tətbiqləri təmin edir.

## Əsas terminlər
- Real-World Cryptography (Real dünya kriptoqrafiyası)
- Standard Library (Standart kitabxana)
- Side-Channel Attack (Yan kanal hücumu)
- Key Management (Açar idarəetməsi)

## Praktik nəticə

Öz şifrəni yaratmaqdan çəkinmək. Profesional kitabxanalardan və açar idarəetmə sistemlərindən istifadə etmək.

## Mənbə

Pages: 295-329
