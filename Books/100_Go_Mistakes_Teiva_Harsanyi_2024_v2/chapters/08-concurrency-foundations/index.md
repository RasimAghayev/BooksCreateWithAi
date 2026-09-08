# Chapters 8-10 — Concurrency Foundations & Practice, Stdlib (#55-#82) (səh. 184-283)

## Bu fəsillər nədən bəhs edir?

(8) Concurrency vs parallelism, G-M-P scheduler, parallel merge sort, kanal vs
mutex, memory model, workload tipləri (CPU/IO-bound), context. (9) Context
çirklənməsi, goroutine lifecycle, loop dəyişəni, select determinizmi, nil
kanallar, kanal ölçüsü, String formatlaşdırma deadlock-u, append race, mutex
hususiyyətləri, WaitGroup, sync.Cond, errgroup, sync kopyası. (10) time
yanlışlıqları, JSON embedded, monotonic clock, SQL, transient resurslar,
DefaultClient.

## Əsas səhvlər və həllər

### #55: Concurrency ≠ Parallelism
- **Concurrency:** struktur (süd maşını nümunəsi — addımların ayrılması);
  **Parallelism:** icra (addım səviyyəsində çox işçi)
- Concurrency PARALLELİZMİ MÜMKÜN edir — amma zəmanət vermir

### #56: Paralellik həmişə sürətli deyil
- **G-M-P scheduler:** G (goroutine) → P (processor, local queue) → M (OS
  thread); M P gözləyir; work-stealing
- Parallel merge sort: 10× YAVAŞ goroutine yaratma + sync xərci dominant oldu;
  həll — GİBRİLTƏ (threshold): `if len(s) < max { sequential }` → sürətli

### #57: Kanal yoxsa mutex?
- **Kanal:** bir goroutine-dən nəticə ÖTÜRÜLÜRÜRSƏ + sahiblik köçürülür (data
  flow). **Mutex:** shared state QORUNURSA + bir neçə goroutine oxu/yazı
- "Kanallar orkestrləyir, mutex serializə edir" — amma seçim kontekstdən

### #58: Race + memory model
- Data race: 2 goroutine eyni dəyişən, ən azı 1 yazma, 0 senxronizasiya →
  NONDETERMINIZM (i++ = oxu-artır-yaz; ikisi 0 oxuyur, hər ikisi 1 yazır)
- Həll: atomic, mutex, channel
- **Happens-before zənciri:** kanal GÖNDƏRİŞİ receive-dən ƏVVƏL; close → bütün
  receive-lərdən sonra; dar amma möhkəm zəmanətlər

### #59: Workload tipləri
- **CPU-bound:** worker pool = GOMAXPROCS (nüvə sayı) — daha çox artıq; runtime.GOMAXPROCS
- **IO-bound:** nüvə sayından ASILI DEYİL — 128 də ola bilər (bloklanma çoxdur)
- Amma gündəlik işdə: "worker sayı — queue-backpressure və qəbul edilən
  gecikməyə görə"

### #60: Context
- 4 element: deadline, cancel siqnalı, values; kontekst QƏBUL edən funksiya
  yazırıqsa: **bloklama əməliyyatları select + ctx.Done() ilə**
```go
select {
case <-ctx.Done():
    return ctx.Err()
case v := <-ch:
    // ...
}
```
- Value yalnız unexported açar tiplərli (çirklənmə qarşısı)
- WithCancel/WithTimeout: main qurur, bütün asılılıqlara ötür

### #61: Context çirklənməsi
- Handler kafka publish üçün HTTP request ctx ötürür → request bitəndə publish
  ləğv! Həll: **detach** — context.Background() və ya custom "cancel() no-op"
  wrapper

### #62: Goroutine lifecycle
- Yaradılan goroutine: NE VAXT dayanacaqı MƏLUM olmalıdır (exit point)
- `newWatcher` → `go w.watch()` → `close()` signal kanalı; yoxsa LEAK

### #63: Loop dəyişəni + closure (klassik!)
```go
for _, i := range s {
    go func() { fmt.Print(i) }()   // hamısı SON i (Go ≤1.21)
}
// həll: i := i  və ya  go func(i int) {...}(i)
// (Go 1.22+ avtomatik)
```

### #64: select + determinizm
- Birdən çox case HAZIRDIRSA → RANDOM seçim!
- Prioritet lazımdısa: **nested select pattern**:
```go
for {
    select {
    case v := <-messageCh:  // ...
    case <-disconnectCh:
        for {
            select {
            case v := <-messageCh:  // öncəlik: mesajları bitir
                // ...
            default:
                return
            }
        }
    }
}
```

### #65: (worker-də tapşırıq bölgüsü — sync həllər)

### #66: nil kanallar
- nil kanaldan oxu/göndərmə = BLOCK (sonsuz) — bu XÜSUSİYYƏT kimi istifadə:
- merge 2 kanal: biri bağlananda → `ch1 = nil` → select artıq yalnız ch2-yə
  baxır — qəliz closed-flag-lərdən təmiz həll
