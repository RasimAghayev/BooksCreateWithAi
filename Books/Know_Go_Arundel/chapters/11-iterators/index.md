# Chapter 11 — Iterators (İteratorlar)

## Bu fəsil nədən bəhs edir?

Go 1.23-də gələn iteratorlar ("range over func"): nə üçün lazımdırlar
(slice qaytaran funksiyaların 3 problemi: bütün slice hazırlanana qədər
gözləmə, tam slice üçün yaddaş, istifadə olunmayan elementlərin israfı),
iterator imzası (`func(yield func(V) bool)`), yield-in true/false protokolu
(false = break/return baş verdi → təmizlənmə şansı; false-dan sonra davam
→ runtime panic), iter.Seq / iter.Seq2 (index+value, value+error), funksiyalar
iterator QAYTARIR (özü iterator deyil), kompozisiya (Primes(Integers()) —
sonsuz ardıcıllıqlar), iterator vs channel (concurrenciyyə məcbur etməmək,
goroutine leak), slices.All/Values/Collect, maps.All/Keys/Values. Sonda
müəllif haqqında və digər kitablar (bu hissə biliyə aiddir, texniki deyil).

## Əsas fikirlər

### 1. Niyə Iterator — Slice Qaytarmağın 3 Problemi
```go
func Items() (items []Item) {
    ... // BÜTÜN items-ləri GENERATE ET, sonra qaytar
}
for _, v := range Items() { ... }
```
- **Problem 1:** ilk elementi emal etmək üçün BÜTÜN slice-in hazırlanmasını
  gözləmək lazım
- **Problem 2:** bir elementi eyni anda istifadə etsək də, tam slice üçün
  yaddaş ayrılır
- **Problem 3:** elementlərin hamısını istifadə etməsək — hesablanan vaxt
  və yaddaş İSRAF olunur
- **Həll:** iterator = tələb olduqca bir element "yield" edən funksiya
  (qısa-sifariş aşpazı metaforası: hamısını əvvəlcədən bişirməz, müştəri
  gəldikcə bişirir)

### 2. İteratorun İmzası və yield Protokolu
```go
func iterateItems(yield func(Item) bool) {
    items := []Item{1, 2, 3}
    for _, v := range items {
        if !yield(v) {
            return      // loop bitdi — təmizlənmə şansı
        }
    }
}   // normal qayıdış = "iteration complete"

for v := range iterateItems {
    ...
}
```
- **Nədir:** iterator = xüsusi imzalı funksiya: `func(yield func(V) bool)`
- **yield true:** "daha da ver" — loop davam edir
- **yield false:** loop break/return ilə çıxdı — iterator DAVAM ETMƏLİDİR:
  - detect edib qayıt (lazımsız hesabdan qaçınmaq)
  - təmizlənmə (database bağlantısı kimi resursların buraxılması) —
    yield false = "Clean up, you're done" siqnalı
- **PANIK halı:** false-dan sonra yenidən yield →
  `range function continued iteration after exit`

### 3. iter.Seq və iter.Seq2 — Standart Tiplər
```go
import "iter"

func Items() iter.Seq[Item] {          // 1 dəyər
    return func(yield func(Item) bool) {
        items := []Item{1, 2, 3}
        for _, v := range items {
            if !yield(v) { return }
        }
    }
}

for v := range Items() { fmt.Println("item", v) }
```
- **Vacib:** Items ÖZÜ iterator DEYİL — o, iterator QAYTARIR; range
  ifadəsi funksiya çağırışının NƏTİCƏSİ üzərindədir
- **Seq2 — iki dəyər:**
```go
func Items() iter.Seq2[int, Item] {   // indeks + dəyər
    return func(yield func(int, Item) bool) {
        items := []Item{1, 2, 3}
        for i, v := range items {
            if !yield(i, v) { return }
        }
    }
}
for i, v := range Items() { ... }   // slice-range ilə eyni görüntü
for _, v := range Items() { ... }   // indeks istenmirsə blank
for range Items() { ... }           // heç biri istenmirsə
```
- **Uyğunsuzluq xətası:** Seq ilə 2 dəyişən →
  `permits only one iteration variable`

### 4. Xətaların İdarəsi — Seq2[value, error]
```go
func Lines(file string) iter.Seq2[string, error]
```
- **Problemlər:** iterator səbəbsiz dayansa çağıran loop bunu BİLMİR
- **Həll:** yield (value, error) — uğurda error=nil; xətada error=dəyər;
  loop daxilində error yoxlaması adi `(T, error)` funksiyası kimidir
