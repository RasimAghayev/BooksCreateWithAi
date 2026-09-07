# Chapter 10 — Dependence day (Asılılıq Günü)

## Bu fəsil nədən bəhs edir?

Dependency-lərin testə təsiri: toxic codependency, tight coupling; həllər —
komponentlərin birləşdirilməsi, dependency scope-un azaldılması
(EmailUserIfAccountIsNearExpiry → IsNearExpiry + EmailRenewalReminder),
dependency injection-ə şübhə (CreateUserInDB(user, db) — "moving the furniture"),
test-induced damage (GetForecast(location, apiURL) — API-nin səthini 2 qat
artırır), chunking (FormatURL + ParseResponse — magic package, network yoxdur),
router nümunəsi (switch → map[RequestType]Handler — map datası test olunur),
adapter pattern (ambassador — DB/PostgresStore/Store interface, mapStore fake,
sqlmock stub, DATA-DOG), fakes/stubs/spies/mocks/crock terminologiyası,
mock-lara qarşı müdafiə (brittle, implementation-i test edir, interface
pollution), singleton + time.Now (OneHourAgo flaky 00:00-01:00 — NetHack
analogiyası), time-as-data (parametr və ya `var Now = time.Now` funksiya
dəyişəni + seam), wall clock vs monotonic clock, Sub/Abs delta muqayisəsi
(== YOX), "ən yaxşı test lazımsız funksiyanı göstərəndir".

## Əsas fikirlər

### 1. Dependency Nədir və Niyə Problem
Komponent işini İZOLYASİYADA görə bilmir — başqa komponent/xidmət/resurs (DB,
API, terminal) lazımdır. Test üçün: bütün asılılıqları qane etmək lazımdır —
bəzən çətin.

**Mülayim asılılıqlar** — sadəcə konstrue et və ver (print üçün buffer).
**Toxic codependency** — çoxlu komponentlər bir-birinə sıx bağlı; anlaşılmaz,
dəyişilməz, test olunmaz. **Brittle test** = dizayn əks-sədası: "If an object
is difficult to decouple from its environment because it has many dependencies
or its dependencies are hidden, its tests will fail when distant parts of the
system change." (GOOS)

**Test-driven yazmaq** over-engineering, awkward API və test-olunmaz
funksiyalardan qoruyur — "Just don't write untestable functions" ən yaxşı
strateji, amma legacy-də başqa yollar lazımdır. "Every time you encounter a
testability problem, there is an underlying design problem." (Michael Feathers)

### 2. Həll 1 — Birləşdir (Merge)
Komponent A dəyişir → B-nin testləri fail; B dəyişir → A fail. Bu TIGHT
COUPLING işarəsidir. Birbaşa etiraf: MÜTLƏQ ayrı deyillər — birləşdir!
Kod strukturu onları ayrı göstərsə belə, müstəqil istifadə olunmursa AYRI
komponent deyillər.

### 3. Həll 2 — Scope-u Azalt (Email Nümunəsi)
`account.EmailUserIfAccountIsNearExpiry()` — test üçün kabab: IMAP server,
DNS, spamblock, müvəqqəti dəyişənlər... **Bu problemi HƏLL ETMƏK deyil,
LƏĞV ETMƏK istəyirik.**

"What are we really testing here?" → Vacib olan EMAIL GÖNDƏRİLMƏSİ deyil —
DÜZGÜN QƏRAR (expiry yaxındır?) verilməsidir. Email = qərarın nəticəsi.
```go
if account.IsNearExpiry() {       // PURE decision — 2 fake account ilə sadə test
    account.EmailRenewalReminder() // ayrıca, asılı tərəf
}
```
İkinci fayda: hər komponent ayrılıqda DAHA ASAN test olunur. Email göndərmə
hələ də test ediləcək — amma ARTIQ decoupled halda.

### 4. Dependency Injection-ə Şübhə
`CreateUserInDB(user)` — IMPLICIT dependency: hardwired qlobal *sql.DB /
config dərinliyi / preset address. Test-ə yönəldilməmiş kodun əlamətidir
(prod-da bir DB var, nə üçün başqa yol verilsin?).

**DI = "passing in dependencies as arguments":** `CreateUserInDB(user, db)` —
yaxşılaşma, amma "moving the furniture around": indi test DB-ni KONSTRUE
etməli! Dedicated test DB, Docker, Makefile... İstəmirik.

