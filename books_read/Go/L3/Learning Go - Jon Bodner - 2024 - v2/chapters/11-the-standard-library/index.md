# Chapter 11 — The Standard Library (Standart Kitabxana)

## Bu chapter nədən bəhs edir?

io paketi (Reader/Writer fəlsəfəsi), time (Duration/Time, monotonic saat, formatlaşdırma),
encoding/json (struct tag-lər, Marshal/Unmarshal, Encoder/Decoder, custom parsing) və
net/http (client, server, ServeMux, middleware).

## Əsas fikirlər

### 1. io — Reader və Writer Fəlsəfəsi
**Nədir:** Go I/O-nun ürəyi — 2 ən çox istifadə olunan interface (`error`-dan sonra):

**Kitabdan kod nümunəsi:**
```go
type Reader interface {
    Read(p []byte) (n int, err error)
}
type Writer interface {
    Write(p []byte) (n int, err error)
}
```

**Niyə `Read(p []byte)` (qaytarma ilə yox)?** 3 səbəb:
1. **Bufferi çağıran idarə edir** — bir dəfə yarat, dəfələrlə reuse (hər çağırışda
   heap allocation YOX); daha da azaltmaq üçün buffer pool.
2. **n ilə dolu hissə** — `buf[:n]` subslice emal olunur.
3. **io.EOF** — "bitiş" sentinel error-u; error-u adətən dəyərlərdən ƏVVƏL yoxlasaq da,
   Read-də ƏKSİNƏ — son chunk error-dən əvvəl gələ bilər. Gözlənilməz kəsilmə:
   `io.ErrUnexpectedEOF`.

**İdiomatik oxu pattern:**
```go
func countLetters(r io.Reader) (map[string]int, error) {
    buf := make([]byte, 2048)      // bir dəfə ayrılır
    out := map[string]int{}
    for {
        n, err := r.Read(buf)
        for _, b := range buf[:n] { /* emal */ }
        if err == io.EOF { return out, nil }
        if err != nil { return nil, err }
    }
}
```

**Decorator pattern:** `strings.NewReader`, `gzip.NewReader(os.File)` — hamısı eyni
interface; `buildGZipReader` fayl adı → `*gzip.Reader` + **closer closure** qaytarır
(iki resursun da təmizliyi). Standard helper-lər: `io.Copy`, `io.MultiReader`,
`io.LimitReader`, `io.MultiWriter`.

**Digər io interface-lər:** `Closer` (defer ilə), `Seeker` (random access; whence:
`io.SeekStart/Current/End`); kombinasiyalar: `ReadCloser`, `ReadSeeker`,
`ReadWriteCloser` və s. Funksiyalarınızda `*os.File` yox, konkret interface qəbul edin.

**NopCloser hiyəsi** (interfeysə metod əlavə etmək — embedding pattern):
```go
type nopCloser struct { io.Reader }
func (nopCloser) Close() error { return nil }
func NopCloser(r io.Reader) io.ReadCloser { return nopCloser{r} }
```
ioutil.ReadAll/ReadFile/WriteFile — kiçik data üçün; böyüklər üçün `bufio`.

**Loop şəxsiyyəti:** loop daxilində fayl açırsansa defer YOX — hər iterasiyada aç,
hər hərəkətdə Close çağır (defer funksiya sonuna qədər gecikir).

### 2. time — Vaxt Paketi
**time.Duration:** int64 əsaslı; nanosaniyədən saata qədər tipli konstantlar:
```go
d := 2 * time.Hour + 30 * time.Minute
time.ParseDuration("2h45m")   // "300ms", "-1.5h" formatları
d.Truncate(time.Minute); d.Round(time.Hour)
```

**time.Time:** vaxt zonası daxil; `time.Now()` — lokal vaxt. **`==` İŞLƏTMƏ** —
zamana baxmayaraq müqayisə üçün `.Equal()`. Müqayisə: `After/Before/Equal`;
həcmetmə: `Sub` (Duration qaytarır) / `Add` / `AddDate`; çıxarış: `Year/Month/Day/
Hour/Minute/Second/Weekday/Date/Clock`. Hamısı value receiver — immutable.

**Formatlaşdırma — reference time:** Go strftime YOX; "Jan 2, 2006 3:04:05PM MST"
(1-2-3-4-5-6-7 mnemonikası) şablonu istifadə olunur:
```go
t, _ := time.Parse("2006-02-01 15:04:05 -0700", "2016-13-03 00:00:00 +0000")
fmt.Println(t.Format("January 2, 2006 at 3:04:05PM MST"))
```
Eyni şablon həm parse, həm formatdır; hazır konstantlar (`time.RFC3339` və s.) var.

