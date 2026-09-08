# Chapter 12 — Entropy (səh. 199-212)

## Bu chapter nədən bəhs edir?

Entropiya — məlumatın "həqiqi" təsadüfülük ölçüsü: bit-string-lərin informasiya
məzmunu, sıxışdırma ilə əlaqəsi, Kolmogorov mürəkkəbliyi (K) və kriptoqrafik
açar üçün nə üçün yüksək entropiyanın həyatı-mət olduğunu izah edir. Cloudflare-in
lava lamp "Wall of Entropy"-si real dünya nümunəsidir.

## Əsas fikirlər

### 1. Informasiya və entropiya
**Nədir:** N mümkün mesaj varsa, hər bir seçimin entropiyası = seçimi unikal
ifadə etmək üçün lazım olan **bit sayı** = eyni zamanda "neçə qəpik atmaşı
(coin toss)" lazımdır.

**Vacib:** Mesajın öz bit sayı entropiya deyil — bəzi bitlər **redundant** ola
bilər (ünikal informasiya əlavə etmirlər).

**Nümunə (ingilis mətni):**
```
F U CN RD THIS, U CN BCM A SEC & GT A GD JB W HI PA
```
Vokallar atılıb, amma oxunur — deməli hərflərin bir hissəsi redundantdır.
Müntəzəm ingilis mətni real entropiyası ~1.5 bit/hərf, tam 8 bit yox.

### 2. Şifrələnmiş mətn sıxılmır
**Nədir:** Yaxşı şifrlənmiş mətn **sıxılmaz** (incompressible). Hər bit
plaintext-i bərpa etmək üçün lazımdır — redundansiya yoxdursa, sıxılma alınmır.

**Nəyə lazımdır:** Eve üçün potensial plaintext sayı = onun bilmədiyi şeyin
entropiyası. Kamera freymi 100MiB olsa da, ehtiyat məlumat (image redundancy)
sıxışdırılır → real entropiya daha az → Eve-in axtarış fazası kiçik olur.

**Lossy compression (itkil sıxılma):** domain biliyi tələb edir — insan gözü
(piksel fərqlərini görür), qulağın tezlik cavabı (audio) → "həyatın hansı
hisəsəsini itirmək olar" biliklən məlumat atılır.

### 3. Kolmogorov mürəkkəbliyi (K)
**Nədir:** Sistemin mürəkkəbliyi = **onları yaradan ən qısa proqram**.

**Nümunələr:**
- "1 milyondan kiçik bütün cüt müsbət ədədlər" → 123 baytlıq Go proqramı
  (984 bit) — milyonlarla bayt siyahısından qat-qat qısa
- π-nin rəqəmləri təsadüfü görünür, amma `π/4 = 1 - 1/3 + 1/5 - 1/7 + ...`
  düsturu hər rəqəmi yaradır — "pi" sözü kifayətdir → aşağı K
- Mandelbrot set (mürəkkəb görünür): `z → z² + c` → aşağı K

**Kriptoqrafik dərs:** Uzun və mürəkkəb **görünən**, amma qısa qayda ilə
təsvir olunan açar = **alçaq K = qeyri-təhlükəsiz**. Eve belə açarı asan tapır.

### 4. Randomness = yüksək entropiya
**Nədir:** Həqiqi təsadüfi ardıcıllıq: onu yaradan ən qısa proqram, ardıcıllığın
özündən qısa ola bilməz.

**Nümunələr:** Dünyadakı hər çimərlikdəki hər qum dənəsinin diametri; qalaktikadakı
ulduzların kütlələri — bunları yaratmaq üçün proqram yoxdur, yalnız özünü
göstərmək mümkündür → yüksək K.

**Nə üçün təsadüfi açar yaxşıdır:** 2^(256) mümkün açardan təsadüfi seçim,
statistik olaraq asanlıqla təsvir olunan (alçaq K) açara düşmə ehtimalı
sıfıra yaxındır → çoxlu entropiya → brute-force qeyri-mümkün (Chapter 11
enerji limiti).

### 5. Praktik qeyri-determinizm
**Nədir:** Proqram tərəfindən təsadüf seçmək — kompüter qəpik ata bilmir, amma...

**Deterministik görünən, amma praktik olmayan:** Qəpik atışının nəticəsi əgər
ilkin şərti tam bilsən deterministikdır → "chaotic" (xaotik) amma deterministik
fiziki proseslər praktik olaraq qeyri-prediktabləşdirilə bilər.

**Cloudflare Wall of Entropy:** ~100 lava lampın video görüntüsündən piksel
dataları çıxarılır → kriptoqrafik açar üçün giriş entropiyası. Gözəl görünür,
funksionaldır.

**İnsan beyni yoxsul entropiya mənbəyidir:** Şahmat/jan-jürsiz oyunlarda
"random" olmağa çalışsan, özünü əleyhinə pattern-lərə düşürsən. İnsan tərəfindən
uydurulmuş "təsadüfi" sıralar pattern analizi ilə sınanır (LECTER sitatı:
"random scattering of sites seem desperately random, like the elaborations of a
bad liar").

## Əsas terminlər

- Entropy (entropiya / məlumat miqdarı)
- Redundancy (ehtiyatlıq / lazımsız təkrar)
- Lossless/Lossy Compression (itkisiz/itkili sıxışdırma)
- Kolmogorov Complexity (Kolmoqorov mürəkkəbliyi)
- Run-length Encoding (uzunluq-kodlaşdırma)
- Wall of Entropy (Cloudflare lava lamp divarı)
- Randomness (təsadüfülük)

## Praktik nəticə

- Açar uzunluğu tək başına kifayət deyil — entropiya əsasdır
- Şifrələnmiş faylı sıxmağa çalışmaq (sıxılmır — bunun özü yaxşı şifrə göstəricisidir)
- Alçaq K açarlar (uzun amma təsvir oluna bilən) qeyri-təhlükəsizdir
- Yüksək entropiyalı açar = həqiqi təsadüfi mənbədən (lava lamp, hardware RNG)
- İnsan tərəfindən seçilən açarlar prediktabləşdirilə bilər — `crypto/rand` işlət

## Mənbə

Pages: 199-212 (Chapter 12, Explore Go: Cryptography)