**Əsl problem:** CreateUserInDB-nin vacibliyi istifadəçinin DB-də YARADILMASI
DEYİL — "under the right circumstances it WOULD CAUSE the user to be created".
Real DB-yə bağlanan test BRITTLE-dir: DB işləmir, şəbəkə, credential, permission,
schema, DB bug... — test DAXİLİNDƏN environment, network və "someone else's
database product" test edir. **Injecting the dependency doesn't fix that —
THE DEPENDENCY IS THE PROBLEM.**

**DI-nin məqsədi:** çox asılılıqlı komponentlər yazmağı ASANLAŞDIRMAQ — amma bu
pis ideyanı asanlaşdırmaq istəmirik.

### 5. Test-Induced Damage (Weather Nümunəsi)
`weather.GetForecast(location)` — implicit real-API asılılığı (yavaş, flaky,
pullu, nuisance traffic). Onların API-sini test ETMİRİK — bizim kod deyil.
Həmçinin "onların kodu İŞLƏMƏDİKĐƏ" bizim kodun davranışı lazımdır — API-ni
söndürməyi xahiş edə bilmərik.

**Səhv həll — URL inject:**
```go
fc, err := weather.GetForecast(location, apiURL)   # API səthi 2X
fc, err := weather.GetForecast(location, "https://weather.example.com")  # istifadəçi həmişə eyni dəyəri ötürür
```
Go-da OPTIONAL ARGUMENT YOXDUR → real istifadəçilər həmişə real URL-i yazmalı
— 100% əlavə paperwork, test üçün. **Test-induced damage = API-nin test qoxusu
ilə çirklənməsi.**

### 6. Həll 3 — Chunking (Behavior Parçalama)
GetForecast-u 3 addıma böl:
1. URL/HTTP request KONSTRUE et (location → URL)
2. API-yə göndər
3. Cavabı PARSE et və formatla

Addım 1 və 3 PRİNSIPCƏ asılı deyil! **Magic package:** FormatURL, ParseResponse
"sehirli" çağır:
```go
func TestFormatURL_IncludesLocation(t *testing.T) {
    t.Parallel()
    location := "Nowhere"
    want := "https://weather.example.com/?q=Nowhere"
    got := weather.FormatURL(location)
    if want != got { t.Error(cmp.Diff(want, got)) }
}

func TestParseResponse(t *testing.T) {
    t.Parallel()
    data := []byte(`{"weather":[{"main":"Clouds"}]}`)   // real API-dən tutulmuş blob da ola bilər
    want := weather.Forecast{Summary: "Clouds"}
    got, err := weather.ParseResponse(data)
    ...
}
```
Network ÇAĞIRIŞI YOX. GetForecast-da yenidən yığ:
```go
func GetForecast(location string) (Forecast, error) {
    URL := FormatURL(location)
    ...                                  // http.Get + status check — standart, az səhvriski
    return ParseResponse(resp.Body)
}
```
Mərkəzi (step 2) manual test kifayət — "The goal is not to write tests. The
goal is to write correct software." (Michael Feathers). URL xətaları və parse
xətaları — ən çox səhv edilən yerlər — tam və izolyasiyalı test olunur.

### 7. Həll 4 — Əsas Məntiqi Çıxart (Router Nümunəsi)
Mürəkkəb API server: onlarla endpoint, hər çağırış ekranlarla kod, side effects,
response-lar eyni → HANSI handler çağırıldı bilinmir.

**"What are we really testing here?"** → "requests are passed to the correct
handler for their type."

Switch → MAP çevir:
```go
handler := map[RequestType]Handler{
    TypeA: HandleTypeA,
    TypeB: HandleTypeB,
    TypeC: HandleTypeC,
}
handle := handler[requestType(req)]
handle(req)
```
Map lookup və funksiya çağırışı test OLUNMAZ (dil özü işi). Yalnız MAP DATA
test olunmalı — sadə, izolyasiyalı. End-to-end çətin → routing məntiqi asan
(təcrid düzgün qurulsa).

### 8. Həll 5 — Adapter Pattern (Ambassador)
Komponent = xarici asılılıqla ƏLAQƏLİ BÜTÜN kodun cəmi. Diplomatik səfir:
bizim sorğunu xarici dilə TƏRCÜMƏ edir (outbound), cavabı bizim dilə gətirir
(inbound).

