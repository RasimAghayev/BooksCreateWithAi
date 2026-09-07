# Chapter 7 — Containers (Konteynerlər)

## Bu fəsil nədən bəhs edir?

Generic container tipləri — genericsdən əvvəl yalnız map/slice var idi;
indi user-defined containerlər mümkündür. Əsas nümunə: Set tipi
(`type Set[E comparable] map[E]struct{}`) — map-in unikallıq xüsusiyyəti
PULSUZ gəlir; Add/Contains/All/String/Union/Intersection metodlarının
tam inkişafı; variadic NewSet/Add idiomu (type inference sayəsində
instantiate lüzumsuzlaşır); real həyat nümunəsi (job skills Intersection);
Stack məşqi (LIFO, `Stack[E any]` + data []E, Pop-un (v, ok) patterni,
pointer receiver lazımlılığı). Nəticə mülahizəsi: bir dəfə yazılan
mürəkkəb data struktur (concurrency safety, performans) bütün tiplər
üçün işləyir.

## Əsas fikirlər

### 1. Set — Map Üzərində Qurulmuş Tip
```go
type Set[E comparable] map[E]struct{}

func NewSet[E comparable](vals ...E) Set[E] {
    s := Set[E]{}
    for _, v := range vals {
        s[v] = struct{}{}
    }
    return s
}
```
- **Nədir:** unikal elementlər kolleksiyası — sıra YOX, dublikat YOX
- **Dizayn qərarları:**
  - Açar = E (comparable — unikallıq üçün == lazımdır)
  - Value = struct{} (bool YOX) — sıfır yaddaş tutumu (zero-size type)
  - Köhnə hack: `map[string]bool` + `if validCategory[category]` — işləyir
    amma yalnız string, "hacky"dir
- **Variadic + inference üstünlüyü:** `NewSet(1, 2, 3)` — `[int]` yazmağa
  ehtiyac YOX (E dəyərlərdən çıxarılır)

### 2. Add və Contains
```go
func (s Set[E]) Add(vals ...E) {
    for _, v := range vals {
        s[v] = struct{}{}
    }
}

func (s Set[E]) Contains(v E) bool {
    _, ok := s[v]
    return ok
}

s := NewSet(1, 2, 3)
s.Add(true, true)              // type error — E artıq int
fmt.Println(s.Contains(1))     // true
```
- **Add:** map-ə mənimsətmə — `struct{}{}` literal (boş struct-un yeganə
  mümkün dəyəri; "çoxlu mötərizə" görünüşü normaldır)
- **Contains:** comma-ok idiomu — `_, ok := s[v]`; açar yoxdursa zero value
  (false) qayıdır

### 3. All və String
```go
func (s Set[E]) All() []E {
    result := make([]E, 0, len(s))
    for v := range s {
        result = append(result, v)
    }
    return result
}

func (s Set[E]) String() string {
    return fmt.Sprintf("%v", s.All())
}

s := NewSet(1, 2, 3)
fmt.Println(s)  // [2 3 1] — sıra FƏRLİ ola bilər və bu NORMALDIR
```
- **All:** set üzvlərini slice kimi almaq — map-as-set hack-ində ən çətin
  hissə idi, indi metod kimi təmiz həlldir
- **String:** fmt.Stringer implementasiyası — Println birbaşa gözəl çap edir
- **MÜHÜM:** set-in tərifinə görə üzvlər sıralı DEYİL (map kimi) —
  çapda fərqli sıra görünməsi xəta deyil

