# Chapter 2 — Serializing Data (səh. 31-46)

## Bu fəsil nədən bəhs edir?

Serializasiya (marshaling) nüansları — 7 resept: encoding/gob ilə event
axını, anonim struct ilə kompleks JSON parse, Decoder ilə JSON stream,
itkin sahələrin tutulması (pointer texnikası), custom tipin MarshalJSON/
UnmarshalJSON, dinamik tiplərin mapstructure ilə dispatch-i və struct
tag-lərin reflect ilə parse edilməsi.

## Əsas fikirlər

### Recipe 8 — encoding/gob ilə event axını
**Tapşırıq:** Müştəri davranışı üçün Add/Checkout eventləri — Go↔Go
 kommunikasiyası.

```go
// Kind is event kind.
type Kind string
const (
    AddKind      Kind = "add"
    CheckoutKind Kind = "checkout"
)
type Event interface {
    Kind() Kind
}

// Konkret eventlər:
type Add struct {
    Time     time.Time
    ID       string
    User     string
    Item     int // SKU
}
func (*Add) Kind() Kind { return AddKind }

type Checkout struct {
    Time time.Time
    Cart string
    User string
}
func (*Checkout) Kind() Kind { return CheckoutKind }

// gob qeydiyyatı — interfeys kimi kodlamaq üçün:
func init() {
    gob.Register(&Add{})
    gob.Register(&Checkout{})
}

// Encoder — axın boyu EYNİ encoder istifadə olunmalıdır:
type Encoder struct {
    enc *gob.Encoder
}
func NewEncoder(w io.Writer) *Encoder {
    return &Encoder{gob.NewEncoder(w)}
}
func (e *Encoder) Encode(evt Event) error {
    return e.enc.Encode(&evt)   // İNTERFEYSİN POİNTERİ — vacib!
}

// Handler — decode + dispatch:
func eventHandler(r io.Reader) error {
    dec := gob.NewDecoder(r)
    for {
        var e Event
        err := dec.Decode(&e)
        if err == io.EOF {
            break
        }
        if err != nil {
            return err
        }
        switch e.Kind() {
        case AddKind:
            handleAdd(e.(*Add))
        case CheckoutKind:
            handleCheckout(e.(*Checkout))
        default:
            return fmt.Errorf("unknown event kind: %s", e.Kind())
        }
    }
    return nil
}
```

**Vacib qaydalar:**
- gob Go-specific-dir — digər dillərlə uyğunlaşmaq çətindir
- interfeysi yox, **interfeysin pointerini** (`&evt`) encode et — yoxsa
  konkret tip interfeys məlumatı olmadan kodlanır
- eyni encoder bütün axın üçün (Decoder/Encoder saxlanan sahədə)

### Recipe 9 — anonim struct ilə kompleks JSON
**Tapşırıq:** `{"payments":[{"id":..,"date":..,"amount":..},...]}` —
yalnız cəmi lazımdır.

```go
func totalPayments(r io.Reader) (float64, error) {
    // Anonim struct — YALNIZ lazımi sahələr:
    var reply struct {
        Payments []struct {
            Amount float64
        }
    }
    if err := json.NewDecoder(r).Decode(&reply); err != nil {
        return 0, err
    }
    total := 0.0
    for _, p := range reply.Payments {
        total += p.Amount
    }
    return total, nil
}
```
- encoding/json bilinməyən sahələri İGNOR edir → tam model lazım deyil
- anonim struct → namespace çirklənməsi yoxdur; sahələr yuxarı hərflə
  başlamalı (exported) — anonim olsa belə
- Ad uyğunluğu: `Amount ↔ amount` heuristik; dəqiq idarə üçün tag:
```go
type User struct {
    Name string `json:"login"`
    ID   int    `json:"uid"`
}
```

### Recipe 10 — streaming JSON (Decoder döngüsü)
**Tapşırıq:** socket-dən sətir-sətir JSON log udma.

