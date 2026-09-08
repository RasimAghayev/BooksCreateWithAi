# Data Structures and Algorithms with Go — Xülasə (AZ)

## Kitab kimin üçündür?

DS&A-nı Go dilində öyrənmək istəyən, Go əsaslarını bilən developer-lər üçün
nəzəri-tədris kitabı. Hər struktur/alqoritm: konsept → klassifikasiya → Go
implementasiyası → test sualları. BPB Publications (Hindistan), 1-ci nəşr
2024, 487 səh., 7 fəsil.

## Struktur xətti

| Fəsil | Mövzu | Əhatə |
|---|---|---|
| 1 | Fundamentals | DS xarakteristikaları, yaddaş təmsili, pointer, O-notasiya, alqoritm tarixi/təsnifatı, Go funksiyaları |
| 2 | Arrays | array/slice/matrix, method/interface, sequential+binary search, 4 sort alqoritmi, sort paketi |
| 3 | Lists | node/head/tail, single-linked-in tam implementasiyası, container/list + container/ring |
| 4 | Stack & Queue | LIFO/FIFO, push/pop, defer-in stack-i, enqueue/dequeue, priority queue (container/heap) |
| 5 | Hashing & Maps | hash funksiya metodları, collision həlli (probing/chaining), Go map əməliyyatları |
| 6 | Trees | binary tree növləri, insert/delete, 4 traversal, tree sort, heap sort |
| 7 | Graphs | BFS/DFS, Prim/Kruskal (MST), Warshall, Floyd/Dijkstra, Ford-Fulkerson, topological sort, critical path |

## Ən vacib 5 fikir

1. **Hər DS öz xarakteristikaları ilə təsnif olunur:** linear/non-linear,
   statik/dinamik, homogen/heterogen; seçim işin tələbindən asılıdır
   (array: sürətli index girişi; list: ucuz insert/remove)
2. **O-notasiya qərarvermə alətidir:** O(1) < O(log n) < O(n) < O(n log n)
   < O(n²) < O(k^n) — eyni problemi həll edən alqoritmlərdən ən aşağı
   kompleksliyi seç
3. **Go standart kitabxanası pərakə DS təklif edir:** map/slice daxili;
   container/list, container/ring, container/heap hazır; stack, queue,
   single-linked list, tree, graph — ÖZÜN YAZ
4. **Sort+Search hazır optimallaşdırılıb:** sort paketi (optimallaşdırılmış
   quicksort, Interface: Len/Less/Swap; Stable) və binary search funksiyaları
   — öz tipini sıralamaq üçün yalnız interface implement et
5. **Qraflar real dünyanın modelidir:** yol şəbəkələri (MST), naviqasiya
   (Dijkstra), layihə idarəsi (topological sort + critical path), ötürmə
   şəbəkələri (max flow)

## Kitabın güclü/ zəif tərəfləri

**Güclü:** tam Go kodu (hər alqoritm işləyən implementasiya ilə); kviz +
cavablar + key terms hər fəsildə; tərifə-sadiq qraf implementasiyası
(map[struct]struct{} set simulyasiyası).

**Zəif:** bəzi kodlar sətir səhvi ilə çap olunub (mətndə "rigth" kimi
səhvlər, minEdge()-də weight yox value müqayisəsi — orijinal depoda
düzəldilməli); Big-O bəzən qeyri-standart göstərilir (binary search üçün
"O(n log n)" yazılıb — düzgünü O(log n)); versiya köhnə deyil (any tipi
istifadə olunub).

## Qeydlər

- Kod bundle: https://rebrand.ly/8bb22e (GitHub-da da var)
- Kitabın ön sözü: Go-nun bəzi strukturları standart kitabxanada hazırdır,
  qalanları sıfırdan yazılacaq
- Müəllif: Dušan Stojanović — "Building Server-side and Microservices with
  Go" (2021) müəllifi, Belqrad, senior Go developer
