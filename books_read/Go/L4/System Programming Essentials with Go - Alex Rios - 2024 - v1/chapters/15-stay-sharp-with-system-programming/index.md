# Chapter 15 — Stay Sharp with System Programming (Sistem Proqramlaşdırmada Kəskin Qalmaq)

## Bu chapter nədən bəhs edir?
Go-un real dünya uğur hekayələri (Dropbox, HashiCorp, Grafana Labs, Docker, SoundCloud), sistem proqramlaşdırma ekosistemində necə aktual qalmaq (release notes, community, töhfə) və davamlı öyrənmə üçün klassik kitab resursları.

## Əsas fikirlər

### 1. Real dünya tətbiqləri — 5 uğur hekayəsi

**Dropbox (Python → Go miqrasiyası):**
Python backend böyüdükcə çətinləşdi — hər request üçün yeni thread resurs israfı idi. Go goroutine-lər çox yüngüldür → milyonlarla paralel istifadəçi sərbəst idarə olundu. Konkurensi modeli scalability-ni asanlaşdırdı.

**HashiCorp (birinci gündən Go):**
Terraform, Vault, Consul — hamısı Go ilə. Seçim səbəbləri:
- Sadə və öyrənməsi asan — tez produktivlik
- Kiçik, özünü-özünə yetər tətbiqlər üçün ideal (microservices)
- Zəngin standart kitabxana — çox funksional daxildə
- Paylanmış sistemlər üçün konkurensi modeli

**Grafana Labs (observability backend-i):**
Grafana frontend olsa da, backend + alətlər Go-dadır. Loki (log aqreqasiya — effektiv I/O + sıxışdırma), Tempo (distributed tracing backend — şəbəkə imkanları).
Go-un üstünlükləri (onların siyahısından): sürət (Perl/Ruby-dən qat-qat; C-dən daha sürətli inkişaf), statik binary → sadə deployment, pointer arifmetikası kimi tələlərdən uzaq balanslı azadlıq, cross-platform (linux/darwin/windows/bsd, amd64/arm; plan9, gccgo ekzotikası), güclü profiling (CPU/heap/trace — Go 1.5-dən), daxili konkurensi, güclü interface-lər.

**Docker (konteyner inqilabı):**
Lightweight isolation + portability fəlsəfəsi Go-un minimalizm + statik binary modelinə uyğun. Goroutine/channel-lər çoxlu konteynerin paralel idarəsi üçün; cross-platform (Linux/Windows/macOS) təbii. Komanda libcontainer kimi açıq kitabxanalarla Go icmasına da töhfə verdi.

**SoundCloud (Ruby monolith → Go microservices):**
"Mothership" Ruby on Rails monoliti milyonluq istifadəçi yükündə limitə çatdı → microservices arxitekturasına keçid, dil kimi Go seçildi:
- Performance + konkurensi (goroutine-lər yüksək trafik üçün)
- Sadəlik/WYSIWYG — yeni mühəndislər tez productive olur
- Sürətli compile + static typing → tez iterasiya
- Artan ekosistem və icma
Fazalı keçid: kritik olmayan komponentlərdən başladı, Bazooka (deployment platformu) kimi daxili alətlər yazdı. Nəticə: daha az resursla daha çox yük, ucuz server, dil incəliklərinə deyil, problem domeninə fokuslanmış code review.