```go
func ingestLogs(r io.Reader, handler func(Log)) error {
    dec := json.NewDecoder(r)
    for {
        var l Log               // hər iterasiyada YENİ struct
        err := dec.Decode(&l)
        if err == io.EOF {
            break
        }
        if err != nil {
            return err
        }
        handler(l)
    }
    return nil
}
```
- JSON protokolu axını dəstəkləmir — Decoder **bir obyekt bir dəfə** oxuyur
  ("JSON lines" praktikası standartdır)
- Struct yenidən istifadə etmək yaddaş qənaətidir, amma sıfırlama unudulsa
  xəta — hər dəfə yeni struct daha etibarlı
- Bu "ingestion + handler" pattern-i net/http, gRPC-də də var — Eric
  Raymond-un Rule of Separation: mechanism-i policy-dən ayır

### Recipe 11 — itkin JSON sahələri (pointer yoxlama)
**Tapşırıq:** Level=0 və Time="January 1, 1" — zero value mü patientdir?
Həll: **pointer sahələr** nil yoxlaması ilə.

```go
// Log is a log event.
type Log struct {
    Time    time.Time
    Level   int
    Message string
}

func parseLog(data []byte) (*Log, error) {
    var l struct {
        Time    *time.Time   // POINTER — iştirak etməyibsə nil qalır
        Level   *int
        Message string
    }
    if err := json.Unmarshal(data, &l); err != nil {
        return nil, err
    }
    if l.Time == nil {
        return nil, fmt.Errorf("missing `time` field")
    }
    if l.Level == nil {
        return nil, fmt.Errorf("missing `level` field")
    }
    if l.Message == "" {
        return nil, fmt.Errorf("missing `message` field")
    }
    return &Log{*l.Time, *l.Level, l.Message}, nil
}
```
- Zero value problemi: `var l Log` → Level=0, Time=0001-01-01 — "göndərilməyib"
  ilə "0 göndərilib" ayırd edilə bilmir
- Pointer sahə: JSON-da yoxdursa → nil → xəta
- Mübadilə: iki dəfə kopyalama; amma istifadəçi tərəfdə pointer-siz, təmiz
  Log tipi qalır

### Recipe 12 — custom tiplərin serializasiyası
**Tapşırıq:** Rekursiv Stack (linked list) → JSON.

```go
// Stack is a LIFO data structure.
type Stack struct {
    Value string
    Next  *Stack
}

// MarshalJSON implements json.Marshaler.
func (s *Stack) MarshalJSON() ([]byte, error) {
    var values []string
    for s != nil {
        values = append(values, s.Value)
        s = s.Next
    }
    return json.Marshal(values)      // slices → Array
}

// UnmarshalJSON implements json.Unmarshaler.
func (s *Stack) UnmarshalJSON(data []byte) error {
    var values []string
    if err := json.Unmarshal(data, &values); err != nil {
        return err
    }
    var node *Stack
    for i := len(values) - 1; i >= 0; i-- {   // sondan qur — LIFO
        node = &Stack{values[i], node}
    }
    *s = *node     // POİNTERİN DƏYƏRİNİ yenilə — özünü yox!
    return nil
}
```
- JSON tipləri məhduddur: `[]any`(Array), `bool`, `nil`(null),
  `float64`(Number), `map[string]any`(Object), `string`
- Strategiya: custom tipi ən yaxın JSON tipinə çevir → built-in Marshal
- UnmarshalJSON mütləq POINTER receiver + `*s = *node`
- Eyni yanaşma: YAML → MarshalYAML/UnmarshalYAML

### Recipe 13 — dinamik tiplərin unmarshal-ı
**Tapşırıq:** HTTP-dən gələn JSON-lər `type` sahəsinə görə Login/Message.

