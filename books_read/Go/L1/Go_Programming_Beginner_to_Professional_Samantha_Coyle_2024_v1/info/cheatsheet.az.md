# Go Programming B2P — Cheat Sheet (AZ)

## Dil əsasları (ch1-3)

### `var x T` / `x := v` / `:=` yalnız funksiya daxilində
Zero value-lar: 0 / false / "" / nil.

### Shadowing: `:=` uşaq scope-də YENİ dəyişən; köhnəni dəyişmək üçün `=`

### `&x` / `*p` / `new(T)` — pointer trio
Nil dereference → panic; yoxlamadan keçirməyin.

### wraparound: int8 127+1 → -128; `math/big` həddən artıq üçün

### `[]rune(s)` — simvol əməliyyatları üçün; `len(s)` bayt sayır!

## Kollesiya (ch4)

### `make([]T, len, cap)` + `append` — cap dolanda yeni massiv
```go
s2 := append([]int{}, s1...)    // MÜSTƏQİL kopya
s2 := append(s1[:0:0], s1...)   // ən effektiv kopya
// s2 := s1 → LINK! gizli massiv paylaşılır
```

### `v, ok := m[k]` / `delete(m, k)`
nil map-a yazma → panic; make/literal şərt.

### `any` + `switch v := i.(type)`
Çoxlu-tip case-də assertion lazım.

## Funksiya/xəta (ch5-6)

### `func f(a, b int) (int, error)` — error-sonuncu

### `nums ...int` qəbul; `slice...` ötür

### `defer` — FILO; arqumentlər İNDİ donur

### `var ErrX = errors.New("...")` + `if err != nil { return err }`

### panic yalnız bərpaolunmaz; sargı:
```go
defer func() {
    if e := recover(); e != nil { err = fmt.Errorf("%v", e) }
}()
```

### `fmt.Errorf("context: %w", err)` — wrap zənciri

## İnterfeys/generics (ch7-8)

### "Accept interfaces, return structs"
`func load(r io.Reader)` — string/fayl/HTTP hamısı uyğun.

### `type Speaker interface { Speak() string }` — implisit, -er adı

### `func Max[T Number](s []T) (T, error)` + `var zero T` + `~int` tilde
Map açarı: `[K comparable, V int|float64]`.

## Modullar/paketlər (ch9-10)

### `go mod init` → `go get pkg` → `go mod tidy`
go.sum = SHA-256 qoruması.

### `go work init` + `go work use ./mod` — lokal multi-modul (replace əvəzi)

### Exported = BÖYÜK hərf; init sırası: import → var → init → main

### cmd/ (main) + pkg/ (kitabxana)

## Vaxt (ch12)

### `start := time.Now()` → `end.Sub(start).Seconds()`

### `time.Parse(time.RFC3339, s)` ↔ `t.Format(time.ANSIC)`

### `time.LoadLocation("America/Los_Angeles")` + `t.In(loc)`

### `t.Add(-10 * time.Minute)` — geri; `AddDate(1, 2, 3)` — irəli

## CLI/fayl (ch13-14)

### `flag.String("name", "Sam", "help")` + `flag.Parse()` + `*nameFlag`

### `stat, _ := os.Stdin.Stat(); (stat.Mode() & os.ModeCharDevice) == 0` → pip edilir

### `os.Create` (truncate!) / `os.OpenFile(f, O_APPEND|O_CREATE|O_WRONLY, 0644)` — 0644 OKTAL!

### `os.Stat` + `os.IsNotExist(err)` — mövcudluq

### `os.ReadFile` (kiçik) vs `f.Read(buf)` loop (böyük)

### `signal.Notify(sigs, syscall.SIGINT)` + goroutine dinləyici + cleanup
defer Ctrl+C-də çağırılmır!

### `//go:embed templates` + `var f embed.FS`

## SQL (ch15)

### `db, _ := sql.Open("postgres", dsn)` + `db.Ping()` + `defer db.Close()`

