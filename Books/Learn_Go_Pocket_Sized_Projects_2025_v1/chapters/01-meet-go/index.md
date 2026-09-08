# Chapter 1 — Meet Go (səh. 33-45)

## Bu fəsil nədən bəhs edir?

Kitabın metodu (pocket-sized layihələrlə öyrənmə — John Dewey 1897 "doing
is the best way to learn"), Go-nun mənşəyi (2007, Google: Griesemer/Thompson/
Pike — yavaş build, asılılıq idarəetməsi, mürəkkəb kod problemləri), 25
rezerv açar söz sadəliyi, zəngin alət dəsti (compile/format/dependency/
static-analysis/test/doc/profiling/LSP/trace), 2024 Developer Survey (1.
API/RPC, 2. CLI, 3. kitabxana/framework, HTML web, avtomatlaşdırma, agents;
7% embedded, 4% oyun, 4% AI), karyera sahələri (fintech/medtech/aerospace/
satellite), goroutine-lərin üstünlüyü (resizable bounded stack, bir neçə
KB başlanğıc → yüzlərlə min eyni adres fəzasında), composition + implicit
interfeyslər (miras YOX), generics 1.18 (slice/map filter boilerplate),
compile-time sintaksis xətaları, cloud dəstəyi (Kubernetes/Docker Go-da
yazılıb), Go-nun ZƏİF tərəfləri (GC = tam yaddaş nəzarəti YOX — C/cgo;
binari kitabxana "painfully achievable"; yenidən qurulma məcburiyyəti; böyük
binary — TinyGo; "Googling for Go" → golang axtarış hiyləsi), dil müqayisə
cedvəli (C++/Python/Java/Go: errors are values, goroutines+kanallar, implicit
interfeys, built-in test/bench/fuzz), pocket-layihə fəlsəfəsi (grammar →
service-ə qədər; interfeys implisitliyi = mock/DI dünyası; goroutine-lərə
az diqqət — sadəcə 1 layihə; errors as values), unit test hər yerdə (ch12
istisna), fuzzing (random dəyərlər = vulnerability yoxlama), clean code
(Eagleson's Law — 6 ay sonra öz kodun yad kodur; domain-driven təşkilat),
arxitektura layihələri (HTML/HTTP vs Protobuf/gRPC — sevimli protokolu
SEÇ; logger; anti-corruption layer; keş), alət dəsti + ch12 (microcontrollers,
WebAssembly), side quests (əlavə məşqlər).

## Əsas fikirlər

### 1. Go Niyə Yaradıldı?
- **2007, Google (Griesemer/Thompson/Pike):** yavaş proqram qurulması +
  asılılıq idarəetməsi + mürəkkəblik + cross-language çətinliyi
- **Həll fəlsəfəsi:** memory idarəetməsini aradan qaldır + paralel kodu
  SADƏ et + zəngin alət dəsti (compile→trace hamısı daxili)
- **25 rezerv söz** (2024 noyabr): sadəlik = sürətli öyrənmə

### 2. 2024 İstifadə Sırası
1. API/RPC servis (ən çox!) 2. CLI 3. Kitabxana/framework 4. HTML web
5. Avtomatlaşdırma/data 6. Agents/daemons; 7% embedded; 4% oyun; 4% AI/ML
- Karyera: fintech/medtech/foodtech/gaming/music/e-commerce/aerospace/
  satellite — tələb YÜKSƏK

### 3. Goroutine Üstünlüyü
| Thread (OS) | Goroutine (app səviyyəsi) |
|---|---|
| OS-dan asılı; CPU limitli | runtime idarəli |
| böyük stack | **resizable bounded:** bir neçə KB başlanğıc, böyüyüb-kiçilir |
| minik sayda | **yüzlər min** eyni adres fəzasında |

### 4. OOP-siz OOP xüsusiyyətləri
- Miras YOX → **composition + embedded tiplər** (miras-vari davranış,
  mürəkkəbliyi olmadan)
- **Implicit interfeyslər:** metodu İMPLEMENT etmək kifayətdir (elan YOX) —
  "bilmədən interfeys implement edirsən"; mock/stub/DI dünyası açılır
- **Generics (1.18):** filter/map boilerplate-i kəsdi; type-safe reuse

### 5. Güclü Tərəflər
- Compile-time xətalar (runtime YOX) + FAST build
- API/cloud-native: Kubernetes, Docker, lambda-lar — hamısı Go
- Sürətli oxuma-yazma: turnover ~1 il → başqasının kodu = SƏNİN problemən

### 6. Zəif Tərəflər (namuslu siyahı!)
| Problem | Detal / Həll |
|---|---|
| GC | tam yaddaş nəzarəti YOX → C ailəsi; cgo wrapper (kitabda YOX) |
| Library binary | "painfully achievable"; asılılıq = yenidən build + mənbə |
| OS yazmaq | GC vaxtı/şəkli səndə deyil |
| Böyük binari | cloud-da problem YOX; embedded → TinyGo (ch12) |
| "Googling Go" | axtarışda **golang** işlət |
| Hire çətinliyi | developer üçün YAXŞI (tələb artır) |

### 7. Dil Müqayisəsi (C++/Python/Java/Go)
| | Go |
|---|---|
| Dizayn | prosedural+OOP-vari, çoxparadigmlı |
| Xətalar | **errors are values** (exception YOX) |
| Tiplər | statik |
| Compile | birbaşa binari (VM YOX) |
| Concurrency | goroutine + kanal |
| İnterfeys | implicit (və explicit) |
| Yaddaş | GC |
| Test | **daxili**: test/bench/fuzz |
| İstifadə | web API + cloud |

### 8. Pocket-Layihə Fəlsəfəsi
- **Dewey 1897:** "doing is the best way to learn" — BÜTÜN kitabın əsası
- Sıra: hello world → syntax → ... → **cloud-deploy servis**
- **Goroutine az diqqət:** "Go-da goroutinesiz də EFFEKTLİ proqramlaşdırma
  olar" — kitabda YALNIZ 1 layihədə
- Side quests: optional dərinləşdirmə məşqlər

### 9. Kitabın Öyrənmə Vədləri
- **Grammar:** eyni loop açar sözü; switch-lərdə implicit break; exposal
  (böyük hərf) — Java public/private müqabilində
- **Test:** hər fəsildə unit (ch12 istisna); benchmark (commit-də performans
  regresiya yoxlaması); **fuzzing** (random input = vulnerability tapmaq)
- **Clean code:** Eagleson's Law; **domain-driven** qovluq təşkilatı;
  "nəyi EXPOSE etmək" sualı
- **Arxitektura:** HTTP+HTML VƏ gRPC+Protobuf — öz sevimlini SEÇ; logger
  (stdlib-dən irəli); **anti-corruption layer** (Gordle); keş (generics ilə)

## Əsas terminlər
- Goroutine (resizable stack) — yüngül, miqyaslana bilən paralellik
- Composition/embedding — mirassız davranış birləşdirmə
- Implicit interface — metod icrası = implementasiya (bəyansız)
- Errors are values — exception əvəzinə dəyər-kimi xətalar
- Fuzzing — random inputlarla vulnerability testi
- Anti-corruption layer — xarici API-nin domendən izolyasiyası
- cgo — Go↔C körpüsü (kitabda istifadə olunmur)
- TinyGo — embedded/microcontroller Go compiler (ch12)
- Domain-driven design — biznes domenlərinə görə kod təşkilatı
- Side quest — kitabın optional dərin məşqi
- golang — axtarış motoru üçün Go-nun ləqəbi

## Praktik nəticə

1. **Go seçimi:** API/CLI/cloud üçün ideal; GC tam nəzarət tələb edən
   (OS, real-time) üçün YOX — namuslu gözləntilər qur.
2. **Öyrənmə strategiyası:** bu kitab = layihə əsaslı (pocket-sized);
   unit-test HƏR layihədə; benchmark commit-regressiya aləti kimi.
3. **İnterfeys intuisiyası:** implicit implementasiya mock/DI-ni
   "pulsuz" edir — kitab boyu istifadə olunacaq.
4. **Concurrency gözləntisi:** goroutine güclüdür, amma ESAS deyil —
   sadəlik əvvəl.
5. **Axtarış hiyləsi:** Go sualı → "golang" yaz.

## Mənbə
Pages: 33-45 (PDF 34-46)
