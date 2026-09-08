# Chapter 2 — Hello, earth! (səh. 46-78)

## Bu fəsil nədən bəhs edir?

İlk Go proqramı və onun ətrafında BÜTÜN əsas praktikalar: go mod init (module
= kod repozitoriya yolu), package main + fmt.Println ("print with new line"),
adlandırma qaydaları (qısa scope = 1-2 hərf; camelCase; Hungarian notation
YOX; rəqəmlə başlaya bilməz), BÖYÜK hərf = exposal qaydası (dəyişən, konstant,
funksiya, tip — hamısı), tab-indentasiya, Example funksiyası (stdout yoxlaması
+ go doc sənədi; "// Output:" şərhi = test runner; olmazsa yalnız compile),
go test (PASS/ok modul adı), Dijkstra ("testing proves presence of bugs,
not absence"), internal vs external test (main_internal_test.go — package
main = unexposed-ə giriş; {pkg}_test = istifadəçi nöqteyi-nəzəri), 4 fazalı
test strukturu (Preparation → Execution → Decision → Teardown; defer ilə),
want/got konvensiyası + t.Errorf (%q formatı), refactoring greet çıxarışı
(funskiya dar scope = testable/debuggable/explicit), Test{Func}_{Scenario}
adlandırma (unsupported Akkadian halı — safety net testi), switch (implicit
break; error/pointer/bool istisna), type language string ( CLARITY THROUGH
TYPING — parametr qarışmasını maneə törədir), map = hash table (phrasebook;
multi-literal Unicode dil salamı), comma-ok idiomu (greeting, ok :=
phrasebook[l]; zero value təhlükəsi!), multiple return 4 halı (map, range,
channel <-, funksiya), table-driven test (testCase struct + map[string]
testCase + t.Run(name, ...) — subtest adları; slice/map variantları), 3
sitat növü ("..." interpritasiyalı; `...` raw/escape-siz — JSON üçün ideal;
'...' rune = Unicode code point), flag package (os.Args positional YOX;
StringVar(&lang, ad, default, help) + flag.Parse; String = pointer qaytarır;
"Sunset begins at the Parse-fic Ocean" mnemonik; & adres / * indireksiya;
copy-by-value → pointer dəyişiklik üçün; pointer arifmetikası YOX), CLI
testi (go run main.go -lang=el).

## Əsas fikirlər

### 1. İki Sətirlik Proqramın Anatomiyası
```go
package main        // hər fayl paket adı ilə BAŞLAYIR; main = xüsusi
import "fmt"        // URL-ə oxşamayan import = STANDART KİTABXANA
func main() {       // main() = proqramın GİRİŞ nöqtəsi
    fmt.Println("Hello world")   // Println = "print with new line"
}
```
- **package main iki xüsusiyyəti:** (1) qovluq adı ilə uyğun gəlmir;
  (2) main() funksiyasının yerini compiler-a bildirir
- **Böyük hərf qaydası:** Println böyük P — packaədən XARİCİ giriş;
  kiçik hərf/underscore = gizli (dəyişən/konstant/funskiya/tip hamısına)

### 2. Example Funksiyası — stdout Testi + Sənəd
```go
func ExampleMain() {
    main()
    // Output:
    // Hello world
}
```
- **Ad konvensiyası:** Example<FunctionName> — stdout yoxlamaya ELİGİB
- **"// Output:" şərhi:** ondan sonrakı sətirlər = gözlənilən çıxış; ŞƏRH
  YOXDURSA → compile olunur amma İCRA OLUNMUR (yalnız go doc sənədi)
- **Niyə Example:** stdout yazan funksiyalar nadirdir; əsas test vasitəsi
  bu DEYİL — amma ilk funksiya üçün asan yol

### 3. 4 Fazalı Test Strukturu
```go
func TestGreet(t *testing.T) {
    want := "Hello world"        // 1. PREPARATION
    got := greet()               // 2. EXECUTION
    if got != want {             // 3. DECISION
        t.Errorf("expected: %q, got: %q", want, got)
    }
}                                 // 4. TEARDOWN — defer ilə (bu halda lazımsız)
```
- **want/got** — Go konvensiyası; t.Errorf = FAIL işarəsi (Errorf/Printf
  imzası)
- **Dijkstra:** "Testing can prove the presence of bugs, but not their
  absence!" — çox test = etibarlı kod
- **Internal test faylı:** main_internal_test.go — package main (gizli
  funksiyalara giriş); əksər testlər main() DEYİL, onun ÇAĞIRDIĞI
  funksiyalardadır

### 4. Test Adlandırma — Ssenari Bölgüsü
```go
func TestGreet_English(t *testing.T) { lang := language("en"); ... }
func TestGreet_French(t *testing.T)  { lang := language("fr"); ... }
func TestGreet_Akkadian(t *testing.T) { want := "" ... }   // SAFETY NET!
```
- **Test{Func}_{Scenario}** konvensiyası; **unsupported input testi** =
  "good" input-dən daha DƏYƏRLİ — safety netlərin yoxlanması

### 5. type language string — Aydınlıq Tipi
```go
type language string
func greet(l language) string
```
- **Clarity through typing:** string/int/index URL qarışması İMKANSIZlaşır;
  adi string konstantı avtomatik language-ə çevrilir ("en" imzasına görə)
- **switch:** case-lər arası break İMLSİDIR (error/pointer/bool istisna);
  if yalnız 1-2 hal üçün — artıq olanda switch/map