**Weather adapter:** outbound = FormatURL (URL = API-nin dili), inbound =
ParseResponse (JSON → Forecast struct). Effekt: API haqqında BÜTÜN bilik bir
komponentdə, sistemin qalanı xəbərsiz.

**DB adapter — CRUD Widget:**
Ilkin (pis): `Create(db *sql.DB, w Widget)` — SQL + biznes məntiqi qarışıq
(Single Responsibility pozulub), test real DB tələb edir (Makefile/Docker/
testcontainers — "sumo, not judo").

**Interface ilə qırma:**
```go
type Store interface {
    Store(Widget) (string, error)
}

func Create(s Store, w Widget) (ID string, err error) {
    ID, err = s.Store(w)          // biznes məntiqi burada — DB haqqında bilik YOX
    if err != nil { return "", err }
    return ID, nil
}
```
**Test fake — mapStore (in-memory, real amma non-persistent):**
```go
type mapStore struct {
    m    *sync.Mutex
    data map[string]widget.Widget
}
func (ms *mapStore) Store(w widget.Widget) (string, error) {
    ms.m.Lock(); defer ms.m.Unlock()
    ms.data[w.ID] = w
    return w.ID, nil
}
```
Mutex — paralel testlər arasında paylaşma üçün (kvstore dərsi). Test:
`widget.Create(s, w)` — heç bir xarici server yoxdur, `go test` birbaşa.

**Prod implementation:**
```go
type PostgresStore struct { DB *sql.DB }
func (p *PostgresStore) Store(w Widget) (string, error) { /* horrible SQL */ }
```
mapStore ilə fərq: arxasında 1.3 milyon sətir kod var — amma Create-un
işlədiyini bilmək üçün onların HƏMİSİNİ test ETMİRİK.

**PostgresStore özü necə test olunur?** sqlmock (DATA-DOG) — STUB kimi:
```go
func fakePostgresStore(t *testing.T) widget.PostgresStore {
    db, mock, err := sqlmock.New()
    if err != nil { t.Fatal(err) }
    t.Cleanup(func() { db.Close() })
    query := "SELECT id, name FROM widgets"
    rows := sqlmock.NewRows([]string{"id", "name"}).
        AddRow("widget01", "Acme Giant Rubber Band")
    mock.ExpectQuery(query).WillReturnRows(rows)
    return widget.PostgresStore{DB: db}
}

func TestPostgresStore_Retrieve(t *testing.T) {
    t.Parallel()
    ps := fakePostgresStore(t)
    want := widget.Widget{ID: "widget01", Name: "Acme Giant Rubber Band"}
    got, err := ps.Retrieve("widget01")
    ...
}
```
Retrieve: `QueryRowContext + row.Scan` — row → Widget TƏRCÜMƏSİ test olunur;
Postgres-un özü DEYİL ("let's hope it does"). Etiraf: SQL dəqiqliyi və canned
data real serverlə eynilik YOXDUR → ara-sıra real serverlə 1-2 backstop test
(deploy-dan əvvəl). Dependency scope MİNİMUMA endi.

### 9. Fakes / Stubs / Spies / Mocks / Crock — Terminologiya
Müəllif hamısını "fake" başlığı altında birləşdirir:
- **Stub** — heç nə etmir / sabit cavab; kompilyasiya üçün. DESIGN SMELL:
  stub-la işləyirsə — ya çox bağlıdır, ya asılılıq o qədər də vacib deyil.
- **Fake (dar)** — mapStore kimi: real davranış, sadə texnologiya
  (in-memory DB).
- **Spy** — istifadəni QEYD EDİR (bytes.Buffer kimi "fake" io.Writer — nə
  yazıldığını sonra yoxlayırsan).
- **Mock** — proqramlaşdırılmış EXPECTATION-lar; özü testi FAIL edə bilər.
- **Crock** — bilərəkdən SƏHV işləyən fake (errReader — həmişə error qaytarır).

### 10. Mock-Lara Qarşı Müdafiə
BuyWidget + ödəniş xidməti: stub → "çağırıldımı?" məlum deyil; spy → test özü
sonra yoxlayar (yaxşı variant); mock → *testing.T qəbul edir, dəqiq payment
detalları GÖZLƏYİR, uyğunsuzluqda ÖZÜ fail edir, çağırılmayıbsa ÖZÜ fail edir.

