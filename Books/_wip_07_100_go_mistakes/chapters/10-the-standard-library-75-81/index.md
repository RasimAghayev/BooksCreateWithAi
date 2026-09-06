# Chapter 10 — The standard library (#75-#81)

## Bu chapter nədən bəhs edir?

Bu chapter standart kitabxananın 7 səhvini əhatə edir: `time.Duration` qarışması, `time.After` memory leak-ləri, JSON səhvləri (embedded `time.Time`, monotonic clock, `map[string]any`), SQL səhvləri (sql.Open, connection pool, prepared statements, NULL, rows.Err), `io.Closer` resurslarının bağlanması (HTTP body, sql.Rows, os.File), HTTP handler-da `return` unutması və default HTTP client/server-in production-da istifadəsi.

---

## Əsas fikirlər

### #75: Providing a wrong time duration (Yanlış müddət)

```go
ticker := time.NewTicker(1000)   // 1000 nəsaniyə = 1 MİKROSANİYƏ — 1 saniyə YOX!
```

**Səbəb:** `time.Duration` = int64 alias, NƏSANİYƏ ilə ifadə olunur. Java/JavaScript fonlu developer-lər ms gözləyir.

**Həll — həmişə time API istifadə et:**

```go
ticker = time.NewTicker(time.Microsecond)
// və ya
ticker = time.NewTicker(1000 * time.Nanosecond)
// 1 saniyə üçün:
ticker = time.NewTicker(time.Second)
```

---

### #76: time.After and memory leaks

```go
// PİS — loop-da time.After: hər iterasiya yeni timer yaratmaq deməkdir:
func consumer(ch <-chan Event) {
    for {
        select {
        case event := <-ch:
            handle(event)
        case <-time.After(time.Hour):        // hər mesajda YENİ 1-saatlıq timer!
            log.Println("warning: no messages received")
        }
        }
    }
}
```

**Problem:** time.After qaytardığı channel-i (və resurslarını) **timeout bitənə qədər** yaddaşda saxlayır. Go 1.15-də ~200 bayt/call: 5M mesaj/saat → **1 GB yaddaş!** Qaytarılan `<-chan time.Time` receive-only — proqramatik bağlana bilmir.

**Həllər:**

```go
// 1. Context (hər iterasiyada yaratma — yüngül deyil):
ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
select {
case event := <-ch:
    cancel()
    handle(event)
case <-ctx.Done():
    log.Println("warning: no messages received")
}

// 2. time.NewTimer + Reset (ƏN YAXŞI — allocation YOXDUR):
func consumer(ch <-chan Event) {
    timerDuration := 1 * time.Hour
    timer := time.NewTimer(timerDuration)
    defer timer.Stop()   // production-da exit üçün
    for {
        timer.Reset(timerDuration)    // eyni timer yenidən işə salınır
        select {
        case event := <-ch:
            handle(event)
        case <-timer.C:
            log.Println("warning: no messages received")
        }
    }
}
```

**time.After daxili:** `return NewTimer(d).C` — Reset-ə çıxış yoxdur; ona görə loop-da problem.

**Prinsip:** `time.After` yalnız tək çağırışlarda; loop/HTTP handler (təkrar çağırılan hər yer) → `time.NewTimer` + Reset.

---

### #77: Common JSON-handling mistakes

#### a) Type embedding gözlənilməz davranışı

```go
type Event struct {
    ID int
    time.Time       // EMBEDDED
}
event := Event{ID: 1234, Time: time.Now()}
b, _ := json.Marshal(event)
fmt.Println(string(b))
// Gözlənti: {"ID":1234,"Time":"2021-..."}
// Reallıq: "2021-05-18T21:15:08.381652+02:00"  ← ID İTİB!
```

**Mexanizm:** `time.Time` → `json.Marshaler` (MarshalJSON) implement edir → promote → **Event də Marshaler implement edir** → json.Marshal Event-in ÖZ marshalinqini YOX, time.Time-ınkını işlədir.