### 4. Union və Intersection
```go
func (s Set[E]) Union(s2 Set[E]) Set[E] {
    result := NewSet(s.All()...)
    result.Add(s2.All()...)
    return result
}

func (s Set[E]) Intersection(s2 Set[E]) Set[E] {
    result := NewSet[E]()
    for _, v := range s.All() {
        if s2.Contains(v) {
            result.Add(v)
        }
    }
    return result
}

s1 := NewSet(1, 2, 3)
s2 := NewSet(3, 4, 5)
fmt.Println(s1.Union(s2))         // [1 2 3 4 5] — dublikat avtomatik silindi
fmt.Println(s1.Intersection(s2))  // [3]
```
- **Union:** variadic NewSet/Add sayəsində döngüsüz, bir sətirlik
- **Intersection:** s-nin hər üzvü s2-də varmı? → nəticəyə əlavə et
- **PULSUZ dublikat-təmizləmə:** heç bir xüsusi kod yazmadıq — map-in
  unikal açar xassəsi bunu verir

### 5. Real Həyat Nümunəsi — İşə Qəbul
```go
jobSkills := NewSet("go", "java")
mySkills := NewSet("go", "ruby")
matches := jobSkills.Intersection(mySkills)
if len(matches) > 0 {
    fmt.Println("You're hired!")
}
// You're hired!
```
- **Nəyə lazımdır:** tələb olunan bacarıqlar ∩ namizəd bacarıqları →
  boş deyilsə uyğunluq var

### 6. Stack Məşqi — LIFO Struktur
```go
type Stack[E any] struct {
    data []E
}

func (s Stack[E]) Len() int {
    return len(s.data)
}

func (s *Stack[E]) Push(vals ...E) {
    s.data = append(s.data, vals...)
}

func (s *Stack[E]) Pop() (v E, ok bool) {
    if len(s.data) == 0 {
        return v, false
    }
    v = s.data[len(s.data)-1]
    s.data = s.data[:len(s.data)-1]
    return v, true
}

s := stack.Stack[string]{}
s.Push("a", "b")
got, ok := s.Pop()   // "b", true — LIFO: son giren çıxır
```
- **Nədir:** Last-In-First-Out — Push əlavə edir, Pop sonuncunu çıxarır
- **Dizayn qərarları:**
  - Altında slice — sıra VACİBDIR (set-dən fərqli!) → slice seçilir
  - `Stack[E any]` — E üçün heç bir operator yoxdur → any kifayətdir
  - Push/Pop POINTER receiver — slice dəyişir (append/re-slice), value
    receiver dəyişikliyi itirərdi; Len value receiver (yalnız oxuyur)
  - Pop (v, ok) pattern — boş stack halı; zero value `var v E` avtomatik
- **Boş stack:** ok=false, heç bir panic YOX — idiomatik Go üslubu

## Əsas terminlər
- Set (çoxluq) — unikal, sırasız elementlər toplusu
- Union (birləşmə) — iki setin bütün elementləri
- Intersection (kəsişmə) — iki setin ortaq elementləri
- Stack (yığın) — LIFO struktur: Push/Pop
- LIFO (Last-In-First-Out) — son daxil olan ilk çıxır
- Comma-ok idiomu — `v, ok := m[key]` ikili qaytarma
- Zero-size type — struct{}: yaddaş tutumu sıfır
- Abstract data type (mücərrəd data tipi) — əməliyyatlarla təyin olunan tip

## Praktik nəticə

1. **Set lazımdırsa:** `map[E]struct{}` + comparable — bool value israfdır;
   metodlar toplusu (Add/Contains/All/Union/Intersection) bir dəfə yazılır.
2. **Alt-verilənlər strukturunun seçimi:** sıra lazımdırsa → slice;
   unikallıq lazımdırsa → map; hər ikisi Stack-də data []E kimi birləşir.
3. **Receiver seçimi:** dəyişən metod → pointer receiver (Push/Pop);
   oxuyan metod → value receiver (Len).
4. **"Boş" halın idarəsi:** panic YOX — (v, ok) pattern; istifadəçi tərəfdə
   müvafiq yoxlama.
5. **Generics-in əsl qazancı:** mürəkkəb data strukturları (concurrency
   safe, optimallaşdırılmış) BİR DƏFƏ yazılır — növbəti fəsil: concurrency
   safety.

## Mənbə
Pages: 125-138 (PDF 126-139)
