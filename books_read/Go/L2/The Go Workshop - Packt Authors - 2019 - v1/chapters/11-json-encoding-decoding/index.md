# Chapter 11 — Encoding and Decoding (JSON) (JSON ilə İş)

## Bu fəsil nədən bəhs edir?

JSON formatı (struktur, data tipləri, XML müqayisəsi), json.Unmarshal (struct-a decode,
struct tag-lər, nested object), json.Valid, json.Marshal (encoding), tag atributları
(omitempty, `-`, ad dəyişmə), MarshalIndent (pretty print), naməlum JSON strukturu
(map[string]interface{} + type switch), GOB protokolu (binary, Go-only, rpc) və customer
order activity.

## Əsas fikirlər

### 1. JSON Nədir
**JavaScript Object Notation** — dil-müstəqil TEXT formatı; XML-dən az-verbose,
yüngül, oxunaqlı. REST API-lərdə, DB-lərdə (JSON tip sahəsi), statik web renderində
istifadə olunur.

**Struktur:** `"key": value` cütlükləri; `{}` obyekt, `[]` array; vergüllə ayrılır.

**6 data tipi:**
| Tip | Nümunə |
|---|---|
| String | `{"firstname": "Captain"}` |
| Number | `{"age": 32}` |
| Array | `{"hobbies": ["Go", "Shield"]}` |
| Boolean | `{"ismarried": false}` |
| Null | `{"middlename": null}` |
| Object | struct-a uyğun: `{"person": {...}}` |

**JSON object ≈ Go struct:** sahələr uyğun gəlir; nested object → embedded struct.

### 2. Unmarshal — JSON → Struct
```go
func Unmarshal(data []byte, v interface{}) error
```
**Şərtlər:** data = []byte; v = POINTER (nil/pointer-olmayan → xəta).

**Kitabdan kod nümunəsi:**
```go
type greeting struct {
    SomeMessage string `json:"message"`     // tag = JSON açar adı
}

data := []byte(`{"message": "Greetings fellow gopher!"}`)
var g greeting
err := json.Unmarshal(data, &g)        // &g — POINTER
if err != nil { fmt.Println(err) }
fmt.Println(g.SomeMessage)
```
**Vacib qayda:** struct sahəsi EXPORTED olmalıdır (böyük hərf) — yoxsa marshaler
GÖRMÜR. Tag unexported sahədədirsə → COMPILE XƏTASI.

**Field matching ardıcıllığı:** 1) exported + tag; 2) ad case-SENSITIVE uyğun;
3) case-INSENSITIVE uyğunluq.

### 3. json.Valid — Ön Yoxlama
```go
if !json.Valid(data) {
    fmt.Printf("JSON is not valid: %s", data)
    os.Exit(1)
}
```
Unmarshal-dan ƏVVƏL — xətalı JSON-un təhlili və aydın mesajı.

### 4. Nested Object Decode
**Kitabdan kod nümunəsi:**
```go
type person struct {
    Lastname  string  `json:"lname"`
    Firstname string  `json:"fname"`
    Address   address `json:"address"`    // NESTED struct
}
type address struct {
    Street  string `json:"street"`
    City    string `json:"city"`
    State   string `json:"state"`
    ZipCode int    `json:"zipcode"`
}
var p person
json.Unmarshal(data, &p)      // p.Address nested doldurulur
```

### 5. Marshal — Struct → JSON
```go
func Marshal(v interface{}) ([]byte, error)
```
Tag YOXDURSA sahə ADI açar olur (case-keçirmə YOXDUR — SomeMessage olduğu kimi).

**Tag atributları:**
```go
type book struct {
    ISBN          string `json:"isbn"`
    Title         string `json:"title"`
    YearPublished int    `json:"yearpub,omitempty"`  // BOŞSA JSON-a DÜŞMÜR (0)
    Author        string `json:"author"`
    CoAuthor      string `json:"coauthor,omitempty"` // "" düşmür
    Secret        string `json:"-"`                  // TAMAMİLƏ GİZLİ
}
```
- `omitempty` — zero value-dan JSON-dan YOX
- `json:"-"` — heç vaxt encode olunmur (parol, token!)
- `json:",omitempty"` — ad saxlanır (sahə adı), yalnız boşluq nəzarəti

**Vergüldən SONRA BOŞLUQ YOX** — `json:"yearpub, omitempty"` → go vet xətası.

**Xəta nümunəsi:** marshal xətası olsa BELƏ nəticə qaytarır (yearpub:0 səhv dəyərlə)
— error yoxlaması MÜTLƏQ.

### 6. MarshalIndent — Pretty Print
```go
prettyPrint, err := json.MarshalIndent(p, "", "    ")   // prefix, indent
```
Bir sətirlə yox, hər sahə yeni sətirdə, 4 boşluq girinti — böyük strukturlarda
oxunaqlılıq.

### 7. Naməlum JSON Strukturu
**Problem:** 3-cü tərəf metrikası, tez-tez dəyişən JSON, servisi dayandırmadan uyğunlaşma.

