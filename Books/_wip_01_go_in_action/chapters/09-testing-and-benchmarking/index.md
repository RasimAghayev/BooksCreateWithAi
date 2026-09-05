# Fəsil 9 — Testetmə və benchmarkinq

## Bu fəsil nədən bəhs edir?

Go-nun test framework-u: unit testlər (basic və table test), xarici resurslardan asılılığı aradan qaldıran mockinq (`httptest`), web servis endpoint-lərinin serveri işə salmadan testi, sənədləşmə+test olan Example funksiyaları və benchmarkinq (performans ölçmə).

## Əsas fikirlər

### 1. Unit testin əsasları
**Nədir:** Kodun müəyyən hissəsinin gözlənilən ssenari üzrə düzgün işlədiyini yoxlayan funksiya. Positive-path (normal icra xətasız) və negative-path (düzgün xətanın qaytarılması) testləri.

**Konvensiyalar (MƏCBURİ):**
- Test faylı adı **`_test.go`** ilə bitməlidir — yoxsa `go test` onu görmür
- Test funksiyası **exported** olmalı, **`Test`** prefiksi ilə başlamalıdır
- İmza: `func TestXxx(t *testing.T)` — heç nə qaytarmır

**Kitabdan kod nümunəsi (basic unit test):**
```go
package listing01

import (
    "net/http"
    "testing"
)

const checkMark = "\u2713"
const ballotX = "\u2717"

// TestDownload validates the http Get function can download content.
func TestDownload(t *testing.T) {
    url := "http://www.goinggo.net/feeds/posts/default?alt=rss"
    statusCode := 200

    t.Log("Given the need to test downloading content.")
    {
        t.Logf("\tWhen checking \"%s\" for status code \"%d\"",
            url, statusCode)
        {
            resp, err := http.Get(url)
            if err != nil {
                t.Fatal("\t\tShould be able to make the Get call.",
                    ballotX, err)
            }
            t.Log("\t\tShould be able to make the Get call.",
                checkMark)

            defer resp.Body.Close()

            if resp.StatusCode == statusCode {
                t.Logf("\t\tShould receive a \"%d\" status. %v",
                    statusCode, checkMark)
            } else {
                t.Errorf("\t\tShould receive a \"%d\" status. %v %v",
                    statusCode, ballotX, resp.StatusCode)
            }
        }
    }
}
```

**Sub-kod izahı:**
- `t.Log / t.Logf` → informasiya yazısı (yalnız `-v` ilə görünür və ya test fail olsa)
- `t.Fatal / t.Fatalf` → testi **FAIL** edir + mesaj yazır + **funksiyanı dayandırır** (digər test funksiyaları davam edir)
- `t.Error / t.Errorf` → testi FAIL edir, amma funksiya **davam edir**
- `t.Fatal`/`t.Error` çağırılmayan test **PASS** sayılır
- Müəllifin output konvensiyası: "Given the need… / When checking… / Should…" — test çıxışı sənədləşmə kimi oxunmalıdır

**İşə salma:**
```bash
go test -v    # verbose — bütün output
```

### 2. Table testlər
**Nədir:** Eyni test kodunu bir çox fərqli parametr/nəticə cütlüyü ilə işlədən pattern — dəyərlər "cədvəldə" saxlanılır.

**Kitabdan kod nümunəsi:**
```go
func TestDownload(t *testing.T) {
    var urls = []struct {
        url        string
        statusCode int
    }{
        {
            "http://www.goinggo.net/feeds/posts/default?alt=rss",
            http.StatusOK,
        },
        {
            "http://rss.cnn.com/rss/cnn_topstbadurl.rss",
            http.StatusNotFound,
        },
    }

    t.Log("Given the need to test downloading different content.")
    {
        for _, u := range urls {
            t.Logf("\tWhen checking \"%s\" for status code \"%d\"",
                u.url, u.statusCode)
            {
                resp, err := http.Get(u.url)
                if err != nil {
                    t.Fatal("\t\tShould be able to Get the url.",
                        ballotX, err)
                }
                t.Log("\t\tShould be able to Get the url",
                    checkMark)

                defer resp.Body.Close()

                if resp.StatusCode == u.statusCode {
                    t.Logf("\t\tShould have a \"%d\" status. %v",
                        u.statusCode, checkMark)
                } else {
                    t.Errorf("\t\tShould have a \"%d\" status %v %v",
                        u.statusCode, ballotX, resp.StatusCode)
                }
            }
        }
    }
}
```

