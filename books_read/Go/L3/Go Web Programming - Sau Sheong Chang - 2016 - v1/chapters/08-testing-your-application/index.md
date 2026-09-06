# Chapter 8 — Testing your application

## Bu chapter nədən bəhs edir?

Go testing: unit testlər (testing.T, Skip, parallel, benchmarking), HTTP handler testləri (httptest.NewRecorder + NewRequest), test doubles + dependency injection (Text interfeysi + FakePost), üçüncü tərəf kitabxanalar (gocheck — suite/fixture/assertion, Ginkgo+Gomega — BDD).

## Əsas fikirlər

### 1. Test əsasları
- Fayl adı **`_test.go`** ilə bitməli, test olunan kodun eyni paketində
- `func TestXxx(t *testing.T)` — Xxx böyük hərflə başlayan ixtiyari ad
- İşə salma: `go test` (hamısı), `go test -v` (verbose), `go test -v -cover` (coverage: `coverage: 46.7% of statements`)

### 2. testing.T funksiyaları
**Əsas:**
- `Log` / `Logf` → error log-a yaz (Fail ilə birlikdə görünür)
- `Fail` → testi FAILED et, **davam et**
- `FailNow` → testi FAILED et, **dayandır**

**Kombinasiya (convenience):**
```
Error/Errorf = Log + Fail      (fail + davam)
Fatal/Fatalf = Log + FailNow    (fail + dayan)
Skip                          (testi tam keç)
```

**Kitabdan kod nümunəsi:**
```go
func TestDecode(t *testing.T) {
    post, err := decode("post.json")
    if err != nil {
        t.Error(err)
    }
    if post.Id != 1 {
        t.Error("Wrong id, was expecting 1 but got", post.Id)
    }
    if post.Content != "Hello World!" {
        t.Error("Wrong content, was expecting 'Hello World!' but got", post.Content)
    }
}

func TestEncode(t *testing.T) {
    t.Skip("Skipping encoding for now")
}
```

**Testability haqqında vacib qeyd:** kodu **test üçün dizayn etmək** lazımdır — decode funksiyasını main-dən ayırmaq buna nümunədir: "it's equally important that the code be testable".

### 3. Skip və -short
```go
func TestLongRunningTest(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping long-running test in short mode")
    }
    time.Sleep(10 * time.Second)
}
```
- `go test -v -short` → uzun testlər keçilir (10s → 0.00s)
- `-short` fərqi: testin **hissələrini** şərtli keçmək; `-run` isə hansı testlərin işə düşəcəyini seçir

### 4. Paralel testlər
```go
func TestParallel_1(t *testing.T) {
    t.Parallel()                        // İLK statement olmalıdır
    time.Sleep(1 * time.Second)
}
```
- `go test -v -short -parallel 3` → maksimum 3 paralel
- Nəticə: 1+2+3 = 6 saniyelik iş → **3.006 saniyə** (ən uzun test qədər)

### 5. Benchmarking
```go
func BenchmarkDecode(b *testing.B) {
    for i := 0; i < b.N; i++ {
        decode("post.json")
    }
}
```
- `go test -bench .` → bütün benchmarklar; `-run x` → funksional testləri kənarlaşdır
- Output: `BenchmarkDecode  100000  19480 ns/op` — 100000 iterasiya, ~19.5µs/əməliyyat
- `b.N` — Go özü təyin edir (dəqiqlik üçün); müddəti `-benchtime` ilə təsir etmək olar

**Real müqayisə (Decode vs Unmarshal):**
```
BenchmarkDecode     100000    19577 ns/op
BenchmarkUnmarshal   50000    24532 ns/op   → ~25% daha YAVAŞ
```

### 6. HTTP testing — httptest
**Kitabdan kod nümunəsi (GET testi):**
```go
func TestHandleGet(t *testing.T) {
    mux := http.NewServeMux()                        // 1. mux yarat
    mux.HandleFunc("/post/", handleRequest)          // 2. handler-i bağla
    writer := httptest.NewRecorder()                 // 3. cavab qeydedicisi
    request, _ := http.NewRequest("GET", "/post/1", nil)  // 4. sorğu yarat
    mux.ServeHTTP(writer, request)                   // 5. sorğunu GÖNDƏRMƏDƏN icra et

    if writer.Code != 200 {                          // 6. yoxla
        t.Errorf("Response code is %v", writer.Code)
    }
    var post Post
    json.Unmarshal(writer.Body.Bytes(), &post)
    if post.Id != 1 {
        t.Error("Cannot retrieve JSON post")
    }
}
```

