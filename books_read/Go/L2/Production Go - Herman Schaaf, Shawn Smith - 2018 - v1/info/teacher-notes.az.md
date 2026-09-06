# Production Go — Müəllim qeydləri (Teacher Mode)

📖 Kitab deyir:
Production-ə hazır Go: tooling (gofmt/vet/golint), tam sintaksis incəlikləri (underflow, nil interface), table-driven test + mock interfeyslər, benchmark (benchmem, escape analysis, bitwise), race detector workflow-ları, security headerları (CSRF/HSTS/CSP) və CI (Travis/Makefile).

👨‍🏫 Müəllim qeydi:
Məncə bunu praktikada belə tətbiq etmək daha yaxşı olar:

1. **Kitab 2018-in sonu (Go 1.11 dövrü) — GOPATH əsrlidir.** Bugünkü oxucu üçün ən vacib dəyişikliklər:
   - GOPATH → **go modules** (Go 1.16+ default): `go mod init` + `go get` artıq versioned; `$HOME/go/src` strukturu lazım deyil.
   - gometalinter → **golangci-lint** (industrial standart; ölü layihədir).
   - misspell → golangci-lint daxilində.
   - gofmt+goimports → editor-un LSP-dən (gopls) avtomatik.
   - Goland → bugün ən yetkin; vim-go → vim-go/LSP; GoSublime (Sublime) təkmilləşib.

2. **GOPATH-era "src qovluğu yarat" təlimatı atlanmalıdır** — `mkdir -p $GOPATH/src/...` yox, `go mod init github.com/prodgopher/dinner` ilə başla.

3. **Yarımçıq fəsillər haqqında:** Monitoring (Prometheus/Grafana), Optimization, Deployment, Drone bölməsi — müəlliflərin öz "work in progress" qeydləri var. Bunları kitabdan yox, müasir mənbələrdən öyrən. Kitabın README-də də göstərilir ki, bu Leanpub "living book" idi — son versiya 2018-də dayanıb.

4. **Race detector workflow-ları hələ də QIZIL dəyərindədir.** Bütün nümunələr (Cat Mutex-embed; API WaitGroup) bu gün də eyni. Əlavə edərdim: `-race` yalnız happens-before xətlərini tapa bilən race-ləri yaxalayır — false negative mümkündür; CI-da `go test -race ./...` MÜTLƏQ; production-da Go'nun race detector yalnız 64-bit + OS supportlu platormlarda.

5. **Benchmark dərsləri də aktualdır** (ResetTimer, benchmem, benchcmp), amma: `benchcmp` → `benchstat` (golang.org/x/perf) — statistik cədvəllər + p-values; yeni Go versiyalarda benchmark output formatı sabit qalıb, amma `b.Loop()` (Go 1.24+) — compiler-optimizasiya dərəcələməsi üçün yeni pattern (b.N manual loop → b.Loop() avtomatik, keyfiyyət qorunur).

6. **Modulo/Bitwise dərsi də aktual** — amma müasir compiler-lərdə fərq 2x qədər qalır. Məsələn `n%4` vs `n&3` hələ də asm-də görünür; konstant-fold-cases-lərdə compiler bəzən özü optimallaşdırır.

7. **JSON struct tags** — sənəddə göstərilən kimi hələ aktual. `omitempty`, `json:"-"` da əlavə etmək olardı.

## Ən vacib 5 texnika (bugün üçün)

1. **Race detector workflow** — hər test -race; go run -race + curl production sınaği (kitabın fərqləndiyi yer).
2. **Table-driven test** — bu gün də dominant pattern; anonymous struct dəyişməz.
3. **Mock pattern** — öz interfeys + fixed generator — müəlliflərin randIntGenerator forması Go test dərsi standartıdır.
4. **httptest.Recorder** — HTTP handler testi üçün bu gün də standart.
5. **gofmt/vet CI integration** — Makefile → CI configuration bugün də eyni yanaşma (git hooks + GitHub Actions/GitLab CI müasir ekvivalent).

## Kitabın ən dəyərli hissəsi

**Tooling chapter (s. 113-127)** — race detector-in 2 tam workflow-u + Go Report Card praktikası. Kitabın qalan hissələri müasir versiyaları var (golangci-lint, benchstat və s.), amma race detector bölməsi zaman aşmamışdır.

## Uyğunsuzluq / diqqət qeydləri

⚠️ **Ch 2 (Installing Go):** GOPATH instructions, `mkdir -p $GOPATH/src/...` — modul dövründə tam əvəz olunmalı. gometalinter dead; GoSublime dead (LSP istifadə et).

⚠️ **Ch 7 (Concurrency):** `time.Sleep` nümunələri (goroutine-lərin bitməsini gözləmək üçün) kitabın ÖZ qeyd edir ki, yanlışdır — amma bir neçə nümunədə istifadə olunur. Oxucu üçün: WaitGroup hər yerdə.

⚠️ **Ch 8 (Mocking):** `rand.Seed(time.Now().Unix())` deprecated (Go 1.20+) — `rand.New(rand.NewSource(...))` düzgün. Kitabın öz New() funksiyasında düzgün formadır; test nümunəsində seed qoyulmur.

⚠️ **Ch 9 (Benchmarks):** `go get -u` ilə quraşdırılan `benchcmp` köhnədir — `benchstat` müasir seçim.

⚠️ **Ch 12 (CI):** `.travis.yml` hələ işləyir, amma müasir layihələrin 90% (GitHub Actions/GitLab CI). Drone bugün ölü layihə kimi görünür.

⚠️ **HSTS:** Kitab "go run -race" bənzəri "production sınağı" deyir, amma real production-da HSTS production certificate + test setup tələb edir. Sadə header nümunəsi üçün yaxşı başlanğıcdır.
