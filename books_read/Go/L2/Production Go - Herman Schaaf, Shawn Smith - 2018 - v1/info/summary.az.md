# Production Go — Xülasə (Azərbaycan dilində)

**Müəlliflər:** Herman Schaaf, Shawn Smith (Go Report Card müəllifləri)
**Nəşriyyat:** Leanpub, 2018-12-20 · **Səhifə:** 141
**Struktur:** 14 bölmə — basics → production practices (CI/security/monitoring)

## Kitabın ümumi məqsədi

Production-a hazır Go servisi yazmaq: yalnız sintaksis yox — gofmt/golint/vet idiomları, table-driven test + mock patternləri, benchmark dəqiqliyi (escape analysis, modulo-/bitwise), race detector ilə real həllər, security headerları və CI qurulması. Auditoriya: CS əsaslarını bilən, ilk dəfə production üçün Go yazan mühəndislər.

## Chapter-by-chapter xülasə

### Introduction
PHP→Go real story: 10x az cavab vaxtı, az server xərci, "works → just works". Beginner YOX — CS bilən ilk dəfə Go yazanlara.

### Getting Started (s. 1-6)
Binary quraşdırma + `$HOME/go` default GOPATH (Go 1.8+). Editor inteqrasiyaları: GoLand/Sublime/vim — hamısında **goimports on save**. gometalinter (deadcode, ineffassign, misspell). **go vet mütləqdir**: printf verb yanlış tipi, Println içində directive — hər ikisi compile OLUR amma vet yaxalayır.

### Basics (s. 7-51)
Tam sintaksis: 3 var forma + redeclare fərqi (:= yeni tələb edir). Tiplər: implicit conversion YOX (int32+int64 xəta); constant overflow-i compile saxlayır; **uint underflow** (0-1=4294967295!). Operatorlar (5 kateqoriya, spec-dən). if-scoped declaration, şərtsiz switch. Slice (len/cap/append/copy-capacity tələsi/sort.Slice-comparator), map (nil→panic, sync.Map, sıra qeyri-müəyyən). Pointer (aritmetika YOX). Goroutine+WaitGroup (**Add sayı dəqiq** — azı erkən qayıdış, çoxu deadlock). Kanal (buffered, select+default). Interfeys (implicit) + empty interface + type assertion + **nil interface tələsi** (i = nilPointer → i != nil çünki tip doludur). Error handling (fmt.Errorf, ParseBool). stdin (bufio.Scanner) + fayl yazma (os.Create/ioutil.WriteFile).

### Style & Error Handling (s. 46-51)
gofmt istisnasız (exceptions YOX, hamı eyni — success amili!). gofmt -s (simplify). Qısa dəyişən adları — CodeReviewComments ("Prefer c to lineCount"): <10 sətir span → 1 simvol; funksiya <15 sətir. golint (else-return at, exported comment). Error: `if err != nil` təkrarları idiomatikdir — strictlik xətanın hardeyni gizlətmir. `fmt.Errorf` ilə spesifikləşdir; **kiçik hərflə başla** (log axını).

### Strings (s. 52-62)
`+` çox az string üçün; qarışıq tip → Sprintf. strings paketi: Split/Count (non-overlapping!)/Index/Contains/HasPrefix/HasSuffix/FieldsFunc/EqualFold. Palindrome nümunəsi — FieldsFunc + funksiya dəyər kimi + for-un 3 komponenti (while əvəzi).

### Supporting Unicode (s. 62-67)
ASCII→Unicode→UTF-8 tarixi. **Strings are read-only byte slices** — encoding DAŞIMIR; string(b) fərzsizdir. Printf verbləri: %s/%q/%+q(ASCII-only)/%x/% x/%# x; flag-lər + - # 0 ' '. Genişlik RUNE ilə. Range: index BAYT (buna görə 你-dən sonra 6), dəyər RUNE — UTF-8 yeganə yerdi Go fərzi edir.

### Concurrency (s. 68-81)
Goroutine sırası qeyri-deterministik. WaitGroup düzgün idarə (Sleep YOX). Kanal: 3 send / 2 receive → deadlock. **Handler-da fon goroutine** — handler qayıdır, iş DAVAM edir. **Ticker poller**: NewTicker + select + defer Stop; sonsuz loop → goroutine. **Race**: 2 goroutine + map → data race; sync.Map YALNIZ 2 halda (sənəddən!) → tövsiyə **safeMap** (Mutex embed + Store/Load/Delete).

