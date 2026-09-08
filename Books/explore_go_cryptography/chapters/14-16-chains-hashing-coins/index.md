# Chapters 14-16 — Chains, Hashing, Coins (səh. 226-277)

## Bu fəsillər nədən bəhs edir?

ECB-nin VİZUAL iflası (devil.ppm şəkli: header (3 sətir) + payload ayrılır;
enciphered şəkil GÖRÜNÜR — 32-bayt açar təkrarı = "stippling"; random
noise qızıl standartı), Electronic Code Book (eyni blok → eyni ciphertext;
lookup-table = code book), **Mallory** (oxuya BİLƏN VƏ DƏYİŞƏ BİLƏN
hücumçu — Eve-dən fərqli!): replay attack ($1000 köçürməsini təkrar göndər),
block dropping (son blokları sil — sonu yoxdursa aşkarlanMAZ), block
modifying; CTR mode (counter + key + plaintext; nonce = "for the nonce" —
bir dəfəlik, GİZLİ OLMASINA EHTİYAC YOXDUR), CBC mode (əvvəlki ciphertext =
blok üçün nonce → zəncir; eyni plaintext ≠ eyni ciphertext; replay/drop/
modify AŞKARLANIR; "blockchain" termini buradan!), IV (ilk blokun nonce-u —
random; mesajdan ƏVVƏL açıq göndərilir; eyni mesaj + müxtəlif IV = müxtəlif
ciphertext; IV dəyişsə yalnız İLK blok pozulur), crypto/cipher NewCBCEncrypter/
NewCBCDecrypter ilə tam keçid (encipher: rand.Read(iv) → stdout-a IV +
ciphertext; decipher: ilk blok = IV), hashing: mesaj bütövlüyü (digest —
fiks uzunluq nümayəndə; "bit-flipping/block-dropping" qarşısı), hash
cədvəlləri (buckets, uniform distribution, clustering — T/A/O dolu, X/Z boş;
collisions DOSTDUR amma hücumçu üçün DÜŞMƏN), LenHash (uzunluq = 64-bit
digest; binary.BigEndian.PutUint64 — TRİVİAL preimage: "I hate you, Bob"
eyni hash!), SumHash (bayt cəmi; avalanche YOXDUR — A→B dəyişikliyi hash-i
yalnız 1 vahid dəyişir; preimage hücumu: z×10 + "7" = hədəf hash), MD5
(1-2 saniyə collision — ÖLDÜ), SHA-1 ($100k cloud collision — ÖLDÜ),
SHA-256 (praktik hücum YOXDUR; sha256.Sum256(data); SHA-3 mövcud), password
hashing (zero-knowledge: parol YOX, hash saxlanılır; parol heç vaxt şəbəkədən
keçmir), rainbow table (dictionary + bad passwordsın əvvəlcədən hash-ləri =
code book analogi), salting (hər istifadəçiyə unikal random data + hash;
eyni parol = fərli hash; rainbow table MƏHV — hər parol üçün yenidən
dictionary hash-ləmək), yavaş hashing (parol üçün YAVAŞ = YAXŞI; gündəlik
login 3ms fərq etməz, Eve milyonlarla hash-ləyir; bcrypt/scrypt — SHA-256
parol üçün PİS), Bobcoin (fiat; mərkəzi Bob = single point of failure +
etimadsizlıq), distributed ledger (public + hər nodada; eventually
consistent), double-spending (Sam: pizza + ice cream eyni 10 Bobcoin;
transaction sıralaması lazım), hash blockchain (blok = əvvəlki blokun
hash-i + transaction datası; sıra yalnız hash zənciri ilə doğrulanır),
consensus (2 namizəd blok → növbəti blok hansına zəncirlənibsə O qalib;
doomed stub atılır; UZUN zəncir = majority), 51% attack (Mallory bütün
hashrate-in YARISINDAN çoxunu ələ keçirməlidir — böyük şəbəkədə qeyri-mümkün;
safety in numbers), proof of work (hash < target; nonce guess; ~1/10^71
(BTC: 1/10^64); 10 dəqiqədə bir blok; 120 TWh/il — ölkədən çox; "planeti
yandıran riyaziyyat"), proof of stake (ən çox pulu olan = ən etibarlı?
"capitalism" tənqidi).