```go
func handler(data []byte) error {
    // 1) Əvvəlcə map-ə decode:
    var obj map[string]any
    if err := json.Unmarshal(data, &obj); err != nil {
        return err
    }
    // 2) type sahəsini çıxar:
    val, ok := obj["type"]
    if !ok {
        return fmt.Errorf("missing `type` in %+v", obj)
    }
    typ, ok := val.(string)
    if !ok {
        return fmt.Errorf("`type` is not a string - %v of %T", val, val)
    }
    // 3) Tipə görə dispatch — mapstructure ilə map → struct:
    switch typ {
    case "login":
        var l Login
        if err := mapstructure.Decode(obj, &l); err != nil {
            return err
        }
        return handleLogin(l)
    case "message":
        var m Message
        if err := mapstructure.Decode(obj, &m); err != nil {
            return err
        }
        return handleMessage(m)
    default:
        return fmt.Errorf("unknown message type: %q", typ)
    }
}
```
- `map[string]any` + bir çox type assertion → çirkin kod;
  `github.com/mitchellh/mapstructure` map-i struct-a çevirir
- Qiymət: iki serializasiya — dinamik mesajlarda başqa yol yoxdur

### Recipe 14 — struct tag parse (reflect)
**Tapşırıq:** ORM üçün `db:"kolon"` tag-lərini çıxar.

```go
// Log is a log structure.
type Log struct {
    Time  time.Time `db:"ts"`
    Level int       `db:"level"`
    Text  string    `db:"message"`
}

func parseStructTags(s any) (map[string]string, error) {
    typ := reflect.TypeOf(s)
    if typ.Kind() != reflect.Struct {
        return nil, fmt.Errorf("%s is not a struct", typ)
    }
    m := make(map[string]string)
    for i := 0; i < typ.NumField(); i++ {
        fld := typ.Field(i)
        if dbName := fld.Tag.Get("db"); dbName != "" {
            m[fld.Name] = dbName
        }
    }
    return m, nil
}
// Nəticə: Time→ts, Level→level, Text→message
```
- `any` parametri adətən "code smell" — amma reflection API-da məcburən
- Tag formatı: `key:"value" key:"value" ...` — `Tag.Get(key)` ilə oxunur
- Mürəkkəb dəyərlər (protobuf misalı): `protobuf:"fixed64,1,opt,..."`
  json:"value,omitempty"` — öz formatınıız üçün json/protobuf üslubunu
  kopyalayın, icad etməyin

## Final Thoughts-dən seçmə

Format seçimi meyarları: performans (sürət + bayt), dəstəklənən tiplər
(JSON-da timestamp yoxdur!), dil dəstəyi, schema dəstəyi. Mümkünsə
bir neçə formatla PoC edin. ("CSV-dən qaçın — dəhşətli formatdır.")

## Əsas terminlər
- Serialization/Marshaling — struct → bayt ardıcıllığı
- encoding/gob — Go-özəl binary format
- gob.Register — interfeys tiplərinin qeydiyyatı
- Anonymous struct — adsız, lokal struct (namespace təmizliyi)
- Zero value problem — "0 göndərilib" ↔ "göndərilməyib" ikiliyi
- json.Marshaler/Unmarshaler — custom serializasiya interfeysləri
- mapstructure — map → struct üçün paket
- Struct tag — sahə metadatası (key:"value")
- JSON lines — sətir-sətir JSON axını
- Rule of Separation — mechanism/policy ayrılığı

## Praktik nəticə
JSON-in tip dünyası ilə qarşıdurması: lazımi sahələr üçün anonim struct;
itkin sahələr üçün pointer + nil yoxlaması; custom tiplər üçün
MarshalJSON/UnmarshalJSON (ən yaxın JSON tipinə çevir); dinamik mesajlar
üçün map[string]any + type sahəsi + mapstructure. Go↔Go üçün gob asandır,
amma interfeys pointeri encode edin və eyni encoderi axın boyu saxlayın.
Tag-lər reflect-in Tag.Get API-si ilə oxunur.

## Mənbə
Pages: 31-46 (PDF 31-46)
