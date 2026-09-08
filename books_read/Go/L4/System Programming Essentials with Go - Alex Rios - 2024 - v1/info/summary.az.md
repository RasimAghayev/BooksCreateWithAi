# System Programming Essentials with Go — Xülasə (Azərbaycanca)

> **Kitab:** System Programming Essentials with Go — Alex Rios, Packt, 2024 (ISBN 978-1-83512-772-8) · 372 səh.
> 🚀 Advanced (4/5) · Linux-mərkəzli praktik yol: syscall-lardan paylanmış keşə qədər tam sistem proqramlaşdırma dövrü.

---

## Kitabın yanaşması
5 hissəli: (1) giriş — Go seçimi + konkurensi refresh; (2) OS interaksiyası — syscall, fayllar, hadisələr, pipe, UDS; (3) performans — GC, escape, benchmark, profiling; (4) bağlı tətbiqlər — şəbəkə, telemetriya, paylama; (5) irəli — capstone (distributed cache), effektiv praktikalar, real dünya. Hər chapter-də trade-off analizi; Linux üstünlüklü (strace, inotify, D-Bus). 17 chapter + hardware avtomatlaşdırma appendix-i.

## Fəsl-fəsil xülasə

### Ch 1 — Niyə Go?
Go-un sistem proqramlaşdırmaya uyğunluğu: goroutine konkurensi (CSP), OS tooling (build/test/run/vet/fmt), cross-platform compile. Sistem səviyyəli iş üçün C-dən təmiz, script dillərindən sürətli orta yol.

### Ch 2 — Konkurensi refresh
Goroutine + WaitGroup; data race və həlli (atomic/mutex); channel dərinliyi — unbuffered (senxron, delivery zəmanəti) vs buffered (asinxron, itirmə riski); state/signaling üçün kanal seçim meyarları; sinxronizasiya aləti seçim cədvəli.

### Ch 3 — Sistem çağırışları
Syscall mexanizmi (user→kernel mode), 3 API səviyyəsi (syscall/x/sys/os); strace ilə tracing; standart stream-lər + fayl deskriptorları; redirection; CLI tətbiqi + os.Pipe ilə testability.

### Ch 4 — Fayl əməliyyatları
Permission bitləri + umask; təhlükəsiz permission aşkarlanması; filepath paketi (WalkDir, Ext, Join); symlink/unlink; qovluq ölçüsü hesablama; hash ilə dublikat fayl tapma; filesystem optimizasiyası (directoryWalk, parallel read).

### Ch 5 — Sistem hadisələri
Siqnallar (os/signal.Notify, graceful shutdown); task scheduling; timer siqnalları; fayl monitorinqi — inotify mexanizmi + fsnotify wrapper; fayl rotasiyası; proses idarəsi (exec.CommandContext timeout); paylanmış lock manager qurulumu.

### Ch 6 — Pipes (IPC)
Anonymous pipe mexanikası (os.Pipe, fork-yönümdə valideyn/uşaq); named pipe (mkfifo) — fayl sistemində görünən, byte-stream, bloklanma xüsusiyyətləri; pipe best practice-ləri (error handling, resurs, security); log processing aləti inkişafı.

### Ch 7 — Unix socketlər
UDS yaratma (Listen/Dial "unix"); socket tipləri (STREAM/DGRAM); lsof ilə inspeksiya; chat server (çoxlu client, goroutine per-conn); HTTP over UDS (custom Transport DialContext); performans müqayisəsi (TCP-yə qarşı local IPC üstünlüyü).

### Ch 8 — Yaddaş idarəsi
GC alqoritmi (tri-color mark-sweep, concurrent); stack vs heap bölgüsü; GOGC tuning; GC pacer; GODEBUG tracing; memory ballast texnikası; GOMEMLIMIT (Go 1.19+); arena eksperimenti.

### Ch 9 — Performans analizi
Escape analysis (-gcflags=-m; pointer/stack/heap qaydaları); benchmark yazımı (b.N, -benchmem, alloc metricaları); CPU profiling (pprof top/list/web); memory profiling (inuse/alloc); profil-over-time; trade-off framework-i (hər optimallaşdırma ölçmə ilə).

### Ch 10 — Şəbəkə
net paketi; TCP echo server (goroutine per-conn, buf[:n]); HTTP server/client, verb-switch, status kodları; TLS (sertifikat yaratma, PEM/CRT, CSR, ListenAndServeTLS, tls.Listen); UDP + custom Selective Retransmissions (seq/ACK, BigEndian); WebSocket (gobwas/ws, upgrade).

