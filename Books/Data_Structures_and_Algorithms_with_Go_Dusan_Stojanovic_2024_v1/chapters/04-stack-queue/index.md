# Chapter 4 — Stack and Queue (səh. 203-247)

## Bu fəsil nədən bəhs edir?

İki xətti struktur element giriş qaydaları ilə: **stack** (LIFO — push/pop,
Go-da defer-in daxili stack-i, slice ilə implementasiya) və **queue** (FIFO —
enqueue/dequeue, container/list əsaslı implementasiya) + **priority queue**
(container/heap ilə min-priority nümunəsi).

## Əsas fikirlər

### 1. Stack nədir?
- **Xətti struktur + xüsusi giriş qaydası:** element yalnız **top**-dan
  əlavə/çıxarılır; top-u **Stack Pointer (SP)** göstərir
- **LIFO** (Last In First Out): son daxil olan — ilk çıxan
- Top-un əksi — **bottom**; stack bottom-a nisbətən aşağı/yuxarı böyüyə bilər
  (yeni element top-a yaxın əvvəlki/növbəti yaddaş yerinə)
- **Dinamik** struktur; tip məhdudiyyəti yoxdur, adətən homogen saxlanılır

### 2. Stack əməliyyatları
**Push — element əlavə:**
```go
func Push(v int) {
    sp = sp + 1
    stack[sp] = v
}
```
- SP bir irəliləyir, dəyəri yazılır (nümunə: 25 top-a düşür, SP ona keçir)

**Pop — element çıxar:**
```go
func Pop() int {
    if sp == -1 {        // boş stack — qeyri-normal əməliyyat
        return -1
    }
    v := stack[sp]
    sp = sp - 1
    return v
}
```
- Reallıqda yalnız SP aşağı enir; element fiziki yaddaşda qalır, istinad
  olunmur → yeni push üstünə yazacaq
- **Boş stack-dən pop = qeyri-normal əməliyyat** → əvvəlcədən yoxlanılır,
  qeyri-normal dəyər (-1) qaytarıla bilər

### 3. Stack Go-da — defer
- Go-da developer üçün stack implementasiyası YOXDUR, amma **defer** üçün
  daxili stack var
- **defer:** funksiyanın icrasını əhatə edən funksiya qayıdana qədər təxirə
  salır; arqumentlər dərhal hesablanır
```go
func main() {
    defer fmt.Print("a")
    fmt.Print("b")       // nəticə: "ba"
}
```
- Deferred çağırışlar stack-ə push olunur; funksiya qayıdanda **LIFO**
  sırası ilə pop olunub icra edilir:
```go
func main() {
    defer fmt.Print(4)
    defer fmt.Print(3)
    defer fmt.Print(2)
    defer fmt.Print(1)
}                       // nəticə: 1234
```

### 4. Stack implementasiyası (slice əsaslı)
- Slice sonu = top; ən yüksək index = stack pointer
- SP = len-1; **SP == -1 ⇒ boş stack**
```go
type Stack struct {
    stack        []int
    stackPointer int
}

func (s *Stack) Push(v int) {
    s.stack = append(s.stack, v)
    s.stackPointer = len(s.stack) - 1
}

func (s *Stack) Pop() int {
    if s.stackPointer == -1 {
        return -1
    }
    element := s.stack[s.stackPointer]
    s.stack = s.stack[:s.stackPointer]   // top elementi kəs
    s.stackPointer--
    return element
}
```
- SP-len() ilə hesablana bilərdi, amma kitab "məktəb nümunəsi"nə sadiqdir

### 5. Queue nədir?
- Xətti struktur, **2 giriş ucu:** front (başlanğıc) + rear (son)
- Elementlər front-a daxil olur, rear-dan çıxır
- **FIFO** (First In First Out): ilk daxil olan — ilk çıxan
- Dinamik, (adətən) homogen
- Xüsusi növ: **double-ended queue (deque)** — hər iki ucdan insert/remove

### 6. Queue əməliyyatları
**Enqueue — front-a daxil:**
```go
func Enqueue(v int) {
    front = front - 1
    queue[front] = v
}
```

**Dequeue — rear-dan çıxar:**
```go
func Dequeue() int {
    if front == rear {   // boş queue
        return -1
    }
    v := queue[rear]
    rear = rear - 1
    return v
}
```
- Boş queue-dan dequeue = yanlış əməliyyat → yoxlanılır (front == rear ⇒ boş)

### 7. Queue implementasiyası
Go-da hazır queue yoxdur — 2 yol:

**a) Slice ilə:** enqueue = append (sona), dequeue = index 0 kəsilir.
QARIŞIQ olur: slice başı = queue rear, sonu = queue front.

