# Beyond Effective Go — Part 2 — Xülasə (Azərbaycanca)

**Müəllif:** Corey Scott | **İl:** 2024 | **Səviyyə:** L4 (Advanced) | **ISBN:** 978-0-6455820-8-6

## Kitabın ümumi məqsədi

"Beyond Effective Go" seriyasının 2-ci hissəsi (Chapter 4-9): Effective Go-nun
öldüyü yerdən davam edərək YÜKSƏK KEYFİYYƏTLİ Go kodu yazma sənətini öyrədir.
Müəllif (PluralSight müəllimi, Go veteranı) iddia edir: kod əsasən İNSANLARLA
ünnsiyyət vasitəsidir — performansdan və "ağıllıbuludluq"dan əvvəl
aydınlıq, tutarlılıq və proqnozlaşdırıla bilənlik gəlir. Kitab 6 böyük
mövzunu əhatə edir: software design prinsipləri, code UX, irəli unit-test
texnikaları, developer productivliyi, qeyri-adi pattern-lər və
metaproqramlaşdırma.

## Fəsil-fəsil xülasə

4. **Exploring Software Design Principles (18-68):** Go ideologiyası ilə
   uzlaşan prinsiplər — Unix Philosophy (kiçik, tək işi görən alətlər), DRY
   vs KISS (DRY axiomatik DEYİL — yanlış abstraksiya zərərli), Delegation,
   Composition over inheritance, Accept interfaces/return structs (istifadəçi
   tərəfindən genişlənmə), Single Responsibility/ISP/DIP; 4 klassik pattern-in
   Go-dakı tətbiqi: Singleton (sync.Once, amma adətən YOX), Factory Method
   (closure/constructor funksiyaları), Observer (kanallarla), Adapter
   (struct daxilində interface embed).
5. **Optimizing for Code UX (69-147):** üç sütun — Clarity (aydın adlar,
   qısa funksiyalar, dərin nesting-dən qaçınma), Consistency (layihə daxilində
   vahid üslub; core libraries-dən nümunə götür), Predictability (principle
   of least astonishment; zero value, error contract-ları, API surface
   minimallaşdırma); konkret texnikalar: functional options-un sadə
   alternativləri, context err-inin wrap-lənməsi, defer bloklarının
   oxunaqlığı, struct sahə adlandırma/sıralama.
6. **Advanced Unit Testing (148-236):** test piramidi (70 unit/28 integration/
   2 E2E); behavior yox implementation testi; table-driven testsin full
   nəzəriyyəsi (senari quruluşu, namealanmış struct-lar, paralel mod,
   t.Cleanup); mocks/stubs fərqi (mockery ilə avtomatik generasiya); test
   recorder patterni (call-ların qeydi — davranış yox struktur yoxlaması);
   resilience testləri; az tanınan imkanlar: context timeout testi
   (shortDuration + WaitPlease), testdata/ qovluğu, sync.WaitGroup test
   latch; concurrent kodun testi (race detector, -count; blocking kanal
   patternlərinin idarəsi).
7. **Improving Development Productivity (237-276):** "Be lazy" fəlsəfəsi;
   alətlərin mükəmməl işlədilməsi (gopls, dlv dərinliyi, staticcheck
   class-ları, go doc/vet/coverage/gcflags); IDE bacarığı (refactor
   shortcuts, multi-cursor); git axını (feature branch, rebase vs merge,
   conventional commit-lər, code review praktikası); kiçik commit-lərlə
   işləmək (WIP, interactive rebase); CI avtomatlaşdırma.
8. **Examining Unusual Patterns (277-316):** FP konseptləri Go-da —
   pure functions, immutability (bəzən mümkün), higher-order functions,
   currying (amma Go-da ağır — taxılıq); closures (loop variable, göstərici
   tutma tələsi); function-based pattern-lər: abstract method (interface
   metod-u default implementasiya ilə), middleware zənciri, functional
   options, decorator, function chaining (fluent API), futures (async
   nəticə placeholder — channel ilə); struct fndləri: empty struct{}
   (0 bayt), anonymous struct (lokal ad-hoc tiplər), noCopy (go vet
   üçün marker metod).
9. **Metaprogramming with Go (317-357):** 3 silah: (1) API inteqrasiya
   avtomatlaşdırması — https://api.github.com və s. ilə Go proqramından
   işləmək (net/http, struct tag-lərlə JSON); (2) `os/exec` ilə proqram
   koordinasiyası (Run/Output/CombinedOutput, stdin/stdout pipe-lar,
   environment, exit code idarəetməsi); (3) kod generasiyası — go/ast
   ilə Go kodunun parse-edilməsi + text/template ilə yeni kodun
   generasiyası (mock generasiya, boilerplate azaldılması). Meyarı: alətə
   sərf olunan vaxt < qənaət olunan vaxt.

## Kitabın əsas mesajları

1. **Kod = insan ünsiyyəti:** aydınlıq > ağıllı görünmə; gələcək oxuyan
   (6 ay sonra özün) üçün yaz.
2. **DRY yanlış ola bilər:** təsadüfi oxşarlıq ≠ eynilik; erkən
   abstraksiya = gələcək dəyişikliyin bloklanması.
3. **Test davranışı yoxlamalıdır:** implementation testi refaktoru
   partladır; recorder-lər struktur yox ÇAĞIRIŞ ardıcıllığını yoxlayır.
4. **Be lazy:** alətləri öyrən (gopls/dlv/staticcheck), təkrarı
   avtomatlaşdır, kiçik commit-lərlə işlə.
5. **Metaproqramlaşdırma = son çarə, amma güclü:** API+exec+codegen
   üçlüyü, amma vaxt-hesabı ilə.

## Kitabın güclü və zəif cəhətləri

**Güclü:** hər pattern üçün "niyə/niyə yox" müzakirəsi; real prodakstion
nümunələr; test fəsili müstəqil dərs kimi. **Zəif:** FP materialı Go-nun
ideologiyasına bir qədər "qarşı" — müəllif özü də warn edir; bəzi
texnikalar (currying) Go-da nadir istifadə olunur.
