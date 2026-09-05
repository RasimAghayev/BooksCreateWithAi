# Cheat Sheet — Deciphering

## Subtests

### `t.Run` istifadəsi
```go
func TestDecipher(t *testing.T) {
    cases := []struct{...}{...}
    for _, c := range cases {
        t.Run(fmt.Sprintf("key=%d", c.key), func(t *testing.T) {
            got := Decipher(c.key, c.ciphertext)
            if got != c.expected {
                t.Fatalf("got %q, want %q", got, c.expected)
            }
        })
    }
}
```

**İzah:**
- `t.Run` → Alt test yaradır, adlandırma effektli debug imkanı verir
- `fmt.Sprintf` → Dinamik test adı yaratmaq

**Mənbə:** Chapter 3, page 72
