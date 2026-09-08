# Chapter 10 — Templates (səh. 208-220)

## Bu fəsil nədən bəhs edir?

`text/template` və `html/template` paketləri: struct data ilə şablonların
doldurulması (`{{.Sahə}}`), actions (if/else, range), müqayisə operatorları,
FuncMap ilə xüsusi funksiyalar, nested templates (block/define) və HTML
təhlükəsizliyi (avtomatik escaping).

## Əsas fikirlər

### 1. Şablonun struct ilə doldurulması
**Nədir:** Data əsaslı şablonlar — kodu dəyişmədən mətn nəticələri
generasiya etmək (məktublar, hesabatlar, reset-password e-poçtları).

**Kitabdan kod nümunəsi:**
```go
type User struct {
    Name   string
    UserId string
    Email  string
}

const Msg = `Dear {{.Name}},
You were registered with id {{.UserId}}
and e-mail {{.Email}}.
`

func main() {
    u := User{"John", "John33", "john@gmail.com"}

    t := template.Must(template.New("msg").Parse(Msg))
    err := t.Execute(os.Stdout, u)   // stdout-a yazılır
    if err != nil {
        panic(err)
    }
}
```

**Sub-kod izahı:**
- `{{.Name}}` → data-nın `Name` sahəsi (`.` = cari data)
- `template.New("ad")` → şablon obyekti yaradır
- `.Parse(Msg)` → şablon sətrini təhlil edir (xəta qaytara bilər)
- `template.Must(...)` → parse xətasında panic edən wrapper — təhlükəsiz
  qısayol
- `.Execute(writer, data)` → doldurulmuş nəticəni `io.Writer`-a yazır
  (Chapter 7 Reader/Writer dünyası)

### 2. Actions — if/else
**Kitabdan kod nümunəsi (cinsə görə müraciət):**
```go
type User struct {
    Name   string
    Female bool
}

const Msg = `
{{if .Female}}Mrs.{{- else}}Mr.{{- end}} {{.Name}},
Your package is ready.
Thanks,
`

u1 := User{"John", false}
u2 := User{"Mary", true}
t := template.Must(template.New("msg").Parse(Msg))
t.Execute(os.Stdout, u1)   // Mr. John, ...
t.Execute(os.Stdout, u2)   // Mrs. Mary, ...
```
- `{{if .Şərt}}...{{else}}...{{end}}` → bool sahəyə görə budaqlanma
- `{{- ...}}` →多余的 boşluqları yeyir

### 3. Müqayisə operatorları (Table 10.1)
Şablonlarda Go operatorları əvəzinə adlar istifadə olunur:

| Go | Template |
|---|---|
| `<` | `lt` |
| `>` | `gt` |
| `<=` | `le` |
| `>=` | `ge` |
| `==` | `eq` |
| `!=` | `ne` |

**Kitabdan kod nümunəsi (score-a görə səviyyə):**
```go
type User struct {
    Name  string
    Score uint32
}

const Msg = `
{{.Name}} your score is {{.Score}}
your level is:
{{if le .Score 50}}Amateur
{{else if le .Score 80}}Professional
{{else}}Expert
{{end}}
`
```
- `le .Score 50` → `.Score <= 50` (prefiks notasiyası: operator birinci)

### 4. range — kolleksiya iterasiyası
```go
const msg = `
The musketeers are:
{{range .}}{{print .}} {{end}}
`

musketeers := []string{"Athos", "Porthos", "Aramis", "D'Artagnan"}
t := template.Must(template.New("msg").Parse(msg))
t.Execute(os.Stdout, musketeers)
// Athos Porthos Aramis D'Artagnan
```
- `{{range .}}...{{end}}` → slice/map üzrə döngü; daxildə `.` = cari element
- `{{print .}}` → elementi çap edir

### 5. Şablon funksiyaları
**Hazır funksiyalar:** `{{slice . 3}}` → `x[3]` — index 3 elementi:
```go
const Msg = `
The fourth musketeer is:
{{slice . 3}}
`
// D'Artagnan
```

