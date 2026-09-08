# Chapter 12 — Distributing Your Apps (Tətbiqlərin Paylanması)

## Bu chapter nədən bəhs edir?
Tətbiqi paylama zəncirinə: Go modules (asılılıq idarəetməsi, SemVer, MVS, private repo-lption), module workspaces (go.work), CI (GitHub Actions: test + cache + Staticcheck) və release avtomatlaşdırması (GoReleaser — cross-compile, arxivlər, Docker image-lər).

## Əsas fikirlər

### 1. Go Modules — 3 problemin həlli
Go 1.11-də gəldi; "works on my machine" sindromunu öldürür.

| Problem | Əvvəl | Sonra |
|---|---|---|
| Reliable versioning | versiya qeyri-müəyyən | SemVer (v<major>.<minor>.<patch>) |
| Reproducible builds | repo yoxa çıxsa build sınırdı | module proxy + vendor |
| Dependency bloat | nested asılılıqlar şişirdi | minimal lazımi dəst hesablanır |

**Kitabxana analoqiyası:** repo = kitabxana bölməsi, module = kitab seriyası, package = kitab, source file = fəsil. Versiya = git tag-lər.
**SemVer:** major = qırıcı API dəyişikliyi, minor = geri-uyğun yeni funksiya, patch = bugfix. **v2+ modul yolu major daşıyır:** `github.com/alexrios/timer/v2`.

### 2. Module iş axını
```bash
go mod init github.com/yourusername/mybestappever   # go.mod yaradır
# .go fayllarında import yaz → go build / go test avtomatik əlavə edir:
require github.com/alexrios/timer/v2 v2.0.0
go get foo@v1.2.3                                   # konkret versiya
go get github.com/alexrios/timer@latest             # ən yeniyə
go list -m all                                      # modul + asılılıqlar
go mod tidy                                         # təmizlə: lazımsızı at, çatmağı əlavə et
go install                                          # $GOPATH/bin-ə quraşdır
go install github.com/alexrios/endpoints@v0.5.0     # versiyalı quraşdırma (Go 1.16+)
go install github.com/alexrios/endpoints@latest
```
**go.sum:** modul versiyalarının checksum-ları — build-lər arası tutarlılıq zəmanəti.

### 3. MVS (Minimal Version Selection)
Go-un asılılıq versiya alqoritmi:
1. **Başlanğıc:** öz go.mod-dakı birbaşa asılılıq versiyaları.
2. **Yayılma:** hər asılılığın go.mod-u rekursiv oxunur.
3. **Seçim:** hər modul üçün go.mod-larda tələb olunan ƏN YÜKSƏK versiya götürülür.
4. **Minimalizm:** Bundan yuxarı heç nə seçilmir — lazımsız upgrade və qırıcı dəyişiklik riski yoxdur.

### 4. Private asılılıqlar
```ini
# ~/.gitconfig:
[url "ssh://git@github.com/"]
    insteadOf = https://github.com/          # HTTPS yerinə SSH
```
```bash
go env -w GOPRIVATE="github.com/<org>/<project>"   # tək repo
go env -w GOPRIVATE="github.com/<org>/*"           # bütün org (wildcard)
go env GOPRIVATE                                    # yoxla
ssh-add ~/.ssh/your-ssh-key                         # agent-ə açar
```
**Sub-kod izahı:** GOPRIVATE siyahısındakı modullar public proxy/mirror-dan çəkilmir — birbaşa mənbədən (SSH ilə) gəlir; private repo üçün məcburi tənzimləmə.

### 5. Module workspaces (Go 1.18+) — go.work
Bir layihədə çoxlu modulla iş: go.work faylı modulları birləşdirir, kompilyator onları "peer" kimi görür — bir moduldakı dəyişiklik dərhal digərlərinə tətbiq olunur.
```bash
go work init ./myproject
go work use ./moduleA
go work use ./moduleB
go work sync        # go.mod dəyişəndə go.work-ü sinxronlaşdır
```
```go
// go.work:
go 1.21
use (
    ./path/to/module-a
    ./path/to/module-b
)
```
**`go work sync` niyə vacib:** sync-siz qalarsa — go.mod/go.work uyğunsuzluğu, keşdəki köhnə versiyanın istifadəsi, sınan build-lər, komanda üzvləri arasında divergent dependency vəziyyəti.

### 6. CI — GitHub Actions
Hər push-da test + build avtomatik; xəta erkən tutulur.

**Baza workflow:**
```yaml
name: Go CI on Commit
on: [push]
jobs:
  test-and-dependencies:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v3
        with:
          go-version: '^1.21'
      - run: go mod download
      - run: go test -v ./...
```

