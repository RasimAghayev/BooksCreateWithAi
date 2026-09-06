# Chapter 5 — Utilizing the Request Context and Data Bindings

## Bu chapter nədən bəhs edir?

Bu chapter Echo framework-unun **mərkəzi primitivi** olan `echo.Context` anlayışını dərinləşdirir: request-in emalı zamanı handler və middleware-lər arasında **state (vəziyyət) saxlama probleminin tarixi həlləri**, Echo-nun `sync.Pool` optimizasiyası, request payload-un avtomatik **Binding (Bağlama)** və **Validation (Doğrulama)** mexanizmləri, response-un hazır **rendering** helper-ləri və aşağı səviyyəli `echo.Response` strukturu izah olunur.

---

## Əsas fikirlər

### 1. Context probleminin tarixi (Pre Go 1.7 dövrü)

**Nəyə lazımdır?** Handler-lər arası çağırışlarda (xüsusən middleware zəncirində) yaradılmış state-i (məs., autentifikasiya nəticəsi, request ID) növbəti handler-ə ötürmək lazımdır. Standart kitabxanadakı `http.HandlerFunc` bunun üçün uyğun deyil.

```go
type HandlerFunc func(ResponseWriter, *Request) // http handler function type
type HandlerFunc func(Context) error            // Echo handler function type
```

**Problem:** `http.ResponseWriter` yalnız cavab baytlarını yazmaq üçün interfeysdir, `*http.Request` isə request-i təmsil edir — Go 1.7-yə qədər `http.Request` daxilində əlavə məlumat saxlamaq mexanizmi yox idi. Nəticədə handler-lər çox böyüyürdü, çünki bütün business logic-i özündə daşımalı idi.

**Üç əsas tarixi həll yanaşması:**

| Yanaşma | Nümunə | Üstünlük | Çatışmazlıq |
|---------|--------|----------|-------------|
| Qlobal map (Global mapping) | `gorilla/context` | Standart handler imzası qalır | Hər request üçün memory allocation, lock contention, memory leak təhlükəsi |
| Yeni handler imzası (New signature) | Echo, bir çox framework | Context birinci sinif vətəndaşdır | Vendor lock-in (handler-lər framework-ə bağlanır) |
| Request-içində gizlətmə (Hide in request) | Vestigo router (`husobee/vestigo`) | Standart imza qalır, imza dəyişmir | "Hack" xarakterli daşıyır, dinamik map sahələrindən sui-istifadə |

#### 1a. Qlobal context mapping (Gorilla context)

**Necə işləyir?** `gorilla/context` paketi tək qlobal map saxlayır:

```go
map[*http.Request]map[interface{}]interface{}
```

- `Set(request, key, value)` çağrıldıqda request pointer-i map key kimi istifadə olunur, dəyəri isə yeni yaranan `map[interface{}]interface{}` olur.
- `Get(request, key)` ilə oxunur.

**Niyə qeyri-optimaldır?**
1. Hər gələn request üçün qlobal map-ə yazma → minlərlə konkurent goroutine arasında **lock contention** (rəqabət).
2. Hər request üçün əlavə map allocation.
3. Explicit (açıq) təmizləmə edilməsə **memory leak** (yaddaş sızması) — dokumentasiyada da qeyd olunur.

#### 1b. Yeni handler funksiya tipi (Echo yanaşması)

**Necə işləyir?** Framework öz `Context` tipini yaradır və handler imzasını dəyişir: `func(Context) error`. Echo context-i həm request, həm response, həm də contextual data-yı bir yerdə birləşdirir (53 metod — "hər şey, hətta mətbəx lavabısı da").

**Vendor lock-in problemi:** Handler-lər digər framework-ə köçürülmək istəndikdə hamısı refaktor edilməlidir.

**Echo-nun `sync.Pool` optimizasiyası:** Echo hər request üçün yeni context yaratmaq əvəzinə `sync.Pool`-dan mövcud context instansiyasını götürüb **reuse** (yenidən istifadə) edir. Digər framework-lər hər request üçün yeni allocation edir — Echo bu sayəsində aşağı xərcli abstraksiyadır.

#### 1c. Request daxilində gizlətmə (Vestigo)

**Necə işləyir?** Vestigo URL router-i URL parameter adlarını və dəyərlərini request-in **form field**-ləri içərisində gizlədir. Developer standart handler imzası ilə URL parametrlərinə çıxış əldə edir — nə qlobal map, nə də dəyişdirilmiş imza lazımdır.

