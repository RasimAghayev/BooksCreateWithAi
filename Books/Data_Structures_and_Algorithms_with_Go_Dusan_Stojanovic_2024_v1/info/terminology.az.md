# DS&A with Go — Terminologiya (AZ)

## Ümumi

| Termin | Azərbaycanca | İzah |
|---|---|---|
| Data structure | Data strukturu | datanın təşkilinin təsviri |
| Linear / Non-linear | Xətti / Qeyri-xətti | 2 qonşu / çoxlu əlaqə |
| Static / Dynamic | Statik / Dinamik | ölçü dəyişmir / dəyişir |
| Homogenous / Heterogenous | Homogen / Heterogen | eyni / fərqli tipli elementlər |
| Sequential representation | Ardıcıl təmsil | fasiləsiz yaddaş; fiziki=məntiqi sıra |
| Linked representation | Zəncirvari təmsil | pointer-lə birləşən səpələnmiş elementlər |
| Pointer | Göstərici | yaddaş ünvanı saxlayan dəyişən |
| Algorithm | Alqoritm | girişi çıxışa çevirən proses |
| Invariant | İnvariant | icra boyu dəyişməyən ifadə |
| Proof of correctness | Düzgünlük sübutu | predikat hesabı ilə teorem kimi sübut |
| Turing machine | Turing maşını | hesablamanın riyazi modeli (sonsuz lent + qaydalar) |
| Pseudocode | Pseudokod | təbii dil + proqram konvensiyaları |
| Flowchart | Axın diaqramı | qrafik təmsil (rectangle=proses, diamond=qərar) |
| Big O | Böyük O | komplekslik təxmini (order of approximation) |

## Komplekslik sinifləri

| Sinif | Azərbaycanca | Nümunə |
|---|---|---|
| O(1) | Konstant | list əvvəlinə insert |
| O(log n) | Loqarifmik | binary search |
| O(n) | Xətti | sequential search |
| O(n log n) | Loqarifmik-xətti | merge sort, quicksort (orta) |
| O(n²) | Kvadratik | bubble/selection/insertion sort |
| O(n^k) | Polinomial | k iç-içə loop |
| O(k^n) | Eksponensial | exhaustive search |

## Paradiqmalar

| Termin | Azərbaycanca |
|---|---|
| Divide and conquer | Böl və hökm et |
| Dynamic programming | Dinamik proqramlaşdırma |
| Greedy | Aclıq (o anki ən yaxşı) |
| Randomized | Təsadüfi |
| Genetic | Genetik |
| Heuristic | Evristik |
| Recursive / Iterative | Rekursiv / İterativ |
| Parallel / Distributed | Paralel / Paylanmış |

## Array / Slice

| Termin | Azərbaycanca | İzah |
|---|---|---|
| Index | İndeks | element mövqeyi |
| Underlying array | Alt massiv | slice-in istinad etdiyi array |
| Length / Capacity | Uzunluq / Tutum | len / cap |
| Half-open range | Yarı-açıq aralıq | [low:high) — high çıxarılır |
| Matrix | Matris | 2D array; row-major saxlanma |
| Append | Əlavə etmə | append() — sonda |

## List

| Termin | Azərbaycanca | İzah |
|---|---|---|
| Node | Node | data + pointer |
| Head / Tail | Baş / Quyruq | ilk / son node |
| Single/double-linked list | Tək/İki-zəncirli siyahı | next / next+prev |
| Circular list (ring) | Dairəvi siyahı | son → ilk |
| Synonym (GC-də) | — | John McCarthy ~1959 |
| Orphan | Yetim | istinadsız yaddaş obyekti |
| Detach | Ayırmaq | node.next = nil |
| Concatenate | Birləşdirmək | 2 listin zəncirlənməsi |

## Stack / Queue

| Termin | Azərbaycanca |
|---|---|
| Stack pointer | Stek göstəricisi |
| Top / Bottom | Zirvə / Dip |
| LIFO / FIFO | Son gələn ilk çıxar / İlkin gələn ilk çıxar |
| Push / Pop | Yığ / Çıxar |
| Enqueue / Dequeue | Növbəyə sal / çıxar |
| Front / Rear | Baş / Quyruq (növbə ucları) |
| Deque | İkiuclu növbə |
| Priority queue | Prioritet növbəsi |
| Min / Max priority queue | Min / Maks prioritet |
| defer | Gecikdirilmiş çağırış (LIFO stack) |

