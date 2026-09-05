# 8. Blocks

**Səhifələr:** 127-145

## Bu fəsil nədən bəhs edir?

Blok şifrələri (block ciphers) və iki əsas növ: stream cipher (axın şifrəsi) və block cipher (blok şifrəsi). ECB və CBC rejimləri.

## Əsas fikirlər

### 1. Block Cipher vs Stream Cipher
- **Block Cipher:** Məlumatı sabit ölçülü bloklar halında emal edir (AES, DES)
- **Stream Cipher:** Məlumatı bit/bit və ya bayt/bayt emal edir

### 2. Electronic Codebook (ECB)
Hər blok müstəqil şifrələnir. Eyni plaintext həmişə eyni ciphertext yaradır — təhlükəsiz deyil.

### 3. Cipher Block Chaining (CBC)
Hər blok əvvəlki şifrələnmiş blokla əlaqələndirilir. Initialization Vector (IV) ilə başlanır.

## Əsas terminlər
- Block Cipher (Blok şifrəsi)
- Stream Cipher (Axın şifrəsi)
- ECB (Electronic Codebook)
- CBC (Cipher Block Chaining)
- IV (Initialization Vector)

## Praktik nəticə

Block cipher interface (cipher.Block) öyrənilir. ECB və CBC rejimləri müqayisə edilir.

## Mənbə

Pages: 127-145
