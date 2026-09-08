# Effective Go Recipes — Terminologiya (AZ)

| İngiliscə termin | Azərbaycanca qarşılıq | Fəsil |
|---|---|---|
| bytes.NewReader | bayt oxuyucu (yaddaş içi) | 1 |
| gzip.NewWriter/Reader | sıxışdırma qatı | 1 |
| fs.FS / os.DirFS | fayl sistemi abstraksiyası | 1 |
| io.Pipe (os.Pipe) | Unix borusu (Reader/Writer) | 1 |
| mmap (unix.Mmap) | yaddaşa xəritələnmiş fayl | 1 |
| JSON lines | sətir-sətir JSON | 1, 2 |
| encoding/gob | Go-özəl binary format | 2 |
| Anonymous struct | anonim strukturdur | 2 |
| Zero value problem | sıfır dəyər ikiliyi | 2, 6 |
| Pointer field detection | pointer-lə itkin sahə yoxlaması | 2 |
| json.Marshaler/Unmarshaler | xüsusi serializasiya interfeysi | 2 |
| mapstructure | map→struct çevirici | 2 |
| URL query (url.Values) | təhlükəsiz URL qurumu | 3 |
| io.LimitReader | oxu həddi (DoS qoruması) | 3 |
| Chunked transfer encoding | ölçüsüz HTTP axını | 3 |
| Middleware (HTTP) | handler sargısı | 3 |
| %#v verb | tip-göstərən çap | 4 |
| fmt.Stringer | String() custom çap | 4 |
| strings.EqualFold | Unicode case-insensitive müqayisə | 4 |
| Unicode normalization (NFC/NFD/NFKC/NFKD) | unikod normal forması | 4 |
| First-class functions | funksiya = dəyər | 5 |
| Function registry (dispatch) | açar→funksiya xəritəsi | 5 |
| Functional options | WithX seçim funksiyaları | 5 |
| Closure | bağlı funksiya (mühit yaddaşı) | 5 |
| go:linkname | unexported çıxış direktivi | 5 |
| Comma, ok paradigm | ikili qəbul (map/kanal/assertion) | 6 |
| Slice header (array/len/cap) | slice daxili strukturu | 6 |
| append growth policy | böyümə siyasəti (×2 → 1.3x) | 6 |
| Composite key | struct map açarı | 6 |
| time.LoadLocation | saat qurşağı yükləmə | 6 |
| Epoch / Y2038 | 1970 bazası / 32-bit daşması | 6 |
| Ad hoc interface | kiçik məqsədli interfeys (syncer) | 7 |
| NOP (nopSyncer) | heç-nə-edən tip | 7 |
| Type constraint (~int) | generic məhdudiyyət | 7 |
| var zero T | generic sıfır dəyər | 7 |
| Accept interfaces, return types | idiomatik API qaydası | 7 |
| %w (wrap) | xəta sargısı | 8 |
| Named return values | adlı qaytarma (panic üçün) | 8 |
| recover (deferred) | panic tutucu | 8 |
| Fail fast | tez çökmə fəlsəfəsi | 8 |
| safelyGo | goroutine panic qoruması | 8 |
| errors.Is / errors.As | zəncir yoxlaması | 8 |
| runtime.Callers / CallersFrames | stack PC → frame | 8 |
| Fan-out / Fan-in | yayım / yığılma | 9 |
| Amdahl's law | paralel speedup limiti | 9 |
| Buffered channel semaphore | goroutine həddi | 9 |
| Worker pool | işçi dəsti | 9 |
| Goroutine leak | əbədi bloklanmış goroutine | 9 |
| errors.Join | xəta birləşdirmə | 9 |
| ctx.Done() | ləğv siqnalı | 9 |
| ctxKey idiomu | toqquşmasız kontekst açarı | 9 |
| sync.Once | birdəfəlik icra | 10 |
| sync.WaitGroup | qrup gözləməsi | 10 |
| sync.RWMutex | oxuyan/yazan kilidi | 10 |
| Race detector (-race) | yarış yoxlayıcısı | 10 |
| Memory barrier | cache sinxronu | 10 |
| atomic.Value | atomik saxlanc (cast ilə) | 10 |
| net.Listen / Accept | TCP qəbulu | 11 |
| io.CopyN | dəqiq N bayt kopya | 11 |
| Unix domain socket | lokal IPC | 11 |
| binary.BigEndian | şəbəkə bayt sırası | 11 |
| Fallacies of distributed computing | paylanmış sistem yanlış inancları | 11 |
| os/exec | xarici əmr icrası | 12 |
| StdinPipe / StdoutPipe | proses I/O körpüləri | 12 |
| cgo / import "C" | C inteqrasiyası | 12 |
| C.CString / C.CBytes / C.GoBytes | tip çevrilmələri | 12 |
| "Cgo is not Go" | cgo xəbərdarlığı | 12 |
| t.Skip / t.Fatal | test status idarəsi | 13 |
| Fuzzing / heuristic | random input test / uğur şərti | 13 |
| http.RoundTripper (mock) | transport mock nöqtəsi | 13 |
| TestMain / fixtures | qlobal setup/teardown | 13 |
| t.Cleanup / t.TempDir | avtomatik təmizlik | 13 |
| go/analysis (linter) | kod analiz framework-u | 13 |
| //go:embed | binary-yə fayl daxiletmə | 14 |
| embed.FS | yaddaş fayl sistemi | 14 |
| ldflags -X | compile-vaxtı inyeksiya | 14 |
| CGO_ENABLED=0 | statik build | 14 |
| Build tags (//go:build) | şərti kompilyasiya | 14 |
| goreleaser | release avtomatlaşdırma | 14 |
| //go:generate | kod generasiya direktivi | 14 |
| Rule of Generation | "proqram yazan proqram yaz" | 14 |
| conf.Parse | env/CLI konfiqurasiya | 15 |
| replace direktivi | asılılıq patch | 15 |
| Multistage Docker | build+slim deploy | 15 |
| Graceful shutdown (srv.Shutdown) | yumşaq dayanma | 15 |
| signal.Notify (SIGTERM/SIGINT) | siqnal tutma | 15 |
| zap.Check | səviyyə yoxlaması (bahalı parametr) | 15 |
| expvar (/debug/vars) | built-in metrics | 15 |
| Delve (dlv attach) | canlı proses debugger | 15 |
