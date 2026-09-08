# System Programming Essentials with Go — Müəllim Qeydləri (Teacher Notes)

> Bu sənəd kitabı öyrədərkən istifadə üçün metodiki qeydlər, çətin anlar, müzakirə sualları və praktik tapşırıqlar.

---

## 1. Auditoriya və ön şərtlər

- **Hədəf auditoriya:** Go əsaslarını bilən, OS anlayışlarına maraqlı intermediate developer-lər.
- **Ön şərtlər:** Go sintaksisi + konkurensi əsasları, Linux command line, əsas şəbəkə anlayışları (port, HTTP).
- **Həcm:** 372 səh., 15 fəsil + appendix — 4 aylıq kurs (həftədə 1 fəsil + lab) və ya 2 günlük intensiv mümkün deyil; minimum 5 günlük intensiv seçdi hissələr üçün.
- **Mühit:** Linux (preferably) — strace/inotify/D-Bus bölmələri Windows-da işləmir; WSL2 alternativi.

## 2. Tədris axını — 5 hissə paralelində

| Həftə | Chapter | Lab |
|-------|---------|-----|
| 1 | 1-2 | `-race` ilə data race demo; channel tipli quiz |
| 2 | 3 | strace ilə syscall sayı; CLI cat/echo tətbiqi + pipe testləri |
| 3 | 4 | Uzantı əsaslı fayl təşkilatçısı + dublikat tapıcı |
| 4 | 5 | fsnotify watcher + SIGHUP ilə config reload |
| 5 | 6 | Log processing pipeline (anonim + named pipe) |
| 6 | 7 | UDS chat server + lsof inspeksiya |
| 7 | 8 | GOGC/GOMEMLIMIT təcrübələri; gctrace oxunuşu |
| 8 | 9 | -m escape analizi; benchmark + pprof sessiyası |
| 9 | 10 | TCP echo → HTTP CRUD → TLS server; UDP SACK demo |
| 10 | 11 | slog→zap köçürmə; Prometheus counter + OTel span |
| 11 | 12 | go.work iki modul; CI workflow + GoReleaser config |
| 12 | 13 | Capstone: keş hissə-hissə (RWMutex→TTL→LRU→replika→shard) |
| 13 | 14 | sync.Pool benchmark; tələ təkrarları (time.After, defer loop) |
| 14 | 15+Appendix | RSSI lock və ya D-Bus USB təşkilatçısı layihəsi |

## 3. Çətin anlar və izah üsulları

### a) Buffered channel itirir (Ch 2)
Demo: 3-elementli buffered kanala 5 göndərmə — kod bloklanmır, amma buffer dolub-daşır. "Sənət poçtu vs dəbdəbəli restoran" metaforası: unbuffered = garson masanı gözləyir.

### b) Syscall qatmanları (Ch 3)
Diaqram: app → os → x/sys → syscall → kernel. `os.Open`-i strace ilə izlə — eyni openat syscall-ı görünür. "os = tərcüməçi, syscall = xam imza" müqayisəsi.

### c) Inotify hadisə axını (Ch 5)
fsnotify-də WRITE hadisəsinin təkrar gəlməsi şagirdləri çaşdırır — `echo x > f` = CREATE+WRITE. Hadisə axınını terminalda canlı izləyin (inotifywait -m).

### d) Named pipe bloklanması (Ch 6)
mkfifo açılan proses yazma üçün bloklanır — qiraətçi gələnə qədər. İki terminalda canlı demo: writer blokda "donur", reader açılanda axır. "Telefon xətti: hər iki tərəf gözləyir".

### e) Ballast əleyhinə intuisiya (Ch 8)
"1GB boş slice niyə yaddaşı yandırmır?" — virtual yaddaş, səhifələr toxunulmur. pprof heap qrafikində ballast görünür, amma RSS yox. Go 1.19+ GOMEMLIMIT ballast-ı köhnəlir — tarixi kontekst verin.

### f) Escape qaydaları (Ch 9)
3 nümunə serialı: lokal dəyər (stack) → funksiyadan pointer qaytaran (heap) → interface-ə salınan (heap). `-gcflags="-m"` output-unu birgə oxuyun — compiler-in "səbəb" şərhləri ən yaxşı müəllimdir.

### g) Consistent hashing halqası (Ch 13)
Whiteboard: dairə + 3 node + 6 key. Node B çıxarılır — yalnız B-nin aralığındakı key-lər C-yə keçir. Sonra sadə hash % n müqayisəsi: node sayı dəyişəndə HƏR ŞEY yerindən oynayır. "Kütikdə oturmaq vs halqada" metaforası.

### h) Replikasiya dövrü (Ch 13)
X-Replication-Request header-siz versiyanı canlı işə salın — iki node sonsuz bir-birinə POST göndərir (log seli). Sonra header əlavəsi ilə dayandırma — dərs şagirdin gözü önündə "partlayır".

