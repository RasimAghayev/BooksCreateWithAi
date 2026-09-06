# Go Optimizations 101 — Müəllim Qeydləri (Teacher Notes)

> Bu sənəd kitabı öyrədərkən istifadə üçün metodiki qeydlər, çətin anlar, müzakirə sualları və praktik tapşırıqlar.

---

## 1. Auditoriya və ön şərtlər

- **Hədəf auditoriya:** Go-nu aktiv yazan intermediate+ developer-lər. Başlanğıc üçün YOX — value/pointer, slice/map mexanikası, goroutine əsasları tələb olunur.
- **Ön şərtlər:** Go syntax, slice-in pointer+len+cap modeli, benchmark əsasları (`go test -bench`).
- **Həcm:** 161 səh., 14 fəsil — 3-4 günlük intensiv və ya 6-8 həftəlik kurs üçün material.

## 2. Tədris axını

| Module | Chapter | Çəki | Fokus |
|--------|---------|------|-------|
| A: Təməl model | 2-3 | 25% | Value part + allocation — HƏRŞEY bunun üzərində |
| B: Stack/GC | 4-5 | 25% | Escape + GC pacer — ƏSAS modul |
| C: Data strukturlar | 6-11 | 30% | Loop, slice, string, BCE, map |
| D: Abstraksiya | 12-14 | 20% | Channel, funksiya, interface xərcləri |

**Qızıl qayda:** Kitab 3 alətə əsaslanır — hər dərsdə canlı istifadə et:
1. `go build -gcflags="-m -m"` (escape + inline cost)
2. `go run -gcflags="-d=ssa/check_bce"` (yoxlamalar)
3. `GODEBUG=gctrace=1` + `go test -bench` (GC + performans)

## 3. Çətin anlar və izah üsulları

### a) Direct/indirect part modeli (Ch 2)
**Vizual alət:** İki şəkil — (1) struct: bütün sahələr BİR qutuda; (2) slice: başlıq qutusu + ayrı element qutusu, arada ox. `s2 := s1` → başlıq kopyalanır, element qutusu ŞƏRƏKLİ. Bu şəkil Ch 3, 4, 8-in hamısına istinad nöqtəsidir.

### b) Escape analysis-in "mühakiməsi" (Ch 4)
İnsan məntiqi ilə izah: "compiler təhlükəsizliyi seçir — əmin ola bilmirsə heap". `fmt.Println` vs `println` demosi həmişə təsirli: sadə print əmrinin 1 allocation fərqi! Sonra `-m` ilə canlı yoxla.

### c) 9→10 sahə sıçrayışı (Ch 2)
Benchmark nəticələrini tələbələrə GÖSTƏRMƏ, özlərinə İCRa etdir — `Add4` 2.6ns vs `Add5` 19ns təəccüb effekti yaradır. Dərs: "4 sahə bir maşın register set-inə sığır".

### d) `m[string(k)]` oku vs yazı fərqi (Ch 9/11)
Sınıq anlayış: "oxuma götür-qaytar, yazı isə saxla". `stat()` demosu ilə get=0, set=1, `++`=1 alloc göstər. Sonra pointer-element və index-table həllərini İRLİ SIRAYLA gətir.

### e) BCE hint-lərinin "işləyən/işləməyən" fərqi (Ch 10)
`is = is[:256]` (işləyir) vs `_ = is[:256]` (işləmir) — niyə? Compiler-yalnız **təyinatların** dataflowunu izləyir, atılan ifadələri yox. Bu, optimizer-in məhdudiyyəti haqqında dərsdir.

### f) Aliasing təhlükəsi (Ch 6)
`g` (lokal yığma) funksiyasının `sum = &s[2]` halında DÜZGÜN OLMADIĞINI birgə kəşf et — "5× sürət" və "yanlış nəticə" birgə təhlükədir. Praktik dərs: **semantik ekvivalensiya olmadan transformasiya yoxdur.**

### g) Boxing xərc cədvəli (Ch 14)
Cədvəli divara as: pointer 1.2ns → string 50× → array 500×. Şagirdlərdə "interface{} ucuzdur" yanılgısını bu cədvəl dağıdır.

## 4. Müzakirə sualları

1. Niyə `a [100]int` funksiya parametri hər çağırışda 100-element kopya deməkdir — Ruby/Python-da niyə belə deyil?
2. `GOGC=off` qoysaq nə olur? 2-dəqiqəlik taymer nə üçün var? (finalizer + stack shrink)
3. Ballast 150MiB "pulsuz"dur — niyə? (virtual allocation, Linux lazy pages) Container-də riski nədir?
4. `sync.Pool` GC-də boşalır — nə vaxt problem, nə vaxt xüsusiyyətdir?
5. `for _, v := range` niyə array-də 2 dəfə kopya yaradır? (Ekvivalensiya kodunu birgə yazın.)
6. Interface metod 8× yavaşdır — buna baxmayaraq standart kitabxana niyə interfeyslərlə doludur? (təmiz dizayn > mikro-perf; de-virtualizasiya)
7. Kitab v1.19 üçün yazılıb — 2024+ Go-da hansı fəndlər artıq keçərsizdir? (PGO, yeni inliner, arena eksperimenti...)

## 5. Praktik tapşırıqlar

