# Chapter 7 — Sessions, Authentication, and Authorization (Sessiyalar, Autentifikasiya və Avtorizasiya)

## Bu chapter nədən bəhs edir?
Sessiya anlayışına, sessiya data saxlama variantlarına (client/server/Redis), təhlükəsiz cookie idarəsinə (HTTPOnly/Secure/SameSite/encrypt), istifadəçi autentifikasiyasına (bcrypt qeydiyyat/login/logout), parol qorunmasına (hash/salt/pepper/iterations) və OAuth inteqrasiyasına.

## Əsas fikirlər

### 1. Sessiyalar nədir və niyə vacibdir
**Nədir:** İstifadəçi ilə app araslı müvəqqəti, interaktiv informasiya mübadiləsi — səhifələrarası kontinuitet.

**Vazifələri:** səbətin saxlanması, axtarış tarixçəsi, personalizasiya (tövsiyələr), autentifikasiya vəziyyətinin saxlanması (hər səhifədə login yenidən YOX).

**Scale problemi:** Minlərlə eynizamanı sessiya → storage strategiyası kritikdir.

### 2. Sessiya storage variantları
| Variant | Üstünlük | Çatışmazlıq |
|---|---|---|
| **Client-side (cookie)** | server resursu yox, sadə | məhdud yer, təhlükəsizlik riski |
| **Server memory** | sürətli | restart-da İTİRİLİR (volatile) |
| **File-based** | davamlı | çox sessiyada yavaş |
| **DB (PostgreSQL/MySQL)** | etibarlı, davamlı | DB yükü |
| **Redis** | çox sürətli | external asılılıq |

**Client-side cookie nümunəsi:**
```go
http.SetCookie(w, &http.Cookie{
    Name:    "session_token",
    Value:   "some_token_value",
    Expires: time.Now().Add(72 * time.Hour),
})
```
**Server-side (gorilla/sessions filesystem):**
```go
store := sessions.NewFilesystemStore("", []byte("secret-key"))
session, _ := store.Get(r, "session-name")
session.Values["user_id"] = "12345"
session.Save(r, w)
```
**Redis (redistore):**
```go
store, _ := redistore.NewRediStore(10, "tcp", ":6379", "", []byte("secret-key"))
session, _ := store.Get(r, "session-name")
session.Values["user_id"] = "12345"
session.Save(r, w)
```
**Tövsiyə:** Redis (sürət) + DB (persistence) kombinasiyası — hibrid.

### 3. Təhlükəsiz cookie idarəsi
**Təhdidlər:** XSS (skript cookie oğurlayır), CSRF (3-cü tərəf sorğusu cookie ilə göndərilir).

**HTTPOnly — XSS qoruması:**
```go
http.SetCookie(w, &http.Cookie{
    Name:     "session_token",
    Value:    "some_token_value",
    HTTPOnly: true,   // JavaScript cookie-ə ÇATA BİLMƏZ
})
```
**Secure — yalnız HTTPS:**
```go
http.SetCookie(w, &http.Cookie{
    Name:   "session_token",
    Value:  "some_token_value",
    Secure: true,     // yalnız şifrələnmiş bağlantıda ötürülür
})
```
**SameSite — CSRF qoruması:**
```go
http.SetCookie(w, &http.Cookie{
    Name:     "session_token",
    Value:    "some_token_value",
    SameSite: http.SameSiteStrictMode,   // yalnız eyni saytdan
})
```
**İnzibariyyət qaydaları:**
- Cookie dəyərlərini encrypt et (`crypto/aes`) — intercepted olsa belə oxunmaz
- Expiration təyin et (`Expires: time.Now().Add(2 * time.Hour)`) — sonsuz cookie = uzun istismar pəncərəsi
- HTTPS hər yerdə — man-in-the-middle qarşısı
- Müəllifin sitatı (Bruce Schneier): "Security is not a product, but a process" (Təhlükəsizlik məhsul yox, prosesdir)

### 4. İstifadəçi autentifikasiyası (bcrypt ilə)
**Qeydiyyat — parolu hash-lə saxla:**
```go
import "golang.org/x/crypto/bcrypt"

func RegisterUser(username, password string) error {
    hashedPassword, err := bcrypt.GenerateFromPassword(
        []byte(password), bcrypt.DefaultCost)
    if err != nil { return err }
    // hash-lənmiş parolu DB-yə yaz
}
```
**Login — hash müqayisəsi:**
```go
func LoginUser(username, enteredPassword string) bool {
    storedHashedPassword := /* DB-dən al */
    err := bcrypt.CompareHashAndPassword(
        storedHashedPassword, []byte(enteredPassword))
    return err == nil    // nil = uğurlu
}
```
**Logout:** sessiya tokenini serverdən sil + client-də cookie-ni sil.
**Sessiya idarəsi:** login sonrası token yarat → cookie ilə göndər → növbəti sorğularda token yoxlanılır.

