# Chapter 6 — Trees and Traversal Algorithms (səh. 279-341)

## Bu fəsil nədən bəhs edir?

İlk qeyri-xətti struktur — ağac (tree): node/edge, dərəcə, root/leaf, subtree,
parent/child, path, ancestor/descendant, level/height, növ axtarışları (full/
complete/almost complete/balanced), insert/delete əməliyyatları, binary tree
Go implementasiyası, 4 traversal alqoritmi (preorder/inorder/postorder/
level-order) və ağaclarla array sıralaması (tree sort + heap sort).

## Əsas fikirlər

### 1. Ağacın əsasları
- **Tree:** sonlu, qeyri-boş elementlər dəsti; elementlər = **node**-lar (hər
  tip data saxlaya bilər; adətən homogen)
- Qeyri-xətti struktur; node-lar **edge**-lərlə birləşir
- **Degree (dərəcə):** node-dan çıxan edge sayı; ağacın dərəcəsi = node
  dərəcələrinin maksimumu
- Xüsusi node-lar:
  - **Root:** giriş qolu olmayan node (heç bir node ondan əvvəl gəlmir)
  - **Leaf:** dərəcəsi 0 olan node-lar (çıxış qolu yoxdur)

**Nisbət terminləri:**
| Termin | Təsvir |
|---|---|
| Subtree | root xaric node-ların alt dəsti — özü də ağacdır; edge-lər root-dan subtree-lərə istiqamətlidir |
| Parent | subtree-lərin budaqlandığı node |
| Child | ortaq parent-li node-lar (subtree rootları) |
| Path | node-lar ardıcıllığı — hər növbəti əvvəlkinin child-ıdır |
| Path length | path-dakı edge sayı |
| Ancestors | root-dan həmin node-a qədər path yaratyan bütün node-lar |
| Descendants | node-un subtree-sindəki bütün node-lar |
| Level | node ilə root arasındakı edge sayı |
| Height | yarpaq level-lərinin maksimumu (ən uzun path) |

**Struktur müqayisələri:**
- **Structurally similar:** eyni node/edge sayı + eyni topologiya
- **Identical:** strukturca bənzər + müvafiq node-ların məzmunu eyni
- **Ordered tree:** hər node-un subtree-ləri sıralı dəst təşkil edir
  (uşaqlar nizamlanır — rəqəmsal/əlifba)

**Düsturlar:**
- **Internal path length:** root-dan bütün node-lara path uzunluqları cəmi
  (nᵢ = i səviyyəsindəki node sayı): I = Σ i·nᵢ
  (nümunə: 1 node @0, 3 @1, 4 @2 → I = 0+3+8)
- **Maksimum node sayı** (dərəcə m, hündürlük h; səviyyə 0: 1, 1: m, 2: m²...):
  n_max = (m^(h+1) − 1) / (m − 1)
- Yaddaşda **linked representation**
- İstifadə sahələri: hierarxik təşkilat modelləri — ailə ağaçı, şirkət
  strukturu, idman turnir bracket-i (2022 FIFA Dünya Kuboku nümunəsi)

### 2. Binary tree (ikili ağac)
- Node dərəcəsi MAKSİMUM 2; yarpaq olmayan node-un 2 və ya 1 uşağı var
  (tək uşaq sol/sağ ola bilər)

| Növ | Tərif |
|---|---|
| **Full binary tree** | bütün daxili (non-leaf) node-ların hər iki uşağı var |
| **Complete binary tree** | full + bütün yarpaqlar EYNİ səviyyədə |
| **Almost complete** | son səviyyə qismən dolu, soldan sağa ardıcıl dolur; 2^h-dən az node |
| **Balanced** | hər node üçün sol/sağ subtree hündürlük fərqi ≤ 1 |

- Complete ⇒ həmişə full; almost complete full OLA BİLƏR amma məcbur deyil

### 3. Əməliyyatlar
**Insert:**
- **Yarpaq kimi:** parent-dən yeni node-a edge əlavəsi — ASAN
- **Daxili node kimi:** parent-in edge-i indi parent ↔ yeni node-a bağlanır;
  yeni node-dan köhnə child-a yeni edge lazımdır — MÜRƏKKƏB

**Delete:**
- **Yarpaq:** node + parent-edge silinir — asan
- **Daxili node:** yarpaq silinməsinə REDUKSİYA olunur — node-un subtree-sindən
  bir meyara görə yarpaq seçilir (kitabda: ən sol yarpaq), dəyəri silinən
  node-a köçürülür, yarpaq silinir

