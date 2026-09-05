# Go in Action — Xülasə (Azərbaycan dilində)

**Müəlliflər:** William Kennedy with Brian Ketelsen and Erik St. Martin
**Nəşriyyat:** Manning Publications Co., 2016 · **ISBN:** 9781617291784
**Səhifə:** 236 (+indeks) · **9 chapter**

## Kitabın ümumi məqsədi

Kitab intermediate səviyyəli developer-lər üçün yazılıb — Go dilini intensiv, kompleks və idiomatik şəkildə öyrədir: sintaksis, tip sistemi, concurrency, standart kitabxana, testing. Hər mövzu "nədir / necə işləyir / nəyə lazımdır" prinsipi ilə real kod nümunələri üzərində izah olunur.

## Chapter-by-chapter xülasə

### Ch 1 — Introducing Go (s. 1-8)
Go niyə yaradıldı: "sürətli inkişaf" vs "sürətli icra" dilemma-sını həll edir. 4 əsas sütun: (1) saniyə-altı kompilyasiya — kompilyator yalnız birbaşa import-ları yoxlayır; (2) concurrency — goroutine-lər bir OS thread-də minlərlə icra oluna bilir, logical processor-lərə schedule olunur; (3) hierarşiyasız tip sistemi — inheritance yox, kompozisiya + kiçik interfeyslər (`io.Reader` kimi); (4) garbage collector. Hello World + Go Playground ilə bitir.

### Ch 2 — Go quick-start (s. 9-38)
Tam RSS axtarış proqramı üzərində bütün əsas sintaksis: paketlər, `_` blank import (init-in çağırılması üçün), `init()` registration pattern-i, `:=` vs `var`, multi-return + error handling, `make` ilə map/channel yaratma, anonymous function + `go` ilə goroutine başlatma, WaitGroup sayğacı, closure tələsi (loop dəyişənlərini parametr kimi ötür!), `defer file.Close()`, struct tag ilə JSON/XML decode, `regexp.MatchString` ilə axtarış, map-lərdə 2 dəyərli lookup, `append` funksiyası.

### Ch 3 — Packaging and tooling (s. 39-56)
Paket qaydaları: bir qovluq = bir paket; `main` paketi + `main()` funksiyası = executable. Import axtarış sırası: GOROOT → GOPATH-lər. Remote import + `go get` (rekursiv). Named import (`myfmt "mylib/fmt"`). Go alətləri: `build`, `run`, `clean`, `vet` ( Printf xətaları, struct tag yoxlaması), `fmt` (avtomatik format), `doc`/`godoc`. Sənədləşmə konvensiyası: şərh identifikatorun üstündə, `doc.go` paket səviyyəli sənəd. Kod paylaşma: paket repo root-da. Dependency idarəsi: vendoring + import path rewriting (godep), gb (GOPATH-sız, path rewriting-sız).

### Ch 4 — Arrays, slices, maps (s. 57-87)
**Array:** sabit uzunluq, tipin hissəsidir (`[4]string` ≠ `[5]string`), value kimi kopyalanır — böyüklərində pointer ötür. **Slice:** 3 sahəli (pointer, len, cap) 24-byte görünüş; `make`, literal, nil/empty fərqi; `s[i:j]` → len=j-i, cap=k-i; üçüncü indeks `s[i:j:k]` cap məhdudlaşdırma → `append` detach pattern; `append` capacity 2x (≥1000 → 1.25x) böyüdür; `range` dəyərin **kopyasını** verir; funksiyaya ötürülməsi ucuzdur (underlying array kopyalanmır). **Map:** hash table + bucket; açar slice/funkdsiya ola bilməz; nil map-ə yazmaq panic; 2 dəyərli lookup; iterasiya sırası qeyri-müəyyən; funksiyaya ötürüləndə dəyişikliklər görünür.

### Ch 5 — Go's type system (s. 88-127)
Struct elanı + literal formaları; mövcud tip əsasında yeni tip (`type Duration int64` — fərqli tip, implicit convert yoxdur). Metodlar: value receiver (kopya) vs pointer receiver (paylaşım); kompilyator value/pointer-ı avtomatik uyğunlaşdırır. **Tiplərin təbiəti qaydası:** dəyişməz → value receiver (time.Time), dəyişən → pointer receiver (os.File) — metodun nə etdiyinə yox, tipin təbiətinə bax. İnterfeys internals: 2 sözlük (iTable + dəyər pointer). **Method sets:** value receiver → T və *T interfeysi implement edir; pointer receiver → yalnız *T (value-nun ünvanı həmişə alınmur). Polimorfizm nümunələri. **Type embedding:** tip adını sahə adı kimi yazma → promotion; xarici tip eyni metodu yazsa → override, amma iç tip birbaşa həmişə əlçatandır. Export/unexport: böyük/kiçik hərf; unexported tip + `New` factory funksiyası pattern; unexported iç tipin exported sahələri promotion ilə görünür.

### Ch 6 — Concurrency (s. 128-157)
Arxitektura: process → thread → logical processor (1:1 OS thread) → goroutine; scheduler global/local run queue idarə edir; blocking syscall → thread ayrılır, yeni thread yaranır; network I/O → poller; 10.000 thread limit. **Concurrency ≠ parallelism** — birincisi idarəetmə, ikincisi eyni anda icra. `runtime.GOMAXPROCS`. WaitGroup + `defer wg.Done()`. **Race condition:** kopya-oxu → yield → yazma nəticəsində itki; `go build -race` ilə aşkarlama. Həllər: atomic (AddInt64/LoadInt64/StoreInt64 — shutdown flag pattern), mutex (kritik bölmə), channel (CSP — "data-nı locklama, ötür"). Unbuffered channel: send/receive eyni anda kilidlənir (tennis, estafet nümunələri); buffered channel: receive yalnız boş kanalda, send yalnız dolu kanalda bloklanır (worker pool nümunəsi); qapalı kanaldan receive olmur — boşalana qədər.