**Mock problemləri:**
1. Real komponent qədər mürəkkəb ola bilər — doğruluğu NECƏ bilirsən?
   Səhv mock = yanlış şeyin testi.
2. **Indirect outputs-u test edir** — çağırış SEQUENCE-i + parametrlər →
   İMPLEMENTASİYAYA bağlı → refactor (davranış eyni) mock-testləri QIRIR =
   brittle. Görmək istədiyimiz: user-visible davranış.
3. Mock DİZAYNI presuppose edir — düzgün dizayn davranışdan EMERGE olmalıdır.
4. Mock kodu yazmaq + ömrü boyu real komponentlə sinxron saxlamaq = iş.
5. Interface tələb edir → interface yalnız mock üçün = **interface pollution**
   (test-induced damage).

**Vəziyyət:** mock pis deyil — SON option. Bütün digər yollar tükənəndək yox.
**Hipokriziya check:** sqlmock istifadə etdik — texniki olaraq STUB kimi + hazır
import; sql.DB fake-ləmək çətindir → lesser evil. errReader — mock deyil
(özü fail ETMİR), crock-dur.

### 11. Singleton Problemi + Vaxt
Qlobal DB = singleton — bir DB varsayımı → fake əvəzetmə çətin. Həll:
dependency-ni EKSPLİSİT et → DB sadəcə PARAMETR olur.

**Gizli singleton #2: VAXT.** `time.Now` çağıran proqram günün saatından
asılı → testlər də. **Real layihədən nümunə (flaky!):**
```go
func TestOneHourAgo(t *testing.T) {
    now := time.Now()
    then := past.OneHourAgo()
    if now.Hour() == 0 {
        assert.Equal(t, 23, then.Hour())          // 00:00-01:00 arası XÜSUSİ HAL
        assert.Equal(t, now.Day()-1, then.Day())
    } else {
        assert.Equal(t, now.Hour()-1, then.Hour())
        assert.Equal(t, now.Day(), then.Day())
    }
    ...
}
```
Programmist problemi DOĞRU tapdı (flaky, saatdan asılı) — amma SƏHV həll:
testi günlük saatlara görə PATCH-ləmək. (NetHack: gremlins 10pm-6am oğurlayır,
undead midnight-də 2X zədə — oyunda əyləncəlidir, testlərdə YOX.)

**Həll A — parametr:**
```go
then := past.OneHourAgo(now)     // hər hansı vaxt üçün hesabla
```
Test: `time.Parse(RFC3339, "2022-08-04T23:00:00Z")` + want `22:00:00Z` —
MİDNIGHT xüsusən test OLUNUR (özünü təsdiqləyən müqayisə, assert yığını YOX).
Zərər: istifadəçi həmişə time.Now() konstrue etməli + bütün çağırışları refactor
— kiçik test-induced damage.

**Həll B — funksiya dəyişəni (seam):**
```go
var Now = time.Now          // DİQQƏT: time.Now DEYİL — funksiyanın ÖZÜ

func OneHourAgo() time.Time { return Now().Add(-time.Hour) }
```
Testdə:
```go
past.Now = func() time.Time { return testTime }   // inject
got := past.OneHourAgo()
```
Mock object-dən ÇOX sadə: interface lazım deyil, istifadəçiyə paperwork YOX,
tam transparent. Test PARALEL OLA BİLMƏZ (qlobal dəyişən dəyişir) —
singleton fake-ləyən hər test kimi sequential.

**Hər iki həll = "turn time into data":** vaxt sistemə magic singleton-clock
yox, DATA kimi daxil olur.

### 12. Vaxt Sadəcə Rəqəm DEYİL
"Döşqapı kartlarından fərqli olaraq, time is not just a number." — Doctor Who:
"wibbly-wobbly timey-wimey stuff".

İki saat var: **wall clock** (dayandıra, dəyişə, GERİ gedə bilər — DST, leap
second) və **monotonic clock** (yalnız İRƏLİ). time.Time hər ikisini daşıyır.