- Siqnal üçün: `chan struct{}` — 0 bayt

### #67: Kanal ölçüsü
- Unbuffered = SENXRONİZASİYA zəmanəti; buffered = ASİNXRAN amma gizli
  deadlock mümkün — "buffered default olaraq al, amma nə üçün olduğunu bil"

### #68: String formatlaşdırma side-effect
- `fmt.Sprintf("%v", customer)` — String() metodu çağırır → daxildə mutex →
  çağıran da mutex tutubsa DEADLOCK
- Map açarı kimi ctx formatlaşdırma (Randomization!) — hər çağırışda fərqli

### #69: append race
- İki goroutine `append(s, ...)` eyni backing-ə → kapasite dolubsa yazma
  üst-üstə — KOPYA apararaq izolyasiya

### #70: Mutex + slice/map
- RWMutex: RLock oxu paralel; amma "iterasiya + yazma" funksiyada → BÜTÜN
  funksiya critical section (parçalama deadlock riski)

### #71: WaitGroup
- `wg.Add(1)` HARADA? — spin-updan ƏVVƏL (race: Done Add-dən əvvəl işləyə
  bilər); loop sayı məlumsa `wg.Add(n)` əvvəlcədən

### #72: sync.Cond (unutqanlıq)
- Broadcast/repeat-signaling: bir hadisədən BİR NEÇƏ goroutine xəbərdar et →
  `Cond.Wait/S.Signal/Broadcast`; Wait daxilə mutex TUTULMUŞ gəlir (Wait
  unlock→sleep→relock edir)

### #73: errgroup
```go
g, ctx := errgroup.WithContext(ctx)
for i, circle := range circles {
    g.Go(func() error { return process(ctx, circle) })
}
if err := g.Wait(); err != nil { ... }
```
- Xəta + context avtomatik ləğv + nəticə aqreqasiyası

### #74: sync tipin kopyalanması
- struct mutex sahəsi + value pass → KOPYA mutex = ayrı lock!
- go vet yaxalayır; qayda: sync tipləri (Cond/Map/Mutex/RWMutex/Pool/WaitGroup)
  heç vaxt value kopya olunmasın

### #75: time.Duration
- `time.Duration` int64 alias: `time.Second` MƏCBURİ — `100` YOX (nanosaniyə!)

### #76: time.After + memory leak
- Loop-da `time.After(5s)` → hər iterationda yeni timer; expired olana qədər
  yaddaşda. Həll: `timer := time.NewTimer(...)` + `timer.Reset(...)`

### #77: JSON embedded
- Embedded `time.Time` → sahə adı ÇAKIŞIR (ID yox olur!). Həll: adlandırılmış
  sahə
- time.Now() = wall + MONOTONIC; müqayisə üçün monotonic daxildir (t.Sub doğru);
  JSON-a girəndə yalnız wall çap olunur
- `map[string]any` — ƏDƏD həmişə float64!

### #78: SQL
- `sql.Open` YALNIZ validasiya (bağlantı YOX) → `db.Ping()`
- Prepared statements (təkrar istifadə + injection qoruması)
- NULL sütun: `*string` pointer scan
- Connection pool: SetMaxOpenConns/SetMaxIdleConns/SetConnMaxLifetime(Idle)

### #79: Transient resurslar
- resp.Body, sql.Rows, fayl — HƏMİŞƏ bağla; `resp != nil && err != nil` → body
  yalnız resp non-nil-də mövcuddur: `if resp != nil { defer resp.Body.Close() }`

### #80: HTTP handler-da return unutmaq
- WriteHeader + yazı + KOD DAVAM — superfluous WriteHeader warning

### #81: DefaultClient/Server
- http.Get = DefaultClient — timeout YOX (sonsuz gözləmə, connection yığılması)
- Hər zaman custom: `&http.Client{Timeout: ...}`; server: ReadHeaderTimeout
  (Slowloris hücumu qarşısı)

### #82: Test kateqoriyaları (test piramidi — unit/integration/E2E; -short build
tag-lərlə ayır)

## Əsas terminlər

- G-M-P Scheduler
- Work Stealing
- Happens-Before (memory model)
- Workload CPU/IO-bound
- Context Detach
- Nested Select
- Nil Channel Semantics
- sync.Cond (Broadcast)
- errgroup
- Monotonic Clock
- Slowloris

## Praktik nəticə

- Goroutine = lifecycle müqaviləsi: exit point olmadan başlatma
- Determinizm istəyirsənsə select-i NEST et; bitmiş kanalı nil-ə çevir
- Paralellikdə threshold qoy (kiçik iş sequensial)
- DefaultClient YOX; time.After loop-da YOX; sync tipləri kopyalanmaz

## Mənbə

Pages: 184-283 (Chapters 8-10, 100 Go Mistakes 2nd ed.)