**Sub-kod izahı:**
- `httptest.NewRecorder()` → `ResponseRecorder` — response-u yadda saxlayan saxta ResponseWriter (network yoxdur!)
- `http.NewRequest(method, url, body)` → sorğu obyekti; body üçün `strings.NewReader(...)`
- `mux.ServeHTTP(writer, request)` → mux-u birbaşa çağır → handler icra olunur → cavab recorder-a düşür
- `writer.Code`, `writer.Body.Bytes()` → cavabı yoxla

**PUT testi (body ilə):**
```go
json := strings.NewReader(`{"content":"Updated post","author":"Sau Sheong"}`)
request, _ := http.NewRequest("PUT", "/post/1", json)
```

### 7. TestMain — mərkəzi setup/teardown
```go
func TestMain(m *testing.M) {
    setUp()
    code := m.Run()
    tearDown()
    os.Exit(code)
}

func setUp() {
    mux = http.NewServeMux()
    mux.HandleFunc("/post/", handleRequest)
    writer = httptest.NewRecorder()
}
```
- `setUp/tearDown` **bir dəfə** — bütün testlər üçün; `m.Run()` testləri işə salır, exit code qaytarır
- Qlobal dəyişənlər (mux, writer) testlər arasında paylaşılır → kod təkrarı aradan qalxır

### 8. Problem: DB asılılığı
Ch 7 servisinin testləri **retrīeve funksiyası** → **qlobal sql.DB** → PostgreSQL asılıdır. Test müstəqil deyil (DB söndürülsə fail).

### 9. Dependency injection — həll
**Addım 1 — Interfeys yarat (Text):**
```go
type Text interface {
    fetch(id int) (err error)
    create() (err error)
    update() (err error)
    delete() (err error)
}
```

**Addım 2 — Post interfeysi implement etsin + Db sahəsi:**
```go
type Post struct {
    Db      *sql.DB
    Id      int    `json:"id"`
    Content string `json:"content"`
    Author  string `json:"author"`
}

func (post *Post) fetch(id int) (err error) {
    err = post.Db.QueryRow("select id, content, author from posts where id = $1",
        id).Scan(&post.Id, &post.Content, &post.Author)
    return
}
```
- Db artıq **struct sahəsidir**, qlobal deyil — funksiyalar sql.DB-yə birbaşa asılı deyil

**Addım 3 — Handler-i dəyiş → closure qaytar:**
```go
func handleRequest(t Text) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var err error
        switch r.Method {
        case "GET":
            err = handleGet(w, r, t)
        case "POST":
            err = handlePost(w, r, t)
        case "PUT":
            err = handlePut(w, r, t)
        case "DELETE":
            err = handleDelete(w, r, t)
        }
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
    }
}

func handleGet(w http.ResponseWriter, r *http.Request, post Text) (err error) {
    id, err := strconv.Atoi(path.Base(r.URL.Path))
    if err != nil {
        return
    }
    err = post.fetch(id)          // interfeys üzərindən!
    if err != nil {
        return
    }
    output, err := json.MarshalIndent(post, "", "\t\t")
    if err != nil {
        return
    }
    w.Header().Set("Content-Type", "application/json")
    w.Write(output)
    return
}
```

**Addım 4 — main-da inject et:**
```go
func main() {
    var err error
    db, err := sql.Open("postgres", "user=gwp dbname=gwp password=gwp sslmode=disable")
    if err != nil {
        panic(err)
    }
    server := http.Server{Addr: ":8080"}
    http.HandleFunc("/post/", handleRequest(&Post{Db: db}))   // Dİ!
    server.ListenAndServe()
}
```
- `handleRequest(&Post{Db: db})` → HandlerFunc qaytaran closure → HandleFunc imzasına uyğun (Ch 3 chaining texnikası)