### 6. map — Hash Table
```go
var phrasebook = map[language]string{
    "el": "Χαίρετε Κόσμε",  // Yunan
    "en": "Hello world",
    "fr": "Bonjour le monde",
    "he": "םלוע םולש",      // İvrit
    ...
}
greeting, ok := phrasebook[l]     // COMMA-OK idiomu
if !ok {
    return fmt.Sprintf("unsupported language: %q", l)
}
```
- **Comma-ok:** map çıxışı 2 dəyər — (value, found); ok yoxlanılmazsa
  zero value ("") qayıdır və TAPILMADIQLIĞI BİLMƏMİK olur
- **Multiple return-ın 4 ümumi halı:** (1) map[key] → (v, ok); (2) range →
  (key, value); (3) channel <- → (v, closed); (4) adi funksiya (ən çox:
  errors are values)
- **Qlobal dəyişən:** bu layihədə qəbul olunur (production-da çox vaxt
  anti-pattern); switch-i map ilə əvəz = funksiya QISALIR

### 7. Table-Driven Test (TDT)
```go
func TestGreet(t *testing.T) {
    type testCase struct {
        lang language
        want string
    }
    var tests = map[string]testCase{       // slice də OLAR
        "English":                {lang: "en", want: "Hello world"},
        "French":                 {lang: "fr", want: "Bonjour le monde"},
        "Akkadian, not supported": {lang: "akk", want: `unsupported language: "akk"`},
        "Empty":                   {lang: "",   want: `unsupported language: ""`},
    }
    for name, tc := range tests {
        t.Run(name, func(t *testing.T) {      // SUBTEST = adlı nəticə
            got := greet(tc.lang)
            if got != tc.want {
                t.Errorf("expected: %q, got: %q", tc.want, got)
            }
        })
    }
}
```
- **Motivasiya:** hər ssenari EYNİ cəmi kod — yalnız 2 sətir dəyişir;
  yeni hal = map-ə 1 ENTRY (10 sətir YOX)
- **t.Run faydaları:** fail adı görünür; editorda tək hal İŞLƏDİLİR
- **Akkadian/Empty halları:** dəstəklənməyən dillər DƏ test olunmalıdır

### 8. Sitat İşarələri — 3 Növ
| İşarə | Ad | Xüsusiyyət |
|---|---|---|
| "..." | literal string | escape sekwensiyalar İŞLƏYİR (\n) |
| \`...\` | raw literal | escape YOXDUR — \n literal olaraq qalır; JSON payload üçün ideal |
| '...' | rune | TƏK Unicode code point ('學') |

### 9. flag Package
```go
var lang string
flag.StringVar(&lang, "lang", "en", "The required language, e.g. en, ur...")
flag.Parse()                     // BUDUR parsing edir!
greeting := greet(language(lang))
```
- **os.Args vs flag:** positional + manual parse vs type-hazır (int/float/
  duration/string/bool)
- **StringVar parametrləri:** (&deyisen, "ad", default, "help") — 1) adres
  (& — copy-by-value qaydasına görə dəyişmək üçün POINTER lazımdır);
  2) komanda sətri adı; 3) default; 4) description
- **flag.String alternativi:** pointer YARADIB qaytarır
- **MNEMONİK:** "Sunset begins at the Parse-fic Ocean" — *Var funksiyaları
  doldurmur, flag.Parse doldurur!
- **& / * operatorları:** adres / indireksiya; pointer arifmetikası YOX
  (array-ın 2-ci elementinə pointer-dan keçmək OLMAZ)
- **Validation production-da:** supported dillər siyahısına və ya formatına
  (2 ASCII hərf) yoxlama burada olmalıdır

### 10. CLI Sınaq + Side Quests
```bash
go run main.go -lang=el    # Χαίρετε Κόσμε
```
- Side quests: Urdu/destəksiz/boş/default dilləri sına; ÖZ dilini əlavə et

## Əsas terminlər
- Module (go mod init) — asılılıq + paket konteyneri; ad = repo yolu
- package main — icra edilən proqramın paketi; main() girish nöqtəsi
- Exposal (böyük hərf) — xarici görünənlik; kiçik hərf = gizli
- Example funksiyası — stdout testi + go doc nümunəsi
- Internal/external test — package main / {pkg}_test faylları
- 4 faza — preparation/execution/decision/teardown
- want/got — gözlənilən/alınan konvensiyası
- Comma-ok — map/kanal çıxışının (v, bool) idiomu
- Zero value — tipin default (string = ""); map miss gizli təhlükə
- Table-driven test — testCase struct + map/slice + t.Run
- Subtest — t.Run(name, ...) adlı hal
- Raw literal (\`) — escapesiz string; rune (') — code point
- flag.Parse — *Var-ları DOLDURAN yeganə funksiya
- Safety net testi — unsupported/hata hallarının testi

## Praktik nəticə

1. **Hər proqramın şablonu:** mod init → package main → main() → funksiya
   çıxarışı (dar scope) → test faylı (internal) → 4 faza → want/got.
2. **TDT-ə TEZ keç:** 2+ eyni-quruluşlu test funksiyası görəndə cədvələ
   çevir; map[string]testCase + t.Run.
3. **map çıxışında comma-ok MÜTLƏQ:** zero value səssiz qayıdır.
4. **Custom tip = sənəd + təhlükəsizlik:** type language string parametr
   səhvlərini compile-a çevirir.
5. **flag.Parse-i unutma:** *Var yalnız QEYD edir; dəyər Parse-dan sonra.
6. **Testlərdə safety netlərə öncəlik ver:** "happy path" qədər unsupported
   hallar da yoxlanılır.
7. **Raw string JSON üçün:** `{"key": "value"}` — dırnaq escape problemindən
   azad.

## Mənbə
Pages: 46-78 (PDF 47-79)
