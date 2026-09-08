# Chapter 3 — Utilizing HTTP (səh. 47-57)

## Bu fəsil nədən bəhs edir?

`net/http` paketinin praktik reseptləri: GET sorğuları (URL qurumu,
pagination, timeout, LimitReader), chunked encoding ilə POST axını
(os.Pipe), performans middleware-i, server timeout-ları və eyni serverdə
bir neçə API versiyası.

## Əsas fikirlər

### Recipe 15 — GET sorğuları (pagination ilə)
**Tapşırıq:** GitHub API-dən açıq PR-lərin siyahısı (state=open,
per_page=100, page=N).

**Kitabdan kod nümunəsi:**
```go
// URL qurulumu — əl ilə yox, url.Values ilə:
func buildURL(owner, repo string, page int) string {
    query := url.Values{}
    query.Set("state", "open")
    query.Set("base", "master")
    query.Set("per_page", "100")
    query.Set("page", fmt.Sprintf("%d", page))
    owner, repo = url.PathEscape(owner), url.PathEscape(repo)
    const format = "https://api.github.com/repos/%s/%s/pulls?%s"
    return fmt.Sprintf(format, owner, repo, query.Encode())
}

// Parse — anonim lazımi sahələr:
type PR struct {
    Number int
    Title  string
}
func parseResponse(r io.Reader) ([]PR, error) {
    var prs []PR
    if err := json.NewDecoder(r).Decode(&prs); err != nil {
        return nil, err
    }
    return prs, nil
}

// Səhifənin çəkilməsi:
func pageOpenPRs(owner, repo string, page int) ([]PR, error) {
    url := buildURL(owner, repo, page)
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    req, err := http.NewRequestWithContext(ctx, "GET", url, nil) // body yoxdur → nil
    if err != nil {
        return nil, err
    }
    req.Header.Set("Accept", "application/vnd.github.v3+json")
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return nil, err
    } else if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("bad status: %d - %s", resp.StatusCode, resp.Status)
    }
    defer resp.Body.Close()
    r := io.LimitReader(resp.Body, 10*(1<<20))  // maks 10MB oxu
    return parseResponse(r)
}

// Səhifə-səhifə yığılma:
func openPRs(owner, repo string) ([]PR, error) {
    var prs []PR
    for page := 1; true; page++ {
        pagePRs, err := pageOpenPRs(owner, repo, page)
        if err != nil {
            return nil, err
        }
        if len(pagePRs) == 0 {
            break
        }
        prs = append(prs, pagePRs...)
    }
    return prs, nil
}
```

**Prinsiplər:**
- URL-i əllə fmt ilə YOX — `url.Values` + `url.PathEscape` (xüsusi
  simvollar, təhlükəsizlik)
- Şəbəkə etibarsızdır → Context ilə timeout (10s)
- Şəbəkə təhlükəsizdir → `io.LimitReader` (100GB cavab qoruması)
- Retry məntiyi mürəkkəbdir → `hashicorp/go-retryablehttp` nümunəsi

### Recipe 16 — POST axını (chunked encoding)
**Tapşırıq:** metrikleri (kanaldan gəlir) serverə axıt — ölçüsü əvvəlcədən
bilinmir.

```go
type Metric struct {
    Name  string    `json:"name"`
    Host  string    `json:"host"`
    Time  time.Time `json:"time"`
    Value float64   `json:"value"`
}

// Producer — kanalı JSON-a kodlayıb WriteCloser-a yazır:
func producer(w io.WriteCloser) {
    defer w.Close()                     // "data bitdi" siqnalı
    enc := json.NewEncoder(w)
    for m := range collectMetrics() {
        if err := enc.Encode(m); err != nil {
            log.Printf("error: can't encode %#v - %s", m, err)
            return
        }
    }
}

func updateMetrics() error {
    r, w, err := os.Pipe()              // Reader/Writer cütü
    if err != nil {
        return err
    }
    go producer(w)                       // goroutine-də istehsal
    req, err := http.NewRequest("POST", serverURL, r)   // Reader = body!
    if err != nil {
        return err
    }
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return err
    }
    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("bad reply status: %d %s", resp.StatusCode, resp.Status)
    }
    return nil
}
```
- Go avtomatik **HTTP chunked transfer encoding** işlədəcək (ölçü
  əvvəlcədən verilmir)
- Chunked: hər hissənin ölçüsü + datanın özü göndərilir
- İstifadə halı: ölçü bilinməyən/böyük data — yaddaşda tutmursunuz

