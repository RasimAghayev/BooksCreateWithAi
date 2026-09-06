# Chapter 9 — Xətalar, panic/recover, veb-server (səh. 41-45)

## Bu chapter nədən bəhs edir?

Xətaların error interfeysi ilə detallı ötürülməsi (PathError nümunəsi), type assertion ilə xəta analizinin incəliyi, panic-in yeri (kritik, bərpaedilməz hallar), recover-in goroutine-qoruyucu və paket-daxili parse-error patternləri və tam nəticə nümunəsi — QR-link veb-serveri.

---

## Əsas fikirlər

### 1. Xətalar — error interfeysi

```go
type error interface {
    Error() string
}
```

Kitabxana müəllifi daha zəngin model də tətbiq edə bilər — detal + kontekst. **os.PathError nümunəsi:**

```go
type PathError struct {
    Op   string  // "open", "unlink" ...
    Path string  // fayl yolu
    Err  error   // sistem çağırışının xətası
}
func (e *PathError) Error() string {
    return e.Op + " " + e.Path + ": " + e.Err.Error()
}
// "open /etc/passwx: no such file or directory"
```

Belə xəta uzaqdan çap olunsa belə məlumatlıdır. **Konvensiya:** xəta sətri MƏNBƏYİ tanıtsın — prefiks (op/paket adı): `"image: unknown format"`.

### 2. Xətanın incə analizi — type assertion

Détallı maraqlanan caller type switch/assertion ilə xüsusi xətaları ayırır — bərpa cəhdi:

```go
for try := 0; try < 2; try++ {
    file, err = os.Create(filename)
    if err == nil {
        return
    }
    if e, ok := err.(*os.PathError); ok && e.Err == syscall.ENOSPC {
        deleteTempFiles()    // yer azad et
        continue              // təkrar cəhd
    }
    return
}
```

`e, ok := err.(*os.PathError)` — uğursuzsa ok=false, e=nil.

### 3. Panic — kritik, bərpaolunmaz hallar

Xəta qaytarma norma olsa da, bəzən proqram DAVAM EDƏ BİLMƏZ. `panic(arg)` — runtime xətası, proqram dayanır (növbəti bölmə — tam deyil!).

```go
// Converge olmayan kub kök:
func CubeRoot(x float64) float64 {
    z := x / 3
    for i := 0; i < 1e6; i++ {
        prevz := z
        z -= (z*z*z - x) / (3 * z * z)
        if veryClose(z, prevz) {
            return z
        }
    }
    panic(fmt.Sprintf("CubeRoot(%g) did not converge", x))
}

// Məcburi asılılıq yoxdursa:
var user = os.Getenv("USER")
func init() {
    if user == "" {
        panic("no value for $USER")
    }
}
```

**Prinsip:** Kitabxana panic-i SON ÇARƏ kimi istifadə etməli; problemin ört-basdır/barma yolu varsa davam etmək daha yaxşıdır.

### 4. Recover — goroutine xilası

Panic (daxili runtime xətaları — index-out-of-range, uğursuz assertion — daxil) funksiyanı dərhal dayandırır və **stack unwind**-ə başlayır — yoldakı bütün defer-ləri icra edir. Goroutine stack zirvəsinə çatsa proqam ölür. **`recover` unwind-i dayandırır və panic arqumentini qaytarır** — yalnız defer daxilində işlək.

**Serverdə çökan goroutine-in izolyasiyası:**

```go
func server(workChan <-chan *Work) {
    for work := range workChan {
        go safelyDo(work)
    }
}
func safelyDo(work *Work) {
    defer func() {
        if err := recover(); err != nil {
            log.Println("work failed:", err)   // qeyd et və tərk et
        }
    }()
    do(work)
}
```

Bir goroutine çöker — digərləri davam edir. recover defer-siz çağrılsa həmişə nil qaytarır → defer içindəki kod çağıran kitabxanaların öz panic/recover-undan qorxmur.

### 5. Parse-error patterni (paket-daxili panic+recover)

Regexp parser nümunəsi — daxildə panic, xaricə error:

```go
// Yerli xəta tipi:
type Error string
func (e Error) Error() string { return string(e) }

func (regexp *Regexp) error(err string) {
    panic(Error(err))    // daxili səhv siqnalı
}

func Compile(str string) (regexp *Regexp, err error) {
    regexp = new(Regexp)
    defer func() {
        if e := recover(); e != nil {
            regexp = nil       // named nəticələr defer-dən DƏYİŞİLƏ BİLƏR
            err = e.(Error)    // parse xətası deyilsə YENİDƏN panic (re-panic)
        }
    }()
    return regexp.doParse(str), nil
}
```

