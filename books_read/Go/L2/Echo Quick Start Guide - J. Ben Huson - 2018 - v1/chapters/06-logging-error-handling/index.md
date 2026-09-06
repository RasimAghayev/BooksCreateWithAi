# Chapter 6 — Performing Logging and Error Handling

## Bu chapter nədən bəhs edir?

Bu chapter veb-tətbiqetmənin ən çox gözardı edilən, amma kritik iki aspektini işıqlandırır: **Logging (Jurnallama)** və **Error handling (Xətaların idarəsi)**. Echo-nun `Logger` interfeysi, log səviyyələrinin düzgün istifadəsi, Logger middleware, handler-lərdən qaytan xətaların Echo tərəfindən avtomatik HTTP cavabına çevrilməsi (`HTTPErrorHandler`), custom error handler və panic-lərdən qorunma (`Recover` middleware) izah olunur.

> Uğur halının (success case) funksionallığı üzərində çox işlənir; "uğursuz" hallar isə asanlıqla unudulur. Xətalar yaxşı idarə edilməsə, tətbiqetmə **resilient (davamlı)** və **fail with grace (layiqcə uğursuz olacaq)** olmayacaq.

---

## Əsas fikirlər

### 1. Logger — Echo-nun daxili jurnal sistemi

**Nədir?** Echo-nun `Logger` interfeysi `c.Logger()` (context-dən) və ya `e.Logger` (instance-dan) ilə əldə edilir. Handler içərisində istifadə nümunəsi:

```go
// HealthCheck - Health Check Handler
func HealthCheck(c echo.Context) error {
    if reqID, ok := c.Get(middlewares.RequestIDContextKey).(uuid.UUID); ok {
        c.Logger().Debugf("RequestID: %s", reqID.String())
    }
    resp := renderings.HealthCheckResponse{
        Message: "Everything is good!",
    }
    return c.JSON(http.StatusOK, resp)
}
```

`c.Logger()` kontekstdəki Logger implementasiyasını qaytarır; `Debugf` isə debug səviyyəsində formatlı jurnal yazır.

### 2. Echo Logger interfeysi

Interfeysin əsas metodları (gommon `log` paketinə əsaslanır — `github.com/labstack/gommon/log`):

```go
type Logger interface {
    Output() io.Writer
    SetOutput(w io.Writer)
    Prefix() string
    SetPrefix(p string)
    Level() log.Lvl
    SetLevel(v log.Lvl)

    Print(i ...interface{})
    Printf(format string, args ...interface{})
    Printj(j log.JSON)

    Debug(i ...interface{})
    Debugf(format string, args ...interface{})
    Debugj(j log.JSON)

    Info(i ...interface{})
    Infof(format string, args ...interface{})
    Infoj(j log.JSON)

    Warn(i ...interface{})
    Warnf(format string, args ...interface{})
    Warnj(j log.JSON)

    Error(i ...interface{})
    Errorf(format string, args ...interface{})
    Errorj(j log.JSON)

    Fatal(i ...interface{}) / Fatalf(format, args) / Fatalj(j log.JSON)
    Panic(i ...interface{}) / Panicf(format, args) / Panicj(j log.JSON)
}
```

**Xüsusiyyətlər:**
- Hər səviyyənin 3 variantı var: `Xxx` (birbaşa dəyər), `Xxxf` (formatlı), `Xxxj` (JSON strukturlu — structured logging).
- `SetOutput` ilə çıxışın yönünü dəyişmək, `SetPrefix` ilə prefiks, `SetLevel` ilə səviyyə tənzimlənir.
- Interfeys olduğundan **tam əvəz edilə bilər** (extend oluna bilər) — öz logger implementasiyanızı qoşa bilərsiniz.

Başlanğıcda səviyyəni qurmaq (`cmd/service/main.go`):

```go
func main() {
    // create a new echo instance
    e := echo.New()
    e.Logger.SetLevel(log.DEBUG)
    e.Validator = new(bindings.Validator)
    // ...
}
```

---

### 3. Log səviyyələri və düzgün istifadəsi

**Problem:** Developer-lər ya çox az, ya da həddindən artıq log yazır. Yanlış səviyyədə log:
- Vacib xəta mesajları aşağı səviyyədə itə bilər (production-da o səviyyə söndürülür),
- Secret (gizli) məlumatların (parollar, token-lər) log-a düşməsi — **information leakage (məlumat sızması)** təhlükəsi yaradır.

