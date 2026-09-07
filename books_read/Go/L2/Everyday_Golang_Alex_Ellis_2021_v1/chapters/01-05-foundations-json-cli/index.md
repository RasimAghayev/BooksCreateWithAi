# Chapters 1-5 — Giriş, ilk proqramlar, cross-compile, JSON, CLI (səh. 4-33)

## Bu fəsillər nədən bəhs edir?

Kitabın məqsədi (production alətlərindən götürülmüş gündəlik pattern-lər),
Go-nun güclü tərəfləri (statik compile, cross-compile, stdlib zənginliyi),
Go quraşdırma + GOPATH strukturu, Go modules və vendoring (dep → modules
keçidi, `go mod init/vendor`, `-mod=vendor`), ilk proqram + VS Code, xarici
asılılıq (hmac sign/validate), flag parsing (`flag.StringVar` + `TrimSpace`
validasiyası), çoxpaketli layihə (main + cmd paketi, böyük hərf = export),
cross-compilation (GOOS/GOARCH matrisi, `-ldflags "-s -w"` strip, CGO_QADAĞASI),
HTTP+JSON kliyent (Unmarshal struct tag-lərlə, http.Client{Timeout},
User-Agent etikası, funksiyaya çıxarma), 5 CLI prinsipi (Go seç, flags→
Cobra, avtomatlaşdırma, paket menecerləri, geri bağlantı).

## Əsas fikirlər

### 1. Go-nun Gündəlik Dəyəri (Ch1)
- **Niyə Go:** statik binar (runtime YOX — JS/Python-dan fərqli), cross-
  compile bir əmrlə, konkurrent model, vahid stil (codebase-lər arası keçid
  asandır), stdlib: crypto/compress/http/net/encoding/json/text/template/os
- **Məşhur Go layihələri:** Docker, Kubernetes, Traefik, Consul/Terraform,
  Caddy, OpenFaaS
- **GOPATH strukturu:** `$HOME/go/` → `pkg/` (module cache), `bin/`
  (go install binariləri — PATH-ə əlavə et), `src/` (köhnə konvensiya)
- **go install/go get:** `go install path@v1.3.0` (binari); `go get`
  modul üçün (1.17-dən binari quraşdırma DEPRECATE); versiyasız modullarda
  pseudo-version `v0.0.0-...-SHAhort`

### 2. Go Modules vs Vendoring (Ch1)
- **Keçid tarixi:** vendor/ + dep (lock faylları) → Go 1.13: **modules** —
  asılılıq kodu artıq repo-ya commit OLUNMUR, module cache-də
- **`go mod init`:** GOPATH daxilində ad avtomatik; kənarda mütləq arqument:
  `go mod init github.com/alexellis/hash-gen`
- **Vendoring hələ mövcuddur:** `go mod vendor` → vendor/ yarat;
  `go build -mod vendor` ilə build — private asılılıqlar + CI konteynerləri
  üçün ən asan yol; review-lər isə modules ilə təmizdir (dəyişən fayl az)

### 3. İlk Proqram + Xarici Asılılıq (Ch2)
```go
// go mod init → asılılıq əlavəsi avtomatik go.mod/go.sum-a düşür
digest := hmac.Sign(input, secret)          // GitHub-dan çəkilir
err := hmac.Validate(input, fmt.Sprintf("sha1=%x", digest), string(secret))
```
- `%x` — hex çap; `go.sum` versiyanı PIN-ləyir (hash ilə)
- **HMAC konsepti:** simmetrik açar; göndərən hash qoyur, qəbul edən
  yenidən hesablayıb müqayisə edir — webhook doğrulamasının əsası

### 4. Flags (Ch2)
```go
flag.StringVar(&inputVar, "message", "", "message to create a digest from")
flag.Parse()
if len(strings.TrimSpace(secretVar)) == 0 {   // boş + whitespace hər ikisi
    panic("--secret is required")
}
```
- `%q` — dırnaq içində çap (gizli whitespace-i ifşa edir)
- **Strategiya:** hər yeni CLI standart `flag` paketi ilə BAŞLAYIR; həddə
  çatanda Cobra-ya keç (Docker/K8s/OpenFaaS bunu edir)

