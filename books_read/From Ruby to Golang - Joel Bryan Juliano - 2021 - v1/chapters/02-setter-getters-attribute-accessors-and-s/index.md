# Chapter 2 — Setter/Getters, Attribute Accessors və Structs (book səh. 3-28)

## Bu chapter nədən bəhs edir?

Ruby instans dəyişənlərinin (@var), attr_accessor-un və irs (inheritance) anlayışlarının Go struktur (struct) dünyasına tərcüməsi: private/public struct, funksiyaya struct bağlama (metodlar), pointer vs value receiver, struct irsi (embedding), anonymous struct və sahələr.

---

## əsas fikirlər

### 1. İnstans dəyişəni (Ruby @) → struct (Go)

**Ruby:** `@breed` — sinif daxilində tək istinadla çatılaşan dəyişən; initializer ilə doldurulur, `kind` metodu ilə oxunur.

**Go:** struct — sahələrin deklarativ qrupu, tək istinadla çatılaşır. **"Instans dəyişəninin Go ekvivalenti = struct"dır.**

```go
// PRIVATE struct (birinci hərf kiçik — YALNIZ paket daxilində görünür):
type dog struct {
    breed string
}

kind := dog{breed: "Rottweiler"}
fmt.Println(kind.breed)   // Rottweiler
```

**Public struct — BÖYÜK hərf:**

```go
type Dog struct {
    Breed string      // sahə də böyük — hər ikisi xaricdən görünür
}
kind := Dog{Breed: "Rottweiler"}
```

**Qayda:** Böyük hərf = public (struct, sahə, funksiya, dəyişən — hamısı üçün). Çoxsahəli struct: `name`, `breed`, `age` + `fmt.Printf("%+v", pet)` → `{name:Maximus breed:Rottweiler age:5}`.

### 2. Funksiyaya struct bağlama — metod (Ruby sinif metodu → Go receiver)

**Ruby:** Basket sinifi — @produce instans dəyişəni + `add_item`/`change_item`/`items` metodları.

**Go:** Funksiyanın receiver-i olaraq struct tipi təyin edilir:

```go
type basket []produce
type produce struct {
    name    string
    flavour string
    kind    string
}

func (p *basket) add_item(entry produce) {     // POINTER receiver
    *p = append(*p, entry)
    fmt.Printf("Entry %s created!\n", entry.name)
}

func (p basket) change_item(name string, entry produce) {   // VALUE receiver
    for key, val := range p {
        if val.name == name {
            p[key] = entry     // mövcud elementi dəyişir
        }
    }
}

func (p basket) items() {
    fmt.Printf("There are %d items\n", len(p))
    for _, val := range p { /* çap */ }
}

func main() {
    basket := new(basket)
    basket.add_item(produce{name: "apple", ...})
    basket.items()
}
```

**Mahiyyət:** Go-nun birinci sinif vətəndaşları FUNKSİYALAR-dır (sinif YOXDUR) — struct-ı funksiyaya "bağlayaraq" metod effekti alınır. Bu, Ruby-dəki dot-notasiya çağrışının (basket.add_item) ekvivalentidir.

### 3. Pass-by-value vs Pass-by-reference

**Ruby özü pass-by-value-dır** (obyekt referansının dəyəri kopyalanır). Go-da İKİ semantika:

| Receiver | Keçirmə | Davranış |
|----------|---------|----------|
| **Value receiver** `(p basket)` | Pass-by-value | Hər çağrışda YENİ MÜSTƏQİL KOPYA yaranır — orijinala təsir yoxdur |
| **Pointer receiver** `(p *basket)` | Pass-by-reference | Yaddaş ÜNVANI ötürülür — dəyişikliklər ORİJİNALA yazılır |

**Pointer receiver həm setter həm getterdir** — Ruby `attr_accessor`-un ekvivalenti:

```go
func (receiverName *receiverType) funcName(paramName paramType) {
    *receiverName = paramName    // SET — birbaşa mutasiya
    fmt.Println(receiverName)     // GET
}
```

**Value receiver nüansı:** Mövcud elementləri dəyişə bilir (slice elementləri ünvanla), amma **receiver struktunun ÖZÜNÜ birbaşa mutasiya edə bilmir** — yeni kopyada işləyir.

