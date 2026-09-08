# Chapter 14 — Chains (səh. 226-247)

## Bu chapter nədən bəhs edir?

Uzun açar belə şifrəni gizlətmirsə nə olur: ECB (Electronic Code Book) modunun
zəiflikləri — pattern sızması, code book hücumu, Mallory-in replay/drop/redaktə
hücumları — və həll: CBC (Cipher Block Chaining) modu, IV anlayışı və
`crypto/cipher` ilə implementasiya.

## Əsas fikirlər

### 1. Məqsəd: ciphertext = təsadüfi şüma
**Qızıl standart:** Şifrələnmiş şəkil tam random noise-a bənzəməlidir —
piksel-lərin "hər şeyi" ehtiva edən kasıb görüntüsü.

**Təsadüfi şəkil generatoru:**
```go
for x := range width {
    for y := range height {
        img.Set(x, y, color.RGBA{
            R: uint8(rand.IntN(256)),
            G: uint8(rand.IntN(256)),
            B: uint8(rand.IntN(256)),
            A: 255,
        })
    }
}
```

### 2. Sınaq: ECB ilə şifrələnmiş devil gopher
**Metod:** PNG → PPM formatına çevir (header 3 sətir + payload ayrılır) →
`tail` ilə header-i ayır → payload-ı şifrələ → header-i geri birləşdir →
nəticəyə bax.

**Nəticə (aşağı açar, "01 01 01"):** Şəkil DEMƏK OLAR DƏYİŞMƏYİB — fon qara,
gözlər fırlanan, amma devil gopher aydın seçilir!

**Nəticə (32 baytlıq qvant açar, qrand ilə):** Rənglər qarışıb, amma kontur
hələ də görünür.

**Səbəb:** Açar təkrarlanır — 32 baytlıq açar ~megabaytlıq şəklin cəmi 0.05%-i
qədər məlumatı əhatə edir. Açarı uzatmaq həll deyil — **problem moddadir**.

### 3. ECB — "code book" modu
**Nədir:** Hər blok müstəqil şəkildə yalnız açar ilə birləşdirilir → **eyni
input blok = həmişə eyni output blok**.

**Zəifliklər:**
1. **Pattern sızması:**plaintext strukturu ciphertext-də qalır (devil konturu)
2. **Code book hücumu:** bütün input→output xəritələrini əvvəlcədən cədvəl
   kimi yaradıla bilər → Eve mesajı sadə lookup ilə oxuyur

### 4. Mallory-in ortadakı hücumları (active attacker)
**Replay attack (təkrar hücumu):** Mallory Alice-in "Alice-dən Bob-a 1000$
köçür" blokunu yaxalayır, sonra təkrar göndərir. ECB-də eyni mesaj = eyni
ciphertext → bank fərqi bilmir.

**Block dropping (blok silmə):** Mallory laptop alır, ödəniş blokunu tapıb
silir — Bob heç nə almır, Alice isə ödəniş getdiyinə əmindir.

**Block modification (blok redaktəsi):** Mallory blokların yerini dəyişə,
məzmununu redaktə edə bilər — CBC-də autentifikasiya olmadan bu qismən qalır
(göstərilir: bütöv mesaj replay-i CBC-də də mümkündür — ilk blok eyni
başlayanda bütün zəncir eyni nəticə verir).

### 5. CTR — counter modu (qısa baxış)
**İdea:** Hər blokla birləşdirilən **sayğaç (counter)** daxil edilir — hər blok
fərqli "spice" alır, code book problemi aradan qalxır.

### 6. CBC — zəncir modu (əsas həll)
**Nədir:** Hər blokun şifrələnməsinə **əvvəlki blokun CIPHERTEXT-i** daxil
olunur:

```
input_i = plaintext_i ⊕ ciphertext_(i-1)
ciphertext_i = E(key, input_i)
```

**Sub-kod izahı:**
- Eyni plaintext iki yerdə → fərqli əvvəlki ciphertext → fərqli nəticə → pattern sızması yox
- Sayğaç və ya nonce yaratmağa ehtiyac yoxdur — zəncir özü təchiz edir
- **Zəif qalan:** bütöv mesaj replay — bütövlükdə eyni mesaj eyni zənciri verir

### 7. IV — Initialization Vector (başlanğıc vektoru)
**Problem:** İlk blokun "əvvəlki ciphertext-i" yoxdur.

**Həll:** Bir blok qədər **təsadüfi junk** yaradılır — məzmunu əhəmiyyətsizdir,
yalnız zənciri başlatmalıdır. `crypto/rand` ilə:
```go
iv := make([]byte, shift.BlockSize)
_, err = rand.Read(iv)
```
IV ciphertext-in **birinci bloku kimi** yazılır → decipher tərəfi onu oxuyub
çıxarır. Mesajı replay etsələr belə, yeni IV yeni zəncir verir.

### 8. Implementasiya — crypto/cipher CBC
**Encipher:**
```go
block, err := shift.NewCipher(key)          // bizim blok şifrəmiz (cipher.Block)
// ... read plaintext ...
iv := make([]byte, shift.BlockSize)
rand.Read(iv)                                // təsadüfi IV
mode := cipher.NewCBCEncrypter(block, iv)    // CBC encrypter
padded := pad(plaintext, shift.BlockSize)    // blok həcminə tamamla
ciphertext := make([]byte, len(padded))
mode.CryptBlocks(ciphertext, padded)         // CBC ilə şifrələ
// yaz: iv || ciphertext  (480→512 bayt: +1 blok IV)
```

**Decipher:**
```go
ciphertext, _ := io.ReadAll(os.Stdin)
iv := ciphertext[:shift.BlockSize]           // ilk blok = IV
rest := ciphertext[shift.BlockSize:]         // qalan şifrələnmiş hissə
block, _ := shift.NewCipher(key)
mode := cipher.NewCBCDecrypter(block, iv)
plaintext := make([]byte, len(rest))
mode.CryptBlocks(plaintext, rest)            // CBC ilə aç
plaintext = unpad(plaintext)                 // pad-i çıxar
```

**Nəticə:** Şifrələnmiş devil gopher artıq random noise-a bənzəyir — CBC
uğurla işlədi.

## Əsas terminlər

- ECB (Electronic Code Book / elektron şifrə kitabı)
- CBC (Cipher Block Chaining / şifr blok zənciri)
- CTR (Counter Mode / sayğac modu)
- IV (Initialization Vector / başlanğıc vektor)
- Replay Attack (təkrar hücumu)
- Mallory (aktiv hücumçu — ortadakı adam)
- Code Book (blok xəritəsi)

## Praktik nəticə

- ECB heç vaxt istifadə etmə — pattern sızması, code book, replay hücumları
- CBC: hər blok əvvəlki ciphertext-ilə zəncirlənir → təkrarlanma ölür
- IV: hər mesaj üçün yeni təsadüfi blok, `crypto/rand`-dən; ciphertext-ə
  prefiks kimi yazılır (gizli saxlanmır!)
- `cipher.NewCBCEncrypter/NewCBCDecrypter` + `CryptBlocks` — mod kodu
  standart kitabxanadan gəlir, özün yazma
- CBC məxfiliyi (confidentiality) verir, amma **autentifikasiya yoxdur** —
  blok kəsilməsinə qarşı MAC/AEAD lazımdır (sonrakı chapter-lər)

## Mənbə

Pages: 226-247 (Chapter 14, Explore Go: Cryptography)
