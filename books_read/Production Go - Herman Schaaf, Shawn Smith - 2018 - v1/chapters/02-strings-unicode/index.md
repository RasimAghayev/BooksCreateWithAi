# Chapters 5-7 — Strings, Supporting Unicode, Concurrency

## Bu bölmələr nədən bəhs edir?

strings paketi (Split/Count/Index/Contains/FieldsFunc/EqualFold), Unicode və UTF-8 dərinliyi (byte slice semantikası, rune, printf verb-ləri), concurrency (WaitGroup, kanallar, web handler-larda goroutine, ticker poller, race condition + sync.Map vs mutex-ətri safeMap).

**PDF səhifələr:** 58-87 (Strings: 58-68, Unicode: 68-73, Concurrency: 74-87)

## Əsas fikirlər

### Ch 5 — Strings

**String yaratma — 3 yol:**
```go
s := "I am a string - 你好"          // literal (Çin heroqlifləri də LEGAL)
greeting := "Hello, my name is "
greeting += "Inigo Montoya"           // + / += — sadə hallarda
sentence := fmt.Sprintf("Hello, my name is %s.", name)   // TÖVSİYƏ
```
- `+` yalnız FEW string, hot path YOX (optimization chapterda müzakirə)
- **`+` int-lə işləməz:** `"I am" + 32` → `cannot convert "I am" to type int` — Sprintf istifadə et

**fmt.Sprintf verb-ləri:**
```go
fmt.Printf("Hello, my name is %s, age %d, weight %.2fkg", name, age, weight)
// %s string, %d int, %f float, %v Go özü seçir, %.2f 2 onluq
```

**strings paketi əsasları:**
```go
strings.Split("a,b,c", ",")              // ["a" "b" "c"]
strings.Count("banana", "a")             // 3
strings.Count("banana", "ana")          // 1 — NON-OVERLAPPING! (ana 2x var amma overlap)
strings.Contains(str, "moon")            // true — bool
strings.HasPrefix / strings.HasSuffix
strings.Index("banana", "an")            // 1; yoxdursa -1
strings.LastIndex                        // sonuncu
strings.ToLower / ToUpper / Trim / Join
```
**Qeyd:** Go-da overloading YOX — tək hərf də STRING kimi ötür.

**Advanced nümunə — palindrome:**
```go
func isNotLetter(c rune) bool {
    return !unicode.IsLetter(c)
}

func isPalindromicSentence(s string) bool {
    // FieldsFunc — funksiya İLE split (funksiyalar dəyər kimi!)
    w := strings.FieldsFunc(s, isNotLetter)     // qeyri-hərf simvollarla böl → SÖZLƏR

    l := len(w)
    for i := 0; i < l/2; i++ {
        fw := w[i]       // front word
        bw := w[l-i-1]   // back word
        if !strings.EqualFold(fw, bw) {          // case-insensitive müqayisə
            return false
        }
    }
    return true
}

// for loop-un 3 komponentindən istifadə (while əvəzi):
for l := getInput(); l != ""; l = getInput() { ... }
```

**String range — index tələsi:**
```go
s := "ABC你好"
for i, r := range s {
    fmt.Printf("%q(%d) ", r, i)
}
// 'A'(0) 'B'(1) 'C'(2) '你'(3) '好'(6)  ← 3-dən 6-ya SıÇRAY!
```
Niyə? → Unicode chapter cavabı: **range RUNE-larla gəzir, index BAYT yeridir** — 你 3 baytdır (3,4,5), 好 Inhalt 6-dan.

### Ch 6 — Supporting Unicode

**Encoding tarixçəsi:**
- ASCII: 127 simvol, 7 bit — yalnız İngilis hərfləri
- Sonra: ölkə-əsaslı ENCODING xaosu
- **Unicode:** hər simvola nömrə (code point, U+0041 formatında) — lakin binary NECƏ saxlanacağını demir
- **UTF-8:** Unicode-un ən məşhur encodingi; 0-127 = 1 bayt (ASCII ilə EYNİ!); 128+ = 2-6 bayt
- Tövsiyə: Joel Spolsky "Absolute Minimum... Unicode" (2003)

**Strings are byte slices:**
```go
b := []byte{65, 66, 67}
s := string(b)          // "ABC" — Go heç bir encoding fərzi etmədi!
```
**Vacib:** Go string tipi encoding məlumatı DAŞIMIR — sadəcə read-only bayt slice. Çap edəndə terminal öz encodingi ilə render edir. Digər encodinglər: golang.org/x/text/encoding.

**Printf verb cheat (strings üçün):**
```go
s := "ABC 你好"
%s     → ABC 你好                  // kopyası kimi
%q     → "ABC 你好"                // quote ilə
%+q    → "ABC \u4f60\u597d"        // ASCII-ONLY output!
%x     → 41424320e4bda0e5a5bd      // hex (1 bayt = 2 simvol)
% x    → 41 42 43 20 e4 bd a0 e5 a5 bd   // boşluqlu
%# x   → 0x41 0x42 ...              // 0x prefiksli
```

**Ümumi verb-lər:** `%v` (default), `%+v` (struct sahə adları), `%#v` (Go representasiya), `%T` (tip), `%%` (faiz).
**Flag-lər:** `+` (işarət/ASCII %q), `-` (sol-justifi), `#` (0x/0 prefiks), `' '` (baytlar arası boşluq), `0` (sıfır pad).
**Genişlik RUNE ilə ölçülür** (C-də bayt ilə!) — fərq beynəlxalq mətndə əhəmiyyətli.

**Runes:** int32 alias — 1 Unicode simvol; range avtomatik decode edir (buna görə index bayt-enişli). `range` UTF-8 fərzi YEGANƏ yerdir — string tipinin özü encodingdən asılı deyil.

