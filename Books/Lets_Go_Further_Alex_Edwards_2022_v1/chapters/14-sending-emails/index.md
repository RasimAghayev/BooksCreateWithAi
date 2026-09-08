# Chapter 14 — Sending Emails

## Bu fəsil nədən bəhs edir?

SMTP (Mailtrap) ilə email göndərmə, `html/template` + `//go:embed` ilə
template-lərin binary-yə daxil edilməsi, go-mail/mail ilə Mailer paketi,
background goroutine-də email + panic recovery + `sync.WaitGroup` ilə
graceful shutdown koordinasiyası.

## Əsas fikirlər

### 1. SMTP setup — Mailtrap
**Nədir:** Development üçün SMTP sandbox — göndərilən email-lər real
alıcıya deyil, Mailtrap inbox-a düşür. Alternativlər: Postmark, SendGrid,
Amazon SES, öz SMTP server.

### 2. Email template-ləri — html/template
**Kitabdan kod nümunəsi:**
```
{{define "subject"}}Welcome to Greenlight!{{end}}
{{define "plainBody"}}
Hi,
Thanks for signing up for a Greenlight account...
For future reference, your user ID number is {{.ID}}.
{{end}}
{{define "htmlBody"}}
<!doctype html>
<html>...
    <p>For future reference, your user ID number is {{.ID}}.</p>
...</html>
{{end}}
```

**Sub-kod izahı:**
- Bir faylda 3 named template: subject / plainBody / htmlBody
- `{{.ID}}` → dinamik data (User struct) render olunur
- Alternativ: template-lər DB-də string kimi (tez-tez dəyişən/user-editable
  hallarda); fayl sadə başlanğıcdır

### 3. Mailer paketi — go:embed + go-mail
**Kitabdan kod nümunəsi:**
```go
// internal/mailer/mailer.go
//go:embed "templates"
var templateFS embed.FS

type Mailer struct {
    dialer *mail.Dialer
    sender string
}

func New(host string, port int, username, password, sender string) Mailer {
    dialer := mail.NewDialer(host, port, username, password)
    dialer.Timeout = 5 * time.Second
    return Mailer{dialer: dialer, sender: sender}
}

func (m Mailer) Send(recipient, templateFile string, data any) error {
    tmpl, err := template.New("email").ParseFS(templateFS, "templates/"+templateFile)
    if err != nil {
        return err
    }

    subject := new(bytes.Buffer)
    err = tmpl.ExecuteTemplate(subject, "subject", data)
    // ... plainBody, htmlBody eyni qayda ilə ...

    msg := mail.NewMessage()
    msg.SetHeader("To", recipient)
    msg.SetHeader("From", m.sender)
    msg.SetHeader("Subject", subject.String())
    msg.SetBody("text/plain", plainBody.String())
    msg.AddAlternative("text/html", htmlBody.String())

    return m.dialer.DialAndSend(msg)
}
```

**Sub-kod izahı:**
- **`//go:embed`** qaydaları: yalnız package-level global var-larda işləyir
  (funksiya daxilində compile xətası); path directive-i olan fayla nisbidir;
  `.`/`..`/başlanğın `/` qadağandır; directory → rekursiv embed (`.`/`_`
  başlayanlar istisna — wildcard `templates/*` lazımdır); forward slash
  Windows-da belə
- `net/smtp` (stdlib) FROZEN-dur, attachment dəstəyi yoxdur → go-mail/mail
  üstünlük
- `AddAlternative` həmişə `SetBody`-dən SONRA — HTML alternative kimi əlavə
  olunur
