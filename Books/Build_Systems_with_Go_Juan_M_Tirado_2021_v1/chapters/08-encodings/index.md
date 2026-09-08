# Chapter 8 — Encodings (səh. 157-182)

## Bu fəsil nədən bəhs edir?

Data serialization formatları: CSV (`encoding/csv`), JSON
(`encoding/json`), XML (`encoding/xml`) və YAML (go-yaml paketi) üzrə
marshal/unmarshal əməliyyatları, struct tag-lərin rolu, XML üçün custom
marshaller yazılması və öz encodinq-imizin (tag əsaslı) qurulması.

## Əsas fikirlər

### 1. CSV (encoding/csv)
**Nədir:** Cədvəli datanın flat mətn formatı; hər sətir = record, vergüllə
 ayrılmış elementlər.

**Oxuma (csv.Reader):**
```go
in := `user_id,score,password
"Gopher",1000,"admin"
"BigJ",10,"1234"
"GGBoom",,"1111"
`
r := csv.NewReader(strings.NewReader(in))   // hər hansı Reader üstünə

for {
    record, err := r.Read()
    if err == io.EOF {
        break
    }
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(record)   // []string kimi: [Gopher 1000 admin]
}
```

**Yazma (csv.Writer):**
```go
out := [][]string{
    {"user_id", "score", "password"},
    {"Gopher", "1000", "admin"},
    {"BigJ", "10", "1234"},
    {"GGBoom", "", "1111"},
}
writer := csv.NewWriter(os.Stdout)
for _, rec := range out {
    err := writer.Write(rec)
    if err != nil {
        panic(err)
    }
}
writer.Flush()   // buferdəki datanı çıxar — unutma!
```

### 2. JSON (encoding/json)
**Nədir:** Yüngül, insan oxunaqlı data mübadilə formatı. `Marshal` (Go →
[]byte JSON) və `Unmarshal` ([]byte JSON → Go) əsas funksiyalardır.

**Əsas tiplərin marshal olunması:**
```go
number, err := json.Marshal(42)                       // "42"
float, _ := json.Marshal(3.14)                       // "3.14"
msg, _ := json.Marshal("This is a msg!!!")           // "\"This is a msg!!!\""
numbers, _ := json.Marshal([]int{1,1,2,3,5,8})       // "[1,1,2,3,5,8]"
aMap, _ := json.Marshal(map[string]int{"one":1,"two":2})  // {"one":1,"two":2}
```

**Unmarshal — pointer tələb edir:**
```go
var recoveredNumber int = -1
err := json.Unmarshal(aNumber, &recoveredNumber)   // & → pointer mütləqdir

recoveredMap := make(map[string]int)
err = json.Unmarshal(aMap, &recoveredMap)
```

**Struct + JSON tag-lər:**
```go
type User struct {
    UserId   string `json:"userId,omitempty"`
    Score    int    `json:"score,omitempty"`
    password string `json:"password,omitempty"`  // unexported — HEÇ VAXT çıxmır
}

userC := User{UserId: "GGBoom", password: "1111"}  // Score = 0
db := []User{userA, userB, userC}

dbJson, err := json.Marshal(&db)
// [{"userId":"Gopher","score":1000},{"userId":"BigJ","score":10},{"userId":"GGBoom"}]
// GGBoom-un score-u YOXDUR — omitempty onu buraxdı

var recovered []User
json.Unmarshal(dbJson, &recovered)   // Score = 0 (zero value) kimi qayıdır
```

**Sub-tag izahı:**
- `json:"userId"` → JSON-dakı sahə adı
- `omitempty` → zero value-dursa marshaldə burax
- `password` (kiçik hərf) → unexported → JSON-ə heç düşmür — parolun
  gizlədilməsi üçün praktik üsul
- **Vacib:** `omitempty` yalnız marshaldə işləyir; unmarshaldə buraxılmış
  sahələr zero value alır

**Indent (oxunaqlı çap):**
```go
var indented bytes.Buffer
json.Indent(&indented, dbJson, "", "    ")
fmt.Println(indented.String())
```

### 3. XML (encoding/xml)
**Nədir:** W3C markup dili; JSON-a bənzər amma daha "verbose" (sözülü).

**Məhdudiyyət — map marshaldə supported DEYİL:**
```go
xml.Marshal(42)                          // <int>42</int>
xml.Marshal(3.14)                        // <float64>3.14</float64>
xml.Marshal("This is a msg!!!")          // <string>This is a msg!!!</string>
xml.Marshal([]int{1,2,2,3,5,8})         // <int>1</int><int>2</int>...
xml.Marshal(map[string]int{"one":1})    // XƏTA: xml: unsupported type: map[string]int
```
Map-in tək XML təsviri yoxdur — buna görə paket onu dəstəkləmir.

**Custom marshaller/unmarshaller (MarshalXML/UnmarshalXML):**
```go
type MyMap map[string]string

func (s MyMap) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    tokens := []xml.Token{start}
    for key, value := range s {
        t := xml.StartElement{Name: xml.Name{"", key}}
        tokens = append(tokens, t, xml.CharData(value), xml.EndElement{t.Name})
    }
    tokens = append(tokens, xml.EndElement{start.Name})
    for _, t := range tokens {
        if err := e.EncodeToken(t); err != nil { return err }
    }
    return e.Flush()
}

func (a MyMap) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
    key := ""
    val := ""
    for {
        t, _ := d.Token()
        switch tt := t.(type) {
        case xml.StartElement:
            key = tt.Name.Local
        case xml.CharData:
            val = string(tt)
        case xml.EndElement:
            if len(key) != 0 {
                a[key] = val
                key, val = "", ""
            }
            if tt.Name == start.Name {
                return nil    // öz elementinin bağlanması → bitdi
            }
        default:
            return errors.New(fmt.Sprintf("unknown %T", t))
        }
    }
}

// Nəticə: <MyMap><one>1</one><two>2</two><three>3</three></MyMap>
```