**FuncMap — öz funksiyalarınız:**
```go
const Msg = `
The musketeers are:
{{join . ", "}}
`

funcs := template.FuncMap{"join": strings.Join}   // strings.Join → "join"

t, err := template.New("msg").Funcs(funcs).Parse(Msg)
if err != nil {
    panic(err)
}
t.Execute(os.Stdout, musketeers)
// Athos, Porthos, Aramis, D'Artagnan
```
- `template.FuncMap{"ad": fn}` → şablonda çağırıla bilən funksiya xəritəsi
- `.Funcs(funcs)` → parse-dən əvvəl qeydiyyat olunmalıdır

### 6. Nested templates (block/define)
**Nədir:** Şablon içində şablon — hissələrin təkrar istifadəsi.

**Kitabdan kod nümunəsi:**
```go
const Header = `
{{block "hello" .}}Hello and welcome{{end}}`

const Welcome = `
{{define "hello"}}
{{range .}}{{print .}} {{end}}
{{end}}
`

func main() {
    musketeers := []string{"Athos", "Porthos", "Aramis", "D'Artagnan"}

    helloMsg, err := template.New("start").Parse(Header)
    if err != nil {
        panic(err)
    }

    welcomeMsg, err := template.Must(helloMsg.Clone()).Parse(Welcome)
    if err != nil {
        panic(err)
    }

    if err := helloMsg.Execute(os.Stdout, musketeers); err != nil {
        panic(err)
    }
    if err := welcomeMsg.Execute(os.Stdout, musketeers); err != nil {
        panic(err)
    }
}
```

**Sub-kod izahı:**
- `{{block "ad" .}}...{{end}}` → "ad" şablonunu işə salır (yoxdursa
  default gövdə işlədilir)
- `{{define "ad"}}...{{end}}` → "ad" şablonunun tərifini verir
- `helloMsg.Clone()` → birincinin klonu ikinci parse üçün — hər ikisi
  eyni "hello" adına baxır
- Nəticə: "Hello and welcome" + musketorların siyahısı

### 7. HTML şablonları (html/template)
**Nədir:** `text/template`-in HTML təhlükəsizliyi ilə variantı — **code
injection** hücumlarının qarşısını avtomatik escaping ilə alır.

**Kitabdan kod nümunəsi:**
```go
const Page = `
<html>
<head>
    <title>{{.Name}}'s Languages</title>
</head>
<body>
    <ul>
    {{range .Languages}}<li>{{print .}}</li>{{end}}
    </ul>
</body>
</html>
`

type UserExperience struct {
    Name      string
    Languages []string
}

func main() {
    languages := []string{"Go", "C++", "C#"}
    u := UserExperience{"John", languages}

    t := template.Must(template.New("web").Parse(Page))
    t.Execute(os.Stdout, u)
}
```

**Nəticə (avtomatik escaping diqqətlə baxın):**
```html
<html>
<head><title>John's Languages</title></head>
<body>
<ul>
<li>Go</li><li>C&#43;&#43;</li><li>C#</li>
</ul>
</body>
</html>
```
- `C++` → `C&#43;&#43;` — `+` simvolu HTML entity-yə çevrildi
- API funksiyaları `text/template` ilə eynidir; fərq yalnız təhlükəsizlik
  qatındadır

## Əsas terminlər
- Template (şablon) — data ilə doldurulan mətn forması
- Action — `{{...}}` daxilindəki əmr (if, range, print...)
- FuncMap — şablona qeyd olunan xüsusi funksiyalar xəritəsi
- block/define — nested (iç-içe) şablon tərifləri
- Escaping — təhlükəli simvolların HTML entity-ə çevrilməsi
- html/template — HTML-safe şablon paketi

## Praktik nəticə
Şablonlar kodu datadan ayırır: mətn dəyişəndə proqramı yenidən
kompilyasiya etmək lazım deyil. Workflow: `template.New` → `Funcs` (lazımdırsa)
→ `Must(Parse(...))` → `Execute(writer, data)`. Müqayisələr üçün lt/gt/le/ge/eq/ne;
iterasiya üçün range; mürəkkəb üçün block/define + Clone. **Web üçün həmişə
`html/template`** — avtomatik escaping XSS qoruması verir.

## Mənbə
Pages: 208-220 (PDF 208-220)
