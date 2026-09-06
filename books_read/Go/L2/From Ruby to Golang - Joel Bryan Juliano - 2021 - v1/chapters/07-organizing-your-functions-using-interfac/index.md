# Chapter 7 — Funksiyaların İnterfeys ilə Təşkili (book səh. 101-111)

## Bu chapter nədən bəhs edir?

Ruby module-mixin yanaşmasının Go interfeyslərinə tərcüməsi: interfeys API kontraktı kimi (self-documenting reference), tip kontraktı kimi (heterojen kolleksiya) və qayıdış dəyərlərinin şərtlənməsi kimi (return type contract) — 3 istifadə forması tam kodlarla.

---

## Əsas fikirlər

### 1. Niyə interfeys?

Funksiyaları xüsusi implementasiyalara təşkil etmək + komponent idarəsini təmizləşdirmək: **interfeys = implementasiyalar ARASINDA API kontraktı.**

**Ruby müqayisəsi:** Go-da interfeysin DİREKT analoqu YOXDUR; ən yaxın — modullarla metodların təşkili:

```ruby
require 'set'
module Basket
  attr_accessor :items
  def setup_basket; @items = Set.new; end
  def add_item(item); @items << item; end
  def remove_item(item); @items.delete(item); end
end

class ProduceBasket < String        # superclass = TİP KONTRAKTI
  include Basket
  def initialize; super; setup_basket; end
end
```

### 2. İnterfeys = Self-Documenting API Reference

Funksiyaları ƏVVƏLCƏ interfeysdə deklarativ təyin et — sonra implementasiya:

```go
type ProduceBasket interface {
    AddItem(entry Produce)
    RemoveItem(entry Produce)
    Items()
}

type Basket []Produce
type Produce string

func (p *Basket) AddItem(entry Produce) {
    *p = append(*p, entry)
    fmt.Printf("%s Added\n", entry)
}

func (p *Basket) RemoveItem(entry Produce) {
    s := *p
    for i, v := range s {
        if v == entry {
            s = append(s[:i], s[i+1:]...)    // slice-dan element silmə idiomu
            break
        }
    }
    *p = s
    fmt.Printf("%s Removed\n", entry)
}

func (p Basket) Items() {
    for _, v := range p { fmt.Println(v) }
}
```

**İstifadə — interfeys dəyişəni kimi yüklə:**

```go
var basket ProduceBasket = &Basket{}   // interfeysə pointer-assign

basket.AddItem("Apple")      // Apple Added
basket.AddItem("Mango")      // Mango Added
basket.RemoveItem("Mango")   // Mango Removed
basket.Items()               // Apple ...
```

Interfeys təyini = funksiya siyahısının SƏNƏDİ — hansı davranış gözlənilir bir baxışda görünür.

### 3. İnterfeys = Type Contract (heterojen kolleksiya)

`ProduceBasket` tətbiq edən İSTƏNILƏN tip eyni kolleksiyada yaşaya bilər:

```go
fruits := new(Basket)
fruits.AddItem("Apple")
fruits.AddItem("Mango")

veggies := new(Basket)
veggies.AddItem("Broccolli")

fruits.RemoveItem("Mango")

var items []ProduceBasket             // İNTERFEYS tipində massiv!
items = append(items, fruits)         // hər iki basket-tip uyuşur
items = append(items, veggies)

for _, v := range items {
    v.Items()        // Apple / Broccolli — polymorphic çağırış
}
```

Bu, Ruby-dəki `ProduceBasket < String` (superclass kontraktı) + include kombinasiyasının Go tərcüməsidir: **tip koerası interfeyslə deklarativ təmin olunur.**

### 4. İnterfeys = Return Value Contract

Qayıdış dəyərinin kontraktı — müxtəlif strukturlar eyni interfeysi qaytara bilər:

```go
type Produce interface {
    Flavour() string
    Kind() string
}

type Item struct {
    Name string
}

func (i Item) Flavour() string {
    flavour := map[string]string{
        "Apple": `It's a little sour and bitter, mostly sweet...`,
        "Kale":  `It boasts deep, earthy flavors...`,
    }
    return fmt.Sprintf("Item %s flavour is: %s\n", i.Name, flavour[i.Name])
}

func (i Item) Kind() string {
    kind := map[string]string{"Apple": "Fruit", "Kale": "Veggies"}
    return fmt.Sprintf("Item %s is a %s\n", i.Name, kind[i.Name])
}

// İnterfeys massivi — struct-lar interfeysi ŞƏRTLƏYİR:
var basket []Produce
basket = append(basket, Item{Name: "Apple"})
basket = append(basket, Item{Name: "Kale"})

for _, v := range basket {
    fmt.Println(v.Flavour())
    fmt.Println(v.Kind())
}
```

**Müqayisə cədvəli (Ruby mixin vs Go interface):**

| Aspect | Ruby | Go |
|--------|------|-----|
| Təşkilat prinsipi | module + include | interface + metod tətbiqi |
| Kontrakt tipi | superclass (`< String`) | interfeys tipi |
| Çoxlu davranış | bir neçə module include | bir neçə interfeys şərtləmə |
| Polymorf kolleksiya | eyni superclass obyektləri | eyni interfeys tipləri |

---

## Əsas terminlər

| Termin | İzah |
|--------|------|
| API kontraktı | İnterfeys = implementasiyalar arası müqavilə |
| Self-documenting | İnterfeys funksiya siyahısının sənədi |
| `var x Iface = &Impl{}` | İnterfeys dəyişəninə implementasiya yükləmə |
| Type contract | Heterojen tiplərin interfeys massivində birləşməsi |
| Return value contract | Funksiya interfeys qaytarır — struktur dəyişəndə davranış sabit |
| Implicit satisfaction | `implements` YOX — metodlar kifayət |
| `append(s[:i], s[i+1:]...)` | Slice-dan element silmə idiomu |

---

## Praktik nəticə

1. **Metodları əvvəlcə interfeysdə təyin et** — API-ni sənədləş + gələcək implementasiyalara yol göstər.
2. **Heterojen kolleksiyalar üçün `[]Iface`:** Fərqli strukturlar (fruits/veggies baskets) eyni loop-da işləyir.
3. **Qayıdış tiplərini interfeys et:** Çağıran struktur bilir, konkret tipdən asılı olmur — yeni Item tipləri sorğusuz qoşulur.
4. **Ruby module → Go interface xəritəsi:** include = implicit tətbiq; attr_accessor = pointer receiver; superclass = interfeys tipi.
5. **Slice element silmə:** `s = append(s[:i], s[i+1:]...)` — kitabın RemoveItem-dan praktik idiom.

---

## Mənbə

- Kitab: *From Ruby to Golang* — Joel Bryan Juliano, Leanpub, 2021
- Chapter 7: Organizing your Functions using Interface, book səh. 101-111
- PDF səhifələri: 108-117
