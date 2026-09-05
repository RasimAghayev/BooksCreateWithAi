# Chapter 2 — A solid foundation (Technique 1-9)

## Bu chapter nədən bəhs edir?

4 foundational sahə: CLI flag-lar (flag paketi + GNU-style kitabxanalar + cli.go framework), konfiqurasiya (JSON/YAML/INI fayllar + environment variables), web server start/shutdown (init daemon + graceful shutdown - manners), URL routing (multi-handler → path wildcard → regex → 3rd-party router-lar).

## Əsas fikirlər

### 1. Go flag paketi — Plan 9 üslubu
**Fərq:** Go flag sistemi Plan 9-dan gəlir — GNU/Linux/BSD-dən FƏRLİ:
- `ls -la` → Go-da `-la` TƏK flag kimi oxunur (qruplaşdırma YOX)
- Uzun seçələr: GNU `--color` (iki tire) ↔ Go `-color` (tək tire, qısa/uzun ayrı-seçilməz)

**Kitabdan kod nümunəsi:**
```go
package main

import (
    "flag"
    "fmt"
)

var name = flag.String("name", "World", "A name to say hello to.")
var spanish bool

func init() {
    flag.BoolVar(&spanish, "spanish", false, "Use Spanish language.")
    flag.BoolVar(&spanish, "s", false, "Use Spanish language.")   // qısa ad — EYNİ dəyişənə!
}

func main() {
    flag.Parse()                     // parse — dəyərlər dəyişənlərə düşür
    if spanish == true {
        fmt.Printf("Hola %s!\n", *name)      // pointer dereference!
    } else {
        fmt.Printf("Hello %s!\n", *name)
    }
}
```

**2 üsul:**
1. `flag.String("name", "World", "...")` → pointer qaytarır (`*name` ilə çıxış)
2. `flag.BoolVar(&var, ...)` → mövcud dəyişənə yazar; qısa+uzun ad üçün 2 dəfə çağır

**Help imkanları:**
- `flag.PrintDefaults()` → avtomatik help: `-name string / A name to say hello to. (default "World")`
- `flag.VisitAll(func(f *flag.Flag) {...})` → custom help
- Arqumentlər (flag olmayan): `flag.Args()` / `flag.Arg(i)`

### TECHNIQUE 1: GNU/UNIX-style flag-lar
**Problem:** İstifadəçilər UNIX üslubu gözləyir (`-la` qrup, `--long-opt`).

**Həll — 2 kitabxana:**

**a) launchpad.net/gnuflag** — flag paketi ilə API-uyğun (drop-in əvəz):
- `-f` (qısa), `-fg` (qrup!), `--flag`, `--flag x`, `-f x`, `-fx`
- Fərq: `Parse(true|false, args)` — true = flag-lar hər yerdə axtarılır (`foo -bar` -dən sonra da)

**b) github.com/jessevdk/go-flags** — tam fərqli API, struct əsaslı:
```go
var opts struct {
    Name    string `short:"n" long:"name" default:"World" description:"A name to say hello to."`
    Spanish bool   `short:"s" long:"spanish" description:"Use Spanish Language"`
}

func main() {
    flags.Parse(&opts)           // reflection ilə struct tag-lərdən oxuyur!
    fmt.Printf("Hello %s!\n", opts.Name)
}
```
+ Option groups, yaxşı help, `-p/usr/local` / `-p /usr/local` / `-p=/usr/local`, eyni optionun çoxdəfə → slice.

### TECHNIQUE 2: CLI framework — cli.go
**Problem:** hər yeni CLI tool-da eyni boilerplate.

**Həll:** github.com/urfave/cli — Docker, Cloud Foundry, Drone istifadə edir.