**b) container/list ilə (daha təmiz):**
```go
type Queue struct {
    queue *list.List
}

func New() *Queue {
    return &Queue{queue: list.New()}
}

func (q *Queue) Enqueue(v int) {
    q.queue.PushFront(v)       // front-a daxil
}

func (q *Queue) Dequeue() int {
    if q.queue.Len() == 0 {
        return -1              // boş — qeyri-normal dəyər
    }
    element := q.queue.Back()  // rear-dan çıxar
    q.queue.Remove(element)
    return element.Value.(int)
}

func (q *Queue) Print() {
    fmt.Print("[")
    for e := q.queue.Front(); e != nil; e = e.Next() {
        if e.Next() != nil {
            fmt.Printf("%v, ", e.Value)
        } else {
            fmt.Print(e.Value)
        }
    }
    fmt.Println("]")
}
```
- Double-linked list başlanğıc/son pointer-ləri artıq saxlayır → bir neçə
  düzüşmə kifayət edir

### 8. Priority queue
- Dequeue sırası yalnız gəliş sırasından ASILI DEYİL — **prioritet** üzrə
- **Prioritet:** hər elementə təyin olunan dəyər; elementin öz dəyəri də ola
  bilər (xüsusi dəyər məcburi deyil)
- Növlər:
  - **Min priority queue:** ən kiçik prioritet dəyəri ilk çıxır
  - **Max priority queue:** ən böyük prioritet dəyəri ilk çıxır
- Eyni prioritetdə olanlar üçün gəliş sırası qorunur (ilk gələn — ilk çıxan)

### 9. Priority queue Go-da — container/heap
- Hazır priority queue yoxdur; `container.heap` interface-i var (heap — PQ-in
  ən çox yayılmış implementasiya yolu, haqqında Ch6):
```go
type Interface interface {
    sort.Interface
    Push(x any)
    Pop() any
}
```

**Min priority queue (int dəyər = prioritet):**
```go
type Element struct{ value int }
type PriorityQueue []Element

func (pq PriorityQueue) Len() int { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool {
    return pq[i].value < pq[j].value   // min-PQ; max üçün > (adla zidd olsa da arzuolunan davranış)
}
func (pq PriorityQueue) Swap(i, j int) {
    if pq.Len() == 0 {                 // index out-of-bound qarşısı
        return
    }
    pq[i], pq[j] = pq[j], pq[i]
}
func (pq *PriorityQueue) Push(v any) {
    element := Element{value: v.(int)}  // any → int conversion
    *pq = append(*pq, element)
}
func (pq *PriorityQueue) Pop() any {
    if pq.Len() == 0 {
        return -1
    }
    queue := *pq
    n := pq.Len() - 1
    element := queue[n]
    *pq = queue[0:n]
    return element
}
```

**Vacib nüans:** metodlar özü prioritet yoxlamır — **heap.Push() / heap.Pop()**
istifadə edilməlidir (heap interface-ini implement edən tiplərdə işləyir,
arxa planda bu metodları çağırır):
```go
func main() {
    pq := make(PriorityQueue, 0)
    heap.Push(&pq, 27)
    heap.Push(&pq, 5)
    heap.Push(&pq, 1)
    heap.Push(&pq, 18)
    fmt.Println(heap.Pop(&pq))   // 1 — min prioritet ilk çıxır
    // ...
}
```

## Termindirmə (AZ)
- Stack — Stek (yalnız top-dan add/remove)
- Stack Pointer (SP) — Stek göstəricisi
- Top / Bottom — Zirvə / Dip
- LIFO — Son Gələn — İlk Çıxan
- Push / Pop — Yığ / Çıxar
- Queue — Növbə
- Front / Rear — Baş / Quyruq (növbə ucları)
- FIFO — İlk Gələn — İlk Çıxan
- Enqueue / Dequeue — Növbəyə sal / Növbədən çıxar
- Double-ended queue (deque) — İkiuclu növbə
- Priority queue — Prioritet növbəsi
- Min / Max priority queue — Min / Maks prioritet növbəsi
- Priority — Prioritet (çıxarma sırasını təyin edən dəyər)
- defer — Təxirə salma (funksiya qayıdana qədər)
- Deferred stack — Gecikdirilmiş çağırışların stack-i (LIFO)
- container/heap — Heap əməliyyatları paketi (PQ üçün)
- heap.Interface — sort.Interface + Push/Pop

## Kviz sualları
1. Stack-də SP = -1 nə deməkdir? (stack boşdur — pop qeyri-normal əməliyyatdır)
2. defer-dən 4 dəfə istifadə edilsə icra sırası nə olur? (LIFO — son
   deferlənən birinci icra olunur)
3. Pop-da element fiziki yaddaşdan silinirmi? (xeyr — SP aşağı enir;
   element istinadsız qalır, növbəti push üstünə yazır)
4. Queue-da front və rear hansı əməliyyatlarda iştirak edir? (enqueue
   front-a, dequeue rear-dan)
5. Front == rear nəyi göstərir? (queue boşdur)
6. Max priority queue üçün Less() necə yazılır? (pq[i] > pq[j] — metod adı
   ilə zidd olsa da, istənilən davranış verir)
7. Niyə PQ metodları özbaşına çağırılmır? (heap.Push/heap.Pop heap
   quruluşunu saxlayaraq prioriteti təmin edir)
