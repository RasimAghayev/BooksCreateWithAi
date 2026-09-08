# Chapters 1-4 — Ciphers, Enciphering, Deciphering, Cracking (səh. 20-84)

## Bu fəsillər nədən bəhs edir?

Kitabın TDD ilə qurulan bünövrəsi: codes vs ciphers fərqi (Morse = ötürmə
effektivliyi; "TOMATO/PARMESAN" kod kitabı = saxlanma problemi), shift
(Caesar) cipher — hər baytı +k ilə əvəz et; []byte API (arbitrary data);
behaviour-cümlədən-testə axını ("Encipher transforms HAL to IBM" →
TestEncipherTransformsHALToIBM); smoke-detector metaforası — testi SİNDİRMƏDƏN
əvvəl məqsədli bug qoy (null implementation); smarty-pants həlli ("return
[]byte("IBM")") → daha böyük test lazımdır; table test (testCase struct,
t.Run subtest adları "HAL to IBM"); encipher/decipher/crack CLI alətləri
(filter pattern: stdin→stdout), frequency analysis (Ə hərfi = ən çox F →
key=1), variable key (flag.Int, byte(*key)), 255 faydalı açar (0 = weak),
Decipher = Encipher(ciphertext, -key) refaktoru, ortaq test cədvəli (var
cases package səviyyəsində — plaintext/ciphertext sahə adları), brute-force
Crack (crib ilə: for guess := range 256 {Decipher + müqayisə }), errors.New
("no key found"), magic number cribs (PNG = 0x89 + "PNG" — printf '\x89PNG'
shell escape ilə), binary faylların (şəkil!) crack-i.

## Əsas fikirlər

