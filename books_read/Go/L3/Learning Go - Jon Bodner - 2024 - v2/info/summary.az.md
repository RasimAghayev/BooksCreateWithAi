# Learning Go, Second Edition — Xülasə (Azərbaycanca)

**Müəllif:** Jon Bodner | **Nəşriyyat:** O'Reilly Media | **İl:** 2024 | **Səviyyə:** L3 (Intermediate)

## Kitabın ümumi məqsədi

"Learning Go" — Go dilini idiomatik yazmağı öyrədən, 15 fəsillik fundamental kitab.
Məqsəd sadəcə sintaksis yox — Go-nun **niyə** belə qurulduğunu və Go icmasının hansı
praktikaları müdafiə etdiyini göstərməkdir. Kitabın əsas tezisi: düzgün yazılmış Go
kodu "sıxdırıcı"dır (boring) — birbaşa, bəzən təkrarlanan — və bu, böyük komandalar
üzerindən illər keçən saxlanıla bilən proqram təminatının açarıdır.

## Fəsil-fəsil xülasə

1. **Mühitin qurulması:** go tool, GOPATH workspace-i, `go run`/`go build`, `go install`,
   `go fmt`/`goimports` format məcburiyyəti, lint/vet alətləri, Makefile build zənciri,
   versiya uyğunluq vədi.
2. **Primitiv tiplər:** zero value konsepsiyası, untyped literal-lar, int/float seçim
   qaydaları (`int` və `float64` default; pul üçün float QADAĞAN), `var` vs `:=`, untyped
   const-un çevikliyi, unused variable compile xətası, camelCase.
3. **Kompozit tiplər:** array-lərin sərtliyi (ölçü tipin hissəsi), slice-lər (len/cap/append/
   make), slice paylaşım tələləri və 3-hissəli ifadə, string=rune+bytes+UTF-8, map (hash),
   comma-ok, map-i set kimi, struct-lar (anonim daxil).
4. **Bloklar və idarə:** shadowing (universe block daxil!), `if`-in scoped bəyanatı, `for`-un
   4 forması, for-range (map təsadüfi, string rune üzrə, dəyər=kopya), switch fall-through
   yoxdur, blank switch, `goto`-nun məhdud qanuni yeri.
5. **Funksiyalar:** variadic, çoxlu qaytarma (error sonda), named return (blank return
   QADAĞAN), funksiya=dəyər, closure-lar (sort.Slice, middleware), `defer` (LIFO, adlı
   return + commit/rollback pattern), call-by-value (map dəyişir, slice len-i dəyişməz).
6. **Pointer-lər:** &/*, nil dereference panic, class davranışı = pointer davranışı,
   pointer = mutable parametr siqnalı, nil pointer-i "diriltmək" olmaz, map=runtime pointer,
   slice=3 sahəli struct, buffer pattern, stack/heap/escape analysis, mechanical sympathy.
7. **Tiplər, metodlar, interfeyslər:** receiver qaydaları (dəyişir/nil → pointer), nil
   receiver üçün kodlaşdırma (ağac), method value/expression, embedding=kompozisiya (dynamic
   dispatch YOX), implicit interface = type-safe duck typing, "accept interfaces return
   structs", interface nil tələsi, type assertion/switch (comma-ok mütləq), DI frameworksüz.
8. **Xətalar:** error interface, sentinel (nadir), custom tiplər (nil interface tələsi!),
   `%w` wrap zənciri, `errors.Is/As` (== və assertion YOX), defer ilə ortaq wrap, panic/recover
   (library sərhəddi), stack trace paketləri.
9. **Modullar:** repo→modul→paket, böyük hərf=export, paket adı=funksional ad (`util` YOX),
   cmd/pkg struktur, internal, init-dən qaçma, cycle həlli, type alias, SemVer, minimum
   version selection, /v2 qaydası, proxy + sum database təhlükəsizliyi.
10. **Paralellik:** CSP modeli, goroutine, channel cədvəli (nil/ bağlı / buffered), select
    (random → starvation yoxdur), loop dəyişəni tələsi, done channel, cancel funksiyası,
    backpressure (token buffer), nil channel ilə case söndürmə, WaitGroup (çoxyazanlı
    close), sync.Once, mutex vs channel qərar ağacı, RWMutex, sync.Map məhdudiyyəti.
11. **Standart kitabxana:** io.Reader fəlsəfəsi (caller-buffer, EOF-ə qədər data), decorator,
    time (Duration, monotonic, 2006 şablonu), encoding/json (struct tag, Encoder/Decoder,
    custom Marshaler, iki strukturlu ayrılıq), net/http (timeout mütləq, ServeMux, middleware
    closure zənciri).
12. **Context:** ilk parametr konvensiyası, WithCancel/WithTimeout (defer cancel mütləq),
    Done/Err, uşaq-valideyn deadline büdcəsi, WithValue (yalnız unexported açar, middleware),
    GUID tracking pattern.
13. **Testlər:** _test.go, Error vs Fatal, TestMain, t.Cleanup, testdata, packagename_test
    (black-box), go-cmp, table testlər, coverage (100% bug YOXDUR demək deyil), benchmark
    (b.N, blackhole, -benchmem), stub patternləri (embed, funksiya sahəli), httptest, build
    tag ilə integration, -short əleyhinə arqumentlər, race checker.
14. **Reflect/unsafe/cgo:** reflection 3 anlayışı (Type/Kind/Value), marshaler quruluşu,
    MakeFunc, Filter benchmarkı (30-70x yavaş — sərhəd xarici QADAĞAN), unsafe binary
    konversiya (endianness, KeepAlive), cgo (pointer qaydaları, 29x yavaş — yalnız əvəzsiz
    C kitabxana).
15. **Generics:** `[T any]`, `comparable`, generic interfeys (Orderable[T]), type list
    (operatorlar), Map/Reduce/Filter, nə YOXDUR (operator overloading, parametrli metodlar,
    variadic tiplər), idiomatik təsir.

Ən vacib 5 fikir
1. Aydınlıq qısalıqdan üstündür: açıq tip çevirmələri, açıq error yoxlaması, açıq
   concurrency — Go-nun bütün "verbositesi" oxunaqlıq üçün seçilib.
2. Implicit interface + kiçik client-side interfeyslər = decoupling + dependency
   injection-in Go təbii forması.
3. Slice-in 3 sahəli strukturu və paylaşım semantikası — ən çox bug mənbəyi; 3-hissəli
   slice ifadəsi və copy ilə müdafiə.
4. Error-lar dəyərdir: yoxlanılır, wrap olunur (%w), zəncirdə axtarılır (Is/As) — exception
   yox, şəffaf data axını.
5. Concurrency struktur alətidir: done/cancel ilə hər goroutine-in çıxış yolu + channel
   data axını üçün, mutex isə struct sahə qoruması üçün.

## Kitabın ən dəyərli hissəsi

Chapter 7 (interfeyslər) və Chapter 10 (concurrency) — Go-nun digər dillərdən fundamental
fərqləndiyi, ən çox yanlış başa düşülən iki mövzu. Bu iki fəsli anlayan kəs Go-nu "düşünə"
başlayır.
