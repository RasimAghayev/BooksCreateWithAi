# Chapter 4 — Working with Text (səh. 58-72)

## Bu fəsil nədən bəhs edir?

Mətn emalı — fmt format verb-ləri (%v/%+v/%#v), fmt.Stringer, encoding
müəyyənləşdirmə (Content-Type + charset), regexp ilə camelCase→snake_case,
strings.EqualFold (Unicode case-insensitive) və Unicode normalizasiyası
(NFC/NFD/NFKC/NFKD).

## Əsas fikirlər

### Recipe 20 — format verb-ləri
**Tapşırıq:** LocationEvent-i loglamaq — hər hansı tipi oxunaqlı çap.

```go
type LocationEvent struct {
    ID        string
    Latitude  float64
    Longitude float64
}

evt := &LocationEvent{ID: "McQueen95", Latitude: 40.7060361, Longitude: -74.0110143}
log.Printf("info: loc: %#v", evt)
```

**Verb müqayisəsi:**
```go
fmt.Printf("%v\n", evt)   // &{McQueen95 40.7060361 -74.0110143}      — dəyərlər
fmt.Printf("%+v\n", evt)  // &{ID:McQueen95 Latitude:... Longitude:...} — sahə adları
fmt.Printf("%#v\n", evt)  // &main.LocationEvent{ID:"McQueen95", ...}  — tip də
```
- `log` və `fmt` eyni verb-ləri paylaşır; log istifadəçiyə çıxışın
  yönləndirilməsini verir
- **%#v debugging üçün ən yaxşısı:** `id1=1 id2="1"` — int və string fərqi
  görünür (%v ilə ikisi də `1` kimi görünür!)
- log formatını `log.SetFlags()` ilə idarə edin

### Recipe 21 — fmt.Stringer
**Tapşırıq:** Priority uint8 — çapda 20 yox "medium" görünsün.

```go
type Priority uint8

const (
    Low    Priority = 10
    Medium Priority = 20
    High   Priority = 30
)

// String implements the fmt.Stringer interface.
func (p Priority) String() string {
    switch p {
    case Low:
        return "low"
    case Medium:
        return "medium"
    case High:
        return "high"
    }
    return fmt.Sprintf("<%d>", p)   // naməlum dəyər üçün
}

type Bug struct {
    Title    string
    Priority Priority
}

bug := Bug{"Bug level is printed as number", Medium}
fmt.Printf("%+v\n", bug)
// {Title:Bug level is printed as number Priority:medium}
```
- fmt reflection ilə Stringer implementasiyasını yoxlayır → String()
  avtomatik istifadə olunur
- String() daxilində %s/%v işlətsəniz **sonsuz rekursiya** (String →
  printf → String → ...)
- %#v üçün ayrıca: `fmt.GoStringer` / `fmt.Formatter` (0x14 görünür əks halda)

### Recipe 22 — encoding aşkarlanması
**Tapşırıq:** crawler-dən gələn []byte — hansı encoding ilə string-ə
çevirmək?

```go
// Yol 1: Content-Type header:
func ctypeEncoding(ctype string) string {
    _, params, err := mime.ParseMediaType(ctype)
    if err != nil {
        return ""
    }
    return params["charset"]   // "text/html; charset=utf-8" → utf-8
}

// Yol 2: data-dan təxmin:
func dataEncoding(data []byte) string {
    _, name, certain := charset.DetermineEncoding(data, "text/plain")
    if certain {
        return name
    }
    return ""
}

// İstifadə sırası — əvvəl header, sonra təxmin:
resp, err := http.Get(url)
defer resp.Body.Close()
enc := ctypeEncoding(resp.Header.Get("Content-Type"))
if enc != "" {
    fmt.Printf("Content-Type encoding is %s\n", enc)
    return
}
data, _ := io.ReadAll(resp.Body)
enc = dataEncoding(data)   // golang.org/x/net/html/charset
```
- Proqramın kənarlarında []byte alırsınız — string-ə çevirmək üçün
  encoding BİLİNMƏLİDİR
- Screen scraping son çarədir — API axtarın (köhnəlmiş sayt = qırılan kod)

### Recipe 23 — regexp ilə camelCase → snake_case
```go
var (
    // Example: match "rN" in "userName"
    camelRe = regexp.MustCompile(`[a-z][A-Z]`)
)

// fixCase gets a string in the format "aB" and return "a_b"
func fixCase(s string) string {
    return fmt.Sprintf("%c_%c", s[0], unicode.ToLower(rune(s[1])))
}

// camelToLower turns "camelCase" to "camel_case"
func camelToLower(s string) string {
    return camelRe.ReplaceAllStringFunc(s, fixCase)
}
```
- `MustCompile` — var blokunda error yoxlanıla bilmədiyindən; testlər
  xətanı production-dan əvvəl tutur
- `ReplaceAllStringFunc` — hər match-də funksiya çağırır
- Regex kənarları çoxdur — funksiya əsaslı substitusiya kodu sadələşdirir;
  regex101.com kimi alətlərlə test edin

### Recipe 24 — strings.EqualFold (case-insensitive)
**Tapşırıq:** "Gdańsk" = "gdańsk" axtarışı.

```go
func findTours(db []*Tour, city string) []*Tour {
    var tours []*Tour
    for _, t := range db {
        if strings.EqualFold(t.City, city) {   // Unicode-aware
            tours = append(tours, t)
        }
    }
    return tours
}
```
- ToLower/ToUpper yalnız ingiliscə üçün etibarlıdır; yunan sigma: Σ/σ/ς
- EqualFold: Unicode-fərqli, + yeni string YARADMIR → sürətli, yaddaş-
  qənaətli
- Artıq fold edilmiş saxlamaq lazımdırsa: `golang.org/x/text/cases`
  paketinin Caser-i

### Recipe 25 — Unicode normalizasiyası
**Problem:** "Kraków" ≠ "Kraków" — biri NFC (7 bayt), digəri NFD (8 bayt)!

```go
// Təsvir: len("Kraków") == 7, len("Kraków") == 8 — gözə görünməz fərq

// normString normalizes string in NFKC format
func normString(s string) string {
    return norm.NFKC.String(s)     // golang.org/x/text/unicode/norm
}

// Giriş nöqtələrində normallaşdır:
func NewTour(city, name string, time time.Time) *Tour {
    return &Tour{City: normString(city), Name: name, Time: time}
}

// Axtarışda da:
func findTours(db []*Tour, city string) []*Tour {
    city = normString(city)         // sorğunu da normallaşdır!
    // ... EqualFold müqayisəsi
}
```
- 4 normal form: NFC, NFD, NFKC, NFKD (Unicode TR15)
- **Qızıl qayda:** proqramın kənarlarında bir formaya normallaşdır,
  daxildə tək formada işlə
- ASCII → 7 bit; LATIN-1 → 8 bit; UTF-8 — bütün dillər (Rob Pike həmmüəllifi!)

## Final Thoughts-dən

Mətn universal interfeysdir; müasir emal Unicode bilik tələb edir. Go
standart kitabxanada + golang.org/x/text-də əla dəstək verir.

## Əsas terminlər
- Format verb (%v/%+v/%#v) — fmt/log çap əmrləri
- fmt.Stringer — String() metodu ilə xüsusi çap
- Sonsuz rekursiya təhlükəsi — String daxilində %v
- mime.ParseMediaType — Content-Type parse
- charset.DetermineEncoding — encoding təxmini
- regexp.MustCompile — var-də panic edən kompilyasiya
- ReplaceAllStringFunc — funksiya əsaslı əvəzetmə
- strings.EqualFold — Unicode case-insensitive müqayisə
- Unicode normalizasiyası (NFC/NFD/NFKC/NFKD) — ekvivalent formaların
  vahidləşdirilməsi
- ASCII/LATIN-1/UTF-8 — tarixi encoding sxemləri

## Praktik nəticə
Debugging üçün %#v (tiplər görünür!); istifadəçiyə görünən çap üçün
Stringer implement edin (rekursiyaya diqqət). []byte→string keçidində
əvvəl Content-Type charset-i, yoxsa charset.DetermineEncoding. Case-
insensitive müqayisə — ToLower yox, EqualFold; çoxdilli sistemlərdə bütün
girişləri NFKC (və ya tək form) ilə normallaşdırın — yoxsa "görünməz"
fərqlər axtarışları səhvə salacaq.

## Mənbə
Pages: 58-72 (PDF 58-72)
