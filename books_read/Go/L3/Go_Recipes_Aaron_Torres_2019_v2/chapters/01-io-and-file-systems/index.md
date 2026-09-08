# Chapter 1 — Ввод-вывод и файловые системы (I/O və fayl sistemləri)

## Bu chapter nədən bəhs edir?

Go-nun əsas I/O interfeysləri (io.Reader/io.Writer), bytes və strings
paktları, fayl/kataloq əməliyyatları (os paketi), CSV formatı, temp
fayllar (ioutil) və şablonlar (text/template, html/template).

## Əsas fikirlər

### 1. Ümumi I/O interfeysləri (io.Reader / io.Writer)
**Nədir:** Go standart kitabxanasının bütün məlumat axını əməliyyatlarının
əsasını təşkil edən iki interfeys.

**Necə işləyir:** Konkret tip deyil, interfeys qəbul edilir — funksiyaya file,
buffer, socket, stdout — hamısı eyni qaydada ötürülə bilər.

**Nəyə lazımdır:** Məlumat mənbəyindən asılı olmayaq ümumi kod yazmaq; şəbəkə
trafiyi və fayl sistemi ilə işdə "stream" (axın) modeli.

**Kitabdan kod nümunəsi:**
```go
type Reader interface {
    Read(p []byte) (n int, err error)
}
type Writer interface {
    Write(p []byte) (n int, err error)
}
// Interfeyslərin kombinasiyası:
type ReadSeeker interface {
    Reader
    Seeker
}
// io.Pipe — yaddaşda konveyer birləşdirir:
func Pipe() (*PipeReader, *PipeWriter)
```

**Copy funksiyası — MultiWriter + buffer:**
```go
func Copy(in io.ReadSeeker, out io.Writer) error {
    w := io.MultiWriter(out, os.Stdout)   // iki yazıçını birləşdirir
    if _, err := io.Copy(w, in); err != nil {
        return err
    }
    in.Seek(0, 0)                          // axını əvvələ qaytar
    buf := make([]byte, 64)                // 64 baytlıq buffer
    if _, err := io.CopyBuffer(w, in, buf); err != nil {
        return err
    }
    return nil
}
```

**Sub-kod izahı:**
- `io.MultiWriter(out, os.Stdout)` → yazılan datanı eyni anda iki yerə yönləndirir
- `io.Copy(w, in)` → Reader-dən Writer-ə bütün axını köçürür (stream kimi)
- `in.Seek(0, 0)` → ikinci oxu üçün axının başına qayıdır (buna görə ReadSeeker lazımdır)
- `io.CopyBuffer` → böyük datanı yaddaşa sığmayanda buffer ilə parça-parça yazır

**PipeExample — bloklanan pipe:**
```go
r, w := io.Pipe()
go func() {          // ayrı goroutine-da yazmaq MÜTLƏQ — pipe bloklayır
    w.Write([]byte("test\n"))
    w.Close()
}()
io.Copy(os.Stdout, r)  // oxu tərəfi yazılanı gözləyir
```

**Üstünlükləri:** Təmiz abstraksiya, axın təhlükəsizliyi (thread-safe pipe),
yaddaş effektivliyi (stream).
**Çatışmamazlıqları:** Reader istifadə etsən Seek ola bilməz; böyük datada
buffersiz io.Copy təhlükəlidir.

### 2. bytes və strings paketləri
**Nədir:** String ↔ []byte çevirməsi, buffer yaratma və string üzərində
əməliyyatlar üçün köməkçi paketlər.

**Kitabdan kod nümunəsi:**
```go
// Buffer yaratmanın 3 yolu:
rawBytes := []byte(rawString)
var b = new(bytes.Buffer)
b.Write(rawBytes)                       // 1) boş buffer-ə yaz
b = bytes.NewBuffer(rawBytes)           // 2) baytlardan
b = bytes.NewBufferString(rawString)    // 3) stringdən birbaşa

// io.Reader-i tam stringə çevir:
func toString(r io.Reader) (string, error) {
    b, err := ioutil.ReadAll(r)
    return string(b), err
}

// Reader + Scanner ilə tokenləşdirmə:
reader := bytes.NewReader([]byte(rawString))
scanner := bufio.NewScanner(reader)
scanner.Split(bufio.ScanWords)    // söz-söz ayır
for scanner.Scan() {
    fmt.Print(scanner.Text())
}
```

**String funksiyaları:**
```go
strings.Contains(s, "this")      // alt string var mı
strings.ContainsAny(s, "abc")     // hərflərdən biri var mı
strings.HasPrefix(s, "this")     // başlanğıc
strings.HasSuffix(s, "test")     // sonluq
strings.Split(s, " ")            // hissələrə böl
strings.Title(s)                 // hər sözün baş hərfi böyük
strings.TrimSpace(s)             // kənar boşluqları sil
r := strings.NewReader(s)        // stringdən io.Reader
```

