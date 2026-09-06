# Chapter 1 — Go: Simple to learn but hard to master

## Bu chapter nədən bəhs edir?

Bu giriş chapter-i kitabın fəlsəfəsini izah edir: Go niyə populyardır, "simple" ilə "easy" arasındakı kritik fərq, səhvlərdən öyrənmənin elmi əsası və kitabda təsvir olunacaq 7 səhv kateqoriyası.

---

## Əsas fikirlər

### 1. Go outline — Go niyə yaradıldı?

**Tarixi kontekst:** 2007-ci ildə Google tərəfindən yaradılıb. Müasir sistemlər tək şəxs deyil, yüzlərlə/Minlərlə proqramçıdan ibarət komandalar tərəfindən yazılır — kod **readable (oxunaqlı), expressive (ifadəli), maintainable (dəstəklənə bilən)** olmalıdır. Go-nun istifadə sahələri: API-lər, avtomatlaşdırma, databaselər, CLI-lər. Bir çoxları Go-nu **"cloud-un dili"** adlandırır.

**Funksiya baxımından məhrumiyyətlər (bilərəkdən):** Go-da type inheritance, exceptions, macros, partial functions, lazy evaluation, immutability, operator overloading, pattern matching YOXDUR. Go FAQ: *"Why does Go not have feature X? Your favorite feature may be missing because it doesn't fit, because it affects compilation speed or clarity of design, or because it would make the fundamental system model too difficult."*

**Dilin 4 əsas xüsusiyyəti:**

| Xüsusiyyət | İzah |
|-----------|------|
| **Stability (Stabillik)** | Tez-tez yenilənsə də (improvements + security patches) dil stabil qalır — dilin ən yaxşı xüsusiyyətlərindən biri sayılır |
| **Expressivity (İfadəlilik)** | Kodu nə qədər təbii/intuitive yazıb-oxumaq mümkündür; az keyword + problemi həll etməyin məhdud yolları → böyük codebase-lər üçün ifadəli dil |
| **Compilation (Kompilyasiya)** | Sürətli build — dil dizaynerlərinin şüurlu hədəfi; productivity-ni təmin edir |
| **Safety (Təhlükəsizlik)** | Strong, statically typed → strict compile-time qaydaları → kod çox halda type-safe |

**Concurrency primitivləri** (goroutines, channels) dizayndan gəlir — səmərəli konkurent tətbiqlər üçün xarici kitabxanaya ehtiyac azdır.

---

### 2. Simple doesn't mean easy — əsas tezis

**Fərq:**
- **Simple** = öyrənmək/anlamaq mürəkkəb deyil
- **Easy** = az səylə nail olmaq mümkündür

Go-nun əsas xüsusiyyətləri bir gündə öyrənilə bilər → **simple to learn**. Amma **easy to master** deyil.

**Konkurentlik nümunəsi:** 2019-cu ildə "Understanding Real-World Concurrency Bugs in Go" tədqiqatı (ASPLOS 2019) Docker, gRPC, Kubernetes kimi populyar Go repositoriyalarını analiz etdi. **Əsas nəticə:** blocking bug-ların əksəriyyəti channels vasitəsilə message-passing paradigmasının qeyri-dəqiq istifadəsindən yaranır — buna baxmayaraq, ümumi inanc message-passing-in shared memory-dən daha asan və az xətaya meyilli olduğu istiqamətindədir.

**Dərs:** Bu, message passing vs shared memory mübarizəsi deyil. Go developer-i kimi konkurentliyi dərindən başa düşməli, müasir prosessorlara təsirini, hansı yanaşmanın hansı halda üstünlük təşkil etdiyini və ümumi tələlərdən qaçmağı bilməliyik. Channels/goroutines öyrənmək simple-dır, praktikada isə easy deyil.

Bu leitmotif bütün dil aspektlərinə şamil olunur — proficiency (peşəkarlıq) üçün vaxt, səy və **səhvlər** tələb olunur.

---

### 3. Səhvlərdən öyrənmənin elmi əsası

- **2011 (neyroelm):** Beynin ən çox böyüdüyü an səhvlə qarşılaşdığımız andır (Moser, Schroder və b. — "Mind Your Errors", Psychological Science).
- **Janet Metcalfe ("Learning from Errors"):** Səhvlərin **facilitative effect (asanlaşdırıcı təsiri)** var — yalnız xətanın özünü deyil, onu əhatə edən **konteksti də xatırlayırıq**. Odur ki, səhvlərdən öyrənmə bu qədər effektivdir.

Kitab bu effekti gücləndirmək üçün hər səhvi **real-world nümunələrlə** müşayiət edir — nəzəriyyə deyil, şüurlu qərarların əsası olan rationale (dəyərləndirmə əsası) verilir.

---

### 4. 7 səhv kateqoriyası

