# Chapters 8-10 — Blocks, Modes, Padding (səh. 127-177)

## Bu fəsillər nədən bəhs edir?

Stream vs block cipher (bit-bit/time-varying vs bloklar/fixed; hardware vs
software), blokların səbəbi (10GiB fayl = io.ReadAll panic — yaddaş; blok =
hissə-hissə emal), `cipher.Block` interfeysi (BlockSize/Encrypt/Decrypt —
"plug and socket" konsepsiya; Encrypt nəticəni QAYTARMIR, dst-ə YAZIR; açar
parametr YOX — receiver struct-da!), shiftCipher-in blok versiyası (key
[32]byte fixed; BlockSize=32), açar genişləndirməsinin MÜQAVİLƏSİ (SHA-256
ilə qısa açarı uzatmaq = entropi DƏYİŞMİR — "böyük hərflərlə yazılan parol";
tough love: tam 256-bit tələb et), NewCipher konstruktoru (sentinel
ErrKeySize + fmt.Errorf "%w" wrap + errors.Is testi — string yoxlamadan
ROBUST), `cipher.BlockMode` interfeysi (CryptBlocks; encrypter/decrypter
TIPLƏRİ — Encrypt/Decrypt-in əvəzinə TƏK CryptBlocks), blok-döngüsü
(`for len(src) > 0 { Encrypt; src = src[n:]; dst = dst[n:] }` — re-slice
progressi), block-aligned olmayan data → panic ("input not full blocks" —
interface imzada error YOXDUR → proqramçı səhvi = panic konvensiyası; dst <
src → "output smaller than input"), panic testləmənin MÜMKÜNSÜZLÜYÜ
(recover ilə test yazıla bilər amma mənasızdır — hər iki halda panic var,
test nəyi yoxlayır?), PKCS#7 padding (həmişə pad et — "illegal pixel"
problemindən qaçınmaq; N qısa = N dənə N baytı; tam blok = TAM BLOK padding
[32×32]; son bayt = sayğac), Pad/Unpad TDD (4-baytlıq kiçik blockSize ilə
test; bytes.Repeat; Unpad = son baytın dəyəri qədər kəs), RFC 2315 (block
size < 256 məhdudiyyəti — 1 bayt sayğac), tam encipher/decipher pipeline
(Pad → CryptBlocks → Unpad), decrypter = encrypterin surəti (Decrypt çağırışı
istisna — "duplication is better than complication").

## Əsas fikirlər

### 1. Stream vs Block Cipher (Ch8)
|  | Stream | Block |
|---|---|---|
| Vahid | 1 bit | N bayt (blok) |
| Transformasiya | time-varying | fixed |
| Yer | hardware, telekom dövrələri | software, arbitrary mesajlar |
- **Shift cipher:** faktiki "mesaj boyu tək blok" — 10GiB faylda
  io.ReadAll → OOM panic; blok-blok emal = yaddaş dostu
- **Praktik blok ölçüsü:** ~32 bayt çox ciphers üçün

### 2. cipher.Block İnterfeysi (Ch8)
```go
type Block interface {
    BlockSize() int
    Encrypt(dst, src []byte)   // dst-yə YAZ — qaytarmır!
    Decrypt(dst, src []byte)
}
```
- **Kontrakt fərqləri Encipher-dən:** nəticə RETURN YOX (dst parametr);
  AÇAR parametr YOX — struct sahəsində (receiver c)
- **İnterfeys = socket:** istənilən cipher adapter ilə PLUG-IN olur —
  standart kitabxana konneptoru (io.Reader/Writer kimi)

### 3. shiftCipher — Fiks Açar Ölçüsü (Ch8)
```go
const BlockSize = 32
type shiftCipher struct {
    key [BlockSize]byte     // []byte YOX — FIXED 32
}
func (c *shiftCipher) Encrypt(dst, src []byte) {
    for i, b := range src {
        dst[i] = b + c.key[i]    // açar artıq receiver-də
    }
}
```
- **Niyə fixed 32:** minimum uzunluq MƏCBURİ (zəif qısa açarlardan qoruma);
  32 = praktik maksimum, artıq lazımsız
