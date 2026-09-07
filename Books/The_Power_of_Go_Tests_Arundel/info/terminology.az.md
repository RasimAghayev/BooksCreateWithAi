# The Power of Go: Tests — Terminologiya (Azərbaycanca)

Format: **English Term (Azərbaycanca qarşılıq)** — qısa izah

## Test Əsasları (Ch1-3)
- **TDD (Test-Driven Development)** — testdən əvvəl kod; red-green-refactor tsikli
- **Red-Green-Refactor** — fail edən test → keçən test → təmizləmə
- **Fake It Till You Make It** — sadəcə testi keçəcək kod, sonra generalizasiya
- **Want vs Got** — gözlənilən və alınan dəyərin müqayisəsi
- **Table-Driven Test** — input/expected cütləri şəklinə yığılan test
- **Subtest (t.Run)** — test daxilində test; hər table sətri = subtest
- **t.Parallel** — testin paralel icrası; istismarı = qlobal state yoxdur
- **t.Helper** — helper funksiyanı test kimi göstərmə (fayl/sətir düzgün olsun)
- **t.Cleanup** — defer-in test versiyası; test bitəndə icra
- **t.TempDir** — hər test üçün unikal müvəqqəti qovluq (təmizlənmə avtomatik)
- **t.Setenv** — test boyu env dəyişəni (yalnız non-parallel testlərdə)
- **cmp.Diff** — fərqi göstərən human-readable muqayisə (go-cmp)
- **Test Helper** — setup kodunu paylaşan funksiya; özü test DEYİL

## Error Testləri (Ch4)
- **Sentinel Error** — `var ErrX = errors.New(...)` — müqayisə üçün adlandırılmış
- **Type-Wrapped Error** — xüsusi tip; `errors.As` ilə çıxarılır
- **errors.Is / errors.As** — error zəncirində müqayisə / tip çıxarışı
- **Error Value** — Go-da error = qaytarılan DƏYƏR, exception deyil

## Invalid Input (Ch5)
- **Inapproach (yalnış yanaşma)** — funksiyanın mənfi istifadəsi
- **Sərhəd Dəyərlər** — boş, minimal, maksimal, nil, hədd inputları

## Fuzz (Ch6)
- **Fuzz Testing** — random generasiya inputlarla avtomatik bug ovu (Go 1.18+)
- **Fuzz Target** — f.Fuzz-ə ötürülən funksiya: *testing.T + inputlar
- **Seed Corpus (f.Add)** — başlanğıc "training data" inputları
- **testdata/fuzz** — fail inputlarının avtomatik yazıldığı regression faylları
- **Minimizasiya (Shrink)** — fail input-un minimal forması (27→0 bayt)
- **Oracle Pattern** — gözlənilən nəticəni etibarlı mənbədən almaq
- **Property-Based Testing** — dəqiq dəyər yox, invariant yoxlama
- **Invariant** — input-dan asılı olmayaraq dəyişməyən xassə
- **Deterministic Seed** — rand üçün sabit başlanğıc → təkrarlana bilən test
- **rand.Perm** — 0..n dəyərlərinin RANDOM SIRASI (tam örtük + sıra təsadüfü)
- **%[1]x Explicit Argument Index** — eyni arqumentin 2-ci istifadəsi formatta
- **t.Skip** — keçilməli input; pass/fail deyil

## Coverage və Mutation (Ch7)
- **Test Coverage** — testlərin işə saldığı statement %-i
- **-coverprofile + go tool cover -html** — profil faylı + qırmızı/yaşıl görünüş
- **Coverage Instrumentation** — compiler-in statement sayğacları (fuzz ilə ortaq)
- **Goodhart's Law** — ölçü hədəf olsa yaxşı ölçü olmaqdan çıxır
- **Cobra Effect** — stimul → əks nəticə (ilan yetişdirmə)
- **Coverage Ratchet** — coverage-i SALAN check-in qadağası
- **Bebugging (Seeding)** — bilərəkdən bug əkib tapılma % ölçmək (Weinberg)
- **Feeble Test** — kodu cover edir, amma kifayət qədər test ETMİR
- **Mutation Testing** — avtomatik bebugging; syntax tree mutasiyası
- **go-mutesting / gremlins** — Go mutation alətləri
- **Mutation Score** — öldürülən mutant nisbəti (1.0 = hamısı tutulub)
- **Survivable Mutant** — testlərin TUTMADIĞI dəyişiklik = problem

## Untestable (Ch8)
- **Walking Skeleton** — uçtan-uca işləyən ən nazik funksional dilim
- **Hello, Production** — ən sadə versiyanın dərhal prod-a çıxarılması
- **Functional Options** — WithX(...) variadic konfiqurasiya (sadə hal qısa)
- **External / Internal Test** — başqa package (black-box) / _internal_test.go
- **Close-only Channel** — bağlanma özü mesajdır (ctx.Done())
- **Port 0 Trick** — kernel-in random BOŞ port təyin etməsi
- **Wait for Success** — retry loop + timeout (fixed sleep əvəzinə)
- **Flickering Test** — ara-sıra fail; timeout-a yaxın olan test
- **Race Detector (-race)** — data race-lərin icra-zamanı aşkarı
- **Smoke Test** — paralel çağırışla race yaratmaq cəhdi
- **context.WithTimeout / DeadlineExceeded** — standart timeout mexanizmi

