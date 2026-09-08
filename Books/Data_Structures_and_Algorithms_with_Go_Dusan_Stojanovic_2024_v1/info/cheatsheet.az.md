# DS&A with Go — Cheat Sheet (AZ)

## Slice əsasları

```go
a := [8]int{1, 18, 5, 27, 8, 25, 9, 21}
s := a[1:4]        // [18 5 27]; low=0, high=len defolt
len(s); cap(s)     // 3 / 7
s = append(s, 3, 21)  // sonda; cap çatmazsa yeni array
make([]int, 0, 5)     // type, len, cap
matrix := [3][3]int{{1,2,3},{4,5,6},{7,8,9}}  // row-major
```

## Axtarış

```go
// Sequential — O(n); sıralmamış üçün yeganə
func SeqSearch(array []int, key int) int {
    for i, elem := range array {
        if key == elem { return i }
    }
    return -1
}

// Binary — O(log n); YALNIZ sıralı
func BinSearch(array []int, key int) int {
    low, high := 0, len(array)-1
    for low <= high {
        mid := (low + high) / 2
        if key == array[mid] { return mid }
        if key < array[mid] { high = mid - 1 } else { low = mid + 1 }
    }
    return -1
}

// sort paketi: SearchInts/SearchFloat64s/SearchStrings + Search(n, f)
// match yoxdursa: daxil edilməli index
```

## Sıralama

```go
// Bubble O(n²) — qonşu müqayisə + dəyişmə
for i := 0; i < n-1; i++ {
    for j := 0; j < n-i-1; j++ {
        if a[j] > a[j+1] { a[j], a[j+1] = a[j+1], a[j] }
    }
}

// Selection O(n²) — min tap, yerləri dəyiş
for i := 0; i < len(a)-1; i++ {
    min, pos := a[i], i
    for j := i+1; j < len(a); j++ { if a[j] < min { min, pos = a[j], j } }
    a[pos], a[i] = a[i], min
}

// Insertion O(n²) — sorted alt-hissəyə daxil et
for i := 1; i < len(a); i++ {
    k, j := a[i], i-1
    for j >= 0 && a[j] > k { a[j+1] = a[j]; j-- }
    a[j+1] = k
}

// Quick sort — orta O(n log n), ən pis O(n²)
func QuickSort(a []int, low, high int) {
    if low < high {
        a, j := partition(a, low, high)
        QuickSort(a, low, j-1)
        QuickSort(a, j+1, high)
    }
}

// sort.Interface: Len() / Less(i,j) / Swap(i,j) → Sort() / Stable()
sort.Ints(a)  // optimallaşdırılmış quicksort, ən pis O(n log n)
```

## Linked list

```go
type Node struct{ next *Node; value int }
type List struct{ head *Node; len int }

func (l *List) Insert(v int) {          // əvvələ — O(1)
    node := Node{value: v}
    if l.head != nil { node.next = l.head }
    l.head = &node
    l.len++
}
func (l *List) Find(v int) *Node {      // O(n)
    for t := l.head; t != nil; t = t.next {
        if t.value == v { return t }
    }
    return nil
}
// Remove: əvvəlki node-un next-ini ötürüncəyə yönəlt; node.next = nil (GC)
// container/list: PushFront/PushBack/InsertBefore/InsertAfter/Remove/Front/Back
// container/ring: New(n)/Next/Prev/Link/Do
```

## Stack / Queue / Priority queue

