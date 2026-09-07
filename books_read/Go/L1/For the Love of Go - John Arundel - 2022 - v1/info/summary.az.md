# For the Love of Go — Xülasə (Azərbaycanca)

**Müəllif:** John Arundel | **Nəşriyyat:** Bitfield Consulting | **İl:** 2022 | **Səviyyə:** L1 (Beginner)

## Kitabın ümumi məqsədi

John Arundelin "For the Love of Go" — tam başlanğıc səviyyəsindən Go öyrənmək
üçün interaktiv "ssenarili" kitabdır. Oxucu Texio Instronics (kalkulyator) və
Happy Fun Books (onlayn kitab mağazası) şirkətlərində "ilk iş günü"ndən başlayaraq
real layihələr üzərində işləyir. Kitabın əsas fərqi: hər mövzu TEST YAZMAQLA
öyrədilir — 1-ci fəsildən etibarən TDD tsikli (test → compile error → null
implementation → FAIL → real implementation → PASS), want/got pattern, table
testlər, valid/invalid davranış ayrılığı. İki layihə təbii təkamül yolu ilə
dilin bütün əsaslarını açır: modullar, funksiyalar, tiplər, error-lar, struct-lar,
slice/map kolleksiyaları, metodlar, pointer-lər, validasiya (unexported +
accessor), konstruktorlar, map-as-set, konstantlar/iota, control flow (if/switch/
for/range/continue/break), closure-lar, defer, variadic, binary compilation
(GOOS/GOARCH cross-compilation). Final fəsil — "The Tao of Go": kindness,
simplicity, humility, not striving fəlsəfəsi.

Kitabın tərzi: hər fəsildə GOAL məşqləri (özün yaz), compiler xətalarının
"tərcüməsi", null implementation dərsi ("testin bug tutduğunu SÜBUT et"),
"One behaviour, one test" qanunu, "What are we really testing here?" sualı.

## Fəsil-fəsil xülasə

1. **Testing times:** go mod init, calculator layihəsi, go test, gofmt,
   funksiya sintaksisi, FAIL çıxışının oxunması, Subtract bug-ı, testing paketi,
   test imzası (_test.go, Test prefiks, *testing.T), want-and-got, if/bool.
2. **Go forth and multiply:** Multiply TDD ilə (kopyala-dəyiştir), null impl,
   table testlər (testCase struct + range), Divide invalid input problemi,
   "something and error" pattern, "One behaviour, one test".
3. **Errors and expectations:** (float64, error) imza, compiler xətaları
   (assignment mismatch, not enough arguments), errors.New, TestDivideInvalid,
   blank identifier, zero value konvensiyası, red-green-refactor, closeEnough
   (float tolerance), go run/build, cmd/ strukturu.
4. **Happy Fun Books:** string/int/bool tipləri, value/variable/type, zero
   values, var vs `:=`, struct, type definition, doc comments, exported/
   unexported, core package bookstore, compile-only test (`_ = bookstore.Book{}`).
5. **Story time:** user stories, core stories (buy/list/details), davranış-
   əsaslı dizayn, dot notation, ++/--/+=/-=, TestBuy, null impl (Book{}),
   coverage alətləri, test-last təhlükəsi (zero copies → -1), error-lu Buy,
   covered vs tested (TestBuyTrivial), 80-90% hədəf.
6. **Slicing & dicing:** slice ([]Book, literal, index, len, append), go-cmp
   (cmp.Equal/Diff, go get -t), GetAllBooks, unikal ID problemi, GetBook
   range-axtarışı, "Crime doesn't pay" (catalog[0] tələsi → 2 elementli test).
7. **Map mischief:** map[int]Book, lookup/assign, sahə yazışı qadağası
   (çıxar-dəyiş-geri yaz), zero value lookup, comma-ok (b, ok :=), error-lu
   GetBook (fmt.Errorf), GetAllBooks map-dən slice, RANDOM map sırası →
   sort.Slice + function literal, flaky test həlli.
8. **Objects behaving badly:** obyekt = data + davranış, PriceCents (int sent),
   metod/receiver, funksiya→metod, non-local tip qadağası, type MyInt int,
   type conversion, MyString.Len, type Catalog map[int]Book, metodlaşdırma.
