# Chapter 4 — Go Generics (Go Generics)

## Bu chapter nədən bəhs edir?

Generics-in sintaksisi (Go 1.18+), constraints (any, comparable, xüsusi interfeys
constraints, ~ supertype operatoru), generic data tipləri və generic struct-lar (linked
list), cmp/slices/maps standart paketləri, shallow vs deep copy və generics-in nə vaxt
istifadə edilməsi qərarı.

## Əsas fikirlər

### 1. Generics Nədir və Nəyə Lazımdır
**Nədir:** Funksiyanın parametr tipini sonradan (compile vaxtı, çağırış kodundan)
müəyyən etmək imkanı. `[T any]` — T hər hansı tip ola bilər; tip runtime-da yox,
COMPILE vaxtında təyin+enforced olunur.

**Üstünlükləri:**
- Eyni funksiya birdən çox tip ilə işləyir — kopyala-yapışdırıq funksiyalar yox
- İnterfeysdən fərqi: tip aşkarlama kodu (type switch/ assertion) TƏLƏB ETMİR

**Çatışmamazlıqları:**
- Statik tipə görə kod daha SÜRTƏTLİ işləyir — generics flexibility qarşılığında sürət
  itirir; compile vaxtı da uzanır
- Məcburi DEYİL — generics-siz də mükəmməl Go yazmaq olur (cəmiyyətin böyük hissəsi
  belə edir)

**Kitabdan kod nümunəsi:**
```go
func PrintSlice[T any](s []T) {
    for _, v := range s {
        fmt.Print(v, " ")
    }
    fmt.Println()
}

PrintSlice([]int{1, 2, 3})
PrintSlice([]string{"a", "b", "c"})
PrintSlice([]float64{1.2, -2.33, 4.55})
// Bir funksiya, 3 tip — compiler hər tip üçün ayrıca versiya yaradır
```
Çoxlu tip parametri: `[T, U, W any]`.

### 2. Constraints — Tip Məhdudiyyətləri
**Nədir:** Generic funksiyanın qəbul etdiyi tiplərin siyahısı — məntiqi səhvlərin
compile vaxtında tutulması.

**comparable (built-in):** == və != ilə müqayisə oluna bilən tiplər.

**Kitabdan kod nümunəsi:**
```go
func Same[T comparable](a, b T) bool {
    return a == b
}
Same(4, 3)         // false
Same("aa", "aa")   // true
Same([]int{1}, []int{1})   // COMPILE XƏTASI: []int does not satisfy comparable
```
- Slice-lar comparable DEYİL (array-lar müqayisə olunur!)

**Xüsusi constraint (interfeys + union):**
```go
type Numeric interface {
    int | int8 | int16 | int32 | int64 | float64   // union — | operatoru
}

func Add[T Numeric](a, b T) T {
    return a + b
}
Add(4, 3)      // 7
Add(4.1, 3.2)  // 7.3
Add(4.1, 3)    // COMPILE XƏTASI: default type int of 3 does not match
                // inferred type float64 — hər iki parametr EYNİ tip olmalı
```

### 3. Supertype (~) — Underlying Tipə İcazə
**Problem:** `type AnotherInt int` → AnotherInt Numeric-də YOXdur (nominal typing).
**Həll:** `~int` — underlying tipi int olan BÜTÜN tiplər (aliaslar daxil).

**Kitabdan kod nümunəsi:**
```go
type AnotherInt int
type AllInts interface {
    ~int              // underlying int olan hər tip
}

func AddElements[T AllInts](s []T) T {
    sum := T(0)       // T-nin zero value constructor forması
    for _, v := range s {
        sum = sum + v
    }
    return sum
}

s := []AnotherInt{0, 1, 2}
AddElements(s)   // 3 — ~int sayəsində AnotherInt qəbul edilir
```

### 4. Slice Constraint-inin Qısa Yazılışı
```go
func f1[S interface{ ~[]E }, E interface{}](x S) int { return len(x) }  // uzun
func f2[S ~[]E, E interface{}](x S) int { return len(x) }               // interface{} buraxılır
func f3[S ~[]E, E any](x S) int { return len(x) }                        // any ilə — ən sadə
```
Üçü ekvivalentdir — `~[]E` = "underlying-i slice olan tip".

### 5. Generic Data Tiplər
**Kitabdan kod nümunəsi:**
```go
type TreeLast[T any] []T

func (t TreeLast[T]) replaceLast(element T) (TreeLast[T], error) {
    if len(t) == 0 {
        return t, errors.New("This is empty!")
    }
    t[len(t)-1] = element
    return t, nil
}

tempStr := TreeLast[string]{"aa", "bb"}
tempStr.replaceLast("cc")

tempInt := TreeLast[int]{12, -3}
tempInt.replaceLast(0)
```
**Sub-kod izahı:** metod receiver-i `TreeLast[T]` — metod da generic olur; tip arqumenti
istifadə yerində (`TreeLast[string]`) verilir.

