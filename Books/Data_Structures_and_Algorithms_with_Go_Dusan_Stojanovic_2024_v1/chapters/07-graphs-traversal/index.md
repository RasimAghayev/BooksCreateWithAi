# Chapter 7 — Graphs and Traversal Algorithms (səh. 342-487)

## Bu fəsil nədən bəhs edir?

Qraf nəzəriyyəsi: təriflər və növlər, Go implementasiyası, traversal (BFS/DFS),
spanning tree (Prim, Kruskal), transitive closure (Warshall), ən qısa yollar
(Floyd, Dijkstra), maksimal axın (Ford-Fulkerson), topoloji sıralama və
kritik yol — hamısı tam Go kodu ilə. Kitabın yekun fəsili.

## Əsas fikirlər

### 1. Qrafın əsasları
- **G = (V, E):** V — sonlu, qeyri-boş node dəsti; E — elementlər arası
  ikili münasibətlər; qeyri-xətti münasibətlərin modelləşdirilməsi üçün ideal
- **Order** = node sayı; **Size** = edge sayı
- Edge yalnız bir cüt node-u birləşdirir — həmin node-lara **incident**-dir

| Növ | Xassə |
|---|---|
| Directed (istiqamətli) | node cütləri SIRALIDIR (haradan haraya məlumdur) |
| Undirected | cütlər sıralanmamış |
| Mixed | hər ikisi var |

- **Degree** = incident edge sayı; directed: **indegree** (giriş) + **outdegree** (çıxış)
- **Loop:** eyni node-dan başlayıb yalnız ona bitişən edge
- **Parallel edges:** eyni node cütünü birləşdirən çoxlu edge (praktikada çox
  vaxt qadağandır)
- **Multigraph:** loop + parallel edge saxlayır; **Simple graph:** saxlamır
- **Subgraph** G'=(V',E'): V'⊆V, E'⊆E
- **Weighted graph:** edge-lərə çəki təyin olunub; istiqamətli weighted =
  **network**
- **Path:** (v₀, v₁, ..., vₖ) node ardıcıllığı — edge-lər ardıcıl birləşir;
  path mövcuddursa node **reachable**
- **Simple path:** bütün node-lar fərqli (ilk/son istisna ola bilər)
- **Cycle:** eyni node-dan başlayıb bitən path; cyclic/acyclic
- **Complete graph:** hər node cütü edge ilə birləşib; edge sayı =
  n(n-1)/2 (maksimum)
- **Dense / Sparse:** edge sayı maksimuma yaxın / az
- **Connected:** (undirected) hər cüt arasında path var; əks halda bir neçə
  connected component; **Strongly connected:** (directed) hər node hər
  node-dan reachable

### 2. Əməliyyatlar
- **Insert node:** node + edge-lər dəstlərə əlavə
- **Delete node:** node + BÜTÜN incident edge-lər silinir
- **Traverse:** mürəkkəb — alqoritmlər aşağıda

### 3. Qraf Go-da (tərifə sadiq implementasiya)
```go
type Node struct{ value int }
type Edge struct{ u, v Node }
type Graph struct {
    nodes map[Node]struct{}     // set simulyasiyası: map + empty struct
    edges map[Edge]struct{}
}
func New() *Graph {
    return &Graph{
        nodes: make(map[Node]struct{}),
        edges: make(map[Edge]struct{}),
    }
}
func (g *Graph) AddNode(n Node)   { g.nodes[n] = struct{}{} }
func (g *Graph) AddEdge(u, v Node) {
    g.edges[Edge{u, v}] = struct{}{}
}
func (g *Graph) RemoveNode(n Node) {
    delete(g.nodes, n)
    for e := range g.edges {
        if e.u == n || e.v == n {   // incident edge-lər
            delete(g.edges, e)
        }
    }
}
```

### 4. Traversal alqoritmləri
- Bütün node-lar mütəşəkkil sıra ilə 1 dəfə ziyarət; root YOXDUR → başlanğıc
  node seçilir; node-lar **visited** kimi işarələnir (map ilə visit vektoru)
- Traversal traversal ağacı yaradır; edge-lər 4 kateqoriyaya bölünür:
  - **Tree edges:** parent-child (child bununla ziyarət olunub) — yalnız bunlar
    traversal ağacına daxildir
  - **Forward:** artıq ziyarət olunmuş descendant-ə
  - **Back:** artıq ziyarət olunmuş ancestor-a
  - **Cross:** ancestor-descendant münasibətində olmayan 2 node

