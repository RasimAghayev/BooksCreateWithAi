# Cheat Sheet — Ciphers

## Shift Cipher Kodu

### `caesar.go`
```go
func Encipher(key int, plaintext string) string {
    var result strings.Builder
    for _, r := range plaintext {
        if r >= 'A' && r <= 'Z' {
            result.WriteRune(shiftRune(r, key))
        } else {
            result.WriteRune(r)
        }
    }
    return result.String()
}
```

**İzah:**
- `key int` → Açar dəyəri (neçə hərf dəyişdirmək)
- `strings.Builder` → Sətirləri birləşdirmək üçün effektiv strukt
- `shiftRune` → Hərfi açar qədər dəyişdirən köməkçi funksiya

**Mənbə:** Chapter 1, page 22