9. **Wrapper's delight:** strings.Builder, metod MİRASI YOXDUR, wrapping
   (struct + Contents sahəsi), StringUppercaser, pass-by-value (Double
   puzzle: want 24, got 12), pointer (& sharing), *int ≠ int, dereference (*),
   nil pointer panic, pointer receiver (`func (input *MyInt) Double()`).
10. **Very valid values:** SetPriceCents value receiver tələsi (SA4005),
    automatic dereferencing, always valid field (unexported category +
    Category/SetCategory), cmpopts.IgnoreUnexported, always valid struct
    (creditcard: unexported card + New constructor), map-as-set (validCategory),
    constants (http.StatusOK), type Category int, iota, table-driven kateqoriyalar.
11. **Opening statements:** statement/declaration/assignment, `:=`, tuple
    assignment, blank identifier, ++ statement (ifadə DEYİL), happy path
    refaktoru (flip + early return), else nadir, &&/||/!, inclusive or, if ok,
    compound if (shadowing təhlükəsi, uzunluq qaydası).
12. **Switch which?:** switch/case/default, ilk-match, fallthrough, switch
    expression + comma case, break, for (conditional/forever/range), 3-hissəli
    for, continue (happy path sola), break, labels (nadir — refaktor siqnalı),
    goto (avoid).
13. **Fun with functions:** signature, parameter/result listlər, çağırış =
    control flow, functions are values (TestCase.function sahəsi, sort.Slice
    less), function literal, closures (bubble), loop variable tələsi (3,3,3!),
    defer (resource leak, LIFO stacking), named results (documentation),
    naked returns harmful, deferred closure (err = closeErr — yalnız fail-da),
    variadic (...float64 → slice, AddMany/DivideMany).
14. **Building blocks:** executable binary, package main + main(), init
    (magic — ÜSTÜNLÜK VERİLMEZ; əvəzinə main-də et / var x = initFn()),
    go build (file: Mach-O/PE32+/ELF, ~2 MiB asılılıqsız), GOOS "goose"/
    GOARCH "gorch" cross-compilation, go tool dist list, os.Exit(0/1
    konvensiyası), "yalnız main-dən exit et" qaydası.
15. **The Tao of Go:** Tao = təbiətlə mübarizəsizlik ("thrashing yox, surfing");
    Kindness (istifadəçi/İşlədən/oxuyan/özün üçün kod, mikro-refaktorlar),
    Simplicity ("Simplify, simplify!", bir işi yaxşı, extensibility tələsi),
    Humility (clever YOX, pre-engineering YOX, öz kodunu review, sənəd+nümunələr),
    Not Striving (problem-eliminating, "ən yaxşı optimallaşdırma = etməmək",
    programming ≠ typing, DELETE düyməsi, buffalo qanunu).

## Ən vacib 5 fikir

1. **Test-first öyrənmə:** hər konsepsiya testlə təqdim olunur — "We simply
   can't write untestable functions when the test comes first"; null
   implementation testin bug tutduğunu SÜBUT edir.
2. **"One behaviour, one test":** valid və invalid davranışlar ayrı testlər;
   funksiya yox, DAVRANIŞ test olunur; "What are we really testing here?"
   hər dizayn qərarının açarıdır.
3. **Always valid data:** unexported sahə/tip + accessor/constructor —
   "Making it impossible to compile incorrect programs is the best kind of
   validation!"
4. **Pointer dərsləri:** pass-by-value (kopya!), pointer receiver (setter
   mütləq), automatic dereferencing, closure-un çağırış-vaxtı semantikası
   (3,3,3 tələsi).
5. **Tao fəlsəfəsi:** kindness (insanlar üçün kod), simplicity (bir işi yaxşı,
   defaults), humility (óbvio > clever), not striving (problemi ləğv et —
   həll etmə; "ən yaxşı optimallaşdırma — işi etməmək").

## Kitabın auditoriyası

- Proqramlaşdırmaya YENİ başlayanlar (Go ilk dil kimi) — ən uyğun
- Başqa dildən gələnlər üçün "Go-nun yolu" — test mədəniyyəti ilə
- TDD-ni praktikada görmək istəyənlər — 15 fəsil hamısı test-əsaslı
- "The Power of Go" seriyasına giriş pilləsi (Tests/Tools/Generics əvvəli)
