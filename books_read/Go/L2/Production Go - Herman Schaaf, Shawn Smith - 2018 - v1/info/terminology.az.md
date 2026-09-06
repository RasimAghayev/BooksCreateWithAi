# Production Go — Terminologiya lüğəti (AZ)

`English (Azərbaycanca)` formatında — SYSTEM_PROMPT qayda 4 üzrə.

## Tooling
- goimports (format + avtomatik import idarəsi)
- gometalinter / deadcode / ineffassign / misspell
- go vet (printf verb yoxlaması / formatting directive)
- gofmt -s (simplify — məntiq toxunulmadan)
- golint (else-return xəbərdarlığı / exported comment)
- Go Report Card (paket qiymətləndirmə)
- godoc / doc.go / "Package [name] ..."
- Go Guru (:GoImplements / :GoReferrers / :GoCallees / :GoCallers)
- golang-announce [security]

## Sintaksis
- Zero Value / Short Declaration (`:=`)
- Redeclare (:= yeni dəyişən tələb edir; var — YOX)
- Mismatched Types (int32 + int64 → cast)
- Constant Overflow (`128` int8-ə sığmaz)
- **uint Underflow** (0-1 → max value)
- Anonymous Struct (table-driven test cases)
- Struct Tag (`json:"ad"`) — sahələr EXPORTED
- Nil Interface Tələsi ((tip, dəyər) cütü)
- sync.Map (write-once / disjoint keys)
- select + default (kanal doludursa ötür)
- len(ch) (kanal doluluğu)

## Strings / Unicode
- strings.Split / Count / Index / LastIndex / Contains / HasPrefix / HasSuffix
- FieldsFunc (rune funksiyası ilə split)
- EqualFold (case-insensitive müqayisə)
- ASCII / Unicode / Code Point / UTF-8
- Rune (int32 = 1 Unicode simvol)
- Range: index BAYT, dəyər RUNE (UTF-8 fərzi)
- %+q (ASCII-only output), %x/%X (hex), %# x (0x prefiksli)
- Escape Analysis (-gcflags=-m; stack vs heap)
- n % m == n & (m-1) (2-nin qüvvəti modullo qısaldması)
- go tool compile -S (ASM çıxarışı)

## Concurrency
- Goroutine (lightweight thread)
- WaitGroup (Add/Done/Wait; Add sayı dəqiq!)
- Ticker / time.NewTicker / ticker.C
- Poller (fondda API yoxlaması)
- Handler-da fon goroutine
- Race Condition / Data Race
- -race Detector (WARNING: DATA RACE)
- safeMap (Mutex embed + Store/Load/Delete)

## Testing
- Table-driven Test (cases anonymous struct)
- actual != expected (konvensiya)
- t.Error / t.Fatal (davam / dayandır)
- httptest.NewRecorder / http.HandlerFunc adapter
- Mock / Fixed Generator (calledWithN qeydi)
- -cover / -coverprofile / go tool cover -html
- ExampleX / ExampleT_M / // Output:
- Benchmark / b.N / -benchmem (B/op, allocs/op)
- benchcmp (old vs new delta)
- b.ResetTimer / b.SetBytes

## Security
- CSRF (Cross-Site Request Forgery) / per-session token
- nosurf / gorilla/csrf
- HSTS (Strict-Transport-Security) / max-age / includeSubDomains
- headerWrap (wrapper handler — mərkəzi header)
- CSP (Content-Security-Policy)
- bluemonday (HTML sanitizer) / unrolled/secure

## CI / Deploy
- Travis CI (.travis.yml)
- Drone (.drone.yml — Docker əsaslı)
- Jenkins (Jenkinsfile)
- Makefile target-ləri (lint / build / test)
- Prometheus / Grafana (monitoring — kitabda yarımçıq)
