# Chapter 7 — Serving and receiving assets and forms (Technique 39-48)

## Bu chapter nədən bəhs edir?

Static fayl servisi (FileServer, StripPrefix, custom error pages, memory cache, binary embed, CDN), HTML form emalı (ParseForm/ParseMultipartForm, çoxdəyərli sahələr), fayl yükləmə (tək/çox fayl, MIME type yoxlaması, incremental stream upload).

## Əsas fikirlər

### 1. Static fayl servisi əsasları
Go app daxili web serverdir — Apache/Nginx ARXASINDA DURMAQ ZƏRURİ DEYİL. CGI/FastCGI paketləri mövcud amma tövsiyə olunmur (hər request üçün yeni proses).

```go
dir := http.Dir("./files")
http.ListenAndServe(":8080", http.FileServer(dir))     // http.Dir → FileSystem

// Və ya tək fayl:
http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
    http.ServeFile(res, req, "./files/readme.txt")
})
```
- **FileServer:** ağıllı — `If-Modified-Since` → **304 Not Modified** cavabı
- **ServeFile:** tək fayl; eyni 304 dəstəyi

### TECHNIQUE 39: Subdirectory + StripPrefix
```go
dir := http.Dir("./files/")
handler := http.StripPrefix("/static/", http.FileServer(dir))   // /static/XX → files/XX
http.Handle("/static/", handler)
http.HandleFunc("/", homePage)
http.ListenAndServe(":8080", nil)
```
- **StripPrefix** — URL-dən `/static/` prefiksini kəsir, yoxsa fayl axtarışı `files/static/...` edərdi

**Gotcha 1 — path resolution:** `pathResolver` (Ch 2, path.Match) istifadə etsən `*` bir directory səviyyəsində dayanır — subdirectory-lərdəki fayllar SERV OLUNMAZ. Regex router (T8) istifadə et.

**Gotcha 2 — error pages:** FileServer/ServeFile-in error cavabları **dəyişdirilə bilməz** (plain text).

### TECHNIQUE 40: Custom error pages
FileServer-in daxilində http.Error/NotFound **baked-in** — dəyişmək üçün fork lazımdır.

```go
fs "github.com/Masterminds/go-fileserver"   // kitab üçün yazılmış fork

fs.NotFoundHandler = func(w http.ResponseWriter, req *http.Request) {
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    fmt.Fprintln(w, "The requested page could not be found.")   // istənilən HTML
}
dir := http.Dir("./files")
http.ListenAndServe(":8080", fs.FileServer(dir))
```

### TECHNIQUE 41: Memory cache + ServeContent
**Kitabdan kod nümunəsi:**
```go
type cacheFile struct {
    content io.ReadSeeker         // Read+Seek — ServeContent üçün mütləq
    modTime time.Time             // If-Modified-Since üçün
}

var cache map[string]*cacheFile
var mutex = new(sync.RWMutex)     // RWMutex: çox oxu + tək yazı

func serveFiles(res http.ResponseWriter, req *http.Request) {
    mutex.RLock()                  // OXU kilidi — paralel oxuyanlar OK
    v, found := cache[req.URL.Path]
    mutex.RUnlock()
    if !found {
        mutex.Lock()               // YAZI kilidi — bir yazıçı
        defer mutex.Unlock()
        fileName := "./files" + req.URL.Path
        f, err := os.Open(fileName)
        defer f.Close()
        if err != nil {
            http.NotFound(res, req)
            return
        }
        var b bytes.Buffer
        io.Copy(&b, f)             // fayl → buffer
        r := bytes.NewReader(b.Bytes())   // → ReadSeeker
        info, _ := f.Stat()
        v := &cacheFile{content: r, modTime: info.ModTime()}
        cache[req.URL.Path] = v
    }
    http.ServeContent(res, req, req.URL.Path, v.modTime, v.content)
}
```

