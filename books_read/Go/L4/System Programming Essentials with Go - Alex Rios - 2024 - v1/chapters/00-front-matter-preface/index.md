# Front Matter / Preface (Ön Sözlər)

## Bu chapter nədən bəhs edir?
Kitabın məqsədi, hədəf auditoriyası, chapter-by-chapter xülasəsi, texniki tələblər (Go 1.16+, Linux üstünlük), kod konvensiyaları və müəllif haqqında.

## Kitabın kimə lazımdır?
- Software engineer, architect, developer — sistem dizayn biliklərini dərinləşdirmək istəyənlər
- İş yerində kompleks dizayn problemlərini həll edənlər
- Low-level proqramlaşdırmaya maraqlı olanlar
- Tələb: proqramlaşdırma əsasları + ən azı bir dil təcrübəsi (Go əsasları burada təkrar olunmur, amma Chapter 2 konkurensini refresh edir)

## Kitabın strukturu (5 hissə)

### Part 1: Introduction
- **Ch 1 Why Go?** — Go-un sistem proqramlaşdırmaya uyğunluğu: konkurensi modeli, networking/I/O, low-level nəzarət, syscall-lar, cross-platform, tooling (go build/test/run/vet/fmt)
- **Ch 2 Refreshing Concurrency and Parallelism** — goroutine, WaitGroup, data race idarəsi (atomic/mutex), channel-lər (buffered/unbuffered, delivery zəmanəti, latency, state/signaling)

### Part 2: Interaction with the OS
- **Ch 3 Understanding System Calls** — syscall kataloqu, os/syscall/x/sys paketləri, portabilitet, syscall tracing, standart stream-lər, fayl deskriptorları, CLI tətbiqi + testability
- **Ch 4 File and Directory Operations** — təhlükəsiz permission aşkarlanması, filepath, qovluq gəzinti, symlink, unlink, qovluq ölçüsü, dublikat fayl tapmaq, filesystem optimizasiyası
- **Ch 5 Working with System Events** — signal-lar (os/signal), task scheduling, timer, fayl monitorinqi (inotify/fsnotify), fayl rotasiyası, proses idarəsi (timeout), distributed lock manager
- **Ch 6 Pipes in IPC** — anonymous pipe mexanikası, named pipe (mkfifo), pipe best practice-ləri, log processing aləti
- **Ch 7 Unix Sockets** — socket yaradılması, lsof ilə inspeksiya, chat server, HTTP over UNIX domain socket, performans

### Part 3: Performance
- **Ch 8 Memory Management** — GC alqoritmi, stack/heap, GOGC, GC pacer, GODEBUG, memory ballast, GOMEMLIMIT, arena-lar
- **Ch 9 Analyzing Performance** — escape analysis, benchmark yazımı, alloc metricaları, CPU/memory profiling, trade-off hazırlığı

### Part 4: Connected Apps
- **Ch 10 Networking** — net paketi, TCP/UDP, HTTP server/client, verb/status, TLS sertifikatları
- **Ch 11 Telemetry** — logs (slog/zap), traces, metrics (Prometheus), OTel
- **Ch 12 Distributing Your Apps** — Go modules, CI (GitHub Actions + cache + Staticcheck), GoReleaser release

### Part 5: Going Beyond
- **Ch 13 Capstone: Distributed Cache** — Memcached/Redis tipli sistem: sharding, eviction, replikasiya, trade-off-lar
- **Ch 14 Effective Coding Practices** — sync.Pool, sync.Once, singleflight, mmap, performans tələləri
- **Ch 15 Stay Sharp** — real dünya case study-ləri + davamlı öyrənmə resursları
- **Appendix: Hardware Automation** — USB/Bluetooth avtomatlaşdırması, D-Bus, XDG

## Texniki tələblər
| Tələb | Detal |
|---|---|
| Golang | 1.16+ |
| OS | Windows, macOS, Linux (preferably **Linux**) |

Kod nümunələri: https://github.com/PacktPublishing/System-Programming-Essentials-with-Go

## Müəllif haqqında
**Alex Rios** — Braziliyalı software engineer, 15 il təcrübə; Go ixtisası ilə yüksək-throughput sistemlər (fintech, telecom, gaming). Stone Co.-da staff engineer. Curitiba Go meetup təşkilatçısı, GopherCon Brazil spikeri. Bu, onun ilk kitabıdır (Elton Minetto-nun təşviqi ilə). Technical reviewer: **Natan Streppel** (2017-dən Go, distributed systems, open source).

## Nəşr məlumatları
- Packt Publishing, iyun 2024
- ISBN: 978-1-83763-413-2
- 372 səhifə, 17 chapter + appendix

## Praktik nəticə
1. Kitab Linux mərkəzlidir (syscall/inotify/UDisks2/D-Bus) — Windows-dan izləyənlər üçün WSL/Linux VM tövsiyə olunur.
2. Oxu ardıcıllığı təbii axın izləyir: OS interaksiyası → performans → şəbəkə → paylanmış sistem — amma hər chapter müstəqil də istifadə oluna bilər.
3. Capstone (Ch 13) bütün bilikləri birləşdirir — ora qədər praktikaları yerinə yetirmək faydalıdır.

## Mənbə
Pages: xiii-xx (PDF səh. 2-21)
