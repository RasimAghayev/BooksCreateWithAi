# Chapter 9 — Testing and Debugging (Testləmə və Debug)

## Bu chapter nədən bəhs edir?
Go testing paketinə (T/B tipləri, Errorf/Fatalf/Skipf), table-driven testlərə, subtest/parallel/helper-lərə, mock asılılıqlarına (interface əsaslı), log/tracing-ə (logrus), performans profilinqinə (pprof: CPU/block/memory) və 10 klassik veb xətasına + həllərinə.

## Əsas fikirlər

### 1. Testing nəyə lazımdır
Testing = proaktiv keyfiyyət qoruması (unit/integration/E2E); Debugging = reaktiv xəta ovu (Delve debugger). Performance testləri: response time, load, resilience.

### 2. Go testing paketi — əsaslar
**Sadə test:**
```go
func TestSum(t *testing.T) {
    result := sum(2, 3)
    if result != 5 {
        t.Errorf("Expected 5 but got %d", result)
    }
}
```
**Axın idarəsi:**
| Metod | Davranış |
|-------|----------|
| `t.Errorf` | xəta qeyd et, test DAVAM etsin |
| `t.Fatalf` | KRİTİK — test dərhal dayansın |
| `t.Skipf` | bu testi ötür (mühit şərtinə görə) |

**Testify assertion nümunəsi:**
```go
func TestFetchBookDetails(t *testing.T) {
    book, err := FetchBookDetails("123456789")
    assert.Nil(t, err)
    assert.Equal(t, "Go Programming", book.Title)
}
```
**Benchmark:**
```go
func BenchmarkSum(b *testing.B) {
    for i := 0; i < b.N; i++ { sum(2, 3) }
}
```
**İcra:** `go test` — test faylları kodun yanında (`book.go` ↔ `book_test.go`).

### 3. Table-driven testing — Go-nun imza pattern-i
```go
func TestAdd(t *testing.T) {
    tests := []struct {
        name string
        a, b, want int
    }{
        {"add positives", 3, 4, 7},
        {"add negatives", -3, -4, -7},
        {"add mixed", -3, 4, 1},
        {"add zero", 0, 4, 4},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {      // subtest — adlı, müstəqil
            got := Add(tt.a, tt.b)
            if got != tt.want {
                t.Errorf("Add(%d, %d) = %d; want %d", tt.a, tt.b, got, tt.want)
            }
        })
    }
}
```
**Sub-kod izahı:** Anonim struct slice-u = ssenari cədvəli; `t.Run` hər case-i adlı subtest kimi icra edir — hansı case düşdü dərhal görünür.

### 4. Helper funksiyalar + paralel + subtest
```go
func setupMockDatabase() *MockDB { /* ... */ return db }
func tearDownMockDatabase(db *MockDB) { /* cleanup */ }

func TestSomeFeature(t *testing.T) {
    t.Parallel()                 // paralel icra (asılılıq YOXDURSA!)
    // ...
}

func TestBookFeatures(t *testing.T) {
    db := setupMockDatabase()           // shared setup
    defer tearDownMockDatabase(db)
    t.Run("test fetch", func(t *testing.T) { /* ... */ })
    t.Run("test update", func(t *testing.T) { /* ... */ })
}
```
**Xəbərdarlıq:** Paralel testlər paylaşılan state olsa qeyri-müəyyən düşür — izolyasiya mütləqdir.

### 5. Mock asılılıqlar — interface əsaslı mocking
**Nədir:** Real DB/API əvəzinə simulyasiya — unit testi xarici asılılıqlardan izolyasiya edir (yoxsa bu integration test olar).

**Kitabdan kod nümunəsi:**
```go
type Database interface {
    FetchBookByID(id int) (Book, error)
}

type MockDatabase struct{ books []Book }

func (m *MockDatabase) FetchBookByID(id int) (Book, error) {
    for _, book := range m.books {
        if book.ID == id { return book, nil }
    }
    return Book{}, fmt.Errorf("Book not found")
}

// Testdə istifadə:
func TestFetchBookDetails(t *testing.T) {
    mockDB := &MockDatabase{
        books: []Book{{ID: 1, Name: "Go Basics", Author: "A. Developer"}},
    }
    book, err := FetchBookDetails(mockDB, 1)   // interface parametrinə mock ötür
    if err != nil || book.Name != "Go Basics" { t.Fail() }
}
```
**Xarici API mock-u:** `ReviewAPI` interface + `MockReviewAPI` — predefined data qaytarır.
**Qeyd:** Mock hər nüansı tutmur — integration testlər də lazımdır.

