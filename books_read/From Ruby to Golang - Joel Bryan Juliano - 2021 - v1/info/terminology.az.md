# From Ruby to Golang — Terminologiya (Azərbaycanca)

> Kitab boyu rast gəlinən terminlər. Ruby mənşəli terminlər Go qarşılığı ilə birgə.

## A

**Anonymous struct** — Adsız inline struct: `struct{...}{values}` — tez çağırış, namespace çirklənməməsi; `interface{}`-ə ucuz alternativ.

**Anonymous struct field** — Yalnız tipdən ibarət sahə (`int`, `string`); ad = tip adı; hər tip yalnız 1 dəfə.

**any? (Ruby)** — Predicate metod: hər hansı element şərti qarşılayırsa true.

**all? (Ruby)** — Predicate metod: HAMISI qarşılayırsa true.

**Array (fixed)** — `[N]T` — sabit ölçülü, ölçü TİPİN hissəsi, value type (assign = tam kopya).

**Assignment (map)** — `m := map[K]V{...}` — bir ifadədə yarat+ilkinləşdir.

**Attribute accessor (Ruby)** — `attr_accessor` — həm oxu həm yazı; Go: pointer receiver.

## B

**Big O notation** — Komplekslik təsviri; array O(n), map O(1).

**Bundler / Gemfile (Ruby)** — Asılılıq idarəsi; Go: Go Modules / go.mod.

## C

**Capacity (make 3-cü parametr)** — Slicing-dən sonra qalan tutum; `make(T, len, cap)`.

**Comma-ok** — `v, ok := m[k]` — mövcudluğu zero-value-dan ayıran forma.

**collect/map (Ruby)** — Transformasiya nəticələrini yeni massivə yığır.

**cycle (Ruby)** — Kolleksiyanı n dəfə təkrarla.

## D

**Declaration (map)** — `var m map[K]V` + sonra `make`/literal — allocate sonraya.

**Deep copy** — `make` + `copy` — orijinaldan müstəqil nüsxə.

**delete(m, k)** — Map açarı silmə; mövcud olmasa təhlükəsiz.

**detect/find (Ruby)** — İlk uyğun element; Go: şərt + break.

**Double-splat `**kwargs` (Ruby)** — Keyword arqumentləri hash-ə; Go emulyasiyası: `...interface{}` + assertion.

**drop / drop_while (Ruby)** — İlk n-i at / şərt true olduqca at; Go: `a[n:]` slicing.

**Duck typing** — "Gədirsə, ördəkdir" — tip yoxlamadan davranış; Go: `interface{}` (və ya anonymous struct).

## E

**each (Ruby)** — İterasiya; Go: `for range` 5 semantik forma.

**Embedding (struct)** — `*Parent` daxilində — irs qohumu; sahə promote; çoxlu embed OK.

**Enumerable (Ruby)** — Kolleksiya builtin-metodlar ailəsi; Go-da manual.

**Enumerator** — İterasiya obyekti.

## G

**GOPATH** — Paket/kitabxana axtarış yolu; `$GOPATH/src`; Go Modules ilə arxa plana keçir.

**go.mod** — Modul faylı: module adı + require-lər (Gemfile analoqu).

**go get** — Paket yükləmə (RubyGems-in əməli funksiyası).

**GO111MODULE** — Modul rejimi keçidi: on/off/auto.

**Gopkg.toml / Gopkg.lock** — dep konfiq/kilid faylları; TOML Array of Tables.

## H

**Hash (Ruby)** — `map[K]V`-nin qohumu; `{}` və ya `Hash.new`.

## I

**instance variable (Ruby)** — `@var` — sinif daxili tək-istinad dəyişən; Go: struct sahəsi.

**interface{} (empty)** — Hər tip qəbul edir — dinamik tiplər; tip itkisi qiyməti var.

**Interface (Go)** — Metod dəsti müqaviləsi; implicit təmin; module+include-ün Go tərcüməsi.

**// indirect** — go.mod-da sub-paket asılılıq işarəsi.

## L

**Literal type assignment** — `var m = map[K]V{...}` — elan+ilkinləşdirmə bir sətirdə.

## M

**map[K][]V** — Array dəyərli map — bir açara çox dəyər.

**make hint size** — Map üçün ölçü tövsiyəsi — performans; avtomatik idarə olunur.

**mixin (Ruby)** — `include` ilə modul əlavəsi; Go: interfeyslər + embedding.

**module (Ruby)** — Metod namespace-i; Go: package.

## N

**Namespace pollution** — Lazımsız tip adlarının çirklənməsi — anonymous struct ilə qarşısı.

## O

**O(n) / O(1)** — Ardıcıl emal / bir addımlı çıxış (array/map).

**Operator: `...T` vs `x...`** — Variadic parametr / çağırışda array açılışı.

## P

**Package (Go)** — Funksiya/struct namespace-i; böyük hərf = export.

**Pass-by-value / Pass-by-reference** — Dəyər kopyası / yaddaş ünvanı ötürməsi.

**Pointer receiver** — `(p *T) m()` — orijinalı mutasiya edir; attr_accessor ekvivalenti.

**Predicate method (Ruby)** — Boolean qaytaran metod (all?, any?).

**Private/Public (Go)** — Kiçik/BÖYÜK ilk hərf — paket-daxili / export.

**ProduceBasket pattern** — Kitabın nümunə interfeysi: AddItem/RemoveItem/Items.

## R

**Reference type** — Slice — assign = eyni yaddaşın paylaşımı (fixed array-in əksi).

**Receiver** — `func (p T) m()` — struct-a bağlı funksiya = metod.

**RubyGems** — Ruby paket meneceri; go get qohumu.

## S

**Self-documenting API** — İnterfeysin funksiya siyahısı sənəd kimi.

**Sliced Array** — `[]T` — dinamik, referans; standart API-nin hamısı.

**Slicing `a[:3]`** — Pəncərə: eyni yaddaş + capacity-nin aktivləşdirilməsi.

**Splat `*args` (Ruby)** — Catch-all arqument; Go: `...T`.

**Struct** — Sahələrin deklarativ qrupu — instans dəyişəninin Go ekvivalenti.

## T

**TOML** — Gopkg.toml formatı; `[[constraint]]` Array of Tables.

**Type assertion** — `x.(T)` — interface{}-dən konkret tip çıxarma.

**Type contract** — İnterfeyslə heterojen tiplərin eyni kolleksiyada birləşməsi.

## V

**Value receiver** — `(p T) m()` — kopya ilə işləyir; thread-safe.

**Value type** — Fixed array: assign = tam element kopyası.

**Variadic function** — `func f(args ...T)` — dəyişən saylı arqument.

**vendor/** — Asılılıqların layihədaxili nüsxəsi (dep).

## Z

**Zero value** — Mövcud olmayan map açarının tip defaultu: false/0/0.0/""/nil.
