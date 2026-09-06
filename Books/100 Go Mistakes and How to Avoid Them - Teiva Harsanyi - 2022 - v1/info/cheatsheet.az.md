# 100 Go Mistakes and How to Avoid Them — Cheatsheet (Azərbaycanca)

> **Kitab:** 100 Go Mistakes and How to Avoid Them — Teiva Harsanyi, Manning 2022 · 100 səhv, 12 chapter · Sürətli istinad.

---

## Code & Project Organization (#1-#16)

```go
// #1 Shadowing — daxili blokda := YENİ dəyişən yaradır:
var client *http.Client; var err error
client, err = createClient()      // = istifadə et (:= YOX)

// #3 init yalnız uğura bilən statik konfiq üçün; DB üçün adi funksiya:
func createClient(dsn string) (*sql.DB, error)

// #5/#6/#7 Interface-lər: abstractions should be DISCOVERED, not created.
// Interface consumer tərəfdə (minimal, unexported); qaytar STRUCT, qəbul et INTERFACE (Postel's law)

// #8 any ifadəsizdir — tipə xas metodlar yaz (json.Marshal istisna)

// #9 Generics: data strukturlar + slice/map/chan funksiyaları + davranış factor-out.
// Metodlarda type parametr YOX. Tipləri çağıran metod üçün generics YOX.

// #10 Embed: yalnız davranış promote üçün; sync.Mutex EMBED ETMƏ (mu sync.Mutex)

// #11 Functional options:
type Option func(*options) error
func WithPort(port int) Option { return func(o *options) error { /* validate+set */ } }
func NewServer(addr string, opts ...Option) (*Server, error)

// #13 util/common/shared paketlər YOX — stringset kimi ifadəli ad

// #15 Hər exported element dokumentləşdir: // Name does X; paket: // Package name ...
// #16 go vet + errcheck + gocyclo + gofmt + golangci-lint (CI/precommit)
```

---

## Data Types (#17-#29)

```go
// #17 010 = 8 (oktal!) → 0o644 yaz
// #18 Overflow silent-dir:  if counter == math.MaxInt { panic("overflow") }
//      Add: a > math.MaxInt-b;  Mul: result/b != a

// #19 Float == YOX — delta müqayisə; vurma/bölmə əvvəl; böyüklük-qruplu toplama

// #20 Slice = ptr+len+cap; append len<cap olanda ORTAQ array-ə yazır!
// #21 make([]T, 0, n) və ya make([]T, n) — 400% fərq
// #22 nil slice allocationsız: var s []string; JSON-da: nil→null, empty→[]
// #23 Boşluq yoxlaması: len(s) != 0 (nil/empty hər ikisi)
// #24 copy min(len(dst), len(src)) — dst əvvəlcədən make(len(src))
// #25 Ortak array: s[:2:2] full slice expression və ya kopya
// #26 Böyük slice-dan kiçik saxlayırsansa KOPYALA (cap leak); pointerli elementləri nil-lə

// #27 make(map[K]V, n) — 60% sürətli
// #28 Map yalnız BÖYÜYÜR — pikdən sonra: re-create və ya map[int]*[128]byte (38MB vs 293MB)

// #29 == comparable üçün; slice/map YOX → errors.As/Is + reflect.DeepEqual (test, ~100x yavaş)
//      və ya custom equal() (~96x sürətli)
```

---

## Control / Strings / Functions (#30-#47)

```go
// #30 Range value = KOPYA → accounts[i].balance (indekslə)
// #31 Range ifadəsi 1 dəfə qiymətlənir; klassik for-da len(s) canlıdır
// #32 Loop dəyişəninə pointer saxlaMA: current := customer / &customers[i]
// #33 Map sırası UNPSECIFIED; insert zamanı iteration → kopya üzərində yaz
// #34 break = innermost! Loop üçün: break loop (label — idiomatikdir, goto deyil)
// #35 Loop-da defer YOX → hər iterasiya üçün ayrı funksiya/closure

// #36 rune = int32 code point; len(s) = BAYT; utf8.RuneCountInString(s) = rune
// #37 for i, r := range s → r rune-un ÖZÜ; i = bayt başlanğıcı
// #38 TrimRight/Left = SET (təkrar); TrimSuffix/Prefix = tək dəfə
// #39 5+ string birləşdirəndə strings.Builder{ Grow(total); WriteString }
// #40 I/O-da bytes istifadə et — konversiya = kopya + alloc
// #41 log[:36] substring BÜTÜN mesajı yaşadır → string([]byte(...)) / strings.Clone

// #42 Receiver: mutasiya/kopyalanmayan sahə → pointer; map/func/chan → value; şübhə → pointer
// #43 Named result: eyni tipli nəticələr üçün (lat, lng float32); tək err üçün YOX
// #44 Named + return 0,0,err → err zero-value NIL-dir! (if err := ctx.Err(); err != nil)
// #45 Nil pointer → error interface = NON-NIL! Sonda: if m != nil { return m }; return nil
// #46 Filename YOX — io.Reader qəbul et (strings.NewReader ilə test)
// #47 Defer arqumentləri DƏRHAL qiymətlənir → pointer ötür / closure istifadə
```