### 4. Binary tree — Go implementasiyası
```go
type Node struct {
    value  int
    left   *Node
    right  *Node
    parent *Node   // mütləq deyil, amma delete-i asanlaşdırır
}

type BinaryTree struct {
    root *Node
}
```

**Insert(node, side, value):**
```go
func (bt *BinaryTree) Insert(node *Node, side string, value int) *Node {
    newNode := &Node{value, nil, nil, nil}
    if bt.root == nil {          // boş ağac → root
        bt.root = newNode
        return bt.root
    }
    if side == "left" {
        if node.left == nil {
            node.left = newNode          // yarpaq kimi
        } else {
            newNode.left = node.left      // daxili kimi: köhnə uşağı miras alır
            node.left = newNode
        }
    } else {
        if node.right == nil {
            node.right = newNode
        } else {
            newNode.right = node.right
            node.right = newNode
        }
    }
    newNode.parent = node
    return newNode
}
```

**DeleteLeaf:**
```go
func (n *Node) DeleteLeaf() {
    if n.parent.left == n {
        n.parent.left = nil
    } else {
        n.parent.right = nil
    }
}
```

**FindLeftmost (subtree-nin ən sol node-u):**
```go
func (n *Node) findLeftmost() *Node {
    node, next := n, n
    for next != nil {
        node = next
        if next.left != nil {
            next = next.left
        } else {
            next = next.right
        }
    }
    return node
}
```

**Delete (yarpaq → birbaşa; daxili → reduksiya):**
```go
func (bt *BinaryTree) Delete(node *Node) {
    if node.left == nil && node.right == nil {
        if node == bt.root {          // ağacda tək node
            bt.root = nil
        } else {
            node.DeleteLeaf()
        }
    } else {
        leftmostNode := node.FindLeftmost()
        node.value = leftmostNode.value
        leftmostNode.DeleteLeaf()
    }
}
```

### 5. Traversal alqoritmləri
**Traversal:** node-lara sistematik sıra ilə çıxış; hər node YALNIZ BİR
dəfə ziyarət olunur; konvensiya — sol subtree sağdan ƏVVƏL.

Nümunə ağacı: root 18, uşaqları 8 (5,9) və 25 (21,27):

| Alqoritm | Sıra | Nəticə |
|---|---|---|
| **Preorder** | root → sol → sağ | 18, 8, 5, 9, 25, 21, 27 |
| **Inorder** | sol → root → sağ | 5, 8, 9, 18, 21, 25, 27 (SIRALI!) |
| **Postorder** | sol → sağ → root | 5, 9, 8, 21, 27, 25, 18 (root SON) |
| **Level-order** | səviyyə-səviyyə, soldan sağa | 18, 8, 25, 5, 9, 21, 27 |

**Preorder (rekursiv):**
```go
func Preorder(node *Node) {
    if node != nil {
        fmt.Println(node.value)
        Preorder(node.left)
        Preorder(node.right)
    }
}
```

**Inorder (rekursiv):**
```go
func Inorder(node *Node) {
    if node != nil {
        Inorder(node.left)
        fmt.Println(node.value)
        Inorder(node.right)
    }
}
```

**Postorder (rekursiv):**
```go
func Postorder(node *Node) {
    if node != nil {
        Postorder(node.left)
        Postorder(node.right)
        fmt.Println(node.value)
    }
}
```

**Level-order (queue + loop, daha az istifadə olunur):**
```go
func Levelorder(node *Node) {
    queue := New()
    queue.Enqueue(node)
    for !queue.IsEmpty() {
        next := queue.Dequeue()
        fmt.Println(next.value)
        if next.left != nil {
            queue.Enqueue(next.left)
        }
        if next.right != nil {
            queue.Enqueue(next.right)
        }
    }
}
```
- Queue: Ch4-dəki container/list implementasiyası, element tipi `*Node`:
```go
func (q *Queue) Enqueue(node *Node) { q.queue.PushFront(node) }
func (q *Queue) Dequeue() *Node {
    if q.queue.Len() == 0 {
        return nil
    }
    element := q.queue.Back()
    q.queue.Remove(element)
    return element.Value.(*Node)
}
func (q *Queue) IsEmpty() bool { return q.queue.Len() == 0 }
```