**Thread-safe:** Value receiver konkurent tətbiqlər üçün idealdir — orijinalı dəyişmədiyindən kopya ilə işləmək təhlükəsizdir. (Rəsmi semantika üçün: Go wiki CodeReviewComments#receiver-type.)

### 4. Struct irsi — embedding (Ruby inheritance → Go)

**İrs üçün strukturları kiçik hissələrə ayır** (modularity → reuse → maintainability). Struct daxilində digər struct adı = İRS:

```go
type Animal struct {
    Kind string
    Diet  string
}
type Owner struct {
    Name    string
    Country string
}

type Dog struct {
    *Animal    // POINTER embed — sahələrə birbaşa çıxış
    *Owner
    Name  string
    Breed string
}

var dog Dog
dog.Name = "Maximus"
dog.Animal = &Animal{Kind: "Dog", Diet: "Omnivorous"}   // daxili struct ayrıca ilkinləşdirilir
dog.Owner = &Owner{Name: "John", Country: "USA"}

fmt.Printf("%s is a %s with an %s diet\n", dog.Name, dog.Breed, dog.Animal.Diet)
```

Asterisk (`*Animal`) ilə irs = pointer — sahə dəyərlərini birbaşa dəyişmə imkanı. Çoxlu irs (Dog həm Animal həm Owner) qanunidir — Ruby-də fərqli olaraq istənilən sayda.

### 5. Anonymous struct

Adsız, inline elan + ilkinləşdirmə — "type alias çirklənməsindən" qurtuluş:

```go
Animal := struct {
    Kind string
    Diet string
}{"Dog", "Omnivorous"}    // elan + doldurma bir sətirdə

fmt.Println(Animal.Kind, "-", Animal.Diet)
```

**İstifadə yerləri:** ad lazım deyil; tez çağırış; `interface{}`-ə ucuz alternativ. **Qadağa:** dərin nested anonymous struct-lar → oxunmazlıq → baxım xətası.

### 6. Anonymous struct sahələri

Sahə adı YOX, yalnız tip — ** sahə adı = tipin adı:**

```go
animals := struct {
    int
    string
}{1, "Dog"}

fmt.Println(animals.int)     // 1
fmt.Println(animals.string)  // Dog
```

**Məhdudiyyət:** hər tip YALNIZ BİR dəfə — fərqləndirmə tip adı ilə olduğundan eyni tip iki dəfə mümkünsüz.

---

## Əsas terminlər

| Termin | İzah |
|--------|------|
| İnstans dəyişəni (@var) | Ruby: sinif daxili tək-istinadlı dəyişən; Go ekvivalenti struct |
| Public/private | Böyük hərf = export; kiçik = paket-daxili |
| Receiver | `(p *basket)` — funksiyaya bağlı struct = metod |
| Pointer receiver | attr_accessor ekvivalenti — setter+getter, birbaşa mutasiya |
| Value receiver | Kopya ilə işləyir — thread-safe, orijinal toxunulmaz |
| Embedding (irs) | `*Animal` daxilində — sahə promote; çoxlu irs OK |
| Anonymous struct | `struct{...}{values}` — inline, adsız |
| Anonymous field | Yalnız tip; ad = tip adı; hər tip bir dəfə |
| `%+v` | Struct sahə adları ilə çap |

---

## Praktik nəticə

1. **Ruby class → Go struct + receiver funksiyaları:** İnstans dəyişəni = struct sahəsi; metod = `(p *basket) method(...)`.
2. **attr_accessor lazımdırsa pointer receiver:** `(p *basket) add_item` — `*p = append(*p, ...)` orijinalı böyüdür.
3. **Value receiver seç:** kopya kifayətdirsə / konkurentlik vacibdirsə (read-only əməliyyatlar).
4. **İrs = embed:** `*Animal` birja göstərici kimi — daxili structları ayrıca initialize et (`dog.Animal = &Animal{...}`).
5. **Anonymous struct:** tez bir dəfəlik data — amma nested-dən qaçın; `interface{}`-dən ucuz və tipli.
6. **Export qaydasını yadda saxla:** D ölçüsü hər yerdə eynidir — struct adı, sahə, funksiya, konstant.

---

## Mənbə

- Kitab: *From Ruby to Golang* — Joel Bryan Juliano, Leanpub, 2021
- Chapter 2: Setter/Getters, Attribute Accessors and Structs, book səh. 3-28
- PDF səhifələri: 10-34
- İstinad: Go wiki CodeReviewComments (receiver-type semantikası)