### 5. Çoxpaketli Layihə (Ch2)
```
multiple-packages/
├── main.go        ← package main (giriş)
└── cmd/ls.go      ← package cmd (kitabxana)
    func ExecuteLs(path string) (string, error)   // BÖYÜK hərf = export
```
- Testlər EYNİ paketdə yazılır (cmd/ üçün cmd/ daxilində)
- `go install` → binari GOPATH/bin-ə düşür

### 6. Cross-Compile (Ch3)
```bash
GOOS=linux   go build -o first-go
GOOS=darwin  go build -o first-go-darwin
GOOS=windows go build -o first-go.exe
GOARCH=arm   GOOS=linux go build -o first-go-arm    # Raspberry Pi 32-bit
GOARCH=arm64 GOOS=linux go build -o first-go-arm64  # ARM server (Graviton)
```
- **Azaltma:** `go build -ldflags "-s -w"` — debug simvolları sil (2.0M →
  1.4M); `file` komandası ilə format yoxla
- **CGO_QADAĞASI:** `CGO_ENABLED=0` — C kitabxanalarına bağlanmayıbsa,
  portabiliti üçün SÖNDÜR (SQLite, bəzi crypto paketləri CGO istəyir —
  cross-compile-i çətinləşdirir)

### 7. HTTP + JSON (Ch4)
```go
type people struct {
    Number int      `json:"number"`
    Person []person `json:"people"`
}
spaceClient := http.Client{Timeout: time.Second * 2}   // timeout MÜTLƏQ
req, _ := http.NewRequest(http.MethodGet, url, nil)
req.Header.Set("User-Agent", "spacecount-tutorial")    // vətandaşlıq etikası
res, err := spaceClient.Do(req)
defer res.Body.Close()
body, _ := io.ReadAll(res.Body)
err = json.Unmarshal(body, &p)                          // & — UNUTMA!
```
- **Unmarshal qaydaları:** struct sahəsi BÖYÜK hərflə (export); `json:"name"`
  tag-i JSON açarı ilə birləşdirir; `[]byte` qəbul edir; `&p` — kopyaya yox,
  originala yaz
- **Refaktorinq dərsi:** funksiyaya çıxar (getAstros) → log.Fatal MÜXTƏLİF
  çağıranlara uyğun deyil; inline if: `if err := ...; err != nil {}` — err
  scope-u if-də qalır

### 8. 5 CLI Prinsipi (Ch5)
1. **Go seç:** statik binar, sürətli başlanğıc (Node.js 1.5-3s lag), vahid
   stil; npm ilə də paylana bilər (binari yükləyən paket)
2. **flags → Cobra:** sub-komandalara böyüyəndə; verb-noun sintaksisi:
   `-deploy` → `faas-cli deploy --image=...` — UX fluidliyi
3. **Avtomatlaşdır:** GitHub Actions + releases (hər platformaya binari)
4. **Paket menecerləri:** brew, apt, npm — istifadəçiyə tanış quraşdırma
5. **Geri bağlantı topla:** issue-lar, istifadəçi təcrübəsi

## Əsas terminlər
- Statically-linked binary — runtime tələbi olmayan, tək fayl proqram
- GOPATH — Go iş kataloqu (pkg/bin/src)
- Pseudo-version — tag-siz modul üçün `v0.0.0-tarix-SHA` versiyası
- Vendoring — asılılıqların vendor/ kataloqunda saxlanması
- Marshal/Unmarshal — seriyalaşdırma/əks-seriyalaşdırma
- Struct tag — `json:"field"` metadata annotasiyası
- HMAC — simmetrik açarla mesaj doğrulaması
- CGO — Go↔C körpüsü; cross-compile düşməni
- Verb-noun syntax — `kubectl get pods` üslubu (Cobra)
- User-Agent — HTTP sorğunun identifikasiya başlığı (etiket vacibdir)

## Praktik nəticə

1. **Yeni layihə ritualı:** `go mod init` → standart flag → gərəksə
   Cobra; vendoring yalnız private deps/CI ehtiyacında.
2. **Hər HTTP sorğusu:** client Timeout + User-Agent + defer Body.Close()
   — üçlü vətənborluq.
3. **Cross-compile matrisi:** ən çox 4 hədəf (linux/amd64, darwin, windows,
   linux/arm) — GitHub Actions-da avtomatik.
4. **Binarini strip et:** `-ldflags "-s -w"` — ~30% qənaət.
5. **CGO söndür:** `CGO_ENABLED=0` cross-compile planında varsa.

## Mənbə
Pages: 4-33 (PDF 5-34)