**Breadth-first search (BFS):** başlanğıc → qonşular → onların ziyarət
olunmamış qonşuları... səviyyə-səviyyə. **Queue** ilə:
```go
func BFS(g *Graph, start *Node) {
    visit := make(map[int]bool)
    for n := range g.nodes {
        visit[n.value] = false
    }
    visit[start.value] = true
    queue := NewQueue()
    queue.Enqueue(start)
    for !queue.IsEmpty() {
        u := queue.Dequeue()
        fmt.Println(u.value)
        for edge := range g.edges {
            if edge.u.value == u.value && !visit[edge.v.value] {
                visit[edge.v.value] = true
                n := edge.v
                queue.Enqueue(&n)
            }
        }
    }
}
```

**Depth-first search (DFS):** bir istiqamətdə mümkün qədər dərin; yox çıxanda
əvvəlki node-a qayıdıb yeni path. **Rekursiya** ilə:
```go
func DFS(g *Graph, start *Node) {
    visit := make(map[int]bool)
    for n := range g.nodes {
        visit[n.value] = false
    }
    dfsVisit(g, start, visit)
}
func dfsVisit(g *Graph, u *Node, visit map[int]bool) {
    visit[u.value] = true
    fmt.Println(u.value)
    for edge := range g.edges {
        if edge.u.value == u.value && !visit[edge.v.value] {
            dfsVisit(g, &edge.v, visit)
        }
    }
}
```

### 5. Spanning tree
- Birləşmiş undirected G üçün: ST bütün node-ları saxlayır (U=V), dövrsüz
  birləşmə üçün lazım olan minimum edge sayı ilə — UNİKAL DEYİL (traversal
  alqoritmləri ilə yaradıla bilər)
- **Cost:** tree edge-lərinin çəkiləri cəmi; **minimal-cost spanning tree
  (MST):** ən ucuz spanning tree (şəhər/yol şəbəkəsi nümunəsi)

**Weighted graph implementasiyası:**
```go
type WeightedEdge struct{ u, v Node; weight int }
type WeightedGraph struct {
    nodes map[Node]struct{}
    edges map[WeightedEdge]struct{}
}
type MST struct {
    nodes map[Node]struct{}
    edges map[WeightedEdge]struct{}
}
```

**Prim (greedy; əsl müəllif Vojtech Jarnik — Prim sonradan yenidən nəşr edib):**
connected component-i artırmaqla; hər addımda minimal çəkili edge:
```go
func Prim(wg *WeightedGraph, start *Node) *MST {
    treeNodes := make(map[Node]struct{})
    treeEdges := make(map[WeightedEdge]struct{})
    treeNodes[*start] = struct{}{}
    for len(treeNodes) != len(wg.nodes) {
        node, minEdge := minEdge(wg, treeNodes)
        treeNodes[node] = struct{}{}
        treeEdges[minEdge] = struct{}{}
    }
    return &MST{treeNodes, treeEdges}
}
// minEdge: bir ucu tree-də, digəri xaricində olan minimal edge
```

**Kruskal (greedy; forest-dən başlayır):** ən kiçik çəkili edge hər addımda
2 ayrı komponenti birləşdirirsə əlavə olunur. **Priority queue** (weight
prioritet) + forest:
```go
func Kruskal(wg *WeightedGraph) *MST {
    treeEdges := make(map[WeightedEdge]struct{})
    pq := make(PriorityQueue, 0)
    for edge := range wg.edges {
        heap.Push(&pq, edge)
    }
    forest := make(map[int][]Node)   // hər node ayrı komponent
    // ...
    // edge pop → findInForest i≠j → birləşdir (append + delete), n++
    // MST edge sayı = node sayı - 1
}
func findInForest(forest map[int][]Node, node int) int { ... }
```

### 6. Transitive closure
- **Reachability matrix P:** p[i,j]=1 path varsa; =0 yoxdursa → qrafın
  **transitive closure**-u
- Giriş: **adjacency matrix** (a[i,j]=1 qonşuluqda; node-lar 0-dan nömrələnir)
- **Warshall alqoritmi** (Stephen Warshall): adjacency-dən başlayır; k
  vasitəsilə i→j path yoxlanılır — O(n³):
```go
func Warshall(a [][]bool) (p [][]bool) {
    p = a
    for k := 0; k < len(p); k++ {
        for i := 0; i < len(p); i++ {
            for j := 0; j < len(p); j++ {
                p[i][j] = p[i][j] || (p[i][k] && p[k][j])
            }
        }
    }
    return
}
```
- Undirected üçün matris simmetrikdir; cycle varsa diaqonaldal elementlər 1

### 7. Shortest paths
- Path weight = edge çəkiləri cəmi; **d(i,j) = min{w(p)}**; unreachable → ∞
  (bəzi alqoritmlər mənfi çəkini dəstəkləmir)

**Floyd (Robert Floyd) — bütün cütlər üçün, relaxation tipli:**
- Giriş: **cost matrix W** (diaqonal 0; edge yoxdursa ∞); Çıxış: **distance
  matrix D**
