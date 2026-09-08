# System Programming Essentials with Go — Terminologiya (Azərbaycanca)

> Texniki terminlər `English (Azərbaycanca qarşılığı)` formatında — kitabın ardıcıllığı ilə.

## Giriş + Konkurensi (Ch 1-2)
- **System programming (sistem proqramlaşdırması)** — OS/hardware ilə yaxın işləyən effektiv kod
- **CSP (Communicating Sequential Processes)** — Go-un konkurensi modeli
- **Goroutine** — yüngül icra vahidi (KB stack, dinamik böyümə)
- **WaitGroup** — işçi qrupu gözləməsi (Add/Done/Wait)
- **Data race (data yarışı)** — paralel yazma qarşılıqlı pozuntusu (`-race` ilə tutulur)
- **Atomic operation** — bölünməz əməliyyat (AddInt64/CompareAndSwap)
- **Mutex / RWMutex** — eksklüziv / oxu-çox-yazma-1 kilid
- **Unbuffered channel** — senxron; göndərici+qəbulçu eyni anda
- **Buffered channel** — asinxron; buffer dolana qədər blok yoxdur
- **Delivery guarantee (çatdırılma zəmanəti)** — qəbul edilməmiş göndərmə "itmir" (buffered-də itir!)
- **Latency vs throughput** — gecikmə vs axın
- **State / Signaling** — channel-in 2 vəzifəsi

## Sistem çağırışları (Ch 3)
- **System call (sistem çağırışı)** — kernel xidmətinə keçid (kullanıcı→kernel modu)
- **Syscall interface** — xidmət kataloqu (nömrə ilə)
- **syscall / os / golang.org/x/sys** — 3 səviyyəli API
- **strace** — syscall trace aləti
- **File descriptor (fayl deskriptoru)** — açıq resurs nömrəsi
- **stdin / stdout / stderr** — standart axınlar (0/1/2)
- **Redirection (yenidən yönləndirmə)** — `>`, `2>&1`
- **os.Pipe** — stdini sınaqdan keçirmə üçün suni daxiletmə

## Fayllar (Ch 4)
- **File permission bitləri** — rwx (owner/group/other), 0755 kimi
- **umask** — default permission istisnası
- **filepath.WalkDir** — qovluq ağacı gəzintisi
- **Symlink (simvolik keçid)** — fayl qısayolu
- **Unlink** — fayl/symlink silmə
- **Inode** — faylın metadata qeydi
- **Hard link** — eyni inode-a ikinci ad
- **SHA-256** — dublikat aşkarlama hash-i

## Sistem hadisələri (Ch 5)
- **Signal (siqnal)** — kernel→proses asinxron bildiriş (SIGTERM/SIGINT/SIGHUP)
- **os/signal.Notify** — siqnal kanalı
- **Cron-like scheduling** — vaxtplanlı icra
- **time.Ticker / Timer** — təkrarən / tək anlı işıq
- **Inotify** — Linux fayl hadisə subsistemi
- **fsnotify** — inotify-nin Go wrapper-i
- **File rotation (fayl rotasiyası)** — lumberjack kimi ölçü/vaxt əsaslı bölünmə
- **exec.CommandContext** — timeout idarəli prosses
- **Distributed lock (paylanmış kilid)** — konsensus əsaslı mutex

## Pipes və IPC (Ch 6)
- **IPC (Inter-Process Communication)** — proseslərarası kommunikasiya
- **Anonymous pipe (anonim boru)** — valideyn→uşaq axını (os.Pipe)
- **Named pipe / FIFO (adlı boru)** — mkfifo ilə fayl sisteminde görünən
- **os/exec** — xarici proses idarəsi
- **Pipeline (boru xətti)** — stdout→stdin zənciri
- **log processing tool** — pipe əsaslı log filtri

## Unix socketlər (Ch 7)
- **Unix domain socket (UDS)** — lokal IPC socket növü
- **SOCK_STREAM / SOCK_DGRAM** — etibarlı / datagram UDS
- **net.Listen("unix")** — UDS server
- **lsof -U** — açıq unix socketlər
- **HTTP over UDS** — reverse proxy arxitekturası
- **net/http/textproto** — HTTP başlıq emalı
- **struct{} channel** — zero-byte siqnal kanalı

## Yaddaş (Ch 8)
- **Garbage collector (zibil toplayıcı)** — avtomatik yaddaş geri qazanması
- **Stack / Heap** — funksiya frame / paylanmış yaddaş
- **Tri-color mark-sweep** — Go GC alqoritmi (write barrier + concurrent)
- **GOGC** — heap artım faizi (default 100)
- **GC pacer** — hədəf pause hesablayıcısı
- **GODEBUG=gctrace=1** — GC statistikası
- **Memory ballast** — süni heap (GC tezliyini azaltma hiyləsi)
- **GOMEMLIMIT** — yumşaq yaddaş sərhədi (Go 1.19+)
- **Arena** — manual lifecycle bölgəsi (eksperimental)

## Performans (Ch 9)
- **Escape analysis (qacış analizi)** — stack↔heap qərarı
- **`-gcflags="-m"`** — escape səbəblərini göstər
- **Benchmark** — `go test -bench` (ns/op, allocs/op)
- **b.N adaptiv dövrü** — benchmark dəqiqlik mexanizmi
- **CPU profiling** — isti funksiya profili
- **Memory profiling** — alloc profili (inuse_space vs alloc_space)
- **pprof** — profil vizuallaşdırıcı (top/list/web)
- **Premature optimization** — ölçmədən optimallaşdırma qadağı

