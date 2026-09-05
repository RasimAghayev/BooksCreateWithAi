# Cheat Sheet — Enciphering

## Go Test Kodu

### `encipher_test.go`
```go
func TestEncipher(t *testing.T) {
    cases := []struct {
        key      int
        plaintext string
        expected string
    }{
        {1, "A", "B"},
        {2, "ABC", "CDE"},
    }
    for _, c := range cases {
        got := Encipher(c.key, c.plaintext)
        if got != c.expected {
            t.Errorf("Encipher(%d, %s) = %s, want %s", c.key, c.plaintext, got, c.expected)
        }
    }
}
```

**İzah:**
- `struct { ... }` → Test case-lərin saxlanması üçün strukt
- `t.Errorf` → Test uğursuz olduqda xəta mesajı yazdır
- `for _, c := range cases` → Cədvəli iterator

**Mənbə:** Chapter 2, page 49