**Struct tag-lər XML-də də:**
```go
type User struct {
    UserId   string `xml:"userId,omitempty"`
    Score    int    `xml:"score,omitempty"`
    password string `xml:"password,omitempty"`
}

type UsersArray struct {
    Users []User `xml:"users,omitempty"`
}

db := UsersArray{[]User{userA, userB, userC}}
dbXML, _ := xml.Marshal(&db)
// <UsersArray><users><userId>Gopher</userId><score>1000</score></users>...</UsersArray>
```
- Massiv üçün əhatə edici struct (`UsersArray`) lazımdır — kök element
  tələbi

### 4. YAML (go-yaml — üçüncü tərəf)
**Nədir:** JSON-un superset-i, konfiqurasiya fayllarının standartı;
indentation ilə nesting. Standart kitabxanada YOXDUR:

```bash
go get gopkg.in/yaml.v2
```

```go
import "gopkg.in/yaml.v2"

yaml.Marshal(42)                          // "42\n"
yaml.Marshal([]int{1,1,2,3,5,8})          // "- 1\n- 1\n- 2\n..."
yaml.Marshal(map[string]int{"one":1})     // "one: 1\ntwo: 2\n"

type User struct {
    UserId   string `yaml:"userId,omitempty"`
    Score    int    `yaml:"score,omitempty"`
    password string `yaml:"password,omitempty"`
}
db := []User{userA, userB, userC}
dbYaml, _ := yaml.Marshal(&db)
// - userId: Gopher
//   score: 1000
// - userId: BigJ
// ...
var recovered []User
yaml.Unmarshal(dbYaml, &recovered)
```

### 5. Tag-lərlə öz encodinq (custom encoding)
**Nədir:** Chapter 5-dəki reflection tag mexanizmi ilə tam öz serialize
formatı. Nümunə: `pretty` tag-i — string sahələri upper/lower edir.

**Tag tərifi:**
```go
type User struct {
    UserId   string `pretty:"upper"`
    Email    string `pretty:"lower"`
    password string `pretty:"lower"`   // unexported → encode EDİLMƏZ
}
type Record struct {
    Name    string `pretty:"lower" json:"name"`   // tag-lar birgə işləyir
    Surname string `pretty:"upper" json:"surname"`
    Age     int    `pretty:"other" json:"age"`
}
```

**Qaydalar:** unexported yox; yalnız `upper`/`lower`; yalnız string;
 `SahəAdı:Dəyər` formatı; `{...}` mötərizə; rekursiya/kolleksiya yox.

**Marshal implementasiyası (reflection ilə):**
```go
func Marshal(input interface{}) ([]byte, error) {
    var buffer bytes.Buffer
    t := reflect.TypeOf(input)
    v := reflect.ValueOf(input)

    buffer.WriteString("{")
    for i := 0; i < t.NumField(); i++ {
        encodedField, err := encodeField(t.Field(i), v.Field(i))
        if err != nil {
            return nil, err
        }
        if len(encodedField) != 0 {
            if i > 0 && i <= t.NumField()-1 {
                buffer.WriteString(", ")
            }
            buffer.WriteString(encodedField)
        }
    }
    buffer.WriteString("}")
    return buffer.Bytes(), nil
}

func encodeField(f reflect.StructField, v reflect.Value) (string, error) {
    if f.PkgPath != "" {              // unexported sahə süzgəci
        return "", nil
    }
    if f.Type.Kind() != reflect.String {   // yalnız string
        return "", nil
    }
    tag, found := f.Tag.Lookup("pretty")
    if !found {
        return "", nil
    }
    result := f.Name + ":"
    switch tag {
    case "upper":
        result = result + strings.ToUpper(v.String())
    case "lower":
        result = result + strings.ToLower(v.String())
    default:
        return "", errors.New("invalid tag value")
    }
    return result, nil
}

// Nəticə:
// pretty user {UserId:JOHN, Email:john@gmail.com}
// pretty rec {Name:john, Surname:JOHNSON}
// json rec {"name":"John","surname":"Johnson","age":33}
```

**Sub-kod izahı:**
- `f.PkgPath != ""` → unexported sahə deməkdir (exported sahələrdə boşdur)
- `f.Tag.Lookup("pretty")` → tag varsa dəyəri al
- `v.String()` → sahənin cari dəyəri
- Age (int) və password (unexported) encoddə iştirak etmir

## Əsas terminlər
- CSV (Comma Separated Values) — vergüllə ayrılmış cədvəl formatı
- JSON (JavaScript Object Notation) — yüngül data mübadilə formatı
- XML (Extensible Markup Language) — markup formatı
- Marshal (seriyalaşdırma) — Go tipi → format baytları
- Unmarshal (deserializasiya) — format baytları → Go tipi
- omitempty — zero value sahələrin buraxılması
- Tag — struct sahəsi üçün encoding metadatası
- go-yaml — YAML üçün üçüncü tərəf paket (gopkg.in/yaml.v2)

## Praktik nəticə
Bütün encoding-lər eyni nümunəni izləyir: `Marshal`/`Unmarshal` + struct
tag-lər. Unmarshal həmişə pointer (`&var`) tələb edir. Unexported sahələr
heç bir formata düşmür — parolları belə qoruyun. XML map tipini
dəstəkləmir — custom MarshalXML lazımdır. YAML üçün `go get
gopkg.in/yaml.v2`. Öz encoding-inizi reflection + öz tag açarınızla
qura bilərsiniz (`pretty` nümunəsi kimi).

## Mənbə
Pages: 157-182 (PDF 157-182)