| Səviyyə | Nə vaxt istifadə olunur | Nümunə |
|---------|------------------------|--------|
| **DEBUG** | Developer-in troubleshoot (problem axtarışı) üçün; tətbiqetmənin necə işlədiyini izləmək | Data strukturu daxili vəziyyəti, `RequestID` |
| **INFO** | Xəta bildirməyən, amma əməliyyat üçün faydalı məlumat | Client bağlantı məlumatları, remote address, request metodu, tam path |
| **WARN** | Bərpa olunmuş, amma araşdırılmalı xəta | Bərpa edilmiş bağlantı kəsilməsi |
| **ERROR** | Sistemin bərpa etdiyi, amma təhqiqat tələb edən tətbiq xətası | DB sorğusu xətası |
| **FATAL** | Bərpa mümkün olmayan xəta — tətbiqetmə dayanır | Başlanğıcda DB-yə qoşula bilməmək |
| **PANIC** | Tətbiqetmə davam edə bilmir, panic baş verir | Port-a bind oluna bilməməsi |

**Qaydalar:**
- **DEBUG:** Production-da adətən söndürülür → hər şeyi DEBUG-a qoymaq yaxşı fikir deyil; istisna hallarda vacib məlumat itə bilər.
- **INFO:** Production log səviyyəsi adətən burada durur — debug-ın verbosluğundan qorunur, WARN/ERROR isə buradan yuxarı buraxılır.
- **WARN və yuxarı:** Səviyyə INFO-ya tənzimlənərsə, WARN və ERROR avtomatik emitted olunur (səviyyə iyerarxiyası).
- **FATAL/PANIC:** Yalnız başlanğıc zamanı alınmaması mümkün olan xətalar üçün (məs., DB konnekti, port bind).

---

### 4. Logger middleware — request-lərin avtomatik jurnalı

```go
e.Use(middleware.Logger())
```

Hər request üçün standartlaşdırılmış JSON jurnal sətri yazar:

```json
{"time":"2018-03-27T22:09:44.084450941-04:00","id":"","remote_ip":"::1",
 "host":"localhost:8080","method":"GET","uri":"check","status":200,
 "latency":424727,"latency_human":"424.727µs","bytes_in":0,"bytes_out":33}
```

**Sahələr:** zaman, request ID (`RequestID` middleware ilə doldurulur), klient IP, host, HTTP metodu, URI, status kodu, gecikmə (latency), giriş/çıxış bayt sayı. Bu məlumat operations komandası üçün debugging və load analizində əsasdır.

---

### 5. Error handling — handler-dən qaytan xəta nə olur?

Go konvensiyası (Go blog "Error handling and Go"): xətalar **exception throw etmək yox, explicit (açıq) yoxlanılır**. Echo-da hər handler `error` qaytarır:

```go
package handlers

import (
    "errors"
    "github.com/labstack/echo"
)

// Error - Example Error Handler
func Error(c echo.Context) error {
    return errors.New("failure!")
}
```

Handler `error` qaytardıqda Echo bunu avtomatik HTTP cavabına çevirir:

```bash
curl localhost:8080/error -D -
HTTP/1.1 500 Internal Server Error
Content-Type: application/json; charset=UTF-8
...
{"message":"Internal Server Error"}
```

**Necə işləyir?** Echo instance-ındakı `HTTPErrorHandler` metodu handler/middleware-dən gələn `error`-u qəbul edir və leqal HTTP cavabına çevirir. Default: `echo.DefaultHTTPErrorHandler`.

### 6. Custom HTTP Error Handler

`HTTPErrorHandler` **tam əvəz edilə bilər**:

```go
func myHTTPErrorHandler(err error, c echo.Context) {
    code := http.StatusInternalServerError
    if httpErr, ok := err.(*echo.HTTPError); ok {
        code = httpErr.Code
    }
    c.String(code, ...)
}

e.HTTPErrorHandler = myHTTPErrorHandler
```

**`echo.HTTPError` tipindən istifadə:** Custom error handler yazmasanız belə, handler-də `echo.HTTPError` qaytarsanız (status kodu + mesaj dolu şəkildə), `DefaultHTTPErrorHandler` onu tanıyır və düzgün status + message ilə cavab verir:

```go
return echo.NewHTTPError(http.StatusNotFound, "user not found")
// → 404 {"message":"user not found"}
```

---

### 7. Panic-lərin idarəsi — Recover middleware

**Problem:** Panic baş verərsə və tutulmasa, service istehlakçıları tətbiqetmənin **tam stack trace**-ini cavab olaraq alar — həm utancverici, həm də təhlükəlidir.

**Həll:** Echo-nun `Recover` middleware-i handler-ləri wrap edərək Go-nun built-in `recover()` funksiyası ilə panic-i tutur:

```go
defer func() {
    if r := recover(); r != nil {
        //...
    }
}()
```

