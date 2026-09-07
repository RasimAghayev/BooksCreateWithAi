# Chapter 9 — Packages (Paketlər)

## Bu fəsil nədən bəhs edir?

Generics ilə gələn YENİ standart kitabxana paketləri və ekosistem
dəyişiklikləri: `cmp` (Ordered constraint + Or funksiyası), `slices`
(Equal/EqualFunc, Compare/CompareFunc, Index/Contains, Max/Min, Insert/Delete,
Clone, Compact, Grow/Clip, Sort/SortFunc/SortStableFunc/IsSorted, Reverse,
Replace, BinarySearch), `maps` (Equal, DeleteFunc, Clear, Clone, Copy),
köhnə əl ilə yazılan idiomların yeni birxətliq ekvivalentləri, Merge məşqi
(`[M ~map[K]V, K comparable, V any]` — derived map tipləri üçün ~ idiomu
+ maps.Copy-dan istifadə).

## Əsas fikirlər

### 1. cmp Paketi
```go
x := cmp.Or(userX, "default X value")
```
- **cmp.Ordered:** bütün ordered built-in tiplər + derived (ch4-dən tanış)
- **cmp.Or:** arqumentlərdən İLK non-zero olanı qaytarır — default value
  üçün birxətliq həll

### 2. slices — Müqayisə Funksiyaları
```go
s1 := []int{1, 2, 3}
s2 := []int{1, 2, 3}
fmt.Println(slices.Equal(s1, s2))       // true — slice-larda == YOXDUR!

fmt.Println(slices.EqualFunc([]string{"a"}, []string{"A"}, strings.EqualFold))
// true — custom equality

fmt.Println(slices.EqualFunc([]int{0, 1, 2}, []rune{'a', 'b', 'c'},
    func(n int, r rune) bool { return n == int(r-'a') }))
// true — FƏRLİ element tipləri müqayisə olunur!
```
- **Equal:** == üçün (element comparable olmalıdır)
- **EqualFunc:** custom funksiya ilə — İKİ FƏRLİ slice tipi (E və T)
  müqayisə edilə bilər
- **Compare/CompareFunc:** -1/0/+1 — lekoqrafik sıra; qısa slice kiçik
  sayılır; ilk fərqlənən element qərar verir

### 3. slices — Axtarış
```go
s := []int{1, 2, 3}
slices.Index(s, 2)      // 1 — ilk rastgəlinin indeksi, yoxsa -1
slices.IndexFunc(s, negative)
slices.Contains(s, 1)   // true — bizim Contains-a ehtiyac QALMADI
slices.ContainsFunc(s, negative)
```

### 4. slices — Max/Min və Dəyişikliklər
```go
slices.Max([]int{1, 3, 2})  // 3 — bizim Greatest funksiyası əvəzinə
slices.Min(s)                // ən kiçik
slices.MaxFunc/MinFunc(s, customCmp)

s := []string{"a", "c"}
s = slices.Insert(s, 1, "b")        // [a b c] — indeksə daxil et
s = slices.Delete(s, 0, 2)          // [c] — [0, 2) aralığı sil (2 DAHIL DEYIL)
s = slices.Clone(s)                 // shallow kopya
slices.Compact([]int{1, 1, 1, 2, 3})     // [1 2 3] — ardıcıl dublikatlar (uniq)
slices.Compact([]int{1, 2, 1, 3, 1})      // [1 2 1 3 1] — QONŞU OLMAYANLAR SAXLANILIR
slices.CompactFunc(s, strings.EqualFold)
```
- **Compact diqqət:** yalnız ARDICIL dublikatlar — Unix uniq kimi;
  sıralanmamış slice-də bütün dublikatları silmək üçün Sort + Compact

### 5. slices — Yaddaş İdarəetməsi
```go
s := []int{1, 2}
fmt.Println(cap(s))          // 2
s = slices.Grow(s, 10)
fmt.Println(cap(s))          // 12 — bir dəfədə böyüt (copy overhead-ini azalt)

s = slices.Delete(s, 0, 100)
s = slices.Clip(s)           // kapasiteyi length-ə kəs — yaddaşı geri qaytar
```
- **Grow:** çoxsaylı append-in hər dəfə kopya yaratmasının qarşısı
- **Clip:** boş qalan slotları yaddaşa qaytarır