**Sadə app:**
```go
package main

import (
    "fmt"
    "os"
    "gopkg.in/urfave/cli.v1"
)

func main() {
    app := cli.NewApp()
    app.Name = "hello_cli"
    app.Usage = "Print hello world"
    app.Flags = []cli.Flag{
        cli.StringFlag{
            Name:  "name, n",
            Value: "World",
            Usage: "Who to say hello to.",
        },
    }
    app.Action = func(c *cli.Context) error {
        name := c.GlobalString("name")
        fmt.Printf("Hello %s!\n", name)
        return nil
    }
    app.Run(os.Args)
}
```
- `--help` / `--version` AVTOMATİK; help text AVTOMATİK
- `cli.NewExitError` → nonzero exit code

**Komanda + subkomanda (git kimi):**
```go
app.Commands = []cli.Command{
    {
        Name:      "up",
        ShortName: "u",
        Usage:     "Count Up",
        Flags: []cli.Flag{
            cli.IntFlag{Name: "stop, s", Usage: "Value to count up to", Value: 10},
        },
        Action: func(c *cli.Context) error {
            start := c.Int("stop")     // komanda flag-ı — Global YOX
            for i := 1; i <= start; i++ {
                fmt.Println(i)
            }
            return nil
        },
    },
    // ... "down" komandası
}
```
- Struktur: `$ app [global options] command sub-command [command options] [arguments]`
- Global flag-lar → `c.GlobalString`; komanda flag-ları → `c.String`
- Default Action = help göstər

### TECHNIQUE 3: Konfiqurasiya faylları (3 format)

**JSON:**
```go
type configuration struct {
    Enabled bool
    Path    string
}

file, _ := os.Open("conf.json")
defer file.Close()
decoder := json.NewDecoder(file)
conf := configuration{}
err := decoder.Decode(&conf)
fmt.Println(conf.Path)     // /usr/local
```
+ Standart kitabxana, etcd; − şərhlər YOX.

**YAML** (github.com/kylelemons/go-gypsy/yaml):
```go
config, err := yaml.ReadFile("conf.yaml")
fmt.Println(config.Get("path"))
fmt.Println(config.GetBool("enabled"))
```
+ Şərhlər, oxunaqlı; müəllif tövsiyəsi.

**INI** (gopkg.in/gcfg.v1):
```go
config := struct {
    Section struct {
        Enabled bool
        Path    string
    }
}{}
err := gcfg.ReadFileInto(&config, "conf.ini")
fmt.Println(config.Section.Enabled)
```
+ io.Reader qəbul edir, tag dəstəyi.

### TECHNIQUE 4: Environment variables (12-factor)
**12-factor:** konfiq PaaS-larda fayl YOX — ENV-də:
```go
http.ListenAndServe(":"+os.Getenv("PORT"), nil)
```
- `os.Getenv` → tapılmadıqda boş string; çevirmə: `strconv.ParseInt`
- **XƏBƏRDARLIQ:** subprocess-lər ENV-i GÖRÜR — sensitive data ehtiyatlı!

**12 faktor (qısa):** 1 codebase; 2 asılılıqlar declare; 3 konfiq ENV-də; 4 backing services; 5 build/run ayrı; 6 stateless proseslər; 7 port binding; 8 horizontal scale; 9 sürətli startup + graceful shutdown; 10 dev/staging/prod bənzər; 11 log = event stream; 12 admin task ayrıca.

### 2.3 Web server start/shutdown
**Doğru yol — init daemon** (systemd/upstart/init/launchd):
```bash
$ systemctl start myapp.service
$ systemctl stop myapp.service
```
**Qayda:** Go app daemon YAZMA — init daemon tərəfindən idarə olun!

**ANTIPATTERN — callback URL (/shutdown):**
```go
http.HandleFunc("/shutdown", func(res http.ResponseWriter, req *http.Request) {
    os.Exit(0)             // dərhal — data İTİR!
})
```
− production-da silinməli (security); − davam edən işlər kəsilir; − ops tooling-i bypass edir.