**Sub-kod izahı:**
- **RWMutex:** `RLock/RUnlock` — çoxsaylı oxucular paralel; `Lock/Unlock` — tək yazıçı, hamıyı blokləyır
- `bytes.NewReader` — Read+Seek interfeysləri (ServeContent tələbi)
- **ServeContent:** MIME type təyinatı, 304 Not Modified, Content-Length — hamısı AVTOMATİK
- Xatırlatma: Varnish (reverse proxy) çox vaxt daha yaxşı; in-memory cache — özəl hal; **memory monitoring** vacibdir (OOM qorxusu). groupcache — serverlərarası shared cache.

### TECHNIQUE 42: Faylları binary-yə embed etmək
**Məqsəd:** 1 binary distribusiya — fayllar yox.

```go
box := rice.MustFindBox("../files/")
httpbox := box.HTTPBox()                        // http.FileSystem interfeysi
http.ListenAndServe(":8080", http.FileServer(httpbox))
```
- **go.rice** (github.com/GeertJohan/go.rice): `go run` → fayl sistemindən; binary → daxilə kömürülüb
- Build: `rice embed-go` (faylları Go fayllarına çevirir) + `go build`
- Template-lər də mümkün: `templateString, err := box.String("example.html")`
- Minification (CSS/JS whitespace təmizliyi) binary ölçüsünü azaldır
- Xatırlatma: `os.Walk` symlink-ləri gəzmir

### TECHNIQUE 43: Alternative location (CDN)
- Hər mühitdə (dev/test/prod) ÖZ asset nüsxəsi
- Lokasiya → konfiqurasiya (flag/env/file/etcd) → template-ə ötür

```go
var l = flag.String("location", "http://localhost:8080", "A location.")
var tpl = `... <link rel="stylesheet" href="{{.Location}}/styles.css"> ...`

func servePage(res http.ResponseWriter, req *http.Request) {
    data := struct{ Location *string }{Location: l}
    t.Execute(res, data)
}
```
- Productionda: defolt dəyər yoxdursa panic et — yanlış URL-yə işarə etmə
- **HTTP/2 (RFC 7540):** server push — səhifə + asset-lər eyni connection, sorğu olmadan əvvəl.

### 2. Form əsasları
```go
name := r.FormValue("name")     // avtomatik parse + İLK dəyər

// Və ya açıq şəkildə:
err := r.ParseForm()              // yalnız TEXT sahələr
if err != nil { fmt.Println(err) }
name := r.FormValue("name")
```

**Parse lokasiyaları:**
- `r.Form` — URL query + POST/PUT body; hər açar = DƏYƏR SLICE
- `r.PostForm` — yalnız POST/PUT body (query YOX)
- `FormValue` → Form-un ilkini; `PostFormValue` → PostForm-un ilkini

### TECHNIQUE 44: Çoxdəyərli sahələr (checkboxes)
```go
maxMemory := 16 << 20                      // 16 MB — fayl hissələri üçün (artığı diskə)
err := r.ParseMultipartForm(maxMemory)     // FormValue default 32 MB istifadə edir
if err != nil { fmt.Println(err) }

for k, v := range r.PostForm["names"] {     // TAM slice-i iterasiya et
    fmt.Println(v)
}
```
- FormValue yalnız birincini verir — **slice-ı birbaşa oxu**
- CSRF token tövsiyəsi: en.wikipedia.org/wiki/Cross-site_request_forgery

### TECHNIQUE 45: Tək fayl upload
```html
<form action="/" method="POST" enctype="multipart/form-data">   <!-- MÜTLƏQ multipart -->
  <input type="file" name="file" id="file">
  <button type="submit" name="submit">Submit</button>
</form>
```

