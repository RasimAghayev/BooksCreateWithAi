# Chapters 5-7 — Keys, Cribs, Passwords (səh. 85-126)

## Bu fəsillər nədən bəhs edir?

Keyspace anlayışı (hər əlavə bayt = 256x böyük açar sahəsi; 32 bayt/256 bit
= kainatdakı atomlar ~10^77; supernova enerjisi — 256-bit açarı brute-force
etmək 137 milyard supernova = qalaktika yandırmaq; "Galactic Emperor Eve"),
multi-byte açarlar (key[i%len(key)] — modul operatoru ilə açar təkrarı),
Encipher/Decipher []byte açar imzalarına keçidi (kompilyasiya xətalarının
tədricən həlli), test hallarının düzgünlüyü (all-zero plaintext = zəif detektor
— 0,1,2 qarışığı), hex açar flag-i (encoding/hex.DecodeString; DEADBEEF),
crack-in multi-byte versiyası (bayt-bayt guess; MaxKeyLen=32; min(32,
len(ciphertext)); slices müqayisəsi = bytes.Equal), crib uzunluq qaydası
(crib ≥ açar uzunluğu — "The" kifayət etmədi, "The " (4 bayt) bəs etdi;
yanlış deciphering təhlükəsi), parollar: dictionary attack (40% hesab;
25 söz = 10% parol), obfuscations (s3cr3t), passphrases ("Colourless green
ideas sleep furiously"; şeir baş hərfləri SIctasd?Tamlamt), password
POLİTİKALARI = zərər (rəqəm tələbi keyspace-i 2x AZALDIR; hər məhdudiyyət
Eve-ə MƏLUMATDIR), keyspace = K^N, Unicode (150k simvol; amma proqram
dəstəyi buggy), minimum 14 simvol + paste icazəsi; random anlayışı
(determinizm problemi; predictable sequences; coin toss = bit; məlumat
bitlərlə — GO < 10 bit; bilik = gözəçıxaran alt setlər "1111..." = 1 bit).

## Əsas fikirlər

### 1. Keyspace Riyaziyyatı (Ch5)
- **Hər əlavə bayt:** keyspace × 256 → orta crack vaxtı × 256 (EKSZONENSIAL:
  7 simvol = 6 dəqiqə; 14 simvol = 200 milyon İL)
- **256-bit (32 bayt) = praktik maksimum:** 2^256 ≈ 10^77 = görünən
  kainatın atomları; brute-force enerjiyə görə MÜMKÜNSÜZ (Schneier: 2^256
  sayğac = 137 milyard supernova enerjisi; qara dəlik təhlükəsi — kompüter
  çox sıx olarsa)
- **Nəticə:** uzun açarlarda brute-force QİYMƏTLİ hücum deyil — açarı OĞRURLAMAQ
  ucuzdur ($5 wrench attack)
- **Moore qanunu qarşı:** Eve gücü 2x artırsa — Alice 1 bit əlavə edir; açar
  uzatma ƏN ucuz təkmilləşmədir

### 2. Multi-Byte Açar — % Operatoru (Ch5)
```go
func Encipher(plaintext []byte, key []byte) (ciphertext []byte) {
    ciphertext = make([]byte, len(plaintext))
    for i, b := range plaintext {
        ciphertext[i] = b + key[i%len(key)]   // açarı DÖVRİ istifadə et
    }
    return ciphertext
}
func Decipher(ciphertext []byte, key []byte) (plaintext []byte) {
    // eynisi, amma minus: plaintext[i] = b - key[i%len(key)]
}
```
- **Modul (clock arithmetic):** 17 mod 12 = 5 — saat 12-dən "wrap"; Go-da `%`
- **Təkamül yolu:** magic := i%len(key) — "index out of range" xətasından
  → məhdud indeksə
- **Decipher artıq Encipher(-k) DEYİL:** multi-byte-də mənfi açar mənasızdı
  → loop-vari versiya qayıdır
- **Nəticə:** eyni plaintext hərfi artıq FƏRLİ şifrə hərfləri verir ("ll" → "[M")
  — frequency analysis çətinləşir (tam öldürmür)

### 3. Təkmilləşmiş Test Halları (Ch5)
```go
// ZƏİF hall (bug maskası):
{key: []byte{1,2,3}, plaintext: []byte{0,0,0}, ciphertext: []byte{1,2,3}}
// "ciphertext = key birbaşa" bugu KEÇİR (plaintext 0 olduğundan)!

// GÜCLÜ hall:
{key: []byte{1,2}, plaintext: []byte{0,1,2}, ciphertext: []byte{1,3,3}}
// plaintext-in İSTİRAKİ + açarın TƏKRARI hər ikisi yoxlanılır
```
- **Dərs:** test halı təsadüfi keçə bilər — hər uydurma bug-un testdən
  DÜŞDÜYÜNÜ təsəvvür et (bug-oriented test design)

### 4. Hex Açar Bayrağı (Ch5)
```go
keyHex := flag.String("key", "01", "key in hexadecimal (for example 'FF')")
key, err := hex.DecodeString(*keyHex)   // "DEADBEEF" → 4 bayt
```
- **Niyə hex:** 32-bayt açar int64-a SİĞMIR (8 bayt); 1 bayt = tam 2 hex
  rəqəm; base-16 (0-F; 0x3B = 59)
- flag.IntSlice YOXDUR — string + decode pattern; decode xətası mümkün →
  err yoxla

### 5. Multi-Byte Crack (Ch6)
```go
const MaxKeyLen = 32
func Crack(ciphertext, crib []byte) (key []byte, err error) {
    for k := range min(MaxKeyLen, len(ciphertext)) {  // 2 limit!
        for guess := range 256 {
            result := ciphertext[k] - byte(guess)
            if result == crib[k] {
                key = append(key, byte(guess))
                break
            }
        }
        if bytes.Equal(crib, Decipher(ciphertext[:len(crib)], key)) {
            return key, nil
        }
    }
    return nil, errors.New("no key found")
}
```
- **İki uzunluq limiti:** MaxKeyLen=32 (praktik maksimum) VƏ len(ciphertext)
  (uzun açar = sadəcə israf — şifrəmətninə təsir etmir)
- **Bayt-bayt hücum:** hər açar baytını ayrıca guess et (plaintext[k] +
  key[k] = ciphertext[k] → guess = ciphertext[k] − crib[k])
- **Bitmə yoxlaması:** hər baytdan sonra TAM crib-i decipher et — tutdusa qayıt
- **Slice müqayisəsi:** `!=` İŞLƏMİR (slice yalnız nil ilə müqayisə olunar)
  → bytes.Equal (və ya slices.Equal — ümumi)
- **Köhnə error-test SİLİNDİ:** istənilən ciphertext+crib üçün HƏMİŞƏ bir açar
  var (çoxbaytlı dünyada) — dəyərsiz test = texniki borc; "keep them lean"

### 6. Crib Uzunluq Qaydası (Ch6)
- **QAYDA: crib ≥ açar uzunluğu** (bu scheme-də)
- "The" (3 bayt) + DEADBEEF (4 bayt) → yanlış qısa açar DE AD BE tapılır →
  decipher DÜZGÜN DEYİL (ilk 3 hərf düz, qalanı zibil)
- "The " (4 bayt, space DAHİL) → düzgün
- **Niyə:** crib-dən qısа açar yalnız crib özünü bərpa edə bilər; qısa
  crib tam açar haqqında TAM MƏLUMAT vermir (qismən azaldır — search space
  kiçilir, amma kifayət deyil)

### 7. Parollar — Növlər və Hücumlar (Ch7)
| Növ | Nümunə | Risk |
|---|---|---|
| Dictionary word | "abdominoplasty" | 40% hesablar sözlə tapılır |
| Ən çox istifadə (25 söz) | password, qwerty, iloveyou, 123456 | 10% BÜTÜN parollar |
| Obfuscation | s3cr3t-passw0rd | Eve-nin 3-cü siyahısı; xatırlatma ÇƏTİN |
| Passphrase | "Colourless green ideas sleep furiously" | uzun (38) + yaddaşda qalır |
| Baş hərflər | SIctasd?Tamlamt (Şekspir) | uzunluq limiti olan sistemlər üçün |
| Random + manager | 32 simvol | ƏN YAXŞI; password manager istifadə et |

- **Eve-in axtarış sırası:** ən çox istifadə → dictionary → obfuscations →
  cüt sözlər (60 milyard — amma keyspace yanında "peanuts") → üçlü → brute

### 8. Password Politikalası = ZƏRƏR (Ch7)
- **Riyaziyyat:** 5-simvol, 100 simvol set = 10^10; "ən azı 1 rəqəm" tələbi:
  Eve 4-simvol keyspace (10^8) × 50 variant = 5×10^9 — YARI keyspace!
  "1 böyük hərf" daha da azaldır — HƏR məhdudiyyət Eve-ə MƏLUMAT verir
- **K^N düsturu:** K (set ölçüsü) yarımdursa keyspace 2^N azalır
- **Müsbət istisna:** MINİMUM uzunluq (14) — Eve qısa parolları atlayır, amma
  keyspace ekszponensial BÖYÜYÜR — xalis fayda
- **Regular dəyişmə = fəlakət:** güclü parol dəyişməməli; zəif dəyişmək
  kömək etmir; yeganə effekt: monitora yapışdırılmış qeyd
- **Duzmlu praktika:** dictionary lookup (yoxdursa qəbul et) + güc göstəricisi
  + minimum 14 + paste İCAZƏSİ (manager üçün)
- **Unicode:** 150k simvol → 8×10^25 (5 simvolda!) — amma REAL dünyada
  buggy dəstək (ASCII fərziyyəsi) → etibarsız; uzunluq daha etibarlı alət
- **Uzunluq TRUNCATION = cinayət:** ilk N simvolun saxlanması (köhnə Unix
  parolları) — uzun parol görünür, AMMA deyil

### 9. Random — Nə Deməkdir? (Ch7)
```go
func random() int { return 7 }   // "random" DEYİL — çünki SEQUENCE prediktabldır
```
- **Random = sequence xüsusiyyəti:** 17,86,2,17... / 1,2,4,8... / 2,3,5,7,11...
  (primes!) / 0,1,1,2,3,5 (Fibonacci!) — hamısı GÜCLÜ pattern; proqramlar
  deterministikdir → paradox
- **Coin toss = bit:** 5 toss = 32 ehtimal ≥ 26 hərf; 256-bit açar = 256 toss;
  3 toss = 8 → həftənin günü / göy qurşağı rəngi
- **Məlumat = bit sayı:** "GO" — 26²=676 kombinasiya < 2^10 → GO < 10 bit;
  artıq bitlər keyspace BÖYÜTMÜR
- **Bilik gözəçıxar:** "1111...1" (16 bit görünür) — Eve bilirsə Alice yalnız
  all-0/all-1 seçir → CƏMİ 1 bit!; həmişə all-1 → 0 bit (məlumat YOX)
- **Entropi = BİLMƏDİYİMİZ** — gözəçıxar alt setlər entropini öldürür
  (Ch12-də davam)

## Əsas terminlər
- Keyspace — bütün mümkün açarların çoxluğu
- Weak key — mətni dəyişməyən açar (all-zero)
- Modular reduction (%) — qalıq; "clock arithmetic"
- Dictionary attack — sözlər siyahısı ilə parol axtarışı
- Obfuscation — s3cr3t tipli əvəzetmələr
- Passphrase — çoxsözlü yaddaşda qalan parol
- Password policy — simvol tələbləri (adətən ZƏRƏRLİ)
- K^N — N simvol, K ehtimal/simvol → keyspace düsturu
- Truncation — parolun kəsilməsi (ilk N simvol)
- Deterministic program — eyni giriş → eyni nəticə (random problemi)
- Coin toss scheme — bit-əsaslı təsadüfi seçim
- Bit (məlumat vahidi) — unikal identifikasiya üçün binary seçim
- Bug-oriented test design — uydurma bug-un testdən düşməsini təsəvvür et

## Praktik nəticə

1. **Açar uzunluğu = birinci müdafiə xətti:** 256 bit praktik limit; hər bit
   Eve üçün 2x iş; brute-force-dan daha ucuz hücumlar HƏMİŞƏ var (oğurluq).
2. **[]byte açar API dizaynı:** single-byte → multi-byte keçidində % len(key)
   idiomu; testləri QARİŞIQ plaintext-lərlə yoxla (all-zero maskaları).
3. **Crib ≥ açar uzunluğu:** qısa crib yanlış (qısa) açar qaytara bilər;
   tam yoxlama hər baytdan sonra.
4. **Parol sistemləri üçün ƏMR:** minimum 14; dictionary reject; paste icazəsi;
   simvol TƏLƏBLƏRİ YOX; regular dəyişmə YOX; truncation YOX.
5. **Random üçün:** deterministik proqram təbii random DEYİL — entropiya
   mənbəyi lazımdır (Ch13-də hardware); sequence pattern-ə GÖRƏ hakim.

## Mənbə
Pages: 85-126 (PDF 86-127)
