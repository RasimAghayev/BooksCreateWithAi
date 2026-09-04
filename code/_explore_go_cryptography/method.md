# Explore Go: Cryptography Method

## Book Information
- **Author**: John Arundel
- **Publisher**: Bitfield Consulting
- **Date**: February 13, 2026
- **Copyright**: © 2025 John Arundel
- **GitHub**: https://github.com/bitfield/eg-crypto

## Book Overview
This book teaches cryptography by building a shift cipher in Go from scratch using TDD (Test-Driven Development). It starts with simple Caesar cipher concepts and progressively builds to modern cryptographic techniques including AES-GCM, SHA-256, HMAC, RSA, Diffie-Hellman-Merkle key exchange, blockchain, and quantum cryptography.

---

## Table of Contents

1. **Ciphers** (страницы 1-50) — Introduction to cryptography, shift cipher, TDD approach
2. **Enciphering** (страницы 51-110) — Variable keys, command-line flags, filter programs
3. **Deciphering** (страницы 111-150) — Brute force, frequency analysis, cribs, cracking
4. **Cracking** (страницы 151-200) — Crack function, CLI cracker tool, binary data
5. **Keys** (страницы 201-260) — Keyspace, multi-byte keys, key length
6. **Cribs** (страницы 261-310) — Cracking long keys, brute force
7. **Passwords** (страницы 311-330) — Password security, rainbow tables, salting
8. **Blocks** (страницы 331-390) — Block ciphers, cipher.Block interface
9. **Modes** (страницы 391-450) — Block modes, CBC mode, CryptBlocks
10. **Padding** (страницы 451-500) — Block alignment, PKCS padding
11. **Enumeration** (страницы 501-560) — Key enumeration, benchmarking
12. **Entropy** (страницы 561-620) — Entropy, compression, Kolmogorov complexity
13. **Randomness** (страницы 621-680) — PRNG, seeding, hardware entropy
14. **Chains** (страницы 681-740) — ECB weaknesses, CBC mode, CTR mode
15. **Hashing** (страницы 741-810) — Hash functions, MD5/SHA-1/SHA-256
16. **Coins** (страницы 811-870) — Cryptocurrency, blockchain, proof of work
17. **Authentication** (страницы 871-930) — MAC, HMAC, key exchange, RSA
18. **Cryptography** (страницы 931-990) — AES, implementation, quantum computing
19. **Conclusion** (страницы 991-1024) — Future of cryptography

---

## Chapter Details

### Chapter 1: Ciphers
**Pages 1-50.** Introduces cryptography concepts starting with the shift cipher (Caesar cipher). Key concepts:
- **Shift cipher**: Each byte is shifted by a key value. With key=1, 'A'→'B', 'H'→'I', etc.
- **TDD approach**: Write tests first, then implement code. Tests verify behavior: `Encipher("HAL")` → `"IBM"`.
- **Table tests**: Multiple test cases in a struct, looped with `t.Run` for subtests.
- **Test philosophy**: "See the test fail first" — prove it detects bugs before trusting it.
- **Decipher**: Inverse of encipher: `Decipher(ciphertext, key)` = `Encipher(ciphertext, -key)`.
- **Characters**: Alice (sender), Bob (receiver), Eve (eavesdropper), Mallory (attacker).

### Chapter 2: Enciphering
**Pages 51-110.** Adding variable keys to the shift cipher.
- **Variable keys**: Instead of fixed key=1, Alice can choose any byte value (0-255) as the key.
- **Key flag**: Using Go's `flag` package to define `-key` command-line flag.
- **Hexadecimal keys**: Key specified as hex string (e.g., `DEADBEEF`), decoded via `hex.DecodeString`.
- **Filter pattern**: Read from stdin, transform, write to stdout. Enables piping: `echo "text" | go run ./cmd/encipher -key 5`.
- **Multi-byte keys**: Key is a `[]byte` that repeats cyclically. Key `DEADBEEF` (4 bytes) applied to 5-byte input uses bytes D, E, A, D, E.
- **Frequency analysis**: Same plaintext letter no longer produces same ciphertext letter with multi-byte keys.

### Chapter 3: Deciphering
**Pages 111-150.** Breaking ciphers.
- **Brute force**: Try all 255 possible keys.
- **Frequency analysis**: Eve guesses most common ciphertext letter represents 'E'.
- **Cribs**: Known plaintext fragments used to verify guessed keys. E.g., if plaintext starts with "GO", test each key against first 2 bytes.
- **Crack function signature**: `Crack(ciphertext, crib []byte) (key byte, err error)`.
- **Magic numbers**: PNG files start with `0x89` + "PNG". Used as cribs for binary files.
- **Binary data**: Enciphering images works same as text — cipher operates on raw bytes.