Parser daxilində sadə çağırış:

```go
if pos == 0 {
    re.error("'*' illegal at start of expression")
}
```

Bu pattern **YALNIZ paket daxilində** — xarici klientlərə panic SİZMƏZ; parse funksiyası panic→error çevirir. Re-panic orijinalın report-da görünməsinə mane olmur (hər ikisi çap olunur); yalnız orijinal lazımdırsa — filtr + orijinal ilə yenidən panic.

### 6. Veb-server — tam nümunə (QR generator)

```go
package main

import (
    "flag"
    "html/template"
    "log"
    "net/http"
)

var addr = flag.String("addr", ":1718", "http service address") // Q=17, R=18
var templ = template.Must(template.New("qr").Parse(templateStr))

func main() {
    flag.Parse()
    http.Handle("/", http.HandlerFunc(QR))    // funksiya→HandlerFunc çevirməsi!
    err := http.ListenAndServe(*addr, nil)
    if err != nil {
        log.Fatal("ListenAndServe:", err)
    }
}

func QR(w http.ResponseWriter, req *http.Request) {
    templ.Execute(w, req.FormValue("s"))
}

const templateStr = `
<html><head><title>QR Link Generator</title></head><body>
{{if .}}
<img src="http://chart.apis.google.com/chart?chs=300x300&cht=qr&choe=UTF-8&chl={{.}}" />
<br>{{.}}<br><br>
{{end}}
<form action="/" name=f method="GET">
    <input maxLength=1024 size=70 name=s value="" title="Text to QR Encode">
    <input type=submit value="Show QR" name=qr>
</form>
</body></html>
`
```

**Öyrənilən elementlər:** flag (`:1718` — telefon hərfləri Q=17/R=18!); `template.Must` (parse-time panic məqbuldur — başlanğıcda xəta); HandlerFunc çevirməsi (ch6); html/template avtomatik escapinqlə təhlükəsiz `{{.}}`; `{{if .}}` — boş formada img gizlənir.

Əlbəttə ki, bu bir neçə sətirlək faydalı veb-server + HTML şablonudur — Go az sətirlə çox iş görür. (Qeyd: chart.apis.google.com API artıq mövcud deyil — QR generasiyası işləmir, amma pattern aktual.)

---

## Əsas terminlər

| Termin | İzah |
|--------|------|
| error interfeysi | `Error() string` — istənilən zəngin modellə genişlənə bilər |
| PathError | Op + Path + Err — mənbəyi tanıdan xəta strukturu |
| Prefiks konvensiyası | "image: unknown format" — mənbə paketi |
| Xəta-filtri assertion | `err.(*os.PathError)` + ENOSPC yoxlaması → retry |
| panic | Runtime ölümcül xəta — unwind başlayır |
| Stack unwinding | Panic-in yuxarı qaçışı — yolda defer-lər icra |
| recover | Defer-daxili — unwind dayandırır, panic dəyəri qaytarır |
| safelyDo | Serverdə goroutine izolyasiyası — log + davam |
| Named result + defer | Compile: `regexp = nil; err = e.(Error)` |
| Re-panic | Beklenməyən panic-i buraxmaq — yalnız öz Error-un yutmaq |
| template.Must | Başlanğıcda xəta = panic məqbul (sürətli uğursuzluq) |

---

## Praktik nəticə

1. **Xətaları kontekstlə qaytar:** PathError kimi struct — op+path+sistem xətası; uzaq çapda belə məlumatlı. Sətir prefiksi mənbəyi göstərsin.
2. **Xəta tipini ayıraraq bərpa et:** comma-ok assertion + daxili Err yoxlaması → müvəqqəti xətalarda (ENOSPC) təmizlə+retry.
3. **Panic = son çarə:** Kitabxanalarda bərpa yolu varsa panic ETMƏ; yalnız davam mümkünsüzdür (init asılılıq, converge xətası).
4. **Server goroutine-lərini qoru:** `defer recover()` + log — bir işin çöküşü serveri öldürməz.
5. **Paket-daxili panic idiomu:** Parser sadə `re.error(...)` çağırsın; public funksiya defer+recover ilə error-a çevirsin; yalnız öz Error tipini yudla (re-panic).
6. **template.Must + flag:** Başlanğıc zamanı ölümcül xəta → panic MƏQBULDUR; HandlerFunc çevirməsi funksiyanı handler edir.
7. **html/template avtomatik escapinq:** İstifadəçi mətni səhifəyə təhlükəsiz keçir — XSS qoruması.

---

## Mənbə

- Sənəd: *Effective Go* (rus tərcüməsi), 2009
- PDF səhifələri: 41-45