- `DialAndSend` → bağlan + göndər + qapat; 5s dialer timeout ("dial tcp:
  i/o timeout")

**Retry pattern (Additional Info):**
```go
for i := 1; i <= 3; i++ {
    err = m.dialer.DialAndSend(msg)
    if nil == err {
        return nil
    }
    time.Sleep(500 * time.Millisecond)
}
return err
```
- `if nil == err` yazılışı (Yoda style) — `if err != nil` ilə vizual
  qarışıklığın qarşısını almaq üçün; background-da sleep daha uzun ola bilər

### 4. SMTP konfiqurasiyası + istifadə
```go
type config struct {
    // ...
    smtp struct {
        host     string
        port     int
        username string
        password string
        sender   string
    }
}

flag.StringVar(&cfg.smtp.host, "smtp-host", "smtp.mailtrap.io", "SMTP host")
// ... 5 flag

mailer: mailer.New(cfg.smtp.host, cfg.smtp.port, ...)
```
- registerUserHandler-da: `app.mailer.Send(user.Email, "user_welcome.tmpl", user)`
- Senxron göndəriş latency: ~2.3s (client gözləyir) → növbədi addım bunu
  aradan qaldırır

### 5. Background email — 202 Accepted
**Kitabdan kod nümunəsi:**
```go
go func() {
    err = app.mailer.Send(user.Email, "user_welcome.tmpl", user)
    if err != nil {
        app.logger.PrintError(err, nil) // serverErrorResponse YOX!
    }
}()
err = app.writeJSON(w, http.StatusAccepted, envelope{"user": user}, nil)
```

**Sub-kod izahı:**
- **202 Accepted** → "sorğu qəbul edildi, emal tamamlanmayıb" — async
  əməliyyatın düzgün status kodu
- Xəta idarəsi MÜTLƏQ `PrintError` — `serverErrorResponse` ikinci HTTP cavab
  yazmağa cəhd edərdi: `"http: superfluous response.WriteHeader call"`
- Closure `user` və `app`-i yaxınlaşdırır — dəyişsələr hər iki tərəfə
  görünür (biz dəyişmirik, amma diqqət!)
- Latency: 2.33s → 0.27s

### 6. background() helper — panic recovery ilə
**Problem:** Background goroutine-dəki panic recoverPanic middleware və
http.Server tərəfindən TUTULMUR → app çökür (Chapter 10-un vədi burada
yerinə yetir).

**Kitabdan kod nümunəsi:**
```go
func (app *application) background(fn func()) {
    app.wg.Add(1)             // WaitGroup artım (14.5-də əlavə olundu)
    go func() {
        defer app.wg.Done()   // bitəndə azalt
        defer func() {
            if err := recover(); err != nil {
                app.logger.PrintError(fmt.Errorf("%s", err), nil)
            }
        }()
        fn()
    }()
}
```

**Handler-də istifadə:**
```go
app.background(func() {
    err = app.mailer.Send(user.Email, "user_welcome.tmpl", user)
    if err != nil {
        app.logger.PrintError(err, nil)
    }
})
```

**Sub-kod izahı:**
- First-class functions: `fn func()` — istənilən funksiya parametr kimi
  ötürülür; panic recovery bir yerdə cəmlənir, təkrar-təkrar yazılmır,
  unudulmur

### 7. sync.WaitGroup — background task-ların shutdown-ı
**Problem:** Graceful shutdown background goroutine-ləri GÖZLƏMİR → email
yarımçıq qala bilər.

**WaitGroup konsepti:** counter — Add(1) artır, Done() azaltır, Wait()
sıfırlanana qədər bloklayır.

**Qayda:** `wg.Add(1)` goroutine BAŞLAMAZDAN ƏVVƏL çağrılmalı — goroutine
daxilində çağrılsa Wait() onsuz da işə düşə bilər (race).

**Kitabdan kod nümunəsi:**
```go
type application struct {
    config config
    logger *jsonlog.Logger
    models data.Models
    mailer mailer.Mailer
    wg     sync.WaitGroup // zero-value istifadəyə hazırdır
}
```

**serve() daxilində shutdown zənciri:**
```go
err := srv.Shutdown(ctx)
if err != nil {
    shutdownError <- err
}
app.logger.PrintInfo("completing background tasks", map[string]string{"addr": srv.Addr})
app.wg.Wait()               // bütün background task-lar bitsin
shutdownError <- nil        // uğurlu shutdown siqnalı
```

**Sub-kod izahı:**
- Shutdown-in öz xətası KANALA YALNIZ xəta halında gedir; uğurda isə
  `wg.Wait()`-dən SONRA `nil` göndərilir — ana goroutine hər halda cavab
  alır (deadlock yoxdur)
- Log ardıcıllığı: "caught signal" → "completing background tasks" → (2s
  email gözləmə) → "stopped server"
- Shutdown context 5s (20s-dən endirilib — email + in-flight üçün)

## Bu fəsildə layihəyə nə əlavə olundu? (project continuity)

**Yeni fayllar:**
- `internal/mailer/mailer.go` — Mailer, New, Send, embed.FS
- `internal/mailer/templates/user_welcome.tmpl`

**Dəyişdirilən:**
- main.go → config.smtp + application.mailer + application.wg
- users.go → background email + 202 Accepted
- helpers.go → background() helper (wg.Add/Done + recover)
- server.go → Shutdown sonrası wg.Wait()

**Yeni asılılıq:** `github.com/go-mail/mail/v2`

## Əsas terminlər

- SMTP (Simple Mail Transfer Protocol) — email göndərmə protokolu
- Embedded File System (embed.FS) — binary-yə daxil edilmiş fayllar
- //go:embed directive — compile zamanı faylları binary-yə yerləşdirən
  komanda
- First-Class Function — dəyişən/parametr kimi ötürülə bilən funksiya
- 202 Accepted — emal qəbul edilib, tamamlanmayıb
- sync.WaitGroup — goroutine kolleksiyasının bitməsini gözləmə primitivi
- Superfluous WriteHeader — ikinci cavab yazma xətası (background xətasının
  əsəbkeşliyi)

## Praktik nəticə

Async işlərin tam recepi: `background()` helper-i (recover + WaitGroup) +
202 Accepted + shutdown-da `wg.Wait()` — üçü bir yerdə həm client
təcrübəsini (0.27s), həm data bütövlüyünü (email həmişə gedir), həm də app
stabilliyini (panic tutulur) təmin edir. Bu pattern istənilən async
task üçün (SMS, webhook, image processing) kopyalanabiləndir.

## Mənbə
Pages: 294-322 (raw 294-322)
