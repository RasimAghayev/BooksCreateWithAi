# Chapter 8 — Working with Errors (səh. 125-137)

## Bu fəsil nədən bəhs edir?

Xətaların dəyər kimi idarəolunması: %w ilə wrapping, panic-in recover ilə
tutulması (named returns), goroutine panic-ləri (safelyGo), errors.Is/As
ilə xəta yoxlaması və runtime.Callers ilə stack əlavə edən Wrapper.

## Əsas fikirlər

### Recipe 42 — xətaların qaytarılması
**Tapşırıq:** PID faylını oxu → prosesi öldür → faylı sil.

```go
func killServer(pidFile string) error {
    file, err := os.Open(pidFile)
    if err != nil {
        return fmt.Errorf("can't open PID file: %w", err)   // %w = WRAP
    }
    defer file.Close()

    var pid int
    _, err = fmt.Fscanf(file, "%d", &pid)
    if err != nil {
        return fmt.Errorf("bad PID in %q: %w", pidFile, err)
    }

    if err := os.Remove(pidFile); err != nil {
        log.Printf("can't remove %q - %s", pidFile, err)   // warning — xəta YOX
    }

    return kill(pid)
}

func kill(pid int) error {
    proc, err := os.FindProcess(pid)
    if err != nil {
        return err
    }
    return proc.Kill()
}
```
- `%w` → orijinal xətanı saxlayır → `errors.Is/As` işləyir
- Fayl silinməsi warning-dır (design decision) — diskussiya edin
- `_` ilə ignore — güclü tövsiyə olunmur

### Recipe 43 — panic-in tutulması (named returns)
**Tapşırıq:** xarici `dist.Edit` panic edir — patch etmək riskli; sargı
ilə error-a çevir.

```go
// EditDistance Returns the edit (Levenshtein) distance between s1 and s2.
// It wraps dist.Edit against panics.
func EditDistance(s1, s2 string) (distance int, err error) {   // NAMED!
    defer func() {
        if e := recover(); e != nil {
            err = fmt.Errorf("%v", e)    // any → error
        }
    }()
    return dist.Edit(s1, s2), nil
}
```
- **Named return** mütləqdir — defer daxilindən `err`-i dəyişmək
  başqa yolu yoxdur
- Uğurda err=nil qalır; panic-də recover e-ni err-ə yazır
- Fail-fast məktəbi: monitorinq + recovery ilə birlikdə çox etibarlıdır
  (Joe Armstrong "Systems that Run Forever Self-Heal and Scale")

### Recipe 44 — goroutine panic-ləri (safelyGo)
**Problem:** `go handler(msg)` panic edəndə BÜTÜN proqram çökür!

```go
// safelyGo will run fn in a goroutine, and guard it from panics
func safelyGo(fn func()) {
    go func() {
        defer func() {
            if err := recover(); err != nil {
                log.Printf("error: %s", err)
            }
        }()
        fn()
    }()
}

func drain(ch <-chan Message, handler func(Message)) {
    for msg := range ch {
        msg.Time = time.Now()
        safelyGo(func() {
            handler(msg)          // panic burada tutulur
        })
    }
}
```

**Niyə goroutine panic-i proqramı öldürür?**
- Digər dillərdə thread çökəndə proqram pis vəziyyətdə DAVAM edir →
  sonrakı çökmənin səbəbi tapmaq çətinləşir
- Go fəlsəfəsi: pis vəziyyətdə davam etməkdənsə çökmək yaxşıdır
- Həll: safelyGo sargısı + Docker/Kubernetes restart qatı + monitorinq
- **Qeyd:** handler özü YENİ goroutine başlatdığında (http nümunəsindəki
  kimi) sargı kömür etmir — heç bir bulletproof yol yoxdur

### Recipe 45 — xətaların yoxlanması (errors.Is)
**Tapşırıq:** fayllar siyahısından PID oxu-öldür; Fayl yoxdursa IGNORE,
digər xətalar qayıtsın.