## Əsas fikirlər

### 1. ECB-nin Vizual İflası (Ch14)
- **PPM formatı:** header (3 sətir) + payload — head/tail ilə ayrılır;
  şəkli encipher edib header-ə birləşdir → VİZUAL nəticə
- **Nəticə:** "all ones" açar → şəkil GÖRÜNÜR (fon qara, gözlər burulğan);
  random 32-bayt açar → hələ də KONTUR GÖRÜNÜR — 32-bayt təkrar =
  0.05% perturbasiya; struktur (uzun eyni-piksel zolaqları) QALIR
- **Qızıl standart:** random noise şəkli — ciphertext ona bənzəməlidir

### 2. Mallory — Aktiv Hücumçu (Ch14)
| Hücum | Mexanizm | ECB-də nəticə |
|---|---|---|
| **Replay** | yazılmış "pay $1000" blokunu yenidən göndər | bank YENİ əmr sanır — jackpot |
| **Block dropping** | transfer blokunu sil | aşkarlanmaz (son bloklardan sonra yoxdur) |
| **Modifying** | bit flip | schemesiz asılı; ECB heç bir qoruma YOX |

- **Eve vs Mallory:** Eve YALNIZ oxuyur; Mallory OXUYUR + DƏYİŞİR

### 3. CTR və Nonce (Ch14)
- **CTR:** counter + key + plaintext → ciphertext; counter əvəzinə hər
  dəyişən dəyər OLAR (vaxt, random)
- **Nonce ("for the nonce"):** bir-dəfəlik dəyər; GİZLİ DEYİL — Bob
  bilməlidir (sayğac tutaraq, hər blokdan əvvəl göndərərək); təkrar
  istifadə = ECB-yə QAYIDIŞ

### 4. CBC — Blockchain (Ch14)
- **Hər blokun girişi:** key + plaintext bu blok üçün + CIPHERTEXT əvvəlki
  blokdən → zəncir; eyni plaintext ≠ eyni ciphertext (əvvəlki blokdan
  asılı); replay → Bob decrypt EDED BİLMİR → aşkar
- **"Blockchain" sözü buradan gəlir!** (cipher block chain)
- **IV (initialization vector):** ilk blokun nonce-u; random blok
  (/dev/random — crypto/rand.Read); açıq göndərilir (IV + ciphertext);
  eyni mesaj + unikal IV = fərli ciphertext → tam mesaj replay qarşısı
- **IV tampering:** IV dəyişsə → yalnız İLK blok pozulur (qalan bloklar
  əvvəlki CIPHERTEXT-dən asılıdır — dəyişməyib)

### 5. CBC Keçidi — crypto/cipher (Ch14)
```go
// encipher:
iv := make([]byte, shift.BlockSize)
rand.Read(iv)                        // crypto/rand
os.Stdout.Write(iv)                  // mesajın başına IV!
enc := cipher.NewCBCEncrypter(block, iv)
plaintext = shift.Pad(plaintext, shift.BlockSize)
enc.CryptBlocks(ciphertext, plaintext)
os.Stdout.Write(ciphertext)          // 468→480→512 (IV bloku + 12 pad)

// decipher:
iv := ciphertext[:shift.BlockSize]  // ilk blok = IV
plaintext := make([]byte, len(ciphertext)-shift.BlockSize)
dec := cipher.NewCBCDecrypter(block, iv)
dec.CryptBlocks(plaintext, ciphertext[shift.BlockSize:])
plaintext = shift.Unpad(plaintext, shift.BlockSize)
```
- **Öz encrypter/decrypter-lərimiz SİLİNDİ** — standart kitabxana CBC-ni
  təqdim edir (kod təmizlənməsi: dəyərsiz testlər/kod çıxarılır)

### 6. Hashing — Bütövlük (Ch15)
- **Problem:** "Bob aldığı ciphertext, Alice-in göndərdiyi İLƏ eynidirmi?"
  — uzunluq prefiksi YARAMIR (Mallory dəyişər; crib verər)
