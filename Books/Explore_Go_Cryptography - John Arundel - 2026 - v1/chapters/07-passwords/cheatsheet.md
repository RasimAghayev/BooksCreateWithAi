# Cheat Sheet — Passwords

## Key Stretching

### PBKDF2 nümunəsi
```go
import (
    "golang.org/x/crypto/pbkdf2"
    "crypto/sha256"
)

func PasswordKey(password string, salt []byte) []byte {
    return pbkdf2.Key([]byte(password), salt, 100000, 32, sha256.New)
}
```

**İzah:**
- `pbkdf2.Key` → Şifrəni təkrarlı hash edərək açar yaradır
- `100000` → İterasiya sayı (daha yüksək = daha yavaş = daha təhlükəsiz)
- `32` → 256 bit açar

**Mənbə:** Chapter 7, page 117
