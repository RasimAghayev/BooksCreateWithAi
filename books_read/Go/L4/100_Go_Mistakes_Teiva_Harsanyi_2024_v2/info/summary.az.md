# 100 Go Mistakes and How to Avoid Them (2nd ed.) — Xülasə (AZ)

## Kitab kimin üçündür?

Go bilikli developer-lər üçün — "Effective Java-nın Go ekvivalenti" (Nutan
rəyi). Kitabın tezisi: beyin səhvlərdən öyrənir (neyroelm 2011); 100 real
səhv — bug, needless complexity, readability, suboptimal organization, API
convenience, optimization, productivity kateqoriyalarında.

## 100 səhv qruplaşması (12 fəsil)

1. **Kod təşkilatı (#1-16):** shadowing, nested code, init, interface
   pollution, producer-side interfeys, any, generics vaxtı, embedding,
   functional options, util paketlər, dokumentasiya, linters
2. **Məlumat tipləri (#17-29):** overflow, float, slice len/cap, nil vs
   empty, append yan təsiri, slice/map leak, müqayisə
3. **Nəzarət (#30-38):** range kopyası, ifadə bir dəfə, pointer element,
   map sırasız, break switch, defer loop
4. **Stringlər (#39-43):** rune/bayt, Builder, konversiya, substring leak
5. **Funksiyalar (#44-47):** receiver seçimi, named results, nil receiver,
   io.Reader API, defer qiymətləndirmə
6. **Xətalar (#48-54):** panic, wrap, Is/As, iki dəfə handle, defer xətası
7. **Paralellik əsasları (#55-60):** concurrency≠parallelism, G-M-P, kanal vs
   mutex, memory model, workload, context
8. **Paralellik praktikası (#61-74):** ctx detach, goroutine lifecycle, loop
   var, select determinizm, nil kanal, kanal ölçüsü, String deadlock, append
   race, WaitGroup, Cond, errgroup, sync kopyası
9. **Standart kitabxana (#75-82):** Duration, time.After leak, JSON embedded
   + monotonic, SQL, transient resurs, handler return, DefaultClient
10. **Test (#83-90):** kateqoriyalar, -race, parallel/shuffle, TDT, sleep-siz,
    time injeksiya, httptest/iotest, benchmark tələləri
11. **Optimallaşdırma (#91-100):** CPU cache, false sharing, ILP, alignment,
    stack/heap, escape, sync.Pool, inlining, pprof, GC, CFS

## Ən vacib 5 fikir

1. **Sadəlik ≠ asanlıq** — Go-nun 25 açar sözü aldadıcı; mənimsəmək nuansları
   hər səhvdə gizlidir
2. **Range = kopya, map = qaydasız, select = random** — bu üç "görünməz"
   davranış 100 səhvin böyük hissəsinin köküdür
3. **Paralellik struktur, parallelizm icradır** — goroutine yaratmaq ucuz,
   amma lifecycle, race və scheduler biliyi olmadan = bug fabriki
4. **Yaddaş sızması statistik deyil, strukturdur:** slice backing paylaşımı,
   capacity leak, false sharing, escape — hamısı QARIDA görülə bilər
5. **Ölç → dəyiş:** benchmark + benchstat + pprof olmadan heç bir
   optimallaşdırma "hissi" etibarlı deyil

## Kitabın ən dəyərli hissəsi

Chapter 12 (Optimizations) — CPU cache mexanikasından CFS quota-ya qədər
"hardware-dən Kubernetes-ə" tam zəncir; 512 vs 513 sütun (50% fərq!) və false
sharing nümunələri unudulmaz dərslərdir.

## 2nd edition yenilikləri

- #9 (generics confusion), #73 (errgroup), #74 (sync kopyası) kimi yeni
  mövzular; Go 1.22 loop semantikası, 1.25 container-aware GOMAXPROCS
  qeydləri aktualdır
