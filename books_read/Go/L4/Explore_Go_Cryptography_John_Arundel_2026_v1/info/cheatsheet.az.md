# Explore Go: Cryptography — Cheat Sheet (AZ)

## Şifrələmə (müasir standart — hər zaman bu)

```go
// AES-256-GCM: gizlilik + bütövlük + padding BİR YERDƏ
block, err := aes.NewCipher(key)          // key: 32 bayt (AES-256)
gcm, err := cipher.NewGCM(block)
nonce := make([]byte, gcm.NonceSize())
rand.Read(nonce)                          // crypto/rand!
ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

// Açma (dəyişilibsə SƏSLİ xəta):
result, err := gcm.Open(nil, nonce, ciphertext, nil)
if err != nil { /* tampered */ }
```

## Açar yaratmaq

```go
key := make([]byte, 32)
rand.Read(key)        // crypto/rand — MƏCBURİ (math/rand YOX!)
// hex ilə ötür: encoding/hex.EncodeToString / DecodeString
```

## CBC (minimal, bütövlük YOXDUR — MAC əlavə et!)

```go
block, _ := shift.NewCipher(key)             // istənilən cipher.Block
iv := make([]byte, block.BlockSize())
rand.Read(iv)
os.Stdout.Write(iv)                          // IV əvvəlcədən yaz
enc := cipher.NewCBCEncrypter(block, iv)
padded := Pad(plaintext, block.BlockSize())  // PKCS#7
enc.CryptBlocks(dst, padded)
```

## PKCS#7 Padding

```go
func Pad(data []byte, blockSize int) []byte {
    n := blockSize - len(data)%blockSize       // həmişə 1..blockSize
    return append(data, bytes.Repeat([]byte{byte(n)} , n)...)
}
func Unpad(data []byte, blockSize int) []byte {
    n := int(data[len(data)-1])                // son bayt = sayğac
    return data[:len(data)-n]
}
```

## Hashing

```go
digest := sha256.Sum256(message)    // [32]byte
// QADAĞAN: MD5, SHA-1 (collisions sındırılıb)
```

## Parol saxlama

```
saxla: salt (random 16+) + hash(password + salt)
tövsiyə: bcrypt/scrypt (qəsdən yavaş)
EVE-ə QARŞI: rainbow table salt ilə İŞSİZ
```

## HMAC (mesaj bütövlüyü)

```go
mac := hmac.New(sha256.New, key)
mac.Write(message)
tag := mac.Sum(nil)
// Göndər: message + tag
```

## cipher.Block interfeysi (öz cipher-u "pluggable" edir)

```go
type Block interface {
    BlockSize() int
    Encrypt(dst, src []byte)   // dst-yə YAZ — qaytarmır
    Decrypt(dst, src []byte)   // açar receiver struct-dadır
}
// Tələb: tam bloklar; dst ≥ src; yazılı şəkildə "full blocks"
```

## Blok-döngüsü (BlockMode olmadan)

```go
for len(src) > 0 {
    block.Encrypt(dst[:n], src[:n])
    src = src[n:]      // re-slice progress
    dst = dst[n:]
}
```

## Benchmark (crack/hız ölçmə)

```go
func BenchmarkCrack(b *testing.B) {
    for i := 0; i < b.N; i++ { Crack(ciphertext, crib) }
}
// çıxış: ns/op — hər +1 bayt açar ≈ ×256 vaxt
```

## Keyspace enumeration (Next)

```go
// {0,0,0} → {1,0,0} → ... → {255,0,0} → {0,1,0} (carry)
// hamısı 255 → error "no more keys"
// İLK bayt = ən az əhəmiyyətli (little-endian increment)
```

## Crack CLI pattern

```go
crib := flag.String("crib", "", "crib text")
flag.Parse()
ciphertext, _ := io.ReadAll(os.Stdin)
key, err := shift.Crack(ciphertext, *crib)   // tanış plaintext ilə açar tap
```

## Mod seçim cədvəli

| Mod | Gizlilik | Bütövlük | Padding | Status |
|---|---|---|---|---|
| ECB | ✓ | ✗ | lazım | QADAĞAN |
| CBC | ✓ | ✗ | PKCS#7 | minimal |
| CTR | ✓ | ✗ | lazım deyil | orta |
| **GCM** | ✓ | **✓** | **daxili** | **STANDART** |

## Qadağalar (kitabın toplum mesajı)

- Öz cipher icad ETMƏ — "perpetuum mobile kimi"
- math/rand açar üçün — YOX
- ECB — YOX (Zoom ləkəsi, IRS səhvi)
- Parolsuz/ aşağı-entropiyalı açar (qısa qayda, pi rəqəmləri)
- Qısa açar (SHA-256 ilə "uzatmaq" entropi DƏYİŞMİRh)