**Sub-kod izahı:**
- `[]struct{ url string; statusCode int }{...}` → anonim struct slice — cədvəl: giriş + gözlənilən nəticə
- `for _, u := range urls` → hər sətir üçün eyni test kodu icra olunur
- Yeni ssenari əlavə etmək = cədvələ yeni sətir yazmaq; test kodunun özü dəyişmir
- `http.StatusOK` / `http.StatusNotFound` → status kod konstantları (maqik rəqəmlərdən qaçın)

### 3. Mocking — httptest paketi
**Problem:** İnternetə çıxışı olmayan CI mühitində, yaxud sahibi olmadığın serverdən asılı testlər işə yaramır → deployment bloklanır.

**Həll:** `httptest.NewServer` — lokal mock server.

**Kitabdan kod nümunəsi:**
```go
// feed is mocking the XML document we except to receive.
var feed = `<?xml version="1.0" encoding="UTF-8"?>
<rss>
  <channel>
    <title>Going Go Programming</title>
    ...
  </channel>
</rss>`

// mockServer returns a pointer to a server to handle the get call.
func mockServer() *httptest.Server {
    f := func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(200)
        w.Header().Set("Content-Type", "application/xml")
        fmt.Fprintln(w, feed)
    }

    return httptest.NewServer(http.HandlerFunc(f))
}

func TestDownload(t *testing.T) {
    statusCode := http.StatusOK

    server := mockServer()
    defer server.Close()

    t.Log("Given the need to test downloading content.")
    {
        t.Logf("\tWhen checking \"%s\" for status code \"%d\"",
            server.URL, statusCode)
        {
            resp, err := http.Get(server.URL)
            if err != nil {
                t.Fatal("\t\tShould be able to make the Get call.",
                    ballotX, err)
            }
            t.Log("\t\tShould be able to make the Get call.",
                checkMark)

            defer resp.Body.Close()

            if resp.StatusCode != statusCode {
                t.Fatalf("\t\tShould receive a \"%d\" status. %v %v",
                    statusCode, ballotX, resp.StatusCode)
            }
            t.Logf("\t\tShould receive a \"%d\" status. %v",
                statusCode, checkMark)
        }
    }
}
```

**Sub-kod izahı:**
- `http.HandlerFunc(f)` → adapter: adi funksiyanı (`func(ResponseWriter, *Request)`) HTTP handler-ə çevirir
- `httptest.NewServer(handler)` → localhost-da təsadüfi portda (hər run-da dəyişir) real HTTP server qaldırır
- Handler daxilində: status kod set et → Content-Type set et → body yaz (`fmt.Fprintln(w, feed)`)
- `server.URL` → mock URL; `http.Get(server.URL)` onu real internet sorğusu kimi görür
- `defer server.Close()` → test bitəndə server söndürülür
- Nəticə: internet olmadan test keçir; goinggo.net-ə heç bir sorğu getmir

### 4. Endpoint testi — serveri işə salmadan
**Nədir:** Web API handler-larını `http.ListenAndServe` işə salmadan test etmək.

**Test olunan web servis (handlers paketi):**
```go
// Package handlers provides the endpoints for the web service.
package handlers

import (
    "encoding/json"
    "net/http"
)

// Routes sets the routes for the web service.
func Routes() {
    http.HandleFunc("/sendjson", SendJSON)
}

// SendJSON returns a simple JSON document.
func SendJSON(rw http.ResponseWriter, r *http.Request) {
    u := struct {
        Name  string
        Email string
    }{
        Name:  "Bill",
        Email: "bill@ardanstudios.com",
    }

    rw.Header().Set("Content-Type", "application/json")
    rw.WriteHeader(200)
    json.NewEncoder(rw).Encode(&u)
}
```