### Testing (s. 82-100)
Niyə test: double-entry bookkeeping analoqu — checks-and-balances. Sadə → **table-driven** (anonymous struct; edge case 1 sətirdə). Error mesajı: FUNKSİYA(p) = actual, want expected. t.Error vs t.Fatal. **HTTP handler test**: NewRequest + httptest.NewRecorder + handler.ServeHTTP — server YOX. **Mock**: ÖZ interfeys yaz (randIntGenerator; 1 metod) → fixedRandIntGenerator (randomNum qaytar + calledWithN qeyd et) — kitabxana implicit implement edir, seed fiksle → spesifik nəticə test etmək MÜMKÜN. Coverage: -cover → -coverprofile → **go tool cover -html** (62.5% → 100%). **Examples**: `// Output:` şərhi test KİMİ yoxlanılır + godoc-da görünür (ExampleF / ExampleT_M).

### Benchmarks (s. 101-112)
Fibonacci: recursive 1,255,534 ns/op → iterativ `a, b = b, a+b` 20.3 ns/op (**~60,000x!**). benchcmp old/new. **b.ResetTimer** (setup-i sayma). **-benchmem**: B/op + allocs/op; HighMem variant 132 B/op (slice) — escape analysis (-gcflags=-m) stack vs heap. **Modulo vs Bitwise-and**: n%m == n&(m-1) (m=2^n) — ASM-də IDIVQ vs ANDQ → **2x dəqiqlik fərqi** (16.7 → 7.40 ns/op) — benchmark-un öz overhead-i!

### Tooling (s. 113-127)
godoc (CLI/HTTP; şərh konvensiyası; doc.go). **Go Guru** (vim: :GoImplements — interfeysi kim implement edir). **Race detector 2 tam dərs**: Cat (SetNoise fondda + Noise oxu → dəqiq sətri göstərir; həll: struct Mutex embed) və API (normalizeCountry fondda → go run -race + curl — production ssenarisi; həll: WaitGroup gözlə). Go Report Card.

### Security (s. 128-131)
Go versiya yenilə (golang-announce). CSRF → per-session token (nosurf/gorilla/csrf). **HSTS** → headerWrap middleware: `max-age=31536000; includeSubDomains`. CSP (yarımçıq). bluemonday/unrolled/secure.

### CI (s. 132-134)
Travis (.travis.yml; 5 Go versiyası + make lint/test) / Drone (Docker; yarımçıq). **Makefile nümunəsi** (Go Report Card-dan): all→lint+build+test.

### Deployment/Monitoring (s. 135-139)
Deployment — yarımçıq (infrastruktur asılı). Prometheus + Grafana xatırladılır. Optimization — TODO. Gotchas: nil interface (Ch 3-də). Further reading: Spec, Effective Go, Golang Weekly, Go Blog, Gophers Slack.

## Kitabın əsas mesajları

1. **Production-ə hazır = tooling + test + ölçmə + təhlükəsizlik** — "kod yaz" hələ başlanğıcdır.
2. **Deterministik alətlər:** gofmt (format müzakirəsiz), vet, golint — hamısı CI-a asanlıqla inteqrasiya olunur.
3. **Biliyi ölç:** coverage (HTML) + benchmark (ns/op, B/op) + race detector — real rəqəmlər.
4. **Benchmark özü də baqlanır** (modulo tələsi) — mikro-ölçmədə instrumentasiya xərcini bil.
5. **Test = sənəd:** Example + // Output: həm godoc, həm regression test.

## Ən dəyərli 5 fikir

1. **Add sayı = goroutine sayı** — WaitGroup-da az/çox addım erkən qayıdış/deadlock; bu xəta production-da tapılması ən çətin sinifdir.
2. **Nil interface tələsi:** nil pointer-i interface-ə qoymaq → interface nil OLMUR — Go-nun ən məşhur tələsi; book Ch 3 + Gotchas-da 2 dəfə vurğulanır.
3. **n % m == n & (m-1)** (m = 2^n): benchmark dəqiqliyi üçün yox, ÜMUMI optimizasiya üçün də aktual — asm səviyyəsində IDIVQ vs ANDQ fərqi.
4. **Mock pattern:** öz 1-metodlu interfeys + fixed generator + çağırış arqumentlərinin qeydi (calledWithN) — bir dəfə öyrənilən, hər yerdə tətbiq olunan forma.
5. **go run -race production-da:** race yalnız testdə yox — serveri -race build edib real curl ilə yoxla (API nümunəsi).

## Kitabın ən dəyərli hissəsi

**Tooling chapterdə Race Detector bölməsi (s. 115-127)** — çünki:
- 2 TAM workflow: test race + production (go run -race + curl) — hər ikisi step-by-step output ilə
- Xəta çıxışının OXUNMASI öyrədilir (Write at ... by goroutine → dəqiq sətir) — debug bacarığı
- 2 həll müqayisəsi: Mutex-embed vs WaitGroup-un vaxtında gözləmə — hər birinin yeri aydın
- "Test pass = kod düzdür" fikrinin YANLIŞ olduğu göstərilir (pass edən testin race-i görünür)

Bu 13 səhifə kitabın ən practical hissəsidir — müəlliflərin Go Report Card (10,000+ istifadə olunan alət) təcrübəsindən gəlir.