**Nəticə:** time.Time == ilə müqayisə ETMƏ, riyaziyyat DÜZ etmə — time API işlət:
```go
func TestOneHourAgo(t *testing.T) {
    t.Parallel()
    now := time.Now()
    want := now.Add(-time.Hour)
    got := past.OneHourAgo()
    delta := want.Sub(got).Abs()          // Duration
    if delta > 10*time.Microsecond {       // komputasiya vaxtı toleransı
        t.Errorf("want %v, got %v", want, got)
    }
}
```
Equal dəqiq işləməz — hesablama SIFIR-dan böyük vaxt aparır. Sub → Duration →
Abs → kiçikdeltadan az olmalı. **Final dərs:** bu test göstərir ki, OneHourAgo
FUNKSIYASI (və testi) lazım deyil — `now.Add(-time.Hour)` birbaşa çağır.
"Ən yaxşı test bəzən funksiyanın lazımsızlığını göstərəndir."

## Əsas terminlələr
- Dependency — izolyasiyalı işləməyə mane olan komponent/xidmət
- Toxic Codependency — qarşılıqlı sıx bağlı komponent dəsti
- Brittle Test — uzaq dəyişikliklərdən qırılan test (dizayn siqnalı)
- Merge — ayrı görünən, amma ayrı olmayan komponentlərin birləşdirilməsi
- Scope Reduction — IsNearExpiry + EmailRenewalReminder bölgüsü
- Dependency Injection — asılılığı arqument kimi ötürmək (faydası ŞÜBHƏLİ)
- Test-Induced Damage — API-nin yalnız-test-parametrləri ilə zədələnməsi
- Interface Pollution — yalnız mock üçün interface
- Chunking — davranışı test-olunan dilimlərə bölən (FormatURL/ParseResponse)
- Magic Package — mövcud olmayan funksiyanı təsəvvür edib test yazmaq
- Adapter / Ambassador — xarici asılılıq kodunun bir komponentdə cəmi
- Store interface — abstract yaddaş; mapStore / PostgresStore implementasiyaları
- sqlmock — sql.DB stub-u (DATA-DOG); canned rows
- testcontainers — testlərdə konteyner başlatma ("sumo")
- Fake / Stub / Spy / Mock / Crock — test ikamediləri şkalaı
- Indirect Outputs — mock-un yoxladığı çağırış sequence-i (brittle!)
- Singleton — "yalnız bir nümunə" varsayımı (qlobal DB, time.Now)
- Seam — fake inject oluna bilən yer (parametr, `var Now = time.Now`)
- Turn Time Into Data — vaxtı gizli singleton-dan data-girişə çevir
- Wall Clock vs Monotonic Clock — geri gedə bilən / yalnız irəli
- Sub().Abs() + delta — == əvəzinə toleranslı vaxt müqayisəsi

## Praktik nəticə
(1) Asılılıq problemi görsən — DİZAYN problemidir: birləşdir / scope azalt /
chunk-la / adapter-la / çıxart. (2) "What are we really testing here?" — email
GÖNDƏRİLMƏSİ yox, QƏRAR; DB-yə YAZILMA yox, "would cause". (3) DI panaseya
deyil: asılılığı ötürmək onu həll etmir — dependency-nin ÖZÜ problemidir.
(4) Yalnız-test parametri əlavə etmə (apiURL) = API zədəsi; Go-da optional
argument yoxdur. (5) Weather: FormatURL + ParseResponse chunk-la; http.Get
qalığını manual testə burax. (6) Router: switch → map; map DATANI yoxla.
(7) Adapter: xarici bilik bir komponentdə; Store interface ilə biznes məntiqi
DB-dən təcrid; mapStore fake + mutex; sqlmock stub-u ilə SQL tərcümə testi;
real serverlə ara-sıra backstop. (8) Stub = smell; spy yaxşıdır; mock = son
option (indirect output, brittle, presuppose); interface-i YALNIZ mock üçün
YARATMA. (9) time.Now = gizli singleton: flaky midnight testlərini patch-ləmə —
parametr və ya `var Now = time.Now` seam işlət (paralel testlərdə diqqət!). (10)
Vaxt rəqəm deyil: == YOX, Sub/Abs + delta toleransı; wall/monotonic fərqi.
(11) Ən yaxşı test bəzən funksiyanın LAZIMSIZLIĞINI göstərir — OneHourAgo
yerinə now.Add(-time.Hour).

## Mənbə
Pages: 309-345 (PDF 321-357)
