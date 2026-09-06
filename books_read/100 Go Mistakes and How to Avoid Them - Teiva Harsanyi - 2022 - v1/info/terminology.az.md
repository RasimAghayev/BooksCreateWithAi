# 100 Go Mistakes and How to Avoid Them — Terminologiya (Azərbaycanca)

> Kitab boyu rast gəlinən terminlərin Azərbaycanca izahları. Format: `English Term (Azərbaycanca qarşılıq)`.

## A

**Abstractions should be discovered, not created** (Abstraksiyalar kəşf olunur, yaradılmır) — Rob Pike prinsipi: interface/generics yalnız konkret ehtiyacda.

**Alignment (Düzləndirmə)** — dəyişən adresinin öz ölçüsünün multiple-u olması; int32 → 4-ün multiple-u.

**any** — `interface{}` alias (Go 1.18) — bütün tip məlumatını itirir.

**Append race** — dolmamış slice üzərində konkurrent append eyni indeksə yazır.

**ariane 5** — 1996 raket qəzası: float64→int16 konversiya overflow-unun tarixi nümunəsi.

## B

**Backing array (Arxa massiv)** — slice-in/map-in altında duran massiv.

**Benchmark trap** — ResetTimer, compiler inline, observer effect kimi dəqiqlik pozan hallar.

**Blank identifier iqnoru** — `_ = f()` niyyətli xəta iqnorunun aşkar forması.

**Branch prediction (Budaq proqnozu)** — control hazard həlli; yanlış proqnoz = 10-20 dövrü flush.

**Build tag (`//go:build`)** — fayl-şərtli kompilyasiya; test kateqoriyalaşdırma.

**Busy-waiting (Boş-gözləmə)** — fasiləsiz loop yoxlama — CPU israfı (sync.Cond ilə həll).

## C

**Cache line (Keş sətri)** — 64 baytlıq kopyalanan yaddaş bloku.

**Cache placement policy** — blokların cache set-lərinə yerləşdirilməsi (set-associative).

**Chaining (`Use` sırası)** — middleware-lərin nested kompozisiyası.

**Channel size default = 1** — buffered kanal üçün magic number-lər əvəzinə default dəyər.

**CFS throttling** — Kubernetes kvotası aşılanda tətbiqin tam dondurulması.

**Closure loop tələsi** — closure-un bədəndən kənar loop dəyişəninə istinadı.

**Compulsory miss** — ilk dəfə çıxışdan yaranan qaçılmaz cache miss.

**Conflict miss** — set dolandıqda əvəzləmədən yaranan miss (critical stride).

**Concurrent mark-and-sweep** — Go GC: işarələ + təmizlə, əsasən konkurrent + 2 qısa STW.

**Consumer/producer side interface** — interface-in istifadə olunduğu / implementasiya yanında olması.

**Context detach** — ləğv siqnalı söndürülmüş, dəyərləri saxlayan custom context.

**CPU-bound / I/O-bound** — CPU / I/O sürəti ilə məhdud iş; worker pool ölçüsünü müəyyən edir.

**Critical stride** — eyni cache set-inə düşən addım (512 int64 = 4KB L1D-də).

## D

**Data hazard** — instraksiya nəticə asılılığı — ILP maneəsi.

**Data race vs race condition** — eyni yaddaş+yazma / sıra-timing asılılığı.

**Deadline (context)** — `WithTimeout` ilə vaxt limiti.

**Done() channel** — context ləğv channel-ı; bağlanma = bütün istehlakçılara yayımlanır.

**Detach context** — #61 həlli: Deadline/Done/Err söndürülür, Value saxlanılır.

## E

**Escape analysis (Çıxış analizi)** — kompilyatorun stack/heap qərarı.

**Errors wrapping** — %w (source açıq) vs %v (transform, source bağlı).

**Execution tracer** — hadisə əsaslı goroutine runtime vizualı.

## F

**Facilitative effect** — səhvlərin kontekstlə birlikdə yaddaşda saxlanması.

**False sharing** — ayrı dəyişənlərin eyni cache line-da olması → MESI invalidasiyası.

**Fast-path inlining** — slow path-i ayıraraq fast path-in inline olmasına imkan.

**Flaky test** — kod dəyişmədən keçən/fail olan test.

**Full slice expression** — `s[low:high:max]` — cap = max - low; append qoruması.

**Functional options pattern** — `Option func(*options) error` + With-prefiks closure-lar.

## G

**G/M/P** — Goroutine / OS thread (machine) / CPU core (processor).

**GOGC** — heap böyümə faizi GC tetikleyicisi (default 100).

**GOMAXPROCS** — paralel user-kod M limiti; default CPU core sayı (host-a görə!).

**Goroutine leak** — exit planı olmayan goroutine (2KB+ + resurs tutulumu).

## H

**Happens-before** — memory model hadisə sıralama zəmanəti.

**Hyper-Threading** — 1 fiziki core-un 2 məntiqi core-a bölünməsi.

## I