**Sub-kod izahı:**
- `bytes.Buffer` → io.Reader interfeysini realizə edir, striminq üçün flexibil
- `bufio.NewScanner` + `Split(ScanWords)` → axını tokenlərə böler
- `strings.NewReader` → stringi Reader kimi təqdim edir

### 3. Kataloqlar və fayllarla iş (os paketi)
**Nədir:** Fayl yaratma/açma/oxuma, kataloq yaratma, CRUD əməliyyatları —
platformalararası (Windows/Unix eyni API).

**Kitabdan kod nümunəsi:**
```go
func Operate() error {
    if err := os.Mkdir("example_dir", os.FileMode(0755)); err != nil {
        return err
    }
    if err := os.Chdir("example_dir"); err != nil {   // kataloqa keç
        return err
    }
    f, err := os.Create("test.txt")      // yarat + yazma üçün aç
    if err != nil {
        return err
    }
    value := []byte("hello\n")
    count, err := f.Write(value)
    if count != len(value) {
        return errors.New("incorrect length returned from write")
    }
    f.Close()
    f, err = os.Open("test.txt")         // oxu üçün aç
    io.Copy(os.Stdout, f)
    f.Close()
    os.Chdir("..")
    return os.RemoveAll("example_dir")   // təmizlə — DİQQƏTLİ!
}
```

**Faylı emal edib başqa fayla yaz (Capitalizer):**
```go
func Capitalizer(f1 *os.File, f2 *os.File) error {
    if _, err := f1.Seek(0, io.SeekStart); err != nil {
        return err
    }
    var tmp = new(bytes.Buffer)
    io.Copy(tmp, f1)                       // fayl → yaddaş
    s := strings.ToUpper(tmp.String())
    io.Copy(f2, strings.NewReader(s))      // string → fayl
    return nil
}
```

**Sub-kod izahı:**
- `os.FileMode(0755)` → Unix-style icazələr (Chown-a bənzər)
- `os.File` → Reader və Writer interfeyslərini realizə edir (açılış bitindən asılı)
- `os.RemoveAll` → təhlükəlidir — root/user input ilə ehtiyatlı!
- Fayl işləri stream kimi — buffer nümunələri ilə eyni interfeyslər

### 4. CSV formatı (encoding/csv)
**Nədir:** io.Reader/io.Writer üzərində işləyən CSV parser/writer.

**Necə işləyir:** Reader-dən oxuyur, Writer-ə yazır — mənbə/ismətləçər
istənilən ola bilər (fayl, buffer, socket).

**Kitabdan kod nümunəsi:**
```go
func ReadCSV(b io.Reader) ([]Movie, error) {
    r := csv.NewReader(b)
    r.Comma = ';'     // ayırıcı: nöqtə-vergül
    r.Comment = '-'   // şərh sətri prefiksi
    _, err := r.Read() // header-i oxu və nəzərə almayaraq keç
    if err != nil && err != io.EOF {
        return nil, err
    }
    var movies []Movie
    for {
        record, err := r.Read()
        if err == io.EOF {
            break
        } else if err != nil {
            return nil, err
        }
        year, err := strconv.ParseInt(record[2], 10, 64)
        m := Movie{record[0], record[1], int(year)}
        movies = append(movies, m)
    }
    return movies, nil
}

// Yazma tərəfi:
func (books *Books) ToCSV(w io.Writer) error {
    n := csv.NewWriter(w)
    n.Write([]string{"Author", "Title"})   // header
    for _, book := range *books {
        n.Write([]string{book.Author, book.Title})
    }
    n.Flush()           // buffer-i boşalt
    return n.Error()     // yığılmış xətanı qaytar
}
```

**Sub-kod izahı:**
- `r.Comma`/`r.Comment` → formatın konfiqurasiyası
- `n.Error()` → Flush-dan sonra yığılmış xətaları yoxlamaq üçün
- Record-record oxuma → az yaddaş, böyük fayllar üçün uyğun

**Praktik məqam:** CSV paketi interfeyslərin gücünü göstərir — bir sətirlə
mənbəyi/ismətləçəri dəyişmək olar; paralel emal (pipeline/worker pool) ilə
birləşdirilə bilər.

### 5. Müvəqqəti fayllar (ioutil)
**Nədir:** Ad toqquşması və təmizləmə problemi olmadan müvəqqəti fayl/kataloq.

**Kitabdan kod nümunəsi:**
```go
func WorkWithTemp() error {
    t, err := ioutil.TempDir("", "tmp")   // os.TempDir() yerində yaradır
    if err != nil {
        return err
    }
    defer os.RemoveAll(t)                  // funksiya çıxanda hamısı silinir
    tf, err := ioutil.TempFile(t, "tmp")   // temp kataloq İÇINDƏ temp fayl
    if err != nil {
        return err
    }
    fmt.Println(tf.Name())
    return nil
}
```

