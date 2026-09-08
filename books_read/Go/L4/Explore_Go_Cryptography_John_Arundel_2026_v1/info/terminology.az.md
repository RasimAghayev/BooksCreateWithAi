# Explore Go: Cryptography — Terminologiya (AZ)

| Termin | Azərbaycanca qarşılıq | İzah |
|---|---|---|
| Cipher | şifrə | plaintext-i ciphertext-ə çevirən alqoritm |
| Plaintext / Ciphertext | açıq mətn / şifrələnmiş mətn | input / output |
| Encipher / Decipher | şifrələmək / dəşifrə etmək | yönlü transformasiya |
| Crack | sındırmaq | açarsız plaintext bərpası |
| Crib | tanış parça | bilinən plaintext fraqmenti (sındırma üçün) |
| Keyspace | açarlar fəzası | bütün mümkün açarlar çoxluğu |
| Brute force | qüvvə hücumu | bütün açarları ardıcıl sınamaq |
| Block cipher | blok şifrəsi | fiks ölçülü bloklarla işləyir |
| Stream cipher | axın şifrəsi | bit-bit, hardware dünyası |
| BlockSize | blok ölçüsü | cipher.Block metodlarından biri |
| ECB | elektron kod kitabı | eyni input → eyni output; QADAĞAN |
| CBC | zəncirli blok şifrələmə | hər blok əvvəlkidən asılı |
| IV | başlanğıc vektoru | CBC-nin ilk "əvvəlki blok" random dəyəri |
| Nonce | tək-istifadə dəyəri | bir dəfəlik rəqəm (replay qarşısı) |
| Padding | doldurma | bloğa tam almaq üçün artıq baytlar (PKCS#7) |
| Entropy | entropiya | informasiyanın həqiqi bit miqdarı |
| Kolmogorov complexity | Kolmogorov kompleksliyi | ən qısa təsvirin ölçüsü |
| PRNG / TRNG | psevdo / həqiqi təsadüfi generator | deterministik / fiziki |
| Digest | həzm | hash funksiyasının çıxışı |
| Hashing | hashing | sabit uzunluqlu barmaq izi |
| Preimage attack | ilkin-obraz hücumu | hədəf hash-ə uyğun mesaj tapmaq |
| Collision | toqquşma | eyni hash-li iki fərqli mesaj |
| Avalanche effect | sel effekti | kiçik input dəyişikliyi → böyük hash fərqi |
| Salt | duz | parol hash-inə random əlavə (rainbow table qarşısı) |
| MAC | mesaj autentifikasiya kodu | hash(key + msg) — bütövlük |
| HMAC | nested MAC standartı | H(key + H(key + msg)) |
| Replay attack | təkrar hücumu | saxta blokun/mesajın yenidən göndərilməsi |
| Key exchange | açar mübadiləsi | təhlükəsiz kanalsız ortaq açar |
| Diffie-Hellman | Diffi-Hellman | g^ab mod p söhbəti |
| Discrete logarithm | diskret logaritm | g^x = y → x tapmaq (çətin istiqamət) |
| Public / private key | açıq / xüsusi açar | RSA cütü |
| Certificate Authority | sertifikat orqanı | açarların etimad zənciri |
| S-box | substitusiya cədvəli | AES-in bayt əvəzetmə cədvəli |
| Confusion / diffusion | qarışdırma / yayılma | AES-in 2 prinsipi |
| Round key | raund açarı | açar genişləndirməsinin addımları |
| AEAD | autentifikasiyalı şifrələmə | encryption + integrity birlikdə |
| AES-GCM | tövsiyə olunan AES modu | Seal/Open, padding daxildir |
| Proof of work | iş sübutu | hash < target əməliyyatı (mining) |
| Double-spend | ikiqat xərcləmə | eyni coin-i iki dəfə xərcləmək |
| Eventually consistent | sonunda ardıcıl | distributed ledger xüsusiyyəti |
| Quantum superposition | kvant superpozisiyası | qubit 0 və 1 "eyni anda" |
| Post-quantum crypto | post-kvant kriptoqrafiya | kvanta davamlı alqoritmlər |