```go
func fileForm(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
        t, _ := template.ParseFiles("file.html")
        t.Execute(w, nil)
    } else {
        f, h, err := r.FormFile("file")      // (multipart.File, *FileHeader, err)
        if err != nil {
            panic(err)
        }
        defer f.Close()
        filename := "/tmp/" + h.Filename      // productionda file store lazımdır
        out, err := os.Create(filename)
        if err != nil {
            panic(err)
        }
        defer out.Close()
        io.Copy(out, f)
        fmt.Fprint(w, "Upload complete")
    }
}
```
- **r.FormFile("ad")** — parse (lazımsa) + fayl obyekti + **FileHeader** (Filename, Header metadata)
- `h.Filename` — orijinal fayl adı

### TECHNIQUE 46: Çox fayl upload (`multiple`)
```html
<input type="file" name="files" id="files" multiple>
```

```go
err := r.ParseMultipartForm(16 << 20)
if err != nil { ... }
data := r.MultipartForm
files := data.File["files"]        // []*multipart.FileHeader slice!
for _, fh := range files {
    f, err := fh.Open()             // FileHeader → fayl obyekti
    defer f.Close()
    out, err := os.Create("/tmp/" + fh.Filename)
    defer out.Close()
    io.Copy(out, f)
}
```
- `r.MultipartForm.File["ad"]` — FileHeader slice; **FormFile yalnız birincini verir**
- Hər FileHeader → `.Open()` → fayl → io.Copy

### TECHNIQUE 47: Fayl tipi yoxlaması
**3 üsul — etibar azalan sıra ilə:**

**a) Content yoxlama (ƏN ETİBARLI):**
```go
buffer := make([]byte, 512)
_, err = file.Read(buffer)                 // ilk 512 bayt
filetype := http.DetectContentType(buffer)  // daxili magic sniffing
// HTML, text, XML, PDF, PS, images, RAR/Zip/GZip, WAV, WebM
// tanınmayanda: application/octet-stream
```

**b) Header-dən:**
```go
file, header, err := r.FormFile("file")
contentType := header.Header["Content-Type"][0]   // header-lər multi-value — [0]
// client tərəfindən qoyulur — ETİBARSIZ
```

**c) Uzantıdan:**
```go
extension := filepath.Ext(header.Filename)
type := mime.TypeByExtension(extension)     // uzantı dəyişdirilə bilər
```
- **accept** HTML atributu — brauzer dəstəyi qeyri-mümkün, client manipulyasiyası — yalnız UX
- Geniş format lazımdırsa → **libmagic** binding-ləri
- MIME sniffing spec: mimesniff.spec.whatwg.org

### TECHNIQUE 48: Incremental save (streaming upload)
**Problem:** böyük fayllar API serverdən KEÇİR (proxy/storage) — ParseMultipartForm hər şeyi tmp-də yığır.

**Həll:** `r.MultipartReader()` — raw stream; upload VAXTINDA emal et.

```go
func fileForm(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
        t, _ := template.ParseFiles("file_plus.html")
        t.Execute(w, nil)
    } else {
        mr, err := r.MultipartReader()      // *multipart.Reader!
        if err != nil {
            panic("Failed to read multipart message")
        }

        values := make(map[string][]string)          // text sahələri burada
        maxValueBytes := int64(10 << 20)             // text-lər üçün 10 MB limit

        for {
            part, err := mr.NextPart()               // növbəti hissə
            if err == io.EOF {
                break                                // bitdi
            }
            name := part.FormName()
            if name == "" {
                continue
            }
            filename := part.FileName()
            var b bytes.Buffer

            if filename == "" {                       // TEXT sahəsi
                n, err := io.CopyN(&b, part, maxValueBytes)
                if err != nil && err != io.EOF {
                    fmt.Fprint(w, "Error processing form")
                    return
                }
                maxValueBytes -= n
                if maxValueBytes == 0 {
                    fmt.Fprint(w, "multipart message too large")
                    return
                }
                values[name] = append(values[name], b.String())
                continue
            }

            // FAYL hissəsi — incrementally yaz:
            dst, err := os.Create("/tmp/" + filename)
            defer dst.Close()
            if err != nil {
                return
            }
            for {
                buffer := make([]byte, 100000)         // 100KB chunk!
                cBytes, err := part.Read(buffer)
                if err == io.EOF {
                    break
                }
                dst.Write(buffer[0:cBytes])            // gələn kimi diskə
            }
        }
        fmt.Fprint(w, "Upload complete")
    }
}
```