**Modul cache (sürət):**
```yaml
      - name: Cache Go modules
        uses: actions/cache@v3
        with:
          path: ~/go/pkg/mod                    # modul anbarı
          key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
          restore-keys: ${{ runner.os }}-go-    # dəqiq tapılmasa fallback
```
**Sub-kod izahı:** cache açarı go.sum hash-indən düzəlir — asılılıq dəyişsə yeni cache yaradılır; `~/go/pkg/mod` (module cache) saxlanılır, download-lər təkrar olmur.

**Staticcheck (static analysis):**
```yaml
      - uses: dominikh/staticcheck-action@v1.2.0
```
Linting-dən dərinə: bug-lar, ineffektiv pattern-lər, style problemləri. Action versiyası pin edilir — CI davranışı stabil qalır.

### 7. Release — GoReleaser
Cross-compilation, arxivlər, Docker image, GitHub release — bir konfiqurasiyadan.

**.goreleaser.yml:**
```yaml
builds:
- ldflags:
  - -s -w                              # symbol/debug info soyulur
  - -extldflags "-static"              # tam statik link
  env:
  - CGO_ENABLED=0                      # saf Go — hər platformada
  goos: [linux, windows, darwin]
  goarch: [amd64]
  mod_timestamp: '{{ .CommitTimestamp }}'
archives:
- name_template: "{{ .ProjectName }}_{{ .Version }}_{{ .Os }}_{{ .Arch }}"
  wrap_in_directory: true
  format: binary
  format_overrides:
  - goos: windows
    format: zip                        # Windows üçün zip
dockers:
- image_templates:
  - "ghcr.io/alexrios/endpoints:{{ .Tag }}"
  - "ghcr.io/alexrios/endpoints:v{{ .Major }}"
  - "ghcr.io/alexrios/endpoints:v{{ .Major }}.{{ .Minor }}"
  - "ghcr.io/alexrios/endpoints:latest"
```

**Tag-triggered release workflow:**
```yaml
name: GoReleaser
on:
  push:
    tags: ['*']
jobs:
  goreleaser:
    runs-on: ubuntu-latest
    permissions:
      packages: write
      contents: write
    steps:
      - uses: actions/checkout@v2
        with: { fetch-depth: 0 }          # tam tarix — changelog üçün
      - uses: actions/setup-go@v2
        with: { go-version: 1.21 }
      - uses: goreleaser/goreleaser-action@v2
        with:
          version: latest
          args: release --rm-dist --clean
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```
**Sub-kod izahı:** tag push olunca workflow işə düşür; `fetch-depth: 0` bütün commit tarixini verir (release notes üçün); permissions GHCR push-a icazə verir; ldflags `-s -w` binary həcmini kiçildir.

## Əsas terminlər
- Go Modules / go.mod / go.sum (checksum)
- SemVer (Semantic Versioning) — major/minor/patch
- Semantic import versioning — /v2 modul yolu
- MVS (Minimal Version Selection)
- go mod init / go get / go list -m all / go mod tidy
- go install module@version — versiyalı quraşdırma
- GOPRIVATE / .gitconfig insteadOf (HTTPS→SSH)
- go work init / use / sync — module workspaces (Go 1.18+)
- CI (Continuous Integration) — push-da avtomatik test/build
- GitHub Actions: actions/cache, hashFiles(go.sum), ~/go/pkg/mod
- Staticcheck / dominikh/staticcheck-action — static analysis
- GoReleaser / goreleaser-action
- ldflags: -s -w (strip) / -extldflags "-static" / CGO_ENABLED=0
- goos/goarch — cross-compilation hədəfləri
- GHCR (GitHub Container Registry) image tag-ləri
- fetch-depth: 0 — tam git tarixi

## Praktik nəticə
1. go.mod-u əl ilə redaktə etmək əvəzinə go get / go mod tidy işlət — alətlər faylı düzgün saxlayır.
2. go.sum-u commit et — checksum-lar reproducible build-in zəmanətidir.
3. Private repo üçün: GOPRIVATE + SSH insteadOf — əks halda proxy 404/410 qaytarır.
4. Çoxmodullu layihədə hər go.mod dəyişikliyindən sonra go work sync — uyğunsuzluq build qəzası gətirir.
5. CI-də module cache açarı kimi go.sum hash-i işlət; Staticcheck action versiyasını pin et.
6. Release: CGO_ENABLED=0 + static ldflags → tək konfiqdən bütün OS/arch binary-ləri; tag push = release.

## Mənbə
Pages: 247-263 (PDF səh. 268-285)
