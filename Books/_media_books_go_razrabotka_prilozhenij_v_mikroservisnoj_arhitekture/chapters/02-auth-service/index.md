# Глава 2 — Avtorizasiya və autentifikasiya mikroservisi (Auth) (səh. 87-110)

## Bu fəsil nədən bəhs edir?

Identifikasiya / autentifikasiya / avtorizasiya terminlərinin dəqiq ayrımı;
5 autentifikasiya üsulu (parol, sertifikat, OTP/2FA, access key, token —
JWT/SWT/SAML), HTTP autentifikasiya sxemləri (Basic/Digest/NTLM), data
ötürmə yerlərinin təhlükəsizlik müqayisəsi, ümumi zəifliklər və müdafiə
(buffer overflow, race condition, input validation, XSS/CSRF/clickjacking);
Auth mikroservisinin implementasiyası: kontrakt (Register/Login/Refresh/
Validate/Logout → TokenPair), JWT konfiqurasiyası (secret, TTL-lər), users +
refresh_tokens cədvəlləri (FK ON DELETE CASCADE, revoked_at), bcrypt
(DefaultCost=10, salt), access JWT (HS256, sub/exp) + refresh SHA-256 token,
sessiya idarəetməsi.

## Əsas fikirlər

### 1. Üç Termin — Dəqiq Ayrım
| Termin | Nədir | Nümunə |
|---|---|---|
| Идентификация (identifikasiya) | "Mən KİMƏM" — özünü tanıtma | login formunu doldurmaq (save basmadan belə) |
| Аутентификация (autentifikasiya) | "Sən həqiqətən o-san" — sübut | parol yoxlaması, sənəd+şəkil müqayisəsi |
| Авторизация (avtorizasiya) | "Sənə icazə VAR" — resurs çıxışı | moderator rolu → şərh silmək |

### 2. HTTP Parol Autentifikasiyası
- **Axın:** 401 + `WWW-Authenticate` başlığı → login səhifəsi → hər sorğuda
  `Authorization` başlığı → server yoxlayır + avtorizasiya
- **Sxemlər:**
  - **Basic** — login:parol base64 (ŞİFRƏLƏNMƏYİB!; HTTPS-də qəbul oluna
    bilər)
  - **Digest** — challenge/response: server nonce (vaxt damğası) göndərir,
    client MD5(parol, nonce) qaytarır
  - **NTLM/Negotiate** — Windows domeni üçün, az yayılıb
- **Parol ötürmə yeri:** URL query — TƏHLÜKƏLİ (brauzer/web-server
  yaddaşında qalır, intercept); Request body — POST/PUT/PATCH üçün OK;
  **HTTP header — OPTİMAL** (Authorization və ya custom)
- **Çatışmamazlıq:** HTTP auth-da standart "logout" yoxdur (brauzer
  pəncərələrini bağlamaqdan başqa)

### 3. Sertifikat Autentifikasiyası
- **CA (Certificate Authority):** pasport verən orqan kimi — sertifikatın
  doğruluğuna zəmanət; X.509 formatı, SSL/TLS daxilində
- **Yoxlama qaydaları:** (1) etibarlı CA tərəfindən imzalı; (2) müddət
  etibarlıdır; (3) ləğv edilməyib (revocation list)
- **Gücü:** rəqəmsal imza → inkar edilməzlik (non-repudiation); zəif cəhəti:
  yayılma/idarəetmə çətinliyi

### 4. OTP / İki Faktorlu (2FA)
- **Nədir:** bildiklərin (parol) + sahiblik (token) kombinasiyası
- **Mənbələr:** hardware token (USB dongle, smart kart, biometrik),
  SMS/kod (SIM faktoru), TOTP (Google Authenticator — <1 dəqiqəlik kod),
  scratch card (nömrələnmiş çap siyahısı)
- **İstifadə:** maliyyə əməliyyatları, VPN, kritik hesab dəyişiklikləri

### 5. Access Key (API Key)
- **Nədir:** uzun unikal string — login+parolu ƏVƏZ edir; cihaz/servis
  autentifikasiyası üçün
- **Müddət + səviyyə limiti:** yaradılanda有效期 və access scope qoyula bilər
- **Açıq şəbəkə üçün imza sxemi:** açar = public (identifikasiya) + secret
  (imza); server nonce/timestamp göndərir → client HMAC/Hash(nonce, secret)
  qaytarır → tam açar ötürülmür + **replay attack** qarşısı

### 6. Token Autentifikasiyası — JWT
- **Model:** IP (Identity Provider — kimlik təchizatçısı) + SP (Service
  Provider): client IP-də token alır, SP-yə təqdim edir