```go
func killFromFiles(logFiles []string) error {
    for _, logFile := range logFiles {
        err := killServer(logFile)
        if err == nil {
            return nil                      // uğur — bitdi
        }
        if !errors.Is(err, os.ErrNotExist) { // fayl yoxdurmu?
            return err                       // xeyr — real xəta
        }
        // fayl yoxdur → növbəti fayla keç
    }
    files := strings.Join(logFiles, ", ")
    return fmt.Errorf("no existing file found: %v", files)
}
```
- `errors.Is(err, os.ErrNotExist)` → wrap zənciri boyunca axtarır
- `errors.As` → xətanın konkret TİP olub-olmamasını yoxlayır
- errors paketi sənədləri + nümunələrlə oynayın

### Recipe 46 — xətalara stack əlavə edilməsi
**Tapşırıq:** Python-un call stack-i kimi kontekst — hansı fayl/sətir?

```go
// Wrapper wraps an error with call stack information.
type Wrapper struct {
    error                  // embed → error interfeysi hazır
    stack []uintptr        // PC-lər (gecikmiş)
}

func (w *Wrapper) Unwrap() error {
    return w.error          // errors.Is/As uyğunluğu
}

// Frame is a call location.
type Frame struct {
    Function string
    File     string
    Line     int
}

func (f Frame) String() string {
    return fmt.Sprintf("%s:%d: %s", trimPath(f.File, 3), f.Line, f.Function)
}

// Wrap wraps an error with call stack information.
func Wrap(err error) error {
    const depth = 32
    var pcs [depth]uintptr
    n := runtime.Callers(2, pcs[:])    // 2 = Wrap+Callers-u ötür
    w := Wrapper{
        error: err,
        stack: pcs[:n],
    }
    return &w
}

// Stack returns the call stack, innermost frame first.
func (w *Wrapper) Stack() []Frame {
    locs := make([]Frame, 0, len(w.stack))
    frames := runtime.CallersFrames(w.stack)
    for {
        frame, more := frames.Next()
        loc := Frame{Function: frame.Function, File: frame.File, Line: frame.Line}
        locs = append(locs, loc)
        if !more {
            break
        }
    }
    return locs
}
```

**Dizayn qərarları:**
- Wrap yalnız PC-ləri saxlayır (ucuz) — fayl/sətir çözümü Stack()
  çağırılınca (runtime.CallersFrames) — performans həssas yerlər üçün
- `runtime.Callers(2, ...)` — 2 daxili funksiyanı ötür (0 olarsa
  istifadəçiləri qarışdırar)
- Stack trace-i LOG-a yazın, istifadəçiyə yalnız mesaj — "crash" təəssüratı
  yaratmasın
- pkg/errors (Dave Cheney) — Go 1.13 errors dizaynına təsir etdi; arxivlənib

## Final Thoughts-dən — xəta haqqında 5 sual

1. Bu xətanı özüm həll edim, yoxsa yuxarı ötürüm?
2. Aydın mesaj necə formatlıyım?
3. Orijinalı wrap edim? (adətən BƏLİ)
4. Qayıtmazdan əvvəl bağlanmalı resurs varmı? (defer)
5. Panic-i burada tutmalıyam?

## Əsas terminlər
- Errors are values — Go proverbü; xətalar ucuz, dəyər kimi ötürülür
- %w verb — wrap edən format
- errors.Is / errors.As — zəncir yoxlaması
- Named return — defer-dən cavabı dəyişmək üçün
- recover — yalnız deferred daxilində işləyir
- Fail fast — xətalı vəziyyətdə davam etməmək fəlsəfəsi
- safelyGo — goroutine panic qoruması
- runtime.Callers / CallersFrames — PC → stack frame
- Deferred resolution — ağır işi sonra saxlamaq
- Wrapper pattern — error + stack

## Praktik nəticə
Hər xətada: %w ilə wrap + kontekst; errors.Is/As ilə konkret səbəb yoxlaması
(fayl yoxdur kimi halları süzmək üçün); qeyri-kritik uğursuzluqlar (fayl
silinməsi) log-la warning. Xarici kodun panic-indən qorunmaq üçün named
return + defer/recover sargısı; goroutine-lərdə safelyGo. Stack konteksti
lazımdırsa Wrap idiomu — amma yalnız lazım olan yerlərdə və trace-i log-da.

## Mənbə
Pages: 125-137 (PDF 125-137)
