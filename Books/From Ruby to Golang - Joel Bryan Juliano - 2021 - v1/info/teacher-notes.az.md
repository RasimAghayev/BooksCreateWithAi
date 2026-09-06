# From Ruby to Golang — Müəllim Qeydləri (Teacher Notes)

> Bu sənəd kitabı öyrədərkən istifadə üçün metodiki qeydlər, çətin anlar, müzakirə sualları və praktik tapşırıqları birləşdirir.

---

## 1. Auditoriya və ön şərtlər

- **Hədəf auditoriya:** Ruby (və ya digər dinamik tip dil — Python/JS) bilən, Go-ya keçən developer-lər. Tam başlanğıc üçün YOX — Ruby nümunələri ön-şərt kimi oxunur.
- **Ön şərtlər:** Ruby əsasları (class, module, Hash, Array, block), ümumi proqramlaşdırma anlayışları.
- **Həcm:** 160 səh. (46-sı lüğət) — real öyrəniləcək məzmun ~114 səh. — 8-12 saatlıq kurs və ya 2 günlük intensiv üçün ideal.

---

## 2. Tədris axını

| Module | Chapter | Süret | Fokus |
|--------|---------|------|-------|
| 1 | Ch 1-2 | 25% | Struct + receiver semantikası (ƏSAS module) |
| 2 | Ch 3 | 20% | Map + variadic + zero value |
| 3 | Ch 4 | 20% | Array/slice value-vs-reference (ƏN çətin module) |
| 4 | Ch 5 | 10% | Paket idarəsi (nümayiş yönümlü) |
| 5 | Ch 6 | 10% | Enumerable patternlər |
| 6 | Ch 7 | 15% | Interface 3 forması |
| — | Ch 8 | — | Lüğət — ev tapşırığı/istinad |

**Dərs formatı (kitabın öz quruluşu):** Hər mövzu = Ruby nümunə → Go tərcümə → Chapter Questions. Dərsdə eyni ritmi saxlayın: (1) Ruby versiyanı birgə oxuyun; (2) "Go-da bu necə olardı?" müzakirəsi; (3) Go kodu; (4) suallar.

---

## 3. Çətin anlar və izah üsulları

### a) Value vs Reference (Ch 2 və 4 — iki dəfə qayıdılır, ən çətin mövzu)
**Vizual alət:** Whiteboard-da 2 diaqram: (1) fixed array assign → İKİ ayrı qutular (kopya); (2) slice assign → İKİ ox, BİR qutu. basket1/2/3 nümunəsini canlı işə salın — "mən basket3-ü dəyişdim, nə üçün basket1 də dəyişdi?!" sualı dərsin ən yaddaqalan anı olacaq.

### b) Pointer receiver nə vaxt? (Ch 2)
**Ruby körpüsü:** `attr_accessor` istəyirsən → pointer; yalnız oxuyursan → value. change_item tələsini müzakirə edin: value receiver slice ELEMENTLƏRİNİ dəyişə bilir, amma receiver-in ÖZÜNÜ (uzunluq) yox — "mövcud elementləri redaktə edir, amma yeni sətir əlavə edə bilmir."

### c) Zero value + comma-ok (Ch 3)
"Ruby-də h[?]=nil alırsan, Go-da 0 — amma 0 QANUNİ dəyər ola bilər!" — comma-ok-un zorunluluğu buradan doğur. `basket["cabbage"]++` sayaclar üçün niyə gözəldir (mövcud olmayan açar 0-dan başlayır).

### d) make(type, len, cap) (Ch 4)
Kitabın ən yaxşı izahlarından: `make([]string, 2, 3)` + `array[:3]` — "capacity = slicingdən sonra qalan yer". 3 addımlı canlı demo: [2]-yə yazma xətası → slicing → yazma OK.

### e) `**kwargs` → `...interface{}` (Ch 3)
Bu, kitabın ən "dünyəvi" həlləridir — production-da struct arqumenti daha yaxşıdır. Müzakirə sualı: "bu kod real layihədə hansı problemlər doğurar?" (tip itkisi, runtime panic, oxunmazlıq). Ardıcıllıq üçün öyrədin, amma idiomatik alternativi GÖSTƏRİN.

### f) Interface implicit təmin (Ch 7)
Ruby mixin-dən gələn gözlənti: "harada include yazılır?" — Go-da YER YOXDUR. `var b ProduceBasket = &Basket{}` sətri təmin edir. `[]ProduceBasket` heterojen kolleksiyası polimorfizmin "a ha!" anıdır.

