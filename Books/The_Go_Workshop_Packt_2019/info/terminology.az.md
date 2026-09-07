# The Go Workshop — Terminologiya (Azərbaycanca)

Format: **English Term (Azərbaycanca qarşılıq)** — qısa izah

## Dəyişənlər və Tiqlər (Ch1-3)
- **Short Variable Declaration (qısa təyinat)** — `:=` — tip çıxarışı ilə qısa bəyan
- **Zero Value (sıfır dəyəri)** — tipin default dəyəri (0/false/""/nil)
- **Shadowing (kölgələmə)** — daxili scope-da eyni adlı YENİ dəyişən
- **Escape Analysis (çıxış analizi)** — compiler-in stack/heap qərarı
- **Wraparound (sarğırma)** — tip həddi aşınca min dəyərə qayıtma
- **math/big** — int64-dən böyük ədədlər üçün paket
- **Rune (simvol kodu)** — int32 əsaslı tək Unicode simvol
- **Raw/Interpreted Literal** — backquote (escape-siz) / dırnaq (escape-li)
- **len(s) = bayt, len([]rune(s)) = simvol** — multi-byte fərq

## Kolleksiyalar (Ch4)
- **Slice Header** — pointer + len + cap üçlüyü
- **Capacity (tutum)** — hidden array ölçüsü; append həddi
- **`append(s1[:0:0], s1...)`** — ən effektiv müstəqil kopya
- **Comma-Ok (`v, ok := m[k]`)** — map açarının mövcudluğu
- **delete(m, k)** — map elementinin silinməsi
- **Type Assertion (`v.(T)`)** — interface{}-dən concrete tip
- **Type Switch (`switch v := x.(type)`)** — çoxtipli yoxlama

## Funksiya və Xətalar (Ch5-6)
- **Named Return (adlı qaytarma)** — imzada adlı nəticə; naked return tələsi
- **Pack/Unpack (`...T` / `s...`)** — variadic parametr / slice açma
- **Closure (bağlanma)** — xarici scope dəyişənini saxlayan anonim funksiya
- **FILO defer** — son daxil, ilk icra; arqumentlər DEFER ANINDA bağlanır
- **Err Prefiksi** — xəta dəyişənlərinin idiomatik adı
- **panic/recover** — çökdürmə / YALNIZ defer daxilində bərpa
- **os.Exit vs panic** — defer-ləri öldürür / işlədir

## İnterfeyslər (Ch7)
- **Implicit Satisfaction (bəyansız satisfaksiya)** — implements YOXDUR
- **Duck Typing (ördək yazısı)** — metodlara görə tip uyğunluğu
- **Stringer** — String() metodu ilə fmt çapının fərdiləşdirilməsi
- **Accept Interfaces, Return Structs** — API dizayn prinsipi
- **io.Reader/io.Writer** — universal oxuma/yazma interfeysləri

## Paketlər (Ch8)
- **Exported/Unexported** — BÖYÜK/kiçik hərf görünürlüyü
- **Blank Import (`_ "pkg"`)** — yalnız side-effect (driver qeydiyyatı)
- **Init Order** — import → variables → init → main
- **GOPATH src/bin/pkg** — mənbə / binary / object qovluqları

## Vaxt və JSON (Ch10-11)
- **time.Duration** — vaxt fərqi tipi; 6 rezolyusiya
- **Parse/Format** — string↔time; RFC3339/ANSIC sabitləri
- **Referens Time (`01/02 03:04:05`)** — custom format şablonu
- **Struct Tag (`json:"ad,omitempty"`)** — JSON açar nəzarəti
- **`json:"-"`** — sahə tamamilə gizli (parol/token)
- **map[string]interface{}** — naməlum JSON strukturun Go forması
- **GOB** — Go-only binary protokolu; model uyğunsuzluq tolerantlığı

## Fayllar və Sistem (Ch12-13)
- **O_APPEND|O_CREATE|O_WRONLY** — append üçün əs üçlük
- **os.IsNotExist** — Stat xətasının TİPLƏ yoxlanması
- **signal.Notify** — SIGINT/SIGTSTP tutma; graceful cleanup
- **Prepare + $1/$2** — SQL injection qorunması
- **rows.Next + Scan** — nəticə iterasiyası (& ilə doldurma)
- **RowsAffected** — əməliyyat təsirinin təsdiqi

## Web (Ch14-15)
- **http.Get/Post** — sadə client qısa yolları
- **Multipart Form (CreateFormFile)** — fayl upload formatı
- **FormDataContentType** — multipart başlığı MÜTLƏQ
- **http.Client{Timeout}** — custom klient + vaxt həddi
- **Handler (ServeHTTP) vs HandleFunc** — state / sadə funksiya
- **`{{if}}...{{end}}` / `{{.Field}}`** — template şərti / struct istinadı
- **FileServer + StripPrefix** — statik qovluq xidməti + maskalama
- **r.ParseForm / r.Form.Get** — POST form oxunuşu

## Concurrency (Ch16)
- **WaitGroup (Add/Done/Wait)** — sayğaclı goroutine sinxronizasiyası
- **Race Condition / `-race`** — data yarışı / detector flag
- **sync/atomic** — kilidsiz atomik əməliyyatlar
- **Mutex Lock/Unlock** — kritik bölmə; MINIMAL saxla
- **Buffered/Unbuffered Channel** — tutumlu (non-blocking) / sinxron
- **Deadlock** — hamı gözləyir; unbuffered tək rutində
- **close(ch) + `for range ch`** — kanal sonu; worker pool mexanizmi
- **Fan-out/Fan-in** — N worker-ə paylama / birləşdirmə
- **Done Channel** — bitmə siqnalı (WaitGroup alternativi)
- **context.WithCancel/Done** — rutin ləğvi (HTTP timeout standartı)
- **"Share by communicating"** — paylaşmaq üçün kanal, mutex yalnız lazımda

## Alətlər və Təhlükəsizlik (Ch17-19)
- **gofmt/goimports (-w)** — format / import idarəsi
- **go vet** — statik analiz; Printf args, unmarshal pointer
- **go doc -all** — kommentlərdən sənəd
- **SQL Injection / Prepare** — concat hücumu / placeholder qorunması
- **XSS / html-template** — script inyeksiyası / avtomatik escape
- **Nonce** — tək-istifadə random; AES-GCM Seal-in dst-sində
- **EncryptOAEP/DecryptOAEP** — RSA public/private əməliyyatları
- **crypto/rand vs math/rand** — security / simulyasiya randomu
- **x509/PEM/CA** — sertifikat / ASCII format / imzalayan orqan
- **bcrypt** — parol hash; GenerateFromPassword/CompareHashAndPassword
- **Build Tag / Filename Suffix** — `// +build linux` / `_GOOS_GOARCH`
- **GOOS/GOARCH** — cross-compile dəyişənləri
- **reflect.TypeOf/ValueOf/DeepEqual** — runtime inspekt / dərin müqayisə
- **unsafe.Pointer / cgo** — bit çevirmə / C FFI; C.free MÜTLƏQ