#### 1d. Go 1.7-dən sonra

Go 1.7 `http.Request`-ə `context.Context` əlavə etdi. Artıq framework seçmək məcburiyyəti yoxdur — context həmişə request-in özündədir. Bu context "bare bones" (sondar) olsa da cancellation (ləğvetmə) dəstəyi var.

---

### 2. Echo Context — əsas metodları

`echo.Context` request haqqında məlumat almaq və state manipulyasiya etmək üçün helper-lər toplusudur:

```go
// URL parametrləri
Param(name string) string          // URL-də :name ilə definiya edilmiş parametrin dəyəri
QueryParam(name string) string      // Query string parametri (?name=value)

// Form və fayl
FormValue(name string) string       // Form payload parametri
FormFile(name string) (*multipart.FileHeader, error) // Multipart fayl

// Cookie-lər
Cookie(name string) (*http.Cookie, error)  // Request-dən cookie oxu
SetCookie(cookie *http.Cookie)             // Set-Cookie response header-i əlavə et

// Context state (middleware ↔ handler körpüsü)
Get(key string) interface{}          // Middleware-in qoyduğu dəyəri oxu
Set(key string, val interface{})     // Növbəti middleware/handler üçün dəyər qoy

// Digər
Handler() HandlerFunc               // Router-in tapdığı handler funksiyası
Logger() Logger                     // Echo logger instansiyası
```

**Nəyə lazımdır?** `Set`/`Get` cütü middleware-dən handler-ə state ötürməyin əsas yoludur (məs., auth middleware user adını qoyur, handler `c.Get("username")` ilə oxuyur). `Param` isə `/users/:id` kimi rotalardan `id` almağın birinci yolu.

---

### 3. Request Binding (Request-in bağlanması)

**Nədir?** `Context.Bind(i interface{})` metodu HTTP request-in body-sini və `Content-Type` header-ini oxuyub göstərilən struktura avtomatik deserializə edir:

```go
// Bind binds the request body into provided type `i`. The default binder
// does it based on Content-Type header.
Bind(i interface{}) error
```

**Necə işləyir?** Default binder `Content-Type`-a baxaraq düzgün deserializasiya edir — dəstəklənən tiplər:
- `application/json`
- `application/xml`
- `application/x-www-form-urlencoded`

Login handler nümunəsi (`handlers/login.go`):

```go
func Login(c echo.Context) error {
    resp := renderings.LoginResponse{}
    lr := new(bindings.LoginRequest)
    if err := c.Bind(lr); err != nil {
        resp.Success = false
        resp.Message = "Unable to bind request for login"
        return c.JSON(http.StatusBadRequest, resp)
    }
    // ...
}
```

Binding strukturu (`bindings/login.go`):

```go
type LoginRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
}
```

Test üçün curl əmrləri (hər ikisi eyni handler-ə işləyir):

```bash
curl -XPOST -H"Content-Type: application/json" localhost:8080/login \
  -d'{"username":"test","password":"test"}'

curl -XPOST -H"Content-Type: application/xml" localhost:8080/login \
  -d'<LoginRequest><Username>test</Username>...</LoginRequest>'
```

**Üstünlükləri:**
1. Developer deserializasiya kodu yazmır — Echo yerinə yetirir.
2. **Content negotiation (Məzmun danışıqları):** Eyni API JSON və XML dəstəkləyir → daha çox inteqrasiya imkanı.
3. Handler təmiz qalır — yalnız business logic.

---

### 4. Binding Validation (Doğrulama)

**Nədir?** Binding-dən sonra input-un etibarlılığını yoxlamaq üçün mərkəzləşdirilmiş sxem. `Context.Validate(i)` strukturun `Validate()` metodunu çağırır.

**Necə qurulur?** 3 addım:

**Addım 1 — Validator interfeysi** (`bindings/common.go`):

```go
type Validatable interface {
    Validate() error
}

var ErrNotValidatable = errors.New("Type is not validatable")

type Validator struct{}

func (v *Validator) Validate(i interface{}) error {
    if validatable, ok := i.(Validatable); ok {
        return validatable.Validate()
    }
    return ErrNotValidatable
}
```

**Addım 2 — Hər binding öz Validate metodunu daşıyır** (`bindings/login.go`):

