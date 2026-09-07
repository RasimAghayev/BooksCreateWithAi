# Know Go — Xülasə (Azərbaycanca)

**Müəllif:** John Arundel | **Nəşriyyat:** Bitfield Consulting | **İl:** 2025 | **Səviyyə:** L3 (Intermediate)

## Kitabın ümumi məqsədi

"Know Go" — John Arundelin generics mövzusunda dərin, amma dostcasına yazılmış
kitabıdır. Kitab Go-nun generics mexanizmini SIFIRDAN — interfeyslərlə
müqayisədən başlayaraq — tam dərinliyə qədər aparır: type parameter,
constraint, type set, approximation (~), generic tiplər və onların metodları,
container-lərin (Set, Stack) tikilməsi, onlara concurrency safety əlavəsi,
slices/maps/cmp/iter kimi yeni standart paketlər və Go 1.23 iteratorları.
Hər fəsil "Nədir / Necə işləyir / Nəyə lazımdır" sualları və GOAL → HINT →
SOLUTION məşqləri ilə gedir; bütün kod nümunələri testlərlə təsdiqlənir.

Kitabın fərqi: sadə "necə yazılır" YOX — "niyə belə dizayn edilib" sualı.
Generics-in əsası olan compile-time instantiation ("stencilling"),
runtime-da generic tipin OLMAMASI, interface indirection ilə müqayisədə
səmərəlilik, ~ approximation-un derived tiplər üçün zərurəti — hamısı
səhv mesajları vasitəsilə təbii şəkildə öyrədilir.

## Fəsil-fəsil xülasə

1. **Interfaces:** specific vs generic programming, implicit implementation
   (bəyansız), io.Writer, metod dəsti müqaviləsi, polymorphism, interfeysin
   məhdudiyyəti (yalnız metod adı+imzası, davranış YOX), any-də `x + y`
   xətası, type assertion/switch dublikat problemi, math paketinin
   float64-only problemi, generics tarixi (2009 → Ian Lance Taylor → ~10 il).
2. **Type parameters:** T placeholder, parameterised funksiya, instantiation
   (Go T-ni çıxarır), explicit instantiation (nəticə T olanda), stencilling
   (hər tip üçün ayrıca maşın kodu, indirection YOX), Identity (any-versiya
   ilə müqayisə), composite generic tiplər (Bunch[E]), generic funksiya
   tipləri (type idFunc[T any] func(T) T), + operatorun any-də QADAĞASI.