---

## Errors (#48-#54)

```go
// #48 Panic yalnız: programmer error + məcburi asılılıq (MustCompile)
// #49 %w = wrap (source açıq — coupling); %v = transform (source bağlı)
// #50 Tip yoxlaması: errors.As(err, &transientError{})  // pointer hədəf ŞƏRT
// #51 Sentinel dəyər: errors.Is(err, sql.ErrNoRows)    // == YOX
//      Expected → sentinel value; unexpected → custom type
// #52 Log YAXUD return — heç vaxt hər ikisi (context-i %w wrap ilə ver)
// #53 İqnor: _ = notify()  + səbəb kommenti
// #54 Defer Close: minimum _ =; daha yaxşısı log; propagasiya = named err + closeErr ayrılığı
```

---

## Concurrency (#55-#74)

```go
// #55 Concurrency = STRUKTUR; Parallelism = İCRA
// #56 Konkurrentlik həmişə sürətli DEYİL — threshold: if len(s) <= 2048 { sequentialMergesort(s) }
//      G/M/P; GOMAXPROCS = M limiti; work stealing; Go 1.14+ preemptive
// #57 Paralel goroutine → MUTEX; konkurrent → CHANNEL
// #58 Data race = eyni yaddaş + yazma; race condition = sıra asılılığı.
//      Memory model: go create < icra; unbuffered: receive < send-tamamlanma

// #61 HTTP context response yazılanca ÖLÜR → detach / context.Background
// #62 Hər goroutine-in EXIT planı olsun; resurs üçün blocking close + defer
// #63 Closure + loop dəyişəni: val := i  /  go func(val int){...}(i)
// #64 Select RANDOM seçir! Prioritet: inner select + default
// #65 Siqnal kanalı: chan struct{} (0 bayt)
// #66 Qapalı kanalı nil-lə → select case-i SÖNDÜR (merge pattern)
// #67 Sinxronizasiya = unbuffered; buffered default ölçü = 1 (magic 40 YOX)
// #68 %v gizli Stringer çağırır → deadlock/race; lock-dan ƏVVƏL validasiya
// #69 Ortak slice-da append = race (dolmamışsa) — kopya istifadə et
// #70 Map/slice assign = SHALLOW — ya tam funksiya lock, ya dərin kopya
// #71 wg.Add valideyndə, spin-dən ƏVVƏL; Done uşaqda
// #72 Təkrar broadcast çoxlu goroutine-yə → sync.Cond (Wait: unlock-suspend-relock)
// #73 errgroup.WithContext + g.Go(func() error) + g.Wait() — error + ctx ləğvi birlikdə
// #74 Sync tiplər KOPYALANMAZ — pointer receiver / *sync.Mutex sahə
```

---

## Stdlib / Testing (#75-#90)

