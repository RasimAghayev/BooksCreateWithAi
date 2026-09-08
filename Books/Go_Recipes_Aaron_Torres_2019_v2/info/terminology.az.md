# Go Recipes (Torres) — Terminoloji lüğət

Format: **English Term (Azərbaycanca qarşılıq)**

## Chapter 1 — I/O və fayllar

- io.Reader (oxuma interfeysi)
- io.Writer (yazma interfeysi)
- io.MultiWriter (çoxlu yazıçı birləşməsi)
- io.Pipe (yaddaş konveyeri)
- Buffer / bytes.Buffer (yaddaş tamponu)
- Stream (məlumat axını)
- Tokenizer / Scanner (tokenləyici)
- text/template / html/template (mətn/HTML şablonu)
- FuncMap (şablon funksiyaları xəritəsi)
- Context-aware Escaping (kontekst-uyğun escape)
- Temporary File (müvəqqəti fayl)

## Chapter 2 — CLI

- Flag / FlagSet (bayraq / bayraq dəsti)
- flag.Value interfeysi (custom flag tipi)
- Positional Argument (mövqe arqumenti)
- Subcommand (alt komanda)
- Environment Variable (mühit dəyişəni)
- Marshal / Unmarshal (seriyalaşdırma)
- Unix Pipe (boru)
- Signal / SIGINT / SIGTERM (siqnal)
- Graceful Shutdown (təmiz bağlanma)
- ANSI Escape Codes (rəng kodları)

## Chapter 3 — Data transformasiyası

- Type Conversion / Type Assertion (tip çevirməsi / təsdiqi)
- Type Switch (tip seçimi)
- comma-ok idiom (səhvsiz tip yoxlaması)
- big.Int (dəyişən uzunluqlu tam ədəd)
- Memoization (nəticə yaddaşları)
- SQL NullTypes (sql.NullInt64 və s.)
- gob Encoding (Go-nun binary formatı)
- Base64 (Std/URLEncoding)
- Struct Tag (struktur etiketi)
- Reflection (refleksiya — reflect paketi)
- Closure (qapanış)
- Map / Filter (funksional əməliyyatlar)

## Chapter 4 — Xəta emalı

- Error Interface (xəta interfeysi)
- Error Value / Error Type (xəta dəyəri/tipi)
- Error Wrapping (xəta örtülməsi)
- errors.Cause (orijinal xəta)
- Stack Trace (stek izi)
- Log Levels (səviyyələr)
- Structured Logging (strukturlaşdırılmış logging)
- Hook / Handler / Formatter (log komponentləri)
- sync.Once (bir dəfəlik icra)
- Panic / Recover (panik / bərpa)

## Chapter 5 — Şəbəkə

- TCP/IP / UDP (protokollar)
- Broadcast (yayım)
- net.Listen / net.Dial / net.Accept (bağlantı API)
- DNS / CNAME (domen sistemi)
- WebSocket / Upgrade (iki tərəfli bağlantı)
- RPC (Remote Procedure Call)
- RFC5322 (e-poçt standartı)
- Stringer interfeysi (String() metodu)

## Chapter 6 — Bazalar

- database/sql / Driver (SQL sürücüsü)
- Blank Import (yalnız side-effect importu)
- Connection Pool (bağlantı pulu)
- MaxOpenConns / MaxIdleConns (pul limitləri)
- Transaction / Commit / Rollback (transaksiya)
- BeginTx (context-aware transaksiya)
- TTL (Time To Live — yaşam müddəti)
- redis.Nil (açar tapılmadı xətası)
- BSON (MongoDB binary JSON)
- Storage Interface (saxlanc interfeysi)

## Chapter 7 — Veb klientlər

- http.Client / http.DefaultClient
- Transport / RoundTripper (nəqliyyat qatı)
- NopTransport (test transportu)
- Async / Parallel sorğular
- Buffered Channel (buferli kanal)
- OAuth2 / Authorization Code / Token Exchange
- Refresh Token (yeniləmə tokeni)
- TokenSource / Token Storage
- Decorator / Middleware (dekorator / araqat)
- gRPC / Protobuf / protoc
- unary RPC (tək çağırış)
- twirp RPC (HTTP/1.1 üzərində RPC)