### Ch 11 — Telemetriya
Struktur loglama (slog vs zap, JSON seçim meyarları, nə logla/nə loglama); tracing (runtime/trace, go tool trace, HTTP TraceHandler, distributed tracing 4 konsepti); metrics (Prometheus, 4 metric tipi, /metrics, rate()); OTel (vendor-neutral, otlptracehttp, semconv, otelhttp).

### Ch 12 — Paylama
Go modules (SemVer, MVS, go.sum checksum); private repo (GOPRIVATE + SSH insteadOf); go install@version; module workspaces (go.work, go work sync); CI (GitHub Actions: test + modul cache + Staticcheck); release (GoReleaser: cross-compile, statik binary, arxivlər, GHCR docker, tag-triggered workflow).

### Ch 13 — Capstone: Distributed cache
Tələblər → dizayn trade-off-ları: thread safety (RWMutex seçimi), interface (HTTP vs TCP/gRPC — sadəlik/standart/stateless), eviction (TTL ticker + LRU container/list), replikasiya (P2P — SPOF yoxdur; X-Replication-Request ilə dövr qoruması), sharding (consistent hashing — hash ring, SHA-1, binary search; X-Forwarded-For loop qoruması). İrəliləyiş: LIRS/ARC, compression, connection pooling, metrics+profiling.

### Ch 14 — Effektiv praktikalar
sync.Pool (3 ssenari: buffer, şəbəkə server, JSON marshal); sync.Once + Go 1.21 OnceFunc/OnceValue/OnceValues; singleflight (cache stampede, dedup, throttling); mmap (x/exp/mmap, Msync, PROT/MAP flaq-ləri); 6 performans tələsi: time.After leak, defer-in-loop, map kiçilməməsi, bağlanmayan resurslar, HTTP body, channel idarəsizliyi.

### Ch 15 — Kəskin qalmaq
Real dünya: Dropbox (Python→Go miqrasiya), HashiCorp (birinci gündən Go), Grafana Labs (Loki/Tempo), Docker (libcontainer), SoundCloud (Ruby monolith→Go microservices); CNCF ekosistemi. aktual qalma: release notes, icma, töhfə. 9 klassik kitab tövsiyəsi (Stevens, Tanenbaum, Raymond).

### Appendix — Hardware avtomatlaşdırma
USB: /proc/mounts oxuma, uzantı-əsaslı fayl təşkilatçısı; D-Bus (session vs system bus, UDisks2 InterfacesAdded, match rule, Properties.GetAll MountPoints); sistem bildirişləri (Notifications.Notify spec); Bluetooth: RSSI (-70 dBm threshold) ilə smartwatch-a əsaslanan ekran kilidi; XDG/Wayland məhdudiyyəti.

## Kitabın əsas mesajları
1. **Sistem proqramlaşdırma = trade-off idarəsi:** hər qərar (RWMutex vs sync.Map, TCP vs HTTP, LRU vs TTL, P2P vs consensus) ölçülmüş şərtlərə əsaslanır.
2. **Alçaq səviyyə bilik portativdir:** syscall mexanizmi, pipe, UDS, mmap anlayışları Go-da olsa da Unix fundamentalına keçir.
3. **Performans ölçməsiz iddia deyil:** escape analysis, benchmark, pprof — hər optimallaşdırmanın öncül şərti.
4. **Resurs intizamı:** hər açılan resurs defer Close; time.After/defer-loop/map/channel tələləri production killer-dır.
5. **Vendor-neutral telemetriya:** OTel ilə backend dəyişir, instrumentasiya qalır.
6. **Distribution = modullar + CI + avtomatik release:** MVS, go.sum, workspaces, GoReleaser zənciri.

## Kitabdan sonra
- Kitabda qeyd olunan klassiklər: APUE, Unix Network Programming (Stevens), Modern Operating Systems (Tanenbaum)
- Distributed cache-i genişləndir: LIRS/ARC, compression, ristretto müqayisəsi
- Linux daxililəri: kernel modulları, eBPF (Mastering Embedded Linux tövsiyəsi)
- go101.org optimallaşdırmalar kitabı ilə GC/escape dərinləşməsi (bizim kitab 10)