### `stmt, _ := db.Prepare("... VALUES($1, $2)")` + `stmt.Exec(a, b)` — injection qoruması

### `rows.Next() + rows.Scan(&a, &b)` — çoxsətir; `QueryRow().Scan` — tək

### `res.RowsAffected()` — UPDATE/DELETE uğur sayı

### GORM: `gorm.Open(...)` → `AutoMigrate(&T{})` → `Create/First/Find` + `tx.Error`

## Web (ch16-17)

### `http.HandleFunc("/x", fn)` / `http.Handle("/", handler{})`
Handler struct = vəziyyət saxlayır (pointer receiver!).

### Middleware: `func(HandlerFunc) HandlerFunc` → `Hello(Function1)`

### `r.URL.Query()` → `map[string][]string` + ok yoxlaması

### `template.New("t").Parse(s)` main-də; `tmpl.Execute(w, data)`

### `http.FileServer(http.Dir("./public"))` + `http.StripPrefix("/statics/", ...)`

### Klient: `http.Client{Timeout: 11s}` + `http.NewRequest` + `req.Header.Set` + `client.Do(req)`
DefaultClient production-da YOX; fayl upload: multipart.NewWriter +
FormDataContentType.

## Konkurensiya (ch18)

### `wg.Add(n)` + `go f(wg, &res)` (daxilində `wg.Done()`) + `wg.Wait()`

### `atomic.AddInt32(&s, int32(i))` — rəqəm yarışlarına; `go test -race` — yoxla!

### `mtx.Lock()` ... `mtx.Unlock()` — qısa kritik sahə

### `ch := make(chan T)` / `make(chan T, n)`; `ch <- v` / `<-ch` / `close(ch)`
Unbuffered tək-goroutine-də deadlock!

### `for i := range ch` — close-a qədər; `out <- "done"` — bitmə bildirişi

### `cl, stop := context.WithCancel(ctx)` + `select { case <-cl.Done(): return }` + `stop()`

### `sync.NewCond(&sync.Mutex{})` + `Wait()/Signal()` — bloklı növbə

### `sync.Map` — `LoadOrStore(k, 0)` + `Store(k, v.(int)+1)` + `Range`

## Test (ch19)

### Table-driven + subtest:
```go
for _, test := range tests {
    test := test               // closure kopyası!
    t.Run(test.name, func(t *testing.T) {
        got := f(test.inputs...)
        assert.Equal(t, test.want, got)
    })
}
```

### sqlmock: `mock.ExpectExec("INSERT...").WithArgs(...).WillReturnResult(...)` + `mock.ExpectationsWereMet()`

### httptest.NewServer(handler) + `authService.URL`

### `go test -fuzz .` / `-bench .` / `-cover` (~80%) / `-v -json > report.json`

### TestMain: `setup(); defer teardown(); m.Run()`

## Alətlər (ch20)

### `gofmt -w` / `goimports -w` / `go vet` / `go run --race` / `go doc -all`

## Cloud (ch21)

### Prometheus:
```go
healthzCounter := prometheus.NewCounter(prometheus.CounterOpts{
    Name: "healthz_calls_total", Help: "...",
})
prometheus.MustRegister(healthzCounter)
http.Handle("/metrics", promhttp.Handler())
// handler-da: healthzCounter.Inc()
```

### OTel: `otlptracegrpc.New(ctx, WithEndpoint(...))` + `sdktrace.NewTracerProvider(WithBatcher(...))` + `otel.SetTracerProvider(tp)` + `otelhttp.NewHandler(handler, "HTTPServer")` + `logger.Infow("msg", "key", val, ...)`

### Docker multi-stage:
```dockerfile
FROM golang:latest AS builder
ENV CGO_ENABLED=0
COPY go.mod go.sum ./
RUN go mod download
COPY main.go ./
RUN go build -o app .
FROM scratch
COPY --from=builder /app/app /
EXPOSE 8080
CMD ["./app"]
```

### K8s: /healthz + readiness probe + graceful shutdown + Deployment/Service manifestləri