**Necə işləyir?** Middleware handler-ləri wrap etdiyindən (Ch4-də öyrənilən prinsip), handler və ya istənilən nested middleware-də yaranan panic `Recover` tərəfindən yoxlanılır:
1. Panic tutulur (`recover()`),
2. Echo framework-ünə generic error qaytarılır,
3. Echo default error handler ilə istehlakçıya layiqli cavab (500) göndərilir — stack trace yox.

**Tövsiyə:** Hər Echo tətbiqetməsini **həmişə** `Recover` middleware ilə wrap edin. Əlavə oxu: Go blog "Defer, Panic and Recover".

---

## Əsas terminlər

| Termin | İzah |
|--------|------|
| Logging (Jurnallama) | Tətbiqetmənin vəziyyəti haqqında strukturlaşdırılmış iz qeydləri — troubleshooting üçün |
| `Logger` interface | Echo-nun jurnal abstraksiyası: səviyyələr, output, prefiks idarəsi; əvəz edilə biləndir |
| `log.Lvl` / `SetLevel` | Jurnal səviyyəsi (DEBUG, INFO, WARN, ERROR, FATAL, PANIC) və onun tənzimi |
| `Xxxj` metodları | JSON strukturlu jurnal yazımı (structured logging) |
| Information leakage | Secret-lərin (parol, token) təsadüfən log-a düşməsi |
| Logger middleware | Hər request üçün standart JSON sətri yazar (method, uri, status, latency) |
| Latency (Gecikmə) | Request-in emal müddəti — Logger middleware-in json sahəsi |
| Explicit error check | Go konvensiyası: exception throw yox, xətanın yaranma yerində açıq yoxlanılması |
| `HTTPErrorHandler` | Handler-dən qaytan error-u HTTP cavabına çevirən, əvəz edilə bilən funksiya |
| `echo.HTTPError` | Status kodu + mesaj daşıyan xəta tipi — default handler tərəfindən tanınır |
| `Recover` middleware | `recover()` ilə panic-i tutub generic error-a çevirən qoruyucu middleware |
| Stack trace (Stek izi) | Panic tutulmasa istehlakçıya düşən daxili kod izi |
| Fail with grace | Xəta halında layiqcə (düzgün cavabla) uğursuz olmaq |
| Resilient (Davamlı) | Xətalara baxmayaraq işləməyə davam edə bilmə qabiliyyəti |

---

## Praktik nəticə

1. **Logger-i vahid nöqtədən idarə edin:** `e.Logger.SetLevel(...)` ilə başlanğıcda bir dəfə qurun; handler-lərdə `c.Logger()` istifadə edin — jurnal infrastrukturunu özünüz yazmayın.
2. **Səviyyə intizamına riayət edin:** DEBUG = troubleshoot, INFO = production default, WARN = bərpa olunmuş problem, ERROR = təhqiqat lazım, FATAL/PANIC = yalnız başlanğıc xətaları. Hər şeyi DEBUG-a qoymayın, secret-ləri heç vaxt log-lamayın.
3. **Logger middleware hər servisdə olsun:** `e.Use(middleware.Logger())` + `RequestID` middleware birlikdə requests-in audit izini avtomatik yaradır.
4. **Handler-lərdə error qaytarmaq kifayətdir:** Echo onu avtomatik 500-ə çevirir; dəqiq status lazımdırsa `echo.NewHTTPError(code, msg)` qaytarın.
5. **Xəta formatını özelleştirin:** `e.HTTPErrorHandler = myHTTPErrorHandler` ilə bütün xətaların cavab formatını vahidləşdirin.
6. **`Recover` mütləq əlavə edin:** `e.Use(middleware.Recover())` — əks halda panic zamanı klientə stack trace düşür.
7. **Handler imzası `func(echo.Context) error`** — error qaytarma Echo-nun idarəetmə zəncirinə xətanı ötürür; özünüz `http.Error` ilə yazmayın.

---

## Mənbə

- Kitab: *Echo Quick Start Guide* — J. Ben Huson, Packt Publishing, 2018 (ISBN 9781789340849)
- Chapter 6: "Performing Logging and Error Handling", book səh. 114–133
- PDF səhifələri: 126–141
- Kod: https://github.com/PacktPublishing/Echo-Essentials/tree/master/chapter6
- Video: https://goo.gl/owjFvG
- Əlavə: Go blog "Error handling and Go" (https://go.dev/blog/error-handling-and-go), Go blog "Defer, Panic and Recover"
- Default error handler kodu: https://github.com/labstack/echo/blob/60f88a7a.../echo.go#L317