### 6. Ağacla array sıralaması
**Tree sort — O(n log n):**
1. Sıralanmamış array elementləri binary tree-yə daxil edilir — mövqe DƏYƏR
   ilə müəyyən olunur (kiçik → sol subtree, böyük → sağ); array[0] = root;
   hər yeni node həmişə yarpaq kimi daxil olur
2. Ağac **inorder** ilə traverl edilir → sıralı çıxış

```go
func ArrayToTree(array []int) *Node {
    var bt BinaryTree
    root := bt.Insert(nil, "left", array[0])
    for i := 1; i < len(array); i++ {
        side, node := find(array[i], root)
        bt.Insert(node, side, array[i])
    }
    return root
}

func find(value int, root *Node) (side string, node *Node) {
    next := root
    for next != nil {
        node = next
        if value <= next.value {
            side = "left"
            next = next.left
        } else {
            side = "right"
            next = next.right
        }
    }
    return
}

// sıralama:
array := []int{18, 27, 5, 21, 1, 9}
root := binarytree.ArrayToTree(array)
binarytree.Inorder(root)
```

**Heap — təkmilləşmiş variant:**
- Heap = complete/almost complete binary tree, parent-in dəyəri həmişə
  child-dən BÖYÜK BƏRABƏRDIR → ən böyük dəyər ROOT-dadır
- `container/heap` ilə implementasiya:
```go
type IntHeap []int

func (h IntHeap) Len() int           { return len(h) }
func (h IntHeap) Less(i, j int) bool { return h[i] > h[j] }  // maksimum root-da!
func (h IntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *IntHeap) Push(v any)        { *h = append(*h, v.(int)) }
func (h *IntHeap) Pop() any {
    old := *h
    n := len(old)
    v := old[n-1]
    *h = old[0 : n-1]
    return v
}

// Heap sort:
func Heapsort(array *IntHeap) []int {
    heap.Init(array)
    n := array.Len()
    sortedArray := make([]int, n)
    for i := n - 1; array.Len() > 0; i-- {
        sortedArray[i] = heap.Pop(array).(int)   // kök götürülür, heap yenidən təşkil olunur
    }
    return sortedArray
}
```

## Termindirmə (AZ)
- Tree — Ağac (qeyri-xətti struktur)
- Node / Edge — Node / Qol (birləşdirici)
- Degree — Dərəcə (çıxan qolların sayı)
- Root — Kök (giriş qolu olmayan node)
- Leaf — Yarpaq (dərəcə 0)
- Subtree — Alt ağac
- Parent / Child — Valideyn / Övlad node
- Path — Yol (ardıcıl valideyn-övlad node-ları)
- Ancestors / Descendants — Əcdadlar / Nəsillər
- Level / Height — Səviyyə / Hündürlük
- Internal path length — Daxili yol uzunluğu
- Binary tree — İkili ağac (dərəcə ≤ 2)
- Full / Complete / Almost complete / Balanced — Tam / Bitmiş / Yarı-bitmiş /
  Balanslı ikili ağac
- Traversal — Traversal (sistematik ziyarət)
- Preorder / Inorder / Postorder / Level-order — ön-sıra / orta-sıra /
  son-sıra / səviyyə-sıra
- Tree sort — Ağac sıralaması (inorder çıxışı)
- Heap — Heap (parent ≥ child; max root-da)
- Heap sort — Heap sıralaması
- container/heap — heap paketi

## Kviz sualları
1. Ağacın hündürlüyü nədir? (yarpaq səviyyələrinin maksimumu — ən uzun path)
2. Complete və almost complete fərqi? (complete: bütün yarpaqlar eyni
   səviyyədə; almost: son səviyyə soldan qismən dolu)
3. Hansı traversal root-u SON ziyarət edir? (postorder)
4. Hansı traversal sıralı nəticə verir (BST üçün)? (inorder: sol→root→sağ)
5. Level-order hansı data struktur ilə implement olunur? (queue — node
   dequeued, uşaqları enqueue)
6. Daxili node-un silinməsi necə asanlaşdırılır? (ən sol yarpağın dəyəri
   silinən node-a köçürülür, yarpaq silinir)
7. Heap-də Less() niyə `>` işarəsidir (heapsort üçün)? (ən böyük element
   root-da olsun — sondan başa düzülüş)
8. Tree sort-un kompleksliyi? (O(n log n))
