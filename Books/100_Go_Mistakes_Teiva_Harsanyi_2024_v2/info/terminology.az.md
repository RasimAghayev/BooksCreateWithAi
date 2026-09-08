# 100 Go Mistakes — Terminologiya (AZ)

| Termin | Azərbaycanca | İzah |
|---|---|---|
| Variable shadowing | dəyişən kölgələməsi | daxili scope xarici adı gizlədir (#1) |
| Line of sight | görüş xətti | happy path solda düz axın (#2) |
| Interface pollution | interfeys çirklənməsi | lazımsız interfeyslər (#5) |
| Producer side | istehsalçı tərəfi | interfeys yeri: istehlakçıda olmalıdır (#6) |
| Functional options | funksional seçənəklər | variadic closure konfiq (#11) |
| Utility package | util paketi | mənasız ad — antipattern (#13) |
| Octal literal | səkkizlik literal | 0-prefix: 010 = 8 (#17) |
| Backing array | arxa massiv | slice-ın altındakı real yaddaş |
| Capacity leak | kapasitet sızıntısı | reslicing böyük array-i saxlayır (#26) |
| reflect.DeepEqual | dərin müqayisə | universal amma yavaş (#29) |
| Named result | adlı nəticə | naked return təhlükəli ola bilər (#43) |
| nil receiver | nil qəbuledici | nil pointerdə metod çağırıla bilər (#45) |
| Error wrapping | xəta sarılması | %w + errors.Is/As (#49-51) |
| G-M-P scheduler | G-M-P | goroutine-processor-thread modeli (#56) |
| Work stealing | iş oğurlama | boş P dolu P-dən goroutine götürür |
| Happens-before | baş-vermə-öncəsi | memory model zəmanətləri (#58) |
| CPU-bound / IO-bound | prosessor/GİRİX-çıxış asılı | worker sayı seçiminin əsası (#59) |
| Context detach | kontekst ayırma | cancel-dən təcrid olunmuş ctx (#61) |
| Nested select | iç-içə seçim | prioritet üçün pattern (#64) |
| Nil channel | nil kanal | block edir — merge-də istifadə olunur (#66) |
| sync.Cond | şərt dəyişəni | broadcast/repeat signal (#72) |
| errgroup | errgroup | xətalı WaitGroup (x/sync) (#73) |
| Monotonic clock | monoton saat | wall clock dəyişə bilər; time.Sub üçün (#77) |
| Transient resource | keçici resurs | resp.Body, rows, fayl — bağlanmalı (#79) |
| Slowloris | Slowloris | yavaş header hücumu — ReadHeaderTimeout (#81) |
| Build tag | build teqi | //go:build integration — test ayrımı (#82) |
| -race | -race | race detektor — CI-də məcburi (#83) |
| -shuffle | -shuffle | test sırasını qarışdır (Go 1.17+) (#84) |
| Eventually | Eventually | testify retry-assertion (sleep əvəzi) (#86) |
| Cache line | keş xətti | 64 bayt — spatial locality vahidi (#91) |
| Critical stride | kritik addım | 512 sütun = conflict miss (#91) |
| False sharing | yalan paylaşım | eyni cache line-da ayrı dəyişənlər (#92) |
| Padding | doldurma | sahələri ayrı line-lara sökür (#92) |
| Data hazard | data təhlükəsi | RAW asılılığı — ILP azaldır (#93) |
| Alignment | düzləndirmə | struct sahə sırası padding-i dəyişir (#94) |
| Escape analysis | kaçış analizi | stack/heap qərarı (#95) |
| sync.Pool | sync.Pool | obyekt təkrar istifadə hovuzu (#96) |
| Fast-path inlining | fast-path inline | kiçik yol inline + slow path ayrı (#97) |
| pprof | pprof | profil: cpu/heap/goroutine/mutex/block (#98) |
| GOGC | GOGC | GC tetikleme əmsalı (default 100) (#99) |
| CFS quota | CFS kvotası | K8s CPU limiti + GOMAXPROCS uyuşmazlığı (#100) |

## Səhv → həll sürətli cədvəli (seçmə)

| # | Səhv | Bir sətirlə həll |
|---|---|---|
| 1 | `:=` kölgələmə | xaricdə elan, daxildə `=` |
| 30 | range kopyası | `accounts[i]` və ya pointer slice |
| 32 | `&customer` hamısı son | `customer := customer` / `&customers[i]` |
| 35 | defer loop-da | helper funksiya / closure |
| 43 | naked return + defer | nəticəni açıq qaytar |
| 45 | nil MultiError | `if m != nil { return m }; return nil` |
| 47 | defer status kopyası | `defer func(){ notify(status) }()` |
| 63 | loop closure | `i := i` (1.22+ avtomatik) |
| 69 | append race | kopya çıxarıb işlə |
| 76 | time.After loop | NewTimer + Reset |
| 81 | DefaultClient | `&http.Client{Timeout: x}` |
| 100 | GOMAXPROCS hostda | runtime.GOMAXPROCS(n) / automaxprocs |
