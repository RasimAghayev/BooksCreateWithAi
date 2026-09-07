# Chapter 10 — Sending and receiving data (Data göndərmə və qəbul)

## Bu chapter nədən bəhs edir?

Statik fayl servis (FileServer/ServeFile, subdirectory, custom error səhifələri, go:embed
binary daxil, CDN), advanced form (ParseForm/ParseMultipartForm, çoxdəyərli sahələr,
fayl upload — tək/çox), MIME type yoxlaması (3 üsul) və incremental upload (MultipartReader).

## Əsas fikirlər

### 1. Statik Fayl Servisi
**FileServer:** `http.Dir` + `http.FileServer` — If-Modified-Since görə 304 qaytarır:
```go
dir := http.Dir("./files")
http.ListenAndServe(":8080", http.FileServer(dir))
// və ya route ilə:
fileServer := http.FileServer(http.Dir("./static/"))
http.Handle("/static/", http.StripPrefix("/static", fileServer))
```
Qeydlər: mövcud olmayan qovluq ERROR VERMİR — sadəcə 404; yoxlama əlavə et. Quraşdırma
yolu dəyişə bilər — env variable + fallback.

**ServeFile:** Tək fayl handler daxilindən (`http.ServeFile(res, req, "./files/readme.txt")`)
— fayl adı user-dən gəlirsə sanitize; route+PathValue variantında FileServer daha təhlükəsiz.

### 2. Custom Error Səhifələri
Standard FileServer 404-ü brauzer formatında verir, dəyişməz (ServeContent → Error/NotFound
private funksiyaları). **Masterminds/go-fileserver** (kitab üçün yaradılıb) — standart
fork-u:
```go
fs.NotFoundHandler = func(w http.ResponseWriter, req *http.Request) {
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    fmt.Fprintln(w, "The requested page could not be found.")
}
http.ListenAndServe(":8080", fs.FileServer(http.Dir("./files")))
```

### 3. go:embed (Go 1.16+) — Binary Daxil Fayllar
```go
import "embed"

//go:embed files          // compiler directive — funksiya daxilində OLMAZ (pragma kimi)
var f embed.FS
http.ListenAndServe(":8088", http.FileServer(http.FS(f)))   // http.FS wrapper
```
Tək fayl variantı:
```go
import _ "embed"          // side-effect import

//go:embed files/example.html
var myString string       // string/[]byte birbaşa dəyər kimi
```
Fayda: tək binary paylanması — konteynerlər üçün ideal. Qeyd: böyük binar fayllar binary-ni
şişirdir — CDN nəzərdən keçir.

### 4. Alternativ Yer (CDN)
Aktivlər ayrı serve olunursa — yer konfiqurasiyadan gəlir (flag/env/etcd), şablonda
istifadə olunur:
```go
var l = flag.String("location", "http://localhost:8080", "A location.")
data := struct{ Location *string }{Location: l}
// tpl: <link rel="stylesheet" href="{{.Location}}/styles.css">
```
Qaydalar: hər environment öz nüsxəsi; default zərərlidir — dəyər yoxdursa log/panic;
global config obyektinə bağla. HTTP/2/3 push imkanları app+asset eyni server-dən tələb edir.

### 5. Form Əsasları (təkrar + dərinləşmə)
- `r.ParseForm()` — mətn sahələri; `r.ParseMultipartForm(maxMemory)` — fayl/binary
  (FormValue bunu lazım olanda avtomatik çağırır, default 32MB)
- **Form vs PostForm:** `r.Form` = URL query + POST body; `r.PostForm` = yalnız body.
  `FormValue`/`PostFormValue` = hər birinin İLK dəyəri.

**Çoxdəyərli sahə (checkbox qrupu kimi):**
```go
maxMemory := 16 << 20                        // 16MB; artığı diskə
r.ParseMultipartForm(maxMemory)
for k, v := range r.PostForm["names"] { ... }   // BÜTÜN dəyərlər
```
CSRF token forması daxil et (OWASP).

### 6. Fayl Upload — Tək Fayl
```html
<form action="/" method="POST" enctype="multipart/form-data">   <!-- məcburi enctype -->
  <input type="file" name="file" id="file">
```
```go
f, h, err := r.FormFile("file")          // (multipart.File, *multipart.FileHeader, err)
defer f.Close()
out, err := os.Create("/tmp/" + h.Filename)   // production: fayl store
defer out.Close()
io.Copy(out, f)
```

