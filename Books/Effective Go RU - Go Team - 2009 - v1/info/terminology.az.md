# Effective Go (RU) — Terminologiya (Azərbaycanca)

> Sənəd boyu terminlər. Format: `English (Azərbaycanca qarşılıq)`.

## A

**Ad gizlətmə qaydası** — embedding-də xarici sahə/metod eyni adlı daxilini örtür.

**Assertion (type assertion)** — `value.(T)` — interfeysdən konkret tip/dəyər çıxarma; `v, ok` forması təhlükəsizdir.

## B

**Blank identifier (`_`)** — yazıla bilən /dev/null: dəyər/istifadəsiz import/iqnorun aşkar forması.

**Builtin (dilə daxil funksiya)** — append/new/make — kompilyator dəstəyi tələb edən (generic öncəsi).

**Buffered channel** — bufer dolana qədər göndərən bloklamır; semafor kimi istifadə.

## C

**Composite literal** — `&File{fd: fd}` — konstruktoru əvəz edən ifadə; lokalin adresi qanuni.

**Comma-ok idiomu** — `v, ok = m[k]` / `v, ok = value.(T)` — uğursuzluğu ayıran ikili mənimsətmə.

**Concurrensy vs parallelism** — struktur (müstəqil komponentlər) / çoxcore-lu hesab icrası.

**CSP** — Communicating Sequential Processes: Go kanallarının nəzəri kökü.

**Canonical signature (kanonik imza)** — Read/Write/String-in standart semantikası; eyni adda metod uyğun olmalıdır.

## D

**Defer** — əhatəedən funksiyanın return-unda icra; arqumentlər DƏRHAL qiymətlənir; LIFO.

**Doc comment** — yuxarı səviyyə elandan əvvəlki comment — paket sənədi.

**DoSome/DoAll pattern** — Vector paralelləşdirməsi: hissələr + c kanalı sayğacı.

## E

**Embedding (daxilə yerləşdirmə)** — interfeys/struct tiplərin daxilə adsız yerləşdirməsi; metod promote.

**error interfeysi** — `Error() string`; PathError kimi zəngin strukturlarla tətbiq olunur.

**Export (böyük hərf)** — ilk hərf böyükdürsə paketdən görünən ad.

## F

**Fallthrough (yoxluğu)** — Go switch-də avtomatik keçid YOXDUR; vergüllü çoxlu case var.

**Functional literal** — anonim funksiya; Go-da closure — dəyişənlər aktivlikcə yaşayır.

## G

**GOMAXPROCS(0)** — cari paralel core limiti (env üstünlüyünə hörmət edir).

**Goroutine** — yüngül paralel funksiya; stack heapdan böyüyür; OS thread-lərə multipleks.

**gofmt** — standart formatter — bütün format müzakirələrini bitirən alət.

## I

**init funksiyası** — paket ilkinləşdirməsi; sıra: import → var → init.

**iota** — konstant enumerator; ifadədə iştirak edir, implicit təkrar.

## K

**Kanal-kanal (channel of channels)** — Request içində resultChan — asinxron RPC.

## L

**Label break** — `break Loop` — nested switch-dən loop qırma.

**Leaky buffer** — freeList channel + select/default — GC idarəli buffer pool.

**LIFO** — defer-lərin icra sırası (son giren, ilk çıxan).

## M

**make vs new** — initializə olunmuş T (slice/map/channel) / zero-value *T pointer.

**MixedCaps** — çoxsözlü ad konvensiyası; alt xətt YOX.

## N

**Named result parameter** — adlı nəticə: zero başlanğıc + naked return + defer-də mutasiya.

**Naked return** — arqumentsiz `return` — cari named dəyərləri qaytarır.

## P

**Panic** — runtime ölümcül xəta; stack unwind başlayır; kitabxanalarda son çarə.

**PathError** — Op+Path+Err: mənbəyi tanıdan xəta strukturu.

**Pointer receiver** — caller-ın dəyərini mutasiya edən metod; yalnız pointer üzərində çağrılır.

**Prefiks konvensiyası** — xəta sətri mənbəyi göstərsin: "image: unknown format".

## R

**recover** — unwind-i dayandırır, panic dəyərini qaytarır; YALNIZ defer daxilində.

**Redeclaration (`:=`)** — eyni scope-da mövcud dəyişən + yeni dəyişən: elan YOX, mənimsətmə.

**Re-panic** — yalnız öz Error tipini yutma; beklenməyən panic-i buraxmaq.

**Reverse (array)** — paralel mənimsətmə ilə: `i, j = i+1, j-1`.

## S

**Semafor (channel)** — `sem <- 1` / `<-sem` — MaxOutstanding limiti.

**Semicolon insertion** — lekser qaydası: sətir sonu tokeni → avtomatik `;`.

**sync.Mutex zero value** — açıq mutex — initializasiyasız istifadə.

**Serving HTTP (4 tip)** — struct/int/channel/func Handler kimi — metod hər adlı tipə.

**Stringer** — `String() string` interfeysi; %v/-s formatları işə salır.

**Stack unwinding** — panic-in yuxarı qaçması — yolda defer-lər icra olunur.

## T

**template.Must** — başlanğıcda xəta → panic məqbul (sürətli uğursuzluq).

**Trace pattern** — `defer un(trace("a"))` — dərhal arqument qiymətlənməsinin istifadəsi.

**Type switch** — dinamik tipə görə budaqlanma; hər case öz tipi ilə eyni adlı var.

## V

**Value semantics (massiv)** — assign = tam kopya; ölçü tipin hissəsi.

**Variadic (`...`)** — `func Min(a ...int)`; ötürmə: `Sprintln(v...)`.

## W

**Worker pool** — fiksiləşmiş N handler goroutine — hər request-ə goroutine-dən qənaətcil.

**Zero value hazırlığı** — tip elə qurulmalı ki, new/var sonrası dərhal işləsin (Buffer/Mutex).