### Chapter 4: Cracking
**Pages 151-200.**
- **Command-line cracker**: `cmd/crack` — reads ciphertext from stdin, uses `-crib` flag, outputs plaintext.
- **Error handling**: `errors.New("no key found")` when brute force fails.
- **Table tests extended**: Same test cases work for Crack as for Encipher/Decipher.

### Chapter 5: Keys
**Pages 201-260.**
- **Keyspace**: With 1-byte keys, only 255 possible keys — easily brute-forced.
- **Exponential scaling**: Each additional byte multiplies keyspace by 256. 2 bytes = 65,536 keys. 32 bytes = 2^256 ≈ 10^77 (more than atoms in universe).
- **Energy cost**: Brute-forcing 256-bit key requires ~10^37 supernovae energy — physically impossible.
- **Multi-byte key implementation**: Key is `[]byte`, indexed cyclically: `c.key[i % len(c.key)]`.
- **Modulus optimization**: Using `%` instead of wrapping with conditional logic.

### Chapter 6: Cribs
**Pages 261-310.**
- **Multi-byte crack**: For each key byte position, try all 256 values against corresponding crib byte.
- **MaxKeyLen constant**: 32-byte limit for key guessing.
- **Key enumeration**: Build key one byte at a time, using `bytes.Equal` to verify.

### Chapter 7: Passwords
**Pages 311-340.** (extracted in page 251-260 partial)
- **Password strength**: Longer passwords exponentially harder to crack (7 chars: 6 min, 14 chars: 200M years).
- **Rainbow tables**: Pre-computed hash lookup tables. Defeated by salting.
- **Salting**: Adding random data to password before hashing. Makes rainbow tables useless.
- **Password hashing**: Use slow algorithms (bcrypt, scrypt) — fast hashes like SHA-256 are too fast for passwords.

### Chapter 8: Blocks
**Pages 341-390.**
- **Block ciphers**: Encrypt fixed-size blocks (32 bytes in our example).
- **cipher.Block interface**: `BlockSize() int`, `Encrypt(dst, src []byte)`, `Decrypt(dst, src []byte)`.
- **shiftCipher struct**: Stores key as `[BlockSize]byte`, implements `cipher.Block`.
- **NewCipher function**: Validates key length, returns `cipher.Block`.

### Chapter 9: Modes
**Pages 391-440.**
- **cipher.BlockMode interface**: `BlockSize() int`, `CryptBlocks(dst, src []byte)`.
- **Simple block mode**: Each block encrypted independently.
- **CBC mode** (Cipher Block Chaining): XOR each plaintext block with previous ciphertext block before encrypting.
- **Initialization Vector (IV)**: Random value XORed with first block — ensures same plaintext produces different ciphertext.
- **Block alignment**: Data must be padded to exact block size before encryption.

### Chapter 10: Padding
**Pages 441-490.**
- **PKCS padding**: Add N bytes of value N to fill last block.
- **Pad function**: `padding := bytes.Repeat([]byte{byte(n)}, n)`.
- **Unpad function**: Read last byte to determine how many padding bytes to remove.
- **Illegal pixel problem**: Block cipher output may not be valid UTF-8/text.

### Chapter 11: Enumeration
**Pages 501-540.**
- **Key enumeration**: Systematic key guessing.
- **Benchmarking**: `go test -bench` to measure crack speed.
- **Realistic crack time estimation**: How long to brute-force N-byte key on typical hardware.
- **Parallelisation**: Splitting work across CPUs — but doesn't fundamentally change exponential complexity.

### Chapter 12: Entropy
**Pages 561-600.**
- **Entropy**: Measure of uncertainty/unpredictability. More entropy = more random = more secure.
- **Kolmogorov complexity**: Length of shortest program that produces the sequence.
- **Compression = low entropy**: Compressible data has patterns. Ciphertext shouldn't compress.
- **Randomness**: True random = high entropy. Pseudo-random = looks random but is deterministic.
- **Key generation**: Use cryptographically secure random (e.g., `/dev/urandom`, hardware RNG).

