# Cheat Sheet — Explore Go: Cryptography

Bu fayl kitabın bütün chapter-lərindən toplanan komandalar, kod blokları və konfiqurasiya nümunələrini ehtiva edir.

## Go Test Commands

### `go test -v`
```bash
go test -v
```
**Nə edir:** Verbose test modunda bütün testləri işə salır.

**Sub-komanda/flag izahı:**
- `-v` → Verbose (detallı) çıxış
- `-run TestName` → Yalnız müəyyən testi işə sal
- `-cover` → Kod coverage faizini göstər

**Mənbə:** Chapter 2, page 45

## Shift Cipher Kodu

### `encipher.go`
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
- `key int` → Açar dəyəri
- `strings.Builder` → Sətirləri birləşdirmək üçün effektiv strukt
- `shiftRune` → Hərfi açar qədər dəyişdirən köməkçi funksiya

**Mənbə:** Chapter 1, page 22

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

**Mənbə:** Chapter 2, page 49

## Subtests

### `t.Run` istifadəsi
```go
t.Run(fmt.Sprintf("key=%d", c.key), func(t *testing.T) {
    got := Decipher(c.key, c.ciphertext)
    if got != c.expected {
        t.Fatalf("got %q, want %q", got, c.expected)
    }
})
```

**İzah:**
- `t.Run` → Alt test yaradır, adlandırma effektli debug imkanı verir

**Mənbə:** Chapter 3, page 72

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

**Mənbə:** Chapter 4, page 80

## Keyspace Ölçüsü

| Açar uzunluğu | Mümkün dəyərlər | Təhlükəsizlik səviyyəsi |
|----------------|------------------|--------------------------|
| 1 byte (8 bit) | 256 | Çox zəif |
| 16 byte (128 bit) | 2¹²⁸ | Təhlükəsiz |
| 32 byte (256 bit) | 2²⁵⁶ | Praktiki olaraq qırılmaz |

**Mənbə:** Chapter 5, page 85

## Key Stretching

### PBKDF2 nümunəsi
```go
import "golang.org/x/crypto/pbkdf2"

func PasswordKey(password string, salt []byte) []byte {
    return pbkdf2.Key([]byte(password), salt, 100000, 32, sha256.New)
}
```

**İzah:**
- `100000` → İterasiya sayı
- `32` → 256 bit açar

**Mənbə:** Chapter 7, page 117

## Block Cipher Interface

```go
import "crypto/cipher"

block, err := aes.NewCipher(key)
plaintext := make([]byte, block.BlockSize())
ciphertext := make([]byte, block.BlockSize())
block.Encrypt(ciphertext, plaintext)
block.Decrypt(plaintext, ciphertext)
```

**İzah:**
- `aes.NewCipher` → AES blok şifrəsi yaradır
- `block.BlockSize()` → Blok ölçüsü (AES üçün 16 bayt)

**Mənbə:** Chapter 8, page 127

## GCM Mode

```go
import "crypto/cipher"

aesGCM, err := cipher.NewGCM(block)
nonce := make([]byte, aesGCM.NonceSize())
ciphertext := aesGCM.Seal(nil, nonce, plaintext, nil)
```

**İzah:**
- `NewGCM` → GCM rejimi yaradır
- `Seal` → Şifrələyir və autentifikasiya teqini əlavə edir

**Mənbə:** Chapter 9, page 149

## PKCS#7 Padding

```go
func Pad(data []byte, blockSize int) []byte {
    padLen := blockSize - len(data)%blockSize
    padding := bytes.Repeat([]byte{byte(padLen)}, padLen)
    return append(data, padding...)
}
```

**İzah:**
- `padLen` → Neçə bayt əlavə edilməli olduğunu hesablayır

**Mənbə:** Chapter 10, page 165

## Parallel Crack

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

**Mənbə:** Chapter 11, page 189

## Entropiya Hesablanması

```
H(X) = -Σ p(x) log₂ p(x)
```

**Nümunə:** 16 simvol və eyni ehtimallı parol üçün:
```
H = log₂(16) = 4 bits
```

**Mənbə:** Chapter 12, page 202

## crypto/rand

```go
import "crypto/rand"

func SecureKey() (uint32, error) {
    var key uint32
    err := binary.Read(rand.Reader, binary.BigEndian, &key)
    return key, err
}
```

**İzah:**
- `crypto/rand` → Kriptoqrafik təhlükəsiz təsadüfiyyət
- `rand.Reader` → OS tərəfindən təchiz edilən təsadüfi mənbə

**Mənbə:** Chapter 13, page 217

## Commitment Scheme

```go
func Commit(value int, nonce int) [32]byte {
    data := []byte(fmt.Sprintf("%d:%d", value, nonce))
    return sha256.Sum256(data)
}
```

**İzah:**
- `Commit` → Dəyəri gizləyir, lakin sonradan açmağı mümkün edir

**Mənbə:** Chapter 16, page 268

## SHA-256 və HMAC

```go
import "crypto/sha256"

func SHA256(data string) string {
    h := sha256.Sum256([]byte(data))
    return hex.EncodeToString(h[:])
}
```

**İzah:**
- `sha256.Sum256` → 256-bit hash dəyəri qaytarır

**Mənbə:** Chapter 15, page 252