**ILP (Instruction-Level Parallelism)** — superscalar CPU-da paralel instraksiya icrası.

**Inlining budget (80)** — funksiyanın inline olma mürəkkəblik həddi.

**Interface pollution** — lazımsız interface-lərlə kod çirklənməsi.

**Interface Segregation Principle** — client istifadə etmədiyi metodlara asılı olmamalıdır.

## J-K-L

**JSON monotonic tələsi** — `time.Now()` marshal→unmarshal == poza bilər (m=+ hissəsi).

**L1/L2/L3** — cache səviyyələri: ~1ns / 4x / 10x; RAM 50-100x yavaş.

**Load factor (map)** — bucket başına element; 6.5-i keçəndə map 2x böyüyür.

**Lock contention (mənafəə rəqabəti)** — global map-lərdə çoxlu goroutine yazma rəqabəti.

## M

**MESI** — cache coherency protokolu: Modified/Exclusive/Shared/Invalid.

**Mechanical sympathy** — sistemin dizaynına uyğun işlərək optimal performans.

**Memory model (Go)** — yaz/oxu sıralama zəmanətlərinin spesifikasiyası.

**Mid-stack inlining** — non-leaf funksiyaların inline olması (Go 1.9+).

**Mutex (mutual exclusion)** — kritik seqment qoruyucusu.

## N

**Named result parameters** — adlı nəticələr — zero value başlanğıcı tələ yarada bilər.

**Nil channel** — əbədi bloklama; select case söndürmə aləti.

**Notification channel** — `chan struct{}` — datasız siqnal.

**Naked return** — arqumentsiz return; qısa funksiyalarda.

## O

**Observer effect** — təkrar istifadə CPU cache qazandırır → hər iterasiyada yeni data.

**Overflow (silent)** — run-time integer aşımı panic YOXDUR — MaxInt yoxla.

## P

**Padding** — ayrıcı baytlar (alignment və ya false sharing qarşısı).

**Postel's law** — qəbulda liberal (interface), etməkdə konservativ (struct).

**pprof** — Go profiler: CPU/heap/goroutine/block/mutex.

**Preemptive scheduling** — Go 1.14+: 10ms-dən uzun goroutine məcburi paylaşır.

**Prepared statement** — precompiled SQL — effektivlik + injection qoruması.

**Producer-side interface** — implementasiya yanında interface (əksər hallarda qeyri-idiomatik).

## R

**Race detector (-race)** — runtime instrumentasiya; vector clock; false positive YOXDUR.

**Runtime tracer** — go tool trace — GC/paralellik/goroutine vizualı.

## S

**Sentinel error** — qlobal error dəyəri (ErrFoo, sql.ErrNoRows, io.EOF) — expected xətalar.

**Set-associative cache** — set-lərə bölünmüş cache; blok → set adresin index bitləri ilə.

**Sharing down/up** — pointer qəbul etmək (stack) / pointer qaytarmaq (heap).

**Stack frame** — funksiyaya məxsus yaddaş — self-cleaning.

**Striding** — unit (kontigu) / constant (addımlı) / non-unit (linked list — proqnozsuz).

**Struct padding** — sahə sırası yanlış olsa 14 bayt israf (24 vs 16 bayt).

**Struct-of-slices** — `Bar{a, b []int64}` — slice-of-structs-dan daha yaxşı spatial locality.

**sync.Cond** — şərt dəyişəni — təkrar broadcast üçün (Wait = unlock-suspend-relock).

**sync.Pool** — obyekt reuse hovuzu — GC-də boşalır, cache deyil.

**sync tip kopyası** — Mutex/WaitGroup və s. heç vaxt kopyalanmaz.

## T

**Table-driven test** — map + subtests — dublikatsız çoxsaylı case.

**Temporal/spatial locality** — eyni / qonşu yaddaşın təkrar istifadəsi.

**Threshold (paralellik)** — bu ölçüdən kiçik işlər sequensial; benchmark ilə tapılır.

**Time.After leak** — loop-da hər çağırış timer yaddaşı saxlayır (200B × milyonlar).

**TrimRight/Left vs TrimSuffix/Prefix** — SET (təkrar) / dəqiq affiks (bir dəfə).

**Truncate(0)** — time.Time-dan monotonic hissəni silmə.

## V

**Vector clock** — qismi sıralama strukturu — race detectorun əsası.

**Virtual heap hilesi** — `make(1GB)` + mmap — fiziki yaddaş istehlak etmədən GC sıxlığını azaltmaq.

## W

**Wall vs monotonic clock** — NTP-sıçrayanlı / həmişə-irəli; `m=+` hissəsi.

**Work stealing** — boş P-nin digər P-lərdən goroutine oğurlaması.

**Worker pooling** — fiksiləşmiş goroutine hovuzu — ölçü workload tipinə görə (CPU→GOMAXPROCS, IO→xarici sistem).

**Wrapping (error)** — source error-u wrapper-də açıq saxlamaq; coupling riski.
