# Chapter 12 — Graceful Shutdown

## Bu fəsil nədən bəhs edir?

Shutdown siqnalları (SIGINT/SIGTERM/SIGQUIT/SIGKILL), `signal.Notify` ilə
tutulması və `http.Server.Shutdown()` ilə in-flight sorğuların tamamlanmasına
imkan verən tədrici dayandırma.

## Əsas fikirlər

### 1. Unix siqnalları
| Signal | Təsvir | Qısayol | Catchable? |
|---|---|---|---|
| SIGINT | Klaviatura interrupt | Ctrl+C | Bəli |
| SIGQUIT | Quit + stack dump | Ctrl+\ | Bəli (default saxlanılır) |
| SIGTERM | Sıralı terminasiya | `pkill -SIGTERM` | Bəli |
| SIGKILL | Dərhal öldür | `pkill -SIGKILL` | **YOX** |

**Sub-kod izahı:**
- Catchable → proqram tərəfindən tutulub idarə oluna bilər (graceful
  shutdown tetikler); SIGKILL heç vaxt tutulmur — həmişə dərhal öldürür
- SIGQUIT default davranışda saxlanılır — "əvvəlki davranış" qorunur
  (Ctrl+\ → stack dump ilə dərhal çıxış; debug üçün faydalı)
- `pgrep -l api` / `pkill -SIGTERM api` → process tap/siqnal göndər

### 2. serve() metodunun ayrılması
Server kodu main()-dən ayrılır (Chapter 20-də bunun faydası görünəcək):

```go
// cmd/api/server.go
func (app *application) serve() error {
    srv := &http.Server{
        Addr:         fmt.Sprintf(":%d", app.config.port),
        Handler:      app.routes(),
        IdleTimeout:  time.Minute,
        ReadTimeout:  10 * time.Second,
        WriteTimeout: 30 * time.Second,
    }
    app.logger.PrintInfo("starting server", map[string]string{
        "addr": srv.Addr,
        "env":  app.config.env,
    })
    return srv.ListenAndServe()
}
```
- main() → `err = app.serve()` + PrintFatal

### 3. Siqnal tutma — signal.Notify
**Kitabdan kod nümunəsi:**
```go
go func() {
    quit := make(chan os.Signal, 1) // BUFFERED — mütləq!
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    s := <-quit // siqnal gələnədək bloklanır
    app.logger.PrintInfo("shutting down server", map[string]string{
        "signal": s.String(),
    })
    // ...
}()
```

**Sub-kod izahı:**
- `signal.Notify(quit, SIGINT, SIGTERM)` → sadəcə göstərilən siqnallar
  kanala yönləndirilir; digərləri default davranışda qalır
- **Buffer size 1 MÜTLƏQDİR:** signal.Notify göndirəndə receiver gözləmir —
  unbuffered kanalda siqnal receiver hazır deyilsə İTİRİLƏR
- Background goroutine → ömrü proqram boyu; `<-quit` gələnədək bloklu
- `s.String()` → "interrupt" / "terminated"

### 4. Graceful shutdown — Shutdown() + ErrServerClosed
**Kitabdan kod nümunəsi (tam serve):**
```go
func (app *application) serve() error {
    srv := &http.Server{ /* ... */ }

    shutdownError := make(chan error)

    go func() {
        quit := make(chan os.Signal, 1)
        signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
        s := <-quit
        app.logger.PrintInfo("shutting down server", map[string]string{
            "signal": s.String(),
        })

        ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
        defer cancel()

        shutdownError <- srv.Shutdown(ctx)
    }()

    app.logger.PrintInfo("starting server", map[string]string{
        "addr": srv.Addr,
        "env":  app.config.env,
    })

    err := srv.ListenAndServe()
    if !errors.Is(err, http.ErrServerClosed) {
        return err // "əsl" xəta (port busy və s.)
    }

    err = <-shutdownError // Shutdown bitənədək gözlə
    if err != nil {
        return err
    }

    app.logger.PrintInfo("stopped server", map[string]string{"addr": srv.Addr})
    return nil
}
```

**Sub-kod izahı:**
- `srv.Shutdown(ctx)`: listener-ləri bağlayır → idle bağlantıları qapatır →
  aktiv bağlantılar idle-ya düşüb qapanmağa "sonsuz" gözləyir — bizim 20s
  context bunu məhdudlaşdırır
- **`http.ErrServerClosed` yararsız deyil:** Shutdown çağırılanda
  `ListenAndServe()` dərhal bu xətanı qaytarır → graceful shutdown
  BAŞLADIĞININ göstəricisidir; yalnız DİGƏR xətalar error sayılır
  (`errors.Is` ilə yoxlanır)
- `shutdownError` kanalı → Shutdown-in nəticəsi ana goroutine-ə çatdırılır
  (deadline vurulsa xəta gəlir)
- Axın: siqnal → Shutdown başlayır → ListenAndServe ErrServerClosed ilə
  qayıdır → `<-shutdownError` gözləyir → "stopped server" → main()
  təmiz çıxır (defer db.Close() işləyir!)

**Nə baş vermir (vacib limit):**
- Background task-lar gözləMİR
- Hijacked uzunömürlü bağlantılar (WebSocket) qapanMIR
- Bunlar üçün öz koordinasiya məntiqiniz lazımdır (sonrakı fəsillərdə)

**Demo:** healthcheck-ə 4s `time.Sleep` qoyub `curl ... & pkill -SIGTERM api`:
"shutting down server" dərhal, "stopped server" 4 saniyə sonra — in-flight
sorğu tamamlanır, client cavab ALIR.

## Bu fəsildə layihəyə nə əlavə olundu? (project continuity)

**Yeni fayllar:**
- `cmd/api/server.go` — serve() + signal goroutine + Shutdown + ErrServerClosed triage

**Dəyişdirilən:** main.go → sadeleşdi (app.serve() çağırışı)

## Əsas terminlər

- Signal (siqnal) — OS-dən prosesə asinxron bildiriş
- Graceful Shutdown (tədrici dayandırma) — aktiv işlərin tamamlanması gözlənilərək
  dayanma
- In-flight Request (işlənməkdə olan sorğu) — hələ cavablanmamış sorğu
- Buffered Channel (buferli kanal) — göndərici gözləməyən kanal (siqnal
  itkisinin qarşısı)
- Grace Period (gözəl müddət) — shutdown üçün icazə verilən vaxt (20s context)
- Hijacked Connection (ələ keçirilmiş bağlantı) — server idarəsindən çıxmış
  (WebSocket kimi) bağlantı

## Praktik nəticə

`ErrServerClosed`-in "xəta olmadığını" bilmək və buffered signal kanalı —
bu fəslin iki ən vacib texniki detalıdır. 20s grace period deployment
strateqiyasının hissəsidir: systemd/Kubernetes SIGTERM göndərir, app təmiz
şəkildə bitər. SIGQUIT/SIGKILL isə "panic button" olaraq qalır.

## Mənbə
Pages: 263-275 (raw 263-275)