**Həllər:**

```go
// 1. Embedding-i adlandır:
type Event struct {
    ID   int
    Time time.Time      // adlı sahə
}
// → {"ID":1234,"Time":"2021-..."}

// 2. Custom MarshalJSON (embed saxlanılmalısa):
func (e Event) MarshalJSON() ([]byte, error) {
    return json.Marshal(struct {
        ID   int
        Time time.Time
    }{ID: e.ID, Time: e.Time})
}
```

#### b) Monotonic clock

**2 saat tipi:** Wall clock (gündüz vaxtı — NTP ilə sıçrayır) vs Monotonic (həmişə irəli — sıçrayışdan qorunur).

`time.Now()` hər İKİSİNİ daşıyır: `2021-01-10 17:13:08.852061 +0100 CET m=+0.000338660` (son hissə = monotonic).

**Tələ:** marshal→unmarshal monotonic-ı ATIR → `event1 == event2` → **false!**

**Həllər:**

```go
// 1. Equal metodu (monotonic-ı nəzərə almır):
fmt.Println(event1.Time.Equal(event2.Time))   // true

// 2. Truncate(0) — monotonic hissini sil:
event1 := Event{Time: t.Truncate(0)}
// indi event1 == event2 → true
```

**Location qeydi:** JSON nəticyəsi timezone-a bağlı — sabitləmək üçün `time.Now().UTC()` və ya `time.Now().In(location)`.

#### c) Map of any

```go
var m map[string]any
json.Unmarshal(b, &m)   // {"id": 32, "name": "foo"}
fmt.Printf("%T\n", m["id"])
// float64 — int YOX!
```

**Qayda:** `map[string]any` ilə unmarshal-da BÜTÜN numeric dəyərlər (kəsr olmasa belə) `float64` olur — yanlış pressumpla goroutine panic yarada bilər.

---

### #78: Common SQL mistakes

#### a) sql.Open bağlantı AÇMIR (zəruri deyil)

```go
db, err := sql.Open("mysql", dsn)   // yalnız arqumentləri validasiya edə bilər!
```

Bağlantı driver-dən asılı olaraq **lazily** açılır. Servisin dependency-lər hazır olanda başlaması üçün:

```go
db, err := sql.Open("mysql", dsn)
if err != nil { return err }
if err := db.Ping(); err != nil { return err }   // (və ya PingContext)
```

#### b) Connection pooling unutması

`*sql.DB` = tək bağlantı YOX, **bağlantı HOVUZU** (istifadədə/idle). 4 konfiqurasiya parametri:

| Parametr | Default | Niyi dəyişməli |
|----------|---------|----------------|
| `SetMaxOpenConns` | limitsiz | Production MÜTLƏQ — DB-nin dözdüyü həddə |
| `SetMaxIdleConns` | 2 | Yüksək konkurrentlikdə artır (əks halda tez-tez reconnect) |
| `SetConnMaxIdleTime` | limitsiz | Burst-dən sonra bağlantılar azad olunsun |
| `SetConnMaxLifetime` | limitsiz | Load-balanced DB-lərdə çox uzun istifadəyin qarşısı |

#### c) Prepared statements

```go
stmt, err := db.Prepare("SELECT * FROM ORDER WHERE ID = ?")
if err != nil { return err }
defer stmt.Close()
rows, err := stmt.Query(id)
```

**2 fayda:** Effektivlik (təkrar kompilyasiya YOX: parse+optimize+translate bir dəfə) + Təhlükəsizlik (SQL injection azalır). Təkrarlanan/etibarsız kontekstdəki sorğular üçün şərt. Context variantları: `PrepareContext`, `QueryContext`.

#### d) NULL dəyərləri

```go
// PİS — "converting NULL to string is unsupported":
var department string
rows.Scan(&department, &age)

// Həll 1 — pointer:
var department *string   // NULL → nil

// Həll 2 — sql.NullXXX (niyyət daha aşkar — Russ Cox):
var department sql.NullString   // .String + .Valid sahələri
```

