# Chapter 19 — Testing (səh. 608-627)

## Bu fəsil nədən bəhs edir?

Test növləri: unit (table-driven + subtest + assert), integration
(sqlmock ilə DB), E2E, HTTP testing (httptest.NewServer), fuzz testing
(f.Fuzz), benchmark (b.N), test suite (TestMain + setup/teardown), JSON
test report və code coverage (-cover).

## Əsas fikirlər

### 1. Unit testlər + table-driven pattern
**Qaydalar:** _test.go suffiksi; Test prefiks; testing.T parametri;
main() YOXDUR — hər test müstəqil ardıcıl icra olunur.

**Table-driven + assert (gotest.tools):**
```go
package main

import (
    "testing"
    "gotest.tools/assert"
)

func add(x, y int) int {
    return x + y
}

func TestAdd(t *testing.T) {
    tests := []struct {
        name   string
        inputs []int
        want   int
    }{
        {"Test Case 1", []int{5, 6}, 11},
        {"Test Case 2", []int{11, 7}, 18},
        {"Test Case 3", []int{1, 8}, 9},
        {"Test Case 4 (intentional failure)", []int{2, 3}, 0},  // SƏHRƏN 5
    }

    for _, test := range tests {
        got := add(test.inputs[0], test.inputs[1])
        assert.Equal(t, test.want, got)
    }
}
// Çıxış: assertion failed: 0 (test.want int) != 5 (got int) → FAIL
```

**Subtests (t.Run) — hansı case-in düşdüyü görünür:**
```go
for _, test := range tests {
    test := test                     // LOKAL KOPYA — closure tələbi!
    t.Run(test.name, func(t *testing.T) {
        got := add(test.inputs[0], test.inputs[1])
        assert.Equal(t, test.want, got)
    })
}
// --- FAIL: TestAdd/Test_Case_4_(intentional_failure)
//     → dəqiq case adı göstərilir!
```
- `test := test` — loop dəyişəni closure-da REFERENCE ilə tutulur;
  kopyasız linter xətası verir / gözlənilməz davranış

**Subtest faydaları:** izolyasiya, aydın çıxış, paralel icra imkanı,
struktur təşkilat. Müsbət + mənfi + edge case-lər.

### 2. Integration testlər (sqlmock)
**Nədir:** komponentlərin ƏLAQƏSİNİ yoxlayır; unit-dən çox setup
tələb edir. Mock/stub alətləri: gomock, testify/mock, go-sqlmock.

```go
import (
    "context"
    "database/sql"
    "testing"
    "github.com/DATA-DOG/go-sqlmock"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

type Record struct {
    ID    int
    Name  string
    Value string
}

type Database struct {
    conn *sql.DB
}

func NewDatabase(conn *sql.DB) *Database {
    return &Database{conn: conn}
}

func (d *Database) InsertRecord(ctx context.Context, record Record) error {
    _, err := d.conn.ExecContext(ctx,
        "INSERT INTO records (id, name, value) VALUES ($1, $2, $3)",
        record.ID, record.Name, record.Value)
    return err
}

func (d *Database) GetRecordByID(ctx context.Context, id int) (Record, error) {
    var record Record
    row := d.conn.QueryRowContext(ctx,
        "SELECT id, name, value FROM records WHERE id = $1", id)
    err := row.Scan(&record.ID, &record.Name, &record.Value)
    return record, err
}

func TestDatabaseIntegration(t *testing.T) {
    db, mock, err := sqlmock.New()          // in-memory SQL mock
    require.NoError(t, err)
    defer db.Close()

    testRecord := Record{ID: 1, Name: "TestRecord", Value: "TestValue"}

    // EXPECTATIONS — hansı sorğular gələcək + nə qaytaracaq:
    mock.ExpectExec("INSERT INTO records").
        WithArgs(testRecord.ID, testRecord.Name, testRecord.Value).
        WillReturnResult(sqlmock.NewResult(1, 1))

    rows := sqlmock.NewRows([]string{"id", "name", "value"}).
        AddRow(testRecord.ID, testRecord.Name, testRecord.Value)
    mock.ExpectQuery("SELECT id, name, value FROM records").
        WillReturnRows(rows)

    // TEST:
    dbInstance := NewDatabase(db)
    err = dbInstance.InsertRecord(context.Background(), testRecord)
    assert.NoError(t, err, "Error inserting record into the database")

    retrievedRecord, err := dbInstance.GetRecordByID(context.Background(), 1)
    assert.NoError(t, err)
    assert.Equal(t, testRecord, retrievedRecord)

    assert.NoError(t, mock.ExpectationsWereMet())   // hamısı çağırıldı?
}
```
- sqlmock: real DB lazım deyil; gözləntilər səniddə səhv → test düşür

