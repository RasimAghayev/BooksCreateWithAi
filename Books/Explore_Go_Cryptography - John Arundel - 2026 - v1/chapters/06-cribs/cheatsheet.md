# Cheat Sheet — Cribs

## Known-Plaintext Attack

### Crib yoxlaması
```go
func CrackWithCrib(ciphertext string, crib string) int {
    for key := 0; key < 256; key++ {
        decrypted := Decipher(key, ciphertext)
        if strings.Contains(decrypted, crib) {
            return key
        }
    }
    return -1
}
```

**İzah:**
- `strings.Contains` → Dekodlanmış mətndə crib axtarışı
- `-1` → Açar tapılmasa xəta dəyəri

**Mənbə:** Chapter 6, page 107
