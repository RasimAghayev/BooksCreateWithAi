# Cheat Sheet — Cracking

## Scoring Sistemi

### `Score` funksiyası
```go
func Score(text string) int {
    score := 0
    for _, r := range text {
        if freq[r] > 0 {
            score += freq[r]
        }
    }
    return score
}
```

**İzah:**
- `freq[r]` → İngilis dilində hərfin tezliyinə görə skor
- Yüksək skor = daha çox tanımlı hərflər = daha ehtimallı plaintext

**Mənbə:** Chapter 4, page 80
