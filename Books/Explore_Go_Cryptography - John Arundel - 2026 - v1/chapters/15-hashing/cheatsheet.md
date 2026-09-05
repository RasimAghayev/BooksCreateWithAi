# Cheat Sheet — Hashing

## SHA-256 və HMAC

### SHA-256 hash hesablanması
```go
import (
    "crypto/sha256"
    "encoding/hex"
)

func SHA256(data string) string {
    h := sha256.Sum256([]byte(data))
    return hex.EncodeToString(h[:])
}
```

### HMAC
```go
import "crypto/hmac"

func HMAC(key, data string) string {
    h := hmac.New(sha256.New, []byte(key))
    h.Write([]byte(data))
    return hex.EncodeToString(h.Sum(nil))
}
```

**İzah:**
- `sha256.Sum256` → 256-bit hash dəyəri qaytarır
- `hmac.New` → Açar və hash funksiyası ilə HMAC yaradır

**Mənbə:** Chapter 15, page 252
