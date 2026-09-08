# Chapters 10-13 — Scheduler, BoltDB, CLI, Next (səh. 260-404)

## Bu fəsillər nədən bəhs edir?

(10) Scheduler interfeysi + E-PVM alqoritmi (Borg-un əcdadı). (11) Store
interfeysi, in-memory → BoltDB persistent datastore. (12) Cobra CLI: worker/
manager/run/stop/status/node komandaları. (13) Nə öyrəndik, hara getməli.

## Əsas fikirlər

### 1. Scheduler interfeysi (Ch10)
```go
type Scheduler interface {
    SelectCandidateNodes(t task.Task) []*node.Node   // feasibility: resurs filtr
    Score(t task.Task, nodes []*node.Node) map[string]float64  // ballandırma
    PickNode(scores map[string]float64) string       // ən yaxşı seç
}
```
- İstənilən scheduler = bu 3 metod; K8s profiles, Nomad 4 scheduler tipi —
  eyni ideya
- RoundRobin interfeysə adaptasiya: Score hamısına 1.0 (fərqsiz), PickNode
  round-robin

### 2. E-PVM alqoritmi (Borg-un atası)
**Məqsəd:** CPU yükünü minimuma saxlayan node seçmək.

```go
// feasibility: disk yoxla
if taskDisk > node.Stats.DiskFree { skip }

// scoring: hər node üçün marginaal yük hesabla
cpuLoad := calculateLoad(cpuUsage, math.Pow(2, 0.8))  // nüvə sayına görə norma
memoryAllocated := float64(node.Stats.MemUsedKb()) + float64(node.MemoryAllocated)
memoryPercentAllocated := memoryAllocated / float64(node.Memory)
newMemPercent := (calculateLoad(memUsage, node.Memory) + memoryPercentAllocated) / 2
// 2 ölçünün ortası — köhnə node.MemoryAllocated + yeni task tələbi

// picking: minimum score
```
- **CPU istifadəsi 2 snapshotla:** /proc/stat 3 saniyə arayla —
  `total := stat2Total - stat1Total; idle := stat2Idle - stat1Idle; usage = (total-idle)/total`
- utils.HTTPWithRetry — API çağırışlarına exponential retry qatı

### 3. Store interfeysi (Ch11)
```go
type Store interface {
    Put(key string, value interface{}) error
    Get(key string) (interface{}, error)
    List() (interface{}, error)
    Count() (int, error)
    // Delete YOX — datastore = tarixi qeyd
}
```
- İn-memory → BoltDB dəyişməsi YALNIZ `New()` switch-ində ("memory"/"persistent")
- Manager: `TaskDb store.Store` + `EventDb store.Store`; worker də eyni

### 4. BoltDB
- **Embedded key-value** (server YOX) — Go-də yazılmış; fayl = DB
- Buckets (namespace) + transaction-lar:
```go
db.View(func(tx *bolt.Tx) error {        // read-only tx
    b := tx.Bucket([]byte("tasks"))
    v := b.Get([]byte(key))              // []byte key/value
    // Count: b.ForEach ilə say
})
db.Update(func(tx *bolt.Tx) error {      // yazma tx
    b := tx.Bucket([]byte(t.Bucket))
    data, _ := json.Marshal(task)         // value = JSON serialize
    return b.Put([]byte(key), data)
})
```
- Siçan yox: məcburi transaction modeli — consistency avtomatik

### 5. Cobra CLI (Ch12)
```go
var workerCmd = &cobra.Command{
    Use:   "worker",
    Short: "...",
    Run: func(cmd *cobra.Command, args []string) {
        host, _ := cmd.Flags().GetString("host")
        port, _ := cmd.Flags().GetInt("port")
        // worker.New + go RunTasks/UpdateTasks + api.Start()
    },
}
func init() {
    rootCmd.AddCommand(workerCmd)
    workerCmd.Flags().StringP("host", "H", "0.0.0.0", "...")
}
```
- `cobra-cli init` + `cobra-cli add worker` — iskelet avtomatik
- Komandalar: `worker` (start worker), `manager` (start manager), `run`
  (-f task.json → POST /tasks), `stop {taskID}` (DELETE), `status` (GET
  tasks), `node` (GET /nodes)
- Flag pattern: qısa/uzun (n/name), default, help auto-generate — docker/kubectl
  istifadəçi gözləntisi

### 6. Nə öyrəndik (Ch13)
- Manager-worker pattern: workflow sistemlərində (Airflow-tipli), integration
  sistemlərində, CI-da — orchestration konsepti konteynerdən genişdir
- Növbəti addımlar: K8s source oxu, xüsusi workflow sistemi yaz, open
  source töhfə

## Əsas terminlər

- Feasibility / Scoring / Picking
- E-PVM (Enhanced Parallel Virtual Machine — marginal cost)
- Marginal Cost (kənar yük)
- Store Interface
- BoltDB (embedded KV, buckets, transactions)
- Cobra (CLI framework)
- PersistentFlags vs Flags
- cobra-cli (scaffold)

## Praktik nəticə

- Scheduler = interfeys + 3 mərhələ; yeni alqoritm = yeni struct, manager
  dəyişməz
- Datastore = interfeys; dev-də memory, prod-da persistent — switch bir
  yerdə
- CLI = istifadəçi üçün API üzərində qat; Cobra bütün böyük Go layihələrinin
  standartı
- E-PVM sadə amma REAL (Google-un prod tarixi) — ballandırma = ölçülərin
  normallaşdırılmış ortası

## Mənbə

Pages: 260-404 (Chapters 10-13, Build an Orchestrator in Go)
