# Chapter 2 — Scaffolding a Go Microservice (səh. 13-54)

## Bu fəsil nədən bəhs edir?

Idiomatik Go kodu yazma qaydaları, proyekt strukturu standartları və kitabın
boyunca istifadə ediləcək **Movie aplikasiyasının** (metadata, rating, movie — 3
mikroservis) sıfırdan scaffolding-i.

## Əsas fikirlər

### 1. Core principles (əsas prinsiplər)
- Rəsmi qidlara uyğun ol: **Effective Go**, **CodeReviewComments**
- Standart kitabxananın (context, net) stilini izlə
- Başqa dillərin ideyalarını Go-ya MƏCBURİƏT etmə — Go fəlsəfəsini başa düş

### 2. İdiomatik Go qaydaları
**Adlandırma (Naming):**
- Export olunan adlar böyük hərflə; istinadda paket adı daxildir (bytes.Buffer)
- Paket adı ilə prefiks YOX (xml.Reader, xml.XMLReader deyil)
- Getter-lərdə Get prefiksi YOX (`Age()` düz, `GetAge()` səhv), Set olar
- Tək-metodlu interfeys: metod + er → `Writer`
- Akronimlər ardıcıl: `URL`, `ID` (Url, Id YOX); qısa dəyişən adları (i, loop)

**Şərhlər (Comments):** hər paket və export olunan ad üçün tam cümlə, nöqtə ilə
bitən; ilk cümlə export olunan adla başlayır (`// ErrNotFound is returned when...`)

**Xətalar (Errors):**
```go
return fmt.Errorf("upload failed: %w", err)  // wrap: %w
if errors.Is(err, ErrNotFound) {...}          // sentinel müqayisə
var e *QueryError
if errors.As(err, &e) {...}                   // tip müqayisəsi
```
- Panic yalnız həqiqətən istisna halda; `_` ilə xəta atma YOX
- Xəta stringləri kiçik hərflə, durğu işarəsisiz

**İnterfeyslər:** istifadə olunmamış interfeys əvvəlcədən təyin etmə; konkret
tip qaytar (interfeys yox)

**Context:**
- context.Context — ilk arqument; I/O edən hər funksiyaya ötür
- İmmutabledir, metadata ilə klonlanır; cancel/timeout/metadata (tracing)
- Struct-a bağlama YOX

### 3. Proyekt strukturu
| Qovluq | Mənası |
|---|---|
| `internal/` | YALNIZ eyni kök daxilindən import oluna bilər — artıq asılılıqdan qorunma (spaghettification-a qarşı) |
| `pkg/` | Xarici istifadə üçün açıq (rəsmi tövsiyə deyil, geniş yayılıb) |
| `cmd/` | main() olan executable paketlər (cmd/indexer, cmd/crawler) |
| `api/`, `testdata/`, `web/` | şema/protokol faylları, test datası, web assetləri |

- Fayllar: `main.go`, `doc.go`, `*_test.go`, `README.md`, `LICENSE`
- Qranullik balansı: çox erkən bölmə YOX, amma nəhəng tək paket də YOX

### 4. Movie aplikasiyasının dizaynı
- **Metadata servisi:** statik data, ID ilə oxu → key-value/DB
- **Rating servisi:** dinamik, append/delete + **aqreqasiya** → ayrı saxlanma
- Gələcək üçün düşün: rating recordType sahəsi ilə generikləşdi (film,
  aktyor ifası, soundtrack) — "6-12 ayda lazım olacaqmı?" testi
- Bölünmə meyarları: logic loosely coupled, data modelləri fərqli, data müstəqil

**3 servis:**
1. **metadata** (:8081) — film metadata DB-si
2. **rating** (:8082) — record rating-ləri + aqreqasiya
3. **movie** (:8083) — client-facing API, ikisini birləşdirir (aggregator)

### 5. Hər servisin qatları
```
cmd/                    ← main: komponentləri yığıb HTTP server başladır
internal/controller/    ← biznes məntiqi
internal/handler/http/  ← API handler (JSON encode/decode, status kodlar)
internal/repository/    ← DB məntiqi (memory implementasiyası test üçün)
internal/gateway/       ← başqa servislərin çağırılması (yalnız movie-də)
pkg/model/              ← export olunan tiplər (shared)
```
- Handler DB-yə birbaşa MÜRACİƏT ETMİR — məntiq controller-də
- HTTP-dən gRPC-ə keçidsə məntiq iki dəfə yazılmır

### 6. Kod detalları
- memory repo: `sync.RWMutex` + map (metadata: `map[string]*Metadata`;
  rating: nested `map[RecordType]map[RecordID][]Rating`)
- Xüsusi tiplər `RecordID`, `RecordType`, `UserID` — özünü izah edən açarlar,
  əlavə tip qoruması
- Aqreqasiya: ratings cəmi / say = ortalamа; boş → ErrNotFound
- movie controller: metadata tapılmadı → ErrNotFound; rating tapılmadı → **əvvəlki
  addımla davam** (rating mütləq deyil), `Rating *float64` optional
- ErrNotFound hər komponentdə AYRI təyin olunur (rating vs metadata qarışmasın)
- Gateway: non-2xx → error, 404 → ErrNotFound map-lənir
- Statik unvanlar (localhost) — cloud-da işləməz → növbəti fəsil: discovery

## Termindirmə (AZ)
- Scaffolding — Skelet/İlkin Quruluş (proyektin ilkin strukturunu yaratma)
- Idiomatic Code — İdiomatik (dilə xas) kod
- Repository Pattern — Repozitoriya Naxışı (data girişi qatı)
- Gateway — Keçid (başqa servislə çağırış qatı)
- Aggregation — Aqreqasiya (toplama/birləşdirmə)
- Spaghettification — Spagetti Asılılıq (nəzarətsiz paket asılılıqları)

## Kviz sualları
1. `internal/` qovluğunun semantik mənası nədir? (Yalnız eyni kökdən import
   oluna bilər — xarici paketlər üçün bağlıdır)
2. Niyə handler və controller ayrılıdır? (Biznes məntiqi API növündən asılı
   olmasın — HTTP→gRPC keçidi məntiqi təkrarlamasın)
3. Rating-ə `RecordType` niyə əlavə edildi? (Gələcək rating növləri — aktyor,
   soundtrack — API-ni dəyişmədən dəstəklənsin)
4. Movie servisdə rating tapılmayanda nə olur? (ErrNotFound xətası DEYİL —
   rating-siz cavab qaytarılır)