```go
func (lr *LoginRequest) Validate() error {
    errs := new(RequestErrors)
    if lr.Username == "" {
        errs.Append(ErrUsernameEmpty)
    }
    if lr.Password == "" {
        errs.Append(ErrPasswordEmpty)
    }
    if errs.Len() == 0 {
        return nil
    }
    return errs
}
```

**Addım 3 — Echo-ya Validator qeydiyyatdan keçirilir** (`cmd/service/main.go`):

```go
func main() {
    e := echo.New()
    e.Logger.SetLevel(log.INFO)
    e.Validator = new(bindings.Validator)
    // ...
}
```

Handler-də istifadə — Bind-dən dərhal sonra:

```go
if err := c.Bind(lr); err != nil { /* 400 */ }
if err := c.Validate(lr); err != nil {
    resp.Success = false
    resp.Message = err.Error()
    return c.JSON(http.StatusBadRequest, resp)
}
```

**Nəyə lazımdır?** Validation kodu strukturun yanında yaşayır → hər input strukturu öz qaydalarını özündə daşıyır, handler-lər validation "minutiyası" ilə clutter (daşqalaq) olmur.

---

### 5. Response Rendering (Cavabın render edilməsi)

`echo.Context`-də hazır cavab helper-ləri — developer-in özü serializasiya etməsinə ehtiyac yoxdur:

| Metod | Təyinat |
|-------|---------|
| `HTML(code, html)` / `HTMLBlob(code, b)` | HTML cavab |
| `String(code, s)` | Düz mətn |
| `JSON(code, i)` / `JSONPretty(code, i, indent)` / `JSONBlob(code, b)` | JSON (adi / gözəl formatlı / hazır blob) |
| `JSONP(code, callback, i)` / `JSONPBlob(...)` | JSONP |
| `XML(code, i)` / `XMLPretty(...)` / `XMLBlob(...)` | XML |
| `Blob(code, contentType, b)` | İstənilən content-type ilə xam baytlar |
| `Stream(code, contentType, r io.Reader)` | `io.Reader`-dən axın cavabı |
| `File(file)` | Fayl məzmununu serve et |
| `Attachment(file, name)` | Yükləmə kimi göndər |
| `Inline(file, name)` | Inline göstər |
| `NoContent(code)` | Boş body |
| `Redirect(code, url)` | Yönləndirmə |

**Necə işləyir (JSON nümunəsi)?** `c.JSON(http.StatusOK, resp)` çağrılanda:
1. `code` → `ResponseWriter`-a status kimi yazılır.
2. Status və header-lər yazıldıqdan sonra struktura JSON serializasiya edilir.
3. Nəticə `ResponseWriter.Write` ilə göndərilir.

```go
return c.JSON(http.StatusOK, resp)
```

---

### 6. echo.Response — aşağı səviyyə cavab manipulyasiyası

**Nədir?** `Context.Response()` çağıraraq `echo.Response` strukturuna çıxış əldə edirsiniz — əslində `http.ResponseWriter`-ın wrapper-i:

```go
type (
    Response struct {
        echo        *Echo
        beforeFuncs []func()
        afterFuncs  []func()
        Writer      http.ResponseWriter
        Status      int
        Size        int64
        Committed   bool
    }
)
```

**Əsas sahələr:**
- `Writer` — əsl `http.ResponseWriter`
- `Status` — yazılan status kodu
- `Size` — yazılan bayt sayı
- `Committed` — cavabın artıq göndərilib-göndərilmədiyi flag-i

**Header əlavə etmək:**

```go
ctx.Response().Header().Add("X-My-Header", "ValueOfHeader")
```

**`WriteHeader` wrapper-i — Committed qoruması:**

```go
func (r *Response) WriteHeader(code int) {
    if r.Committed {
        r.echo.Logger.Warn("response already committed")
        return
    }
    for _, fn := range r.beforeFuncs {
        fn()
    }
    r.Status = code
    r.Writer.WriteHeader(code)
    r.Committed = true
}
```

- Cavab iki dəfə yazılmasının qarşısını alır — `Committed == true` olduqda yalnız warn log qeyd olunur.
- Yazılmadan əvvəl bütün `beforeFuncs` icra olunur.

**`Write` wrapper-i:**