**Test faylı (handlers_test paketi!):**
```go
// Sample test to show how to test the execution of an
// internal endpoint.
package handlers_test

import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/goinaction/code/chapter9/listing17/handlers"
)

func init() {
    handlers.Routes()
}

// TestSendJSON testing the sendjson internal endpoint.
func TestSendJSON(t *testing.T) {
    t.Log("Given the need to test the SendJSON endpoint.")
    {
        req, err := http.NewRequest("GET", "/sendjson", nil)
        if err != nil {
            t.Fatal("\tShould be able to create a request.",
                ballotX, err)
        }
        t.Log("\tShould be able to create a request.",
            checkMark)

        rw := httptest.NewRecorder()
        http.DefaultServeMux.ServeHTTP(rw, req)

        if rw.Code != 200 {
            t.Fatal("\tShould receive \"200\"", ballotX, rw.Code)
        }
        t.Log("\tShould receive \"200\"", checkMark)

        u := struct {
            Name  string
            Email string
        }{}

        if err := json.NewDecoder(rw.Body).Decode(&u); err != nil {
            t.Fatal("\tShould decode the response.", ballotX)
        }
        t.Log("\tShould decode the response.", checkMark)

        if u.Name == "Bill" {
            t.Log("\tShould have a Name.", checkMark)
        } else {
            t.Error("\tShould have a Name.", ballotX, u.Name)
        }

        if u.Email == "bill@ardanstudios.com" {
            t.Log("\tShould have an Email.", checkMark)
        } else {
            t.Error("\tShould have an Email.", ballotX, u.Email)
        }
    }
}
```

**Sub-kod izahı:**
- `package handlers_test` → **`_test` paketi** — eyni qovluqda olsa belə yalnız **exported** identifikatorlara çıxış verir (black-box test)
- `init()` → `handlers.Routes()` — routingsiz test `StatusNotFound` alar
- `http.NewRequest("GET", "/sendjson", nil)` → sorğu obyekti yarat (nil body — GET-dir)
- `httptest.NewRecorder()` → `ResponseRecorder` — **cavabı yadda saxlayan** saxta ResponseWriter (real network yoxdur!)
- `http.DefaultServeMux.ServeHTTP(rw, req)` → mux-u birbaşa çağır — sorğu handler-a gedir, server dayanır
- `rw.Code`, `rw.Body` → yazılmış status və body-ni yoxla
- Sonuncu 2 yoxlama `t.Error` istifadə edir — hər iki sahə həmişə yoxlanılsın deyə

### 5. Example funksiyaları — test + sənədləşmə
**Nədir:** `Example` prefiksi ilə başlayan, godoc-a düşən və **eyni zamanda test** kimi işləyən nümunə kodlar.

**Qaydalar:**
- Adı mövcud **exported** funksiya/metod adına əsaslanmalıdır: `ExampleSendJSON` → `SendJSON`
- Çıxış `// Output:` şərhi ilə göstərilir — framework stdout-u bu şərhlə müqayisə edir

**Kitabdan kod nümunəsi:**
```go
// ExampleSendJSON provides a basic example.
func ExampleSendJSON() {
    r, _ := http.NewRequest("GET", "/sendjson", nil)
    rw := httptest.NewRecorder()
    http.DefaultServeMux.ServeHTTP(rw, r)

    var u struct {
        Name  string
        Email string
    }

    if err := json.NewDecoder(w.Body).Decode(&u); err != nil {
        log.Println("ERROR:", err)
    }

    // Use fmt to write to stdout to check the output.
    fmt.Println(u)
    // Output:
    // {Bill bill@ardanstudios.com}
}
```

**Sub-kod izahı:**
- Gövdə — funksiyanın istifadə nümunəsidir (sənədləşmə rolunu oynayır)
- `// Output:` şərhi → test framework stdout-u müqayisə edir: uyğun gəlsə PASS, gəlməsə FAIL (fərq göstərilir)
- `godoc -http=":3000"` → local godoc-da paketin sənədləri yanında nümunə görünür
- `go test -run="ExampleSendJSON"` → adını `-run` regex filtri ilə işə sal

### 6. Benchmarking — performans ölçmə
**Nədir:** Kodun performansını ölçmək üçün test framework-in imkanı — müxtəlif həlləri müqayisə etmək, CPU/allokasiya problemlərini tapmaq.

**Konvensiyalar:**
- Funksiya adı **`Benchmark`** prefiksi ilə başlayır, parametr: `b *testing.B`
- Ölçülən kod **`for i := 0; i < b.N; i++`** loop-unun DAXİLİNDƏ olmalıdır
- Framework funksiyanı ən azı 1 saniyə çalışdırır, hər çağırışda `b.N` artırır