- **Digest:** mesajdan hesablanan FİKS uzunluqlu dəyər; Alice hesablayır
  göndərir → Bob yenidən hesablayır → müqayisə
- **Hash cədvəlləri (kripto XARİC):** buckets + distribution; clustering
  (T/A/O vs X/Z — filing cabinet "VWXYZ" drawer); collisions HİSSƏVƏN
  (həddindən artıq = tam scan, həddindən az = fayda yox)

### 7. Zəif Hash Dərsləri (Ch15)
```go
// LenHash — uzunluq:
binary.BigEndian.PutUint64(digest, uint64(len(input)))
// Preimage TRİVİAL: HƏR HANSİ eyni-uzunluqlu mesaj ("I hate you, Bob")!

// SumHash — bayt cəmi:
var sum uint64
for _, b := range input { sum += uint64(b) }
// Preimage: "zzzzzzzzzz7" (z=122×10 + '7'=55 + newline 10 = hədəf)
```
- **LenHash distribution:** əksər bitlər 0 (2^64 simvol = 17 trilyon GB
  olmalı — YOX); 1GB mesajlar → faktiki 30-bit fəza (milyard dəfə kiçik!)
- **Avalanche effect:** girişin KİÇİK dəyişikliyi → çıxışın BÖYÜK
  dəyişikliyi (qartopu → çığı); SumHash: A→B = hash ±1 — YOXDUR;
  təkmilləşdirmə = kifayət qədər YAXINLAŞ + kiçik tweak-lər
- **ProductHash:** distribution biraz yaxşı — amma hələ uzaq

### 8. Real Hash Alqoritmləri (Ch15)
| Alqoritm | Status |
|---|---|
| MD5 | 1-2 saniyə brute-force collision — MİLLİ |
| SHA-1 | $100k cloud collision (2015) — MİLLİ |
| **SHA-256** | praktik hücum YOXDUR — DE FAKTO standart |
| SHA-3 | fərqli ailə, yeni |
```go
hash := sha256.Sum256(data)   // [32]byte
```
- 256 bit = kifayət (az = qeyri-təhlükəsiz, çox = lazımsız); "bir gün
  hücum tapılacaq — bugün üçün kifayət"

### 9. Password Hashing (Ch15)
- **Zero-knowledge proof:** parolu SAXBANMA — HASH saxlanılır; login:
  hash(giriş) == saxlanılan hash; parol heç vaxt şəbəkədən/stordan keçmir
- **Rainbow table:** əvvəlcədən hash-lənmiş dictionary + bad passwords =
  lookup — ECB code book ANALOGİYASI (eyni parol → eyni digest)
- **Salting:** hər parola unikal random data qarışdır → hər hash fərli;
  rainbow table MƏHV (hər hədəf üçün dictionary YENİDƏN hash-lənməli);
  salt açıq saxlanılır (onsuz yoxlama olmaz) — açıqlıq ZƏRƏR DEYİL
- **Parol üçün YAVAŞ hash:** normal istifadə gündə 1 hash (3ms fərq
  etməz); Eve milyonlarla → YAVAŞLIQ = müdafiə; SHA-256 SÜRƏTLİ GPU-larla =
  parol üçün PİS; **bcrypt/scrypt** — bilə-bilə yavaş və yaddaş-ac

### 10. Bobcoin — Blockchain (Ch16)
- **Fiat → mərkəzi problem:** Bob = single point of failure + "nə üçün
  Bob-a etibad?" → **çözüm: hamıya açıq + paylanmış** (hər node ledgerin
  nüsxəsi; eventually consistent)
- **Double-spending:** Sam 10 Bobcoin ilə pizza + ice cream (yayılma
  gecikməsindən istifadə) → **TRANSAKSİYA SIRALAMASI** lazım
- **Hash zənciri:** blok = hash(əvvəlki blok + transactionlar) → sıra
  yalnız bir yolla doğrulanır (CBC analogiyası — amma ŞİFRƏLMƏ YOX, hamı
  GÖRÜR)