Wrapper-lər: `sql.NullString/NullBool/NullInt32/NullInt64/NullFloat64/NullTime`.

#### e) rows.Err unutması

```go
for rows.Next() {   // loop 2 səbəbdən qırılır: ya sətirlər bitdi, YAXUD xəta!
    rows.Scan(&department, &age)
}
if err := rows.Err(); err != nil {   // MUTLƏQ — fərqi ayırır
    return "", 0, err
}
```

---

### #79: Not closing transient resources (Keçici resursların bağlanması)

**Qayda:** `io.Closer` (`Close() error`) implement edən hər struktur nə vaxtsa bağlanmalıdır.

#### HTTP body

```go
// PİS — resp.Body bağlanmır → yaddaş + TCP bağlantı reuse qadağası:
resp, err := h.client.Get(h.url)
body, err := io.ReadAll(resp.Body)
return string(body), nil

// DÜZGÜN:
resp, err := h.client.Get(h.url)
if err != nil { return "", err }
defer func() {
    err := resp.Body.Close()
    if err != nil {
        log.Printf("failed to close response: %v\n", err)
    }
}()
body, err := io.ReadAll(resp.Body)
```

**Nüanslar:**
- Server tərəfdə request body-ni bağlamaq LAZIM DEYİL (server avtomatik edir).
- Body OXUNMASA BELƏ bağlanmalıdır (yalnız status kod lazımdırsa).
- **Oxumadan bağlansa** — default transport bağlantını BAĞLAYA bilər; **oxuyub bağlansa** — keep-alive reuse mümkün. Status kod funksiyası tez-tez çağırılırsa:

```go
_, _ = io.Copy(io.Discard, resp.Body)   // oxu + at — ReadAll-dan effektiv
```

- `if resp != nil { defer resp.Body.Close() }` — LAZIMSIZ yoxlama: *"On error, any Response can be ignored. A non-nil Response with a non-nil error only occurs when CheckRedirect fails, and even then, the returned Response.Body is already closed."*

#### sql.Rows

```go
rows, err := db.Query("SELECT * FROM CUSTOMERS")
if err != nil { return err }
defer func() {
    if err := rows.Close(); err != nil {
        log.Printf("failed to close rows: %v\n", err)
    }
}()
```

Bağlamamaq = connection leak → bağlantı hovuza QAYITMIR. (`*sql.DB` də Closer-dır amma adətən uzunömürlüdür — bağlamaq nadir.)

#### os.File

```go
f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
if err != nil { return err }
defer func() {
    if err := f.Close(); err != nil {
        log.Printf("failed to close file: %v\n", err)
    }
}()
```

GC avtomatik bağlaya bilər, amma nə vaxt — bəlli deyil → explicit Close (+ error nəzarəti).

**Yazıla bilən fayllarda Close error-ü propagasiya et** (buffer-lənmiş yazma commit olmamış ola bilər — close(2) BSD manual):

```go
func writeToFile(filename string, content []byte) (err error) {
    // Open file
    defer func() {
        closeErr := f.Close()
        if err == nil {
            err = closeErr      // write uğurlu olsa close error-u qaytar
        }
    }()
    _, err = f.Write(content)
    return
}
```

**Dayanıqlıq (durability) kritikdirsə:** `f.Sync()` — diskə sync commit; onda Close error-u iqnor edilə bilər:

```go
defer func() { _ = f.Close() }()
_, err = f.Write(content)
if err != nil { return err }
return f.Sync()   // məzmun diskdir — performance təsiri VAR
```

---

### #80: Forgetting the return statement after replying to an HTTP request

```go
// PİS — http.Error handler-i DAYANDIRMIR:
func handler(w http.ResponseWriter, req *http.Request) {
    err := foo(req)
    if err != nil {
        http.Error(w, "foo", http.StatusInternalServerError)
        // return YOXDUR!
    }
    _, _ = w.Write([]byte("all good"))
    w.WriteHeader(http.StatusCreated)
}
```

