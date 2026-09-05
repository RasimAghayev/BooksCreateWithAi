# Cheat Sheet — Randomness

## crypto/rand

### Təhlükəsiz təsadüfi açar
```go
import (
    "crypto/rand"
    "encoding/binary"
)

func SecureKey() (uint32, error) {
    var key uint32
    err := binary.Read(rand.Reader, binary.BigEndian, &key)
    return key, err
}
```

**İzah:**
- `crypto/rand` → Kriptoqrafik təhlükəsiz təsadüfiyyət
- `binary.Read` → Baytları tam ədədə çevirir
- `rand.Reader` → OS tərəfindən təchiz edilən təsadüfi mənbə

**Mənbə:** Chapter 13, page 217
