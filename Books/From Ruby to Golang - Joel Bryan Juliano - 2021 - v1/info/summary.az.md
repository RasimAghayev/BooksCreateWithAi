# From Ruby to Golang — Xülasə (Azərbaycanca)

> **Kitab:** From Ruby to Golang: A Ruby Programmer's Guide to Learning Go — Joel Bryan Juliano, Leanpub, 2021 (ISBN 978-1080944002, 160 səh., 2021-06-10 buraxılışı)
> **Səviyyə:** 🌱 Elementary (2/5) · **Dillər:** Ruby → Go
> **Xüsusiyyət:** Hər mövzu Ruby nümunəsi ilə AÇILIR, sonra Go tərcüməsi gəlir. Hər fəslin sonunda "Chapter Questions" (özünüyoxlama) bölməsi var.

---

## Kitabın məqsədi və yanaşması

Müəllif — 8+ illik Ruby təcrübəsi olan mühəndis — 2018-də Go istifadə edən şirkətə keçəndə öyrənmə qeydlərini kitaba çevirib. **Metod: analogiya** — bildiyin Ruby anlayışını Go-yə xəritələmək. Kitab Rubyist üçün yazılıb, amma strukturu istənilən dinamik-dil istifadəçısına (Python/JS) uyğun gəlir.

## Chapter-by-chapter xülasə

### Ch 1-2 — Structs: instans dəyişənindən struktura
Ruby `@var` → Go struct sahələri; `class` + metod → struct + receiver funksiyaları. Private/public = kiçik/böyük hərf. **Pointer receiver = attr_accessor** (setter+getter, orijinalı mutasiya); **value receiver** — kopya, thread-safe. İrs = embedding (`*Animal`, çoxlu irs OK). Anonymous struct — inline, adsız; anonymous sahələr — tip adı ilə çıxış.

### Ch 3 — Hash/Map + Variadic
Ruby Hash → `map[K]V`. İki ilkinləşdirmə ailəsi: declaration (var+make/literal) və assignment (`:=` bir sətir). `make` hint size. `map[K][]V` array dəyərlər. `interface{}` açarlar (duck-typing — pulu var, ucuz alternativ anonymous struct). `delete` təhlükəsiz. Mövcud olmayan açar → zero value (0/""/false/nil); comma-ok ayırıcı; `m[k]++` sayğaclar üçün. `**kwargs` → `...interface{}` + `args[0].(map[string]string)` assertion; `%T` = `is_a?`.

### Ch 4 — Arrays/Slices + naviqasiya
İki klasifikasiya: **Fixed Array** (ölçü tipin hissəsi, `[...]` avto-ölçü, value-copy assign — böyük datasda bahalı) və **Sliced Array** (referans — 4 dəyişən eyni yaddaş!). `make(T, len, cap)` — capacity slicing-lə aktivləşir (`s[:3]`); ortaq yaddaş tələsi (basket1/2/3 eyni array). **Deep copy = make+copy.** append capacity aşımında böyüyür. Splat → `...T`; `basket(slice...)` açılışı. `[]interface{}` qarışıq massivlər. Ruby `each` → for range-in 5 semantik forması (C-style, index+value, muted `_`, index-only, pointer — pis praktika).

### Ch 5 — Package Management
RubyGems/Bundler → go get/Go Modules. GOPATH ənənəvi; modullar GOPATH-dən azad edir. Ruby module+include → Go package+import (export = böyük hərf). go.mod: module + require; manual, `go mod -sync` (import-lardan) və dep-dən generasiya. `// indirect` sub-paket həlli (versiya tapmaq: curl+jq). Legacy dep: `dep init/ensure`, Gopkg.toml (TOML `[[constraint]]`: name/version/branch/source-fork).

### Ch 6 — Enumerable → Go patternlər
Ruby builtin-ləri Go-da YOXDUR — manual loop patternlər: **all?** (bool massivi + hamısı yoxlaması), **any?** (bir true tapanda break), **collect** (nəticə massivi), **cycle** (nested loop), **detect** (şərt+break), **drop** (slicing bir sətirdə!), **drop_while** (şərt+kəsmə). Dərs: kolleksiyaya ümumi sual verə bilmək.

### Ch 7 — Interface
Ruby module-mixin → Go interface. 3 forma: **(1) Self-documenting API** — `ProduceBasket` interfeysi AddItem/RemoveItem/Items; `var b ProduceBasket = &Basket{}`; **(2) Type contract** — `[]ProduceBasket` heterojen kolleksiyası (fruits+veggies); **(3) Return value contract** — `Produce` interfeysi: Item struktur Flavour()/Kind() qaytarır. Slice-dan silmə idiomu: `append(s[:i], s[i+1:]...)`.

### Ch 8 — Glossary (112-158)
46 səhifəlik ensiklopedik lüğət — simvollar (`%+v`, `&&`, `$GOPATH`), operatorlar, proqramlaşdırma anlayışları (Wiktionary mənbəli). Kitabın 30%-i — istinad xarakterli.

---

## Kitabın əsas mesajları

1. **Analogiya keçidin sürətləndirir:** Ruby bildiklərin — instans dəyişəni, module, Hash, each, splat, mixin — Go-da struct, package, map, range, variadic, interface qarşılıqları var.
2. **Value vs reference semantikası kritikdir:** Fixed array kopya, slice referans; value receiver təhlükəsiz kopya, pointer receiver birbaşa mutasiya.
3. **Zero value Go-nun sürprizsizlik prinsipidir:** Mövcud olmayan açar 0/""/nil — comma-ok ilə dəqiqləşdir.
4. **Duck-typing Go-da mümkündür, amma PULLU:** `interface{}` + type assertion; tipli alternativlər (anonymous struct, interface) daha idiomatikdir.
5. **Enumerable öz looplarınla qurulur** — və ya hazır patternlərdən (gobyexample) götürülür.
6. **Paket idarəsi ekosistem fərqi:** Gemfile→go.mod; module+include→package+import.
7. **Kitabın dəyəri:** keçid dövründə "Ruby-də bu necə idi?" sualına dərhal cavab + başlanğıc səviyyəli qısa təriflər.

## Kitabdan sonra

- Effective Go (idiomatik dərinlik) —本书-dən təbii növbəti addım
- Go-də sinonim olmayan anlayışlar: error handling, defer, goroutine/channel — kitab bu mövzulara GİRİRMİR (və bu, boşluqdur)
- Go Modules müasir axını (kitabın dep bölməsi tarixi maraqdır)
