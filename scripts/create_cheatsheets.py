import os

base_dir = r"D:\Tasks\ai_code\books_ai\Books\Explore_Go_Cryptography - John Arundel - 2026 - v1"

cheatsheets = {
    "00-introduction": """# Cheat Sheet — Introduction

## Go Test Commands

### `go test -v`
Verbose test modunda bütün testləri işə salır.

**Nə edir:** Testlərin hər birinin uğur/ uğursuzluq məlumatını göstərir.

**Sub-komanda/flag izahı:**
- `-v` → Verbose (detallı) çıxış, hər testin nəticəsini göstər
- `-run TestName` → Yalnız müəyyən testi işə sal
- `-cover` → Kod coverage faizini göstər

**Mənbə:** Chapter 2, page 45
""",
    "01-ciphers": """# Cheat Sheet — Ciphers

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
""",
    "02-enciphering": """# Cheat Sheet — Enciphering

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
""",
    "03-deciphering": """# Cheat Sheet — Deciphering

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
""",
    "04-cracking": """# Cheat Sheet — Cracking

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
""",
    "05-keys": """# Cheat Sheet — Keys

## Keyspace Ölçüsü

### Açar uzunluğu və mümkün dəyərlər
| Açar uzunluğu | Mümkün dəyərlər | Təhlükəsizlik səviyyəsi |
|----------------|------------------|--------------------------|
| 1 byte (8 bit) | 256 | Çox zəif |
| 16 byte (128 bit) | 2¹²⁸ | Təhlükəsiz |
| 32 byte (256 bit) | 2²⁵⁶ | Praktiki olaraq qırılmaz |

**Mənbə:** Chapter 5, page 85
""",
    "06-cribs": """# Cheat Sheet — Cribs

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
""",
    "07-passwords": """# Cheat Sheet — Passwords

## Key Stretching

### PBKDF2 nümunəsi
```go
import (
    "golang.org/x/crypto/pbkdf2"
    "crypto/sha256"
)

func PasswordKey(password string, salt []byte) []byte {
    return pbkdf2.Key([]byte(password), salt, 100000, 32, sha256.New)
}
```

**İzah:**
- `pbkdf2.Key` → Şifrəni təkrarlı hash edərək açar yaradır
- `100000` → İterasiya sayı (daha yüksək = daha yavaş = daha təhlükəsiz)
- `32` → 256 bit açar

**Mənbə:** Chapter 7, page 117
""",
    "08-blocks": """# Cheat Sheet — Blocks

## Block Cipher Interface

### `cipher.Block`
```go
import "crypto/cipher"

block, err := aes.NewCipher(key)
if err != nil {
    log.Fatal(err)
}

plaintext := make([]byte, block.BlockSize())
ciphertext := make([]byte, block.BlockSize())

block.Encrypt(ciphertext, plaintext)
block.Decrypt(plaintext, ciphertext)
```

**İzah:**
- `aes.NewCipher` → AES blok şifrəsi yaradır
- `block.BlockSize()` → Blok ölçüsü (AES üçün 16 bayt)
- `Encrypt/Decrypt` → Blokları şifrələyir/açır

**Mənbə:** Chapter 8, page 127
""",
    "09-modes": """# Cheat Sheet — Modes

## GCM Mode

### AES-GCM şifrələməsi
```go
import "crypto/cipher"

aesGCM, err := cipher.NewGCM(block)
nonce := make([]byte, aesGCM.NonceSize())
ciphertext := aesGCM.Seal(nil, nonce, plaintext, nil)
```

**İzah:**
- `NewGCM` → GCM rejimi yaradır
- `NonceSize` → Təsadüfi nonce ölçüsü
- `Seal` → Şifrələyir və autentifikasiya teqini əlavə edir

**Mənbə:** Chapter 9, page 149
""",
    "10-padding": """# Cheat Sheet — Padding

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
""",
    "11-enumeration": """# Cheat Sheet — Enumeration

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
""",
    "12-entropy": """# Cheat Sheet — Entropy

## Entropiya Hesablanması

### Şannon Entropiyası
```
H(X) = -Σ p(x) log₂ p(x)
```
- `p(x)` — Hər simvolun müşahidə edilmə ehtimalı
- `H(X)` — Məlumatın entropiyası (bitlərlə)

**Nümunə:** 16 simvol və eyni ehtimallı parol üçün:
```
H = log₂(16) = 4 bits
```

**Mənbə:** Chapter 12, page 202
""",
    "13-randomness": """# Cheat Sheet — Randomness

## crypto/rand

### Təhlükəsiz təsadüfi açar
```go
import (
    "crypto/rand"
    "encoding/binary"
)

func SecureKey() (uint32, error) {
    var key uint32
    err := binary.Read(rand.Reader, binary.BigEndian, &key)
    return key, err
}
```

**İzah:**
- `crypto/rand` → Kriptoqrafik təhlükəsiz təsadüfiyyət
- `binary.Read` → Baytları tam ədədə çevirir
- `rand.Reader` → OS tərəfindən təchiz edilən təsadüfi mənbə

**Mənbə:** Chapter 13, page 217
""",
    "14-chains": """# Cheat Sheet — Chains

## Hash Chain

### Zəncir hash quruluşu
```
H(Block₁) = H₀
H(Block₂) = H(H₀ + Data₂)
H(Block₃) = H(H₁ + Data₃)
```

**İzah:**
Hər yeni blok əvvəlki hash dəyərini və öz məlumatını götürərək yeni hash yaradır.

**Mənbə:** Chapter 14, page 226
""",
    "15-hashing": """# Cheat Sheet — Hashing

## SHA-256 və HMAC

### SHA-256 hash hesablanması
```go
import (
    "crypto/sha256"
    "encoding/hex"
)

func SHA256(data string) string {
    h := sha256.Sum256([]byte(data))
    return hex.EncodeToString(h[:])
}
```

### HMAC
```go
import "crypto/hmac"

func HMAC(key, data string) string {
    h := hmac.New(sha256.New, []byte(key))
    h.Write([]byte(data))
    return hex.EncodeToString(h.Sum(nil))
}
```

**İzah:**
- `sha256.Sum256` → 256-bit hash dəyəri qaytarır
- `hmac.New` → Açar və hash funksiyası ilə HMAC yaradır

**Mənbə:** Chapter 15, page 252
""",
    "16-coins": """# Cheat Sheet — Coins

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
""",
    "17-authentication": """# Cheat Sheet — Authentication

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
""",
    "18-cryptography": """# Cheat Sheet — Cryptography

## Təhlükəsizlik Qaydaları

### Edilməməli olanlar
- ❌ Öz şifrəni yaratma (roll your own crypto)
- ❌ ECB rejimindən istifadə etmək
- ❌ Təsadüfi olmayan nonce istifadə etmək
- ❌ Açarı kod daxilində saxlamaq

### Edilməli olanlar
- ✅ Standart kitabxanalardan (`crypto/aes`, `crypto/rsa`) istifadə etmək
- ✅ Açar idarəetmə sistemlərindən (KMS) istifadə etmək
- ✅ Təhlükəsiz (random) nonce istifadə etmək
- ✅ Testləri əhatəli yazmaq

**Mənbə:** Chapter 18, page 300
"""
}

for folder, content in cheatsheets.items():
    path = os.path.join(base_dir, "chapters", folder, "cheatsheet.md")
    with open(path, "w", encoding="utf-8") as f:
        f.write(content)
    print(f"Created: {path}")

print(f"Created {len(cheatsheets)} cheatsheet files")
