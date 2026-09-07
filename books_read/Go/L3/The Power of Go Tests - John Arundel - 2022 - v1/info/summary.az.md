# The Power of Go: Tests — Xülasə (Azərbaycanca)

**Müəllif:** John Arundel | **Nəşriyyat:** Bitfield Consulting | **İl:** 2022 | **Səviyyə:** L3 (Intermediate)

## Kitabın ümumi məqsədi

John Arundelin "The Power of Go: Tests" kitabı Go-da test YAZMAĞI deyil, testlə DÜŞÜNMƏYİ
öyrədir. Fəsil-fəsil inkişaf edən mövzu xətti: sadə "want vs got" müqayisəsindən başlayır,
table-driven testlər, subtestlər, t.Parallel, effektiv fail mesajları (cmp.Diff),
error testləri (sentinel/type-wrapped/`errors.Is`), invalid-input testləri (table-driven
inappropları), fuzz testing (Go 1.18+ builtin: Fuzz*, f.Add, corpus, testdata/fuzz,
minimizasiya, oracle pattern, %[1]x explicit argument index), mutation testing
(go-mutesting, bebugging, feeble test anlayışı, coverage = siqnal hədəf deyil),
"untestable"ların testi (walking skeleton, Hello Production, functional options,
syncron API üstünlüyü, wait-for-success + time.Timer, port 0 fəndi, race detector,
context timeout, io.Reader/Writer parametrləri, "Mary Jo" boşluq bug-ı, bufio.Scanner,
main→Main(os.Args) delegate), testscript DSL (exec/stdout/!/cmp/txtar/stdin/conditions/
env/Setup/&-wait, standalone runner, txtar-c, bug repro-lar), dependency strategiyaları
(merge, scope azalt, chunking, adapter pattern, Store interface + mapStore, sqlmock,
mock-lara qarşı müdafiə, `var Now = time.Now` seam, "time into data", wall vs
monotonic clock) və sonda bütün suite-in sağlamlığı (guerilla testing, legacy adalar,
10-addımlıq review checklist, optimistic testlər + preconditionlər, persnickety testlər,
fingerprint property testi, delta coverage, "import testing", flaky/brittle fərqi,
shared worlds, zero defects, bug bankruptcy, Beck time = 10 dəqiqə).

Kitabın fəlsəfəsi: **test = dizayn alətidir.** Test yazmaq çətindirsə — dizaynda problem
var ("tests are a canary in a coal mine"). Hər fəsildə təkrarlanan sual: "What are we
really testing here?" — user-visible DAVRANIŞI yoxla, implementasiyanı yox.

## Fəsil-fəsil xülasə

1. **Programming with confidence:** TDD fəlsəfəsi, red-green-refactor, kiçik addımlar,
   "fake it till you make it", davranış spesifikasiyası kimi testlər, inam psixologiyası.
2. **Tools for testing:** go test, go vet, go-cmp (cmp.Diff), test coverage, benchmark,
   t.Helper, t.Cleanup, t.TempDir, t.Setenv, test helper dizaynı.
3. **Communicating with tests:** adlar (TestX_Y formatı), table-driven testlər, subtestlər
   (t.Run), t.Parallel, data-driven dizayn, maraq balansı, ölçü qərarları, ad + input
   göstərən fail mesajları, DRY qəddarcasına yoxlamalar.
4. **Errors expected, exceptions prevented:** error DƏYƏRİ kimi, error testləri,
   sentinel errors, type-wrapped, errors.Is/As, "error var" testləri, error yoxluğu testləri.
5. **Users shouldn't do that:** invalid/empty/nil/huge input testləri, table-driven
   inapproplar, sərhəd dəyərlər, paranoid olmaq sənəti, əyləncəli xətalar.
6. **Fuzzy thinking:** rand determinizmi (seed=1), rand.Perm, property-based testing
   (invariantlar: Square non-negative), fuzz (Fuzz*, f.Add, fuzz target, -fuzztime,
   testdata/fuzz faylları = regression testlər, minimizasiya 27→0 bayt, FirstRune
   rune/byte bug-ı, oracle pattern, t.Skip, %[1]x).
7. **Wandering mutants:** coverage (go test -cover, -coverprofile, go tool cover -html),
   compiler instrumentationu, Goodhart + Cobra Effect, Add→42 = 100% coverage 0% test,
   coverage ratchet, bebugging (!= → ==), feeble/unreachable kod, mutation testing
   (go-mutesting, score, IsEven memoizasiya nümunəsi, "təkrar input" dərsi).
