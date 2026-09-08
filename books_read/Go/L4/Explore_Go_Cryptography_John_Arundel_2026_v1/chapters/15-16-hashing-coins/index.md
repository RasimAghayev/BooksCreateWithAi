# Chapters 15-16 — Hashing, Coins (səh. 248-277)

## Bu fəsillər nədən bəhs edir?

Mesaj bütövlüyü: hash funksiyaları (LenHash → SumHash → MD5 → SHA-256),
dictionary/salt hücumları, preimage attack. Sonra blockchain: "Bobcoin" —
distributed ledger, double-spend, proof-of-work, mining, 51% hücumu.

## Əsas fikirlər

### 1. Niyə hashing? (Ch15)
CBC zənciri noise səbəbli öz-özünə korlanır — Mallory-in dəyişikliyini
fərqləndirmək olmur. Həll: mesajın **digest**-i (sabit uzunluqda "barmaq izi")
ayrıca ötürülür.

**Hash funksiyasının tələbləri:**
- İstənilən ölçülü input → sabit ölçülü output (məs. 64/256 bit)
- Bərabər paylanım (uniform) — "shy teenagers" klasterlənməsi OLMA
- Kiçik input dəyişikliyi → böyük hash dəyişikliyi (avalanche / sel effekti)

### 2. Sadə təcrübi hash-lərin ifşası

**LenHash** (uzunluq = hash):
```go
// input 16 simvol → 0000000000000010 (hex 16)
func LenHash(input []byte) []byte {
    // len(input) int → binary.BigEndian.PutUint64 → 8 bayt
}
```
- Hücum: eyni uzunluqlu istənilən mesaj eyni hash — Mallory dərhal saxta mesaj
  qurur (preimage attack trivial)

**SumHash** (bayt cəmi = hash):
```go
// 2 fərqli mesaj → fərqli hash: ...0505, ...04f1 — yaxşı başlanğıc
```
- Hücum: hədəf 0505 (1285): 10 dənə 'z' (122×10=1220) + qalan 65 → asan
  saxta mesaj; kiçik dəyişiklik → kiçik hash dəyişikliyi (avalance YOX)

**Real alqoritmlər:**
- **MD5:** blok-blok "tortuous bit-twiddling" — illər boyu standart, amma
  hazırda collisions brute-force ilə saatlarda tapılır (BROKEN)
- **SHA-256:** hazırda praktik hücum YOXDUR; de-facto standart. Narahatsan → SHA-3
```go
digest := sha256.Sum256(message) // [32]byte
```

### 3. Passwords + salt
**Problem:** parolun hash-i saxlanılır; Eve hash siyahısını oğurlayır →
dictionary attack: bütün lüğət sözlərini bir dəfə hash-ləyib cədvəlləşdirir
(rainbow table), hamısına qarşı yoxlayır.

**Həll — salt (duz):**
```go
// saxlanılan: salt (random) + hash(password + salt)
```
- Hər istifadəçi AYRI salt → Eve lüğəti hər parol üçün YENİDƏN hash-ləməli
  (one-size-fits-all cədvəz işə yaramır)
- Bonus: bcrypt/scrypt — xatırlatmaq üçün qəsdən YAVAŞ funksiyalar

### 4. Coins — Bobcoin və blockchain (Ch16)

**Fiat valyuta problemi:** Bob-un mərkəzi bankı = single point of failure +
"Bob-a niyə etibad?" (Bob-un serveri elektriksiz → ticarət dayanır; Bob
milyonları oğurlaya bilər)

**Distributed ledger:** hər iştirakçıda nüsxə; **eventually consistent** —
bir neçə saniyədən sonra hamı razılaşır.

**Double-spend problemi (Sam):** 10 Bobcoin pizza + 10 Bobcoin dondurma —
 hansı transaction "birinci"? Transaction-ları SIRA-LAMAQ lazımdır.

**Həll — zəncir (CBC-dən ilhamla):**
- Hər blok əvvəlki blokun HASH-ini ehtiva edir → sıra yalnız bir tərfdən
  doğrulanır → double-spend aşkar olunur
- İki node fərqli "növbəti blok" versiyası göndərirsə → zəncir SPLIT olur
  → **ən uzun zəncir qalib** (doomed stub atılır)

**Proof of Work (mining):**
- Blok qəbul olunsun: hash(block) < target difficulty (məs. hash 0 ilə
  başlasın) → milyonlarla cəhd (nonce dəyişdir → yenidən hash)
- Uğurlu node → transaction fee-lərdən mükafat
- **51% hücumu:** Mallory yarısından çox hash gücünü toplasa saxta zənciri
  "daha uzun" edə bilər — amma həqiqi zəncir daha çox güclə tez üstün olur;
  mükafat strukturu səni düzgün maqnitləməyə yönəldir
- Bitcoin target-i: qəbul şansı ~1/10^64; vaxt keçdikcə difficulty artır

**Proof of Stake (alternativ):** səs haqqı = investisiya — "ən varlılar ən az
əvvəl səni aldadır". Amma bu da sistemi zənginlərin lehinə "rigging" edir.

**Torpaq duyğusu:** "mən Bitcoin haqqında nə düşündümsə... ən yaxşısı öz
research" — müəllifin mövqeyi: texnologi maraqlıdır, ideologiya şübhəli.

## Əsas terminlər

- Digest (həzm / barmaq izi)
- Preimage Attack (ilkin-obraz hücumu)
- Collision (toqquşma)
- Avalanche Effect (sel effekti)
- Salt (duz)
- Distributed Ledger (paylanmış hesab dəftəri)
- Double-Spend (ikiqat xərcləmə)
- Proof of Work (iş sübutu)
- 51% Attack (51% hücumu)

## Praktik nəticə

- Öz hash funksiyanızı YAZMAYIN — SHA-256 istifadə edin
- Parol saxlama: salt + (bcrypt/scrypt), MD5/SHA-1 QADAĞAN
- Hash = bütövlük, encryption = gizlilik; ikisi bir yerdə (MAC/HMAC — Ch17)
- Blockchain mahiyyətcə: hash-zəncirli, PoW ilə sıralanan transaction jurnalıdır

## Mənbə

Pages: 248-277 (Chapters 15-16, Explore Go: Cryptography)