```go
// #75 time.NewTicker(time.Second) — çılpaq 1000 YOX (nanosaniyə!)
// #76 Loop-da time.After YOX → time.NewTimer + Reset (+ defer Stop)
// #77 Embedded time.Time Event-i MarshalJSON-la ƏZİR — sahəni adlandır
//      time.Now() == müqayisə monotonic poza bilər → .Equal / .Truncate(0)
//      map[string]any numeric → HAMISI float64
// #78 sql.Open bağlantı AÇMIR — Ping ilə doğrula; pool: SetMaxOpen/Idle/IdleTime/Lifetime
//      Prepare (injection + effektivlik); NULL → *string / sql.NullString; rows.Err() ŞƏRT
// #79 io.Closer hamısı: HTTP body (oxunmasa belə!), sql.Rows, os.File
//      Keep-alive: io.Copy(io.Discard, resp.Body); yazıla bilən faylda close-err propagate / Sync
// #80 http.Error saxlamır — return ƏLAVƏ ET
// #81 Default client/server QADAĞAN: 4 client timeout + MaxIdleConnsPerHost(2!);
//      server: ReadHeaderTimeout + ReadTimeout + TimeoutHandler + IdleTimeout

// #82 Kateqoriya: //go:build integration / env+t.Skip / -short
// #83 go test -race — vector-clock, false positive YOXDUR, loop-da təkrarla
// #84 t.Parallel + -parallel N; -shuffle=on (seed reproduce)
// #85 Table-driven: map[string]struct{in, exp} + t.Run; paraleldə tt := tt!
// #86 Sleep = flaky → channel sinxronizasiyası > retry assert > sleep
// #87 time.Now testi: type now func() time.Time sahəsi / client ötürməsi
// #88 httptest.NewRecorder/NewRequest (handler); httptest.NewServer (client)
//      iotest.TestReader / TimeoutReader / HalfReader
// #89 Benchmark: ResetTimer/Stop-Start; -count=10 + benchstat; local→global pattern
//      (inline qarşı); hər iterasiyada YENİ data (observer effect)
// #90 -coverpkg; package x_test; utility(t *testing.T); t.Cleanup; TestMain(m)
```

---

## Optimizations (#91-#100)

```go
// #91 Cache line 64B; RAM 50-100x L1-dən yavaş; struct-of-slices > slice-of-structs
//      Linked list = non-unit stride (~70% yavaş); critical stride = 512 int64 (4KB)
// #92 False sharing: pad: _ [56]byte ayrıcı / channel kommunikasiyası (~40%)
// #93 ILP: v := s[0]; s[0] = v+1; if v%2 != 0 {...}  (~20% — hazard azaltma)
// #94 Sahələr böyükdən kiçiyə: i int64; b1, b2 byte → 24→16 bayt
// #95 Sharing up = heap (pointer return, 12x yavaş); sharing down = stack
//      go build -gcflags "-m=2"
// #96 m[string(bytes)] — konversiya-atlama; sync.Pool{New} + Get + [:0] + defer Put
// #97 Fast-path inlining: slow path-i ayrı funksiyaya çıxar (sync.Mutex.Lock pattern)
// #98 pprof: CPU (SIGPROF 10ms), heap diff (leak), goroutine dump, block, mutex
//      tracer: go tool trace — paralellik/GC vizual; runtime/trace = task bölgüsü
// #99 GOGC=100: heap 2x olanda GC; pik üçün artır; virtual heap hilesi: make(1GB) mmap
// #100 K8s: GOMAXPROCS host core-larına görədir — CFS throttling!
//      import _ "go.uber.org/automaxprocs"  → kvotaya uyğunlaşdır
```

---

## Sürətli yaddaş cədvəli

| Ehtiyac | Həll |
|---------|------|
| Assignment daxili blokda | `=` + əvvəlcədən `var err error` |
| Interface qaytarma | Struct qaytar, interface qəbul et |
| Optional konfiq | Functional options (`WithPort(8080)`) |
| Overflow yoxlama | `math.MaxInt` qarşılaşdırmaları |
| Boşluq yoxlaması | `len(s) != 0` |
| Slice kopyası | `make(len(src))` + `copy` / `append([]T(nil), s...)` |
| Append qoruması | `s[:2:2]` full slice expression |
| String birləşdirmə | `strings.Builder` + `Grow` |
| Substring saxlama | `strings.Clone` |
| Nil receiver → error | `if m != nil { return m }; return nil` |
| Error tipi/dəyəri | `errors.As` / `errors.Is` |
| Select prioritet | Inner select + `default` |
| Qapalı kanal idarəsi | `v, open := <-ch` → kanalı `nil`-lə |
| Paralel + error | `errgroup.WithContext` |
| Loop-da vaxt aşımı | `time.NewTimer` + `Reset` |
| HTTP production client | 4 timeout + custom Transport |
| Handler xətası | `http.Error` + `return` |
| Flaky test | Channel mock / retry assert |
| Vaxt asılı test | `type now func() time.Time` inyeksiya |
| Benchmark dəqiqliyi | Reset timer, local→global, yeni data, benchstat |
| Heap azalt | Sharing-down API, `m[string(b)]`, `sync.Pool` |
| Leak tapma | pprof heap diff (`-diff_base`) |
| K8s deploy | `_ "go.uber.org/automaxprocs"` |