- **Aktiv client** (mobil tətbiq — özü addımları icra edir) vs **passiv
  client** (brauzer — IP/SP arasında yönləndirilir)
- **SP yoxlamaları:** (1) etibarlı IP verib; (2) bu SP üçündür; (3) müddəti
  keçməyib; (4) dəyişdirilməyib (imza)
- **Formatlar:**
  - **SWT** — sadə ad/dəyər; HMACSHA256 simmetrik açar (IP+SP hər ikisində)
  - **JWT** — 3 blok: header.payload.signature (ilk ikisi JSON+base64);
    claims (iss, aud, exp...); simmetrik VƏ YA asimmetrik imza; OAuth və
    OpenID Connect-in də əsası
  - **SAML** — XML; asimmetrik imza; ownership sübutu → MitM müdafiəsi
- **Keycloak:** hazır OpenID Connect həlli — kiçik layihələr üçün güclü;
  böyük/kommersiya layihələrdə limitlər, yüksək sessiya sayında
  qeyri-stabilik xəbər verilir; çox sadə layihə üçün isə overkill

### 7. Zəifliklər və Müdafiələr
| Hücum / zəiflik | Təbiəti | Müdafiə |
|---|---|---|
| Buffer overflow | 200 simvol üçün ayrılmış yaddaşa 10 000 daxil | GİRİŞ VALIDASİYASI — sorğu biznes məntiqinə ÇATMAMAZDAN yoxla (11 rəqəmli telefon = dəqiq 11) |
| Race condition | paralel yazılar eyni data-ya; icra sırası → nondeterminizm | düzgün sinxronizasiya |
| Input validation attacks | xüsusi simvollar/ekranlaşdırma dəyərləri ilə daxili yaddaşa baxış, icazəsiz yazı | sahə-uyğun format yoxlaması |
| Autentifikasiya hücumları | bruteforce: qwerty01 → ANINDA; S7tvZsVm1q (10 simvol, qarışıq registr) → 20+ İL | mürəkkəb parol tələbi; autentifikasiya YALNIZ SERVER tərəfində — client "hər şey qaydasındadır" mesajı SAXTALANA bilər |
| Avtorizasiya hücumları | icazəsiz resurs çıxışı | avtorizasiya serverdə; minimum imtiyaz prinsipi; HƏR privileged əməliyyətdə yoxlama |
| XSS | səhifəyə gizlədilmiş skript — baxışda icra olunur | brauzer aktual saxlamaq, NoScript kimi alətlər |
| CSRF | gizli link — avtorizasiya olmuş səhifədə avtomatik əməliyyat (səbətə əlavə, pul köçürmə) | brauzer müdafiələri |
| Clickjacking | "Ətraflı" düymə altında gizli "İndi al" düyməsi | brauzer müdafiələri |

### 8. Auth Servisinin Yerləşmə Qərarı
- **Variant A (Auth User-ə sorğu göndərir):** hər yoxlamada ekstra şəbəkə
  sorğusu → performans xərçi
- **Variant B (seçilmiş — data ayrılığı):** Account = biznes data; Auth =
  parol/hash/token — öz DB, ekstra sorğu YOX
- **Token ömrü trade-off:** hər sorğuda DB-dən status (bloklanıb?) yoxlama
  (təhlükəsizlik, amma yük) VS qısa TTL 1-2 saat (seçilən yanaşma)

### 9. Auth Kontraktı (gRPC)
```protobuf
service Auth {
  rpc Register(RegisterRequest) returns (google.protobuf.Empty);
  rpc Login(LoginRequest) returns (TokenPair);          // login VEYA email + parol
  rpc Refresh(RefreshRequest) returns (TokenPair);      // refresh_token → yeni cüt
  rpc Validate(ValidateRequest) returns (ValidateResponse;  // access_token → user_id
  rpc Logout(RefreshRequest) returns (google.protobuf.Empty); // sessiyanı bitir
}
message TokenPair { string access_token = 1; string refresh_token = 2; }
```
- Makefile-a `auth` kataloqu əlavə olunur (gen/clean siyahılarına); yeni
  versiya: `git tag v1.1.0` (semver — funksiya əlavəsi, breaking YOX)