### Ch 7 — Concurrency patterns (s. 158-183)
**Runner:** tapşırıqların ömür müddəti idarəsi — 3 kanal (interrupt buffered(1) + signal.Notify, complete unbuffered, `time.After` timeout), `select` 2 hadisəni gözləyir, `select+default` ilə nonblocking interrupt peek; ErrTimeout/ErrInterrupt error dəyişənləri + `switch err` + exit kodları. **Pool:** buffered channel resurs hovuzu — `Acquire` (select/default: boş resurs və ya factory), `Release` (select/default: geri qoy və ya Close), `Close` (mutex + closed flag + əvvəl close, sonra drain — deadlock qarşısı); type assertion ilə konkret tipə qayıdış. **Work:** unbuffered channel işçi pool — `Run` qayıdanda işin götürüldüyü **qarantiya** olunur (backpressure); `Shutdown` = close + Wait.

### Ch 8 — Standard library (s. 184-210)
100+ paket/38 kateqoriya, backward-compatibility zəmanəti. **log:** SetPrefix/SetFlags; iota + bit shift (`1 << iota`), flag-lərin `|` birləşməsi; `Println`/`Fatalln` (Exit(1))/`Panicln`; custom logger-lər — `log.New(io.Writer, prefix, flags)`, 4 səviyyə (Trace→Discard, Info→stdout, Error→MultiWriter(fayl, stderr)); logger-lər multigoroutine-safe. **json:** struct tag mapinqi; `NewDecoder(stream).Decode(&v)`; `Unmarshal([]byte, &v)`; struct bilmiriksə `map[string]interface{}` + type assertion; `MarshalIndent` (pretty) vs `Marshal` (sıx). **io:** Writer/Reader interfeysləri və davranış qaydaları (n < len(p) → error; n=0+err=nil qadağan); `fmt.Fprintf(&b, ...)` — Buffer da, File da Writer-dir; `io.MultiWriter`, `io.Copy`; sadə curl nümunəsi — 20 sətirdə URL → stdout + fayl.

### Ch 9 — Testing and benchmarking (s. 211-236)
Test konvensiyaları: `_test.go` faylı, `Test*` funksiyası, `*testing.T`; `t.Fatal` (dayandır) vs `t.Error` (davam et); "Given/When/Should" output formatı. Table test — anonim struct slice + `for range`. **Mocking:** `httptest.NewServer(http.HandlerFunc(f))` — internet olmadan HTTP test; handler: status + Content-Type + body. **Endpoint test:** `package handlers_test` (black-box), `http.NewRequest` + `httptest.NewRecorder` + `ServeHTTP` — server qaldırmadan handler test. **Examples:** `ExampleSendJSON` + `// Output:` şərhi — godoc-da görünür + çıxış yoxlanılır. **Benchmarks:** `Benchmark*` + `b *testing.B`, `for i := 0; i < b.N; i++`, `b.ResetTimer()`; `-bench`, `-benchtime=3s`, `-benchmem`; int→string müqayisəsi: `strconv.FormatInt` 45.9 ns/op vs `fmt.Sprintf` 258 ns/op (~5x) + allokasiya fərqi (2 B vs 16 B).

## Kitabın əsas mesajları

1. **Go sadəliyi dəbdən deyil, mühəndislikdən gəlir** — az açar söz, tək build yolu, vahid format (`go fmt`).
2. **Concurrency dilin génindədir** — `go` açar sözü + channel; shared memory + lock əvəzinə CSP mesaj modeli.
3. **Interfeyslər kiçikdir** — 1 metodlu `-er` interfeysləri kod reuse-un mühərrikidir; implementasiya elan tələb etmir.
4. **Kompozisiya irsiyyəti əvəz edir** — embedding + promotion, tipin təbiəti receiver seçimini müəyyən edir.
5. **Standard library etibarlıdır** — network-dən şifrələməyə qədər hər şey; onun source kodu idiomatik Go dərsliyidir.
6. **Testing inkişafın hissəsidir** — mock, table test, example, benchmark hamisi daxili tooling-lə.

## Ən dəyərli 5 fikir (az oxucu üçün)

1. **Race condition-lər görünməzdir** — concurrent kodda `-race` flag-i vərdiş olmalıdır; atomic/mutex/channel seçimi: int üçün atomic, blok üçün mutex, məntiq üçün channel.
2. **Closure tələsi:** loop daxilində goroutine başladanda dəyişənləri parametr kimi ötür — yoxsa bütün goroutine-lər son loop dəyərini görəcək.
3. **Slice paylaşımı:** iki slice eyni underlying array-i görür; `s[i:j:j]` detach pattern append-in mənbəyi korlamamasının qarantiyasıdır.
4. **Method sets qanunu:** pointer receiver → interfeysə yalnız pointer; bunu bilmədən "user does not implement..." compile xətası başa düşülməz qalır.
5. **İnterfeys-le kod yaz, konkret tip-lə yox:** `io.Copy`, `json.NewDecoder`, `log.New` hamisi Writer/Reader qəbul edir — öz tipinə bu interfeysləri ver, bütün ekosistem pulsuz açılır.
