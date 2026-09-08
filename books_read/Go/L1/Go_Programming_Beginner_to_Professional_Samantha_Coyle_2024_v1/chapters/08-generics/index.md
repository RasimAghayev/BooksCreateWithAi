# Chapter 8 — Generic Algorithm Superpowers (səh. 310-325)

## Bu fəsil nədən bəhs edir?

Go 1.18 generics: nə zaman istifadə etmək, type parametrləri `[Num int | float64]`,
type constraint interfeysləri (Number, comparable), type inference (avtomatik
tip çıxarma), generics vs interfaces seçimi və best practices.

## Əsas fikirlər

### 1. Nə zaman generics?
**Motivasiya:** təkrarlanan (boilerplate) kod — eyni məntiq, fərqli tiplər.

```go
// Generics OLMADAN — hər tip üçün kopya:
func findMaxInt(nums []int) int {
    if len(nums) == 0 {
        return -1
    }
    max := nums[0]
    for _, num := range nums {
        if num > max {
            max = num
        }
    }
    return max
}

func findMaxFloat(nums []float64) float64 { /* EYNİ kod! */ }

// GENERICS İLƏ — bir funksiya:
func findMaxGeneric[Num int | float64](nums []Num) Num {
    if len(nums) == 0 {
        return -1
    }
    max := nums[0]
    for _, num := range nums {
        if num > max {
            max = num
        }
    }
    return max
}
```

**Qərar meyarları:**
1. Normal Go kodu yaz, tip DİZAYN ETMƏ — generics alətdir, fəlsəfə deyil
2. Boilerplate görəndə tətbiq et (təkrar = siqnal)
3. Mürəkkəb strukturlarda abstraksiya üçün
4. Çoxşamlı tiplər üçün çeviklik
5. Gələcək tip genişlənmələrini ön-gör
6. **ERKƏN abstraksiya etmə — yalnız ehtiyac yarandanda!**

### 2. Type parametrləri
```go
func findMaxGeneric[Num int | float64](nums []Num) Num { ... }
//          tip parametri ↑  kvadrat mötərizə ↑    istifadə ↑
```
- `[T any]` — T hər tip ola bilər
- Tip parametrləri adətən BÖYÜK hərflə (tip olduqlarını vurğulamaq üçün)
- **Type set** = icazə verilən tiplərin birliyi (int | float64)
- **Instantiation** — konkret tip arqumentlərinin verilməsi; tiplər
  avtomatik İNFERENCE olunur:
```go
maxGenericInt := findMaxGeneric([]int{1, 32, 5, 8, 10, 11})
// string ötürsən → COMPILE XƏTASI:
// string does not satisfy int | float64
```

### 3. Type constraint interfeysləri
**Nədir:** tip parametrinin icazəli tiplərini təyin edən interfeys
("meta-tip"); constraint mövqeyində (tip parametri siyahısında) durur.

```go
// Constraint interfeysi — Daha OXUNAQLI:
type Number interface {
    int | float64
}

func findMaxGeneric[Num Number](nums []Num) Num {
    if len(nums) == 0 {
        return -1
    }
    max := nums[0]
    for _, num := range nums {
        if num > max {
            max = num
        }
    }
    return max
}

maxGenericInt := findMaxGeneric([]int{1, 32, 5, 8, 10, 11})        // 32
maxGenericFloat := findMaxGeneric([]float64{1.1, 32.1, 5.1})         // 32.1
```
- Standart kitabxanada eksperimental `constraints` paketi var

### 4. comparable constraint
**Nədir:** `==` və `!=` operatorlarını dəstəkləyən tiplər; map AÇARLARI
üçün MƏCBURİDİR.