### 10. JWT Konfiqurasiyası və DB Strukturu
```go
JwtSecret            string `env:"JWT_SECRET" required:"true"`
AccessTokenTTLMinutes int  `env:"ACCESS_TOKEN_TTL_MINUTES" default:"60"`
RefreshTokenTTLDays   int  `env:"REFRESH_TOKEN_TTL_DAYS" default:"30"`
```
```sql
CREATE TABLE users (id BIGSERIAL PRIMARY KEY, login TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE, password_hash TEXT NOT NULL, ...);
CREATE TABLE refresh_tokens (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token TEXT NOT NULL UNIQUE, expires_at TIMESTAMP NOT NULL,
    revoked_at TIMESTAMP NULL, created_at ...);
```
- **İKİ cədvəl bir tranzaksiyada** — uğursuzluqda hər ikisi rollback
- **ON DELETE CASCADE:** istifadəçi silinəndə tokenlər də silinir
- **revoked_at:** NULL = aktiv; vaxt damğası = ləğv edilmiş

### 11. Service Məntiqi — bcrypt + JWT + SHA-256
```go
// Register: parol → hash
hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

// Login: yoxla + token cürü burax
user, err := s.repo.GetUserByLoginOrEmail(ctx, req.LoginOrEmail)  // login OR email
bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
return s.issueTokens(ctx, user.ID)

// Refresh: ləğv/müddət yoxlaması
rt, _ := s.repo.GetRefreshToken(ctx, refreshToken)
if rt.RevokedAt != nil || time.Now().After(rt.ExpiresAt) { ... }

// Validate: JWT imza + sub çıxarma
jwt.ParseWithClaims(accessToken, claims, func(t *jwt.Token) (interface{}, error) {
    return []byte(s.cfg.JwtSecret), nil
})
uidFloat, ok := claims["sub"].(float64)  // JSON rəqəmlər = float64!
```
- **Access token:** JWT HS256; claims: `sub` (user ID), `exp` (Unix vaxtı)
- **Refresh token:** SHA-256(userID:unixnano:secret) → hex; DB-də saxlanılır
- **Logout:** `Update("revoked_at", gorm.Expr("NOW()"))`
- **bcrypt əsasları:**
  - Açıq parol DB-də saxlanmaz — leak halında bütün istifadəçilər təhlükədə
  - **Salt** avtomatik qarışdırılır → eyni parollar = FƏRLİ hash-lər
  - **DefaultCost = 10** — hesablama bilərəkdən YAVAŞ → bruteforce praktiki
    olaraq mümkünsüz, proqram üçün isə hiss olunmur
- **GORM qeydi:** Create-də Model göstərilmir (model məlumatla birlikdə
  ötürülür → INSERT avtomatik qurulur); First(&user) receiver-dən model
  çıxarır; axtarışda Model() QİYMƏTLİNDİRİLMƏLİ (aydınlıq üçün)

## Əsas terminlər
- Identifikasiya / Autentifikasiya / Avtorizasiya — kim / sübut / icazə
- Nonce — təkrar istifadəni bloklayan unikal dəyər (vaxt damğası)
- Replay attack (təkrar oxutma hücumu) — tutulmuş sorğunun yenidən göndərilməsi
- 2FA (iki faktorlu autentifikasiya) — biliş + sahiblik
- Identity Provider (IP) / Service Provider (SP) — token verən / qəbul edən
- Claims — JWT payload sahələri (iss, aud, exp, sub)
- Bearer token — `Authorization: Bearer [token]` sxemi
- Salt (duz) — hash-ə qarışdırılan təsadüfi baytlar
- Cost factor — bcrypt-in bilərəkdən yavaşlığı (DefaultCost=10)
- Revocation (ləğv) — refresh token-in revoked_at ilə söndürülməsi

## Praktik nəticə

1. **Data ayrılığı:** Auth = parol/hash/token ÖZ bazasında — hər dəfə User
   servisinə sorğu göndərmə (performans) VƏ YA qısa access TTL (1-2 saat).
2. **Parol heç vaxt açıq saxlanmır:** bcrypt + DefaultCost; hash
   müqayisəsi yalnız bcrypt.CompareHashAndPassword ilə.
3. **Token cütü:** access = qısaömürlü JWT (stateless, DB-siz yoxlama);
   refresh = uzunömürlü DB-da (revocation imkanı) — Logout revoked_at.
4. **Validasiya qapıda:** sorğu biznes məntiqinə çatmazdan ƏVVƏL — buffer
   overflow və input hücumlarının kök həlli.
5. **Autentifikasiya/avtorizasiya YALNIZ server tərəfdə:** client tərəfli
   yoxlama = saxta "OK" mesajı ilə keçilə bilər.
6. **claims["sub"] float64-dir** — JSON rəqəm semantic-ləri; uint64-ə
   çevirməyi unutma.

## Mənbə
Pages: 87-110 (PDF 88-111)