**Addım 5 — Test double (FakePost):**
```go
type FakePost struct {
    Id      int
    Content string
    Author  string
}

func (post *FakePost) fetch(id int) (err error) {
    post.Id = id        // test üçün sadəcə id qaytar
    return
}
func (post *FakePost) create() (err error) { return }
func (post *FakePost) update() (err error) { return }
func (post *FakePost) delete() (err error) { return }
```

**Addım 6 — Test-də FakePost inject et:**
```go
mux.HandleFunc("/post/", handleRequest(&FakePost{}))
```
→ DB söndürülmüş olsa belə test KEÇİR — handler müstəqil test olunur. (fetch-in öz DB testləri ayrıca suite-də.)

### 10. gocheck (gopkg.in/check.v1)
**Xüsusiyyətlər:** suite əsaslı qruplaşdırma, test fixture-ləri, assertion-lar (genişlənən checker), daha yaxşı error reporting, testing ilə sıx inteqrasiya.

**Quraşdırma:** `go get gopkg.in/check.v1`

**Kitabdan kod nümunəsi:**
```go
import (
    . "gopkg.in/check.v1"    // dot import — qualifier-siz
)

type PostTestSuite struct {
    mux    *http.ServeMux
    post   *FakePost
    writer *httptest.ResponseRecorder
}

func init() {
    Suite(&PostTestSuite{})         // suite-i qeydiyyata al
}

func Test(t *testing.T) { TestingT(t) }   // testing ilə körpü

func (s *PostTestSuite) SetUpTest(c *C) {  // hər testdən ƏVVƏL
    s.post = &FakePost{}
    s.mux = http.NewServeMux()
    s.mux.HandleFunc("/post/", handleRequest(s.post))
    s.writer = httptest.NewRecorder()
}

func (s *PostTestSuite) TestGetPost(c *C) {
    request, _ := http.NewRequest("GET", "/post/1", nil)
    s.mux.ServeHTTP(s.writer, request)
    c.Check(s.writer.Code, Equals, 200)       // assertion: fail + DAVAM
    var post Post
    json.Unmarshal(s.writer.Body.Bytes(), &post)
    c.Check(post.Id, Equals, 1)
}

func (s *PostTestSuite) TestPutPost(c *C) {
    json := strings.NewReader(`{"content":"Updated post","author":"Sau Sheong"}`)
    request, _ := http.NewRequest("PUT", "/post/1", json)
    s.mux.ServeHTTP(s.writer, request)
    c.Check(s.writer.Code, Equals, 200)
    c.Check(s.post.Id, Equals, 1)
    c.Check(s.post.Content, Equals, "Updated post")
}
```

**Sub-kod izahı:**
- `Suite(&PostTestSuite{})` → qeydiyyat; `TestXxx` metodları test kimi icra olunur
- `c.Check(dəyər, Equals, gözlənilən)` → assertion; `c.Assert` → fail + **dayandır** (Check fail + davam)
- Fixture lifecycle: `SetUpSuite`/`TearDownSuite` (suite başında/sonunda 1 dəfə), `SetUpTest`/`TearDownTest` (hər testin ətrafında)
- İşə salma: `go test -check.vv` → ətraflı output; error halında:
```
... obtained int = 0
... expected int = 1
```

### 11. Ginkgo + Gomega (BDD)
**BDD user story nümunəsi:**
```
Story: Get a post
In order to display a post to the user
As a calling program
I want to get a post
Scenario 1: using an id
  Given a post id 1
  When I send a GET request with the id
  Then I should get a post
Scenario 2: using a non-integer id
  Given a post id "hello"
  When I send a GET request with the id
  Then I should get a HTTP 500 response
```

**Quraşdırma:**
```bash
go get github.com/onsi/ginkgo/ginkgo   # CLI daxil
go get github.com/onsi/gomega          # matcher kitabxanası
```

**Mövcud testləri çevir:** `ginkgo convert .` — suite faylı + testləri yenidən yazır

**Sıfırdan:** `ginkgo bootstrap` (suite) + `ginkgo generate` (skelet)