### TECHNIQUE 5: Graceful shutdown — manners
**Problem:** standard http dərhal ölür — data yazılmır, connection-lar kəsilir.

**Həll:** github.com/braintree/manners (PayPal) — WaitGroup ilə connection-ları izləyir.

**Kitabdan kod nümunəsi:**
```go
func main() {
    handler := newHandler()
    ch := make(chan os.Signal)
    signal.Notify(ch, os.Interrupt, os.Kill)     // OS siqnallarını dinlə
    go listenForShutdown(ch)
    manners.ListenAndServe(":8080", handler)      // http ilə EYNİ interface
}

func listenForShutdown(ch <-chan os.Signal) {
    <-ch                       // siqnal gözlə
    manners.Close()            // yeni connection ALMA → mövcudları bitir → qapat
}
```

**Üstünlüklər:** request-lər tamamlanır; port sərbəst → yeni versiya başlaya bilər (zero-downtime!).
**Məhdudiyyət:** yalnız HTTP; socket handoff YOX; ayrıca goroutine-lər üçün öz WaitGroup lazımdır.

### TECHNIQUE 6: Multi-handler routing (Sadə)
```go
http.HandleFunc("/hello", hello)          // YALNIZ /hello
http.HandleFunc("/goodbye/", goodbye)     // /goodbye/ + alt-path-lər
http.HandleFunc("/", homePage)            // fallback (ən az spesifik)

func goodbye(res http.ResponseWriter, req *http.Request) {
    path := req.URL.Path
    parts := strings.Split(path, "/")
    name := parts[2]                       // /goodbye/Buttercup → Buttercup
    ...
}
func homePage(res http.ResponseWriter, req *http.Request) {
    if req.URL.Path != "/" {
        http.NotFound(res, req)           // 404 helper
        return
    }
    fmt.Fprint(res, "The homepage.")
}
```
**Qaydalar:**
- Ən spesifikdən az spesifikə resolve
- `/goodbye/` → `/goodbye` redirect (query string İTƏ BİLƏR!)
- `/hello` → yalnız dəqiq `/hello`
- Query: `req.URL.Query().Get("name")` — yoxdursa ""
- Method ayrımı: `req.Request.Method` yoxla (GET/POST eyni handler-a düşür!)

**Pro:** sənədləşmiş, sadə. **Con:** method+path ayrıseçilməz, wildcard yox, 404 hər handler-da əl ilə.

### TECHNIQUE 7: path paketi — wildcard router
```go
pr := newPathResolver()
pr.Add("GET /hello", hello)
pr.Add("* /goodbye/*", goodbye)      // method wildcard + path wildcard
http.ListenAndServe(":8080", pr)

type pathResolver struct {
    handlers map[string]http.HandlerFunc
}

func (p *pathResolver) ServeHTTP(res http.ResponseWriter, req *http.Request) {
    check := req.Method + " " + req.URL.Path      // "GET /hello"
    for pattern, handlerFunc := range p.handlers {
        if ok, err := path.Match(pattern, check); ok && err == nil {
            handlerFunc(res, req)
            return
        }
    }
    http.NotFound(res, req)
}
```
- `pathResolver` özü Handler-dir (ServeHTTP implement) → Server-a birbaşa verilir
- `path.Match` → POSIX pattern; `*` bir seqment qədər (`foo/*` ≠ `foo/bar/baz`)
- `* /goodbye/*` → istənilən method

**Pro:** standart kitabxana. **Con:** `*` slash-dan öncə dayanır; trailing slash xüsusi hal (yalnız `/goodbye/`, `/goodbye` 404).

