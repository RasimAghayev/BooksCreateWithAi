# 10. Padding

**Səhifələr:** 162-177

## Bu fəsil nədən bəhs edir?

Blok şifrələrində blok həcmi tam olmayan mətnlərin necə işlənəcəyi. Padding (doldurma) sxemləri, xüsusilə PKCS#7 standartı.

## Əsas fikirlər

### 1. Niyə padding lazımdır?
Blok şifrələri sabit ölçülü bloklar üzərində işləyir (məsələn, AES — 128 bit = 16 bayt). Son blok tam deyilsə, padding ilə doludur.

### 2. PKCS#7
Hər padding baytı, əlavə edilən padding baytlarının sayını göstərir. Məsələn, 3 bayt padding lazımdırsa, hər 3 padding baytı dəyəri 3 olur.

### 3. Padding təmizləmə
Deşifrələmə zamanı son baytdan padding sayını öyrənib lazımsız baytları silmək.

## Əsas terminlər
- Padding (Doldurma)
- PKCS#7
- Block Size (Blok ölçüsü)
- Unpadding (Doldurmanı silmək)

## Praktik nəticə

Padding və unpadding funksiyaları yazılır və testlənir.

## Mənbə

Pages: 162-177