### 6. slices — Sıralama və Axtarış
```go
slices.Sort([]int{3, 1, 2})           // [1 2 3] — in-place
slices.SortFunc(s, cmp.Compare)        // custom
slices.SortStableFunc(s, cmp)          // stabil — bərabər elementlərin sırası qorunur
slices.IsSorted(s) / IsSortedFunc(s, f)

slices.Reverse(s)                       // in-place tərs çevir
s = slices.Replace(s, 1, 2, "bat", "bee", "bison")
// [1,2) aralığını yeni elementlərlə əvəz et

slices.BinarySearch([]int{1, 2, 3}, 2)    // 1 true
slices.BinarySearch([]int{1, 2, 3}, 999)  // 3 false — tapılmasa DAHİL OLACAĞI indeks
slices.BinarySearch([]string{"a", "c"}, "b")  // 1 false
```
- **BinarySearch:** sıralı slice-də O(log n) axtarış — lüğətdə söz axtarmaq
  metaforası (hər baxış axtarış sahəsini yarıya bölür)
- **Tapılmadıqda:** düzgün yerləşəcəyi indeks qaytarılır — insertion point

### 7. maps Paketi
```go
fmt.Println(maps.Equal(m1, m2))     // == ilə (açarlar həmişə comparable)
maps.EqualFunc(m1, m2, valuesEq)     // custom value müqayisəsi

maps.DeleteFunc(m, func(k string, v bool) bool { return !v })  // şərti sil
maps.Clear(m)                       // hamısını sil

m2 := maps.Clone(m1)                // shallow kopya
maps.Copy(to, from)                 // from-dan to-ya — mövcud açarlar ÜZƏRİNƏ yazılır
```
- **Equal:** açar üçün custom funksiya LAZIM DEYİL (map açarı həmişə ==)
- **Copy arqument sırası:** Copy(destination, source) — Copy(to, from)!
- delete() built-in qalır; DeleteFunc — funksiya true olan girişləri silir

### 8. Yeni İdiomlar — Köhnə Əl Əməyinin Sonu
```go
// KÖHNƏ                            // YENİ
b := make([]T, len(a)); copy(b, a)   b := slices.Clone(a)
s = append(s[:1], s[3:]...)          s = slices.Delete(s, 1, 3)
for _, v := range s { if v == 2 …}   slices.Contains(s, 2)
```

### 9. Merge Məşqi — Derived Map İdiomu
```go
func Merge[M ~map[K]V, K comparable, V any](ms ...M) M {
    result := M{}
    for _, m := range ms {
        maps.Copy[map[K]V](result, m)
    }
    return result
}
```
- **Tələblər:** istənilən sayda map, eyni tip; konflikt halında SONRA
  gələn qalib
- **`~map[K]V`** — derived map tiplərini (type menu map[int]string) də qəbul
  et — StringList idiomasının map versiyası
- **maps.Copy[map[K]V]** — açıq instantiate: result M, amma Copy-ya
  adlandırılmamış map[K]V lazımdır (derived → underlying dönüşü)
- **Fəlsəfə:** "standard library take the strain" — nəsə hazır varsa
  onu istifadə et

## Əsas terminlər
- Shallow copy (seyrin kopya) — yalnız struktur kopyalanır, elementlərin ÖZLƏRİ yox
- Binary search (ikili axtarış) — sıralı məlumatda yarıya-bölma axtarışı
- Stable sort (stabil sıralama) — bərabər elementlərin ilkin sırası qorunur
- Compact (sıxlşdırma) — ardıcıl dublikatların silinməsi (Unix uniq)
- Insertion point (yerləşdirmə nöqtəsi) — elementin düzüləcəyi indeks
- Capacity (tutum) — slice-in slot sayı; length — dolu slot sayı

## Praktik nəticə

1. **Slice əməliyyatlarından ƏVVƏL slices paketinə bax:** Contains, Index,
   Clone, Delete, Sort, Max, BinarySearch — 90% hal hazır var.
2. **Custom müqayisə lazımdırsa:** -Func variantları (EqualFunc,
   CompareFunc, SortFunc, ContainsFunc) — hətta fərqli element tipləri
   müqayisə olunur.
3. **Performans:** böyüyəcək slice üçün slices.Grow; boşalmış üçün Clip.
4. **Map köçürmə:** maps.Copy(to, from) — sıra qarışdırma!
5. **Generic funksiya yazarkən derived tipləri unutma:** `[M ~map[K]V,
   K comparable, V any]` — StringList dərsi map-lərə də aiddir.

## Mənbə
Pages: 156-176 (PDF 157-177)
