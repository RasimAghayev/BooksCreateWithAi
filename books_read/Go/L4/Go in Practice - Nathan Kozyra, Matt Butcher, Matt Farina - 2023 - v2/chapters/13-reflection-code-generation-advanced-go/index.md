# Chapter 13 — Reflection, code generation, and advanced Go (Refleksiya, kod generasiyası və qabaqcıl Go)

## Bu chapter nədən bəhs edir?

Refleksiyanın 3 anlayışı (value/type/kind), type switch vs kind switch, interfeys
implementasiya yoxlaması (nil pointer hiyləsi), struct gəzinti (rekursiv walk), öz
struct tag-lərin (ini: format) ilə Marshal/Unmarshal qurulması, go:generate ilə kod
generasiyası (queue generator nümunəsi) və C interop (cgo, Ch1-dən tanış).

## Əsas fikirlər

### 1. Value / Type / Kind
- **Value** — dəyişənin məzmunu (`reflect.Value`); **Type** — konkret tip
  (`reflect.Type`); **Kind** — primitiv quruluş (struct/ptr/int/string/slice/func...)
  (`reflect.Kind`). MyInt tipi: type=MyInt, kind=int.

### 2. Type Switch vs Kind Switch
**Type switch limiti:** `type MyInt int64` — int64 case-inə DÜŞMÜR (default-a).
**Kind switch həlli:** tipləri qruplaşdırır — 10 int case-i 2-yə endirir:
```go
func sum(v ...interface{}) float64 {
    var res float64
    for _, val := range v {
        ref := reflect.ValueOf(val)
        switch ref.Kind() {
        case reflect.Int, reflect.Int64:
            res += float64(ref.Int())     // Int() — int ailəsinin ən böyüyünə
        case reflect.Uint8:
            res += float64(ref.Uint())
        case reflect.String:
            a, _ := strconv.ParseFloat(ref.String(), 64)
            res += a
        default:
            fmt.Printf("Unsupported type %T. Ignoring.\n", val)
        }
    }
    return res
}
```
Type switch error ayırdetməsində (catch blokları kimi) daha uyğundur; say yoxlama
— kind.

### 3. İnterfeys İmplementasiya Yoxlaması
**Statik:** type assertion — `_, ok := v.(fmt.Stringer)`.
**Runtime (refleksiya ilə):** interfeysə birbaşa refleksiya YOXDUR (reflect.Interface
yoxdur) — **nil pointer hiyləsi**:
```go
stringer := (*fmt.Stringer)(nil)          // tipi bilinən nil
writer := (*io.Writer)(nil)
func implements(concrete interface{}, target interface{}) bool {
    iface := reflect.TypeOf(target).Elem()   // pointer-in hədəf TİPİ = interfeys tipi
    t := reflect.ValueOf(concrete).Type()
    return t.Implements(iface)               // Type.Implements()
}
```
`*main.Name is a Stringer / not a Writer` — Go-da interfeys "qarşı müqayisə təsviri"
dir; declaration YOX (composition, inheritance deyil).

### 4. Struct Rekursiv Gəzinti
```go
func walk(u interface{}, depth int) {
    val := reflect.Indirect(reflect.ValueOf(u))   // pointer-i SAMİYƏTƏN aç (pointer
    t := val.Type()                                 // deyilsə özünü qaytarır)
    fmt.Printf("%sValue is type %q (%s)\n", tabs, t, val.Kind())
    if val.Kind() == reflect.Struct {
        for i := 0; i < t.NumField(); i++ {
            field := t.Field(i)                     // StructField: ad/tip/TAG
            fieldVal := reflect.Indirect(val.Field(i))   // EYNİ indeks = dəyər
            if fieldVal.Kind() == reflect.Struct {
                walk(fieldVal.Interface(), depth+1)
            }
        }
    }
}
```
Açar: **Type.Field(i) (metadata) + Value.Field(i) (dəyər) — eyni i**. Digər alətlər:
MethodByName, AssignableTo, Comparable. Diqqət: reflect metodları error YOX — PANIC
edir (Ch4 recover patternləri tətbiq et).

### 5. Öz Struct Tag-lər — ini: Encoder/Decoder
Tag = backtick-lı freeform string; `NAME:"VALUE"` formatı de-fakto standartdır və
`StructField.Tag.Get("ini")` ilə parse olunur. (Validasiya regex tag-ləri də mümkün:
`validate:"^[a-z]+$"`.)