## testscript (Ch9)
- **testscript** — rogpeppe/go-internal-dən CLI test DSL-i
- **txtar** — "text archive": skript + daxili fayllar bir faylda
- **exec** — proqram işə sal + exit-0 assert; `! exec` = fail gözlə
- **stdout/stderr** — regex match assertləri
- **TestMain + RunMain** — custom binary-lərin $PATH qeydiyyatı
- **Delegate Main (Main() int)** — main-in test-olunan proxy-si
- **total coverage** — subprocess daxil coverage rəqəmi
- **cmp / cmpenv** — golden fayl muqayisəsi / env-expand variantı
- **grep -count=N** — dəqiq match sayı asserti
- **[condition] / skip** — exec:sh, go1.x, darwin, !arm64 şərtləri
- **$WORK / $exe** — workdir yolu / Windows .exe suffiksi
- **Setup + env.Setenv** — testdən skriptə dinamik dəyər ötürülməsi
- **& + wait** — fon proqramları + gözləmə
- **txtar-c** — kataloqu + skripti txtar-a birləşdirən alət
- **Repro** — bug report-da txtar test case

## Dependency (Ch10)
- **Dependency** — izolyasiyalı işləyə bilməyən komponentin ehtiyacı
- **Toxic Codependency** — qarşılıqlı sıx bağlı komponent dəsti
- **Merge** — ayrı olmayan komponentlərin birləşdirilməsi
- **Scope Reduction** — asılılığı funksiya sərhədindən çıxarmaq
- **Dependency Injection** — asılılığı arqument kimi ötürmək (kitab: şübhəli)
- **Test-Induced Damage** — yalnız-test parametrləri ilə API zədəsi
- **Interface Pollution** — yalnız mock üçün yaradılmış interface
- **Chunking** — davranışı asılılıqsız funksiya dilimlərinə bölmə
- **Magic Package** — mövcud olmayan funksiyanı təsəvvür edib test yazmaq
- **Adapter / Ambassador** — xarici asılılıq kodunun bir komponentdə cəmi
- **sqlmock** — sql.DB üçün stub paketi
- **Stub / Spy / Mock / Fake / Crock** — heç-nə-etməyən / qeyd edən /
  expectation yoxlayan / yüngül-real / bilərəkdən-səhv test ikamediləri
- **Indirect Outputs** — mock-un yoxladığı çağırış zənciri (brittle!)
- **Singleton** — "yalnız bir nüsxə" varsayımı (qlobal DB, time.Now)
- **Seam** — fake-in inject oluna biləcəyi yer
- **`var Now = time.Now`** — vaxt üçün funksiya-dəyişəni seam-i
- **Wall Clock vs Monotonic Clock** — geri gedə bilən / yalnız irəli saat
- **Turn Time Into Data** — vaxtı gizli singleton-dan data-girişə çevirmək

## Suite Sağlamlığı (Ch11)
- **Guerilla Testing** — icazəsiz, toxunduğun yerdə test yazmaq
- **Test What You Touch** — Feathers-in legacy adalar strategiyası
- **Legacy Code** — testsiz + testability-nəzərəsiz kod
- **Edit and Pray / Cover and Modify** — testsiz / testli dəyişiklik
- **Rescue Proseduru** — istənilən funksiya testi + yeni impl + köhnə silinir
- **Zero Defects** — fail-ə icazə yox; dərhal düzəlt
- **Broken Windows** — "hamısı pisdir, mən də pis yazaram" sindromu
- **Bug Bankruptcy / Amnesty** — backlog-un tam ləğvi; vaciblər yenidən açılır
- **Precondition / Postcondition** — ƏMƏLİYYATDAN ƏVVƏLKİ / SONRAKI gözlənti
- **Mirage Test** — keçir amma heç nə test etmir (Create {})
- **Implicit Postcondition** — yazılmamış gözlənti (Bob silinməməlidir)
- **"Halt and Catch Fire"** — test edilməli olmayan aşkar şeylər
- **Seçici Assert** — yalnız scenario-ya aid dəyərlərin yoxlanması
- **cmpopts.EquateApprox** — float epsilon müqayisəsi
- **Scaffolding Tests** — inşaat boyu dəstək, sonra silinən testlər
- **Delta Coverage** — testin unikal əlavə etdiyi örtük
- **Beck Time** — 10 dəqiqəlik psixoloji suite limiti
- **Flaky vs Brittle** — istədiyində fail / əlaqəsiz dəyişiklikdə fail
- **cmpopts.SortSlices** — sıradan asılı olmayan slice müqayisəsi
- **Shared World / Private World** — paylaşılan / hər testə məxsus environment
- **Transaction + Rollback** — DB test izolyasiyası üsulu
- **Prod-safe Tests** — "bir gün prod-a qarşı işləyəcək" prinsipi