**Sub-kod izahı:**
- `ioutil.TempDir("", "tmp")` → unikal adlı kataloq; "" = sistem temp yeri
- `defer os.RemoveAll(t)` → kataloqu bütün içi ilə silir — fərdi silmə lazım deyil
- Test yazarkən müvəqqəti fayllar tövsiyə olunur

### 6. Şablonlar (text/template və html/template)
**Nədir:** Məlumatı mətnə çevirmək üçün şablon mühərriki — dəyişənlər,
şərtlər, tsikllər, bloklar, funksiyalar.

**Necə işləyir:** `{{ }}` mötərizələri içində məntiq; `.Field` ilə dataya
çıxış; FuncMap ilə xüsusi funksiyalar.

**Kitabdan kod nümunəsi:**
```go
const sampleTemplate = `
This template demonstrates printing a {{ .Variable | printf "%#v" }}.
{{if .Condition}}
If condition is set, we'll print this
{{else}}
Otherwise, we'll print this instead
{{end}}
{{range $index, $item := .Items}}
{{$index}}: {{$item}}
{{end}}
{{ range $index, $item := split .Words ","}}
{{$index}}: {{$item}}
{{end}}
{{ block "block_example" .}}
No Block defined!
{{end}}
{{/* multi-line comment */}}
`
// Funksiya qeydiyyatı:
funcmap := template.FuncMap{"split": strings.Split}
t := template.New("example").Funcs(funcmap)   // zəncirlənə bilər
t, err := t.Parse(sampleTemplate)
// Blok overriding — Clone + ikinci şablon parse:
t2, err := t.Clone()
t2, err = t2.Parse(secondTemplate)   // "block_example" yenidən təyin edilir
err = t2.Execute(os.Stdout, &data)
```

**Fayllardan şablon toplama (ParseGlob):**
```go
pattern := filepath.Join(tempdir, "*.tmpl")
tmpl, err := template.ParseGlob(pattern)   // bütün .tmpl fayllarını birləşdirir
tmpl.Execute(os.Stdout, map[string]string{  // map ilə də işləyir
    "Var1": "Var1!!", "Var2": "Var2!!", "Var3": "Var3!!",
})
```

**html/template — təhlükəsizlik:**
```go
// html/template avtomatik escape edir (JS injection-a qarşı):
err = t.Execute(os.Stdout, map[string]string{
    "Name": "<script>alert('Can you see me?')</script>",
})
// <script> NEUTRALIZƏ OLUNUR — context-aware escape
// Manual escaper-lər:
template.JSEscaper(`example <example@example.com>`)
template.HTMLEscaper(`example <example@example.com>`)
template.URLQueryEscaper(`example <example@example.com>`)
```

**Sub-kod izahı:**
- `{{if}}/{{else}}/{{end}}` → şərti render
- `{{range $index, $item := .Items}}` → tsikl
- `{{ block "ad" .}}` → standart blok; `{{ define "ad" }}` ilə override
- `t.Clone()` → əvvəlki parse-ləri saxlayaraq yeni blok əlavə etmək
- `template.Must(...)` → xətada panic (qısa başlanğıc üçün)
- html/template = text/template wrapper; veb üçün html, qalan üçün text

**Üstünlükləri:** Nesting, FuncMap, kontekst-aware escape (XSS qoruması).
**Çatışmamazlıqları:** Əvvəl qorxuducu görünə bilər; sintaksis məhdudluğu.

## Əsas terminlər

- io.Reader / io.Writer (oxuma/yazma interfeysləri)
- Stream (məlumat axını)
- io.MultiWriter (çoxlu yazıçı birləşdirməsi)
- io.Pipe (yaddaş konveyeri)
- Buffer (yaddaş tamponu)
- Tokenizer / Scanner (tokenləyici)
- Template (şablon)
- FuncMap (şablon funksiyaları xəritəsi)
- Context-aware escaping (kontekst-uyğun escape)
- Temporary File (müvəqqəti fayl)

## Praktik nəticə

- Funksiya imzalarında konkret tip YOX — io.Reader/io.Writer qəbul et
- Böyük axınlar üçün io.CopyBuffer; yenidən oxuma lazımdırsa ReadSeeker
- Fayl əməliyyatları stream kimi — os.File həm Reader həm Writer-dir
- CSV record-record emal edir — yaddaş effektivdir
- Temp fayllar testlər və artifact-lər üçün idealdır
- Veb render üçün mütləq html/template (XSS qoruması), plain mətn üçün text/template

## Mənbə

Pages: 23-57 (Chapter 1, Go Programming Cookbook 2nd ed)
