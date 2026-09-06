# Go in Action — Terminologiya lüğəti (AZ)

Texniki terminlər `English (Azərbaycanca)` formatında — SYSTEM_PROMPT qayda 4 üzrə.

## Dilin əsasları (Ch 1-2)
- Goroutines (qorutinlər) — yüngül, runtime tərəfindən idarə olunan paralel funksiyalar
- Channels (kanallar) — goroutine-lər arasında tipli mesajlaşma quruluşu
- Concurrency (eyni-vaxtlılıq) — çox işin idarə edilməsi
- Parallelism (paralellik) — çox işin eyni anda fiziki icrası
- Composition (kompozisiya) — tiplərin daxil edilməsi ilə qurulması
- Duck Typing (ördək tipi) — davranışa əsaslanan tip uyğunluğu
- Garbage Collector / GC (zibil toplayıcı) — yaddaşın avtomatik idarəsi
- Static Typing (statik tip sistemi) — tipin kompilyasiya vaxtı bilinməsi
- Package (paket) — kompilyasiya vahidi, namespace
- Blank Identifier (boş identifikator, `_`) — dəyəri ignorə etmək üçün
- Zero Value (sıfır dəyər) — initiazlaşdırılmayan dəyişənin default dəyəri
- Anonymous Function (anonim funksiya) — adsız funksiya (closure üçün əsas)
- Closure (klouzer) — xarici scope dəyişənlərinə birbaşa çıxış
- Pass by Value (dəyərlə ötürmə) — Go-da bütün ötürmələrin təbiəti
- Struct Tag (struktur teqi) — `` `json:"site"` `` — sahə mapinqi metadatası

## Data strukturları (Ch 4)
- Array (massiv) — sabit uzunluqlu ardıcıl yaddaş bloğu
- Slice (hissə) — pointer + length + capacity-dən ibarət dinamik görünüş
- Underlying Array (altta yatan massiv) — slice-ın real data saxlayan hissəsi
- Length / Capacity (uzunluq / tutum) — slice-ın 2 ölçüsü
- Reference Type (istinad tipi) — header value saxlayan tiplər (slice, map, channel, interface, func)
- Header Value (başlıq dəyəri) — reference tipinin daxili strukturu
- Hash Table (xəş cədvəli) — map-in implementasiyası
- Bucket (çən) — xəş cədvəlinin saxlama vahidi
- Multidimensional (çoxölçülü) — strukturların kompozisiyası ilə əldə olunur

## Tip sistemi (Ch 5)
- User-defined Type (istifadəçi tipi) — `struct` və ya base type əsasında
- Base Type (əsas tip) — `type Duration int64`-dəki int64
- Method / Receiver (metod / qəbuledici) — funksiyaya tip bağlayan parametr
- Value / Pointer Receiver (dəyər / göstərici qəbuledici) — kopya vs paylaşım
- Primitive / Nonprimitive Nature (ilkin / qeyri-ilkin təbiət) — dəyişməzlik qərarının meyarı
- Interface (interfeys) — davranış müqaviləsi
- iTable (interfeys cədvəli) — tip məlumatı + metod siyahısı saxlayan daxili struktur
- Concrete Type (konkret tip) — interfeys daxilində saxlanan real tip
- Method Set (metod dəsti) — hansı dəyərlərin interfeys uyğunluğu qaydaları
- Polymorphism (çoxşəkillilik) — eyni interfeys üzərindən fərqli davranış
- Type Embedding (tip daxiletməsi) — kompozisiya yolu ilə genişləndirmə
- Inner Type Promotion (iç tipin qaldırılması) — iç identifikatorların xaricə çıxması
- Exported / Unexported (ixrac / qeyri-ixrac) — böyük / kiçik hərflə görünüş qaydası
- Factory Function (fabrik funksiyası) — `New` konvensiyası ilə konstruktor

## Concurrency (Ch 6-7)
- Logical Processor (məntiqi prosessor) — OS thread-ə bağlı scheduler vahidi
- Run Queue (icra növbəsi) — goroutine-lərin gözləmə növbəsi (global + local)
- Preemption (zorla dayandırma) — goroutine-in processor-dən qaldırılması
- Race Condition (yarış şəraiti) — sinxronizasiyasız paylaşılan resursa eyni anda çıxış
- Race Detector (`-race`) — data race aşkarlayan build flag-i
- Atomic Functions (atomik funksiyalar) — `sync/atomic` paketi
- Mutex / Mutual Exclusion (qarşılıqlı istisna) — `sync.Mutex` kritik bölmə qoruyucusu
- Critical Section (kritik bölmə) — bir vaxtda bir goroutine üçün kod hissəsi
- CSP / Communicating Sequential Processes (sequnial proseslərin kommunikasiyası) — mesaj ötürmə modeli
- Unbuffered / Buffered Channel (bufersiz / bufferli kanal) — sinxron / asinxron mübadilə
- WaitGroup (gözləmə qrupu) — sayğaclı semafor
- Backpressure (geri təzyiq) — dolu kanalın göndərəni bloklaması
- Network Poller (şəbəkə sorğulayıcısı) — I/O gözləyən goroutine-lərin idarəsi
- Deadlock (qarşılıqlı bloklanma) — qapalı kanaldan drain etməmək kimi səhvlər

## Standart kitabxana (Ch 8)
- Standard Library (standart kitabxana) — 100+ paket, backward-compatibility zəmanəti
- GOROOT / GOPATH (quraşdırma / workspace yolları)
- iota — konstant blokunda auto-artan sayğac
- Bit Shift / Bitwise OR (`<<` / `|`) — flag pattern-inin əsası
- stdout / stderr (standart çıxış / xəta çıxışı)
- io.Writer / io.Reader (yazıcı / oxuyucu) — data axını interfeysləri
- MultiWriter (çox-Writer) — bir yazını bir neçə destinasiyaya yönləndirir
- Discard (uducu) — yazılan hər şeyi məhv edən Writer
- Marshal / Unmarshal (seriyalaşdırma / geri-çevirmə) — JSON ↔ Go dəyəri
- MarshalIndent (formatlı seriyalaşdırma) — pretty-print JSON
- Reflection (refleksiya) — runtime tip məlumatı (empty interface ilə işləyir)
- Empty Interface (`interface{}`) — hər tipi qəbul edən tip
- Type Assertion (`x.(T)`) — interfeysdən konkret tipə çevirmə
- EOF (End Of File) — axın sonu

## Test (Ch 9)
- Unit Test (vahid test) — `Test` prefiksli funksiya
- Table Test (cədvəl testi) — çoxparametrli test pattern
- Positive / Negative Path (müsbət / mənfi ssenari)
- Mocking (mock etmə) — saxta resurs simulasiyası
- httptest.Server — lokal mock HTTP server
- ResponseRecorder (cavab qeydedicisi) — httptest.NewRecorder — cavabı yadda saxlayan Writer
- HandlerFunc — adi funksiyanı HTTP handler-ə çevirən adapter
- Black-box Testing (`_test` paketi) — yalnız exported API test olunur
- Example Function (nümunə funksiyası) — sənədləşmə + test (Output: markeri)
- Benchmark (`Benchmark` prefiksli funksiya) — performans ölçmə
- ns/op / B/op / allocs/op — nanosaniyə / byte / allokasiya — əməliyyat başına