**Monotonic saat:** wall clock (DST, leap second, NTP ilə sıçrayır) + monotonic clock
(boot-dan sayır). Go `time.Now()` ilə yaradılan Time-larda monotonic oxu da saxlayır —
timer-lər və `Sub` avtomatik monotonic istifadə edir. Sıçrayış bug-ları (Cloudflare
hadisəsi) belə qarşısı alınır.

**Timer-lər:** `time.After` (tək atış channel), `time.Tick` (təkrar) — Tick-in ticker-i
söndürülməz/GC-olunmaz: yalnız trivial proqramlarda; əks halda `time.NewTicker`
(Reset/Stop metodları ilə). `time.AfterFunc` — gecikmiş funksiya.

### 3. encoding/json
**Anlayışlar:** marshaling = Go → JSON; unmarshaling = JSON → Go.

**Struct tag-lər (metadata):**
```go
type Order struct {
    ID          string    `json:"id"`
    DateOrdered time.Time `json:"date_ordered"`
    CustomerID  string    `json:"customer_id"`
    Items       []Item    `json:"items"`
}
```
- Format: `tagName:"tagValue"` — boşluqla ayrılmış cütlər; tək sətir; `go vet` yoxlayır.
- Default davranış (tagsiz): unmarshal-da case-insensitive field match; marshal-da
  export adı olduğu kimi. Amma **həmişə açış tag yaz** (idiomatik).
- `json:"-"` — iqnor; `json:",omitempty"` — boşsa çıxarma (zərif yer: zero struct
  "boş" SAYILMIR, sıfır-uzunluq slice/map sayılır).
- Sadə maşınlar yox — struct tag-lər yalnız funksiyaya ötürüləndə emal olunur (Java
  annotation fəlsəfəsindən fərqli — explicit qalır).

**Unmarshal/Marshal:**
```go
var o Order
err := json.Unmarshal([]byte(data), &o)     // parametri doldurur (Reader kimi)
out, err := json.Marshal(o)                  // []byte qaytarır
```
Pointer ötürülür: (1) strukturun reuse idarəsi; (2) generics olmadan tip seçiminin
yeganə yolu (reflection ilə — Ch14).

**json.Encoder/json.Decoder — stream I/O:**
```go
err := json.NewEncoder(tmpFile).Encode(toFile)        // io.Writer-a birbaşa
err = json.NewDecoder(tmpFile2).Decode(&fromFile)     // io.Reader-dan birbaşa
```
Böyük fayl/HTTP body üçün ReadAll+Unmarshal zəncirindən üstündür.

**JSON stream (çoxdəyərli axın):**
```go
dec := json.NewDecoder(strings.NewReader(data))
for dec.More() {
    err := dec.Decode(&t)   // obyekt-obyekt
}
```
Böyük array-ları tam yaddaşa yükləmədən element-element oxumaq da mümkündür.

**Custom parsing:** `json.Marshaler` + `json.Unmarshaler` interfeysləri:
```go
type RFC822ZTime struct{ time.Time }        // embed — time metodları qalır
func (rt RFC822ZTime) MarshalJSON() ([]byte, error) { ... }    // value receiver
func (rt *RFC822ZTime) UnmarshalJSON(b []byte) error { ... }   // pointer receiver
```
Çatışmazlıq: wire format biznes strukturlarının tipinə təsir edir. **Ən yaxşı həll —
iki struktur:** JSON-universal struct ↔ biznes struct; konvertasiya kənarla əlaqəni
kəsir (bir az dublikatiq, amma biznes məntiqi protokoldan təcrid olunur).
`map[string]interface{}` — yalnız kəşfiyyət mərhələsində; konkret tipə keç.

**QADAĞA:** `encoding/gob` + `net/rpc` — Go-spesifik; gRPC kimi standart protokolları işlət
(başqa dillərə bağlılıq üçün).

### 4. net/http — Client
**Qayda:** `http.DefaultClient` (timeout-suz) production-da QADAĞAN — öz client-ini yarat
(tək instans kifayətdir, goroutine-lər arasında təhlükəsizdir):
```go
client := &http.Client{ Timeout: 30 * time.Second }
req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)  // body: io.Reader və ya nil
req.Header.Add("X-My-Client", "Learning Go")
res, err := client.Do(req)
defer res.Body.Close()
if res.StatusCode != http.StatusOK { ... }
var data Todo
err = json.NewDecoder(res.Body).Decode(&data)
```
`http.Get/Head/Post` paket funksiyaları — default client istifadə edir → QAÇIN.