### 5. Parol qorunması — 5 qat
**Hashing:** biryollu transformasiya — hash-dən parol geri qaytarıla bilməz; login zamanı hash(hash(girilən)) müqayisəsi.

**Salting (duzlama):** Hər parola unikal təsadüfi string əlavə → rainbow table hücumlarını məhv edir; eyni parollar fərqli hash-lər verir.
```go
func hashPasswordWithSalt(plainPassword, salt string) (string, error) {
    combinedPassword := salt + plainPassword
    hashed, err := bcrypt.GenerateFromPassword(
        []byte(combinedPassword), bcrypt.DefaultCost)
    return string(hashed), err
}
```
**Pepper (istiot):** Hash əvvəlinə gizli sabit dəyər — DB-dən AYRI (offline/env) saxlanılır → DB oğurluğu kifayət etmir:
```go
const pepper = "your-secret-pepper"
combinedPassword := salt + plainPassword + pepper
```
**Iterations:** minlərlə dəfə hash → brute-force-i praktik olmayan yavaşladır (istifadəçi üçün hiss olunmaz).
**Up-to-date alqoritmlər:** bcrypt bu komplekslikləri özü idarə edir; crypto standartları izlə, miqrasiyaya hazır ol.

**Əlavə qaydalar:** HTTPS, security headers (XSS qarşı), asılılıq auditləri.

### 6. OAuth — "Login with Google"
**Nədir:** Token-əsaslı standart protokol — parol paylaşılmadan 3-cü tərəf data çıxışı.

**Axın:** app → provider-in login səhifəsinə redirect → istifadəçi login edir → app-ə authorization code ilə qayıdır → code access token-a exchange olunur → token ilə user profili alınır.

**Konfiqurasiya:**
```go
import (
    "golang.org/x/oauth2"
    "golang.org/x/oauth2/google"
)

conf := &oauth2.Config{
    ClientID:     "YOUR_CLIENT_ID",        // provider-dən qeydiyyatda
    ClientSecret: "YOUR_CLIENT_SECRET",
    RedirectURL:  "https://gitforgits.com/oauth/callback",
    Scopes:       []string{"profile", "email"},
    Endpoint:     google.Endpoint,
}
```
**Login + callback handler-ləri:**
```go
http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
    url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
    http.Redirect(w, r, url, http.StatusFound)    // provider-ə yönləndir
})

http.HandleFunc("/oauth/callback", func(w http.ResponseWriter, r *http.Request) {
    code := r.URL.Query().Get("code")
    token, err := conf.Exchange(oauth2.NoContext, code)   // code → token
    // token ilə user profilini al, bookstore-da auth et
})
```
**Token ömrü:** Access token qısaömürlüdür → **refresh token** ilə yenilənir (redirect-siz). OAuth-da həmişə HTTPS.

## Əsas terminlər
- Session (sessiya) — istifadəçi-app vəziyyət davamlılığı
- Cookie — brauzerdə saxlanan data fraqmenti
- HTTPOnly / Secure / SameSite — cookie qoruma atributları
- bcrypt — parol hash kitabxanası (DefaultCost)
- Hashing / Salting / Pepper / Iterations — parol qoruma qatları
- Rainbow table — hazır hash cədvəli hücumu
- CSRF (Cross-Site Request Forgery) — saxta sorğu hücumu
- OAuth (Open Authorization) — token-əsaslı 3-cü tərəf auth
- Access/Refresh token — qısa/uzunömürlü token cütü
- AuthCodeURL/Exchange — OAuth kod-token mübadiləsi

## Praktik nəticə
1. Cookie-lərdə HƏMİŞƏ üçlük: HTTPOnly + Secure + SameSite — XSS, sniffing, CSRF-ni birlikdə bağlayır.
2. Parollar yalnız bcrypt hash + salt; pepper DB-dən ayrı saxla; plain text = cinayət.
3. Sessiya storage scale-edə bilsin: Redis + DB hibrid; memory-only restart-da itir.
4. OAuth-da redirect + callback + code-exchange axınını düz qur; refresh token ilə sessiyanı uzat.
5. Cookie-lərdə expiration mütləq təyin et — sonsuz cookie uzun hücum pəncərəsidir.
6. Təhlükəsizlik prosesdir: asılılıqları audit et, crypto standartlarını izlə, miqrasiyaya hazır ol.

## Mənbə
Pages: 186-205 (PDF səh. 186-205)