### 7. Çoxfayllı Upload (multiple atributu)
`<input type="file" name="files" multiple>` — FormFile yalnız İLKİNİ verir:
```go
r.ParseMultipartForm(16 << 20)
data := r.MultipartForm
files := data.File["files"]               // []*multipart.FileHeader
for _, fh := range files {
    f, err := fh.Open()
    defer f.Close()
    out, err := os.Create("/tmp/" + fh.Filename)
    defer out.Close()
    io.Copy(out, f)
}
```

### 8. MIME Type Yoxlaması — 3 Etibar Səviyyəsi
1. **Content-Type header** (upload edən tərəfindən qoyulur — az etibarlı):
   `header.Header["Content-Type"][0]` — `image/png` və ya `application/octet-stream`
2. **Fayl uzantısı** (user tərəfindən dəyişilə bilər):
   `filepath.Ext(header.Filename)` → `mime.TypeByExtension(ext)`
3. **Məzmun sniffing** (ən etibarlı, ən bahalı) — 512 bayt oxu:
```go
buffer := make([]byte, 512)
file.Read(buffer)
filetype := http.DetectContentType(buffer)
```
http.DetectContentType: HTML/text/XML/PDF/PostScript/şəkil/zip/RAR/gz/WAV/WebM.
Word/MP4 üçün — libmagic (Go binding-ləri mövcuddur; WHATWG mimesniff spec).
HTML accept atributu — bəzəkdür; yoxlama server-də MÜTLƏQ.

### 9. Incremental Upload ( böyük fayllar)
**Problem:** ParseMultipartForm temp diskdə yığıb gözləyir — böyük/paralel uploadlarda
server disk full. API pass-through modelində yaddaşda saxlama lazımsız.

**Həll — r.MultipartReader() (raw stream):**
```go
mr, err := r.MultipartReader()
values := make(map[string][]string)
maxValueBytes := int64(10 << 20)          // mətn sahələri üçün 10MB qorunma
for {
    part, err := mr.NextPart()
    if err == io.EOF { break }
    name := part.FormName()
    if name == "" { continue }
    filename := part.FileName()
    if filename == "" {                   // MƏTN sahəsi
        n, err := io.CopyN(&b, part, maxValueBytes)
        maxValueBytes -= n
        if maxValueBytes == 0 { return }   // hədd aşımı
        values[name] = append(values[name], b.String())
        continue
    }
    // FAYL sahəsi — gələn kimi yaz:
    dst, _ := os.Create("/tmp/" + filename)   // və ya cloud storage writer
    defer dst.Close()
    for {
        buffer := make([]byte, 100000)
        cBytes, err := part.Read(buffer)
        if err == io.EOF { break }
        dst.Write(buffer[0:cBytes])
    }
}
```
Fayldan/mətndən ayırma meyarı: FileName() boşdursa mətn. Fayl yazmaqla paralel —
upload bitməmiş diskə axır; hədəf bulud storage ola bilər.

## Əsas terminlələr
- FileServer / ServeFile / ServeContent — statik servis qatları
- http.Dir — fayl sistemi adapteri (FileSystem interfeysi)
- go:embed — compile-time fayl daxil etmə direktivi
- http.FS — embed.FS → http.FileSystem çevirici
- ParseForm / ParseMultipartForm — mətn / multipart parse
- Form vs PostForm — query+body / yalnız body
- multipart.FileHeader — upload metadata (Filename, Header)
- DetectContentType — 512-bayt MIME sniffing
- MultipartReader / NextPart — stream upload emalı

## Praktik nətidə

Statik data qərarları: (1) qovluq serve — FileServer + StripPrefix; (2) xüsusi 404 —
go-fileserver fork; (3) paylanma sadəliyi — go:embed (tək binary); böyük assetlər — CDN
(yer config-dən); (4) mətn forma — ParseForm; fayl — FormFile; çoxfayl — MultipartForm.File
slice; (5) tip yoxlaması — məzmun sniffing (512 bayt); header/extension yalnız müşayiətçi;
(6) böyük fayl / pass-through — MultipartReader + hissə-hissə yaz (temp yığılması YOX);
(7) bütün input sanitize + CSRF token; (8) upload fayl adlarını sanitize et (path traversal).

## Mənbə
Pages: 244-265 (PDF 265-286)
