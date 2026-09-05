# 15. Hashing

**Səhifələr:** 248-267

## Bu fəsil nədən bəhs edir?

Hash funksiyaları, SHA-256 və HMAC (Hash-based Message Authentication Code). Mətn bütövlüyünü yoxlamaq və autentifikasiya.

## Əsas fikirlər

### 1. Hash Function
Hər hansı bir məlumatdan sabit uzunluqlu hash dəyəri çıxaran riyazi funksiya. Tək yönlüdür (geri qayıtmaz).

### 2. SHA-256
256-bit (32 bayt) hash dəyəri yaradan təhlükəsiz hash funksiyası. Bitcoin və bir çox sistemlərdə istifadə olunur.

### 3. HMAC
Hash funksiyasından istifadə edərək mesaj autentifikasiyasını təmin edən mexanizm. Açar və məsləhət (secret) istifadə edir.

## Əsas terminlər
- Hash Function (Hash funksiyası)
- SHA-256
- HMAC (Hash-based Message Authentication Code)
- Digest (Hesaslama)
- Message Integrity (Mesaj bütövlüyü)

## Praktik nəticə

Go-da `crypto/sha256` və `crypto/hmac` paketləri ilə hash hesablama və HMAC yaratmaq.

## Mənbə

Pages: 248-267