- **Açar genişləndirmə TƏLƏSİK:** SHA-256 ilə 5 hərfi 32 bayta ÇEKMƏK
  entropi ƏLAVƏ ETMİR — Eve eyni hash-i tətbiq edib 26^5-i sınayır;
  "nəhəng hərflə parol" — YERİ YOXDUR; tam 256-bit istə (tough love)

### 4. NewCipher + Sentinel Error Wrap (Ch8)
```go
var ErrKeySize = errors.New("shift: invalid key size")
func NewCipher(key []byte) (cipher.Block, error) {
    if len(key) != BlockSize {
        return nil, fmt.Errorf("%w %d (must be %d)", ErrKeySize, len(key), BlockSize)
    }
    return &shiftCipher{key: [BlockSize]byte(key)}, nil
}
```
- **Üç səviyyəli test təkamülü:** (1) err != nil yoxlaması; (2) err ==
  ErrKeySize (sentinel); (3) **errors.Is(err, ErrKeySize)** — wrap-ə davamlı!
- **Niyə interfeys tipi qaytarır:** `cipher.Block` (concrete *shiftCipher
  YOX) — (a) niyyət siqnalı; (b) interfeysi POZAN dəyişiklik KOMPİLYASİYA
  xətası verir ("missing method BlockSize" — dərhal tutulur!)
- **Mütləq array çevirmə:** [BlockSize]byte(key) — slice→array (uzunluq
  yoxlanılıb artıq)

### 5. Test Açarı İnkşafı (Ch8)
```go
var testKey = bytes.Repeat([]byte{1}, shift.BlockSize)  // 32× 0x01
var cipherCases = []struct{ plaintext, ciphertext []byte }{
    {plaintext: []byte{0, 1, 2, 3, 4, 5}, ciphertext: []byte{1, 2, 3, 4, 5, 6}},
}
```
- 32 hərflək açarı hər halda yazmaq əziyyət — TEK açar + Repeat; qısa
  mesajlar kifayət (açarın ilk N baytı istifadə olunur)
- Köhnə Encipher/Decipher/Crack testləri SİLİNDİ (əvəz olundular) — ölü
  kod = texniki borc

### 6. cipher.BlockMode İnterfeysi (Ch9)
```go
type BlockMode interface {
    BlockSize() int
    CryptBlocks(dst, src []byte)   // İSTƏNILƏN uzunluq; Encrypt/Decrypt YOX
}
```
- **DİZAYN:** mode = "hər hansı Block ver, mən blok-blok işləyərəm" —
  cipher-dən MÜSTƏQİL (istənilən blok cipher ilə işləyir)
- **İki tip:** encrypter (Encrypt çağırır) + decrypter (Decrypt) —
  interfeys tək metod, istiqamət tipə bağlı
- **Blok döngüsü kanonu:**
```go
type encrypter struct {
    block     cipher.Block
    blockSize int
}
func NewEncrypter(block cipher.Block) cipher.BlockMode {
    return &encrypter{block: block, blockSize: block.BlockSize()}  // SORUŞ!
}
func (e *encrypter) CryptBlocks(dst, src []byte) {
    if len(src)%e.blockSize != 0 { panic("encrypter: input not full blocks") }
    if len(dst) < len(src) { panic("encrypter: output smaller than input") }
    for len(src) > 0 {
        e.block.Encrypt(dst[:e.blockSize], src[:e.blockSize])
        src = src[e.blockSize:]   // re-slice ilə progress
        dst = dst[e.blockSize:]
    }
}
```
- **blockSize SORUŞULUR, hardwire YOX:** hər hansı cipher üçün işləməli —
  block.BlockSize() çağırışı

### 7. Panic Konvensiyası (Ch9)
- **Niyə panic, error YOX:** interface imzası FIXED (CryptBlocks error
  qaytarmır) — block-aligned olmayan giriş = PROQRAMÇI səhvi (istifadəçi
  yox) → Go konvensiyası: panic
- **Faydalı panic mesajı:** "slice bounds out of range [:32] with capacity
  31" YOX → "encrypter: input not full blocks" — SƏBƏBİ de
- **Panic testlənməsi MÜMKÜNSÜZDÜR (praktiki):** recover() ilə "panic gözlə"
  testi yazmaq olar, amma yoxlama silinsə BELƏ panic baş verir (yeni
  mesajla) → test nəyi yoxlayır? → panic-ləri adətən TEST ETMƏ; real
  istifadə özü aşkar edəcək

