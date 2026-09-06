# Chapter 2 — Go quick-start

## Bu chapter nədən bəhs edir?

Tam funksional bir Go proqramının (RSS feed axtarış sistemi) kodunu sətir-sətir izah edir: paketlər, import, `init`, `main`, variable declaration, pointer-lər, anonymous function + goroutine, WaitGroup, channel, defer, interface implementasiyası, JSON/XML decoding, regexp axtarışı. Bu chapter Go-nun bütün əsas sintaksisini real layihə üzərində göstərir.

Layihə strukturu (kitab reposu: https://github.com/goinaction/code/tree/master/chapter2/sample):
```
cd $GOPATH/src/github.com/goinaction/code/chapter2
- sample
  - data
    data.json        ← data feed-lərin siyahısı
  - matchers
    rss.go           ← RSS feed axtarışı üçün matcher
  - search
    default.go       ← default matcher
    feed.go          ← JSON data faylının oxunması
    match.go         ← interfeys dəstəyi
    search.go        ← əsas proqram məntiqi
  main.go            ← giriş nöqtəsi
```

## Əsas fikirlər

### 1. Paketlər və `main` paketi
**Nədir:** Go-da hər kod faylı bir paketə məxsusdur. Paket = kompilyasiya olunan kod vahidi, namespace rolunu oynayır.

**Necə işləyir:** İcra olunan proqram üçün `main` funksiyası **`main` paketində** olmalıdır, yoxsa build alətləri executable (icra olunan fayl) hasil etmir. Eyni qovluqdakı bütün fayllar eyni paket adını daşıyır (adətən qovluq adı ilə eyni).

### 2. Import və blank identifier
**Kitabdan kod nümunəsi (main.go):**
```go
package main

import (
    "log"
    "os"

    _ "github.com/goinaction/code/chapter2/sample/matchers"
    "github.com/goinaction/code/chapter2/sample/search"
)

// init is called prior to main.
func init() {
    // Change the device for logging to stdout.
    log.SetOutput(os.Stdout)
}

// main is the entry point for the program.
func main() {
    // Perform the search for the specified term.
    search.Run("president")
}
```

**Sub-kod izahı:**
- `package main` → icra olunan proqramın paketi
- `import (...)` → xarici kod istifadəsi; standart kitabxana üçün yalnız paket adı, xarici kod üçün tam path
- `_ "…/matchers"` → **blank identifier**: paketdən heç bir identifikator istifadə olunmasa da import qəbul edilir və paketin `init` funksiyaları çağırılır. Go kompilyatoru istifadə olunmayan import-a icazə vermir — `_` bu qadağanı aradan qaldırır
- `func init()` → `main`-dən əvvəl avtomatik çağırılır; burada logger-in çıxışını stdout-a yönləndirir (default stderr-dir)
- `search.Run("president")` → biznes məntiqi buradan başlayır; Run qayıdanda proqram bitir

**Vacib:** Bütün `init` funksiyaları (istənilən paketdə) `main`-dən əvvəl çağırılır. Bu, registration (qeydiyyat) pattern-inin təməlidir.

### 3. Variable declaration — `var`, `:=`, zero value, `make`
**Nədir:** Go-da dəyişən elan etməyin bir neçə yolu var; initiazlaşdırılmayan hər dəyişən **zero value** alır.

**Necə işləyir:**
- Rəqəmsel tiplər → `0`; string → `""`; bool → `false`; pointer → `nil`
- Map, slice, channel **reference type**-dır — istifadədən əvvəl `make` ilə yaradılmalıdır, yoxsa nil-dirlər

**Kitabdan kod nümunəsi:**
```go
// A map of registered matchers for searching.
var matchers = make(map[string]Matcher)
```

**Sub-kod izahı:**
- `var matchers` → package-level dəyişən (funksiya xaricində)
- `make(map[string]Matcher)` → map-i yaradir; `make`-siz map istifadəsi runtime xəta verir
- `matchers` kiçik hərflə başlayır → **unexported** (yalnız paket daxilində görünür). Böyük hərf = exported. Amma unexported tip də funksiyadan return oluna bilər

**Qayda:** zero value üçün `var` istifadə et; funksiya çağırışı/ilkin dəyər varsa `:=` (short variable declaration) istifadə et.

### 4. Multi-return və error handling
**Nədir:** Go funksiyaları birdən çox dəyər qaytara bilər; konvensiya — sonuncu dəyər `error` olur.

**Kitabdan kod nümunəsi:**
```go
feeds, err := RetrieveFeeds()
if err != nil {
    log.Fatal(err)
}
```

**Sub-kod izahı:**
- `feeds, err := RetrieveFeeds()` → `:=` eyni anda elan + initialize edir; tipi kompilyator çıxarır
- `log.Fatal(err)` → xətanı loglayıb proqramı **sonlandırır**
- **Qayda:** xəta olanda qaytarılan digər dəyərlərə etibar etmə — ya görmezlikən, ya da xətanı dərhal emal et

### 5. Goroutine-lərin başladılması + WaitGroup
**Nədir:** `go` açar sözü ilə funksiya (anonymous də olar) paralel işə salınır. `sync.WaitGroup` — sayğaclı semafor — bütün goroutine-lər bitənə qədər gözləmək üçün.

**Necə işləyir:** `waitGroup.Add(n)` sayğacı qoyur, hər goroutine işi bitincə `waitGroup.Done()` çağırır (sayğacı azaldır), `waitGroup.Wait()` sayğac sıfırlanana qədər bloklayır. `main` qayıdanda proqram bitir — hələ işləyən goroutine-lər runtime tərəfindən öldürülür, ona görə də təmiz shutdown vacibdir.

**Kitabdan kod nümunəsi (Run funksiyası — əsas kontrol məntiqi):**
```go
func Run(searchTerm string) {
    // Retrieve the list of feeds to search through.
    feeds, err := RetrieveFeeds()
    if err != nil {
        log.Fatal(err)
    }

    // Create a unbuffered channel to receive match results.
    results := make(chan *Result)

    // Setup a wait group so we can process all the feeds.
    var waitGroup sync.WaitGroup

    // Set the number of goroutines we need to wait for while
    // they process the individual feeds.
    waitGroup.Add(len(feeds))

    // Launch a goroutine for each feed to find the results.
    for _, feed := range feeds {
        matcher, exists := matchers[feed.Type]
        if !exists {
            matcher = matchers["default"]
        }

        // Launch the goroutine to perform the search.
        go func(matcher Matcher, feed *Feed) {
            Match(matcher, feed, searchTerm, results)
            waitGroup.Done()
        }(matcher, feed)
    }

    // Launch a goroutine to monitor when all the work is done.
    go func() {
        // Wait for everything to be processed.
        waitGroup.Wait()

        // Close the channel to signal to the Display
        // function that we can exit the program.
        close(results)
    }()

    // Start displaying results as they are available and
    // return after the final result is displayed.
    Display(results)
}
```

**Sub-kod izahı:**
- `results := make(chan *Result)` → unbuffered channel — nəticələrin goroutine-lərdən main-ə axını
- `waitGroup.Add(len(feeds))` → hər feed üçün 1 sayğac
- `for _, feed := range feeds` → `range` slice üzərində index + value qaytarır; `_` blank identifier index-i ignorə edir
- `matcher, exists := matchers[feed.Type]` → map lookup-un 2 dəyərli forması: value + "açar var mı?" bool. Açar yoxdursa map zero value qaytarır; burada default matcher-ə fallback edilir
- `go func(matcher Matcher, feed *Feed) { … }(matcher, feed)` → anonymous function (adsız funksiya) **parametrlərlə** goroutine kimi başladılır
- `(matcher, feed)` → dəyərlər funksiyaya **kopya** olunaraq ötürülür — bu, closure tələsindən qaçmaq üçün vacibdir
- `waitGroup.Wait()` + `close(results)` → bütün axtarışlar bitəndə channel bağlanır

**Closure tələsi (kitabın vacib qeydi):** `searchTerm` və `results` closure vasitəsilə anonymous funksiya üçün əlçatandır (birbaşa, kopyasız). Amma `matcher` və `feed` loop daxilində hər iterasiyada **dəyişir** — closure istifadə olunsaydı, bütün goroutine-lər eyni (ən sonuncu) feed-i emal edərdi. Ona görə parametr kimi ötürülür.

### 6. Pointer-lər və pass-by-value
**Nədir:** Go-da bütün dəyişənlər funksiyalara **dəyər üzrə** (by value) ötürülür. Pointer dəyişənin dəyəri = yaddaş ünvanıdır.

**Necə işləyir:** Pointer-lər funksiyalar və goroutine-lər arasında dəyişən paylaşmağa imkan verir — funksiya başqa funksiyanın scope-unda elan olunmuş dəyişənin vəziyyətini dəyişə bilər.

**Qayda:** Dəyişdiriləcək state üçün pointer receiver, yalnız davranış lazım olan tiplər üçün value receiver.

### 7. Interface — Matcher nümunəsi
**Nədir:** İnterfeys — davranış müqaviləsi. Tiptən tələb olunan metodları bildirir; implementasiya heç bir elan tələb etmir (duck typing).

**Kitabdan kod nümunəsi:**
```go
// Result contains the result of a search.
type Result struct {
    Field   string
    Content string
}

// Matcher defines the behavior required by types that want
// to implement a new search type.
type Matcher interface {
    Search(feed *Feed, searchTerm string) ([]*Result, error)
}
```

**Sub-kod izahı:**
- `type Matcher interface { … }` → interfeys tipi; tək `Search` metodu tələb edir
- Adlandırma konvensiyası: tək metodlu interfeys `-er` şəkilçisi ilə bitir (`Reader`, `Writer`, `Matcher`); çox metodlu interfeysin adı ümumi davranışı əks etdirir
- `Search(feed *Feed, searchTerm string) ([]*Result, error)` → pointer Feed qəbul edir, Result pointer slice + error qaytarır

### 8. defaultMatcher — interfeys implementasiyası
**Kitabdan kod nümunəsi:**
```go
// defaultMatcher implements the default matcher.
type defaultMatcher struct{}

// init registers the default matcher with the program.
func init() {
    var matcher defaultMatcher
    Register("default", matcher)
}

// Search implements the behavior for the default matcher.
func (m defaultMatcher) Search(feed *Feed, searchTerm string) ([]*Result, error) {
    return nil, nil
}
```

**Sub-kod izahı:**
- `type defaultMatcher struct{}` → **empty struct** (boş struktur) — 0 byte yaddaş ayırır; state lazım olmayanda idealdır
- `func init()` → proqram başlananda default matcher-i qeydiyyata alır
- `func (m defaultMatcher) Search(...)` → **value receiver** ilə metod — `Search` artıq `defaultMatcher` tipinə bağlıdır

**Receiver qaydaları (kitabın vacib cədvəli):**
| Receiver növü | Value-dan çağırış | Pointer-dan çağırış | Interface-də value | Interface-də pointer |
|---|---|---|---|---|
| Value receiver `func (m T) M()` | ✅ | ✅ (deref edir) | ✅ | ✅ |
| Pointer receiver `func (m *T) M()` | ✅ (ref edir) | ✅ | ❌ compile xətası | ✅ |

Pointer receiver ilə elan edilmiş metodu value interface-ə qoymaq **compile xətası verir**:
```
> go build
cannot use dm (type defaultMatcher) as type Matcher in assignment
```

**Ən yaxşı praktika:** ümumiyyətlə pointer receiver istifadə et; yalnız zero-allocation tiplər üçün value receiver.

### 9. Match və Display — channel istifadəsi
**Kitabdan kod nümunəsi:**
```go
// Match is launched as a goroutine for each individual feed to run
// searches concurrently.
func Match(matcher Matcher, feed *Feed, searchTerm string,
    results chan<- *Result) {

    // Perform the search against the specified matcher.
    searchResults, err := matcher.Search(feed, searchTerm)
    if err != nil {
        log.Println(err)
        return
    }

    // Write the results to the channel.
    for _, result := range searchResults {
        results <- result
    }
}

// Display writes results to the terminal window as they
// are received by the individual goroutines.
func Display(results chan *Result) {
    // The channel blocks until a result is written to the channel.
    // Once the channel is closed the for loop terminates.
    for result := range results {
        fmt.Printf("%s:\n%s\n\n", result.Field, result.Content)
    }
}
```

**Sub-kod izahı:**
- `results chan<- *Result` → **send-only channel** — bu funksiya yalnız channel-a yaza bilər
- `matcher.Search(...)` → interfeys üzərindən konkret implementasiya çağırılır (polimorfizm)
- `results <- result` → channel-a nəticə göndərir; alan tərəf yoxdursa bloklanır
- `for result := range results` → channel-da nəticə gələnə qədər bloklanır; **channel bağlananda loop bitir** — sonsuz loop görünür, amma `close(results)` ilə qurtarır
- `range` kanallarla da işləyir — array, slice, map, string, channel-in hamısında

### 10. Register — init-time qeydiyyat
**Kitabdan kod nümunəsi:**
```go
// Register is called to register a matcher for use by the program.
func Register(feedType string, matcher Matcher) {
    if _, exists := matchers[feedType]; exists {
        log.Fatalln(feedType, "Matcher already registered")
    }

    log.Println("Register", feedType, "matcher")
    matchers[feedType] = matcher
}
```

**Sub-kod izahı:**
- Eyni tip iki dəfə qeydiyyatdan keçsə, `log.Fatalln` proqramı dayandırır
- `matchers[feedType] = matcher` → map-ə əlavə; bütün qeydiyyat `main`-dən əvvəl `init`-lər vasitəsilə baş verir

### 11. JSON decoding — feed.go
**Kitabdan kod nümunəsi (data.json):**
```json
[
    {
        "site" : "npr",
        "link" : "http://www.npr.org/rss/rss.php?id=1001",
        "type" : "rss"
    },
    {
        "site" : "cnn",
        "link" : "http://rss.cnn.com/rss/cnn_world.rss",
        "type" : "rss"
    }
]
```

**Kitabdan kod nümunəsi (struct + RetrieveFeeds):**
```go
const dataFile = "data/data.json"

// Feed contains information we need to process a feed.
type Feed struct {
    Name string `json:"site"`
    URI  string `json:"link"`
    Type string `json:"type"`
}

// RetrieveFeeds reads and unmarshals the feed data file.
func RetrieveFeeds() ([]*Feed, error) {
    // Open the file.
    file, err := os.Open(dataFile)
    if err != nil {
        return nil, err
    }

    // Schedule the file to be closed once
    // the function returns.
    defer file.Close()

    // Decode the file into a slice of pointers
    // to Feed values.
    var feeds []*Feed
    err = json.NewDecoder(file).Decode(&feeds)

    // We don't need to check for errors, the caller can do this.
    return feeds, err
}
```

**Sub-kod izahı:**
- `` `json:"site"` `` → **struct tag**: JSON-dakı `site` açarını `Name` sahəsinə map edir
- `const dataFile` → tip kompilyator tərəfindən sağ tərəfdən çıxarılır; kiçik hərf = unexported
- `os.Open(dataFile)` → fayl açır, `*File` + error qaytarır
- `defer file.Close()` → **defer**: funksiya qayıdandan dərhal sonra icra olunacaq — panic belə olsa Close qaçmır; açılışla yaxın yazıldığı üçün oxunaqlılıq artır
- `json.NewDecoder(file).Decode(&feeds)` → faylı streaming şəkildə decode edir; `Decode` `interface{}` (boş interfeys — istənilən tip) qəbul edir, reflection ilə işləyir

### 12. RSS matcher — XML decode + HTTP + regexp
**Kitabdan kod nümunəsi (XML struct-ları):**
```go
type (
    // item defines the fields associated with the item tag
    // in the rss document.
    item struct {
        XMLName     xml.Name `xml:"item"`
        PubDate     string   `xml:"pubDate"`
        Title       string   `xml:"title"`
        Description string   `xml:"description"`
        Link        string   `xml:"link"`
        GUID        string   `xml:"guid"`
        GeoRssPoint string   `xml:"georss:point"`
    }

    // channel defines the fields associated with the channel tag
    // in the rss document.
    channel struct {
        XMLName xml.Name `xml:"channel"`
        Title   string   `xml:"title"`
        // ... digər sahələr
        Image   image    `xml:"image"`
        Item    []item   `xml:"item"`
    }

    // rssDocument defines the fields associated with the rss document
    rssDocument struct {
        XMLName xml.Name `xml:"rss"`
        Channel channel  `xml:"channel"`
    }
)
```

**Sub-kod izahı:**
- XML decoding JSON ilə eyni məntiqlədir — struct tag-lər `xml:"..."` formatında
- `XMLName xml.Name` → sahənin uyğun gəldiyi XML elementini göstərir
- `Item []item` → nested slice — bir channel-da çox item ola bilər

**Kitabdan kod nümunəsi (retrieve — HTTP GET):**
```go
func (m rssMatcher) retrieve(feed *search.Feed) (*rssDocument, error) {
    if feed.URI == "" {
        return nil, errors.New("No rss feed URI provided")
    }

    // Retrieve the rss feed document from the web.
    resp, err := http.Get(feed.URI)
    if err != nil {
        return nil, err
    }

    // Close the response once we return from the function.
    defer resp.Body.Close()

    // Check the status code for a 200 so we know we have received a
    // proper response.
    if resp.StatusCode != 200 {
        return nil, fmt.Errorf("HTTP Response Error %d\n", resp.StatusCode)
    }

    // Decode the rss feed document into our struct type.
    var document rssDocument
    err = xml.NewDecoder(resp.Body).Decode(&document)
    return &document, err
}
```

**Sub-kod izahı:**
- `http.Get(feed.URI)` → HTTP sorğusu; `*Response` + error qaytarır
- `defer resp.Body.Close()` → body mütləq bağlanmalıdır — defer bunu qarantiya edir
- `resp.StatusCode != 200` → status yoxlanışı; `fmt.Errorf` formatlanmış xəta yaradır
- `xml.NewDecoder(resp.Body).Decode(&document)` → HTTP cavabını birbaşa XML struct-a decode edir

**Kitabdan kod nümunəsi (Search — regexp axtarışı):**
```go
func (m rssMatcher) Search(feed *search.Feed, searchTerm string) ([]*search.Result, error) {
    var results []*search.Result

    log.Printf("Search Feed Type[%s] Site[%s] For Uri[%s]\n", feed.Type, feed.Name, feed.URI)

    // Retrieve the data to search.
    document, err := m.retrieve(feed)
    if err != nil {
        return nil, err
    }

    for _, channelItem := range document.Channel.Item {
        // Check the title for the search term.
        matched, err := regexp.MatchString(searchTerm, channelItem.Title)
        if err != nil {
            return nil, err
        }

        // If we found a match save the result.
        if matched {
            results = append(results, &search.Result{
                Field:   "Title",
                Content: channelItem.Title,
            })
        }

        // Check the description for the search term.
        matched, err = regexp.MatchString(searchTerm, channelItem.Description)
        if err != nil {
            return nil, err
        }

        // If we found a match save the result.
        if matched {
            results = append(results, &search.Result{
                Field:   "Description",
                Content: channelItem.Description,
            })
        }
    }

    return results, nil
}
```

**Sub-kod izahı:**
- `var results []*search.Result` → nil slice elanı
- `regexp.MatchString(searchTerm, …)` → regex qarşılaşdırması; `matched, err :=` iki dəyərli qayıdış
- `results = append(results, &search.Result{…})` → `append` slice-in uzunluğunu tutumunu avtomatik böyüdür; `&` yeni struct-literal-in ünvanını slice-a qoyur
- Həm `Title`, həm `Description` sahələri yoxlanılır — uyğun gələnlər `Field` adı ilə nəticəyə düşür

## Əsas terminlər
- Package (paket — kod vahidi/namespace)
- Import (idxal)
- Blank Identifier (boş identifikator `_`)
- Zero Value (sıfır dəyər)
- Reference Type (istinad tipi — map, slice, channel)
- Anonymous Function (anonim funksiya)
- Closure (klouzer — xarici scope dəyişənlərinə birbaşa çıxış)
- WaitGroup (gözləmə qrupu — sayğaclı semafor)
- Unbuffered Channel (bufersız kanal)
- Interface (interfeys)
- Value/Pointer Receiver (dəyər/göstərici qəbuledici)
- Struct Tag (struktur teqi)
- defer (təxirə salınmış çağırış)
- Pass by Value (dəyərlə ötürmə)

## Praktik nəticə
- Framework-ləri interfeys + init-registration pattern-i ilə qur: yeni matcher əlavə etmək üçün `matchers/` qovluğuna yeni fayl + `init()` bəsdir, `search` paketi dəyişmir.
- Döngü daxilində goroutine başladanda dəyişənləri **parametr kimi ötür** — closure tələsindən qaç.
- `main` qayıtmazdan əvvəl bütün goroutine-lərin bitməsini WaitGroup + channel close ilə təmin et.
- Resurs açılanda `defer`-i dərhal yanına yaz.
- Xətaları dərhal yoxla; xəta varsa qaytarılan dəyərlərə etibar etmə.

## Mənbə
Pages: 30-59 (PDF), book pages 9-38
