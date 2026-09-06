# Go in Practice — Terminologiya lüğəti (AZ)

`English (Azərbaycanca)` formatında — SYSTEM_PROMPT qayda 4 üzrə.

## CLI
- Plan 9 flag üslubu (tək dash; qruplaşdırma yox)
- flag.String / flag.BoolVar / flag.Parse / PrintDefaults / VisitAll
- gnuflag / go-flags (struct tag reflection)
- cli.go / Commands / ShortName / cli.Context / GlobalString / NewExitError
- Variadic / args...

## Konfiqurasiya
- 12-Factor App (konfiq ENV-də; 12 amil)
- gcfg (INI) / go-gypsy (YAML)
- etcd (distributed konfiq; Raft)
- Ansible / Chef / Puppet (ops platformaları)

## Web server
- Init Daemon (systemd/upstart/init/launchd)
- Callback Shutdown URL (antipattern)
- Graceful Shutdown / manners
- signal.Notify / os.Interrupt / os.Kill
- Zero-downtime deploy
- pathResolver / regexResolver / Compile Cache
- httprouter / gorilla/mux / pat (Sinatra-inspired)

## Concurrency
- CSP model
- runtime.Gosched (yield)
- Loop Variable Closure Tələsi (parametr ötür!)
- sync.WaitGroup (Add/Done/Wait)
- sync.Mutex embedding / defer Unlock
- sync.RWLock (RLock/RUnlock)
- `--race` Detector / Fatal: concurrent map writes
- Kanal: send/receive/direction / time.After / select / default
- Done-channel Pattern (close yalnız sender)
- Zero Value on Closed Channel (sonsuz loop)
- Buffered Channel Lock (size 1 = semafor)

## Error/Panic
- Error as Last Return / Named Return + recover
- Error Variables (ErrTimeout, == müqayisə)
- Custom Error Type (Error() string)
- panic(errors.New(...)) idiomu
- recover / Deferred Closure Scope
- Type Assertion (`r.(error)`)
- Goroutine Function Stack (panic keçmir!)
- GoDoer / safely.Go

## Logging/Testing
- log.Logger / Ldate / Ltime / Lmicroseconds / Llongfile / Lshortfile
- Fatal* (defer-i keçir!) vs Panicln
- Back Pressure (ACK gözləmə)
- UDP vs TCP logging
- Syslog: facility / severity / priority / LOG_LOCAL3 / syslog.Dial
- Log Levels: Trace/Debug/Info/Warn/Error/Critical
- Netcat (nc -lk)
- testing.T: Error/Fatal
- Mock/Stub (interfeys + yadda saxla)
- Canary Test ("kömürdəki kanarya"; compile-time)
- Generative Testing / quick.Check / MaxCount
- testing.B / b.N / RunParallel / pb.Next / -cpu=1,2,4
- debug.PrintStack / runtime.Stack (all)
- Delve / Godebug (debugger-lər)

## Templates
- html/template (context-aware) vs text/template
- FuncMap (Funcs → Parse sırası!)
- Parse Cache (package-level; ~10x)
- Buffer Pattern (atomik cavab)
- Nested Templates / ExecuteTemplate
- define / block (Go 1.6) / Template Inheritance
- template.HTML (safe)
- SMTP / smtp.PlainAuth / SendMail

## Static/Formlar
- http.FileServer / http.Dir / ServeFile / ServeContent
- If-Modified-Since / 304 Not Modified
- StripPrefix
- RWMutex
- go-fileserver (custom error pages)
- go.rice (binary embed; rice embed-go)
- CDN / per-mühit asset
- HTTP/2 Server Push (RFC 7540)
- ParseForm / ParseMultipartForm / FormValue / PostForm
- multipart.File / FileHeader (Filename/Open)
- FormFile (tək) vs MultipartForm.File (çox)
- DetectContentType (512 bayt) / TypeByExtension / libmagic
- multipart.Reader / NextPart / chunk streaming
- CSRF Token

## REST
- http.DefaultClient / Client{Timeout} (body daxil!)
- net.Error / url.Error / OpError / Timeout()
- Range Header / Accept-Ranges / Resumable Download
- JSON Error Envelope / `json:"-"`
- Arbitrary JSON / interface{} / Type Map (float64 vs s.)
- API Versioning: URL vs vnd. Content Type / IANA
- Accept Header negotiation / fallthrough

## Cloud
- IaaS / PaaS / SaaS
- Vendor Lock-In / İnterfeys abstraksiyası / LocalFile
- Package Error Variables (divergent errors)
- runtime-detect / os.Hostname / net.LookupHost
- exec.LookPath
- GOOS / GOARCH / gox / Build Tags / foo_windows.go
- runtime.NumGoroutine / ReadMemStats / MemStats
- Remediation / Cloud-Native
- Container vs VM (host kernel)

## Kommunikasiya
- Connection Reuse / Keep-Alive
- TCP Slow-Start
- HTTP vs TCP keep-alive / DisableKeepAlives
- Body Close Tələsi (defer-in zərəri)
- Pipelining vs Multiplexing
- codecgen / JsonHandle / `//go:generate`
- Protocol Buffers / protoc / protoc-gen-go
- proto.String/Int32 (pointer) / proto.Marshal
- gRPC / proto3 / service / plugins=grpc
- RegisterHelloServer / NewHelloClient
- context.Background / WithCancel / ctx.Done() / ctx.Err()

## Reflection/Generasiya
- reflect.Value / Type / Kind (üçlük)
- Type Switch vs Kind Switch (MyInt dərsi)
- ref.Int()/Uint() (böyük formaya)
- Nil Pointer Trick + Type.Implements()
- reflect.Indirect / NumField / StructField
- Struct Tag (`NAME:"VALUE,DATA"`) / Tag.Get
- `json:"-"` (ignore)
- v.Field(i).Set (runtime yazma)
- go generate / $GOPACKAGE
- text/template ilə kod generasiyası
- go/ast / Metaprogramming
- Helm (nümunə)