### 3. E2E testlər
**Nədir:** BÜTÜN sistemi istifadəçi ssenariləri kimi yoxlayır (UI+API+DB).
- Realistik ssenarilər; çoxkomponent əlaqəsi; production-a yaxın mühit
- Alətlər: Selenium, Cypress, xüsusi HTTP klientlər, testing + httptest
- Best practices: izolyasiya, avtomatlaşdırma (CI), aydın ssenarilər,
  data idarəetməsi (test sabitliyi)

### 4. HTTP testing (httptest.NewServer)
```go
func TestAuthenticationIntegration(t *testing.T) {
    // TEST SERVER — auth servisinin simulyasiyası:
    authService := httptest.NewServer(http.HandlerFunc(
        func(w http.ResponseWriter, r *http.Request) {
            if r.Header.Get("Authorization") == "Bearer valid_token" {
                w.WriteHeader(http.StatusOK)
                w.Write([]byte(`{"user_id": "123", "username": "testuser"}`))
            } else {
                w.WriteHeader(http.StatusUnauthorized)
            }
        }))
    defer authService.Close()

    app := NewApplication(authService.URL)    // URL avtomatik verilir!
    token := "valid_token"
    gotUser, err := app.AuthenticateUser(token)
    assert.NoError(t, err)
    assert.Equal(t, "123", gotUser.UserID)
    assert.Equal(t, "testuser", gotUser.Username)
}
```
- HTTP testing: functional validation + integration verification +
  error handling + security (auth yoxlamaları)

### 5. Fuzz testing
```go
func FuzzAdd(f *testing.F) {
    f.Fuzz(func(t *testing.T, i int, j int) {
        got := add(i, j)
        assert.Equal(t, i+j, got)
    })
}
```
- Fuzz prefiks + testing.F; f.Fuzz — random input generator
- `go test -fuzz .` — fərdi input-lardan MİLYONLARLA test:
```
fuzz: elapsed: 20s, execs: 2237555 (107504/sec), new interesting: 0
```
- Uğursuz nümunə seed corpus-a yazılır → səhv düzəltmək üçün

### 6. Benchmark
```go
func BenchmarkAdd(b *testing.B) {
    for i := 0; i < b.N; i++ {
        add(1, 2)
    }
}
```
```bash
go test -bench .
# BenchmarkAdd-10   1000000000   0.3444 ns/op
```
- Benchmark prefiks + testing.B; runner b.N-i avtomatik artırır
- ns/op — əməliyyat başına nano saniyə; real input istifadə et;
  böyük funksiyaları parçalara böl

### 7. Test suite (TestMain)
```go
func setup() {
    log.Println("setup() running")
}
func teardown() {
    log.Println("teardown() running")
}

func TestMain(m *testing.M) {
    setup()
    defer teardown()     // m.Run qayıdandan SONRA icra
    m.Run()              // bütün testləri işə sal
}

func TestA(t *testing.T) { log.Println("TestA running") }
func TestB(t *testing.T) { log.Println("TestB running") }
func TestC(t *testing.T) { log.Println("TestC running") }
```
- TestMain — qlobal setup/teardown; shared resurslar üçün

### 8. Test report + coverage
```bash
# JSON report:
go test . -v -json > test-report.json

# Coverage:
go test . -cover
# ok ... coverage: 100.0% of statements

# Coverage fayla:
go test . -coverprofile=coverage.out
```
- assert mesajlarını zənginləşdir: `assert.Equal(t, i+j, got, "should sum properly")`
- Sənaye standartı: ~80% coverage hədəfi; CI-də PR-larda drop
  yoxlaması
- Yeni funksiya əlavə → coverage düşür → test əlavə et → geri qalxır

## Əsas terminlər
- _test.go / Test prefiks / testing.T
- Table-driven test — anonim struct cədvəli
- Subtest (t.Run) + loop variable copy (test := test)
- assert.Equal / require.NoError (gotest.tools / testify)
- Positive/negative/edge case
- Integration test — komponent əlaqələri
- sqlmock — ExpectExec/ExpectQuery/ExpectationsWereMet
- E2E — bütün sistem, istifadəçi axını
- httptest.NewServer — test HTTP serveri
- Fuzz — Fuzz prefiks, f.Fuzz, random input
- Benchmark — testing.B, b.N, ns/op
- TestMain — qlobal setup/teardown
- -json — machine-readable report
- -cover / coverprofile — örtük faizi (80% standart)

## Praktik nəticə
Unit testləri table-driven + subtest formatında yaz (name/inputs/want
struct + t.Run + test := test kopyası!) — pozulan case dərhal görünür.
DB asılılığını sqlmock ilə (gözləntilər → icra → ExpectationsWereMet),
HTTP asılılığını httptest.NewServer ilə yoxla. Performans üçün
-bench, anomaliyalar üçün -fuzz. Qlobal fixture-lər TestMain-da.
Report üçün -json, keyfiyyət göstəricisi üçün -cover (~80% hədəf;
hər yeni funksiya yeni test tələb edir). Assert mesajlarını mümkün
qədər izahlı yaz — reportda səhvin SƏBƏBİ oxunsun.

## Mənbə
Pages: 608-627 (PDF 608-627)