---

## 4. Müzakirə sualları

1. change_item value receiver-dır, amma basket elementlərini dəyişir — bu, setterdir? (element: bəli; struktura: xeyr)
2. Niyə `[4]string ≠ [6]string`? Ruby-də massiv ölçüsü tipdirmi? (dinamik dildə yox — statik tipin qiyməti)
3. `interface{}` açarlar rahatdır — hansı hallarda PUL ödənilir? (runtime assertion, panic, tip səhvləri)
4. drop = `a[n:]` bir sətirdir — Ruby drop_while Go-da nə üçün bir sətir deyil? (lazy şərtlənmə)
5. dep nə üçün öldü? (rəsmi Modules standartı) — Gemfile.lock-a bənzər nə var? (go.sum)
6. Kitab hansı mövzuları KƏSİB? (error handling, defer, goroutine — ən vacib boşluqlar; müzakisə dərsi bitirən sual)

---

## 5. Praktik tapşırıqlar

**Tapşırıq 1 (struct):** Ruby Basket sinifini Go struct+receiver-a tərcümə edin — add_item (pointer), items (value). `%-v` ilə çap edin.

**Tapşırıq 2 (embedding):** Animal+Owner+Dog üçlüsünü qurun; birdən çox irs sınaqına baxın — Ruby-də eyni mümkündürmü?

**Tapşırıq 3 (map):** 3 üsulla map yaradın (make/literal/empty); struct dəyərli map; `map[string][]string`; comma-ok ilə mövcudluq; delete.

**Tapşırıq 4 (slice tələləri):** basket1→basket2→basket3 zənciri qurun, hamısının dəyişdiyini görün; deep copy ilə müstəqillik qazanın; append-lə cap aşımını sınayın.

**Tapşırıq 5 (variadic):** `sum(nums ...int)` yazın; mövcud slice ötürün (`nums...`); `...interface{}` qarışıq arqumentləri `%T` ilə filtrləyin.

**Tapşırıq 6 (enumerable):** `All`, `Any`, `Detect` funksiyalarını GENERIC şəkildə yazın (müasir Go) — kitabın konkret-tipli versiyaları ilə müqayisə edin.

**Tapşırıq 7 (package):** coffee paketi + go.mod qurun; module adını github.com şəklinə çevirin; `// indirect` ssenarisini yaradın.

**Tapşırıq 8 (interface):** ProduceBasket interfeysi + 2 fərqli Basket tipi; `[]ProduceBasket` polimorfik loop; Produce return-contract nümunəsi.

---

## 6. Sınav sualları (nümunə)

**Asan:**
1. `type dog` ilə `type Dog` fərqi? (private/public)
2. `m["yoxdur"]` nə qaytarır int map-də? (0)
3. `[...]string{...}` nə edir? (ölçünü literal sayından alır)

**Orta:**
4. Pointer receiver nə üçün attr_accessor-a bənzərdir?
5. `make([]string, 2, 3)`-də [2] yazmaq niyə xətadır və necə düzəldilir? (len 2; `s[:3]`)
6. `func f(a ...int)` funksiyasına slice ötürmək? (`f(slice...)`)

**Çətin:**
7. Value receiver change_item nəyi dəyişə BİLİR, nəyi YOX? (elementləri bəli; strukturu/uzunluğu yox)
8. basket3 slicing-dən sonra basket1 niyə dəyişdi? (eyni backing array — referans paylaşımı)
9. `[]ProduceBasket`-ə fərqli tiplər niyə qoşula bilir? (hamısı interfeysi implicit təmin edir)

---

## 7. Kollektiv layihə ideyası

**"Ruby-to-Go Rosetta Stone"** — tələbələr Ruby layihəsindən (sadə blog/basket app) bir modulu qrupla Go-ya köçürürlər. Hər çətinlikdə kitabın xəritəsindən istifadə: class→struct, each→range, Hash→map, module→package. Təqdimat: hər qrup öz "keçid qeydləri"ni (analogiya + sürpriz anları) təqdim edir.

---

## 8. Əlavə resurslar

- Kitabın səhifəsi: leanpub.com/rb2go
- Go By Example — collection-functions (kitabın tövsiyəsi): gobyexample.com
- A Tour of Go: go.dev/tour
- Effective Go (növbəti addım): go.dev/doc/effective_go
- Go Modules rəsmi: go.dev/ref/mod
- dep (arxivləşdirilib): github.com/golang/dep
- Ruby doc — double-splat: ruby-doc.org (kitabın istinadı)