**Kitabdan kod nümunəsi (Gomega matcher-ləri ilə):**
```go
package main_test

import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
    . "gwp/Chapter_8_Testing_Web_Applications/test_ginkgo"
)

var _ = Describe("Get a post", func() {
    var mux *http.ServeMux
    var post *FakePost
    var writer *httptest.ResponseRecorder

    BeforeEach(func() {                 // fixture: hər scenaridən əvvəl
        post = &FakePost{}
        mux = http.NewServeMux()
        mux.HandleFunc("/post/", HandleRequest(post))   // EXPORTED olmalı!
        writer = httptest.NewRecorder()
    })

    Context("Get a post using an id", func() {
        It("should get a post", func() {
            request, _ := http.NewRequest("GET", "/post/1", nil)
            mux.ServeHTTP(writer, request)
            Expect(writer.Code).To(Equal(200))
            var post Post
            json.Unmarshal(writer.Body.Bytes(), &post)
            Expect(post.Id).To(Equal(1))
        })
    })

    Context("Get an error if post id is not an integer", func() {
        It("should get a HTTP 500 response", func() {
            request, _ := http.NewRequest("GET", "/post/hello", nil)
            mux.ServeHTTP(writer, request)
            Expect(writer.Code).To(Equal(500))
        })
    })
})
```

**Sub-kod izahı:**
- `package main_test` → main paketindən **izolyasiya** → HandleRequest **exported** olmalı (böyük hərf)
- `var _ = Describe(...)` → init funksiyası olmadan çağırış (circuit trick)
- `Describe` (story) → `Context` (scenario) → `It` (behavior) iyerarxiyası
- `BeforeEach` → fixture
- `Expect(x).To(Equal(y))` → Gomega matcher assertion-u
- İşə salma: `ginkgo -v` → formatlaşdırılmış output

## Test alətləri müqayisəsi
| Alət | Funksiya | Assertion | Fixture |
|---|---|---|---|
| testing (std) | `TestXxx(t *testing.T)` | `t.Error/Fatal` | TestMain |
| httptest | Recorder + NewRequest | — | setUp funksiyaları |
| gocheck | Suite + TestXxx(c *C) | `c.Check/Assert` | SetUpTest/SetUpSuite |
| Ginkgo | Describe/Context/It | `Expect().To()` (Gomega) | BeforeEach/AfterEach |

## Əsas terminlər
- testing.T / testing.B / testing.M
- Fail / FailNow / Error / Fatal / Skip
- Coverage (`-cover`)
- Short mode (`-short` + `testing.Short()`)
- Parallel testing (`t.Parallel()`, `-parallel N`)
- Benchmark (`-bench`, `b.N`, ns/op)
- ResponseRecorder (cavab qeydedicisi)
- http.NewRequest (test sorğusu)
- TestMain (mərkəzi lifecycle)
- Test Double (test ikiləsi)
- Dependency Injection (asılılığın inyeksiyası)
- Interface injection (Text)
- Test Fixture (sınaq şəraiti)
- SetUpTest / TearDownTest / SetUpSuite / TearDownSuite
- Check / Assert (gocheck)
- BDD (Behavior-Driven Development)
- User Story / Scenario
- Describe / Context / It / BeforeEach
- Matcher (Gomega) / Expect

## Praktik nəticə
- Kodu test-olunacaq şəkildə dizayn et: funksiyaları kiçik, ayrı və parametrli yaz.
- Handler testi üçün 6 addım: mux → HandleFunc → Recorder → NewRequest → ServeHTTP → yoxla.
- `-short` + `testing.Short()` ilə uzun testləri CI-də qısaltma.
- DB asılılığını aradan qaldır: interfeys (Text) + sahə (Db) + closure (handleRequest(t Text) http.HandlerFunc) + test double (FakePost).
- gocheck: suite + fixture + zəngin assertion; Ginkgo: BDD story → Describe/Context/It strukturuna.
- Benchmark-lə müqayisə et (Decode vs Unmarshal = 25% fərq) — "hissələr üstünlük" iddialarını ölç.

## Mənbə
Pages: 211-242 (PDF), book pages 190-221
