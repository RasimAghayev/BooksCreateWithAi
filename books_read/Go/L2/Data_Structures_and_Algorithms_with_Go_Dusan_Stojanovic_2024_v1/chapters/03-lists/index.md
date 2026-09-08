# Chapter 3 — Lists (səh. 160-202)

## Bu fəsil nədən bəhs edir?

List anlayışı: node-lardan ibarət xətti kolleksiya; list növləri (single-linked,
double-linked, circular), əməliyyatlar (insert/remove/find/concatenate),
single-linked list-in Go-da sıfırdan implementasiyası, list vs array müqayisəsi
və standart kitabxananın `container/list` (double) + `container/ring`
(circular) paketləri.

## Əsas fikirlər

### 1. List nədir?
- **Linear kolleksiya**, elementlər **pointer**-lə birləşir
- Element = **node** = data + pointer(lər); pointer sayı list tipindən asılıdır
- List adətən ilk node-u göstərən pointer ilə idarə olunur (bu pointer
  listin ÖZÜ deyil)
- İlk node = **head**, son node = **tail**
- **Ordered list:** node sırası vacibdirsə (məs. artan sıra); əks halda
  **unordered**
- Yaddaşda **linked representation**: qonşu node-lar ardıcıl yaddaş
  sahələrində deyil, yaddaş boyu səpələnə bilər

### 2. List növləri
| Tip | Qovluq | Xassə |
|---|---|---|
| Single-linked (tək-zəncirli) | next | yalnız irəli gediş |
| Double-linked (iki-zəncirli) | next + prev | hər iki istiqamət |
| Circular (dairəvi) | son node → ilk node | "ring" də adlanır |

### 3. Əməliyyatlar — hamısı pointer manipulyasiyasıdır
**Insert — unordered list (əvvələ):**
1. Yeni node yaradılır
2. Yeni node-un next-i list pointer-inin göstərdiyi node-u göstərir
3. List pointer yeni node-u göstərir

**Insert — ordered list:**
- Əlavə **temp** pointer: yeni node-un hansı node-dan SONRA daxil
  olacağını tapana qədər iterasiya
- Yeni node-un next-i temp-in next-unu göstərir; temp-in next-i yeni node-u

**Remove:**
- temp silinən node-dan ƏVVƏLKİ node-u göstərir
- temp-in next-i silinən node-dan sonrakına yönləndirilir → node artıq
  istinad olunmur → **garbage collector** təmizləyir
- GC: John McCarthy tərəfindən ~1959-da yaradılıb; proqramla paralel
  işləyib istifadəsiz obyektləri (orphan-ları) axtarır və yaddaşı azad edir
- Memory leak-in qarşısını almaq üçün: silinən node-un next-ini nil etmək
  düzgün detach'dir

**Search:** pointer-lərlə iterasiya; ordered list-də uğursuz axtarışı daha
tez aşkar etmək olar

**Concatenate:** birinci listin son node-u ikinci listin ilk node-unu
göstərir; ikincinin head-i nil olur

### 4. Single-linked list — Go implementasiyası
```go
type Node struct {
    next  *Node
    value int
}

type List struct {
    head *Node
    len  int
}

func New() List {
    return List{head: nil, len: 0}
}
```

**Insert (unordered, əvvələ) — O(1):**
```go
func (l *List) Insert(v int) {
    node := Node{next: nil, value: v}
    if l.head != nil {
        node.next = l.head
    }
    l.head = &node
    l.len++
}
```

**InsertOrdered (sıralı daxil etmə):**
```go
func (l *List) InsertOrdered(v int) {
    node := Node{next: nil, value: v}
    if l.head == nil {          // edge case: boş list
        l.head = &node
        l.len++
        return
    }
    temp := l.head
    for temp.next != nil && temp.next.value < v {
        temp = temp.next
    }
    node.next = temp.next
    temp.next = &node
    l.len++
}
```
- Yaxşı praktika: edge case-ləri əvvəlcədən işlə

**Remove — edge case-lər: boş list; tək node-lu list:**
```go
func (l *List) Remove(v int) {
    if l.head == nil {
        return
    }
    if l.len == 1 {          // tək node
        l.head = nil
        l.len--
        return
    }
    for temp := l.head; temp != nil; temp = temp.next {
        if temp.next.value == v {
            node := temp.next
            temp.next = node.next
            node.next = nil    // düzgün detach
            l.len--
            return
        }
    }
}
```

**Find:**
```go
func (l *List) Find(v int) *Node {
    for temp := l.head; temp != nil; temp = temp.next {
        if temp.value == v {
            return temp
        }
    }
    return nil
}
```

