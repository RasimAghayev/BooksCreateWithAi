# Chapter 8 — Memory Management (Yaddaş İdarəsi)

## Bu chapter nədən bəhs edir?
Go GC-nin daxilinə: stack vs heap ayrımı, tri-color concurrent mark-and-sweep alqoritmi, GOGC/GODEBUG tuning, GC pacer, memory ballast (Twitch hadisesi), GOMEMLIMIT (Go 1.20+) və experimental memory arenas.

## Əsas fikirlər

### 1. GC nədir — "memory inference" problemi
**Qabaqcadan (GC-siz):** memory leak, dangling pointer, double free.
**GC-nin işi:** heap allocation-ları izlə → lazımsızları azad et → istifadədəkiləri saxla.

**2 strategiya:** tracing (istinad zənciri ilə çatılabilənlər CANLI) və reference counting. **Go = tracing GC.**

**Tarixi nailiyyət:** GC cycle 300ms (Go 1.0) → **0.5ms** (müasir).

**Mif dağıtma:** "GC avtomatikdir, yaddaşa fikir vermə" — Roomba evi təmizləmir; linter kodu refaktor etmir; GC dərk etmək performansın açarıdır.

### 2. Stack vs Heap
| | Stack | Heap |
|---|---|---|
| Ömür | funksiya çağırışına bağlı, proqnozlaşdırıla bilən | dinamik, funksiyadan sonra da yaşaya bilər |
| Mexanizm | LIFO — pointer yuxarı/aşağı | kompleks bookkeeping |
| Sürət | ÇOX sürətli | yavaş + GC overhead |
| Limit | kiçik — stack overflow riski | böyük, global erişim |

**Escape analysis (compiler):** dəyişən funksiyadan qaçmırsa → stack; istinad ötürülürsə/qaytarılırsa → heap "escapes".

**Qayda:** Scope-u dar saxla; lazımsız pointer/referanslardan ehtiyat et.

### 3. Tri-color concurrent mark-and-sweep
**Concurrent:** proqram işləyərkən paralel (real-zaman sistemlər üçün kritik).

**Tri-color (traffic light):**
- **Ağ:** talebi qərarsız (başlanğıc)
- **Boz:** araşdırılmalı
- **Qara:** təmiz işlənib, İSTİFADƏDƏ

**Mark fazası:**
1. **STW (stop-the-world) <0.3ms:** root set təyini (stack, qlobal, xüsusi yerlər) — "suya dalmazdan əvvəl nəfəs"
2. **Concurrent marking:** root-dan çatılan ağlar boz → emal → qara

**GC resurs büdcəsi:** 25% CPU (əsas GC işi) + 5% (mark assists — allocate edən goroutine-lar GC geridə qalsa kömək edir).

**Sweep:** ağ qalanlar = zibil → deallocate.

### 4. GOGC — GC termostatı
**Default 100:** GC-dən sonra heap-in 100% böyüməsinə icazə → yeni cycle.
- **Aşağı (50):** tez-tez cycle → kiçik heap, çox CPU
- **Yüksək (200):** seyrək cycle → az CPU, böyük yaddaş
- **off:** GC tam söndürülür (qısaömürlü proqramlar üçün; risk — nəzarətsiz artım)
- **0-dan böyük istənilən tam dəyər**

### 5. GC pacer — orkestr dirijoru
**Vəzifə:** GC cycle-larının ZAMANINI tənzimlə — heap ölçüsü + allocation rate + GOGC hədəfinə əsasən.

**Adaptiv:** sürətli allocate → tez-tez cycle; yavaş allocate → gec cycle.

**Simbioz:** GOGC faiz təyin edir; pacer onu thresholdda istifadə edir. Pis tuned pacer → ya çox cycle (performans itkisi), ya gecikmiş cycle (yaddaş artımı).

### 6. GODEBUG=gctrace=1 — GC diaqnostikası
```
gc 1 @0.019s 2%: 0.015+2.5+0.003 ms clock, 0.061+0.5/2.0/3.0+0.012 ms cpu, 4->4->1 MB, 5 MB goal, 4 P
```
| Sahə | Məna |
|---|---|
| gc 1 @0.019s | cycle № + başlanğıc vaxtı |
| 2% | proqram vaxtının GC payı (YÜKSƏK = problem!) |
| 0.015+2.5+0.003 ms clock | STW-sweep-term + concurrent mark + STW-mark-term |
| 4->4->1 MB | heap: başlanğıc / orta / son |
| 5 MB goal | növbəti hədəf heap |
| 4 P | processor sayı |

**Diaqnostika:** yüksək %-faiz; uzun STW; heap böyüyür-düşmür = leak; yüksək CPU = ineffektiv yaddaş istifadəsi.

### 7. Memory ballast — Twitch hadisesi (2019)
**Nədir:** Heç vaxt istifadə olunmayan BÖYÜK yaddaş bloku — GC-i aldadır.
```go
ballast := make([]byte, 10<<30)    // 10 GiB!
```
**Mexanizm:** İstifadə olunmasa da CANLI sayılır (istinad var) → heap bazası şişir → GC gec trigger → daha az, daha yumşaq cycle-lar.