3. **Constraints:** basic interface (metod dəsti), any-nin zəifliyi ("bigger
   the interface, weaker the abstraction"), type element, union (|),
   kompozisiya (Number = Integer | Float | Complex), intersection,
   empty type set xətası, struct literal constraint, sahə çıxışı limitasiyası,
   "constraints are not classes" (Farm[Animal] YOX), ~ approximation
   (MyInt, Point), interface literal (T ~int ixtisarı), constraint daxilində
   T-yə istinad (Equal(T) bool — yalnız literal ilə mümkün).
4. **Operations:** operator → constraint cədvəli (+ → Number, > →
   cmp.Ordered, == → comparable), Real (Complex çıxarılıb), max/min
   built-in-ləri, T və U ayrı tiplərdir (operator aralarında İŞLƏMİZ),
   derived slice problemi ([]string ≠ StringList) → [S ~[]E, E any] idiomu
   (slices.Clone imzası kimi), comparable-in sonsuz tip dəsti (predeclared,
   compiler daxili), abstract type (body-də var top E), E(0) zero value,
   type switch qadağası → any(x).(type) həlli.
5. **Types:** named basic tiplər (Age/HeightCM — məna mismatch), generic
   basic tip QADAĞASI, Bunch[E] homogen vs []any heterogen ("runtime-da
   generic tip YOXDUR"), generic map (K comparable — açar == tələb edir),
   generic struct (NestedThing — daxili instantiate), metodlar receiver-da
   E saxlayır (First()), parameterised metod QADAĞASI (PrintWith[T] YOX),
   generic interface (Equaler[T]), generic channel (myChan[E]).
6. **Functions:** Contains (comparable), Reverse (any), Sort (cmp.Ordered),
   first-class functions + generics: mapFunc/keepFunc/reduceFunc tip adları,
   Map (type inference dərsi: mapFunc[int] mismatch), Filter (generic
   keepFunc: IsEven[int] explicit instantiate, constraints.Integer), Reduce
   (any kifayətdir — operator BİZİM funksiyamızdadır), konkurrent
   Map/Filter (embarrassingly parallel), FuncMap (dynamic dispatch),
   Compose (f(g(v)), 3 type parameter).
7. **Containers:** Set = map[E]struct{} (bool YOX — zero-size), variadic
   NewSet/Add (inference sayəsində instantiate lüzumsuz), All/String,
   Union/Intersection (dublikat-təmizləmə PULSUZ), işə qəbul nümunəsi,
   Stack (LIFO, data []E, Pop (v, ok) pattern, pointer receiver).
8. **Concurrency:** data race, 3 müdafiə yolu (shared data yox > guard
   goroutine+channel > mutex), mutex problemləri (contention, deadlock,
   opt-in), SetC (mutex *sync.RWMutex — pointer, kopya qorunması), Lock vs
   RLock, 1000-lik smoke test, mutexsiz "concurrent map write" fatal,
   -race detector, Intersection optimizasiyası (bir RLock vs granulyar),
   atomic.Uint64 (mutex paperwork-siz sayaç), Channel məşqi.
9. **Packages:** cmp.Or, slices — Equal/EqualFunc (2 fərqli tip mümkün!),
   Compare, Index/Contains, Max/Min, Insert/Delete, Clone, Compact
   (yalnız ARDICIL dublikatlar), Grow/Clip, Sort/SortStableFunc/IsSorted,
   Reverse, Replace, BinarySearch (tapılmasa insertion point), maps —
   Equal/DeleteFunc/Clear/Clone/Copy(to, from!), yeni idiomlar (Clone
   əvəzi 3 sətirlik make+copy), Merge məşqi ([M ~map[K]V] + maps.Copy).
10. **Questions:** nə qədər bilmək lazım (çox az), kod tərzi dəyişmir,
    mövcud kod dəyişməz ("if it isn't broken"), abstraksiya xərci
    (Griesemer: "move cautiously"), performans (runtime EYNY; interface/
    reflection əvəzsə 2x SÜRƏTLİ — tidwall/btree; compile fərqsiz),
    niyə [] (<> parser ambiguity), olmayanlar (option, enum, union, makros,
    parameterised metod, generic paket, covariance), "Go 2 olmayacaq"
    (breaking change = fəlakət: Python/IPv6 dərsləri), amma varis dil olacaq.
11. **Iterators:** Go 1.23 "range over func": slice qaytarmağın 3 problemi
    (gözləmə/yaddaş/israf), func(yield func(V) bool) imzası, yield false =
    təmizlənmə (davam etsə panic), iter.Seq/Seq2, funksiya iterator QAYTARIR,
    (value, error) Seq2 pattern, kompozisiya (Primes(Integers()) — sonsuz
    ardıcıllıq + filtrlər), iterator vs channel (concurrenciyyə məcbur
    etmə, goroutine leak), slices.All/Values/Collect, maps.All/Keys/Values.

## Kitabın əsas mesajları

1. Generics = compile-time fenomen: runtime-da generic tip YOXDUR, hər
   şey konkret tipə instantiate olunur.
2. Constraint = zəmanət mübadiləsi: operator istəyirsiniz → tip dəstini
   daraldın (Number/Ordered/comparable). "Bigger the interface, weaker
   the abstraction."
3. Derived tipləri heç vaxt unutma: ~ approximation, [S ~[]E, E any],
   [M ~map[K]V] idiomları.
4. Generic container + concurrency safety BİR DƏFƏ yazılır — hamısı üçün
   işləyir; amma abstraksiya xərci var — konkret fayda olmadan yazma.
