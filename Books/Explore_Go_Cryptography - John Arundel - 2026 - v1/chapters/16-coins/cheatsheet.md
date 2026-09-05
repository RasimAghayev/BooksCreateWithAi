# Cheat Sheet — Coins

## Commitment Scheme

### Rəqəmsal imdatlaşdırma
```go
func Commit(value int, nonce int) [32]byte {
    data := []byte(fmt.Sprintf("%d:%d", value, nonce))
    return sha256.Sum256(data)
}

func Open(commitment [32]byte, value int, nonce int) bool {
    data := []byte(fmt.Sprintf("%d:%d", value, nonce))
    return commitment == sha256.Sum256(data)
}
```

**İzah:**
- `Commit` → Dəyəri gizləyir, lakin sonradan açmağı mümkün edir
- `Open` → İmdatlaşdırmanı doğrular

**Mənbə:** Chapter 16, page 268
