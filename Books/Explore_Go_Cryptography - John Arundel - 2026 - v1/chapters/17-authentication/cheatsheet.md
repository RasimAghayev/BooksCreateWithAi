# Cheat Sheet — Authentication

## HMAC və AEAD

### HMAC yoxlaması
```go
func VerifyMAC(key, message, mac string) bool {
    expected := HMAC(key, message)
    return hmac.Equal([]byte(mac), []byte(expected))
}
```

**İzah:**
- `hmac.Equal` → Təhlükəsiz müqayisə (zaman hücumuna qarşı)

**Mənbə:** Chapter 17, page 278