### Ch 7 — Concurrency

**Goroutine əsasları:**
```go
go func() { fmt.Println("hello, world") }()
// main çıxsa goroutine ölür → görüntü YOX
// 3 goroutine → SIRASI TƏSDÜFİ (hello world 2, 1, 3 ola bilər)
```
**Sleep → WaitGroup (düzgün yolu):**
```go
var wg sync.WaitGroup
wg.Add(1)
go func() {
    defer wg.Done()
    ...
}()
wg.Wait()
```
**Qeyd:** Add hər goroutine üçün ayrıca çağırıla bilər (Add(3) əvəzinə 3x Add(1)).

**Kanallar — 3 sorğu, 2 gözləmə:**
```go
ch := make(chan string)
go func(ch chan string) { ch <- "hello, world 1" }(ch)   // ×3
a, b, c := <-ch, <-ch, <-ch    // 3 oxu
```
2 göndər + 3 oxu → `fatal error: all goroutines are asleep - deadlock!`

**Goroutine web handler-da:**
```go
http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
    go func() {
        time.Sleep(3 * time.Second)
        log.Println("hello, world")    // handler qayıtdiqdan SONRA çap olunur!
    }()
    return                            // cavab DƏRHAL; goroutine fondda davam
})
```
→ Handler qayıdır, goroutine server process-də yaşayır — background task patterninin əsası.

**Ticker poller (production pattern):**
```go
func poll() {
    ticker := time.NewTicker(5 * time.Second)   // 5 saniyəlik ticker
    defer ticker.Stop()
    for {
        select {
        case <-ticker.C:                        // hər 5 saniyə
            resp, err := http.Get(bsaEndpoint)  // BART API
            if err != nil { log.Println("ERROR: ...", err) }
            defer resp.Body.Close()
            b, _ := ioutil.ReadAll(resp.Body)
            var br bsaResponse
            json.Unmarshal(b, &br)              // CDATA tag də istifadə olunur:
            for _, adv := range br.Root.Advisories {
                log.Println(adv.Station, adv.Description.Text)
            }
        }
    }
}

go poll()                            // sonsuz loop → goroutine ŞƏRT
log.Fatal(http.ListenAndServe(...))  // main web server
```
JSON tag: `` `json:"#cdata-section"` `` — XML CDATA-sının Go-dakı qarşılığı.

**Race condition:**
```go
m := map[string]string{}
go func() { m["test"] = "hello, world 1" }()
go func() { m["test"] = "hello, world 2" }()    // DATA RACE!
```

**Həll 1 — sync.Map:**
```go
var m = &sync.Map{}
m.Store("test", "hello, world 1")
val, _ := m.Load("test")
```
**sync.Map YALNIZ 2 halda (sənəddən!):**
1. Açar BİR DƏFƏ yazılır, çox oxunur (write-once, read-many cache)
2. Goroutine-lər AYRI açar dəstləri üzərində işləyir

**Həll 2 — Mutex + map (TÖVSİYƏ):**
```go
type safeMap struct {
    sync.Mutex
    m map[string]string
}

func (sm *safeMap) Store(key, val string) {
    sm.Lock()
    defer sm.Unlock()
    sm.m[key] = val
}

func (sm *safeMap) Load(key string) (val string, exists bool) {
    sm.Lock()
    defer sm.Unlock()
    val, ok := sm.m[key]
    return val, ok
}

func (sm *safeMap) Delete(key string) {
    sm.Lock()
    defer sm.Unlock()
    delete(sm.m, key)
}
```
→ `go run -race main.go` — no data races found ✓

## Əsas terminlər
- strings.Split / Count / Contains / HasPrefix / HasSuffix / Index / LastIndex
- FieldsFunc (funksiya ilə split)
- EqualFold (case-insensitive müqayisə)
- unicode.IsLetter
- for-un init/cond/post komponentləri (while əvəzi)
- ASCII / Unicode / Code Point (U+XXXX) / UTF-8
- Strings are read-only byte slices (encoding daşınmır)
- rune = int32 (4 bayt, 1 Unicode simvol)
- range-in UTF-8 fərzi (index BAYT sayır, dəyər RUNE)
- %+q (ASCII-only), %x/%X (hex), %# x (0x prefiksli)
- Genişlik RUNE ilə (C-də bayt ilə)
- WaitGroup (Add-per-goroutine)
- Kanal deadlock (send≠receive sayı)
- Handler-da fon goroutine
- time.NewTicker / ticker.C / ticker.Stop
- sync.Map (write-once / disjoint-keys)
- safeMap (Mutex embed + Store/Load/Delete)

## Praktik nəticə
- String birləşdirmə: bir neçə string üçün +; qarışıq tiplər və ya tez-tez → Sprintf; production hot-path → strings.Builder (benchmark chapterdə).
- Count NON-OVERLAPPING-di — "ananada ana 2x" gözləntisi yanlışdır.
- Range-in index-i baytdır, dəyəri rune — UTF-8 pozuntularında index sıçrayış normaldır.
- `+` operatoru tipləri qarışdırmır — int string-ə çevrilməz; Sprintf %d istifadə et.
- Web handler-da fon iş üçün handler qayıtmadan goroutine başlat — amma graceful shutdown-u unutma.
- Sonsuz loop (poller) həmişə goroutine-də; ticker + select + defer Stop əsası.
- sync.Map yalnız 2 sənəd halında; əks halda Mutex-ətri safeMap yarat (Store/Load/Delete metodları) — kitabın tövsiyəsi.
- Hər concurrent kodu `-race` ilə sına — "işləyir" görünən race kod əslində qırıqdır.

## Mənbə
Pages: 58-87 (PDF), book pages 52-81