```go
// Stack (slice): push=append; pop=kəs; SP=-1 → boş
type Stack struct{ stack []int; sp int }
func (s *Stack) Push(v int) { s.stack = append(s.stack, v); s.sp = len(s.stack)-1 }
func (s *Stack) Pop() int {
    if s.sp == -1 { return -1 }
    v := s.stack[s.sp]
    s.stack = s.stack[:s.sp]
    s.sp--
    return v
}

// Queue (container/list): Enqueue=PushFront; Dequeue=Back+Remove
func (q *Queue) Enqueue(v int)   { q.queue.PushFront(v) }
func (q *Queue) Dequeue() int {
    if q.queue.Len() == 0 { return -1 }
    e := q.queue.Back()
    q.queue.Remove(e)
    return e.Value.(int)
}

// Priority queue: heap.Interface = sort.Interface + Push/Pop
// min-PQ: Less → pq[i].value < pq[j].value; İSTİFADƏ: heap.Push/heap.Pop
type PriorityQueue []Element
func (pq *PriorityQueue) Push(v any) { *pq = append(*pq, Element{v.(int)}) }
func (pq *PriorityQueue) Pop() any {
    old := *pq; n := len(old)-1; e := old[n]; *pq = old[:n]; return e
}
```

## Map / Hashing

```go
m := make(map[string]int)     // nil map-ə yazma OLMAZ
m["k"] = 1                    // insert/update
v := m["k"]                   // yoxdursa zero dəyər
v, ok := m["k"]               // mövcudluq
delete(m, "k")
for k, v := range m { ... }   // iterasiya

// hash: h(k)=k mod n (n primplə bağlı!); multiplication c≈0.618 (qızıl nisbət)
// collision: open addressing (linear/quadratic/double) | separate chaining
// load factor α = element / slot sayı
```

## Binary tree + traversal

```go
type Node struct{ value int; left, right, parent *Node }
// Insert(node, "left"/"right", v); Delete: yarpaq→sil; daxili→ən sol
// yarpağın dəyərini köçür + yarpağı sil

func Inorder(n *Node)  { if n != nil { Inorder(n.left); visit(n); Inorder(n.right) } }
func Preorder(n *Node) { if n != nil { visit(n); Preorder(n.left); Preorder(n.right) } }
func Postorder(n *Node){ if n != nil { Postorder(n.left); Postorder(n.right); visit(n) } }
// Level-order: queue; node dequeue → uşaqları enqueue
// inorder BST-də SIRALI nəticə verir; tree sort: array→BST→inorder O(n log n)
// Heap: parent≥child; heapsort: heap.Init + Pop → sondan yaz
```

## Qraf (set simulyasiyası ilə)

```go
type Graph struct {
    nodes map[Node]struct{}
    edges map[Edge]struct{}
}
func (g *Graph) AddEdge(u, v Node) { g.edges[Edge{u, v}] = struct{}{} }
// RemoveNode: node + incident edge-lər

// BFS (queue, səviyyəli) / DFS (rekursiv, dərin)
// MST: Prim (bir komponent genişlənir) | Kruskal (forest + PQ birləşdirir)
// Warshall: p[i][j] = p[i][j] || (p[i][k] && p[k][j])  — O(n³)
// Floyd: if d[i][j] > d[i][k]+d[k][j] { d[i][j] = ... }  — O(n³)
// Dijkstra: min proqnozlu node-u S-ə al; qalanları relax et
// Ford-Fulkerson: augmenting path tap → min residual capacity qədər axın artır
// TopSort: indegree=0 node tap → sil → təkrarla (DAG)
// Critical path: EST=irəli max(predecessors); LST=geri min(successors); L=0 → kritik
```

## Go method/interface tez referans

```go
func (r Rectangle) Area() int { return r.a * r.b }   // value receiver
func (r *Rectangle) Scale(f int) { r.a *= f }        // pointer receiver — dəyişir
type Shape interface { Area() int }                   // implicit implementasiya
```

## Komplekslik xülasəsi

| Əməliyyat | Array | List | Map |
|---|---|---|---|
| Giriş (index) | O(1) | O(n) | O(1) orta |
| Insert/remove (random) | O(n) sürüşdürmə | O(1) pointer | — |
| Axtarış (sıralmamış) | O(n) | O(n) | — |
| Axtarış (sıralı) | O(log n) | O(n) | — |