**Twitch nəticələri:** GC cycle-ları ~99% azaldı; CPU ~30% düşdü; p99 latency ~45% azaldı; virtual memory-də yaşayır (ucuz).

**Məhdudiyyətlər:** memory-sensitive mühitlər (konteynerlər), dinamik istifadə, low-latency (proqnozlaşdırıla bilən pauza prioritet), kiçik heap tətbiqləri, GOGC tuning kifayət olanda, arxa kod problemini maskalayanda.

**STATUS:** Go 1.19-a qədər aktual; 1.20+-dan GOMEMLIMIT əvəz edir.

### 8. GOMEMLIMIT (Go 1.20+)
**Nədir:** Runtime yaddaşına SOFT limit — heap + runtime idarəli yaddaş (binary mapping, C yaddaşı, OS-held HARİC).

```bash
GOMEMLIMIT=2GiB      # vahid: B, KiB, MiB, GiB, TiB (IEC 80000-13)
```
- Default: math.MaxInt64 (deaktiv)
- Runtime-da: `runtime/debug.SetMemoryLimit()`

**Soft cap:** Limit yaxınlaşınca GC daha aqressiv işləyir — amma mütləq maneə DEYİL ("sürət xəbərdarlıq nişanı — avtomobili fiziki yavaşlandıra bilməz").

**Ballastla birgə istifadə = 2 saat taqvimi — redundant.**

### 9. Memory arenas (Go 1.20, experimental)
**Qurulum:** `GOEXPERIMENT=arenas` env + `import "arena"`.
**Xəbərdarlıq:** rəsmi dəstək/garantiya YOXDUR — gələcəkdə silinə bilər.

**Konsept:** Kontiguus yaddaş regionından obyekt yarat → hamısını BİRGƏ, bir anda azad et — minimal GC overhead.

```go
mem := arena.NewArena()

p := arena.New[Person](mem)               // arenadan referans AL (adi "yarat-qoy" DEYİL)
slice := arena.MakeSlice[string](mem, 100, 100)

p2 := arena.Clone(p1)                      // arenadan heap-ə köçür (GC-ə qaytar)

mem.Free()                                 // hamısı BİRGƏ azad olunur
```

**Təhlükə + müdafiə:** Free-dən SONRA istifadə = bug:
```go
mem.Free()
o.Num = 123          // <- problem!
```
```bash
go run -asan main.go
# accessed data from freed user arena 0x40c0007ff7f8
```

**İdeallar:** gRPC (hər RPC çox obyekt allocate edir; C++ versiyası artıq arenas işlədir), JSON (fastjson — allegedly 15× standart kitabxanadan sürətli).

**Deciding questions (müəllifin sualları):**
1. Data-n varmı? (Yoxdursa — tahmin edirsən)
2. Çoxmu allocation? (Az isə arena YADDAŞI ARTIRIR — rule of thumb: arena = 8MB)
3. Eyni kiçik struktur? → bəlkə **sync.Pool** düzgün alətdir
4. Hot path? → yoxsa premature optimization — əvvəl GC/GOMEMLIMIT kombinasiyaları

## Əsas terminlər
- Memory inference — "hansı yaddaşı azad etməli"
- Tracing vs reference counting GC
- Escape analysis — stack/heap qərarı
- STW (stop-the-world) — qısa pauza (<0.3ms)
- Tri-color: white/gray/black
- Mark assists — allocate edən goroutine-ların GC köməyi
- 25% + 5% CPU GC büdcəsi
- GOGC (default 100; off)
- GC pacer — adaptiv cycle zamani
- GODEBUG=gctrace=1 — GC jurnal
- Memory ballast — 10GiB süni heap (Twitch)
- GOMEMLIMIT — soft limit (IEC vahidləri)
- arena.NewArena/New[T]/MakeSlice/Clone/Free
- GOEXPERIMENT=arenas; -asan (address sanitizer)

## Praktik nəticə
1. Stack-də qalmaq üçün scope dar, pointerlardan ehtiyatlı — escape = GC xərci.
2. gctrace=1 ilə GC%-faiz + STW + heap trend izlə: yüksək %, uzun STW, düşməyən heap = 3 əsas siqnal.
3. GOGC tuning-ə əvvəl bax (müdaxiləsi az), sonra GOMEMLIMIT; ballast artıq legacy (1.20+).
4. Konteynerdə GOMEMLIMIT təyin et — soft limit OOM qorunması + ballastın müasir alternativi.
5. Arenas üçün: data ilə sübut + çox allocation + hot path — yoxsa sync.Pool və ya heç nə.
6. Arena Free-dən sonra istifadə — asan bug; -asan ilə yoxla.
7. Performans tahmin oyunu deyil — ölç (növbəti fəsil: profiling).

## Mənbə
Pages: 147-159 (PDF səh. 168-181)
