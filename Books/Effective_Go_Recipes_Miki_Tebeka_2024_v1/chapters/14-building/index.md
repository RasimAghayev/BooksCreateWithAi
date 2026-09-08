# Chapter 14 — Building Applications (səh. 218-233)

## Bu fəsil nədən bəhs edir?

go build-in imkanları: embed ilə asset-lərin binary-yə daxil edilməsi,
ldflags ilə versiya inyeksiyası, CGO_ENABLED=0 ilə statik build, build
tags ilə şərti kompilyasiya, goreleaser ilə cross-platform və go:generate
ilə kod generasiyası.

## Əsas fikirlər

### Recipe 72 — embed (asset-lər binary-də)
**Problem:** SQL kopyalama xətası — SQL Go kodunda yox, .sql faylında
olsun (syntax highlight + SQL alətləri).

```go
// user.sql:
// SELECT id, name, email
// FROM users
// WHERE id = @id;

import (
    "database/sql"
    _ "embed"
)

//go:embed user.sql      // direktiv — dəyişəndən DƏRHAL əvvəl
var userSQL string      // build vaxtı dolur!

func UserByID(db *sql.DB, id string) (User, error) {
    row := db.QueryRow(userSQL, sql.Named("id", id))
    var u User
    err := row.Scan(&u.ID, &u.Name, &u.Email)
    return u, err
}
```
- Go 1.16+; fayl tapılmasa build düşür
- Qovluq embed: `//go:embed static` + `var staticFS embed.FS` — FS kimi
  açılır; http.FileServer ilə birbaşa xidmət

### Recipe 73 — ldflags ilə versiya
```go
var version = "<unknown>"        // default

var showVersion bool
flag.BoolVar(&showVersion, "version", false, "show version and exit")
flag.Parse()
if showVersion {
    fmt.Printf("agent version %s\n", version)
}
```

**build.go utiliti (_scripts qovluğunda):**
```go
func run(args ...string) (string, error) {
    cmd := exec.Command(args[0], args[1:]...)
    out, err := cmd.CombinedOutput()
    if err != nil {
        return "", err
    }
    return strings.TrimSpace(string(out)), nil
}

// Tag YOXSA commit:
version, err := run("git", "tag", "--points-at", "HEAD")
if version == "" {
    version, err = run("git", "rev-parse", "--short", "HEAD")
}

// İNYEKSİYA:
ldflags := fmt.Sprintf("-ldflags=-X main.version=%s", version)
run("go", "build", ldflags, "-o", "agent")
```
- `-X importpath.name=value` → compile vaxtı dəyişən dəyəri
- _scripts prefiksi → go build iqnor edir; + `//go:build ignore`
  təhlükəsizlik üçün
- Alternativlər (sed/awk, version.go overwrite) daha zəifdir
- Semantik versiya: v1.2.3 (patch/minor/major — major dəyişməzdən əvvəl
  düşün)

### Recipe 74 — statik build (CGO_ENABLED=0)
**Bug:** Alpine konteynerdə `./agent: not found` — fayl VAR amma işləmir!

**Diaqnostika:**
```bash
$ file agent
# dynamically linked, interpreter /lib64/ld-linux-x86-64.so.2
$ ldd agent
# libc.so.6 => ...
```
- **Səbəb:** net/http DNS resolver-i cgo ilə getaddrinfo çağırır →
  dinamik link; Alpine isə libc YOX, musl istifadə edir
- Pure Go resolver default (bloklanan DNS = goroutine; bloklanan C çağırış
  = OS thread!) — amma şərtlər varsa cgo resolver keçir

**Fix:**
```bash
$ CGO_ENABLED=0 go build -o agent agent.go
$ ldd agent
not a dynamic executable          # ✓
```
```dockerfile
ENV CGO_ENABLED=0
```
- Trade-off: bəzi DNS halları fərqli nəticə verə bilər — qərar sizindir

### Recipe 75 — build tags
**Tapşırıq:** default build profil endpoint-i OLMASIN (təhlükəsizlik);
bəzi serverlərdə OLSUN.

```go
// prof.go — ayrı fayl:
//go:build prof

package main

import (
    _ "net/http/pprof"   // blank import: /debug/pprof qeydiyyatı
)
```
```bash
go build                 # profil-siz
go build -tags prof      # profilli
```
- Platform tag-ləri: `//go:build windows`, `//go:build linux`;
  GOOS/GOARCH hazır taglərdir
- Fayl suffiksi alternativi: `_windows.go`, `_arm64.go`
- İstifadə: profiling, metrics, client-a görə funksionallıq

