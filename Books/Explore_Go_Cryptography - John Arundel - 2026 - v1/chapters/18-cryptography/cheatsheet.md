# Cheat Sheet — Cryptography

## Təhlükəsizlik Qaydaları

### Edilməməli olanlar
- ❌ Öz şifrəni yaratma (roll your own crypto)
- ❌ ECB rejimindən istifadə etmək
- ❌ Təsadüfi olmayan nonce istifadə etmək
- ❌ Açarı kod daxilində saxlamaq

### Edilməli olanlar
- ✅ Standart kitabxanalardan (`crypto/aes`, `crypto/rsa`) istifadə etmək
- ✅ Açar idarəetmə sistemlərindən (KMS) istifadə etmək
- ✅ Təhlükəsiz (random) nonce istifadə etmək
- ✅ Testləri əhatəli yazmaq

**Mənbə:** Chapter 18, page 300