### 8. PKCS#7 Padding (Ch10)
**Problem:** zero-padding → son sıfırlar data-dırmı, padding-dirmi?
AYRI DEDİRMƏK OLMAZ ("illegal pixel" — "şəkilizin sağ alt pikseli qara,
QADAĞANDIR"!)
**Həll — HƏMİŞƏ pad et:**
```go
func Pad(data []byte, blockSize int) []byte {
    n := blockSize - len(data)%blockSize      // 1..blockSize (0 YOX!)
    padding := bytes.Repeat([]byte{byte(n)}, n)
    return append(data, padding...)
}
func Unpad(data []byte, blockSize int) []byte {
    n := int(data[len(data)-1])               // son bayt = sayğac
    return data[:len(data)-n]
}
```
- **Qaydalar:** N qısa → N dənə N; TAM blok → TAM BLOK (32×32) — 3% halda
  "israf", amma ambiguity YOXDUR; boş data → tam blok padding
- **RFC 2315 məhdudiyyəti:** blockSize < 256 ZƏRURİ (sayğac 1 baytdır);
  255+ bloklar üçün başqa scheme
- **Test dizaynı:** 4-baytlıq blockSize (32 yazmaq əziyyəti YOX); maraqlı
  hallar: 1-qısa/2-qısa/3-qısa/tam/boş; padCases raw/padded sahə adları
  (want/got YOX — hər iki istiqamətdə istifadə)

### 9. Tam Pipeline (Ch10)
```go
// encipher:
key → NewCipher → NewEncrypter → ReadAll → Pad(bs) → CryptBlocks → stdout
// decipher:
key → NewCipher → NewDecrypter → ReadAll → CryptBlocks → Unpad(bs) → stdout
```
- Yoxlama: 468 baytlıq tiger → 480 (12 əlavə = 0D 0D ... hexdump-da GÖRÜNÜR!)
- **Decrypter = Encrypter-in surəti:** Encrypt→Decrypt fərqi; "duplication
  is better than complication" — sadə kodu ABSTRAKT etməyin

## Əsas terminlər
- Stream cipher — bit-bit, time-varying (hardware/telekom)
- Block cipher — fiks bloklar, fixed transformasiya (software)
- cipher.Block / cipher.BlockMode — standart interfeyslər (plug/socket)
- Mode of operation — ardıcıl blokların encipher qaydası
- Block-aligned — uzunluq blockSize-un qatları
- Sentinel error — adlandırılmış sabit error (ErrKeySize)
- Error wrapping (%w) — sentinel + runtime məlumat; errors.Is ilə yoxlama
- Re-slicing progress — src = src[n:] ilə döngə
- PKCS#7 (RFC 2315) — N qısa → N×N padding; həmişə pad
- Illegal pixel problem — data/padding ayırd etməzlik ambiguity-si
- Tough love API — zəif girişi QƏBUL ETMƏYƏN dizayn
- Panic konvensiyası — proqramçı səhvi interface-də error yoxdursa

## Praktik nəticə

1. **cipher.Block implementasiya et:** struct (key sahəsi) + 3 metod;
   NewCipher → (cipher.Block, error) — INTERFEYS tipi qaytar (pozuntular
   compile-da tutulsun).
2. **Sentinel + wrap:** var Err... + fmt.Errorf("%w ...") + errors.Is —
   error STRING müqayisəsi heç vaxt.
3. **CryptBlocks qoruma kodu:** 2 panic yoxlaması (full blocks + dst≥src)
   — KÖMƏKLİ mesajlarla; panic-ləri test etməyə çalışma.
4. **PKCS#7:** həmişə pad; blockSize < 256; Unpad = son bayt sayğacı.
5. **Blok döngüsü:** for len(src) > 0 + Encrypt + re-slice — mənimsəməli
   idiom.
6. **Açar uzatma YALTAQLIĞI:** hash ilə "uzatmaq" entropi saxtalaşdırır —
   qısa açarları rədd et; istifadəçini düzgün alətə (KDF/password hashing)
   yönləndir (Ch15-də görünəcək).

## Mənbə
Pages: 127-177 (PDF 128-178)