### Chapter 13: Randomness
**Pages 621-680.**
- **PRNG (Pseudorandom Number Generator)**: Deterministic but looks random. Needs good seed.
- **Fibonacci generator**: Uses recurrence relation. Periodicity issues.
- **Seeding**: `rand.Seed(time.Now().UnixNano())` — time-based seed.
- **Cryptographic RNG**: `crypto/rand` package, reads from system entropy pool.
- **Quantum randomness**: Physically unpredictable — photon polarization, etc.
- **Side-channel attacks**: Even if math is perfect, implementation can leak information.

### Chapter 14: Chains
**Pages 681-740.**
- **ECB (Electronic Codebook)**: Each block encrypted independently. Same plaintext block → same ciphertext block. **INSECURE**.
- **CBC (Cipher Block Chaining)**: XOR with previous ciphertext. Requires IV.
- **CTR (Counter mode)**: Uses counter + key, XOR with plaintext. Turns block cipher into stream cipher.
- **Nonce**: Number used once — ensures uniqueness in CTR/GCM.
- **Block replay**: ECB pattern reveals structure.

### Chapter 15: Hashing
**Pages 741-810.**
- **Hash functions**: Map arbitrary data → fixed-size digest.
- **Properties**: Deterministic, fast, avalanche effect (tiny input change → big output change).
- **MD5**: Broken — collisions trivial to find. Don't use.
- **SHA-1**: Weakened — collisions found with ~$100K cloud compute. Don't use.
- **SHA-256**: Current standard. No known practical attacks.
- **Preimage attack**: Constructing input that hashes to a given value.
- **LenHash**: Hash = length of input. Trivially broken.
- **SumHash**: Hash = sum of input bytes. Better but still broken.

### Chapter 16: Coins
**Pages 811-870.**
- **Cryptocurrency**: Digital currency using blockchain.
- **Double-spending problem**: How to prevent spending the same coin twice.
- **Blockchain**: Public ledger of all transactions.
- **Proof of Work**: Computational puzzle that's hard to solve but easy to verify.
- **Proof of Stake**: Alternative — validators chosen by stake.
- **Bobcoin example**: Walkthrough of digital currency creation.

### Chapter 17: Authentication
**Pages 871-930.**
- **Message Authentication Code (MAC)**: Authenticate message + verify integrity.
- **Length extension attacks**: Attacker can append to message if they know the hash of `message + key`.
- **HMAC**: Hash-based MAC. `H(key + H(key + message))` — prevents length extension.
- **Key exchange**: Alice and Bob need to agree on a key.
- **Diffie-Hellman-Merkle**: Allows key exchange over insecure channel.
  - Alice and Bob agree on public (p, g)
  - Alice sends g^a mod p, Bob sends g^b mod p
  - Shared secret: g^(ab) mod p
- **RSA**: Public-key cryptosystem based on difficulty of factoring large primes.

### Chapter 18: Cryptography
**Pages 931-990.**
- **AES (Advanced Encryption Standard)**: Modern symmetric cipher.
  - Block size: 128 bits (16 bytes)
  - Key sizes: 128, 192, or 256 bits
  - 14 rounds for 256-bit keys
  - Uses S-boxes, ShiftRows, MixColumns, AddRoundKey
- **AES-CBC**: AES in CBC mode with IV.
- **AES-GCM (Galois Counter Mode)**: Authenticated encryption. Provides both confidentiality and integrity.
- **Implementation fails**: Even secure algorithms can be misused (ECB mode, zero IV, shared keys).
- **Quantum computing future**: Threat to RSA/DH (Shor's algorithm), but AES-256 is post-quantum secure.

---

## Code Project Structure
```
shift/     - Shift cipher implementation
  cmd/encipher/ - Encipher command-line tool
  cmd/decipher/ - Decipher command-line tool  
  cmd/crack/    - Crack command-line tool
  shift.go      - Encipher, Decipher, Crack functions
  shift_test.go - Table-driven tests

hash/      - Hash function implementations (LenHash, SumHash)
  cmd/hash/
  hash.go
  hash_test.go

aes/       - AES encryption tools
  cmd/encipher/
  cmd/decipher/
  aes.go (Pad, Unpad functions)
```

## Key Go Concepts Covered
- `t.Parallel()` for parallel test execution
- Table-driven tests with `t.Run` subtests
- `bytes.Equal` for slice comparison
- `flag` package for CLI flags
- `encoding/hex` for hex encoding/decoding
- `crypto/sha256` for SHA-256 hashing
- `crypto/aes` and `crypto/cipher` for AES
- `crypto/rand` for secure random
- `golang.org/x/crypto/bcrypt` for password hashing