| # | Kateqoriya | Mahiyyəti | Nümunə mövzular |
|---|-----------|-----------|-----------------|
| 1 | **Bugs** | Proqram xətaları — data races, leaks, logic errors, digər defektlər | 2020 Synopsys hesablaması: yalnız ABŞ-da software bug-ların xərci **$2 trillion-dan çox**; Therac-25 (radiasiya terapiyası maşını) race condition səbəbindən gözləniləndən yüzlərlə dəfə böyük doza verib, 3 xəstənin ölümünə səbəb olub |
| 2 | **Needless complexity (Lazımsız mürəkkəblik)** | Xəyali gələcəklər üçün abstraksiya qurmaq — konkret problemi indi həll etmək əvəzinə "evolyasion" proqram yazmaq | Gələcək ehtiyaclar üçün interface/generics dizaynı |
| 3 | **Weaker readability (Zəif oxunaqlılıq)** | Kodu oxumaq vaxtı yazmaqdan 10:1 nisbətində çoxdur (Robert C. Martin, Clean Code) | Nested code, data type representations, named result parameters-in istifadəsizliyi |
| 4 | **Suboptimal or unidiomatic organization (Suboptimal/qeyri-idiomatik təşkilat)** | Layihə/kod təşkilinin qeyri-optimal qurulması | Project structure, utility packages, init functions |
| 5 | **Lack of API convenience (API rahatlığının çatışmazlığı)** | İstifadəçi üçün qeyri-user-friendly API → az ifadəli, başa düşülməsi çətin, xətaya meylli | `any` tiplərinin həddindən artıq istifadəsi, yanlış creational pattern, OOP standartlarının kor-koranə tətbiqi |
| 6 | **Under-optimized code (Optimallaşdırılmamış kod)** | Performance və accuracy itkiləri | Floating-point dəqiqliyi, zəif paralelləşdirmə, allocation azaltma bilməməsi, data alignment |
| 7 | **Lack of productivity (Məhsuldarlıq çatışmazlığı)** | Dilə hakim olmamaq → yavaş iş | Effektiv test yazma, standard library-dən istifadə, profiling tools, linter-lər |

**Bugs haqqında vacib qeyd:** Testlər bug-ları erkən tapmağın yoludur, amma vaxt çatışmazlığı/mürəkkəblik səbəbindən bəzi hallar buraxıla bilər — odur ki, ümumi bug-lardan qaçmaq üçrüstdür vacibdir. Bug-lar yalnız pul deyil — işimizin nə qədər təsirli olduğunu xatırlamalıyıq.

---

## Əsas terminlər

| Termin | İzah |
|--------|------|
| Simple vs. easy | Simple = öyrənmək mürəkkəb deyil; easy = az səylə nail olmaq — Go simple-dır, amma easy deyil |
| Stability (Stabillik) | Diliin tez yenilənməsinə baxmayaraq kompatibilliyin qorunması |
| Expressivity (İfadəlilik) | Kodun nə dərəcədə təbii yazılıb oxunması — az keyword + həll yollarının məhdudluğu |
| Facilitative effect (Asanlaşdırıcı təsir) | Səhvlərin konteksti birlikdə yaddaşda saxlanması effekti |
| Message passing vs. sharing memory | İki konkurentlik paradiqması — kitabda ikisinin düzgün seçimi müzakirə olunur |
| Needless complexity (Lazımsız mürəkkəblik) | Xəyali gələcək ehtiyaclar üçün yaradılan abstraksiya |
| Readability ratio (Oxuma nisbəti) | Oxuma:vaxt yazma = 10:1 (Clean Code) — kod oxunaqlılığı zaman ölçüsündə proqramlaşdırma |
| Race condition (Yarış şəraiti) | Therac-25 faciəsinin səbəbi — konkurrentlik xətası növü |
| Proficiency (Peşəkarlıq) | Dilə dərindən hakim olma səviyyəsi — vaxt + səy + səhvlər tələb edir |

---

## Praktik nəticə

1. **Feature sayı dil keyfiyyəti deyil:** Go-nun məhrumiyyətləri (inheritance yoxdur, exceptions yoxdur və s.) bilərəkdəndir — kompilyasiya sürəti və dizayn aydınlığı üçün.
2. **4 meyar:** Komanda şəraitində dil seçərkən stability, expressivity, compilation, safety prizmasından baxın.
3. **Simple ≠ easy:** Channels/goroutines bir gündə öyrənilir, amma doğru istifadəsi (blocks, races, deadlocks) dərin biliyi tələb edir — kitabın 7-9-cu chapter-ləri məhz bunun üçündür.
4. **Səhvlərdən öyrənin:** Xətanı konteksti ilə birgə xatırlayırıq — bu kitabın metodologiyasının elmi əsasıdır.
5. **7 kateqoriya öz-özünü diaqnoz etmək üçündür:** Kodu nəzərdən keçirərkən bu kateqoriyaları check-list kimi istifadə edin — bug riski, mürəkkəblik, oxunaqlılıq, təşkilat, API rahatlığı, optimallıq, məhsuldarlıq.

---

## Mənbə

- Kitab: *100 Go Mistakes and How to Avoid Them* — Teiva Harsanyi, Manning, 2022 (ISBN 9781617299599)
- Chapter 1: "Go: Simple to learn but hard to master", book səh. 1–6
- PDF səhifələri: 21–26
- İstinadlar: Go FAQ (https://go.dev/doc/faq); T. Tu, X. Liu et al., "Understanding Real-World Concurrency Bugs in Go" (ASPLOS 2019); Moser et al. (Psychological Science, 2011); Metcalfe, "Learning from Errors" (2017); Synopsys, "The Cost of Poor Software Quality in the US" (2020); R. C. Martin, *Clean Code* (2008)