**Sub-kod izahı:**
- **Handler request BAŞLAYANDA çağırılır** (bitməyib!) — upload zamanı emal mümkündür
- `mr.NextPart()` — növbəti multipart hissə; io.EOF → son
- `part.FileName() == ""` → text sahəsi; dolu → fayl
- Fayl: 100KB chunk-larla oxu → diskə yaz — **bütün fayl yaddaşda heç vaxt olmur**
- `io.CopyN(&b, part, maxValueBytes)` — text həcmi məhdudlaşdırma
- Pass-through model: yaddaş/disk cache YOX — direkt hədəfə

## Fayl/form emalı qərar cədvəli
| Hal | Alət |
|---|---|
| Sadə text form | ParseForm + FormValue |
| Çoxdəyərli text | ParseMultipartForm + r.PostForm["ad"] slice |
| Tək fayl | r.FormFile("ad") |
| Çox fayl | ParseMultipartForm + MultipartForm.File["ad"] |
| Böyük fayl / pass-through | MultipartReader + NextPart + chunk Read |

## Static serving qərar cədvəli
| Hal | Alət |
|---|---|
| Sadə qovluq | FileServer + http.Dir |
| Subpath | StripPrefix |
| Custom 404 | Masterminds/go-fileserver fork |
| Sürət | Memory cache + ServeContent (Varnish alternativi) |
| Tək binary | go.rice embed |
| CDN / ayrı lokasiya | Konfiq + template {{.Location}} |

## Əsas terminlər
- http.FileServer / http.Dir / http.ServeFile
- If-Modified-Since / 304 Not Modified
- http.StripPrefix
- http.ServeContent (ReadSeeker + modTime; MIME/Length avtomatik)
- RWMutex (RLock/RUnlock — çox oxu; Lock/Unlock — tək yazı)
- go-fileserver (custom error pages)
- go.rice (binary embed; rice embed-go)
- rice.MustFindBox / HTTPBox / http.FileSystem
- CDN / mühit-əsaslı asset ayrılığı
- HTTP/2 Server Push (RFC 7540)
- ParseForm / ParseMultipartForm(maxMemory) / FormValue / PostFormValue
- r.Form / r.PostForm / r.MultipartForm.File
- multipart.File / *multipart.FileHeader (Filename/Header/Open)
- DetectContentType (512 bayt sniffing)
- mime.TypeByExtension / filepath.Ext
- libmagic
- multipart.Reader / MultipartReader / NextPart / FormName / FileName
- Chunked upload (streaming, incremental save)
- Pass-through / API proxy model
- CSRF Token

## Praktik nəticə
- Static fayl üçün `FileServer + http.Dir + StripPrefix` başlanğıc nöqtəsidir; custom 404 lazımdırsa fork.
- ParseMultipartForm faylları tmp-də saxlayır — böyük fayl/pass-through üçün MultipartReader + chunk Read (yaddaş sabit qalır).
- FormValue yalnız İLK dəyəri verir — çoxdəyərli hallarda Form/PostForm slice-ı birbaşa oxu.
- RWMutex oxu-yazı ssenarisində: oxu çox, yazı nadir → RLock paralelizm qazandırır.
- MIME type: extension/header ETKİBARSIZ — yüklənən faylın ÖZ baytlarını yoxla (DetectContentType / libmagic).
- go.rice ilə tək binary + asset; minification ilə ölçünü azalt.
- Asset lokasiyasını konfiqurasiyadan gətir, template-ə ötür — dev/test/prod fərqli CDN-lər.

## Mənbə
Pages: 191-213 (PDF), book pages 168-193