**Tapşırıq 1 (ölçü):** Öz layihəndə 3 struct seç; `unsafe.Sizeof` ilə ölç; sahə sırasını dəyişib fərqi hesabla; minlik massivdə ümumi qənaəti proqnozlaşdır.

**Tapşırıq 2 (escape):** 5 funksiyalı mini-layihə: (a) lokal slice qaytar, (b) fmt.Println çağır, (c) interface metoda pointer ötür, (d) `make([]T, n)` runtime ölçü, (e) closure-da dəyişən. Hamısını `-m` ilə proqnozlaşdır → yoxla.

**Tapşırıq 3 (allocation azalt):** Kitabın MergeWithOneLoop funksiyasını optimallaşdır; `AllocsPerRun` ilə 4→1 alloc sübut et; overflow qoruması əlavə et.

**Tapşırıq 4 (BCE):** `-d=ssa/check_bce` ilə öz kodunda 10 yoxlama tap; 3-ünə hint tətbiq et (yenidən təyin, ən-böyük-indeks-əvvəl, qlobal→lokal).

**Tapşırıq 5 (map):** Söz-sayğac üçlüsü yaz: A (birbaşa `m[string]++`), B (pointer element), C (index table) — 3 benchmark-i müqayisə et, GC təzyiqini müzakirə et.

**Tapşırıq 6 (interface boxing):** `Box` funksiyaları yaz: int16, string, `[100]int`, `*T`, `int8` — kitabın cədvəlini ÖZ maşında yenidən qur. Sonra lookup-table ilə uint16-nı 20× ucuzlaşdır.

**Tapşırıq 7 (GC terapiyası):** "Zibil istehsalçısı" proqramı yaz (kitabın garbageProducer): gctrace=1 ilə GOGC=100 vs 1000 vs GOMEMLIMIT ölç; ballast variantını əlavə et; 3 strategiyanın heap/interval qrafikini müqayisə et.

**Tapşırıq 8 (kompleks):** Metin emalı funksiyası yaz (split+count+merge): 5 kitab fəndini tətbiq et — in-place, pre-allocate, EqualFold, `[N]byte` açar, index table. Optimallaşdırma ƏVVƏLİ/SONRAKİ benchmark raportunu təqdim et.

## 6. Sınav sualları

**Asan:**
1. `string` neçə word-dur? Slice? (2; 3)
2. Niyə bütün qlobal dəyərlər heap-dədir? (package-level, hər yerdən görünür)
3. `for i := range s` vs `for _, v := range s` — hansı element kopyası yaradır? (ikincisi)

**Orta:**
4. `m[k]++` vs `m[k] = m[k] + 1` — fərq nədir? (hashləmə sayı)
5. `*(*[64]byte)(d) = *(*[64]byte)(s)` nə edir, nə vaxt üstündür? (slice→array kopyası, ≤64B)
6. `_ = is[:256]` niyə BCE hint-i deyil? (atılan ifadənin dataflowu izlənmir)

**Çətin:**
7. Aliasing varsa `n := *sum` transformasiyası niyə yanlış nəticə verir? Nümunə göstər.
8. Boxing cədvəlində string-in pointer-dən 50× yavaş olmasının kök səbəbi nədir? (non-constant string dublikasiya+heap; string sabitdirsə 1.2ns)
9. Ballast vs GOMEMLIMIT — hər birinin üstünlük/tradeoff'u nədir? (ballast: virtual pulsuz, v1.18- klassik; limit: soft, konteyner-uyğun, seçim çətindir)

## 7. Kollektiv layihə ideyası

**"Hot Path Klinikası":** Hər tələbə real layihəsindən (və ya açıq-mənbədən) bir hot funksiya gətirir. Qrup birgə: (1) benchmark qurur, (2) `-m`/`check_bce`/`gctrace` diaqnozu, (3) 3 fənd tətbiqi, (4) sübutlu raport. Ən yaxşı "əvvəl/sonra" təqdimatları mükafatlandırılır.

## 8. Əlavə resurslar

- Kitabın səhifəsi: go101.org/optimizations/101.html
- Müəllifin əsas kitabı: go101.org (Go 101)
- Rəsmi benchmark dok: pkg.go.dev/testing#B
- Escape analysis istinadı: github.com/golang/go/wiki/CompilerOptimizations
- GOMEMLIMIT rəsmi bələdçi: go.dev/doc/gc-guide
- The Go Memory Model / runtime source: github.com/golang/go/src/runtime
- Go 101 (AZ kontekstdə): bu seriyadakı Go 101 kitabı da oxunub — "From Ruby to Golang" ilə birlikdə əsas yaxşı təşkil edir

## 9. Kitabın dəyəri və limitləri (müəllim qiyməti)

**Güclü:** Mexanizm-əsaslı yanaşma; hər iddianın benchmark sübutu; compiler qaynaqlarına yaxınlıq (Keith Randall kimi Go komandası ilə dialog); "niyə" sualına cavab.

**Zəif/həssas:** v1.19 asılılığı (PGO, arenas, yeni inliner yoxdur); bəzi fəndlər künklü davranışlardan asılı (composite-literal stack, split-concat) və qeyd olunur; dili öyrətmir — yalnız optimallaşdırır.