### Recipe 76 — goreleaser (cross-platform)
```bash
$ go install github.com/goreleaser/goreleaser@latest
```
```yaml
# .goreleaser.yaml
project_name: agent
builds:
- env:
  - CGO_ENABLED=0
  targets:
  - linux_amd64
  - darwin_arm64
  - darwin_amd64
  - windows_amd64
```
```bash
$ GORELEASER_CURRENT_TAG=v1.2.3 goreleaser build --rm-dist --skip-validate
# dist/agent_{linux_amd64,darwin_amd64,darwin_arm64,windows_amd64}/...
$ ./dist/agent_linux_amd64_v1/agent -version
agent version 1.2.3
```
- Versiyanı özü inyeksiya edir (ldflags); tag default — env ilə override
- Daha çox: CI/CD inteqrasiya, GitHub release upload, build-əvvəli əmrlər
- **Yadda saxla:** GOOS/GOARCH üçün build yetmir — TEST də lazımdır
- `go tool dist list` — bütün etibarlı kombinasiyalar

### Recipe 77 — go:generate (kod generasiyası)
**Tapşırıq:** icazəli IP siyahısı faylı → kod; kənar IP → 401 + WARNING.

**Generator (_scripts/gen_ips.go):**
```go
var header = `
package main

// IPAllowed returns true if connections from ip are allowed
func IPAllowed(ip string) bool {
    return allowedIPs[ip]
}

var allowedIPs = map[string]bool{`

in, err := os.Open(os.Args[1])
defer in.Close()
out, err := os.Create(os.Args[2])
defer out.Close()

fmt.Fprintln(out, header)
s := bufio.NewScanner(in)
for s.Scan() {
    fmt.Fprintf(out, "%q: true,\n", s.Text())   // "127.0.0.1": true,
}
if err := s.Err(); err != nil {
    log.Fatalf("error: %s", err)
}
fmt.Fprintln(out, "}")
```

**İstifadə (httpd.go daxilində direktivlər):**
```go
//go:generate go run _scripts/gen_ips.go _scripts/allowed_ips.txt ips.go
//go:generate go fmt ips.go

func requestIP(r *http.Request) string {
    fields := strings.Split(r.RemoteAddr, ":")
    if len(fields) != 2 {
        return ""
    }
    return fields[0]
}

// Handler-də:
if ip := requestIP(r); !IPAllowed(ip) {
    log.Printf("WARNING: (sec) access from disallowed IP: %s", ip)
    http.Error(w, "Unknown IP", http.StatusUnauthorized)
    return
}
```
```bash
$ go generate ./...   # bütün direktivləri icra et
```

**Üstünlüklər (Raymond-un Rule of Generation):**
1. Serializasiya YOX — data Go-da hazır
2. Tək fayl deploy — əlavə fayl yoxdur
3. Xətalar TEZ yaxalanır — yararsız Go = compile xətası (embed+runtime
   parse-dən daha erkən)

**Debat:** generasiya olunmuş fayl VCS-də saxlanılsın? Müəllif: BƏLİ
(git tag yenilənməsə də build işləyir). Minus: 2 addımlı build →
Makefile/avtomatlaşdırma şərt.

## Final Thoughts-dən

Tək fayl deploy = dəyərli üstünlük; çoxlu fayl = yer/waxt/axtarış
problemləri. Mümkünsə tək executable-da qalın.

## Əsas terminlər
- //go:embed — binary-yə fayl daxil etmə direktivi
- embed.FS — yaddaş fayl sistemi (qovluq embed)
- ldflags -X — compile-vaxtı dəyişən inyeksiyası
- _scripts + //go:build ignore — build alətlərinin gizlədilməsi
- CGO_ENABLED=0 — statik build məcburiyyəti
- musl vs libc — Alpine'in libc əvəzi
- Build tag (//go:build) — şərti fayl kompilyasiyası
- go install — alət quraşdırma ($GOPATH/bin)
- goreleaser — release avtomatlaşdırması
- GOOS/GOARCH — platform dəyişənləri
- go tool dist list — valid kombinasiyalar
- //go:generate — kod generasiya direktivi
- Rule of Generation — "proqram yazan proqram yaz"

## Praktik nəticə
Assetlər üçün embed (SQL, statik web); versiya üçün -ldflags -X + git
tag/commit; Alpine/musl və ya minimal şəkillər üçün CGO_ENABLED=0 (DNS
trade-off-unu bilin); şərti kod üçün build tags (prof, metrics,
müştəri funksiyaları); çoxplatforma üçün goreleaser + HEDEF kombinasiyalarda
test; datadan kod üçün go:generate + generator _scripts-də + VCS-də
generasiya nəticəsi. Hər halda tək-binary strategiyası üstünlük təşkil edir.

## Mənbə
Pages: 218-233 (PDF 218-233)