**Concatenate:**
```go
func (l *List) Concatenate(l2 List) {
    temp := l.head
    for temp.next != nil {
        temp = temp.next
    }
    temp.next = l2.head
    l2.head = nil
}
```

**Print:**
```go
func (l *List) Print() {
    fmt.Print("[")
    for temp := l.head; temp != nil; temp = temp.next {
        if temp.next != nil {
            fmt.Printf("%v, ", temp.value)
        } else {
            fmt.Print(temp.value)
        }
    }
    fmt.Println("]")
}
```

### 5. Lists vs Arrays
| Mezar | Array | List |
|---|---|---|
| Yaddaş | node başına pointer yoxdur; kapasiti tam ayrılır | node daha çox yer tutur (pointer); kapasiti tədavücə artır |
| Elementə giriş | index ilə DÜZQÜNCA sürətli | node-a çatmaq üçün iterasiya lazım |
| Insert (mövqeyə) | mövqedən başlayaraq hamısı bir sağa sürüşür; doludursa reallokasiya | yalnız pointer manipulyasiyası |
| Remove | mövqedən sonrakılar sola sürüşür | node-u listdən ayır |
| Nə vaxt? | kapasiti əvvəlcədən məlum + sürətli giriş lazımdır | random mövqelərdə çox insert/remove əməliyyatları |

### 6. Lists in Go — standart kitabxana
**Double-linked list — `container/list`:**
```go
type List struct {
    root Element   // root node listin hissəsi DEYİL
    len  int
}
type Element struct {
    next, prev *Element
    list       *List
    Value      any
}
```
```go
l := list.New()
head := l.PushFront(5)        // əvvələ — node qaytarır
tail := l.PushBack(27)        // sona
l.InsertBefore(1, tail)       // node-dan əvvəl
l.InsertAfter(18, head)       // node-dan sonra
l.Remove(tail)
l.Len()

for node := l.Front(); node != nil; node = node.Next() {
    fmt.Print(node.Value)     // irəli iterasiya
}
for node := l.Back(); node != nil; node = node.Prev() {
    fmt.Print(node.Value)     // geri iterasiya
}
```

**Circular list — `container/ring`:**
```go
type Ring struct {
    next, prev *Ring
    Value      any
}
```
- Başlanğıc və sonu YOXDUR → istənilən node bütün listə istinad ola bilər
- Zero value: tək node-lu, Value=nil list
```go
r := ring.New(4)              // 4 node-lu ring
r.Len()
r.Next(); r.Prev()            // qonşu node-lar
for i := 0; i < r.Len(); i++ {
    r.Value = i               // dəyər yaz
    r = r.Next()
}
r3 := r1.Link(r2)             // iki ring-in birləşdirilməsi
r.Do(func(p any) {            // hər node-da funksiya icrası
    fmt.Println(p.(int))
})
```
- Single-linked üçün standart kitabxanada implementasiya YOXDUR →
  özümüz yazdıq (yuxarıda)

## Termindirmə (AZ)
- List — Siyahı (pointer-lə birləşən xətti kolleksiya)
- Node — Node/düyün (list elementi: data + pointer)
- Head / Tail — Baş / Quyruq (ilk / son node)
- Single-linked list — Tək-zəncirli siyahı
- Double-linked list — İki-zəncirli siyahı
- Circular list (ring) — Dairəvi siyahı
- Ordered / Unordered list — Sıralı / Sıralanmamış siyahı
- Linked representation — Zəncirvari yaddaş təmsili
- Garbage collector — Zibil yığıcısı (John McCarthy, ~1959)
- Orphan — Yetim obyekt (istinad olunmayan yaddaş)
- Memory leak — Yaddaş sızması
- Detach — Ayırmaq (node-u listdən düzgün ayırmaq)
- Concatenate — Birləşdirmək (zəncirləmək)
- container/list — ikizəncirli list paketi
- container/ring — dairəvi list paketi

## Kviz sualları
1. Node-un neçə pointer-i ola bilər? (tipdən asılı: single=1, double=2,
   circular=sonuncu ilkinə bağlanır)
2. List-də node silinəndə yaddaş kim təmizləyir? (GC — node istinadsız
   qalır; next-i nil etmək leak qarşısını alır)
3. Ordered list-də insert nə üçün temp pointer istifadə edir? (daxil
   olunacağı yer tapılana qədər iterasiya, sonra pointer dəyişimi)
4. Hansı əməliyyətdə list array-dən üstündür? (random mövqedə insert/remove
   — sürüşdürmə yoxdur, yalnız pointer dəyişir)
5. `container/list`-in root element-i listin hissəsidirmi? (xeyr —
   sentinel/root node, listə daxil deyil)
6. Ring-in zero value nədir? (Value=nil tək node-lu dairəvi list)
