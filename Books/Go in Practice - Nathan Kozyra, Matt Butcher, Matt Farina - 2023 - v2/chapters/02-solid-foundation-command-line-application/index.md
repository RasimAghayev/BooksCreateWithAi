# Chapter 2 — A solid foundation: Building a command-line application (Möhkəm Təməl: Komanda Sətri Tətbiqi)

## Bu chapter nədən bəhs edir?

CLI flag-ları (flag paketi, Cobra framework), enum çatışmazlıqları, konfiqurasiya
(JSON/YAML/INI faylları + environment dəyişənləri, 12-factor), web server-in start/stop
idarəetməsi (init daemonlar, graceful shutdown, OS siqnalları) və URL routing (handler-lar,
ServeMux, regex resolver, üçüncü tərəf routerlər).

## Əsas fikirlər

### 1. Command-line Flag-ları — flag Paketi
**Nədir:** Go-nun flag sistemi Plan 9 stilidir: `-la` BİRLƏŞDİRİLMİR (tək "la" flagı kimi
oxunur); `-` və `--` eynidir; qısa/uzun ayrımı YOXDUR.

**Kitabdan kod nümunəsi:**
```go
var name = flag.String("name", "World", "A name to say hello to.")  // pointer qaytarır
var inSpanish bool
func init() {
    flag.BoolVar(&inSpanish, "spanish", false, "Use Spanish language.")
    flag.BoolVar(&inSpanish, "s", false, "Use Spanish language.")   // eyni dəyişənə 2 ad
    flag.Parse()
}
func main() {
    fmt.Printf("Hello %s!\n", *name)   // pointer dereference
}
```

**Qaydalar:**
- `flag.String/Bool/Int...` — yeni dəyişən qaytarır (pointer); `flag.StringVar/BoolVar...`
  — mövcud dəyişənə yazır
- Hər tip üçün funksiya var; `time.Duration` istisna — integer-dan avtomatik cast
- Yanlış flag → avtomatik usage çapı + exit; `flag.PrintDefaults()`, `flag.VisitAll(cb)`
  custom help üçün; `flag.Args()` — parse olunmamış qalıq arqumentlər
- Eyni dəyişənə qısa+uzun verdikdə SON deklarasiya qalib gəlir (sıra ilə parse)

### 2. Enum Qeyri-ciddiliyi
Go-da enum keyword YOXDUR; `type language = string` + konstantlar sadəcə aliasdır —
`--lang="xy"` kimi yanlış dəyər SÖZSÜZ keçir. **Həll:** keçərli dəyərlər slice-da +
`init`-də validasiya:
```go
validLanguages = []string{"en", "sp", "fr", "de"}
if !slices.Contains(validLanguages, userLanguage) {
    log.Fatalf("Invalid language %s. Please use one of %v", userLanguage, validLanguages)
}
```

### 3. Cobra Framework — Real CLI
**Nədir:** Docker/Kubernetes-in istifadə etdiyi CLI framework; subcommand, short+long
flag, help, autocomplete.

**Kitabdan kod nümunəsi (hesablayıcı, add/sub/mul/div subcommand-ları):**
```go
var cmdAdd = &cobra.Command{
    Use:   "add",
    Short: "Add two numbers",
    Run: func(cmd *cobra.Command, args []string) {
        result, err := operation("add", args[0], args[1])
        if err != nil { log.Fatal(err) }
        log.Printf("result: %f\n", result)
    },
    Args: cobra.ExactArgs(2),        // arqument sayı validasiyası
}
func main() {
    calculator := &cobra.Command{Use: "calculator"}
    calculator.AddCommand(cmdAdd, cmdSub, cmdMul, cmdDiv)
    calculator.Execute()
}
```
Struktur: `$ app [global options] command [command options] [arguments...]` — global
flag-lar `Global*` getter-ləri, komanda flag-ları `cmd.Flags().GetString("name")` ilə.
`MarkFlagRequired("name")` — məcburi flag. Alternativ: urfave/cli.

### 4. Konfiqurasiya — JSON / YAML / INI
**12-factor prinsipi:** config environment-da (fayl YOX) — PaaS/konteyner dünyasında
fayl sistemi hər instance-da fərqlidir.

**JSON (standard):**
```go
type configuration struct {
    Enabled bool    // BÖYÜK hərf məcburi — unexported sahə unmarshal OLUNMAZ
    Path    string
}
file, _ := os.Open("config.json")
defer file.Close()
conf := configuration{}
json.NewDecoder(file).Decode(&conf)
```
Çatışmazlıq: komment YOX, rigidlir.

**YAML (go-gypsy):** komment + oxunaqlı; `yaml.ReadFile` → `config.Get("path")`,
`config.GetBool("enabled")` (typed getterlər). Kitabın tövsiyəsi.

**INI (go-ini/ini):** `ini.Load("conf.ini")` → `config.Section("Section").Key("path").String()`;
Bool oxunuşunda "yoxdur" vs "false" fərqi error ilə ayırd edilir.

### 5. Environment Dəyişənləri
```go
if port = os.Getenv("PORT"); port == "" {
    panic("environment variable PORT has not been set!")
}
http.ListenAndServe(":"+port, nil)
```
- Bütün dəyərlər stringdir → strconv ilə çevir
- **Təhlükəsizlik:** subprocess-lər env-i GÖRÜR; namespace et (`MYAPP_PORT`); build/run
  mərhələlərini ayır (Ch12)

