# Chapter 18 — OpenFaaS və Go (səh. 126-157)

## Bu fəsil nədən bəhs edir?

Serverless konsepsiyası (funksiya = ixtisaslaşmış mikro servis: kiçik, tək
məsuliyyət, avtomatik miqyas/monitorinq, şablonla qurulur); OpenFaaS tarixi
(2016, Dockercon; AWS Lambda Go dəstəyi yox idi); cloud functions-dan fərqi
(lock-in YOX, öz limitlərin, white-label, extend); Cons (tight cloud
integration yoxdursa Lambda daha uyğun; K8s bilmirsən — faasd); 3 işləmə
platforması (OpenFaaS/K8s, faasd tək-host, Google Cloud Run); 3 Go şablonu
(go legacy, golang-http req/resp struct, golang-middleware standart
http.HandlerFunc); faas-cli axını (new → build → push → deploy = up);
bcrypt-fn funksiyası (GenerateFromPassword, X-Duration-Seconds başlığı);
handler unit-testləri (httptest; build zamanı DA icra olunur — 90%
coverage); lokal iterasiya (go build; docker run; --shrinkwrap); çoxfayllı
funksiya (eyni package function + sub-package types + go mod edit -replace);
secrets (faas-cli secret create → /var/openfaas/secrets/NAME faylı) vs env
(konfidensial OLMAYAN); DB bağlantısı init()-də (reqreslər arası AÇIQ
qalır); todo-fn REST routing (POST=create, GET + /todo/{id}=get, GET=list).

## Əsas fikirlər

### 1. Funksiya = İxtisaslaşmış Mikro Servis
- **Mikro servisdən fərqi:** daha kiçik; tək use-case; avtomatik scale +
  monitorinq; HTTP server boilerplate-i ŞABLONDA gizli
- **Eyni işləri edir:** DB oxu/yazı, səhifə servisi, batch, data emalı
- **Dəyər:** bütün servisləri BİR yolla idarə → təkrar azalır → sürətli ship

### 2. OpenFaaS vs Cloud Functions
| Ölçü | OpenFaaS | AWS Lambda |
|---|---|---|
| Lock-in | YOX — public/private/on-prem | AWS-ə bağlı |
| Limitlər | ÖZÜN təyin edirsən | Vendor təyin edir |
| White-label | Waylay, LivePerson kimi qutu kimi satılır | Mümkün deyil |
| Extend/fork | Kod açıqdır | Qapalı |
| Cloud-managed inteqrasiya | Zəif | Çox güclü (tight) |

- **3 platforma:** OpenFaaS@K8s (managed control-plane + spot instanslar),
  faasd (tək host, aşağı resurs), Cloud Run (容器 only, managed billing)
- Go şablonları: `golang-middleware` tövsiyə — adi `http.HandlerFunc`
  (kitabda da əsas bu)

### 3. Funksiya Yaratma Axını
```bash
faas-cli template store pull golang-middleware
faas-cli new --lang golang-middleware bcrypt-fn   # → bcrypt-fn.yml + handler.go
export OPENFAAS_PREFIX=docker.io/alexellis2       # image prefix
faas-cli up        # = build + push + deploy
# build: DOCKER_BUILDKIT=1 — paralel + smart cache
```
- stack.yml: functions.bcrypt-fn.{lang, handler, image}
- Çağırış: `curl --data-binary "functions" http://.../function/bcrypt-fn`
  → `$2a$10$...` + başlıqlar: X-Duration-Seconds (FUNKSIYA vaxtı),
  X-Call-Id (log korrelyasiyası), X-Start-Time

### 4. Handler + Test (httptest)
```go
func Handle(w http.ResponseWriter, r *http.Request) {
    body, _ := io.ReadAll(r.Body)
    res, err := bcrypt.GenerateFromPassword(body, bcrypt.DefaultCost)
    // 500 VƏYA 200 + hash
}
func Test_Handler_WithKnownInput(t *testing.T) {
    w := httptest.NewRecorder()
    r := httptest.NewRequest("POST", "/", bytes.NewBufferString("functions"))
    Handle(w, r)
    bcrypt.CompareHashAndPassword(w.Body.Bytes(), []byte("functions"))  // doğrula!
}
```
- **Build testləri də işə salır:** faas-cli build → RUN go test -cover
  (90% coverage logda görünür)

