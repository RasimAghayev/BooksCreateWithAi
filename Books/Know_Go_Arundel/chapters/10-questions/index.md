# Chapter 10 — Questions (Suallar)

## Bu fəsil nədən bəhs edir?

Generics-in praktiki tətbiqi ilə bağlı ən çox verilən suallara cavablar:
nə qədər öyrənmək lazımdır (çox az — istifadəçi gözləmir, yazan az bilir),
kod yazış tərzini dəyişəcəkmi (xeyr — slices/maps yeni API, amma məcburi
deyil), mövcud kodu dəyişmək lazımdırmı (xeyr — "if it isn't broken, don't
fix it"; abstraksiya PULSUZ deyil), performans təsiri (runtime: EYNY —
compile-time fenomen; interface/reflection əvəz edilsə DAHA SÜRƏTLİ —
tidwall/btree 2x; compile: nəzərəçarpmaz), niyə <> deyil [] (parser
ambiguity: `w < x, y > (z)`), hələ olmayanlar (option/Maybe tiplər, enum-lar,
tagged union, makroslar, parameterised metodlar, generic paketlər,
covariance, currying), "Go 2" olmayacaq (breaking change fəlakəti —
Python/IPv6 dərsləri; "Go 2" label = Go 1-ə alınmayacaq), amma bir gün Go-nun
varisi ola bilər (Go → C əlaqəsi kimi).

## Əsas fikirlər

### 1. Nə Qədər Bilmək Lazımdır?
- **İstifadəçi üçün:** demək olar HEÇ NƏ — generic funksiya çağırmaq/
  generic tip istifadə etmək çox vaxt onların generic olduğunu bilmədən
  olur (type inference)
- **Yazan üçün:** type parameter sintaksisi + əsas constraint-lər
  (comparable, cmp.Ordered) — başlamaq üçün kifayətdir
- **Master üçün:** çox kod YOX, çox kod OXU — oxuyan generics-i bilməlidir,
  amma type inference qaydalarını əzbərləmək lazım deyil (compiler yazanlar
  üçün "postgraduate-level")
- **Nəticə:** "Go bilmək" artıq generics bilməyi də tələb edir — istifadə
  olmasa belə

### 2. Kod Tərzi və Mövcud Kod
- **Tərz dəyişikliyi:** yalnız slices/maps kimi yeni kiçik API-lər — vaxtın
  qənaəti, amma məcburi deyil; goroutine kimi: lazım olanda gözəldir,
  məcbur etmək olmaz
