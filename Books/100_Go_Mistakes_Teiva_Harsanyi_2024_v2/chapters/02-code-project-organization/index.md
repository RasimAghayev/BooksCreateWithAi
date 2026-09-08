# Chapter 2 — Code and Project Organization (#1-#27) (səh. 28-77)

## Bu fəsil nədən bəhs edir?

Səhvlər #1-#16 (kod təşkilatı): variable shadowing, nested code, init funksiyaları,
getters/setters, interface pollution, producer-side interface, any tipi, generics
istifadəsi, type embedding, functional options, project misorganization, utility
paketlər, dokumentasiya, linter-lər.

## Əsas səhvlər və həllər

### #1: Unintended variable shadowing (qeyri-ixtiyari kölgələmə)
**Səhv:** `if` daxilində `client, err := ...` xarici `client`-i kölgələyir —
if blokundan sonra client nil qalır.
**Həll:** xaricdə `var client *http.Client; var err error` elan et, daxildə
`=` (assign) istifadə et, `:=` YOX.

### #2: Unnecessary nested code (lazımsız yuvalanma)
```go
// PIS: 4 səviyyə else
if s1 != "" {
    if s2 != "" { ... } else { return ... }
} else { return ... }

// YAXŞI: happy path SOLDA düz axın (Mat Ryer "line of sight")
if s1 == "" {
    return "", errors.New("s1 is empty")
}
if s2 == "" {
    return "", errors.New("s2 is empty")
}
return concatenate(s1, s2)
```

### #3: Misusing init functions
- init: `var` → `init()` → `main()` sırası; asılı paketlərin init-ləri ƏVVƏL
- **Problemlər:** xəta idarəsi yoxdur (app crash), çağıran tərəf retry/fallback
  edə bilmir, test çətin, global state
- **Yalnız:** xəta verə bilməyən statik konfiqurasiya (məs. flaglardan asılı
  olmayan handler qeydiyyatı) üçün

### #4: Overusing getters and setters
- Getter/setter hər sahə üçün MƏCBURİ deyil (Java fərqi)
- Dəyəri VARSA: interception/debug nöqtəsi, dəyəri dəyişə bilən gələcək
  compatibility

### #5: Interface pollution
İnterfeys 3 halda dəyərli:
1. **Orta davranışı faktorlaşdırmaq** (məs. io.Reader/Writer, sort.Interface —
   abstraksiyanın DOĞRU səviyyəsi)
2. **Decoupling** — asılılıq abstraksiyadan (test/mock mümkün)
3. **Davranışı məhdudlaşdırmaq** — tipi konkret davranışa söndürmək
   (IntConfig struct yalnız Get təklif edir)

### #6: Interface on the producer side
- İnterfeys İSTEHLAKÇI tərəfində yerləşir, producer YOX
- "Abstractions should be discovered, not created" — istehsalçı öz tipini
  qaytarsın; client lazım olanda ÖZ interfeysini yaradır
- İstisna: standart kitabxana (encodings) — əvvəlcədən məlum universal ehtiyac

### #7: Returning interfaces
- Funksiya interfeys YOX, konkret struct qaytarmalı (accept interfaces, return
  structs); client interfeyslə birləşdirmə azadlığı

### #8: any says nothing
- `any` (interface{}) heç nə demir (Rob Pike) — type safety itir; runtime
  assertion-lar + xətalar; yalnız REAL ehtiyacda (json.Marshal kimi)

### #9: Being confused about when to use generics
- **Nə vaxt:** boilerplate faktorlaşdırmaq (getKeys map[K]V), custom constraint:
```go
type customConstraint interface{ ~int | ~string }

func getKeys[K customConstraint, V any](m map[K]V) []K { ... }
```
- **Nə vaxt YOX:** "bəlkə gərək olar" deyə əvvəlcədən — "abstraksiya
  kəşf olunur, yaradılmır"; konkret ehtiyac + boilerplate görünəndə

### #10: Type embedding problemləri
- Embedding = subclassing DEYİL: metodu alan RESİVER embed edilən tipdir
- Sync.Mutex-i embed etsən → xarici client Lock/Unlock çağıra bilər (sərhəd
  pozulur); unexported sahə kimi saxla
- Düzgün: logger → io.WriteCloser embed (Write metodu avtomatik təklif olunur)

### #11: Not using the functional options pattern
```go
type Option func(*httptest.Server)

func WithPort(port int) Option {
    return func(s *Server) {
        if port <= 0 { panic(...) }  // yaxud error daşıyan Option
        s.port = port
    }
}

func NewServer(addr string, opts ...Option) (*http.Server, error)```
- Config struct: pointer-lər lazım (`*int`), opt-in xəta yoxdur
- Builder pattern: metod zənciri, amma port kimi xətalı hal çətin
- **Functional options:** variadic + closure — API dostu, təkmilləşdirilə bilər

### #12-#15: Project structure, utility packages, collisions, docs
- Mandatory layout YOXDUR (Russ Cox tənqidi: "golang-standards" rəsmi deyil)
- **Paket adı NƏ TƏMİN EDİR:** util/common/shared/base QADAĞAN —
  `stringset` kimi ifadəli ad
- Ad toqquşması: dəyişən adı = paket adı (`x := mysql.Conn` kimi) → alias
- **Hər export olunan element dokumentləşməlidir** — godoc formatı; dəyişən
  üçün: MƏQSƏD public şərhdə, MƏZMUN private

### #16: Not using linters
- Daily: `go vet`, errcheck, staticcheck, golangci-lint (aggregator)

## Əsas terminlər

- Variable Shadowing (kölgələmə)
- Line of Sight (görüş xətti — happy path solda)
- Interface Pollution (interfeys çirklənməsi)
- Producer/Consumer side interface
- Functional Options Pattern
- Type Embedding vs Subclassing
- Utility Package (antipattern)

## Praktik nəticə

- `:=` vs `=` — scope xəritəsini başında qur
- Happy path solda, xətalar yuxarıda return
- İnterfeys: istehlakçıda, kəşf olunanda; any-ə qaçış YOX
- util paketi əvəzinə ifadəli ad; hər export = şərh

## Mənbə

Pages: 28-77 (Chapter 2, 100 Go Mistakes 2nd ed.)
