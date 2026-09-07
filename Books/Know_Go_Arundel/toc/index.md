# Know Go — Mündəricat və Naviqasiya

**Müəllif:** John Arundel | **Nəşriyyat:** Bitfield Consulting | **İl:** 2025 | **Səviyyə:** L3 (Intermediate)

Kitabın mövzusu: Go generics — tam dərinlikdə. Type parameter-lərdən
iteratorlara qədər: interfeyslərlə müqayisə, constraint-lər, type set-lər,
generic tiplər və funksiyalar, container-lər (Set/Stack), concurrency safety,
slices/maps/cmp/iter paketləri, performans və ekosistem sualları.

## Fəsillər

| # | Fəsil | Səhifə | Əsas mövzu |
|---|-------|--------|------------|
| 1 | [Interfaces](../chapters/01-interfaces/index.md) | 14-24 | İnterfeyslər: implicit implementation, metod dəstləri, any-nin limitləri, type assertion/switch, generics tarixi |
| 2 | [Type parameters](../chapters/02-type-parameters/index.md) | 25-43 | T placeholder, instantiation (implicit/explicit), stencilling, generic funksiya tipləri |
| 3 | [Constraints](../chapters/03-constraints/index.md) | 44-65 | Basic interface, type set, union/intersection, ~ approximation, interface literal-lər, Equal(T) pattern |
| 4 | [Operations](../chapters/04-operations/index.md) | 66-90 | Number/Ordered/Real, cmp.Ordered, comparable, [S ~[]E, E any] idiomu, abstract type, zero value |
| 5 | [Types](../chapters/05-types/index.md) | 91-103 | Generic tiplər: slice/map/struct/interface/channel; metodlar; parameterised metod QADAĞASI |
| 6 | [Functions](../chapters/06-functions/index.md) | 104-124 | Contains/Reverse/Sort, Map/Filter/Reduce, FuncMap, Compose |
| 7 | [Containers](../chapters/07-containers/index.md) | 125-138 | Set (map[E]struct{}), Add/Contains/Union/Intersection, Stack (LIFO) |
| 8 | [Concurrency](../chapters/08-concurrency/index.md) | 139-155 | Data race, RWMutex, SetC, race detector (-race), sync/atomic |
| 9 | [Packages](../chapters/09-packages/index.md) | 156-176 | cmp, slices (Equal/Sort/BinarySearch/...), maps (Equal/Copy/...) |
| 10 | [Questions](../chapters/10-questions/index.md) | 177-191 | Performans, sintaksis (<> vs []), olmayan xüsusiyyətlər, "Go 2" sualı |
| 11 | [Iterators](../chapters/11-iterators/index.md) | 192-209 | iter.Seq/Seq2, yield protokolu, kompozisiya, iterator vs channel |

## Oxu ardıcıllığı tövsiyəsi

1-3: bünövrə (interfeys → type parameter → constraint) — ardıcıl oxunmalı.
4-6: əməliyyatlar və funksional üslub. 7-8: praktik container-lər və
concurrency. 9: standart kitabxana. 10-11: kontekst və yeni istiqamət.

## Kitabdakı məşqlər (GOAL/HINT/SOLUTION formatında)

- Ch2: PrintAnythingTo, Group
- Ch3: stringy (StringifyTo), intish (Intish ~int), greater (Equal(T)-variand interface literal)
- Ch4: product (Product[T Number]), dupes (Dupes[E comparable] map-idiomu)
- Ch5: empty (Sequence[E].Empty())
- Ch6: funcmap (FuncMap[T, U] dynamic dispatch), compose (Compose[T, U, V])
- Ch7: stack (Stack[E] LIFO)
- Ch8: channel (instrumented Channel + atomic sayaclar)
- Ch9: merge (Merge[M ~map[K]V])