**Həll — map[string]interface{}:**
```go
jsonData := []byte(`{"checkNum":123,"amount":200,"category":["gift","clothing"]}`)
var v interface{}
json.Unmarshal(jsonData, &v)      // v = map[string]interface{}

// Tipləri aşkar et (type switch):
data := v.(map[string]interface{})
for k, v := range data {
    switch value := v.(type) {
    case string:   fmt.Println("(string):", k, value)
    case float64:  fmt.Println("(float64):", k, value)   // JSON rəqəmlər həmişə float64!
    case bool:     fmt.Println("(bool):", k, value)
    case []interface{}:                            // array → interface slice
        for i, j := range value { fmt.Println("    ", i, j) }
    default:       fmt.Println("(unknown):", k, value)
    }
}
```
**Vacib:** JSON-da NUMBER → Go-da həmişə float64; map sırası RANDOM.

### 8. GOB — Go-nun Öz Binary Protokolu
**Nədir:** Go↔Go üçün binary encoding; JSON-un string məhdudiyyətləri YOX — yüksək
performans + effektiv (space/processing).

**Xüsusiyyətləri:**
- Yalnız Go (digər dillərlə işləməz — yeganə məhdudiyyət)
- Konfiqurasiya/setup YOXDUR
- Data model TAM uyğunluq tələb etmir: uyğun gələn sahələr istifadə olunur,
  qalanlar ATILIR → legacy uyğunluq
- Pointer/value, int/float — hamısını çevirir
- rpc paketi DEFAULT olaraq gob işlədir — network arxitekturası hazır
- İstifadə: servislərarası (low latency → microservices), FAYLLAR (restart davamlılığı:
  transaction backup, cold-start cache priming — DB-stampede qorunması)

**Kitabdan kod nümunəsi (client≠server modelləri!):**
```go
// Client: pointer user, float64 amount
type TxClient struct {
    ID          string
    User        *UserClient     // pointer
    Amount      float64
}
// Server: value user, *float32 amount, Name sahəsi YOXDUR
type TxServer struct {
    ID          string
    User        UserServer      // value
    Amount      *float32        // pointer + fərqli tip!
}

var net bytes.Buffer                    // dummy şəbəkə
enc := gob.NewEncoder(&net)
enc.Encode(clientTx)                    // ENCODE

func sendToServer(net io.Reader) (*TxServer, error) {
    tx := &TxServer{}
    dec := gob.NewDecoder(net)
    err := dec.Decode(tx)             // DECODE — fərqli struct-a!
    return tx, err
}
```
Gob pointer↔value, float64→float32, Name olmamasını avtomatik idarə edir.

### 9. Customer Order Activity
address/item/order/customer struct-ları; Description omitempty, Fragile omitempty,
Password/Token `json:"-"`; Valid yoxlaması; decode → 2 item əlavə + TotalPrice hesabı
+ pretty print (MarshalIndent).

## Əsas terminlələr
- JSON — key-value text formatı; dil-müstəqil
- Unmarshal — JSON []byte → Go struct (pointer tələb)
- Marshal — Go struct → JSON []byte
- Struct Tag — `json:"ad"` — açar adı nəzarəti
- omitempty — zero value JSON-a düşmür
- `json:"-"` — sahə tamamilə gizli
- json.Valid — JSON düzgünlüyünün ön yoxlaması
- MarshalIndent — pretty print (girintili)
- map[string]interface{} — naməlum schema həlli
- Type Switch — interface{} dəyərlərinin tip aşkarlanması
- float64 — JSON rəqəmlərinin Go-da universal forması
- GOB — Go-only binary protokolu; model uyğunsuzluq tolerantlığı
- rpc — servislərarası çağırış; default gob
- Cache Priming — fayldan başlanğıc cache yükləməsi

## Praktik nətidə

(1) Marshal/Unmarshal sahələri EXPORTED tələb edir — tag unexported sahədə compile
xətası verir (yaxşı müdafiə). (2) Tag açar adını dəyişmə azadlığı verir — struct adı
Go konvensiyası, JSON açarı API konvensiyası. (3) omitempty boş sahələri təmizləyir;
`-` parol/token kimi sahələri TAM gizlədir. (4) Vergüldən sonra BOŞLUQ = tag
səhvi — go vet tutur. (5) Marshal xətası nəticəni ləğv ETMİR — error yoxlanmasız
səhv data istifadə olunar. (6) Naməlum JSON: interface{} unmarshal → map + type
switch; rəqəmlər float64. (7) json.Valid girişi — xətalı məlumatı aydın report
etmək üçün. (8) MarshalIndent API cavabları/debug çapları üçün. (9) Go↔Go
kommunikasiyada gob seç — daha sürətli, model dəyişikliyinə davamlı; rpc onsuz da
gob işlədir. (10) Gob fayllarda: transaction backup + cache priming — restart
davamlılığının sadə yolu.

## Mənbə
Pages: 375-418 (PDF 408-453)