**Kitabdan kod nümunəsi (Marshal):**
```go
type Processes struct {
    Total    int     `ini:"total"`
    Running  int     `ini:"running"`
    Load     float64 `ini:"load"`
}
func fieldName(field reflect.StructField) string {
    if t := field.Tag.Get("ini"); t != "" { return t }
    return field.Name                              // tag-siz fallback
}
func Marshal(v interface{}) ([]byte, error) {
    var b bytes.Buffer
    val := reflect.Indirect(reflect.ValueOf(v))
    if val.Kind() != reflect.Struct { return nil, errors.New("structs only") }
    t := val.Type()
    for i := 0; i < t.NumField(); i++ {
        fmt.Fprintf(&b, "%s=%v\n", fieldName(t.Field(i)), val.Field(i).Interface())
    }
    return b.Bytes(), nil
}
```

**Unmarshal:** sətir-sətir splitN("=") → setField (kind switch: Int/Float64/String/Bool
strconv ilə; `v.Field(i).Set(reflect.ValueOf(converted))`). Naməlum kind — skip (error
deyil). `-` tag = iqnor konvensiyası (json-dan). Bu, encoding/json-in özünün necə
işlədiyinin miniatürüdür.

### 6. go:generate — Kod Generasiyası
**Motivasiya:** reflection runtime xərclidir; codegen = compile-olunmuş tip-təhlükəsiz
kod — daha sürətli və sadə. Generics-dən əvvəl preskripsiya idi; hələ də niş rol oynayır
(protobuf/gRPC/Thrift/SQL kitabxanaları hamısı generator-dur).

**Mexanizm:** faylın ilk sətrində `//go:generate COMMAND ARGS` → `go generate` komandanı
icra edir:
```go
//go:generate echo hello        // ən sadə nümunə
```

**Queue generator (template əsaslı):**
```go
var tpl = `package {{.Package}}
type {{.MyType}}Queue struct { q []{{.MyType}} }
func New{{.MyType}}Queue() *{{.MyType}}Queue { ... }
func (o *{{.MyType}}Queue) Insert(v {{.MyType}}) { o.q = append(o.q, v) }
func (o *{{.MyType}}Queue) Remove() {{.MyType}} { ... }`

func main() {
    tt := template.Must(template.New("queue").Parse(tpl))
    for i := 1; i < len(os.Args); i++ {
        dest := strings.ToLower(os.Args[i]) + "_queue.go"
        file, _ := os.Create(dest)
        tt.Execute(file, map[string]string{
            "MyType":  os.Args[i],
            "Package": os.Getenv("GOPACKAGE"),   // go generate təyin edir
        })
        file.Close()
    }
}
```
İstifadə: `//go:generate ./queue MyInt` → `go generate` → `myint_queue.go` yaranır →
`go run myint.go myint_queue.go`.

**Dövr qaydaları:** (1) generasiya DEVELOPMENT dövrünə aiddir — build/runtime-a YOX;
(2) generasiya olunmuş kodu VCS-ə commit et — istifadəçi generator işlətməməli;
(3) generator kodun compile olmasını tələb ETMİR — build-dən əvvəl işlət;
(4) generator $PATH-də və ya lokal `./queue` kimi çağrıla bilər; $GOPACKAGE və b.
env-lər hazır gəlir. İstifadə sahələri: type-safe kolleksiyalar, DB→struct, JSON schema→kod.

## Əsas terminlələr
- reflect.Value / Type / Kind — refleksiya üçlüsü
- Kind switch — quruluşa görə birləşdirilmiş switch
- Nil pointer hiyləsi — interfeys tipinə refleksiya üçün (*Iface)(nil)
- reflect.Indirect — pointer-samiiyətə açma (pointer deyilsə no-op)
- StructField.Tag.Get — tag parser
- go:generate / $GOPACKAGE — generasiya direktivi / mühit dəyişəni
- text/template ilə kod generasiyası
- Codegen vs reflection — compile-time vs runtime tip emalı

## Praktik nətidə

Qabaqcıl qərarlar: (1) tiplər AİLƏ üzrə emal — kind switch; konkret tip — type switch
(və ya 1.18+ generics!); (2) runtime interfeys yoxlaması — (*Iface)(nil) + Type.Implements;
(3) struct introspeksiya — Indirect + NumField/Field(i) (tip+dəyər eyni indeks);
(4) öz tag formatın — `name:"value"` + Tag.Get; Marshal/Unmarshal patternini izlə;
(5) reflection ağır/məkrli (panic!) — alternativ: text/template generator + go:generate;
(6) generasiya = dev-time; nəticə VCS-də; (7) protobufların ast-ninə (go/ast) da bax —
amma 80% hallarda template kifayətdir.

## Mənbə
Pages: 316-342 (PDF 337-363)
