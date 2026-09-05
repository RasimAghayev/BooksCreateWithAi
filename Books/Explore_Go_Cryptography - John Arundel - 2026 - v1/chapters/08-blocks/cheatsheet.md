# Cheat Sheet — Blocks

## Block Cipher Interface

### `cipher.Block`
```go
import "crypto/cipher"

block, err := aes.NewCipher(key)
if err != nil {
    log.Fatal(err)
}

plaintext := make([]byte, block.BlockSize())
ciphertext := make([]byte, block.BlockSize())

block.Encrypt(ciphertext, plaintext)
block.Decrypt(plaintext, ciphertext)
```

**İzah:**
- `aes.NewCipher` → AES blok şifrəsi yaradır
- `block.BlockSize()` → Blok ölçüsü (AES üçün 16 bayt)
- `Encrypt/Decrypt` → Blokları şifrələyir/açır

**Mənbə:** Chapter 8, page 127
