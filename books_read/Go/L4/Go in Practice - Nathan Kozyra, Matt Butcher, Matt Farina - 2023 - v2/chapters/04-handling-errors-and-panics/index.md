# Chapter 4 — Handling errors and panics (Xətaların və panic-lərin idarəsi)

## Bu chapter nədən bəhs edir?

Error konvensiyası (son qaytarma, zero value ilə birgə), custom error tipləri, sentinel
error dəyişənləri, error wrapping (%w, Unwrap/Is), panic vs error fəlsəfəsi, panic-ə error
ötürmək, defer/recover, deferred closure scope qaydaları, adlı return ilə cleanup,
goroutine-lərdə panic (handler server idiom, safely.Go).

## Əsas fikirlər

### 1. Error Konvensiyası — Son Qaytarma + Zero Value
**Kitabdan kod nümunəsi:**
```go
func Concat(parts ...string) (string, error) {
    if len(parts) == 0 {
        return "", errors.New("No strings supplied")   // zero value + error
    }
    return strings.Join(parts, " "), nil
}
```
- `error` HƏMİŞƏ sonuncu qaytarmadır (köhnə kodda pozulmuş ola bilər)
- **Nil best practice:** xəta zamanı mümkünsə İŞLƏNƏBİLƏN dəyər qaytar ("" Concat üçün —
  strings.Join/Split kimi). İstifadəçi error-u iqnor etsə belə kod işləyə bilər; amma ""
  həm xəta, həm də qanuni input olduğu üçün ayırdetmə YALNIZ error yoxlaması ilə mümkündür.
- Qaydalı iqnor: `result, _ := Concat(args...)` — yalnız dəyər əhəmiyyətli olanda.
  Heç nə faydalı qaytarıla bilmirsə — nil.
- Şərh konvensiyası: normal + xəta davranışını sənədləşdir.
- Go-dev-lər custom error TİP-lərini az istifadə edir — "error is an error"; ConcatError
  tipi heç nə qazandırmır.

### 2. Custom Error Tipləri
**Nə vaxt:** Xətanın detalı istifadəçini FƏRQLİ kod yazmağa sövq etməlidir:
```go
type ParseError struct {
    Message    string
    Line, Char int
}
func (p *ParseError) Error() string {
    return fmt.Sprintf("%s on Line %d, Char %d", p.Message, p.Line, p.Char)
}
```
Parse xətasının YERİ (sətir/simvol) — xəta sətrini işıqlandırmaq kimi UI imkanları verir.

### 3. Sentinel Error Dəyişənləri
**Nədir:** Paket-səviyyəli `errors.New` dəyişənləri; məna daşıyan konkret instansiyalar
(`io.EOF` kimi).

**Kitabdan kod nümunəsi (timeout/retry pattern):**
```go
var ErrTimeout = errors.New("The request timed out")
var ErrRejected = errors.New("The request was rejected")

response, err := SendRequest("Hello")
if errors.Is(err, ErrTimeout) {
    timeouts := 0
    for err == ErrTimeout {
        timeouts++
        fmt.Println("Timeout. Retrying.")
        if timeouts == MAX_TIMEOUTS {
            panic("too many timeouts!")
        }
        response, err = SendRequest("Hello")
    }
}
```
Üstünlüklər: bir dəfə instantiate; `==`/`errors.Is` ilə sadə müqayisə; Java-nın exception
class-atma modelinin Go-vari effektli əvəzi. Qeyd: Go 1.20+ `rand.Seed` avtomatikdir.

### 4. Wrapping — %w, Unwrap, Is
```go
return "", fmt.Errorf("we got an error: %w ", ErrTimeout)   // %w = wrap
```
- Wrap-dan sonra `err == ErrTimeout` MÜTLƏQ `false`! (daha çox "annotation" tipidir)
- `errors.Unwrap(err)` — bir səviyyə açır (nil qaytarə bilər)
- **`errors.Is(err, ErrTimeout)`** — rekursiv zəncir axtarışı; praktikada HƏMİŞƏ bunu işlət

### 5. Panic vs Error Fəlsəfəsi
- **Error:** gözlənilən problemi bildirir — sənədləşdirilə, kodda görünə; developer
  cavabdehdır; iqnor edilirsə Go heç nə etmir (sonrakı addımlara yanlış state düşə bilər).
- **Panic:** sistemin (sub)sistem davam edə BİLMƏYİCƏYİ hal; stack unwind edir, recover
  yoxdursa proqramı ÖLDÜRÜR + stack trace verir.
- Divide nümunəsi: precheckDivide (error qaytarır) vs `divide(2, 0)` → runtime panic
  "integer divide by zero" — kontrolsüz hal sistemin bacara bilmədiyi vəziyyətə düşür.
- **Qızıl qayda:** indiki kontekstdə aydın idarə yolu yoxdursa panic; mümkündür — error.

