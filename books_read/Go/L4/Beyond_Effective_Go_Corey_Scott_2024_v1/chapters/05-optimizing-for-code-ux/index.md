# Chapter 5 — Optimizing for Code UX (səh. 69-147)

## Bu chapter nədən bəhs edir?

Code UX (kod istifadəçi təcrübəsi) — kodun **oxunaqlılığı, tutarlılığı və
proqnozlaşdırıla bilənliyi**. Müəllifin tezisi: kod əsasən insanlarla ünsiyyət
vasitəsidir; performansdan, ağıllıbuludluqdan və innovasiyadan əvvəl
istifadəçilik gəlir. Üç sütun: **Clarity** (aydınlıq), **Consistency**
(tutarlılıq), **Predictability** (proqnozlaşdırma).

## Əsas fikirlər

### 1. Clarity — adlandırma və format
**Nədir:** Kodun tez və intuitiv başa düşülməsi.

**Ad qaydaları (kitabın 5 prinsipi):**
- **Meaningful (mənalı):** `cur` → `create` → **`userCreationRequest`**;
  `get()` → **`getUserDecision()`**
- **Concise (qısa, amma oxunula bilən):** `order` ✓, `customerOrder` artıqdır
- **Consistent (tutarlı):** bütün kod eyni üslubda adlandırsın
- **Context-aware (kontekstli):** paket adı kontekst verir → `user.User` ✓,
  `user.UserRequest` "stutter" edir; `AddressStreet` struktur daxilində artıqdır
  (`Street` kifayətdir)
- **Original olmaya bilər:** standart adlar təkrar istifadə olunsun

**Whitespace (boşluq):** Məntiqi qrupları boş sətirlə ayır; amma `if err != nil`
blokun özündən AYIRMA — xəta yoxlaması və emalı bir qrupdur.

**Kitabdan kod nümunəsi (early return):**
```go
func validateItems(items []item) error {
    if len(items) <= maxItems {
        return nil
    }
    return errors.New("too many items")
}
```
**Sub-kod izahı:**
- `else` bloku silinib → "validate, process, return" axını
- Happy path funksiyanın bütün uzunluğu boyu düz axır

### 2. Shadowing (kölgələmə) və yad üslublar
**Shadowing nədir:** Daxili `for err := range resultCh` xarici `err`-i gizlədir
— hansının qaytarıldığı qeyri-müəyyən olur. Həll: ya yeni ad, ya `=` ilə təyinat,
ya da funksiyanı parçala.

**Yad idiomatic olmayan üslublar (digər dillərdən gələn):**
- Strict singleton → əvəzinə **Dependency Injection**: `NewEncoder(pool ObjectPool)`
  — pool-un yaradılmasını Encoder-dən ayırır, concurrency-test mümkün olur
- Artıq formalizm (Java-tipli FileBuilder struct-lar) → paket səviyyəli
  funksiyalar Go-da daha təbii
- Yanlış yerdə interfeys → "təmin edirəm" yox, **"ehtiyacım var"** perspektivi

### 3. Errors are not exceptional (xətalar istisna deyil)
**Nədir:** Go-da error adi dəyərdir — çağıran tərəf axına təsir qərarını verir.

**-1/-2/-3 return kodları QADAĞANDIR.** Əvəzi:
```go
var (
    ErrBadID        = errors.New("bad ID")
    ErrInvalidFormat = errors.New("invalid format")
    ErrParseFailed  = errors.New("failed to parse version")
)

func ExtractVersionFromID(id string) (int, error) {
    if id == "" {
        return 0, ErrBadID
    }
    // ...
    return version, nil
}
```

**Error handling pattern:**
```go
switch err {
case ErrBadID:
    resp.WriteHeader(http.StatusBadRequest)
    return
// ...
default:
    resp.WriteHeader(http.StatusInternalServerError) // system errors
}
```
**Sub-kod izahı:**
- Adlandırılmış xətalar → caller fərqli reaksiya verə bilər
- `default` → proqnozlaşdırıla bilməyən "system errors" üçün həmişə olmalıdır

**Wrapping (Go 1.13+):**
```go
return nil, fmt.Errorf("%w - %s", ErrBadID, err)
```
- `%w` → sarılan (wrapped) xətanı işarələyir; `errors.Is` ilə yoxlanıla bilər
- Aşağı səviyyəli xətalara business kontekst əlavə etmək üçün də istifadə olunur

### 4. Globals, init() və panic()
- **Globals/init() QADAĞANDIR:** magic state, data race, test sırası asılılığı
  yaradır. Funksiyanın gizli state-i varsa → struct metoduna çevir.
- **İstisnalar:** named errors, enum-lar (`const North Direction = "north"`),
  init-dən sonra dəyişməyən private map-lər.
- **panic() QADAĞANDIR** (xidməti kodda): defer/recover pattern-i kodu qısaltmış
  kimi görünür, amma call chain-i oxunmaz edir. `main()`-də `log.Fatalf`,
  testlərdə `t.Fatalf` istifadə olunur.

