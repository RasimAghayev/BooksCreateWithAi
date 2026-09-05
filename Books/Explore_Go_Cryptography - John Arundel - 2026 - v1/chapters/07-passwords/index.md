# 7. Passwords

**Səhifələr:** 114-126

## Bu fəsil nədən bəhs edir?

Şifrə (password) əsaslı şifrələmə (password-based encryption). Şifrələri təhlükəsiz açar halına çevirmək üçün key stretching və salt mexanizmləri.

## Əsas fikirlər

### 1. Password-Based Keys
İnsanların yaddaşda saxlaya biləcəyi şifrələri kriptoqrafik açar halına çevirmək.

### 2. Key Stretching (Açar Uzatma)
Şifrənin təhlükəsizliyini artırmaq üçün təkrarlı hash funksiyaları tətbiq etmək (PBKDF2, scrypt, bcrypt).

### 3. Salt (Duz)
Eyni şifrə ilə eyni hash alınmasının qarşısını almaq üçün təsadüfi məlumat əlavə etmək.

## Əsas terminlər
- Password-Based Encryption (Şifrə əsaslı şifrələmə)
- Key Stretching (Açar uzatma)
- Salt (Duz)
- PBKDF2
- Hash Function (Hash funksiyası)

## Praktik nəticə

Şifrə əsaslı şifrələmə üçün Go funksiyaları yazılır. Salt və key stretching tətbiq edilir.

## Mənbə

Pages: 114-126
