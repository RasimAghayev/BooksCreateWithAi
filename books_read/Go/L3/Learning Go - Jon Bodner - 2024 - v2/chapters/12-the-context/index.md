# Chapter 12 — The Context (Kontekst)

## Bu chapter nədən bəhs edir?

Context-in məqsədi (request metadata), kontekst-in ötürülmə konvensiyası, HTTP
middleware/client ilə inteqrasiya, cancellation, timeout/deadline, öz kodda cancellation
dəstəyi, WithValue ilə metadata və unexported key pattern.

## Əsas fikirlər

### 1. Context Nədir?
**Problem:** Serverlər request-in meta-məlumatına ehtiyac duyar: (a) işləmək üçün lazım
olan data (tracking ID), (b) nə vaxt dayanmaq (timeout/cancel). Digər dillər threadlocal
istifadə edir — Go-da goroutine identity-i YOXDUR və threadlocal "magiya"dır (dəyər bir
yerdə girir, başqa yerdə peyda olur).

**Həll:** `context.Context` — adi interface; funksiyanın **İLK parametri** kimi açıq
ötürülür (error-un sonuncu olması kimi konvensiya; adi adı `ctx`). Gözəl Go fəlsəfəsi:
yeni dil xüsusiyyəti YOX — sadəcə tip + konvensiya.

```go
func logic(ctx context.Context, info string) (string, error) { ... }
ctx := context.Background()   // giriş nöqtəsi (CLI main və s.) — boş başlanğıc
```
`context.TODO` — inkişaf zamanı yer tutucusu (production-da QALMASIN).

**İmmutabil wrap pattern:** Context-ə məlumat əlavə etmək = valideyni uşaq context-lə
örərək yeni instansiya yaratmaq. Context məlumatı **dərinə** ötürür — kənara YOX.

### 2. HTTP Server: Middleware → Handler
`http.Handler` interfeysi context qəbul etmir (compatibility promise) — ona görə
`http.Request`-ə 2 metod əlavə edilib: `Context()` (oxu) və `WithContext(ctx)`
(yeni request qaytarır — request də immutabildir):

**Kitabdan kod nümunəsi:**
```go
func Middleware(handler http.Handler) http.Handler {
    return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
        ctx := req.Context()
        // ctx-i wrap et (dəyərlər/cancel/timeout)
        req = req.WithContext(ctx)
        handler.ServeHTTP(rw, req)
    })
}
```
Handler daxilində: `ctx := req.Context()` → biznes məntiqinə birinci parametr kimi ötür.

**Outgoing HTTP:** başqa servisi çağırarkən də `req.WithContext(ctx)` və ya
`http.NewRequestWithContext(ctx, ...)` — context sorğu ilə birgə axır.

### 3. Cancellation
**Nədir:** Bəzi goroutine-lər işə düşsə və biri uğursuz olsa — digərlərini dayandırmaq.

**context.WithCancel:**
```go
ctx, cancel := context.WithCancel(parentCtx)
defer cancel()   // MÜTLƏQ — yoxsa RESURS LEAK (yaddaş + goroutine)
```
- Qaytarılan ctx uşaq contextdir; `cancel` — hamıya "dayan" siqnalı.
- Cancel-i birdən çox çağırmaq zərərsizdir (2+-ci çağırışlar iqnor).
- **defer cancel həmişə yaz** — error olsun-olmasın.

**Kitabdan kod nümunəsi (iki paralel çağırış, biri error verəndə hamı dayanır):**
```go
func callBoth(ctx context.Context, errVal, slowURL, fastURL string) {
    ctx, cancel := context.WithCancel(ctx)
    defer cancel()
    var wg sync.WaitGroup
    wg.Add(2)
    go func() {
        defer wg.Done()
        if err := callServer(ctx, "slow", slowURL); err != nil {
            cancel()          // digər goroutine-i də dayandırır
        }
    }()
    go func() {
        defer wg.Done()
        if err := callServer(ctx, "fast", fastURL+"?error="+errVal); err != nil {
            cancel()
        }
    }()
    wg.Wait()
}
// slow server tərəfində: "context canceled" error-u ilə request dərhal qırılır
```

### 4. Timers — Timeout və Deadline
**Load idarəetməsinin 4 yolu** (ilk 3-ü Go-da var): eyni vaxtlı request sayı (goroutine
limiti), növbə (buffered channel), request müddəti (context), resurslar (özün yaz).

**context.WithTimeout / WithDeadline:**
```go
ctx, cancel := context.WithTimeout(parentCtx, 2*time.Second)
defer cancel()
ctx2, cancel2 := context.WithDeadline(parentCtx, deadlineTime)
```
Keçmiş vaxt verilsə — context əvvəlcədən canceled yaranır. `Deadline()` metodu —
(time.Time, bool) comma-ok: timeout qurulubmu?

**Uşaq valideyndən böyük timeout verə bilər, amma valideyn müddəti üstünlük təşkil edir:**
valideyn 2s + uşaq 3s → uşaq 2s-də bitir (done channel 2s sonra bağlanır). Bu, servislər
zəncirində vaxt büdcəsinin bölünməsinə imkan verir (cəmi 50ms → çağırış A 20ms, B 10ms...).