### 6. Generic Struct — Linked List
**Kitabdan kod nümunəsi:**
```go
type node[T any] struct {
    Data T
    next *node[T]      // öz-özünə istinad — eyni T
}

type list[T any] struct {
    start *node[T]
}

func (l *list[T]) add(data T) {
    n := node[T]{Data: data, next: nil}
    if l.start == nil {
        l.start = &n
        return
    }
    if l.start.next == nil {
        l.start.next = &n
        return
    }
    temp := l.start
    l.start = l.start.next
    l.add(data)
    l.start = temp
}

var myList list[int]
myList.add(12)
// Travers: cur := myList.start; for { ...; cur = cur.next }
```
**Qayda:** Bir list daxilində bütün node-lar EYNİ tipdə olmalıdır; amma string-list,
int-list, struct-list — hamısı eyni kod ilə.

### 7. cmp Paketi (Go 1.21+)
```go
cmp.Compare(5, 4)     // 1   (x>y)
cmp.Compare(4, 5)     // -1  (x<y)
cmp.Less(4, 5.1)      // true — int-i float64-ə ÖZÜ çevirib müqayisə edir!
```
Ordered dəyərlərin müqayisəsi üçün; slices/maps paketlərinin əsasını təşkil edir.

### 8. Shallow vs Deep Copy
- **Shallow copy (dayaz nüsxə):** yeni dəyişən + dəyərlərin adi assignment ilə köçürülməsi.
  Referens tiplərində (pointer, struct içində struct) daxili obyektlər PAYLAŞILIR.
- **Deep copy (dərin nüsxə):** bütün dəyərlər REKURSİV kopyalanır — struct/sahə/pointer
  zəncirinin sonuna qədər; sonsuz dövrə (circular reference) diqqət tələb edir.

### 9. slices Paketi (Go 1.21+)
**Kitabdan kod nümunəsi:**
```go
s2 := slices.Clone(s1)              // shallow copy — müstəqil slice
s1 = slices.Compact(s1)             // ardıcıl dublikatları birə endirir
                                    // (sorted slice-də ən yaxşı işləyir)
slices.Contains(s1, 2)              // element varmı?
s4 = slices.Clip(s4)                // cap → len (yaddaş qənaəti)
slices.Min(s1); slices.Max(s1)      // min/max element
s2 = slices.Replace(s2, 1, 3, 100, 200)  // s2[1:3] → 100, 200
slices.Sort(s2)                     // artan sıralama
```

### 10. maps Paketi (Go 1.21+)
**Kitabdan kod nümunəsi:**
```go
// Şərtli silmə — bütün tək dəyərləri sil:
func delete(k string, v int) bool { return v%2 != 0 }
maps.DeleteFunc(m, delete)

n := maps.Clone(m)                  // shallow clone
maps.Equal(m, n)                    // açar-dəyər bərabərliyi
maps.Copy(m, n)                     // n-dəki cütləri m-ə köçür (mövcud açarlar üstünə yazır)

// Fərqli dəyər tipli map-lərin funksiya ilə müqayisəsi:
func equal(v1 int, v2 float64) bool { return float64(v1) == v2 }
eq := maps.EqualFunc(t, mFloat, equal)   // true
```

### 11. Nə Vaxt Generics İstifadə Etməli
**Kitabın tövsiyələri:**
- Kod birdən çox data tipi ilə işləməli olduqda
- İnterfeys/reflection həlli kodu MÜRƏKKƏBLƏŞDİRDİKDƏ
- Gələcəkdə yeni tiplər dəstəklənəcəksə
- Məqsəd sadəlik və maintenance-dir — imkan nümayişi YOX
- Razı deyilsənz — istifadə ETMƏYİN; məcburiyyət yoxdur

## Əsas terminlər
- Generics — tip parametrli proqramlaşdırma (Go 1.18)
- Type Parameter (tip parametri) — `[T any]`-dəki T
- Constraint (məhdudiyyət) — icazə verilən tiplərin interfeys-union siyahısı
- comparable — == ilə müqayisə edilə bilən tiplər üçün built-in constraint
- Union (union) — `int | float64` tip birləşməsi
- Supertype (~) — underlying tipə görə qəbul (`~int`)
- any — interface{}-nin qısa adı
- Shallow Copy (dayaz nüsxə) — assignment ilə səth kopyası
- Deep Copy (dərin nüsxə) — rekursiv tam kopya
- slices/cmp/maps — generics əsaslı Go 1.21+ standart paketləri

## Praktik nəticə

(1) Generics flexibility verir, amma statik tipə nisbətən yavaşdır və compile uzanır —
qərarı xərc anlayışı ilə ver. (2) `any` hər yerdə işlətmə — constraint-lər məntiqi
səhvləri compile-a keçirir. (3) Add(4.1, 3) xətası: generic parametrlər eyni T —
qarışıq tip çağırışları ayrı constraint-lərlə mümkün. (4) ~ (supertype) olmadan
`type MyInt int` constraint-dən kənar qalır. (5) Standart paketlər (slices/maps/cmp)
generics-in praktik faydasının ən yaxşı nümunəsidir — təkəri özün icad etməzdən əvvəl
onlara bax. (6) Linked list kimi strukturlar generics ilə bir dəfə yazılıb bütün
tiplər üçün işləyir. (7) Generics məcburiyyət deyil — interfeyssiz də, genericsiz də
düzgün Go mümkündür; meyar sadəlikdir.

## Mənbə
Pages: 133-152 (PDF 164-183)