### i) time.After leak (Ch 14)
1000 goroutine × 10 dəqiqəlik time.After — heap qrafikində timerlərin yığılması. NewTimer+Stop versiyası ilə müqayisə. "Sığorta müqaviləsini ləğv etmək" metaforası.

### j) /proc/mounts \\040 (Appendix)
"My Drive" adlı USB — mount sətrində /media/My\040Drive. Parse edib fayl tapa bilməyən kodun debug-i canlı göstərin: "kernel boşluğu saxtalaşdırır, çünki sahə ayırıcı space-dir".

## 4. Müzakirə sualları

1. sync.Mutex vs RWMutex vs sync.Map: oxu/yaz nisbəti 9:1 olan keş üçün hansı? 1:1 üçün?
2. Siqnal handler-da mətn yazmaq niyə təhlükəlidir? (async-signal-safety)
3. Pipe vs UDS: hansı IPC hansı ssenaridə? (valideyn-uşaq vs müstəqil proseslər)
4. GOGC=50 vs GOGC=400: CPU xərci vs yaddaş — hansı trade-off?
5. TCP keş interfeysi üçün HTTP-dən 5x sürətli ola bilər — kitab niyə HTTP seçir? Nə vaxt dəyişmək lazımdır?
6. LRU keşdə Get niyə yazma kilidi tələb edir? (MoveToFront = sıra dəyişməsi)
7. P2P replikasiyada iki node eyni anda eyni açarı fərqli yazsa nə olur? (konflikt həlli yoxdur — eventual consistency)
8. time.After-i select-default-da işlətsək hər iteration-da yeni timer yaranır — bu leak nə vaxt real problem olur?
9. RSSI -70 niyə "uzaq" deməkdir, -70 dBm nə vaxt "yaxın" ola bilər? (kalibrasiya, interferensiya)
10. Wayland-də xdg-screensaver işləmir — cross-desktop avtomatlaşdırma üçün nə etməli?

## 5. Praktik tapşırıqlar (lab)

1. **Ch 2:** data-race'li counter yazın (2 goroutine × 1000 artım) — `-race` ilə tutun, atomic + mutex ilə düzəldin; hər 3 versiyanı benchmark edin.
2. **Ch 3:** strace -c ilə "Hello World" proqramının kaç syscall etdiyini hesablayın; stdini pipe ilə əvəz edən test yazın.
3. **Ch 4:** 0777 permission'lı faylları tapıb xəbərdarlıq edən auditor; dublikatları SHA-256 ilə siyahılayan alət.
4. **Ch 5:** fsnotify ilə config watcher — dəyişiklikdə SIGHUP gözləyən worker (graceful reload simulyasiyası).
5. **Ch 6:** mkfifo + iki proses: writer JSON logları yazır, reader filtrləyib stdout-a.
6. **Ch 7:** UDS echo server + klient; lsof ilə socket tipini göstərin; TCP versiyası ilə benchmark müqayisəsi.
7. **Ch 8:** GOGC=100 vs 400 ilə eyni benchmark — GC sayı və pause müqayisəsi (gctrace).
8. **Ch 9:** 3 versiyalı funksiya (stack/heap/pointer) — -m output + benchmark alloc fərqi.
9. **Ch 10:** TLS echo server: openssl ilə self-signed, Go client ilə qoşulma; Wireshark-da şifrələnmiş axını görün.
10. **Ch 11:** Prometheus counter + rate() sorğusu; OTel span ilə 2-layiqli trace (server→DB simulyasiya).
11. **Ch 13:** Capstone hissə-hissə: 1) map+RWMutex 2) TTL ticker 3) LRU list 4) 2-node replikasiya 5) hash ring. Hər mərhələdə test.
12. **Ch 14:** sync.Pool'lu vs pool'suz JSON marshal benchmark; 1000 fayllıq defer-in-loop tələsini ölçün.
13. **Appendix:** /proc/mounts parse edən fayllisteleme; genişlənmə: D-Bus InterfacesAdded + avtomatik təşkilatçısı.

## 6. Qiymətləndirmə

- **Capstone (Ch 13) 40%:** işlək distributed cache — thread-safe, TTL+LRU, replikasiya, sharding; testlər mövcud.
- **Labs 30%** · **Final imtahan 20%** (trade-off sualları) · **İştirak 10%** (müzakirələr).
- Bonus: capstone-a singleflight cache stampede qoruması əlavə edən şagird.

## 7. Müəllifin tərzi haqqında qeyd

Kitab ağır ironiya/sarkazm ilə yazılıb ("herding cats", "bureaucracy") — tədrisdə bu nümunələri sərtləşdirin; texniki məzmun dəqiq və praktikdir. Hər chapter-in "Technical requirements" + GitHub kodu (PacktPublishing/System-Programming-Essentials-with-Go) var — kodu oxutmaq müzakirə üçün əla materialdır.