- Hər an yuxarı hədd saxlanılır; k inter-node-u vasitəsilə yoxlama:
```go
func Floyd(w [][]int) (d [][]int) {
    d = w
    for k := 0; k < len(d); k++ {
        for i := 0; i < len(d); i++ {
            for j := 0; j < len(d); j++ {
                if d[i][j] > d[i][k]+d[k][j] {
                    d[i][j] = d[i][k] + d[k][j]
                }
            }
        }
    }
    return
}
// INF = 99999 konstantası ilə modelləşdirilir
```
- **Eccentricity:** node-un digərlərinə maksimum məsafə; **Center:** minimum
  eccentricity-li node

**Dijkstra (Edsger Dijkstra) — bir node-dan hamısına; mənfi çəki YOX:**
- **S:** qısa məsafəsi TAPILMIŞ node-lar; **V:** qalanlar (carxi proqnozlar);
  **d vektoru:** carxi qiymətlər → yekunda qısa yollar:
```go
func Dijkstra(start int, w [][]int) (d map[int]int) {
    n := len(w)
    s := make(map[int]struct{})
    s[start] = struct{}{}
    v := make(map[int]struct{})
    d = make(map[int]int)
    for i := 0; i < n; i++ {
        if i != start {
            v[i] = struct{}{}
            d[i] = w[start][i]
        }
    }
    for len(v) != 0 {
        i := findMin(d, v)     // V-də ən kiçik proqnoz
        s[i] = struct{}{}
        delete(v, i)
        for j := range v {
            if d[i]+w[i][j] < d[j] {   // relaxation
                d[j] = d[i] + w[i][j]
            }
        }
    }
    return
}
```

### 8. Flow in graphs
- **Flow network:** hər edge-də qeyri-mənfi çəki = **capacity**; **source**
  (giriş yoxdur) + **target** (çıxış yoxdur); su boruları nümunəsi
- **f(u,v)** axın funksiyası 3 məhdudiyyətlə:
  1. **Capacity:** f(u,v) ≤ c(u,v)
  2. **Symmetry:** f(u,v) = −f(v,u)
  3. **Conservation:** daxil = xaric (source/target xaric) — ümumi node axını 0
- **Residual network:** tam istifadə olunmamış edge-lər; **residual capacity:**
  c_f(u,v) = c(u,v) − f(u,v)
- **Augmenting path:** source-dan target-a hamısı müsbət residual
  capacity-li path; path axını = ən kiçik residual capacity

**Ford-Fulkerson (Ford & Fulkerson):** sıfır axından başlayır; augmenting
path tapılır, axın artırılır; path qurtaranda maksimum flow:
```go
type FlowEdge struct {
    u, v     Node
    capacity int
    flow     int
}
func FordFulkerson(graph *FlowGraph, source, target *Node) (maxFlow int) {
    // axınları sıfırla
    ok, path := AugmentingPath(graph, source, target)
    for ok {
        cf := pathFlow(path)             // path-in minimal residual capacity-si
        for _, edge := range path {
            delete(graph.edges, edge)     // map key-də flow dəyişə bilməz
            edge.flow += cf
            graph.edges[edge] = struct{}{}
        }
        maxFlow += cf
        ok, path = AugmentingPath(graph, source, target)
    }
    return
}
// AugmentingPath: BFS modifikasiyası — target çatanda path qaytarır
```

### 9. Topological sorting
- Layihə modelləşdirmə: node-lar = hadisələr, edge-lər = aktivliklər;
  asılılıq sırası = topoloji sıra
- **DAG üçün:** bütün node-ların xətti düzülüşü — hər (u,v) edge üçün u,
  v-dən ƏVVƏL gəlir
- Alqoritm (təkrarlanan): indegree=0 node tap → T-ə əlavə et → node-u
  çıxan edge-lərlə birlikdə sil; acyclic ⇒ belə node HƏMİŞƏ var; **sıra
  unikal olmaya bilər**
```go
func TopSort(graph *Graph) []Node {
    // hər iterasiyada findZeroIndegreeNode + delete(node + çıxan edge-lər)
}
func findZeroIndegreeNode(...) Node {
    // heç bir edge-in hədəfi olmayan node
}
```

### 10. Critical path
- Layihə qrafı: edge çəkisi = aktivlik müddəti; **source** (indegree 0) =
  layihə başlanğıcı, **target** (outdegree 0) = bitiş
- Layihə müddəti = source→target ƏN UZUN path (**critical path**); onun
  aktivlikləri = **critical activities** (gecikmə olmaz!)
- Hər node üçün:
  - **EST (earliest start time):** EST[j] = max{EST[i]+w(i,j)}, i∈P(j) —
    predecessors üzrə
  - **LST (latest start time):** LST[j] = min{LST[i]−w(j,i)}, i∈S(j) —
    successors üzrə
  - **Latency:** L[i] = LST[i] − EST[i]; **L=0 olan node-lar kritik yoldadır**