**Kitabdan kod nümunəsi (int → string, 3 üsul):**
```go
// BenchmarkSprintf provides performance numbers for the
// fmt.Sprintf function.
func BenchmarkSprintf(b *testing.B) {
    number := 10

    b.ResetTimer()

    for i := 0; i < b.N; i++ {
        fmt.Sprintf("%d", number)
    }
}

// BenchmarkFormat provides performance numbers for the
// strconv.FormatInt function.
func BenchmarkFormat(b *testing.B) {
    number := int64(10)

    b.ResetTimer()

    for i := 0; i < b.N; i++ {
        strconv.FormatInt(number, 10)
    }
}

// BenchmarkItoa provides performance numbers for the
// strconv.Itoa function.
func BenchmarkItoa(b *testing.B) {
    number := 10

    b.ResetTimer()

    for i := 0; i < b.N; i++ {
        strconv.Itoa(number)
    }
}
```

**İşə salma əmrləri:**
```bash
# Yalnız benchmark, unit testləri burax:
go test -v -run="none" -bench="BenchmarkSprintf"

# Daha dəqiq nəticə üçün 3 saniyə:
go test -run="none" -bench="BenchmarkSprintf" -benchtime=3s

# Allokasiya statistikası ilə:
go test -run="none" -bench=. -benchmem
```

**Nəticənin oxunması:**
```
BenchmarkSprintf-8   5000000    258 ns/op   16 B/op   2 allocs/op
BenchmarkFormat-8  30000000   45.9 ns/op    2 B/op   1 allocs/op
BenchmarkItoa-8    30000000   49.4 ns/op    2 B/op   1 allocs/op
```
- **iterasiya sayı** → kodun neçə dəfə icra olunduğu
- **ns/op** → əməliyyat başına nanosaniyə — əsas performans göstəricisi
- **B/op** → əməliyyat başına ayrılan byte
- **allocs/op** → əməliyyat başına heap alloksiyası sayı

**Sub-kod izahı:**
- `b.ResetTimer()` → setup vaxtını hesaba salmamaq üçün taymeri sıfırla
- `-run="none"` → heç bir unit test işə düşməsin (regex filtri)
- `-bench=.` → bütün benchmark-lar
- `-benchtime=3s` → daha uzun ölçmə (adətən 3s-dən çoxu fərq vermir)
- `-benchmem` → allokasiya statistikası əlavə olunur
- **Nəticə:** `strconv.FormatInt`/`Itoa` ~5x daha sürətlidir və 2 byte vs 16 byte alloksiya edir — `fmt.Sprintf` universal amma bahalıdır

## Test funksiyaları cədvəli
| Növ | Prefiks | Parametr | Məqsəd |
|---|---|---|---|
| Unit test | `Test` | `t *testing.T` | Düzgünlük yoxlaması |
| Benchmark | `Benchmark` | `b *testing.B` | Performans ölçmə |
| Example | `Example` | (yoxdur) | Sənədləşmə + çıxış testi |

## Əsas terminlər
- Unit Test (vahid test)
- Positive/Negative Path (müsbət/mənfi ssenari)
- Table Test (cədvəl testi)
- Mocking ( saxta mühit yaratma)
- httptest.Server (mock HTTP server)
- ResponseRecorder (cavab qeydedicisi)
- HandlerFunc (handler adapteri)
- ServeMux (marşrutlaşdırıcı)
- Black-box Test (`_test` paketi)
- Example Function (nümunə funksiyası)
- Output Marker (`// Output:`)
- Benchmark / ns/op / B/op / allocs/op
- b.N (iterasiya sayəcısı)
- Continuous Integration (fasiləsiz inteqrasiya)

## Praktik nəticə
- Testləri inkişaf prosesinin sonunda yox, ONUN İÇİNDƏ yaz — `go test` build qədər təbii addımdır.
- Xarici servislərdən asılı testləri `httptest` ilə mock-la — CI internet olmadan da işləməlidir.
- Bir funksiyanın çox ssenarisini table test ilə ört — yeni hal əlavə etmək kod yazmaq demək deyil.
- Handler-ları `NewRecorder` + `ServeHTTP` ilə serveri qaldırmadan test et.
- Public paketlərə Example funksiyaları yaz — istifadəçi sənədləri və çıxış qarantiyası eyni anda.
- Performans həssas kodda benchmark + `-benchmem` ilə həlləri müqayisə et; mikro-optimizasiyalarda `strconv` > `fmt` (bu fəsildə 5x fərq sübut olundu).
- `t.Fatal` = dayandır, `t.Error` = davam et — mümkün qədər çox yoxlama aparmaq üçün `t.Error` seç.

## Mənbə
Səhifələr: 232-257 (PDF), kitab səhifələri 211-236