### 5. Consistency — stil, paketlər, fayl təşkili
**Stil:** `gofmt -s` həmişə; qalanı team qərarı (məs. channel adlarında `Ch`
soneksi, assertion library). Məqsəd — dəbdən çox, komfort və naviqasiya.

**Paketlər:**
- Feature/module üzrə təşkil olunmalı, **`util`, `dto`, `errors`, `constants`
  adlı paketlər QADAĞANDIR** (paket adı kontekst verməlidir)
- Self-contained: request/response tipləri (dto) ayrı paketdə deyil, feature-in
  özündə
- Paket şərhi bir cümlədə məqsəd + scope (`// Package httptest provides
  utilities for HTTP testing.`)

**Fayl daxili sıra (makrodan mikroya):**
```
1. Constants
2. Global variables (named errors)
3. Constructor
4. Struct definition
5. Public methods/functions
6. Private methods/functions
7. Input interfaces + private types
```

### 6. Funksiyalar və arqumentlər
- **Stateless funksiya, state-li metod:** funksiya state istəyirsə → struct metodu
- **Uzunluq ~30 sətir** (bir ekrana sığmalı)
- **Arqument sayı:** ideal 0; kompozit (>3) → qruplaşdır
- **Scope ayırımı:** init-scoped data → konstruktor/member; request-scoped →
  metod arqumenti
- **Qruplaşdırma:** `func (s Sender) Send(from, to string, email Email) error`
- **Sıra:** `context.Context` həmişə birinci, `error` həmişə sonuncu

### 7. Boolean arqumentlər təhlükəlidir
```go
result := user.CheckGender(true) // true nədir??
```
Həll: `IsMale()` / `IsFemale()` kimi iki spesifik metod. Callback-in "tapıldı"
bool-u əvəzinə named error (`ErrOrderNotFound`) — caller maraqlı deyilsə eyni
reaksiya verir.

### 8. Konstruktor arqumentlərinin azaldılması
- **Yalnız test üçün ötürülən asılılıq QADAĞANDIR** (məs. `now func() time.Time`
  hamı eyni default verir → type-ə daxil et, test üçün setter)
- Hamı eyni dəyəri ötürürsə → DI məcburiyyəti yoxdur
- Konfiqurasiya çoxalırsa → **Config injection**: `NewUserManager(cfg Config, ...)`
  — yeni konfiq gələndə konstruktor dəyişmir, Config interfeysinə metod əlavə
  olunur
- **Konstruktor error qaytarırsa — code smell:** xarici sistemin
  əlçatanlığını app başlanğıcına bağlayır (resilience itir)

### 9. Predictability — API və encapsulation
- **Export only what you must:** hər şeyi private et, zərurət yaranana qədər.
  İstisna: `internal/` direktivası ilə modul daxilində məhdudlaşdır (fictive
  MyStore nümunəsi: `api/internal/httputil` yalnız `api` moduluna açıqdır).
- **Encapsulation = information hiding:** struct level-də internal state,
  paket level-də private funksiyalar gizlədilir. **Sızma nümunəsi:**
  `Storage.DoQuery(q string, ...)` — SQL implementasiya detalı interfeysə
  sızmış; əvəzinə domain-ə uyğun metodlar (`LoadUser`, `SaveOrder`).
- **Constructors:** sadə və yalnız init etsin; variantlı konstruktor →
  functional options pattern (Chapter 8-də).

### 10. A little copying vs. a little dependency
Hər asılılığın qiyməti var (build time, coupling, komplekslik). Tək funksiya
qədər kod üçün **kopyalamaq** asılılıq əlavə etməkdən yaxşıdır. "Shared" kodu
`util`/`common` paketinə atmaq problemi gücləndirir.

### 11. Conflicting goals (ziddiyyətli məqsədlər)
Default prioritet: **code UX**. Performans optimallaşdırması yalnız ölçülmüş
ehtiyacda. "Hansı UX pozulsun" sualına cavab: **əsas istifadəçiyə** (reader)
fayda verən seçilər. Nümunə: login endpoint-i gündə minlərlə çağırılır →
100ms-dən 80ms — report upload-un 10s-dən 7s-dən daha qiymətlidir.

## Əsas terminlər

- Code UX (kod istifadəçi təcrübəsi)
- Idiomatic (dilə xas təbiilik)
- Early Return (erkən qayıdış)
- Variable Shadowing (dəyişən kölgələməsi)
- Named Error (adlandırılmış xəta)
- Error Wrapping (xəta sarılması)
- Dependency Injection (asılılıq inyeksiyası)
- Information Hiding (informasiya gizlətmə)
- Functional Options Pattern (funksional seçənəklər nümunəsi)

## Praktik nəticə

- `gofmt -s` + adların 5 qaydası + early return = clarity-nin 80%-i
- Xətalar dəyər kimi qaytarılır, named + wrapped (%w) olur; bool əvəzinə named error
- Global/init/panic-dən qaçın; state → struct
- Paket adı kontekst verir; `util`/`dto` paketləri yaratma
- API-ni minimal export et, `internal/` ilə sərhəd çək
- Kiçik kod üçün copy > dependency

## Mənbə

Pages: 69-147 (Chapter 5, Beyond Effective Go Part 2)