```go
func (r *Response) Write(b []byte) (n int, err error) {
    if !r.Committed {
        r.WriteHeader(http.StatusOK) // implicit WriteHeader
    }
    n, err = r.Writer.Write(b)
    r.Size += int64(n)
    for _, fn := range r.afterFuncs {
        fn()
    }
    return
}
```

- `WriteHeader` açıq çağrılmayıbsa ilk `Write` avtomatik `WriteHeader(200)` edir.
- Yazıdan sonra `afterFuncs` icra olunur, `Size` sayğacı artır.

**Before/After hook-ləri:** `Response.Before(fn)` və `Response.After(fn)` ilə cavabdan əvvəl/sonra işləyəcək funksiyalar qeydiyyata alınır — məs., cavab bitdikdən sonra state cleanup (təmizlənmə) aparmaq üçün.

---

## Əsas terminlər

| Termin | İzah |
|--------|------|
| Request Context (Sorğu konteksti) | Request-in emalı boyu middleware/handler-lər arasında ötürülmüş state konteyneri |
| Binding (Bağlama) | Request body-nin Content-Type-a görə struktura avtomatik deserializasiyası |
| Validation (Doğrulama) | Bound input-un qaydalara uyğunluğunun yoxlanması |
| `echo.Context` | Echo handler-lərinin tək parametri; request/response + state + helper-lər (53 metod) |
| `sync.Pool` | Müvəqqəti instansların reuse edilə bilən hovuzu — Echo context allocation xərcini sıfıra endirir |
| Global context mapping | Request pointer-inin qlobal map key kimi istifadəsi (gorilla/context) |
| Vendor lock-in | Kodun yalnız bir framework-ə bağlı olması, daşınmazlığı |
| Content negotiation | API-nın JSON/XML/form kimi birdən çox wire format dəstəkləməsi |
| `Committed` flag | Response-un artıq göndərildiyini bildirən bayraq — double-write qoruması |
| `beforeFuncs`/`afterFuncs` | Response yazılışından əvvəl/sonra icra olunan hook-lər |
| Memory leak (Yaddaş sızıntısı) | Təmizlənməyən qlobal map-lərdə yığılan istifadə olunmayan yaddaş |
| Lock contention | Çoxlu goroutine-in eyni map-ə yazma rəqabəti, performance azalması |
| Validatable interface | `Validate() error` metodu daşıyan strukturlar üçün universal abstraksiya |
| `echo.Response` | `http.ResponseWriter`-ın Status/Size/Committed/hook-lərlə zənginləşdirilmiş wrapper-i |

---

## Praktik nəticə

1. **Echo context = framework-un mərkəzi primitivi:** Handler imzası `func(echo.Context) error` olduğundan request/response-ə çıxış, state keçirmə və render etmə hamısı bir nöqtədən idarə olunur; `sync.Pool` sayəsində bu abstraksiya praktiki olaraq pulsuzdur.
2. **Bind + Validate cütü handler-i təmizləyir:** `c.Bind(&x)` → `c.Validate(&x)` ardıcıllığı deserializasiya və doğrulamanı strukturların yanında saxlayır; `e.Validator = new(bindings.Validator)` bir dəfə qurulur və bütün binding-lər avtomatik Validatable olur.
3. **Content-Type agnostik API:** Bind default olaraq JSON/XML/form-urlencoded fərq edir → eyni endpoint birdən çox klient formatına xidmət edir.
4. **Response render helper-ləri serializasiyanı öz üzərinə götürür:** `c.JSON/String/XML/File/Stream/Redirect...` — developer yalnız status kodu + datanı verir.
5. **Aşağı səviyyə nəzarət lazım olsa:** `c.Response()` → `Header().Add()`, `WriteHeader`, `Write`, `Before/After` hook-ləri — `Committed` qoruması ilə təhlükəsizdir.
6. **Köhnə qlobal map həllərindən (gorilla/context) qaçın:** Go 1.7+ request-də context dəstəkləyir; Echo `Set/Get` isə lock-free və pool-lu variantdır.

---

## Mənbə

- Kitab: *Echo Quick Start Guide* — J. Ben Huson, Packt Publishing, 2018 (ISBN 9781789340849)
- Chapter 5: "Utilizing the Request Context and Data Bindings", book səh. 89–113
- PDF səhifələri: 106–125
- Kod: https://github.com/PacktPublishing/Echo-Essentials/tree/master/chapter5
- Video: https://goo.gl/3gDXrq