8. **Testing the untestable:** walking skeleton, Hello Production, make-problem-easy,
   functional options (NewLoadTester), external vs internal testlər, sinxron API
   (ListDirectory 3 variant, callback/fs.WalkDir), randomLocalAddr (port 0), eager
   uğursuzluğu, sleep flakiness, wait-for-success, waitForServer, t.Fatal goroutine
   məhdudiyyəti, race detector (-race), global map data race → Store+mutex,
   context.WithTimeout (DeadlineExceeded), Greet refaktoru (Fprint/Fscanln, io.Reader/
   Writer, "Mary Jo" boşluq bug-ı, bufio.Scanner), CLI: main→timer.Main(os.Args),
   ParseArgs/Sleep, NewTimerFromArgs.
9. **Flipping the script:** testscript (go-internal), exec/stdout/stderr/!, qoşma
   pravilaları (double quote literal!), hello.Main delegate + TestMain/RunMain +
   os.Exit tələsi, total coverage, cmp/cmpenv golden fayllar, exists/grep/-count,
   txtar formatı, stdin (fayl/stdin stdout), cp/mv/mkdir/cd/rm/symlink, shell fərqləri
   (exec sh -c), fazalar (şərhlər), conditions ([go1.18], [darwin], [$exe]), env
   ($WORK, HOME=/no-home), Setup+env.Setenv, & fon + wait, standalone runner (shebang),
   txtar-c, "! exec go test" fail-olmalı nümunələr.
10. **Dependence day:** toxic codependency, merge, EmailUserIfAccountIsNearExpiry →
    IsNearExpiry + EmailRemewalReminder, DI-ə şübhə ("dependency özü problemidir"),
    test-induced damage (apiURL), chunking (FormatURL/ParseResponse), router switch→map,
    adapter pattern (ambassador), Store interface + mapStore (mutex) + PostgresStore +
    sqlmock stub, fakes/stubs/spies/mocks/crock, mock-lara qarşı (indirect outputs,
    interface pollution), time.Now singleton, OneHourAgo flaky midnight, `var Now =
    time.Now` seam (paralel olmayan test), wall vs monotonic clock, Sub().Abs() delta.
11. **Suite smells:** no tests (guerilla testing, Uncle Bob), legacy adalar
    (test-what-you-touch, rescue proseduru, Edit and Pray), Feathers eksperimenti,
    tests-first/last (dogma yox), 10-addımlıq review checklist, optimistic testlər
    (precondition: Alice ƏVVƏL yoxdur; Bob sağ qalır; createUserOrFail), persnickety
    (halt-and-catch-fire, seçici assert), fingerprint property (MD5→SHA256 refactor-proof),
    float epsilon (EquateApprox, tuning), delta coverage, scaffolding, import "testing"
    (assert.Equal premature abstraction), flaky vs brittle (map sırası, SortSlices,
    "thou shalt not suffer a flaky test to live"), shared worlds (private world,
    transaction rollback, prod-safe), zero defects, broken windows, bug bankruptcy,
    Beck time 10 dəq + 7 sürət resepti.

## Ən vacib 5 fikir

1. **"What are we really testing here?"** — user-visible davranışı müəyyənləşdir;
   implementasiya detallarını yox. Bu sual hər çətin testə açar verir (email göndərilməsi
   yox, qərar; DB-yə yazılma yox, "would cause").
2. **Test çətinliyi = dizayn siqnalı:** flaky/brittle/untestable — hamısı decoupling
   tələb edir. Test-first yazsan untestable funksiya YAZA BİLMƏZSƏN ("We simply can't
   write untestable functions when the test comes first").
3. **Fuzz + mutation = dərinlik:** icra (coverage) ≠ düzgünlük; fuzz gözlənilməz
   inputları, mutation feeble testləri üzə çıxarır; hər ikisi adi `go test`-ə
   regression kimi qayıdır (testdata/fuzz, score).
4. **Asılılıqlarla mübarizə arsenalı:** merge, scope azalt, chunk (FormatURL/
   ParseResponse), adapter (Store interface), funksiya seam (`var Now = time.Now`),
   testscript (binary səviyyəsində). Mock — SON option; interface-i yalnız mock üçün
   YARATMA (interface pollution).
5. **Suite sağlamlığı prosesdir:** precondition/postcondition intizamı, seçici
   assertlər, delta coverage ilə budama, flaky testlərə toleranssızlıq ("bad tests
   are worse than no tests"), zero defects, Beck time ≤ 10 dəq, və icazə gözləmədən
   test yazmaq — "guerilla testing".

## Kitabın auditoriyası

- Orta səviyyəli Go proqramçıları: test yazırlar, amma testlərin KEYFİYYƏTİNİ
  artırmaq istəyirlər
- TDD-yə yeni başlayanlar: fəlsəfə + praktika bir yerdə
- Legacy kod üzərində işləyənlər: adalar strategiyası, rescue, refaktor-to-testable
- Team lead/review-lər: 10-addımlıq checklist, suite qoxularının diaqnozu