### 5. net/http — Server
**Struktur:** `http.Server` + `http.Handler` interfeysi:
```go
type Handler interface {
    ServeHTTP(http.ResponseWriter, *http.Request)
}
```
ResponseWriter çağırış sırası: `Header()` ( başlıqlar) → `WriteHeader(status)`
(200-də kənarlaşdırıla bilər) → `Write(body)`.

**Server konfiqurasiyası:**
```go
s := http.Server{
    Addr:         ":8080",
    ReadTimeout:  30 * time.Second,   // zəruri — default timeout YOXDUR!
    WriteTimeout: 90 * time.Second,
    IdleTimeout:  120 * time.Second,
    Handler:      HelloHandler{},
}
err := s.ListenAndServe()   // http.ErrServerClosed — normal shutdown
```

**ServeMux (router):**
```go
mux := http.NewServeMux()
mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) { ... })
```
**QADAĞA:** `http.Handle/HandleFunc/ListenAndServe` paket səviyyəli funksiyaları
`http.DefaultServeMux`-u istifadə edir — server konfiqurasiyası mümkün deyil + üçüncü
tərəf kitabxanalar oraya öz handler-lərini qeyd edə bilər (gizli shared state).

**Nested mux (sub-router):**
```go
mux.Handle("/person/", http.StripPrefix("/person", personMux))
mux.Handle("/dog/", http.StripPrefix("/dog", dogMux))
```
`StripPrefix` — valideyn mux-un artıq emal etdiyi hissəni kəsir.

### 6. Middleware Pattern
**Nədir:** Cross-cutting concern-lər (auth, timing, logging) — `func(http.Handler) http.Handler`:

**Kitabdan kod nümunəsi:**
```go
func RequestTimer(h http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        h.ServeHTTP(w, r)
        end := time.Now()
        log.Printf("request time for %s: %v", r.URL.Path, end.Sub(start))
    })
}

func TerribleSecurityProvider(password string) func(http.Handler) http.Handler {
    return func(h http.Handler) http.Handler {          // konfiqurasiya → middleware
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            if r.Header.Get("X-Secret-Password") != password {
                w.WriteHeader(http.StatusUnauthorized)
                w.Write(securityMsg)
                return
            }
            h.ServeHTTP(w, r)
        })
    }
}
```
Axın: setup/check → keçmədisə cavab yaz + return → keçdisə `h.ServeHTTP(w, r)` →
sonrakı cleanup. Zəncirləmə: `terribleSecurity(RequestTimer(handler))` — ən xarici
birinci icra olunur. Bütün mux-a tətbiq: `wrappedMux := terribleSecurity(RequestTimer(mux))`.
Katmanlararası data — context (Ch12). Fonksiyon zənciri xoşlamırsansa — `alice` modulu;
router ehtiyacları üçün — `gorilla/mux` və ya `chi` (hər ikisi http.Handler ekosistemi ilə
idiomatik).

## Əsas terminlər
- io.Reader/Writer — buffer parametrli oxuma/yazma interfeysləri
- io.EOF — axın sonu sentinel error-u (data ilə birlikdə gələ bilər)
- Decorator pattern — Reader-in Reader ilə örtülməsi (gzip və s.)
- Monotonic clock — sıçramayan, boot-dan sayılan saat
- Reference time — "Jan 2 15:04:05 2006 MST" şablon dili
- Struct tag — `json:"name"` metadata sətri
- Marshaling/Unmarshaling — Go↔JSON çevrilməsi
- Middleware — `func(http.Handler) http.Handler` sarğı funksiyası
- ServeMux — standart HTTP router

## Praktik nəticə

Standart kitabxana "batteries included"dir və özündə ən yaxşı dizayn nümunələrini daşıyır:
(1) funksiyalar interface qəbul etsin (`io.Reader`-ı `*os.File`-dan üstün tut);
(2) bufferləri sahibi idarə etsin; (3) Read-də error-dən əvvəl datanı emal et;
(4) HTTP client/server-da həmişə timeout qur; DefaultClient/DefaultServeMux-dan qaç;
(5) middleware üçün closure zənciri və ya alice; (6) JSON strukturunu biznes strukturu
ayır; (7) time müqayisəsi — `Equal`, format — 2006 şablonu; (8) Tick yerinə Ticker.

## Mənbə
Pages: 317-345