### Recipe 17 — performans middleware
```go
// addLogging — handler-ı büküb vaxt ölçür:
func addLogging(name string, handler http.Handler) http.Handler {
    wrapper := func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        handler.ServeHTTP(w, r)
        duration := time.Since(start)
        log.Printf("%s took %s", name, duration)
    }
    return http.HandlerFunc(wrapper)
}

// İstifadə:
hdlr := addLogging("query", http.HandlerFunc(queryHandler))
http.Handle("/query", hdlr)
```
- Middleware = http.Handler alır, http.Handler qaytarır; əvvəl/sonra
  funksionallıq əlavə edir (auth, logging, recovery)
- Handler təmiz qalır — məntiq ilə infrastruktur ayrılır
- chi kimi framework-lərdə hazır middleware dəsti var

### Recipe 18 — server timeout-ları
```go
srv := http.Server{
    Addr:         ":8080",
    Handler:      http.DefaultServeMux,
    ReadTimeout:  3 * time.Second,   // oxuma limiti
    WriteTimeout: 2 * time.Second,   // yazma limiti
}
if err := srv.ListenAndServe(); err != nil {
    log.Fatalf("error: %s", err)
}
```
- **Hər şəbəkə əməliyyatına timeout qoyun** — asılıqalmaya və zərərli
  istifadəyə qarşı
- `ReadTimeout`/`WriteTimeout` — ən çox işlədilən ikili
- `http.ListenAndServe` yerinə konfiqurlu `srv.ListenAndServe()`

### Recipe 19 — eyni serverdə API versiyaları
**Tapşırıq:** v1 və v2 API paralel xidmət etməli (köçürmə dövrü).
**Qərar:** `X-API-ver: 2` header-i v2-yə yönləndirir.

```go
func v1Mux() *http.ServeMux {
    mux := http.NewServeMux()
    mux.HandleFunc("/health", v1HealthHandler)
    return mux
}
func v2Mux() *http.ServeMux {
    mux := http.NewServeMux()
    mux.HandleFunc("/_/health", v2HealthHandler)
    return mux
}

func main() {
    v1 := v1Mux()
    v2 := v2Mux()
    versionRouter := func(w http.ResponseWriter, r *http.Request) {
        if r.Header.Get("X-API-ver") == "2" {
            v2.ServeHTTP(w, r)
        } else {
            v1.ServeHTTP(w, r)
        }
    }
    http.HandleFunc("/", versionRouter)   // hamısı bu router-dən keçir
    addr := ":8080"
    log.Printf("info: server ready on %s", addr)
    if err := http.ListenAndServe(addr, nil); err != nil {
        log.Fatalf("error: %s", err)
    }
}
```

**Test:**
```bash
$ curl http://localhost:8080/health           # OK (v1)
$ curl http://localhost:8080/_/health         # 404
$ curl -H 'X-API-ver: 2' http://localhost:8080/_/health  # OK (v2)
$ curl -H 'X-API-ver: 2' http://localhost:8080/health    # 404
```
- ServeMux-lər ayrıdır → yollar ayrıca tərzdə ola bilər (/_/health)
- API dəyişmək ciddi qərardır — istifadəçilərə uzun köçürmə müddəti verin
  (Go-nun compatibility promise uğurunun səbəbi; Python 2→3 = 8 il!)
- Bir çox developer ilk gündən `/v1` prefiksi qoyur

## Final Thoughts-dən

net/http çoxluğu örtür; fancy routing (metodlar üzrə, path dəyişənləri)
lazım olsa chi kimi paketlər — amma xarici asılılıq riski ilə dəyərini
tartın.

## Əsas terminlər
- url.Values/PathEscape — təhlükəsiz URL qurulumu
- Pagination — səhifələnmiş API nəticələri
- io.LimitReader — oxu həddi (DoS qoruması)
- NewRequestWithContext — kontekstli sorğu
- Chunked transfer encoding — ölçüsüz HTTP axını
- os.Pipe — producer/consumer körpüsü
- Middleware — handler sargısı
- ReadTimeout/WriteTimeout — server şəbəkə limitləri
- API versioning — paralel versiya xidməti

## Praktik nəticə
GET üçün: url.Values ilə qur → Context timeout → LimitReader → status
yoxla → defer Close. Axın POST üçün: os.Pipe + goroutine producer +
NewRequest(POST, url, reader) — chunked avtomatik. Serverdə middleware
= Handler→Handler funksiyası; timeout-lar mütləq; API versiyalashdırma
header/mux ayrılığı ilə təmiz həll olunur.

## Mənbə
Pages: 47-57 (PDF 47-57)
