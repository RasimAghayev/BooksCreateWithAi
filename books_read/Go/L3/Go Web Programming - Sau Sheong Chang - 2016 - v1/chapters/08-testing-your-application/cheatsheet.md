# Chapter 8 — Testing your application Cheatsheet

## `go test` komandası və flag-ləri

**Nə edir:** `_test.go` fayllarındakı `TestXxx` funksiyalarını icra edir.

**Parametrlər:**
- `-v` — verbose çıxış
- `-cover` — kod örtüyü (coverage)
- `-short` — `-short` rejimi aktivləşdirir (`testing.Short()` true olur)
- `-bench .` — benchmark-ləri regex ilə seçir
- `-run x` — yalnız `x` adlı test-ləri icra edir
- `-parallel N` — neçə paralel test icra oluna bilər

```bash
go test -v -cover -short
go test -bench .
```

**Mənbə:** Chapter 8, page 194

---

## `t.Error()` / `t.Errorf()` / `t.Fatal()` / `t.Fatalf()`

**Nə edir:** Test xətası bildirir. `Error` davam edir, `Fatal` dayanır.

```go
if post.Id != 1 {
    t.Error("Wrong id, was expecting 1 but got", post.Id)
}
if writer.Code != 200 {
    t.Errorf("Response code is %v", writer.Code)
}
```

**Mənbə:** Chapter 8, page 194

---

## `t.Skip("mesaj")` və `testing.Short()`

**Nə edir:** Test-i keçir. `-short` flag ilə birlikdə uzun testləri skip etmək üçün.

```go
func TestEncode(t *testing.T) {
    t.Skip("Skipping encoding for now")
}

func TestLongRunning(t *testing.T) {
    if testing.Short() {
        t.Skip("Skip in short mode")
    }
    time.Sleep(10 * time.Second)
}
```

**Mənbə:** Chapter 8, page 195

---

## `t.Parallel()` ilə paralel test

**Nə edir:** Test-i paralel qrupa əlavə edir. `-parallel N` ilə maksimum say təyin olunur.

```go
func TestParallel_1(t *testing.T) {
    t.Parallel()
    time.Sleep(1 * time.Second)
}
```

**Mənbə:** Chapter 8, page 196

---

## `BenchmarkXxx(b *testing.B)`

**Nə edir:** Performans ölçür. `b.N` iterasiyaları Go təyin edir.

```go
func BenchmarkDecode(b *testing.B) {
    for i := 0; i < b.N; i++ {
        decode("post.json")
    }
}
// go test -bench .
```

**Mənbə:** Chapter 8, page 198

---

## `TestMain(m *testing.M)` — suite setup/teardown

**Nə edir:** Bütün test-lər üçün bir dəfə setup və teardown işlədir.

```go
func TestMain(m *testing.M) {
    setUp()
    code := m.Run()
    os.Exit(code)
}

func setUp() {
    mux = http.NewServeMux()
    mux.HandleFunc("/post/", handleRequest)
    writer = httptest.NewRecorder()
}
```

**Mənbə:** Chapter 8, page 203

---

## `httptest.NewRecorder()` + `http.NewRequest()` + `mux.ServeHTTP()`

**Nə edir:** Real server qaldırmadan handler test edir. `ResponseRecorder` cavabı yadda saxlayır.

```go
writer := httptest.NewRecorder()
request, _ := http.NewRequest("GET", "/post/1", nil)
mux.ServeHTTP(writer, request)
if writer.Code != 200 { /* ... */ }
var post Post
json.Unmarshal(writer.Body.Bytes(), &post)
```

**Mənbə:** Chapter 8, page 202

---

## Test double: interface + FakePost

**Nə edir:** Handler-i verilənlər bazasından asılı etmədən test etmək.

```go
type Text interface {
    fetch(id int) (err error)
    create() (err error)
    update() (err error)
    delete() (err error)
}

type FakePost struct { Id, Content, Author string /* ... */ }
func (post *FakePost) fetch(id int) error { post.Id = id; return nil }

mux.HandleFunc("/post/", handleRequest(&FakePost{}))
```

**Mənbə:** Chapter 8, page 209

---

## `gocheck` — suite, fixture, assertion

**Nə edir:** Standart testing-i suite/fixture/assertion ilə genişləndirir.

```go
import . "gopkg.in/check.v1"

type PostTestSuite struct{}
func init() { Suite(&PostTestSuite{}) }
func Test(t *testing.T) { TestingT(t) }

func (s *PostTestSuite) SetUpTest(c *C)  { /* hər test-dən əvvəl */ }
func (s *PostTestSuite) TearDownTest(c *C) { /* hər test-dən sonra */ }
func (s *PostTestSuite) SetUpSuite(c *C) { /* bir dəfə */ }
func (s *PostTestSuite) TearDownSuite(c *C) { /* bir dəfə */ }

func (s *PostTestSuite) TestGetPost(c *C) {
    c.Check(s.writer.Code, Equals, 200)
    c.Check(post.Id, Equals, 1)  // davam edir
    c.Assert(post.Id, Equals, 1) // dayanır
}
```

**Mənbə:** Chapter 8, page 211

---

## Ginkgo BDD strukturu

**Nə edir:** User story-ləri `Describe`/`Context`/`It` blokları ilə test-ə çevirir. Gomega ilə matcher-lər.

```go
import (
    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
    "testing"
)
func TestGinkgo(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, "Ginkgo Suite")
}

var _ = Describe("Get a post", func() {
    var post *FakePost
    BeforeEach(func() { post = &FakePost{} })

    Context("Get a post using an id", func() {
        It("should get a post", func() {
            Expect(post.Id).To(Equal(1))
        })
    })
})
```

**Mənbə:** Chapter 8, page 220

---

## `Expect(val).To(Equal(expected))` — Gomega matcher

**Nə edir:** Oxunaqlı assertion sintaksisi. `Equal`, `BeNil`, `HaveLen` və s. matcher-lər var.

```go
Expect(writer.Code).To(Equal(200))
Expect(post.Id).To(Equal(1))
Expect(err).To(BeNil())
```

**Mənbə:** Chapter 8, page 220