## Şəbəkə (Ch 10)
- **net.Conn / net.Listener** — socket abstraksiyaları
- **TCP handshake** — 3-yollu bağlantı qurulması
- **HTTP verbs / status codes** — GET/POST/PUT/DELETE + 200/201/204/405
- **TLS (Transport Layer Security)** — şifrələmə + autentifikasiya
- **PEM / CRT / CSR / CN** — sertifikat formatları və axını
- **UDP / connectionless** — zəmanətsiz datagram
- **Go-Back-N / SACK (Selective Retransmissions)** — retransmit strategiyaları
- **Big-endian (şəbəkə bayt sırası)** — MSB birinci
- **WebSocket upgrade** — HTTP→uzunmüddətli ikiistiqamətli kanal

## Telemetriya (Ch 11)
- **Observability (müşahidə olunma)** — log+trace+metric üçlüyü
- **Structured logging** — JSON/key-value loglar
- **slog / zap** — standart / uber log kitabxanaları
- **runtime/trace + go tool trace** — icra trace UI
- **Distributed tracing** — trace ID + propagation + span-lar
- **Counter / Gauge / Histogram / Summary** — Prometheus metric tipləri
- **/metrics endpoint** — scrape nöqtəsi
- **rate()** — artım sürəti sorğusu
- **OTel (OpenTelemetry) / OTLP** — vendor-neutral telemetriya standartı
- **TracerProvider / WithBatcher / semconv** — OTel konfiqurasiyası
- **otelhttp** — HTTP avtomatik instrumentasiya

## Paylama (Ch 12)
- **Go modules / go.mod / go.sum** — asılılıq idarəsi + checksum
- **SemVer (major.minor.patch)** — mənalı versiya
- **Semantic import versioning** — /v2 modul yolu
- **MVS (Minimal Version Selection)** — minimal versiya alqoritmi
- **GOPRIVATE** — private modul proxy-dən çıxarış
- **insteadOf (gitconfig)** — HTTPS→SSH yönləndirmə
- **Module workspace / go.work** — çoxmodullu layihə (go 1.18+)
- **CI (Continuous Integration)** — push-da avtomatik test
- **actions/cache + hashFiles(go.sum)** — modul keşi
- **Staticcheck** — statik analiz
- **GoReleaser** — release avtomatlaşdırması
- **ldflags -s -w / -extldflags "-static" / CGO_ENABLED=0** — kiçik statik binary
- **GHCR** — GitHub Container Registry

## Capstone: Distributed cache (Ch 13)
- **CAP theorem** — consistency/availability/partition tolerance
- **Eviction policy** — TTL / LRU / FIFO
- **container/list (doubly-linked list)** — LRU sıra strukturu
- **Primary-replica / P2P / Pub-Sub / Raft-Paxos** — replikasiya strategiyaları
- **X-Replication-Request** — replikasiya dövrünün kəsilməsi
- **Sharding: range / hash / consistent hashing** — data bölgüsü
- **Hash ring** — dairəvi hash məkanı (node + key hash-ləri)
- **X-Forwarded-For** — forward loop aşkarlanması
- **Hit/miss ratio** — keş effektivliyi
- **LIRS / ARC** — inkişaf etmiş replacement alqoritmləri

## Effektiv praktikalar (Ch 14)
- **sync.Pool** — müvəqqəti obyekt hovuzu
- **sync.Once / OnceFunc / OnceValue / OnceValues** — dəfəyə icra
- **singleflight** — paralel çağırış deduplikasiyası (cache stampede dərmanı)
- **golang.org/x/sync** — Go komandasının eksperimental modul
- **mmap (memory map)** — fayl↔yaddaş xəritəsi
- **Msync / MS_SYNC / MS_ASYNC** — disk sinxronizasiyası
- **PROT_READ/WRITE/EXEC** — səhifə qoruma flaq-ləri
- **MAP_SHARED / MAP_PRIVATE** — paylaşılan / şəxsi xəritə
- **time.After leak** — dayandırılmaz timer yaddaş saxlaması
- **defer accumulation** — loop-daxili defer yığılması
- **Map non-shrinking** — delete() yaddaşı qaytarmır
- **Resource leak** — bağlanmamış fayl/socket/body
- **Channel mismanagement** — qəbulçusuz göndərən bloklanması

## Hardware avtomatlaşdırma (Appendix)
- **freedesktop.org / XDG** — desktop interop standartları
- **D-Bus** — Linux message bus (session vs system bus)
- **godbus/dbus/v5** — Go D-Bus kitabxanası
- **Match rule (type/sender/interface/path)** — siqnal filtri
- **UDisks2** — disk idarə D-Bus servisi
- **InterfacesAdded** — yeni device siqnalı
- **/proc/mounts** — kernel-in mount virtual faylı (\\040 = boşluq)
- **Partition / Block / Device / Disk** — saxlama anlayışları
- **org.freedesktop.Notifications.Notify** — sistem bildirişi
- **go-bluetooth api** — BlueZ adapter
- **RSSI** — siqnal gücü (0…-100 dBm)
- **xdg-screensaver** — ekran kilidi (Wayland-da işləmir)