## Chapter 8 — Mikroservislər

- Handler / HandlerFunc / ServeHTTP
- ResponseWriter / Request
- Dependency Injection (asılılıq inyeksiyası)
- Closure Handler (qapanış handler)
- Validation / Typed Error (Verror)
- Content Negotiation (məzmun danışığı)
- Middleware zənciri (ApplyMiddleware)
- context.WithValue / r.WithContext
- Reverse Proxy (tərs proksi)
- Director (proksi yönləndiricisi)
- RequestURI (proksidə boşaldılmalı)

## Chapter 9 — Test

- Mock / Stub / Test Double (sınaq ikiqatı)
- Patch / Restore (əvəzetmə / geri qaytarma)
- gomock / mockgen / EXPECT (mock generasiyası)
- Table-driven Tests (cədvəlli testlər)
- Subtest (t.Run)
- Test Coverage (örtüklük) / coverprofile
- gotests (test generatoru)
- Assertion (iddia) / Convey / So
- BDD / Gherkin / Cucumber / godog
- httptest.NewRecorder (ResponseWriter mock)

## Chapter 10 — Concurrency

- Goroutine (qo-rotin)
- Channel / Buffered Channel (kanal)
- select operatoru
- time.Tick (taymer — selectdə sıfırlanma tələsi)
- sync.WaitGroup (gözləmə qrupu)
- Race Condition (yarış şəraiti)
- Mutex / RWMutex (qarşılıqlı istisna)
- sync/atomic (Store/Load/Add/CAS)
- Quorum ((n/2)+1 səs)
- Worker Pool (işçi hovuzu)
- Pipeline (konveyer)
- State Channel (vəziyyət kanalı — WorkRequest/WorkResponse)

## Chapter 11 — Distributiv sistemlər

- Service Discovery (servis kəşfi)
- Consul / Health Check
- Consensus / Raft (konsensus alqoritmi)
- Leader Election (lider seçkisi)
- FSM (Finite State Machine)
- In-memory Transport (yaddaş nəqliyyatı)
- Dockerfile / ENTRYPOINT / EXPOSE
- ldflags -X (build-time dəyişənlər)
- Docker Compose / vendor (offline build)
- Prometheus / Pull Model / scrape_interval
- Registry / Counter / Timer / Percentile (metric-lər)

## Chapter 12 — Reaktiv

- Reactive Programming (reaktiv proqramlaşdırma)
- Dataflow / Graph (data axını / qraf)
- Goflow Component / Connect / MapInPort
- Kafka: Topic / Partition / Offset
- SyncProducer / AsyncProducer
- Successes / Errors kanalları
- OffsetNewest (ən yeni mesajlardan)
- GraphQL: Schema / Field / Resolve / RootQuery / Args
- Overfetching (artıq data çəkilməsi)

## Chapter 13 — Serverless

- Serverless (serversiz arxitektura)
- AWS Lambda / Handler imza qaydaları
- Apex (deploy/invoke/logs/rollback)
- IAM Role / Policy
- CloudWatch (Lambda logları)
- Billed Duration (hesablanan müddət)
- App Engine / Datastore (IncompleteKey/Query/Order)
- Firebase / Firestore / Collection / Document
- Service Account (xidmət hesabı tokeni)
- Eventual Consistency (sonrakı tutarlılıq)

## Chapter 14 — Performans

- pprof / net/http/pprof (profilinq)
- CPU/Heap Profile
- Hot Path (isti yol)
- Benchmark / b.N / RunParallel
- Lock Contention (kilid rəqabəti)
- -benchmem / B/op / allocs/op
- String Concatenation vs strings.Join
- fasthttp / fasthttprouter / RequestCtx
- UserValue (route parametri)
- HTTP/2 (fasthttp-da dəstəklənmir)