- **Mövcud kod:** backwards compatibility zəmanəti — heç nə SİNMİR;
  dəyişmək yalnız fayda xərci örtəndə ("if it isn't broken, you don't
  need to fix it")
- **Abstraksiya xərci (Griesemer):** "Genericity introduces abstraction,
  and needless abstraction introduces complexity. Move cautiously!"
  - io.Reader kimi interface YƏTİR — onu generic-lə əvəz etməyin faydası YOX
  - Generic-lə interface-i əvəz etmək YERİNƏ YETİR: tip təhlükəsizliyi
    artır (konkret tip, interface indirection yox) və ya performans

### 3. Performans — Dəqiq Cavab
- **Runtime:** generics-in ÖZÜ proqramı nə sürətləndirir, nə yavaşıdır —
  pure compile-time fenomen; hər generic instantiate olunur, binary
  "fractionally" böyük ola bilər, amma eyni sürətlə işləyir
- **SÜRƏTLƏNMƏ MÜMKÜN:** interface value/reflection ARADAN ÇIXARILSA —
  tidwall/btree: generic versiya any versiyasından ~2x sürətli (runtime
  indirection yoxdur)
- **Compile time (Dan Scales):** instantiate üçün "very small extra work" —
  normal funksiya əlavə etmək kimi, "hardly noticeable in a full compile"
- **Yekun cavab:** "faster, slower, or about the same — sualdan asılı"

### 4. Niyə <> Yox, [] Belə?
- **Orijinal təklif:** mötərizələr `()` — funksiya parametrləri ilə simmetriya;
  amma `func Foo(T, U any)(p T, q U) (T, U)` — ÇOX mötərizə, oxuma çətin
- **Niyə <> MÜMKÜN DEYİL:** `<` və `>` onsuz da operatorlardır —
  parser ambiguity:
```go
a, b = w < x, y > (z)
// (w < x) && (y > (z))?  YOXSA  w<x, y>(z) — generic instantiation?
```
  Tək cavab YOXDUR → compiler unambiguous tələbi pozulur → [] seçildi
- **Nəticə:** "Worst. Syntax. Ever." şikayətlərinə baxmayaraq — square
  brackets FİNE, alışacaqsan

### 5. Generics-in Həll ETMƏdiyi Problemlər
| Problem | Vəziyyət | Yaxlaşdırma |
|---|---|---|
| Option/Maybe type | YOXDUR | pointer+nil; "something and error" idiomu; optional paketlər (generics yazmağını asudə edir) |
| Enum | YOXDUR | sabitlər (amma qadağan edə bilmirsən); struct+validating accessor |
| Tagged union | YOXDUR | pointer bir növ union (düz göstərici VƏ ya nil) |
| Makros/template metaprogramming | YOXDUR və İSTƏNMİR | generics yalnız tip instantiate edir; compile-time hesab YOX; specialization YOX |
| Parameterised metod | YOXDUR | facilitators pattern (Jaana Dogan) |
| Generic package | YOXDUR | paket instantiate oluna bilmir; uses "packages-ə bölünmür" |
| Operator metodları | YOXDUR | struct ilə +, > istifadə etmək olmaz — "too late" |
| Covariance/contravariance | YOXDUR | inheritance yoxdursa ehtiyac da yoxdur |
| Currying/variadic type params | YOXDUR | — |

- **Option-type dərinliyi:** "something and error" cütünü bir növ option
  tipi kimi düşünmək olar; struct-based option (optional paketi) generics
  ilə asanlaşır

### 6. "Go 2" Olmayacaq
- **Ian Lance Taylor:** "There are no plans for anything called Go 2."
- **"Go 2" label-in mənası:** rədd edilən yox — Go 1-ə alınMAYACAQ, amma
  nəzərdən keçirməyə dəyər təkliflər
- **Niyə versiya 2 fəlakətdir (semver tərifinə görə versiya 2 = versiya 1
  ilə UYĞUNSUZ):**
  1. Ekosistemi məhv edir — mövcud proqramlar işləməz, avtomatik köçürmə
     mümkün deyil (Python 2→3 dərsi)
  2. "Daha yaxşı" ≠ avtomatik qəbul (IPv6 dərsi) — yeni gələnlər gözləyir
  3. Köhnə versiya heç vaxt tam ölmür → cəmiyyət parçalanır
  - Nəticə: "eyni adlı iki fərqli dil" — yeni dilsə YENİ AD ver
- **Amma VARİS olacaq (bir gün):** Go-nun ən yaxşı ideyaları + ən yaxşı
  təkmilləşdirmələr + dil dizaynının son nailiyyətləri → Go-dan Go-nun
  C-dən fərqləndiyi qədər fərqli YENİ dil. "It just won't be Go 2."
- **Konservativlik səbəbi:** tənbəllik YOX — consistency/stability dəyəri;
  dəyişən dil/API üzərində proqram saxlayan hər kəs stabilliyi qiymətləndirir

## Əsas terminlər
- Backward compatibility (geriyə uyğunluq) — 1.x heç nə pozmur
- Abstraction cost (abstraksiya xərci) — oxunaqlılıq sadəliyi qurbanı
- Parser ambiguity (təhlil ikiliyi) — <> operator kimi oxunardı
- Option/Maybe type — dəyər VƏ ya "dəyər yox" tipi
- Tagged union (nişanlı birlik) — icazəli TİPLƏR dəstindən bir dəyər
- Template metaprogramming — compile-time kod icrası (Go-da YOX)
- Facilitators — parameterised metod üçün workaround pattern
- Semantic versioning (semver) — major 2 = uyğunsuzluq

## Praktik nəticə

1. **Generics yazmaq üçün səbəb:** konkret fayda (tip təhlükəsizliyi,
   interface/reflection aradan qalxması, 2x sürət) — yoxsa interface/adi
   kod qalsın; "move cautiously".
2. **Heç nəyi MƏCBURİ dəyişmə:** köhnə kod işləyirsə — işləsin; slices/maps
   yalnız ZAMAN qənaətidir.
3. **Sualların qısa cheat-i:** performans "eyni və ya sürətli"; compile
   "fərqsiz"; sintaksis "[] — çünün <> parse olunmur"; Go 2 "olmayacaq,
   amma varis olacaq".
4. **Dil təkmilləşməsi Gözləntiləri:** breaking change = yeni dil;
   non-breaking = Go 1.x içində tədricən (error wrapping, modules,
   generics kimi).

## Mənbə
Pages: 177-191 (PDF 178-192)
