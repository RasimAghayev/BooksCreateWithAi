# Chapter 18 — Security (Təhlükəsizlik)

## Bu fəsil nədən bəhs edir?

SQL injection (concat vs Prepare placeholder), command injection (exec.Command
tələsi), XSS (text/template vs html/template), hashing (MD5/SHA/BLAKE), simmetrik
şifrələmə (AES-GCM Seal/Open + nonce), asimmetrik (RSA EncryptOAEP/DecryptOAEP),
random generatorlar (math vs crypto/rand), TLS (self-signed x509 sertifikat
generasiyası, HTTPS server/client handshake) və bcrypt parol idarəsi.

## Əsas fikirlər

### 1. Təhlükəsizlik Fəlsəfəsi
- Security SONRADAN DEYİL, dizayn hissəsidir — gündəlik praktika (code kata)
- Vulnerability-lərin kökü: hücum vektorlarını BİLMƏMƏK + deploy-dan əvvəl review
  OLMAMASI
- Tanınmış vektorlardan (SQL injection və s.) qorunsan — hücumların BÖYÜK hissəsi
  dəf olunur
- OWASP Top 10 — ən common zəiflik reyestri

### 2. SQL Injection — Concat Tələsi
**Təhlükəli kod:**
```go
query := `SELECT CARD_NUMBER FROM USER_DETAILS WHERE USER_ID = ` + userID
// Input: "" OR '1' == '1'
// → SELECT ... WHERE USER_ID = "" OR '1' == '1'  — HAMISININ kartı sızır!
```
**Həll — Prepare + placeholder:**
```go
func GetCardNumberSecure(db *sql.DB, userID string) (resp string, err error) {
    stmt, err := db.Prepare(
        `SELECT CARD_NUMBER FROM USER_DETAILS WHERE USER_ID = ?`)
    if err != nil { return resp, err }
    defer stmt.Close()
    row := stmt.QueryRow(userID)         // input DƏYƏR kimi — sorğu DƏYİŞMİR
    ...
}
```
**Qayda:** istifadəçi input-u HƏMİŞƏ placeholder; concat heç vaxt.

### 3. Command Injection — OS Əmrləri
**Təhlükəli kod:**
```go
func listFiles(path string) (string, error) {
    cmd := exec.Command("bash", "-c", "ls"+path)    // input ƏMR kimi gİRİR!
    ...
}
// Input: " .; cat /etc/hosts"  → ls SADECE DEYİL, hosts faylı DA OXUNUR!
```
**Səbəblər:** sanitizasiya YOX; istənilən string QƏBUL; əlavə komandalar icra
olunur. **Həll:** istifadəçi inputunu OS komandasına DİREKT VERMƏ; lazımdırsa
ağ siyahı/validasiya.

### 4. XSS — Script İnjection
**Təhlükə:** `<script>alert("Hello")</script>` comment kimi göndərilir →
qurbanın brauzerində İCRA OLUNUR (data oğurlama, izləmə).

**Kitabdan həll — template PAKETİ dəyiş:**
```go
// TƏHLÜKƏLİ:
import "text/template"      // escape ETMİR — script KOD kimi işləyir

// TƏHLÜKƏSİZ:
import "html/template"      // AVTOMATİK ESCAPE — script MƏTN kimi görünür
```
Yalnız bir söz fərqi — html/template bütün HTML entity-ləri escape edir.

### 5. Hashing — Bir Yönlü Çevirmə
**Nədir:** plaintext → unikal sabit uzunluq; geri QAYTARILA BİLMİR (one-way);
collision ehtimalı minimal.

**Alqoritm seçimi:**
| Alqoritm | Status |
|---|---|
| MD5, SHA-1 | TƏHLÜKƏSİZ DEYİL (brute-force asan) — yalnız checksum |
| SHA-256/512 | Standart |
| SHA3, BLAKE2s/2b | Müasir (x/crypto paketindən) |