## Hashing / Map

| Termin | Azərbaycanca |
|---|---|
| Hash function | Hash funksiyası (açar → index) |
| Hash table | Hash cədvəli |
| Collision | Toqquşma |
| Synonyms | Sinonimlər (eyn index) |
| Equivalence class | Ekvivalentlik sinfi |
| Load factor (α) | Yükləmə amili |
| Open addressing | Açıq ünvanlama |
| Probe sequence | Zond ardıcıllığı |
| Rehashing | Yenidən hash |
| Linear/Quadratic probing | Xətti / Kvadratik zondlama |
| Double hashing | İkili hash |
| Separate chaining | Ayrı zəncirləmə |
| Division/Multiplication/Mid-square/Digit folding/Radix/Perfect | hash üsulları |
| Digit analysis | Rəqəm təhlili |
| Golden ratio | Qızıl nisbət (c≈0.618) |
| Map | Xəritə (key-value) |
| "comma ok" | Açar mövcudluğu idiomu (v, ok := m[k]) |

## Tree

| Termin | Azərbaycanca |
|---|---|
| Root / Leaf | Kök / Yarpaq |
| Degree | Dərəcə (çıxan edge sayı) |
| Subtree | Alt ağac |
| Parent / Child | Valideyn / Övlad |
| Path / Length | Yol / uzunluğu |
| Ancestor / Descendant | Əcdad / Nəsil |
| Level / Height | Səviyyə / Hündürlük |
| Internal path length | Daxili yol uzunluğu |
| Binary tree | İkili ağac (dərəcə≤2) |
| Full / Complete / Almost complete / Balanced | Tam / Bitmiş / Yarı-bitmiş / Balanslı |
| Preorder / Inorder / Postorder / Level-order | Ön / Orta / Son / Səviyyə sıralı traversal |
| Heap | Heap (parent ≥ child) |
| Heap sort | Heap sıralaması |

## Graph

| Termin | Azərbaycanca |
|---|---|
| Graph G=(V,E) | Qraf (node + edge dəstləri) |
| Order / Size | Node sayı / Edge sayı |
| Incident | Bitişik |
| Directed / Undirected / Mixed | İstiqamətli / -siz / Qarışıq |
| Indegree / Outdegree | Giriş / Çıxış dərəcəsi |
| Loop / Parallel edges | Loop / Paralel |
| Multigraph / Simple | Multiqraf / Sadə qraf |
| Subgraph | Alt qraf |
| Weighted / Network | Çəkili / Şəbəkə |
| Path / Simple path / Cycle | Yol / Sadı / Dövr |
| Complete / Dense / Sparse | Tam / Sıx / Seyri |
| Connected / Strongly connected | Birləşmiş / Güclü |
| Tree/Forward/Back/Cross edges | traversal edge kateqoriyaları |
| BFS / DFS | Eninə / Dərinlik axtarışı |
| Spanning tree / MST | Spanning ağacı / Minimal spanning tree |
| Forest | Meşək (birləşməmiş komponentlər dəsti) |
| Adjacency matrix | Qonşuluq matrisi |
| Reachability matrix | Çatımlıq matrisi (transitive closure) |
| Cost / Distance matrix | Qiymət / Məsafə matrisi |
| Warshall / Floyd / Dijkstra | alqoritm adları |
| Eccentricity / Center | Ekssentrisitet / Mərkəz |
| Flow network | Axın şəbəkəsi |
| Source / Target | Mənbə / Hədəf |
| Capacity / Flow | Tutum / Axın |
| Residual capacity | Qalıq tutum |
| Augmenting path | Artırıcı yol |
| Ford-Fulkerson | Maksimal axın alqoritmi |
| Topological sorting | Topoloji sıralama |
| Critical path / activity | Kritik yol / aktivlik |
| EST / LST / Latency | Ən erkən / ən gec başlanğıc / gecikmə payı |
| DAG | İstiqamətli dövrsüz qraf |