### 1. Codes vs Ciphers (Ch1)
- **Code:** söz/ibarənin başqa təyinatla əvəzi ("tomato" = "enemy in
  sight"); kod KİTABI böyük + təhlükəli; məqsəd gizlilik olmaya bilər
  (Morse = bandwith qənaəti)
- **Cipher:** hərf/bayt transformasiyası (rəqəm→rəqəm); açar (key) + düz
  mətn → şifrə mətni;**Alice/Bob/Eve** (göndərən/qəbul edən/dinləyici)
- **Shift cipher:** hər bayta +k əlavə; decoder ring metaforası; Caesar

### 2. Behaviour → Test Kanonu (Ch1)
```
"Encipher transforms HAL to IBM."
    ↓ (cümlə = test adı)
TestEncipherTransformsHALToIBM
    input []byte("HAL") → want []byte("IBM") → got := Encipher(input)
    bytes.Equal(want, got) → t.Errorf("want %q, got %q")
```
- **Spherical cow texniqi:** mürəkkəb problemi sadələşdir — key=1, data
  nəqli YOX; sadə işləyəndə realistik et
- **Test = bug detektoru:** işə salmazdan əvvəl SİNDİR — null Encipher
  (return nil) qoy → fail göstər → yalnız sonra düzəlt; "test keçir" =
  detektor İŞLƏYİR demək deyil, yalnız "bu bugu tutur" deyir
- **GOAL → HINT → SOLUTION:** kitabın TDD formatı — hər addımda oxucu
  əvvəl ÖZÜ cəhd edir

### 3. Table Test İnkişafı (Ch1)
```go
tcs := []struct {
    key         byte
    input, want []byte
}{
    {key: 1, input: []byte("HAL"), want: []byte("IBM")},
    {key: 3, input: []byte("PERK"), want: []byte("SHUN")},
    {key: 7, input: []byte("CHEER"), want: []byte("JOLLY")},
    {input: []byte{0, 1, 2, 3, 255}, want: []byte{1, 2, 3, 4, 0}},  // wrap!
}
for _, tc := range tcs {
    name := fmt.Sprintf("%s to %s", tc.input, tc.want)
    t.Run(name, func(t *testing.T) { ... })   // subtest = aydın FAIL adı
}
```
- **Smarty-pants dərsi:** `return []byte("IBM")` testi KEÇİR — testlərin
  zəifliyi; yeni hal ("ADD"→"BEE") əlavə et → brain-dead həllər istisna
- **255 wrap halı:** bayt 255+1 = 0 — sərhəd şərtləri testdə MÜTLƏQ

### 4. Filter Proqramı (Ch2)
```go
// cmd/encipher/main.go — Unix filter: stdin → stdout
func main() {
    key := flag.Int("key", 1, "shift value")   // adlı flag - aydın UX
    flag.Parse()
    plaintext, err := io.ReadAll(os.Stdin)
    if err != nil { fmt.Fprintln(os.Stderr, err); os.Exit(1) }
    ciphertext := shift.Encipher(plaintext, byte(*key))
    os.Stdout.Write(ciphertext)
}
```
- `echo "HEY BOB U UP?" | go run ./cmd/encipher -key 5`
- **flag.Int → *int pointer:** flag.Parse() sonra dəyər içində; byte(*key)
  çevirməsi (flag.Byte YOXDUR)

### 5. Frequency Analysis — Eve (Ch2-3)
- **Hücum:** ən çox rast gəlinən şifrə hərfi = E (ingiliscə ən çox);
  fərq = açar; T/A/I/N/O/S növbəti tutumlar; dil cədvəlləri fərqli amma
  oxşar
- **Açar zəifliyi kök deyil — determinizm:** eyni düz hərf → həmişə eyni
  şifrə hərfi = statistik barmaq izi
- **Müalicə istiqamətləri:** açar sayını artır; bir açarla encipher
  olunan data miqdarını AZALT
- **Shannon:** "the enemy knows the system" — scheme-i gizli SAXLA, amma
  onun bilinməsi sistemin iflası OLMAMALI (Kerckhoffs prinsipi)

### 6. Brute-Force Riyaziyyatı (Ch3)
- 3-rəqəmli seyf: 1000 kombinasiya, ~1 saniyə/cəhd → 16 dəqiqə TAM,
  8 dəqiqə 50% ehtimal — orta = yarı yolu
- Bayt açarı: 256 dəyər (0 = weak → 255 faydalı) → bir kelibayt = BUGÜNKÜ
  kompüterlə üçün ANLIQ crack

### 7. Decipher — Refaktor Dərsi (Ch3)
```go
// v1 (işləyir, amma təkrar):
func Decipher(ciphertext []byte, key byte) (plaintext []byte) {
    plaintext = make([]byte, len(ciphertext))
    for i, b := range ciphertext { plaintext[i] = b - key }
    return plaintext
}
// v2 (enciphering = mənfi açarla deciphering):
func Decipher(ciphertext []byte, key byte) []byte {
    return Encipher(ciphertext, -key)
}
```
- **Test implementation-a bağlı DEYİL:** bugünkü Decipher = Encipher(-k);
  sabah döngü ilə — hər ikisi DAVRANIŞ testini keçməli
- **Ortaq test cədvəli:** package-level `var cases` — plaintext/ciphertext
  sahə adları (input/want YOX — istiqamətdən asılı); TestEncipher və
  TestDecipher HƏR İKİSİ eyni cədvəldən oxuyur
- **Duplikasiya > komplikasiya:** encipher/decipher ayrı proqramlar — bir
  "hər şeyi edən" alət YOX ("encipher -decipher" ikili mənalı UX-dən pis)

### 8. Crack — Brute Force + Cribs (Ch4)
```go
func Crack(ciphertext, crib []byte) (key byte, err error) {
    for guess := range 256 {
        result := Decipher(ciphertext[:len(crib)], byte(guess))
        if bytes.Equal(result, crib) {
            return byte(guess), nil
        }
    }
    return 0, errors.New("no key found")   // 0 qaytarmaq Go-like DEYİL
}
```
- **Crib:** bilinən düz-mətn fraqmenti (mesajın başlanğıcı); "GO" →
  ilk 2 baytı yoxla; TAM mesajı decipher etmək İSRAFDIR — yalnız
  len(crib) qədər
- **Stereotip = crib:** Go faylı → "package"; PNG → 0x89+"PNG" (magic
  number/file header)
- **Error halı:** heç bir açar tutmadısa (başqa scheme!) → error; test:
  TestCrackReturnsErrorWhenKeyNotFound
- **CLI:** crack -crib 'The tiger' <enciphered.bin → plaintext; PNG
  üçün: `printf '\x89PNG'` shell escape $(...) substitusiya ilə

### 9. Binary Data Crack (Ch4)
- Şəkil encipher et (devil.png → devil.bin) → crack magic-number crib ilə
  → bayt-bayt ORİJİNAL fayl bərpa olunur — API []byte olduğu üçün hər
  format AVTOMATİK işləyir

## Əsas terminlər
- Code vs cipher — söz-əvəzi vs hərf-transformasiya
- Plaintext/ciphertext/key — düz mətn / şifrə mətni / açar
- Alice/Bob/Eve — göndərən/qəbul edən/dinləyici
- Shift (Caesar) cipher — +k bayt sürüşdürmə
- Weak key — key=0 (mətni dəyişmir)
- Frequency analysis — hərf tezliyi statistikası ilə crack
- Brute force — bütün açarların sınanması
- Crib — bilinən fraqment ilə açar yoxlaması
- Magic number — fayl formatının başlıq baytları (0x89 PNG)
- Behaviour sentence → test — "X transforms Y to Z" cümləsinin testə çevrilməsi
- Null implementation — qəsdən səhv (return nil) — detektor sınağı
- Table test/subtest — cədvəl + t.Run adlı hallar
- Stochastic debugging — təsadüfi kod dəyişmə (ANTİ-PATTERN)
- Shannon maxim — "enemy knows the system"

## Praktik nəticə

1. **TDD ardıcıllığı (kitabın kanonu):** behaviour cümləsi → test → testi
   SİNDİR (null impl) → işləyən minimal həll → refaktor (test yaşıl qalsın)
2. **Hər yeni funksiya imzası üçün:** kompilyasiya xətasını daxil et (test
   "xəyali" çağırış), sonra imzanı düzəlt, sonra məntiqi — "too many
   arguments" messajlarını OXU (have/want formatı)
3. **Ortaq test cədvəli:** plaintext/ciphertext adları ilə package-level
   var — Encipher/Decipher/Crack hamısı eyni həqiqətdən yoxlanılır
4. **Filter CLI idiomu:** stdin→ReadAll→transform→stdout; flag adı =
   məqsəd (encipher 5 YOX, encipher -key 5)
5. **Deterministik cipher statistik üçün VULNERABLDIR:** uzun mesaj,
  eyni açar = Eve qazanır; həll gələcək fəsillərdə (bloklar/zəncirlər)
6. **Sərhəd şərtlərini testə sal:** 255→0 wrap; boş crib; tapılmayan açar
  → error (magic 0 YOX)

## Mənbə
Pages: 20-84 (PDF 21-85)