**Kitabdan kod nümunəsi:**
```go
func getHash(input string, hashType string) string {
    switch hashType {
    case "MD5":
        return fmt.Sprintf("%x", md5.Sum([]byte(input)))
    case "SHA256":
        return fmt.Sprintf("%x", sha256.Sum256([]byte(input)))
    case "SHA512":
        return fmt.Sprintf("%x", sha512.Sum512([]byte(input)))
    case "SHA3_512":
        return fmt.Sprintf("%x", sha3.Sum512([]byte(input)))   // x/crypto
    case "BLAKE2s_256":
        return fmt.Sprintf("%x", blake2s.Sum256([]byte(input)))
    default:
        return fmt.Sprintf("%x", sha256.Sum256([]byte(input)))
    }
}
```
golang.org/x/* — Go layihəsi, amma standart quraşdırmadan XARİC → go get lazım.

### 6. Symmetric Encryption — AES-GCM
**Nədir:** EYNİ açar şifrələyir və açır (parol əsaslı).

**Kitabdan kod nümunəsi:**
```go
func encrypt(data []byte, key string) (resp []byte, err error) {
    block, err := aes.NewCipher([]byte(key))     // açardan cipher bloku
    gcm, err := cipher.NewGCM(block)             // GCM wrapper (auth-lı şifrələmə)
    nonce := make([]byte, gcm.NonceSize())        // tək-istifadə random
    rand.Read(nonce)                               // crypto/rand!
    return gcm.Seal(nonce, nonce, data, []byte("test")), nil
    //            ^^^ dst = nonce → nonce CIPHER-İN BAŞINDA saxlanılır!
}

func decrypt(data []byte, key string) (resp []byte, err error) {
    block, _ := aes.NewCipher([]byte(key))
    gcm, _ := cipher.NewGCM(block)
    ciphertext := data[gcm.NonceSize():]        // başdan nonce
    nonce := data[:gcm.NonceSize()]              // sonrası cipher
    resp, err = gcm.Open(nil, nonce, ciphertext, []byte("test"))
    return resp, err
}
```
**Nonce pattern:** Seal-in dst parametrinə nonce ver → nəticə [nonce|ciphertext]
bir byte massivi; decrypt başdan ayırır. additionalData decryptdə EYNİ olmalıdır.

### 7. Asymmetric Encryption — RSA
**Nədir:** açar CÜTÜ — public (hamıya paylaşılan, ŞİFRƏLƏMƏ) + private (səndə,
AÇMA). EncryptOAEP/DecryptOAEP.

**Kitabdan kod nümunəsi:**
```go
privateKey, err := rsa.GenerateKey(rand.Reader, 1024)   // açar cütü
publicKey := privateKey.PublicKey                        // public çıxarılır

// PUBLIC ilə şifrələ:
ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader,
    &publicKey, text, nil)

// PRIVATE ilə aç:
decrypted, err := rsa.DecryptOAEP(sha256.New(), rand.Reader,
    privateKey, ciphertext, nil)
```

### 8. Random Generatorlar — math vs crypto
| math/rand | crypto/rand |
|---|---|
| Sürətli | Yavaş (OS entropiyası — /dev/urandom) |
| Təhlükəsizlik üçün YOX | Security-sensitive üçün MÜTLƏQ |
| Simulyasiya/oyun | Session ID, açar, nonce |

```go
// crypto/rand — big.Int qaytarır:
data, _ := rand.Int(rand.Reader, big.NewInt(1000))

// math/rand — alias ilə ayır:
import math "math/rand"
math.Intn(1000)
```
**Niyə vacib:** session ID pattern-li olsa → attacker NÖVBƏTİNİ TAPAR →
sessiya oğurluğu.

### 9. TLS — Nəqliyyat Təhlükəsizliyi
**TLS 4 təminatı:** Identity (sertifikat), Integrity (message digest — dəyişmə
yoxdur), Authentication (public-key), Confidentiality (şifrə).

**Self-signed sertifikat generasiyası (kitabdan):**
```go
func generate() (cert []byte, privateKey []byte, err error) {
    serialNumber, _ := rand.Int(rand.Reader, big.NewInt(27))
    ca := &x509.Certificate{
        SerialNumber: serialNumber,
        Subject:      pkix.Name{Organization: []string{"example.com"}},
        NotBefore:    time.Now(),
        NotAfter:     time.Now().Add(365 * 24 * time.Hour),  // 1 il
        // ...
    }
    rsaKey, _ := rsa.GenerateKey(rand.Reader, 2048)
    DER, _ := x509.CreateCertificate(rand.Reader, ca, ca,
        &rsaKey.PublicKey, rsaKey)                    // self-signed: ca özünü imzalayır

    // DER (binary) → PEM (ASCII):
    cert = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: DER})
    privateKey = pem.EncodeToMemory(&pem.Block{
        Type:  "RSA PRIVATE KEY",
        Bytes: x509.MarshalPKCS1PrivateKey(rsaKey)})
    return cert, privateKey, nil
}
```

**HTTPS server + client handshake (kitabdan):**
```go
// Server — TLS tələb edir:
func runServer(certFile, key string, clientCert []byte) error {
    http.HandleFunc("/", hello)
    server := &http.Server{Addr: ":443"}
    cert, err := tls.LoadX509KeyPair(certFile, key)
    // tls.Config{Certificates: ..., ClientAuth: tls.RequireAndVerifyClientCert}
    // ClientCAs: client pool — CLIENT də doğrulanır!
    server.ListenAndServeTLS(certFile, key)
}

// Client — CA + öz sertifikatı:
func client(caCert []byte, ClientCert tls.Certificate) error {
    certPool := x509.NewCertPool()
    certPool.AppendCertsFromPEM(caCert)             // server-i doğrulamaq üçün
    client := &http.Client{
        Transport: &http.Transport{
            TLSClientConfig: &tls.Config{
                RootCAs:      certPool,               // CA etibarı
                Certificates: []tls.Certificate{ClientCert},  // öz identitası
            },
        },
    }
    resp, err := client.Get("https://127.0.0.1:443")
}
```
**Mutual TLS:** server client-i, client server-i doğrulayır — bank səviyyəli
əlaqə. CA = Certificate Authority (sertifikat imzalayıcı).

### 10. Password Management — bcrypt
**Qızıl qayda:** parolu plaintext HEÇ YERDƏ saxlama (yaddaş və ya DB).
**Həll:** one-way hash + müqayisə.

**Kitabdan kod nümunəsi:**
```go
import "golang.org/x/crypto/bcrypt"

password := "mysecretpassword"
encrypted, _ := bcrypt.GenerateFromPassword([]byte(password), 10)  // cost 10
// DB-yə SADECE encrypted yazılır

err := bcrypt.CompareHashAndPassword([]byte(encrypted), []byte(password))
if err == nil {
    fmt.Println("Password matched")    // parolun ÖZÜ heç vaxt saxlanmır
}
```
bcrypt niyə: yavaş (brute-force bahalı), avtomatik salt, parol sahəsi üçün standart.
(Sadelə SHA hash parollar üçün KƏSİNLİKLə azdır — rainbow table-larla saniyələrdə
qırılır; bcrypt bu qoruma daxil edir.)

## Əsas terminlələr
- Attack Vector — hücum yolu (SQL/command injection, XSS)
- SQL Injection — concat input sorğunu DƏYİŞİR
- Prepared Statement/Placeholder — `?` ilə təhlükəsiz parametr
- Command Injection — OS əmrinə input sızması
- XSS — `<script>` web səhifəsinə inyeksiya; qurban brauzerində icra
- text/template vs html/template — escape YOX / AVTOMATİK escape
- OWASP Top 10 — ümumi zəiflik reyestri
- Hash — one-way; MD5/SHA1 qırılmış, SHA256+/BLAKE/SHA3 etibarlı
- Symmetric Encryption — eyni açar (AES-GCM)
- Nonce — tək-istifadə random; Seal-in dst-də gizlədilir
- Asymmetric/Public-Key — public şifrələr + private açır (RSA OAEP)
- math/rand vs crypto/rand — oyun vs security (OS entropiyası)
- TLS — Identity+Integrity+Authentication+Confidentiality
- x509 Sertifikat — digital identity; PEM ASCII format
- Self-Signed — sertifikat özünü imzalayır (dev/test)
- CA / Certificate Authority — sertifikat imzalayıcı
- Mutual TLS — hər iki tərəf sertifikat doğrulayır
- bcrypt — parol hash; GenerateFromPassword/CompareHashAndPassword

## Praktik nətidə

(1) İstifadəçi input-u: SQL-də placeholder, OS əmrində HEÇ VAXT; HTML-də
html/template (avtomatik escape). (2) SQL-də bir concat belə riskdir — "OR '1'='1"
klassikası bütün datanı sızdırır. (3) exec.Command-a input qatmaq = serverdə
İSTƏDİYİN əmri icra etmək hüququ vermək. (4) text/template YALNIZ mətn üçün — web
səhifəsində html/template MÜTLƏQ. (5) Hash seçimi: MD5/SHA1 yalnız checksum;
sənədləşdirmədə SHA256+. (6) AES-GCM: nonce Seal-in dst-ndə gizlət — decrypt
başdan ayırır; additionalData hər iki tərəfdə eyni. (7) RSA: public hamıya, private
SİNDƏ — EncryptOAEP(public)/DecryptOAEP(private). (8) Security-sensitive random =
crypto/rand — session ID pattern-li olsa oğurluq. (9) TLS 4 təminat: kimlik,
bütövlük, autentikasiya, konfidansiallıq — prod üçün mütələq. (10) Parollar:
plaintext YOX, bcrypt YOXSA plain SHA — rainbow table qoruyucusu bcrypt-də VAR.
(11) Öz şifrələmə alqoritm ixtira ETMƏ — güc riyaziyyatda, gizlilikdə DEYİL.

## Mənbə
Pages: 625-659 (PDF 658-689)
