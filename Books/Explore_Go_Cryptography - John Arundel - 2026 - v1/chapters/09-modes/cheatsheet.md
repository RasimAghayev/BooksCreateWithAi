# Cheat Sheet — Modes

## GCM Mode

### AES-GCM şifrələməsi
```go
import "crypto/cipher"

aesGCM, err := cipher.NewGCM(block)
nonce := make([]byte, aesGCM.NonceSize())
ciphertext := aesGCM.Seal(nil, nonce, plaintext, nil)
```

**İzah:**
- `NewGCM` → GCM rejimi yaradır
- `NonceSize` → Təsadüfi nonce ölçüsü
- `Seal` → Şifrələyir və autentifikasiya teqini əlavə edir

**Mənbə:** Chapter 9, page 149