### 6. Logging və tracing
**Standart log:**
```go
log.Printf("Error fetching book with ID %d: %v", bookID, err)
```
**Structured logging (logrus):**
```go
import log "github.com/sirupsen/logrus"

func init() {
    log.SetFormatter(&log.JSONFormatter{})   // JSON — ELK/Splunk-uyğun
    log.SetLevel(log.InfoLevel)
}

log.WithFields(log.Fields{
    "bookID": bookID,
    "error":  err,
}).Error("Failed to fetch book details")
```
**Praktikalar:** Request-ID middleware (context ilə) — hər loqda sorğu izi; OpenTracing — span-larla sorğu yolu; Logstash/Fluentd — mərkəzləşdirilmiş toplus.

### 7. Performans profilinqi (pprof)
**CPU:**
```go
import _ "net/http/pprof"      // və ya manual:
f, _ := os.Create("cpu.pprof")
pprof.StartCPUProfile(f)
defer pprof.StopCPUProfile()
// go tool pprof cpu.pprof → top
```
**Block (sinxronizasiya):**
```go
runtime.SetBlockProfileRate(1)
// sonra:
pprof.Lookup("block").WriteTo(f, 0)
```
**Memory:**
```go
f, _ := os.Create("memory.pprof")
pprof.WriteHeapProfile(f)
```
**Təhlil:** `go tool pprof` + `top` — ən çox CPU/alloc edən funksiyaları tap.

### 8. 10 klassik xəta + həlləri
| Xəta | Simptom | Həll |
|------|---------|------|
| **N+1 query** | list + hər element üçün ayrı sorğu | JOIN və ya batch fetch |
| **Data race** | eynizamanı oxu/yazı | sync.Mutex / channel |
| **Memory leak** | qlobal dəyişənlər, açılmamış kanallar, zombí goroutine | pprof ilə tap; kanalları bağla |
| **Broken auth** | validasiyasız input, güvensiz token | etibarlı kitabxanalar + dəbli validasiya |
| **İneffektiv data** | pagination-sız böyük dataset | LIMIT/OFFSET + DB səviyyəsində filter |
| **Input validation** | SQL injection, XSS | hər input yoxla; raw input-ı sorğuya/səhifəyə qoyma |
| **Səhv error handling** | app çöker / stack trace istifadəçiyə görünür | generic mesaj user-ə, detal log-da |
| **Dependency xətası** | köhnə/ləğv kitabxana | `go get -u` + vulnerability scanner |
| **Hardcoded config** | koddа DB parol/API key | env var: `os.Getenv("DB_PASSWORD")` |
| **Unsafe concurrency** | kilidsiz paylaşılan sayğac | mutex/channel — hər paylaşılan dəyişən |

**Nümunə həll kodları:**
```go
// N+1 → JOIN
booksWithAuthors := db.GetBooksWithAuthors()

// Input validation
if len(bookTitle) < 3 { return errors.New("book title is too short") }

// Error handling
if err := db.SaveBook(book); err != nil {
    log.Println("Error saving book:", err)      // detal logda
    return errors.New("failed to save the book") // user-ə generic
}
```

## Əsas terminlər
- testing.T / testing.B — test/benchmark tipləri
- Table-driven test — cədvəl əsaslı ssenari testi
- t.Run subtest — adlı alt-test
- t.Parallel — paralel test icrası
- Mock — interface əsaslı simulyasiya
- logrus/zap — structured logging
- OpenTracing — span əsaslı izləmə
- pprof — CPU/block/memory profiler
- N+1 query — performans antipattern-i
- Data race — eynizamanı yazı zərəri

## Praktik nəticə
1. Hər funksiya üçün `Xxx_test.go` yanında; table-driven + t.Run adlı subtest-lərlə ssenari əhatəsi.
2. Asılılıqları HƏMİŞƏ interface-ə sar — mock test-də, real integration-da.
3. Fatal xətalar Fatalf, qeyri-kritiklər Errorf; mühit asılı testlər Skipf.
4. Structured JSON logging + request-ID — producda axtarış 5 saniyə, 5 saat YOX.
5. pprof üçlüyünü (CPU/block/memory) optimizasiya başlanğıcında işlət — fərziyyə yox, ölçü.
6. Xəta mesajları: user-ə generic, log-a detallı; heç vaxt stack trace göstərmə.

## Mənbə
Pages: 231-261 (PDF səh. 231-261)