**CNCF:** bulud-native layihələrin əksəriyyəti Go-da yazılıb (https://www.cncf.io/).

### 2. Ekosistemdə aktual qalmaq
- **Release notes:** hər Go versiyasında runtime/memory optimizations — dindarxasə izlə
- **Go blog:** https://blog.golang.org/ — rəsmi xəbərlər
- **Community:** golang-nuts, golang-dev (mailing lists), X/Twitter #golang, subreddit, Slack, GopherCon konfransı
- **Töhfə:** GitHub-da Go repo və populyar kitabxanaları izlə; kiçik bug-fix-lərdən başlayaraq open source-a contribute et
- **Eksperiment:** yeni feature/kitabxanaları əldən keçir — hands-on təcrübə əvəzolunmazdır

### 3. Sistem proqramlaşdırmanın təbiəti
Son texnologiya xəsisində deyil — fundamental qatlarda dərinləşməkdədir: OS + hardware + sistem kitabxanalarının qarşılıqlı əlaqəsi. C (və bəzən Assembly) bilmək tələb olunur; memory management, process scheduling, filesystem implementasiyası; CPU əməliyyatları, cache mexanizmləri, I/O prosesləri. Fundamental anlayışlar framework-lərdən uzunömürlüdür — köhnəlmir.

### 4. Klassik kitab resursları
| Kitab | Müəllif | Nə verir |
|---|---|---|
| **Advanced Programming in the UNIX Environment (APUE)** | W. Richard Stevens | Unix sistem proqramlaşdırmanın "bibliyası": file I/O, process control, signal, IPC |
| **Learn C Programming (2nd ed.)** | — | C-yə yumşaq giriş — sistem proqramlaşdırmada qaçınılmaz dil üçün |
| **Linux Kernel Programming (2nd ed.)** | — | Kernel daxililəri, modul yazımı, kernel sinxronizasiyası |
| **Linux System Programming Techniques** | — | Expert recipes: fork, zombie, daemon/systemd, shared library, POSIX thread, GDB/Valgrind, signal/pipe/IPC |
| **Operating Systems: Design and Implementation** | Andrew S. Tanenbaum | OS nəzəriyyə + praktika — MINIX üzərində əməli təcrübə |
| **Unix Network Programming** | W. Richard Stevens | Şəbəkə proqramlaşdırma kanonu: socket, TCP/IP, UDP, raw socket, multicast, select/poll, non-blocking I/O |
| **Mastering Embedded Linux Programming (3rd ed.)** | — | Toolchain, bootloader, kernel, rootfs; Buildroot/Yocto, Mender/Balena OTA, perf/ftrace/eBPF profiling |
| **Modern Operating Systems** | Andrew S. Tanenbaum | Müasir OS-lərin mexanikası: Windows/Linux/Unix case study-lər, distributed/multimedia/real-time |
| **The Art of UNIX Programming** | Eric S. Raymond | Unix fəlsəfəsi: kiçik, tək-şeyi-yaxşı-edən komponentlər, text stream-lər, OSS; dizayn patternləri + case study-lər |

### 5. Mentorluq
Təcrübəli Go developer-lərdən rəhbərlik al — sistem proqramlaşdırma kontekstində insight və istiqamət dəyərlidir. Bu, öyrənmə yoluğun sonuncu, amma əhəmiyyətsiz olmayan addımıdır.

## Əsas terminlər
- Goroutine vs thread — yüngül vs resurs-ağır paralellik
- Monolith → microservices miqrasiyası — fazalı keçid strategiyası
- WYSIWYG (What You See Is What You Get) — Go-un oxunaqlılıq fəlsəfəsi
- Static binary / cross-compilation — deployment sadəliyi
- CNCF (Cloud Native Computing Foundation) — Go-un bulud-native ekosistemi
- golang-nuts / golang-dev / GopherCon — icma kanalları
- APUE / Unix Network Programming — Stevens klassikləri
- MINIX — Tanenbaum-un tədris üçün yazdığı OS
- libcontainer — Docker-ın Go töhfəsi
- Loki / Tempo — Grafana Labs-ın Go alətləri
- Buildroot / Yocto Project — embedded Linux build sistemləri
- perf / ftrace / eBPF — Linux performance alətləri

## Praktik nəticə
1. Real layihələrdə Go seçimi 4 amilə görədir: konkurensi, sadəlik, statik binary, cross-platform — arxitektura qərarlarında bunları istinad kimi işlət.
2. Monolitdən microservices-ə keçiddə kritik olmayan komponentlərdən başla — riski azaldır.
3. Go release notes + blog-u izlə — runtime optimizasiyalar bilavasitə sistemin performansına təsir edir.
4. Open source-a kiçik bug-fix-lərlə başla — öyrənmə + icma kapitalı.
5. Sistem proqramlaşdırma fundamental bilikdir: Stevens (APUE, UNP) + Tanenbaum (OS) cütlüyü əsas kitabxana.
6. Mentor tap — təcrübəli rəhbərlik səhvən qaçışmaq müddətini qısaldır.

## Mənbə
Pages: 331-344 (PDF səh. 352-365)