**Nəticə:**
- Response body hər ikisini daşıyır: `foo\nall good`
- Yalnız İLK status kodu (500) göndərilir; sonra `superfluous response.WriteHeader call` warning.
- Davam edən icra → foo pointer qaytarırsa **nil pointer dereference → goroutine panic!**

```go
// DÜZGÜN:
if err != nil {
    http.Error(w, "foo", http.StatusInternalServerError)
    return     // MÜTLƏQ
}
```

Test coverage bu səhvi tutmalıdır.

---

### #81: Using the default HTTP client and server

#### Client

`&http.Client{}` == `http.Get` == `http.DefaultClient` — **timeout YOXDUR** → sonsuz request-lər → resurs tükənməsi.

**HTTP request-in 5 addımı:** Dial → TLS handshake → Request göndər → Response header oxu → Response body oxu.

**4 əsas timeout:**

| Timeout | Nəyi əhatə edir |
|---------|-----------------|
| `net.Dialer.Timeout` | TCP dial (addım 1) |
| `http.Transport.TLSHandshakeTimeout` | TLS handshake (addım 2) |
| `http.Transport.ResponseHeaderTimeout` | Response header (addım 4) |
| `http.Client.Timeout` | QLOBAL — addım 1-dən 5-ə qədər hamısı |

**Production client:**

```go
client := &http.Client{
    Timeout: 5 * time.Second,
    Transport: &http.Transport{
        DialContext: (&net.Dialer{
            Timeout: time.Second,
        }).DialContext,
        TLSHandshakeTimeout:    time.Second,
        ResponseHeaderTimeout: time.Second,
    },
}
```

**Connection pooling parametrləri:**
- `http.Transport.IdleConnTimeout` — idle bağlantının hovuzda qalma müddəti (default 90 s)
- `http.Transport.MaxIdleConns` — ümumi limit (default 100)
- `http.Transport.MaxIdleConnsPerHost` — **HOST BAŞINA default 2!** — 100 request→1 host → yalnız 2 bağlantı hovuzda qalır → 98 reconnect. Paralel eyni-host trafikində latency təsiri.

**Timeout error-i:** `net/http: request canceled (Client.Timeout exceeded while awaiting headers)` — header ilk gözlənilən addım olduğundan.

#### Server

`&http.Server{}` / `http.ListenAndServe` — default, timeout-suz.

**Server response 5 addımı:** Client-i gözlə → TLS handshake → Request header oxu → Request body oxu → Response yaz.

**3 əsas timeout:**

| Timeout | Nəyi əhatə edir |
|---------|-----------------|
| `http.Server.ReadHeaderTimeout` | Request header oxunuşu |
| `http.Server.ReadTimeout` | Tam request oxunuşu (header+body) |
| `http.TimeoutHandler(handler, t, msg)` | Handler icra müddəti — limit aşılarsa **503** + ctx ləğv |

**`http.Server.WriteTimeout` haqqında:** TimeoutHandler (Go 1.8+) buraxıldıqdan sonra lazımsız — TLS-dən asılı davranış + TCP-ni düzgün kod olmadan qirir + ctx ləğvini handler-a çatdırmır.

**Production server:**

```go
s := &http.Server{
    Addr:              ":8080",
    ReadHeaderTimeout: 500 * time.Millisecond,
    ReadTimeout:       500 * time.Millisecond,
    Handler:           http.TimeoutHandler(handler, time.Second, "foo"),
    IdleTimeout:       time.Second,
}
```

IdleTimeout qoyulmazsa ReadTimeout istifadə olunur; heç biri yoxdursa — bağlantılar klient qədər AÇIQ (etibarsız klientlər resurs tükədir).

---

## Əsas terminlər

