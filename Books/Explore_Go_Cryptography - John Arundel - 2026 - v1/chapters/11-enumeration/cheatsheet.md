# Cheat Sheet — Enumeration

## Parallel Crack

### Goroutine ilə paralel kırma
```go
func CrackParallel(ciphertext string, crib string) int {
    var wg sync.WaitGroup
    result := make(chan int)
    
    for key := 0; key < 256; key++ {
        wg.Add(1)
        go func(k int) {
            defer wg.Done()
            if strings.Contains(Decipher(k, ciphertext), crib) {
                result <- k
            }
        }(key)
    }
    
    wg.Wait()
    close(result)
    return <-result
}
```

**İzah:**
- `sync.WaitGroup` → Goroutine-lərin bitməsini gözləyir
- `go func()` → Paralel goroutine yaradır
- `result <- k` → Tapılan açarı kanala göndərir

**Mənbə:** Chapter 11, page 189
