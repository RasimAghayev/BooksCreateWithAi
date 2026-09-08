# Chapter 13 — Randomness (səh. 213-225)

## Bu chapter nədən bəhs edir?

Təsadüfülük necə yaradılır: pseudo-random generator-lar (PRNG — Fibonacci
nümunəsi, seed, periodicity, parite tələsi), çevrə şüması (environmental noise),
hardware TRNG, `math/rand` vs `crypto/rand` fərqi və kvant random generator
(ANU qrand). "Təhlükəsizlik nisbidir" tezisi.

## Əsas fikirlər

### 1. PRNG — psevdo-təsadüfi generator
**Nədir:** Deterministik proqram olub təsadüfü görünən nəticə verir ("pseudo" =
yalançı).

**Fibonacci nümunəsi:** İki seed ədədinin cəmi, mod 10:
```
seed: 1, 2 → 1, 2, 3, 5, 8, 3, 1, 4...  (son rəqəm mod 10)
```
**Sub-kod izahı:**
- Hər addım: əvvəlki iki ədədin cəmi → mod 10 (son rəqəm)
- "Artan" xarakter mod 10 ilə qırılır → təsadüfü görünür
- Seed-lər saatdan götürülür (məs. dəqiqə + saniyə)

**Nə üçün oyunlara yeter:** Eve seed-i bilmirsə, seriyaları proqnozlaşdıra
bilmir. Elite (BBC Micro, 1984) oyunu: partlayış hissəciklərinin koordinatları
hər dəfə eyni seed-lə **yenidən yaradılır** — yaddaş qənaəti + təkrarlanma
(repeatability) + təsadüfü görünüş. Gameplay-dan seed götürülür (məs. ən yaxın
gəminin ekran koordinatı) → hər oyun fərqlənir.

### 2. PRNG-nin ölümcül zəiflikləri
**Periodicity (dövrilik):** seed-lərdən asılı olaraq çox qısa tsiklərə düşür:
```
seed: 0, 5 → 0, 5, 5, 0 → 0, 5, 5, 0 ...  (4 rəqəmdən ibarət döngü!)
```

**Statistik tələs (parite):** seed 2 və 4 → `6, 0, 6, 6, 2, 8, 0, 8, 8, 6, 4`
— cüt + cüt = cüt → **heç vaxt tək rəqəm çıxmır**. Tsiklə düşməsə belə,
seriya bütün rəqəm fəzasının yarısında həmişəlik həbs qalır.

**Nəticə:** `math/rand` bu sinifdəndir (deterministik + vaxt seed) →
kriptoqrafiya üçün QADAĞANDIR. "Pseudo" sözü bunu xatırladır.

### 3. Environmental noise — həqiqi entropiya mənbəyi
**Nədir:** Kompüterin ətrafından fiziki hadisələr toplanır: mouse hərəkəti,
klaviatura vaxtı, disk şüması, network gecikmələri. Bunlar deterministik
olaraq modelləşdirilə bilsə də, praktik olaraq proqnozlaşdırılmazdır.

**İşləmə sxemi:** noise → deterministik "mixing" (qarışdırma) alqoritmi →
çıxış. Yəni "fake randomness" amma kifayət qədər təhlükəsiz.

### 4. Təhlükəsizlik nisbidir
**Nədir:** "Təhlükəsiz/qeyri-təhlükəsiz" ikiliyi yoxdur — hücumçunun sərf
edəcəyi səylə nisbi təhlükəsizlik var.

**Eve-in real hücum yolları (açıq yerdən keçənlər):**
- Malware inyeksiyası (şifrə açarını oxuyur — nə generator, nə şifrə əhəmiyyətli)
- Sosial mühəndislik ("parolu ADMIN edin")
- Fiziki təhdid (açarlı; wrench attack)
- Phishing meme linki
- Bob-un evinə girib çözülmüş mesajları oxumaq

**Ayı nağılı:** İki turist ayıdan qaçmalı — ayını deyil, bir-birini ötməli.
Kriptoqrafiya də: açarların "əbədi sınmaz" olması lazım deyil — hücumçunun
sərf etməyə hazır olduğu səydən bir az daha çox davamlı olmalıdır.

### 5. TRNG — hardware true RNG
**Nədir:** CPU-lara daxilindəki termal şüama (thermal noise) əsaslı həqiqi
təsadüfülük generatorları + oxumaq üçün machine instruction-lar.

- Yoxdursa: USB qurğu kimi xarici electronical/optical RNG alətləri
- "God Emperor Eve" məqamı: termal şüama belə, əgər Eve fərdi atom/elektron
  kinetik enerjisini ölçə bilsəydi, nəzəri cəhətdən deterministik idi →
  növbəti səviyyə: kvant.

### 6. Kvant təsadüfülük
**Nədir:** Foton polarizasiyası ölçməsi — superpozisiya: ölçmədən əvvəl nəticə
prinsipcə bilinə bilməz (gizli switch yoxdur; "hər iki halda, amma fərqli
reallıqlarda").

**Nəticə:** İstənilən qədər irəliləmiş "alien texnologiyalı" Eve belə kvant
nəticəsini əvvəlcədən bilməyəcək.

**ANU qrand (qrand.NewReader(apiKey)):**
```go
q := qrand.NewReader(apiKey)
buf := make([]byte, 32)
_, err = q.Read(buf) // 32 baytlıq kvant açarı
```

**Amma diqqət:** ANU serverindən HTTP ilə gələn təsadüfülük **optik tapa bilər**
(Eve kaboldan oxuya bilər / data mərkəzində yaradılarkən görə bilər). Yerli
entropy pool az-az təsadüfü olsa belə, **daha təhlükəsizdir** — intercept
imkanları azdır.

## Əsas terminlər

- PRNG (pseudo-random number generator / psevdo-təsadüfi generator)
- Seed (toxum / başlanğıc dəyər)
- Periodicity (dövrilik)
- Environmental Noise (mühit şüaması)
- TRNG (true random number generator / həqiqi təsadüfülük generatoru)
- Thermal Noise (istilik şüaması)
- Superposition (kvant superpozisiyası)
- Wrench Attack (açarlı hücum)

## Praktik nəticə

- Oyun/simulyasiya üçün PRNG kifayətdir; kriptoqrafik açar üçün YOX
- `math/rand` vaxt-seedlidir, periodicity + statistik tələlərə meyllidir
- Kriptoqrafik açarlar: `crypto/rand` (yerli entropy pool) — uzaqdan gələn
  "daha təsadüfü" məlumatdan təhlükəsizdir
- Təhlükəsizlik mütləq deyil — hücumçunun səy səviyyəsinə görə qiymətləndirilir
- Ən zəif halqa adətən alqoritm deyil — insan/infrastruktur tərəfdir

## Mənbə

Pages: 213-225 (Chapter 13, Explore Go: Cryptography)