| Termin | İzah |
|--------|------|
| `time.Duration` | int64 alias — NƏSANİYƏ vahidli |
| `time.NewTimer` + `Reset` | Loop-larda time.After əvəzi — allocationsız yenidən istifadə |
| json.Marshaler promote | Embedded tip interfeysi bütün struct-a keçir — marshalinq əzilir |
| Wall vs monotonic clock | NTP-sıçrayanlı / həmişə-irəli saat; `m=+` hissəsi monotonic |
| `time.Time.Equal` / `Truncate(0)` | Monotonic-fiksiz müqayisə yolları |
| `map[string]any` numeric → float64 | Bütün ədədlər float64 — int gözləmək panic |
| `sql.Open` lazy | Bağlantı zəruri deyil — `Ping`/`PingContext` ilə doğrula |
| Connection pool (`*sql.DB`) | 4 parametr: MaxOpen/MaxIdle/MaxIdleTime/MaxLifetime |
| Prepared statement | Precompiled sorğu — effektivlik + SQL injection qoruması |
| `*string` / `sql.NullXXX` | NULL kolonlarının 2 həlli |
| `rows.Err()` | Loop-un səbəbini ayırır: bitdi YOXSA xəta |
| `io.Closer` / `io.ReadCloser` | Bağlanmalı keçici resurs interfeysləri |
| `io.Discard` | Oxu + atma writer-ı — keep-alive üçün body tüketmə |
| `superfluous WriteHeader` | İkiqat status yazma warning-i — return unutmuşsan |
| `http.DefaultClient` | Timeout-suz default — production üçün təhlükə |
| `MaxIdleConnsPerHost` | Default 2 — eyni-host keep-alive bottleneck |
| `http.TimeoutHandler` | Handler müddət limiti — 503 + ctx ləğv |
| `f.Sync()` | Disk commit — dayanıqlı yazma (Close error-u sonra iqnor) |

---

## Praktik nəticə

1. **Duration heç vaxt "çılpaq" rəqəm:** `time.Second`, `100 * time.Millisecond` — vahidli API.
2. **Loop/handler-da `time.After` YOX:** `time.NewTimer` + `Reset` (200 bayt/call yığma; 1GB-lər qarşısı). `defer timer.Stop()`.
3. **Embedded struct + JSON:** Embedded tip `Marshaler` implement edirsə bütün struct-ın davranışı DƏYİŞİR — sahəni adlandır.
4. **`time.Now()` müqayisəsi:** == monotonic-ı da yoxlayır — `Equal` və ya `Truncate(0)`; JSON round-trip == pozur.
5. **`any` map-i numeric = float64:** int gözləyən type assertion panic edir.
6. **SQL açılışı doğrula:** `sql.Open` + `Ping`; pool parametrlərini production üçün konfiqurasiya et.
7. **Təkrar/etibarsız sorğu = prepared statement:** `Prepare` + `stmt.Query` + `defer stmt.Close()`.
8. **NULL üçün `*string` və ya `sql.NullString`;** iterasiyadan sonra MÜTLƏQ `rows.Err()`.
9. **`io.Closer` olan hər şey bağlanır:** HTTP body (oxunmasa belə!), sql.Rows, os.File — defer + error idarəsi; yazılacaq faylda close error-u propagasiya et, durability üçün `Sync`.
10. **`http.Error` saxlamır — `return` yaz:** əks halda qarışıq body + superfluous WriteHeader + nil dereference riski.
11. **Default client/server QADAĞANDIR:** client üçün 4 timeout (dial/TLS/header/global) + pooling (MaxIdleConnsPerHost=2 diqqət!); server üçün ReadHeaderTimeout + ReadTimeout + TimeoutHandler + IdleTimeout.

---

## Mənbə

- Kitab: *100 Go Mistakes and How to Avoid Them* — Teiva Harsanyi, Manning, 2022 (ISBN 9781617299599)
- Chapter 10: "The standard library", book səh. 234–261
- PDF səhifələri: 254–281
- İstinadlar: database/sql docs; net/http docs; Russ Cox NullString sitatı (mng.bz/rJNX); BSD close(2) manual