### 6. Server Start/Stop — İnit Daemonlar
**Qayda:** Öz applikasiyanı daemon olaraq YOX — systemd/upstart/init/supervisor/launchd
kimi init daemon ona baxsın (restart, health). NGINX reverse proxy Go üçün adətən lazımsız
— daemon birbaşa servisi idarə etsin. Kubernetes = konteyner səviyyəli "daemon".

### 7. Graceful Shutdown — OS Siqnalları
**Problem:** `http` paketi default olaraq dərhal ölür — açıq connection-lar, yazılmamış
data itir.

**Kitabdan kod nümunəsi:**
```go
server := &http.Server{Addr: ":8080", Handler: handleFunc}
ch := make(chan os.Signal, 1)
signal.Notify(ch, os.Interrupt, os.Kill)   // siqnal kanalı
go func() {
    server.ListenAndServe()                 // bloklayan çağırış — goroutine-də
}()
<-ch                                        // Ctrl-C / kill gözlə
if err := server.Shutdown(nil); err != nil { // mövcud requestləri bitir, yeniləri ALMA
    panic(err)
}
```

**Üstünlüklər:** cari requestlər tamamlanır; TCP port dərhal boşalır (yeni versiyanın
port-ə qoşulması — zero-downtime deploy); in-flight data (upload, DB multi-insert)
bitir. **Çatışmazlıq:** uzun-ömürlü socket-lər (websocket) başqa instansiyaya ötürülə
BİLMİR. Qeyd: Shutdown yalnız HTTP handler-ləri gözləyir — öz goroutine-ləriniz üçün
WaitGroup lazımdır.

### 8. Routing — Çoxlu Handler
```go
http.HandleFunc("/hello", helloHandler)       // query: ?name=...
http.HandleFunc("/goodbye/", goodbyeHandler)   // /goodbye/Buttercup — strings.Split ilə
http.HandleFunc("/", homePageHandler)          // fallback; http.NotFound(res, req)
```
- Spesifikdən spesifiksizə həll olunur
- `req.URL.Query().Get("name")` — query string; `req.URL.Path` — yol
- HTTP metod ayrımı: `req.Method` yoxlaması (GET/POST...)
- İstifadəçi məzmununu SANITIZE et (XSS) — Ch9 templating

### 9. Go 1.22 ServeMux — Metod + Path Dəyişənləri
**Nədir:** 1.22-dən routerdə metod prefiksi + `{name}` dəyişənləri built-in oldu:
```go
mux := http.NewServeMux()
mux.HandleFunc("GET /goodbye/{name}", goodbyeHandler)
name := req.PathValue("name")        // dəyişəni oxu
```
Precedence: longest-match-wins. Əvvəllər üçüncü tərəf mux kitabxanası tələb edən
 hallar (metod ayrımı, path dəyişəni) artıq standarddır. Məhdudiyyət: `{bar}` bir seqment
 (`foo/{bar}/baz` ayrıca yazılmalı).

### 10. Custom Regex Router (nümayiş üçün)
```go
rr.Add("GET /hello", helloHandler)
rr.Add("(GET|HEAD) /goodbye(/?[A-Za-z0-9]*)?", goodbyeHandler)
// ServeHTTP: method+" "+path key; cache map (miss-i cache-ləmə — DoS qorunması);
// compile olunmuş regexp-lər əvvəlcədən saxlanılır (hər requestdə recompile YOX)
```
Trades-off: güclü, amma oxunmaz + test tələb edir; map iterasiyası nondeterministikdir —
qismən uyğunluqlar düzüntü tələb edir.

### 11. Üçüncü Tərəf Routerlər
- **httprouter** — sürət + minimal yaddaş; case-insensitive, path cleanup
- **gorilla/mux** — host, scheme, header matching; deprecation-dan dirilib
- **gin** — zero-allocation router, framework
- **pat / gorilla/pat** — Sinatra-vari `/user/:name` sintaksisi
Standart 1.22 çox halları örtür; sürət/mürəkkəblik ehtiyacında üçüncü tərəf.

## Əsas terminlər
- Flag/Option — CLI parametri (Plan 9 stili: birleşdirmə YOX)
- Cobra — subcommand CLI framework
- 12-factor app — config environment-da saxlayan metodologiya
- Init daemon — servis həyat dövrünü idarə edən sistem (systemd və s.)
- Graceful shutdown — mövcud requestləri bitirərək dayanma
- Signal.Notify — OS siqnallarını kanala yönləndirmə
- ServeMux — standart router; PathValue — {name} dəyişəni (Go 1.22)
- Precedence — route uyğunluq qaydası (longest-match-wins)

## Praktik nətidə

CLI layihə qərarları: (1) sadə alət → flag paketi (yalnız `-name` formasını qəbul et);
(2) subcommand-lar → Cobra (Docker/K8s tərzi); (3) konfiqurasiya: 12-factor env vars
(namespace ilə), YAML kommentli konfiq üçün, JSON standard kitabxana ilə; (4) enum yerinə
slice + slices.Contains validasiyası; (5) production-da server-i init daemon ilə idarə
et + signal.Notify + server.Shutdown ilə zero-downtime; (6) routing: default ServeMux
(1.22 metod+dəyişən dəstəyi), güclü ehtiyac — httprouter/gin.

## Mənbə
Pages: 29-60 (PDF 50-81)