- Edge gecikmə payı: l(i,j) = LST[j] − EST[i] − w(i,j)
```go
func CriticalPath(graph WeightedGraph) ([]int, []int, []int) {
    t := TopSortWG(graph)          // topo sıra (map-larin copy-si ilə!)
    est[0] = 0; ileri findMax      // predecessors max
    lst[n-1] = est[n-1]; geri findMinLST  // successors min
    l[i] = lst[i] - est[i]
    return est, lst, l
}
```

### 11. Kitabın yekunu — alqoritmlər xülasəsi (Table 7.3)
| Alqoritm | Komplekslik | İstifadə |
|---|---|---|
| BFS / DFS | O(n+e) | traversal, spanning tree |
| Prim / Kruskal | O(n²) / O(e·log n) | MST |
| Warshall | O(n³) | transitive closure |
| Floyd | O(n³) | bütün cütlər qısa yol |
| Dijkstra | O(n²) | bir mənbədən qısa yol (mənfi YOX) |
| Ford-Fulkerson | O(e·maxflow) | maksimal axın |
| Topological sort | O(n²) (bu impl.) | asılılıq sırası |

- Web servis praktikasında ən çox array/slice/map; tree və graph — simulyasiya,
  AI kimi mürəkkəb problemlərdə
- Nümunələr %100 optimal deyil — orijinal ideyaya yaxın təqdim olunub;
  təkmilləşdirmək azaddır

## Termindirmə (AZ)
- Graph (V, E) — Qraf (node dəsti + edge dəsti)
- Order / Size — Qrafın dərəcəsi (node sayı) / Ölçüsü (edge sayı)
- Incident — Bitişik (edge node-a bağlıdır)
- Directed / Undirected / Mixed — İstiqamətli / İstiqamətsiz / Qarışıq
- Indegree / Outdegree — Giriş / Çıxış dərəcəsi
- Loop / Parallel edges — Loop / Paralel edge-lər
- Multigraph / Simple graph — Multiqraf / Sadə qraf
- Subgraph — Alt qraf
- Weighted graph / Network — Çəkili qraf / Şəbəkə
- Path / Simple path / Cycle — Yol / Sadə yol / Dövr
- Complete / Dense / Sparse — Tam / Sıx / Seyri qraf
- Connected / Strongly connected — Birləşmiş / Güclü birləşmiş
- Visit vector — Ziyarət vektoru (visited map)
- Tree / Forward / Back / Cross edges — traversal edge kateqoriyaları
- BFS / DFS — Eninə / Dərinlik axtarışı
- Spanning tree / MST — Spanning ağacı / Minimal spanning tree
- Prim / Kruskal alqoritmləri — MST alqoritmləri (greedy)
- Forest — Meşək (birləşməmiş komponentlər)
- Adjacency / Reachability / Cost / Distance matrix — Qonşuluq / Çatımlıq /
  Qiymət / Məsafə matrisi
- Warshall / Floyd / Dijkstra — transitive closure / bütün cütlər / tək
  mənbə qısa yol alqoritmləri
- Flow network — Axın şəbəkəsi
- Source / Target — Mənbə / Hədəf node
- Capacity / Flow — Tutum / Axın
- Residual capacity / Residual network / Augmenting path — Qalıq tutum /
  qalıq şəbəkə / artırıcı yol
- Ford-Fulkerson — maksimal axın alqoritmi
- Topological sorting — Topoloji sıralama (DAG xətti düzülüş)
- Critical path / activity — Kritik yol / aktivlik
- EST / LST / Latency — Ən erkən / Ən gec başlanğıc vaxtı / Gecikmə payı
- Eccentricity / Center — Ekssentrisitet / Mərkəz

## Kviz sualları
1. Qrafın order və size nədir? (order = node sayı; size = edge sayı)
2. Complete qrafda neçə edge var? (n(n-1)/2)
3. Traversal edge kateqoriyaları? (tree, forward, back, cross — 4)
4. BFS hansı struktur istifadə edir, DFS hansı? (BFS: queue; DFS: rekursiya/stack)
5. Prim ilə Kruskalın fərqi? (Prim: bir komponenti artırır; Kruskal:
   forest-dən birləşdirir, PQ ilə)
6. Warshall-ın girişi/çıxışı? (adjacency → reachability matrix)
7. Dijkstra mənfi çəkiləri dəstəkləyirmi? (xeyr)
8. Ford-Fulkerson nə vaxta qədər davam edir? (augmenting path qurtarana qədər)
9. Topoloji sıra hansı qrafda mövcuddur? (yalnız DAG-da — cycle olarsa
   indegree=0 node tapılmaz)
10. Kritik yol necə tapılır? (L=0 node-lar: EST=max predecessor, LST=min
    successor; fərqi sıfır olanlar)