- **Consensus:** 2 namizəd blok → node ikisini də saxlayır → növbəti blok
  hansına zəncirlənibsə O avtoritetdir; digəri "doomed stub" — atılır;
  **uzun zəncir qalib** (more work = majority)
- **Reward:** qalib nodun öz fee-transaction-ı bloka daxil edilir → şəbəkə
  tərəfindən doğrulanır

### 11. 51% Attack və Proof of Work (Ch16)
- **Fırəng node qarşısı #1 — imzalar:** yalnız Alice "Alice ödəyir"
  transaction-ı yarada bilər (Ch17); #2 — hashing PUL ALIR: saxta zənciri
  saxlamaq Mallory-yə honest şəbəkədən baha başa gəlir; #3 — 51%: ümumi
  hashrate-in yarısından çoxu ələ keçirmək (böyük şəbəkədə qeyri-mümkün —
  safety in numbers)
- **Proof of work:** hash(data + nonce) < target; nonce TAHMİN ET (preimage-
  resistant → brute-force); BTC: ~1/10^64 şans/cəhd; 10 dəq/blok; difficulty
  zamanla ARTIR (güc artımını kompensasiya)
- **Qiymət:** 120 TWh/il (ölkədən çox); "planeti yandıran, bəşəriyyətə
  faydasız riyaziyyat"; itirən minlərlə nodun işi KÜLƏ DÖNÜR
- **Proof of stake:** ən varlı nodlar = ən etibarlı? — "capitalizm artıq
  var" tənqidi; adalet-effizient həll HƏLƏ YOXDUR

## Əsas terminlər
- ECB (Electronic Code Book) — eyni blok → eyni ciphertext; QADAĞAN
- Mallory — aktiv hücumçu (oxu + dəyişmə)
- Replay attack — yazılmış blokun yenidən göndərilməsi
- Nonce — bir-dəfəlik dəyər; gizli DEYİL
- CTR mode — counter əsaslı mode
- CBC mode — əvvəlki ciphertext zənciri; "blockchain" mənşəyi
- IV (initialization vector) — ilk blokun random nonce-u; açıq göndərilir
- Digest — mesajın fiks uzunluqlu nümayəndəsi
- Hash table/bucket/clustering — lookup struktur/düyün/yığılma
- Preimage — hədəf hash-i verən mesaj; tapılması ÇƏTİN olmalıdır
- Avalanche effect — kiçik giriş → böyük çıxış dəyişikliyi
- Rainbow table — əvvəlcədən hash-lənmiş parol lüğəti
- Salt — hər parola unikal random qatqı
- bcrypt/scrypt — bilərəkdən yavaş parol hash-ləri
- Eventually consistent — paylanmış ledgerin nəhayət uzaqlaşması
- Double-spending — eyni pulun 2 dəfə xərclənməsi
- Doomed stub — konsensusda itirən zəncir qolu
- 51% attack — hashrate yarı-əksəriyyətinin ələ keçirilməsi
- Proof of work — hash < target şərti; enerji israfı
- Proof of stake — sərvət əsaslı etibar

## Praktik nəticə

1. **ECB heç vaxt:** struktur saxlanır (şəkil konturları!); CBC və ya
   CTR — öz mode YAZMA, cipher.NewCBCEncrypter işlət.
2. **CBC pipeline formatı:** IV + ciphertext; decipher-də ilk blok = IV;
   Pad/Unpad qalır.
3. **Mesaj bütövlüyü = ayrıca problem:** CBC Mallory-ni AŞKARLADIR amma
   QARŞISINI ALMIR — hash/MAC lazım (Ch17-də HMAC).
4. **Parol saxlama:** SHA-256 YOX; bcrypt/scrypt + unikal salt;
   rainbow-table-ə qarşı YALNIZ salt xilas edir.
5. **Hash seçimi:** MD5/SHA-1 MİLLİ; SHA-256 standart; parol vs digest
   ehtiyacları ƏKS-dir (sürət vs yavaşlıq).
6. **Zəncir = sıra doğrulaması:** hash(predecessor + data) idiomu —
   blockchain-dən logging-ə qədər universaldır.

## Mənbə
Pages: 226-277 (PDF 227-278)