### 5. Lokal İterasiya (3 səviyyə)
1. `go test` + `go build` (handler qovluğunda) — modul/sintaksis xətaları
2. `docker run -p 8080:8080 image` — konteyner local, secrets üçün -v
3. `faas-cli build --shrinkwrap` → ./build/bcrypt-fn → go mod init +
   build → ./handler (port 8082) — SÜRƏTLİ redaktə dövrü

### 6. Çoxfayllı + Sub-package
- **Eyni package (function):** handler.go + types.go hər ikisi `package
  function` → qarşılıqlı giriş; testlər də
- **Sub-package:** types/ qovluğu → import "handler/function/types";
  **go.mod hiyləsi:** `go mod edit -replace=handler/function=./` — local
  yola yönləndir

### 7. Secrets vs Environment
- **Qayda:** konfidansial (db parol/host) → SECRET; qeyri-konfidansial
  (debug, output format) → env
```yaml
secrets:
- db-user
- db-pass
environment:
  db_host: "..."
  db_ssl_mode: "require"
```
```go
func readSecret(name string) (string, error) {
    data, _ := os.ReadFile("/var/openfaas/secrets/" + name)
    return strings.TrimSpace(string(data)), nil
}
```

### 8. DB Bağlantısı — init() İçində
```go
var dbConn *sql.DB
func init() {
    dbUser, _ = readSecret("db-user")
    dbPass, _ = readSecret("db-pass")
    host = os.Getenv("db_host")         // + LookupEnv db_port (Atoi yoxla!)
    dbConn, err = connect(...)           // BAĞLANTI BİR DƏFƏ — sorqular ARASI AÇIQ
    if err != nil { log.Fatalf(...) }    // xəta → konteyner loglarına bax
}
```
- **REST routing (HTTP-dən):**
```go
if r.Method == "POST" { action = "create"; json.Unmarshal(body, &todo) }
else if r.Method == "GET" {
    if strings.HasPrefix(r.URL.Path, "/todo/") {
        id, _ := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/todo/"))
        action = "get"
    } else { action = "list" }
}
```
- POST → 201 Created; GET /todo/{id} → JSON VƏ YA 404; GET → siyahı JSON

## Əsas terminlər
- Serverless/FaaS — funksiya = kiçik idarə olunan mikro servis
- Template — Dockerfile + boilerplate gizlədən skeleton (template store)
- faas-cli up — build + push + deploy birləşdirici
- OPENFAAS_PREFIX — funksiya image-lərinin registry prefiksi
- X-Duration-Seconds / X-Call-Id — çağırış ölçüsü / izləmə ID
- shrinkwrap — build-in Docker-sız variantı (lokal sürətli iterasiya)
- Secret — /var/openfaas/secrets/ faylı (konfidansial data)
- init() funksiyası — handler-dən ƏVVƏL: konfiq + DB bağlantısı
- White-label — OpenFaaS-in başqa platforma daxilində qutu kimi satılması
- faasd — Kubernetes-sız tək-host OpenFaaS

## Praktik nəticə

1. **Funksiyanı adi HTTP handler kimi yaz:** golang-middleware şablonu =
   öyrənilən biliklər birbaşa keçir (mux, httptest, RED metrikaları).
2. **Testlər CI-ə ƏMƏL EDƏCƏK:** build prosesində icra olunur — test yazmaq
   məcburiyyəti yox, amma faydası dərhal görünür.
3. **Secret fayldır, env YOX:** readSecret + TrimSpace; mount avtomatik.
4. **DB bağlantısını init()-də aç:** hər sorquda YENİDƏN qoşulma — anti-
   pattern; xətada log.Fatalf → konteyner restart edir.
5. **Şablon statik yoxlanmır:** testlə; --shrinkwrap sürətli debug üçün.
6. **Multi-env config:** eyni funksiya, fərqli YAML (dev/prod) — env +
   secrets dəyişir, kod YOX.

## Mənbə
Pages: 126-157 (PDF 127-158)
