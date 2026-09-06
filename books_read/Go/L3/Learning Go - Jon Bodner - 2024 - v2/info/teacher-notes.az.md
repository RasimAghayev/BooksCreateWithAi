# Learning Go, Second Edition — Müəllim Qeydləri (Azərbaycanca)

## 📖 Kitab deyir:

Kitab Go-nu "practical" dil kimi təqdim edir: hər texniki seçimin arxasında aydınlıq,
saxlanıla bilənlik və böyük komandalar üçün uyğunluq dayanır. 15 fəsil boyunca:
mühitdən başlayır (fmt/vet/Makefile mədəniyyəti), primitiv və kompozit tiplərlə dili qurur,
funksiya/pointer/metod/interfeys ilə abstraksiya qurur, error/modul/concurrency ilə real
layihə həyatına keçir, standart kitabxana və context ilə servis yazmağa gətirir, testlə
keyfiyyəti garantiliya alır və reflect/unsafe/cgo "əjdahaları" ilə sərhədləri göstərir.

Kitabın ən məqamlı tövsiyələri:
- "Types are executable documentation" — tip adları özü-özünü izah edir
- "Accept interfaces, return structs"
- "Share memory by communicating; do not communicate by sharing memory"
- "Make sure your program benefits from concurrency" — 5 mərhələli paralellik xəyal qırığı
- "Code coverage is necessary, but it is not sufficient"
- "If your function returns values, never use a blank return"

## 👨‍🏫 Müəllim qeydi:

Məncə kitabın praktikada tətbiqində bir neçə incəlik var ki, kitab səthi toxunur, amma
real layihələrdə həlledici olur:

1. **Interface-ləri "böyük layihə quranda" təyin et, əvvəlcədən YOX.** Kitab "client
   tərəfdə təyin et" deyir — bu düzgündür, amma yeni developer bunu "hər funksiya üçün
   əvvəlcədən interface yaz" kimi başa düşür. Əsl həll: konkret tiplərlə yaz, ehtiyac
   (mock/test/ikinci implementasiya) yarananda interface çıxart — Go-nun implicit
   interfeysləri bunu sıfır xərclə edir. Bu, "rule of three"-ün Go tərcüməsidir.

2. **Error handling-də qatmanı kitabdan sonra da dərinləşdir.** Kitab %w/Is/As öyrədir,
   amma real sistemlərdə hər error-u wrap etmək zəruri deyil — kitab da deyir. Praktik
   tövsiyəm: layihə sərhədlərində (HTTP handler, queue consumer) error-a kontekst əlavə
   et, daxili zəncirlərdə isə sentinel/typed error-ları olduğu kimi burax ki, yuxarı
   təbəqələr onları dəqiq tanısın.

3. **Generics fəsli kitabın ən köhnəlmiş hissəsidir.** Kitab Go 1.18-dən əvvəl yazılıb
   ("expected in Go 1.18"). Bugünkü oxucular üçün: type parametrləri artıq dilin sabit
   hissəsidir; `slices` və `maps` standart paketləri, generic data strukturları üçün
   `cmp.Ordered` type list-lərin qanuni formasıdır. Fəslin konseptual dərsləri (any vs
   comparable, metod tipli parametr qadağası) hələ də doğrudur.

4. **Shadowing-dən tam qorunma yoxdur — build-in zorla.** Kitab shadow linter deyir, amma
   univarsal-blok kölgələməsini heç bir alət tutmur. CI-da `golangci-lint`-ə `govet
   shadow` + `gochecknoglobals` + `errcheck` əlavə et — kitabın tövsiyələrinin 80%-ni
   avtomatikləşdirər.

5. **Chapter 10-dakı "concurrency is not parallelism" mövzusu tələbələrə ən çətin gəlir.**
   Bunun ən yaxşı izahı: concurrency = strukturlaşdırma aləti (məsələn, 3 asılı olmayan
   şəbəkə çağırışı paralel NƏSƏ olsa da, struktur sadəcə daha aydındır), parallelism =
   eyni anda icra (hardware qərarı). Tədrisdə "nə vaxt faydalıdır" 3 şərtini vurğula:
   I/O + asılılıq olmama + vaxt limiti.

## Ən vacib 5 fikir

1. Düzgün Go "sıxdırıcı"dır — sadəlik uzunmüddətli saxlanıla bilənlik üçün dəyərli xüsusiyyətdir.
2. Go-nun tipo sistemi (implicit interface, açıq conversion, composition) — özünü sənədləşdirən kod.
3. Error = dəyər, exception = gizli nəzarət axını —前者 always wins for large teams.
4. Channel-lər data axınını, mutex-lər dəyər qorumasını təmsil edir — hansını seçim
   problemi strukturuna başlanğıc verir.
5. Paralellik həmişə faydalı deyil — ölç (benchmark), sonra inan.

## Kitabın ən dəyərli hissəsi

Chapter 7 "Types, Methods, and Interfaces" — embedding-in inheritance olmadığını, implicit
interfeysin duck typing ilə static typing-ni birləşdirdiyini və DI-nin frameworksüz necə
təbii olduğunun göstərilməsi. Bu, kitabın fəlsəfi nüvəsidir. Chapter 10 (concurrency) isə
ən praktik tətbiqi tamamlayır.
