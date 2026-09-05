# Cheat Sheet — Padding

## PKCS#7 Padding

### Padding əlavə etmək
```go
func Pad(data []byte, blockSize int) []byte {
    padLen := blockSize - len(data)%blockSize
    padding := bytes.Repeat([]byte{byte(padLen)}, padLen)
    return append(data, padding...)
}
```

**İzah:**
- `padLen` → Neçə bayt əlavə edilməli olduğunu hesablayır
- `bytes.Repeat` → `padLen` dəyərli baytları təkrarlayır
- `append` → Padding-i əsas məlumata əlavə edir

**Mənbə:** Chapter 10, page 165