- **İnfallible iteratorlar:** kolleksiya tipləri (Set.All) — xəta
  qaynağı yoxdur; Seq kifayətdir

### 5. Kompozisiya — İteratordan İterator
```go
func Integers() iter.Seq[int] {
    return func(yield func(int) bool) {
        for i := range math.MaxInt {        // SONSUZ ardıcıllıq!
            if !yield(i) { return }
        }
    }
}

func Primes(seq iter.Seq[int]) iter.Seq[int] {
    return func(yield func(int) bool) {
        for n := range seq {                // iteratoru range ET
            if isPrime(n) {
                if !yield(n) { return }
            }
        }
    }
}

for p := range Primes(Integers()) { ... }   // zəncir: sonsuz → sadələr
```
- **Güc:** sonsuz məlumat mənbəyi + filtrləyici zəncirlər — yalnız
  istənilən qədər hesablanır (lazy evaluation)
- **PrintAll[V any](seq iter.Seq[V]):** iteratorun NƏYƏ iteratoru olması
  önəmsizdir — generic + range universal

### 6. İterator vs Channel
| | İterator | Channel |
|---|---|---|
| Konkurentlik | YOX — sadə ardıcıllıq | Bəli — goroutine lazımdır |
| Erken çıxış (break) | yield false → təmiz bitiş | göndərən goroutine BLOKLANIR → LEAK |
| Xəta riski | aşağı | deadlock/crash ehtimalları |
| Nə zaman | ardıcıllıq kifayətdirsə | proqram onsuz da konkurrentdirsə |

- **Channel həllinin 2 problemi:** (1) proqram konkurent olur — "concurrent
  programming is just hard"; (2) loop erkən çıxarsa göndərən goroutine
  bloklanır və yaddaşda asılı qalır (leak) — context ləğvi kimi vasitələr
  var, amma "annoying"
- **Qərar:** ardıcıllıq lazımdırsa İTERATOR; konkurrentli arcsa CHANNEL

### 7. Standart Kitabxana Yenilikləri
```go
// slices:
for i, v := range slices.All(s) { ... }    // (index, value) iteratoru
for v := range slices.Values(s) { ... }     // yalnız dəyərlər
s := slices.Collect(slices.Values([]int{1, 2, 3}))  // [1 2 3]
// Collect: iteratorun NƏTİCƏLƏRİNİ slice-a yığır (geri dönüş)

// maps:
for k, v := range maps.All(m) { ... }      // (açar, dəyər)
for k := range maps.Keys(m) { ... }         // açarlar
for v := range maps.Values(m) { ... }       // dəyərlər
// xəbərdarlıq: sıralanmamış (unordered) — map təbiəti
```
- **Collect əhəmiyyəti:** iterator → slice — API-lər arasında körpü

## Əsas terminlər
- Iterator (itorator) — tələbə görə element verən funksiya
- Yield (vermə) — növbəti dəyərin loop-a ötürülməsi
- iter.Seq / iter.Seq2 — standart iterator tipləri (1 və 2 dəyər)
- Lazy evaluation (tənbəl hesablama) — yalnız istənilən qədər hesabla
- Goroutine leak (goroutine sızması) — bloklanmış goroutine-in yaddaşda qalması
- Infallible (yanılmaz) — xətər qaynağı olmayan iterator
- Composition (kompozisiya) — iteratorların zəncirlənməsi

## Praktik nəticə

1. **API dizaynı:** kolleksiya qaytaran funksiya YERİNƏ iterator qaytar —
   `iter.Seq[T]` (1 dəyər) / `iter.Seq2[K, V]` (2 dəyər; indeks və ya error).
2. **Iterator yazarkən:** hər yield-dən sonra `if !yield(...) { return }` —
   həm performans, həm təmizlənmə, həm panic qarşısı.
3. **Sonsuz/iri ardıcıllıqlar:** iterator + kompozisiya — slices.Collect
   lazım olanda slice-a yığır.
4. **Channel YALNIZ konkurrentlik tələb edildikdə:** ardıcıllıq üçün
   iterator daha təhlükəsizdir (leak riski yoxdur).
5. **Mövcud kodda range slice/map:** slices.All/Values, maps.All/Keys/Values
   ilə eyni sintaksis — keçid məqsədli.

## Mənbə
Pages: 192-209 (PDF 193-210; kitabın texniki hissəsi ~səh. 207-də bitir,
sonrası müəllif/barrier/acknowledgement bölməsidir)