```go
// K = comparable (map açarı), V = int|float64 (dəyər):
func FindLargestRanchStock[K comparable, V int | float64](m map[K]V) K {
    var stock V
    var name K
    for k, v := range m {
        if v > stock {
            stock = v
            name = k
        }
    }
    return name
}

animalStock := map[string]int{
    "Chicken": 5,
    "Cattle":  20,
    "Horses":  4,
}
miscStock := map[string]float64{
    "Hay":        5.5,
    "Feed":       1.2,
    "Fertilizer": 4.5,
}

fmt.Printf("The largest stocked item on the ranch is %s\n",
    FindLargestRanchStock(animalStock))   // Cattle
fmt.Printf("The largest stocked item on the ranch is %s\n",
    FindLargestRanchStock(miscStock))     // Hay
```
- comparable OLMASAYDI — `map[K]V` referansı compile keçməzdi

### 5. Type inference
**Nədir:** kompilyatorun tip arqumentlərini Funksiya arqumentlərindən
çıxarması — çox vaxt açıq tip yazmaq lazım deyil.

```go
// İnference (gizli):
largestStockOnRanchInt := FindLargestRanchStock(animalStock)
// Ekvivalent açıq forma:
largestStockOnRanchInt := FindLargestRanchStock[string, int](animalStock)
```

**Açıq tipin lazım olduğu hallar:**
- **Ambiguous types** — bir neçə tip constraint-i qane edir
  (PrintType("Hello") / PrintType(42))
- **Multiple type parameters** — PrintTwoTypes[int, string](42, "Hello")
- **Chained calls** — növbəti funksiyaya tip axını lazımdır
- Oxunaqlılıq üçün istənilən halda açıq yazıla bilər

### 6. Generics vs Interfaces
| Mezar | Seçim |
|---|---|
| Abstraksiya qatı, başqaları implement edəcək | **Interfaces** |
| Fərqli DAVRANIŞLARI tutmaq (metod dəsti) | **Interfaces** |
| Tip-agnostik funksiyalar, COMPILE-vaxtı təhlükəsizlik | **Generics** |
| Məlumat strukturu + tip parametri | **Generics** |
| Polimorfizm | Interfaces |

- Interfeyslər də texniki olaraq generic proqramlaşdırma formasıdır —
  ortaq davranışları metod kəsb edir
- Generics: saxlama daha effektiv ola bilər; type assertion YOX;
  tam compile-time yoxlama

### 7. Best practices
1. **Metod yerinə funksiya** — funksiya tipə bağlı deyil → daha çevik
2. **Metodu funksiyaya çevirmək asandır** — əksinə daha çətin
3. **Metod tələb edən constraint-lərdən QAÇIN** — `comparable` kimi
   ümumi constraint-lər gələcəkdə daha çox tipə şərait verir
4. Sadəliyi qoru — generics alətdir, hər yerdə deyil

## Activity icmalı
- **8.01:** findMinGeneric[Num int | float64] — minimum tapımı (max-ın
  tərsi; `num < min` müqayisəsi)

## Əsas terminlər
- Generics (Go 1.18+) — tip parametrli kod
- Type parameter — `[Num ...]` kvadrat mötərizədə
- Type set — icazəli tiplər birliyi (`int | float64`)
- Type argument / instantiation — konkret tipin qoyulması
- Type constraint — tip parametrinin meta-tipi (interfeys)
- Constraint position — tip parametri siyahısındakı yeri
- Number interfeysi — `int | float64` constraint
- comparable — ==/!= müqayisə oluna bilən; map açarı tələbi
- Type inference — tiplərin avtomatik çıxarılması
- Boilerplate — təkrarlanan eyni məntiq
- Premature abstraction — erkən/hədsiz abstraksiya (qaçın!)

## Praktik nəticə
Generics-i yalnız təkrarlanan kodu görəndə tətbiq et: hər tip üçün
eyni funksiya yazırsansa → `[Num Constraint]` parametrli tək funksiya.
Constraint-ləri adlandırılmış interfeyslərə çıxar (`type Number interface
{ int | float64 }`) — oxunaqlılıq. Map açarı kimi istifadə olunan tip
parametri mütləq `comparable`. Tip arqumentləri adətən inference olunur;
bulanıq/növbəli hallarda açıq yaz. Davranış abstraksiyası → interfeys;
tip-safety ilə təkrarsızlıq → generics. Funksiyalar metodlardan üstün.

## Mənbə
Pages: 310-325 (PDF 310-325)