### 5. Öz Kodda Cancellation Dəstəyi
Context-in 2 metodu:
- **Done()** — `chan struct{}`; cancel/timer zamanında **bağlanır** (bağlı channel dərhal
  zero value verir). Cancel-olunmayan context-də Done() **nil qaytarır** — nil channel
  oxumaq sonsuz asılır → yalnız select case-daxilində işlət!
- **Err()** — nil (aktiv), `context.Canceled` (explicit cancel), `context.DeadlineExceeded`
  (timeout).

**Pattern (buffered nəticə channel + select):**
```go
func longRunningThingManager(ctx context.Context, data string) (string, error) {
    type wrapper struct {
        result string
        err    error
    }
    ch := make(chan wrapper, 1)     // buffer 1 — goroutine də cancel olsa yaza bilir → leak yox
    go func() {
        result, err := longRunningThing(ctx, data)
        ch <- wrapper{result, err}
    }()
    select {
    case data := <-ch:
        return data.result, data.err
    case <-ctx.Done():
        return "", ctx.Err()
    }
}
```

### 6. Values — Request Metadata
**Default: açıq parametrlər.** Context dəyər ötürmək yalnız açıq ötürmək MÜMKÜN OLMAYANDA —
klassik hal: HTTP handler və middleware sabit imzalıdır (`http.Handler` 2 parametr).

**context.WithValue / Value:**
```go
ctx = context.WithValue(ctx, key, value)   // uşaq context qaytarır
v := ctx.Value(key)                          // zəncirdə axtarış (linear!), tapılmasa nil
```
Onlarla dəyər saxlamaq — O(n) axtarış + refactor siqnalı.

**Unexported key pattern (toqquşma qorunması) — MÜTLƏQ idiom:**
```go
type userKey int          // unexported tip
const key userKey = 1     // unexported konstanta (iota ilə bir neçə dəyər üçün)

func ContextWithUser(ctx context.Context, user string) context.Context {
    return context.WithValue(ctx, key, user)
}
func UserFromContext(ctx context.Context) (string, bool) {
    user, ok := ctx.Value(key).(string)   // comma-ok assertion
    return user, ok
}
```
Adlandırma konvensiyası: `ContextWithX` (yaz) / `XFromContext` (oxu). String açarlar /
public tiplər FƏRQLİ paketlərdə toqquşa bilər →调试 çətin bug-lar.

**Qayda:** Biznes datanı handler-də çıxar → biznes məntiqinə **açıq parametr** kimi ötür —
context "dəyər qaçırma" yolu DEYİL.

**İstisna — sistem idarəetmə məlumatı (tracking GUID):** biznes state-i deyil, inteqrasiya
parametri artırmır, üçüncü tərəf kod görmür. Logger + RequestDecorator pattern:
```go
type Logger struct{}
func (Logger) Log(ctx context.Context, message string) {
    if guid, ok := guidFromContext(ctx); ok {
        message = fmt.Sprintf("GUID: %s - %s", guid, message)
    }
    fmt.Println(message)
}
func Request(req *http.Request) *http.Request {
    if guid, ok := guidFromContext(req.Context()); ok {
        req.Header.Add("X-GUID", guid)   // servislərarası izləmə
    }
    return req
}
```
DI ilə birləşdiriləndə biznes məntiqi tracking-dən tam xəbərsiz qalır: məntiq yalnız
`bl.Logger.Log(ctx, ...)` və `bl.RequestDecorator(req)` çağırır.

**Tövsiyə (kitab):** Standart API-lərdən dəyər keçirmək üçün context işlət; biznes emal
üçün lazım olan dəyərləri context-dən açıq parametrə köçür; sistem texniki məlumatı
(GUID) birbaşa context-dən oxu.

## Əsas terminlər
- Context — request metadata daşıyan, ilk parametr kimi ötürülən interface
- Cancellation — uğursuz budaqların digər paralel işləri dayandırması
- CancelFunc — context-i ləğv edən funksiya (defer ilə mütləq çağırılmalı)
- Deadline — avtomatik ləğv anı (time.Time)
- Done channel — ləğv zamanı bağlanan struct{} kanalı
- Unexported key pattern — dəyər açarlarının unexported tip/konstantası
- Request decorator — outgoing request-i (GUID header kimi) zənginləşdirən funksiya

## Praktik nəticə

Context qaydaları: (1) `ctx` həmişə ilk parametr, açıq ötür; (2) hər With* funksiyasının
cancel-ını `defer` ilə qoru — leak yaratma; (3) HTTP: `req.Context()` / `req.WithContext()`
/ `NewRequestWithContext`; (4) timeout-u valideyndən uşağa böl — uşaq deadline-i heç vaxt
valideyndən gec bitmir; (5) öz uzun əməliyyatlarında select + Done/Err; (6) WithValue
yalnız açıq ötürmək mümkün olmayan yerlərdə (middleware), unexported açarla, `ContextWith/
FromContext` API-si ilə; (7) biznes dəyəri context-dən çıxarıb parametr et, GUID kimi
texniki metadata context-də qalsın.

## Mənbə
Pages: 346-366