### 6. Panic-ə Nə Ötürməli?
`panic(interface{})` — hər şey olar (nil, string, error). **İdiomatik: error ötür**:
1. Intuitivdir — səbəb error-dursa, error
2. Recover-də error kimi tutmaq asanlaşır (`fmt.Printf("Error: %s", thePanic)`)
Error bubble-up zəncirinin sonu: `Initialize() → StartServer() → main()` — main error-u
panic-ə çevirir.

### 7. defer və Closure Scope
```go
defer func() { ... }()   // defer + anonim funksiya (closure) + dərhal çağırış ()
```
**Scope qaydaları:**
- Closure özündən ƏVVƏL elan olunan dəyişənləri GÖRÜR (runtime dəyəri ilə):
```go
var msg string
defer func() { fmt.Println(msg) }()
msg = "Hello world"     // "Hello world" çap olunur
```
- Sonradan elan olunanı GÖRMÜR — compile xətası:
```go
defer func() { fmt.Println(msg) }()
msg := "Hello world"    // compile error!
```

### 8. Recover + Cleanup + Adlı Return (OpenCSV Pattern)
**Kitabdan kod nümunəsi:**
```go
func OpenCSV(filename string) (file *os.File, err error) {   // ADLI return!
    defer func() {
        if r := recover(); r != nil {
            file.Close()          // panic olsa belə resurs təmizliyi
            err = r.(error)       // panic-i error-a ÇEVİR və adlı return-la qaytar
        }
    }()
    file, err = os.Open(filename)
    if err != nil { return file, err }
    RemoveEmptyLines(file)       // panic edir
    return file, err
}
```
Adlı return mütləqdir — defer daxilisində err yenilənir və qaytarma dəyəri dəyişir.
Library sərhədində panic-i error-a çevirmək — ən yaxşı təcrübə.

**defer qaydaları (siyahı):**
1. defer-i funksiyanın ən BAŞINA yaxın qoy
2. Inline olmasa da olar (metod çağırışı kimi)
3. Sadə təyinatlar (`foo := 1`) defer-dən əvvəl ola bilər
4. Mürəkkəb dəyişənlər elan ƏVVƏL, init SONRA (var myFile io.Reader)
5. Çoxlu defer-lərdən qaç; mütləqdisə — TƏRS sıra ilə icra olunur
6. Fayl/şəbəkə/DB resurslarını defer-də bağla

### 9. Goroutine-lərdə Panic — Server Idiom
**Problem:** Goroutine öz funksiya stack-ına malikdir; oradakı recover-süz panic BAŞQA
stack-ə keçə BİLMƏZ → proqram ölür. listen-də recover yazar versevimiz işləməz.

**Həll 1 — handler-da recover (echo server):**
```go
func handle(conn net.Conn) {
    defer func() {
        if err := recover(); err != nil {
            fmt.Printf("Fatal error: %s", err)
        }
        conn.Close()            // hər halda bağla
    }()
    reader := bufio.NewReader(conn)
    data, err := reader.ReadBytes('\n')
    ...
    response(data, conn)         // panic etsə də server YAŞAYIR
}
```

**Həll 2 — server handler idiom (interfeys + wrap):** net/http.Server kimi — istifadəçi
handler-i təmin edir, kitabxana onu recover-li wrapper ilə goroutine-də icra edir
(handlerChain: "if r := recover(); r != nil { log }"; handler-i goroutine-də çağır).

**Həll 3 — safely.Go (Masterminds/cookoo):**
```go
safely.Go(message)   // goroutine + daxili panic trap + log
```
Qayda: kitabxanalar panic-in proqramı öldürməsinə structura qur; developer-in yaddaşına
güvənmə. Qeyd: closure ilə dəyişən ötürsən race riski — `-race` flag (go run/build)
konkurent kod üçün MÜTLƏQ alət.

## Əsas terminlələr
- Sentinel error — mənalı paket-səviyyəli error instansiyası
- `%w` — wrap verb (Unwrap/Is-i mümkün edən)
- errors.Is — wrap zəncirində rekursiv axtarış
- Panic/recover — ölümcül hal + yalnız defer-də tutulma
- Deferred closure — defer edilmiş anonim funksiya; əvvəlki scope-u görür
- Adlı return (named return) — defer-in qaytarma dəyərini dəyişməsi üçün
- Handler server idiom — istifadəçi handler-i recover-li wrapper-də icra
- GoDoer/safely.Go — panic-trapping goroutine starter

## Praktik nətidə

Xəta qərar ağacı: (1) gözlənilən problem → error (son qaytarma, zero-value və ya işlənəbilən
dəyərlə); (2) ayrı-seçkilə lazım oldu → sentinel var (errors.Is ilə); (3) detal zərurəti →
custom struct error; (4) kontekst əlavəsi → %w wrap (müqayisə üçün errors.Is, == YOX);
(5) davam mümkünsüz → panic(YALNIZ error dəyəri ilə); (6) kitabxanada panic → adlı return
+ defer recover + error-a çevir; (7) goroutine başlan hər yerdə recover wrapper (safely.Go
kimi); (8) defer qaydalarına riayət — resurs sızması qarşısı; (9) konkurent kodu -race
ilə test et.

## Mənbə
Pages: 87-114 (PDF 108-135)