### TECHNIQUE 8: Regex router
```go
rr.Add("GET /hello", hello)
rr.Add("(GET|HEAD) /goodbye(/?[A-Za-z0-9]*)?", goodbye)   // GET|HEAD + optional ad

type regexResolver struct {
    handlers map[string]http.HandlerFunc
    cache    map[string]*regexp.Regexp          // compile cache!
}

func (r *regexResolver) Add(regex string, handler http.HandlerFunc) {
    r.handlers[regex] = handler
    cache, _ := regexp.Compile(regex)
    r.cache[regex] = cache                       // BİR dəfə compile
}

func (r *regexResolver) ServeHTTP(res http.ResponseWriter, req *http.Request) {
    check := req.Method + " " + req.URL.Path
    for pattern, handlerFunc := range r.handlers {
        if r.cache[pattern].MatchString(check) {    // cache-dən match
            handlerFunc(res, req)
            return
        }
    }
    http.NotFound(res, req)
}
```
- `(GET|HEAD)`, optional group `(/?[A-Za-z0-9]*)?` → `/goodbye`, `/goodbye/`, `/goodbye/name` hamısı
- **Compile cache** — hər request-də yenidən compile YOX (performans!)
- 1-ci match qalib (sıra vacib)

**Pro:** tam nəzarət. **Con:** oxunmaz → test LAZIM.

### TECHNIQUE 9: 3rd-party router-lar
| Paket | Xüsusiyyət |
|---|---|
| **httprouter** (julienschmidt) | ən sürətli; minimal memory; case-insensitive; `../` cleanup; optional trailing `/` |
| **gorilla/mux** | versatile: host, scheme, header, path, query matching |
| **pat** (bmizerany) | Sinatra-inspired; `/user/:name` named params; gorilla/pat fork |

Sinatra → 50+ framework-ə ilham vermiş Ruby library.

## Router müqayisəsi
| Yanaşma | Method ayrımı | Wildcard | Test | Tövsiyə |
|---|---|---|---|---|
| Multi-handler (T6) | ✗ (manual) | ✗ | asan | sadə saytlar |
| path.Match (T7) | ✓ (`* /...`) | ✓ (1 seqment) | orta | kiçik REST |
| Regex (T8) | ✓ | ✓✓ (tam) | çətin | kompleks URL |
| 3rd-party (T9) | ✓ | ✓ | — | production REST |

## Əsas terminlər
- Plan 9 flag üslubu (tək dash, qrup yox)
- flag.String / flag.BoolVar / flag.Parse / PrintDefaults / VisitAll
- gnuflag / go-flags (struct tag reflection)
- cli.go: NewApp / Commands / Action / GlobalString / cli.Context
- 12-Factor App (konfiq ENV-də)
- Init Daemon (systemd/upstart/init/launchd)
- Callback Shutdown URL (antipattern)
- Graceful Shutdown / manners / WaitGroup connection tracking
- signal.Notify / os.Interrupt / os.Kill
- Handler (ServeHTTP interfeysi) / pathResolver / regexResolver
- path.Match (POSIX) vs regexp
- Compile Cache
- httprouter / gorilla/mux / pat
- http.NotFound / http.Error

## Praktik nəticə
- İstifadəçilərin UNIX gözləntisi varsa — flag paketi YOX, gnuflag/go-flags/cli istifadə et; struct-tag approach ən təmiz.
- CLI layihəsi böyüyəcəksə BAŞLANĞICDA cli.go ilə başla — komanda/subkomanda/help/versiya pulsuz.
- Konfiq: 12-factor ENV (PaaS üçün) yoxsad fayl (YAML şərhlərlə ən rahat) — hər ikisini dəstəklə.
- `/shutdown` URL YAZMA — init daemon + graceful shutdown (manners / müasir: http.Server.Shutdown) doğru yol.
- Öz router yazarkən regex-ləri ƏVVƏLCƏDƏN compile et və cache-lə — request başına compile tələsdir.
- Sadə app → multi-handler; REST API → method-aware router (T7/T8/3rd-party); production → httprouter/mux kimi sınaqdan keçmiş paket.

## Mənbə
Pages: 50-81 (PDF), book pages 27-58
