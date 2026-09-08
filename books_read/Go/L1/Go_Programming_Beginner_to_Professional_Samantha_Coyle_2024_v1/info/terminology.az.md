# Go Programming B2P — Terminologiya (AZ)

| İngiliscə termin | Azərbaycanca qarşılıq | Fəsil |
|---|---|---|
| Short variable declaration (:=) | qısa elan | 1 |
| Zero value | sıfır dəyər | 1 |
| Pointer / dereference / new | işarətçi / dəyərə keçid / new | 1 |
| Escape analysis (stack/heap) | yaddaş yerləşdirmə analizi | 1 |
| Shadowing | kölgələmə (ad üstünü örtmə) | 1 |
| iota | enum sayğacı | 1 |
| Initial statement | if/switch başlanğıc ifadəsi | 2 |
| Expressionless switch | ifadəsiz switch | 2 |
| fallthrough | case-dən keçid | 2 |
| Range (loop) | aralıq iterasiyası | 2 |
| Wraparound | tip həddindən daşma (max→min) | 3 |
| math/big | böyük ədədlər paketi | 3 |
| Raw literal (`) | xam string | 3 |
| rune | simvol (int32) | 3 |
| []rune vs len(string) | simvol sayı vs bayt sayı | 3 |
| Slice header (ptr/len/cap) | slice daxili quruluşu | 4 |
| Hidden/backing array | gizli (dəstəkləyən) massiv | 4 |
| make / append / copy | yaratma / böyütmə / kopyalama | 4 |
| Ok idiom (map) | mövcudluq yoxlaması | 4 |
| Struct embedding / promotion | daxiletmə / irəli çəkmə | 4 |
| Type assertion / switch | tip iddiası / seçimi | 4 |
| Variadic (pack/unpack) | dəyişən saylı (topla/aç) | 5 |
| Anonymous function / closure | anonim / bağlı funksiya | 5 |
| Function type | funksiya tipi | 5 |
| defer + FILO + freeze | təxirə salma + value dondurma | 5 |
| Named return / naked return | adlı / çılpaq qaytarma | 5 |
| Error interface / errors.New | xəta interfeysi / konstruktor | 6 |
| Panic / recover | qəza / bərpa | 6 |
| Err naming convention | Err... prefiksi | 6 |
| %w wrapping | xəta sarması | 6 |
| Implicit implementation | implisit (imzasız) tətbiq | 7 |
| Duck typing | ördək tiplemesi | 7 |
| Stringer | String() çap interfeysi | 7 |
| "Accept interfaces, return structs" | interfeys qəbul, struct qaytar | 7 |
| Type constraint / comparable | tip məhdudiyyəti / müqayisəli | 8 |
| Type inference | tip çıxarması | 8 |
| go.mod / go.sum | modul tərif / checksum faylı | 9 |
| go mod init / go get / tidy | modul əmrləri | 9 |
| Workspace (go.work) | çoxmodullu iş sahəsi | 9 |
| Exported / unexported | ixrac olunan / olunmayan | 10 |
| init() order | başlanğıc sırası | 10 |
| cmd/ + pkg/ | icra + kitabxana strukturu | 10 |
| Format verbs (%v/%T/%#v) | çap fiilləri | 11 |
| log.SetFlags (Ldate/Llongfile) | log konfiqurasiyası | 11 |
| Fatal vs Panic (log) | çıxışlı vs bərpalı log | 11 |
| time.Parse / Format / Sub | vaxt oxu / göstər / fərq | 12 |
| Duration resolutions | müddət ölçüləri | 12 |
| LoadLocation / In | saat qurşağı | 12 |
| flag package | bayraq paketi | 13 |
| Rot13 | 13-sürüşmə şifri | 13 |
| ModeCharDevice (pip detection) | pip aşkarlaması | 13 |
| Exit codes ($?) | çıxış kodları | 13 |
| SIGINT/SIGTERM/SIGKILL | dayandırma siqnalları | 13, 14 |
| signal.Notify | siqnal qeydiyyatı | 14 |
| os.Create / OpenFile (O_APPEND) | fayl yaratma / əlavə modu | 14 |
| Stat / IsNotExist | mövcudluq yoxlaması | 14 |
| //go:embed + embed.FS | binary-yə daxil etmə | 14 |
| database/sql + driver | vahid SQL API | 15 |
| SQL injection + Prepare ($1) | inyeksiya + parametrləmə | 15 |
| Query / QueryRow / Scan | oxu sorğuları | 15 |
| Exec + RowsAffected | dəyişdirici + təsir sayı | 15 |
| GORM / AutoMigrate | ORM / avtomigrasiya | 15 |
| Handler interface / ServeHTTP | işləyici interfeysi | 16 |
| Middleware | arasofta qat | 16 |
| html/template ({{.Field}}) | HTML şablonları | 16 |
| http.FileServer / StripPrefix | statik fayl xidməti | 16 |
| http.Client{Timeout} | vaxt limitli klient | 17 |
| multipart form | çoxhissəli forma (upload) | 17 |
| Goroutine (`go`) | qorutin | 18 |
| WaitGroup (Add/Done/Wait) | gözləmə qrupu | 18 |
| Race condition / atomic | yarış / bölünməz əməliyyat | 18 |
| -race flag | yarış detektoru | 18 |
| Mutex (Lock/Unlock) | qarşılıqlı istisna kilidi | 18 |
| Buffered/unbuffered channel | buferli / bufer-siz kanal | 18 |
| close(ch) / range kanal | kanal bağlama / iterasiya | 18 |
| Done channel | bitmə bildirişi | 18 |
| Pipeline / fan-out/fan-in | kəmər / yayım-yığım | 18 |
| context.WithCancel / Done() | ləğv konteksti / siqnalı | 18 |
| sync.Cond (Wait/Signal) | şərt dəyişəni | 18 |
| sync.Map (LoadOrStore) | konkurrent map | 18 |
| Table-driven test | cədvəl əsaslı test | 19 |
| Subtest (t.Run) + variable copy | alt-test + dəyişən kopyası | 19 |
| sqlmock (ExpectExec) | DB mock | 19 |
| httptest.NewServer | test HTTP serveri | 19 |
| Fuzz (f.Fuzz) | təsadüfi input testi | 19 |
| Benchmark (b.N) | performans ölçmə | 19 |
| TestMain | qlobal fixture | 19 |
| Coverage (-cover, ~80%) | kod örtüyü | 19 |
| gofmt / goimports (-w) | format / import aləti | 20 |
| go vet | statik analiz | 20 |
| go doc -all | sənəd generasiyası | 20 |
| Prometheus (counter/gauge/histogram) | monitoring metrikləri | 21 |
| /metrics + promhttp.Handler | metrik endpoint-i | 21 |
| OpenTelemetry (OTel) | izləmə sistemi | 21 |
| Tracing / span / TracerProvider | iz / aralıq / təminatçı | 21 |
| Structured logging (zap) | strukturlu jurnal | 21 |
| Multi-stage Dockerfile (scratch) | çoxmərhələli image | 21 |
| CGO_ENABLED=0 | statik build şərti | 21 |
| Kubernetes (Deployment/Service) | orkestrator resursları | 21 |
| Readiness probe / graceful shutdown | hazırlıq yoxlaması / yumşaq dayanma | 21 |
